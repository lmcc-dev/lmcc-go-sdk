/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	ginpkg "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin"
)

// RecoveryMiddleware Gin Recovery中间件 (Gin Recovery middleware)
type RecoveryMiddleware struct {
	config *server.RecoveryMiddlewareConfig
}

// NewRecoveryMiddleware 创建新的Recovery中间件 (Create new Recovery middleware)
func NewRecoveryMiddleware(config *server.RecoveryMiddlewareConfig) server.Middleware {
	return &RecoveryMiddleware{
		config: config,
	}
}

// Process 处理Recovery中间件逻辑 (Process Recovery middleware logic)
func (m *RecoveryMiddleware) Process(ctx server.Context, next func() error) error {
	// 尝试获取Gin上下文 (Try to get Gin context)
	var ginCtx *gin.Context
	
	// 如果是GinContext类型，获取底层gin.Context (If it's GinContext type, get underlying gin.Context)
	if ginContext, ok := ctx.(*ginpkg.GinContext); ok {
		ginCtx = ginContext.GetGinContext()
	} else {
		// 如果不是Gin上下文，跳过Recovery处理 (If not Gin context, skip Recovery processing)
		return next()
	}

	// 创建Recovery处理器 (Create Recovery handler)
	var recoveryHandler gin.HandlerFunc
	
	if m.config != nil && m.config.PrintStack {
		// 使用带堆栈跟踪的Recovery (Use Recovery with stack trace)
		recoveryHandler = gin.RecoveryWithWriter(gin.DefaultWriter, m.createRecoveryFunc())
	} else {
		// 使用默认Recovery (Use default Recovery)
		recoveryHandler = gin.Recovery()
	}

	// 执行Recovery中间件 (Execute Recovery middleware)
	recoveryHandler(ginCtx)

	// 如果请求被中止，返回错误 (If request is aborted, return error)
	if ginCtx.IsAborted() {
		return errors.New("request was aborted by recovery middleware")
	}

	// 继续执行下一个中间件 (Continue to next middleware)
	return next()
}

// createRecoveryFunc 创建自定义Recovery函数 (Create custom Recovery function)
func (m *RecoveryMiddleware) createRecoveryFunc() gin.RecoveryFunc {
	return func(c *gin.Context, err interface{}) {
		// 设置HTTP状态码 (Set HTTP status code)
		c.AbortWithStatus(http.StatusInternalServerError)
		
		// 记录错误信息 (Log error information)
		if m.config.PrintStack {
			fmt.Printf("Recovery middleware caught panic: %v\n", err)
		}
		
		// 返回错误响应 (Return error response)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": "An unexpected error occurred",
		})
	}
}

// RecoveryMiddlewareFactory Recovery中间件工厂 (Recovery middleware factory)
func RecoveryMiddlewareFactory(config interface{}) (server.Middleware, error) {
	recoveryConfig, ok := config.(*server.RecoveryMiddlewareConfig)
	if !ok {
		return nil, ErrInvalidConfig
	}
	return NewRecoveryMiddleware(recoveryConfig), nil
} 