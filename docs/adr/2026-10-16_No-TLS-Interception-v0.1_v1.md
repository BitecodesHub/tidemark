# No TLS interception in v0.1

- **Status:** Proposed (Phase 0 draft, revised 2026-10-16 after review and
  2026-10-17 in a cross-file consistency pass; founder review in week 3)
- **Date:** 2026-10-16
- **Phase:** 0 (Foundations and kill gate); revisited in Phase 5 (v0.2)
- **Deciders:** Founder

## Context

The gateway mediates every network exit. A design question is whether, in v0.1,
the gateway should terminate TLS (act as a man in the middle with a per-session
CA) so it can see and decide on the plaintext of HTTPS requests, or whether it
should decide per destination domain and pass HTTPS `CONNECT` tunnels through
opaquely.

Evidence:

- The concept document lists "TLS interception in v0.1" as a non-goal and places
  opt-in TLS termination for declared destinations in Phase 5, using a pinned
  `elazarl/goproxy`. Source: `docs/research/2026-10-16_Tidemark-Concept_v1.md`,
  non-goals and Phase 5.
- The clearance guarantee does **not** need plaintext. The decision is
  content-blind: it compares the session's current label with the destination's
  clearance before every forwarded chunk. Whether the chunk is a byte of
  ciphertext or plaintext does not change the subset test. Source:
  `spec/clearance-model.md`.
- The pre-opened-tunnel bypass (open a `CONNECT` tunnel while unlabelled, read a
  secret, stream it out) is closed by the per-chunk re-check plus tunnel reset on
  label rise, again without reading the tunnel's contents. Source: the red-team
  findings in the concept document and `spec/clearance-model.md` section 3.1.
- The **model** channel does not need interception either, because the broker is
  a local OpenAI-compatible endpoint: it sees prompts in plaintext on the
  loopback link before they leave, so it can strip provider-side tools and place
  by clearance with no certificate work.
- TLS termination carries real cost and risk for a solo maintainer: a per-session
  CA in each session's trust store, a proxy that parses attacker-influenced TLS
  and HTTP, clients that pin certificates and break, and a larger trusted
  computing base. mitmproxy is not packaged: the Fedora package metadata API
  (`mdapi.fedoraproject.org`) returned "not found" for `mitmproxy` on the
  `rawhide`, `f44`, `epel9` and `epel10` branches on 2026-10-16, while returning
  a record for `curl` on `rawhide` and `f44`. (mdapi had no `f45` branch that
  day.)

## Decision

Do **not** intercept TLS in v0.1.

- The gateway implements HTTP `CONNECT` and plain HTTP forwarding. It resolves
  DNS itself, only after the clearance check (G2), and refuses IP literals and
  non-HTTP protocols; the v0 policy schema rejects IP-literal host patterns
  (`spec/policy-v0.schema.json`, `hostPattern`).
- Decisions are made **per destination domain**. v0.1 does not decide by
  method or path. Inside TLS, the HTTP `Host` header and the HTTP/2
  `:authority` are encrypted, so the gateway checks only what it can see:
  - **CONNECT:** the gateway takes the host from the `CONNECT` authority (which
    the session chooses), checks the clearance, resolves the name itself, dials
    the resolved address, and then reads the client's first TLS handshake
    message (a bounded buffer, possibly spanning records) before forwarding
    anything. It refuses the tunnel when that message is not a TLS
    ClientHello, when the SNI is missing, when the SNI differs from the
    authority, or when the ClientHello offers Encrypted Client Hello.
  - **Plain HTTP:** the gateway cross-checks the request target against the
    `Host` header and refuses a mismatch.
  - Domain fronting (an inner `Host` that differs from the SNI) and HTTP/2
    connection coalescing across names that share a certificate remain
    possible inside TLS and stay documented non-guarantees.
- The clearance subset test runs before every forwarded chunk of ciphertext, and
  tunnels are reset when the label rises. Correctness does not depend on reading
  the bytes.
- The model broker sees request bodies in plaintext on the local link, which is
  how provider-side tools are dropped and placement is enforced, with no
  certificate machinery.
- Human release in v0.1 applies **only to broker requests**, as G9 states,
  because only there are the exact outgoing bytes visible in plaintext on the
  local link. Releases for `CONNECT` tunnels and plain HTTP wait for v0.2. (See
  G9 and `spec/clearance-model.md` section 6.)
- Opt-in TLS termination for a small set of declared destinations (git
  smart-HTTP and two declared APIs), each with its own per-session CA, arrives in
  Phase 5 with a pinned `goproxy`.

## Consequences

- The v0.1 trusted computing base stays small: no CA key management per session,
  no TLS termination, and only a bounded ClientHello parser (which must be
  fuzzed) exposed to session-chosen bytes; fewer client compatibility breaks.
- The confidentiality guarantee holds without interception, because labels and
  the per-chunk subset test, not content, decide every exit.
- v0.1 **cannot** apply clearances by method or path, and cannot inject
  credentials for a cleared destination, because it never sees the plaintext of
  an HTTPS request. These are explicit v0.1 limitations and the reason the
  "credential paradox" (a push token that must live in the session) is handled by
  removing the credentialed-push task from the v0.1 benign suite until v0.2.
- Human release in v0.1 is narrower: content-hashed release for broker requests
  only.
- Refusing any ClientHello that offers Encrypted Client Hello also refuses
  clients that send a GREASE ECH extension by default (Chromium-based browsers,
  for example a headless browser an agent drives). Phase 2 measures this. The
  alternative, accepting ECH when the outer SNI equals the authority, leaves the
  inner name unverified and would join domain fronting as a non-guarantee; the
  founder decides in Phase 2.
- Domain-level decisions are coarse. CDN domain fronting through a cleared domain
  remains possible and is documented as a non-guarantee; guidance tells users
  which destinations are safe to clear.
- The Phase 5 revisit is scoped tightly: TLS termination applies only to declared
  destinations, each with its own CA, so the blast radius of the added TCB is
  bounded.

## Alternatives considered

- **Full TLS interception in v0.1 (per-session CA for all HTTPS).** Rejected.
  It is not needed for the confidentiality guarantee, it enlarges the TCB and the
  attack surface, it breaks certificate-pinning clients broadly, and it is too
  much for a solo maintainer to build and secure in the v0.1 budget.
- **No gateway at all, rely on network namespace plus allowlist only.** Rejected.
  The per-chunk label re-check requires a component that reads the kernel label
  synchronously on the data path; a static allowlist cannot revoke a tunnel when
  the label rises mid-connection.
- **Decide by the `CONNECT` authority alone, with no ClientHello SNI check.**
  Rejected: the authority is chosen by the session, and a ClientHello naming a
  different server on a shared address would pass unnoticed. Checking the SNI
  costs one bounded read before forwarding and needs no decryption.

## References

- `docs/research/2026-10-16_Tidemark-Concept_v1.md` (non-goals; Phase 5 TLS
  termination; red-team items on pre-opened tunnels, releasing HTTPS bytes, and
  the credential paradox).
- `docs/research/2026-10-16_OS-Feasibility_v1.md` (egress proxy landscape).
- Fedora package metadata API, `https://mdapi.fedoraproject.org/<branch>/pkg/<name>`,
  queried 2026-10-16.
- `spec/clearance-model.md` (content-blind per-chunk rule; release semantics).
