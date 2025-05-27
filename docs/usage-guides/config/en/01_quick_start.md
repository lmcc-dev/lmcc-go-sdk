# Quick Start Guide

This guide will help you get started with the configuration module in just a few minutes.

## Installation

The configuration module is part of the lmcc-go-sdk. If you haven't already, add it to your project:

```bash
go mod init your-project
go get github.com/lmcc-dev/lmcc-go-sdk
```

## Basic Usage

### Step 1: Define Your Configuration Structure

Create a struct that represents your application's configuration:

```go
package main

import (
    "log"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
)

type AppConfig struct {
    Server struct {
        Host string `mapstructure:"host" default:"localhost"`
        Port int    `mapstructure:"port" default:"8080"`
    } `mapstructure:"server"`
    
    Database struct {
        URL      string `mapstructure:"url" default:"postgres://localhost/mydb"`
        MaxConns int    `mapstructure:"max_connections" default:"10"`
    } `mapstructure:"database"`
    
    Debug bool `mapstructure:"debug" default:"false"`
}
```

### Step 2: Create a Configuration File

Create a `config.yaml` file in your project root:

```yaml
server:
  host: "0.0.0.0"
  port: 3000

database:
  url: "postgres://user:pass@localhost/production_db"
  max_connections: 25

debug: false
```

### Step 3: Load Configuration

```go
func main() {
    var cfg AppConfig
    
    // Simple configuration loading (no hot-reload)
    err := config.LoadConfig(
        &cfg,
        config.WithConfigFile("config.yaml", ""),
        config.WithEnvPrefix("APP"),
        config.WithEnvVarOverride(true),
    )
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }
    
    // Use your configuration
    log.Printf("Server will run on %s:%d", cfg.Server.Host, cfg.Server.Port)
    log.Printf("Database URL: %s", cfg.Database.URL)
    log.Printf("Debug mode: %t", cfg.Debug)
}
```

## Environment Variable Override

You can override any configuration value using environment variables. With the prefix `APP`, the environment variables would be:

```bash
export APP_SERVER_HOST=production.example.com
export APP_SERVER_PORT=443
export APP_DEBUG=true
```

## Hot Reload Example

For applications that need dynamic configuration updates:

```go
func main() {
    var cfg AppConfig
    
    // Load configuration with hot-reload support
    cm, err := config.LoadConfigAndWatch(
        &cfg,
        config.WithConfigFile("config.yaml", ""),
        config.WithEnvPrefix("APP"),
        config.WithHotReload(true),
        config.WithEnvVarOverride(true),
    )
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }
    
    // Register callback for configuration changes
    cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
        appCfg := currentCfg.(*AppConfig)
        log.Printf("Configuration updated! New port: %d", appCfg.Server.Port)
        
        // Update your application components here
        // For example: restart HTTP server, reconnect to database, etc.
        
        return nil
    })
    
    // Your application logic here
    log.Printf("Application started with config: %+v", cfg)
    
    // Keep the application running to observe hot-reload
    select {}
}
```

## Configuration File Formats

The module supports multiple file formats. Here are examples:

### YAML (Recommended)
```yaml
server:
  host: localhost
  port: 8080
debug: true
```

### JSON
```json
{
  "server": {
    "host": "localhost",
    "port": 8080
  },
  "debug": true
}
```

### TOML
```toml
debug = true

[server]
host = "localhost"
port = 8080
```

## Default Values

Default values are specified using struct tags and will be used when:
- No configuration file is provided
- A field is missing from the configuration file
- Environment variables are not set

```go
type Config struct {
    // Will default to "localhost" if not specified
    Host string `mapstructure:"host" default:"localhost"`
    
    // Will default to 8080 if not specified
    Port int `mapstructure:"port" default:"8080"`
    
    // Will default to false if not specified
    Debug bool `mapstructure:"debug" default:"false"`
}
```

## Error Handling

Always handle configuration loading errors appropriately:

```go
err := config.LoadConfig(&cfg, options...)
if err != nil {
    // Log the error with context
    log.Printf("Configuration error: %v", err)
    
    // Decide how to handle the error:
    // 1. Use default configuration
    // 2. Exit the application
    // 3. Retry loading
    
    log.Fatal("Cannot continue without valid configuration")
}
```

## Next Steps

- [Configuration Options](02_configuration_options.md) - Learn about all available options
- [Hot Reload](03_hot_reload.md) - Implement dynamic configuration updates
- [Best Practices](04_best_practices.md) - Follow recommended patterns
- [Integration Examples](05_integration_examples.md) - See real-world examples 