#!/bin/bash
# Copyright 2025 lmcc Authors. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# examples-category.sh
# Run all examples in a specific category

set -euo pipefail

EXAMPLES_DIR="${1:-examples}"
CATEGORY="${2:-}"

if [ -z "$CATEGORY" ]; then
    echo "ERROR: 'examples-category' requires CATEGORY variable."
    echo "Available categories: basic-usage, config-features, error-handling, integration, logging-features"
    exit 1
fi

echo "===========> Running examples in category: $CATEGORY"

# Find examples in the category
if [ "$CATEGORY" = "basic-usage" ]; then
    # Special case for basic-usage which is at the root level
    if [ -f "$EXAMPLES_DIR/basic-usage/main.go" ]; then
        category_examples=("basic-usage")
    else
        category_examples=()
    fi
else
    # Find examples in subdirectory
    category_examples=($(find "$EXAMPLES_DIR/$CATEGORY" -name "main.go" -type f 2>/dev/null | xargs dirname | sed "s|$EXAMPLES_DIR/||g" | sort || true))
fi

if [ ${#category_examples[@]} -eq 0 ]; then
    echo "ERROR: No examples found for category '$CATEGORY'"
    echo "Available categories: basic-usage, config-features, error-handling, integration, logging-features"
    exit 1
fi

echo "Examples: ${category_examples[*]}"
echo ""

for example in "${category_examples[@]}"; do
    echo "===========> Running: $example"
    
    example_dir="$EXAMPLES_DIR/$example"
    if [ ! -d "$example_dir" ]; then
        echo "✗ Failed: Example directory '$example_dir' not found"
        continue
    fi
    
    if [ ! -f "$example_dir/main.go" ]; then
        echo "✗ Failed: main.go not found in '$example_dir'"
        continue
    fi
    
    echo "Working directory: $example_dir"
    
    if cd "$example_dir" && go run main.go; then
        echo "✓ Completed: $example"
    else
        echo "✗ Failed: $example"
    fi
    
    cd - > /dev/null
    echo ""
done

echo "===========> Category '$CATEGORY' execution completed!"