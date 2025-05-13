# LMCC Go SDK

[![Go Report Card](https://goreportcard.com/badge/github.com/lmcc-dev/lmcc-go-sdk)](https://goreportcard.com/report/github.com/lmcc-dev/lmcc-go-sdk)
[![Go Reference](https://pkg.go.dev/badge/github.com/lmcc-dev/lmcc-go-sdk.svg)](https://pkg.go.dev/github.com/lmcc-dev/lmcc-go-sdk)
<!-- Add other badges like build status, coverage later -->

[**中文说明**](./README_zh.md)

`lmcc-go-sdk` is a Go software development kit designed to provide foundational components and utilities for building robust applications.

## ✨ Features

*   **Configuration Management (`pkg/config`):** Flexible loading from files (YAML, TOML), environment variables, and struct tag defaults with hot-reloading capabilities.
*   **Logging (`pkg/log`):** Comprehensive logging capabilities including structured logging (with `zap`), configurable levels, formats (text, JSON), and output paths (console, files). Supports log rotation and dynamic reconfiguration via `pkg/config` for hot-reloading of logging settings. Context-aware logging for enhanced traceability.
*   **(More features to be added)**

## 🚀 Getting Started

### Installation

```bash
go get github.com/lmcc-dev/lmcc-go-sdk
```

### Quick Start Example (Configuration)

```go
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	sdkconfig "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
	sdklog "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
	"time"
)

// Define your application's configuration structure
type ServerConfig struct {
	Host string        `mapstructure:"host" default:"localhost"`
	Port int           `mapstructure:"port" default:"8080"`
	Timeout time.Duration `mapstructure:"timeout" default:"5s"`
}

type AppConfig struct {
	sdkconfig.Config // Embed the base SDK config (optional but recommended)
	Server *ServerConfig `mapstructure:"server"`
	Debug  bool          `mapstructure:"debug" default:"false"`
}

var MyConfig AppConfig

func main() {
	configFile := flag.String("config", "config.yaml", "Path to configuration file (e.g., config.yaml)")
	flag.Parse()

	// Load configuration
	err := sdkconfig.LoadConfig(
		&MyConfig,
		sdkconfig.WithConfigFile(*configFile, ""), // Load from file (type inferred)
		// sdkconfig.WithEnvPrefix("MYAPP"),      // Optionally override default env prefix "LMCC"
		// sdkconfig.WithHotReload(),             // Optionally enable hot reload
	)
	if err != nil {
		// Handle specific error types if needed, e.g., config file not found
		log.Printf("WARN: Failed to load configuration from file '%s', using defaults and env vars: %v\n", *configFile, err)
		// Decide if this is a fatal error or if proceeding with defaults is acceptable
	} else {
		log.Printf("Configuration loaded successfully from %s\n", *configFile)
	}

	// Access configuration values
	fmt.Printf("Server Host: %s\n", MyConfig.Server.Host)
	fmt.Printf("Server Port: %d\n", MyConfig.Server.Port)
	fmt.Printf("Server Timeout: %s\n", MyConfig.Server.Timeout)
	fmt.Printf("Debug Mode: %t\n", MyConfig.Debug)

	// --- SDK Logging Quick Start ---
	// Initialize a simple logger
	logOpts := sdklog.NewOptions()
	logOpts.Level = "info"       // Set desired level (e.g., "debug", "info", "warn")
	logOpts.Format = "console"    // Choose "console" (human-readable) or "json"
	logOpts.OutputPaths = []string{"stdout"} // Log to standard output. Can also be a file path e.g., ["./app.log"]
	logOpts.EnableColor = true // For console output, makes it more readable
	sdklog.Init(logOpts)
	// Important: Defer Sync to flush logs before application exits.
	// It's good practice to call this, especially for file-based logging.
	defer func() {
		if err := sdklog.Sync(); err != nil {
			// Handle logger sync error, e.g., log to standard logger
			// This is unlikely to happen with stdout but good for file logs.
			fmt.Fprintf(os.Stderr, "Failed to sync sdk logger: %v\n", err)
		}
	}()

	sdklog.Info("SDK Logger initialized. This is an INFO message.")
	sdklog.Debugw("This is a DEBUG message with structured fields (won't be visible if level is 'info').", "userID", "user123", "action", "attempt_debug")
	sdklog.Errorw("This is an ERROR message.", "operation", "database_connection", "attempt", 3, "success", false)

	// Contextual logging example
	ctx := context.Background()
	// Typically, a trace ID would come from an incoming request or be generated.
	ctxWithTrace := sdklog.ContextWithTraceID(ctx, "trace-id-example-xyz789") 
	sdklog.Ctx(ctxWithTrace).Infow("Processing payment.", "customerID", "cust999", "amount", 100.50)

	// Note: For advanced logging (file rotation, hot-reloading via pkg/config),
	// please see the detailed pkg/log usage guide in `docs/usage-guides/log/log_usage_en.md`
	// and the comprehensive example in `examples/simple-config-app/main.go`.

	// Example config.yaml:
	/*
	server:
	  host: "127.0.0.1"
	  port: 9090
	debug: true
	*/

	// Example Environment Variables (assuming default prefix LMCC):
	// export LMCC_SERVER_PORT=9999
	// export LMCC_DEBUG=true
}

```

## 📚 Usage Guides

For detailed information on specific modules, please refer to the [Usage Guides](./docs/usage-guides/index_en.md).

## 🤝 Contributing

Contributions are welcome! Please refer to the `CONTRIBUTING.md` file (to be added) for guidelines.

## 📄 License

This project is licensed under the MIT License - see the `LICENSE` file (to be added) for details. 