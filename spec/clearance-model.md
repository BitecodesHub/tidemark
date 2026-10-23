# Clearance model (v0)

**Status: Phase 0 draft, revised 2026-10-16 after review. Not implemented and
not tested.** This document states the rule that the kernel label program, the
gateway and the model broker are intended to enforce. No component enforces it
yet. It describes design intent for the Phase 0 kill gate, and it stays within
the "not claimed" list of the concept document
(`docs/research/2026-10-16_Tidemark-Concept_v1.md`). In particular, this model
does not prevent prompt injection, does not resist kernel or root compromise,
does not close covert or timing channels to cleared destinations, and does not
inspect content. Every guarantee below holds only while the preconditions in
section 4 hold.

Normative words (MUST, MUST NOT, SHOULD) apply to the future implementation.

## 1. Labels and the lattice

A **label value** is a finite set of **label names** declared in the policy
lattice (`policy/default.yaml`, validated by `spec/policy-v0.schema.json`). The
default policy declares `private`, `secret` and the compartment `mail`. Each
declared name owns one stable bit of a 64-bit kernel **label word**, so a policy
declares at most 64 names, and file attributes store bits, not names.

Order between two label values is **subset inclusion**:

```
A <= B   if and only if   A is a subset of B
```

- **Bottom** is the empty set `{}`: *unlabelled*. **"public" is the reserved
  display name for `{}`**; it is not a label name and cannot be declared.
- The **join** of two values is their union; the **meet** is their intersection.
- Names are independent elements, not points on a line. `{private}` and
  `{secret}` are incomparable; a session that reads both reaches
  `{private, secret}`, which dominates each.

**A clearance is an upper bound.** It lists the labels a destination may
receive. It can never *require* a label: `{}` is a subset of every clearance, so
an unlabelled session satisfies every exit's clearance.

## 2. The session label (high-water mark)

Each session has one current label value, `L_session`, held in the session's
slot of the kernel label word. It is **raise-only** for the life of the session.
Three events raise it.

1. **Session start with `--label`.** tidemarkd sets the slot to the start label
   before the shim execs the agent. A session started without `--label` starts
   at `{}`.
2. **Open or receipt of a labelled file.** When a session process opens a file
   for reading, or receives a readable descriptor, the kernel label program ORs
   the file's label into the slot **before the call returns**. A file's label is
   the union of its `security.bpf.tidemark` attribute and every directory rule
   that matches its session-visible path (see the schema's `directory_rules`).
3. **Receipt from a private source.** Before the gateway relays the first
   response byte from a destination that declares a `source_label`, it asks
   tidemarkd to raise `L_session` by that set, and waits for completion.

**Raise-only, normatively.**

- Every raise MUST be an atomic fetch-OR on the 64-bit word: by the BPF program,
  and by tidemarkd through a writable mapping that only tidemarkd holds.
- No component MAY store a whole value into a live slot. The only whole-value
  stores are slot initialisation before the session's first process starts, and
  zeroing after every session process has exited and the gateway and broker
  have dropped the session.
- The gateway and broker hold read-only mappings and only load the word.

**Capture at open time.** A file's label is captured when it is opened or its
descriptor is received. A later change to the file's attribute, to a directory
rule or to the policy does not reach descriptors that are already open (see
precondition P4).

The agent cannot lower `L_session`: it cannot set `security.bpf.*` attributes
(that requires `CAP_SYS_ADMIN`, which sessions lack), cannot reach tidemarkd's
control socket, and holds no mapping of the label word.

## 3. The clearance rule

Every exit carries a **clearance**. The one rule, applied identically at the
network gateway and the model broker:

```
Bytes may leave the session to a destination D, at a given moment,
only while  L_session  is a subset of  clearance(D).
```

- **Per-chunk re-check.** The gateway MUST load the label word **after** the
  `read()` that returned a chunk from the session socket and **before** the
  upstream `write()` of that chunk. This ordering is what makes the check sound:
  the raise completes before `open()` returns, the session's write to the socket
  follows that return, and the gateway's read follows the session's write, so a
  load placed after the read observes the raise. A load placed before the read
  could miss it.
- **Broker.** The broker MUST load the word after it has read the full request
  body and before it writes the first byte upstream.
- **Fail-closed and content-blind.** When the subset test fails, the exit
  stops. The decision never depends on the bytes, so encoding, compression,
  splitting or paraphrase change nothing.

### 3.1 Why a pre-opened tunnel does not help the attacker

The classic bypass opens a tunnel to a cleared public host while unlabelled,
reads a secret, then streams it through the open tunnel. The model closes it in
two ways:

1. **Per-chunk subset test** (section 3), which is where correctness rests.
2. **Tunnel reset on label change.** When `L_session` rises, the gateway resets
   every live connection whose clearance no longer covers it. Change events can
   race, so this is defence in depth only.

The Phase 2 bypass matrix must demonstrate this. It is not demonstrated yet.

### 3.2 Enforcement modes

- **enforce:** the rule above. Denied chunks are not forwarded, uncovered
  tunnels are reset, and the broker fails closed when no cleared backend fits.
- **audit:** labels are raised and recorded exactly as in enforce mode, but the
  exits forward bytes the rule would deny and record the decision
  `audit_allow` in the receipt; they never reset tunnels. The broker still
  prefers a cleared backend; when none fits, it places the request on the first
  backend in policy order and records `audit_allow`. **Audit mode provides no
  confidentiality.** It exists for measurement and as benchmark configuration
  (c).

## 4. Preconditions and write rule

The rule in section 3 governs only the mediated exits. It is sound only if no
other output path exists and if written data carries labels forward. Breaking
any precondition below voids the guarantee; `tidemark doctor` checks those it
can, and the receipt records the doctor result.

- **P1. Only two unmediated outputs.** The only paths from a session to anything
  outside it are its two per-session sockets (gateway and broker). Mechanisms,
  per the concept's session principal: `PrivateNetwork` (loopback only, in a
  per-session network namespace), `PrivateIPC`, `PrivateTmp`, a mount view with
  `ProtectSystem=strict` and `ProtectHome=tmpfs`, Landlock at ABI 9 or later
  with `LANDLOCK_ACCESS_FS_RESOLVE_UNIX` (new in ABI 9, Linux 7.1) limiting
  pathname Unix sockets to the two per-session sockets, Landlock scoping of
  abstract Unix sockets and signals, and seccomp denial of `ptrace`,
  `process_vm_readv/writev`, `pidfd_getfd`, `io_uring_setup`, `bpf`, `mount`
  and `keyctl`. **Open item:** the concept's unit properties do not name a
  private `/dev/shm`; Phase 1 must show that `/dev/shm`, FIFOs and other
  writable shared locations are private to the session or unwritable.
- **P2. Every write is staged.** Every path the session can write is a
  per-session overlay upper or a per-session private tmpfs, invisible to other
  sessions until commit (G6). Landlock write, `make_fifo` and `make_sock`
  rights are granted only beneath those paths. Doctor refuses a profile that
  declares a writable path without an overlay.
- **P3. Commit stamps with the final label and never lowers one.** Commit sets
  each committed file's label to the union of its existing label (the attribute
  that copy-up carried into the upper) and the session's **final**
  `L_session`. The agent-state delta merges into the delta for the final
  `L_session`. Neither the start label nor any intermediate value is used.
- **P4. Label and policy changes apply forward.** A policy change applies only
  to sessions started after it; each session records the policy hash it started
  with. Setting or adding a file label while a session holds that file open
  does not affect reads through descriptors opened earlier; label files before
  sessions can open them. Raising every live session conservatively when labels
  are added is a Phase 1 option.
- **P5. Rule-labelled files cannot be moved out of their rule.** A directory
  rule labels a path, not an inode. A session that could rename or hard-link a
  rule-labelled file to an unruled path inside its own writable overlay would
  read it unlabelled. Therefore either (a) every existing rule path is visible
  to sessions only through read-only binds, so `rename` and `link` fail with
  `EROFS` or `EXDEV`, and doctor refuses a profile whose writable paths contain
  a rule path; or (b) the labeller denies `rename` and `link` whose source
  matches a rule and whose destination does not (the `path_rename` and
  `path_link` BPF LSM hooks exist at Linux 7.2). The choice is a Phase 1
  decision; the Labeller ADR lists the test.
- **P6. Session-visible paths equal host paths.** Read paths and writable
  overlays are bound at their own absolute paths, and sessions do not use
  `RootDirectory=`, so a directory rule written for the host path matches the
  path a session sees.

## 5. The model call is an exit

For a request from a session at `L_session`, the broker:

1. filters the declared backends to those whose clearance covers `L_session`;
2. picks the first survivor in policy order (static preference);
3. if none survives, **fails closed** with a structured reason in enforce mode.
   Nothing falls back to a cloud backend silently.

With the default policy, `local-llama` is first and cleared for every declared
label, so every request lands on the local model while it is available. For
labelled sessions the broker also drops request fields outside an allow-list
and any tool that is not a plain function, which removes provider-side web
fetch, web search and remote connectors. The broker holds the only copy of
backend API keys; sessions see placeholders, which carry no label.

## 6. Declassification (human release only)

The session can never move data to a destination its label does not cover. The
only way is a **human release**, performed from outside the agents slice. In
v0.1 a release applies **only to broker requests**, because only there are the
exact outgoing bytes visible in plaintext on the local link (G9). A release:

- is bound to the **content hash** of the exact bytes the broker will send
  upstream after allow-list filtering, to the exact backend, and to the session;
- is a **one-time token** consumed by the agent's own retry; Tidemark never
  replays anything itself;
- is recorded in the receipt, with its consumption.

Releases for opaque HTTPS `CONNECT` tunnels and for plain HTTP wait for v0.2
(see `docs/adr/2026-10-16_No-TLS-Interception-v0.1_v1.md`).

## 7. Worked examples

Notation: `L` is the session label after each step; `clr(D)` is a clearance.

### 7.1 Baseline public task

1. Session starts without `--label`. `L = {}`.
2. Reads an unlabelled repository file. `L = {}`.
3. Sends to `public-paste`, `clr = {}`. `{}` is a subset of `{}`. **Allowed.**

### 7.2 Read, then send, and the mail compartment

1. `L = {}`.
2. Reads `~/Private/notes.txt` (rule `{private}`). `L = {private}` before
   `open()` returns.
3. Sends to `public-paste`, `clr = {}`. **Denied.**
4. Sends to `mail-api`, `clr = {mail}`. `{private}` is not a subset of
   `{mail}`. **Denied.** Without the compartment, private documents could be
   mailed to any recipient.
5. A different session reads only the inbox through `mail-api`: `L = {mail}`.
   Replying through `mail-api` is **allowed**; v0.1 cannot see recipients, so
   that session could mail inbox content to anyone (a stated non-guarantee).

### 7.3 Tunnel opened before the read (the key case)

1. `L = {}`. Opens a `CONNECT` tunnel to `public-paste`; `{}` passes.
2. Reads `~/.ssh/id_ed25519` (rule `{secret}`). `L = {secret}` before the
   `open()` returns.
3. Writes the key bytes into the open tunnel. After reading that chunk and
   before forwarding it, the gateway loads `L = {secret}`; the test against `{}`
   fails. The chunk is **not forwarded** and the tunnel is reset.

### 7.4 Model placement under a private label

1. `L = {}`. Reads `~/Private/notes.txt`. `L = {private}`.
2. Calls `tidemark/auto`. Only `local-llama` covers `{private}`. The request
   lands locally; if it cannot, the broker fails closed.

### 7.5 Memory laundering across sessions

1. Session A starts with `--label private` (`L_A = {private}`), reads
   `~/.ssh/config` (`L_A = {private, secret}`) and writes its contents into
   agent memory. At commit, the delta is stamped with the **final** label
   `{private, secret}` (P3), not with the start label.
2. Session B starts without `--label` (`L_B = {}`). Its state view stacks only
   deltas whose label is a subset of `{}`, so the delta is invisible.
3. Session C starts with `--label private`. `{private, secret}` is not a subset
   of `{private}`, so the delta is still invisible. Only a session started at
   `{private, secret}` or above sees it, and that session cannot reach a public
   host.

The Phase 3 laundering test must demonstrate this. It is not demonstrated yet.

### 7.6 Source label from an API

1. `L = {}`. Opens `CONNECT` to `mail-api`; `{}` passes.
2. Before the first response byte, the gateway raises `L` to `{mail}`.
3. A later send to `public-paste` fails. **Denied.**

### 7.7 Moving a labelled file out of its rule

1. `L = {}`. The session runs `mv ~/Private/notes.txt ~/work/notes.txt`.
2. Under P5(a) `~/Private` is a read-only bind and the rename fails with
   `EROFS` or `EXDEV`. Under P5(b) the labeller denies it. Without P5, a later
   `open()` of `~/work/notes.txt` would match no rule and read unlabelled data.

## 8. Scope and limits (restated)

This model governs only labelled data and only mediated exits. It does not label
data the user never marked, does not stop a cleared destination from misusing
what it receives, and does not remove covert or timing channels toward a cleared
destination. Every guarantee here is a Phase 0 design claim subject to the
kernel kill gate.
