/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Echo路由组适配器实现 (Echo route group adapter implementation)
 */

package echo

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
)

// EchoRouteGroup Echo路由组适配器 (Echo route group adapter)
// 实现server.RouteGroup接口 (Implements server.RouteGroup interface)
type EchoRouteGroup struct {
	group            *echo.Group               // Echo路由组 (Echo route group)
	server           *EchoServer               // Echo服务器 (Echo server)
	serviceContainer services.ServiceContainer // 服务容器 (Service container)
	prefix           string                    // 路由前缀 (Route prefix)
	logger           services.Logger           // 日志服务 (Logger service)
}

// NewEchoRouteGroup 创建Echo路由组适配器 (Create Echo route group adapter)
func NewEchoRouteGroup(group *echo.Group, server *EchoServer) server.RouteGroup {
	return &EchoRouteGroup{
		group:            group,
		server:           server,
		serviceContainer: server.serviceContainer,
		prefix:           "",  // 需要从外部传入或通过其他方式获取 (Need to pass from external or get through other means)
		logger:           server.serviceContainer.GetLogger(),
	}
}

// NewEchoRouteGroupWithPrefix 创建带前缀的Echo路由组适配器 (Create Echo route group adapter with prefix)
func NewEchoRouteGroupWithPrefix(group *echo.Group, server *EchoServer, prefix string) server.RouteGroup {
	return &EchoRouteGroup{
		group:            group,
		server:           server,
		serviceContainer: server.serviceContainer,
		prefix:           prefix,
		logger:           server.serviceContainer.GetLogger(),
	}
}

// RegisterRoute 在组内注册路由 (Register route within group)
func (g *EchoRouteGroup) RegisterRoute(method, path string, handler server.Handler) error {
	// 将统一处理器转换为Echo处理器 (Convert unified handler to Echo handler)
	echoHandler := g.wrapHandler(handler)
	
	// 注册路由到路由组 (Register route to route group)
	g.group.Add(method, path, echoHandler)
	
	// 计算完整路径 (Calculate full path)
	fullPath := g.prefix + path
	
	g.logger.Debugw("Route registered in group",
		"method", method,
		"path", path,
		"full_path", fullPath,
		"group_prefix", g.prefix,
		"framework", "echo",
	)
	
	return nil
}

// RegisterMiddleware 注册组级中间件 (Register group-level middleware)
func (g *EchoRouteGroup) RegisterMiddleware(middleware server.Middleware) error {
	// 将统一中间件转换为Echo中间件 (Convert unified middleware to Echo middleware)
	echoMiddleware := g.wrapMiddleware(middleware)
	
	// 添加中间件到路由组 (Add middleware to route group)
	g.group.Use(echoMiddleware)
	
	g.logger.Debugw("Middleware registered in group",
		"group_prefix", g.prefix,
		"framework", "echo",
	)
	
	return nil
}

// Group 创建子路由组 (Create sub route group)
func (g *EchoRouteGroup) Group(prefix string, middlewares ...server.Middleware) server.RouteGroup {
	// 创建Echo子路由组 (Create Echo sub route group)
	subGroup := g.group.Group(prefix)
	
	// 添加中间件到子路由组 (Add middleware to sub route group)
	for _, mw := range middlewares {
		echoMiddleware := g.wrapMiddleware(mw)
		subGroup.Use(echoMiddleware)
	}
	
	// 创建子路由组适配器 (Create sub route group adapter)
	return NewEchoRouteGroup(subGroup, g.server)
}

// wrapHandler 包装统一处理器为Echo处理器 (Wrap unified handler to Echo handler)
func (g *EchoRouteGroup) wrapHandler(handler server.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 创建统一上下文 (Create unified context)
		ctx := NewEchoContext(c, g.serviceContainer)
		
		// 调用统一处理器 (Call unified handler)
		return handler.Handle(ctx)
	}
}

// wrapMiddleware 包装统一中间件为Echo中间件 (Wrap unified middleware to Echo middleware)
func (g *EchoRouteGroup) wrapMiddleware(middleware server.Middleware) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 创建统一上下文 (Create unified context)
			ctx := NewEchoContext(c, g.serviceContainer)
			
			// 定义下一个处理器 (Define next handler)
			nextHandler := func() error {
				return next(c)
			}
			
			// 调用统一中间件 (Call unified middleware)
			return middleware.Process(ctx, nextHandler)
		}
	}
}

// GetPrefix 获取路由组前缀 (Get route group prefix)
func (g *EchoRouteGroup) GetPrefix() string {
	return g.prefix
}

// GetRoutes 获取路由组中的所有路由 (Get all routes in the route group)
func (g *EchoRouteGroup) GetRoutes() []RouteInfo {
	routes := make([]RouteInfo, 0)
	
	// 获取Echo引擎的路由信息 (Get route information from Echo engine)
	if g.server != nil {
		if engine := g.server.GetEchoEngine(); engine != nil {
			for _, route := range engine.Routes() {
				// 检查路由是否属于当前组 (Check if route belongs to current group)
				if strings.HasPrefix(route.Path, g.prefix) {
					routes = append(routes, RouteInfo{
						Method:      route.Method,
						Path:        route.Path,
						HandlerName: route.Name,
					})
				}
			}
		}
	}
	
	return routes
}

// Use 添加中间件到路由组 (Add middleware to route group)
// 这是RegisterMiddleware的别名，提供更简洁的API (This is an alias for RegisterMiddleware, providing a more concise API)
func (g *EchoRouteGroup) Use(middlewares ...server.Middleware) error {
	for _, middleware := range middlewares {
		if err := g.RegisterMiddleware(middleware); err != nil {
			return err
		}
	}
	return nil
}

// Handle 注册任意HTTP方法的路由 (Register route for any HTTP method)
func (g *EchoRouteGroup) Handle(method, path string, handler server.Handler) error {
	return g.RegisterRoute(method, path, handler)
}

// GET 注册GET路由的便捷方法 (Convenience method for registering GET routes)
func (g *EchoRouteGroup) GET(path string, handler server.Handler) error {
	return g.RegisterRoute(http.MethodGet, path, handler)
}

// POST 注册POST路由的便捷方法 (Convenience method for registering POST routes)
func (g *EchoRouteGroup) POST(path string, handler server.Handler) error {
	return g.RegisterRoute(http.MethodPost, path, handler)
}

// PUT 注册PUT路由的便捷方法 (Convenience method for registering PUT routes)
func (g *EchoRouteGroup) PUT(path string, handler server.Handler) error {
	return g.RegisterRoute(http.MethodPut, path, handler)
}

// DELETE 注册DELETE路由的便捷方法 (Convenience method for registering DELETE routes)
func (g *EchoRouteGroup) DELETE(path string, handler server.Handler) error {
	return g.RegisterRoute(http.MethodDelete, path, handler)
}

// PATCH 注册PATCH路由的便捷方法 (Convenience method for registering PATCH routes)
func (g *EchoRouteGroup) PATCH(path string, handler server.Handler) error {
	return g.RegisterRoute(http.MethodPatch, path, handler)
}

// OPTIONS 注册OPTIONS路由的便捷方法 (Convenience method for registering OPTIONS routes)
func (g *EchoRouteGroup) OPTIONS(path string, handler server.Handler) error {
	return g.RegisterRoute(http.MethodOptions, path, handler)
}

// HEAD 注册HEAD路由的便捷方法 (Convenience method for registering HEAD routes)
func (g *EchoRouteGroup) HEAD(path string, handler server.Handler) error {
	return g.RegisterRoute(http.MethodHead, path, handler)
}

// Any 注册所有HTTP方法的路由 (Register route for all HTTP methods)
func (g *EchoRouteGroup) Any(path string, handler server.Handler) error {
	methods := []string{
		http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete,
		http.MethodPatch, http.MethodOptions, http.MethodHead,
	}
	
	for _, method := range methods {
		if err := g.RegisterRoute(method, path, handler); err != nil {
			return g.serviceContainer.GetErrorHandler().Wrap(err, 
				fmt.Sprintf("failed to register %s route", method))
		}
	}
	
	return nil
}

// RouteInfo 路由信息结构 (Route information structure)
type RouteInfo struct {
	Method      string `json:"method"`      // HTTP方法 (HTTP method)
	Path        string `json:"path"`        // 路由路径 (Route path)
	HandlerName string `json:"handler_name"` // 处理器名称 (Handler name)
}

// String 返回路由信息的字符串表示 (Return string representation of route info)
func (r RouteInfo) String() string {
	return fmt.Sprintf("%s %s -> %s", r.Method, r.Path, r.HandlerName)
}

// 验证EchoRouteGroup实现了server.RouteGroup接口 (Verify EchoRouteGroup implements server.RouteGroup interface)
var _ server.RouteGroup = (*EchoRouteGroup)(nil) 