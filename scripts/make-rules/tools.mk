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

# tools.install: Install all required tools.
.PHONY: tools.install
tools.install: $(addprefix tools.install., $(TOOLS_REQUIRED)) ## Install all required Go tools.

# tools.install.%: Generic rule to install a specific tool.
.PHONY: tools.install.%
tools.install.%:
	@echo "===========> Installing $*..."
	@$(MAKE) install.$*

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

# install.goimports: Install goimports.
.PHONY: install.goimports
install.goimports:
	@$(GO) install golang.org/x/tools/cmd/goimports@latest

# TODO: Add install rules for other tools as needed.
# Example for golines:
#.PHONY: install.golines
#install.golines:
#	@$(GO) install github.com/segmentio/golines@latest

# Example for mockgen:
#.PHONY: install.mockgen
#install.mockgen:
#	@$(GO) install github.com/golang/mock/mockgen@latest

# Example for gotests:
#.PHONY: install.gotests
#install.gotests:
#	@$(GO) install github.com/cweill/gotests/gotests@latest 