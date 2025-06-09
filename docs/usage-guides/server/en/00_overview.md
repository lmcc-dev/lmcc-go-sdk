# Server Module Overview

The lmcc-go-sdk server module provides a unified, plugin-based web framework abstraction that enables developers to build scalable web applications using popular Go frameworks like Gin, Echo, and Fiber. This module is designed with flexibility, performance, and maintainability in mind.

## Design Philosophy

### Framework Agnostic Architecture

The server module follows a **framework-agnostic** approach, allowing developers to:

- **Switch frameworks** without changing application logic
- **Write portable code** that works across different frameworks
- **Leverage framework-specific features** when needed through the native engine access
- **Future-proof applications** against framework evolution

### Plugin-Based System

The module implements a **plugin architecture** that:

- **Separates concerns** between core functionality and framework-specific implementations
- **Enables extensibility** through custom plugin development
- **Maintains consistency** across different framework implementations
- **Simplifies testing** with mock plugins and framework isolation

### Unified Interface Design

All supported frameworks are accessed through **consistent interfaces**:

- **WebFramework**: Core server operations
- **Handler**: Request processing logic
- **Middleware**: Request/response interception
- **Context**: Request/response abstraction
- **RouteGroup**: Route organization

## Core Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                        │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │              Your Application Code                      │ │
│  │    (Routes, Handlers, Business Logic)                  │ │
│  └─────────────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                 Unified Server Interface                    │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────────┐ │
│  │  WebFrame   │ │  Handler    │ │       Context           │ │
│  │  work       │ │  Interface  │ │     Interface           │ │
│  └─────────────┘ └─────────────┘ └─────────────────────────┘ │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────────┐ │
│  │ Middleware  │ │ RouteGroup  │ │    ServerManager        │ │
│  │ Interface   │ │ Interface   │ │     (Lifecycle)         │ │
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
│  │ Gin Engine  │ │Echo Instance│ │    Fiber App            │ │
│  │  (v1.9+)    │ │   (v4.0+)   │ │    (v2.0+)              │ │
│  └─────────────┘ └─────────────┘ └─────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### Component Architecture

#### 1. Server Manager
The `ServerManager` coordinates the server lifecycle:

```go
type ServerManager interface {
    Start(ctx context.Context) error    // Start the server
    Stop(ctx context.Context) error     // Stop the server gracefully
    IsRunning() bool                    // Check if server is running
    GetFramework() WebFramework         // Access framework instance
    GetConfig() *ServerConfig           // Get current configuration
}
```

#### 2. Framework Plugin System
Each framework is implemented as a plugin:

```go
type FrameworkPlugin interface {
    Name() string                       // Plugin identifier
    Version() string                    // Plugin version
    Description() string                // Plugin description
    CreateFramework(config, services) (WebFramework, error)
    ValidateConfig(config) error        // Configuration validation
    GetConfigSchema() interface{}       // Configuration schema
}
```

#### 3. Service Integration
The module integrates with other SDK components:

```go
type ServiceContainer interface {
    GetLogger() Logger                  // Logging service
    GetConfig() ConfigManager           // Configuration service
    GetErrorHandler() ErrorHandler      // Error handling service
    // ... other services
}
```

## Core Interfaces

### WebFramework Interface

The central interface that all framework plugins implement:

```go
type WebFramework interface {
    // Lifecycle management
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    
    // Route registration
    RegisterRoute(method, path string, handler Handler) error
    
    // Middleware management
    RegisterMiddleware(middleware Middleware) error
    
    // Route grouping
    Group(prefix string, middlewares ...Middleware) RouteGroup
    
    // Framework access
    GetNativeEngine() interface{}
    GetConfig() *ServerConfig
}
```

**Key Benefits:**
- **Consistent API** across all frameworks
- **Type-safe operations** with compile-time checks
- **Native access** when framework-specific features are needed
- **Lifecycle management** built into the interface

### Handler Interface

Unified request handling across frameworks:

```go
type Handler interface {
    Handle(ctx Context) error
}

// Function type for convenience
type HandlerFunc func(ctx Context) error
```

**Usage Examples:**
```go
// Simple handler
handler := server.HandlerFunc(func(ctx server.Context) error {
    return ctx.JSON(200, map[string]string{
        "message": "Hello World",
    })
})

// Handler with error handling
handler := server.HandlerFunc(func(ctx server.Context) error {
    data, err := fetchData()
    if err != nil {
        return err  // Will be handled by recovery middleware
    }
    return ctx.JSON(200, data)
})
```

### Context Interface

Framework-agnostic request/response operations:

```go
type Context interface {
    // HTTP operations
    Request() *http.Request
    Response() http.ResponseWriter
    
    // Parameter access
    Param(key string) string           // Path parameters
    Query(key string) string           // Query parameters
    Header(key string) string          // Request headers
    
    // Response methods
    JSON(code int, obj interface{}) error
    String(code int, format string, values ...interface{}) error
    Data(code int, contentType string, data []byte) error
    
    // Context storage
    Set(key string, value interface{})
    Get(key string) (interface{}, bool)
    GetString(key string) string
    GetInt(key string) int
    GetBool(key string) bool
    
    // Request binding
    Bind(obj interface{}) error
    
    // Utilities
    ClientIP() string
    UserAgent() string
    Method() string
    Path() string
}
```

### Middleware Interface

Unified middleware processing:

```go
type Middleware interface {
    Process(ctx Context, next func() error) error
}

// Function type for convenience
type MiddlewareFunc func(ctx Context, next func() error) error
```

**Middleware Chain Example:**
```go
// Logger middleware
logger := server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
    start := time.Now()
    err := next()  // Call next middleware/handler
    duration := time.Since(start)
    log.Printf("%s %s - %v", ctx.Method(), ctx.Path(), duration)
    return err
})

// Authentication middleware
auth := server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
    token := ctx.Header("Authorization")
    if !isValidToken(token) {
        return ctx.JSON(401, map[string]string{"error": "Unauthorized"})
    }
    return next()
})
```

## Plugin System

### Framework Plugins

The module currently supports three major Go web frameworks:

#### Gin Plugin
- **Performance**: High-performance HTTP framework
- **Ecosystem**: Rich middleware ecosystem
- **Features**: Built-in rendering, binding, validation
- **Best For**: High-throughput APIs, traditional web applications

#### Echo Plugin
- **Design**: Minimalist and fast
- **Architecture**: Clean and extensible
- **Features**: Automatic TLS, HTTP/2 support
- **Best For**: RESTful APIs, microservices

#### Fiber Plugin
- **Performance**: Ultra-fast (inspired by Express.js)
- **Memory**: Low memory footprint
- **Features**: Built-in rate limiting, caching
- **Best For**: Real-time applications, high-concurrency services

### Plugin Registration

Plugins are automatically registered during import:

```go
import (
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin"
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/echo"
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/fiber"
)
```

Or manually registered:

```go
factory := server.NewServerFactory()
err := factory.RegisterPlugin(myCustomPlugin)
```

### Custom Plugin Development

Create custom framework adapters:

```go
type MyFrameworkPlugin struct{}

func (p *MyFrameworkPlugin) Name() string {
    return "myframework"
}

func (p *MyFrameworkPlugin) CreateFramework(config interface{}, services services.ServiceContainer) (server.WebFramework, error) {
    // Implement framework adapter
    return &MyFrameworkAdapter{}, nil
}

// Register plugin
server.RegisterPlugin(&MyFrameworkPlugin{})
```

## Service Integration

### Logger Integration

```go
import "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"

// Logger is automatically available in handlers
handler := server.HandlerFunc(func(ctx server.Context) error {
    log.Info("Processing request", "path", ctx.Path())
    return ctx.JSON(200, response)
})
```

### Configuration Integration

```go
import "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"

type AppConfig struct {
    Server server.ServerConfig `mapstructure:"server"`
}

var cfg AppConfig
config.LoadConfig(&cfg, config.WithConfigFile("config.yaml", ""))
manager, _ := server.CreateServerManager(cfg.Server.Framework, &cfg.Server)
```

### Error Handling Integration

```go
import "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"

handler := server.HandlerFunc(func(ctx server.Context) error {
    if someCondition {
        return errors.New("BUSINESS_ERROR", "Something went wrong")
    }
    return ctx.JSON(200, response)
})
```

## Middleware System

### Built-in Middleware

The server module provides several built-in middleware components:

1. **Logger Middleware**: Request/response logging
2. **Recovery Middleware**: Panic recovery and error handling
3. **CORS Middleware**: Cross-origin resource sharing
4. **Rate Limiting Middleware**: Request rate control
5. **Authentication Middleware**: JWT and basic auth
6. **Compression Middleware**: Response compression
7. **Static File Middleware**: Static asset serving
8. **Security Middleware**: Security headers and protection
9. **Metrics Middleware**: Performance monitoring

### Middleware Chain Processing

```
Request
   ↓
┌─────────────────┐
│ Logger          │
├─────────────────┤
│ Recovery        │
├─────────────────┤
│ CORS            │
├─────────────────┤
│ Authentication  │
├─────────────────┤
│ Rate Limiting   │
├─────────────────┤
│ Your Handler    │
└─────────────────┘
   ↓
Response
```

## Configuration Management

### Configuration Hierarchy

1. **Default Configuration**: Sensible defaults for all options
2. **File Configuration**: YAML/JSON configuration files
3. **Environment Variables**: Runtime configuration override
4. **Programmatic Configuration**: Direct code configuration

### Configuration Validation

```go
config := server.DefaultServerConfig()
config.Framework = "gin"
config.Port = 8080

// Automatic validation
if err := config.Validate(); err != nil {
    log.Fatal("Invalid configuration:", err)
}
```

### Hot Reload Support

```go
configManager, _ := config.LoadConfig(&cfg,
    config.WithConfigFile("config.yaml", ""),
    config.WithHotReload(true),
)

// Configuration changes are automatically applied
```

## Performance Characteristics

### Framework Performance Comparison

| Framework | Requests/sec | Memory Usage | Cold Start | Ecosystem |
|-----------|-------------|--------------|------------|-----------|
| **Gin**   | ~47,000     | Moderate     | Fast       | Rich      |
| **Echo**  | ~45,000     | Low          | Very Fast  | Growing   |
| **Fiber** | ~50,000     | Very Low     | Fastest    | Modern    |

*Benchmarks are approximate and depend on use case and configuration*

### Optimization Features

- **Zero-allocation routing** (framework dependent)
- **Efficient middleware chains** with minimal overhead
- **Connection pooling** and keep-alive support
- **Graceful shutdown** with connection draining
- **Memory-efficient context handling**

## Testing Support

### Mock Framework

```go
type MockFramework struct {
    mock.Mock
}

func (m *MockFramework) Start(ctx context.Context) error {
    args := m.Called(ctx)
    return args.Error(0)
}

// Use in tests
mockFramework := &MockFramework{}
manager := server.NewServerManager(mockFramework, config)
```

### Test Utilities

```go
// Test server lifecycle
func TestServerStartStop(t *testing.T) {
    config := server.DefaultServerConfig()
    config.Mode = "test"
    
    manager, err := server.CreateServerManager("gin", config)
    assert.NoError(t, err)
    
    // Test lifecycle
    err = manager.Start(context.Background())
    assert.NoError(t, err)
    
    err = manager.Stop(context.Background())
    assert.NoError(t, err)
}
```

## Security Features

### Built-in Security

- **TLS/HTTPS support** with automatic certificate management
- **CORS protection** with configurable policies
- **Rate limiting** to prevent abuse
- **Security headers** (HSTS, CSP, etc.)
- **Input validation** and sanitization
- **Authentication middleware** with JWT support

### Security Best Practices

1. **Always use HTTPS** in production
2. **Configure CORS** restrictively
3. **Enable rate limiting** for public APIs
4. **Validate all inputs** in handlers
5. **Use secure authentication** methods
6. **Keep dependencies updated**

## Deployment Considerations

### Production Deployment

```go
config := &server.ServerConfig{
    Framework: "gin",
    Mode:     "release",
    Host:     "0.0.0.0",
    Port:     8080,
    
    TLS: server.TLSConfig{
        Enabled: true,
        AutoTLS: true,
        Domains: []string{"yourdomain.com"},
    },
    
    GracefulShutdown: server.GracefulShutdownConfig{
        Enabled: true,
        Timeout: 30 * time.Second,
    },
}
```

### Monitoring Integration

```go
// Metrics middleware for monitoring
metrics := middleware.Metrics()
framework.RegisterMiddleware(metrics)

// Health check endpoint
framework.RegisterRoute("GET", "/health", healthHandler)
```

## Migration Guide

### From Framework-Specific Code

**Before (Gin-specific):**
```go
router := gin.Default()
router.GET("/users", func(c *gin.Context) {
    c.JSON(200, users)
})
router.Run(":8080")
```

**After (Framework-agnostic):**
```go
config := server.DefaultServerConfig()
config.Framework = "gin"
manager, _ := server.CreateServerManager("gin", config)

framework := manager.GetFramework()
framework.RegisterRoute("GET", "/users", 
    server.HandlerFunc(func(ctx server.Context) error {
        return ctx.JSON(200, users)
    }),
)

manager.Start(context.Background())
```

### Benefits of Migration

1. **Framework flexibility** - switch frameworks without code changes
2. **Consistent interfaces** - unified API across frameworks
3. **Better testing** - mock frameworks and isolated tests
4. **Enhanced features** - built-in middleware and configuration
5. **Future-proofing** - protection against framework evolution

## Next Steps

Now that you understand the server module architecture, explore these guides:

1. **[Quick Start Guide](01_quick_start.md)** - Get hands-on experience
2. **[Configuration Guide](02_configuration.md)** - Learn about configuration options
3. **[Framework Plugins](03_framework_plugins.md)** - Deep dive into framework specifics
4. **[Middleware System](04_middleware_system.md)** - Understand middleware architecture
5. **[Server Management](05_server_management.md)** - Master lifecycle management
6. **[Integration Examples](06_integration_examples.md)** - See real-world patterns
7. **[Best Practices](07_best_practices.md)** - Production recommendations

## Contributing

The server module is designed for extensibility. Contributions are welcome in:

- **New framework plugins**
- **Additional middleware components**
- **Performance optimizations**
- **Documentation improvements**
- **Testing enhancements**

See the [Contributing Guide](../../../CONTRIBUTING.md) for details.
