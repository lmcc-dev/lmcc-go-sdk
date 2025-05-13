/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package config

import (
	"log"
	"reflect"
	// "time" // No longer needed after removing GetDuration
	// "github.com/spf13/viper" // No longer needed after removing viper dependent accessors
)

// updateGlobalCfg 尝试根据传入的配置结构体 `cfg` 更新全局配置变量 `config.Cfg`。
// 它会检查 `cfg` 是否直接是 `*config.Config`，或者是否嵌入了 `config.Config`，
// 或者是否包含一个指向已初始化 `*config.Config` 的指针字段。
// (updateGlobalCfg attempts to update the global configuration variable `config.Cfg` based on the provided config struct `cfg`.)
// (It checks if `cfg` is directly `*config.Config`, or if it embeds `config.Config`,
// or if it contains an initialized pointer field to `*config.Config`.)
// Parameters:
//   cfg: 指向用户加载的配置结构体的指针。
//        (A pointer to the user-loaded configuration struct.)
func updateGlobalCfg[T any](cfg *T) {
	targetVal := reflect.ValueOf(cfg).Elem() // Assuming cfg is always a pointer to the user's config struct

	// 检查 cfg 本身是否就是 *config.Config (Check if cfg itself is *config.Config)
	if baseCfg, ok := interface{}(cfg).(*Config); ok {
		Cfg = baseCfg
		log.Println("Debug: Updated global Cfg directly.")
		return
	}

	// 检查 cfg 是否嵌入了 config.Config (Check if cfg embeds config.Config)
	// 查找名为 "Config" 的字段，其类型为 config.Config (Find field named "Config" of type config.Config)
	field := targetVal.FieldByName("Config")
	if field.IsValid() && field.Type() == reflect.TypeOf(Config{}) {
		// 确保我们可以获取嵌入字段的地址 (Ensure we can get the address of the embedded field)
		if field.CanAddr() {
			embeddedCfgPtr, ok := field.Addr().Interface().(*Config)
			if ok {
				Cfg = embeddedCfgPtr // 设置全局 Cfg 指向嵌入的 Config 实例 (Set global Cfg to point to the embedded Config instance)
				log.Println("Debug: Updated global Cfg from embedded field.")
				return
			}
		} else {
			log.Println("Warning: Found embedded 'Config' field, but cannot take its address to update global Cfg.")
		}
	}

	// 如果 cfg 包含一个指向 config.Config 的指针字段
	// (If cfg contains a pointer field to config.Config)
	for i := 0; i < targetVal.NumField(); i++ {
		fieldVal := targetVal.Field(i)
		fieldType := targetVal.Type().Field(i)
		if fieldVal.Kind() == reflect.Ptr && fieldVal.Type().Elem() == reflect.TypeOf(Config{}) {
			// 检查指针是否已初始化 (Check if the pointer is initialized)
			if !fieldVal.IsNil() {
				if baseCfgPtr, ok := fieldVal.Interface().(*Config); ok {
					Cfg = baseCfgPtr
					log.Printf("Debug: Updated global Cfg from pointer field '%s'.", fieldType.Name)
					return
				}
			} else {
				log.Printf("Warning: Found pointer field '%s' of type *config.Config, but it is nil.", fieldType.Name)
			}
		}
	}

	log.Println("Warning: Global Cfg variable was not updated. Provided type is not *config.Config and does not embed config.Config, nor does it contain an initialized pointer field to *config.Config.")
}

// Note: The following accessor functions (GetString, GetInt, GetBool, GetDuration, GetStringSlice, IsSet, AllSettings)
// have been removed. They relied on Viper's global instance and were not recommended for use.
// Please access configuration values directly through the global `Cfg` variable (obtained via `GetGlobalCfg()`)
// or through the specific configuration struct instance returned by `LoadConfigAndWatch`.
