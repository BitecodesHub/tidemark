# AgentOS Tidemark (working name)

> **Status: Phase 0, as of 23 October 2026. No product is implemented; the kernel mechanism is measured, not shipped.**
> This repository holds research, plans, governance drafts, specifications and one measurement
> instrument: the Phase 0 kernel spike in [`spike/`](spike/README.md), which builds and runs inside a
> Fedora 45 virtual machine and passed 16 of 16 scenarios on aarch64 and, in CI under QEMU/KVM, on
> x86_64 on 2026-10-23 (results in [`docs/spike-results/`](docs/spike-results/)). There is **no product to
> install, no daemon, no image to download and no clearance enforcement running anywhere.** The image
> pipeline does build and boot both Fedora variants in CI, and by decision
> [`PD-17`](docs/adr/2026-10-23_Phase0-Founder-Decisions_v1.md) it publishes nothing: no GHCR push and
> no signing happen until publication is deliberately turned on. Every
> feature below is a design target, not a shipped capability. The name "AgentOS Tidemark" is a
> working name pending trademark clearance (see [TRADEMARK.md](TRADEMARK.md)). The repository is
> currently private; the interview coding rules, spike evidence and decisions are all in it and are
> published with the launch.

## The claim (target for the v0.1 preview)

**What an agent session has read decides where its data may go next, including which model sees it,
and anyone can check this with a canary they generate on their own machine.**

The claim is meant to hold for unmodified agents from different vendors (Hermes Agent and OpenCode
in v0.1) and for every MCP server, hook and subprocess they start. It is a confidentiality claim
only. The list of things Tidemark does **not** claim follows below and is part of the claim.

## What Tidemark is meant to be

Tidemark is planned as an open-source (Apache-2.0), headless, signed operating-system image built
with [bootc](https://github.com/bootc-dev/bootc) on Fedora (`quay.io/fedora/fedora-bootc:45`,
pinned by digest, with a `:44` fallback build until launch). The software will also ship as a
Fedora RPM that the image installs. Full guarantees will be claimed only on the signed image, and
only when `tidemark doctor` passes. The image is designed with no telemetry, and the default policy
keeps private data on the local machine.

The design integrates existing containment and adds one new kernel piece. The planned primitives
are:

| Primitive | Planned mechanism | Status |
|---|---|---|
| Session principal | One systemd transient unit per agent run: cgroup subtree, loopback-only network namespace, strict Landlock (ABI 9 or later), seccomp deny list, no namespace creation | Not started |
| Kernel-recorded read-time labels | Sleepable BPF LSM programs on `file_open` and `file_receive` raise the session's label word before the labelled `open()` returns; a file's label is the union of its `security.bpf.*` xattr and every matching directory rule | Phase 0 kernel spike passed 16 of 16 scenarios on aarch64 (Fedora 45 Beta, kernel 7.2.0-61, three identical runs, 2026-10-23; 12 of 12 on 2026-10-18) and 16 of 16 on x86_64 (same kernel version, five consecutive GitHub Actions QEMU/KVM runs, 2026-10-23), including the 10,000-iteration raise-before-return race in both the opener and a separate observer reading and the overlayfs case on each architecture. Measured on an instrument that only observes: no production labeller exists |
| Clearance-gated exits | `tidemark-gw`: one Unix socket per session; the session's current label is re-checked against the destination's clearance before **every forwarded chunk**; tunnels are reset when the label rises; DNS, UDP and IP literals are refused | Not started |
| Label-constrained model broker | `tidemark-broker`: OpenAI Chat Completions; a labelled session is placed only on backends cleared for its label (by default the local llama-server), otherwise it fails closed with a reason; the broker alone holds API keys | Not started |
| Staged writes with label stamping | overlayfs staging; `commit` (run outside the session) stamps written files with the session's final label; auto-exec, hook and skill files wait for `accept`; agent memory is partitioned by label; `abort` restores the workspace byte for byte | Not started |
| Agent-unreachable signed receipts | `tidemark-ledgerd`: hash-chained JSONL plus an in-toto/DSSE Ed25519 statement per session; tamper-evident against the agent, not against root; offline `tidemark-verify` for macOS, Linux and Windows | Not started |
| Release and global stop | One-time release tokens bound to the exact broker request; `tidemark stop` freezes every session and denies every exit | Not started |

The preview guarantees G1 to G9 (containment, single exit, label before read, clearance on every
chunk, model channel, staged writes carry labels, evidence against the agent, stop, and
declassification) are defined in the threat model. They are targets until the test matrices that
check them pass in CI.

## Measured overhead (Phase 0 spike, 2026-10-23)

The spike instrument records `open(2)` latency around the label hook, 20,000 iterations per case
after 2,000 warmup. These are floors measured on an observing instrument, not a production
labeller, and they are the numbers any future optimiser has to beat:

| Environment | `open()` outside a session (p50) | `open()` inside a session, unlabelled file (p50) | Labeller delta (p50) |
|---|---|---|---|
| aarch64, Lima VM (Fedora 45) | 835 ns | 1,211 ns | about **+376 ns** |
| x86_64, CI guest under QEMU/KVM | 2,535 ns (median of four runs) | 4,098 ns | about **+1.56 µs** (the guest is CPU-oversubscribed) |

The xattr-labelled steady state, the raise path and the directory-rule path are measured too, and
the full five-run tables with p99 and mean are in
[`docs/spike-results/`](docs/spike-results/), with the method in
[`docs/BENCHMARK_METHODOLOGY.md`](docs/BENCHMARK_METHODOLOGY.md).

## Not claimed

Tidemark does not claim, and will not claim in v0.1:

- prevention of prompt injection, or correctness of agent output;
- resistance to kernel or root compromise;
- protection against covert or timing channels to destinations the user has cleared;
- protection against misuse by destinations the user has cleared, including a cleared cloud model;
- protection of unlabelled data;
- content inspection (v0.1 decides per destination domain);
- holding or undoing external side effects (holds for declared APIs are planned for v0.2);
- third-party proof against the machine's owner (TPM attestation is a conditional Phase 7 goal);
- always-on messaging-gateway modes (Telegram, Discord, Slack, WhatsApp, Signal), VPS hosting and
  OpenClaw;
- labels on files copied off the machine or rewritten by tools outside sessions (directory rules
  excepted);
- protection of data that arrives through APIs, unless the destination is tagged as a private
  source;
- containment or labels for agents started outside `tidemark run`;
- elimination of covert or timing channels between concurrent sessions;
- detection of encoded or disguised content in a release that a human approves;
- truth of receipt content when a Tidemark event source is already compromised (receipts detect
  later tampering only);
- availability (quotas limit resource abuse but do not eliminate it);
- fine-grained flow control: once a session reads labelled data, it loses public egress for the
  rest of its life;
- protection through macOS or other host model endpoints, which are recorded as unattested.

Until v0.2, keep each agent's own approval prompts on for destructive tools.

## Bring Your Own Canary (future deliverable, Phase 4)

The planned public proof is a demo a stranger can run on their own x86_64 Linux host with KVM. It
does not exist yet. As designed, `./demo/byoc.sh` will:

1. download the signed qcow2 and verify it with cosign;
2. generate a random canary on the host and start sinks outside the VM (an attacker sink, a public
   paste host, a public search host, a private-source mail API and a cloud-model puppet);
3. boot the VM with the canary placed in a private notes file and a fake SSH key;
4. run an unmodified Hermes Agent, driven by a scripted, fully compromised puppet, through eight
   exfiltration channels in **audit** mode, where the canary must appear in the host sinks;
5. repeat with a fresh canary in **enforce** mode, where every host log must show zero matches,
   private inference must land on the local backend, and the receipt must verify while a tampered
   copy fails at the exact entry;
6. show the cross-vendor step: Hermes output committed with the private label makes OpenCode
   private when it reads that output;
7. show staging and stop: 200 edits and a planted git hook are held, `abort` restores a
   byte-identical tree, and `stop` freezes everything (target: 1 s or less).

Targets: 20 minutes or less including downloads, 10 minutes or less after them. The proof is the
grep on the host; receipts are a secondary check.

## Platform requirements (planned; verified in Phase 0)

- **Proposed kernel floor: Linux 7.2 or later, with Landlock ABI 9 or later** (labeller ADR, status
  Proposed; checklist decision FD-G). Landlock ABI 9 first ships in Linux 7.1
  (`landlock_abi_version = 9` in `security/landlock/syscalls.c` at the v7.1 tag), but
  `bpf_real_data_inode`, which keys the optional overlay-safe label cache, first appears at the
  **v7.2** tag and is absent at v7.1. Support for 7.1 kernels is deferred; the labeller ADR
  revisits the floor if checklist item C2 selects the no-cache design after the Phase 0 spike.
- BPF LSM active in the boot LSM list, with sleepable programs on `file_open` and `file_receive`
  (both are in the sleepable hook set in `kernel/bpf/bpf_lsm.c`).
- The fs kfuncs `bpf_get_file_xattr`, `bpf_path_d_path` and `bpf_real_data_inode`
  (`fs/bpf_fs_kfuncs.c`). Fedora 44 and 45 kernels were in the 7.2.x series when checked on
  16 October 2026; the Phase 0 probe confirms what the pinned image actually ships.
- Development hosts: the aarch64 kernel spike is planned to run in a Lima VM on an Apple silicon
  Mac; free GitHub x86_64 KVM runners are the only gating environment. No nested virtualization is required.

## Repository map

Present today (16 October 2026; every file is a Phase 0 draft, and nothing here is built or tested
yet):

| Path | Contents |
|---|---|
| `docs/2026-10-16_Project-Plan_v1.md` | Approved project plan |
| `docs/2026-10-16_Phase0-Checklist_v1.md` | Phase 0 checklist |
| `docs/THREAT_MODEL.md`, `docs/BENCHMARK_METHODOLOGY.md` | Threat model v0 with G1 to G9; benchmark pre-registration draft |
| `docs/adr/` | Architecture decision records and their index |
| `docs/devlog/` | Development log |
| `docs/interviews/` | Interview guide, screener and synthesis template (the anonymised synthesis follows the interviews) |
| `docs/research/` | Concept, research atlas, research briefs, bootable-OS feasibility study, raw outputs |
| `spec/` | Policy schema v0, clearance model and receipt predicate v0 |
| `policy/default.yaml` | Default policy (no component reads it yet) |
| `docs/2026-10-16_Image-Pipeline_v1.md` | Image pipeline design: base pinning, self-test, CI workflows, verification |
| `image/` | Containerfiles (Fedora 45 and the Fedora 44 fallback), first-boot and self-test units for the image pipeline test (no enforcement) |
| `spike/` | Phase 0 kernel spike: BPF probe, Go loader, 16 scenario scripts and Lima template (passed 16 of 16 on aarch64 and on x86_64 in CI on 2026-10-23, 12 of 12 on 2026-10-18; results in `docs/spike-results/`) |
| `.github/workflows/` | Image build, boot test, kernel spike, lint and Scorecard workflows (not yet run) |
| `renovate.json` | Dependency update rules for pinned digests and actions |
| `LICENSE`, `NOTICE` | Apache License 2.0 and project notice |
| `README.md`, `SECURITY.md`, `CONTRIBUTING.md`, `GOVERNANCE.md`, `CODE_OF_CONDUCT.md`, `TRADEMARK.md` | Governance drafts |
| `.gitignore`, `.editorconfig` | Ignore rules (including generated BPF objects and key material) and editor formatting rules |

Planned layout (directories appear as work lands; see the concept for the full tree):

| Path | Planned contents | Phase |
|---|---|---|
| `bpf/` | Production label program (the Phase 0 probe lives in `spike/bpf/`) | 1 |
| `cmd/`, `pkg/`, `internal/` | Daemons, CLI and verifier in Go; `pkg/labels` and `pkg/clearance` as importable libraries | 1 to 3 |
| `profiles/`, `packaging/rpm/` | Hermes Agent and OpenCode profiles; RPM spec (built in CI and installed into the image; the same source goes to COPR) | 1 to 2 |
| `test/`, `bench/`, `demo/` | Matrix suites; `tidemark-bench`; the BYOC demo | 1 to 4 |

The Go module path is `github.com/BitecodesHub/tidemark`, set on 2026-10-23 when the repository was
created (Phase 0 decisions ADR, PD-10). The repository is private until launch; the
links in these documents point at it.

## Roadmap

About 10 hours a week, solo. Every phase gates on its acceptance criteria.

| Phase | Dates | Goal |
|---|---|---|
| 0 Foundations and kill gate | 19 Oct to 8 Nov 2026 | Governance, threat model v0, benchmark methodology v0, specs, ADRs, x86_64 image pipeline, kernel spike on stock Fedora 45, 8 to 10 observational interviews, written go or no-go decision |
| 1 Confine and label | to 20 Dec 2026 | `tidemarkd`, `tidemark run`, production label engine, `stop`/`ps`/`doctor`, label-creep experiment |
| 2 Gate every exit | to 31 Jan 2027 | Gateway with per-chunk re-check, model broker, llama-server Quadlet, Hermes Agent and OpenCode profiles |
| 3 Stage, carry labels, sign | to 7 Mar 2027 | overlayfs staging, commit/accept/abort, label-partitioned state, ledger, verifier |
| 4 Benchmark and v0.1 preview | to 4 Apr 2027 (slip limit mid-May 2027) | Pre-registered `tidemark-bench`, BYOC demo, signed release |
| 5 v0.2 | Apr to Jul 2027 | Credentials and synchronous holds for irreversible effects at declared APIs |
| 6 v0.3 | Jul to Oct 2027 | OpenShell adapter, more packages, more agents |
| 7 v1.0 candidate (conditional) | Oct 2027 to Apr 2028 | TPM-quote attestation, frozen schemas, neutral home for the specification |

Kill criteria are public. For example, if kernel label raising proves unreliable on Fedora 45 or
through overlayfs at the end of Phase 0, the kernel-label claim is dropped and the project ships as
an audit tool, or stops. See the plan for the full list.

## Documents

- [Project plan](docs/2026-10-16_Project-Plan_v1.md)
- [Concept: primitives, threat model, phases, benchmark methodology, kill criteria](docs/research/2026-10-16_Tidemark-Concept_v1.md)
- [Bootable-OS feasibility research](docs/research/2026-10-16_OS-Feasibility_v1.md)
- [Research atlas](docs/research/2026-10-16_Research-Atlas_v1.md) and [research briefs](docs/research/2026-10-16_Research-Briefs_v1.md)
- [Threat model](docs/THREAT_MODEL.md) (Phase 0 draft)
- [Benchmark methodology](docs/BENCHMARK_METHODOLOGY.md) (Phase 0 draft)
- [Architecture decision records](docs/adr/)

## How the project is run

- Licence: [Apache-2.0](LICENSE), intended to be permanent (pending founder confirmation; see
  [GOVERNANCE.md](GOVERNANCE.md)).
- Security reports: private GitHub security advisories, 90-day coordinated disclosure
  ([SECURITY.md](SECURITY.md)).
- Contributions: DCO sign-off (proposed; the founder confirms before the first outside
  contribution), disclosure of AI assistance, and two-person review for security-critical areas
  (until reviewers are confirmed, the founder's review alone applies and each change needs a second
  review before launch) ([CONTRIBUTING.md](CONTRIBUTING.md)).
- Governance: one maintainer, so the bus factor is 1 ([GOVERNANCE.md](GOVERNANCE.md)).
- Conduct: [Contributor Covenant 2.1](CODE_OF_CONDUCT.md).
- Development uses AI-assisted coding. The maintainer reviews every diff in the BPF program, the
  gateway copy loop, the staging merge and the ledger line by line, and named reviewers are sought
  before launch.

## Credits and prior art

Tidemark is planned to build on these projects (licences as reported by their repositories on
16 October 2026):

- [bootc](https://github.com/bootc-dev/bootc) (MIT or Apache-2.0) and the
  [Fedora bootc base images](https://gitlab.com/fedora/bootc/base-images): the image foundation.
- [Landlock](https://landlock.io/) in the Linux kernel, through
  [go-landlock](https://github.com/landlock-lsm/go-landlock) (MIT).
- [cilium/ebpf](https://github.com/cilium/ebpf) (MIT): loading the BPF programs from Go.
- [llama.cpp](https://github.com/ggml-org/llama.cpp) (MIT): the one validated local inference
  backend.
- [Hermes Agent](https://github.com/NousResearch/hermes-agent) (MIT) and
  [OpenCode](https://github.com/anomalyco/opencode) (MIT): the unmodified agents the v0.1 profiles
  target.

The ideas are not new on their own. Tidemark integrates prior work and claims none of it:

- the high-water mark (ADEPT-50, 1969), SELinux MLS/MCS ([SELinux project](https://github.com/SELinuxProject/selinux))
  and Solaris Trusted Extensions, which enforce static labels on exits;
- HiStar, Flume and Asbestos, research systems for OS-level information-flow control;
- [ActPlane](https://github.com/eunomia-bpf/ActPlane) (MIT), which propagates kernel labels along
  file edges with BPF-LSM;
- CaMeL ([Defeating Prompt Injections by Design](https://arxiv.org/abs/2503.18813);
  [code](https://github.com/google-research/camel-prompt-injection), Apache-2.0) and FIDES, which
  enforce information flow inside one agent;
- [NVIDIA OpenShell](https://github.com/NVIDIA/OpenShell) (Apache-2.0),
  [Anthropic sandbox-runtime](https://github.com/anthropics/sandbox-runtime) (Apache-2.0),
  [nono](https://github.com/nolabs-ai/nono) (Apache-2.0) and
  [agentsh](https://github.com/canyonroad/agentsh) (Apache-2.0), which already ship agent
  containment, egress control, credential handling, staging or audit features.

What Tidemark adds is a narrow combination: floating labels on the process trees of unmodified
agents from different vendors, the model API treated as a labelled exit, a label re-check on every
forwarded chunk, label-carrying writes and label-partitioned memory, shipped as a signed image with
a demo strangers can check. Naming a project here does not imply its endorsement of, or affiliation
with, Tidemark.

## Licence

Copyright 2026 The AgentOS Tidemark Authors. Licensed under the [Apache License, Version 2.0](LICENSE).
See [NOTICE](NOTICE) for the project notice.

The source repository does not currently include third-party software that requires attribution.
No binary release exists yet. The planned binary distributions (the RPM, the OCI images and the
disk images) will contain third-party software under its own licences. Each release is planned to
ship a software bill of materials together with the licence texts and notices that those licences
require; an SBOM alone does not satisfy those requirements.

The Apache License does not grant permission to use the project's names or marks. See
[TRADEMARK.md](TRADEMARK.md).
