// SPDX-License-Identifier: Apache-2.0
//
// Extra measurement helpers added for the Phase 0 items that the first run
// left open: the observer-side half of the race test, the labelled/unlabelled
// open() latency baseline (checklist D5), xattr removal (checklist D4 "xattr
// invalidation"), and the "session cannot move itself out" attempt required by
// acceptance criterion AC3.
//
// These helpers only observe and report; like every other subcommand they never
// deny a syscall.

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/cilium/ebpf"
	"golang.org/x/sys/unix"
)

// cmdClearSlot zeroes one session slot's label and cause words without
// re-registering the cgroup. The race and invalidation scenarios call it between
// iterations so each open() under test starts from an unraised label; a label
// left over from the previous iteration would make the raise trivially true and
// the test meaningless.
func cmdClearSlot(args []string) error {
	fs := flag.NewFlagSet("clear-slot", flag.ExitOnError)
	slot := fs.Uint("slot", 0, "label slot index")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *slot >= maxSlots {
		return fmt.Errorf("--slot %d out of range", *slot)
	}
	for _, name := range []string{"labels", "causes"} {
		m, err := openPinnedMap(name)
		if err != nil {
			return err
		}
		err = m.Update(uint32(*slot), uint64(0), ebpf.UpdateAny)
		m.Close()
		if err != nil {
			return fmt.Errorf("clear %s slot %d: %w", name, *slot, err)
		}
	}
	fmt.Printf("INFO cleared slot %d\n", *slot)
	return nil
}

// cmdRemoveXattr removes an xattr as root, outside any session, and verifies the
// removal. The invalidation scenario needs a removal that definitely happened, so
// a stale label is attributable to the kernel, not to a failed setup step.
func cmdRemoveXattr(args []string) error {
	fs := flag.NewFlagSet("removexattr", flag.ExitOnError)
	path := fs.String("path", "", "file")
	name := fs.String("name", "", "xattr name")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *path == "" || *name == "" {
		return errors.New("--path and --name are required")
	}
	if err := unix.Removexattr(*path, *name); err != nil {
		return fmt.Errorf("setup: removexattr %s on %s: %s", *name, *path, errnoName(err))
	}
	if _, err := unix.Getxattr(*path, *name, make([]byte, 8)); !errors.Is(err, unix.ENODATA) {
		return fmt.Errorf("setup: %s still readable after removal (err=%v)", *name, err)
	}
	fmt.Printf("removexattr %s on %s: ok (absent after removal)\n", *name, *path)
	return nil
}

// cmdTryMove writes this process's own pid into the target cgroup's
// cgroup.procs and reports the outcome. Acceptance criterion AC3 requires that a
// session process cannot move itself out of the session cgroup when the scope
// runs with Delegate=no; the scenario runs this as the capability-less uid 65534
// principal and requires the kernel to refuse.
func cmdTryMove(args []string) error {
	fs := flag.NewFlagSet("try-move", flag.ExitOnError)
	cgroup := fs.String("cgroup", "", "target cgroupfs directory")
	require := fs.String("require", "", "required outcome: ok or an errno name")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *cgroup == "" {
		return errors.New("--cgroup is required")
	}
	if err := checkSessionContext(); err != nil {
		return err
	}

	dst := *cgroup + "/cgroup.procs"
	outcome := "ok"
	if err := os.WriteFile(dst, []byte(strconv.Itoa(os.Getpid())), 0); err != nil {
		outcome = errnoName(err)
		fmt.Printf("try-move %s: errno=%s\n", dst, outcome)
	} else {
		fmt.Printf("try-move %s: ok\n", dst)
	}

	// Report where the process actually is, so a move that "succeeded" without
	// changing the cgroup cannot be misread.
	if own, err := ownCgroup(); err == nil {
		fmt.Printf("try-move own cgroup now: %s\n", own)
	}
	if *require != "" && *require != outcome {
		return fmt.Errorf("try-move %s: outcome %s, required %s", dst, outcome, *require)
	}
	return nil
}

// cmdMkdirTry attempts to create a directory inside the target cgroup as the
// calling (unprivileged) process. Delegate=no leaves the cgroupfs subtree
// root-owned, so a session must not be able to restructure its own placement.
func cmdMkdirTry(args []string) error {
	fs := flag.NewFlagSet("try-mkdir", flag.ExitOnError)
	dir := fs.String("dir", "", "directory to create")
	require := fs.String("require", "", "required outcome: ok or an errno name")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *dir == "" {
		return errors.New("--dir is required")
	}
	outcome := "ok"
	err := unix.Mkdir(*dir, 0755)
	if err != nil {
		outcome = errnoName(err)
		fmt.Printf("try-mkdir %s: errno=%s\n", *dir, outcome)
	} else {
		fmt.Printf("try-mkdir %s: ok\n", *dir)
	}
	if *require != "" && *require != outcome {
		return fmt.Errorf("try-mkdir %s: outcome %s, required %s", *dir, outcome, *require)
	}
	return nil
}

// latencyStat holds the reported distribution of one open()/close() loop.
type latencyStat struct {
	n      int
	minNs  int64
	maxNs  int64
	p50Ns  int64
	p99Ns  int64
	meanNs float64
}

func summarize(samples []int64) latencyStat {
	sort.Slice(samples, func(i, j int) bool { return samples[i] < samples[j] })
	s := latencyStat{n: len(samples)}
	if len(samples) == 0 {
		return s
	}
	var sum int64
	for i, v := range samples {
		sum += v
		if i == 0 {
			s.minNs = v
		}
		s.maxNs = v
	}
	s.meanNs = float64(sum) / float64(len(samples))
	idx := func(q float64) int64 {
		i := int(q * float64(len(samples)-1))
		if i < 0 {
			i = 0
		}
		if i >= len(samples) {
			i = len(samples) - 1
		}
		return samples[i]
	}
	s.p50Ns = idx(0.50)
	s.p99Ns = idx(0.99)
	return s
}

// cmdBenchOpen measures the latency of open(O_RDONLY)+close() on one path. The
// scenario runs it four times - outside a session on an unlabelled file, inside a
// session on an unlabelled file, inside a session on an xattr-labelled file, and
// inside a session under a directory rule - which is the checklist D5 baseline
// (no target is set in Phase 0; the numbers are recorded for the C2 cache
// decision).
//
// The first open() of a session is reported separately: it is the event that
// raises the label, and the sleepable BPF programs do the xattr and path lookups
// inside it.
func cmdBenchOpen(args []string) error {
	fs := flag.NewFlagSet("bench-open", flag.ExitOnError)
	path := fs.String("path", "", "file to open repeatedly")
	slot := fs.Uint("slot", 0, "label slot index (for the raised check)")
	iterations := fs.Int("iterations", 2000, "number of open()/close() pairs")
	warmup := fs.Int("warmup", 50, "iterations excluded from the distribution")
	checkRaised := fs.Bool("check-raised", false, "after each open, read the label and require it raised")
	clearEach := fs.Bool("clear-each", false, "clear the slot before each measured open (raises every time)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *path == "" {
		return errors.New("--path is required")
	}
	if *iterations <= 0 {
		return errors.New("--iterations must be positive")
	}
	if err := checkSessionContext(); err != nil {
		return err
	}

	var firstNs int64
	var samples []int64
	var raisedSamples []int64
	misses := 0
	total := *warmup + *iterations
	for i := 0; i < total; i++ {
		measured := i >= *warmup
		if *clearEach && *checkRaised && measured {
			if err := clearSlotOnly(uint32(*slot)); err != nil {
				return err
			}
		}
		t0 := time.Now()
		fd, err := unix.Open(*path, unix.O_RDONLY|unix.O_CLOEXEC, 0)
		elapsed := time.Since(t0).Nanoseconds()
		if err != nil {
			return fmt.Errorf("setup: open %s: %s", *path, errnoName(err))
		}
		if *checkRaised {
			tR := time.Now()
			v, lerr := readLabelRO(uint32(*slot))
			raisedSamples = append(raisedSamples, time.Since(tR).Nanoseconds())
			if lerr != nil {
				return fmt.Errorf("setup: read label: %w", lerr)
			}
			if measured && v == labelPublic {
				misses++
			}
		}
		if err := unix.Close(fd); err != nil {
			return fmt.Errorf("setup: close: %s", errnoName(err))
		}
		switch {
		case i == 0:
			firstNs = elapsed
		case measured:
			samples = append(samples, elapsed)
		}
	}

	st := summarize(samples)
	fmt.Printf("BENCH path=%s iterations=%d warmup=%d first_open_ns=%d min_ns=%d p50_ns=%d p99_ns=%d max_ns=%d mean_ns=%.0f\n",
		*path, st.n, *warmup, firstNs, st.minNs, st.p50Ns, st.p99Ns, st.maxNs, st.meanNs)
	if *checkRaised {
		rs := summarize(raisedSamples)
		fmt.Printf("BENCH_READLABEL slot=%d p50_ns=%d p99_ns=%d misses=%d of %d\n",
			*slot, rs.p50Ns, rs.p99Ns, misses, st.n)
		if misses > 0 {
			return fmt.Errorf("%d of %d measured opens did not show a raised label", misses, st.n)
		}
	}
	return nil
}

// clearSlotOnly is cmdClearSlot without the human-readable line, for the bench
// loop's per-iteration reset.
func clearSlotOnly(slot uint32) error {
	for _, name := range []string{"labels", "causes"} {
		m, err := openPinnedMap(name)
		if err != nil {
			return err
		}
		err = m.Update(slot, uint64(0), ebpf.UpdateAny)
		m.Close()
		if err != nil {
			return fmt.Errorf("clear %s slot: %w", name, err)
		}
	}
	return nil
}

// slotView holds the handles the race loops need for the whole run: a read-only
// mmap of the labels array and read-write handles of the label and cause arrays.
// readLabelRO and clearSlotOnly reopen and remap on every call, which is the
// honest access pattern for a single check but would dominate a 10,000-iteration
// loop; the race test reuses the handles instead.
type slotView struct {
	mapping []byte
	labels  *ebpf.Map
	causes  *ebpf.Map
	slot    uint32
}

func newSlotView(slot uint32) (*slotView, error) {
	if slot >= maxSlots {
		return nil, fmt.Errorf("slot %d out of range", slot)
	}
	m, err := openPinnedMap("labels")
	if err != nil {
		return nil, err
	}
	b, err := unix.Mmap(m.FD(), 0, os.Getpagesize(), unix.PROT_READ, unix.MAP_SHARED)
	if err != nil {
		m.Close()
		return nil, fmt.Errorf("mmap labels: %s", errnoName(err))
	}
	c, err := openPinnedMap("causes")
	if err != nil {
		unix.Munmap(b)
		m.Close()
		return nil, err
	}
	return &slotView{mapping: b, labels: m, causes: c, slot: slot}, nil
}

func (v *slotView) close() {
	if v == nil {
		return
	}
	if v.mapping != nil {
		_ = unix.Munmap(v.mapping)
	}
	if v.labels != nil {
		_ = v.labels.Close()
	}
	if v.causes != nil {
		_ = v.causes.Close()
	}
}

func (v *slotView) label() uint64 {
	off := int(v.slot) * 8
	return nativeEndian.Uint64(v.mapping[off : off+8])
}

func (v *slotView) clear() error {
	if err := v.labels.Update(v.slot, uint64(0), ebpf.UpdateAny); err != nil {
		return fmt.Errorf("clear label slot: %w", err)
	}
	if err := v.causes.Update(v.slot, uint64(0), ebpf.UpdateAny); err != nil {
		return fmt.Errorf("clear cause slot: %w", err)
	}
	return nil
}

// cmdRaceOpen is the self-half of the labeller ADR kill test: for N iterations,
// clear the session slot, open a labelled file, and read the label back in the
// SAME process immediately after open() returns. Any iteration that observes an
// unraised label is a miss, and the ADR requires zero misses. When --notify is
// given, the process also connects once and writes one byte per iteration after
// its own sample, so a separate observer process (cmdRaceObserve, standing in for
// the gateway) can sample the same raise. It runs inside the session as a
// grandchild of the process that joined.
func cmdRaceOpen(args []string) error {
	fs := flag.NewFlagSet("race-open", flag.ExitOnError)
	path := fs.String("path", "", "labelled file to open once per iteration")
	slot := fs.Uint("slot", 0, "label slot index")
	iterations := fs.Int("iterations", 10000, "iterations")
	expect := fs.String("expect", "0x1", "expected label after open() returns (hex)")
	notify := fs.String("notify", "", "unix socket of the observer process")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *path == "" {
		return errors.New("--path is required")
	}
	if *iterations <= 0 {
		return errors.New("--iterations must be positive")
	}
	if err := checkSessionContext(); err != nil {
		return err
	}
	want, ok, err := parseHex("--expect", *expect)
	if err != nil {
		return err
	}
	if !ok {
		want = labelPublic
	}

	view, err := newSlotView(uint32(*slot))
	if err != nil {
		return err
	}
	defer view.close()

	var s int
	if *notify != "" {
		s, err = unix.Socket(unix.AF_UNIX, unix.SOCK_STREAM|unix.SOCK_CLOEXEC, 0)
		if err != nil {
			return fmt.Errorf("setup: notify socket: %w", err)
		}
		defer unix.Close(s)
		var cerr error
		for i := 0; i < 400; i++ {
			cerr = unix.Connect(s, &unix.SockaddrUnix{Name: *notify})
			if cerr == nil || (!errors.Is(cerr, unix.ECONNREFUSED) && !errors.Is(cerr, unix.ENOENT)) {
				break
			}
			time.Sleep(50 * time.Millisecond)
		}
		if cerr != nil {
			return fmt.Errorf("setup: connect notify %s: %s", *notify, errnoName(cerr))
		}
	}

	var lat []int64
	misses := 0
	for i := 0; i < *iterations; i++ {
		if err := view.clear(); err != nil {
			return err
		}
		t0 := time.Now()
		fd, err := unix.Open(*path, unix.O_RDONLY|unix.O_CLOEXEC, 0)
		elapsed := time.Since(t0).Nanoseconds()
		if err != nil {
			return fmt.Errorf("setup: open %s: %s", *path, errnoName(err))
		}
		got := view.label()
		if err := unix.Close(fd); err != nil {
			return fmt.Errorf("setup: close: %s", errnoName(err))
		}
		if got != want {
			misses++
			if misses <= 5 {
				fmt.Printf("MISS iteration=%d label=0x%x expected=0x%x\n", i+1, got, want)
			}
		}
		lat = append(lat, elapsed)
		if s != 0 {
			if _, err := unix.Write(s, []byte{'r'}); err != nil {
				return fmt.Errorf("setup: notify write: %s", errnoName(err))
			}
			// Wait for the observer's ack before the next iteration clears the
			// slot. Without it the socket buffers several iterations ahead, the
			// observer samples a slot the writer has already cleared, and the
			// harness - not the kernel - reports a miss.
			ack := make([]byte, 1)
			if n, err := unix.Read(s, ack); err != nil {
				return fmt.Errorf("setup: notify ack: %s", errnoName(err))
			} else if n == 0 {
				return errors.New("setup: observer closed the connection")
			}
		}
	}
	st := summarize(lat)
	fmt.Printf("RACE_SELF iterations=%d misses=%d path=%s open_min_ns=%d open_p50_ns=%d open_p99_ns=%d open_max_ns=%d\n",
		*iterations, misses, *path, st.minNs, st.p50Ns, st.p99Ns, st.maxNs)
	if misses > 0 {
		return fmt.Errorf("%d of %d iterations did not show the label raised immediately after open()", misses, *iterations)
	}
	return nil
}

// cmdRaceObserve is the observer half: one long-lived process OUTSIDE the session
// (standing in for tidemark-gw) that accepts a single connection, then for N
// iterations waits for the in-session writer's byte and reads the label
// immediately afterwards. Zero misses required, same as the self half.
func cmdRaceObserve(args []string) error {
	fs := flag.NewFlagSet("race-observe", flag.ExitOnError)
	sock := fs.String("sock", "", "unix socket path to bind")
	slot := fs.Uint("slot", 0, "label slot index")
	iterations := fs.Int("iterations", 10000, "iterations")
	expect := fs.String("expect", "0x1", "expected label (hex)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *sock == "" {
		return errors.New("--sock is required")
	}
	want, ok, err := parseHex("--expect", *expect)
	if err != nil {
		return err
	}
	if !ok {
		want = labelPublic
	}

	view, err := newSlotView(uint32(*slot))
	if err != nil {
		return err
	}
	defer view.close()

	l, err := unix.Socket(unix.AF_UNIX, unix.SOCK_STREAM|unix.SOCK_CLOEXEC, 0)
	if err != nil {
		return fmt.Errorf("setup: socket: %w", err)
	}
	defer unix.Close(l)
	_ = os.Remove(*sock)
	defer os.Remove(*sock)
	if err := unix.Bind(l, &unix.SockaddrUnix{Name: *sock}); err != nil {
		return fmt.Errorf("setup: bind %s: %s", *sock, errnoName(err))
	}
	if err := unix.Listen(l, 1); err != nil {
		return fmt.Errorf("setup: listen: %w", err)
	}
	c, _, err := unix.Accept4(l, unix.SOCK_CLOEXEC)
	if err != nil {
		return fmt.Errorf("setup: accept: %s", errnoName(err))
	}
	defer unix.Close(c)

	misses := 0
	buf := make([]byte, 1)
	for i := 0; i < *iterations; i++ {
		n, err := unix.Read(c, buf)
		if err != nil {
			return fmt.Errorf("setup: read notify byte at iteration %d: %s", i+1, errnoName(err))
		}
		if n == 0 {
			return fmt.Errorf("setup: notify connection closed at iteration %d of %d", i+1, *iterations)
		}
		got := view.label()
		if _, err := unix.Write(c, []byte{'a'}); err != nil {
			return fmt.Errorf("setup: ack write: %s", errnoName(err))
		}
		if got != want {
			misses++
			if misses <= 5 {
				fmt.Printf("MISS observer iteration=%d label=0x%x expected=0x%x\n", i+1, got, want)
			}
		}
	}
	fmt.Printf("RACE_OBSERVE iterations=%d misses=%d slot=%d\n", *iterations, misses, *slot)
	if misses > 0 {
		return fmt.Errorf("observer saw %d of %d iterations with an unraised label", misses, *iterations)
	}
	return nil
}
