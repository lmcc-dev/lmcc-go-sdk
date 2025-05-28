#!/bin/bash
# Copyright 2025 lmcc Authors. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# examples-run-all.sh
# Run all examples sequentially

set -euo pipefail

EXAMPLES_DIR="${1:-examples}"

echo "===========> Running all examples..."

# Find all example directories with main.go
example_dirs=($(find "$EXAMPLES_DIR" -name "main.go" -type f | xargs dirname | sort))
total=${#example_dirs[@]}
current=1

for dir in "${example_dirs[@]}"; do
    example_name=$(echo "$dir" | sed "s|$EXAMPLES_DIR/||g")
    echo ""
    echo "===========> [$current/$total] Running: $example_name"
    echo "Directory: $dir"
    echo "---"
    
    if cd "$dir" && go run main.go; then
        echo "✓ Completed: $example_name"
    else
        echo "✗ Failed or timed out: $example_name"
    fi
    
    cd - > /dev/null
    current=$((current + 1))
done

echo ""
echo "===========> All examples execution completed!"