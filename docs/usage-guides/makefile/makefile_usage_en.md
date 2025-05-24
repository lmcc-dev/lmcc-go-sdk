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
    -   Runs a sequence of common tasks: `format`, `lint`, `test-unit`, and `tidy`. This is useful for a quick check before committing code. (Note: `test-unit` is run by default, not all tests).

-   **`make format`**
    -   Formats Go source code using standard tools (`gofmt`, `goimports`).
    -   Automatically checks if `goimports` is installed and installs it if missing (via `make tools.verify.goimports`).

-   **`make lint`**
    -   Runs code linters to check for style issues and potential errors.
    -   Currently uses `go vet` and `golangci-lint`.
    -   Automatically checks if `golangci-lint` is installed and installs it if missing (via `make tools.verify.golangci-lint`).
    -   *Note:* You might need to configure `.golangci.yaml` for project-specific rules.

-   **`make test-unit [PKG=...] [RUN=...]`**
    -   Runs unit tests.
    -   Includes the `-race` flag to detect race conditions.
    -   By default, runs unit tests in all packages (excluding `examples/`, `vendor/`, and paths containing `test/integration`).
    -   **Optional `PKG`**: Specifies the package(s) for unit tests. Use relative paths (e.g., `make test-unit PKG=./pkg/log`). If `PKG` is specified, only tests within that package are run.
    -   **Optional `RUN`**: Filters tests to run based on a regular expression matching the test function name (e.g., `make test-unit RUN=TestMyFunction`, `make test-unit PKG=./pkg/log RUN=^TestLog`).

-   **`make test-integration [RUN=...]`**
    -   Runs integration tests. These are typically located in `test/integration/`.
    -   Includes the `-race` flag.
    -   The `PKG` parameter is not typically used with this target as it's designed to run all integration tests.
    -   **Optional `RUN`**: Filters integration tests to run based on a regular expression.

-   **`make cover [PKG=...]`**
    -   Runs unit tests (similar to `test-unit`) and generates code coverage reports.
    -   Saves a text profile to `_output/coverage/coverage.out`.
    -   Saves an HTML report to `_output/coverage/coverage.html` which can be opened in a browser for detailed analysis.
    -   **Optional `PKG`**: Specifies the package(s) for coverage. If `PKG` is specified, coverage is generated only for those packages. Otherwise, it covers all packages subject to unit testing.
    -   *Note:* This target focuses on unit test coverage.

-   **`make tidy`**
    -   Runs `go mod tidy` to ensure the `go.mod` and `go.sum` files are consistent with the source code dependencies.

-   **`make clean`**
    -   Removes generated files, including the `_output` directory (which contains build artifacts, coverage reports) and the downloaded tools directory (`_output/tools`).

-   **`make tools [GOLANGCI_LINT_STRATEGY=...]`**
    -   Installs all required development tools listed in `scripts/make-rules/tools.mk` (currently `golangci-lint`, `goimports`, and `godoc`). This is useful for setting up the development environment initially.
    -   **golangci-lint Version Strategies**:
        - `stable` (default): Uses tested stable version v1.64.8 for reproducible builds
        - `latest`: Always uses the latest available version for cutting-edge features
        - `auto`: Auto-selects based on Go version (Go 1.24+ uses latest, older versions use stable)
    -   Examples:
        - `make tools` (uses stable strategy)
        - `make tools GOLANGCI_LINT_STRATEGY=latest`
        - `make tools GOLANGCI_LINT_STRATEGY=auto`
        - `make tools GOLANGCI_LINT_VERSION=v1.65.0` (override specific version)

-   **`make tools.version`**
    -   Shows versions of all installed tools and current strategy settings.

-   **`make tools.help`**
    -   Shows detailed help for tools management strategies and usage examples.

### `make doc-view PKG=./path/to/package`

Displays the Go documentation for the specified package in the terminal. The `PKG` variable is **required**.

*   **`PKG`**: Specifies the path to the Go package for which to display documentation (e.g., `./pkg/log`, `./pkg/config`).

Example: `make doc-view PKG=./pkg/config`

### `make doc-serve`

Starts a local `godoc` HTTP server (typically on `http://localhost:6060`) to browse HTML documentation for all packages in your Go workspace, including your current project. It will automatically install `godoc` if it's not already present.

Press `Ctrl+C` in the terminal to stop the server.

## 3. Customization (Optional)

-   **Adding Tools**: Edit `scripts/make-rules/tools.mk` to add new tools to `TOOLS_REQUIRED` and provide corresponding `install.<toolname>` rules.
-   **Verbose Output**: Run any command with `V=1` for more detailed output (e.g., `make test V=1`). 