# Quick Start Guide

This guide will help you get up and running with the lmcc-go-sdk server module in minutes. The server module provides a unified abstraction for multiple web frameworks including Gin, Echo, and Fiber.

## Prerequisites

- Go 1.21 or later
- Basic knowledge of Go web development
- Familiarity with at least one supported framework (optional)

## Installation

The server module is part of the lmcc-go-sdk package:

```bash
go get github.com/lmcc-dev/lmcc-go-sdk
```

## Basic Usage

### 1. Simple Server with Default Configuration

The fastest way to create a server is using the default configuration:

```go
package main

import (
    "context"
    "log"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
)

func main() {
    // Create default configuration
    config := server.DefaultServerConfig()
    config.Framework = "gin"    // Choose framework: gin, echo, or fiber
    config.Port = 8080
    
    // Create server manager
    manager, err := server.CreateServerManager("gin", config)
    if err != nil {
        log.Fatal("Failed to create server:", err)
    }
    
    // Register a simple health check route
    err = manager.GetFramework().RegisterRoute("GET", "/health", 
        server.HandlerFunc(func(ctx server.Context) error {
            return ctx.JSON(200, map[string]string{
                "status": "healthy",
                "framework": "gin",
            })
        }),
    )
    if err != nil {
        log.Fatal("Failed to register route:", err)
    }
    
    // Start the server
    log.Printf("Starting server on port %d", config.Port)
    if err := manager.Start(context.Background()); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
```

### 2. Quick Start Function

For the simplest possible setup, use the `QuickStart` function:

```go
package main

import (
    "log"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
)

func main() {
    config := server.DefaultServerConfig()
    config.Framework = "gin"
    config.Port = 8080
    
    // QuickStart automatically registers plugins and starts the server
    if err := server.QuickStart("gin", config); err != nil {
        log.Fatal("QuickStart failed:", err)
    }
}
```

## Framework-Specific Examples

### Gin Framework

```go
package main

import (
    "context"
    "log"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
)

func main() {
    // Gin-specific configuration
    config := &server.ServerConfig{
        Framework: "gin",
        Host:     "0.0.0.0",
        Port:     8080,
        Mode:     "release",
        GracefulShutdown: server.GracefulShutdownConfig{
            Enabled: true,
            Timeout: 30, // seconds
        },
    }
    
    manager, err := server.CreateServerManager("gin", config)
    if err != nil {
        log.Fatal(err)
    }
    
    // Register routes
    framework := manager.GetFramework()
    
    // Simple route
    _ = framework.RegisterRoute("GET", "/", 
        server.HandlerFunc(func(ctx server.Context) error {
            return ctx.JSON(200, map[string]string{
                "message": "Hello from Gin!",
            })
        }),
    )
    
    // Route with parameter
    _ = framework.RegisterRoute("GET", "/user/:id", 
        server.HandlerFunc(func(ctx server.Context) error {
            id := ctx.Param("id")
            return ctx.JSON(200, map[string]string{
                "user_id": id,
                "framework": "gin",
            })
        }),
    )
    
    // Route group
    api := framework.Group("/api")
    _ = api.RegisterRoute("GET", "/status", 
        server.HandlerFunc(func(ctx server.Context) error {
            return ctx.JSON(200, map[string]string{
                "status": "API is running",
            })
        }),
    )
    
    log.Println("Starting Gin server on :8080")
    if err := manager.Start(context.Background()); err != nil {
        log.Fatal(err)
    }
}
```

### Echo Framework

```go
package main

import (
    "context"
    "log"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
)

func main() {
    config := &server.ServerConfig{
        Framework: "echo",
        Host:     "localhost",
        Port:     8080,
        Mode:     "production",
    }
    
    manager, err := server.CreateServerManager("echo", config)
    if err != nil {
        log.Fatal(err)
    }
    
    framework := manager.GetFramework()
    
    // Echo-style routes
    _ = framework.RegisterRoute("GET", "/", 
        server.HandlerFunc(func(ctx server.Context) error {
            return ctx.JSON(200, map[string]interface{}{
                "message": "Hello from Echo!",
                "timestamp": ctx.Get("timestamp"),
            })
        }),
    )
    
    _ = framework.RegisterRoute("POST", "/users", 
        server.HandlerFunc(func(ctx server.Context) error {
            var user map[string]interface{}
            if err := ctx.Bind(&user); err != nil {
                return ctx.JSON(400, map[string]string{
                    "error": "Invalid JSON",
                })
            }
            
            return ctx.JSON(201, map[string]interface{}{
                "message": "User created",
                "user": user,
            })
        }),
    )
    
    log.Println("Starting Echo server on :8080")
    if err := manager.Start(context.Background()); err != nil {
        log.Fatal(err)
    }
}
```

### Fiber Framework

```go
package main

import (
    "context"
    "log"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
)

func main() {
    config := &server.ServerConfig{
        Framework: "fiber",
        Host:     "0.0.0.0",
        Port:     3000, // Fiber commonly uses port 3000
        Mode:     "production",
    }
    
    manager, err := server.CreateServerManager("fiber", config)
    if err != nil {
        log.Fatal(err)
    }
    
    framework := manager.GetFramework()
    
    // Fiber-style routes
    _ = framework.RegisterRoute("GET", "/", 
        server.HandlerFunc(func(ctx server.Context) error {
            return ctx.JSON(200, map[string]string{
                "message": "Hello from Fiber!",
                "speed": "⚡️ Ultra fast",
            })
        }),
    )
    
    _ = framework.RegisterRoute("GET", "/api/*", 
        server.HandlerFunc(func(ctx server.Context) error {
            path := ctx.Param("*")
            return ctx.JSON(200, map[string]string{
                "wildcard_path": path,
                "framework": "fiber",
            })
        }),
    )
    
    log.Println("Starting Fiber server on :3000")
    if err := manager.Start(context.Background()); err != nil {
        log.Fatal(err)
    }
}
```

## Adding Middleware

The server module provides a unified middleware system that works across all frameworks:

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/middleware"
)

func main() {
    config := server.DefaultServerConfig()
    config.Framework = "gin"
    config.Port = 8080
    
    manager, err := server.CreateServerManager("gin", config)
    if err != nil {
        log.Fatal(err)
    }
    
    framework := manager.GetFramework()
    
    // Add built-in middleware
    _ = framework.RegisterMiddleware(middleware.Logger())
    _ = framework.RegisterMiddleware(middleware.Recovery())
    _ = framework.RegisterMiddleware(middleware.CORS())
    
    // Add custom middleware
    _ = framework.RegisterMiddleware(
        server.MiddlewareFunc(func(next server.Handler) server.Handler {
            return server.HandlerFunc(func(ctx server.Context) error {
                start := time.Now()
                
                // Call next handler
                err := next.Handle(ctx)
                
                // Log request duration
                duration := time.Since(start)
                log.Printf("Request to %s took %v", ctx.Path(), duration)
                
                return err
            })
        }),
    )
    
    // Register routes
    _ = framework.RegisterRoute("GET", "/", 
        server.HandlerFunc(func(ctx server.Context) error {
            return ctx.JSON(200, map[string]string{
                "message": "Hello with middleware!",
            })
        }),
    )
    
    log.Println("Starting server with middleware on :8080")
    if err := manager.Start(context.Background()); err != nil {
        log.Fatal(err)
    }
}
```

## Graceful Shutdown

The server module supports graceful shutdown for production environments:

```go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
)

func main() {
    config := &server.ServerConfig{
        Framework: "gin",
        Host:     "0.0.0.0",
        Port:     8080,
        Mode:     "production",
        GracefulShutdown: server.GracefulShutdownConfig{
            Enabled: true,
            Timeout: 30, // 30 seconds timeout
        },
    }
    
    manager, err := server.CreateServerManager("gin", config)
    if err != nil {
        log.Fatal(err)
    }
    
    framework := manager.GetFramework()
    
    // Register routes
    _ = framework.RegisterRoute("GET", "/health", 
        server.HandlerFunc(func(ctx server.Context) error {
            return ctx.JSON(200, map[string]string{
                "status": "healthy",
            })
        }),
    )
    
    // Start server in a goroutine
    go func() {
        log.Println("Starting server on :8080")
        if err := manager.Start(context.Background()); err != nil {
            log.Printf("Server start error: %v", err)
        }
    }()
    
    // Wait for interrupt signal to gracefully shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    log.Println("Shutting down server...")
    
    // Create shutdown context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := manager.Stop(ctx); err != nil {
        log.Printf("Server shutdown error: %v", err)
    } else {
        log.Println("Server gracefully stopped")
    }
}
```

## Server Factory Pattern

For applications that need to manage multiple server instances or plugins:

```go
package main

import (
    "context"
    "log"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
)

func main() {
    // Create a server factory
    factory := server.NewServerFactory()
    
    // List available plugins
    plugins := factory.ListPlugins()
    log.Printf("Available plugins: %v", plugins)
    
    // Get plugin information
    if len(plugins) > 0 {
        info, err := factory.GetPluginInfo(plugins[0])
        if err == nil {
            log.Printf("Plugin %s: %s", info.Name, info.Description)
        }
    }
    
    // Create server using factory
    config := server.DefaultServerConfig()
    config.Framework = "gin"
    config.Port = 8080
    
    manager, err := factory.CreateServer("gin", config)
    if err != nil {
        log.Fatal(err)
    }
    
    // Register a route
    _ = manager.GetFramework().RegisterRoute("GET", "/", 
        server.HandlerFunc(func(ctx server.Context) error {
            return ctx.JSON(200, map[string]string{
                "message": "Hello from factory-created server!",
            })
        }),
    )
    
    log.Println("Starting factory-created server on :8080")
    if err := manager.Start(context.Background()); err != nil {
        log.Fatal(err)
    }
}
```

## Configuration Overview

Here are the key configuration options:

```go
config := &server.ServerConfig{
    Framework: "gin",           // Required: gin, echo, or fiber
    Host:     "0.0.0.0",       // Server host (default: localhost)
    Port:     8080,             // Server port (default: 8080)
    Mode:     "production",     // Server mode: development, production, test
    
    // TLS Configuration (optional)
    TLS: server.TLSConfig{
        Enabled:  false,
        CertFile: "",
        KeyFile:  "",
    },
    
    // Graceful shutdown (optional)
    GracefulShutdown: server.GracefulShutdownConfig{
        Enabled: true,
        Timeout: 30, // seconds
    },
    
    // Middleware configuration (optional)
    Middleware: server.MiddlewareConfig{
        EnableLogger:   true,
        EnableRecovery: true,
        EnableCORS:     true,
    },
}
```

## Testing Your Server

The server module provides utilities for testing:

```go
package main

import (
    "context"
    "testing"
    "time"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
    "github.com/stretchr/testify/assert"
)

func TestServerStartStop(t *testing.T) {
    config := &server.ServerConfig{
        Framework: "gin",
        Host:     "localhost",
        Port:     8080,
        Mode:     "test",
        GracefulShutdown: server.GracefulShutdownConfig{
            Enabled: false, // Disable for faster tests
        },
    }
    
    manager, err := server.CreateServerManager("gin", config)
    assert.NoError(t, err)
    assert.NotNil(t, manager)
    
    // Test server lifecycle
    ctx := context.Background()
    
    // Start server in goroutine
    done := make(chan error, 1)
    go func() {
        done <- manager.Start(ctx)
    }()
    
    // Wait for server to start
    time.Sleep(10 * time.Millisecond)
    assert.True(t, manager.IsRunning())
    
    // Stop server
    err = manager.Stop(ctx)
    assert.NoError(t, err)
    assert.False(t, manager.IsRunning())
    
    // Wait for start to return
    select {
    case <-done:
        // Normal completion
    case <-time.After(1 * time.Second):
        t.Fatal("Start method did not return within timeout")
    }
}
```

## Next Steps

Now that you have a basic server running, explore these advanced features:

1. **[Configuration Guide](02_configuration.md)** - Deep dive into server configuration options
2. **[Framework Plugins](03_framework_plugins.md)** - Learn about framework-specific features
3. **[Middleware System](04_middleware_system.md)** - Understand the unified middleware architecture
4. **[Server Management](05_server_management.md)** - Advanced lifecycle management
5. **[Integration Examples](06_integration_examples.md)** - Real-world integration patterns

## Common Issues

### Plugin Not Found
If you get a "plugin not found" error, make sure to import the plugin package:

```go
import (
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin"
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/echo"
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/fiber"
)
```

### Port Already in Use
Change the port in your configuration:

```go
config.Port = 8081 // Use a different port
```

### Framework Dependencies
Make sure to install the required framework dependencies:

```bash
# For Gin
go get github.com/gin-gonic/gin

# For Echo
go get github.com/labstack/echo/v4

# For Fiber
go get github.com/gofiber/fiber/v2
```
