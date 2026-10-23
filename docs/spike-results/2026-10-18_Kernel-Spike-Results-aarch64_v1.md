# Kernel spike results: aarch64, Fedora 45 Beta (2026-10-18)

**Status:** evidence for the Phase 0 kill gate, aarch64 half only. The x86_64
half (`.github/workflows/kernel-spike.yml`) has not run yet, because the
repository has not been pushed to GitHub.

## Executive summary

On stock Fedora 45 Beta (kernel 7.2.0-61.fc45, aarch64, SELinux enforcing) the
read-time label mechanism worked: **12 of 12 scenarios passed in three
consecutive, identical runs** (normalised results compared line by line). Every
in-session helper proved it was a grandchild of the process that joined the
session, and every pass shows the intended mechanism in its cause bits, not
only the final label. The first real build exposed three kernel-level issues
that reading sources had missed; all three are fixed in the probe and recorded
below. The three gaps the spike was designed to record (inherited descriptors,
namespace aliases, hard links and renames) were all confirmed.

**Kill-gate reading (criterion 1, end of week 3):** label raising is reliable
on stock Fedora 45 on aarch64, including through overlayfs. The criterion is
not yet closed: it requires both architectures, and criterion 2 (interviews)
is still open.

## Environment

| Item | Value |
|---|---|
| Host | Apple M1 Pro, 16 GB, macOS 27, Lima 2.1.3 (vz), 4 vCPU, 6 GiB |
| Guest image | Fedora-Cloud-Base-Generic-45_Beta-1.3.aarch64.qcow2 (sha256 2bacfd79...169c) |
| Kernel | 7.2.0-61.fc45.aarch64 |
| Active LSMs | lockdown, capability, yama, selinux, bpf, landlock, ipe, ima, evm |
| Landlock ABI | 10 (floor 9) |
| `bpf_real_data_inode` | present in kernel BTF |
| Lockdown | none |
| SELinux | Enforcing (selinux-policy-45.15-2.fc45) |
| Toolchain | clang 23.1.3, bpftool 7.7.0, Go 1.27.2, systemd 261.2 |
| Libraries | cilium/ebpf v0.22.0, golang.org/x/sys v0.43.0 (`spike/go.sum`) |

Raw evidence: `2026-10-18_Spike-Run-aarch64_v1.log` (run 3, full output) and
`2026-10-18_Spike-Environment-aarch64_v1.txt`.

## Results

| Scenario | Assumption | Result | Evidence in the log |
|---|---|---|---|
| 01 grandchild open | A1 / S7 | PASS | label `0x1`, cause `open|xattr`, sampled before `read()` |
| 02 cgroup keying | A2 / S2 | PASS | outside `0x0`; in scope `0x1`; child cgroup `0x1`; unregistered sibling scope `0x0` |
| 03 overlay | A3 / S1 | PASS | gate: `0x1`, cause `open|dir`, no `xattr` |
| 04 directory rules | A4 | PASS | `<root>X` `0x0`; under root `0x1` (`dir`); overlapping rules `0x3` |
| 05 SCM_RIGHTS | A5 / S3 | PASS | `0x1`, cause `receive|xattr`, no `open` |
| 06 label mapping | A6 / S11 (part) | PASS | `mprotect(PROT_WRITE)` refused `EACCES`; writable map of read-only fd refused `EACCES` |
| 07 setxattr | A7 / S5 | PASS | positive control `user.*` ok; `security.bpf.*` `EPERM`; root ok |
| 08 open then mmap | A8 / S4 | PASS | `0x1`, cause `open|mmap|xattr` |
| 09 unprivileged receiver | A5b / S3 | PASS | exactly `0x4`, cause `receive|error|eacces` (fail-closed) |
| 10 inherited descriptor | A8b / S4 | PASS | `mmap` alone: `0x1`, cause `mmap|xattr`, no `open` |
| 11 namespace alias | A4b (observation) | PASS | gap confirmed: alias path `0x0` |
| 12 link and rename | A4c (observation) | PASS | gap confirmed: hard link `0x0`, rename `0x0` |

## Issues the first real run exposed (all fixed)

1. **`vmlinux.h` needs `-fms-extensions`.** Kernel 7.x BTF contains tagged
   anonymous struct members (for example `struct ns_tree;` inside a struct),
   which clang rejects under `-Werror`. The Makefile now passes
   `-fms-extensions -Wno-microsoft-anon-tag`.
2. **`bpf_dynptr_from_mem` refuses stack memory on 7.2** ("Unsupported reg type
   fp"). The xattr buffer moved into the per-task storage element.
3. **Fedora's 7.2 kernel types the `mmap_file` argument as nullable**
   (`trusted_ptr_or_null_file`), contrary to the upstream header the reviewer
   read. The program now tests both the scalar argument and the pointer.

## Design inputs for Phase 1

- **Overlay:** the backing-file event reads the lower file's xattr (cause
  `xattr|backing`), and the data-inode cache never contributed. A no-cache
  labeller is viable, which bears on checklist decision C2, and it removes the
  dependency on the brand-new `bpf_real_data_inode` kfunc.
- **Label transport:** a `BPF_F_RDONLY` map descriptor cannot be mapped
  `MAP_SHARED` at all (`EINVAL`). The exits must either receive a read-write
  descriptor and map it read-only (the mapping cannot be upgraded), which
  leaves the descriptor itself writable through the bpf syscall, or receive
  the label through another channel. This needs an ADR before Phase 2.
- **Inherited descriptors:** `read()` on a descriptor opened before the
  session started raised nothing. The shim must close every inherited
  descriptor before exec (G3).
- **Namespace aliases:** a capability-less principal created a user and mount
  namespace and escaped a directory rule. Sessions must not create namespaces
  (G1).
- **Hard links and renames:** both escaped directory rules without any new
  namespace. This needs the founder's decision in the labeller ADR (labels
  follow the location or the content).
- **SELinux enforcing** did not interfere with loading, labelling or setxattr.

## What is still open for the kill gate

- Run `.github/workflows/kernel-spike.yml` on x86_64 (requires pushing the
  repository to GitHub), and confirm the same 12 of 12 result.
- Re-run on the Fedora 45 GA kernel when it is released (the Beta kernel is
  pinned on purpose for reproducibility).
- The 10,000-iteration race test (checklist AC1 versus the labeller ADR) and
  `open()` latency measurements (checklist D5 and C2) are not part of this run.
- Interview kill criterion 2 is independent of this result.
