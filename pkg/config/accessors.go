/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package config

import (
	"log"
	"reflect"
	"time"

	"github.com/spf13/viper"
)

// updateGlobalCfg 尝试更新全局 Cfg 变量
// (updateGlobalCfg tries to update the global Cfg variable)
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

// GetString 从加载的配置中获取字符串值 (Gets a string value from the loaded config)
// 注意：这些 Getters 依赖于 LoadConfig 正确更新了全局 Cfg
// (Note: These Getters rely on LoadConfig having correctly updated the global Cfg)
// 警告：访问的是 Viper 的默认实例，可能不准确。推荐直接访问 Cfg 变量。
// (Warning: Accesses Viper's default instance, which might be inaccurate. Recommend accessing the Cfg variable directly.)
func GetString(key string) string {
	log.Println("Warning: GetString(key) accesses the default Viper instance, which might not reflect the loaded Cfg struct accurately after initial load or hot reload. Access fields via the Cfg variable directly.")
	return viper.GetString(key) // Accesses viper's global instance, potentially inaccurate
}

// GetInt 从加载的配置中获取整数值 (Gets an integer value from the loaded config)
// 警告：访问的是 Viper 的默认实例，可能不准确。推荐直接访问 Cfg 变量。
// (Warning: Accesses Viper's default instance, which might be inaccurate. Recommend accessing the Cfg variable directly.)
func GetInt(key string) int {
	log.Println("Warning: GetInt(key) accesses the default Viper instance, which might not reflect the loaded Cfg struct accurately. Access fields via the Cfg variable directly.")
	return viper.GetInt(key)
}

// GetBool 从加载的配置中获取布尔值 (Gets a boolean value from the loaded config)
// 警告：访问的是 Viper 的默认实例，可能不准确。推荐直接访问 Cfg 变量。
// (Warning: Accesses Viper's default instance, which might be inaccurate. Recommend accessing the Cfg variable directly.)
func GetBool(key string) bool {
	log.Println("Warning: GetBool(key) accesses the default Viper instance, which might not reflect the loaded Cfg struct accurately. Access fields via the Cfg variable directly.")
	return viper.GetBool(key)
}

// GetDuration 从加载的配置中获取时间段值 (Gets a time duration value from the loaded config)
// 警告：访问的是 Viper 的默认实例，可能不准确。推荐直接访问 Cfg 变量。
// (Warning: Accesses Viper's default instance, which might be inaccurate. Recommend accessing the Cfg variable directly.)
func GetDuration(key string) time.Duration {
	log.Println("Warning: GetDuration(key) accesses the default Viper instance, which might not reflect the loaded Cfg struct accurately. Access fields via the Cfg variable directly.")
	return viper.GetDuration(key)
}

// GetStringSlice 从加载的配置中获取字符串切片值 (Gets a string slice value from the loaded config)
// 警告：访问的是 Viper 的默认实例，可能不准确。推荐直接访问 Cfg 变量。
// (Warning: Accesses Viper's default instance, which might be inaccurate. Recommend accessing the Cfg variable directly.)
func GetStringSlice(key string) []string {
	log.Println("Warning: GetStringSlice(key) accesses the default Viper instance, which might not reflect the loaded Cfg struct accurately. Access fields via the Cfg variable directly.")
	return viper.GetStringSlice(key)
}

// IsSet 检查某个键是否在 Viper 默认实例中设置 (Checks if a key is set in the Viper default instance)
// 警告：访问的是 Viper 的默认实例，可能不准确。推荐直接访问 Cfg 变量。
// (Warning: Accesses Viper's default instance, which might be inaccurate. Recommend accessing the Cfg variable directly.)
func IsSet(key string) bool {
	log.Println("Warning: IsSet(key) accesses the default Viper instance, which might not reflect the loaded Cfg struct accurately. Access fields via the Cfg variable directly.")
	return viper.IsSet(key)
}

// AllSettings 获取 Viper 默认实例的所有设置 (Gets all settings from the Viper default instance)
// 警告：访问的是 Viper 的默认实例，可能不准确。推荐直接访问 Cfg 变量。
// (Warning: Accesses Viper's default instance, which might be inaccurate. Recommend accessing the Cfg variable directly.)
func AllSettings() map[string]interface{} {
	log.Println("Warning: AllSettings() accesses the default Viper instance, which might not reflect the loaded Cfg struct accurately. Access fields via the Cfg variable directly.")
	return viper.AllSettings()
}
