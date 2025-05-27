# Log Module Documentation

The log module provides a powerful, flexible, and high-performance logging solution for Go applications. Built on top of Zap, it offers structured logging with multiple output formats, log rotation, and context-aware logging capabilities.

## Quick Links

- **[中文文档](README_zh.md)** - Chinese documentation
- **[Quick Start Guide](en/01_quick_start.md)** - Get started in minutes
- **[Configuration Options](en/02_configuration_options.md)** - All available options
- **[Output Formats](en/03_output_formats.md)** - JSON, text, and key=value formats
- **[Context Logging](en/04_context_logging.md)** - Context-aware logging
- **[Log Rotation](en/05_log_rotation.md)** - File rotation and management
- **[Best Practices](en/06_best_practices.md)** - Recommended patterns
- **[Integration Examples](en/07_integration_examples.md)** - Real-world examples
- **[Troubleshooting](en/08_troubleshooting.md)** - Common issues and solutions

## Features

### 🚀 High Performance
- Built on Uber's Zap logger for maximum performance
- Zero-allocation logging in hot paths
- Efficient structured logging with minimal overhead

### 📝 Multiple Output Formats
- **JSON**: Machine-readable structured logs
- **Text**: Human-readable console output with colors
- **Key=Value**: Simple key-value pair format

### 🔄 Flexible Output Destinations
- Console output (stdout/stderr)
- File output with automatic rotation
- Multiple simultaneous outputs
- Custom output destinations

### 🎯 Context-Aware Logging
- Request ID tracking
- User context preservation
- Automatic field inheritance
- Structured context propagation

### ⚙️ Easy Configuration
- YAML/JSON configuration files
- Environment variable overrides
- Hot reload support (with config module)
- Sensible defaults for quick setup

### 🔧 Advanced Features
- Log level filtering per output
- Caller information (file, line, function)
- Stack traces for errors
- Log sampling for high-volume scenarios
- Custom field encoders

## Quick Example

```go
package main

import (
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

func main() {
    // Initialize with default settings
    log.Init(nil)
    
    // Basic logging
    log.Info("Application started")
    log.Warn("This is a warning")
    log.Error("Something went wrong")
    
    // Structured logging
    log.Infow("User logged in",
        "user_id", 12345,
        "username", "john_doe",
        "ip", "192.168.1.100",
    )
    
    // Context logging
    ctx := log.WithContext(context.Background())
    ctx = log.WithValues(ctx, "request_id", "req-123")
    
    log.InfoContext(ctx, "Processing request")
}
```

## Installation

The log module is part of the lmcc-go-sdk:

```bash
go get github.com/lmcc-dev/lmcc-go-sdk
```

## Basic Configuration

### Simple Configuration

```go
import "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"

// Use defaults (console output, info level)
log.Init(nil)
```

### Custom Configuration

```go
opts := &log.Options{
    Level:            "debug",
    Format:           "json",
    OutputPaths:      []string{"stdout", "/var/log/app.log"},
    EnableColor:      false,
    DisableCaller:    false,
    DisableStacktrace: false,
}

log.Init(opts)
```

### YAML Configuration

```yaml
# config.yaml
log:
  level: "info"
  format: "json"
  output_paths: ["stdout", "/var/log/app.log"]
  enable_color: true
  disable_caller: false
  disable_stacktrace: false
  stacktrace_level: "error"
  
  # Log rotation settings
  log_rotate_max_size: 100      # MB
  log_rotate_max_backups: 5     # files
  log_rotate_max_age: 30        # days
  log_rotate_compress: true
```

## Output Format Examples

### JSON Format
```json
{
  "level": "info",
  "timestamp": "2024-01-15T10:30:45.123Z",
  "caller": "main.go:25",
  "message": "User logged in",
  "user_id": 12345,
  "username": "john_doe",
  "ip": "192.168.1.100"
}
```

### Text Format
```
2024-01-15T10:30:45.123Z	INFO	main.go:25	User logged in	{"user_id": 12345, "username": "john_doe", "ip": "192.168.1.100"}
```

### Key=Value Format
```
timestamp=2024-01-15T10:30:45.123Z level=info caller=main.go:25 message="User logged in" user_id=12345 username=john_doe ip=192.168.1.100
```

## Integration with Config Module

The log module integrates seamlessly with the configuration module for hot reload:

```go
import (
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

type AppConfig struct {
    Log log.Options `mapstructure:"log"`
    // ... other config
}

func main() {
    var cfg AppConfig
    
    // Load configuration with hot reload
    cm, err := config.LoadConfigAndWatch(&cfg, 
        config.WithConfigFile("config.yaml", ""),
        config.WithHotReload(true),
    )
    if err != nil {
        panic(err)
    }
    
    // Initialize logging
    log.Init(&cfg.Log)
    
    // Register for hot reload
    log.RegisterConfigHotReload(cm)
    
    log.Info("Application started with hot reload logging")
}
```

## Performance Characteristics

The log module is designed for high-performance applications:

- **Zero allocations** for disabled log levels
- **Minimal allocations** for enabled logs
- **Efficient JSON encoding** using Zap's optimized encoder
- **Buffered I/O** for file outputs
- **Async logging** option for maximum throughput

### Benchmarks

```
BenchmarkLogInfo-8           	 5000000	       230 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogInfow-8          	 3000000	       450 ns/op	      64 B/op	       1 allocs/op
BenchmarkLogJSON-8           	 2000000	       680 ns/op	     128 B/op	       2 allocs/op
```

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Application   │───▶│   Log Module     │───▶│   Zap Logger    │
│                 │    │                  │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                │                        │
                                ▼                        ▼
                       ┌──────────────────┐    ┌─────────────────┐
                       │  Configuration   │    │   Output Sinks  │
                       │    Options       │    │  (Console/File) │
                       └──────────────────┘    └─────────────────┘
```

## Use Cases

### Web Applications
- Request/response logging with correlation IDs
- Error tracking and debugging
- Performance monitoring
- Security audit trails

### Microservices
- Distributed tracing correlation
- Service-to-service communication logging
- Health check and metrics logging
- Configuration change tracking

### CLI Applications
- User action logging
- Error reporting
- Debug information
- Progress tracking

### Background Services
- Job processing logs
- Scheduled task execution
- System monitoring
- Data processing pipelines

## Getting Started

1. **[Quick Start Guide](en/01_quick_start.md)** - Basic setup and usage
2. **[Configuration Options](en/02_configuration_options.md)** - Detailed configuration
3. **[Output Formats](en/03_output_formats.md)** - Choose the right format
4. **[Context Logging](en/04_context_logging.md)** - Advanced context features
5. **[Best Practices](en/06_best_practices.md)** - Production recommendations

## Contributing

We welcome contributions! Please see our [Contributing Guide](../../CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details. 