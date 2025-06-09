/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: 服务器工厂和管理器单元测试 (Server factory and manager unit tests)
 */

package server

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWebFramework 模拟Web框架用于测试 (Mock web framework for testing)
type MockWebFramework struct {
	mock.Mock
	stopChan chan struct{} // 用于控制Start方法的阻塞 (Used to control blocking in Start method)
}

func (m *MockWebFramework) Start(ctx context.Context) error {
	args := m.Called(ctx)
	if err := args.Error(0); err != nil {
		return err
	}
	
	// 如果没有错误，则阻塞直到stopChan被关闭 (If no error, block until stopChan is closed)
	if m.stopChan == nil {
		m.stopChan = make(chan struct{})
	}
	
	select {
	case <-m.stopChan:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (m *MockWebFramework) Stop(ctx context.Context) error {
	args := m.Called(ctx)
	if m.stopChan != nil {
		close(m.stopChan)
		m.stopChan = nil
	}
	return args.Error(0)
}

func (m *MockWebFramework) RegisterRoute(method, path string, handler Handler) error {
	args := m.Called(method, path, handler)
	return args.Error(0)
}

func (m *MockWebFramework) RegisterMiddleware(middleware Middleware) error {
	args := m.Called(middleware)
	return args.Error(0)
}

func (m *MockWebFramework) Group(prefix string, middlewares ...Middleware) RouteGroup {
	args := m.Called(prefix, middlewares)
	return args.Get(0).(RouteGroup)
}

func (m *MockWebFramework) GetNativeEngine() interface{} {
	args := m.Called()
	return args.Get(0)
}

func (m *MockWebFramework) GetConfig() *ServerConfig {
	args := m.Called()
	return args.Get(0).(*ServerConfig)
}

// MockFrameworkPlugin 模拟框架插件用于测试 (Mock framework plugin for testing)
type MockFrameworkPlugin struct {
	name string
}

func (m *MockFrameworkPlugin) Name() string {
	return m.name
}

func (m *MockFrameworkPlugin) Version() string {
	return "1.0.0"
}

func (m *MockFrameworkPlugin) Description() string {
	return "Mock plugin for testing"
}

func (m *MockFrameworkPlugin) DefaultConfig() interface{} {
	return &ServerConfig{
		Framework: m.name,
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}
}

func (m *MockFrameworkPlugin) CreateFramework(config interface{}, services services.ServiceContainer) (WebFramework, error) {
	framework := &MockWebFramework{}
	if serverConfig, ok := config.(*ServerConfig); ok {
		framework.On("GetConfig").Return(serverConfig)
	}
	return framework, nil
}

func (m *MockFrameworkPlugin) ValidateConfig(config interface{}) error {
	return nil
}

func (m *MockFrameworkPlugin) GetConfigSchema() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"framework": map[string]interface{}{"type": "string"},
			"host":      map[string]interface{}{"type": "string"},
			"port":      map[string]interface{}{"type": "integer"},
			"mode":      map[string]interface{}{"type": "string"},
		},
	}
}

// TestNewServerManager 测试创建服务器管理器 (Test creating server manager)
func TestNewServerManager(t *testing.T) {
	framework := &MockWebFramework{}
	config := &ServerConfig{
		Framework: "test",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}

	manager := NewServerManager(framework, config)

	assert.NotNil(t, manager)
	assert.Equal(t, framework, manager.GetFramework())
	assert.Equal(t, config, manager.GetConfig())
	assert.False(t, manager.IsRunning())
}

// TestServerManager_Start 测试启动服务器 (Test starting server)
func TestServerManager_Start(t *testing.T) {
	framework := &MockWebFramework{}
	config := &ServerConfig{
		Framework: "test",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
		GracefulShutdown: GracefulShutdownConfig{
			Enabled: false,
		},
	}

	manager := NewServerManager(framework, config)
	ctx := context.Background()

	// 期望框架启动成功，返回nil表示立即启动成功 (Expect framework to start successfully, return nil means immediate success)
	framework.On("Start", ctx).Return(nil)

	// 由于Start()会阻塞，我们需要在goroutine中运行它 (Since Start() will block, we need to run it in a goroutine)
	done := make(chan error, 1)
	go func() {
		done <- manager.Start(ctx)
	}()

	// 检查服务器是否正在运行 (Check if server is running)
	time.Sleep(10 * time.Millisecond) // 给一点时间让Start方法设置running状态 (Give some time for Start method to set running state)
	assert.True(t, manager.IsRunning())
	
	// 停止服务器以使Start方法返回 (Stop server to make Start method return)
	framework.On("Stop", ctx).Return(nil)
	err := manager.Stop(ctx)
	assert.NoError(t, err)
	
	// 等待Start方法返回 (Wait for Start method to return)
	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(1 * time.Second):
		t.Fatal("Start method did not return within timeout")
	}
	
	framework.AssertExpectations(t)
}

// TestServerManager_StartAlreadyRunning 测试重复启动服务器 (Test starting server when already running)
func TestServerManager_StartAlreadyRunning(t *testing.T) {
	framework := &MockWebFramework{}
	config := &ServerConfig{
		Framework: "test",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
		GracefulShutdown: GracefulShutdownConfig{
			Enabled: false,
		},
	}

	manager := NewServerManager(framework, config)
	ctx := context.Background()

	// 第一次启动 (First start) - 在goroutine中运行因为它会阻塞 (Run in goroutine because it will block)
	framework.On("Start", ctx).Return(nil)
	done := make(chan error, 1)
	go func() {
		done <- manager.Start(ctx)
	}()

	// 等待服务器启动 (Wait for server to start)
	time.Sleep(10 * time.Millisecond)
	assert.True(t, manager.IsRunning())

	// 第二次启动应该失败 (Second start should fail)
	err := manager.Start(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already running")

	// 停止服务器清理测试 (Stop server to clean up test)
	framework.On("Stop", ctx).Return(nil)
	_ = manager.Stop(ctx)

	// 等待第一次Start返回 (Wait for first Start to return)
	select {
	case <-done:
		// Start方法正常返回 (Start method returned normally)
	case <-time.After(1 * time.Second):
		t.Fatal("Start method did not return within timeout")
	}
}

// TestServerManager_StartInvalidConfig 测试使用无效配置启动 (Test starting with invalid config)
func TestServerManager_StartInvalidConfig(t *testing.T) {
	framework := &MockWebFramework{}
	config := &ServerConfig{
		Framework: "",    // 无效框架名 (Invalid framework name)
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}

	manager := NewServerManager(framework, config)
	ctx := context.Background()

	// 配置验证失败时不会调用框架的Start方法 (Framework Start method won't be called when config validation fails)
	err := manager.Start(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid configuration")
	assert.False(t, manager.IsRunning())
}

// TestServerManager_StartFrameworkError 测试框架启动失败 (Test framework start failure)
func TestServerManager_StartFrameworkError(t *testing.T) {
	framework := &MockWebFramework{}
	config := &ServerConfig{
		Framework: "test",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
		GracefulShutdown: GracefulShutdownConfig{
			Enabled: false,
		},
	}

	manager := NewServerManager(framework, config)
	ctx := context.Background()

	// 期望框架启动失败 (Expect framework start to fail)
	framework.On("Start", ctx).Return(errors.New("framework start error"))

	err := manager.Start(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to start framework")
	assert.False(t, manager.IsRunning())
}

// TestServerManager_Stop 测试停止服务器 (Test stopping server)
func TestServerManager_Stop(t *testing.T) {
	framework := &MockWebFramework{}
	config := &ServerConfig{
		Framework: "test",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
		GracefulShutdown: GracefulShutdownConfig{
			Enabled: false,
		},
	}

	manager := NewServerManager(framework, config)
	ctx := context.Background()

	// 先启动服务器 (Start server first) - 在goroutine中运行因为它会阻塞 (Run in goroutine because it will block)
	framework.On("Start", ctx).Return(nil)
	done := make(chan error, 1)
	go func() {
		done <- manager.Start(ctx)
	}()

	// 等待服务器启动 (Wait for server to start)
	time.Sleep(10 * time.Millisecond)
	assert.True(t, manager.IsRunning())

	// 期望框架停止成功 (Expect framework to stop successfully)
	framework.On("Stop", ctx).Return(nil)

	err := manager.Stop(ctx)
	assert.NoError(t, err)
	assert.False(t, manager.IsRunning())

	// 等待Start方法返回 (Wait for Start method to return)
	select {
	case <-done:
		// Start方法正常返回 (Start method returned normally)
	case <-time.After(1 * time.Second):
		t.Fatal("Start method did not return within timeout")
	}

	framework.AssertExpectations(t)
}

// TestServerManager_StopNotRunning 测试停止未运行的服务器 (Test stopping server when not running)
func TestServerManager_StopNotRunning(t *testing.T) {
	framework := &MockWebFramework{}
	config := &ServerConfig{
		Framework: "test",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}

	manager := NewServerManager(framework, config)
	ctx := context.Background()

	err := manager.Stop(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not running")
}

// TestServerManager_StopFrameworkError 测试框架停止失败 (Test framework stop failure)
func TestServerManager_StopFrameworkError(t *testing.T) {
	framework := &MockWebFramework{}
	config := &ServerConfig{
		Framework: "test",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
		GracefulShutdown: GracefulShutdownConfig{
			Enabled: false,
		},
	}

	manager := NewServerManager(framework, config)
	ctx := context.Background()

	// 先启动服务器 (Start server first) - 在goroutine中运行因为它会阻塞 (Run in goroutine because it will block)
	framework.On("Start", ctx).Return(nil)
	done := make(chan error, 1)
	go func() {
		done <- manager.Start(ctx)
	}()

	// 等待服务器启动 (Wait for server to start)
	time.Sleep(10 * time.Millisecond)
	assert.True(t, manager.IsRunning())

	// 期望框架停止失败 (Expect framework stop to fail)
	framework.On("Stop", ctx).Return(errors.New("framework stop error"))

	err := manager.Stop(ctx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to stop framework")

	// 等待Start方法返回 (Wait for Start method to return)
	select {
	case <-done:
		// Start方法正常返回 (Start method returned normally)
	case <-time.After(1 * time.Second):
		t.Fatal("Start method did not return within timeout")
	}
}

// TestServerManager_GetFramework 测试获取框架实例 (Test getting framework instance)
func TestServerManager_GetFramework(t *testing.T) {
	framework := &MockWebFramework{}
	config := &ServerConfig{
		Framework: "test",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}

	manager := NewServerManager(framework, config)

	assert.Equal(t, framework, manager.GetFramework())
}

// TestServerManager_GetConfig 测试获取配置 (Test getting configuration)
func TestServerManager_GetConfig(t *testing.T) {
	framework := &MockWebFramework{}
	config := &ServerConfig{
		Framework: "test",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}

	manager := NewServerManager(framework, config)

	assert.Equal(t, config, manager.GetConfig())
}

// TestNewServerFactory 测试创建服务器工厂 (Test creating server factory)
func TestNewServerFactory(t *testing.T) {
	factory := NewServerFactory()

	assert.NotNil(t, factory)
	// 通过公共方法验证工厂功能而不是访问私有字段 (Verify factory functionality through public methods instead of accessing private fields)
	plugins := factory.ListPlugins()
	assert.NotNil(t, plugins)
}

// TestServerFactory_RegisterPlugin 测试注册插件 (Test registering plugin)
func TestServerFactory_RegisterPlugin(t *testing.T) {
	// 清理全局注册表状态 (Clear global registry state)
	ClearFrameworks()
	factory := NewServerFactory()
	plugin := &MockFrameworkPlugin{name: "test"}

	err := factory.RegisterPlugin(plugin)
	assert.NoError(t, err)

	plugins := factory.ListPlugins()
	assert.Contains(t, plugins, "test")
}

// TestServerFactory_UnregisterPlugin 测试注销插件 (Test unregistering plugin)
func TestServerFactory_UnregisterPlugin(t *testing.T) {
	// 清理全局注册表状态 (Clear global registry state)
	ClearFrameworks()
	
	factory := NewServerFactory()
	plugin := &MockFrameworkPlugin{name: "test"}

	// 先注册插件 (Register plugin first)
	err := factory.RegisterPlugin(plugin)
	assert.NoError(t, err)

	// 然后注销插件 (Then unregister plugin)
	err = factory.UnregisterPlugin("test")
	assert.NoError(t, err)

	plugins := factory.ListPlugins()
	assert.NotContains(t, plugins, "test")
}

// TestServerFactory_CreateServer 测试创建服务器 (Test creating server)
func TestServerFactory_CreateServer(t *testing.T) {
	// 清理全局注册表状态 (Clear global registry state)
	ClearFrameworks()
	factory := NewServerFactory()
	plugin := &MockFrameworkPlugin{name: "test"}

	// 注册插件 (Register plugin)
	err := factory.RegisterPlugin(plugin)
	assert.NoError(t, err)

	config := &ServerConfig{
		Framework: "test",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}

	manager, err := factory.CreateServer("test", config)

	assert.NoError(t, err)
	assert.NotNil(t, manager)
	assert.Equal(t, config, manager.GetConfig())
}

// TestServerFactory_CreateServerInvalidPlugin 测试使用无效插件创建服务器 (Test creating server with invalid plugin)
func TestServerFactory_CreateServerInvalidPlugin(t *testing.T) {
	factory := NewServerFactory()
	config := &ServerConfig{
		Framework: "invalid",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}

	manager, err := factory.CreateServer("invalid", config)

	assert.Error(t, err)
	assert.Nil(t, manager)
	assert.Contains(t, err.Error(), "failed to create framework instance")
}

// TestServerFactory_GetPluginInfo 测试获取插件信息 (Test getting plugin info)
func TestServerFactory_GetPluginInfo(t *testing.T) {
	// 清理全局注册表状态 (Clear global registry state)
	ClearFrameworks()
	
	factory := NewServerFactory()
	plugin := &MockFrameworkPlugin{name: "test"}

	// 注册插件 (Register plugin)
	err := factory.RegisterPlugin(plugin)
	assert.NoError(t, err)

	info, err := factory.GetPluginInfo("test")
	assert.NoError(t, err)
	assert.NotNil(t, info)
}

// TestServerFactory_GetAllPluginInfo 测试获取所有插件信息 (Test getting all plugin info)
func TestServerFactory_GetAllPluginInfo(t *testing.T) {
	// 清理全局注册表状态 (Clear global registry state)
	ClearFrameworks()
	// 清理全局注册表状态 (Clear global registry state)
	ClearFrameworks()
	
	factory := NewServerFactory()
	plugin := &MockFrameworkPlugin{name: "test"}

	// 注册插件 (Register plugin)
	err := factory.RegisterPlugin(plugin)
	assert.NoError(t, err)

	allInfo := factory.GetAllPluginInfo()
	assert.NotNil(t, allInfo)
	assert.Contains(t, allInfo, "test")
}

// TestGlobalFunctions 测试全局函数 (Test global functions)
func TestGlobalFunctions(t *testing.T) {
	// 保存原始的全局工厂 (Save original global factory)
	originalFactory := globalFactory

	// 创建新的测试工厂 (Create new test factory)
	globalFactory = NewServerFactory()

	// 测试完成后恢复原始工厂 (Restore original factory after test)
	defer func() {
		globalFactory = originalFactory
	}()

	plugin := &MockFrameworkPlugin{name: "global-test"}

	// 测试全局注册 (Test global registration)
	err := RegisterPlugin(plugin)
	assert.NoError(t, err)

	// 测试列出插件 (Test listing plugins)
	plugins := ListPlugins()
	assert.Contains(t, plugins, "global-test")

	// 测试获取插件信息 (Test getting plugin info)
	info, err := GetPluginInfo("global-test")
	assert.NoError(t, err)
	assert.NotNil(t, info)

	// 测试获取所有插件信息 (Test getting all plugin info)
	allInfo := GetAllPluginInfo()
	assert.Contains(t, allInfo, "global-test")

	// 测试创建服务器管理器 (Test creating server manager)
	config := &ServerConfig{
		Framework: "global-test",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}

	manager, err := CreateServerManager("global-test", config)
	assert.NoError(t, err)
	assert.NotNil(t, manager)

	// 测试全局注销 (Test global unregistration)
	err = UnregisterPlugin("global-test")
	assert.NoError(t, err)

	plugins = ListPlugins()
	assert.NotContains(t, plugins, "global-test")
}

// TestQuickStart 测试快速启动 (Test quick start)
func TestQuickStart(t *testing.T) {
	// 这个测试需要一个完整的插件实现，在实际环境中会更复杂
	// (This test would need a complete plugin implementation, more complex in real environment)
	
	// 保存原始的全局工厂 (Save original global factory)
	originalFactory := globalFactory

	// 创建新的测试工厂 (Create new test factory)
	globalFactory = NewServerFactory()

	// 测试完成后恢复原始工厂 (Restore original factory after test)
	defer func() {
		globalFactory = originalFactory
	}()

	config := &ServerConfig{
		Framework: "invalid",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}

	// 由于没有注册有效插件，应该返回错误 (Should return error since no valid plugin is registered)
	err := QuickStart("invalid", config)
	assert.Error(t, err)
}

// BenchmarkServerManager_Start 基准测试服务器启动 (Benchmark server start)
func BenchmarkServerManager_Start(b *testing.B) {
	framework := &MockWebFramework{}
	config := &ServerConfig{
		Framework: "test",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
		GracefulShutdown: GracefulShutdownConfig{
			Enabled: false,
		},
	}

	ctx := context.Background()
	framework.On("Start", ctx).Return(nil)
	framework.On("Stop", ctx).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager := NewServerManager(framework, config)
		_ = manager.Start(ctx)
		_ = manager.Stop(ctx)
	}
}

// BenchmarkServerFactory_CreateServer 基准测试创建服务器 (Benchmark server creation)
func BenchmarkServerFactory_CreateServer(b *testing.B) {
	factory := NewServerFactory()
	plugin := &MockFrameworkPlugin{name: "bench-test"}
	_ = factory.RegisterPlugin(plugin)

	config := &ServerConfig{
		Framework: "bench-test",
		Host:      "localhost",
		Port:      8080,
		Mode:      "test",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = factory.CreateServer("bench-test", config)
	}
} 