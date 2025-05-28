#!/bin/bash
# Copyright 2025 lmcc Authors. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# examples-build.sh
# Build all examples script

set -euo pipefail

EXAMPLES_DIR="${1:-examples}"
OUTPUT_DIR="${2:-_output/examples}"

# Create output directory
mkdir -p "$OUTPUT_DIR"

echo "===========> Building all examples..."

# Find all example directories with main.go
find "$EXAMPLES_DIR" -name "main.go" -type f | while read -r main_file; do
    example_dir=$(dirname "$main_file")
    echo "Building $example_dir..."
    
    # Generate output filename
    example_name=$(echo "$example_dir" | sed "s|$EXAMPLES_DIR/||g" | tr '/' '-')
    output_file="$OUTPUT_DIR/$example_name"
    
    # Build the example
    if go build -o "$output_file" "./$example_dir"; then
        echo "✓ Built: $output_file"
    else
        echo "✗ Failed to build: $example_dir"
        exit 1
    fi
done

echo "===========> All examples built successfully!"