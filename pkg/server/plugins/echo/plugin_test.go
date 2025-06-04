/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Echo插件测试 (Echo plugin tests)
 */

package echo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
)

func TestPlugin_BasicInfo(t *testing.T) {
	plugin := NewPlugin()
	
	assert.Equal(t, "echo", plugin.Name())
	assert.Equal(t, "v1.0.0", plugin.Version())
	assert.Contains(t, plugin.Description(), "Echo")
}

func TestPlugin_DefaultConfig(t *testing.T) {
	plugin := NewPlugin()
	config := plugin.DefaultConfig()
	
	serverConfig, ok := config.(*server.ServerConfig)
	require.True(t, ok, "Default config should be *server.ServerConfig")
	
	assert.Equal(t, "echo", serverConfig.Framework)
	assert.Equal(t, "0.0.0.0", serverConfig.Host)
	assert.Equal(t, 8080, serverConfig.Port)
	assert.Equal(t, "debug", serverConfig.Mode)
	assert.True(t, serverConfig.CORS.Enabled)
	assert.True(t, serverConfig.Middleware.Logger.Enabled)
	assert.True(t, serverConfig.Middleware.Recovery.Enabled)
	assert.True(t, serverConfig.GracefulShutdown.Enabled)
}

func TestPlugin_ValidateConfig_Valid(t *testing.T) {
	plugin := NewPlugin()
	config := &server.ServerConfig{
		Framework: "echo",
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}
	
	err := plugin.ValidateConfig(config)
	assert.NoError(t, err)
}

func TestPlugin_ValidateConfig_InvalidFramework(t *testing.T) {
	plugin := NewPlugin()
	config := &server.ServerConfig{
		Framework: "gin", // Wrong framework
		Host:      "localhost", 
		Port:      8080,
		Mode:      "debug",
	}
	
	err := plugin.ValidateConfig(config)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidFramework, err)
}

func TestPlugin_ValidateConfig_InvalidConfig(t *testing.T) {
	plugin := NewPlugin()
	
	err := plugin.ValidateConfig("invalid")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidConfig, err)
}

func TestPlugin_ValidateConfig_EmptyFramework(t *testing.T) {
	plugin := NewPlugin()
	config := &server.ServerConfig{
		Framework: "", // Empty framework
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}
	
	err := plugin.ValidateConfig(config)
	assert.Error(t, err)
}

func TestPlugin_ValidateConfig_InvalidPort(t *testing.T) {
	plugin := NewPlugin()
	config := &server.ServerConfig{
		Framework: "echo",
		Host:      "localhost",
		Port:      0, // Invalid port
		Mode:      "debug",
	}
	
	err := plugin.ValidateConfig(config)
	assert.Error(t, err)
}

func TestPlugin_GetConfigSchema(t *testing.T) {
	plugin := NewPlugin()
	schema := plugin.GetConfigSchema()
	
	schemaMap, ok := schema.(map[string]interface{})
	require.True(t, ok, "Schema should be a map")
	
	assert.Equal(t, "object", schemaMap["type"])
	
	properties, ok := schemaMap["properties"].(map[string]interface{})
	require.True(t, ok, "Properties should exist")
	
	// 验证framework属性 (Verify framework property)
	frameworkProp, ok := properties["framework"].(map[string]interface{})
	require.True(t, ok, "Framework property should exist")
	assert.Equal(t, "string", frameworkProp["type"])
	assert.Equal(t, "echo", frameworkProp["default"])
	
	// 验证required字段 (Verify required fields)
	required, ok := schemaMap["required"].([]string)
	require.True(t, ok, "Required should be a string array")
	assert.Contains(t, required, "framework")
	assert.Contains(t, required, "host")
	assert.Contains(t, required, "port")
}

func TestPlugin_CreateFramework_Success(t *testing.T) {
	plugin := NewPlugin()
	
	// 创建模拟服务容器 (Create mock service container)
	serviceContainer := services.NewServiceContainerWithDefaults()
	
	config := &server.ServerConfig{
		Framework:      "echo",
		Host:           "localhost",
		Port:           8080,
		Mode:           "debug",
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
		CORS: server.CORSConfig{
			Enabled: true,
		},
		Middleware: server.MiddlewareConfig{
			Logger: server.LoggerMiddlewareConfig{
				Enabled: true,
			},
			Recovery: server.RecoveryMiddlewareConfig{
				Enabled: true,
			},
		},
		GracefulShutdown: server.GracefulShutdownConfig{
			Enabled: true,
			Timeout: 30 * time.Second,
		},
	}
	
	framework, err := plugin.CreateFramework(config, serviceContainer)
	assert.NoError(t, err)
	assert.NotNil(t, framework)
	
	// 验证框架类型 (Verify framework type)
	echoServer, ok := framework.(*EchoServer)
	require.True(t, ok, "Framework should be *EchoServer")
	
	assert.Equal(t, config, echoServer.GetConfig())
	assert.NotNil(t, echoServer.GetEchoEngine())
}

func TestPlugin_CreateFramework_InvalidConfig(t *testing.T) {
	plugin := NewPlugin()
	serviceContainer := services.NewServiceContainerWithDefaults()
	
	// 测试无效配置类型 (Test invalid config type)
	_, err := plugin.CreateFramework("invalid", serviceContainer)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid configuration type")
}

func TestPlugin_CreateFramework_ValidationFailed(t *testing.T) {
	plugin := NewPlugin()
	serviceContainer := services.NewServiceContainerWithDefaults()
	
	config := &server.ServerConfig{
		Framework: "gin", // Wrong framework
		Host:      "localhost",
		Port:      8080,
		Mode:      "debug",
	}
	
	_, err := plugin.CreateFramework(config, serviceContainer)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
} 