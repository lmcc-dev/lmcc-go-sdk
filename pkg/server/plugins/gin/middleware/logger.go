/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Gin日志中间件实现 (Gin logger middleware implementation)
 */

package middleware

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	ginpkg "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
)

// GinLoggerMiddleware Gin日志中间件 (Gin logger middleware)
type GinLoggerMiddleware struct {
	config      *LoggerConfig
	skipPaths   map[string]bool
	services    services.ServiceContainer
	includeBody bool
	maxBodySize int
}

// LoggerConfig 日志中间件配置 (Logger middleware configuration)
type LoggerConfig struct {
	// SkipPaths 跳过记录的路径列表 (List of paths to skip logging)
	SkipPaths []string `json:"skip_paths" yaml:"skip_paths"`
	
	// Format 日志格式 (Log format)
	Format string `json:"format" yaml:"format"`
	
	// IncludeBody 是否包含请求体 (Whether to include request body)
	IncludeBody bool `json:"include_body" yaml:"include_body"`
	
	// MaxBodySize 最大请求体大小 (Maximum request body size)
	MaxBodySize int `json:"max_body_size" yaml:"max_body_size"`
}

// DefaultLoggerConfig 返回默认日志配置 (Return default logger configuration)
func DefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		SkipPaths:   []string{"/health", "/metrics"},
		Format:      "json",
		IncludeBody: false,
		MaxBodySize: 1024, // 1KB
	}
}

// NewGinLoggerMiddleware 创建Gin日志中间件 (Create Gin logger middleware)
// 保留向后兼容性 (Maintain backward compatibility)
func NewGinLoggerMiddleware(config *LoggerConfig) server.Middleware {
	return NewGinLoggerMiddlewareWithServices(config, nil)
}

// NewGinLoggerMiddlewareWithServices 创建带服务容器的Gin日志中间件 (Create Gin logger middleware with service container)
func NewGinLoggerMiddlewareWithServices(config *LoggerConfig, serviceContainer services.ServiceContainer) *GinLoggerMiddleware {
	if config == nil {
		config = DefaultLoggerConfig()
	}
	
	// 如果没有提供服务容器，创建默认的 (If no service container provided, create default)
	if serviceContainer == nil {
		serviceContainer = services.NewServiceContainerWithDefaults()
	}
	
	// 创建跳过路径映射 (Create skip paths map)
	skipPaths := make(map[string]bool)
	for _, path := range config.SkipPaths {
		skipPaths[path] = true
	}
	
	return &GinLoggerMiddleware{
		config:      config,
		skipPaths:   skipPaths,
		services:    serviceContainer,
		includeBody: config.IncludeBody,
		maxBodySize: config.MaxBodySize,
	}
}

// Process 实现Middleware接口的Process方法 (Implement Process method of Middleware interface)
func (m *GinLoggerMiddleware) Process(ctx server.Context, next func() error) error {
	// 检查是否跳过此路径 (Check if skip this path)
	if m.skipPaths[ctx.Path()] {
		return next()
	}
	
	// 记录开始时间 (Record start time)
	start := time.Now()
	
	// 记录请求体 (Record request body)
	var requestBody []byte
	if m.includeBody && ctx.Method() != "GET" && ctx.Method() != "HEAD" {
		if body, err := m.readRequestBody(ctx); err == nil {
			requestBody = body
		}
	}
	
	// 执行下一个处理器 (Execute next handler)
	err := next()
	
	// 计算处理时间 (Calculate processing time)
	latency := time.Since(start)
	
	// 获取响应状态码 (Get response status code)
	statusCode := 200 // 默认状态码 (Default status code)
	if ginCtx, ok := m.getGinContext(ctx); ok {
		statusCode = ginCtx.Writer.Status()
	}
	
	// 构建日志字段 (Build log fields)
	fields := map[string]interface{}{
		"method":     ctx.Method(),
		"path":       ctx.Path(),
		"full_path":  ctx.FullPath(),
		"status":     statusCode,
		"latency":    latency.String(),
		"latency_ms": latency.Milliseconds(),
		"client_ip":  ctx.ClientIP(),
		"user_agent": ctx.UserAgent(),
		"timestamp":  start.Format(time.RFC3339),
	}
	
	// 添加请求体 (Add request body)
	if len(requestBody) > 0 {
		fields["request_body"] = string(requestBody)
	}
	
	// 添加错误信息 (Add error information)
	if err != nil {
		fields["error"] = err.Error()
	}
	
	// 根据状态码确定日志级别 (Determine log level based on status code)
	logLevel := m.getLogLevel(statusCode, err)
	
	// 记录日志 (Log the request)
	message := fmt.Sprintf("%s %s %d %s", ctx.Method(), ctx.Path(), statusCode, latency)
	
	// 使用服务容器的日志器 (Use service container's logger)
	logger := m.services.GetLogger()
	
	switch logLevel {
	case "error":
		logger.Errorw(message, m.fieldsToKeyValue(fields)...)
	case "warn":
		logger.Warnw(message, m.fieldsToKeyValue(fields)...)
	case "info":
		logger.Infow(message, m.fieldsToKeyValue(fields)...)
	default:
		logger.Debugw(message, m.fieldsToKeyValue(fields)...)
	}
	
	return err
}

// readRequestBody 读取请求体 (Read request body)
func (m *GinLoggerMiddleware) readRequestBody(ctx server.Context) ([]byte, error) {
	request := ctx.Request()
	if request.Body == nil {
		return nil, nil
	}
	
	// 读取请求体 (Read request body)
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}
	
	// 限制请求体大小 (Limit request body size)
	if len(body) > m.maxBodySize {
		body = body[:m.maxBodySize]
	}
	
	// 恢复请求体 (Restore request body)
	request.Body = io.NopCloser(bytes.NewBuffer(body))
	
	return body, nil
}

// getGinContext 获取Gin上下文 (Get Gin context)
func (m *GinLoggerMiddleware) getGinContext(ctx server.Context) (*gin.Context, bool) {
	// 如果是我们的适配器，尝试获取原生上下文 (If it's our adapter, try to get native context)
	if adapter, ok := ctx.(*ginpkg.GinContext); ok {
		return adapter.GetGinContext(), true
	}
	
	return nil, false
}

// getLogLevel 根据状态码和错误确定日志级别 (Determine log level based on status code and error)
func (m *GinLoggerMiddleware) getLogLevel(statusCode int, err error) string {
	if err != nil {
		return "error"
	}
	
	switch {
	case statusCode >= 500:
		return "error"
	case statusCode >= 400:
		return "warn"
	case statusCode >= 300:
		return "info"
	default:
		return "debug"
	}
}

// fieldsToKeyValue 将字段映射转换为键值对 (Convert fields map to key-value pairs)
func (m *GinLoggerMiddleware) fieldsToKeyValue(fields map[string]interface{}) []interface{} {
	result := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		result = append(result, k, v)
	}
	return result
}

// Apply 应用中间件到Gin引擎 (Apply middleware to Gin engine)
// 这是一个便利方法，用于直接应用到Gin引擎 (This is a convenience method for direct application to Gin engine)
func (m *GinLoggerMiddleware) Apply(engine *gin.Engine) {
	engine.Use(m.CreateGinHandler())
}

// CreateGinHandler 创建Gin处理器 (Create Gin handler)
func (m *GinLoggerMiddleware) CreateGinHandler() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		// 创建适配器 (Create adapter)
		ctx := ginpkg.NewGinContext(ginCtx)
		
		// 执行中间件逻辑 (Execute middleware logic)
		err := m.Process(ctx, func() error {
			ginCtx.Next()
			return nil
		})
		
		// 如果有错误，设置到Gin上下文 (If there's an error, set it to Gin context)
		if err != nil {
			ginCtx.Error(err)
		}
	}
}