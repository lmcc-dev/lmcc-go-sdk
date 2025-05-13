#!/bin/bash

# scripts/format-func-coverage.sh
# Description: Filters and formats go tool cover -func output for a specific package.
# Arguments:
#   $1: PKG - The package path (e.g., ./pkg/config or ./...)
#   $2: COVERAGE_FILE - Path to the coverage.out file
#   $3: ROOT_PACKAGE - The root Go package name (e.g., github.com/lmcc-dev/lmcc-go-sdk)

PKG_PATH="$1"
COVERAGE_FILE="$2"
ROOT_PACKAGE_NAME="$3"
GO_CMD="go" # Assuming go is in PATH

# --- Input Validation ---
if [ -z "$PKG_PATH" ] || [ -z "$COVERAGE_FILE" ] || [ -z "$ROOT_PACKAGE_NAME" ]; then
    echo "Usage: $0 <package_path> <coverage_file> <root_package_name>"
    exit 1
fi

if [ ! -f "$COVERAGE_FILE" ]; then
    echo "Error: Coverage file not found: $COVERAGE_FILE"
    echo "Please run 'make cover PKG=$PKG_PATH' first."
    exit 1
fi

# --- Filtering Logic ---
# Construct the full path prefix to grep for
# Remove leading ./ if present
PKG_REL_PATH="${PKG_PATH#./}"
# Handle the ./... case - don't filter if PKG is ./...
if [ "$PKG_REL_PATH" == "..." ]; then
    FILTER_PREFIX="$ROOT_PACKAGE_NAME/"
    echo "Displaying coverage for all packages..."
else
    FILTER_PREFIX="$ROOT_PACKAGE_NAME/$PKG_REL_PATH/"
    echo "Displaying coverage for package: $PKG_PATH ($FILTER_PREFIX)..."
fi

# --- Processing and Formatting ---
# Run go tool cover and filter/format the output
COVER_OUTPUT=$($GO_CMD tool cover -func="$COVERAGE_FILE")
TOTAL_LINE=$(echo "$COVER_OUTPUT" | grep '^total:')

# Use awk to filter, parse, and format into a Markdown table
# Pass TOTAL_LINE as an awk variable using -v
FORMATTED_TABLE=$(echo "$COVER_OUTPUT" | grep "^${FILTER_PREFIX}" | awk -v total_line_val="$TOTAL_LINE" -F'[: \t]+' '
BEGIN {
  print "| 文件 (File) | 函数 (Function) | 覆盖率 (Coverage) |"
  print "| :---------- | :-------------- | :---------------- |"
}
{
  # Extract relevant parts, handle potential path variations
  # $1 is the full path like github.com/lmcc-dev/lmcc-go-sdk/pkg/config/accessors.go
  # We want the relative path from ROOT_PACKAGE_NAME
  sub("'$ROOT_PACKAGE_NAME'/", "", $1); # Remove root package prefix
  filename = $1
  funcname = $3 # Function name is usually the 3rd field after splitting
  coverage = $NF # Coverage percentage is the last field

  # Find the function name; it might shift if filename has spaces (unlikely in Go paths)
  # A more robust way: find the last field (coverage), the field before it is the func name
  for(i=NF-1; i>1; i--) {
      if ($i != "") {
          funcname = $i
          break
      }
  }

  # Print the Markdown table row
  printf "| %s | %s | %s |\n", filename, funcname, coverage
}
END {
  # Optionally print the total line if it exists, using the awk variable
  if (total_line_val != "") {
      print "" # Add a newline before total
      # We can clean up the value within awk if needed, but printing directly is usually fine
      # gsub(/^[ \t]+|[ \t]+$/, "", total_line_val); # Optional: Trim whitespace
      print "**全局总覆盖率 (Global Total Coverage):**"
      print "```"
      print total_line_val # Use the awk variable directly
      print "```"
  }
}
')

# --- Output ---
echo "" # Add a newline before the table
echo "$FORMATTED_TABLE"

exit 0
