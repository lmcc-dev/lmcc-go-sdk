# Copyright 2025 lmcc Authors. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# ==============================================================================
# Makefile helper functions for tools

# Define the list of tools required by the project.
# Add more tools here as needed.
TOOLS_REQUIRED := \
	golangci-lint \
	goimports
#	 golines # Uncomment if needed
#	 mockgen # Example: Uncomment if needed
#	 gotests # Example: Uncomment if needed

# Check if asdf is available
ASDF_AVAILABLE := $(shell command -v asdf >/dev/null 2>&1 && echo "yes" || echo "no")

# Check if Go is managed by asdf
GO_MANAGED_BY_ASDF := $(shell command -v go >/dev/null 2>&1 && (which go | grep -q "asdf" && echo "yes" || echo "no"))

# Get current Go version if managed by asdf
ifeq ($(GO_MANAGED_BY_ASDF),yes)
  GO_VERSION := $(shell asdf current golang | awk 'NR==2 {print $$2}' 2>/dev/null || echo "")
endif

# Function to run asdf reshim if needed
define run_asdf_reshim
	@if [ "$(GO_MANAGED_BY_ASDF)" = "yes" ] && [ -n "$(GO_VERSION)" ]; then \
		echo "Go is managed by asdf, running asdf reshim golang $(GO_VERSION)..."; \
		asdf reshim golang $(GO_VERSION); \
	elif [ "$(ASDF_AVAILABLE)" = "yes" ]; then \
		echo "Asdf is available, checking if $(1) is managed by asdf..."; \
		if TOOL_VERSION=$$(asdf current $(1) 2>/dev/null | awk 'NR==2 {print $$2}'); then \
			if [ -n "$$TOOL_VERSION" ]; then \
				echo "$(1) is managed by asdf (version $$TOOL_VERSION), running asdf reshim $(1) $$TOOL_VERSION..."; \
				asdf reshim $(1) $$TOOL_VERSION; \
			else \
				echo "$(1) is managed by asdf but could not determine version, running asdf reshim $(1)..."; \
				asdf reshim $(1); \
			fi; \
		else \
			echo "$(1) is not managed by asdf, skipping reshim"; \
		fi; \
	fi
endef

# tools.install: Install all required tools.
.PHONY: tools.install
tools.install: $(addprefix tools.install., $(TOOLS_REQUIRED)) ## Install all required Go tools.

# tools.install.%: Generic rule to install a specific tool.
.PHONY: tools.install.%
tools.install.%:
	@echo "===========> Installing $*..."
	@$(MAKE) install.$*
	$(call run_asdf_reshim,$*)

# tools.verify.%: Generic rule to verify if a tool is installed, installing if necessary.
.PHONY: tools.verify.%
tools.verify.%:
	@if ! which $* &>/dev/null; then \
		echo "Tool $* not found, installing..."; \
		$(MAKE) install.$* && echo "Tool $* installed successfully."; \
	else \
		echo "Tool $* found: $(shell which $*)"; \
	fi

# ==============================================================================
# Specific tool installation rules

# Default version for golangci-lint (can be overridden in main Makefile or common.mk)
GOLANGCI_LINT_VERSION ?= v1.58.1

# install.golangci-lint: Install golangci-lint.
.PHONY: install.golangci-lint
install.golangci-lint:
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	$(call run_asdf_reshim,golang)

# install.goimports: Install goimports.
.PHONY: install.goimports
install.goimports:
	@$(GO) install golang.org/x/tools/cmd/goimports@latest
	$(call run_asdf_reshim,golang)

# TODO: Add install rules for other tools as needed.
# Example for golines:
#.PHONY: install.golines
#install.golines:
#	@$(GO) install github.com/segmentio/golines@latest
#	$(call run_asdf_reshim,golang)

# Example for mockgen:
#.PHONY: install.mockgen
#install.mockgen:
#	@$(GO) install github.com/golang/mock/mockgen@latest
#	$(call run_asdf_reshim,golang)

# Example for gotests:
#.PHONY: install.gotests
#install.gotests:
#	@$(GO) install github.com/cweill/gotests/gotests@latest
#	$(call run_asdf_reshim,golang) 