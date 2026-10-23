#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
#
# A1: a grandchild process opens an xattr-labelled file, and the session label
# is already raised when read in the same process immediately after open()
# returns, before any read(). Phase 0 draft; built and run on aarch64 on 2026-10-18.
#
# Proves: the file_open sleepable LSM program raises the label synchronously,
# before open() returns, for a process two forks below the process that joined
# the session (the helper asserts its lineage and cgroup before acting).
# Does not prove: anything about overlayfs, directory rules or fd passing.
# Failure maps to the Phase 0 kill gate: "label raising in the kernel is
# unreliable on stock Fedora 45" -> ship audit-only or stop.

# shellcheck disable=SC1091
source "$(dirname "$0")/_lib.sh"

echo "DESC: A1 label rises before open() returns for a grandchild opening an xattr-labelled file"

f="${SCRATCH}/secret.txt"
echo "canary-placeholder" >"${f}"
label_file "${f}" p

session_start

# Grandchild open; the helper samples the label right after open() returns.
run_in_session "${PROBE}" open --path "${f}" --slot "${SLOT}" --expect 0x1
"${PROBE}" readcause --slot "${SLOT}" --has open,xattr --not error

pass
