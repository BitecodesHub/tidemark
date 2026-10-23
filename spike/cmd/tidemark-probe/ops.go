// SPDX-License-Identifier: Apache-2.0
//
// File operations performed by the scenario scripts inside or outside a
// session scope: open, mmap read, descriptor passing over SCM_RIGHTS, the
// read-only label mapping check, and the "join the session holding an
// inherited descriptor" case. Phase 0 draft; built and run on aarch64 on 2026-10-18.
//
// Label sampling: helpers that run as root sample the session label word
// through a read-only mmap of the pinned array (readLabelRO) IMMEDIATELY after
// the operation under test returns, before any further I/O, so a PASS shows
// the label was raised synchronously by that operation. Helpers that run as an
// unprivileged uid cannot open the root-only bpffs pins; they take --no-label,
// perform the operation, and the scenario script then reads the label as root.

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/cilium/ebpf"
	"golang.org/x/sys/unix"
)

// opFlags are the common flags of the label-sampling helpers.
type opFlags struct {
	slot    *uint
	expect  *string
	nonzero *bool
	noLabel *bool
}

func addOpFlags(fs *flag.FlagSet) opFlags {
	return opFlags{
		slot:    fs.Uint("slot", 0, "label slot index"),
		expect:  fs.String("expect", "", "expected label (hex), optional"),
		nonzero: fs.Bool("expect-nonzero", false, "require a raised label"),
		noLabel: fs.Bool("no-label", false, "do not sample the label (unprivileged callers)"),
	}
}

// sampleAndCheck reads the label right away and enforces the expectations,
// unless --no-label was given.
func (o opFlags) sampleAndCheck(what string) error {
	if *o.noLabel {
		fmt.Printf("%s: label not sampled (--no-label); the script checks it as root\n", what)
		return nil
	}
	got, err := readLabelRO(uint32(*o.slot))
	if err != nil {
		return err
	}
	return checkExpect(what, uint32(*o.slot), got, *o.expect, *o.nonzero)
}

func errnoName(err error) string {
	var errno syscall.Errno
	if errors.As(err, &errno) {
		if n := unix.ErrnoName(errno); n != "" {
			return n
		}
	}
	return err.Error()
}

// cmdOpen opens a file read-only, samples the label immediately after open()
// returns (before read()), then reads a few bytes.
func cmdOpen(args []string) error {
	fs := flag.NewFlagSet("open", flag.ExitOnError)
	path := fs.String("path", "", "file to open")
	o := addOpFlags(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *path == "" {
		return errors.New("--path is required")
	}
	if err := checkSessionContext(); err != nil {
		return err
	}

	fd, err := unix.Open(*path, unix.O_RDONLY|unix.O_CLOEXEC, 0)
	if err != nil {
		return fmt.Errorf("setup: open %s: %s", *path, errnoName(err))
	}
	defer unix.Close(fd)

	// Sample first: the claim under test is "raised before open() returns".
	checkErr := o.sampleAndCheck("open")

	buf := make([]byte, 64)
	n, rerr := unix.Read(fd, buf)
	if rerr != nil {
		fmt.Printf("open %s: read failed: %s\n", *path, errnoName(rerr))
	} else {
		fmt.Printf("open %s: read %d bytes after sampling\n", *path, n)
	}
	return checkErr
}

// fileSize returns the size of an open file, at least one byte for mmap.
func fileSize(fd int) (int, error) {
	var st unix.Stat_t
	if err := unix.Fstat(fd, &st); err != nil {
		return 0, err
	}
	if st.Size <= 0 {
		return 0, errors.New("file is empty; mmap needs at least one byte")
	}
	return int(st.Size), nil
}

// cmdMmapRead opens a file, maps it PROT_READ, touches the first byte, then
// samples the label. The open alone already raises the label through
// file_open; scenario 10 isolates the mmap_file hook.
func cmdMmapRead(args []string) error {
	fs := flag.NewFlagSet("mmapread", flag.ExitOnError)
	path := fs.String("path", "", "file to map")
	o := addOpFlags(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *path == "" {
		return errors.New("--path is required")
	}
	if err := checkSessionContext(); err != nil {
		return err
	}

	fd, err := unix.Open(*path, unix.O_RDONLY|unix.O_CLOEXEC, 0)
	if err != nil {
		return fmt.Errorf("setup: open %s: %s", *path, errnoName(err))
	}
	defer unix.Close(fd)
	size, err := fileSize(fd)
	if err != nil {
		return fmt.Errorf("setup: %w", err)
	}
	b, err := unix.Mmap(fd, 0, size, unix.PROT_READ, unix.MAP_PRIVATE)
	if err != nil {
		return fmt.Errorf("setup: mmap %s: %s", *path, errnoName(err))
	}
	defer unix.Munmap(b)
	fmt.Printf("mmapread %s: first byte 0x%02x\n", *path, b[0])
	return o.sampleAndCheck("mmapread")
}

// cmdFdSend opens a file (normally outside any session) and passes the
// descriptor over a connected unix stream socket with SCM_RIGHTS.
func cmdFdSend(args []string) error {
	fs := flag.NewFlagSet("fd-send", flag.ExitOnError)
	path := fs.String("path", "", "file to open and send")
	sock := fs.String("sock", "", "unix socket path of the receiver")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *path == "" || *sock == "" {
		return errors.New("--path and --sock are required")
	}

	fd, err := unix.Open(*path, unix.O_RDONLY|unix.O_CLOEXEC, 0)
	if err != nil {
		return fmt.Errorf("setup: open %s: %s", *path, errnoName(err))
	}
	defer unix.Close(fd)

	s, err := unix.Socket(unix.AF_UNIX, unix.SOCK_STREAM|unix.SOCK_CLOEXEC, 0)
	if err != nil {
		return fmt.Errorf("setup: socket: %w", err)
	}
	defer unix.Close(s)
	// The socket file appears at bind(), before the receiver calls listen(),
	// so retry a refused or missing socket for up to about two seconds.
	var cerr error
	for i := 0; i < 40; i++ {
		cerr = unix.Connect(s, &unix.SockaddrUnix{Name: *sock})
		if cerr == nil || (!errors.Is(cerr, unix.ECONNREFUSED) && !errors.Is(cerr, unix.ENOENT)) {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if cerr != nil {
		return fmt.Errorf("setup: connect %s: %s", *sock, errnoName(cerr))
	}
	if err := unix.Sendmsg(s, []byte{'x'}, unix.UnixRights(fd), nil, 0); err != nil {
		return fmt.Errorf("setup: sendmsg: %s", errnoName(err))
	}
	fmt.Printf("fd-send %s: descriptor sent over %s\n", *path, *sock)
	return nil
}

// cmdFdRecv binds and listens on a unix socket, accepts one connection,
// receives one descriptor with SCM_RIGHTS (the file_receive hook fires inside
// recvmsg), samples the label, then reads from the received descriptor.
func cmdFdRecv(args []string) error {
	fs := flag.NewFlagSet("fd-recv", flag.ExitOnError)
	sock := fs.String("sock", "", "unix socket path to bind")
	o := addOpFlags(fs)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *sock == "" {
		return errors.New("--sock is required")
	}
	if err := checkSessionContext(); err != nil {
		return err
	}

	l, err := unix.Socket(unix.AF_UNIX, unix.SOCK_STREAM|unix.SOCK_CLOEXEC, 0)
	if err != nil {
		return fmt.Errorf("setup: socket: %w", err)
	}
	defer unix.Close(l)
	_ = os.Remove(*sock)
	if err := unix.Bind(l, &unix.SockaddrUnix{Name: *sock}); err != nil {
		return fmt.Errorf("setup: bind %s: %s", *sock, errnoName(err))
	}
	defer os.Remove(*sock)
	if err := unix.Listen(l, 1); err != nil {
		return fmt.Errorf("setup: listen: %w", err)
	}
	c, _, err := unix.Accept4(l, unix.SOCK_CLOEXEC)
	if err != nil {
		return fmt.Errorf("setup: accept: %w", err)
	}
	defer unix.Close(c)

	buf := make([]byte, 1)
	oob := make([]byte, unix.CmsgSpace(4))
	_, oobn, _, _, err := unix.Recvmsg(c, buf, oob, unix.MSG_CMSG_CLOEXEC)
	if err != nil {
		return fmt.Errorf("setup: recvmsg: %s", errnoName(err))
	}
	msgs, err := unix.ParseSocketControlMessage(oob[:oobn])
	if err != nil || len(msgs) == 0 {
		return fmt.Errorf("setup: no control message received (%v)", err)
	}
	fds, err := unix.ParseUnixRights(&msgs[0])
	if err != nil || len(fds) == 0 {
		return fmt.Errorf("setup: no descriptor in SCM_RIGHTS (%v)", err)
	}
	rfd := fds[0]
	defer unix.Close(rfd)

	checkErr := o.sampleAndCheck("fd-recv")

	data := make([]byte, 16)
	if n, rerr := unix.Read(rfd, data); rerr != nil {
		fmt.Printf("fd-recv: read on received descriptor failed: %s\n", errnoName(rerr))
	} else {
		fmt.Printf("fd-recv: read %d bytes from the received descriptor\n", n)
	}
	return checkErr
}

// cmdCheckROMmap checks the label-transport properties the exits rely on:
//  1. a PROT_READ, MAP_SHARED mapping of the labels array reads the value;
//  2. that mapping cannot be made writable with mprotect (bpf_map_mmap clears
//     VM_MAYWRITE on a read-only mapping), checked in the kernel, not in Go;
//  3. a BPF_F_RDONLY descriptor refuses a writable shared mapping;
//  4. OBSERVE: whether a BPF_F_RDONLY descriptor can be mapped read-only at all
//     (expected to fail with EINVAL at v7.2; see mapLabelsRO). This is a
//     design input for the Phase 1 label transport, not a gate.
func cmdCheckROMmap(args []string) error {
	fs := flag.NewFlagSet("check-ro-mmap", flag.ExitOnError)
	slot := fs.Uint("slot", 0, "label slot index")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *slot >= maxSlots {
		return fmt.Errorf("--slot %d out of range", *slot)
	}
	page := os.Getpagesize()
	failed := false

	b, err := mapLabelsRO()
	if err != nil {
		return err
	}
	defer unix.Munmap(b)
	off := int(*slot) * 8
	fmt.Printf("check-ro-mmap: PROT_READ mapping reads label=0x%x\n", nativeEndian.Uint64(b[off:off+8]))

	if err := unix.Mprotect(b, unix.PROT_READ|unix.PROT_WRITE); err == nil {
		fmt.Println("check-ro-mmap: FAIL mprotect(PROT_WRITE) succeeded; the mapping could be made writable")
		failed = true
	} else {
		fmt.Printf("check-ro-mmap: mprotect(PROT_WRITE) refused: %s\n", errnoName(err))
	}

	ro, err := ebpf.LoadPinnedMap(pinPath("maps", "labels"), &ebpf.LoadPinOptions{ReadOnly: true})
	if err != nil {
		return fmt.Errorf("open labels map read-only: %w", err)
	}
	defer ro.Close()
	if w, err := unix.Mmap(ro.FD(), 0, page, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED); err == nil {
		_ = unix.Munmap(w)
		fmt.Println("check-ro-mmap: FAIL a BPF_F_RDONLY descriptor accepted a writable shared mapping")
		failed = true
	} else {
		fmt.Printf("check-ro-mmap: writable mapping of a BPF_F_RDONLY descriptor refused: %s\n", errnoName(err))
	}
	if r, err := unix.Mmap(ro.FD(), 0, page, unix.PROT_READ, unix.MAP_SHARED); err == nil {
		_ = unix.Munmap(r)
		fmt.Println("OBSERVE a BPF_F_RDONLY descriptor CAN be mapped PROT_READ|MAP_SHARED on this kernel")
	} else {
		fmt.Printf("OBSERVE a BPF_F_RDONLY descriptor cannot be mapped PROT_READ|MAP_SHARED: %s (label transport must hand exits a read-write descriptor or use another channel)\n", errnoName(err))
	}

	if failed {
		return errors.New("label mapping is not read-only")
	}
	return nil
}

// cmdJoinThenUse opens a file OUTSIDE the session, then moves itself into the
// session cgroup while still holding the descriptor, then (i) calls read() and
// reports the label (expected 0x0: read() on a pre-session descriptor fires no
// hook this probe attaches, a documented gap that the production shim closes
// with close_range before exec), and (ii) maps the same descriptor and checks
// that mmap_file raised the label. Root only.
func cmdJoinThenUse(args []string) error {
	fs := flag.NewFlagSet("join-then-use", flag.ExitOnError)
	path := fs.String("path", "", "file to open before joining")
	cg := fs.String("cgroup", "", "session cgroupfs directory to join")
	slot := fs.Uint("slot", 0, "label slot index")
	expectMmap := fs.String("expect-after-mmap", "0x1", "expected label after mmap (hex)")
	skipMmap := fs.Bool("skip-mmap-gate", false, "report the mmap result but do not gate on it")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *path == "" || *cg == "" {
		return errors.New("--path and --cgroup are required")
	}

	fd, err := unix.Open(*path, unix.O_RDONLY|unix.O_CLOEXEC, 0)
	if err != nil {
		return fmt.Errorf("setup: open %s: %s", *path, errnoName(err))
	}
	defer unix.Close(fd)

	// Join the session. Writing the pid to cgroup.procs moves every thread of
	// this (multi-threaded Go) process.
	if err := os.WriteFile(*cg+"/cgroup.procs", []byte(strconv.Itoa(os.Getpid())), 0); err != nil {
		return fmt.Errorf("setup: join %s: %w", *cg, err)
	}
	if own, err := ownCgroup(); err == nil {
		fmt.Printf("join-then-use: now in cgroup %s\n", own)
	}

	buf := make([]byte, 16)
	if _, err := unix.Read(fd, buf); err != nil {
		return fmt.Errorf("setup: read on inherited descriptor: %s", errnoName(err))
	}
	l1, err := readLabelRO(uint32(*slot))
	if err != nil {
		return err
	}
	fmt.Printf("OBSERVE inherited-descriptor read(): label=0x%x (0x0 is the documented gap)\n", l1)

	size, err := fileSize(fd)
	if err != nil {
		return fmt.Errorf("setup: %w", err)
	}
	b, err := unix.Mmap(fd, 0, size, unix.PROT_READ, unix.MAP_PRIVATE)
	if err != nil {
		return fmt.Errorf("setup: mmap inherited descriptor: %s", errnoName(err))
	}
	defer unix.Munmap(b)
	_ = b[0]
	l2, err := readLabelRO(uint32(*slot))
	if err != nil {
		return err
	}
	if *skipMmap {
		fmt.Printf("OBSERVE inherited-descriptor mmap: label=0x%x (gate skipped)\n", l2)
		return nil
	}
	return checkExpect("join-then-use mmap", uint32(*slot), l2, *expectMmap, false)
}
