/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: CORS middleware for Gin web framework plugin system / Gin Web框架插件系统的CORS中间件
 */

package middleware

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
)

var (
	// ErrInvalidConfig 无效配置错误 (Invalid configuration error)
	ErrInvalidConfig = errors.New("invalid configuration")
)

// CORSMiddleware Gin CORS中间件实现 (Gin CORS middleware implementation)
type CORSMiddleware struct {
	config *server.CORSConfig
}

// NewCORSMiddleware 创建新的CORS中间件实例 (Create new CORS middleware instance)
func NewCORSMiddleware(config *server.CORSConfig) server.Middleware {
	return &CORSMiddleware{
		config: config,
	}
}

// isOriginAllowed 检查origin是否被允许 (Check if origin is allowed)
func (m *CORSMiddleware) isOriginAllowed(origin string) bool {
	if len(m.config.AllowOrigins) == 0 {
		return true // 没有限制时允许所有 (Allow all when no restrictions)
	}
	
	for _, allowedOrigin := range m.config.AllowOrigins {
		if allowedOrigin == "*" || allowedOrigin == origin {
			return true
		}
	}
	return false
}

// getMaxAgeString 获取MaxAge的字符串形式 (Get MaxAge as string)
func (m *CORSMiddleware) getMaxAgeString() string {
	maxAge := m.config.MaxAge
	if maxAge <= 0 {
		maxAge = 12 * time.Hour // 默认12小时 (Default 12 hours)
	}
	return strconv.Itoa(int(maxAge.Seconds()))
}

// getDefaultMethods 获取默认允许的HTTP方法 (Get default allowed HTTP methods)
func (m *CORSMiddleware) getDefaultMethods() []string {
	return []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
}

// getDefaultHeaders 获取默认允许的头部 (Get default allowed headers)
func (m *CORSMiddleware) getDefaultHeaders() []string {
	return []string{
		"Origin",
		"Content-Length", 
		"Content-Type",
		"Accept",
		"Accept-Encoding",
		"Authorization",
		"Cache-Control",
		"X-Requested-With",
		"X-CSRF-Token",
	}
}

// handleCORS 处理CORS逻辑 (Handle CORS logic)
func (m *CORSMiddleware) handleCORS(c *gin.Context) {
	origin := c.Request.Header.Get("Origin")
	
	// 处理简单请求和实际请求 (Handle simple requests and actual requests)
	if origin != "" {
		if m.isOriginAllowed(origin) {
			// 设置允许的源 (Set allowed origin)
			if len(m.config.AllowOrigins) == 1 && m.config.AllowOrigins[0] == "*" && !m.config.AllowCredentials {
				c.Header("Access-Control-Allow-Origin", "*")
			} else {
				c.Header("Access-Control-Allow-Origin", origin)
			}
			
			// 设置是否允许凭证 (Set credentials allowance)
			if m.config.AllowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
			
			// 设置暴露的头部 (Set exposed headers)
			if len(m.config.ExposeHeaders) > 0 {
				c.Header("Access-Control-Expose-Headers", strings.Join(m.config.ExposeHeaders, ","))
			}
		} else {
			// Origin不被允许，返回403 (Origin not allowed, return 403)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
	}
	
	// 处理预检请求 (Handle preflight requests)
	if c.Request.Method == "OPTIONS" {
		if origin != "" && m.isOriginAllowed(origin) {
			// 设置允许的方法 (Set allowed methods)
			methods := m.config.AllowMethods
			if len(methods) == 0 {
				methods = m.getDefaultMethods()
			}
			c.Header("Access-Control-Allow-Methods", strings.Join(methods, ", "))
			
			// 设置允许的头部 (Set allowed headers)
			headers := m.config.AllowHeaders
			if len(headers) == 0 {
				headers = m.getDefaultHeaders()
			}
			c.Header("Access-Control-Allow-Headers", strings.Join(headers, ", "))
			
			// 设置预检缓存时间 (Set preflight cache time)
			c.Header("Access-Control-Max-Age", m.getMaxAgeString())
			
			c.AbortWithStatus(http.StatusNoContent)
			return
		} else {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
	}
}

// Process 处理CORS中间件逻辑 (Process CORS middleware logic)
func (m *CORSMiddleware) Process(ctx server.Context, next func() error) error {
	// 检查是否启用CORS (Check if CORS is enabled)
	if m.config == nil || !m.config.Enabled {
		return next()
	}

	// 使用统一接口处理CORS (Handle CORS using unified interface)
	origin := ctx.Request().Header.Get("Origin")
	method := ctx.Request().Method

	// 检查是否为预检请求 (Check if it's a preflight request)
	if method == "OPTIONS" {
		// 设置CORS头部 (Set CORS headers)
		m.setCORSHeaders(ctx, origin)
		return ctx.String(http.StatusNoContent, "")
	}

	// 设置CORS头部 (Set CORS headers)
	m.setCORSHeaders(ctx, origin)

	// 继续处理请求 (Continue processing request)
	return next()
}

// GetGinHandler 返回Gin兼容的处理器 (Return Gin compatible handler)
func (m *CORSMiddleware) GetGinHandler() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// 检查是否启用CORS (Check if CORS is enabled)
		if m.config == nil || !m.config.Enabled {
			c.Next()
			return
		}
		
		// 处理CORS (Handle CORS)
		m.handleCORS(c)
		
		// 如果没有被中止，继续执行 (If not aborted, continue)
		if !c.IsAborted() {
			c.Next()
		}
	})
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
	methods := m.config.AllowMethods
	if len(methods) == 0 {
		methods = m.getDefaultMethods()
	}
	ctx.SetHeader("Access-Control-Allow-Methods", strings.Join(methods, ", "))

	// 设置允许的头部 (Set allowed headers)
	headers := m.config.AllowHeaders
	if len(headers) == 0 {
		headers = m.getDefaultHeaders()
	}
	ctx.SetHeader("Access-Control-Allow-Headers", strings.Join(headers, ", "))

	// 设置暴露的头部 (Set exposed headers)
	if len(m.config.ExposeHeaders) > 0 {
		ctx.SetHeader("Access-Control-Expose-Headers", strings.Join(m.config.ExposeHeaders, ", "))
	}

	// 设置凭证支持 (Set credentials support)
	if m.config.AllowCredentials {
		ctx.SetHeader("Access-Control-Allow-Credentials", "true")
	}

	// 设置最大缓存时间 (Set max age)
	if m.config.MaxAge > 0 {
		ctx.SetHeader("Access-Control-Max-Age", m.getMaxAgeString())
	}
}

// CORSMiddlewareFactory CORS中间件工厂函数 (CORS middleware factory function)
func CORSMiddlewareFactory(config interface{}) (server.Middleware, error) {
	corsConfig, ok := config.(*server.CORSConfig)
	if !ok {
		return nil, ErrInvalidConfig
	}
	return NewCORSMiddleware(corsConfig), nil
} 