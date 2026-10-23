#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
#
# A1r: the labeller ADR kill test "where timing matters" - 10,000 iterations with
# ZERO misses. Checklist AC1 and threat-model S7 checked this once, on one open;
# the ADR requires the repeated form, and the first 2026-10-18 run did not have
# the observer process it also names.
#
# Proves, once per iteration (default 10,000, override with TM_RACE_ITERATIONS):
#   (a) SELF - after clearing the session slot, a labelled open() by a grandchild
#       that joined the session leaves the label raised when that same process
#       reads it immediately after open() returns; and
#   (b) OBSERVER - a SEPARATE process outside the session, standing in for
#       tidemark-gw, reads the label immediately after receiving the byte the
#       opener writes after its open() returned, and sees the same raise.
# Each iteration clears the slot first, so the raise under test is always caused
# by that iteration's own open(), never by an earlier one.
# Does not prove: anything about the gateway's own clearance decisions (Phase 1).
# Failure maps to the kill gate: a miss is unreliable label raising -> criterion 1.

# TIMEOUT: 1800

# shellcheck disable=SC1091
source "$(dirname "$0")/_lib.sh"

ITER="${TM_RACE_ITERATIONS:-10000}"
echo "DESC: A1r ${ITER} iterations, zero misses: label raised before open() returns, seen by the opener and by a separate observer process"

f="${SCRATCH}/race.txt"
echo "canary-race" >"${f}"
label_file "${f}" p

session_start

sock="${SCRATCH}/race.sock"
obs_out="${SCRATCH}/observer.out"

# The observer runs outside any session, as root, for the whole loop.
"${PROBE}" race-observe --sock "${sock}" --slot "${SLOT}" \
	--iterations "${ITER}" --expect 0x1 >"${obs_out}" 2>&1 &
obs=$!
HOLDERS+=("${obs}")

run_in_session "${PROBE}" race-open --path "${f}" --slot "${SLOT}" \
	--iterations "${ITER}" --expect 0x1 --notify "${sock}"

if ! wait "${obs}"; then
	cat "${obs_out}"
	fail "A1r(b): the observer process reported a miss or exited early"
fi
cat "${obs_out}"
"${PROBE}" readcause --slot "${SLOT}" --has open,xattr

pass
