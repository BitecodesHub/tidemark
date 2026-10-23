#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
#
# Install the Phase 0 kernel-spike build dependencies on a Fedora 45 VM.
# STATUS: Phase 0 draft. Run as root inside the VM (Lima provisioning and the CI
# cloud-init both run it as root). It installs the clang/LLVM BPF toolchain,
# bpftool, the Go toolchain, libbpf headers and a few utilities the scenario
# scripts use, then prints the kernel facts the spike depends on.

set -euo pipefail

echo ">> installing kernel-spike dependencies"
dnf -y install \
	clang \
	llvm \
	bpftool \
	golang \
	libbpf \
	libbpf-devel \
	elfutils-libelf-devel \
	zlib-devel \
	attr \
	util-linux \
	make

echo ">> toolchain versions"
# sed -n 1p reads all input, so the producer never gets SIGPIPE, which would
# fail the pipeline under pipefail.
uname -r
clang --version | sed -n 1p
bpftool version | sed -n 1p
go version

echo ">> kernel facts the spike depends on (for the serial log)"
echo "lsm: $(cat /sys/kernel/security/lsm 2>/dev/null || echo unreadable)"
if [ -r /sys/kernel/btf/vmlinux ]; then
	echo "btf: /sys/kernel/btf/vmlinux present"
	# Report whether the overlay data-inode kfunc is present in this kernel.
	# Capture the dump first: "grep -q" would exit on the first match and the
	# SIGPIPE to bpftool would make the pipeline fail under pipefail, printing
	# ABSENT for a kfunc that is present.
	btf_raw="$(bpftool btf dump file /sys/kernel/btf/vmlinux format raw 2>/dev/null || true)"
	if grep -q "FUNC 'bpf_real_data_inode'" <<<"${btf_raw}"; then
		echo "kfunc: bpf_real_data_inode present"
	else
		echo "kfunc: bpf_real_data_inode ABSENT (probe will use the fallback)"
	fi
else
	echo "btf: /sys/kernel/btf/vmlinux MISSING (CONFIG_DEBUG_INFO_BTF off?)"
fi

echo ">> dependencies installed"
