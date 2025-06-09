# 快速开始指南

本指南将帮助您在几分钟内快速上手 lmcc-go-sdk 服务器模块。服务器模块为多种 Web 框架（包括 Gin、Echo 和 Fiber）提供了统一的抽象接口。

## 前置要求

- Go 1.21 或更高版本
- 基本的 Go Web 开发知识
- 熟悉至少一种支持的框架（可选）

## 安装

服务器模块是 lmcc-go-sdk 包的一部分：

```bash
go get github.com/lmcc-dev/lmcc-go-sdk
```

## 基本用法

### 1. 使用默认配置的简单服务器

创建服务器的最快方法是使用默认配置：

```go
package main

import (
    "context"
    "log"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
)

func main() {
    // 创建默认配置 (Create default configuration)
    config := server.DefaultServerConfig()
    config.Framework = "gin"    // 选择框架: gin, echo, or fiber (Choose framework)
    config.Port = 8080
    
    // 创建服务器管理器 (Create server manager)
    manager, err := server.CreateServerManager("gin", config)
    if err != nil {
        log.Fatal("Failed to create server:", err)
    }
    
    // 注册一个简单的健康检查路由 (Register a simple health check route)
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
    
    // 启动服务器 (Start the server)
    log.Printf("Starting server on port %d", config.Port)
    if err := manager.Start(context.Background()); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
```

### 2. 快速启动函数

对于最简单的设置，使用 `QuickStart` 函数：

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
    
    // QuickStart 自动注册插件并启动服务器 (QuickStart automatically registers plugins and starts the server)
    if err := server.QuickStart("gin", config); err != nil {
        log.Fatal("QuickStart failed:", err)
    }
}
```

## 框架特定示例

### Gin 框架

```go
package main

import (
    "context"
    "log"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
)

func main() {
    // Gin 特定配置 (Gin-specific configuration)
    config := &server.ServerConfig{
        Framework: "gin",
        Host:     "0.0.0.0",
        Port:     8080,
        Mode:     "release",
        GracefulShutdown: server.GracefulShutdownConfig{
            Enabled: true,
            Timeout: 30, // 秒 (seconds)
        },
    }
    
    manager, err := server.CreateServerManager("gin", config)
    if err != nil {
        log.Fatal(err)
    }
    
    // 注册路由 (Register routes)
    framework := manager.GetFramework()
    
    // 简单路由 (Simple route)
    _ = framework.RegisterRoute("GET", "/", 
        server.HandlerFunc(func(ctx server.Context) error {
            return ctx.JSON(200, map[string]string{
                "message": "Hello from Gin!",
            })
        }),
    )
    
    // 带参数的路由 (Route with parameter)
    _ = framework.RegisterRoute("GET", "/user/:id", 
        server.HandlerFunc(func(ctx server.Context) error {
            id := ctx.Param("id")
            return ctx.JSON(200, map[string]string{
                "user_id": id,
                "framework": "gin",
            })
        }),
    )
    
    // 路由组 (Route group)
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

### Echo 框架

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
    
    // Echo 风格的路由 (Echo-style routes)
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

### Fiber 框架

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
        Port:     3000, // Fiber 通常使用端口 3000 (Fiber commonly uses port 3000)
        Mode:     "production",
    }
    
    manager, err := server.CreateServerManager("fiber", config)
    if err != nil {
        log.Fatal(err)
    }
    
    framework := manager.GetFramework()
    
    // Fiber 风格的路由 (Fiber-style routes)
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

## 添加中间件

服务器模块提供了统一的中间件系统，可在所有框架中工作：

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
    
    // 添加内置中间件 (Add built-in middleware)
    _ = framework.RegisterMiddleware(middleware.Logger())
    _ = framework.RegisterMiddleware(middleware.Recovery())
    _ = framework.RegisterMiddleware(middleware.CORS())
    
    // 添加自定义中间件 (Add custom middleware)
    _ = framework.RegisterMiddleware(
        server.MiddlewareFunc(func(next server.Handler) server.Handler {
            return server.HandlerFunc(func(ctx server.Context) error {
                start := time.Now()
                
                // 调用下一个处理器 (Call next handler)
                err := next.Handle(ctx)
                
                // 记录请求持续时间 (Log request duration)
                duration := time.Since(start)
                log.Printf("Request to %s took %v", ctx.Path(), duration)
                
                return err
            })
        }),
    )
    
    // 注册路由 (Register routes)
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

## 优雅关闭

服务器模块支持生产环境的优雅关闭：

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
            Timeout: 30, // 30 秒超时 (30 seconds timeout)
        },
    }
    
    manager, err := server.CreateServerManager("gin", config)
    if err != nil {
        log.Fatal(err)
    }
    
    framework := manager.GetFramework()
    
    // 注册路由 (Register routes)
    _ = framework.RegisterRoute("GET", "/health", 
        server.HandlerFunc(func(ctx server.Context) error {
            return ctx.JSON(200, map[string]string{
                "status": "healthy",
            })
        }),
    )
    
    // 在 goroutine 中启动服务器 (Start server in a goroutine)
    go func() {
        log.Println("Starting server on :8080")
        if err := manager.Start(context.Background()); err != nil {
            log.Printf("Server start error: %v", err)
        }
    }()
    
    // 等待中断信号以优雅关闭 (Wait for interrupt signal to gracefully shutdown)
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    log.Println("Shutting down server...")
    
    // 创建带超时的关闭上下文 (Create shutdown context with timeout)
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := manager.Stop(ctx); err != nil {
        log.Printf("Server shutdown error: %v", err)
    } else {
        log.Println("Server gracefully stopped")
    }
}
```

## 服务器工厂模式

对于需要管理多个服务器实例或插件的应用程序：

```go
package main

import (
    "context"
    "log"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
)

func main() {
    // 创建服务器工厂 (Create a server factory)
    factory := server.NewServerFactory()
    
    // 列出可用插件 (List available plugins)
    plugins := factory.ListPlugins()
    log.Printf("Available plugins: %v", plugins)
    
    // 获取插件信息 (Get plugin information)
    if len(plugins) > 0 {
        info, err := factory.GetPluginInfo(plugins[0])
        if err == nil {
            log.Printf("Plugin %s: %s", info.Name, info.Description)
        }
    }
    
    // 使用工厂创建服务器 (Create server using factory)
    config := server.DefaultServerConfig()
    config.Framework = "gin"
    config.Port = 8080
    
    manager, err := factory.CreateServer("gin", config)
    if err != nil {
        log.Fatal(err)
    }
    
    // 注册路由 (Register a route)
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

## 配置概览

以下是关键的配置选项：

```go
config := &server.ServerConfig{
    Framework: "gin",           // 必需：gin, echo, 或 fiber (Required)
    Host:     "0.0.0.0",       // 服务器主机 (默认：localhost) (Server host)
    Port:     8080,             // 服务器端口 (默认：8080) (Server port)
    Mode:     "production",     // 服务器模式：development, production, test (Server mode)
    
    // TLS 配置 (可选) (TLS Configuration - optional)
    TLS: server.TLSConfig{
        Enabled:  false,
        CertFile: "",
        KeyFile:  "",
    },
    
    // 优雅关闭 (可选) (Graceful shutdown - optional)
    GracefulShutdown: server.GracefulShutdownConfig{
        Enabled: true,
        Timeout: 30, // 秒 (seconds)
    },
    
    // 中间件配置 (可选) (Middleware configuration - optional)
    Middleware: server.MiddlewareConfig{
        EnableLogger:   true,
        EnableRecovery: true,
        EnableCORS:     true,
    },
}
```

## 测试您的服务器

服务器模块提供了测试工具：

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
            Enabled: false, // 对于更快的测试禁用 (Disable for faster tests)
        },
    }
    
    manager, err := server.CreateServerManager("gin", config)
    assert.NoError(t, err)
    assert.NotNil(t, manager)
    
    // 测试服务器生命周期 (Test server lifecycle)
    ctx := context.Background()
    
    // 在 goroutine 中启动服务器 (Start server in goroutine)
    done := make(chan error, 1)
    go func() {
        done <- manager.Start(ctx)
    }()
    
    // 等待服务器启动 (Wait for server to start)
    time.Sleep(10 * time.Millisecond)
    assert.True(t, manager.IsRunning())
    
    // 停止服务器 (Stop server)
    err = manager.Stop(ctx)
    assert.NoError(t, err)
    assert.False(t, manager.IsRunning())
    
    // 等待启动返回 (Wait for start to return)
    select {
    case <-done:
        // 正常完成 (Normal completion)
    case <-time.After(1 * time.Second):
        t.Fatal("Start method did not return within timeout")
    }
}
```

## 下一步

现在您有一个基本的服务器运行，探索这些高级功能：

1. **[配置指南](02_configuration.md)** - 深入了解服务器配置选项
2. **[框架插件](03_framework_plugins.md)** - 了解框架特定功能
3. **[中间件系统](04_middleware_system.md)** - 理解统一中间件架构
4. **[服务器管理](05_server_management.md)** - 高级生命周期管理
5. **[集成示例](06_integration_examples.md)** - 真实世界集成模式

## 常见问题

### 找不到插件

如果遇到"plugin not found"错误，请确保导入插件包：

```go
import (
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin"
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/echo"
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/fiber"
)
```

### 端口已在使用

在配置中更改端口：

```go
config.Port = 8081 // 使用不同的端口 (Use a different port)
```

### 框架依赖

确保安装所需的框架依赖：

```bash
# 对于 Gin (For Gin)
go get github.com/gin-gonic/gin

# 对于 Echo (For Echo)
go get github.com/labstack/echo/v4

# 对于 Fiber (For Fiber)
go get github.com/gofiber/fiber/v2
```