#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
#
# A4c (observation): a hard link or a rename OUT of a labelled root, made
# inside the session's own mount namespace (no new namespace needed), gives the
# file a path that no directory rule matches. d_path names the dentry that was
# opened, linkat and renameat fire no hook this probe attaches, and rule-derived
# labels are not cached per inode (labeller ADR). Phase 0 draft; built and run on aarch64 on 2026-10-18.
#
# This scenario records the gap; it is NOT a gate and passes unless setup
# fails. Open founder decision (labeller ADR): should a directory label follow
# the location or the content? If the content, Phase 1 must bind labels to the
# inode, for example by stamping the xattr from inode_link and path_rename
# hooks, or deny link and rename out of labelled roots. Files that carry an
# xattr label are not affected, because the xattr follows the inode.

# shellcheck disable=SC1091
source "$(dirname "$0")/_lib.sh"

echo "DESC: A4c observation: a hard link or rename out of a labelled root escapes path-based rules (gap recorded)"

dir="${SCRATCH}/Private"
mkdir -p "${dir}"
echo "canary-placeholder" >"${dir}/linkme.txt"  # rule-labelled only, no xattr
echo "canary-placeholder" >"${dir}/moveme.txt"  # rule-labelled only, no xattr

session_start
"${PROBE}" set-root --index 0 --prefix "${dir}" --label private

observe() {
	local what="$1" label cause
	label="$(label_of_slot)"
	cause="$(cause_of_slot)"
	case "${cause}" in
	*dir*) echo "OBSERVE ${what}: still labelled by a rule (label=${label} ${cause}); no gap" ;;
	*) echo "OBSERVE GAP ${what}: not labelled by any rule (label=${label} ${cause})" ;;
	esac
}

# Hard link made inside the session, then opened through the link.
run_in_session ln "${dir}/linkme.txt" "${SCRATCH}/linked.txt"
run_in_session "${PROBE}" open --path "${SCRATCH}/linked.txt" --slot "${SLOT}"
observe "hard link ${SCRATCH}/linked.txt"

# Rename made inside the session, then opened at the new path.
reset_slot
run_in_session mv "${dir}/moveme.txt" "${SCRATCH}/moved.txt"
run_in_session "${PROBE}" open --path "${SCRATCH}/moved.txt" --slot "${SLOT}"
observe "rename to ${SCRATCH}/moved.txt"

pass
