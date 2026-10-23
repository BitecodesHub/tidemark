// SPDX-License-Identifier: Apache-2.0
//
// Session-context assertions for in-session helpers. Phase 0 draft; built and run on aarch64 on 2026-10-18.
//
// tests/_lib.sh run_in_session / run_in_cgroup export TM_MOVED_PID (the pid of
// the subshell that was written into cgroup.procs) and TM_SESSION_CG (the
// cgroupfs directory it joined). When both are set, every helper checks, before
// doing anything else, that:
//   * it is a GRANDCHILD of the moved process (pid != moved, ppid != moved,
//     parent of ppid == moved), so scenario A1 really exercises a descendant
//     several levels below the process that joined the session, rather than
//     the joined process itself (bash can elide forks); and
//   * its own cgroup (from /proc/self/cgroup) is the expected cgroup.
// A mismatch is a setup error, reported as such, so it cannot be mistaken for
// a kernel hook failure.

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// parentOf returns the parent pid of pid, read from /proc/<pid>/stat.
func parentOf(pid int) (int, error) {
	b, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return 0, err
	}
	s := string(b)
	// The comm field is in parentheses and may contain spaces; parse after
	// the last ')'. The fields that follow are: state ppid ...
	i := strings.LastIndexByte(s, ')')
	if i < 0 {
		return 0, fmt.Errorf("malformed /proc/%d/stat", pid)
	}
	f := strings.Fields(s[i+1:])
	if len(f) < 2 {
		return 0, fmt.Errorf("short /proc/%d/stat", pid)
	}
	return strconv.Atoi(f[1])
}

// ownCgroup returns this process's cgroup v2 path relative to the cgroupfs
// root, for example "/tmspike.slice/tm-sess-1.scope".
func ownCgroup() (string, error) {
	b, err := os.ReadFile("/proc/self/cgroup")
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(strings.TrimSpace(string(b)), "\n") {
		if strings.HasPrefix(line, "0::") {
			return strings.TrimPrefix(line, "0::"), nil
		}
	}
	return "", fmt.Errorf("no cgroup v2 entry in /proc/self/cgroup")
}

// checkSessionContext enforces the lineage and cgroup expectations described
// at the top of this file. It is a no-op when TM_MOVED_PID is not set.
func checkSessionContext() error {
	movedStr := os.Getenv("TM_MOVED_PID")
	if movedStr == "" {
		return nil
	}
	moved, err := strconv.Atoi(movedStr)
	if err != nil {
		return fmt.Errorf("setup: bad TM_MOVED_PID %q: %w", movedStr, err)
	}
	pid := os.Getpid()
	ppid := os.Getppid()
	gpid, err := parentOf(ppid)
	if err != nil {
		return fmt.Errorf("setup: read parent of %d: %w", ppid, err)
	}
	fmt.Printf("LINEAGE pid=%d ppid=%d grandparent=%d moved=%d\n", pid, ppid, gpid, moved)
	if pid == moved || ppid == moved || gpid != moved {
		return fmt.Errorf("setup: helper is not a grandchild of the moved process " +
			"(a fork was elided); the grandchild claim would not be exercised")
	}

	if cg := os.Getenv("TM_SESSION_CG"); cg != "" {
		own, err := ownCgroup()
		if err != nil {
			return fmt.Errorf("setup: %w", err)
		}
		want := strings.TrimPrefix(cg, "/sys/fs/cgroup")
		fmt.Printf("CGROUP own=%s want=%s\n", own, want)
		if own != want {
			return fmt.Errorf("setup: helper runs in cgroup %s, expected %s", own, want)
		}
	}
	return nil
}
