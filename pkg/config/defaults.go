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

	"github.com/spf13/viper" // Add viper import
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
func setDefaultsFromTags(v *viper.Viper, config interface{}, keyPrefix string) {
	val := reflect.ValueOf(config)
	typ := reflect.TypeOf(config)

	// Handle pointers correctly (正确处理指针)
	if typ.Kind() == reflect.Ptr {
		if val.IsNil() {
			log.Printf("Warning: Encountered nil pointer at prefix '%s', cannot set defaults for it.", keyPrefix)
			return // Cannot proceed with nil pointer
		}
		val = val.Elem()
		typ = val.Type()
	}

	if typ.Kind() != reflect.Struct {
		log.Printf("Warning: Expected a struct or pointer to struct at prefix '%s', got %s. Skipping defaults.", keyPrefix, typ.Kind())
		return // Only process structs
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
			setDefaultsFromTags(v, fieldVal.Addr().Interface(), fullKey)
		} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct {
			// Recurse into struct pointer type (递归指针类型结构体)
			if !fieldVal.IsNil() {
				setDefaultsFromTags(v, fieldVal.Interface(), fullKey)
			} else {
				// If the pointer is nil, we might still want to process defaults defined *within* that struct type.
				// Create a temporary instance to traverse its tags.
				tempInstance := reflect.New(field.Type.Elem()).Interface()
				setDefaultsFromTags(v, tempInstance, fullKey)
			}
		}

		// Set the default value in Viper if tag exists and key is not already set
		// (如果标签存在且键尚未设置，则在 Viper 中设置默认值)
		if defaultValue != "" {
			// Viper's SetDefault only sets if the key is not already defined
			// We attempt to parse the default value to the correct type for Viper
			parsedVal, err := parseStringToType(defaultValue, field.Type)
			if err != nil {
				log.Printf("Warning: Failed to parse default tag value '%s' for key '%s' (field %s): %v. Setting as string.", defaultValue, fullKey, field.Name, err)
				v.SetDefault(fullKey, defaultValue) // Fallback to setting as string
			} else {
				v.SetDefault(fullKey, parsedVal)
			}
		}
	}
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
	kind := targetType.Kind()

	// Handle pointer types by looking at the element type
	// (通过查看元素类型来处理指针类型)
	if kind == reflect.Ptr {
		targetType = targetType.Elem()
		kind = targetType.Kind()
	}

	switch kind {
	case reflect.String:
		return value, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if targetType == reflect.TypeOf(time.Duration(0)) {
			return time.ParseDuration(value)
		}
		parsedInt, err := strconv.ParseInt(value, 0, targetType.Bits())
		if err != nil {
			return nil, fmt.Errorf("invalid integer format for '%s': %w", value, err)
		}
		// Convert the int64 to the specific target integer type using reflection
		resultVal := reflect.New(targetType).Elem()
		resultVal.SetInt(parsedInt)
		return resultVal.Interface(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		parsedUint, err := strconv.ParseUint(value, 0, targetType.Bits())
		if err != nil {
			return nil, fmt.Errorf("invalid unsigned integer format for '%s': %w", value, err)
		}
		resultVal := reflect.New(targetType).Elem()
		resultVal.SetUint(parsedUint)
		return resultVal.Interface(), nil
	case reflect.Float32, reflect.Float64:
		parsedFloat, err := strconv.ParseFloat(value, targetType.Bits())
		if err != nil {
			return nil, fmt.Errorf("invalid float format for '%s': %w", value, err)
		}
		resultVal := reflect.New(targetType).Elem()
		resultVal.SetFloat(parsedFloat)
		return resultVal.Interface(), nil
	case reflect.Bool:
		return strconv.ParseBool(value)
	case reflect.Slice:
		// Handle string slices specifically (特别处理字符串切片)
		if targetType.Elem().Kind() == reflect.String {
			if value == "" {
				return []string{}, nil // Empty string means empty slice
			}
			// Allow comma or space separation, trim whitespace (允许逗号或空格分隔，修剪空白)
			splitFunc := func(c rune) bool {
				return c == ',' || c == ' '
			}
			parts := strings.FieldsFunc(value, splitFunc)
			// Trim whitespace from each part
			for i, p := range parts {
				parts[i] = strings.TrimSpace(p)
			}
			return parts, nil
		}
		return nil, fmt.Errorf("parsing default values for non-string slices (type %v) is not supported", targetType)
	default:
		return nil, fmt.Errorf("unsupported type '%v' for default value parsing", targetType)
	}
}

// applyDefaultsToZeroFields 递归地遍历目标结构体 `target`，
// 对于值为零值 (或 nil 指针) 且具有 `default` 标签的字段，尝试解析标签值并设置该字段。
// 这个函数设计在 Viper/mapstructure 解码 *之后* 调用，用于填充未被其他方式覆盖的字段。
// (applyDefaultsToZeroFields recursively traverses the target struct `target`,
// and for fields that are zero-valued (or nil pointers) and have a `default` tag,
// it attempts to parse the tag value and set the field.)
// (This function is designed to be called *after* Viper/mapstructure decoding
// to fill fields not overridden by other means.)
// Parameters:
//   target: 指向要应用默认值的结构体的非 nil 指针。
//           (A non-nil pointer to the struct to apply defaults to.)
// Returns:
//   error: 在递归或应用默认值过程中发生的任何错误。
//          (Any error occurring during recursion or default application.)
func applyDefaultsToZeroFields(target interface{}) error {
	val := reflect.ValueOf(target)

	// Must be a non-nil pointer to a struct
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return fmt.Errorf("target must be a non-nil pointer, got %T", target)
	}
	elem := val.Elem()
	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to a struct, got pointer to %s", elem.Kind())
	}
	typ := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		fieldVal := elem.Field(i)
		fieldType := typ.Field(i)

		if !fieldType.IsExported() {
			continue
		}

		// Recurse first for nested structs (pointer or value)
		if fieldVal.Kind() == reflect.Struct && fieldVal.CanAddr() {
			if err := applyDefaultsToZeroFields(fieldVal.Addr().Interface()); err != nil {
				return fmt.Errorf("error applying defaults to nested struct field %s: %w", fieldType.Name, err)
			}
		} else if fieldVal.Kind() == reflect.Ptr && fieldVal.Type().Elem().Kind() == reflect.Struct {
			// Ensure pointer is initialized (should be by initializeNilPointers)
			if !fieldVal.IsNil() {
				if err := applyDefaultsToZeroFields(fieldVal.Interface()); err != nil {
					return fmt.Errorf("error applying defaults to nested pointer field %s: %w", fieldType.Name, err)
				}
			}
			// If nil, we can't set defaults inside it. initializeNilPointers should handle creation.
			// We don't need to create it here again just for defaults.
		}

		// Apply default if the field is currently zero and has a default tag
		defaultValueTag := fieldType.Tag.Get("default")
		if defaultValueTag != "" && fieldVal.CanSet() {
			isZero := false
			// Use IsZero() for most types, check IsNil() for pointers
			if fieldVal.Kind() == reflect.Ptr {
				isZero = fieldVal.IsNil()
			} else {
				// We need to be careful with IsZero for complex types like time.Time or structs.
				// Relying on IsZero might be okay if we assume defaults are for basic types, slices, time.Duration.
				isZero = fieldVal.IsZero()
			}

			if isZero {
				// Parse the default value string to the field's type
				targetFieldType := fieldType.Type // Get the type of the field itself
				parsedValue, err := parseStringToType(defaultValueTag, targetFieldType)
				if err != nil {
					log.Printf("Warning: applyDefaultsToZeroFields: Failed to parse default value for field %s from tag '%s': %v. Skipping default.", fieldType.Name, defaultValueTag, err)
					continue
				}

				parsedReflectValue := reflect.ValueOf(parsedValue)

				// Handle setting pointer vs non-pointer fields
				if fieldVal.Kind() == reflect.Ptr {
					// Ensure the parsed value's type matches the pointer's element type
					if fieldVal.Type().Elem() == parsedReflectValue.Type() {
						// Create a new pointer to the parsed value and set the field
						newPtr := reflect.New(fieldVal.Type().Elem())
						newPtr.Elem().Set(parsedReflectValue)
						fieldVal.Set(newPtr)
					} else {
						// This might happen if parseStringToType returns a value type for a pointer field request (e.g., int for *int).
						// Let's try setting the element if possible.
						if fieldVal.IsNil() { // Defensive check, should have been initialized
							log.Printf("Warning: applyDefaultsToZeroFields: Cannot set nil pointer field %s directly with non-pointer default value.", fieldType.Name)
							continue
						}
						if fieldVal.Elem().CanSet() && fieldVal.Elem().Type() == parsedReflectValue.Type() {
							fieldVal.Elem().Set(parsedReflectValue)
						} else {
							log.Printf("Warning: applyDefaultsToZeroFields: Type mismatch for pointer field %s: default tag parsed type %T, field element type %s. Skipping.",
								fieldType.Name, parsedValue, fieldVal.Type().Elem())
						}
					}
				} else { // Field is not a pointer
					// Ensure the parsed value's type can be assigned to the field's type
					if parsedReflectValue.Type().AssignableTo(fieldVal.Type()) {
						fieldVal.Set(parsedReflectValue)
					} else {
						// Attempt conversion if possible (e.g., int64 to int)
						if parsedReflectValue.CanConvert(fieldVal.Type()) {
							convertedValue := parsedReflectValue.Convert(fieldVal.Type())
							fieldVal.Set(convertedValue)
						} else {
							log.Printf("Warning: applyDefaultsToZeroFields: Type mismatch for field %s: default tag parsed type %T, field type %s. Skipping.",
								fieldType.Name, parsedValue, fieldVal.Type())
						}
					}
				}
			}
		}
	}
	return nil
}

// DefaultConfigLoader implements the Loader interface for setting defaults.
// (DefaultConfigLoader 实现 Loader 接口用于设置默认值。)
type DefaultConfigLoader struct {
	defaults interface{} // A pointer to the struct defining the defaults (指向定义默认值的结构体的指针)
}

// NewDefaultConfigLoader 创建一个新的 DefaultConfigLoader 实例。
// 这个加载器用于从结构体标签中读取 `default` 值并设置到 Viper 中。
// (NewDefaultConfigLoader creates a new DefaultConfigLoader instance.)
// (This loader is used to read `default` values from struct tags and set them in Viper.)
// Parameters:
//   defaultsStructPtr: 一个指向定义了 `default` 标签的结构体实例的非 nil 指针。
//                      (A non-nil pointer to a struct instance where `default` tags are defined.)
// Returns:
//   *DefaultConfigLoader: 指向新创建的加载器的指针。
//                         (Pointer to the newly created loader.)
func NewDefaultConfigLoader(defaultsStructPtr interface{}) *DefaultConfigLoader {
	// Ensure it's a pointer to a struct
	v := reflect.ValueOf(defaultsStructPtr)
	if v.Kind() != reflect.Ptr || v.IsNil() || v.Elem().Kind() != reflect.Struct {
		log.Panicf("NewDefaultConfigLoader requires a non-nil pointer to a struct, got %T", defaultsStructPtr)
	}
	return &DefaultConfigLoader{defaults: defaultsStructPtr}
}

// Load 使用内部存储的默认配置结构体，调用 setDefaultsFromTags 将默认值设置到提供的 Viper 实例中。
// (Load uses the internally stored default config struct to call setDefaultsFromTags, applying defaults to the provided Viper instance.)
// Parameters:
//   v: 要设置默认值的 Viper 实例。
//      (The Viper instance to set defaults on.)
// Returns:
//   error: 目前总是返回 nil，但保留以符合接口。
//          (Currently always returns nil, but reserved for interface compliance.)
func (l *DefaultConfigLoader) Load(v *viper.Viper) error {
	log.Println("Applying default configuration values from struct tags to Viper...") // Clarify action
	setDefaultsFromTags(v, l.defaults, "")
	log.Println("Finished applying default configuration values to Viper.") // Clarify action
	return nil
}

// Name 返回加载器的名称。
// (Name returns the name of the loader.)
// Returns:
//   string: 加载器的名称 ("DefaultTagLoader")。
//           (The name of the loader ("DefaultTagLoader").)
func (l *DefaultConfigLoader) Name() string {
	return "DefaultTagLoader"
}
