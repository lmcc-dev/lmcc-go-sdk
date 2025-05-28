# Basic Usage Example

[中文版本](README_zh.md)

This example demonstrates the fundamental integration of all three core modules in the LMCC Go SDK:
- **Config**: Configuration loading with defaults and environment variable override
- **Errors**: Error creation, wrapping, and code assignment
- **Log**: Structured logging with context support

## What This Example Shows

1. **Configuration Management**
   - Loading configuration from YAML file
   - Using default values from struct tags
   - Environment variable override support
   - Embedding SDK base configuration

2. **Error Handling**
   - Creating errors with automatic stack traces
   - Wrapping errors with additional context
   - Assigning error codes for API responses
   - Extracting error information

3. **Logging Integration**
   - Initializing logger from configuration
   - Structured logging with key-value pairs
   - Context-aware logging
   - Different log levels (info, warn, error)

4. **Business Logic Integration**
   - Retry mechanisms with configuration
   - Error propagation through call stacks
   - Context passing for request tracking

## File Structure

```
basic-usage/
├── main.go         # Main application code
├── config.yaml     # Configuration file
├── README.md       # This file
└── README_zh.md    # Chinese version
```

## Running the Example

### Prerequisites
- Go 1.21 or later
- No external dependencies required

### Run with Default Configuration

```bash
cd examples/basic-usage
go run main.go
```

### Run with Environment Variables

You can override any configuration value using environment variables with the `BASIC_EXAMPLE_` prefix:

```bash
# Override application settings
export BASIC_EXAMPLE_APP_NAME="Custom App Name"
export BASIC_EXAMPLE_APP_MAX_RETRIES=5
export BASIC_EXAMPLE_LOG_LEVEL="debug"

go run main.go
```

### Run with Different Configuration File

```bash
# Create a custom config file
cp config.yaml custom-config.yaml
# Edit custom-config.yaml as needed

# Run with custom config
go run main.go -config custom-config.yaml
```

## Expected Output

The example will output step-by-step execution showing:

1. **Configuration Loading**
   ```
   1. Loading configuration...
   ✓ Configuration loaded successfully
     App Name: LMCC SDK Basic Example
     Version: 1.0.0
     Environment: development
     Timeout: 30s
     Max Retries: 3
   ```

2. **Logger Initialization**
   ```
   2. Initializing logger...
   ✓ Logger initialized
   ```

3. **Normal Operations**
   ```
   4. Demonstrating normal operations...
   ✓ Data processing completed successfully
   ```

4. **Error Handling**
   ```
   5. Demonstrating error handling...
   ✓ Expected error caught: failed to process data: invalid data format detected
     Error Code: 100001
     Error Type: ConfigFileReadError
   ```

5. **Retry Logic**
   ```
   6. Demonstrating retry logic...
   ✓ Simple operation succeeded
   ✓ Flaky operation succeeded after retries
   ```

6. **Structured Logging**
   ```
   7. Demonstrating structured logging...
   ✓ Structured logging demonstrated
   ```

## Key Learning Points

### 1. Configuration Integration
```go
// Embed SDK base config
type AppConfig struct {
    config.Config                    // SDK base
    App *AppSpecificConfig          // Your app config
}

// Load with options
err := config.LoadConfig(&cfg,
    config.WithConfigFile("config.yaml", "yaml"),
    config.WithEnvPrefix("BASIC_EXAMPLE"),
    config.WithEnvVarOverride(true),
)
```

### 2. Error Handling Patterns
```go
// Create error with stack trace
err := errors.New("operation failed")

// Wrap with context
err = errors.Wrap(err, "failed to process data")

// Add error code
err = errors.WithCode(err, errors.ErrConfigFileRead)

// Extract error information
if coder := errors.GetCoder(err); coder != nil {
    fmt.Printf("Code: %d, Type: %s", coder.Code(), coder.String())
}
```

### 3. Logging Best Practices
```go
// Initialize from config
logOpts := createLogOptions(cfg.Log)
log.Init(logOpts)

// Structured logging
logger.Infow("Operation completed",
    "duration", "100ms",
    "items_processed", 42,
    "success", true,
)

// Context-aware logging
ctx = log.ContextWithRequestID(ctx, "req-123")
logger.CtxInfof(ctx, "Processing request")
```

## Configuration Options

The example uses these configuration sections:

- **app**: Application-specific settings
- **server**: HTTP server configuration (from SDK)
- **log**: Logging configuration (from SDK)
- **database**: Database connection settings (from SDK)
- **tracing**: Distributed tracing settings (from SDK)
- **metrics**: Metrics collection settings (from SDK)

## Environment Variable Override

Any configuration value can be overridden using environment variables:

Format: `{PREFIX}_{SECTION}_{FIELD}`

Examples:
- `BASIC_EXAMPLE_APP_NAME` → `app.name`
- `BASIC_EXAMPLE_LOG_LEVEL` → `log.level`
- `BASIC_EXAMPLE_SERVER_PORT` → `server.port`

## Next Steps

After understanding this basic example, explore:

1. **Config Features** (`../config-features/`) - Advanced configuration patterns
2. **Error Handling** (`../error-handling/`) - Comprehensive error management
3. **Logging Features** (`../logging-features/`) - Advanced logging capabilities
4. **Integration Examples** (`../integration/`) - Real-world application patterns

## Common Issues

### Issue: Configuration file not found
**Solution**: Ensure `config.yaml` exists in the same directory as `main.go`

### Issue: Environment variables not working
**Solution**: Check the prefix format: `BASIC_EXAMPLE_SECTION_FIELD`

### Issue: Log output not appearing
**Solution**: Verify log level settings in configuration file

## Related Documentation

- [Configuration Module Guide](../../docs/usage-guides/config/)
- [Error Handling Guide](../../docs/usage-guides/errors/)
- [Logging Guide](../../docs/usage-guides/log/) 