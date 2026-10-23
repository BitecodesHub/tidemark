#!/usr/bin/bash
# SPDX-License-Identifier: Apache-2.0
#
# AgentOS Tidemark boot self-test.
# Status: Phase 0 draft (2026-10-16, revised after review). Not a product
# component. It has not yet run on Linux.
#
# What it does
#   Reports whether the running system provides the kernel and filesystem
#   features that later Tidemark phases depend on. It enforces nothing: it does
#   not start, block, label or confine any session, and it sends nothing over
#   the network. In Phase 1, tidemarkd is expected to refuse to start sessions
#   when this self-test fails; that behaviour does not exist yet.
#
# Output
#   Human-readable lines prefixed "tidemark-selftest:", then exactly one
#   machine-readable line, always the last line, of the form
#     TIDEMARK_SELFTEST v=1 result=PASS phase=0 failed=none uptime_s=... kernel=... key=value ...
#   Values never contain spaces. The same key=value pairs are also written,
#   one per line, to /run/tidemark/selftest.result.
#
# Required checks (any failure gives result=FAIL and exit status 1)
#   kernel_floor  uname -r is Linux 7.2 or later. Landlock ABI 9 first ships in
#                 Linux 7.1 (security/landlock/syscalls.c at tag v7.1). The 7.2
#                 floor is a policy choice, not a kfunc dependency: the spike
#                 measured that the label is correct with the inode-keyed cache
#                 compiled out, so the labeller keeps no cache
#                 (docs/adr/2026-10-23_Phase0-Founder-Decisions_v1.md, PD-1) and
#                 no longer calls bpf_real_data_inode. 7.2 stays the floor
#                 because Phase 0 has no evidence on 7.1 and the project does
#                 not claim untested support; 7.1 is Phase 6 (portable mode).
#                 Fedora 44 and 45 both ship 7.2.x.
#   lsm_bpf       /sys/kernel/security/lsm lists bpf
#   lsm_landlock  /sys/kernel/security/lsm lists landlock
#   landlock      landlock_create_ruleset(NULL, 0, LANDLOCK_CREATE_RULESET_VERSION)
#                 returns 9 or more (raw syscall 444 through Python ctypes,
#                 because the kernel exposes the ABI version only through
#                 this query; 444 is the number on x86_64 and on the
#                 asm-generic table used by aarch64)
#   cgroup2       /sys/fs/cgroup is the unified cgroup v2 hierarchy
#   overlayfs     an overlay mount inside a private mount namespace, with its
#                 upper and work directories on /var, accepts a write that
#                 lands in the upper directory
#   var_xattr     a user.* extended attribute can be set and read back on /var
#   var_append    FS_APPEND_FL (chattr +a) can be set on a file on /var, a
#                 truncating open is then refused, and appending still works
#   btf           /sys/kernel/btf/vmlinux is readable
#   kfunc_bpf_get_file_xattr, kfunc_bpf_path_d_path
#                 each name is present as an exact string in
#                 /sys/kernel/btf/vmlinux. This is a proxy only: it shows that
#                 the kernel carries BTF for the function, not that a sleepable
#                 BPF LSM program may call it. Confirming each kfunc by loading
#                 a program is the job of the Phase 0 kernel probe and, from
#                 Phase 1, of tidemarkd.
#
# Informational checks (reported only; they never change the result in Phase 0)
#   uptime_s (seconds since boot when the self-test started, from
#   /proc/uptime), arch, lsm, landlock_abi, cgroup2_controllers,
#   cgroup2_freeze, var_fstype, xattr_security_bpf (set and read back a
#   security.bpf.* attribute as root), selinux.
#
# Environment
#   TIDEMARK_SELFTEST_DIR          directory on the filesystem under test (default
#                                  /var/lib/tidemark/selftest; must be absolute). The
#                                  probes run in a private run.XXXXXX subdirectory
#                                  that is removed afterwards.
#   TIDEMARK_SELFTEST_RESULT_FILE  result file (default /run/tidemark/selftest.result)
#
# Exit status: 0 PASS, 1 FAIL, 2 usage error.
# Root is required for the overlayfs and var_append checks; without root they
# are reported as failed with the reason "requires_root".

set -u
export LC_ALL=C
PATH=/usr/sbin:/usr/bin:/sbin:/bin
export PATH

readonly SCHEMA_VERSION=1
readonly MIN_KERNEL_MAJOR=7
readonly MIN_KERNEL_MINOR=2
readonly MIN_LANDLOCK_ABI=9
readonly SYS_LANDLOCK_CREATE_RULESET=444
readonly -a REQUIRED_KFUNCS=(bpf_get_file_xattr bpf_path_d_path)

state_dir="${TIDEMARK_SELFTEST_DIR:-/var/lib/tidemark/selftest}"
result_file="${TIDEMARK_SELFTEST_RESULT_FILE:-/run/tidemark/selftest.result}"

declare -a keys=()
declare -A values=()
declare -a failed=()
append_probe=""
overlay_base=""
work_dir=""

usage() {
  cat <<'EOF'
Usage: tidemark-selftest.sh [--help]

Runs the AgentOS Tidemark Phase 0 boot self-test and prints one line that
starts with the token TIDEMARK_SELFTEST. See the header of this script for
the list of checks and the environment variables it reads.
EOF
}

log() {
  printf 'tidemark-selftest: %s\n' "$*"
}

# record KEY VALUE: store a key=value pair for the result line.
record() {
  local key=$1 value=$2
  value=${value//[[:space:]]/_}
  [[ -n $value ]] || value=empty
  if [[ -z ${values[$key]+set} ]]; then
    keys+=("$key")
  fi
  values[$key]=$value
}

# require KEY STATUS: STATUS is "ok" or a short failure reason.
require() {
  record "$1" "$2"
  if [[ $2 == ok ]]; then
    log "ok   $1"
  else
    failed+=("$1")
    log "FAIL $1: $2"
  fi
}

# info KEY VALUE: informational only.
info() {
  record "$1" "$2"
  log "info $1: $2"
}

have() {
  command -v "$1" >/dev/null 2>&1
}

is_root() {
  [[ ${EUID} -eq 0 ]]
}

cleanup() {
  if [[ -n $append_probe && -e $append_probe ]]; then
    chattr -a "$append_probe" >/dev/null 2>&1
    rm -f -- "$append_probe"
  fi
  if [[ -n $overlay_base && -d $overlay_base ]]; then
    rm -rf -- "$overlay_base"
  fi
  if [[ -n $work_dir && -d $work_dir ]]; then
    rm -rf -- "$work_dir"
  fi
}

# Creates a private per-run directory (mode 0700, from mktemp) inside the
# state directory. Probes run there, so the permissions of a caller-supplied
# TIDEMARK_SELFTEST_DIR are never changed.
prepare_work_dir() {
  if [[ $state_dir != /* ]]; then
    log "TIDEMARK_SELFTEST_DIR must be an absolute path: $state_dir"
    return 1
  fi
  mkdir -p -- "$state_dir" 2>/dev/null || return 1
  work_dir=$(mktemp -d "$state_dir/run.XXXXXX" 2>/dev/null) || return 1
  [[ -d $work_dir && -w $work_dir ]]
}

# Records the seconds since boot at which the self-test started, so that the
# guest reports its own boot time (the boot test records it for target P2).
check_uptime() {
  local up
  if read -r up _ </proc/uptime 2>/dev/null && [[ -n $up ]]; then
    info uptime_s "$up"
  else
    info uptime_s unknown
  fi
}

check_kernel() {
  local release major rest minor
  release=$(uname -r)
  info kernel "$release"
  info arch "$(uname -m)"
  major=${release%%.*}
  rest=${release#*.}
  minor=${rest%%[!0-9]*}
  if [[ $major =~ ^[0-9]+$ && $minor =~ ^[0-9]+$ ]] &&
    ((major > MIN_KERNEL_MAJOR || (major == MIN_KERNEL_MAJOR && minor >= MIN_KERNEL_MINOR))); then
    require kernel_floor ok
  else
    require kernel_floor "below_${MIN_KERNEL_MAJOR}.${MIN_KERNEL_MINOR}"
  fi
}

check_lsm() {
  local lsm_file=/sys/kernel/security/lsm lsm name
  if [[ ! -r $lsm_file ]]; then
    info lsm unreadable
    require lsm_bpf securityfs_unavailable
    require lsm_landlock securityfs_unavailable
    return
  fi
  lsm=$(<"$lsm_file")
  info lsm "$lsm"
  for name in bpf landlock; do
    if [[ ,$lsm, == *",$name,"* ]]; then
      require "lsm_$name" ok
    else
      require "lsm_$name" not_active
    fi
  done
}

# Prints the Landlock ABI version, or error_<ERRNO> when the query fails.
landlock_abi_query() {
  python3 -I -S - "$SYS_LANDLOCK_CREATE_RULESET" <<'PY'
import ctypes
import errno
import sys

LANDLOCK_CREATE_RULESET_VERSION = 1 << 0
nr = int(sys.argv[1])
libc = ctypes.CDLL(None, use_errno=True)
libc.syscall.restype = ctypes.c_long
ret = libc.syscall(
    ctypes.c_long(nr),
    ctypes.c_void_p(None),
    ctypes.c_size_t(0),
    ctypes.c_uint32(LANDLOCK_CREATE_RULESET_VERSION),
)
if ret < 0:
    err = ctypes.get_errno()
    print("error_" + errno.errorcode.get(err, str(err)))
else:
    print(ret)
PY
}

check_landlock() {
  local abi
  if ! have python3; then
    info landlock_abi unknown
    require landlock missing_python3
    return
  fi
  abi=$(landlock_abi_query 2>/dev/null) || abi=error_python
  info landlock_abi "$abi"
  if [[ $abi =~ ^[0-9]+$ ]]; then
    if ((abi >= MIN_LANDLOCK_ABI)); then
      require landlock ok
    else
      require landlock "abi_below_${MIN_LANDLOCK_ABI}"
    fi
  else
    require landlock "$abi"
  fi
}

check_cgroup2() {
  local fstype
  fstype=$(stat -f -c %T /sys/fs/cgroup 2>/dev/null) || fstype=unknown
  if [[ $fstype == cgroup2fs && -r /sys/fs/cgroup/cgroup.controllers ]]; then
    require cgroup2 ok
    info cgroup2_controllers "$(tr ' ' ',' </sys/fs/cgroup/cgroup.controllers)"
    if [[ -e /sys/fs/cgroup/system.slice/cgroup.freeze ]]; then
      info cgroup2_freeze present
    else
      info cgroup2_freeze absent
    fi
  else
    require cgroup2 "fstype_${fstype}"
  fi
}

check_overlayfs() {
  if ! is_root; then
    require overlayfs requires_root
    return
  fi
  if ! have unshare || ! have mount; then
    require overlayfs missing_util_linux
    return
  fi
  overlay_base="$work_dir/overlay"
  rm -rf -- "$overlay_base"
  if ! mkdir -p -- "$overlay_base/lower" "$overlay_base/upper" "$overlay_base/work" "$overlay_base/merged"; then
    require overlayfs mkdir_failed
    return
  fi
  printf 'lower\n' >"$overlay_base/lower/from-lower"
  # The inner script runs in a private mount namespace, so the overlay mount
  # disappears when it exits. $1 is expanded by the inner shell on purpose.
  # shellcheck disable=SC2016
  if unshare --mount --propagation private -- bash -c '
      set -e
      mount -t overlay tidemark-selftest \
        -o "lowerdir=$1/lower,upperdir=$1/upper,workdir=$1/work" "$1/merged"
      grep -qx lower "$1/merged/from-lower"
      printf "upper\n" >"$1/merged/from-merged"
    ' tidemark-overlay "$overlay_base" >/dev/null 2>&1 &&
    [[ -f $overlay_base/upper/from-merged ]]; then
    require overlayfs ok
  else
    require overlayfs mount_or_write_failed
  fi
  rm -rf -- "$overlay_base"
  overlay_base=""
}

check_var_xattr() {
  local probe="$work_dir/xattr-probe" got
  if ! have setfattr || ! have getfattr; then
    require var_xattr missing_attr_tools
    info xattr_security_bpf not_tested
    return
  fi
  if ! : >"$probe" 2>/dev/null; then
    require var_xattr create_failed
    info xattr_security_bpf not_tested
    return
  fi
  if setfattr -n user.tidemark.selftest -v ok -- "$probe" 2>/dev/null &&
    got=$(getfattr --only-values -n user.tidemark.selftest -- "$probe" 2>/dev/null) &&
    [[ $got == ok ]]; then
    require var_xattr ok
  else
    require var_xattr user_xattr_unsupported
  fi
  # security.bpf.* is the namespace that bpf_get_file_xattr may read
  # (fs/bpf_fs_kfuncs.c). Whether root may set it from user space here is
  # recorded for the Phase 0 spike; it does not affect the result.
  if setfattr -n security.bpf.tidemark_selftest -v 1 -- "$probe" 2>/dev/null &&
    got=$(getfattr --only-values -n security.bpf.tidemark_selftest -- "$probe" 2>/dev/null) &&
    [[ $got == 1 ]]; then
    info xattr_security_bpf ok
  else
    info xattr_security_bpf set_or_get_failed
  fi
  rm -f -- "$probe"
}

check_var_append() {
  local flags reason
  if ! is_root; then
    require var_append requires_root
    return
  fi
  if ! have chattr || ! have lsattr; then
    require var_append missing_e2fsprogs
    return
  fi
  append_probe="$work_dir/append-probe"
  if ! printf 'first\n' >"$append_probe" 2>/dev/null; then
    require var_append create_failed
    append_probe=""
    return
  fi
  # chattr does not accept "--"; the path is absolute, so it cannot be
  # mistaken for an option.
  if ! chattr +a "$append_probe" 2>/dev/null; then
    rm -f -- "$append_probe"
    append_probe=""
    require var_append chattr_refused
    return
  fi
  flags=$(lsattr "$append_probe" 2>/dev/null)
  flags=${flags%% *}
  if [[ $flags != *a* ]]; then
    reason=flag_not_reported
  elif (: >"$append_probe") 2>/dev/null; then
    reason=truncate_allowed
  elif ! printf 'second\n' >>"$append_probe" 2>/dev/null; then
    reason=append_refused
  else
    reason=ok
  fi
  chattr -a "$append_probe" 2>/dev/null
  rm -f -- "$append_probe"
  append_probe=""
  require var_append "$reason"
}

check_var_fstype() {
  local fstype
  fstype=$(stat -f -c %T "$work_dir" 2>/dev/null) || fstype=unknown
  info var_fstype "$fstype"
}

check_btf_kfuncs() {
  local btf=/sys/kernel/btf/vmlinux name
  if [[ ! -r $btf ]]; then
    require btf absent
    for name in "${REQUIRED_KFUNCS[@]}"; do
      require "kfunc_$name" btf_absent
    done
    return
  fi
  require btf ok
  # BTF stores names as NUL-terminated strings, so match whole NUL-separated
  # records exactly.
  for name in "${REQUIRED_KFUNCS[@]}"; do
    if grep -a -z -x -F -q -- "$name" "$btf"; then
      require "kfunc_$name" ok
    else
      require "kfunc_$name" absent
    fi
  done
}

check_selinux() {
  local enforce=/sys/fs/selinux/enforce mode
  if [[ -r $enforce ]]; then
    mode=$(<"$enforce")
    if [[ $mode == 1 ]]; then
      info selinux enforcing
    else
      info selinux permissive
    fi
  else
    info selinux disabled
  fi
}

emit_result() {
  local result=PASS failed_list=none line key tmp
  if ((${#failed[@]} > 0)); then
    result=FAIL
    failed_list=$(
      IFS=,
      printf '%s' "${failed[*]}"
    )
  fi
  line="TIDEMARK_SELFTEST v=$SCHEMA_VERSION result=$result phase=0 failed=$failed_list"
  for key in "${keys[@]}"; do
    line+=" $key=${values[$key]}"
  done

  if mkdir -p -- "${result_file%/*}" 2>/dev/null; then
    tmp="$result_file.tmp.$$"
    if {
      printf 'v=%s\nresult=%s\nphase=0\nfailed=%s\n' "$SCHEMA_VERSION" "$result" "$failed_list"
      for key in "${keys[@]}"; do
        printf '%s=%s\n' "$key" "${values[$key]}"
      done
    } >"$tmp" 2>/dev/null; then
      mv -f -- "$tmp" "$result_file" 2>/dev/null || rm -f -- "$tmp"
    fi
  fi

  # The machine-readable line is always the last line of output.
  printf '%s\n' "$line"
  [[ $result == PASS ]]
}

main() {
  case "${1:-}" in
    "") ;;
    -h | --help)
      usage
      return 0
      ;;
    *)
      usage >&2
      return 2
      ;;
  esac

  trap cleanup EXIT
  log "AgentOS Tidemark Phase 0 boot self-test; it reports only and enforces nothing."

  check_uptime
  check_kernel
  check_lsm
  check_landlock
  check_cgroup2

  if prepare_work_dir; then
    check_var_fstype
    check_overlayfs
    check_var_xattr
    check_var_append
  else
    info var_fstype unknown
    require overlayfs work_dir_unavailable
    require var_xattr work_dir_unavailable
    require var_append work_dir_unavailable
  fi

  check_btf_kfuncs
  check_selinux

  if emit_result; then
    return 0
  fi
  return 1
}

main "$@"
