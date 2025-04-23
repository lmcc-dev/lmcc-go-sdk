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
GOLANGCI_LINT_CMD ?= $(GOLANGCI_LINT) run --config=$(GOLANGCI_LINT_CONFIG) --verbose

# go.build: Build packages individually.
.PHONY: go.build
go.build:
	@echo "Building packages..."
	@$(GO) build $(GO_BUILD_FLAGS) $(PKGS)

# go.test: Run unit tests.
.PHONY: go.test
go.test:
	@echo "Running unit tests..."
	@$(GO) test -race -count=1 $(GO_BUILD_FLAGS) $(TEST_PKGS)

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
	@$(GOLANGCI_LINT_CMD) $(PKGS)
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