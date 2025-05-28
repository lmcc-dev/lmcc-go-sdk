/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * Microservice example demonstrating gRPC services with observability features.
 */

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

// ServiceConfig 微服务配置
// (ServiceConfig represents microservice configuration)
type ServiceConfig struct {
	Service struct {
		Name    string `yaml:"name" default:"user-service"`
		Version string `yaml:"version" default:"v1.0.0"`
		Port    int    `yaml:"port" default:"50051"`
	} `yaml:"service"`

	HTTP struct {
		Port            int `yaml:"port" default:"8080"`
		HealthCheckPath string `yaml:"health_check_path" default:"/health"`
		MetricsPath     string `yaml:"metrics_path" default:"/metrics"`
	} `yaml:"http"`

	Database struct {
		Host         string `yaml:"host" default:"localhost"`
		Port         int    `yaml:"port" default:"5432"`
		Name         string `yaml:"name" default:"userdb"`
		MaxConns     int    `yaml:"max_conns" default:"10"`
		ConnTimeout  int    `yaml:"conn_timeout" default:"5"`
	} `yaml:"database"`

	Logging struct {
		Level       string `yaml:"level" default:"info"`
		Format      string `yaml:"format" default:"json"`
		OutputPaths []string `yaml:"output_paths"`
	} `yaml:"logging"`

	Observability struct {
		TracingEnabled bool   `yaml:"tracing_enabled" default:"true"`
		MetricsEnabled bool   `yaml:"metrics_enabled" default:"true"`
		ServiceMesh    string `yaml:"service_mesh" default:"istio"`
	} `yaml:"observability"`
}

// User 用户数据模型
// (User represents user data model)
type User struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Status   string    `json:"status"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

// UserRequest 用户请求
// (UserRequest represents user operation request)
type UserRequest struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

// UserResponse 用户响应
// (UserResponse represents user operation response)
type UserResponse struct {
	Success bool   `json:"success"`
	User    *User  `json:"user,omitempty"`
	Error   string `json:"error,omitempty"`
	TraceID string `json:"trace_id"`
}

// UserService 用户微服务
// (UserService represents the user microservice)
type UserService struct {
	config  *ServiceConfig
	logger  log.Logger
	metrics *MetricsCollector
	tracer  *TracingService
	db      *DatabaseService
}

// NewUserService 创建用户微服务
// (NewUserService creates a new user microservice)
func NewUserService(cfg *ServiceConfig) *UserService {
	// 设置日志默认值 (Set logging defaults)
	if len(cfg.Logging.OutputPaths) == 0 {
		cfg.Logging.OutputPaths = []string{"stdout"}
	}

	// 初始化日志 (Initialize logging)
	opts := log.NewOptions()
	opts.Level = cfg.Logging.Level
	opts.Format = cfg.Logging.Format
	opts.EnableColor = cfg.Logging.Format == "text"
	opts.DisableCaller = false
	opts.DisableStacktrace = cfg.Logging.Level != "debug"
	opts.OutputPaths = cfg.Logging.OutputPaths

	log.Init(opts)

	logger := log.Std().WithValues(
		"service", cfg.Service.Name,
		"version", cfg.Service.Version,
		"component", "microservice")

	// 初始化组件 (Initialize components)
	metrics := NewMetricsCollector(logger)
	tracer := NewTracingService(cfg, logger)
	db := NewDatabaseService(cfg, logger)

	logger.Infow("User microservice initialized",
		"service_name", cfg.Service.Name,
		"service_version", cfg.Service.Version,
		"grpc_port", cfg.Service.Port,
		"http_port", cfg.HTTP.Port,
		"tracing_enabled", cfg.Observability.TracingEnabled,
		"metrics_enabled", cfg.Observability.MetricsEnabled)

	return &UserService{
		config:  cfg,
		logger:  logger,
		metrics: metrics,
		tracer:  tracer,
		db:      db,
	}
}

// GetUser 获取用户
// (GetUser retrieves a user by ID)
func (s *UserService) GetUser(ctx context.Context, req *UserRequest) (*UserResponse, error) {
	// 生成链路追踪ID (Generate trace ID)
	traceID := s.tracer.GenerateTraceID()
	ctx = s.tracer.WithTraceID(ctx, traceID)

	// 创建带追踪信息的日志记录器 (Create logger with tracing info)
	logger := s.logger.WithValues("trace_id", traceID, "operation", "get_user")

	// 记录指标 (Record metrics)
	s.metrics.RecordRequest("get_user")
	startTime := time.Now()

	logger.Infow("Processing get user request",
		"user_id", req.ID)

	// 验证请求 (Validate request)
	if req.ID == "" {
		err := errors.New("user ID is required")
		s.metrics.RecordError("get_user", "validation_error")
		logger.Errorw("Validation failed", "error", err)
		return &UserResponse{
			Success: false,
			Error:   err.Error(),
			TraceID: traceID,
		}, nil
	}

	// 从数据库获取用户 (Get user from database)
	user, err := s.db.GetUser(ctx, req.ID)
	if err != nil {
		s.metrics.RecordError("get_user", "database_error")
		logger.Errorw("Database operation failed", "error", err)
		return &UserResponse{
			Success: false,
			Error:   err.Error(),
			TraceID: traceID,
		}, nil
	}

	// 记录成功指标 (Record success metrics)
	s.metrics.RecordSuccess("get_user", time.Since(startTime))

	logger.Infow("Get user request completed successfully",
		"user_id", user.ID,
		"username", user.Username,
		"duration", time.Since(startTime))

	return &UserResponse{
		Success: true,
		User:    user,
		TraceID: traceID,
	}, nil
}

// CreateUser 创建用户
// (CreateUser creates a new user)
func (s *UserService) CreateUser(ctx context.Context, req *UserRequest) (*UserResponse, error) {
	traceID := s.tracer.GenerateTraceID()
	ctx = s.tracer.WithTraceID(ctx, traceID)

	logger := s.logger.WithValues("trace_id", traceID, "operation", "create_user")

	s.metrics.RecordRequest("create_user")
	startTime := time.Now()

	logger.Infow("Processing create user request",
		"username", req.Username,
		"email", req.Email)

	// 验证请求 (Validate request)
	if req.Username == "" || req.Email == "" {
		err := errors.New("username and email are required")
		s.metrics.RecordError("create_user", "validation_error")
		logger.Errorw("Validation failed", "error", err)
		return &UserResponse{
			Success: false,
			Error:   err.Error(),
			TraceID: traceID,
		}, nil
	}

	// 创建用户 (Create user in database)
	user, err := s.db.CreateUser(ctx, req.Username, req.Email)
	if err != nil {
		s.metrics.RecordError("create_user", "database_error")
		logger.Errorw("Database operation failed", "error", err)
		return &UserResponse{
			Success: false,
			Error:   err.Error(),
			TraceID: traceID,
		}, nil
	}

	s.metrics.RecordSuccess("create_user", time.Since(startTime))

	logger.Infow("Create user request completed successfully",
		"user_id", user.ID,
		"username", user.Username,
		"email", user.Email,
		"duration", time.Since(startTime))

	return &UserResponse{
		Success: true,
		User:    user,
		TraceID: traceID,
	}, nil
}

// MetricsCollector 指标收集器
// (MetricsCollector collects service metrics)
type MetricsCollector struct {
	logger       log.Logger
	requests     map[string]int64
	errors       map[string]map[string]int64
	latencies    map[string][]time.Duration
	mu           sync.RWMutex
}

// NewMetricsCollector 创建指标收集器
// (NewMetricsCollector creates a new metrics collector)
func NewMetricsCollector(logger log.Logger) *MetricsCollector {
	return &MetricsCollector{
		logger:    logger.WithValues("component", "metrics"),
		requests:  make(map[string]int64),
		errors:    make(map[string]map[string]int64),
		latencies: make(map[string][]time.Duration),
	}
}

// RecordRequest 记录请求
// (RecordRequest records a request)
func (mc *MetricsCollector) RecordRequest(operation string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.requests[operation]++
	mc.logger.Debugw("Request recorded", "operation", operation, "count", mc.requests[operation])
}

// RecordError 记录错误
// (RecordError records an error)
func (mc *MetricsCollector) RecordError(operation, errorType string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	if mc.errors[operation] == nil {
		mc.errors[operation] = make(map[string]int64)
	}
	mc.errors[operation][errorType]++
	mc.logger.Warnw("Error recorded", "operation", operation, "error_type", errorType)
}

// RecordSuccess 记录成功响应
// (RecordSuccess records a successful response)
func (mc *MetricsCollector) RecordSuccess(operation string, duration time.Duration) {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.latencies[operation] = append(mc.latencies[operation], duration)
	mc.logger.Debugw("Success recorded", "operation", operation, "duration", duration)
}

// GetMetrics 获取指标
// (GetMetrics returns current metrics)
func (mc *MetricsCollector) GetMetrics() map[string]interface{} {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	metrics := map[string]interface{}{
		"requests":  make(map[string]int64),
		"errors":    make(map[string]map[string]int64),
		"latencies": make(map[string]map[string]interface{}),
	}

	// 复制请求计数 (Copy request counts)
	for op, count := range mc.requests {
		metrics["requests"].(map[string]int64)[op] = count
	}

	// 复制错误计数 (Copy error counts)
	for op, errors := range mc.errors {
		metrics["errors"].(map[string]map[string]int64)[op] = make(map[string]int64)
		for errorType, count := range errors {
			metrics["errors"].(map[string]map[string]int64)[op][errorType] = count
		}
	}

	// 计算延迟统计 (Calculate latency statistics)
	for op, durations := range mc.latencies {
		if len(durations) > 0 {
			var total time.Duration
			min := durations[0]
			max := durations[0]

			for _, d := range durations {
				total += d
				if d < min {
					min = d
				}
				if d > max {
					max = d
				}
			}

			avg := total / time.Duration(len(durations))
			metrics["latencies"].(map[string]map[string]interface{})[op] = map[string]interface{}{
				"count":   len(durations),
				"avg":     avg,
				"min":     min,
				"max":     max,
			}
		}
	}

	return metrics
}

// TracingService 链路追踪服务
// (TracingService provides distributed tracing)
type TracingService struct {
	config *ServiceConfig
	logger log.Logger
}

// NewTracingService 创建链路追踪服务
// (NewTracingService creates a new tracing service)
func NewTracingService(cfg *ServiceConfig, logger log.Logger) *TracingService {
	return &TracingService{
		config: cfg,
		logger: logger.WithValues("component", "tracing"),
	}
}

// GenerateTraceID 生成链路追踪ID
// (GenerateTraceID generates a new trace ID)
func (ts *TracingService) GenerateTraceID() string {
	traceID := fmt.Sprintf("trace_%d_%d", time.Now().UnixNano(), os.Getpid())
	ts.logger.Debugw("Generated trace ID", "trace_id", traceID)
	return traceID
}

// WithTraceID 在上下文中添加追踪ID
// (WithTraceID adds trace ID to context)
func (ts *TracingService) WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, "trace_id", traceID)
}

// GetTraceID 从上下文获取追踪ID
// (GetTraceID retrieves trace ID from context)
func (ts *TracingService) GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value("trace_id").(string); ok {
		return traceID
	}
	return ""
}

// DatabaseService 数据库服务
// (DatabaseService provides database operations)
type DatabaseService struct {
	config *ServiceConfig
	logger log.Logger
	users  map[string]*User
	mu     sync.RWMutex
}

// NewDatabaseService 创建数据库服务
// (NewDatabaseService creates a new database service)
func NewDatabaseService(cfg *ServiceConfig, logger log.Logger) *DatabaseService {
	ds := &DatabaseService{
		config: cfg,
		logger: logger.WithValues("component", "database"),
		users:  make(map[string]*User),
	}

	// 预填充一些测试数据 (Pre-populate with test data)
	ds.seedData()

	return ds
}

// seedData 填充测试数据
// (seedData populates test data)
func (ds *DatabaseService) seedData() {
	testUsers := []*User{
		{
			ID:       "user_001",
			Username: "alice",
			Email:    "alice@example.com",
			Status:   "active",
			Created:  time.Now().Add(-30 * 24 * time.Hour),
			Updated:  time.Now().Add(-1 * time.Hour),
		},
		{
			ID:       "user_002",
			Username: "bob",
			Email:    "bob@example.com",
			Status:   "active",
			Created:  time.Now().Add(-15 * 24 * time.Hour),
			Updated:  time.Now().Add(-2 * time.Hour),
		},
		{
			ID:       "user_003",
			Username: "charlie",
			Email:    "charlie@example.com",
			Status:   "inactive",
			Created:  time.Now().Add(-7 * 24 * time.Hour),
			Updated:  time.Now().Add(-3 * time.Hour),
		},
	}

	for _, user := range testUsers {
		ds.users[user.ID] = user
	}

	ds.logger.Infow("Database seeded with test data", "user_count", len(testUsers))
}

// GetUser 获取用户
// (GetUser retrieves a user by ID)
func (ds *DatabaseService) GetUser(ctx context.Context, userID string) (*User, error) {
	traceID := ds.getTraceID(ctx)
	logger := ds.logger.WithValues("trace_id", traceID, "operation", "db_get_user")

	start := time.Now()

	logger.Debugw("Database query started", "user_id", userID)

	// 模拟数据库延迟 (Simulate database latency)
	time.Sleep(20 * time.Millisecond)

	ds.mu.RLock()
	user, exists := ds.users[userID]
	ds.mu.RUnlock()

	duration := time.Since(start)

	if !exists {
		err := errors.New("user not found")
		logger.Warnw("User not found", "user_id", userID, "duration", duration)
		return nil, err
	}

	// 返回用户副本 (Return copy of user)
	userCopy := *user
	logger.Debugw("Database query completed", "user_id", userID, "duration", duration)

	return &userCopy, nil
}

// CreateUser 创建用户
// (CreateUser creates a new user)
func (ds *DatabaseService) CreateUser(ctx context.Context, username, email string) (*User, error) {
	traceID := ds.getTraceID(ctx)
	logger := ds.logger.WithValues("trace_id", traceID, "operation", "db_create_user")

	start := time.Now()

	logger.Debugw("Database insert started", "username", username, "email", email)

	// 模拟数据库延迟 (Simulate database latency)
	time.Sleep(50 * time.Millisecond)

	// 检查用户名是否已存在 (Check if username already exists)
	ds.mu.RLock()
	for _, user := range ds.users {
		if user.Username == username {
			ds.mu.RUnlock()
			err := errors.New("username already exists")
			logger.Warnw("Username conflict", "username", username)
			return nil, err
		}
		if user.Email == email {
			ds.mu.RUnlock()
			err := errors.New("email already exists")
			logger.Warnw("Email conflict", "email", email)
			return nil, err
		}
	}
	ds.mu.RUnlock()

	// 创建新用户 (Create new user)
	user := &User{
		ID:       fmt.Sprintf("user_%d", time.Now().Unix()),
		Username: username,
		Email:    email,
		Status:   "active",
		Created:  time.Now(),
		Updated:  time.Now(),
	}

	ds.mu.Lock()
	ds.users[user.ID] = user
	ds.mu.Unlock()

	duration := time.Since(start)
	logger.Debugw("Database insert completed",
		"user_id", user.ID,
		"username", username,
		"duration", duration)

	return user, nil
}

// getTraceID 从上下文获取追踪ID
// (getTraceID retrieves trace ID from context)
func (ds *DatabaseService) getTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value("trace_id").(string); ok {
		return traceID
	}
	return "unknown"
}

// HealthChecker 健康检查器
// (HealthChecker provides health check functionality)
type HealthChecker struct {
	service *UserService
	logger  log.Logger
}

// NewHealthChecker 创建健康检查器
// (NewHealthChecker creates a new health checker)
func NewHealthChecker(service *UserService) *HealthChecker {
	return &HealthChecker{
		service: service,
		logger:  service.logger.WithValues("component", "health"),
	}
}

// Check 执行健康检查
// (Check performs health check)
func (hc *HealthChecker) Check(ctx context.Context) map[string]interface{} {
	hc.logger.Debugw("Performing health check")

	health := map[string]interface{}{
		"status":     "healthy",
		"timestamp":  time.Now(),
		"service":    hc.service.config.Service.Name,
		"version":    hc.service.config.Service.Version,
		"checks":     make(map[string]interface{}),
	}

	// 检查数据库连接 (Check database connection)
	dbHealth := hc.checkDatabase(ctx)
	health["checks"].(map[string]interface{})["database"] = dbHealth

	// 检查内存使用 (Check memory usage)
	memHealth := hc.checkMemory()
	health["checks"].(map[string]interface{})["memory"] = memHealth

	// 检查服务响应时间 (Check service response time)
	responseHealth := hc.checkResponseTime(ctx)
	health["checks"].(map[string]interface{})["response_time"] = responseHealth

	// 判断整体健康状态 (Determine overall health status)
	if !dbHealth["healthy"].(bool) || !memHealth["healthy"].(bool) || !responseHealth["healthy"].(bool) {
		health["status"] = "unhealthy"
	}

	hc.logger.Infow("Health check completed", "status", health["status"])

	return health
}

// checkDatabase 检查数据库健康状态
// (checkDatabase checks database health)
func (hc *HealthChecker) checkDatabase(ctx context.Context) map[string]interface{} {
	start := time.Now()

	// 尝试获取一个测试用户 (Try to get a test user)
	_, err := hc.service.db.GetUser(ctx, "user_001")

	duration := time.Since(start)
	healthy := err == nil && duration < 100*time.Millisecond

	return map[string]interface{}{
		"healthy":       healthy,
		"response_time": duration,
		"error":         func() string { if err != nil { return err.Error() }; return "" }(),
	}
}

// checkMemory 检查内存使用
// (checkMemory checks memory usage)
func (hc *HealthChecker) checkMemory() map[string]interface{} {
	// 这里可以添加实际的内存检查逻辑 (Add actual memory checking logic here)
	// 目前返回模拟数据 (Currently returning mock data)
	return map[string]interface{}{
		"healthy":    true,
		"usage_mb":   156.7,
		"limit_mb":   512.0,
		"usage_pct":  30.6,
	}
}

// checkResponseTime 检查服务响应时间
// (checkResponseTime checks service response time)
func (hc *HealthChecker) checkResponseTime(ctx context.Context) map[string]interface{} {
	start := time.Now()

	// 执行一个轻量级操作 (Perform a lightweight operation)
	req := &UserRequest{ID: "user_001"}
	_, err := hc.service.GetUser(ctx, req)

	duration := time.Since(start)
	healthy := err == nil && duration < 200*time.Millisecond

	return map[string]interface{}{
		"healthy":       healthy,
		"response_time": duration,
		"threshold":     "200ms",
	}
}

// HTTPServer HTTP服务器
// (HTTPServer provides HTTP endpoints)
type HTTPServer struct {
	service       *UserService
	healthChecker *HealthChecker
	logger        log.Logger
}

// NewHTTPServer 创建HTTP服务器
// (NewHTTPServer creates a new HTTP server)
func NewHTTPServer(service *UserService) *HTTPServer {
	return &HTTPServer{
		service:       service,
		healthChecker: NewHealthChecker(service),
		logger:        service.logger.WithValues("component", "http"),
	}
}

// Start 启动HTTP服务器
// (Start starts the HTTP server)
func (hs *HTTPServer) Start() error {
	mux := http.NewServeMux()

	// 健康检查端点 (Health check endpoint)
	mux.HandleFunc(hs.service.config.HTTP.HealthCheckPath, hs.healthHandler)

	// 指标端点 (Metrics endpoint)
	mux.HandleFunc(hs.service.config.HTTP.MetricsPath, hs.metricsHandler)

	// 用户API端点 (User API endpoints)
	mux.HandleFunc("/api/users/", hs.getUserHandler)
	mux.HandleFunc("/api/users", hs.createUserHandler)

	addr := fmt.Sprintf(":%d", hs.service.config.HTTP.Port)
	hs.logger.Infow("Starting HTTP server", "address", addr)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return server.ListenAndServe()
}

// healthHandler 健康检查处理器
// (healthHandler handles health check requests)
func (hs *HTTPServer) healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	health := hs.healthChecker.Check(ctx)

	w.Header().Set("Content-Type", "application/json")
	if health["status"] == "healthy" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	fmt.Fprintf(w, `{"status": "%s", "timestamp": "%s", "service": "%s", "version": "%s"}`,
		health["status"], health["timestamp"], health["service"], health["version"])
}

// metricsHandler 指标处理器
// (metricsHandler handles metrics requests)
func (hs *HTTPServer) metricsHandler(w http.ResponseWriter, r *http.Request) {
	metrics := hs.service.metrics.GetMetrics()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// 简化的指标输出 (Simplified metrics output)
	fmt.Fprintf(w, `{"requests": %v, "timestamp": "%s"}`,
		metrics["requests"], time.Now().Format(time.RFC3339))
}

// getUserHandler 获取用户处理器
// (getUserHandler handles get user requests)
func (hs *HTTPServer) getUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 从URL提取用户ID (Extract user ID from URL)
	userID := r.URL.Path[len("/api/users/"):]
	if userID == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	req := &UserRequest{ID: userID}
	resp, err := hs.service.GetUser(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if resp.Success {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprintf(w, `{"success": %t, "user": %v, "trace_id": "%s"}`,
		resp.Success, resp.User, resp.TraceID)
}

// createUserHandler 创建用户处理器
// (createUserHandler handles create user requests)
func (hs *HTTPServer) createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 简化的请求解析 (Simplified request parsing)
	username := r.FormValue("username")
	email := r.FormValue("email")

	req := &UserRequest{Username: username, Email: email}
	resp, err := hs.service.CreateUser(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if resp.Success {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	fmt.Fprintf(w, `{"success": %t, "user": %v, "trace_id": "%s"}`,
		resp.Success, resp.User, resp.TraceID)
}

// runDemo 运行演示
// (runDemo runs the microservice demonstration)
func runDemo(service *UserService) {
	fmt.Println("=== Running Microservice Demonstration ===")
	fmt.Println()

	ctx := context.Background()

	// 演示用户操作 (Demonstrate user operations)
	fmt.Println("1. Testing User Operations:")

	// 获取存在的用户 (Get existing user)
	getUserReq := &UserRequest{ID: "user_001"}
	resp, err := service.GetUser(ctx, getUserReq)
	if err != nil {
		fmt.Printf("   ❌ Get user failed: %v\n", err)
	} else if resp.Success {
		fmt.Printf("   ✅ Get user successful: %s (%s)\n", resp.User.Username, resp.TraceID)
	} else {
		fmt.Printf("   ❌ Get user failed: %s\n", resp.Error)
	}

	// 获取不存在的用户 (Get non-existent user)
	getUserReq2 := &UserRequest{ID: "user_999"}
	resp2, err := service.GetUser(ctx, getUserReq2)
	if err != nil {
		fmt.Printf("   ❌ Get non-existent user failed: %v\n", err)
	} else if !resp2.Success {
		fmt.Printf("   ✅ Get non-existent user correctly failed: %s\n", resp2.Error)
	}

	// 创建新用户 (Create new user)
	createReq := &UserRequest{Username: "david", Email: "david@example.com"}
	resp3, err := service.CreateUser(ctx, createReq)
	if err != nil {
		fmt.Printf("   ❌ Create user failed: %v\n", err)
	} else if resp3.Success {
		fmt.Printf("   ✅ Create user successful: %s (%s)\n", resp3.User.ID, resp3.TraceID)
	} else {
		fmt.Printf("   ❌ Create user failed: %s\n", resp3.Error)
	}

	fmt.Println()

	// 演示健康检查 (Demonstrate health check)
	fmt.Println("2. Testing Health Check:")
	healthChecker := NewHealthChecker(service)
	health := healthChecker.Check(ctx)
	fmt.Printf("   Service Status: %s\n", health["status"])
	fmt.Printf("   Timestamp: %s\n", health["timestamp"])
	fmt.Println()

	// 显示指标 (Show metrics)
	fmt.Println("3. Service Metrics:")
	metrics := service.metrics.GetMetrics()
	requests := metrics["requests"].(map[string]int64)
	for operation, count := range requests {
		fmt.Printf("   %s: %d requests\n", operation, count)
	}
	fmt.Println()

	fmt.Println("=== Microservice Demonstration Completed ===")
}

func main() {
	fmt.Println("=== Microservice Integration Example ===")
	fmt.Println("This example demonstrates a microservice with gRPC, observability, and health checks.")
	fmt.Println()

	// 加载配置 (Load configuration)
	cfg := &ServiceConfig{}
	if err := config.LoadConfig(cfg); err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		// 使用默认配置继续 (Continue with default configuration)
	}

	// 创建用户微服务 (Create user microservice)
	userService := NewUserService(cfg)

	// 运行演示 (Run demonstration)
	runDemo(userService)

	// 可选：启动HTTP服务器进行交互测试 (Optional: start HTTP server for interactive testing)
	fmt.Println("Starting HTTP server for additional testing...")
	httpServer := NewHTTPServer(userService)

	// 在goroutine中启动HTTP服务器 (Start HTTP server in goroutine)
	go func() {
		if err := httpServer.Start(); err != nil {
			userService.logger.Errorw("HTTP server failed", "error", err)
		}
	}()

	// 等待一段时间让服务器启动 (Wait for server to start)
	time.Sleep(1 * time.Second)

	fmt.Printf("HTTP server running on port %d\n", cfg.HTTP.Port)
	fmt.Println("Available endpoints:")
	fmt.Printf("  GET  http://localhost:%d%s\n", cfg.HTTP.Port, cfg.HTTP.HealthCheckPath)
	fmt.Printf("  GET  http://localhost:%d%s\n", cfg.HTTP.Port, cfg.HTTP.MetricsPath)
	fmt.Printf("  GET  http://localhost:%d/api/users/user_001\n", cfg.HTTP.Port)
	fmt.Println()

	// 运行一些HTTP测试 (Run some HTTP tests)
	runHTTPTests(cfg)

	userService.logger.Infow("Microservice example completed successfully")
	fmt.Println("=== Example completed successfully ===")
}

// runHTTPTests 运行HTTP测试
// (runHTTPTests runs HTTP endpoint tests)
func runHTTPTests(cfg *ServiceConfig) {
	fmt.Println("=== Running HTTP Endpoint Tests ===")
	fmt.Println()

	baseURL := fmt.Sprintf("http://localhost:%d", cfg.HTTP.Port)

	// 测试健康检查端点 (Test health check endpoint)
	fmt.Println("Testing health check endpoint:")
	resp, err := http.Get(baseURL + cfg.HTTP.HealthCheckPath)
	if err != nil {
		fmt.Printf("   ❌ Health check failed: %v\n", err)
	} else {
		fmt.Printf("   ✅ Health check successful: %d\n", resp.StatusCode)
		resp.Body.Close()
	}

	// 测试指标端点 (Test metrics endpoint)
	fmt.Println("Testing metrics endpoint:")
	resp, err = http.Get(baseURL + cfg.HTTP.MetricsPath)
	if err != nil {
		fmt.Printf("   ❌ Metrics endpoint failed: %v\n", err)
	} else {
		fmt.Printf("   ✅ Metrics endpoint successful: %d\n", resp.StatusCode)
		resp.Body.Close()
	}

	// 测试用户API (Test user API)
	fmt.Println("Testing user API:")
	resp, err = http.Get(baseURL + "/api/users/user_001")
	if err != nil {
		fmt.Printf("   ❌ User API failed: %v\n", err)
	} else {
		fmt.Printf("   ✅ User API successful: %d\n", resp.StatusCode)
		resp.Body.Close()
	}

	fmt.Println()
	fmt.Println("=== HTTP Tests Completed ===")
	fmt.Println()
} 