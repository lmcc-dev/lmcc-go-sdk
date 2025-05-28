/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This documentation was collaboratively developed by Martin and AI Assistant.
 */

# LMCC Go SDK Examples

[中文版本](README_zh.md)

This directory contains comprehensive examples demonstrating the usage of the three core modules in the LMCC Go SDK:

- **Config**: Configuration management with hot-reload support
- **Errors**: Enhanced error handling with stack traces and error codes  
- **Log**: Structured logging with context support

## Directory Structure

```
examples/
├── README.md                       # This file
├── README_zh.md                    # Chinese version
├── basic-usage/                    # Basic usage of all three modules
│   ├── main.go
│   ├── config.yaml
│   └── README.md
├── config-features/                # Configuration module demonstrations
│   ├── 01-simple-config/           # Simple configuration loading
│   ├── 02-hot-reload/              # Hot-reload demonstration
│   ├── 03-env-override/            # Environment variable override
│   ├── 04-default-values/          # Default values demonstration
│   ├── 05-multiple-formats/        # Multiple file format support
│   └── README.md
├── error-handling/                 # Error handling examples
│   ├── 01-basic-errors/            # Basic error creation
│   ├── 02-error-wrapping/          # Error wrapping
│   ├── 03-error-codes/             # Error codes usage
│   ├── 04-stack-traces/            # Stack trace demonstration
│   ├── 05-error-groups/            # Error aggregation
│   └── README.md
├── logging-features/               # Logging module examples
│   ├── 01-basic-logging/           # Basic logging
│   ├── 02-structured-logging/      # Structured logging
│   ├── 03-context-logging/         # Context-aware logging
│   ├── 04-log-rotation/            # Log rotation
│   ├── 05-multiple-outputs/        # Multiple output destinations
│   └── README.md
├── integration/                    # Integration examples
│   ├── web-app/                    # Web application example
│   ├── microservice/               # Microservice example
│   ├── cli-tool/                   # Command-line tool example
│   └── README.md
└── simple-config-app/              # Legacy example (deprecated)
    ├── main.go
    ├── config.yaml
    └── README.md
```

## Quick Start

### 1. Basic Usage
Start with the `basic-usage/` example to see how all three modules work together:

```bash
cd examples/basic-usage
go run main.go
```

### 2. Module-Specific Examples
Explore individual module features:

```bash
# Configuration examples
cd examples/config-features/01-simple-config
go run main.go

# Error handling examples  
cd examples/error-handling/01-basic-errors
go run main.go

# Logging examples
cd examples/logging-features/01-basic-logging
go run main.go
```

### 3. Integration Examples
See real-world usage patterns:

```bash
# Web application
cd examples/integration/web-app
go run main.go

# CLI tool
cd examples/integration/cli-tool
go run main.go --help
```

## Prerequisites

- Go 1.21 or later
- Basic understanding of Go modules

## Installation

Each example is self-contained. To run any example:

1. Navigate to the example directory
2. Run `go mod tidy` (if needed)
3. Run `go run main.go`

## Examples Overview

### Basic Usage (`basic-usage/`)
Demonstrates the fundamental integration of all three modules:
- Configuration loading with defaults
- Error handling with proper wrapping
- Structured logging with configuration

### Configuration Features (`config-features/`)
Shows various configuration capabilities:
- **01-simple-config**: Basic configuration file loading
- **02-hot-reload**: Real-time configuration updates
- **03-env-override**: Environment variable precedence
- **04-default-values**: Default value mechanisms
- **05-multiple-formats**: YAML, JSON, TOML support

### Error Handling (`error-handling/`)
Demonstrates error management patterns:
- **01-basic-errors**: Creating and formatting errors
- **02-error-wrapping**: Adding context to errors
- **03-error-codes**: Using typed error codes
- **04-stack-traces**: Stack trace capture and formatting
- **05-error-groups**: Aggregating multiple errors

### Logging Features (`logging-features/`)
Shows logging capabilities:
- **01-basic-logging**: Different log levels and formats
- **02-structured-logging**: Key-value pair logging
- **03-context-logging**: Context-aware logging
- **04-log-rotation**: File rotation configuration
- **05-multiple-outputs**: Console and file output

### Integration Examples (`integration/`)
Real-world application patterns:
- **web-app**: HTTP server with middleware
- **microservice**: gRPC service with observability
- **cli-tool**: Command-line application

## Best Practices Demonstrated

1. **Configuration Management**
   - Use struct tags for defaults
   - Implement hot-reload for dynamic services
   - Use environment variables for deployment flexibility

2. **Error Handling**
   - Always wrap errors with context
   - Use typed error codes for API responses
   - Preserve stack traces for debugging

3. **Logging**
   - Use structured logging for machine parsing
   - Include context in all log messages
   - Configure appropriate log levels for environments

## Contributing

To add a new example:

1. Create a new directory in the appropriate category
2. Include a `main.go` with proper documentation
3. Add configuration files if needed
4. Write a comprehensive `README.md`
5. Update this main README.md

## Support

For questions or issues with these examples, please refer to:
- [Configuration Module Documentation](../docs/usage-guides/config/)
- [Error Handling Documentation](../docs/usage-guides/errors/)
- [Logging Documentation](../docs/usage-guides/log/)

## License

These examples are part of the LMCC Go SDK and follow the same license terms. 