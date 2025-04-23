# Copyright 2025 lmcc Authors. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# ==============================================================================
# Common variables and settings

# Determine the root directory based on the location of this file.
ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))/../..
# Default output directory.
OUTPUT_DIR ?= $(ROOT_DIR)/_output
# Tools directory.
TOOLS_DIR ?= $(OUTPUT_DIR)/tools
TOOLS_BIN_DIR ?= $(TOOLS_DIR)/bin

# Ensure the output and tools bin directories exist.
$(shell mkdir -p $(OUTPUT_DIR))
$(shell mkdir -p $(TOOLS_BIN_DIR))

# Go tool paths.
GO := go

# Determine OS and Architecture.
OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)

# Verbose setting.
V ?= 0
GO_BUILD_FLAGS ?= 
ifeq ($(V), 1)
	GO_BUILD_FLAGS += -v
	Q = # Q is non-empty for non-verbose builds
else
	GO_BUILD_FLAGS += 
	Q = @
endif

# Add the tools bin directory to the PATH for make commands.
export PATH := $(TOOLS_BIN_DIR):$(PATH)

# Print separator line.
define PRINT_SEP
	@echo "=============================================================================="
endef 