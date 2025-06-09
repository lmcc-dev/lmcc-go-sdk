# 框架插件指南

lmcc-go-sdk 服务器模块使用插件架构来支持多种 Web 框架。本指南涵盖如何使用内置框架插件（Gin、Echo、Fiber）以及如何开发自定义插件。

## 插件架构概述

### 核心组件

插件系统由几个关键组件组成：

1. **FrameworkPlugin 接口**: 定义插件契约
2. **插件注册表**: 管理插件注册和发现
3. **框架适配器**: 实现框架特定逻辑
4. **服务集成**: 连接插件与其他 SDK 模块

### 插件接口

每个框架插件都实现 `FrameworkPlugin` 接口：

```go
type FrameworkPlugin interface {
    Name() string                                    // 插件标识符 (Plugin identifier)
    Version() string                                 // 插件版本 (Plugin version)
    Description() string                             // 插件描述 (Plugin description)
    DefaultConfig() interface{}                      // 默认配置 (Default configuration)
    CreateFramework(config, services) (WebFramework, error)  // 工厂方法 (Factory method)
    ValidateConfig(config interface{}) error         // 配置验证 (Configuration validation)
    GetConfigSchema() interface{}                    // 配置模式 (Configuration schema)
}
```

### 插件注册

插件在包导入时自动注册：

```go
import (
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin"
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/echo"
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/fiber"
)
```

或手动注册：

```go
plugin := myplugin.NewPlugin()
err := server.RegisterFramework(plugin)
if err != nil {
    panic("Failed to register plugin: " + err.Error())
}
```

## 内置框架插件

### Gin 插件

Gin 插件提供高性能 HTTP 路由和丰富的生态系统。

#### 特性
- **高性能**: 快速路由和最小内存占用
- **丰富生态**: 广泛的中间件库
- **JSON/XML 绑定**: 内置请求/响应绑定
- **模板渲染**: HTML 模板支持
- **测试支持**: 内置测试工具

#### 配置

```go
config := &server.ServerConfig{
    Framework: "gin",
    Host:     "0.0.0.0",
    Port:     8080,
    Mode:     "release", // gin 模式: debug, release, test (gin modes)
    
    Plugins: map[string]interface{}{
        "gin": map[string]interface{}{
            "trusted_proxies":           []string{"127.0.0.1", "10.0.0.0/8"},
            "redirect_trailing_slash":   true,
            "redirect_fixed_path":       false,
            "handle_method_not_allowed": false,
            "max_multipart_memory":     32 << 20, // 32MB
            "use_h2c":                  false,
            "forward_by_client_ip":     true,
            "use_raw_path":             false,
            "unescape_path_values":     true,
        },
    },
}
```

#### 基本用法

```go
package main

import (
    "context"
    "log"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin"
)

func main() {
    // 使用 Gin 插件创建服务器 (Create server with Gin plugin)
    config := server.DefaultServerConfig()
    config.Framework = "gin"
    config.Port = 8080
    
    manager, err := server.CreateServerManager("gin", config)
    if err != nil {
        log.Fatal("Failed to create server:", err)
    }
    
    // 获取框架实例 (Get framework instance)
    framework := manager.GetFramework()
    
    // 注册路由 (Register routes)
    framework.RegisterRoute("GET", "/hello", server.HandlerFunc(func(ctx server.Context) error {
        return ctx.JSON(200, map[string]string{
            "message": "Hello from Gin!",
            "framework": "gin",
        })
    }))
    
    // 启动服务器 (Start server)
    if err := manager.Start(context.Background()); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
```

#### Gin 特定功能

访问 Gin 的原生引擎以使用框架特定功能：

```go
framework := manager.GetFramework()
ginEngine := framework.GetNativeEngine().(*gin.Engine)

// 使用 Gin 特定功能 (Use Gin-specific features)
ginEngine.HTMLRender = &ginrender.HTMLDebug{Glob: "templates/*"}
ginEngine.Static("/static", "./static")

// 自定义 Gin 中间件 (Custom Gin middleware)
ginEngine.Use(func(c *gin.Context) {
    // Gin 特定中间件逻辑 (Gin-specific middleware logic)
    c.Header("X-Framework", "Gin")
    c.Next()
})
```

#### 性能调优

```go
config.Plugins["gin"] = map[string]interface{}{
    "trusted_proxies":       []string{},  // 如果不在代理后面则禁用 (Disable if not behind proxy)
    "redirect_trailing_slash": false,     // 禁用以获得更好性能 (Disable for better performance)
    "use_raw_path":         true,         // 跳过路径清理 (Skip path cleaning)
    "max_multipart_memory": 1 << 20,     // 减少内存使用 (Reduce memory usage)
}
```

### Echo 插件

Echo 插件提供极简主义和可扩展的 Web 框架。

#### 特性
- **极简设计**: 清洁简单的架构
- **高性能**: 快速路由，零内存分配
- **中间件**: 广泛的中间件支持
- **数据绑定**: JSON、XML、表单数据绑定
- **自动 TLS**: 内置 Let's Encrypt 支持

#### 配置

```go
config := &server.ServerConfig{
    Framework: "echo",
    Host:     "0.0.0.0",
    Port:     8080,
    Mode:     "production",
    
    Plugins: map[string]interface{}{
        "echo": map[string]interface{}{
            "hide_banner":              true,
            "hide_port":                false,
            "debug":                   false,
            "disable_startup_message":  false,
            "http_error_handler":       nil,
            "binder":                  nil,
            "validator":               nil,
            "renderer":                nil,
            "logger":                  nil,
        },
    },
}
```

#### 基本用法

```go
package main

import (
    "context"
    "log"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/echo"
)

func main() {
    // 使用 Echo 插件创建服务器 (Create server with Echo plugin)
    config := server.DefaultServerConfig()
    config.Framework = "echo"
    config.Port = 8080
    
    manager, err := server.CreateServerManager("echo", config)
    if err != nil {
        log.Fatal("Failed to create server:", err)
    }
    
    framework := manager.GetFramework()
    
    // 注册路由 (Register routes)
    framework.RegisterRoute("GET", "/hello", server.HandlerFunc(func(ctx server.Context) error {
        return ctx.JSON(200, map[string]string{
            "message": "Hello from Echo!",
            "framework": "echo",
        })
    }))
    
    // REST API 示例 (REST API example)
    api := framework.Group("/api/v1")
    api.RegisterRoute("GET", "/users", server.HandlerFunc(getUsersHandler))
    api.RegisterRoute("POST", "/users", server.HandlerFunc(createUserHandler))
    
    // 启动服务器 (Start server)
    if err := manager.Start(context.Background()); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}

func getUsersHandler(ctx server.Context) error {
    // 模拟用户数据 (Simulate user data)
    users := []map[string]interface{}{
        {"id": 1, "name": "John Doe"},
        {"id": 2, "name": "Jane Smith"},
    }
    return ctx.JSON(200, users)
}

func createUserHandler(ctx server.Context) error {
    var user map[string]interface{}
    if err := ctx.Bind(&user); err != nil {
        return ctx.JSON(400, map[string]string{"error": "Invalid JSON"})
    }
    
    // 模拟用户创建 (Simulate user creation)
    user["id"] = 123
    return ctx.JSON(201, user)
}
```

#### Echo 特定功能

```go
framework := manager.GetFramework()
echoInstance := framework.GetNativeEngine().(*echo.Echo)

// 自定义错误处理器 (Custom error handler)
echoInstance.HTTPErrorHandler = func(err error, c echo.Context) {
    code := http.StatusInternalServerError
    message := "Internal Server Error"
    
    if he, ok := err.(*echo.HTTPError); ok {
        code = he.Code
        message = he.Message.(string)
    }
    
    c.JSON(code, map[string]string{"error": message})
}

// 自定义验证器 (Custom validator)
echoInstance.Validator = &CustomValidator{}

// 模板渲染器 (Template renderer)
echoInstance.Renderer = &TemplateRenderer{
    templates: template.Must(template.ParseGlob("templates/*.html")),
}
```

### Fiber 插件

Fiber 插件提供受 Express.js 启发的高性能框架。

#### 特性
- **Express 风格**: Node.js 开发者熟悉的 API
- **高性能**: 基于 fasthttp 构建，最大速度
- **低内存**: 最小内存占用
- **内置功能**: 压缩、缓存、限流
- **WebSocket 支持**: 内置 WebSocket 支持

#### 配置

```go
config := &server.ServerConfig{
    Framework: "fiber",
    Host:     "0.0.0.0",
    Port:     8080,
    Mode:     "production",
    
    Plugins: map[string]interface{}{
        "fiber": map[string]interface{}{
            "prefork":              false,
            "strict_routing":       false,
            "case_sensitive":       false,
            "immutable":           false,
            "unescaped_path":      false,
            "etag":                false,
            "body_limit":          4 * 1024 * 1024, // 4MB
            "concurrency":         256 * 1024,      // 256k 并发连接 (concurrent connections)
            "read_timeout":        "30s",
            "write_timeout":       "30s",
            "idle_timeout":        "60s",
            "read_buffer_size":    4096,
            "write_buffer_size":   4096,
            "compress_body":       false,
            "network":             "tcp4",
        },
    },
}
```

#### 基本用法

```go
package main

import (
    "context"
    "log"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/fiber"
)

func main() {
    // 使用 Fiber 插件创建服务器 (Create server with Fiber plugin)
    config := server.DefaultServerConfig()
    config.Framework = "fiber"
    config.Port = 8080
    
    manager, err := server.CreateServerManager("fiber", config)
    if err != nil {
        log.Fatal("Failed to create server:", err)
    }
    
    framework := manager.GetFramework()
    
    // 注册路由 (Register routes)
    framework.RegisterRoute("GET", "/hello", server.HandlerFunc(func(ctx server.Context) error {
        return ctx.JSON(200, map[string]string{
            "message": "Hello from Fiber!",
            "framework": "fiber",
        })
    }))
    
    // 文件上传示例 (File upload example)
    framework.RegisterRoute("POST", "/upload", server.HandlerFunc(fileUploadHandler))
    
    // WebSocket 端点（使用原生 Fiber 功能） (WebSocket endpoint using native Fiber features)
    setupWebSocket(framework)
    
    // 启动服务器 (Start server)
    if err := manager.Start(context.Background()); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}

func fileUploadHandler(ctx server.Context) error {
    // 文件上传处理 (File upload handling)
    request := ctx.Request()
    if err := request.ParseMultipartForm(32 << 20); err != nil { // 32MB 最大值 (32MB max)
        return ctx.JSON(400, map[string]string{"error": "Failed to parse form"})
    }
    
    file, header, err := request.FormFile("upload")
    if err != nil {
        return ctx.JSON(400, map[string]string{"error": "No file uploaded"})
    }
    defer file.Close()
    
    return ctx.JSON(200, map[string]interface{}{
        "filename": header.Filename,
        "size":     header.Size,
        "status":   "uploaded",
    })
}

func setupWebSocket(framework server.WebFramework) {
    // 访问 Fiber 的原生引擎以使用 WebSocket (Access Fiber's native engine for WebSocket)
    fiberApp := framework.GetNativeEngine().(*fiber.App)
    
    fiberApp.Get("/ws", websocket.New(func(c *websocket.Conn) {
        for {
            mt, message, err := c.ReadMessage()
            if err != nil {
                break
            }
            if err := c.WriteMessage(mt, message); err != nil {
                break
            }
        }
    }))
}
```

#### 高性能配置

```go
config.Plugins["fiber"] = map[string]interface{}{
    "prefork":          true,           // 为生产启用 prefork (Enable prefork for production)
    "concurrency":      512 * 1024,     // 增加并发连接 (Increase concurrent connections)
    "compress_body":    true,           // 启用压缩 (Enable compression)
    "etag":            true,            // 启用 ETag (Enable ETag)
    "read_buffer_size": 8192,          // 增加缓冲区大小 (Increase buffer size)
    "write_buffer_size": 8192,         // 增加缓冲区大小 (Increase buffer size)
    "network":          "tcp4",        // 仅使用 IPv4 (Use IPv4 only)
}
```

## 插件管理

### 列出可用插件

```go
// 列出所有注册的插件 (List all registered plugins)
plugins := server.ListFrameworks()
fmt.Println("Available frameworks:", plugins)

// 获取插件信息 (Get plugin information)
for _, name := range plugins {
    info, err := server.GetFrameworkInfo(name)
    if err != nil {
        continue
    }
    fmt.Printf("Plugin: %s v%s - %s\n", 
        info["name"], info["version"], info["description"])
}
```

### 插件发现

```go
// 获取所有插件信息 (Get all plugin information)
allInfo := server.GetAllFrameworkInfo()
for name, info := range allInfo {
    fmt.Printf("Framework: %s\n", name)
    fmt.Printf("  Version: %s\n", info["version"])
    fmt.Printf("  Description: %s\n", info["description"])
}
```

### 默认插件管理

```go
// 获取默认插件 (Get default plugin)
defaultPlugin, err := server.GetDefaultFramework()
if err != nil {
    log.Fatal("No default plugin available")
}
fmt.Println("Default framework:", defaultPlugin.Name())

// 设置默认插件 (Set default plugin)
err = server.SetDefaultFramework("gin")
if err != nil {
    log.Fatal("Failed to set default framework:", err)
}
```

### 插件验证

```go
// 验证插件配置 (Validate plugin configuration)
plugin, err := server.GetFramework("gin")
if err != nil {
    log.Fatal("Plugin not found:", err)
}

config := &server.ServerConfig{
    Framework: "gin",
    Port:     8080,
}

if err := plugin.ValidateConfig(config); err != nil {
    log.Fatal("Invalid configuration:", err)
}
```

## 框架比较

### 性能比较

| 特性          | Gin        | Echo       | Fiber      |
|--------------|------------|------------|------------|
| **请求/秒**   | ~47,000    | ~45,000    | ~50,000    |
| **内存使用**  | 中等       | 低         | 很低       |
| **冷启动**    | 快         | 很快       | 最快       |
| **二进制大小** | 中等       | 小         | 中等       |
| **CPU 使用**  | 低         | 低         | 很低       |

### 功能比较

| 功能              | Gin | Echo | Fiber |
|------------------|-----|------|-------|
| **路由器**        | ✅   | ✅    | ✅     |
| **中间件**        | ✅   | ✅    | ✅     |
| **JSON 绑定**     | ✅   | ✅    | ✅     |
| **模板引擎**      | ✅   | ✅    | ✅     |
| **静态文件**      | ✅   | ✅    | ✅     |
| **WebSocket**     | ❌   | ✅    | ✅     |
| **HTTP/2**        | ✅   | ✅    | ❌     |
| **自动 TLS**      | ❌   | ✅    | ❌     |
| **Prefork**       | ❌   | ❌    | ✅     |

### 使用案例建议

#### 选择 Gin 当：
- 您需要**成熟的生态系统**和广泛的中间件
- 构建**传统 Web 应用程序**和模板
- 需要 **JSON/XML 绑定**和验证
- 团队有 **Go Web 开发经验**
- 需要**稳定性和社区支持**

#### 选择 Echo 当：
- 构建 **RESTful API** 和微服务
- 需要使用 Let's Encrypt 的**自动 TLS**
- 希望**最小内存占用**
- 偏好**清洁简单的架构**
- 构建**中间件密集型应用程序**

#### 选择 Fiber 当：
- 需要**最大性能**和速度
- 构建**实时应用程序**
- 来自 **Node.js/Express 背景**
- 需要 **WebSocket 支持**
- 构建**高并发服务**

## 自定义插件开发

### 插件接口实现

创建自定义框架插件：

```go
package myframework

import (
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
)

type MyFrameworkPlugin struct {
    name        string
    version     string
    description string
}

func NewPlugin() server.FrameworkPlugin {
    return &MyFrameworkPlugin{
        name:        "myframework",
        version:     "1.0.0",
        description: "Custom framework plugin for lmcc-go-sdk",
    }
}

func (p *MyFrameworkPlugin) Name() string {
    return p.name
}

func (p *MyFrameworkPlugin) Version() string {
    return p.version
}

func (p *MyFrameworkPlugin) Description() string {
    return p.description
}

func (p *MyFrameworkPlugin) DefaultConfig() interface{} {
    config := server.DefaultServerConfig()
    config.Framework = p.name
    return config
}

func (p *MyFrameworkPlugin) CreateFramework(config interface{}, services services.ServiceContainer) (server.WebFramework, error) {
    serverConfig, ok := config.(*server.ServerConfig)
    if !ok {
        return nil, fmt.Errorf("invalid configuration type")
    }
    
    // 创建并返回您的框架适配器 (Create and return your framework adapter)
    return NewMyFrameworkAdapter(serverConfig, services), nil
}

func (p *MyFrameworkPlugin) ValidateConfig(config interface{}) error {
    serverConfig, ok := config.(*server.ServerConfig)
    if !ok {
        return fmt.Errorf("invalid configuration type")
    }
    
    return serverConfig.Validate()
}

func (p *MyFrameworkPlugin) GetConfigSchema() interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "framework": map[string]interface{}{
                "type": "string",
                "enum": []string{p.name},
            },
            // 在这里添加您的配置模式 (Add your configuration schema here)
        },
        "required": []string{"framework"},
    }
}

// 自动注册插件 (Auto-register plugin)
func init() {
    server.RegisterFramework(NewPlugin())
}
```

## 下一步

- **[中间件系统](04_middleware_system.md)** - 了解中间件架构
- **[服务器管理](05_server_management.md)** - 掌握服务器生命周期管理
- **[集成示例](06_integration_examples.md)** - 查看真实世界集成模式
- **[最佳实践](07_best_practices.md)** - 生产部署指南