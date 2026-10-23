#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
#
# D5: the labelled and unlabelled open() latency baseline. Phase 0 sets no
# target; the numbers are recorded as the input to checklist item C2, which
# decides whether the labeller needs an inode cache at all (the 2026-10-18 run
# showed the cache never contributed to a label).
#
# Measures, with open(O_RDONLY)+close() timed inside the kernel-visible path:
#   1. outside any session (the probe's early exit for non-session processes);
#   2. inside a session on an unlabelled file with no matching rule (the cost the
#      labeller adds to every open of an ordinary file in a session);
#   3. inside a session on an xattr-labelled file, label already raised (steady
#      state: the xattr read still happens on every event);
#   4. inside a session on an xattr-labelled file, slot cleared before EVERY open
#      (the raise path, with the per-iteration label check; zero misses required,
#      which doubles as a self-race at the bench iteration count);
#   5. inside a session under a matching directory rule (the d_path lookup cost).
# Also reports the file's first open after a clear, which is the event that
# raises the label.
# Does not prove: throughput of a real workload (npm install, compiles) - Phase 1
# benchmark methodology.
# Failure maps to the kill gate: only case 4's miss count; slow numbers are a
# design input, not a failure.

# TIMEOUT: 1200

# shellcheck disable=SC1091
source "$(dirname "$0")/_lib.sh"

ITER="${TM_BENCH_ITERATIONS:-20000}"
WARMUP="${TM_BENCH_WARMUP:-2000}"
echo "DESC: D5 open() latency baseline, outside session vs in-session unlabelled vs xattr-labelled vs raise-path vs directory rule (${ITER} iterations each, ${WARMUP} warmup)"

ruledir="${SCRATCH}/benchdir"
mkdir -p "${ruledir}"

clean="${SCRATCH}/bench-clean.txt"
echo "clean" >"${clean}"
labelled="${SCRATCH}/bench-labelled.txt"
echo "secret" >"${labelled}"
ruled="${ruledir}/bench-ruled.txt"
echo "ruled" >"${ruled}"
label_file "${labelled}" p

session_start
"${PROBE}" set-root --index 3 --prefix "${ruledir}" --label secret

echo "BASELINE note: probe loaded, kernel $(uname -r), arch $(uname -m); Phase 0 records numbers and sets no target."

# 1. Outside any session, unlabelled file.
"${PROBE}" bench-open --path "${clean}" --iterations "${ITER}" --warmup "${WARMUP}"

# 2. Inside the session, unlabelled file, no rule.
run_in_session "${PROBE}" bench-open --path "${clean}" --iterations "${ITER}" --warmup "${WARMUP}"

# 3. Inside the session, xattr-labelled file (steady state).
run_in_session "${PROBE}" bench-open --path "${labelled}" --iterations "${ITER}" --warmup "${WARMUP}"

# 4. Inside the session, the raise path: clear the slot, open, check the label.
run_in_session "${PROBE}" bench-open --path "${labelled}" --slot "${SLOT}" \
	--check-raised --clear-each --iterations "${ITER}" --warmup "${WARMUP}"

# 5. Inside the session, a matching directory rule (no xattr on this file).
reset_slot
run_in_session "${PROBE}" bench-open --path "${ruled}" --iterations "${ITER}" --warmup "${WARMUP}"

# The rule still labels the file after 60,000 measured opens.
reset_slot
run_in_session "${PROBE}" open --path "${ruled}" --slot "${SLOT}" --expect 0x2

pass
