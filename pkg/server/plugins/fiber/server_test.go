/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Fiber服务器适配器测试 (Fiber server adapter tests)
 */

package fiber

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
)

func TestNewFiberServer(t *testing.T) {
	config := &server.ServerConfig{
		Framework:    "fiber",
		Host:         "localhost",
		Port:         8080,
		Mode:         "debug",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	serviceContainer := services.NewServiceContainer()
	fiberServer, err := NewFiberServer(config, serviceContainer)
	
	assert.NoError(t, err)
	assert.NotNil(t, fiberServer)
	assert.Equal(t, config, fiberServer.GetConfig())
	assert.NotNil(t, fiberServer.GetFiberApp())
}

func TestFiberServer_GetNativeEngine(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "fiber",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}
	
	serviceContainer := services.NewServiceContainer()
	fiberServer, err := NewFiberServer(config, serviceContainer)
	require.NoError(t, err)
	
	engine := fiberServer.GetNativeEngine()
	assert.NotNil(t, engine)
	
	// 验证返回的是Fiber应用实例
	fiberApp := fiberServer.GetFiberApp()
	assert.Equal(t, fiberApp, engine)
}

func TestFiberServer_RegisterRoute(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "fiber",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}
	
	serviceContainer := services.NewServiceContainer()
	fiberServer, err := NewFiberServer(config, serviceContainer)
	require.NoError(t, err)
	
	// 创建测试处理器
	handler := server.HandlerFunc(func(ctx server.Context) error {
		return ctx.JSON(200, map[string]string{"message": "test"})
	})
	
	// 注册路由
	err = fiberServer.RegisterRoute("GET", "/test", handler)
	assert.NoError(t, err)
}

func TestFiberServer_RegisterMiddleware(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "fiber",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}
	
	serviceContainer := services.NewServiceContainer()
	fiberServer, err := NewFiberServer(config, serviceContainer)
	require.NoError(t, err)
	
	// 创建测试中间件
	middleware := server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
		ctx.SetHeader("X-Test", "middleware")
		return next()
	})
	
	// 注册中间件
	err = fiberServer.RegisterMiddleware(middleware)
	assert.NoError(t, err)
}

func TestFiberServer_Group(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "fiber",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}
	
	serviceContainer := services.NewServiceContainer()
	fiberServer, err := NewFiberServer(config, serviceContainer)
	require.NoError(t, err)
	
	// 创建路由组
	group := fiberServer.Group("/api")
	assert.NotNil(t, group)
	
	// 验证路由组类型
	fiberGroup, ok := group.(*FiberRouteGroup)
	assert.True(t, ok)
	assert.Equal(t, "/api", fiberGroup.GetPrefix())
}

func TestFiberServer_GroupWithMiddleware(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "fiber",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}
	
	serviceContainer := services.NewServiceContainer()
	fiberServer, err := NewFiberServer(config, serviceContainer)
	require.NoError(t, err)
	
	// 创建测试中间件
	middleware := server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
		ctx.SetHeader("X-Group", "middleware")
		return next()
	})
	
	// 创建带中间件的路由组
	group := fiberServer.Group("/api", middleware)
	assert.NotNil(t, group)
	
	// 验证路由组类型
	fiberGroup, ok := group.(*FiberRouteGroup)
	assert.True(t, ok)
	assert.Equal(t, "/api", fiberGroup.GetPrefix())
}

func TestFiberServer_Stop_WithoutStart(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "fiber",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}
	
	serviceContainer := services.NewServiceContainer()
	fiberServer, err := NewFiberServer(config, serviceContainer)
	require.NoError(t, err)
	
	// 测试在未启动的情况下停止服务器
	ctx := context.Background()
	err = fiberServer.Stop(ctx)
	assert.NoError(t, err)
} 