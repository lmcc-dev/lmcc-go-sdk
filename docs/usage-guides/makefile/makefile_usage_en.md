# Makefile Usage Guide

This guide explains how to use the `Makefile` included in the `lmcc-go-sdk` project to manage common development tasks.

## 1. Introduction

The `Makefile` provides a standardized way to build, test, lint, format, and clean the project. It leverages includes from the `scripts/make-rules/` directory for modularity, inspired by `marmotedu` practices.

Core goals:
- Simplify common development workflows.
- Ensure consistency in building and testing.
- Automate code quality checks and formatting.
- Manage development tool dependencies.
- Provide centralized examples management.

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

## 3. Examples Management

The project includes a comprehensive examples management system that allows you to build, run, test, and debug the 19 included examples across 5 categories.

### Available Categories

- **basic-usage**: Basic integration examples (1 example)
- **config-features**: Configuration module demonstrations (5 examples)
- **error-handling**: Error handling patterns (5 examples) 
- **integration**: Full integration scenarios (3 examples)
- **logging-features**: Logging module features (5 examples)

### Examples Commands

-   **`make examples-list`**
    -   Lists all available examples with their numbers and categories.
    -   Shows total count of examples.

-   **`make examples-build`**
    -   Builds all examples into the `_output/examples/` directory.
    -   Each example becomes a standalone executable.
    -   Supports parallel building for better performance.

-   **`make examples-clean`**
    -   Removes all built example binaries from `_output/examples/`.
    -   Useful for cleaning up build artifacts.

-   **`make examples-run EXAMPLE=<name>`**
    -   Runs a specific example by name.
    -   **Required `EXAMPLE`**: Specifies which example to run (e.g., `basic-usage`, `config-features/01-simple-config`).
    -   Validates example existence before running.
    -   Examples:
        - `make examples-run EXAMPLE=basic-usage`
        - `make examples-run EXAMPLE=config-features/01-simple-config`
        - `make examples-run EXAMPLE=error-handling/02-error-wrapping`

-   **`make examples-run-all`**
    -   Runs all examples sequentially with progress tracking.
    -   Shows execution status for each example.
    -   Useful for comprehensive testing of all examples.

-   **`make examples-test`**
    -   Performs comprehensive testing of all examples:
        - Step 1: Lints all example code using `golangci-lint`
        - Step 2: Builds all examples to verify compilation
    -   Ensures code quality and buildability of examples.

-   **`make examples-debug EXAMPLE=<name>`**
    -   Starts an interactive debugging session for a specific example using `delve`.
    -   **Required `EXAMPLE`**: Specifies which example to debug.
    -   Automatically installs `dlv` if not present.
    -   Shows helpful debugger commands on start.
    -   Examples:
        - `make examples-debug EXAMPLE=basic-usage`
        - `make examples-debug EXAMPLE=config-features/02-hot-reload`

-   **`make examples-analyze`**
    -   Provides comprehensive analysis of the examples structure:
        - Directory tree visualization
        - Statistics by category
        - Code metrics (total lines, Go files count)
    -   Useful for understanding project structure and scope.

-   **`make examples-category CATEGORY=<name>`**
    -   Runs all examples within a specific category.
    -   **Required `CATEGORY`**: Specifies which category to run.
    -   Available categories: `basic-usage`, `config-features`, `error-handling`, `integration`, `logging-features`
    -   Examples:
        - `make examples-category CATEGORY=config-features`
        - `make examples-category CATEGORY=error-handling`

-   **`make examples`**
    -   Default examples target that runs `examples-test`.
    -   Equivalent to running lint and build verification for all examples.

### Examples Usage Tips

1. **Start with listing**: Use `make examples-list` to see all available examples.
2. **Run individual examples**: Use `make examples-run EXAMPLE=<name>` to test specific functionality.
3. **Category exploration**: Use `make examples-category CATEGORY=<name>` to explore related examples.
4. **Development workflow**: Use `make examples-test` to ensure all examples are working before commits.
5. **Debugging**: Use `make examples-debug EXAMPLE=<name>` when you need to step through example code.

## 4. Customization (Optional)

-   **Adding Tools**: Edit `scripts/make-rules/tools.mk` to add new tools to `TOOLS_REQUIRED` and provide corresponding `install.<toolname>` rules.
-   **Verbose Output**: Run any command with `V=1` for more detailed output (e.g., `make test V=1`).
-   **Adding Examples**: Create new examples in the `examples/` directory with a `main.go` file, and they will automatically be detected by the examples management system. 