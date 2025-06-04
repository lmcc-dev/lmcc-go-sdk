/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Fiber路由适配器实现 (Fiber router adapter implementation)
 */

package fiber

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
)

// FiberRouteGroup Fiber路由组适配器 (Fiber route group adapter)
// 实现server.RouteGroup接口 (Implements server.RouteGroup interface)
type FiberRouteGroup struct {
	app              fiber.Router                 // Fiber路由器 (Fiber router)
	prefix           string                       // 路由前缀 (Route prefix)
	serviceContainer services.ServiceContainer    // 服务容器 (Service container)
	routes           []RouteInfo                  // 路由信息 (Route information)
}

// NewFiberRouteGroup 创建Fiber路由组适配器 (Create Fiber route group adapter)
func NewFiberRouteGroup(app fiber.Router, serviceContainer services.ServiceContainer) server.RouteGroup {
	return &FiberRouteGroup{
		app:              app,
		prefix:           "",
		serviceContainer: serviceContainer,
		routes:           make([]RouteInfo, 0),
	}
}

// NewFiberRouteGroupWithPrefix 创建带前缀的Fiber路由组适配器 (Create Fiber route group adapter with prefix)
func NewFiberRouteGroupWithPrefix(app fiber.Router, prefix string, serviceContainer services.ServiceContainer) server.RouteGroup {
	return &FiberRouteGroup{
		app:              app,
		prefix:           prefix,
		serviceContainer: serviceContainer,
		routes:           make([]RouteInfo, 0),
	}
}

// RegisterRoute 在组内注册路由 (Register route within group)
func (r *FiberRouteGroup) RegisterRoute(method, path string, handler server.Handler) error {
	// 包装处理器
	fiberHandler := r.wrapHandler(handler)
	
	// 注册路由
	r.app.Add(method, path, fiberHandler)
	
	// 记录路由信息
	r.routes = append(r.routes, RouteInfo{
		Method:   method,
		Path:     r.prefix + path,
		Handler:  "fiber.handler",
		Metadata: map[string]interface{}{
			"framework": "fiber",
			"prefix":    r.prefix,
		},
	})
	
	return nil
}

// RegisterMiddleware 注册组级中间件 (Register group-level middleware)
func (r *FiberRouteGroup) RegisterMiddleware(middleware server.Middleware) error {
	r.app.Use(r.wrapMiddleware(middleware))
	return nil
}

// Group 创建子路由组 (Create sub route group)
func (r *FiberRouteGroup) Group(prefix string, middlewares ...server.Middleware) server.RouteGroup {
	// 创建Fiber子组
	group := r.app.Group(prefix)
	
	// 应用中间件
	for _, mw := range middlewares {
		group.Use(r.wrapMiddleware(mw))
	}
	
	// 计算完整前缀
	fullPrefix := r.prefix + prefix
	
	return NewFiberRouteGroupWithPrefix(group, fullPrefix, r.serviceContainer)
}

// Use 注册中间件 (Register middleware)
func (r *FiberRouteGroup) Use(middleware ...server.MiddlewareFunc) {
	for _, mw := range middleware {
		r.app.Use(r.wrapMiddlewareFunc(mw))
	}
}

// Handle 注册路由处理器 (Register route handler)
func (r *FiberRouteGroup) Handle(method, path string, handler server.HandlerFunc, middleware ...server.MiddlewareFunc) {
	// 包装处理器
	fiberHandler := r.wrapHandlerFunc(handler)
	
	// 应用中间件
	handlers := make([]fiber.Handler, 0, len(middleware)+1)
	for _, mw := range middleware {
		handlers = append(handlers, r.wrapMiddlewareFunc(mw))
	}
	handlers = append(handlers, fiberHandler)
	
	// 注册路由
	r.app.Add(method, path, handlers...)
	
	// 记录路由信息
	r.routes = append(r.routes, RouteInfo{
		Method:   method,
		Path:     r.prefix + path,
		Handler:  "fiber.handler",
		Metadata: map[string]interface{}{
			"framework": "fiber",
			"prefix":    r.prefix,
		},
	})
}

// GET 注册GET路由 (Register GET route)
func (r *FiberRouteGroup) GET(path string, handler server.HandlerFunc, middleware ...server.MiddlewareFunc) {
	r.Handle("GET", path, handler, middleware...)
}

// POST 注册POST路由 (Register POST route)
func (r *FiberRouteGroup) POST(path string, handler server.HandlerFunc, middleware ...server.MiddlewareFunc) {
	r.Handle("POST", path, handler, middleware...)
}

// PUT 注册PUT路由 (Register PUT route)
func (r *FiberRouteGroup) PUT(path string, handler server.HandlerFunc, middleware ...server.MiddlewareFunc) {
	r.Handle("PUT", path, handler, middleware...)
}

// DELETE 注册DELETE路由 (Register DELETE route)
func (r *FiberRouteGroup) DELETE(path string, handler server.HandlerFunc, middleware ...server.MiddlewareFunc) {
	r.Handle("DELETE", path, handler, middleware...)
}

// PATCH 注册PATCH路由 (Register PATCH route)
func (r *FiberRouteGroup) PATCH(path string, handler server.HandlerFunc, middleware ...server.MiddlewareFunc) {
	r.Handle("PATCH", path, handler, middleware...)
}

// HEAD 注册HEAD路由 (Register HEAD route)
func (r *FiberRouteGroup) HEAD(path string, handler server.HandlerFunc, middleware ...server.MiddlewareFunc) {
	r.Handle("HEAD", path, handler, middleware...)
}

// OPTIONS 注册OPTIONS路由 (Register OPTIONS route)
func (r *FiberRouteGroup) OPTIONS(path string, handler server.HandlerFunc, middleware ...server.MiddlewareFunc) {
	r.Handle("OPTIONS", path, handler, middleware...)
}

// GetRoutes 获取路由信息 (Get route information)
func (r *FiberRouteGroup) GetRoutes() []RouteInfo {
	return r.routes
}

// GetPrefix 获取路由前缀 (Get route prefix)
func (r *FiberRouteGroup) GetPrefix() string {
	return r.prefix
}

// GetFiberRouter 获取原生Fiber路由器 (Get native Fiber router)
func (r *FiberRouteGroup) GetFiberRouter() fiber.Router {
	return r.app
}

// wrapHandler 包装处理器接口 (Wrap handler interface)
func (r *FiberRouteGroup) wrapHandler(handler server.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := NewFiberContext(c, r.serviceContainer)
		return handler.Handle(ctx)
	}
}

// wrapHandlerFunc 包装处理器函数 (Wrap handler function)
func (r *FiberRouteGroup) wrapHandlerFunc(handler server.HandlerFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := NewFiberContext(c, r.serviceContainer)
		return handler(ctx)
	}
}

// wrapMiddleware 包装中间件接口 (Wrap middleware interface)
func (r *FiberRouteGroup) wrapMiddleware(middleware server.Middleware) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := NewFiberContext(c, r.serviceContainer)
		return middleware.Process(ctx, func() error {
			return c.Next()
		})
	}
}

// wrapMiddlewareFunc 包装中间件函数 (Wrap middleware function)
func (r *FiberRouteGroup) wrapMiddlewareFunc(middleware server.MiddlewareFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := NewFiberContext(c, r.serviceContainer)
		return middleware(ctx, func() error {
			return c.Next()
		})
	}
}

// RouteInfo 路由信息结构 (Route information structure)
type RouteInfo struct {
	Method   string                 // HTTP方法 (HTTP method)
	Path     string                 // 路由路径 (Route path)
	Handler  string                 // 处理器名称 (Handler name)
	Metadata map[string]interface{} // 元数据 (Metadata)
}

// String 返回路由信息的字符串表示 (Return string representation of route info)
func (r RouteInfo) String() string {
	return fmt.Sprintf("%s %s -> %s", r.Method, r.Path, r.Handler)
} 