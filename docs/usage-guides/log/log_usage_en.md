# Logging (`pkg/log`) Usage Guide

[Switch to Chinese (切换到中文)](./log_usage_zh.md)

This guide explains how to use the `pkg/log` module within the `lmcc-go-sdk` for robust and configurable logging.

## 1. Feature Introduction

The `pkg/log` module provides a flexible and powerful logging solution based on `zap`. Key features include:

-   **Multiple Log Levels:** Supports standard levels like Debug, Info, Warn, Error, Fatal.
-   **Structured Logging:** Output logs in JSON or human-readable text format.
-   **Configurable Outputs:** Direct logs to `stdout`, `stderr`, or one or more files. Configured via `outputPaths` and `errorOutputPaths`.
-   **Log Rotation:** Automatic log file rotation based on size, age, and number of backups.
-   **Hot Reload:** Dynamically update logger configuration (level, format, output, etc.) when the application's configuration changes, by integrating with `pkg/config`.
-   **Contextual Logging:** Automatically include fields from `context.Context` (like Trace ID, Request ID, and custom keys configured via `contextKeys`) in log messages.
-   **Caller Information:** Optionally include file and line number of the log call site (enabled via `disableCaller: false`).
-   **Errors and Stack Traces:**
    -   For Error level and above, `zap` attempts to append a stack trace by default (can be disabled for `zap`'s own stack trace via `disableStacktrace: true`).
    -   When integrated with `github.com/marmotedu/errors`, if a logged error is wrapped by `marmotedu/errors`, its detailed stack trace is included in the `errorVerbose` field of JSON logs.
-   **Colorized Output:** Enable colorized output for different log levels when the format is `text` (via `enableColor: true`).
-   **Development Mode:** `development: true` configures more developer-friendly logging formats and behaviors.
-   **Logger Naming:** A name can be specified for the logger instance via the `name` field.

## 2. Integration Guide

This section demonstrates how to integrate `pkg/log` with `pkg/config` for a typical application setup.
For a complete runnable example, please refer to the `examples/simple-config-app` directory.

### 2.1. Configuration (`config.yaml`)

First, define the log settings in your `config.yaml` (or other supported config file format). The `log` section in your configuration file should correspond to the fields in `sdkconfig.LogConfig`.

```yaml
# Example config.yaml snippet
server:
  port: 9091

log:
  level: "debug"       # e.g., debug, info, warn, error, fatal
  format: "json"      # \"json\" or \"text\"
  enableColor: true   # Enables color output when format is \"text\" and terminal supports it
  outputPaths:        # List of log output paths
    - "stdout"
    - "./logs/app.log"
  errorOutputPaths:   # List of output paths for internal errors and PANIC logs (defaults to stderr)
    - "stderr"
    - "./logs/app_error.log"
  # filename: \"./logs/app.log\" # Older single file output, can be ignored if using outputPaths or used for specific rotation config
  maxSize: 100        # Max size in MB before rotation
  maxBackups: 3       # Max old log files to keep
  maxAge: 7           # Max days to keep old log files
  compress: false     # Compress rotated files
  disableCaller: false # false means output caller information (file and line number)
  disableStacktrace: false # false means zap will attempt to attach stacktrace for Error level and above (not the marmotedu/errors stack)
  development: false  # true enables more developer-friendly log configurations (e.g., more readable stack traces)
  name: "example-app" # Name for the logger
  contextKeys:        # List of additional keys to extract from context.Context and include in logs
    - "customKey1"
    - "user_id"
```

**Note:**
-   `sdkconfig.LogConfig` (in `pkg/config/types.go`) is the structure that `pkg/config` uses to unmarshal the `log` section from your config file.
-   The example in `examples/simple-config-app/main.go` uses a helper function (`createLogOpts`) to map these fields to `sdklog.Options` which `sdklog.Init()` expects. Ensure your `createLogOpts` or similar translation logic handles all fields you wish to read from the config.

### 2.1.1. JSON Log Format Key Names

When `format` is set to `"json"`, the logger uses concise key names for core fields to optimize performance and log size. The default key names for these core fields in JSON output are:

-   **`L`**: Log Level (e.g., "DEBUG", "INFO")
-   **`T`**: Timestamp (e.g., "2023-10-27T10:00:00.123Z")
-   **`M`**: Message
-   **`N`**: Logger Name (if configured via `log.name`)
-   **`C`**: Caller (e.g., "module/file.go:123")
-   **`stacktrace`**: Stacktrace (for ERROR, PANIC, FATAL levels, or when `errorVerbose` is present)

Contextual fields like `trace_id`, `request_id`, and keys specified in `contextKeys` will retain their names as configured or defined. The `errorVerbose` field also retains its name when present.

### 2.2. Application Code (`main.go` - Key Parts)

Here are the key parts of how to initialize and use the logger in your application. Please refer to `examples/simple-config-app/main.go` for the complete code.

```go
package main

import (
	"context"
	"flag"
	"fmt"
	stdlog "log" // Standard log for initial setup errors
	"os"
	"os/signal"
	"syscall"
	"time"

	sdkconfig "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
	sdklog "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
	"github.com/spf13/viper"
	merrors "github.com/marmotedu/errors" // For demonstrating stack traces
)

// MyAppConfig defines your application's config structure
type MyAppConfig struct {
	sdkconfig.Config // Embed SDK base config
	// Add other custom app config fields here
}

var AppCfg MyAppConfig

// createLogOpts converts sdkconfig.LogConfig to sdklog.Options
func createLogOpts(cfg *sdkconfig.LogConfig) *sdklog.Options {
	if cfg == nil {
		sdklog.Warn("Log configuration section is nil, creating default log options.")
		return sdklog.NewOptions() 
	}
	opts := sdklog.NewOptions() // Start with defaults

	opts.Level = cfg.Level
	opts.Format = cfg.Format
	opts.EnableColor = cfg.EnableColor // New: Pass color config

	// Output paths
	if len(cfg.OutputPaths) > 0 {
		opts.OutputPaths = cfg.OutputPaths
	} else if cfg.Output != "" { // Backward compatibility for old 'output' field
	if cfg.Output == "stdout" {
		opts.OutputPaths = []string{"stdout"}
	} else if cfg.Output == "stderr" {
		opts.OutputPaths = []string{"stderr"}
		} else {
			// If Filename also exists, it might take precedence over Output in older logic
			// Simplified here: if OutputPaths is empty, Output (if not stdout/stderr) acts as single file path
			filePath := cfg.Output
			if cfg.Filename != "" { // If Filename exists, prioritize it
				filePath = cfg.Filename
			}
			opts.OutputPaths = []string{filePath}
		}
	} else {
		opts.OutputPaths = []string{"stdout"} // Default
	}

	// Error output paths
	if len(cfg.ErrorOutputPaths) > 0 {
		opts.ErrorOutputPaths = cfg.ErrorOutputPaths
	} else if cfg.ErrorOutput != "" { // Backward compatibility for old 'errorOutput' field
	    if cfg.ErrorOutput == "stdout" {
			opts.ErrorOutputPaths = []string{"stdout"}
		} else if cfg.ErrorOutput == "stderr" {
			opts.ErrorOutputPaths = []string{"stderr"}
		} else {
			opts.ErrorOutputPaths = []string{cfg.ErrorOutput}
		}
	} else {
	    opts.ErrorOutputPaths = []string{"stderr"} // Default
	}

	opts.LogRotateMaxSize = cfg.MaxSize
	opts.LogRotateMaxBackups = cfg.MaxBackups
	opts.LogRotateMaxAge = cfg.MaxAge
	opts.LogRotateCompress = cfg.Compress

	opts.DisableCaller = cfg.DisableCaller         // New
	opts.DisableStacktrace = cfg.DisableStacktrace // New
	opts.Development = cfg.Development           // New
	opts.Name = cfg.Name                           // New

	if len(cfg.ContextKeys) > 0 { // New
		opts.ContextKeys = make([]any, len(cfg.ContextKeys))
		for i, k := range cfg.ContextKeys {
			opts.ContextKeys[i] = k
		}
	}
	return opts
}

// deeperErrorFunction simulates a function generating an error
func deeperErrorFunction() error {
    return merrors.Wrap(merrors.New("underlying database error"), "service layer processing failed")
}

func main() {
	configFile := flag.String("config", "config.yaml", "Path to config file")
	flag.Parse()

	configManager, err := sdkconfig.LoadConfigAndWatch(
		&AppCfg,
		sdkconfig.WithConfigFile(*configFile, "yaml"),
		sdkconfig.WithHotReload(true),
	)
	if err != nil {
		stdlog.Fatalf("FATAL: Failed to load initial configuration: %v\\n", err)
	}
	stdlog.Println("Initial configuration loaded successfully.")

	if AppCfg.Log == nil {
		stdlog.Fatalln("FATAL: Log configuration section is missing.")
	}

	logOpts := createLogOpts(AppCfg.Log)
	sdklog.Init(logOpts)
	sdklog.Info("SDK Logger initialized with initial config.")
	sdklog.Infof("Initial log settings: Level=%s, Format=%s, OutputPaths=%v, EnableColor=%t",
		logOpts.Level, logOpts.Format, logOpts.OutputPaths, logOpts.EnableColor)

	if configManager != nil {
		configManager.RegisterCallback(func(v *viper.Viper, currentCfgAny any) error {
			currentTypedCfg, ok := currentCfgAny.(*MyAppConfig)
			if !ok { /* ... error handling ... */ return fmt.Errorf("config type error")}
			if currentTypedCfg.Log == nil { /* ... error handling ... */ return fmt.Errorf("log config missing")}
			
			sdklog.Info("Configuration reloaded. Re-initializing logger...")
			newLogOpts := createLogOpts(currentTypedCfg.Log)
			sdklog.Init(newLogOpts) // Re-initialize
			sdklog.Infof("SDK Logger re-initialized. New settings: Level=%s, Format=%s, OutputPaths=%v, EnableColor=%t",
				newLogOpts.Level, newLogOpts.Format, newLogOpts.OutputPaths, newLogOpts.EnableColor)
			// Demo color output after re-init
			if newLogOpts.Format == sdklog.FormatText && newLogOpts.EnableColor {
				sdklog.Info("\\033[32mThis INFO message should be green.\\033[0m")
				sdklog.Warn("\\033[33mThis WARN message should be yellow.\\033[0m")
			}
			return nil
		})
		sdklog.Info("Callback for logger updates registered.")
	}

	// --- Demonstrate Logging --- 
	sdklog.Debug("This is a debug message.")
	sdklog.Infow("User logged in", "username", "martin", "sessionID", 12345)
	
	// Demonstrate error and stack trace
	errWithStack := deeperErrorFunction()
	sdklog.Errorw("An error occurred, with stack trace from marmotedu/errors", "error", errWithStack, "relevant_id", "id-123")

	// Contextual logging
	ctx := context.Background()
	ctx = sdklog.ContextWithTraceID(ctx, "trace-abc-123")
	ctx = sdklog.ContextWithRequestID(ctx, "req-def-456")
	// Assume "customKey1" is configured in ContextKeys
	ctx = context.WithValue(ctx, "customKey1", "customValueForLog")
	sdklog.Ctx(ctx, "Processing request with trace, request ID, and customKey1.")

	sdklog.Info("Application running. Modify config.yaml to test hot reload. Press Ctrl+C to exit.")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	sdklog.Info("Application shutting down.")
	if err := sdklog.Sync(); err != nil {
		stdlog.Printf("Error syncing logger: %v\\n", err)
	}
}
```

### 2.3. Running the Example

1.  Refer to the complete example code and `config.yaml` in the `examples/simple-config-app` directory.
2.  Run `go run examples/simple-config-app/main.go -config examples/simple-config-app/config.yaml`.
3.  While it's running, modify `log.level`, `log.format`, etc., in `config.yaml` and save. Observe the log output reflecting the changes.

## 3. API Reference

### 3.1. Core Functions

-   **`Init(opts *Options)`**: Initializes or re-initializes the global logger with the given options. This function is thread-safe.
-   **`NewOptions() *Options`**: Returns a new `Options` struct populated with default values.
-   **`NewLogger(opts *Options) Logger`**: Creates a new logger instance with the given options. Useful if you need multiple distinct logger instances (though typically the global logger `Std()` is sufficient).
-   **`Std() Logger`**: Returns the global singleton logger instance. It is initialized with default options if `Init` has not been called.
-   **`Sync() error`**: Flushes any buffered log entries from the global logger. It's important to call this before application exit.

### 3.2. Logging Methods (Global and on `Logger` instance)

These are available as global functions (e.g., `sdklog.Info(...)`) which use the global logger, and as methods on a `Logger` instance (e.g., `myLogger.Info(...)`).

-   `Debug(args ...any)` / `Debugf(template string, args ...any)` / `Debugw(msg string, keysAndValues ...any)`
-   `Info(args ...any)` / `Infof(template string, args ...any)` / `Infow(msg string, keysAndValues ...any)`
-   `Warn(args ...any)` / `Warnf(template string, args ...any)` / `Warnw(msg string, keysAndValues ...any)`
-   `Error(args ...any)` / `Errorf(template string, args ...any)` / `Errorw(msg string, keysAndValues ...any)`
-   `Fatal(args ...any)` / `Fatalf(template string, args ...any)` / `Fatalw(msg string, keysAndValues ...any)` (Note: Fatal logs then calls `os.Exit(1)`)

### 3.3. Contextual Logging

The `pkg/log` module allows you to enrich your log messages with data from `context.Context`. This is particularly useful for request tracing and associating logs with specific operations.

-   **`ContextWithTraceID(ctx context.Context, traceID string) context.Context`**: Returns a new context with the TraceID set. The `pkg/log` module internally uses an unexported key type (`pkg/log.TraceIDKey`) to store this value.
-   **`ContextWithRequestID(ctx context.Context, requestID string) context.Context`**: Returns a new context with the RequestID set. Similar to TraceID, an internal key (`pkg/log.RequestIDKey`) is used.
-   **`TraceIDFromContext(ctx context.Context) (string, bool)`**: Extracts TraceID from the context, if present.
-   **`RequestIDFromContext(ctx context.Context) (string, bool)`**: Extracts RequestID from the context, if present.

-   **`Ctx(ctx context.Context, args ...any)`**: Logs a message at InfoLevel. It automatically extracts recognized fields from the context:
    -   **Trace ID**: Extracted using `TraceIDFromContext` and logged with the field key **`"trace_id"`**.
    -   **Request ID**: Extracted using `RequestIDFromContext` and logged with the field key **`"request_id"`**.
    -   **Custom Keys**: If a list of keys is configured via `Options.ContextKeys` (corresponding to `log.contextKeys` in the config file), their corresponding values will also be extracted from the context and logged.

    **Example:**
    ```go
    package main

    import (
    	"context"
    	"fmt"
    	sdklog "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
    	"github.com/google/uuid" // For generating unique IDs
    )

    // Define a custom key type (as an example for advanced context usage)
    type myCustomKey struct{}

    func main() {
    	// Initialize logger (assuming opts are set, e.g., to output JSON to stdout)
    	opts := sdklog.NewOptions()
    	opts.Format = "json"
    	opts.OutputPaths = []string{"stdout"}
    	opts.Level = "info"
    	// To extract 'myCustomKey', you would add it to opts.ContextKeys:
    	// opts.ContextKeys = []any{myCustomKey{}} 
    	sdklog.Init(opts)

    	traceID := uuid.NewString()
    	requestID := uuid.NewString()
    	customValue := "my_custom_context_data"

    	ctx := context.Background()
    	ctx = sdklog.ContextWithTraceID(ctx, traceID)
    	ctx = sdklog.ContextWithRequestID(ctx, requestID)
    	ctx = context.WithValue(ctx, myCustomKey{}, customValue) // Store custom value

    	// Log with context. `pkg/log` will automatically pick up trace_id and request_id.
    	// If myCustomKey{} was added to Options.ContextKeys, it would also be picked up.
    	// The field name for myCustomKey{} in the log would be its string representation, e.g., "{}"
    	sdklog.Ctx(ctx, "Processing user request with contextual data.")
    	
    	// Example of how you might verify this (conceptual, for a test):
    	// Log output (JSON):
    	// {
    	//   "level": "info",
    	//   "timestamp": "...",
    	//   "caller": "...",
    	//   "message": "Processing user request with contextual data.",
    	//   "trace_id": "generated-trace-id",  // Automatically added
    	//   "request_id": "generated-request-id", // Automatically added
    	//   "{}": "my_custom_context_data"  // If myCustomKey{} was in ContextKeys, and its string form is "{}"
    	// }
    }
    ```

-   Similar context-aware logging functions are available for other levels:
    -   `CtxDebugf(ctx context.Context, template string, args ...interface{})`
    -   `CtxInfof(ctx context.Context, template string, args ...interface{})`
    -   `CtxWarnf(ctx context.Context, template string, args ...interface{})`
    -   `CtxErrorf(ctx context.Context, template string, args ...interface{})`
    -   `CtxFatalf(ctx context.Context, template string, args ...interface{})` (also exits)
    -   `CtxPanicf(ctx context.Context, template string, args ...interface{})` (also panics)

### 3.4. Logger Manipulation

-   **`WithName(name string) Logger`**: Returns a new logger instance with the specified name appended to its existing name.
-   **`WithValues(keysAndValues ...any) Logger`**: Returns a new logger instance with the given key-value pairs added to its structured context.

### 3.5. `Options` Struct

Key fields in `pkg/log/Options` (refer to `pkg/log/options.go` for all fields and defaults):

-   `Level`: `string` (e.g., "debug", "info")
-   `Format`: `string` ("json" or "text")
-   `OutputPaths`: `[]string` (e.g., `["stdout"]`, `["./logs/app.log"]`)
-   `ErrorOutputPaths`: `[]string` (for internal logger errors)
-   `DisableCaller`: `bool`
-   `DisableStacktrace`: `bool`
-   `EnableColor`: `bool` (for text format)
-   `Development`: `bool`
-   `Name`: `string` (logger name)
-   `LogRotateMaxSize`: `int` (MB)
-   `LogRotateMaxBackups`: `int`
-   `LogRotateMaxAge`: `int` (days)
-   `LogRotateCompress`: `bool`
-   `ContextKeys`: `[]any` (list of custom keys to extract from context)

## 4. Advanced Features

### 4.1. Error Stack Traces

When you log an `error` type object:

-   **`zap` Default Behavior**: If the log level is Error or higher, and `log.disableStacktrace` is `false` (default), `zap` attempts to append a stack trace. This stack trace is typically logged in a field named `stacktrace` (JSON format).
-   **Integration with `marmotedu/errors`**: If the error you are logging was wrapped using `github.com/marmotedu/errors`, the more detailed stack trace information captured by `marmotedu/errors` will be logged under the `errorVerbose` field name in JSON logs. This is often more readable than `zap`'s own stack trace as it focuses on the path of error generation.
    -   To leverage this, ensure your errors are created with functions like `merrors.New`, `merrors.Errorf`, `merrors.Wrap`, etc.
    -   Log using `sdklog.Errorw("message", "error", yourMarmotError)` or `sdklog.Errorf("message: %+v", yourMarmotError)` (when format is text, `%+v` prints the stack).
    -   The `log.disableStacktrace: true` setting **does not** disable the `marmotedu/errors` stack in `errorVerbose`.

**Example (`log.format: "json"` in `config.yaml`)**:

```go
// In your code:
import merrors "github.com/marmotedu/errors"
// ...
func doSomething() error {
    return merrors.New("something bad happened")
}

err := doSomething()
if err != nil {
    sdklog.Errorw("Operation failed", "error", err, "user_id", 123)
}
```

**Possible JSON Output Snippet**:
```json
{
  "L": "ERROR", "T": "...", "C": "...", "N": "example-app",
  "M": "Operation failed",
  "user_id": 123,
  "error": "something bad happened",
  "errorVerbose": "something bad happened\\n    main.doSomething\\n        /path/to/your/app/main.go:XX\\n    main.main\\n        /path/to/your/app/main.go:YY\\n    ...",
  "stacktrace": "main.main\\n\\t/path/to/your/app/main.go:YY\\nruntime.main..." // zap's stack (if disableStacktrace=false)
}
```
As shown, `errorVerbose` provides the stack trace generated by `marmotedu/errors`.

### 4.2. Colorized Log Output

For better readability during development or local debugging, you can enable colorized log output.

-   **Configuration**:
    -   Set `format: "text"` in the `log` section of your `config.yaml`.
    -   Set `enableColor: true`.
-   **Effect**: When logs are output to a terminal that supports ANSI color escape sequences, different log levels will be displayed in different colors (e.g., Error in red, Warn in yellow, Info in green, etc.).
-   **Note**: JSON formatted logs are inherently structured data and do not contain color information. This feature is primarily for text-formatted console output.

**Example (`config.yaml`)**:
```yaml
log:
  format: "text"
  enableColor: true
  level: "debug"
  outputPaths: ["stdout"]
```
When you run your application and view the terminal output, you will see colorized log entries.
`examples/simple-config-app` also demonstrates how to check color settings after a config hot reload and print specifically colored messages.

## 5. Log Rotation

The `pkg/log` module supports log rotation via `lumberjack.v2`. Relevant configuration options:
-   `maxSize`: Maximum size in megabytes of a single log file.
-   `maxBackups`: Maximum number of old log files to retain.
-   `maxAge`: Maximum number of days to retain old log files.
-   `compress`: Whether to compress rotated log files (e.g., using gzip).

These settings take effect when logging to a file (e.g., `outputPaths: ["./logs/app.log"]`).

## 6. Hot Reload

Through integration with the `pkg/config` module (`sdkconfig.LoadConfigAndWatch` and `configManager.RegisterCallback`), `pkg/log` can dynamically respond to changes in the configuration file at runtime. This means you can modify the `log` section in `config.yaml` (e.g., change `level`, `format`, `enableColor`, or output paths), and the application's logging behavior will update accordingly without requiring a restart.

A full demonstration can be found in `examples/simple-config-app/main.go`.

## 7. Best Practices

-   **Sync Logs**: Always call `sdklog.Sync()` before your application exits to ensure all buffered logs are written.
-   **Prefer Structured Logs**: For production environments, prefer the `json` format as it's easier for machines to parse and integrate into log management systems.
-   **Appropriate Log Levels**: Choose log levels wisely based on the importance and frequency of the information. Avoid excessive Debug level logging in production.
-   **Use Contextual Logging**: For request-related logs, always pass the `context.Context` to include tracing information.
-   **Error Handling**: Log errors using `Errorw` or `Errorf`, providing as much context as possible. If using `marmotedu/errors`, its stack trace will be automatically captured and logged.
-   **Configuration Management**: Manage log settings via configuration files and leverage the hot reload feature for dynamic adjustments.

This guide provides a comprehensive overview of the `pkg/log` module. For more details and advanced usage, please refer to the source code and the official `zap` documentation.
