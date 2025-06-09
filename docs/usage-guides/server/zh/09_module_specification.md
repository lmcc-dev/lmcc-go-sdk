# 模块规范

服务器模块为多种 Web 框架提供统一接口，具有全面的配置管理、中间件支持和插件架构。

## 目录

1. [核心接口](#核心接口)
2. [配置类型](#配置类型)
3. [服务器管理](#服务器管理)
4. [插件系统](#插件系统)
5. [错误类型](#错误类型)
6. [常量](#常量)

## 核心接口

### WebFramework

提供框架无关服务器操作的统一 Web 框架接口。

```go
type WebFramework interface {
    // Start 启动服务器 (Start starts the server)
    Start(ctx context.Context) error
    
    // Stop 停止服务器 (Stop stops the server)
    Stop(ctx context.Context) error
    
    // RegisterRoute 注册路由 (RegisterRoute registers a route)
    RegisterRoute(method, path string, handler Handler) error
    
    // RegisterMiddleware 注册全局中间件 (RegisterMiddleware registers global middleware)
    RegisterMiddleware(middleware Middleware) error
    
    // Group 创建路由组 (Group creates a route group)
    Group(prefix string, middlewares ...Middleware) RouteGroup
    
    // GetNativeEngine 返回原生框架实例 (GetNativeEngine returns the native framework instance)
    GetNativeEngine() interface{}
    
    // GetConfig 返回服务器配置 (GetConfig returns the server configuration)
    GetConfig() *ServerConfig
}
```

### Handler

用于请求处理的统一处理器接口。

```go
type Handler interface {
    // Handle 处理请求 (Handle processes the request)
    Handle(ctx Context) error
}

type HandlerFunc func(ctx Context) error

func (f HandlerFunc) Handle(ctx Context) error {
    return f(ctx)
}
```

### Middleware

用于请求/响应处理的统一中间件接口。

```go
type Middleware interface {
    // Process 处理中间件逻辑 (Process handles middleware logic)
    Process(ctx Context, next func() error) error
}

type MiddlewareFunc func(ctx Context, next func() error) error

func (f MiddlewareFunc) Process(ctx Context, next func() error) error {
    return f(ctx, next)
}
```

### Context

提供框架无关请求/响应操作的统一请求上下文接口。

```go
type Context interface {
    // 请求/响应 (Request/Response)
    Request() *http.Request
    Response() http.ResponseWriter
    
    // 参数 (Parameters)
    Param(key string) string
    Query(key string) string
    Header(key string) string
    SetHeader(key, value string)
    
    // 响应方法 (Response methods)
    JSON(code int, obj interface{}) error
    String(code int, format string, values ...interface{}) error
    Data(code int, contentType string, data []byte) error
    
    // 上下文值 (Context values)
    Set(key string, value interface{})
    Get(key string) (interface{}, bool)
    GetString(key string) string
    GetInt(key string) int
    GetBool(key string) bool
    
    // 请求处理 (Request processing)
    Bind(obj interface{}) error
    ClientIP() string
    UserAgent() string
    Method() string
    Path() string
    FullPath() string
}
```

### RouteGroup

用于组织具有共享中间件和前缀的路由的路由组接口。

```go
type RouteGroup interface {
    // RegisterRoute 在组内注册路由 (RegisterRoute registers a route within the group)
    RegisterRoute(method, path string, handler Handler) error
    
    // RegisterMiddleware 注册组级中间件 (RegisterMiddleware registers group-level middleware)
    RegisterMiddleware(middleware Middleware) error
    
    // Group 创建子路由组 (Group creates a sub route group)
    Group(prefix string, middlewares ...Middleware) RouteGroup
}
```

### FrameworkPlugin

用于注册 Web 框架实现的框架插件接口。

```go
type FrameworkPlugin interface {
    // 插件信息 (Plugin information)
    Name() string
    Version() string
    Description() string
    
    // 配置 (Configuration)
    DefaultConfig() interface{}
    ValidateConfig(config interface{}) error
    GetConfigSchema() interface{}
    
    // 框架创建 (Framework creation)
    CreateFramework(config interface{}, services services.ServiceContainer) (WebFramework, error)
}
```

## 配置类型

### ServerConfig

主服务器配置结构。

```go
type ServerConfig struct {
    // 基本设置 (Basic settings)
    Framework        string        `yaml:"framework"`
    Host            string        `yaml:"host"`
    Port            int           `yaml:"port"`
    Mode            string        `yaml:"mode"`
    
    // 超时 (Timeouts)
    ReadTimeout     time.Duration `yaml:"read-timeout"`
    WriteTimeout    time.Duration `yaml:"write-timeout"`
    IdleTimeout     time.Duration `yaml:"idle-timeout"`
    MaxHeaderBytes  int           `yaml:"max-header-bytes"`
    
    // 功能配置 (Feature configurations)
    CORS            CORSConfig            `yaml:"cors"`
    Middleware      MiddlewareConfig      `yaml:"middleware"`
    TLS             TLSConfig             `yaml:"tls"`
    GracefulShutdown GracefulShutdownConfig `yaml:"graceful-shutdown"`
    
    // 插件配置 (Plugin configurations)
    Plugins         map[string]interface{} `yaml:"plugins"`
}
```

**方法:**
- `Validate() error` - 验证配置
- `GetAddress() string` - 返回服务器地址 (host:port)
- `IsDebugMode() bool` - 检查服务器是否处于调试模式
- `IsReleaseMode() bool` - 检查服务器是否处于发布模式
- `IsTestMode() bool` - 检查服务器是否处于测试模式

### CORSConfig

CORS 配置结构。

```go
type CORSConfig struct {
    Enabled          bool          `yaml:"enabled"`
    AllowOrigins     []string      `yaml:"allow-origins"`
    AllowMethods     []string      `yaml:"allow-methods"`
    AllowHeaders     []string      `yaml:"allow-headers"`
    ExposeHeaders    []string      `yaml:"expose-headers"`
    AllowCredentials bool          `yaml:"allow-credentials"`
    MaxAge           time.Duration `yaml:"max-age"`
}
```

### MiddlewareConfig

中间件配置结构。

```go
type MiddlewareConfig struct {
    Logger    LoggerMiddlewareConfig    `yaml:"logger"`
    Recovery  RecoveryMiddlewareConfig  `yaml:"recovery"`
    RateLimit RateLimitMiddlewareConfig `yaml:"rate-limit"`
    Auth      AuthMiddlewareConfig      `yaml:"auth"`
}
```

### LoggerMiddlewareConfig

日志中间件配置。

```go
type LoggerMiddlewareConfig struct {
    Enabled     bool     `yaml:"enabled"`
    SkipPaths   []string `yaml:"skip-paths"`
    Format      string   `yaml:"format"`         // "json" 或 "text"
    IncludeBody bool     `yaml:"include-body"`
    MaxBodySize int      `yaml:"max-body-size"`
}
```

### RecoveryMiddlewareConfig

恢复中间件配置。

```go
type RecoveryMiddlewareConfig struct {
    Enabled             bool `yaml:"enabled"`
    PrintStack          bool `yaml:"print-stack"`
    DisableStackAll     bool `yaml:"disable-stack-all"`
    DisableColorConsole bool `yaml:"disable-color-console"`
}
```

### RateLimitMiddlewareConfig

限流中间件配置。

```go
type RateLimitMiddlewareConfig struct {
    Enabled bool    `yaml:"enabled"`
    Rate    float64 `yaml:"rate"`    // 每秒请求数 (Requests per second)
    Burst   int     `yaml:"burst"`   // 突发请求数 (Burst requests)
    KeyFunc string  `yaml:"key-func"` // "ip", "user", 或 "custom"
}
```

### AuthMiddlewareConfig

身份验证中间件配置。

```go
type AuthMiddlewareConfig struct {
    Enabled   bool     `yaml:"enabled"`
    Type      string   `yaml:"type"`       // "jwt", "basic", 或 "custom"
    SkipPaths []string `yaml:"skip-paths"`
    JWT       JWTConfig `yaml:"jwt"`
}
```

### JWTConfig

JWT 身份验证配置。

```go
type JWTConfig struct {
    Secret         string        `yaml:"secret"`
    Issuer         string        `yaml:"issuer"`
    Audience       string        `yaml:"audience"`
    ExpirationTime time.Duration `yaml:"expiration-time"`
    RefreshTime    time.Duration `yaml:"refresh-time"`
}
```

### TLSConfig

TLS 配置结构。

```go
type TLSConfig struct {
    Enabled  bool     `yaml:"enabled"`
    CertFile string   `yaml:"cert-file"`
    KeyFile  string   `yaml:"key-file"`
    AutoTLS  bool     `yaml:"auto-tls"`
    Domains  []string `yaml:"domains"`
}
```

### GracefulShutdownConfig

优雅关闭配置。

```go
type GracefulShutdownConfig struct {
    Enabled  bool          `yaml:"enabled"`
    Timeout  time.Duration `yaml:"timeout"`
    WaitTime time.Duration `yaml:"wait-time"`
}
```

## 服务器管理

### ServerManager

主服务器管理接口。

```go
type ServerManager interface {
    // 服务器生命周期 (Server lifecycle)
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    Restart(ctx context.Context) error
    
    // 配置 (Configuration)
    GetConfig() *ServerConfig
    UpdateConfig(config *ServerConfig) error
    
    // 状态 (Status)
    IsRunning() bool
    GetStatus() ServerStatus
    
    // 框架访问 (Framework access)
    GetFramework() WebFramework
}
```

### ServerFactory

用于创建服务器实例的工厂。

```go
type ServerFactory interface {
    // Create 创建新的服务器实例 (Create creates a new server instance)
    Create(config *ServerConfig, services services.ServiceContainer) (ServerManager, error)
    
    // CreateWithPlugin 使用特定插件创建服务器 (CreateWithPlugin creates a server using a specific plugin)
    CreateWithPlugin(pluginName string, config interface{}, services services.ServiceContainer) (ServerManager, error)
    
    // GetSupportedFrameworks 返回支持的框架列表 (GetSupportedFrameworks returns list of supported frameworks)
    GetSupportedFrameworks() []string
}
```

**函数:**
- `NewServerFactory() ServerFactory` - 创建新的服务器工厂
- `DefaultServerConfig() *ServerConfig` - 返回默认服务器配置

## 插件系统

### Registry

用于管理框架插件的插件注册表。

```go
type Registry interface {
    // 插件管理 (Plugin management)
    Register(plugin FrameworkPlugin) error
    Unregister(name string) error
    Get(name string) (FrameworkPlugin, error)
    List() []FrameworkPlugin
    
    // 插件查询 (Plugin queries)
    Exists(name string) bool
    GetSupportedFrameworks() []string
}
```

**函数:**
- `NewRegistry() Registry` - 创建新的插件注册表
- `DefaultRegistry() Registry` - 返回默认全局注册表

### 内置插件

模块包含流行 Web 框架的内置插件：

#### Gin 插件
- **名称**: `gin`
- **版本**: 依赖框架版本
- **配置**: Gin 特定设置

#### Echo 插件
- **名称**: `echo`
- **版本**: 依赖框架版本
- **配置**: Echo 特定设置

#### Fiber 插件
- **名称**: `fiber`
- **版本**: 依赖框架版本
- **配置**: Fiber 特定设置

## 错误类型

### 常见错误

```go
var (
    ErrServerNotRunning    = errors.New("server is not running")
    ErrServerAlreadyRunning = errors.New("server is already running")
    ErrInvalidConfig       = errors.New("invalid server configuration")
    ErrPluginNotFound      = errors.New("plugin not found")
    ErrFrameworkNotSupported = errors.New("framework not supported")
    ErrMiddlewareRegistration = errors.New("middleware registration failed")
    ErrRouteRegistration   = errors.New("route registration failed")
)
```

### 错误类别

1. **配置错误**: 无效或缺失的配置值
2. **插件错误**: 插件注册、加载或执行失败
3. **框架错误**: 框架特定的初始化或操作错误
4. **中间件错误**: 中间件注册或执行失败
5. **路由错误**: 路由注册或处理失败

## 常量

### 框架名称

```go
const (
    FrameworkGin   = "gin"
    FrameworkEcho  = "echo"
    FrameworkFiber = "fiber"
)
```

### 服务器模式

```go
const (
    ModeDebug   = "debug"
    ModeRelease = "release"
    ModeTest    = "test"
)
```

### HTTP 方法

```go
const (
    MethodGET     = "GET"
    MethodPOST    = "POST"
    MethodPUT     = "PUT"
    MethodPATCH   = "PATCH"
    MethodDELETE  = "DELETE"
    MethodHEAD    = "HEAD"
    MethodOPTIONS = "OPTIONS"
)
```

### 默认值

```go
const (
    DefaultHost            = "0.0.0.0"
    DefaultPort            = 8080
    DefaultMode            = ModeDebug
    DefaultFramework       = FrameworkGin
    DefaultReadTimeout     = 15 * time.Second
    DefaultWriteTimeout    = 15 * time.Second
    DefaultIdleTimeout     = 60 * time.Second
    DefaultMaxHeaderBytes  = 1 << 20 // 1MB
    DefaultShutdownTimeout = 30 * time.Second
)
```

## 包函数

### 配置

```go
// DefaultServerConfig 返回具有默认值的服务器配置 (DefaultServerConfig returns a server configuration with default values)
func DefaultServerConfig() *ServerConfig

// ValidateServerConfig 验证服务器配置 (ValidateServerConfig validates a server configuration)
func ValidateServerConfig(config *ServerConfig) error
```

### 工厂函数

```go
// NewServerFactory 创建新的服务器工厂 (NewServerFactory creates a new server factory)
func NewServerFactory() ServerFactory

// NewServerManager 使用给定配置创建新的服务器管理器 (NewServerManager creates a new server manager with the given configuration)
func NewServerManager(config *ServerConfig, services services.ServiceContainer) (ServerManager, error)

// QuickStart 创建并启动具有最小配置的服务器 (QuickStart creates and starts a server with minimal configuration)
func QuickStart(framework string, port int) (ServerManager, error)
```

### 注册表函数

```go
// NewRegistry 创建新的插件注册表 (NewRegistry creates a new plugin registry)
func NewRegistry() Registry

// DefaultRegistry 返回默认全局注册表 (DefaultRegistry returns the default global registry)
func DefaultRegistry() Registry

// RegisterPlugin 向默认注册表注册插件 (RegisterPlugin registers a plugin with the default registry)
func RegisterPlugin(plugin FrameworkPlugin) error

// GetPlugin 从默认注册表检索插件 (GetPlugin retrieves a plugin from the default registry)
func GetPlugin(name string) (FrameworkPlugin, error)
```

### 上下文工具

```go
// NewContext 从 HTTP 请求和响应创建新的统一上下文 (NewContext creates a new unified context from an HTTP request and response)
func NewContext(req *http.Request, resp http.ResponseWriter) Context

// WrapContext 将框架特定上下文包装到统一接口中 (WrapContext wraps a framework-specific context into the unified interface)
func WrapContext(frameworkCtx interface{}) (Context, error)
```

## 类型别名

```go
// 方便使用的常见类型别名 (Common type aliases for convenience)
type (
    HTTPHandler = http.HandlerFunc
    HTTPMethod  = string
    HTTPStatus  = int
)
```

## 集成点

### 服务容器集成

服务器模块与服务容器系统集成以进行依赖注入：

```go
type ServiceContainer interface {
    Get(name string) (interface{}, error)
    Set(name string, service interface{}) error
    Has(name string) bool
}
```

### 日志集成

与 SDK 的日志模块集成：

```go
// 来自日志模块的日志接口 (Logger interface from log module)
type Logger interface {
    Debug(msg string, fields ...Field)
    Info(msg string, fields ...Field)
    Warn(msg string, fields ...Field)
    Error(msg string, fields ...Field)
    Fatal(msg string, fields ...Field)
}
```

### 配置集成

与 SDK 的配置模块集成：

```go
// 来自配置模块的配置接口 (Config interface from config module)
type Config interface {
    Get(key string) interface{}
    GetString(key string) string
    GetInt(key string) int
    GetBool(key string) bool
    Unmarshal(v interface{}) error
}
```

## 版本信息

- **模块版本**: v1.0.0
- **Go 版本**: 1.21+
- **支持的框架**:
  - Gin: v1.9.0+
  - Echo: v4.10.0+
  - Fiber: v2.48.0+

## 迁移指南

### 从 v0.x 到 v1.0

1. **配置结构更改**: 更新配置文件以使用新结构
2. **接口更新**: 更新自定义中间件和处理器以使用新接口
3. **插件注册**: 更新插件注册调用
4. **错误处理**: 更新错误处理以使用新错误类型

有关详细的迁移说明，请参阅 [迁移指南](../migration.md)。

## 示例用法

### 基本服务器创建

```go
package main

import (
    "context"
    "log"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin"
)

func main() {
    // 创建默认配置 (Create default configuration)
    config := server.DefaultServerConfig()
    config.Framework = "gin"
    config.Port = 8080
    
    // 创建服务器管理器 (Create server manager)
    manager, err := server.CreateServerManager("gin", config)
    if err != nil {
        log.Fatal("Failed to create server:", err)
    }
    
    // 获取框架并注册路由 (Get framework and register routes)
    framework := manager.GetFramework()
    framework.RegisterRoute("GET", "/health", server.HandlerFunc(func(ctx server.Context) error {
        return ctx.JSON(200, map[string]string{"status": "healthy"})
    }))
    
    // 启动服务器 (Start server)
    if err := manager.Start(context.Background()); err != nil {
        log.Fatal("Server failed:", err)
    }
}
```

### 中间件注册

```go
// 注册全局中间件 (Register global middleware)
framework.RegisterMiddleware(middleware.NewLogger())
framework.RegisterMiddleware(middleware.NewRecovery())

// 创建路由组并添加组特定中间件 (Create route group and add group-specific middleware)
api := framework.Group("/api/v1")
api.RegisterMiddleware(middleware.NewAuth())

// 在组中注册路由 (Register routes in group)
api.RegisterRoute("GET", "/users", userHandler.GetUsers)
api.RegisterRoute("POST", "/users", userHandler.CreateUser)
```

### 自定义插件

```go
type MyFrameworkPlugin struct{}

func (p *MyFrameworkPlugin) Name() string {
    return "myframework"
}

func (p *MyFrameworkPlugin) Version() string {
    return "1.0.0"
}

func (p *MyFrameworkPlugin) Description() string {
    return "Custom framework plugin"
}

func (p *MyFrameworkPlugin) DefaultConfig() interface{} {
    return &MyFrameworkConfig{
        Debug: false,
        Port:  8080,
    }
}

func (p *MyFrameworkPlugin) CreateFramework(config interface{}, services services.ServiceContainer) (server.WebFramework, error) {
    cfg := config.(*MyFrameworkConfig)
    return NewMyFramework(cfg), nil
}

// 注册自定义插件 (Register custom plugin)
err := server.RegisterPlugin(&MyFrameworkPlugin{})
if err != nil {
    log.Fatal("Failed to register plugin:", err)
}
```

## 下一步

- **[快速开始](01_quick_start.md)** - 开始使用服务器模块
- **[配置指南](02_configuration.md)** - 详细配置选项
- **[最佳实践](07_best_practices.md)** - 生产部署指南