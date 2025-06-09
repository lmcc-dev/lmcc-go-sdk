/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Fiber框架插件注册和工厂 (Fiber framework plugin registration and factory)
 */

package fiber

import (
	"errors"
	"time"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
)

// 预定义错误 (Predefined errors)
var (
	ErrInvalidConfig    = errors.New("invalid configuration type")
	ErrInvalidFramework = errors.New("invalid framework name")
)

// Plugin Fiber框架插件 (Fiber framework plugin)
type Plugin struct{}

// NewPlugin 创建Fiber插件实例 (Create Fiber plugin instance)
func NewPlugin() server.FrameworkPlugin {
	return &Plugin{}
}

// Name 返回插件名称 (Return plugin name)
func (p *Plugin) Name() string {
	return "fiber"
}

// Version 返回插件版本 (Return plugin version)
func (p *Plugin) Version() string {
	return "v1.0.0"
}

// Description 返回插件描述 (Return plugin description)
func (p *Plugin) Description() string {
	return "Fiber web framework plugin for lmcc-go-sdk (Fiber Web框架插件)"
}

// DefaultConfig 返回插件的默认配置 (Return plugin default configuration)
func (p *Plugin) DefaultConfig() interface{} {
	return &server.ServerConfig{
		Framework:      "fiber",
		Host:           "localhost",
		Port:           8080,
		Mode:           "debug",
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
		CORS: server.CORSConfig{
			Enabled:      false,
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders: []string{"*"},
		},
		Middleware: server.MiddlewareConfig{
			Logger: server.LoggerMiddlewareConfig{
				Enabled: true,
			},
			Recovery: server.RecoveryMiddlewareConfig{
				Enabled:    true,
				PrintStack: true,
			},
		},
		GracefulShutdown: server.GracefulShutdownConfig{
			Enabled: true,
			Timeout: 30 * time.Second,
		},
	}
}

// CreateFramework 创建框架实例 (Create framework instance)
func (p *Plugin) CreateFramework(config interface{}, serviceContainer services.ServiceContainer) (server.WebFramework, error) {
	// 验证配置类型 (Validate configuration type)
	serverConfig, ok := config.(*server.ServerConfig)
	if !ok {
		return nil, ErrInvalidConfig
	}

	// 验证配置 (Validate configuration)
	if err := p.ValidateConfig(serverConfig); err != nil {
		return nil, err
	}

	// 创建Fiber服务器实例 (Create Fiber server instance)
	fiberServer, err := NewFiberServer(serverConfig, serviceContainer)
	if err != nil {
		return nil, err
	}

	// 记录插件创建成功 (Log plugin creation success)
	logger := serviceContainer.GetLogger()
	if logger != nil {
		logger.Infow("Fiber framework plugin created successfully",
			"plugin_name", p.Name(),
			"plugin_version", p.Version(),
			"server_host", serverConfig.Host,
			"server_port", serverConfig.Port,
			"server_mode", serverConfig.Mode,
		)
	}

	return fiberServer, nil
}

// ValidateConfig 验证配置 (Validate configuration)
func (p *Plugin) ValidateConfig(config interface{}) error {
	serverConfig, ok := config.(*server.ServerConfig)
	if !ok {
		return ErrInvalidConfig
	}

	// 验证框架名称 (Validate framework name)
	if serverConfig.Framework != "fiber" {
		return ErrInvalidFramework
	}

	// 使用ServerConfig的Validate方法进行验证 (Use ServerConfig's Validate method for validation)
	return serverConfig.Validate()
}

// GetConfigSchema 获取配置模式 (Get configuration schema)
func (p *Plugin) GetConfigSchema() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"framework": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"fiber"},
				"description": "Web framework name (must be 'fiber')",
			},
			"host": map[string]interface{}{
				"type":        "string",
				"default":     "localhost",
				"description": "Server host address",
			},
			"port": map[string]interface{}{
				"type":        "integer",
				"minimum":     1,
				"maximum":     65535,
				"default":     8080,
				"description": "Server port number",
			},
			"mode": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"debug", "release"},
				"default":     "debug",
				"description": "Server running mode",
			},
		},
		"required":             []string{"framework"},
		"additionalProperties": true,
	}
}

// init 自动注册Fiber插件 (Auto-register Fiber plugin)
func init() {
	_ = server.RegisterFramework(NewPlugin())
} 