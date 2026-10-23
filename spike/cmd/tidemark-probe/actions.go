// SPDX-License-Identifier: Apache-2.0
//
// Setup and inspection helpers invoked by the scenario scripts (register a
// session, install a directory rule, read the label or cause word, set an
// xattr). The file operations live in ops.go. Phase 0 draft; built and run on aarch64 on 2026-10-18.
// Helpers that take --expect, --expect-nonzero, --has, --not or --require
// return a non-zero exit on mismatch, so a script can rely on the exit code.

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/cilium/ebpf"
	"golang.org/x/sys/unix"
)

// parseHex parses an optional hex/decimal value. ok is false when s is empty.
func parseHex(flagName, s string) (val uint64, ok bool, err error) {
	if s == "" {
		return 0, false, nil
	}
	v, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		return 0, false, fmt.Errorf("parse %s %q: %w", flagName, s, err)
	}
	return v, true, nil
}

// checkExpect prints the label (and, when readable, the slot's cause word) and
// enforces --expect (exact match) and --expect-nonzero.
func checkExpect(what string, slot uint32, got uint64, expectStr string, nonzero bool) error {
	want, ok, err := parseHex("--expect", expectStr)
	if err != nil {
		return err
	}
	line := fmt.Sprintf("%s label=0x%x", what, got)
	if c, cerr := readCause(slot); cerr == nil {
		line += fmt.Sprintf(" cause=0x%x(%s)", c, causeString(c))
	}
	fmt.Println(line)
	if ok && got != want {
		return fmt.Errorf("%s: label 0x%x, expected 0x%x", what, got, want)
	}
	if nonzero && got == labelPublic {
		return fmt.Errorf("%s: label 0x0, expected a raised label", what)
	}
	return nil
}

// cmdRegisterSession maps a session cgroup to a label slot and clears the slot.
func cmdRegisterSession(args []string) error {
	fs := flag.NewFlagSet("register-session", flag.ExitOnError)
	cgroup := fs.String("cgroup", "", "cgroupfs directory of the session scope")
	slot := fs.Uint("slot", 0, "label slot index")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *cgroup == "" {
		return errors.New("--cgroup is required")
	}
	if *slot >= maxSlots {
		return fmt.Errorf("--slot %d out of range", *slot)
	}

	// The cgroup depth below /sys/fs/cgroup must equal session_level, or
	// bpf_get_current_ancestor_cgroup_id(session_level) never returns this
	// cgroup's id. This catches the systemd dashed-slice nesting trap (a slice
	// named a-b.slice nests under a.slice). A mismatch is a hard error.
	rel, err := filepath.Rel("/sys/fs/cgroup", filepath.Clean(*cgroup))
	if err != nil || strings.HasPrefix(rel, "..") {
		return fmt.Errorf("--cgroup %s is not under /sys/fs/cgroup", *cgroup)
	}
	depth := len(strings.Split(rel, string(os.PathSeparator)))
	if uint32(depth) != sessionLevel {
		return fmt.Errorf("cgroup %s has depth %d but session_level is %d; keying would miss (check the slice name)",
			*cgroup, depth, sessionLevel)
	}

	// On cgroup v2 the cgroup id equals the inode number of the cgroup
	// directory (64-bit kernels), which is what
	// bpf_get_current_ancestor_cgroup_id returns.
	var st unix.Stat_t
	if err := unix.Stat(*cgroup, &st); err != nil {
		return fmt.Errorf("stat cgroup %s: %w", *cgroup, err)
	}
	cid := uint64(st.Ino)

	sessions, err := openPinnedMap("sessions")
	if err != nil {
		return err
	}
	defer sessions.Close()
	if err := sessions.Update(cid, uint32(*slot), ebpf.UpdateAny); err != nil {
		return fmt.Errorf("register session: %w", err)
	}

	// Clear the label and cause words for a clean run.
	for _, name := range []string{"labels", "causes"} {
		m, err := openPinnedMap(name)
		if err != nil {
			return err
		}
		uerr := m.Update(uint32(*slot), uint64(0), ebpf.UpdateAny)
		m.Close()
		if uerr != nil {
			return fmt.Errorf("clear %s slot: %w", name, uerr)
		}
	}

	fmt.Printf("INFO registered cgroup %s id=%d -> slot %d\n", *cgroup, cid, *slot)
	return nil
}

// cmdSetRoot installs a labelled directory rule in slot --index. A trailing
// slash is removed so that the BPF side can require a path-component boundary
// right after the prefix. --clear empties the slot instead.
func cmdSetRoot(args []string) error {
	fs := flag.NewFlagSet("set-root", flag.ExitOnError)
	index := fs.Uint("index", 0, "root slot index")
	prefix := fs.String("prefix", "", "absolute path prefix")
	labelName := fs.String("label", "private", "public|private|secret")
	clearSlot := fs.Bool("clear", false, "empty this root slot")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *index >= maxRoots {
		return fmt.Errorf("--index %d out of range", *index)
	}

	var r tmRoot
	if !*clearSlot {
		p := strings.TrimRight(*prefix, "/")
		if !strings.HasPrefix(p, "/") {
			return errors.New("--prefix must be an absolute path other than /")
		}
		if len(p) >= rootLen {
			return fmt.Errorf("--prefix too long (max %d bytes)", rootLen-1)
		}
		lbl, err := labelFromName(*labelName)
		if err != nil {
			return err
		}
		copy(r.Prefix[:], p)
		r.Label = lbl
		r.Len = uint32(len(p))
	}

	roots, err := openPinnedMap("roots")
	if err != nil {
		return err
	}
	defer roots.Close()
	if err := roots.Update(uint32(*index), &r, ebpf.UpdateAny); err != nil {
		return fmt.Errorf("set root: %w", err)
	}
	if *clearSlot {
		fmt.Printf("INFO root[%d] cleared\n", *index)
	} else {
		fmt.Printf("INFO root[%d] = %q label=%s\n", *index, string(r.Prefix[:r.Len]), *labelName)
	}
	return nil
}

// cmdReadLabel prints the session label word for a slot and optionally checks
// it. Root only (it opens the pinned map).
func cmdReadLabel(args []string) error {
	fs := flag.NewFlagSet("readlabel", flag.ExitOnError)
	slot := fs.Uint("slot", 0, "label slot index")
	expect := fs.String("expect", "", "expected label (hex), optional")
	nonzero := fs.Bool("expect-nonzero", false, "require a raised label")
	quiet := fs.Bool("raw", false, "print only label=0x..")
	if err := fs.Parse(args); err != nil {
		return err
	}
	v, err := readLabelRO(uint32(*slot))
	if err != nil {
		return err
	}
	if *quiet {
		fmt.Printf("label=0x%x\n", v)
		return nil
	}
	return checkExpect("readlabel", uint32(*slot), v, *expect, *nonzero)
}

// parseCauseList turns "open,xattr" or "0x11" into a cause bitmask.
func parseCauseList(flagName, s string) (uint64, error) {
	if s == "" {
		return 0, nil
	}
	if v, ok, err := parseHex(flagName, s); err == nil && ok {
		return v, nil
	}
	byName := map[string]uint64{
		"open": causeHookOpen, "receive": causeHookReceive, "mmap": causeHookMmap,
		"xattr": causeXattr, "dir": causeDir, "cache": causeCache,
		"error": causeError, "backing": causeBacking,
		"eacces": causeEacces, "nobuf": causeNobuf,
	}
	var out uint64
	for _, n := range strings.Split(s, ",") {
		n = strings.TrimSpace(n)
		b, ok := byName[n]
		if !ok {
			return 0, fmt.Errorf("%s: unknown cause %q", flagName, n)
		}
		out |= b
	}
	return out, nil
}

// cmdReadCause prints the cause word for a slot and checks that every bit in
// --has is set and no bit in --not is set. Root only.
func cmdReadCause(args []string) error {
	fs := flag.NewFlagSet("readcause", flag.ExitOnError)
	slot := fs.Uint("slot", 0, "label slot index")
	hasStr := fs.String("has", "", "cause bits that must be set (names or hex)")
	notStr := fs.String("not", "", "cause bits that must NOT be set (names or hex)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	has, err := parseCauseList("--has", *hasStr)
	if err != nil {
		return err
	}
	not, err := parseCauseList("--not", *notStr)
	if err != nil {
		return err
	}
	c, err := readCause(uint32(*slot))
	if err != nil {
		return err
	}
	fmt.Printf("readcause cause=0x%x(%s)\n", c, causeString(c))
	if c&has != has {
		return fmt.Errorf("cause 0x%x is missing required bits %s", c, causeString(has&^c))
	}
	if c&not != 0 {
		return fmt.Errorf("cause 0x%x has forbidden bits %s", c, causeString(c&not))
	}
	return nil
}

// cmdSetxattr tries to set an xattr and prints "setxattr <name>: ok" or
// "setxattr <name>: errno=<NAME>". With --require ok|<ERRNO> it exits non-zero
// unless the outcome matches, so scripts cannot mistake a setup failure for a
// kernel result. It never checks capabilities itself; the kernel decides.
func cmdSetxattr(args []string) error {
	fs := flag.NewFlagSet("setxattr", flag.ExitOnError)
	path := fs.String("path", "", "file")
	name := fs.String("name", "", "xattr name")
	value := fs.String("value", "", "xattr value")
	require := fs.String("require", "", "required outcome: ok or an errno name such as EPERM")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *path == "" || *name == "" {
		return errors.New("--path and --name are required")
	}
	if err := checkSessionContext(); err != nil {
		return err
	}

	outcome := "ok"
	if err := unix.Setxattr(*path, *name, []byte(*value), 0); err != nil {
		var errno syscall.Errno
		if errors.As(err, &errno) {
			outcome = unix.ErrnoName(errno)
		} else {
			outcome = err.Error()
		}
		fmt.Printf("setxattr %s: errno=%s\n", *name, outcome)
	} else {
		fmt.Printf("setxattr %s: ok\n", *name)
	}

	if *require != "" && *require != outcome {
		return fmt.Errorf("setxattr %s on %s: outcome %s, required %s", *name, *path, outcome, *require)
	}
	return nil
}
