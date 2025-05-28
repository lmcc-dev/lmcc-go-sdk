/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * Web application example demonstrating full integration of config, logging, and error handling.
 */

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

// AppConfig Web应用配置
// (AppConfig represents web application configuration)
type AppConfig struct {
	Server struct {
		Port         int    `yaml:"port" default:"8080"`
		Host         string `yaml:"host" default:"0.0.0.0"`
		ReadTimeout  int    `yaml:"read_timeout" default:"30"`
		WriteTimeout int    `yaml:"write_timeout" default:"30"`
		IdleTimeout  int    `yaml:"idle_timeout" default:"120"`
	} `yaml:"server"`

	Database struct {
		Host         string `yaml:"host" default:"localhost"`
		Port         int    `yaml:"port" default:"5432"`
		Name         string `yaml:"name" default:"webapp"`
		User         string `yaml:"user" default:"webapp_user"`
		Password     string `yaml:"password" default:"webapp_pass"`
		MaxConns     int    `yaml:"max_conns" default:"10"`
		ConnTimeout  int    `yaml:"conn_timeout" default:"5"`
	} `yaml:"database"`

	Logging struct {
		Level           string   `yaml:"level" default:"info"`
		Format          string   `yaml:"format" default:"json"`
		OutputPaths     []string `yaml:"output_paths"`
		EnableCaller    bool     `yaml:"enable_caller" default:"true"`
		EnableStacktrace bool    `yaml:"enable_stacktrace" default:"false"`
	} `yaml:"logging"`

	API struct {
		RateLimit    int    `yaml:"rate_limit" default:"100"`
		AuthRequired bool   `yaml:"auth_required" default:"true"`
		ApiVersion   string `yaml:"api_version" default:"v1"`
	} `yaml:"api"`
}

// User 用户数据模型
// (User represents user data model)
type User struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Created  time.Time `json:"created"`
	LastSeen time.Time `json:"last_seen"`
}

// APIResponse 标准API响应格式
// (APIResponse represents standard API response format)
type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	RequestID string      `json:"request_id"`
	Timestamp time.Time   `json:"timestamp"`
}

// WebApp Web应用结构
// (WebApp represents the web application structure)
type WebApp struct {
	config *AppConfig
	logger log.Logger
	server *http.Server
}

// NewWebApp 创建Web应用实例
// (NewWebApp creates a new web application instance)
func NewWebApp(cfg *AppConfig) *WebApp {
	// 设置日志输出路径默认值 (Set default log output paths)
	if len(cfg.Logging.OutputPaths) == 0 {
		cfg.Logging.OutputPaths = []string{"stdout"}
	}

	// 基于配置初始化日志 (Initialize logging based on configuration)
	opts := log.NewOptions()
	opts.Level = cfg.Logging.Level
	opts.Format = cfg.Logging.Format
	opts.EnableColor = cfg.Logging.Format == "text"
	opts.DisableCaller = !cfg.Logging.EnableCaller
	opts.DisableStacktrace = !cfg.Logging.EnableStacktrace
	opts.OutputPaths = cfg.Logging.OutputPaths

	log.Init(opts)

	logger := log.Std().WithValues(
		"component", "web-app",
		"version", "v1.0.0",
		"environment", "development")

	logger.Infow("Web application initialized",
		"server_port", cfg.Server.Port,
		"api_version", cfg.API.ApiVersion,
		"log_level", cfg.Logging.Level)

	return &WebApp{
		config: cfg,
		logger: logger,
	}
}

// LoggingMiddleware 日志中间件
// (LoggingMiddleware provides request logging)
type LoggingMiddleware struct {
	logger log.Logger
}

// NewLoggingMiddleware 创建日志中间件
// (NewLoggingMiddleware creates logging middleware)
func NewLoggingMiddleware(logger log.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger.WithValues("middleware", "logging"),
	}
}

// Handler 日志中间件处理器
// (Handler logging middleware handler)
func (m *LoggingMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 生成请求ID (Generate request ID)
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

		// 将请求信息添加到上下文 (Add request info to context)
		ctx := context.WithValue(r.Context(), "request_id", requestID)
		ctx = context.WithValue(ctx, "logger", requestLogger)

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

// ErrorMiddleware 错误处理中间件
// (ErrorMiddleware provides error handling)
type ErrorMiddleware struct {
	logger log.Logger
}

// NewErrorMiddleware 创建错误处理中间件
// (NewErrorMiddleware creates error handling middleware)
func NewErrorMiddleware(logger log.Logger) *ErrorMiddleware {
	return &ErrorMiddleware{
		logger: logger.WithValues("middleware", "error"),
	}
}

// Handler 错误处理中间件处理器
// (Handler error handling middleware handler)
func (m *ErrorMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				requestID := r.Context().Value("request_id")
				
				m.logger.Errorw("Panic recovered",
					"request_id", requestID,
					"panic", err,
					"path", r.URL.Path,
					"method", r.Method)

				// 返回500错误 (Return 500 error)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// UserService 用户服务
// (UserService provides user operations)
type UserService struct {
	logger log.Logger
	config *AppConfig
}

// NewUserService 创建用户服务
// (NewUserService creates user service)
func NewUserService(logger log.Logger, cfg *AppConfig) *UserService {
	return &UserService{
		logger: logger.WithValues("service", "user"),
		config: cfg,
	}
}

// GetUser 获取用户
// (GetUser retrieves a user)
func (s *UserService) GetUser(ctx context.Context, userID string) (*User, error) {
	requestID := ctx.Value("request_id")
	logger := s.logger.WithValues("request_id", requestID, "operation", "get_user")

	logger.Debugw("Starting user retrieval",
		"user_id", userID)

	if userID == "" {
		err := errors.New("user ID cannot be empty")
		logger.Errorw("User retrieval validation failed", "error", err)
		return nil, err
	}

	// 模拟数据库查询 (Simulate database query)
	if err := s.simulateDBQuery(ctx, "SELECT", userID); err != nil {
		wrappedErr := errors.Wrap(err, "failed to retrieve user from database")
		logger.Errorw("Database query failed", "error", wrappedErr)
		return nil, wrappedErr
	}

	// 模拟用户数据 (Simulate user data)
	user := &User{
		ID:       userID,
		Username: fmt.Sprintf("user_%s", userID),
		Email:    fmt.Sprintf("%s@example.com", userID),
		Created:  time.Now().Add(-30 * 24 * time.Hour),
		LastSeen: time.Now().Add(-2 * time.Hour),
	}

	logger.Infow("User retrieved successfully",
		"user_id", userID,
		"username", user.Username)

	return user, nil
}

// CreateUser 创建用户
// (CreateUser creates a new user)
func (s *UserService) CreateUser(ctx context.Context, username, email string) (*User, error) {
	requestID := ctx.Value("request_id")
	logger := s.logger.WithValues("request_id", requestID, "operation", "create_user")

	logger.Infow("Starting user creation",
		"username", username,
		"email", email)

	// 验证输入 (Validate input)
	if username == "" {
		err := errors.New("username cannot be empty")
		logger.Errorw("User creation validation failed", "error", err)
		return nil, err
	}

	if email == "" {
		err := errors.New("email cannot be empty")
		logger.Errorw("User creation validation failed", "error", err)
		return nil, err
	}

	// 检查用户名是否已存在 (Check if username already exists)
	if username == "admin" || username == "root" {
		err := errors.New("username is reserved")
		logger.Warnw("User creation failed", "error", err, "reason", "reserved_username")
		return nil, err
	}

	// 模拟数据库插入 (Simulate database insert)
	userID := fmt.Sprintf("user_%d", time.Now().Unix())
	if err := s.simulateDBInsert(ctx, userID, username, email); err != nil {
		wrappedErr := errors.Wrap(err, "failed to create user in database")
		logger.Errorw("Database insert failed", "error", wrappedErr)
		return nil, wrappedErr
	}

	user := &User{
		ID:       userID,
		Username: username,
		Email:    email,
		Created:  time.Now(),
		LastSeen: time.Now(),
	}

	logger.Infow("User created successfully",
		"user_id", userID,
		"username", username,
		"email", email)

	return user, nil
}

// simulateDBQuery 模拟数据库查询
// (simulateDBQuery simulates database query)
func (s *UserService) simulateDBQuery(ctx context.Context, operation, userID string) error {
	requestID := ctx.Value("request_id")
	logger := s.logger.WithValues("request_id", requestID, "component", "database")

	start := time.Now()

	logger.Debugw("Database query started",
		"operation", operation,
		"user_id", userID)

	// 模拟数据库延迟 (Simulate database latency)
	time.Sleep(50 * time.Millisecond)

	// 模拟一些错误情况 (Simulate some error conditions)
	if userID == "error" {
		err := errors.New("database connection timeout")
		logger.Errorw("Database query failed", "error", err)
		return err
	}

	if userID == "notfound" {
		err := errors.New("user not found")
		logger.Warnw("User not found", "error", err)
		return err
	}

	duration := time.Since(start)
	logger.Debugw("Database query completed",
		"operation", operation,
		"duration", duration)

	return nil
}

// simulateDBInsert 模拟数据库插入
// (simulateDBInsert simulates database insert)
func (s *UserService) simulateDBInsert(ctx context.Context, userID, username, email string) error {
	requestID := ctx.Value("request_id")
	logger := s.logger.WithValues("request_id", requestID, "component", "database")

	start := time.Now()

	logger.Debugw("Database insert started",
		"user_id", userID,
		"username", username)

	// 模拟数据库延迟 (Simulate database latency)
	time.Sleep(100 * time.Millisecond)

	// 模拟重复邮箱错误 (Simulate duplicate email error)
	if email == "duplicate@example.com" {
		err := errors.New("email address already exists")
		logger.Errorw("Database constraint violation", "error", err)
		return err
	}

	duration := time.Since(start)
	logger.Debugw("Database insert completed",
		"user_id", userID,
		"duration", duration)

	return nil
}

// APIHandler API处理器
// (APIHandler handles API requests)
type APIHandler struct {
	userService *UserService
	logger      log.Logger
}

// NewAPIHandler 创建API处理器
// (NewAPIHandler creates API handler)
func NewAPIHandler(userService *UserService, logger log.Logger) *APIHandler {
	return &APIHandler{
		userService: userService,
		logger:      logger.WithValues("component", "api"),
	}
}

// writeJSONResponse 写入JSON响应
// (writeJSONResponse writes JSON response)
func (h *APIHandler) writeJSONResponse(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}, err error) {
	requestID := r.Context().Value("request_id").(string)

	response := APIResponse{
		Success:   err == nil,
		Data:      data,
		RequestID: requestID,
		Timestamp: time.Now(),
	}

	if err != nil {
		response.Error = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		h.logger.Errorw("Failed to encode JSON response",
			"request_id", requestID,
			"error", encodeErr)
	}
}

// GetUserHandler 获取用户处理器
// (GetUserHandler handles get user requests)
func (h *APIHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeJSONResponse(w, r, http.StatusMethodNotAllowed, nil, 
			errors.New("method not allowed"))
		return
	}

	// 从URL路径中提取用户ID (Extract user ID from URL path)
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		h.writeJSONResponse(w, r, http.StatusBadRequest, nil, 
			errors.New("user ID is required"))
		return
	}

	userID := pathParts[2] // /api/users/{userID}

	user, err := h.userService.GetUser(r.Context(), userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			h.writeJSONResponse(w, r, http.StatusNotFound, nil, err)
		} else {
			h.writeJSONResponse(w, r, http.StatusInternalServerError, nil, err)
		}
		return
	}

	h.writeJSONResponse(w, r, http.StatusOK, user, nil)
}

// CreateUserHandler 创建用户处理器
// (CreateUserHandler handles create user requests)
func (h *APIHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSONResponse(w, r, http.StatusMethodNotAllowed, nil, 
			errors.New("method not allowed"))
		return
	}

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSONResponse(w, r, http.StatusBadRequest, nil, 
			errors.Wrap(err, "invalid JSON request"))
		return
	}

	user, err := h.userService.CreateUser(r.Context(), req.Username, req.Email)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") || 
		   strings.Contains(err.Error(), "reserved") {
			h.writeJSONResponse(w, r, http.StatusConflict, nil, err)
		} else {
			h.writeJSONResponse(w, r, http.StatusInternalServerError, nil, err)
		}
		return
	}

	h.writeJSONResponse(w, r, http.StatusCreated, user, nil)
}

// HealthHandler 健康检查处理器
// (HealthHandler handles health check requests)
func (h *APIHandler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeJSONResponse(w, r, http.StatusMethodNotAllowed, nil, 
			errors.New("method not allowed"))
		return
	}

	health := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"uptime":    time.Since(time.Now().Add(-2 * time.Hour)).String(),
		"version":   "v1.0.0",
	}

	h.writeJSONResponse(w, r, http.StatusOK, health, nil)
}

// setupRoutes 设置路由
// (setupRoutes sets up HTTP routes)
func (app *WebApp) setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// 创建服务 (Create services)
	userService := NewUserService(app.logger, app.config)
	apiHandler := NewAPIHandler(userService, app.logger)

	// 创建中间件 (Create middleware)
	loggingMiddleware := NewLoggingMiddleware(app.logger)
	errorMiddleware := NewErrorMiddleware(app.logger)

	// 设置路由 (Setup routes)
	mux.HandleFunc("/api/health", apiHandler.HealthHandler)
	mux.HandleFunc("/api/users/", apiHandler.GetUserHandler)   // GET /api/users/{id}
	mux.HandleFunc("/api/users", apiHandler.CreateUserHandler) // POST /api/users

	// 应用中间件 (Apply middleware)
	var handler http.Handler = mux
	handler = loggingMiddleware.Handler(handler)
	handler = errorMiddleware.Handler(handler)

	return mux
}

// Start 启动Web应用
// (Start starts the web application)
func (app *WebApp) Start() error {
	mux := app.setupRoutes()

	app.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", app.config.Server.Host, app.config.Server.Port),
		Handler:      mux,
		ReadTimeout:  time.Duration(app.config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(app.config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(app.config.Server.IdleTimeout) * time.Second,
	}

	// 应用中间件到服务器 (Apply middleware to server)
	loggingMiddleware := NewLoggingMiddleware(app.logger)
	errorMiddleware := NewErrorMiddleware(app.logger)

	var handler http.Handler = mux
	handler = loggingMiddleware.Handler(handler)
	handler = errorMiddleware.Handler(handler)

	app.server.Handler = handler

	app.logger.Infow("Starting web server",
		"address", app.server.Addr,
		"read_timeout", app.config.Server.ReadTimeout,
		"write_timeout", app.config.Server.WriteTimeout)

	if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return errors.Wrap(err, "failed to start web server")
	}

	return nil
}

// Stop 停止Web应用
// (Stop stops the web application)
func (app *WebApp) Stop(ctx context.Context) error {
	app.logger.Infow("Stopping web server")

	if app.server != nil {
		if err := app.server.Shutdown(ctx); err != nil {
			return errors.Wrap(err, "failed to stop web server")
		}
	}

	app.logger.Infow("Web server stopped successfully")
	return nil
}

func main() {
	fmt.Println("=== Web Application Integration Example ===")
	fmt.Println("This example demonstrates a complete web application with integrated logging, error handling, and configuration management.")
	fmt.Println()

	// 加载配置 (Load configuration)
	cfg := &AppConfig{}
	if err := config.LoadConfig(cfg); err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// 创建Web应用 (Create web application)
	app := NewWebApp(cfg)

	// 启动服务器在后台 (Start server in background)
	serverErrChan := make(chan error, 1)
	go func() {
		if err := app.Start(); err != nil && err != http.ErrServerClosed {
			serverErrChan <- err
		}
	}()

	// 等待服务器启动 (Wait for server to start)
	time.Sleep(2 * time.Second)

	app.logger.Infow("Web application started successfully",
		"pid", os.Getpid(),
		"port", cfg.Server.Port)

	fmt.Printf("Web server is running on http://%s:%d\n", cfg.Server.Host, cfg.Server.Port)
	fmt.Println("Available endpoints:")
	fmt.Println("  GET  /api/health")
	fmt.Println("  GET  /api/users/{id}")
	fmt.Println("  POST /api/users")
	fmt.Println()

	// 运行自动化测试 (Run automated tests)
	runAutomatedTests(cfg, app.logger)

	// 检查是否有服务器错误 (Check for server errors)
	select {
	case err := <-serverErrChan:
		app.logger.Errorw("Server error occurred", "error", err)
		os.Exit(1)
	default:
		// 继续正常关闭流程 (Continue with normal shutdown)
	}

	// 优雅关闭 (Graceful shutdown)
	app.logger.Infow("Shutting down web application")
	fmt.Println("Shutting down web application...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Stop(ctx); err != nil {
		app.logger.Errorw("Error during shutdown", "error", err)
		os.Exit(1)
	}

	app.logger.Infow("Web application stopped successfully")
	fmt.Println("=== Example completed successfully ===")
}

// runAutomatedTests 运行自动化测试
// (runAutomatedTests runs automated tests against the web application)
func runAutomatedTests(cfg *AppConfig, logger log.Logger) {
	fmt.Println("=== Running Automated Tests ===")
	fmt.Println()

	baseURL := fmt.Sprintf("http://localhost:%d", cfg.Server.Port)
	testLogger := logger.WithValues("component", "automated-tests")

	// 测试1: 健康检查 (Test 1: Health Check)
	fmt.Println("1. Testing Health Check Endpoint:")
	if err := testHealthEndpoint(baseURL, testLogger); err != nil {
		testLogger.Errorw("Health check test failed", "error", err)
		fmt.Printf("   ❌ FAILED: %v\n", err)
	} else {
		fmt.Printf("   ✅ PASSED: Health check endpoint working\n")
	}
	fmt.Println()

	// 测试2: 获取用户（成功案例）(Test 2: Get User - Success Case)
	fmt.Println("2. Testing Get User Endpoint (Success):")
	if err := testGetUserSuccess(baseURL, testLogger); err != nil {
		testLogger.Errorw("Get user success test failed", "error", err)
		fmt.Printf("   ❌ FAILED: %v\n", err)
	} else {
		fmt.Printf("   ✅ PASSED: Get user endpoint working\n")
	}
	fmt.Println()

	// 测试3: 获取用户（错误案例）(Test 3: Get User - Error Case)
	fmt.Println("3. Testing Get User Endpoint (Error Handling):")
	if err := testGetUserError(baseURL, testLogger); err != nil {
		testLogger.Errorw("Get user error test failed", "error", err)
		fmt.Printf("   ❌ FAILED: %v\n", err)
	} else {
		fmt.Printf("   ✅ PASSED: Error handling working correctly\n")
	}
	fmt.Println()

	// 测试4: 创建用户（成功案例）(Test 4: Create User - Success Case)
	fmt.Println("4. Testing Create User Endpoint (Success):")
	if err := testCreateUserSuccess(baseURL, testLogger); err != nil {
		testLogger.Errorw("Create user success test failed", "error", err)
		fmt.Printf("   ❌ FAILED: %v\n", err)
	} else {
		fmt.Printf("   ✅ PASSED: Create user endpoint working\n")
	}
	fmt.Println()

	// 测试5: 创建用户（验证错误）(Test 5: Create User - Validation Error)
	fmt.Println("5. Testing Create User Endpoint (Validation Error):")
	if err := testCreateUserValidation(baseURL, testLogger); err != nil {
		testLogger.Errorw("Create user validation test failed", "error", err)
		fmt.Printf("   ❌ FAILED: %v\n", err)
	} else {
		fmt.Printf("   ✅ PASSED: Validation error handling working\n")
	}
	fmt.Println()

	// 测试6: 不存在的端点 (Test 6: Non-existent Endpoint)
	fmt.Println("6. Testing Non-existent Endpoint:")
	if err := testNotFoundEndpoint(baseURL, testLogger); err != nil {
		testLogger.Errorw("Not found test failed", "error", err)
		fmt.Printf("   ❌ FAILED: %v\n", err)
	} else {
		fmt.Printf("   ✅ PASSED: 404 handling working correctly\n")
	}
	fmt.Println()

	fmt.Println("=== Automated Tests Completed ===")
	fmt.Println()
}

// testHealthEndpoint 测试健康检查端点
// (testHealthEndpoint tests the health check endpoint)
func testHealthEndpoint(baseURL string, logger log.Logger) error {
	url := baseURL + "/api/health"
	logger.Debugw("Testing health endpoint", "url", url)

	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "failed to make health check request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("expected status 200, got %d", resp.StatusCode))
	}

	var response APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return errors.Wrap(err, "failed to decode response")
	}

	if !response.Success {
		return errors.New("response indicates failure")
	}

	logger.Infow("Health check test completed successfully",
		"status_code", resp.StatusCode,
		"success", response.Success)

	return nil
}

// testGetUserSuccess 测试获取用户成功案例
// (testGetUserSuccess tests successful user retrieval)
func testGetUserSuccess(baseURL string, logger log.Logger) error {
	url := baseURL + "/api/users/123"
	logger.Debugw("Testing get user success", "url", url)

	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "failed to make get user request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("expected status 200, got %d", resp.StatusCode))
	}

	var response APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return errors.Wrap(err, "failed to decode response")
	}

	if !response.Success {
		return errors.New("response indicates failure")
	}

	// 验证用户数据 (Validate user data)
	userData, ok := response.Data.(map[string]interface{})
	if !ok {
		return errors.New("invalid user data format")
	}

	if userData["id"] != "123" {
		return errors.New("incorrect user ID returned")
	}

	logger.Infow("Get user success test completed",
		"user_id", userData["id"],
		"username", userData["username"])

	return nil
}

// testGetUserError 测试获取用户错误案例
// (testGetUserError tests user retrieval error handling)
func testGetUserError(baseURL string, logger log.Logger) error {
	url := baseURL + "/api/users/notfound"
	logger.Debugw("Testing get user error", "url", url)

	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "failed to make get user request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		return errors.New(fmt.Sprintf("expected status 404, got %d", resp.StatusCode))
	}

	var response APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return errors.Wrap(err, "failed to decode response")
	}

	if response.Success {
		return errors.New("response should indicate failure")
	}

	if response.Error == "" {
		return errors.New("error message should be present")
	}

	logger.Infow("Get user error test completed",
		"status_code", resp.StatusCode,
		"error", response.Error)

	return nil
}

// testCreateUserSuccess 测试创建用户成功案例
// (testCreateUserSuccess tests successful user creation)
func testCreateUserSuccess(baseURL string, logger log.Logger) error {
	url := baseURL + "/api/users"
	logger.Debugw("Testing create user success", "url", url)

	requestBody := map[string]string{
		"username": "john_doe",
		"email":    "john@example.com",
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return errors.Wrap(err, "failed to marshal request body")
	}

	resp, err := http.Post(url, "application/json", strings.NewReader(string(jsonBody)))
	if err != nil {
		return errors.Wrap(err, "failed to make create user request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return errors.New(fmt.Sprintf("expected status 201, got %d", resp.StatusCode))
	}

	var response APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return errors.Wrap(err, "failed to decode response")
	}

	if !response.Success {
		return errors.New("response indicates failure")
	}

	// 验证用户数据 (Validate user data)
	userData, ok := response.Data.(map[string]interface{})
	if !ok {
		return errors.New("invalid user data format")
	}

	if userData["username"] != "john_doe" {
		return errors.New("incorrect username returned")
	}

	logger.Infow("Create user success test completed",
		"username", userData["username"],
		"email", userData["email"])

	return nil
}

// testCreateUserValidation 测试创建用户验证错误
// (testCreateUserValidation tests user creation validation errors)
func testCreateUserValidation(baseURL string, logger log.Logger) error {
	url := baseURL + "/api/users"
	logger.Debugw("Testing create user validation", "url", url)

	requestBody := map[string]string{
		"username": "admin", // 保留用户名 (Reserved username)
		"email":    "admin@example.com",
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return errors.Wrap(err, "failed to marshal request body")
	}

	resp, err := http.Post(url, "application/json", strings.NewReader(string(jsonBody)))
	if err != nil {
		return errors.Wrap(err, "failed to make create user request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusConflict {
		return errors.New(fmt.Sprintf("expected status 409, got %d", resp.StatusCode))
	}

	var response APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return errors.Wrap(err, "failed to decode response")
	}

	if response.Success {
		return errors.New("response should indicate failure")
	}

	if !strings.Contains(response.Error, "reserved") {
		return errors.New("error message should mention reserved username")
	}

	logger.Infow("Create user validation test completed",
		"status_code", resp.StatusCode,
		"error", response.Error)

	return nil
}

// testNotFoundEndpoint 测试不存在的端点
// (testNotFoundEndpoint tests non-existent endpoint handling)
func testNotFoundEndpoint(baseURL string, logger log.Logger) error {
	url := baseURL + "/api/nonexistent"
	logger.Debugw("Testing not found endpoint", "url", url)

	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "failed to make request to non-existent endpoint")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		return errors.New(fmt.Sprintf("expected status 404, got %d", resp.StatusCode))
	}

	logger.Infow("Not found endpoint test completed",
		"status_code", resp.StatusCode)

	return nil
} 