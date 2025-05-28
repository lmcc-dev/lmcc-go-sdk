# LMCC Go SDK

[![Go Report Card](https://goreportcard.com/badge/github.com/lmcc-dev/lmcc-go-sdk)](https://goreportcard.com/report/github.com/lmcc-dev/lmcc-go-sdk)
[![Go Reference](https://pkg.go.dev/badge/github.com/lmcc-dev/lmcc-go-sdk.svg)](https://pkg.go.dev/github.com/lmcc-dev/lmcc-go-sdk)

A comprehensive Go SDK providing foundational components and utilities for building robust applications.

## Quick Links

- **[中文文档](README_zh.md)** - Chinese documentation
- **[📚 Usage Guides](./docs/usage-guides/)** - Comprehensive module documentation
- **[API Reference](https://pkg.go.dev/github.com/lmcc-dev/lmcc-go-sdk)** - Go package documentation
- **[Examples](./examples/)** - Working code examples

## Features

### 📦 Core Modules
- **Configuration Management**: Hot-reload support with multiple sources
- **Structured Logging**: High-performance logging with multiple formats
- **Error Handling**: Enhanced error handling with codes and stack traces

### 🚀 Developer Experience
- **Type Safety**: Strong typing through user-defined structs
- **Hot Reload**: Dynamic configuration updates without restart
- **Multiple Formats**: JSON, YAML, TOML configuration support
- **Environment Integration**: Automatic environment variable binding

## Quick Example

```go
package main

import (
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

func main() {
	// Initialize logging
	log.Init(nil)
	log.Info("Hello, LMCC Go SDK!")
	
	// Load configuration
	var cfg MyConfig
	err := config.LoadConfig(&cfg)
	if err != nil {
		log.Error("Failed to load config", "error", err)
	}
}
```

## Installation

```bash
go get github.com/lmcc-dev/lmcc-go-sdk
```

## Available Modules

| Module | Description | Documentation |
|--------|-------------|---------------|
| **config** | Configuration management with hot-reload | [📖 Guide](./docs/usage-guides/config/) |
| **log** | High-performance structured logging | [📖 Guide](./docs/usage-guides/log/) |
| **errors** | Enhanced error handling with codes | [📖 Guide](./docs/usage-guides/errors/) |

## Development Tools

This project includes a comprehensive Makefile for development workflows and examples management.

### Quick Commands

```bash
# Development workflow
make help              # Show all available commands
make all               # Format, lint, test, and tidy (recommended before commits)
make format            # Format Go source code
make lint              # Run code linters
make test-unit         # Run unit tests
make cover             # Generate coverage reports

# Examples management (19 examples across 5 categories)
make examples-list                        # List all available examples
make examples-run EXAMPLE=basic-usage    # Run a specific example
make examples-test                       # Test all examples
make examples-build                      # Build all examples
make examples-debug EXAMPLE=basic-usage  # Debug with delve

# Documentation
make doc-serve         # Start local documentation server
make doc-view PKG=./pkg/log  # View package docs in terminal
```

### Examples Categories

The project includes **19 practical examples** across **5 categories**:

- **basic-usage** (1): Basic integration patterns
- **config-features** (5): Configuration management demos
- **error-handling** (5): Error handling patterns  
- **integration** (3): Full integration scenarios
- **logging-features** (5): Logging capabilities

**📖 Complete Makefile Documentation**: [docs/usage-guides/makefile/](./docs/usage-guides/makefile/)

## Getting Started

1. **[Browse all modules](./docs/usage-guides/)** in the usage guides directory
2. **Choose a module** that fits your needs
3. **Follow the Quick Start guide** for that module
4. **Explore advanced features** using the detailed documentation
5. **Check Best Practices** for production-ready patterns

## Contributing

Contributions are welcome! Please see our [contributing guidelines](./CONTRIBUTING.md).

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details. 