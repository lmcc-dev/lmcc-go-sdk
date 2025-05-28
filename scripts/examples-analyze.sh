#!/bin/bash
# Copyright 2025 lmcc Authors. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# examples-analyze.sh
# Analyze examples (dependencies, metrics)

set -euo pipefail

EXAMPLES_DIR="${1:-examples}"

echo "===========> Analyzing examples..."
echo "Examples structure:"
if command -v tree >/dev/null 2>&1; then
    tree "$EXAMPLES_DIR"
else
    find "$EXAMPLES_DIR" -type d | sed 's|[^/]*/|  |g'
fi

echo ""
echo "Examples statistics:"

# Count examples by category
basic_count=$(find "$EXAMPLES_DIR" -path "*/basic-usage/*" -name "main.go" | wc -l | tr -d ' ')
config_count=$(find "$EXAMPLES_DIR" -path "*/config-features/*" -name "main.go" | wc -l | tr -d ' ')
error_count=$(find "$EXAMPLES_DIR" -path "*/error-handling/*" -name "main.go" | wc -l | tr -d ' ')
integration_count=$(find "$EXAMPLES_DIR" -path "*/integration/*" -name "main.go" | wc -l | tr -d ' ')
logging_count=$(find "$EXAMPLES_DIR" -path "*/logging-features/*" -name "main.go" | wc -l | tr -d ' ')
total_count=$(find "$EXAMPLES_DIR" -name "main.go" | wc -l | tr -d ' ')

# Handle basic-usage special case (it's directly in the root)
if [ -f "$EXAMPLES_DIR/basic-usage/main.go" ]; then
    basic_count=$((basic_count + 1))
fi

echo "  Total examples: $total_count"
echo "  Categories:"
echo "    - basic-usage: $basic_count examples"
echo "    - config-features: $config_count examples"
echo "    - error-handling: $error_count examples" 
echo "    - integration: $integration_count examples"
echo "    - logging-features: $logging_count examples"

echo ""
echo "Code metrics:"
total_lines=$(find "$EXAMPLES_DIR" -name "*.go" -exec wc -l {} + 2>/dev/null | tail -1 | awk '{print $1}' || echo "0")
total_files=$(find "$EXAMPLES_DIR" -name "*.go" | wc -l | tr -d ' ')

echo "  Total lines of code: $total_lines"
echo "  Total Go files: $total_files"

# Build sizes (if built)
output_dir="_output/examples"
if [ -d "$output_dir" ]; then
    echo ""
    echo "Build artifacts:"
    echo "  Built examples: $(ls "$output_dir" 2>/dev/null | wc -l | tr -d ' ')"
    if [ "$(ls "$output_dir" 2>/dev/null | wc -l)" -gt 0 ]; then
        total_size=$(du -sh "$output_dir" 2>/dev/null | awk '{print $1}' || echo "0")
        echo "  Total build size: $total_size"
    fi
fi