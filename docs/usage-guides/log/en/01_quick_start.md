# Quick Start Guide

This guide will help you get started with the log module in just a few minutes.

## Installation

The log module is part of the lmcc-go-sdk. If you haven't already, add it to your project:

```bash
go mod init your-project
go get github.com/lmcc-dev/lmcc-go-sdk
```

## Basic Usage

### Step 1: Initialize the Logger

The simplest way to get started is with the default configuration:

```go
package main

import (
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

func main() {
    // Initialize with default settings
    // - Level: info
    // - Format: text (with colors in console)
    // - Output: stdout
    log.Init(nil)
    
    log.Info("Logger initialized successfully!")
}
```

### Step 2: Basic Logging

```go
func basicLogging() {
    // Different log levels
    log.Debug("This is a debug message")   // Won't show with default 'info' level
    log.Info("Application started")        // Will show
    log.Warn("This is a warning")          // Will show
    log.Error("Something went wrong")      // Will show
    
    // Fatal and Panic (use with caution)
    // log.Fatal("Critical error")  // Calls os.Exit(1)
    // log.Panic("Panic message")   // Calls panic()
}
```

### Step 3: Structured Logging

The real power comes with structured logging using key-value pairs:

```go
func structuredLogging() {
    // Using the 'w' variants for structured logging
    log.Infow("User logged in",
        "user_id", 12345,
        "username", "john_doe",
        "ip", "192.168.1.100",
        "timestamp", time.Now(),
    )
    
    log.Errorw("Database connection failed",
        "error", "connection timeout",
        "database", "users_db",
        "retry_count", 3,
        "duration", "5.2s",
    )
    
    log.Warnw("High memory usage detected",
        "memory_usage", "85%",
        "threshold", "80%",
        "process_id", 1234,
    )
}
```

## Custom Configuration

### Step 1: Configure Output Format

```go
func customConfiguration() {
    opts := &log.Options{
        Level:  "debug",  // Show debug messages
        Format: "json",   // Use JSON format instead of text
    }
    
    log.Init(opts)
    
    log.Debug("This debug message will now appear")
    log.Infow("User action",
        "action", "login",
        "user_id", 123,
    )
}
```

### Step 2: Configure Output Destinations

```go
func multipleOutputs() {
    opts := &log.Options{
        Level:       "info",
        Format:      "json",
        OutputPaths: []string{
            "stdout",              // Console output
            "/var/log/app.log",    // File output
        },
    }
    
    log.Init(opts)
    
    log.Info("This message goes to both console and file")
}
```

### Step 3: Enable Log Rotation

```go
func withLogRotation() {
    opts := &log.Options{
        Level:       "info",
        Format:      "json",
        OutputPaths: []string{"stdout", "/var/log/app.log"},
        
        // Log rotation settings
        LogRotateMaxSize:    100,  // 100 MB per file
        LogRotateMaxBackups: 5,    // Keep 5 backup files
        LogRotateMaxAge:     30,   // Keep files for 30 days
        LogRotateCompress:   true, // Compress old files
    }
    
    log.Init(opts)
    
    log.Info("Logging with rotation enabled")
}
```

## Context-Aware Logging

### Basic Context Usage

```go
import (
    "context"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

func contextLogging() {
    // Create a context with logging fields
    ctx := context.Background()
    ctx = log.WithValues(ctx, "request_id", "req-123", "user_id", 456)
    
    // Log with context - fields are automatically included
    log.InfoContext(ctx, "Processing user request")
    log.WarnContext(ctx, "Request taking longer than expected")
    
    // Add more fields to the context
    ctx = log.WithValues(ctx, "operation", "database_query")
    log.InfoContext(ctx, "Executing database query")
}
```

### HTTP Request Context

```go
import (
    "net/http"
    "github.com/google/uuid"
)

func httpHandler(w http.ResponseWriter, r *http.Request) {
    // Create request context with correlation ID
    requestID := uuid.New().String()
    ctx := log.WithValues(r.Context(), 
        "request_id", requestID,
        "method", r.Method,
        "path", r.URL.Path,
        "remote_addr", r.RemoteAddr,
    )
    
    log.InfoContext(ctx, "Request started")
    
    // Process request...
    processRequest(ctx)
    
    log.InfoContext(ctx, "Request completed")
}

func processRequest(ctx context.Context) {
    // All logs in this function will include the request context
    log.InfoContext(ctx, "Validating request")
    log.InfoContext(ctx, "Querying database")
    log.InfoContext(ctx, "Generating response")
}
```

## Output Format Examples

### JSON Format Output

When using `Format: "json"`, logs look like this:

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

### Text Format Output

When using `Format: "text"` (default), logs look like this:

```
2024-01-15T10:30:45.123Z	INFO	main.go:25	User logged in	{"user_id": 12345, "username": "john_doe", "ip": "192.168.1.100"}
```

### Key=Value Format Output

When using `Format: "keyvalue"`, logs look like this:

```
timestamp=2024-01-15T10:30:45.123Z level=info caller=main.go:25 message="User logged in" user_id=12345 username=john_doe ip=192.168.1.100
```

## Complete Example

Here's a complete example that demonstrates various features:

```go
package main

import (
    "context"
    "time"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

func main() {
    // Configure the logger
    opts := &log.Options{
        Level:            "debug",
        Format:           "json",
        OutputPaths:      []string{"stdout", "app.log"},
        EnableColor:      false,  // Disable colors for file output
        DisableCaller:    false,  // Include caller information
        DisableStacktrace: false, // Include stack traces for errors
        StacktraceLevel:  "error", // Only show stack traces for errors and above
    }
    
    log.Init(opts)
    
    // Basic logging
    log.Info("Application starting up")
    
    // Structured logging
    log.Infow("Configuration loaded",
        "config_file", "app.yaml",
        "log_level", opts.Level,
        "output_format", opts.Format,
    )
    
    // Context logging
    ctx := context.Background()
    ctx = log.WithValues(ctx, "component", "database", "version", "1.0.0")
    
    log.InfoContext(ctx, "Connecting to database")
    
    // Simulate some work
    time.Sleep(100 * time.Millisecond)
    
    log.InfoContext(ctx, "Database connection established")
    
    // Error logging with stack trace
    log.Errorw("Failed to process user data",
        "user_id", 123,
        "error", "validation failed",
        "field", "email",
    )
    
    log.Info("Application startup completed")
}
```

## Integration with Configuration Module

For production applications, you'll often want to load log configuration from files:

```go
import (
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

type AppConfig struct {
    Log log.Options `mapstructure:"log"`
    // ... other configuration
}

func main() {
    var cfg AppConfig
    
    // Load configuration
    err := config.LoadConfig(&cfg,
        config.WithConfigFile("config.yaml", ""),
        config.WithEnvPrefix("APP"),
    )
    if err != nil {
        panic(err)
    }
    
    // Initialize logging with loaded configuration
    log.Init(&cfg.Log)
    
    log.Info("Application started with configuration-based logging")
}
```

With a `config.yaml` file like:

```yaml
log:
  level: "info"
  format: "json"
  output_paths: ["stdout", "/var/log/app.log"]
  enable_color: false
  disable_caller: false
  log_rotate_max_size: 100
  log_rotate_max_backups: 5
  log_rotate_max_age: 30
  log_rotate_compress: true
```

## Best Practices for Getting Started

1. **Start Simple**: Begin with `log.Init(nil)` and add configuration as needed
2. **Use Structured Logging**: Prefer `log.Infow()` over `log.Info()` for better searchability
3. **Include Context**: Use context logging for request tracing and correlation
4. **Choose the Right Level**: Use appropriate log levels (debug for development, info for production)
5. **Configure Rotation**: Always set up log rotation for file outputs in production

## Next Steps

- [Configuration Options](02_configuration_options.md) - Learn about all available options
- [Output Formats](03_output_formats.md) - Understand different output formats
- [Context Logging](04_context_logging.md) - Master context-aware logging
- [Best Practices](06_best_practices.md) - Production-ready patterns 