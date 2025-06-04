/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Echo路由组适配器测试 (Echo route group adapter tests)
 */

package echo

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
)

func TestNewEchoRouteGroup(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "echo",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}

	serviceContainer := services.NewServiceContainerWithDefaults()
	echoServer, err := NewEchoServer(config, serviceContainer)
	require.NoError(t, err)

	e := echoServer.GetEchoEngine()
	group := e.Group("/api")

	// 创建路由组适配器 (Create route group adapter)
	routeGroup := NewEchoRouteGroup(group, echoServer)
	assert.NotNil(t, routeGroup)

	// 验证接口实现 (Verify interface implementation)
	_, ok := routeGroup.(server.RouteGroup)
	assert.True(t, ok)

	// 验证类型转换 (Verify type conversion)
	echoRouteGroup, ok := routeGroup.(*EchoRouteGroup)
	assert.True(t, ok)
	assert.Equal(t, echoServer, echoRouteGroup.server)
}

func TestNewEchoRouteGroupWithPrefix(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "echo",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}

	serviceContainer := services.NewServiceContainerWithDefaults()
	echoServer, err := NewEchoServer(config, serviceContainer)
	require.NoError(t, err)

	e := echoServer.GetEchoEngine()
	group := e.Group("/api")

	// 创建带前缀的路由组适配器 (Create route group adapter with prefix)
	routeGroup := NewEchoRouteGroupWithPrefix(group, echoServer, "/api")
	assert.NotNil(t, routeGroup)

	echoRouteGroup, ok := routeGroup.(*EchoRouteGroup)
	assert.True(t, ok)
	assert.Equal(t, "/api", echoRouteGroup.GetPrefix())
}

func TestEchoRouteGroup_RegisterRoute(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "echo",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}

	serviceContainer := services.NewServiceContainerWithDefaults()
	echoServer, err := NewEchoServer(config, serviceContainer)
	require.NoError(t, err)

	e := echoServer.GetEchoEngine()
	group := e.Group("/api")
	routeGroup := NewEchoRouteGroupWithPrefix(group, echoServer, "/api")

	// 创建测试处理器 (Create test handler)
	handler := &testHandler{message: "group-test"}

	// 注册路由 (Register route)
	err = routeGroup.RegisterRoute("GET", "/users", handler)
	assert.NoError(t, err)

	// 验证路由已注册 (Verify route is registered)
	routes := e.Routes()
	found := false
	for _, route := range routes {
		if route.Method == "GET" && route.Path == "/api/users" {
			found = true
			break
		}
	}
	assert.True(t, found, "Route should be registered in group")
}

func TestEchoRouteGroup_RegisterMiddleware(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "echo",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}

	serviceContainer := services.NewServiceContainerWithDefaults()
	echoServer, err := NewEchoServer(config, serviceContainer)
	require.NoError(t, err)

	e := echoServer.GetEchoEngine()
	group := e.Group("/api")
	routeGroup := NewEchoRouteGroupWithPrefix(group, echoServer, "/api")

	// 创建测试中间件 (Create test middleware)
	middleware := &testMiddleware{name: "group-middleware"}

	// 注册中间件 (Register middleware)
	err = routeGroup.RegisterMiddleware(middleware)
	assert.NoError(t, err)
}

func TestEchoRouteGroup_Group(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "echo",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}

	serviceContainer := services.NewServiceContainerWithDefaults()
	echoServer, err := NewEchoServer(config, serviceContainer)
	require.NoError(t, err)

	e := echoServer.GetEchoEngine()
	group := e.Group("/api")
	routeGroup := NewEchoRouteGroupWithPrefix(group, echoServer, "/api")

	// 创建子路由组 (Create sub route group)
	subGroup := routeGroup.Group("/v1")
	assert.NotNil(t, subGroup)

	// 验证子路由组类型 (Verify sub route group type)
	echoSubGroup, ok := subGroup.(*EchoRouteGroup)
	assert.True(t, ok)
	assert.Equal(t, echoServer, echoSubGroup.server)
}

func TestEchoRouteGroup_GroupWithMiddleware(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "echo",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}

	serviceContainer := services.NewServiceContainerWithDefaults()
	echoServer, err := NewEchoServer(config, serviceContainer)
	require.NoError(t, err)

	e := echoServer.GetEchoEngine()
	group := e.Group("/api")
	routeGroup := NewEchoRouteGroupWithPrefix(group, echoServer, "/api")

	// 创建测试中间件 (Create test middleware)
	middleware := &testMiddleware{name: "sub-group-middleware"}

	// 创建带中间件的子路由组 (Create sub route group with middleware)
	subGroup := routeGroup.Group("/v1", middleware)
	assert.NotNil(t, subGroup)

	echoSubGroup, ok := subGroup.(*EchoRouteGroup)
	assert.True(t, ok)
	assert.Equal(t, echoServer, echoSubGroup.server)
}

func TestEchoRouteGroup_GetPrefix(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "echo",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}

	serviceContainer := services.NewServiceContainerWithDefaults()
	echoServer, err := NewEchoServer(config, serviceContainer)
	require.NoError(t, err)

	e := echoServer.GetEchoEngine()
	group := e.Group("/api")
	routeGroup := NewEchoRouteGroupWithPrefix(group, echoServer, "/api")

	// 测试GetPrefix方法 (Test GetPrefix method)
	echoRouteGroup, ok := routeGroup.(*EchoRouteGroup)
	require.True(t, ok)
	prefix := echoRouteGroup.GetPrefix()
	assert.Equal(t, "/api", prefix)
}

func TestEchoRouteGroup_GetRoutes(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "echo",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}

	serviceContainer := services.NewServiceContainerWithDefaults()
	echoServer, err := NewEchoServer(config, serviceContainer)
	require.NoError(t, err)

	e := echoServer.GetEchoEngine()
	group := e.Group("/api")
	routeGroup := NewEchoRouteGroupWithPrefix(group, echoServer, "/api")

	// 注册一些路由 (Register some routes)
	handler := &testHandler{message: "test"}
	err = routeGroup.RegisterRoute("GET", "/users", handler)
	require.NoError(t, err)
	err = routeGroup.RegisterRoute("POST", "/users", handler)
	require.NoError(t, err)

	// 获取路由信息 (Get route information)
	echoRouteGroup, ok := routeGroup.(*EchoRouteGroup)
	require.True(t, ok)
	routes := echoRouteGroup.GetRoutes()
	assert.GreaterOrEqual(t, len(routes), 2)

	// 验证路由信息 (Verify route information)
	foundGet := false
	foundPost := false
	for _, route := range routes {
		if route.Method == "GET" && route.Path == "/api/users" {
			foundGet = true
		}
		if route.Method == "POST" && route.Path == "/api/users" {
			foundPost = true
		}
	}
	assert.True(t, foundGet, "GET route should be found")
	assert.True(t, foundPost, "POST route should be found")
}

func TestEchoRouteGroup_Use(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "echo",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}

	serviceContainer := services.NewServiceContainerWithDefaults()
	echoServer, err := NewEchoServer(config, serviceContainer)
	require.NoError(t, err)

	e := echoServer.GetEchoEngine()
	group := e.Group("/api")
	routeGroup := NewEchoRouteGroupWithPrefix(group, echoServer, "/api")

	// 创建测试中间件 (Create test middleware)
	middleware1 := &testMiddleware{name: "middleware1"}
	middleware2 := &testMiddleware{name: "middleware2"}

	// 使用Use方法添加多个中间件 (Use the Use method to add multiple middleware)
	echoRouteGroup, ok := routeGroup.(*EchoRouteGroup)
	require.True(t, ok)
	err = echoRouteGroup.Use(middleware1, middleware2)
	assert.NoError(t, err)
}

func TestEchoRouteGroup_HTTPMethods(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "echo",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}

	serviceContainer := services.NewServiceContainerWithDefaults()
	echoServer, err := NewEchoServer(config, serviceContainer)
	require.NoError(t, err)

	e := echoServer.GetEchoEngine()
	group := e.Group("/api")
	routeGroup := NewEchoRouteGroupWithPrefix(group, echoServer, "/api")

	echoRouteGroup, ok := routeGroup.(*EchoRouteGroup)
	require.True(t, ok)

	handler := &testHandler{message: "test"}

	// 测试所有HTTP方法 (Test all HTTP methods)
	tests := []struct {
		method   string
		function func(string, server.Handler) error
	}{
		{"GET", echoRouteGroup.GET},
		{"POST", echoRouteGroup.POST},
		{"PUT", echoRouteGroup.PUT},
		{"DELETE", echoRouteGroup.DELETE},
		{"PATCH", echoRouteGroup.PATCH},
		{"OPTIONS", echoRouteGroup.OPTIONS},
		{"HEAD", echoRouteGroup.HEAD},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			path := "/" + tt.method + "_test"
			err := tt.function(path, handler)
			assert.NoError(t, err)

			// 验证路由已注册 (Verify route is registered)
			routes := e.Routes()
			found := false
			for _, route := range routes {
				if route.Method == tt.method && route.Path == "/api"+path {
					found = true
					break
				}
			}
			assert.True(t, found, "Route should be registered for method "+tt.method)
		})
	}
}

func TestEchoRouteGroup_Handle(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "echo",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}

	serviceContainer := services.NewServiceContainerWithDefaults()
	echoServer, err := NewEchoServer(config, serviceContainer)
	require.NoError(t, err)

	e := echoServer.GetEchoEngine()
	group := e.Group("/api")
	routeGroup := NewEchoRouteGroupWithPrefix(group, echoServer, "/api")

	echoRouteGroup, ok := routeGroup.(*EchoRouteGroup)
	require.True(t, ok)

	handler := &testHandler{message: "handle-test"}

	// 测试Handle方法 (Test Handle method)
	err = echoRouteGroup.Handle("CUSTOM", "/custom", handler)
	assert.NoError(t, err)

	// 验证路由已注册 (Verify route is registered)
	routes := e.Routes()
	found := false
	for _, route := range routes {
		if route.Method == "CUSTOM" && route.Path == "/api/custom" {
			found = true
			break
		}
	}
	assert.True(t, found, "Custom method route should be registered")
}

func TestEchoRouteGroup_Any(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "echo",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}

	serviceContainer := services.NewServiceContainerWithDefaults()
	echoServer, err := NewEchoServer(config, serviceContainer)
	require.NoError(t, err)

	e := echoServer.GetEchoEngine()
	group := e.Group("/api")
	routeGroup := NewEchoRouteGroupWithPrefix(group, echoServer, "/api")

	echoRouteGroup, ok := routeGroup.(*EchoRouteGroup)
	require.True(t, ok)

	handler := &testHandler{message: "any-test"}

	// 测试Any方法 (Test Any method)
	err = echoRouteGroup.Any("/any", handler)
	assert.NoError(t, err)

	// 验证所有方法的路由都已注册 (Verify routes for all methods are registered)
	routes := e.Routes()
	expectedMethods := []string{
		http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete,
		http.MethodPatch, http.MethodOptions, http.MethodHead,
	}

	for _, method := range expectedMethods {
		found := false
		for _, route := range routes {
			if route.Method == method && route.Path == "/api/any" {
				found = true
				break
			}
		}
		assert.True(t, found, "Route should be registered for method "+method)
	}
}

func TestEchoRouteGroup_WrapHandler(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "echo",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}

	serviceContainer := services.NewServiceContainerWithDefaults()
	echoServer, err := NewEchoServer(config, serviceContainer)
	require.NoError(t, err)

	e := echoServer.GetEchoEngine()
	group := e.Group("/api")
	routeGroup := NewEchoRouteGroupWithPrefix(group, echoServer, "/api")

	echoRouteGroup, ok := routeGroup.(*EchoRouteGroup)
	require.True(t, ok)

	// 创建测试处理器 (Create test handler)
	handler := &testHandler{message: "wrap-test"}

	// 包装处理器 (Wrap handler)
	echoHandler := echoRouteGroup.wrapHandler(handler)
	assert.NotNil(t, echoHandler)
}

func TestEchoRouteGroup_WrapMiddleware(t *testing.T) {
	config := &server.ServerConfig{
		Framework: "echo",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}

	serviceContainer := services.NewServiceContainerWithDefaults()
	echoServer, err := NewEchoServer(config, serviceContainer)
	require.NoError(t, err)

	e := echoServer.GetEchoEngine()
	group := e.Group("/api")
	routeGroup := NewEchoRouteGroupWithPrefix(group, echoServer, "/api")

	echoRouteGroup, ok := routeGroup.(*EchoRouteGroup)
	require.True(t, ok)

	// 创建测试中间件 (Create test middleware)
	middleware := &testMiddleware{name: "wrap-test"}

	// 包装中间件 (Wrap middleware)
	echoMiddleware := echoRouteGroup.wrapMiddleware(middleware)
	assert.NotNil(t, echoMiddleware)
}

func TestRouteInfo_String(t *testing.T) {
	routeInfo := RouteInfo{
		Method:      "GET",
		Path:        "/api/users",
		HandlerName: "getUsersHandler",
	}

	expected := "GET /api/users -> getUsersHandler"
	assert.Equal(t, expected, routeInfo.String())
} 