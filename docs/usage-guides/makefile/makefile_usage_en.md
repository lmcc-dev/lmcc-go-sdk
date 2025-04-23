# Makefile Usage Guide

This guide explains how to use the `Makefile` included in the `lmcc-go-sdk` project to manage common development tasks.

## 1. Introduction

The `Makefile` provides a standardized way to build, test, lint, format, and clean the project. It leverages includes from the `scripts/make-rules/` directory for modularity, inspired by `marmotedu` practices.

Core goals:
- Simplify common development workflows.
- Ensure consistency in building and testing.
- Automate code quality checks and formatting.
- Manage development tool dependencies.

## 2. Common Commands

You can run these commands from the root directory of the project.

-   **`make help`**
    -   Displays a help message listing all available targets and their descriptions. This is the best way to see all available commands.

-   **`make all`** (Default Goal)
    -   Runs a sequence of common tasks: `format`, `lint`, `test`, and `tidy`. This is useful for a quick check before committing code.

-   **`make format`**
    -   Formats Go source code using standard tools (`gofmt`, `goimports`).
    -   Automatically checks if `goimports` is installed and installs it if missing (via `make tools.verify.goimports`).

-   **`make lint`**
    -   Runs code linters to check for style issues and potential errors.
    -   Currently uses `go vet` and `golangci-lint`.
    -   Automatically checks if `golangci-lint` is installed and installs it if missing (via `make tools.verify.golangci-lint`).
    -   *Note:* You might need to configure `.golangci.yaml` for project-specific rules.

-   **`make test`**
    -   Runs all unit tests in the project (excluding `examples/`).
    -   Includes the `-race` flag to detect race conditions.

-   **`make cover`**
    -   Runs all unit tests and generates code coverage reports.
    -   Saves a text profile to `_output/coverage/coverage.out`.
    -   Saves an HTML report to `_output/coverage/coverage.html` which can be opened in a browser for detailed analysis.

-   **`make tidy`**
    -   Runs `go mod tidy` to ensure the `go.mod` and `go.sum` files are consistent with the source code dependencies.

-   **`make clean`**
    -   Removes generated files, including the `_output` directory (which contains build artifacts, coverage reports) and the downloaded tools directory (`_output/tools`).

-   **`make tools`**
    -   Installs all required development tools listed in `scripts/make-rules/tools.mk` (currently `golangci-lint` and `goimports`). This is useful for setting up the development environment initially.

## 3. Customization (Optional)

-   **Adding Tools**: Edit `scripts/make-rules/tools.mk` to add new tools to `TOOLS_REQUIRED` and provide corresponding `install.<toolname>` rules.
-   **Verbose Output**: Run any command with `V=1` for more detailed output (e.g., `make test V=1`). 