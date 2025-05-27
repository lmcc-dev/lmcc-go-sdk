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