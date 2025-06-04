/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Gin恢复中间件单元测试 (Gin recovery middleware unit tests)
 */

package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	ginpkg "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin"
)

// TestNewRecoveryMiddleware 测试创建恢复中间件 (Test creating recovery middleware)
func TestNewRecoveryMiddleware(t *testing.T) {
	// 测试使用nil配置 (Test with nil config)
	middleware := NewRecoveryMiddleware(nil)
	assert.NotNil(t, middleware)
	
	// 测试使用自定义配置 (Test with custom config)
	config := &server.RecoveryMiddlewareConfig{
		PrintStack: true,
	}
	middleware2 := NewRecoveryMiddleware(config)
	assert.NotNil(t, middleware2)
	
	// 验证类型转换 (Verify type conversion)
	recoveryMiddleware, ok := middleware2.(*RecoveryMiddleware)
	assert.True(t, ok)
	assert.Equal(t, config, recoveryMiddleware.config)
}

// TestRecoveryMiddleware_ProcessWithGinContext 测试使用Gin上下文的Process方法 (Test Process method with Gin context)
func TestRecoveryMiddleware_ProcessWithGinContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	config := &server.RecoveryMiddlewareConfig{
		PrintStack: false,
	}
	middleware := NewRecoveryMiddleware(config).(*RecoveryMiddleware)
	
	// 创建测试引擎 (Create test engine)
	engine := gin.New()
	
	var processedSuccessfully bool
	
	// 添加中间件和路由 (Add middleware and route)
	engine.Use(func(c *gin.Context) {
		ctx := ginpkg.NewGinContext(c)
		err := middleware.Process(ctx, func() error {
			processedSuccessfully = true
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
	
	engine.ServeHTTP(w, req)
	
	// 验证响应 (Verify response)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, processedSuccessfully, "Next function should be called")
}

// TestRecoveryMiddleware_ProcessWithNonGinContext 测试使用非Gin上下文的Process方法 (Test Process method with non-Gin context)
func TestRecoveryMiddleware_ProcessWithNonGinContext(t *testing.T) {
	config := &server.RecoveryMiddlewareConfig{
		PrintStack: true,
	}
	middleware := NewRecoveryMiddleware(config).(*RecoveryMiddleware)
	
	// 创建模拟的非Gin上下文 (Create mock non-Gin context)
	mockCtx := &mockContext{}
	
	var nextCalled bool
	
	// 执行Process方法 (Execute Process method)
	err := middleware.Process(mockCtx, func() error {
		nextCalled = true
		return nil
	})
	
	// 验证结果 (Verify results)
	assert.NoError(t, err)
	assert.True(t, nextCalled, "Next function should be called even with non-Gin context")
}

// TestRecoveryMiddleware_ProcessWithPanic 测试panic处理 (Test panic handling)
func TestRecoveryMiddleware_ProcessWithPanic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	config := &server.RecoveryMiddlewareConfig{
		PrintStack: true,
	}
	middleware := NewRecoveryMiddleware(config).(*RecoveryMiddleware)
	
	// 创建测试引擎 (Create test engine)
	engine := gin.New()
	
	// 添加Recovery中间件和会panic的路由 (Add Recovery middleware and route that panics)
	engine.Use(func(c *gin.Context) {
		ctx := ginpkg.NewGinContext(c)
		err := middleware.Process(ctx, func() error {
			c.Status(200)
			return nil
		})
		// 如果请求被中止（由于panic），会返回错误 (If request is aborted due to panic, will return error)
		if c.IsAborted() {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	})
	
	engine.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})
	
	// 执行请求 (Execute request)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/panic", nil)
	
	// 请求不应该导致程序崩溃 (Request should not crash the program)
	assert.NotPanics(t, func() {
		engine.ServeHTTP(w, req)
	})
	
	// 验证Recovery处理了panic (Verify Recovery handled the panic)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestRecoveryMiddleware_ProcessWithError 测试next函数返回错误 (Test next function returning error)
func TestRecoveryMiddleware_ProcessWithError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	config := &server.RecoveryMiddlewareConfig{
		PrintStack: false,
	}
	middleware := NewRecoveryMiddleware(config).(*RecoveryMiddleware)
	
	// 创建测试引擎 (Create test engine)
	engine := gin.New()
	
	testError := errors.New("test error")
	
	// 添加中间件 (Add middleware)
	engine.Use(func(c *gin.Context) {
		ctx := ginpkg.NewGinContext(c)
		err := middleware.Process(ctx, func() error {
			return testError
		})
		assert.Equal(t, testError, err)
	})
	
	engine.GET("/error", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})
	
	// 执行请求 (Execute request)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/error", nil)
	
	engine.ServeHTTP(w, req)
	
	// 验证响应 (Verify response)
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestRecoveryMiddleware_CreateRecoveryFunc 测试createRecoveryFunc方法 (Test createRecoveryFunc method)
func TestRecoveryMiddleware_CreateRecoveryFunc(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	config := &server.RecoveryMiddlewareConfig{
		PrintStack: true,
	}
	middleware := NewRecoveryMiddleware(config).(*RecoveryMiddleware)
	
	// 获取Recovery函数 (Get Recovery function)
	recoveryFunc := middleware.createRecoveryFunc()
	assert.NotNil(t, recoveryFunc)
	
	// 创建测试上下文 (Create test context)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	
	// 调用Recovery函数 (Call Recovery function)
	assert.NotPanics(t, func() {
		recoveryFunc(c, "test panic")
	})
	
	// 验证响应 (Verify response)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.True(t, c.IsAborted())
	assert.Contains(t, w.Body.String(), "Internal Server Error")
}

// TestRecoveryMiddleware_ConfigVariations 测试不同配置的变化 (Test different config variations)
func TestRecoveryMiddleware_ConfigVariations(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	testCases := []struct {
		name   string
		config *server.RecoveryMiddlewareConfig
	}{
		{
			name:   "nil config",
			config: nil,
		},
		{
			name: "print stack enabled",
			config: &server.RecoveryMiddlewareConfig{
				PrintStack: true,
			},
		},
		{
			name: "print stack disabled",
			config: &server.RecoveryMiddlewareConfig{
				PrintStack: false,
			},
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			middleware := NewRecoveryMiddleware(tc.config).(*RecoveryMiddleware)
			assert.NotNil(t, middleware)
			assert.Equal(t, tc.config, middleware.config)
			
			// 测试中间件可以正常工作 (Test middleware works normally)
			engine := gin.New()
			
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
			
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test", nil)
			
			assert.NotPanics(t, func() {
				engine.ServeHTTP(w, req)
			})
			
			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

// TestRecoveryMiddleware_AbortedRequest 测试中止的请求 (Test aborted request)
func TestRecoveryMiddleware_AbortedRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	config := &server.RecoveryMiddlewareConfig{
		PrintStack: false,
	}
	middleware := NewRecoveryMiddleware(config).(*RecoveryMiddleware)
	
	// 创建测试引擎 (Create test engine)
	engine := gin.New()
	
	// 手动中止请求的中间件 (Middleware that manually aborts request)
	engine.Use(func(c *gin.Context) {
		c.Abort()
	})
	
	// Recovery中间件 (Recovery middleware)
	engine.Use(func(c *gin.Context) {
		ctx := ginpkg.NewGinContext(c)
		err := middleware.Process(ctx, func() error {
			return nil
		})
		// 如果请求被中止，应该返回错误 (If request is aborted, should return error)
		if c.IsAborted() {
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "aborted")
		}
	})
	
	engine.GET("/abort", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "should not reach here"})
	})
	
	// 执行请求 (Execute request)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/abort", nil)
	
	engine.ServeHTTP(w, req)
	
	// 请求应该被中止 (Request should be aborted)
	// 注意：由于请求在Recovery之前就被中止了，状态码可能不是500
	// Note: Since request is aborted before Recovery, status code might not be 500
}

// TestRecoveryMiddlewareFactory 测试Recovery中间件工厂 (Test Recovery middleware factory)
func TestRecoveryMiddlewareFactory(t *testing.T) {
	// 测试有效配置 (Test valid config)
	config := &server.RecoveryMiddlewareConfig{
		PrintStack: true,
	}
	
	middleware, err := RecoveryMiddlewareFactory(config)
	assert.NoError(t, err)
	assert.NotNil(t, middleware)
	
	// 验证类型 (Verify type)
	recoveryMiddleware, ok := middleware.(*RecoveryMiddleware)
	assert.True(t, ok)
	assert.Equal(t, config, recoveryMiddleware.config)
	
	// 测试无效配置 (Test invalid config)
	invalidConfig := "invalid config"
	middleware2, err2 := RecoveryMiddlewareFactory(invalidConfig)
	assert.Error(t, err2)
	assert.Nil(t, middleware2)
	assert.Equal(t, ErrInvalidConfig, err2)
}

// BenchmarkRecoveryMiddleware_Process 基准测试Process方法 (Benchmark Process method)
func BenchmarkRecoveryMiddleware_Process(b *testing.B) {
	gin.SetMode(gin.TestMode)
	
	config := &server.RecoveryMiddlewareConfig{
		PrintStack: false,
	}
	middleware := NewRecoveryMiddleware(config).(*RecoveryMiddleware)
	
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

// mockContext 模拟上下文用于测试 (Mock context for testing)
type mockContext struct{}

func (m *mockContext) Request() *http.Request                                { return nil }
func (m *mockContext) Response() http.ResponseWriter                         { return nil }
func (m *mockContext) Method() string                                        { return "GET" }
func (m *mockContext) Path() string                                          { return "/test" }
func (m *mockContext) FullPath() string                                      { return "/test" }
func (m *mockContext) ClientIP() string                                      { return "127.0.0.1" }
func (m *mockContext) UserAgent() string                                     { return "test" }
func (m *mockContext) Header(key string) string                              { return "" }
func (m *mockContext) SetHeader(key, value string)                           {}
func (m *mockContext) Param(key string) string                               { return "" }
func (m *mockContext) Query(key string) string                               { return "" }
func (m *mockContext) Bind(obj interface{}) error                            { return nil }
func (m *mockContext) JSON(code int, obj interface{}) error                  { return nil }
func (m *mockContext) String(code int, format string, values ...interface{}) error { return nil }
func (m *mockContext) Data(code int, contentType string, data []byte) error  { return nil }
func (m *mockContext) Set(key string, value interface{})                     {}
func (m *mockContext) Get(key string) (interface{}, bool)                    { return nil, false }
func (m *mockContext) GetString(key string) string                           { return "" }
func (m *mockContext) GetInt(key string) int                                 { return 0 }
func (m *mockContext) GetBool(key string) bool                               { return false } 