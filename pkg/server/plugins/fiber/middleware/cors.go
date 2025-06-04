/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Fiber CORS中间件实现 (Fiber CORS middleware implementation)
 */

package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
)

// CORSConfig CORS中间件配置 (CORS middleware configuration)
type CORSConfig struct {
	// Enabled 是否启用CORS (Whether to enable CORS)
	Enabled bool `yaml:"enabled" mapstructure:"enabled"`

	// AllowOrigins 允许的源列表 (List of allowed origins)
	AllowOrigins []string `yaml:"allow-origins" mapstructure:"allow-origins"`

	// AllowMethods 允许的HTTP方法 (Allowed HTTP methods)
	AllowMethods []string `yaml:"allow-methods" mapstructure:"allow-methods"`

	// AllowHeaders 允许的请求头 (Allowed request headers)
	AllowHeaders []string `yaml:"allow-headers" mapstructure:"allow-headers"`

	// ExposeHeaders 暴露的响应头 (Exposed response headers)
	ExposeHeaders []string `yaml:"expose-headers" mapstructure:"expose-headers"`

	// AllowCredentials 是否允许凭证 (Whether to allow credentials)
	AllowCredentials bool `yaml:"allow-credentials" mapstructure:"allow-credentials"`

	// MaxAge 预检请求缓存时间（秒） (Preflight request cache time in seconds)
	MaxAge int `yaml:"max-age" mapstructure:"max-age"`
}

// DefaultCORSConfig 默认CORS配置 (Default CORS configuration)
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		Enabled:          true,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD", "PATCH"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{},
		AllowCredentials: false,
		MaxAge:           86400, // 24小时
	}
}

// CORSMiddleware Fiber CORS中间件 (Fiber CORS middleware)
type CORSMiddleware struct {
	config           *CORSConfig              // CORS配置 (CORS configuration)
	serviceContainer services.ServiceContainer // 服务容器 (Service container)
	logger           services.Logger          // 日志服务 (Logger service)
}

// NewCORSMiddleware 创建CORS中间件 (Create CORS middleware)
func NewCORSMiddleware(config *CORSConfig, serviceContainer services.ServiceContainer) *CORSMiddleware {
	if config == nil {
		config = DefaultCORSConfig()
	}

	return &CORSMiddleware{
		config:           config,
		serviceContainer: serviceContainer,
		logger:           serviceContainer.GetLogger(),
	}
}

// Handler 返回Fiber中间件处理器 (Return Fiber middleware handler)
func (m *CORSMiddleware) Handler() fiber.Handler {
	if !m.config.Enabled {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	// 创建Fiber CORS配置 (Create Fiber CORS configuration)
	corsConfig := cors.Config{
		AllowOrigins:     m.getAllowOrigins(),
		AllowMethods:     m.getAllowMethods(),
		AllowHeaders:     m.getAllowHeaders(),
		ExposeHeaders:    m.getExposeHeaders(),
		AllowCredentials: m.config.AllowCredentials,
		MaxAge:           m.config.MaxAge,
	}

	// 记录CORS配置 (Log CORS configuration)
	if m.logger != nil {
		m.logger.Debugw("Fiber CORS middleware configured",
			"allow_origins", m.config.AllowOrigins,
			"allow_methods", m.config.AllowMethods,
			"allow_headers", m.config.AllowHeaders,
			"allow_credentials", m.config.AllowCredentials,
			"max_age", m.config.MaxAge,
		)
	}

	return cors.New(corsConfig)
}

// Process 实现统一中间件接口 (Implement unified middleware interface)
func (m *CORSMiddleware) Process(ctx server.Context, next func() error) error {
	// 对于统一接口，我们需要手动处理CORS逻辑
	if !m.config.Enabled {
		return next()
	}

	// 获取请求信息 (Get request information)
	req := ctx.Request()
	origin := req.Header.Get("Origin")
	method := req.Method

	// 检查是否为预检请求 (Check if it's a preflight request)
	if method == "OPTIONS" {
		// 设置CORS头部 (Set CORS headers)
		m.setCORSHeaders(ctx, origin)
		return ctx.String(204, "")
	}

	// 设置CORS头部 (Set CORS headers)
	m.setCORSHeaders(ctx, origin)

	// 继续处理请求 (Continue processing request)
	return next()
}

// setCORSHeaders 设置CORS响应头 (Set CORS response headers)
func (m *CORSMiddleware) setCORSHeaders(ctx server.Context, origin string) {
	// 检查源是否被允许 (Check if origin is allowed)
	if m.isOriginAllowed(origin) {
		ctx.SetHeader("Access-Control-Allow-Origin", origin)
	} else if len(m.config.AllowOrigins) == 1 && m.config.AllowOrigins[0] == "*" {
		ctx.SetHeader("Access-Control-Allow-Origin", "*")
	}

	// 设置允许的方法 (Set allowed methods)
	if len(m.config.AllowMethods) > 0 {
		methods := ""
		for i, method := range m.config.AllowMethods {
			if i > 0 {
				methods += ", "
			}
			methods += method
		}
		ctx.SetHeader("Access-Control-Allow-Methods", methods)
	}

	// 设置允许的头部 (Set allowed headers)
	if len(m.config.AllowHeaders) > 0 {
		headers := ""
		for i, header := range m.config.AllowHeaders {
			if i > 0 {
				headers += ", "
			}
			headers += header
		}
		ctx.SetHeader("Access-Control-Allow-Headers", headers)
	}

	// 设置暴露的头部 (Set exposed headers)
	if len(m.config.ExposeHeaders) > 0 {
		headers := ""
		for i, header := range m.config.ExposeHeaders {
			if i > 0 {
				headers += ", "
			}
			headers += header
		}
		ctx.SetHeader("Access-Control-Expose-Headers", headers)
	}

	// 设置凭证支持 (Set credentials support)
	if m.config.AllowCredentials {
		ctx.SetHeader("Access-Control-Allow-Credentials", "true")
	}

	// 设置最大缓存时间 (Set max age)
	if m.config.MaxAge > 0 {
		ctx.SetHeader("Access-Control-Max-Age", string(rune(m.config.MaxAge)))
	}
}

// isOriginAllowed 检查源是否被允许 (Check if origin is allowed)
func (m *CORSMiddleware) isOriginAllowed(origin string) bool {
	if origin == "" {
		return false
	}

	for _, allowedOrigin := range m.config.AllowOrigins {
		if allowedOrigin == "*" || allowedOrigin == origin {
			return true
		}
	}

	return false
}

// getAllowOrigins 获取允许的源字符串 (Get allowed origins string)
func (m *CORSMiddleware) getAllowOrigins() string {
	if len(m.config.AllowOrigins) == 0 {
		return "*"
	}
	
	origins := ""
	for i, origin := range m.config.AllowOrigins {
		if i > 0 {
			origins += ","
		}
		origins += origin
	}
	return origins
}

// getAllowMethods 获取允许的方法字符串 (Get allowed methods string)
func (m *CORSMiddleware) getAllowMethods() string {
	if len(m.config.AllowMethods) == 0 {
		return "GET,POST,PUT,DELETE,OPTIONS"
	}
	
	methods := ""
	for i, method := range m.config.AllowMethods {
		if i > 0 {
			methods += ","
		}
		methods += method
	}
	return methods
}

// getAllowHeaders 获取允许的头部字符串 (Get allowed headers string)
func (m *CORSMiddleware) getAllowHeaders() string {
	if len(m.config.AllowHeaders) == 0 {
		return "*"
	}
	
	headers := ""
	for i, header := range m.config.AllowHeaders {
		if i > 0 {
			headers += ","
		}
		headers += header
	}
	return headers
}

// getExposeHeaders 获取暴露的头部字符串 (Get exposed headers string)
func (m *CORSMiddleware) getExposeHeaders() string {
	if len(m.config.ExposeHeaders) == 0 {
		return ""
	}
	
	headers := ""
	for i, header := range m.config.ExposeHeaders {
		if i > 0 {
			headers += ","
		}
		headers += header
	}
	return headers
}

// GetConfig 获取CORS配置 (Get CORS configuration)
func (m *CORSMiddleware) GetConfig() *CORSConfig {
	return m.config
}

// SetConfig 设置CORS配置 (Set CORS configuration)
func (m *CORSMiddleware) SetConfig(config *CORSConfig) {
	if config != nil {
		m.config = config
	}
} 