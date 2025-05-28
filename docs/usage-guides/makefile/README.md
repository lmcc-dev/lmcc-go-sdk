# Makefile Documentation

This directory contains comprehensive documentation for the `lmcc-go-sdk` Makefile system.

## 📖 Documentation Files

- **[English Guide](makefile_usage_en.md)** - Complete Makefile usage guide in English
- **[中文指南](makefile_usage_zh.md)** - Complete Makefile usage guide in Chinese

## 🚀 Quick Reference

### Core Development Commands
```bash
make help              # Show all available commands
make all               # Format, lint, test, and tidy (default)
make format            # Format Go source code
make lint              # Run code linters
make test-unit         # Run unit tests
make test-integration  # Run integration tests
make cover             # Generate coverage reports
make tidy              # Tidy go modules
make clean             # Remove build artifacts
```

### Examples Management (19 examples across 5 categories)
```bash
make examples-list                           # List all examples
make examples-build                          # Build all examples
make examples-run EXAMPLE=basic-usage       # Run specific example
make examples-run-all                       # Run all examples
make examples-test                          # Test all examples (lint + build)
make examples-debug EXAMPLE=basic-usage     # Debug with delve
make examples-analyze                       # Analyze examples structure
make examples-category CATEGORY=config-features  # Run examples by category
make examples-clean                         # Clean example binaries
```

### Documentation Commands
```bash
make doc-view PKG=./pkg/log    # View package docs in terminal
make doc-serve                 # Start local documentation server
```

### Tool Management
```bash
make tools                     # Install all required tools
make tools.version            # Show tool versions
make tools.help               # Show tools help
```

## 📊 Examples Overview

The project includes **19 examples** organized into **5 categories**:

| Category | Count | Description |
|----------|-------|-------------|
| `basic-usage` | 1 | Basic integration examples |
| `config-features` | 5 | Configuration module demonstrations |
| `error-handling` | 5 | Error handling patterns |
| `integration` | 3 | Full integration scenarios |
| `logging-features` | 5 | Logging module features |

## 🎯 Common Workflows

### Development Workflow
```bash
# 1. Format and lint code
make format lint

# 2. Run tests
make test-unit

# 3. Check coverage
make cover

# 4. Test examples
make examples-test

# 5. Final check before commit
make all
```

### Example Exploration
```bash
# 1. See what's available
make examples-list

# 2. Try a basic example
make examples-run EXAMPLE=basic-usage

# 3. Explore a category
make examples-category CATEGORY=config-features

# 4. Debug if needed
make examples-debug EXAMPLE=config-features/01-simple-config
```

### Documentation Generation
```bash
# 1. View specific package docs
make doc-view PKG=./pkg/config

# 2. Start documentation server for browsing
make doc-serve
# Open http://localhost:6060 in browser
```

## 🔧 Variables

### Package Selection
- `PKG=./path/to/package` - Specify target package for testing/coverage

### Example Management
- `EXAMPLE=<name>` - Specify example to run/debug
- `CATEGORY=<name>` - Specify example category to run

### Tool Configuration
- `GOLANGCI_LINT_STRATEGY=stable|latest|auto` - Linter version strategy
- `V=1` - Enable verbose output

## 📁 Directory Structure

```
scripts/make-rules/
├── common.mk          # Common variables and functions
├── golang.mk          # Go-specific build rules
└── tools.mk           # Tool installation and management

scripts/
├── examples-build.sh     # Build all examples
├── examples-run-all.sh   # Run all examples sequentially
├── examples-analyze.sh   # Analyze examples structure
├── examples-category.sh  # Run examples by category
└── format-func-coverage.sh  # Format coverage output
```

## 🎪 Features

### ✅ Code Quality
- Automatic code formatting (`gofmt`, `goimports`)
- Comprehensive linting (`golangci-lint`)
- Race condition detection
- Code coverage reporting

### ✅ Testing
- Unit tests with package selection
- Integration tests
- Coverage reports (text + HTML)
- Function-level coverage details

### ✅ Examples Management
- Automatic discovery of examples
- Parallel building
- Category-based execution
- Interactive debugging
- Comprehensive analysis

### ✅ Documentation
- Terminal documentation viewing
- Local documentation server
- Automatic tool installation

### ✅ Development Tools
- Automatic tool installation
- Version management strategies
- Tool version reporting

## 📚 Learn More

For detailed usage instructions and examples, see:
- [English Documentation](makefile_usage_en.md)
- [中文文档](makefile_usage_zh.md)

## 🤝 Contributing

When adding new Makefile targets:
1. Add appropriate help text with `## comment`
2. Follow existing naming conventions
3. Update documentation
4. Test thoroughly

For adding new examples:
1. Create directory under `examples/`
2. Include `main.go` file
3. Add category-appropriate documentation
4. Examples are automatically detected 