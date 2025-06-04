/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Fiber路由器适配器测试 (Fiber router adapter tests)
 */

package fiber

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
)

func TestNewFiberRouteGroup(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	
	group := NewFiberRouteGroup(app, serviceContainer)
	assert.NotNil(t, group)
	
	fiberGroup, ok := group.(*FiberRouteGroup)
	assert.True(t, ok)
	assert.Equal(t, "", fiberGroup.GetPrefix())
	assert.Equal(t, serviceContainer, fiberGroup.serviceContainer)
}

func TestNewFiberRouteGroupWithPrefix(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	prefix := "/api"
	
	group := NewFiberRouteGroupWithPrefix(app, prefix, serviceContainer)
	assert.NotNil(t, group)
	
	fiberGroup, ok := group.(*FiberRouteGroup)
	assert.True(t, ok)
	assert.Equal(t, prefix, fiberGroup.GetPrefix())
	assert.Equal(t, serviceContainer, fiberGroup.serviceContainer)
}

func TestFiberRouteGroup_RegisterRoute(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	group := NewFiberRouteGroup(app, serviceContainer)
	
	handler := server.HandlerFunc(func(ctx server.Context) error {
		return ctx.JSON(200, map[string]string{"message": "test"})
	})
	
	err := group.RegisterRoute("GET", "/test", handler)
	assert.NoError(t, err)
	
	// 验证路由信息
	fiberGroup := group.(*FiberRouteGroup)
	routes := fiberGroup.GetRoutes()
	assert.Len(t, routes, 1)
	assert.Equal(t, "GET", routes[0].Method)
	assert.Equal(t, "/test", routes[0].Path)
}

func TestFiberRouteGroup_RegisterMiddleware(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	group := NewFiberRouteGroup(app, serviceContainer)
	
	middleware := server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
		ctx.SetHeader("X-Test", "middleware")
		return next()
	})
	
	err := group.RegisterMiddleware(middleware)
	assert.NoError(t, err)
}

func TestFiberRouteGroup_Group(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	group := NewFiberRouteGroup(app, serviceContainer)
	
	subGroup := group.Group("/api")
	assert.NotNil(t, subGroup)
	
	fiberSubGroup, ok := subGroup.(*FiberRouteGroup)
	assert.True(t, ok)
	assert.Equal(t, "/api", fiberSubGroup.GetPrefix())
}

func TestFiberRouteGroup_GroupWithMiddleware(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	group := NewFiberRouteGroup(app, serviceContainer)
	
	middleware := server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
		ctx.SetHeader("X-Group", "middleware")
		return next()
	})
	
	subGroup := group.Group("/api", middleware)
	assert.NotNil(t, subGroup)
	
	fiberSubGroup, ok := subGroup.(*FiberRouteGroup)
	assert.True(t, ok)
	assert.Equal(t, "/api", fiberSubGroup.GetPrefix())
}

func TestFiberRouteGroup_Handle(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	group := NewFiberRouteGroup(app, serviceContainer)
	
	handler := server.HandlerFunc(func(ctx server.Context) error {
		return ctx.JSON(200, map[string]string{"message": "test"})
	})
	
	fiberGroup := group.(*FiberRouteGroup)
	fiberGroup.Handle("POST", "/test", handler)
	
	routes := fiberGroup.GetRoutes()
	assert.Len(t, routes, 1)
	assert.Equal(t, "POST", routes[0].Method)
	assert.Equal(t, "/test", routes[0].Path)
}

func TestFiberRouteGroup_HandleWithMiddleware(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	group := NewFiberRouteGroup(app, serviceContainer)
	
	handler := server.HandlerFunc(func(ctx server.Context) error {
		return ctx.JSON(200, map[string]string{"message": "test"})
	})
	
	middleware := server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
		ctx.SetHeader("X-Middleware", "test")
		return next()
	})
	
	fiberGroup := group.(*FiberRouteGroup)
	fiberGroup.Handle("POST", "/test", handler, middleware)
	
	routes := fiberGroup.GetRoutes()
	assert.Len(t, routes, 1)
	assert.Equal(t, "POST", routes[0].Method)
	assert.Equal(t, "/test", routes[0].Path)
}

func TestFiberRouteGroup_GET(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	group := NewFiberRouteGroup(app, serviceContainer)
	
	handler := server.HandlerFunc(func(ctx server.Context) error {
		return ctx.JSON(200, map[string]string{"message": "get"})
	})
	
	fiberGroup := group.(*FiberRouteGroup)
	fiberGroup.GET("/test", handler)
	
	routes := fiberGroup.GetRoutes()
	assert.Len(t, routes, 1)
	assert.Equal(t, "GET", routes[0].Method)
}

func TestFiberRouteGroup_POST(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	group := NewFiberRouteGroup(app, serviceContainer)
	
	handler := server.HandlerFunc(func(ctx server.Context) error {
		return ctx.JSON(200, map[string]string{"message": "post"})
	})
	
	fiberGroup := group.(*FiberRouteGroup)
	fiberGroup.POST("/test", handler)
	
	routes := fiberGroup.GetRoutes()
	assert.Len(t, routes, 1)
	assert.Equal(t, "POST", routes[0].Method)
}

func TestFiberRouteGroup_PUT(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	group := NewFiberRouteGroup(app, serviceContainer)
	
	handler := server.HandlerFunc(func(ctx server.Context) error {
		return ctx.JSON(200, map[string]string{"message": "put"})
	})
	
	fiberGroup := group.(*FiberRouteGroup)
	fiberGroup.PUT("/test", handler)
	
	routes := fiberGroup.GetRoutes()
	assert.Len(t, routes, 1)
	assert.Equal(t, "PUT", routes[0].Method)
}

func TestFiberRouteGroup_DELETE(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	group := NewFiberRouteGroup(app, serviceContainer)
	
	handler := server.HandlerFunc(func(ctx server.Context) error {
		return ctx.JSON(200, map[string]string{"message": "delete"})
	})
	
	fiberGroup := group.(*FiberRouteGroup)
	fiberGroup.DELETE("/test", handler)
	
	routes := fiberGroup.GetRoutes()
	assert.Len(t, routes, 1)
	assert.Equal(t, "DELETE", routes[0].Method)
}

func TestFiberRouteGroup_PATCH(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	group := NewFiberRouteGroup(app, serviceContainer)
	
	handler := server.HandlerFunc(func(ctx server.Context) error {
		return ctx.JSON(200, map[string]string{"message": "patch"})
	})
	
	fiberGroup := group.(*FiberRouteGroup)
	fiberGroup.PATCH("/test", handler)
	
	routes := fiberGroup.GetRoutes()
	assert.Len(t, routes, 1)
	assert.Equal(t, "PATCH", routes[0].Method)
}

func TestFiberRouteGroup_HEAD(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	group := NewFiberRouteGroup(app, serviceContainer)
	
	handler := server.HandlerFunc(func(ctx server.Context) error {
		return ctx.JSON(200, map[string]string{"message": "head"})
	})
	
	fiberGroup := group.(*FiberRouteGroup)
	fiberGroup.HEAD("/test", handler)
	
	routes := fiberGroup.GetRoutes()
	assert.Len(t, routes, 1)
	assert.Equal(t, "HEAD", routes[0].Method)
}

func TestFiberRouteGroup_OPTIONS(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	group := NewFiberRouteGroup(app, serviceContainer)
	
	handler := server.HandlerFunc(func(ctx server.Context) error {
		return ctx.JSON(200, map[string]string{"message": "options"})
	})
	
	fiberGroup := group.(*FiberRouteGroup)
	fiberGroup.OPTIONS("/test", handler)
	
	routes := fiberGroup.GetRoutes()
	assert.Len(t, routes, 1)
	assert.Equal(t, "OPTIONS", routes[0].Method)
}

func TestFiberRouteGroup_Use(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	group := NewFiberRouteGroup(app, serviceContainer)
	
	middleware := server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
		ctx.SetHeader("X-Use", "middleware")
		return next()
	})
	
	fiberGroup := group.(*FiberRouteGroup)
	fiberGroup.Use(middleware)
	
	// 测试中间件是否正确注册（通过实际请求验证）
	handler := server.HandlerFunc(func(ctx server.Context) error {
		return ctx.JSON(200, map[string]string{"message": "test"})
	})
	
	fiberGroup.GET("/test", handler)
	
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestFiberRouteGroup_GetFiberRouter(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	group := NewFiberRouteGroup(app, serviceContainer)
	
	fiberGroup := group.(*FiberRouteGroup)
	router := fiberGroup.GetFiberRouter()
	assert.NotNil(t, router)
	assert.Equal(t, app, router)
}

func TestFiberRouteGroup_NestedGroups(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	group := NewFiberRouteGroup(app, serviceContainer)
	
	// 创建嵌套路由组
	apiGroup := group.Group("/api")
	v1Group := apiGroup.Group("/v1")
	
	fiberV1Group := v1Group.(*FiberRouteGroup)
	assert.Equal(t, "/api/v1", fiberV1Group.GetPrefix())
	
	// 在嵌套组中注册路由
	handler := server.HandlerFunc(func(ctx server.Context) error {
		return ctx.JSON(200, map[string]string{"message": "nested"})
	})
	
	err := v1Group.RegisterRoute("GET", "/users", handler)
	assert.NoError(t, err)
	
	routes := fiberV1Group.GetRoutes()
	assert.Len(t, routes, 1)
	assert.Equal(t, "GET", routes[0].Method)
	assert.Equal(t, "/api/v1/users", routes[0].Path)
}

func TestRouteInfo_String(t *testing.T) {
	routeInfo := RouteInfo{
		Method:  "GET",
		Path:    "/api/users",
		Handler: "fiber.handler",
		Metadata: map[string]interface{}{
			"framework": "fiber",
		},
	}
	
	expected := "GET /api/users -> fiber.handler"
	assert.Equal(t, expected, routeInfo.String())
} 