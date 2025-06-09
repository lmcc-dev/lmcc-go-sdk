# Server Management Guide

The lmcc-go-sdk server module provides comprehensive server lifecycle management through the `ServerManager` component. This guide covers server creation, startup, shutdown, monitoring, and best practices for production deployment.

## Server Management Overview

### Core Components

The server management system consists of:

1. **ServerManager**: Core component managing server lifecycle
2. **ServerFactory**: Factory for creating server instances
3. **Graceful Shutdown**: Clean shutdown handling with signal management
4. **Health Monitoring**: Server state monitoring and reporting
5. **Configuration Management**: Runtime configuration and validation

### Server Lifecycle

```
Creation → Configuration → Startup → Running → Shutdown → Cleanup
    ↓           ↓           ↓         ↓         ↓         ↓
   Factory   Validation   Listen   Serving   Graceful   Done
```

## Server Creation

### Using ServerManager

```go
package main

import (
    "context"
    "log"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin"
)

func main() {
    // Create configuration
    config := server.DefaultServerConfig()
    config.Framework = "gin"
    config.Host = "0.0.0.0"
    config.Port = 8080
    
    // Create server manager
    manager, err := server.CreateServerManager("gin", config)
    if err != nil {
        log.Fatal("Failed to create server:", err)
    }
    
    // Configure routes and middleware
    framework := manager.GetFramework()
    framework.RegisterRoute("GET", "/health", server.HandlerFunc(healthHandler))
    
    // Start server
    ctx := context.Background()
    if err := manager.Start(ctx); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}

func healthHandler(ctx server.Context) error {
    return ctx.JSON(200, map[string]interface{}{
        "status": "healthy",
        "timestamp": time.Now().Unix(),
    })
}
```

### Using ServerFactory

```go
// Create factory
factory := server.NewServerFactory()

// Create server using factory
manager, err := factory.CreateServer("gin", config)
if err != nil {
    log.Fatal("Failed to create server:", err)
}

// Register custom plugins with factory
customPlugin := myplugin.NewPlugin()
err = factory.RegisterPlugin(customPlugin)
if err != nil {
    log.Fatal("Failed to register plugin:", err)
}
```

### Configuration Validation

```go
// Create configuration
config := &server.ServerConfig{
    Framework: "gin",
    Host:     "0.0.0.0",
    Port:     8080,
    Mode:     "production",
}

// Validate configuration before creating server
if err := config.Validate(); err != nil {
    log.Fatal("Invalid configuration:", err)
}

// Create server with validated configuration
manager, err := server.CreateServerManager("gin", config)
if err != nil {
    log.Fatal("Failed to create server:", err)
}
```

## Server Startup

### Basic Startup

```go
func startServer(manager *server.ServerManager) error {
    ctx := context.Background()
    
    // Start server (blocking call)
    if err := manager.Start(ctx); err != nil {
        return fmt.Errorf("server startup failed: %w", err)
    }
    
    return nil
}
```

### Non-blocking Startup

```go
func startServerAsync(manager *server.ServerManager) error {
    ctx := context.Background()
    
    // Start server in goroutine
    go func() {
        if err := manager.Start(ctx); err != nil {
            log.Printf("Server error: %v", err)
        }
    }()
    
    // Wait for server to be ready
    if err := waitForServerReady(manager); err != nil {
        return fmt.Errorf("server not ready: %w", err)
    }
    
    return nil
}

func waitForServerReady(manager *server.ServerManager) error {
    timeout := time.After(30 * time.Second)
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()
    
    for {
        select {
        case <-timeout:
            return fmt.Errorf("timeout waiting for server to be ready")
        case <-ticker.C:
            if manager.IsRunning() {
                return nil
            }
        }
    }
}
```

### Startup with Context

```go
func startServerWithTimeout(manager *server.ServerManager) error {
    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Start server with timeout
    if err := manager.Start(ctx); err != nil {
        return fmt.Errorf("server startup failed: %w", err)
    }
    
    return nil
}
```

### Error Handling During Startup

```go
func startServerWithRetry(manager *server.ServerManager, maxRetries int) error {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        ctx := context.Background()
        
        if err := manager.Start(ctx); err != nil {
            lastErr = err
            log.Printf("Startup attempt %d failed: %v", i+1, err)
            
            // Wait before retry
            time.Sleep(time.Duration(i+1) * time.Second)
            continue
        }
        
        return nil // Success
    }
    
    return fmt.Errorf("failed to start server after %d attempts: %w", maxRetries, lastErr)
}
```

## Server Shutdown

### Graceful Shutdown Configuration

```go
config := &server.ServerConfig{
    Framework: "gin",
    Host:     "0.0.0.0",
    Port:     8080,
    
    GracefulShutdown: server.GracefulShutdownConfig{
        Enabled:  true,
        Timeout:  30 * time.Second,  // Maximum time to wait for shutdown
        WaitTime: 5 * time.Second,   // Time to wait for ongoing requests
    },
}
```

### Manual Shutdown

```go
func shutdownServer(manager *server.ServerManager) error {
    if !manager.IsRunning() {
        return fmt.Errorf("server is not running")
    }
    
    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Stop server
    if err := manager.Stop(ctx); err != nil {
        return fmt.Errorf("failed to stop server: %w", err)
    }
    
    log.Println("Server stopped successfully")
    return nil
}
```

### Signal-based Shutdown

```go
func startServerWithSignalHandling(manager *server.ServerManager) error {
    // Create signal channel
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    // Start server in goroutine
    serverErrChan := make(chan error, 1)
    go func() {
        ctx := context.Background()
        serverErrChan <- manager.Start(ctx)
    }()
    
    // Wait for signal or server error
    select {
    case sig := <-sigChan:
        log.Printf("Received signal %v, shutting down server...", sig)
        
        // Create shutdown context
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        
        // Stop server
        if err := manager.Stop(ctx); err != nil {
            log.Printf("Error during shutdown: %v", err)
            return err
        }
        
        log.Println("Server shutdown completed")
        return nil
        
    case err := <-serverErrChan:
        if err != nil {
            log.Printf("Server error: %v", err)
            return err
        }
        return nil
    }
}
```

### Automatic Graceful Shutdown

The server manager includes built-in graceful shutdown when enabled:

```go
config := &server.ServerConfig{
    Framework: "gin",
    Host:     "0.0.0.0",
    Port:     8080,
    
    GracefulShutdown: server.GracefulShutdownConfig{
        Enabled:  true,                // Enable automatic signal handling
        Timeout:  30 * time.Second,    // Shutdown timeout
        WaitTime: 5 * time.Second,     // Wait for ongoing requests
    },
}

manager, err := server.CreateServerManager("gin", config)
if err != nil {
    log.Fatal("Failed to create server:", err)
}

// The server will automatically handle SIGINT and SIGTERM
// No additional signal handling needed
if err := manager.Start(context.Background()); err != nil {
    log.Fatal("Server failed:", err)
}
```

## Server Monitoring

### Health Checks

```go
// Built-in health check
func setupHealthCheck(framework server.WebFramework) {
    framework.RegisterRoute("GET", "/health", server.HandlerFunc(func(ctx server.Context) error {
        return ctx.JSON(200, map[string]interface{}{
            "status":    "healthy",
            "timestamp": time.Now().Unix(),
            "uptime":    time.Since(startTime).Seconds(),
        })
    }))
}

// Detailed health check
func setupDetailedHealthCheck(framework server.WebFramework, manager *server.ServerManager) {
    framework.RegisterRoute("GET", "/health", server.HandlerFunc(func(ctx server.Context) error {
        config := manager.GetConfig()
        
        health := map[string]interface{}{
            "status":      "healthy",
            "timestamp":   time.Now().Unix(),
            "server": map[string]interface{}{
                "framework": config.Framework,
                "host":      config.Host,
                "port":      config.Port,
                "mode":      config.Mode,
                "running":   manager.IsRunning(),
            },
            "system": map[string]interface{}{
                "goroutines": runtime.NumGoroutine(),
                "memory":     getMemoryStats(),
                "uptime":     time.Since(startTime).Seconds(),
            },
        }
        
        return ctx.JSON(200, health)
    }))
}

func getMemoryStats() map[string]interface{} {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    return map[string]interface{}{
        "allocated":      m.Alloc,
        "total_alloc":    m.TotalAlloc,
        "sys":           m.Sys,
        "num_gc":        m.NumGC,
        "gc_cpu_fraction": m.GCCPUFraction,
    }
}
```

### Readiness Checks

```go
// Readiness check for Kubernetes
func setupReadinessCheck(framework server.WebFramework, manager *server.ServerManager) {
    framework.RegisterRoute("GET", "/ready", server.HandlerFunc(func(ctx server.Context) error {
        if !manager.IsRunning() {
            return ctx.JSON(503, map[string]string{
                "status": "not ready",
                "reason": "server not running",
            })
        }
        
        // Check dependencies (database, cache, etc.)
        if err := checkDependencies(); err != nil {
            return ctx.JSON(503, map[string]interface{}{
                "status": "not ready",
                "reason": "dependency check failed",
                "error":  err.Error(),
            })
        }
        
        return ctx.JSON(200, map[string]string{
            "status": "ready",
        })
    }))
}

func checkDependencies() error {
    // Check database connection
    if err := checkDatabase(); err != nil {
        return fmt.Errorf("database check failed: %w", err)
    }
    
    // Check cache connection
    if err := checkCache(); err != nil {
        return fmt.Errorf("cache check failed: %w", err)
    }
    
    // Check external services
    if err := checkExternalServices(); err != nil {
        return fmt.Errorf("external service check failed: %w", err)
    }
    
    return nil
}
```

### Metrics Collection

```go
// Setup metrics endpoints
func setupMetrics(framework server.WebFramework) {
    // Prometheus metrics
    framework.RegisterRoute("GET", "/metrics", server.HandlerFunc(func(ctx server.Context) error {
        // Return Prometheus metrics
        metrics := collectMetrics()
        ctx.SetHeader("Content-Type", "text/plain")
        return ctx.String(200, metrics)
    }))
    
    // Custom metrics
    framework.RegisterRoute("GET", "/stats", server.HandlerFunc(func(ctx server.Context) error {
        stats := map[string]interface{}{
            "requests_total":    requestCounter.Value(),
            "requests_per_sec":  getRequestsPerSecond(),
            "response_time_avg": getAverageResponseTime(),
            "active_connections": getActiveConnections(),
            "memory_usage":      getMemoryUsage(),
        }
        
        return ctx.JSON(200, stats)
    }))
}
```

## Advanced Server Management

### Multiple Server Instances

```go
// Manage multiple servers
type MultiServerManager struct {
    servers map[string]*server.ServerManager
    configs map[string]*server.ServerConfig
}

func NewMultiServerManager() *MultiServerManager {
    return &MultiServerManager{
        servers: make(map[string]*server.ServerManager),
        configs: make(map[string]*server.ServerConfig),
    }
}

func (msm *MultiServerManager) AddServer(name, framework string, config *server.ServerConfig) error {
    manager, err := server.CreateServerManager(framework, config)
    if err != nil {
        return fmt.Errorf("failed to create server %s: %w", name, err)
    }
    
    msm.servers[name] = manager
    msm.configs[name] = config
    
    return nil
}

func (msm *MultiServerManager) StartAll(ctx context.Context) error {
    var wg sync.WaitGroup
    errChan := make(chan error, len(msm.servers))
    
    for name, manager := range msm.servers {
        wg.Add(1)
        go func(name string, mgr *server.ServerManager) {
            defer wg.Done()
            
            if err := mgr.Start(ctx); err != nil {
                errChan <- fmt.Errorf("server %s failed: %w", name, err)
            }
        }(name, manager)
    }
    
    // Wait for all servers to start
    go func() {
        wg.Wait()
        close(errChan)
    }()
    
    // Check for errors
    for err := range errChan {
        return err
    }
    
    return nil
}

func (msm *MultiServerManager) StopAll(ctx context.Context) error {
    var wg sync.WaitGroup
    errChan := make(chan error, len(msm.servers))
    
    for name, manager := range msm.servers {
        wg.Add(1)
        go func(name string, mgr *server.ServerManager) {
            defer wg.Done()
            
            if err := mgr.Stop(ctx); err != nil {
                errChan <- fmt.Errorf("failed to stop server %s: %w", name, err)
            }
        }(name, manager)
    }
    
    // Wait for all servers to stop
    go func() {
        wg.Wait()
        close(errChan)
    }()
    
    // Collect errors
    var errors []error
    for err := range errChan {
        errors = append(errors, err)
    }
    
    if len(errors) > 0 {
        return fmt.Errorf("shutdown errors: %v", errors)
    }
    
    return nil
}
```

### Load Balancing

```go
// Simple load balancer
type LoadBalancer struct {
    servers []*server.ServerManager
    current int64
    mutex   sync.RWMutex
}

func NewLoadBalancer() *LoadBalancer {
    return &LoadBalancer{
        servers: make([]*server.ServerManager, 0),
    }
}

func (lb *LoadBalancer) AddServer(manager *server.ServerManager) {
    lb.mutex.Lock()
    defer lb.mutex.Unlock()
    
    lb.servers = append(lb.servers, manager)
}

func (lb *LoadBalancer) GetServer() *server.ServerManager {
    lb.mutex.RLock()
    defer lb.mutex.RUnlock()
    
    if len(lb.servers) == 0 {
        return nil
    }
    
    // Round-robin selection
    index := atomic.AddInt64(&lb.current, 1) % int64(len(lb.servers))
    return lb.servers[index]
}

func (lb *LoadBalancer) GetHealthyServers() []*server.ServerManager {
    lb.mutex.RLock()
    defer lb.mutex.RUnlock()
    
    healthy := make([]*server.ServerManager, 0)
    for _, server := range lb.servers {
        if server.IsRunning() {
            healthy = append(healthy, server)
        }
    }
    
    return healthy
}
```

### Configuration Hot Reload

```go
type ConfigWatcher struct {
    manager    *server.ServerManager
    configPath string
    lastMod    time.Time
    stopChan   chan bool
}

func NewConfigWatcher(manager *server.ServerManager, configPath string) *ConfigWatcher {
    return &ConfigWatcher{
        manager:    manager,
        configPath: configPath,
        stopChan:   make(chan bool),
    }
}

func (cw *ConfigWatcher) Start() {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            if cw.checkConfigChange() {
                if err := cw.reloadConfig(); err != nil {
                    log.Printf("Failed to reload config: %v", err)
                }
            }
        case <-cw.stopChan:
            return
        }
    }
}

func (cw *ConfigWatcher) Stop() {
    close(cw.stopChan)
}

func (cw *ConfigWatcher) checkConfigChange() bool {
    info, err := os.Stat(cw.configPath)
    if err != nil {
        return false
    }
    
    if info.ModTime().After(cw.lastMod) {
        cw.lastMod = info.ModTime()
        return true
    }
    
    return false
}

func (cw *ConfigWatcher) reloadConfig() error {
    // Load new configuration
    newConfig, err := loadConfigFromFile(cw.configPath)
    if err != nil {
        return fmt.Errorf("failed to load config: %w", err)
    }
    
    // Validate new configuration
    if err := newConfig.Validate(); err != nil {
        return fmt.Errorf("invalid config: %w", err)
    }
    
    // Apply configuration changes
    return cw.applyConfigChanges(newConfig)
}

func (cw *ConfigWatcher) applyConfigChanges(newConfig *server.ServerConfig) error {
    currentConfig := cw.manager.GetConfig()
    
    // Check if restart is required
    if cw.requiresRestart(currentConfig, newConfig) {
        log.Println("Configuration change requires server restart")
        
        // Graceful restart
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        
        if err := cw.manager.Stop(ctx); err != nil {
            return fmt.Errorf("failed to stop server: %w", err)
        }
        
        // Update configuration
        *currentConfig = *newConfig
        
        if err := cw.manager.Start(ctx); err != nil {
            return fmt.Errorf("failed to start server: %w", err)
        }
        
        log.Println("Server restarted with new configuration")
    } else {
        // Apply non-breaking changes
        log.Println("Applying configuration changes without restart")
        // Update only safe-to-change configuration
    }
    
    return nil
}

func (cw *ConfigWatcher) requiresRestart(old, new *server.ServerConfig) bool {
    // Changes that require restart
    return old.Host != new.Host ||
           old.Port != new.Port ||
           old.Framework != new.Framework ||
           old.Mode != new.Mode
}
```

## Production Deployment

### Docker Integration

```dockerfile
# Dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/server .
COPY --from=builder /app/configs ./configs

EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./server"]
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: lmcc-server
spec:
  replicas: 3
  selector:
    matchLabels:
      app: lmcc-server
  template:
    metadata:
      labels:
        app: lmcc-server
    spec:
      containers:
      - name: server
        image: lmcc-server:latest
        ports:
        - containerPort: 8080
        env:
        - name: SERVER_HOST
          value: "0.0.0.0"
        - name: SERVER_PORT
          value: "8080"
        - name: SERVER_MODE
          value: "production"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "64Mi"
            cpu: "250m"
          limits:
            memory: "128Mi"
            cpu: "500m"
---
apiVersion: v1
kind: Service
metadata:
  name: lmcc-server-service
spec:
  selector:
    app: lmcc-server
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: LoadBalancer
```

### Process Management with systemd

```ini
# /etc/systemd/system/lmcc-server.service
[Unit]
Description=LMCC Go Server
After=network.target

[Service]
Type=simple
User=lmcc
Group=lmcc
WorkingDirectory=/opt/lmcc-server
ExecStart=/opt/lmcc-server/server
ExecReload=/bin/kill -HUP $MAINPID
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=lmcc-server

# Security settings
NoNewPrivileges=yes
PrivateTmp=yes
ProtectSystem=strict
ProtectHome=yes
ReadWritePaths=/opt/lmcc-server/logs

[Install]
WantedBy=multi-user.target
```

### Environment Configuration

```go
// Production server setup
func setupProductionServer() (*server.ServerManager, error) {
    config := &server.ServerConfig{
        Framework: getEnv("SERVER_FRAMEWORK", "gin"),
        Host:     getEnv("SERVER_HOST", "0.0.0.0"),
        Port:     getEnvInt("SERVER_PORT", 8080),
        Mode:     getEnv("SERVER_MODE", "production"),
        
        ReadTimeout:    getEnvDuration("SERVER_READ_TIMEOUT", 30*time.Second),
        WriteTimeout:   getEnvDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
        IdleTimeout:    getEnvDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
        MaxHeaderBytes: getEnvInt("SERVER_MAX_HEADER_BYTES", 1<<20), // 1MB
        
        GracefulShutdown: server.GracefulShutdownConfig{
            Enabled:  getEnvBool("GRACEFUL_SHUTDOWN_ENABLED", true),
            Timeout:  getEnvDuration("GRACEFUL_SHUTDOWN_TIMEOUT", 30*time.Second),
            WaitTime: getEnvDuration("GRACEFUL_SHUTDOWN_WAIT", 5*time.Second),
        },
        
        TLS: server.TLSConfig{
            Enabled:  getEnvBool("TLS_ENABLED", false),
            CertFile: getEnv("TLS_CERT_FILE", ""),
            KeyFile:  getEnv("TLS_KEY_FILE", ""),
        },
    }
    
    return server.CreateServerManager(config.Framework, config)
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
    if value := os.Getenv(key); value != "" {
        if boolValue, err := strconv.ParseBool(value); err == nil {
            return boolValue
        }
    }
    return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
    if value := os.Getenv(key); value != "" {
        if duration, err := time.ParseDuration(value); err == nil {
            return duration
        }
    }
    return defaultValue
}
```

## Testing Server Management

### Unit Testing

```go
func TestServerManager_Lifecycle(t *testing.T) {
    // Create test configuration
    config := server.DefaultServerConfig()
    config.Framework = "gin"
    config.Port = 0 // Use random port
    
    // Create mock framework
    framework := &MockWebFramework{}
    manager := server.NewServerManager(framework, config)
    
    // Test initial state
    assert.False(t, manager.IsRunning())
    
    // Test startup
    ctx := context.Background()
    err := manager.Start(ctx)
    assert.NoError(t, err)
    assert.True(t, manager.IsRunning())
    
    // Test shutdown
    err = manager.Stop(ctx)
    assert.NoError(t, err)
    assert.False(t, manager.IsRunning())
}

func TestServerManager_GracefulShutdown(t *testing.T) {
    config := server.DefaultServerConfig()
    config.GracefulShutdown.Enabled = true
    config.GracefulShutdown.Timeout = 5 * time.Second
    
    framework := &MockWebFramework{}
    manager := server.NewServerManager(framework, config)
    
    // Start server
    go func() {
        ctx := context.Background()
        manager.Start(ctx)
    }()
    
    // Wait for server to start
    time.Sleep(100 * time.Millisecond)
    assert.True(t, manager.IsRunning())
    
    // Send shutdown signal
    p, _ := os.FindProcess(os.Getpid())
    p.Signal(syscall.SIGTERM)
    
    // Wait for graceful shutdown
    time.Sleep(1 * time.Second)
    assert.False(t, manager.IsRunning())
}
```

### Integration Testing

```go
func TestServerIntegration(t *testing.T) {
    // Create server
    config := server.DefaultServerConfig()
    config.Framework = "gin"
    config.Port = 0 // Random port
    
    manager, err := server.CreateServerManager("gin", config)
    assert.NoError(t, err)
    
    // Setup routes
    framework := manager.GetFramework()
    framework.RegisterRoute("GET", "/test", server.HandlerFunc(func(ctx server.Context) error {
        return ctx.JSON(200, map[string]string{"status": "ok"})
    }))
    
    // Start server in goroutine
    go func() {
        ctx := context.Background()
        manager.Start(ctx)
    }()
    
    // Wait for server to be ready
    time.Sleep(100 * time.Millisecond)
    assert.True(t, manager.IsRunning())
    
    // Make HTTP request
    port := getServerPort(manager)
    resp, err := http.Get(fmt.Sprintf("http://localhost:%d/test", port))
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
    
    // Cleanup
    ctx := context.Background()
    manager.Stop(ctx)
}
```

## Best Practices

### Resource Management

1. **Connection Limits**: Set appropriate connection limits
2. **Memory Management**: Monitor and limit memory usage
3. **File Descriptors**: Manage file descriptor limits
4. **Goroutine Management**: Prevent goroutine leaks

### Security

1. **TLS Configuration**: Use proper TLS settings for production
2. **Security Headers**: Implement security middleware
3. **Input Validation**: Validate all inputs
4. **Rate Limiting**: Implement rate limiting
5. **Access Control**: Secure administrative endpoints

### Monitoring

1. **Health Checks**: Implement comprehensive health checks
2. **Metrics Collection**: Collect and monitor metrics
3. **Logging**: Implement structured logging
4. **Alerting**: Set up alerting for critical issues

### Performance

1. **Keep-Alive**: Enable HTTP keep-alive
2. **Compression**: Enable response compression
3. **Caching**: Implement appropriate caching
4. **Connection Pooling**: Use connection pools for external services

## Troubleshooting

### Common Issues

1. **Port Already in Use**
   ```go
   // Check for port conflicts
   if err := checkPortAvailable(config.Port); err != nil {
       return fmt.Errorf("port %d is not available: %w", config.Port, err)
   }
   ```

2. **Graceful Shutdown Timeout**
   ```go
   // Increase shutdown timeout
   config.GracefulShutdown.Timeout = 60 * time.Second
   ```

3. **Memory Leaks**
   ```go
   // Monitor goroutines
   go func() {
       ticker := time.NewTicker(10 * time.Second)
       for range ticker.C {
           fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
       }
   }()
   ```

### Debug Information

```go
// Server debug information
func getServerDebugInfo(manager *server.ServerManager) map[string]interface{} {
    config := manager.GetConfig()
    
    return map[string]interface{}{
        "server": map[string]interface{}{
            "running":   manager.IsRunning(),
            "framework": config.Framework,
            "host":      config.Host,
            "port":      config.Port,
            "mode":      config.Mode,
        },
        "runtime": map[string]interface{}{
            "goroutines": runtime.NumGoroutine(),
            "memory":     getMemoryStats(),
            "gc_stats":   getGCStats(),
        },
        "config": config,
    }
}
```

## Next Steps

- **[Integration Examples](06_integration_examples.md)** - See real-world integration patterns
- **[Best Practices](07_best_practices.md)** - Production deployment guidelines
- **[Troubleshooting](08_troubleshooting.md)** - Common issues and solutions
- **[API Reference](09_api_reference.md)** - Complete API documentation
