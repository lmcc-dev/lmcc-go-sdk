/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: 服务模块单元测试 (Services module unit tests)
 */

package services

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

// TestNewLoggerImpl 测试创建日志服务实现 (Test creating logger service implementation)
func TestNewLoggerImpl(t *testing.T) {
	// 测试使用nil logger (Test with nil logger)
	logger := NewLoggerImpl(nil)
	assert.NotNil(t, logger)
	
	// 测试使用具体logger (Test with concrete logger)
	testLogger, _ := log.NewLogger(log.NewOptions())
	logger2 := NewLoggerImpl(testLogger)
	assert.NotNil(t, logger2)
}

// TestLoggerImpl_AllMethods 测试日志器的所有方法 (Test all logger methods)
func TestLoggerImpl_AllMethods(t *testing.T) {
	// 使用测试模式避免实际输出 (Use test mode to avoid actual output)
	opts := log.NewOptions()
	opts.Level = "error" // 设置高级别减少输出 (Set high level to reduce output)
	testLogger, _ := log.NewLogger(opts)
	
	logger := NewLoggerImpl(testLogger)
	
	// 测试所有日志方法不会panic (Test all log methods don't panic)
	assert.NotPanics(t, func() {
		logger.Debug("debug message")
		logger.Debugf("debug message %s", "formatted")
		logger.Debugw("debug message", "key", "value")
		
		logger.Info("info message")
		logger.Infof("info message %s", "formatted")
		logger.Infow("info message", "key", "value")
		
		logger.Warn("warn message")
		logger.Warnf("warn message %s", "formatted")
		logger.Warnw("warn message", "key", "value")
		
		logger.Error("error message")
		logger.Errorf("error message %s", "formatted")
		logger.Errorw("error message", "key", "value")
	})
}

// TestNewErrorHandlerImpl 测试创建错误处理服务实现 (Test creating error handler service implementation)
func TestNewErrorHandlerImpl(t *testing.T) {
	handler := NewErrorHandlerImpl()
	assert.NotNil(t, handler)
}

// TestErrorHandlerImpl_New 测试New方法 (Test New method)
func TestErrorHandlerImpl_New(t *testing.T) {
	handler := NewErrorHandlerImpl()
	
	err := handler.New("test error")
	assert.NotNil(t, err)
	assert.Equal(t, "test error", err.Error())
}

// TestErrorHandlerImpl_Wrap 测试Wrap方法 (Test Wrap method)
func TestErrorHandlerImpl_Wrap(t *testing.T) {
	handler := NewErrorHandlerImpl()
	
	originalErr := errors.New("original error")
	wrappedErr := handler.Wrap(originalErr, "wrapped message")
	
	assert.NotNil(t, wrappedErr)
	assert.Contains(t, wrappedErr.Error(), "wrapped message")
	assert.Contains(t, wrappedErr.Error(), "original error")
}

// TestErrorHandlerImpl_Wrapf 测试Wrapf方法 (Test Wrapf method)
func TestErrorHandlerImpl_Wrapf(t *testing.T) {
	handler := NewErrorHandlerImpl()
	
	originalErr := errors.New("original error")
	wrappedErr := handler.Wrapf(originalErr, "wrapped %s", "message")
	
	assert.NotNil(t, wrappedErr)
	assert.Contains(t, wrappedErr.Error(), "wrapped message")
	assert.Contains(t, wrappedErr.Error(), "original error")
}

// TestErrorHandlerImpl_WithCode 测试WithCode方法 (Test WithCode method)
func TestErrorHandlerImpl_WithCode(t *testing.T) {
	handler := NewErrorHandlerImpl()
	
	// 测试有效的Coder (Test valid Coder)
	testCoder := errors.NewCoder(404, 404, "NOT_FOUND", "Resource not found")
	originalErr := errors.New("test error")
	errorWithCode := handler.WithCode(originalErr, testCoder)
	
	assert.NotNil(t, errorWithCode)
	
	// 测试无效的code类型 (Test invalid code type)
	errorWithInvalidCode := handler.WithCode(originalErr, "invalid")
	assert.Equal(t, originalErr, errorWithInvalidCode)
}

// TestErrorHandlerImpl_GetCode 测试GetCode方法 (Test GetCode method)
func TestErrorHandlerImpl_GetCode(t *testing.T) {
	handler := NewErrorHandlerImpl()
	
	// 测试带错误码的错误 (Test error with code)
	testCoder := errors.NewCoder(404, 404, "NOT_FOUND", "Resource not found")
	originalErr := errors.New("test error")
	errorWithCode := errors.WithCode(originalErr, testCoder)
	
	code := handler.GetCode(errorWithCode)
	assert.NotNil(t, code)
	assert.Equal(t, 404, code)
	
	// 测试没有错误码的错误 (Test error without code)
	normalErr := errors.New("normal error")
	nilCode := handler.GetCode(normalErr)
	assert.Nil(t, nilCode)
}

// TestErrorHandlerImpl_IsCode 测试IsCode方法 (Test IsCode method)
func TestErrorHandlerImpl_IsCode(t *testing.T) {
	handler := NewErrorHandlerImpl()
	
	// 创建带错误码的错误 (Create error with code)
	testCoder := errors.NewCoder(404, 404, "NOT_FOUND", "Resource not found")
	originalErr := errors.New("test error")
	errorWithCode := errors.WithCode(originalErr, testCoder)
	
	// 测试匹配的错误码 (Test matching error code)
	isMatch := handler.IsCode(errorWithCode, testCoder)
	assert.True(t, isMatch)
	
	// 测试不匹配的错误码 (Test non-matching error code)
	otherCoder := errors.NewCoder(500, 500, "INTERNAL_ERROR", "Internal error")
	isNotMatch := handler.IsCode(errorWithCode, otherCoder)
	assert.False(t, isNotMatch)
	
	// 测试无效的code类型 (Test invalid code type)
	isInvalid := handler.IsCode(errorWithCode, "invalid")
	assert.False(t, isInvalid)
}

// TestErrorHandlerImpl_GetStackTrace 测试GetStackTrace方法 (Test GetStackTrace method)
func TestErrorHandlerImpl_GetStackTrace(t *testing.T) {
	handler := NewErrorHandlerImpl()
	
	err := errors.New("test error")
	stackTrace := handler.GetStackTrace(err)
	
	assert.NotEmpty(t, stackTrace)
	assert.Contains(t, stackTrace, "test error")
}

// TestNewConfigManagerImpl 测试创建配置管理服务实现 (Test creating config manager service implementation)
func TestNewConfigManagerImpl(t *testing.T) {
	// 测试使用nil manager (Test with nil manager)
	configManager := NewConfigManagerImpl(nil)
	assert.NotNil(t, configManager)
	
	// 测试使用具体manager (Test with concrete manager)
	mockManager := &mockConfigManager{viper: viper.New()}
	configManager2 := NewConfigManagerImpl(mockManager)
	assert.NotNil(t, configManager2)
}

// TestConfigManagerImpl_GetMethods 测试配置获取方法 (Test config get methods)
func TestConfigManagerImpl_GetMethods(t *testing.T) {
	// 创建临时配置 (Create temporary config)
	v := viper.New()
	v.Set("test_string", "hello")
	v.Set("test_int", 42)
	v.Set("test_bool", true)
	v.Set("test_float", 3.14)
	v.Set("test_slice", []string{"a", "b", "c"})
	
	// 创建mock manager (Create mock manager)
	mockManager := &mockConfigManager{viper: v}
	configManager := NewConfigManagerImpl(mockManager)
	
	// 测试各种获取方法 (Test various get methods)
	assert.Equal(t, "hello", configManager.Get("test_string"))
	assert.Equal(t, "hello", configManager.GetString("test_string"))
	assert.Equal(t, 42, configManager.GetInt("test_int"))
	assert.Equal(t, true, configManager.GetBool("test_bool"))
	assert.Equal(t, 3.14, configManager.GetFloat64("test_float"))
	assert.Equal(t, []string{"a", "b", "c"}, configManager.GetStringSlice("test_slice"))
	
	// 测试不存在的键 (Test non-existent keys)
	assert.Nil(t, configManager.Get("non_existent"))
	assert.Equal(t, "", configManager.GetString("non_existent"))
	assert.Equal(t, 0, configManager.GetInt("non_existent"))
	assert.Equal(t, false, configManager.GetBool("non_existent"))
	assert.Equal(t, 0.0, configManager.GetFloat64("non_existent"))
	assert.Nil(t, configManager.GetStringSlice("non_existent"))
}

// TestConfigManagerImpl_SetAndIsSet 测试Set和IsSet方法 (Test Set and IsSet methods)
func TestConfigManagerImpl_SetAndIsSet(t *testing.T) {
	v := viper.New()
	mockManager := &mockConfigManager{viper: v}
	configManager := NewConfigManagerImpl(mockManager)
	
	// 测试Set方法 (Test Set method)
	configManager.Set("new_key", "new_value")
	assert.Equal(t, "new_value", configManager.GetString("new_key"))
	
	// 测试IsSet方法 (Test IsSet method)
	assert.True(t, configManager.IsSet("new_key"))
	assert.False(t, configManager.IsSet("non_existent_key"))
}

// TestConfigManagerImpl_NilViper 测试nil viper的情况 (Test nil viper case)
func TestConfigManagerImpl_NilViper(t *testing.T) {
	configManager := NewConfigManagerImpl(nil)
	
	// 所有Get方法应该返回零值 (All Get methods should return zero values)
	assert.Nil(t, configManager.Get("any_key"))
	assert.Equal(t, "", configManager.GetString("any_key"))
	assert.Equal(t, 0, configManager.GetInt("any_key"))
	assert.Equal(t, false, configManager.GetBool("any_key"))
	assert.Equal(t, 0.0, configManager.GetFloat64("any_key"))
	assert.Nil(t, configManager.GetStringSlice("any_key"))
	assert.False(t, configManager.IsSet("any_key"))
	
	// Set方法应该不会panic (Set method should not panic)
	assert.NotPanics(t, func() {
		configManager.Set("any_key", "any_value")
	})
}

// TestConfigManagerImpl_Unmarshal 测试Unmarshal方法 (Test Unmarshal method)
func TestConfigManagerImpl_Unmarshal(t *testing.T) {
	v := viper.New()
	v.Set("name", "test")
	v.Set("age", 25)
	
	mockManager := &mockConfigManager{viper: v}
	configManager := NewConfigManagerImpl(mockManager)
	
	type TestStruct struct {
		Name string `mapstructure:"name"`
		Age  int    `mapstructure:"age"`
	}
	
	var result TestStruct
	err := configManager.Unmarshal(&result)
	assert.NoError(t, err)
	assert.Equal(t, "test", result.Name)
	assert.Equal(t, 25, result.Age)
	
	// 测试nil viper (Test nil viper)
	nilConfigManager := NewConfigManagerImpl(nil)
	err = nilConfigManager.Unmarshal(&result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "viper instance is nil")
}

// TestConfigManagerImpl_UnmarshalKey 测试UnmarshalKey方法 (Test UnmarshalKey method)
func TestConfigManagerImpl_UnmarshalKey(t *testing.T) {
	v := viper.New()
	v.Set("user.name", "test")
	v.Set("user.age", 25)
	
	mockManager := &mockConfigManager{viper: v}
	configManager := NewConfigManagerImpl(mockManager)
	
	type User struct {
		Name string `mapstructure:"name"`
		Age  int    `mapstructure:"age"`
	}
	
	var user User
	err := configManager.UnmarshalKey("user", &user)
	assert.NoError(t, err)
	assert.Equal(t, "test", user.Name)
	assert.Equal(t, 25, user.Age)
	
	// 测试nil viper (Test nil viper)
	nilConfigManager := NewConfigManagerImpl(nil)
	err = nilConfigManager.UnmarshalKey("user", &user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "viper instance is nil")
}

// TestServiceContainerImpl 测试服务容器实现 (Test service container implementation)
func TestServiceContainerImpl(t *testing.T) {
	// 测试创建空容器 (Test creating empty container)
	container := NewServiceContainer()
	assert.NotNil(t, container)
	
	// 所有服务应该为nil (All services should be nil)
	assert.Nil(t, container.GetLogger())
	assert.Nil(t, container.GetErrorHandler())
	assert.Nil(t, container.GetConfigManager())
}

// TestServiceContainerWithDefaults 测试带默认服务的容器 (Test container with default services)
func TestServiceContainerWithDefaults(t *testing.T) {
	container := NewServiceContainerWithDefaults()
	assert.NotNil(t, container)
	
	// 所有服务都应该有默认实现 (All services should have default implementations)
	assert.NotNil(t, container.GetLogger())
	assert.NotNil(t, container.GetErrorHandler())
	assert.NotNil(t, container.GetConfigManager())
}

// TestServiceContainer_SetMethods 测试服务容器的Set方法 (Test service container Set methods)
func TestServiceContainer_SetMethods(t *testing.T) {
	container := NewServiceContainer()
	
	// 创建测试服务 (Create test services)
	logger := NewLoggerImpl(nil)
	errorHandler := NewErrorHandlerImpl()
	configManager := NewConfigManagerImpl(nil)
	
	// 设置服务 (Set services)
	container.SetLogger(logger)
	container.SetErrorHandler(errorHandler)
	container.SetConfigManager(configManager)
	
	// 验证服务已设置 (Verify services are set)
	assert.Equal(t, logger, container.GetLogger())
	assert.Equal(t, errorHandler, container.GetErrorHandler())
	assert.Equal(t, configManager, container.GetConfigManager())
}

// BenchmarkLoggerImpl_Info 基准测试日志器Info方法 (Benchmark logger Info method)
func BenchmarkLoggerImpl_Info(b *testing.B) {
	opts := log.NewOptions()
	opts.Level = "error" // 设置高级别减少IO (Set high level to reduce IO)
	testLogger, _ := log.NewLogger(opts)
	logger := NewLoggerImpl(testLogger)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark test message")
	}
}

// BenchmarkErrorHandlerImpl_New 基准测试错误处理器New方法 (Benchmark error handler New method)
func BenchmarkErrorHandlerImpl_New(b *testing.B) {
	handler := NewErrorHandlerImpl()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = handler.New("benchmark test error")
	}
}

// mockConfigManager 模拟配置管理器 (Mock config manager)
type mockConfigManager struct {
	viper *viper.Viper
}

func (m *mockConfigManager) GetViperInstance() *viper.Viper {
	return m.viper
}

func (m *mockConfigManager) RegisterCallback(callback func(v *viper.Viper, cfg any) error) {
	// 空实现 (Empty implementation)
}

func (m *mockConfigManager) RegisterSectionChangeCallback(sectionKey string, callback config.SectionChangeCallback) {
	// 空实现 (Empty implementation for interface compliance)
} 