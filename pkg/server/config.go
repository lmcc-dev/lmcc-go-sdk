/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: 服务器配置结构定义
 */

package server

import (
	"fmt"
	"time"
)

// ServerConfig 服务器配置结构 (Server configuration structure)
// 定义Web服务器的所有配置选项 (Defines all configuration options for web server)
type ServerConfig struct {
	// Framework 框架名称 (Framework name)
	// 支持的值: gin, echo, fiber (Supported values: gin, echo, fiber)
	Framework string `yaml:"framework" mapstructure:"framework" json:"framework"`
	
	// Host 监听主机地址 (Host address to listen on)
	Host string `yaml:"host" mapstructure:"host" json:"host"`
	
	// Port 监听端口 (Port to listen on)
	Port int `yaml:"port" mapstructure:"port" json:"port"`
	
	// Mode 运行模式 (Running mode)
	// 支持的值: debug, release, test (Supported values: debug, release, test)
	Mode string `yaml:"mode" mapstructure:"mode" json:"mode"`
	
	// ReadTimeout 读取超时时间 (Read timeout duration)
	ReadTimeout time.Duration `yaml:"read-timeout" mapstructure:"read-timeout" json:"read_timeout"`
	
	// WriteTimeout 写入超时时间 (Write timeout duration)
	WriteTimeout time.Duration `yaml:"write-timeout" mapstructure:"write-timeout" json:"write_timeout"`
	
	// IdleTimeout 空闲超时时间 (Idle timeout duration)
	IdleTimeout time.Duration `yaml:"idle-timeout" mapstructure:"idle-timeout" json:"idle_timeout"`
	
	// MaxHeaderBytes 最大请求头字节数 (Maximum request header bytes)
	MaxHeaderBytes int `yaml:"max-header-bytes" mapstructure:"max-header-bytes" json:"max_header_bytes"`
	
	// CORS 跨域配置 (CORS configuration)
	CORS CORSConfig `yaml:"cors" mapstructure:"cors" json:"cors"`
	
	// Middleware 中间件配置 (Middleware configuration)
	Middleware MiddlewareConfig `yaml:"middleware" mapstructure:"middleware" json:"middleware"`
	
	// TLS TLS配置 (TLS configuration)
	TLS TLSConfig `yaml:"tls" mapstructure:"tls" json:"tls"`
	
	// Plugins 插件特定配置 (Plugin-specific configuration)
	// 键为插件名称，值为插件配置 (Key is plugin name, value is plugin configuration)
	Plugins map[string]interface{} `yaml:"plugins" mapstructure:"plugins" json:"plugins"`
	
	// GracefulShutdown 优雅关闭配置 (Graceful shutdown configuration)
	GracefulShutdown GracefulShutdownConfig `yaml:"graceful-shutdown" mapstructure:"graceful-shutdown" json:"graceful_shutdown"`
}

// CORSConfig CORS配置结构 (CORS configuration structure)
type CORSConfig struct {
	// Enabled 是否启用CORS (Whether to enable CORS)
	Enabled bool `yaml:"enabled" mapstructure:"enabled" json:"enabled"`
	
	// AllowOrigins 允许的源 (Allowed origins)
	AllowOrigins []string `yaml:"allow-origins" mapstructure:"allow-origins" json:"allow_origins"`
	
	// AllowMethods 允许的HTTP方法 (Allowed HTTP methods)
	AllowMethods []string `yaml:"allow-methods" mapstructure:"allow-methods" json:"allow_methods"`
	
	// AllowHeaders 允许的请求头 (Allowed request headers)
	AllowHeaders []string `yaml:"allow-headers" mapstructure:"allow-headers" json:"allow_headers"`
	
	// ExposeHeaders 暴露的响应头 (Exposed response headers)
	ExposeHeaders []string `yaml:"expose-headers" mapstructure:"expose-headers" json:"expose_headers"`
	
	// AllowCredentials 是否允许凭证 (Whether to allow credentials)
	AllowCredentials bool `yaml:"allow-credentials" mapstructure:"allow-credentials" json:"allow_credentials"`
	
	// MaxAge 预检请求缓存时间 (Preflight request cache duration)
	MaxAge time.Duration `yaml:"max-age" mapstructure:"max-age" json:"max_age"`
}

// MiddlewareConfig 中间件配置结构 (Middleware configuration structure)
type MiddlewareConfig struct {
	// Logger 日志中间件配置 (Logger middleware configuration)
	Logger LoggerMiddlewareConfig `yaml:"logger" mapstructure:"logger" json:"logger"`
	
	// Recovery 恢复中间件配置 (Recovery middleware configuration)
	Recovery RecoveryMiddlewareConfig `yaml:"recovery" mapstructure:"recovery" json:"recovery"`
	
	// RateLimit 限流中间件配置 (Rate limit middleware configuration)
	RateLimit RateLimitMiddlewareConfig `yaml:"rate-limit" mapstructure:"rate-limit" json:"rate_limit"`
	
	// Auth 认证中间件配置 (Auth middleware configuration)
	Auth AuthMiddlewareConfig `yaml:"auth" mapstructure:"auth" json:"auth"`
}

// LoggerMiddlewareConfig 日志中间件配置 (Logger middleware configuration)
type LoggerMiddlewareConfig struct {
	// Enabled 是否启用 (Whether to enable)
	Enabled bool `yaml:"enabled" mapstructure:"enabled" json:"enabled"`
	
	// SkipPaths 跳过记录的路径 (Paths to skip logging)
	SkipPaths []string `yaml:"skip-paths" mapstructure:"skip-paths" json:"skip_paths"`
	
	// Format 日志格式 (Log format)
	// 支持的值: json, text (Supported values: json, text)
	Format string `yaml:"format" mapstructure:"format" json:"format"`
	
	// IncludeBody 是否包含请求体 (Whether to include request body)
	IncludeBody bool `yaml:"include-body" mapstructure:"include-body" json:"include_body"`
	
	// MaxBodySize 最大请求体大小 (Maximum request body size)
	MaxBodySize int `yaml:"max-body-size" mapstructure:"max-body-size" json:"max_body_size"`
}

// RecoveryMiddlewareConfig 恢复中间件配置 (Recovery middleware configuration)
type RecoveryMiddlewareConfig struct {
	// Enabled 是否启用 (Whether to enable)
	Enabled bool `yaml:"enabled" mapstructure:"enabled" json:"enabled"`
	
	// PrintStack 是否打印堆栈信息 (Whether to print stack trace)
	PrintStack bool `yaml:"print-stack" mapstructure:"print-stack" json:"print_stack"`
	
	// DisableStackAll 是否禁用所有堆栈信息 (Whether to disable all stack traces)
	DisableStackAll bool `yaml:"disable-stack-all" mapstructure:"disable-stack-all" json:"disable_stack_all"`
	
	// DisableColorConsole 是否禁用彩色控制台输出 (Whether to disable color console output)
	DisableColorConsole bool `yaml:"disable-color-console" mapstructure:"disable-color-console" json:"disable_color_console"`
}

// RateLimitMiddlewareConfig 限流中间件配置 (Rate limit middleware configuration)
type RateLimitMiddlewareConfig struct {
	// Enabled 是否启用 (Whether to enable)
	Enabled bool `yaml:"enabled" mapstructure:"enabled" json:"enabled"`
	
	// Rate 每秒请求数 (Requests per second)
	Rate float64 `yaml:"rate" mapstructure:"rate" json:"rate"`
	
	// Burst 突发请求数 (Burst requests)
	Burst int `yaml:"burst" mapstructure:"burst" json:"burst"`
	
	// KeyFunc 键函数类型 (Key function type)
	// 支持的值: ip, user, custom (Supported values: ip, user, custom)
	KeyFunc string `yaml:"key-func" mapstructure:"key-func" json:"key_func"`
}

// AuthMiddlewareConfig 认证中间件配置 (Auth middleware configuration)
type AuthMiddlewareConfig struct {
	// Enabled 是否启用 (Whether to enable)
	Enabled bool `yaml:"enabled" mapstructure:"enabled" json:"enabled"`
	
	// Type 认证类型 (Authentication type)
	// 支持的值: jwt, basic, custom (Supported values: jwt, basic, custom)
	Type string `yaml:"type" mapstructure:"type" json:"type"`
	
	// SkipPaths 跳过认证的路径 (Paths to skip authentication)
	SkipPaths []string `yaml:"skip-paths" mapstructure:"skip-paths" json:"skip_paths"`
	
	// JWT JWT配置 (JWT configuration)
	JWT JWTConfig `yaml:"jwt" mapstructure:"jwt" json:"jwt"`
}

// JWTConfig JWT配置 (JWT configuration)
type JWTConfig struct {
	// Secret JWT密钥 (JWT secret)
	Secret string `yaml:"secret" mapstructure:"secret" json:"secret"`
	
	// Issuer 签发者 (Issuer)
	Issuer string `yaml:"issuer" mapstructure:"issuer" json:"issuer"`
	
	// Audience 受众 (Audience)
	Audience string `yaml:"audience" mapstructure:"audience" json:"audience"`
	
	// ExpirationTime 过期时间 (Expiration time)
	ExpirationTime time.Duration `yaml:"expiration-time" mapstructure:"expiration-time" json:"expiration_time"`
	
	// RefreshTime 刷新时间 (Refresh time)
	RefreshTime time.Duration `yaml:"refresh-time" mapstructure:"refresh-time" json:"refresh_time"`
}

// TLSConfig TLS配置 (TLS configuration)
type TLSConfig struct {
	// Enabled 是否启用TLS (Whether to enable TLS)
	Enabled bool `yaml:"enabled" mapstructure:"enabled" json:"enabled"`
	
	// CertFile 证书文件路径 (Certificate file path)
	CertFile string `yaml:"cert-file" mapstructure:"cert-file" json:"cert_file"`
	
	// KeyFile 私钥文件路径 (Private key file path)
	KeyFile string `yaml:"key-file" mapstructure:"key-file" json:"key_file"`
	
	// AutoTLS 是否启用自动TLS (Whether to enable auto TLS)
	AutoTLS bool `yaml:"auto-tls" mapstructure:"auto-tls" json:"auto_tls"`
	
	// Domains 自动TLS域名列表 (Auto TLS domain list)
	Domains []string `yaml:"domains" mapstructure:"domains" json:"domains"`
}

// GracefulShutdownConfig 优雅关闭配置 (Graceful shutdown configuration)
type GracefulShutdownConfig struct {
	// Enabled 是否启用优雅关闭 (Whether to enable graceful shutdown)
	Enabled bool `yaml:"enabled" mapstructure:"enabled" json:"enabled"`
	
	// Timeout 关闭超时时间 (Shutdown timeout duration)
	Timeout time.Duration `yaml:"timeout" mapstructure:"timeout" json:"timeout"`
	
	// WaitTime 等待时间 (Wait time)
	WaitTime time.Duration `yaml:"wait-time" mapstructure:"wait-time" json:"wait_time"`
}

// DefaultServerConfig 返回默认服务器配置 (Return default server configuration)
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Framework:      "gin",
		Host:           "0.0.0.0",
		Port:           8080,
		Mode:           "debug",
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
		CORS: CORSConfig{
			Enabled:          true,
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			AllowCredentials: false,
			MaxAge:           12 * time.Hour,
		},
		Middleware: MiddlewareConfig{
			Logger: LoggerMiddlewareConfig{
				Enabled:     true,
				Format:      "json",
				IncludeBody: false,
				MaxBodySize: 1024,
			},
			Recovery: RecoveryMiddlewareConfig{
				Enabled:             true,
				PrintStack:          true,
				DisableStackAll:     false,
				DisableColorConsole: false,
			},
			RateLimit: RateLimitMiddlewareConfig{
				Enabled: false,
				Rate:    100,
				Burst:   200,
				KeyFunc: "ip",
			},
			Auth: AuthMiddlewareConfig{
				Enabled: false,
				Type:    "jwt",
				JWT: JWTConfig{
					ExpirationTime: 24 * time.Hour,
					RefreshTime:    7 * 24 * time.Hour,
				},
			},
		},
		TLS: TLSConfig{
			Enabled: false,
			AutoTLS: false,
		},
		GracefulShutdown: GracefulShutdownConfig{
			Enabled:  true,
			Timeout:  30 * time.Second,
			WaitTime: 5 * time.Second,
		},
		Plugins: make(map[string]interface{}),
	}
}

// Validate 验证配置的有效性 (Validate configuration validity)
func (c *ServerConfig) Validate() error {
	// 严格验证关键字段 (Strict validation for critical fields)
	if c.Framework == "" {
		return fmt.Errorf("framework name cannot be empty")
	}
	
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535, got %d", c.Port)
	}
	
	// 自动修正其他字段 (Auto-correct other fields)
	if c.Host == "" {
		c.Host = "0.0.0.0"
	}
	
	if c.Mode == "" {
		c.Mode = "debug"
	}
	
	if c.ReadTimeout <= 0 {
		c.ReadTimeout = 30 * time.Second
	}
	
	if c.WriteTimeout <= 0 {
		c.WriteTimeout = 30 * time.Second
	}
	
	if c.IdleTimeout <= 0 {
		c.IdleTimeout = 60 * time.Second
	}
	
	if c.MaxHeaderBytes <= 0 {
		c.MaxHeaderBytes = 1 << 20 // 1MB
	}
	
	return nil
}

// GetAddress 获取监听地址 (Get listen address)
func (c *ServerConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// IsDebugMode 是否为调试模式 (Whether in debug mode)
func (c *ServerConfig) IsDebugMode() bool {
	return c.Mode == "debug"
}

// IsReleaseMode 是否为发布模式 (Whether in release mode)
func (c *ServerConfig) IsReleaseMode() bool {
	return c.Mode == "release"
}

// IsTestMode 是否为测试模式 (Whether in test mode)
func (c *ServerConfig) IsTestMode() bool {
	return c.Mode == "test"
}