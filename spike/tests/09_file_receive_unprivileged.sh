#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
#
# A5b: an UNPRIVILEGED in-session receiver gets a descriptor to a root-owned
# 0600 labelled file over SCM_RIGHTS. Phase 0 draft; built and run on aarch64 on 2026-10-18.
#
# Why: at v7.2 bpf_get_file_xattr checks inode_permission(MAY_READ) with the
# CURRENT task's credentials for security.bpf.* names, so the receiver cannot
# read the xattr of a file it could not open itself, yet it can read() the
# received descriptor. The probe must fail CLOSED: the label must still rise
# (TM_LABEL_UNKNOWN, 0x4, cause "error|eacces"), never stay public.
# Proves: the fail-closed path for the file_receive hook under an unprivileged
# principal.
# Does not prove: which exact production label an unknown read should map to
# (the clearance model treats unknown as at least secret).
# Failure means a capability-less session can read labelled data received over
# a socket without any label -> narrow the claim and redesign the xattr read
# (for example read the label from a root-owned map, not the xattr, at receive).

# shellcheck disable=SC1091
source "$(dirname "$0")/_lib.sh"

echo "DESC: A5b unprivileged receiver of a root-only labelled fd: label still rises (fail-closed)"

f="${SCRATCH}/rootonly.txt"
echo "canary-placeholder" >"${f}"
chmod 0600 "${f}"
label_file "${f}" p

rx="${SCRATCH}/rx"
mkdir "${rx}"
chmod 1777 "${rx}" # the unprivileged receiver must be able to bind here
sock="${rx}/fd.sock"

session_start

run_in_session unpriv "${PROBE}" fd-recv --sock "${sock}" --no-label &
recv_pid=$!
for _ in $(seq 1 50); do
	[ -S "${sock}" ] && break
	sleep 0.1
done
[ -S "${sock}" ] || fail "setup: unprivileged receiver did not bind ${sock}"

"${PROBE}" fd-send --path "${f}" --sock "${sock}"
wait "${recv_pid}"

# Root reads the slot after the receiver exits. The exact value and causes
# show the fail-closed EACCES mechanism, not merely a raised label.
"${PROBE}" readlabel --slot "${SLOT}" --expect 0x4
"${PROBE}" readcause --slot "${SLOT}" --has receive,error,eacces --not xattr

pass
