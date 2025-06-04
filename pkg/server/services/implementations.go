/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: 服务接口实现 (Service interface implementations)
 */

package services

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

// LoggerImpl 日志服务实现 (Logger service implementation)
type LoggerImpl struct {
	logger log.Logger
}

// NewLoggerImpl 创建日志服务实现 (Create logger service implementation)
func NewLoggerImpl(logger log.Logger) Logger {
	if logger == nil {
		logger = log.Std() // 使用标准日志器 (Use standard logger)
	}
	return &LoggerImpl{logger: logger}
}

// Debug 实现Logger接口 (Implement Logger interface)
func (l *LoggerImpl) Debug(msg string) {
	l.logger.Debug(msg)
}

func (l *LoggerImpl) Debugf(template string, args ...interface{}) {
	l.logger.Debugf(template, args...)
}

func (l *LoggerImpl) Debugw(msg string, keysAndValues ...interface{}) {
	l.logger.Debugw(msg, keysAndValues...)
}

func (l *LoggerImpl) Info(msg string) {
	l.logger.Info(msg)
}

func (l *LoggerImpl) Infof(template string, args ...interface{}) {
	l.logger.Infof(template, args...)
}

func (l *LoggerImpl) Infow(msg string, keysAndValues ...interface{}) {
	l.logger.Infow(msg, keysAndValues...)
}

func (l *LoggerImpl) Warn(msg string) {
	l.logger.Warn(msg)
}

func (l *LoggerImpl) Warnf(template string, args ...interface{}) {
	l.logger.Warnf(template, args...)
}

func (l *LoggerImpl) Warnw(msg string, keysAndValues ...interface{}) {
	l.logger.Warnw(msg, keysAndValues...)
}

func (l *LoggerImpl) Error(msg string) {
	l.logger.Error(msg)
}

func (l *LoggerImpl) Errorf(template string, args ...interface{}) {
	l.logger.Errorf(template, args...)
}

func (l *LoggerImpl) Errorw(msg string, keysAndValues ...interface{}) {
	l.logger.Errorw(msg, keysAndValues...)
}

func (l *LoggerImpl) Fatal(msg string) {
	l.logger.Fatal(msg)
}

func (l *LoggerImpl) Fatalf(template string, args ...interface{}) {
	l.logger.Fatalf(template, args...)
}

func (l *LoggerImpl) Panic(msg string) {
	l.logger.DPanic(msg) // 使用DPanic代替Panic (Use DPanic instead of Panic)
}

func (l *LoggerImpl) Panicf(template string, args ...interface{}) {
	l.logger.DPanicf(template, args...) // 使用DPanicf代替Panicf (Use DPanicf instead of Panicf)
}

// ErrorHandlerImpl 错误处理服务实现 (Error handler service implementation)
type ErrorHandlerImpl struct{}

// NewErrorHandlerImpl 创建错误处理服务实现 (Create error handler service implementation)
func NewErrorHandlerImpl() ErrorHandler {
	return &ErrorHandlerImpl{}
}

// New 实现ErrorHandler接口 (Implement ErrorHandler interface)
func (e *ErrorHandlerImpl) New(message string) error {
	return errors.New(message)
}

func (e *ErrorHandlerImpl) Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

func (e *ErrorHandlerImpl) Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

func (e *ErrorHandlerImpl) WithCode(err error, code interface{}) error {
	// 检查code是否实现了Coder接口 (Check if code implements Coder interface)
	if coder, ok := code.(errors.Coder); ok {
		return errors.WithCode(err, coder)
	}
	// 如果不是Coder，返回原错误 (If not Coder, return original error)
	return err
}

func (e *ErrorHandlerImpl) GetCode(err error) interface{} {
	// 使用GetCoder获取Coder (Use GetCoder to get Coder)
	if coder := errors.GetCoder(err); coder != nil {
		return coder.Code()
	}
	return nil
}

func (e *ErrorHandlerImpl) IsCode(err error, code interface{}) bool {
	// 检查code是否实现了Coder接口 (Check if code implements Coder interface)
	if coder, ok := code.(errors.Coder); ok {
		return errors.IsCode(err, coder)
	}
	return false
}

func (e *ErrorHandlerImpl) GetStackTrace(err error) string {
	// pkg/errors没有直接的GetStackTrace函数，我们使用fmt格式化 (pkg/errors doesn't have direct GetStackTrace, use fmt formatting)
	return fmt.Sprintf("%+v", err)
}

// ConfigManagerImpl 配置管理服务实现 (Config manager service implementation)
type ConfigManagerImpl struct {
	manager config.Manager
	viper   *viper.Viper
}

// NewConfigManagerImpl 创建配置管理服务实现 (Create config manager service implementation)
func NewConfigManagerImpl(manager config.Manager) ConfigManager {
	var v *viper.Viper
	if manager != nil {
		v = manager.GetViperInstance()
	}
	
	return &ConfigManagerImpl{
		manager: manager,
		viper:   v,
	}
}

// Get 实现ConfigManager接口 (Implement ConfigManager interface)
func (c *ConfigManagerImpl) Get(key string) interface{} {
	if c.viper != nil {
		return c.viper.Get(key)
	}
	return nil
}

func (c *ConfigManagerImpl) GetString(key string) string {
	if c.viper != nil {
		return c.viper.GetString(key)
	}
	return ""
}

func (c *ConfigManagerImpl) GetInt(key string) int {
	if c.viper != nil {
		return c.viper.GetInt(key)
	}
	return 0
}

func (c *ConfigManagerImpl) GetBool(key string) bool {
	if c.viper != nil {
		return c.viper.GetBool(key)
	}
	return false
}

func (c *ConfigManagerImpl) GetFloat64(key string) float64 {
	if c.viper != nil {
		return c.viper.GetFloat64(key)
	}
	return 0.0
}

func (c *ConfigManagerImpl) GetStringSlice(key string) []string {
	if c.viper != nil {
		return c.viper.GetStringSlice(key)
	}
	return nil
}

func (c *ConfigManagerImpl) Set(key string, value interface{}) {
	if c.viper != nil {
		c.viper.Set(key, value)
	}
}

func (c *ConfigManagerImpl) IsSet(key string) bool {
	if c.viper != nil {
		return c.viper.IsSet(key)
	}
	return false
}

func (c *ConfigManagerImpl) GetViperInstance() *viper.Viper {
	return c.viper
}

func (c *ConfigManagerImpl) RegisterCallback(callback func(v *viper.Viper, cfg any) error) {
	if c.manager != nil {
		c.manager.RegisterCallback(callback)
	}
}

func (c *ConfigManagerImpl) Unmarshal(rawVal interface{}) error {
	if c.viper != nil {
		return c.viper.Unmarshal(rawVal)
	}
	return errors.New("viper instance is nil")
}

func (c *ConfigManagerImpl) UnmarshalKey(key string, rawVal interface{}) error {
	if c.viper != nil {
		return c.viper.UnmarshalKey(key, rawVal)
	}
	return errors.New("viper instance is nil")
}

// ServiceContainerImpl 服务容器实现 (Service container implementation)
type ServiceContainerImpl struct {
	logger       Logger
	errorHandler ErrorHandler
	configManager ConfigManager
}

// NewServiceContainer 创建服务容器 (Create service container)
func NewServiceContainer() ServiceContainer {
	return &ServiceContainerImpl{}
}

// NewServiceContainerWithDefaults 创建带默认服务的服务容器 (Create service container with default services)
func NewServiceContainerWithDefaults() ServiceContainer {
	return &ServiceContainerImpl{
		logger:       NewLoggerImpl(nil), // 使用默认日志器 (Use default logger)
		errorHandler: NewErrorHandlerImpl(),
		configManager: NewConfigManagerImpl(nil), // 没有配置管理器 (No config manager)
	}
}

// GetLogger 实现ServiceContainer接口 (Implement ServiceContainer interface)
func (s *ServiceContainerImpl) GetLogger() Logger {
	return s.logger
}

func (s *ServiceContainerImpl) GetErrorHandler() ErrorHandler {
	return s.errorHandler
}

func (s *ServiceContainerImpl) GetConfigManager() ConfigManager {
	return s.configManager
}

func (s *ServiceContainerImpl) SetLogger(logger Logger) {
	s.logger = logger
}

func (s *ServiceContainerImpl) SetErrorHandler(handler ErrorHandler) {
	s.errorHandler = handler
}

func (s *ServiceContainerImpl) SetConfigManager(manager ConfigManager) {
	s.configManager = manager
} 