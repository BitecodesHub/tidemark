# Phase 0 Checklist (v1)

| Field | Value |
|---|---|
| Status | Phase 0 draft, revised 16 October 2026 after review, 17 October 2026 in a cross-file consistency pass, and 23 October 2026 after the aarch64 spike and the founder-decision record. Boxes are ticked when their listed evidence exists in the repository. Items marked "AI-drafted" still await founder line-by-line review; the 23 October tick was made under the founder's delegation that day (see `docs/adr/2026-10-23_Phase0-Founder-Decisions_v1.md`, Deciders). |
| Version | v1, prepared Friday 16 October 2026 |
| Phase | Phase 0, Foundations and kill gate: weeks 1 to 3, Monday 19 October to Sunday 8 November 2026 |
| Gate date | **Sunday 8 November 2026** (end of week 3): written go or no-go decision |
| Budget | About 30 founder hours planned (3 weeks at about 10 hours per week). The estimates below total 32.75 hours; **D2 and D7 are now deferred to Phase 1 (PD-6), which puts Phase 0 at 31.75 hours.** |
| Sources | `docs/2026-10-16_Project-Plan_v1.md` (roadmap, Phase 0 row, verification); `docs/research/2026-10-16_Tidemark-Concept_v1.md` (Phase 0 deliverables, acceptance criteria, tests, security checks, kill criteria) |

## Executive summary
- This checklist lists every Phase 0 deliverable and acceptance criterion from the concept document, with an owner, an hour estimate and a target week. The estimates sum to 32.75 hours: 11.20 in week 1, 11.20 in week 2 and 10.35 in week 3; **31.75 hours after the PD-6 deferral of D2 and D7.**
- **The overrun is settled (PD-6):** D2 and D7 move to Phase 1 week 4 and the remaining 1.75-hour overrun is accepted, with actual hours tracked weekly. A realistic range stays 32 to 40 hours if the AI drafts need substantial rework.
- The concept splits Phase 0 into 21 building hours and 9 other hours. This checklist allocates 17.5 building hours and 15.25 other hours, because 8 interviews with their per-interview records need about 10.5 hours.
- Two kill criteria are decided on Sunday 8 November 2026: kernel label raising (criterion 1) and interview evidence (criterion 2). Section 6 holds the decision-record template. **The concept's full wording governs both criteria and a triggered criterion cannot be overridden (PD-7, recorded 2026-10-23, before the x86_64 run).** Criterion 1 is 16 of 16 on aarch64 as of 2026-10-23 and pending on x86_64; criterion 2 stands at zero of eight interviews.
- One kernel fact changed the plan's wording and then changed back: `bpf_real_data_inode` exists from Linux 7.2, not 7.1, and the labeller used to depend on it. **It no longer does (PD-1: no inode-keyed cache), so the 7.2 floor is now a testing-coverage choice rather than a kfunc requirement (PD-2).** Fedora 44 and 45 ship 7.2 kernels (section 10).

## 1. Calendar
| Week | Dates (2026) | Focus |
|---|---|---|
| 1 | Mon 19 Oct to Sun 25 Oct | Governance and threat-model review; coding-rule freeze; recruiting; first interviews (target 3); Lima VM; probe drafting starts; signed image build pipeline; reviewer invitations |
| 2 | Mon 26 Oct to Sun 1 Nov | Most interviews (target 5); probe and core spike tests on aarch64; devlog 2 |
| 3 | Mon 2 Nov to Sun 8 Nov | Spill-over interviews only if needed; remaining spike tests; x86_64 spike in CI; boot smoke test; specifications and ADR reviews; repository settings; synthesis by Fri 6 Nov; go or no-go on Sun 8 Nov; devlog 3 |

## 2. How to read this list
- **Owner: Founder.** The founder does the work and makes the decision.
- **Owner: AI-assisted.** An AI coding or writing assistant drafts; the founder reviews every line before acceptance. For the BPF program, the founder reviews line by line.
- **Hours** are founder hours, including review. Where an AI draft already exists, the hours cover review, correction and verification, not authoring.
- **Evidence** names the file or CI artifact that proves the item. Tick a box only when that evidence exists.
- The concept's weekly cadence before launch is 7 hours building, 2 on docs and devlog, and 1 on interviews and community. Phase 0 departs from it because the interviews fall here.

## 3. Deliverables

### A. Repository governance (0.75 h)
- [x] **A1. LICENSE and NOTICE.** Apache-2.0 full text; NOTICE with copyright line. Owner: AI-assisted. Hours: 0.10. Week: 1. Evidence: `LICENSE`, `NOTICE`. **Present 2026-10-23; `NOTICE` also records the BPF dual-licence exception (PD-3).**
- [x] **A2. TRADEMARK.md.** "AgentOS" and "Tidemark" are working names; trademark search pending before Phase 4. Owner: AI-assisted. Hours: 0.10. Week: 1. Evidence: `TRADEMARK.md`. **Present 2026-10-23.**
- [x] **A3. CONTRIBUTING.md.** DCO sign-off (the plan's choice); the founder records the DCO or CLA decision. Owner: Founder decision, AI-assisted draft. Hours: 0.10. Week: 1. Evidence: `CONTRIBUTING.md`. **Present; DCO decided under PD-11.**
- [x] **A4. SECURITY.md.** Private advisories and a 90-day coordinated disclosure window. Owner: AI-assisted. Hours: 0.10. Week: 1. Evidence: `SECURITY.md`. **Present 2026-10-23.**
- [x] **A5. GOVERNANCE.md.** Assumes a bus factor of 1. Owner: AI-assisted. Hours: 0.05. Week: 1. Evidence: `GOVERNANCE.md`. **Present 2026-10-23.**
- [x] **A6. CODE_OF_CONDUCT.md.** Owner: AI-assisted. Hours: 0.05. Week: 1. Evidence: `CODE_OF_CONDUCT.md`. **Present 2026-10-23.**
- [x] **A7. README.md.** What the project is, what it does not claim, Phase 0 status, and the placeholder module path `github.com/agentos-tidemark/tidemark`. Owner: AI-assisted. Hours: 0.25. Week: 1. Evidence: `README.md`. **Present; the module path is now `github.com/BitecodesHub/tidemark` (PD-10) and the status lines record the 16/16 aarch64 spike with x86_64 pending.**

### B. Threat model, benchmark methodology and specifications (1.75 h)
- [x] **B1. docs/THREAT_MODEL.md v0.** One document: threat model, non-goals, guarantees G1 to G9 and the hosts that provide each guarantee; kernel floor checked against section 10. Owner: AI-assisted (founder line-by-line review). Hours: 0.75. Week: 1. Evidence: `docs/THREAT_MODEL.md`. **Present; revised 2026-10-23 for PD-1, PD-2, PD-4 and the 16-scenario spike coverage. Founder line-by-line review still outstanding.**
- [x] **B2. docs/BENCHMARK_METHODOLOGY.md v0.** Pre-registration draft: 18 gating channel classes, 2 non-gating incumbent-favourable classes, configurations (a) to (d), fairness rules, host-side canary handling. Owner: AI-assisted. Hours: 0.50. Week: 3. Evidence: `docs/BENCHMARK_METHODOLOGY.md`. **Present 2026-10-23.**
- [x] **B3. spec/policy-v0.** Labels, directory rules, clearances and source labels. Owner: AI-assisted. Hours: 0.25. Week: 3. Evidence: `spec/policy-v0.schema.json`. **Present 2026-10-23.**
- [x] **B4. spec/receipt-predicate-v0.** in-toto and DSSE receipt predicate. Owner: AI-assisted. Hours: 0.25. Week: 3. Evidence: `spec/receipt-predicate-v0.md`. **Present 2026-10-23.**

### C. Architecture decision records (0.5 h)
Named `docs/adr/YYYY-MM-DD_Subject_v1.md`. AI drafts of all five exist; the hours cover founder review.
- [x] **C1. Base image.** `quay.io/fedora/fedora-bootc:45` pinned by digest with the resolution date; `:44` fallback build until launch. Owner: AI-assisted. Hours: 0.10. Week: 1. Evidence: `docs/adr/2026-10-16_BaseImage_v1.md`.
- [x] **C2. Labeller.** BPF design; overlay behaviour observed; lookup method chosen (`bpf_get_file_xattr` per open versus a cache keyed by `bpf_real_data_inode`, available from Linux 7.2); fanotify fallback coverage gaps. Owner: Founder (after the spike). Hours: 0.25. Week: 3. Evidence: `docs/adr/2026-10-16_Labeller-BPF-LSM_v1.md`, updated with the spike results. **Closed 2026-10-23: the lookup method is no cache (PD-1), chosen on the measured latency baseline and on the measured declassification failure of the cache; the overlay behaviour is recorded from scenario 03; the ADR now carries the 16/16 aarch64 results and PD-4's link/rename decision. The ADR stays status Proposed until the founder promotes it, because the x86_64 half of the kill test has not run.**
- [x] **C3. Go as the language**, with the Go version pinned. Owner: AI-assisted. Hours: 0.05. Week: 3. Evidence: `docs/adr/2026-10-16_Language-Go_v1.md`.
- [x] **C4. No TLS interception in v0.1.** Owner: AI-assisted. Hours: 0.05. Week: 3. Evidence: `docs/adr/2026-10-16_No-TLS-Interception-v0.1_v1.md`.
- [x] **C5. Delivery as an RPM built in CI and installed into the signed image, with the same source on COPR.** Owner: AI-assisted. Hours: 0.05. Week: 3. Evidence: `docs/adr/2026-10-16_Delivery-RPM-Plus-Image_v1.md`.

### D. Kernel spike on stock Fedora 45 Cloud (12.5 h)
- [x] **D1. Lima VM.** Lima (vz) template for stock Fedora 45 Cloud, aarch64; boot; record kernel version, LSM list and Landlock ABI. One VM at a time (16 GB host). Owner: Founder (AI-drafted template). Hours: 0.75. Week: 1. Evidence: the Lima template `spike/lima/fedora-45.yaml` (the plan's layout places the long-term template under `demo/lima/`) and the captured output.
- [ ] **D2. Lima cidata on fedora-bootc.** Record whether Lima's cidata provisioning works on a fedora-bootc disk image. Owner: Founder. Hours: 0.50. Week: 3. Dependency: an aarch64 fedora-bootc disk image (FD-B). **Deferred to Phase 1 week 4 on 2026-10-23 (PD-6); the FD-B dependency goes with it. Threat-model assumption S8 stays untested, so Lima cidata behaviour on bootc is not claimed.**
- [x] **D3. BPF probe (`spike/bpf/label_probe.bpf.c`; the plan names it `bpf/label.bpf.c`).** Sleepable `lsm.s/file_open` and `lsm.s/file_receive` (the draft also attaches `lsm.s/mmap_file`); session keyed by its ancestor cgroup at a fixed depth with `Delegate=no` (tested in D4); a file's label is the union of its `security.bpf.*` xattr (through `bpf_get_file_xattr`) and every matching directory rule (through `bpf_path_d_path`); label word in an mmapable BPF array; early exit for processes outside sessions; bounded loops; failed kfunc calls handled by failing closed (an unknown-label bit). Owner: Founder, line-by-line review of an AI draft. Hours: 4.75 (2.00 in week 1, 2.75 in week 2). Evidence: `spike/bpf/label_probe.bpf.c` and review notes.
- [x] **D4. Go loader and spike tests.** cilium/ebpf with bpf2go, built inside the Lima VM or in CI (Go is not installed on the host). Tests: map entry after a child's and a grandchild's `open()`; **sub-cgroup test: a grandchild placed by the test harness in a nested cgroup below the session cgroup raises the session's slot (ancestor-cgroup keying at a fixed depth), and a session process cannot move itself out of the session cgroup under `Delegate=no`**; overlay path comparing `bpf_get_file_xattr` with the `bpf_real_data_inode` cache; descriptor passed with SCM_RIGHTS; directory rule; xattr invalidation; read-only mmap of the label word; EPERM on `security.bpf.*` inside a session. The checked-in draft now has 16 scenarios, mapped to assumptions in `spike/README.md` (a passing run ends with `SUMMARY: ALL PASS (16/16)`): scenarios 13 to 16, added 2026-10-23, are the move-out case AC3 requires, xattr and directory-rule invalidation, the 10,000-iteration race in both the opener and the separate-observer reading, and the `open()` latency baseline. Owner: AI-assisted. Hours: 2.75 (2.00 in week 2 for the tests AC1 needs; 0.75 in week 3 for the remaining tests, run on both architectures with D6). Evidence: test source (`spike/cmd/tidemark-probe/`, `spike/tests/`) and output. **Status 2026-10-18:** built and run on aarch64; 12 of 12 scenarios passed in three identical runs. **Status 2026-10-23: 16 of 16 in three consecutive runs**, plus 16 of 16 with the `Dual BSD/GPL` licence tag and 16 of 16 with the inode cache compiled out (`TM_DISABLE_INODE_CACHE=1`). Evidence: `docs/spike-results/2026-10-23_Kernel-Spike-Results-aarch64_v2.md` and the five run logs beside it. The x86_64 run (D6) is the remaining half.
- [x] **D5. aarch64 run.** Run the probe in Lima; capture its output; record labelled and unlabelled `open()` latency as a baseline (no target). Owner: Founder. Hours: 1.00. Week: 2. Evidence: output file committed or attached to devlog 2. **Done 2026-10-23: three runs plus the two configuration builds are committed under `docs/spike-results/`, and scenario 16 records the latency baseline (20,000 iterations per case): 835 ns p50 outside a session, 1,211 ns inside a session on an unlabelled file, 1,211 ns on an xattr-labelled file, 1,711 ns on the raise path, 1,252 ns with a matching directory rule - about +376 ns p50 for the labeller on a session open.**
- [ ] **D6. x86_64 run.** Run the probe under QEMU/KVM on a GitHub-hosted x64 runner with a stock Fedora 45 Cloud image; publish the output as a CI artifact (`.github/workflows/kernel-spike.yml`). Owner: AI-assisted. Hours: 1.75. Week: 3. Evidence: CI run link.
- [ ] **D7. Shared-folder behaviour.** Record how xattrs and the append-only attribute behave on virtiofs and 9p shares. Owner: Founder. Hours: 0.50. Week: 3. **Deferred to Phase 1 week 4 on 2026-10-23 (PD-6). Until it is recorded, macOS users whose data lives in virtiofs or 9p shares stay outside every v0.1 claim (assumption S8).**
- [ ] **D8. Fanotify fallback coverage.** Document what the fanotify fallback sees and misses for child, grandchild, mmap and overlay reads, before it can be selected. Owner: AI-assisted, founder review. Hours: 0.50. Week: 3. **Contingent empirical test:** if any kernel row in the kill-criterion-1 table (section 6) that maps to the kill gate fails (see the spike mapping note under that table) and the concept's wording governs (FD-I), D8 becomes a test, not a document: child, grandchild, mmap and overlay reads under fanotify permission events (`FAN_OPEN_PERM`, plus `FAN_ACCESS_PERM` or `FAN_PRE_ACCESS` where the kernel offers them) on both architectures. A paper analysis cannot clear the criterion. Contingent hours: about 3.0, not in the totals; they come from the slip plan (section 9).

### E. x86_64 image pipeline (5.0 h)
- [x] **E1. Containerfiles.** `image/Containerfile` FROM `quay.io/fedora/fedora-bootc:45@sha256:<digest>` with the digest and resolution date recorded (the tag is rebuilt frequently), and `image/Containerfile.f44` for the fallback. Owner: AI-assisted. Hours: 0.50. Week: 1. Evidence: `image/Containerfile`, `image/Containerfile.f44`. **Both present with the 2026-10-16 digests; the image label namespace now matches the real owner (`io.github.bitecodeshub.tidemark.phase`, PD-10).**
- [ ] **E2. build-image.yml.** Builds both the `:45` image and the `:44` fallback; pushes to `ghcr.io/<owner>/tidemark` (owner lowercased); cosign keyless signing and `cosign verify`; syft SBOM; trivy in report-only mode; every Action pinned by commit SHA; build in 30 minutes or less. Owner: AI-assisted. Hours: 2.00. Week: 1. Evidence: green CI run with the verify step for both builds.
- [ ] **E3. boot-test.yml.** bootc-image-builder qcow2; QEMU/KVM serial-console smoke test that runs the probe; boot to multi-user in 90 seconds or less. Owner: AI-assisted. Hours: 2.00. Week: 3. Evidence: CI run with serial log.
- [ ] **E4. Repository settings and CI hygiene.** Branch protection, secret scanning, DCO check, gofmt, staticcheck and the Scorecard workflow. Owner: Founder. Hours: 0.50. Week: 3. Dependency: the GitHub repository exists (FD-C).

### F. Observational interviews (10.5 h)
- [ ] **F1. Prepare.** Finalise the guide and screener; freeze the pre-registered coding rules, record decisions FD-1 to FD-6 (synthesis template, section A) and publish the freeze (public commit, or SHA-256 in devlog 1) by Monday 19 October, before interview P01. Owner: Founder. Hours: 0.50. Week: 1. Evidence: freeze date and published link in the synthesis template. **Rules frozen 2026-10-23 (PD-12): every default in the synthesis template section A.4 is accepted unchanged, because each default is the stricter reading. The freeze is two weeks later than the checklist asked and is recorded as a deviation; it still precedes P01, which has not happened. The SHA-256 of the frozen file is published in devlog 1; the repository is private for now, so the SHA route is the one taken.**
- [ ] **F2. Recruit.** Post manually in each channel after reading its rules; screen; schedule. The devlog does not recruit. Owner: Founder. Hours: 0.75 (0.50 in week 1, 0.25 in week 2).
- [ ] **F3. Interview.** Eight interviews of 45 minutes (ten add 1.5 hours). Owner: Founder. Hours: 6.00 (2.25 in week 1 for 3 interviews, 3.75 in week 2 for 5; 0 in week 3, which is for spill-over only, until Wed 4 Nov, recorded as a deviation and paid for from the slip plan).
- [ ] **F4. Per-interview records and administration.** Section B record written from the raw notes within 24 hours of each interview; notes checked against any recording before it is deleted; scheduling and consent messages. Owner: Founder. Hours: 2.00 (8 × 0.25; 0.75 in week 1, 1.25 in week 2). This allowance is tight for a form of about 60 fields. Evidence: records held outside the repository; their count noted in the synthesis.
- [ ] **F5. Synthesise.** Aggregate tables; small cells merged; publication checklist; kill criterion 2 evaluated. Section B records stay outside the repository and are processed only by local tools; cloud AI tools see only the published synthesis. Owner: Founder. Hours: 1.25. Week: 3 (by Fri 6 Nov). Evidence: `docs/interviews/2026-11-06_Interview-Synthesis_v1.md` (proposed name).

### G. Pre-launch reviewers (0.75 h)
- [ ] **G1. Reviewers.** Two or three pre-launch security reviewers recruited and confirmed by name, with consent recorded. Consent records stay outside the repository; names are published only with each reviewer's agreement. Owner: Founder. Hours: 0.75 (0.25 in week 1 to identify candidates and send invitations; 0.50 in week 3 to follow up, confirm and record consent).

### H. Devlogs (0.5 h)
- [ ] **H1. Devlog 1** reviewed and published after the coding-rule freeze, with the SHA-256 of the frozen file filled in (drafted 16 October: `docs/devlog/2026-10-16_Devlog_v1.md`). Owner: Founder. Hours: 0.10. Week: 1.
- [x] **H2. Devlog 2**, about Sunday 1 November: interview progress and aarch64 spike results. Owner: AI-assisted. Hours: 0.20. Week: 2. Evidence: `docs/devlog/2026-10-18_Devlog_v2.md`, written 2026-10-18; it reports no interviews because none had happened by then.
- [ ] **H3. Devlog 3**, Sunday 8 November: the gate decision and its evidence. Owner: AI-assisted. Hours: 0.20. Week: 3. **The third published entry went out early on 2026-10-23 (`docs/devlog/2026-10-23_Devlog_v3.md`) with the four new scenarios and the decision record; the gate-decision entry is still owed on Sunday 8 November and becomes devlog 4, so this box stays unchecked.**

### I. Gate (0.5 h)
- [ ] **I1. Go or no-go decision record** written and published using section 6. Owner: Founder. Hours: 0.50. Week: 3 (Sunday 8 November 2026).

## 4. Acceptance criteria
From the concept document's Phase 0 acceptance criteria and deliverables, and the plan's Phase 0 row. Hours are covered by the tasks in brackets.
- [x] **AC1. aarch64 probe.** In Lima, the probe prints the kernel version, an LSM list containing `bpf` and `landlock`, a Landlock ABI of 9 or later, and the label word already raised when read in the same process immediately after a grandchild's `open()` returns. [D1, D3, D4, D5] **Met 2026-10-18:** kernel 7.2.0-61.fc45, LSM list includes `bpf` and `landlock`, Landlock ABI 10, grandchild label `0x1` sampled right after `open()` (scenario 01); evidence `docs/spike-results/2026-10-18_Kernel-Spike-Results-aarch64_v1.md`. **The 10,000-iteration race variant named in the labeller ADR now runs (added 2026-10-23): scenario 15 reports `RACE_SELF iterations=10000 misses=0` and `RACE_OBSERVE iterations=10000 misses=0` - the second reading from a separate observer process outside the session - in three consecutive runs. See `docs/spike-results/2026-10-23_Kernel-Spike-Results-aarch64_v2.md`.**
- [ ] **AC2. x86_64 probe.** The same four outputs under QEMU/KVM in CI. [D6] **Blocked on the repository, which now exists (PD-10): `.github/workflows/kernel-spike.yml` runs the same 16 scenarios under QEMU/KVM, so this criterion is closed by the first green CI run and its serial-log artifact.**
- [x] **AC3. Raise paths and keying.** The probe proves a label raise before `open()` returns from a grandchild process, through overlayfs, through a descriptor passed with SCM_RIGHTS, and through a directory rule; ancestor-cgroup keying holds (a grandchild in a nested cgroup below the session raises the session's slot, and a session process cannot move itself out of the session cgroup under `Delegate=no`); `security.bpf.*` returns EPERM inside a session. [D3, D4] **Met on aarch64 2026-10-23: scenarios 01, 02, 03, 04, 05, 07 and 13, in three runs. Scenario 13 is the move-out case - an unprivileged in-session process gets `EACCES` moving itself to a sibling scope and to the cgroupfs root, `EACCES` creating a sub-cgroup under its scope, and still raises `0x1` afterwards. One asymmetry is recorded rather than hidden: root *can* write `cgroup.procs`, which is why a session is granted no privilege (guarantee G1) and not why the kernel refuses it. x86_64 stays open under AC2.**
- [x] **AC4. Overlay decision.** The overlay behaviour and the chosen lookup method are recorded in an ADR. [C2, D4] **Done 2026-10-23: the labeller ADR records the measured overlay behaviour (two `file_open` events, layer-relative backing path, rules skipped for it) and chooses no inode-keyed cache (PD-1), with the no-cache 16/16 run - overlay scenario included - as the evidence that the design does not need the kfunc.**
- [ ] **AC5. Fallback documented.** The fanotify fallback's coverage for child, grandchild, mmap and overlay reads is documented before it can be selected, and tested empirically if section 6 requires it. [D8]
- [ ] **AC6. Signed image.** A signed x86_64 image is on GHCR, and `cosign verify` passes in CI. [E1, E2]
- [ ] **AC7. Interviews.** The interview synthesis is published, and kill criterion 2 is evaluated at the end of week 3. [F1 to F5]
- [ ] **AC8. Decision.** A written go or no-go decision is recorded. [I1]

## 5. Tests, security checks, compatibility, performance and documentation
From the concept document's Phase 0 section and the plan's "every phase also includes" list.

**Automated tests**
- [ ] T1. x86_64 image build in CI, for both `:45` and the `:44` fallback. [E2]
- [ ] T2. QEMU/KVM boot smoke test that runs the probe. [E3]
- [ ] T3. BPF spike tests: map entry after a child's and a grandchild's open; sub-cgroup keying under `Delegate=no`; overlay path; descriptor passed with SCM_RIGHTS; directory rule; xattr invalidation; read-only mmap of the label word; EPERM on `security.bpf.*`. [D4]
- [ ] T4. Landlock ABI and LSM-list assertions. [D4, E3]
- [ ] T5. gofmt, staticcheck, the DCO check and the Scorecard workflow. [E4]

**Security checks**
- [ ] S1. Cosign keyless signing and verification; no long-lived keys in the repository. [E2]
- [ ] S2. Actions pinned by commit SHA; branch protection; secret scanning. [E2, E4]
- [ ] S3. SBOM generated; trivy scan in report-only mode. [E2]
- [ ] S4. BPF program reviewed for bounded loops, early exit for processes outside sessions, and handling of failed kfunc calls. [D3]
- [ ] S5. gosec, govulncheck, CodeQL and a review list of root operations. The plan requires them in every phase; this checklist proposes deferring them to Phase 1 week 4, when `tidemarkd` and its root operations exist (FD-J). Not in the totals.

**Compatibility**
- [ ] K1. `fedora-bootc:45` (multi-architecture manifest) with `:44` as the fallback; stock Fedora 45 Cloud images for the spike. [C1, E1, E2, D1, D6]
- [ ] K2. QEMU with KVM or HVF; Lima (vz). [D1, D6, E3]
- [ ] K3. M1 Pro host; no nested virtualization required. [D1]

**Performance targets**
- [ ] P1. Image build of 30 minutes or less on x86_64. [E2]
- [ ] P2. Boot to multi-user in 90 seconds or less under KVM. [E3]
- [ ] P3. Baseline latency of labelled and unlabelled `open()` recorded; no target yet. [D5]

**Documentation**
- [ ] DOC1. README stating what the project is and what it does not claim. [A7]
- [ ] DOC2. THREAT_MODEL v0 and BENCHMARK_METHODOLOGY v0 draft. [B1, B2]
- [ ] DOC3. ADRs, the interview synthesis and devlogs 1 to 3. [C1 to C5, F5, H1 to H3]

## 6. Go or no-go decision record (template)
Copy into `docs/adr/2026-11-08_Phase0-GoNoGo_v1.md` (proposed name) on Sunday 8 November 2026.

```markdown
# Phase 0 go or no-go decision (v1)

| Field | Value |
|---|---|
| Date | Sunday 8 November 2026 |
| Decider | Founder |
| Status | [Proposed or Accepted] |
| Evidence cut-off | [date and time] |

## Kill criterion 1: kernel label raising
| Check | aarch64 (Lima) | x86_64 (CI QEMU/KVM) | Evidence link |
|---|---|---|---|
| Kernel version | [ ] | [ ] | [ ] |
| LSM list contains bpf and landlock | [ ] | [ ] | [ ] |
| Landlock ABI (9 or later required) | [ ] | [ ] | [ ] |
| Label raised before a grandchild's open() returns | [ ] | [ ] | [ ] |
| Ancestor-cgroup keying: grandchild in a nested sub-cgroup raises the session slot; move-out fails under Delegate=no | [ ] | [ ] | [ ] |
| Raise through overlayfs (method used) | [ ] | [ ] | [ ] |
| Raise through SCM_RIGHTS (file_receive; spike A5, and A5b for an unprivileged receiver) | [ ] | [ ] | [ ] |
| Raise through a directory rule (bpf_path_d_path) | [ ] | [ ] | [ ] |
| Raise through a file mapping (mmap_file; spike A8, and A8b for a descriptor opened before the session) | [ ] | [ ] | [ ] |
| Read-only mmap of the label word (spike A6) | [ ] | [ ] | [ ] |
| EPERM on security.bpf.* inside a session | [ ] | [ ] | [ ] |
| Fanotify fallback, empirical test (D8; required if a row above that maps to the kill gate fails and the concept wording governs): child / grandchild / mmap / overlay reads seen | [ ] | [ ] | [ ] |

Spike mapping (`spike/README.md`, "How a failure maps to the kill criteria"): a failure of the grandchild, keying, directory-rule or file-mapping row (spike A1, A2, A4 or A8) maps to this criterion directly, and a failure of the overlay row (A3) maps to it for the overlay case. A failure of A5, A5b or A8b narrows the public claim, and a failure of the label-word mapping row (A6) requires a redesign of the label transport; none of these is a kill on its own. A failure of the EPERM row (A7) means label integrity is broken. Scenario A4b (a mount-namespace alias) is an observation and never fails the suite.

The two sources word this criterion differently:
- Concept: label raising is unreliable on stock Fedora 45 on either architecture or through overlayfs, AND the fanotify fallback misses reads by child processes or through mmap. The fallback condition is judged only by the empirical D8 test, never by a document.
- Approved plan (summary row): label raising in the kernel is unreliable on Fedora 45 or through overlayfs.
Governing wording, recorded before the spike ran (FD-I): **the concept's full
wording governs both criteria, and no triggered criterion may be overridden
(PD-7). Recorded 2026-10-23 - after the first aarch64 spike runs, before the
x86_64 CI run, and before any gate evaluation - in
[`docs/adr/2026-10-23_Phase0-Founder-Decisions_v1.md`](adr/2026-10-23_Phase0-Founder-Decisions_v1.md).
Consequence recorded honestly: criterion 1's fallback clause ("and the fanotify
fallback misses reads by child processes or through mmap") is judged only by the
empirical D8 test, which has not run; D8 becomes a test, not a document, if any
gated row fails.**
Result: [TRIGGERED or NOT TRIGGERED]

## Kill criterion 2: interviews
From the synthesis, table C3: N = [ ] (reduced sample: [yes or no]); P1-Linux-ready (gating by default) = [ ]; P1-strict (reported) = [ ]; prong 2 gating reading (D + X + F) = [ ] of [ ] (denominator 1) and [ ] of [ ] (denominator 2); designed-workflow reading (D + X) = [ ] and [ ]; exposure check = [same or different].
Coding rules frozen and published on [date]: [link].
Result: [TRIGGERED or NOT TRIGGERED]

## Other Phase 0 acceptance criteria
| Criterion | Met? | Evidence or reason |
|---|---|---|
| AC4 Overlay decision in an ADR | [ ] | [ ] |
| AC5 Fallback documented (and tested, if required) | [ ] | [ ] |
| AC6 Signed image and cosign verify in CI | [ ] | [ ] |
| AC7 Synthesis published | [ ] | [ ] |

## Capacity
- Hours spent in Phase 0: [ ] of 32.75 estimated (about 30 budgeted).
- Weeks below 6 hours: [ ] (founder-capacity kill criterion: below 6 hours a week for 4 consecutive weeks).

## Options considered
1. GO: start Phase 1 on Monday 9 November 2026 as planned. Available only if neither criterion triggered.
2. GO WITH CHANGES: start Phase 1 with the changes listed below. Available only if neither criterion triggered.
3. PIVOT TO AUDIT TOOL (kill criterion 1): drop the kernel-label claim; ship the gateway, broker and receipts as an audit tool.
4. PIVOT TO BENCHMARK-FIRST (kill criterion 2): re-scope to a benchmark-first project for security researchers before Phase 1.
5. STOP.
If both criteria trigger, state how options 3, 4 and 5 combine.

## Decision
[Option number and name]

## Rationale
[Three to five sentences tied to the evidence above.]

## Overrides
None permitted. A triggered criterion cannot be overridden (synthesis template FD-6 for criterion 2; the same rule is proposed for criterion 1 under FD-I). If a criterion triggered, confirm that the decision is option 3 or 5 for criterion 1, and option 4 or 5 for criterion 2: [confirmation, or "No criterion triggered"].

## Changes for Phase 1
- [For example: kernel floor stated as 7.2; planner/worker experiment moved earlier; default labels changed by interview evidence.]

## What would change this decision
- [ ]

## Publication
- Devlog 3: [link]
- Next checkpoint: end of week 9, Sunday 20 December 2026 (label-creep kill criterion).
```

## 7. Hours summary
| Section | Hours | Building? |
|---|---|---|
| A. Governance | 0.75 | No |
| B. Threat model, methodology, specifications | 1.75 | No |
| C. ADRs | 0.50 | No |
| D. Kernel spike | 12.50 | Yes |
| E. Image pipeline | 5.00 | Yes |
| F. Interviews | 10.50 | No |
| G. Reviewers | 0.75 | No |
| H. Devlogs | 0.50 | No |
| I. Gate decision | 0.50 | No |
| **Total** | **32.75** | 17.50 building, 15.25 other |

| Week | Planned hours | If D2 and D7 are deferred (FD-H) |
|---|---|---|
| Week 1 (19 to 25 Oct) | 11.20 | 11.20 |
| Week 2 (26 Oct to 1 Nov) | 11.20 | 11.20 |
| Week 3 (2 to 8 Nov) | 10.35 | 9.35 |
| **Total** | **32.75** | **31.75** |

No contingency buffer is included. The contingent D8 test (about 3.0 hours) and the deferred checks in S5 are outside these totals.

## 8. Dependencies and founder decisions
**Every row below was decided on 2026-10-23 and is recorded with its reasoning in
[`docs/adr/2026-10-23_Phase0-Founder-Decisions_v1.md`](adr/2026-10-23_Phase0-Founder-Decisions_v1.md).**
The mapping is FD-A to PD-9, FD-B and FD-H to PD-6, FD-C to PD-10, FD-D to PD-11,
FD-E to PD-12, FD-F to PD-13, FD-G to PD-2, FD-I to PD-7 and FD-J to PD-8; the
labeller's lookup method became PD-1, its licence question PD-3, the aliasing
question PD-4 and the exec-only question PD-5. The recommendation column below
states what was chosen, except where the decision ADR records a change.

| ID | Decision | Why it matters | Recommendation (decided 2026-10-23) |
|---|---|---|---|
| FD-A | Install Go and podman on the macOS host, or keep all Phase 0 builds inside the Lima VM and CI | Neither is installed on the host; the plan lists them as Homebrew installs that need approval | Keep Phase 0 builds in the Lima VM and CI; no host install is needed except possibly for FD-B |
| FD-B | How to obtain an aarch64 fedora-bootc disk image for D2 | Phase 0 builds only x86_64 images; the cidata check needs aarch64 | Defer D2 to Phase 1 (FD-H); otherwise build the image on a GitHub arm64 runner as a CI artifact |
| FD-C | Create the GitHub repository (owner and name) | E2, E4, D6 and AC6 need it; the module path placeholder changes with it | Create it in week 1 |
| FD-D | DCO or CLA | Required before the first outside contribution | DCO, as the plan states |
| FD-E | Interview coding rules FD-1 to FD-6 | Must be frozen and published before P01 | Accept the defaults in the synthesis template, section A.4, or amend them, by Monday 19 October, before P01; publish the freeze |
| FD-F | Screener form tool; no cash incentive | Data minimisation and local-only handling | A self-hosted or privacy-respecting form; no incentive |
| FD-G | Kernel floor wording | The plan says "kernel 7.1 or later"; `bpf_real_data_inode` needs 7.2 | State 7.2 as the v0.1 floor, as the labeller ADR proposes. Revisit 7.1 only if C2 selects the no-cache design or the spike shows that a 7.1 path is needed; otherwise 7.1 support waits for Phase 6 (portable mode) |
| FD-H | Hours: accept 32.75 hours (2.75 over budget, no buffer, 17.5 building hours instead of 21), or defer D2 and D7 (1.0 hour) to Phase 1 week 4 now, for 31.75 hours | Interviews and their records need about 10.5 hours; D2 and D7 are records, not gate criteria | Defer D2 and D7 now (this also removes the FD-B dependency), accept the remaining 1.75-hour overrun, and track actual hours weekly |
| FD-I | Which wording of kill criterion 1 governs | The approved plan's summary triggers on unreliable kernel raising alone; the concept adds "and the fanotify fallback misses reads by child processes or through mmap". The plan's summary of criterion 2 also omits the second prong and says "of 10" | Record before the spike runs (by Sunday 25 October). Recommendation: the concept's full wording governs both criteria, because the plan rows are summaries; the fallback condition is judged only by the empirical D8 test; and no triggered criterion is overridden |
| FD-J | Defer gosec, govulncheck, CodeQL and the review list of root operations to Phase 1 | The plan lists them for every phase; Phase 0 has only the spike loader | Defer to Phase 1 week 4 and record the departure in devlog 2. Adding govulncheck and gosec beside staticcheck in E4 now would cost about 0.25 hours |

## 9. Proposed slip plan within Phase 0 (not approved; founder decision)
If the plan is more than 3 hours behind at the end of week 2 (Sunday 1 November), or if D8 becomes an empirical test:
1. Hold interviews at 8.
2. Move D7 (virtiofs and 9p behaviour) to Phase 1 week 4, unless already deferred under FD-H; macOS shared-folder users stay outside the v0.1 claims until it is recorded.
3. Move D2 (Lima cidata on fedora-bootc) to Phase 1 week 4, unless already deferred under FD-H; use UTM or QEMU with HVF for the boot loop, as the concept allows.
4. Merge devlog 2 into devlog 3.
5. Move B3, B4 and the ADR reviews C3 to C5 (0.65 hours) to Phase 1 week 4; they do not affect the gate.

Do not slip AC1, AC2, AC3, AC6, AC7 or AC8; the gate depends on them. If any of them cannot be met by 8 November, the decision record says so, and the founder either decides on the available evidence or moves the gate by at most one week, recorded publicly.

## 10. Kernel and base-image facts verified for the spike (16 October 2026)
| Fact | Source |
|---|---|
| Landlock ABI is 9 in Linux v7.1, 10 in v7.2 and 11 in the current mainline (7.3-rc7) | `security/landlock/syscalls.c` at tags v7.1 and v7.2 and on master, torvalds/linux |
| `bpf_lsm_file_open` and `bpf_lsm_file_receive` are in the `sleepable_lsm_hooks` set | `kernel/bpf/bpf_lsm.c`, master |
| `bpf_get_file_xattr` (KF_SLEEPABLE), `bpf_path_d_path` and `bpf_real_data_inode` (KF_SLEEPABLE, KF_RET_NULL) are registered in the fs kfunc set | `fs/bpf_fs_kfuncs.c`, master |
| `bpf_real_data_inode` is absent at v7.0 and v7.1 and present at v7.2; it returns the inode hosting a regular file's data (the upper or lower inode on overlayfs) | `fs/bpf_fs_kfuncs.c` at tags v7.0, v7.1 and v7.2 |
| `FAN_OPEN_PERM`, `FAN_ACCESS_PERM` and `FAN_PRE_ACCESS` are defined in the fanotify UAPI header (relevant to the contingent D8 test) | `include/uapi/linux/fanotify.h`, master |
| Fedora 45 kernel 7.2.9 stable since 6 October 2026 (7.2.10 in testing); Fedora 44 kernel 7.2.10 stable since 14 October 2026 | Fedora Bodhi updates API |
| `quay.io/fedora/fedora-bootc:45` is a manifest list, rebuilt daily (different digests on 14, 15 and 16 October). Digest observed on 16 October 2026 at 11:30 UTC: `sha256:86cc9d98ef210b71a6c4eeec4b7d04e7fef0a4415184c2322a14c879df189f45`. The base-image ADR is the authoritative record. | quay.io repository tag API |
