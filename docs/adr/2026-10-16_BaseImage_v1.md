# Base image: fedora-bootc:45 pinned by digest, with a :44 fallback

- **Status:** Accepted (Phase 0 draft, revised 2026-10-16 after review,
  2026-10-17 in a cross-file consistency pass, and 2026-10-23 for the digest
  re-pin record and the kfunc correction under decisions PD-1, PD-2 and PD-14;
  founder review in week 3)
- **Date:** 2026-10-16
- **Phase:** 0 (Foundations and kill gate)
- **Deciders:** Founder; the 2026-10-23 re-pin was decided by the AI assistant
  under the founder's explicit 2026-10-23 delegation
  ([2026-10-23_Phase0-Founder-Decisions_v1.md](2026-10-23_Phase0-Founder-Decisions_v1.md),
  PD-14) and has had no founder line-by-line review yet

## Context

Tidemark ships as a signed, headless bootc image. The base image must meet four
needs at once:

1. A kernel new enough for the primitives. Landlock ABI 9 (pathname Unix socket
   scoping) requires Linux 7.1, and the read-time labeller depends on the fs
   kfuncs `bpf_get_file_xattr` and `bpf_path_d_path` and on the sleepable BPF
   LSM hooks `file_open`, `file_receive` and `mmap_file`.
   `bpf_real_data_inode` first appears in Linux 7.2 and was on this list until
   the Phase 0 spike: decision PD-1 removed the inode-storage cache it served,
   so it is no longer a dependency (see
   `docs/adr/2026-10-16_Labeller-BPF-LSM_v1.md`). The v0.1 floor stays at
   **Linux 7.2 or later with Landlock ABI 9 or later** as a policy choice, not
   as a kfunc constraint (PD-2).
2. A multi-architecture base, because the founder develops on an Apple M1 Pro
   (aarch64) and gates releases on free x86_64 KVM CI runners.
3. A base a solo, part-time maintainer can track without running a distribution.
4. A base that lets an existing Fedora Atomic user adopt Tidemark with
   `bootc switch`.

Evidence verified against primary sources on 2026-10-16:

- **fedora-bootc:45 is a multi-architecture manifest list.** The quay.io API tag
  record for `:45` is manifest digest
  `sha256:86cc9d98ef210b71a6c4eeec4b7d04e7fef0a4415184c2322a14c879df189f45`
  (last modified 2026-10-16), an OCI image index with amd64, arm64, ppc64le and
  s390x children. The amd64 child is
  `sha256:4e5c285af0fbf9698525a3c5ff6b6dbfc5c1d6102b0f1f46a48b2c437ffe533c`
  and the arm64 child is
  `sha256:78d94543f47b1b147af00a83c5ecd79023db5073cb0c90fc74f4d00cae9a1ed4`.
  Source: `https://quay.io/api/v1/repository/fedora/fedora-bootc/tag/?specificTag=45`
  and the quay.io registry v2 manifest API.
- **The :45 kernel is 7.2.9.** The amd64 child config label `ostree.linux` reads
  `7.2.9-300.fc45.x86_64` and the arm64 child reads `7.2.9-300.fc45.aarch64`
  (image version `45.20261016.0`). A 7.2 kernel meets the 7.2 floor, provides
  Landlock ABI 10 (`landlock_abi_version` is 8, 9 and 10 at the v7.0, v7.1 and
  v7.2 tags), and carries both kfuncs the labeller now depends on
  (`bpf_get_file_xattr`, `bpf_path_d_path`).
- **The :45 tag is not new.** The quay.io tag history for `:45` (including
  inactive tags) goes back to 2026-06-01; the tag is rebuilt frequently, which
  is why the digest and its resolution date are recorded.
- **fedora-bootc:44 is also multi-arch and current.** Tag `:44` is manifest
  `sha256:9db34195bc0b6be43e17c20ca840be5bc153093c94b62f6792a76fb199eeaced`
  (2026-10-16), kernel `7.2.10-200.fc44.x86_64`. Fedora 44 therefore also meets
  the kernel floor, which makes it a viable fallback build.
- **The Fedora kernel enables the needed LSMs.** The shipped Fedora 45 kernel
  configuration sets `CONFIG_LSM="lockdown,yama,integrity,selinux,bpf,landlock,ipe"`,
  `CONFIG_BPF_LSM=y`, `CONFIG_SECURITY_LANDLOCK=y` and `CONFIG_DEBUG_INFO_BTF=y`
  for both x86_64 and aarch64. bpf and landlock are in the active LSM list
  without a boot-parameter change. Source: `src.fedoraproject.org/rpms/kernel`,
  branch `f45`, files `kernel-x86_64-fedora.config` and
  `kernel-aarch64-fedora.config`, checked 2026-10-16.
- **Universal Blue bases are amd64-only.** The OS-feasibility research recorded
  that the ublue-os images (base-main, bazzite, bluefin, aurora) are amd64 only,
  and that the reference CentOS Stream 10 model (Bluefin GDX/LTS) was restructured
  and archived during 2026. Only `quay.io/fedora/fedora-bootc` provides arm64 for
  an M1-based workflow. Source:
  `docs/research/2026-10-16_OS-Feasibility_v1.md`.

## Decision

Build the image `FROM quay.io/fedora/fedora-bootc:45`, pinned by the digest
`sha256:86cc9d98ef210b71a6c4eeec4b7d04e7fef0a4415184c2322a14c879df189f45`
resolved on 2026-10-16.

### Re-pin record

A pinned digest identifies a compose; it does not keep it retrievable. quay.io
drops a compose's manifests once the daily rebuild has moved the tag, so a pin
must be refreshed by a recorded change.

| Date (UTC) | Tag | Digest that was pinned | Still retrievable? | Replacement index digest | Compose / amd64 kernel |
|---|---|---|---|---|---|
| 2026-10-16 | 45 | `sha256:86cc9d98…df189f45` | no (`manifest unknown` on 2026-10-23) | `sha256:4c71411136f85da97cfc43d55291e13eec8cb2b2ff25ddf5ff5709ed093929da` | 45.20261023.0 / `7.2.11-300.fc45.x86_64` |
| 2026-10-16 | 44 | `sha256:9db34195…9eeaced` | no (`manifest unknown` on 2026-10-23) | `sha256:03808453feb5369ed6b074421117b56734f3a6dff8e25d9e0ddff0e48e4d6d2e` | 44.20261023.0 / `7.2.11-200.fc44.x86_64` |

Both replacements were resolved on 2026-10-23 at about 15:20 UTC through the
Docker Registry v2 API (`GET /v2/fedora/fedora-bootc/manifests/<tag>`, reading
the `Docker-Content-Digest` response header for the OCI image index), with the
per-architecture children read from that index and the compose version and
kernel from each child's image config labels. Recorded under PD-14; the first
build against the 2026-10-16 pins failed in CI (build-image run 39965540155)
before this was known.

- Pin by digest, not by the floating `:45` tag. Renovate proposes digest bumps,
  which merge only after the boot test and the read matrix pass.
- CI checks reachability before building:
  `.github/actions/check-base-digest` reads the pin out of the Containerfile,
  asks the registry whether that digest still resolves, prints the tag's current
  digest, and fails with that name if the pin is gone. It never rewrites a pin.
  Auto re-pinning was rejected because an image labelled with one digest must not
  be built from another compose's content.
- Keep a parallel fallback build `FROM quay.io/fedora/fedora-bootc:44`
  (digest `sha256:9db34195bc0b6be43e17c20ca840be5bc153093c94b62f6792a76fb199eeaced`,
  2026-10-16, re-pinned 2026-10-23 as above) in `image/Containerfile.f44`,
  maintained until launch, so a regression in the 45 line does not block the
  project.
- The image build file paths are owned by the `image_pipeline` work, not by this
  ADR; this ADR fixes only the base and the pinning policy.
- A boot self-test asserts, at first boot, a kernel of 7.2 or later, Landlock
  ABI >= 9, that the active LSM list contains bpf and landlock, and that the
  required kfuncs are present. C2 decided on 2026-10-23 (PD-1): the cache is
  out, so `bpf_real_data_inode` is not in the required list and the self-test
  checks `bpf_get_file_xattr` and `bpf_path_d_path`. From Phase 1, sessions
  refuse to start if any assertion fails; in Phase 0 the self-test only reports.
  The Phase 0 self-test script (`image/selftest/tidemark-selftest.sh`, owned by
  the image pipeline) already checks the 7.2 floor and the two kfunc names.

## Consequences

- The kernel is pinned to a known-good point, which matters because the fs
  kfuncs are not a stable ABI. CI assertions catch a base bump that drops or
  changes them.
- Both shipped architectures come from one base, so the M1 development loop and
  the x86_64 CI gate use the same kernel package line.
- Users of bootc-based Fedora Atomic hosts can adopt Tidemark with
  `bootc switch`, because the base is the stock Fedora bootc image. A
  package-mode Fedora Server host cannot run `bootc switch`; it can take over
  with `bootc install to-existing-root`, which replaces the installed system.
- The `:45` tag is rebuilt frequently and the kfunc surface is unstable, so the
  project carries base churn risk. The `:44` fallback and digest pinning reduce
  it.
- The image build cannot be exercised natively on the M1 for x86_64; x86_64
  builds and the gating boot test run in CI. This is accepted and recorded in the
  development-environment notes.
- A future move to the bootc composefs or sealed-image backend requires a
  reinstall, because the composefs backend has no in-place move from ostree
  (verified in the OS-feasibility research). Measured boot stays deferred to
  Phase 7, and the ESP layout is chosen to allow a later switch.

## Alternatives considered

- **CentOS Stream 10 (the Bluefin GDX/LTS model).** Rejected. The reference
  images that made this attractive were restructured and the LTS line archived
  during 2026, and the countme data shows almost no users on the CentOS path.
  Its 6.12-based kernel is below the 7.2 floor, so the labeller's kfunc and
  Landlock ABI 9 would depend on vendor backports that a solo maintainer cannot
  track.
- **A Universal Blue base (bazzite/bluefin).** Rejected for v0.1 because these
  are amd64 only, which breaks the aarch64 development loop on the M1.
- **NixOS or mkosi/ParticleOS.** Rejected. NixOS has no cached CUDA for stable
  branches and a steep learning curve for a solo maintainer; ParticleOS is
  explicitly unstable and ships prebuilt images for x86-64 only. Neither fits a
  part-time schedule, though mkosi/ParticleOS remain references for the Phase 7
  verified-boot design.
- **Pin by the floating `:45` tag.** Rejected. A floating tag silently changes
  the kernel and the kfunc surface; digest pinning plus CI assertions is required
  for a project whose core depends on an unstable kfunc ABI.

## References

- `docs/research/2026-10-16_OS-Feasibility_v1.md` (bootc tooling, Universal Blue
  amd64-only finding, Bluefin restructuring, composefs reinstall cost).
- `docs/research/2026-10-16_Tidemark-Concept_v1.md` (L0 image layer, base-image
  rationale).
- quay.io registry API and v2 manifest API for `fedora/fedora-bootc:45` and
  `:44` (digests and `ostree.linux` labels, checked 2026-10-16).
- Linux `security/landlock/syscalls.c` at tags v7.0, v7.1, v7.2
  (`landlock_abi_version` = 8, 9, 10 respectively).
- `src.fedoraproject.org/rpms/kernel`, branch `f45`,
  `kernel-{x86_64,aarch64}-fedora.config` (`CONFIG_LSM`, `CONFIG_BPF_LSM`,
  `CONFIG_SECURITY_LANDLOCK`, `CONFIG_DEBUG_INFO_BTF`), checked 2026-10-16.
- quay.io tag history API for `fedora/fedora-bootc` tag `45` with inactive tags
  (earliest entry 2026-06-01), checked 2026-10-16.
- `fs/bpf_fs_kfuncs.c` at v7.1 and v7.2 (`bpf_real_data_inode` absent and
  present).
