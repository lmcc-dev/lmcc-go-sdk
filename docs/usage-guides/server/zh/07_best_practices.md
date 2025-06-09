# 最佳实践指南

本指南涵盖了有效且安全地使用 lmcc-go-sdk 服务器模块的生产就绪实践。

## 架构最佳实践

### 分层架构

```go
// 推荐的项目结构 (Recommended project structure)
project/
├── cmd/                    // 应用程序入口点 (Application entry points)
│   └── server/
│       └── main.go
├── internal/              // 私有应用程序代码 (Private application code)
│   ├── api/              // API 处理器 (API handlers)
│   ├── service/          // 业务逻辑 (Business logic)
│   ├── repository/       // 数据访问 (Data access)
│   └── config/          // 配置 (Configuration)
├── pkg/                  // 公共库 (Public libraries)
└── configs/             // 配置文件 (Configuration files)
```

### 依赖注入

```go
type Container struct {
    Config     *config.Config
    Logger     *log.Logger
    DB         *gorm.DB
    UserRepo   repository.UserRepository
    UserSvc    service.UserService
}

func NewContainer() (*Container, error) {
    cfg := config.Load()
    
    logger := log.New()
    
    db, err := gorm.Open(postgres.Open(cfg.DatabaseURL))
    if err != nil {
        return nil, err
    }
    
    userRepo := repository.NewUserRepository(db)
    userSvc := service.NewUserService(userRepo)
    
    return &Container{
        Config:   cfg,
        Logger:   logger,
        DB:       db,
        UserRepo: userRepo,
        UserSvc:  userSvc,
    }, nil
}

func (c *Container) SetupRoutes(framework server.WebFramework) {
    api := framework.Group("/api/v1")
    
    userHandler := api.NewUserHandler(c.UserSvc)
    api.RegisterRoute("GET", "/users", userHandler.GetUsers)
    api.RegisterRoute("POST", "/users", userHandler.CreateUser)
}
```

## 配置管理

### 基于环境的配置

```go
type Config struct {
    Server   ServerConfig   `yaml:"server"`
    Database DatabaseConfig `yaml:"database"`
    Redis    RedisConfig    `yaml:"redis"`
    Log      LogConfig      `yaml:"log"`
}

type ServerConfig struct {
    Framework string        `yaml:"framework" env:"SERVER_FRAMEWORK" default:"gin"`
    Host      string        `yaml:"host" env:"SERVER_HOST" default:"0.0.0.0"`
    Port      int           `yaml:"port" env:"SERVER_PORT" default:"8080"`
    Mode      string        `yaml:"mode" env:"SERVER_MODE" default:"release"`
    Timeout   time.Duration `yaml:"timeout" env:"SERVER_TIMEOUT" default:"30s"`
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    
    // 从 YAML 文件加载 (Load from YAML file)
    if err := loadFromYAML(cfg, "config.yaml"); err != nil {
        return nil, err
    }
    
    // 用环境变量覆盖 (Override with environment variables)
    if err := loadFromEnv(cfg); err != nil {
        return nil, err
    }
    
    // 验证配置 (Validate configuration)
    if err := cfg.Validate(); err != nil {
        return nil, err
    }
    
    return cfg, nil
}

func (c *Config) Validate() error {
    if c.Server.Port < 1 || c.Server.Port > 65535 {
        return fmt.Errorf("invalid port: %d", c.Server.Port)
    }
    
    if c.Server.Framework == "" {
        return fmt.Errorf("framework is required")
    }
    
    return nil
}
```

### 多环境支持

```yaml
# config/development.yaml
server:
  framework: gin
  host: localhost
  port: 8080
  mode: debug
  
database:
  url: postgres://user:pass@localhost/dev_db

log:
  level: debug
  format: text

---
# config/production.yaml
server:
  framework: gin
  host: 0.0.0.0
  port: 8080
  mode: release
  
database:
  url: ${DATABASE_URL}

log:
  level: info
  format: json
```

## 安全最佳实践

### 输入验证

```go
type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

func (h *UserHandler) CreateUser(ctx server.Context) error {
    var req CreateUserRequest
    if err := ctx.Bind(&req); err != nil {
        return ctx.JSON(400, ErrorResponse{
            Error: "Invalid request format",
            Code:  "INVALID_FORMAT",
        })
    }
    
    // 验证输入 (Validate input)
    if err := h.validator.Struct(req); err != nil {
        return ctx.JSON(400, ErrorResponse{
            Error: "Validation failed",
            Code:  "VALIDATION_ERROR",
            Details: formatValidationErrors(err),
        })
    }
    
    // 清理输入 (Sanitize input)
    req.Name = html.EscapeString(strings.TrimSpace(req.Name))
    req.Email = strings.ToLower(strings.TrimSpace(req.Email))
    
    // 处理请求 (Process request)
    user, err := h.userService.CreateUser(req)
    if err != nil {
        return ctx.JSON(500, ErrorResponse{
            Error: "Failed to create user",
            Code:  "CREATION_FAILED",
        })
    }
    
    return ctx.JSON(201, user)
}
```

### 限流

```go
func RateLimitMiddleware(limit int, window time.Duration) server.Middleware {
    limiter := rate.NewLimiter(rate.Every(window), limit)
    
    return server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
        if !limiter.Allow() {
            return ctx.JSON(429, map[string]string{
                "error": "Rate limit exceeded",
                "retry_after": window.String(),
            })
        }
        
        return next()
    })
}

// 使用 (Usage)
framework.RegisterMiddleware(RateLimitMiddleware(100, time.Minute))
```

### CORS 配置

```go
func CORSMiddleware(allowedOrigins []string) server.Middleware {
    return server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
        origin := ctx.Header("Origin")
        
        // 检查来源是否被允许 (Check if origin is allowed)
        allowed := false
        for _, allowedOrigin := range allowedOrigins {
            if origin == allowedOrigin || allowedOrigin == "*" {
                allowed = true
                break
            }
        }
        
        if allowed {
            ctx.SetHeader("Access-Control-Allow-Origin", origin)
            ctx.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            ctx.SetHeader("Access-Control-Allow-Headers", "Content-Type, Authorization")
            ctx.SetHeader("Access-Control-Allow-Credentials", "true")
        }
        
        if ctx.Method() == "OPTIONS" {
            return ctx.Status(204)
        }
        
        return next()
    })
}
```

### 安全头

```go
func SecurityHeadersMiddleware() server.Middleware {
    return server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
        // 防止 XSS 攻击 (Prevent XSS attacks)
        ctx.SetHeader("X-XSS-Protection", "1; mode=block")
        
        // 防止内容类型嗅探 (Prevent content type sniffing)
        ctx.SetHeader("X-Content-Type-Options", "nosniff")
        
        // 防止点击劫持 (Prevent clickjacking)
        ctx.SetHeader("X-Frame-Options", "DENY")
        
        // 强制 HTTPS (Force HTTPS)
        ctx.SetHeader("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        
        // 内容安全策略 (Content Security Policy)
        ctx.SetHeader("Content-Security-Policy", "default-src 'self'")
        
        return next()
    })
}
```

## 性能优化

### 连接池

```go
func setupDatabase(cfg *DatabaseConfig) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(cfg.URL), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return nil, err
    }
    
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }
    
    // 连接池设置 (Connection pool settings)
    sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)    // 最大打开连接数 (Maximum open connections)
    sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)    // 最大空闲连接数 (Maximum idle connections)
    sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime) // 连接生命周期 (Connection lifetime)
    sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime) // 连接空闲时间 (Connection idle time)
    
    return db, nil
}
```

### 缓存策略

```go
type CachedUserService struct {
    userRepo repository.UserRepository
    cache    cache.Cache
    ttl      time.Duration
}

func (s *CachedUserService) GetUser(id int) (*User, error) {
    // 先尝试缓存 (Try cache first)
    key := fmt.Sprintf("user:%d", id)
    if cached := s.cache.Get(key); cached != nil {
        if user, ok := cached.(*User); ok {
            return user, nil
        }
    }
    
    // 从数据库获取 (Fetch from database)
    user, err := s.userRepo.GetByID(id)
    if err != nil {
        return nil, err
    }
    
    // 缓存结果 (Cache the result)
    s.cache.Set(key, user, s.ttl)
    
    return user, nil
}

func (s *CachedUserService) UpdateUser(id int, updates *User) (*User, error) {
    user, err := s.userRepo.Update(id, updates)
    if err != nil {
        return nil, err
    }
    
    // 使缓存失效 (Invalidate cache)
    key := fmt.Sprintf("user:%d", id)
    s.cache.Delete(key)
    
    return user, nil
}
```

### 响应压缩

```go
func CompressionMiddleware() server.Middleware {
    return server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
        // 检查客户端是否接受压缩 (Check if client accepts compression)
        acceptEncoding := ctx.Header("Accept-Encoding")
        if !strings.Contains(acceptEncoding, "gzip") {
            return next()
        }
        
        // 设置压缩头 (Set compression headers)
        ctx.SetHeader("Content-Encoding", "gzip")
        ctx.SetHeader("Vary", "Accept-Encoding")
        
        // 用 gzip writer 包装响应写入器 (Wrap response writer with gzip writer)
        gz := gzip.NewWriter(ctx.Response())
        defer gz.Close()
        
        // 替换响应写入器 (Replace response writer)
        originalWriter := ctx.Response()
        ctx.SetResponse(&gzipResponseWriter{
            ResponseWriter: originalWriter,
            Writer:        gz,
        })
        
        return next()
    })
}
```

## 错误处理

### 结构化错误响应

```go
type ErrorResponse struct {
    Error   string            `json:"error"`
    Code    string            `json:"code"`
    Details map[string]string `json:"details,omitempty"`
    TraceID string            `json:"trace_id,omitempty"`
}

type AppError struct {
    Message string
    Code    string
    Status  int
    Cause   error
}

func (e *AppError) Error() string {
    return e.Message
}

func ErrorHandlerMiddleware() server.Middleware {
    return server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
        err := next()
        if err == nil {
            return nil
        }
        
        // 生成跟踪 ID (Generate trace ID)
        traceID := generateTraceID()
        
        // 记录错误 (Log error)
        log.WithFields(log.Fields{
            "trace_id": traceID,
            "path":     ctx.Path(),
            "method":   ctx.Method(),
            "error":    err.Error(),
        }).Error("Request failed")
        
        // 处理不同错误类型 (Handle different error types)
        switch e := err.(type) {
        case *AppError:
            return ctx.JSON(e.Status, ErrorResponse{
                Error:   e.Message,
                Code:    e.Code,
                TraceID: traceID,
            })
        case *ValidationError:
            return ctx.JSON(400, ErrorResponse{
                Error:   "Validation failed",
                Code:    "VALIDATION_ERROR",
                Details: e.Details,
                TraceID: traceID,
            })
        default:
            return ctx.JSON(500, ErrorResponse{
                Error:   "Internal server error",
                Code:    "INTERNAL_ERROR",
                TraceID: traceID,
            })
        }
    })
}
```

### Panic 恢复

```go
func RecoveryMiddleware() server.Middleware {
    return server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
        defer func() {
            if r := recover(); r != nil {
                // 记录带堆栈跟踪的 panic (Log panic with stack trace)
                log.WithFields(log.Fields{
                    "panic": r,
                    "stack": debug.Stack(),
                    "path":  ctx.Path(),
                    "method": ctx.Method(),
                }).Error("Panic recovered")
                
                // 返回错误响应 (Return error response)
                ctx.JSON(500, ErrorResponse{
                    Error: "Internal server error",
                    Code:  "PANIC_RECOVERED",
                })
            }
        }()
        
        return next()
    })
}
```

## 日志最佳实践

### 结构化日志

```go
func LoggingMiddleware(logger *log.Logger) server.Middleware {
    return server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
        start := time.Now()
        
        // 向上下文添加请求 ID (Add request ID to context)
        requestID := generateRequestID()
        ctx.Set("request_id", requestID)
        
        // 创建请求日志器 (Create request logger)
        reqLogger := logger.WithFields(log.Fields{
            "request_id": requestID,
            "method":     ctx.Method(),
            "path":       ctx.Path(),
            "user_agent": ctx.Header("User-Agent"),
            "remote_ip":  ctx.ClientIP(),
        })
        
        // 记录请求开始 (Log request start)
        reqLogger.Info("Request started")
        
        // 执行下一个中间件 (Execute next middleware)
        err := next()
        
        // 计算持续时间 (Calculate duration)
        duration := time.Since(start)
        
        // 记录请求完成 (Log request completion)
        logFields := log.Fields{
            "duration": duration.Milliseconds(),
            "status":   ctx.Status(),
        }
        
        if err != nil {
            logFields["error"] = err.Error()
            reqLogger.WithFields(logFields).Error("Request failed")
        } else {
            reqLogger.WithFields(logFields).Info("Request completed")
        }
        
        return err
    })
}
```

### 日志级别

```go
// 开发：调试级别以获取详细信息 (Development: Debug level for detailed information)
if config.Environment == "development" {
    log.SetLevel(log.DebugLevel)
}

// 生产：信息级别用于正常操作 (Production: Info level for normal operations)
if config.Environment == "production" {
    log.SetLevel(log.InfoLevel)
    log.SetFormatter(&log.JSONFormatter{})
}

// 使用适当的日志级别 (Use appropriate log levels)
func (s *UserService) CreateUser(user *User) error {
    log.Debug("Creating user", "email", user.Email)
    
    if err := s.repo.Create(user); err != nil {
        log.Error("Failed to create user", "error", err, "email", user.Email)
        return err
    }
    
    log.Info("User created successfully", "id", user.ID, "email", user.Email)
    return nil
}
```

## 测试最佳实践

### 单元测试

```go
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name    string
        user    *User
        wantErr bool
        setup   func(*MockUserRepository)
    }{
        {
            name: "successful creation",
            user: &User{Name: "John", Email: "john@example.com"},
            wantErr: false,
            setup: func(repo *MockUserRepository) {
                repo.On("Create", mock.AnythingOfType("*User")).Return(nil)
            },
        },
        {
            name: "duplicate email",
            user: &User{Name: "John", Email: "existing@example.com"},
            wantErr: true,
            setup: func(repo *MockUserRepository) {
                repo.On("Create", mock.AnythingOfType("*User")).Return(errors.New("duplicate email"))
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := &MockUserRepository{}
            tt.setup(repo)
            
            service := NewUserService(repo)
            err := service.CreateUser(tt.user)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
            }
            
            repo.AssertExpectations(t)
        })
    }
}
```

### 集成测试

```go
func TestUserAPI_Integration(t *testing.T) {
    // 设置测试数据库 (Setup test database)
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    // 创建服务器 (Create server)
    config := server.DefaultServerConfig()
    config.Port = 0 // 随机端口 (Random port)
    
    manager, err := server.CreateServerManager("gin", config)
    require.NoError(t, err)
    
    // 设置依赖 (Setup dependencies)
    userRepo := repository.NewUserRepository(db)
    userService := service.NewUserService(userRepo)
    userHandler := api.NewUserHandler(userService)
    
    // 注册路由 (Register routes)
    framework := manager.GetFramework()
    api := framework.Group("/api/v1")
    api.RegisterRoute("POST", "/users", userHandler.CreateUser)
    api.RegisterRoute("GET", "/users/:id", userHandler.GetUser)
    
    // 启动服务器 (Start server)
    go manager.Start(context.Background())
    defer manager.Stop(context.Background())
    
    // 等待服务器启动 (Wait for server to start)
    time.Sleep(100 * time.Millisecond)
    
    // 测试 API 端点 (Test API endpoints)
    baseURL := fmt.Sprintf("http://localhost:%d", getServerPort(manager))
    
    t.Run("create user", func(t *testing.T) {
        user := map[string]string{
            "name":  "John Doe",
            "email": "john@example.com",
        }
        
        resp, err := postJSON(baseURL+"/api/v1/users", user)
        require.NoError(t, err)
        assert.Equal(t, 201, resp.StatusCode)
    })
    
    t.Run("get user", func(t *testing.T) {
        resp, err := http.Get(baseURL + "/api/v1/users/1")
        require.NoError(t, err)
        assert.Equal(t, 200, resp.StatusCode)
    })
}
```

## 监控和可观测性

### 健康检查

```go
type HealthChecker struct {
    db    *gorm.DB
    redis *redis.Client
}

func (hc *HealthChecker) CheckHealth() map[string]interface{} {
    health := map[string]interface{}{
        "status":    "healthy",
        "timestamp": time.Now().Unix(),
        "checks":    make(map[string]interface{}),
    }
    
    // 数据库健康 (Database health)
    if err := hc.checkDatabase(); err != nil {
        health["status"] = "unhealthy"
        health["checks"].(map[string]interface{})["database"] = map[string]interface{}{
            "status": "failed",
            "error":  err.Error(),
        }
    } else {
        health["checks"].(map[string]interface{})["database"] = map[string]interface{}{
            "status": "healthy",
        }
    }
    
    // Redis 健康 (Redis health)
    if err := hc.checkRedis(); err != nil {
        health["status"] = "unhealthy"
        health["checks"].(map[string]interface{})["redis"] = map[string]interface{}{
            "status": "failed",
            "error":  err.Error(),
        }
    } else {
        health["checks"].(map[string]interface{})["redis"] = map[string]interface{}{
            "status": "healthy",
        }
    }
    
    return health
}

func (hc *HealthChecker) checkDatabase() error {
    sqlDB, err := hc.db.DB()
    if err != nil {
        return err
    }
    
    return sqlDB.Ping()
}

func (hc *HealthChecker) checkRedis() error {
    return hc.redis.Ping(context.Background()).Err()
}
```

### 指标收集

```go
type Metrics struct {
    requestsTotal    prometheus.CounterVec
    requestDuration  prometheus.HistogramVec
    activeConnections prometheus.Gauge
}

func NewMetrics() *Metrics {
    return &Metrics{
        requestsTotal: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "http_requests_total",
                Help: "Total number of HTTP requests",
            },
            []string{"method", "path", "status"},
        ),
        requestDuration: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "http_request_duration_seconds",
                Help: "HTTP request duration in seconds",
            },
            []string{"method", "path"},
        ),
        activeConnections: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Name: "http_active_connections",
                Help: "Number of active HTTP connections",
            },
        ),
    }
}

func (m *Metrics) MetricsMiddleware() server.Middleware {
    return server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
        start := time.Now()
        m.activeConnections.Inc()
        defer m.activeConnections.Dec()
        
        err := next()
        
        duration := time.Since(start).Seconds()
        status := strconv.Itoa(ctx.Status())
        
        m.requestsTotal.WithLabelValues(ctx.Method(), ctx.Path(), status).Inc()
        m.requestDuration.WithLabelValues(ctx.Method(), ctx.Path()).Observe(duration)
        
        return err
    })
}
```

## 部署最佳实践

### Docker 最佳实践

```dockerfile
# 多阶段构建以获得更小的镜像 (Multi-stage build for smaller images)
FROM golang:1.24-alpine AS builder

# 安装依赖 (Install dependencies)
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# 首先复制 go mod 文件以获得更好的缓存 (Copy go mod files first for better caching)
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码 (Copy source code)
COPY . .

# 构建应用程序 (Build the application)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# 最终阶段 (Final stage)
FROM alpine:latest

# 为 HTTPS 安装 ca-certificates (Install ca-certificates for HTTPS)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从构建器阶段复制二进制文件 (Copy binary from builder stage)
COPY --from=builder /app/server .

# 复制配置文件 (Copy configuration files)
COPY --from=builder /app/configs ./configs

# 创建非 root 用户 (Create non-root user)
RUN adduser -D -s /bin/sh appuser
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./server"]
```

### Kubernetes 配置

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  config.yaml: |
    server:
      framework: gin
      host: 0.0.0.0
      port: 8080
      mode: release
    database:
      url: postgres://user:pass@postgres:5432/app

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
      - name: app
        image: myapp:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: database-url
        volumeMounts:
        - name: config-volume
          mountPath: /root/configs
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
      volumes:
      - name: config-volume
        configMap:
          name: app-config
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

## 安全清单

### 生产安全检查

1. **身份验证和授权** ✓
   - [ ] 实现强身份验证
   - [ ] 使用 RBAC（基于角色的访问控制）
   - [ ] 保护敏感端点

2. **数据保护** ✓
   - [ ] 加密传输中的数据（HTTPS）
   - [ ] 加密静态数据
   - [ ] 清理敏感日志

3. **输入验证** ✓
   - [ ] 验证所有用户输入
   - [ ] 清理 HTML/SQL 注入
   - [ ] 实现请求大小限制

4. **网络安全** ✓
   - [ ] 配置防火墙
   - [ ] 使用 VPN/私有网络
   - [ ] 实现 DDoS 保护

5. **监控和审计** ✓
   - [ ] 记录安全事件
   - [ ] 监控异常行为
   - [ ] 定期安全审计

## 下一步

- **[故障排除](08_troubleshooting.md)** - 常见问题和解决方案
- **[API 参考](09_api_reference.md)** - 完整 API 文档