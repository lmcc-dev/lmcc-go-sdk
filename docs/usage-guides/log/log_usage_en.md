# Logging (`pkg/log`) Usage Guide

[Switch to Chinese (切换到中文)](./log_usage_zh.md)

This guide explains how to use the `pkg/log` module within the `lmcc-go-sdk` for robust and configurable logging.

## 1. Feature Introduction

The `pkg/log` module provides a flexible and powerful logging solution based on `zap`. Key features include:

-   **Multiple Log Levels:** Supports standard levels like Debug, Info, Warn, Error, Fatal.
-   **Structured Logging:** Output logs in JSON or human-readable text format.
-   **Configurable Outputs:** Direct logs to `stdout`, `stderr`, or one or more files.
-   **Log Rotation:** Automatic log file rotation based on size, age, and number of backups.
-   **Hot Reload:** Dynamically update logger configuration (level, format, output) when the application's configuration changes, by integrating with `pkg/config`.
-   **Contextual Logging:** Automatically include fields from `context.Context` (like Trace ID, Request ID) in log messages.
-   **Caller Information:** Optionally include file and line number of the log call site.
-   **Stack Traces:** Automatically include stack traces for error-level logs.

## 2. Integration Guide

This section demonstrates how to integrate `pkg/log` with `pkg/config` for a typical application setup.

### 2.1. Configuration (`config.yaml`)

First, define the log settings in your `config.yaml` (or other supported config file format). The `log` section in your configuration file should correspond to the fields in `sdkconfig.LogConfig` (which itself maps to fields in `sdklog.Options`).

```yaml
# Example config.yaml snippet
server:
  port: 9091

log:
  level: "info"       # e.g., debug, info, warn, error, fatal
  format: "json"      # "json" or "text"
  output: "stdout"    # "stdout", "stderr", or a file path like "./logs/app.log"
  # filename: "./logs/app.log" # Alternative to 'output' if you only specify a file and want rotation
  maxSize: 100        # Max size in MB before rotation
  maxBackups: 3       # Max old log files to keep
  maxAge: 7           # Max days to keep old log files
  compress: false     # Compress rotated files
  # disableCaller: false
  # disableStacktrace: false 
  # enableColor: true # if format is text
  # development: false
  # name: "my-app-logger"
  # errorOutputPaths: ["stderr", "./logs/app_error.log"]
  # contextKeys: ["customKey1", "customKey2"] # If you have custom keys to extract from context
```

**Note:** `sdkconfig.LogConfig` (in `pkg/config/types.go`) is the structure that `pkg/config` uses to unmarshal the `log` section from your config file. The `examples/simple-config-app/main.go` then uses a helper function (`createLogOpts`) to map these fields to `sdklog.Options` which `sdklog.Init()` expects.

### 2.2. Application Code (`main.go`)

Here's how to initialize and use the logger in your application:

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
)

// Define your application's config structure
// (Define your application's config structure)
type MyAppConfig struct {
	sdkconfig.Config // Embed SDK base config
	// Add other custom app config fields here
}

var AppCfg MyAppConfig

// createLogOpts converts sdkconfig.LogConfig to sdklog.Options
// (createLogOpts converts sdkconfig.LogConfig to sdklog.Options)
func createLogOpts(cfg *sdkconfig.LogConfig) *sdklog.Options {
	if cfg == nil {
		return sdklog.NewOptions() 
	}
	opts := sdklog.NewOptions() 
	opts.Level = cfg.Level
	opts.Format = cfg.Format
	if cfg.Output == "stdout" {
		opts.OutputPaths = []string{"stdout"}
	} else if cfg.Output == "stderr" {
		opts.OutputPaths = []string{"stderr"}
	} else if cfg.Output != "" {
		if cfg.Filename != "" {
			opts.OutputPaths = []string{cfg.Filename}
		} else {
			opts.OutputPaths = []string{cfg.Output}
		}
	}
	opts.LogRotateMaxSize = cfg.MaxSize
	opts.LogRotateMaxBackups = cfg.MaxBackups
	opts.LogRotateMaxAge = cfg.MaxAge
	opts.LogRotateCompress = cfg.Compress
	// To use other fields from sdklog.Options, ensure they are present in
	// sdkconfig.LogConfig and map them here. For example:
	// opts.DisableCaller = cfg.DisableCaller // Assuming DisableCaller exists in sdkconfig.LogConfig
	// opts.Name = cfg.Name                 // Assuming Name exists in sdkconfig.LogConfig
	return opts
}

func main() {
	configFile := flag.String("config", "config.yaml", "Path to config file")
	flag.Parse()

	// Load configuration and watch for changes
	// (Load configuration and watch for changes)
	configManager, err := sdkconfig.LoadConfigAndWatch(
		&AppCfg,
		sdkconfig.WithConfigFile(*configFile, "yaml"),
		sdkconfig.WithHotReload(true),
	)
	if err != nil {
		stdlog.Fatalf("FATAL: Failed to load initial configuration: %v\n", err)
	}
	stdlog.Println("Initial configuration loaded successfully.")

	if AppCfg.Log == nil {
		stdlog.Fatalln("FATAL: Log configuration section is missing.")
	}

	// Initialize logger with loaded config
	// (Initialize logger with loaded config)
	logOpts := createLogOpts(AppCfg.Log)
	sdklog.Init(logOpts) // sdklog.Init does not return an error
	sdklog.Info("SDK Logger initialized with initial config.")

	// Register callback for log config changes
	// (Register callback for log config changes)
	if configManager != nil {
		configManager.RegisterCallback(func(v *viper.Viper, currentCfgAny any) error {
			currentTypedCfg, ok := currentCfgAny.(*MyAppConfig)
			if !ok {
				sdklog.Error("Config type assertion failed in callback")
				return fmt.Errorf("config type assertion error")
			}
			if currentTypedCfg.Log == nil {
				sdklog.Warn("Log configuration section missing after reload.")
				return fmt.Errorf("log config missing after reload")
			}
			sdklog.Info("Configuration reloaded. Re-initializing logger...")
			newLogOpts := createLogOpts(currentTypedCfg.Log)
			sdklog.Init(newLogOpts)
			sdklog.Infof("SDK Logger re-initialized. New level: %s, New format: %s", 
				newLogOpts.Level, newLogOpts.Format)
			return nil
		})
		sdklog.Info("Callback for logger updates registered.")
	}

	// --- Demonstrate Logging --- 
	sdklog.Debug("This is a debug message.")
	sdklog.Infow("User logged in", "username", "martin", "sessionID", 12345)
	sdklog.Warn("Potential issue detected.")
	sdklog.Error("An error occurred", "errorDetails", fmt.Errorf("database connection failed"))

	// Contextual logging
	// (Contextual logging)
	ctx := context.Background()
	ctx = sdklog.ContextWithTraceID(ctx, "trace-abc-123")
	ctx = sdklog.ContextWithRequestID(ctx, "req-def-456")
	sdklog.Ctx(ctx, "Processing request with trace and request ID.")

	sdklog.Info("Application running. Modify config.yaml to test hot reload. Press Ctrl+C to exit.")

	// Keep app running
	// (Keep app running)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	sdklog.Info("Application shutting down.")
	if err := sdklog.Sync(); err != nil {
		stdlog.Printf("Error syncing logger: %v\n", err)
	}
}
```

### 2.3. Running the Example

1.  Save the Go code above as `main.go` in a new directory.
2.  Create a `config.yaml` in the same directory with the log settings.
3.  Ensure `lmcc-go-sdk` is in your `go.mod`.
4.  Run `go mod tidy`.
5.  Run `go run main.go -config config.yaml`.
6.  While it's running, modify `log.level` or `log.format` in `config.yaml` and save. Observe the log output reflecting the changes.

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
    -   **Custom Keys**: If `Options.ContextKeys` is populated during logger initialization, those keys will also be extracted. For custom context keys that are not simple strings (e.g., struct types), `pkg/log` typically uses `fmt.Sprintf("%v", key)` to generate the field name in the log output.

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

## 4. Relevant Makefile Commands

-   `make test-unit PKG=./pkg/log`: Runs unit tests for `pkg/log`.
-   `make cover PKG=./pkg/log`: Runs unit tests with coverage for `pkg/log`.
-   `make test-integration`: Runs all integration tests (relevant if `pkg/log` has specific integration tests, or is tested as part of broader integration tests).
-   `make lint`: Lints the codebase, including `pkg/log`.

This guide provides a comprehensive overview of using the `pkg/log` module. For more detailed information, refer to the source code and specific function documentation within the `pkg/log` directory.
