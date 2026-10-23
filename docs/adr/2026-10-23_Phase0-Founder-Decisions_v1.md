# Phase 0 founder decisions: the open list, closed

- **Status:** Accepted (recorded under the founder's explicit delegation; see Deciders)
- **Date:** 2026-10-23
- **Phase:** Phase 0 (kernel spike and gate week 3)
- **Deciders:** The founder delegated every open Phase 0 decision to the Qoder
  agent on 2026-10-23 ("you decide everything and stop when done"). The agent
  made each choice and is the author of this record; the founder has not
  reviewed it line by line. That delegation is the only reason this ADR carries
  **Accepted** rather than **Proposed**, and it is stated here so no reader
  mistakes the status for a founder review. Any decision below may be reversed
  by the founder in a later ADR; the evidence cited is reproducible.

## Context

Phase 0 leaves a set of decisions that block other work. They are listed in
[../2026-10-16_Phase0-Checklist_v1.md](../2026-10-16_Phase0-Checklist_v1.md)
section 8 as FD-A to FD-J, in the synthesis template section A.4 as FD-1 to
FD-6, and inside the labeller ADR as the lookup-method question (C2), the
hard-link/rename aliasing question, and the BPF licence question.

The naming in this record avoids both existing ID spaces: `PD-n` (Phase 0
decision). The FD ids it replaces are named in each row.

New evidence became available on 2026-10-23 and changes several answers, which
is why these decisions are taken now rather than in week 1. The 16-scenario
aarch64 spike passes 16 of 16 in three consecutive runs, and two configuration
experiments were run inside that suite:

- the BPF object rebuilt with `SEC("license") = "Dual BSD/GPL"` still has its
  kfunc calls accepted by the verifier (16/16), so dual-licensing costs nothing;
- the loader switch `TM_DISABLE_INODE_CACHE=1` compiles out every inode-label
  cache read and write, and the whole suite - including the overlayfs scenario -
  still passes 16/16 without `bpf_real_data_inode` being used.

The cache experiment also produced the decisive negative result. With the cache
on, removing a file's `security.bpf.tidemark` attribute leaves the label raised
on the next open (`label=0x3 cause=open|cache`), because the cache is a floor
and only the per-inode generation counter in the ADR's invalidation design could
lift it. With the cache off, the same operation reads `label=0x0 cause=none`.
Evidence:
[../spike-results/2026-10-23_Kernel-Spike-Results-aarch64_v2.md](../spike-results/2026-10-23_Kernel-Spike-Results-aarch64_v2.md),
logs `2026-10-23_Spike-Run-aarch64_v1.log` through `_v3.log`,
`-dual-licence_v1.log` and `-no-cache_v1.log`.

## Decision

### Summary table

| ID | Closes | Decision |
|---|---|---|
| PD-1 | C2, checklist section C lookup method | The production labeller keeps **no inode-keyed label cache**. Labels are computed per open from the file's xattr plus matching directory rules. |
| PD-2 | FD-G | **Linux 7.2 or later, Landlock ABI 9 or later** stays the v0.1 floor. `bpf_real_data_inode` is removed from the required-kfunc list; 7.1 support remains deferred to Phase 6. |
| PD-3 | licence question in `spike/README.md` and the labeller ADR | BPF C sources are **`GPL-2.0-only OR BSD-2-Clause`**, SPDX line to match, kernel tag `char LICENSE[] SEC("license") = "Dual BSD/GPL"`. Everything else in the repository stays Apache-2.0. |
| PD-4 | hard-link/rename aliasing open question | A directory label follows the **location**, not the content. The Phase 1 labeller **denies** in-session `link` and `rename` that would move a path out of a ruled root (`inode_link`, `path_rename` return `-EPERM`) instead of stamping labels onto inodes. |
| PD-5 | exec-only opens | An `open()` used only to execute counts as a read and may raise the label. Over-labelling an execute-only binary is accepted; the design fails closed. |
| PD-6 | FD-H, FD-B | **D2 and D7 deferred to Phase 1 week 4.** Phase 0 runs at 31.75 planned hours, and the FD-B dependency (an aarch64 fedora-bootc disk image) disappears with D2. |
| PD-7 | FD-I | **The concept's full wording governs both kill criteria**, and no triggered criterion is overridden. |
| PD-8 | FD-J | gosec, govulncheck, CodeQL and the root-operations review list are **deferred to Phase 1 week 4**, recorded as a departure from the plan's "every phase" list. |
| PD-9 | FD-A | **No Go or podman on the macOS host.** All Phase 0 builds happen in the Lima VM and in CI. |
| PD-10 | FD-C | The repository is **created private under `BitecodesHub/tidemark`**, the Go module path is updated from the placeholder `github.com/agentos-tidemark/tidemark`, and moving it to an organisation or making it public is a launch decision, not a Phase 0 one. |
| PD-11 | FD-D | **DCO**, not a CLA. |
| PD-12 | FD-E | The interview coding rules **keep every default in synthesis template section A.4**, which are the stricter reading in each case; the freeze is published as the SHA-256 of the frozen file in devlog 1. |
| PD-13 | FD-F | The screener is a **self-hosted, privacy-respecting form with no incentive** offered. |
| PD-14 | found in CI | **Base-image digests are re-pinned by hand, and CI checks reachability before building.** The 16 October digests are gone from the registry, so a pin is an identity, not a retention guarantee. |
| PD-15 | found in CI | **OpenSSF Scorecard is non-gating while the repository is private**, and no long-lived personal token is added as a repository secret to make it work. |
| PD-16 | found in CI | **The x86_64 spike is driven over SSH** with a per-run ephemeral key; the serial log is boot and oops evidence only, and the source ISO is built from `git archive`. |
| PD-17 | found in CI, blocks AC6 | **Image publication is opt-in.** `push` to GHCR and keyless cosign signing run only while the repository variable `PUBLISH_IMAGES` is `push`; a passing build and boot test no longer publish by themselves. |
| PD-18 | E4, S2 | **Branch protection stays off while the repository is private and solo**, and secret scanning is not requested: the settings API returns 404 for this repository, so there is nothing to enable. DCO, gofmt, staticcheck, pinned Actions, SBOM and the report-only vulnerability scan are all in place. |

### PD-1: no inode cache

The labeller ADR allowed two lookup methods and gated the choice on C2: an
inode-keyed label cache plus per-inode generation counters, or reading the
attribute on every open. The ADR's own text already warned that the cache is a
floor that must never hide a fresh higher label, and that absence caching needs
a compare-and-swap generation counter.

This record goes further than "the cache is risky". The measured result is that
the cache as designed **breaks declassification**, and the cost of fixing it is
paid by every open: two extra LSM hooks (`inode_post_setxattr`,
`inode_post_removexattr`), a per-inode counter, a compare-and-swap on the read
path, and a bounded map with eviction policy that must fail closed when full.
Against that, the benefit measured on the only hot path that the cache touches
is small: an xattr-labelled open inside a session costs **1,211 ns at p50**, the
same as an unlabelled open in a session (also 1,211 ns), and the labeller's
total overhead over the outside-session baseline is **376 ns at p50**. A cache
cannot buy back a meaningful fraction of 376 ns.

So the production labeller computes the label per open from
`bpf_get_file_xattr` plus directory-rule matching and stores nothing per inode.

### PD-2: the kernel floor, restated

The floor was 7.2 because `bpf_real_data_inode` needs 7.2. PD-1 removes that
dependency, so 7.1 could in principle return. It does not: the sleepable LSM
kfuncs, `bpf_get_file_xattr` and `bpf_path_d_path` are all verified present on
7.2, and Phase 0 has no 7.1 evidence at all. Supporting an untested kernel is a
claim without proof, which is the thing this project is trying not to do.
**v0.1 requires Linux 7.2 or later and Landlock ABI 9 or later.** Phase 1 doctor
must check the kfuncs it actually uses, not `bpf_real_data_inode`.

### PD-3: BPF licence

The reviewer's recommendation was already in `spike/README.md`: the kernel
grants the BPF LSM kfuncs only to GPL-compatible objects, so Apache-2.0 alone
cannot carry the BPF C sources, and the checked-in object declared `"GPL"` while
its SPDX line said Apache-2.0. That contradiction is a launch-blocking defect in
a project whose whole product is trustworthy metadata. Dual licence
`GPL-2.0-only OR BSD-2-Clause` with the conventional `"Dual BSD/GPL"` tag
resolves it, keeps recipients free to take either side, and is the standard
pattern in kernel-adjacent code. It was not accepted on argument alone: the
rebuilt object passed the full suite.

`NOTICE` records the exception and the file list it applies to.

### PD-4: location, with link and rename denied

Scenario `spike/tests/12_link_rename_alias.sh` shows that a hard link or a
rename out of a labelled root inside the session's mount namespace gives the
file a path that no rule matches, so a directory-derived label silently
disappears. The two candidate fixes are content binding (stamp the xattr from
`inode_link`/`path_rename`, so the label follows the inode) and location
binding with enforcement (the label follows the path, and the session is denied
the operations that would break the promise).

Location wins, for three reasons. First, content binding makes the label
survive the very operation a user may intend: `mv ~/Private/report.txt
~/shared/` under content binding keeps a private label on a file the user has
deliberately published, and a label the user cannot remove by ordinary means is
not a label the user can reason about. Second, content binding requires the
labeller to write xattrs from BPF hooks on filesystems that may not support
them, on inodes owned by other users, which is a much larger authority than the
project's threat model grants. Third, directory rules are already path-based
objects in the policy language; a rule that means "everything under this
directory" is honoured exactly by refusing to let a session move content out
from under it.

The cost is real and is stated: inside a session, `ln` or `mv` out of a ruled
root fails with `EPERM` and the user sees a refusal they may not expect. The
gateway must surface which rule caused it. Phase 1 acceptance requires the
error path to name the rule.

### PD-5, PD-6, PD-7, PD-8

PD-5: `open()` with `O_RDONLY` and an `execve()` following it are
indistinguishable at the hook without tracking later execution, and the
conservative reading is that executing a file reads it. Over-labelling a binary
can block a session; under-labelling a private file can leak it. Fail closed.

PD-6: D2 (Lima cidata on fedora-bootc) and D7 (virtiofs and 9p xattr and
append-only behaviour) are records, not gate criteria, and neither is needed to
evaluate either kill criterion. Deferring them buys the FD-B dependency and
brings the plan inside the budget's spirit while still accepting a 1.75-hour
overrun. macOS shared-folder users stay outside every v0.1 claim until D7 is
recorded; that restriction is stated in the README.

PD-7: the concept is the document the public reads, so its wording is the
standard the project can be held to. The plan's shorter rows are summaries.
Adopting the concept's wording means criterion 1 also requires the fanotify
fallback to miss reads by child processes or through mmap before the failure
kills the project, which makes D8 an empirical test rather than a document if
any gated row fails (contingent ~3.0 hours from the slip plan). And a triggered
criterion cannot be overridden: the concept's action for a triggered criterion
is taken as the decision, so no option keeps Phase 1 as planned.

PD-8: Phase 0 ships no product code that runs as root and no network listener.
`spike/` is a probe. gosec, govulncheck and CodeQL earn their keep against
`tidemarkd` and the gateway, which do not exist yet. Deferring them is a
departure from the plan's "every phase" list and is recorded in devlog 2 rather
than hidden.

### PD-9 through PD-13

No host toolchain install (PD-9) because the spike already builds and runs in
the Lima VM and CI builds x86_64; adding Homebrew Go and podman to the M1 host
buys nothing and creates two more version-pinned things to keep honest.

PD-10 makes the repository exist so that D6, E2, E3 and E4 can run. Private is
the reversible choice: the CI, image pipeline and gate evidence get built where
they can be run, and publication remains the founder's launch decision. The
module path placeholder is not load-bearing and is updated to the real path.

PD-12 accepts the interview defaults because each default is the stricter
option; none of them makes kill criterion 2 easier to pass. FD-3 in particular
gates on "abandon a session", the criterion's own words, so a fresh session
counts as abandonment, and the designed-workflow reading is reported alongside
it rather than substituted for it. The freeze is published before P01 as the
SHA-256 of the synthesis template recorded in devlog 1, which is what F1 asks
for and does not require a public commit of the whole file.

### PD-14, PD-15 and PD-16: found in the first CI run

The first CI run on the new repository failed in all four workflows, and each
failure was infrastructure rather than the kernel mechanism, so each needed a
recorded call.

**PD-14, base-image digests.** The image pipeline pins
`quay.io/fedora/fedora-bootc` by the digest of its multi-architecture index. On
23 October the first real build reported `manifest unknown` for both 16 October
pins: index `sha256:86cc9d98…189f45` for tag 45 and `sha256:9db34195…9eeaced`
for tag 44. A pinned digest therefore identifies a compose; it does not keep it
retrievable, because quay.io drops a compose's manifests once the daily rebuild
has moved the tag. Decision: re-pin by hand to the digests resolved on
23 October (tag 45 index `sha256:4c714111…93929da`, version 45.20261023.0, amd64
kernel `7.2.11-300.fc45.x86_64`; tag 44 index `sha256:03808453…8e4d6d2e`, version
44.20261023.0, kernel `7.2.11-200.fc44.x86_64`, both read from the image config's
`ostree.linux` label), record every re-pin in the base-image ADR, and add a
pre-build check that resolves the tag, tests whether the pinned digest is still
retrievable, prints the tag's current digest, and fails with that name instead
of letting podman fail two minutes later. Rejected: build against the moving
tag, or re-pin automatically inside CI. Either would let an image labelled with
one digest contain another compose's content, which is the exact failure mode
the pin exists to prevent.

**PD-15, Scorecard.** Scorecard aborts on a private repository before analysis:
`scorecard had an error: internal error: ListCommits:error during
graphqlHandler.setup: githubv4.Query: Resource not accessible by integration`
(run 39965539746). The documented remedy is a personal access token with `repo`
scope stored as a repository secret. Decision: keep the workflow, run it on
schedule and by hand only, set `publish_results: false` and make the job
non-gating, and record the gap honestly rather than claim a supply-chain
score. Reasoning: minting a long-lived personal credential and storing it in the
repository is a change to the founder's account security posture with a far
wider blast radius than the missing metric, and two of Scorecard's three
outputs are unavailable on this plan anyway - publishing needs a public
repository, and code-scanning SARIF upload needs GitHub Advanced Security,
which this private plan does not have (the secret-scanning endpoints answer
404). Restore gating and publishing at launch, when the repository becomes
public.

Verified rather than assumed later the same day: dispatch run 39984799845
reproduced the identical `ListCommits` error, the job concluded `failure` and
the workflow run still concluded `success`. That is the whole of what
`continue-on-error: true` buys here - the gap stays visible in the job while
nothing that should gate is blocked by it. What it does not buy is a score, and
none is claimed anywhere in the repository.

**PD-16, how the x86_64 spike is driven.** The first version ran the whole
suite inside a cloud-init `runcmd` whose output was redirected to `/dev/ttyS0`,
and gated on a marker grepped from the serial log. Run 39965539646 left the
runner with `qemu exit code: 0`, no marker, and no explanation: while dnf was
still running, `serial-getty@ttyS0` restarted, its `vhangup()` replaced the
already-open stdout and stderr of the runcmd shell with the hung-up tty's file
operations, and every later write - the build log, the scenario results and the
result marker itself - returned `-EIO` and was discarded. The run's outcome was
unknowable, which is the one thing a gate may not produce. Decision: keep QEMU,
KVM and the pinned Beta compose, but drive the guest over SSH through QEMU
user-mode networking with a host forward; generate an ed25519 key pair per run
and put only the public half in the seed ISO; keep the serial log as an artifact
for boot failures and kernel oops, never as the verdict channel; and build the
source ISO from `git archive HEAD spike` rather than the working directory, so
the aarch64 build products that live in the local tree (`label_probe.bpf.o`,
`tidemark-probe`, `bpf/vmlinux.h`, whose mtimes are newer than the sources)
cannot be carried into an x86_64 guest and reused as if already built for that
architecture. Cost: the job now depends on the guest's sshd, which is present
and enabled in the Fedora Cloud image; if it ever does not answer, the job says
so and prints the serial tail. This is instrument methodology, not a product
change.

**PD-17, how the image pipeline publishes.** With the first green build of
`tidemark-image.yml`, the `publish` job became reachable: its guard was
`image`-path changes plus a push to `main`, and the matrix boot tests were
network-dependent, not part of the job's requirements beyond `needs`. On the
first run that completed the build and the Fedora 45 boot test, publication was
about to start on the Fedora 44 boot test finishing - which would have pushed
`ghcr.io/bitecodesHub/tidemark-bootc` images and, through keyless cosign,
written **public, permanent entries to the Sigstore Rekor transparency log**
naming this private, unlaunched repository and its layer digests. A Rekor entry
cannot be retracted; the log is mirrored and searched. The repository is private
under PD-10 precisely because the launch decision is the founder's, and this
side effect is the least reversible one in the pipeline.

Decision: the workflow still builds and boots on every change to the image
inputs, but publishing requires the repository variable `PUBLISH_IMAGES` to equal
`push`, and the publish job now also requires `needs.boot-test.result ==
'success'` explicitly rather than relying on `needs` alone. A new
`publish-decision` job runs on every image-path run and writes to the run summary
what it found - the boot-test result, the trigger, the current `PUBLISH_IMAGES`
value, and one of five named reasons for skipping - so "why did this not publish"
is answered by the log rather than by reading the workflow. Enabling publication
is one command (`gh variable set PUBLISH_IMAGES --body push`) and is the
founder's call, made when the repository goes public or when a consumer needs a
pulled image. AC6 is therefore **deferred by one deliberate command**, not met. The guard has since been observed working rather than assumed: run 39982982294 completed the whole graph - both builds, both boot tests green, including the Fedora 44 fallback's first CI boot - and `sign and publish` came out `skipped`, with the decision job reporting `BOOT=success`, `EVENT=push`, `REF=refs/heads/main`, `PUBLISH_IMAGES` unset.

**PD-18, branch protection and secret scanning.** E4 asks for both, and neither
is a code change, so the honest reading of "who owns this" is a decision rather
than a task. Branch protection with required status checks is turned off on
purpose: Phase 0 is one person pushing documentation, evidence transcripts and
spike fixes straight to `main`, often in the same hour a CI run finishes, and a
required-check rule would make every one of those commits wait on a merge
through a PR whose only reviewer would be the same person. It is also the one
E4 item that can be added later with no loss: a single `gh api` call that PUTs
protection on `main` requiring `lint`, `kernel-spike` and `dco` takes a minute,
and the right moment is the first outside contributor or the move to public,
whichever comes first. Secret scanning
and push protection are not available: both endpoints answer HTTP 404 for this
private repository on this plan, so the item cannot be ticked by clicking
anything. What the plan can and does get now is the parts that are code:
`Signed-off-by` enforcement, gofmt and staticcheck, every Action pinned by
commit SHA, a per-build SBOM and a report-only vulnerability scan, and the
dependency graph with vulnerability alerts enabled. E4 and S2 stay open on their
named items rather than being reworded until they pass.

## Consequences

**Becomes easier.** The labeller ADR's kernel-floor dependency drops out; C2 and
AC4 have a recorded answer; the licence contradiction is gone, with a passing
build as proof; declassification by removing an attribute works with no extra
Phase 1 machinery; the aliasing hole has an enforcement answer that fits the
existing policy language; the Phase 0 plan fits inside 31.75 hours; and the
repository exists, so AC2, AC6, D6, E2, E3 and E4 are no longer blocked.

**Becomes harder.** Every labelled open now reads the attribute through
`bpf_get_file_xattr`; if a Phase 1 benchmark shows the per-open read is the
dominant cost on a real workload, PD-1 is reopened with that data and the
invalidation design must be built and paid for then. `ln`/`mv` out of a ruled
root failing with `EPERM` is a user-visible behaviour that needs wording, error
reporting, and a rule identifier in the refusal. Phase 0 carries no security
scanner beyond gofmt and staticcheck, so Phase 1 week 4 has a named block of
work rather than a trickle. A private repository means the devlog's "everything
public" default is deferred, not abandoned, and the gate record's publication
step still has to happen.

**Risks accepted.** The latency baseline is measured with the probe loaded and
in a single Go process, so it brackets `open(2)` only and understates end-to-end
session cost; it is a floor, not a product number. PD-4's enforcement depends on
`inode_link` and `path_rename` being attachable and reliable under mount
namespaces, which Phase 0 does not test; if Phase 1 finds those hooks cannot see
the operation, the fallback is to deny the operation by other means or to
narrow the directory-rule claim, and either outcome is a documented limitation
rather than a silent hole. Kill criterion 1 still needs the x86_64 half.

**Amended later the same day.** The x86_64 half has now run (PD-16): 16 of 16 in
four consecutive CI runs on the same scripts, and a fifth green run after them on
unchanged test code, reported in
`docs/spike-results/2026-10-23_Kernel-Spike-Results-x86_64_v1.md`, so the
"either architecture" condition of criterion 1 is measured rather than assumed.
The latency caveat also sharpened with the second architecture: the labeller
delta is about +0.38 µs p50 on the aarch64 Lima VM but about +1.56 µs p50 on the
CI guest, which QEMU gives 4 vCPU on a runner whose KVM recommends 2. Both are
floors measured in one process around `open(2)`, and neither is a product
number. The `inode_link`/`path_rename` question above is unchanged: still
untested in Phase 0. One further decision came from the same CI session and is
recorded here rather than in a second ADR: PD-17, publication is opt-in, after
keyless signing was found to be one boot test away from writing an irreversible
public transparency-log entry for a private repository.

## Alternatives considered

- **Keep the inode cache with generation counters (PD-1).** Rejected: measured
  to break declassification as designed, to add two LSM hooks plus a CAS
  read path plus eviction policy, and to buy back at most a fraction of a 376 ns
  overhead that the xattr read itself contributes nothing measurable to.
- **Keep Apache-2.0 for BPF sources and rely on the `"GPL"` tag (PD-3).**
  Rejected: the SPDX line and the kernel tag disagree, which is exactly the kind
  of metadata that this project promises to make trustworthy.
- **Make BPF sources GPL-only (PD-3).** Rejected: it removes a recipient's
  option to take the object under a permissive licence and complicates the
  project's own Apache-2.0 distribution story, for no gain over the standard
  dual licence.
- **Bind directory labels to content by stamping xattrs from BPF (PD-4).**
  Rejected for the three reasons in PD-4: contrary to user intent on publish,
  larger write authority than the threat model grants, and a policy language
  that is already path-based.
- **Keep the full 32.75-hour Phase 0 plan (PD-6).** Rejected: it needs an
  aarch64 bootc image for no gate criterion, and the overrun buys records the
  gate does not read.
- **Adopt the plan's shorter wording of criterion 1 (PD-7).** Rejected: it would
  let the project claim a gate it never tested, since the fallback condition is
  part of what the public concept promises.
- **Add gosec/govulncheck/CodeQL now (PD-8).** Rejected: there is no product
  code for them to review, and ~0.25 hours spent now scans a spike loader.
- **Install Go and podman on the host (PD-9).** Rejected: no Phase 0 artifact
  requires them outside the VM.
- **Create the repository public immediately (PD-10).** Rejected for now: the
  launch decision is the founder's and publication is irreversible; private
  unblocks every CI-dependent criterion and can be flipped at any time without
  changing history.
- **Pay participants (PD-13).** Rejected: an incentive biases the sample that
  kill criterion 2 is built on.
- **Gate FD-3 on the designed-workflow reading (PD-12).** Rejected: it is the
  looser option, and the freeze must not be chosen to make a criterion pass.
- **Build `FROM quay.io/fedora/fedora-bootc:45` by tag, or re-pin it in CI
  (PD-14).** Rejected: the digest is the identity that makes "this commit built
  this image" mean something, and an automatic refresh would let an image
  labelled with one digest carry another compose's content while pretending
  otherwise.
- **Store a personal access token as a repository secret to make Scorecard work
  on the private repository (PD-15).** Rejected: it converts a short-lived,
  least-privilege CI credential into a long-lived one on the founder's account,
  for a metric whose two public outputs (badge and code-scanning SARIF) are not
  available on this plan anyway.
- **Keep the x86_64 spike on the serial console and hope the marker survives
  (PD-16).** Rejected after measuring why it did not: a getty restart silently
  swallows the rest of the transcript, so the gate would sometimes be unable to
  tell a pass from a lost log.
- **Ship the x86_64 spike from the local working tree as before (PD-16).**
  Rejected: `git archive` is the only way to guarantee the guest builds from
  sources, because the local tree carries aarch64 objects newer than every
  source file.
- **Let the first green build publish to GHCR and sign automatically (PD-17).**
  Rejected after reading what keyless signing does: it writes a public,
  permanent Rekor entry naming the repository and its digests, which is
  irreversible, visible to anyone searching the log, and inconsistent with a
  private, unlaunched repository whose publication is an open founder decision.
  Publishing stays one variable away, so nothing is lost but the surprise.
- **Turn on branch protection with required checks today (PD-18).** Rejected for
  Phase 0: it would gate every evidence commit behind a self-review on a solo
  repository, and it is fully recoverable in a minute at the moment collaborators
  appear. The compensating controls that matter for a spike - pinned Actions,
  sign-off checks, lint and a report-only scan - are all already running.

## References

- [2026-10-16_Labeller-BPF-LSM_v1.md](2026-10-16_Labeller-BPF-LSM_v1.md) - the
  ADR PD-1, PD-2, PD-4 and PD-5 amend (revised note of 2026-10-23).
- [2026-10-16_BaseImage_v1.md](2026-10-16_BaseImage_v1.md) - kernel floor
  wording, unchanged by PD-2 in substance.
- [../2026-10-16_Phase0-Checklist_v1.md](../2026-10-16_Phase0-Checklist_v1.md) -
  section 8 FD-A to FD-J, section 4 AC1-AC8, section 6 decision template.
- [../spike-results/2026-10-23_Kernel-Spike-Results-aarch64_v2.md](../spike-results/2026-10-23_Kernel-Spike-Results-aarch64_v2.md) -
  the 16/16 runs, the two configuration experiments, the latency table.
- [../spike-results/2026-10-23_Spike-Environment-aarch64_v2.txt](../spike-results/2026-10-23_Spike-Environment-aarch64_v2.txt) -
  kernel, toolchain and object digest for the runs.
- [../spike-results/2026-10-23_Kernel-Spike-Results-x86_64_v1.md](../spike-results/2026-10-23_Kernel-Spike-Results-x86_64_v1.md) -
  the five CI runs that closed the x86_64 half of criterion 1, and the latency
  comparison that amended the caveat above.
- [../../.github/workflows/build-image.yml](../../.github/workflows/build-image.yml) -
  where PD-14's reachability check and PD-17's `PUBLISH_IMAGES` guard live.
- [../../spike/README.md](../../spike/README.md) - scenario-to-assumption
  mapping, licence note, `TM_DISABLE_INODE_CACHE`, `# TIMEOUT:`.
- [../interviews/2026-10-19_Synthesis-Template_v1.md](../interviews/2026-10-19_Synthesis-Template_v1.md) -
  section A.4 interview rules FD-1 to FD-6, accepted unchanged under PD-12.
- [../devlog/2026-10-16_Devlog_v1.md](../devlog/2026-10-16_Devlog_v1.md) - where
  the freeze SHA is published.
