# Interview Screener (v1)

| Field | Value |
|---|---|
| Status | Phase 0 draft, revised 16 October 2026 after review. Not yet published as a form. No responses collected. |
| Version | v1, prepared 2026-10-16 for use from Monday 19 October 2026 |
| Owner | Founder. AI-assisted drafting. |
| Used with | `2026-10-19_Interview-Guide_v1.md` (recruiting, consent, session plan) and `2026-10-19_Synthesis-Template_v1.md` (tables C1 and C3) |

## Executive summary
- Fourteen short questions, about four minutes to answer. Questions 1 and 2 decide eligibility; questions 3 to 9 decide priority and quotas; question 10 excludes acquaintances; questions 11 and 12 record the channel and flag respondents who already know the project; questions 13 and 14 handle scheduling.
- The screener selects on agent use and hosting, not on access to private data. Selecting on private data would make the first prong of kill criterion 2 pass by construction. Question 8 is recorded only to report a base rate.
- The form collects no name. The contact field (question 14) is stored separately from the answers. Contact details of respondents who are not interviewed are deleted when recruiting closes on Wednesday 4 November 2026.

## 1. Form implementation notes
- Use a form tool that can export answers without contact details, or keep two forms: one for questions 1 to 13 (anonymous) and one for question 14 (contact), linked by a random respondent code that the form shows to the respondent.
- Disable collection of email addresses, IP addresses and sign-in identity in the form tool where it offers that setting.
- Store exports outside the repository. Only aggregate counts (tables C1 and C2 of the synthesis template) enter the repository.
- Founder decision: choose the form tool. A self-hosted form keeps answers local, in line with the project's local-only default.

## 2. Introduction text shown to respondents (external-facing)
"Thank you for your interest. These questions take about four minutes. They help me find people whose setups differ from one another; there are no right answers. I do not ask for your name. If you qualify, I will ask for a contact method at the end so that I can schedule a 45-minute conversation between 21 October and 1 November 2026. Your answers are stored privately, and only aggregate counts are published. If I do not invite you to a conversation, I delete your contact details when recruiting closes on 4 November 2026. You can ask me to delete your answers at any time by quoting the code shown at the end of the form; requests made by 5 November 2026 also remove them from the published counts."

## 3. Questions

**Q1. Are you 18 years of age or older?**
- Yes
- No (end the form: "Thank you; this study is limited to adults.")

**Q2. In the past 30 days, have you used an AI agent that can read files or run commands on a computer (for example Hermes Agent, OpenCode, Claude Code, Codex or OpenClaw)?**
- Yes
- No, but I used one and stopped in the past 12 months (please answer the remaining questions for the time when you last used one)
- No (end the form: "Thank you; this study is about people who use such agents today.")

**Q3. Which agents have you used in the past 30 days, or when you last used one? (Select all that apply.)**
- Hermes Agent
- OpenCode
- Claude Code
- Codex (CLI or app)
- OpenClaw
- Other (please name)

**Q4. Which one agent do you use most, or did you use most? (Select one.)**
- Hermes Agent
- OpenCode
- Claude Code
- Codex (CLI or app)
- OpenClaw
- Other (please name)

**Q5. Where does the agent process usually run? (Select all that apply.)**
- A physical Linux machine I own (desktop, laptop or server)
- A VM on hardware I own (for example Proxmox, libvirt, ESXi, Hyper-V)
- A container on hardware I own (for example a Proxmox LXC container, Docker or Podman)
- A cloud VM where I choose the operating-system image
- A rented VPS from a hosting provider
- Directly on macOS
- Directly on Windows or in WSL
- A machine managed by my employer
- Other (please describe)

**Q6. Do you send the agent tasks through a chat or messaging app, such as Telegram, Discord, Slack, WhatsApp or Signal?**
- Yes, regularly
- Occasionally
- No
- Not sure

**Q7. Which models does the agent use?**
- Only local models (running on hardware I control)
- Only cloud models (through a provider's API or subscription)
- Both local and cloud models
- Not sure

**Q8. In the past 30 days, or when you last used one, which kinds of material could the agent read? (Select all that apply.)**
- Public or open-source code
- Documentation or web pages
- Source code I or my team own
- Personal notes or documents
- Credentials, keys or configuration files with secrets (for example SSH keys or `.env` files)
- Email, calendar or messages
- Documents in online services (for example Notion or Google Drive)
- None of these, or not sure

**Q9. Which best describes your professional field?**
- Software development
- Operations or infrastructure
- Security
- Research or academia
- Other
- Prefer not to say

**Q10. Do you know the interviewer personally (friend, colleague or collaborator)?**
- Yes
- No

**Q11. Where did you see this invitation? (Select one.)**
- r/LocalLLaMA
- r/selfhosted
- Nous Research Discord
- Proxmox Support Forum
- Fedora Discussion
- Universal Blue Discourse
- OpenCode community
- The interviewer's devlog or project page
- Someone referred me
- Other (please name)

**Q12. Before seeing this invitation, had you read anything the interviewer has published about their own project (for example a devlog or a code repository)?**
- Yes
- No
- Not sure

**Q13. Which region and time zone are you in, and when are you usually available between 21 October and 1 November 2026?**
- Region only (for example "Western Europe, UTC+1"); free text for availability.

**Q14. How should I contact you to schedule the conversation?** (Stored separately from your answers. If you are not invited to a conversation, it is deleted when recruiting closes on 4 November 2026. If you take part, it is deleted by 8 December 2026, unless you tick the box below.)
- Free text (for example a Reddit, Discord or forum username, or an email address).
- Optional checkbox: "You may contact me again about testing a future build. Only my contact method is kept, for at most 12 months."

## 4. Selection rules
These rules are applied by the founder in the order shown. They are deliberately independent of Q8.

### 4.1 Eligibility
- Eligible: Q1 = Yes, and Q2 = Yes or "stopped in the past 12 months". Respondents who stopped are scheduled only under the disconfirming quota (section 4.3).
- Not eligible: Q10 = Yes (personal acquaintance). Thank the respondent and do not schedule.
- Respondents whose only host is an employer-managed machine (Q5) are eligible only if they agree to describe personal or freely shareable use; record this in the raw notes.
- Referrals complete this screener like any other respondent; the same rules apply.

### 4.2 Priority
| Priority | Rule (from Q3 and Q5) |
|---|---|
| A | Q5 includes a physical Linux machine they own, a VM or container on hardware they own, or a cloud VM with their own image. |
| B | Q5 is only macOS, Windows or WSL on their own machine; or only a rented VPS. |
| C | Q5 is only an employer-managed machine, or Q3 is only agents outside the list with no Linux hosting. |

Schedule priority A first, then fill the quotas below from B. Use C only to reach the minimum of 8. Within the same priority, schedule respondents who are not flagged as exposed (section 4.4) first.

### 4.3 Quotas (from the interview guide, section 2.2)
| Attribute | Target | Screener field |
|---|---|---|
| Hermes Agent or OpenCode as the most-used agent | at least 3 | Q4 |
| Proxmox or other homelab hypervisor | at least 2 | Q5 (confirmed in the interview) |
| Dedicated Linux VM, local or cloud, own image | at least 2 | Q5 |
| Agent directly on macOS or Windows | at most 3 | Q5 |
| OpenClaw as the most-used agent, or tasks through a messaging app | at most 2 | Q4, Q6 |
| Professional field: security | at most 2 | Q9 |
| Stopped using an agent in the past 12 months | at least 1 if available, at most 2 | Q2 |

### 4.4 Exposure flag
A respondent is flagged as exposed if Q11 is "the interviewer's devlog or project page" or Q12 is "Yes". Exposed respondents know the design and the kill criterion, so their answers may be primed. They remain eligible, but the synthesis reports prong 1 and prong 2 with and without them (synthesis template, section A.3). "Not sure" in Q12 is recorded but does not set the flag.

### 4.5 Why Q8 does not decide selection
Kill criterion 2 asks whether enough real users run agents over local private files or secrets. If the screener admitted only respondents who tick private material in Q8, at least three interviewees would almost certainly qualify, and the gate would test nothing. The screener therefore records Q8 and reports, in table C1 of the synthesis template, the share of eligible respondents who ticked personal notes, credentials, or email. The interviews then verify actual use by observation; configured access alone does not count (synthesis template, section A.2). Founder decision: confirm this rule before the first post (synthesis template, FD-5).

## 5. Records kept
| Record | Location | Retention |
|---|---|---|
| Answers to Q1 to Q13 | Form export, outside the repository | Deleted by 8 December 2026; aggregate counts kept |
| Answer to Q14, respondents not interviewed | Separate encrypted file, outside the repository | Deleted when recruiting closes, Wednesday 4 November 2026 |
| Answer to Q14, participants | Separate encrypted file, outside the repository | Deleted by 8 December 2026. If the participant opted in, only the contact method is kept, without answers or participant ID, for at most 12 months after the interview |
| Aggregate counts (respondents, eligible, priority, channel, Q8 base rate, exposure) | Synthesis, tables C1 and C2, in the repository | Kept; channels with fewer than 3 interviewed participants merged into "Referral or other" |

## Tone Check
Scope: the respondent-facing introduction (section 2), the question wording (section 3) and the end-of-form messages. The "Warm but Authoritative" brand guideline document was not available to the drafter; re-score against it before publishing the form.

| Criterion | Score (1 to 5) | Notes |
|---|---|---|
| Warmth | 4 | Thanks the respondent; "there are no right answers" lowers pressure; polite end-of-form messages. |
| Authority | 4 | Exact dates, retention periods and withdrawal deadline. |
| Clarity | 4 | Closed questions with examples; select-one and select-all instructions stated. Fourteen questions is at the upper limit for a four-minute form. |
| Neutrality | 4 | Q8 mixes neutral categories (public code, documentation) with private ones, and Q9 asks about professional field without naming AI security. Q8 still hints at the topic; it comes after the hosting questions and is not used for selection. |
| Data minimisation | 5 | No name; region only; contact stored separately and deleted early for respondents not interviewed; no question invites disclosure of employer-policy breaches. |
