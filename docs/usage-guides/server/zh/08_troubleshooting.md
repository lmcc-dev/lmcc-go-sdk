# 故障排除指南

本指南涵盖使用 lmcc-go-sdk 服务器模块时的常见问题及其解决方案。

## 常见服务器问题

### 服务器启动失败

#### 端口已被使用

**问题**: 服务器无法绑定到指定端口。

```
Error: listen tcp :8080: bind: address already in use
```

**解决方案**:

1. **检查正在运行的进程**:
   ```bash
   # 查找使用端口的进程 (Find process using the port)
   lsof -i :8080
   netstat -tulpn | grep :8080
   ```

2. **终止冲突进程**:
   ```bash
   kill -9 <PID>
   ```

3. **使用不同端口**:
   ```go
   config := server.DefaultServerConfig()
   config.Port = 8081  // 使用不同端口 (Use different port)
   ```

4. **使用动态端口分配**:
   ```go
   config := server.DefaultServerConfig()
   config.Port = 0  // 让系统分配可用端口 (Let system assign available port)
   ```

#### 配置无效

**问题**: 服务器配置验证失败。

```
Error: invalid configuration: port must be between 1 and 65535
```

**解决方案**:

1. **手动验证配置**:
   ```go
   config := server.DefaultServerConfig()
   if err := config.Validate(); err != nil {
       log.Fatal("Configuration error:", err)
   }
   ```

2. **检查环境变量**:
   ```bash
   echo $SERVER_PORT
   echo $SERVER_HOST
   ```

3. **使用配置默认值**:
   ```go
   config := server.DefaultServerConfig()
   // 只覆盖特定值 (Only override specific values)
   config.Framework = "gin"
   ```

#### 框架插件未找到

**问题**: 指定的框架未注册。

```
Error: framework plugin 'gin' not found
```

**解决方案**:

1. **导入框架插件**:
   ```go
   import _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin"
   ```

2. **检查可用框架**:
   ```go
   frameworks := server.ListFrameworks()
   fmt.Println("Available frameworks:", frameworks)
   ```

3. **手动注册自定义框架**:
   ```go
   err := server.RegisterFramework(myCustomPlugin)
   if err != nil {
       log.Fatal("Failed to register framework:", err)
   }
   ```

### 框架特定问题

#### Gin 问题

**问题**: Gin 路由无法按预期工作。

**解决方案**:

1. **检查路由注册顺序**:
   ```go
   // 在路由之前注册中间件 (Register middleware before routes)
   framework.RegisterMiddleware(middleware.Logger())
   framework.RegisterRoute("GET", "/users", handler)
   ```

2. **验证原生引擎访问**:
   ```go
   ginEngine := framework.GetNativeEngine().(*gin.Engine)
   ginEngine.GET("/native", func(c *gin.Context) {
       c.JSON(200, gin.H{"message": "direct gin route"})
   })
   ```

#### Echo 问题

**问题**: Echo 中间件未正确应用。

**解决方案**:

1. **在正确级别注册中间件**:
   ```go
   // 全局中间件 (Global middleware)
   framework.RegisterMiddleware(middleware.CORS())
   
   // 分组中间件 (Group middleware)
   api := framework.Group("/api")
   api.RegisterMiddleware(middleware.Auth())
   ```

2. **检查 Echo 上下文处理**:
   ```go
   framework.RegisterRoute("GET", "/echo", server.HandlerFunc(func(ctx server.Context) error {
       echoCtx := ctx.(*echo.Context)  // 如果需要，转换为 Echo 上下文 (Cast to Echo context if needed)
       return ctx.JSON(200, map[string]string{"status": "ok"})
   }))
   ```

#### Fiber 问题

**问题**: Fiber 性能问题或意外行为。

**解决方案**:

1. **配置 Fiber 设置**:
   ```go
   config := server.DefaultServerConfig()
   config.Framework = "fiber"
   config.ReadTimeout = 30 * time.Second
   config.WriteTimeout = 30 * time.Second
   ```

2. **检查 Fiber 应用配置**:
   ```go
   fiberApp := framework.GetNativeEngine().(*fiber.App)
   // 配置 Fiber 特定设置 (Configure Fiber-specific settings)
   ```

## 配置问题

### 环境变量问题

**问题**: 环境变量未正确加载。

**解决方案**:

1. **检查变量名**:
   ```bash
   # 正确的变量名 (Correct variable names)
   export SERVER_HOST=0.0.0.0
   export SERVER_PORT=8080
   export SERVER_FRAMEWORK=gin
   ```

2. **验证变量加载**:
   ```go
   fmt.Println("HOST:", os.Getenv("SERVER_HOST"))
   fmt.Println("PORT:", os.Getenv("SERVER_PORT"))
   ```

3. **使用默认值**:
   ```go
   func getEnv(key, defaultValue string) string {
       if value := os.Getenv(key); value != "" {
           return value
       }
       return defaultValue
   }
   
   config.Host = getEnv("SERVER_HOST", "localhost")
   ```

### 配置文件问题

**问题**: YAML 配置未正确解析。

**解决方案**:

1. **验证 YAML 语法**:
   ```bash
   # 检查 YAML 语法 (Check YAML syntax)
   yamllint config.yaml
   ```

2. **检查文件权限**:
   ```bash
   ls -la config.yaml
   chmod 644 config.yaml
   ```

3. **调试配置加载**:
   ```go
   configData, err := os.ReadFile("config.yaml")
   if err != nil {
       log.Fatal("Cannot read config file:", err)
   }
   fmt.Println("Config content:", string(configData))
   ```

## 中间件问题

### 中间件顺序问题

**问题**: 中间件未按预期顺序执行。

**解决方案**: 按正确顺序注册中间件:

```go
// 典型 Web 应用程序的正确顺序 (Correct order for typical web application)
framework.RegisterMiddleware(middleware.Recovery())    // 第一：panic 恢复 (First: panic recovery)
framework.RegisterMiddleware(middleware.Logger())      // 第二：请求日志 (Second: request logging)
framework.RegisterMiddleware(middleware.CORS())        // 第三：CORS 处理 (Third: CORS handling)
framework.RegisterMiddleware(middleware.Auth())        // 第四：身份验证 (Fourth: authentication)
framework.RegisterMiddleware(middleware.RateLimit())   // 第五：限流 (Fifth: rate limiting)
```

### 自定义中间件问题

**问题**: 自定义中间件无法正常工作。

**解决方案**:

1. **检查中间件实现**:
   ```go
   func MyMiddleware() server.Middleware {
       return server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
           // 预处理 (Pre-processing)
           fmt.Println("Before handler")
           
           // 调用下一个中间件/处理器 (Call next middleware/handler)
           err := next()
           
           // 后处理 (Post-processing)
           fmt.Println("After handler")
           
           return err
       })
   }
   ```

2. **验证中间件注册**:
   ```go
   // 在路由之前注册中间件 (Register middleware before routes)
   framework.RegisterMiddleware(MyMiddleware())
   
   // 然后注册路由 (Then register routes)
   framework.RegisterRoute("GET", "/test", handler)
   ```

3. **调试中间件执行**:
   ```go
   func DebugMiddleware() server.Middleware {
       return server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
           log.Printf("Middleware executing for: %s %s", ctx.Method(), ctx.Path())
           return next()
       })
   }
   ```

## 性能问题

### 内存泄漏

**问题**: 服务器内存使用持续增长。

**解决方案**:

1. **监控 goroutines**:
   ```go
   go func() {
       ticker := time.NewTicker(10 * time.Second)
       for range ticker.C {
           fmt.Printf("Goroutines: %d, Memory: %s\n", 
               runtime.NumGoroutine(), 
               formatBytes(getMemUsage()))
       }
   }()
   ```

2. **检查未关闭的资源**:
   ```go
   // 始终关闭 HTTP 响应体 (Always close HTTP response bodies)
   resp, err := http.Get("http://example.com")
   if err != nil {
       return err
   }
   defer resp.Body.Close()
   
   // 始终关闭数据库连接 (Always close database connections)
   rows, err := db.Query("SELECT * FROM users")
   if err != nil {
       return err
   }
   defer rows.Close()
   ```

3. **使用 pprof 进行分析**:
   ```go
   import _ "net/http/pprof"
   
   go func() {
       log.Println(http.ListenAndServe("localhost:6060", nil))
   }()
   ```

### 高 CPU 使用率

**问题**: 服务器消耗过多 CPU。

**解决方案**:

1. **分析 CPU 使用**:
   ```bash
   go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
   ```

2. **优化热点路径**:
   ```go
   // 为数据库使用连接池 (Use connection pooling for databases)
   db.SetMaxOpenConns(25)
   db.SetMaxIdleConns(25)
   db.SetConnMaxLifetime(5 * time.Minute)
   ```

3. **实现缓存**:
   ```go
   // 缓存频繁访问的数据 (Cache frequently accessed data)
   cache := make(map[string]interface{})
   var cacheMutex sync.RWMutex
   
   func getCachedValue(key string) (interface{}, bool) {
       cacheMutex.RLock()
       defer cacheMutex.RUnlock()
       value, exists := cache[key]
       return value, exists
   }
   ```

### 响应时间慢

**问题**: API 响应太慢。

**解决方案**:

1. **添加请求计时中间件**:
   ```go
   func TimingMiddleware() server.Middleware {
       return server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
           start := time.Now()
           err := next()
           duration := time.Since(start)
           
           if duration > 1*time.Second {
               log.Printf("Slow request: %s %s took %v", 
                   ctx.Method(), ctx.Path(), duration)
           }
           
           return err
       })
   }
   ```

2. **优化数据库查询**:
   ```go
   // 使用适当的索引 (Use proper indexes)
   // 限制结果集 (Limit result sets)
   // 使用预处理语句 (Use prepared statements)
   stmt, err := db.Prepare("SELECT * FROM users WHERE id = ?")
   defer stmt.Close()
   ```

3. **实现响应压缩**:
   ```go
   framework.RegisterMiddleware(middleware.Gzip())
   ```

## 网络问题

### 连接超时

**问题**: 请求频繁超时。

**解决方案**:

1. **调整超时设置**:
   ```go
   config := server.DefaultServerConfig()
   config.ReadTimeout = 30 * time.Second
   config.WriteTimeout = 30 * time.Second
   config.IdleTimeout = 60 * time.Second
   ```

2. **配置客户端超时**:
   ```go
   client := &http.Client{
       Timeout: 30 * time.Second,
       Transport: &http.Transport{
           DialTimeout:           5 * time.Second,
           TLSHandshakeTimeout:   5 * time.Second,
           ResponseHeaderTimeout: 10 * time.Second,
       },
   }
   ```

### TLS/SSL 问题

**问题**: HTTPS 配置无法工作。

**解决方案**:

1. **验证证书文件**:
   ```bash
   # 检查证书有效性 (Check certificate validity)
   openssl x509 -in cert.pem -text -noout
   
   # 验证私钥 (Verify private key)
   openssl rsa -in key.pem -check
   ```

2. **正确配置 TLS**:
   ```go
   config := server.DefaultServerConfig()
   config.TLS.Enabled = true
   config.TLS.CertFile = "path/to/cert.pem"
   config.TLS.KeyFile = "path/to/key.pem"
   ```

3. **测试 TLS 连接**:
   ```bash
   # 测试 SSL 连接 (Test SSL connection)
   openssl s_client -connect localhost:8443 -servername localhost
   ```

## 调试技巧

### 启用调试日志

```go
config := server.DefaultServerConfig()
config.Mode = "debug"  // 启用调试模式 (Enable debug mode)

// 添加详细日志中间件 (Add detailed logging middleware)
framework.RegisterMiddleware(func() server.Middleware {
    return server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
        log.Printf("Request: %s %s %s", 
            ctx.Method(), ctx.Path(), ctx.Header("User-Agent"))
        
        err := next()
        
        log.Printf("Response: %d", ctx.Status())
        return err
    })
}())
```

### 健康检查端点

```go
// 添加综合健康检查 (Add comprehensive health check)
framework.RegisterRoute("GET", "/health", server.HandlerFunc(func(ctx server.Context) error {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    health := map[string]interface{}{
        "status":     "healthy",
        "timestamp":  time.Now().Unix(),
        "version":    "1.0.0",
        "uptime":     time.Since(startTime).Seconds(),
        "goroutines": runtime.NumGoroutine(),
        "memory": map[string]interface{}{
            "allocated": formatBytes(m.Alloc),
            "total":     formatBytes(m.TotalAlloc),
            "sys":       formatBytes(m.Sys),
        },
    }
    
    return ctx.JSON(200, health)
}))
```

### 请求跟踪

```go
func TracingMiddleware() server.Middleware {
    return server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
        traceID := generateTraceID()
        ctx.Set("trace_id", traceID)
        
        log.Printf("[%s] Request started: %s %s", 
            traceID, ctx.Method(), ctx.Path())
        
        err := next()
        
        log.Printf("[%s] Request completed: %d", 
            traceID, ctx.Status())
        
        return err
    })
}
```

## 错误恢复

### 优雅降级

```go
func ResilientHandler(ctx server.Context) error {
    // 尝试主要服务 (Try primary service)
    result, err := primaryService.GetData()
    if err != nil {
        log.Printf("Primary service failed: %v", err)
        
        // 回退到缓存 (Fall back to cache)
        result, err = cache.GetData()
        if err != nil {
            log.Printf("Cache failed: %v", err)
            
            // 返回默认数据 (Return default data)
            result = getDefaultData()
        }
    }
    
    return ctx.JSON(200, result)
}
```

### 断路器模式

```go
type CircuitBreaker struct {
    failures    int
    lastFailure time.Time
    threshold   int
    timeout     time.Duration
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    if cb.failures >= cb.threshold {
        if time.Since(cb.lastFailure) < cb.timeout {
            return fmt.Errorf("circuit breaker open")
        }
        cb.failures = 0 // 超时后重置 (Reset after timeout)
    }
    
    err := fn()
    if err != nil {
        cb.failures++
        cb.lastFailure = time.Now()
        return err
    }
    
    cb.failures = 0 // 成功时重置 (Reset on success)
    return nil
}
```

## 常见错误消息

### "framework plugin not found"
- **原因**: 框架插件未导入
- **解决方案**: 添加导入语句: `import _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin"`

### "context canceled"
- **原因**: 请求上下文被取消
- **解决方案**: 检查上下文超时设置并正确处理取消

### "connection refused"
- **原因**: 服务器未运行或地址错误
- **解决方案**: 验证服务器正在运行并检查主机/端口配置

### "method not allowed"
- **原因**: 路由未注册 HTTP 方法
- **解决方案**: 为路由注册所有必需的 HTTP 方法

### "handler not found"
- **原因**: 路由未注册或路径不匹配
- **解决方案**: 检查路由注册和路径模式

## 性能监控

### 要监控的关键指标

1. **响应时间**: 平均值、95th 百分位、99th 百分位
2. **吞吐量**: 每秒请求数
3. **错误率**: 失败请求百分比
4. **资源使用**: CPU、内存、文件描述符
5. **连接统计**: 活跃连接、连接池使用

### 监控工具集成

```go
// Prometheus 指标 (Prometheus metrics)
import "github.com/prometheus/client_golang/prometheus"

var (
    requestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{Name: "http_requests_total"},
        []string{"method", "path", "status"},
    )
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{Name: "http_request_duration_seconds"},
        []string{"method", "path"},
    )
)

func MetricsMiddleware() server.Middleware {
    return server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
        start := time.Now()
        err := next()
        duration := time.Since(start).Seconds()
        
        status := strconv.Itoa(ctx.Status())
        requestsTotal.WithLabelValues(ctx.Method(), ctx.Path(), status).Inc()
        requestDuration.WithLabelValues(ctx.Method(), ctx.Path()).Observe(duration)
        
        return err
    })
}
```

## 获取帮助

### 要收集的调试信息

报告问题时请包含:

1. **服务器配置**
2. **错误消息和堆栈跟踪**
3. **Go 版本和依赖**
4. **操作系统和版本**
5. **网络配置**
6. **资源使用统计**

### 有用的命令

```bash
# 检查 Go 版本 (Check Go version)
go version

# 检查模块依赖 (Check module dependencies)
go mod list -m all

# 检查系统资源 (Check system resources)
top
free -h
df -h

# 检查网络连接 (Check network connections)
netstat -tulpn
ss -tulpn

# 检查服务器日志 (Check server logs)
journalctl -u myserver.service -f
tail -f /var/log/myserver.log
```

## 下一步

- **[API 参考](09_module_specification.md)** - 完整 API 文档