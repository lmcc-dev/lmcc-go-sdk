/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: 统一中间件接口定义
 */

package middleware

import (
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
)

// LoggerMiddleware 日志中间件接口 (Logger middleware interface)
// 定义请求日志记录的标准接口 (Defines standard interface for request logging)
type LoggerMiddleware interface {
	server.Middleware
	
	// SetSkipPaths 设置跳过记录的路径 (Set paths to skip logging)
	SetSkipPaths(paths []string)
	
	// SetFormat 设置日志格式 (Set log format)
	SetFormat(format string)
	
	// SetIncludeBody 设置是否包含请求体 (Set whether to include request body)
	SetIncludeBody(include bool)
}

// RecoveryMiddleware 恢复中间件接口 (Recovery middleware interface)
// 定义panic恢复处理的标准接口 (Defines standard interface for panic recovery)
type RecoveryMiddleware interface {
	server.Middleware
	
	// SetPrintStack 设置是否打印堆栈信息 (Set whether to print stack trace)
	SetPrintStack(print bool)
	
	// SetDisableStackAll 设置是否禁用所有堆栈信息 (Set whether to disable all stack traces)
	SetDisableStackAll(disable bool)
	
	// SetCustomRecoveryHandler 设置自定义恢复处理器 (Set custom recovery handler)
	SetCustomRecoveryHandler(handler func(server.Context, interface{}))
}

// CORSMiddleware CORS中间件接口 (CORS middleware interface)
// 定义跨域资源共享处理的标准接口 (Defines standard interface for CORS handling)
type CORSMiddleware interface {
	server.Middleware
	
	// SetAllowOrigins 设置允许的源 (Set allowed origins)
	SetAllowOrigins(origins []string)
	
	// SetAllowMethods 设置允许的HTTP方法 (Set allowed HTTP methods)
	SetAllowMethods(methods []string)
	
	// SetAllowHeaders 设置允许的请求头 (Set allowed request headers)
	SetAllowHeaders(headers []string)
	
	// SetAllowCredentials 设置是否允许凭证 (Set whether to allow credentials)
	SetAllowCredentials(allow bool)
}

// RateLimitMiddleware 限流中间件接口 (Rate limit middleware interface)
// 定义请求限流处理的标准接口 (Defines standard interface for request rate limiting)
type RateLimitMiddleware interface {
	server.Middleware
	
	// SetRate 设置每秒请求数 (Set requests per second)
	SetRate(rate float64)
	
	// SetBurst 设置突发请求数 (Set burst requests)
	SetBurst(burst int)
	
	// SetKeyFunc 设置键函数 (Set key function)
	SetKeyFunc(keyFunc func(server.Context) string)
}

// AuthMiddleware 认证中间件接口 (Auth middleware interface)
// 定义身份认证处理的标准接口 (Defines standard interface for authentication)
type AuthMiddleware interface {
	server.Middleware
	
	// SetSkipPaths 设置跳过认证的路径 (Set paths to skip authentication)
	SetSkipPaths(paths []string)
	
	// SetAuthFunc 设置认证函数 (Set authentication function)
	SetAuthFunc(authFunc func(server.Context) (interface{}, error))
	
	// SetUnauthorizedHandler 设置未授权处理器 (Set unauthorized handler)
	SetUnauthorizedHandler(handler func(server.Context))
}

// CompressionMiddleware 压缩中间件接口 (Compression middleware interface)
// 定义响应压缩处理的标准接口 (Defines standard interface for response compression)
type CompressionMiddleware interface {
	server.Middleware
	
	// SetLevel 设置压缩级别 (Set compression level)
	SetLevel(level int)
	
	// SetMinLength 设置最小压缩长度 (Set minimum compression length)
	SetMinLength(length int)
	
	// SetExcludedExtensions 设置排除的文件扩展名 (Set excluded file extensions)
	SetExcludedExtensions(extensions []string)
}

// SecurityMiddleware 安全中间件接口 (Security middleware interface)
// 定义安全头处理的标准接口 (Defines standard interface for security headers)
type SecurityMiddleware interface {
	server.Middleware
	
	// SetXSSProtection 设置XSS保护 (Set XSS protection)
	SetXSSProtection(enable bool)
	
	// SetContentTypeNosniff 设置内容类型嗅探保护 (Set content type nosniff protection)
	SetContentTypeNosniff(enable bool)
	
	// SetFrameOptions 设置框架选项 (Set frame options)
	SetFrameOptions(options string)
	
	// SetHSTSMaxAge 设置HSTS最大年龄 (Set HSTS max age)
	SetHSTSMaxAge(maxAge int)
}

// MetricsMiddleware 指标中间件接口 (Metrics middleware interface)
// 定义请求指标收集的标准接口 (Defines standard interface for request metrics collection)
type MetricsMiddleware interface {
	server.Middleware
	
	// SetMetricsPath 设置指标路径 (Set metrics path)
	SetMetricsPath(path string)
	
	// SetSkipPaths 设置跳过收集指标的路径 (Set paths to skip metrics collection)
	SetSkipPaths(paths []string)
	
	// SetCustomLabels 设置自定义标签 (Set custom labels)
	SetCustomLabels(labels map[string]string)
}

// MiddlewareFactory 中间件工厂接口 (Middleware factory interface)
// 定义创建各种中间件的标准接口 (Defines standard interface for creating various middleware)
type MiddlewareFactory interface {
	// CreateLogger 创建日志中间件 (Create logger middleware)
	CreateLogger(config server.LoggerMiddlewareConfig) LoggerMiddleware
	
	// CreateRecovery 创建恢复中间件 (Create recovery middleware)
	CreateRecovery(config server.RecoveryMiddlewareConfig) RecoveryMiddleware
	
	// CreateCORS 创建CORS中间件 (Create CORS middleware)
	CreateCORS(config server.CORSConfig) CORSMiddleware
	
	// CreateRateLimit 创建限流中间件 (Create rate limit middleware)
	CreateRateLimit(config server.RateLimitMiddlewareConfig) RateLimitMiddleware
	
	// CreateAuth 创建认证中间件 (Create auth middleware)
	CreateAuth(config server.AuthMiddlewareConfig) AuthMiddleware
	
	// CreateCompression 创建压缩中间件 (Create compression middleware)
	CreateCompression() CompressionMiddleware
	
	// CreateSecurity 创建安全中间件 (Create security middleware)
	CreateSecurity() SecurityMiddleware
	
	// CreateMetrics 创建指标中间件 (Create metrics middleware)
	CreateMetrics() MetricsMiddleware
}

// MiddlewareChain 中间件链 (Middleware chain)
// 管理中间件的执行顺序 (Manages middleware execution order)
type MiddlewareChain struct {
	middlewares []server.Middleware
}

// NewMiddlewareChain 创建中间件链 (Create middleware chain)
func NewMiddlewareChain() *MiddlewareChain {
	return &MiddlewareChain{
		middlewares: make([]server.Middleware, 0),
	}
}

// Add 添加中间件 (Add middleware)
func (mc *MiddlewareChain) Add(middleware server.Middleware) *MiddlewareChain {
	mc.middlewares = append(mc.middlewares, middleware)
	return mc
}

// Execute 执行中间件链 (Execute middleware chain)
func (mc *MiddlewareChain) Execute(ctx server.Context, handler server.Handler) error {
	if len(mc.middlewares) == 0 {
		return handler.Handle(ctx)
	}
	
	// 创建执行链 (Create execution chain)
	var execute func(int) error
	execute = func(index int) error {
		if index >= len(mc.middlewares) {
			return handler.Handle(ctx)
		}
		
		return mc.middlewares[index].Process(ctx, func() error {
			return execute(index + 1)
		})
	}
	
	return execute(0)
}

// GetMiddlewares 获取所有中间件 (Get all middleware)
func (mc *MiddlewareChain) GetMiddlewares() []server.Middleware {
	return mc.middlewares
}

// Clear 清空中间件链 (Clear middleware chain)
func (mc *MiddlewareChain) Clear() {
	mc.middlewares = mc.middlewares[:0]
}

// Count 获取中间件数量 (Get middleware count)
func (mc *MiddlewareChain) Count() int {
	return len(mc.middlewares)
}