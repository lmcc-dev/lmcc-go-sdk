/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Gin路由管理实现 (Gin router management implementation)
 */

package gin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
)

// GinRouteGroup Gin路由组实现 (Gin route group implementation)
type GinRouteGroup struct {
	ginGroup    *gin.RouterGroup
	prefix      string
	middlewares []server.Middleware
	server      *GinServer
	services    services.ServiceContainer
}

// NewGinRouteGroup 创建Gin路由组 (Create Gin route group)
// 保留向后兼容性 (Maintain backward compatibility)
func NewGinRouteGroup(ginGroup *gin.RouterGroup, prefix string, ginServer *GinServer) *GinRouteGroup {
	return NewGinRouteGroupWithServices(ginGroup, prefix, ginServer, nil)
}

// NewGinRouteGroupWithServices 创建带服务容器的Gin路由组 (Create Gin route group with service container)
func NewGinRouteGroupWithServices(ginGroup *gin.RouterGroup, prefix string, ginServer *GinServer, serviceContainer services.ServiceContainer) *GinRouteGroup {
	// 如果没有提供服务容器，尝试从服务器获取 (If no service container provided, try to get from server)
	if serviceContainer == nil && ginServer != nil {
		serviceContainer = ginServer.GetServices()
	}
	
	// 如果还是没有，创建默认的 (If still none, create default)
	if serviceContainer == nil {
		serviceContainer = services.NewServiceContainerWithDefaults()
	}
	
	return &GinRouteGroup{
		ginGroup:    ginGroup,
		prefix:      prefix,
		middlewares: make([]server.Middleware, 0),
		server:      ginServer,
		services:    serviceContainer,
	}
}

// RegisterRoute 在组内注册路由 (Register route within group)
func (g *GinRouteGroup) RegisterRoute(method, path string, handler server.Handler) error {
	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}
	
	// 转换为Gin处理器 (Convert to Gin handler)
	ginHandler := g.adaptHandler(handler)
	
	// 注册路由 (Register route)
	switch strings.ToUpper(method) {
	case http.MethodGet:
		g.ginGroup.GET(path, ginHandler)
	case http.MethodPost:
		g.ginGroup.POST(path, ginHandler)
	case http.MethodPut:
		g.ginGroup.PUT(path, ginHandler)
	case http.MethodPatch:
		g.ginGroup.PATCH(path, ginHandler)
	case http.MethodDelete:
		g.ginGroup.DELETE(path, ginHandler)
	case http.MethodHead:
		g.ginGroup.HEAD(path, ginHandler)
	case http.MethodOptions:
		g.ginGroup.OPTIONS(path, ginHandler)
	default:
		return fmt.Errorf("unsupported HTTP method: %s", method)
	}
	
	return nil
}

// RegisterMiddleware 注册组级中间件 (Register group-level middleware)
func (g *GinRouteGroup) RegisterMiddleware(middleware server.Middleware) error {
	if middleware == nil {
		return fmt.Errorf("middleware cannot be nil")
	}
	
	// 添加到中间件列表 (Add to middleware list)
	g.middlewares = append(g.middlewares, middleware)
	
	// 转换为Gin中间件并应用 (Convert to Gin middleware and apply)
	ginMiddleware := g.adaptMiddleware(middleware)
	g.ginGroup.Use(ginMiddleware)
	
	return nil
}

// Group 创建子路由组 (Create sub route group)
func (g *GinRouteGroup) Group(prefix string, middlewares ...server.Middleware) server.RouteGroup {
	// 创建Gin子组 (Create Gin sub group)
	ginSubGroup := g.ginGroup.Group(prefix)
	
	// 创建子路由组 (Create sub route group)
	subGroup := NewGinRouteGroupWithServices(ginSubGroup, g.prefix+prefix, g.server, g.services)
	
	// 应用中间件 (Apply middlewares)
	for _, middleware := range middlewares {
		if err := subGroup.RegisterMiddleware(middleware); err != nil {
			// 记录错误但不中断创建过程 (Log error but don't interrupt creation process)
			if g.services != nil {
				logger := g.services.GetLogger()
				logger.Warnf("Failed to register middleware in route group: %v", err)
			}
			continue
		}
	}
	
	return subGroup
}

// adaptHandler 将统一Handler适配为Gin处理器 (Adapt unified Handler to Gin handler)
func (g *GinRouteGroup) adaptHandler(handler server.Handler) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		// 创建统一上下文 (Create unified context)
		ctx := NewGinContext(ginCtx)
		
		// 调用处理器 (Call handler)
		if err := handler.Handle(ctx); err != nil {
			// 处理错误 (Handle error)
			g.handleError(ginCtx, err)
		}
	}
}

// adaptMiddleware 将统一Middleware适配为Gin中间件 (Adapt unified Middleware to Gin middleware)
func (g *GinRouteGroup) adaptMiddleware(middleware server.Middleware) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		// 创建统一上下文 (Create unified context)
		ctx := NewGinContext(ginCtx)
		
		// 调用中间件 (Call middleware)
		err := middleware.Process(ctx, func() error {
			ginCtx.Next()
			return nil
		})
		
		// 处理错误 (Handle error)
		if err != nil {
			g.handleError(ginCtx, err)
		}
	}
}

// handleError 处理错误 (Handle error)
func (g *GinRouteGroup) handleError(ginCtx *gin.Context, err error) {
	// 使用服务容器的错误处理器和日志器 (Use service container's error handler and logger)
	if g.services != nil {
		errorHandler := g.services.GetErrorHandler()
		logger := g.services.GetLogger()
		
		// 记录错误 (Log error)
		logger.Errorf("Route group error: %v", err)
		
		// 设置错误到Gin上下文 (Set error to Gin context)
		ginCtx.Error(err)
		
		// 如果响应还没有写入，发送错误响应 (If response not written yet, send error response)
		if !ginCtx.Writer.Written() {
			ginCtx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal Server Error",
				"message": errorHandler.GetStackTrace(err),
			})
		}
	} else {
		// 回退到基本错误处理 (Fallback to basic error handling)
		ginCtx.Error(err)
		if !ginCtx.Writer.Written() {
			ginCtx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Internal Server Error",
				"message": err.Error(),
			})
		}
	}
	
	// 中止请求处理 (Abort request processing)
	ginCtx.Abort()
}

// GetPrefix 获取路由组前缀 (Get route group prefix)
func (g *GinRouteGroup) GetPrefix() string {
	return g.prefix
}

// GetMiddlewares 获取路由组中间件 (Get route group middlewares)
func (g *GinRouteGroup) GetMiddlewares() []server.Middleware {
	// 返回副本以防止外部修改 (Return copy to prevent external modification)
	result := make([]server.Middleware, len(g.middlewares))
	copy(result, g.middlewares)
	return result
}

// GetServices 获取服务容器 (Get service container)
func (g *GinRouteGroup) GetServices() services.ServiceContainer {
	return g.services
}

// GetGinGroup 获取原生Gin路由组 (Get native Gin route group)
// 用于访问Gin特定功能 (Used to access Gin-specific features)
func (g *GinRouteGroup) GetGinGroup() *gin.RouterGroup {
	return g.ginGroup
}

// RouteInfo 路由信息 (Route information)
type RouteInfo struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	HandlerName string `json:"handler_name"`
}

// GetRoutes 获取路由组中的所有路由 (Get all routes in the route group)
func (g *GinRouteGroup) GetRoutes() []RouteInfo {
	routes := make([]RouteInfo, 0)
	
	// 获取Gin引擎的路由信息 (Get route information from Gin engine)
	if g.server != nil {
		if engine := g.server.GetGinEngine(); engine != nil {
			for _, route := range engine.Routes() {
				// 检查路由是否属于当前组 (Check if route belongs to current group)
				if strings.HasPrefix(route.Path, g.prefix) {
					routes = append(routes, RouteInfo{
						Method:      route.Method,
						Path:        route.Path,
						HandlerName: route.Handler,
					})
				}
			}
		}
	}
	
	return routes
}

// Use 添加中间件到路由组 (Add middleware to route group)
// 这是RegisterMiddleware的别名，提供更简洁的API (This is an alias for RegisterMiddleware, providing a more concise API)
func (g *GinRouteGroup) Use(middlewares ...server.Middleware) error {
	for _, middleware := range middlewares {
		if err := g.RegisterMiddleware(middleware); err != nil {
			return err
		}
	}
	return nil
}

// Handle 注册任意HTTP方法的路由 (Register route for any HTTP method)
func (g *GinRouteGroup) Handle(method, path string, handler server.Handler) error {
	return g.RegisterRoute(method, path, handler)
}

// GET 注册GET路由的便捷方法 (Convenience method for registering GET routes)
func (g *GinRouteGroup) GET(path string, handler server.Handler) error {
	return g.RegisterRoute(http.MethodGet, path, handler)
}

// POST 注册POST路由的便捷方法 (Convenience method for registering POST routes)
func (g *GinRouteGroup) POST(path string, handler server.Handler) error {
	return g.RegisterRoute(http.MethodPost, path, handler)
}

// PUT 注册PUT路由的便捷方法 (Convenience method for registering PUT routes)
func (g *GinRouteGroup) PUT(path string, handler server.Handler) error {
	return g.RegisterRoute(http.MethodPut, path, handler)
}

// PATCH 注册PATCH路由的便捷方法 (Convenience method for registering PATCH routes)
func (g *GinRouteGroup) PATCH(path string, handler server.Handler) error {
	return g.RegisterRoute(http.MethodPatch, path, handler)
}

// DELETE 注册DELETE路由的便捷方法 (Convenience method for registering DELETE routes)
func (g *GinRouteGroup) DELETE(path string, handler server.Handler) error {
	return g.RegisterRoute(http.MethodDelete, path, handler)
}

// HEAD 注册HEAD路由的便捷方法 (Convenience method for registering HEAD routes)
func (g *GinRouteGroup) HEAD(path string, handler server.Handler) error {
	return g.RegisterRoute(http.MethodHead, path, handler)
}

// OPTIONS 注册OPTIONS路由的便捷方法 (Convenience method for registering OPTIONS routes)
func (g *GinRouteGroup) OPTIONS(path string, handler server.Handler) error {
	return g.RegisterRoute(http.MethodOptions, path, handler)
} 