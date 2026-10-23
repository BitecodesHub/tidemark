// SPDX-License-Identifier: GPL-2.0-only OR BSD-2-Clause
//
// AgentOS Tidemark - Phase 0 kernel spike: read-time label probe.
//
// STATUS: Phase 0 measurement instrument. First compiled, verified and loaded
// on 2026-10-18 inside a Fedora 45 Beta aarch64 VM (kernel 7.2.0-61.fc45,
// clang 23.1.3, bpftool 7.7.0); all spike scenarios passed in three identical
// runs (12 then, 16 from 2026-10-23). Not yet run on x86_64. Written against
// torvalds/linux tag
// v7.2 sources (see spike/README.md for the facts and citations). It is a
// measurement instrument for the Phase 0 kill gate, not production code.
//
// What this program does (observe only, it never denies an operation):
//   * On file_open, file_receive and mmap_file for any process that belongs to
//     a registered session (identified by its ancestor cgroup id at a fixed
//     depth), it computes the file's label and ORs it into a per-session label
//     word held in a BPF_F_MMAPABLE array. The label word never falls.
//   * A file's label is the UNION of its security.bpf.tidemark xattr label and
//     the labels of EVERY directory rule (labelled root) that prefix-matches the
//     file's resolved path, as docs/adr/2026-10-16_Labeller-BPF-LSM_v1.md and
//     spec/clearance-model.md define it. Both sources are evaluated on every
//     event.
//   * Errors fail closed: any lookup error other than "no such xattr" or
//     "xattrs not supported" sets TM_LABEL_UNKNOWN, which consumers must treat
//     as at least SECRET.
//   * An inode-storage cache keyed by the real data inode (bpf_real_data_inode)
//     holds XATTR-derived bits only (rule-derived bits are never cached, per the
//     labeller ADR) and acts as a FLOOR, never as an answer: the xattr and the
//     directory rules are evaluated on every event and the cached bits are ORed
//     on top. The cache therefore cannot hide a raised xattr or a stricter rule.
//     It exists to observe whether an overlay open and a direct open of the same
//     lower file share one data inode. It never lowers a label; invalidation on
//     inode_post_setxattr and inode_post_removexattr is Phase 1 work (concept,
//     mechanism step 4). The cache is compiled in only when the loader reports,
//     from kernel BTF, that bpf_real_data_inode exists (have_real_data_inode
//     gate).
//   * A per-slot cause bitmask records which hook and which mechanism produced
//     each non-public label, so scenarios can assert the mechanism under test
//     rather than only the final label.
//
// The programs always return 0. Returning non-zero from an LSM hook would deny
// the operation; the Tidemark label engine only raises a label, it does not
// block reads. Blocking is done later at the exits (out of scope for Phase 0).

// Keep this file's own kfunc declarations authoritative. bpftool 7.x emits
// kfunc prototypes into vmlinux.h unless this macro is defined, and those
// prototypes would otherwise conflict with (or silently replace) the __weak
// declaration of bpf_real_data_inode below on 7.0, 7.1 and 7.2 headers.
#define BPF_NO_KFUNC_PROTOTYPES
#include "vmlinux.h"
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_tracing.h>
#include <bpf/bpf_core_read.h>

// LICENCE NOTE (decided 2026-10-23, docs/adr/2026-10-23_Phase0-Founder-Decisions_v1.md
// FD-3): the kernel only lets GPL-compatible objects call the BPF LSM kfuncs used
// here, so a kernel-facing object must declare a GPLv2-compatible licence tag.
// Kernel-facing BPF C sources are therefore dual-licensed
// (GPL-2.0-only OR BSD-2-Clause) with the tag "Dual BSD/GPL", which the kernel's
// licence check accepts as GPL-compatible; the rest of the repository stays
// Apache-2.0 and the exception is recorded in NOTICE.
char LICENSE[] SEC("license") = "Dual BSD/GPL";

// ---------------------------------------------------------------------------
// Label bit layout (spike only; the production lattice lives in pkg/labels).
// ---------------------------------------------------------------------------
#define TM_LABEL_PUBLIC   0x0ULL
#define TM_LABEL_PRIVATE  0x1ULL  // directory rule or xattr 'p'
#define TM_LABEL_SECRET   0x2ULL  // directory rule or xattr 's'
#define TM_LABEL_UNKNOWN  0x4ULL  // fail-closed: label could not be determined
#define TM_LABEL_CACHEABLE (TM_LABEL_PRIVATE | TM_LABEL_SECRET)

// Cause bits, ORed into the per-slot causes word. Mirrored in
// spike/cmd/tidemark-probe/load.go.
#define TM_CAUSE_HOOK_OPEN    0x001ULL // raised from lsm.s/file_open
#define TM_CAUSE_HOOK_RECEIVE 0x002ULL // raised from lsm.s/file_receive
#define TM_CAUSE_HOOK_MMAP    0x004ULL // raised from lsm.s/mmap_file
#define TM_CAUSE_XATTR        0x010ULL // security.bpf.tidemark xattr
#define TM_CAUSE_DIR          0x020ULL // directory rule (bpf_path_d_path)
#define TM_CAUSE_CACHE        0x040ULL // inode-storage cache added bits
#define TM_CAUSE_ERROR        0x080ULL // a lookup failed; TM_LABEL_UNKNOWN set
#define TM_CAUSE_BACKING      0x100ULL // file->f_mode has FMODE_BACKING
#define TM_CAUSE_EACCES       0x200ULL // xattr read refused with -EACCES (fail-closed)
#define TM_CAUSE_NOBUF        0x400ULL // no per-task path buffer (fail-closed)

#define TM_MAX_SLOTS       64
#define TM_MAX_ROOTS        8
#define TM_ROOT_LEN       128
#define TM_PATH_MAX      4096 // PATH_MAX; d_path never returns more
#define TM_XATTR_LEN        8

#define TM_XATTR_NAME "security.bpf.tidemark"

// errno values used by the fs kfuncs. Both spike targets (x86_64, aarch64) use
// the asm-generic numbering. vmlinux.h carries no errno macros.
#define TM_EACCES      13
#define TM_ENODATA     61
#define TM_EOPNOTSUPP  95

// File type bits from include/uapi/linux/stat.h (macros, so not in BTF).
#define TM_S_IFMT   0170000
#define TM_S_IFREG  0100000

// FMODE_BACKING is a macro, so it is not in BTF. Value from
// include/linux/fs.h at tag v7.2 (line 169, read 2026-10-17). Used only to
// record an observation bit; a wrong value cannot change a label.
#define TM_FMODE_BACKING (1U << 24)

// Fallback definition in case the generated vmlinux.h does not expose the flag.
#ifndef BPF_LOCAL_STORAGE_GET_F_CREATE
#define BPF_LOCAL_STORAGE_GET_F_CREATE (1ULL << 0)
#endif

// Return values of the label sources.
#define TM_SRC_NONE  0 // source has no opinion (no xattr, no rule matched)
#define TM_SRC_HIT   1 // source produced a label (possibly public)
#define TM_SRC_ERROR 2 // source failed; caller fails closed
#define TM_SRC_EACCES 3 // xattr read refused (-EACCES); caller fails closed

// ---------------------------------------------------------------------------
// Configuration, set by the user-space loader before load.
// const volatile so the verifier treats them as known constants and can prune
// branches (for example the bpf_real_data_inode branch when the kfunc is
// absent). See spike/cmd/tidemark-probe/load.go.
// ---------------------------------------------------------------------------

// Ancestor cgroup level at which a session is identified. For the spike a
// session is a transient scope under tmspike.slice (dash-free on purpose), so
// the scope sits at level 2 (root=0, slice=1, scope=2).
const volatile __u32 session_level = 2;

// Set to 0 by the loader when kernel BTF does not expose bpf_real_data_inode,
// which removes the overlay data-inode cache path from the verified program.
const volatile __u8 have_real_data_inode = 1;

// ---------------------------------------------------------------------------
// Maps.
// ---------------------------------------------------------------------------

// Registered sessions: ancestor cgroup id -> slot index into the label array.
struct {
	__uint(type, BPF_MAP_TYPE_HASH);
	__uint(max_entries, TM_MAX_SLOTS);
	__type(key, __u64);   // cgroup id
	__type(value, __u32); // slot
} sessions SEC(".maps");

// Per-session label word. BPF_F_MMAPABLE so the loader (and, in production, the
// gateway and broker) can hold a read-only mmap of it. Index is the slot.
struct {
	__uint(type, BPF_MAP_TYPE_ARRAY);
	__uint(max_entries, TM_MAX_SLOTS);
	__type(key, __u32);
	__type(value, __u64);
	__uint(map_flags, BPF_F_MMAPABLE);
} labels SEC(".maps");

// Per-session cause bitmask (TM_CAUSE_*), for test assertions only.
struct {
	__uint(type, BPF_MAP_TYPE_ARRAY);
	__uint(max_entries, TM_MAX_SLOTS);
	__type(key, __u32);
	__type(value, __u64);
} causes SEC(".maps");

// Labelled roots for directory rules. Each entry is a NUL-padded absolute path
// prefix of length len, without a trailing slash; entries with len == 0 are
// empty slots and are skipped.
struct tm_root {
	char prefix[TM_ROOT_LEN];
	__u64 label;
	__u32 len;
};

struct {
	__uint(type, BPF_MAP_TYPE_ARRAY);
	__uint(max_entries, TM_MAX_ROOTS);
	__type(key, __u32);
	__type(value, struct tm_root);
} roots SEC(".maps");

// Inode storage label cache, keyed by the file's real data inode. Only used
// when have_real_data_inode is set. local_storage requires BPF_F_NO_PREALLOC.
struct {
	__uint(type, BPF_MAP_TYPE_INODE_STORAGE);
	__uint(map_flags, BPF_F_NO_PREALLOC);
	__type(key, int);
	__type(value, __u64);
} inode_cache SEC(".maps");

// Per-task scratch space: a PATH_MAX buffer for bpf_path_d_path and a small
// buffer for the xattr value. Task storage is used instead of a per-CPU array
// because sleepable programs run with preemption enabled: two tasks on one CPU
// could interleave and corrupt a shared per-CPU buffer. Task storage is private
// to the current task and is freed when the task exits. The xattr buffer lives
// here, not on the stack, because at v7.2 bpf_dynptr_from_mem refuses a stack
// pointer ("Unsupported reg type fp for bpf_dynptr_from_mem data", observed on
// 7.2.0-61.fc45 on 2026-10-18).
struct tm_pathbuf {
	char p[TM_PATH_MAX];
	__u8 x[TM_XATTR_LEN];
};

struct {
	__uint(type, BPF_MAP_TYPE_TASK_STORAGE);
	__uint(map_flags, BPF_F_NO_PREALLOC);
	__type(key, int);
	__type(value, struct tm_pathbuf);
} pathbuf SEC(".maps");

// Ring buffer of label-change events, emitted only when a label actually rises.
struct tm_event {
	__u64 cgroup_id;
	__u64 old_label;
	__u64 new_label;
	__u64 cause; // TM_CAUSE_* bits for this event
	__u32 slot;
	char comm[16];
};

struct {
	__uint(type, BPF_MAP_TYPE_RINGBUF);
	__uint(max_entries, 1 << 16);
} events SEC(".maps");

// ---------------------------------------------------------------------------
// kfunc declarations. bpf_real_data_inode is declared __weak so the object
// still loads on kernels that lack it; the have_real_data_inode gate keeps the
// call out of the verified path when it is absent.
// Signatures taken verbatim from fs/bpf_fs_kfuncs.c at tag v7.2 (lines 79,
// 160 and 377).
// ---------------------------------------------------------------------------
extern int bpf_get_file_xattr(struct file *file, const char *name__str,
			      struct bpf_dynptr *value_p) __ksym;
extern int bpf_path_d_path(const struct path *path, char *buf,
			   size_t buf__sz) __ksym;
extern struct inode *bpf_real_data_inode(struct file *file) __ksym __weak;

// ---------------------------------------------------------------------------
// Helpers.
// ---------------------------------------------------------------------------

// Resolve the session this task belongs to. Returns 1 and fills *slot on hit.
static __always_inline int session_slot(__u32 *slot)
{
	__u64 cid = bpf_get_current_ancestor_cgroup_id(session_level);
	if (cid == 0)
		return 0;
	__u32 *s = bpf_map_lookup_elem(&sessions, &cid);
	if (!s)
		return 0;
	*slot = *s;
	return 1;
}

// OR cause bits into the slot's cause word.
static __always_inline void record_cause(__u32 slot, __u64 cause)
{
	__u64 *c = bpf_map_lookup_elem(&causes, &slot);
	if (c)
		__sync_fetch_and_or(c, cause);
}

// OR new bits into the session's label word, emit an event if it actually rose.
static __always_inline void raise_label(__u32 slot, __u64 add, __u64 cause)
{
	if (add == 0)
		return;
	__u64 *word = bpf_map_lookup_elem(&labels, &slot);
	if (!word)
		return;

	// The fetch form needs BPF ISA v3 atomics; the Makefile passes -mcpu=v3.
	__u64 old = __sync_fetch_and_or(word, add);
	__u64 now = old | add;
	if (now == old)
		return; // no change, stay quiet

	struct tm_event *e = bpf_ringbuf_reserve(&events, sizeof(*e), 0);
	if (!e)
		return;
	e->cgroup_id = bpf_get_current_ancestor_cgroup_id(session_level);
	e->old_label = old;
	e->new_label = now;
	e->cause = cause;
	e->slot = slot;
	bpf_get_current_comm(&e->comm, sizeof(e->comm));
	bpf_ringbuf_submit(e, 0);
}

// Fetch (creating on first use) the current task's scratch space. Creating the
// 4 KiB element can fail transiently (the v7.2 allocation path uses a
// trylock), so retry once; a second failure is reported with its own cause bit
// so it cannot be mistaken for a lookup error.
static __always_inline struct tm_pathbuf *get_scratch(void)
{
	struct task_struct *task = bpf_get_current_task_btf();
	struct tm_pathbuf *pb;

	pb = bpf_task_storage_get(&pathbuf, task, NULL,
				  BPF_LOCAL_STORAGE_GET_F_CREATE);
	if (!pb)
		pb = bpf_task_storage_get(&pathbuf, task, NULL,
					  BPF_LOCAL_STORAGE_GET_F_CREATE);
	return pb;
}

// Read the security.bpf.tidemark xattr.
// Returns TM_SRC_NONE when the file has no such xattr (-ENODATA) or the
// filesystem has no xattr support (-EOPNOTSUPP); TM_SRC_HIT with *label set
// for a recognised value; TM_SRC_ERROR for every other outcome. Note that at
// v7.2 bpf_get_file_xattr applies inode_permission(MAY_READ) with the CURRENT
// task's credentials (fs/bpf_fs_kfuncs.c, bpf_xattr_read_permission), so a
// session principal that receives or maps a file it cannot read gets -EACCES
// here; that case returns TM_SRC_EACCES and fails closed. One consequence is
// over-labelling of execute-only binaries (for example a mode 4111 sudo): an
// exec open of such a file sets TM_LABEL_UNKNOWN. -ERANGE (value longer than 8
// bytes) also fails closed.
static __always_inline int label_from_xattr(struct file *file,
					    struct tm_pathbuf *pb, __u64 *label)
{
	struct bpf_dynptr ptr;

	__builtin_memset(pb->x, 0, sizeof(pb->x));
	if (bpf_dynptr_from_mem(pb->x, sizeof(pb->x), 0, &ptr))
		return TM_SRC_ERROR;

	int n = bpf_get_file_xattr(file, TM_XATTR_NAME, &ptr);
	if (n == -TM_ENODATA || n == -TM_EOPNOTSUPP)
		return TM_SRC_NONE;
	if (n == -TM_EACCES)
		return TM_SRC_EACCES;
	if (n != 1)
		return TM_SRC_ERROR; // other errors, an empty value, or extra bytes

	// The spike stores the label as exactly one ASCII byte: 'p' private,
	// 's' secret, '-' explicitly public. Any other value fails closed.
	// Production uses a richer encoding.
	switch (pb->x[0]) {
	case 's':
		*label = TM_LABEL_SECRET;
		return TM_SRC_HIT;
	case 'p':
		*label = TM_LABEL_PRIVATE;
		return TM_SRC_HIT;
	case '-':
		*label = TM_LABEL_PUBLIC;
		return TM_SRC_HIT;
	default:
		return TM_SRC_ERROR;
	}
}

// Match the file's resolved path against the directory rules.
// Returns TM_SRC_HIT with *label set to the union of the labels of EVERY
// matching root, TM_SRC_NONE when no root matches, and TM_SRC_ERROR when the
// path cannot be resolved (for example -ENAMETOOLONG).
//
// A root matches only on a path-component boundary: the byte after the prefix
// must be '/' or the terminating NUL, so a root /a/Private does not match
// /a/PrivateX/file.
//
// The path is resolved by d_path against current->fs->root, so the result
// depends on the caller's mount namespace and chroot, and it names the dentry
// that was opened: a hard link or a rename out of a labelled root yields a path
// no rule matches (scenarios 11 and 12 record these gaps). See spike/README.md.
static __always_inline int label_from_dir(struct file *file,
					  struct tm_pathbuf *pb, __u64 *label)
{
	// plen includes the terminating NUL (fs/bpf_fs_kfuncs.c, v7.2 line 79).
	int plen = bpf_path_d_path(&file->f_path, pb->p, sizeof(pb->p));
	if (plen <= 0)
		return TM_SRC_ERROR;

	// Bounded loop over every root slot; empty slots (len == 0) are skipped.
	// Every matching root contributes its label (union), so overlapping rules
	// such as ~/Private (private) and ~/Private/keys (secret) both apply.
	int hit = 0;
	__u64 acc = TM_LABEL_PUBLIC;
	for (__u32 i = 0; i < TM_MAX_ROOTS; i++) {
		__u32 key = i;
		struct tm_root *r = bpf_map_lookup_elem(&roots, &key);
		if (!r)
			continue;
		__u32 rlen = r->len;
		// rlen < plen guarantees pb->p[rlen] is inside the string (it is
		// at most the NUL terminator).
		if (rlen == 0 || rlen >= TM_ROOT_LEN || (int)rlen >= plen)
			continue;

		// Manual bounded prefix compare; the verifier needs the bounds.
		int match = 1;
		for (__u32 j = 0; j < TM_ROOT_LEN; j++) {
			if (j >= rlen)
				break;
			if (pb->p[j & (TM_PATH_MAX - 1)] !=
			    r->prefix[j & (TM_ROOT_LEN - 1)]) {
				match = 0;
				break;
			}
		}
		if (!match)
			continue;

		char next = pb->p[rlen & (TM_PATH_MAX - 1)];
		if (next != '/' && next != '\0')
			continue; // not a path-component boundary

		acc |= r->label;
		hit = 1;
	}
	if (!hit)
		return TM_SRC_NONE;
	*label = acc;
	return TM_SRC_HIT;
}

// Compute and apply the label for a file reference in the current session.
// hook is one of TM_CAUSE_HOOK_*.
//
// file is never NULL here: the file_open and file_receive hooks always pass a
// file, and tm_mmap_file filters anonymous mappings before calling this. Do
// NOT add an "if (!file)" check here as a safety net: the verifier types LSM
// arguments as trusted, non-NULL pointers and removes such a check as dead
// code (see tm_mmap_file).
static __always_inline void handle_file(struct file *file, __u64 hook)
{
	__u32 slot;
	if (!session_slot(&slot))
		return;

	// Only regular files carry data labels. Skipping everything else avoids
	// fail-closed over-labelling from pipes, sockets and anon-inode files
	// (eventfd, epoll, bpf maps) that another uid created, whose xattr read
	// would return -EACCES.
	struct inode *ino = file->f_inode;
	if (!ino || (ino->i_mode & TM_S_IFMT) != TM_S_IFREG)
		return;

	__u64 label = TM_LABEL_PUBLIC;
	__u64 cause = 0;
	int backing = (file->f_mode & TM_FMODE_BACKING) != 0;

	struct tm_pathbuf *pb = get_scratch();
	if (!pb) {
		// No scratch space: neither source can be evaluated. Fail closed.
		label |= TM_LABEL_UNKNOWN;
		cause |= TM_CAUSE_NOBUF | hook;
		record_cause(slot, cause);
		raise_label(slot, label, cause);
		return;
	}

	// 1. Directory rules, ALWAYS evaluated (except on overlay backing-file
	// events, see below). The label is the union of every matching rule and
	// the xattr label (labeller ADR), so an explicit '-' (public) xattr cannot
	// override a rule. Rules run BEFORE the xattr so the verifier explores the
	// bounded match loop once per program instead of once per xattr outcome.
	//
	// The backing-file event of an overlay open (FMODE_BACKING) carries an
	// f_path on a private clone mount whose root is the layer directory, so
	// d_path yields a path relative to the layer, which could match an
	// unrelated rule. Rules are therefore skipped for that event; the overlay
	// file's own event carries the merged path and is matched normally.
	if (!backing) {
		__u64 dl = TM_LABEL_PUBLIC;
		int dr = label_from_dir(file, pb, &dl);
		if (dr == TM_SRC_HIT) {
			label |= dl;
			if (dl != TM_LABEL_PUBLIC)
				cause |= TM_CAUSE_DIR;
		} else if (dr == TM_SRC_ERROR) {
			label |= TM_LABEL_UNKNOWN;
			cause |= TM_CAUSE_ERROR;
		}
	}

	// 2. xattr.
	__u64 xl = TM_LABEL_PUBLIC;
	int xr = label_from_xattr(file, pb, &xl);
	if (xr == TM_SRC_HIT) {
		label |= xl;
		if (xl != TM_LABEL_PUBLIC)
			cause |= TM_CAUSE_XATTR;
	} else if (xr == TM_SRC_EACCES) {
		label |= TM_LABEL_UNKNOWN;
		cause |= TM_CAUSE_ERROR | TM_CAUSE_EACCES;
	} else if (xr == TM_SRC_ERROR) {
		label |= TM_LABEL_UNKNOWN;
		cause |= TM_CAUSE_ERROR;
	}

	// 3. Cache floor, keyed by the real data inode, when available. Only
	// definitive XATTR-derived bits are cached (never rule-derived bits, per
	// the labeller ADR, and never TM_LABEL_UNKNOWN, which is per-event).
	if (have_real_data_inode) {
		struct inode *di = bpf_real_data_inode(file);
		if (di) {
			__u64 *cached = bpf_inode_storage_get(&inode_cache, di,
							      NULL, 0);
			if (cached && (*cached & ~label)) {
				label |= *cached;
				cause |= TM_CAUSE_CACHE;
			}
			__u64 add = (xr == TM_SRC_HIT ? xl : 0) & TM_LABEL_CACHEABLE;
			if (add && (!cached || (*cached | add) != *cached)) {
				__u64 *v = bpf_inode_storage_get(&inode_cache, di,
						NULL, BPF_LOCAL_STORAGE_GET_F_CREATE);
				if (v)
					__sync_fetch_and_or(v, add);
			}
		}
	}

	if (label == TM_LABEL_PUBLIC)
		return; // public, nothing to raise

	cause |= hook;
	if (backing)
		cause |= TM_CAUSE_BACKING;

	record_cause(slot, cause);
	raise_label(slot, label, cause);
}

// ---------------------------------------------------------------------------
// Sleepable LSM programs. Section names match cilium/ebpf's "lsm.s+" matcher
// (elf_sections.go, v0.22.0) and the sleepable_lsm_hooks set in
// kernel/bpf/bpf_lsm.c at tag v7.2: file_open, file_receive and mmap_file are
// all sleepable.
// ---------------------------------------------------------------------------

// file_open fires inside do_dentry_open (fs/open.c) via security_file_open,
// before the fd is returned, so the label is up before open() returns. For an
// overlayfs file it fires twice: once for the overlay file (merged path) and
// once for the backing (real) file that ovl_open opens through
// backing_file_open. The backing event's f_path lives on a private clone mount
// rooted at the layer directory, so d_path gives a layer-relative path, not
// the absolute lower path; handle_file skips directory rules for it.
SEC("lsm.s/file_open")
int BPF_PROG(tm_file_open, struct file *file)
{
	handle_file(file, TM_CAUSE_HOOK_OPEN);
	return 0;
}

// file_receive fires in receive_fd (fs/file.c), which the SCM_RIGHTS path
// reaches through scm_detach_fds (net/core/scm.c) and pidfd_getfd reaches
// through kernel/pid.c.
SEC("lsm.s/file_receive")
int BPF_PROG(tm_file_receive, struct file *file)
{
	handle_file(file, TM_CAUSE_HOOK_RECEIVE);
	return 0;
}

// mmap_file covers a read performed by mapping a file, including a descriptor
// that was opened before the process joined the session (scenario 10).
//
// NULL HAZARD: security_mmap_file passes file == NULL for anonymous mappings
// (mm/util.c), which the Go runtime and ld.so create at startup. The verifier,
// however, types this argument as a trusted, never-NULL pointer (it has no
// __nullable suffix in lsm_hook_defs.h at v7.2), so a plain "if (!file)" test
// is removed as dead code and NULL would reach the fs kfuncs and oops the
// kernel. The test below therefore fetches the raw argument as a scalar with
// bpf_get_func_arg, which LSM trampoline programs may call (kernel/trace/
// bpf_trace.c, tracing_prog_func_proto, bpf_prog_has_trampoline) and which the
// verifier cannot fold. Unlike bpf_probe_read_kernel it needs no
// allow_ptr_leaks and is not withheld under lockdown=confidentiality. This
// must be confirmed on the kernel under test (scenario 10). The argument
// arguably should be marked __nullable upstream; that report is a Phase 1
// follow-up.
//
// Set TM_ATTACH_MMAP_FILE=0 in the loader's environment to remove this program
// from the collection entirely (neither loaded nor attached) if it fails to
// verify or misbehaves on the kernel under test.
SEC("lsm.s/mmap_file")
int BPF_PROG(tm_mmap_file, struct file *file, unsigned long reqprot,
	     unsigned long prot, unsigned long flags)
{
	__u64 raw = 0;

	if (bpf_get_func_arg(ctx, 0, &raw))
		return 0;
	if (raw == 0)
		return 0; // anonymous mapping: no file
	// Observed on Fedora 7.2.0-61.fc45 (2026-10-18): this kernel already types
	// the argument as trusted_ptr_or_null, so the verifier requires an explicit
	// NULL test before any dereference. On kernels that type it non-NULL this
	// test is folded away, and the scalar test above is the real guard.
	if (!file)
		return 0;
	handle_file(file, TM_CAUSE_HOOK_MMAP);
	return 0;
}
