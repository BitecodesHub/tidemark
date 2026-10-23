# Tidemark Phase 0 kernel spike

Status: Phase 0 measurement instrument, written 2026-10-16. **Built and run on
2026-10-18** in a Lima VM on the M1 Pro (Fedora 45 Beta, kernel
7.2.0-61.fc45, aarch64, SELinux enforcing): all 12 scenarios passed in three
identical runs. **On 2026-10-23 four scenarios were added (13 to 16) and the
suite passes 16 of 16 in three consecutive runs**, with two further configuration
runs: rebuilt with the `Dual BSD/GPL` licence tag, and with the inode cache
compiled out via `TM_DISABLE_INODE_CACHE=1`. Evidence:
`docs/spike-results/2026-10-23_Kernel-Spike-Results-aarch64_v2.md` and the five
logs beside it; the earlier run is
`docs/spike-results/2026-10-18_Kernel-Spike-Results-aarch64_v1.md`. The x86_64
CI job has not run yet. This is a measurement instrument for the Phase 0 kill
gate, not product code. Build and run happen only inside a Fedora 45 virtual
machine.

This spike proves, or fails to prove, the read-time label mechanism on **stock
Fedora 45** (no custom image), on aarch64 in Lima on the M1 and on x86_64 under
QEMU/KVM in CI. The product label engine is built later, in Phase 1, only if
this gate passes.

## What the probe does

`bpf/label_probe.bpf.c` attaches three sleepable BPF LSM programs:

- `lsm.s/file_open` raises a per-session label before `open()` returns.
- `lsm.s/file_receive` raises it when a descriptor arrives over SCM_RIGHTS (and
  through `pidfd_getfd`, which also reaches `security_file_receive` via
  `receive_fd`).
- `lsm.s/mmap_file` raises it when a file is mapped, including a descriptor
  that was opened before the process joined the session. Anonymous mappings
  pass `file == NULL` even though the verifier types the argument as a
  non-NULL trusted pointer, so the program fetches the raw argument as a scalar
  with `bpf_get_func_arg` and tests it before touching the pointer (no
  `allow_ptr_leaks` needed, not withheld under lockdown; to be confirmed on the
  kernel under test). `TM_ATTACH_MMAP_FILE=0` removes this program from the
  collection before load, so if it fails to verify the rest of the suite still
  runs.

A session is any process whose level-2 ancestor cgroup is registered. Only
**regular files** are considered; pipes, sockets and anon-inode files are
skipped. A file's label is the **union** of its `security.bpf.tidemark` xattr
label (exactly one byte: `p`, `s` or `-`) and the labels of **every** directory
rule (labelled root) whose prefix matches the file's resolved path on a
path-component boundary, as the labeller ADR and `spec/clearance-model.md`
define it. Directory rules are evaluated first, so the verifier explores the
bounded match loop once per program. Rules are skipped for the overlay
backing-file event, whose path is relative to the layer root (see scenario 03).

Lookups fail **closed**: any error other than "no such xattr" (`-ENODATA`) or
"xattrs not supported" (`-EOPNOTSUPP`, treated as unlabelled, unlike the ADR's
stricter wording) sets the `TM_LABEL_UNKNOWN` bit, which consumers treat as at
least secret. `-EACCES` (at v7.2 `bpf_get_file_xattr` checks read permission
with the caller's credentials) carries its own cause bit `eacces`; a missing
per-task path buffer, retried once, carries `nobuf`. Known over-labelling cases
follow from the `-EACCES` rule: a write-only open of a file the caller cannot
read, and an exec of an execute-only binary (Fedora ships `sudo` as mode 4111).
**Decided 2026-10-23 (decisions ADR PD-5): an exec-only open counts as a read
and may raise the label; over-labelling an execute-only binary is accepted,
because under-labelling a private file would leak it.** The probe does not strip
the `" (deleted)"` suffix that d_path appends for unlinked files, so such files
never match a rule.

An inode-storage cache keyed by the real data inode (`bpf_real_data_inode`,
present from kernel 7.2) holds **xattr-derived** bits only and acts as a
**floor**: both sources are evaluated on every event and the cached bits are
ORed on top, so the cache can neither hide a raised xattr nor a stricter rule,
and it never lowers a label. A per-slot **cause** bitmask records which hook
(`open`, `receive`, `mmap`) and which mechanism (`xattr`, `dir`, `cache`,
`error`, `backing`) produced each label, so the scenarios assert the mechanism
under test, not only the final value.

**The cache is not part of the product design.** `TM_DISABLE_INODE_CACHE=1`
zeroes `have_real_data_inode` at load time, which compiles out every cache read
and write; the full suite passes that way. Scenario 14 measures what the cache
costs: with it on, removing a file's label xattr leaves the label raised on the
next open (`label=0x3 cause=open|cache`), because a floor can never lower a
label and declassification lowers; with it off the same operation reads
`label=0x0 cause=none`. That is the evidence behind decisions ADR PD-1 - the
Phase 1 labeller keeps no inode-keyed label state and therefore needs neither
`bpf_real_data_inode` nor the `inode_post_setxattr`/`inode_post_removexattr`
invalidation that was previously listed as Phase 1 work.

The label word lives in a `BPF_F_MMAPABLE` array. Every program returns 0; the
probe only observes, it never denies.

`cmd/tidemark-probe` loads and pins the object (removing stale pins first),
then runs the scenario scripts under `tests/`. Each scenario runs in its own
process group with a timeout (`TM_SCENARIO_TIMEOUT`, default 120 s), and its
output is streamed to the console as it happens, so a hang or a kernel oops
still leaves its last lines in the serial log. A scenario may raise its own
limit with a `# TIMEOUT: <seconds>` line in its first 60 lines - the race and
latency scenarios need 1800 s and 1200 s. SIGINT and SIGTERM unload the probe.
The run prints one `PASS`/`FAIL` line per scenario and a final `SUMMARY:` line,
and exits non-zero on any failure.

## How to run

### aarch64 on the M1 (Lima)

```
cd spike
limactl start --name tidemark-spike lima/fedora-45.yaml
limactl shell tidemark-spike -- sudo make -C /opt/tidemark-spike test SUDO=
```

macOS asks once for permission before `limactl` can share a folder under
`~/Desktop`. The build writes `bpf/vmlinux.h`, the object, the binary, `go.sum`
and `.cache/` into the shared `spike/` directory; `spike/.gitignore` excludes
everything except `go.sum`, which was reviewed and committed after the first
successful run.

### x86_64 in CI (QEMU/KVM, driven over SSH)

The `.github/workflows/kernel-spike.yml` workflow runs on `workflow_dispatch`
and on any change under `spike/`. It boots a stock Fedora 45 Cloud image (the
pinned Beta compose, kernel 7.2.0-61.fc45, deliberately not the newest 7.2.x,
for reproducibility) under QEMU/KVM on a free `ubuntu-24.04` runner, then builds
and runs `make test` **over SSH** through QEMU user-mode networking with a host
forward, and gates on `SUMMARY: ALL PASS` in the captured transcript plus a zero
exit status from the guest. Re-pin to the GA image within a week of the Fedora
45 release.

Three things about that design are load-bearing, and all three were learned from
a failed run rather than from reading:

- **The verdict is never read off the serial console.** The first version ran
  the whole suite inside a cloud-init `runcmd` with stdout redirected to
  `/dev/ttyS0` and grepped the log for a completion marker. In run
  39965539646 (2026-10-23) `serial-getty@ttyS0` restarted while dnf was still
  working; its `vhangup()` replaced the runcmd shell's already-open stdout and
  stderr with the hung-up tty's file operations, so every later write returned
  `-EIO` and was discarded. The guest powered off cleanly, `qemu` exited 0, and
  the transcript stopped at the point of the restart: the run could not be read
  as either a pass or a fail. The serial log is still uploaded, but only as
  evidence about boot and kernel oops.
- **The SSH key is generated per run.** `ssh-keygen` makes an ed25519 pair in
  the runner's temp directory and only the public half goes into the NoCloud
  seed, for a `spikeci` user in the `sudo` group with `NOPASSWD`. Nothing is
  committed and nothing is a secret.
- **The source reaches the guest as `git archive HEAD spike`, not as the working
  tree.** The local `spike/` directory holds aarch64 products (`bpf/vmlinux.h`,
  `label_probe.bpf.o`, `tidemark-probe`) whose mtimes are newer than every
  source file, so a copy of the tree would let the guest's `make` treat an
  arm64 object and binary as already built for x86_64. The workflow asserts the
  absence of all three in the staging directory before it builds the ISO.

`make test` installs a copy of the helper to `/usr/local/libexec/tidemark-probe`
and runs that copy, because the helper is executed inside sessions and the
Lima shared folder is virtiofs, whose xattr errors could otherwise set the
fail-closed bit on the helper binary itself.

### By hand inside any Fedora 45 VM

```
sudo bash provision/fedora-deps.sh
make vmlinux   # generate bpf/vmlinux.h from the running kernel's BTF
make build     # compile the BPF object (clang -mcpu=v3) and the Go probe, then go vet
sudo make test SUDO=   # load, attach and run every scenario
```

A passing run ends with `SUMMARY: ALL PASS (16/16)`.

## What each test proves, and what it does not

Each scenario maps to one Phase 0 assumption. "Proves" means the scenario, if it
passes, is evidence for that assumption on the kernel under test. "Does not
prove" lists what the scenario deliberately leaves out. Every in-session helper
first asserts that it is a **grandchild** of the process that joined the
session (pid, parent and grandparent are printed) and that it runs in the
expected cgroup; a mismatch is reported as a setup error, never as a kernel
result. Every labelled file is verified with `getfattr` before use.

| Script | Assumption | Proves | Does not prove |
|---|---|---|---|
| `01_label_before_open_grandchild.sh` | A1 | A grandchild's `open()` of an xattr-labelled file has raised the label when sampled immediately after `open()` returns, before any `read()`; cause is `open|xattr`. | Overlays, directory rules, fd passing. |
| `02_session_cgroup_keying.sh` | A2 | Out-of-session opens do not raise the slot; in-session opens do; a process in a child cgroup below the scope still counts (level-2 ancestor); an unregistered sibling scope does not. A scope at the wrong depth is refused at registration. | The label value beyond raised or not raised. |
| `03_overlay.sh` | A3 | **Gate:** a lower file with no xattr, opened through the overlay, is labelled by a directory rule on the merged path, with cause `dir` and no `xattr`. **Observation:** whether an xattr-labelled lower file opened through the overlay is labelled with no matching rule. The backing-file open also reaches `security_file_open`, with `f_path` on a private clone mount rooted at the layer directory (so its d_path is layer-relative and rules are skipped for it), and its xattr read is reported with cause `backing`; whether the cache contributes is reported too. | That the overlay inode itself exposes the lower xattr. |
| `04_directory_rule.sh` | A4 | A rule matches only on a path-component boundary (`<root>X` does not match `<root>`); an unlabelled file under a root is labelled with cause `dir`; overlapping rules combine (`0x3`). | Aliases created in another mount namespace (scenario 11). |
| `05_file_receive_scm.sh` | A5 | A descriptor opened outside the session and passed in over SCM_RIGHTS raises the label by the time `recvmsg()` returns; cause `receive|xattr`, not `open`. | An unprivileged receiver (scenario 09). |
| `06_mmapable_label.sh` | A6 | A `PROT_READ`, `MAP_SHARED` mapping returns the raised value; the kernel refuses `mprotect(PROT_WRITE)` on it; a `BPF_F_RDONLY` descriptor refuses a writable shared mapping. Observes whether a `BPF_F_RDONLY` descriptor can be mapped read-only at all (expected not to at v7.2). | Mapping from a non-root daemon UID (Phase 1). |
| `07_setxattr_eperm.sh` | A7 | A capability-less in-session process gets `EPERM` for `security.bpf.tidemark`, while the same uid sets `user.tm_probe` on the same file (positive control), so `EPERM` comes from the capability check; root outside the session succeeds. | Filesystems a session mounts in its own user namespace (mitigated by G1; see scenario 11). |
| `08_mmap_read_label.sh` | A8 | The open-then-map read path leaves the label raised. It does **not** isolate `mmap_file`; the open already raises the label. | `read()` on an inherited descriptor (scenario 10). |
| `09_file_receive_unprivileged.sh` | A5b | An unprivileged receiver of a root-only (0600) labelled descriptor still raises the label: the xattr read fails with `EACCES` and the probe fails closed, exactly `0x4` with cause `receive|error|eacces` and no `xattr`. | The production mapping of "unknown" to a clearance. |
| `10_inherited_descriptor.sh` | A8b | **Gate:** `mmap_file` alone raises the label for a descriptor opened before the process joined the session (cause `mmap`, not `open`). **Recorded gap:** `read()` on that descriptor fires no attached hook and leaves the label at `0x0`. | Any defence for the read gap; that is the shim's `close_range` before exec (Phase 1). |
| `11_userns_alias.sh` | A4b (observation) | Records whether a capability-less principal (uid 65534 with `unshare --user --map-root-user --mount`) can reach a rule-labelled file through a bind-mount alias that no rule matches. Concludes "no gap" only if the cause shows `dir`. Not a gate. | Anything about xattr-labelled files, which follow the inode. |
| `12_link_rename_alias.sh` | A4c (observation) | Records whether a hard link or a rename out of a labelled root, made inside the session's own mount namespace, escapes path-based rules. Not a gate. | A fix. The decision is recorded (decisions ADR PD-4: labels follow the location, Phase 1 denies link/rename out of a ruled root); the probe never denies, so the enforcement itself is untested here. |
| `13_session_move_out.sh` | A2c | **Gate (AC3):** an unprivileged in-session process cannot write its own pid into another cgroup's `cgroup.procs` (a sibling scope and the cgroupfs root both give `EACCES`) and cannot create a sub-cgroup under its scope (`EACCES`); the session still raises its label on the next labelled open, and a root `try-move` on the same destination succeeds as the positive control. | That a *privileged* process cannot leave the scope - root can, which is why the session is granted no privilege (guarantee G1). |
| `14_xattr_and_rule_invalidation.sh` | A12 | **Gate (D4, AC3):** raising a file's attribute from `p` to `s` is honoured on the next open (the secret bit appears, so the cache is a floor and never a hide); adding a directory rule labels the next open of a pre-existing file (`0x2`, cause `dir`, no `xattr`); clearing the rule un-labels it (`0x0`). Records the attribute-removal result as the cache's declassification gap. | A production deny path; and, under PD-1, nothing needs invalidating - the no-cache run reads `0x0` after removal. |
| `15_race_10k_iterations.sh` | A1r | **Gate (AC1 strict form, the labeller ADR kill test):** 10,000 iterations of clear-slot then in-session `open()` of a labelled file, asserting the label is already raised when sampled immediately after `open()` returns (`RACE_SELF … misses=0`), and, in parallel, that a **separate observer process outside the session** sees it raised after receiving an ack'd byte from the opener (`RACE_OBSERVE … misses=0`). `TM_RACE_ITERATIONS` overrides the count. | Load from a real agent workload; the observer is a stand-in for the gateway, not the gateway. |
| `16_open_latency_baseline.sh` | D5 (baseline, no gate) | Records `open()`+`close()` latency over 20,000 measured iterations (2,000 warmup) for five cases: outside a session, in-session unlabelled, in-session xattr-labelled steady state, the raise path with the slot cleared before each open, and a matching directory rule. Prints min/p50/p99/max/mean per case and asserts the rule still labels at the end. | End-to-end session overhead; it brackets `open(2)` only, with the probe loaded, in one Go process. |

Assumption IDs are local to this spike. They map to `docs/THREAT_MODEL.md`
Section 10.2 as follows: A1 and A1r = S7, A2 and A2c = S2, A3 = S1, A12 = the
invalidation part of S1 (now a no-cache design, PD-1) plus checklist AC3, A4 =
the directory-rule part of checklist criterion AC3 (S9, rule matching inside a
session's own mount view, is covered only by the observation in A4b), A5 and A5b
= S3, A6 = part of S11 (the array is mapped as root, not yet from a non-root
daemon UID), A7 = S5, and A8 and A8b = S4. A4b and A4c are observations for S9.
S6 is covered by the load step, the `env` banner and the boot self-test; S8 and
S10 are not covered by this spike, and S16 (latency) maps to no assumption - it
is the D5 baseline.

### Session identity and the fixed cgroup depth

A spike session is a transient systemd scope under the slice **`tmspike.slice`**.
The slice name is deliberately free of dashes: systemd nests a slice named
`a-b.slice` under `a.slice`, which would change the depth. With a dash-free
slice, the scope cgroup sits at depth 2 below `/sys/fs/cgroup`
(root = 0, `tmspike.slice` = 1, scope = 2), so the loader sets `session_level`
to 2 and `bpf_get_current_ancestor_cgroup_id(2)` returns the scope's cgroup id
however deep below the scope a process sits. `register-session` reads the
scope's live cgroup id (the inode number of its cgroupfs directory) and refuses
a cgroup at any other depth. The real risks are a scope created at the wrong
level and a process leaving the scope; nesting below the scope does not break
keying (scenario 02 case c).

### The CAP_SYS_ADMIN emulation

A real Tidemark session runs without CAP_SYS_ADMIN. The spike's systemd scope
runs as root, so scenarios 07, 09 and 11 use `setpriv` (uid 65534, empty
bounding set, no_new_privs) where an unprivileged principal is needed. Phase 1
drops the capability through the systemd unit properties of the real session.

### Known gaps this spike records rather than closes

- **Inherited descriptors.** A descriptor or mapping inherited across fork,
  exec or a cgroup move is never seen by `file_open` or `file_receive`, and
  `read()` is not hooked (no `file_permission` program; that hook is not
  sleepable). Requirement for Phase 1: every session starts from a fresh exec
  with no inherited descriptors above stderr (the shim's `close_range`).
- **Path-based rules and namespaces.** `bpf_path_d_path` resolves against the
  caller's root and mount namespace. A session that can create a mount or user
  namespace can reach a rule-labelled file through an alias, and inside its own
  user namespace it holds CAP_SYS_ADMIN over filesystems it mounts there.
  Mitigation: sessions cannot create namespaces (guarantee G1), and important
  files carry xattr labels, which follow the inode.
- **Hard links and renames out of a labelled root.** These need no new
  namespace: d_path names the opened dentry, `linkat` and `renameat` fire no
  attached hook, and rule-derived bits are not cached per inode. Scenario 12
  records the gap. **Decided 2026-10-23 (decisions ADR PD-4): a directory label
  follows the location, not the content.** The Phase 1 labeller therefore
  attaches sleepable `path_link` and `path_rename` and returns `-EPERM` for an
  in-session link or rename whose destination is not covered by the same rule as
  the source; it does not stamp an xattr onto the destination inode. The refusal
  must name the rule that caused it. This spike cannot test the enforcement,
  because the probe never denies.

## How a failure maps to the kill criteria

The Phase 0 kill gate (concept document, "Kill criteria", end of week 3) reads:
"kernel label raising is unreliable on stock Fedora 45 on either architecture or
through overlayfs, and the fanotify fallback misses reads by child processes or
through mmap. Action: drop the claim of kernel-recorded labels and ship the
gateway, broker and receipts as an audit tool, or stop."

- **Kill gate directly:** A1, A2, A4 or A8 failing on either architecture. The
  core label raise is unreliable. Ship audit-only, or stop.
- **Kill gate for the overlay case:** A3 failing (no label raise through the
  overlay by the directory rule). The staged-writes design depends on labels
  surviving the overlay.
- **Kill gate, session containment half:** A2c failing. A session process that
  could move itself out of its scope escapes the labeller entirely, so the
  label-before-read guarantee would hold only for processes that stay put.
- **Kill gate, ordering half:** A1r failing on any iteration. A miss means the
  label is not always raised before `open()` returns, which is the exact claim
  the concept makes. A12 failing to honour a raised attribute or an added or
  removed rule means the labeller cannot be trusted to reflect policy changes.
- **Narrow the claim, no kill:** A5, A5b or A8b failing. Fd passing, unprivileged
  receipt or mapped pre-session descriptors would be unlabelled read paths,
  which the public claim must then exclude.
- **Label integrity broken:** A7 failing. A session could rewrite its own labels.
- **Transport redesign, no kill:** A6 failing. The exits need another way to
  observe the label word; the label raise itself is unaffected.
- **Observation only:** A4b and A4c never fail the suite; they inform the
  containment requirements, and A4c informed the PD-4 decision.
- **Not a gate:** the latency scenario. It records the baseline PD-1 was decided
  on and gives Phase 1 a number to compare against.

**Which wording of the kill criterion governs is decided (PD-7):** the concept
document's full wording, not the plan's summary row, and a triggered criterion
cannot be overridden. Where that makes this list stricter than checklist AC1 or
AC3, the stricter reading applies.

If the kernel label raise fails but a fallback is wanted, the concept names a
fanotify `FAN_OPEN_PERM` vault, whose gaps for mmap reads, grandchild reads, and
overlay reads must be documented in the labeller ADR before it can be selected.
This spike does not implement that fallback; it only produces the evidence that
decides whether it is needed.

## Licence of the BPF source (decided 2026-10-23, decisions ADR PD-3)

The kernel allows the BPF LSM kfuncs used here only for GPL-compatible programs.
The repository is Apache-2.0, which is not GPLv2-compatible, so BPF C sources
are licensed **`GPL-2.0-only OR BSD-2-Clause`** with the kernel licence string
`"Dual BSD/GPL"`; the exception and the file list are recorded in `NOTICE`. This
is not taken on argument alone: the object was rebuilt with the dual tag and the
full suite still passes 16 of 16
(`docs/spike-results/2026-10-23_Spike-Run-aarch64-dual-licence_v1.log`), so
kfunc access is preserved. Phase 1 BPF code follows the same pattern.

## Verified facts and sources

All kernel facts were read from the `torvalds/linux` tree at tag **v7.2** on
2026-10-16, cross-checked against the Fedora 45 package versions below. Tag v7.2
is the kernel series Fedora 45 ships, so it is the correct reference for a
stock-Fedora-45 spike.

- **Fedora 45 ships kernel 7.2.** The Fedora 45 Beta compose carries
  `kernel-7.2.0-61.fc45`; the `development/45` tree carries `kernel-7.2.9-300.fc45`
  (package listings on `dl.fedoraproject.org`, 2026-10-16). As of 2026-10-16
  Fedora 45 is still at **Beta** (45_Beta 1.3, composed 2026-09-05); only Fedora
  44 is GA. The Lima template and the CI workflow therefore pin the 45_Beta 1.3
  Generic Cloud images; re-pin to the GA images and digests when Fedora 45 is
  released.
- **Landlock ABI on kernel 7.2 is 10, which satisfies the ABI 9 floor.** The
  ABI-to-kernel mapping is 7.0 = 8, 7.1 = 9, 7.2 = 10, 7.3-rc = 11
  (`security/landlock/syscalls.c`; recorded in
  `docs/research/2026-10-16_OS-Feasibility_v1.md`). The probe reports the ABI it
  measures at runtime.
- **Sleepable BPF LSM hooks.** In `kernel/bpf/bpf_lsm.c` at v7.2, the
  `sleepable_lsm_hooks` set contains `file_open`, `file_receive`, `mmap_file`,
  `inode_setxattr`, `inode_post_setxattr`, and `inode_post_removexattr`.
  `file_permission` is not sleepable and is not used.
- **fs kfuncs.** In `fs/bpf_fs_kfuncs.c` at v7.2: `bpf_get_file_xattr`
  (`KF_SLEEPABLE`), `bpf_get_dentry_xattr` (`KF_SLEEPABLE`), `bpf_path_d_path`,
  and `bpf_real_data_inode` (`KF_SLEEPABLE | KF_RET_NULL`, signature
  `struct inode *bpf_real_data_inode(struct file *file)`). All are restricted to
  `BPF_PROG_TYPE_LSM` by `bpf_fs_kfuncs_filter`.
- **`bpf_real_data_inode` is new in the 7.2 cycle.** It was added as
  `bpf_real_inode(dentry)` (commit `9af8c8a54f6e`, vfs-7.2-rc1, 2026-05-26) and
  then renamed with a `struct file *` signature to `bpf_real_data_inode`
  (commit `3f8c65b06faf`, committed 2026-07-01). It does **not** exist in v7.0 or v7.1. The
  rename commit states the kfunc "landed this cycle and has no released users."
  Fedora 45's 7.2 kernel should therefore carry it, but because it is brand new
  and is not a stable ABI, the probe detects it at load time from kernel BTF and,
  when it is absent, compiles the data-inode cache out of the verified program
  (the `have_real_data_inode` gate). See the fallback note below.
- **`bpf_get_current_ancestor_cgroup_id` is a helper available to LSM programs**,
  under `CONFIG_CGROUPS`, returning the cgroup id of the ancestor at the given
  level, or 0 if there is none (`kernel/bpf/helpers.c`; exposed to LSM through
  `bpf_lsm_func_proto` in `kernel/bpf/bpf_lsm.c`).
- **`file_open` runs before `open()` returns.** `do_dentry_open` (`fs/open.c`)
  calls `security_file_open(f)` before the descriptor is installed.
- **`file_receive` covers SCM_RIGHTS.** The SCM_RIGHTS receive path
  (`net/core/scm.c` -> `receive_fd` in `include/linux/file.h`) calls
  `security_file_receive` before installing the descriptor.
- **Overlay data inode.** `bpf_real_data_inode` resolves to the lower/upper data
  inode via `d_real_inode(file_dentry(file))`. `bpf_get_file_xattr` uses
  `file_dentry(file)`, which for an overlay file is the overlay dentry; whether
  the overlay passes the lower `security.bpf.*` xattr through is the behaviour
  `03_overlay.sh` measures at runtime.
- **`security.bpf.*` requires CAP_SYS_ADMIN.** `cap_inode_setxattr`
  (`security/commoncap.c`) returns EPERM for any `security.*` xattr without
  `CAP_SYS_ADMIN` in the filesystem's user namespace. `XATTR_NAME_BPF_LSM` is
  `"security.bpf."` (`include/uapi/linux/xattr.h`).
- **Inode storage in sleepable LSM.** `BPF_MAP_TYPE_INODE_STORAGE` is in the
  verifier's sleepable map list at v7.2; `bpf_inode_storage_get` accepts a
  `struct inode *` that may be NULL.
- **Library versions.** `github.com/cilium/ebpf` latest release is **v0.22.0**
  (2026-06-29); its `go.mod` requires `go 1.25.0` and
  `golang.org/x/sys v0.43.0`. The loader uses `link.AttachLSM`, `unix.Mmap` on a
  read-write map descriptor for the `BPF_F_MMAPABLE` array (a `BPF_F_RDONLY`
  descriptor cannot be mapped `MAP_SHARED` at v7.2), `LoadPinnedMap` with `ReadOnly` for checks, the
  `"lsm.s+"` sleepable-LSM section matcher, the weak-kfunc poisoning in
  `linker.go`, and `btf.Spec.TypeByName` for the kfunc detection. `go.sum` is
  not committed in this draft; the Makefile `build` target (`GOFLAGS=-mod=mod`)
  writes it inside the VM.
- **Fedora 45 toolchain (development/45 tree, 2026-10-16):** `golang-1.27.1`,
  `clang-23.1.2`, `bpftool-7.7.0`, `libbpf` and `libbpf-devel` available.
- **GitHub Actions pinned by commit SHA:** `actions/checkout` v7.0.1 =
  `3d3c42e5aac5ba805825da76410c181273ba90b1`; `actions/upload-artifact` v7.0.2 =
  `cf430e030ddbb5b0abf93d22962f4752f3646cd9`.
- **Fedora 45 Beta Cloud Generic images (verified against the GPG-signed
  CHECKSUM):**
  - aarch64 qcow2 `Fedora-Cloud-Base-Generic-45_Beta-1.3.aarch64.qcow2`,
    sha256 `2bacfd793c2e769a325520da64a0bfcd38a7bbd755561d3d332a35b5ae83169c`.
  - x86_64 qcow2 `Fedora-Cloud-Base-Generic-45_Beta-1.3.x86_64.qcow2`,
    sha256 `06bd4382a3bc5e5f94cd520f9ee15bdfb38d7bb72552b27fe06d581ba9a84292`.
  - Fallback, Fedora 44 GA Generic (from the Lima upstream `fedora-44`
    template): aarch64 sha256
    `55c60a3b80d3616a08705afd0459e75fe9f03c54aba7a46e4002a41a72fa0d5b`,
    x86_64 sha256
    `28680fe5b371a5a82ebf43a31926e086a168e59949d03969c5093e7071f90b7f`.
- **Base bootc image (for cross-reference; the image pipeline is outside this
  spike):** `quay.io/fedora/fedora-bootc:45` resolved on 2026-10-16 to
  `sha256:86cc9d98ef210b71a6c4eeec4b7d04e7fef0a4415184c2322a14c879df189f45`;
  the `:44` fallback to
  `sha256:9db34195bc0b6be43e17c20ca840be5bc153093c94b62f6792a76fb199eeaced`.
  Both were quarantined by quay.io manifest GC on 2026-10-23 and are
  superseded; `docs/adr/2026-10-16_BaseImage_v1.md` is the authoritative and
  dated re-pin record.

## The `bpf_real_data_inode` fallback

When kernel BTF lacks `bpf_real_data_inode` (any kernel before 7.2), the loader
sets `have_real_data_inode = 0` and the verifier prunes the branch that calls
the kfunc, so the cache is inactive. Because the cache is only a floor over
xattr-derived bits, labels are still computed from the xattr and the directory
rules on every event; only the observation of shared data inodes is lost.
`03_overlay.sh` reports which mechanisms labelled the overlay open, which told
the labeller ADR whether the kfunc is needed at all.

**It is not needed.** The ADR's question is answered: the no-cache path is the
production design (decisions ADR PD-1), so `bpf_real_data_inode` is dropped from
the required-kfunc list and the loader's runtime detection is kept only as the
spike's own experiment switch (`TM_DISABLE_INODE_CACHE=1`). The practical upshot
is that the "absent on 7.1" argument for the 7.2 kernel floor disappears; the
floor stays because Phase 0 tested no other kernel (PD-2).

## Assumptions that only the in-VM run can settle

- Whether an xattr-labelled lower file is labelled when opened through an
  overlay, and by which event (overlay open or backing-file open). Reported by
  `03_overlay.sh`.
- Whether the `mmap_file` NULL guard verifies and behaves as intended (scenario
  10; `TM_ATTACH_MMAP_FILE=0` is the escape hatch).
- Whether the BPF program verifies and loads on the actual Fedora 45 kernel, and
  whether every scenario passes. **Settled for aarch64 on 2026-10-23: 16 of 16
  in three runs**, and the same for the dual-licence and no-cache builds. The
  x86_64 answer is still pending CI.
- Whether Fedora 45's shipped kernel enables `bpf` in the active LSM list by
  default. The probe asserts it from `/sys/kernel/security/lsm`.
- The exact Landlock ABI reported at runtime (expected 10).
- xattr and `FS_APPEND_FL` behaviour on virtiofs and 9p mounts (recorded as
  observations in later phases, not gated here).

## Placeholders to change

- The Go module path `github.com/agentos-tidemark/tidemark/spike` is a
  placeholder. **Decided 2026-10-23 (PD-10): the repository is
  `github.com/BitecodesHub/tidemark`, created private, and the module path is
  updated to it.** Moving it to an organisation, or making the repository
  public, is a launch decision.
- The image pins and digests change when Fedora 45 reaches GA.
