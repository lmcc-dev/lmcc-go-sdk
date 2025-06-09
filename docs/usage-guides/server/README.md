# Server Module

The `pkg/server` module provides a unified web framework abstraction that enables plugin-based server architecture with support for multiple popular Go web frameworks including Gin, Echo, and Fiber. This module is designed to simplify server development while maintaining flexibility and performance.

## Quick Links

- **[中文文档](README_zh.md)** - Chinese documentation
- **[Quick Start Guide](en/01_quick_start.md)** - Get started in minutes
- **[Framework Plugins](en/03_framework_plugins.md)** - Gin, Echo, Fiber support
- **[Middleware System](en/04_middleware_system.md)** - Unified middleware architecture
- **[Server Management](en/05_server_management.md)** - Lifecycle management
- **[Integration Examples](en/06_integration_examples.md)** - Real-world examples
- **[Best Practices](en/07_best_practices.md)** - Production recommendations
- **[Module Specification](en/09_module_specification.md)** - Complete API reference

## Features

### 🔌 Plugin-Based Architecture
- **Framework Abstraction**: Unified interface for different web frameworks
- **Hot-Swappable Plugins**: Support for Gin, Echo, and Fiber frameworks
- **Custom Plugin Development**: Extensible plugin system for new frameworks
- **Framework-Agnostic Code**: Write once, run on any supported framework

### 🚀 High Performance
- **Efficient Routing**: Optimized for high-throughput applications
- **Minimal Overhead**: Lightweight abstraction layer
- **Connection Pooling**: Built-in connection management
- **Graceful Shutdown**: Zero-downtime deployments

### 🛠️ Comprehensive Middleware
- **Built-in Middleware**: Logger, Recovery, CORS, Rate Limiting, Auth, and more
- **Unified Interface**: Consistent middleware API across all frameworks
- **Custom Middleware**: Easy development of application-specific middleware
- **Middleware Chains**: Flexible middleware composition and ordering

### ⚙️ Advanced Configuration
- **Flexible Configuration**: Comprehensive server configuration options
- **Hot Reload**: Dynamic configuration updates without restart
- **Environment-Specific**: Support for different deployment environments
- **Validation**: Built-in configuration validation and error handling

### 🔗 Deep Integration
- **Config Module**: Seamless integration with `pkg/config`
- **Logging Module**: Built-in integration with `pkg/log`
- **Error Handling**: Native support for `pkg/errors`
- **Service Container**: Dependency injection and service management

## Quick Example

```go
package main

import (
    "context"
    "log"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
)

func main() {
    // Create server configuration
    config := server.DefaultServerConfig()
    config.Framework = "gin"
    config.Host = "0.0.0.0"
    config.Port = 8080
    
    // Create server manager
    manager, err := server.CreateServerManager("gin", config)
    if err != nil {
        log.Fatal(err)
    }
    
    // Register a simple route
    err = manager.GetFramework().RegisterRoute("GET", "/health", 
        server.HandlerFunc(func(ctx server.Context) error {
            return ctx.JSON(200, map[string]string{
                "status": "healthy",
                "framework": "gin",
            })
        }),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Start server
    log.Println("Starting server on :8080")
    if err := manager.Start(context.Background()); err != nil {
        log.Fatal(err)
    }
}
```

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                        │
├─────────────────────────────────────────────────────────────┤
│                 Unified Server Interface                    │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────────┐ │
│  │   Router    │ │ Middleware  │ │    Server Manager       │ │
│  │   Group     │ │   Chain     │ │   (Lifecycle)           │ │
│  └─────────────┘ └─────────────┘ └─────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                   Plugin Architecture                       │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────────┐ │
│  │ Gin Plugin  │ │Echo Plugin  │ │    Fiber Plugin         │ │
│  │ (Adapter)   │ │ (Adapter)   │ │    (Adapter)            │ │
│  └─────────────┘ └─────────────┘ └─────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                Framework Native Layer                       │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────────┐ │
│  │Gin Engine   │ │Echo Instance│ │    Fiber App            │ │
│  └─────────────┘ └─────────────┘ └─────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## Supported Frameworks

| Framework | Version | Status | Features |
|-----------|---------|--------|----------|
| **Gin** | v1.9+ | ✅ Production Ready | High performance, rich middleware ecosystem |
| **Echo** | v4.0+ | ✅ Production Ready | Minimalist, fast, extensible |
| **Fiber** | v2.0+ | ✅ Production Ready | Express-inspired, ultra-fast |

## Installation

The server module is part of the lmcc-go-sdk:

```bash
go get github.com/lmcc-dev/lmcc-go-sdk
```

## Framework-Specific Quick Start

### Gin Framework
```go
config := server.DefaultServerConfig()
config.Framework = "gin"
manager, _ := server.CreateServerManager("gin", config)
```

### Echo Framework
```go
config := server.DefaultServerConfig()
config.Framework = "echo"
manager, _ := server.CreateServerManager("echo", config)
```

### Fiber Framework
```go
config := server.DefaultServerConfig()
config.Framework = "fiber"
manager, _ := server.CreateServerManager("fiber", config)
```

## Integration with Other Modules

The server module seamlessly integrates with other SDK modules:

```go
import (
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

type AppConfig struct {
    Server server.ServerConfig `mapstructure:"server"`
    Log    log.Options         `mapstructure:"log"`
}

func main() {
    var cfg AppConfig
    
    // Load configuration
    config.LoadConfig(&cfg, 
        config.WithConfigFile("config.yaml", ""),
        config.WithHotReload(true),
    )
    
    // Initialize logging
    log.Init(&cfg.Log)
    
    // Create server with integrated configuration
    manager, _ := server.CreateServerManager(cfg.Server.Framework, &cfg.Server)
    
    // Start server
    manager.Start(context.Background())
}
```

## Getting Started

1. **[Quick Start Guide](en/01_quick_start.md)** - Basic setup and usage
2. **[Configuration Guide](en/02_configuration.md)** - Server configuration options
3. **[Framework Plugins](en/03_framework_plugins.md)** - Choose your framework
4. **[Middleware System](en/04_middleware_system.md)** - Add functionality with middleware

## Contributing

Please read our [contributing guidelines](../../../CONTRIBUTING.md) before submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](../../../LICENSE) file for details. 