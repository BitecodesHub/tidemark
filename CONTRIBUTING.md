# Contributing to AgentOS Tidemark

> **Status: Phase 0 draft, 16 October 2026.** The project has no code yet. This document sets the
> rules that will apply from the first outside contribution. The founder confirms the sign-off model
> (DCO, proposed in this draft) before that contribution is accepted.

Thank you for your interest. Tidemark makes a security claim, so the contribution rules are
stricter than usual in a few places. They exist to protect users, not to keep people out.

## What helps most right now

During Phase 0 (19 October to 8 November 2026), the most useful contributions are:

- review of the [threat model](docs/THREAT_MODEL.md), the
  [benchmark methodology](docs/BENCHMARK_METHODOLOGY.md) and the specifications under `spec/`;
- exfiltration channels or bypass ideas the design does not cover (report anything that may be a
  vulnerability privately, as described in [SECURITY.md](SECURITY.md));
- evidence about how people actually run agents over private data (issues, not personal data);
- corrections to factual claims in the documents, with a primary source.

Large code contributions are unlikely to be merged before the Phase 0 go or no-go decision.

## Ground rules

- Follow the [Code of Conduct](CODE_OF_CONDUCT.md).
- Never report a vulnerability in a public issue or pull request. Use [SECURITY.md](SECURITY.md).
- Never include real secrets, real personal data or real canaries in issues, fixtures or logs.
- Contributions that add telemetry, phone-home behaviour or network calls outside the documented
  exits are not accepted.
- Keep claims inside the threat model. Documentation must not describe a feature as implemented or
  tested unless CI shows it. Every measured result in the README links to the CI artifact that
  produced it, and targets are labelled as targets.

## Developer Certificate of Origin (DCO)

Every commit must carry a sign-off certifying the
[Developer Certificate of Origin, version 1.1](https://developercertificate.org/):

```
Signed-off-by: Your Name <you@example.org>
```

Add it with `git commit -s`. The name must be your real name or a known identity you use
consistently; anonymous contributions cannot be accepted. A DCO check, planned for Phase 0 CI, will
reject pull requests with unsigned commits.

By signing off, you certify that you have the right to submit the work under the project's licence.
Contributions are accepted under the [Apache License, Version 2.0](LICENSE) (inbound licence equals
outbound licence, as in Section 5 of the licence). This draft proposes no contributor licence
agreement and a commitment not to relicense contributions. Both are pending the founder's decision
on DCO or CLA (see [GOVERNANCE.md](GOVERNANCE.md)).

## AI-assisted contributions

AI coding tools are allowed, and the founder uses them. They must be disclosed.

1. **Say so in the pull request.** Name the tools and describe what they produced (for example,
   "first draft of the parser and its tests").
2. **Add a trailer** to each commit that contains substantial AI-generated content:

   ```
   Assisted-by: <tool name and version>
   ```

   A `Co-authored-by:` trailer that names the AI tool also satisfies this rule. It does not make
   the tool an author or a DCO signatory.

3. **A human signs off.** An AI tool cannot certify the DCO. The person who signs off has read every
   line, understands it, and is responsible for it as if they had written it.
4. **Verify claims.** Generated code that calls kernel interfaces, kfuncs, library functions or
   GitHub Actions must be checked against current upstream sources. Note the source in the pull
   request.
5. Undisclosed AI-generated contributions, once identified, may be reverted.

The same rules apply to the founder's own commits.

## Review requirements

### Two-person review for security-critical areas

Changes in the following areas require review by **two humans**: the human submitter, who has
reviewed the diff line by line, and at least one other approving reviewer. An AI tool never counts
as a reviewer.

| Area | Paths (planned) |
|---|---|
| BPF programs and their loaders | `bpf/`, the Phase 0 spike (`spike/bpf/`, `spike/cmd/`), generated loader code, label-engine code in `pkg/labels/` |
| Gateway copy loop and clearance checks | `cmd/tidemark-gw/`, `internal/egress/`, `pkg/clearance/` |
| Model broker placement and field filtering | `cmd/tidemark-broker/`, `internal/broker/` |
| Staging merge, commit, accept and abort | `internal/stage/`, `internal/stateview/` |
| Ledger, receipts and verifier | `cmd/tidemark-ledgerd/`, `cmd/tidemark-verify/`, `internal/receipt/`, `spec/` |
| Root control plane | `cmd/tidemarkd/`, `internal/session/` |
| Build scripts and release pipeline | `image/`, `packaging/`, `.github/workflows/`, Makefiles and scripts |
| Test fixtures and benchmark definitions | `test/`, `bench/prereg/`, `bench/scenarios/`, `bench/baselines/`, the Phase 0 spike scenarios (`spike/tests/`) |

Test fixtures and benchmark definitions are on this list because weakening a test is the easiest way
to hide a regression.

**Honest interim rule.** The project has one maintainer today, so a second reviewer is not always
available. Until named reviewers are confirmed, changes in these areas may merge with the founder's
line-by-line review only, labelled `needs-second-review`. The v0.1 release gate requires a second
review of every `needs-second-review` change in the BPF programs, the gateway copy loop, the
staging merge and the ledger. Changes in the other listed areas get a second review where reviewer
capacity allows; any left unreviewed are listed in the release notes. This interim rule ends once
two reviewers are confirmed.

### All other changes

Changes from outside contributors need one approving review from the maintainer. Until a second
maintainer or committer exists, the maintainer's own changes outside the security-critical areas
merge through a pull request once CI and the DCO check pass. Branch protection on `main` requires
pull requests and passing status checks, with zero required approvals until a second reviewer
exists. Every change needs a DCO sign-off on every commit.

## Commit rights

Commit rights are granted only after **six months of verified contribution**: a sustained record of
reviewed, merged work and constructive review of other people's changes, by an identity the project
can verify. Rights are granted by the founder, require two-factor authentication on the account, and
can be withdrawn after long inactivity or a breach of these rules. Commit rights for the areas that
need two-person review are granted separately. See [GOVERNANCE.md](GOVERNANCE.md).

## Pull request checklist

- [ ] Every commit has a `Signed-off-by:` line.
- [ ] AI assistance is disclosed in the description and in `Assisted-by:` trailers (or
      `Co-authored-by:` trailers that name the tool).
- [ ] The change has tests, or the description explains why it cannot be tested.
- [ ] A confirmed bypass fix includes a regression test or benchmark scenario.
- [ ] New dependencies are justified, pinned and licence-compatible with Apache-2.0.
- [ ] GitHub Actions are pinned by full commit SHA; container images are pinned by digest.
- [ ] Documentation and the threat model are updated if behaviour or guarantees change.
- [ ] The change does not widen a claim beyond what the tests measure.

## Style

- **Go:** `gofmt`, `go vet` and `staticcheck` clean; `gosec` and `govulncheck` findings addressed
  or explained. The Go module path is `github.com/BitecodesHub/tidemark` (set 2026-10-23,
  decisions ADR PD-10).
- **BPF C:** bounded loops, early exit for processes outside sessions, and every kfunc failure
  handled. Generated files (`vmlinux.h`, `*_bpfel.go`, `*_bpfeb.go`, `*.o`) are not committed.
- **Shell:** `shellcheck` clean. **Workflows:** `actionlint` clean.
- **Formatting:** follow [.editorconfig](.editorconfig).
- **Documentation:** formal English, active voice, no contractions and no emojis. New dated
  documents use the file name format `YYYY-MM-DD_Subject_Version.md`, for example
  `docs/adr/2026-10-16_BaseImage_v1.md`. Conventional files keep their conventional names.
- **Commit messages:** an imperative summary line of 72 characters or fewer, a blank line, then the
  reason for the change.

## Benchmark scenarios

A new adversarial scenario counts only if it leaks in the audit or unsandboxed configuration, which
proves the channel is real. Scenarios must not contain canaries; the host generates one per run.
Baseline configurations must use each tool's strongest documented configuration. See
[docs/BENCHMARK_METHODOLOGY.md](docs/BENCHMARK_METHODOLOGY.md).

## Questions

Open a GitHub discussion or issue for anything that is not a security report.
