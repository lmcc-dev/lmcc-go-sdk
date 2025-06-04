/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Echo CORS中间件测试 (Echo CORS middleware tests)
 */

package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
)

func TestDefaultCORSConfig(t *testing.T) {
	config := DefaultCORSConfig()
	
	assert.True(t, config.Enabled)
	assert.Equal(t, []string{"*"}, config.AllowOrigins)
	assert.Contains(t, config.AllowMethods, "GET")
	assert.Contains(t, config.AllowMethods, "POST")
	assert.Equal(t, []string{"*"}, config.AllowHeaders)
	assert.False(t, config.AllowCredentials)
	assert.Equal(t, 86400, config.MaxAge)
}

func TestNewCORSMiddleware(t *testing.T) {
	// 创建模拟服务容器 (Create mock service container)
	serviceContainer := services.NewServiceContainerWithDefaults()
	
	// 测试使用默认配置 (Test with default config)
	middleware := NewCORSMiddleware(nil, serviceContainer)
	assert.NotNil(t, middleware)
	assert.NotNil(t, middleware.config)
	assert.True(t, middleware.config.Enabled)
	
	// 测试使用自定义配置 (Test with custom config)
	customConfig := &CORSConfig{
		Enabled:      false,
		AllowOrigins: []string{"http://localhost:3000"},
	}
	middleware = NewCORSMiddleware(customConfig, serviceContainer)
	assert.NotNil(t, middleware)
	assert.False(t, middleware.config.Enabled)
	assert.Equal(t, []string{"http://localhost:3000"}, middleware.config.AllowOrigins)
}

func TestCORSMiddleware_Handler(t *testing.T) {
	serviceContainer := services.NewServiceContainerWithDefaults()
	
	tests := []struct {
		name     string
		config   *CORSConfig
		enabled  bool
	}{
		{
			name:    "enabled middleware",
			config:  DefaultCORSConfig(),
			enabled: true,
		},
		{
			name: "disabled middleware",
			config: &CORSConfig{
				Enabled: false,
			},
			enabled: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := NewCORSMiddleware(tt.config, serviceContainer)
			handler := middleware.Handler()
			
			assert.NotNil(t, handler)
			
			// 创建Echo实例进行测试 (Create Echo instance for testing)
			e := echo.New()
			e.Use(handler)
			e.GET("/test", func(c echo.Context) error {
				return c.String(http.StatusOK, "test")
			})
			
			// 创建测试请求 (Create test request)
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("Origin", "http://localhost:3000")
			rec := httptest.NewRecorder()
			
			// 执行请求 (Execute request)
			e.ServeHTTP(rec, req)
			
			// 验证响应 (Verify response)
			assert.Equal(t, http.StatusOK, rec.Code)
			
			if tt.enabled {
				// 验证CORS头部存在 (Verify CORS headers exist)
				assert.NotEmpty(t, rec.Header().Get("Access-Control-Allow-Origin"))
			}
		})
	}
}

func TestCORSMiddleware_PreflightRequest(t *testing.T) {
	serviceContainer := services.NewServiceContainerWithDefaults()
	config := DefaultCORSConfig()
	middleware := NewCORSMiddleware(config, serviceContainer)
	
	// 创建Echo实例 (Create Echo instance)
	e := echo.New()
	e.Use(middleware.Handler())
	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})
	
	// 创建预检请求 (Create preflight request)
	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "POST")
	rec := httptest.NewRecorder()
	
	// 执行请求 (Execute request)
	e.ServeHTTP(rec, req)
	
	// 验证预检响应 (Verify preflight response)
	assert.Equal(t, http.StatusNoContent, rec.Code)
	assert.NotEmpty(t, rec.Header().Get("Access-Control-Allow-Origin"))
	assert.NotEmpty(t, rec.Header().Get("Access-Control-Allow-Methods"))
}

func TestCORSMiddleware_IsOriginAllowed(t *testing.T) {
	serviceContainer := services.NewServiceContainerWithDefaults()
	
	tests := []struct {
		name         string
		allowOrigins []string
		origin       string
		expected     bool
	}{
		{
			name:         "wildcard allows all",
			allowOrigins: []string{"*"},
			origin:       "http://localhost:3000",
			expected:     true,
		},
		{
			name:         "specific origin allowed",
			allowOrigins: []string{"http://localhost:3000", "https://example.com"},
			origin:       "http://localhost:3000",
			expected:     true,
		},
		{
			name:         "origin not allowed",
			allowOrigins: []string{"https://example.com"},
			origin:       "http://localhost:3000",
			expected:     false,
		},
		{
			name:         "empty origin",
			allowOrigins: []string{"*"},
			origin:       "",
			expected:     false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &CORSConfig{
				Enabled:      true,
				AllowOrigins: tt.allowOrigins,
			}
			middleware := NewCORSMiddleware(config, serviceContainer)
			
			result := middleware.isOriginAllowed(tt.origin)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCORSMiddleware_GetConfig(t *testing.T) {
	serviceContainer := services.NewServiceContainerWithDefaults()
	config := DefaultCORSConfig()
	middleware := NewCORSMiddleware(config, serviceContainer)
	
	retrievedConfig := middleware.GetConfig()
	assert.Equal(t, config, retrievedConfig)
}

func TestCORSMiddleware_SetConfig(t *testing.T) {
	serviceContainer := services.NewServiceContainerWithDefaults()
	middleware := NewCORSMiddleware(nil, serviceContainer)
	
	// 设置新配置 (Set new config)
	newConfig := &CORSConfig{
		Enabled:      false,
		AllowOrigins: []string{"https://example.com"},
	}
	middleware.SetConfig(newConfig)
	
	assert.Equal(t, newConfig, middleware.config)
	
	// 测试设置nil配置 (Test setting nil config)
	middleware.SetConfig(nil)
	assert.Equal(t, newConfig, middleware.config) // 应该保持不变
} 