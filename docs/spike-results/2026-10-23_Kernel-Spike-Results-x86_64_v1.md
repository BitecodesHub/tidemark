# Kernel spike results: x86_64 in CI, the second half of kill criterion 1 (2026-10-23)

**Status:** evidence for the Phase 0 kill gate, x86_64 half. Same suite, same
scripts, same source as
[2026-10-23_Kernel-Spike-Results-aarch64_v2.md](2026-10-23_Kernel-Spike-Results-aarch64_v2.md),
run on a different architecture inside GitHub Actions instead of on the host.
Environment: `2026-10-23_Spike-Environment-x86_64_v1.txt`; transcripts:
`2026-10-23_Spike-Run-x86_64-ci_v1.log` and the four repeats beside it.

## Executive summary

**All 16 scenarios pass, 16 of 16, in four consecutive CI runs** under QEMU/KVM
on a stock Fedora 45 Beta x86_64 guest (kernel `7.2.0-61.fc45.x86_64`, SELinux
enforcing, Landlock ABI 10, lockdown `[none]`, BTF present). A fifth
consecutive green run followed the same day on the same test code; it is
described in "A fifth run, and what it changes" below rather than merged into
the four-run tables. The four
transcripts are identical in every PASS/FAIL, RESULT, DESC and race line; only
latency numbers differ, as expected. Three of them were separate
`workflow_dispatch` runs after the push-triggered one, so the result is not an
artifact of one runner at one moment.

This closes the x86_64 half of kill criterion 1. The aarch64 half closed on the
same day; the two together mean the criterion's wording - reliable label raising
**on both architectures** and **through overlayfs** - has now been measured on
both, not asserted on one.

## Cross-architecture agreement

Every behaviour the gate depends on reproduces exactly on x86_64:

| Observation | aarch64 (three runs) | x86_64 (four runs) |
|---|---|---|
| Grandchild `open()` raises before return, opener and observer both see it | 10,000 / 10,000 iterations, 0 misses | 10,000 / 10,000 iterations, 0 misses |
| Overlayfs, xattr path | `label=0x1 cause=0x111(open\|xattr\|backing)` | identical, two `file_open` events, the backing one layer-relative |
| Overlayfs, directory rule on an unlabelled lower file | `label=0x1 cause=0x21(open\|dir)` | identical |
| SCM_RIGHTS receive into a session | `label=0x1 cause=0x12(receive\|xattr)` | identical |
| Unprivileged receiver of a root-only fd (fail-closed) | `label=0x4 cause=0x282(receive\|error\|eacces)` | identical |
| Session cannot move out of, or restructure, its `Delegate=no` scope | `EACCES` on sibling scope, cgroupfs root and `try-mkdir`; root control `ok` | identical |
| Directory rules: component boundary, add, clear | `PrivateX` 0x0, `Private` 0x1, `Private/keys` 0x3; add 0x2, clear 0x0 | identical |
| Attribute removal with the cache on (the measured defect behind PD-1) | `label=0x3 cause=0x41(open\|cache)` | identical |
| `security.bpf.*` `setxattr` from a capability-less session process | `EPERM` (positive control `user.*` ok) | identical |
| Read-only mmap of the label word | `PROT_WRITE` refused `EACCES`; writable mapping of a `BPF_F_RDONLY` fd refused | identical |
| Mount-namespace alias and hard-link/rename escapes | gaps recorded, `label=0x0` | identical |

The architecture-specific build detail that matters: the CI job regenerates
`bpf/vmlinux.h` from the running guest's BTF and compiles with `-target bpf
-mcpu=v3 -D__TARGET_ARCH_x86`, so these runs prove the object builds and loads
from source on x86_64 rather than replaying an aarch64 artifact.

## open() latency baseline (scenario 16, 20,000 measured iterations each, 2,000 warmup)

Median of the four x86_64 runs, `p50` in ns, with the range of each statistic across runs:

| Case | min | p50 (range) | p99 (range) | mean (range) |
|---|---|---|---|---|
| outside a session, unlabelled file (probe early exit) | 2,454 | 2,535 (2,504-2,575) | 8,776-9,779 | 3,064-3,147 |
| inside a session, unlabelled file, no rule | 3,977 | 4,098 (4,047-4,158) | 18,755-22,803 | 4,847-5,251 |
| inside a session, xattr-labelled file (steady state) | 4,227 | 4,344 (4,268-4,479) | 22,902-25,698 | 5,224-5,554 |
| inside a session, raise path (slot cleared and label checked before each open) | 5,179 | 5,651 (5,591-5,681) | 39,113-43,672 | 7,352-7,747 |
| inside a session, matching directory rule | 4,097 | 4,208 (4,158-4,259) | 22,351-23,434 | 5,044-5,296 |

`BENCH_READLABEL` - the harness's own map read, which brackets the label-slot read
rather than the hook - is p50 19,136-20,047 ns across the four runs, and every
run reports `misses=0 of 20000`. Scenario 15's own `open()` timings are p50
12,173-12,754 ns, which include the harness's slot-clearing and label-reading
syscalls and are not a hook measurement.

Read as: on x86_64 the attached labeller adds about **1.56 µs at p50** to a
session process's open (row 2 minus row 1, about +62% of the un-hooked open),
reading the attribute adds about 0.25 µs beyond that, a matching directory rule
about 0.11 µs, and the path that actually raises the label about 1.55 µs more
than the steady state. The same shape as aarch64 (where the delta was ~0.38 µs,
+45%) at a higher absolute cost - which is expected, because the two absolute
floors differ (835 ns vs 2,535 ns) and the CI guest is CPU-oversubscribed: QEMU
is given 4 vCPU on a runner whose KVM recommends 2, so the hook's work competes
with the harness for host CPU, and the fourth run overlapped two image builds on
the same runner.

Phase 0 sets no target for these numbers, and the checklist says so; they are a
reproducible baseline, not a product figure. The same method caveats as the
aarch64 report apply, and one more here: the probe stays loaded for the whole
suite, the loop runs in one Go process and brackets only `open(2)`, and the CI
guest shares its host.

## The verification channel, and why it is SSH

The first attempt at this job printed its verdict to the serial console and read
back nothing after the 31-second mark: the guest had run the suite, powered off
cleanly and exited 0, and the gate failed on "no marker". The cause is not
Tidemark-specific: when `serial-getty@ttyS0` restarts during the dnf
transaction, its `vhangup()` replaces the already-open stdout of the cloud-init
`runcmd` shell with hung-up tty file operations, and every later write returns
`-EIO` and is dropped. A gate that cannot read its own result is worse than a
failing gate, so the job now boots the guest with a per-run ephemeral ed25519
key, runs the suite over SSH, and keeps the serial log only as boot and oops
evidence. Recorded as decision PD-16 in
[../../docs/adr/2026-10-23_Phase0-Founder-Decisions_v1.md](../../docs/adr/2026-10-23_Phase0-Founder-Decisions_v1.md)
and as the red warning block in `spike/README.md`.

## A fifth run, and what it changes

A fifth consecutive green run landed later the same day:
[39982982413](https://github.com/BitecodesHub/tidemark/actions/runs/39982982413),
push-triggered at commit `7b02d80`. Between `3f7c091` (the four runs above) and
`7b02d80` the only file under `spike/` that differs is `spike/README.md`, which
is not compiled, so the probe object, the loader and the test binaries are the
same code. Its transcript is committed beside the others as
`2026-10-23_Spike-Run-x86_64-ci_v5.log`.

It is recorded as a fifth repetition and deliberately NOT folded into the tables
above, because those tables say "median of the four runs" and re-deriving them
would silently invalidate the range columns a reader can check against the
transcripts. What the fifth run does say:

- The 50 decision lines (RESULT, PASS/FAIL, DESC, LABEL, CAUSE, RACE, SUMMARY)
  are byte-identical across all five transcripts once pids, temp-directory
  names, cgroup ids and timings are normalised. Five runs, five
  `SUMMARY: ALL PASS (16/16)`.
- Scenario 16's five cases keep the same order, and the hook's cost is the same:
  in-session unlabelled p50 4,138 ns minus outside-session 2,564 ns = +1,574
  ns, against the four-run figure of +1,563 ns.
- Five statistics in three of its rows fall outside the four-run ranges in the
  table, each by less than a microsecond: the in-session unlabelled mean 4,792
  against a range starting at 4,847 (-55 ns), the raise-path p50 5,751 against a
  range ending at 5,681 (+70 ns), and the directory-rule run's p50 4,268, p99
  21,931 and mean 4,996 against ranges of 4,158-4,259, 22,351-23,434 and
  5,044-5,296 (+9 ns, -420 ns, -48 ns). Read correctly, those range columns are
  the spread of four samples, not prediction intervals, and a fifth sample
  landing just outside three of them is what a 0.1-1.9% run-to-run spread looks
  like. The claim being made here - that the hook adds about 1.5-1.6 µs at p50
  and never misses - is unaffected; the point is that the table must not be
  quoted as a bound.

## What is now closed, and what is not

Closed by these runs:
- **AC2 / D6** - the x86_64 probe run, published as a CI artifact.
- The x86_64 half of **kill criterion 1**, and with it the "both architectures"
  clause that the concept's wording requires and that PD-7 refused to soften.
- **AC3** on x86_64 (scenario 13 had only been met on aarch64).

Not closed, and independent of this result:
- **Fedora 45 GA kernel.** These runs pin the Beta kernel `7.2.0-61` on purpose,
  for reproducibility. Re-run when GA ships, within a week of its release per
  the base-image ADR.
- **Criterion 2 (interviews):** still zero of eight. Founder-owned.
- **The gate decision itself (AC8)** and the founder's review of the ADRs that
  this evidence feeds.
- Scenario 14's attribute-removal gap stays open in the spike because the spike
  still ships the cache; PD-1 removes it from the product, and the aarch64
  no-cache run is the evidence that the suite passes without it.
