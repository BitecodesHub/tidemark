# Observational Interview Guide (v1)

| Field | Value |
|---|---|
| Status | Phase 0 draft, revised 16 October 2026 after review. Not yet piloted. No interview has taken place. |
| Version | v1, prepared 2026-10-16 for use from Monday 19 October 2026 |
| Owner | Founder (interviewer). AI-assisted drafting; the founder reviews and adapts before use. |
| Interview window | Wednesday 21 October to Sunday 1 November 2026 (weeks 1 and 2). Spill-over to Wednesday 4 November is permitted only if fewer than 8 interviews are complete, and is recorded as a deviation. |
| Decision date | Kill criterion 2 is evaluated at the week-3 gate, Sunday 8 November 2026. |
| Companion files | `2026-10-19_Interview-Screener_v1.md`, `2026-10-19_Synthesis-Template_v1.md`, `../2026-10-16_Phase0-Checklist_v1.md` |
| Sources | `../2026-10-16_Project-Plan_v1.md`; `../research/2026-10-16_Tidemark-Concept_v1.md` (sections "First users", "Kill criteria", "Open questions for the founder") |

## Executive summary
- The founder holds 8 to 10 observational interviews of 45 minutes each with homelab, Proxmox and dedicated-VM users who run coding or research agents (Hermes Agent, OpenCode, Claude Code, Codex, OpenClaw) over local files.
- Every question asks about past behaviour and asks the participant to show, not to predict. Unprompted questions come first; questions about private data come later; the one hypothetical probe comes last. The interviewer never describes Tidemark before minute 40, and only if the participant wants to hear it.
- The interviews decide kill criterion 2 at the week-3 gate and answer the plan's open question on where private data actually lives.
- No recording happens without explicit consent. The repository holds only the aggregate synthesis. Per-interview records, raw notes, recordings and contact details stay outside the repository, are processed only on the founder's machine, and are deleted on a fixed schedule.

## 1. What the interviews must decide

### 1.1 Kill criterion 2 (verbatim from the concept document)
> Week 3. Trigger: of the 8 to 10 observational interviews, fewer than 3 participants run agents over local private files or secrets on a machine or VM they control, or a majority would abandon a session rather than accept denied public-web access after it reads private data. Action: re-scope to a benchmark-first project for security researchers before building Phase 1.

The coding rules that turn interview evidence into these two counts are pre-registered in section A of the synthesis template. The founder freezes them before the first interview, publishes the freeze with a timestamp, and logs any later change as a deviation.

### 1.2 Open questions the interviews inform
| Question (source) | Interview section that informs it |
|---|---|
| Where do target users' private data live: local files, or APIs such as Gmail, IMAP and Notion? (plan, open questions; concept, founder questions) | 6.6 Where private data lives |
| Will first users accept denied public-web access after reading private data, or must the planner/worker split ship in the MVP? (concept, founder questions) | 6.8 Prompts and blocked access; 6.9 Scenario probe; 6.7 Destinations |
| Should the default policy let users clear a cloud provider for private data, or keep private data local-only? (concept, founder questions) | 6.4 Agent setup (model mix); 6.7 Destinations |
| Is a headless VM appliance acceptable, and will people boot an image at all? (concept, founder questions and risks) | 6.5 Where the agent runs |
| Do v0.1 exclusions match reality: VPS hosting, messaging-gateway modes, OpenClaw, macOS virtiofs shared folders? (concept, "First users") | 6.4 and 6.5 |
| Is per-action approval failing as the research claims (about 93% of prompts approved; Omarchy disables prompts for 11 of 15 agents)? (plan, research conclusions) | 6.8 Prompts and blocked access |

Questions the interviews do not answer: the project name and trademark, DCO or CLA, and the names of pre-launch reviewers.

## 2. Who to interview

### 2.1 Target profile
- Adults (18 or older) who, in the past 30 days, ran a coding or research agent that can read files or run commands.
- Agents of interest: Hermes Agent, OpenCode, Claude Code, Codex, OpenClaw. Other agents are acceptable when the participant also meets the hosting profile.
- Hosting of interest: a physical Linux machine they own, a VM or container on hardware they own (for example Proxmox), or a dedicated cloud VM whose image they choose.

### 2.2 Sample targets (soft quotas)
| Attribute | Target |
|---|---|
| Completed interviews | 8 minimum, 10 maximum |
| Hermes Agent or OpenCode as the most-used agent (screener Q4) | at least 3 |
| Proxmox or other homelab hypervisor host (screener Q5) | at least 2 |
| Dedicated Linux VM, local or cloud, own image (screener Q5) | at least 2 |
| Agent runs directly on macOS or Windows (screener Q5) | at most 3; under the default coding rules they do not count toward the gating prong-1 count (synthesis template, FD-2) |
| OpenClaw as the most-used agent (Q4), or tasks sent through a messaging app (Q6) | at most 2 |
| Professional field: security (Q9) | at most 2 (they know the vocabulary and may bias answers) |
| Stopped using an agent in the past 12 months (Q2), or restricted it because of a privacy or safety concern | at least 1 if available, at most 2 (disconfirming evidence) |
| Personal acquaintances of the founder (Q10) | none |
| Recruited from the devlog or project page, or already read the project's material (Q11, Q12) | not sought; flagged as exposed if they arrive |

### 2.3 Selection bias to record, not to hide
- The recruiting channels favour local-first enthusiasts, which inflates the first prong of kill criterion 2. The screener therefore does not select on private-data access (see the screener, section 4), and the synthesis reports the screener base rate next to the interview counts.
- Record the recruiting channel for every participant.
- The devlog and the repository describe the design and publish kill criterion 2, so anyone who has read them knows which answers keep the project alive. The devlog does not recruit. Participants who arrive from it, or who say they have read the project's material, are flagged as exposed, and the synthesis reports the results with and without them.

## 3. Method rules
These rules adapt the past-behaviour interviewing approach described in Rob Fitzpatrick's book "The Mom Test". They are paraphrased here; consult the book for the full method.

1. **Ask about their work, not about the idea.** The subject is how they run agents today.
2. **Ask about specific past events.** "Tell me about the last time..." beats "Do you usually...". "Would you..." questions are not evidence.
3. **Ask to be shown.** A configuration on screen outranks a description; a description of a specific event outranks an opinion.
4. **Treat compliments and enthusiasm as noise.** Record facts, workarounds, time and money already spent.
5. **Follow emotion and workarounds.** When the participant sighs, laughs or says "annoying", ask what happened and what they did about it.
6. **Let silence work.** Wait after each answer before the next question.
7. **Do not teach.** Do not correct the participant's security understanding during the interview.
8. **No pitch before minute 40.** If the participant asks what the founder is building, answer: "I will happily tell you at the end, so that I do not colour your answers."
9. **Unprompted before prompted.** Ask open questions about surprises and setup before any question that names private data, secrets or blocked access.

### 3.1 Words the interviewer does not use before minute 40
Tidemark, label, taint, clearance, exfiltration, leak, canary, sandbox (unless the participant uses it first), "security product", "would you use", "would you pay".

### 3.2 Leading versus observational phrasing
| Avoid (leading or hypothetical) | Use (observational) |
|---|---|
| "Are you worried about agents leaking your data?" | "Which folders could your agent read in the last session you ran? Can you show me where that is set?" |
| "Would you want the agent to lose web access after reading secrets?" | "Tell me about the last time a tool blocked network access in the middle of a task. What did you do next?" |
| "Do you think permission prompts are useless?" | "In your last session, how many prompts did you see, and which ones, if any, did you deny?" |
| "Would you trust a cloud model with your notes?" | "Which model handled the last task that touched your notes? How did you choose it?" |

## 4. Consent and privacy protocol

### 4.1 Before the call
- Send the consent summary (section 4.5) with the scheduling message.
- Ask the participant to prepare their agent configuration with every API key, token, password, hostname and username removed or covered. Ask to see where a secret lives (environment variable, file, keychain), never the secret itself.
- Ask the participant to share one application window, not the whole screen, with notifications silenced.
- Disable meeting-platform AI note-takers, automatic recording and cloud transcription before the call starts. This follows the project's local-only default for private data.

### 4.2 At the start of the call (scripted; read it)
"Thank you for your time. Before we start, four points. First, I am an independent open-source developer researching how people run AI agents on their own machines. I am working on a project in this area, and I will describe it only at the end if you want to hear about it. Second, I will take written notes. I will record audio only if you say yes now; the recording stays on my own machine, is transcribed locally, and is deleted within 14 days. My notes also stay on my machine, and no cloud AI service processes them. Third, I will not publish your name, handle, employer, location or anything that identifies you; only an anonymised summary of all the conversations goes into the public project repository. Fourth, you can skip any question or stop at any time, and you can ask me to delete your notes at any time before I delete them anyway on 8 December 2026. If you ask by 5 November 2026, your answers are also left out of the published summary. After that date, I still delete your notes and remove your answers from any later version, but the counts in the summary published around 8 November stay in the repository history. Please share one window rather than your whole screen, and please do not show me any password, key or token; if one appears on screen, I will tell you and I will not write it down. Do you agree to take part? And may I record audio, yes or no?"

Record in the raw notes: consent given (yes or no), recording permitted (yes or no), and the time. Do not proceed without a clear yes to taking part.

### 4.3 During the call
- Do not take screenshots or photos of the participant's screen.
- Generalise identifying details as you write: "personal notes vault in the home directory" instead of a real path; "a mid-range GPU" instead of a model number; "a small consultancy" instead of an employer name.
- If a secret, a private message or a third party's personal information appears on screen, say so, look away from it in the notes, and do not record it.
- Do not ask about the content of private documents; ask only about their kind and location.

### 4.4 After the call
| Item | Where it lives | Retention |
|---|---|---|
| Raw notes keyed by participant ID (P01 to P10) | Encrypted local storage on the founder's machine, outside the repository | Deleted by 8 December 2026 (gate plus 30 days) |
| Audio recording (only with consent) | Encrypted local storage outside the repository | Deleted when the notes have been checked against it, and no later than 14 days after the interview |
| Transcript (only with consent) | Produced locally (for example with a local speech-to-text model), stored with the raw notes | Deleted with the raw notes |
| Per-interview record (synthesis template, section B) | Encrypted local storage with the raw notes; never committed | Deleted with the raw notes by 8 December 2026 |
| Contact key (participant ID to contact handle) | A separate encrypted file, outside the repository | Deleted by 8 December 2026. If the participant opted in to future contact, only the contact handle is kept, without the participant ID or any notes, for at most 12 months after the interview |
| Contact details of screener respondents who were not interviewed | A separate encrypted file, outside the repository | Deleted when recruiting closes, Wednesday 4 November 2026 |
| Published synthesis | `docs/interviews/` in the repository | Kept; contains the coding rules and aggregate tables only (synthesis template, sections A and C to F), with no personal identifiers |

- Write the per-interview record from the raw notes within 24 hours of each interview, while memory is fresh.
- **AI tools.** Raw notes, recordings, transcripts, contact details and per-interview records are processed only by tools that run locally on the founder's machine. No cloud AI service sees them. Cloud AI tools may process the published synthesis after its publication checklist passes.
- Never commit raw notes, recordings, transcripts, per-interview records, screenshots, configuration files, names, handles, email addresses, employers, cities, hostnames, IP addresses or secrets to the repository.
- Small samples are easy to re-identify. Published tables merge any cell with fewer than 3 participants into "Other" (synthesis template, "Handling rules"), and no kernel versions, hardware models or per-person rows are published.
- **Withdrawal.** A request received by Thursday 5 November 2026 removes the participant's answers from the synthesis before it is published. A later request deletes the raw notes, the per-interview record and the contact key at once, and removes the participant from any later version of the synthesis. The first published synthesis (around 8 November) is not changed, and its counts remain in the repository history.
- Founder decision: confirm whether any data-protection law applies to participants in your jurisdiction or theirs (for example, the GDPR for participants in the EU). This guide collects minimal data in either case; it is not legal advice.

### 4.5 Consent summary (sent with the scheduling message)
"This is a 45-minute conversation about how you run AI agents today. Taking part is voluntary. I take written notes; audio recording happens only if you agree at the start, stays on my machine and is deleted within 14 days. My notes stay on my machine, and no cloud AI service processes them. Your name and other identifying details are never published. An anonymised summary of all conversations will be published in an open-source project repository. You can skip questions or stop at any time, and you can ask me to delete your notes at any time before I delete them on 8 December 2026. Requests made by 5 November 2026 also keep your answers out of the published summary; after that date, the counts in the summary published around 8 November stay in the repository history. Please share one window rather than your whole screen, and remove keys, tokens and passwords from anything you show me. There is no payment; I will send you the published summary if you would like it."

## 5. Recruiting

### 5.1 Channels
| Channel | Notes |
|---|---|
| r/LocalLLaMA | Read the subreddit rules on self-promotion and research requests before posting; message the moderators first if the rules require it. |
| r/selfhosted | Same as above. Frame the post as research, not as a product announcement. |
| Hermes Agent community | The Hermes Agent repository (NousResearch/hermes-agent) has GitHub Discussions disabled and about 15,800 open issues and 33,700 open pull requests (checked 16 October 2026); do not post in the issue tracker. Use the Nous Research Discord linked from the project README, in the channel its rules designate. |
| Proxmox community | The Proxmox Support Forum (forum.proxmox.com); use the off-topic or community area if one exists, per forum rules. |
| Fedora and Universal Blue | Fedora Discussion (discussion.fedoraproject.org) and the Universal Blue Discourse (universal-blue.discourse.group); use the categories their rules designate for community or research requests. |
| OpenCode community | Optional; use the official community channel named in the project's README. |
| Devlog or project page | Not a recruiting channel. Devlog 1 does not invite interviewees. If a respondent names it (screener Q11), flag the participant as exposed, and report the results with and without exposed participants. |
| Referrals | Ask each participant at the end of the interview for one or two people who run agents differently from them. Referrals complete the screener like any other respondent. |

- Posting is a public action. The founder posts manually after reading each community's rules; no automated posting.
- Do not offer cash incentives in Phase 0; a payment would require collecting payment details. Offer the published synthesis and, with consent, an acknowledgement. Founder decision: confirm.
- Track per channel: posts made, screener responses, qualified respondents, interviews held.

### 5.2 Recruiting post (draft; external-facing)
Title: "Research: how do you run AI coding or research agents on your own hardware? (45-minute chat, nothing to buy)"

"I am an independent open-source developer, and I am trying to understand how people actually run agents such as Hermes Agent, OpenCode, Claude Code, Codex or OpenClaw on machines they control: homelab servers, Proxmox VMs, dedicated Linux VMs or workstations.

I am looking for 8 to 10 people for a 45-minute voice or video conversation between 21 October and 1 November 2026. I will ask about your setup and about things that have actually happened, not about a product. There is nothing to buy and no sales follow-up.

Your name and identifying details will never be published. Notes are anonymised, and only an aggregated summary will be posted publicly; I will send it to everyone who takes part.

If you are interested, please answer a few short questions (about four minutes) here: [screener link]. To keep your setup private, please use the form rather than replying here with details of how you run your agents. Thank you."

### 5.3 Scheduling message (draft; external-facing)
"Thank you for answering the questions. I would like to talk with you for 45 minutes. Could you pick a time here: [scheduling link]? The consent summary is below. Before the call, please remove or cover any keys, tokens, passwords, hostnames and usernames in the configuration you might show me, and plan to share one window with notifications silenced; I only need to see how things are set up, never the secrets themselves."

## 6. Session plan (45 minutes)
The order is deliberate: open questions about surprises and setup come before any question that names private data, and the one hypothetical probe comes last, so that earlier answers are not primed.

| Time | Section | Purpose |
|---|---|---|
| 0:00 to 0:03 | 6.1 Opening and consent | Consent script; recording decision |
| 0:03 to 0:06 | 6.2 Context | What they use agents for, and how often |
| 0:06 to 0:11 | 6.3 The last surprise | Real incidents, unprompted |
| 0:11 to 0:18 | 6.4 Show me the agent setup | Configuration, models, tools, MCP servers |
| 0:18 to 0:22 | 6.5 Where the agent runs | Host, VM or container; image habits |
| 0:22 to 0:29 | 6.6 Where private data lives | Local files versus APIs; actual use; kill criterion prong 1 |
| 0:29 to 0:34 | 6.7 Destinations | Which services receive what |
| 0:34 to 0:38 | 6.8 Prompts and blocked access | Past behaviour; kill criterion prong 2 |
| 0:38 to 0:40 | 6.9 Scenario probe | The one anchored hypothetical; prong 2 fallback |
| 0:40 to 0:45 | 6.10 Close | Commitment, referrals, optional description of the project |

If time runs short, protect sections 6.6, 6.7, 6.8 and 6.9, because they decide the kill criterion. Shorten 6.2 and 6.4 first, then 6.5. Keep at least two minutes for 6.3, because it is the only unprompted incident question.

### 6.1 Opening and consent (3 minutes)
- Read the script in section 4.2.
- Use the participant ID (P01 to P10) in the notes from the first line; never write the participant's name or handle in the notes.

### 6.2 Context (3 minutes)
- "What did you last use an agent for? Walk me through that task."
- "Which agents have you used in the past month, and roughly how often?"
- "Is this personal use, work use, or both?" (If work use is under an employer policy, ask the participant to describe only what they are free to share.)

### 6.3 The last surprise (5 minutes)
- "Tell me about the last time an agent did something you did not expect."
- "What did you notice first? What did you do next? What did you change afterwards?"
- "Did it cost you anything: time, money, data, a credential you had to rotate?"
- If there is no surprise: "What is the closest it came?" Do not supply examples.

### 6.4 Show me the agent setup (7 minutes)
Ask the participant to share the configuration window, or to read it out, with secrets covered.
- "Could you show me how your main agent is configured? Where does that configuration live?"
- "Which models does it use? Which tasks go to a local model and which to a cloud model? How did you decide that?"
- "Which tools, plugins, skills or MCP servers are enabled? Which of them did you install from someone else?"
- "Which folders or files can it read? Where is that set? When did you last change it, and why?"
- "Which approval or permission settings are on right now?" (Observe whether auto-approve, "yolo" or equivalent modes are enabled.)
- "Do you connect it to any chat or messaging app, such as Telegram, Discord, Slack, WhatsApp or Signal?"

Record: agents; model mix; tool and MCP list (as categories); readable paths (generalised); approval settings as observed; messaging modes.

### 6.5 Where the agent runs (4 minutes)
- "Where does the agent process actually run: this laptop, a server at home, a Proxmox VM or container, a cloud VM, a VPS?"
- "If it runs in Proxmox, is it a VM or an LXC container? Show me how it is defined." (An LXC container cannot load its own kernel security programs; a VM can boot a custom image.)
- "When did you last rebuild or reinstall that machine? How did you do it?"
- "Have you ever booted a downloaded VM image or switched to an image-based Linux system? What happened?"
- "If you share files between a Mac host and a Linux VM, how?" (Shared folders over virtiofs or 9p matter for later work.)

Do not record kernel versions or hardware models.

### 6.6 Where private data lives (7 minutes)
- "Has the agent ever touched something you would not want published? Tell me about the last time."
- "Show me where your notes, documents or personal files live on this machine." (Generalise the location in notes.)
- "Where do keys and credentials for this machine live: files, environment variables, a password manager, a keychain?" (Location only.)
- "Does the agent read email, calendars, Notion, Google Drive or similar services? Through what: an MCP server, an API key, an export, IMAP?"
- "Which of these has the agent actually read, or have you pointed it at, in the past 30 days? Can you show me where that is visible, for example in its history or logs?"

Record for each category of private material: local file, API, or both; whether the agent is configured with access; whether it actually read the material or was pointed at it in the past 30 days; evidence strength (observed, specific past event, or general statement). Configured access alone is not actual use (synthesis template, A.2).

### 6.7 Destinations (5 minutes)
- "Which services did your agent contact in the last week? If you can, show me where that is visible: logs, a proxy, a firewall."
- For each destination class the participant names (local model, cloud model provider, code host, package registry, search engine, paste or sharing site, chat app, email or document API), ask: "Have you ever knowingly sent private material there through the agent? Is that acceptable to you? Why?"
- "Is there any service your agent contacts that you did not know about until later? How did you find out?"

Record a three-way classification per destination class: acceptable for private material, not acceptable, or unsure, with the reason in the participant's words.

### 6.8 Prompts and blocked access (4 minutes)
- "In your last session, roughly how many permission prompts did you see? Which, if any, did you deny?"
- "Have you ever turned prompts off, or switched to an auto-approve mode? Tell me about the moment you decided that."
- "Tell me about the last time a tool, a firewall or a setting blocked network access in the middle of a task. What did you do next?" Follow up: "Before that happened, had that session read any of the private material we talked about?"
- "Have you ever split one job into separate sessions on purpose, for example one session for private material and another for web research? Walk me through it."

Code prong 2 from these answers first, using the mapping table in the synthesis template, section A.3.

### 6.9 Scenario probe (2 minutes)
First ask: "Is there a session you remember in which your agent opened your private notes or keys?" If no such session exists, skip the probe and code U.

Otherwise, read the probe exactly:
"Think of the last session in which your agent opened your private notes or keys. Suppose that, from that moment on, any request from that session to the public web simply failed, with a message explaining why. What would you have done next?"
- Follow up once with "What would that have cost you?" Do not suggest options.
- Code the answer after the call as C (continue without web), F (start a fresh session for the web part), D (disable or bypass the restriction, or switch tools), X (abandon the task) or U (unclear), marked E1. The probe decides the prong-2 code only when no past-behaviour source in section 6.8 applies (synthesis template, section A.3).

### 6.10 Close (5 minutes)
1. "Is there anything about how you run agents that I should have asked about?"
2. Commitment signals (record exactly what the participant agrees to; do not push):
   - "Could I contact you again when I have something to test in your setup, probably in December?"
   - "Is there one person who runs agents quite differently from you whom you could introduce me to?" (Referrals complete the screener.)
3. Optional description of the project, only if the participant wants it: two minutes, neutral, from the concept's claim and "not claimed" list. Record the reaction, but weight it low; reactions to a description are opinions about the future.
4. Thank the participant, confirm the withdrawal dates (5 November 2026 to be left out of the published summary; notes deleted by 8 December 2026) and how they will receive the summary.

## 7. Note-taking
- Use the per-interview record in the synthesis template, section B, filled from the raw notes within 24 hours and stored with the raw notes outside the repository.
- Separate three kinds of statement in the notes: **observed** (shown on screen), **past event** (a specific occurrence), **opinion** (general or hypothetical).
- Write short paraphrases. Keep a verbatim quote only when the exact words matter, and scrub it before publication.

## 8. Pilot and changes
- P01 is also the pilot. P01 counts toward N and the kill criterion unless the guide changes materially after P01; in that case, report P01 separately, exclude it from N, and log the change in the deviation log (synthesis template, sections A.5 and A.6).
- Any change to the coding rules after the freeze is a deviation and must be justified in writing before the next interview.
- If N is below 8 when recruiting closes on Wednesday 4 November, apply the reduced-sample rule in synthesis template section A.5. The threshold of 3 is never lowered.

## 9. Timeline
| Date (2026) | Step |
|---|---|
| Mon 19 Oct | Freeze the coding rules (synthesis template, section A) and publish the freeze (public commit, or SHA-256 in devlog 1). Post in two or three channels after reading their rules. |
| Tue 20 Oct to Fri 30 Oct | Screen respondents; schedule interviews; post in remaining channels if responses are slow. |
| Wed 21 Oct to Sun 1 Nov | Interviews (target: 3 in week 1, 5 in week 2); each per-interview record written within 24 hours. |
| Mon 2 Nov to Wed 4 Nov | Spill-over interviews only if fewer than 8 are complete; recorded as a deviation. |
| Wed 4 Nov | Recruiting closes. Contact details of respondents who were not interviewed are deleted. If N is below 8, apply synthesis template A.5. |
| Thu 5 Nov | Last day for a withdrawal to keep a participant's answers out of the published synthesis. |
| By Fri 6 Nov | Synthesis drafted from the per-interview records; kill criterion 2 evaluated. |
| Sun 8 Nov | Week-3 gate: written go or no-go decision (see the Phase 0 checklist). |
| Tue 8 Dec | Raw notes, transcripts, per-interview records and contact keys deleted. Opted-in contact handles are kept without participant IDs for at most 12 months after the interview. |

## 10. Known limitations of this guide
- Eight to ten interviews give a coarse signal, not a statistically meaningful estimate.
- The founder is the interviewer and has a stake in the outcome. The pre-registered coding rules, the published freeze and the verbatim scenario probe reduce, but do not remove, this bias.
- The scenario probe in 6.9 is hypothetical, even though it is anchored to a real session. The synthesis ranks it below observed and past-event evidence.
- The recruiting channels over-represent local-first users.
- The project is developed in public. Some participants may have read the devlog or repository and know the hypothesis; the screener flags them, but it relies on self-report.

## Tone Check
Scope: the external-facing drafts in sections 4.2, 4.5, 5.2 and 5.3. The "Warm but Authoritative" brand guideline document was not available to the drafter; the scores below apply the two named qualities directly and should be re-scored against the guideline before posting.

| Criterion | Score (1 to 5) | Notes |
|---|---|---|
| Warmth | 4 | Thanks the reader, states there is nothing to buy, offers the summary. Could add one sentence on why their experience matters. |
| Authority | 4 | Specific dates, duration, retention periods and withdrawal deadlines give credibility. No claims about the project are made. |
| Clarity | 4 | One call to action per message. The revised consent script is longer because it states the withdrawal limits honestly; consider reading it slowly. |
| Respect for privacy | 5 | Recording only with consent, local processing only, one window instead of the whole screen, deletion dates, secrets explicitly excluded, public replies discouraged. |
| Neutrality (no pitch) | 5 | The project is mentioned only as existing; no features or vocabulary that would prime answers. |
