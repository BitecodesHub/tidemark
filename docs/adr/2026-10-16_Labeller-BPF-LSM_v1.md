# Read-time labels via sleepable BPF LSM, with a fanotify fallback

- **Status:** Proposed (Phase 0 draft, revised 2026-10-16 after review,
  2026-10-17 in a cross-file consistency pass, and 2026-10-23 after the spike
  runs on both architectures). The aarch64 half of the kill test passes 16 of 16
  in three runs and the x86_64 half passes 16 of 16 in four CI runs the same
  day. The lookup method is now decided (no cache), so the
  C2 gate that held this ADR open is closed on the evidence in
  [../spike-results/2026-10-23_Kernel-Spike-Results-aarch64_v2.md](../spike-results/2026-10-23_Kernel-Spike-Results-aarch64_v2.md)
  and
  [../spike-results/2026-10-23_Kernel-Spike-Results-x86_64_v1.md](../spike-results/2026-10-23_Kernel-Spike-Results-x86_64_v1.md);
  the decisions are recorded in
  [2026-10-23_Phase0-Founder-Decisions_v1.md](2026-10-23_Phase0-Founder-Decisions_v1.md).
  The ADR stays Proposed because the founder has not reviewed and promoted it,
  not because evidence is missing. The kernel-label claim is still dropped if
  the kill test fails on either architecture.
- **Date:** 2026-10-16
- **Phase:** 0 (Foundations and kill gate)
- **Deciders:** Founder; the 2026-10-23 revision records decisions made under
  the founder's explicit delegation that day.

## Context

The one new kernel piece in Tidemark is the read-time label: when any process in
a session opens a labelled file for reading, or receives a readable descriptor
to one, the session's label must rise **before the call returns**, across the
whole process tree, and the agent must not be able to lower it. Everything else
(containment, exits, staging, receipts) is integration. If this mechanism is not
reliable on the stock Fedora kernel, the kernel-label claim cannot ship.

Facts verified against the mainline kernel source on 2026-10-16 (files at the
`torvalds/linux` tags named):

- **The kfuncs exist.** In `fs/bpf_fs_kfuncs.c` at v7.2: `bpf_get_file_xattr`,
  `bpf_get_dentry_xattr`, `bpf_set_dentry_xattr` and `bpf_remove_dentry_xattr`
  are `KF_SLEEPABLE`; `bpf_real_data_inode` is `KF_SLEEPABLE | KF_RET_NULL`;
  `bpf_path_d_path` carries no sleepable flag and is restricted to LSM programs.
  `bpf_get_file_xattr` reads only `user.*` and `security.bpf.*` names and checks
  `MAY_READ` permission of the current task on the inode.
- **Sleepable hooks.** In `kernel/bpf/bpf_lsm.c` at v7.2 the sleepable set
  contains `file_open`, `file_receive`, `inode_post_setxattr`,
  `inode_post_removexattr`, `inode_rename`, `path_rename` and `path_link`, but
  not `file_permission`.
- **`bpf_real_data_inode`** returns `d_real_inode(file_dentry(file))`: the
  upper or lower inode that stores a regular file's data on overlayfs, and the
  attached inode elsewhere. It is **absent at v7.1** and present at v7.2 and
  v7.3-rc7. It was added in the 7.2 cycle (commit `9af8c8a54f6e`) and changed to
  take a `struct file` in commit `3f8c65b06faf`, whose message says it had "no
  released users". fedora-bootc:45 ships kernel 7.2.9 and fedora-bootc:44 ships
  7.2.10 (image labels, 2026-10-16).
- **Overlay opens produce two `file_open` events.** The overlay-level open runs
  `security_file_open` in `do_dentry_open` before `ovl_open`, which calls
  `ovl_open_realfile` and then `backing_file_open` (`fs/overlayfs/file.c`,
  `fs/backing-file.c`). The backing file is allocated with `FMODE_BACKING`
  (`fs/file_table.c`) and its `f_path` lies on a layer mount that
  `clone_private_mount` created with `mnt_ns = MNT_NS_INTERNAL`
  (`fs/overlayfs/super.c`, `fs/namespace.c`). `d_path` on such a path treats the
  mount as detached and returns a layer-relative path such as `/notes.txt`
  (`fs/d_path.c`). Directory rules therefore cannot be matched on the backing
  event.
- **`d_path` appends `" (deleted)"`** to the path of an unlinked dentry
  (`fs/d_path.c`).
- **Sessions cannot set the label attribute.** `cap_inode_setxattr` requires
  `CAP_SYS_ADMIN` for any `security.*` name (`security/commoncap.c`).

## Decision

Implement read-time labels as a sleepable BPF LSM program, and keep a fanotify
vault as a documented fallback.

**Primary design (BPF LSM):**

1. Attach sleepable programs to `file_open` and `file_receive`.
2. Identify the session by its ancestor cgroup id at the fixed slice-relative
   depth (`bpf_get_current_ancestor_cgroup_id`), with `Delegate=no` on every
   session unit. Exit immediately if the process is not in a session cgroup.
3. Skip files whose `f_mode` lacks `FMODE_READ`: such a descriptor cannot read
   or map data. Skip backing-file events (`FMODE_BACKING`) for rule matching.
4. At the **overlay-level (session-visible) open**, compute the label as the
   union of:
   - the `security.bpf.tidemark` attribute, read with `bpf_get_file_xattr`
     (overlayfs passes `security.bpf.*` reads through to the real inode, and
     copy-up copies the attribute; the spike's overlay scenario records the
     first as an observation, not a gate, and does not yet test copy-up), and
   - every directory rule matching the session-visible path from
     `bpf_path_d_path`, relative to the session root. A trailing
     `" (deleted)"` is removed first. A prefix matches when the path equals it
     or continues it with `/`. A candidate implementation looks up each
     `/`-boundary prefix of the path in a hash map, in a bounded loop.
   This relies on preconditions P5 and P6 in `spec/clearance-model.md`:
   overlays and read paths are bound at their own absolute paths, sessions do
   not use `RootDirectory=`, and rule-labelled files cannot be renamed or linked
   out of their rule.
5. **Fail closed.** If the attribute read fails with any error other than
   `-ENODATA`, or the path cannot be resolved into the buffer, deny the open or
   receive with `-EACCES`. Never treat an error as "unlabelled".
6. Atomically OR the label bits into the session's slot of a `BPF_F_MMAPABLE`
   array, read `bpf_ktime_get_ns` after the OR, and emit a ring-buffer event
   only when the value changed. Count failed ring-buffer reservations so
   tidemarkd can record a `resync` (see `spec/receipt-predicate-v0.md`).

**Lookup method and cache: decided on 2026-10-23 - no cache (PD-1).** The
production labeller evaluates the attribute and the directory rules on every
readable open and keeps no inode-keyed label state. Checklist item C2 required
the choice to follow a latency measurement; the measurement is
`spike/tests/16_open_latency_baseline.sh` (20,000 iterations per case, aarch64,
kernel 7.2.0-61): an open inside a session costs 1,211 ns at p50 whether or not
the file carries an attribute, and the labeller's total overhead over the
outside-session open is about 376 ns at p50. A cache cannot buy back a
meaningful fraction of that.

The deciding evidence is not latency but correctness. With the cache active, the
spike's invalidation scenario shows that removing a file's
`security.bpf.tidemark` attribute leaves the label raised on the next open
(`label=0x3 cause=open|cache`), because a cache that only ORs bits upward is a
floor: it can never lower a label, and declassification lowers. With
`TM_DISABLE_INODE_CACHE=1` the same operation reads `label=0x0 cause=none`, and
the full 16-scenario suite - including the overlayfs scenario - still passes.
So the no-cache design closes the removal gap by construction instead of by
adding the invalidation machinery the previous version of this section required.

What a future cache would still have to satisfy, if Phase 1 benchmarking on a
real workload reopens this decision:

- It caches only the attribute component, keyed by `bpf_real_data_inode(file)`.
  Directory-rule results are never cached per inode: a rule labels a path, and
  one inode can have several paths (hard links), so no epoch scheme makes an
  inode-keyed rule cache sound.
- It never caches an error.
- It caches absence (`-ENODATA`) only with a per-inode generation counter that
  `inode_post_setxattr` and `inode_post_removexattr` increment, stored with
  compare-and-swap against the value read before the attribute read. Caching
  positive results only would leave the hot path (unlabelled files) uncached
  and is not worth its complexity.
- It is invalidated on attribute removal, which the last two constraints are the
  minimum for and which the spike measured as missing.

**Kernel floor.** v0.1 requires **Linux 7.2 or later and Landlock ABI 9 or
later** (PD-2). From Phase 1, doctor and tidemarkd refuse to start sessions when
a required kfunc is missing; the Phase 0 boot self-test reports it as a failure.
With PD-1 the required-kfunc list **drops `bpf_real_data_inode`**, because the
labeller no longer calls it; the list is the kfuncs the design actually uses
(`bpf_get_file_xattr`, `bpf_path_d_path`, and the sleepable LSM attach path).
The floor stays at 7.2 anyway: Phase 0 has no evidence on 7.1, and this project
does not claim support it has not tested. Any support for 7.1 kernels is
deferred to Phase 6 (portable mode).

**Phase 0 kill test (the gate for this ADR).** On stock Fedora 45 (aarch64 in
Lima on the M1, x86_64 under QEMU/KVM in CI), the probe must show, each case for
10,000 iterations with zero misses where timing matters:

- the label word is already raised when the **grandchild** that called `open()`
  reads it immediately after the call returns, and when a separate **observer**
  process (standing in for the gateway) reads it after receiving a byte the
  grandchild wrote after `open()` returned;
- ancestor-cgroup keying holds with `Delegate=no`: a grandchild in a nested
  cgroup raises the session's slot, and a session process cannot move itself out;
- **overlay, attribute:** a lower file labelled by attribute, opened through a
  workspace overlay, raises the label; compare `bpf_get_file_xattr` at the
  overlay-level event with the backing event;
- **overlay, rule:** a lower file labelled **only** by a directory rule, opened
  through a workspace overlay, raises the label;
- **rule boundaries:** `Private` labels `Private/x` but not `PrivateX/x`;
- **moves and links:** open F (unlabelled), `mv` F into `~/Private`, reopen:
  labelled. `mv` a directory into `~/Private`: its files are labelled on the
  next open. Reload a policy that adds a rule: new opens are labelled.
  Inside a session, `mv` or `ln` of a rule-labelled file to an unruled path
  fails or is labelled (precondition P5);
- `file_receive` raises the label for a descriptor passed with `SCM_RIGHTS`;
- a read-only `mmap` of the label word reflects the raise;
- `setxattr` on `security.bpf.*` returns `EPERM` inside a session and the map
  pins cannot be opened.

This kill test is stricter than the other Phase 0 records. As of 2026-10-23 the
gap is closed on aarch64: the spike now runs each case once **and** runs the
10,000-iteration race in both readings this section requires - the opening
grandchild sampling its own label, and a separate observer process sampling
after receiving a byte (`spike/tests/15_race_10k_iterations.sh`, zero misses in
both, three runs) - plus the move-out case
(`spike/tests/13_session_move_out.sh`) and the move/link and invalidation cases
(`spike/tests/14_xattr_and_rule_invalidation.sh`). Threat-model assumption S7
and the plan still place the race test in Phase 1; the spike runs it now because
this ADR asks for it.

Two of the cases above sit below the kill threshold in the spike's failure
mapping (`spike/README.md`): a failure of the `SCM_RIGHTS` case narrows the
public claim, and a failure of the label-word mapping case requires a redesign of
the label transport; neither is a kill on its own. **PD-7 settles the resulting
ambiguity in favour of this ADR and the concept:** where the records differ, the
stricter wording gates week 3, and a triggered criterion cannot be overridden.
`SCM_RIGHTS` and the label-word mapping therefore gate the kernel-label claim as
written in `spike/README.md` (they narrow the claim rather than kill the project)
only because that mapping is itself the recorded decision; the kill condition
text above is what the public concept promises, and both are now reported.

If label raising is unreliable on Fedora 45 or through overlayfs, and the
fanotify fallback cannot cover the gaps below, the project drops the
kernel-label claim and ships the gateway, broker and receipts as an audit tool,
or stops. This is the week-3 kill criterion.

## Decided: the aliasing question from the Phase 0 spike review (2026-10-17, closed 2026-10-23)

Directory rules are path-based. A hard link or a rename out of a labelled root,
made inside the session's own mount namespace, gives the file a path that no
rule matches; `linkat` and `renameat` fire no hook the labeller attaches, and
rule-derived bits are not cached per inode. Scenario
`spike/tests/12_link_rename_alias.sh` records this.

**PD-4 decides that a directory label follows the location, not the content, and
that the session is denied the operations that would break the promise.** From
Phase 1 the labeller attaches sleepable `path_link` and `path_rename` (both in
the v7.2 sleepable set, verified 2026-10-16 against `kernel/bpf/bpf_lsm.c`) and returns `-EPERM` for an in-session link or rename
whose destination is not covered by the same rule that covers the source. It
does not stamp an xattr onto the destination inode.

Reasons, in short: content binding would keep a private label on a file the user
deliberately published, would require the labeller to write xattrs on inodes the
session may not own (more authority than `spec/threat-model.md` grants), and
would contradict a policy language whose rules are path prefixes. `mv` and `ln`
out of a ruled root failing with `EPERM` is the honest, user-visible form of the
same guarantee. The Phase 1 acceptance condition is that the refusal names the
rule that caused it; a bare `EPERM` is not shipped.

The same review found that the overlay backing-file event carries a
layer-relative path, so the spike skips directory rules for that event; the
production labeller does the same, as step 3 and step 4 already say, by ignoring
`FMODE_BACKING` events for rule matching rather than by matching on root device
and inode.

The spike does not test `path_link`/`path_rename` enforcement: Phase 0 has no
deny-mode program (the probe is observe-only). That is Phase 1 work, and the
risk it leaves open is recorded in the Consequences section below.

## Consequences

- The mechanism rides on mainline kernel primitives; the new kernel code is the
  small label program (about 250 lines of BPF C plus about 300 lines of Go
  loader with `cilium/ebpf` and `bpf2go`).
- Labelling at `file_open` (not `file_permission`) means one evaluation per
  open, so descriptors must be covered by `file_receive` and by `close_range` in
  the shim before exec, so no labelled descriptor is inherited untracked.
- Per-open `bpf_path_d_path` and attribute reads cost time on open-heavy
  workloads (npm install, large compiles). D5 records a baseline; Phase 1
  measures a real workload. The cache stays optional for this reason.
- kfuncs are not a stable ABI. CI asserts their presence, and Renovate digest
  bumps merge only after the read matrix passes.
- The 7.2 floor supersedes the earlier 7.1 fallback. `docs/THREAT_MODEL.md` (G3
  row, Section 7 and Section 10.1) and `image/selftest/tidemark-selftest.sh`
  (kernel floor 7.2, `bpf_real_data_inode` required) now match this ADR; the
  threat model keeps the 7.1 approach only as a record for a later decision.
- The checked-in spike (`spike/bpf/label_probe.bpf.c`) follows step 4 and the
  cache constraints above: it computes the label as the union of the attribute
  and every matching directory rule, caches only attribute-derived bits keyed
  by `bpf_real_data_inode(file)` (as a floor ORed onto a fresh evaluation of
  both sources, never as an answer), and never caches an error. It fails closed
  by setting an unknown-label bit (`TM_LABEL_UNKNOWN`) on a lookup error. As an
  observe-only measurement instrument, it still differs from this design in
  these respects: it never denies an operation (step 5 denies with `-EACCES`);
  it treats `-EOPNOTSUPP` (no xattr support) like `-ENODATA`; it does not skip
  write-only opens or backing-file events (step 3); it does not remove a
  trailing `" (deleted)"` before rule matching (step 4); it also attaches
  `mmap_file`; and its cache is not invalidated on `inode_post_setxattr` or
  `inode_post_removexattr` (Phase 1 work). The Phase 1 label engine follows
  this ADR, not the spike.

### fanotify FAN_OPEN_PERM fallback and its coverage gaps

If the BPF path fails the gate, the fallback is a `fanotify` `FAN_OPEN_PERM`
vault. Its gaps must be documented before it can be selected:

- **Overlay and namespace marks.** The overlay-level open does fire
  `FAN_OPEN_PERM` (`do_dentry_open` calls `fsnotify_open_perm_and_set_mode` for
  every opened file). The gap is placement: a mount mark is per mount, so a mark
  on the real lower or upper mount misses backing opens, which happen on cloned
  private mounts, and a mark on the host's overlay mount misses the copy of that
  mount in each session's mount namespace. The vault would need filesystem
  marks (`FAN_MARK_FILESYSTEM`) on every overlay and backing filesystem, which
  route every open on those filesystems through the daemon.
- **mmap.** `FAN_OPEN_PERM` fires at open, not at page fault. In the mainline
  source, `mmap` generates pre-content events only for HSM descriptors.
- **Inherited and received descriptors.** A process that reads through a
  descriptor it inherited or received via `SCM_RIGHTS` issues no `open()`.
- **Session attribution.** fanotify reports a PID, not a cgroup; the vault needs
  an auxiliary cgroup lookup that races with process exit.
- **Performance and races.** Every open on a marked filesystem round-trips
  through a userland daemon, which conflicts with the overhead goal for
  unlabelled opens and adds a userland race between the read and a later send.

fanotify is therefore a fallback only, at a stated lower guarantee level.

## Alternatives considered

- **Hook `file_permission`.** Rejected: it is not sleepable, so it cannot call
  the attribute kfuncs.
- **Cache directory-rule results per inode with a global epoch bumped by
  rename, link and policy load.** Rejected: hard links that existed before the
  cache filled give one inode two paths with different rule results.
- **Match rules on the backing-file event.** Rejected: its path is
  layer-relative.
- **Write labels forward with `bpf_set_dentry_xattr`.** Deferred: v0.1 stamps on
  commit instead, to keep the Phase 0 surface small.
- **fanotify as primary.** Rejected for the gaps above.
- **ptrace or an LD_PRELOAD shim.** Rejected: neither covers arbitrary
  unmodified process trees, and both are evaded easily.

## References

- `docs/research/2026-10-16_Tidemark-Concept_v1.md` (read-time labels; red-team
  items on overlay, cgroup keying, incomplete taint).
- `spec/clearance-model.md` (preconditions P1 to P6) and
  `spec/policy-v0.schema.json` (`directory_rules`).
- Linux source at tags v7.1, v7.2 and v7.3-rc7 (checked 2026-10-16):
  `fs/bpf_fs_kfuncs.c`, `kernel/bpf/bpf_lsm.c`, `fs/overlayfs/file.c`,
  `fs/overlayfs/super.c`, `fs/backing-file.c`, `fs/file_table.c`, `fs/open.c`,
  `fs/namespace.c`, `fs/d_path.c`, `fs/notify/fsnotify.c`,
  `include/linux/fs.h`, `security/commoncap.c`.
- Kernel commits `9af8c8a54f6e` and `3f8c65b06faf`.
