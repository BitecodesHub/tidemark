#!/usr/bin/bash
# SPDX-License-Identifier: Apache-2.0
#
# AgentOS Tidemark first-boot script.
# Status: Phase 0 draft (2026-10-16, revised after review). Not a product
# component.
#
# Runs once, from tidemark-firstboot.service. It prints the Phase 0 banner,
# runs the boot self-test, and writes a stamp so that it does not run again.
# On later boots, tidemark-selftest.service runs the self-test instead, so
# every boot prints exactly one TIDEMARK_SELFTEST result line. The CI boot test
# (.github/workflows/boot-test.yml) is designed to check this on two
# consecutive boots of the same disk; it has not run yet.
#
# It creates no accounts, applies no labels, provisions no keys and opens no
# network connections. Those first-boot duties arrive with the Tidemark RPM
# in Phase 1 and are not implemented here.
#
# Exit status: the exit status of the self-test (0 PASS, 1 FAIL).

set -u
export LC_ALL=C
PATH=/usr/sbin:/usr/bin:/sbin:/bin
export PATH

readonly STAMP=/var/lib/tidemark/firstboot.done
readonly SELFTEST=/usr/libexec/tidemark/tidemark-selftest.sh

os_name="unknown"
if [[ -r /usr/lib/os-release ]]; then
  # shellcheck disable=SC1091
  os_name=$(. /usr/lib/os-release && printf '%s' "${PRETTY_NAME:-unknown}")
fi

cat <<EOF
======================================================================
 AgentOS Tidemark: Phase 0 draft image
======================================================================
 Status: build and boot pipeline test only.
 This image contains no Tidemark enforcement yet: no read-time labels,
 no clearance-gated exits, no model broker, no staged writes and no
 receipts. Do not rely on it to protect private data.
 Base: ${os_name}
 Kernel: $(uname -r)
 The boot self-test runs now and prints one machine-readable result
 line. Documentation: docs/2026-10-16_Image-Pipeline_v1.md in the
 project repository.
======================================================================
EOF

rc=0
if [[ -x $SELFTEST ]]; then
  "$SELFTEST" || rc=$?
else
  printf 'tidemark-firstboot: self-test not found at %s\n' "$SELFTEST"
  rc=1
fi

if mkdir -p -- "${STAMP%/*}"; then
  printf 'firstboot_completed=%s\nselftest_exit=%s\n' \
    "$(date -u +%Y-%m-%dT%H:%M:%SZ)" "$rc" >"$STAMP"
  # Flush the file system that holds the stamp, so that the stamp survives
  # an unclean power-off (for example, the boot test stopping QEMU) and the
  # next boot does not repeat the first-boot banner.
  sync -f -- "$STAMP" || printf 'tidemark-firstboot: could not sync %s\n' "$STAMP"
else
  printf 'tidemark-firstboot: could not create %s\n' "${STAMP%/*}"
fi

exit "$rc"
