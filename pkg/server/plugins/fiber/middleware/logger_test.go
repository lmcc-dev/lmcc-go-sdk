/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Tests for Fiber Logger middleware / Fiber Logger中间件测试
 */

package middleware

import (
	"testing"
	"time"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
	"github.com/stretchr/testify/assert"
)

// TestNewLoggerMiddleware 测试创建Logger中间件 (Test creating Logger middleware)
func TestNewLoggerMiddleware(t *testing.T) {
	config := &LoggerConfig{
		Enabled:   true,
		Format:    "${time} ${status} ${method} ${path}",
		TimeZone:  "UTC",
		TimeFormat: time.RFC3339,
	}
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	middleware := NewLoggerMiddleware(config, serviceContainer)
	assert.NotNil(t, middleware)
}

// TestLoggerMiddleware_GetConfig 测试获取Logger配置 (Test getting Logger configuration)
func TestLoggerMiddleware_GetConfig(t *testing.T) {
	config := &LoggerConfig{
		Enabled:   true,
		Format:    "${time} ${status} ${method} ${path}",
		TimeZone:  "Asia/Shanghai",
		TimeFormat: time.RFC822,
	}
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	middleware := NewLoggerMiddleware(config, serviceContainer)
	
	retrievedConfig := middleware.GetConfig()
	assert.Equal(t, config, retrievedConfig)
}

// TestLoggerMiddleware_SetConfig 测试设置Logger配置 (Test setting Logger configuration)
func TestLoggerMiddleware_SetConfig(t *testing.T) {
	originalConfig := &LoggerConfig{
		Enabled: true,
		Format:  "${method} ${path}",
	}
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	middleware := NewLoggerMiddleware(originalConfig, serviceContainer)
	
	newConfig := &LoggerConfig{
		Enabled:   false,
		Format:    "${time} ${status} ${method} ${path} ${latency}",
		TimeZone:  "America/New_York",
		TimeFormat: time.Kitchen,
	}
	
	middleware.SetConfig(newConfig)
	assert.Equal(t, newConfig, middleware.GetConfig())
	
	// 测试设置nil配置 (Test setting nil config)
	middleware.SetConfig(nil)
	assert.Equal(t, newConfig, middleware.GetConfig()) // SetConfig(nil) 不会改变配置，仍然是newConfig
} 