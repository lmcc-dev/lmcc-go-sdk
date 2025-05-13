# Copyright 2025 lmcc Authors. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# ==============================================================================
# Go related rules

# Default packages to run commands against.
PKGS ?= $(shell $(GO) list ./... | grep -v /vendor/)
# Packages for tests.
TEST_PKGS ?= $(shell $(GO) list ./... | grep -v /vendor/ | grep -v /examples/)

# Golangci-lint configuration.
# TODO: Create a .golangci.yaml configuration file.
GOLANGCI_LINT_CONFIG ?= $(ROOT_DIR)/.golangci.yaml
# 确保 golangci-lint 命令能正确获取
GOLANGCI_LINT ?= golangci-lint
# 检查配置文件是否存在，如果不存在则不使用 --config 参数
GOLANGCI_CONFIG_PARAM := $(shell test -f $(GOLANGCI_LINT_CONFIG) && echo "--config=$(GOLANGCI_LINT_CONFIG)" || echo "")
GOLANGCI_LINT_CMD = $(GOLANGCI_LINT) run $(GOLANGCI_CONFIG_PARAM) --verbose

# go.build: Build packages individually.
.PHONY: go.build
go.build:
	@echo "Building packages..."
	@$(GO) build $(GO_BUILD_FLAGS) $(PKGS)

# go.test: Run unit tests. Accepts optional PKG and RUN variables.
# Example: make test PKG=./pkg/log RUN=TestMyFunction
# Example: make test PKG=./pkg/config
# Example: make test RUN=TestSpecificFeature
.PHONY: go.test
go.test:
	@echo "Running unit tests..."
	@# Using a single shell invocation for variable assignments and execution
	_TEST_TARGETS='$(TEST_PKGS)'; \
	if [ -n "$(PKG)" ]; then \
	    _TEST_TARGETS='$(PKG)'; \
	    echo "--> Testing specific package(s): $(PKG)"; \
	fi; \
	_RUN_FILTER=''; \
	if [ -n "$(RUN)" ]; then \
	    _RUN_FILTER="-run=$(RUN)"; \
	    echo "--> Filtering tests with pattern: $(RUN)"; \
	fi; \
	echo "--> Executing: $(GO) test -v -race -count=1 $(GO_BUILD_FLAGS) $$_RUN_FILTER $$_TEST_TARGETS"; \
	$(GO) test -v -race -count=1 $(GO_BUILD_FLAGS) $$_RUN_FILTER $$_TEST_TARGETS

# go.test.cover: Run unit tests and generate coverage report.
COVER_DIR := $(OUTPUT_DIR)/coverage
COVER_PROFILE := $(COVER_DIR)/coverage.out
COVER_HTML_PROFILE := $(COVER_DIR)/coverage.html
.PHONY: go.test.cover
go.test.cover:
	@echo "Running unit tests with coverage..."
	@mkdir -p $(COVER_DIR)
	@$(GO) test -race -count=1 -coverprofile=$(COVER_PROFILE) -covermode=atomic $(GO_BUILD_FLAGS) $(TEST_PKGS)
	@echo "Coverage profile saved to $(COVER_PROFILE)"
	@$(GO) tool cover -html=$(COVER_PROFILE) -o $(COVER_HTML_PROFILE)
	@echo "HTML coverage report saved to $(COVER_HTML_PROFILE)"

# go.lint: Run linters.
.PHONY: go.lint
go.lint:
	@echo "Running golangci-lint..."
	@if ! which $(GOLANGCI_LINT) >/dev/null 2>&1; then \
		echo "golangci-lint not found, please install it first"; \
		exit 1; \
	fi
	@cd $(ROOT_DIR) && $(GOLANGCI_LINT_CMD) --path-prefix="" --modules-download-mode=readonly ./... || { \
		echo "Linting failed. Fixing automatically with 'golangci-lint run --fix' when possible..."; \
		cd $(ROOT_DIR) && $(GOLANGCI_LINT_CMD) --path-prefix="" --modules-download-mode=readonly --fix ./... 2>/dev/null || true; \
		exit 1; \
	}
	@echo "Running go vet..."
	@$(GO) vet $(PKGS)

# go.format: Format Go source code.
.PHONY: go.format
go.format:
	@echo "Formatting with gofmt..."
	@gofmt -s -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")
	@echo "Formatting with goimports..."
	@goimports -w -local $(ROOT_PACKAGE) $(shell find . -type f -name '*.go' -not -path "./vendor/*")
	# TODO: Consider adding golines if needed and install it via tools.mk
	# @echo "Formatting with golines..."
	# @golines -w --max-len=120 --reformat-tags --shorten-comments --ignore-generated $(shell find . -type f -name '*.go' -not -path "./vendor/*")
	@echo "Formatting go.mod..."
	@$(GO) mod edit -fmt 