#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
#
# A8b: a process opens a labelled file OUTSIDE the session, joins the session
# while holding the descriptor, then (i) read()s it and (ii) mmap()s it.
# Phase 0 draft; built and run on aarch64 on 2026-10-18.
#
# Proves (gate): mmap_file raises the label for a mapping of a descriptor that
# file_open never saw inside the session, i.e. the mmap hook in isolation.
# Records (not a gate): read() on such a descriptor fires no hook this probe
# attaches (file_permission is not sleepable and is not used), so the label
# stays 0x0. This is a known gap: production sessions must start from a fresh
# exec with no inherited descriptors (the tidemark-shim calls close_range
# before exec), which the Phase 1 containment matrix tests.
# If TM_ATTACH_MMAP_FILE=0 the mmap result is reported but not gated.
# Failure of the gate means mapped reads of pre-session descriptors are
# unlabelled -> the read-coverage claim must exclude them.

# shellcheck disable=SC1091
source "$(dirname "$0")/_lib.sh"

echo "DESC: A8b mmap of a pre-session descriptor raises the label (gated); read() on it is a recorded gap"

f="${SCRATCH}/secret.txt"
echo "canary-placeholder-with-enough-bytes-to-map" >"${f}"
label_file "${f}" p

session_start

extra=()
if [ "${TM_ATTACH_MMAP_FILE:-1}" = "0" ]; then
	echo "OBSERVE lsm.s/mmap_file is detached (TM_ATTACH_MMAP_FILE=0); the mmap gate is skipped"
	extra=(--skip-mmap-gate)
fi

"${PROBE}" join-then-use --path "${f}" --cgroup "${CG}" --slot "${SLOT}" "${extra[@]}"
if [ "${#extra[@]}" -eq 0 ]; then
	"${PROBE}" readcause --slot "${SLOT}" --has mmap,xattr --not open
fi

pass
