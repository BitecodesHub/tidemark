# Governance

> **Status: Phase 0 draft, 16 October 2026.** This document describes how AgentOS Tidemark (working
> name) is governed today. It will change as the project changes, through the process described
> below.

## Summary

- Tidemark is run by one founder and maintainer, who makes final decisions (a "benevolent dictator
  for life", or BDFL, model).
- **The bus factor is 1.** If the founder stops, the project stops unless someone forks it. This
  document says so plainly and explains what happens in that case.
- Significant decisions are recorded as architecture decision records (ADRs).
- The licence is Apache-2.0. A commitment never to relicense is proposed, pending the founder's
  decision on DCO or CLA (an open question in the project plan).

## Roles

| Role | Who | Responsibilities |
|---|---|---|
| Founder and maintainer | One person, working about 10 hours a week | Final decisions, releases, security triage, roadmap, granting commit rights |
| Security reviewers | Two or three named people, recruited in Phase 0 with their consent recorded (none confirmed yet) | Second review of security-critical changes; private attack on the image and benchmark before launch |
| Committers | None yet | Merge rights in agreed areas, granted only after six months of verified contribution (see [CONTRIBUTING.md](CONTRIBUTING.md)) |
| Conduct alternate | One independent person, named with consent (none yet) | Receives and decides conduct reports that concern the maintainer (see [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)) |
| Contributors | Anyone | Issues, reviews, scenarios, documentation and code under the sign-off model in [CONTRIBUTING.md](CONTRIBUTING.md) (DCO proposed) |

A second maintainer would be welcome. It is a bonus, not a plan the project depends on.

## Honest constraints

- The project is part-time and solo. Response times in [SECURITY.md](SECURITY.md) reflect that.
- Most code is written with AI assistance. The founder reviews every security-critical diff line by
  line, and named reviewers provide a second review of the BPF programs, the gateway copy loop, the
  staging merge and the ledger before launch (see [CONTRIBUTING.md](CONTRIBUTING.md)).
- The project has public kill criteria in the
  [project plan](docs/2026-10-16_Project-Plan_v1.md) and the
  [concept document](docs/research/2026-10-16_Tidemark-Concept_v1.md). Hitting one can narrow the
  claim, change the scope or stop the project. Each such decision is recorded in an ADR.

## How decisions are made

1. **Routine changes** are decided in issues and pull requests. Silence after a reasonable review
   period counts as agreement (lazy consensus), but the founder decides when there is disagreement.
2. **Significant decisions** are recorded as ADRs in `docs/adr/`, named
   `YYYY-MM-DD_Subject_Version.md` (for example `docs/adr/2026-10-16_BaseImage_v1.md`). An ADR is
   required for:
   - any change to a guarantee (G1 to G9), a non-guarantee or the threat model's trust boundaries;
   - any change to the policy or receipt specifications;
   - kernel interface choices, including the labeller design and its fallback;
   - the base image, language, packaging and delivery choices;
   - benchmark methodology changes after pre-registration (logged as deviations);
   - go or no-go decisions at phase gates and kill-criterion outcomes;
   - changes to this document, [CONTRIBUTING.md](CONTRIBUTING.md), [SECURITY.md](SECURITY.md) or
     [TRADEMARK.md](TRADEMARK.md).
3. Each ADR states its context, the decision, the alternatives considered and the consequences. Its
   status is one of: Proposed, Accepted, Rejected or Superseded. Accepted ADRs are not edited;
   a new ADR supersedes them.
4. Proposed ADRs that change a guarantee or a governance document stay open for public comment for
   at least 7 days before acceptance, except for urgent security fixes.
5. Until the repository is first made public, the founder may edit governance drafts directly.
   Filling a placeholder contact is a routine change.

## Licence

The project's code, specifications, benchmark and documentation are licensed under the
[Apache License, Version 2.0](LICENSE), except where a file states otherwise. Quotations and
excerpts from third-party sources in `docs/research/` remain under their owners' terms.

**Proposed commitment, pending the founder's decision on DCO or CLA (plan open question).**

- The project will not relicense to a source-available, non-commercial or more restrictive
  licence, and will not adopt an "open core" model in which guarantees depend on proprietary code.
- The project will not use a contributor licence agreement that would allow relicensing.
  Contributions are accepted under Apache-2.0 with a DCO sign-off, so each contributor keeps their
  copyright.

## Conflicts of interest

Maintainers and reviewers disclose relevant conflicts of interest in the affected issue or pull
request, for example employment by a vendor whose product is a benchmark baseline. Each baseline
project receives the benchmark methodology at least 7 days before results are published, and its
reply is published with them.

## Succession and continuity

With a bus factor of 1, continuity depends on what is written down now.

- **Successor contact.** `[SUCCESSOR CONTACT: to be named by the founder, with that person's
  consent]`. The successor holds offline recovery material and two-factor authentication, is listed
  publicly, and uses that material only under the conditions below. The successor has no standing
  access to the GitHub organisation or the package repositories; the recovery mechanism is a
  founder decision, recorded in an ADR.
- **Inactivity.** If the founder does not respond to an open security report for 30 days past its
  disclosure deadline, or makes no project activity for 6 months, the successor contact may:
  - add a notice to the README stating that the project is unmaintained and that its guarantees
    are no longer being defended;
  - publish any pending security advisories with the information available;
  - archive the repository, or transfer it to a maintainer or a neutral organisation willing to
    take it on.
- **Signing.** Releases use keyless cosign signing tied to the repository's workflow identity, so
  there is no long-lived project signing key to hand over. A transfer or rename changes the
  identity users pass to `cosign verify`; the successor publishes the new identity in the README
  and in an advisory. Receipt signing keys belong to each user's machine, not to the project.
- **Stopping deliberately.** If a kill criterion ends image maintenance, the plan is to keep the
  RPM, the benchmark and the specifications available, mark the image as unmaintained, and record
  the decision in an ADR.
- **Forks.** Anyone may fork under Apache-2.0. Forks must use a different name (see
  [TRADEMARK.md](TRADEMARK.md)).
- **Longer term.** Phase 7 (conditional, from October 2027) includes proposing the specification to
  a neutral home such as an OpenSSF or Agentic AI Foundation working group, publishing a
  bus-factor plan, and a written decision on the project's long-term form (community-only, a hosted
  attestation service, or a foundation).

## Transparency

- The roadmap, acceptance criteria and kill criteria are public.
- Progress is recorded in devlogs under `docs/devlog/`.
- Interview findings are published only as an anonymised synthesis under `docs/interviews/`.
- The project collects no telemetry.

## Changing this document

Changes to this document follow the ADR process above, including the pre-publication exception in
item 5.
