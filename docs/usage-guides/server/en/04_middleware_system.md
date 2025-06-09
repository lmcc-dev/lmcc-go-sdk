# Middleware System Guide

The lmcc-go-sdk server module provides a comprehensive middleware system that enables request/response processing in a modular and composable way. This guide covers the unified middleware interface, built-in middleware components, and custom middleware development.

## Middleware Architecture

### Core Concepts

The middleware system is built around several key concepts:

1. **Unified Interface**: All middleware implements the same interface across frameworks
2. **Chain Processing**: Middleware is executed in a defined order
3. **Context Passing**: Request context flows through the middleware chain
4. **Framework Agnostic**: Middleware works with any supported framework
5. **Service Integration**: Middleware integrates with other SDK modules

### Middleware Interface

All middleware implements the unified `Middleware` interface:

```go
type Middleware interface {
    Process(ctx Context, next func() error) error
}

// Function type for convenience
type MiddlewareFunc func(ctx Context, next func() error) error

func (f MiddlewareFunc) Process(ctx Context, next func() error) error {
    return f(ctx, next)
}
```

### Middleware Chain Execution

```
Request
   ↓
┌─────────────────┐
│ Logger          │ ← Log request details
├─────────────────┤
│ Recovery        │ ← Recover from panics
├─────────────────┤
│ CORS            │ ← Handle cross-origin requests
├─────────────────┤
│ Security        │ ← Add security headers
├─────────────────┤
│ Compression     │ ← Compress responses
├─────────────────┤
│ Authentication  │ ← Validate user credentials
├─────────────────┤
│ Rate Limiting   │ ← Control request rate
├─────────────────┤
│ Metrics         │ ← Collect performance metrics
├─────────────────┤
│ Your Handler    │ ← Process business logic
└─────────────────┘
   ↓
Response
```

> **⚠️ Implementation Status Notice**: The middleware system is currently under active development. Please check the implementation status below before using.

## Implementation Status

| Middleware | Status | Gin | Echo | Fiber | Notes |
|------------|--------|-----|------|-------|-------|
| **Logger** | ⚠️ Partial | ❌ | ✅ | ✅ | Gin implementation pending |
| **Recovery** | ⚠️ Partial | ❌ | ✅ | ✅ | Gin implementation pending |
| **CORS** | ✅ Complete | ✅ | ✅ | ✅ | Fully implemented |
| **Rate Limiting** | 🚧 In Development | ❌ | ❌ | ❌ | Interface only |
| **Authentication** | 🚧 In Development | ❌ | ❌ | ❌ | Interface only |
| **Security** | 🚧 In Development | ❌ | ❌ | ❌ | Interface only |
| **Compression** | 🚧 In Development | ❌ | ❌ | ❌ | Interface only |
| **Metrics** | 🚧 In Development | ❌ | ❌ | ❌ | Interface only |

### Currently Available (Fully Implemented)
- **CORS Middleware**: Complete cross-framework implementation
- **Logger Middleware**: Available for Echo and Fiber
- **Recovery Middleware**: Available for Echo and Fiber

### Coming Soon (Under Development)
The following middleware are planned for implementation in the next releases:
- Security Headers Middleware (ETA: 2025-01-12)
- Rate Limiting Middleware (ETA: 2025-01-13)
- Authentication Middleware (ETA: 2025-01-14)
- Compression Middleware (ETA: 2025-01-14)
- Metrics Collection Middleware (ETA: 2025-01-15)

## Built-in Middleware

### Logger Middleware

The logger middleware records request and response information.

#### Configuration

```go
config.Middleware.Logger = server.LoggerMiddlewareConfig{
    Enabled:     true,
    Format:      "json",        // json or text
    SkipPaths:   []string{"/health", "/favicon.ico"},
    IncludeBody: false,         // Include request/response body
    MaxBodySize: 1024,          // Maximum body size to log (bytes)
}
```

#### Usage Examples

```go
// Basic logger middleware
logger := middleware.NewLogger(server.LoggerMiddlewareConfig{
    Enabled: true,
    Format:  "json",
})

framework.RegisterMiddleware(logger)
```

#### Custom Logger Configuration

```go
// Detailed logging for development
logger := middleware.NewLogger(server.LoggerMiddlewareConfig{
    Enabled:     true,
    Format:      "text",
    SkipPaths:   []string{},     // Log all requests
    IncludeBody: true,           // Include request bodies
    MaxBodySize: 4096,           // Log up to 4KB
})

// Production logging
logger := middleware.NewLogger(server.LoggerMiddlewareConfig{
    Enabled:     true,
    Format:      "json",
    SkipPaths:   []string{"/health", "/metrics", "/static/*"},
    IncludeBody: false,          // Don't log request bodies
    MaxBodySize: 0,              // No body logging
})
```

#### Log Output Format

**JSON Format:**
```json
{
    "timestamp": "2024-01-09T10:30:00Z",
    "method": "GET",
    "path": "/api/users",
    "status": 200,
    "duration": "15ms",
    "client_ip": "192.168.1.100",
    "user_agent": "curl/7.68.0",
    "request_id": "req_123456"
}
```

**Text Format:**
```
2024-01-09 10:30:00 | 200 | 15ms | 192.168.1.100 | GET /api/users
```

### Recovery Middleware

The recovery middleware catches panics and converts them to proper HTTP error responses.

#### Configuration

```go
config.Middleware.Recovery = server.RecoveryMiddlewareConfig{
    Enabled:             true,
    PrintStack:          true,   // Print stack trace
    DisableStackAll:     false,  // Disable all stack traces
    DisableColorConsole: false,  // Disable colored output
}
```

#### Usage Examples

```go
// Basic recovery middleware
recovery := middleware.NewRecovery(server.RecoveryMiddlewareConfig{
    Enabled:    true,
    PrintStack: true,
})

framework.RegisterMiddleware(recovery)
```

#### Custom Recovery Handler

```go
// Custom recovery with error reporting
recovery := middleware.NewRecovery(server.RecoveryMiddlewareConfig{
    Enabled:    true,
    PrintStack: false,  // Don't print to console
})

// Set custom recovery handler
recovery.SetCustomRecoveryHandler(func(ctx server.Context, err interface{}) {
    // Log error to monitoring service
    log.Errorf("Panic recovered: %v", err)
    
    // Send to error tracking service
    errorTracker.ReportPanic(ctx.Request().Context(), err)
    
    // Return custom error response
    ctx.JSON(500, map[string]string{
        "error": "Internal server error",
        "request_id": ctx.GetString("request_id"),
    })
})

framework.RegisterMiddleware(recovery)
```

#### Production vs Development Configuration

```go
// Development configuration
if config.IsDebugMode() {
    recovery := middleware.NewRecovery(server.RecoveryMiddlewareConfig{
        Enabled:             true,
        PrintStack:          true,   // Show full stack traces
        DisableColorConsole: false,  // Colored output
    })
} else {
    // Production configuration
    recovery := middleware.NewRecovery(server.RecoveryMiddlewareConfig{
        Enabled:             true,
        PrintStack:          false,  // No stack traces in logs
        DisableStackAll:     true,   // Completely disable stack traces
        DisableColorConsole: true,   // No colored output
    })
}
```

### CORS Middleware

The CORS middleware handles Cross-Origin Resource Sharing requests.

#### Configuration

```go
config.CORS = server.CORSConfig{
    Enabled:          true,
    AllowOrigins:     []string{"https://example.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:          12 * time.Hour,
}
```

#### Usage Examples

```go
// Basic CORS for API
cors := middleware.NewCORS(server.CORSConfig{
    Enabled:      true,
    AllowOrigins: []string{"https://myapp.com", "https://admin.myapp.com"},
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders: []string{"Content-Type", "Authorization"},
})

framework.RegisterMiddleware(cors)
```

#### Development CORS (Allow All)

```go
// Development CORS - allow everything
cors := middleware.NewCORS(server.CORSConfig{
    Enabled:      true,
    AllowOrigins: []string{"*"},
    AllowMethods: []string{"*"},
    AllowHeaders: []string{"*"},
    AllowCredentials: false,  // Must be false when origins is "*"
})
```

#### Restrictive Production CORS

```go
// Production CORS - restrictive settings
cors := middleware.NewCORS(server.CORSConfig{
    Enabled:          true,
    AllowOrigins:     []string{"https://yourdomain.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Content-Type", "Authorization", "X-Requested-With"},
    ExposeHeaders:    []string{"X-Total-Count", "X-Page-Count"},
    AllowCredentials: true,
    MaxAge:          1 * time.Hour,  // Cache preflight for 1 hour
})
```

#### Dynamic CORS Configuration

```go
// Dynamic CORS based on request
cors := middleware.NewCORS(server.CORSConfig{
    Enabled: true,
})

// Set dynamic origin validation
cors.SetOriginValidator(func(origin string) bool {
    // Check against database or configuration
    allowedOrigins := getAllowedOrigins()
    for _, allowed := range allowedOrigins {
        if origin == allowed {
            return true
        }
    }
    return false
})

framework.RegisterMiddleware(cors)
```

### Rate Limiting Middleware

The rate limiting middleware controls the rate of incoming requests.

#### Configuration

```go
config.Middleware.RateLimit = server.RateLimitMiddlewareConfig{
    Enabled: true,
    Rate:    100.0,      // 100 requests per second
    Burst:   200,        // Allow burst of 200 requests
    KeyFunc: "ip",       // Rate limit by IP address
}
```

#### Usage Examples

```go
// Basic rate limiting by IP
rateLimit := middleware.NewRateLimit(server.RateLimitMiddlewareConfig{
    Enabled: true,
    Rate:    50.0,   // 50 requests per second
    Burst:   100,    // Allow burst of 100 requests
    KeyFunc: "ip",
})

framework.RegisterMiddleware(rateLimit)
```

#### Custom Rate Limiting Key

```go
// Rate limit by user ID
rateLimit := middleware.NewRateLimit(server.RateLimitMiddlewareConfig{
    Enabled: true,
    Rate:    10.0,   // 10 requests per second per user
    Burst:   20,     // Allow burst of 20 requests
})

// Set custom key function
rateLimit.SetKeyFunc(func(ctx server.Context) string {
    userID := ctx.GetString("user_id")
    if userID == "" {
        return ctx.ClientIP()  // Fallback to IP
    }
    return "user:" + userID
})

framework.RegisterMiddleware(rateLimit)
```

#### Different Limits for Different Endpoints

```go
// API rate limiting
apiRateLimit := middleware.NewRateLimit(server.RateLimitMiddlewareConfig{
    Enabled: true,
    Rate:    100.0,
    Burst:   200,
})

// Upload rate limiting (more restrictive)
uploadRateLimit := middleware.NewRateLimit(server.RateLimitMiddlewareConfig{
    Enabled: true,
    Rate:    5.0,    // 5 uploads per second
    Burst:   10,     // Allow burst of 10 uploads
})

// Apply different rate limits to different route groups
api := framework.Group("/api")
api.RegisterMiddleware(apiRateLimit)

upload := framework.Group("/upload")
upload.RegisterMiddleware(uploadRateLimit)
```

### Authentication Middleware

The authentication middleware validates user credentials and manages user sessions.

#### Configuration

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
```

#### JWT Authentication

```go
// JWT authentication middleware
auth := middleware.NewAuth(server.AuthMiddlewareConfig{
    Enabled: true,
    Type:    "jwt",
    SkipPaths: []string{"/auth/login", "/auth/register", "/public/*"},
})

// Set JWT validation function
auth.SetAuthFunc(func(ctx server.Context) (interface{}, error) {
    token := ctx.Header("Authorization")
    if token == "" {
        return nil, errors.New("missing authorization header")
    }
    
    // Remove "Bearer " prefix
    if strings.HasPrefix(token, "Bearer ") {
        token = token[7:]
    }
    
    // Validate JWT token
    claims, err := validateJWTToken(token)
    if err != nil {
        return nil, err
    }
    
    return claims, nil
})

framework.RegisterMiddleware(auth)
```

#### Basic Authentication

```go
// Basic authentication middleware
auth := middleware.NewAuth(server.AuthMiddlewareConfig{
    Enabled: true,
    Type:    "basic",
})

auth.SetAuthFunc(func(ctx server.Context) (interface{}, error) {
    username, password, ok := ctx.Request().BasicAuth()
    if !ok {
        return nil, errors.New("missing basic auth")
    }
    
    // Validate credentials
    user, err := validateCredentials(username, password)
    if err != nil {
        return nil, err
    }
    
    return user, nil
})

framework.RegisterMiddleware(auth)
```

#### Custom Unauthorized Handler

```go
auth := middleware.NewAuth(server.AuthMiddlewareConfig{
    Enabled: true,
    Type:    "custom",
})

// Set custom unauthorized response
auth.SetUnauthorizedHandler(func(ctx server.Context) {
    ctx.JSON(401, map[string]interface{}{
        "error": "Unauthorized",
        "code":  "AUTH_REQUIRED",
        "message": "Please provide valid authentication credentials",
    })
})
```

### Security Middleware

The security middleware adds security headers to protect against common vulnerabilities.

#### Basic Security Headers

```go
// Basic security middleware
security := middleware.NewSecurity()

security.SetXSSProtection(true)               // X-XSS-Protection: 1; mode=block
security.SetContentTypeNosniff(true)          // X-Content-Type-Options: nosniff
security.SetFrameOptions("DENY")              // X-Frame-Options: DENY
security.SetHSTSMaxAge(31536000)              // Strict-Transport-Security: max-age=31536000

framework.RegisterMiddleware(security)
```

#### Comprehensive Security Configuration

```go
security := middleware.NewSecurity()

// XSS Protection
security.SetXSSProtection(true)

// Content Type Options
security.SetContentTypeNosniff(true)

// Frame Options
security.SetFrameOptions("DENY")

// HSTS (HTTP Strict Transport Security)
security.SetHSTSMaxAge(31536000)  // 1 year
security.SetHSTSIncludeSubdomains(true)

// Content Security Policy
security.SetContentSecurityPolicy("default-src 'self'; script-src 'self' 'unsafe-inline'")

// Referrer Policy
security.SetReferrerPolicy("strict-origin-when-cross-origin")

// Permissions Policy
security.SetPermissionsPolicy("geolocation=(), microphone=(), camera=()")

framework.RegisterMiddleware(security)
```

### Compression Middleware

The compression middleware compresses HTTP responses to reduce bandwidth usage.

#### Basic Compression

```go
// Basic compression middleware
compression := middleware.NewCompression()

compression.SetLevel(6)           // Compression level (1-9)
compression.SetMinLength(1024)    // Minimum response size to compress
compression.SetExcludedExtensions([]string{".jpg", ".png", ".gif", ".zip"})

framework.RegisterMiddleware(compression)
```

#### Advanced Compression Configuration

```go
compression := middleware.NewCompression()

// Compression settings
compression.SetLevel(6)                    // Balance between speed and compression
compression.SetMinLength(1024)             // Only compress responses > 1KB

// Exclude already compressed content
compression.SetExcludedExtensions([]string{
    ".jpg", ".jpeg", ".png", ".gif", ".webp",  // Images
    ".mp4", ".avi", ".mov",                    // Videos
    ".zip", ".rar", ".7z", ".tar.gz",          // Archives
    ".pdf",                                    // Documents
})

// Exclude specific content types
compression.SetExcludedContentTypes([]string{
    "image/*",
    "video/*",
    "application/zip",
    "application/pdf",
})

framework.RegisterMiddleware(compression)
```

### Metrics Middleware

The metrics middleware collects performance and usage metrics.

#### Basic Metrics Collection

```go
// Basic metrics middleware
metrics := middleware.NewMetrics()

metrics.SetMetricsPath("/metrics")  // Prometheus metrics endpoint
metrics.SetSkipPaths([]string{"/health", "/metrics"})

framework.RegisterMiddleware(metrics)
```

#### Custom Metrics Configuration

```go
metrics := middleware.NewMetrics()

// Metrics configuration
metrics.SetMetricsPath("/metrics")
metrics.SetSkipPaths([]string{"/health", "/metrics", "/favicon.ico"})

// Custom labels
metrics.SetCustomLabels(map[string]string{
    "service": "api-server",
    "version": "1.0.0",
    "environment": "production",
})

// Custom metric collectors
metrics.AddHistogram("request_duration_seconds", "HTTP request duration", []string{"method", "endpoint", "status"})
metrics.AddCounter("request_total", "Total HTTP requests", []string{"method", "endpoint", "status"})
metrics.AddGauge("active_connections", "Active HTTP connections", []string{})

framework.RegisterMiddleware(metrics)
```

## Middleware Chain Management

### Creating Middleware Chains

```go
// Create middleware chain
chain := middleware.NewMiddlewareChain()

// Add middleware in order
chain.Add(middleware.NewLogger(loggerConfig))
chain.Add(middleware.NewRecovery(recoveryConfig))
chain.Add(middleware.NewCORS(corsConfig))
chain.Add(middleware.NewSecurity())
chain.Add(middleware.NewAuth(authConfig))

// Register chain with framework
for _, mw := range chain.GetMiddlewares() {
    framework.RegisterMiddleware(mw)
}
```

### Conditional Middleware

```go
// Add middleware conditionally
if config.IsDebugMode() {
    // Development middleware
    chain.Add(middleware.NewDetailedLogger())
    chain.Add(middleware.NewDebugHeaders())
} else {
    // Production middleware
    chain.Add(middleware.NewProductionLogger())
    chain.Add(middleware.NewSecurity())
}

// Always add essential middleware
chain.Add(middleware.NewRecovery(recoveryConfig))
chain.Add(middleware.NewCORS(corsConfig))
```

### Route-Specific Middleware

```go
// Global middleware
framework.RegisterMiddleware(middleware.NewLogger(loggerConfig))
framework.RegisterMiddleware(middleware.NewRecovery(recoveryConfig))

// API-specific middleware
api := framework.Group("/api")
api.RegisterMiddleware(middleware.NewAuth(authConfig))
api.RegisterMiddleware(middleware.NewRateLimit(rateLimitConfig))

// Admin-specific middleware
admin := framework.Group("/admin")
admin.RegisterMiddleware(middleware.NewAuth(adminAuthConfig))
admin.RegisterMiddleware(middleware.NewAuditLog())

// Public routes (no additional middleware)
public := framework.Group("/public")
```

## Custom Middleware Development

### Simple Custom Middleware

```go
// Request ID middleware
func RequestIDMiddleware() server.Middleware {
    return server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
        // Generate unique request ID
        requestID := generateRequestID()
        
        // Store in context
        ctx.Set("request_id", requestID)
        
        // Add to response header
        ctx.SetHeader("X-Request-ID", requestID)
        
        // Continue to next middleware
        return next()
    })
}

// Usage
framework.RegisterMiddleware(RequestIDMiddleware())
```

### Advanced Custom Middleware

```go
// Custom audit logging middleware
type AuditMiddleware struct {
    auditLogger *log.Logger
    skipPaths   []string
}

func NewAuditMiddleware(auditLogger *log.Logger) *AuditMiddleware {
    return &AuditMiddleware{
        auditLogger: auditLogger,
        skipPaths:   []string{"/health", "/metrics"},
    }
}

func (am *AuditMiddleware) Process(ctx server.Context, next func() error) error {
    // Skip audit logging for certain paths
    path := ctx.Path()
    for _, skipPath := range am.skipPaths {
        if path == skipPath {
            return next()
        }
    }
    
    // Record audit log entry
    start := time.Now()
    userID := ctx.GetString("user_id")
    
    // Execute next middleware/handler
    err := next()
    
    // Log audit entry
    am.auditLogger.Info("Audit log",
        "user_id", userID,
        "method", ctx.Method(),
        "path", path,
        "duration", time.Since(start),
        "status", getResponseStatus(ctx),
        "error", err,
    )
    
    return err
}

// Usage
auditMiddleware := NewAuditMiddleware(auditLogger)
framework.RegisterMiddleware(auditMiddleware)
```

### Middleware with Configuration

```go
// Cache middleware with configuration
type CacheConfig struct {
    TTL         time.Duration
    MaxSize     int64
    ExcludePaths []string
}

type CacheMiddleware struct {
    config CacheConfig
    cache  *cache.Cache
}

func NewCacheMiddleware(config CacheConfig) *CacheMiddleware {
    return &CacheMiddleware{
        config: config,
        cache:  cache.New(config.TTL, config.TTL*2),
    }
}

func (cm *CacheMiddleware) Process(ctx server.Context, next func() error) error {
    // Only cache GET requests
    if ctx.Method() != "GET" {
        return next()
    }
    
    // Check if path should be excluded
    path := ctx.Path()
    for _, excludePath := range cm.config.ExcludePaths {
        if strings.HasPrefix(path, excludePath) {
            return next()
        }
    }
    
    // Generate cache key
    cacheKey := fmt.Sprintf("%s:%s", ctx.Method(), path)
    
    // Check cache
    if cachedResponse, found := cm.cache.Get(cacheKey); found {
        // Return cached response
        response := cachedResponse.(CachedResponse)
        for key, value := range response.Headers {
            ctx.SetHeader(key, value)
        }
        ctx.SetHeader("X-Cache", "HIT")
        return ctx.Data(response.StatusCode, response.ContentType, response.Body)
    }
    
    // Execute handler and capture response
    responseRecorder := NewResponseRecorder(ctx.Response())
    
    err := next()
    if err != nil {
        return err
    }
    
    // Cache the response
    cachedResponse := CachedResponse{
        StatusCode:  responseRecorder.StatusCode,
        Headers:     responseRecorder.Headers,
        ContentType: responseRecorder.ContentType,
        Body:        responseRecorder.Body,
    }
    
    cm.cache.Set(cacheKey, cachedResponse, cm.config.TTL)
    ctx.SetHeader("X-Cache", "MISS")
    
    return nil
}

// Usage
cacheConfig := CacheConfig{
    TTL:         5 * time.Minute,
    MaxSize:     100 * 1024 * 1024, // 100MB
    ExcludePaths: []string{"/api/", "/admin/"},
}

cacheMiddleware := NewCacheMiddleware(cacheConfig)
framework.RegisterMiddleware(cacheMiddleware)
```

## Middleware Testing

### Testing Individual Middleware

```go
func TestRequestIDMiddleware(t *testing.T) {
    // Create mock context
    ctx := &MockContext{}
    
    // Create middleware
    middleware := RequestIDMiddleware()
    
    // Test execution
    called := false
    err := middleware.Process(ctx, func() error {
        called = true
        
        // Verify request ID is set
        requestID := ctx.GetString("request_id")
        assert.NotEmpty(t, requestID)
        assert.Equal(t, requestID, ctx.GetHeader("X-Request-ID"))
        
        return nil
    })
    
    assert.NoError(t, err)
    assert.True(t, called)
}
```

### Testing Middleware Chains

```go
func TestMiddlewareChain(t *testing.T) {
    // Create middleware chain
    chain := middleware.NewMiddlewareChain()
    chain.Add(RequestIDMiddleware())
    chain.Add(middleware.NewLogger(loggerConfig))
    chain.Add(middleware.NewAuth(authConfig))
    
    // Create mock context
    ctx := &MockContext{
        headers: map[string]string{
            "Authorization": "Bearer valid-token",
        },
    }
    
    // Create test handler
    handler := server.HandlerFunc(func(ctx server.Context) error {
        return ctx.JSON(200, map[string]string{"status": "ok"})
    })
    
    // Execute chain
    err := chain.Execute(ctx, handler)
    assert.NoError(t, err)
    
    // Verify middleware effects
    assert.NotEmpty(t, ctx.GetString("request_id"))
    assert.NotEmpty(t, ctx.GetString("user_id"))
}
```

### Integration Testing

```go
func TestMiddlewareIntegration(t *testing.T) {
    // Create test server
    config := server.DefaultServerConfig()
    config.Framework = "gin"
    config.Port = 0 // Random port
    
    manager, err := server.CreateServerManager("gin", config)
    assert.NoError(t, err)
    
    framework := manager.GetFramework()
    
    // Register middleware
    framework.RegisterMiddleware(RequestIDMiddleware())
    framework.RegisterMiddleware(middleware.NewAuth(authConfig))
    
    // Register test route
    framework.RegisterRoute("GET", "/test", server.HandlerFunc(func(ctx server.Context) error {
        return ctx.JSON(200, map[string]interface{}{
            "request_id": ctx.GetString("request_id"),
            "user_id":   ctx.GetString("user_id"),
        })
    }))
    
    // Start server
    ctx := context.Background()
    err = manager.Start(ctx)
    assert.NoError(t, err)
    defer manager.Stop(ctx)
    
    // Make test request
    client := &http.Client{}
    req, _ := http.NewRequest("GET", "http://localhost:"+getPort(manager)+"/test", nil)
    req.Header.Set("Authorization", "Bearer valid-token")
    
    resp, err := client.Do(req)
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
    
    // Verify response headers
    assert.NotEmpty(t, resp.Header.Get("X-Request-ID"))
}
```

## Best Practices

### Middleware Ordering

Order middleware based on their purpose and dependencies:

1. **Logger** - First to log all requests
2. **Recovery** - Early to catch all panics
3. **CORS** - Before authentication for preflight requests
4. **Security** - Add security headers early
5. **Compression** - Before content generation
6. **Authentication** - Before business logic
7. **Rate Limiting** - After authentication
8. **Metrics** - Monitor authenticated requests
9. **Business Logic** - Your handlers

### Performance Considerations

1. **Minimize Allocations**: Avoid creating unnecessary objects in hot paths
2. **Early Returns**: Return early from middleware when possible
3. **Efficient String Operations**: Use string builders for concatenation
4. **Context Management**: Don't store large objects in context
5. **Resource Cleanup**: Always clean up resources in defer statements

### Error Handling

```go
func SafeMiddleware() server.Middleware {
    return server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
        defer func() {
            if r := recover(); r != nil {
                // Log panic
                log.Errorf("Middleware panic: %v", r)
                
                // Return error instead of panicking
                ctx.JSON(500, map[string]string{
                    "error": "Internal server error",
                })
            }
        }()
        
        // Validate inputs
        if ctx == nil {
            return errors.New("context cannot be nil")
        }
        
        // Execute next middleware
        return next()
    })
}
```

### Security Best Practices

1. **Validate Inputs**: Always validate middleware configuration
2. **Sanitize Headers**: Clean user-provided headers
3. **Rate Limiting**: Prevent abuse with rate limiting
4. **Authentication**: Secure endpoints with proper authentication
5. **Error Messages**: Don't expose sensitive information in errors

## Troubleshooting

### Common Issues

1. **Middleware Not Executing**
   ```go
   // Ensure middleware is registered correctly
   framework.RegisterMiddleware(myMiddleware)
   
   // Check middleware order
   middlewares := chain.GetMiddlewares()
   for i, mw := range middlewares {
       fmt.Printf("Middleware %d: %T\n", i, mw)
   }
   ```

2. **Context Values Not Available**
   ```go
   // Ensure values are set before accessing
   if requestID := ctx.GetString("request_id"); requestID == "" {
       log.Warning("Request ID not found in context")
   }
   ```

3. **Performance Issues**
   ```go
   // Profile middleware performance
   start := time.Now()
   err := next()
   duration := time.Since(start)
   
   if duration > time.Millisecond*100 {
       log.Warningf("Slow middleware: %v", duration)
   }
   ```

### Debug Information

```go
// Middleware debugging
func DebugMiddleware() server.Middleware {
    return server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
        fmt.Printf("Middleware: %s %s\n", ctx.Method(), ctx.Path())
        fmt.Printf("Headers: %+v\n", ctx.Request().Header)
        
        start := time.Now()
        err := next()
        duration := time.Since(start)
        
        fmt.Printf("Duration: %v, Error: %v\n", duration, err)
        return err
    })
}
```

## Next Steps

- **[Server Management](05_server_management.md)** - Learn about server lifecycle management
- **[Integration Examples](06_integration_examples.md)** - See real-world integration patterns
- **[Best Practices](07_best_practices.md)** - Production deployment guidelines
- **[Troubleshooting](08_troubleshooting.md)** - Common issues and solutions
