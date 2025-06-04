/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Tests for Fiber Recovery middleware / Fiber Recovery中间件测试
 */

package middleware

import (
	"testing"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
	"github.com/stretchr/testify/assert"
)

// TestNewRecoveryMiddleware 测试创建Recovery中间件 (Test creating Recovery middleware)
func TestNewRecoveryMiddleware(t *testing.T) {
	config := &RecoveryConfig{
		Enabled:            true,
		PrintStack:         true,
		StackSize:          4096,
		DisableStackAll:    false,
		DisableColorOutput: false,
	}
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	middleware := NewRecoveryMiddleware(config, serviceContainer)
	assert.NotNil(t, middleware)
}

// TestRecoveryMiddleware_GetConfig 测试获取Recovery配置 (Test getting Recovery configuration)
func TestRecoveryMiddleware_GetConfig(t *testing.T) {
	config := &RecoveryConfig{
		Enabled:            true,
		PrintStack:         false,
		StackSize:          8192,
		DisableStackAll:    true,
		DisableColorOutput: true,
	}
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	middleware := NewRecoveryMiddleware(config, serviceContainer)
	
	retrievedConfig := middleware.GetConfig()
	assert.Equal(t, config, retrievedConfig)
}

// TestRecoveryMiddleware_SetConfig 测试设置Recovery配置 (Test setting Recovery configuration)
func TestRecoveryMiddleware_SetConfig(t *testing.T) {
	originalConfig := &RecoveryConfig{
		Enabled:    true,
		PrintStack: false,
	}
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	middleware := NewRecoveryMiddleware(originalConfig, serviceContainer)
	
	newConfig := &RecoveryConfig{
		Enabled:            false,
		PrintStack:         true,
		StackSize:          2048,
		DisableStackAll:    false,
		DisableColorOutput: true,
	}
	
	middleware.SetConfig(newConfig)
	assert.Equal(t, newConfig, middleware.GetConfig())
	
	// 测试设置nil配置 (Test setting nil config)
	middleware.SetConfig(nil)
	assert.Equal(t, newConfig, middleware.GetConfig()) // SetConfig(nil) 不会改变配置，仍然是newConfig
} 