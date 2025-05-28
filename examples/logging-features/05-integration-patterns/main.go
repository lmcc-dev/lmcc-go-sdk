/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * Integration patterns example demonstrating various logging integration patterns.
 */

package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

// AppConfig 应用配置结构
// (AppConfig represents application configuration structure)
type AppConfig struct {
	Server struct {
		Port         int    `yaml:"port" default:"8080"`
		ReadTimeout  int    `yaml:"read_timeout" default:"30"`
		WriteTimeout int    `yaml:"write_timeout" default:"30"`
	} `yaml:"server"`
	
	Database struct {
		Host         string `yaml:"host" default:"localhost"`
		Port         int    `yaml:"port" default:"5432"`
		Name         string `yaml:"name" default:"appdb"`
		MaxConns     int    `yaml:"max_conns" default:"10"`
		ConnTimeout  int    `yaml:"conn_timeout" default:"5"`
	} `yaml:"database"`
	
	Logging struct {
		Level  string `yaml:"level" default:"info"`
		Format string `yaml:"format" default:"json"`
		Output string `yaml:"output" default:"stdout"`
	} `yaml:"logging"`
}

// LoggingMiddleware HTTP日志中间件
// (LoggingMiddleware provides HTTP logging middleware)
type LoggingMiddleware struct {
	logger log.Logger
}

// NewLoggingMiddleware 创建日志中间件
// (NewLoggingMiddleware creates a new logging middleware)
func NewLoggingMiddleware(logger log.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

// Handler 中间件处理函数
// (Handler middleware handler function)
func (m *LoggingMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// 创建请求ID (Create request ID)
		requestID := fmt.Sprintf("req_%d", time.Now().UnixNano())
		
		// 创建带请求信息的日志记录器 (Create logger with request info)
		requestLogger := m.logger.WithValues(
			"request_id", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent())
		
		// 记录请求开始 (Log request start)
		requestLogger.Infow("HTTP request started",
			"query_params", r.URL.RawQuery,
			"content_length", r.ContentLength)
		
		// 创建响应写入器包装器 (Create response writer wrapper)
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		// 将日志记录器添加到上下文 (Add logger to context)
		ctx := context.WithValue(r.Context(), "logger", requestLogger)
		ctx = context.WithValue(ctx, "request_id", requestID)
		
		// 执行下一个处理器 (Execute next handler)
		next.ServeHTTP(wrapped, r.WithContext(ctx))
		
		// 记录请求完成 (Log request completion)
		duration := time.Since(start)
		
		logLevel := requestLogger.Infow
		if wrapped.statusCode >= 400 && wrapped.statusCode < 500 {
			logLevel = requestLogger.Warnw
		} else if wrapped.statusCode >= 500 {
			logLevel = requestLogger.Errorw
		}
		
		logLevel("HTTP request completed",
			"status_code", wrapped.statusCode,
			"duration", duration,
			"response_size", wrapped.bytesWritten)
	})
}

// responseWriter 响应写入器包装器
// (responseWriter wrapper for capturing response details)
type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(data)
	rw.bytesWritten += n
	return n, err
}

// UserService 用户服务（集成日志和错误处理）
// (UserService with integrated logging and error handling)
type UserService struct {
	logger log.Logger
	config *AppConfig
}

// NewUserService 创建用户服务
// (NewUserService creates a new user service)
func NewUserService(logger log.Logger, cfg *AppConfig) *UserService {
	return &UserService{
		logger: logger.WithValues("component", "user-service"),
		config: cfg,
	}
}

// CreateUser 创建用户（集成日志和错误处理）
// (CreateUser creates a user with integrated logging and error handling)
func (s *UserService) CreateUser(ctx context.Context, userID, email string) error {
	// 从上下文获取请求ID (Get request ID from context)
	requestID := ctx.Value("request_id")
	
	// 创建带上下文的日志记录器 (Create context-aware logger)
	logger := s.logger.WithValues("request_id", requestID, "operation", "create_user")
	
	logger.Infow("Starting user creation",
		"user_id", userID,
		"email", email)
	
	// 验证输入 (Validate input)
	if userID == "" {
		err := errors.New("user ID cannot be empty")
		logger.Errorw("User creation validation failed",
			"error", err,
			"validation_field", "user_id")
		return err
	}
	
	if email == "" {
		err := errors.New("email cannot be empty")
		logger.Errorw("User creation validation failed",
			"error", err,
			"validation_field", "email")
		return err
	}
	
	// 模拟数据库操作 (Simulate database operation)
	if err := s.simulateDBOperation(ctx, "INSERT", "users", map[string]interface{}{
		"user_id": userID,
		"email":   email,
	}); err != nil {
		// 使用集成的错误处理 (Use integrated error handling)
		wrappedErr := errors.Wrap(err, "failed to create user in database")
		logger.Errorw("Database operation failed",
			"error", wrappedErr,
			"operation", "INSERT",
			"table", "users")
		return wrappedErr
	}
	
	logger.Infow("User created successfully",
		"user_id", userID,
		"email", email)
	
	return nil
}

// GetUser 获取用户
// (GetUser retrieves a user)
func (s *UserService) GetUser(ctx context.Context, userID string) (map[string]interface{}, error) {
	requestID := ctx.Value("request_id")
	logger := s.logger.WithValues("request_id", requestID, "operation", "get_user")
	
	logger.Debugw("Starting user retrieval",
		"user_id", userID)
	
	if userID == "" {
		err := errors.New("user ID cannot be empty")
		logger.Errorw("User retrieval validation failed",
			"error", err)
		return nil, err
	}
	
	// 模拟数据库查询 (Simulate database query)
	userData := map[string]interface{}{
		"user_id": userID,
		"email":   fmt.Sprintf("%s@example.com", userID),
		"status":  "active",
		"created_at": time.Now().Add(-30 * 24 * time.Hour).Format(time.RFC3339),
	}
	
	if err := s.simulateDBOperation(ctx, "SELECT", "users", map[string]interface{}{
		"user_id": userID,
	}); err != nil {
		wrappedErr := errors.Wrap(err, "failed to retrieve user from database")
		logger.Errorw("Database operation failed",
			"error", wrappedErr,
			"operation", "SELECT",
			"table", "users")
		return nil, wrappedErr
	}
	
	logger.Infow("User retrieved successfully",
		"user_id", userID)
	
	return userData, nil
}

// simulateDBOperation 模拟数据库操作
// (simulateDBOperation simulates database operation)
func (s *UserService) simulateDBOperation(ctx context.Context, operation, table string, params map[string]interface{}) error {
	requestID := ctx.Value("request_id")
	logger := s.logger.WithValues("request_id", requestID, "component", "database")
	
	start := time.Now()
	
	logger.Debugw("Database operation started",
		"operation", operation,
		"table", table,
		"params", params)
	
	// 模拟数据库延迟 (Simulate database latency)
	time.Sleep(time.Duration(50+len(params)*10) * time.Millisecond)
	
	duration := time.Since(start)
	
	// 模拟一些错误情况 (Simulate some error conditions)
	if table == "users" && operation == "INSERT" {
		if email, ok := params["email"].(string); ok && email == "duplicate@example.com" {
			err := errors.New("duplicate email address")
			logger.Errorw("Database constraint violation",
				"error", err,
				"constraint", "unique_email",
				"duration", duration)
			return err
		}
	}
	
	if duration > 100*time.Millisecond {
		logger.Warnw("Slow database operation detected",
			"operation", operation,
			"table", table,
			"duration", duration,
			"threshold", "100ms")
	} else {
		logger.Debugw("Database operation completed",
			"operation", operation,
			"table", table,
			"duration", duration)
	}
	
	return nil
}

// ConfigurationIntegration 配置集成示例
// (ConfigurationIntegration demonstrates configuration integration)
type ConfigurationIntegration struct {
	logger log.Logger
	config *AppConfig
}

// NewConfigurationIntegration 创建配置集成示例
// (NewConfigurationIntegration creates configuration integration example)
func NewConfigurationIntegration(cfg *AppConfig) *ConfigurationIntegration {
	// 根据配置初始化日志 (Initialize logging based on configuration)
	opts := log.NewOptions()
	opts.Level = cfg.Logging.Level
	opts.Format = cfg.Logging.Format
	opts.EnableColor = cfg.Logging.Format == "text"
	opts.DisableCaller = false
	opts.DisableStacktrace = cfg.Logging.Level != "debug"
	opts.OutputPaths = []string{cfg.Logging.Output}
	
	log.Init(opts)
	
	logger := log.Std().WithValues(
		"service", "configuration-integration",
		"version", "v1.0.0")
	
	logger.Infow("Configuration-based logging initialized",
		"log_level", cfg.Logging.Level,
		"log_format", cfg.Logging.Format,
		"log_output", cfg.Logging.Output)
	
	return &ConfigurationIntegration{
		logger: logger,
		config: cfg,
	}
}

// ProcessConfiguration 处理配置
// (ProcessConfiguration processes configuration)
func (ci *ConfigurationIntegration) ProcessConfiguration(ctx context.Context) {
	ci.logger.Infow("Processing application configuration",
		"server_port", ci.config.Server.Port,
		"db_host", ci.config.Database.Host,
		"db_port", ci.config.Database.Port,
		"max_connections", ci.config.Database.MaxConns)
	
	// 验证配置 (Validate configuration)
	if ci.config.Server.Port < 1024 || ci.config.Server.Port > 65535 {
		ci.logger.Warnw("Invalid server port configuration",
			"port", ci.config.Server.Port,
			"valid_range", "1024-65535")
	}
	
	if ci.config.Database.MaxConns > 100 {
		ci.logger.Warnw("High database connection pool size",
			"max_conns", ci.config.Database.MaxConns,
			"recommendation", "consider reducing for better resource management")
	}
	
	ci.logger.Infow("Configuration processing completed",
		"config_valid", true)
}

// ErrorHandlingIntegration 错误处理集成示例
// (ErrorHandlingIntegration demonstrates error handling integration)
type ErrorHandlingIntegration struct {
	logger log.Logger
}

// NewErrorHandlingIntegration 创建错误处理集成示例
// (NewErrorHandlingIntegration creates error handling integration example)
func NewErrorHandlingIntegration(logger log.Logger) *ErrorHandlingIntegration {
	return &ErrorHandlingIntegration{
		logger: logger.WithValues("component", "error-handling"),
	}
}

// ProcessWithErrorHandling 带错误处理的处理过程
// (ProcessWithErrorHandling processing with error handling)
func (ehi *ErrorHandlingIntegration) ProcessWithErrorHandling(ctx context.Context, data map[string]interface{}) error {
	requestID := ctx.Value("request_id")
	logger := ehi.logger.WithValues("request_id", requestID)
	
	logger.Infow("Starting error handling integration demo",
		"data_size", len(data))
	
	// 模拟多层错误处理 (Simulate multi-layer error handling)
	if err := ehi.validateData(ctx, data); err != nil {
		wrappedErr := errors.Wrap(err, "data validation failed")
		logger.Errorw("Validation error occurred",
			"error", wrappedErr,
			"step", "validation")
		return wrappedErr
	}
	
	if err := ehi.processData(ctx, data); err != nil {
		wrappedErr := errors.Wrap(err, "data processing failed")
		logger.Errorw("Processing error occurred",
			"error", wrappedErr,
			"step", "processing")
		return wrappedErr
	}
	
	if err := ehi.persistData(ctx, data); err != nil {
		wrappedErr := errors.Wrap(err, "data persistence failed")
		logger.Errorw("Persistence error occurred",
			"error", wrappedErr,
			"step", "persistence")
		return wrappedErr
	}
	
	logger.Infow("Error handling integration demo completed successfully")
	return nil
}

// validateData 验证数据
// (validateData validates data)
func (ehi *ErrorHandlingIntegration) validateData(ctx context.Context, data map[string]interface{}) error {
	logger := ehi.logger.WithValues("operation", "validate")
	
	if len(data) == 0 {
		err := errors.New("data cannot be empty")
		logger.Errorw("Data validation failed", "error", err)
		return err
	}
	
	if _, ok := data["id"]; !ok {
		err := errors.New("missing required field: id")
		logger.Errorw("Data validation failed", "error", err)
		return err
	}
	
	logger.Debugw("Data validation completed successfully")
	return nil
}

// processData 处理数据
// (processData processes data)
func (ehi *ErrorHandlingIntegration) processData(ctx context.Context, data map[string]interface{}) error {
	logger := ehi.logger.WithValues("operation", "process")
	
	// 模拟处理错误 (Simulate processing error)
	if data["id"] == "error_case" {
		err := errors.New("processing error for error_case")
		logger.Errorw("Data processing failed", "error", err)
		return err
	}
	
	logger.Debugw("Data processing completed successfully")
	return nil
}

// persistData 持久化数据
// (persistData persists data)
func (ehi *ErrorHandlingIntegration) persistData(ctx context.Context, data map[string]interface{}) error {
	logger := ehi.logger.WithValues("operation", "persist")
	
	// 模拟持久化错误 (Simulate persistence error)
	if data["id"] == "persist_error" {
		err := errors.New("database connection timeout")
		logger.Errorw("Data persistence failed", "error", err)
		return err
	}
	
	logger.Debugw("Data persistence completed successfully")
	return nil
}

// demonstrateMiddlewareIntegration 演示中间件集成
// (demonstrateMiddlewareIntegration demonstrates middleware integration)
func demonstrateMiddlewareIntegration() {
	fmt.Println("=== Demonstrating Middleware Integration ===")
	fmt.Println()
	
	// 初始化日志 (Initialize logging)
	opts := log.NewOptions()
	opts.Level = "info"
	opts.Format = "json"
	opts.EnableColor = false
	opts.DisableCaller = false
	opts.DisableStacktrace = true
	opts.OutputPaths = []string{"stdout"}
	
	log.Init(opts)
	
	logger := log.Std().WithValues("component", "http-server")
	
	// 创建中间件 (Create middleware)
	loggingMiddleware := NewLoggingMiddleware(logger)
	
	// 创建HTTP处理器 (Create HTTP handler)
	mux := http.NewServeMux()
	
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		// 从上下文获取日志记录器 (Get logger from context)
		if ctxLogger, ok := r.Context().Value("logger").(log.Logger); ok {
			ctxLogger.Infow("Processing API request",
				"endpoint", "/api/users",
				"handler", "users")
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Users endpoint"}`))
	})
	
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		if ctxLogger, ok := r.Context().Value("logger").(log.Logger); ok {
			ctxLogger.Infow("Health check requested")
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy"}`))
	})
	
	// 应用中间件 (Apply middleware)
	handler := loggingMiddleware.Handler(mux)
	
	// 模拟HTTP请求 (Simulate HTTP requests)
	fmt.Println("Simulating HTTP requests with logging middleware:")
	
	requests := []struct {
		method string
		path   string
		status int
	}{
		{"GET", "/api/users", 200},
		{"POST", "/api/users", 201},
		{"GET", "/api/health", 200},
		{"GET", "/api/nonexistent", 404},
	}
	
	for _, req := range requests {
		fmt.Printf("Processing: %s %s\n", req.method, req.path)
		
		// 这里只是演示，实际中会有真正的HTTP请求 (This is just demonstration, real HTTP requests would be made)
		// 创建模拟请求 (Create mock request)
		r, _ := http.NewRequest(req.method, req.path, nil)
		r.RemoteAddr = "192.168.1.100:12345"
		r.Header.Set("User-Agent", "Integration-Test/1.0")
		
		// 创建模拟响应写入器 (Create mock response writer)
		w := &mockResponseWriter{statusCode: req.status}
		
		// 处理请求 (Handle request)
		handler.ServeHTTP(w, r)
		fmt.Println()
	}
}

// mockResponseWriter 模拟响应写入器
// (mockResponseWriter mock response writer for demonstration)
type mockResponseWriter struct {
	statusCode   int
	headers      http.Header
	bytesWritten int
}

func (m *mockResponseWriter) Header() http.Header {
	if m.headers == nil {
		m.headers = make(http.Header)
	}
	return m.headers
}

func (m *mockResponseWriter) Write(data []byte) (int, error) {
	m.bytesWritten += len(data)
	return len(data), nil
}

func (m *mockResponseWriter) WriteHeader(statusCode int) {
	m.statusCode = statusCode
}

// demonstrateConfigurationIntegration 演示配置集成
// (demonstrateConfigurationIntegration demonstrates configuration integration)
func demonstrateConfigurationIntegration() {
	fmt.Println("=== Demonstrating Configuration Integration ===")
	fmt.Println()
	
	// 创建配置 (Create configuration)
	cfg := &AppConfig{}
	config.LoadConfig(cfg)
	
	// 基于配置创建集成示例 (Create integration example based on configuration)
	ci := NewConfigurationIntegration(cfg)
	
	ctx := context.Background()
	ci.ProcessConfiguration(ctx)
	
	fmt.Println()
}

// demonstrateServiceIntegration 演示服务集成
// (demonstrateServiceIntegration demonstrates service integration)
func demonstrateServiceIntegration() {
	fmt.Println("=== Demonstrating Service Integration ===")
	fmt.Println()
	
	// 初始化日志 (Initialize logging)
	opts := log.NewOptions()
	opts.Level = "debug"
	opts.Format = "json"
	opts.EnableColor = false
	opts.DisableCaller = false
	opts.DisableStacktrace = false
	opts.OutputPaths = []string{"stdout"}
	
	log.Init(opts)
	
	logger := log.Std().WithValues("service", "user-management")
	
	// 创建配置 (Create configuration)
	cfg := &AppConfig{}
	config.LoadConfig(cfg)
	
	// 创建用户服务 (Create user service)
	userService := NewUserService(logger, cfg)
	
	// 模拟服务操作 (Simulate service operations)
	operations := []struct {
		name     string
		userID   string
		email    string
		shouldFail bool
	}{
		{"Successful user creation", "user123", "user123@example.com", false},
		{"User creation with duplicate email", "user456", "duplicate@example.com", true},
		{"User creation with empty ID", "", "test@example.com", true},
		{"User retrieval", "user123", "", false},
	}
	
	for _, op := range operations {
		fmt.Printf("Operation: %s\n", op.name)
		
		ctx := context.WithValue(context.Background(), "request_id", fmt.Sprintf("req_%d", time.Now().UnixNano()))
		
		if op.email != "" {
			// 创建用户操作 (Create user operation)
			err := userService.CreateUser(ctx, op.userID, op.email)
			if err != nil {
				fmt.Printf("Result: FAILED - %v\n", err)
			} else {
				fmt.Printf("Result: SUCCESS\n")
			}
		} else {
			// 获取用户操作 (Get user operation)
			userData, err := userService.GetUser(ctx, op.userID)
			if err != nil {
				fmt.Printf("Result: FAILED - %v\n", err)
			} else {
				fmt.Printf("Result: SUCCESS - %v\n", userData)
			}
		}
		fmt.Println()
	}
}

// demonstrateErrorHandlingIntegration 演示错误处理集成
// (demonstrateErrorHandlingIntegration demonstrates error handling integration)
func demonstrateErrorHandlingIntegration() {
	fmt.Println("=== Demonstrating Error Handling Integration ===")
	fmt.Println()
	
	logger := log.Std().WithValues("integration", "error-handling")
	ehi := NewErrorHandlingIntegration(logger)
	
	// 测试不同的错误场景 (Test different error scenarios)
	testCases := []struct {
		name string
		data map[string]interface{}
	}{
		{
			name: "Successful processing",
			data: map[string]interface{}{
				"id":   "success_case",
				"name": "Test Data",
				"value": 123,
			},
		},
		{
			name: "Empty data error",
			data: map[string]interface{}{},
		},
		{
			name: "Missing ID error",
			data: map[string]interface{}{
				"name": "Test Data",
				"value": 123,
			},
		},
		{
			name: "Processing error",
			data: map[string]interface{}{
				"id":   "error_case",
				"name": "Error Test",
			},
		},
		{
			name: "Persistence error",
			data: map[string]interface{}{
				"id":   "persist_error",
				"name": "Persist Test",
			},
		},
	}
	
	for _, tc := range testCases {
		fmt.Printf("Test case: %s\n", tc.name)
		
		ctx := context.WithValue(context.Background(), "request_id", fmt.Sprintf("req_%d", time.Now().UnixNano()))
		
		err := ehi.ProcessWithErrorHandling(ctx, tc.data)
		if err != nil {
			fmt.Printf("Result: FAILED - %v\n", err)
		} else {
			fmt.Printf("Result: SUCCESS\n")
		}
		fmt.Println()
	}
}

func main() {
	fmt.Println("=== Integration Patterns Example ===")
	fmt.Println("This example demonstrates various logging integration patterns.")
	fmt.Println()
	
	// 1. 演示中间件集成 (Demonstrate middleware integration)
	demonstrateMiddlewareIntegration()
	
	// 2. 演示配置集成 (Demonstrate configuration integration)
	demonstrateConfigurationIntegration()
	
	// 3. 演示服务集成 (Demonstrate service integration)
	demonstrateServiceIntegration()
	
	// 4. 演示错误处理集成 (Demonstrate error handling integration)
	demonstrateErrorHandlingIntegration()
	
	// 最终日志 (Final log)
	opts := log.NewOptions()
	opts.Level = "info"
	opts.Format = "text"
	opts.EnableColor = true
	opts.OutputPaths = []string{"stdout"}
	log.Init(opts)
	
	log.Std().Info("Integration patterns example completed successfully")
	fmt.Println("=== Example completed successfully ===")
} 