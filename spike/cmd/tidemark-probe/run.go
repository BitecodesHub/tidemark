// SPDX-License-Identifier: Apache-2.0
//
// The `run` orchestrator and the `env` banner. Phase 0 draft; built and run on aarch64 on 2026-10-18.
//
// `run` loads and pins the probe, prints the environment banner the Phase 0
// acceptance criteria call for (kernel version, LSM list, Landlock ABI), then
// executes every scenario script under tests/ in sorted order, printing one
// PASS/FAIL line per assumption and a final SUMMARY line. It exits non-zero if
// any assumption fails, so CI can gate on the exit code and on "SUMMARY: ALL
// PASS".
//
// Robustness: each scenario runs in its own process group with a timeout
// (TM_SCENARIO_TIMEOUT, default 120 s), its output is streamed to the serial
// log as it is produced (so a hang or a kernel oops still leaves the
// scenario's last lines on the console), and SIGINT/SIGTERM unload the probe
// before exiting.

package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

// cmdEnv prints the environment facts the acceptance criteria require.
func cmdEnv(_ []string) error {
	ok := true

	if rel, err := os.ReadFile("/proc/sys/kernel/osrelease"); err == nil {
		fmt.Printf("ENV kernel: %s\n", strings.TrimSpace(string(rel)))
	} else {
		fmt.Printf("ENV kernel: unknown (%v)\n", err)
		ok = false
	}

	lsm, err := os.ReadFile("/sys/kernel/security/lsm")
	if err != nil {
		fmt.Printf("ENV lsm: unreadable (%v)\n", err)
		ok = false
	} else {
		list := strings.TrimSpace(string(lsm))
		fmt.Printf("ENV lsm: %s\n", list)
		for _, want := range []string{"bpf", "landlock"} {
			if !containsCSV(list, want) {
				fmt.Printf("ENV lsm: MISSING %q in active LSM list\n", want)
				ok = false
			}
		}
	}

	abi, _, errno := unix.Syscall(unix.SYS_LANDLOCK_CREATE_RULESET, 0, 0,
		uintptr(unix.LANDLOCK_CREATE_RULESET_VERSION))
	if errno != 0 {
		fmt.Printf("ENV landlock: unavailable (%s)\n", unix.ErrnoName(errno))
		ok = false
	} else {
		fmt.Printf("ENV landlock_abi: %d\n", int(abi))
		if int(abi) < 9 {
			fmt.Printf("ENV landlock_abi: BELOW required floor 9\n")
			ok = false
		}
	}

	fmt.Printf("ENV bpf_real_data_inode: %v\n", kernelHasRealDataInode())

	if !ok {
		return fmt.Errorf("environment does not meet the Phase 0 floor")
	}
	return nil
}

// prefixWriter prefixes every line written through it, for the serial log.
type prefixWriter struct {
	mu     sync.Mutex
	w      io.Writer
	prefix []byte
	bol    bool
}

func newPrefixWriter(w io.Writer, prefix string) *prefixWriter {
	return &prefixWriter{w: w, prefix: []byte(prefix), bol: true}
}

func (p *prefixWriter) Write(b []byte) (int, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, c := range b {
		if p.bol {
			if _, err := p.w.Write(p.prefix); err != nil {
				return 0, err
			}
			p.bol = false
		}
		if _, err := p.w.Write([]byte{c}); err != nil {
			return 0, err
		}
		if c == '\n' {
			p.bol = true
		}
	}
	return len(b), nil
}

func scenarioTimeout() time.Duration {
	if v := os.Getenv("TM_SCENARIO_TIMEOUT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return time.Duration(n) * time.Second
		}
	}
	return 120 * time.Second
}

// scriptTimeout reads a "# TIMEOUT: <seconds>" comment from a scenario script.
// The race scenario needs far more than the default budget; without this it
// would be reported as a kernel failure when it is only slow.
func scriptTimeout(path string) (int, bool) {
	f, err := os.Open(path)
	if err != nil {
		return 0, false
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for n := 0; sc.Scan() && n < 60; n++ {
		line := strings.TrimSpace(sc.Text())
		if !strings.HasPrefix(line, "# TIMEOUT:") {
			continue
		}
		secs, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "# TIMEOUT:")))
		if err != nil || secs <= 0 {
			return 0, false
		}
		return secs, true
	}
	return 0, false
}

// cmdRun is the top-level orchestrator.
func cmdRun(_ []string) error {
	self, err := os.Executable()
	if err != nil {
		return fmt.Errorf("find own path: %w", err)
	}
	if v := os.Getenv("TM_PROBE_BIN"); v != "" {
		self = v
	}

	fmt.Println("=== Tidemark Phase 0 kernel spike (Phase 0 measurement instrument; results are from this run) ===")

	// Environment banner. A failure here is reported and gates the summary.
	envErr := cmdEnv(nil)

	var unloadOnce sync.Once
	unload := func() {
		unloadOnce.Do(func() {
			if uerr := cmdUnload(nil); uerr != nil {
				fmt.Fprintf(os.Stderr, "WARN unload: %v\n", uerr)
			}
		})
	}
	defer unload()

	// The running scenario's process group, so a signal can be forwarded.
	var curMu sync.Mutex
	curPgid := 0

	// SIGINT/SIGTERM: stop the running scenario (it runs in its own process
	// group, so the terminal signal does not reach it), then detach the
	// programs before exiting.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-sigs
		fmt.Printf("INFO received %v; stopping the scenario and unloading the probe\n", s)
		curMu.Lock()
		if curPgid > 0 {
			_ = syscall.Kill(-curPgid, syscall.SIGTERM)
		}
		curMu.Unlock()
		time.Sleep(2 * time.Second)
		unload()
		os.Exit(130)
	}()

	// Load and attach once for the whole suite. cmdLoad removes its own pins
	// if it fails partway.
	if err := cmdLoad(nil); err != nil {
		return fmt.Errorf("load probe: %w", err)
	}

	testsDir := envOr("TM_TESTS_DIR", "./tests")
	scripts, err := filepath.Glob(filepath.Join(testsDir, "*.sh"))
	if err != nil {
		return fmt.Errorf("list scenarios: %w", err)
	}
	// Skip shared helpers (basename starting with "_"); only numbered scenario
	// scripts are run.
	filtered := scripts[:0]
	for _, s := range scripts {
		if strings.HasPrefix(filepath.Base(s), "_") {
			continue
		}
		filtered = append(filtered, s)
	}
	scripts = filtered
	sort.Strings(scripts)
	if len(scripts) == 0 {
		return fmt.Errorf("no scenario scripts found in %s", testsDir)
	}

	type result struct {
		name string
		desc string
		pass bool
	}
	var results []result

	for i, script := range scripts {
		name := filepath.Base(script)
		// A scenario may ask for a longer budget with a "# TIMEOUT: <seconds>"
		// comment, which the 10,000-iteration race test needs; the default is
		// TM_SCENARIO_TIMEOUT.
		perScript := scenarioTimeout()
		if secs, ok := scriptTimeout(script); ok {
			perScript = time.Duration(secs) * time.Second
		}
		fmt.Printf("--- %s (timeout %s)\n", name, perScript)

		ctx, cancel := context.WithTimeout(context.Background(), perScript)
		cmd := exec.CommandContext(ctx, "bash", script)
		cmd.Env = append(os.Environ(),
			"TM_PROBE_BIN="+self,
			"TM_BPFFS="+bpffsRoot(),
			fmt.Sprintf("TM_SLOT=%d", i+1), // slot 0 reserved; scenarios use 1..N
		)
		// Own process group, so a timeout stops the whole scenario tree:
		// SIGTERM first (the scenario's EXIT trap cleans up scopes, mounts and
		// scratch files), SIGKILL after a grace period.
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
		cmd.Cancel = func() error {
			pgid := cmd.Process.Pid
			err := syscall.Kill(-pgid, syscall.SIGTERM)
			go func() {
				time.Sleep(5 * time.Second)
				_ = syscall.Kill(-pgid, syscall.SIGKILL)
			}()
			return err
		}
		cmd.WaitDelay = 10 * time.Second

		var captured bytes.Buffer
		stream := io.MultiWriter(newPrefixWriter(os.Stdout, "    | "), &captured)
		cmd.Stdout = stream
		cmd.Stderr = stream
		var runErr error
		if runErr = cmd.Start(); runErr == nil {
			curMu.Lock()
			curPgid = cmd.Process.Pid
			curMu.Unlock()
			runErr = cmd.Wait()
			curMu.Lock()
			curPgid = 0
			curMu.Unlock()
		}
		timedOut := ctx.Err() == context.DeadlineExceeded
		cancel()
		if timedOut {
			cleanupLeftovers(i + 1)
		}

		desc := name
		for _, line := range strings.Split(captured.String(), "\n") {
			if strings.HasPrefix(line, "DESC:") {
				desc = strings.TrimSpace(strings.TrimPrefix(line, "DESC:"))
			}
		}

		pass := runErr == nil
		results = append(results, result{name: name, desc: desc, pass: pass})
		switch {
		case pass:
			fmt.Printf("PASS %s - %s\n", name, desc)
		case timedOut:
			fmt.Printf("FAIL %s - %s (timed out after %s)\n", name, desc, perScript)
		default:
			fmt.Printf("FAIL %s - %s (%v)\n", name, desc, runErr)
		}
	}

	passed := 0
	for _, r := range results {
		if r.pass {
			passed++
		}
	}
	total := len(results)

	// The environment floor is itself a gate.
	if envErr != nil {
		fmt.Printf("FAIL environment - %v\n", envErr)
	}

	if passed == total && envErr == nil {
		fmt.Printf("SUMMARY: ALL PASS (%d/%d)\n", passed, total)
		return nil
	}
	fmt.Printf("SUMMARY: FAIL (%d/%d scenarios passed; environment ok=%v)\n",
		passed, total, envErr == nil)
	return fmt.Errorf("one or more Phase 0 assumptions failed")
}

// cleanupLeftovers removes what a killed scenario may have left behind: its
// transient scopes and any mounts under the spike's scratch directories.
func cleanupLeftovers(slot int) {
	_ = exec.Command("systemctl", "stop",
		fmt.Sprintf("tm-sess-%d.scope", slot), fmt.Sprintf("tm-other-%d.scope", slot)).Run()
	if b, err := os.ReadFile("/proc/self/mountinfo"); err == nil {
		for _, line := range strings.Split(string(b), "\n") {
			f := strings.Fields(line)
			if len(f) > 4 && strings.HasPrefix(f[4], "/var/tmp/tm-spike.") {
				_ = exec.Command("umount", "-l", f[4]).Run()
			}
		}
	}
	fmt.Printf("INFO cleaned up leftovers of slot %d after a timeout\n", slot)
}

func containsCSV(list, want string) bool {
	for _, f := range strings.Split(list, ",") {
		if strings.TrimSpace(f) == want {
			return true
		}
	}
	return false
}
