# Receipt predicate (v0)

**Status: Phase 0 draft, revised 2026-10-16 after review. Not implemented and
not tested.** This document specifies the format of a Tidemark session receipt
and the offline verification algorithm. No daemon writes receipts yet and no
verifier reads them yet.

**What a receipt proves, and what it does not.** A receipt is **tamper-evident
against the agent only**: a session process cannot reach the recording daemon,
so it cannot add, remove, reorder or alter entries without detection by the
offline verifier, provided the signing key has not been extracted. A receipt is
**not** proof against root or the machine's owner, who can sign anything, until
TPM quotes arrive (Phase 7). The authoritative public check for the v0.1 preview
is the host-side canary grep in the Bring Your Own Canary demo. The receipt is a
secondary, offline, after-the-fact check.

## 1. Envelope and statement

Each session produces one **ledger segment file** and, at its end, one signed
statement:

- an **in-toto Statement v1** as the payload,
- wrapped in a **DSSE envelope** (secure-systems-lab/dsse v1.0.2),
- signed with **pure Ed25519** (RFC 8032; not Ed25519ph).

In v0.1 the signing key is a systemd credential readable only by the ledger
daemon's UID. TPM2 sealing of that key arrives in v0.2 and protects it at rest
only.

### 1.1 DSSE envelope

```json
{
  "payload": "<base64(statement bytes)>",
  "payloadType": "application/vnd.in-toto+json",
  "signatures": [ { "keyid": "<hint only>", "sig": "<base64(Ed25519 signature)>" } ]
}
```

- `payloadType` is fixed. A verifier rejects any other value **before** any
  other processing.
- Writers emit standard, padded base64 (RFC 4648 section 4). Verifiers accept
  standard or URL-safe base64, as DSSE allows.
- The signature covers the Pre-Authentication Encoding of the payload bytes
  exactly as decoded, never a re-serialisation:
  `PAE(type, body) = "DSSEv1" SP LEN(type) SP type SP LEN(body) SP body`.
- `keyid` is an unauthenticated hint and is never used for a security decision.

### 1.2 in-toto Statement

```json
{
  "_type": "https://in-toto.io/Statement/v1",
  "subject": [
    {
      "name": "ledger/<session-id>.jsonl",
      "digest": { "sha256": "<SHA-256 of the segment file's bytes>" }
    }
  ],
  "predicateType": "https://agentos-tidemark.example/attestations/session/v0",
  "predicate": { "...": "section 2" }
}
```

The subject is the session's ledger segment, an immutable artifact matched by
digest, as in-toto requires. There is exactly one segment file per session.
`predicateType` is a **placeholder URI** that changes when the real
specification home is created.

### 1.3 Checkpoints

A crashed or unended session yields no final statement. To make truncation
visible, ledgerd also signs a **checkpoint** statement every 256 lines or 60
seconds, whichever comes first, with predicate type
`https://agentos-tidemark.example/attestations/checkpoint/v0` (placeholder), the
subject name `ledger/<session-id>.jsonl` with the digest of the segment's first
`line_count` lines, and the predicate
`{ "session_id", "line_count", "chain_head" }`. The final statement supersedes
all checkpoints.

## 2. Predicate fields

Hashes are lower-case hex SHA-256 unless named `*_hmac`. Fields named `*_hmac`
are HMAC-SHA-256 under a per-installation key that tidemarkd holds outside
sessions as a systemd credential; the key never appears in a receipt, so a
receipt holder cannot confirm guessed paths or prompts. Timestamps `at` are RFC
3339 UTC (informational). `mono_ns` is `CLOCK_MONOTONIC` in nanoseconds, as a
decimal **string**, because it can exceed 2^53.

```json
{
  "session_id": "<opaque id>",
  "profile": "<profile name>",
  "agent": { "identity": "<name and version>", "version_hash": "<sha256 of the agent artifact>" },
  "image_digest": "<per-architecture image manifest digest, from bootc status>",
  "policy_hash": "<sha256 of the active policy file bytes>",
  "enforcement": "audit | enforce",
  "start_label": ["..."],
  "doctor_result": {
    "passed": true,
    "guarantee_level": "image | rpm-on-stock | degraded",
    "checks": { "landlock_abi": 10, "bpf_lsm": true, "kfuncs": true, "overlay_xattr": true }
  },
  "key_kind": "systemd-creds | tpm2-sealed | tpm2-quote",
  "started": "<RFC3339>",
  "ended": "<RFC3339>",
  "label_raises": [
    {
      "seq": 3, "at": "<RFC3339>", "mono_ns": "<decimal>",
      "cause": "file | source | resync",
      "path_hmac": "<present when cause=file>",
      "source": "<destination name, present when cause=source>",
      "from_label": ["..."], "to_label": ["..."]
    }
  ],
  "exit_decisions": [
    {
      "seq": 4, "at": "<RFC3339>", "mono_ns": "<decimal, read before the label load>",
      "exit": "gateway | broker",
      "point": "open | deny | close | request",
      "conn_id": "<gateway connection id or broker request id>",
      "destination": "<destination or backend name, or none>",
      "label": ["..."], "clearance": ["..."],
      "decision": "allow | deny | audit_allow",
      "reason": "<code, present when decision is deny or audit_allow>",
      "content_hmac": "<broker only, required: the exact bytes sent upstream after allow-list filtering>",
      "release_id": "<present when the allow relies on a release>",
      "digest_pinned": true,
      "bytes_up": 0
    }
  ],
  "tunnel_resets": [
    { "seq": 5, "at": "<RFC3339>", "mono_ns": "<decimal>", "conn_id": "<id>", "destination": "<name>", "new_label": ["..."] }
  ],
  "releases": [
    {
      "seq": 6, "at": "<RFC3339>", "mono_ns": "<decimal, when minted>",
      "release_id": "<opaque id>", "session_id": "<id>",
      "content_hmac": "<the exact released bytes>", "destination": "<backend name>",
      "approver_uid": 1000
    }
  ],
  "commit": {
    "performed": true, "at": "<RFC3339>",
    "final_label": ["..."], "state_delta_label": ["..."],
    "files_stamped": 0, "files_held": 0,
    "tree_hash_before": "<sha256>", "tree_hash_after": "<sha256>"
  },
  "abort": { "performed": false, "at": null, "tree_hash_restored": null },
  "line_count": 0,
  "chain_head": "<SHA-256 of the last ledger line's bytes>"
}
```

Field notes:

- **Recording points.** The gateway records an exit decision when a connection
  or plain-HTTP request is first allowed (`open`), on every denial (`deny`),
  and when a connection closes (`close`), carrying the label from the load for
  the last forwarded chunk and that load's `mono_ns`. Because labels only rise,
  the close record covers every earlier chunk. The broker records one
  `request` decision per request; it doubles as the placement record
  (`destination` is the chosen backend, or `none` when it fails closed).
- **`mono_ns` placement.** Exits read the clock immediately **before** loading
  the label word. The BPF program and tidemarkd read it immediately **after**
  the atomic OR. This is what makes the staleness check in section 4 sound.
- **`resync`.** If the kernel ring buffer drops a label event, tidemarkd reads
  the word directly and records a `resync` raise. The verifier accepts it as a
  chain step but reports the receipt as degraded (section 4, step 8.6).
- **Informational fields.** `key_kind` and `doctor_result.guarantee_level` are
  asserted by the writer inside the signed payload. They describe the writer's
  view; the strength of the evidence is the verifier's trust in the key it was
  given.
- `image_digest`: `bootc status` is expected to report the per-architecture
  manifest digest, not the multi-architecture index digest. The spike confirms
  this, and the image pipeline signs every child image (`cosign sign
  --recursive`) so that the per-architecture digest carries a signature.
- `approver_uid` is the numeric UID of the human who released, never a user name.

## 3. Ledger segment format

- One segment per session at `ledger/<session-id>.jsonl`, carrying the kernel
  append-only attribute (`FS_APPEND_FL`) set by tidemarkd. ledgerd can append
  but not truncate or rewrite it. Root still can (a stated non-guarantee).
- A **line** is one UTF-8 JSON object without a byte-order mark, followed by a
  single LF. "The line's bytes" means the bytes before that LF. A final line
  without its LF is a torn write and makes the segment incomplete.
- Writers emit compact JSON with keys in the order shown, no floating-point
  numbers anywhere, and integers between 0 and 2^53 - 1. Verifiers never depend
  on key order or byte layout, because they never re-serialise.

Each line:

```json
{"seq":0,"session_id":"<id>","ts":"<RFC3339>","mono_ns":"<decimal>","kind":"session_start","body":{},"prev":"<64 hex>"}
```

`kind` is one of `session_start`, `label_raise`, `exit_decision`,
`tunnel_reset`, `release`, `commit`, `abort`, `session_end`. `body` has the
shape of the matching predicate entry; `session_start.body` carries `profile`,
`agent`, `image_digest`, `policy_hash`, `enforcement`, `start_label`,
`doctor_result` and `key_kind`.

The chain rule has no self-hash field:

```
line[0].prev = "0" * 64
line[i].prev = hex(SHA-256(bytes of line[i-1]))      for i >= 1
chain_head   = hex(SHA-256(bytes of line[n-1]))
```

## 4. Verification algorithm

`tidemark-verify` runs offline on macOS, Linux and Windows. Inputs: a DSSE
envelope (or, for an unended session, the latest checkpoint), the segment file,
a trusted public key, and optionally the policy file and a cosign-verified image
digest. Lines are decoded with a number-preserving decoder (in Go,
`json.Decoder.UseNumber`), and a line with a duplicate object key is rejected.

1. **payloadType.** Reject unless it equals `application/vnd.in-toto+json`.
2. **Signature.** Decode the payload; verify the Ed25519 signature over
   `PAE(payloadType, payload bytes)` with the trusted key. Reject on failure.
3. **Statement shape.** Check `_type`, the expected `predicateType`, and exactly
   one subject named `ledger/<session_id>.jsonl`.
4. **Subject digest.** SHA-256 of the supplied segment bytes (for a checkpoint,
   of its first `line_count` lines) must equal the subject digest.
5. **Framing.** Every line parses; `seq` equals the line index; every line's
   `session_id` equals the predicate's; line 0 has kind `session_start`.
6. **Chain.** Check the chain rule for every line and report the **exact line
   number** of the first mismatch. Check `chain_head` and `line_count`. A
   complete segment ends with `session_end`.
7. **Predicate equals ledger.** Each predicate array equals, in order, the
   bodies of the ledger lines of the matching kind, compared as decoded values;
   the top-level session fields equal the `session_start` body.
8. **Replay of the clearance invariants** (`spec/clearance-model.md`):
   1. **Raise chain.** Ledger order is arrival order, not causal order, so raises
      are ordered by `to_label` under subset inclusion. They must form one
      chain: each `to_label` strictly contains its `from_label`, no two raises
      share a `to_label`, the first `from_label` equals `start_label`, and each
      later `from_label` equals the previous `to_label`. The final value
      `L_final` is the last `to_label`, or `start_label` if there are none.
   2. **Exit labels are chain values.** Every exit decision's `label` equals
      `start_label` or some `to_label`.
   3. **No stale label.** For every exit decision, every raise whose `mono_ns`
      is smaller than the decision's `mono_ns` has a `to_label` contained in the
      decision's `label`.
   4. **Decisions.** `allow` requires `label` to be a subset of `clearance`, or a
      `release_id` naming a release with equal `content_hmac`, `destination` and
      `session_id`, minted earlier (smaller `mono_ns`), and named by no other
      decision. `audit_allow` requires `enforcement` to be `audit` and `label`
      not to be a subset of `clearance`. `deny` requires `label` not to be a
      subset of `clearance`, or a non-label reason code.
   5. **Commit.** `commit.final_label` and `commit.state_delta_label` equal
      `L_final`.
   6. **Gaps.** If any raise has cause `resync`, steps 1 and 2 still apply, but
      step 3 is unverifiable around it and the result is reported as degraded.
   7. **Policy (when supplied).** If SHA-256 of the supplied policy file equals
      `policy_hash`, each decision's `clearance` must equal that policy's
      clearance for the named destination or backend, and each `source` raise
      must name a destination with a `source_label` and have a `to_label` equal
      to its `from_label` united with that `source_label`.
9. **Image digest (optional).** If given a cosign-verified per-architecture
   digest, check it equals `image_digest`.

**Results.** `PASS` when steps 1 to 8 pass; `PASS (degraded)` when step 8.6
applies; `FAIL` with the step and line; `INCOMPLETE` when only a checkpoint
exists, meaning lines up to its `line_count` verified and any later lines did
not. The verifier prints `key_kind` and `guarantee_level` as writer-asserted
information and a reminder that the receipt is tamper-evident against the agent
only, and that the authoritative external check is the host-side canary.

## 5. Placeholders and open items

- The `predicateType` URIs (`https://agentos-tidemark.example/...`) are
  placeholders.
- The Go module path is `github.com/BitecodesHub/tidemark` (resolved 2026-10-23,
  decisions ADR PD-10).
- The writer encoding rules in section 3 MUST be frozen before the first
  receipt is written. Verification does not depend on them, but independent
  writers must agree.
- The ring-buffer drop handling (`resync`) and the clock reads around the
  atomic OR are design intent; Phase 3 tests them.
- The v1 schema freeze and TPM-quote binding are Phase 7 work.
