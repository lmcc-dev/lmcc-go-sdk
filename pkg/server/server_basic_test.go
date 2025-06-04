/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: 服务器基本功能单元测试 (Server basic functionality unit tests)
 */

package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestServerConfig_Validate 测试配置验证 (Test configuration validation)
func TestServerConfig_Validate(t *testing.T) {
	// 测试有效配置 (Test valid configuration)
	validConfig := &ServerConfig{
		Framework: "gin",
		Host:      "localhost",
		Port:      8080,
		Mode:      "development",
	}
	
	err := validConfig.Validate()
	assert.NoError(t, err)
	
	// 测试无效配置 - 空框架名 (Test invalid configuration - empty framework name)
	invalidConfig := &ServerConfig{
		Framework: "",
		Host:      "localhost",
		Port:      8080,
		Mode:      "development",
	}
	
	err = invalidConfig.Validate()
	assert.Error(t, err)
	
	// 测试无效配置 - 无效端口 (Test invalid configuration - invalid port)
	invalidPortConfig := &ServerConfig{
		Framework: "gin",
		Host:      "localhost",
		Port:      0,
		Mode:      "development",
	}
	
	err = invalidPortConfig.Validate()
	assert.Error(t, err)
}

// TestNewServerFactory_Basic 测试创建服务器工厂基本功能 (Test server factory basic functionality)
func TestNewServerFactory_Basic(t *testing.T) {
	factory := NewServerFactory()
	
	assert.NotNil(t, factory)
	assert.NotNil(t, factory.registry)
	
	// 测试列出插件 (Test listing plugins)
	plugins := factory.ListPlugins()
	assert.NotNil(t, plugins)
	
	// 测试获取所有插件信息 (Test getting all plugin info)
	allInfo := factory.GetAllPluginInfo()
	assert.NotNil(t, allInfo)
}

// TestGlobalFactory_Basic 测试全局工厂基本功能 (Test global factory basic functionality)
func TestGlobalFactory_Basic(t *testing.T) {
	// 测试全局函数可以调用 (Test global functions can be called)
	plugins := ListPlugins()
	assert.NotNil(t, plugins)
	
	allInfo := GetAllPluginInfo()
	assert.NotNil(t, allInfo)
}

// TestServerManager_Basic 测试服务器管理器基本属性 (Test server manager basic properties)
func TestServerManager_Basic(t *testing.T) {
	// 使用nil framework进行基本测试 (Use nil framework for basic testing)
	config := &ServerConfig{
		Framework: "test",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}
	
	manager := NewServerManager(nil, config)
	
	assert.NotNil(t, manager)
	assert.Equal(t, config, manager.GetConfig())
	assert.Nil(t, manager.GetFramework())
	assert.False(t, manager.IsRunning())
} 