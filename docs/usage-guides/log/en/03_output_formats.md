# Output Formats

This document provides detailed information about the three output formats supported by the log module: text, JSON, and key-value formats.

## Format Overview

The log module supports three main output formats:

1. **Text** - Human-readable format, suitable for development and debugging
2. **JSON** - Structured format, suitable for production environments and log aggregation
3. **KeyValue** - Traditional key-value pair format, suitable for legacy log systems

## Text Format

### Basic Structure

The text format provides human-readable log output, particularly suitable for development environments.

```go
opts := &log.Options{
    Format: "text",
    EnableColor: true,  // Enable color in terminal
}
log.Init(opts)
```

### Output Examples

```go
log.Info("Application started")
log.Infow("User login", "user_id", 123, "username", "john_doe")
log.Errorw("Database connection failed", "error", "connection timeout", "retry_count", 3)
```

**Output:**
```
2024-01-15T10:30:45.123Z	INFO	main.go:25	Application started
2024-01-15T10:30:45.124Z	INFO	main.go:26	User login	{"user_id": 123, "username": "john_doe"}
2024-01-15T10:30:45.125Z	ERROR	main.go:27	Database connection failed	{"error": "connection timeout", "retry_count": 3}
```

### Color Support

In color-capable terminals, different log levels display in different colors:

- **DEBUG** - Blue
- **INFO** - Green
- **WARN** - Yellow
- **ERROR** - Red
- **FATAL** - Red (bold)
- **PANIC** - Red (bold)

```go
opts := &log.Options{
    Format:      "text",
    EnableColor: true,
}
```

### Caller Information

Text format can include caller information (filename and line number):

```go
opts := &log.Options{
    Format:        "text",
    DisableCaller: false,  // Show caller information
}
```

**Output:**
```
2024-01-15T10:30:45.123Z	INFO	main.go:25	Message content
```

## JSON Format

### Basic Structure

JSON format provides structured log output, ideal for production environments and automated log processing.

```go
opts := &log.Options{
    Format: "json",
}
log.Init(opts)
```

### Output Examples

```go
log.Info("Application started")
log.Infow("User login", "user_id", 123, "username", "john_doe", "ip", "192.168.1.100")
log.Errorw("Database connection failed", "error", "connection timeout", "retry_count", 3, "duration", "5.2s")
```

**Output:**
```json
{"level":"info","timestamp":"2024-01-15T10:30:45.123Z","caller":"main.go:25","message":"Application started"}
{"level":"info","timestamp":"2024-01-15T10:30:45.124Z","caller":"main.go:26","message":"User login","user_id":123,"username":"john_doe","ip":"192.168.1.100"}
{"level":"error","timestamp":"2024-01-15T10:30:45.125Z","caller":"main.go:27","message":"Database connection failed","error":"connection timeout","retry_count":3,"duration":"5.2s"}
```

### JSON Field Description

Standard JSON logs contain the following fields:

- **level** - Log level (debug, info, warn, error, fatal, panic)
- **timestamp** - ISO 8601 formatted timestamp
- **caller** - Caller information (filename:line)
- **message** - Log message
- **other fields** - Structured fields added through `log.Infow()` etc.

### Context Fields

When using context logging, context fields are automatically included in JSON output:

```go
ctx := context.Background()
ctx = log.WithValues(ctx, "request_id", "req-123", "user_id", 456)

log.InfoContext(ctx, "Processing request")
```

**Output:**
```json
{"level":"info","timestamp":"2024-01-15T10:30:45.123Z","caller":"main.go:25","message":"Processing request","request_id":"req-123","user_id":456}
```

### Error Stack Traces

In JSON format, stack traces are included as string fields:

```go
opts := &log.Options{
    Format:            "json",
    DisableStacktrace: false,
    StacktraceLevel:   "error",
}
```

**Output:**
```json
{"level":"error","timestamp":"2024-01-15T10:30:45.123Z","caller":"main.go:25","message":"Critical error","stacktrace":"goroutine 1 [running]:\nmain.main()\n\t/path/to/main.go:25 +0x123"}
```

## KeyValue Format

### Basic Structure

Key-value format provides traditional log format, suitable for integration with existing log systems.

```go
opts := &log.Options{
    Format: "keyvalue",
}
log.Init(opts)
```

### Output Examples

```go
log.Info("Application started")
log.Infow("User login", "user_id", 123, "username", "john_doe")
log.Errorw("Database connection failed", "error", "connection timeout", "retry_count", 3)
```

**Output:**
```
timestamp=2024-01-15T10:30:45.123Z level=info caller=main.go:25 message="Application started"
timestamp=2024-01-15T10:30:45.124Z level=info caller=main.go:26 message="User login" user_id=123 username=john_doe
timestamp=2024-01-15T10:30:45.125Z level=error caller=main.go:27 message="Database connection failed" error="connection timeout" retry_count=3
```

### Field Format Rules

In key-value format:

- **String values** - Surrounded by quotes if they contain spaces or special characters
- **Numeric values** - Output directly without quotes
- **Boolean values** - Output as `true` or `false`
- **Special characters** - Special characters in values are escaped

### Complex Value Handling

For complex data types, key-value format performs appropriate serialization:

```go
log.Infow("User information",
    "user", map[string]interface{}{
        "id":   123,
        "name": "John Doe",
        "tags": []string{"admin", "active"},
    },
    "timestamp", time.Now(),
)
```

**Output:**
```
timestamp=2024-01-15T10:30:45.123Z level=info caller=main.go:25 message="User information" user="{\"id\":123,\"name\":\"John Doe\",\"tags\":[\"admin\",\"active\"]}" timestamp=2024-01-15T10:30:45.123Z
```

## Format Comparison

### Readability

| Format | Human Readability | Machine Readability | Use Case |
|--------|-------------------|---------------------|----------|
| Text | ⭐⭐⭐⭐⭐ | ⭐⭐ | Development, debugging |
| JSON | ⭐⭐ | ⭐⭐⭐⭐⭐ | Production, log aggregation |
| KeyValue | ⭐⭐⭐ | ⭐⭐⭐⭐ | Legacy systems, monitoring |

### Performance

| Format | Serialization Speed | File Size | Parse Speed |
|--------|-------------------|-----------|-------------|
| Text | Fast | Medium | Slow |
| JSON | Medium | Large | Fast |
| KeyValue | Fast | Small | Medium |

### Storage Efficiency

```go
// Size comparison for the same log message in different formats
message := "User login"
fields := map[string]interface{}{
    "user_id": 123,
    "username": "john_doe",
    "ip": "192.168.1.100",
    "timestamp": time.Now(),
}

// Text: ~120 bytes
// JSON: ~180 bytes  
// KeyValue: ~100 bytes
```

## Custom Formatting

### Timestamp Format

While timestamp format cannot be directly customized, it can be adjusted through post-processing:

```go
// All formats use ISO 8601 format
// 2024-01-15T10:30:45.123Z
```

### Field Order

In JSON and KeyValue formats, field order is fixed:

1. Standard fields (level, timestamp, caller, message)
2. Context fields (in order of addition)
3. Structured fields (in order of addition)

## Format Selection Guide

### Development Environment

```go
opts := &log.Options{
    Level:       "debug",
    Format:      "text",
    OutputPaths: []string{"stdout"},
    EnableColor: true,
}
```

**Advantages:**
- Easy to read and debug
- Color coding improves readability
- Quick problem identification

### Production Environment

```go
opts := &log.Options{
    Level:        "info",
    Format:       "json",
    OutputPaths:  []string{"/var/log/app.log"},
    EnableColor:  false,
}
```

**Advantages:**
- Structured data for easy querying
- Compatible with log aggregation systems
- Supports complex filtering and analysis

### Legacy System Integration

```go
opts := &log.Options{
    Level:       "warn",
    Format:      "keyvalue",
    OutputPaths: []string{"/var/log/syslog"},
}
```

**Advantages:**
- Compatible with existing tools
- Compact output format
- Easy to parse and process

## Format Conversion

### Runtime Format Switching

```go
// Can dynamically switch formats at runtime
func switchToJSONFormat() {
    opts := &log.Options{
        Level:  "info",
        Format: "json",
    }
    log.Init(opts)
}

func switchToTextFormat() {
    opts := &log.Options{
        Level:       "debug",
        Format:      "text",
        EnableColor: true,
    }
    log.Init(opts)
}
```

### Multi-Format Output

While a single logger instance can only use one format, multiple outputs can be configured:

```go
// This would need to be implemented at the application level
// Example: output to console (text) and file (json) simultaneously
```

## Best Practices

### 1. Environment-Specific Formats

```go
func getLogFormat() string {
    env := os.Getenv("APP_ENV")
    switch env {
    case "development":
        return "text"
    case "production":
        return "json"
    case "testing":
        return "keyvalue"
    default:
        return "text"
    }
}
```

### 2. Performance Optimization

```go
// High-performance scenario: use keyvalue format
opts := &log.Options{
    Level:             "error",
    Format:            "keyvalue",
    DisableCaller:     true,
    DisableStacktrace: true,
}
```

### 3. Debug-Friendly

```go
// Debug scenario: use text format
opts := &log.Options{
    Level:       "debug",
    Format:      "text",
    EnableColor: true,
}
```

### 4. Production Monitoring

```go
// Production monitoring: use json format
opts := &log.Options{
    Level:  "info",
    Format: "json",
    OutputPaths: []string{
        "/var/log/app.log",
        "stdout",  // For container log collection
    },
}
```

## Next Steps

- [Context Logging](04_context_logging.md) - Master context-aware logging
- [Performance](05_performance.md) - Optimize logging performance
- [Best Practices](06_best_practices.md) - Production-ready patterns 