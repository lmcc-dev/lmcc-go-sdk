# 中间件系统指南

lmcc-go-sdk 服务器模块提供了全面的中间件系统，支持以模块化和可组合的方式处理请求/响应。本指南涵盖统一中间件接口、内置中间件组件和自定义中间件开发。

## 中间件架构

### 核心概念

中间件系统围绕几个关键概念构建：

1. **统一接口**: 所有中间件在各框架中实现相同接口
2. **链式处理**: 中间件按定义顺序执行
3. **上下文传递**: 请求上下文在中间件链中流动
4. **框架无关**: 中间件适用于任何支持的框架
5. **服务集成**: 中间件与其他 SDK 模块集成

### 中间件接口

所有中间件实现统一的 `Middleware` 接口：

```go
type Middleware interface {
    Process(ctx Context, next func() error) error
}

// 便利的函数类型 (Function type for convenience)
type MiddlewareFunc func(ctx Context, next func() error) error

func (f MiddlewareFunc) Process(ctx Context, next func() error) error {
    return f(ctx, next)
}
```

### 中间件链执行

```
请求 (Request)
   ↓
┌─────────────────┐
│ 日志记录器       │ ← 记录请求详情 (Log request details)
├─────────────────┤
│ 恢复中间件       │ ← 从 panic 恢复 (Recover from panics)
├─────────────────┤
│ CORS           │ ← 处理跨域请求 (Handle cross-origin requests)
├─────────────────┤
│ 安全中间件       │ ← 添加安全头 (Add security headers)
├─────────────────┤
│ 压缩中间件       │ ← 压缩响应 (Compress responses)
├─────────────────┤
│ 身份验证        │ ← 验证用户凭据 (Validate user credentials)
├─────────────────┤
│ 限流中间件       │ ← 控制请求速率 (Control request rate)
├─────────────────┤
│ 指标收集        │ ← 收集性能指标 (Collect performance metrics)
├─────────────────┤
│ 您的处理器       │ ← 处理业务逻辑 (Process business logic)
└─────────────────┘
   ↓
响应 (Response)
```

> **⚠️ 实现状态提醒**: 中间件系统当前正在积极开发中。使用前请检查下方的实现状态。

## 实现状态

| 中间件 | 状态 | Gin | Echo | Fiber | 备注 |
|-------|------|-----|------|-------|------|
| **日志中间件** | ⚠️ 部分完成 | ❌ | ✅ | ✅ | Gin实现待开发 |
| **恢复中间件** | ⚠️ 部分完成 | ❌ | ✅ | ✅ | Gin实现待开发 |
| **CORS中间件** | ✅ 完整实现 | ✅ | ✅ | ✅ | 完全实现 |
| **限流中间件** | 🚧 开发中 | ❌ | ❌ | ❌ | 仅接口定义 |
| **认证中间件** | 🚧 开发中 | ❌ | ❌ | ❌ | 仅接口定义 |
| **安全中间件** | 🚧 开发中 | ❌ | ❌ | ❌ | 仅接口定义 |
| **压缩中间件** | 🚧 开发中 | ❌ | ❌ | ❌ | 仅接口定义 |
| **指标中间件** | 🚧 开发中 | ❌ | ❌ | ❌ | 仅接口定义 |

### 当前可用 (完整实现)
- **CORS中间件**: 完整的跨框架实现
- **日志中间件**: 可用于Echo和Fiber
- **恢复中间件**: 可用于Echo和Fiber

### 即将推出 (开发中)
以下中间件计划在后续版本中实现：
- 安全头中间件 (预计: 2025-01-12)
- 限流中间件 (预计: 2025-01-13)
- 认证中间件 (预计: 2025-01-14)
- 压缩中间件 (预计: 2025-01-14)
- 指标收集中间件 (预计: 2025-01-15)

## 内置中间件

### 日志中间件

日志中间件记录请求和响应信息。

#### 配置

```go
config.Middleware.Logger = server.LoggerMiddlewareConfig{
    Enabled:     true,
    Format:      "json",        // json 或 text (json or text)
    SkipPaths:   []string{"/health", "/favicon.ico"},
    IncludeBody: false,         // 包含请求/响应体 (Include request/response body)
    MaxBodySize: 1024,          // 记录的最大体大小（字节） (Maximum body size to log in bytes)
}
```

#### 使用示例

```go
// 基本日志中间件 (Basic logger middleware)
logger := middleware.NewLogger(server.LoggerMiddlewareConfig{
    Enabled: true,
    Format:  "json",
})

framework.RegisterMiddleware(logger)
```

#### 自定义日志配置

```go
// 开发详细日志 (Detailed logging for development)
logger := middleware.NewLogger(server.LoggerMiddlewareConfig{
    Enabled:     true,
    Format:      "text",
    SkipPaths:   []string{},     // 记录所有请求 (Log all requests)
    IncludeBody: true,           // 包含请求体 (Include request bodies)
    MaxBodySize: 4096,           // 记录最多 4KB (Log up to 4KB)
})

// 生产日志 (Production logging)
logger := middleware.NewLogger(server.LoggerMiddlewareConfig{
    Enabled:     true,
    Format:      "json",
    SkipPaths:   []string{"/health", "/metrics", "/static/*"},
    IncludeBody: false,          // 不记录请求体 (Don't log request bodies)
    MaxBodySize: 0,              // 不记录体 (No body logging)
})
```

#### 日志输出格式

**JSON 格式:**
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

**文本格式:**
```
2024-01-09 10:30:00 | 200 | 15ms | 192.168.1.100 | GET /api/users
```

### 恢复中间件

恢复中间件捕获 panic 并将其转换为适当的 HTTP 错误响应。

#### 配置

```go
config.Middleware.Recovery = server.RecoveryMiddlewareConfig{
    Enabled:             true,
    PrintStack:          true,   // 打印堆栈跟踪 (Print stack trace)
    DisableStackAll:     false,  // 禁用所有堆栈跟踪 (Disable all stack traces)
    DisableColorConsole: false,  // 禁用彩色输出 (Disable colored output)
}
```

#### 使用示例

```go
// 基本恢复中间件 (Basic recovery middleware)
recovery := middleware.NewRecovery(server.RecoveryMiddlewareConfig{
    Enabled:    true,
    PrintStack: true,
})

framework.RegisterMiddleware(recovery)
```

#### 自定义恢复处理器

```go
// 带错误报告的自定义恢复 (Custom recovery with error reporting)
recovery := middleware.NewRecovery(server.RecoveryMiddlewareConfig{
    Enabled:    true,
    PrintStack: false,  // 不打印到控制台 (Don't print to console)
})

// 设置自定义恢复处理器 (Set custom recovery handler)
recovery.SetCustomRecoveryHandler(func(ctx server.Context, err interface{}) {
    // 记录错误到监控服务 (Log error to monitoring service)
    log.Errorf("Panic recovered: %v", err)
    
    // 发送到错误跟踪服务 (Send to error tracking service)
    errorTracker.ReportPanic(ctx.Request().Context(), err)
    
    // 返回自定义错误响应 (Return custom error response)
    ctx.JSON(500, map[string]string{
        "error": "Internal server error",
        "request_id": ctx.GetString("request_id"),
    })
})

framework.RegisterMiddleware(recovery)
```

#### 生产 vs 开发配置

```go
// 开发配置 (Development configuration)
if config.IsDebugMode() {
    recovery := middleware.NewRecovery(server.RecoveryMiddlewareConfig{
        Enabled:             true,
        PrintStack:          true,   // 显示完整堆栈跟踪 (Show full stack traces)
        DisableColorConsole: false,  // 彩色输出 (Colored output)
    })
} else {
    // 生产配置 (Production configuration)
    recovery := middleware.NewRecovery(server.RecoveryMiddlewareConfig{
        Enabled:             true,
        PrintStack:          false,  // 日志中无堆栈跟踪 (No stack traces in logs)
        DisableStackAll:     true,   // 完全禁用堆栈跟踪 (Completely disable stack traces)
        DisableColorConsole: true,   // 无彩色输出 (No colored output)
    })
}
```

### CORS 中间件

CORS 中间件处理跨域资源共享请求。

#### 配置

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

#### 使用示例

```go
// API 基本 CORS (Basic CORS for API)
cors := middleware.NewCORS(server.CORSConfig{
    Enabled:      true,
    AllowOrigins: []string{"https://myapp.com", "https://admin.myapp.com"},
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders: []string{"Content-Type", "Authorization"},
})

framework.RegisterMiddleware(cors)
```

#### 开发 CORS（允许所有）

```go
// 开发 CORS - 允许一切 (Development CORS - allow everything)
cors := middleware.NewCORS(server.CORSConfig{
    Enabled:      true,
    AllowOrigins: []string{"*"},
    AllowMethods: []string{"*"},
    AllowHeaders: []string{"*"},
    AllowCredentials: false,  // 当 origins 是 "*" 时必须为 false (Must be false when origins is "*")
})
```

#### 限制性生产 CORS

```go
// 生产 CORS - 限制性设置 (Production CORS - restrictive settings)
cors := middleware.NewCORS(server.CORSConfig{
    Enabled:          true,
    AllowOrigins:     []string{"https://yourdomain.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Content-Type", "Authorization", "X-Requested-With"},
    ExposeHeaders:    []string{"X-Total-Count", "X-Page-Count"},
    AllowCredentials: true,
    MaxAge:          1 * time.Hour,  // 缓存预检请求 1 小时 (Cache preflight for 1 hour)
})
```

#### 动态 CORS 配置

```go
// 基于请求的动态 CORS (Dynamic CORS based on request)
cors := middleware.NewCORS(server.CORSConfig{
    Enabled: true,
})

// 设置动态源验证 (Set dynamic origin validation)
cors.SetOriginValidator(func(origin string) bool {
    // 检查数据库或配置 (Check against database or configuration)
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

### 限流中间件

限流中间件控制传入请求的速率。

#### 配置

```go
config.Middleware.RateLimit = server.RateLimitMiddlewareConfig{
    Enabled: true,
    Rate:    100.0,      // 每秒 100 个请求 (100 requests per second)
    Burst:   200,        // 允许突发 200 个请求 (Allow burst of 200 requests)
    KeyFunc: "ip",       // 按 IP 地址限流 (Rate limit by IP address)
}
```

#### 使用示例

```go
// 按 IP 基本限流 (Basic rate limiting by IP)
rateLimit := middleware.NewRateLimit(server.RateLimitMiddlewareConfig{
    Enabled: true,
    Rate:    50.0,   // 每秒 50 个请求 (50 requests per second)
    Burst:   100,    // 允许突发 100 个请求 (Allow burst of 100 requests)
    KeyFunc: "ip",
})

framework.RegisterMiddleware(rateLimit)
```

#### 自定义限流键

```go
// 按用户 ID 限流 (Rate limit by user ID)
rateLimit := middleware.NewRateLimit(server.RateLimitMiddlewareConfig{
    Enabled: true,
    Rate:    10.0,   // 每用户每秒 10 个请求 (10 requests per second per user)
    Burst:   20,     // 允许突发 20 个请求 (Allow burst of 20 requests)
})

// 设置自定义键函数 (Set custom key function)
rateLimit.SetKeyFunc(func(ctx server.Context) string {
    userID := ctx.GetString("user_id")
    if userID == "" {
        return ctx.ClientIP()  // 回退到 IP (Fallback to IP)
    }
    return "user:" + userID
})

framework.RegisterMiddleware(rateLimit)
```

#### 不同端点的不同限制

```go
// API 限流 (API rate limiting)
apiRateLimit := middleware.NewRateLimit(server.RateLimitMiddlewareConfig{
    Enabled: true,
    Rate:    100.0,
    Burst:   200,
})

// 上传限流（更严格） (Upload rate limiting - more restrictive)
uploadRateLimit := middleware.NewRateLimit(server.RateLimitMiddlewareConfig{
    Enabled: true,
    Rate:    5.0,    // 每秒 5 次上传 (5 uploads per second)
    Burst:   10,     // 允许突发 10 次上传 (Allow burst of 10 uploads)
})

// 对不同路由组应用不同限流 (Apply different rate limits to different route groups)
api := framework.Group("/api")
api.RegisterMiddleware(apiRateLimit)

upload := framework.Group("/upload")
upload.RegisterMiddleware(uploadRateLimit)
```

### 身份验证中间件

身份验证中间件验证用户凭据并管理用户会话。

#### 配置

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

#### JWT 身份验证

```go
// JWT 身份验证中间件 (JWT authentication middleware)
auth := middleware.NewAuth(server.AuthMiddlewareConfig{
    Enabled: true,
    Type:    "jwt",
    SkipPaths: []string{"/auth/login", "/auth/register", "/public/*"},
})

// 设置 JWT 验证函数 (Set JWT validation function)
auth.SetAuthFunc(func(ctx server.Context) (interface{}, error) {
    token := ctx.Header("Authorization")
    if token == "" {
        return nil, errors.New("missing authorization header")
    }
    
    // 移除 "Bearer " 前缀 (Remove "Bearer " prefix)
    if strings.HasPrefix(token, "Bearer ") {
        token = token[7:]
    }
    
    // 验证 JWT 令牌 (Validate JWT token)
    claims, err := validateJWTToken(token)
    if err != nil {
        return nil, err
    }
    
    return claims, nil
})

framework.RegisterMiddleware(auth)
```

#### 基本身份验证

```go
// 基本身份验证中间件 (Basic authentication middleware)
auth := middleware.NewAuth(server.AuthMiddlewareConfig{
    Enabled: true,
    Type:    "basic",
})

auth.SetAuthFunc(func(ctx server.Context) (interface{}, error) {
    username, password, ok := ctx.Request().BasicAuth()
    if !ok {
        return nil, errors.New("missing basic auth")
    }
    
    // 验证凭据 (Validate credentials)
    user, err := validateCredentials(username, password)
    if err != nil {
        return nil, err
    }
    
    return user, nil
})

framework.RegisterMiddleware(auth)
```

#### 自定义未授权处理器

```go
auth := middleware.NewAuth(server.AuthMiddlewareConfig{
    Enabled: true,
    Type:    "custom",
})

// 设置自定义未授权响应 (Set custom unauthorized response)
auth.SetUnauthorizedHandler(func(ctx server.Context) {
    ctx.JSON(401, map[string]interface{}{
        "error": "Unauthorized",
        "code":  "AUTH_REQUIRED",
        "message": "Please provide valid authentication credentials",
    })
})
```

### 安全中间件

安全中间件添加安全头以防止常见漏洞。

#### 基本安全头

```go
// 基本安全中间件 (Basic security middleware)
security := middleware.NewSecurity()

security.SetXSSProtection(true)               // X-XSS-Protection: 1; mode=block
security.SetContentTypeNosniff(true)          // X-Content-Type-Options: nosniff
security.SetFrameOptions("DENY")              // X-Frame-Options: DENY
security.SetHSTSMaxAge(31536000)              // Strict-Transport-Security: max-age=31536000

framework.RegisterMiddleware(security)
```

#### 全面安全配置

```go
security := middleware.NewSecurity()

// XSS 保护 (XSS Protection)
security.SetXSSProtection(true)

// 内容类型选项 (Content Type Options)
security.SetContentTypeNosniff(true)

// 框架选项 (Frame Options)
security.SetFrameOptions("DENY")

// HSTS（HTTP 严格传输安全） (HSTS - HTTP Strict Transport Security)
security.SetHSTSMaxAge(31536000)  // 1 年 (1 year)
security.SetHSTSIncludeSubdomains(true)

// 内容安全策略 (Content Security Policy)
security.SetContentSecurityPolicy("default-src 'self'; script-src 'self' 'unsafe-inline'")

// 引用者策略 (Referrer Policy)
security.SetReferrerPolicy("strict-origin-when-cross-origin")

// 权限策略 (Permissions Policy)
security.SetPermissionsPolicy("geolocation=(), microphone=(), camera=()")

framework.RegisterMiddleware(security)
```

### 压缩中间件

压缩中间件压缩 HTTP 响应以减少带宽使用。

#### 基本压缩

```go
// 基本压缩中间件 (Basic compression middleware)
compression := middleware.NewCompression()

compression.SetLevel(6)           // 压缩级别（1-9） (Compression level 1-9)
compression.SetMinLength(1024)    // 压缩的最小响应大小 (Minimum response size to compress)
compression.SetExcludedExtensions([]string{".jpg", ".png", ".gif", ".zip"})

framework.RegisterMiddleware(compression)
```

#### 高级压缩配置

```go
compression := middleware.NewCompression()

// 压缩设置 (Compression settings)
compression.SetLevel(6)                    // 速度和压缩之间的平衡 (Balance between speed and compression)
compression.SetMinLength(1024)             // 只压缩 > 1KB 的响应 (Only compress responses > 1KB)

// 排除已压缩的内容 (Exclude already compressed content)
compression.SetExcludedExtensions([]string{
    ".jpg", ".jpeg", ".png", ".gif", ".webp",  // 图片 (Images)
    ".mp4", ".avi", ".mov",                    // 视频 (Videos)
    ".zip", ".rar", ".7z", ".tar.gz",          // 档案 (Archives)
    ".pdf",                                    // 文档 (Documents)
})

// 排除特定内容类型 (Exclude specific content types)
compression.SetExcludedContentTypes([]string{
    "image/*",
    "video/*",
    "application/zip",
    "application/pdf",
})

framework.RegisterMiddleware(compression)
```

### 指标中间件

指标中间件收集性能和使用指标。

#### 基本指标收集

```go
// 基本指标中间件 (Basic metrics middleware)
metrics := middleware.NewMetrics()

metrics.SetMetricsPath("/metrics")  // Prometheus 指标端点 (Prometheus metrics endpoint)
metrics.SetSkipPaths([]string{"/health", "/metrics"})

framework.RegisterMiddleware(metrics)
```

#### 自定义指标配置

```go
metrics := middleware.NewMetrics()

// 指标配置 (Metrics configuration)
metrics.SetMetricsPath("/metrics")
metrics.SetSkipPaths([]string{"/health", "/metrics", "/favicon.ico"})

// 自定义标签 (Custom labels)
metrics.SetCustomLabels(map[string]string{
    "service": "api-server",
    "version": "1.0.0",
    "environment": "production",
})

// 自定义指标收集器 (Custom metric collectors)
metrics.AddHistogram("request_duration_seconds", "HTTP request duration", []string{"method", "endpoint", "status"})
metrics.AddCounter("request_total", "Total HTTP requests", []string{"method", "endpoint", "status"})
metrics.AddGauge("active_connections", "Active HTTP connections", []string{})

framework.RegisterMiddleware(metrics)
```

## 中间件链管理

### 创建中间件链

```go
// 创建中间件链 (Create middleware chain)
chain := middleware.NewMiddlewareChain()

// 按顺序添加中间件 (Add middleware in order)
chain.Add(middleware.NewLogger(loggerConfig))
chain.Add(middleware.NewRecovery(recoveryConfig))
chain.Add(middleware.NewCORS(corsConfig))
chain.Add(middleware.NewSecurity())
chain.Add(middleware.NewAuth(authConfig))

// 向框架注册链 (Register chain with framework)
for _, mw := range chain.GetMiddlewares() {
    framework.RegisterMiddleware(mw)
}
```

### 条件中间件

```go
// 有条件地添加中间件 (Add middleware conditionally)
if config.IsDebugMode() {
    // 开发中间件 (Development middleware)
    chain.Add(middleware.NewDetailedLogger())
    chain.Add(middleware.NewDebugHeaders())
} else {
    // 生产中间件 (Production middleware)
    chain.Add(middleware.NewProductionLogger())
    chain.Add(middleware.NewSecurity())
}

// 始终添加基本中间件 (Always add essential middleware)
chain.Add(middleware.NewRecovery(recoveryConfig))
chain.Add(middleware.NewCORS(corsConfig))
```

### 路由特定中间件

```go
// 全局中间件 (Global middleware)
framework.RegisterMiddleware(middleware.NewLogger(loggerConfig))
framework.RegisterMiddleware(middleware.NewRecovery(recoveryConfig))

// API 特定中间件 (API-specific middleware)
api := framework.Group("/api")
api.RegisterMiddleware(middleware.NewAuth(authConfig))
api.RegisterMiddleware(middleware.NewRateLimit(rateLimitConfig))

// 管理员特定中间件 (Admin-specific middleware)
admin := framework.Group("/admin")
admin.RegisterMiddleware(middleware.NewAuth(adminAuthConfig))
admin.RegisterMiddleware(middleware.NewAuditLog())

// 公共路由（无额外中间件） (Public routes - no additional middleware)
public := framework.Group("/public")
```

## 自定义中间件开发

### 简单自定义中间件

```go
// 请求 ID 中间件 (Request ID middleware)
func RequestIDMiddleware() server.Middleware {
    return server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
        // 生成唯一请求 ID (Generate unique request ID)
        requestID := generateRequestID()
        
        // 存储在上下文中 (Store in context)
        ctx.Set("request_id", requestID)
        
        // 添加到响应头 (Add to response header)
        ctx.SetHeader("X-Request-ID", requestID)
        
        // 继续到下一个中间件 (Continue to next middleware)
        return next()
    })
}

// 使用 (Usage)
framework.RegisterMiddleware(RequestIDMiddleware())
```

### 高级自定义中间件

```go
// 自定义审计日志中间件 (Custom audit logging middleware)
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
    // 跳过某些路径的审计日志 (Skip audit logging for certain paths)
    path := ctx.Path()
    for _, skipPath := range am.skipPaths {
        if path == skipPath {
            return next()
        }
    }
    
    // 记录审计日志条目 (Record audit log entry)
    start := time.Now()
    userID := ctx.GetString("user_id")
    
    // 执行下一个中间件/处理器 (Execute next middleware/handler)
    err := next()
    
    // 记录审计条目 (Log audit entry)
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

// 使用 (Usage)
auditMiddleware := NewAuditMiddleware(auditLogger)
framework.RegisterMiddleware(auditMiddleware)
```

## 最佳实践

### 中间件排序

根据其目的和依赖关系排序中间件：

1. **日志记录器** - 首先记录所有请求
2. **恢复** - 早期捕获所有 panic
3. **CORS** - 在身份验证之前处理预检请求
4. **安全** - 早期添加安全头
5. **压缩** - 在内容生成之前
6. **身份验证** - 在业务逻辑之前
7. **限流** - 在身份验证之后
8. **指标** - 监控已认证请求
9. **业务逻辑** - 您的处理器

### 性能考虑

1. **最小化分配**: 避免在热路径中创建不必要的对象
2. **早期返回**: 尽可能从中间件早期返回
3. **高效字符串操作**: 使用字符串构建器进行连接
4. **上下文管理**: 不要在上下文中存储大对象
5. **资源清理**: 始终在 defer 语句中清理资源

### 错误处理

```go
func SafeMiddleware() server.Middleware {
    return server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
        defer func() {
            if r := recover(); r != nil {
                // 记录 panic (Log panic)
                log.Errorf("Middleware panic: %v", r)
                
                // 返回错误而不是 panic (Return error instead of panicking)
                ctx.JSON(500, map[string]string{
                    "error": "Internal server error",
                })
            }
        }()
        
        // 验证输入 (Validate inputs)
        if ctx == nil {
            return errors.New("context cannot be nil")
        }
        
        // 执行下一个中间件 (Execute next middleware)
        return next()
    })
}
```

### 安全最佳实践

1. **验证输入**: 始终验证中间件配置
2. **清理头部**: 清理用户提供的头部
3. **限流**: 通过限流防止滥用
4. **身份验证**: 使用适当的身份验证保护端点
5. **错误消息**: 不要在错误中暴露敏感信息

## 故障排除

### 常见问题

1. **中间件未执行**
   ```go
   // 确保中间件正确注册 (Ensure middleware is registered correctly)
   framework.RegisterMiddleware(myMiddleware)
   
   // 检查中间件顺序 (Check middleware order)
   middlewares := chain.GetMiddlewares()
   for i, mw := range middlewares {
       fmt.Printf("Middleware %d: %T\n", i, mw)
   }
   ```

2. **上下文值不可用**
   ```go
   // 确保在访问之前设置值 (Ensure values are set before accessing)
   if requestID := ctx.GetString("request_id"); requestID == "" {
       log.Warning("Request ID not found in context")
   }
   ```

3. **性能问题**
   ```go
   // 分析中间件性能 (Profile middleware performance)
   start := time.Now()
   err := next()
   duration := time.Since(start)
   
   if duration > time.Millisecond*100 {
       log.Warningf("Slow middleware: %v", duration)
   }
   ```

### 调试信息

```go
// 中间件调试 (Middleware debugging)
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

## 下一步

- **[服务器管理](05_server_management.md)** - 了解服务器生命周期管理
- **[集成示例](06_integration_examples.md)** - 查看真实世界集成模式
- **[最佳实践](07_best_practices.md)** - 生产部署指南
- **[故障排除](08_troubleshooting.md)** - 常见问题和解决方案