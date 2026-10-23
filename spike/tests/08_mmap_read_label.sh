#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
#
# A8 (end to end): opening a labelled file inside the session and reading it
# through a PROT_READ mapping ends with the label raised. Phase 0 draft,
# Built and run on aarch64 on 2026-10-18.
#
# Proves: the open-then-map path leaves the label raised. The open alone is
# caught by file_open; this scenario does NOT isolate the mmap_file hook.
# Scenario 10 isolates mmap_file with a descriptor opened before the process
# joined the session.
# Does not prove: coverage of read() on a descriptor inherited across a cgroup
# move or fork/exec (scenario 10 records that gap).
# Failure means the basic open-and-map read path is unlabelled -> kill gate.

# shellcheck disable=SC1091
source "$(dirname "$0")/_lib.sh"

echo "DESC: A8 an open-then-mmap read of a labelled file leaves the label raised"

f="${SCRATCH}/secret.txt"
echo "canary-placeholder-with-enough-bytes-to-map" >"${f}"
label_file "${f}" p

session_start
run_in_session "${PROBE}" mmapread --path "${f}" --slot "${SLOT}" --expect 0x1
"${PROBE}" readcause --slot "${SLOT}" --has open,xattr

pass
