/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package config

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	lmccerrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors" // SDK errors package (SDK 错误包)
	"github.com/spf13/viper"                                // Viper configuration management
)

// initializeNilPointers 递归地初始化给定结构体指针 target 内部所有为 nil 的结构体指针字段。
// 这确保了在后续处理（如设置默认值或解码）之前，所有嵌套结构体都已分配内存。
// (initializeNilPointers recursively initializes all nil struct pointer fields within the given struct pointer target.)
// (This ensures that all nested structs are allocated before subsequent processing like setting defaults or decoding.)
// Parameters:
//   target: 指向要初始化的结构体的指针。
//           (A pointer to the struct to initialize.)
func initializeNilPointers(target interface{}) {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr {
		return // Must be a pointer
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return // Must point to a struct
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := v.Type().Field(i)

		// 只初始化结构体指针 (Only initialize struct pointers)
		if field.Kind() == reflect.Ptr && field.IsNil() && field.Type().Elem().Kind() == reflect.Struct {
			if field.CanSet() {
				newStructPtr := reflect.New(fieldType.Type.Elem())
				field.Set(newStructPtr)
				initializeNilPointers(newStructPtr.Interface()) // 递归初始化新创建的结构体 (Recursively initialize the newly created struct)
			}
		} else if field.Kind() == reflect.Struct && field.CanAddr() {
			// 递归处理嵌套的值类型结构体 (如果它们是指针的目标)
			// (Recursively handle nested value structs if they are addressable)
			initializeNilPointers(field.Addr().Interface())
		} else if field.Kind() == reflect.Ptr && !field.IsNil() && field.Elem().Kind() == reflect.Struct {
			// 如果指针非 nil，递归处理它指向的结构体
			// (If pointer is not nil, recurse into the struct it points to)
			initializeNilPointers(field.Interface())
		}
	}
}

// Removed unused functions: applyDefaultsFromStructTags, parseDefaultValueToType
// (移除了未使用的函数：applyDefaultsFromStructTags, parseDefaultValueToType)

// setDefaultsFromTags 递归地遍历配置结构体 `config`，读取字段上的 `default` 标签，
// 并使用 `v.SetDefault` 将其设置到 Viper 实例 `v` 中。
// `keyPrefix` 用于在递归时构建嵌套的 Viper 键。
// (setDefaultsFromTags recursively traverses the configuration struct `config`, reads the `default` tag on fields,
// and sets them into the Viper instance `v` using `v.SetDefault`.)
// (`keyPrefix` is used to build nested Viper keys during recursion.)
// Parameters:
//   v: 要设置默认值的 Viper 实例。
//      (The Viper instance to set defaults on.)
//   config: 包含 `default` 标签的配置结构体实例（或指向它的指针）。
//           (The configuration struct instance (or pointer to it) containing `default` tags.)
//   keyPrefix: 当前递归层级的 Viper 键前缀。
//              (The Viper key prefix for the current recursion level.)
// Returns:
//   error: 解析或设置默认值过程中发生的任何错误。
//          (Any error that occurs during parsing or setting defaults.)
func setDefaultsFromTags(v *viper.Viper, config interface{}, keyPrefix string) error {
	val := reflect.ValueOf(config)
	typ := reflect.TypeOf(config)

	// Handle pointers correctly (正确处理指针)
	if typ.Kind() == reflect.Ptr {
		if val.IsNil() {
			log.Printf("Warning: Encountered nil pointer at prefix '%s', cannot set defaults for it.", keyPrefix)
			return nil // Cannot proceed with nil pointer, but not an error for the caller of setDefaultsFromTags itself.
		}
		val = val.Elem()
		typ = val.Type()
	}

	if typ.Kind() != reflect.Struct {
		log.Printf("Warning: Expected a struct or pointer to struct at prefix '%s', got %s. Skipping defaults.", keyPrefix, typ.Kind())
		return nil // Not an error for the caller, just skipping this part of the config.
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		// Skip unexported fields (跳过未导出字段)
		if !field.IsExported() {
			continue
		}

		// Build the viper key using mapstructure tag first, then field name (构建 viper key)
		mapKey := field.Tag.Get("mapstructure")
		if mapKey == "" {
			mapKey = field.Name // Use field name if mapstructure tag is absent (如果 mapstructure 标签不存在，则使用字段名)
		}

		// Handle the case where mapKey might be "-" to skip the field
		if mapKey == "-" {
			continue
		}

		// Handle omitted fields (`mapstructure:",omitempty"`)
		// 直接使用 TrimSuffix，如果后缀不存在，则不会有任何影响
		// (Directly use TrimSuffix, it has no effect if the suffix doesn't exist)
		mapKey = strings.TrimSuffix(mapKey, ",omitempty")

		// Construct the full key path (构建完整的 key 路径)
		fullKey := mapKey
		if keyPrefix != "" {
			fullKey = keyPrefix + "." + mapKey
		}

		// Get the default value tag (获取默认值标签)
		defaultValue := field.Tag.Get("default")

		// Handle nested structs (recursively) *before* setting default for the current level
		// This ensures viper keys for nested fields exist before setting a default on the parent struct itself (if applicable)
		// (在设置当前级别的默认值*之前*递归处理嵌套结构体)
		if field.Type.Kind() == reflect.Struct {
			// Recurse into struct value type (递归值类型结构体)
			if err := setDefaultsFromTags(v, fieldVal.Addr().Interface(), fullKey); err != nil {
				return err // Propagate error
			}
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
			// Recurse into struct pointer type (递归指针类型结构体)
			if !fieldVal.IsNil() {
				if err := setDefaultsFromTags(v, fieldVal.Interface(), fullKey); err != nil {
					return err // Propagate error
				}
			} // Do not create a tempInstance and recurse if fieldVal is nil, as this would set defaults for a non-existent path in Viper.
			  // (如果 fieldVal 为 nil，则不创建 tempInstance 并递归，因为这会在 Viper 中为不存在的路径设置默认值。)
		}

		// Set the default value in Viper if tag exists and key is not already set
		// (如果标签存在且键尚未设置，则在 Viper 中设置默认值)
		if defaultValue != "" {
			// Viper's SetDefault only sets if the key is not already defined
			// We attempt to parse the default value to the correct type for Viper
			parsedVal, err := parseStringToType(defaultValue, field.Type)
			if err != nil {
				// Instead of logging and continuing, return the error
				// (不再是记录日志并继续，而是返回错误)
				return lmccerrors.WithCode(
					lmccerrors.Wrapf(err, "failed to parse default tag value '%s' for key '%s' (field %s)", defaultValue, fullKey, field.Name),
					lmccerrors.ErrConfigDefaultTagParse,
				)
			}
			v.SetDefault(fullKey, parsedVal)
		}
	}
	return nil // Return nil if no errors occurred (如果没有发生错误则返回 nil)
}

// parseStringToType 将字符串值 `value` 解析为 `targetType` 指定的 Go 类型。
// 支持基本类型 (string, int*, uint*, float*, bool), time.Duration, 以及 string 切片 (逗号或空格分隔)。
// (parseStringToType parses the string `value` into the Go type specified by `targetType`.)
// (Supports basic types (string, int*, uint*, float*, bool), time.Duration, and string slices (comma or space separated).)
// Parameters:
//   value: 要解析的字符串值。
//          (The string value to parse.)
//   targetType: 目标 Go 类型的反射类型。
//               (The reflection type of the target Go type.)
// Returns:
//   interface{}: 解析后的值。
//                (The parsed value.)
//   error: 解析过程中发生的错误，或类型不受支持。
//          (Any error during parsing, or if the type is unsupported.)
func parseStringToType(value string, targetType reflect.Type) (interface{}, error) {
	// Check if targetType itself is nil (e.g. reflect.TypeOf(nil) returns nil)
	// (检查 targetType 本身是否为 nil（例如 reflect.TypeOf(nil) 返回 nil）)
	if targetType == nil {
		return nil, lmccerrors.NewWithCode(lmccerrors.ErrConfigInternal, "targetType cannot be nil for parsing") // Cannot parse to a nil type (无法解析为 nil 类型)
	}
	kind := targetType.Kind()

	// Handle pointer types by looking at the element type
	// (通过查看元素类型来处理指针类型)
	if kind == reflect.Ptr {
		// Get the element type. If targetType represents a nil pointer type itself (which is unusual for a Type),
		// or if Elem() somehow results in a nil Type (also unusual), it's an internal error.
		// (获取元素类型。如果 targetType 本身代表一个 nil 指针类型（这对于 Type 来说不寻常），
		// 或者如果 Elem() 不知何故导致了 nil Type（也不寻常），则这是一个内部错误。)
		elemType := targetType.Elem()
		if elemType == nil { // Check if the element type is nil
			return nil, lmccerrors.NewWithCode(lmccerrors.ErrConfigInternal, "targetType is a pointer to an undefined element type")
		}
		targetType = elemType
		kind = targetType.Kind()
	}

	switch kind {
	case reflect.String:
		return value, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if targetType == reflect.TypeOf(time.Duration(0)) {
					d, err := time.ParseDuration(value)
		if err != nil {
			return nil, lmccerrors.WithCode(
				lmccerrors.Wrapf(err, "failed to parse duration string '%s'", value),
				lmccerrors.ErrConfigDefaultTagParse,
			)
		}
			return d, nil
		}
		parsedInt, err := strconv.ParseInt(value, 0, targetType.Bits())
		if err != nil {
			return nil, lmccerrors.WithCode(
				lmccerrors.Wrapf(err, "invalid integer format for '%s'", value),
				lmccerrors.ErrConfigDefaultTagParse,
			)
		}
		// Convert the int64 to the specific target integer type using reflection
		// (使用反射将 int64 转换为特定的目标整数类型)
		resultVal := reflect.New(targetType).Elem()
		resultVal.SetInt(parsedInt)
		return resultVal.Interface(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		parsedUint, err := strconv.ParseUint(value, 0, targetType.Bits())
		if err != nil {
			return nil, lmccerrors.WithCode(
				lmccerrors.Wrapf(err, "invalid unsigned integer format for '%s'", value),
				lmccerrors.ErrConfigDefaultTagParse,
			)
		}
		resultVal := reflect.New(targetType).Elem()
		resultVal.SetUint(parsedUint)
		return resultVal.Interface(), nil
	case reflect.Float32, reflect.Float64:
		parsedFloat, err := strconv.ParseFloat(value, targetType.Bits())
		if err != nil {
			return nil, lmccerrors.WithCode(
				lmccerrors.Wrapf(err, "invalid float format for '%s'", value),
				lmccerrors.ErrConfigDefaultTagParse,
			)
		}
		resultVal := reflect.New(targetType).Elem()
		resultVal.SetFloat(parsedFloat)
		return resultVal.Interface(), nil
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return nil, lmccerrors.WithCode(
				lmccerrors.Wrapf(err, "invalid boolean format for '%s'", value),
				lmccerrors.ErrConfigDefaultTagParse,
			)
		}
		return b, nil
	case reflect.Slice:
		// Handle string slices specifically (特别处理字符串切片)
		if targetType.Elem().Kind() == reflect.String {
			if value == "" {
				return []string{}, nil // Empty string means empty slice (空字符串表示空切片)
			}
			// Allow comma or space separation, trim whitespace (允许逗号或空格分隔，修剪空白)
			splitFunc := func(c rune) bool {
				return c == ',' || c == ' '
			}
			parts := strings.FieldsFunc(value, splitFunc)
			// Trim whitespace from each part (修剪每个部分的空白)
			for i, p := range parts {
				parts[i] = strings.TrimSpace(p)
			}
			return parts, nil
		}
		return nil, lmccerrors.NewWithCode(lmccerrors.ErrConfigDefaultTagParse, fmt.Sprintf("unsupported slice type for default tag: %s. Only []string is supported.", targetType.Elem().Kind()))
	default:
		return nil, lmccerrors.NewWithCode(lmccerrors.ErrConfigDefaultTagParse, fmt.Sprintf("unsupported type for default tag: %s", kind))
	}
}

// applyDefaultsToZeroFieldsWithViper 递归地将结构体中带有 `default` 标签的零值字段设置为其默认值，
// 但只有当该字段在配置文件中实际不存在时才设置。这样可以避免覆盖从配置文件显式设置的值。
// (applyDefaultsToZeroFieldsWithViper recursively sets zero-value fields with `default` tags in a struct to their default values,
// but only when the field doesn't actually exist in the config file. This avoids overriding values explicitly set from config files.)
// Parameters:
//   target: 指向要应用默认值的结构体的指针。
//           (A pointer to the struct to apply defaults to.)
//   v: Viper 实例，用于检查键是否存在。
//      (Viper instance to check if keys exist.)
//   keysFromConfigFile: 配置文件中实际存在的键的映射。
//                       (Map of keys actually present in the config file.)
// Returns:
//   error: 如果解析默认标签或设置值时发生错误。
//          (Error if parsing default tag or setting value fails.)
func applyDefaultsToZeroFieldsWithViper(target interface{}, v *viper.Viper, keysFromConfigFile map[string]bool) error {
	return applyDefaultsToZeroFieldsWithViperInternal(target, v, keysFromConfigFile, "")
}

// applyDefaultsToZeroFieldsWithViperInternal 是内部递归函数
// (applyDefaultsToZeroFieldsWithViperInternal is the internal recursive function)
func applyDefaultsToZeroFieldsWithViperInternal(target interface{}, v *viper.Viper, keysFromConfigFile map[string]bool, keyPrefix string) error {
	val := reflect.ValueOf(target)

	if val.Kind() != reflect.Ptr || val.IsNil() || val.Elem().Kind() != reflect.Struct {
		return lmccerrors.NewWithCode(lmccerrors.ErrConfigInternal, "applyDefaultsToZeroFieldsWithViper expects a non-nil pointer to a struct")
	}

	structVal := val.Elem()
	structType := structVal.Type()

	for i := 0; i < structVal.NumField(); i++ {
		fieldVal := structVal.Field(i)
		fieldType := structType.Field(i)

		if !fieldVal.CanSet() || !fieldType.IsExported() {
			continue
		}

		// 获取 mapstructure 标签来构建 Viper 键
		// (Get mapstructure tag to build Viper key)
		mapstructureTag := fieldType.Tag.Get("mapstructure")
		
		// 处理嵌入字段：如果没有 mapstructure 标签且字段是匿名的，则递归处理
		// (Handle embedded fields: if no mapstructure tag and field is anonymous, recurse)
		if mapstructureTag == "" && fieldType.Anonymous {
			// 对于嵌入字段，使用当前的 keyPrefix 而不是构建新的键
			// (For embedded fields, use current keyPrefix instead of building new keys)
			if fieldVal.Kind() == reflect.Struct && fieldVal.CanAddr() {
				if err := applyDefaultsToZeroFieldsWithViperInternal(fieldVal.Addr().Interface(), v, keysFromConfigFile, keyPrefix); err != nil {
					return err
				}
			} else if fieldVal.Kind() == reflect.Ptr && !fieldVal.IsNil() && fieldVal.Type().Elem().Kind() == reflect.Struct {
				if err := applyDefaultsToZeroFieldsWithViperInternal(fieldVal.Interface(), v, keysFromConfigFile, keyPrefix); err != nil {
					return err
				}
			}
			continue
		}
		
		if mapstructureTag == "" || mapstructureTag == "-" {
			continue
		}

		// 构建完整的 Viper 键路径
		// (Build complete Viper key path)
		var fullKey string
		if keyPrefix == "" {
			fullKey = mapstructureTag
		} else {
			fullKey = keyPrefix + "." + mapstructureTag
		}
		
		// Viper 将所有键转换为小写，所以我们也需要将 fullKey 转换为小写进行比较
		// (Viper converts all keys to lowercase, so we need to convert fullKey to lowercase for comparison)
		fullKeyLower := strings.ToLower(fullKey)

		defaultTag := fieldType.Tag.Get("default")
		isZero := fieldVal.IsZero()

		// Handle pointers to structs: if nil, initialize and recurse
		// (处理结构体指针：如果为 nil，则初始化并递归)
		if fieldVal.Kind() == reflect.Ptr && fieldVal.Type().Elem().Kind() == reflect.Struct {
			if fieldVal.IsNil() {
				// Only initialize if there's a default tag anywhere within the nested struct or if the pointer field itself has a default.
				// (仅当嵌套结构体内部任何位置有默认标签，或者指针字段本身有默认值时才初始化。)
				hasAnyDefault := defaultTag != "" || hasDefaultsDefined(fieldVal.Type().Elem())
				if hasAnyDefault {
					newStructPtr := reflect.New(fieldVal.Type().Elem())
					fieldVal.Set(newStructPtr)
					if err := applyDefaultsToZeroFieldsWithViperInternal(fieldVal.Interface(), v, keysFromConfigFile, fullKey); err != nil {
						return err // Propagate error from recursive call (从递归调用传播错误)
					}
					isZero = false // No longer considered zero for the purpose of applying a default *to the pointer itself*
				}
			} else {
				// If not nil, recurse to apply defaults to its fields
				// (如果不为 nil，则递归以将默认值应用于其字段)
				if err := applyDefaultsToZeroFieldsWithViperInternal(fieldVal.Interface(), v, keysFromConfigFile, fullKey); err != nil {
					return err // Propagate error
				}
			}
		} else if fieldVal.Kind() == reflect.Struct && fieldVal.CanAddr() {
			// Recurse for non-pointer struct fields to handle their nested defaults
			// (对非指针结构体字段进行递归以处理其嵌套的默认值)
			if err := applyDefaultsToZeroFieldsWithViperInternal(fieldVal.Addr().Interface(), v, keysFromConfigFile, fullKey); err != nil {
				return err // Propagate error
			}
		}

		// Apply default value if the field is zero, has a default tag, AND the key was not present in config file
		// (如果字段为零、存在默认标签且该键在配置文件中不存在，则应用默认值)
		if isZero && defaultTag != "" && !keysFromConfigFile[fullKeyLower] {
			parsedVal, err := parseStringToType(defaultTag, fieldVal.Type())
			if err != nil {
				return lmccerrors.WithCode(
					lmccerrors.Wrapf(err, "error parsing default tag '%s' for field '%s.%s'", defaultTag, structType.Name(), fieldType.Name),
					lmccerrors.ErrConfigDefaultTagParse,
				)
			}

			targetVal := reflect.ValueOf(parsedVal)
			if fieldVal.Kind() == reflect.Ptr {
				// If the field is a pointer, create a new pointer to the parsed value
				// (如果字段是指针，则创建指向解析值的新指针)
				ptr := reflect.New(fieldVal.Type().Elem())
				// Ensure targetVal is assignable to ptr.Elem()
				if ptr.Elem().Type() != targetVal.Type() {
					// This can happen if parseStringToType returns, e.g., int, but field is *int64.
					// We need to convert targetVal to the element type of the pointer.
					if !targetVal.CanConvert(ptr.Elem().Type()) {
						return lmccerrors.NewWithCode(lmccerrors.ErrConfigInternal, 
							fmt.Sprintf("type mismatch: cannot convert parsed default value of type %s to field %s.%s's element type %s", 
							targetVal.Type(), structType.Name(), fieldType.Name, ptr.Elem().Type()))
					}
					targetVal = targetVal.Convert(ptr.Elem().Type())
				}
				ptr.Elem().Set(targetVal)
				fieldVal.Set(ptr)
			} else {
				// Ensure targetVal is assignable to fieldVal
				if fieldVal.Type() != targetVal.Type() {
					if !targetVal.CanConvert(fieldVal.Type()) {
						return lmccerrors.NewWithCode(lmccerrors.ErrConfigInternal, 
							fmt.Sprintf("type mismatch: cannot convert parsed default value of type %s to field %s.%s's type %s", 
							targetVal.Type(), structType.Name(), fieldType.Name, fieldVal.Type()))
					}
					targetVal = targetVal.Convert(fieldVal.Type())
				}
				fieldVal.Set(targetVal)
			}
		}
	}
	return nil
}

// applyDefaultsToZeroFields 递归地将结构体中带有 `default` 标签的零值字段设置为其默认值。
// 注意：此函数直接修改传入的结构体。
// (applyDefaultsToZeroFields recursively sets zero-value fields with `default` tags in a struct to their default values.)
// (Note: This function modifies the passed-in struct directly.)
// Parameters:
//   target: 指向要应用默认值的结构体的指针。
//           (A pointer to the struct to apply defaults to.)
// Returns:
//   error: 如果解析默认标签或设置值时发生错误。
//          (Error if parsing default tag or setting value fails.)
func applyDefaultsToZeroFields(target interface{}) error {
	v := reflect.ValueOf(target)

	if v.Kind() != reflect.Ptr || v.IsNil() || v.Elem().Kind() != reflect.Struct {
		return lmccerrors.NewWithCode(lmccerrors.ErrConfigInternal, "applyDefaultsToZeroFields expects a non-nil pointer to a struct")
	}

	structVal := v.Elem()
	structType := structVal.Type()

	for i := 0; i < structVal.NumField(); i++ {
		fieldVal := structVal.Field(i)
		fieldType := structType.Field(i)

		if !fieldVal.CanSet() || !fieldType.IsExported() {
			continue
		}

		defaultTag := fieldType.Tag.Get("default")
		isZero := fieldVal.IsZero()

		// Handle pointers to structs: if nil, initialize and recurse
		// (处理结构体指针：如果为 nil，则初始化并递归)
		if fieldVal.Kind() == reflect.Ptr && fieldVal.Type().Elem().Kind() == reflect.Struct {
			if fieldVal.IsNil() {
				// Only initialize if there's a default tag anywhere within the nested struct or if the pointer field itself has a default.
				// (仅当嵌套结构体内部任何位置有默认标签，或者指针字段本身有默认值时才初始化。)
				// For simplicity in this function, we check if the pointer field itself has a default OR if any sub-field has one.
				hasAnyDefault := defaultTag != "" || hasDefaultsDefined(fieldVal.Type().Elem())
				if hasAnyDefault {
					newStructPtr := reflect.New(fieldVal.Type().Elem())
					fieldVal.Set(newStructPtr)
					if err := applyDefaultsToZeroFields(fieldVal.Interface()); err != nil {
						return err // Propagate error from recursive call (从递归调用传播错误)
					}
					isZero = false // No longer considered zero for the purpose of applying a default *to the pointer itself*
				}
			} else {
				// If not nil, recurse to apply defaults to its fields
				// (如果不为 nil，则递归以将默认值应用于其字段)
				if err := applyDefaultsToZeroFields(fieldVal.Interface()); err != nil {
					return err // Propagate error
				}
			}
		} else if fieldVal.Kind() == reflect.Struct && fieldVal.CanAddr() {
			// Recurse for non-pointer struct fields to handle their nested defaults
			// (对非指针结构体字段进行递归以处理其嵌套的默认值)
			if err := applyDefaultsToZeroFields(fieldVal.Addr().Interface()); err != nil {
				return err // Propagate error
			}
		}

		// Apply default value if the field is zero and a default tag exists
		// (如果字段为零且存在默认标签，则应用默认值)
		if isZero && defaultTag != "" {
			parsedVal, err := parseStringToType(defaultTag, fieldVal.Type())
			if err != nil {
				return lmccerrors.WithCode(
					lmccerrors.Wrapf(err, "error parsing default tag '%s' for field '%s.%s'", defaultTag, structType.Name(), fieldType.Name),
					lmccerrors.ErrConfigDefaultTagParse,
				)
			}

			targetVal := reflect.ValueOf(parsedVal)
			if fieldVal.Kind() == reflect.Ptr {
				// If the field is a pointer, create a new pointer to the parsed value
				// (如果字段是指针，则创建指向解析值的新指针)
				ptr := reflect.New(fieldVal.Type().Elem())
				// Ensure targetVal is assignable to ptr.Elem()
				if ptr.Elem().Type() != targetVal.Type() {
					// This can happen if parseStringToType returns, e.g., int, but field is *int64.
					// We need to convert targetVal to the element type of the pointer.
					if !targetVal.CanConvert(ptr.Elem().Type()) {
						return lmccerrors.NewWithCode(lmccerrors.ErrConfigInternal, 
							fmt.Sprintf("type mismatch: cannot convert parsed default value of type %s to field %s.%s's element type %s", 
							targetVal.Type(), structType.Name(), fieldType.Name, ptr.Elem().Type()))
					}
					targetVal = targetVal.Convert(ptr.Elem().Type())
				}
				ptr.Elem().Set(targetVal)
				fieldVal.Set(ptr)
					} else {
				// Ensure targetVal is assignable to fieldVal
				if fieldVal.Type() != targetVal.Type() {
					if !targetVal.CanConvert(fieldVal.Type()) {
						return lmccerrors.NewWithCode(lmccerrors.ErrConfigInternal, 
							fmt.Sprintf("type mismatch: cannot convert parsed default value of type %s to field %s.%s's type %s", 
							targetVal.Type(), structType.Name(), fieldType.Name, fieldVal.Type()))
					}
					targetVal = targetVal.Convert(fieldVal.Type())
				}
				fieldVal.Set(targetVal)
			}
		}
	}
	return nil
}

// flattenViperKeys 递归地将 Viper 的 AllSettings() 返回的嵌套映射扁平化为点分隔的键
// (flattenViperKeys recursively flattens the nested map returned by Viper's AllSettings() into dot-separated keys)
func flattenViperKeys(settings map[string]interface{}) map[string]bool {
	result := make(map[string]bool)
	flattenViperKeysRecursive(settings, "", result)
	return result
}

// flattenViperKeysRecursive 是递归辅助函数
// (flattenViperKeysRecursive is the recursive helper function)
func flattenViperKeysRecursive(m map[string]interface{}, prefix string, result map[string]bool) {
	for key, value := range m {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}
		
		// 记录当前键 (Record current key)
		result[fullKey] = true
		
		// 如果值是嵌套映射，递归处理 (If value is nested map, recurse)
		if nestedMap, ok := value.(map[string]interface{}); ok {
			flattenViperKeysRecursive(nestedMap, fullKey, result)
		}
	}
}

// hasDefaultsDefined checks if a struct type or any of its nested struct types have default tags.
// (hasDefaultsDefined 检查结构体类型或其任何嵌套结构体类型是否具有默认标签。)
func hasDefaultsDefined(structType reflect.Type) bool {
	if structType.Kind() != reflect.Struct {
		return false
	}
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if field.Tag.Get("default") != "" {
			return true
		}
		if field.Type.Kind() == reflect.Struct {
			if hasDefaultsDefined(field.Type) {
				return true
			}
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
			if hasDefaultsDefined(field.Type.Elem()) {
				return true
			}
		}
	}
	return false
}

// DefaultConfigLoader is a loader that applies defaults from struct tags to Viper.
// (DefaultConfigLoader 是一个加载器，它将结构体标签中的默认值应用于 Viper。)
type DefaultConfigLoader struct {
	defaults interface{} // A pointer to the struct defining the defaults (指向定义默认值的结构体的指针)
}

// NewDefaultConfigLoader creates a new loader for applying struct tag defaults.
// `defaultsStructPtr` must be a pointer to a struct.
// (NewDefaultConfigLoader 创建一个新的加载器，用于应用结构体标签默认值。)
// (`defaultsStructPtr` 必须是指向结构体的指针。)
func NewDefaultConfigLoader(defaultsStructPtr interface{}) *DefaultConfigLoader {
	val := reflect.ValueOf(defaultsStructPtr)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		panic(lmccerrors.NewWithCode(lmccerrors.ErrConfigInternal, "NewDefaultConfigLoader expects a non-nil pointer to a struct"))
	}
	if val.Elem().Kind() != reflect.Struct {
		panic(lmccerrors.NewWithCode(lmccerrors.ErrConfigInternal, fmt.Sprintf("NewDefaultConfigLoader expects a pointer to a struct, not a pointer to %s", val.Elem().Kind())))
	}
	return &DefaultConfigLoader{defaults: defaultsStructPtr}
}

// Load applies the defaults defined in the struct (via `default` tags) to the Viper instance.
// (Load 将结构体中定义的默认值（通过 `default` 标签）应用于 Viper 实例。)
func (l *DefaultConfigLoader) Load(v *viper.Viper) error {
	if l.defaults == nil {
		return lmccerrors.NewWithCode(lmccerrors.ErrConfigInternal, "DefaultConfigLoader.defaults is nil, cannot load")
	}
	// Set defaults in Viper from struct tags
	// (从结构体标签在 Viper 中设置默认值)
	return setDefaultsFromTags(v, l.defaults, "")
}

// Name returns the name of the loader.
// (Name 返回加载器的名称。)
func (l *DefaultConfigLoader) Name() string {
	return "DefaultTagLoader"
}
