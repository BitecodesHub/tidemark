#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
#
# A4b (observation): directory rules are path-based. bpf_path_d_path resolves
# against the caller's root and mount namespace, so a capability-less session
# that can create a user and mount namespace can reach a rule-labelled file
# through an ALIAS path (a bind mount) that no rule matches. The alias step runs
# as uid 65534 with unshare --user --map-root-user --mount, as an unprivileged
# session principal would. Phase 0 draft; built and run on aarch64 on 2026-10-18.
#
# This scenario records the gap; it is NOT a gate and passes unless setup
# fails. Mitigation in the product: sessions cannot create namespaces
# (RestrictNamespaces=yes and a seccomp deny list, guarantee G1), and files that
# matter carry xattr labels, which follow the inode rather than the path. The
# Phase 1 containment matrix tests that namespace creation fails.
# A related gap, also mitigated by G1: inside its own user namespace a session
# holds CAP_SYS_ADMIN over filesystems it mounts there, so it could set
# security.bpf.* on files it created (those files are its own, so this cannot
# lower the label of existing data).

# shellcheck disable=SC1091
source "$(dirname "$0")/_lib.sh"

echo "DESC: A4b observation: a mount-namespace alias path escapes path-based directory rules (gap recorded; mitigated by G1)"

dir="${SCRATCH}/Private"
alias_dir="${SCRATCH}/alias"
mkdir -p "${dir}" "${alias_dir}"
echo "canary-placeholder" >"${dir}/notes.txt" # rule-labelled only, no xattr

session_start
"${PROBE}" set-root --index 0 --prefix "${dir}" --label private

# Control: the direct path is labelled by the rule.
run_in_session "${PROBE}" open --path "${dir}/notes.txt" --slot "${SLOT}" --expect 0x1

# Alias: a private mount namespace binds the labelled directory elsewhere.
reset_slot
# shellcheck disable=SC2016 # expanded by the inner bash, on purpose
if ! run_in_session unpriv unshare --user --map-root-user --mount --propagation private bash -c \
	'mount --bind "$1" "$2" && exec "$3" open --path "$2/notes.txt" --no-label' _ \
	"${dir}" "${alias_dir}" "${PROBE}"; then
	echo "OBSERVE unprivileged user+mount namespace creation failed on this kernel; the alias path could not be built"
	pass
fi
alias_label="$(label_of_slot)"
alias_cause="$(cause_of_slot)"
case "${alias_cause}" in
*dir*) echo "OBSERVE alias path was labelled by a directory rule (label=${alias_label} ${alias_cause}); no gap on this kernel" ;;
*) echo "OBSERVE GAP alias path ${alias_dir}/notes.txt was not labelled by any rule (label=${alias_label} ${alias_cause}); mitigation: deny namespace creation in sessions (G1)" ;;
esac

pass
