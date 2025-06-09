/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Panic恢复中间件接口定义 (Panic recovery middleware interface definitions)
 */

package middleware

import (
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
)

// RecoveryConfig 恢复中间件配置 (Recovery middleware configuration)
type RecoveryConfig struct {
	// EnableStackTrace 是否启用堆栈跟踪 (Whether to enable stack trace)
	EnableStackTrace bool `json:"enable_stack_trace" yaml:"enable_stack_trace"`
	
	// LogStackTrace 是否记录堆栈跟踪到日志 (Whether to log stack trace)
	LogStackTrace bool `json:"log_stack_trace" yaml:"log_stack_trace"`
	
	// ResponseWithStack 是否在响应中包含堆栈信息 (Whether to include stack in response)
	ResponseWithStack bool `json:"response_with_stack" yaml:"response_with_stack"`
	
	// PrintStack 是否打印堆栈到控制台 (Whether to print stack to console)
	PrintStack bool `json:"print_stack" yaml:"print_stack"`
	
	// DisableStackAll 是否禁用所有堆栈信息 (Whether to disable all stack traces)
	DisableStackAll bool `json:"disable_stack_all" yaml:"disable_stack_all"`
	
	// CustomRecoveryHandler 自定义恢复处理器 (Custom recovery handler)
	CustomRecoveryHandler func(server.Context, interface{}) `json:"-" yaml:"-"`
	
	// ErrorResponseHandler 错误响应处理器 (Error response handler)
	ErrorResponseHandler func(server.Context, interface{}, []byte) `json:"-" yaml:"-"`
}

// DefaultRecoveryConfig 返回默认恢复配置 (Return default recovery configuration)
func DefaultRecoveryConfig() *RecoveryConfig {
	return &RecoveryConfig{
		EnableStackTrace:      true,
		LogStackTrace:         true,
		ResponseWithStack:     false, // 生产环境不暴露堆栈信息 (Don't expose stack in production)
		PrintStack:            false,
		DisableStackAll:       false,
		CustomRecoveryHandler: nil,
		ErrorResponseHandler:  nil,
	}
}

// PanicInfo Panic信息 (Panic information)
type PanicInfo struct {
	// Value panic的值 (Panic value)
	Value interface{} `json:"value"`
	
	// Stack 堆栈跟踪 (Stack trace)
	Stack []byte `json:"stack,omitempty"`
	
	// Message 错误消息 (Error message)
	Message string `json:"message"`
	
	// RequestInfo 请求信息 (Request information)
	RequestInfo RequestInfo `json:"request_info"`
}

// RequestInfo 请求信息 (Request information)
type RequestInfo struct {
	// Method HTTP方法 (HTTP method)
	Method string `json:"method"`
	
	// Path 请求路径 (Request path)
	Path string `json:"path"`
	
	// ClientIP 客户端IP (Client IP)
	ClientIP string `json:"client_ip"`
	
	// UserAgent 用户代理 (User agent)
	UserAgent string `json:"user_agent"`
	
	// Headers 请求头 (Request headers)
	Headers map[string]string `json:"headers,omitempty"`
}

// PanicRecovery Panic恢复中间件接口 (Panic recovery middleware interface)
type PanicRecovery interface {
	server.Middleware
	
	// SetConfig 设置配置 (Set configuration)
	SetConfig(config *RecoveryConfig)
	
	// GetConfig 获取配置 (Get configuration)
	GetConfig() *RecoveryConfig
	
	// SetPrintStack 设置是否打印堆栈信息 (Set whether to print stack trace)
	SetPrintStack(print bool)
	
	// SetDisableStackAll 设置是否禁用所有堆栈信息 (Set whether to disable all stack traces)
	SetDisableStackAll(disable bool)
	
	// SetCustomRecoveryHandler 设置自定义恢复处理器 (Set custom recovery handler)
	SetCustomRecoveryHandler(handler func(server.Context, interface{}))
	
	// SetErrorResponseHandler 设置错误响应处理器 (Set error response handler)
	SetErrorResponseHandler(handler func(server.Context, interface{}, []byte))
	
	// HandlePanic 处理panic (Handle panic)
	HandlePanic(ctx server.Context, recovered interface{})
	
	// FormatPanicInfo 格式化panic信息 (Format panic information)
	FormatPanicInfo(info PanicInfo) string
}

// RecoveryFactory 恢复中间件工厂接口 (Recovery middleware factory interface)
type RecoveryFactory interface {
	// CreateRecovery 创建恢复中间件 (Create recovery middleware)
	CreateRecovery(config *RecoveryConfig) PanicRecovery
	
	// CreateDefaultRecovery 创建默认恢复中间件 (Create default recovery middleware)
	CreateDefaultRecovery() PanicRecovery
}

// RecoveryMode 恢复模式 (Recovery mode)
type RecoveryMode int

const (
	// RecoveryModeDefault 默认模式 (Default mode)
	RecoveryModeDefault RecoveryMode = iota
	
	// RecoveryModeDebug 调试模式 (Debug mode)
	// 包含详细的堆栈信息 (Includes detailed stack information)
	RecoveryModeDebug
	
	// RecoveryModeProduction 生产模式 (Production mode)
	// 隐藏敏感信息 (Hides sensitive information)
	RecoveryModeProduction
	
	// RecoveryModeSilent 静默模式 (Silent mode)
	// 不输出任何信息 (No output)
	RecoveryModeSilent
)

// String 返回恢复模式字符串 (Return recovery mode string)
func (r RecoveryMode) String() string {
	switch r {
	case RecoveryModeDefault:
		return "default"
	case RecoveryModeDebug:
		return "debug"
	case RecoveryModeProduction:
		return "production"
	case RecoveryModeSilent:
		return "silent"
	default:
		return "unknown"
	}
}

// ErrorResponse 错误响应 (Error response)
type ErrorResponse struct {
	// Error 错误信息 (Error message)
	Error string `json:"error"`
	
	// Message 详细消息 (Detailed message)
	Message string `json:"message"`
	
	// Code 错误代码 (Error code)
	Code int `json:"code,omitempty"`
	
	// Timestamp 时间戳 (Timestamp)
	Timestamp string `json:"timestamp"`
	
	// RequestID 请求ID (Request ID)
	RequestID string `json:"request_id,omitempty"`
	
	// Stack 堆栈信息 (Stack trace)
	Stack string `json:"stack,omitempty"`
}

// DefaultErrorResponse 创建默认错误响应 (Create default error response)
func DefaultErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{
		Error:   "Internal Server Error",
		Message: message,
		Code:    500,
	}
}

// PanicHandler Panic处理器函数类型 (Panic handler function type)
type PanicHandler func(ctx server.Context, recovered interface{}, stack []byte)

// DefaultPanicHandler 默认panic处理器 (Default panic handler)
func DefaultPanicHandler(ctx server.Context, recovered interface{}, stack []byte) {
	// 创建错误响应 (Create error response)
	response := DefaultErrorResponse("An unexpected error occurred")
	
	// 发送JSON响应 (Send JSON response)
	_ = ctx.JSON(500, response)
}

// DebugPanicHandler 调试模式panic处理器 (Debug mode panic handler)
func DebugPanicHandler(ctx server.Context, recovered interface{}, stack []byte) {
	// 创建详细错误响应 (Create detailed error response)
	response := &ErrorResponse{
		Error:   "Internal Server Error",
		Message: "Panic recovered",
		Code:    500,
		Stack:   string(stack),
	}
	
	// 发送JSON响应 (Send JSON response)
	_ = ctx.JSON(500, response)
}

// ProductionPanicHandler 生产模式panic处理器 (Production mode panic handler)
func ProductionPanicHandler(ctx server.Context, recovered interface{}, stack []byte) {
	// 创建简化错误响应 (Create simplified error response)
	response := &ErrorResponse{
		Error:   "Internal Server Error",
		Message: "An unexpected error occurred",
		Code:    500,
	}
	
	// 发送JSON响应 (Send JSON response)
	_ = ctx.JSON(500, response)
}

// GetPanicHandlerByMode 根据模式获取panic处理器 (Get panic handler by mode)
func GetPanicHandlerByMode(mode RecoveryMode) PanicHandler {
	switch mode {
	case RecoveryModeDebug:
		return DebugPanicHandler
	case RecoveryModeProduction:
		return ProductionPanicHandler
	case RecoveryModeSilent:
		return func(ctx server.Context, recovered interface{}, stack []byte) {
			// 静默模式不做任何处理 (Silent mode does nothing)
		}
	default:
		return DefaultPanicHandler
	}
} 