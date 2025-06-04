/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: 服务器工厂和管理功能
 */

package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ServerManager 服务器管理器 (Server manager)
// 提供服务器生命周期管理功能 (Provides server lifecycle management functionality)
type ServerManager struct {
	framework WebFramework
	config    *ServerConfig
	server    *http.Server
	running   bool
}

// NewServerManager 创建服务器管理器 (Create server manager)
func NewServerManager(framework WebFramework, config *ServerConfig) *ServerManager {
	return &ServerManager{
		framework: framework,
		config:    config,
		running:   false,
	}
}

// Start 启动服务器 (Start server)
func (sm *ServerManager) Start(ctx context.Context) error {
	if sm.running {
		return fmt.Errorf("server is already running")
	}
	
	// 验证配置 (Validate configuration)
	if err := sm.config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}
	
	// 启动框架 (Start framework)
	if err := sm.framework.Start(ctx); err != nil {
		return fmt.Errorf("failed to start framework: %w", err)
	}
	
	sm.running = true
	
	// 如果启用了优雅关闭，设置信号处理 (If graceful shutdown is enabled, set up signal handling)
	if sm.config.GracefulShutdown.Enabled {
		go sm.handleGracefulShutdown()
	}
	
	return nil
}

// Stop 停止服务器 (Stop server)
func (sm *ServerManager) Stop(ctx context.Context) error {
	if !sm.running {
		return fmt.Errorf("server is not running")
	}
	
	// 停止框架 (Stop framework)
	if err := sm.framework.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop framework: %w", err)
	}
	
	sm.running = false
	return nil
}

// IsRunning 检查服务器是否正在运行 (Check if server is running)
func (sm *ServerManager) IsRunning() bool {
	return sm.running
}

// GetFramework 获取底层框架实例 (Get underlying framework instance)
func (sm *ServerManager) GetFramework() WebFramework {
	return sm.framework
}

// GetConfig 获取服务器配置 (Get server configuration)
func (sm *ServerManager) GetConfig() *ServerConfig {
	return sm.config
}

// handleGracefulShutdown 处理优雅关闭 (Handle graceful shutdown)
func (sm *ServerManager) handleGracefulShutdown() {
	// 创建信号通道 (Create signal channel)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	// 等待信号 (Wait for signal)
	<-sigChan
	
	fmt.Println("Received shutdown signal, starting graceful shutdown...")
	
	// 创建关闭上下文 (Create shutdown context)
	ctx, cancel := context.WithTimeout(context.Background(), sm.config.GracefulShutdown.Timeout)
	defer cancel()
	
	// 等待一段时间让正在处理的请求完成 (Wait for ongoing requests to complete)
	if sm.config.GracefulShutdown.WaitTime > 0 {
		fmt.Printf("Waiting %v for ongoing requests to complete...\n", sm.config.GracefulShutdown.WaitTime)
		time.Sleep(sm.config.GracefulShutdown.WaitTime)
	}
	
	// 停止服务器 (Stop server)
	if err := sm.Stop(ctx); err != nil {
		fmt.Printf("Error during graceful shutdown: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println("Server shutdown completed")
	os.Exit(0)
}

// ServerFactory 服务器工厂 (Server factory)
// 提供创建和管理服务器实例的功能 (Provides functionality to create and manage server instances)
type ServerFactory struct {
	registry *PluginRegistry
}

// NewServerFactory 创建服务器工厂 (Create server factory)
func NewServerFactory() *ServerFactory {
	return &ServerFactory{
		registry: NewPluginRegistry(),
	}
}

// RegisterPlugin 注册插件 (Register plugin)
func (sf *ServerFactory) RegisterPlugin(plugin FrameworkPlugin) error {
	return sf.registry.Register(plugin)
}

// UnregisterPlugin 注销插件 (Unregister plugin)
func (sf *ServerFactory) UnregisterPlugin(name string) error {
	return sf.registry.Unregister(name)
}

// CreateServer 创建服务器实例 (Create server instance)
func (sf *ServerFactory) CreateServer(frameworkName string, config *ServerConfig) (*ServerManager, error) {
	// 创建框架实例 (Create framework instance)
	framework, err := sf.registry.CreateServer(frameworkName, config, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create framework instance: %w", err)
	}
	
	// 创建服务器管理器 (Create server manager)
	manager := NewServerManager(framework, config)
	
	return manager, nil
}

// ListPlugins 列出所有插件 (List all plugins)
func (sf *ServerFactory) ListPlugins() []string {
	return sf.registry.List()
}

// GetPluginInfo 获取插件信息 (Get plugin information)
func (sf *ServerFactory) GetPluginInfo(name string) (map[string]string, error) {
	return sf.registry.GetPluginInfo(name)
}

// GetAllPluginInfo 获取所有插件信息 (Get all plugin information)
func (sf *ServerFactory) GetAllPluginInfo() map[string]map[string]string {
	return sf.registry.GetAllPluginInfo()
}

// 全局服务器工厂实例 (Global server factory instance)
var globalFactory = NewServerFactory()

// 全局函数 - 使用全局工厂 (Global functions - using global factory)

// RegisterPlugin 注册插件到全局工厂 (Register plugin to global factory)
func RegisterPlugin(plugin FrameworkPlugin) error {
	return globalFactory.RegisterPlugin(plugin)
}

// UnregisterPlugin 从全局工厂注销插件 (Unregister plugin from global factory)
func UnregisterPlugin(name string) error {
	return globalFactory.UnregisterPlugin(name)
}

// CreateServerManager 使用全局工厂创建服务器管理器 (Create server manager using global factory)
func CreateServerManager(frameworkName string, config *ServerConfig) (*ServerManager, error) {
	return globalFactory.CreateServer(frameworkName, config)
}

// ListPlugins 列出全局工厂中的所有插件 (List all plugins in global factory)
func ListPlugins() []string {
	return globalFactory.ListPlugins()
}

// GetPluginInfo 获取全局工厂中的插件信息 (Get plugin information from global factory)
func GetPluginInfo(name string) (map[string]string, error) {
	return globalFactory.GetPluginInfo(name)
}

// GetAllPluginInfo 获取全局工厂中的所有插件信息 (Get all plugin information from global factory)
func GetAllPluginInfo() map[string]map[string]string {
	return globalFactory.GetAllPluginInfo()
}

// QuickStart 快速启动服务器 (Quick start server)
// 这是一个便捷函数，用于快速创建和启动服务器 (This is a convenience function for quickly creating and starting a server)
func QuickStart(frameworkName string, config *ServerConfig) error {
	// 如果没有提供配置，使用默认配置 (If no config provided, use default config)
	if config == nil {
		config = DefaultServerConfig()
		if frameworkName != "" {
			config.Framework = frameworkName
		}
	}
	
	// 创建服务器管理器 (Create server manager)
	manager, err := CreateServerManager(frameworkName, config)
	if err != nil {
		return fmt.Errorf("failed to create server manager: %w", err)
	}
	
	// 启动服务器 (Start server)
	ctx := context.Background()
	if err := manager.Start(ctx); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	
	fmt.Printf("Server started successfully on %s\n", config.GetAddress())
	fmt.Printf("Framework: %s\n", config.Framework)
	fmt.Printf("Mode: %s\n", config.Mode)
	
	// 如果启用了优雅关闭，等待关闭信号 (If graceful shutdown is enabled, wait for shutdown signal)
	if config.GracefulShutdown.Enabled {
		// 阻塞等待关闭信号 (Block waiting for shutdown signal)
		select {}
	}
	
	return nil
}