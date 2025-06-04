/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: 定义Web框架插件化支持的核心接口
 */

package server

import (
	"context"
	"net/http"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
)

// WebFramework 统一的Web框架接口 (Unified web framework interface)
// 提供框架无关的服务器操作接口 (Provides framework-agnostic server operation interface)
type WebFramework interface {
	// Start 启动服务器 (Start the server)
	Start(ctx context.Context) error
	
	// Stop 停止服务器 (Stop the server)
	Stop(ctx context.Context) error
	
	// RegisterRoute 注册路由 (Register a route)
	RegisterRoute(method, path string, handler Handler) error
	
	// RegisterMiddleware 注册全局中间件 (Register global middleware)
	RegisterMiddleware(middleware Middleware) error
	
	// Group 创建路由组 (Create a route group)
	Group(prefix string, middlewares ...Middleware) RouteGroup
	
	// GetNativeEngine 获取原生框架实例 (Get native framework instance)
	// 用于访问框架特定功能 (Used to access framework-specific features)
	GetNativeEngine() interface{}
	
	// GetConfig 获取服务器配置 (Get server configuration)
	GetConfig() *ServerConfig
}

// Handler 统一的处理器接口 (Unified handler interface)
// 定义请求处理的标准方法 (Defines standard method for request handling)
type Handler interface {
	// Handle 处理请求 (Handle the request)
	Handle(ctx Context) error
}

// HandlerFunc 处理器函数类型 (Handler function type)
// 允许使用函数作为处理器 (Allows using functions as handlers)
type HandlerFunc func(ctx Context) error

// Handle 实现Handler接口 (Implement Handler interface)
func (f HandlerFunc) Handle(ctx Context) error {
	return f(ctx)
}

// Middleware 统一的中间件接口 (Unified middleware interface)
// 定义中间件的标准处理方法 (Defines standard processing method for middleware)
type Middleware interface {
	// Process 处理中间件逻辑 (Process middleware logic)
	// next: 调用下一个中间件或处理器 (Call next middleware or handler)
	Process(ctx Context, next func() error) error
}

// MiddlewareFunc 中间件函数类型 (Middleware function type)
// 允许使用函数作为中间件 (Allows using functions as middleware)
type MiddlewareFunc func(ctx Context, next func() error) error

// Process 实现Middleware接口 (Implement Middleware interface)
func (f MiddlewareFunc) Process(ctx Context, next func() error) error {
	return f(ctx, next)
}

// RouteGroup 路由组接口 (Route group interface)
// 提供路由分组和嵌套功能 (Provides route grouping and nesting functionality)
type RouteGroup interface {
	// RegisterRoute 在组内注册路由 (Register route within group)
	RegisterRoute(method, path string, handler Handler) error
	
	// RegisterMiddleware 注册组级中间件 (Register group-level middleware)
	RegisterMiddleware(middleware Middleware) error
	
	// Group 创建子路由组 (Create sub route group)
	Group(prefix string, middlewares ...Middleware) RouteGroup
}

// Context 统一的请求上下文接口 (Unified request context interface)
// 提供框架无关的请求/响应操作 (Provides framework-agnostic request/response operations)
type Context interface {
	// Request 获取HTTP请求对象 (Get HTTP request object)
	Request() *http.Request
	
	// Response 获取HTTP响应写入器 (Get HTTP response writer)
	Response() http.ResponseWriter
	
	// Param 获取路径参数 (Get path parameter)
	Param(key string) string
	
	// Query 获取查询参数 (Get query parameter)
	Query(key string) string
	
	// Header 获取请求头 (Get request header)
	Header(key string) string
	
	// SetHeader 设置响应头 (Set response header)
	SetHeader(key, value string)
	
	// JSON 返回JSON响应 (Return JSON response)
	JSON(code int, obj interface{}) error
	
	// String 返回字符串响应 (Return string response)
	String(code int, format string, values ...interface{}) error
	
	// Data 返回原始数据响应 (Return raw data response)
	Data(code int, contentType string, data []byte) error
	
	// Set 设置上下文值 (Set context value)
	Set(key string, value interface{})
	
	// Get 获取上下文值 (Get context value)
	Get(key string) (interface{}, bool)
	
	// GetString 获取字符串类型的上下文值 (Get string type context value)
	GetString(key string) string
	
	// GetInt 获取整数类型的上下文值 (Get integer type context value)
	GetInt(key string) int
	
	// GetBool 获取布尔类型的上下文值 (Get boolean type context value)
	GetBool(key string) bool
	
	// Bind 绑定请求数据到结构体 (Bind request data to struct)
	Bind(obj interface{}) error
	
	// ClientIP 获取客户端IP地址 (Get client IP address)
	ClientIP() string
	
	// UserAgent 获取用户代理 (Get user agent)
	UserAgent() string
	
	// Method 获取请求方法 (Get request method)
	Method() string
	
	// Path 获取请求路径 (Get request path)
	Path() string
	
	// FullPath 获取完整路径模式 (Get full path pattern)
	FullPath() string
}

// FrameworkPlugin 框架插件接口 (Framework plugin interface)
// 定义了插件的基本信息和创建方法 (Defines basic plugin information and creation methods)
type FrameworkPlugin interface {
	// Name 返回插件名称 (Return plugin name)
	Name() string
	
	// Version 返回插件版本 (Return plugin version)
	Version() string
	
	// Description 返回插件描述 (Return plugin description)
	Description() string
	
	// DefaultConfig 返回插件的默认配置 (Return plugin default configuration)
	DefaultConfig() interface{}
	
	// CreateFramework 创建框架实例 (Create framework instance)
	// 现在接受服务容器作为参数 (Now accepts service container as parameter)
	CreateFramework(config interface{}, services services.ServiceContainer) (WebFramework, error)
	
	// ValidateConfig 验证配置 (Validate configuration)
	ValidateConfig(config interface{}) error
	
	// GetConfigSchema 获取配置模式 (Get configuration schema)
	// 返回JSON Schema或其他格式的配置描述 (Return JSON Schema or other format config description)
	GetConfigSchema() interface{}
}