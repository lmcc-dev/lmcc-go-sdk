# Copyright 2025 lmcc Authors. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# ==============================================================================
# Makefile helper functions for tools

# Define the list of tools required by the project.
# Add more tools here as needed.
TOOLS_REQUIRED := \
	golangci-lint \
	goimports \
	godoc
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

# tools: Alias for tools.install for convenience.
# tools: tools.install 的便捷别名。
.PHONY: tools
tools: tools.install ## Install all required Go tools. Use GOLANGCI_LINT_STRATEGY=latest|auto|stable to control golangci-lint version.

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

# tools.version: Show versions of all installed tools.
# tools.version: 显示所有已安装工具的版本。
.PHONY: tools.version
tools.version: ## Show versions of all installed tools.
	@echo "===========> Tool Versions"
	@echo "Go version: $(GO_VERSION_FULL) ($(GO_VERSION_MAJOR_MINOR))"
	@echo "golangci-lint strategy: $(GOLANGCI_LINT_STRATEGY)"
	@echo "Recommended golangci-lint version: $(GOLANGCI_LINT_VERSION_DEFAULT)"
	@echo "Current golangci-lint version: $$(golangci-lint version 2>/dev/null | head -1 || echo 'Not installed')"
	@echo "Current goimports version: $$(goimports -h 2>&1 | head -1 || echo 'Not installed')"
	@echo "Current godoc version: $$(godoc -h 2>&1 | head -1 || echo 'Not installed')"
	@echo ""
	@echo "To use latest golangci-lint: make tools GOLANGCI_LINT_STRATEGY=latest"
	@echo "To use auto-detection: make tools GOLANGCI_LINT_STRATEGY=auto"

# tools.help: Show detailed help for tools management.
# tools.help: 显示工具管理的详细帮助。
.PHONY: tools.help
tools.help: ## Show detailed help for tools management strategies.
	@echo "===========> Tools Management Help"
	@echo ""
	@echo "Available strategies for golangci-lint version selection:"
	@echo "可用的 golangci-lint 版本选择策略："
	@echo ""
	@echo "1. stable (default): Use tested stable version v1.64.8"
	@echo "   稳定版（默认）：使用经过测试的稳定版本 v1.64.8"
	@echo "   - Ensures reproducible builds across team"
	@echo "   - 确保团队间构建的可重现性"
	@echo "   Usage: make tools"
	@echo ""
	@echo "2. latest: Always use the latest available version"
	@echo "   最新版：始终使用最新可用版本"
	@echo "   - Gets cutting-edge features and bug fixes"
	@echo "   - 获得前沿特性和错误修复"
	@echo "   - May introduce breaking changes"
	@echo "   - 可能引入破坏性变更"
	@echo "   Usage: make tools GOLANGCI_LINT_STRATEGY=latest"
	@echo ""
	@echo "3. auto: Auto-select based on Go version"
	@echo "   自动选择：根据 Go 版本自动选择"
	@echo "   - Go 1.24+: latest"
	@echo "   - Go 1.23 and below: v1.64.8"
	@echo "   Usage: make tools GOLANGCI_LINT_STRATEGY=auto"
	@echo ""
	@echo "You can also override the specific version:"
	@echo "您也可以覆盖特定版本："
	@echo "   make tools GOLANGCI_LINT_VERSION=v1.65.0"

# ==============================================================================
# Specific tool installation rules

# Determine golangci-lint version based on Go version
# 根据 Go 版本确定 golangci-lint 版本
GO_VERSION_FULL := $(shell go version | awk '{print $$3}' | sed 's/go//')
GO_VERSION_MAJOR_MINOR := $(shell go version | awk '{print $$3}' | sed 's/go//' | cut -d. -f1,2)

# Determine the best golangci-lint version for current Go version
# 确定当前 Go 版本的最佳 golangci-lint 版本
# 
# Strategy options (can be overridden by setting GOLANGCI_LINT_STRATEGY):
# 策略选项（可以通过设置 GOLANGCI_LINT_STRATEGY 来覆盖）：
# - "stable": Use tested stable version (default for reproducible builds)
# - "latest": Always use latest version (for cutting-edge features)
# - "auto": Auto-select based on Go version
#
# - "stable": 使用经过测试的稳定版本（默认，用于可重现构建）
# - "latest": 始终使用最新版本（用于前沿特性）
# - "auto": 根据 Go 版本自动选择
GOLANGCI_LINT_STRATEGY ?= stable

# Determine golangci-lint version based on strategy
# 根据策略确定 golangci-lint 版本
ifeq ($(GOLANGCI_LINT_STRATEGY),latest)
	GOLANGCI_LINT_VERSION_DEFAULT := latest
else ifeq ($(GOLANGCI_LINT_STRATEGY),auto)
	ifeq ($(GO_VERSION_MAJOR_MINOR),1.24)
		GOLANGCI_LINT_VERSION_DEFAULT := latest
	else ifeq ($(GO_VERSION_MAJOR_MINOR),1.25)
		GOLANGCI_LINT_VERSION_DEFAULT := latest
	else
		GOLANGCI_LINT_VERSION_DEFAULT := v1.64.8
	endif
else
	# Default to stable strategy
	# 默认使用稳定策略
	GOLANGCI_LINT_VERSION_DEFAULT := v1.64.8
endif

# Allow override of golangci-lint version (can be set in main Makefile or common.mk)
# 允许覆盖 golangci-lint 版本（可以在主 Makefile 或 common.mk 中设置）
GOLANGCI_LINT_VERSION ?= $(GOLANGCI_LINT_VERSION_DEFAULT)

# install.golangci-lint: Install golangci-lint with version auto-detection.
# install.golangci-lint: 安装 golangci-lint 并自动检测版本。
.PHONY: install.golangci-lint
install.golangci-lint:
	@echo "===========> Detected Go version: $(GO_VERSION_MAJOR_MINOR)"
	@echo "===========> Installing golangci-lint version: $(GOLANGCI_LINT_VERSION)"
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	$(call run_asdf_reshim,golang)

# install.goimports: Install goimports.
.PHONY: install.goimports
install.goimports:
	@$(GO) install golang.org/x/tools/cmd/goimports@latest
	$(call run_asdf_reshim,golang)

# install.godoc: Install godoc.
.PHONY: install.godoc
install.godoc:
	@echo "Installing godoc (golang.org/x/tools/cmd/godoc)..."
	@$(GO) install golang.org/x/tools/cmd/godoc@latest
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