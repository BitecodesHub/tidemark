# BPF probe review: bounded loops, early exit, failed kfunc handling (v1, 2026-10-23)

- **Checklist items:** S4 (security checks) against D3 (`spike/bpf/label_probe.bpf.c`).
- **Reviewed object:** `spike/bpf/label_probe.bpf.c` (595 lines) at the same commit as
  the aarch64 and x86_64 16-of-16 runs, plus the checked-in build product
  `spike/bpf/label_probe.bpf.o`.
- **Reviewer:** the Qoder agent that drafted the program, under the founder's
  2026-10-23 delegation ("you decide everything and stop when done").
- **Status:** the three properties S4 names are verified in the source and are
  listed below with the lines that establish them. **This is a self-review by the
  original author, not an independent one**, which is weaker than what the concept
  document means by "reviewed". An independent read is recorded as a Phase 1
  week-4 item beside S5, and the founder's line-by-line pass on D3 is still owed.

## 1. Bounded loops

There are exactly two loops in the program, both in `label_from_dir`, and both are
bounded by compile-time constants rather than by data:

- the outer loop over root slots runs `i < TM_MAX_ROOTS` with `TM_MAX_ROOTS` 8
  ([spike/bpf/label_probe.bpf.c:380](../spike/bpf/label_probe.bpf.c#L380)), over an
  `ARRAY` map of `max_entries = TM_MAX_ROOTS` (line 177);
- the inner prefix compare runs `j < TM_ROOT_LEN` with `TM_ROOT_LEN` 128 and an
  explicit `if (j >= rlen) break` (lines 393-401), and every memory access is
  index-masked (`pb->p[j & (TM_PATH_MAX - 1)]`,
  `r->prefix[j & (TM_ROOT_LEN - 1)]`) so the verifier can prove the bounds without
  reasoning about `rlen`.

No loop iterates over a path, a filename or a label value: the path buffer is
written by the kfunc, not walked character by character across its 4,096 bytes, and
label computation is a single pass per source. Worst-case instruction count is
therefore 8 x 128 byte comparisons plus three map lookups plus two sleepable kfunc
calls, which is what the verifier accepted on 7.2.0-61 on both architectures. There
is no recursion and no `bpf_tail_call`.

One deliberate consequence of the same constants: `TM_MAX_ROOTS = 8` means a ninth
directory rule silently does not exist for the probe. That is a spike limit, not a
product one, and it is stated here rather than left to be discovered.

## 2. Early exit for processes outside sessions

`handle_file` ([spike/bpf/label_probe.bpf.c:426](../spike/bpf/label_probe.bpf.c#L426))
returns before doing any work that can sleep or allocate, in this order:

1. `session_slot(&slot)` resolves the current task's ancestor cgroup id
   (`bpf_get_current_ancestor_cgroup_id(session_level)`) and looks it up in the
   `sessions` hash; `cid == 0` or a map miss returns 0 and the program exits
   (lines 244-254, 429-430). No kfunc call, no scratch allocation and no ring
   buffer touch happens for a process outside a session - which is every other
   process on the machine, because these programs are attached globally.
2. Non-regular files exit next: `i_mode` is compared against `S_IFREG` (line 437),
   so pipes, sockets, and anon-inode files (eventfd, epoll, BPF map fds) are
   skipped. This is not only a speed choice: an xattr read on a file another uid
   created returns `-EACCES`, and `handle_file` fails closed on `-EACCES`, so
   without this filter ordinary IPC would raise `TM_LABEL_UNKNOWN` in unrelated
   sessions.

The measured cost of the attached-but-not-raising path is the latency baseline in
scenario 16: an unlabelled open inside a session costs the same 1,211 ns at p50 as
an xattr-labelled one on aarch64, and the labeller's whole delta over the
outside-session baseline is about +376 ns p50 there (+1.56 µs on the
CPU-oversubscribed CI guest). Numbers and caveats:
`docs/spike-results/2026-10-23_Kernel-Spike-Results-x86_64_v1.md`.

## 3. Handling of failed kfunc calls

Every sleepable kfunc result is checked, and every unrecognised outcome fails
closed by raising `TM_LABEL_UNKNOWN` with its own cause bit, so a measurement can
never read "unlabelled" where the real answer was "could not tell".

| kfunc | outcome | handling |
|---|---|---|
| `bpf_get_file_xattr` | `-ENODATA`, `-EOPNOTSUPP` | `TM_SRC_NONE` - genuinely unlabelled (line 329) |
| | `-EACCES` | `TM_SRC_EACCES` raises `TM_LABEL_UNKNOWN` with `error\|eacces` (lines 331, 485-487) |
| | `!= 1`, unrecognised value byte, `-ERANGE` | `TM_SRC_ERROR` raises `TM_LABEL_UNKNOWN` with `error` (lines 333, 349-350, 488-490) |
| | `bpf_dynptr_from_mem` fails | `TM_SRC_ERROR` before the kfunc is called (line 325) |
| `bpf_path_d_path` | `<= 0` (including `-ENAMETOOLONG`) | `TM_SRC_ERROR` raises `TM_LABEL_UNKNOWN` (lines 371-373, 472-474) |
| `bpf_real_data_inode` | NULL, or the kfunc absent | the cache step is skipped entirely; `__weak` extern plus `const volatile have_real_data_inode` keeps the object loadable on a kernel without it (lines 134, 237, 496-498) |
| `bpf_task_storage_get` (scratch) | NULL after one retry | `TM_CAUSE_NOBUF` plus `TM_LABEL_UNKNOWN` and an immediate return (lines 295-306, 444-452) |
| `bpf_ringbuf_reserve` | NULL | the label is still raised; only the event is dropped (lines 279-281) - telemetry loss, never a label loss |
| `bpf_get_func_arg` (mmap NULL guard) | non-zero return | the program returns 0 without touching `file` (lines 583-585) |

`raise_label` uses `__sync_fetch_and_or` on the label word and emits an event only
when the value actually changed (lines 274-277), so a concurrent raise from
another hook cannot lose a bit and a repeated open cannot flood the ring buffer.
No program ever returns an error to the kernel: these are observers, and a
non-zero return from an LSM hook would deny the operation.

## What this review found, stated as gaps

- **The probe never denies.** S4 asks for review of a program that only observes.
  The production labeller is supposed to refuse `link`/`rename` out of a ruled root
  (PD-4), and nothing in this file does or can test that - the risk is recorded in
  the founder-decisions ADR and stays open for Phase 1.
- **The directory-rule match is path-based and namespace-relative.** `d_path`
  resolves against the caller's `fs->root`, so a mount-namespace alias or a hard
  link out of a ruled root yields a path no rule matches; scenarios 11 and 12
  reproduce it on both architectures with `label=0x0`. Documented, not fixed.
- **Overlay backing-file events skip directory rules on purpose** (lines 460-465):
  their `f_path` is layer-relative and could match an unrelated rule. The merged
  overlay event is matched normally, which is why scenario 03 raises with
  `cause=open|xattr|backing`.
- **Execute-only opens are over-labelled.** `bpf_get_file_xattr` reads with the
  current task's credentials, so a mode-4111 binary opened for exec returns
  `-EACCES` and raises `TM_LABEL_UNKNOWN`. Accepted as PD-5; the failure direction
  is fail-closed.
- **The inode cache is a floor and it survived review as designed.** It cannot
  lower a label, which is exactly the declassification break scenario 14 measures;
  PD-1 removes the cache from the production labeller, and this review is the
  reason that decision needed measuring rather than arguing.
- **Self-review.** The author checked the program against the three properties the
  checklist names and against the kernel sources cited in its comments, and the
  16 scenarios are the empirical cross-check on those claims. It is still the
  author. Independence is owed.
