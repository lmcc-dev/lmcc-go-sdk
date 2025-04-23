/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	// Needed by defaults.go and accessors.go
	// "strconv" // Needed by defaults.go
	"strings"
	// Needed by types.go, defaults.go, accessors.go

	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// 全局配置实例 (Global configuration instance)
// Type Config is defined in types.go
// Getter functions (like GetString, GetInt) are defined in accessors.go
// Option types and With... functions are defined in options.go
// Default handling helpers (like initializeNilPointers, applyDefaultsFromStructTags) are defined in defaults.go
var Cfg *Config // This will be set after LoadConfig is called.

// LoadConfig 加载配置 (Load configuration)
// 现在的加载顺序：
// 1. 应用选项
// 2. 初始化 Nil 指针
// 3. 应用结构体标签默认值
// 4. 配置 Viper (Env, 文件)
// 5. 读取配置文件
// 6. Viper Unmarshal
// 7. 配置热重载
// 8. 更新全局 Cfg
// (Current loading order:
// 1. Apply options
// 2. Initialize nil pointers
// 3. Apply struct tag defaults
// 4. Configure Viper (Env, file)
// 5. Read config file
// 6. Viper Unmarshal
// 7. Configure hot reload
// 8. Update global Cfg)
func LoadConfig[T any](cfg *T, opts ...Option) error {
	// 应用默认选项和用户提供的选项 (Apply default options and user-provided options)
	options := defaultOptions // from options.go
	for _, opt := range opts {
		opt(&options)
	}

	// --- 默认值处理提前 --- (Default value handling moved earlier)
	// 1. 初始化 cfg 中的 nil 指针字段 (Initialize nil pointer fields in cfg)
	initializeNilPointers(cfg) // from defaults.go

	// --- Viper 配置 --- (Viper configuration)
	v := viper.New()

	// 3. 配置 Viper 从环境变量读取 (Configure Viper to read from environment variables)
	if options.enableEnvVarOverride {
		replacer := strings.NewReplacer(".", "_", "-", "_")
		v.SetEnvPrefix(options.envPrefix)
		v.SetEnvKeyReplacer(replacer)
		v.AutomaticEnv()

		// 显式绑定环境变量 (Explicitly bind environment variables)
		bindEnvs(v, replacer, cfg) // <-- 恢复调用
	}

	// 4. 设置并读取配置文件 (Set and read the config file)
	if options.configFilePath != "" {
		v.SetConfigFile(options.configFilePath)
		if options.configFileType == "" {
			ext := filepath.Ext(options.configFilePath)
			if len(ext) > 1 {
				configType := strings.ToLower(ext[1:])
				v.SetConfigType(configType)
			} else {
				log.Printf("Warning: Could not infer config type from file extension '%s'...", options.configFilePath)
			}
		} else {
			v.SetConfigType(strings.ToLower(options.configFileType))
		}

		err := v.ReadInConfig()
		if err != nil {
			var configFileNotFoundError viper.ConfigFileNotFoundError
			if errors.As(err, &configFileNotFoundError) || os.IsNotExist(err) {
				log.Printf("Info: Config file '%s' not found...", options.configFilePath)
			} else {
				return fmt.Errorf("failed to read config file '%s': %w", options.configFilePath, err)
			}
		} else {
			log.Printf("Info: Successfully read config file '%s'.", options.configFilePath)
		}
	} else {
		log.Println("Info: No config file path provided...")
	}

	// --- 新增步骤: 从结构体标签设置 Viper 默认值 ---
	// --- New Step: Set Viper defaults from struct tags ---
	// 必须在 ReadInConfig 和 AutomaticEnv/bindEnvs 之后调用，以避免覆盖用户设置
	// (Must be called AFTER ReadInConfig and AutomaticEnv/bindEnvs to avoid overriding user settings)
	if err := setDefaultsFromTags(v, "", cfg); err != nil {
		// 通常这里的错误是转换错误，前面会有日志打印，可能不需要返回错误
		// (Usually errors here are conversion errors, logged previously, maybe not return error)
		log.Printf("Warning: Error setting defaults from tags: %v", err)
	}
	// --- 结束新增步骤 ---

	// 5. 将 Viper 配置解组到结构体中 (Unmarshal the Viper config into the struct)

	// 直接使用 mapstructure 库，更精确地控制解码过程
	// (Use mapstructure library directly for more precise control over decoding)
	decoderConfig := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		),
		WeaklyTypedInput: true,
		TagName:          "mapstructure",
		Result:           cfg,
		Squash:           true,
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return fmt.Errorf("failed to create mapstructure decoder: %w", err)
	}

	if err := decoder.Decode(v.AllSettings()); err != nil {
		return fmt.Errorf("failed to unmarshal config from mapstructure: %w", err)
	}

	// 7. 配置热重载（如果启用）(Configure hot reload if enabled)
	if options.enableHotReload && options.configFilePath != "" {
		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			log.Printf("Config file changed: %s. Reloading...", e.Name)
			// 热重载时是否需要重新应用默认值？这是一个复杂的问题。
			// (Should defaults be reapplied on hot reload? This is a complex question.)
			// 当前逻辑：不重新应用默认值，只用 Viper 的新值覆盖。
			// (Current logic: Defaults are NOT reapplied; only new Viper values overwrite.)

			// 使用与主流程相同的 mapstructure 解码方法，而不是 v.Unmarshal
			// (Use the same mapstructure decoding method as the main process, rather than v.Unmarshal)
			newDecoderConfig := &mapstructure.DecoderConfig{
				WeaklyTypedInput: true,
				TagName:          "mapstructure",
				Result:           cfg,
				Squash:           true,
				// 添加相同的类型转换钩子
				// (Add the same type conversion hooks)
				DecodeHook: mapstructure.ComposeDecodeHookFunc(
					mapstructure.StringToTimeDurationHookFunc(),
					mapstructure.StringToSliceHookFunc(","),
				),
			}

			newDecoder, errDecoder := mapstructure.NewDecoder(newDecoderConfig)
			if errDecoder != nil {
				log.Printf("Error creating decoder during hot reload: %v", errDecoder)
				return
			}

			if errUnmarshal := newDecoder.Decode(v.AllSettings()); errUnmarshal != nil {
				log.Printf("Error re-unmarshalling config during hot reload: %v", errUnmarshal)
				return
			}

			log.Println("Config reloaded successfully.")
			updateGlobalCfg(cfg) // from accessors.go
		})
		log.Printf("Hot reload enabled for config file: %s", options.configFilePath)
	}

	// 8. 更新全局 Cfg 变量 (Update the global Cfg variable)
	updateGlobalCfg(cfg) // from accessors.go

	return nil
}
