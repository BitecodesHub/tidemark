// SPDX-License-Identifier: Apache-2.0
//
// Load, attach and pin the BPF object, plus the shared map/label helpers used
// by the action subcommands. Phase 0 draft; built and run on aarch64 on 2026-10-18.

package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/btf"
	"github.com/cilium/ebpf/link"
	"golang.org/x/sys/unix"
)

// nativeEndian is the host byte order. BPF maps and their mmap views use host
// native endianness; both spike targets (aarch64, x86_64) are little-endian.
var nativeEndian = binary.NativeEndian

// Label bits, mirroring label_probe.bpf.c.
const (
	labelPublic  uint64 = 0x0
	labelPrivate uint64 = 0x1
	labelSecret  uint64 = 0x2
	labelUnknown uint64 = 0x4 // fail-closed bit
)

// Cause bits, mirroring TM_CAUSE_* in label_probe.bpf.c.
const (
	causeHookOpen    uint64 = 0x001
	causeHookReceive uint64 = 0x002
	causeHookMmap    uint64 = 0x004
	causeXattr       uint64 = 0x010
	causeDir         uint64 = 0x020
	causeCache       uint64 = 0x040
	causeError       uint64 = 0x080
	causeBacking     uint64 = 0x100
	causeEacces      uint64 = 0x200
	causeNobuf       uint64 = 0x400
)

// causeString renders a cause word for the serial log.
func causeString(c uint64) string {
	names := []struct {
		bit  uint64
		name string
	}{
		{causeHookOpen, "open"}, {causeHookReceive, "receive"},
		{causeHookMmap, "mmap"}, {causeXattr, "xattr"}, {causeDir, "dir"},
		{causeCache, "cache"}, {causeError, "error"}, {causeBacking, "backing"},
		{causeEacces, "eacces"}, {causeNobuf, "nobuf"},
	}
	var parts []string
	for _, n := range names {
		if c&n.bit != 0 {
			parts = append(parts, n.name)
		}
	}
	if len(parts) == 0 {
		return "none"
	}
	return strings.Join(parts, "|")
}

// Layout constants, mirroring label_probe.bpf.c.
const (
	maxSlots = 64
	maxRoots = 8
	rootLen  = 128
	// sessionLevel is the ancestor cgroup level at which a spike session is
	// identified: a transient scope under tmspike.slice (dash-free, so systemd
	// does not nest it under a parent slice) sits at level 2 (root=0,
	// slice=1, scope=2). See spike/README.md.
	sessionLevel uint32 = 2
)

// tmRoot mirrors `struct tm_root` in label_probe.bpf.c. The trailing pad makes
// the Go size equal the C size (128 + 8 + 4, padded to an 8-byte boundary = 144)
// so cilium/ebpf marshals it into the array value without a size mismatch.
type tmRoot struct {
	Prefix [rootLen]byte
	Label  uint64
	Len    uint32
	_      [4]byte
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func bpffsRoot() string { return envOr("TM_BPFFS", "/sys/fs/bpf/tidemark") }
func bpfObjPath() string {
	return envOr("TM_BPF_OBJ", "./label_probe.bpf.o")
}

// attachMmapFile reports whether the mmap_file program should be loaded and
// attached. TM_ATTACH_MMAP_FILE=0 removes it from the collection before load
// (escape hatch for the NULL guard described in label_probe.bpf.c), so a
// verifier rejection of that one program cannot block the whole suite.
func attachMmapFile() bool { return envOr("TM_ATTACH_MMAP_FILE", "1") != "0" }

func pinPath(parts ...string) string {
	return filepath.Join(append([]string{bpffsRoot()}, parts...)...)
}

// labelFromName converts a human label name to its bits.
func labelFromName(name string) (uint64, error) {
	switch name {
	case "public":
		return labelPublic, nil
	case "private":
		return labelPrivate, nil
	case "secret":
		return labelSecret, nil
	default:
		return 0, fmt.Errorf("unknown label %q (want public|private|secret)", name)
	}
}

// kernelHasRealDataInode reports whether the running kernel exposes the
// bpf_real_data_inode kfunc, by looking it up in kernel BTF. This drives the
// have_real_data_inode gate so the verified program never contains a call to a
// kfunc the kernel lacks. See spike/README.md for the fallback behaviour.
func kernelHasRealDataInode() bool {
	spec, err := btf.LoadKernelSpec()
	if err != nil {
		// If BTF cannot be read, assume the kfunc is absent and take the
		// fallback path; a kernel without BTF cannot load the probe anyway.
		return false
	}
	var fn *btf.Func
	if err := spec.TypeByName("bpf_real_data_inode", &fn); err != nil {
		return false
	}
	return true
}

// cmdLoad loads the object, configures the config globals, attaches the
// sleepable LSM programs, and pins the maps and links under bpffs. It first
// removes any stale pins left by an interrupted run, so programs from an older
// run are detached and Pin does not fail with EEXIST.
func cmdLoad(_ []string) (err error) {
	// On any failure, remove whatever was already pinned so no program stays
	// attached after an aborted load.
	defer func() {
		if err != nil {
			_ = os.RemoveAll(bpffsRoot())
		}
	}()
	if err := os.RemoveAll(bpffsRoot()); err != nil {
		return fmt.Errorf("remove stale pins under %s: %w", bpffsRoot(), err)
	}
	for _, d := range []string{"", "maps", "links"} {
		if err := os.MkdirAll(pinPath(d), 0o700); err != nil {
			return fmt.Errorf("mkdir pin dir %q: %w", d, err)
		}
	}

	spec, err := ebpf.LoadCollectionSpec(bpfObjPath())
	if err != nil {
		return fmt.Errorf("load spec %s: %w", bpfObjPath(), err)
	}

	have := kernelHasRealDataInode()
	// TM_DISABLE_INODE_CACHE=1 measures the no-cache labeller: the BPF side
	// keeps the cache path compiled but gates it on this variable, so the whole
	// suite can be run with the cache off. That is the evidence for the C2
	// lookup-method decision (docs/adr/2026-10-23_Phase0-Founder-Decisions_v1.md
	// FD-1); it also shows whether a Linux 7.1 kernel, which lacks the kfunc,
	// could run the labeller unchanged.
	if envOr("TM_DISABLE_INODE_CACHE", "0") == "1" {
		have = false
		fmt.Printf("INFO TM_DISABLE_INODE_CACHE=1: running without the inode label cache\n")
	}
	fmt.Printf("INFO bpf_real_data_inode present in kernel BTF: %v (cache active: %v)\n",
		kernelHasRealDataInode(), have)

	cfg := map[string]any{
		"session_level":        sessionLevel,
		"have_real_data_inode": boolToU8(have),
	}
	for name, val := range cfg {
		v, ok := spec.Variables[name]
		if !ok {
			return fmt.Errorf("config variable %q missing from object", name)
		}
		if err := v.Set(val); err != nil {
			return fmt.Errorf("set %s: %w", name, err)
		}
	}

	if !attachMmapFile() {
		delete(spec.Programs, "tm_mmap_file")
		fmt.Println("INFO tm_mmap_file removed from the collection (TM_ATTACH_MMAP_FILE=0)")
	}

	coll, err := ebpf.NewCollectionWithOptions(spec, ebpf.CollectionOptions{})
	if err != nil {
		var ve *ebpf.VerifierError
		if errors.As(err, &ve) {
			return fmt.Errorf("new collection (verifier log follows):\n%+v", ve)
		}
		return fmt.Errorf("new collection: %w", err)
	}
	defer coll.Close()

	// Pin the maps the helpers and the mmap reader need.
	for _, name := range []string{"sessions", "labels", "causes", "roots", "events"} {
		m, ok := coll.Maps[name]
		if !ok {
			return fmt.Errorf("map %q missing from object", name)
		}
		if err := m.Pin(pinPath("maps", name)); err != nil {
			return fmt.Errorf("pin map %s: %w", name, err)
		}
	}

	// Attach and pin the sleepable LSM programs so they stay attached after
	// this process exits (pinned bpf_links survive the loader).
	progs := []string{"tm_file_open", "tm_file_receive"}
	if attachMmapFile() {
		progs = append(progs, "tm_mmap_file")
	}
	for _, name := range progs {
		p, ok := coll.Programs[name]
		if !ok {
			return fmt.Errorf("program %q missing from object", name)
		}
		l, err := link.AttachLSM(link.LSMOptions{Program: p})
		if err != nil {
			return fmt.Errorf("attach %s: %w", name, err)
		}
		if err := l.Pin(pinPath("links", name)); err != nil {
			l.Close()
			return fmt.Errorf("pin link %s: %w", name, err)
		}
		// The pin keeps the link (and thus the attachment) alive after this
		// process exits; our own fd can be released.
		l.Close()
	}

	fmt.Println("INFO probe loaded, maps and links pinned under", bpffsRoot())
	return nil
}

func boolToU8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

// cmdUnload removes every pin, which detaches the programs.
func cmdUnload(_ []string) error {
	root := bpffsRoot()
	if err := os.RemoveAll(root); err != nil {
		return fmt.Errorf("remove pins: %w", err)
	}
	fmt.Println("INFO probe unloaded, pins removed from", root)
	return nil
}

// openPinnedMap opens a pinned map read-write.
func openPinnedMap(name string) (*ebpf.Map, error) {
	return ebpf.LoadPinnedMap(pinPath("maps", name), nil)
}

// mapLabelsRO returns a PROT_READ, MAP_SHARED mapping of the pinned labels
// array. The caller must unix.Munmap it.
//
// Why a READ-WRITE descriptor: at v7.2, do_mmap (mm/mmap.c, MAP_SHARED case)
// clears VM_SHARED when the file lacks FMODE_WRITE, and bpf_map_mmap
// (kernel/bpf/syscall.c) then rejects the mapping with -EINVAL because it
// requires VM_SHARED. A descriptor obtained with BPF_F_RDONLY (cilium/ebpf
// LoadPinOptions{ReadOnly: true}) therefore cannot be mapped at all, and
// cilium/ebpf's Map.Memory() PROT_READ fallback (memory.go, v0.22.0) does not
// help, because it targets frozen or BPF_F_RDONLY_PROG maps, not read-only
// descriptors. With a read-write descriptor, a PROT_READ mapping keeps
// VM_SHARED, and bpf_map_mmap clears VM_MAYWRITE, so the MAPPING cannot later
// be made writable. The descriptor itself is not read-only; check-ro-mmap
// records both facts. This is a label-transport design input for Phase 1.
func mapLabelsRO() ([]byte, error) {
	m, err := openPinnedMap("labels")
	if err != nil {
		return nil, fmt.Errorf("open labels map: %w", err)
	}
	// The mapping keeps its own reference to the map after Close.
	defer m.Close()
	// 64 slots x 8 bytes = 512 bytes, inside one page on both targets.
	b, err := unix.Mmap(m.FD(), 0, os.Getpagesize(), unix.PROT_READ, unix.MAP_SHARED)
	if err != nil {
		return nil, fmt.Errorf("mmap labels map PROT_READ: %w", err)
	}
	return b, nil
}

// readLabelRO reads the label word for a slot through a read-only mmap of the
// pinned BPF_F_MMAPABLE array, which is the access pattern the gateway and
// broker are meant to use in production (see mapLabelsRO for the descriptor
// caveat).
func readLabelRO(slot uint32) (uint64, error) {
	if slot >= maxSlots {
		return 0, fmt.Errorf("slot %d out of range", slot)
	}
	b, err := mapLabelsRO()
	if err != nil {
		return 0, err
	}
	defer unix.Munmap(b)
	off := int(slot) * 8
	return nativeEndian.Uint64(b[off : off+8]), nil
}

// readCause reads the cause bitmask for a slot (root only; plain map lookup).
func readCause(slot uint32) (uint64, error) {
	if slot >= maxSlots {
		return 0, fmt.Errorf("slot %d out of range", slot)
	}
	m, err := ebpf.LoadPinnedMap(pinPath("maps", "causes"), &ebpf.LoadPinOptions{ReadOnly: true})
	if err != nil {
		return 0, fmt.Errorf("open causes map: %w", err)
	}
	defer m.Close()
	var v uint64
	if err := m.Lookup(slot, &v); err != nil {
		return 0, fmt.Errorf("lookup cause slot %d: %w", slot, err)
	}
	return v, nil
}
