# Framework Plugins Guide

The lmcc-go-sdk server module uses a plugin architecture to support multiple web frameworks. This guide covers how to use the built-in framework plugins (Gin, Echo, Fiber) and how to develop custom plugins.

## Plugin Architecture Overview

### Core Components

The plugin system consists of several key components:

1. **FrameworkPlugin Interface**: Defines the plugin contract
2. **Plugin Registry**: Manages plugin registration and discovery
3. **Framework Adapters**: Implement framework-specific logic
4. **Service Integration**: Connects plugins with other SDK modules

### Plugin Interface

Every framework plugin implements the `FrameworkPlugin` interface:

```go
type FrameworkPlugin interface {
    Name() string                                    // Plugin identifier
    Version() string                                 // Plugin version
    Description() string                             // Plugin description
    DefaultConfig() interface{}                      // Default configuration
    CreateFramework(config, services) (WebFramework, error)  // Factory method
    ValidateConfig(config interface{}) error         // Configuration validation
    GetConfigSchema() interface{}                    // Configuration schema
}
```

### Plugin Registration

Plugins are automatically registered during package import:

```go
import (
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin"
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/echo"
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/fiber"
)
```

Or manually registered:

```go
plugin := myplugin.NewPlugin()
err := server.RegisterFramework(plugin)
if err != nil {
    panic("Failed to register plugin: " + err.Error())
}
```

## Built-in Framework Plugins

### Gin Plugin

The Gin plugin provides high-performance HTTP routing with a rich ecosystem.

#### Features
- **High Performance**: Fast routing and minimal memory footprint
- **Rich Ecosystem**: Extensive middleware library
- **JSON/XML Binding**: Built-in request/response binding
- **Template Rendering**: HTML template support
- **Testing Support**: Built-in testing utilities

#### Configuration

```go
config := &server.ServerConfig{
    Framework: "gin",
    Host:     "0.0.0.0",
    Port:     8080,
    Mode:     "release", // gin modes: debug, release, test
    
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

#### Basic Usage

```go
package main

import (
    "context"
    "log"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin"
)

func main() {
    // Create server with Gin plugin
    config := server.DefaultServerConfig()
    config.Framework = "gin"
    config.Port = 8080
    
    manager, err := server.CreateServerManager("gin", config)
    if err != nil {
        log.Fatal("Failed to create server:", err)
    }
    
    // Get framework instance
    framework := manager.GetFramework()
    
    // Register routes
    framework.RegisterRoute("GET", "/hello", server.HandlerFunc(func(ctx server.Context) error {
        return ctx.JSON(200, map[string]string{
            "message": "Hello from Gin!",
            "framework": "gin",
        })
    }))
    
    // Start server
    if err := manager.Start(context.Background()); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}
```

#### Gin-Specific Features

Access Gin's native engine for framework-specific features:

```go
framework := manager.GetFramework()
ginEngine := framework.GetNativeEngine().(*gin.Engine)

// Use Gin-specific features
ginEngine.HTMLRender = &ginrender.HTMLDebug{Glob: "templates/*"}
ginEngine.Static("/static", "./static")

// Custom Gin middleware
ginEngine.Use(func(c *gin.Context) {
    // Gin-specific middleware logic
    c.Header("X-Framework", "Gin")
    c.Next()
})
```

#### Performance Tuning

```go
config.Plugins["gin"] = map[string]interface{}{
    "trusted_proxies":       []string{},  // Disable if not behind proxy
    "redirect_trailing_slash": false,     // Disable for better performance
    "use_raw_path":         true,         // Skip path cleaning
    "max_multipart_memory": 1 << 20,     // Reduce memory usage
}
```

### Echo Plugin

The Echo plugin provides a minimalist and extensible web framework.

#### Features
- **Minimalist Design**: Clean and simple architecture
- **High Performance**: Fast routing with zero memory allocation
- **Middleware**: Extensive middleware support
- **Data Binding**: JSON, XML, form data binding
- **Automatic TLS**: Built-in Let's Encrypt support

#### Configuration

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

#### Basic Usage

```go
package main

import (
    "context"
    "log"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/echo"
)

func main() {
    // Create server with Echo plugin
    config := server.DefaultServerConfig()
    config.Framework = "echo"
    config.Port = 8080
    
    manager, err := server.CreateServerManager("echo", config)
    if err != nil {
        log.Fatal("Failed to create server:", err)
    }
    
    framework := manager.GetFramework()
    
    // Register routes
    framework.RegisterRoute("GET", "/hello", server.HandlerFunc(func(ctx server.Context) error {
        return ctx.JSON(200, map[string]string{
            "message": "Hello from Echo!",
            "framework": "echo",
        })
    }))
    
    // REST API example
    api := framework.Group("/api/v1")
    api.RegisterRoute("GET", "/users", server.HandlerFunc(getUsersHandler))
    api.RegisterRoute("POST", "/users", server.HandlerFunc(createUserHandler))
    
    // Start server
    if err := manager.Start(context.Background()); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}

func getUsersHandler(ctx server.Context) error {
    // Simulate user data
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
    
    // Simulate user creation
    user["id"] = 123
    return ctx.JSON(201, user)
}
```

#### Echo-Specific Features

```go
framework := manager.GetFramework()
echoInstance := framework.GetNativeEngine().(*echo.Echo)

// Custom error handler
echoInstance.HTTPErrorHandler = func(err error, c echo.Context) {
    code := http.StatusInternalServerError
    message := "Internal Server Error"
    
    if he, ok := err.(*echo.HTTPError); ok {
        code = he.Code
        message = he.Message.(string)
    }
    
    c.JSON(code, map[string]string{"error": message})
}

// Custom validator
echoInstance.Validator = &CustomValidator{}

// Template renderer
echoInstance.Renderer = &TemplateRenderer{
    templates: template.Must(template.ParseGlob("templates/*.html")),
}
```

### Fiber Plugin

The Fiber plugin provides an Express.js-inspired framework with high performance.

#### Features
- **Express-inspired**: Familiar API for Node.js developers
- **High Performance**: Built on fasthttp for maximum speed
- **Low Memory**: Minimal memory footprint
- **Built-in Features**: Compression, caching, rate limiting
- **WebSocket Support**: Built-in WebSocket support

#### Configuration

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
            "concurrency":         256 * 1024,      // 256k concurrent connections
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

#### Basic Usage

```go
package main

import (
    "context"
    "log"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/fiber"
)

func main() {
    // Create server with Fiber plugin
    config := server.DefaultServerConfig()
    config.Framework = "fiber"
    config.Port = 8080
    
    manager, err := server.CreateServerManager("fiber", config)
    if err != nil {
        log.Fatal("Failed to create server:", err)
    }
    
    framework := manager.GetFramework()
    
    // Register routes
    framework.RegisterRoute("GET", "/hello", server.HandlerFunc(func(ctx server.Context) error {
        return ctx.JSON(200, map[string]string{
            "message": "Hello from Fiber!",
            "framework": "fiber",
        })
    }))
    
    // File upload example
    framework.RegisterRoute("POST", "/upload", server.HandlerFunc(fileUploadHandler))
    
    // WebSocket endpoint (using native Fiber features)
    setupWebSocket(framework)
    
    // Start server
    if err := manager.Start(context.Background()); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}

func fileUploadHandler(ctx server.Context) error {
    // File upload handling
    request := ctx.Request()
    if err := request.ParseMultipartForm(32 << 20); err != nil { // 32MB max
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
    // Access Fiber's native engine for WebSocket
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

#### High-Performance Configuration

```go
config.Plugins["fiber"] = map[string]interface{}{
    "prefork":          true,           // Enable prefork for production
    "concurrency":      512 * 1024,     // Increase concurrent connections
    "compress_body":    true,           // Enable compression
    "etag":            true,            // Enable ETag
    "read_buffer_size": 8192,          // Increase buffer size
    "write_buffer_size": 8192,         // Increase buffer size
    "network":          "tcp4",        // Use IPv4 only
}
```

## Plugin Management

### Listing Available Plugins

```go
// List all registered plugins
plugins := server.ListFrameworks()
fmt.Println("Available frameworks:", plugins)

// Get plugin information
for _, name := range plugins {
    info, err := server.GetFrameworkInfo(name)
    if err != nil {
        continue
    }
    fmt.Printf("Plugin: %s v%s - %s\n", 
        info["name"], info["version"], info["description"])
}
```

### Plugin Discovery

```go
// Get all plugin information
allInfo := server.GetAllFrameworkInfo()
for name, info := range allInfo {
    fmt.Printf("Framework: %s\n", name)
    fmt.Printf("  Version: %s\n", info["version"])
    fmt.Printf("  Description: %s\n", info["description"])
}
```

### Default Plugin Management

```go
// Get default plugin
defaultPlugin, err := server.GetDefaultFramework()
if err != nil {
    log.Fatal("No default plugin available")
}
fmt.Println("Default framework:", defaultPlugin.Name())

// Set default plugin
err = server.SetDefaultFramework("gin")
if err != nil {
    log.Fatal("Failed to set default framework:", err)
}
```

### Plugin Validation

```go
// Validate plugin configuration
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

## Framework Comparison

### Performance Comparison

| Feature          | Gin        | Echo       | Fiber      |
|------------------|------------|------------|------------|
| **Requests/sec** | ~47,000    | ~45,000    | ~50,000    |
| **Memory Usage** | Moderate   | Low        | Very Low   |
| **Cold Start**   | Fast       | Very Fast  | Fastest    |
| **Binary Size**  | Medium     | Small      | Medium     |
| **CPU Usage**    | Low        | Low        | Very Low   |

### Feature Comparison

| Feature              | Gin | Echo | Fiber |
|---------------------|-----|------|-------|
| **Router**          | ✅   | ✅    | ✅     |
| **Middleware**      | ✅   | ✅    | ✅     |
| **JSON Binding**    | ✅   | ✅    | ✅     |
| **Template Engine** | ✅   | ✅    | ✅     |
| **Static Files**    | ✅   | ✅    | ✅     |
| **WebSocket**       | ❌   | ✅    | ✅     |
| **HTTP/2**          | ✅   | ✅    | ❌     |
| **Auto TLS**        | ❌   | ✅    | ❌     |
| **Prefork**         | ❌   | ❌    | ✅     |

### Use Case Recommendations

#### Choose Gin When:
- You need a **mature ecosystem** with extensive middleware
- Building **traditional web applications** with templates
- Requiring **JSON/XML binding** and validation
- Team has **Go web development experience**
- Need **stability and community support**

#### Choose Echo When:
- Building **RESTful APIs** and microservices
- Need **automatic TLS** with Let's Encrypt
- Want **minimal memory footprint**
- Prefer **clean and simple architecture**
- Building **middleware-heavy applications**

#### Choose Fiber When:
- Need **maximum performance** and speed
- Building **real-time applications**
- Coming from **Node.js/Express background**
- Need **WebSocket support**
- Building **high-concurrency services**

## Custom Plugin Development

### Plugin Interface Implementation

Create a custom framework plugin:

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
    
    // Create and return your framework adapter
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
            // Add your configuration schema here
        },
        "required": []string{"framework"},
    }
}

// Auto-register plugin
func init() {
    server.RegisterFramework(NewPlugin())
}
```

### Framework Adapter Implementation

Implement the WebFramework interface:

```go
type MyFrameworkAdapter struct {
    config   *server.ServerConfig
    services services.ServiceContainer
    engine   *MyNativeFramework
    server   *http.Server
}

func NewMyFrameworkAdapter(config *server.ServerConfig, services services.ServiceContainer) *MyFrameworkAdapter {
    return &MyFrameworkAdapter{
        config:   config,
        services: services,
        engine:   NewMyNativeFramework(),
    }
}

func (a *MyFrameworkAdapter) Start(ctx context.Context) error {
    // Implement server startup logic
    addr := fmt.Sprintf("%s:%d", a.config.Host, a.config.Port)
    
    a.server = &http.Server{
        Addr:           addr,
        Handler:        a.engine,
        ReadTimeout:    a.config.ReadTimeout,
        WriteTimeout:   a.config.WriteTimeout,
        IdleTimeout:    a.config.IdleTimeout,
        MaxHeaderBytes: a.config.MaxHeaderBytes,
    }
    
    // Start server in goroutine
    go func() {
        if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            a.services.GetLogger().Errorw("Server failed to start", "error", err)
        }
    }()
    
    a.services.GetLogger().Infow("Server started", "address", addr)
    return nil
}

func (a *MyFrameworkAdapter) Stop(ctx context.Context) error {
    if a.server == nil {
        return nil
    }
    
    return a.server.Shutdown(ctx)
}

func (a *MyFrameworkAdapter) RegisterRoute(method, path string, handler server.Handler) error {
    // Implement route registration logic
    a.engine.AddRoute(method, path, func(w http.ResponseWriter, r *http.Request) {
        // Create unified context
        ctx := NewMyFrameworkContext(w, r)
        
        // Call handler
        if err := handler.Handle(ctx); err != nil {
            // Handle error using error handling service
            a.services.GetErrorHandler().Handle(ctx, err)
        }
    })
    
    return nil
}

func (a *MyFrameworkAdapter) RegisterMiddleware(middleware server.Middleware) error {
    // Implement middleware registration
    a.engine.Use(func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            ctx := NewMyFrameworkContext(w, r)
            
            middleware.Process(ctx, func() error {
                next(w, r)
                return nil
            })
        }
    })
    
    return nil
}

func (a *MyFrameworkAdapter) Group(prefix string, middlewares ...server.Middleware) server.RouteGroup {
    // Implement route grouping
    return NewMyFrameworkGroup(a.engine, prefix, middlewares...)
}

func (a *MyFrameworkAdapter) GetNativeEngine() interface{} {
    return a.engine
}

func (a *MyFrameworkAdapter) GetConfig() *server.ServerConfig {
    return a.config
}
```

### Context Adapter Implementation

Implement the unified Context interface:

```go
type MyFrameworkContext struct {
    writer  http.ResponseWriter
    request *http.Request
    params  map[string]string
    data    map[string]interface{}
}

func NewMyFrameworkContext(w http.ResponseWriter, r *http.Request) *MyFrameworkContext {
    return &MyFrameworkContext{
        writer:  w,
        request: r,
        params:  make(map[string]string),
        data:    make(map[string]interface{}),
    }
}

func (c *MyFrameworkContext) Request() *http.Request {
    return c.request
}

func (c *MyFrameworkContext) Response() http.ResponseWriter {
    return c.writer
}

func (c *MyFrameworkContext) Param(key string) string {
    return c.params[key]
}

func (c *MyFrameworkContext) Query(key string) string {
    return c.request.URL.Query().Get(key)
}

func (c *MyFrameworkContext) Header(key string) string {
    return c.request.Header.Get(key)
}

func (c *MyFrameworkContext) SetHeader(key, value string) {
    c.writer.Header().Set(key, value)
}

func (c *MyFrameworkContext) JSON(code int, obj interface{}) error {
    c.SetHeader("Content-Type", "application/json")
    c.writer.WriteHeader(code)
    
    return json.NewEncoder(c.writer).Encode(obj)
}

func (c *MyFrameworkContext) String(code int, format string, values ...interface{}) error {
    c.SetHeader("Content-Type", "text/plain")
    c.writer.WriteHeader(code)
    
    _, err := fmt.Fprintf(c.writer, format, values...)
    return err
}

func (c *MyFrameworkContext) Bind(obj interface{}) error {
    // Implement request binding logic
    contentType := c.Header("Content-Type")
    
    if strings.Contains(contentType, "application/json") {
        return json.NewDecoder(c.request.Body).Decode(obj)
    }
    
    return fmt.Errorf("unsupported content type: %s", contentType)
}

// Implement remaining Context interface methods...
```

### Plugin Registration

Register your custom plugin:

```go
package main

import (
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
    "myapp/plugins/myframework"
)

func main() {
    // Register custom plugin
    plugin := myframework.NewPlugin()
    if err := server.RegisterFramework(plugin); err != nil {
        panic("Failed to register custom plugin: " + err.Error())
    }
    
    // Use custom plugin
    config := server.DefaultServerConfig()
    config.Framework = "myframework"
    
    manager, err := server.CreateServerManager("myframework", config)
    if err != nil {
        panic("Failed to create server: " + err.Error())
    }
    
    // ... rest of application
}
```

## Testing Framework Plugins

### Plugin Testing

```go
func TestMyFrameworkPlugin(t *testing.T) {
    plugin := myframework.NewPlugin()
    
    // Test plugin metadata
    assert.Equal(t, "myframework", plugin.Name())
    assert.Equal(t, "1.0.0", plugin.Version())
    assert.NotEmpty(t, plugin.Description())
    
    // Test default configuration
    config := plugin.DefaultConfig()
    assert.NotNil(t, config)
    
    serverConfig, ok := config.(*server.ServerConfig)
    assert.True(t, ok)
    assert.Equal(t, "myframework", serverConfig.Framework)
    
    // Test configuration validation
    err := plugin.ValidateConfig(serverConfig)
    assert.NoError(t, err)
    
    // Test invalid configuration
    err = plugin.ValidateConfig("invalid")
    assert.Error(t, err)
}
```

### Framework Integration Testing

```go
func TestFrameworkIntegration(t *testing.T) {
    config := server.DefaultServerConfig()
    config.Framework = "myframework"
    config.Port = 0 // Use random port for testing
    
    // Create mock service container
    services := &MockServiceContainer{}
    
    // Create framework
    plugin := myframework.NewPlugin()
    framework, err := plugin.CreateFramework(config, services)
    assert.NoError(t, err)
    assert.NotNil(t, framework)
    
    // Test route registration
    err = framework.RegisterRoute("GET", "/test", server.HandlerFunc(func(ctx server.Context) error {
        return ctx.JSON(200, map[string]string{"status": "ok"})
    }))
    assert.NoError(t, err)
    
    // Test middleware registration
    err = framework.RegisterMiddleware(server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
        ctx.SetHeader("X-Test", "true")
        return next()
    }))
    assert.NoError(t, err)
    
    // Test server lifecycle
    ctx := context.Background()
    
    err = framework.Start(ctx)
    assert.NoError(t, err)
    
    // Make test request
    // ... test HTTP requests
    
    err = framework.Stop(ctx)
    assert.NoError(t, err)
}
```

## Best Practices

### Plugin Development Guidelines

1. **Follow Interface Contracts**: Implement all interface methods correctly
2. **Handle Errors Gracefully**: Provide meaningful error messages
3. **Support Configuration**: Accept and validate configuration properly
4. **Service Integration**: Use provided service container
5. **Thread Safety**: Ensure plugin is thread-safe
6. **Resource Cleanup**: Properly clean up resources in Stop method

### Performance Optimization

1. **Minimize Allocations**: Reduce memory allocations in hot paths
2. **Connection Pooling**: Reuse connections and resources
3. **Efficient Routing**: Use efficient routing algorithms
4. **Middleware Ordering**: Order middleware for optimal performance
5. **Caching**: Implement appropriate caching strategies

### Security Considerations

1. **Input Validation**: Validate all inputs and configurations
2. **Error Handling**: Don't expose sensitive information in errors
3. **Resource Limits**: Implement proper resource limits
4. **Security Headers**: Add appropriate security headers
5. **TLS Support**: Support TLS/HTTPS properly

## Troubleshooting

### Common Plugin Issues

1. **Plugin Not Found**
   ```go
   // Ensure plugin is imported
   import _ "path/to/your/plugin"
   
   // Check plugin registration
   plugins := server.ListFrameworks()
   fmt.Println("Available plugins:", plugins)
   ```

2. **Configuration Validation Errors**
   ```go
   // Validate configuration before use
   if err := plugin.ValidateConfig(config); err != nil {
       log.Printf("Configuration error: %v", err)
   }
   ```

3. **Service Integration Issues**
   ```go
   // Ensure service container is properly passed
   framework, err := plugin.CreateFramework(config, serviceContainer)
   if err != nil {
       log.Printf("Service integration error: %v", err)
   }
   ```

### Debug Information

```go
// Get detailed plugin information
info := server.GetAllFrameworkInfo()
for name, details := range info {
    fmt.Printf("Plugin: %s\n", name)
    fmt.Printf("  Version: %s\n", details["version"])
    fmt.Printf("  Description: %s\n", details["description"])
}

// Test plugin functionality
plugin, err := server.GetFramework("gin")
if err != nil {
    log.Fatal("Plugin not available:", err)
}

schema := plugin.GetConfigSchema()
fmt.Printf("Config schema: %+v\n", schema)
```

## Next Steps

- **[Middleware System](04_middleware_system.md)** - Learn about middleware architecture
- **[Server Management](05_server_management.md)** - Master server lifecycle management
- **[Integration Examples](06_integration_examples.md)** - See real-world integration patterns
- **[Best Practices](07_best_practices.md)** - Production deployment guidelines
