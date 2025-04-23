/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package config

// Options 结构体定义了配置加载的可选参数
// (Options struct defines optional parameters for config loading)
type Options struct {
	configFilePath       string // 配置文件路径 (Configuration file path)
	configFileType       string // 配置文件类型 (Configuration file type)
	envPrefix            string // 环境变量前缀 (Environment variable prefix)
	enableEnvVarOverride bool   // 是否启用环境变量覆盖 (Whether to enable environment variable override)
	enableHotReload      bool   // 是否启用热重载 (Whether to enable hot reload)
}

// Option 是一个函数类型，用于修改 Options 结构体
// (Option is a function type used to modify the Options struct)
type Option func(*Options)

// 默认配置选项 (Default configuration options)
var defaultOptions = Options{
	configFilePath:       "",     // 默认无配置文件 (No config file by default)
	configFileType:       "",     // 默认无配置文件类型 (No config file type by default)
	envPrefix:            "LMCC", // 默认前缀 (Default prefix)
	enableEnvVarOverride: true,   // 默认启用环境变量覆盖 (Enable env var override by default)
	enableHotReload:      false,  // 默认禁用热重载 (Disable hot reload by default)
}

// WithConfigFile 设置配置文件路径和类型 (Sets the config file path and type)
func WithConfigFile(path string, fileType string) Option {
	return func(o *Options) {
		o.configFilePath = path
		o.configFileType = fileType
	}
}

// WithEnvPrefix 设置环境变量前缀 (Sets the environment variable prefix)
func WithEnvPrefix(prefix string) Option {
	return func(o *Options) {
		if prefix != "" {
			o.envPrefix = prefix
		}
	}
}

// WithEnvVarOverride 启用或禁用环境变量覆盖 (Enables or disables environment variable override)
func WithEnvVarOverride(enable bool) Option {
	return func(o *Options) {
		o.enableEnvVarOverride = enable
	}
}

// WithHotReload 启用或禁用配置热重载 (Enables or disables config hot reload)
func WithHotReload(enable bool) Option {
	return func(o *Options) {
		o.enableHotReload = enable
	}
}
