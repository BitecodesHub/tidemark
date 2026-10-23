#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
#
# A6: user space reads the per-session label word through a read-only mmap of a
# BPF_F_MMAPABLE array. Phase 0 draft; built and run on aarch64 on 2026-10-18.
#
# Proves: a PROT_READ, MAP_SHARED mapping of the label array returns the raised
# value; the KERNEL refuses to make that mapping writable (mprotect fails); and
# a BPF_F_RDONLY descriptor refuses a writable shared mapping. It also OBSERVES
# whether a BPF_F_RDONLY descriptor can be mapped read-only at all (expected
# not to, at v7.2), which decides how Phase 1 hands the label word to the
# unprivileged gateway and broker.
# Does not prove: performance, or mapping from a non-root daemon UID (Phase 1).
# Failure means the exits cannot cheaply and safely observe the label ->
# redesign the label transport. This is NOT a kill condition for the label
# raise itself.

# shellcheck disable=SC1091
source "$(dirname "$0")/_lib.sh"

echo "DESC: A6 the label array maps read-only, cannot be made writable, and reads the raised label"

f="${SCRATCH}/secret.txt"
echo "canary-placeholder" >"${f}"
label_file "${f}" p

session_start
run_in_session "${PROBE}" open --path "${f}" --slot "${SLOT}" --expect 0x1

"${PROBE}" check-ro-mmap --slot "${SLOT}"

lbl="$(label_of_slot)"
[ "${lbl}" = "0x1" ] || fail "read-only mmap returned label ${lbl}, expected 0x1"

pass
