# Trademark and Name Use Policy

> **Status: Phase 0 draft, 16 October 2026. Provisional, pending trademark clearance.**
> "AgentOS Tidemark" and "Tidemark" are working names. No trademark has been searched, cleared,
> registered or applied for. The name may change before the v0.1 preview launch, planned for April
> 2027. This document is not legal advice.

## Why this policy exists

The [Apache License, Version 2.0](LICENSE) grants broad rights to the code, but Section 6 of that
licence grants no permission to use the project's names or marks. This policy fills that gap.

The name matters for security. Tidemark's guarantees are claimed only on the project's **signed
image**, when `tidemark doctor` passes. A modified build that uses the name could mislead users into
trusting guarantees that the build does not provide.

## Current status of the name

- "AgentOS" is a crowded term already used by several unrelated projects and products.
- "Tidemark" needs a trademark search before public launch. The founder plans to self-search the
  USPTO, EUIPO and IP Australia databases first, and to consider professional clearance before
  v1.0.
- No domain name will be bought until the search is complete.
- There is no logo yet.

If the name changes, this policy applies to the new name.

## What you may do without asking

- Refer to the project truthfully by name, for example "a fork of Tidemark", "works with
  Tidemark", "compared against Tidemark" or "a Tidemark benchmark result".
- Redistribute the official, unmodified signed artifacts (image, RPM and OCI images) under the
  project's name, provided the signatures still verify.
- Use the name in articles, talks, research papers, tutorials and benchmark reports.
- Publish independent "Bring Your Own Canary" verifications that name the project, as long as you
  state the image digest you verified and the conditions of your run.
- Package unmodified upstream source releases for a Linux distribution under the project's name,
  with only the packaging changes the distribution requires, provided the package does not claim
  the signed image's guarantees.

## What requires a different name

- **Modified builds.** If you change the code, the policy defaults, the image contents or the
  build process and distribute the result, do not call it "Tidemark" or "AgentOS Tidemark". Say
  instead that it is "based on Tidemark". Distribution packages that meet the packaging condition
  above are not modified builds for this purpose.
- **Forks** must use a different name.
- **Products and services** must not use the name in a way that suggests they are the project or
  are endorsed by it.

## What is never allowed

- Implying that the project endorses, certifies or audits your product or service.
- Claiming Tidemark guarantees for a build that is not the signed image, or for a host where
  `tidemark doctor` does not pass.
- Using the name in a domain name, package name or account name that could be confused with the
  project's official channels.

## Questions and permission requests

Contact: `[TRADEMARK CONTACT: to be set by the founder]`. Until a contact is set, open an issue
labelled `trademark` in the project repository once it is public.

This policy will be reviewed after the trademark search, and again if the project moves to a neutral
home (see [GOVERNANCE.md](GOVERNANCE.md)).
