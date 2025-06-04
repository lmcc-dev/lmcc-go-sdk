/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Echo服务器适配器测试 (Echo server adapter tests)
 */

package echo

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
)

func TestNewEchoServer_Success(t *testing.T) {
	config := &server.ServerConfig{
		Framework:      "echo",
		Host:           "localhost",
		Port:           8080,
		Mode:           "debug",
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
		CORS: server.CORSConfig{
			Enabled: true,
		},
		Middleware: server.MiddlewareConfig{
			Logger: server.LoggerMiddlewareConfig{
				Enabled: true,
			},
			Recovery: server.RecoveryMiddlewareConfig{
				Enabled: true,
			},
		},
		GracefulShutdown: server.GracefulShutdownConfig{
			Enabled: true,
			Timeout: 30 * time.Second,
		},
	}

	serviceContainer := services.NewServiceContainerWithDefaults()
	
	echoServer, err := NewEchoServer(config, serviceContainer)
	assert.NoError(t, err)
	assert.NotNil(t, echoServer)
	
	// 验证配置 (Verify configuration)
	assert.Equal(t, config, echoServer.GetConfig())
	
	// 验证Echo引擎 (Verify Echo engine)
	engine := echoServer.GetEchoEngine()
	assert.NotNil(t, engine)
	assert.True(t, engine.HideBanner)
	assert.True(t, engine.HidePort)
	assert.True(t, engine.Debug) // 因为是debug模式 (Because it's debug mode)
}

func TestNewEchoServer_ReleaseMode(t *testing.T) {
	config := &server.ServerConfig{
		Framework:      "echo",
		Host:           "localhost",
		Port:           8080,
		Mode:           "release", // 发布模式 (Release mode)
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	serviceContainer := services.NewServiceContainerWithDefaults()
	
	echoServer, err := NewEchoServer(config, serviceContainer)
	assert.NoError(t, err)
	assert.NotNil(t, echoServer)
	
	engine := echoServer.GetEchoEngine()
	assert.False(t, engine.Debug) // 发布模式下应该关闭debug (Debug should be off in release mode)
}

func TestEchoServer_RegisterRoute(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "echo",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}

	serviceContainer := services.NewServiceContainerWithDefaults()
	echoServer, err := NewEchoServer(config, serviceContainer)
	require.NoError(t, err)

	// 创建测试处理器 (Create test handler)
	handler := &testHandler{message: "test"}
	
	// 注册路由 (Register route)
	err = echoServer.RegisterRoute("GET", "/test", handler)
	assert.NoError(t, err)
	
	// 验证路由已注册 (Verify route is registered)
	engine := echoServer.GetEchoEngine()
	routes := engine.Routes()
	
	found := false
	for _, route := range routes {
		if route.Method == "GET" && route.Path == "/test" {
			found = true
			break
		}
	}
	assert.True(t, found, "Route should be registered")
}

func TestEchoServer_RegisterMiddleware(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "echo",
		Host:      "localhost", 
		Port:      8080,
		Mode:      "debug",
	}

	serviceContainer := services.NewServiceContainerWithDefaults()
	echoServer, err := NewEchoServer(config, serviceContainer)
	require.NoError(t, err)

	// 创建测试中间件 (Create test middleware)
	middleware := &testMiddleware{name: "test"}
	
	// 注册中间件 (Register middleware)
	err = echoServer.RegisterMiddleware(middleware)
	assert.NoError(t, err)
}

func TestEchoServer_Group(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "echo",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}

	serviceContainer := services.NewServiceContainerWithDefaults()
	echoServer, err := NewEchoServer(config, serviceContainer)
	require.NoError(t, err)

	// 创建路由组 (Create route group)
	group := echoServer.Group("/api")
	assert.NotNil(t, group)
	
	// 验证路由组类型 (Verify route group type)
	echoGroup, ok := group.(*EchoRouteGroup)
	assert.True(t, ok)
	assert.Equal(t, "/api", echoGroup.GetPrefix())
}

func TestEchoServer_GroupWithMiddleware(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "echo",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}

	serviceContainer := services.NewServiceContainerWithDefaults()
	echoServer, err := NewEchoServer(config, serviceContainer)
	require.NoError(t, err)

	// 创建测试中间件 (Create test middleware)
	middleware := &testMiddleware{name: "group-test"}
	
	// 创建带中间件的路由组 (Create route group with middleware)
	group := echoServer.Group("/api", middleware)
	assert.NotNil(t, group)
	
	echoGroup, ok := group.(*EchoRouteGroup)
	assert.True(t, ok)
	assert.Equal(t, "/api", echoGroup.GetPrefix())
}

func TestEchoServer_GetNativeEngine(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "echo",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}

	serviceContainer := services.NewServiceContainerWithDefaults()
	echoServer, err := NewEchoServer(config, serviceContainer)
	require.NoError(t, err)

	// 获取原生引擎 (Get native engine)
	nativeEngine := echoServer.GetNativeEngine()
	assert.NotNil(t, nativeEngine)
	
	// 验证类型 (Verify type)
	engine, ok := nativeEngine.(*echo.Echo)
	assert.True(t, ok)
	assert.Equal(t, echoServer.GetEchoEngine(), engine)
}

func TestEchoServer_StartStop(t *testing.T) {
	config := &server.ServerConfig{
		Framework:      "echo",
		Host:           "127.0.0.1",
		Port:           9999, // 使用不常用端口避免冲突 (Use uncommon port to avoid conflicts)
		Mode:           "debug",
		ReadTimeout:    1 * time.Second,
		WriteTimeout:   1 * time.Second,
		IdleTimeout:    1 * time.Second,
		MaxHeaderBytes: 1 << 10,
		GracefulShutdown: server.GracefulShutdownConfig{
			Enabled: true,
			Timeout: 5 * time.Second,
		},
	}

	serviceContainer := services.NewServiceContainerWithDefaults()
	echoServer, err := NewEchoServer(config, serviceContainer)
	require.NoError(t, err)

	// 启动服务器 (Start server)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	startErr := make(chan error, 1)
	go func() {
		startErr <- echoServer.Start(ctx)
	}()
	
	// 等待启动 (Wait for startup)
	time.Sleep(100 * time.Millisecond)
	
	// 停止服务器 (Stop server)
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer stopCancel()
	
	err = echoServer.Stop(stopCtx)
	assert.NoError(t, err)
	
	// 验证启动没有错误 (Verify startup had no errors)
	select {
	case err := <-startErr:
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("Unexpected start error: %v", err)
		}
	case <-time.After(3 * time.Second):
		// Timeout is acceptable
	}
}

func TestEchoServer_StopNotStarted(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "echo",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}

	serviceContainer := services.NewServiceContainerWithDefaults()
	echoServer, err := NewEchoServer(config, serviceContainer)
	require.NoError(t, err)

	// 停止未启动的服务器应该没有错误 (Stopping non-started server should have no error)
	ctx := context.Background()
	err = echoServer.Stop(ctx)
	assert.NoError(t, err)
}

func TestEchoServer_WrapHandler(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "echo",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}

	serviceContainer := services.NewServiceContainerWithDefaults()
	echoServer, err := NewEchoServer(config, serviceContainer)
	require.NoError(t, err)

	// 创建测试处理器 (Create test handler)
	handler := &testHandler{message: "wrap-test"}
	
	// 包装处理器 (Wrap handler)
	echoHandler := echoServer.wrapHandler(handler)
	assert.NotNil(t, echoHandler)
}

func TestEchoServer_WrapMiddleware(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "echo",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}

	serviceContainer := services.NewServiceContainerWithDefaults()
	echoServer, err := NewEchoServer(config, serviceContainer)
	require.NoError(t, err)

	// 创建测试中间件 (Create test middleware)
	middleware := &testMiddleware{name: "wrap-test"}
	
	// 包装中间件 (Wrap middleware)
	echoMiddleware := echoServer.wrapMiddleware(middleware)
	assert.NotNil(t, echoMiddleware)
}

func TestEchoServer_SetupMiddleware(t *testing.T) {
	tests := []struct {
		name   string
		config *server.ServerConfig
	}{
		{
			name: "CORS enabled",
			config: &server.ServerConfig{
				Framework: "echo",
				Host:      "localhost",
				Port:      8080,
				Mode:      "debug",
				CORS: server.CORSConfig{
					Enabled:      true,
					AllowOrigins: []string{"*"},
					AllowMethods: []string{"GET", "POST"},
				},
				Middleware: server.MiddlewareConfig{
					Recovery: server.RecoveryMiddlewareConfig{
						Enabled: true,
					},
				},
			},
		},
		{
			name: "CORS disabled",
			config: &server.ServerConfig{
				Framework: "echo",
				Host:      "localhost",
				Port:      8080,
				Mode:      "debug",
				CORS: server.CORSConfig{
					Enabled: false,
				},
				Middleware: server.MiddlewareConfig{
					Recovery: server.RecoveryMiddlewareConfig{
						Enabled: false,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serviceContainer := services.NewServiceContainerWithDefaults()
			echoServer, err := NewEchoServer(tt.config, serviceContainer)
			assert.NoError(t, err)
			assert.NotNil(t, echoServer)
		})
	}
}

// 测试辅助类型 (Test helper types)

type testHandler struct {
	message string
}

func (h *testHandler) Handle(ctx server.Context) error {
	return ctx.String(200, h.message)
}

type testMiddleware struct {
	name string
}

func (m *testMiddleware) Process(ctx server.Context, next func() error) error {
	// 在请求前添加标记 (Add marker before request)
	ctx.Set("middleware_"+m.name, true)
	
	// 调用下一个处理器 (Call next handler)
	return next()
} 