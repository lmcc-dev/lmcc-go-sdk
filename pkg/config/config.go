/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package config

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	lmccerrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Note: ConfigChangeCallback, configManager, newConfigManager, RegisterCallback, notifyCallbacks, GetViperInstance
// are now defined in manager.go
// SectionChangeCallback and Manager interface are defined in types.go

// 全局配置实例 (Global configuration instance)
// 注意：这个变量将在 LoadConfig 或 LoadConfigAndWatch 成功执行后被设置。
// (Note: This variable will be set after LoadConfig or LoadConfigAndWatch executes successfully.)
var Cfg *Config

// 全局 Cfg 的互斥锁 (Mutex for global Cfg) - 用于保护对全局 Cfg 的并发访问。
// (Mutex for global Cfg - used to protect concurrent access to the global Cfg.)
var cfgMux sync.RWMutex

// updateGlobalCfg is now defined in accessors.go

// GetGlobalCfg 安全地返回全局 Cfg 变量。
// (GetGlobalCfg safely returns the global Cfg variable.)
// Returns:
//   *Config: 指向全局配置实例的指针 (如果已加载)。
//           (A pointer to the global configuration instance, if loaded.)
func GetGlobalCfg() *Config {
	cfgMux.RLock()
	defer cfgMux.RUnlock()
	return Cfg
}

// LoadConfigAndWatch 加载配置并启动监控以实现运行时更新。
// 这是推荐的主要配置加载函数。
// (LoadConfigAndWatch loads the configuration and starts watching for runtime updates.)
// (This is the recommended primary function for loading configuration.)
//
// Parameters:
//   cfg:  指向要填充的配置结构体的指针。该结构体应使用 `mapstructure` 和 `default` 标签。
//         (A pointer to the configuration struct to be populated. It should use `mapstructure` and `default` tags.)
//   opts: 一个或多个配置选项 (Option)，用于自定义加载行为 (例如，配置文件路径、环境变量前缀、热重载)。
//         (One or more configuration options (Option) to customize loading behavior (e.g., config file path, env var prefix, hot-reload).)
//
// Returns:
//   *configManager[T]: 一个配置管理器实例，可用于注册回调或获取内部 Viper 实例。
//                      (A config manager instance that can be used to register callbacks or get the internal Viper instance.)
//   error: 加载或监控过程中发生的任何错误。
//          (Any error that occurred during loading or watching.)
func LoadConfigAndWatch[T any](cfg *T, opts ...Option) (Manager, error) {
	cm := newConfigManager(cfg, opts...) // newConfigManager is defined in manager.go

	// 1. 初始化 cfg 中的 nil 指针字段 (Initialize nil pointer fields in cfg)
	// Assuming initializeNilPointers is defined elsewhere (e.g., defaults.go)
	initializeNilPointers(cm.cfg)

	// 2. 配置 Viper 从环境变量读取 (Configure Viper to read from environment variables)
	if cm.options.enableEnvVarOverride {
		replacer := strings.NewReplacer(".", "_", "-", "_")
		cm.v.SetEnvPrefix(cm.options.envPrefix)
		cm.v.SetEnvKeyReplacer(replacer)
		cm.v.AutomaticEnv()
		// Assuming bindEnvs is defined elsewhere (e.g., env.go or defaults.go)
		bindEnvs(cm.v, replacer, cm.cfg)
	}

	// 3. 设置并读取配置文件 (Set and read the config file)
	configFileUsed := ""
	var keysFromConfigFile map[string]bool // 记录配置文件中实际存在的键 (Record keys actually present in config file)
	if cm.options.configFilePath != "" {
		cm.v.SetConfigFile(cm.options.configFilePath)
		if cm.options.configFileType == "" {
			ext := filepath.Ext(cm.options.configFilePath)
			if len(ext) > 1 {
				configType := strings.ToLower(ext[1:])
				cm.v.SetConfigType(configType)
			} else {
				log.Printf("Warning: Could not infer config type from file extension '%s'...", cm.options.configFilePath)
			}
		} else {
			cm.v.SetConfigType(strings.ToLower(cm.options.configFileType))
		}

		err := cm.v.ReadInConfig()
		if err != nil {
			var configFileNotFoundError viper.ConfigFileNotFoundError
			if errors.As(err, &configFileNotFoundError) || os.IsNotExist(err) {
				// 文件未找到，也应该是一个错误，而不仅仅是日志
				// (File not found should also be an error, not just a log)
				return nil, lmccerrors.WithCode(
					lmccerrors.Wrapf(err, "config file '%s' not found", cm.options.configFilePath),
					lmccerrors.ErrConfigFileRead,
				)
			} else {
				return nil, lmccerrors.WithCode(
					lmccerrors.Wrapf(err, "failed to read config file '%s'", cm.options.configFilePath),
					lmccerrors.ErrConfigFileRead,
				)
			}
		} else {
			configFileUsed = cm.options.configFilePath
			log.Printf("Info: Successfully read config file '%s'.", configFileUsed)
			
			// 记录配置文件中实际存在的键 (Record keys actually present in config file)
			keysFromConfigFile = flattenViperKeys(cm.v.AllSettings())
		}
	} else {
		log.Println("Info: No config file path provided...")
		keysFromConfigFile = make(map[string]bool) // 空映射 (Empty map)
	}

	// 4. 从结构体标签设置 Viper 默认值 (Set Viper defaults from struct tags)
	// Assuming setDefaultsFromTags is defined elsewhere (e.g., defaults.go)
	if err := setDefaultsFromTags(cm.v, cm.cfg, ""); err != nil {
		return nil, lmccerrors.WithCode(
			lmccerrors.Wrap(err, "failed to set defaults from struct tags"),
			lmccerrors.ErrConfigSetup,
		)
	}

	// 5. 将 Viper 配置解组到结构体中 (Unmarshal the Viper config into the struct)
	decoderConfig := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		),
		WeaklyTypedInput: true,
		TagName:          "mapstructure",
		Result:           cm.cfg,
		Squash:           true,
	}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return nil, lmccerrors.WithCode(
			lmccerrors.Wrap(err, "failed to create mapstructure decoder"),
			lmccerrors.ErrConfigSetup,
		)
	}
	if err := decoder.Decode(cm.v.AllSettings()); err != nil {
		return nil, lmccerrors.WithCode(
			lmccerrors.Wrap(err, "failed to unmarshal config from mapstructure"),
			lmccerrors.ErrConfigSetup,
		)
	}

	// 6. 在解码后应用默认值到零值字段 (Apply defaults to zero-value fields after decoding)
	// 使用改进版本的函数，它能够区分显式设置的值和真正的零值
	// (Use improved version of the function that can distinguish explicitly set values from true zero values)
	if err := applyDefaultsToZeroFieldsWithViper(cm.cfg, cm.v, keysFromConfigFile); err != nil {
		return nil, lmccerrors.WithCode(
			lmccerrors.Wrap(err, "failed to apply defaults to zero fields after initial load"),
			lmccerrors.ErrConfigSetup,
		)
	}

	// 7. 配置并启动监控（如果启用）(Configure and start watching if enabled)
	if cm.options.enableHotReload && configFileUsed != "" {
		cm.v.WatchConfig()
		// 使用 OnConfigChange 来处理 Viper 内部的文件变更通知
		// (Use OnConfigChange to handle Viper's internal file change notifications)
		cm.v.OnConfigChange(func(e fsnotify.Event) {
			// 检查事件类型，避免不必要的重载（例如 CHMOD）
			// Check event type to avoid unnecessary reloads (e.g., CHMOD)
			if e.Op&fsnotify.Write != fsnotify.Write && e.Op&fsnotify.Create != fsnotify.Create {
				log.Printf("Info: Config watcher received non-write/create event (%s), skipping reload.", e.Op)
				return
			}

			log.Printf("Config file changed: %s. Reloading...", e.Name)

			// 重新读取配置 (Re-read the config)
			if errRead := cm.v.ReadInConfig(); errRead != nil {
				// 如果文件在监控期间被删除，ReadInConfig 会报错，这是可能的场景
				// (If the file is deleted during watch, ReadInConfig will error, which is possible)
				log.Printf("Error reading config during hot reload: %v", errRead)
				// Consider if we should reset config or keep old one? Keep old one for now.
				return // Skip update and callbacks if re-read fails
			}

			// 重新解码配置到 cm.cfg (Re-decode the configuration into cm.cfg)
			newDecoderConfig := &mapstructure.DecoderConfig{
				WeaklyTypedInput: true,
				TagName:          "mapstructure",
				Result:           cm.cfg, // Update the existing config object
				Squash:           true,
				DecodeHook: mapstructure.ComposeDecodeHookFunc(
					mapstructure.StringToTimeDurationHookFunc(),
					mapstructure.StringToSliceHookFunc(","),
				),
			}
			newDecoder, errDecoder := mapstructure.NewDecoder(newDecoderConfig)
			if errDecoder != nil {
				log.Printf("Error creating decoder during hot reload: %v", errDecoder)
				return // Skip notifying callbacks on decoder error
			}

			if errUnmarshal := newDecoder.Decode(cm.v.AllSettings()); errUnmarshal != nil {
				log.Printf("Error re-unmarshalling config during hot reload: %v", errUnmarshal)
				return // Skip notifying callbacks on unmarshal error
			}

			// 在热重载解码后应用默认值 (Apply defaults after hot reload decoding)
			// 使用改进版本的函数，它能够区分显式设置的值和真正的零值
			// (Use improved version of the function that can distinguish explicitly set values from true zero values)
			// 重新构建配置文件键映射 (Rebuild config file keys map)
			hotReloadKeysFromConfigFile := flattenViperKeys(cm.v.AllSettings())
			if errApplyDefaults := applyDefaultsToZeroFieldsWithViper(cm.cfg, cm.v, hotReloadKeysFromConfigFile); errApplyDefaults != nil {
				log.Printf("Error applying defaults to zero fields during hot reload: %v", errApplyDefaults)
				// Decide if we should skip callbacks or proceed. For now, proceed.
			}

			log.Println("Config reloaded successfully.")
			// 调用 accessors.go 中的 updateGlobalCfg (Call updateGlobalCfg from accessors.go)
			updateGlobalCfg(cm.cfg)

			// 通知所有注册的回调 (Notify all registered callbacks)
			cm.notifyCallbacks() // notifyCallbacks is defined in manager.go
		})
		log.Printf("Hot reload enabled for config file: %s", configFileUsed)
	} else if cm.options.enableHotReload {
		log.Println("Warning: Hot reload enabled but no config file was used, watcher not started.")
	}

	// 首次加载后更新全局 Cfg 变量 (Update the global Cfg variable after initial load)
	// 调用 accessors.go 中的 updateGlobalCfg (Call updateGlobalCfg from accessors.go)
	updateGlobalCfg(cm.cfg)

	return cm, nil
}

// LoadConfig 是一个简化的包装器，用于加载配置（不带热重载监控）。
// (LoadConfig is a simplified wrapper for loading configuration without hot-reload watching.)
// 推荐使用 LoadConfigAndWatch 来获取完整的运行时更新功能。
// (Using LoadConfigAndWatch is recommended for full runtime update capabilities.)
//
// Parameters:
//   cfg:  指向要填充的配置结构体的指针。该结构体应使用 `mapstructure` 和 `default` 标签。
//         (A pointer to the configuration struct to be populated. It should use `mapstructure` and `default` tags.)
//   opts: 一个或多个配置选项 (Option)，用于自定义加载行为。注意：WithHotReload 选项会被忽略。
//         (One or more configuration options (Option) to customize loading behavior. Note: The WithHotReload option will be ignored.)
//
// Returns:
//   error: 加载过程中发生的任何错误。
//          (Any error that occurred during loading.)
func LoadConfig[T any](cfg *T, opts ...Option) error {
	// Filter out the WithHotReload option if present, as this function doesn't support it.
	// (如果存在 WithHotReload 选项，则将其过滤掉，因为此函数不支持它。)
	filteredOpts := []Option{}
	for _, opt := range opts {
		// We need a way to identify the hot reload option. Assuming it sets a specific flag.
		// Let's temporarily assume WithHotReload sets options.enableHotReload to true.
		// A better approach might be to have Option return an identifier or use type assertion.
		tempOpts := defaultOptions // Apply option to temp struct to check its effect
		opt(&tempOpts)
		if !tempOpts.enableHotReload {
			filteredOpts = append(filteredOpts, opt)
		} else {
			log.Println("Info: LoadConfig ignores the WithHotReload option.")
		}
	}
	_, err := LoadConfigAndWatch(cfg, filteredOpts...)
	return err
}
