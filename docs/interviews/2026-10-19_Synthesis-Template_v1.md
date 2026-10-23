# Interview Synthesis Template (v1)

| Field | Value |
|---|---|
| Status | Phase 0 template, revised 2026-10-16 after review. Empty: no interview has taken place. Replace every `[ ]` placeholder when filling it in. |
| Version | v1, prepared 2026-10-16 |
| Owner | Founder. Filled section B records are processed only by tools that run locally on the founder's machine (see "Handling rules"). |
| Published copy | When complete, copy the executive summary, sections A and C to F, and the Tone Check into `docs/interviews/2026-11-06_Interview-Synthesis_v1.md` (proposed name). Filled section B records never enter the repository. Keep this template unchanged. |
| Gate | Kill criterion 2 is evaluated from table C3 at the week-3 gate, Sunday 8 November 2026. |
| Inputs | `2026-10-19_Interview-Guide_v1.md`, `2026-10-19_Interview-Screener_v1.md`, raw notes and section B records held outside the repository |

## Executive summary (fill in last; five bullets at most)
- [Kill criterion 2 result: triggered or not triggered, with the gating counts and the reported readings]
- [Where private data lives: the dominant pattern]
- [Destinations: what participants would and would not send private material to]
- [Hosting: share that could boot an image today]
- [The most important disconfirming finding]

## Handling rules
- **Section B stays outside the repository.** Section B below is a blank form kept here for reference. Filled copies live with the raw notes in encrypted local storage, are never committed, and are deleted with the raw notes by 8 December 2026.
- **AI tools.** Raw notes, recordings, transcripts, contact details and filled section B records are processed only by tools that run locally on the founder's machine (for example, a local speech-to-text model for transcripts made with consent). No cloud AI service sees them. Cloud AI tools may process the published synthesis after the publication checklist (section F) passes.
- **Small cells.** In every published table, merge any channel, agent, host class or incident class with fewer than 3 participants into "Other", or report it as "fewer than 3". The gating counts in table C3 are the only exception: they are published exactly, because the gate depends on them and they describe the whole sample.
- **No per-participant attributes in the published copy.** Kernel versions, hardware models, interview dates and per-person rows are never published.

---

## A. Pre-registered coding rules
Freeze this section before interview P01 and publish the freeze with a timestamp, either as a public commit of this file in the project repository or as the SHA-256 of the frozen file posted in devlog 1. The full text is published with the synthesis. A founder choice that departs from a default in A.4 is valid only if it appears in that published freeze. Any change after the freeze is a deviation (A.6).

**Frozen on:** [date and time, UTC] **SHA-256 of the frozen file, or commit:** [ ] **Published at:** [public commit link, or devlog 1 link] **By:** founder

### A.1 Evidence strength
| Code | Meaning | Counts toward the kill criterion? |
|---|---|---|
| E3 | Observed: shown on screen during the interview (configuration, history, logs, transcript) | Yes |
| E2 | Specific past event: the participant described a particular occurrence | Yes |
| E1 | General statement, opinion or hypothetical | Prong 1: no. Prong 2: only through the scenario probe when no past-behaviour episode applies (A.3), and flagged |

When E2 or E3 evidence contradicts E1 evidence, the E2 or E3 evidence wins.

### A.2 Prong 1: "runs agents over local private files or secrets on a machine or VM they control"
**Private material categories**
| Code | Category |
|---|---|
| M1 | Personal notes or documents |
| M2 | Credentials or secrets: SSH or GPG keys, `.env` files, API tokens, cloud credentials, password-store exports |
| M3 | Local copies of email, messages or calendars |
| M4 | Private source code the participant would object to having published |
| M5 | Private material reached only through an API (for example Gmail, IMAP, Notion, Google Drive) |

**Host classes**
| Code | Host | Counts as "a machine or VM they control" (literal wording)? | Linux-ready? |
|---|---|---|---|
| H1 | Physical machine they own, running Linux | Yes | Yes |
| H2 | VM on hardware they own (for example Proxmox) | Yes | Yes |
| H3 | Cloud VM where they choose the image | Yes | Yes |
| H4 | Container on hardware they own (for example Proxmox LXC) | Yes | No (cannot load its own kernel programs; would need a VM) |
| H5 | Rented VPS without a custom-image option | No | No |
| H6 | Their own macOS or Windows machine, agent running directly on it | Yes | No |
| H7 | Employer-managed machine | No | No |

**Actual use, not capability.** A participant meets a prong-1 definition only through actual use: the agent read the material during a task (E3: shown in history, logs or a transcript; E2: the participant describes a specific task in which it did), or the participant deliberately pointed the agent at it (for example, asked it to summarise their notes or to fix a `.env` file). Configured access alone does not count. Home-directory access, or the presence of `~/.ssh` or `.env` files on a host the agent can read, is P1-capability.

**Definitions**
- **P1-Linux-ready (gating by default; FD-2):** in the past 30 days, with E2 or E3 evidence, an agent the participant runs actually read, or the participant deliberately pointed it at, material in M1, M2 or M3 stored as local files, on a host in H1, H2 or H3.
- **P1-strict (reported; the criterion's literal host scope):** the same condition on a host in H1, H2, H3, H4 or H6.
- **P1-broad (reported only):** P1-strict, or the same condition with M4.
- **P1-capability (reported only):** the agent was configured with read access to M1, M2 or M3 local files on a host in H1 to H4 or H6, but no E2 or E3 evidence shows actual use. Never counted toward the gate.
- **M5-only participants** do not meet any prong-1 definition, because the criterion names local files. They are reported in table C4, because API-sourced private data is in scope for source labels.

**Rationale**
- *Actual use.* Almost any agent with home-directory access on a machine that holds SSH keys has configured access. A capability test would pass almost by construction and would undo the screener's independence from private-data access.
- *Linux-ready hosts by default.* The criterion's wording admits macOS and Windows hosts (H6) and containers (H4), but Tidemark v0.1 protects none of them without a VM. Gating on H1 to H3 measures the users v0.1 can serve; P1-strict is reported beside it.
- *M1 to M3 only.* The default Tidemark policy leaves repository files unlabelled to limit label creep, so private code alone (M4) is not protected by default labels. Counting M4 would overstate the fit.

### A.3 Prong 2: "a majority would abandon a session rather than accept denied public-web access after it reads private data"
**Response codes** (assigned after the interview from guide sections 6.8 and 6.9)
| Code | Meaning | Gating reading (default, FD-3) | Designed-workflow reading (reported) |
|---|---|---|---|
| C | Kept working in the same session without public web | Accept | Accept |
| F | Ended the session and started a fresh one, without the private material, for the web part; includes planned splits | Abandon (the session is abandoned) | Accept (this is the designed workflow) |
| D | Disabled or bypassed the restriction, allowlisted the destination, or switched to a tool without it | Abandon | Abandon |
| X | Gave up the task | Abandon | Abandon |
| U | Unclear or insufficient evidence; or no applicable episode and no private-read session for the probe | Excluded from the denominator; count reported | Same |

**Mapping from past behaviour (guide section 6.8).** Use the first source that applies, in this order: (1) the most recent episode in which network access was blocked mid-task after the session had read private material; (2) a planned split for private material; (3) the most recent episode in which network access was blocked mid-task without a private read; (4) the scenario probe; (5) U. Record which source was used.

| Past behaviour (E2 or E3) | Code |
|---|---|
| Network access blocked mid-task by a tool, firewall or setting; kept working in the same session without the web | C |
| Network access blocked mid-task; restarted in a new session for the web part | F |
| Split one job into separate sessions on purpose, before any block, to keep private material away from web access | F (marked "planned split") |
| Network access blocked mid-task; disabled the block, allowlisted the destination, or switched agent or tool | D |
| Network access blocked mid-task; gave up the task | X |
| Turned permission prompts off, or used auto-approve, for convenience | Not codable for prong 2 (record in C7) |
| Denied a permission prompt that was not about network access | Not codable for prong 2 |

**Scenario probe (guide section 6.9).** It decides the code only when no past-behaviour source applies, and the code is then marked E1. Map the answer with the same codes: carry on without the web (C), start a fresh session for the web part (F), turn the restriction off or switch tools (D), give up the task (X), anything else (U). If the participant has never had a session in which the agent opened their private notes or keys, skip the probe and code U.

**Trigger rule (defaults; FD-3 and FD-4).** Prong 2 triggers if D + X + F is a majority (more than half) under either denominator:
- **Denominator 1:** all completed interviews counted in N (A.5), minus U.
- **Denominator 2:** participants who meet the gating prong-1 definition (P1-Linux-ready by default), minus U.

**Designed-workflow reading (reported only).** Also report D + X over each denominator. This reading treats F as acceptance, because a fresh session for the web part is the workflow Tidemark is designed around. It replaces the gating reading only if the founder chose it in the published freeze (FD-3).

**Recruiting-exposure check (reported only).** A participant is "exposed" if they were recruited from the devlog or project page (screener Q11) or said they had read the project's published material (screener Q12). Report the prong-1 gating count and both prong-2 gating figures again without exposed participants. If any result differs, record both in C3 and discuss the difference in the decision record.

### A.4 Founder decisions required before freezing
**FROZEN 2026-10-23.** Every default below is accepted unchanged (Phase 0
decisions ADR, PD-12); each default is the stricter option, so none of them
makes kill criterion 2 easier to pass. The freeze is published as the SHA-256 of
this file in `docs/devlog/2026-10-16_Devlog_v1.md`, and the freeze date is
recorded as a deviation in A.6 because it fell after Monday 19 October. It still
precedes interview P01, which has not been conducted.

An empty "Founder choice" cell means the default applies. Every default below is the stricter option: none makes kill criterion 2 easier to pass.

| ID | Decision | Default proposed by this template | Founder choice |
|---|---|---|---|
| FD-1 | Private material categories for the prong-1 gate | M1 to M3; M4 is reported as P1-broad | Default accepted 2026-10-23 (PD-12) |
| FD-2 | Host scope for the prong-1 gate | H1 to H3 (P1-Linux-ready). P1-strict, which adds H4 and H6 as the criterion's literal wording allows, is reported | Default accepted 2026-10-23 (PD-12) |
| FD-3 | Is F (a fresh session) acceptance or abandonment? | Gate on the stricter reading: F counts as abandonment, because the criterion says "abandon a session". The designed-workflow reading is reported. The founder may gate on the designed-workflow reading only by recording that choice in the published freeze before P01 | Default accepted 2026-10-23 (PD-12): the stricter reading gates, F counts as abandonment |
| FD-4 | Prong 2 denominator | Trigger under either denominator | Default accepted 2026-10-23 (PD-12) |
| FD-5 | The screener does not select on private-data access (screener, section 4.5) | Confirmed | Confirmed 2026-10-23 (PD-12) |
| FD-6 | Can a triggered criterion be overridden? | No. If kill criterion 2 triggers, the decision is the concept's action: re-scope to a benchmark-first project before Phase 1 (decision-record option 4). Stopping (option 5) remains possible. No option continues Phase 1 as planned | Default accepted 2026-10-23 (PD-12, and PD-7 for criterion 1) |

### A.5 Sample rules
- **Pilot.** P01 is also the pilot. It counts toward N and the kill criterion unless the guide changed materially after P01. In that case P01 is reported separately, excluded from N, and the change is logged in A.6.
- **Referrals** complete the screener like any other respondent, and the same selection rules apply.
- **Disconfirming quota.** Participants who stopped using an agent (screener Q2) count toward N.
- **Reduced sample.** If N is below 8 when recruiting closes on Wednesday 4 November 2026, the founder either (a) evaluates with the threshold unchanged at 3 and marks the result "reduced sample" in C3 and in the decision record, or (b) moves the gate by at most one week, recorded publicly (Phase 0 checklist, section 9). The threshold is never lowered.
- **Withdrawals.** A withdrawal received by Thursday 5 November 2026 removes the participant from N before publication. A later withdrawal removes the participant from later versions only (guide, section 4.4).

### A.6 Deviation log
| Date | Change | Reason | Interviews affected |
|---|---|---|---|
| 2026-10-23 | Coding rules FD-1 to FD-6 frozen with every default accepted (PD-12); the freeze is published as this file's SHA-256 in devlog 1 rather than as a public commit, because the repository is private until launch. | The checklist asked for the freeze by Monday 19 October. No interview has been conducted, so the freeze still precedes P01. | None. |
| [ ] | [ ] | [ ] | [ ] |

---

## B. Per-interview record (blank form; filled copies never enter the repository)
Fill in one copy per participant from the raw notes within 24 hours of the interview, and store it with the raw notes outside the repository. It contains no names, handles, employers, locations finer than a region, hostnames, IP addresses, real paths, secrets or verbatim configuration.

### B.[Pxx]
**Administrative**
| Field | Value |
|---|---|
| Participant ID | P[xx] |
| Interview week | [week 1, 2 or 3] |
| Recruiting channel (screener Q11) | [r/LocalLLaMA, r/selfhosted, Nous Research Discord, Proxmox forum, Fedora Discussion, Universal Blue Discourse, OpenCode community, devlog or project page, referral, other] |
| Read the project's published material before screening (screener Q12) | [yes, no or not sure] |
| Exposed (devlog or project-page channel, or Q12 = yes) | [yes or no] |
| Screener priority | [A, B or C] |
| Disconfirming quota (stopped using an agent, screener Q2) | [yes or no] |
| Duration (minutes) | [ ] |
| Consent to take part confirmed | [yes] |
| Audio recorded | [yes or no]; if yes, notes verified against it on [ ] and deleted on [ ] |
| Pilot interview | [yes or no] |
| Withdrawal received | [no, or the date] |

**Profile**
| Field | Value | Evidence |
|---|---|---|
| Agents used (past 30 days) | [ ] | [E1 to E3] |
| Most-used agent | [ ] | |
| Main use | [coding, research, personal admin, other] | |
| Frequency | [daily, weekly, monthly] | |
| Model mix | [local only, cloud only, both]; which tasks go where: [ ] | |
| Tools, skills or MCP servers (as categories) | [ ] | |
| Messaging-gateway mode in use | [none, or which class] | |
| Approval settings as observed | [prompts on, partial auto-approve, prompts off] | |

**The last surprise** (guide section 6.3; paraphrase, no identifying details)
- What happened: [ ]
- What they did next, and what they changed: [ ]
- Cost (time, money, data, rotated credential): [ ]
- Incident class (for C8): [prompt injection, malicious or compromised tool, deceptive agent behaviour, laundering across sessions or tools, accidental destructive action, unexpected destination, cost overrun, other, none]

**Hosting** (guide section 6.5)
| Field | Value |
|---|---|
| Host class (H1 to H7) | [ ] |
| Proxmox: VM or LXC | [ ] |
| Last rebuild and method | [ ] |
| Has booted a downloaded VM image or used an image-based Linux system | [yes or no; what happened] |
| Mac-to-VM shared folders (virtiofs, 9p, other, none) | [ ] |

**Where private data lives** (guide section 6.6)
| Category | Local file, API, both, or not used | Agent configured with read access? | Actual use in past 30 days (read, or pointed at)? | Evidence |
|---|---|---|---|---|
| M1 Notes or documents | [ ] | [ ] | [ ] | [ ] |
| M2 Credentials or secrets | [ ] | [ ] | [ ] | [ ] |
| M3 Email, messages or calendar (local copies) | [ ] | [ ] | [ ] | [ ] |
| M4 Private source code | [ ] | [ ] | [ ] | [ ] |
| M5 API-only private sources (name the service class) | [ ] | [ ] | [ ] | [ ] |

**Destinations** (guide section 6.7)
| Destination class | Contacted by the agent in the past week? | Private material acceptable there? (yes, no, unsure) | Reason, paraphrased |
|---|---|---|---|
| Local model | [ ] | [ ] | [ ] |
| Cloud model provider | [ ] | [ ] | [ ] |
| Code host | [ ] | [ ] | [ ] |
| Package registry | [ ] | [ ] | [ ] |
| Search engine or web fetch | [ ] | [ ] | [ ] |
| Paste or file-sharing site | [ ] | [ ] | [ ] |
| Chat or messaging service | [ ] | [ ] | [ ] |
| Email or document API | [ ] | [ ] | [ ] |
| A destination they did not know about until later | [ ] | | [ ] |

**Permission prompts and blocked access** (guide sections 6.8 and 6.9)
| Field | Value | Evidence |
|---|---|---|
| Prompts seen and denied in the last session | [ ] | [ ] |
| Turned prompts off or used auto-approve; the reason | [ ] | [ ] |
| Last time network access was blocked mid-task; what they did | [ ] | [ ] |
| Had that session already read private material? | [yes, no or unknown] | [ ] |
| Split work into separate sessions on purpose for private material | [ ] | [ ] |
| Has had a session in which the agent opened private notes or keys (probe precondition) | [yes or no] | [ ] |
| Scenario probe answer (paraphrase), or "skipped" | [ ] | E1 |

**Kill-criterion coding** (apply section A exactly)
| Item | Value |
|---|---|
| P1-Linux-ready (gating by default) | [yes or no]; categories [ ]; host class [ ]; evidence [ ] |
| P1-strict (reported) | [yes or no] |
| P1-broad (reported) | [yes or no] |
| P1-capability only (reported) | [yes or no] |
| Prong 2 code | [C, F, D, X or U]; source [(1) block after a private read, (2) planned split, (3) other block, (4) scenario probe, (5) none]; evidence [ ] |

**Commitment signals**
| Signal | Value |
|---|---|
| Agreed to future contact for testing | [yes or no] |
| Offered a referral | [yes or no] |
| Other concrete offer | [ ] |

**Reaction to the optional project description** (low weight)
- [ ]

**Interviewer notes on bias or problems in this interview**
- [ ]

---

## C. Aggregate tables
Apply the small-cell rule (see "Handling rules") to every table except the gating counts in C3.

### C1. Recruiting funnel and screener base rate
| Channel | Posts | Screener responses | Eligible | Priority A | Priority B | Priority C | Interviewed |
|---|---|---|---|---|---|---|---|
| r/LocalLLaMA | [ ] | [ ] | [ ] | [ ] | [ ] | [ ] | [ ] |
| r/selfhosted | [ ] | [ ] | [ ] | [ ] | [ ] | [ ] | [ ] |
| Nous Research Discord (Hermes Agent) | [ ] | [ ] | [ ] | [ ] | [ ] | [ ] | [ ] |
| Proxmox forum | [ ] | [ ] | [ ] | [ ] | [ ] | [ ] | [ ] |
| Fedora Discussion | [ ] | [ ] | [ ] | [ ] | [ ] | [ ] | [ ] |
| Universal Blue Discourse | [ ] | [ ] | [ ] | [ ] | [ ] | [ ] | [ ] |
| Devlog or project page (not a recruiting channel; exposed) | none | [ ] | [ ] | [ ] | [ ] | [ ] | [ ] |
| Referral or other | [ ] | [ ] | [ ] | [ ] | [ ] | [ ] | [ ] |
| **Total** | [ ] | [ ] | [ ] | [ ] | [ ] | [ ] | [ ] |

Before publication, merge every channel row with fewer than 3 interviewed participants into "Referral or other".

**Screener base rate (Q8):** [n] of [eligible] eligible respondents ([%]) reported that the agent could read personal notes, credentials, or email in the past 30 days. Self-reported; not used for selection.

**Exposure (Q11 and Q12):** [n] of [eligible] eligible respondents and [n] of N participants were exposed.

### C2. Sample composition (counts, not per-person rows, to limit re-identification)
| Attribute | Count (of N = [ ]) |
|---|---|
| Hermes Agent | [ ] |
| OpenCode | [ ] |
| Claude Code | [ ] |
| Codex | [ ] |
| OpenClaw | [ ] |
| Other agents | [ ] |
| Model mix: local only / cloud only / both | [ ] / [ ] / [ ] |
| Host class H1 / H2 / H3 / H4 / H5 / H6 / H7 | [ ] / [ ] / [ ] / [ ] / [ ] / [ ] / [ ] |
| Messaging-gateway mode in use | [ ] |
| Professional field: security (screener Q9) | [ ] |
| Disconfirming quota (stopped using an agent) | [ ] |

### C3. Kill criterion 2 evaluation
| Measure | Value | Threshold | Result |
|---|---|---|---|
| Completed interviews counted in N (A.5) | [ ] | 8 to 10 | [met, or reduced sample] |
| P1-Linux-ready count (gating by default, FD-2) | [ ] | fewer than 3 triggers | [triggered or not] |
| P1-strict count (reported; adds H4 and H6) | [ ] | none | |
| P1-broad count (reported) | [ ] | none | |
| P1-capability only (reported) | [ ] | none | |
| Prong 2 codes C / F / D / X / U (all N) | [ ] / [ ] / [ ] / [ ] / [ ] | | |
| Prong 2 gating reading, denominator 1: (D + X + F) of (N minus U) | [ ] of [ ] | majority triggers | [triggered or not] |
| Prong 2 gating reading, denominator 2: (D + X + F) of (P1-gating participants minus U) | [ ] of [ ] | majority triggers | [triggered or not] |
| Designed-workflow reading: (D + X) over each denominator | [ ] ; [ ] | reported only; gates only if chosen in the published freeze (FD-3) | |
| Exposure check: prong-1 gating count, and prong-2 gating figures, without exposed participants | [ ] ; [ ] ; [ ] | reported only | [same or different] |
| Prong 2 codes resting only on E1 evidence | [ ] | more than 2 means weak evidence | [ ] |
| **Kill criterion 2** | | either prong triggers | **[TRIGGERED or NOT TRIGGERED]** |

Action if triggered (concept document): re-scope to a benchmark-first project for security researchers before building Phase 1. No override is possible (FD-6). Record the decision in the go/no-go record (Phase 0 checklist, section 6).

### C4. Where private data lives (answers the plan's open question)
| Category | Local file | API | Both | Not used | Configured access only | Actual use in past 30 days |
|---|---|---|---|---|---|---|
| M1 Notes or documents | [ ] | [ ] | [ ] | [ ] | [ ] | [ ] |
| M2 Credentials or secrets | [ ] | [ ] | [ ] | [ ] | [ ] | [ ] |
| M3 Email, messages or calendar | [ ] | [ ] | [ ] | [ ] | [ ] | [ ] |
| M4 Private source code | [ ] | [ ] | [ ] | [ ] | [ ] | [ ] |
| M5 API services named (by class) | | [ ] | | | [ ] | [ ] |

Implications to state: default directory rules and xattr labels in `policy/default.yaml`; which API destinations need source labels in v0.1; whether the contingent cut "source labels moved to v0.2" is acceptable.

### C5. Destination comfort (informs default clearances and the cloud-for-private question)
| Destination class | Acceptable for private material | Not acceptable | Unsure | Contacted in past week |
|---|---|---|---|---|
| Local model | [ ] | [ ] | [ ] | [ ] |
| Cloud model provider | [ ] | [ ] | [ ] | [ ] |
| Code host | [ ] | [ ] | [ ] | [ ] |
| Package registry | [ ] | [ ] | [ ] | [ ] |
| Search engine or web fetch | [ ] | [ ] | [ ] | [ ] |
| Paste or file-sharing site | [ ] | [ ] | [ ] | [ ] |
| Chat or messaging service | [ ] | [ ] | [ ] | [ ] |
| Email or document API | [ ] | [ ] | [ ] | [ ] |
| Unknown destinations discovered later (participants) | | | | [ ] |

### C6. Hosting and image readiness
| Measure | Count (of N) |
|---|---|
| Linux-ready host (H1 to H3) | [ ] |
| Proxmox VM / Proxmox LXC | [ ] / [ ] |
| Rebuilt the agent host in the past 6 months | [ ] |
| Has booted a downloaded VM image or used an image-based Linux system | [ ] |
| Uses Mac-to-VM shared folders (virtiofs / 9p) | [ ] / [ ] |
| VPS-only (outside v0.1 claims) | [ ] |

### C7. Permission prompts and blocked access
| Measure | Count (of N) |
|---|---|
| Prompts fully on (observed) | [ ] |
| Partial auto-approve (observed) | [ ] |
| Prompts off (observed) | [ ] |
| Denied at least one prompt in the last session | [ ] |
| Has split work into separate sessions on purpose | [ ] |
| Past episode of blocked network access mid-task | [ ] |
| Of those, continued / restarted / bypassed / gave up | [ ] / [ ] / [ ] / [ ] |
| Of those, the blocked session had already read private material | [ ] |

### C8. Surprise incidents mapped to the threat model
| Incident class | Count | Threat-model adversary or benchmark channel class |
|---|---|---|
| Prompt injection | [ ] | Prompt injector |
| Malicious or compromised tool, skill or MCP server | [ ] | Malicious skill, MCP server, plugin or hook |
| Deceptive agent behaviour (spoofed calls, log tampering) | [ ] | Agent or model acting deceptively; evidence class |
| Laundering across sessions or tools | [ ] | Second agent or later session; laundering classes |
| Unexpected destination | [ ] | Containment and label channel classes |
| Accidental destructive action | [ ] | Not held in v0.1 (non-guarantee); Phase 5 theme |
| Cost overrun | [ ] | Out of scope (quotas only) |
| Other or none | [ ] | |

### C9. Agents and profile priorities
| Agent | Participants | Planned status |
|---|---|---|
| Hermes Agent | [ ] | v0.1 profile |
| OpenCode | [ ] | v0.1 profile |
| Claude Code | [ ] | Phase 6 profile |
| OpenClaw | [ ] | Phase 6 profile, short sessions only |
| Codex | [ ] | Not planned; generic `tidemark run` only |
| Other | [ ] | Generic `tidemark run` only |

### C10. Open questions mapped to evidence
| Open question (source) | Evidence summary | Confidence (low, medium, high) | Implication for the plan |
|---|---|---|---|
| Where does private data live: local files or APIs? (plan; concept) | [from C4] | [ ] | [default labels; source-label list] |
| Will users accept denied public web after reading private data, or must planner/worker ship in the MVP? (concept) | [from C3, C7] | [ ] | [keep or move the planner/worker experiment] |
| Clear a cloud provider for private data by default, or local-only? (concept) | [from C5] | [ ] | [default policy; hero message] |
| Is a headless VM appliance acceptable; will people boot an image? (concept) | [from C6] | [ ] | [image versus RPM emphasis] |
| Are the v0.1 exclusions right: VPS, messaging gateways, OpenClaw, macOS virtiofs? (concept) | [from C2, C6] | [ ] | [ ] |
| Is per-action approval failing as claimed? (plan, research conclusions) | [from C7] | [ ] | [ ] |
| Is the core asset the image or the specification plus benchmark? (concept) | [indirect, from C6 and commitment signals] | [ ] | [OpenShell adapter timing] |
| Name and trademark; DCO or CLA; reviewer names (plan) | Not answered by interviews | n/a | n/a |

### C11. Disconfirming evidence
List every finding that contradicts the plan or the concept, even if it appeared in one interview only. Generalise it so that it does not identify the participant.
- [ ]

### C12. Commitment signals
| Signal | Count (of N) |
|---|---|
| Agreed to future contact for testing | [ ] |
| Offered a referral | [ ] |
| Other concrete offers | [ ] |

---

## D. Findings (narrative)
1. [Finding, with the tables that support it]
2. [ ]
3. [ ]

**Recommended input to the go/no-go decision:** [one paragraph; state whether kill criterion 2 triggered under the gating reading, what the reported readings show, and what the founder should weigh]

## E. Limitations
- N = [ ]; a coarse signal, not a statistical estimate. [State "reduced sample" if A.5 applied.]
- The interviewer is the founder and has a stake in the result.
- Recruiting channels favour local-first users; see C1.
- Participants who had read the project's devlog or repository know what the interviews test; see the exposure check in C3.
- Prong 2 relies partly on an anchored hypothetical; see the E1 count in C3.
- [Any other limitation observed]

## F. Publication checklist (complete before committing the filled synthesis)
- [ ] Only the executive summary, sections A and C to F, and the Tone Check are published. No filled section B record is in the repository.
- [ ] The freeze of section A was published before P01, and the link is recorded in section A.
- [ ] No names, handles, email addresses, employers, cities, hostnames, IP addresses, real paths, hardware serials, kernel versions, interview dates or secrets.
- [ ] Every channel, agent, host-class or incident-class cell with fewer than 3 participants is merged into "Other" or reported as "fewer than 3" (C3 gating counts excepted).
- [ ] No attribute combination that is unique to one participant, unless generalised.
- [ ] Every quote is short, scrubbed and attributed only to a participant ID.
- [ ] Withdrawals received by Thursday 5 November 2026 are excluded. Later withdrawals are recorded for removal from later versions.
- [ ] No cloud AI tool processed raw notes, recordings, transcripts, contact details or section B records.
- [ ] Raw notes, recordings, transcripts and section B records remain outside the repository.
- [ ] Recordings deleted (no later than 14 days after each interview).
- [ ] Contact details of screener respondents who were not interviewed deleted on Wednesday 4 November 2026.
- [ ] Deletion deadline recorded: raw notes, section B records and contact keys by 8 December 2026; opted-in contact handles kept without participant IDs for at most 12 months.
- [ ] Coding rules applied exactly as frozen; deviations listed in A.6.
- [ ] Founder has reviewed the whole document.
- [ ] Tone Check completed below.

## Tone Check
Complete before publication. The "Warm but Authoritative" brand guideline document was not available when this template was drafted; score against it when filling in the synthesis.

| Criterion | Score (1 to 5) | Notes |
|---|---|---|
| Warmth (respect for participants) | [ ] | [ ] |
| Authority (claims supported by tables) | [ ] | [ ] |
| Clarity | [ ] | [ ] |
| Honesty about limitations | [ ] | [ ] |
