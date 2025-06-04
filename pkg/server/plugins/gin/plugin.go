/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Gin框架插件实现
 */

package gin

import (
	"fmt"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
)

// Plugin Gin框架插件 (Gin framework plugin)
// 实现FrameworkPlugin接口 (Implements FrameworkPlugin interface)
type Plugin struct {
	name        string
	version     string
	description string
}

// NewPlugin 创建Gin插件实例 (Create Gin plugin instance)
func NewPlugin() server.FrameworkPlugin {
	return &Plugin{
		name:        "gin",
		version:     "1.0.0",
		description: "Gin HTTP web framework plugin with high performance and minimal memory footprint",
	}
}

// Name 返回插件名称 (Return plugin name)
func (p *Plugin) Name() string {
	return p.name
}

// Version 返回插件版本 (Return plugin version)
func (p *Plugin) Version() string {
	return p.version
}

// Description 返回插件描述 (Return plugin description)
func (p *Plugin) Description() string {
	return p.description
}

// DefaultConfig 返回插件的默认配置 (Return plugin default configuration)
func (p *Plugin) DefaultConfig() interface{} {
	config := server.DefaultServerConfig()
	config.Framework = p.name
	
	// Gin特定的默认配置 (Gin-specific default configuration)
	config.Mode = "debug" // Gin默认为debug模式 (Gin defaults to debug mode)
	
	// Gin插件特定配置 (Gin plugin specific configuration)
	if config.Plugins == nil {
		config.Plugins = make(map[string]interface{})
	}
	
	config.Plugins["gin"] = map[string]interface{}{
		"trusted_proxies":           []string{},
		"use_h2c":                  false,
		"redirect_trailing_slash":   true,
		"redirect_fixed_path":       false,
		"handle_method_not_allowed": false,
		"forward_by_client_ip":      true,
		"use_raw_path":             false,
		"unescape_path_values":     true,
		"max_multipart_memory":     32 << 20, // 32MB
	}
	
	return config
}

// CreateFramework 创建框架实例 (Create framework instance)
func (p *Plugin) CreateFramework(config interface{}, serviceContainer services.ServiceContainer) (server.WebFramework, error) {
	// 验证配置 (Validate configuration)
	if err := p.ValidateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}
	
	// 转换配置类型 (Convert configuration type)
	var serverConfig *server.ServerConfig
	switch cfg := config.(type) {
	case *server.ServerConfig:
		serverConfig = cfg
	case server.ServerConfig:
		serverConfig = &cfg
	case nil:
		serverConfig = p.DefaultConfig().(*server.ServerConfig)
	default:
		return nil, fmt.Errorf("unsupported configuration type: %T", config)
	}
	
	// 确保框架名称正确 (Ensure framework name is correct)
	serverConfig.Framework = p.name
	
	// 创建Gin服务器适配器 (Create Gin server adapter)
	ginServer := NewGinServerWithServices(serverConfig, serviceContainer)
	
	return ginServer, nil
}

// ValidateConfig 验证配置 (Validate configuration)
func (p *Plugin) ValidateConfig(config interface{}) error {
	if config == nil {
		return nil // nil配置是允许的，会使用默认配置 (nil config is allowed, will use default config)
	}
	
	switch cfg := config.(type) {
	case *server.ServerConfig:
		return cfg.Validate()
	case server.ServerConfig:
		return cfg.Validate()
	default:
		return fmt.Errorf("unsupported configuration type: %T, expected *server.ServerConfig", config)
	}
}

// GetConfigSchema 获取配置模式 (Get configuration schema)
func (p *Plugin) GetConfigSchema() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"framework": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"gin"},
				"description": "Framework name, must be 'gin'",
			},
			"host": map[string]interface{}{
				"type":        "string",
				"default":     "0.0.0.0",
				"description": "Host address to listen on",
			},
			"port": map[string]interface{}{
				"type":        "integer",
				"minimum":     1,
				"maximum":     65535,
				"default":     8080,
				"description": "Port to listen on",
			},
			"mode": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"debug", "release", "test"},
				"default":     "debug",
				"description": "Running mode",
			},
			"plugins": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"gin": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"trusted_proxies": map[string]interface{}{
								"type":        "array",
								"items":       map[string]interface{}{"type": "string"},
								"description": "List of trusted proxy IP addresses",
							},
							"redirect_trailing_slash": map[string]interface{}{
								"type":        "boolean",
								"default":     true,
								"description": "Whether to redirect trailing slash",
							},
							"max_multipart_memory": map[string]interface{}{
								"type":        "integer",
								"default":     33554432, // 32MB
								"description": "Maximum memory for multipart forms",
							},
						},
					},
				},
			},
		},
		"required": []string{"framework"},
	}
}

// init 自动注册Gin插件 (Auto-register Gin plugin)
func init() {
	// 自动注册到全局注册表 (Auto-register to global registry)
	if err := server.RegisterFramework(NewPlugin()); err != nil {
		// 如果注册失败，可能是因为已经注册过了 (If registration fails, it might be already registered)
		// 这里不做处理，让用户手动处理 (Don't handle here, let user handle manually)
	}
}