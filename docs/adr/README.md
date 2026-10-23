# Architecture decision records

This directory holds the architecture decision records (ADRs) for AgentOS
Tidemark. Each ADR captures one significant decision: its context, the choice
made, the consequences, and the alternatives rejected.

**Status: Phase 0.** All current ADRs are Phase 0 drafts with status
**Proposed**. They record design intent before any product code exists. They do
not claim that any feature is implemented or tested.

## Lifecycle

- **Proposed:** a draft. It may be edited in place, and every edit adds a
  dated "revised" note to the status line.
- **Accepted:** the founder has reviewed the ADR line by line and promoted it.
  An accepted ADR is immutable. A later decision that changes it is a new ADR
  that marks the earlier one superseded.
- **Superseded** or **Deprecated:** kept for the record and never deleted.

**Promotion rule.** The founder promotes an ADR from Proposed to Accepted only
after review (the Phase 0 target is the week-3 review, checklist section C). An
ADR whose decision depends on an experiment stays Proposed until the experiment
passes and its results are written into the ADR. The Labeller ADR therefore
stays Proposed until the Phase 0 kill gate passes and checklist item C2 records
the spike results and the chosen lookup method.

## Naming

New ADRs follow the repository date convention:

```
YYYY-MM-DD_Subject_Version.md
```

For example `2026-10-16_BaseImage_v1.md`. The date is the date the decision was
drafted. The version increments only when the same decision is re-drafted; a
different decision gets a new subject and its own file.

## Index

| ADR | Decision | Status |
|---|---|---|
| [2026-10-16_BaseImage_v1.md](2026-10-16_BaseImage_v1.md) | Base image: `quay.io/fedora/fedora-bootc:45` pinned by digest, `:44` fallback; kernel floor Linux 7.2 with Landlock ABI 9 or later | Proposed (Phase 0) |
| [2026-10-16_Labeller-BPF-LSM_v1.md](2026-10-16_Labeller-BPF-LSM_v1.md) | Read-time labels via sleepable BPF LSM on `file_open` and `file_receive`, rules matched on the session-visible path, no cache by default; fanotify fallback | Proposed (Phase 0), gated by the Phase 0 kill test |
| [2026-10-16_Language-Go_v1.md](2026-10-16_Language-Go_v1.md) | Go for all daemons, the CLI, the verifier and the benchmark (minimum `go 1.25.0`, toolchain `go1.27.2`); BPF in C | Proposed (Phase 0) |
| [2026-10-16_No-TLS-Interception-v0.1_v1.md](2026-10-16_No-TLS-Interception-v0.1_v1.md) | No TLS interception in v0.1; per-domain decisions with a ClientHello SNI check; releases for broker requests only | Proposed (Phase 0) |
| [2026-10-16_Delivery-RPM-Plus-Image_v1.md](2026-10-16_Delivery-RPM-Plus-Image_v1.md) | An RPM built in the same CI run is installed into the signed image by SHA-256; the same source RPM goes to COPR | Proposed (Phase 0) |
| [2026-10-23_Phase0-Founder-Decisions_v1.md](2026-10-23_Phase0-Founder-Decisions_v1.md) | The open Phase 0 decisions, closed: no inode cache (PD-1), 7.2 floor restated as a policy choice (PD-2), BPF sources dual-licensed (PD-3), directory labels follow location with link/rename denied (PD-4), D2/D7 deferred and the concept's wording of both kill criteria governing (PD-6, PD-7), repository created private as `BitecodesHub/tidemark` (PD-10) | Accepted under the founder's 2026-10-23 delegation; no founder line-by-line review yet |

## Template

Copy the block below into a new `YYYY-MM-DD_Subject_Version.md` file.

```markdown
# <Short decision title>

- **Status:** Proposed | Accepted | Superseded by <ADR> | Deprecated
- **Date:** YYYY-MM-DD
- **Phase:** <phase number and name>
- **Deciders:** <names or roles>

## Context

What problem forces a decision now. State the constraints (solo part-time
maintainer, M1 Pro development host, local-only privacy default, and so on) and
cite the evidence from the research documents under `docs/research/` that bears
on the decision. Note anything verified against primary sources, with the source
and the date checked.

## Decision

The decision, stated in the active voice. Be specific: name versions, digests,
APIs and file paths where they matter.

## Consequences

What becomes easier and what becomes harder. Include the follow-on work the
decision creates, the risks it accepts, and the kill or review gate that can
overturn it.

## Alternatives considered

Each rejected option, with the reason it was rejected, grounded in the same
evidence.

## References

Links to the research documents, primary sources, and related ADRs.
```
