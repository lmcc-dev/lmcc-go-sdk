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

// WithConfigFile 返回一个 Option，用于设置要加载的配置文件的路径和可选的文件类型。
// 如果 fileType 为空字符串，将尝试从文件扩展名推断类型。
// (WithConfigFile returns an Option to set the path and optional type of the configuration file to load.)
// (If fileType is an empty string, the type will be inferred from the file extension.)
// Parameters:
//   path: 配置文件的完整路径。
//         (The full path to the configuration file.)
//   fileType: 配置文件的类型 (例如 "yaml", "json", "toml")。如果为空则自动推断。
//             (The type of the configuration file (e.g., "yaml", "json", "toml"). Auto-inferred if empty.)
// Returns:
//   Option: 应用此配置的 Option 函数。
//           (The Option function to apply this configuration.)
func WithConfigFile(path string, fileType string) Option {
	return func(o *Options) {
		o.configFilePath = path
		o.configFileType = fileType
	}
}

// WithEnvPrefix 返回一个 Option，用于设置查找环境变量时使用的前缀。
// 例如，如果前缀为 "APP"，则会查找如 APP_SERVER_PORT 这样的变量。
// (WithEnvPrefix returns an Option to set the prefix used when looking up environment variables.)
// (For example, if the prefix is "APP", variables like APP_SERVER_PORT will be looked up.)
// Parameters:
//   prefix: 要使用的环境变量前缀。
//           (The environment variable prefix to use.)
// Returns:
//   Option: 应用此配置的 Option 函数。
//           (The Option function to apply this configuration.)
func WithEnvPrefix(prefix string) Option {
	return func(o *Options) {
		if prefix != "" {
			o.envPrefix = prefix
		}
	}
}

// WithEnvVarOverride 返回一个 Option，用于控制是否允许环境变量覆盖配置文件或默认值。
// (WithEnvVarOverride returns an Option to control whether environment variables are allowed to override config file or default values.)
// Parameters:
//   enable: true 表示启用覆盖，false 表示禁用。默认为 true。
//           (true to enable override, false to disable. Defaults to true.)
// Returns:
//   Option: 应用此配置的 Option 函数。
//           (The Option function to apply this configuration.)
func WithEnvVarOverride(enable bool) Option {
	return func(o *Options) {
		o.enableEnvVarOverride = enable
	}
}

// WithHotReload 返回一个 Option，用于启用或禁用配置文件的热重载功能。
// 如果启用，当配置文件发生更改时，配置将自动重新加载，并触发已注册的回调。
// (WithHotReload returns an Option to enable or disable the hot-reload feature for the configuration file.)
// (If enabled, the configuration will be automatically reloaded upon file changes, and registered callbacks will be triggered.)
// Parameters:
//   enable: true 表示启用热重载，false 表示禁用。默认为 false。
//           (true to enable hot-reload, false to disable. Defaults to false.)
// Returns:
//   Option: 应用此配置的 Option 函数。
//           (The Option function to apply this configuration.)
func WithHotReload(enable bool) Option {
	return func(o *Options) {
		o.enableHotReload = enable
	}
}
