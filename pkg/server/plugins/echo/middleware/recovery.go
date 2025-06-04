/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Echo Recovery中间件实现 (Echo Recovery middleware implementation)
 */

package middleware

import (
	"fmt"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
)

// RecoveryConfig Recovery中间件配置 (Recovery middleware configuration)
type RecoveryConfig struct {
	// Enabled 是否启用Recovery中间件 (Whether to enable recovery middleware)
	Enabled bool `yaml:"enabled" mapstructure:"enabled"`

	// PrintStack 是否打印堆栈信息 (Whether to print stack trace)
	PrintStack bool `yaml:"print-stack" mapstructure:"print-stack"`

	// StackSize 堆栈大小限制 (Stack size limit)
	StackSize int `yaml:"stack-size" mapstructure:"stack-size"`

	// DisableStackAll 是否禁用所有goroutine的堆栈 (Whether to disable stack of all goroutines)
	DisableStackAll bool `yaml:"disable-stack-all" mapstructure:"disable-stack-all"`

	// DisableColorOutput 是否禁用颜色输出 (Whether to disable colored output)
	DisableColorOutput bool `yaml:"disable-color-output" mapstructure:"disable-color-output"`
}

// DefaultRecoveryConfig 默认Recovery配置 (Default Recovery configuration)
func DefaultRecoveryConfig() *RecoveryConfig {
	return &RecoveryConfig{
		Enabled:            true,
		PrintStack:         true,
		StackSize:          4 << 10, // 4KB
		DisableStackAll:    false,
		DisableColorOutput: false,
	}
}

// RecoveryMiddleware Echo Recovery中间件 (Echo Recovery middleware)
type RecoveryMiddleware struct {
	config           *RecoveryConfig          // Recovery配置 (Recovery configuration)
	serviceContainer services.ServiceContainer // 服务容器 (Service container)
	logger           services.Logger          // 日志服务 (Logger service)
	errorHandler     services.ErrorHandler    // 错误处理服务 (Error handler service)
}

// NewRecoveryMiddleware 创建Recovery中间件 (Create Recovery middleware)
func NewRecoveryMiddleware(config *RecoveryConfig, serviceContainer services.ServiceContainer) *RecoveryMiddleware {
	if config == nil {
		config = DefaultRecoveryConfig()
	}

	return &RecoveryMiddleware{
		config:           config,
		serviceContainer: serviceContainer,
		logger:           serviceContainer.GetLogger(),
		errorHandler:     serviceContainer.GetErrorHandler(),
	}
}

// Handler 返回Echo中间件处理器 (Return Echo middleware handler)
func (m *RecoveryMiddleware) Handler() echo.MiddlewareFunc {
	if !m.config.Enabled {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return next
		}
	}

	// 创建Echo Recovery配置 (Create Echo Recovery configuration)
	recoveryConfig := middleware.RecoverConfig{
		StackSize:                 m.config.StackSize,
		DisableStackAll:           m.config.DisableStackAll,
		DisablePrintStack:         !m.config.PrintStack,
		DisableErrorHandler:       false,
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			// 使用我们的日志服务记录错误 (Use our logger service to log errors)
			if m.logger != nil {
				m.logger.Errorw("Panic recovered",
					"error", err.Error(),
					"method", c.Request().Method,
					"uri", c.Request().RequestURI,
					"client_ip", c.RealIP(),
					"user_agent", c.Request().UserAgent(),
					"framework", "echo",
				)

				if m.config.PrintStack && len(stack) > 0 {
					m.logger.Errorw("Stack trace", "stack", string(stack))
				}
			}
			return nil
		},
	}

	// 记录Recovery配置 (Log Recovery configuration)
	if m.logger != nil {
		m.logger.Debugw("Echo Recovery middleware configured",
			"enabled", m.config.Enabled,
			"print_stack", m.config.PrintStack,
			"stack_size", m.config.StackSize,
			"disable_stack_all", m.config.DisableStackAll,
		)
	}

	return middleware.RecoverWithConfig(recoveryConfig)
}

// Process 实现统一中间件接口 (Implement unified middleware interface)
func (m *RecoveryMiddleware) Process(ctx server.Context, next func() error) error {
	if !m.config.Enabled {
		return next()
	}

	// 使用defer recover来捕获panic (Use defer recover to catch panic)
	defer func() {
		if err := recover(); err != nil {
			// 获取堆栈信息 (Get stack trace)
			var stack []byte
			if m.config.PrintStack {
				stack = make([]byte, m.config.StackSize)
				length := runtime.Stack(stack, !m.config.DisableStackAll)
				stack = stack[:length]
			}

			// 记录panic信息 (Log panic information)
			if m.logger != nil {
				req := ctx.Request()
				m.logger.Errorw("Panic recovered",
					"error", fmt.Sprintf("%v", err),
					"method", req.Method,
					"uri", req.RequestURI,
					"client_ip", ctx.ClientIP(),
					"user_agent", req.UserAgent(),
					"framework", "echo",
				)

				if m.config.PrintStack && len(stack) > 0 {
					m.logger.Errorw("Stack trace", "stack", string(stack))
				}
			}

			// 返回500错误 (Return 500 error)
			if err := ctx.JSON(500, map[string]interface{}{
				"error":   "Internal Server Error",
				"message": "An unexpected error occurred",
			}); err != nil && m.logger != nil {
				m.logger.Errorw("Failed to send error response", "error", err)
			}
		}
	}()

	// 执行下一个处理器 (Execute next handler)
	return next()
}

// GetConfig 获取Recovery配置 (Get Recovery configuration)
func (m *RecoveryMiddleware) GetConfig() *RecoveryConfig {
	return m.config
}

// SetConfig 设置Recovery配置 (Set Recovery configuration)
func (m *RecoveryMiddleware) SetConfig(config *RecoveryConfig) {
	if config != nil {
		m.config = config
	}
} 