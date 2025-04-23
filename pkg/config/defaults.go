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

// initializeNilPointers 初始化结构体中的 nil 指针字段
// (initializeNilPointers initializes nil pointer fields in the struct)
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
				initializeNilPointers(newStructPtr.Interface())
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

// applyDefaultsFromStructTags 递归地将结构体标签中的默认值直接应用到结构体字段（如果字段为零值）
// (applyDefaultsFromStructTags recursively applies default values from struct tags directly to struct fields if they are zero)
func applyDefaultsFromStructTags(target interface{}) error {
	val := reflect.ValueOf(target)

	// 必须是指针 (Must be a pointer)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		// Allow non-pointer structs passed by value initially, but recurse with pointers
		if val.Kind() == reflect.Struct && val.CanAddr() {
			val = val.Addr() // Get pointer for subsequent checks
			if val.Kind() != reflect.Ptr {
				return nil
			} // Should now be a pointer
		} else {
			return nil // Or return error?
		}
	}

	elem := val.Elem() // Get the struct value

	// 必须指向结构体 (Must point to a struct)
	if elem.Kind() != reflect.Struct {
		return nil
	}

	typ := elem.Type()

	for i := 0; i < elem.NumField(); i++ {
		fieldVal := elem.Field(i)
		fieldType := typ.Field(i)

		// 跳过未导出的字段 (Skip unexported fields)
		if !fieldType.IsExported() {
			continue
		}

		// 递归处理嵌套结构体 (指针或值类型)
		// (Recursively handle nested structs - pointer or value type)
		// 需要在应用当前字段的默认值 *之前* 递归，以处理内层默认值
		// (Need to recurse *before* applying default to current field to handle inner defaults first)
		if fieldVal.Kind() == reflect.Ptr && fieldVal.Type().Elem().Kind() == reflect.Struct {
			// 确保指针已初始化 (Ensure pointer is initialized - initializeNilPointers should have done this)
			if !fieldVal.IsNil() {
				if err := applyDefaultsFromStructTags(fieldVal.Interface()); err != nil {
					return fmt.Errorf("error applying defaults to pointer field %s: %w", fieldType.Name, err)
				}
			}
		} else if fieldVal.Kind() == reflect.Struct {
			if fieldVal.CanAddr() {
				if err := applyDefaultsFromStructTags(fieldVal.Addr().Interface()); err != nil {
					return fmt.Errorf("error applying defaults to value field %s: %w", fieldType.Name, err)
				}
			}
		}

		// 处理当前字段的默认值标签 (Process the default value tag for the current field)
		defaultValueTag := fieldType.Tag.Get("default")
		if defaultValueTag != "" {
			if fieldVal.CanSet() {
				isZero := false
				// 使用 reflect.Value.IsZero() 来判断是否为零值
				// (Use reflect.Value.IsZero() to check if it's the zero value)
				// 对于指针，我们需要检查它是否为 nil
				if fieldVal.Kind() == reflect.Ptr {
					isZero = fieldVal.IsNil()
				} else {
					// 对于非指针类型，IsZero() 是可靠的
					isZero = fieldVal.IsZero()
				}

				if isZero {
					typedValue, err := parseDefaultValueToType(fieldVal, defaultValueTag)
					if err != nil {
						log.Printf("Warning: Failed to parse default value for field %s from tag '%s': %v. Skipping default application.", fieldType.Name, defaultValueTag, err)
						continue // Skip setting default if parsing fails
					}

					// 需要创建一个正确类型的值来设置 (Need to create a value of the correct type to set)
					parsedVal := reflect.ValueOf(typedValue)

					// 如果字段是指针类型，需要创建指针 (If the field is a pointer type, need to create a pointer)
					if fieldVal.Kind() == reflect.Ptr {
						// 确保解析后的值类型与指针元素类型匹配
						if fieldVal.Type().Elem() == parsedVal.Type() {
							// 创建一个新的指针指向解析后的值
							newPtr := reflect.New(fieldVal.Type().Elem())
							newPtr.Elem().Set(parsedVal)
							fieldVal.Set(newPtr)
						} else {
							log.Printf("Warning: Parsed default value type (%T) does not match pointer element type (%s) for field %s. Skipping.", typedValue, fieldVal.Type().Elem(), fieldType.Name)
						}
					} else {
						// 如果字段不是指针，直接设置值 (If the field is not a pointer, set the value directly)
						// 确保解析后的值类型与字段类型匹配
						if fieldVal.Type() == parsedVal.Type() {
							fieldVal.Set(parsedVal)
						} else {
							log.Printf("Warning: Parsed default value type (%T) does not match field type (%s) for field %s. Skipping.", typedValue, fieldVal.Type(), fieldType.Name)
						}
					}

				} else {
				}
			} else {
			}
		}
	}
	return nil
}

// parseDefaultValueToType 尝试将字符串默认值解析为字段的实际类型
// (parseDefaultValueToType tries to parse a string default value into the field's actual type)
func parseDefaultValueToType(field reflect.Value, defaultValue string) (interface{}, error) {
	kind := field.Kind()
	fieldType := field.Type()

	// 处理指针类型 (Handle pointer types - should primarily look at the element type)
	if kind == reflect.Ptr {
		// If the field is a nil pointer, create a zero value of the element type to determine its kind
		if field.IsNil() {
			fieldType = field.Type().Elem()
			kind = field.Type().Elem().Kind()
		} else {
			// If not nil, get the kind of the element it points to
			fieldType = field.Elem().Type()
			kind = field.Elem().Kind()
		}
	}

	switch kind {
	case reflect.String:
		return defaultValue, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if fieldType == reflect.TypeOf(time.Duration(0)) {
			parsedDuration, err := time.ParseDuration(defaultValue)
			if err != nil {
				return nil, fmt.Errorf("invalid time duration format '%s': %w", defaultValue, err)
			}
			return parsedDuration, nil // Return time.Duration directly
		} else {
			val, err := strconv.ParseInt(defaultValue, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid integer format '%s': %w", defaultValue, err)
			}
			// Create a value of the correct integer type
			resultVal := reflect.New(fieldType).Elem()
			if resultVal.OverflowInt(val) {
				return nil, fmt.Errorf("integer overflow for '%s'", defaultValue)
			}
			resultVal.SetInt(val)
			return resultVal.Interface(), nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(defaultValue, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid unsigned integer format '%s': %w", defaultValue, err)
		}
		resultVal := reflect.New(fieldType).Elem()
		if resultVal.OverflowUint(val) {
			return nil, fmt.Errorf("unsigned integer overflow for '%s'", defaultValue)
		}
		resultVal.SetUint(val)
		return resultVal.Interface(), nil
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(defaultValue, fieldType.Bits())
		if err != nil {
			return nil, fmt.Errorf("invalid float format '%s': %w", defaultValue, err)
		}
		resultVal := reflect.New(fieldType).Elem()
		if resultVal.OverflowFloat(val) {
			return nil, fmt.Errorf("float overflow for '%s'", defaultValue)
		}
		resultVal.SetFloat(val)
		return resultVal.Interface(), nil
	case reflect.Bool:
		val, err := strconv.ParseBool(defaultValue)
		if err != nil {
			return nil, fmt.Errorf("invalid boolean format '%s': %w", defaultValue, err)
		}
		return val, nil
	case reflect.Slice:
		// Handle string slices specifically, common use case
		if fieldType.Elem().Kind() == reflect.String {
			if defaultValue == "" {
				return []string{}, nil // Empty string means empty slice
			}
			return strings.Split(defaultValue, ","), nil // Split by comma
		}
		return nil, fmt.Errorf("default values for non-string slices (type %v) are not supported via tags", fieldType)
	default:
		return nil, fmt.Errorf("unsupported type '%v' for default value parsing from tag", kind)
	}
}

// setDefaultsFromTags 递归地遍历结构体，读取 'default' 标签，并将它们设置到 Viper 实例中
// (setDefaultsFromTags recursively traverses the struct, reads 'default' tags, and sets them in the Viper instance)
func setDefaultsFromTags(v *viper.Viper, prefix string, obj interface{}) error {
	val := reflect.ValueOf(obj)

	// 如果是指针，获取它指向的值 (If it's a pointer, get the value it points to)
	if val.Kind() == reflect.Ptr {
		// 如果指针为 nil，无法设置默认值 (If the pointer is nil, cannot set defaults)
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}

	// 必须是结构体才能继续 (Must be a struct to proceed)
	if val.Kind() != reflect.Struct {
		return nil
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		// 获取 mapstructure 标签名，用于构建 Viper key
		// (Get mapstructure tag name to build Viper key)
		mapKey := field.Tag.Get("mapstructure")
		if mapKey == "" || mapKey == "-" {
			mapKey = field.Name // Use field name if mapstructure tag is missing or "-"
		}

		// 处理 ",squash"
		// (Handle ",squash")
		parts := strings.Split(mapKey, ",")
		mapKey = parts[0]
		isSquash := false
		for _, part := range parts {
			if part == "squash" {
				isSquash = true
				break
			}
		}

		// 构建 Viper key (Build Viper key)
		currentKey := mapKey
		if prefix != "" && !isSquash {
			currentKey = prefix + "." + mapKey
		} else if prefix != "" && isSquash {
			currentKey = prefix // For squashed fields, use the parent prefix
		}

		// 如果字段是匿名的（嵌入式结构体），递归处理，不添加字段名到 key
		// (If field is anonymous (embedded struct), recurse without adding field name to key)
		if field.Anonymous && fieldValue.Kind() == reflect.Struct {
			if err := setDefaultsFromTags(v, prefix, fieldValue.Addr().Interface()); err != nil {
				return err
			}
			continue // 继续下一个字段 (Continue to the next field)
		}

		// 检查 'default' 标签 (Check for 'default' tag)
		defaultTag := field.Tag.Get("default")

		// 递归处理嵌套结构体或指针 (Recurse into nested structs or pointers)
		fieldType := fieldValue.Type()
		if fieldType.Kind() == reflect.Struct {
			if err := setDefaultsFromTags(v, currentKey, fieldValue.Addr().Interface()); err != nil {
				return err
			}
		} else if fieldType.Kind() == reflect.Ptr && fieldType.Elem().Kind() == reflect.Struct {
			// 如果指针为 nil，尝试初始化它 (If pointer is nil, try initializing it)
			if fieldValue.IsNil() {
				fieldValue.Set(reflect.New(fieldType.Elem()))
			}
			if err := setDefaultsFromTags(v, currentKey, fieldValue.Interface()); err != nil {
				return err
			}
		} else if defaultTag != "" {
			// 只为非结构体/指针且有 default 标签的字段设置默认值
			// (Only set default for non-struct/pointer fields with a default tag)

			// 检查 Viper 中是否已设置该值 (Check if the value is already set in Viper)
			// 注意：这需要在 ReadInConfig 和 AutomaticEnv 之后调用，以避免覆盖文件/环境变量的值
			// (Note: This needs to be called AFTER ReadInConfig and AutomaticEnv to avoid overwriting file/env values)
			if v.IsSet(currentKey) {
				continue
			}

			// 尝试将 'default' 标签字符串转换为字段的实际类型
			// (Try converting the 'default' tag string to the field's actual type)
			var defaultValue interface{}
			var err error
			switch fieldValue.Kind() {
			case reflect.String:
				defaultValue = defaultTag
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				// 特别处理 time.Duration，因为它底层是 int64
				// (Special handling for time.Duration as its underlying type is int64)
				if fieldValue.Type() == reflect.TypeOf(time.Duration(0)) {
					defaultValue, err = time.ParseDuration(defaultTag)
					if err != nil {
						log.Printf("Warning: Failed to parse duration default tag '%s' for key '%s': %v. Skipping default.", defaultTag, currentKey, err)
						continue
					}
				} else {
					defaultValue, err = strconv.ParseInt(defaultTag, 0, 64)
					if err != nil {
						log.Printf("Warning: Failed to parse int default tag '%s' for key '%s': %v. Skipping default.", defaultTag, currentKey, err)
						continue
					}
					// 需要转换为字段的具体 int 类型 (Need to convert to the specific int type of the field)
					defaultValue = reflect.ValueOf(defaultValue).Convert(fieldValue.Type()).Interface()
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				defaultValue, err = strconv.ParseUint(defaultTag, 0, 64)
				if err != nil {
					log.Printf("Warning: Failed to parse uint default tag '%s' for key '%s': %v. Skipping default.", defaultTag, currentKey, err)
					continue
				}
				defaultValue = reflect.ValueOf(defaultValue).Convert(fieldValue.Type()).Interface()
			case reflect.Float32, reflect.Float64:
				defaultValue, err = strconv.ParseFloat(defaultTag, 64)
				if err != nil {
					log.Printf("Warning: Failed to parse float default tag '%s' for key '%s': %v. Skipping default.", defaultTag, currentKey, err)
					continue
				}
				defaultValue = reflect.ValueOf(defaultValue).Convert(fieldValue.Type()).Interface()
			case reflect.Bool:
				defaultValue, err = strconv.ParseBool(defaultTag)
				if err != nil {
					log.Printf("Warning: Failed to parse bool default tag '%s' for key '%s': %v. Skipping default.", defaultTag, currentKey, err)
					continue
				}
			case reflect.Slice:
				// 对于切片，我们假设默认值是逗号分隔的字符串
				// (For slices, assume the default is a comma-separated string)
				// 注意：这不会自动转换元素类型，它会将整个字符串切片存储起来
				// (Note: This doesn't auto-convert element types, it stores the string slice)
				// 这与 mapstructure 的 StringToSliceHookFunc 行为不同
				// (This differs from mapstructure's StringToSliceHookFunc behavior)
				// 可能需要更复杂的逻辑来处理不同类型的切片元素
				// (More complex logic might be needed for different slice element types)
				if fieldType.Elem().Kind() == reflect.String {
					defaultValue = strings.Split(defaultTag, ",")
					// 去除可能存在的空格 (Trim potential spaces)
					stringSlice := defaultValue.([]string)
					for j := range stringSlice {
						stringSlice[j] = strings.TrimSpace(stringSlice[j])
					}
					defaultValue = stringSlice
				} else {
					log.Printf("Warning: Default tag for non-string slice key '%s' is not supported yet. Skipping default.", currentKey)
					continue
				}
			default:
				log.Printf("Warning: Default tag for type '%s' (key: '%s') is not supported yet. Skipping default.", fieldValue.Kind(), currentKey)
				continue
			}

			// 设置默认值 (Set the default value)
			v.SetDefault(currentKey, defaultValue)
		}
	}
	return nil
}
