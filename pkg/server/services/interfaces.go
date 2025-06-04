/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: 服务接口定义 (Service interface definitions)
 */

package services

import (
	"github.com/spf13/viper"
)

// Logger 统一的日志服务接口 (Unified logger service interface)
type Logger interface {
	// Debug 记录调试级别日志 (Log debug level message)
	Debug(msg string)
	Debugf(template string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	
	// Info 记录信息级别日志 (Log info level message)
	Info(msg string)
	Infof(template string, args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	
	// Warn 记录警告级别日志 (Log warn level message)
	Warn(msg string)
	Warnf(template string, args ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	
	// Error 记录错误级别日志 (Log error level message)
	Error(msg string)
	Errorf(template string, args ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	
	// Fatal 记录致命错误并退出 (Log fatal error and exit)
	Fatal(msg string)
	Fatalf(template string, args ...interface{})
	
	// Panic 记录错误并panic (Log error and panic)
	Panic(msg string)
	Panicf(template string, args ...interface{})
}

// ErrorHandler 统一的错误处理服务接口 (Unified error handling service interface)
type ErrorHandler interface {
	// New 创建新错误 (Create new error)
	New(message string) error
	
	// Wrap 包装错误 (Wrap error)
	Wrap(err error, message string) error
	
	// Wrapf 格式化包装错误 (Wrap error with format)
	Wrapf(err error, format string, args ...interface{}) error
	
	// WithCode 为错误添加错误码 (Add error code to error)
	WithCode(err error, code interface{}) error
	
	// GetCode 获取错误码 (Get error code)
	GetCode(err error) interface{}
	
	// IsCode 检查错误码 (Check error code)
	IsCode(err error, code interface{}) bool
	
	// GetStackTrace 获取堆栈跟踪 (Get stack trace)
	GetStackTrace(err error) string
}

// ConfigManager 统一的配置管理服务接口 (Unified config management service interface)
type ConfigManager interface {
	// Get 获取配置值 (Get config value)
	Get(key string) interface{}
	
	// GetString 获取字符串配置 (Get string config)
	GetString(key string) string
	
	// GetInt 获取整数配置 (Get int config)
	GetInt(key string) int
	
	// GetBool 获取布尔配置 (Get bool config)
	GetBool(key string) bool
	
	// GetFloat64 获取浮点数配置 (Get float64 config)
	GetFloat64(key string) float64
	
	// GetStringSlice 获取字符串切片配置 (Get string slice config)
	GetStringSlice(key string) []string
	
	// Set 设置配置值 (Set config value)
	Set(key string, value interface{})
	
	// IsSet 检查配置是否设置 (Check if config is set)
	IsSet(key string) bool
	
	// GetViperInstance 获取底层Viper实例 (Get underlying Viper instance)
	GetViperInstance() *viper.Viper
	
	// RegisterCallback 注册配置变更回调 (Register config change callback)
	RegisterCallback(callback func(v *viper.Viper, cfg any) error)
	
	// Unmarshal 将配置解析到结构体 (Unmarshal config to struct)
	Unmarshal(rawVal interface{}) error
	
	// UnmarshalKey 将指定键的配置解析到结构体 (Unmarshal specific key config to struct)
	UnmarshalKey(key string, rawVal interface{}) error
}

// ServiceContainer 服务容器接口 (Service container interface)
// 提供统一的服务访问入口 (Provides unified service access point)
type ServiceContainer interface {
	// GetLogger 获取日志服务 (Get logger service)
	GetLogger() Logger
	
	// GetErrorHandler 获取错误处理服务 (Get error handler service)
	GetErrorHandler() ErrorHandler
	
	// GetConfigManager 获取配置管理服务 (Get config manager service)
	GetConfigManager() ConfigManager
	
	// SetLogger 设置日志服务 (Set logger service)
	SetLogger(logger Logger)
	
	// SetErrorHandler 设置错误处理服务 (Set error handler service)
	SetErrorHandler(handler ErrorHandler)
	
	// SetConfigManager 设置配置管理服务 (Set config manager service)
	SetConfigManager(manager ConfigManager)
} 