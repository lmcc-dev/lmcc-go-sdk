/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Fiber Logger中间件实现 (Fiber Logger middleware implementation)
 */

package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
)

// LoggerConfig Logger中间件配置 (Logger middleware configuration)
type LoggerConfig struct {
	// Enabled 是否启用日志中间件 (Whether to enable logger middleware)
	Enabled bool `yaml:"enabled" mapstructure:"enabled"`

	// Format 日志格式 (Log format)
	Format string `yaml:"format" mapstructure:"format"`

	// SkipPaths 跳过记录的路径 (Paths to skip logging)
	SkipPaths []string `yaml:"skip-paths" mapstructure:"skip-paths"`

	// EnableColors 是否启用颜色输出 (Whether to enable colored output)
	EnableColors bool `yaml:"enable-colors" mapstructure:"enable-colors"`

	// TimeFormat 时间格式 (Time format)
	TimeFormat string `yaml:"time-format" mapstructure:"time-format"`

	// TimeZone 时区 (Time zone)
	TimeZone string `yaml:"time-zone" mapstructure:"time-zone"`
}

// DefaultLoggerConfig 默认Logger配置 (Default Logger configuration)
func DefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Enabled:      true,
		Format:       "[${time}] ${status} - ${method} ${path} ${latency}\n",
		SkipPaths:    []string{"/health", "/metrics"},
		EnableColors: false,
		TimeFormat:   "15:04:05",
		TimeZone:     "Local",
	}
}

// LoggerMiddleware Fiber Logger中间件 (Fiber Logger middleware)
type LoggerMiddleware struct {
	config           *LoggerConfig            // Logger配置 (Logger configuration)
	serviceContainer services.ServiceContainer // 服务容器 (Service container)
	logger           services.Logger          // 日志服务 (Logger service)
}

// NewLoggerMiddleware 创建Logger中间件 (Create Logger middleware)
func NewLoggerMiddleware(config *LoggerConfig, serviceContainer services.ServiceContainer) *LoggerMiddleware {
	if config == nil {
		config = DefaultLoggerConfig()
	}

	return &LoggerMiddleware{
		config:           config,
		serviceContainer: serviceContainer,
		logger:           serviceContainer.GetLogger(),
	}
}

// Handler 返回Fiber中间件处理器 (Return Fiber middleware handler)
func (m *LoggerMiddleware) Handler() fiber.Handler {
	if !m.config.Enabled {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	// 创建Fiber Logger配置 (Create Fiber Logger configuration)
	loggerConfig := logger.Config{
		Format:     m.config.Format,
		TimeFormat: m.config.TimeFormat,
		TimeZone:   m.config.TimeZone,
		Next: func(c *fiber.Ctx) bool {
			// 检查是否跳过此路径 (Check if this path should be skipped)
			path := c.Path()
			for _, skipPath := range m.config.SkipPaths {
				if path == skipPath {
					return true
				}
			}
			return false
		},
	}

	// 记录Logger配置 (Log Logger configuration)
	if m.logger != nil {
		m.logger.Debugw("Fiber Logger middleware configured",
			"enabled", m.config.Enabled,
			"format", m.config.Format,
			"skip_paths", m.config.SkipPaths,
			"enable_colors", m.config.EnableColors,
			"time_format", m.config.TimeFormat,
		)
	}

	return logger.New(loggerConfig)
}

// Process 实现统一中间件接口 (Implement unified middleware interface)
func (m *LoggerMiddleware) Process(ctx server.Context, next func() error) error {
	if !m.config.Enabled {
		return next()
	}

	// 检查是否跳过此路径 (Check if this path should be skipped)
	path := ctx.Request().URL.Path
	for _, skipPath := range m.config.SkipPaths {
		if path == skipPath {
			return next()
		}
	}

	// 记录请求开始时间 (Record request start time)
	start := time.Now()

	// 执行下一个处理器 (Execute next handler)
	err := next()

	// 计算请求处理时间 (Calculate request processing time)
	latency := time.Since(start)

	// 获取请求信息 (Get request information)
	req := ctx.Request()
	method := req.Method
	uri := req.RequestURI
	userAgent := req.UserAgent()
	clientIP := ctx.ClientIP()

	// 获取响应状态码 (Get response status code)
	status := 200 // 默认状态码
	if err != nil {
		status = 500 // 错误状态码
	}

	// 记录请求日志 (Log request)
	if m.logger != nil {
		m.logger.Infow("HTTP Request",
			"method", method,
			"uri", uri,
			"status", status,
			"latency", latency.String(),
			"client_ip", clientIP,
			"user_agent", userAgent,
			"framework", "fiber",
		)
	}

	return err
}

// GetConfig 获取Logger配置 (Get Logger configuration)
func (m *LoggerMiddleware) GetConfig() *LoggerConfig {
	return m.config
}

// SetConfig 设置Logger配置 (Set Logger configuration)
func (m *LoggerMiddleware) SetConfig(config *LoggerConfig) {
	if config != nil {
		m.config = config
	}
} 