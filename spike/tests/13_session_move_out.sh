#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
#
# A2c: a session process cannot move itself out of the session cgroup when the
# scope runs with Delegate=no, and cannot restructure that cgroup. Checklist
# AC3 requires this case; the first 2026-10-18 run did not test it.
#
# Proves:
#   (a) the capability-less uid 65534 principal cannot move itself into a sibling
#       scope's cgroup, and cannot move itself to the cgroupfs root;
#   (b) it cannot create its own child cgroup under the session scope;
#   (c) the session still raises its label afterwards, so the refusals did not
#       break the session's placement (the process is still in the scope);
#   (d) positive control - a root writer CAN write the same cgroup.procs file, so
#       (a) is a permission result and not a broken helper. The threat model
#       relies on this asymmetry: a session that held CAP_SYS_ADMIN could relocate
#       itself, which is why sessions get no privilege (G1) and why doctor checks
#       the delegation.
# Does not prove: that Landlock or seccomp restrict anything (Phase 1 shim).
# Failure maps to the kill gate: keying escape -> labeller design is unsound.

# shellcheck disable=SC1091
source "$(dirname "$0")/_lib.sh"

echo "DESC: A2c an unprivileged session process cannot move itself out of, or restructure, its Delegate=no scope; the session still raises labels afterwards"

f="${SCRATCH}/secret.txt"
echo "canary-placeholder" >"${f}"
label_file "${f}" p

session_start
start_scope "${OTHER_UNIT}" "${OTHER_CG}"

# (a) Move out: sibling scope, then the cgroupfs root. Either refusal is the
#     required result; the errno is printed by the helper and recorded here.
for dst in "${OTHER_CG}" "/sys/fs/cgroup"; do
	out="$(run_in_session unpriv "${PROBE}" try-move --cgroup "${dst}")" \
		|| fail "A2c(a): the unprivileged move attempt against ${dst} crashed the helper"
	echo "A2c(a) ${dst}: ${out}"
	if ! grep -Eq 'errno=(EACCES|EPERM)' <<<"${out}"; then
		fail "A2c(a): unprivileged move to ${dst} was not refused as EACCES/EPERM: ${out}"
	fi
done

# (b) Restructure: no child cgroup of its own under the scope.
out="$(run_in_session unpriv "${PROBE}" try-mkdir --dir "${CG}/escape")" \
	|| fail "A2c(b): the unprivileged mkdir attempt crashed the helper"
echo "A2c(b) ${CG}/escape: ${out}"
if ! grep -Eq 'errno=(EACCES|EPERM)' <<<"${out}"; then
	fail "A2c(b): unprivileged mkdir under the session scope was not refused: ${out}"
fi

# (c) The session is intact: its next labelled open still raises the slot.
reset_slot
run_in_session "${PROBE}" open --path "${f}" --slot "${SLOT}" --expect 0x1

# (d) Positive control: root may write the session scope's own cgroup.procs (the
#     mechanism works), which is exactly the privilege sessions must not have.
run_in_session "${PROBE}" try-move --cgroup "${CG}" --require ok

pass
