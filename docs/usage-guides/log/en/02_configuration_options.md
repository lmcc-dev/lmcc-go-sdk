# Configuration Options

This document provides detailed descriptions of all configuration options available in the log module and how to use them effectively.

## Options Structure

The log module uses the `log.Options` structure to configure all logging behavior:

```go
type Options struct {
    // Basic configuration
    Level            string   `mapstructure:"level" default:"info"`
    Format           string   `mapstructure:"format" default:"text"`
    OutputPaths      []string `mapstructure:"output_paths" default:"[\"stdout\"]"`
    ErrorOutputPaths []string `mapstructure:"error_output_paths" default:"[\"stderr\"]"`
    
    // Display options
    EnableColor      bool `mapstructure:"enable_color" default:"true"`
    DisableCaller    bool `mapstructure:"disable_caller" default:"false"`
    DisableStacktrace bool `mapstructure:"disable_stacktrace" default:"false"`
    StacktraceLevel  string `mapstructure:"stacktrace_level" default:"error"`
    
    // Log rotation configuration
    LogRotateMaxSize    int  `mapstructure:"log_rotate_max_size" default:"100"`
    LogRotateMaxBackups int  `mapstructure:"log_rotate_max_backups" default:"5"`
    LogRotateMaxAge     int  `mapstructure:"log_rotate_max_age" default:"30"`
    LogRotateCompress   bool `mapstructure:"log_rotate_compress" default:"true"`
}
```

## Basic Configuration Options

### Level (Log Level)

Controls which log messages will be output.

**Available values:**
- `"debug"` - Show all log messages
- `"info"` - Show info, warn, error, fatal, panic
- `"warn"` - Show warn, error, fatal, panic
- `"error"` - Show error, fatal, panic
- `"fatal"` - Show fatal, panic
- `"panic"` - Show only panic

**Examples:**
```go
opts := &log.Options{
    Level: "debug",  // Development environment
}

opts := &log.Options{
    Level: "warn",   // Production environment
}
```

### Format (Output Format)

Controls the output format of log messages.

**Available values:**
- `"text"` - Human-readable text format (default)
- `"json"` - Structured JSON format
- `"keyvalue"` - Key-value pair format

**Examples:**

```go
// Text format (suitable for development)
opts := &log.Options{
    Format: "text",
}
// Output: 2024-01-15T10:30:45.123Z	INFO	main.go:25	User login	{"user_id": 123}

// JSON format (suitable for production and log aggregation)
opts := &log.Options{
    Format: "json",
}
// Output: {"level":"info","timestamp":"2024-01-15T10:30:45.123Z","caller":"main.go:25","message":"User login","user_id":123}

// Key-value format (suitable for traditional log systems)
opts := &log.Options{
    Format: "keyvalue",
}
// Output: timestamp=2024-01-15T10:30:45.123Z level=info caller=main.go:25 message="User login" user_id=123
```

### OutputPaths (Output Paths)

Specifies the output destinations for log messages.

**Available values:**
- `"stdout"` - Standard output
- `"stderr"` - Standard error
- File paths - e.g., `"/var/log/app.log"`

**Examples:**
```go
// Output to console only
opts := &log.Options{
    OutputPaths: []string{"stdout"},
}

// Output to both console and file
opts := &log.Options{
    OutputPaths: []string{"stdout", "/var/log/app.log"},
}

// Output to file only
opts := &log.Options{
    OutputPaths: []string{"/var/log/app.log"},
}
```

### ErrorOutputPaths (Error Output Paths)

Specifies the output destinations for error-level logs.

**Example:**
```go
opts := &log.Options{
    OutputPaths:      []string{"stdout", "/var/log/app.log"},
    ErrorOutputPaths: []string{"stderr", "/var/log/error.log"},
}
```

## Display Options

### EnableColor (Enable Color)

Enables color coding in terminal output.

**Example:**
```go
// Development environment: enable color
opts := &log.Options{
    Format:      "text",
    EnableColor: true,
}

// Production environment or file output: disable color
opts := &log.Options{
    Format:      "json",
    EnableColor: false,
}
```

### DisableCaller (Disable Caller Information)

Controls whether to include caller information (filename and line number) in logs.

**Example:**
```go
// Include caller information (default)
opts := &log.Options{
    DisableCaller: false,
}
// Output: 2024-01-15T10:30:45.123Z	INFO	main.go:25	Message

// Exclude caller information
opts := &log.Options{
    DisableCaller: true,
}
// Output: 2024-01-15T10:30:45.123Z	INFO	Message
```

### DisableStacktrace (Disable Stack Trace)

Controls whether to include stack traces in error logs.

**Example:**
```go
// Enable stack trace (default)
opts := &log.Options{
    DisableStacktrace: false,
    StacktraceLevel:   "error",
}

// Disable stack trace
opts := &log.Options{
    DisableStacktrace: true,
}
```

### StacktraceLevel (Stack Trace Level)

Specifies from which level to include stack traces.

**Available values:**
- `"debug"`, `"info"`, `"warn"`, `"error"`, `"fatal"`, `"panic"`

**Example:**
```go
// Show stack trace only for error and above
opts := &log.Options{
    StacktraceLevel: "error",
}

// Show stack trace for warn and above
opts := &log.Options{
    StacktraceLevel: "warn",
}
```

## Log Rotation Configuration

### LogRotateMaxSize (Maximum File Size)

Maximum size of a single log file (MB).

**Example:**
```go
opts := &log.Options{
    OutputPaths:         []string{"/var/log/app.log"},
    LogRotateMaxSize:    100,  // 100 MB
}
```

### LogRotateMaxBackups (Maximum Backup Count)

Number of old log files to retain.

**Example:**
```go
opts := &log.Options{
    LogRotateMaxBackups: 5,  // Keep 5 backup files
}
```

### LogRotateMaxAge (Maximum Retention Days)

Maximum number of days to retain log files.

**Example:**
```go
opts := &log.Options{
    LogRotateMaxAge: 30,  // Keep for 30 days
}
```

### LogRotateCompress (Compress Old Files)

Whether to compress rotated log files.

**Example:**
```go
opts := &log.Options{
    LogRotateCompress: true,  // Compress old files to save space
}
```

## Configuration Examples

### Development Environment Configuration

```go
func developmentConfig() *log.Options {
    return &log.Options{
        Level:             "debug",
        Format:            "text",
        OutputPaths:       []string{"stdout"},
        EnableColor:       true,
        DisableCaller:     false,
        DisableStacktrace: false,
        StacktraceLevel:   "error",
    }
}
```

### Production Environment Configuration

```go
func productionConfig() *log.Options {
    return &log.Options{
        Level:               "info",
        Format:              "json",
        OutputPaths:         []string{"/var/log/app.log"},
        ErrorOutputPaths:    []string{"/var/log/error.log"},
        EnableColor:         false,
        DisableCaller:       false,
        DisableStacktrace:   false,
        StacktraceLevel:     "error",
        LogRotateMaxSize:    100,
        LogRotateMaxBackups: 10,
        LogRotateMaxAge:     30,
        LogRotateCompress:   true,
    }
}
```

### Test Environment Configuration

```go
func testConfig() *log.Options {
    return &log.Options{
        Level:             "warn",
        Format:            "text",
        OutputPaths:       []string{"stdout"},
        EnableColor:       false,
        DisableCaller:     true,
        DisableStacktrace: true,
    }
}
```

### High Performance Configuration

```go
func highPerformanceConfig() *log.Options {
    return &log.Options{
        Level:             "error",  // Log errors only
        Format:            "json",   // Faster serialization
        OutputPaths:       []string{"/var/log/app.log"},
        EnableColor:       false,    // Disable color processing
        DisableCaller:     true,     // Disable caller lookup
        DisableStacktrace: true,     // Disable stack trace
    }
}
```

## Loading from Configuration Files

### YAML Configuration

```yaml
# config.yaml
log:
  level: "info"
  format: "json"
  output_paths: ["stdout", "/var/log/app.log"]
  error_output_paths: ["stderr", "/var/log/error.log"]
  enable_color: false
  disable_caller: false
  disable_stacktrace: false
  stacktrace_level: "error"
  log_rotate_max_size: 100
  log_rotate_max_backups: 5
  log_rotate_max_age: 30
  log_rotate_compress: true
```

```go
import (
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

type AppConfig struct {
    Log log.Options `mapstructure:"log"`
}

func main() {
    var cfg AppConfig
    
    err := config.LoadConfig(&cfg,
        config.WithConfigFile("config.yaml", ""),
    )
    if err != nil {
        panic(err)
    }
    
    log.Init(&cfg.Log)
}
```

### JSON Configuration

```json
{
  "log": {
    "level": "info",
    "format": "json",
    "output_paths": ["stdout", "/var/log/app.log"],
    "enable_color": false,
    "log_rotate_max_size": 100
  }
}
```

### Environment Variable Override

```bash
# Environment variables can override configuration file settings
export APP_LOG_LEVEL=debug
export APP_LOG_FORMAT=text
export APP_LOG_OUTPUT_PATHS=stdout
export APP_LOG_ENABLE_COLOR=true
```

## Dynamic Configuration Updates

### Using Configuration Hot Reload

```go
type AppConfig struct {
    Log log.Options `mapstructure:"log"`
}

func main() {
    var cfg AppConfig
    
    cm, err := config.LoadConfigAndWatch(&cfg,
        config.WithConfigFile("config.yaml", ""),
        config.WithHotReload(true),
    )
    if err != nil {
        panic(err)
    }
    
    // Initialize logging
    log.Init(&cfg.Log)
    
    // Register log configuration change callback
    cm.RegisterSectionChangeCallback("log", func(v *viper.Viper) error {
        var newLogOpts log.Options
        if err := v.UnmarshalKey("log", &newLogOpts); err != nil {
            return err
        }
        
        log.Info("Updating log configuration")
        log.Init(&newLogOpts)
        log.Info("Log configuration updated")
        
        return nil
    })
    
    // Application logic...
}
```

## Configuration Validation

### Validation Function

```go
func (opts *log.Options) Validate() error {
    // Validate log level
    validLevels := []string{"debug", "info", "warn", "error", "fatal", "panic"}
    if !contains(validLevels, opts.Level) {
        return fmt.Errorf("invalid log level: %s", opts.Level)
    }
    
    // Validate output format
    validFormats := []string{"text", "json", "keyvalue"}
    if !contains(validFormats, opts.Format) {
        return fmt.Errorf("invalid log format: %s", opts.Format)
    }
    
    // Validate output paths
    if len(opts.OutputPaths) == 0 {
        return fmt.Errorf("at least one output path is required")
    }
    
    // Validate rotation configuration
    if opts.LogRotateMaxSize <= 0 {
        return fmt.Errorf("log rotate max size must be greater than 0")
    }
    
    if opts.LogRotateMaxBackups < 0 {
        return fmt.Errorf("log rotate max backups cannot be negative")
    }
    
    return nil
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}
```

### Using Validation

```go
func initializeLogging(opts *log.Options) error {
    // Validate configuration
    if err := opts.Validate(); err != nil {
        return fmt.Errorf("log configuration validation failed: %w", err)
    }
    
    // Initialize logging
    log.Init(opts)
    
    return nil
}
```

## Best Practices

1. **Environment-specific configuration**: Use different configurations for different environments
2. **Sensitive information**: Avoid logging sensitive information
3. **Performance considerations**: Adjust log level and options in high-load environments
4. **Rotation configuration**: Always configure log rotation in production
5. **Monitoring**: Monitor log file size and disk usage

## Next Steps

- [Output Formats](03_output_formats.md) - Learn about different output formats
- [Context Logging](04_context_logging.md) - Master context-aware logging
- [Performance](05_performance.md) - Optimize logging performance
- [Best Practices](06_best_practices.md) - Production-ready patterns 