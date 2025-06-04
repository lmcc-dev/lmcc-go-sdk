/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Gin服务器适配器单元测试 (Gin server adapter unit tests)
 */

package gin

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
)

// TestNewGinServer 测试创建Gin服务器 (Test creating Gin server)
func TestNewGinServer(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "gin",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}
	
	ginServer := NewGinServer(config)
	
	if ginServer == nil {
		t.Fatal("NewGinServer() returned nil")
	}
	
	if ginServer.GetConfig().Framework != "gin" {
		t.Error("Server framework should be 'gin'")
	}
	
	if ginServer.GetGinEngine() == nil {
		t.Error("Gin engine should not be nil")
	}
	
	if ginServer.GetHTTPServer() == nil {
		t.Error("HTTP server should not be nil")
	}
}

// TestNewGinServerWithServices 测试创建带服务容器的Gin服务器 (Test creating Gin server with services)
func TestNewGinServerWithServices(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "gin",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	ginServer := NewGinServerWithServices(config, serviceContainer)
	
	if ginServer == nil {
		t.Fatal("NewGinServerWithServices() returned nil")
	}
	
	if ginServer.GetServices() == nil {
		t.Error("Service container should not be nil")
	}
	
	if ginServer.GetServices() != serviceContainer {
		t.Error("Service container should match the provided one")
	}
}

// TestGinServerRegisterRoute 测试注册路由 (Test registering routes)
func TestGinServerRegisterRoute(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "gin",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}
	
	ginServer := NewGinServer(config)
	
	tests := []struct {
		name     string
		method   string
		path     string
		handler  server.Handler
		wantErr  bool
	}{
		{
			name:   "Valid GET route",
			method: "GET",
			path:   "/users",
			handler: server.HandlerFunc(func(ctx server.Context) error {
				return ctx.JSON(http.StatusOK, map[string]string{"message": "users"})
			}),
			wantErr: false,
		},
		{
			name:   "Valid POST route",
			method: "POST",
			path:   "/users",
			handler: server.HandlerFunc(func(ctx server.Context) error {
				return ctx.JSON(http.StatusCreated, map[string]string{"message": "created"})
			}),
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ginServer.RegisterRoute(tt.method, tt.path, tt.handler)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterRoute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	
	// 测试路由 (Test route)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	ginServer.GetGinEngine().ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

// TestGinServerRegisterMiddleware 测试注册中间件 (Test registering middleware)
func TestGinServerRegisterMiddleware(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "gin",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}
	
	ginServer := NewGinServer(config)
	
	// 创建测试中间件 (Create test middleware)
	middlewareCalled := false
	middleware := server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
		middlewareCalled = true
		return next()
	})
	
	// 注册中间件 (Register middleware)
	err := ginServer.RegisterMiddleware(middleware)
	if err != nil {
		t.Fatalf("RegisterMiddleware() failed: %v", err)
	}
	
	// 注册测试路由 (Register test route)
	handler := server.HandlerFunc(func(ctx server.Context) error {
		executed, _ := ctx.Get("middleware_executed")
		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"middleware_executed": executed,
		})
	})
	
	ginServer.RegisterRoute("GET", "/test", handler)
	
	// 测试中间件 (Test middleware)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	ginServer.GetGinEngine().ServeHTTP(w, req)
	
	if !middlewareCalled {
		t.Error("Middleware should have been called")
	}
}

// TestGinServerGroup 测试路由组 (Test route groups)
func TestGinServerGroup(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "gin",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}
	
	ginServer := NewGinServer(config)
	
	// 创建路由组 (Create route group)
	apiGroup := ginServer.Group("/api")
	
	if apiGroup == nil {
		t.Fatal("Group() returned nil")
	}
	
	// 在组中注册路由 (Register route in group)
	handler := server.HandlerFunc(func(ctx server.Context) error {
		return ctx.JSON(http.StatusOK, map[string]string{"message": "api test"})
	})
	
	err := apiGroup.RegisterRoute("GET", "/test", handler)
	if err != nil {
		t.Fatalf("RegisterRoute() in group failed: %v", err)
	}
	
	// 测试组路由 (Test group route)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/test", nil)
	ginServer.GetGinEngine().ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

// TestGinServerGroupWithMiddleware 测试带中间件的路由组 (Test route groups with middleware)
func TestGinServerGroupWithMiddleware(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "gin",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}
	
	ginServer := NewGinServer(config)
	
	// 创建测试中间件 (Create test middleware)
	middlewareCalled := false
	middleware := server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
		middlewareCalled = true
		return next()
	})
	
	// 创建带中间件的路由组 (Create route group with middleware)
	apiGroup := ginServer.Group("/api", middleware)
	
	// 在组中注册路由 (Register route in group)
	handler := server.HandlerFunc(func(ctx server.Context) error {
		return ctx.JSON(http.StatusOK, map[string]string{"message": "api test"})
	})
	
	apiGroup.RegisterRoute("GET", "/test", handler)
	
	// 测试组路由和中间件 (Test group route and middleware)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/test", nil)
	ginServer.GetGinEngine().ServeHTTP(w, req)
	
	if !middlewareCalled {
		t.Error("Group middleware should have been called")
	}
}

// TestGinServerGetNativeEngine 测试获取原生引擎 (Test getting native engine)
func TestGinServerGetNativeEngine(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "gin",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}
	
	ginServer := NewGinServer(config)
	
	nativeEngine := ginServer.GetNativeEngine()
	if nativeEngine == nil {
		t.Fatal("GetNativeEngine() returned nil")
	}
	
	// 检查是否为Gin引擎 (Check if it's a Gin engine)
	ginEngine, ok := nativeEngine.(*gin.Engine)
	if !ok {
		t.Error("Native engine should be *gin.Engine")
	}
	
	if ginEngine != ginServer.GetGinEngine() {
		t.Error("Native engine should match Gin engine")
	}
}

// TestGinServerStartStop 测试服务器启动和停止 (Test server start and stop)
func TestGinServerStartStop(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "gin",
		Host:      "localhost",
		Port:      0, // 使用随机端口 (Use random port)
		Mode:      "test",
	}
	
	ginServer := NewGinServer(config)
	
	// 启动服务器 (Start server)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	err := ginServer.Start(ctx)
	if err != nil {
		t.Fatalf("Start() failed: %v", err)
	}
	
	// 给服务器一点时间启动 (Give server some time to start)
	time.Sleep(100 * time.Millisecond)
	
	// 停止服务器 (Stop server)
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer stopCancel()
	
	err = ginServer.Stop(stopCtx)
	if err != nil {
		t.Fatalf("Stop() failed: %v", err)
	}
}

// TestGinServerApplyGinConfig 测试应用Gin特定配置 (Test applying Gin-specific configuration)
func TestGinServerApplyGinConfig(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "gin",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
		Plugins: map[string]interface{}{
			"gin": map[string]interface{}{
				"trusted_proxies":           []string{"127.0.0.1"},
				"redirect_trailing_slash":   false,
				"redirect_fixed_path":       true,
				"handle_method_not_allowed": true,
				"max_multipart_memory":      int64(64 << 20), // 64MB
			},
		},
	}
	
	ginServer := NewGinServer(config)
	engine := ginServer.GetGinEngine()
	
	// 检查配置是否应用 (Check if configuration is applied)
	if engine.RedirectTrailingSlash != false {
		t.Error("RedirectTrailingSlash should be false")
	}
	
	if engine.RedirectFixedPath != true {
		t.Error("RedirectFixedPath should be true")
	}
	
	if engine.HandleMethodNotAllowed != true {
		t.Error("HandleMethodNotAllowed should be true")
	}
	
	if engine.MaxMultipartMemory != 64<<20 {
		t.Error("MaxMultipartMemory should be 64MB")
	}
}

// BenchmarkGinServerRegisterRoute 基准测试路由注册 (Benchmark route registration)
func BenchmarkGinServerRegisterRoute(b *testing.B) {
	config := &server.ServerConfig{
		Framework: "gin",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}
	
	ginServer := NewGinServer(config)
	handler := server.HandlerFunc(func(ctx server.Context) error {
		return ctx.JSON(http.StatusOK, map[string]string{"message": "test"})
	})
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		path := "/test" + string(rune(i))
		err := ginServer.RegisterRoute("GET", path, handler)
		if err != nil {
			b.Fatalf("RegisterRoute() failed: %v", err)
		}
	}
}

// BenchmarkGinServerHandleRequest 基准测试请求处理 (Benchmark request handling)
func BenchmarkGinServerHandleRequest(b *testing.B) {
	config := &server.ServerConfig{
		Framework: "gin",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}
	
	ginServer := NewGinServer(config)
	handler := server.HandlerFunc(func(ctx server.Context) error {
		return ctx.JSON(http.StatusOK, map[string]string{"message": "test"})
	})
	
	ginServer.RegisterRoute("GET", "/test", handler)
	
	req, _ := http.NewRequest("GET", "/test", nil)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		ginServer.GetGinEngine().ServeHTTP(w, req)
	}
} 