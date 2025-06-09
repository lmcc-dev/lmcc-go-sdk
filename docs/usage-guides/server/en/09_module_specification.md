# Module Specification

The Server module provides a unified interface for multiple web frameworks with comprehensive configuration management, middleware support, and plugin architecture.

## Table of Contents

1. [Core Interfaces](#core-interfaces)
2. [Configuration Types](#configuration-types)
3. [Server Management](#server-management)
4. [Plugin System](#plugin-system)
5. [Error Types](#error-types)
6. [Constants](#constants)

## Core Interfaces

### WebFramework

The unified web framework interface providing framework-agnostic server operations.

```go
type WebFramework interface {
    // Start starts the server
    Start(ctx context.Context) error
    
    // Stop stops the server
    Stop(ctx context.Context) error
    
    // RegisterRoute registers a route
    RegisterRoute(method, path string, handler Handler) error
    
    // RegisterMiddleware registers global middleware
    RegisterMiddleware(middleware Middleware) error
    
    // Group creates a route group
    Group(prefix string, middlewares ...Middleware) RouteGroup
    
    // GetNativeEngine returns the native framework instance
    GetNativeEngine() interface{}
    
    // GetConfig returns the server configuration
    GetConfig() *ServerConfig
}
```

### Handler

Unified handler interface for request processing.

```go
type Handler interface {
    // Handle processes the request
    Handle(ctx Context) error
}

type HandlerFunc func(ctx Context) error

func (f HandlerFunc) Handle(ctx Context) error {
    return f(ctx)
}
```

### Middleware

Unified middleware interface for request/response processing.

```go
type Middleware interface {
    // Process handles middleware logic
    Process(ctx Context, next func() error) error
}

type MiddlewareFunc func(ctx Context, next func() error) error

func (f MiddlewareFunc) Process(ctx Context, next func() error) error {
    return f(ctx, next)
}
```

### Context

Unified request context interface providing framework-agnostic request/response operations.

```go
type Context interface {
    // Request/Response
    Request() *http.Request
    Response() http.ResponseWriter
    
    // Parameters
    Param(key string) string
    Query(key string) string
    Header(key string) string
    SetHeader(key, value string)
    
    // Response methods
    JSON(code int, obj interface{}) error
    String(code int, format string, values ...interface{}) error
    Data(code int, contentType string, data []byte) error
    
    // Context values
    Set(key string, value interface{})
    Get(key string) (interface{}, bool)
    GetString(key string) string
    GetInt(key string) int
    GetBool(key string) bool
    
    // Request processing
    Bind(obj interface{}) error
    ClientIP() string
    UserAgent() string
    Method() string
    Path() string
    FullPath() string
}
```

### RouteGroup

Route group interface for organizing routes with shared middleware and prefixes.

```go
type RouteGroup interface {
    // RegisterRoute registers a route within the group
    RegisterRoute(method, path string, handler Handler) error
    
    // RegisterMiddleware registers group-level middleware
    RegisterMiddleware(middleware Middleware) error
    
    // Group creates a sub route group
    Group(prefix string, middlewares ...Middleware) RouteGroup
}
```

### FrameworkPlugin

Framework plugin interface for registering web framework implementations.

```go
type FrameworkPlugin interface {
    // Plugin information
    Name() string
    Version() string
    Description() string
    
    // Configuration
    DefaultConfig() interface{}
    ValidateConfig(config interface{}) error
    GetConfigSchema() interface{}
    
    // Framework creation
    CreateFramework(config interface{}, services services.ServiceContainer) (WebFramework, error)
}
```

## Configuration Types

### ServerConfig

Main server configuration structure.

```go
type ServerConfig struct {
    // Basic settings
    Framework        string        `yaml:"framework"`
    Host            string        `yaml:"host"`
    Port            int           `yaml:"port"`
    Mode            string        `yaml:"mode"`
    
    // Timeouts
    ReadTimeout     time.Duration `yaml:"read-timeout"`
    WriteTimeout    time.Duration `yaml:"write-timeout"`
    IdleTimeout     time.Duration `yaml:"idle-timeout"`
    MaxHeaderBytes  int           `yaml:"max-header-bytes"`
    
    // Feature configurations
    CORS            CORSConfig            `yaml:"cors"`
    Middleware      MiddlewareConfig      `yaml:"middleware"`
    TLS             TLSConfig             `yaml:"tls"`
    GracefulShutdown GracefulShutdownConfig `yaml:"graceful-shutdown"`
    
    // Plugin configurations
    Plugins         map[string]interface{} `yaml:"plugins"`
}
```

**Methods:**
- `Validate() error` - Validates the configuration
- `GetAddress() string` - Returns the server address (host:port)
- `IsDebugMode() bool` - Checks if server is in debug mode
- `IsReleaseMode() bool` - Checks if server is in release mode
- `IsTestMode() bool` - Checks if server is in test mode

### CORSConfig

CORS configuration structure.

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

Middleware configuration structure.

```go
type MiddlewareConfig struct {
    Logger    LoggerMiddlewareConfig    `yaml:"logger"`
    Recovery  RecoveryMiddlewareConfig  `yaml:"recovery"`
    RateLimit RateLimitMiddlewareConfig `yaml:"rate-limit"`
    Auth      AuthMiddlewareConfig      `yaml:"auth"`
}
```

### LoggerMiddlewareConfig

Logger middleware configuration.

```go
type LoggerMiddlewareConfig struct {
    Enabled     bool     `yaml:"enabled"`
    SkipPaths   []string `yaml:"skip-paths"`
    Format      string   `yaml:"format"`         // "json" or "text"
    IncludeBody bool     `yaml:"include-body"`
    MaxBodySize int      `yaml:"max-body-size"`
}
```

### RecoveryMiddlewareConfig

Recovery middleware configuration.

```go
type RecoveryMiddlewareConfig struct {
    Enabled             bool `yaml:"enabled"`
    PrintStack          bool `yaml:"print-stack"`
    DisableStackAll     bool `yaml:"disable-stack-all"`
    DisableColorConsole bool `yaml:"disable-color-console"`
}
```

### RateLimitMiddlewareConfig

Rate limiting middleware configuration.

```go
type RateLimitMiddlewareConfig struct {
    Enabled bool    `yaml:"enabled"`
    Rate    float64 `yaml:"rate"`    // Requests per second
    Burst   int     `yaml:"burst"`   // Burst requests
    KeyFunc string  `yaml:"key-func"` // "ip", "user", or "custom"
}
```

### AuthMiddlewareConfig

Authentication middleware configuration.

```go
type AuthMiddlewareConfig struct {
    Enabled   bool     `yaml:"enabled"`
    Type      string   `yaml:"type"`       // "jwt", "basic", or "custom"
    SkipPaths []string `yaml:"skip-paths"`
    JWT       JWTConfig `yaml:"jwt"`
}
```

### JWTConfig

JWT authentication configuration.

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

TLS configuration structure.

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

Graceful shutdown configuration.

```go
type GracefulShutdownConfig struct {
    Enabled  bool          `yaml:"enabled"`
    Timeout  time.Duration `yaml:"timeout"`
    WaitTime time.Duration `yaml:"wait-time"`
}
```

## Server Management

### ServerManager

Main server management interface.

```go
type ServerManager interface {
    // Server lifecycle
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    Restart(ctx context.Context) error
    
    // Configuration
    GetConfig() *ServerConfig
    UpdateConfig(config *ServerConfig) error
    
    // Status
    IsRunning() bool
    GetStatus() ServerStatus
    
    // Framework access
    GetFramework() WebFramework
}
```

### ServerFactory

Factory for creating server instances.

```go
type ServerFactory interface {
    // Create creates a new server instance
    Create(config *ServerConfig, services services.ServiceContainer) (ServerManager, error)
    
    // CreateWithPlugin creates a server using a specific plugin
    CreateWithPlugin(pluginName string, config interface{}, services services.ServiceContainer) (ServerManager, error)
    
    // GetSupportedFrameworks returns list of supported frameworks
    GetSupportedFrameworks() []string
}
```

**Functions:**
- `NewServerFactory() ServerFactory` - Creates a new server factory
- `DefaultServerConfig() *ServerConfig` - Returns default server configuration

## Plugin System

### Registry

Plugin registry for managing framework plugins.

```go
type Registry interface {
    // Plugin management
    Register(plugin FrameworkPlugin) error
    Unregister(name string) error
    Get(name string) (FrameworkPlugin, error)
    List() []FrameworkPlugin
    
    // Plugin queries
    Exists(name string) bool
    GetSupportedFrameworks() []string
}
```

**Functions:**
- `NewRegistry() Registry` - Creates a new plugin registry
- `DefaultRegistry() Registry` - Returns the default global registry

### Built-in Plugins

The module includes built-in plugins for popular web frameworks:

#### Gin Plugin
- **Name**: `gin`
- **Version**: Framework version dependent
- **Configuration**: Gin-specific settings

#### Echo Plugin
- **Name**: `echo`
- **Version**: Framework version dependent
- **Configuration**: Echo-specific settings

#### Fiber Plugin
- **Name**: `fiber`
- **Version**: Framework version dependent
- **Configuration**: Fiber-specific settings

## Error Types

### Common Errors

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

### Error Categories

1. **Configuration Errors**: Invalid or missing configuration values
2. **Plugin Errors**: Plugin registration, loading, or execution failures
3. **Framework Errors**: Framework-specific initialization or operation errors
4. **Middleware Errors**: Middleware registration or execution failures
5. **Route Errors**: Route registration or handling failures

## Constants

### Framework Names

```go
const (
    FrameworkGin   = "gin"
    FrameworkEcho  = "echo"
    FrameworkFiber = "fiber"
)
```

### Server Modes

```go
const (
    ModeDebug   = "debug"
    ModeRelease = "release"
    ModeTest    = "test"
)
```

### HTTP Methods

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

### Default Values

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

## Package Functions

### Configuration

```go
// DefaultServerConfig returns a server configuration with default values
func DefaultServerConfig() *ServerConfig

// ValidateServerConfig validates a server configuration
func ValidateServerConfig(config *ServerConfig) error
```

### Factory Functions

```go
// NewServerFactory creates a new server factory
func NewServerFactory() ServerFactory

// NewServerManager creates a new server manager with the given configuration
func NewServerManager(config *ServerConfig, services services.ServiceContainer) (ServerManager, error)

// QuickStart creates and starts a server with minimal configuration
func QuickStart(framework string, port int) (ServerManager, error)
```

### Registry Functions

```go
// NewRegistry creates a new plugin registry
func NewRegistry() Registry

// DefaultRegistry returns the default global registry
func DefaultRegistry() Registry

// RegisterPlugin registers a plugin with the default registry
func RegisterPlugin(plugin FrameworkPlugin) error

// GetPlugin retrieves a plugin from the default registry
func GetPlugin(name string) (FrameworkPlugin, error)
```

### Context Utilities

```go
// NewContext creates a new unified context from an HTTP request and response
func NewContext(req *http.Request, resp http.ResponseWriter) Context

// WrapContext wraps a framework-specific context into the unified interface
func WrapContext(frameworkCtx interface{}) (Context, error)
```

## Type Aliases

```go
// Common type aliases for convenience
type (
    HTTPHandler = http.HandlerFunc
    HTTPMethod  = string
    HTTPStatus  = int
)
```

## Integration Points

### Service Container Integration

The server module integrates with the service container system for dependency injection:

```go
type ServiceContainer interface {
    Get(name string) (interface{}, error)
    Set(name string, service interface{}) error
    Has(name string) bool
}
```

### Logger Integration

Integration with the SDK's logging module:

```go
// Logger interface from log module
type Logger interface {
    Debug(msg string, fields ...Field)
    Info(msg string, fields ...Field)
    Warn(msg string, fields ...Field)
    Error(msg string, fields ...Field)
    Fatal(msg string, fields ...Field)
}
```

### Configuration Integration

Integration with the SDK's configuration module:

```go
// Config interface from config module
type Config interface {
    Get(key string) interface{}
    GetString(key string) string
    GetInt(key string) int
    GetBool(key string) bool
    Unmarshal(v interface{}) error
}
```

## Version Information

- **Module Version**: v1.0.0
- **Go Version**: 1.21+
- **Supported Frameworks**:
  - Gin: v1.9.0+
  - Echo: v4.10.0+
  - Fiber: v2.48.0+

## Migration Guide

### From v0.x to v1.0

1. **Configuration Structure Changes**: Update configuration files to use new structure
2. **Interface Updates**: Update custom middleware and handlers to use new interfaces
3. **Plugin Registration**: Update plugin registration calls
4. **Error Handling**: Update error handling to use new error types

For detailed migration instructions, see the [Migration Guide](../migration.md).
