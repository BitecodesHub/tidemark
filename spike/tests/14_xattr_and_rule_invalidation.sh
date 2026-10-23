#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
#
# A12: label invalidation. Checklist D4 lists "xattr invalidation" as a test that
# the 2026-10-18 draft did not include, and the labeller ADR kill test requires
# "Reload a policy that adds a rule: new opens are labelled".
#
# The inode cache (cause bit `cache`) holds XATTR-derived bits only, acts as a
# FLOOR that is ORed on top of the fresh xattr read, and is never invalidated on
# setxattr or removexattr (that is Phase 1 work). Therefore:
#   (a) a first labelled open raises the label by the xattr;
#   (b) raising the attribute from 'p' to 's' on the SAME inode is picked up by
#       the very next open - a stricter label can never be hidden by the cache
#       (this is the gate the ADR's floor property claims);
#   (c) OBSERVE removing the attribute: whether the cached bit still raises the
#       label on the next open. Recorded, not gated: this is the Phase 1
#       invalidation work (inode_post_setxattr / inode_post_removexattr
#       generation counter), and the no-cache design decision in C2 removes the
#       problem entirely if C2 selects it;
#   (d) a directory rule added AFTER a file was created labels that file's next
#       open (rule-derived bits are never cached);
#   (e) clearing that rule un-labels the next open of the same file, with the
#       label returning to exactly public - a rule change is never stale.
# Failure maps to the kill gate: (b), (d) or (e) failing means policy changes are
# not honoured, which is a labeller design failure.

# shellcheck disable=SC1091
source "$(dirname "$0")/_lib.sh"

echo "DESC: A12 a raised attribute is honoured on the next open (cache is a floor, never a hide); added and removed directory rules take effect on the next open; attribute removal is recorded as the Phase 1 invalidation gap"

label_bits() { # hex label -> bit test helper
	local raw="$1"
	echo $((raw))
}

dir1="${SCRATCH}/ruledir"
mkdir -p "${dir1}"

f1="${SCRATCH}/raise.txt"
echo "canary-1" >"${f1}"
f2="${dir1}/rule.txt"
echo "canary-2" >"${f2}"

session_start

# (a) attribute 'p' -> private on the first open.
label_file "${f1}" p
reset_slot
run_in_session "${PROBE}" open --path "${f1}" --slot "${SLOT}" --expect 0x1
"${PROBE}" readcause --slot "${SLOT}" --has xattr

# (b) raise the same inode to 's'; the next open must show the secret bit.
label_file "${f1}" s
reset_slot
run_in_session "${PROBE}" open --path "${f1}" --slot "${SLOT}"
raw="$(label_of_slot)"
got="$(label_bits "${raw}")"
echo "A12(b) after raising p to s: label=${raw}"
if [ "$((got & 0x2))" != "2" ]; then
	fail "A12(b): raising the attribute to 's' was not honoured on the next open (label=${raw}; the cache must never hide a stricter label)"
fi
"${PROBE}" readcause --slot "${SLOT}" --has xattr

# (c) OBSERVE removal: is the cached bit still raised?
"${PROBE}" removexattr --path "${f1}" --name security.bpf.tidemark
reset_slot
run_in_session "${PROBE}" open --path "${f1}" --slot "${SLOT}"
echo "A12(c) OBSERVE after removing the attribute: label=$(label_of_slot) cause=$(cause_of_slot) (a raised label here is the documented Phase 1 invalidation gap; 0x0 means the cache contributed nothing for this inode)"

# (d) add a directory rule over a file that already exists, then open it.
reset_slot
run_in_session "${PROBE}" open --path "${f2}" --slot "${SLOT}" --expect 0x0
"${PROBE}" set-root --index 2 --prefix "${dir1}" --label secret
reset_slot
run_in_session "${PROBE}" open --path "${f2}" --slot "${SLOT}" --expect 0x2
"${PROBE}" readcause --slot "${SLOT}" --has dir --not xattr

# (e) remove the rule: the next open is unlabelled again.
"${PROBE}" set-root --index 2 --clear
reset_slot
run_in_session "${PROBE}" open --path "${f2}" --slot "${SLOT}" --expect 0x0

pass
