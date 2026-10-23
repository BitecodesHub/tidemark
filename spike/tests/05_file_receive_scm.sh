#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
#
# A5: file_receive coverage for descriptors passed with SCM_RIGHTS. Phase 0
# draft; built and run on aarch64 on 2026-10-18.
#
# Proves: when a descriptor to a labelled file is opened OUTSIDE the session and
# passed into a session process over a unix socket (SCM_RIGHTS), the receiving
# in-session process raises the session label via the file_receive hook, and
# the label is already raised when sampled right after recvmsg() returns. The
# sender's own open raises nothing because the sender is not in a session.
# pidfd_getfd also reaches security_file_receive through receive_fd, so it is
# covered by the same hook (not exercised here; production sessions also deny
# it in seccomp).
# Does not prove: an unprivileged receiver of a file it could not open itself
# (scenario 09).
# Failure means fd passing is an unlabelled read path -> narrow the claim.

# shellcheck disable=SC1091
source "$(dirname "$0")/_lib.sh"

echo "DESC: A5 file_receive raises the label when a labelled fd is passed into the session via SCM_RIGHTS"

f="${SCRATCH}/secret.txt"
echo "canary-placeholder" >"${f}"
label_file "${f}" p

session_start

sock="${SCRATCH}/fd.sock"

# Receiver INSIDE the session, in the background; it binds the socket and waits.
run_in_session "${PROBE}" fd-recv --sock "${sock}" --slot "${SLOT}" --expect 0x1 &
recv_pid=$!

for _ in $(seq 1 50); do
	[ -S "${sock}" ] && break
	sleep 0.1
done
[ -S "${sock}" ] || fail "setup: receiver did not bind ${sock}"

# Sender OUTSIDE the session opens the file and passes its fd in.
"${PROBE}" fd-send --path "${f}" --sock "${sock}"

# The receiver's exit status carries the --expect result.
wait "${recv_pid}"
"${PROBE}" readcause --slot "${SLOT}" --has receive,xattr --not open

pass
