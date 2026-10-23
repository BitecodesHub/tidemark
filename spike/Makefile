# SPDX-License-Identifier: Apache-2.0
#
# AgentOS Tidemark - Phase 0 kernel spike build and test.
# STATUS: Phase 0 draft. These targets run INSIDE a Fedora 45 VM, never on the
# macOS host (the host has no Go or BPF toolchain). Run provision/fedora-deps.sh
# first to install clang, llvm, bpftool, golang and libbpf-devel.
#
# Targets:
#   make vmlinux   generate bpf/vmlinux.h from the running kernel's BTF
#   make build     compile the BPF object and the Go probe binary
#   make test      build, then load the probe and run every scenario (root)
#   make clean     remove build products
#
# `make test` must run as root (it loads BPF LSM programs and writes cgroupfs).
# SUDO defaults to `sudo`; pass SUDO= when already root (the CI and the
# provisioning path do this).

SHELL := /bin/bash

ARCH := $(shell uname -m)
ifeq ($(ARCH),x86_64)
  BPF_ARCH := x86
else ifeq ($(ARCH),aarch64)
  BPF_ARCH := arm64
else
  BPF_ARCH := $(ARCH)
endif

CLANG   ?= clang
BPFTOOL ?= bpftool
GO      ?= go
SUDO    ?= sudo

# Go needs a build cache and a module cache. cloud-init's final stage and some
# sudo configurations run with HOME unset, so pin both caches inside the spike
# tree (ignored by git) and never download a different Go toolchain.
export GOCACHE     ?= $(CURDIR)/.cache/go-build
export GOMODCACHE  ?= $(CURDIR)/.cache/gomod
export GOPATH      ?= $(CURDIR)/.cache/gopath
export GOTOOLCHAIN ?= local

BPF_SRC := bpf/label_probe.bpf.c
BPF_OBJ := label_probe.bpf.o
VMLINUX := bpf/vmlinux.h
PROBE   := tidemark-probe
BPFFS   ?= /sys/fs/bpf/tidemark
# The helper is exec'd inside sessions, which fires file_open and mmap_file on
# the binary itself. On Lima the spike tree is a virtiofs mount whose xattr
# errors could set the fail-closed bit, so tests run a copy on a local
# filesystem instead.
PROBE_INSTALL ?= /usr/local/libexec/tidemark-probe

.PHONY: all vmlinux build test clean help

all: build

help:
	@echo "targets: vmlinux build test clean (see comments in this Makefile)"

vmlinux: $(VMLINUX)

# Regenerated whenever the running kernel's BTF is newer than the header, so a
# kernel update cannot leave a stale vmlinux.h behind.
$(VMLINUX): /sys/kernel/btf/vmlinux
	@echo ">> generating $(VMLINUX) from /sys/kernel/btf/vmlinux ($$(uname -r))"
	$(BPFTOOL) btf dump file /sys/kernel/btf/vmlinux format c > $(VMLINUX).tmp
	mv $(VMLINUX).tmp $(VMLINUX)

build: $(BPF_OBJ) $(PROBE)

$(BPF_OBJ): $(BPF_SRC) $(VMLINUX)
	@echo ">> compiling BPF object for $(BPF_ARCH)"
	# -mcpu=v3: __sync_fetch_and_or with a used result needs BPF ISA v3 atomics.
	# -fms-extensions: kernel 7.x BTF contains tagged anonymous struct members,
	# which bpftool reproduces in vmlinux.h (for example "struct ns_tree;"
	# inside a struct); clang rejects them under -Werror without the extension.
	$(CLANG) -O2 -g -Wall -Werror \
		-fms-extensions -Wno-microsoft-anon-tag \
		-target bpf -mcpu=v3 -D__TARGET_ARCH_$(BPF_ARCH) \
		-Ibpf -c $(BPF_SRC) -o $(BPF_OBJ)

$(PROBE): $(wildcard cmd/tidemark-probe/*.go) go.mod
	@echo ">> building Go probe"
	@mkdir -p $(GOCACHE) $(GOMODCACHE) $(GOPATH)
	GOFLAGS=-mod=mod CGO_ENABLED=0 $(GO) build -o $(PROBE) ./cmd/tidemark-probe
	$(GO) vet ./...

test: build
	@echo ">> running the Phase 0 kernel spike (root)"
	$(SUDO) install -D -m 0755 $(PROBE) $(PROBE_INSTALL)
	$(SUDO) env \
		TM_ATTACH_MMAP_FILE="$(or $(TM_ATTACH_MMAP_FILE),1)" \
		TM_BPF_OBJ="$(CURDIR)/$(BPF_OBJ)" \
		TM_TESTS_DIR="$(CURDIR)/tests" \
		TM_BPFFS="$(BPFFS)" \
		TM_PROBE_BIN="$(PROBE_INSTALL)" \
		"$(PROBE_INSTALL)" run

clean:
	rm -f $(BPF_OBJ) $(PROBE) $(VMLINUX) $(VMLINUX).tmp
	rm -rf .cache
