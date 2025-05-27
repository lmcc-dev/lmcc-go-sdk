# Log Module Overview

The `pkg/log` module provides a comprehensive and high-performance logging solution for Go applications. Built on top of proven logging libraries, it offers structured logging, multiple output formats, context-aware logging, and hot-reload capabilities.

## Key Features

### 🚀 **High Performance**
- Optimized for high-throughput applications
- Minimal allocation overhead
- Configurable log levels to reduce unnecessary processing
- Efficient serialization for different output formats

### 📊 **Multiple Output Formats**
- **Text Format**: Human-readable output for development
- **JSON Format**: Structured output for production and log aggregation
- **Key-Value Format**: Traditional format for legacy systems

### 🎯 **Context-Aware Logging**
- Automatic context propagation through Go's `context.Context`
- Request tracing and correlation IDs
- Structured field inheritance across function calls

### ⚙️ **Flexible Configuration**
- Environment-specific configurations
- Hot-reload support for dynamic updates
- Integration with the `pkg/config` module
- Comprehensive configuration validation

### 🔄 **Log Rotation**
- Automatic log file rotation based on size, age, or count
- Compression of archived log files
- Configurable retention policies

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Application   │───▶│   Log Module     │───▶│   Output        │
│                 │    │                  │    │   Destinations  │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                              │                         │
                              ▼                         ▼
                       ┌──────────────┐         ┌──────────────┐
                       │ Configuration│         │ • stdout     │
                       │ Hot Reload   │         │ • Files      │
                       └──────────────┘         │ • Network    │
                                               └──────────────┘
```

## Core Components

### Logger Interface
The main logging interface provides methods for different log levels and structured logging:

```go
// Basic logging methods
log.Debug("Debug message")
log.Info("Info message")
log.Warn("Warning message")
log.Error("Error message")

// Structured logging with fields
log.Infow("User action", "user_id", 123, "action", "login")

// Context-aware logging
log.InfoContext(ctx, "Request processed")
```

### Configuration System
Comprehensive configuration options for all aspects of logging:

```go
type Options struct {
    Level            string   `mapstructure:"level"`
    Format           string   `mapstructure:"format"`
    OutputPaths      []string `mapstructure:"output_paths"`
    EnableColor      bool     `mapstructure:"enable_color"`
    DisableCaller    bool     `mapstructure:"disable_caller"`
    // ... more options
}
```

### Context Integration
Seamless integration with Go's context system for request tracing:

```go
ctx := log.WithValues(ctx, "request_id", "req-123")
log.InfoContext(ctx, "Processing request") // Automatically includes request_id
```

## Use Cases

### Development Environment
- Human-readable text format with colors
- Debug-level logging for detailed information
- Caller information for easy debugging

### Production Environment
- JSON format for structured logging
- Info-level logging for performance
- Log rotation and compression
- Integration with log aggregation systems

### High-Performance Applications
- Error-level logging only
- Disabled caller information and stack traces
- Key-value format for minimal overhead
- Asynchronous logging patterns

## Integration Points

### With pkg/config
```go
type AppConfig struct {
    Log log.Options `mapstructure:"log"`
}

// Automatic configuration loading and hot-reload
cm, err := config.LoadConfigAndWatch(&cfg, ...)
log.Init(&cfg.Log)
```

### With Web Frameworks
```go
// Gin middleware example
func LoggingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx := log.WithValues(c.Request.Context(),
            "request_id", generateRequestID(),
            "method", c.Request.Method,
            "path", c.Request.URL.Path,
        )
        c.Request = c.Request.WithContext(ctx)
        c.Next()
    }
}
```

### With Error Handling
```go
import "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"

if err := someOperation(); err != nil {
    log.ErrorContext(ctx, "Operation failed",
        "error", err,
        "error_code", errors.GetCoder(err),
    )
}
```

## Performance Characteristics

| Feature | Impact | Recommendation |
|---------|--------|----------------|
| Log Level | High | Use "info" or higher in production |
| Output Format | Medium | JSON for production, text for development |
| Caller Info | Medium | Disable in high-performance scenarios |
| Stack Traces | High | Only for error level and above |
| Context Fields | Low | Efficient for request tracing |

## Getting Started

1. **[Quick Start](01_quick_start.md)** - Get up and running in minutes
2. **[Configuration Options](02_configuration_options.md)** - Detailed configuration guide
3. **[Output Formats](03_output_formats.md)** - Choose the right format
4. **[Context Logging](04_context_logging.md)** - Master request tracing
5. **[Performance](05_performance.md)** - Optimize for your use case
6. **[Best Practices](06_best_practices.md)** - Production-ready patterns

## Next Steps

- Explore the [Quick Start Guide](01_quick_start.md) for immediate hands-on experience
- Review [Configuration Options](02_configuration_options.md) for detailed setup
- Check [Best Practices](06_best_practices.md) for production deployment guidance 