#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
#
# A7: setting security.bpf.* fails with EPERM for a capability-less process
# inside a session. Phase 0 draft; built and run on aarch64 on 2026-10-18.
#
# Proves: an in-session process without CAP_SYS_ADMIN (as a real Tidemark
# session principal runs) cannot set or overwrite security.bpf.tidemark, so a
# session cannot forge or clear its own labels. A positive control, run as the
# SAME unprivileged uid on the SAME file, sets a user.* xattr successfully, so
# the EPERM can only come from the capability check (cap_inode_setxattr), not
# from path lookup or write permission. A root control outside the session
# must succeed.
# Emulation note: the spike's scope runs as root, so the test drops to uid
# 65534 with no capabilities via setpriv to stand in for the session principal.
# Does not prove: protection against a session that mounts its own filesystem
# inside a user namespace, where cap_inode_setxattr checks CAP_SYS_ADMIN in
# that namespace (production sessions deny namespace creation; see scenario 11
# and the README); cache invalidation on a privileged setxattr (Phase 1).
# Failure means sessions could rewrite labels -> the label integrity claim fails.

# shellcheck disable=SC1091
source "$(dirname "$0")/_lib.sh"

echo "DESC: A7 security.bpf.* setxattr is EPERM for a capability-less in-session process (positive control: user.* succeeds)"

f="${SCRATCH}/target.txt"
echo "canary-placeholder" >"${f}"
chmod 0666 "${f}" # write access is not the barrier; only the capability is

session_start

# Positive control: same uid, same file, a user.* xattr must succeed.
run_in_session unpriv "${PROBE}" setxattr --path "${f}" --name user.tm_probe --value x --require ok

# The claim: security.bpf.* must be refused with EPERM.
run_in_session unpriv "${PROBE}" setxattr --path "${f}" --name security.bpf.tidemark --value p --require EPERM

# Root control outside the session must succeed.
"${PROBE}" setxattr --path "${f}" --name security.bpf.tidemark --value p --require ok

pass
