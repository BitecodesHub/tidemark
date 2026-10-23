# AgentOS Tidemark (working name; trademark search required before launch)

> What an agent session has read decides where it may send data, including which model sees it. Check it yourself with a canary you generate on your own machine.

> Produced by the concept-design workflow on 2026-10-16: five independent architects, three judges, a synthesis,
> three red-team reviewers and a final revision. The approved project plan summarizes this document.
> Note: a naming mismatch split the judges' totals; normalized, the security-first 'Escrow' concept scored 122 and
> Tidemark 116. Tidemark absorbs Escrow's core within a scope one part-time maintainer can deliver.

## Executive summary
AgentOS Tidemark is a headless, signed bootc image built FROM quay.io/fedora/fedora-bootc:45, pinned by digest, with a :44 fallback build kept until launch. The software ships as a Fedora RPM (COPR) that the image installs. Fedora Atomic users can therefore also adopt it with `bootc switch`. Full guarantees are claimed only on the signed image.

v0.1 is a confidentiality preview for unmodified Hermes Agent and OpenCode, run with their own read and network prompts turned off. It integrates existing containment and adds one new kernel piece:
(1) Read-time labels. A sleepable BPF LSM program raises the session's label before a labelled open returns. Labels come from file xattrs and from directory rules, and the session is identified by its ancestor cgroup.
(2) Exits re-checked on every chunk. The gateway and the model broker forward bytes only while the session's current label lies inside the destination's clearance. This closes tunnels opened before a read.
(3) The model call is an exit. A labelled session is placed only on backends cleared for its label (by default the local llama-server). If none fits, the request fails closed with a reason.
(4) Labels move forward. Every writable path is staged. Commit stamps written files with the session's final label, and agent memory is partitioned by label.
(5) Receipts are signed by a daemon the agent cannot reach. They are tamper-evident against the agent, not proof against the machine's owner.

The public proof is 'Bring Your Own Canary': a stranger generates the secret on their own host and greps logs kept on that host, outside the VM. v0.1 does not hold side effects, so the release notes advise keeping each agent's own approval on for destructive tools until v0.2.

The schedule has been re-baselined:
- Phase 0 (3 weeks) ends in a combined kernel and interview kill gate.
- The MVP runs from week 4 to week 24: about 210 hours, 147 of them building. Default cuts are already applied. Launch target is early April 2027, mid-May at the latest.
- After launch, each quarter has one theme, chosen from user evidence.
- Training, multi-engine serving and hardware adapters stay deferred, each with an explicit trigger.

## Revolutionary claim
The claim, as a v0.1 preview: Tidemark makes what a session has read decide where that session's data may go next. It covers unmodified agents from different vendors: Hermes Agent and OpenCode in v0.1, and Claude Code from Phase 6. Every MCP server, hook and subprocess those agents start is covered too.
- **The kernel records each labelled read before the read returns.** This covers the whole process tree, and the agent cannot lower the label.
- **Every exit re-checks the session's current label before it forwards each chunk.** A tunnel opened before the read is cut. Encoding, compression, splitting or paraphrase make no difference, because the decision never depends on content.
- **The model call is an exit.** A labelled session's requests go only to backends cleared for its label. Otherwise they fail closed with a structured reason, and nothing silently falls back to the cloud.
- **What a labelled session writes carries its label forward.** Committed files are stamped with the label. Agent memory written in a private session is visible only to sessions started as private.
- **Private data read through declared APIs raises the label**, exactly as a file read does.

How to check it:
1. Generate a random canary on your own host.
2. Boot the cosign-verified image with the canary injected.
3. Run the same scripted, fully compromised agent in audit mode and then in enforce mode.
4. Grep the sink logs on your host, outside the VM. Enforce mode must show zero matches on every in-scope channel.

Receipts are a secondary check. They are tamper-evident against the agent and can be verified offline as long as the signing key has not been extracted. They are not proof against the machine's owner or root until TPM attestation arrives (Phase 7).

What is not new:
- The high-water mark dates to ADEPT-50 (1969).
- Shipped operating systems have enforced static confidentiality labels on network exits for decades: SELinux MLS/MCS and Solaris Trusted Extensions.
- HiStar, Flume and Asbestos did OS-level information-flow control as research systems.
- ActPlane already propagates kernel labels along file read and write edges and can block connect() through BPF-LSM.
- CaMeL and FIDES enforce information-flow control inside one rewritten agent. Invariant Guardrails enforces data-flow rules at the tool layer.
- Containment, credential proxies, staging, rollback, audit chains and kill switches already ship in OpenShell, sandbox-runtime, nono, AgentFS, agentsh and Microsoft AGT. Tidemark integrates these ideas and claims none of them.

What is new is a narrow combination, as far as public READMEs show on 16 Oct 2026:
- floating labels on the process trees of unmodified agents from different vendors;
- the model API treated as a labelled exit, with placement constrained by the label;
- per-chunk re-checks at every exit;
- label-carrying writes and label-partitioned agent memory;
- all of this shipped as a reproducible image with a demo strangers can check.

Not claimed:
- prevention of prompt injection;
- resistance to kernel or root compromise;
- protection against covert or timing channels to cleared destinations;
- protection against misuse by destinations the user cleared;
- protection of unlabelled data;
- content inspection (v0.1 decides per destination domain);
- holding or undoing external side effects in v0.1;
- third-party proof against the machine's owner.

## First users
**Primary users (narrowed).** Homelab and Proxmox users, and people running a dedicated Linux VM, who run coding or research agents (Hermes Agent, OpenCode) over local files: code, notes and secrets. They mix a local llama.cpp model with cheap cloud models.

**Why this group, from the evidence:**
- Hermes Agent has 254k stars and installs hardware-matched llama.cpp builds.
- Local users already mix in cloud models.
- Per-action approval has already failed: users approve about 93% of prompts, and Omarchy launches 11 of 15 agents with prompts disabled.
- This group already runs VM images.
- The breakout pattern of llama.cpp, Ollama and Open WebUI started with this audience and one checkable demo.

**Excluded from v0.1 claims:**
- Hosts such as a USD 5 VPS. Hermes Agent's own README suggests running it on one, but such hosts cannot run llama-server and rarely accept custom images.
- Always-on messaging-gateway modes (Telegram, Discord, Slack, WhatsApp, Signal). Clearing a chat domain at domain level would also permit sending to any chat.
- OpenClaw, whose profile arrives in Phase 6.
- Mac users whose data lives in virtiofs shared folders, until Phase 0 confirms how xattrs behave there.

**Secondary users, as adversarial co-developers.** Agent-security researchers reach the project through the standalone benchmark container and the cross-platform verifier, without booting the image. Each bypass they report becomes a regression scenario.

**Validation.** In weeks 1 and 2, hold 8 to 10 observational interviews (not leading questions). Each participant shows:
- their agent configuration;
- where their private data actually lives (local files or APIs);
- which destinations they would clear;
- where they host the agent.
The kill decision falls at week 3.

**Deliberately not first:**
- enterprise CISOs, because Agent 365 and AgentCore are bundled;
- regulated deployers, because evidence mapping is only a backlog item;
- GPU and training users, because their work cannot be tested on the M1.

Willingness to pay is unvalidated and not needed for v0.1.

## Positioning

### NVIDIA OpenShell (Apache-2.0, v0.1.4)
- Relationship: Benchmark baseline in v0.2. In Phase 6 it becomes a distribution channel through a supervisor-middleware adapter built on Tidemark's pkg/labels and pkg/clearance libraries.
- Differentiation: Its docs (checked 16 Oct 2026) list the following:
- Landlock, seccomp and network-namespace sandboxing;
- default-deny egress with request-level read/write rules;
- credential placeholders;
- a policy prover;
- OCSF events;
- supervisor middleware: a gRPC service over TCP that inspects, transforms or denies allowed HTTP and WebSocket traffic after policy evaluation and before credential injection.
Its managed inference router has been removed. It has no read-time labels and no label-driven model placement.

Tidemark v0.1 is weaker on L7 rules and credential placeholders, and says so. It differs only on labels, placement, label-carrying writes and evidence. NVIDIA could add a taint check as middleware in one release, so Tidemark treats that middleware as a delivery route rather than only a threat.

### NVIDIA NemoClaw (alpha)
- Relationship: Main competitor in the self-hoster segment.
- Differentiation: It runs OpenClaw and Hermes Agent inside OpenShell, with local Ollama or vLLM or with hosted providers. Local models need no NVIDIA account, and it runs on Linux with rootless Podman and on Windows through WSL. Account-free and hardware-neutral therefore do not set Tidemark apart. Its Model Router documents no routing criteria tied to what the session has read. Tidemark's placement is decided by labels.

### nono (Apache-2.0, about 4.5k stars)
- Relationship: Benchmark baseline in v0.2; a v0.1 stretch goal if setting it up takes 4 hours or less. Prior art for rollback, integrity logging and signed instruction files.
- Differentiation: README v0.61.2 lists:
- Landlock and Seatbelt sandboxing;
- a credential proxy that keeps keys outside the sandbox;
- host and endpoint allowlists enforced through a local proxy;
- content-addressable rollback with Merkle integrity;
- append-only integrity hashing;
- Sigstore-signed instruction files.
The README describes no labels, taint or information-flow control. Tidemark adds read-time labels that decide both exits and placement.

### ActPlane (eunomia-bpf, MIT, about 104 stars, no releases)
- Relationship: The closest prior art for the mechanism. Tidemark studies it and offers contributions upstream, but does not depend on it.
- Differentiation: Corrected. ActPlane propagates kernel labels along fork/exec and file read/write edges. It hooks open, unlink, rename and connect, and can block through BPF-LSM, for Claude Code, Codex or any command.

Tidemark's differences are narrower than the earlier plan stated:
- the model API as a labelled exit, with placement;
- per-chunk re-checks on proxied connections;
- staged writes that carry labels, and agent memory partitioned by label;
- a demo strangers can check from their own host;
- an owned image.
ActPlane's README does not describe these (checked 16 Oct 2026).

### Anthropic sandbox-runtime (Apache-2.0, about 5.6k stars)
- Relationship: Benchmark baseline (b), in its strongest documented configuration. Phase 6 adds a label-aware filterRequest example for it.
- Differentiation: Corrected. srt sandboxes the entire process tree, local MCP servers included. It has HTTP and SOCKS proxies over Unix sockets, a domain allowlist, credential masking, a filterRequest hook and TLS termination. Only Claude Code's built-in sandbox is limited to shell commands. srt has no read-time labels and no model placement. Tidemark's per-session socket design follows srt's pattern.

### Invariant Guardrails (Apache-2.0)
- Relationship: Complementary. It is the main alternative approach to information-flow control for agents.
- Differentiation: It enforces data-flow rules over sequences of tool calls (for example, reading an inbox and then sending an email) through an MCP or LLM proxy, so it sees tool names and arguments. Tidemark sees file reads and network exits in the kernel and the gateway, whatever the tool is called, but it cannot see what a tool means. The two can be stacked.

### Microsoft Agent Governance Toolkit (MIT, public preview)
- Relationship: Prior art for the audit log, the kill switch and evidence mapping.
- Differentiation: It ships a Merkle audit log, EU AI Act mapping, a kill switch, execution rings and OPA or Cedar policy. Tidemark's receipts and global stop are therefore integration, not novelty. The evidence only Tidemark produces is per-session label raises and placement decisions, bound to an image digest.

### agentsh (Apache-2.0)
- Relationship: Prior art for the broker, approvals and holds.
- Differentiation: It ships an LLM proxy with DLP redaction, approvals by TOTP or WebAuthn, soft_delete and redirect decisions, checkpoints and audit-integrity keys. Its decisions depend on content. Tidemark decides on what the session has read, and does no content inspection in v0.1.

### Staging and rollback prior art: AgentFS, Dagger container-use, Rubrik Agent Rewind, Fly.io Sprites, Claude Code checkpoints
- Relationship: Prior art. Tidemark's staging is integration over overlayfs.
- Differentiation: These products offer copy-on-write isolation, discardable branches, snapshots and rollback. Rubrik has offered rollback since August 2025, and Claude Code checkpoints cover only edits made by its file tools. Staging's only distinct role in Tidemark is to stamp committed files and agent memory with the session's final label.

### Qubes OS
- Relationship: Prior art for an owned OS with controlled egress.
- Differentiation: Qubes isolates user-chosen domains in separate VMs with per-domain firewalls. That isolation is stronger than a shared kernel. Its domains are static: the user decides in advance where each piece of data lives. Tidemark raises a session's label automatically as the session reads, inside one kernel, and constrains the model call.

### SELinux MLS/MCS and Solaris Trusted Extensions
- Relationship: Prior art. Tidemark coexists with Fedora's SELinux targeted policy.
- Differentiation: They have shipped mandatory confidentiality labels with labelled networking for decades. A process's level is fixed when it starts and does not rise as it reads. That is why Tidemark uses a small BPF high-water mark instead of an MLS policy.

### Microsoft Execution Containers (MXC, MIT)
- Relationship: Containment on other operating systems. A policy-export candidate in the backlog.
- Differentiation: Process containers under one JSON policy, without labels or placement. The repository URL could not be re-verified on 16 Oct 2026.

### Docker Sandboxes
- Relationship: A complementary developer-desktop product. Possible community-contributed baseline.
- Differentiation: It offers microVMs, network policies, a credential proxy and a paid governance tier. It documents no labels and no signed receipts.

### Hermes Agent, OpenCode, Claude Code, OpenClaw
- Relationship: Workloads wrapped unmodified, not competitors. Hermes Agent and OpenCode are supported in v0.1. Claude Code, and OpenClaw in short sessions, arrive in Phase 6.
- Differentiation: Tidemark treats all of them as untrusted. Hermes Agent already ships seven terminal backends and an option to approve memory writes, so Nous could add label-like gating upstream. Tidemark will offer Hermes a 'tidemark' backend and its label specification.

### Omarchy
- Relationship: No longer a named benchmark baseline. A neutral 'unsandboxed' configuration replaces the earlier 'Omarchy launcher flags' baseline.
- Differentiation: Omarchy holds the 'agentic distro' mindshare, with a desktop and agents launched with prompts disabled. Tidemark does not compete for desktop polish or that audience.

### Microsoft Agent 365, AWS Bedrock AgentCore, Google Agent Identity, Okta Cross App Access
- Relationship: Not competitors. Exporting evidence to them is a backlog item with a trigger.
- Differentiation: They govern identity, inventory and cloud policy inside their own estates. Homelab and on-premises Linux agent VMs fall outside those estates.

### Ollama, and Ollama with Fly.io Sprites
- Relationship: Ollama is a supported broker endpoint and is not bundled. Sprites is an adjacent hosted persistent-VM product.
- Differentiation: Ollama serves models but does not govern the agent's process tree or egress. Sprites restores files and VM state and keeps credentials out, but has no labels and cannot undo external effects.

### CaMeL and FIDES (research)
- Relationship: Complementary.
- Differentiation: They track information flow per token inside a rewritten agent, at a utility cost of 3 to 32 points. Tidemark is coarser, working per session, but needs no change to the agent and adds no token overhead.

### Fedora bootc, Universal Blue, Bluefin
- Relationship: Base image and peers. Tidemark's OCI image can be adopted with `bootc switch`.
- Differentiation: Tidemark is one Containerfile plus an RPM on the multi-architecture fedora-bootc base, not a new distribution. Universal Blue images are amd64-only.

## Primitives

### Session principal (integration)
Turns each agent run, with every descendant (MCP servers, hooks, LSP servers, shells, double-forked daemons), into one contained principal. The principal has its own cgroup subtree, a private mount view, a network namespace with loopback only, quotas and a receipt identity. v0.1 allows one active session per agent profile.

**Mechanism.** tidemarkd starts a systemd transient unit through D-Bus StartTransientUnit in the agents slice, with these properties:
- PrivateNetwork, PrivateTmp, PrivateIPC, ProtectSystem=strict;
- ProtectHome=tmpfs, with BindPaths only for the staged overlays and the read paths the profile declares;
- NoNewPrivileges;
- RestrictNamespaces=yes, which also denies user namespaces;
- Delegate=no;
- a SystemCallFilter that also denies mount, bpf, keyctl, ptrace, process_vm_readv/writev, pidfd_getfd and io_uring_setup;
- CPUQuota, MemoryMax, TasksMax and KillMode=control-group.

Each agent profile runs under one fixed UID, so persistent state stays writable. Sessions are isolated from each other by cgroup, mount namespace and Landlock domain, not by UID.

tidemark-shim then:
1. calls close_range on every inherited descriptor above stderr;
2. applies go-landlock in strict mode at ABI 9 or later, with abstract Unix sockets and signals scoped and pathname Unix sockets limited to the two per-session sockets;
3. waits until tidemarkd confirms the session's label slot exists, and fails closed if it does not;
4. starts loopback forwarders. These are untrusted, because enforcement happens at the gateway;
5. execs the agent.

**Reuse.** systemd transient units and sandboxing, Landlock, cgroup v2, namespaces, go-landlock (MIT), coreos/go-systemd, godbus. The Unix-socket proxy bridge follows Anthropic sandbox-runtime's pattern.

**Build.** About 600 lines of Go.

### Read-time labels (floating high-water mark; the one new kernel piece)
When any session process opens a labelled file, or receives a descriptor to one, the session's label rises before the call returns. The label never falls during the session. A file is labelled by its own xattr or by a directory rule. Directory rules keep default sensitive folders labelled even when an editor or other tool drops xattrs.

**Mechanism.** Sleepable BPF LSM programs on file_open and file_receive. Both hooks were verified as sleepable in kernel/bpf/bpf_lsm.c; file_permission is not.
1. Identify the session by its ancestor cgroup id at a fixed depth (bpf_get_current_ancestor_cgroup_id). Exit at once if the process is not in a session.
2. Read security.bpf.tidemark with bpf_get_file_xattr.
3. If no xattr is found, resolve the path with the bpf_path_d_path kfunc and prefix-match it against the labelled roots.
4. Cache the result in BPF inode storage keyed by bpf_real_data_inode(file), so a labelled lower file stays labelled when opened through overlayfs. inode_post_setxattr and inode_post_removexattr invalidate the cache. All three kfuncs were verified in fs/bpf_fs_kfuncs.c.
5. Atomically OR the label bits into a BPF_F_MMAPABLE array indexed by session slot. tidemarkd also emits a ring-buffer event, but only when a label actually changes.

The gateway and broker hold read-only mappings of that array. Setting security.bpf.* needs CAP_SYS_ADMIN, which sessions lack.

Source labels: the gateway raises a session's label through tidemarkd's raise-only call before it relays the first byte from a private-source destination.

Fallback, chosen only if the Phase 0 gate fails: a fanotify FAN_OPEN_PERM vault. Its gaps for mmap, grandchild and overlay reads must be documented before it can be selected.

**Reuse.** Mainline BPF LSM and the fs kfuncs (bpf_get_file_xattr, bpf_path_d_path, bpf_real_data_inode), BPF inode storage, cilium/ebpf with bpf2go (MIT). Prior art: ADEPT-50's high-water mark and ActPlane's kernel label propagation.

**Build.** About 250 lines of BPF C and 300 lines of Go.

### Clearance-gated exits with per-chunk re-check
One rule for every exit: bytes leave only while the session's current labels are a subset of the destination's clearance. The check runs before every forwarded chunk, not once per connection. A tunnel, keep-alive connection, HTTP/2 connection or WebSocket opened before a labelled read therefore cannot carry data once the label has risen. Denials happen without prompts.

**Mechanism.** tidemark-gw is unprivileged and runs under its own UID. It owns one Unix socket per session, and the socket identifies the session.

What it does:
- implements HTTP CONNECT and plain HTTP forwarding;
- resolves DNS itself and refuses IP literals and non-HTTP protocols unless policy names them.

How the check works:
- Before each write of client bytes upstream, it atomically loads the session's label word and compares it with the destination's clearance. The client's write to the socket orders this load after the label raise, because the raise completes before open() returns.
- When the label word changes, it also resets every tunnel the label no longer covers.
- Plain HTTP is re-checked per request after the body is read.
- Refusals return HTTP 403 with JSON giving the reason, the label and the destination.

Policy is YAML. Each destination has a clearance and an optional source label.

**Reuse.** Go net/http. Credits: sandbox-runtime's proxy-over-Unix-socket pattern and OpenShell's egress design. The pkg/clearance evaluator is packaged so an OpenShell middleware adapter can reuse it.

**Build.** About 600 lines of Go.

### Label-constrained model broker
Treats the model call as an exit. Each session gets one local OpenAI Chat Completions endpoint. Each request is placed on a backend whose clearance covers the session label. If no cleared backend fits, the request fails closed with a structured reason, and nothing falls back to the cloud silently. The broker is the only holder of model API keys.

**Mechanism.** tidemark-broker is unprivileged, runs under its own UID and listens on the per-session socket.
- It reads the full request body (size-capped) before it checks the label.
- Placement filters backends by clearance, then applies a static preference order.
- For labelled sessions, an allowlist parser drops unknown fields and any tool not of type function. This closes provider-side web fetch, web search and remote connectors.
- Streaming responses pass through unchanged.
- Keys arrive as systemd credentials readable only by the broker. Sessions see placeholders, which do not taint.
- A release token, bound to a content hash and destination, lets the agent's own retry through once.

Default clearances:
- local llama-server: all labels;
- host endpoints (macOS Metal or Ollama): unattested local;
- cloud providers: public only, and the user may add private.

**Reuse.** llama.cpp llama-server (MIT) as the one validated backend; Podman Quadlet; systemd-creds. Prior art: agentsh's LLM proxy and generic OpenAI-compatible proxies.

**Build.** About 700 lines of Go. Responses and Anthropic Messages arrive in Phase 6.

### Staged writes with label stamping, and label-partitioned agent state
Every writable path is staged, so a labelled session's writes carry its label forward. This holds even for files opened for writing before the label rose. Agent memory cannot launder data between sessions. Abort restores the workspace byte for byte.

**Mechanism.** Each writable path is an overlayfs mount with a per-session upper directory under /var/lib/tidemark/stage. tidemarkd creates it and binds it into the session. Doctor refuses any profile that declares a writable path without an overlay.

Workspaces (manual commit by default):
1. After the session stops, `tidemark commit` stamps every file in the upper directory with the session's final label.
2. It holds any file whose lower copy changed since the session started.
3. It holds auto-exec paths (.git/hooks, .envrc, editor task files, shell rc files, autostart entries, systemd user units, agent skill and hook directories) until someone runs `tidemark accept`.
4. It then merges.
Commit and accept are refused from inside the agents slice. Abort discards the upper directory.

Agent state (for example Hermes memories and skills): the view for a label set L stacks the deltas of every label set contained in L. At session end, the state diff merges into the delta for the session's final label set. A session started without --label therefore never sees memory written by private sessions. A session started with --label private sees it, and is private from its first instruction.

**Reuse.** overlayfs. Prior art: container-use, AgentFS, nono rollback, Rubrik Agent Rewind and Claude Code checkpoints.

**Build.** About 700 lines of Go. Contingent cut: discard state writes from labelled sessions, with a notice.

### Agent-unreachable receipts (integration)
Records every session start, label raise (path hash and cause), exit decision, tunnel reset, placement, release, commit and abort. The process doing the recording is one the session cannot reach. Receipts are tamper-evident against session processes and can be verified offline as long as the key has not been extracted. They are not proof against root or the owner.

**Mechanism.** tidemark-ledgerd runs under its own UID and accepts events only from the peer UIDs of tidemarkd, the gateway and the broker, checked with SO_PEERCRED.
- Ledger segments carry the kernel append-only attribute (FS_APPEND_FL), set by tidemarkd. The ledger daemon can append but cannot truncate or rewrite them. Root still can.
- Events form a hash-chained JSONL log.
- Each session ends with one in-toto Statement in a DSSE envelope, signed with Ed25519. The key is a systemd credential in v0.1 and is sealed with TPM2 in v0.2.
- tidemark-verify ships as static binaries for macOS, Linux and Windows. It checks the signature, the chain, the image digest and the clearance invariants.

**Reuse.** in-toto, DSSE through go-securesystemslib, crypto/ed25519, systemd-creds. Prior art: Microsoft AGT's Merkle audit and nono's integrity hashing.

**Build.** About 400 lines of Go, including the verifier.

### Release and global stop (integration)
The human has exactly two ways to intervene. They can release one exact model request after seeing its bytes, or they can stop everything.

**Mechanism.** In v0.1, a release applies only to broker requests, because only there are the exact bytes visible.
- `tidemark queue show <id>` prints the request body and its destination.
- `tidemark release <id>` mints a one-time token bound to the content hash and destination. The agent's own retry consumes it, and Tidemark never replays anything.
- Grants for CONNECT destinations wait for TLS termination in v0.2.

`tidemark stop` writes cgroup.freeze=1 on the agents slice and switches the gateway and broker to deny-all.

Release, stop, commit and accept are accepted only from the human user's UID or root, and only from outside the agents slice.

**Reuse.** The cgroup v2 freezer. Prior art: the Microsoft AGT kill switch and agentsh approvals.

**Build.** About 150 lines of Go.

## Architecture layers

### L0 Image (owned OS) and delivery
- Containerfile FROM quay.io/fedora/fedora-bootc:45@sha256 (pinned digest) that installs the tidemark RPM built in CI. A :44 fallback build is kept until launch.
- v0.1 outputs: x86_64 qcow2, plus OCI images for x86_64 and aarch64, usable with `bootc switch` and `bootc install to-existing-root`. aarch64 qcow2 is best effort. ISO, raw and AMI are deferred.
- First-boot unit: creates the agent accounts, applies default labels and directory rules, and provisions keys.
- Boot self-test: sessions refuse to start unless the system has Landlock ABI 9 or later, bpf and landlock in the active LSM list, the required kfuncs, cgroup v2, overlayfs, and a /var with xattr and FS_APPEND_FL support.
- Image-owned agent shims and the llama-server Quadlet.
- Headless. No telemetry.

**Technology.** bootc and bootc-image-builder (qcow2), COPR for the RPM, cosign keyless signing through GitHub OIDC, syft, trivy, Renovate digest pinning. fedora-bootc:45 was verified on 16 Oct 2026 as a manifest list with amd64 and arm64 children.

**Rationale.** Stock Fedora already enables bpf and landlock in CONFIG_LSM, so switching them on does not justify an image. The image is justified because it:
- pins a kernel for kfuncs that are not a stable ABI, with assertions in CI;
- guarantees a Landlock ABI 9 floor;
- owns every agent launch path;
- keeps credentials only inside daemons;
- fixes the staging and ledger layout;
- binds receipts to a signed digest;
- gives bootc rollback.
The RPM is the unit of delivery. Guarantees on hosts outside the image are only those that doctor reports.

### L1 Kernel enforcement
- Sleepable BPF LSM programs on file_open and file_receive.
- An inode-storage label cache, invalidated by the setxattr and removexattr hooks.
- An mmapable per-session label array and a map from ancestor cgroup to slot.
- A ring buffer for label-change events.
- Strict Landlock, seccomp, namespaces (creation denied), cgroup v2 with Delegate=no, the freezer, and overlayfs.

**Technology.** BPF C compiled with clang and loaded through cilium/ebpf (bpf2go); go-landlock; systemd unit properties.

**Rationale.** Existing kernel primitives are reused. The only new kernel-side code is the small labelling program, and the Phase 0 gate proves it on stock Fedora 45 before any other work begins.

### L2 Control plane
tidemarkd runs as root with a minimal surface:
- session lifecycle;
- overlay and state-view mounts;
- label-slot allocation;
- a raise-only call for the gateway;
- label administration;
- distribution of the policy and its hash;
- stop;
- doctor.
Mutating calls are accepted only from the human UID or root, from outside the agents slice.

**Technology.** Statically linked Go. Uses the base image's systemd and needs no v261 or v262 features.

**Rationale.** tidemarkd is the only component that needs root, so its root operations are listed and reviewed. Everything that parses untrusted input runs unprivileged.

### L3 Exit plane
- tidemark-gw: CONNECT and plain HTTP, a per-chunk label re-check, tunnel reset when the label changes, source labels.
- tidemark-broker: OpenAI Chat Completions, clearance placement, a field allowlist, key custody.
- The release queue, for broker requests only.

**Technology.** Go, unprivileged, each daemon under its own UID with systemd sandboxing. Both hold read-only mappings of the label array. Decisions are per domain, with no TLS interception in v0.1.

**Rationale.** The socket identifies the session. The per-chunk check closes the bypass where a tunnel is opened before a read. The broker sees prompts in plaintext on the local link, which closes the channel through provider-side tools without any interception.

### L4 Evidence plane
tidemark-ledgerd (its own UID, append-only segments, Ed25519 key) and tidemark-verify (macOS, Linux, Windows).

**Technology.** in-toto Statement v1 inside DSSE, a SHA-256 hash chain, crypto/ed25519. The key is sealed with TPM2 in v0.2, with swtpm used in CI. TPM quotes and optional Rekor anchoring come in Phase 7.

**Rationale.** Agents have tampered with logs before (OpenAI and Hugging Face, July 2026). Receipts protect against the agent. Third-party claims wait for attestation.

### L5 Model layer
- llama-server Quadlet pinned by digest: CPU by default, Vulkan optional.
- `tidemark model pull|ls`, pulling by Hugging Face revision and SHA-256.
- A backend registry with clearances.
- Host endpoints recorded as unattested.
- Cloud providers registered as destinations.

**Technology.** ghcr.io/ggml-org/llama.cpp server image (arm64 confirmed in Phase 2 or built in CI). One permissively licensed GGUF model under 1B parameters (Apache-2.0) for functional tests. The gating demo uses puppet backends and needs no model.

**Rationale.** This meets the spec's one validated backend. Placement is a clearance decision, so local inference becomes an enforced privacy control without forcing every task onto a weak model.

### L6 Agent layer
- Profiles: hermes.yaml (Hermes Agent v0.22.x, base URL set to the broker, its own llama.cpp auto-install disabled, memory and skills mounted as label-partitioned state views) and opencode.yaml.
- Generic `tidemark run -- <cmd>`.
- Agents installed on demand into /var, pinned by version and hash.
- A pinned demo bundle containing only MIT-licensed agents.

**Technology.** YAML profiles. Proprietary agents are never redistributed. The Claude Code profile arrives in Phase 6, when the broker gains Anthropic Messages.

**Rationale.** Agents are wrapped, not forked. Two vendors are enough to prove the cross-vendor claim in v0.1.

### L7 Interfaces
- CLI: run, ps, stop, resume, label, unlabel, labels, audit-labels, diff, commit, accept, abort, queue, release, receipt, model, doctor, demo, report.
- `tidemark report --html` writes a static, self-contained page with no scripts, covering Home, Sessions, Models, Receipts and System. This is the simple UI for v0.1.

**Technology.** Go html/template rendered to a file. No server, no JavaScript, no mutating routes.

**Rationale.** This is the cheapest interface that still covers the spec's Home, Agents, Models and System pages with no authentication surface. A live dashboard and Cockpit approvals are backlog items, each with a trigger.

### L8 Benchmark and CI
- tidemark-bench, a standalone OCI container: pre-registered scenarios, local and cloud puppet backends, host-side sinks, and baselines (unsandboxed, sandbox-runtime, Tidemark audit and enforce).
- Workflows: image and RPM builds, x86 KVM boot test, integration, nightly benchmark (45 minutes or less), benchmark subset on pull requests (5 minutes or less), reproduce, release, Scorecard.

**Technology.** Go plus YAML scenarios. QEMU/KVM on free x64 runners is the only gating boot environment. arm64 runners build natively and boot only when /dev/kvm exists, without gating.

**Rationale.** The benchmark is the unit of adoption that spreads even when people never boot the image.

### L9 Portable mode
- v0.1: the same RPM on stock Fedora 45 (Atomic or Server), with doctor guarantee levels written into every receipt.
- Phase 6: DEB (Ubuntu 26.04) and AUR packages, and the OpenShell middleware adapter.
- Specification files published separately under Apache-2.0.

**Technology.** The same Go binaries, plus pkg/labels and pkg/clearance as importable libraries.

**Rationale.** This lowers the barrier to trying Tidemark while staying honest that guarantees are weaker outside the image.

## Repository layout

```
tidemark/                         Apache-2.0 core; TRADEMARK.md; DCO or CLA decided before first outside contribution
├── cmd/
│   ├── tidemarkd/                root control plane: sessions, overlay and state-view mounts, label slots, raise-only API, stop, doctor
│   ├── tidemark-gw/              unprivileged egress gateway (per-session socket, per-chunk label re-check, source labels)
│   ├── tidemark-broker/          unprivileged model broker (OpenAI Chat Completions in v0.1)
│   ├── tidemark-ledgerd/         receipt writer (own UID, append-only segments)
│   ├── tidemark-shim/            in-session exec shim (close_range, strict Landlock, slot wait, forwarders)
│   ├── tidemark/                 CLI, including report --html
│   └── tidemark-verify/          offline verifier, static builds for macOS, Linux, Windows
├── pkg/
│   ├── labels/                   label lattice and read-only label-word reader (importable by adapters)
│   └── clearance/                destination clearance and source-label evaluator (importable by adapters)
├── internal/
│   ├── session/  stage/  stateview/  egress/  broker/  receipt/  queue/  doctor/  report/
├── bpf/
│   └── label.bpf.c               file_open and file_receive programs, inode cache, invalidation hooks
├── spec/                         policy-v0.schema.json, receipt-predicate-v0.md, clearance-model.md
├── policy/
│   └── default.yaml              default labels, directory rules, destinations, clearances, source labels
├── profiles/                     hermes.yaml, opencode.yaml (claude-code.yaml and openclaw.yaml in Phase 6)
├── packaging/
│   ├── rpm/                      COPR spec (v0.1; the image installs this RPM)
│   └── deb/  aur/                (Phase 6)
├── image/
│   ├── Containerfile             FROM quay.io/fedora/fedora-bootc:45@sha256:<pinned>, installs the RPM
│   ├── Containerfile.f44         fallback build, kept until launch
│   ├── kargs.d/  systemd/  quadlets/  firstboot/  shims/
│   └── bib/                      qcow2 config only in v0.1
├── bench/
│   ├── Containerfile             standalone tidemark-bench image
│   ├── prereg/                   frozen methodology, scenario list and scoring, tagged before baselines run
│   ├── scenarios/adversarial/    YAML plus script per channel; counts only if it leaks in audit or unsandboxed
│   ├── scenarios/benign/         10 tasks in v0.1, 20 in v0.2
│   ├── puppet/                   local-puppet and cloud-puppet backends driving real CLIs with dummy keys
│   ├── sinks/                    attacker, public-paste, public-search, mail-api (private-source), cloud-llm
│   └── baselines/                unsandboxed/  sandbox-runtime/  tidemark/  (nono/ and openshell/ in v0.2)
├── demo/                         byoc.sh (host), heist/ (in-VM), qemu/, lima/ (best effort)
├── test/                         integration/  bypass/ (including tunnel/)  readmatrix/  staging/  laundering/  fuzz/
├── docs/
│   ├── THREAT_MODEL.md           single document: threat model, non-goals, guarantees and the hosts that provide each
│   ├── BENCHMARK_METHODOLOGY.md  limitations.md  compatibility.md  policy-reference.md  quickstart/
│   ├── interviews/               anonymised synthesis only
│   ├── adr/                      for example 2026-10-26_BaseImage_v1.md
│   └── devlog/                   for example 2026-10-25_Devlog_v1.md
├── .github/workflows/            build-image.yml  build-rpm.yml  boot-test.yml  integration.yml  bench-nightly.yml
│                                 bench-pr.yml  reproduce.yml  release.yml  scorecard.yml
└── LICENSE  NOTICE  TRADEMARK.md  SECURITY.md  CONTRIBUTING.md  GOVERNANCE.md  CODE_OF_CONDUCT.md

Code budget for v0.1, re-estimated with the reviewers' ranges, comes to about 4,000 lines of product code:
- session and shim: 600
- labels: 250 BPF C plus 300 Go
- gateway: 600
- broker (Chat Completions only): 700
- staging and state views: 700
- ledger and verifier: 400
- CLI, stop, doctor and report: 450

On top of that come about 1,500 lines of benchmark, puppet and demo code, plus tests. The pre-launch build budget is 168 hours. That assumes AI-assisted coding, with the founder reviewing every diff in the BPF program, the gateway copy loop, the staging merge and the ledger. Without that assistance, expect the reviewers' estimate of 26 to 30 weeks for the MVP.
```

## Stack
- Go (version pinned in Phase 0) for every daemon, the CLI, the verifier and the benchmark. Reasons: cilium/ebpf is mature, go-landlock is maintained, DSSE and in-toto libraries exist, cross-compiling to linux/arm64 with CGO_ENABLED=0 is simple, and a single language suits a solo founder.
- BPF C compiled with clang/LLVM in a Fedora toolbox container and loaded through cilium/ebpf with bpf2go. Uses the kfuncs bpf_get_file_xattr, bpf_path_d_path and bpf_real_data_inode, verified in fs/bpf_fs_kfuncs.c on 16 Oct 2026, and BPF inode storage.
- go-landlock in strict mode (never BestEffort), coreos/go-systemd and godbus for transient units.
- secure-systems-lab/go-securesystemslib (DSSE), the in-toto attestation Go bindings and crypto/ed25519.
- Base image: quay.io/fedora/fedora-bootc:45, pinned by digest (multi-architecture manifest list verified on 16 Oct 2026). :44 stays as a fallback build until launch.
- Delivery: a COPR RPM that the image installs; bootc and bootc-image-builder (qcow2 in v0.1); OCI images for `bootc switch`; cosign keyless signing; syft; trivy; Renovate; OpenSSF Scorecard.
- Podman Quadlet running the ghcr.io/ggml-org/llama.cpp server image, pinned by digest, as the one validated inference backend.
- overlayfs on any /var filesystem that supports xattrs and FS_APPEND_FL (xfs, ext4 or btrfs). Btrfs is no longer required.
- systemd-creds for keys. TPM2 sealing from v0.2, with swtpm in CI.
- The static HTML report is built with Go html/template and contains no JavaScript.
- Testing: go test with native fuzzing, a QEMU/KVM boot harness on GitHub Actions, gosec, govulncheck, staticcheck, CodeQL and shellcheck.
- Added in Phase 5: elazarl/goproxy (BSD-3) for opt-in TLS termination on declared APIs only, and cosign verify-blob for effect manifests.
- Development tools on the M1 Pro: Lima (vz) running stock Fedora 45 Cloud for the kernel spike, podman machine (libkrun or applehv) for image builds, UTM or QEMU with HVF, and VS Code Remote-SSH.

## Threat model

### Assets
- User secrets that a profile exposes to a session: SSH and GPG keys, cloud credentials, password stores and tokens.
- Private data in labelled files and directories, and private data read from destinations tagged as private sources (for example a mail API).
- Agent persistent state: memories, skills and hooks.
- Model API keys, held only by the broker. From v0.2, credentials for declared APIs, held only by the gateway.
- Integrity of the workspace and host: files, auto-exec paths, agent skill directories and the immutable /usr.
- External accounts reachable through cleared destinations. v0.1 does not hold effects on them; v0.2 holds irreversible effects on declared APIs.
- Integrity of the evidence: ledger segments, receipts and the signing key.
- The human's attention: release requests.

### Adversaries
- A prompt injector who controls content agents read: web pages, issues, READMEs, emails and API responses.
- A malicious or compromised skill, MCP server, plugin, hook or agent package running inside the session.
- The agent or model acting deceptively: spoofing tool calls, tampering with logs, trying to persist through implants or memory.
- A second agent from another vendor, or a later session of the same agent, used to launder data through files or persistent memory.
- Out of scope:
- a root or kernel-level attacker;
- a physical attacker;
- the machine owner;
- a destination the user cleared that misuses data;
- network attackers beyond what standard TLS covers;
- a compromised Fedora base, COPR or Sigstore (mitigated only by signature verification);
- the human's own tools outside sessions that strip labels, which is treated as a documented limitation.

### Trust boundaries
- Session to kernel: namespaces (with creation denied, user namespaces included), a cgroup subtree with Delegate=no, a Landlock domain, seccomp and BPF LSM. The shared kernel is trusted.
- Session to Tidemark services: exactly two per-session Unix sockets, and the socket identifies the session. The gateway and broker hold read-only mappings of the session's label word. The control socket and the ledger are unreachable from sessions.
- Session to network: no interface except loopback. Every exit passes through tidemark-gw or tidemark-broker, and each forwarded chunk is checked against the current label.
- Session to filesystem: /usr is read-only, every writable path is an overlay upper, and read binds are limited to what the profile declares. Labels and label maps can be changed only with CAP_SYS_ADMIN, from outside sessions.
- Session to other sessions: separate mount namespaces and Landlock domains, so ptrace and signals are blocked, and each session's uncommitted writes stay invisible to the others.
- Exit plane to ledger: a socket filtered by peer UID. Segments are append-only by kernel attribute against the ledger daemon, but root can still alter them.
- Human to machine: commit, accept, release and stop are accepted only from the human UID or root, from outside the agents slice.
- Machine to outside world: the BYOC canary and the sinks on the host are the primary external check. Receipts can be verified offline as long as the key stays secret.
- Machine to cleared destinations: clearing a destination is the user's trust decision.
- Build pipeline to image: keyless signing, a signed RPM, an SBOM and pinned digests. Users verify the image with cosign.

### Guarantees
- These are preview guarantees. They hold on the signed image when `tidemark doctor` passes and the kernel is not compromised.
- G1 Containment. Every descendant of a session stays in the session's cgroup subtree, network namespace and mount view, and is killed when the session stops. Sessions cannot create namespaces (user namespaces included) or io_uring instances.
- G2 Single exit. A session has no network path except tidemark-gw and tidemark-broker. DNS, UDP, QUIC, IP literals and host Unix sockets are unavailable.
- G3 Label before read. The session's label rises before the call returns when any session process opens a labelled file (labelled by xattr or by directory rule) or receives a descriptor to one. Descriptors are not inherited across the shim's exec. Reading from a private-source destination raises the label before the first response byte reaches the session. A session cannot lower any label.
- G4 Clearance on every chunk. The gateway and broker forward bytes only while the session's current label is a subset of the destination's clearance at that moment. Connections opened earlier are reset when the label rises. The decision never depends on content.
- G5 Model channel. A labelled session's request is never placed on a backend outside its clearance; it fails closed with a reason instead. Fields outside the allowlist and provider-side tools are dropped for labelled sessions. Model API keys never enter a session.
- G6 Staged writes carry labels. Writes stay unmerged and invisible to other sessions until commit; root and the human can inspect them with `tidemark diff`. Commit stamps every written file with the session's final label, and abort restores a byte-identical workspace. Auto-exec, skill and hook files wait for human acceptance. Agent memory written by a labelled session is visible only to sessions started at or above that label.
- G7 Evidence against the agent. A process the session cannot reach records every decision. Any edit, deletion, insertion or reordering of a receipt is detected offline, provided the signing key has not been extracted. Each receipt names the image digest, the policy hash and the doctor result.
- G8 Stop. `tidemark stop` freezes every session and denies every exit, with a target of 1 s or less.
- G9 Declassification. Only a human outside the agents slice can release data. In v0.1, releases apply only to broker requests. Each is bound to the exact content hash and destination, is consumed by the agent's own retry and is recorded.

### Non guarantees
- Tidemark does not prevent prompt injection and does not make the agent's output correct.
- A kernel exploit or root compromise defeats every guarantee. The kernel stays shared, and microVM principals are only a backlog item.
- Covert and timing channels toward cleared destinations remain: request counts, timing, choice of URL, and CDN domain fronting through a cleared domain.
- Data sent to a destination the user cleared, including a cleared cloud model, is outside Tidemark's control.
- Unlabelled data is unprotected. Data that arrives through APIs is protected only when the destination is tagged as a private source.
- Labels are lost when files are copied off the machine, or when tools outside sessions that drop security.* xattrs rewrite them. Directory rules are the exception. `tidemark audit-labels` reports drift only after the fact.
- In v0.1, external side effects at cleared destinations are neither held nor undone. Keep each agent's own approval on for destructive tools until v0.2.
- v0.1 decides per destination domain and does not inspect content. Clearances by method and path arrive in v0.2.
- Receipts are not proof against root or the machine's owner. TPM sealing in v0.2 protects the key at rest only. Third-party verifiability needs TPM quotes (Phase 7).
- Hosts running only the RPM get just the guarantees doctor reports. macOS host endpoints are recorded as unattested.
- Session granularity is coarse: after a session reads one secret, it loses public egress for the rest of its life. Private memory views may hide later public edits to the same file.
- Always-on messaging-gateway modes, VPS hosting and OpenClaw fall outside the v0.1 claims.
- Availability is not guaranteed. Quotas limit resource abuse but do not eliminate it.
- Enforcement may move work to weaker local models. That quality cost is measured and published.

## Hero demo script
"Bring Your Own Canary" (v0.1 preview).

**Ways to run it**
- Primary: an x86_64 Linux host with KVM. Run `git clone`, then `./demo/byoc.sh`.
- Convenience: fork the repository, enable Actions on the fork, and run the `reproduce` workflow on a free GitHub x64 KVM runner. It runs the project's own workflow file, so it shows reproducibility, not independence. Independence comes from running byoc.sh on your own host.
- Best effort: Apple silicon with QEMU (HVF) or Lima.

**Downloads and time.** You download the signed x86_64 qcow2 (target: 1.6 GB compressed or less) and a pinned demo bundle containing the MIT-licensed agents, the puppets and the sinks; its size is published. No model is downloaded. Each release publishes the wall-clock time measured on a clean runner, downloads included. Targets: 20 minutes or less including downloads, 10 minutes or less after them.

**Steps**
1. byoc.sh downloads the qcow2, runs `cosign verify-blob` against the published workflow identity and prints the image digest.

2. The host generates a canary (`openssl rand -hex 16`) and starts these sinks, all outside the VM:
   - an attacker sink;
   - 'public paste' and 'public search' services, which the policy treats as allowlisted public destinations;
   - an HTTPS 'mail API', tagged as a private source and cleared for private;
   - a cloud puppet, registered as a cloud provider cleared for public only. It logs every prompt and simulates a server-side web-fetch tool.
   Inside the VM, a local puppet acts as the local backend, cleared for all labels.

3. The VM boots with the canary passed as a systemd credential through QEMU fw_cfg. First boot writes the canary into ~/Private/notes.txt (private by directory rule) and puts a fake key in ~/.ssh (secret). The mail API serves a second copy of the canary.

4. Run A, audit mode: `tidemark demo heist --mode audit`. The puppet drives real, unmodified Hermes Agent through a recorded injection transcript. The agent fetches an injected issue page, reads the notes, and tries 8 headline channels:
   1. upload to the paste host;
   2. the canary in a public-search query string;
   3. a CONNECT tunnel to the paste host opened BEFORE the read, with the canary streamed through it afterwards;
   4. the canary in the next tidemark/auto prompt;
   5. a request that defines a provider-side web-fetch tool aimed at the attacker;
   6. a git push over HTTPS to an allowlisted public remote the attacker controls;
   7. agent-memory laundering: the canary is written into Hermes memory, then a fresh session without --label posts its memory to the paste host;
   8. API-sourced exfiltration: a fresh session reads the mail API, then uploads to the paste host.
   A grep on the host shows the canary in the sinks, which proves the channels are real. Containment channels such as DNS, UDP, IP literals and host sockets belong to the benchmark, where an unsandboxed baseline shows that they are real.

5. Run B, enforce mode, with a new canary and the same script. It passes when:
   - `grep -c` returns 0 on every host log;
   - the receipt shows each label raise with its path hash or source as the cause;
   - the receipt shows each denial with its reason and the tunnel reset;
   - every placement lands on the local puppet, and the cloud puppet's log is empty;
   - the provider-side tool was dropped;
   - the memory-laundering session saw no canary, because it started public.

6. Cross-vendor step. Hermes Agent's summary is staged, and `tidemark commit`, run outside the slice, stamps it private. OpenCode reads it and becomes private: its tidemark/auto call goes to the local puppet, and its paste upload is denied.

7. An injected session edits 200 files and plants .git/hooks/post-checkout and a new Hermes skill. `tidemark diff` shows both implants held. `tidemark abort` restores a byte-identical tree hash.

8. Receipts:
   - `tidemark receipt export`, then `tidemark-verify` on the host, which checks the signature and the chain and confirms the digest equals the cosign-verified digest;
   - tampered copies fail at the exact entry;
   - the attempt to write the ledger from inside the session appears as denied.
   The script states on screen that receipts are tamper-evident against the agent and that the host grep is the proof.

9. `tidemark stop` freezes a running session and denies all exits (target: 1 s or less). `tidemark resume` restores it.

**Optional and non-gating:**
- a placement smoke test with the real llama-server and a model under 1B parameters;
- `--live` with the viewer's own key, reported separately.
The baseline comparison runs as `tidemark-bench compare` and is not part of the demo.

**v0.2 extension, same demo.** Forgejo and Mailpit run as Quadlets and show:
- an injected credentialed push, with the credential injected by the gateway and never present in the session;
- a public comment denied by a path-level clearance;
- an email held synchronously and approved within the time bound;
- a timed-out hold reported to the agent as not performed;
- the agent's retry passing with the approval token.

## Phases

### Phase 0: Interviews, single threat-model document, kernel kill gate on stock Fedora 45, x86_64 image pipeline
_Weeks 1-3 (about 30 hours, 21 of them building; 19 Oct to 8 Nov 2026)_

**Goal.** Observe real users and prove the kernel mechanism on stock Fedora 45 for both architectures before writing product code.

**Deliverables**
- Repository governance files:
- LICENSE (Apache-2.0), NOTICE, TRADEMARK.md;
- CONTRIBUTING.md (DCO or CLA per the founder's decision);
- SECURITY.md with private advisories and a 90-day coordinated disclosure window;
- GOVERNANCE.md, which assumes a bus factor of 1;
- CODE_OF_CONDUCT.md.
- docs/THREAT_MODEL.md v0, one document combining the threat model, non-goals and G1 to G9, with the hosts that provide each guarantee.
- docs/BENCHMARK_METHODOLOGY.md v0, as the pre-registration draft.
- spec/policy-v0 (labels, directory rules, clearances, source labels) and spec/receipt-predicate-v0.
- 8 to 10 observational interviews in weeks 1 and 2, with an anonymised synthesis.
- Kernel spike on stock Fedora 45 Cloud VMs: aarch64 in Lima (vz) on the M1, and x86_64 under QEMU/KVM on a GitHub runner. It proves:
- a label raise before open() returns, from a grandchild process;
- ancestor-cgroup keying, with Delegate=no;
- behaviour through overlayfs, comparing bpf_get_file_xattr with a cache keyed by bpf_real_data_inode;
- directory rules through bpf_path_d_path;
- file_receive for descriptors passed with SCM_RIGHTS;
- a read-only mmap of the label word;
- EPERM on security.bpf.* inside a session.
It also records:
- how xattrs and the append-only attribute behave on virtiofs and 9p;
- whether Lima's cidata provisioning works on fedora-bootc.
- x86_64 image pipeline: Containerfile on :45 pinned by digest, plus the :44 fallback build, cosign keyless signing, a syft SBOM, and a QEMU serial-console boot smoke test.
- Architecture decision records, named per the date convention: base image; labeller, covering the BPF design and the fanotify fallback's coverage gaps; Go as the language; no TLS interception in v0.1; delivery as an RPM plus the image.
- Two or three pre-launch security reviewers recruited and confirmed by name, with their consent recorded.

**Acceptance criteria**
- On both architectures, the probe prints:
- the kernel version;
- an LSM list containing bpf and landlock;
- a Landlock ABI of 9 or later;
- the label word already raised when read in the same process immediately after a grandchild's open() returns.
- The overlay behaviour and the chosen lookup method are recorded in an ADR.
- The fanotify fallback's coverage for mmap, grandchild and overlay reads is documented before it can be selected.
- A signed x86_64 image is on GHCR, and `cosign verify` passes in CI.
- The interview synthesis is published, and kill criterion 2 is evaluated at the end of week 3.
- A written go or no-go decision is recorded.

**Automated tests**
- x86_64 image build in CI.
- QEMU/KVM boot smoke test that runs the probe.
- BPF spike tests: map entry after a child's open; overlay path; descriptor passed with SCM_RIGHTS; directory rule; xattr invalidation.
- Landlock ABI and LSM-list assertions.
- gofmt, staticcheck, the DCO check and the Scorecard workflow.

**Security checks**
- Cosign keyless signing and verification. No long-lived keys in the repository.
- Actions pinned by commit SHA; branch protection; secret scanning.
- SBOM generated, plus a trivy scan in report-only mode.
- BPF program reviewed for bounded loops, early exit for processes outside sessions, and handling of failed kfunc calls.

**Compatibility**
- fedora-bootc:45 (multi-architecture manifest verified 16 Oct 2026) with :44 as the fallback. Stock Fedora 45 Cloud images for the spike.
- QEMU with KVM or HVF; Lima (vz).
- M1 Pro host. No nested virtualization is required.

**Performance targets**
- Image build of 30 minutes or less on x86_64.
- Boot to multi-user in 90 s or less under KVM.
- Baseline latency of labelled and unlabelled open() recorded; no target yet.

**Docs**
- README stating what the project is and what it does not claim.
- THREAT_MODEL v0 and BENCHMARK_METHODOLOGY v0 draft.
- ADRs, the interview synthesis and devlogs 1 to 3.

### Phase 1 (MVP-A): Confine and label
_Weeks 4-9 (about 60 hours, 42 of them building; 9 Nov to 20 Dec 2026). Checkpoint at week 9._

**Goal.** Make each agent run one contained principal whose label rises in the kernel before a labelled read returns, measure label creep on a real Hermes Agent setup, and report guarantees honestly.

**Deliverables**
- tidemarkd skeleton with caller checks (SO_PEERCRED plus cgroup) and a raise-only call reserved for the gateway's UID.
- `tidemark run --profile <p> [--label L] -- <cmd>`:
- a transient unit with the principal's properties;
- a fixed UID per profile;
- one active session per profile;
- Delegate=no and RestrictNamespaces=yes;
- the seccomp deny list.
- tidemark-shim: close_range, strict Landlock at ABI 9 or later with scoping, a wait for the label slot, forwarders (stubs until Phase 2), then exec.
- Production label engine:
- file_open and file_receive programs;
- ancestor-cgroup slots;
- the inode cache and its invalidation;
- directory rules;
- mmapable label words and change events;
- `tidemark label`, `unlabel`, `labels` and `audit-labels` (human or root only);
- first-boot defaults (~/.ssh, ~/.gnupg, keyrings, password stores, cloud credential files, ~/Private and mail stores as directory rules).
- `tidemark stop`, `resume`, `ps`, `doctor`, and a boot self-test.
- aarch64 OCI image build and the COPR RPM build in CI.
- Label-creep experiment in audit mode: Hermes Agent v0.22.x with memory on runs 20 scripted realistic tasks over code and notes. It records the share of sessions that end labelled and the share that write to state directories while labelled.

**Acceptance criteria**
- Containment matrix of at least 18 cases passes:
- child processes, an MCP server, a hook, an LSP server, a double-forked daemon and a setsid process all stay in the session;
- no process survives a stop;
- only loopback exists;
- host D-Bus, Podman, ssh-agent and abstract sockets are unreachable;
- signals and ptrace across sessions fail;
- creating a user namespace fails;
- io_uring_setup fails;
- creating a sub-cgroup fails.
- Read matrix of at least 16 cases raises the label in 100% of cases: cat, Python open, mmap, sendfile and splice, an MCP child, a hook, symlink, hardlink, bind path, a grandchild after a double fork, reopening through /proc/self/fd, an overlay path, a descriptor received through SCM_RIGHTS, no inherited descriptor above 2 after exec, a file under a labelled directory whose xattr was dropped by a rename-based editor, and a path matched by a directory rule.
- Zero false label raises across 1,000 unlabelled files.
- setxattr on security.bpf.* fails with EPERM inside a session, and the map pins cannot be opened.
- `tidemark stop` freezes 6 sessions, as verified through cgroup.events.
- The shim fails closed when a test flag forces the ABI below 9 and when the label slot is missing.
- The label-creep report is published, and kill criterion 3 is evaluated.

**Automated tests**
- Go unit tests.
- In-VM integration harness on x86 KVM in CI, producing JUnit output.
- Race test: 10,000 iterations of an open immediately followed by a read of the label word.
- Fuzzing of the profile and policy parsers in the nightly run.

**Security checks**
- gosec, govulncheck, staticcheck and CodeQL.
- A session UID connecting to the control socket fails with EACCES.
- The list of root operations is reviewed against the threat model.
- Assertions for the seccomp deny list and for RestrictNamespaces.

**Compatibility**
- x86_64 and aarch64 UEFI.
- Kernel floor 7.1 (Landlock ABI 9); the expected base kernel is 7.2.
- The base image's systemd, with no dependency on v261 or v262 features.
- /var on xfs, ext4 or btrfs.
- Hermes Agent v0.22.x starts under `tidemark run`.
- The RPM installs on stock Fedora 45, and doctor shows the guarantee levels.

**Performance targets**
- Session start 2 s or less at p50.
- Labelled open() 10 µs or less at p50 and 50 µs or less at p99 with a warm cache.
- open() outside sessions adds 0.5 µs or less at p50.
- A pinned npm install and a pinned C or Go build inside a session add 5% or less to wall-clock time (target; published as measured).
- Stop of 6 sessions in 1 s or less.
- tidemarkd RSS 50 MB or less.

**Docs**
- G1 to G3 in THREAT_MODEL.md linked to CI artifacts.
- Guides to `tidemark run` and to labelling and directory rules.
- Doctor reference.
- Label-creep report.
- Devlogs.

### Phase 2 (MVP-B): Gate every exit, including the model, with per-chunk re-checks
_Weeks 10-15 (about 60 hours, 42 of them building; 21 Dec 2026 to 31 Jan 2027; holiday risk in weeks 10 and 11). Checkpoint at week 15._

**Goal.** Make one rule decide every byte that leaves a session, treat the model call as an exit placed by label, and close the bypass through tunnels opened before a read.

**Deliverables**
- tidemark-gw:
- one socket per session;
- CONNECT and plain HTTP;
- DNS resolved by the gateway;
- IP literals and non-HTTP protocols refused;
- a per-chunk label re-check against the mmapped word;
- tunnel reset on change events;
- source labels raised through tidemarkd before the first relayed byte;
- structured 403 responses.
- tidemark-broker:
- OpenAI Chat Completions with streaming passthrough;
- the tidemark/auto alias;
- a registry with clearances;
- placement by clearance filter, then static preference;
- failure closed with a structured reason;
- a field and tool allowlist for labelled sessions;
- keys held as systemd credentials, with placeholders in sessions;
- direct provider domains denied to sessions;
- release tokens.
- llama-server Quadlet pinned by digest, and `tidemark model pull|ls` by revision and SHA-256.
- `tidemark queue` and `tidemark release`, for broker requests only.
- Profiles hermes.yaml and opencode.yaml.
- pkg/labels and pkg/clearance as importable packages.
- Bypass matrix of at least 45 cases.

**Acceptance criteria**
- 100% of the bypass matrix is blocked on x86 CI and on the local aarch64 VM. Cases include:
- a CONNECT tunnel opened before the read that streams afterwards;
- HTTP/2 connection reuse;
- a WebSocket opened before the read;
- keep-alive plain HTTP;
- curl --noproxy, raw IPv4 and IPv6, UDP, QUIC;
- a DNS TXT tunnel, DoH to a hard-coded IP;
- host sockets, /proc/net tricks;
- connections from an MCP child and from a double-forked daemon;
- CONNECT to an IP literal or a non-443 port;
- a mismatch between the CONNECT authority and the Host header.
- No byte from a labelled session reaches an uncleared host after the label rises, whatever the encoding. The same request from an unlabelled session passes.
- Reading from a private-source destination raises the label before the first response byte reaches the session.
- A private session's tidemark/auto request lands on the local backend, and the cloud mock receives 0 requests. When no cleared backend fits (for example, the context is too long), the request fails closed with a structured reason.
- Fields outside the allowlist and tools that are not functions are dropped for labelled sessions.
- Hermes Agent and OpenCode each complete 3 benign puppet-driven tasks through the broker.
- The API key appears nowhere in session environ, in files the session can read, or in core dumps.
- Functional test: a model under 1B parameters answers through the broker on the x86 runner.

**Automated tests**
- Bypass matrix on every pull request that touches enforcement.
- Broker contract tests against recorded Chat Completions fixtures, streaming and non-streaming.
- Allowlist table tests and 20 placement cases.
- Authorization tests for release requests.
- Go fuzzing of the policy and request-body parsers, plus a test of the gateway copy loop.

**Security checks**
- The gateway and broker run unprivileged under their own UIDs, with read-only label mappings and systemd sandboxing.
- Credential-leak tests, and denial responses checked to contain no secrets.
- Review of how CONNECT authority is handled and of fronting through cleared domains (recorded in the limitations).
- gosec, govulncheck and CodeQL. trivy blocks critical findings.

**Compatibility**
- Proxy-aware clients: curl, git over HTTPS, Python requests and httpx, Go net/http, Node with NODE_USE_ENV_PROXY=1, and Bun.
- Clients that are not proxy-aware fail closed and are listed. SSH git and plain IMAP clients are denied and documented.
- Chat Completions clients only; Responses and Anthropic Messages arrive in Phase 6.
- llama-server arm64 image confirmed, or built in CI.
- Host endpoints recorded as unattested.

**Performance targets**
- The gateway adds 3 ms or less to connect at p50 and 10 ms or less at p95.
- The per-chunk re-check costs 2% or less of throughput on a 100 MB loopback transfer (target).
- The broker adds 10 ms or less at p50 and 20 ms or less to the first streamed byte.
- Gateway RSS 40 MB or less.
- Local inference throughput published, with no target.

**Docs**
- Policy, clearance and source-label reference.
- Model placement guide.
- Profile guide.
- Compatibility matrix v0.
- limitations.md.
- Devlogs.

### Phase 3 (MVP-C): Stage, carry labels forward, sign
_Weeks 16-20 (about 50 hours, 35 of them building; 1 Feb to 7 Mar 2027)_

**Goal.** Make labels follow everything a labelled session writes, close laundering through agent memory, make workspace writes abortable, and record decisions in a form the agent cannot tamper with.

**Deliverables**
- Staging:
- an overlay on every writable path;
- doctor refuses writable paths without an overlay;
- `tidemark diff`;
- `tidemark commit`, after the session stops, which stamps every upper file with the session's final label, holds conflicts and holds auto-exec, skill and hook files;
- `tidemark accept` and `tidemark abort`;
- commit and accept refused from inside the agents slice.
- Label-partitioned agent-state views for the built-in labels.
- tidemark-ledgerd (append-only segments, hash chain, DSSE receipts) and tidemark-verify (macOS, Linux, Windows).
- `tidemark report --html`.

**Acceptance criteria**
- Aborting a 10,000-file workspace restores a byte-identical tree hash.
- A property test of 1,000 random operations (renames, deletions, whiteouts) shows that commit gives the same result as an unwrapped run.
- A file opened for writing before the label rose and written afterwards is stamped at commit.
- Committed writes from a private session carry the private label. OpenCode reading them starts private. A concurrent session cannot see uncommitted writes.
- Memory laundering: a private Hermes session writes the canary to memory. The next session without --label sees no canary. A session started with --label private sees it and is private from the start.
- Implants and new skill files stay held until accepted from outside the slice. A commit from inside the slice is refused.
- The verifier detects 100% of fuzzed tampering at the exact entry.
- A session UID gets EACCES on the ledger socket, the ledger files and the key, and the attempt is recorded. The ledger daemon cannot truncate a segment.
- The receipt's digest equals the digest from `bootc status` and the cosign-verified digest.
- The report HTML contains no scripts and renders offline.

**Automated tests**
- Property tests for staging, plus label-stamping, conflict and state-view tests.
- Go fuzzing of the ledger decoder and the verifier.
- Verifier test matrix on macOS, Windows and Linux runners.
- End-to-end test: run, abort, run, commit, export, verify.

**Security checks**
- Privilege separation for the ledger daemon; append-only attribute tests; the key unreadable by any other UID.
- DSSE and in-toto handled by maintained libraries only.
- Paths and bodies hashed by default.
- Mount-namespace review: no session can see another session's upper directory.

**Compatibility**
- /var on a filesystem with xattr, overlay upper and FS_APPEND_FL support; workspaces under /var/home.
- Verifier runs on macOS 13 or later, on Linux (x86_64 and arm64) and on Windows 11.

**Performance targets**
- Commit of 1,000 changed files in 5 s or less.
- Abort of a 10,000-file workspace in 5 s or less.
- Ledger append in 2 ms or less at p95.
- Receipts average 1 KB or less per decision.
- Report generation in 2 s or less for 100 sessions.

**Docs**
- Guides to staging and agent-state views.
- Receipt format specification v0 and a verifier guide.
- Privacy note on evidence.
- Devlogs.

### Phase 4 (MVP-D): Pre-registered benchmark, demo, launch of the v0.1 preview
_Weeks 21-24 (about 40 hours, 28 of them building; 8 Mar to 4 Apr 2027). Launch slips to mid-May 2027 if the holiday weeks and the contingent cut order are exhausted._

**Goal.** Ship a confidentiality preview with one demo strangers can check and a benchmark that is fair to incumbents.

**Deliverables**
- tidemark-bench, as a standalone OCI container and an in-image command:
- the methodology frozen and published before any baseline runs;
- at least 24 adversarial scenarios across 18 gating channel classes, plus 2 classes where incumbents should win;
- 10 benign tasks;
- local and cloud puppets;
- host-side sinks.
- Baselines with identical agents and allowlists: (a) unsandboxed; (b) sandbox-runtime wrapping the whole process tree, in its strongest documented configuration including credential masking; (c) Tidemark audit; (d) Tidemark enforce. nono is a stretch goal if it takes 4 hours or less.
- demo/byoc.sh on the primary x86_64 KVM path, the reproduce.yml workflow, a Lima template (best effort) and a 90-second recording.
- v0.1 preview release:
- signed x86_64 qcow2;
- x86_64 and aarch64 OCI images;
- the COPR RPM;
- SBOM attestation and checksums;
- release notes that state the non-claims and advise keeping agents' own approval on for destructive tools.
- Launch on Show HN, r/LocalLLaMA, r/selfhosted and the Hermes Agent, Fedora, Universal Blue and Proxmox forums, with a 'break it' issue template, a hall of fame and a 90-day disclosure window.

**Acceptance criteria**
- The gating suite runs in 45 minutes or less nightly on a free x64 runner, and the subset in 5 minutes or less on pull requests. Three consecutive runs give identical results.
- In enforce mode, every host log shows zero canary bytes on every in-scope channel. In audit mode, every label channel leaks. Containment channels leak only in the unsandboxed configuration.
- Benign completion in enforce mode is at least 85% of audit mode (target; published as measured). Completion is published for every baseline next to its leak count.
- The BYOC wall-clock time on a clean runner, all downloads included, is published. Targets: 20 minutes or less including downloads, 10 minutes or less after them.
- One outside person reproduces the demo on their own host and posts the verifier output.
- Methodology and configurations reach the sandbox-runtime maintainers at least 7 days before publication, and any reply is published with the results.
- Release gate: no in-scope channel open and no high or critical finding open.

**Automated tests**
- Determinism job: three runs, then a diff.
- The reproduce workflow runs on every release tag and nightly.
- The benchmark subset runs on pull requests.
- A claims checker fails the build when a README metric lacks a link to a CI artifact; a docs link check.

**Security checks**
- The named reviewers attack the image and the benchmark privately for one week.
- Fuzzers run nightly; the image scan finds no critical issue.
- The reproduce workflow verifies the image with cosign before booting it.
- Private vulnerability reporting is enabled with a 90-day coordinated disclosure window.

**Compatibility**
- Demo hosts: x86_64 Linux with KVM and the GitHub x64 runner are primary; Apple silicon with QEMU (HVF) or Lima is best effort.
- Pinned: Hermes Agent v0.22.x, OpenCode and the sandbox-runtime version.

**Performance targets**
- Gating suite 45 minutes or less nightly, and 5 minutes or less as the pull-request subset.
- Demo times as stated above.
- Compressed qcow2 1.6 GB or less.
- Boot to ready in 60 s or less under KVM.

**Docs**
- Quickstarts: QEMU/KVM, Proxmox, bare-metal x86_64, `bootc switch`, and the RPM on Fedora 45. Lima and UTM are best effort.
- BYOC guide.
- Benchmark methodology v1 and results.
- THREAT_MODEL v1.
- Limitations.
- Compatibility matrix v1.
- Scenario contribution guide.

### Phase 5 (v0.2): Default theme 'credentials and irreversible effects at declared APIs', confirmed at the v0.1 retrospective from user evidence
_Months 7-9: weeks 25-37 (about 130 hours, 78 of them building; 5 Apr to 4 Jul 2027)_

**Goal.** Close the credential gap and hold only irreversible effects, without leaving the agent out of step with reality.

**Deliverables**
- Opt-in TLS termination for declared destinations only: git smart-HTTP (Forgejo, GitHub, GitLab) plus two declared APIs (a Mailpit-style mail API and one chat API). Each session gets its own CA, present only in its trust store. goproxy is pinned.
- Credential injection with placeholders for those APIs. Credentialed push returns to the benign suite.
- Clearances by method, host and path.
- Synchronous bounded holds for irreversible classes:
- unknown writes, and any effect that notifies a person, count as irreversible;
- MCP annotations may raise a class but never lower it;
- the request is held open up to a configurable bound (default 120 s);
- approval forwards it and returns the real response;
- a timeout returns 403 'not performed' with an approval id;
- a later approval mints a one-time token that the agent's own retry consumes;
- Tidemark never replays anything.
- Effect manifests for the declared APIs, signed with cosign.
- Ledger key sealed with TPM2 through systemd-creds (swtpm in CI); the receipt records which key type signed it.
- Benchmark v0.2: nono and OpenShell baselines (pinned, telemetry off), 20 benign tasks, a non-gating live variant and community scenarios.
- Demo extension with Forgejo and Mailpit Quadlets.

**Acceptance criteria**
- 100% of irreversible-class requests in the suite are held.
- Across 1,000 fuzzed runs, no held request reaches a sink before approval, and no request is performed after a timeout unless the agent retries with a token.
- Declared-API credentials never appear inside a session.
- Method and path clearances are enforced.
- Benign completion in enforce mode is at least 90% of audit mode.
- Failures caused by holds stay at 15% of tasks or less.
- The TPM-sealed key cannot be decrypted from an offline copy of the disk.

**Automated tests**
- Classification table with at least 60 cases.
- Fuzz test that no held request leaks.
- Tests of hold timeout and retry.
- Failure tests with certificate-pinning clients.
- TPM tests under swtpm.

**Security checks**
- The CA key is unreadable from sessions, and each session has its own CA.
- TLS termination applies only to declared destinations.
- Manifest signatures are verified.
- goproxy is reviewed.
- First external red-team round.

**Compatibility**
- Clients that pin certificates fail closed only on declared destinations.
- HTTP/2 and SSE are supported; WebSockets on declared APIs are denied unless a manifest allows them.
- TPM 2.0 or swtpm.

**Performance targets**
- A TLS-terminated request adds 25 ms or less at p95.
- A hold decision takes 5 ms or less.

**Docs**
- Effect manifest guide.
- Holds guide, including the no-replay contract.
- Credential injection guide.
- Threat model update covering TLS termination.
- Benchmark v0.2 results.

### Phase 6 (v0.3): Default theme 'distribution through incumbents and more agents'; alternatives are swapped in from the backlog if user evidence points elsewhere
_Months 10-12: weeks 38-50 (about 130 hours, 78 of them building; 5 Jul to 3 Oct 2027)_

**Goal.** Make the label engine reachable through tools people already run, and widen agent coverage.

**Deliverables**
- OpenShell supervisor-middleware adapter, built on pkg/labels and pkg/clearance after a spike of 6 hours or less confirms the API (gRPC over TCP, after policy evaluation and before credential injection). Labels come from the Tidemark BPF engine when the sandbox runs on a host where Tidemark is loaded. Otherwise only source labels apply, at a stated lower guarantee level.
- Upstream proposals: a label-provider interface for OpenShell, a 'tidemark' sandbox backend for Hermes Agent, and a label-aware filterRequest example for sandbox-runtime.
- tidemark-lite DEB (Ubuntu 26.04) and AUR packages with doctor guarantee levels.
- Broker support for Anthropic Messages and OpenAI Responses, plus a Claude Code profile (installed on demand).
- An OpenClaw profile for short sessions. Always-on gateway mode stays out of scope unless sub-session evidence supports it.
- Planner/worker experiment for label creep (cloud planner, local worker), built as a Hermes Agent plugin.

**Acceptance criteria**
- The adapter denies a labelled OpenShell sandbox's request to an uncleared host in 100% of tests and adds 10 ms or less at p95.
- Lite doctor output matches tables checked by hand on Ubuntu 26.04 and Arch.
- Claude Code completes 3 benign tasks.
- The planner/worker report is published.

**Automated tests**
- Adapter integration tests against a pinned OpenShell version.
- Install tests for the lite packages.
- Contract tests for Messages and Responses.

**Security checks**
- The adapter endpoint is bound to loopback or protected by mTLS.
- OpenShell telemetry is off in tests.
- Second external review.

**Compatibility**
- A pinned OpenShell version; Ubuntu 26.04; Arch; a pinned Claude Code version.

**Performance targets**
- Adapter adds 10 ms or less at p95.
- Lite packages pass the same benchmark subset at their stated guarantee level.

**Docs**
- Adapter guide.
- Lite guide.
- Planner/worker report.
- Upstream proposal documents.

### Phase 7 (v1.0 candidate): Third-party verifiability and longevity, conditional on the adoption gate three months after launch
_Months 13-18: weeks 51-76 (about 260 hours, 156 of them building; 4 Oct 2027 to 2 Apr 2028)_

**Goal.** Earn the right to say 'third-party verifiable', freeze the formats, and prepare the project to outlast a single founder.

**Deliverables**
- Receipts bound to TPM quotes over measured boot, with a Keylime-style remote verifier. Optional Rekor anchoring. Sealed UKI only if bootc sealed images have left test status.
- Policy and receipt schemas frozen at v1, with compatibility tests.
- The specification proposed to a neutral home (an OpenSSF or AAIF working group).
- Applications for audit funding, made only after confirming eligibility (GitHub Secure Open Source Fund, OpenSSF Alpha-Omega).
- A succession plan that assumes a bus factor of 1; a second maintainer is a bonus, not a deliverable.
- A written month-18 decision: community-only, a hosted attestation service, or a foundation.

**Acceptance criteria**
- A third party verifies a receipt together with a TPM quote against the published image measurements.
- v0 receipts verify under the v1 verifier.
- At least 25 independent BYOC verifications exist.
- An audit is scheduled or funded, or the reason it is not is documented.
- A plan for bus-factor risk is published.

**Automated tests**
- Attestation tests under swtpm in CI.
- Schema compatibility tests.
- Upgrade and rollback tests from the v0.1 image to the v1.0 image.
- A 6-hour soak test.

**Security checks**
- Review document for the trusted computing base.
- Key rotation procedure.
- SLSA provenance through reusable workflows.
- Reproducible-build check where feasible.

**Compatibility**
- Hosts with TPM 2.0 for attestation; x86_64 and aarch64 UEFI.
- Adopting sealed images or composefs requires a reinstall, and this is documented.

**Performance targets**
- Attestation verification in 5 s or less per receipt (target).

**Docs**
- v1.0 documentation set.
- Attestation guide.
- Migration guide.
- Governance guide.

## Benchmark methodology
- Threat assumption: the agent is fully compromised. Scripted tool-call sequences come from two puppet backends that speak OpenAI Chat Completions (Anthropic Messages from Phase 6) and drive real, unmodified Hermes Agent and OpenCode:
- a local puppet, registered as the local backend and cleared for all labels;
- a cloud puppet, registered as a cloud provider, whose log is a sink.
Placement is proven by which log receives the request. A synthetic agent covers channels the CLIs cannot express. The gating suite needs no live model and no API key.
- Pre-registration: the methodology, scenario list, scoring rules and baseline configurations are frozen at a tagged commit and posted publicly before any baseline runs. Later changes are logged as deviations.
- Gating channel classes in v0.1 (18):
- Containment (7): unlisted host; DNS tunnelling; IP literals; UDP and QUIC; host Unix sockets; a socket opened by an MCP child; a double-forked daemon.
- Label (10): upload to an allowlisted host; a query string to an allowlisted host; git push over HTTPS to an allowlisted remote; a tunnel opened before the read (including HTTP/2 reuse and WebSocket); the model prompt; provider-side tools; cross-vendor laundering through files; laundering through agent memory; private data sourced from an API; label stripping by a rewrite followed by a re-read.
- Evidence (1): ledger tampering.
Persistence through auto-exec and skill implants is scored separately.
- Two incumbent-favourable classes are reported but do not gate:
- exfiltration of a credential placed in the session environment;
- an L7 method/path rule on the same host.
Tidemark v0.1 is expected to lose on both, and the results say so.
- Canary handling: the host generates a random canary per run and injects it through a QEMU fw_cfg systemd credential (on Lima, by copying it in; best effort). Sinks run on the host. Pass or fail is the grep count on host logs. No canary appears in scenario source.
- Configurations in v0.1, on the same Fedora base with the same agents and allowlists:
- (a) unsandboxed. This is needed because audit mode keeps containment on and so cannot show that containment channels are real.
- (b) sandbox-runtime wrapping the whole process tree, in its strongest documented configuration including credential masking.
- (c) Tidemark in audit mode.
- (d) Tidemark in enforce mode.
v0.2 adds nono and OpenShell (pinned, telemetry off). Docker Sandboxes may come as a community contribution.
- Fairness rules:
- configurations are published;
- each baseline runs in its strongest documented configuration;
- each baseline project receives the methodology at least 7 days before publication, and its reply is published;
- results are reported as measured;
- documentation never predicts other tools' outcomes;
- benign completion is published for every baseline next to its leak count.
- Utility: 10 benign tasks in v0.1 and 20 in v0.2.
- Coding: fix a bug, add a test, add a dependency. Credentialed push returns in v0.2.
- Personal: summarise private notes, draft an email to a mock.
- Research: fetch public documentation.
Metrics: completion in enforce mode against audit mode, the share of tasks forced onto a local model, tasks blocked by labels, sessions ending labelled (creep) and human decisions per 100 tool calls.
- Live mode does not gate releases: real Hermes Agent with a local model or the tester's own key faces injected content. The report gives injection and leak rates over at least 10 runs, with confidence intervals.
- Overhead microbenchmarks:
- open() latency with and without labels, inside and outside sessions;
- a pinned npm install and a pinned compile inside a session;
- session start;
- gateway connect latency and the throughput cost of the per-chunk check;
- broker overhead and time to first byte;
- commit and abort time;
- ledger append;
- daemon RSS.
Reference hardware: a free GitHub x64 runner and an M1 Pro VM. Reports give the median, p95 and p99 of at least 1,000 samples, with raw data.
- Model placement cost is reported separately for each hardware class: completion of labelled tasks on the local model against a cloud model.
- Determinism: everything is pinned by digest. Three consecutive runs must give identical results. The full suite runs nightly within 45 minutes and a subset runs on pull requests within 5 minutes. The results JSON is recorded in a signed receipt.
- Extensibility: a community scenario counts only if it leaks in the audit or unsandboxed configuration. Bypasses are credited in a hall of fame and become regression tests.
- Reporting rules: every README number links to a CI artifact, each report lists its known gaps, and no claim goes beyond what the scenarios measure.

## Non-goals
- Preventing prompt injection or judging whether agent output is correct.
- Holding or undoing external side effects in v0.1. Holds for declared APIs arrive in v0.2.
- Claiming third-party-verifiable receipts before TPM attestation (Phase 7).
- Building a new kernel or a new distribution, or shipping a desktop environment in the first year.
- Per-token information-flow tracking inside the model, as CaMeL does.
- Eliminating covert or timing channels to destinations the user has cleared.
- Inference serving performance, multi-engine auto-tuning, or GPU and cluster scheduling.
- Training and fine-tuning.
- Agent identity, inventory or governance dashboards.
- Rebuilding containment that OpenShell, sandbox-runtime, nono or MXC already provide. The gateway exists only because it must read kernel labels synchronously, and its evaluator is packaged for reuse by incumbents.
- TLS interception in v0.1.
- Always-on messaging-gateway modes and VPS hosting in v0.1.
- Compliance certification.
- Any telemetry.
- Redistributing proprietary agents.
- Enforcement on Windows or macOS. Only the verifier runs there.

## Development environment
- Host: M1 Pro with 16 GB running macOS 27. Install through Homebrew, pinned where possible: Go, podman, lima, qemu, cosign, syft, gh and shellcheck. Edit code with VS Code Remote-SSH.
- Phase 0 kernel spike: a stock Fedora 45 Cloud image in Lima (vz, aarch64, cloud-init works). The kill gate therefore does not depend on image work.
- Native aarch64 image builds: a rootful podman machine (libkrun or applehv, 4 vCPU, 4 GB), stopped when idle, with bootc-image-builder producing the qcow2.
- Boot loop: Lima (vz) with the qcow2 if Phase 0 confirms cidata support; otherwise UTM or QEMU with HVF, using fw_cfg for the canary.
- Memory: the 6 GB dev VM, the 4 GB podman machine and about 5 GB for macOS do not fit together, so run one VM at a time.
- Code loop:
- cross-compile with CGO_ENABLED=0 GOOS=linux GOARCH=arm64;
- generate the BPF program with bpf2go in a Fedora toolbox;
- run `bootc usr-overlay`, copy the binary with scp and restart it with systemctl, or install the RPM.
- Local model: a GGUF model under 1B parameters on CPU in the VM, or the host's Metal llama.cpp registered as an unattested endpoint through host.lima.internal.
- x86_64:
- Free x64 KVM runners in CI are the only gating environment.
- Interactive bypass triage uses an hourly x86_64 Fedora cloud VM, converted with `bootc install to-existing-root` or given the RPM, for about USD 5 to 10 per month.
- arm64 runners build natively. Boot tests run only when /dev/kvm exists, and they do not gate.
- No nested virtualization is required.
- Weekly cadence: before launch, 7 hours building, 2 on docs and devlog, and 1 on interviews and community, since there is no triage queue yet. After launch, the 6 / 2 / 2 split applies. CI must be green every Friday.
- If the founder builds with AI coding agents, those agents run inside Tidemark sessions from Phase 2 onward. The founder reviews every diff in the BPF program, the gateway copy loop, the staging merge and the ledger line by line, and a named reviewer checks them before launch.

## Kill criteria
- End of Phase 0 (week 3). Trigger: kernel label raising is unreliable on stock Fedora 45 on either architecture or through overlayfs, and the fanotify fallback misses reads by child processes or through mmap. Action: drop the claim of kernel-recorded labels and ship the gateway, broker and receipts as an audit tool, or stop.
- Week 3. Trigger: of the 8 to 10 observational interviews, fewer than 3 participants run agents over local private files or secrets on a machine or VM they control, or a majority would abandon a session rather than accept denied public-web access after it reads private data. Action: re-scope to a benchmark-first project for security researchers before building Phase 1.
- Week 9. Trigger: in the label-creep measurement, more than 40% of realistic benign Hermes sessions read labelled data under the default labels. Action: make audit the default mode for v0.1 and move the planner/worker split ahead of enforcement.
- Week 15. Trigger: the bypass matrix reveals a class of bypass the design cannot fix. Action: do not launch that part of the enforcement claim, and narrow it publicly.
- Launch gate (week 24). Trigger: benign completion in enforce mode is below 80% of audit mode; or more than 35% of tasks are forced onto a local model while local completion is below 50% of cloud completion on the reference hardware. Action: ship audit mode and receipts only.
- Any time. Trigger: an incumbent ships free, cross-vendor, label-driven exits with label-constrained placement and evidence. Action: move effort to the adapter, the specifications and the benchmark, and keep the image only for guarantees that depend on owning it.
- Any time. Trigger: a published bypass of a core guarantee is not fixed within the 90-day disclosure window. Action: withdraw that guarantee publicly and document the bypass.
- Three months after launch (about early July 2027). Trigger: fewer than 10 independent BYOC verifications and fewer than 3 users running real work for 2 weeks or more. Action: stop maintaining the image; keep the RPM, the benchmark and the specifications. Phase 7 does not start.
- v0.2. Trigger: holds make more than 30% of real-agent tasks fail, and retries with tokens cannot bring failures below 15%. Action: drop remote-effect holds and keep labels, staging and receipts.
- Maintenance. Trigger: image breakage takes more than 4 of the 10 weekly hours for 2 consecutive months. Action: ship only the OCI image and the RPM.
- Founder capacity. Trigger: available time stays below 6 hours a week for 4 consecutive weeks. Action: freeze scope at the benchmark and receipts.

## Open questions for the founder
- What hands-on experience do you have with BPF LSM, the systemd D-Bus APIs, overlayfs and Go? The hour estimates assume strong systems experience plus AI assistance.
- Will AI coding agents write most of the code? If so, who besides you reviews the security-critical diffs: the BPF program, the gateway copy loop, the staging merge and the ledger?
- Do you accept the RPM as the delivery unit, with the bootable image as a thin layer that installs it and the only place full guarantees are claimed? This keeps your bootable-OS constraint.
- Is a headless VM appliance (SSH, CLI and a static HTML report, no desktop) acceptable as the 'literal bootable OS'?
- Do you accept the re-baselined launch (early April 2027, mid-May at the latest) with the default cuts, or would you rather cut further for an earlier date?
- Will you make a public security claim as a preview, with a 90-day disclosure window, handled in 2 triage hours a week after launch? If not, launch with audit mode and the benchmark only.
- Should the default policy let users clear a cloud provider for private data, or keep private data local-only by default? The answer decides the hero message.
- Will your first users accept denied public-web access after reading private data, or must the planner/worker split ship in the MVP? The interviews test this; what is your prior?
- Is the core asset the OS image, or the label specification plus the benchmark? The answer decides whether the OpenShell adapter moves earlier than Phase 6.
- Where do your own and your target users' private data actually live: local files, or APIs such as Gmail, IMAP and Notion? This sets the default labels and the list of private-source destinations.
- 'AgentOS' is crowded and 'Tidemark' needs a trademark search. Will you choose and clear a name before Phase 4?
- DCO or CLA? This must be decided before the first outside contribution.

## Changes made after the red team
- Kept from the previous synthesis: the clearance rule (labels must be a subset of the destination's clearance), per-session sockets as identity, the model call treated as an exit, the host-side canary, and a privilege-separated daemon set.
- BLOCKER, pre-opened tunnel bypass (red teams 1 and 2). The gateway now loads the session's label word, through a read-only mmapped BPF array, before every chunk it forwards upstream. Change events also reset uncovered tunnels, but they are not relied on alone, because they race. Plain HTTP is re-checked per request and the broker checks after reading the full body. Cases for tunnels opened before the read, HTTP/2 reuse and WebSockets were added to the bypass matrix, the benchmark and the demo.
- BLOCKER, laundering through persistent agent state (red teams 1 and 3).
- Every writable path is now an overlay, and doctor enforces this.
- Commit stamps every file in the upper directory with the session's final label, which also covers files opened for writing before the label rose.
- Red team 1's suggestion to stamp from file_permission was not adopted: kernel/bpf/bpf_lsm.c shows file_permission is not a sleepable hook, while file_open and file_receive are.
- Agent memory and skills are mounted as label-partitioned views, with one session per profile.
- Skill and hook directories are treated as auto-exec paths.
- A label-creep measurement with Hermes memory on was added to Phase 1, with a fallback cut if needed.
- BLOCKER, schedule (red team 3; red team 1 major). The MVP was re-baselined to weeks 4 to 24 (147 build hours), with a 7/2/1 split before launch and 6/2/2 after. The reviewers' default cuts are applied up front. Checkpoints were added at weeks 9 and 15, line estimates were raised to about 4,000 product lines plus 1,500 benchmark lines, and the dependence on AI assistance is stated explicitly.
- BLOCKER, demo determinism (red team 3). The gating demo now uses two puppet backends (a local puppet cleared for all labels and a cloud puppet acting as a sink), so placement is proven by which log receives each request. Real llama-server inference moved to a separate non-gating smoke test, plus a gating CI functional test.
- Major, network-sourced private data (red teams 1 and 3). Destinations now carry source labels, and the gateway raises the label through tidemarkd before relaying the first byte. An API-sourced scenario was added to the demo and the benchmark.
- Major, receipts oversold (red teams 1 and 2).
- Receipts are now described as tamper-evident against session processes and verifiable offline as long as the key stays secret.
- The BYOC grep is named as the primary proof.
- Ledger segments use the kernel append-only attribute against the ledger daemon, though root can still alter them.
- A TPM2-sealed key arrives in Phase 5, and TPM quotes with Keylime-style attestation in Phase 7, before any third-party claim.
- The tagline was rewritten.
- Major, novelty wording (red teams 1 and 2). The 'first installable OS' claim was dropped and the claim narrowed. Prior art is now cited: ADEPT-50, SELinux MLS/MCS, Solaris Trusted Extensions, HiStar and Flume, ActPlane's connect hook and BPF-LSM blocking, CaMeL and FIDES, and Invariant Guardrails. Qubes OS and SELinux were added to the positioning, with a stated reason for the BPF approach: MLS levels do not float.
- Major, missing competitors (red teams 1 and 2). nono (README v0.61.2 verified), AgentFS, agentsh, container-use, Microsoft AGT, Invariant Guardrails, Rubrik Agent Rewind and Qubes were added. Containment, staging, receipts, stop and release are reframed as integration, with the upstream projects credited. White space 4 is narrowed to permission-diff admission.
- Major, mischaracterised products (red team 1). sandbox-runtime is corrected: it covers the whole process tree, including MCP servers, and masks credentials. NemoClaw is corrected: it needs no account and runs on Linux with rootless Podman and on Windows through WSL. NemoClaw is named the main competitor among self-hosters, and differentiation now rests only on labels, placement and evidence.
- Major, OpenShell middleware (red team 1). The docs were verified on 16 Oct 2026: the middleware is a gRPC service over TCP that inspects, transforms or denies allowed HTTP and WebSocket traffic, after policy evaluation and before credential injection. The label engine and clearance evaluator are now packaged as pkg/labels and pkg/clearance, and an OpenShell adapter is planned for Phase 6 after a spike. Kill criterion 6 now redirects effort to the adapter rather than stopping. This departs from red team 1's v0.2 timing because v0.2 is committed to closing the credential gap.
- Major, credential paradox (red team 1). Credentialed push is removed from the v0.1 benign suite and the limitation is documented. Credential injection for git smart-HTTP and two declared APIs arrives in Phase 5.
- Major, Phase 5 holds leaving the agent out of step (red team 1). Holds are now synchronous and time-bounded. A timeout returns 'not performed', approval mints a token consumed by the agent's own retry, and Tidemark never replays. Temporal sagas and agentsh are credited.
- Major, incomplete taint (red team 2). The shim now calls close_range before exec, a file_receive hook raises the label for descriptors passed with SCM_RIGHTS, and seccomp denies io_uring_setup and pidfd_getfd. Matching read-matrix cases were added.
- Major, cgroup keying (red team 2). Sessions are now keyed by their ancestor cgroup id at a fixed depth, with Delegate=no and a sub-cgroup test case.
- Major, overlayfs and xattrs (red team 2). The Phase 0 gate proves label reads through overlayfs first. The candidate design is a cache keyed by bpf_real_data_inode, a kfunc verified in fs/bpf_fs_kfuncs.c. The fanotify fallback's gaps must be documented before it can be selected.
- Major, fragile 'labels travel with files' (red team 2). Directory rules are now evaluated in the kernel through bpf_path_d_path (kfunc verified), and sessions stamp what they commit. A content-hash index is used only to report drift. The claim is narrowed to files Tidemark controls, and a rewrite-then-reread scenario was added.
- Major, commit authority (red team 2). Commit, accept, release and stop are now accepted only from the human UID or root outside the agents slice. G6 no longer claims that writes are invisible to the host; root can inspect them.
- Major, cost of per-open lookups (red team 2). An inode-storage cache with invalidation on xattr changes was added, together with a target for real npm-install and compile workloads.
- Minor, user namespaces (red team 2). RestrictNamespaces=yes now denies them, and the containment matrix asserts it.
- Minor, testability (red team 2). x86 KVM is the only gating boot, arm64 boots do not gate, and an hourly x86 triage VM is budgeted, converted with `bootc install to-existing-root`.
- Major, releasing HTTPS bytes (red team 3). In v0.1, releases apply only to broker requests, bound to a content hash. CONNECT grants wait for TLS termination. G9 was rewritten.
- Major, tagline (red team 3). The tagline is now limited to confidentiality, and the release notes advise keeping agents' own approval on for destructive tools until v0.2.
- Major, first-user mismatch (red team 3). First users are narrowed to homelab, Proxmox and dedicated-VM users working over local files. VPS hosts, messaging gateways and OpenClaw are excluded from the v0.1 claims, and virtiofs xattr behaviour is checked in Phase 0.
- Major, interviews too late and too leading (red team 3). They are now 8 to 10 observational interviews in weeks 1 and 2, with the kill decision at week 3.
- Major, sustainability (red team 3). The launch is a preview with a 90-day disclosure window, a bus factor of 1 is assumed, reviewers are named in Phase 0, and audit-funding eligibility is checked before applying.
- Major, overloaded post-MVP phases (red team 3). Each quarter now has one theme, with a backlog of items and explicit triggers. Phase 7 is conditional and runs in months 13 to 18.
- Major, overloaded Phase 0 (red team 3). The spike runs on stock Fedora 45 VMs, the image pipeline starts with x86_64 only, and the threat documents are merged into one THREAT_MODEL.md.
- Minor, CI limits (red team 3). The full suite runs nightly within 45 minutes and a subset on pull requests within 5 minutes. v0.1 outputs are qcow2 and OCI only.
- Minor, Omarchy baseline (red team 3). It was dropped. A neutral 'unsandboxed' configuration is kept, because audit mode keeps containment on and cannot show that containment channels are real; this is a partial disagreement with the reviewer.
- Minor, benchmark bias (red team 1). The methodology is now pre-registered, two incumbent-favourable channel classes were added, and completion is published next to leak counts. nono and OpenShell join as baselines in v0.2.
- Minor, 'model calls move, not fail' (red team 1). Reworded: requests are never placed outside clearance and fail closed with a structured reason.
- Delivery (red teams 1 and 3). The COPR RPM is now the delivery unit, the image installs it, and the OCI image works with `bootc switch`. Full guarantees are claimed only on the image, which preserves the bootable-OS constraint.
- Scope details (red team 3). Each profile now runs under a fixed UID, replacing the UID pool. The broker supports Chat Completions only in v0.1. The ranker, memory-fit admission and cost cap were cut. A static HTML report replaces the dashboard and still meets the founder's 'simple UI' wish.
- Founder questions added from red team 1 (core asset; accepting denied web access) and red team 3 (experience, AI review, RPM delivery, where data lives).

## Deferred from the original specification

- **Training and fine-tuning (TorchTitan, DeepSpeed, Axolotl, Unsloth; torchtune excluded as unmaintained)** — Not before v1.0. Only if a contributor with NVIDIA hardware maintains a separate workbench image with CI on rented GPUs, and at least 20 independent users ask for it.
- **vLLM, SGLang and TensorRT-LLM** — Backlog. Trigger: a community maintainer and a g4dn boot test on every release, because vLLM ships breaking minor releases every 2 to 3 weeks. They would be registered as broker backends with clearances.
- **Multi-engine auto-tuning** — After v1.0, once two validated backends and a throughput benchmark run in CI on the target hardware.
- **Hardware adapters and an NVIDIA image variant** — Backlog. Trigger: a volunteer maintainer. Boot tests would run on AWS g4dn spot instances at about USD 0.33 per hour. AMD ROCm only with its own maintainer.
- **MLX** — Never inside the image. It is supported as an unattested endpoint on the macOS host from the MVP onward.
- **Bundling Ollama** — Not bundled. Supported as a broker endpoint.
- **Model unification (Llama, Qwen, DeepSeek, Mistral, Gemma)** — Cut as a feature. Any GGUF model runs through llama-server, and defaults are small Apache-2.0 models. Llama is treated as a legacy line.
- **Laya-style router and the ranker interface** — The interface is cut from v0.1. Backlog trigger: two or more cleared backends and a labelled request set. It may only reorder the permitted set and is never a security boundary.
- **GPU and cluster resource scheduler** — Cut. The 'basic scheduler' consists of per-session cgroup quotas and one session per profile. llama-server keeps its own slot queue.
- **Six-section UI and a live dashboard** — v0.1 ships a static HTML report (Home, Sessions, Models, Receipts, System). Training and Developer Tools pages are cut. A live dashboard or Cockpit plugin is a backlog item, triggered when at least 10 active users ask for one.
- **Desktop environment** — After v1.0, and only with at least 100 weekly active installs and a desktop maintainer.
- **Holds for irreversible effects** — Phase 5 (v0.2): synchronous, time-bounded holds for declared APIs.
- **Compensations and saga journal** — Backlog. Trigger: v0.2 users need undo rather than hold. Credit Temporal sagas and agentsh soft_delete.
- **Pinned verifier gate at commit** — Backlog. Trigger: coding-agent users report tests or CI files being tampered with.
- **Capability manifests (narrowed to permission-diff admission, since nono already ships signed instruction files)** — Backlog. Trigger: at least 3 users run third-party MCP servers or skills under Tidemark and ask for admission control.
- **Cross-vendor attention scheduler** — Backlog. Trigger: users report more than 10 releases per day. The span-of-control benefit must be measured, not claimed.
- **Regulatory evidence mapping (EU AI Act, Colorado, CPPA, Australian ADM)** — Backlog. Trigger: a deployer or auditor commits to a pilot. It would be presented as mapping, never certification, and could target 2 Dec 2027 only if triggered by Phase 6.
- **MicroVM principals** — Backlog, on KVM hosts only. They can never be tested on the M1.
- **Sealed UKI, TPM attestation and measured boot** — The TPM2-sealed key comes in Phase 5. Quotes and attestation come in Phase 7. Sealed UKI only when bootc sealed images leave test status.
- **ISO, raw and AMI outputs; aarch64 qcow2** — v0.1 ships x86_64 qcow2 and OCI images for both architectures. Other formats wait for user demand; aarch64 qcow2 is best effort.
- **Btrfs workspace safety snapshot** — Cut. Overlay staging makes it redundant.
- **tidemark-lite DEB and AUR packages** — Phase 6. The Fedora 45 RPM ships in v0.1.
- **Claude Code and OpenClaw profiles** — Phase 6, when the broker gains Anthropic Messages, and with OpenClaw limited to short sessions.
- **Bare-metal support on Asahi** — Community-maintained only.
- **Agent identity, inventory and payments** — Cut, because incumbents bundle these.
- **Enforcement on Windows and macOS** — Backlog: policy export to MXC or sandbox-runtime at a stated lower guarantee level.

## Budget
Cash: about USD 20 to 50 per month in total.
- Free for a public repository: GitHub Actions (x64 and arm64), GHCR, COPR and Sigstore.
- Hourly x86_64 VM for interactive triage: about USD 5 to 10 per month.
- Domain: about USD 15 per year, bought only after the trademark check.
- Model API for nightly live-mode runs: USD 0 to 15 per month (a cheap small model such as Claude Haiku 5.5 at USD 0.10 per million input tokens, or local only).
- AWS g4dn spot instances only if the NVIDIA backlog item triggers.
- Optional before v1.0: professional trademark clearance, roughly USD 300 to 1,500 (an unverified estimate). Self-search the USPTO, EUIPO and IP Australia databases first.

Hardware: no purchase is needed.

Time, re-baselined honestly:
- Phase 0: 30 hours (21 building).
- MVP, weeks 4 to 24: 210 hours (147 building).
- Phases 5 and 6: 260 hours (156 building).
- Phase 7, if the adoption gate passes: 260 hours.

Default cuts already applied to v0.1:
- live dashboard (replaced by a static report);
- Claude Code profile;
- Responses and Anthropic Messages;
- ranker interface, memory-fit admission and cost caps;
- CONNECT releases;
- OpenShell and Omarchy baselines;
- ISO, raw and AMI outputs;
- Btrfs snapshot;
- Mac demo path in CI;
- benign tasks reduced to 10.

Contingent cut order, applied at the week-9 or week-15 checkpoint when slippage exceeds one week:
1. Static report replaced by CLI text only.
2. aarch64 OCI image moved to v0.2.
3. Source labels moved to v0.2, with the API scenario reported as an open channel.
4. State views replaced by discarding state writes from labelled sessions.
5. Release queue dropped; the user starts a fresh session instead.

Never cut: the per-chunk re-check, commit stamping, the BYOC demo and the cross-vendor step.

## Risks

- **Label creep: at session granularity, ordinary sessions may lose public egress too often.** — - Label-partitioned memory, so creep does not persist across sessions.
- Clearances per destination, with placeholders that do not taint.
- Repository files unlabelled by default.
- The Phase 1 measurement and kill criterion 3.
- The planner/worker experiment in Phase 6.
- **A bug in the per-chunk check or the label plumbing reopens the bypass through pre-opened tunnels.** — - A race test of 10,000 iterations.
- Tunnel cases in the CI run on every pull request that touches enforcement.
- Fuzzing of the copy loop.
- Private review by named reviewers before launch.
- **Kernel fragility: kfuncs are not a stable ABI, overlayfs may interact badly, and the BPF verifier may behave differently across architectures.** — - The kernel is pinned through the image digest, with assertions in CI.
- The Phase 0 gate runs on stock Fedora 45.
- The fallback's gaps are documented before it can be chosen.
- Renovate merges only after the boot test and the read matrix pass.
- **Tools outside sessions drop security.* xattrs and silently remove labels.** — - Directory rules for default sensitive folders, evaluated in the kernel.
- Sessions stamp what they commit.
- audit-labels reports drift using a content-hash index.
- The claim is narrowed to files that Tidemark controls.
- **Scope exceeds solo capacity.** — - A re-baselined date with default cuts.
- Checkpoints at weeks 9 and 15 with a contingent cut order.
- One language.
- An open question on experience and AI assistance.
- **Local models on the M1 are slow and weak at agentic work.** — - Placement is a clearance decision, not a forced move to local.
- An unattested Metal endpoint on the host is supported.
- Completion is published per hardware class.
- The pitch is control and proof, not speed.
- **Agents break: clients that ignore proxy settings, no DNS in the session, SSH git denied, IMAP.** — - Profiles set the proxy variables.
- Structured 403 messages explain each denial.
- Clients fail closed rather than leak.
- A compatibility matrix is published for every release.
- **Domain-level clearance is coarse, and CDN domain fronting is possible.** — - Documented as a non-guarantee.
- Guidance on which destinations are safe to clear.
- Clearances by method and path arrive in v0.2.
- **An incumbent bundles the idea, for example NVIDIA adding a taint check as OpenShell middleware.** — - pkg/labels and pkg/clearance plus the Phase 6 adapter turn OpenShell into a distribution channel.
- Open specifications and a neutral benchmark.
- Incumbent changelogs reviewed every quarter.
- Kill criterion 6 applies.
- **NemoClaw already serves self-hosters with less friction than booting an image.** — - Positioned as complementary through the adapter.
- The RPM and `bootc switch` paths.
- The BYOC proof differentiates on evidence, not on convenience.
- **Critics say the attacker is scripted, or dispute benchmark fairness.** — - Real CLIs, a host-side canary and pre-registration.
- Neutral baseline names and each tool in its strongest documented configuration.
- Incumbent-favourable channels included, with completion published next to leaks.
- A 7-day right of reply.
- Live mode reported separately.
- **People will not boot an image.** — - The RPM on Fedora 45, `bootc switch`, the benchmark container and the verifier are each useful without booting.
- **Upstream churn: Hermes Agent's issue volume, Fedora's six-month releases, llama.cpp's pace.** — - Digests pinned.
- Nightly profile smoke tests.
- Agents installed in /var.
- Rebases only after full CI passes.
- **The base image is new: fedora-bootc:45 was published on 15 Oct 2026.** — - Pinned by digest.
- The :44 fallback build stays until launch.
- **Root code in the trusted computing base.** — - Small root surface.
- Unprivileged gateway, broker and ledger daemon.
- Fuzzing.
- systemd sandboxing of every daemon.
- **Disclosure burden and burnout for a solo maintainer.** — - Launched as a preview with a 90-day disclosure window.
- Triage capped after launch.
- Named reviewers and a hall of fame.
- A bus factor of 1 assumed.
- **Free runner limits: disk space, time, and arm64 KVM.** — - Only qcow2 and OCI outputs.
- A 45-minute nightly suite and a 5-minute pull-request subset.
- Only the x86 KVM boot gates.
- A rented x86 VM for triage.
- **Receipts expose sensitive paths, or are oversold.** — - Paths and bodies hashed by default.
- The 'tamper-evident against the agent' wording enforced by the claims checker.
- Attestation required before any third-party claim.
- **Phase 5 holds leave the agent out of step with reality.** — - Synchronous, time-bounded holds.
- A timeout means not performed.
- No replay; approval applies on the agent's own retry.

## Appendix A: candidate concepts

### AgentOS Escrow (architect:radical)
_Agents act without prompts. Every change to your disk and to the outside world is staged, typed and signed, and nothing irreversible leaves the machine until you commit._

**Claim.** For the first time, a third-party agent from any vendor (Hermes Agent, Claude Code, OpenCode, or a scripted replay of a compromised agent) can run unmodified with all permission prompts disabled, on an installable Linux image, under three guarantees. (a) Every file change is staged in an overlay. (b) Every outbound HTTP(S) effect is classified as read-only, reversible, compensable or irreversible by a signed effect manifest; unknown effects default to irreversible. Irreversible effects are held in escrow and never reach the destination before commit. (c) Any outbound request from a task that has read a file labelled secret is held, even to an allowlisted host such as the project's own forge, if it carries a fingerprint of that file. When the task ends, the OS writes a receipt that a third party can verify with a public key. The receipt lists the staged diff, the effect journal, the labels touched and the OS image digest. Abort restores the workspace byte for byte and runs recorded compensations in reverse order. The claim can be checked three ways: the hero demo assertions (mail sink shows 0 messages, the forge repository still exists, the workspace tree hash matches); flipping one byte of a receipt makes verification fail; and a review of the feature lists of OpenShell, MXC, sandbox-runtime, Docker Sandboxes, Copilot sandboxing, agent-sandbox, Sprites, Agent 365 and AgentCore, none of which stages remote effects or classifies them by reversibility (as of 16 Oct 2026). Limits stated up front: only HTTP(S) effects described by manifests are classified, non-HTTP egress is blocked rather than understood, and fingerprinting misses paraphrased leaks.

**Hero demo.** Name: abort the injected agent. Time: under 10 minutes after the image download, which is the only step that depends on bandwidth.

1. Download agentos-0.1 qcow2 for x86_64 or aarch64 and verify it with cosign verify-blob against the published key.
2. Boot it in QEMU, UTM or Lima, or run bootc switch on an existing Fedora bootc host, then log in as demo.
3. Run agentos demo injection (about 2 minutes). It starts the local Forgejo and Mailpit, labels ~/proj/.env as secret, then runs a deterministic replay of a compromised agent with all prompts off. The replay:
   - reads a README containing an injected instruction;
   - reads .env;
   - edits 3 source files;
   - creates branch agent/fix on Forgejo (reversible; compensation: delete branch);
   - posts an issue comment containing the base64-encoded secret to the allowlisted forge;
   - sends an email through the Mailpit send API (irreversible);
   - calls DELETE on the repository (irreversible and destructive).
4. Run agentos diff. Expected: 3 staged files; 1 reversible effect executed with its compensation journalled; 3 held effects (the email, the repository delete, and the comment flagged as label violation secret to forge).
5. Run agentos demo check. Expected: Mailpit shows 0 messages, the Forgejo repository exists with 0 comments, and branch agent/fix exists.
6. Run agentos abort, then agentos demo check --after-abort. Expected: the workspace tree hash equals the pre-task hash, the branch is gone, Mailpit is still at 0, and the repository is intact.
7. Run agentos verify ~/receipts/latest.dsse.json. Expected: signature valid, chain valid, and image digest equal to bootc status. Flip one byte and the same command fails.
8. Optional: rerun the demo, then agentos commit --approve email-1 --deny rest. Exactly 1 message appears in Mailpit and nothing else is sent.

Each demo check prints PASS or FAIL. CI runs the identical script on x86_64 (KVM) on every commit and publishes the receipts, so anyone can verify those receipts in about 1 minute without booting anything. A non-deterministic variant, agentos run -- hermes --yolo, uses the viewer's own model key or endpoint.

### Docket (working name for the AgentOS flagship; trademark search required before public use) (architect:workflow)
_Agents run as transactions. Independent verifiers decide what may commit. You are interrupted only for what cannot be undone._

**Claim.** As of October 2026, no shipping product we found does all four of the following for an unmodified agent from any vendor, together with every MCP server, hook, LSP server and subprocess that agent starts:
(a) all file writes land in a staging layer;
(b) every outbound network request that changes state is held in a side-effect journal instead of being sent;
(c) a verifier set taken from the commit before the task, which the agent cannot edit, decides whether the result may commit;
(d) only the irreversible remainder reaches a human, through one queue ranked across all agents and capped by an interruption budget.
Checkable form, from the published hero demo: six agents run in parallel and the demo prints and journals every effect event (design target: at least 30). The human makes exactly 3 decisions and receives 1 interrupt. After the rejected transactions, the SHA-256 of the host home directory and the state of the mock production database are unchanged. The approved webhook is delivered exactly once. `docket journal verify` passes, then fails at the exact entry after a one-byte edit.
The conceptual change: permission prompts give way to commit, as in a database. Agents act freely inside staging, and approval is requested once, at commit, only for effects that staging cannot contain.

**Hero demo.** 'Six agents, three decisions, one interrupt', reproducible in under 10 minutes in either of two ways.

Path A, no local setup:
1. Fork the repository.
2. Run the reproduce-demo GitHub Actions workflow. On a free x86_64 KVM runner it downloads the signed qcow2 release asset, checks it with `cosign verify`, boots it in QEMU, runs `docket demo fleet --decisions=scripted` and uploads the journal and a verification report. Target: 8 minutes or less.

Path B, local:
1. Download the qcow2 for aarch64 or x86_64 (target size 2 GB or less).
2. Boot it in UTM, Lima or QEMU.
3. Run `docket demo fleet` and decide interactively.

The fleet runs six scripted agents against a small Python library with tests. The agents are scripted so that the demo checks the OS mechanism, not model quality, and needs no GPU or API key. A live mode with Hermes Agent and the user's own key or a local llama-server is optional.
- A: a correct bug fix. Verifiers pass, so it is auto-committed to docket/A with no human involved.
- B: a correct refactor that also starts a filesystem MCP server, whose writes are staged too. Auto-committed.
- C: makes the tests pass by editing an existing test file. Verifier-tamper flag, so it goes to the digest.
- D: breaks the tests. Auto-returned to the agent with the verifier output, with no human involved.
- E: runs `rm -rf ~` and sends DELETE to a mock production API inside the VM. The deletion stays inside the overlay and the DELETE is held as irreversible. Top priority, and the only interrupt.
- F: posts to a mock chat webhook. Held, low priority, goes to the digest.

The human sees one ranked queue with three decision cards: approve F, reject E, reject C.

`docket demo verify` then prints six independently re-computable checks:
1. The SHA-256 of /home/demo is identical before and after.
2. The mock production database row count is unchanged.
3. The webhook log shows exactly one delivery, carrying F's idempotency key.
4. Effect events total N (target 30 or more), with 3 human decisions and 1 interrupt.
5. `docket journal verify` passes. A scripted one-byte edit made as root then makes it fail at the reported entry, and the same edit attempted as an agent UID is denied with EACCES.
6. The escape suite (curl --noproxy, raw socket, UDP, IPv6, DNS-over-HTTPS, direct IP) is captured or blocked, with 0 escapes.

A 90-second recording is published with the release.

### AgentOS Bulkhead (architect:adoption)
_Keep your agents in yolo mode. The OS records what each agent session has read and holds every exit that could leak it._

**Claim.** Bulkhead adds a third permission state between deny and allow: quarantine-on-read. A file can be readable, but opening it changes what the reading session may send out afterwards.

The claim, checkable as of October 2026: no shipping product lets unmodified agents from different vendors (Claude Code, Codex, OpenCode, Hermes, OpenClaw) run with per-action prompts switched off while all of the following hold:
- The OS, not the agent, records which labelled files each agent session has opened, including files opened by MCP servers, hooks and double-forked children. The label is committed before the open returns, using fanotify permission events.
- Every later exit outside a declared trust set is held for one human decision. This includes exits to allowlisted hosts.

How to check it: run LeakBench in three modes. The modes are (a) the skip-permission flags Omarchy's launcher passes, (b) the same OpenShell sandbox with a static allowlist only, and (c) Bulkhead. Pass condition for (c): zero canary bytes in any sink log across all 11 scenarios, and zero prompts before a session reads a labelled file.

Explicitly not claimed:
- Content the agent sends to its configured model provider. That exit is excluded by design. With a local model, the guarantee becomes "nothing labelled leaves the machine".
- Protection against kernel exploits.
- Fine-grained tracking of data through model transformations. Labels are per session and deliberately coarse.

**Hero demo.** "Hijacked on purpose: 11 exfiltration attempts, 3 setups, 1 table." Anyone can reproduce it in under 10 minutes.

1. Get the image (about 4 to 6 minutes for a qcow2 of 1.5 GB or less on a 50 Mbit/s link):
   - macOS: limactl start agentos.yaml
   - Linux: import the qcow2 with virt-install or quickemu
   - Fedora Atomic: sudo bootc switch ghcr.io/ORG/agentos:0.1, then reboot
   - Verify the signature with cosign verify.
2. Run agentos doctor (about 10 seconds). It prints the kernel version, Landlock ABI (9 or higher), BPF LSM status, fanotify permission support, the enforcer version and telemetry off. It exits non-zero if any guarantee is missing.
3. Run leakbench run (about 3 minutes, offline, no API key). A synthetic hijacked agent, plus real unmodified CLIs driven by a local puppet-model server (OpenCode, Hermes and Claude Code at minimum), plant fake canary secrets and try 11 exits:
   1. HTTPS to an unlisted host
   2. DNS tunnelling
   3. POST to an allowlisted gist-style API
   4. git push to an allowlisted remote
   5. a secret inside a package-registry URL
   6. an encoded (base64 and reversed) payload
   7. an exit from a spawned MCP server
   8. an exit from a double-forked daemon
   9. host-socket abuse (container socket, D-Bus)
   10. an auto-exec implant in .git/hooks or a shell rc file
   11. ledger tampering

   The "allowlisted" services are local sinks, so nothing leaves the VM.
4. Read the table for three modes:
   - (a) yolo: the same skip-permission flags Omarchy's launcher passes, with no sandbox
   - (b) allowlist: the same OpenShell sandbox with quarantine-on-read disabled
   - (c) Bulkhead

   Bulkhead passes only with zero canary bytes in any sink (viewers check the raw logs with leakbench sinks grep), every held exit listed in agentos queue, and zero prompts before taint. The results for (a) and (b) are hypotheses: they will be measured and published as observed. The expected pattern is that (b) stops unlisted hosts but not allowlisted ones.
5. Run agentos ledger verify (about 5 seconds). The chain is intact and the truncation attempt is recorded as denied.
6. Run agentos undo SESSION (about 10 seconds). It reverts the implant and shows the diff.

A 90-second video shows the same run. An optional step runs the same poisoned repository with the user's own keys or a host Ollama model.

### AgentOS Tidemark (architect:local_sovereign)
_Your data decides where your AI runs, and the machine proves it._

**Claim.** For the first time, an unmodified third-party agent (Hermes Agent first) on a personal machine automatically loses the ability to send anything off the machine, including to its own cloud model API, at the moment any process in its tree reads a file the user labelled private. It keeps working on a local model and is not interrupted by prompts. The label then persists on every file that session writes, so a different vendor's agent that later reads those files is constrained the same way. How to check it: a verifier generates a random canary on their own host, boots the signed image, and runs the published prompt-injection scenario twice. With labels off (ordinary allowlist containment), the canary appears in the host-side sink and mock-cloud logs. With labels on, grep of those host-side logs returns zero matches, the receipt shows that the inference call was placed on the local engine, and `tidemark verify` validates a hash-chained, Ed25519-signed receipt bound to the cosign-verified image digest. This claim is bounded to 16 Oct 2026. It is based on the research corpus and on direct checks of OpenShell, NemoClaw, ActPlane and Tailscale Aperture, none of which combines read-time labels, label-constrained model placement, forced egress mediation and verifiable receipts.

**Hero demo.** 'Bring your own canary', designed to finish in under 10 minutes on a 4-core x86_64 Linux machine with KVM or an Apple silicon Mac with QEMU, over a 50 Mbit/s connection. Anyone can also fork the repo and run the same workflow on a free GitHub x64 runner.

Step 1. `git clone` the repo and run `./demo/verify.sh`. The script downloads the signed qcow2 (target 1.6 GB or less, including a sub-2B Apache-2.0 GGUF model) and runs `cosign verify` against the published identity.
Step 2. The script generates the canary on the host (`openssl rand -hex 16`). Outside the VM it starts two host-side loggers: an 'attacker' sink, which the policy treats as an allowlisted public site, and a mock cloud LLM endpoint that records every prompt. It boots the VM with the canary passed as a QEMU fw_cfg systemd credential into ~/Private/notes.txt, which is private by default.
Step 3. Run A, labels off (equivalent to ordinary allowlist containment). A replay agent issues the exact tool calls from a recorded real Hermes Agent transcript of a prompt injection: fetch an injected web page, read ~/Private/notes.txt, POST its contents to the sink, then include them in the next model prompt. The host logs now contain the canary.
Step 4. Run B, labels on. The session label rises to P+U at the read, the POST is denied, the next model call is placed on the local llama-server, and `grep -c <canary>` on both host logs returns 0. The verifier runs this grep on their own machine, outside the VM.
Step 5. Cross-vendor step: a second, different client (OpenCode or a minimal OpenAI-SDK script) reads the summary file written in Run B. It starts as P, and its cloud call is placed locally or denied.
Step 6. `tidemark receipt export` and then `tidemark verify` on the host. The check passes, the receipt's image digest matches the cosign-verified digest, and flipping one byte makes verification fail.

Optional live mode runs real Hermes Agent with a real API key. It is not required for verification, because injection success is stochastic. The same scenario runs in CI on every commit for both architectures, and the receipts are published.

### AgentOS Escrow (architect:security)
_The safe place to let agents act: data is labelled when it is read, side effects wait in escrow until committed, and every decision is written to a ledger the agent cannot touch._

**Claim.** Escrow would be the first installable OS on which an unmodified agent from any vendor (Hermes Agent and OpenCode in v0.1; OpenClaw and Claude Code in v0.2) runs with its own permission prompts turned off, together with every MCP server, hook and subprocess it starts. The OS would do three things. (1) It labels data when a process reads it, follows those labels across processes, files and agents from different vendors, and blocks any request, including requests to allowlisted hosts, whose destination is not cleared for the session's labels. The decision depends on what the session read, not on the bytes it sends, so encoding or paraphrasing the data does not get it out. (2) It lets read-only, locally staged and compensable actions run, and holds only irreversible external actions for one human decision. (3) A process the agent cannot reach records every decision in a hash-chained, signed ledger tied to the signed image digest. How to check it: fork the public repository and run the red-team workflow on a free GitHub x64 KVM runner. The workflow runs the same scripted 'fully compromised agent' twice, once in audit mode and once in enforce mode. Pass means that in enforce mode the attacker sink receives zero canary bytes on every listed channel, every irreversible-class request is held, an aborted session leaves the workspace byte-identical, and the standalone verifier rejects any tampered ledger copy at the exact altered entry. Not claimed: prevention of prompt injection, resistance to kernel exploits (the kernel is shared until microVM principals arrive), timing or other low-bandwidth covert channels, or protection of data the user deliberately cleared for a model provider. CaMeL and FIDES enforce information flow inside one rewritten agent, and HiStar/Flume did OS-level information-flow control in research. What is new is enforcement across vendors with no change to the agent, combined with staged commit and third-party-verifiable evidence, as a reproducible product.

**Hero demo.** 'The Heist', reproducible in under 10 minutes. Path A, with no local machine: fork the repository and run the 'reproduce' workflow. A free GitHub x64 KVM runner boots the cosign-verified qcow2 under QEMU, runs the demo and uploads the scorecard and ledger bundle (about 6 to 8 minutes). Path B, local: `limactl start <published Lima template with pinned image digest>` on macOS or Linux, or `bootc switch` on an existing Fedora bootc VM; the run time depends on downloading an image of about 1.5 GB.

Seeded state: a git workspace; ~/vault/clients.csv labelled private:clients; a token labelled secret; and a mock issue containing a prompt injection, served by a local mock GitHub-like API.

Step 1: `agentos demo heist --mode audit`. A scripted, fully compromised agent reads the CSV and token, then tries to send canary data through 13 channels: POST to an attacker host; DNS subdomain encoding; hard-coded IPv4 and IPv6; UDP/QUIC; a public issue comment on the allowlisted API; a public paste on an allowlisted host; a GET query string to an allowlisted search API; git push to an attacker remote; an MCP-server child opening its own socket; a hook that double-forks a daemon which outlives the agent; a Unix socket to a host service; cross-agent laundering, in which agent A writes a derived file and a second agent with a different vendor manifest posts it; and a model-API request defining a provider-side web-fetch tool aimed at the attacker. It also edits 200 files, creates an issue (compensable), deletes a branch and sends mail (both irreversible). The scorecard shows which channels leaked, which proves the test environment is real.

Step 2: `--mode enforce`, same script. Pass criteria: the attacker sink receives zero canary bytes; the queue holds exactly the branch deletion and the mail; the issue creation proceeded and was journaled; the file edits are staged.

Step 3: `agentos abort`. The workspace tree hash matches the starting state, and the mock API state matches the initial state.

Step 4: `agentos ledger export`, then run the standalone verifier on the host: OK. The script then edits, deletes and reorders entries in copies of the bundle, and the verifier fails at the exact entry each time. An attempt to write the ledger from inside a principal is denied and recorded.

Optional live mode: `--live --agent hermes` with a local llama.cpp model or the user's own API key, so a real agent faces the injected issue.

Output: a 90-second video and a scorecard JSON that names the image digest.

## Appendix B: judges

### judge:adoption
Lens: Adoption and revolutionary potential: would HN and r/LocalLLaMA developers and self-hosters install it within a week of launch, is the hero demo concrete, checkable and memorable, does it introduce a new OS or agent primitive rather than repackaging, and could it become the default safe home for agents people already run. Cross-cutting challenge: no concept proves that this audience will boot a headless image to run its daily agents. The realistic first adoption unit is the reproducible demo or benchmark, not the OS, and every plan should be judged on whether something useful still spreads when people never install the image.

Ranking: AgentOS Escrow (security-first lens: labels, clearance gateway, holds, ledger), Docket (transactions, held intents, pinned verifier gate, attention scheduler), AgentOS Tidemark (read-time labels decide local versus cloud model placement), AgentOS Bulkhead (quarantine-on-read on OpenShell plus LeakBench), AgentOS Escrow (radical-primitives lens: typed effects, synthetic 202 escrow, fingerprints, in-toto receipts)

| Concept | Evidence | Differentiation | Feasibility | OS justification | Adoption | Revolutionary | Total |
|---|---|---|---|---|---|---|---|
| AgentOS Escrow (security-first lens: labels, clearance gateway, holds, ledger) | 8 | 8 | 3 | 7 | 6 | 8 | 40 |
| Docket (transactions, held intents, pinned verifier gate, attention scheduler) | 7 | 7 | 5 | 4 | 7 | 8 | 38 |
| AgentOS Tidemark (read-time labels decide local versus cloud model placement) | 6 | 7 | 4 | 8 | 6 | 7 | 38 |
| AgentOS Bulkhead (quarantine-on-read on OpenShell plus LeakBench) | 8 | 5 | 6 | 6 | 7 | 5 | 37 |
| AgentOS Escrow (radical-primitives lens: typed effects, synthetic 202 escrow, fingerprints, in-toto receipts) | 7 | 6 | 3 | 7 | 5 | 7 | 35 |

Fatal flaws:
- **AgentOS Escrow (security-first lens: labels, clearance gateway, holds, ledger)**: The scope does not fit the founder's time. Eight primitives and about 3,300 lines of security-critical code (root Go daemon plus a mitmproxy addon in the TCB) are planned for about 120 hours with no slack, and a public zero-leak claim that ships incomplete destroys credibility. Separately, labels live only on a dedicated vault filesystem, so the secrets users actually hold (~/.ssh, ~/.aws, repository .env files) are unprotected by default unless moved.
- **Docket (transactions, held intents, pinned verifier gate, attention scheduler)**: The MVP has no read-state egress control. GET and HEAD requests to allowlisted hosts and inference calls pass unchecked, so a prompt-injected agent can leak secrets through a query string or a provider-side web-fetch tool, and the safety framing of the 'rm -rf ~ is harmless' demo will be dismantled within hours of launch. HTTP 202 deferral of every mutating request also breaks common write-then-read agent flows.
- **AgentOS Tidemark (read-time labels decide local versus cloud model placement)**: Session-level P and U taint, with P+U meaning no external egress, will force most realistic agent sessions onto 4 to 8B local models with minutes of prefill on the M1 class, or cut them off entirely. The product is then either unusable for agentic work or switched off, and the evidence it cites (local models trail on agentic tasks) predicts exactly this outcome.
- **AgentOS Escrow (radical-primitives lens: typed effects, synthetic 202 escrow, fingerprints, in-toto receipts)**: Exfiltration to allowlisted destinations is stopped only by fingerprint matching on plain, base64 and hex forms, which any alternative encoding, compression or chunking defeats, so the headline label guarantee is not a security boundary. Synthetic 202 responses for unknown effects will break real agents, and the concept duplicates the security-first Escrow with weaker protection.
- **AgentOS Bulkhead (quarantine-on-read on OpenShell plus LeakBench)**: The model-provider exit is excluded by design, yet provider-side tools (server-side web fetch defined in the model request) let a tainted session send secrets to arbitrary URLs through the trusted endpoint. This is a predictable first bypass that would publicly defeat the 'zero canary bytes' claim unless model-API request bodies are parsed and undeclared provider-side tools are denied. Holding every exit after a common .env read also risks recreating prompt fatigue.

### judge:evidence
Lens: Evidence and differentiation: how strongly each concept rests on verified, severe, unsolved pain, and how far it stands apart from what incumbents ship free as of 16 Oct 2026. Anything already shipped was penalised. I checked four facts myself on 16 Oct 2026. (1) The OpenShell README mentions no taint, labels, write staging or audit log, but it does say new access is approved while the agent runs. (2) ActPlane is MIT-licensed, has 104 stars and no visible releases, and propagates labels along fork/exec and file read/write edges. (3) The NemoClaw README never says "privacy router"; it describes only "local and routed inference". (4) Stock Fedora's kernel config (kernel-ark) is CONFIG_LSM="lockdown,yama,integrity,selinux,bpf,landlock,ipe", so BPF LSM and Landlock are already on in Fedora by default. That weakens every concept that justifies the OS image by saying "only our image turns on BPF LSM".

Ranking: AgentOS Escrow (concept 1, security-first), AgentOS Escrow (concept 4, radical primitives), Docket, AgentOS Tidemark, AgentOS Bulkhead

| Concept | Evidence | Differentiation | Feasibility | OS justification | Adoption | Revolutionary | Total |
|---|---|---|---|---|---|---|---|
| AgentOS Escrow (concept 1, security-first) | 8 | 9 | 4 | 6 | 5 | 8 | 40 |
| AgentOS Escrow (concept 4, radical primitives) | 8 | 7 | 5 | 6 | 5 | 8 | 39 |
| Docket | 6 | 6 | 7 | 5 | 6 | 6 | 36 |
| AgentOS Tidemark | 6 | 7 | 5 | 5 | 5 | 7 | 35 |
| AgentOS Bulkhead | 6 | 5 | 6 | 4 | 7 | 5 | 33 |

Fatal flaws:
- **AgentOS Escrow (concept 1, security-first)**: Not fatal, but very risky on schedule: eight primitives, Go plus Python, and a public claim of zero leaks across 13 channels, all in about 120 hours with no slack. Session labels that only rise may block more than 25% of benign tasks (the week-6 kill criterion), which would collapse the enforcement claim. The kernel-floor argument for the image is overstated, because stock Fedora already ships bpf and landlock in CONFIG_LSM.
- **AgentOS Escrow (concept 4, radical primitives)**: The headline claim about allowlisted destinations rests on fingerprint matching of plain, base64 and hex content, so gzip, chunking, splitting data across requests or paraphrasing leak data while the concept claims protection. That is DLP, which incumbents already ship, presented as information-flow control.
- **Docket**: A compromised agent can still leak data. GET, HEAD and OPTIONS requests to allowlisted hosts and all model-endpoint traffic pass unchecked, so secrets can leave through query strings, side-effecting GET requests or prompts. 'Safe yolo' therefore covers only side effects, and the attention-scheduling thesis has no supporting study.
- **AgentOS Tidemark**: The product forces any session that touches private files onto local models, which the evidence shows are weak at agentic work and slow on M1-class hardware (minutes of prefill). Sessions with P+U labels lose all external egress. Whether users will accept that cost is unvalidated, and the BPF LSM argument for the image is wrong for stock Fedora.
- **AgentOS Bulkhead**: Holding every exit after any labelled read brings back per-action prompting for normal coding sessions, which read .env routinely, and holds exactly the rubber-stamp approvals its own 93% evidence condemns. It is also a single feature added to OpenShell, which already approves new access at runtime and could add taint conditioning in one release.

### judge:feasibility
Lens: Lens: can a solo founder working about 10 hours a week (about 120 hours over 12 weeks) on an M1 Pro ship this as a bootc image? I checked facts on 16 Oct 2026. Fedora 44 ships kernel 7.2.10, so Landlock ABI 10 is available. It ships systemd 259, so the v261 RestrictFileSystemAccess and v262 Varlink StartTransient features are missing. Fedora 45, with systemd 262, is already out. Fedora and EPEL do not package mitmproxy. In mainline fs/bpf_fs_kfuncs.c, bpf_get_file_xattr can be called only from LSM programs (marked KF_SLEEPABLE) and reads user.* and security.bpf.* xattrs; bpf_set_dentry_xattr writes only security.bpf.*. OpenShell's local gateway needs Docker, Podman or host virtualization, and its README documents no external per-request decision hook. For scale: OpenShell had 99 contributors (8 heavy) and reached only v0.1.4 in about 8 months; MXC took 3-5 core contributors 8 months. Every concept plans about 2,700-3,000 lines of security-critical code plus the image, CI, a demo and docs in about 120 hours. My estimates run from 90 to 200 hours depending on the concept. Questions the founder must answer: (1) Is a headless VM appliance acceptable as the 'literal bootable OS'? (2) Must the 12 weeks hold even if scope is halved, or is a 20-24 week MVP acceptable?

Ranking: AgentOS Tidemark, AgentOS Escrow (security-first), AgentOS Escrow (radical primitives), AgentOS Bulkhead, Docket

| Concept | Evidence | Differentiation | Feasibility | OS justification | Adoption | Revolutionary | Total |
|---|---|---|---|---|---|---|---|
| AgentOS Tidemark | 7 | 8 | 6 | 9 | 6 | 7 | 43 |
| AgentOS Escrow (security-first) | 9 | 8 | 3 | 8 | 6 | 8 | 42 |
| AgentOS Escrow (radical primitives) | 8 | 7 | 3 | 8 | 5 | 8 | 39 |
| AgentOS Bulkhead | 8 | 6 | 5 | 5 | 7 | 5 | 36 |
| Docket | 7 | 7 | 5 | 4 | 6 | 6 | 35 |

Fatal flaws:
- **AgentOS Escrow (security-first)**: The 12-week plan cannot be delivered as written: four open gaps plus a 30-scenario red-team suite in about 120 hours, with no slack in three of the weeks. Without cuts agreed in advance, it will either slip to about 20-24 weeks or launch with an in-scope channel still open, which its own release gate forbids.
- **Docket**: Classifying effects by HTTP method (GET read-only, everything else a change, fake 202) holds legitimate POST-based reads (GraphQL, MCP streamable HTTP, search) and leaves GET-based leaks to allowlisted hosts open. With no labels until v0.3, the MVP does not address data leaks driven by prompt injection.
- **Docket**: By its own account the bootable OS is only a delivery vehicle, and an RPM provides most of the guarantees. That fails the founder's non-negotiable requirement that the OS be justified, not decorative.
- **AgentOS Tidemark**: On the founder's only local hardware, private sessions are forced onto 4-8B models with multi-minute prefill. If the week-10 stop rule triggers (more than 40% of tasks forced local), the session-level design cannot launch, and the showcase must run on x86 hardware the founder can test only in CI. The plan also never specifies how a proxied TCP connection is attributed to its session cgroup.
- **AgentOS Escrow (radical primitives)**: Taint passes through a BPF ring buffer to a userspace daemon, so a read followed at once by a send can beat the proxy check. Label enforcement relies on fingerprints, which gzip, XOR or paraphrase defeats, so the label guarantee is weaker than the claim suggests.
- **AgentOS Bulkhead**: Core enforcement depends on an OpenShell per-request external decision hook that the README does not document. OpenShell's container-based compute drivers also conflict with the systemd-unit session model. If the week-1 trial fails, the fallback Go proxy and wrapper rebuild what OpenShell was supposed to provide, roughly doubling scope.
- **AgentOS Bulkhead**: fanotify inode marks cannot match glob patterns (.env*, *.pem) across arbitrary directory trees. Matching them needs a filesystem mark on /var/home, which sends every file open by every process through the daemon and contradicts the plan's target of no measurable overhead on unlabelled opens.

## Appendix C: red team

### redteam:scope — verdict: sound_with_fixes

- **blocker**: The MVP cannot ship in about 12 weeks at 10 hours a week. Build capacity is overstated by a factor of 2 to 3. Evidence: The plan's own 6/2/2 split leaves 14 x 6 = 84 build hours for Phase 0 plus the MVP. In that time it must deliver about 2,900 lines of product code, about 700 lines of bench and demo code, tests for more than 100 matrix cases with property tests, fuzzers and contract tests, 8 CI workflows, image builds for 2 architectures in 4 formats, 4 demo paths and about 25 documents. That is roughly 6,000 to 8,000 lines, or 70 to 95 tested lines per build hour, of kernel and systemd integration code. The prior feasibility brief (bn2heqr9p.txt line 186) advised holding v0.1 to under 1,000 lines of the founder's own code. Line 55 notes that even Universal Blue's small team had to cut image variants. Reviewer estimate: 250 to 400 hours at the full scope, which is 40 to 65 weeks at 6 build hours a week. Fix: Re-baseline publicly. Either keep this scope with a target of April to May 2027, or apply the default cuts below (not contingent ones) to fit 14 weeks. Make the cut order the default plan rather than a week-9 fallback. Add a checkpoint at week 5, not only week 9.
- **major**: Several components are underestimated by factors of 3 to 5. Evidence: The broker is budgeted at 500 lines for OpenAI Chat Completions, Responses and Anthropic Messages with SSE streaming, per-protocol stripping of provider-side tools, model aliasing, credential injection, placement, memory-fit admission and FIFO fairness. Proxies that do this run to thousands of lines. Staging is budgeted at 250 lines, yet it must handle overlay whiteouts, opaque directories, renames, xattr stamping and conflict detection, and pass a property test requiring parity with an unwrapped run over 1,000 operations. The bench and puppet code is budgeted at 700 lines, but every agent has different tool schemas (Hermes terminal, OpenCode bash, Claude Code Bash), so transcripts are needed per agent and per scenario. The session principal is budgeted at 450 lines, but a per-session UID pool conflicts with persistent agent state: ~/.hermes written by one UID cannot be written by the next session's UID, so idmapped mounts or a fixed per-agent UID are needed. Fix: Re-estimate with these ranges: broker 1,500 to 2,500 lines, staging 800 to 1,200, bench and puppet at least 2,000, session at least 800. Accept only OpenAI Chat Completions in v0.1, since Hermes Agent and OpenCode both speak it. Drop memory-fit admission, the cost cap and the ranker stub. Use one fixed UID per agent profile with per-session cgroups.
- **blocker**: Hermes Agent's persistent memory launders labelled data across sessions. This breaks G3 and G4, the cross-vendor claim and the demo. Evidence: The Hermes documentation says the agent writes to ~/.hermes/memories/MEMORY.md and USER.md itself ('The agent saves automatically'), and that this content loads into the system prompt at the next session start. In the plan, labels are stamped only when the workspace is committed. Profile BindPaths (agent state directories) are neither staged nor labelled. A private session can therefore write the canary into MEMORY.md, and the next, unlabelled session sends it to an uncleared cloud model or a public host. OpenCode and Claude Code keep state of the same kind. Fix: Stamp the session label on every file a labelled session creates or modifies on any writable path, not only on commit. Two options: a write-side BPF hook using the security.bpf xattr setter kfunc, or a stamping pass at session end before the unit is released. Add 'agent-memory laundering' to the bypass matrix and to the hero demo. Budget 15 to 25 extra hours.
- **blocker**: The hero demo contradicts itself: once the label rises, the puppet stops driving the agent. Evidence: After the read, the broker must send the private session's model calls to the local llama-server. The real, unmodified Hermes Agent then gets its next responses from a model under 2B parameters, not from the scripted puppet. The remaining exit attempts become stochastic, which conflicts with the gate of 3 identical runs, and a CPU prefill of agent system prompts on a 4 vCPU runner would take tens of seconds per call. Step 6, where OpenCode's call is placed locally, has the same problem. Fix: In the gating suite, register two puppet backends: a local puppet cleared for all labels, and a cloud puppet whose log is the sink. Placement is then proven by which log receives the request. Show real llama-server inference only in a separate, non-gating placement smoke test, and say so in the demo text.
- **major**: A stranger cannot reproduce the demo in under 10 minutes as specified. Evidence: The stated 1.6 GB download leaves out the separate model (about 1 GB), agent installs from PyPI and npm, the sink and bench containers, and the baseline tools. Step 10 reruns steps 4 and 5 under three baselines, which is the whole benchmark run with its own 10-minute target, so the demo takes at least 20 minutes. There are four host paths (GitHub workflow, Linux QEMU, Mac QEMU with HVF, Lima), each injecting the canary differently. Lima provisions guests through a cloud-init NoCloud 'cidata' path, and the plan does not check whether the bootc image supports it. Many Macs have 8 GB of RAM but the VM needs 4 to 6 GB. On Path A the fork runs the project's own workflow YAML, so it is weak evidence of independence, and Actions are disabled on forks by default. Fix: Run v0.1's demo on one primary path: x86_64 Linux with KVM, plus the fork workflow. Show 6 to 8 headline channels and run no baselines in the demo; move them to a separate `tidemark-bench compare`. State the total download size honestly. Gate the release on a measured wall-clock run on a clean runner that includes every download. Document the Mac and Lima paths as best effort, and check cidata support in Phase 0.
- **major**: Releasing a request 'after seeing the exact bytes' is impossible for HTTPS in v0.1. Evidence: v0.1 has no TLS interception: the gateway passes CONNECT tunnels through opaquely. Yet the release primitive and G9 promise that the human sees the exact bytes and that each release records a content hash. For a CONNECT, a release can only grant the whole tunnel, with unbounded content. Fix: In v0.1, allow content-hashed releases only for broker requests (plaintext on the local link) and plain HTTP. Record a CONNECT release as a per-session destination grant without a content hash, and rewrite G9 to match. Alternatively, cut releases entirely and keep only stop and re-run.
- **major**: The tagline promises more than v0.1 delivers. Evidence: 'Agents run without prompts' paired with v0.1's non-guarantee that external side effects are neither held nor undone means an unlabelled session with prompts off can still delete a production database or a repository at a cleared destination: the Replit and PocketOS pattern the plan cites as pain. Holds arrive only in v0.2. Fix: Limit the v0.1 tagline to confidentiality, for example 'what a session has read decides where it may send'. Recommend keeping agent prompts on for tools with side effects until v0.2 holds ship, and say so in the release notes.
- **major**: The first-user choice does not match where those users run agents or where their private data lives. Evidence: Hermes' README says 'Run it on a $5 VPS'. Such a host cannot run llama-server and rarely accepts custom images. Its main mode is an always-on gateway to Telegram, Discord, Slack, WhatsApp and Signal. With domain-level clearance, clearing api.telegram.org for private data also permits sending to any chat, so the guarantee becomes hollow for this mode. Private mail, calendar and notes mostly arrive through APIs or IMAP, which the gateway marks as untrusted but never labels as confidential. Mac users would have to move their data into the VM, and security.* xattrs on virtiofs shared folders are unverified. The OpenClaw profile is deferred to v0.2, yet OpenClaw users are named as primary users. Fix: Narrow the first users to homelab and Proxmox users who run coding or research agents in a dedicated VM over local files. Add destination source labels to the MVP (about 50 to 100 lines): reading from imap.gmail.com, for example, raises the session to private. Exclude VPS hosts, always-on messaging gateways and OpenClaw from the v0.1 claims. Test labels and staging on virtiofs or 9p in Phase 0.
- **major**: The user validation comes too late and asks leading questions. Evidence: The five interviews finish at week 6, after Phase 1 (weeks 3 to 5) is built. 'Will you label your secrets?' invites polite yes answers. The plan also concedes that willingness-to-pay evidence across the research is vendor-survey based. Fix: Hold 8 to 10 observational interviews in weeks 1 and 2. Ask participants to show their agent configuration, data sources, cleared destinations and hosting. Move kill criterion 2 to week 3.
- **major**: The open-source sustainability plan is not credible for a project that makes a public security claim. Evidence: Triage is capped at 2 hours a week, yet the launch invites 'break it' reports, promises a 14-day fix target and a 7-day right of reply to NVIDIA, Anthropic and Omarchy, and withdraws a guarantee if a fix misses 14 days. Reproducing a bypass needs the VM loop on a 16 GB M1, which cannot run the podman machine and the dev VM together. Founders write 66 to 100% of commits even in breakout projects (bn2heqr9p.txt line 170), so 'a second maintainer after 6 months' is unlikely. The distribution advantage of Omarchy-style projects came from an existing audience (line 221). STF excludes prototypes and NLnet excludes AI projects, and no audit is funded before v1.0. Fix: Launch enforce mode as a 'preview', with a 90-day coordinated disclosure window instead of 14 days. Assume a bus factor of 1 through month 12, and treat a second maintainer as a bonus, not a deliverable. Before launch, verify one audit-funding path whose eligibility rules fit the project, and recruit the two or three pre-launch reviewers by name in Phase 0.
- **major**: The post-MVP phases are overloaded by a factor of 3 or more. Evidence: Phase 6 packs capability manifests, an attention queue, a Cockpit plugin, evidence mapping for four jurisdictions, OTel and OCSF export, Rekor anchoring, a policy compiler for two targets, an NVIDIA variant, the Laya evaluation and new BPF hooks into about 130 hours, which is about 78 build hours. Phases 5 and 7 are similar. Fix: Plan one theme per quarter, chosen by user evidence at each milestone. Move the regulatory mapping, Cockpit, the NVIDIA variant, the policy compiler, the attention queue and microVMs into a backlog, each with an explicit trigger.
- **major**: Phase 0 is overloaded, and the kernel gate depends on image work that has not yet been done. Evidence: In about 12 build hours, Phase 0 must deliver a BPF LSM spike on two architectures including overlayfs, a dual-architecture bootc pipeline with cosign, an SBOM and a QEMU smoke test, four ADRs, five documents, two specifications and six governance files. The aarch64 spike runs in Lima on the bootc image, which first needs Lima provisioning to work. Fix: Run the spike on a stock Fedora 45 Cloud VM in Lima (same kernel package, cloud-init included) and on a stock x86 runner. Build the image pipeline for x86_64 only at first. Merge the threat model, non-goals and guarantees into one document v0.
- **minor**: The benchmark and CI targets do not fit the free runners. Evidence: The full benchmark must finish in 10 minutes or less across 5 configurations x (20 adversarial + 20 benign) scenarios, about 200 runs at roughly 2.4 s each, including VM boot and agent installs. Standard runners have a 14 GB SSD, and Bluefin reported that they lacked space for its image (bn2heqr9p.txt line 38). The plan builds qcow2, raw, ISO and AMI images for 2 architectures and also needs room for a model and agents. Fix: Run the full benchmark nightly with a 30 to 45 minute budget, and keep a 5-minute subset on pull requests. Ship only qcow2 and OCI images (OCI for `bootc switch`) in v0.1, and defer ISO, raw and AMI.
- **minor**: The Omarchy baseline adds controversy but no information. Evidence: Baseline (a), the launcher flags with no sandbox, measures the same thing as (d), Tidemark in audit mode: an unenforced agent. Publishing 'Omarchy leaks everything' invites a dispute with a large audience that a solo founder cannot absorb within 2 triage hours a week. OpenShell needs Docker, Podman or host virtualization plus a local gateway, and its default sandbox has 'no agent installed', so the 4-hour cap will likely be exceeded. Fix: In v0.1, compare only Tidemark audit and enforce against sandbox-runtime. Move OpenShell to v0.2. Drop the Omarchy comparison or describe it neutrally as 'unsandboxed'.

### redteam:technical — verdict: sound_with_fixes

- **blocker**: Egress guarantee G4 has a connect-then-taint-then-send TOCTOU: a session that is clean when it opens an HTTPS CONNECT tunnel keeps that tunnel usable after the session is later tainted, so labelled data exits through an already-authorised opaque tunnel. Evidence: The gateway 'implements HTTP CONNECT and plain HTTP forwarding' and reads the label 'synchronously at every decision', but a CONNECT tunnel is decided once at connect time and is thereafter an opaque TLS pipe the gateway cannot re-evaluate. The high-water mark rises only when open() is called, which can happen after the tunnel is established. Nearly all real destinations are HTTPS, so this is the common case. The scripted heist reads the canary before attempting exits, so the demo never exercises this ordering. Fix: Make the gateway event-driven on label rise: subscribe to a teardown-only BPF ring buffer (fail-safe direction, so the map-update race they already avoided does not reappear) and immediately reset every live connection whose destination clearance no longer dominates the new session label. Re-evaluate plain HTTP per request. Add a benchmark scenario that opens the tunnel BEFORE the labelled read.
- **major**: Taint is incomplete because file_open is the only taint point: file descriptors inherited across the shim's exec (systemd credentials, pre-opened files) or received via SCM_RIGHTS let a session read labelled data without ever triggering the hook. Evidence: G3 asserts the label rises 'when any session process opens a labelled file', but inherited and passed fds are read via read()/mmap() on an already-open fd and never call security_file_open in the agent. The shim 'execs the agent', and systemd credentials are delivered as files/fds. The read matrix lists /proc/self/fd reopen but not exec-inherited fds or SCM_RIGHTS fds. Fix: The shim must close_range() all non-essential fds before exec; any inherited or credential fd pointing at a labelled inode must raise the session label at session start; add exec-inherited-fd and SCM_RIGHTS cases to the read matrix.
- **major**: The per-cgroup label map keyed by the immediate cgroup id breaks taint and quota enforcement if a session can create or be placed in a child cgroup. Evidence: Primitive 2 states a 'pinned map keyed by cgroup' and the hook 'exits unless the current cgroup belongs to a session'. If a profile sets Delegate=yes, or the agent creates a sub-cgroup, descendants run under a different cgroup id that the map does not recognise, so their reads are treated as non-session (no taint) and they escape CPUQuota/MemoryMax. Fix: Key the map by the session's ancestor cgroup id (bpf_get_current_ancestor_cgroup_id at the slice-relative level) rather than the exact cgroup id; set Delegate=no explicitly in every profile unless delegation is strictly required; add a sub-cgroup-spawning agent to the containment and read matrices.
- **major**: The whole product rests on the label hook reading security.bpf.tidemark through overlayfs, which may not surface the lower inode's xattr; reads through the staged workspace could therefore go untainted. Evidence: The plan itself lists this as a Phase 0 risk ('overlayfs may interact badly') and bpf_get_file_xattr operates on the opened dentry, which under overlayfs is an overlay inode, not the labelled lower inode. The read matrix includes 'an overlay path', but the fanotify fallback is only defined for a dedicated vault subvolume, not for overlay-mediated reads by children or via mmap. Fix: The Phase 0 kill gate must prove read-through-overlay tainting before any other code; resolve and consult the real underlying inode/path rather than the overlay dentry; and define, before committing to the fallback, how fanotify FAN_OPEN_PERM covers mmap and grandchild reads, since fanotify permission events have known gaps there.
- **major**: The cross-vendor claim 'labels travel with files' is fragile because security.* xattrs are not preserved by ordinary tools, so a committed labelled file loses its label as soon as it is rewritten, and the next agent reads it untainted. Evidence: G6 and the hero demo depend on a committed file carrying the private label to a second vendor's session. Editors (write-temp-then-rename), cp, git checkout, tar, and rsync drop security.* xattrs by default. The scripted demo commits once and reads immediately, so it passes while realistic rewrite workflows silently launder the label away. Fix: Do not rely on per-inode xattrs for travelling labels: maintain a side index keyed by content hash and/or confine labelled data to a dedicated labelled subtree whose whole-tree membership drives taint; narrow the public claim to files Tidemark itself controls; add a rewrite-then-reread laundering scenario to the benchmark.
- **major**: Authority for commit and auto-exec acceptance is unspecified, so a compromised in-session agent may be able to stamp labels, merge into the real workspace, or accept its own auto-exec implant; and writes to any writable bind path not covered by the overlay bypass staging entirely. Evidence: release and stop are explicitly 'refused from inside the agents slice' via SO_PEERCRED, but tidemark diff|commit|abort and the 'explicit acceptance flag' for auto-exec paths carry no such restriction. G6 claims writes are 'invisible to the host until commit', yet only the overlay-covered workspace is staged; a writable profile bind path outside the overlay is a direct, unstaged, unabortable write. Fix: Restrict commit and auto-exec acceptance to processes outside the agents slice (SO_PEERCRED plus cgroup check, as for release); require that every writable path be overlay-staged, and have doctor fail if a profile declares a writable non-overlaid path.
- **major**: Per-open xattr lookups will be costly for open-heavy workloads, and the 'first installable OS' claim does not address mature in-kernel incumbents. Evidence: A sleepable file_open hook doing an xattr syscall on every open inside a session pays that cost 10^5+ times for npm installs or compilers; the 10us p50 / 50us p99 targets assume a warm cache, while a sleepable hook can fault in the xattr. Separately, positioning omits Qubes OS (owned OS with kernel-enforced per-domain isolation and per-domain network control) and SELinux MLS/MCS (in-kernel label-based MAC with SECMARK network control already shipping on Fedora). Fix: Add a mount/path prefilter or a cached 'inode is unlabelled' decision so only files under labelled trees pay the xattr cost; benchmark a real npm/compile workload inside a session. In docs, differentiate explicitly from Qubes and SELinux and give the one-line justification for a bespoke BPF high-water mark (SELinux MLS has no dynamic float, which is a real and statable reason).
- **minor**: Receipt wording overclaims for v0.1: 'tamper-evident, third-party-verifiable' and 'append-only' rest on an on-machine signing key and daemon-enforced (not kernel-enforced) append-only storage. Evidence: The Ed25519 key lives in systemd-creds readable by the ledger UID, and Btrfs offers no owner-proof WORM mode, so a compromise of the ledger daemon (which parses structured events) forges a self-consistent chain that the offline verifier accepts. The plan acknowledges this as a non-guarantee, but the revolutionary_claim and guarantees text still present receipts as externally verifiable evidence. Fix: Present the bring-your-own-canary host-side grep as the primary external check and the receipt as secondary until TPM sealing or Rekor anchoring lands; soften the public 'tamper-evident, third-party-verifiable' wording for v0.1.
- **minor**: Testability gaps: x86 enforcement is only ever exercised in CI, and arm64 GitHub runners may lack /dev/kvm. Evidence: The founder develops on M1 and runs aarch64 under vz, with x86 'only in CI'. BPF verifier and kernel-config differences can produce x86-only bypasses the founder cannot reproduce interactively, slowing the 14-day fix target. The plan adds a step to 'record whether /dev/kvm exists' on arm64 runners, conceding it is uncertain; without KVM, TCG emulation can exceed the 90s boot target. Fix: Budget occasional rented x86 KVM time (consistent with the stated USD 15-40/month) for interactive x86 bypass triage; define a slower boot-time target or a KVM-required gate for the arm64 CI path.
- **minor**: Session units do not explicitly deny user-namespace creation, which widens the single in-scope-excluded risk (kernel LPE). Evidence: The session properties list RestrictNamespaces and a SystemCallFilter but do not state that CLONE_NEWUSER is denied. Unprivileged user namespaces are the dominant Linux local-privilege-escalation vector, and the threat model treats kernel compromise as out of scope, so leaving user-ns open makes that excluded risk materially easier to reach from inside a session. Fix: Set RestrictNamespaces to deny user namespace creation (and all namespaces the agent does not need) in session units; assert this in the containment matrix.

### redteam:incumbent — verdict: sound_with_fixes

- **blocker**: G4 can be bypassed with a tunnel opened before the label rises. The gateway decides once, when a CONNECT arrives, and v0.1 has no TLS interception. A session can open a CONNECT, keep-alive or WebSocket tunnel to an allowlisted public host while it is still unlabelled, then read the secret, then stream the secret through the tunnel that is already open. Evidence: The plan reads the BPF map synchronously at each decision, but decisions happen only per connection, at the domain level. None of the 40-case bypass matrix or the 14 channel classes covers a tunnel that already exists when the label is raised. Fix: On every forwarded chunk from the client, check the session's current label against the tunnel's clearance, or subscribe to label-raise events and close every tunnel whose clearance no longer covers the label. Add 'pre-opened tunnel', 'HTTP/2 connection reuse' and 'WebSocket opened before the read' to the bypass matrix and the benchmark.
- **blocker**: Data can be laundered through persistent agent state, or labels creep for good, and either outcome hits exactly the first-user segment. Hermes Agent memory and skills and OpenClaw state live in profile bind paths. Those paths are not the staged workspace, so writes to them are never stamped with a label. A later session, from any vendor, reads that memory unlabelled and can send it out. If the state directories are labelled instead, every always-on session becomes private permanently. Evidence: G6 stamps labels only when the staged workspace is committed. Profiles bind-mount the agent's own state paths. The first users are defined as people running always-on Hermes or OpenClaw against personal data. ActPlane already propagates labels along file write edges. Fix: Propagate labels in the kernel when a labelled session writes to any persistent path. Use bpf_set_dentry_xattr, which can write security.bpf.* only, from a sleepable LSM hook on file open with write mode and on file permission with MAY_WRITE. Run each user task as a short sub-session with its own label, and treat the agent's long-lived memory as a labelled store. In Phase 1, measure label creep on a real Hermes setup with memory turned on, before building Phase 2.
- **major**: Personal data that arrives over the network never raises a label. Labels rise only when a file is opened. For the target users, mail, notes and finance usually come through APIs (Gmail or IMAP, Notion, banking connectors reached through MCP), not local files. Evidence: The label raise is the file_open BPF hook. The gateway sets only the 'untrusted' integrity flag on content it receives. No destination carries a 'source label'. Fix: Add source labels to destinations. The gateway raises the session label as soon as a CONNECT is allowed to a domain tagged private-source. The broker does the same when a cleared tool or MCP result comes from a labelled source. Add 'read inbox through an API, then exfiltrate' to the hero demo and the benchmark.
- **major**: The receipts are oversold as proof that 'anyone can check from outside'. The signing key sits on the machine, and the image digest in each receipt is a value the machine reports about itself. The machine owner and root are out of scope as adversaries, yet either can sign whatever they like. Microsoft AGT already ships a Merkle audit log, and nono already ships append-only integrity hashing. Evidence: The plan's non-guarantees list already says receipts are 'only as trustworthy as that machine until TPM sealing or external anchoring' is added. The tagline and revolutionary_claim still promise third-party proof. Fix: For v0.1, describe receipts as 'tamper-evident against the agent'. Make the BYOC host-side grep the public proof. Add a TPM2-backed key early: swtpm in QEMU plus systemd-creds --with-key=tpm2. Add TPM quotes or Keylime-style remote attestation before claiming third-party verifiability.
- **major**: The novelty wording is wrong. Shipped operating systems have enforced kernel confidentiality labels on network exits for decades: SELinux MLS (Fedora ships selinux-policy-mls; RHEL earned LSPP certification) and Solaris Trusted Extensions, with labelled networking. The high-water mark itself goes back to ADEPT-50 in 1969. ActPlane hooks connect and can block through BPF-LSM. Evidence: ActPlane's README lists the hooks 'fork / exec / exit / open / unlink / rename / connect', says labels propagate along 'file read/write edges', and says BPF-LSM lets rules 'block before they commit'. The plan says ActPlane has 'no forced exit mediation' and cites only research systems as prior art. Fix: Narrow the claim to: high-water-mark labels on whole process trees of unmodified AI agents from different vendors, with the model API handled as a labelled exit and placement constrained by the label. Cite SELinux MLS, ADEPT-50 and ActPlane's connect hook as prior art.
- **major**: Several direct competitors are missing from the positioning, and the plan treats as unshipped some features that already ship. Evidence: nono (Apache-2.0, about 4.5k stars) ships Landlock/Seatbelt sandboxing, a credential proxy with phantom tokens, allowlist filtering by host and endpoint, 'Content-addressable rollback with SHA-256 dedup and Merkle tree integrity', append-only integrity hashing, and Sigstore-signed instruction files. AgentFS ships copy-on-write isolation and a SQLite audit log. agentsh ships an LLM proxy with DLP, approvals, soft_delete and checkpoints. container-use gives each agent a branch that can be discarded. Rubrik Agent Rewind (August 2025) rolls back files, databases and repositories. Invariant Guardrails ships data-flow rules across tool calls through an MCP or LLM proxy. Fix: Add all of these to the positioning. Add nono as a benchmark baseline. Reframe white space 4 as 'permission-diff admission' only, since signed instruction files already ship. Describe staging and receipts as integration, not novelty.
- **major**: The plan mischaracterises Anthropic's sandbox-runtime and NVIDIA's NemoClaw. Evidence: The srt README says it sandboxes 'the entire process tree', including local MCP servers, and has credential masking and a filterRequest hook. Only Claude Code's built-in sandbox is limited to shell commands. NemoClaw needs no NVIDIA account for local Ollama or vLLM and runs on Linux with rootless Podman and on Windows through WSL, so 'no account, hardware-neutral' does not set Tidemark apart. Fix: Base the differentiation only on label-driven exits, model placement and evidence. Treat NemoClaw as the main competitor in the self-hoster segment.
- **major**: OpenShell now documents an extension point that both threatens Tidemark and duplicates its gateway. Supervisor middleware 'inspects, transforms, or denies allowed HTTP and WebSocket traffic'. NVIDIA could therefore add a taint check as middleware in one release. Meanwhile, building tidemark-gw contradicts Tidemark's own non-goal of not rebuilding containment. OpenShell already has credential placeholders, request-level read/write limits and OCSF events, so Tidemark v0.1 is weaker on L7 rules and credentials. Evidence: docs.nvidia.com/openshell/extensibility/overview and /observability. The plan dropped OpenShell as 'undocumented'. Fix: Build the label engine and the clearance evaluator as a library with two front ends: the gateway inside the image, and an OpenShell supervisor-middleware adapter in v0.2. The incumbent then becomes a distribution channel rather than a kill criterion.
- **major**: v0.1 has a credential paradox. A git push to a cleared remote is in the benign suite, but without TLS termination the token has to live inside the session. If the token is labelled, it taints the session and the push is denied. If it is not labelled, an injected agent can exfiltrate it. OpenShell, srt, nono, Docker and Sprites already keep credentials out of the session. Evidence: Credential placeholder substitution is deferred to v0.2, together with goproxy. Fix: Move minimal credential injection into the MVP: TLS termination for git smart-HTTP and two declared APIs. Otherwise, remove the push task from the benign suite and document the limitation.
- **major**: Phase 5 holds desynchronise the agent from reality. An irreversible request is held, the agent receives a 403, and the request is replayed after a human approves it. The agent believes the action failed, so it may retry by another route and duplicate the effect. Temporal's sagas and agentsh's soft_delete or approve decisions are prior art for this pattern. Evidence: Phase 5 replays approved requests with an idempotency key after a timeout. Fix: Hold the request synchronously with a bounded wait. Alternatively, define an explicit 'queued' response contract in the profiles. Present the effect journal as integration on top of saga patterns.
- **minor**: The benchmark is built around the one property only Tidemark has. OpenShell and srt allow allowlisted hosts by design, so the leak table is predetermined. Evidence: Of the 14 channel classes, 4 (upload, query string, push, model prompt) are leaks only under label semantics. The suite has no channels where the incumbents should win. Fix: Pre-register the methodology. Add channels where incumbents lead (L7 method and path rules, credential exfiltration, completion rate). Report benign completion for every baseline next to its leak count.
- **minor**: The claim that 'model calls move, not fail' is overclaimed. Evidence: When no cleared backend fits the context or the memory limit (on the M1, prefilling 30K tokens takes 2 to 4 minutes), the request fails. Fix: Reword as 'never silently placed outside clearance; fails closed with a structured reason'.
- **major**: The calendar does not hold. The 6/2/2 weekly split gives 72 build hours in the 12-week MVP, not 120. Those hours must cover about 3,600 lines of product and benchmark code plus security tests. Evidence: Budget section and weekly cadence. Fix: Budget the MVP in build hours, or cut it to labels, the gateway with tunnel revocation, the broker and the BYOC demo, moving staging and the dashboard to v0.2.

