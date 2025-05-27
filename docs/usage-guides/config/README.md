# Configuration Management Module

The `pkg/config` module provides flexible and robust configuration management capabilities for Go applications, inspired by best practices seen in ecosystems like Marmotedu.

## Quick Links

- **[中文文档](README_zh.md)** - Chinese documentation
- **[Quick Start Guide](en/01_quick_start.md)** - Get started in minutes
- **[Configuration Options](en/02_configuration_options.md)** - All available options
- **[Hot Reload](en/03_hot_reload.md)** - Dynamic configuration updates
- **[Best Practices](en/04_best_practices.md)** - Recommended patterns
- **[Integration Examples](en/05_integration_examples.md)** - Real-world examples
- **[Troubleshooting](en/06_troubleshooting.md)** - Common issues and solutions
- **[Module Specification](en/07_module_specification.md)** - Complete API reference

## Features

### 🚀 High Performance
- Built on Viper library for efficient configuration management
- Minimal overhead for configuration access
- Optimized for high-frequency configuration reads

### 📝 Multiple Configuration Sources
- **Files**: YAML, JSON, TOML, and more
- **Environment Variables**: Automatic binding with prefix support
- **Default Values**: Set defaults using struct tags
- **Command Line**: Integration with flag packages

### 🔄 Dynamic Configuration
- **Hot Reload**: Automatic configuration reloading when files change
- **Callback System**: Register callbacks for configuration changes
- **Watch Mode**: Real-time configuration monitoring
- **Graceful Updates**: Non-disruptive configuration updates

### 🎯 Type Safety
- **Strong Typing**: User-defined structs for configuration
- **Validation**: Built-in validation support
- **Automatic Unmarshaling**: Direct mapping to Go structs
- **Error Handling**: Comprehensive error reporting

### ⚙️ Easy Integration
- **Simple API**: Minimal setup required
- **Framework Agnostic**: Works with any Go application
- **Middleware Support**: Easy integration with web frameworks
- **Testing Friendly**: Mock and test configuration easily

## Quick Example

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
    Debug bool `mapstructure:"debug" default:"false"`
}

func main() {
    var cfg AppConfig
    
    cm, err := config.LoadConfigAndWatch(
        &cfg,
        config.WithConfigFile("config.yaml", ""),
        config.WithEnvPrefix("APP"),
        config.WithHotReload(true),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Register callback for changes
    cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
        log.Println("Configuration updated!")
        return nil
    })
    
    // Use configuration
    log.Printf("Server running on %s:%d", cfg.Server.Host, cfg.Server.Port)
}
```

## Installation

The config module is part of the lmcc-go-sdk:

```bash
go get github.com/lmcc-dev/lmcc-go-sdk
```

## Basic Configuration

### Simple Configuration

```go
import "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"

var cfg MyConfig
err := config.LoadConfig(&cfg)
```

### Advanced Configuration

```go
cm, err := config.LoadConfigAndWatch(
    &cfg,
    config.WithConfigFile("config.yaml", ""),
    config.WithEnvPrefix("APP"),
    config.WithHotReload(true),
)
```

### YAML Configuration

```yaml
# config.yaml
server:
  host: "localhost"
  port: 8080
  timeout: "30s"
database:
  host: "localhost"
  port: 5432
  name: "myapp"
debug: false
```

## Integration with Other Modules

The config module integrates seamlessly with other SDK modules:

```go
import (
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

type AppConfig struct {
    Log    log.Options    `mapstructure:"log"`
    Server ServerConfig   `mapstructure:"server"`
}

func main() {
    var cfg AppConfig
    
    // Load configuration
    cm, err := config.LoadConfigAndWatch(&cfg, 
        config.WithConfigFile("config.yaml", ""),
        config.WithHotReload(true),
    )
    if err != nil {
        panic(err)
    }
    
    // Initialize logging with config
    log.Init(&cfg.Log)
    
    // Register for hot reload
    log.RegisterConfigHotReload(cm)
    
    log.Info("Application started with integrated configuration")
}
```

## Getting Started

1. **[Quick Start Guide](en/01_quick_start.md)** - Basic setup and usage
2. **[Configuration Options](en/02_configuration_options.md)** - Detailed configuration
3. **[Hot Reload](en/03_hot_reload.md)** - Dynamic updates
4. **[Best Practices](en/04_best_practices.md)** - Production recommendations

## Contributing

Please read our [contributing guidelines](../../../CONTRIBUTING.md) before submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](../../../LICENSE) file for details.