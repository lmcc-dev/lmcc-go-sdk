# 服务器管理指南

lmcc-go-sdk 服务器模块通过 `ServerManager` 组件提供全面的服务器生命周期管理。本指南涵盖服务器创建、启动、关闭、监控以及生产部署的最佳实践。

## 服务器管理概述

### 核心组件

服务器管理系统包括：

1. **ServerManager**: 管理服务器生命周期的核心组件
2. **ServerFactory**: 创建服务器实例的工厂
3. **优雅关闭**: 带信号管理的清洁关闭处理
4. **健康监控**: 服务器状态监控和报告
5. **配置管理**: 运行时配置和验证

### 服务器生命周期

```
创建 → 配置 → 启动 → 运行 → 关闭 → 清理
  ↓      ↓      ↓      ↓      ↓      ↓
工厂   验证   监听   服务   优雅   完成
```

## 服务器创建

### 使用 ServerManager

```go
package main

import (
    "context"
    "log"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin"
)

func main() {
    // 创建配置 (Create configuration)
    config := server.DefaultServerConfig()
    config.Framework = "gin"
    config.Host = "0.0.0.0"
    config.Port = 8080
    
    // 创建服务器管理器 (Create server manager)
    manager, err := server.CreateServerManager("gin", config)
    if err != nil {
        log.Fatal("Failed to create server:", err)
    }
    
    // 配置路由和中间件 (Configure routes and middleware)
    framework := manager.GetFramework()
    framework.RegisterRoute("GET", "/health", server.HandlerFunc(healthHandler))
    
    // 启动服务器 (Start server)
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

### 使用 ServerFactory

```go
// 创建工厂 (Create factory)
factory := server.NewServerFactory()

// 使用工厂创建服务器 (Create server using factory)
manager, err := factory.CreateServer("gin", config)
if err != nil {
    log.Fatal("Failed to create server:", err)
}

// 向工厂注册自定义插件 (Register custom plugins with factory)
customPlugin := myplugin.NewPlugin()
err = factory.RegisterPlugin(customPlugin)
if err != nil {
    log.Fatal("Failed to register plugin:", err)
}
```

### 配置验证

```go
// 创建配置 (Create configuration)
config := &server.ServerConfig{
    Framework: "gin",
    Host:     "0.0.0.0",
    Port:     8080,
    Mode:     "production",
}

// 在创建服务器之前验证配置 (Validate configuration before creating server)
if err := config.Validate(); err != nil {
    log.Fatal("Invalid configuration:", err)
}

// 使用验证的配置创建服务器 (Create server with validated configuration)
manager, err := server.CreateServerManager("gin", config)
if err != nil {
    log.Fatal("Failed to create server:", err)
}
```

## 服务器启动

### 基本启动

```go
func startServer(manager *server.ServerManager) error {
    ctx := context.Background()
    
    // 启动服务器（阻塞调用） (Start server - blocking call)
    if err := manager.Start(ctx); err != nil {
        return fmt.Errorf("server startup failed: %w", err)
    }
    
    return nil
}
```

### 非阻塞启动

```go
func startServerAsync(manager *server.ServerManager) error {
    ctx := context.Background()
    
    // 在 goroutine 中启动服务器 (Start server in goroutine)
    go func() {
        if err := manager.Start(ctx); err != nil {
            log.Printf("Server error: %v", err)
        }
    }()
    
    // 等待服务器准备就绪 (Wait for server to be ready)
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

### 带上下文的启动

```go
func startServerWithTimeout(manager *server.ServerManager) error {
    // 创建带超时的上下文 (Create context with timeout)
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // 带超时启动服务器 (Start server with timeout)
    if err := manager.Start(ctx); err != nil {
        return fmt.Errorf("server startup failed: %w", err)
    }
    
    return nil
}
```

### 启动期间错误处理

```go
func startServerWithRetry(manager *server.ServerManager, maxRetries int) error {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        ctx := context.Background()
        
        if err := manager.Start(ctx); err != nil {
            lastErr = err
            log.Printf("Startup attempt %d failed: %v", i+1, err)
            
            // 重试前等待 (Wait before retry)
            time.Sleep(time.Duration(i+1) * time.Second)
            continue
        }
        
        return nil // 成功 (Success)
    }
    
    return fmt.Errorf("failed to start server after %d attempts: %w", maxRetries, lastErr)
}
```

## 服务器关闭

### 优雅关闭配置

```go
config := &server.ServerConfig{
    Framework: "gin",
    Host:     "0.0.0.0",
    Port:     8080,
    
    GracefulShutdown: server.GracefulShutdownConfig{
        Enabled:  true,
        Timeout:  30 * time.Second,  // 等待关闭的最大时间 (Maximum time to wait for shutdown)
        WaitTime: 5 * time.Second,   // 等待正在进行请求的时间 (Time to wait for ongoing requests)
    },
}
```

### 手动关闭

```go
func shutdownServer(manager *server.ServerManager) error {
    if !manager.IsRunning() {
        return fmt.Errorf("server is not running")
    }
    
    // 创建带超时的上下文 (Create context with timeout)
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // 停止服务器 (Stop server)
    if err := manager.Stop(ctx); err != nil {
        return fmt.Errorf("failed to stop server: %w", err)
    }
    
    log.Println("Server stopped successfully")
    return nil
}
```

### 基于信号的关闭

```go
func startServerWithSignalHandling(manager *server.ServerManager) error {
    // 创建信号通道 (Create signal channel)
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    // 在 goroutine 中启动服务器 (Start server in goroutine)
    serverErrChan := make(chan error, 1)
    go func() {
        ctx := context.Background()
        serverErrChan <- manager.Start(ctx)
    }()
    
    // 等待信号或服务器错误 (Wait for signal or server error)
    select {
    case sig := <-sigChan:
        log.Printf("Received signal %v, shutting down server...", sig)
        
        // 创建关闭上下文 (Create shutdown context)
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        
        // 停止服务器 (Stop server)
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

### 自动优雅关闭

启用时，服务器管理器包含内置的优雅关闭：

```go
config := &server.ServerConfig{
    Framework: "gin",
    Host:     "0.0.0.0",
    Port:     8080,
    
    GracefulShutdown: server.GracefulShutdownConfig{
        Enabled:  true,                // 启用自动信号处理 (Enable automatic signal handling)
        Timeout:  30 * time.Second,    // 关闭超时 (Shutdown timeout)
        WaitTime: 5 * time.Second,     // 等待正在进行的请求 (Wait for ongoing requests)
    },
}

manager, err := server.CreateServerManager("gin", config)
if err != nil {
    log.Fatal("Failed to create server:", err)
}

// 服务器将自动处理 SIGINT 和 SIGTERM (The server will automatically handle SIGINT and SIGTERM)
// 无需额外的信号处理 (No additional signal handling needed)
if err := manager.Start(context.Background()); err != nil {
    log.Fatal("Server failed:", err)
}
```

## 服务器监控

### 健康检查

```go
// 内置健康检查 (Built-in health check)
func setupHealthCheck(framework server.WebFramework) {
    framework.RegisterRoute("GET", "/health", server.HandlerFunc(func(ctx server.Context) error {
        return ctx.JSON(200, map[string]interface{}{
            "status":    "healthy",
            "timestamp": time.Now().Unix(),
            "uptime":    time.Since(startTime).Seconds(),
        })
    }))
}

// 详细健康检查 (Detailed health check)
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

### 就绪检查

```go
// Kubernetes 就绪检查 (Readiness check for Kubernetes)
func setupReadinessCheck(framework server.WebFramework, manager *server.ServerManager) {
    framework.RegisterRoute("GET", "/ready", server.HandlerFunc(func(ctx server.Context) error {
        if !manager.IsRunning() {
            return ctx.JSON(503, map[string]string{
                "status": "not ready",
                "reason": "server not running",
            })
        }
        
        // 检查依赖（数据库、缓存等） (Check dependencies - database, cache, etc.)
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
    // 检查数据库连接 (Check database connection)
    if err := checkDatabase(); err != nil {
        return fmt.Errorf("database check failed: %w", err)
    }
    
    // 检查缓存连接 (Check cache connection)
    if err := checkCache(); err != nil {
        return fmt.Errorf("cache check failed: %w", err)
    }
    
    // 检查外部服务 (Check external services)
    if err := checkExternalServices(); err != nil {
        return fmt.Errorf("external service check failed: %w", err)
    }
    
    return nil
}
```

### 指标收集

```go
// 设置指标端点 (Setup metrics endpoints)
func setupMetrics(framework server.WebFramework) {
    // Prometheus 指标 (Prometheus metrics)
    framework.RegisterRoute("GET", "/metrics", server.HandlerFunc(func(ctx server.Context) error {
        // 返回 Prometheus 指标 (Return Prometheus metrics)
        metrics := collectMetrics()
        ctx.SetHeader("Content-Type", "text/plain")
        return ctx.String(200, metrics)
    }))
    
    // 自定义指标 (Custom metrics)
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

## 高级服务器管理

### 多服务器实例

```go
// 管理多个服务器 (Manage multiple servers)
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
    
    // 等待所有服务器启动 (Wait for all servers to start)
    go func() {
        wg.Wait()
        close(errChan)
    }()
    
    // 检查错误 (Check for errors)
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
    
    // 等待所有服务器停止 (Wait for all servers to stop)
    go func() {
        wg.Wait()
        close(errChan)
    }()
    
    // 收集错误 (Collect errors)
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

### 负载均衡

```go
// 简单负载均衡器 (Simple load balancer)
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
    
    // 轮询选择 (Round-robin selection)
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

### 配置热重载

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
    // 加载新配置 (Load new configuration)
    newConfig, err := loadConfigFromFile(cw.configPath)
    if err != nil {
        return fmt.Errorf("failed to load config: %w", err)
    }
    
    // 验证新配置 (Validate new configuration)
    if err := newConfig.Validate(); err != nil {
        return fmt.Errorf("invalid config: %w", err)
    }
    
    // 应用配置更改 (Apply configuration changes)
    return cw.applyConfigChanges(newConfig)
}

func (cw *ConfigWatcher) applyConfigChanges(newConfig *server.ServerConfig) error {
    currentConfig := cw.manager.GetConfig()
    
    // 检查是否需要重启 (Check if restart is required)
    if cw.requiresRestart(currentConfig, newConfig) {
        log.Println("Configuration change requires server restart")
        
        // 优雅重启 (Graceful restart)
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        
        if err := cw.manager.Stop(ctx); err != nil {
            return fmt.Errorf("failed to stop server: %w", err)
        }
        
        // 更新配置 (Update configuration)
        *currentConfig = *newConfig
        
        if err := cw.manager.Start(ctx); err != nil {
            return fmt.Errorf("failed to start server: %w", err)
        }
        
        log.Println("Server restarted with new configuration")
    } else {
        // 应用非破坏性更改 (Apply non-breaking changes)
        log.Println("Applying configuration changes without restart")
        // 只更新安全更改的配置 (Update only safe-to-change configuration)
    }
    
    return nil
}

func (cw *ConfigWatcher) requiresRestart(old, new *server.ServerConfig) bool {
    // 需要重启的更改 (Changes that require restart)
    return old.Host != new.Host ||
           old.Port != new.Port ||
           old.Framework != new.Framework ||
           old.Mode != new.Mode
}
```

## 生产部署

### Docker 集成

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

# 健康检查 (Health check)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./server"]
```

### Kubernetes 部署

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

### systemd 进程管理

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

# 安全设置 (Security settings)
NoNewPrivileges=yes
PrivateTmp=yes
ProtectSystem=strict
ProtectHome=yes
ReadWritePaths=/opt/lmcc-server/logs

[Install]
WantedBy=multi-user.target
```

### 环境配置

```go
// 生产服务器设置 (Production server setup)
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

## 最佳实践

### 资源管理

1. **连接限制**: 设置适当的连接限制
2. **内存管理**: 监控和限制内存使用
3. **文件描述符**: 管理文件描述符限制
4. **Goroutine 管理**: 防止 goroutine 泄漏

### 安全

1. **TLS 配置**: 为生产使用适当的 TLS 设置
2. **安全头**: 实现安全中间件
3. **输入验证**: 验证所有输入
4. **限流**: 实现限流
5. **访问控制**: 保护管理端点

### 监控

1. **健康检查**: 实现全面的健康检查
2. **指标收集**: 收集和监控指标
3. **日志记录**: 实现结构化日志
4. **告警**: 为关键问题设置告警

### 性能

1. **Keep-Alive**: 启用 HTTP keep-alive
2. **压缩**: 启用响应压缩
3. **缓存**: 实现适当的缓存
4. **连接池**: 为外部服务使用连接池

## 故障排除

### 常见问题

1. **端口已在使用**
   ```go
   // 检查端口冲突 (Check for port conflicts)
   if err := checkPortAvailable(config.Port); err != nil {
       return fmt.Errorf("port %d is not available: %w", config.Port, err)
   }
   ```

2. **优雅关闭超时**
   ```go
   // 增加关闭超时 (Increase shutdown timeout)
   config.GracefulShutdown.Timeout = 60 * time.Second
   ```

3. **内存泄漏**
   ```go
   // 监控 goroutines (Monitor goroutines)
   go func() {
       ticker := time.NewTicker(10 * time.Second)
       for range ticker.C {
           fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
       }
   }()
   ```

### 调试信息

```go
// 服务器调试信息 (Server debug information)
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

## 下一步

- **[集成示例](06_integration_examples.md)** - 查看真实世界集成模式
- **[最佳实践](07_best_practices.md)** - 生产部署指南
- **[故障排除](08_troubleshooting.md)** - 常见问题和解决方案
- **[API 参考](09_api_reference.md)** - 完整 API 文档