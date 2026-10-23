// SPDX-License-Identifier: Apache-2.0
//
// tidemark-probe: AgentOS Tidemark Phase 0 kernel-spike driver.
//
// STATUS: Phase 0 measurement instrument. First built (go build, go vet) and
// run on 2026-10-18 inside a Fedora 45 Beta aarch64 VM (kernel 7.2.0-61.fc45):
// all 12 scenarios passed in three identical runs. Not yet run on x86_64. It
// uses the cilium/ebpf v0.22.0 and golang.org/x/sys v0.43.0 APIs (see
// spike/README.md). Build and run happen only inside a Fedora 45 VM.
//
// The driver has two roles:
//   - As an orchestrator (`tidemark-probe run`): it loads and attaches the BPF
//     object once, pins the maps and links, configures the labelled roots, then
//     executes each scenario script under tests/, printing one PASS/FAIL line
//     per assumption and a final summary line. It exits non-zero if any
//     assumption fails.
//   - As a set of in-scenario helpers (open, mmapread, fd-send, fd-recv,
//     setxattr, readlabel, register-session, ...): the scenario scripts invoke
//     these inside and outside systemd-run scopes to perform the actual
//     operations and to read the per-session label word back.
//
// Every operation in this program only observes and reports. It never denies a
// syscall; the BPF programs return 0 unconditionally.
package main

import (
	"fmt"
	"os"
)

const usage = `tidemark-probe - Tidemark Phase 0 kernel spike driver (Phase 0 measurement instrument)

Orchestration:
  run                      Load+pin the probe, run every tests/*.sh, summarise.
  env                      Print kernel version, LSM list and Landlock ABI.

Lifecycle helpers (used by scripts and by run):
  load                     Load, attach and pin the BPF object.
  unload                   Remove all pins (detach).
  register-session         Map a session cgroup to a label slot.
                           --cgroup <cgroupfs-dir> --slot <n>
  set-root                 Install (or --clear) a labelled directory rule.
                           --index <i> --prefix <abs-path> --label <public|private|secret>

Action helpers (run inside or outside a session scope):
  open                     open(), sample the label, then read().
                           --path <f> --slot <n> [--expect <hex>] [--no-label]
  mmapread                 open()+mmap(PROT_READ), touch, then sample the label.
                           --path <f> --slot <n> [--expect <hex>] [--no-label]
  fd-send                  Open a file and pass its fd over a unix socket.
                           --path <f> --sock <unix-path>
  fd-recv                  Receive an fd (SCM_RIGHTS), sample the label, read.
                           --sock <unix-path> --slot <n> [--expect <hex>] [--no-label]
  join-then-use            Open a file, join a session cgroup, read(), mmap().
                           --path <f> --cgroup <dir> --slot <n> [--expect-after-mmap <hex>]
  setxattr                 Try to set an xattr; print ok or the errno.
                           --path <f> --name <xattr> --value <v> [--require ok|EPERM|...]
  readlabel                Print the session label word for a slot (hex).
                           --slot <n> [--expect <hex>] [--expect-nonzero] [--raw]
  readcause                Print the cause bitmask for a slot and check it.
                           --slot <n> [--has open,xattr,...] [--not dir,...]
  check-ro-mmap            Check the label-array mapping cannot be made writable.
                           --slot <n>
  clear-slot               Zero a slot's label and cause words.
                           --slot <n>
  removexattr              Remove an xattr as root and verify it is gone.
                           --path <f> --name <xattr>
  try-move                 Write own pid into a cgroup's cgroup.procs; report errno.
                           --cgroup <dir> [--require ok|EACCES|...]
  try-mkdir                Attempt to create a cgroup subdirectory; report errno.
                           --dir <path> [--require ok|EACCES|...]
  bench-open               Measure open()+close() latency on one path.
                           --path <f> [--iterations <n>] [--warmup <n>]
                           [--slot <n>] [--check-raised] [--clear-each]
  race-open                In-session self race: clear, open, sample own label.
                           --path <f> --slot <n> [--iterations 10000]
                           [--expect 0x1] [--notify <sock>]
  race-observe             Observer race outside the session: wait for the byte,
                           then sample the label.
                           --sock <path> --slot <n> [--iterations 10000]
                           [--expect 0x1]

Environment:
  TM_BPF_OBJ    path to label_probe.bpf.o   (default ./label_probe.bpf.o)
  TM_BPFFS      bpffs pin root              (default /sys/fs/bpf/tidemark)
  TM_TESTS_DIR  scenario directory          (default ./tests)
  TM_PROBE_BIN  this binary, for scripts    (default the running executable)
  TM_ATTACH_MMAP_FILE  0 leaves lsm.s/mmap_file detached (default 1)
  TM_SCENARIO_TIMEOUT  per-scenario timeout in seconds (default 120)
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}

	var err error
	switch os.Args[1] {
	case "run":
		err = cmdRun(os.Args[2:])
	case "env":
		err = cmdEnv(os.Args[2:])
	case "load":
		err = cmdLoad(os.Args[2:])
	case "unload":
		err = cmdUnload(os.Args[2:])
	case "register-session":
		err = cmdRegisterSession(os.Args[2:])
	case "set-root":
		err = cmdSetRoot(os.Args[2:])
	case "open":
		err = cmdOpen(os.Args[2:])
	case "mmapread":
		err = cmdMmapRead(os.Args[2:])
	case "fd-send":
		err = cmdFdSend(os.Args[2:])
	case "fd-recv":
		err = cmdFdRecv(os.Args[2:])
	case "setxattr":
		err = cmdSetxattr(os.Args[2:])
	case "readlabel":
		err = cmdReadLabel(os.Args[2:])
	case "readcause":
		err = cmdReadCause(os.Args[2:])
	case "join-then-use":
		err = cmdJoinThenUse(os.Args[2:])
	case "check-ro-mmap":
		err = cmdCheckROMmap(os.Args[2:])
	case "clear-slot":
		err = cmdClearSlot(os.Args[2:])
	case "removexattr":
		err = cmdRemoveXattr(os.Args[2:])
	case "try-move":
		err = cmdTryMove(os.Args[2:])
	case "try-mkdir":
		err = cmdMkdirTry(os.Args[2:])
	case "bench-open":
		err = cmdBenchOpen(os.Args[2:])
	case "race-open":
		err = cmdRaceOpen(os.Args[2:])
	case "race-observe":
		err = cmdRaceObserve(os.Args[2:])
	case "-h", "--help", "help":
		fmt.Fprint(os.Stdout, usage)
		return
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand %q\n\n%s", os.Args[1], usage)
		os.Exit(2)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "tidemark-probe %s: %v\n", os.Args[1], err)
		os.Exit(1)
	}
}
