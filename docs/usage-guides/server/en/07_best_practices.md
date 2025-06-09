# Best Practices Guide

This guide covers production-ready practices for using the lmcc-go-sdk server module effectively and securely.

## Architecture Best Practices

### Layered Architecture

```go
// Recommended project structure
project/
├── cmd/                    // Application entry points
│   └── server/
│       └── main.go
├── internal/              // Private application code
│   ├── api/              // API handlers
│   ├── service/          // Business logic
│   ├── repository/       // Data access
│   └── config/          // Configuration
├── pkg/                  // Public libraries
└── configs/             // Configuration files
```

### Dependency Injection

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

## Configuration Management

### Environment-Based Configuration

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
    
    // Load from YAML file
    if err := loadFromYAML(cfg, "config.yaml"); err != nil {
        return nil, err
    }
    
    // Override with environment variables
    if err := loadFromEnv(cfg); err != nil {
        return nil, err
    }
    
    // Validate configuration
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

### Multiple Environment Support

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

## Security Best Practices

### Input Validation

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
    
    // Validate input
    if err := h.validator.Struct(req); err != nil {
        return ctx.JSON(400, ErrorResponse{
            Error: "Validation failed",
            Code:  "VALIDATION_ERROR",
            Details: formatValidationErrors(err),
        })
    }
    
    // Sanitize input
    req.Name = html.EscapeString(strings.TrimSpace(req.Name))
    req.Email = strings.ToLower(strings.TrimSpace(req.Email))
    
    // Process request
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

### Rate Limiting

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

// Usage
framework.RegisterMiddleware(RateLimitMiddleware(100, time.Minute))
```

### CORS Configuration

```go
func CORSMiddleware(allowedOrigins []string) server.Middleware {
    return server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
        origin := ctx.Header("Origin")
        
        // Check if origin is allowed
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

### Security Headers

```go
func SecurityHeadersMiddleware() server.Middleware {
    return server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
        // Prevent XSS attacks
        ctx.SetHeader("X-XSS-Protection", "1; mode=block")
        
        // Prevent content type sniffing
        ctx.SetHeader("X-Content-Type-Options", "nosniff")
        
        // Prevent clickjacking
        ctx.SetHeader("X-Frame-Options", "DENY")
        
        // Force HTTPS
        ctx.SetHeader("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        
        // Content Security Policy
        ctx.SetHeader("Content-Security-Policy", "default-src 'self'")
        
        return next()
    })
}
```

## Performance Optimization

### Connection Pooling

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
    
    // Connection pool settings
    sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)    // Maximum open connections
    sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)    // Maximum idle connections
    sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime) // Connection lifetime
    sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime) // Connection idle time
    
    return db, nil
}
```

### Caching Strategy

```go
type CachedUserService struct {
    userRepo repository.UserRepository
    cache    cache.Cache
    ttl      time.Duration
}

func (s *CachedUserService) GetUser(id int) (*User, error) {
    // Try cache first
    key := fmt.Sprintf("user:%d", id)
    if cached := s.cache.Get(key); cached != nil {
        if user, ok := cached.(*User); ok {
            return user, nil
        }
    }
    
    // Fetch from database
    user, err := s.userRepo.GetByID(id)
    if err != nil {
        return nil, err
    }
    
    // Cache the result
    s.cache.Set(key, user, s.ttl)
    
    return user, nil
}

func (s *CachedUserService) UpdateUser(id int, updates *User) (*User, error) {
    user, err := s.userRepo.Update(id, updates)
    if err != nil {
        return nil, err
    }
    
    // Invalidate cache
    key := fmt.Sprintf("user:%d", id)
    s.cache.Delete(key)
    
    return user, nil
}
```

### Response Compression

```go
func CompressionMiddleware() server.Middleware {
    return server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
        // Check if client accepts compression
        acceptEncoding := ctx.Header("Accept-Encoding")
        if !strings.Contains(acceptEncoding, "gzip") {
            return next()
        }
        
        // Set compression headers
        ctx.SetHeader("Content-Encoding", "gzip")
        ctx.SetHeader("Vary", "Accept-Encoding")
        
        // Wrap response writer with gzip writer
        gz := gzip.NewWriter(ctx.Response())
        defer gz.Close()
        
        // Replace response writer
        originalWriter := ctx.Response()
        ctx.SetResponse(&gzipResponseWriter{
            ResponseWriter: originalWriter,
            Writer:        gz,
        })
        
        return next()
    })
}
```

## Error Handling

### Structured Error Responses

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
        
        // Generate trace ID
        traceID := generateTraceID()
        
        // Log error
        log.WithFields(log.Fields{
            "trace_id": traceID,
            "path":     ctx.Path(),
            "method":   ctx.Method(),
            "error":    err.Error(),
        }).Error("Request failed")
        
        // Handle different error types
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

### Panic Recovery

```go
func RecoveryMiddleware() server.Middleware {
    return server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
        defer func() {
            if r := recover(); r != nil {
                // Log panic with stack trace
                log.WithFields(log.Fields{
                    "panic": r,
                    "stack": debug.Stack(),
                    "path":  ctx.Path(),
                    "method": ctx.Method(),
                }).Error("Panic recovered")
                
                // Return error response
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

## Logging Best Practices

### Structured Logging

```go
func LoggingMiddleware(logger *log.Logger) server.Middleware {
    return server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
        start := time.Now()
        
        // Add request ID to context
        requestID := generateRequestID()
        ctx.Set("request_id", requestID)
        
        // Create request logger
        reqLogger := logger.WithFields(log.Fields{
            "request_id": requestID,
            "method":     ctx.Method(),
            "path":       ctx.Path(),
            "user_agent": ctx.Header("User-Agent"),
            "remote_ip":  ctx.ClientIP(),
        })
        
        // Log request start
        reqLogger.Info("Request started")
        
        // Execute next middleware
        err := next()
        
        // Calculate duration
        duration := time.Since(start)
        
        // Log request completion
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

### Log Levels

```go
// Development: Debug level for detailed information
if config.Environment == "development" {
    log.SetLevel(log.DebugLevel)
}

// Production: Info level for normal operations
if config.Environment == "production" {
    log.SetLevel(log.InfoLevel)
    log.SetFormatter(&log.JSONFormatter{})
}

// Use appropriate log levels
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

## Testing Best Practices

### Unit Testing

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

### Integration Testing

```go
func TestUserAPI_Integration(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    // Create server
    config := server.DefaultServerConfig()
    config.Port = 0 // Random port
    
    manager, err := server.CreateServerManager("gin", config)
    require.NoError(t, err)
    
    // Setup dependencies
    userRepo := repository.NewUserRepository(db)
    userService := service.NewUserService(userRepo)
    userHandler := api.NewUserHandler(userService)
    
    // Register routes
    framework := manager.GetFramework()
    api := framework.Group("/api/v1")
    api.RegisterRoute("POST", "/users", userHandler.CreateUser)
    api.RegisterRoute("GET", "/users/:id", userHandler.GetUser)
    
    // Start server
    go manager.Start(context.Background())
    defer manager.Stop(context.Background())
    
    // Wait for server to start
    time.Sleep(100 * time.Millisecond)
    
    // Test API endpoints
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

## Monitoring and Observability

### Health Checks

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
    
    // Database health
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
    
    // Redis health
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

### Metrics Collection

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

## Deployment Best Practices

### Docker Best Practices

```dockerfile
# Multi-stage build for smaller images
FROM golang:1.24-alpine AS builder

# Install dependencies
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/server .

# Copy configuration files
COPY --from=builder /app/configs ./configs

# Create non-root user
RUN adduser -D -s /bin/sh appuser
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./server"]
```

### Kubernetes Configuration

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

## Next Steps

- **[Troubleshooting](08_troubleshooting.md)** - Common issues and solutions
- **[API Reference](09_api_reference.md)** - Complete API documentation
