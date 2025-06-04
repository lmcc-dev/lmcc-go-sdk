/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Echo Logger中间件实现 (Echo Logger middleware implementation)
 */

package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

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
}

// DefaultLoggerConfig 默认Logger配置 (Default Logger configuration)
func DefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Enabled:      true,
		Format:       "${time_rfc3339} ${status} ${method} ${uri} ${latency_human} ${bytes_in}/${bytes_out}\n",
		SkipPaths:    []string{"/health", "/metrics"},
		EnableColors: false,
		TimeFormat:   time.RFC3339,
	}
}

// LoggerMiddleware Echo Logger中间件 (Echo Logger middleware)
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

// Handler 返回Echo中间件处理器 (Return Echo middleware handler)
func (m *LoggerMiddleware) Handler() echo.MiddlewareFunc {
	if !m.config.Enabled {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return next
		}
	}

	// 创建Echo Logger配置 (Create Echo Logger configuration)
	loggerConfig := middleware.LoggerConfig{
		Format: m.config.Format,
		Skipper: func(c echo.Context) bool {
			path := c.Request().URL.Path
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
		m.logger.Debugw("Echo Logger middleware configured",
			"enabled", m.config.Enabled,
			"format", m.config.Format,
			"skip_paths", m.config.SkipPaths,
			"enable_colors", m.config.EnableColors,
		)
	}

	return middleware.LoggerWithConfig(loggerConfig)
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
			"framework", "echo",
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