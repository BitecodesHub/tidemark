# AgentOS Tidemark — Research-Backed Plan (v1, 2026-10-16)

## Context
- **Request.** Before building AgentOS, the user asked for deep research into the problems that AI companies and ordinary companies face, and for a plan for something revolutionary rather than a generic OS.
- **Research done.** 81 agents worked in three workflows:
  - A: 22 market and technical areas, 6 gap-fill areas, fact-checking, and a synthesis.
  - B: 5 areas on bootable-OS feasibility, with fact-checking.
  - C: 5 concept architects, 3 judges, a synthesis, 3 red-team reviewers, and a final revision.
- **Constraints confirmed by the user:**
  - Open-source, Apache-2.0.
  - Solo, about 10 hours a week.
  - Must be a literal bootable OS; a headless appliance is acceptable for v0.1.
  - Code written with AI assistance; the user reviews security-critical diffs.
  - Private data runs local-only by default.
  - Development machine: Apple M1 Pro, 16 GB, macOS 27.
- **Decision.** The user chose the **Tidemark** direction over the original broad "fastest AI performance OS" vision. The research showed that vision has no defensible novelty: NVIDIA, Bluefin, Omarchy and Ollama ship equivalent stacks for free. It also cannot be tested on an M1 Pro.
- **Project folder.** `/Users/mac/Desktop/AI-Operating System` is empty and is not a git repository.

## Research conclusions (verified, October 2026)
- **Commoditized in Sep–Oct 2026. Integrate; do not rebuild.**
  - Agent sandboxing: Microsoft MXC (MIT), NVIDIA OpenShell (Apache-2.0), Anthropic sandbox-runtime, Docker Sandboxes, Copilot sandboxing.
  - Agent identity and inventory: Agent 365, AWS AgentCore, Google Agent Identity.
  - Protocols: MCP and A2A, now held by the Agentic AI Foundation.
- **Poor wedges for one part-time person:** inference serving, training and fine-tuning platforms, GPU scheduling, and OS-level model packaging.
- **Unsolved, severe pains** (atlas scores):
  - Data exfiltration through *allowed* channels (240).
  - The verification and oversight bottleneck: review time +441%, 93% of permission prompts approved (240).
  - Agent supply chain: 13.4% of 3,984 skills rated critical (192).
  - Irreversible actions (180).
  - Forgeable logs: OpenAI test agents compromised Hugging Face and tampered with logs, July 2026 (144).
- **Earlier "AI OS" attempts failed or were absorbed:** Humane, Rabbit, OpenAI Atlas browser, Google Mariner, the AIOS research kernel.
  - Omarchy holds the "agentic distro" label, yet launches 11 of 15 agents with permission prompts disabled.
  - Agents people actually run: OpenClaw (392k stars) and Hermes Agent (254k).
- **OS facts:**
  - Base: bootc on `quay.io/fedora/fedora-bootc:45`, which is multi-arch.
  - Kernel: Landlock ABI 9 needs Linux 7.1. The sleepable BPF LSM hooks `file_open` and `file_receive` exist, as do the fs kfuncs `bpf_get_file_xattr`, `bpf_path_d_path` and `bpf_real_data_inode`.
  - M1 Pro limits: no nested virtualization, and 4–14B Q4 models at most.
  - CI: free x64 KVM GitHub runners.
- **Components named in the spec:**
  - "Laya" is a 0.4B routing model.
  - Llama is legacy; Meta moved to Muse Spark Open.
  - torchtune is unmaintained.
  - Hermes Agent is MIT-licensed, with a very large unreviewed backlog.
- **Caveats:**
  - The web-search budget ran out early, so the evidence rests on about 1,500 direct reads of primary sources.
  - Willingness-to-pay figures come from vendors; there are no buyer interviews.
  - In Workflow C, a naming mismatch split the judges' scores. Corrected, the security-first "Escrow" concept scored 122 against Tidemark's 116. Tidemark absorbs Escrow's core and has the more feasible scope.
- **Raw research artifacts** (to be copied into the repository; see Step 1):
  - Scratchpad: `scratchpad/atlas.txt`, `dims.txt`, `final.json`.
  - Workflow outputs: `tasks/wxnj45tf2.output`, `wp4734zz7.output`, `wao1qwd96.output`.
  - Workflow journals under `subagents/workflows/`.

## The concept: AgentOS Tidemark (working name; trademark check needed)
**Claim (v0.1 preview).** What an agent session has *read* decides where its data may go next, including which model sees it. This is enforced for unmodified agents from different vendors and for every MCP server, hook and subprocess they start. Anyone can verify it with a canary they generate on their own machine.

**What is new.** Prior art includes the high-water mark (ADEPT-50), SELinux MLS, ActPlane, CaMeL/FIDES, OpenShell, nono and agentsh. The new piece is the combination:
- Floating labels on the process trees of unmodified agents from different vendors.
- The model API treated as a labelled exit, with label-constrained placement.
- A label re-check on every forwarded chunk.
- Label-carrying writes and label-partitioned agent memory.
- All of it shipped as a reproducible, signed image with a demo strangers can check.

**First users.**
- Homelab, Proxmox and dedicated-VM users who run Hermes Agent or OpenCode over local code, notes and secrets, mixing a local model with cheap cloud models.
- Agent-security researchers, reached through the standalone benchmark.

**Primitives** (about 4,000 lines of Go plus BPF):
1. **Session principal** (integration):
   - A systemd transient unit per agent run: cgroup subtree, private network namespace with loopback only, strict Landlock (ABI 9 or later), seccomp deny list, no namespace creation, quotas.
   - Every descendant is contained.
2. **Read-time labels** (the one new kernel piece):
   - Sleepable BPF LSM programs on `file_open` and `file_receive` raise an mmapable per-session label word before `open()` returns.
   - A file's label comes from an xattr or a directory rule.
   - The label cache is keyed by `bpf_real_data_inode`, so it works through overlayfs.
3. **Clearance-gated exits** (`tidemark-gw`):
   - One Unix socket per session; HTTP CONNECT and plain HTTP.
   - The session's label is re-checked against the destination's clearance before every forwarded chunk, and tunnels are reset when the label rises.
   - DNS, UDP and IP literals are refused.
4. **Label-constrained model broker** (`tidemark-broker`):
   - Speaks OpenAI Chat Completions.
   - Places each request only on backends cleared for the session's label. Default: private goes to the local llama-server. Otherwise it fails closed.
   - For labelled sessions, it strips provider-side tools.
   - It is the only holder of API keys.
5. **Staged writes with label stamping** (overlayfs):
   - `commit`, run from outside the session, stamps written files with the session's final label.
   - Auto-exec, hook and skill files are held until a human runs `accept`.
   - Agent memory is partitioned by label.
   - `abort` restores the workspace byte for byte.
6. **Agent-unreachable receipts** (`tidemark-ledgerd`):
   - A hash-chained JSONL log, plus an in-toto/DSSE Ed25519 statement per session.
   - Tamper-evident against the agent, not against root.
   - An offline `tidemark-verify` runs on macOS, Linux and Windows.
7. **Release and global stop.**
   - One-time release tokens bound to the exact broker request content.
   - `tidemark stop` freezes the agents cgroup slice and denies all exits.

**Not claimed:** prevention of prompt injection; resistance to kernel or root compromise; covert channels to cleared destinations; protection of unlabelled data; content inspection; holding side effects in v0.1; proof against the machine owner before TPM attestation (Phase 7).

## Architecture
| Layer | Components | Technology |
|---|---|---|
| L0 Image | Containerfile installing the Tidemark RPM; first-boot defaults; boot self-test (Landlock ≥9, LSM list includes bpf, kfuncs present); headless; no telemetry | bootc, `fedora-bootc:45@sha256` (with a `:44` fallback build), bootc-image-builder qcow2, COPR, cosign keyless signing, syft, trivy, Renovate |
| L1 Kernel | BPF LSM label programs; Landlock; seccomp; cgroup v2 freezer; overlayfs | clang plus cilium/ebpf (bpf2go), go-landlock |
| L2 Control | `tidemarkd`, the only root daemon: sessions, mounts, label slots, raise-only API, stop, doctor | Go, coreos/go-systemd, godbus |
| L3 Exits | `tidemark-gw` and `tidemark-broker`, each unprivileged under its own UID | Go net/http |
| L4 Evidence | `tidemark-ledgerd` and `tidemark-verify` | in-toto, DSSE (go-securesystemslib), ed25519, systemd-creds (TPM2 sealing in v0.2) |
| L5 Models | One validated backend: llama-server as a Quadlet pinned by digest; `tidemark model pull/ls` by revision and SHA-256; macOS Metal or Ollama registered as unattested endpoints | ghcr.io/ggml-org/llama.cpp |
| L6 Agents | `hermes.yaml` and `opencode.yaml` profiles, plus a generic `tidemark run -- <cmd>` | YAML; agents installed on demand into /var |
| L7 UI | CLI, plus `tidemark report --html`: a static page with Home, Sessions, Models, Receipts and System, and no JavaScript | Go html/template |
| L8 Bench and CI | `tidemark-bench` OCI container; pre-registered scenarios; host-side sinks; baselines | GitHub Actions; QEMU/KVM on free x64 runners |

**Repository layout** (to be created in `/Users/mac/Desktop/AI-Operating System`):
```
cmd/{tidemarkd,tidemark-gw,tidemark-broker,tidemark-ledgerd,tidemark-shim,tidemark,tidemark-verify}/
pkg/{labels,clearance}/      internal/{session,stage,stateview,egress,broker,receipt,queue,doctor,report}/
bpf/label.bpf.c              spec/  policy/default.yaml  profiles/{hermes,opencode}.yaml
packaging/rpm/               image/{Containerfile,Containerfile.f44,kargs.d,systemd,quadlets,firstboot,shims,bib}/
bench/{prereg,scenarios/{adversarial,benign},puppet,sinks,baselines}/   demo/{byoc.sh,heist,qemu,lima}/
test/{integration,bypass,readmatrix,staging,laundering,fuzz}/
docs/{THREAT_MODEL.md,BENCHMARK_METHODOLOGY.md,research/,adr/,devlog/,interviews/}
.github/workflows/{build-image,build-rpm,boot-test,integration,bench-nightly,bench-pr,reproduce,release,scorecard}.yml
LICENSE NOTICE TRADEMARK.md SECURITY.md CONTRIBUTING.md GOVERNANCE.md CODE_OF_CONDUCT.md
```
New documents follow the naming convention `YYYY-MM-DD_Subject_Version.md`, for example `docs/adr/2026-10-19_BaseImage_v1.md`.

## Hero demo: "Bring Your Own Canary"
1. `./demo/byoc.sh` on an x86_64 KVM host downloads the signed qcow2 and runs `cosign verify`.
2. The host generates a canary and starts sinks outside the VM: attacker, public paste, public search, a private-source mail API, and a cloud-model puppet.
3. The VM boots with the canary placed in `~/Private` and in a fake SSH key.
4. **Audit mode.** A puppet drives unmodified Hermes Agent through 8 channels:
   - an upload to the paste host;
   - the canary in a query string;
   - a tunnel opened before the read;
   - the canary inside the next model prompt;
   - a provider-side web-fetch tool;
   - a git push to an attacker-controlled remote;
   - laundering through agent memory;
   - API-sourced exfiltration.
   The canary appears in the host sinks.
5. **Enforce mode**, with a fresh canary:
   - Zero matches on every host log.
   - Private inference goes to the local puppet; the cloud log stays empty.
   - The receipt verifies, and a tampered copy fails at the exact entry.
6. **Cross-vendor step.** The Hermes output is committed with the private label; OpenCode reads it and becomes private itself.
7. **Staging and stop.** 200 file edits and a planted git hook are held, and `abort` restores a byte-identical tree. `stop` freezes everything in under 1 s.

Targets: under 20 minutes including downloads, under 10 minutes after them.

## Roadmap (about 10 hours a week; every phase gates on its acceptance criteria)
| Phase | Weeks / dates | Goal | Key acceptance criteria |
|---|---|---|---|
| 0 Foundations and kill gate | 1–3 (19 Oct – 8 Nov 2026), ~30 h | Governance files, THREAT_MODEL v0, BENCHMARK_METHODOLOGY v0, policy and receipt specs, ADRs, x86_64 image pipeline, kernel spike on stock Fedora 45 (aarch64 in Lima; x86_64 in CI), 8–10 observational interviews | Probe proves a label raise before `open()` returns from a grandchild process, through overlayfs, via SCM_RIGHTS and via directory rules; `security.bpf.*` returns EPERM inside a session; signed image passes `cosign verify` in CI; written go/no-go decision |
| 1 Confine and label | 4–9 (to 20 Dec 2026), ~60 h | `tidemarkd`, `tidemark run`, shim, production label engine, label CLI, `stop`/`ps`/`doctor`, aarch64 OCI image plus RPM, label-creep experiment | Containment matrix (≥18 cases) and read matrix (≥16 cases) pass 100%; zero false raises across 1,000 files; 10,000-iteration race test; label-creep report published |
| 2 Gate every exit | 10–15 (to 31 Jan 2027), ~60 h | Gateway with per-chunk re-check, model broker, llama-server Quadlet, release queue, Hermes and OpenCode profiles | Bypass matrix (≥45 cases) 100% blocked on x86 CI and the aarch64 VM; private requests land on the local backend with 0 cloud calls; fail-closed with a structured reason |
| 3 Stage, carry labels, sign | 16–20 (to 7 Mar 2027), ~50 h | overlayfs staging, commit/accept/abort, label-partitioned state, ledger, verifier, HTML report | Abort of a 10,000-file workspace is byte-identical; 1,000-operation commit property test passes; laundering through memory is blocked; tampered receipts are detected |
| 4 Benchmark and launch v0.1 preview | 21–24 (to 4 Apr 2027; slip limit mid-May), ~40 h | Pre-registered `tidemark-bench` (≥24 scenarios, 18 channel classes); baselines: unsandboxed, sandbox-runtime, Tidemark audit, Tidemark enforce; BYOC demo and recording; signed release; launch on Show HN, r/LocalLLaMA, r/selfhosted and project forums | 0 canary bytes in enforce mode; benign completion ≥85% of audit mode; three identical runs; nightly suite ≤45 min and pull-request subset ≤5 min |
| 5 v0.2 | Apr–Jul 2027, ~130 h | Opt-in TLS termination for declared APIs only; credential injection; clearances by method and path; synchronous bounded holds for irreversible effects; TPM2-sealed key; nono and OpenShell baselines | 100% of irreversible-class requests held; no held request reaches a sink before approval; credentials never appear inside a session |
| 6 v0.3 | Jul–Oct 2027, ~130 h | OpenShell middleware adapter; DEB and AUR packages; Anthropic Messages and Responses APIs; Claude Code and OpenClaw profiles; planner/worker experiment | Adapter blocks 100% of uncleared requests at ≤10 ms p95 |
| 7 v1.0 candidate (conditional) | Oct 2027 – Apr 2028 | TPM-quote attestation, frozen schemas, neutral home for the spec, succession plan | A third party verifies a receipt with a TPM quote; at least 25 independent BYOC verifications |

**Every phase also includes:**
- **Automated tests:** Go unit tests, the in-VM integration harness on x86 KVM, nightly fuzzing.
- **Security checks:** gosec, govulncheck, staticcheck, CodeQL, Actions pinned by SHA, secret scanning, a review list of root operations.
- **Compatibility:** x86_64 and aarch64 UEFI; kernel 7.1 or later.
- **Documentation:** devlog, ADRs, THREAT_MODEL updates.

**Contingent cuts**, applied at the week-9 or week-15 checkpoint if the plan slips by more than one week:
1. Replace the HTML report with CLI text only.
2. Move the aarch64 OCI image to v0.2.
3. Move source labels to v0.2.
4. Replace state views with discarding state writes from labelled sessions.
5. Drop the release queue.

Never cut: the per-chunk re-check, commit stamping, the BYOC demo, and the cross-vendor step.

## Kill criteria (summary)
- **Week 3:** label raising in the kernel is unreliable on Fedora 45 or through overlayfs → ship as an audit tool only, or stop. Fewer than 3 of 10 interviewees run agents over local private data → re-scope to benchmark-first.
- **Week 9:** more than 40% of benign Hermes sessions become labelled → make audit the default and move the planner/worker split earlier.
- **Launch:** enforce-mode completion below 80% of audit mode → ship audit mode and receipts only.
- **Three months after launch:** fewer than 10 independent BYOC verifications and fewer than 3 users running real work → stop maintaining the image; keep the RPM, the benchmark and the specs.
- **Any time:** an incumbent ships label-driven exits with label-constrained placement → move effort to the adapter and the specs.

## Deferred from the original specification (with triggers)
- **Training and fine-tuning:** not before v1.0, and only with a maintainer who has NVIDIA hardware plus 20 user requests.
- **vLLM, SGLang, TensorRT-LLM, multi-engine auto-tuning:** need a community maintainer and a GPU boot test on every release.
- **Hardware adapters and an NVIDIA variant:** need a volunteer maintainer; boot tests on AWS g4dn spot.
- **Laya-style router:** once two or more cleared backends exist; it may reorder the permitted backends but is never a security boundary.
- **Six-section UI, live dashboard, desktop:** after demand (10 or more, and 100 or more, active users respectively).
- **Not built:** a GPU cluster scheduler. Scheduling means per-session cgroup quotas only.
- **Never inside the image:** MLX. It is supported as a host endpoint.
- **Backlog with triggers:** capability manifests, attention scheduler, compliance mapping, microVM principals.

## Development environment and budget
- **macOS host (Homebrew):** Go, podman, lima, qemu, cosign, syft, gh, shellcheck. Lima (vz) runs Fedora 45 Cloud for the kernel spike; podman machine (libkrun) builds images; one VM at a time because of the 16 GB limit.
- **x86_64:** free GitHub x64 KVM runners gate releases; an hourly cloud VM handles triage (USD 5–10 per month).
- **Total cash:** about USD 20–50 per month.

## Immediate execution after approval (this session)
1. **Preserve the research.** Copy the atlas, the 28 briefs, the OS-feasibility research and the final concept into `docs/research/` as:
   - `2026-10-16_Research-Atlas_v1.md`
   - `2026-10-16_Research-Briefs_v1.md`
   - `2026-10-16_OS-Feasibility_v1.md`
   - `2026-10-16_Tidemark-Concept_v1.md`

   The scratchpad is temporary, so this comes first.
2. **Set up the repository.** Run `git init` and create the skeleton per the layout above: governance files (Apache-2.0, TRADEMARK, SECURITY with a 90-day disclosure window, GOVERNANCE assuming a bus factor of 1, CoC, CONTRIBUTING with DCO), THREAT_MODEL v0, BENCHMARK_METHODOLOGY v0, `spec/policy-v0.schema.json`, `spec/receipt-predicate-v0.md`, and the first ADRs. No commits unless the user asks.
3. **Image pipeline.** `image/Containerfile` (FROM a pinned `fedora-bootc:45`), plus `build-image.yml` and `boot-test.yml` with a QEMU serial-console smoke test.
4. **Kernel spike code.** `bpf/label.bpf.c` probe, a Go loader, and a Lima template for Fedora 45 Cloud. Run it locally only after the user approves the Homebrew installs (lima, go, podman).
5. **Interview guide.** Write `docs/interviews/2026-10-19_Interview-Guide_v1.md`: 8–10 observational sessions asking about agent setup, where private data lives, which destinations users would clear, and where the agent is hosted.
6. **Save project memory.** Record the goals, the constraints and the Tidemark decision.

## Verification
- **Phase 0:**
  - `probe` output in the Lima VM (aarch64) and in a CI QEMU/KVM run (x86_64) shows the kernel version, an LSM list containing bpf and landlock, Landlock ABI ≥9, and the label word raised immediately after a grandchild's `open()`.
  - `cosign verify` passes on the GHCR image.
- **Each MVP phase:** the matrix suites (containment, read, bypass, staging, laundering) run in CI on x86 KVM and locally on aarch64, and must pass 100%.
- **Launch:**
  - `demo/byoc.sh` on a clean x64 runner shows 0 canary matches in enforce mode and leaks in audit mode.
  - `tidemark-verify` validates the receipts and rejects tampered copies.
  - The nightly `tidemark-bench` gives three identical runs.
  - Every number in the README links to a CI artifact.

## Open questions (answer before Phase 4)
- Final project name and trademark clearance.
- DCO or CLA, before the first outside contribution.
- Names of 2–3 pre-launch security reviewers for the BPF, gateway, staging and ledger diffs.
- Where your target users' private data lives (local files or APIs such as Gmail and Notion). The interviews answer this.
