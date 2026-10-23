# Language: Go for daemons, CLI, verifier and benchmark; C for BPF

- **Status:** Proposed (Phase 0 draft, revised 2026-10-16 after review; founder
  review in week 3)
- **Date:** 2026-10-16
- **Phase:** 0 (Foundations and kill gate)
- **Deciders:** Founder

## Context

Tidemark is about 4,000 lines of product code plus roughly 1,500 lines of
benchmark, puppet and demo code, built and reviewed by one part-time maintainer
with AI assistance. The maintainer personally reviews every security-critical
diff (the BPF program, the gateway copy loop, the staging merge, the ledger). The
language choice must serve a solo founder, cross-compile cleanly to linux/arm64
and linux/amd64, talk to systemd over D-Bus, load BPF, and produce a small
cross-platform verifier that runs on macOS, Linux and Windows.

Constraints and evidence:

- The components needed are: a BPF loader, Landlock in strict mode, systemd
  transient units over D-Bus, an HTTP proxy, an OpenAI-compatible broker, DSSE
  and in-toto handling, and a static offline verifier for three operating
  systems.
- The ecosystem the concept relies on is Go-first for these needs:
  `cilium/ebpf` with `bpf2go` (MIT), `go-landlock` (MIT), `coreos/go-systemd`
  and `godbus`, and `secure-systems-lab/go-securesystemslib` for DSSE plus the
  in-toto attestation Go bindings. Source:
  `docs/research/2026-10-16_Tidemark-Concept_v1.md`, stack section.
- The verifier must ship as static binaries for macOS, Linux and Windows. Go
  cross-compiles to all three from one toolchain with `CGO_ENABLED=0`.
- The feasibility research recommended "a small Rust or Go daemon with
  declarative policy files" for the egress enforcer, so both languages were
  viable; the deciding factors are the Go-first BPF, Landlock and attestation
  libraries and the single-language simplicity for a solo maintainer. Source:
  `docs/research/2026-10-16_OS-Feasibility_v1.md`.
- BPF enforcement code itself is written in C and compiled with clang/LLVM; no
  high-level language replaces that. The question is only the host-side language.

## Decision

Use **Go** for every daemon (`tidemarkd`, `tidemark-gw`, `tidemark-broker`,
`tidemark-ledgerd`, `tidemark-shim`), the `tidemark` CLI, the `tidemark-verify`
verifier and the benchmark harness. Write the BPF enforcement program in **C**,
compiled with clang in a Fedora toolbox container and loaded through
`cilium/ebpf` with `bpf2go`.

- **Go version pin.** Every module declares `go 1.25.0` as its minimum
  language version (the spike module already does, and `cilium/ebpf` v0.22.0
  requires it) and `toolchain go1.27.2` as the toolchain that builds it. On
  2026-10-16 the supported releases listed by `https://go.dev/dl/?mode=json`
  are go1.27.2 and go1.26.9. Renovate proposes toolchain updates, which merge
  only after CI passes. The spike's `go.mod` does not yet carry the `toolchain`
  line; its owner adds it.
- **Where Go runs.** Go is not installed on the macOS development host, and
  this decision does not install it there. Builds and tests run inside the
  Fedora Lima VM (aarch64) and in CI (x86_64).
- Daemons build with `CGO_ENABLED=0` and cross-compile to linux/arm64 and
  linux/amd64. The verifier additionally builds for darwin and windows.
- The Go module path is `github.com/BitecodesHub/tidemark` (resolved 2026-10-23
  when the repository was created; decisions ADR PD-10). It was
  `github.com/agentos-tidemark/tidemark` while the repository was a placeholder.
- Reusable evaluators live in `pkg/labels` and `pkg/clearance` so an OpenShell
  middleware adapter can import them in Phase 6.

## Consequences

- One language across daemons, CLI, verifier and benchmark lowers context
  switching for a solo maintainer and keeps the security review surface uniform.
- The BPF loader, Landlock, systemd and attestation libraries are mature in Go,
  so little has to be written from scratch.
- The static verifier ships to all three desktop operating systems from one
  toolchain, which matters because the public proof path includes running the
  verifier off the machine.
- Go's garbage collector is acceptable here: the daemons are small and the hot
  path (the gateway copy loop, the broker body read) is bounded and reviewed.
  The per-chunk label read is a plain atomic load from an mmapped array, not an
  allocation.
- A dual-language codebase (C for BPF, Go for the rest) is unavoidable because
  BPF is C. This is the only language split, and the C surface is deliberately
  small (about 250 lines) and reviewed line by line.
- Memory-safety guarantees are weaker than Rust's. The mitigation is a small root
  surface (only `tidemarkd` is root), unprivileged daemons under their own UIDs
  with systemd sandboxing, fuzzing of every parser, and `gosec`, `govulncheck`,
  `staticcheck` and CodeQL in CI.

## Alternatives considered

- **Rust.** Rejected for v0.1, though it was a reasonable candidate (and is what
  OpenShell uses). The BPF loader, Landlock, systemd D-Bus and in-toto/DSSE
  ecosystems the concept depends on are more mature and better documented in Go,
  cross-compilation of the three-OS verifier is simpler, and a single maintainer
  moves faster in one familiar language. The memory-safety advantage does not
  outweigh these for a 4,000-line codebase with a small, heavily reviewed unsafe
  surface. Rust stays open for a future contributor-led component.
- **C or C++ for the daemons.** Rejected: far larger memory-safety risk across
  the whole codebase, slower to write, and no gain over Go given BPF is already
  isolated to its own small C program.
- **Python or another scripting language for the CLI and benchmark.** Rejected:
  it would add a second runtime and dependency surface on the image, and the
  benchmark must be reproducible and pinned by digest, which static Go binaries
  serve better.

## References

- `docs/research/2026-10-16_Tidemark-Concept_v1.md` (stack section, reuse notes,
  per-primitive build estimates).
- `docs/research/2026-10-16_OS-Feasibility_v1.md` (the "small Rust or Go daemon"
  recommendation).
- `docs/adr/2026-10-16_Labeller-BPF-LSM_v1.md` (the BPF C program this ADR keeps
  in C).
- Go release list, `https://go.dev/dl/?mode=json`, checked 2026-10-16.
- `spike/go.mod` (module minimum `go 1.25.0`).
