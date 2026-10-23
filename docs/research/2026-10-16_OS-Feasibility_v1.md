# AgentOS Bootable-OS Feasibility Research (Workflow B)

> Generated on 2026-10-16 by the AgentOS research workflows. Each brief was produced by a research agent and then
> attacked by an independent fact-checking agent; read the verification section before relying on any figure.
> Caveat: the shared web-search budget ran out early, so most evidence comes from direct reads of primary sources.
> Willingness-to-pay figures are mostly vendor-sponsored. This document is research input, not a statement of fact.


## Contents

- How solo or part-time maintainers build successful open-source infrastructure (case-derived rules for a solo, part-time AgentOS maintainer)
- What an AI OS can do only by owning the OS layer: agent-native kernel and system primitives
- Existing AI-focused Linux distributions and OS images: inventory, traction and gaps (as of 2026-10-16)
- Image-based Linux build tooling for a bootable AI OS maintained by one person
- Developing and testing a bootable AI OS on an Apple M1 Pro with 16 GB

## How solo or part-time maintainers build successful open-source infrastructure (case-derived rules for a solo, part-time AgentOS maintainer)
_Workflow label: `oss_solo_sustainability`_

### Executive summary
GitHub API data pulled on 16 October 2026 shows that the most successful "solo" AI and OS infrastructure projects still depend heavily on their founders, even at very large scale. The founder wrote 68.6% of Open WebUI's commits (155,856 stars), 73.7% of Coolify's (63,038 stars), 68.7% of Omarchy's (44,921 stars), 66.1% of Bazzite's and 100% of Pinokio's. Gaining contributors did not lift the founder's load. llama.cpp is the one exception: Gerganov wrote 17.4% of 11,691 commits, spread across roughly 2,100 contributors, because hardware vendors and Hugging Face staff contribute backends.

Every breakout case started as a narrow wedge built on someone else's heavy lifting. llama.cpp targeted one model on one laptop ("4-bit quantization on a MacBook"), and Simon Willison had it running one day after release. Ollama wrapped llama.cpp behind a one-line curl install. Open WebUI began as a front end for Ollama. Omarchy and Universal Blue configure existing distributions: Arch, and Fedora Atomic via bootc and GitHub Actions. Asahi Linux is the counter-case because it did hardware bring-up and kernel upstreaming. Its lead resigned in February 2025, citing burnout, user entitlement and harassment, with more than 1,000 downstream patches still outstanding.

Donations do not fund an OS. Asahi's Open Collective has raised about USD 156,000 in roughly 20 months for a seven-person team, and Tidelift reports that 60% of maintainers are unpaid. Large sums came only from four sources:
- venture capital (Ollama, USD 88M)
- angel investors (ggml)
- celebrity patronage (Omarchy, USD 23.65M pledged, self-reported)
- acquisition (ggml by Hugging Face in February 2026; Astral by OpenAI in March 2026)

Public grants are small or restrictive. NLnet first grants run EUR 5,000-50,000, and its programmes are in transition. The Sovereign Tech Fund requires a scope above EUR 50,000 and excludes prototypes.

Permissive licences dominate. Open WebUI's non-OSI branding clause and Redis's SSPL episode show that restrictive licences cost trust. Jan moved from AGPL to Apache 2.0 to make adoption easier.

The founder's premise has a structural tension. A solo, part-time founder wants an operating system that is both bootable and "revolutionary," while developing on an M1 Pro that cannot test CUDA fine-tuning paths. The evidence points to a bootc image layer rather than a new distribution, one safety demo that users can verify, Apache-2.0 plus a trademark policy, and keeping a day job.

### Recommendations
- Do not build a distribution. Make AgentOS a single Containerfile layered on a Fedora bootc or Universal Blue base, built and cosign-signed in GitHub Actions from the ublue-os image-template, with ISOs built by bootc-image-builder. Hold v0.1 to one repository, one image and under 1,000 lines of your own code.
- Write a NON-GOALS section before the first commit, as uv did at launch: no custom kernel, no new installer, no own inference engine (use llama.cpp, Ollama or vLLM), no desktop environment work, and a single CPU architecture for v0.1.
- Choose one demo users can verify and build only toward it, for example: boot from USB or switch with 'bootc switch', run an agent against a local model in a sandbox, then undo everything it changed with one rollback command. Record a video under 90 seconds. llama.cpp spread within a day because its claim was concrete and checkable.
- Match targets to your test hardware. Develop and test the aarch64 image in a VM on the M1 Pro, label x86_64/NVIDIA and fine-tuning paths as community-tested, and recruit at least two hardware testers before claiming GPU fine-tuning support.
- Launch on one day in sequence: Show HN, then r/LocalLLaMA, then the Universal Blue Discourse. Submit to the Ollama and Open WebUI integration lists in the same week. Before launch, publish a weekly devlog for 8 to 12 weeks to build the audience you currently lack (Omarchy's advantage is distribution, not code).
- Credit every upstream prominently in the README and boot splash (llama.cpp, Fedora, Universal Blue). Ollama's attribution backlash shows the cost of not doing so.
- License your own code under Apache-2.0 (with its patent grant) and accept contributions under a DCO. Protect the name with a TRADEMARK.md policy, not licence clauses. Avoid custom or non-OSI terms: Open WebUI's branding clause made it non-OSI, which would also bar Sovereign Tech Fund eligibility.
- Expect to write 70% or more of commits for years, as the founders of Open WebUI, Coolify, Omarchy and Bazzite do. Budget contributor work as review and triage, not as relief.
- Cap time strictly. At about 10 hours per week, allocate roughly 6 to building, 2 to docs and demos, and 2 to triage. Use issue forms that require reproduction steps, convert support questions to Discussions, and publish a policy requiring disclosure of AI-generated PRs and security reports (the curl lesson).
- Harden governance against xz-style takeover. Grant commit or release rights only after six or more months of verified contribution, require two-person review for build scripts and test fixtures, and use GitHub private security advisories.
- Keep your day job and treat donations as a bonus. Asahi raised about USD 156,000 in 20 months for seven people. Open GitHub Sponsors on day one, apply to NLnet's next call (3 Nov 2026 deadline, or the following odd-month deadline) for the sandbox component, and apply to the GitHub Secure Open Source Fund once you have traction.
- Decide your exit stance in writing now: an independent foundation, a hosted-service company, or openness to acquisition. The 2026 acquisitions of ggml and Astral show that labs buy projects that become infrastructure, and contributors will judge you on your stated intent.

### Pain points

**Founder remains the dominant committer regardless of popularity** — severity 5, breadth 5, willingness to pay 2, confidence high
- Segment: Solo founders of application-layer and distro-layer open-source projects
- Description: Contributor growth does not transfer the core workload. Even projects with 44k-156k stars show the founder authoring roughly two-thirds or more of all commits, so bus factor stays near one.
- Evidence: GitHub API, 16 Oct 2026: Open WebUI tjbck 13,393 of 19,511 commits (68.6%); Coolify andrasbacsai 13,153 of 17,846 (73.7%); Omarchy dhh 4,829 of 7,026 (68.7%); Bazzite KyleGospo 5,693 of 8,613 (66.1%); Pinokio cocktailpeanut 677 of 677 (100%). Counter-example: llama.cpp ggerganov 2,032 of 11,691 (17.4%).
- Quantified impact: Founder share 66-100% in 5 of 7 sampled successful projects; only library-layer llama.cpp diffused authorship.
- Current solutions: Hiring a first paid developer from sponsor or cloud revenue (Coolify, mid-2024); shared leadership after founder exit (Asahi seven-developer model); joining a larger company (ggml to Hugging Face).
- Why still unsolved: Opinionated products require one coherent vision; drive-by contributors fix edges, not architecture; review cost of outside PRs often exceeds writing the code.
- Sources: [GitHub REST API repository and contributor endpoints (queried for open-webui, coolify, omarchy, bazzite, pinokio, llama.cpp)](https://api.github.com/repos/open-webui/open-webui/contributors) (2026-10-16); [Coolify sponsorships page (team history)](https://coolify.io/sponsorships) (2026-10-16)

**Donations are small and decay after the launch peak** — severity 5, breadth 5, willingness to pay 3, confidence high
- Segment: Bootable OS and hardware-enablement projects; volunteer maintainers
- Description: Community donations peak at launch and then decline, and even well-loved OS projects raise less than one Western salary per year from crowdfunding.
- Evidence: Asahi Open Collective: USD 156,337.62 raised, estimated annual budget USD 117,653.87, for a seven-developer team (Oct 2026). Hector Martin: 'The all-time peak monthly donation volume was the very first month or two' (Feb 2025). Tidelift: 60% of maintainers are unpaid hobbyists, 12% earn most or all income from maintenance (Nov 2024).
- Quantified impact: About USD 118k per year spread across seven developers for a high-profile OS project.
- Current solutions: Fiscal hosting (Open Source Collective), GitHub Sponsors, sponsor logo tiers, corporate GPU or CI donations (LocalAI via Spectro Cloud).
- Why still unsolved: Users treat infrastructure as free; corporate users rarely pay for indirect dependencies; donation platforms offer no recurring value exchange.
- Sources: [Asahi Linux Open Collective](https://opencollective.com/asahilinux) (2026-10-15); [Resigning as Asahi Linux project lead (Hector Martin)](https://marcan.st/2025/02/resigning-as-asahi-linux-project-lead/) (2025-02-14); [The 2024 Tidelift Maintainer Impact Report](https://www.sonarsource.com/the-2024-tidelift-maintainer-impact-report.pdf) (2024-11)

**User entitlement, support load and issue backlog** — severity 4, breadth 5, willingness to pay 2, confidence high
- Segment: Desktop and OS distributions; local-AI runners
- Description: Popular projects accumulate large open-issue queues and feature demands that a part-time maintainer cannot service, which drives burnout.
- Evidence: Open issues (GitHub API, 16 Oct 2026): Ollama 4,296; Bazzite 1,206; Pinokio 475 with a single author. Hector Martin: 'No matter how much we did... people always wanted more.'
- Quantified impact: Hundreds to thousands of open issues per project; Asahi lead resignation attributed partly to demands.
- Current solutions: Discord-based community support (Bazzite, Coolify 20k+ members), aggressive issue closing (Open WebUI shows 170 open issues against 155k stars), Discussions forums.
- Why still unsolved: Hardware diversity in OS projects creates unbounded bug surface; users expect product-grade support from volunteers.
- Sources: [GitHub REST API repository endpoints](https://api.github.com/repos/ublue-os/bazzite) (2026-10-16); [Resigning as Asahi Linux project lead](https://marcan.st/2025/02/resigning-as-asahi-linux-project-lead/) (2025-02-14)

**Upstream friction for low-level OS work** — severity 5, breadth 3, willingness to pay 2, confidence high
- Segment: Hardware enablement and kernel-level OS projects
- Description: Projects that must change kernels or drivers carry large downstream patch stacks and depend on slow, sometimes hostile upstream review.
- Evidence: Asahi downstream tree held over 1,000 Apple Silicon patches not yet upstream; M3 and M4 support deferred until upstreaming and CI progress (Feb 2025). Martin: 'Upstreaming to Linux is a terrible experience in C too.' May 2026 addendum: he will forgo participation in the community entirely.
- Quantified impact: 1,000+ patch backlog; new hardware generations blocked.
- Current solutions: Shared maintainer teams, CI on device farms, prioritising upstreaming over features.
- Why still unsolved: Kernel governance is outside any single project's control; reviewer bandwidth upstream is scarce.
- Sources: [Passing the torch (Asahi Linux)](https://asahilinux.org/2025/02/passing-the-torch/) (2025-02); [Resigning as Asahi Linux project lead (with May 2026 update)](https://marcan.st/2025/02/resigning-as-asahi-linux-project-lead/) (2026-05)

**Commercial free-riding and rebranding prompt licence changes that erode trust** — severity 4, breadth 3, willingness to pay 3, confidence high
- Segment: Popular AI front-ends and infrastructure with enterprise demand
- Description: Companies repackage permissively licensed work without contributing, pushing maintainers toward restrictive licences that then trigger community backlash or forks.
- Evidence: Open WebUI added a branding clause in v0.6.6 (19 Apr 2025) because 'some groups were stripping out the branding, repackaging our work as their own, and monetizing it'; licence is not OSI-approved and post-v0.6.5 contributions require a CLA. Redis moved to SSPL in March 2024, which 'hurt our relationship with the Redis community' and led AWS and Google to maintain a fork; Redis added AGPLv3 on 1 May 2025.
- Quantified impact: Loss of OSI status; ineligibility for funds requiring OSI/FSF licences (Sovereign Tech Fund).
- Current solutions: Branding clauses, enterprise licences, SSPL/AGPL, trademark enforcement.
- Why still unsolved: Permissive licences maximise adoption but permit rebranding; restrictive licences reduce adoption and contributor goodwill.
- Sources: [Open WebUI License](https://docs.openwebui.com/license/) (2025-04-19); [Redis is now available under the AGPLv3 open source license](https://redis.io/blog/agplv3/) (2025-05-01)

**Social-engineering and AI-generated noise target exhausted maintainers** — severity 5, breadth 4, willingness to pay 3, confidence high
- Segment: All solo maintainers, especially security-relevant system software
- Description: Overloaded maintainers are vulnerable to pressure campaigns for co-maintainer access and are increasingly swamped by low-quality AI-generated reports.
- Evidence: xz: maintainer cited 'longterm mental health issues' and 'an unpaid hobby project' (8 Jun 2022) while likely sockpuppet accounts pressed for a co-maintainer who later gained control. curl: about 20% of 2025 security submissions AI-generated, about 5% valid, each report engaging three to four people for 30 minutes to 3 hours (14 Jul 2025).
- Quantified impact: curl security team of seven spends multi-person hours per invalid report; xz led to a backdoor in a core OS package.
- Current solutions: Bans and reputation systems, private security advisories, considering removal of bounties.
- Why still unsolved: New accounts are cheap; AI makes plausible-looking submissions nearly free to produce.
- Sources: [Timeline of the xz open source attack (Russ Cox)](https://research.swtch.com/xz-timeline) (2024-04); [Death by a thousand slops (Daniel Stenberg)](https://daniel.haxx.se/blog/2025/07/14/death-by-a-thousand-slops/) (2025-07-14)

**Public grants are small, slow, or exclude early-stage and user-facing projects** — severity 3, breadth 4, willingness to pay 4, confidence high
- Segment: Pre-traction open-source founders in or outside Europe
- Description: Non-dilutive public funding exists but its eligibility, size and timing rarely fit a new, solo, user-facing OS project.
- Evidence: Sovereign Tech Fund: minimum estimated cost above EUR 50,000, about six months to contract, 'We do not finance the development of prototypes', user-facing applications not sought. NLnet NGI Zero Commons Fund: first grants EUR 5,000-50,000; final call closed 1 Jun 2026; a 12 Jun 2026 notice says open calls temporarily paused during the transition from NGI to an Open Internet Stack; new calls listed with a 3 Nov 2026 deadline. GitHub Secure Open Source Fund: USD 10,000 per project in tranches.
- Quantified impact: Typical first grant EUR 5k-50k or USD 10k; six-month lead times.
- Current solutions: NLnet calls, GitHub Secure Open Source Fund, a16z Open Source AI Grants (2023 cohort included text-generation-webui and vLLM).
- Why still unsolved: Funders prefer reusable base technology with proven usage; programmes are cyclical and politically dependent.
- Sources: [Sovereign Tech Fund programme page](https://www.sovereign.tech/programs/fund) (2026-10-16); [NGI Zero Commons Fund](https://nlnet.nl/commonsfund/) (2026-10-15); [NLnet news](https://nlnet.nl/news/) (2026-10-15); [GitHub Secure Open Source Fund](https://github.com/open-source/github-secure-open-source-fund) (2026-10-16)

**Trust backlash when a wrapper commercialises without crediting upstream** — severity 3, breadth 3, willingness to pay 2, confidence medium
- Segment: Wrapper and orchestration projects built on llama.cpp or similar engines
- Description: Thin layers over other projects grow fast but face criticism over attribution, lock-in and cloud pivots once they raise venture money.
- Evidence: Practitioner critique (15 Apr 2026) links Ollama issue #3697 (llama.cpp omitted from README for over a year), a July 2025 GUI shipped without a licence and open-sourced in November 2025 (PR #12933), and cloud models possibly forwarding prompts to third parties (issue #14279). Ollama has raised USD 88M (9 Jul 2026) and introduced per-token pricing (31 Aug 2026).
- Quantified impact: Reputational; motivates competing projects and forks.
- Current solutions: Adding upstream credit, open-sourcing components after criticism.
- Why still unsolved: VC growth incentives favour proprietary differentiation over upstream credit.
- Sources: [Friends Don't Let Friends Use Ollama (practitioner opinion)](https://sleepingrobots.com/dreams/stop-using-ollama/) (2026-04-15); [Ollama: all aboard open models](https://ollama.com/blog/all-aboard-open-models) (2026-07-09)

**Discovery is dominated by founders with existing audiences** — severity 4, breadth 5, willingness to pay 2, confidence medium
- Segment: Unknown solo founders launching OS or desktop projects
- Description: Distribution reach, not technical novelty, explains much of the gap between celebrity-led and unknown-founder projects.
- Evidence: Omarchy (DHH) self-reports 1,280,643 ISO downloads in year one and USD 23,650,000 pledged to the Omacom Foundation, reaching 44,921 stars in about 16 months; Pinokio, a technically ambitious solo project, has 8,260 stars after three years. llama.cpp spread via Simon Willison's write-up one day after release.
- Quantified impact: Order-of-magnitude difference in adoption between audience-backed and unknown launches.
- Current solutions: Building in public, demos, Hacker News and Reddit launches, newsletter coverage, integrations with popular tools.
- Why still unsolved: Attention is concentrated; demo virality is unpredictable.
- Sources: [Omarchy homepage](https://omarchy.org/) (2026-10-16); [Large language models are having their Stable Diffusion moment (Simon Willison)](https://simonwillison.net/2023/Mar/11/llama/) (2023-03-11); [GitHub REST API (pinokio, omarchy)](https://api.github.com/repos/pinokiocomputer/pinokio) (2026-10-16)

### Key statistics

| Claim | Value | Source type | Date | Source |
|---|---|---|---|---|
| llama.cpp GitHub stars and approximate contributors | 131,342 stars; about 2,126 contributors (anonymous included); MIT; created 2023-03-10 | primary_research | 2026-10-16 | https://api.github.com/repos/ggml-org/llama.cpp |
| Share of llama.cpp commits authored by Georgi Gerganov | 2,032 of 11,691 commits (17.4%) | primary_research | 2026-10-16 | https://api.github.com/repos/ggml-org/llama.cpp/contributors |
| Share of Open WebUI commits authored by founder | 13,393 of 19,511 commits (68.6%); 155,856 stars | primary_research | 2026-10-16 | https://api.github.com/repos/open-webui/open-webui/contributors |
| Share of Coolify commits authored by founder | 13,153 of 17,846 commits (73.7%); 63,038 stars | primary_research | 2026-10-16 | https://api.github.com/repos/coollabsio/coolify/contributors |
| Share of Omarchy commits authored by DHH | 4,829 of 7,026 commits (68.7%); 44,921 stars; created 2025-06-01; repository moved to omacom/omarchy | primary_research | 2026-10-16 | https://api.github.com/repos/omacom/omarchy/contributors |
| Share of Bazzite commits authored by top maintainer; open issues | 5,693 of 8,613 commits (66.1%); 1,206 open issues and PRs | primary_research | 2026-10-16 | https://api.github.com/repos/ublue-os/bazzite |
| Pinokio authorship concentration | 677 of 677 commits by one author; 8,260 stars; 475 open issues | primary_research | 2026-10-16 | https://api.github.com/repos/pinokiocomputer/pinokio |
| Omarchy self-reported ISO downloads and foundation pledges | 1,280,643 ISO downloads in year one; USD 23,650,000 pledged to the Omacom Foundation; USD 100,000 bug bounty fund | vendor | 2026-10-16 | https://omarchy.org/ |
| Ollama total funding raised and self-reported reach | USD 88M raised; 8.9 million developers; used by 85% of Fortune 500 | vendor | 2026-07-09 | https://ollama.com/blog/all-aboard-open-models |
| Ollama community integrations listed in README | About 154 entries across 12 categories | primary_research | 2026-10-16 | https://github.com/ollama/ollama |
| Asahi Linux Open Collective totals | USD 156,337.62 raised; USD 117,653.87 estimated annual budget; 967 contributors | primary_research | 2026-10-15 | https://opencollective.com/asahilinux |
| Asahi downstream patch backlog at leadership change | Over 1,000 Apple Silicon patches not upstream; seven developers share leadership | official | 2025-02 | https://asahilinux.org/2025/02/passing-the-torch/ |
| Share of maintainers unpaid | 60% unpaid hobbyists; 12% earn most or all income from maintenance | survey | 2024-11 | https://www.sonarsource.com/the-2024-tidelift-maintainer-impact-report.pdf |
| curl AI-generated security report share and validity rate | About 20% AI-generated; about 5% valid in 2025; 7-person security team; 30 minutes to 3 hours per report for 3-4 people | practitioner | 2025-07-14 | https://daniel.haxx.se/blog/2025/07/14/death-by-a-thousand-slops/ |
| Open WebUI branding licence change | Effective v0.6.6 (19 Apr 2025); branding changes allowed only at 50 or fewer users per 30 days, with permission, or enterprise licence; not OSI-approved; CLA required after v0.6.5 | official | 2025-04-19 | https://docs.openwebui.com/license/ |
| Jan licence change | AGPL to Apache 2.0 (May 2025), with a non-binding attribution request; copyright Menlo Research; 72 contributors at the time raised consent questions | primary_research | 2025-05 | https://github.com/janhq/jan/blob/main/LICENSE |
| Redis licence trajectory | SSPL/RSAL in March 2024; AGPLv3 option added 1 May 2025 with Redis 8 | vendor | 2025-05-01 | https://redis.io/blog/agplv3/ |
| ggml.ai joins Hugging Face | Announced 20 Feb 2026; team stays full-time on ggml/llama.cpp with stated technical autonomy | official | 2026-02-20 | https://huggingface.co/blog/ggml-joins-hf |
| Astral (uv, Ruff, ty) agrees to join OpenAI | Announced 19 Mar 2026; tools reached 'hundreds of millions of monthly downloads' | official | 2026-03-19 | https://astral.sh/blog/openai |
| uv launch scope and speed claims | Launched 15 Feb 2024 as drop-in for pip/pip-tools; 8-10x faster cold, 80-115x warm cache | vendor | 2024-02-15 | https://astral.sh/blog/uv |
| llama.cpp diffusion speed after release | Released 10 Mar 2023; running on a Raspberry Pi by 12 Mar and a Pixel 6 by 13 Mar | practitioner | 2023-03-11 | https://simonwillison.net/2023/Mar/11/llama/ |
| Sovereign Tech Fund eligibility thresholds | Minimum project cost above EUR 50,000; about 6 months to contract; no prototypes; OSI/FSF licence required | official | 2026-10-16 | https://www.sovereign.tech/programs/fund |
| NLnet NGI Zero Commons Fund grant size and status | First grant EUR 5,000-50,000; EUR 21.6M planned through 2027; final call closed 1 Jun 2026; next listed deadline 3 Nov 2026 for new calls | official | 2026-10-15 | https://nlnet.nl/commonsfund/ |
| GitHub Secure Open Source Fund award | USD 10,000 per project (6k/2k/2k tranches) plus three-week security programme | official | 2026-10-16 | https://github.com/open-source/github-secure-open-source-fund |
| Ollama founding and accelerator | Founded 2021; Y Combinator Winter 2021; earlier founders built Kitematic (acquired by Docker) | other | 2026 | https://aiwiki.ai/wiki/ollama |

### Players

- **llama.cpp / ggml (Georgi Gerganov, now Hugging Face)** (Inference engine library): Dependency-free C/C++ inference with quantization across CPU, Metal, CUDA, Vulkan; GGUF format. Traction: 131,342 stars, about 2,100 contributors (Oct 2026); team joined Hugging Face 20 Feb 2026. Gap: Open-core business model never defined; independence now rests on Hugging Face goodwill without formal governance; whisper.cpp future questioned by community
- **Ollama** (Local model runner, VC-backed company): One-line install, 'ollama run' UX, model library, local API, cloud models with per-token pricing. Traction: 183,386 stars; USD 88M raised; about 154 community integrations. Gap: Attribution and lock-in criticism; cloud pivot raises privacy questions; not an OS and not a sandbox
- **Open WebUI (Timothy Jaeryang Baek)** (Self-hosted AI front end, founder-led company): Docker/pip one-liner chat UI for Ollama and OpenAI-compatible APIs; enterprise licence. Traction: 155,856 stars; 889 contributors; founder 68.6% of commits. Gap: Non-OSI branding licence and CLA; heavy founder dependence
- **Pinokio (cocktailpeanut)** (One-click AI app launcher): JSON scripts that install and run open-source AI apps in a local directory. Traction: 8,260 stars; 677 commits all by one author. Gap: Isolation is convention-based; README admits scripts can deviate; single-maintainer bus factor
- **Universal Blue (Bazzite, Bluefin, Aurora)** (Image-based Linux built with cloud-native tooling): bootc/OCI images on Fedora Atomic built by GitHub Actions, cosign-signed, 90-day rollback. Traction: Bazzite 9,184 stars, 216 contributors; Bluefin 2,629 stars. Gap: Concentrated maintainers (66.1% top committer on Bazzite); 1,206 open items; no AI or agent focus
- **Omarchy (DHH, Omacom Foundation)** (Opinionated Arch-based desktop OS): ISO install, Hyprland desktop, first-boot agent setup, agent skills for apps and themes. Traction: 44,921 stars in about 16 months; self-reported 1.28M ISO downloads; USD 23.65M pledged. Gap: Agents run with broad user privileges; no evident sandboxing or rollback story; celebrity-dependent distribution
- **Asahi Linux** (Hardware-enablement Linux distribution for Apple Silicon): Fedora Asahi Remix with reverse-engineered GPU, audio and display drivers. Traction: Open Collective USD 156k raised; 967 contributors to collective. Gap: Founder burnout and exit; 1,000+ downstream patches; M3/M4 support deferred
- **uv / Astral (now OpenAI)** (Developer tooling company): Rust-based Python package and project manager. Traction: 90,802 stars; hundreds of millions of monthly downloads across Astral tools. Gap: Future governance under OpenAI not specified
- **LocalAI (Ettore Di Giacinto)** (Solo-led OpenAI-compatible local inference server): Multi-backend local API for text, audio, image, P2P inference. Traction: 49,568 stars; founder 40.8% of commits. Gap: Maintainer funded by consulting; sponsorship mainly covers hardware
- **Jan (Menlo Research)** (Desktop local-AI app, company-backed): Offline chat app with local and remote models. Traction: 44,963 stars; 185 contributors. Gap: Relicensing from AGPL raised contributor-consent questions
- **Coolify (Andras Bacsai)** (Self-hosted PaaS, self-funded company): Free self-hosted core plus paid Cloud. Traction: 63,038 stars; 765 contributors; reported about USD 15,700 monthly Cloud gross (aggregator, Feb 2026). Gap: Founder 73.7% of commits; revenue figure unverified

### Trends
- Acquisition by AI labs and platforms is now a common way for successful AI-infrastructure open-source projects to end up: ggml.ai joined Hugging Face (20 Feb 2026) and Astral agreed to join OpenAI (19 Mar 2026).
- 'OS for agents' positioning has gone mainstream. Omarchy markets itself as 'The malleable OS for the age of agents', with default-agent setup on first boot (Oct 2026). Ollama and Fly.io promote running agents safely (13 Oct 2026).
- Image-based Linux built with bootc, OCI and GitHub Actions has cut the maintenance cost of derived distributions. Universal Blue images are cosign-signed with 90-day rollback, and ublue-os/main narrowed to fewer base images in September 2025.
- Local-first tools are turning to cloud revenue: Ollama cloud models (preview September 2025), per-token pricing (31 Aug 2026), and Coolify Cloud.
- Licences cycle between restriction and permissiveness. Redis went to SSPL in 2024 and added AGPL in 2025. Jan moved from AGPL to Apache 2.0 in 2025. Open WebUI added a branding clause in 2025.
- AI-generated reports and PRs are adding to maintainer load: about 20% of curl's 2025 security submissions were AI-generated.
- European public funding for open source is in transition. The NGI Zero Commons final call closed 1 Jun 2026, open calls were paused on 12 Jun 2026, and new calls carry a 3 Nov 2026 deadline.
- Corporate-patron foundations are emerging for desktop operating systems: the Omacom Foundation lists Meta Superintelligence Labs, DigitalOcean, Alibaba Cloud and others as patrons (Sep-Oct 2026).
- Shared or company governance increasingly follows a founder's exit or acqui-hire: Asahi moved to seven-developer shared leadership, and ggml moved under Hugging Face with stated autonomy.

### Opportunities

- **Make the AgentOS wedge an 'agent action that is always reversible': agents run in confined sandboxes on an atomic bootc base, and any system change an agent makes can be rolled back.** Why now: Agent-first desktops such as Omarchy give agents broad control, and safety demand is visible (Ollama with Fly.io, October 2026). Atomic images already provide 90-day rollback. Evidence: Pinokio's README admits risk when a script 'intentionally tries to deviate'. Bazzite offers rollback to any build in the last 90 days. Omarchy's agent features show no sandboxing claims on its homepage. Risk: Universal Blue, Omarchy or Fedora could add similar features quickly. Sandboxing GPU workloads is technically hard and may be impossible to test on an M1 Pro.
- **Ship AgentOS as a derived bootc image (FROM a Universal Blue or Fedora bootc base) instead of a new distribution, distributed via 'bootc switch' and ISOs built by bootc-image-builder.** Why now: The ublue-os image-template packages the build, signing and ISO pipeline in GitHub Actions. Evidence: Universal Blue states it is a way to consume existing distributions, not a new one, and its images are built entirely by GitHub Actions with cosign signing. Risk: Depends on upstream base-image decisions (for example, ublue-os/main reduced its images in September 2025). The product may look like 'just another spin' unless the agent layer is distinctive.
- **Become an integration target for existing ecosystems (the Ollama integrations list, Open WebUI, llama.cpp packaging, MCP clients).** Why now: Hugging Face's stated post-acquisition goal is to make llama.cpp 'ubiquitous and readily available everywhere'. Ollama lists about 154 integrations. Evidence: Open WebUI grew as a front end for Ollama. Ollama grew by wrapping llama.cpp. Risk: Commoditisation, and the attribution backlash that Ollama experienced if upstreams are not credited.
- **Seek non-dilutive funding for a reusable security component (an agent sandbox or capability policy) rather than for 'the OS'.** Why now: New NLnet calls carry a 3 Nov 2026 deadline. The GitHub Secure Open Source Fund accepts applications on a rolling basis. Evidence: NLnet first grants are EUR 5,000-50,000. The GitHub fund gives USD 10,000. The Sovereign Tech Fund excludes prototypes and user-facing applications but funds base technologies. Risk: Programme transitions and pauses; six-month lead times; eligibility may require residence or legal-entity conditions that remain unverified.
- **Later, offer a paid hosted service (remote GPU fine-tuning or fleet management) while keeping the OS fully open.** Why now: Precedents show hosted revenue funding the free core: Coolify Cloud, and Ollama's cloud models with per-token pricing. Evidence: Coolify funded its first hire in mid-2024 from sponsors and paid users. Ollama has raised USD 88M on a local-plus-cloud model. Risk: Trust erosion if the cloud is pushed by default. Running GPU infrastructure is capital-intensive for a solo founder.

### Contrarian signals
- Restrictive licensing did not visibly stall Open WebUI. It reached 155,856 stars and 889 contributors by October 2026 despite its non-OSI branding clause, so the 'restrict and die' narrative is not universal.
- Contributor growth did not reduce founder load in most cases. Founders still author 66-100% of commits in five of the seven sampled successful projects, so expecting the community to share the work is unsupported.
- Much of the success is distribution, not technology. Omarchy's self-reported 1.28M downloads and USD 23.65M in pledges rest on DHH's existing audience, and an unknown solo founder cannot replicate that model.
- The most celebrated 'solo' successes ended in venture funding or acquisition (Ollama USD 88M; ggml to Hugging Face; Astral to OpenAI), which conflicts with a purely independent, community-owned AgentOS vision.
- Hardware-level OS ambition is the documented burnout path. Asahi's founder resigned in February 2025 and in May 2026 withdrew from the community entirely, which supports an image-layer approach over bespoke hardware enablement.
- The 'AI-native OS' framing is already taken by a well-funded player. Omarchy calls itself 'The malleable OS for the age of agents', so 'revolutionary' must mean a specific capability, not positioning.
- Grant funding is shrinking or in transition just as the founder would need it: NGI Zero Commons' final call closed in June 2026, open calls were paused, and the Sovereign Tech Fund excludes prototypes and user-facing applications.
- Thin-wrapper strategies invite backlash: Ollama faced attribution, lock-in and privacy criticism, which suggests that wrapping upstream engines without visible value-add carries reputational risk.
- The founder's premise contains an internal inconsistency. Fast, safe local fine-tuning on CUDA hardware cannot be developed or verified on an M1 Pro with 16 GB, and Asahi does not yet support M3/M4, so bare-metal Apple Silicon is a narrowing target.

### Open questions
- Methodology limit: the shared web-search budget ran out after 13 distinct searches, below the 20-search target. The remaining evidence comes from about 25 primary sources read directly and live GitHub API queries. Remaining gaps include the full Tidelift 2024 survey data, Sovereign Tech Agency Fellowship terms, Ollama's 2023 Show HN date, and the Mozilla and a16z cohorts in 2025-2026.
- Clarifying question for the founder: what single capability would make AgentOS 'revolutionary' and that Omarchy, Universal Blue or Ollama could not add within six months? Without a precise answer, the bootable-OS requirement adds cost without differentiation.
- Clarifying question for the founder: should fine-tuning run locally on the user's GPU, remotely, or both? The M1 Pro cannot test CUDA paths, and the answer determines whether v0.1 targets aarch64 VMs, x86_64 NVIDIA, or Apple Silicon via Asahi.
- Clarifying question for the founder: is the long-term intent an independent foundation, a hosted-service company, or openness to acquisition? This determines whether to use a DCO or a CLA, and how to word governance commitments.
- How well does GPU or accelerated inference work inside Linux VMs on Apple Silicon (for example, virtio-gpu or Vulkan paths)? This is unverified here and governs whether the founder can dogfood local models in AgentOS on the M1 Pro.
- Which NLnet call (Restack or CodeSupply) fits an agent-sandbox component, and what are its eligibility rules for applicants outside the EU? This is not verified.
- Are Omarchy's 1.28M downloads and USD 23.65M in pledges independently auditable, and how much of the pledged amount has been received?
- Will Hugging Face and OpenAI keep ggml and uv governance autonomous in practice? Neither announcement describes a binding mechanism.

### Adversarial verification

- **PARTIALLY_CONFIRMED** — Founder commit concentration: Open WebUI 13,393 of 19,511 (68.6%), Coolify 13,153 of 17,846 (73.7%), Omarchy 4,829 of 7,026 (68.7%), Bazzite 5,693 of 8,613 (66.1%), Pinokio 677 of 677, llama.cpp 2,032 of 11,691 (17.4%). The brief infers that gaining contributors did not lift the founder's load, and recommends expecting to write 70% or more of commits.
  - Corrected: The counts reproduce exactly (llama.cpp now shows 11,693 commits). Commit counts mostly measure how often each person commits, not how much work they do. Open WebUI's founder made 100 commits between 12 and 14 Oct 2026, and DHH made 100 between 11 and 14 Oct. The data show concentration. They do not show that contributors failed to reduce load, and Coolify's paid hire, with 1,858 commits, did take on real work.
  - Evidence: GitHub API queries on 16 Oct 2026 (authenticated gh and unauthenticated curl) reproduced every star count and contributor count. Merge commits do not inflate the founder shares: they were 2 of tjbck's last 100 commits, 3 of andrasbacsai's, 0 of dhh's and 0 of ggerganov's. Founders still wrote most recent commits: 92 of Coolify's last 100 and 81 of Omarchy's last 100. Omarchy has 5,458 open items (2,250 issues and 3,208 PRs), which supports the review-burden point but was not cited in the brief.
  - Sources: [GitHub REST API contributors and commits endpoints (open-webui, coolify, omarchy, bazzite, pinokio, llama.cpp)](https://api.github.com/repos/open-webui/open-webui/contributors) (2026-10-16)
- **CONFIRMED** — Ollama has raised USD 88M (blog, 9 Jul 2026) and introduced per-token pricing (31 Aug 2026); it claims 8.9 million developers and use by 85% of the Fortune 500.
  - Corrected: Ollama announced USD 88M from Benchmark, Theory Ventures and 8VC on 9 Jul 2026, without naming the round. On 31 Aug 2026, Ollama moved its cloud Pro, Max and Team plans to per-token pricing, and existing subscribers keep their current plans until they change tier. The announcement does not cover local use. The 8.9M developer and 85% Fortune 500 figures are self-reported and undefined.
  - Evidence: The July post quotes the USD 88M raise and names the investors. The 31 Aug post says that Pro, Max and Team plans 'now use transparent per-token pricing.' The July post also contrasts local models with 'runaway per-token bills', so the pricing change applies to the cloud only.
  - Sources: [Ollama: all aboard open models](https://ollama.com/blog/all-aboard-open-models) (2026-07-09); [Ollama's transparent pricing](https://ollama.com/blog/transparent-pricing) (2026-08-31)
- **PARTIALLY_CONFIRMED** — Exit via acquisition: ggml joined Hugging Face (20 Feb 2026) and Astral was acquired by OpenAI (March 2026); the brief's trend and recommendation sections rely on these as completed acquisitions.
  - Corrected: Hugging Face announced on 20 Feb 2026 that ggml.ai was joining it. The team keeps '100%' of its time on llama.cpp with 'full autonomy'. No financial terms, licence or governance mechanism were stated. On 19 Mar 2026, Astral said it had 'entered into an agreement to join' OpenAI. The page does not confirm that the deal has closed, and treating it as a completed acquisition is not supported by the source.
  - Evidence: The Hugging Face post confirms the date and the autonomy wording, and promises to make llama.cpp 'ubiquitous and readily available everywhere'. It does not mention whisper.cpp. The Astral post confirms the date and the claim of 'hundreds of millions of downloads per month'. It describes joining Codex as a future step and gives no closing conditions or licence statement.
  - Sources: [GGML and llama.cpp join HF](https://huggingface.co/blog/ggml-joins-hf) (2026-02-20); [Astral to join OpenAI](https://astral.sh/blog/openai) (2026-03-19)
- **DISPUTED** — Funding recommendation: apply to NLnet's next call (3 Nov 2026 deadline) for an agent-sandbox component; NGI Zero Commons final call closed 1 Jun 2026, calls paused 12 Jun 2026, first grants EUR 5,000-50,000; Sovereign Tech Fund requires cost above EUR 50,000, about 6 months, no prototypes, OSI/FSF licence.
  - Corrected: The individual facts are correct, but the recommendation does not fit. NLnet's Restack call (deadline 3 Nov 2026, EUR 5k-50k) says 'AI-related projects are not within scope' unless they already have more than 1 million active human users. It also requires a strong European dimension: normally, team members in Horizon Europe-eligible countries receive at least 60% of the funds, and other proposals are considered only as exceptions. CodeSupply funds software supply-chain metadata only. An agent-sandbox application from a founder outside Europe is therefore unlikely to qualify unless the project is reframed as non-AI security infrastructure and has European participation. The Sovereign Tech Fund terms are confirmed.
  - Evidence: The NLnet Commons Fund page states that the 'thirteenth and final call... closed on June 1st 2026' and that grants run EUR 5,000-50,000 from a EUR 21.6M budget. The news page lists the 12 Jun 2026 pause notice, the 3 Aug 2026 move to odd-month deadlines and the 3 Sep 2026 call for 3 Nov. The Restack scope and eligibility pages contain the AI exclusion and the European-dimension rule. The Sovereign Tech Fund page states that costs 'must exceed EUR 50,000', the process takes about six months, it does not fund prototypes, it is not seeking user-facing applications, and code must use an OSI or FSF licence.
  - Sources: [NGI Zero Commons Fund](https://nlnet.nl/commonsfund/) (2026-10-16); [NLnet news](https://nlnet.nl/news/) (2026-10-15); [NLnet Restack call](https://nlnet.nl/restack/) (2026-10-16); [NLnet Restack eligibility](https://nlnet.nl/restack/eligibility/) (2026-10-16); [NLnet CodeSupply call](https://nlnet.nl/codesupply/) (2026-10-16); [Sovereign Tech Fund](https://www.sovereign.tech/programs/fund) (2026-10-16)
- **REFUTED** — Recommendation: build AgentOS as a Containerfile FROM a Fedora bootc or Universal Blue base using the ublue-os image-template, and develop and test the aarch64 image in a VM on the M1 Pro.
  - Corrected: The Universal Blue route and the aarch64 plan conflict. As of 16 Oct 2026, the Universal Blue images (base-main, bazzite, bluefin and aurora) are amd64 only, and the image-template's default base is bazzite:stable (amd64). Only Fedora's own bootc image (quay.io/fedora/fedora-bootc:44) is multi-architecture and includes arm64. An aarch64 build developed on an M1 Pro must therefore start FROM fedora-bootc, not a Universal Blue image, or depend on slow amd64 emulation. The template's disk-build workflow does support arm64 runners.
  - Evidence: Image configs fetched from GHCR give architecture 'amd64' for ublue-os/bazzite:stable, bluefin:stable and base-main:latest. The bluefin:lts index lists only amd64. The Quay manifest list for fedora/fedora-bootc:44 lists arm64, ppc64le, s390x and amd64. The image-template Containerfile starts 'FROM ghcr.io/ublue-os/bazzite:stable', and build-disk.yml offers amd64 and arm64 platforms.
  - Sources: [ublue-os/image-template repository (Containerfile, build-disk.yml)](https://github.com/ublue-os/image-template) (2026-10-16); [GHCR manifest for ublue-os/bazzite:stable](https://ghcr.io/v2/ublue-os/bazzite/manifests/stable) (2026-10-16); [Quay manifest for fedora/fedora-bootc:44](https://quay.io/v2/fedora/fedora-bootc/manifests/44) (2026-10-16)
- **REFUTED** — Asahi Linux does not yet support M3/M4 (M3/M4 support deferred), so bare-metal Apple Silicon is a narrowing target and hardware-level OS work is the documented burnout path.
  - Corrected: Asahi Linux announced official installer support for M3, M3 Pro and M3 Max MacBooks and iMacs on 6 Sep 2026, and on 15 Sep 2026 stopped requiring expert mode. GPU acceleration is not yet usable on M3, the M3 Ultra is unsupported, and sleep and HDMI do not work. M4 and M5 remain unsupported. Since the leadership change the project has also upstreamed its Mesa drivers (Aug 2025) and its SMC drivers (Oct-Dec 2025). The burnout of the founder is accurate, but the claim that the project stalled is outdated.
  - Evidence: The Asahi blog lists 'M2: Episode 1 (or, Asahi Linux on M3)', which states that Asahi Linux 'now officially supports Macs with an M3 series SoC'. The same blog has progress reports for Linux 6.14 through 7.2. The resignation post's May 2026 update confirms that Martin has left the community entirely.
  - Sources: [M2: Episode 1 (or, Asahi Linux on M3)](https://asahilinux.org/2026/09/m2-episode-1/) (2026-09-06); [Asahi Linux blog index](https://asahilinux.org/blog/) (2026-10-16); [Resigning as Asahi Linux project lead (May 2026 update)](https://marcan.st/2025/02/resigning-as-asahi-linux-project-lead/) (2026-05)
- **PARTIALLY_CONFIRMED** — Asahi's Open Collective raised about USD 156,000 in roughly 20 months for a seven-person team, showing donations do not fund an OS.
  - Corrected: The page reports USD 156,337.62 raised, USD 117,653.87 as the estimated annual budget, 967 contributors and seven listed team members. Only USD 43,419.98 has been disbursed, leaving a balance of USD 112,917.64, so the collective is not paying salaries to seven people. The collective began in February 2025, when the founder redirected his personal Patreon and GitHub Sponsors income to it. Its totals therefore exclude years of earlier personal-channel funding and cannot by themselves show that donations decay.
  - Evidence: The Open Collective page gives the totals, the disbursed amount and the balance, and its earliest contributions date from February 2025. The resignation post says the author is pausing his Patreon and asks supporters to redirect funding to the OpenCollective. It also states that donation volume peaked in the first month or two and then kept slowly decreasing.
  - Sources: [Asahi Linux Open Collective](https://opencollective.com/asahilinux) (2026-10-16); [Resigning as Asahi Linux project lead](https://marcan.st/2025/02/resigning-as-asahi-linux-project-lead/) (2025-02-14)
- **CONFIRMED** — Open WebUI branding clause effective v0.6.6 (19 Apr 2025), 50-user threshold, not OSI-approved, CLA after v0.6.5; Jan moved AGPL to Apache 2.0 in May 2025.
  - Corrected: Both licence facts are confirmed. Open WebUI's page states that 'v0.6.6+ is not an OSI-approved' licence and that code up to v0.6.5 remains BSD-3-Clause. Jan's LICENSE was changed to Apache on 20 May 2025 (PR #5042), but GitHub still reports Jan's licence as NOASSERTION, so its LICENSE file contains more than the standard Apache text.
  - Evidence: The Open WebUI licence page gives the v0.6.6 effective date and the 50-user, 30-day threshold, and says the CLA applies from v0.6.6. Its TL;DR and body word the threshold slightly differently. The janhq/jan LICENSE commit history shows 'Jan's code is now under the Apache license (#5042)' dated 2025-05-20.
  - Sources: [Open WebUI License](https://docs.openwebui.com/license/) (2025-04-19); [janhq/jan LICENSE history (GitHub API)](https://github.com/janhq/jan/blob/main/LICENSE) (2025-05-20)
- **CONFIRMED** — Omarchy self-reports 1,280,643 ISO downloads in year one, USD 23.65M pledged to the Omacom Foundation, a USD 100,000 bug bounty, and markets itself as 'The malleable OS for the age of agents' without sandboxing claims.
  - Corrected: The homepage reports these figures. All are self-reported, pledges are not receipts, and the page gives no breakdown that adds up to USD 23.65M. OpenAI is also listed as a corporate patron, in addition to Meta Superintelligence Labs, DigitalOcean and Alibaba Cloud. The page describes no sandboxing, permissions or rollback.
  - Evidence: omarchy.org quotes the download count, pledge total and bounty fund, the bounties are administered through HackerOne, and the page lists its patron tiers.
  - Sources: [Omarchy homepage](https://omarchy.org/) (2026-10-16)

**Missed counter-evidence**
- NLnet's Restack call, the main 3 Nov 2026 option, excludes AI-related projects unless they already have more than 1 million active users, and normally requires 60% of funds to go to team members in Horizon Europe-eligible countries. The brief left eligibility outside the EU as an open question, but this rule directly undermines its NLnet recommendation.
- Universal Blue images (base-main, bazzite, bluefin, aurora) are amd64 only as of 16 Oct 2026. Only quay.io/fedora/fedora-bootc:44 provides arm64, so 'FROM Universal Blue' and 'test aarch64 on the M1 Pro' cannot both apply.
- Asahi Linux officially supports M3, M3 Pro and M3 Max installs since September 2026 (without GPU acceleration), and has upstreamed its Mesa and SMC drivers since the founder left. The brief's statement that Asahi does not support M3 is outdated. The M1 Pro remains Asahi's best-supported generation, so a bare-metal aarch64 test target is available alongside VMs.
- Ollama added an MLX engine for Apple Silicon (preview 30 Mar 2026, performance update 11 Jun 2026), and Ollama 0.30 added GGUF support through llama.cpp (5 Jun 2026). This matters to an M1 Pro developer: native macOS inference is now strong, while Linux VMs on the same machine have no Metal access. The brief's open question about GPU in VMs should be weighed against this.
- curl's bug-bounty page now states that there is no bug bounty, so the brief's 'considering removal of bounties' is out of date. In October 2026, curl also credited the AI-assisted Aisle audit with 35 published CVEs and acknowledged an Anthropic-found use-after-free that became remote code execution. A policy that requires disclosure of AI use is more defensible than an anti-AI policy.
- Omarchy's own backlog is 5,458 open items (2,250 issues and 3,208 PRs), which is stronger review-burden evidence than the Ollama, Bazzite and Pinokio figures the brief cited.
- Asahi's Open Collective has disbursed only USD 43,420 of USD 156,338, so unspent funds, not only low donations, limit how far the money goes. The brief presents the total as if it funds seven developers.
- The Astral announcement describes an agreement, not a completed deal, and Hugging Face describes ggml as 'joining'. Neither discloses terms, so calling both 'acquisitions' overstates the record.
- Minor: Simon Willison ran llama.cpp the night it was released (10 Mar 2023), not one day later.

**Reliability note:** Summary: the brief's numbers are highly accurate, but some recommendations rest on facts it never checked or that are now out of date. Every GitHub star, commit and contributor figure, and every figure from Ollama, Omarchy, Open WebUI, the Sovereign Tech Fund and NLnet, matched the primary sources on 16 Oct 2026. The weaker parts are interpretive or practical. Commit counts are treated as workload. Astral and ggml are presented as completed acquisitions. Asahi's M3 status is out of date. The NLnet recommendation ignores Restack's exclusion of AI projects and its European-dimension rule. The bootc recommendation pairs amd64-only Universal Blue bases with aarch64 development on an M1 Pro. Method: web search was unavailable because the shared budget was exhausted, so all checks used direct fetches of primary pages and live GitHub, GHCR and Quay API queries. Some findings, such as whether the Astral deal closed and how Restack eligibility applies outside Europe, remain unconfirmed and need a follow-up search. Two questions for the orchestrator: does the founder have European collaborators who could meet the NLnet rule, and is x86_64 with NVIDIA, rather than aarch64, the real v0.1 target?


---

## What an AI OS can do only by owning the OS layer: agent-native kernel and system primitives
_Workflow label: `os_level_primitives`_

### Executive summary
Owning the OS layer gives a real advantage, but a narrower one than the "agent kernel" papers suggest. The flagship research does not run in the kernel. AIOS (COLM 2025) and its semantic file system LSFS (ICLR 2025) run as Python userland layers. AIOS reports up to 2.1x throughput, measured on one RTX A5000, and its last tagged release is v0.2.2 from March 2025. Several other ideas also work in userland serving engines: MemGPT context paging, KV-cache virtual memory (vAttention, kvcached) and token-level GPU pooling (Aegaeon, which Alibaba says cut GPU use by 82% in a beta deployment). Today's coding-agent sandboxes (Anthropic sandbox-runtime, OpenAI Codex) already isolate commands with bubblewrap or Seatbelt without a custom OS. Anthropic reports 84% fewer permission prompts from this approach.

The OS image genuinely matters in five areas:
- **Guaranteed kernel features.** Landlock has reached ABI 11 (UDP rules, TSYNC). BPF LSM needs the lsm= boot parameter, and sched_ext and kprobe override need specific kernel config. Userland tools cannot assume these, so they fall back to weaker modes: Codex demoted Landlock to legacy, and Claude Code needs AppArmor workarounds and a weaker nested mode inside containers.
- **Covering every agent process.** Claude Code's sandbox covers shell commands only. MCP servers, hooks and LSP servers run with full access, and an unsandboxed retry exists.
- **Per-agent identity and network control.** The kernel can identify each agent and enforce its egress rules without relying on proxy environment variables.
- **System-wide undo.** A Btrfs root laid out at install time allows one-step rollback of the whole system, not just files.
- **Checkpoint/restore and patched GPU drivers.** gpu_ext, for example, needs patched NVIDIA modules on Linux 6.15.

The closest existing product is Windows' preview On-device Agent Registry. It runs MCP servers under separate agent accounts, but grants file access per host app rather than per server. No product found combines all of the following on Linux: agent identities enforced by the kernel, eBPF enforcement that knows the agent's intent, per-agent egress with credential injection, and one-command system rollback.

Hardware constraint: on the M1 Pro, the Asahi GPU driver is outside mainline Linux, the Neural Engine driver is out of tree, and all GPU-partitioning research found is NVIDIA-specific.

Method caveat: the shared web-search budget ran out after 7 searches, short of the 20 requested. The remaining evidence comes from about 32 primary sources fetched directly.

### Recommendations
- Do not write a kernel. Build an image-based distribution (bootc or mkosi with systemd v258 or later) whose differentiator is secure-by-default policy composition. Kernel floor: Linux 6.15 or later (Landlock audit) with CONFIG_BPF_LSM, lsm= including landlock and bpf, CONFIG_SCHED_CLASS_EXT, CONFIG_CHECKPOINT_RESTORE and CONFIG_BPF_KPROBE_OVERRIDE enabled.
- Make the first public demo 'agent undo': a Btrfs root, a subvolume per task, automatic snapshots around each agent action, and one command to roll back. It is the most visible capability that userland tools cannot match system-wide, and it carries the least kernel risk.
- Build the agent principal layer on top of sandbox-runtime or bubblewrap rather than replacing them. Wrap every agent, MCP server and hook in a systemd transient unit with Landlock, a seccomp profile and a cgroup-scoped egress policy, and publish the policy format so other tools can target it.
- Begin enforcement in audit-only mode: ship AgentSight-style eBPF correlation first, then promote rules to BPF LSM deny decisions once false-positive rates are measured and published.
- Do not invest in kernel-level GPU scheduling (gpu_ext, LithOS) or the dmem cgroup on the M1 Pro: none apply to Apple GPUs under Linux. Use userland serving (llama.cpp, MLX on a macOS host, Pie or kvcached on NVIDIA) and keep GPU policy in user space.
- Develop primarily in a VM guest on the M1 Pro (libkrun or krunkit, or Apple Containerization) and treat Asahi bare-metal boot as a periodic integration target. Publish an x86_64 ISO for users with NVIDIA GPUs, where checkpointing (cuda-checkpoint) and kvcached work.
- Treat semantic file systems (LSFS) and context paging (MemGPT) as ordinary services, not kernel features. They do not require owning the OS, and the evidence for kernel placement is weak.
- Publish reproducible benchmarks from day one (snapshot latency, sandbox overhead, egress-bypass test results), so the project separates itself from 'AI OS' marketing with hard evidence.
- Write the threat model before writing code: decide whether AgentOS protects the host from the agent, the user's data from exfiltration, or both. The egress and TLS-interception design depends on that decision.

### Pain points

**Userland agent sandboxes leave much of the agent outside the boundary** — severity 4, breadth 5, willingness to pay 3, confidence high
- Segment: Developers running coding agents and MCP tool ecosystems on Linux and macOS workstations
- Description: Current coding-agent sandboxes wrap shell commands, not the whole agent system. MCP servers, hooks, LSP servers, status-line helpers and plugin processes run with full user access. An escape hatch lets the model retry a blocked command unsandboxed, and docker is incompatible with the sandbox.
- Evidence: The Claude Code docs say the sandbox covers shell commands only, and that hooks and local MCP servers run with full access. The docs also describe the dangerouslyDisableSandbox retry and state that docker is incompatible with the sandbox. Anthropic's answer is to wrap the entire Claude Code process in sandbox-runtime, which is a beta research preview.
- Quantified impact: Anthropic reports that sandboxing cut permission prompts by 84%, which shows how much friction it removes for whatever it covers. No figure exists for the risk carried by the unsandboxed components.
- Current solutions: Wrapping the whole process with sandbox-runtime (srt), dev containers with firewall scripts, Docker Sandboxes microVMs, and contained MCP servers in Windows ODR (preview).
- Why still unsolved: A userland tool can only sandbox the processes it launches. Other software starts MCP servers, IDE helpers and browsers under the user's full identity, and only the init system and kernel can apply one policy to every process an agent starts.
- Sources: [Claude Code: Configure the sandboxed Bash tool](https://code.claude.com/docs/en/sandboxing) (accessed 2026-10-16); [Anthropic Engineering: Claude Code sandboxing](https://www.anthropic.com/engineering/claude-code-sandboxing) (2025-10-20); [anthropic-experimental/sandbox-runtime](https://github.com/anthropic-experimental/sandbox-runtime) (accessed 2026-10-16)

**Kernel feature fragmentation forces weaker sandboxes** — severity 3, breadth 4, willingness to pay 2, confidence high
- Segment: Agent tool vendors and security-conscious developers on mixed Linux distributions
- Description: Agent sandboxes must run on whatever kernel and security settings the host provides, so they fall back to the weakest common mode. Codex now defaults to bubblewrap and keeps Landlock only as a legacy option. Ubuntu 24.04 blocks unprivileged user namespaces unless an AppArmor profile or sysctl change allows them. Inside unprivileged containers, bubblewrap cannot mount /proc, so Claude Code offers an explicitly weaker nested mode.
- Evidence: The Codex linux-sandbox README makes bubblewrap the default and marks Landlock as legacy because Landlock cannot isolate app-server Unix sockets. The OpenAI sandboxing docs and the Claude Code docs both describe the AppArmor userns restriction on Ubuntu 24.04. The kernel documentation shows Landlock at ABI 11, with UDP rules only from ABI 10, which only very recent kernels have.
- Quantified impact: No quantified impact found. The visible effect is that sandboxes skip available kernel features (Landlock UDP and TSYNC, BPF LSM) and use more privileged mechanisms such as user namespaces.
- Current solutions: bubblewrap plus seccomp, warning messages at startup, distribution-specific AppArmor profiles, and weaker fallback modes.
- Why still unsolved: Tool vendors cannot control which kernel version, kernel config (CONFIG_BPF_LSM, CONFIG_SCHED_CLASS_EXT, CONFIG_BPF_KPROBE_OVERRIDE) or boot parameters (lsm=) a user has. An OS image controls all three.
- Sources: [Codex linux-sandbox README](https://github.com/openai/codex/blob/main/codex-rs/linux-sandbox/README.md) (accessed 2026-10-16); [Landlock kernel documentation](https://docs.kernel.org/userspace-api/landlock.html) (accessed 2026-10-16); [OpenAI Codex sandboxing docs](https://learn.chatgpt.com/docs/sandboxing) (accessed 2026-10-16)

**Network egress control is hostname-based and can be bypassed** — severity 4, breadth 4, willingness to pay 4, confidence high
- Segment: Enterprises and individuals worried about prompt-injection-driven data exfiltration
- Description: Per-agent egress control usually works through an HTTP or SOCKS proxy that checks the hostname the client supplies, without inspecting TLS. Allowing broad domains such as github.com opens exfiltration paths, and domain fronting can reach hosts outside the allowlist. Docker Sandboxes cannot chain non-HTTP TCP traffic to an upstream proxy. The Windows alpha of sandbox-runtime does not fence DNS.
- Evidence: The Claude Code docs warn about domain fronting and state that TLS is not inspected by default. Docker Sandboxes' architecture page says proxy chaining applies only to HTTP and HTTPS. The srt README says DNS is not fenced on Windows.
- Quantified impact: Not quantified in the sources reviewed.
- Current solutions: Host-side proxies with domain allowlists, experimental TLS termination with credential masking, Docker Sandboxes' sentinel-credential substitution, and custom MITM proxies.
- Why still unsolved: Strong per-agent egress needs four things at once: kernel-level identity for each agent (cgroup or socket hooks), transparent redirection that does not depend on proxy environment variables, DNS control, and a trusted CA installed in the system store. All four are OS-image concerns that no single tool controls.
- Sources: [Claude Code sandboxing: Security limitations](https://code.claude.com/docs/en/sandboxing) (accessed 2026-10-16); [Docker Sandboxes architecture](https://docs.docker.com/ai/sandboxes/architecture/) (accessed 2026-10-16); [sandbox-runtime README](https://github.com/anthropic-experimental/sandbox-runtime) (accessed 2026-10-16)

**No system-wide transactional undo for agent actions** — severity 4, breadth 4, willingness to pay 3, confidence medium
- Segment: Developers who let agents run autonomously; operators of long-running agents
- Description: When an agent breaks something, users can roll back files in a workspace, but not the system state the agent changed (installed packages, config, services). Existing options are file-only or tied to a single app.
- Evidence: Fly.io Sprites offers whole-filesystem checkpoints in about a second, but restores roll back files only, not running processes. Turso's AgentFS (3.5k stars) offers copy-on-write isolation and SQL audit over FUSE or NFS, but is in beta. Btrfs snapshots are instant metadata copies, but they are not recursive across nested subvolumes and are not backups.
- Quantified impact: Not quantified. The demand signal is AgentFS at 3.5k stars and Sprites building its product around checkpoints.
- Current solutions: App-level checkpoints in coding agents, git, AgentFS, cloud VM snapshots (Sprites, Firecracker), and manual Btrfs or ZFS snapshots with snapper.
- Why still unsolved: Atomic rollback of a workspace together with the system it changed requires choosing the filesystem and subvolume layout at install time, plus an immutable or layered root (bootc, systemd-sysext). Only the OS image can guarantee that.
- Sources: [Fly.io Sprites](https://fly.io/sprites) (accessed 2026-10-16); [tursodatabase/agentfs](https://github.com/tursodatabase/agentfs) (accessed 2026-10-16); [Btrfs Subvolumes documentation](https://btrfs.readthedocs.io/en/latest/Subvolumes.html) (accessed 2026-10-16)

**Tools either observe agents or enforce rules, but none links the two to intent** — severity 3, breadth 3, willingness to pay 3, confidence medium
- Segment: Security teams and power users auditing autonomous agents
- Description: eBPF observability tools for agents link LLM prompts to system calls but cannot block anything. eBPF enforcement tools can block actions but do not know which prompt or tool call caused them. Tetragon's SIGKILL alone does not guarantee that a write is stopped; it must be combined with return-value override.
- Evidence: AgentSight (PACMI 2025) captures TLS LLM traffic and correlates it with process and file events at under 3% overhead, but it only observes and needs root. Tetragon's docs state that SIGKILL alone may let a write complete and that Override must be added.
- Quantified impact: AgentSight reports under 3% overhead. No measurement of detection coverage was found.
- Current solutions: AgentSight, Tetragon, Falco (detection), and proxy logs.
- Why still unsolved: An intent-aware deny decision needs BPF LSM enabled at boot, kprobe override in the kernel config, and a trusted daemon linking agent identity to its LLM session. Distributions do not ship that combination by default.
- Sources: [AgentSight paper (arXiv 2508.02736)](https://arxiv.org/abs/2508.02736) (2025-08-15); [eunomia-bpf/agentsight](https://github.com/eunomia-bpf/agentsight) (accessed 2026-10-16); [Tetragon enforcement concepts](https://tetragon.io/docs/concepts/enforcement/) (accessed 2026-10-16)

**Sharing one GPU between interactive inference and background fine-tuning lacks OS controls** — severity 3, breadth 3, willingness to pay 3, confidence medium
- Segment: Local AI developers and small labs with one or two GPUs
- Description: Running interactive agents, model serving and training on one accelerator suffers from head-of-line blocking and weak isolation. The kernel has little GPU control: the dmem cgroup only limits device memory for TTM-based drivers (Intel xe first), and research schedulers need patched NVIDIA drivers.
- Evidence: LithOS reports 4.7x lower tail latency than NVIDIA MPS for mixed inference and training. gpu_ext reports up to 4.8x decode throughput but needs instrumented NVIDIA open kernel modules on Linux 6.15. The kernel cgroup-v2 docs show dmem regions only in the form drm/<PCI>/vram0 for xe.
- Quantified impact: The research claims above (13x and 4.7x tail-latency gains against MPS) imply large untapped gains, but none of it ships in products.
- Current solutions: NVIDIA MPS, MIG (datacenter GPUs), time-slicing, and userland serving tricks (kvcached, Aegaeon-style token-level preemption).
- Why still unsolved: GPU drivers are proprietary and differ by vendor. Kernel GPU cgroups are new and partial, and the research needs custom kernels and drivers. Apple Silicon under Linux has none of these mechanisms.
- Sources: [LithOS (arXiv 2504.15465)](https://arxiv.org/abs/2504.15465) (2025-04-21); [gpu_ext (arXiv 2512.12615)](https://arxiv.org/html/2512.12615v2) (2025-12); [Linux cgroup-v2 admin guide (dmem)](https://docs.kernel.org/admin-guide/cgroup-v2.html) (accessed 2026-10-16)

**Checkpoint/restore of agent processes that use GPUs is partial and risky** — severity 3, breadth 2, willingness to pay 3, confidence medium
- Segment: Agent platforms needing pause/resume, branching and preemption
- Description: Pausing, migrating or forking a running agent along with its model state is only partly supported. NVIDIA's cuda-checkpoint does not support UVM memory. Cloning microVM snapshots duplicates randomness, tokens and identifiers unless uniqueness is enforced.
- Evidence: cuda-checkpoint integrates with CRIU 4.0+ from driver 570, adds GPU migration in 580, and does not support UVM memory. CRIUgpu (Feb 2025) claims transparent CUDA and ROCm checkpointing but gives only qualitative claims in its abstract. Firecracker's docs warn that resuming the same snapshot more than once is insecure without uniqueness handling.
- Quantified impact: Not quantified in the sources reviewed.
- Current solutions: CRIU with cuda-checkpoint, Firecracker snapshots, Sprites (files only), and app-level state serialization.
- Why still unsolved: CRIU needs kernel config and privileges, GPU support depends on the driver, and Apple and Asahi GPUs have no checkpoint path.
- Sources: [NVIDIA cuda-checkpoint](https://github.com/NVIDIA/cuda-checkpoint) (accessed 2026-10-16); [CRIUgpu (arXiv 2502.16631)](https://arxiv.org/abs/2502.16631) (2025-02-23); [Firecracker snapshot support](https://github.com/firecracker-microvm/firecracker/blob/main/docs/snapshotting/snapshot-support.md) (accessed 2026-10-16)

### Key statistics

| Claim | Value | Source type | Date | Source |
|---|---|---|---|---|
| Reduction in permission prompts after Claude Code adopted OS-level sandboxing (internal usage) | 84% | vendor | 2025-10-20 | https://www.anthropic.com/engineering/claude-code-sandboxing |
| AIOS throughput improvement (Reflexion, Llama-3.1-8B, HumanEval) on a single RTX A5000 | up to 2.1x | academic | 2025-05 | https://arxiv.org/html/2403.16971v4 |
| AIOS latest tagged release and GitHub stars | v0.2.2 (2025-03-12); 6.5k stars | primary_research | accessed 2026-10-16 | https://github.com/agiresearch/AIOS |
| LSFS retrieval accuracy and speed vs a traditional file system | >=15% higher accuracy; 2.1x faster semantic retrieval | academic | 2025-04 (ICLR 2025) | https://proceedings.iclr.cc/paper_files/paper/2025/file/517eb19e99947f60afff0cf93e451825-Paper-Conference.pdf |
| Highest Landlock ABI documented; UDP restriction added in ABI 10 | ABI 11 | official | accessed 2026-10-16 (docs build 7.3-rc) | https://docs.kernel.org/userspace-api/landlock.html |
| Landlock audit support and IPC scoping kernel releases | Linux 6.15 (2025-05-25) audit; Linux 6.12 (2024-11-17) abstract Unix socket and signal scoping | official | 2025-05-25 | https://landlock.io/news/ |
| Coding-agent sandbox library adoption | sandbox-runtime: 5.6k stars, releases up to v0.0.79 | vendor | accessed 2026-10-16 | https://github.com/anthropic-experimental/sandbox-runtime |
| gpu_ext eBPF GPU driver policies: throughput and tail-latency gains | up to 4.8x throughput; up to 2x lower tail latency; 95% lower P99 launch latency for latency-critical work when sharing a GPU between tenants | academic | 2025-12 | https://arxiv.org/html/2512.12615v2 |
| LithOS tail latency vs NVIDIA MPS (inference stacking / hybrid inference-training) | 13x / 4.7x lower | academic | 2025-04-21 | https://arxiv.org/abs/2504.15465 |
| Aegaeon token-level multi-model GPU pooling in Alibaba's beta deployment | 1,192 to 213 GPUs (82% reduction) | news | 2025-10 (SOSP 2025) | https://www.datacenterdynamics.com/en/news/alibaba-cloud-claims-it-can-reduce-gpu-use-by-82-with-pooling-system/ |
| kvcached TTFT reduction for three Llama-3.1-8B models on one A100-80G | 2-28x | primary_research | accessed 2026-10-16 (v0.1.6) | https://github.com/ovg-project/kvcached |
| vAttention throughput vs PagedAttention kernels using CUDA virtual memory APIs | up to 1.23x | academic | 2025-01-29 (ASPLOS 2025) | https://arxiv.org/abs/2405.04437 |
| AgentSight eBPF agent observability overhead | <3% CPU | academic | 2025-08-15 | https://arxiv.org/abs/2508.02736 |
| cgroup v2 dmem controller for GPU device memory limits | Merged for Linux 6.14; initially wired to Intel xe (regions drm/<PCI>/vram0) | news | 2025-01 | https://www.phoronix.com/news/DMEM-cgroup-vRAM-Control |
| systemd v258 release adding PrivateBPF, BPFDelegate*, unprivileged nspawn and vmspawn; cgroup v1 removed | Released 2025-09-17 | official | 2025-09-17 | https://github.com/systemd/systemd/releases/tag/v258 |
| Windows runs contained MCP servers under a separate agent user account and Windows session | Default for packaged MCP servers (preview); file consent is granted per host app, not per server | official | 2025-11-11 | https://learn.microsoft.com/en-us/windows/ai/mcp/servers/mcp-containment |
| Asahi Linux status on M1 Pro/Max | GPU driver in linux-asahi (outside mainline); Neural Engine driver out of tree; M3 support added 2026-09-06 | official | accessed 2026-10-16 | https://asahilinux.org/docs/platform/feature-support/m1/ |
| Fly.io Sprites whole-filesystem checkpoint time | about 1 second; restores files, not processes | vendor | accessed 2026-10-16 | https://fly.io/sprites |
| cuda-checkpoint CRIU integration and limits | CRIU 4.0+ integration from driver 570; GPU migration from 580; UVM memory not supported | vendor | accessed 2026-10-16 | https://github.com/NVIDIA/cuda-checkpoint |
| WASI 0.3 (native async component model) shipped | 2026-06-11 (0.3.1 on 2026-08-11) | official | 2026-08 | https://wasi.dev/roadmap |
| Kubernetes SIG Apps agent-sandbox adoption | v1.0.2; 4.3k stars | official | accessed 2026-10-16 | https://github.com/kubernetes-sigs/agent-sandbox |

### Players

- **AIOS / Cerebrum (Rutgers, agiresearch)** (Academic agent OS (userland)): Agent 'kernel' with a scheduler (FIFO/RR), context snapshots (text or logits), memory, storage and tool managers, and the LSFS semantic file system. Traction: 6.5k stars; COLM 2025, ICLR 2025 (LSFS), NAACL 2025 (Cerebrum). Gap: Python userland on top of an ordinary OS; last release v0.2.2 (March 2025); not bootable; no kernel-enforced isolation.
- **Anthropic sandbox-runtime (srt)** (Userland agent sandbox): bubblewrap plus seccomp on Linux, Seatbelt on macOS, a host proxy for domain allowlists, and a log of policy violations. Traction: 5.6k stars; built into Claude Code. Gap: Beta; built-in Bash sandbox excludes MCP servers and hooks; hostname-only egress checks without TLS inspection; Windows support is alpha.
- **OpenAI Codex sandbox** (Userland agent sandbox): bubblewrap (default) plus seccomp on Linux, Seatbelt on macOS, native sandbox on Windows; read-only and writable-root modes; .git and .codex kept read-only. Traction: Bundled with Codex CLI, IDE extension and desktop app. Gap: Landlock demoted to legacy; depends on unprivileged user namespaces and an AppArmor workaround on Ubuntu 24.04.
- **Microsoft Windows On-device Agent Registry (ODR) and agent sessions** (OS vendor agent primitives): MCP servers run contained under separate agent accounts and Windows sessions; Intune controls; audit logging; File Explorer MCP connector. Traction: Public preview (docs dated 2025-11 to 2026-06). Gap: Windows-only and proprietary; file consent per host app, not per server; no rollback or GPU scheduling.
- **Docker Sandboxes** (MicroVM agent isolation): MicroVM per sandbox with a private Docker daemon, host proxy, sentinel-credential injection and an MCP gateway. Traction: Free local CLI; paid central policy. Gap: Only HTTP and HTTPS can be chained to an upstream proxy; a product layered on the host, not an OS.
- **Fly.io Sprites** (Cloud stateful agent VMs): Persistent Linux machines for agents with automatic copy-on-write checkpoints in about 1 second. Traction: Commercial, usage-based pricing. Gap: Cloud-only; restores files only, not processes; not local or bootable.
- **Kubernetes agent-sandbox (SIG Apps)** (Orchestration): Sandbox custom resource, warm pools, gVisor or Kata runtimes. Traction: v1.0.2; 4.3k stars. Gap: Cluster-scale only; snapshots and hibernation not confirmed.
- **Turso AgentFS** (Agent filesystem): SQLite-backed copy-on-write filesystem with key-value store, SQL audit of tool calls, and FUSE or NFS mounting. Traction: 3.5k stars. Gap: Beta; FUSE overhead; does not cover system-wide state.
- **AgentSight (eunomia-bpf)** (eBPF agent observability): TLS payload capture plus syscall tracing, linked to agent intent; OpenTelemetry export. Traction: 745 stars; PACMI 2025. Gap: Observation only; needs root; statically linked agents need extra configuration.
- **Cilium Tetragon / Falco** (eBPF runtime enforcement and detection): Kprobe and LSM policies with Override and SIGKILL actions (Tetragon); rule-based detection (Falco). Traction: Widely used in Kubernetes. Gap: No concept of an agent or its LLM intent; Override needs specific kernel config.
- **kvcached (ovg-project)** (KV-cache memory virtualization): Elastic, on-demand GPU memory for KV caches in vLLM and SGLang; multiple models on one GPU; colocation with training. Traction: 1.7k stars; v0.1.6. Gap: CUDA and ROCm only; no Apple Metal support.
- **Pie (SOSP 2025)** (Programmable LLM serving): Wasm 'inferlets' that own their KV cache and control the agent loop; CUDA, Metal, Vulkan and WebGPU backends. Traction: 226 stars; pre-release. Gap: Early stage; userland engine, not an OS.
- **gpu_ext and LithOS (research)** (Kernel and driver-level GPU scheduling): eBPF hooks in the NVIDIA driver for memory and scheduling (gpu_ext); per-TPC scheduling and kernel atomization (LithOS). Traction: arXiv and LPC 2025; not productized. Gap: NVIDIA-only; needs patched drivers and kernels.
- **NVIDIA cuda-checkpoint + CRIU / CRIUgpu** (Checkpoint/restore): Transparent checkpoint of processes using CUDA (and ROCm in CRIUgpu). Traction: Shipped in NVIDIA drivers 550+. Gap: No UVM support; no Apple GPU path.
- **Letta (MemGPT) and MemOS** (Agent memory 'OS'): Virtual context paging between memory tiers (MemGPT); MemCube memory units with lifecycle management (MemOS). Traction: MemGPT widely cited; MemOS v4 December 2025. Gap: Purely application-layer; nothing requires kernel ownership.
- **Apple Containerization / libkrun (krunkit)** (VM substrate on Apple Silicon): One lightweight VM per container (Apple); HVF virtualization with virtio-gpu Venus or native context (libkrun). Traction: Apple Containerization 9.0k stars. Gap: Apple requires macOS 26 and does not mention GPU support; libkrun treats guest and VMM as one security context.
- **Asahi Linux** (Hardware enablement): Linux on Apple Silicon M1 to M3; upstream Mesa GPU drivers. Traction: Seven-person governance since 2025-02; M3 support since 2026-09. Gap: M1 Pro GPU kernel driver still outside mainline; Neural Engine out of tree; no GPU partitioning or checkpoint support.

### Trends
- The isolation boundary is moving to a microVM per agent: Docker Sandboxes, Apple Containerization (one VM per container), Kubernetes agent-sandbox with gVisor or Kata, Firecracker-based clouds, and libkrun on macOS.
- OS vendors are adding agent identities to the OS: Windows runs MCP servers under separate agent accounts and sessions, with admin policy and audit (preview, November 2025 onward).
- eBPF is becoming the policy layer for everything: sched_ext (CPU scheduling), BPF LSM (access control), gpu_ext (GPU memory and scheduling), AgentSight (agent tracing), and systemd v258's PrivateBPF and BPFDelegate options.
- OS-style virtual memory for KV caches is happening in userland, not in kernels: vAttention (ASPLOS 2025), kvcached, Aegaeon token-level pooling (SOSP 2025), and Pie inferlets.
- Landlock is maturing fast (ABI 11, including UDP, TSYNC, audit and IPC scoping), yet the major coding agents rely on bubblewrap because kernel versions vary across users. Owning the OS image closes that gap.
- Checkpoints are becoming a product feature: Sprites (about 1 second, files only), CRIU with cuda-checkpoint, Firecracker snapshots, and AgentFS time-travel.
- Wasm is becoming a capability sandbox for tools: WASI 0.3 shipped in June 2026, Microsoft Wassette runs MCP tools as Wasm components, and Pie runs Wasm inferlets.
- 'Agent OS' research mostly runs above the kernel (AIOS, LSFS, MemGPT, MemOS). The truly kernel-level research (gpu_ext, LithOS) targets NVIDIA datacenter workloads.
- Apple Silicon Linux is stabilizing: Asahi Mesa drivers went upstream in August 2025 and M3 support arrived in September 2026, but GPU compute and the Neural Engine remain second-class.

### Opportunities

- **Agent principals as a first-class OS object: every agent, tool and MCP server launches as its own systemd transient unit or cgroup, with stacked Landlock rules, a BPF LSM policy and per-agent egress identity. A user should be able to ask 'what can agent X touch?' and get a kernel-enforced answer.** Why now: Windows shipped agent accounts and sessions in preview, which validates the concept. Linux has every needed piece (Landlock ABI 11, BPF LSM, cgroup socket hooks, systemd v258 delegation), but no distribution combines them by default. Userland sandboxes leave MCP servers and hooks unsandboxed. Evidence: Claude Code docs (MCP servers and hooks run outside the sandbox); Windows MCP containment docs (separate agent user account); Landlock kernel docs (ABI 11); systemd v258 release notes. Risk: Duplicates srt and bubblewrap if done shallowly. The fix is to build the policy compiler on top of srt and bwrap rather than replacing them, so the difference is OS-wide default enforcement.
- **Transactional agent workspaces with system-wide undo: a Btrfs subvolume per task, automatic snapshots before and after each action, an immutable root with layered extensions (bootc or systemd-sysext), and a single 'agent undo' that rolls back files, installed packages and config together.** Why now: Agents increasingly act autonomously. Sprites and AgentFS show demand, but both work at the file or FUSE level, and only an OS that chooses the root filesystem layout at install time can include system state. Evidence: Sprites checkpoints in about 1 second with file-only restores; AgentFS has 3.5k stars and is in beta; Btrfs snapshots are instant metadata copies. Risk: Btrfs snapshots are not recursive and not backups. Databases and running services need quiescing. Needs careful subvolume design.
- **Intent-aware enforcement: combine AgentSight-style TLS and syscall correlation with BPF LSM deny decisions, so every blocked action is linked to the prompt and tool call that caused it, plus a replayable audit trail.** Why now: AgentSight (2025) showed the correlation costs under 3% overhead. Tetragon proved in-kernel enforcement. Owning the kernel command line (lsm=) and config (kprobe override) removes the deployment barrier. Evidence: AgentSight paper and repository; Tetragon enforcement docs; BPF LSM kernel docs. Risk: TLS probes break on statically linked or pinned TLS libraries, and false positives may annoy users. Begin in audit-only mode.
- **Per-agent egress firewall with credential injection: transparent redirection by cgroup (no reliance on proxy environment variables), DNS fencing, and a system CA so the proxy can inspect TLS and substitute credentials for allowed hosts.** Why now: Hostname-only proxies can be bypassed by domain fronting, and Docker Sandboxes already proved that sentinel-credential injection works in a microVM. Evidence: Claude Code security limitations; Docker Sandboxes architecture; srt Windows DNS caveat. Risk: Intercepting TLS creates a high-value target, and certificate pinning in some tools will break.
- **Ship one image as both a bootable OS and a VM guest: boot on x86 with NVIDIA and on Asahi, and run as a libkrun or Apple Containerization guest on macOS so the founder can develop on the M1 Pro without rebooting.** Why now: libkrun supports HVF and virtio-gpu (Venus or native context). Apple Containerization ships minimal kernels from 6.14.9. Asahi now covers M1 to M3. Evidence: libkrun README; Apple Containerization README; Asahi blog and feature table. Risk: GPU performance inside the VM, and Asahi's GPU kernel driver living outside mainline, may weaken the local-model story on Apple hardware.
- **Memory governor for model memory on unified-memory machines: use pressure-stall information (PSI) to load, unload and offload KV caches across agents, giving model memory its own quota so a 16 GB machine can multiplex models and agents.** Why now: kvcached and Aegaeon proved that elastic KV and model memory pays off, but only for CUDA and ROCm; unified-memory machines (Apple Silicon, AMD APUs) lack an equivalent. Evidence: kvcached README (CUDA and ROCm only); Aegaeon token-level scaling; cgroup v2 PSI interfaces. Risk: Achievable in userland, so it is a weak moat. Its value comes from being on by default in the OS rather than from kernel exclusivity.

### Contrarian signals
- The best-known 'agent operating system' (AIOS) is a Python userland framework that calls itself an abstraction layer over the host kernel. It shows that agent scheduling and context management do not need OS ownership. Its last tagged release is March 2025.
- Anthropic achieved an 84% cut in permission prompts and OpenAI ships cross-platform sandboxes without owning any OS, so much of the safety benefit is already available as a userland download.
- The industry is moving the isolation boundary into microVMs (Docker Sandboxes, Apple Containerization, agent-sandbox with gVisor or Kata, Sprites). The valuable product may be a guest image, not a bootable host OS.
- Both sandbox-runtime and the Claude Code docs note that a sandbox sharing the host kernel can be escaped through a kernel vulnerability. Owning the OS does not remove that risk; VMs reduce it.
- Every kernel-level GPU scheduling result found (gpu_ext, LithOS) is NVIDIA-specific and needs patched drivers, and dmem cgroup support began with Intel xe only. On the founder's M1 Pro, none of these mechanisms exist.
- On the M1 Pro, Asahi's GPU kernel driver is outside mainline and the Neural Engine driver is out of tree. Linux may therefore run local models worse than macOS with MLX on the same hardware, which undermines a 'fast local models' pitch on this hardware.
- Snapshots are not undo for the outside world: Btrfs snapshots are non-recursive and not backups, Sprites restores only files, and no filesystem rollback can reverse an email sent or an API call made.
- Cloning checkpoints for agent branching duplicates randomness and secrets (Firecracker warning), and cuda-checkpoint does not support UVM. 'Fork a running agent' is harder than it sounds.
- LithOS and gpu_ext headline numbers are measured against weak baselines (MPS, framework offloading). Gains against the best prior systems are much smaller (for example, 1.18x for LithOS in mixed inference and training).

### Open questions
- Clarifying question for the founder: is the primary target hardware Apple Silicon (Asahi) or x86 with NVIDIA? GPU scheduling, checkpointing and KV-cache primitives are available only on the latter.
- Clarifying question for the founder: would a single image that boots on bare metal and also runs as a VM guest satisfy 'bootable OS'? A VM guest captures most of the security value with far less hardware friction.
- Which threat model comes first: host integrity, data exfiltration, or protecting users from destructive actions? Each points to a different first feature (sandbox, egress firewall, or undo).
- Unverified because the search budget ran out: details of Windows Agent Workspace beyond the MCP containment docs; Falco's agent-specific rules; seL4 or Capsicum applied to agents; MIG supported-GPU list; Llumnix v1 architecture; Kernel version numbers for Landlock ABIs 8 to 11 (unverified).
- How well does LLM inference perform on Asahi's Vulkan (Honeykrisp) driver compared with MLX on macOS for the same M1 Pro? No 2025-2026 benchmark was found.
- Can the BPF LSM and Landlock policies be generated automatically from an agent's declared tool manifest (for example, MCP tool schemas) without overwhelming users with prompts?
- Will Windows ODR's agent-account model reach general availability and set the de facto design pattern that Linux users expect?

### Adversarial verification

- **REFUTED** — AIOS is effectively dormant: its last tagged release is v0.2.2 from March 2025 (2025-03-12), with 6.5k stars and no 2026 activity. The brief uses this as a contrarian signal.
  - Corrected: AIOS published v0.3.0 on 2026-01-22. The tag was created 2025-11-24, and GitHub marks it 'Latest'. v0.2.2 was published on 2025-03-08, not 2025-03-12; the 03-12 date comes from the README news list. Commits continued until 2026-07-20, including per-user memory isolation work. The repository has 6,478 stars. Two claims hold: the README news list has no 2026 items, and the Rust rewrite (aios-rs) has placeholder code only and says it is 'NOT feature-parity yet'.
  - Evidence: The GitHub releases API lists v0.3.0 published 2026-01-22T14:07:28Z and v0.2.2 published 2025-03-08T17:11:40Z. The repository API shows 6,478 stars and a last push on 2026-07-20. The v0.3.0 notes include an MCP computer-use environment, Qdrant support and the aios-rs scaffold. The README still calls AIOS 'an abstraction layer over the operating system kernel', so the claim that it runs in userland is correct.
  - Sources: [agiresearch/AIOS releases](https://github.com/agiresearch/AIOS/releases) (2026-01-22); [GitHub API: AIOS releases](https://api.github.com/repos/agiresearch/AIOS/releases) (accessed 2026-10-16); [agiresearch/AIOS README](https://github.com/agiresearch/AIOS) (accessed 2026-10-16)
- **PARTIALLY_CONFIRMED** — Landlock has reached ABI 11 (UDP rules, TSYNC). The recommendation sets a kernel floor of Linux 6.15 (Landlock audit) so the OS image can use these features.
  - Corrected: ABI 11 appears only in Linux 7.3-rc, which was unreleased on 2026-10-16 (the latest tag is v7.3-rc7). ABI 11 adds LANDLOCK_RESTRICT_SELF_NO_NEW_PRIVS, not UDP or TSYNC. The mapping is:
- Linux 6.15 to 6.19: ABI 7 (audit).
- Linux 7.0: ABI 8 (TSYNC).
- Linux 7.1: ABI 9 (pathname Unix socket control, FS_RESOLVE_UNIX).
- Linux 7.2: ABI 10 (UDP).
A 6.15 floor therefore gives neither UDP rules nor TSYNC. The floor must be at least 7.2 for UDP and 7.3 for ABI 11.
  - Evidence: In kernel source, security/landlock/syscalls.c sets landlock_abi_version to 7 at v6.15, v6.16, v6.18 and v6.19. It is 8 at v7.0, 9 at v7.1, 10 at v7.2 and 11 at v7.3-rc7. The docs.kernel.org page is built from 7.3.0-rc7. It lists TSYNC under ABI 8, UDP under ABI 10 and NO_NEW_PRIVS under ABI 11, and it does not map ABIs to kernel releases. The landlock.io news page stops at 6.15. This resolves the brief's open question about which kernels carry ABIs 8 to 11.
  - Sources: [Landlock userspace API documentation](https://docs.kernel.org/userspace-api/landlock.html) (accessed 2026-10-16 (build 7.3.0-rc7)); [Linux v7.2 security/landlock/syscalls.c](https://raw.githubusercontent.com/torvalds/linux/v7.2/security/landlock/syscalls.c); [Linux v7.3-rc7 security/landlock/syscalls.c](https://raw.githubusercontent.com/torvalds/linux/v7.3-rc7/security/landlock/syscalls.c); [Landlock news](https://landlock.io/news/) (2025-05-25)
- **CONFIRMED** — The Claude Code sandbox covers shell commands only. MCP servers, hooks and LSP servers run with full access, an unsandboxed retry exists, and docker is incompatible. Anthropic reports 84% fewer permission prompts.
  - Corrected: Confirmed, with three omissions:
- The sandbox is off by default.
- The retry can be disabled with allowUnsandboxedCommands:false, and failIfUnavailable stops unsandboxed fallback. On Linux, CLAUDE_CODE_SUBPROCESS_ENV_SCRUB sandboxes every command, including shell-mode commands.
- The 84% figure is 'in our internal usage' with no methodology. sandbox-runtime now lives at anthropics/sandbox-runtime.
  - Evidence: The docs say 'The sandbox covers shell commands only.' They list command hooks, local MCP servers, plugin monitors, LSP servers, the status line and apiKeyHelper as running 'with your full access'. They also document the dangerouslyDisableSandbox retry, say docker is incompatible, warn about domain fronting, and state that TLS is not inspected by default. They add that 'Stronger TLS-aware network isolation is an active area of development.' The engineering post is dated 2025-10-20.
  - Sources: [Claude Code: Configure the sandboxed Bash tool](https://code.claude.com/docs/en/sandboxing) (accessed 2026-10-16); [Anthropic Engineering: Claude Code sandboxing](https://www.anthropic.com/engineering/claude-code-sandboxing) (2025-10-20); [anthropics/sandbox-runtime](https://github.com/anthropics/sandbox-runtime) (v0.0.79, 2026-10-07)
- **PARTIALLY_CONFIRMED** — Codex demoted Landlock to legacy and defaults to bubblewrap because kernel versions vary across users, and owning the OS image closes that gap.
  - Corrected: The Codex README makes bubblewrap the default and keeps 'the legacy Landlock filesystem implementation'. It gives one reason: Landlock is rejected for filesystem-restricted policies 'because it cannot isolate app-server Unix sockets'. It does not mention kernel-version variance. Pathname Unix socket control arrived in Landlock ABI 9 (Linux 7.1). An OS image could address the stated reason, but only with a floor of 7.1 or later, not 6.15.
  - Evidence: The README states 'Bubblewrap is the default filesystem sandbox.' The OpenAI sandboxing page confirms bwrap on Linux, the Ubuntu 24.04 user-namespace warning (with the bwrap-userns-restrict profile), Seatbelt on macOS and a native Windows sandbox. It does not mention Landlock or seccomp.
  - Sources: [Codex linux-sandbox README](https://github.com/openai/codex/blob/main/codex-rs/linux-sandbox/README.md) (accessed 2026-10-16); [OpenAI Codex sandboxing docs](https://learn.chatgpt.com/docs/sandboxing) (accessed 2026-10-16)
- **PARTIALLY_CONFIRMED** — On the M1 Pro, the Asahi GPU driver is outside mainline and the Neural Engine driver is out of tree. M3 support arrived 2026-09-06, so Asahi now covers M1 to M3.
  - Corrected: The M1 Pro facts are confirmed: the GPU and DCP rows say linux-asahi, and the ANE row says out of tree. The Linux 7.2 report (2026-08-26) says GPU and DCP work 'cannot' be upstreamed yet. M3 support (2026-09-06) is official but limited:
- No performant GPU acceleration.
- Sleep does not work.
- HDMI is disabled.
- The M3 Ultra is unsupported.
'Covers M1 to M3' overstates the case for GPU or local-model work.
  - Evidence: The feature table has no update date. The M3 post warns 'Do not expect performant or power-efficient 3D acceleration right now.' The Mesa upstreaming in August 2025 is confirmed (Linux 6.16 report, 2025-08-07).
  - Sources: [Asahi M1 feature support](https://asahilinux.org/docs/platform/feature-support/m1/) (accessed 2026-10-16); [Asahi Linux on M3](https://asahilinux.org/2026/09/m2-episode-1/) (2026-09-06); [Progress Report: Linux 7.2](https://asahilinux.org/2026/08/progress-report-7-2/) (2026-08-26)

**Missed counter-evidence**
- systemd has moved past v258 to v262 (2026-09-22). These releases add features relevant to the 'agent principal as transient unit' design:
- v262: Varlink StartTransient accepts sandboxing, credential and namespace properties.
- v261: RestrictFileSystemAccess= (BPF LSM, dm-verity-only exec) and nspawn RestrictAddressFamilies=.
- v260: mstack, OCI image pull, and FirewallMark= for resolved delegates.
The 'v258 or later' recommendation and the 'no egress features' caveat are stale.
- The Landlock ABI-to-kernel mapping (7.0 TSYNC, 7.1 pathname Unix sockets, 7.2 UDP, 7.3-rc ABI 11) invalidates the 6.15 kernel floor. It also shows that the reason Codex gives for treating Landlock as legacy (Unix socket isolation) is addressable on 7.1 or later.
- AIOS v0.3.0 (January 2026) and commits through July 2026 weaken the contrarian claim that the flagship agent OS is stagnant.
- sandbox-runtime's Windows alpha already uses a dedicated sandbox user account and a WFP egress fence. A userland vendor tool is moving toward per-agent OS identity, which narrows the 'agent principal' moat.
- Claude Code has added strictAllowlist, failIfUnavailable, tlsTerminate with credential masking, and SUBPROCESS_ENV_SCRUB full sandboxing. Vendors are closing gaps the brief attributes to the lack of OS ownership.
- gpu_ext was evaluated mainly on a consumer RTX 5090, not datacenter GPUs. This contradicts the trend claim and matters for an x86 NVIDIA target.
- Adoption figures have moved: Kubernetes agent-sandbox is at v1.0.7 (2026-10-15), not v1.0.2. AgentFS's latest release is v0.6.4 (2026-03-25).

**Reliability note:** I checked primary sources directly: kernel source tags, GitHub API, vendor docs and the Asahi blog. The shared web-search budget ran out before my first search, so I fetched these sources by URL instead. Three areas remain unverified: whether Windows ODR has reached general availability, new 'AI OS' distributions launched in 2026, and newer Codex sandbox changes. Most of the brief's factual claims hold. The three most consequential errors are AIOS's release status, the Landlock kernel floor, and the M3 support scope. The founder should raise the kernel floor to 7.2 or later before relying on Landlock UDP or TSYNC. Two clarifying questions remain: will the image track mainline 7.x kernels, given Asahi's out-of-tree GPU driver? Is UDP or Unix-socket confinement actually required for the first threat model?


---

## Existing AI-focused Linux distributions and OS images: inventory, traction and gaps (as of 2026-10-16)
_Workflow label: `existing_ai_distros`_

### Executive summary
No product yet treats agents as first-class OS principals. The market splits into four groups, and none combines a bootable image, agent-level permissions and memory-aware local model scheduling.

(1) Vendor appliance stacks. These are NVIDIA DGX OS 7.5.0 for DGX Spark (Ubuntu base, kernel 6.17, CUDA 13.0.2), Red Hat Enterprise Linux AI 3.x (a bootc image, inference only, supporting only NVIDIA CUDA and AMD ROCm) and Lambda Stack. All are tied to one hardware vendor and target enterprises or GPU owners. NVIDIA joins its parts together through more than 100 playbooks (step-by-step setup guides) rather than one designed system. Its issue tracker shows playbooks failing on fresh machines and GB10 throughput dropping 25 to 40 percent over uptime.

(2) Platform vendors adding AI layers. These are Canonical inference snaps (beta, 19 models, 52 GitHub stars), SUSE Linux Enterprise Server 16 with MCP (tech preview) and Fedora's AI/ML SIG (packaging stalled on ROCm and libstdc++ breakage). Red Hat retired InstructLab: the repository is archived, and its parts became separate SDKs.

(3) Opinionated desktops. Omarchy is the breakout. It has 44,921 stars, reports 1.28 million ISO downloads in year one, and has $23.65 million pledged to a foundation. It now calls itself an "agentic" distribution. However, its launcher starts every coding agent with its permission prompts switched off ("yolo", "--dangerously-skip-permissions", "bypassPermissions"). Bluefin and Bluefin GDX (CentOS bootc with NVIDIA drivers, uv and nvtop) and Pop!_OS COSMIC offer GPU-ready desktops but little that is specific to agents.

(4) Agent runtimes and sandboxes outside any distro. These are OpenClaw (392k stars, 722 GitHub security advisories in 2026, 14 critical), NVIDIA OpenShell and NemoClaw (alpha, Apache-2.0), AIOS (a Python runtime, not bootable), RamaLama and Pinokio (scripts "can run anything").

Hard evidence points to three unserved needs:
- **OS-enforced agent permissions:** agents get their own identity, file and network policy, a credential broker and human approval.
- **Model scheduling for unified memory:** DGX Spark's documentation admits that cudaMemGetInfo under-reports and nvidia-smi shows memory as unsupported.
- **Verified, reversible agent execution:** signed action logs tied to atomic rollback.

Native Apple Silicon is also open territory. Omarchy supports it only in a VM, while Fedora Asahi Remix (Fedora 44, Vulkan 1.4, OpenCL 3.0) supports the M1 Pro.

The disconfirming evidence is strong. Omarchy has already taken the "agentic distro" label. OpenShell now runs on macOS, which weakens the case that a bootable OS is required. Even well-funded corporate efforts have stopped: Intel ended Clear Linux, and Red Hat archived InstructLab. Method note: the shared web search budget ran out after 16 searches. I made up for this with more than 30 direct primary-source reads (official docs, READMEs, build scripts, GitHub API data).

### Recommendations
- Do not build a distro from scratch. Derive AgentOS from a bootc base, either Fedora Atomic or CentOS Stream using Universal Blue's build tooling (which Bluefin GDX already builds for aarch64 and x86_64). Keep the whole OS as one Containerfile so that one part-time maintainer can rebuild it in CI each night.
- Pick one hero differentiator and make it the reason for the project: OS-enforced agent permissions (per-agent identity, Landlock and seccomp, network-namespace egress allowlist, credential proxy, desktop approval prompts). Ship it as a small Rust or Go daemon with declarative policy files, and make it able to wrap existing agents (Claude Code, Codex, OpenClaw, Hermes) unchanged.
- Evaluate NVIDIA OpenShell (Apache-2.0) as the enforcement backend before writing your own. Contribute an OS-native, non-Docker driver instead of competing. Verify its actual kernel mechanisms and turn off its default telemetry in your image.
- Make 'agent undo' the second feature: snapshot the home directory and system subvolumes (btrfs or ostree) at the start of each agent session, and record signed, hash-chained action logs. Document plainly that external side effects cannot be rolled back.
- Build the model layer as a thin scheduler over llama.cpp (Vulkan, CUDA, ROCm) and vLLM, served through one OpenAI-compatible endpoint. Use MemAvailable and swap-based admission control with per-agent token and memory quotas. Reuse Canonical's lscompute-style hardware manifests rather than writing a new hardware probe.
- Treat the M1 Pro as both your development machine and a first-class target: test images in a UTM or QEMU aarch64 VM every day, and boot natively through a Fedora Asahi Remix base for hardware tests. Set the default model size for 16 GB machines at about 4B to 8B parameters, quantised.
- Replace playbooks with continuously tested, declarative 'capability packs' (for example, coding agent plus local model plus egress policy). Run them through CI smoke tests on every image build, because NVIDIA's playbook breakage shows what happens without this.
- Do not compete with Omarchy on desktop polish or with Red Hat on enterprise inference. Position AgentOS as 'the safe place to let agents act': publish a threat model and a red-team test suite. Reproducible agent-escape tests become the main marketing material.
- Limit fine-tuning to sandboxed LoRA jobs on small models, using existing tools (Unsloth or LlamaFactory recipes on Linux; note that MLX does not run on Linux). Defer anything larger until a contributor with suitable hardware joins.
- Reduce abandonment risk from day one: Apache-2.0 licence, a minimal set of custom components, upstream contributions to OpenShell, RamaLama and llama.cpp, and a documented path for users to remove the AgentOS layer and keep a working Fedora system.

### Pain points

**Agents run with full user privileges; 'agentic' distros disable safeguards** — severity 5, breadth 4, willingness to pay 3, confidence high
- Segment: Developers and power users running coding agents and autonomous agents locally
- Description: The only mainstream distro with built-in agents (Omarchy 4) starts every supported coding agent with its permission prompts bypassed. No distro offers OS-level, per-agent controls over filesystem, network, credentials or approvals. Agent runtimes themselves have a very large vulnerability surface.
- Evidence: The omarchy-agent script maps agents to flags such as 'claude --permission-mode auto', 'agy --dangerously-skip-permissions', 'copilot --allow-all', 'crush --yolo', 'grok --permission-mode bypassPermissions' and 'hermes --yolo'. Its comment says these agents 'run unattended'. OpenClaw has 722 published GitHub security advisories (14 critical, 249 high) from January to September 2026.
- Quantified impact: 722 advisories in about 9 months for one agent runtime; Omarchy reports 1.28M ISO downloads exposed to permission-bypass defaults
- Current solutions: NVIDIA OpenShell or NemoClaw (alpha), devcontainers (recommended by Bluefin, but that docs section is empty), manual Podman or firejail, cloud sandboxes (E2B)
- Why still unsolved: Sandboxes (OpenShell, NemoClaw, microsandbox, E2B) are add-ons that need Docker or Podman. They are alpha or early-stage and not built into the boot image, init system or desktop. Distros optimise for convenience, and agent CLIs each define their own permission model.
- Sources: [Omarchy omarchy-agent launcher script](https://github.com/basecamp/omarchy/blob/HEAD/bin/omarchy-agent) (2026-10-16); [OpenClaw security advisories](https://github.com/openclaw/openclaw/security/advisories) (2026-09-11); [NVIDIA NemoClaw README (alpha status)](https://github.com/NVIDIA/NemoClaw) (2026-10-16)

**Unified-memory accounting and model scheduling are not handled by the OS** — severity 4, breadth 4, willingness to pay 3, confidence high
- Segment: Local inference users on unified-memory hardware (DGX Spark, Apple Silicon, AMD Strix Halo)
- Description: On unified-memory machines (DGX Spark GB10, Apple Silicon), standard tools misreport memory. No OS component admits, evicts or shares models across agents according to real memory pressure. Users drop caches by hand or estimate from /proc/meminfo.
- Evidence: The DGX Spark known issues say nvidia-smi shows 'Memory-Usage: Not Supported' and cudaMemGetInfo under-reports allocatable memory. The suggested workaround is to estimate from MemAvailable and SwapFree, or to flush caches with drop_caches. The July 2026 release notes list 'improved out-of-memory handling' as a new item.
- Quantified impact: Not quantified publicly; NVIDIA needed a dedicated OOM fix in July 2026 and a BIOS option reserving 2 or 4 GB for the display
- Current solutions: Manual model sizing, Ollama keep-alive settings, by-hand /proc/meminfo checks, drop_caches
- Why still unsolved: Inference engines (llama.cpp, vLLM, Ollama) each manage only their own process. Canonical inference snaps choose an engine for one model at install time, not at runtime across models. Distros treat model serving as an application, not a system resource.
- Sources: [DGX Spark Known Issues](https://docs.nvidia.com/dgx/dgx-spark/known-issues.html) (2026-09-10); [DGX Spark Release Notes](https://docs.nvidia.com/dgx/dgx-spark/release-notes.html) (2026-07); [inference-snaps engine auto selection doc](https://github.com/canonical/inference-snaps/blob/HEAD/docs/explanations/engine-auto-selection.md) (2026-10-16)

**AI OS images are tied to one hardware vendor and integrated through documentation, not design** — severity 4, breadth 4, willingness to pay 3, confidence high
- Segment: Prosumers and small teams with mixed or non-NVIDIA hardware
- Description: Each vendor image serves one silicon family. Integration of tools is left to playbooks that break as upstream components change.
- Evidence: RHEL AI 3.5 'supports NVIDIA CUDA and AMD ROCm AI accelerators only'. DGX OS targets GB10 and DGX only. The Bluefin GDX build adds only NVIDIA drivers, uv and nvtop. Open DGX Spark playbook issues (August to October 2026) include txt2kg failing on a fresh machine, a gpt-oss-20b download returning 'Entry not found', a wrong vLLM tool-call parser, GB10 throughput dropping 25 to 40 percent over uptime, and lost network reachability.
- Quantified impact: Over 100 playbook directories in NVIDIA's repo; reported 25 to 40 percent sustained throughput loss on GB10 until a full power-off
- Current solutions: Vendor playbooks, NGC containers, Ubuntu LTS as the de facto default because vendor guides target it first
- Why still unsolved: Vendors optimise for their own silicon and enterprise support contracts. Community distros lack the testing hardware. Playbooks are cheaper to write than an integrated, continuously tested system.
- Sources: [RHEL AI 3.5 Getting Started](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux_ai/3.5/html/getting_started/pr01) (2026); [NVIDIA dgx-spark-playbooks issues](https://github.com/NVIDIA/dgx-spark-playbooks/issues) (2026-10-09); [Bluefin LTS GDX package script](https://github.com/ublue-os/bluefin-lts/blob/HEAD/build_scripts/overrides/gdx/30-packages.sh) (2026-10-15)

**Fine-tuning has been removed from or split off AI OS images** — severity 3, breadth 3, willingness to pay 3, confidence high
- Segment: Developers who want to adapt models locally and safely
- Description: The leading enterprise AI OS dropped its built-in tuning workflow and became inference only. Fine-tuning now lives in standalone apps and SDKs with no OS-level resource or safety integration.
- Evidence: InstructLab announced on 2025-09-02 that it would split into sdg_hub and training_hub. The repository is archived (last push 2026-03-30). The RHEL AI 3.0 notes describe an inference platform and say 'There is no direct upgrade path from Red Hat Enterprise Linux AI 1.5'. Meanwhile, standalone tuning tools have large followings: Unsloth has 78,097 stars and LlamaFactory 75,502.
- Quantified impact: InstructLab: 1,418 stars when archived, versus 78k for Unsloth, an app-level tool
- Current solutions: Unsloth, LlamaFactory, Axolotl, torchtune, MLX-LM on macOS
- Why still unsolved: Tuning needs a lot of hardware and support. Vendors moved it to Kubernetes and OpenShift AI. Distros do not package fast-changing training stacks.
- Sources: [InstructLab README community announcement](https://github.com/instructlab/instructlab) (2025-09-02); [RHEL AI 3.0 Release Notes](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux_ai/3.0/html-single/release_notes/index) (2026)

**AI toolchain packaging is fragile across distros** — severity 4, breadth 5, willingness to pay 2, confidence high
- Segment: All Linux AI users, especially AMD and aarch64
- Description: PyTorch, ROCm and CUDA versions change faster than distro release cycles. This causes stalled updates, crashes on aarch64 and unclear support matrices.
- Evidence: Fedora AI/ML SIG, 2026-08-27: PyTorch 2.13 for F45 is stalled because a libstdc++ update made ROCm unbuildable, and ROCm 10 will not ship in F45. The January 2026 meeting reported PyTorch segfaults on aarch64. Lambda Stack's page names no supported Ubuntu release. DGX Spark CUDA is 'verified only at software release', so newer features come only through NGC containers.
- Quantified impact: Fedora 43, 44 and 45 ship PyTorch 2.8.0, 2.9.1 and 2.12.0 respectively, so users face three different stacks
- Current solutions: Containers (NGC, RamaLama images), uv or pip virtual environments, Ubuntu LTS with vendor repos
- Why still unsolved: Licensing rules keep CUDA out of community repos. ROCm has ABI churn. Volunteer packagers are scarce (the SIG asks for more contributors).
- Sources: [Fedora ai-ml-sig meeting 2026-08-27](https://discussion.fedoraproject.org/t/fedora-ai-ml-sig-meeting-2026-08-27/200441) (2026-08-27); [Lambda Stack](https://lambda.ai/lambda-stack-deep-learning-software) (2026)

**One-click AI installers execute arbitrary scripts without real isolation** — severity 4, breadth 3, willingness to pay 2, confidence medium
- Segment: Hobbyists and creators installing local AI apps
- Description: Popular launchers install AI apps by running shell commands from community scripts. They offer only directory-level separation.
- Evidence: Pinokio's README script policy says scripts 'can run anything' and that admins may modify featured scripts after publication. Issues include 'Virus found in the latest review of Pinokio' (2026-01-28) and 'Security: please enable private vulnerability reporting' (2026-09-10). Pinokio has 8,260 stars.
- Quantified impact: Not quantified
- Current solutions: Manual review of scripts, running in VMs, Flatpak (incomplete GPU and CUDA support)
- Why still unsolved: Isolation would break many GPU and Python app installs. Launchers favour compatibility, and the OS provides no capability boundary for installed AI apps.
- Sources: [Pinokio repository](https://github.com/pinokiocomputer/pinokio) (2026-09-02); [Pinokio issues](https://github.com/pinokiocomputer/pinokio/issues) (2026-09-10)

**Risk that corporate AI distros and tools will be abandoned** — severity 3, breadth 3, willingness to pay 2, confidence high
- Segment: Early adopters and enterprises evaluating AI distros
- Description: Users are cautious about committing to AI-specific OS projects because well-funded efforts are often shut down or merged into other products.
- Evidence: Intel ended Clear Linux 'effective immediately' with no further security patches (reported 2025-07-18). InstructLab is archived. System76's Tensorman has 206 stars with little activity. blendOS is read-only (last push 2025-10-21). Ollama Rucksack (a bootable Ollama ISO) has 9 stars and no commits since 2025-06-02.
- Quantified impact: Not quantified
- Current solutions: Staying on Ubuntu, Fedora or Arch, with AI tools installed as apps or containers
- Why still unsolved: Business priorities change. Niche distros lack the maintainers to keep up with AI stack churn.
- Sources: [Phoronix: Intel Ends Clear Linux](https://www.phoronix.com/news/Intel-Ends-Clear-Linux) (2025-07-18); [ollama-rucksack](https://github.com/vedmakk/ollama-rucksack) (2025-06-02)

### Key statistics

| Claim | Value | Source type | Date | Source |
|---|---|---|---|---|
| Omarchy ISO downloads in year one (self-reported) | 1,280,643 | vendor | 2026-10-16 | https://omarchy.org/ |
| Omarchy GitHub stars / forks | 44,921 / 5,493 | primary_research | 2026-10-16 | https://github.com/basecamp/omarchy |
| Pledged to the Omacom Foundation (self-reported) | $23,650,000 | vendor | 2026-10-16 | https://omarchy.org/ |
| OpenClaw published GitHub security advisories, Jan to Sep 2026 | 722 total (14 critical, 249 high, 390 medium, 69 low) | primary_research | 2026-09-11 | https://github.com/openclaw/openclaw/security/advisories |
| Current DGX Spark (Founders Edition) software stack | DGX OS 7.5.0, driver 580.159.03, CUDA 13.0.2, kernel 6.17 | official | 2026-07 | https://docs.nvidia.com/dgx/dgx-spark/release-notes.html |
| NVIDIA DGX Spark playbooks repository | 1,436 stars; over 100 playbook directories, including NemoClaw, OpenShell, Hermes and OpenClaw | primary_research | 2026-10-15 | https://github.com/NVIDIA/dgx-spark-playbooks |
| Reported GB10 sustained throughput loss over uptime | 25-40% degradation, restored only by a full power-off | practitioner | 2026-10-09 | https://github.com/NVIDIA/dgx-spark-playbooks/issues |
| NVIDIA NemoClaw traction and status | 22,754 stars; alpha; Apache-2.0; created 2026-03-15 | primary_research | 2026-10-16 | https://github.com/NVIDIA/NemoClaw |
| NVIDIA OpenShell traction and platforms | 16,375 stars; Linux, macOS on Apple Silicon, WSL2 (experimental); needs Docker, Podman or a VM | primary_research | 2026-10-16 | https://github.com/NVIDIA/OpenShell |
| RHEL AI accelerator support | NVIDIA CUDA and AMD ROCm only | official | 2026 | https://docs.redhat.com/en/documentation/red_hat_enterprise_linux_ai/3.5/html/getting_started/pr01 |
| InstructLab status | Split into sdg_hub and training_hub (announced 2025-09-02); repository archived, last push 2026-03-30 | official | 2025-09-02 | https://github.com/instructlab/instructlab |
| Canonical inference snaps catalogue and community traction | 19 snaps; 52 GitHub stars; public beta launched 2025-10-23 | primary_research | 2026-10-16 | https://github.com/canonical/inference-snaps |
| SLES 16 general availability and agentic features | GA 2025-11-04; MCP host and server in tech preview; 16-year lifecycle | news | 2025-11-04 | https://itwire.com/business-it-news/data/ai-assisted-management-drives-latest-suse-linux-release |
| Fedora PyTorch versions by release | F43: 2.8.0, F44: 2.9.1, F45/Rawhide: 2.12.0; PyTorch 2.13 stalled by ROCm and libstdc++ breakage | official | 2026-08-27 | https://discussion.fedoraproject.org/t/fedora-ai-ml-sig-meeting-2026-08-27/200441 |
| Fedora Asahi Remix base and GPU API conformance (M1 Pro supported) | Fedora 44; OpenGL 4.6, OpenCL 3.0, Vulkan 1.4; no Neural Engine mentioned | official | 2026-10-16 | https://asahilinux.org/fedora/ |
| AIOS (agiresearch) traction and nature | 6,478 stars; userspace Python runtime, not bootable; latest listed version v0.2.2 (2025-03-12) | primary_research | 2026-07-20 | https://github.com/agiresearch/AIOS |
| Pinokio traction | 8,260 stars; 841 forks | primary_research | 2026-09-02 | https://github.com/pinokiocomputer/pinokio |
| Universal Blue image traction | Bazzite 9,184 stars; Bluefin 2,629; Aurora 808; Bluefin LTS 114 | primary_research | 2026-10-16 | https://github.com/ublue-os/bluefin |
| Lambda Stack component versions | PyTorch 2.7.0, TensorFlow 2.19.0, JAX 0.6.0, CUDA toolkit 12.8.93 | vendor | 2026 | https://lambda.ai/lambda-stack-deep-learning-software |
| Pop!_OS 24.04 LTS with COSMIC Epoch 1 | Released December 2025; kernel 6.17; NVIDIA ISO for GTX 16-series and newer; ARM64 images | news | 2025-12-12 | https://itsfoss.com/news/pop-os-24-04-review/ |
| Agent runtime popularity (context for demand) | OpenClaw 392,136 stars; Hermes Agent 254,230; Ollama 183,386 | primary_research | 2026-10-16 | https://github.com/openclaw/openclaw |

### Players

- **NVIDIA DGX OS 7 / DGX Spark software stack** (Vendor appliance OS (Ubuntu-based)): Ubuntu 24.04-based DGX OS 7.5.0 with NVIDIA kernel 6.17, CUDA 13, DGX Dashboard, NVIDIA Sync, a Docker container runtime and NGC. Over 100 playbooks cover vLLM, SGLang, TRT-LLM, Unsloth, NemoClaw and OpenShell. Air-gapped and cloud-init ISOs arrived in April 2026.. Traction: Playbooks repo has 1,436 stars; unit sales not public. Gap: Locked to GB10 and DGX hardware. Integration relies on playbooks that break. Unified-memory reporting is a known issue. User reports describe throughput loss and network drops. Agent safety is an optional playbook, not the default.
- **Red Hat Enterprise Linux AI 3.x / Red Hat AI** (Enterprise bootc AI image): A bootc image on RHEL 9 with Red Hat AI Inference Server (vLLM), the Model Optimization Toolkit, and Hugging Face, ModelCar and OCI model formats. Supports NVIDIA CUDA and AMD ROCm.. Traction: Docs at version 3.5 (2026); InstructLab archived with 1,418 stars. Gap: Inference only since InstructLab was retired. No direct upgrade path from 1.5 to 3.0. Subscription-gated. No agent permission model. No Apple Silicon or Intel GPU support.
- **Canonical inference snaps (Ubuntu)** (OS-level model packaging): Snap per model, with engines chosen by a lscompute hardware probe that filters and ranks engine manifests. 19 models, including Gemma 4, Qwen 3.x and Nemotron 3. Local standard API.. Traction: 52 GitHub stars; launched 2025-10-23 with Intel and Ampere. Gap: Beta. Picks an engine per model at install time only, with no runtime scheduling across models or agents. Low community traction. Confinement and API details are thin in the docs.
- **SUSE AI / SLES 16** (Enterprise Linux with agentic administration): SLES 16 (GA 2025-11-04) with MCP host and server in tech preview for Cockpit and CLI. MCP in Rancher Prime and Multi-Linux Manager. AI Factory with NVIDIA in preview, with GA expected later in 2026.. Traction: Not quantified. Gap: Agents are used to manage infrastructure, not agents as workloads. Kubernetes-centric. Analysts say it trails Red Hat on ecosystem depth. I found no openSUSE AI image.
- **Fedora AI/ML SIG** (Community packaging effort): Packages for PyTorch, ROCm, llama.cpp, Ollama (ROCm and Vulkan backends) and whisper.cpp. Meets every two weeks.. Traction: Small volunteer SIG. Gap: Packaging stalls (ROCm and libstdc++), PyTorch segfaults on aarch64, no CUDA, no AI image or agent features.
- **Omarchy (Omacom / DHH)** (Opinionated 'agentic' Arch + Hyprland distro): ISO install. First-boot choice of default agent. An agent panel. Crash notifications that send the agent to analyse coredumps. Multi-account agent usage tracking (Claude, Codex, Grok, Kimi). Installers for Claude, ChatGPT, Hermes, OpenClaw and others.. Traction: 44,921 stars; 1.28M ISO downloads in year one (self-reported); $23.65M pledged. Gap: Agents launch with permissions bypassed and have no sandbox. No local model scheduling. x86 focus, with Apple Silicon only through a VM. A desktop aesthetic rather than a safety model.
- **Universal Blue Bluefin / Bluefin GDX / Aurora** (Atomic bootc developer workstations): Bluefin ships Goose, Docker Model Runner, Podman AI Lab and Alpaca, plus Homebrew AI CLIs (claude-code, codex, gemini-cli). GDX is a CentOS Stream bootc image with NVIDIA open drivers, the container toolkit, uv and nvtop, built for x86_64 and aarch64.. Traction: Bluefin 2,629 stars; Bluefin LTS 114; Bazzite 9,184. Gap: AI support is a set of tools added on top. The docs section on agent isolation is empty. GDX is NVIDIA-only and thin. There is no permission model.
- **Pop!_OS 24.04 with COSMIC (System76)** (Desktop distro with NVIDIA ISO): Ubuntu 24.04 base, COSMIC Epoch 1 desktop written in Rust, preinstalled NVIDIA driver 580, ARM64 images.. Traction: cosmic-epoch 6,922 stars. Gap: No AI or agent features. Tensorman (TensorFlow containers) has little activity. Reviewers report audio and screen-recording bugs.
- **Lambda Stack** (APT meta-package for deep learning): One-command install of drivers, CUDA 12.8, PyTorch 2.7, TensorFlow 2.19, JAX 0.6, NCCL and the container toolkit.. Traction: Claims 'thousands' of customers without figures. Gap: NVIDIA and Ubuntu only. No stated Ubuntu version. Training-framework era, with nothing for agents or local LLM serving.
- **Pinokio** (One-click AI app launcher): JSON and shell scripts that install and run AI apps under ~/pinokio, with a Discover catalogue.. Traction: 8,260 stars. Gap: Scripts 'can run anything'. Separation is only by directory. A virus report and a vulnerability-reporting request are open. Not an OS.
- **NVIDIA OpenShell / NemoClaw** (Agent sandbox runtime and reference stack): Policy-enforced sandboxes covering files, syscalls and network egress. Credential injection only for approved endpoints. Formal verification of policy changes. Container or microVM drivers. GPU passthrough through CDI.. Traction: OpenShell 16,375 stars; NemoClaw 22,754. Gap: Alpha. Needs Docker or Podman. Not integrated with the OS. Isolation mechanisms are not specified in the docs. Telemetry on by default. The VM driver supports one GPU.
- **AIOS (agiresearch)** (Research 'LLM agent OS' runtime): A Python kernel abstraction for LLM, memory, storage and tool management, plus an agent SDK.. Traction: 6,478 stars. Gap: Userspace only, not bootable. The Rust scheduler is a placeholder. Last listed release is March 2025.
- **RamaLama (containers project)** (Containerised local model serving): Detects GPUs (CUDA, ROCm, Vulkan, Intel, Asahi). Runs rootless with read-only model mounts. 'ramalama run' has no network by default.. Traction: 3,087 stars. Gap: Alpha. 'ramalama serve' publishes on 0.0.0.0. MLX bypasses containers. No scheduling across models.
- **openEuler Intelligence (EulerCopilot)** (Chinese OS-integrated LLM platform): RAG, workflow orchestration and an MCP preview on openEuler.. Traction: 32 stars on Gitee (stale mirror). Gap: Migrated to AtomGit. Integration into OS releases is undocumented on the page I read.
- **Fedora Asahi Remix** (Apple Silicon Linux base): Fedora 44 on M1 and M2 (M3 listed in device tables). Vulkan 1.4 and OpenCL 3.0 GPU drivers.. Traction: asahi-installer 952 stars. Gap: No Neural Engine support. No AI integration. Hardware support lags newer chips.
- **Ollama Rucksack and similar single-purpose ISOs** (Community bootable LLM ISO): Debian Live ISO with Ollama, XFCE, optionally bundled models and network drivers blacklisted.. Traction: 9 stars. Gap: Abandoned, with no commits since June 2025. No GPU documentation. A demo rather than an OS.

### Trends
- Agents are moving into the distro layer, but for convenience rather than safety. Omarchy 4 adds first-boot agent selection, coredump-to-agent diagnosis and an agent panel. SLES 16 adds an MCP host and server. Bluefin ships Goose.
- Image-based (bootc or ostree) operating systems are now the standard base for AI images. Examples are RHEL AI 3.x, Bluefin LTS and GDX on CentOS bootc, and Fedora Atomic derivatives. These give atomic updates and rollback that an agent 'undo' feature could build on.
- Hardware-aware model packaging is spreading: Canonical lscompute-ranked engines, RamaLama GPU detection, and AMD Lemonade (5,912 stars). All of it works per model rather than as system-wide scheduling.
- Agent sandboxing is forming as its own layer, outside any OS. Examples are NVIDIA OpenShell and NemoClaw (since February to March 2026), microsandbox (8,706 stars), Rivet agentOS (Wasm, 4,792 stars), E2B, gVisor and Firecracker.
- Enterprise vendors are retreating to inference only and Kubernetes-scale serving (RHEL AI 3, llm-d with 4,837 stars, SUSE AI Factory). Fine-tuning has moved into separate SDKs (training_hub, sdg_hub) and apps such as Unsloth and LlamaFactory.
- Vendors increasingly integrate through playbooks (NVIDIA's 100+), which creates ongoing breakage and opens room for continuously tested, declarative integration.
- Agent runtimes have exploded in popularity and attack surface: OpenClaw has 392k stars and 722 advisories in 2026. This pushes demand for containment that does not rely on each agent's own permission prompts.
- Corporate AI distro projects are consolidating or dying (Clear Linux ended in July 2025, InstructLab archived in March 2026). Community-led opinionated distros backed by large foundations are gaining momentum (Omarchy).

### Opportunities

- **Agents as OS principals: an 'agentd' broker that gives each agent its own identity (uid or systemd slice), Landlock and seccomp filesystem and syscall policy, network-namespace egress allowlists, a credential proxy, and approval prompts handled by the OS rather than by each agent CLI.** Why now: Agent CLIs are being run in 'yolo' mode by default on popular distros. OpenClaw's advisory volume shows that agent-level permissions cannot be trusted. OpenShell shows the policy model works but remains an alpha add-on. Evidence: Omarchy omarchy-agent bypass flags; 722 OpenClaw advisories; OpenShell needs Docker or Podman and does not document its kernel mechanisms; Bluefin's devcontainer isolation docs section is empty. Risk: NVIDIA could ship OpenShell by default in DGX OS, or Omarchy could add sandboxing quickly. Fine-grained policy may annoy users into switching it off, as happened with SELinux.
- **A memory-aware local model scheduler that runs as a system service. It admits, evicts and shares models across agents based on real unified-memory pressure, and exposes one OpenAI-compatible endpoint with per-agent quotas.** Why now: DGX Spark documentation admits that memory reporting breaks on unified memory. Inference snaps and RamaLama choose engines but do not schedule. A 16 GB M1 Pro forces strict admission control. Evidence: DGX Spark known issues (cudaMemGetInfo, nvidia-smi); the July 2026 OOM fix; inference-snaps engine auto-selection runs at install time only. Risk: Ollama or llama.cpp could add a multi-model scheduler upstream, which would remove the OS-level advantage.
- **Verified, reversible agent execution: signed, hash-chained logs of agent actions tied to bootc, ostree or btrfs snapshots, giving a one-command 'undo what the agent did'.** Why now: Image-based operating systems already provide atomic rollback. No product links rollback to agent sessions. Evidence: The RHEL AI and Bluefin bootc bases have rollback. SLES 16 promotes 'Instant Rollback' for administrators, not for agents. No surveyed project offers agent action attestation. Risk: Rollback covers system state only; it cannot undo external side effects such as sent emails or API calls. That limit must be stated honestly.
- **aarch64-first image that boots natively on Apple Silicon (through an Asahi base) and on GB10 or Arm servers, with the same agent and model layer on x86_64 with NVIDIA.** Why now: Omarchy supports Apple Silicon only in a VM. Bluefin GDX already builds for aarch64 and the NVIDIA SBSA repository. Asahi supports the M1 Pro with Vulkan 1.4, and RamaLama supports the Asahi GPU. Evidence: omarchy.org VM-only note; Bluefin LTS build-gdx aarch64 overrides; asahilinux.org/fedora conformance list. Risk: No Neural Engine or CUDA on Asahi. Vulkan llama.cpp performance on an M1 Pro is far below macOS MLX. Asahi lags on M3 and newer chips.
- **Safe fine-tuning jobs as OS-managed workloads: small-model LoRA or QLoRA runs inside sandboxes with resource quotas, using recipes from Unsloth or LlamaFactory and dataset provenance.** Why now: RHEL AI dropped InstructLab, which leaves no OS-integrated tuning path, while tuning apps have more than 75k stars each. Evidence: InstructLab archive notice; Unsloth 78,097 stars; LlamaFactory 75,502 stars. Risk: 16 GB of unified memory limits tuning to roughly 1B to 8B parameter models. Tuning stacks change quickly, which makes packaging expensive for a solo maintainer.

### Contrarian signals
- The 'agentic distro' label is already taken. Omarchy calls itself agentic, reports 1.28M ISO downloads in year one, has 44.9k stars and has $23.65M pledged. Being a novel AI-native distro is not a differentiator.
- Agent safety is becoming portable and OS-independent. OpenShell supports Linux, macOS on Apple Silicon and WSL2, and NemoClaw is tested on Ubuntu, macOS and WSL. A bootable OS may not be needed to solve the main safety problem.
- Practitioner guidance says distro choice matters less than GPU, VRAM and drivers, and recommends plain Ubuntu (OpenSourceFeed, 2026-07-09). Most users get local AI through cross-platform apps (Ollama 183k stars, LM Studio).
- Canonical's inference snaps, backed by a major vendor, have only 52 GitHub stars a year after launch. This suggests weak demand for packaging models at the OS level.
- Well-funded corporate efforts in this area have been shut down or broken up (Intel Clear Linux ended in 2025, InstructLab archived in 2026). This is a warning about sustainability for a solo, part-time project.
- Enterprise leaders moved away from local fine-tuning in the OS towards inference only and Kubernetes serving. This suggests tuning on a 16 GB device is a niche feature.
- Asahi Linux on an M1 Pro has no Neural Engine support and no MLX. Local inference performance on a natively booted AgentOS will be well below macOS on the same machine, which weakens the 'fast' promise on the founder's own hardware.
- One third-party RHEL AI 3.4 guide (lucaberton.com) still lists InstructLab as a core component, which contradicts Red Hat's docs and the archived repository. Secondary coverage in this space is often unreliable.

### Open questions
- Clarifying question for the founder: who is the primary user? Is it a developer running coding agents on a laptop, a home-lab owner with an NVIDIA GPU or DGX Spark, or an organisation that needs auditable agents? Each points to a different base image and threat model.
- Clarifying question for the founder: what does 'revolutionary' mean in measurable terms? Is it agent safety (escape-test pass rate), speed (tokens per second, time to first agent) or zero-configuration setup? The current brief asks for 'safely and fast' without ranking the two, and they conflict on 16 GB hardware.
- Clarifying question for the founder: must the first release boot natively on Apple Silicon, or is a VM on the M1 Pro plus native x86_64 with NVIDIA acceptable? Native Asahi support adds kernel and GPU risk.
- Method limitation: the shared web search budget ran out after 16 of the required 20 searches. I made up for this with more than 30 direct primary-source reads, but some areas remain unverified: openSUSE AI images, Ubuntu 26.04 AI features, Warmwind and other closed 'AI OS' products, and Chinese AI distros (openKylin, deepin UOS AI).
- Which kernel mechanisms does OpenShell actually use (Landlock, seccomp, netns), and is a driver without Docker feasible? The public docs I read do not say.
- What are the real adoption numbers for DGX OS (Spark unit sales), Pop!_OS installs and inference snap installs? None is published.
- What is the release date of RHEL AI 3.0, and does Red Hat's Training Hub return to RHEL AI in a later 3.x release? Official release notes pages are thin.
- Can Vulkan llama.cpp on Asahi with 16 GB deliver usable speed for 7B to 8B models compared with MLX on macOS? Independent benchmarks are needed before committing to native Apple Silicon support.
- Will Omarchy add sandboxing in response to criticism of its permission-bypass defaults? Monitor its issue tracker and AGENTS.md.

### Adversarial verification

- **DISPUTED** — No product yet treats agents as first-class OS principals; no distro offers OS-level, per-agent controls over filesystem, network, credentials or approvals (executive summary and the first pain point).
  - Corrected: No mainstream Linux distribution was found that ships per-agent OS identities by default. Microsoft Windows 11, however, has offered agent accounts and an agent workspace in preview since October to December 2025. Each agent gets its own account and a separate Windows session, and Windows asks for consent per agent and per known folder. The claim holds only if it is limited to Linux, and it should be worded that way.
  - Evidence: The Microsoft support page 'Experimental agentic features' says each agent operates using its own account, distinct from the user's personal account, and that the agent workspace runs in a separate Windows session. The default for known folders (Documents, Downloads, Desktop and so on) is to ask for consent, with options 'Allow Always', 'Ask every time' and 'Never allow'. The page was first published on 2025-10-16 and updated on 2025-11-17 and 2025-12-05 (preview builds 26100.7344+ and 26220.7262+). The feature is off by default and is still in preview. NVIDIA OpenShell also now ships an MXC compute driver that runs agents in Windows AppContainer or isolation sessions (Windows 11 Insider build 26300.8553 or later). I could not run a search for Linux counterexamples because the web search budget was used up.
  - Sources: [Microsoft Support: Experimental agentic features](https://support.microsoft.com/en-us/windows/experimental-agentic-features-a25ede8a-e4c2-4841-85a8-44839191dfb3) (2025-12-05); [OpenShell MXC driver README](https://github.com/NVIDIA/OpenShell/blob/main/crates/openshell-driver-mxc/README.md) (2026-10-16)
- **PARTIALLY_CONFIRMED** — Omarchy has 44,921 stars, 1.28M year-one ISO downloads and $23.65M pledged, calls itself 'agentic', and its launcher starts every coding agent with its permission prompts switched off.
  - Corrected: The adoption figures and the 'agentic' label are correct. However, omarchy-agent does not switch off permissions for every agent. Of its 15 agent entries, 11 pass skip or auto-approve flags (opencode --auto, agy --dangerously-skip-permissions, copilot --allow-all, crush --yolo, grok --permission-mode bypassPermissions, codex --approve-for-me, cursor-agent --yolo --trust, hermes --yolo, muse --approval-mode never, omp --auto-approve, kimi --yolo). Claude is launched in auto mode, which is classifier-reviewed and is Claude Code's own built-in default, not a bypass. OpenClaw, Pi and Ori get no skip flag. Muse keeps its own sandbox.
  - Evidence: The GitHub API returns 44,921 stars, 5,493 forks, MIT licence and repository creation on 2025-06-01. The default branch is 'quattro', not master. omarchy.org says 'Beautiful, fun & agentic Linux', offers version 4.0.4, and reports 1,280,643 year-one ISO downloads and $23,650,000 pledged. It offers Apple Silicon only through a VM. The raw script on the quattro branch shows 'claude --permission-mode auto'. Its comment states that agents launched from the keybinding or menu run unattended. Claude Code documentation says that in auto mode a classifier reviews actions instead of the user, and that auto mode is the built-in starting mode in v2.1.283 or later. bypassPermissions is a separate mode, intended for isolated containers and VMs.
  - Sources: [omarchy-agent (quattro branch)](https://raw.githubusercontent.com/basecamp/omarchy/quattro/bin/omarchy-agent) (2026-10-16); [Omarchy homepage](https://omarchy.org/) (2026-10-16); [GitHub API: basecamp/omarchy](https://api.github.com/repos/basecamp/omarchy) (2026-10-16); [Claude Code permission modes](https://code.claude.com/docs/en/permission-modes) (2026-10-16)
- **CONFIRMED** — OpenClaw has 722 published GitHub security advisories in 2026 (14 critical, 249 high, 390 medium, 69 low), 'January to September', 'about 9 months'; 392k stars.
  - Corrected: Confirmed exactly: 722 published advisories (14 critical, 249 high, 390 medium, 69 low) and 392,136 stars. The advisories run from 2026-01-31 to 2026-09-11, about 7.4 months rather than 9. They are heavily front-loaded: 424 of 722 (59 percent) were published in February and March, none in July or August, and 75 in September. The count reflects a disclosure wave after launch, not a steady monthly rate.
  - Evidence: I paginated the GitHub REST API (security-advisories, state=published, cursor pagination) and found 722 unique GHSA IDs. Severity counts: medium 390, high 249, low 69, critical 14. Advisories by month: January 3, February 220, March 204, April 108, May 67, June 45, September 75. The web UI shows 73 pages, with only 2 entries on the last page. The repository was created on 2025-11-24.
  - Sources: [OpenClaw security advisories](https://github.com/openclaw/openclaw/security/advisories) (2026-09-11); [GitHub API: openclaw/openclaw](https://api.github.com/repos/openclaw/openclaw) (2026-10-16)
- **REFUTED** — NVIDIA OpenShell is alpha, needs Docker or Podman, does not specify its isolation mechanisms in the docs, and has telemetry on by default; the founder should contribute an OS-native, non-Docker driver.
  - Corrected: OpenShell is no longer labelled alpha. It is on a 0.1.x 'stable release cadence', with v0.1.4 released on 2026-10-13. NemoClaw is still alpha. OpenShell documents its kernel mechanisms in detail: Landlock LSM (ABI 3, Linux 6.2 or later, required), seccomp BPF and seccomp user notification, a dedicated network namespace with a veth pair, and privilege drop through setuid and setgid. It does not need Docker. A libkrun MicroVM driver runs on KVM (Linux) and Hypervisor.framework (macOS), alongside Docker, Podman, Kubernetes and a Windows MXC driver. Telemetry is on by default, with an opt-out. A non-Docker driver therefore already exists, so the recommendation should change from 'contribute one' to 'evaluate the MicroVM driver and a host-native driver'.
  - Evidence: The README says 'New in OpenShell 0.1.x: a stable release cadence, new isolation primitives' and lists the requirement as 'Docker, Podman, or host virtualization'. docs/about/support-matrix.mdx lists Landlock and seccomp as Required and marks MicroVM as 'Supported for VM-backed sandboxes'. It also says that on macOS these kernel modules run inside the Docker Desktop Linux VM. docs/security/best-practices.mdx names the Landlock LSM, Seccomp BPF and 'Network Namespace Isolation'. The source contains crates/openshell-sandbox/src/linux/landlock.rs and seccomp.rs. The VM driver README is labelled 'Status: Experimental', which does not match the support matrix. The GitHub API shows 16,376 stars and Apache-2.0.
  - Sources: [OpenShell README](https://github.com/NVIDIA/OpenShell) (2026-10-16); [OpenShell Support Matrix (v0.1.4)](https://docs.nvidia.com/openshell/latest/about/support-matrix) (2026-10-13); [OpenShell security best practices (source)](https://github.com/NVIDIA/OpenShell/blob/main/docs/security/best-practices.mdx) (2026-10-16); [OpenShell releases](https://github.com/NVIDIA/OpenShell/releases) (2026-10-13); [OpenShell VM driver README](https://github.com/NVIDIA/OpenShell/blob/main/crates/openshell-driver-vm/README.md) (2026-10-16)
- **CONFIRMED** — DGX Spark runs DGX OS 7.5.0 (kernel 6.17, driver 580.159.03, CUDA 13.0.2). NVIDIA admits that nvidia-smi shows memory as 'Not Supported' and that cudaMemGetInfo under-reports. The workaround is MemAvailable/SwapFree estimation or drop_caches. The July 2026 release improved OOM handling.
  - Corrected: Confirmed, with framing caveats. NVIDIA describes the nvidia-smi 'Memory-Usage: Not Supported' output as expected behaviour on iGPU platforms, not as a defect. cudaMemGetInfo omits memory that could be reclaimed from swap. NVIDIA labels drop_caches as a workaround for debugging only. The July 2026 release adds better OOM handling and a BIOS option that reserves 2 GB (default) or 4 GB for the display. No release after July 2026 is listed. The software versions apply to the Founders Edition, and partner GB10 systems may lag.
  - Evidence: The Known Issues page was last updated on Sep 10, 2026. It contains sample C code that parses MemAvailable and SwapFree from /proc/meminfo, and the code treats free swap as zero when huge pages are configured. The Release Notes list releases from November 2025 to July 2026, identify 7.5.0 as current, and give 580.159.03, CUDA 13.0.2 and Canonical kernel 6.17 in the version table. The table is not explicitly tied to the July release.
  - Sources: [DGX Spark Known Issues](https://docs.nvidia.com/dgx/dgx-spark/known-issues.html) (2026-09-10); [DGX Spark Release Notes](https://docs.nvidia.com/dgx/dgx-spark/release-notes.html) (2026-07)

**Missed counter-evidence**
- Windows 11 agent accounts and agent workspace (preview since October to December 2025) already give each agent its own OS account, a separate session and per-folder consent. A major OS vendor has therefore already started the 'agents as OS principals' idea, which weakens the novelty of the headline.
- Agent CLIs are building OS-enforced sandboxing into themselves. Claude Code's built-in Bash sandbox uses Seatbelt on macOS and bubblewrap plus socat (optional seccomp) on Linux and WSL2, with a domain-allowlist proxy, and classifier-reviewed auto mode became its default in v2.1.283. This competes directly with the proposed 'agentd' broker.
- OpenShell already provides the non-Docker path the brief recommends building: a libkrun MicroVM driver (KVM and Hypervisor.framework), a Kubernetes driver with user namespaces, and a Microsoft MXC/AppContainer driver for Windows 11 Insider builds. It also ships a formal policy prover. Four tagged releases (v0.1.1 to v0.1.4) came out between 2026-09-26 and 2026-10-13.
- On macOS, OpenShell enforces Landlock and seccomp inside the Docker Desktop Linux VM. A Linux kernel boundary is still required, but delivering a Linux VM image may be enough. This supports a 'Linux image' strategy more than a natively booted OS.
- OpenClaw advisories peaked in February and March 2026 (59 percent of all), with none published in July or August. Using '722 advisories' as evidence of an ongoing risk rate overstates the current trend.
- Fedora Asahi Remix officially supports only M1 and M2 series devices. M3 appears only in feature tables, and M4 is not mentioned. The brief's 'M3 listed in device tables' could be misread as support.
- Inconsistency inside OpenShell: the VM driver README says 'Experimental', while the support matrix says 'Supported'. The founder should treat VM-backed isolation as experimental.
- The founder should answer two questions before relying on this brief. Is the target strictly Linux, or does Windows agent-workspace competition matter? Would an OpenShell-based Linux VM image meet the 'bootable OS' goal?

**Reliability note:** Executive summary: of the five claims checked, two are confirmed (OpenClaw advisories, DGX Spark memory issues), one is partially confirmed (Omarchy), one is disputed (the 'no OS treats agents as principals' headline) and one is refuted (OpenShell maturity and documentation). The brief's numbers are highly accurate: GitHub API values for Omarchy, OpenClaw, NemoClaw, OpenShell, inference-snaps and InstructLab matched exactly or within one star. The quoted wording from RHEL AI 3.5 ('NVIDIA CUDA and AMD ROCm AI accelerators only', a bootc image) and from the DGX Spark docs was verified. The weaknesses are qualitative. The brief overstates Omarchy's permission bypass by counting Claude's default auto mode as a bypass. Its description of OpenShell is out of date: the project is not alpha, does document Landlock, seccomp and netns, and does not need Docker. It also misses Windows 11 agent accounts, which undercut the headline. Recommendations that depend on OpenShell (contribute a non-Docker driver, verify its kernel mechanisms) should be revised. Limitation: the shared web search budget was used up after one attempted search, so all checks used direct primary-source fetches and the GitHub API. Linux-specific counterexamples to the headline, Ubuntu 26.04 AI features and newer Omarchy sandboxing work were not searched. A follow-up message would allow those searches to continue.


---

## Image-based Linux build tooling for a bootable AI OS maintained by one person
_Workflow label: `bootable_tooling`_

### Executive summary
Recommendation: build AgentOS as a bootc (bootable OCI) image based on CentOS Stream 10, set up the way Universal Blue builds its images: a Containerfile, GitHub Actions, cosign signing, and prebuilt signed NVIDIA open kernel modules. Ship heavy AI stacks as pinned OCI containers. Do not bake them into /usr. This is the only base that meets all four needs together. A tiny team already ships at scale on it. NVIDIA and CUDA already work on an immutable system with Secure Boot. x86_64 and aarch64 images can both be built for free in CI. The required skill is writing Containerfiles, not learning a new language.

Evidence: Bazzite, built the Universal Blue way, passed 100,000 weekly active devices (15 Oct 2026). Bluefin GDX already ships CentOS Stream 10 with NVIDIA drivers, CUDA, RamaLama and uv, on x86_64 and aarch64, with Secure Boot (Sep 2025). bootc release 1.17.0 (6 Oct) declared its composefs backend stable, and its CLI and API are stable. It is still only a CNCF Sandbox project. GitHub gives public repositories free x64 and arm64 runners (4 vCPU, 16 GB RAM, 14 GB SSD). AMD lists RHEL 10.x as supported for ROCm, but not Fedora or NixOS, which favours a RHEL-family base.

NixOS is the strongest rival on reproducibility but the weakest fit for a solo AI OS. The default cache.nixos.org does not hold CUDA builds. The community CUDA cache moved and has no coverage for stable releases (Dec 2025). Flox only caches builds from its own fork that follows the unstable branch. nixified-ai deprecated most of its packages for lack of funding. lanzaboote, which gives NixOS Secure Boot, is not yet part of nixpkgs.

Other options are a poor fit:
- Ubuntu Core ties the project to the Snap Store.
- Talos and Flatcar are Kubernetes and server appliances.
- mkosi and ParticleOS follow the ideal verified-boot design, but ParticleOS says it is unstable, and its prebuilt images are x86-64 only.

Heavy stacks: vLLM ships binaries compiled for CUDA 12.9 that do not work with other CUDA or PyTorch builds. Baking it into the base OS therefore guarantees conflicts. Use one OCI image per accelerator, as RamaLama does. Optionally use bootc logically bound images for one default runtime. Keep sysexts for small host tools only.

Critical constraint: the M1 Pro cannot test NVIDIA, x86 Secure Boot or TPM measured boot. Cross-architecture builds under emulation run 8 to 10 times slower. Build in CI, and buy one x86_64 test machine with an NVIDIA GPU. Sealed UKI and composefs images were test-only in April 2026, so postpone measured boot.

Method note: the shared web-search budget ran out after 6 searches. These findings rest on about 40 direct reads of primary sources.

### Recommendations
- Base: use a bootc image derived from quay.io/centos-bootc/centos-bootc:stream10, modelled on Bluefin LTS and GDX. This matches the RHEL 10 support listed by NVIDIA and ROCm and provides a kernel that stays put for the release lifetime. Add a Fedora-based edge variant only if users need GPUs that the 6.12 kernel cannot drive.
- Start from ublue-os/image-template: Containerfile, build.yml, build-disk.yml, cosign key in SIGNING_SECRET, Renovate digest pinning. Publish to GHCR from a public repository so the free x64 and arm64 runners apply.
- Ship only NVIDIA open kernel modules (Turing and newer) as prebuilt signed packages copied in at build time from ublue-os/akmods, or from your own akmods job using your own MOK key. Document the one-time MOK enrollment step. Formally drop support for Maxwell, Pascal and Volta.
- Keep exactly two images and two architectures at launch: agentos (generic, Mesa/ROCm-ready) and agentos-nvidia, each built for x86_64 and aarch64. Universal Blue's 2025 scope cuts show the cost of each extra variant.
- Never install vLLM, SGLang, PyTorch or CUDA userspace into /usr. Ship them as digest-pinned OCI images run with Podman and the NVIDIA Container Toolkit (CDI). Bind at most one default runtime with bootc logically bound images, and keep model weights in /var.
- Use sysexts only for small optional host tools, with SYSEXT_LEVEL set and SHA256SUMS published for systemd-sysupdate, following the Flatcar sysext-bakery pattern.
- Build every architecture natively in CI. Avoid emulated cross-architecture builds on the M1 Pro, which run 8-10x slower. Use podman machine on the Mac only for quick aarch64 builds and qcow2 boot checks in a VM.
- Buy or rent one x86_64 test machine with a Turing-or-newer NVIDIA GPU and TPM 2.0, and run an install, upgrade and rollback test on it before every release.
- Postpone sealed UKI and TPM measured boot until Fedora and bootc sealed images leave test status. In the meantime, choose a systemd-boot-compatible ESP layout so you can switch later.
- Use Nix or Flox only inside containers or developer shells for reproducible agent toolchains. Do not use NixOS as the OS base unless the founder already writes Nix fluently.

### Pain points

**NVIDIA kernel modules and Secure Boot on immutable images** — severity 5, breadth 4, willingness to pay 3, confidence high
- Segment: Maintainers of image-based Linux distributions aimed at GPU and AI users
- Description: Immutable images cannot compile NVIDIA modules on the host at runtime. The modules have to be built and signed at image build time, matched exactly to the image kernel, and their signing key has to be enrolled as a MOK on every machine. When the kernel and driver versions drift apart, systems fail to boot or lose GPU acceleration.
- Evidence: Universal Blue runs a separate akmods pipeline that publishes prebuilt modules daily, grouped as nvidia, nvidia-open and zfs, for six kernel flavours. It sometimes pins the CoreOS kernel when ZFS modules cannot build against it. On Flatcar, building the driver from source takes 5-10 minutes at first boot and cannot work with Secure Boot; only the prebuilt signed sysext can. Fedora users report that a new kernel blocked NVIDIA loading (Jun 2025) and that driver 610 broke package dependencies (May 2026). Talos says its NVIDIA drivers are tied to a specific Talos release.
- Quantified impact: First-boot compile adds 5-10 minutes and makes Secure Boot impossible (Flatcar docs). Every kernel bump needs matching rebuilt modules.
- Current solutions: ublue-os/akmods prebuilt signed modules copied in with COPY --from; Flatcar prebuilt sysext signed with a key thrown away after the build; Talos nonfree-kmod-nvidia extensions; negativo17 akmod, DKMS and RHEL kABI modules; NVIDIA open modules by default from R560.
- Why still unsolved: NVIDIA kernel modules are out-of-tree. The proprietary flavour cannot ship in-tree, and nova_core is not production-ready. Every distribution has to run its own signing and key-enrolment process.
- Sources: [ublue-os/akmods README](https://github.com/ublue-os/akmods); [Flatcar: Using NVIDIA GPUs](https://www.flatcar.org/docs/latest/os-config/host-config/using-nvidia/); [negativo17 NVIDIA driver page and comments](https://negativo17.org/nvidia-driver/) (2026-06); [Talos NVIDIA proprietary drivers](https://docs.siderolabs.com/talos/v1.11/configure-your-talos-cluster/hardware-and-drivers/nvidia-gpu-proprietary)

**Dependency conflicts in heavy AI stacks (CUDA, PyTorch, vLLM, SGLang)** — severity 5, breadth 5, willingness to pay 3, confidence high
- Segment: Users and builders of local LLM serving and fine-tuning
- Description: Each inference engine pins its own versions of CUDA, PyTorch and Python. Installing two engines into one system /usr, or onto the host's Python, breaks one of them. On an image-based OS a host-level conflict also blocks the whole image build.
- Evidence: The vLLM docs say its precompiled binaries target CUDA 12.9 (with 12.8 and 13.0 builds) and are incompatible with other CUDA or PyTorch builds. They recommend uv with a fresh environment and warn against conda-installed PyTorch. The ROCm wheels require Python 3.12.
- Quantified impact: Each engine needs its own environment. Reusing an existing PyTorch install requires building vLLM from source (vLLM docs).
- Current solutions: Per-engine OCI images (vllm/vllm-openai; RamaLama images for CUDA, ROCm, Vulkan and Intel), uv virtual environments, Nix closures and Flox environments.
- Why still unsolved: Upstream projects compile against specific CUDA and PyTorch ABIs and release fast. No distribution can keep one shared set of versions coherent.
- Sources: [vLLM GPU installation](https://docs.vllm.ai/en/latest/getting_started/installation/gpu.html) (2026); [RamaLama README](https://github.com/containers/ramalama)

**Developing on Apple Silicon cannot validate the target hardware** — severity 4, breadth 3, willingness to pay 2, confidence medium
- Segment: Solo founders and indie maintainers who use Macs
- Description: An M1 Pro can build and boot aarch64 Linux in a VM. It cannot test NVIDIA or CUDA, x86 UEFI Secure Boot, or TPM 2.0 measured boot, and x86_64 builds run under slow emulation. The main bootc test tools assume a Linux host with libvirt or QEMU.
- Evidence: The image-builder README marks cross-architecture builds as experimental and 8-10 times slower than native builds. The bcvk article requires a RHEL 10 host with Podman and libvirt. podman-bootc was archived on 3 Jun 2026 in favour of bcvk. Fedora Asahi Remix 44 supports the M1 Pro, but Asahi has no NVIDIA support.
- Quantified impact: Emulated cross-architecture builds run 8-10x slower. GPU and Secure Boot paths cannot be tested on the Mac at all.
- Current solutions: podman machine (applehv) for aarch64 builds; free GitHub arm64 and x64 runners; a separate x86 test machine or a rented cloud GPU.
- Why still unsolved: Apple Silicon does not offer standard UEFI Secure Boot or TPM for Linux, and it cannot host NVIDIA GPUs. The tooling vendors target Linux workstations.
- Sources: [osbuild/image-builder README](https://github.com/osbuild/image-builder); [Test and deploy bootable containers with bcvk](https://developers.redhat.com/articles/2026/10/13/test-deploy-bootable-containers-with-bcvk) (2026-10-13); [podman-bootc (archived)](https://github.com/containers/podman-bootc) (2026-06-03); [Fedora Asahi Remix](https://asahilinux.org/fedora/)

**CUDA binary caching and licensing on NixOS** — severity 4, breadth 3, willingness to pay 3, confidence high
- Segment: NixOS ML users and anyone considering NixOS as an AI OS base
- Description: CUDA is unfree, so cache.nixos.org does not hold CUDA builds. Users must trust third-party caches whose coverage depends on which nixpkgs revision they pin. A cache miss means compiling PyTorch or CUDA packages locally for hours.
- Evidence: The NixOS wiki says the Foundation does not cache CUDA packages and that large builds can freeze weaker machines. On Discourse (Dec 2025), maintainers said the cache moved to cache.nixos-cuda.org and that stable branches are not cached. Flox builds from its own fork that follows the unstable branch, so pinned channel revisions are unlikely to get cache hits. nixified-ai deprecated its InvokeAI and textgen packages for lack of funding.
- Quantified impact: No cache for stable releases (Dec 2025), so systems pinned to stable rebuild CUDA software locally.
- Current solutions: cache.nixos-cuda.org; the Flox CUDA catalog, which has an NVIDIA redistribution agreement; following nixos-unstable.
- Why still unsolved: Funding depends on volunteer maintainers, and NVIDIA's licence limits redistribution to approved parties.
- Sources: [NixOS Wiki: CUDA](https://wiki.nixos.org/wiki/CUDA) (2026-07-19); [CUDA Cache for Nix Community (thread)](https://discourse.nixos.org/t/cuda-cache-for-nix-community/56038) (2025-12-02); [Flox CUDA catalog and cache](https://flox.dev/docs/concepts/cuda-catalog-and-cache); [nixified-ai flake](https://github.com/nixified-ai/flake)

**Churn in image-build tooling** — severity 3, breadth 4, willingness to pay 1, confidence high
- Segment: Small teams building on bootc, OSTree or Nix image tooling
- Description: Tools for building and testing bootable images keep being renamed, moved or archived. Composefs and sealed images are moving from experimental to stable at different speeds in different projects, so documentation contradicts itself.
- Evidence: podman-bootc archived on 3 Jun 2026. bootc-image-builder source moved into osbuild/image-builder. nixos-generators archived on 30 Jan 2026. The OSTree composefs docs say 'do not deploy to production systems yet', while bootc 1.17.0 calls its composefs backend stable.
- Quantified impact: Three major tools archived or relocated within 12 months.
- Current solutions: Pin tool versions, follow the Universal Blue templates, and let Renovate update digests.
- Why still unsolved: The ecosystem is converging fast, and maintainers are consolidating tools.
- Sources: [bootc v1.17.0 release](https://github.com/bootc-dev/bootc/releases/tag/v1.17.0) (2026-10-06); [OSTree composefs](https://ostreedev.github.io/ostree/composefs/); [nixos-generators (archived)](https://github.com/nix-community/nixos-generators) (2026-01-30); [bootc-image-builder](https://github.com/osbuild/bootc-image-builder)

**ROCm support outside Ubuntu and the RHEL family** — severity 4, breadth 3, willingness to pay 2, confidence medium
- Segment: AMD GPU owners running local AI
- Description: AMD officially supports ROCm only on a short list of enterprise distributions. Fedora, NixOS and immutable distributions are not on it, so an AMD-first AI OS takes on support risk.
- Evidence: The ROCm system requirements page (2026-07-15) lists Ubuntu 24.04 and 22.04, RHEL 10.1, 10.0 and 9.x, SLES 15 SP7, Debian 13 and 12, Rocky 9 and Oracle Linux. It does not list Fedora or NixOS. The NixOS wiki says ROCm packages are cached on Hydra but that local compiles can still happen. The vLLM ROCm wheels cover only ROCm 7.0 and 7.2.1.
- Quantified impact: Fedora and NixOS are absent from the official support matrix.
- Current solutions: ROCm containers on the host's in-tree amdgpu driver; Fedora community packaging; nixpkgs rocmSupport.
- Why still unsolved: AMD focuses its validation on enterprise and LTS distributions.
- Sources: [ROCm Linux system requirements](https://rocm.docs.amd.com/projects/install-on-linux/en/latest/reference/system-requirements.html) (2026-07-15); [NixOS Wiki: AMD GPU](https://wiki.nixos.org/wiki/AMD_GPU) (2026-07-24)

**CI disk and time limits for multi-gigabyte AI images** — severity 3, breadth 4, willingness to pay 2, confidence medium
- Segment: Open-source distribution maintainers who rely on free CI
- Description: Free CI runners have little disk space for building and rechunking images that contain the CUDA userspace, and for generating ISOs.
- Evidence: Standard GitHub-hosted Linux runners for public repositories, x64 and arm64, have 4 CPUs, 16 GB RAM and a 14 GB SSD.
- Quantified impact: 14 GB SSD per job.
- Current solutions: Steps that free disk space, rechunking, splitting the image build from the disk-image build (the image-template has build.yml and build-disk.yml), and external storage through rclone to S3.
- Why still unsolved: Free-tier economics.
- Sources: [GitHub-hosted runners reference](https://docs.github.com/en/actions/reference/runners/github-hosted-runners); [ublue-os/image-template](https://github.com/ublue-os/image-template)

**Measured boot and sealed images clash with customisation** — severity 3, breadth 2, willingness to pay 2, confidence medium
- Segment: Security-focused image OS builders
- Description: A boot chain verified with UKI, composefs and TPM limits changes to the kernel command line and the use of third-party kernel modules. The tooling is not yet production-ready.
- Evidence: Fedora's sealed Atomic Desktop images (28 Apr 2026) require Secure Boot, are not signed with Fedora keys, and carry the warning not to use them in production. Kairos Trusted Boot says a change to the kernel command line needs a rebuilt image. The 2022 systemd design described signed PCR policies as future work at that time.
- Quantified impact: No production-ready sealed Fedora or bootc image exists as of April 2026.
- Current solutions: Kairos Trusted Boot, ParticleOS, Ubuntu Core FDE, lanzaboote on NixOS.
- Why still unsolved: Supporting signed third-party modules, managing keys and recovering from mistakes are still hard problems.
- Sources: [Sealed Fedora Atomic Desktop images](https://fedoramagazine.org/sealed-atomic-desktops-test-images/) (2026-04-28); [Kairos Trusted Boot](https://kairos.io/docs/architecture/trustedboot/); [Fitting Everything Together](https://0pointer.net/blog/fitting-everything-together.html) (2022-05-03)

### Key statistics

| Claim | Value | Source type | Date | Source |
|---|---|---|---|---|
| Bazzite, a Universal Blue image built on Fedora Atomic, passed 100,000 weekly active devices | >100,000 weekly active devices | news | 2026-10-15 | https://www.gamingonlinux.com/2026/10/bazzite-linux-continues-growing-now-with-over-100k-weekly-active-devices/ |
| bootc 1.17.0 declares the composefs backend stable and its on-disk format ready for production | Stable as of v1.17.0 | primary_research | 2026-10-06 | https://github.com/bootc-dev/bootc/releases/tag/v1.17.0 |
| bootc was accepted into the CNCF at Sandbox level | Sandbox, accepted 2025-01-21 | official | 2025-01-21 | https://www.cncf.io/projects/bootc/ |
| Fedora published sealed Atomic Desktop test images: systemd-boot, signed UKI, composefs with fs-verity, Secure Boot required, x86_64 and aarch64 | Test only; not for production | official | 2026-04-28 | https://fedoramagazine.org/sealed-atomic-desktops-test-images/ |
| GitHub Linux arm64 hosted runners became free for public repositories | GA 2025-08-07; 4 vCPU | official | 2025-08-07 | https://github.blog/changelog/2025-08-07-arm64-hosted-runners-for-public-repositories-are-now-generally-available/ |
| Standard GitHub-hosted Linux runners (x64 and arm64) for public repositories | 4 CPU, 16 GB RAM, 14 GB SSD | official | 2026 | https://docs.github.com/en/actions/reference/runners/github-hosted-runners |
| NVIDIA made its open GPU kernel modules the default from R560. Turing and newer are recommended, Grace Hopper and Blackwell require them, and Maxwell, Pascal and Volta still need the proprietary modules. | Open by default from R560 / CUDA 12.6 | vendor | 2024-07-17 | https://developer.nvidia.com/blog/nvidia-transitions-fully-towards-open-source-gpu-kernel-modules/ |
| The current NVIDIA driver installation guide covers branch 615 and lists Fedora, RHEL, Rocky, Alma, Ubuntu, Debian and SUSE | Branch 615 | vendor | 2026-09-23 | https://docs.nvidia.com/datacenter/tesla/driver-installation-guide/index.html |
| CUDA is redistributed directly through Canonical, CIQ, SUSE and Flox; Red Hat was added later | 5 redistributors | vendor | 2025-09-10 | https://developer.nvidia.com/blog/developers-can-now-get-cuda-directly-from-their-favorite-third-party-platforms/ |
| The Nix community CUDA cache moved to cache.nixos-cuda.org, and stable NixOS branches are not cached | No stable-branch CUDA cache | practitioner | 2025-12-02 | https://discourse.nixos.org/t/cuda-cache-for-nix-community/56038 |
| nixos-generators was archived and its functions moved into nixpkgs (nixos-rebuild build-image) from NixOS 25.05 | Archived 2026-01-30 | primary_research | 2026-01-30 | https://github.com/nix-community/nixos-generators |
| podman-bootc was archived and development moved to bcvk | Archived 2026-06-03 | primary_research | 2026-06-03 | https://github.com/containers/podman-bootc |
| image-builder cross-architecture builds are experimental and slower than native | 8-10x slower | primary_research | 2026 | https://github.com/osbuild/image-builder |
| vLLM's precompiled binaries target CUDA 12.9 (12.8 and 13.0 also built) and are incompatible with other CUDA or PyTorch builds; Python 3.11-3.14 is supported | CUDA 12.9 default | official | 2026 | https://docs.vllm.ai/en/latest/getting_started/installation/gpu.html |
| ROCm officially supports Ubuntu, RHEL 10.x and 9.x, SLES, Debian, Rocky and Oracle Linux, but not Fedora or NixOS | 0 immutable or Fedora distros listed | vendor | 2026-07-15 | https://rocm.docs.amd.com/projects/install-on-linux/en/latest/reference/system-requirements.html |
| Flatcar's NVIDIA source build at first boot takes minutes and cannot support Secure Boot; only the prebuilt signed sysext can | 5-10 minutes first-boot compile | official | 2026 | https://www.flatcar.org/docs/latest/os-config/host-config/using-nvidia/ |
| Universal Blue narrowed ublue-os/main to base, kinoite and silverblue images and removed the sway, budgie and cosmic images | Scope cut Sep 2025; removals 2025-10-15 | primary_research | 2025-10-15 | https://github.com/ublue-os/main |
| Bluefin LTS and GDX are based on CentOS Stream 10 with EPEL, kernel 6.12 for the release lifetime, aarch64 images, and NVIDIA plus CUDA with Secure Boot (GDX) | Two named maintainers | practitioner | 2025-09-18 | https://raw.githubusercontent.com/projectbluefin/documentation/main/blog/2025-09-18-bluefin-lts-and-gdx.md |
| ParticleOS, systemd's image OS built with mkosi, is still in development with no backward compatibility guarantees | Prebuilt images x86-64 only | primary_research | 2026 | https://github.com/systemd/particleos |
| Kairos (distro-agnostic immutable OS with Trusted Boot) current stable release | v4.3.0 | official | 2026-02-24 | https://kairos.io/docs/ |
| systemd-sysext overlays only /usr and /opt, matches images on host ID, VERSION_ID or SYSEXT_LEVEL and ARCHITECTURE, and has no dependency resolution; mutable mode arrived in v256 | /usr and /opt only | official | 2026 | https://man7.org/linux/man-pages/man8/systemd-sysext.8.html |
| Fedora Asahi Remix is based on Fedora 44 and supports M1 Pro MacBook Pros | Fedora 44 | official | 2026 | https://asahilinux.org/fedora/ |
| nixified-ai deprecated its InvokeAI and textgen packages because of a funding shortfall; ComfyUI is the main maintained output | 3 of 4 outputs deprecated or broken | primary_research | 2026 | https://github.com/nixified-ai/flake |

### Players

- **bootc and Fedora/CentOS bootc base images (Red Hat, Fedora, CentOS)** (Bootable OCI image foundation): Transactional OS updates and rollback delivered as OCI images. Stable CLI and API, ostree default backend, composefs backend stable in 1.17.0. Logically bound images tie app containers to OS updates. Install via to-disk using systemd-repart.. Traction: 1,275 contributors per CNCF; RHEL image mode; Rocky bootc via CIQ (Feb 2026). Gap: CNCF Sandbox. Sealed images are test-only. Tooling churn (podman-bootc archived, BIB moved). Test tools assume a Linux host. Bound images do not support rootless containers or Anaconda.
- **Universal Blue (Bluefin, Aurora, Bazzite, akmods, image-template)** (Tiny-team downstream image factory): Containerfile plus GitHub Actions plus Renovate plus cosign. Prebuilt signed NVIDIA open and proprietary and ZFS modules published as OCI images. Template for custom images with ISO and qcow2 generation.. Traction: Bazzite >100k weekly devices (Oct 2026); Bluefin repo 3,111 commits, 2.6k stars. Gap: Scope cuts in 2025. Secure Boot needs MOK enrollment. Kernel pins when ZFS or NVIDIA lag. Mostly desktop-focused.
- **Bluefin GDX / Bluefin LTS** (Direct competitor: AI workstation image): CentOS Stream 10 plus NVIDIA drivers plus CUDA plus RamaLama plus uv, x86_64 and aarch64, Secure Boot. Traction: GA Sep 2025. Gap: NVIDIA only. No agent orchestration layer. No local layering. HWE branch lacks Secure Boot.
- **NixOS plus Nixpkgs CUDA team** (Declarative reproducible OS): Whole-system closures, generation rollback, build-image for ISO and raw, cudaSupport and rocmSupport flags. Traction: Volunteer CUDA team; community Hydra. Gap: No CUDA on cache.nixos.org; no cache for stable releases; Secure Boot via lanzaboote not yet in nixpkgs; steep learning curve; ML flakes underfunded
- **Flox** (Nix-based environment manager with CUDA redistribution): NVIDIA-approved prebuilt CUDA, PyTorch and ONNX closures in a public cache. Traction: Named NVIDIA CUDA redistributor (Sep 2025). Gap: Builds only from its own fork that follows unstable; depends on volunteer nixpkgs maintainers
- **Ubuntu Core (Canonical)** (Snap-based transactional OS): Immutable system with automatic rollback, FDE, Ubuntu Core 26, CUDA workloads and Jetson guides. Traction: Embedded and IoT focus. Gap: Everything is a snap; a dedicated brand store is needed for a custom distribution; strict confinement makes Python ML stacks awkward
- **mkosi / ParticleOS (systemd)** (Image builder and reference verified-boot OS): UKI, signed dm-verity, systemd-sysupdate, TPM and FIDO2 LUKS, user-owned keys. Traction: 623 stars (ParticleOS). Gap: Experimental; no compatibility guarantees; prebuilt images x86-64 only; builds only on Linux
- **Talos Linux (Sidero Labs)** (Kubernetes appliance OS): Configured only through an API, immutable, NVIDIA via system extensions. Traction: v1.14 latest stable. Gap: Kubernetes only; NVIDIA extensions tied to each Talos release
- **Flatcar Container Linux** (Container host OS): A/B updates, sysext-bakery (nvidia-runtime, ollama, k3s), signed prebuilt NVIDIA sysext. Traction: Bakery publishes via extensions.flatcar.org. Gap: Server-focused; prebuilt NVIDIA only for amd64
- **Kairos (Spectro Cloud, LF Projects)** (Distro-agnostic immutable OS converter): Turns container images into an immutable OS with active, passive and recovery systems and Trusted Boot (UKI, TPM, FDE). Traction: v4.3.0 (Feb 2026). Gap: Edge and Kubernetes framing; kernel command line changes need rebuilds; kernel modules not addressed
- **BlueBuild** (Declarative YAML recipe tooling for atomic and bootc images): YAML recipes compiled to Containerfiles; reusable modules. Traction: Spun out of Universal Blue. Gap: Non-Fedora support experimental; no ISO tooling for non-Fedora images
- **RamaLama (containers project)** (Containerised model serving): Images per accelerator (CUDA, ROCm, Vulkan, Intel); llama.cpp and vLLM; rootless with network disabled. Traction: Bundled in Bluefin GDX. Gap: Alpha; isolation guarantees documented for run, not serve

### Trends
- Bootable OCI images are becoming the shared delivery format across vendors: bootc entered the CNCF in Jan 2025, CIQ shipped Rocky bootc images in Feb 2026, Fedora built sealed test images in Apr 2026, and bootc declared composefs stable in Oct 2026.
- Verified boot chains (signed UKI, composefs or dm-verity, TPM-bound disk unlock) are moving from design documents (systemd 2022) to test images (Fedora 2026). Production readiness is still pending.
- NVIDIA's open kernel modules became the default (R560, 2024), and CUDA is now redistributed by OS vendors and Flox (2025). Packaging is getting simpler, but module signing is still each distribution's job.
- systemd-sysext is spreading as the way to deliver optional host components: Flatcar sysext-bakery, fedora-sysexts, and Talos system extensions.
- Free native arm64 CI (GA Aug 2025) removes the need for cross-architecture emulation for open-source projects.
- Image tooling is consolidating: podman-bootc into bcvk, bootc-image-builder into image-builder, nixos-generators into nixpkgs.
- AI workstation images already exist (Bluefin GDX). Whether the OS layer can be a differentiator is shrinking.
- Small teams are shrinking their scope to survive, for example the 2025 cuts to ublue-os/main.

### Opportunities

- **Containers by default for AI runtimes: the OS detects the GPU at first boot and pulls a matching, signed runtime image (CUDA, ROCm, Vulkan or CPU), pinned by digest and tied to the OS release through bootc logically bound images** Why now: bootc logically bound images exist, RamaLama shows the one-image-per-accelerator pattern, and vLLM's CUDA pinning makes host-level installs unworkable Evidence: bootc LBI docs; vLLM docs warn its binaries are incompatible with other CUDA or PyTorch builds; RamaLama ships CUDA, ROCm, Vulkan and Intel images Risk: Bound images do not support rootless containers. Multi-gigabyte images make updates slow. RamaLama is alpha.
- **Attested agent host: use sealed UKI and composefs plus TPM measured boot so agents and remote model endpoints can verify that the host has not been modified** Why now: Fedora's sealed images and bootc's stable composefs arrived in 2026. Agent security is a fresh concern. Evidence: Fedora sealed test images (Apr 2026); bootc 1.17.0; ParticleOS and Kairos Trusted Boot designs Risk: Not production-ready. Third-party NVIDIA modules complicate a sealed chain. It would take most of a solo maintainer's time.
- **Optional AI capability packs delivered as sysexts with systemd-sysupdate (for example inference CLIs, GPU monitoring and debugging tools), following the Flatcar bakery pattern** Why now: systemd has mature sysext and sysupdate support (mutable mode since v256), and Flatcar already proves distribution through GitHub releases Evidence: Flatcar sysext-bakery includes ollama and nvidia-runtime; fedora-sysexts publishes via sysupdate redirectors Risk: No dependency resolution. Each pack must be rebuilt for every OS VERSION_ID or SYSEXT_LEVEL. Sysexts cannot easily carry kernel modules.
- **A Mac-first workflow for building an OS: local aarch64 builds and VM boots in podman machine, with CI handling x86_64, ISOs and GPU smoke tests** Why now: Free arm64 runners and the aarch64 bootc base images make Apple Silicon a workable authoring platform. Few projects document this. Evidence: GitHub arm64 GA (Aug 2025); Bluefin LTS and GDX aarch64 images; image-builder cross-architecture is 8-10x slower, so builds should be native per architecture Risk: GPU and TPM paths still need real x86 hardware.

### Contrarian signals
- bootc is still a CNCF Sandbox project. CNCF defines Sandbox as not yet widely tested in production.
- Containerfile builds are not bit-for-bit reproducible, because dnf installs pull whatever packages are current at build time. NixOS gives stronger reproducibility, which undercuts the reproducibility argument for bootc.
- The OSTree composefs docs still say not to deploy to production, while bootc 1.17.0 calls composefs stable. Readiness signals conflict.
- Bluefin GDX already ships the base layer of an AI workstation: NVIDIA drivers, CUDA and RamaLama. A new OS base adds little that is new, so differentiation has to come from the agent layer.
- Universal Blue cut image variants in 2025. If an established tiny team has to trim, a solo part-time maintainer should expect even less capacity.
- NVIDIA breakages on Fedora (a kernel change blocking the driver in Jun 2025, driver 610 dependency failures in May 2026) show that fast-moving bases break proprietary GPU stacks.
- Bazzite's 100k figure counts devices through DNF countme, not users, and the project reports it itself.
- The M1 Pro cannot test NVIDIA, x86 Secure Boot or TPM. A Mac-only workflow cannot validate the product's core value.
- Wanting a revolutionary bootable OS conflicts with working solo and part-time. Image-based bases are deliberately conventional, so the innovation budget belongs above the OS layer.
- vLLM's own docs present Docker only as an option. Users may prefer pip or uv on a general-purpose distribution, which weakens the case for a bootable OS.

### Open questions
- Clarifying question: is AgentOS a desktop workstation, a headless agent appliance, or both? The answer decides between deriving from Bluefin LTS and starting from a minimal centos-bootc image.
- Clarifying question: which accelerators must be first-class at launch: NVIDIA only, AMD ROCm, or Apple Silicon through Asahi? Each one multiplies the test matrix.
- Not verified: does ublue-os/akmods publish modules for CentOS Stream 10 kernels, or does Bluefin LTS/GDX get its EL NVIDIA modules from elsewhere?
- Will Fedora and bootc sealed images support third-party signed kernel modules such as NVIDIA's?
- Do standard GitHub-hosted runners expose KVM for bcvk boot and rollback tests? Not verified.
- What exactly are the redistribution terms for bundling NVIDIA CUDA userspace in an open-source OS image?
- Method limitation: the shared web-search budget ran out after 6 searches, below the 20 the brief asked for. Some primary docs (Fedora docs, freedesktop.org, Universal Blue Discourse) blocked automated fetches. AlmaLinux bootc, Ubuntu Core TPM details and RHEL AI were not covered.

### Adversarial verification

- **DISPUTED** — Bluefin GDX already ships CentOS Stream 10 with NVIDIA drivers, CUDA, RamaLama and uv on x86_64 and aarch64 with Secure Boot, and Bluefin LTS/GDX keep kernel 6.12 for the release lifetime. The brief uses this as the reference model behind its CentOS Stream 10 recommendation.
  - Corrected: The September 2025 description is out of date. GDX moved to the Fedora CoreOS (HWE) kernel on 11 Jan 2026. In July 2026 Project Bluefin announced that GDX would stop being a separate image and become a userspace mode, with CUDA delivered through containers, and that LTS would move entirely onto Fedora kernels. Project Bluefin's own LTS repository was later rebranded Uranium (thallium-suite / quill-oci), and that repository was archived on 15 Oct 2026. Only the legacy ublue-os/bluefin-lts repository still builds GDX, for amd64 and arm64. It does this with a Fedora CoreOS kernel and Fedora NVIDIA packages from negativo17 on top of CentOS userspace. Current Secure Boot status for GDX was not verified.
  - Evidence: Sources checked: the GitHub API, the ublue-os/bluefin-lts Justfile, and the projectbluefin blog. In the Justfile, gdx=1 sets AKMODS_BASE=coreos-stable-43. The kernel-swap.sh history includes 'fix: switch GDX to HWE kernel (#1007)', merged 2026-01-11. The GDX script 20-nvidia.sh installs nvidia-driver-cuda from negativo17's fedora-nvidia repository. The Project Bluefin blog post of 2 Jul 2026 says DX and GDX images are merged, GDX moves to 'ujust aimode', CUDA can be consumed through containers, and everything is consolidated onto Fedora's kernels. It gives the reason that GitHub runners lacked the space and resources for GDX. The quill-oci/uranium repository shows archived=true, last pushed 2026-10-15, with the commit 'bump kernel version to 7.2.5 from fedore core os'. The legacy GDX workflow last succeeded on 2026-10-15.
  - Sources: [Organizational Migration for Bluefin LTS/GDX](https://github.com/projectbluefin/documentation/blob/main/blog/2026-07-02-organizational-migration-for-bluefin-lts-gdx.md) (2026-07-02); [fix: switch GDX to HWE kernel (PR #1007)](https://github.com/ublue-os/bluefin-lts/pull/1007) (2026-01-11); [ublue-os/bluefin-lts Justfile](https://github.com/ublue-os/bluefin-lts/blob/main/Justfile) (2026-10-15); [quill-oci/uranium (archived, formerly Bluefin LTS)](https://github.com/quill-oci/uranium) (2026-10-15); [Bluefin LTS docs (now Uranium)](https://docs.projectbluefin.io/lts) (2026-10); [Bluefin LTS and GDX GA post](https://raw.githubusercontent.com/projectbluefin/documentation/main/blog/2025-09-18-bluefin-lts-and-gdx.md) (2025-09-18)
- **PARTIALLY_CONFIRMED** — bootc 1.17.0 (6 Oct) declared its composefs backend stable, and its CLI and API are stable.
  - Corrected: bootc v1.17.0, published 2026-10-06T18:15:02Z, declares the composefs backend stable and calls the 'core on-disk format' ready for production. ostree is still the default backend. The composefs known-issues list says logically bound images are not implemented, there is no in-place move from an ostree system, and upgrades of UKI installs from other releases are not yet tested. The statement that the CLI and API are stable comes from the bootc docs, not from this release, and it covers surface APIs only. bootc is still a CNCF Sandbox project.
  - Evidence: The GitHub releases API shows published_at 2026-10-06T18:15:02Z, which settles the brief's uncertainty about the year. The release body says 'This release declares the composefs backend as stable.' The bootc-composefs(7) page says the backend 'is stable', but it also lists logically bound images as not implemented, notes the missing ostree-to-composefs transition, and lists unimplemented commands (bootc edit, install reset). The bootc docs say the CLI and API are 'considered stable' and that this 'refers to the surface APIs'. The CNCF lists bootc as Sandbox (accepted 21 Jan 2025). This creates a conflict in the brief: it recommends logically bound images, but the composefs or sealed path it plans to adopt later does not support them yet.
  - Sources: [bootc v1.17.0 release](https://github.com/bootc-dev/bootc/releases/tag/v1.17.0) (2026-10-06); [bootc-composefs(7)](https://bootc.dev/bootc/bootc-composefs.7.html) (2026-10); [bootc documentation home](https://bootc.dev/bootc/) (2026-10); [CNCF project page: bootc](https://www.cncf.io/projects/bootc/) (2025-01-21)
- **CONFIRMED** — Bazzite, built the Universal Blue way, passed 100,000 weekly active devices (15 Oct 2026).
  - Corrected: Confirmed. Universal Blue's own countme badge shows Bazzite at '100k', and the project charts this as 'Weekly Active Devices' from Fedora DNF countme data. The founder's post calls them 'users'. Bazzite is Fedora Kinoite-based and mostly x86_64. The same dataset puts the CentOS-based Bluefin LTS at 64, Bluefin at 4.1k and Aurora at 3.1k. The scale proof therefore applies to the Fedora path, not to the CentOS Stream 10 base the brief recommends.
  - Evidence: GamingOnLinux, 15 Oct 2026, quotes founder Kyle Gospodnetich: 'passed 100,000 weekly active users'. The ublue-os/countme badge-endpoints/bazzite.json file reads message '100k' (updated 2026-10-08), and the plot title in countme.py is 'Weekly Active Devices'. A code comment in data_processing.py says Bluefin LTS data 'is not gathered from fedora repos', so the LTS figure may be undercounted. The Bazzite Containerfile builds FROM ublue-os/kinoite-main with akmods, and the akmods Fedora main and ogc targets are x86_64 only.
  - Sources: [Bazzite continues growing, now with over 100k weekly active devices](https://www.gamingonlinux.com/2026/10/bazzite-linux-continues-growing-now-with-over-100k-weekly-active-devices/) (2026-10-15); [ublue-os/countme (badge endpoints and methodology)](https://github.com/ublue-os/countme) (2026-10-08); [ublue-os/akmods images.yaml](https://github.com/ublue-os/akmods/blob/main/images.yaml) (2026-10)
- **CONFIRMED** — GitHub gives public repositories free x64 and arm64 runners (4 vCPU, 16 GB RAM, 14 GB SSD), so x86_64 and aarch64 images can both be built for free in CI.
  - Corrected: Confirmed. Standard Linux runners for public repositories are 4 CPU, 16 GB RAM and 14 GB SSD, for both x64 (ubuntu-24.04, ubuntu-26.04) and arm64 (ubuntu-24.04-arm, ubuntu-26.04-arm). GitHub states that usage is free for public repositories on standard runners, and that larger runners are always charged. Being free does not make them big enough: Project Bluefin reported in July 2026 that these runners lacked the space and resources for its CUDA-bearing GDX image.
  - Evidence: The GitHub-hosted runners reference table lists 'Linux 4 16 GB 14 GB' for both x64 and arm64. The billing docs say: 'GitHub Actions usage is free for self-hosted runners and for public repositories that use standard GitHub-hosted runners.' The runner docs do not say whether KVM or nested virtualisation is available on Linux runners, so that remains unverified. The Bluefin LTS reusable workflow builds amd64 on ubuntu-24.04 and arm64 on ubuntu-24.04-arm, and includes a 'Maximize build space' step.
  - Sources: [GitHub-hosted runners reference](https://docs.github.com/en/actions/reference/runners/github-hosted-runners) (2026-10); [GitHub Actions billing](https://docs.github.com/en/billing/concepts/product-billing/github-actions) (2026-10); [Organizational Migration for Bluefin LTS/GDX](https://github.com/projectbluefin/documentation/blob/main/blog/2026-07-02-organizational-migration-for-bluefin-lts-gdx.md) (2026-07-02)
- **PARTIALLY_CONFIRMED** — AMD lists RHEL 10.x as supported for ROCm, but not Fedora or NixOS, which favours a RHEL-family base; a CentOS Stream 10 base matches the RHEL 10 support listed by NVIDIA and ROCm.
  - Corrected: Fedora and NixOS are indeed absent from ROCm's support list. CentOS Stream and AlmaLinux are absent too, so a CentOS Stream 10 base is not officially supported either. The brief's source (dated 2026-07-15, ROCm 7.14 era) has been replaced by ROCm Core SDK 10.1.0, released 2026-10-05 and built with TheRock. That release lists RHEL 10.2, 10.0, 9.8, 9.6, 9.4 and 8.10, plus Ubuntu 26.04.1, 24.04.5 and 22.04.5, SLES 16.0 and 15.7, Debian, Oracle Linux and Rocky 9. Consumer Radeon support is limited to Ubuntu and selected RHEL releases. An image that swaps in a Fedora CoreOS kernel, as Bluefin LTS now does, gives up any RHEL kernel alignment.
  - Evidence: The system-requirements page (dated 2026-07-15) lists RHEL 10.1, 10.0, 9.7, 9.6, 9.4 and 8.10, but not Fedora, CentOS Stream or AlmaLinux. It limits Radeon and Radeon PRO to Ubuntu 24.04.4, 22.04.5, RHEL 10.1 and 9.7, and carries a banner saying it has moved to the compatibility matrix. That matrix and the release notes show ROCm 10.1.0 (2026-10-05) with the newer distribution list. Neither page explains the jump from 7.x to 10.x, beyond noting that 'Since ROCm 7.14, ROCm uses TheRock'.
  - Sources: [ROCm Linux system requirements](https://rocm.docs.amd.com/projects/install-on-linux/en/latest/reference/system-requirements.html) (2026-07-15); [ROCm compatibility matrix](https://rocm.docs.amd.com/en/latest/compatibility/compatibility-matrix.html) (2026-10); [ROCm Core SDK 10.1.0 release notes](https://rocm.docs.amd.com/en/latest/about/release-notes.html) (2026-10-05)

**Missed counter-evidence**
- Project Bluefin restructured in July 2026. DX and GDX images were merged, AI and CUDA tooling moved to a userspace 'ujust aimode' with CUDA through containers, signing moved to keyless OIDC, and layers are now built with chunkah. This matches the brief's advice to keep CUDA out of /usr, but it removes the GDX reference image the brief relies on.
- Bluefin LTS was renamed Uranium under thallium-suite / quill-oci (the quill-oci organisation was created 12 Oct 2026), and that repository was archived on 15 Oct 2026. The README says Uranium 'has no community infrastructure of its own yet'. A CentOS-based downstream maintained by a tiny team turns out to be organisationally fragile.
- Project Bluefin is moving its effort to new bases. Dakota is GNOME OS built with BuildStream (alpha). Utah is built on Fedora Hummingbird (pre-alpha, nothing published). Bluefin Server, announced in Sep 2026, is a freedesktop-sdk, DDI-first OS with signed UKIs, a dm-verity /usr and systemd-sysupdate (alpha). These are important references for the verified-boot design and a sign that the reference team is leaving CentOS Stream.
- Universal Blue already runs a CentOS 10 NVIDIA module pipeline (Bluefin LTS mounts ghcr.io/ublue-os/akmods-nvidia-open:centos-10). This answers the brief's open question. The kernel and modules are signed with Universal Blue's shared key, so the founder would either depend on that key or run their own pipeline.
- The ublue-os/akmods nvidia-open builds for the Fedora main and ogc flavours (44 and 45) are x86_64 only. aarch64 NVIDIA modules exist only for the coreos-stable flavours, which affects the plan for an aarch64 agentos-nvidia image.
- The brief cited the Nix CUDA cache thread selectively. A follow-up post on 17 Dec 2025 said Flox caches stable channels. Flox's documentation shows stable (monthly) and lts (six-monthly) cadence branches of its own fork. Release channels are still not built, so 'no cache for stable releases' is broadly accurate but overstated.
- Universal Blue's countme data puts Bluefin LTS at 64 active users, against 4.1k for Bluefin and 100k for Bazzite. The data may undercount CentOS hits, but it does not support the claim that a tiny team ships at scale on a CentOS Stream 10 base.
- ROCm 10.1.0 (5 Oct 2026) adds Ubuntu 26.04.1 and SLES 16 and is built with TheRock. The brief's ROCm version details and its vLLM ROCm wheel claims (ROCm 7.0 and 7.2.1) are probably out of date.
- Because the bootc composefs backend has no logically bound images and no in-place move from ostree, a system installed today on the default ostree backend has to be reinstalled to adopt composefs or sealed images later. The brief's advice to 'postpone measured boot' therefore carries a reinstall cost.

**Reliability note:** Summary: the brief's dated facts mostly hold, but its central reference model is out of date. Bluefin GDX, as an image with kernel 6.12 for the release lifetime, was retired or restructured in 2026, and the LTS line has been spun out and archived. That weakens the claim that CentOS Stream 10 is the 'only base that meets all four needs'. Method: the shared web-search budget ran out on the first query, so all verification used direct fetches of primary documents and read-only GitHub API calls. A small model summarised the WebFetch pages, so quoted wording should be spot-checked; the ROCm 10.1.0 version jump in particular is unexplained by AMD. Discourse pages for Universal Blue returned HTTP 403. Clarifying questions for the founder: (1) Does AgentOS need RHEL-kernel ABI stability at all, given that upstream Bluefin abandoned the 6.12 kernel for lack of users? (2) Is a Fedora-based bootc image, which carries the 100k-device evidence, acceptable if ROCm and NVIDIA run only in containers?


---

## Developing and testing a bootable AI OS on an Apple M1 Pro with 16 GB
_Workflow label: `m1_dev_feasibility`_

### Executive summary
Bottom line: the M1 Pro is a good machine for daily AgentOS development, but it should not be the main bare-metal target. It also cannot validate x86_64 plus NVIDIA builds. The best plan uses three tiers: macOS with VMs, Asahi dual-boot and cheap cloud or CI.

A note on method: the shared WebSearch budget ran out after 3 searches. The findings below rest on about 45 direct reads of primary sources, so some 2026 prices remain unverified.

- **Asahi status.** Fedora Asahi Remix 43 shipped on 18 March 2026, and version 44 was targeted for 28 April 2026. M1 Pro covers GPU, suspend, camera, speakers and battery in linux-asahi. Thunderbolt and DP Alt Mode are still WIP. Video encode is TBA, and the Neural Engine driver is out of tree.
- **GPU stack.** The stack is conformant for OpenGL 4.6, Vulkan 1.4 and OpenCL 3.0. However, the 21,000-line kernel GPU driver is still downstream (February 2026). The lead GPU developer left in August 2025, which is a bus-factor risk.
- **Inference on Asahi.** On a base M1, llama.cpp Vulkan reaches 38 t/s prompt processing (pp512) and 12.5 t/s generation (tg128). Metal on macOS reaches 118 and 14.2. Prompt processing on Linux is therefore about three times slower.
- **MLX on Linux.** It runs only on CPU or CUDA, so Asahi gets no MLX GPU path.
- **Inference in a macOS VM.** krunkit/libkrun exposes Vulkan compute (no rendering) through Venus and MoltenVK. Red Hat measured 74–80% of native Metal. A VM on macOS can therefore beat bare-metal Asahi for LLM work.
- **Other VM options.** Lima 2.0 includes an experimental krunkit driver. Apple's container tool (1.5.0, September 2026) needs macOS 26 and has no GPU support. UTM offers only VirGL OpenGL.
- **No nested virtualization.** It requires M3 or later, so KVM-based agent sandboxes can be tested only on bare-metal Asahi or x86 cloud.
- **16 GB memory.** macOS gives the GPU about 66.7% of RAM by default, roughly 10.7 GB. Dense 7–8B Q4 models run well (about 36 t/s on 7B Q4_0, M1 Pro 16-core). 14B is borderline, and 30B MoE models do not fit. QLoRA fine-tuning is realistic only at 8B or below.
- **Testing x86_64 and NVIDIA.** Public GitHub repositories get free x64 runners with KVM for QEMU boot tests, plus free arm64 runners. GitHub GPU runners (T4, $0.052/min) need a paid plan and cannot boot a custom OS. AWS g4dn.xlarge can boot custom AMIs for $0.33/hr spot or $0.53/hr on demand. RunPod RTX 3090/4090 pods ($0.22–0.34/hr) can test userspace only.

Estimated cost: $0 for tooling and about $15–40 per month for cloud.

### Recommendations
- Treat macOS on the M1 Pro as the primary daily loop. Build aarch64 bootc images in a Podman machine using the libkrun provider. Generate qcow2 or raw images with bootc-image-builder, then boot them in krunkit (Lima --vm-type=krunkit) for GPU Vulkan compute tests. Use UTM (QEMU with HVF) for installer ISO walkthroughs. Tooling cost: $0.
- Run inference and fine-tuning experiments on the macOS host with MLX or llama.cpp Metal, and expose them to the guest over HTTP. Do not benchmark AgentOS performance on Asahi or in a VM and present those numbers as representative.
- Dual-boot Fedora Asahi Remix 44 on the same M1 Pro, using about 80–100 GB of disk. Validate weekly: the 16K-page kernel, KVM-based agent sandboxes (impossible in a VM on M1), Vulkan llama.cpp, s2idle suspend and battery. Use the Asahi installer's 'UEFI environment only' option so that custom AgentOS images can boot from USB through U-Boot.
- Ship two aarch64 kernel flavors: a generic 4K-page kernel for VMs and Arm clouds, and the Asahi 16K-page kernel for Mac bare metal. Consume Fedora Asahi kernel packages rather than maintaining a fork.
- Make x86_64 plus NVIDIA the primary product target. In a public GitHub repository, use free ubuntu-latest runners (x64, KVM) to build images and run QEMU boot tests over a serial console. Use free ubuntu-24.04-arm runners for aarch64 builds.
- For NVIDIA boot and driver validation, launch AWS g4dn.xlarge spot instances (about $0.33/hr, T4 16 GB) from bootc-image-builder AMIs at each release, with auto-termination through GitHub Actions and OIDC. Add g6.xlarge spot (about $0.62/hr, L4 24 GB) for newer driver branches.
- Use RunPod community GPUs (RTX 3090 at $0.22/hr, RTX 4090 at $0.34/hr) only for userspace work: CUDA runtime, vLLM, and QLoRA jobs above 8B. Avoid GitHub GPU runners ($3.12/hr, paid plan, no custom OS).
- Default model targets for 16 GB: 3–4B and 7–8B at Q4 (about 36 t/s on 7B Q4_0, M1 Pro 16-core, Metal). Use 14B Q4 only with a raised iogpu.wired_limit_mb and short context. Test gpt-oss-20b MXFP4 as a borderline case. Avoid 30B MoE models. Limit local fine-tuning to MLX QLoRA at 8B or below.
- Ignore the Apple Neural Engine for AgentOS, because the driver is out of tree and has no LLM path.
- Budget: $0 for tooling plus about $15–40 per month for cloud (roughly 20 spot GPU hours, storage and egress). Consider an optional one-time x86 plus NVIDIA test machine with KVM-over-IP for true bare-metal and Secure Boot testing. Its cost is unverified in this research, so get quotes first.

### Pain points

**No first-class GPU ML stack for Linux on Apple Silicon** — severity 4, breadth 3, willingness to pay 2, confidence high
- Segment: Solo developer building an AI OS on M1/M2 hardware
- Description: On Asahi, the only GPU path is Mesa Vulkan (Honeykrisp) and OpenCL. Prompt processing in llama.cpp runs about three times slower than Metal on macOS. MLX has no Linux GPU backend for Apple hardware, and the Neural Engine is out of tree.
- Evidence: Apple M1 on Asahi (Vulkan): pp512 38.29, tg128 12.47 t/s on Llama 2 7B Q4_0, against 117.96/14.15 on macOS Metal for the 8-core M1. A commenter attributes the gap to Honeykrisp. MLX documents only CPU (mlx[cpu]) and CUDA installs for Linux.
- Quantified impact: Prompt processing about 3x slower and token generation about 12% slower than macOS Metal on a base M1.
- Current solutions: Use macOS Metal or MLX for inference and treat Asahi Vulkan as a functional test path only.
- Why still unsolved: Apple publishes no GPU documentation, so the Linux stack depends on reverse engineering. The GPU lead left in August 2025, and team effort has moved to M3 bring-up and upstreaming.
- Sources: [llama.cpp Vulkan performance scoreboard #10879](https://github.com/ggml-org/llama.cpp/discussions/10879); [llama.cpp Apple Silicon performance #4167](https://github.com/ggml-org/llama.cpp/discussions/4167) (2023-11-22); [MLX install docs](https://ml-explore.github.io/mlx/build/html/install.html); [Dissecting the Apple M1 GPU, the end](https://alyssarosenzweig.ca/blog/asahi-gpu-part-n.html) (2025-08-26)

**Bare-metal Apple Silicon support depends on a downstream kernel** — severity 4, breadth 3, willingness to pay 1, confidence high
- Segment: Distro builders targeting Apple Silicon
- Description: A bootable AgentOS image for Apple Silicon needs the Asahi kernel fork, m1n1, Asahi U-Boot, firmware extracted by the Asahi installer, and a 16K-page kernel. It cannot be booted from USB on a fresh machine without first running the installer from macOS.
- Evidence: The Asahi distro policy requires the Asahi kernel fork, m1n1, U-Boot and asahi-firmware. The GPU driver is about 21,000 lines and only began upstreaming in February 2026. The 7.2 report (August 2026) says GPU and DCP work remains too large to upstream. The NixOS Apple Silicon guide shows that USB boot via U-Boot works after a 'UEFI environment only' install.
- Quantified impact: More than 1,000 downstream patches were cited in February 2025, with a 95,000 to 83,000 line reduction reported in February 2026.
- Current solutions: Build on top of Fedora Asahi Remix packages, or follow the nixos-apple-silicon model.
- Why still unsolved: Upstreaming is slow, the project lost its founder (February 2025) and GPU lead (August 2025), and Apple's boot chain requires the installer step from macOS.
- Sources: [Asahi alternative distro policy](https://asahilinux.org/docs/alt/policy/); [Progress Report: Linux 6.19](https://asahilinux.org/2026/02/progress-report-6-19/) (2026-02-15); [Progress Report: Linux 7.2](https://asahilinux.org/2026/08/progress-report-7-2/) (2026-08-26); [nixos-apple-silicon UEFI standalone guide](https://github.com/nix-community/nixos-apple-silicon/blob/main/docs/uefi-standalone.md)

**GPU access in macOS VMs is limited** — severity 3, breadth 4, willingness to pay 2, confidence high
- Segment: Developers testing GPU-dependent Linux images on Macs
- Description: Apple's Virtualization.framework gives Linux guests no 3D acceleration. krunkit/libkrun exposes Vulkan compute only (no rendering) through Venus and MoltenVK, and Lima's krunkit driver is experimental. UTM offers only experimental VirGL OpenGL.
- Evidence: Podman Desktop says macOS GPU support covers only Vulkan compute shaders, not rendering or draw. Red Hat measured 74–80% of native Metal (20.84 t/s for llama3.1:8b on an M4 Pro). Lima's krunkit driver needs Lima 2.0 or later and macOS 14 or later, and is experimental. A 2024 design note set an 8 GB GPU memory cap.
- Quantified impact: 20–26% throughput loss against native Metal, with no graphics rendering.
- Current solutions: krunkit through Podman machine (libkrun provider), Lima, or RamaLama. Alternatively, run inference on the host and expose it to the guest over an API.
- Why still unsolved: Apple offers no paravirtual GPU for Linux guests, and translation layers add overhead. KosmicKrisp (Vulkan on Metal 4) needs macOS 26.
- Sources: [Podman Desktop GPU docs](https://podman-desktop.io/docs/podman/gpu); [How we improved AI inference on macOS Podman containers](https://developers.redhat.com/articles/2025/06/05/how-we-improved-ai-inference-macos-podman-containers) (2025-06-05); [Lima krunkit driver](https://lima-vm.io/docs/config/vmtype/krunkit/); [Enabling containers to access the GPU on macOS](https://sinrega.org/2024-03-06-enabling-containers-gpu-macos/) (2024-03-06)

**No nested virtualization on M1** — severity 4, breadth 3, willingness to pay 2, confidence high
- Segment: AI OS builders using microVM agent sandboxes
- Description: KVM is unavailable inside Linux VMs on the M1 Pro. Firecracker, Kata and microVM-per-agent sandboxes therefore cannot be tested in a VM on the development Mac.
- Evidence: Lima's template states that nested virtualization is supported with vz or krunkit only on Apple M3 or later with macOS 15 or later. GitHub macOS runners also lack nested virtualization.
- Quantified impact: An entire class of isolation features cannot be tested on the primary development VM.
- Current solutions: Run sandboxes on bare-metal Asahi, which has KVM and DRM native context, or on x86 KVM in GitHub Actions.
- Why still unsolved: Hardware support for nested virtualization arrived only with M3.
- Sources: [Lima default.yaml template](https://raw.githubusercontent.com/lima-vm/lima/master/templates/default.yaml); [GitHub-hosted runners reference](https://docs.github.com/en/actions/reference/runners/github-hosted-runners)

**Testing the x86_64 plus NVIDIA boot path cheaply** — severity 4, breadth 4, willingness to pay 3, confidence medium
- Segment: Solo maintainers of bootable AI OS images
- Description: Container GPU clouds such as RunPod and Vast cannot boot a custom kernel or bootloader. GitHub GPU runners are paid, run GitHub's own images and require Team or Enterprise plans. Custom-image GPU VMs and bare-metal servers cost more and need automation.
- Evidence: GitHub's Linux 4-core GPU runner (Tesla T4, 16 GB) costs $0.052/min, and larger runners are always charged, even for public repositories. AWS g4dn.xlarge costs $0.526/hr on demand and $0.325/hr spot. Vast VMs must use docker.io/vastai/kvm images. RunPod lists only pod (container) products.
- Quantified impact: $0.33–0.62 per hour for custom-OS NVIDIA boot tests, against $3.12 per hour for GitHub GPU runners.
- Current solutions: Build AMIs with bootc-image-builder and boot them on spot g4dn or g6 instances. Use RunPod for userspace tests.
- Why still unsolved: GPU clouds optimize for ML containers rather than OS developers, and custom-image GPU capacity is concentrated in hyperscalers.
- Sources: [GitHub Actions runner pricing](https://docs.github.com/en/billing/reference/actions-runner-pricing); [GitHub larger runners](https://docs.github.com/en/actions/reference/runners/larger-runners); [Vantage g4dn.xlarge](https://instances.vantage.sh/aws/ec2/g4dn.xlarge) (2026-10-15); [Vast.ai VMs](https://docs.vast.ai/instances/virtual-machines); [RunPod pricing](https://www.runpod.io/pricing)

**The 16 GB unified memory ceiling** — severity 3, breadth 4, willingness to pay 2, confidence medium
- Segment: Developers on 16 GB Apple Silicon
- Description: By default macOS lets the GPU use about two thirds of RAM on machines with 32 GB or less, roughly 10.7 GB here. This limits local models to small dense models and fine-tuning to QLoRA on 8B or smaller.
- Evidence: The llama.cpp discussion shows the default reserve is 33.3% (25% above 32 GB), and iogpu.wired_limit_mb can raise the limit. Unsloth lists QLoRA minimums of 6 GB for 8B, 8.5 GB for 14B and 26 GB for 32B. gpt-oss-20b (21B total, 3.6B active, MXFP4) claims to run within 16 GB.
- Quantified impact: Roughly 10.7 GB is usable by the GPU by default.
- Current solutions: Raise the wired limit with sysctl, quantize, and offload larger jobs to the cloud.
- Why still unsolved: This is a physical memory limit.
- Sources: [llama.cpp discussion #2182](https://github.com/ggml-org/llama.cpp/discussions/2182); [Unsloth requirements](https://unsloth.ai/docs/get-started/beginner-start-here/unsloth-requirements); [gpt-oss-20b model card](https://huggingface.co/openai/gpt-oss-20b); [mlx-lm README](https://github.com/ml-explore/mlx-lm)

**Hardware-matrix testing is heavy for small projects** — severity 3, breadth 4, willingness to pay 2, confidence medium
- Segment: Small OS projects
- Description: Proper multi-hardware OS testing tools such as openQA and LAVA need dedicated infrastructure. Small bootc projects such as the Universal Blue template provide only manual VM testing.
- Evidence: The ublue image-template offers just build-qcow2, run-vm-qcow2 and spawn-vm, with no automated boot tests. openQA tests full installs via VMs and screen checks and is used by openSUSE and Fedora. LAVA needs a multi-component server with workers.
- Quantified impact: Not quantified.
- Current solutions: QEMU serial-console smoke tests in CI, plus community hardware reports.
- Why still unsolved: Real hardware needs power control and console capture, which is costly to run.
- Sources: [ublue-os image-template](https://github.com/ublue-os/image-template); [openQA docs](https://open.qa/docs/); [LAVA docs](https://docs.lavasoftware.org/lava/index.html)

### Key statistics

| Claim | Value | Source type | Date | Source |
|---|---|---|---|---|
| llama.cpp Vulkan on Asahi, Apple M1, Llama 2 7B Q4_0 (no FA) | pp512 38.29 t/s; tg128 12.47 t/s | practitioner | 2025 (scoreboard, undated rows) | https://github.com/ggml-org/llama.cpp/discussions/10879 |
| llama.cpp Metal on macOS, Llama 7B Q4_0: M1 8-core GPU and M1 Pro 16-core GPU | M1: pp 117.96 / tg 14.15; M1 Pro 16c: pp 266.25 / tg 36.41 t/s | practitioner | 2023-11-22 | https://github.com/ggml-org/llama.cpp/discussions/4167 |
| Containerized llama.cpp via krunkit (Vulkan, Venus, MoltenVK) as a share of native Metal | 74–80%; 20.84 t/s on llama3.1:8b, M4 Pro | vendor | 2025-06-05 | https://developers.redhat.com/articles/2025/06/05/how-we-improved-ai-inference-macos-podman-containers |
| Size of the Asahi GPU kernel driver | About 21,000 lines; upstreaming just begun | primary_research | 2026-02-15 | https://asahilinux.org/2026/02/progress-report-6-19/ |
| Asahi GPU stack conformance | OpenGL 4.6, Vulkan 1.4, OpenCL 3.0 conformant; lead developer departed | practitioner | 2025-08-26 | https://alyssarosenzweig.ca/blog/asahi-gpu-part-n.html |
| Idle power saving on a 14-inch M1 Pro from PMP support | About 0.5 W (about 20% lower idle) | primary_research | 2026-04-26 | https://asahilinux.org/2026/04/progress-report-7-0/ |
| Fedora Asahi Remix 43 general availability | 18 March 2026 | official | 2026-03-18 | https://fedoramagazine.org/fedora-asahi-remix-43-is-now-available/ |
| M1 Pro features not yet working on Asahi | Thunderbolt and DP Alt Mode WIP; video encode TBA; Touch ID TBA; Neural Engine out of tree | official | 2026 | https://asahilinux.org/docs/platform/feature-support/m1/ |
| MLX on Linux backends | CPU (mlx[cpu]) and CUDA only (SM 7.5 or later, driver 550.54.14 or later) | official | 2026 | https://ml-explore.github.io/mlx/build/html/install.html |
| Apple container latest release and requirement | v1.5.0, 29 September 2026; requires macOS 26 | official | 2026-09-29 | https://github.com/apple/container/releases |
| Nested virtualization on Apple hosts | Requires Apple M3 or later and macOS 15 or later | official | 2026 | https://raw.githubusercontent.com/lima-vm/lima/master/templates/default.yaml |
| Default macOS GPU memory share on machines with 32 GB or less | About 66.7% (reserve 33.3%), roughly 10.7 GB on 16 GB | practitioner | 2023-2025 | https://github.com/ggml-org/llama.cpp/discussions/2182 |
| GitHub standard runners in public repositories | Free; x64 and arm64 with 4 vCPU and 16 GB RAM | official | 2026 | https://docs.github.com/en/actions/reference/runners/github-hosted-runners |
| GitHub GPU larger runner | Tesla T4 16 GB, 4 vCPU, 28 GB RAM; $0.052/min | official | 2026 | https://docs.github.com/en/billing/reference/actions-runner-pricing |
| AWS g4dn.xlarge (T4 16 GB) price | $0.526/hr on demand; $0.325/hr spot | other | 2026-10-15 | https://instances.vantage.sh/aws/ec2/g4dn.xlarge |
| AWS g6.xlarge (L4 24 GB) price | $0.805/hr on demand; $0.617/hr spot | other | 2026-10-15 | https://instances.vantage.sh/aws/ec2/g6.xlarge |
| RunPod community cloud pod prices | RTX 3090 $0.22, RTX 4090 $0.34, L4 $0.44, A100 80GB $1.19, H100 PCIe $1.99 per hour | vendor | 2026-10 | https://www.runpod.io/pricing |
| Lambda single-GPU on-demand prices | A6000 $1.09, A10 $1.29, A100 40GB $1.99, H100 PCIe $3.29 per hour | vendor | 2026-10 | https://lambda.ai/pricing |
| Latitude.sh bare-metal H100 | From $1.68/hr; hourly billing | vendor | 2026-10 | https://www.latitude.sh/pricing |
| Unsloth QLoRA minimum VRAM | 7B 5 GB; 8B 6 GB; 14B 8.5 GB; 32B 26 GB | vendor | 2026 | https://unsloth.ai/docs/get-started/beginner-start-here/unsloth-requirements |
| gpt-oss-20b memory footprint | 21B total, 3.6B active; runs within 16 GB (MXFP4) | official | 2025-08 | https://huggingface.co/openai/gpt-oss-20b |
| Asahi M3 support | Officially in the installer in September 2026; GPU not performant; sleep and HDMI not working | primary_research | 2026-09-06 | https://asahilinux.org/2026/09/m2-episode-1/ |

### Players

- **Asahi Linux / Fedora Asahi Remix** (Bare-metal Linux for Apple Silicon): Kernel fork, m1n1, U-Boot, installer and conformant Mesa drivers (Asahi, Honeykrisp); FAR 43 and 44 in 2026. Traction: Flagship Apple Silicon Linux distribution; M3 support added in September 2026. Gap: GPU driver still downstream; slower ML than Metal; Thunderbolt and DP Alt Mode WIP; key developers departed; 16K pages
- **libkrun / krunkit (Red Hat, containers org)** (MicroVM monitor with GPU): HVF-based VMs with virtio-gpu (Venus and native context); drop-in for vfkit in Podman machine. Traction: Integrated into Podman machine, RamaLama and Lima 2.0. Gap: Vulkan compute only; experimental in Lima; 74–80% of Metal
- **Lima** (macOS Linux VM manager): qemu, vz and krunkit drivers; Rosetta with vz; templates. Traction: Widely used CNCF project. Gap: krunkit driver experimental; nested virtualization only on M3 or later
- **Apple container / Virtualization.framework** (First-party virtualization): One lightweight VM per container (v1.5.0, September 2026) with Kubernetes commands. Traction: Rapid release cadence in 2026. Gap: Requires macOS 26; no GPU for Linux guests; partial memory ballooning
- **UTM** (VM GUI (QEMU and Apple VZ)): Boots aarch64 ISOs with HVF; VirGL OpenGL on the QEMU backend. Traction: Popular free macOS VM tool. Gap: VirGL experimental; no Vulkan compute path
- **RamaLama** (Containerized model runner): Selects accelerator images (Asahi, Vulkan, CUDA, ROCm); uses krunkit on macOS or MLX with --nocontainer. Traction: Part of the containers ecosystem. Gap: Alpha with breaking changes
- **MLX / mlx-lm (Apple)** (ML framework): Metal-accelerated inference and LoRA on macOS; CPU and CUDA on Linux. Traction: Ollama moved to an MLX engine on Apple Silicon in 2026. Gap: No Linux GPU path on Apple hardware
- **llama.cpp / ggml** (Inference engine): Metal, Vulkan (including on Asahi and MoltenVK or KosmicKrisp), CUDA. Traction: De facto local inference standard. Gap: Vulkan on Asahi is much slower for prompt processing
- **GitHub Actions** (CI): Free x64 and arm64 runners for public repos; KVM on Linux runners; paid T4 GPU runners. Traction: Default CI for open source. Gap: GPU runners paid and cannot boot a custom OS; 14 GB SSD
- **AWS EC2 (g4dn, g6)** (GPU VMs with custom images): T4 or L4 VMs that boot custom AMIs; spot pricing. Traction: Hyperscaler. Gap: Image import and automation effort; egress and EBS costs
- **RunPod / Vast.ai** (Low-cost GPU containers): RTX 3090/4090 from $0.22–0.34/hr. Traction: Popular with hobbyists. Gap: Cannot boot custom kernels; Vast VMs locked to vastai/kvm images
- **Hetzner GEX45 / Latitude.sh** (Bare-metal GPU): Hetzner RTX PRO 4000 Blackwell 24 GB with no minimum contract; Latitude H100 from $1.68/hr. Traction: Used by small AI teams. Gap: Hetzner pricing not captured; custom ISO or iPXE support unconfirmed
- **bootc / bootc-image-builder / Universal Blue** (Image-based OS tooling): Builds an OS from a Containerfile; outputs qcow2, raw, ISO and AMI; GitHub Actions template. Traction: Powers Bazzite, Bluefin and similar projects. Gap: No automated hardware testing in the template
- **openQA / LAVA** (OS test frameworks): Automated install and boot testing on VMs and hardware. Traction: Used by Fedora, openSUSE and Linaro. Gap: Heavy to operate for one person

### Trends
- Asahi has moved from feature work to upstreaming and M3. The GPU stack is conformant but performance work has slowed since the GPU lead left in August 2025.
- GPU-in-microVM is maturing. DRM native context was merged into virglrenderer (Mesa 25.2), and libkrun supports Venus and native context. Lima 2.0 added a krunkit driver with a Venus Vulkan GPU.
- Apple's own stack is pulling ahead for local inference. Ollama adopted MLX on Apple Silicon in 2026, which widens the gap over Linux Vulkan paths.
- Vulkan-on-Metal translation is improving. LunarG's KosmicKrisp is a conformant Vulkan driver on Metal 4, but it needs macOS 26.
- Apple's container tool iterates monthly (v1.2–1.5 between August and September 2026) and now includes Kubernetes commands, but offers no GPU support.
- Image-based OS delivery (bootc plus GitHub Actions plus cosign) is the standard pattern for small distribution projects.
- Free arm64 GitHub runners (since January 2025) make native aarch64 CI free for public repositories.

### Opportunities

- **Make 'one microVM per agent with GPU access' an AgentOS primitive, built on libkrun or virtio-gpu (native context on Asahi and x86, Venus on macOS).** Why now: Native context is upstream in virglrenderer (2025) and krunkit GPU support reached Lima 2.0 (2026). No distribution packages this pattern for agents. Evidence: The Asahi 6.16 report says DRM native context was merged upstream. The libkrun README lists virtio-gpu with Venus and native context. Lima documents a krunkit driver with Venus. Risk: Venus on macOS supports compute only. Nested virtualization is impossible on M1. The libkrun 2.0 API is unstable.
- **A hybrid development mode: run AgentOS in krunkit on macOS and route heavy inference to a host MLX or Metal server over an OpenAI-compatible socket.** Why now: MLX on the host outperforms every Linux path on this hardware, and RamaLama already exposes both paths. Evidence: Metal reaches 266 t/s prompt processing on an M1 Pro against about 38 t/s on Asahi Vulkan (base M1). krunkit reaches 74–80% of Metal. Risk: This is a development crutch, not the shipped product. Users may confuse it with native performance.
- **Publish a reusable 'solo-maintainer OS CI' recipe: free GitHub x64 runners with KVM for QEMU boot smoke tests, plus spot g4dn AMI boot tests triggered per release.** Why now: Small bootc projects lack automated boot tests, and spot T4 capacity costs about $0.33/hr. Evidence: The ublue template offers only manual VM testing. GitHub KVM has been enabled on 2-vCPU runners since April 2024. g4dn.xlarge spot costs $0.325. Risk: Spot capacity interruptions. AWS image import edge cases (ENA and NVMe drivers).

### Contrarian signals
- Bare-metal Linux on the M1 is slower for LLMs than a VM on macOS. Asahi Vulkan prompt processing is about 3x slower than Metal, while krunkit in a VM keeps 74–80% of Metal. A bootable OS on this specific hardware is therefore worse for AI than the host OS.
- Developing mainly on Apple Silicon risks optimizing for a niche. Most production agent and fine-tuning workloads run on x86_64 with NVIDIA, which the M1 cannot test at all.
- Asahi carries a bus-factor risk. Its founder resigned in February 2025 and its GPU lead left in August 2025. The GPU and DCP drivers were still described as too large to upstream in August 2026.
- Secondary sources conflict with primary ones. One 2026 blog lists Vulkan 1.3 and no VP9/HEVC decode, while primary sources report Vulkan 1.4 conformance and mostly reliable AVC/HEVC/VP9 decode. Rely on asahilinux.org and maintainer posts only.
- On M1 hardware, the Lima nested-virtualization statement conflicts with the stated goal of safe agent sandboxing tested in VMs.
- Red Hat's krunkit benchmark was run on an M4 Pro with 48 GB, so the ratio may not hold on an M1 Pro with 16 GB. The 2024 design note set an 8 GB GPU memory cap, and whether it was lifted is unverified.
- The research itself is limited: the shared WebSearch budget ran out after 3 searches, so pricing for Hetzner, GCP and Vultr and some 2026 status claims could not be confirmed.

### Open questions
- Has the 8 GB GPU memory limit in krunkit/Venus been lifted, and how does performance compare on an M1 Pro with 16 GB?
- Has Asahi Vulkan (Honeykrisp) performance for llama.cpp improved in Mesa 25.x or 26.x? Measure pp512 and tg128 on this exact M1 Pro.
- What is Hetzner GEX45's current monthly and hourly price, and can it boot a custom ISO through the VNC installer or KVM console?
- Was Equinix Metal sunset in June 2026 as previously announced? Verification failed (HTTP 403).
- Will Rosetta for Linux VMs survive Apple's announced Rosetta phase-out after macOS 27? This affects x86_64 userspace testing on the Mac.
- Does Fedora Asahi Remix support bootc or atomic variants well enough to share one AgentOS image pipeline across Mac and x86?
- Do GitHub's free ubuntu-24.04-arm runners expose /dev/kvm for aarch64 boot tests?

### Adversarial verification

- **PARTIALLY_CONFIRMED** — Nested virtualization requires Apple M3 or later (macOS 15+), so KVM-based agent sandboxes (Firecracker, Kata, microVM-per-agent) can be tested on the M1 Pro only on bare-metal Asahi or x86 cloud.
  - Corrected: Hardware-accelerated nested virtualization (Virtualization.framework, krunkit, QEMU with HVF) requires M3 or later and macOS 15+, so the M1 Pro cannot run it. A KVM-in-guest stack can still be tested functionally, but very slowly, under QEMU TCG with '-machine virt,virtualization=on'. Realistic sandbox tests need bare-metal Asahi (KVM) or x86 KVM on GitHub x64 runners, where /dev/kvm must first be enabled with a udev rule. GitHub macOS arm64 runners are M1 machines and also lack nested virtualization.
  - Evidence: Apple's documentation for VZGenericPlatformConfiguration.isNestedVirtualizationSupported (macOS 15.0) states that nested virtualization 'is available for Mac with the M3 chip, and later.' The Lima default.yaml says it is supported with vz or krunkit 'on Apple M3 or later (macOS 15 or later)', and with QEMU-HVF only on QEMU 11.1+ on M3+. The QEMU virt documentation describes 'virtualization=on' as emulating the Arm Virtualization Extensions, and its GIC notes tie this option to TCG. TCG emulation is therefore possible but slow, so the claim that this testing is impossible in a VM is overstated. The GitHub changelog of 2024-04-02 confirms KVM on Linux x64 runners (udev rule 99-kvm4all.rules). The GitHub runner reference lists macOS arm64 runners as M1 with 3 CPU and 7 GB and notes 'Nested-virtualization is not supported'. The arm64 Linux runners do not document /dev/kvm.
  - Sources: [Apple: isNestedVirtualizationSupported](https://developer.apple.com/documentation/virtualization/vzgenericplatformconfiguration/isnestedvirtualizationsupported); [Lima default.yaml template](https://raw.githubusercontent.com/lima-vm/lima/master/templates/default.yaml); [QEMU 'virt' machine docs](https://www.qemu.org/docs/master/system/arm/virt.html); [GitHub changelog: hardware-accelerated Android virtualization](https://github.blog/changelog/2024-04-02-github-actions-hardware-accelerated-android-virtualization-now-available/) (2024-04-02); [GitHub-hosted runners reference](https://docs.github.com/en/actions/reference/runners/github-hosted-runners)
- **PARTIALLY_CONFIRMED** — krunkit/libkrun exposes Vulkan compute only (no rendering) via Venus and MoltenVK at 74-80% of native Metal, so a VM on macOS can beat bare-metal Asahi for LLM work.
  - Corrected: Red Hat measured 74-80% of native Metal overall throughput (RamaLama 0.9.0, llama3.2:1b to llama2:13b) on an M4 Pro with 48 GB under macOS 15.5. Native MoltenVK Vulkan on the host already reached 75-77% of Metal, so most of the loss comes from Vulkan-to-Metal translation, not from virtualization. The path is compute only and is confirmed to work on an M1 Pro. The comparison with Asahi does not use like metrics. On a base M1, Asahi token generation is 88% of Metal (12.47/14.15), which beats the krunkit ratio, so the VM wins there only on prompt processing. On Pro-class chips the VM probably wins on both metrics (see the next check). The 8 GB GPU memory cap from 2024 remains unconfirmed either way.
  - Evidence: The Red Hat article (2025-06-05) reports 20.84 t/s for llama3.1:8b, 'between 74% and 80% of the Metal performance', and states that native ggml-vulkan 'performs at 77%' of ggml-metal. Its metric is overall throughput including TTFT, not pp512/tg128. The Podman Desktop GPU docs say the path 'only supports vulkan compute shaders, not rendering / draw', and the example device is 'Virtio-GPU Venus (Apple M1 Pro)'. The Lima krunkit page says '"krunkit" is experimental' and requires Lima >= 2.0 and macOS >= 14, while its caveats section says Ventura, which contradicts itself. The 2024 sinrega post called the 8 GB cap 'a temporary hard limitation'. The current libkrun README states no cap and warns that the libkrun 2.0 main-branch API 'may change further before the first stable release'.
  - Sources: [How we improved AI inference on macOS Podman containers](https://developers.redhat.com/articles/2025/06/05/how-we-improved-ai-inference-macos-podman-containers) (2025-06-05); [Podman Desktop GPU docs](https://podman-desktop.io/docs/podman/gpu); [Lima krunkit driver](https://lima-vm.io/docs/config/vmtype/krunkit/); [Enabling containers to access the GPU on macOS](https://sinrega.org/2024-03-06-enabling-containers-gpu-macos/) (2024-03-06); [libkrun README](https://github.com/containers/libkrun)
- **PARTIALLY_CONFIRMED** — On Asahi, llama.cpp Vulkan prompt processing is about 3x slower and token generation about 12% slower than Metal on macOS (pp512 38.29/tg128 12.47 vs 117.96/14.15, base M1).
  - Corrected: The base-M1 figures are quoted correctly, but they pair an undated 2025 Vulkan commit with a Metal build from November 2023 (b1550), and they understate the gap on Pro-class chips. On an M2 Pro, Asahi Vulkan reaches pp512 62.70 and tg128 20.95, against 294-341 and 37.9-38.9 on Metal in 2023. That is about 5x slower prompt processing and about 46% slower generation. Current Metal builds are 20-33% faster still. For an M1 Pro (200 GB/s), expect Asahi prompt processing to be roughly 5-6x slower and generation roughly 45-60% slower than current Metal. This is an extrapolation and needs a direct measurement.
  - Evidence: The llama.cpp Vulkan scoreboard (#10879) has these Asahi rows, all without dates and without a Mesa version: M1 38.29/12.47 (commit 2370665), M2 Pro 62.70/20.95, M2 Ultra 205.98/34.34. On 2025-01-14 a maintainer commented: 'I think this is slow on the HoneyKrisp side rather than llama.cpp'. Discussion #4167 (commit 8e672ef, 2023-11-21) lists M1 8c at 117.96/14.15, M1 Pro 16c at 266.25/36.41, and M2 Pro 16c/19c at 294.24-341.19 / 37.87-38.86. A re-run of the M2 Ultra 76c at commit c1d0e7a (2026-08-25) shows 1488.81/125.21, against 1238.48/94.27 in 2023. No M1 Pro row exists for Asahi.
  - Sources: [llama.cpp Vulkan performance scoreboard #10879](https://github.com/ggml-org/llama.cpp/discussions/10879); [llama.cpp Apple Silicon performance #4167](https://github.com/ggml-org/llama.cpp/discussions/4167) (2023-11-22)
- **PARTIALLY_CONFIRMED** — Public GitHub repos get free x64 runners with KVM and free arm64 runners; GPU runners (T4, $0.052/min) require a Team or Enterprise plan, are always charged, and cannot boot a custom OS.
  - Corrected: These points are confirmed: free standard runners for public repositories (x64 and arm64, 4 vCPU, 16 GB RAM, 14 GB SSD), KVM on Linux x64 via a udev rule, a GPU runner at $0.052/min ($3.12/hr) with a T4 (16 GB VRAM, 4 vCPU, 28 GB RAM, 176 GB SSD), and larger runners always charged, even for public repositories. These points are not confirmed: current documentation does not name a Team or Enterprise plan requirement (it requires a payment method and a spending limit above $0, set up at organization or enterprise level), and /dev/kvm on arm64 runners is undocumented. The claim 'cannot boot a custom OS' is overstated. Larger runners support custom images for Linux x64, Linux ARM64 and Windows x64 that can start from 'a clean OS'. GPU support for custom images is undocumented, and these images are still runner VMs, not arbitrary bootloader tests.
  - Evidence: The Actions runner pricing page lists linux_4_core_gpu at $0.052, actions_linux at $0.006 and actions_linux_arm at $0.005 per minute. The billing concepts page says 'Larger runners are always charged for, even when used by public repositories'. The larger-runners reference names Enterprise Cloud only for static IP addresses. The custom images how-to says images can start from 'a clean OS' on Linux x64, Linux ARM64 and Windows x64. The runner reference now also lists ubuntu-26.04 and ubuntu-26.04-arm labels.
  - Sources: [GitHub Actions runner pricing](https://docs.github.com/en/billing/reference/actions-runner-pricing); [GitHub Actions billing concepts](https://docs.github.com/en/billing/concepts/product-billing/github-actions); [Larger runners reference](https://docs.github.com/en/actions/reference/runners/larger-runners); [Use custom images (larger runners)](https://docs.github.com/en/actions/how-tos/manage-runners/larger-runners/use-custom-images); [GitHub changelog: KVM on 2-vCPU runners](https://github.blog/changelog/2024-04-02-github-actions-hardware-accelerated-android-virtualization-now-available/) (2024-04-02)
- **CONFIRMED** — Asahi status: Fedora Asahi Remix 43 shipped 18 March 2026 and FAR 44 was targeted for 28 April 2026; the 21,000-line GPU kernel driver is still downstream (upstreaming began February 2026); M3 support entered the installer in September 2026; the GPU lead left in August 2025; Thunderbolt and DP Alt Mode remain WIP on M1 Pro.
  - Corrected: All parts are confirmed, with updates. FAR 44 has since shipped: asahilinux.org/fedora now says the remix is 'based on Fedora Linux 44', although the exact release date was not verified. The M3 post is titled 'M2: Episode 1 (or, Asahi Linux on M3)', and an update dated 2026-09-15 removed the Expert-mode requirement. The Thunderbolt progress in the 7.2 report applies to M3, not M1, so M1 Pro Thunderbolt and DP Alt Mode are still WIP. The claim that the team has moved away from M1 polish is weak: the April 2026 report added PMP idle-power work for the 14-inch M1 Pro, and FAR 43 added 120 Hz support for the MacBook Pro 14/16.
  - Evidence: Fedora Magazine published the FAR 43 article on 2026-03-18. The 7.0 report (2026-04-26) says the team was 'still on track to release Fedora Asahi Remix 44 either on or within a few days of Fedora Linux 44 on 28 April'. The 6.19 report (2026-02-15) says 'The GPU driver is 21,000 lines by itself' and that upstreaming 'has now begun'; it also reports patches falling from 1232 to 858 and lines from 95,000 to 83,000. The 7.1 report (2026-06-30) says GPU initialisation moved into m1n1 ahead of LKML submission. The 7.2 report (2026-08-26) mentions 'enormous work that we cannot upstream yet (e.g. GPU, DCP)'. Rosenzweig wrote on 2025-08-26: 'I am now stepping away from the Apple ecosystem'; her post says the Mesa drivers are upstream and lists Vulkan 1.4 and OpenCL 3.0 conformance. The M1 feature page lists Thunderbolt WIP, DP Alt Mode WIP, video encode TBA, Touch ID TBA and ANE out of tree.
  - Sources: [Fedora Asahi Remix 43 is now available](https://fedoramagazine.org/fedora-asahi-remix-43-is-now-available/) (2026-03-18); [Progress Report: Linux 7.0](https://asahilinux.org/2026/04/progress-report-7-0/) (2026-04-26); [Progress Report: Linux 6.19](https://asahilinux.org/2026/02/progress-report-6-19/) (2026-02-15); [Progress Report: Linux 7.1](https://asahilinux.org/2026/06/progress-report-7-1/) (2026-06-30); [Progress Report: Linux 7.2](https://asahilinux.org/2026/08/progress-report-7-2/) (2026-08-26); [M2: Episode 1 (or, Asahi Linux on M3)](https://asahilinux.org/2026/09/m2-episode-1/) (2026-09-06); [Fedora Asahi Remix page](https://asahilinux.org/fedora/); [Asahi M1 feature support](https://asahilinux.org/docs/platform/feature-support/m1/); [Dissecting the Apple M1 GPU, the end](https://alyssarosenzweig.ca/blog/asahi-gpu-part-n.html) (2025-08-26)

**Missed counter-evidence**
- The Asahi speed gap is much larger on Pro-class chips than on the base M1. An M2 Pro on Asahi reaches tg128 20.95 against about 38 on Metal, so 'about 12% slower generation' should not be used to plan for an M1 Pro.
- Ollama's MLX engine launched in preview on 30 March 2026 and required 'a Mac with more than 32 GB unified memory'. The June 2026 update gives no guidance for 16 GB machines. The 'Ollama moved to MLX' trend therefore does little for a 16 GB M1 Pro, where mlx-lm or llama.cpp Metal remain the practical paths.
- Red Hat's own data shows native MoltenVK Vulkan reaching only 75-77% of Metal on the host. The krunkit penalty is therefore mostly translation overhead, so KosmicKrisp could narrow it. This deserves a direct test rather than relying on the 2025 ratio.
- This session's host reports Darwin 27.2 (macOS 27). If that host is the M1 Pro, it already meets the macOS 26 requirement for Apple container and KosmicKrisp, so the brief overstates those limitations. Apple also says Rosetta will be available as a general-purpose tool only 'through macOS 27', so the window for Rosetta-based x86_64 userspace testing on the Mac is about one year.
- QEMU TCG with virtualization=on can emulate EL2. Nested KVM is therefore possible on the M1, but only for slow functional smoke tests.
- GitHub larger runners support custom images that can start from 'a clean OS' (Linux x64 and ARM64). Whether GPU runners accept them is undocumented. This weakens the absolute claim that they 'cannot boot a custom OS'.
- GitHub has added ubuntu-26.04 and ubuntu-26.04-arm runner labels. The actions/partner-runner-images repository, which supplied the arm64 images, was archived in June 2026. Whether arm64 runners expose /dev/kvm is still undocumented and should be checked with a CI step.
- Asahi 7.1 (June 2026) moved GPU initialisation into m1n1 to prepare the driver for LKML submission. Asahi 7.2 reports DCP nearly at parity with the macOS 13.5 ABI. A new AVD V4L2 decode driver exists but is not shipped by default in FAR.
- The Asahi 7.0 PMP saving (about 0.5 W on a 14-inch M1 Pro) is disabled by default, and base M1 machines use an incompatible PMP variant.
- On 13 October 2026, Ollama and Fly.io launched Sprites, persistent sandboxed Linux VMs for AI agents in the cloud. Established vendors now market VM-per-agent isolation, which weakens the novelty claim behind the 'microVM per agent' opportunity, although no AgentOS-style distribution was found.
- The '66.7% default GPU share' comes from a 2023 community code snippet in llama.cpp discussion #2182, not from Apple documentation. It should be checked on the M1 Pro against llama.cpp's recommendedMaxWorkingSetSize log line.
- AWS GPU instance families are subject to per-family vCPU quotas, which may need an increase request before the first g4dn or g6 launch. This was not verified in this pass and should be checked before relying on per-release automation.

**Reliability note:** Summary: the brief's dated facts largely hold up against primary sources. Its weak points are performance comparisons that mix chips, builds and metrics, and absolute statements ('only', 'cannot', 'requires Team or Enterprise') that the current documentation does not support. Method: the shared WebSearch budget was already used up at the start of this pass (only one search was attempted, and it was refused). All checks therefore rest on about 40 direct fetches of primary pages (Apple, GitHub, Asahi, Fedora Magazine, Lima, QEMU, Red Hat, llama.cpp, Ollama). Some pages showed self-contradictions (Lima's macOS minimum) or undated rows (llama.cpp scoreboard), and the FAR 44 release date was not confirmed. No files were changed. Two questions for the orchestrator: (1) Does the development M1 Pro run macOS 27, as this host reports? (2) Will AgentOS ship Apple Silicon bare-metal images as a supported target, or only use Asahi for development? Each answer changes how much weight the Asahi findings deserve.


---
