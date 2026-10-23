# AgentOS Tidemark Image Pipeline (v1, 2026-10-16)

- **Status:** Phase 0 draft, revised on 2026-10-16 after review. No workflow in this document has run yet, because the GitHub repository does not exist. Every statement about CI behaviour is a design intent until the first green run is linked here.
- **Phase:** 0 (weeks 1 to 3, 19 October to 8 November 2026), checklist items E1 to E4, T1, T2, P1, P2 and acceptance criterion AC6.
- **Owner:** Founder (reviews every change to the workflows and the self-test).
- **Related records:** `docs/adr/2026-10-16_BaseImage_v1.md` (authoritative base-image decision), `docs/2026-10-16_Phase0-Checklist_v1.md`, `docs/research/2026-10-16_Tidemark-Concept_v1.md` (layer L0 and L8).

## Executive summary

- The image is `quay.io/fedora/fedora-bootc:45`, pinned by digest, plus two Phase 0 units: a first-boot banner and a boot self-test. It contains **no Tidemark enforcement**: no labels, no gateway, no broker, no staging and no receipts. A Fedora 44 fallback is built the same way.
- One run of `build-image.yml` carries each image from build to publication. It builds each x86_64 variant once, saves it as an OCI archive, scans that archive with no registry credentials, and boots that same archive twice through the reusable `boot-test.yml`. On `main`, and only after both builds and both boot tests pass, it pushes that same archive unchanged, signs it, verifies the signature and attaches a signed SBOM.
- Both the Fedora 45 image (`:main`) and the Fedora 44 fallback (`:main-f44`) are pushed, signed and verified, as checklist items E2 and T1 require.
- **Pinning:** actions are pinned by commit SHA and container images by digest; actionlint and the Python wheels are pinned by SHA-256. Not pinned: the Ubuntu packages (QEMU, OVMF, and podman on arm64), the trivy vulnerability database, and the lychee binary that lychee-action downloads.
- The self-test now requires BTF, the three kfunc names and a Linux 7.2 floor (the recommended answer to checklist decision FD-G, pending the founder's confirmation).
- Still open: checklist items E3 and T2 need the kernel probe to run inside the booted image (section 6.2).

## 1. Scope and what is not claimed

Once it runs, this pipeline will prove four things. The image builds from pinned inputs. The image that boots in CI is the image that is published. The published x86_64 images carry verifiable keyless signatures and an SBOM attestation. The image boots to a self-test that reports the kernel and filesystem features later phases need.

It does not claim that the image protects anything. The self-test reports; it does not refuse sessions, because there are no sessions yet. The concept's "not claimed" list applies in full: no prevention of prompt injection, no resistance to kernel or root compromise, no protection of unlabelled data, and no proof against the machine owner before TPM attestation.

## 2. Files

| Path | Purpose |
|---|---|
| `image/Containerfile` | Main image on `fedora-bootc:45@sha256:…`. It copies and enables the Phase 0 units, turns off telemetry, masks automatic updates and ends with `bootc container lint`. |
| `image/Containerfile.f44` | Fallback build on `fedora-bootc:44@sha256:…`; identical build steps. Kept until launch. |
| `image/selftest/tidemark-selftest.sh` | Boot self-test. Installed as `/usr/libexec/tidemark/tidemark-selftest.sh`. |
| `image/systemd/tidemark-selftest.service` | Runs the self-test on every boot after the first. |
| `image/firstboot/tidemark-firstboot.sh` | Prints the Phase 0 banner, runs the self-test once, writes and flushes a stamp. |
| `image/firstboot/tidemark-firstboot.service` | Runs the first-boot script once. |
| `image/bib/config.toml` | bootc-image-builder customizations: serial console kernel arguments only; no users or keys. |
| `.github/workflows/build-image.yml` | Build, scan, boot test (through the reusable workflow), push, sign, verify, attest. |
| `.github/workflows/boot-test.yml` | Reusable workflow: qcow2 build and a two-boot QEMU/KVM serial-console test of one archive. |
| `.github/workflows/lint.yml` | actionlint, shellcheck, policy schema validation, offline Markdown link check. |
| `.github/workflows/scorecard.yml` | OpenSSF Scorecard. |
| `renovate.json` | Digest and action updates (justification in its `description` field and in section 9). |

The Go module path and the documentation URLs in the unit files were placeholders until 2026-10-23, when the repository was created as `github.com/BitecodesHub/tidemark` (Phase 0 decisions ADR PD-10); they now point there. The image name in CI is derived at run time as `ghcr.io/<owner>/tidemark` with the owner lowercased, so it needs no edit.

## 3. Base image pinning

Resolved on 2026-10-16 at about 16:16 UTC through `https://quay.io/api/v1/repository/fedora/fedora-bootc/tag/?specificTag=<tag>` and the quay.io registry v2 manifest API.

| Tag | Pinned index digest | amd64 child | arm64 child | Version | amd64 kernel |
|---|---|---|---|---|---|
| 45 | `sha256:86cc9d98ef210b71a6c4eeec4b7d04e7fef0a4415184c2322a14c879df189f45` | `sha256:4e5c285a…ffe533c` | `sha256:78d94543…9a1ed4` | 45.20261016.0 | 7.2.9-300.fc45 |
| 44 | `sha256:9db34195bc0b6be43e17c20ca840be5bc153093c94b62f6792a76fb199eeaced` | `sha256:f2d97f1e…0136bc` | `sha256:3588110b…c0287` | 44.20261016.0 | 7.2.10-200.fc44 |

Both indexes also list ppc64le and s390x children, which Tidemark does not build. The tags are rebuilt daily, so the digest is the identity. Renovate proposes one grouped digest refresh per week. Each refresh runs both builds and both boot tests before a human merges it. Renovate keeps `Containerfile` on tag 45 and `Containerfile.f44` on tag 44; a move to Fedora 46 needs a new base-image ADR.

## 4. Image contents (Phase 0)

1. **Tool assertions.** The build fails if the base image stops shipping a tool the self-test uses (`python3`, `setfattr`, `getfattr`, `chattr`, `lsattr`, `unshare`, `mount`, `stat`, `grep`, `mktemp`). Phase 0 installs no packages.
2. **Units.** The two scripts go to `/usr/libexec/tidemark/` and the two units to `/usr/lib/systemd/system/`; both units are enabled with `systemctl enable` at build time. Nothing is written to `/var` in the image, as bootc expects; `/var/lib/tidemark` is created at run time by `StateDirectory=`.
3. **No telemetry.** The rpm-ostree countme timer is masked, `countme=1` is changed to `countme=0` in every DNF repository file, and NetworkManager connectivity checks are disabled through a drop-in in `/usr/lib/NetworkManager/conf.d/`.
4. **No unverified automatic updates.** `bootc-fetch-apply-updates.timer` is masked. Without a client-side signature check (section 11), the timer would download whatever the tracked tag points to, stage it and reboot into it. Until `policy.json` enforcement lands in Phase 1, updating is a manual step: verify the new digest with `cosign verify` (section 7), then run `sudo bootc switch ghcr.io/<owner>/tidemark@sha256:<digest>` and reboot. Switching to a digest, rather than upgrading from a tag, avoids pulling content other than the digest that was verified. That `bootc switch` accepts a digest reference is expected but not verified.
5. **Phase 1 placeholder.** A commented block marks where the Tidemark RPM will be installed. Per `docs/adr/2026-10-16_Delivery-RPM-Plus-Image_v1.md`, the image installs the RPM that an earlier job of the same CI run built, after checking its SHA-256, and never pulls Tidemark from COPR or another network repository at build time. That ADR also requires an RPM build job before the image build and, once a multi-architecture index is published, signing with `cosign sign --recursive`; neither exists in Phase 0, which publishes single-architecture x86_64 images only.
6. **Lint.** `bootc container lint` runs last; its errors fail the build, and its warnings are printed.

`fwupd-refresh.timer` (firmware metadata download) stays at the base default, pending a founder decision (section 11).

### First boot and every boot

- On the first boot, `tidemark-selftest.service` is skipped, because its condition requires `/var/lib/tidemark/firstboot.done`. `tidemark-firstboot.service`, ordered after it, prints the banner, runs the self-test, writes the stamp and flushes the file system with `sync -f`, so that the stamp survives an unclean power-off.
- On every later boot, the first-boot unit is skipped and `tidemark-selftest.service` runs the self-test.
- Each boot should therefore print exactly one result line. The boot test is designed to check this on two consecutive boots of the same disk; it has not run yet.
- Both units run as root with `PrivateNetwork=yes`, `PrivateTmp=yes` and `ProtectHome=yes`, and send their output to the journal and the console.

## 5. Boot self-test

The last line of output is always one machine-readable line. For example (illustrative values, not a recorded run):

```
TIDEMARK_SELFTEST v=1 result=PASS phase=0 failed=none uptime_s=14.62 kernel=7.2.9-300.fc45.x86_64 arch=x86_64 kernel_floor=ok lsm=lockdown,capability,yama,selinux,bpf,landlock,ipe lsm_bpf=ok lsm_landlock=ok landlock_abi=10 landlock=ok cgroup2=ok ... btf=ok kfunc_bpf_get_file_xattr=ok kfunc_bpf_path_d_path=ok selinux=enforcing
```

The same pairs are written to `/run/tidemark/selftest.result`. Exit status: 0 for PASS, 1 for FAIL, 2 for a usage error.

| Key | Gating in Phase 0 | How it is checked |
|---|---|---|
| `kernel_floor` | Yes | `uname -r` is 7.2 or later. Landlock ABI 9 first ships in Linux 7.1. The 7.2 floor is a policy choice, not a kfunc dependency: the labeller keeps no inode-keyed cache (decisions ADR PD-1) and does not call `bpf_real_data_inode`. 7.2 remains the floor because Phase 0 has no evidence on 7.1; 7.1 support is Phase 6. |
| `lsm_bpf`, `lsm_landlock` | Yes | Both names appear in `/sys/kernel/security/lsm`. |
| `landlock` (with `landlock_abi`) | Yes | `landlock_create_ruleset(NULL, 0, LANDLOCK_CREATE_RULESET_VERSION)` returns 9 or more. The kernel exposes the ABI only through this query, so the script issues raw syscall 444 through Python `ctypes` (444 on x86_64 and on the asm-generic table used by aarch64). |
| `cgroup2` | Yes | `/sys/fs/cgroup` has file-system type `cgroup2fs` and a readable `cgroup.controllers`. |
| `overlayfs` | Yes | An overlay is mounted in a private mount namespace with its upper and work directories on `/var`; a write through the merged view must land in the upper directory. |
| `var_xattr` | Yes | A `user.*` attribute is set and read back on a file in `/var`. |
| `var_append` | Yes | `chattr +a` succeeds, `lsattr` reports the flag, a truncating open is refused and an append succeeds. |
| `btf` | Yes | `/sys/kernel/btf/vmlinux` is readable. |
| `kfunc_bpf_get_file_xattr`, `kfunc_bpf_path_d_path` | Yes | Each name appears as an exact string in `/sys/kernel/btf/vmlinux`. `bpf_real_data_inode` is no longer required: the labeller keeps no inode-keyed cache (decisions ADR PD-1, 2026-10-23), so the kfunc is unused. This is a proxy: it shows that the kernel carries BTF for the function, not that a sleepable BPF LSM program may call it. Confirming each kfunc by loading a program is the job of the Phase 0 kernel probe and, from Phase 1, of tidemarkd. |
| `uptime_s` | No | Seconds since boot when the self-test started, from `/proc/uptime`; the guest's own figure for target P2. |
| `xattr_security_bpf` | No | Root sets and reads a `security.bpf.*` attribute; recorded for the kernel spike. |
| `cgroup2_controllers`, `cgroup2_freeze`, `var_fstype`, `selinux`, `arch`, `lsm` | No | Recorded for diagnosis. |

The probes run in a private `run.XXXXXX` directory, created with `mktemp` under `TIDEMARK_SELFTEST_DIR` (default `/var/lib/tidemark/selftest`) and removed afterwards. The script therefore never changes the permissions of a caller-supplied directory. The script also runs outside the image, for example in the Lima Fedora 45 Cloud VM used for the kernel spike: `sudo bash image/selftest/tidemark-selftest.sh`. Keep `TIDEMARK_SELFTEST_DIR` on the VM's own disk, not on a virtiofs or 9p share, or the `/var` results describe the share instead.

## 6. CI workflows

`build-image.yml`, `lint.yml`, `scorecard.yml` and the reusable `boot-test.yml` set `permissions: {}` at the top and grant the minimum per job. Every `actions/checkout` step sets `persist-credentials: false`.

### 6.1 build-image.yml

| Job | Runs on | Trigger | Permissions | Gating |
|---|---|---|---|---|
| `detect image changes (build-image)` | ubuntu-24.04 | all | `contents: read` | Yes |
| `build (x86_64, fedora-45)`, `build (x86_64, fedora-44)` | ubuntu-24.04 | all | `contents: read` | Yes |
| `boot test (fedora-45)`, `boot test (fedora-44)` (reusable `boot-test.yml`) | ubuntu-24.04 | all | `contents: read` | Yes |
| `sign and publish (x86_64, fedora-45)`, `sign and publish (x86_64, fedora-44)` | ubuntu-24.04 | push to `main`, manual run on `main` | `packages: write`, `id-token: write` only | Yes |
| `build (aarch64, non-gating)` | ubuntu-24.04-arm | all | `contents: read` | No (`continue-on-error`) |

Triggers: pull requests (always started, so that required checks report), pushes to `main` that touch `image/` or either workflow file, a weekly schedule (Monday 04:17 UTC) and manual runs. On pull requests, the change-detection job compares the pull request with its base. The image jobs are skipped (which GitHub counts as passing) when nothing under `image/`, `build-image.yml` or `boot-test.yml` changed. If the diff cannot be computed, the image jobs run. Scheduled runs build and boot-test but never publish; they exist to notice changes in the runner image, QEMU or OVMF.

**Build job (one per variant, no write permissions).**

1. Builds with rootful `podman build` and saves the result with `podman save --format oci-archive`. It records the image ID, which is the config digest.
2. Unpacks the archive into an OCI layout directory, because trivy `--input` does not read `oci-archive` tars.
3. Runs syft v1.54.1 from its digest-pinned image with `--network none` and no credentials, producing a package-level SPDX JSON SBOM. A missing SBOM fails the job.
4. Runs trivy 0.75.0 from its digest-pinned image with `--exit-code 0` and no credentials. Trivy failures and a missing or empty report become `::warning::` annotations and never fail the job.
5. Uploads the archive and image ID as `tidemark-image-<variant>` (one-day retention), and the SBOM and trivy reports as `tidemark-reports-<variant>` (30 days).

**Publish job (one per variant, `main` only, after every build and boot-test job passes).**

1. Downloads the boot-tested archive and the SBOM. `download-artifact` fails on a digest mismatch by default.
2. Pushes the archive with `skopeo copy --preserve-digests` to `ghcr.io/<owner>/tidemark:main` (Fedora 45) or `:main-f44` (Fedora 44). `--preserve-digests` refuses any manifest change, so the pushed digest is the archive's own manifest digest.
3. Checks that the pushed manifest's config digest equals the image ID that the build job recorded and the boot test loaded. It then adds the tag `sha-<commit>` or `sha-<commit>-f44` by digest.
4. Signs `ghcr.io/<owner>/tidemark@<digest>` twice with cosign v3.0.6, keyless through GitHub OIDC. The first signature is a Sigstore bundle (the cosign v3 default, stored through the OCI referrers API). The second is in the legacy format (`--new-bundle-format=false --use-signing-config=false`), stored as the `sha256-<digest>.sig` tag that podman and bootc read (section 11).
5. Runs `cosign verify` against the exact workflow identity for both formats (acceptance criterion AC6) and checks that the `.sig` tag exists.
6. Attaches the SBOM with `cosign attest --type spdxjson` and verifies it with `cosign verify-attestation`.
7. Uploads the verification output as `tidemark-signing-<variant>`.

**Deviation to confirm.** The Fedora 44 fallback was non-gating in the first draft. It now gates on pull requests and on `main`, because E2 and T1 require both builds to be signed and verified. If the fallback breaks, it blocks publishing of the Fedora 45 image as well, because a job that needs a matrix job waits for every leg of it. If the founder prefers a non-gating fallback, remove it from the three matrices and amend E2 and T1.

### 6.2 boot-test.yml (reusable)

`build-image.yml` calls this workflow once per variant inside the same run. It boots the exact archive that the build job produced and that the publish job later pushes.

1. Enables KVM with the udev rule from the GitHub changelog of 2 April 2024 (`KERNEL=="kvm", GROUP="kvm", MODE="0666", OPTIONS+="static_node=kvm"`) and fails at once if `/dev/kvm` is not writable.
2. Installs `qemu-system-x86`, `qemu-utils` and `ovmf` from the Ubuntu archive (not pinned).
3. Downloads `tidemark-image-<variant>`, loads it with `podman pull oci-archive:…`, and fails unless the loaded image ID equals the recorded one.
4. Runs the digest-pinned bootc-image-builder container with `--privileged`, `--type qcow2`, `--rootfs xfs` and `image/bib/config.toml` mounted at `/config.toml`. The Fedora bootc image sets no default root file system, so `--rootfs` is required.
5. Creates one copy-on-write overlay on top of the qcow2 and boots it with UEFI (OVMF), KVM, 2 vCPUs, 4 GiB, no network (`-nic none`), a serial console written to a log and a QEMU monitor on a Unix socket.
6. **Boot 1.** Polls the log every 2 seconds for up to 600 seconds. A result line counts only once the line is complete up to the space after the `failed=` field, so a line that is still being written cannot end the wait early. After the result line, it keeps polling for up to 30 seconds for systemd's `Reached target … Multi-User System` message. It then sends `system_powerdown` through the monitor and waits up to 120 seconds for an ACPI power-off before it falls back to SIGTERM and SIGKILL.
7. **Boot 2.** Boots the same overlay again and waits for the result line.
8. After each stop, it reads the log once more, so that output written after the last poll counts. The job passes only if:
   - boot 1 printed the first-boot banner and exactly one `TIDEMARK_SELFTEST` line, with `result=PASS`;
   - boot 2 printed exactly one such line, with `result=PASS`, and no banner.
9. Writes a summary table and uploads both serial logs as `tidemark-boot-test-serial-logs-<variant>` on every run, pass or fail. The table shows when the result line and `multi-user.target` were seen (host time since QEMU started), the guest's own `uptime_s`, the result-line counts, and whether each power-off was clean.

**Target P2 (multi-user in 90 seconds or less).** The figure is informational. systemd reaches `multi-user.target` only after both oneshot units finish, so the figure includes the self-test's run time. It has 2-second resolution, and it appears only if systemd prints status messages on the serial console. The guest's `uptime_s` is the second, independent figure.

**Checklist items E3 and T2 stay open.** Both require a boot smoke test "that runs the probe", meaning the D4 loader. The boot test runs only the self-test, which covers T4's LSM-list and Landlock assertions but not the probe. Proposed plan, not implemented:

1. A test-only layer: `FROM localhost/tidemark:boot-test`, then `COPY` in the probe binary and BPF object built by an earlier job of the same run, together with the scenario scripts (`spike/tests/`). The scenarios also call `getfattr` and `setpriv` (packages `attr` and `util-linux`), which the layer installs if the base image lacks them.
2. A oneshot unit that runs the probe over its 12 scenarios and prints a `TIDEMARK_PROBE` line to the console carrying the probe's `SUMMARY:` result (`SUMMARY: ALL PASS (12/12)` when every scenario passes).
3. bootc-image-builder builds the qcow2 from that layer. The layer is never published.

The founder should decide whether the self-test is acceptable as the E3 probe or whether the D4 probe must run in the booted image.

### 6.3 lint.yml

- **actionlint** 1.7.12 is downloaded and checked against a recorded SHA-256 before it runs. It also runs shellcheck on every `run:` block and validates the call to the local reusable workflow.
- **shellcheck** runs on every tracked `*.sh` and `*.bash` file.
- **policy schema** creates a virtual environment, installs `jsonschema` 4.26.0 and `PyYAML` 6.0.3 with all dependencies by pinned SHA-256 (`--require-hashes --only-binary=:all: --no-deps`), and validates `policy/default.yaml` against `spec/policy-v0.schema.json`. The schema's own `$schema` selects the draft. YAML timestamps stay strings, and duplicate keys are an error, because an ambiguous policy file must never pass. The job fails if either file is missing.
- **markdown links** checks relative links offline with lychee. It is non-gating.

### 6.4 scorecard.yml

OpenSSF Scorecard v2.5.0 runs weekly, on pushes to `main` and on branch-protection changes. It uploads SARIF to code scanning and publishes results, which requires a public repository.

### 6.5 Pinned inputs (resolved 2026-10-16)

| Input | Pin | Release date |
|---|---|---|
| actions/checkout | `3d3c42e5aac5ba805825da76410c181273ba90b1` (v7.0.1) | 2026-07-20 |
| actions/upload-artifact | `cf430e030ddbb5b0abf93d22962f4752f3646cd9` (v7.0.2) | 2026-10-07 |
| actions/download-artifact | `9000827ccba6bdab643e8b6fd33ac0654aef8333` (v8.0.2, lightweight tag) | 2026-10-07 |
| sigstore/cosign-installer | `6f9f17788090df1f26f669e9d70d6ae9567deba6` (v4.1.2, installs cosign v3.0.6) | 2026-05-07 |
| ossf/scorecard-action | `fa7117db83c52d787ae5a89f350999701505d78e` (v2.5.0, peeled) | 2026-10-12 |
| github/codeql-action/upload-sarif | `24c54180a607b1449ed407dd24f251e4e9147c8d` (v4.38.3, peeled) | 2026-10-08 |
| lycheeverse/lychee-action | `e7477775783ea5526144ba13e8db5eec57747ce8` (v2.9.0) | 2026-07-09 |
| bootc-image-builder | `quay.io/centos-bootc/bootc-image-builder:latest@sha256:30e7eac70b6e4448da57aaff66d6a89dbd1526572e8b77be31da129e692eb3e9` | 2026-10-15 |
| syft | `ghcr.io/anchore/syft:v1.54.1@sha256:3eb5379ba7b409c3f4069b686110527af0c47df993fa5c10d13e7cf34f49b1aa` | 2026-10-06 |
| trivy | `ghcr.io/aquasecurity/trivy:0.75.0@sha256:af6acf9a6b85dfe389a1941505c0ce9efef52a4719635e1a962f022a3d855daa` | 2026-10-01 |
| actionlint | 1.7.12, SHA-256 `8aca8db9…349a3d8` | 2026-03-30 |

- `anchore/sbom-action` is no longer used. syft runs from its image, which removes one action and its installer script from the trust chain.
- syft v1.54.1 replaces v1.54.0 because it updates the RPM database reader (`go-rpmdb` 0.2.0 to 0.2.2), which matters for a Fedora image.
- `codeql-action` v4.39.0 (released 2026-10-15) and `trivy` 0.76.0 (2026-10-14) were skipped deliberately, because they were one and two days old.
- Renovate applies a 7-day cooldown to GitHub Actions only. Container images get a weekly schedule and human review instead (section 9).

## 7. Verifying the published image

Replace `OWNER` and `REPO` with the GitHub owner and repository, and use the digest from the job summary. The owner in the image name is lowercase; the identity keeps the repository's original case. The same commands apply to the Fedora 44 fallback (tag `main-f44`).

```
IMAGE=ghcr.io/owner/tidemark@sha256:<digest>
IDENTITY=https://github.com/OWNER/REPO/.github/workflows/build-image.yml@refs/heads/main

# Sigstore bundle (cosign v3 default)
cosign verify \
  --certificate-identity "$IDENTITY" \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  "$IMAGE"

# Legacy .sig tag (the format podman and bootc read)
cosign verify --new-bundle-format=false \
  --certificate-identity "$IDENTITY" \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  "$IMAGE"

cosign verify-attestation --type spdxjson \
  --certificate-identity "$IDENTITY" \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  "$IMAGE" > attestation.json
```

- Always verify a digest, not a tag, and always pass the exact identity. Never use `--certificate-identity-regexp '.*'`.
- A successful `verify-attestation` prints the verified in-toto statement, whose predicate is the SBOM. The SBOM file itself is also attached to each build run as a workflow artifact.
- Anonymous verification works only after the GHCR package is made public (section 10).

## 8. Building locally on the M1 Pro (later)

Podman, QEMU and cosign are not installed on the host yet; installing them needs the founder's approval. Lima is installed and is used for the kernel spike, not for image builds. Run one VM at a time: the 4 GiB podman machine and the 6 GiB Lima VM do not fit alongside macOS in 16 GB.

1. Install podman (and later qemu and cosign) with Homebrew.
2. Create a rootful machine on the libkrun provider: `podman machine init --provider libkrun --rootful --cpus 4 --memory 4096 --disk-size 60 tidemark-build`, then `podman machine start tidemark-build`. Check `podman machine init --help` first; older podman releases select the provider through `CONTAINERS_MACHINE_PROVIDER=libkrun` instead of a flag.
3. Build natively for aarch64 from the repository root: `podman build --platform linux/arm64 -f image/Containerfile -t localhost/tidemark:dev image/`. Do not build x86_64 locally: emulated builds run 8 to 10 times slower, and CI builds x86_64 natively for free.
4. Make a qcow2 with the same container as CI, running inside the rootful machine: `podman run --rm --privileged -v ./image/bib/config.toml:/config.toml:ro -v ./output:/output -v /var/lib/containers/storage:/var/lib/containers/storage <BIB_IMAGE from boot-test.yml> --type qcow2 --rootfs xfs localhost/tidemark:dev`. The storage path refers to the machine's storage, and `./output` must be under the home directory, which the machine shares.
5. For an aarch64 boot, copy `config.toml` privately and change the console argument: `console=ttyAMA0` for QEMU's `virt` machine, or `console=hvc0` for Lima on vz. Boot with QEMU and HVF (`qemu-system-aarch64 -machine virt,accel=hvf -cpu host -m 4096 -bios <edk2-aarch64-code.fd> -drive if=virtio,format=qcow2,snapshot=on,file=output/qcow2/disk.qcow2 -nographic -nic none`), or with UTM. Whether Lima's cidata provisioning works on fedora-bootc is checklist item D2.
6. Stop the machine when idle: `podman machine stop tidemark-build`.

## 9. Dependency updates: Renovate

Renovate was chosen over Dependabot for four reasons:

1. It refreshes the digest of an unchanged tag and holds each Containerfile on its own tag.
2. Its regex manager keeps the bootc-image-builder, syft and trivy images pinned inside workflow `env:` blocks, which Dependabot does not parse.
3. It updates SHA-pinned actions together with their version comments.
4. The concept already names it.

Update policy:

- **GitHub Actions** wait 7 days after release (`minimumReleaseAge`). Renovate ages a digest-pinned action against the tag's commit date, so a tag that is force-pushed to new content may pass the cooldown at once. Review every action-update diff.
- **Container images** (the base images and the workflow tool images) get no cooldown, because Quay and GHCR return no release timestamps. Under Renovate's default `minimumReleaseAgeBehaviour` (`timestamp-required`), a cooldown would hold these updates indefinitely. They arrive on a weekly schedule (before 06:00 on Mondays) instead. Each pull request runs both builds and both boot tests, and a human reviews it before merging.
- **Updated by hand:** the actionlint version and hash, the pip hashes, and `COSIGN_VERSION`.

The cost is that the hosted Mend Renovate app needs write access to the repository. The alternative is self-hosting Renovate with a scoped token.

## 10. Founder setup once the repository exists

- Make the GHCR package `tidemark` public after the first push; new packages start private.
- Protect `main` and require the gating checks:
  - `detect image changes (build-image)`;
  - `build (x86_64, fedora-45)` and `build (x86_64, fedora-44)`;
  - `boot test (fedora-45) / boot test (x86_64, QEMU/KVM)` and `boot test (fedora-44) / boot test (x86_64, QEMU/KVM)`;
  - `actionlint`, `shellcheck` and `policy schema`.
- Require the change-detection job too: if it failed, its dependent jobs would be skipped, and a skipped job counts as passing. GitHub names reusable-workflow checks as "caller job / called job"; confirm the exact names after the first run.
- Enable secret scanning and push protection. Keep Actions set to read-only default token permissions.
- Install the Renovate app, or decide to self-host it.
- Keep the repository public; Scorecard publishing and free arm64 runners depend on it.
- Replace the placeholder module path and documentation URLs.

## 11. Known limitations and open risks

- **Not yet run.** Nothing here has executed in CI, and the self-test has not run on a real system. On macOS the scripts have passed only `bash -n` and `shellcheck`, and the workflows only actionlint.
- **aarch64 is not boot-tested.** Phase 0 builds aarch64 natively on `ubuntu-24.04-arm` without gating, and runs no aarch64 boot test, because arm64 runners may lack `/dev/kvm`. Only x86_64 under KVM is a gating boot environment, as the concept states.
- **E3 and T2 are open** until the kernel probe runs in the booted image (section 6.2).
- **bootc-image-builder is deprecated upstream.** The quay.io container is now a compatibility wrapper around `image-builder`. The migration is `ghcr.io/osbuild/image-builder-cli` with `build --bootc-ref <image> --bootc-default-fs xfs --blueprint /blueprint.toml qcow2`. Output paths may differ; the boot test already finds the qcow2 by search.
- **Tag-plus-digest references.** `FROM name:tag@sha256:…` and `podman run name:tag@sha256:…` rely on containers/common dropping the tag, which its current source does. The skopeo docker transport rejects such references, so the workflows never pass one to skopeo. Confirm on the first run with the runner's podman 4.9.
- **No client-side signature enforcement yet; Phase 1 prerequisite.** `bootc upgrade` and `bootc switch` do not yet require a signature, which is why the automatic update timer is masked (section 4). Enforcement needs `/etc/containers/policy.json` with a `sigstoreSigned` requirement bound to the real Fulcio identity, plus a `registries.d` entry with `use-sigstore-attachments: true`. containers/image, which podman and bootc use, finds signatures only through the `sha256-<digest>.sig` tag convention, not through cosign v3's default Sigstore bundles. The publish job therefore also writes a legacy signature and checks that the `.sig` tag exists. A Phase 1 policy must be tested against that tag.
- **Bundle storage on GHCR.** Whether GHCR serves the OCI referrers API for cosign v3 bundles, or cosign falls back to the referrers tag schema, is not verified. The bundle `cosign verify` step will show it on the first run.
- **SBOM size.** The SBOM is package-level (`SYFT_FILE_METADATA_SELECTION=none`) to keep the attestation small. Whether it stays within the transparency-log size limit is not verified.
- **Image identity across jobs.** The design assumes that `podman save --format oci-archive` writes gzip-compressed layers, so that `skopeo copy --preserve-digests` can push the archive without changing its manifest. It also assumes that `podman pull oci-archive:` keeps the config blob, so that the image ID matches. If either assumption fails, the publish or boot-test job fails closed; it never publishes a different image.
- **Artifact size.** Each image archive is expected to be 1 to 2 GB (an estimate, not measured). Archives expire after one day; artifact storage is assumed to be free for a public repository.
- **Clean power-off.** The two-boot test assumes that the guest powers off on an ACPI power-button event through systemd-logind. If it does not, QEMU is stopped with SIGTERM, which QEMU handles cleanly; the first-boot stamp is already flushed by `sync -f`.
- **Weekly runs do not mark the image.** A failed scheduled run means the runner, QEMU or OVMF changed. It produces a failed-workflow notification but does not untag or mark `:main`, which was booted and signed when it was published.
- **Kernel floor.** The self-test now requires Linux 7.2 and `bpf_real_data_inode`, pending the founder's answer to FD-G. Both base images ship 7.2.x.
- **Python wheels.** The pinned wheels cover CPython 3.12 on x86_64 only, which is the `ubuntu-24.04` system Python. A runner image change needs new hashes.
- **Cross-file dependency.** The policy-schema job fails unless `policy/default.yaml` and `spec/policy-v0.schema.json` exist and agree.
- **Not bit-for-bit reproducible.** Inputs are pinned and Phase 0 installs no packages, but layer timestamps differ between builds. The image that is published is the image that was booted, so this no longer affects identity; `--timestamp` or `SOURCE_DATE_EPOCH` can address reproducibility later.
- **Trivy findings.** Base-image CVEs are expected and do not fail the build. Reports are artifacts only and are not uploaded to code scanning. The vulnerability database is fetched unpinned at run time.

## 12. Facts verified on 2026-10-16

| Fact | Source |
|---|---|
| Tag digests, child architectures and kernel labels in section 3 | quay.io tag API and registry v2 manifest and config blobs |
| `landlock_abi_version` is 8 at v7.0, 9 at v7.1 and 10 at v7.2 | `security/landlock/syscalls.c` at those tags, torvalds/linux |
| `LANDLOCK_CREATE_RULESET_VERSION` is `1U << 0`; `__NR_landlock_create_ruleset` is 444 in `syscall_64.tbl` and in asm-generic `unistd.h` | torvalds/linux v7.2 |
| `bpf_get_file_xattr` (KF_SLEEPABLE), `bpf_path_d_path` and `bpf_real_data_inode` (KF_SLEEPABLE, KF_RET_NULL) are registered; `bpf_real_data_inode` is absent at v7.0 and v7.1 | `fs/bpf_fs_kfuncs.c` at v7.0, v7.1, v7.2 and master |
| `bpf_lsm_file_open`, `bpf_lsm_file_receive`, `bpf_lsm_inode_setxattr` and `bpf_lsm_inode_removexattr` are in `sleepable_lsm_hooks` | `kernel/bpf/bpf_lsm.c` at v7.2 |
| `security.bpf.` is `XATTR_NAME_BPF_LSM`; only `user.` and `security.bpf.` names are readable by the xattr kfuncs | `include/uapi/linux/xattr.h` and `fs/bpf_fs_kfuncs.c` at v7.2 |
| Fedora 44 and 45 kernels: `CONFIG_LSM="lockdown,yama,integrity,selinux,bpf,landlock,ipe"`, `CONFIG_BPF_LSM=y`, `CONFIG_SECURITY_LANDLOCK=y`, `CONFIG_DEBUG_INFO_BTF=y`, `CONFIG_OVERLAY_FS=m` | `kernel-x86_64-fedora.config` on the f44 and f45 branches, src.fedoraproject.org/rpms/kernel |
| The fedora-bootc standard image includes `python3-libdnf5`, `python3-rpm` and `attr`; only ELN sets a default `root-fs-type` | gitlab.com/fedora/bootc/base-images manifests |
| bootc-image-builder sources merged into osbuild/image-builder; the quay.io container is a compatibility layer; `/config.toml` and `[customizations.kernel] append` are supported | osbuild/image-builder `bootc-image-builder/README.md` and `doc/20-advanced/20-bootc/` |
| KVM udev rule for GitHub-hosted Linux runners | GitHub changelog, 2 April 2024 |
| Action SHAs, release dates and default tool versions in section 6.5 | `git ls-remote`, GitHub releases API, each action's `action.yml` at the pinned SHA |
| `download-artifact` v8.0.2 fails on a digest mismatch by default (`digest-mismatch: error`); `upload-artifact` v7.0.2 accepts `compression-level` | each action's `action.yml` at the pinned SHA |
| cosign v3.0.6: `sign`, `attest` and `verify` default to `--new-bundle-format=true`; `--use-signing-config` defaults to true and is refused without the new bundle format, so a legacy signature needs both flags set to false | `cmd/cosign/cli/options/sign.go`, `options/verify.go` and `signcommon/common.go` at v3.0.6 |
| syft v1.54.1 (released 2026-10-06) image index digest, `/syft` entrypoint, and its `go-rpmdb` update | ghcr.io registry API; GitHub release notes |
| trivy 0.75.0 `--input` reads a docker-archive tar or an OCI layout directory, not an `oci-archive` tar | `pkg/fanal/image/archive.go` and `oci.go` at v0.75.0 |
| Renovate: `minimumReleaseAgeBehaviour` defaults to `timestamp-required`; GHCR and Quay return no release timestamps, so their digest updates would be held indefinitely; `config:best-practices` adds a release-age rule for npm only | renovatebot/renovate `lib/config/options/index.ts`, `docs/usage/key-concepts/minimum-release-age.md` and `lib/config/presets/internal/` |
| actionlint tarball SHA-256 matches the release checksums file; PyPI wheel hashes | GitHub release assets; PyPI JSON API |
| Renovate's dockerfile manager matches `Containerfile.f44`; `managerFilePatterns` and `custom.regex` are the current option names | renovatebot/renovate source and documentation |
| containers/common strips the tag from tagged and digested references; containers/image's docker transport rejects them | `libimage/normalize.go`, `docker/docker_transport.go` |

## Change log

- v1, 2026-10-16: first Phase 0 draft.
- v1 revised, 2026-10-16, after review:
  - **Pipeline.** One run now builds, boots and publishes the same archive. The Fedora 44 fallback is pushed, signed and verified. The scanners moved out of the job that holds write permissions. A legacy cosign signature was added.
  - **Boot test.** It now boots twice, matches only complete result lines and measures multi-user after the result line.
  - **Self-test and image.** The self-test requires BTF, the kfunc names and a 7.2 floor. The automatic update timer is masked.
  - **Renovate.** The cooldown is limited to GitHub Actions.
  - **Document.** The pinning claims were corrected, and E3 and T2 are recorded as open.
- v1 revised, 2026-10-17, cross-file consistency pass: the Phase 1 RPM placeholder (section 4, item 5, and both Containerfiles) now follows the delivery ADR (same-run RPM checked by SHA-256, no COPR install at build time).

## Tone Check

The "Warm but Authoritative" brand guideline document was not available to the drafter, so the founder should re-score against it before publishing.

| Criterion | Score (1 to 5) | Notes |
|---|---|---|
| Warmth | 3 | Technical reference first; the steps address the reader directly but offer little encouragement. |
| Authority | 5 | Every pin, digest and kernel fact carries a source and a date. |
| Honesty about limits | 5 | States that nothing has run, records the open checklist items and lists every unverified assumption. |
| Clarity | 4 | Tables carry the detail; sections 6 and 11 are long. |
