#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
#
# A2: sessions are identified by the ancestor cgroup id at a fixed depth (2),
# with Delegate=no. Phase 0 draft; built and run on aarch64 on 2026-10-18.
#
# Proves:
#   (a) an open by a process OUTSIDE any session does not raise the slot;
#   (b) an open by a grandchild INSIDE the session scope raises it;
#   (c) a process in a CHILD cgroup below the scope still counts as the session,
#       because bpf_get_current_ancestor_cgroup_id(2) returns the level-2
#       ancestor however deep the process sits;
#   (d) a process in a SIBLING scope that is not registered does not raise it.
# register-session refuses a cgroup whose depth differs from session_level, so
# a scope created at the wrong level is a setup error, not a silent miss.
# Does not prove: correctness of the label value beyond raised/not-raised.
# Failure maps to the kill gate: unreliable kernel keying -> audit-only or stop.

# shellcheck disable=SC1091
source "$(dirname "$0")/_lib.sh"

echo "DESC: A2 session keyed by level-2 ancestor cgroup id; outside and sibling opens do not raise, child-cgroup opens do"

f="${SCRATCH}/secret.txt"
echo "canary-placeholder" >"${f}"
label_file "${f}" p

session_start

# (a) Outside any session: a plain child of this script.
"${PROBE}" open --path "${f}" --slot "${SLOT}" --expect 0x0

# (b) Inside the session scope.
run_in_session "${PROBE}" open --path "${f}" --slot "${SLOT}" --expect 0x1

# (c) Child cgroup below the scope (root may create it; the scope has no
#     controllers enabled, so moving a process into the child is allowed).
reset_slot
sub="${CG}/sub"
mkdir "${sub}"
SUBCGS+=("${sub}")
run_in_cgroup "${sub}" "${PROBE}" open --path "${f}" --slot "${SLOT}" --expect 0x1

# (d) Sibling scope at the same depth that is NOT registered.
reset_slot
start_scope "${OTHER_UNIT}" "${OTHER_CG}"
run_in_cgroup "${OTHER_CG}" "${PROBE}" open --path "${f}" --slot "${SLOT}" --expect 0x0

pass
