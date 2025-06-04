/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: CORS中间件接口定义 (CORS middleware interface definitions)
 */

package middleware

import (
	"time"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
)

// CORSConfig CORS中间件配置 (CORS middleware configuration)
type CORSConfig struct {
	// AllowOrigins 允许的源列表 (List of allowed origins)
	// 支持通配符 "*" 和具体域名 (Supports wildcard "*" and specific domains)
	AllowOrigins []string `json:"allow_origins" yaml:"allow_origins"`
	
	// AllowMethods 允许的HTTP方法 (Allowed HTTP methods)
	AllowMethods []string `json:"allow_methods" yaml:"allow_methods"`
	
	// AllowHeaders 允许的请求头 (Allowed request headers)
	AllowHeaders []string `json:"allow_headers" yaml:"allow_headers"`
	
	// ExposeHeaders 暴露的响应头 (Exposed response headers)
	ExposeHeaders []string `json:"expose_headers" yaml:"expose_headers"`
	
	// AllowCredentials 是否允许凭证 (Whether to allow credentials)
	AllowCredentials bool `json:"allow_credentials" yaml:"allow_credentials"`
	
	// MaxAge 预检请求缓存时间 (Preflight request cache time)
	MaxAge time.Duration `json:"max_age" yaml:"max_age"`
	
	// AllowWildcard 是否允许通配符源 (Whether to allow wildcard origins)
	AllowWildcard bool `json:"allow_wildcard" yaml:"allow_wildcard"`
	
	// AllowBrowserExtensions 是否允许浏览器扩展 (Whether to allow browser extensions)
	AllowBrowserExtensions bool `json:"allow_browser_extensions" yaml:"allow_browser_extensions"`
	
	// AllowWebSockets 是否允许WebSocket (Whether to allow WebSocket)
	AllowWebSockets bool `json:"allow_websockets" yaml:"allow_websockets"`
	
	// AllowFiles 是否允许文件协议 (Whether to allow file protocol)
	AllowFiles bool `json:"allow_files" yaml:"allow_files"`
	
	// CustomSchemas 自定义协议列表 (Custom schema list)
	CustomSchemas []string `json:"custom_schemas" yaml:"custom_schemas"`
	
	// OptionsPassthrough 是否透传OPTIONS请求 (Whether to passthrough OPTIONS requests)
	OptionsPassthrough bool `json:"options_passthrough" yaml:"options_passthrough"`
	
	// Debug 是否启用调试模式 (Whether to enable debug mode)
	Debug bool `json:"debug" yaml:"debug"`
}

// DefaultCORSConfig 返回默认CORS配置 (Return default CORS configuration)
func DefaultCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Length", "Content-Type", "Authorization",
			"Accept", "X-Requested-With", "Cache-Control",
		},
		ExposeHeaders:              []string{},
		AllowCredentials:           false,
		MaxAge:                     12 * time.Hour,
		AllowWildcard:              true,
		AllowBrowserExtensions:     false,
		AllowWebSockets:            false,
		AllowFiles:                 false,
		CustomSchemas:              []string{},
		OptionsPassthrough:         false,
		Debug:                      false,
	}
}

// RestrictiveCORSConfig 返回限制性CORS配置 (Return restrictive CORS configuration)
func RestrictiveCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowOrigins: []string{}, // 需要明确指定 (Need to specify explicitly)
		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Authorization",
		},
		ExposeHeaders:              []string{},
		AllowCredentials:           true,
		MaxAge:                     1 * time.Hour,
		AllowWildcard:              false,
		AllowBrowserExtensions:     false,
		AllowWebSockets:            false,
		AllowFiles:                 false,
		CustomSchemas:              []string{},
		OptionsPassthrough:         false,
		Debug:                      false,
	}
}

// DevelopmentCORSConfig 返回开发环境CORS配置 (Return development CORS configuration)
func DevelopmentCORSConfig() *CORSConfig {
	return &CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://localhost:8080",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:8080",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Length", "Content-Type", "Authorization",
			"Accept", "X-Requested-With", "Cache-Control", "X-CSRF-Token",
		},
		ExposeHeaders:              []string{"X-Total-Count"},
		AllowCredentials:           true,
		MaxAge:                     24 * time.Hour,
		AllowWildcard:              false,
		AllowBrowserExtensions:     true,
		AllowWebSockets:            true,
		AllowFiles:                 true,
		CustomSchemas:              []string{"app", "file"},
		OptionsPassthrough:         false,
		Debug:                      true,
	}
}

// CORSFactory CORS中间件工厂接口 (CORS middleware factory interface)
type CORSFactory interface {
	// CreateCORS 创建CORS中间件 (Create CORS middleware)
	CreateCORS(config *CORSConfig) CORSMiddleware
	
	// CreateDefaultCORS 创建默认CORS中间件 (Create default CORS middleware)
	CreateDefaultCORS() CORSMiddleware
	
	// CreateRestrictiveCORS 创建限制性CORS中间件 (Create restrictive CORS middleware)
	CreateRestrictiveCORS() CORSMiddleware
	
	// CreateDevelopmentCORS 创建开发环境CORS中间件 (Create development CORS middleware)
	CreateDevelopmentCORS() CORSMiddleware
}

// CORSMode CORS模式 (CORS mode)
type CORSMode int

const (
	// CORSModePermissive 宽松模式 (Permissive mode)
	// 允许所有源和方法 (Allows all origins and methods)
	CORSModePermissive CORSMode = iota
	
	// CORSModeRestrictive 限制模式 (Restrictive mode)
	// 严格限制源和方法 (Strictly limits origins and methods)
	CORSModeRestrictive
	
	// CORSModeDevelopment 开发模式 (Development mode)
	// 适合开发环境的配置 (Configuration suitable for development)
	CORSModeDevelopment
	
	// CORSModeProduction 生产模式 (Production mode)
	// 适合生产环境的安全配置 (Secure configuration for production)
	CORSModeProduction
)

// String 返回CORS模式字符串 (Return CORS mode string)
func (c CORSMode) String() string {
	switch c {
	case CORSModePermissive:
		return "permissive"
	case CORSModeRestrictive:
		return "restrictive"
	case CORSModeDevelopment:
		return "development"
	case CORSModeProduction:
		return "production"
	default:
		return "unknown"
	}
}

// GetCORSConfigByMode 根据模式获取CORS配置 (Get CORS configuration by mode)
func GetCORSConfigByMode(mode CORSMode) *CORSConfig {
	switch mode {
	case CORSModePermissive:
		return DefaultCORSConfig()
	case CORSModeRestrictive:
		return RestrictiveCORSConfig()
	case CORSModeDevelopment:
		return DevelopmentCORSConfig()
	case CORSModeProduction:
		config := RestrictiveCORSConfig()
		config.Debug = false
		config.AllowBrowserExtensions = false
		config.AllowFiles = false
		return config
	default:
		return DefaultCORSConfig()
	}
}

// CORSHeaders CORS相关的HTTP头常量 (CORS-related HTTP header constants)
var CORSHeaders = struct {
	// 请求头 (Request headers)
	Origin                        string
	AccessControlRequestMethod    string
	AccessControlRequestHeaders   string
	
	// 响应头 (Response headers)
	AccessControlAllowOrigin      string
	AccessControlAllowMethods     string
	AccessControlAllowHeaders     string
	AccessControlAllowCredentials string
	AccessControlExposeHeaders    string
	AccessControlMaxAge           string
	Vary                          string
}{
	// 请求头 (Request headers)
	Origin:                        "Origin",
	AccessControlRequestMethod:    "Access-Control-Request-Method",
	AccessControlRequestHeaders:   "Access-Control-Request-Headers",
	
	// 响应头 (Response headers)
	AccessControlAllowOrigin:      "Access-Control-Allow-Origin",
	AccessControlAllowMethods:     "Access-Control-Allow-Methods",
	AccessControlAllowHeaders:     "Access-Control-Allow-Headers",
	AccessControlAllowCredentials: "Access-Control-Allow-Credentials",
	AccessControlExposeHeaders:    "Access-Control-Expose-Headers",
	AccessControlMaxAge:           "Access-Control-Max-Age",
	Vary:                          "Vary",
}

// IsPreflightRequest 检查是否为预检请求 (Check if it's a preflight request)
func IsPreflightRequest(ctx server.Context) bool {
	return ctx.Method() == "OPTIONS" &&
		ctx.Header(CORSHeaders.Origin) != "" &&
		ctx.Header(CORSHeaders.AccessControlRequestMethod) != ""
}

// IsCORSRequest 检查是否为CORS请求 (Check if it's a CORS request)
func IsCORSRequest(ctx server.Context) bool {
	return ctx.Header(CORSHeaders.Origin) != ""
} 