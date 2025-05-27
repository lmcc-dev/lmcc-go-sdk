/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package config

import (
	"log"
	"reflect"
	"strings"

	lmccerrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors" // SDK errors package (SDK 错误包)
	"github.com/spf13/viper"
)

// bindEnvs 递归地遍历配置结构体 `iface`，并使用 Viper 实例 `v` 的 `BindEnv` 方法
// 将每个字段绑定到相应的环境变量。
// 它使用 `replacer` 来转换 Viper 键为环境变量名，并处理 `mapstructure` 或 `json` 标签。
// `parts` 用于在递归时构建嵌套的 Viper 键。
// (bindEnvs recursively traverses the configuration struct `iface` and uses the Viper instance `v`'s `BindEnv` method
// to bind each field to its corresponding environment variable.)
// (It uses the `replacer` to convert Viper keys to environment variable names and handles `mapstructure` or `json` tags.)
// (`parts` is used to build nested Viper keys during recursion.)
// Parameters:
//   v: 要绑定环境变量的 Viper 实例。
//      (The Viper instance to bind environment variables on.)
//   replacer: 用于将 Viper 键转换为环境变量名的字符串替换器。
//             (The string replacer used to convert Viper keys to environment variable names.)
//   iface: 当前要处理的配置结构体（或其指针）。
//          (The current configuration struct (or pointer to it) to process.)
//   parts: 构建当前 Viper 键路径的组件。
//          (Components for building the current Viper key path.)
func bindEnvs(v *viper.Viper, replacer *strings.Replacer, iface interface{}, parts ...string) {
	val := reflect.ValueOf(iface)
	typ := reflect.TypeOf(iface)

	// 解引用指针 (Dereference pointer)
	if typ.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	// 必须是结构体才能继续 (Must be a struct to continue)
	if typ.Kind() != reflect.Struct {
		return
	}

	// 获取环境变量前缀 (Get environment variable prefix)
	envPrefix := v.GetEnvPrefix()
	if envPrefix != "" {
		envPrefix += "_"
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		// 跳过未导出的字段 (Skip unexported fields)
		if !fieldVal.CanInterface() {
			continue
		}

		// 获取 mapstructure 标签名，否则使用小写的字段名
		// (Get mapstructure tag name, otherwise use lowercase field name)
		tag := field.Tag.Get("mapstructure")
		if tag == "" {
			// 如果 mapstructure 标签为空，检查是否有 json 标签
			// (If mapstructure tag is empty, check for json tag)
			jsonTag := field.Tag.Get("json")
			if jsonTag != "" && jsonTag != "-" {
				// 取 json 标签逗号前的内容
				// (Take the content before the comma in the json tag)
				tag = strings.Split(jsonTag, ",")[0]
			} else {
				// 都没有则使用小写字段名
				// (If neither exists, use lowercase field name)
				tag = strings.ToLower(field.Name)
			}
		}
		// 忽略 mapstructure 或 json 标签中的 "-"
		// (Ignore "-" in mapstructure or json tags)
		if tag == "-" {
			continue
		}

		// 构建当前键路径 (Build the current key path)
		currentParts := append(parts, tag)
		viperKey := strings.Join(currentParts, ".")

		// 构建环境变量名 (Build the environment variable name)
		// 使用 Viper 的 replacer 来处理 key
		// (Use Viper's replacer to handle the key)
		envVarName := envPrefix + strings.ToUpper(replacer.Replace(viperKey))

		// 递归处理结构体或指向结构体的指针
		// (Recursively handle structs or pointers to structs)
		kind := field.Type.Kind()
		if kind == reflect.Ptr {
			kind = field.Type.Elem().Kind() // Update kind directly from element type
		}

		if kind == reflect.Struct {
			// 如果字段是嵌入的结构体(Anonymous)，我们应该使用当前的 parts 传递给递归调用，
			// 否则（非嵌入结构体字段），我们使用追加了 tag 的 currentParts。
			// (If the field is an embedded struct (Anonymous), we should pass the current parts to the recursive call,
			// otherwise (non-embedded struct field), we use currentParts with the tag appended.)
			recursiveParts := currentParts // 默认用于非嵌入字段 (Default for non-embedded fields)
			if field.Anonymous {
				recursiveParts = parts // 对嵌入字段使用父级 parts (Use parent parts for embedded fields)
			}

			// 如果是指针且为 nil，则不处理 (如果需要处理 nil，需要实例化)
			// (If it's a pointer and nil, don't process (instantiation needed if nil handling required))
			if field.Type.Kind() == reflect.Ptr && fieldVal.IsNil() {
				continue
			}
			// 使用修正后的 recursiveParts 进行递归调用
			// (Use the corrected recursiveParts for the recursive call)
			bindEnvs(v, replacer, fieldVal.Interface(), recursiveParts...)
		} else {
			// 绑定非结构体字段的环境变量
			// (Bind environment variable for non-struct fields)
			if err := v.BindEnv(viperKey, envVarName); err != nil {
				// 通常 BindEnv 不会返回错误，但以防万一
				// (BindEnv usually doesn't return an error, but just in case)
				wrappedErr := lmccerrors.WithCode(
			lmccerrors.Wrapf(err, "failed to bind env var '%s' to key '%s'", envVarName, viperKey),
			lmccerrors.ErrConfigEnvBind,
		)
				log.Printf("Warning: %s: %+v", lmccerrors.ErrConfigEnvBind.String(), wrappedErr)                                                    // 使用标准 log，但错误已包装 (Use standard log, but error is wrapped)
			}
		}
	}
}
