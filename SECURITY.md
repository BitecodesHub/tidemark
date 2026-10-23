# Security Policy

> **Status: Phase 0 draft, 16 October 2026.** No Tidemark code, image or release exists yet. Until
> the first preview release, reports about the design (the threat model, the specifications, the
> benchmark methodology and the build pipeline) are welcome and handled under this policy.

AgentOS Tidemark (working name) intends to make a public security claim. A claim like that is only
worth something if people can report failures safely and see them handled openly. This policy
explains how to report a vulnerability, what counts as one, and what you can expect in return.

## Reporting a vulnerability

**Do not open a public issue, pull request or discussion for a suspected vulnerability.**

Report it privately through GitHub private vulnerability reporting:

1. Open the repository's **Security** tab.
2. Select **Report a vulnerability**.
3. Complete the advisory form.

Use the Security tab of the repository that hosts this file. A direct link will be added once the
project's GitHub organisation is registered.

The maintainer enables private vulnerability reporting in the repository settings before the
repository is made public.

GitHub documents the process in
[Privately reporting a security vulnerability](https://docs.github.com/en/code-security/how-tos/report-and-fix-vulnerabilities/report-privately).

If you cannot use GitHub, contact: `[SECURITY CONTACT: to be set by the founder before the
repository is made public]`.

Please include:

- the guarantee (G1 to G9 below) you believe is broken, or the component affected;
- the image digest or commit, the `tidemark doctor` output and the policy hash, once these exist;
- exact reproduction steps, ideally as a script or a benchmark scenario;
- the observed result and the expected result;
- whether AI tools were used, as required by the AI-generated reports policy below;
- how you would like to be credited, or that you prefer not to be credited.

Never include real secrets or real personal data. Use a canary you generated yourself.

## What to expect

The project has one part-time maintainer (see [GOVERNANCE.md](GOVERNANCE.md)). The timelines below
are honest targets for that capacity, not contractual commitments.

| Step | Target |
|---|---|
| Acknowledgement of the report | within 7 days |
| Initial assessment (in scope or not, severity) | within 14 days |
| Status updates to the reporter | at least every 14 days while the report is open |
| Coordinated disclosure | **90 days** after the report is received |

A core guarantee is any of G1 to G9 (listed under Scope below).

After the v0.1 preview launches, triage is capped at about 2 hours a week. Reports that affect a
core guarantee take priority over all other project work.

### Coordinated disclosure

- The disclosure window is **90 days** from the date the report is received.
- Disclosure can happen earlier when a fix is released and the reporter agrees.
- The window can be extended once, by agreement with the reporter, when a fix is in progress.
- If a vulnerability is being exploited in the wild, the project may publish an advisory and
  mitigations before a fix exists.
- **If a confirmed bypass of a core guarantee is not fixed within 90 days of receipt, the project
  withdraws that guarantee publicly and documents the bypass. An agreed disclosure extension does
  not move this deadline.** If an extension has been agreed, the guarantee is still withdrawn at
  day 90, and the technical details of the bypass are published when the extension ends. This is
  one of the project's published kill criteria.
- Advisories are published as GitHub security advisories. A CVE is requested when the issue
  affects a released artifact.

## Scope

### In scope: the guarantees

Tidemark's preview guarantees are defined in [docs/THREAT_MODEL.md](docs/THREAT_MODEL.md) and in the
[concept document](docs/research/2026-10-16_Tidemark-Concept_v1.md). They are claimed only on the
signed image, when `tidemark doctor` passes and the kernel is not compromised. A report that shows
any of the following, under those conditions, is in scope:

| Guarantee | A vulnerability is |
|---|---|
| G1 Containment | a session descendant escaping the session's cgroup subtree, network namespace or mount view; surviving `stop`; creating a namespace (including a user namespace) or an io_uring instance |
| G2 Single exit | any network path from a session other than `tidemark-gw` and `tidemark-broker`, including DNS, UDP, QUIC, IP literals or host Unix sockets |
| G3 Label before read | a labelled open, or a received descriptor to a labelled file, returning before the session's label rises; a private-source response byte reaching a session before the raise; a session lowering any label |
| G4 Clearance on every chunk | any byte forwarded while the session's current label is not a subset of the destination's clearance, including through a connection opened before the read |
| G5 Model channel | a labelled request placed on a backend outside its clearance; a silent fallback to an uncleared backend; provider-side tools or non-allowlisted fields passed for a labelled session; a model API key reaching a session |
| G6 Staged writes carry labels | staged writes visible to other sessions before commit; a committed file missing the session's final label; `abort` failing to restore a byte-identical workspace; auto-exec, skill or hook files activating without acceptance; labelled memory visible to a lower session |
| G7 Evidence against the agent | a session reaching the ledger; an edit, deletion, insertion or reordering of a receipt that `tidemark-verify` accepts, provided the signing key was not extracted |
| G8 Stop | a session that keeps running, or an exit that stays open, after `tidemark stop` |
| G9 Declassification | data released without a human acting from outside the agents slice; a release token reused, or applied to different content or a different destination |

Also in scope:

- the build and release pipeline: workflow injection, unpinned actions or images, signing identity
  confusion, or anything that lets a third party produce an artifact that passes `cosign verify`
  with the certificate identity and OIDC issuer the project documents for its release workflow;
- `tidemark-verify` accepting a malformed or forged receipt;
- benchmark integrity: a scenario that reports a pass while canary bytes reach a host sink;
- privilege escalation from a session, or from an unprivileged Tidemark daemon, to root;
- receipts exposing paths or bodies that the documentation says are hashed.

### Out of scope: the non-guarantees

These are documented limitations, not vulnerabilities. Reports that demonstrate one of them
convincingly are still welcome as documentation or benchmark contributions.

- Prompt injection itself, and the correctness of agent output.
- Attacks that require kernel compromise or root access.
- Covert and timing channels toward destinations the user has cleared (request counts, timing,
  choice of URL, CDN domain fronting through a cleared domain).
- Misuse of data by a destination the user has cleared, including a cleared cloud model.
- Unlabelled data, and data from APIs not tagged as private sources.
- Labels lost when files are copied off the machine or rewritten by tools outside sessions.
- External side effects at cleared destinations in v0.1.
- Data sent through any URL, method or body on a cleared domain (v0.1 clears per domain and does
  not inspect content).
- Loss of public egress for the rest of a session after it reads labelled data (session
  granularity is coarse by design).
- Always-on messaging-gateway modes, VPS hosting and OpenClaw.
- Agents launched outside `tidemark run`, which are not sessions.
- Covert and timing channels between concurrent sessions on the shared kernel and services.
- Encoded or disguised content in a release that a human approved.
- Receipt content that is false because an event source (`tidemarkd`, `tidemark-gw`,
  `tidemark-broker`) was already compromised. The compromise itself is in scope under the
  guarantee it breaks.
- Host model endpoints recorded as unattested, such as macOS Metal or Ollama.
- Receipts as proof against root or the machine's owner.
- Hosts running only the RPM, beyond the guarantee level `tidemark doctor` reports for them.
- Resource exhaustion within the configured quotas.
- Physical attacks, social engineering of the maintainer, and attacks on GitHub, Sigstore, COPR or
  Fedora infrastructure.
- Vulnerabilities in upstream projects (the Linux kernel, Fedora packages, bootc, llama.cpp,
  Hermes Agent, OpenCode and others). Please report those upstream; tell us as well if Tidemark's
  configuration makes them worse.

## AI-generated reports

AI tools can help find real bugs. They also produce plausible reports about bugs that do not exist,
and each one costs a part-time maintainer hours. The following rules apply:

1. **Disclose AI use.** State whether AI tools helped you find the issue, write the report or write
   the proof of concept, and name the tools.
2. **Reproduce it yourself.** Where a build or kernel exists, run the reproduction on a real build
   or a real kernel before reporting. Say what you ran and what you saw.
3. **Stay accountable.** You must be able to answer follow-up questions about the report without
   relaying them to an AI tool.
4. Reports without a working reproduction, or that cite code, functions or behaviour that do not
   exist, may be closed without further response.
5. Undisclosed AI-generated reports, once identified, may be closed and the reporter may be
   excluded from the hall of fame and from future reports.

Before the first preview release there is no build. For design reports, a reproduction means the
exact document section, the attack steps and the reasoning that shows the guarantee fails. Rules 1
and 3 still apply.

Disclosed, verified, AI-assisted reports are treated exactly like any other report.

## Hall of fame

With their consent, reporters of valid in-scope issues are credited here and in the advisory. Each
confirmed bypass also becomes a regression scenario in `tidemark-bench`, credited to the reporter.

| Date | Reporter | Guarantee or component | Advisory |
|---|---|---|---|
| (none yet) | | | |

There is no bug bounty. The project has no budget for one.

## Safe harbour

The project will not pursue or support legal action against good-faith research that follows this
policy. Good faith means that you:

- test only on systems you own or are authorised to test, such as your own VM;
- do not access, modify or keep other people's data;
- do not degrade services for others;
- report promptly and keep the details private until disclosure.

This statement covers the project's own position only. It cannot authorise testing of systems
operated by others, including cleared destinations or upstream services.

## Supported versions

| Version | Supported |
|---|---|
| Phase 0 drafts (this repository today) | Design reports only |
| v0.1 preview (planned for April 2027) | The latest preview release and the `main` branch |

The v0.1 release is a preview. Its guarantees are preview guarantees, as stated in the threat model.
