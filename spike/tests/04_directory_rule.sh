#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
#
# A4: directory rules via bpf_path_d_path prefix matching. Phase 0 draft; built and run on aarch64 on 2026-10-18.
#
# Proves:
#   (a) a rule matches only on a path-component boundary: a file under
#       <root>X/ (a sibling whose name starts with the root's name) is NOT
#       labelled;
#   (b) a file with NO xattr under a labelled root raises the label, with the
#       directory rule as the cause;
#   (c) overlapping rules combine: a file under <root>/keys, with <root> private
#       and <root>/keys secret, gets private|secret (0x3).
# Does not prove: rule matching for aliases created in another mount namespace
# (see scenario 11), or anything about xattrs or overlayfs.
# Failure maps to the kill gate for kernel label raising; it also means default
# sensitive folders could not be kept labelled when a tool drops xattrs.

# shellcheck disable=SC1091
source "$(dirname "$0")/_lib.sh"

echo "DESC: A4 directory rules: component-boundary match, unlabelled file labelled by rule, overlapping rules combine"

dir="${SCRATCH}/Private"
mkdir -p "${dir}/keys" "${dir}X"
echo "canary-placeholder" >"${dir}/notes.txt"   # deliberately no xattr
echo "canary-placeholder" >"${dir}X/notes.txt"
echo "canary-placeholder" >"${dir}/keys/k.txt"

session_start
"${PROBE}" set-root --index 0 --prefix "${dir}" --label private
"${PROBE}" set-root --index 1 --prefix "${dir}/keys" --label secret

# (a) Boundary: <root>X must not match <root>.
run_in_session "${PROBE}" open --path "${dir}X/notes.txt" --slot "${SLOT}" --expect 0x0

# (b) Unlabelled file under the root.
run_in_session "${PROBE}" open --path "${dir}/notes.txt" --slot "${SLOT}" --expect 0x1
"${PROBE}" readcause --slot "${SLOT}" --has open,dir --not xattr,error

# (c) Overlapping rules.
reset_slot
run_in_session "${PROBE}" open --path "${dir}/keys/k.txt" --slot "${SLOT}" --expect 0x3

pass
