# Delivery: an RPM built in CI that the signed image installs, with the same source on COPR

- **Status:** Proposed (Phase 0 draft, revised 2026-10-16 after review; founder
  review in week 3)
- **Date:** 2026-10-16
- **Phase:** 0 (Foundations and kill gate)
- **Deciders:** Founder

## Context

The founder requires a literal bootable operating system, and a headless
appliance is acceptable for v0.1. At the same time, the feasibility research
advised a solo maintainer to build a derived bootc image rather than a
distribution, and to keep the maintainer's own code small. These two pulls meet
in the question of what the unit of delivery is: the image, the package, or both.

Evidence:

- The concept document resolves this after the red team: "The COPR RPM is now the
  delivery unit, the image installs it, and the OCI image works with
  `bootc switch`. Full guarantees are claimed only on the image, which preserves
  the bootable-OS constraint." Source:
  `docs/research/2026-10-16_Tidemark-Concept_v1.md`, delivery and L0 sections.
- The image is justified on grounds the RPM alone cannot meet: it pins a kernel
  for kfuncs that are not a stable ABI (with CI assertions), guarantees a Landlock
  ABI 9 floor, owns every agent launch path, keeps credentials only inside
  daemons, fixes the staging and ledger layout, binds receipts to a signed image
  digest, and gives bootc rollback. Source: the L0 rationale in the concept
  document.
- Stock Fedora already enables bpf and landlock in `CONFIG_LSM`
  (`lockdown,yama,integrity,selinux,bpf,landlock,ipe`), so merely turning LSMs on
  does not justify an image; the kernel pin and the owned launch paths do.
  Verified from the Fedora 45 kernel configuration in
  `src.fedoraproject.org/rpms/kernel` (branch `f45`) on 2026-10-16.
- The concept's L0 layer states that the Containerfile installs "the tidemark
  RPM built in CI". Installing an unpinned package from a COPR repository at
  image build time would place image content outside the provenance that
  cosign signing and the SBOM describe.
- Running the same software on a stock host is valuable for adoption but offers
  weaker guarantees, because the host may not pin the needed kernel or own the
  launch paths. The concept handles this with `doctor`, which writes a guarantee
  level into every receipt.
- v0.1 image outputs are constrained by free CI: qcow2 plus OCI images for x86_64
  and aarch64; ISO, raw and AMI are deferred. Source: concept L0 and CI sections.

## Decision

Deliver Tidemark as an **RPM built in the same CI run as the image**, installed
into the signed bootc image from the local build output, with the identical
source RPM published to COPR for stock-Fedora users.

- The RPM is the unit of software delivery. It carries the daemons, the CLI, the
  BPF object, the verifier, the profiles, the default policy and the systemd
  units and Quadlets.
- CI builds the binary RPM from a source RPM in a pinned Fedora build
  container, records the RPM's SHA-256, and the image build installs that exact
  file after checking the digest. The image never pulls Tidemark from a network
  repository at build time.
- CI publishes the same source RPM to COPR. COPR rebuilds it, so the COPR
  binary is built from identical source but is not byte-identical to the image's
  RPM unless the build is reproducible; receipts and doctor distinguish the two
  through the guarantee level. The fallback, if a same-run build proves
  impractical, is to install a COPR package pinned by full NEVRA and SHA-256
  with `gpgcheck=1`.
- The image is a thin layer: a `Containerfile` `FROM` the pinned fedora-bootc base
  (see `docs/adr/2026-10-16_BaseImage_v1.md`) that installs the RPM, applies
  first-boot defaults, and runs the boot self-test. **Full guarantees are claimed
  only on the signed image.**
- The OCI image is usable with `bootc switch` on bootc-based hosts, so a Fedora
  Atomic user can adopt Tidemark without a reinstall; a package-mode Fedora
  Server host can take over with `bootc install to-existing-root`, which replaces
  the installed system.
- The COPR RPM installs on stock Fedora 45 (Atomic or Server). On such a host,
  guarantees are only those `doctor` reports, and each receipt records the
  guarantee level (`image`, `rpm-on-stock`, or `degraded`).
- v0.1 image artifacts are qcow2 and OCI only, for x86_64 and aarch64 (aarch64
  qcow2 best effort). ISO, raw and AMI are deferred.
- The build pipeline signs the image with cosign keyless signing through GitHub
  OIDC, generates an SBOM with syft, and scans with trivy. It signs the
  multi-architecture index and every per-architecture child (`cosign sign
  --recursive`), because `bootc status`, and therefore the receipt's
  `image_digest`, is expected to report the per-architecture digest (to be
  confirmed in the spike). The specific workflow files are owned by the
  image-pipeline work, not by this ADR.

## Consequences

- The bootable-OS constraint is satisfied honestly: the full guarantees live on an
  image the user can boot and verify, while the RPM lowers the barrier to trying
  the software.
- The maintainer's own image code stays small (a Containerfile plus first-boot
  units), which matches the solo, part-time budget.
- Adoption has three on-ramps: boot the signed image, `bootc switch` to the OCI
  image, or install the RPM on stock Fedora. Each is useful, and each is honest
  about its guarantee level through `doctor`.
- The image's Tidemark content is covered by the same CI provenance, SBOM and
  cosign signature as the rest of the image. Stock-Fedora users depend on COPR's
  rebuild. A compromised Fedora base, COPR or Sigstore is out of scope,
  mitigated only by signature verification; this is stated in the threat model.
- The image pipeline gains an RPM build step before the image build. The
  workflow files are owned by the image-pipeline work.
- Claiming full guarantees only on the image means the README and docs must be
  careful never to imply the RPM-on-stock path is equivalent. The claims checker
  and `doctor` guarantee levels enforce this.
- Deferring ISO, raw and AMI keeps CI inside free-runner disk and time limits;
  users who need those formats are not served in v0.1.

## Alternatives considered

- **Image only, no separable RPM.** Rejected. It would block adoption by users
  who cannot or will not boot a custom image, and it would prevent the
  `bootc switch` and stock-Fedora on-ramps that widen the audience without
  overclaiming.
- **Install the COPR package, unpinned, at image build time.** Rejected: the
  image would contain whatever COPR served at build time, outside the CI
  provenance.
- **RPM only, no owned image.** Rejected. The RPM alone cannot pin the kernel for
  the unstable kfuncs, guarantee the Landlock ABI floor, own every agent launch
  path, or bind receipts to a signed digest. It also fails the founder's
  requirement that the OS be justified, not decorative.
- **A new distribution.** Rejected. Far beyond a solo part-time budget, and the
  feasibility research is explicit that breakout projects configure an existing
  base rather than build a distribution.
- **Ship every image format (ISO, raw, AMI) in v0.1.** Rejected for v0.1 due to
  free-runner disk and time limits; deferred with a trigger.

## References

- `docs/research/2026-10-16_Tidemark-Concept_v1.md` (delivery decision, L0 image
  rationale, doctor guarantee levels, deferred image formats).
- `docs/research/2026-10-16_OS-Feasibility_v1.md` (derived-image recommendation;
  free-runner constraints; `bootc switch`).
- `docs/adr/2026-10-16_BaseImage_v1.md` (the base the image installs the RPM on).
