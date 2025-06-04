/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Fiber插件测试 (Fiber plugin tests)
 */

package fiber

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
)

func TestFiberPlugin_Name(t *testing.T) {
	plugin := NewPlugin()
	assert.Equal(t, "fiber", plugin.Name())
}

func TestFiberPlugin_Version(t *testing.T) {
	plugin := NewPlugin()
	assert.Equal(t, "v1.0.0", plugin.Version())
}

func TestFiberPlugin_Description(t *testing.T) {
	plugin := NewPlugin()
	assert.Contains(t, plugin.Description(), "Fiber")
}

func TestFiberPlugin_DefaultConfig(t *testing.T) {
	plugin := NewPlugin()
	config := plugin.DefaultConfig()
	assert.NotNil(t, config)
	
	serverConfig, ok := config.(*server.ServerConfig)
	require.True(t, ok)
	assert.Equal(t, "fiber", serverConfig.Framework)
	assert.Equal(t, "localhost", serverConfig.Host)
	assert.Equal(t, 8080, serverConfig.Port)
}

func TestFiberPlugin_ValidateConfig_Valid(t *testing.T) {
	plugin := NewPlugin()
	config := &server.ServerConfig{
		Framework: "fiber",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}
	
	err := plugin.ValidateConfig(config)
	assert.NoError(t, err)
}

func TestFiberPlugin_ValidateConfig_Invalid(t *testing.T) {
	plugin := NewPlugin()
	
	// 测试无效配置类型
	err := plugin.ValidateConfig("invalid")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidConfig, err)
	
	// 测试错误的框架名称
	config := &server.ServerConfig{
		Framework: "gin",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}
	err = plugin.ValidateConfig(config)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidFramework, err)
}

func TestFiberPlugin_CreateFramework_Success(t *testing.T) {
	plugin := NewPlugin()
	config := &server.ServerConfig{
		Framework: "fiber",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}
	
	serviceContainer := services.NewServiceContainer()
	framework, err := plugin.CreateFramework(config, serviceContainer)
	
	assert.NoError(t, err)
	assert.NotNil(t, framework)
	
	// 验证框架类型
	fiberServer, ok := framework.(*FiberServer)
	assert.True(t, ok)
	assert.NotNil(t, fiberServer.GetFiberApp())
}

func TestFiberPlugin_CreateFramework_InvalidConfig(t *testing.T) {
	plugin := NewPlugin()
	serviceContainer := services.NewServiceContainer()
	
	framework, err := plugin.CreateFramework("invalid", serviceContainer)
	
	assert.Error(t, err)
	assert.Nil(t, framework)
	assert.Equal(t, ErrInvalidConfig, err)
}

func TestFiberPlugin_GetConfigSchema(t *testing.T) {
	plugin := NewPlugin()
	schema := plugin.GetConfigSchema()
	assert.NotNil(t, schema)
	
	schemaMap, ok := schema.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "object", schemaMap["type"])
} 