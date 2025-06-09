/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: 框架插件注册和管理机制
 */

package server

import (
	"fmt"
	"sort"
	"sync"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
)

// PluginRegistry 插件注册表 (Plugin registry)
// 管理所有已注册的框架插件 (Manages all registered framework plugins)
type PluginRegistry struct {
	plugins        map[string]FrameworkPlugin
	defaultPlugin  string
	mutex          sync.RWMutex
}

// 全局插件注册表实例 (Global plugin registry instance)
var globalRegistry = NewPluginRegistry()

// NewPluginRegistry 创建新的插件注册表 (Create new plugin registry)
func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{
		plugins: make(map[string]FrameworkPlugin),
	}
}

// Register 注册框架插件 (Register framework plugin)
func (r *PluginRegistry) Register(plugin FrameworkPlugin) error {
	if plugin == nil {
		return fmt.Errorf("plugin cannot be nil")
	}
	
	name := plugin.Name()
	if name == "" {
		return fmt.Errorf("plugin name cannot be empty")
	}
	
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// 检查是否已注册 (Check if already registered)
	if _, exists := r.plugins[name]; exists {
		return fmt.Errorf("plugin '%s' is already registered", name)
	}
	
	r.plugins[name] = plugin
	
	// 如果是第一个插件，设为默认插件 (If it's the first plugin, set as default)
	if r.defaultPlugin == "" {
		r.defaultPlugin = name
	}
	
	return nil
}

// Unregister 注销框架插件 (Unregister framework plugin)
func (r *PluginRegistry) Unregister(name string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.plugins[name]; !exists {
		return fmt.Errorf("plugin '%s' is not registered", name)
	}
	
	delete(r.plugins, name)
	
	// 如果删除的是默认插件，重新选择默认插件 (If removing default plugin, reselect default)
	if r.defaultPlugin == name {
		r.defaultPlugin = ""
		for pluginName := range r.plugins {
			r.defaultPlugin = pluginName
			break
		}
	}
	
	return nil
}

// Get 获取指定名称的插件 (Get plugin by name)
func (r *PluginRegistry) Get(name string) (FrameworkPlugin, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	plugin, exists := r.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin '%s' is not registered", name)
	}
	
	return plugin, nil
}

// GetDefault 获取默认插件 (Get default plugin)
func (r *PluginRegistry) GetDefault() (FrameworkPlugin, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	if r.defaultPlugin == "" {
		return nil, fmt.Errorf("no default plugin available")
	}
	
	return r.plugins[r.defaultPlugin], nil
}

// SetDefault 设置默认插件 (Set default plugin)
func (r *PluginRegistry) SetDefault(name string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	if _, exists := r.plugins[name]; !exists {
		return fmt.Errorf("plugin '%s' is not registered", name)
	}
	
	r.defaultPlugin = name
	return nil
}

// List 列出所有已注册的插件名称 (List all registered plugin names)
func (r *PluginRegistry) List() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	names := make([]string, 0, len(r.plugins))
	for name := range r.plugins {
		names = append(names, name)
	}
	
	sort.Strings(names)
	return names
}

// GetPluginInfo 获取插件信息 (Get plugin information)
func (r *PluginRegistry) GetPluginInfo(name string) (map[string]string, error) {
	plugin, err := r.Get(name)
	if err != nil {
		return nil, err
	}
	
	return map[string]string{
		"name":        plugin.Name(),
		"version":     plugin.Version(),
		"description": plugin.Description(),
	}, nil
}

// GetAllPluginInfo 获取所有插件信息 (Get all plugin information)
func (r *PluginRegistry) GetAllPluginInfo() map[string]map[string]string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	info := make(map[string]map[string]string)
	for name, plugin := range r.plugins {
		info[name] = map[string]string{
			"name":        plugin.Name(),
			"version":     plugin.Version(),
			"description": plugin.Description(),
		}
	}
	
	return info
}

// Clear 清除所有已注册的插件 (Clear all registered plugins)
func (r *PluginRegistry) Clear() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.plugins = make(map[string]FrameworkPlugin)
	r.defaultPlugin = ""
}

// CreateServer 创建服务器实例 (Create server instance)
func (r *PluginRegistry) CreateServer(frameworkName string, config *ServerConfig, serviceContainer services.ServiceContainer) (WebFramework, error) {
	var plugin FrameworkPlugin
	var err error
	
	if frameworkName == "" {
		// 使用默认插件 (Use default plugin)
		plugin, err = r.GetDefault()
		if err != nil {
			return nil, fmt.Errorf("failed to get default plugin: %w", err)
		}
	} else {
		// 使用指定插件 (Use specified plugin)
		plugin, err = r.Get(frameworkName)
		if err != nil {
			return nil, fmt.Errorf("failed to get plugin '%s': %w", frameworkName, err)
		}
	}
	
	// 如果没有提供配置，使用插件的默认配置 (If no config provided, use plugin's default config)
	var pluginConfig interface{}
	if config == nil {
		pluginConfig = plugin.DefaultConfig()
	} else {
		pluginConfig = config
	}
	
	// 验证配置 (Validate configuration)
	if err := plugin.ValidateConfig(pluginConfig); err != nil {
		return nil, fmt.Errorf("invalid configuration for plugin '%s': %w", plugin.Name(), err)
	}
	
	// 如果没有提供服务容器，创建默认的 (If no service container provided, create default one)
	if serviceContainer == nil {
		serviceContainer = services.NewServiceContainerWithDefaults()
	}
	
	return plugin.CreateFramework(pluginConfig, serviceContainer)
}

// 全局函数 - 使用全局注册表 (Global functions - using global registry)

// RegisterFramework 注册框架插件到全局注册表 (Register framework plugin to global registry)
func RegisterFramework(plugin FrameworkPlugin) error {
	return globalRegistry.Register(plugin)
}

// UnregisterFramework 从全局注册表注销框架插件 (Unregister framework plugin from global registry)
func UnregisterFramework(name string) error {
	return globalRegistry.Unregister(name)
}

// GetFramework 从全局注册表获取指定框架插件 (Get specified framework plugin from global registry)
func GetFramework(name string) (FrameworkPlugin, error) {
	return globalRegistry.Get(name)
}

// GetDefaultFramework 从全局注册表获取默认框架插件 (Get default framework plugin from global registry)
func GetDefaultFramework() (FrameworkPlugin, error) {
	return globalRegistry.GetDefault()
}

// SetDefaultFramework 设置全局默认框架插件 (Set global default framework plugin)
func SetDefaultFramework(name string) error {
	return globalRegistry.SetDefault(name)
}

// ListFrameworks 列出全局注册表中的所有框架插件 (List all framework plugins in global registry)
func ListFrameworks() []string {
	return globalRegistry.List()
}

// GetFrameworkInfo 获取指定框架插件信息 (Get specified framework plugin information)
func GetFrameworkInfo(name string) (map[string]string, error) {
	return globalRegistry.GetPluginInfo(name)
}

// GetAllFrameworkInfo 获取所有框架插件信息 (Get all framework plugin information)
func GetAllFrameworkInfo() map[string]map[string]string {
	return globalRegistry.GetAllPluginInfo()
}

// CreateServer 使用全局注册表创建服务器实例 (Create server instance using global registry)
func CreateServer(frameworkName string, config *ServerConfig, serviceContainer services.ServiceContainer) (WebFramework, error) {
	return globalRegistry.CreateServer(frameworkName, config, serviceContainer)
}

// ClearFrameworks 清除全局注册表中的所有框架插件 (Clear all framework plugins in global registry)  
func ClearFrameworks() {
	globalRegistry.Clear()
}