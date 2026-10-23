// AgentOS Tidemark - Phase 0 kernel spike module.
//
// This is a self-contained module so the throwaway spike never pulls into the
// product module. The module path follows the repository
// github.com/BitecodesHub/tidemark (decided 2026-10-23, decisions ADR PD-10);
// record any change in spike/README.md too.
//
// go.sum is committed because the first in-VM build (2026-10-18) resolved and
// verified it; the Makefile `build` target runs `go build` with
// GOFLAGS=-mod=mod. The pinned versions below were read from the upstream
// release metadata on 2026-10-16:
//   * github.com/cilium/ebpf v0.22.0 (released 2026-06-29; go.mod requires
//     go 1.25.0, golang.org/x/sys v0.43.0, golang.org/x/sync v0.20.0).
//   * golang.org/x/sys v0.43.0.
module github.com/BitecodesHub/tidemark/spike

go 1.25.0

require (
	github.com/cilium/ebpf v0.22.0
	golang.org/x/sys v0.43.0
)

// Transitive requirements of cilium/ebpf are added by `go mod tidy` inside the
// VM. They are left out here on purpose so this file stays reviewable by hand.
