# SPDX-License-Identifier: Apache-2.0
# shellcheck shell=bash
#
# Shared helpers for the Tidemark Phase 0 kernel-spike scenario scripts.
# This file is sourced, never run directly; the orchestrator skips files whose
# name starts with "_". Phase 0 draft; built and run on aarch64 on 2026-10-18.
#
# A "session" in the spike is a transient systemd scope under a dash-free slice
# (tmspike.slice), so the scope cgroup sits at a fixed depth of 2 below
# /sys/fs/cgroup (root=0, tmspike.slice=1, scope=2). A dash-free slice name is
# deliberate: a slice named "a-b.slice" would nest under "a.slice" and change
# the depth. session_level in the loader must match this depth (2), and
# register-session refuses a cgroup at any other depth.

set -euo pipefail

PROBE="${TM_PROBE_BIN:?TM_PROBE_BIN not set}"
SLOT="${TM_SLOT:-1}"
SLICE="tmspike.slice"
UNIT="tm-sess-${SLOT}.scope"
CG="/sys/fs/cgroup/${SLICE}/${UNIT}"
OTHER_UNIT="tm-other-${SLOT}.scope"
# shellcheck disable=SC2034 # used by the scenarios that source this file
OTHER_CG="/sys/fs/cgroup/${SLICE}/${OTHER_UNIT}"
MAX_ROOTS=8

# mktemp -d creates a 0700 root-owned directory; make it traversable so the
# unprivileged (uid 65534) steps in scenarios 07, 09 and 11 reach their files.
SCRATCH="$(mktemp -d /var/tmp/tm-spike.XXXXXX)"
chmod 0755 "${SCRATCH}"

HOLDERS=()
MOUNTS=()
SUBCGS=()

cleanup_common() {
	local m c h
	for m in "${MOUNTS[@]:-}"; do
		if [ -n "${m}" ]; then
			umount "${m}" >/dev/null 2>&1 || true
		fi
	done
	# Child cgroups created by a scenario must be emptied and removed, or the
	# scope's own cgroup cannot be removed when systemd stops it.
	for c in "${SUBCGS[@]:-}"; do
		if [ -n "${c}" ] && [ -d "${c}" ]; then
			while read -r p; do kill -9 "${p}" >/dev/null 2>&1 || true; done <"${c}/cgroup.procs"
			rmdir "${c}" >/dev/null 2>&1 || true
		fi
	done
	systemctl stop "${UNIT}" "${OTHER_UNIT}" >/dev/null 2>&1 || true
	for h in "${HOLDERS[@]:-}"; do
		if [ -n "${h}" ]; then
			kill "${h}" >/dev/null 2>&1 || true
		fi
	done
	rm -rf "${SCRATCH}" >/dev/null 2>&1 || true
}
trap cleanup_common EXIT
trap 'echo "RESULT: FAIL (line ${LINENO})"' ERR

# fail <message>: print a FAIL result and exit non-zero.
fail() {
	echo "RESULT: FAIL: $*"
	exit 1
}

# pass: print a PASS result and exit zero.
pass() {
	echo "RESULT: PASS"
	exit 0
}

# start_scope <unit> <cgroup-dir>: start a transient scope holding a sleeper
# and wait until its cgroup exists.
start_scope() {
	local unit="$1" cg="$2" i
	systemd-run --quiet --collect --slice="${SLICE}" --unit="${unit}" \
		--scope -p Delegate=no -- sleep 600 &
	HOLDERS+=("$!")
	for i in $(seq 1 50); do
		[ -d "${cg}" ] && break
		sleep 0.1
	done
	[ -d "${cg}" ] || fail "setup: scope cgroup ${cg} was not created"
}

# clear_roots: empty every directory-rule slot. The probe stays loaded for the
# whole suite, so rules from an earlier scenario must not leak into this one.
clear_roots() {
	local i
	for i in $(seq 0 $((MAX_ROOTS - 1))); do
		"${PROBE}" set-root --index "${i}" --clear >/dev/null
	done
}

# session_start: create the session scope, register it (which clears this
# slot's label and cause words), and clear every directory rule.
session_start() {
	start_scope "${UNIT}" "${CG}"
	"${PROBE}" register-session --cgroup "${CG}" --slot "${SLOT}"
	clear_roots
}

# reset_slot: clear this slot's label and cause words between steps.
reset_slot() {
	"${PROBE}" register-session --cgroup "${CG}" --slot "${SLOT}" >/dev/null
}

# run_in_cgroup <cgroup-dir> <cmd...>: run a command as a GRANDCHILD of a
# process that joined <cgroup-dir>. The moved subshell M forks a bash C0 (not
# the last command, so the fork cannot be elided), and C0 forks the command C1
# (again not the last command). The helper's checkSessionContext asserts
# pid != M, ppid != M and grandparent == M, and that it runs in <cgroup-dir>.
run_in_cgroup() {
	local cg="$1"
	shift
	(
		echo "${BASHPID}" >"${cg}/cgroup.procs"
		export TM_MOVED_PID="${BASHPID}" TM_SESSION_CG="${cg}"
		bash -c '"$@"; rc=$?; exit "${rc}"' _ "$@"
		rc=$?
		exit "${rc}"
	)
}

# run_in_session <cmd...>: run_in_cgroup on this scenario's session scope.
run_in_session() {
	run_in_cgroup "${CG}" "$@"
}

# label_file <path> <value>: set security.bpf.tidemark as root OUTSIDE any
# session and verify it was stored, so a labelling failure is reported as a
# setup error rather than blamed on the kernel hook.
label_file() {
	local f="$1" v="$2" got
	"${PROBE}" setxattr --path "${f}" --name security.bpf.tidemark --value "${v}" --require ok
	got="$(getfattr --only-values -n security.bpf.tidemark "${f}" 2>/dev/null || true)"
	[ "${got}" = "${v}" ] || fail "setup: security.bpf.tidemark on ${f} reads '${got}', expected '${v}'"
}

# label_of_slot: print the current label word (hex) for this session's slot.
label_of_slot() {
	"${PROBE}" readlabel --slot "${SLOT}" --raw | sed -n 's/^label=//p'
}

# cause_of_slot: print the current cause word for this session's slot.
cause_of_slot() {
	"${PROBE}" readcause --slot "${SLOT}" | sed -n 's/^readcause //p'
}

# unpriv <cmd...>: run a command as uid/gid 65534 with no capabilities and
# no_new_privs, standing in for a real (capability-less) session principal.
unpriv() {
	setpriv --reuid=65534 --regid=65534 --clear-groups \
		--no-new-privs --bounding-set=-all -- "$@"
}
# Exported so that the bash child inside run_in_cgroup can call it; setpriv
# then execs the helper, so the helper keeps the grandchild position.
export -f unpriv
