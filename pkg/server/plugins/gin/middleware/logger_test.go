/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Gin日志中间件单元测试 (Gin logger middleware unit tests)
 */

package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	ginpkg "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
)

// TestDefaultLoggerConfig 测试默认日志配置 (Test default logger configuration)
func TestDefaultLoggerConfig(t *testing.T) {
	config := DefaultLoggerConfig()
	
	assert.NotNil(t, config)
	assert.Equal(t, "json", config.Format)
	assert.Equal(t, false, config.IncludeBody)
	assert.Equal(t, 1024, config.MaxBodySize)
	assert.Contains(t, config.SkipPaths, "/health")
	assert.Contains(t, config.SkipPaths, "/metrics")
}

// TestNewGinLoggerMiddleware 测试创建日志中间件 (Test creating logger middleware)
func TestNewGinLoggerMiddleware(t *testing.T) {
	// 测试使用默认配置 (Test with default config)
	middleware := NewGinLoggerMiddleware(nil)
	assert.NotNil(t, middleware)
	
	// 测试使用自定义配置 (Test with custom config)
	config := &LoggerConfig{
		SkipPaths:   []string{"/test"},
		Format:      "text",
		IncludeBody: true,
		MaxBodySize: 2048,
	}
	middleware2 := NewGinLoggerMiddleware(config)
	assert.NotNil(t, middleware2)
}

// TestNewGinLoggerMiddlewareWithServices 测试创建带服务的日志中间件 (Test creating logger middleware with services)
func TestNewGinLoggerMiddlewareWithServices(t *testing.T) {
	config := DefaultLoggerConfig()
	serviceContainer := services.NewServiceContainerWithDefaults()
	
	middleware := NewGinLoggerMiddlewareWithServices(config, serviceContainer)
	assert.NotNil(t, middleware)
	assert.Equal(t, config, middleware.config)
	assert.Equal(t, serviceContainer, middleware.services)
}

// TestGinLoggerMiddleware_Process 测试Process方法 (Test Process method)
func TestGinLoggerMiddleware_Process(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// 创建测试服务容器 (Create test service container)
	serviceContainer := services.NewServiceContainerWithDefaults()
	
	config := &LoggerConfig{
		SkipPaths:   []string{"/skip"},
		Format:      "json",
		IncludeBody: false,
		MaxBodySize: 1024,
	}
	
	middleware := NewGinLoggerMiddlewareWithServices(config, serviceContainer)
	
	// 创建测试引擎和路由 (Create test engine and routes)
	engine := gin.New()
	
	var loggedMessage string
	var loggedFields []interface{}
	
	// 使用自定义日志器来捕获日志输出 (Use custom logger to capture log output)
	mockLogger := &mockLogger{
		onDebugw: func(msg string, keysAndValues ...interface{}) {
			loggedMessage = msg
			loggedFields = keysAndValues
		},
	}
	serviceContainer.SetLogger(mockLogger)
	
	// 添加中间件和路由 (Add middleware and route)
	engine.Use(func(c *gin.Context) {
		ctx := ginpkg.NewGinContext(c)
		err := middleware.Process(ctx, func() error {
			c.Status(200)
			return nil
		})
		assert.NoError(t, err)
	})
	
	engine.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})
	
	// 执行请求 (Execute request)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("User-Agent", "test-agent")
	
	engine.ServeHTTP(w, req)
	
	// 验证响应 (Verify response)
	assert.Equal(t, http.StatusOK, w.Code)
	
	// 验证日志记录 (Verify logging)
	assert.Contains(t, loggedMessage, "GET")
	assert.Contains(t, loggedMessage, "/test")
	assert.Contains(t, loggedMessage, "200")
	
	// 验证日志字段 (Verify log fields)
	assert.Greater(t, len(loggedFields), 0)
}

// TestGinLoggerMiddleware_SkipPaths 测试跳过路径功能 (Test skip paths functionality)
func TestGinLoggerMiddleware_SkipPaths(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	
	config := &LoggerConfig{
		SkipPaths:   []string{"/health", "/metrics"},
		Format:      "json",
		IncludeBody: false,
		MaxBodySize: 1024,
	}
	
	middleware := NewGinLoggerMiddlewareWithServices(config, serviceContainer)
	
	engine := gin.New()
	
	logCalled := false
	mockLogger := &mockLogger{
		onDebugw: func(msg string, keysAndValues ...interface{}) {
			logCalled = true
		},
	}
	serviceContainer.SetLogger(mockLogger)
	
	// 添加中间件 (Add middleware)
	engine.Use(func(c *gin.Context) {
		ctx := ginpkg.NewGinContext(c)
		err := middleware.Process(ctx, func() error {
			c.Status(200)
			return nil
		})
		assert.NoError(t, err)
	})
	
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	
	// 请求跳过的路径 (Request skipped path)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)
	
	engine.ServeHTTP(w, req)
	
	// 验证响应成功但没有记录日志 (Verify response success but no logging)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.False(t, logCalled, "Logger should not be called for skipped paths")
}

// TestGinLoggerMiddleware_WithRequestBody 测试记录请求体 (Test logging request body)
func TestGinLoggerMiddleware_WithRequestBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	
	config := &LoggerConfig{
		SkipPaths:   []string{},
		Format:      "json",
		IncludeBody: true,
		MaxBodySize: 1024,
	}
	
	middleware := NewGinLoggerMiddlewareWithServices(config, serviceContainer)
	
	engine := gin.New()
	
	var loggedFields []interface{}
	mockLogger := &mockLogger{
		onDebugw: func(msg string, keysAndValues ...interface{}) {
			loggedFields = keysAndValues
		},
	}
	serviceContainer.SetLogger(mockLogger)
	
	// 添加中间件 (Add middleware)
	engine.Use(func(c *gin.Context) {
		ctx := ginpkg.NewGinContext(c)
		err := middleware.Process(ctx, func() error {
			c.Status(200)
			return nil
		})
		assert.NoError(t, err)
	})
	
	engine.POST("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "received"})
	})
	
	// 创建带请求体的POST请求 (Create POST request with body)
	requestBody := `{"name": "test", "value": 123}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/test", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	
	engine.ServeHTTP(w, req)
	
	// 验证响应 (Verify response)
	assert.Equal(t, http.StatusOK, w.Code)
	
	// 验证请求体被记录 (Verify request body is logged)
	found := false
	for i := 0; i < len(loggedFields)-1; i += 2 {
		if key, ok := loggedFields[i].(string); ok && key == "request_body" {
			if value, ok := loggedFields[i+1].(string); ok {
				assert.Contains(t, value, "test")
				assert.Contains(t, value, "123")
				found = true
				break
			}
		}
	}
	assert.True(t, found, "Request body should be logged")
}

// TestGinLoggerMiddleware_ErrorHandling 测试错误处理 (Test error handling)
func TestGinLoggerMiddleware_ErrorHandling(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	config := DefaultLoggerConfig()
	middleware := NewGinLoggerMiddlewareWithServices(config, serviceContainer)
	
	engine := gin.New()
	
	var logLevel string
	var loggedMessage string
	mockLogger := &mockLogger{
		onErrorw: func(msg string, keysAndValues ...interface{}) {
			logLevel = "error"
			loggedMessage = msg
		},
	}
	serviceContainer.SetLogger(mockLogger)
	
	// 添加中间件 (Add middleware)
	engine.Use(func(c *gin.Context) {
		ctx := ginpkg.NewGinContext(c)
		err := middleware.Process(ctx, func() error {
			c.Status(500)
			return assert.AnError // 返回错误 (Return error)
		})
		assert.Error(t, err)
	})
	
	engine.GET("/error", func(c *gin.Context) {
		c.JSON(500, gin.H{"error": "test error"})
	})
	
	// 执行请求 (Execute request)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/error", nil)
	
	engine.ServeHTTP(w, req)
	
	// 验证使用了错误级别的日志 (Verify error level logging is used)
	assert.Equal(t, "error", logLevel)
	assert.Contains(t, loggedMessage, "500")
}

// TestGinLoggerMiddleware_MaxBodySize 测试最大请求体限制 (Test max body size limit)
func TestGinLoggerMiddleware_MaxBodySize(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	
	config := &LoggerConfig{
		SkipPaths:   []string{},
		Format:      "json",
		IncludeBody: true,
		MaxBodySize: 10, // 很小的限制 (Very small limit)
	}
	
	middleware := NewGinLoggerMiddlewareWithServices(config, serviceContainer)
	
	engine := gin.New()
	
	var loggedFields []interface{}
	mockLogger := &mockLogger{
		onDebugw: func(msg string, keysAndValues ...interface{}) {
			loggedFields = keysAndValues
		},
	}
	serviceContainer.SetLogger(mockLogger)
	
	// 添加中间件 (Add middleware)
	engine.Use(func(c *gin.Context) {
		ctx := ginpkg.NewGinContext(c)
		err := middleware.Process(ctx, func() error {
			c.Status(200)
			return nil
		})
		assert.NoError(t, err)
	})
	
	engine.POST("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "received"})
	})
	
	// 创建大于限制的请求体 (Create request body larger than limit)
	largeRequestBody := "this is a very long request body that exceeds the limit"
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/test", strings.NewReader(largeRequestBody))
	req.Header.Set("Content-Type", "application/json")
	
	engine.ServeHTTP(w, req)
	
	// 验证响应 (Verify response)
	assert.Equal(t, http.StatusOK, w.Code)
	
	// 验证请求体被截断 (Verify request body is truncated)
	found := false
	for i := 0; i < len(loggedFields)-1; i += 2 {
		if key, ok := loggedFields[i].(string); ok && key == "request_body" {
			if value, ok := loggedFields[i+1].(string); ok {
				assert.LessOrEqual(t, len(value), 10, "Request body should be truncated")
				found = true
				break
			}
		}
	}
	assert.True(t, found, "Request body should be logged even when truncated")
}

// TestGinLoggerMiddleware_Apply 测试Apply方法 (Test Apply method)
func TestGinLoggerMiddleware_Apply(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	config := DefaultLoggerConfig()
	middleware := NewGinLoggerMiddleware(config).(*GinLoggerMiddleware)
	
	engine := gin.New()
	
	// 调用Apply方法 (Call Apply method)
	assert.NotPanics(t, func() {
		middleware.Apply(engine)
	})
}

// TestGinLoggerMiddleware_CreateGinHandler 测试CreateGinHandler方法 (Test CreateGinHandler method)
func TestGinLoggerMiddleware_CreateGinHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	config := DefaultLoggerConfig()
	middleware := NewGinLoggerMiddleware(config).(*GinLoggerMiddleware)
	
	// 创建Gin处理器 (Create Gin handler)
	handler := middleware.CreateGinHandler()
	assert.NotNil(t, handler)
	
	// 测试处理器不会panic (Test handler doesn't panic)
	engine := gin.New()
	engine.Use(handler)
	engine.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})
	
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	
	assert.NotPanics(t, func() {
		engine.ServeHTTP(w, req)
	})
	
	assert.Equal(t, http.StatusOK, w.Code)
}

// BenchmarkGinLoggerMiddleware_Process 基准测试Process方法 (Benchmark Process method)
func BenchmarkGinLoggerMiddleware_Process(b *testing.B) {
	gin.SetMode(gin.TestMode)
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	config := DefaultLoggerConfig()
	middleware := NewGinLoggerMiddlewareWithServices(config, serviceContainer)
	
	engine := gin.New()
	engine.Use(func(c *gin.Context) {
		ctx := ginpkg.NewGinContext(c)
		middleware.Process(ctx, func() error {
			c.Status(200)
			return nil
		})
	})
	
	engine.GET("/benchmark", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "benchmark"})
	})
	
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/benchmark", nil)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Body.Reset()
		w.Code = 0
		engine.ServeHTTP(w, req)
	}
}

// mockLogger 模拟日志器用于测试 (Mock logger for testing)
type mockLogger struct {
	onDebug  func(msg string)
	onDebugf func(template string, args ...interface{})
	onDebugw func(msg string, keysAndValues ...interface{})
	onInfo   func(msg string)
	onInfof  func(template string, args ...interface{})
	onInfow  func(msg string, keysAndValues ...interface{})
	onWarn   func(msg string)
	onWarnf  func(template string, args ...interface{})
	onWarnw  func(msg string, keysAndValues ...interface{})
	onError  func(msg string)
	onErrorf func(template string, args ...interface{})
	onErrorw func(msg string, keysAndValues ...interface{})
	onFatal  func(msg string)
	onFatalf func(template string, args ...interface{})
	onPanic  func(msg string)
	onPanicf func(template string, args ...interface{})
}

func (m *mockLogger) Debug(msg string) {
	if m.onDebug != nil {
		m.onDebug(msg)
	}
}

func (m *mockLogger) Debugf(template string, args ...interface{}) {
	if m.onDebugf != nil {
		m.onDebugf(template, args...)
	}
}

func (m *mockLogger) Debugw(msg string, keysAndValues ...interface{}) {
	if m.onDebugw != nil {
		m.onDebugw(msg, keysAndValues...)
	}
}

func (m *mockLogger) Info(msg string) {
	if m.onInfo != nil {
		m.onInfo(msg)
	}
}

func (m *mockLogger) Infof(template string, args ...interface{}) {
	if m.onInfof != nil {
		m.onInfof(template, args...)
	}
}

func (m *mockLogger) Infow(msg string, keysAndValues ...interface{}) {
	if m.onInfow != nil {
		m.onInfow(msg, keysAndValues...)
	}
}

func (m *mockLogger) Warn(msg string) {
	if m.onWarn != nil {
		m.onWarn(msg)
	}
}

func (m *mockLogger) Warnf(template string, args ...interface{}) {
	if m.onWarnf != nil {
		m.onWarnf(template, args...)
	}
}

func (m *mockLogger) Warnw(msg string, keysAndValues ...interface{}) {
	if m.onWarnw != nil {
		m.onWarnw(msg, keysAndValues...)
	}
}

func (m *mockLogger) Error(msg string) {
	if m.onError != nil {
		m.onError(msg)
	}
}

func (m *mockLogger) Errorf(template string, args ...interface{}) {
	if m.onErrorf != nil {
		m.onErrorf(template, args...)
	}
}

func (m *mockLogger) Errorw(msg string, keysAndValues ...interface{}) {
	if m.onErrorw != nil {
		m.onErrorw(msg, keysAndValues...)
	}
}

func (m *mockLogger) Fatal(msg string) {
	if m.onFatal != nil {
		m.onFatal(msg)
	}
}

func (m *mockLogger) Fatalf(template string, args ...interface{}) {
	if m.onFatalf != nil {
		m.onFatalf(template, args...)
	}
}

func (m *mockLogger) Panic(msg string) {
	if m.onPanic != nil {
		m.onPanic(msg)
	}
}

func (m *mockLogger) Panicf(template string, args ...interface{}) {
	if m.onPanicf != nil {
		m.onPanicf(template, args...)
	}
} 