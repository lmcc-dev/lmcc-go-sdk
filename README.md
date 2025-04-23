# LMCC Go SDK

[![Go Report Card](https://goreportcard.com/badge/github.com/lmcc-dev/lmcc-go-sdk)](https://goreportcard.com/report/github.com/lmcc-dev/lmcc-go-sdk)
[![Go Reference](https://pkg.go.dev/badge/github.com/lmcc-dev/lmcc-go-sdk.svg)](https://pkg.go.dev/github.com/lmcc-dev/lmcc-go-sdk)
<!-- Add other badges like build status, coverage later -->

[**中文说明**](./README_zh.md)

`lmcc-go-sdk` is a Go software development kit designed to provide foundational components and utilities for building robust applications.

## ✨ Features

*   **Configuration Management (`pkg/config`):** Flexible loading from files (YAML, TOML), environment variables, and struct tag defaults with hot-reloading capabilities.
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
	"flag"
	"fmt"
	"log"
	sdkconfig "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
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