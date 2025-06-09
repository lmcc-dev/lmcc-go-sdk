# 配置指南

本指南涵盖了 lmcc-go-sdk 服务器模块中所有可用的配置选项。服务器模块为 Web 服务器设置、中间件管理和框架特定自定义提供了全面的配置功能。

## 概述

服务器模块使用 `ServerConfig` 结构来定义所有配置选项。配置可以从 YAML 文件、环境变量加载，或通过编程方式设置。该模块与 `pkg/config` 模块无缝集成，提供高级配置管理。

## 基础配置

### 默认配置

最简单的入门方法是使用默认配置：

```go
package main

import (
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
)

func main() {
    // 获取默认配置 (Get default configuration)
    config := server.DefaultServerConfig()
    
    // 根据需要自定义 (Customize as needed)
    config.Framework = "gin"
    config.Port = 8080
    config.Mode = "production"
    
    // 使用配置创建服务器 (Create server with configuration)
    manager, err := server.CreateServerManager("gin", config)
    if err != nil {
        panic(err)
    }
    
    // 启动服务器 (Start server)
    // ... rest of the application
}
```

### 完整配置结构

以下是包含所有可用选项的完整结构：

```go
type ServerConfig struct {
    Framework        string                    // 框架: gin, echo, fiber (Framework)
    Host            string                    // 监听主机 (Listen host)
    Port            int                       // 监听端口 (Listen port)
    Mode            string                    // 运行模式: debug, release, test (Running mode)
    ReadTimeout     time.Duration             // HTTP 读取超时 (HTTP read timeout)
    WriteTimeout    time.Duration             // HTTP 写入超时 (HTTP write timeout)
    IdleTimeout     time.Duration             // HTTP 空闲超时 (HTTP idle timeout)
    MaxHeaderBytes  int                       // 最大请求头大小 (Maximum request header size)
    CORS            CORSConfig               // CORS 配置 (CORS configuration)
    Middleware      MiddlewareConfig         // 中间件设置 (Middleware settings)
    TLS             TLSConfig                // TLS/HTTPS 设置 (TLS/HTTPS settings)
    Plugins         map[string]interface{}   // 插件特定配置 (Plugin-specific config)
    GracefulShutdown GracefulShutdownConfig  // 优雅关闭设置 (Graceful shutdown settings)
}
```

## 核心服务器设置

### 框架选择

从支持的 Web 框架中选择：

```go
config := server.DefaultServerConfig()

// Gin 框架 (高性能，丰富生态) (Gin framework - high performance, rich ecosystem)
config.Framework = "gin"

// Echo 框架 (极简，快速) (Echo framework - minimalist, fast)
config.Framework = "echo"

// Fiber 框架 (Express 风格，超快) (Fiber framework - Express-inspired, ultra-fast)
config.Framework = "fiber"
```

### 网络配置

配置主机和端口设置：

```go
config := &server.ServerConfig{
    Framework: "gin",
    Host:     "0.0.0.0",    // 监听所有接口 (Listen on all interfaces)
    Port:     8080,         // 标准 HTTP 端口 (Standard HTTP port)
    Mode:     "production", // 运行模式 (Running mode)
}

// 替代配置 (Alternative configurations)
config.Host = "127.0.0.1"  // 仅本地主机 (Localhost only)
config.Host = "localhost"   // 与 127.0.0.1 相同 (Same as 127.0.0.1)
config.Port = 3000         // 替代端口 (Alternative port)
```

### 运行模式

服务器支持三种运行模式：

```go
// 调试模式 - 详细日志，热重载 (Debug mode - detailed logging, hot reload)
config.Mode = "debug"

// 发布模式 - 为生产优化 (Release mode - optimized for production)
config.Mode = "release"

// 测试模式 - 最小输出，快速启动 (Test mode - minimal output, fast startup)
config.Mode = "test"

// 程序化检查模式 (Check mode programmatically)
if config.IsDebugMode() {
    // 调试特定逻辑 (Debug-specific logic)
}
```

### 超时配置

配置各种超时设置：

```go
config := &server.ServerConfig{
    Framework:      "gin",
    ReadTimeout:    30 * time.Second,  // 读取请求时间 (Time to read request)
    WriteTimeout:   30 * time.Second,  // 写入响应时间 (Time to write response)
    IdleTimeout:    60 * time.Second,  // Keep-alive 超时 (Keep-alive timeout)
    MaxHeaderBytes: 1 << 20,           // 1MB 头部限制 (1MB header limit)
}

// 自定义超时示例 (Custom timeout example)
config.ReadTimeout = 5 * time.Minute   // 用于大文件上传 (For large uploads)
config.WriteTimeout = 2 * time.Minute  // 用于流式响应 (For streaming responses)
```

## CORS 配置

配置跨域资源共享 (CORS) 设置：

```go
config := server.DefaultServerConfig()

// 基本 CORS 设置 (Basic CORS setup)
config.CORS = server.CORSConfig{
    Enabled:          true,
    AllowOrigins:     []string{"https://example.com", "https://app.example.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
    ExposeHeaders:    []string{"Content-Length", "Content-Type"},
    AllowCredentials: true,
    MaxAge:          12 * time.Hour,
}

// 开发 CORS（允许所有） (Development CORS - allow all)
config.CORS = server.CORSConfig{
    Enabled:      true,
    AllowOrigins: []string{"*"},
    AllowMethods: []string{"*"},
    AllowHeaders: []string{"*"},
}

// 禁用 CORS (Disable CORS)
config.CORS.Enabled = false
```

## 中间件配置

配置内置中间件组件：

### 日志中间件

```go
config.Middleware.Logger = server.LoggerMiddlewareConfig{
    Enabled:     true,
    Format:      "json",        // json 或 text (json or text)
    SkipPaths:   []string{"/health", "/metrics"},
    IncludeBody: false,         // 在日志中包含请求体 (Include request body in logs)
    MaxBodySize: 1024,          // 记录的最大请求体大小（字节） (Max body size to log in bytes)
}

// 自定义日志配置 (Custom logger configuration)
config.Middleware.Logger = server.LoggerMiddlewareConfig{
    Enabled:     true,
    Format:      "text",
    SkipPaths:   []string{"/favicon.ico", "/static/*"},
    IncludeBody: true,
    MaxBodySize: 2048,
}
```

### 恢复中间件

```go
config.Middleware.Recovery = server.RecoveryMiddlewareConfig{
    Enabled:             true,
    PrintStack:          true,   // 在 panic 时打印堆栈跟踪 (Print stack trace on panic)
    DisableStackAll:     false,  // 禁用所有堆栈跟踪 (Disable all stack traces)
    DisableColorConsole: false,  // 禁用彩色输出 (Disable colored output)
}

// 生产恢复设置 (Production recovery settings)
config.Middleware.Recovery = server.RecoveryMiddlewareConfig{
    Enabled:             true,
    PrintStack:          false,  // 生产中不打印堆栈 (Don't print stack in production)
    DisableStackAll:     true,
    DisableColorConsole: true,
}
```

### 限流中间件

```go
config.Middleware.RateLimit = server.RateLimitMiddlewareConfig{
    Enabled: true,
    Rate:    100.0,      // 每秒 100 个请求 (100 requests per second)
    Burst:   200,        // 允许突发 200 个请求 (Allow burst of 200 requests)
    KeyFunc: "ip",       // 按 IP 地址限流 (Rate limit by IP address)
}

// 替代限流策略 (Alternative rate limiting strategies)
config.Middleware.RateLimit.KeyFunc = "user"   // 按用户 ID 限流 (Rate limit by user ID)
config.Middleware.RateLimit.KeyFunc = "custom" // 自定义键函数 (Custom key function)
```

### 认证中间件

```go
config.Middleware.Auth = server.AuthMiddlewareConfig{
    Enabled:   true,
    Type:      "jwt",
    SkipPaths: []string{"/auth/login", "/auth/register", "/health"},
    JWT: server.JWTConfig{
        Secret:         "your-secret-key",
        Issuer:         "your-app",
        Audience:       "your-users",
        ExpirationTime: 24 * time.Hour,
        RefreshTime:    7 * 24 * time.Hour,
    },
}

// 基本认证 (Basic authentication)
config.Middleware.Auth.Type = "basic"

// 禁用认证 (Disable authentication)
config.Middleware.Auth.Enabled = false
```

## TLS/HTTPS 配置

配置 HTTPS 和 TLS 设置：

### 基本 TLS 设置

```go
config.TLS = server.TLSConfig{
    Enabled:  true,
    CertFile: "/path/to/certificate.crt",
    KeyFile:  "/path/to/private.key",
}
```

### 自动 TLS (Let's Encrypt)

```go
config.TLS = server.TLSConfig{
    Enabled: true,
    AutoTLS: true,
    Domains: []string{"example.com", "www.example.com"},
}
```

### 开发 TLS

```go
// 用于开发的自签名证书 (Self-signed certificate for development)
config.TLS = server.TLSConfig{
    Enabled:  true,
    CertFile: "./dev-cert.pem",
    KeyFile:  "./dev-key.pem",
}
```

## 优雅关闭配置

配置优雅关闭行为：

```go
config.GracefulShutdown = server.GracefulShutdownConfig{
    Enabled:  true,
    Timeout:  30 * time.Second,  // 最大关闭时间 (Maximum shutdown time)
    WaitTime: 5 * time.Second,   // 开始关闭前等待 (Wait before starting shutdown)
}

// 开发快速关闭 (Quick shutdown for development)
config.GracefulShutdown = server.GracefulShutdownConfig{
    Enabled:  true,
    Timeout:  5 * time.Second,
    WaitTime: 1 * time.Second,
}

// 禁用优雅关闭 (Disable graceful shutdown)
config.GracefulShutdown.Enabled = false
```

## 插件特定配置

配置框架特定选项：

### Gin 插件配置

```go
config.Plugins = map[string]interface{}{
    "gin": map[string]interface{}{
        "disable_color":        false,
        "trust_proxies":        []string{"127.0.0.1"},
        "max_multipart_memory": 32 << 20, // 32MB
        "enable_html_render":   true,
    },
}
```

### Echo 插件配置

```go
config.Plugins = map[string]interface{}{
    "echo": map[string]interface{}{
        "hide_banner":      true,
        "hide_port":        false,
        "debug":           config.IsDebugMode(),
        "disable_startup_message": false,
    },
}
```

### Fiber 插件配置

```go
config.Plugins = map[string]interface{}{
    "fiber": map[string]interface{}{
        "prefork":              false,
        "strict_routing":       false,
        "case_sensitive":       false,
        "immutable":           false,
        "unescaped_path":      false,
        "etag":                false,
        "body_limit":          4 * 1024 * 1024, // 4MB
        "read_buffer_size":    4096,
        "write_buffer_size":   4096,
    },
}
```

## 从文件配置

### YAML 配置

创建 `config.yaml` 文件：

```yaml
server:
  framework: gin
  host: 0.0.0.0
  port: 8080
  mode: production
  read-timeout: 30s
  write-timeout: 30s
  idle-timeout: 60s
  max-header-bytes: 1048576
  
  cors:
    enabled: true
    allow-origins:
      - "https://example.com"
      - "https://app.example.com"
    allow-methods:
      - GET
      - POST
      - PUT
      - DELETE
      - OPTIONS
    allow-headers:
      - Origin
      - Content-Type
      - Accept
      - Authorization
    allow-credentials: true
    max-age: 12h
  
  middleware:
    logger:
      enabled: true
      format: json
      skip-paths:
        - /health
        - /metrics
      include-body: false
      max-body-size: 1024
    
    recovery:
      enabled: true
      print-stack: false
      disable-stack-all: true
      disable-color-console: true
    
    rate-limit:
      enabled: true
      rate: 100.0
      burst: 200
      key-func: ip
    
    auth:
      enabled: true
      type: jwt
      skip-paths:
        - /auth/login
        - /auth/register
        - /health
      jwt:
        secret: your-secret-key
        issuer: your-app
        audience: your-users
        expiration-time: 24h
        refresh-time: 168h
  
  tls:
    enabled: false
    cert-file: ""
    key-file: ""
    auto-tls: false
    domains: []
  
  graceful-shutdown:
    enabled: true
    timeout: 30s
    wait-time: 5s
  
  plugins:
    gin:
      disable_color: false
      trust_proxies:
        - 127.0.0.1
      max_multipart_memory: 33554432
      enable_html_render: true
```

使用配置模块加载配置：

```go
package main

import (
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
)

type AppConfig struct {
    Server server.ServerConfig `mapstructure:"server"`
}

func main() {
    var cfg AppConfig
    
    // 从文件加载配置 (Load configuration from file)
    err := config.LoadConfig(&cfg,
        config.WithConfigFile("config.yaml", ""),
        config.WithHotReload(true),
    )
    if err != nil {
        panic(err)
    }
    
    // 使用加载的配置创建服务器 (Create server with loaded configuration)
    manager, err := server.CreateServerManager(cfg.Server.Framework, &cfg.Server)
    if err != nil {
        panic(err)
    }
    
    // 启动服务器 (Start server)
    // ... rest of application
}
```

### 环境变量

配置也可以从环境变量加载：

```bash
# 服务器设置 (Server settings)
SERVER_FRAMEWORK=gin
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_MODE=production

# CORS 设置 (CORS settings)
SERVER_CORS_ENABLED=true
SERVER_CORS_ALLOW_ORIGINS=https://example.com,https://app.example.com

# 中间件设置 (Middleware settings)
SERVER_MIDDLEWARE_LOGGER_ENABLED=true
SERVER_MIDDLEWARE_LOGGER_FORMAT=json
SERVER_MIDDLEWARE_RECOVERY_ENABLED=true

# TLS 设置 (TLS settings)
SERVER_TLS_ENABLED=false
SERVER_TLS_CERT_FILE=/path/to/cert.pem
SERVER_TLS_KEY_FILE=/path/to/key.pem
```

从环境加载：

```go
func main() {
    var cfg AppConfig
    
    err := config.LoadConfig(&cfg,
        config.WithEnvPrefix("MYAPP"),
        config.WithEnvKeyReplacer(".", "_"),
    )
    if err != nil {
        panic(err)
    }
    
    // 使用加载的配置 (Use loaded configuration)
    manager, err := server.CreateServerManager(cfg.Server.Framework, &cfg.Server)
    // ...
}
```

## 配置验证

服务器模块提供内置配置验证：

```go
config := &server.ServerConfig{
    Framework: "gin",
    Host:     "localhost",
    Port:     8080,
    Mode:     "production",
}

// 验证配置 (Validate configuration)
if err := config.Validate(); err != nil {
    panic(fmt.Sprintf("Invalid configuration: %v", err))
}

// 创建服务器时自动验证配置 (Configuration is automatically validated when creating server)
manager, err := server.CreateServerManager("gin", config)
if err != nil {
    // 处理验证错误 (Handle validation errors)
    panic(err)
}
```

## 配置最佳实践

### 1. 环境特定配置

为不同环境使用不同配置：

```go
// config/development.yaml
server:
  mode: debug
  cors:
    allow-origins: ["*"]
  middleware:
    logger:
      format: text
      include-body: true

// config/production.yaml
server:
  mode: production
  cors:
    allow-origins: ["https://yourdomain.com"]
  middleware:
    logger:
      format: json
      include-body: false
  tls:
    enabled: true
    auto-tls: true
    domains: ["yourdomain.com"]
```

### 2. 安全配置

将敏感信息保存在环境变量中：

```yaml
# config.yaml
server:
  framework: gin
  middleware:
    auth:
      jwt:
        secret: ${JWT_SECRET}  # 来自环境变量 (From environment)
        issuer: ${JWT_ISSUER}
  tls:
    cert-file: ${TLS_CERT_FILE}
    key-file: ${TLS_KEY_FILE}
```

### 3. 配置验证

始终提前验证配置：

```go
func main() {
    config := loadConfiguration()
    
    // 提前验证 (Validate early)
    if err := config.Validate(); err != nil {
        log.Fatalf("Invalid configuration: %v", err)
    }
    
    // 额外验证 (Additional validation)
    if err := validateBusinessRules(config); err != nil {
        log.Fatalf("Configuration validation failed: %v", err)
    }
    
    // 现在可以安全使用配置 (Now safe to use configuration)
    startServer(config)
}
```

### 4. 默认回退

提供合理的默认值：

```go
func loadConfiguration() *server.ServerConfig {
    config := server.DefaultServerConfig()
    
    // 如果可用，使用文件配置覆盖 (Override with file config if available)
    if err := loadFromFile(config, "config.yaml"); err != nil {
        log.Printf("Config file not found, using defaults: %v", err)
    }
    
    // 使用环境变量覆盖 (Override with environment variables)
    loadFromEnv(config)
    
    return config
}
```

## 配置示例

### 微服务配置

```go
config := &server.ServerConfig{
    Framework: "gin",
    Host:     "0.0.0.0",
    Port:     8080,
    Mode:     "production",
    
    Middleware: server.MiddlewareConfig{
        Logger: server.LoggerMiddlewareConfig{
            Enabled: true,
            Format:  "json",
            SkipPaths: []string{"/health", "/metrics"},
        },
        Recovery: server.RecoveryMiddlewareConfig{
            Enabled:             true,
            PrintStack:          false,
            DisableColorConsole: true,
        },
        Auth: server.AuthMiddlewareConfig{
            Enabled:   true,
            Type:      "jwt",
            SkipPaths: []string{"/health", "/metrics"},
        },
    },
    
    GracefulShutdown: server.GracefulShutdownConfig{
        Enabled: true,
        Timeout: 30 * time.Second,
    },
}
```

### API 网关配置

```go
config := &server.ServerConfig{
    Framework: "fiber",
    Host:     "0.0.0.0",
    Port:     8080,
    Mode:     "production",
    
    CORS: server.CORSConfig{
        Enabled:          true,
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"*"},
        AllowHeaders:     []string{"*"},
        AllowCredentials: false,
    },
    
    Middleware: server.MiddlewareConfig{
        Logger: server.LoggerMiddlewareConfig{
            Enabled: true,
            Format:  "json",
        },
        RateLimit: server.RateLimitMiddlewareConfig{
            Enabled: true,
            Rate:    1000,
            Burst:   2000,
            KeyFunc: "ip",
        },
    },
    
    Plugins: map[string]interface{}{
        "fiber": map[string]interface{}{
            "prefork":       true,
            "read_timeout":  "30s",
            "write_timeout": "30s",
        },
    },
}
```

### 开发配置

```go
config := &server.ServerConfig{
    Framework: "gin",
    Host:     "localhost",
    Port:     3000,
    Mode:     "debug",
    
    CORS: server.CORSConfig{
        Enabled:      true,
        AllowOrigins: []string{"*"},
        AllowMethods: []string{"*"},
        AllowHeaders: []string{"*"},
    },
    
    Middleware: server.MiddlewareConfig{
        Logger: server.LoggerMiddlewareConfig{
            Enabled:     true,
            Format:      "text",
            IncludeBody: true,
        },
        Recovery: server.RecoveryMiddlewareConfig{
            Enabled:    true,
            PrintStack: true,
        },
        Auth: server.AuthMiddlewareConfig{
            Enabled: false, // 开发时禁用 (Disabled for development)
        },
    },
    
    GracefulShutdown: server.GracefulShutdownConfig{
        Enabled: false, // 开发快速关闭 (Quick shutdown for development)
    },
}
```

## 配置参考

有关所有配置选项的完整参考，请参阅：

- **[框架插件](03_framework_plugins.md)** - 框架特定配置
- **[中间件系统](04_middleware_system.md)** - 中间件配置详情
- **[服务器管理](05_server_management.md)** - 高级服务器配置
- **[模块规范](09_module_specification.md)** - 完整 API 参考

## 故障排除

### 常见配置问题

1. **无效端口**: 确保端口在 1-65535 之间
2. **找不到框架**: 验证框架插件已导入
3. **TLS 证书问题**: 检查证书路径和权限
4. **CORS 问题**: 验证源模式和允许的方法
5. **中间件冲突**: 检查中间件顺序和兼容性

详细解决方案请参阅 **[故障排除指南](08_troubleshooting.md)**。