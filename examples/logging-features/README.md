/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This documentation was collaboratively developed by Martin and AI Assistant.
 */

# Logging Features Examples

[中文版本](README_zh.md)

This directory contains comprehensive examples demonstrating the logging capabilities of the LMCC Go SDK. Each example focuses on different aspects of the logging system, from basic usage to advanced integration patterns.

## Examples Overview

### 01-basic-logging
**Basic Logging Operations**

Demonstrates fundamental logging operations including:
- Different log levels (Debug, Info, Warn, Error)
- Structured logging with key-value pairs
- Context-aware logging
- Performance logging and timing
- Error logging patterns

```bash
cd 01-basic-logging
go run main.go
```

**Key Features:**
- UserService with authentication flows
- Batch operation logging
- Context propagation
- Performance measurement

### 02-structured-logging
**JSON Format and Structured Fields**

Shows how to implement structured logging for machine-readable logs:
- JSON output format
- Complex data structure logging
- HTTP request/response logging
- Business event tracking
- Performance metrics logging

```bash
cd 02-structured-logging
go run main.go
```

**Key Features:**
- RequestInfo and ResponseInfo structures
- Database operation logging
- Business event patterns
- Large object serialization

### 03-log-levels
**Log Level Control and Dynamic Adjustment**

Demonstrates log level management and filtering:
- Dynamic log level adjustment
- Component-specific log levels
- Performance impact analysis
- Production environment scenarios
- Conditional logging based on thresholds

```bash
cd 03-log-levels
go run main.go
```

**Key Features:**
- Level filtering demonstration
- Performance benchmarking
- Component isolation
- Advanced level usage patterns

### 04-custom-formatters
**Custom Log Formatters and Log Rotation**

Explores different log output formats, customization, and file rotation:
- Text format with colors
- JSON format for structured data
- Key-value format for parsing
- Performance comparison between formats
- Production environment formatting
- **File rotation strategies (size-based, time-based, combined)**

```bash
cd 04-custom-formatters
go run main.go
```

**Key Features:**
- Format comparison and benchmarking
- Environment-specific configurations
- Custom field formatting
- Large data object handling
- **Size-based, time-based, and combined log rotation**
- **Rotation configuration best practices**
- **Disk space management demonstrations**

### 05-integration-patterns
**Logging System Integration Patterns**

Comprehensive example showing real-world integration patterns:
- HTTP middleware logging
- Configuration-driven logging setup
- Service layer integration
- Error handling integration
- Cross-cutting concerns

```bash
cd 05-integration-patterns
go run main.go
```

**Key Features:**
- LoggingMiddleware for HTTP requests
- Configuration integration
- Multi-layer error handling
- Context propagation across services
- Real-world service patterns

## Common Patterns Demonstrated

### 1. Context-Aware Logging
All examples show how to propagate context information (like request IDs, user IDs) through the logging system:

```go
logger := log.Std().WithValues("request_id", requestID, "user_id", userID)
logger.Infow("Processing request", "operation", "user_creation")
```

### 2. Structured Fields
Examples demonstrate consistent use of structured fields for better log parsing:

```go
logger.Infow("Database operation completed",
    "operation", "INSERT",
    "table", "users",
    "duration", duration,
    "rows_affected", rowsAffected)
```

### 3. Error Integration
Shows integration with the lmcc-go-sdk errors package:

```go
if err != nil {
    wrappedErr := errors.Wrap(err, "operation failed")
    logger.Errorw("Database error", "error", wrappedErr)
    return wrappedErr
}
```

### 4. Performance Logging
Demonstrates timing and performance measurement patterns:

```go
start := time.Now()
defer func() {
    duration := time.Since(start)
    logger.Infow("Operation completed", "duration", duration)
}()
```

## Configuration Examples

Each example shows different configuration approaches:

### Development Configuration
```go
opts := log.NewOptions()
opts.Level = "debug"
opts.Format = "text"
opts.EnableColor = true
opts.DisableCaller = false
```

### Production Configuration
```go
opts := log.NewOptions()
opts.Level = "info"
opts.Format = "json"
opts.EnableColor = false
opts.DisableCaller = false
opts.DisableStacktrace = true
```

### Integration with Config Package
```go
type AppConfig struct {
    Logging struct {
        Level  string `yaml:"level" default:"info"`
        Format string `yaml:"format" default:"json"`
        Output string `yaml:"output" default:"stdout"`
    } `yaml:"logging"`
}
```

## Best Practices Demonstrated

1. **Use structured logging**: Always prefer key-value pairs over string formatting
2. **Include context**: Add request IDs, user IDs, and other context information
3. **Measure performance**: Log operation durations for monitoring
4. **Handle errors properly**: Use the errors package for error wrapping and context
5. **Configure for environment**: Use different settings for development vs production
6. **Use appropriate levels**: Debug for development, Info for normal operations, Warn for concerning situations, Error for actual problems

## Running All Examples

To run all examples in sequence:

```bash
# Run each example
for dir in 01-basic-logging 02-structured-logging 03-log-levels 04-custom-formatters 05-integration-patterns; do
    echo "=== Running $dir ==="
    cd $dir
    go run main.go
    cd ..
    echo
done
```

## Integration with Other Modules

These examples demonstrate integration with:
- **config package**: For configuration-driven logging setup
- **errors package**: For enhanced error handling and logging
- **HTTP middleware**: For request/response logging
- **Database operations**: For query and transaction logging

## Next Steps

After exploring these examples, consider:
1. Implementing similar patterns in your applications
2. Customizing log formats for your specific needs
3. Setting up log aggregation and monitoring
4. Integrating with observability platforms
5. Creating your own middleware and service patterns

For more information, see:
- [Log Package Documentation](../../docs/usage-guides/log/)
- [Config Package Documentation](../../docs/usage-guides/config/)
- [Errors Package Documentation](../../docs/usage-guides/errors/) 