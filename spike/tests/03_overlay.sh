#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
#
# A3: behaviour through overlayfs. Phase 0 draft; built and run on aarch64 on 2026-10-18.
#
# Gate: a lower file with NO xattr, opened THROUGH an overlay mount, raises the
# session label through a directory rule on the overlay path (bpf_path_d_path),
# and the cause word shows the directory rule (and no xattr), so the gate
# exercises the directory-rule mechanism in isolation.
# Observation (for the labeller ADR, not a gate): an xattr-labelled lower file
# opened through the overlay with NO matching rule. At v7.2 ovl_open opens the
# real (backing) file through backing_file_open, which runs security_file_open
# a second time with f_path on the real lower path, so the xattr can be read on
# that second event (cause bit "backing") even if the overlay inode exposes
# nothing; whether the data-inode cache contributes is reported too.
# Failure of the gate maps to the kill gate: "unreliable ... through overlayfs"
# -> ship audit-only or stop, or fall back to the fanotify vault with its
# documented gaps.

# shellcheck disable=SC1091
source "$(dirname "$0")/_lib.sh"

echo "DESC: A3 overlayfs: an unlabelled lower file opened via the overlay is labelled by a directory rule (gated); xattr path observed"

lower="${SCRATCH}/lower"
upper="${SCRATCH}/upper"
work="${SCRATCH}/work"
merged="${SCRATCH}/merged"
mkdir -p "${lower}" "${upper}" "${work}" "${merged}"

echo "canary-placeholder" >"${lower}/plain.txt"  # deliberately no xattr
echo "canary-placeholder" >"${lower}/doc.txt"
label_file "${lower}/doc.txt" p

mount -t overlay overlay \
	-o "lowerdir=${lower},upperdir=${upper},workdir=${work}" "${merged}"
MOUNTS+=("${merged}")

session_start

# Gate: directory rule on the merged path, file without xattr.
"${PROBE}" set-root --index 0 --prefix "${merged}" --label private
run_in_session "${PROBE}" open --path "${merged}/plain.txt" --slot "${SLOT}" --expect 0x1
"${PROBE}" readcause --slot "${SLOT}" --has open,dir --not xattr

# Observation: xattr-labelled lower file through the overlay, no matching rule.
reset_slot
clear_roots
run_in_session "${PROBE}" open --path "${merged}/doc.txt" --slot "${SLOT}"
echo "OBSERVE overlay xattr path: label=$(label_of_slot) $(cause_of_slot)"
echo "OBSERVE (label 0x1 with cause xattr|backing means the backing-file open read the lower xattr; cause cache means the data-inode cache contributed)"

pass
