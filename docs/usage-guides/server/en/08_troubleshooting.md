# Troubleshooting Guide

This guide covers common issues and their solutions when using the lmcc-go-sdk server module.

## Common Server Issues

### Server Fails to Start

#### Port Already in Use

**Problem**: Server cannot bind to the specified port.

```
Error: listen tcp :8080: bind: address already in use
```

**Solutions**:

1. **Check for running processes**:
   ```bash
   # Find process using the port
   lsof -i :8080
   netstat -tulpn | grep :8080
   ```

2. **Kill the conflicting process**:
   ```bash
   kill -9 <PID>
   ```

3. **Use a different port**:
   ```go
   config := server.DefaultServerConfig()
   config.Port = 8081  // Use different port
   ```

4. **Use dynamic port allocation**:
   ```go
   config := server.DefaultServerConfig()
   config.Port = 0  // Let system assign available port
   ```

#### Invalid Configuration

**Problem**: Server configuration validation fails.

```
Error: invalid configuration: port must be between 1 and 65535
```

**Solutions**:

1. **Validate configuration manually**:
   ```go
   config := server.DefaultServerConfig()
   if err := config.Validate(); err != nil {
       log.Fatal("Configuration error:", err)
   }
   ```

2. **Check environment variables**:
   ```bash
   echo $SERVER_PORT
   echo $SERVER_HOST
   ```

3. **Use configuration defaults**:
   ```go
   config := server.DefaultServerConfig()
   // Only override specific values
   config.Framework = "gin"
   ```

#### Framework Plugin Not Found

**Problem**: Specified framework is not registered.

```
Error: framework plugin 'gin' not found
```

**Solutions**:

1. **Import framework plugin**:
   ```go
   import _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin"
   ```

2. **Check available frameworks**:
   ```go
   frameworks := server.ListFrameworks()
   fmt.Println("Available frameworks:", frameworks)
   ```

3. **Register custom framework manually**:
   ```go
   err := server.RegisterFramework(myCustomPlugin)
   if err != nil {
       log.Fatal("Failed to register framework:", err)
   }
   ```

### Framework-Specific Issues

#### Gin Issues

**Problem**: Gin routes not working as expected.

**Solutions**:

1. **Check route registration order**:
   ```go
   // Register middleware before routes
   framework.RegisterMiddleware(middleware.Logger())
   framework.RegisterRoute("GET", "/users", handler)
   ```

2. **Verify native engine access**:
   ```go
   ginEngine := framework.GetNativeEngine().(*gin.Engine)
   ginEngine.GET("/native", func(c *gin.Context) {
       c.JSON(200, gin.H{"message": "direct gin route"})
   })
   ```

#### Echo Issues

**Problem**: Echo middleware not applying correctly.

**Solutions**:

1. **Register middleware at correct level**:
   ```go
   // Global middleware
   framework.RegisterMiddleware(middleware.CORS())
   
   // Group middleware
   api := framework.Group("/api")
   api.RegisterMiddleware(middleware.Auth())
   ```

2. **Check Echo context handling**:
   ```go
   framework.RegisterRoute("GET", "/echo", server.HandlerFunc(func(ctx server.Context) error {
       echoCtx := ctx.(*echo.Context)  // Cast to Echo context if needed
       return ctx.JSON(200, map[string]string{"status": "ok"})
   }))
   ```

#### Fiber Issues

**Problem**: Fiber performance issues or unexpected behavior.

**Solutions**:

1. **Configure Fiber settings**:
   ```go
   config := server.DefaultServerConfig()
   config.Framework = "fiber"
   config.ReadTimeout = 30 * time.Second
   config.WriteTimeout = 30 * time.Second
   ```

2. **Check Fiber app configuration**:
   ```go
   fiberApp := framework.GetNativeEngine().(*fiber.App)
   // Configure Fiber-specific settings
   ```

## Configuration Issues

### Environment Variable Problems

**Problem**: Environment variables not loading correctly.

**Solutions**:

1. **Check variable names**:
   ```bash
   # Correct variable names
   export SERVER_HOST=0.0.0.0
   export SERVER_PORT=8080
   export SERVER_FRAMEWORK=gin
   ```

2. **Verify variable loading**:
   ```go
   fmt.Println("HOST:", os.Getenv("SERVER_HOST"))
   fmt.Println("PORT:", os.Getenv("SERVER_PORT"))
   ```

3. **Use default values**:
   ```go
   func getEnv(key, defaultValue string) string {
       if value := os.Getenv(key); value != "" {
           return value
       }
       return defaultValue
   }
   
   config.Host = getEnv("SERVER_HOST", "localhost")
   ```

### Configuration File Issues

**Problem**: YAML configuration not parsing correctly.

**Solutions**:

1. **Validate YAML syntax**:
   ```bash
   # Check YAML syntax
   yamllint config.yaml
   ```

2. **Check file permissions**:
   ```bash
   ls -la config.yaml
   chmod 644 config.yaml
   ```

3. **Debug configuration loading**:
   ```go
   configData, err := os.ReadFile("config.yaml")
   if err != nil {
       log.Fatal("Cannot read config file:", err)
   }
   fmt.Println("Config content:", string(configData))
   ```

## Middleware Issues

### Middleware Order Problems

**Problem**: Middleware not executing in expected order.

**Solution**: Register middleware in correct sequence:

```go
// Correct order for typical web application
framework.RegisterMiddleware(middleware.Recovery())    // First: panic recovery
framework.RegisterMiddleware(middleware.Logger())      // Second: request logging
framework.RegisterMiddleware(middleware.CORS())        // Third: CORS handling
framework.RegisterMiddleware(middleware.Auth())        // Fourth: authentication
framework.RegisterMiddleware(middleware.RateLimit())   // Fifth: rate limiting
```

### Custom Middleware Issues

**Problem**: Custom middleware not working correctly.

**Solutions**:

1. **Check middleware implementation**:
   ```go
   func MyMiddleware() server.Middleware {
       return server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
           // Pre-processing
           fmt.Println("Before handler")
           
           // Call next middleware/handler
           err := next()
           
           // Post-processing
           fmt.Println("After handler")
           
           return err
       })
   }
   ```

2. **Verify middleware registration**:
   ```go
   // Register middleware before routes
   framework.RegisterMiddleware(MyMiddleware())
   
   // Then register routes
   framework.RegisterRoute("GET", "/test", handler)
   ```

3. **Debug middleware execution**:
   ```go
   func DebugMiddleware() server.Middleware {
       return server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
           log.Printf("Middleware executing for: %s %s", ctx.Method(), ctx.Path())
           return next()
       })
   }
   ```

## Performance Issues

### Memory Leaks

**Problem**: Server memory usage keeps increasing.

**Solutions**:

1. **Monitor goroutines**:
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

2. **Check for unclosed resources**:
   ```go
   // Always close HTTP response bodies
   resp, err := http.Get("http://example.com")
   if err != nil {
       return err
   }
   defer resp.Body.Close()
   
   // Always close database connections
   rows, err := db.Query("SELECT * FROM users")
   if err != nil {
       return err
   }
   defer rows.Close()
   ```

3. **Use pprof for profiling**:
   ```go
   import _ "net/http/pprof"
   
   go func() {
       log.Println(http.ListenAndServe("localhost:6060", nil))
   }()
   ```

### High CPU Usage

**Problem**: Server consuming too much CPU.

**Solutions**:

1. **Profile CPU usage**:
   ```bash
   go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
   ```

2. **Optimize hot paths**:
   ```go
   // Use connection pooling for databases
   db.SetMaxOpenConns(25)
   db.SetMaxIdleConns(25)
   db.SetConnMaxLifetime(5 * time.Minute)
   ```

3. **Implement caching**:
   ```go
   // Cache frequently accessed data
   cache := make(map[string]interface{})
   var cacheMutex sync.RWMutex
   
   func getCachedValue(key string) (interface{}, bool) {
       cacheMutex.RLock()
       defer cacheMutex.RUnlock()
       value, exists := cache[key]
       return value, exists
   }
   ```

### Slow Response Times

**Problem**: API responses are too slow.

**Solutions**:

1. **Add request timing middleware**:
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

2. **Optimize database queries**:
   ```go
   // Use proper indexes
   // Limit result sets
   // Use prepared statements
   stmt, err := db.Prepare("SELECT * FROM users WHERE id = ?")
   defer stmt.Close()
   ```

3. **Implement response compression**:
   ```go
   framework.RegisterMiddleware(middleware.Gzip())
   ```

## Networking Issues

### Connection Timeouts

**Problem**: Requests timing out frequently.

**Solutions**:

1. **Adjust timeout settings**:
   ```go
   config := server.DefaultServerConfig()
   config.ReadTimeout = 30 * time.Second
   config.WriteTimeout = 30 * time.Second
   config.IdleTimeout = 60 * time.Second
   ```

2. **Configure client timeouts**:
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

### TLS/SSL Issues

**Problem**: HTTPS configuration not working.

**Solutions**:

1. **Verify certificate files**:
   ```bash
   # Check certificate validity
   openssl x509 -in cert.pem -text -noout
   
   # Verify private key
   openssl rsa -in key.pem -check
   ```

2. **Configure TLS properly**:
   ```go
   config := server.DefaultServerConfig()
   config.TLS.Enabled = true
   config.TLS.CertFile = "path/to/cert.pem"
   config.TLS.KeyFile = "path/to/key.pem"
   ```

3. **Test TLS connection**:
   ```bash
   # Test SSL connection
   openssl s_client -connect localhost:8443 -servername localhost
   ```

## Debugging Techniques

### Enable Debug Logging

```go
config := server.DefaultServerConfig()
config.Mode = "debug"  // Enable debug mode

// Add detailed logging middleware
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

### Health Check Endpoints

```go
// Add comprehensive health check
framework.RegisterRoute("GET", "/health", server.HandlerFunc(func(ctx server.Context) error {
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

### Request Tracing

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

## Error Recovery

### Graceful Degradation

```go
func ResilientHandler(ctx server.Context) error {
    // Try primary service
    result, err := primaryService.GetData()
    if err != nil {
        log.Printf("Primary service failed: %v", err)
        
        // Fall back to cache
        result, err = cache.GetData()
        if err != nil {
            log.Printf("Cache failed: %v", err)
            
            // Return default data
            result = getDefaultData()
        }
    }
    
    return ctx.JSON(200, result)
}
```

### Circuit Breaker Pattern

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
        cb.failures = 0 // Reset after timeout
    }
    
    err := fn()
    if err != nil {
        cb.failures++
        cb.lastFailure = time.Now()
        return err
    }
    
    cb.failures = 0 // Reset on success
    return nil
}
```

## Common Error Messages

### "framework plugin not found"
- **Cause**: Framework plugin not imported
- **Solution**: Add import statement: `import _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin"`

### "context canceled"
- **Cause**: Request context was canceled
- **Solution**: Check context timeout settings and handle cancellation properly

### "connection refused"
- **Cause**: Server not running or wrong address
- **Solution**: Verify server is running and check host/port configuration

### "method not allowed"
- **Cause**: HTTP method not registered for route
- **Solution**: Register all required HTTP methods for the route

### "handler not found"
- **Cause**: Route not registered or path mismatch
- **Solution**: Check route registration and path patterns

## Performance Monitoring

### Key Metrics to Monitor

1. **Response Time**: Average, 95th percentile, 99th percentile
2. **Throughput**: Requests per second
3. **Error Rate**: Percentage of failed requests
4. **Resource Usage**: CPU, memory, file descriptors
5. **Connection Stats**: Active connections, connection pool usage

### Monitoring Tools Integration

```go
// Prometheus metrics
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

## Getting Help

### Debug Information to Collect

When reporting issues, include:

1. **Server configuration**
2. **Error messages and stack traces**
3. **Go version and dependencies**
4. **Operating system and version**
5. **Network configuration**
6. **Resource usage statistics**

### Useful Commands

```bash
# Check Go version
go version

# Check module dependencies
go mod list -m all

# Check system resources
top
free -h
df -h

# Check network connections
netstat -tulpn
ss -tulpn

# Check server logs
journalctl -u myserver.service -f
tail -f /var/log/myserver.log
```

## Next Steps

- **[API Reference](09_api_reference.md)** - Complete API documentation
