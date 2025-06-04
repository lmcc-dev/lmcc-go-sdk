/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Gin路由管理单元测试 (Gin router management unit tests)
 */

package gin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
	"github.com/stretchr/testify/assert"
)

// TestNewGinRouteGroup 测试创建Gin路由组 (Test creating Gin route group)
func TestNewGinRouteGroup(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	ginGroup := engine.Group("/api")
	
	config := &server.ServerConfig{
		Framework: "gin",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}
	
	ginServer := NewGinServer(config)
	routeGroup := NewGinRouteGroup(ginGroup, "/api", ginServer)
	
	assert.NotNil(t, routeGroup)
	assert.Equal(t, "/api", routeGroup.prefix)
	assert.NotNil(t, routeGroup.ginGroup)
	assert.NotNil(t, routeGroup.server)
}

// TestNewGinRouteGroupWithServices 测试创建带服务容器的Gin路由组 (Test creating Gin route group with services)
func TestNewGinRouteGroupWithServices(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	ginGroup := engine.Group("/api")
	
	config := &server.ServerConfig{
		Framework: "gin",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	ginServer := NewGinServerWithServices(config, serviceContainer)
	routeGroup := NewGinRouteGroupWithServices(ginGroup, "/api", ginServer, serviceContainer)
	
	assert.NotNil(t, routeGroup)
	assert.Equal(t, "/api", routeGroup.prefix)
	assert.NotNil(t, routeGroup.services)
	assert.Equal(t, serviceContainer, routeGroup.services)
}

// TestGinRouteGroupRegisterRoute 测试注册路由 (Test registering routes)
func TestGinRouteGroupRegisterRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	ginGroup := engine.Group("/api")
	
	config := &server.ServerConfig{
		Framework: "gin",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}
	
	ginServer := NewGinServer(config)
	routeGroup := NewGinRouteGroup(ginGroup, "/api", ginServer)
	
	// 创建测试处理器 (Create test handler)
	testHandler := server.HandlerFunc(func(ctx server.Context) error {
		return ctx.String(200, "test response")
	})
	
	// 测试注册GET路由 (Test registering GET route)
	err := routeGroup.RegisterRoute("GET", "/users", testHandler)
	assert.NoError(t, err)
	
	// 测试注册POST路由 (Test registering POST route)
	err = routeGroup.RegisterRoute("POST", "/users", testHandler)
	assert.NoError(t, err)
	
	// 测试注册带参数的路由 (Test registering route with parameters)
	err = routeGroup.RegisterRoute("GET", "/users/:id", testHandler)
	assert.NoError(t, err)
}

// TestGinRouteGroupRegisterRouteInvalidMethod 测试注册无效方法的路由 (Test registering route with invalid method)
func TestGinRouteGroupRegisterRouteInvalidMethod(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	ginGroup := engine.Group("/api")
	
	config := &server.ServerConfig{
		Framework: "gin",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}
	
	ginServer := NewGinServer(config)
	routeGroup := NewGinRouteGroup(ginGroup, "/api", ginServer)
	
	testHandler := server.HandlerFunc(func(ctx server.Context) error {
		return ctx.String(200, "test response")
	})
	
	// 测试无效的HTTP方法 (Test invalid HTTP method)
	err := routeGroup.RegisterRoute("INVALID", "/test", testHandler)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported HTTP method")
}

// TestGinRouteGroupRegisterMiddleware 测试注册中间件 (Test registering middleware)
func TestGinRouteGroupRegisterMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	ginGroup := engine.Group("/api")
	
	config := &server.ServerConfig{
		Framework: "gin",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}
	
	ginServer := NewGinServer(config)
	routeGroup := NewGinRouteGroup(ginGroup, "/api", ginServer)
	
	// 创建测试中间件 (Create test middleware)
	testMiddleware := server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
		ctx.SetHeader("X-Test-Middleware", "applied")
		return next()
	})
	
	// 注册中间件 (Register middleware)
	err := routeGroup.RegisterMiddleware(testMiddleware)
	assert.NoError(t, err)
	assert.Len(t, routeGroup.middlewares, 1)
}

// TestGinRouteGroupGroup 测试创建子路由组 (Test creating sub route groups)
func TestGinRouteGroupGroup(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	ginGroup := engine.Group("/api")
	
	config := &server.ServerConfig{
		Framework: "gin",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}
	
	ginServer := NewGinServer(config)
	routeGroup := NewGinRouteGroup(ginGroup, "/api", ginServer)
	
	// 创建子路由组 (Create sub route group)
	subGroup := routeGroup.Group("/v1")
	assert.NotNil(t, subGroup)
	
	// 验证类型转换 (Verify type conversion)
	ginSubGroup, ok := subGroup.(*GinRouteGroup)
	assert.True(t, ok)
	assert.Equal(t, "/api/v1", ginSubGroup.prefix)
}

// TestGinRouteGroupGroupWithMiddleware 测试创建带中间件的子路由组 (Test creating sub route group with middleware)
func TestGinRouteGroupGroupWithMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	ginGroup := engine.Group("/api")
	
	config := &server.ServerConfig{
		Framework: "gin",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}
	
	ginServer := NewGinServer(config)
	routeGroup := NewGinRouteGroup(ginGroup, "/api", ginServer)
	
	// 创建测试中间件 (Create test middleware)
	testMiddleware1 := server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
		ctx.SetHeader("X-Middleware-1", "applied")
		return next()
	})
	
	testMiddleware2 := server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
		ctx.SetHeader("X-Middleware-2", "applied")
		return next()
	})
	
	// 创建带中间件的子路由组 (Create sub route group with middleware)
	subGroup := routeGroup.Group("/v1", testMiddleware1, testMiddleware2)
	assert.NotNil(t, subGroup)
	
	// 验证中间件数量 (Verify middleware count)
	ginSubGroup, ok := subGroup.(*GinRouteGroup)
	assert.True(t, ok)
	assert.Len(t, ginSubGroup.middlewares, 2)
}

// TestGinRouteGroupIntegration 测试路由组的集成功能 (Test route group integration)
func TestGinRouteGroupIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	ginGroup := engine.Group("/api")
	
	config := &server.ServerConfig{
		Framework: "gin",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}
	
	ginServer := NewGinServer(config)
	routeGroup := NewGinRouteGroup(ginGroup, "/api", ginServer)
	
	// 注册中间件 (Register middleware)
	testMiddleware := server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
		ctx.SetHeader("X-API-Version", "v1")
		return next()
	})
	
	err := routeGroup.RegisterMiddleware(testMiddleware)
	assert.NoError(t, err)
	
	// 注册路由 (Register route)
	testHandler := server.HandlerFunc(func(ctx server.Context) error {
		return ctx.JSON(200, map[string]string{"message": "success"})
	})
	
	err = routeGroup.RegisterRoute("GET", "/test", testHandler)
	assert.NoError(t, err)
	
	// 创建子路由组并注册路由 (Create sub group and register route)
	subGroup := routeGroup.Group("/v1")
	err = subGroup.RegisterRoute("POST", "/users", testHandler)
	assert.NoError(t, err)
	
	// 测试HTTP请求 (Test HTTP request)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/test", nil)
	engine.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "v1", w.Header().Get("X-API-Version"))
	assert.Contains(t, w.Body.String(), "success")
}

// TestGinRouteGroupErrorHandling 测试错误处理 (Test error handling)
func TestGinRouteGroupErrorHandling(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	ginGroup := engine.Group("/api")
	
	config := &server.ServerConfig{
		Framework: "gin",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}
	
	ginServer := NewGinServer(config)
	routeGroup := NewGinRouteGroup(ginGroup, "/api", ginServer)
	
	// 测试处理器返回错误 (Test handler returning error)
	errorHandler := server.HandlerFunc(func(ctx server.Context) error {
		return ctx.String(500, "internal error")
	})
	
	err := routeGroup.RegisterRoute("GET", "/error", errorHandler)
	assert.NoError(t, err)
	
	// 测试请求 (Test request)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/error", nil)
	engine.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "internal error")
}

// TestGinRouteGroupNilServer 测试nil服务器的情况 (Test nil server case)
func TestGinRouteGroupNilServer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	ginGroup := engine.Group("/api")
	
	// 使用nil服务器创建路由组 (Create route group with nil server)
	routeGroup := NewGinRouteGroup(ginGroup, "/api", nil)
	
	assert.NotNil(t, routeGroup)
	assert.Nil(t, routeGroup.server)
	assert.NotNil(t, routeGroup.services) // 应该有默认服务 (Should have default services)
}

// BenchmarkGinRouteGroupRegisterRoute 基准测试路由注册 (Benchmark route registration)
func BenchmarkGinRouteGroupRegisterRoute(b *testing.B) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	ginGroup := engine.Group("/api")
	
	config := &server.ServerConfig{
		Framework: "gin",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}
	
	ginServer := NewGinServer(config)
	routeGroup := NewGinRouteGroup(ginGroup, "/api", ginServer)
	
	testHandler := server.HandlerFunc(func(ctx server.Context) error {
		return ctx.String(200, "test")
	})
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		path := "/test" + string(rune(i%1000))
		routeGroup.RegisterRoute("GET", path, testHandler)
	}
}

// BenchmarkGinRouteGroupMiddleware 基准测试中间件性能 (Benchmark middleware performance)
func BenchmarkGinRouteGroupMiddleware(b *testing.B) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	ginGroup := engine.Group("/api")
	
	config := &server.ServerConfig{
		Framework: "gin",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}
	
	ginServer := NewGinServer(config)
	routeGroup := NewGinRouteGroup(ginGroup, "/api", ginServer)
	
	// 注册中间件 (Register middleware)
	testMiddleware := server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
		return next()
	})
	routeGroup.RegisterMiddleware(testMiddleware)
	
	// 注册路由 (Register route)
	testHandler := server.HandlerFunc(func(ctx server.Context) error {
		return ctx.String(200, "benchmark")
	})
	routeGroup.RegisterRoute("GET", "/benchmark", testHandler)
	
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/benchmark", nil)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Body.Reset()
		w.Code = 0
		engine.ServeHTTP(w, req)
	}
} 