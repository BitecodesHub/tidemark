# Kernel spike results: aarch64, four new scenarios and the cache question (2026-10-23)

**Status:** evidence for the Phase 0 kill gate, aarch64 half. This report covers
the four scenarios added on 2026-10-23 (13 to 16), the two configuration
experiments run with them, and the latency baseline. It continues
[2026-10-18_Kernel-Spike-Results-aarch64_v1.md](2026-10-18_Kernel-Spike-Results-aarch64_v1.md),
whose 12 scenarios still pass unchanged.

## Executive summary

**All 16 scenarios pass, 16 of 16, in three consecutive runs** on stock Fedora 45
Beta (kernel 7.2.0-61.fc45, aarch64, SELinux enforcing). The three runs are
identical in every PASS/FAIL, RESULT, DESC and race line; only latency numbers
differ, as expected.

The four additions close the open items from the 2026-10-18 report:

- **AC3 move-out (scenario 13):** an unprivileged in-session process cannot move
  itself out of its `Delegate=no` scope and cannot restructure that scope; both
  attempts return `EACCES`, and the session still raises its label afterwards.
  Root *can* write the same file, which is the positive control and the reason
  G1 gives sessions no privilege.
- **AC1 strict form (scenario 15):** the labeller ADR's 10,000-iteration kill
  test now runs, with zero misses in **both** readings required by the ADR - the
  opening grandchild sampling its own label right after `open()` returns, and a
  **separate observer process** outside the session sampling it after receiving a
  byte from the opener. 10,000 and 10,000, no misses, three runs.
- **Xattr and rule invalidation (scenario 14):** raising a file's attribute from
  `p` to `s` is honoured on the next open (the cache is a floor, never a hide).
  **Removing the attribute is not honoured** - the cached bit keeps the label
  raised (`label=0x3 cause=open|cache`) - which turns the Phase 1 "invalidation"
  task from a theory into a measured defect. Directory rules are never cached:
  adding a rule labels the next open of a pre-existing file, and clearing it
  un-labels the next open exactly.
- **D5 latency baseline (scenario 16):** 20,000 measured `open()`+`close()`
  iterations per case. The labeller costs about **+0.37 µs at p50** for an open
  inside a session on an ordinary unlabelled file, and the read-xattr path costs
  no measurable more than that at p50; the raise path (label cleared before every
  open) is about **+0.5 µs at p50**.

## Two configuration experiments

1. **Licence tag (`Dual BSD/GPL`).** The BPF object previously declared `"GPL"`
   while its SPDX line said Apache-2.0. The licence decision (PD-3 of the
   decisions ADR) made it `GPL-2.0-only OR BSD-2-Clause` with the kernel tag
   `"Dual BSD/GPL"`. The suite was rebuilt and re-run: **16/16 pass**, so the
   kfunc calls are still accepted as GPL-compatible
   (`2026-10-23_Spike-Run-aarch64-dual-licence_v1.log`).
2. **No inode cache (`TM_DISABLE_INODE_CACHE=1`).** The loader was given a
   switch that zeroes `have_real_data_inode`, which compiles out every cache
   read and write. The full 16-scenario suite passes with the cache off
   (`2026-10-23_Spike-Run-aarch64-no-cache_v1.log`), including the overlay
   scenario. Two differences are informative:
   - scenario 14(b) then reads `label=0x2` exactly, not `0x3`: no stale cached
     bit is ORed onto the fresh attribute value;
   - scenario 14(c) then reads `label=0x0 cause=none`: removing the attribute
     takes effect on the next open, which is what declassification needs.
   This is the evidence behind PD-1: **the production labeller keeps no
   inode-keyed label cache**, and `bpf_real_data_inode` is therefore not a
   required kfunc.

## Scenario results (2026-10-23 additions)

| Scenario | Assumption | Result | Evidence in the log |
|---|---|---|---|
| 13 move-out | A2c / AC3 | PASS | unprivileged `try-move` to a sibling scope and to the cgroupfs root both `EACCES`; `try-mkdir` under the scope `EACCES`; next labelled open still `0x1`; root `try-move` on the same file `ok` |
| 14 invalidation | A12 / D4, AC3 | PASS | `p` -> `0x1 (open\|xattr)`; raised to `s` -> `0x3 (open\|xattr\|cache)` with the secret bit set; after `removexattr` -> `0x3 (open\|cache)` recorded as the gap; rule added -> `0x2 (open\|dir)`; rule cleared -> `0x0` |
| 15 race, 10k | A1r / AC1 strict | PASS | `RACE_SELF iterations=10000 misses=0`; `RACE_OBSERVE iterations=10000 misses=0`; cause `open\|xattr` |
| 16 latency | D5 | PASS | five `BENCH` lines below; raise path `misses=0 of 20000`; rule still labels after 60,000 opens |

## open() latency baseline (scenario 16, 20,000 measured iterations each, 2,000 warmup)

| Case | min | p50 | p99 | mean |
|---|---|---|---|---|
| outside a session, unlabelled file (probe early exit) | 751 ns | 835 ns | 960 ns | 862 ns |
| inside a session, unlabelled file, no rule | 1,126 ns | 1,211 ns | 1,503 ns | 1,265 ns |
| inside a session, xattr-labelled file (steady state) | 1,127 ns | 1,211 ns | 1,628 ns | 1,269 ns |
| inside a session, raise path (slot cleared and label checked before each open) | 1,544 ns | 1,711 ns | 7,387 ns | 2,359 ns |
| inside a session, matching directory rule | 1,168 ns | 1,252 ns | 2,212 ns | 1,297 ns |

Read as: the labeller adds about **376 ns at p50** to every open made by a
session process (row 2 minus row 1); reading the attribute on a file that carries
one adds nothing measurable beyond that (row 3 equals row 2 at p50); a matching
directory rule adds about 41 ns (row 5); the first open that actually raises the
label is about 500 ns more expensive than the steady state (row 4). Cold-first
opens in the loop were 6.5 to 13.4 µs and include page-cache and dentry effects,
not only the hook.

Method caveats, stated for the record: the probe stays loaded for the whole
suite, so row 1 is "hook attached but process outside a session", not "no probe";
the loop runs in one Go process and `time.Since` brackets only `open(2)`; the
`p99` of row 4 includes the harness's own map-clear and label-read syscalls
scheduling noise, not kernel work. Phase 0 sets no target for these numbers (the
checklist says so); they are the input to the lookup-method decision and a
reproducible baseline for the Phase 1 benchmark methodology.

## What is still open for the kill gate

- **x86_64 (AC2, D6):** `.github/workflows/kernel-spike.yml` must produce the
  same 16 of 16 under QEMU/KVM in CI. It needs the repository on GitHub.
- **Fedora 45 GA kernel:** these runs pin the Beta kernel 7.2.0-61 on purpose;
  re-run when GA ships.
- **Criterion 2 (interviews):** independent of this result, and still at zero of
  eight interviews.
- Scenario 14's removal gap is closed by PD-1 (no cache), not by Phase 1
  invalidation work; the labeller ADR records that.
