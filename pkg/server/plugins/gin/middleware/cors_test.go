/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Tests for CORS middleware implementation / CORS中间件实现的测试
 */

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	"github.com/stretchr/testify/assert"
)

// setupTestEngine 设置测试用的Gin引擎 (Setup test Gin engine)
func setupTestEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.Default()
}

// TestNewCORSMiddleware 测试创建CORS中间件 (Test creating CORS middleware)
func TestNewCORSMiddleware(t *testing.T) {
	config := &server.CORSConfig{
		Enabled:          true,
		AllowOrigins:     []string{"https://example.com"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           1 * time.Hour,
	}
	
	middleware := NewCORSMiddleware(config)
	assert.NotNil(t, middleware)
}

// TestCORSMiddlewareFactory 测试CORS中间件工厂 (Test CORS middleware factory)
func TestCORSMiddlewareFactory(t *testing.T) {
	config := &server.CORSConfig{
		Enabled:      true,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST"},
	}
	
	middleware, err := CORSMiddlewareFactory(config)
	assert.NoError(t, err)
	assert.NotNil(t, middleware)
	
	// 测试无效配置 (Test invalid config)
	_, err = CORSMiddlewareFactory("invalid config")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidConfig, err)
}

// TestCORSWithAllowAllOrigins 测试允许所有源的配置 (Test allow all origins configuration)
func TestCORSWithAllowAllOrigins(t *testing.T) {
	config := &server.CORSConfig{
		Enabled:          true,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: false, // 关键：不启用凭证 (Key: no credentials)
		MaxAge:           1 * time.Hour,
	}
	
	corsMiddleware := NewCORSMiddleware(config).(*CORSMiddleware)
	engine := setupTestEngine()
	engine.Use(corsMiddleware.GetGinHandler())
	
	engine.GET("/api/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{"message": "test"})
	})
	
	// 测试任意域请求 (Test any domain request)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("Origin", "https://anydomain.com")
	engine.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}

// TestCORSWithSpecificOrigin 测试特定来源的CORS (Test CORS with specific origin)
func TestCORSWithSpecificOrigin(t *testing.T) {
	config := &server.CORSConfig{
		Enabled:      true,
		AllowOrigins: []string{"http://example.com", "https://app.com"}, 
		AllowMethods: []string{"GET", "POST", "PUT", "HEAD"}, 
		AllowHeaders: []string{"Content-Type", "Authorization"}, 
		ExposeHeaders: []string{"X-Total-Count", "X-User-ID"}, 
		AllowCredentials: false,
		MaxAge:       12 * time.Hour, 
	}
	
	corsMiddleware := NewCORSMiddleware(config).(*CORSMiddleware)
	engine := setupTestEngine()
	engine.Use(corsMiddleware.GetGinHandler())
	engine.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "success")
	})
	
	// 测试没有Origin头的请求（非CORS请求）(Test request without Origin header - non-CORS request)
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
	
	// 测试第一个允许的Origin (Test first allowed origin)
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "http://example.com")
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "http://example.com", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "X-Total-Count,X-User-ID", w.Header().Get("Access-Control-Expose-Headers"))

	// 测试第二个允许的Origin (Test second allowed origin)
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "https://app.com")
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "https://app.com", w.Header().Get("Access-Control-Allow-Origin"))

	// 测试不匹配的Origin (Test non-matching origin)
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Origin", "https://evil.com")
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

// TestCORSWithCredentialsAndWildcard 测试凭证+通配符冲突处理 (Test credentials + wildcard conflict handling)
func TestCORSWithCredentialsAndWildcard(t *testing.T) {
	config := &server.CORSConfig{
		Enabled:          true,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true, // 启用凭证与通配符冲突 (Enable credentials conflicts with wildcard)
		MaxAge:           1 * time.Hour,
	}
	
	corsMiddleware := NewCORSMiddleware(config).(*CORSMiddleware)
	engine := setupTestEngine()
	engine.Use(corsMiddleware.GetGinHandler())
	
	engine.GET("/api/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{"message": "test"})
	})
	
	// 测试任意源请求（即使配置了*，由于启用了凭证，应该回显具体origin）
	// (Test any origin request - even with *, since credentials are enabled, should echo specific origin)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	engine.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
}

// TestCORSPreflightRequest 测试预检请求 (Test preflight request)
func TestCORSPreflightRequest(t *testing.T) {
	config := &server.CORSConfig{
		Enabled:      true,
		AllowOrigins: []string{"https://example.com"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		MaxAge:       12 * time.Hour,
	}
	
	corsMiddleware := NewCORSMiddleware(config).(*CORSMiddleware)
	engine := setupTestEngine()
	engine.Use(corsMiddleware.GetGinHandler())
	
	// 自定义CORS实现会自动处理OPTIONS请求 (Custom CORS implementation automatically handles OPTIONS requests)
	engine.POST("/api/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{"message": "post"})
	})
	
	// 测试有效的预检请求 (Test valid preflight request)
	req := httptest.NewRequest("OPTIONS", "/api/test", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "https://example.com", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
	assert.Equal(t, "43200", w.Header().Get("Access-Control-Max-Age")) // 12 * 3600
	
	// 测试不被允许的预检请求 (Test disallowed preflight request)
	req = httptest.NewRequest("OPTIONS", "/api/test", nil)
	req.Header.Set("Origin", "https://evil.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

// TestCORSMiddlewareProcess 测试使用Process方法的CORS中间件 (Test CORS middleware using Process method)
func TestCORSMiddlewareProcess(t *testing.T) {
	config := &server.CORSConfig{
		Enabled:      true,
		AllowOrigins: []string{"https://example.com", "https://app.com"},
		AllowMethods: []string{"GET", "POST", "PUT"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		ExposeHeaders: []string{"X-Total-Count"},
		AllowCredentials: false,
		MaxAge:       2 * time.Hour,
	}
	
	middleware := NewCORSMiddleware(config)
	
	// 创建模拟请求和响应 (Create mock request and response)
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()
	
	// 创建基础上下文 (Create base context)
	ctx := server.NewBaseContext(req, w)
	
	// 执行CORS中间件 (Execute CORS middleware)
	nextCalled := false
	err := middleware.Process(ctx, func() error {
		nextCalled = true
		ctx.SetHeader("X-Total-Count", "100")
		return ctx.JSON(http.StatusOK, map[string]string{"data": "test"})
	})
	
	assert.NoError(t, err)
	assert.True(t, nextCalled)
	assert.Equal(t, "https://example.com", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "X-Total-Count", w.Header().Get("Access-Control-Expose-Headers"))
}

// TestCORSNonGinContext 测试非Gin上下文 (Test non-Gin context)
func TestCORSNonGinContext(t *testing.T) {
	config := &server.CORSConfig{
		Enabled:      true,
		AllowOrigins: []string{"*"},
	}
	
	middleware := NewCORSMiddleware(config)
	
	// 创建模拟请求和响应 (Create mock request and response)
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()
	
	// 创建基础上下文 (Create base context)
	baseCtx := server.NewBaseContext(req, w)
	
	nextCalled := false
	err := middleware.Process(baseCtx, func() error {
		nextCalled = true
		return nil
	})
	
	assert.NoError(t, err)
	assert.True(t, nextCalled)
	assert.Equal(t, "https://example.com", w.Header().Get("Access-Control-Allow-Origin"))
}

// TestCORSDisabled 测试禁用CORS (Test disabled CORS)
func TestCORSDisabled(t *testing.T) {
	config := &server.CORSConfig{
		Enabled: false, // CORS已禁用 (CORS disabled)
	}
	
	corsMiddleware := NewCORSMiddleware(config).(*CORSMiddleware)
	engine := setupTestEngine()
	engine.Use(corsMiddleware.GetGinHandler())
	
	engine.GET("/api/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{"message": "test"})
	})
	
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("Origin", "https://example.com")
	engine.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

// TestCORSWithEmptyOriginsList 测试空Origins列表的默认行为 (Test default behavior with empty origins list)
func TestCORSWithEmptyOriginsList(t *testing.T) {
	config := &server.CORSConfig{
		Enabled:          true,
		AllowOrigins:     []string{}, // 空列表 (Empty list)
		AllowMethods:     []string{"GET", "POST"},
		AllowCredentials: false,
	}
	
	corsMiddleware := NewCORSMiddleware(config).(*CORSMiddleware)
	engine := setupTestEngine()
	engine.Use(corsMiddleware.GetGinHandler())
	
	engine.GET("/api/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{"message": "test"})
	})
	
	// 测试任意域请求 (Test any domain request)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/test", nil)
	req.Header.Set("Origin", "https://anydomain.com")
	engine.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	// 空Origins列表时应允许所有源并回显具体origin (Empty origins list should allow all origins and echo specific origin)
	assert.Equal(t, "https://anydomain.com", w.Header().Get("Access-Control-Allow-Origin"))
}

// BenchmarkCORSMiddleware 基准测试CORS中间件 (Benchmark CORS middleware)
func BenchmarkCORSMiddleware(b *testing.B) {
	config := &server.CORSConfig{
		Enabled:      true,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST"},
		MaxAge:       1 * time.Hour,
	}
	
	corsMiddleware := NewCORSMiddleware(config).(*CORSMiddleware)
	engine := setupTestEngine()
	engine.Use(corsMiddleware.GetGinHandler())
	
	engine.GET("/bench", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{"message": "benchmark"})
	})
	
	req := httptest.NewRequest("GET", "/bench", nil)
	req.Header.Set("Origin", "https://example.com")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
	}
} 