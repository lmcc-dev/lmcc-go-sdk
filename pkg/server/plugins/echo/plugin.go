/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Echo框架插件注册和工厂 (Echo framework plugin registration and factory)
 */

package echo

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

// Plugin Echo框架插件 (Echo framework plugin)
type Plugin struct{}

// NewPlugin 创建Echo插件实例 (Create Echo plugin instance)
func NewPlugin() server.FrameworkPlugin {
	return &Plugin{}
}

// Name 返回插件名称 (Return plugin name)
func (p *Plugin) Name() string {
	return "echo"
}

// Version 返回插件版本 (Return plugin version)
func (p *Plugin) Version() string {
	return "v1.0.0"
}

// Description 返回插件描述 (Return plugin description)
func (p *Plugin) Description() string {
	return "Echo web framework plugin for lmcc-go-sdk (Echo Web框架插件)"
}

// DefaultConfig 返回插件的默认配置 (Return plugin default configuration)
func (p *Plugin) DefaultConfig() interface{} {
	return &server.ServerConfig{
		Framework:      "echo",
		Host:           "0.0.0.0",
		Port:           8080,
		Mode:           "debug",
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
		
		// 启用基本功能 (Enable basic features)
		CORS: server.CORSConfig{
			Enabled: true,
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders: []string{"*"},
		},
		
		Middleware: server.MiddlewareConfig{
			Logger: server.LoggerMiddlewareConfig{
				Enabled: true,
				Format:  "json",
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
}

// CreateFramework 创建Echo框架实例 (Create Echo framework instance)
func (p *Plugin) CreateFramework(config interface{}, serviceContainer services.ServiceContainer) (server.WebFramework, error) {
	// 类型断言获取服务器配置 (Type assertion to get server configuration)
	serverConfig, ok := config.(*server.ServerConfig)
	if !ok {
		return nil, serviceContainer.GetErrorHandler().New("invalid configuration type for Echo plugin")
	}

	// 验证配置 (Validate configuration)
	if err := p.ValidateConfig(serverConfig); err != nil {
		return nil, serviceContainer.GetErrorHandler().Wrap(err, "Echo plugin configuration validation failed")
	}

	// 创建Echo服务器适配器 (Create Echo server adapter)
	echoServer, err := NewEchoServer(serverConfig, serviceContainer)
	if err != nil {
		return nil, serviceContainer.GetErrorHandler().Wrap(err, "failed to create Echo server")
	}

	// 记录插件创建日志 (Log plugin creation)
	serviceContainer.GetLogger().Infow("Echo framework plugin created successfully",
		"plugin_name", p.Name(),
		"plugin_version", p.Version(),
		"server_host", serverConfig.Host,
		"server_port", serverConfig.Port,
		"server_mode", serverConfig.Mode,
	)

	return echoServer, nil
}

// ValidateConfig 验证配置 (Validate configuration)
func (p *Plugin) ValidateConfig(config interface{}) error {
	serverConfig, ok := config.(*server.ServerConfig)
	if !ok {
		return ErrInvalidConfig
	}

	// 验证框架名称 (Validate framework name)
	if serverConfig.Framework != "echo" {
		return ErrInvalidFramework
	}

	// 验证服务器配置 (Validate server configuration)
	return serverConfig.Validate()
}

// GetConfigSchema 获取配置模式 (Get configuration schema)
func (p *Plugin) GetConfigSchema() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"framework": map[string]interface{}{
				"type":    "string",
				"enum":    []string{"echo"},
				"default": "echo",
			},
			"host": map[string]interface{}{
				"type":    "string",
				"default": "0.0.0.0",
			},
			"port": map[string]interface{}{
				"type":    "integer",
				"minimum": 1,
				"maximum": 65535,
				"default": 8080,
			},
			"mode": map[string]interface{}{
				"type": "string",
				"enum": []string{"debug", "release", "test"},
				"default": "debug",
			},
		},
		"required": []string{"framework", "host", "port"},
	}
}

// 注册Echo插件到全局注册表 (Register Echo plugin to global registry)
func init() {
	if err := server.RegisterFramework(NewPlugin()); err != nil {
		panic("failed to register Echo plugin: " + err.Error())
	}
} 