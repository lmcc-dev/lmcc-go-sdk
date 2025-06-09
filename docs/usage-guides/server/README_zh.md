# Server 模块

`pkg/server` 模块提供了统一的 Web 框架抽象，支持基于插件的服务器架构，兼容多种流行的 Go Web 框架，包括 Gin、Echo 和 Fiber。此模块旨在简化服务器开发，同时保持灵活性和高性能。

## 快速链接

- **[英文文档](README.md)** - English documentation
- **[快速开始指南](zh/01_quick_start.md)** - 几分钟内上手
- **[框架插件](zh/03_framework_plugins.md)** - Gin、Echo、Fiber 支持
- **[中间件系统](zh/04_middleware_system.md)** - 统一中间件架构
- **[服务器管理](zh/05_server_management.md)** - 生命周期管理
- **[集成示例](zh/06_integration_examples.md)** - 实际应用示例
- **[最佳实践](zh/07_best_practices.md)** - 生产环境建议
- **[模块规范](zh/09_module_specification.md)** - 完整 API 参考

## 特性

### 🔌 基于插件的架构
- **框架抽象**: 为不同 Web 框架提供统一接口
- **热插拔插件**: 支持 Gin、Echo 和 Fiber 框架
- **自定义插件开发**: 为新框架提供可扩展的插件系统
- **框架无关代码**: 一次编写，在任何支持的框架上运行

### 🚀 高性能
- **高效路由**: 针对高吞吐量应用进行优化
- **最小开销**: 轻量级抽象层
- **连接池**: 内置连接管理
- **优雅关闭**: 零停机部署

### 🛠️ 完整的中间件系统
- **内置中间件**: 日志记录、恢复、跨域、限流、认证等
- **统一接口**: 所有框架中一致的中间件 API
- **自定义中间件**: 轻松开发应用特定的中间件
- **中间件链**: 灵活的中间件组合和排序

### ⚙️ 高级配置
- **灵活配置**: 全面的服务器配置选项
- **热重载**: 无需重启的动态配置更新
- **环境特定**: 支持不同部署环境
- **验证**: 内置配置验证和错误处理

### 🔗 深度集成
- **配置模块**: 与 `pkg/config` 无缝集成
- **日志模块**: 与 `pkg/log` 内置集成
- **错误处理**: 原生支持 `pkg/errors`
- **服务容器**: 依赖注入和服务管理

## 快速示例

```go
package main

import (
    "context"
    "log"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
)

func main() {
    // 创建服务器配置 (Create server configuration)
    config := server.DefaultServerConfig()
    config.Framework = "gin"
    config.Host = "0.0.0.0"
    config.Port = 8080
    
    // 创建服务器管理器 (Create server manager)
    manager, err := server.CreateServerManager("gin", config)
    if err != nil {
        log.Fatal(err)
    }
    
    // 注册简单路由 (Register a simple route)
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
    
    // 启动服务器 (Start server)
    log.Println("Starting server on :8080")
    if err := manager.Start(context.Background()); err != nil {
        log.Fatal(err)
    }
}
```

## 架构概览

```
┌─────────────────────────────────────────────────────────────┐
│                    应用层 (Application Layer)               │
├─────────────────────────────────────────────────────────────┤
│                 统一服务器接口 (Unified Interface)           │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────────┐ │
│  │   路由组    │ │   中间件    │ │    服务器管理器         │ │
│  │  (Router)   │ │ (Middleware)│ │ (Server Manager)        │ │
│  └─────────────┘ └─────────────┘ └─────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                   插件架构 (Plugin Architecture)            │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────────┐ │
│  │ Gin 插件    │ │Echo 插件    │ │    Fiber 插件           │ │
│  │ (Adapter)   │ │ (Adapter)   │ │    (Adapter)            │ │
│  └─────────────┘ └─────────────┘ └─────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                框架原生层 (Framework Native Layer)          │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────────┐ │
│  │Gin Engine   │ │Echo Instance│ │    Fiber App            │ │
│  └─────────────┘ └─────────────┘ └─────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## 支持的框架

| 框架 | 版本 | 状态 | 特性 |
|------|------|------|------|
| **Gin** | v1.9+ | ✅ 生产就绪 | 高性能，丰富的中间件生态系统 |
| **Echo** | v4.0+ | ✅ 生产就绪 | 极简主义，快速，可扩展 |
| **Fiber** | v2.0+ | ✅ 生产就绪 | 受 Express 启发，超快速 |

## 安装

server 模块是 lmcc-go-sdk 的一部分：

```bash
go get github.com/lmcc-dev/lmcc-go-sdk
```

## 框架特定的快速开始

### Gin 框架
```go
config := server.DefaultServerConfig()
config.Framework = "gin"
manager, _ := server.CreateServerManager("gin", config)
```

### Echo 框架
```go
config := server.DefaultServerConfig()
config.Framework = "echo"
manager, _ := server.CreateServerManager("echo", config)
```

### Fiber 框架
```go
config := server.DefaultServerConfig()
config.Framework = "fiber"
manager, _ := server.CreateServerManager("fiber", config)
```

## 与其他模块的集成

server 模块与其他 SDK 模块无缝集成：

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
    
    // 加载配置 (Load configuration)
    config.LoadConfig(&cfg, 
        config.WithConfigFile("config.yaml", ""),
        config.WithHotReload(true),
    )
    
    // 初始化日志 (Initialize logging)
    log.Init(&cfg.Log)
    
    // 使用集成配置创建服务器 (Create server with integrated configuration)
    manager, _ := server.CreateServerManager(cfg.Server.Framework, &cfg.Server)
    
    // 启动服务器 (Start server)
    manager.Start(context.Background())
}
```

## 开始使用

1. **[快速开始指南](zh/01_quick_start.md)** - 基础设置和使用
2. **[配置指南](zh/02_configuration.md)** - 服务器配置选项
3. **[框架插件](zh/03_framework_plugins.md)** - 选择您的框架
4. **[中间件系统](zh/04_middleware_system.md)** - 使用中间件添加功能

## 贡献

在提交拉取请求之前，请阅读我们的[贡献指南](../../../CONTRIBUTING.md)。

## 许可证

此项目采用 MIT 许可证 - 有关详细信息，请参见 [LICENSE](../../../LICENSE) 文件。 