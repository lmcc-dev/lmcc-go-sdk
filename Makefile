# Copyright 2025 lmcc Authors. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# Default target executed when no arguments are given to make.
.DEFAULT_GOAL := all

# ==============================================================================
# Define project variables

# ROOT_PACKAGE :: The root Go package path of the project.
ROOT_PACKAGE := github.com/lmcc-dev/lmcc-go-sdk
# VERSION_PACKAGE :: The Go package path for version information.
# TODO: Create this package or adjust if versioning is handled differently.
# VERSION_PACKAGE := $(ROOT_PACKAGE)/pkg/version

# ==============================================================================
# Includes - Load common variables and specific rules.
# Order matters.

# Include common variables and functions.
include scripts/make-rules/common.mk
# Include Go specific rules.
include scripts/make-rules/golang.mk
# Include tool installation rules.
include scripts/make-rules/tools.mk
# TODO: Include other rules files as needed (e.g., copyright.mk)

# ==============================================================================
# Define Go Test Variables
# ==============================================================================
# GO_TEST_CMD = go test $(GOFLAGS) $(TESTFLAGS)
# PKG_UNIT_LIST_ALL: List of all packages for unit testing (excludes vendor, examples, integration tests)
PKG_UNIT_LIST_ALL = $(shell go list ./... | grep -v -E \'vendor|examples|test/integration\')
# PKG_COVER_LIST_ALL: List of all packages for coverage calculation (usually same as unit tests)
PKG_COVER_LIST_ALL = $(PKG_UNIT_LIST_ALL)

# PKG: Target package for 'make test-unit' or 'make cover'. Defaults to empty. User must provide for specific package.
PKG ?= 
# PKG_COVER: Target package for 'make cover'. Defaults to PKG or all coverage packages if PKG is empty.
ifeq ($(PKG),"")
	PKG_COVER = $(PKG_COVER_LIST_ALL)
	PKG_UNIT = $(PKG_UNIT_LIST_ALL) # Define PKG_UNIT when PKG is empty
else
	PKG_COVER = $(PKG)
	PKG_UNIT = $(PKG) # Define PKG_UNIT when PKG is provided
endif

# Debug variables to see what's happening
DEBUG_PKGS := $(shell go list ./... | grep -v vendor | grep -v examples | grep -v test/integration)

# ==============================================================================
# Targets

# all: Default target. Runs format, lint, and unit tests.
.PHONY: all
all: format lint test-unit tidy ## Run format, lint, unit tests, and tidy.

# format: Format Go source files.
.PHONY: format
format: tools.verify.goimports ## Format Go source files using gofmt, goimports, etc.
	@echo "===========> Formatting codes..."
	@$(MAKE) go.format

# lint: Run linters.
.PHONY: lint
lint: tools.verify.golangci-lint ## Run linters (golangci-lint).
	@echo "===========> Running linters..."
	@# Check if Go is managed by asdf
	@if command -v go >/dev/null 2>&1 && (which go | grep -q "asdf"); then \
		echo "Go is managed by asdf, checking current version..."; \
		GO_VERSION=$$(asdf current golang | awk 'NR==2 {print $$2}' 2>/dev/null); \
		if [ -n "$$GO_VERSION" ]; then \
			echo "Running asdf reshim golang $$GO_VERSION..."; \
			asdf reshim golang $$GO_VERSION; \
		else \
			echo "Could not determine Go version, running asdf reshim golang..."; \
			asdf reshim golang; \
		fi; \
	elif command -v asdf >/dev/null 2>&1; then \
		echo "Checking if golang is managed by asdf..."; \
		if GO_VERSION=$$(asdf current golang 2>/dev/null | awk 'NR==2 {print $$2}'); then \
			if [ -n "$$GO_VERSION" ]; then \
				echo "Golang is managed by asdf (version $$GO_VERSION), running reshim..."; \
				asdf reshim golang $$GO_VERSION; \
			else \
				echo "Golang is managed by asdf but could not determine version, running reshim..."; \
				asdf reshim golang; \
			fi; \
		fi; \
	fi
	@# Run linters
	@$(MAKE) go.lint || { \
		echo "ERROR: Linting failed."; \
		exit 1; \
	}

# test-unit: Run unit tests. Can specify package with PKG=./path/to/pkg
.PHONY: test-unit
test-unit: ## Run unit tests. Usage: make test-unit [PKG=./pkg/log]
	@echo "DEBUG: PKG is [$(PKG)]"
	@echo "DEBUG: DEBUG_PKGS is [$(DEBUG_PKGS)]"
ifeq ($(PKG),)
	@echo "===========> Running unit tests for all packages (excluding integration, examples, vendor)..."
	@echo "Command: go test $(GOFLAGS) $(TESTFLAGS) $(DEBUG_PKGS)"
	@go test $(GOFLAGS) $(TESTFLAGS) $(DEBUG_PKGS)
else
	@echo "===========> Running unit tests for $(PKG)..."
	@echo "Command: go test $(GOFLAGS) $(TESTFLAGS) $(PKG)"
	@go test $(GOFLAGS) $(TESTFLAGS) $(PKG)
endif

# test-integration: Run integration tests.
.PHONY: test-integration
test-integration: ## Run integration tests (found in ./test/integration/...).
	@echo "===========> Running integration tests..."
	@go test $(GOFLAGS) $(TESTFLAGS) ./test/integration/...

# test: Run all unit tests (compatibility). Now defaults to test-unit.
.PHONY: test
test: test-unit ## Run all unit tests (excludes integration tests).

# cover: Run unit tests and generate coverage report. Can specify package with PKG=./path/to/pkg
.PHONY: cover
cover: ## Run unit tests and generate coverage report. Usage: make cover [PKG=./pkg/log]
	@echo "DEBUG: PKG is [$(PKG)]"
	@echo "DEBUG: DEBUG_PKGS is [$(DEBUG_PKGS)]"
	@echo "DEBUG: COVER_DIR is [$(COVER_DIR)]"
ifeq ($(PKG),)
	@echo "===========> Running unit tests with coverage for all packages (excluding integration, examples, vendor)..."
	@echo "Command: go test $(GOFLAGS) $(TESTFLAGS) -cover -coverprofile=$(COVER_DIR)/coverage.out $(DEBUG_PKGS)"
	@go test $(GOFLAGS) $(TESTFLAGS) -cover -coverprofile=$(COVER_DIR)/coverage.out $(DEBUG_PKGS)
else
	@echo "===========> Running unit tests with coverage for $(PKG)..."
	@echo "Command: go test $(GOFLAGS) $(TESTFLAGS) -cover -coverprofile=$(COVER_DIR)/coverage.out $(PKG)"
	@go test $(GOFLAGS) $(TESTFLAGS) -cover -coverprofile=$(COVER_DIR)/coverage.out $(PKG)
endif
	@echo "===========> go tool cover -html=$(COVER_DIR)/coverage.out -o $(COVER_DIR)/coverage.html"
	@go tool cover -html=$(COVER_DIR)/coverage.out -o $(COVER_DIR)/coverage.html
	@echo "===========> Coverage report generated at $(COVER_DIR)/coverage.html"

# cover-func: Show function coverage summary. Requires PKG.
.PHONY: cover-func
cover-func: ## Show function coverage summary for a specific package. Usage: make cover-func PKG=./pkg/log
	@[ -n "$(PKG)" ] || ( echo "ERROR: 'make cover-func' requires the PKG variable to be set to a specific package, e.g., PKG=./pkg/log"; exit 1 )
	@echo "===========> Ensuring coverage report exists for $(PKG)..."
	@$(MAKE) cover PKG=$(PKG) > /dev/null # Ensure coverage report exists for the package
	@echo "===========> Displaying function coverage summary for $(PKG)..."
	@scripts/format-func-coverage.sh "$(PKG)" "$(COVER_DIR)/coverage.out" "$(ROOT_PACKAGE)"

# tidy: Tidy go module files.
.PHONY: tidy
tidy: ## Tidy go module files (go mod tidy).
	@echo "===========> Tidying go module files..."
	@$(GO) mod tidy

# ==============================================================================
# Documentation Targets
# ==============================================================================

# doc-view: View package documentation in the terminal. Requires PKG.
.PHONY: doc-view
doc-view: ## View package documentation in terminal. Usage: make doc-view PKG=./pkg/log
	@[ -n "$(PKG)" ] || ( echo "错误：'make doc-view' 需要设置 PKG 变量来指定一个包，例如：PKG=./pkg/log"; exit 1 )
	@echo "===========> 正在查看 $(ROOT_PACKAGE)/$(PKG) 的文档..."
	@# 如果 PKG 以 ./ 开头，则直接使用，否则拼接 ROOT_PACKAGE
	@if echo "$(PKG)" | grep -q '^\./'; then \
		$(GO) doc $(PKG); \
	else \
		$(GO) doc $(ROOT_PACKAGE)/$(PKG); \
	fi

# doc-serve: Serve HTML documentation locally using godoc.
.PHONY: doc-serve
doc-serve: tools.verify.godoc ## Serve HTML documentation locally (godoc -http=:6060).
	@echo "===========> 正在启动 godoc 服务于 http://localhost:6060"
	@echo "===========> 提供 $(ROOT_PACKAGE) 的文档"
	@echo "===========> 按 Ctrl+C 停止服务。"
	@# 确保 GOBIN 在 PATH 中，或者 godoc (如果通过 go install 安装) 可直接调用
	@# 如果 godoc 安装在 $(TOOLS_BIN_DIR)，你可能需要像 $(GODOC_TOOL) 那样调用它
	@godoc -http=:6060 -play=true

# ==============================================================================
# End Documentation Targets
# ==============================================================================

# ==============================================================================
# Examples Management Targets
# ==============================================================================

# EXAMPLES_DIR :: Examples directory path
EXAMPLES_DIR := examples
# EXAMPLES_OUTPUT_DIR :: Build output directory for examples
EXAMPLES_OUTPUT_DIR := $(OUTPUT_DIR)/examples

# Get all example directories with main.go
EXAMPLE_DIRS := $(shell find $(EXAMPLES_DIR) -name "main.go" -type f | xargs dirname | sort)
# Get all example names for convenience
EXAMPLE_NAMES := $(shell find $(EXAMPLES_DIR) -name "main.go" -type f | xargs dirname | sed 's|$(EXAMPLES_DIR)/||g' | sort)

# examples-list: List all available examples
.PHONY: examples-list
examples-list: ## List all available examples
	@echo "===========> Available Examples:"
	@echo "$(EXAMPLE_NAMES)" | tr ' ' '\n' | nl -v 1 -s '. '
	@echo ""
	@echo "Total: $(words $(EXAMPLE_NAMES)) examples"

# examples-build: Build all examples
.PHONY: examples-build
examples-build: ## Build all examples into $(EXAMPLES_OUTPUT_DIR)
	@./scripts/examples-build.sh $(EXAMPLES_DIR) $(EXAMPLES_OUTPUT_DIR)

# examples-clean: Clean examples build output
.PHONY: examples-clean
examples-clean: ## Clean examples build output
	@echo "===========> Cleaning examples build output..."
	@rm -rf $(EXAMPLES_OUTPUT_DIR)
	@echo "===========> Examples build output cleaned!"

# examples-run: Run a specific example
.PHONY: examples-run
examples-run: ## Run a specific example. Usage: make examples-run EXAMPLE=basic-usage
	@[ -n "$(EXAMPLE)" ] || ( echo "ERROR: 'make examples-run' requires EXAMPLE variable. Use: make examples-run EXAMPLE=basic-usage"; echo "Available examples:"; echo "$(EXAMPLE_NAMES)" | tr ' ' '\n' | nl -v 1 -s '. '; exit 1 )
	@example_dir=$(EXAMPLES_DIR)/$(EXAMPLE); \
	if [ ! -d "$$example_dir" ]; then \
		echo "ERROR: Example '$(EXAMPLE)' not found in $$example_dir"; \
		echo "Available examples:"; \
		echo "$(EXAMPLE_NAMES)" | tr ' ' '\n' | nl -v 1 -s '. '; \
		exit 1; \
	fi; \
	if [ ! -f "$$example_dir/main.go" ]; then \
		echo "ERROR: main.go not found in $$example_dir"; \
		exit 1; \
	fi
	@echo "===========> Running example: $(EXAMPLE)"
	@echo "Working directory: $(EXAMPLES_DIR)/$(EXAMPLE)"
	@cd $(EXAMPLES_DIR)/$(EXAMPLE) && $(GO) run main.go

# examples-run-all: Run all examples sequentially
.PHONY: examples-run-all
examples-run-all: ## Run all examples sequentially
	@./scripts/examples-run-all.sh $(EXAMPLES_DIR)

# examples-test: Test all examples (lint + build)
.PHONY: examples-test
examples-test: ## Test all examples (lint + build verification)
	@echo "===========> Testing all examples..."
	@echo "Step 1: Linting examples..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run $(EXAMPLES_DIR)/.../ --timeout=5m; \
	else \
		echo "WARNING: golangci-lint not found, skipping lint check"; \
	fi
	@echo "Step 2: Building examples..."
	@$(MAKE) examples-build
	@echo "===========> All examples tested successfully!"

# examples-debug: Debug a specific example with delve
.PHONY: examples-debug
examples-debug: tools.verify.dlv ## Debug a specific example with delve. Usage: make examples-debug EXAMPLE=basic-usage
	@[ -n "$(EXAMPLE)" ] || ( echo "ERROR: 'make examples-debug' requires EXAMPLE variable. Use: make examples-debug EXAMPLE=basic-usage"; echo "Available examples:"; echo "$(EXAMPLE_NAMES)" | tr ' ' '\n' | nl -v 1 -s '. '; exit 1 )
	@example_dir=$(EXAMPLES_DIR)/$(EXAMPLE); \
	if [ ! -d "$$example_dir" ]; then \
		echo "ERROR: Example '$(EXAMPLE)' not found in $$example_dir"; \
		echo "Available examples:"; \
		echo "$(EXAMPLE_NAMES)" | tr ' ' '\n' | nl -v 1 -s '. '; \
		exit 1; \
	fi; \
	if [ ! -f "$$example_dir/main.go" ]; then \
		echo "ERROR: main.go not found in $$example_dir"; \
		exit 1; \
	fi
	@echo "===========> Starting debugger for example: $(EXAMPLE)"
	@echo "Working directory: $(EXAMPLES_DIR)/$(EXAMPLE)"
	@echo "Debugger commands:"
	@echo "  b main.main    - Set breakpoint at main function"
	@echo "  c              - Continue execution"
	@echo "  n              - Next line"
	@echo "  s              - Step into"
	@echo "  l              - List source code"
	@echo "  p <var>        - Print variable"
	@echo "  q              - Quit debugger"
	@echo "---"
	@cd $(EXAMPLES_DIR)/$(EXAMPLE) && dlv debug main.go

# examples-analyze: Analyze examples (dependencies, metrics)
.PHONY: examples-analyze
examples-analyze: ## Analyze examples (dependencies, build size, etc.)
	@./scripts/examples-analyze.sh $(EXAMPLES_DIR)

# examples-category: Run all examples in a specific category
.PHONY: examples-category
examples-category: ## Run all examples in a category. Usage: make examples-category CATEGORY=config-features
	@./scripts/examples-category.sh $(EXAMPLES_DIR) $(CATEGORY)

# ==============================================================================
# End Examples Management Targets
# ==============================================================================

# clean: Remove build artifacts.
.PHONY: clean
clean: ## Remove all files generated by the build.
	@echo "===========> Cleaning all build output..."
	@-rm -vrf $(OUTPUT_DIR)
	@echo "===========> Cleaning tools..."
	@-rm -vrf $(TOOLS_DIR)

# help: Show this help message.
.PHONY: help
help: Makefile ## Show this help message.
	@printf "\nUsage: make <TARGETS> [PKG=./path/to/package]\n\nTargets:\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $< | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-25s %s\n", $$1, $$2}'
	@printf "\nPKG Variable:\n"
	@printf "  Used with 'test-unit', 'cover', and 'cover-func' to specify a target package.\n"
	@printf "  If not provided for 'test-unit' or 'cover', it runs on all unit test packages.\n"
	@printf "  Example: make test-unit PKG=./pkg/log\n\n"
	@printf "\nEXAMPLE Variable:\n"
	@printf "  Used with 'examples-run' and 'examples-debug' to specify a target example.\n"
	@printf "  Example: make examples-run EXAMPLE=basic-usage\n"
	@printf "  Example: make examples-debug EXAMPLE=config-features/01-simple-config\n\n"
	@printf "\nCATEGORY Variable:\n"
	@printf "  Used with 'examples-category' to run all examples in a category.\n"
	@printf "  Available categories: basic-usage, config-features, error-handling, integration, logging-features\n"
	@printf "  Example: make examples-category CATEGORY=config-features\n\n"

# tools: Install all required Go tools.
.PHONY: tools
tools: ## Install all required Go tools listed in tools.mk.
	@$(MAKE) tools.install

# examples: Build and test all examples (default examples target)
.PHONY: examples
examples: examples-test ## Build and test all examples (default examples target)

# ==============================================================================
# Removed: Tools section moved to scripts/make-rules/tools.mk
# ==============================================================================

# TOOLS_DIR ?= $(OUTPUT_DIR)/tools
# TOOLS_BIN_DIR ?= $(TOOLS_DIR)/bin
# GOLANGCI_LINT_VERSION := v1.58.1 # Choose a specific version
# GOLANGCI_LINT := $(TOOLS_BIN_DIR)/golangci-lint

# # Ensure golangci-lint is installed
# tools.verify.golangci-lint: $(GOLANGCI_LINT) ## Verify golangci-lint is installed.
# $(GOLANGCI_LINT):
# 	@echo "===========> Installing golangci-lint $(GOLANGCI_LINT_VERSION)..."
# 	@GOBIN=$(TOOLS_BIN_DIR) $(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)