# THREAT_MODEL.md

**Status: v0 draft (Phase 0), written 16 October 2026.** This document is
tagged as v0 at the end of Phase 0 (target Sunday 8 November 2026). Later
revisions are versioned, and each revision adds an entry to the change log in
Section 12. Nothing here claims that a feature is built or tested: every
guarantee below is **designed, not implemented**, and the Status column in
Section 5 names the phase in which each one is to be implemented and measured.
Claims stay within the "not claimed" list of the concept document
(`docs/research/2026-10-16_Tidemark-Concept_v1.md`).

AgentOS Tidemark (working name; a trademark search is required before launch) is
a headless, signed Fedora `bootc` image whose distinctive feature is
kernel-recorded read-time labels for agent sessions, clearance-gated exits
re-checked per chunk, a label-constrained model broker, staged writes with label
stamping, and agent-unreachable signed receipts. The software ships as an RPM
that the image installs; full guarantees are designed for, and will only ever be
claimed on, the signed image.

v0.1 is a confidentiality preview. It does not prevent prompt injection, does
not resist root or kernel compromise, and does not hold external side effects.

---

## 1. Scope and system model

A **session** is one agent run and every descendant it starts: shells, MCP
servers, hooks, LSP servers, subprocesses and double-forked daemons. Tidemark
turns each session into one contained principal with its own cgroup subtree,
mount view, network namespace (loopback only), Landlock domain, two per-session
Unix sockets (one to the gateway, one to the broker), and a receipt identity.
A session exists only when it is started with `tidemark run`; processes started
any other way are not sessions (Section 6).

A **label** is a floating high-water mark on the session. It rises when any
session process reads labelled data and never falls during the session. A
**clearance** is the set of labels a destination is permitted to receive.
Bytes leave a session only while the session's current label is a subset of the
destination's clearance, re-checked before each forwarded chunk.

### 1.1 Trusted computing base, per guarantee

The **root TCB**, shared by every guarantee, is: the Linux kernel (including the
BPF verifier, Landlock, overlayfs and the cgroup freezer); systemd as PID 1,
which creates the session units; `tidemarkd` (the only root daemon) together
with the BPF programs it loads; the signed image and its build pipeline; and the
machine's firmware.

The unprivileged daemons (`tidemark-gw`, `tidemark-broker`, `tidemark-ledgerd`)
each run under their own UID with systemd sandboxing and parse untrusted input
outside root. Running outside root limits what a compromise of one of them can
reach, but it does **not** place them outside the TCB: several guarantees depend
on them, as the table shows.

| Guarantee | Trusted components beyond the root TCB |
|---|---|
| G1 Containment | `tidemark-shim`, which closes inherited descriptors and applies the Landlock domain before it execs the agent |
| G2 Single exit | `tidemark-shim`; `tidemark-gw` for the rule that clearance is evaluated before any DNS query |
| G3 Label before read | None for file reads. `tidemark-gw` for source labels, because it requests the raise before it relays the first byte from a private source |
| G4 Clearance on every chunk | `tidemark-gw` and `tidemark-broker`, including their HTTP, CONNECT and JSON parsers and their read-only view of the label array |
| G5 Model channel | `tidemark-broker` (placement, field and tool allowlist) and the systemd credentials that hold the model keys |
| G6 Staged writes carry labels | The `tidemarkd` commit, stamp, accept and abort paths |
| G7 Evidence against the agent | `tidemark-ledgerd`, the signing key and `tidemark-verify`. The **content** of a receipt is only as true as its event sources, which are `tidemarkd`, `tidemark-gw` and `tidemark-broker`: a compromised source can submit false events, and G7 then signs them faithfully |
| G8 Stop | `tidemark-gw` and `tidemark-broker`, which switch to deny-all |
| G9 Declassification | The `tidemarkd` caller checks and `tidemark-broker`, which binds and consumes the release token |

---

## 2. Assets

| # | Asset | Why it matters |
|---|---|---|
| A1 | User secrets a profile exposes to a session: SSH and GPG keys, cloud credentials, password stores, API tokens | Direct compromise of other systems if exfiltrated |
| A2 | Private data in labelled files and directories | The confidentiality the product exists to protect |
| A3 | Private data read from destinations tagged as private sources (for example a mail API) | Network-sourced private data must taint the session like a file read |
| A4 | Agent persistent state: memories, skills and hooks | A laundering path and a persistence path across sessions |
| A5 | Model API keys, held only by the broker (from v0.2, credentials for declared APIs, held only by the gateway) | Keys must never enter a session |
| A6 | Integrity of the workspace and host: files, auto-exec paths, skill directories, immutable `/usr` | Implants and silent tampering |
| A7 | External accounts reachable through cleared destinations | v0.1 does not hold effects on them; v0.2 holds irreversible effects on declared APIs |
| A8 | Integrity of the evidence: ledger segments, receipts, the signing key | Agents have tampered with logs before |
| A9 | The human's attention | Release requests must be rare and exact, not a prompt flood; a release is only as sound as the human's reading of what is shown (Section 6) |

---

## 3. Adversaries

### 3.1 In scope

- **AD1 Prompt injector.** Controls content the agent reads: web pages, issues,
  READMEs, emails, API responses. The injected instruction is the normal case,
  not an edge case.
- **AD2 Malicious in-session component.** A compromised or hostile skill, MCP
  server, plugin, hook or agent package running inside the session.
- **AD3 Deceptive agent or model.** Spoofs tool calls, attempts to tamper with
  logs, tries to persist through implants or memory, and attempts to lower its
  own label.
- **AD4 Cross-vendor or cross-session launderer.** A second agent from another
  vendor, or a later session of the same agent, used to carry data out through
  files or persistent memory.

The design assumption for the benchmark is the strongest form of AD1 through
AD3 combined: a **fully compromised agent** driven by scripted puppet backends,
with all of the agent's own permission prompts disabled.

### 3.2 Out of scope (v0.1)

- A root or kernel-level attacker. A kernel exploit or root compromise defeats
  every guarantee. The kernel stays shared; microVM principals are a backlog
  item only.
- A physical attacker.
- The machine owner. Receipts are not proof against the owner until TPM
  attestation (Phase 7).
- A destination the user cleared that then misuses the data.
- Network attackers beyond what standard TLS covers.
- A compromised Fedora base image, COPR, GHCR or Sigstore. Mitigated only by
  signature verification and digest pinning, not eliminated.
- The human's own tools **outside** sessions that strip `security.*` xattrs and
  so silently remove labels. This is a documented limitation (see Section 7),
  narrowed by kernel-evaluated directory rules.

---

## 4. Trust boundaries

| Boundary | Enforcement | Residual |
|---|---|---|
| Session to kernel | Namespaces with creation denied (user namespaces included), a cgroup subtree with `Delegate=no`, a Landlock domain at ABI 9 or later, seccomp deny list, BPF LSM programs | The shared kernel is trusted; a kernel bug crosses this boundary |
| Session to Tidemark services | Exactly two per-session Unix sockets; the socket identifies the session; gateway and broker hold read-only mmap of the session label word; the control socket and ledger are unreachable from sessions | Socket peer-credential checks must be correct; gateway and broker parsers face untrusted input (Section 8) |
| Session to network | No interface except loopback; every exit passes through `tidemark-gw` or `tidemark-broker`; clearance is evaluated on the requested hostname before any DNS query; each chunk re-checked against the current label | Covert and timing channels to cleared destinations; the choice and timing of cleared hostnames are visible to the gateway's upstream resolver |
| Session to filesystem | `/usr` read-only; every writable path is an overlay upper; read binds limited to what the profile declares; labels and label maps change only with `CAP_SYS_ADMIN`, from outside sessions | Tools outside sessions can drop xattrs (A6, see Section 7) |
| Session to other sessions | Separate mount namespaces and Landlock domains; ptrace and cross-session signals blocked; uncommitted writes invisible across sessions | One active session per profile in v0.1, but several profiles can run at once, so covert and timing channels between concurrent sessions remain (Section 6) |
| Exit plane to ledger | A socket filtered by peer UID (`SO_PEERCRED`); ledger segments append-only by kernel attribute (`FS_APPEND_FL`) against the ledger daemon | Root can still alter segments; receipt content trusts the event sources (Section 1.1) |
| Human to machine | `commit`, `accept`, `release` and `stop` accepted only from the human UID or root, from outside the agents slice | Depends on correct caller checks |
| Machine to outside world | The BYOC canary and host-side sinks are the primary external check; receipts verify offline while the key stays secret | Third-party proof waits for Phase 7 attestation |
| Machine to cleared destinations | Clearing a destination is the user's trust decision | A cleared destination can misuse what it receives (out of scope, Section 3.2) |
| Build pipeline to image | Keyless `cosign` signing, a signed RPM, an SBOM, pinned digests; users verify with `cosign` | Supply-chain compromise above the signature is out of scope |

---

## 5. Guarantees G1 to G9 and the hosts that provide each

These are **preview** guarantees. In Phase 0 they are **designed only**. Each
becomes a claim only after its MVP phase implements and measures it (G1 to G3 in
Phase 1, G4 and G5 in Phase 2, G6 and G7 in Phase 3, G8 and G9 across Phases 2
and 3). Even then, a guarantee holds only when `tidemark doctor` passes and the
kernel is not compromised.

**Host types.**

- **Signed image** - the cosign-verified `bootc` image. Pins the kernel and
  kfuncs, guarantees a Landlock ABI 9 floor, owns every agent launch path,
  keeps credentials inside daemons, fixes the staging and ledger layout, and
  binds receipts to a signed digest. This is the only host where the full set
  of guarantees is designed to be claimed.
- **RPM on stock Fedora 45** (Atomic or Server, adopted with `bootc switch` or
  `bootc install to-existing-root`, or the RPM alone). Guarantees are only
  those `tidemark doctor` reports for that host: the kernel version, LSM list,
  kfunc availability, cgroup v2, overlayfs and `/var` xattr/`FS_APPEND_FL`
  support are checked at run time, not assumed. The receipt records the doctor
  result, so every receipt states the guarantee level it was produced under.
- **Elsewhere** (macOS, Windows, other Linux without the RPM). Only
  `tidemark-verify` runs, offline. No enforcement guarantee. macOS host model
  endpoints are recorded as **unattested**.

| ID | Guarantee (preview) | Mechanism | Designed level (signed image) | RPM on stock Fedora 45 | Elsewhere | Status |
|---|---|---|---|---|---|---|
| **G1 Containment** | Every descendant stays in the session cgroup subtree, network namespace and mount view, and is killed on stop. No namespace creation (user namespaces included), no io_uring. | systemd transient unit, cgroup v2 `Delegate=no`, `RestrictNamespaces=yes`, seccomp | Full | Doctor-reported (depends on systemd and cgroup v2 config) | None | Designed; not implemented (Phase 1) |
| **G2 Single exit** | A session has no network path except `tidemark-gw` and `tidemark-broker`. DNS, UDP, QUIC, IP literals and machine-local Unix sockets are unavailable. Clearance is evaluated on the requested hostname before any DNS query; a hostname outside the session's current clearance is never resolved. | Private network namespace (loopback only), per-session sockets, Landlock `LANDLOCK_ACCESS_FS_RESOLVE_UNIX`, gateway-side resolution after the clearance check | Full | Doctor-reported | None | Designed; not implemented (Phase 1; gateway rule in Phase 2) |
| **G3 Label before read** | The session label rises before the call returns when any session process opens a labelled file (by xattr or directory rule) or receives a descriptor to one. Reading from a private-source destination raises the label before the first response byte reaches the session. Descriptors are not inherited across the shim exec. A session cannot lower any label. | Sleepable BPF LSM on `file_open` and `file_receive`; `bpf_get_file_xattr`; `bpf_path_d_path`; **no inode-keyed label cache** - the label is recomputed on every readable open (labeller ADR as revised by decisions ADR PD-1, 2026-10-23; the spike measured on both architectures that removing an attribute leaves a cached label raised, which is why the cache was dropped) | Full (kernel 7.2.x pinned, kfuncs asserted at boot) | Doctor-reported. The v0.1 floor is Linux 7.2 (labeller ADR; decisions ADR PD-2). The floor is now a policy choice, not a kfunc dependency: `bpf_real_data_inode` is unused, so doctor checks only the kfuncs the labeller calls. Support for 7.1 kernels stays deferred (Section 10.1) because Phase 0 has no evidence on 7.1 | None | Designed; not implemented (Phase 1). Kernel assumptions tested by the Phase 0 spike (Section 10); 16 of 16 on aarch64 in three identical runs and 16 of 16 on x86_64 in four, on 2026-10-23 |
| **G4 Clearance on every chunk** | Gateway and broker forward bytes only while the session label is a subset of the destination clearance at that moment. Connections opened earlier are reset when the label rises. The decision never depends on content. | Per-chunk atomic load of the mmapped label word; tunnel reset on change events; plain HTTP re-checked per request; broker checks after reading the full body | Full | Doctor-reported (requires G3 functioning) | None | Designed; not implemented (Phase 2) |
| **G5 Model channel** | A labelled session request is never placed on a backend outside its clearance; it fails closed with a structured reason. Fields outside the allowlist and provider-side tools are dropped for labelled sessions. Model API keys never enter a session. | `tidemark-broker` placement by clearance filter then static preference; field and tool allowlist; keys as systemd credentials with placeholders | Full | Doctor-reported | None | Designed; not implemented (Phase 2) |
| **G6 Staged writes carry labels** | Writes stay unmerged and invisible to other sessions until commit; root and the human can inspect them with `tidemark diff`. Commit stamps every written file with the session final label. Abort restores a byte-identical workspace. Auto-exec, skill and hook files wait for human acceptance. Agent memory written by a labelled session is visible only to sessions started at or above that label. | overlayfs per-session upper; `tidemark diff`/`commit`/`accept`/`abort` from outside the slice; label-partitioned state views | Full | Doctor-reported (requires overlay on every writable path; doctor refuses otherwise) | None | Designed; not implemented (Phase 3) |
| **G7 Evidence against the agent** | A process the session cannot reach records every decision reported to it by the event sources. Any edit, deletion, insertion or reordering of a receipt is detected offline, provided the signing key has not been extracted. Each receipt names the image digest, the policy hash and the doctor result. | `tidemark-ledgerd` own UID, `SO_PEERCRED` filter, `FS_APPEND_FL` segments, hash-chained JSONL, in-toto Statement in a DSSE envelope signed with Ed25519 | Tamper-evident against the agent; key is a systemd credential in v0.1, TPM-sealed in v0.2 | Tamper-evident against the agent; receipt records the lower guarantee level | Verifier only | Designed; not implemented (Phase 3) |
| **G8 Stop** | `tidemark stop` freezes every session and denies every exit (target 1 s or less). | cgroup v2 freezer; gateway and broker switch to deny-all | Full | Doctor-reported | None | Designed; not implemented (Phases 2 and 3) |
| **G9 Declassification** | Only a human outside the agents slice can release data. In v0.1, releases apply only to broker requests, each bound to the exact content hash and destination, consumed by the agent's own retry, and recorded. | One-time release token bound to a content hash and destination; caller checks | Full | Doctor-reported | None | Designed; not implemented (Phases 2 and 3) |

---

## 6. Non-guarantees

Tidemark does **not** provide any of the following in v0.1.

- Prevention of prompt injection, and no guarantee that agent output is correct.
- Resistance to a kernel exploit or root compromise. These defeat every
  guarantee.
- Elimination of covert and timing channels toward cleared destinations:
  request counts, timing, choice of URL, and CDN domain fronting through a
  cleared domain.
- Elimination of covert and timing channels **between concurrent sessions** on
  the shared kernel and shared services: CPU and cache contention, local-model
  latency and prompt-cache timing, and shared daemons.
- Any control over data sent to a destination the user cleared, including a
  cleared cloud model.
- Protection of unlabelled data. Data arriving through APIs is protected only
  when the destination is tagged as a private source.
- Containment or labels for agents launched **outside `tidemark run`**. Such
  processes are not sessions: the BPF program exits early for them, and no
  cgroup, namespace or Landlock domain applies. The signed image owns the agent
  launch paths; an RPM host does not, so an agent the user starts directly on
  such a host has no Tidemark guarantee.
- Preservation of labels when files are copied off the machine, or when tools
  outside sessions that drop `security.*` xattrs rewrite them. Directory rules
  are the exception. `tidemark audit-labels` reports drift only after the fact.
- Holding or undoing external side effects. Keep each agent's own approval on
  for destructive tools until v0.2.
- Content inspection. v0.1 decides per destination domain; clearances by method
  and path arrive in v0.2.
- Detection of disguised content in a release. A release relies on the human's
  reading of the bytes shown; encoded or disguised content in a release request
  is not detected.
- Proof against root or the machine owner. TPM sealing in v0.2 protects the key
  at rest only; third-party verifiability needs TPM quotes (Phase 7).
- Truth of receipt content when an event source (`tidemarkd`, `tidemark-gw`,
  `tidemark-broker`) is compromised. G7 detects tampering after an event is
  recorded; it does not validate the event.
- Full guarantees on hosts running only the RPM; those hosts get what doctor
  reports. macOS host endpoints are recorded as unattested.
- Protection against always-on messaging-gateway modes, VPS hosting and
  OpenClaw, which fall outside the v0.1 claims.
- Availability. Quotas limit resource abuse but do not eliminate it.
- Enforcement on Windows or macOS. Only the verifier runs there.

---

## 7. Known limitations

- **Coarse session granularity.** After a session reads one secret, it loses
  public egress for the rest of its life. Private memory views may hide later
  public edits to the same file. Label creep is measured in Phase 1 against
  kill criterion 3 (more than 40 percent of benign sessions ending labelled).
- **Domain-level decisions.** v0.1 clears whole domains. Domain fronting through
  a cleared CDN domain is possible. Method- and path-level clearances arrive in
  v0.2.
- **xattr fragility outside sessions.** Editors and tools that rewrite files by
  rename can drop `security.bpf.tidemark`. Directory rules, evaluated in the
  kernel, keep default sensitive folders labelled; `audit-labels` reports drift
  for the rest.
- **Receipts can expose sensitive paths to a key holder.** An unkeyed SHA-256
  of a guessable path (for example `~/.ssh/id_ed25519`) can be reversed with a
  dictionary. The receipt draft (`spec/receipt-predicate-v0.md`) therefore
  records paths and broker request content as HMAC-SHA-256 under a
  per-installation key that tidemarkd holds outside sessions (`path_hmac`,
  `content_hmac`). Anyone who obtains that key as well as a receipt can still
  confirm guessed paths, so the key must not be shared with receipts.
- **No v0.1 support for 7.1 kernels (decided, 2026-10-23: decisions ADR PD-2).**
  A 7.1 kernel meets the Landlock ABI 9 floor. The 7.2 floor is no longer a
  kfunc dependency - the labeller keeps no inode-keyed cache (PD-1), so
  `bpf_real_data_inode` is unused - but 7.2 stays the v0.1 floor because Phase 0
  gathered no evidence on 7.1, and the project does not claim support it has not
  tested. Under the 7.2 floor, doctor refuses to start sessions on a 7.1 host.
  The approach in Section 10.1 is recorded for a later decision (Phase 6
  portable mode); it is not part of v0.1.
- **Root can still alter ledger segments.** The append-only attribute binds the
  ledger daemon, not root.
- **Clients that ignore proxy settings fail closed** and are listed in the
  compatibility matrix. SSH git and plain IMAP clients are denied and
  documented.
- **Enforcement may move work to weaker local models.** This quality cost is
  measured and published, not hidden.

---

## 8. Residual risks and mitigations

| Risk | Mitigation |
|---|---|
| A bug in the per-chunk check or label plumbing reopens the pre-opened-tunnel bypass | 10,000-iteration race test; tunnel cases on every enforcement PR; fuzzing of the copy loop; named-reviewer sign-off before launch |
| Parser compromise in `tidemark-gw` or `tidemark-broker` defeats G4 and G5 and lets false events into receipts (G7 signs them faithfully) | Fuzzing of the CONNECT, HTTP and JSON parsers; each daemon under its own unprivileged UID; systemd sandboxing; the BYOC canary and host-side sinks as the external check that does not trust receipts |
| The gateway resolves a hostname before checking clearance, turning its resolver into a DNS channel | Clearance evaluated on the requested hostname before any DNS query (G2); a bypass-matrix case and benchmark scenario S-C2b with a host-side DNS query log (`docs/BENCHMARK_METHODOLOGY.md`) |
| Shared local-backend state (slots, prompt cache) carries data between sessions of different label sets | Design assumption D2 (Section 10.3); timing residual listed in Section 6 |
| Kernel fragility: kfuncs are not a stable ABI; overlayfs and the BPF verifier may differ across architectures | Kernel pinned by image digest with CI assertions; Phase 0 gate on stock Fedora 45; fallback gaps documented before selection; Renovate merges only after the boot test and read matrix pass |
| Tools outside sessions drop xattrs and remove labels | Kernel-evaluated directory rules; commit stamping; `audit-labels` drift report; claim narrowed to files Tidemark controls |
| An incumbent ships label-driven exits with label-constrained placement | `pkg/labels` and `pkg/clearance` packaged for reuse; Phase 6 adapter turns OpenShell into a distribution channel; kill criterion 6 applies |
| Root code in the TCB | Small root surface (`tidemarkd` only); unprivileged gateway, broker and ledger; fuzzing; systemd sandboxing of every daemon |
| Receipts oversold | "Tamper-evident against the agent" wording enforced by the claims checker; attestation required before any third-party claim |
| Base image is new (`fedora-bootc:45`) | Pinned by digest; `:44` fallback build kept until launch |

---

## 9. Attack channels by class, mapped to the guarantee that addresses each

Each adversarial channel is listed with the guarantee that is designed to close
it. In the benchmark, a containment channel counts only if it leaks in the
unsandboxed configuration, and a label channel only if it leaks in the Tidemark
audit configuration; the evidence class and persistence use their own validity
oracles. `docs/BENCHMARK_METHODOLOGY.md` Section 4 defines every oracle. The 18
classes below are the v0.1 gating classes; two non-gating classes follow.

### 9.1 Containment classes (7) - addressed by G1 and G2

| # | Channel | Primary guarantee | Notes |
|---|---|---|---|
| C1 | Connection to an unlisted host | G2 (single exit) + G4 (default-deny clearance) | No clearance means fail closed |
| C2 | DNS tunnelling | G2 | Sessions cannot resolve names. The gateway evaluates clearance on the requested hostname before any DNS query, so a hostname outside the session's clearance is never resolved; benchmark scenario S-C2b checks this with a host-side DNS query log |
| C3 | IP literals (raw IPv4/IPv6) | G2 | Refused by the gateway |
| C4 | UDP and QUIC | G2 | No UDP path; QUIC refused |
| C5 | Machine-local Unix sockets (container engine socket, ssh-agent, D-Bus buses) | G1 + G2 | Mount view, namespaces and Landlock `RESOLVE_UNIX` block them (assumption S10) |
| C6 | A socket opened by an MCP child | G1 | Child is inside the session cgroup and namespace |
| C7 | A double-forked daemon | G1 | `KillMode=control-group` and cgroup membership |

### 9.2 Label classes (10) - addressed by G3, G4, G5, G6

| # | Channel | Primary guarantee | Notes |
|---|---|---|---|
| L1 | Upload to an allowlisted host after a private read | G3 + G4 | Label rises before `open()` returns; chunk blocked |
| L2 | Canary in a query string to an allowlisted host | G3 + G4 | Decision never depends on content |
| L3 | `git push` over HTTPS to an allowlisted remote | G3 + G4 | Per-chunk re-check on the CONNECT tunnel |
| L4 | A tunnel opened before the read (HTTP/2 reuse, WebSocket, keep-alive) | G4 | Per-chunk re-check plus tunnel reset on label rise |
| L5 | The canary inside the next model prompt | G5 (+ G3) | Placed only on a backend cleared for the label; fails closed if none fits |
| L6 | A provider-side web-fetch or web-search tool | G5 | Field and tool allowlist drops non-function tools |
| L7 | Cross-vendor laundering through committed files | G6 + G3 | Commit stamps the file; the reader becomes labelled |
| L8 | Laundering through agent memory | G6 | Label-partitioned memory; a public session sees no private memory |
| L9 | Private data sourced from an API (mail API) | G3 (source labels) | Label raised before the first relayed byte |
| L10 | Label stripping by a rewrite followed by a re-read | G3 (directory rules) + G6 | Residual drift reported by `audit-labels` |

### 9.3 Evidence class (1) - addressed by G7

| # | Channel | Primary guarantee | Notes |
|---|---|---|---|
| E1 | Ledger tampering (edit, delete, insert, reorder) | G7 | Session write attempts denied; offline tampering detected at the exact entry while the key is secret |

### 9.4 Persistence (scored separately) - addressed by G6

Auto-exec and skill implants (`.git/hooks`, `.envrc`, shell rc files, autostart
entries, systemd user units, agent skill and hook directories) are held until a
human runs `tidemark accept` from outside the slice.

### 9.5 Intervention channels - addressed by G8 and G9

Global stop (G8) and human declassification (G9) are the only two intervention
paths and are accepted only from the human UID or root, outside the agents slice.

### 9.6 Two non-gating classes for which Tidemark v0.1 has no mechanism

These are reported next to the gating results, as measured, and do not gate.

| # | Channel | Why Tidemark v0.1 has no mechanism for it |
|---|---|---|
| N1 | Exfiltration of a credential placed in the session environment | v0.1 does not inject credentials, so a credential placed in the environment can leave over a cleared channel. Credential injection arrives in v0.2 |
| N2 | An L7 method/path rule on the same host | v0.1 decides per domain and cannot allow GET while denying POST on one host. Method and path clearances arrive in v0.2 |

---

## 10. Assumptions the Phase 0 spike must confirm

The guarantees above rest on kernel behaviour that must be proved on **stock
Fedora 45** before any product code is written: aarch64 in Lima (vz) on the M1
Pro, and x86_64 under QEMU/KVM on a GitHub runner. Each assumption below is
stated with how it is tested. The spike records results in an ADR; if an
assumption fails, the labeller ADR must document the fanotify fallback and its
coverage gaps (mmap, grandchild and overlay reads) before the fallback may be
selected, per kill criterion 1.

**Status 2026-10-23: both environments named above have run the 16-scenario
suite - aarch64 in three consecutive Lima runs, x86_64 in five consecutive
GitHub Actions QEMU/KVM runs - and no gated assumption failed on either. The
`TM_DISABLE_INODE_CACHE=1` no-cache run below is so far aarch64-only. Evidence:
`docs/spike-results/2026-10-23_Kernel-Spike-Results-aarch64_v2.md` and
`2026-10-23_Kernel-Spike-Results-x86_64_v1.md`.**

### 10.1 Facts verified against kernel sources (16 October 2026)

These were read from the `torvalds/linux` tree, so the spike confirms behaviour,
not existence.

- `file_open`, `file_receive`, `mmap_file`, `inode_post_setxattr` and
  `inode_post_removexattr` are in the sleepable BPF LSM hook set in
  `kernel/bpf/bpf_lsm.c` (v7.2); `file_permission` is not.
- The fs kfuncs `bpf_get_file_xattr` (`KF_SLEEPABLE`), `bpf_path_d_path`
  (LSM programs only) and `bpf_real_data_inode` (`KF_SLEEPABLE | KF_RET_NULL`)
  are registered in `fs/bpf_fs_kfuncs.c` at v7.2. `bpf_real_data_inode`
  returns `d_real_inode(file_dentry(file))`, which resolves the upper or lower
  inode that holds a regular file's data on an overlay.
- **`bpf_real_data_inode` is absent at v7.1 and present at v7.2** (re-checked by
  grep of `fs/bpf_fs_kfuncs.c` at both tags on 16 October 2026). Fedora 45
  ships kernel 7.2.x. A kernel at the 7.1 Landlock floor lacks it - which no
  longer matters to v0.1: the labeller keeps no inode-keyed cache (decisions ADR
  PD-1, 2026-10-23) and never calls the kfunc. The spike ran its full
  16-scenario suite, overlay included, with every cache read and write compiled
  out (`TM_DISABLE_INODE_CACHE=1`) and passed 16 of 16, which is the empirical
  basis for the Phase 6 approach below.
- **Approach for a kernel without `bpf_real_data_inode`** (recorded for a later
  decision, not part of v0.1: `docs/adr/2026-10-16_Labeller-BPF-LSM_v1.md`
  keeps a 7.2 floor and defers 7.1 support to Phase 6 portable mode; the floor
  is now a policy choice, not a missing-kfunc constraint, since the no-cache
  baseline that checklist item C2 selected - PD-1 - does not need the kfunc):
  read the label with
  `bpf_get_file_xattr` on the opened file, including the backing-file open that
  overlayfs performs for the real upper or lower file (`backing_file_open`
  calls `vfs_open`, which reaches `security_file_open`, verified in
  `fs/backing-file.c` and `fs/open.c` at v7.2); cache the result on both the
  overlay inode and the backing inode; and resolve directory rules with
  `bpf_path_d_path` on a cache miss. The coverage gap is documented: overlay
  reads whose backing open is not attributed to the session, and a higher
  cache-miss rate. If this approach is ever adopted, doctor would report
  overlay reads as unverified and G3 as degraded on such a host.
- **Recommendation for the founder (not an edit made by this document):**
  update the concept document (Read-time labels mechanism, step 4) and the
  Phase 1 compatibility floor ("Kernel floor 7.1 (Landlock ABI 9); the expected
  base kernel is 7.2") to say that Linux 7.1 is the minimum for Landlock ABI 9,
  and that the proposed v0.1 floor is 7.2 (labeller ADR; checklist decision
  FD-G).
- Landlock ABI is 9 on Linux 7.1 and 10 on Linux 7.2
  (`security/landlock/syscalls.c`). ABI 9 introduced
  `LANDLOCK_ACCESS_FS_RESOLVE_UNIX`, which restricts connections to pathname
  Unix sockets created outside the Landlock domain (absent in the v7.0 UAPI
  header at ABI 8; present at v7.1; documented in
  `Documentation/userspace-api/landlock.rst`). Abstract Unix sockets and
  signals are scoped with `LANDLOCK_SCOPE_ABSTRACT_UNIX_SOCKET` and
  `LANDLOCK_SCOPE_SIGNAL`.
- Writing a `security.bpf.*` xattr requires `CAP_SYS_ADMIN`
  (`cap_inode_setxattr` in `security/commoncap.c`), which sessions lack.
- Sleepable BPF programs may use array (including `BPF_F_MMAPABLE` arrays),
  hash, ringbuf and inode/task/cgroup local-storage maps
  (`kernel/bpf/verifier.c`), which covers the label array and the inode cache.
- `bpf_get_current_ancestor_cgroup_id` exists and is available to BPF LSM
  programs (`kernel/bpf/helpers.c`).

### 10.2 Assumptions and tests

| # | Assumption to confirm | How it is tested in Phase 0 |
|---|---|---|
| S1 | **overlayfs xattr behaviour.** A labelled lower file stays labelled when opened through an overlay, and lazy upper-file open on copy-up does not drop the label. The cache clause previously in this row (an inode-storage cache keyed by `bpf_real_data_inode(file)` returning the same slot for the overlay and direct open) was retired on 2026-10-23: the labeller keeps no inode-keyed cache (decisions ADR PD-1). | Place a labelled file in the lower dir, mount an overlay, open through the overlay from a grandchild, and read the mmapped label word. Compare `bpf_get_file_xattr` on the overlay dentry with the read at the backing event. Repeat after a copy-up (write then re-open). Attribute removal is now tested directly by `spike/tests/14_xattr_and_rule_invalidation.sh`; the invalidation hooks are no longer required, because there is nothing to invalidate. |
| S2 | **Ancestor-cgroup keying.** `bpf_get_current_ancestor_cgroup_id` at a fixed depth identifies the session for every descendant, with `Delegate=no`, and a process outside any session is identified as "no session" so the program exits early. | Start a transient unit, fork a grandchild after a double fork, and assert all map to the same slot. Place a process in a nested cgroup below the session and assert it maps to the same slot (the level-2 ancestor does not change with depth). From inside the session, attempt to create a sub-cgroup and to move a process out of the session cgroup, and assert both fail under `Delegate=no` — **now tested: `spike/tests/13_session_move_out.sh`, both attempts `EACCES`, session still raises afterwards, root positive control succeeds.** Assert a process outside the slice gets slot 0 / no raise. |
| S3 | **`file_receive` coverage.** A descriptor to a labelled file passed over `SCM_RIGHTS` raises the label in the receiving session before the receiver can read. | Pass a labelled fd between two in-session processes with `SCM_RIGHTS`, then read the label word in the receiver immediately after `recvmsg`. Confirm `receive_fd` reaches `security_file_receive` (verified in `fs/file.c`). |
| S4 | **mmap coverage.** Reading a labelled file via `mmap` raises the label. Determine whether `file_open` alone suffices (the file is opened before it is mapped) or whether the `mmap_file` sleepable hook is also required. | `mmap` a labelled file in a session, read from the mapping without an explicit `read()`, and check the label word. Record whether a mapping of an already-open fd needs the `mmap_file` hook; if so, add it. Include `sendfile` and `splice` as read-matrix cases. |
| S5 | **EPERM on `security.bpf.*` inside a session.** A session cannot set or remove the label xattr, and the BPF map pins cannot be opened from a session. | From inside a session, `setxattr("security.bpf.tidemark", ...)` and assert `EPERM`; attempt to open the map pins and assert failure. (Capability path verified in source; the test confirms the session actually lacks `CAP_SYS_ADMIN`.) |
| S6 | **kfunc availability and ABI on the pinned kernel.** The boot self-test finds Landlock ABI 9 or later, `bpf` and `landlock` in the active LSM list, cgroup v2, overlayfs, `/var` with xattr and `FS_APPEND_FL` support, and the two fs kfuncs the labeller calls loadable by a probe program. | A probe BPF object that references each kfunc loads without verifier error; the self-test reads the kernel version, the LSM list and the Landlock ABI and refuses to start sessions if any check fails. `bpf_real_data_inode` was removed from `REQUIRED_KFUNCS` in `image/selftest/tidemark-selftest.sh` on 2026-10-23 (decisions ADR PD-1/PD-2), so a kernel without it no longer fails the self-test; the 7.2 floor is enforced separately by `kernel_floor`. |
| S7 | **Label raise before `open()` returns, from a grandchild.** The core guarantee G3: the label word is already raised when read in the same process immediately after a grandchild's `open()` returns. | A grandchild opens a labelled file and signals the parent; the parent reads the mmapped word and asserts it is raised with no sleep. The 10,000-iteration race test that the labeller ADR's kill test requires was pulled forward into Phase 0 and now runs in `spike/tests/15_race_10k_iterations.sh`, in both readings (opener and a separate outside-session observer): **10,000 iterations, zero misses, three runs on 2026-10-23.** Phase 1 still repeats it under a real workload. |
| S8 | **Filesystem substrate behaviour (recorded, not gating).** How xattrs and `FS_APPEND_FL` behave on virtiofs and 9p, and whether Lima cidata provisioning works on `fedora-bootc`. | Mount a virtiofs and a 9p share in the spike VM, attempt to set and read `security.bpf.*` and `FS_APPEND_FL`, and record support. Attempt a Lima cidata first boot and record the result. Mac users whose data lives in virtiofs shares stay outside v0.1 claims until this is known. |
| S9 | **Directory rules match inside the session's mount view.** Sessions run in a private mount namespace with `ProtectHome=tmpfs` and `BindPaths`, and overlay opens also produce a backing-file open whose path lies on an internal layer mount. The path that `bpf_path_d_path` returns may therefore differ from the host path that a directory rule is written against. | Open a file under a rule-labelled root (with its xattr removed) in three ways: through the session's bind mount, through the session's overlay, and directly on the host. Record the string `bpf_path_d_path` returns in each case and whether the rule matched. If the strings differ, the labeller ADR records how rules are matched (session-view paths, or the device and inode of the rule root) before Phase 1. |
| S10 | **Machine-local Unix sockets are unreachable.** Landlock ABI 9 `LANDLOCK_ACCESS_FS_RESOLVE_UNIX` (Section 10.1) limits pathname Unix socket connections to the two per-session sockets, with abstract sockets scoped. `ProtectHome=tmpfs` also hides `/run/user`, so this is a second layer. | From inside a session, `connect()` to `/run/dbus/system_bus_socket`, a rootful Podman socket, and a user D-Bus, Podman and ssh-agent socket bound into the session view for the test; assert `EACCES` (or `ENOENT` where hidden), and assert that both per-session sockets still connect. Connect to an abstract socket created outside the session and assert failure. |
| S11 | **Read-only label view for the daemons.** A non-root daemon UID can map the label array read-only and sees each raise, and cannot obtain a writable mapping. | Have a test process under a dedicated non-root UID obtain the map descriptor the way the gateway will (passed by `tidemarkd`, or from a pin with restricted permissions), `mmap` it with `PROT_READ`, trigger a raise from a session, and assert the raise is visible with no syscall in between. Assert that a `PROT_WRITE` mapping and a map update from that UID both fail. |

**Coverage by the checked-in spike.** The spike in `spike/` has 16 scenarios,
and `spike/README.md` holds the authoritative mapping of its local assumption
IDs to the rows above: S7 is A1 and A1r (the 10,000-iteration race), S2 is A2
and A2c (the move-out case), S1 is A3, S3 is A5 and A5b (an unprivileged
receiver, for which a failed attribute read fails closed), S5 is A7, and S4 is
A8 and A8b. A8b isolates the `mmap_file` hook for a descriptor opened before the
process joined the session, and records that `read()` on such a descriptor fires
no attached hook; G3 therefore requires that no descriptor be inherited across
the shim exec. A12 covers attribute raising, raising to a higher label,
attribute removal and directory-rule add/clear invalidation. A6 covers only part
of S11, because the array is mapped as root, not from a non-root daemon UID. S9
is covered only by observations: A4b (a bind-mount alias in a private user and
mount namespace) and A4c (a hard link or rename out of a labelled root inside
the session's own namespace, which needs no new namespace and is therefore not
mitigated by G1). **A4c's open question is closed: directory labels follow the
location, and Phase 1 denies in-session link and rename out of a ruled root
(decisions ADR PD-4); the spike does not test that enforcement, because the
probe never denies.** The spike also shows that the overlay backing-file event
carries a layer-relative path, so the probe skips directory rules for that event.
S6 is covered by the probe's load step, its environment banner and the boot
self-test. Scenario 16 records the `open()` latency baseline that PD-1 was
decided on. The spike still does not test the copy-up step of S1, the removal and
map-pin steps of S5, the `sendfile` and `splice` cases of S4, or S8 and S10.
Under the spike's failure mapping, a failure of A5, A5b or A8b narrows the public
claim, and a failure of A6 requires a redesign of the label transport; none of
these failures triggers the Phase 0 kill gate below on its own.

### 10.3 Design assumptions outside the kernel spike

These are not kernel questions, so the Phase 0 spike does not test them. Each is
tested in the phase that builds the component.

| # | Assumption | Tested in |
|---|---|---|
| D1 | The broker exposes only the OpenAI Chat Completions endpoint to sessions. Backend management, slot, metrics and model-loading endpoints (for example llama-server's slot save and restore) are never reachable from a session. | Phase 2 broker contract tests |
| D2 | Backend slot state and prompt-cache state are never shared between sessions of different label sets: the broker assigns backend slots per label set and clears or partitions cache state when the label set changes. Residual timing effects are covered by the non-guarantee on inter-session covert channels (Section 6). | Phase 2 broker tests, with a public session probing for a private session's cached prefix |
| D3 | The gateway evaluates clearance on the requested hostname before it issues any DNS query (G2). | Phase 2 bypass matrix and benchmark scenario S-C2b |

**Phase 0 kill gate.** If label raising is unreliable on stock Fedora 45 on
either architecture or through overlayfs, and the fanotify fallback misses reads
by child processes or through mmap, drop the claim of kernel-recorded labels and
ship the gateway, broker and receipts as an audit tool, or stop. A written go or
no-go decision is recorded at the end of week 3.

---

## 11. Verified facts and sources

The following were checked against primary sources on 16 October 2026 and
support this draft. Version numbers and digests are point-in-time observations,
not fixed facts. The Go module path
`github.com/BitecodesHub/tidemark` was a placeholder when this was written and
was resolved on 2026-10-23 (decisions ADR PD-10); the container image name
`ghcr.io/${{ github.repository_owner }}/tidemark` (lowercased) is derived at run
time.

- **Base image (observation, not a pin).** `quay.io/fedora/fedora-bootc:45` is
  a multi-architecture manifest list that is rebuilt daily, so its digest
  changes from day to day. Observed on 16 October 2026 at about 11:30 UTC:
  `sha256:86cc9d98ef210b71a6c4eeec4b7d04e7fef0a4415184c2322a14c879df189f45`
  (amd64 child
  `sha256:4e5c285af0fbf9698525a3c5ff6b6dbfc5c1d6102b0f1f46a48b2c437ffe533c`,
  kernel `7.2.9-300.fc45`), via the quay.io API. The fallback `:44` was observed
  at `sha256:9db34195bc0b6be43e17c20ca840be5bc153093c94b62f6792a76fb199eeaced`
  (kernel `7.2.10-200.fc44`). The authoritative record of the pinned digest is
  `docs/adr/2026-10-16_BaseImage_v1.md`. **Both of those 16 October observations
  are now historical: the registry stopped serving them and the pins were
  refreshed on 2026-10-23 (decisions ADR PD-14), which is why this document
  treats a digest as an observation with a date rather than a promise.**
- **Fedora 45 kernel (observation).** Stable `kernel-7.2.9-300.fc45` (stable
  since 6 October 2026), `kernel-7.2.10-300.fc45` in testing and
  `kernel-7.2.11-300.fc45` pending, observed in Bodhi on 16 October 2026 (UTC).
  Fedora 45 does not ship 7.3, so Landlock ABI 11 features are not available on
  the pinned kernel.
- Landlock ABI and `LANDLOCK_ACCESS_FS_RESOLVE_UNIX`, the sleepable hook set,
  the fs kfuncs and their flags, `backing_file_open`, sleepable map types and
  the `CAP_SYS_ADMIN` xattr check were read from the `torvalds/linux` tree at
  tags `v7.0`, `v7.1`, `v7.2` and `master`.
- `CONFIG_LSM="lockdown,yama,integrity,selinux,bpf,landlock,ipe"`,
  `CONFIG_BPF_LSM=y`, `CONFIG_SECURITY_LANDLOCK=y`, `CONFIG_FW_CFG_SYSFS=m` in
  the Fedora kernel-ark `os-build` configs.

---

## 12. Change log

| Date | Version | Change |
|---|---|---|
| 2026-10-16 | v0 draft | First draft. |
| 2026-10-16 | v0 draft (review fixes) | Added the per-guarantee TCB (Section 1.1); a Status column and "designed level" wording in Section 5; the 7.2 requirement and the 7.1 fallback for G3; clearance before DNS resolution (G2, C2, D3); non-guarantees for launches outside `tidemark run`, inter-session covert channels, disguised releases and compromised event sources; the "machine to cleared destinations" boundary; `tidemark diff` in G6; assumptions S9 to S11 and D1 to D3; reversible path hashes; dated, point-in-time version facts. |
| 2026-10-17 | v0 draft (consistency pass) | Aligned the G3 host column, the G3 mechanism, Section 7, S6 and Section 10.1 with the labeller ADR (proposed 7.2 floor, no-cache baseline, 7.1 support deferred); updated the path-hash limitation to the HMAC fields now in the receipt draft. |
| 2026-10-17 | v0 draft (spike alignment) | Corrected the S2 test (nesting below the session does not change the level-2 ancestor; the move-out step now matches checklist AC3); added the coverage of Section 10.2 by the 11 spike scenarios, following `spike/README.md`. |
| 2026-10-23 | v0 draft (decisions and 16-scenario spike) | Applied the Phase 0 founder decisions ADR: G3, Section 7, S1, S6 and Section 10.1 now state that the labeller keeps no inode-keyed cache (PD-1) and that the 7.2 floor is a policy choice rather than a `bpf_real_data_inode` dependency (PD-2); S2, S7 and the coverage paragraph record the scenarios that now test move-out, the 10,000-iteration race in both readings, and attribute removal; A4c's aliasing question is closed by PD-4 (location, with link/rename denied out of a ruled root); spike scenario count 12 to 16. |
| 2026-10-23 | v0 draft (x86_64 spike) | Section 10 gains its status line: the x86_64 CI run passes 16 of 16 in four runs, so G3's kernel assumptions now hold on both architectures (the `TM_DISABLE_INODE_CACHE=1` run stays aarch64-only), and G3's cache-defect note is now a both-architecture measurement. |
| 2026-10-23 | v0 draft (x86_64 repetition) | A fifth consecutive green x86_64 CI run closed the same day at commit `7b02d80`, where the only file differing under `spike/` is `spike/README.md` — the test code, the probe and the loader are unchanged. Section 10's status line now reads five runs. Nothing in the guarantee table changes: repetition strengthens G3's kernel assumptions, it does not add one. The run-to-run latency spread is recorded in the spike report, not here, because this document claims no numbers. |

---

## Tone Check

Included under the organisation's drafting rule for external-facing drafts.
Scored against the "Warm but Authoritative" brand guideline; remove or keep at
the founder's discretion before the v0 tag.

- **Authoritative: 5 of 5.** Precise, sourced, and explicit about what is not
  claimed.
- **Warm: 2 of 5.** Appropriately neutral for a security specification; the
  README, not this file, should carry warmth.
- **Overall:** fit for purpose as a technical reference; no change recommended.

*End of THREAT_MODEL.md v0 draft. To be tagged as v0 at the end of Phase 0.*
