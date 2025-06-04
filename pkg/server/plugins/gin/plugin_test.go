/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Gin插件单元测试 (Gin plugin unit tests)
 */

package gin

import (
	"testing"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
)

// TestNewPlugin 测试创建Gin插件 (Test creating Gin plugin)
func TestNewPlugin(t *testing.T) {
	plugin := NewPlugin()
	
	if plugin == nil {
		t.Fatal("NewPlugin() returned nil")
	}
	
	if plugin.Name() != "gin" {
		t.Errorf("Expected plugin name 'gin', got '%s'", plugin.Name())
	}
	
	if plugin.Version() == "" {
		t.Error("Plugin version should not be empty")
	}
	
	if plugin.Description() == "" {
		t.Error("Plugin description should not be empty")
	}
}

// TestPluginDefaultConfig 测试默认配置 (Test default configuration)
func TestPluginDefaultConfig(t *testing.T) {
	plugin := NewPlugin()
	config := plugin.DefaultConfig()
	
	if config == nil {
		t.Fatal("DefaultConfig() returned nil")
	}
	
	// 转换为ServerConfig类型 (Convert to ServerConfig type)
	serverConfig, ok := config.(*server.ServerConfig)
	if !ok {
		t.Fatal("DefaultConfig() should return *server.ServerConfig")
	}
	
	if serverConfig.Framework != "gin" {
		t.Errorf("Expected framework 'gin', got '%s'", serverConfig.Framework)
	}
	
	if serverConfig.Mode != "debug" {
		t.Errorf("Expected mode 'debug', got '%s'", serverConfig.Mode)
	}
	
	if serverConfig.Plugins == nil {
		t.Error("Plugins configuration should not be nil")
	}
	
	ginConfig, ok := serverConfig.Plugins["gin"]
	if !ok {
		t.Error("Gin specific configuration should exist")
	}
	
	if ginConfig == nil {
		t.Error("Gin configuration should not be nil")
	}
}

// TestPluginValidateConfig 测试配置验证 (Test configuration validation)
func TestPluginValidateConfig(t *testing.T) {
	plugin := NewPlugin()
	
	// 测试有效配置 (Test valid configuration)
	validConfig := &server.ServerConfig{
		Framework: "gin",
		Host:      "localhost",
		Port:      8080,
	}
	
	err := plugin.ValidateConfig(validConfig)
	if err != nil {
		t.Errorf("ValidateConfig() should pass for valid config: %v", err)
	}
	
	// 测试nil配置 (Test nil configuration)
	err = plugin.ValidateConfig(nil)
	if err != nil {
		t.Errorf("ValidateConfig() should pass for nil config: %v", err)
	}
	
	// 测试无效配置类型 (Test invalid configuration type)
	err = plugin.ValidateConfig("invalid")
	if err == nil {
		t.Error("ValidateConfig() should fail for invalid config type")
	}
}

// TestPluginGetConfigSchema 测试配置模式 (Test configuration schema)
func TestPluginGetConfigSchema(t *testing.T) {
	plugin := NewPlugin()
	schema := plugin.GetConfigSchema()
	
	if schema == nil {
		t.Fatal("GetConfigSchema() returned nil")
	}
	
	// 检查模式是否为map (Check if schema is a map)
	schemaMap, ok := schema.(map[string]interface{})
	if !ok {
		t.Fatal("GetConfigSchema() should return map[string]interface{}")
	}
	
	// 检查必要字段 (Check required fields)
	if schemaMap["type"] != "object" {
		t.Error("Schema type should be 'object'")
	}
	
	properties, ok := schemaMap["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("Schema should have properties")
	}
	
	if properties["framework"] == nil {
		t.Error("Schema should have framework property")
	}
}

// TestPluginCreateFramework 测试创建框架实例 (Test creating framework instance)
func TestPluginCreateFramework(t *testing.T) {
	plugin := NewPlugin()
	config := plugin.DefaultConfig()
	serviceContainer := services.NewServiceContainerWithDefaults()
	
	framework, err := plugin.CreateFramework(config, serviceContainer)
	if err != nil {
		t.Fatalf("CreateFramework() failed: %v", err)
	}
	
	if framework == nil {
		t.Fatal("CreateFramework() returned nil framework")
	}
	
	// 检查是否为Gin服务器 (Check if it's a Gin server)
	ginServer, ok := framework.(*GinServer)
	if !ok {
		t.Error("Created framework is not a GinServer")
	}
	
	if ginServer.GetConfig().Framework != "gin" {
		t.Error("Server config framework should be 'gin'")
	}
	
	// 检查服务容器 (Check service container)
	if ginServer.GetServices() == nil {
		t.Error("Server should have service container")
	}
}

// TestPluginCreateFrameworkWithNilConfig 测试使用nil配置创建 (Test creating with nil config)
func TestPluginCreateFrameworkWithNilConfig(t *testing.T) {
	plugin := NewPlugin()
	serviceContainer := services.NewServiceContainerWithDefaults()
	
	framework, err := plugin.CreateFramework(nil, serviceContainer)
	if err != nil {
		t.Fatalf("CreateFramework() with nil config failed: %v", err)
	}
	
	if framework == nil {
		t.Fatal("CreateFramework() returned nil framework")
	}
	
	ginServer := framework.(*GinServer)
	if ginServer.GetConfig().Framework != "gin" {
		t.Error("Server config framework should be 'gin'")
	}
}

// TestPluginCreateFrameworkWithNilServices 测试使用nil服务容器创建 (Test creating with nil service container)
func TestPluginCreateFrameworkWithNilServices(t *testing.T) {
	plugin := NewPlugin()
	config := plugin.DefaultConfig()
	
	framework, err := plugin.CreateFramework(config, nil)
	if err != nil {
		t.Fatalf("CreateFramework() with nil services failed: %v", err)
	}
	
	if framework == nil {
		t.Fatal("CreateFramework() returned nil framework")
	}
	
	ginServer := framework.(*GinServer)
	if ginServer.GetServices() == nil {
		t.Error("Server should have default service container")
	}
}

// TestPluginCreateFrameworkWithInvalidConfig 测试使用无效配置创建 (Test creating with invalid config)
func TestPluginCreateFrameworkWithInvalidConfig(t *testing.T) {
	plugin := NewPlugin()
	serviceContainer := services.NewServiceContainerWithDefaults()
	
	// 测试无效配置类型 (Test invalid config type)
	_, err := plugin.CreateFramework("invalid", serviceContainer)
	if err == nil {
		t.Error("CreateFramework() should fail with invalid config type")
	}
	
	// 测试nil配置应该成功（使用默认配置） (Test nil config should succeed with default config)
	_, err = plugin.CreateFramework(nil, serviceContainer)
	if err != nil {
		t.Errorf("CreateFramework() with nil config should succeed: %v", err)
	}
}

// TestPluginRegistration 测试插件注册 (Test plugin registration)
func TestPluginRegistration(t *testing.T) {
	// 创建新的注册表用于测试 (Create new registry for testing)
	registry := server.NewPluginRegistry()
	
	plugin := NewPlugin()
	err := registry.Register(plugin)
	if err != nil {
		t.Fatalf("Failed to register plugin: %v", err)
	}
	
	// 检查插件是否已注册 (Check if plugin is registered)
	registeredPlugin, err := registry.Get("gin")
	if err != nil {
		t.Errorf("Plugin should be registered: %v", err)
		return
	}
	
	if registeredPlugin.Name() != "gin" {
		t.Error("Registered plugin name mismatch")
	}
}

// TestPluginDuplicateRegistration 测试重复注册 (Test duplicate registration)
func TestPluginDuplicateRegistration(t *testing.T) {
	registry := server.NewPluginRegistry()
	plugin := NewPlugin()
	
	// 第一次注册应该成功 (First registration should succeed)
	err := registry.Register(plugin)
	if err != nil {
		t.Fatalf("First registration failed: %v", err)
	}
	
	// 第二次注册应该失败 (Second registration should fail)
	err = registry.Register(plugin)
	if err == nil {
		t.Error("Duplicate registration should fail")
	}
}

// TestPluginIntegrationWithRegistry 测试与注册表的集成 (Test integration with registry)
func TestPluginIntegrationWithRegistry(t *testing.T) {
	// 注册插件 (Register plugin)
	err := server.RegisterFramework(NewPlugin())
	if err != nil {
		// 可能已经注册过了，这是正常的 (Might already be registered, this is normal)
		t.Logf("Plugin registration returned error (might be already registered): %v", err)
	}
	
	// 创建服务器 (Create server)
	config := &server.ServerConfig{
		Framework: "gin",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	
	framework, err := server.CreateServer("gin", config, serviceContainer)
	if err != nil {
		t.Fatalf("Failed to create server through registry: %v", err)
	}
	
	if framework == nil {
		t.Fatal("Created server is nil")
	}
	
	// 检查是否为Gin服务器 (Check if it's a Gin server)
	ginServer, ok := framework.(*GinServer)
	if !ok {
		t.Error("Created server is not a GinServer")
	}
	
	if ginServer.GetConfig().Framework != "gin" {
		t.Error("Server framework should be 'gin'")
	}
}

// BenchmarkPluginCreate 基准测试插件创建 (Benchmark plugin creation)
func BenchmarkPluginCreate(b *testing.B) {
	plugin := NewPlugin()
	config := plugin.DefaultConfig()
	serviceContainer := services.NewServiceContainerWithDefaults()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		framework, err := plugin.CreateFramework(config, serviceContainer)
		if err != nil {
			b.Fatalf("CreateFramework() failed: %v", err)
		}
		_ = framework
	}
}

// BenchmarkPluginValidateConfig 基准测试配置验证 (Benchmark configuration validation)
func BenchmarkPluginValidateConfig(b *testing.B) {
	plugin := NewPlugin()
	config := &server.ServerConfig{
		Framework: "gin",
		Host:      "localhost",
		Port:      8080,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := plugin.ValidateConfig(config)
		if err != nil {
			b.Fatalf("ValidateConfig() failed: %v", err)
		}
	}
} 