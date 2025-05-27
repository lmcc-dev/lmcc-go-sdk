/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * Contains tests for default value parsing and application logic.
 */

package config

import (
	"reflect"
	"testing"
	"time"

	// lmccerrors
	stdErrors "errors" // Standard library errors for IsCode replacement

	lmccerrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParseStringToType tests the internal default value parsing logic.
// 测试默认值解析逻辑
func TestParseStringToType(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		targetType  reflect.Type
		expectedVal interface{}
		expectError bool
	}{
		{"String", "hello", reflect.TypeOf(""), "hello", false},
		{"Int", "123", reflect.TypeOf(0), int(123), false},
		{"Int64", "-456", reflect.TypeOf(int64(0)), int64(-456), false},
		{"Uint", "789", reflect.TypeOf(uint(0)), uint(789), false},
		{"Float32", "1.23", reflect.TypeOf(float32(0)), float32(1.23), false},
		{"Float64", "-4.56", reflect.TypeOf(float64(0)), float64(-4.56), false},
		{"BoolTrue", "true", reflect.TypeOf(false), true, false},
		{"BoolFalse", "false", reflect.TypeOf(false), false, false},
		{"Duration", "5s", reflect.TypeOf(time.Duration(0)), 5 * time.Second, false},
		{"DurationComplex", "1h30m", reflect.TypeOf(time.Duration(0)), (1 * time.Hour) + (30 * time.Minute), false},
		{"StringSliceComma", "a,b, c", reflect.TypeOf([]string{}), []string{"a", "b", "c"}, false},
		{"StringSliceSpace", "a b  c ", reflect.TypeOf([]string{}), []string{"a", "b", "c"}, false},
		{"StringSliceEmpty", "", reflect.TypeOf([]string{}), []string{}, false},
		{"PointerString", "ptr_hello", reflect.TypeOf(new(string)), "ptr_hello", false}, // Should parse to element type
		{"PointerInt", "99", reflect.TypeOf(new(int)), int(99), false},                 // Should parse to element type
		{"InvalidInt", "abc", reflect.TypeOf(0), nil, true},
		{"InvalidFloat", "def", reflect.TypeOf(float64(0)), nil, true},
		{"InvalidBool", "maybe", reflect.TypeOf(false), nil, true},
		{"InvalidDuration", "xyz", reflect.TypeOf(time.Duration(0)), nil, true},
		{"UnsupportedSlice", "1,2,3", reflect.TypeOf([]int{}), nil, true}, // Only string slices supported for default tags
		{"UnsupportedType", "{}", reflect.TypeOf(struct{}{}), nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedVal, err := parseStringToType(tt.value, tt.targetType)
			if tt.expectError {
				assert.Error(t, err, "Expected an error for %s", tt.name)
			} else {
				assert.NoError(t, err, "Did not expect an error for %s", tt.name)
				// Use assert.Equal because reflect types can be tricky
				assert.Equal(t, tt.expectedVal, parsedVal, "Parsed value mismatch for %s", tt.name)
			}
		})
	}
}

// TestApplyDefaultsToZeroFields tests applying defaults after initial load.
// 测试应用默认值到零值字段
func TestApplyDefaultsToZeroFields(t *testing.T) {
	type NestedDefaults struct {
		NestedVal string `default:"nested_default"`
		NestedInt *int   `default:"55"`
	}
	type ConfigWithDefaults struct {
		StrField    string          `default:"default_string"`
		IntField    int             `default:"42"`
		BoolField   bool            `default:"true"`
		DurField    time.Duration   `default:"3m"`
		SliceField  []string        `default:"one, two"`
		PtrField    *string         `default:"ptr_default"`
		NestedField *NestedDefaults // Field itself has no default, but members do
		ZeroDur     time.Duration   // No default tag, should remain zero
		NoDefault   string
	}

	// Scenario 1: Apply defaults to a completely zero struct
	cfg1 := &ConfigWithDefaults{}
	initializeNilPointers(cfg1) // Important to initialize nested pointers
	err := applyDefaultsToZeroFields(cfg1)
	require.NoError(t, err)

	assert.Equal(t, "default_string", cfg1.StrField)
	assert.Equal(t, 42, cfg1.IntField)
	assert.True(t, cfg1.BoolField)
	assert.Equal(t, 3*time.Minute, cfg1.DurField)
	assert.Equal(t, []string{"one", "two"}, cfg1.SliceField)
	require.NotNil(t, cfg1.PtrField)
	assert.Equal(t, "ptr_default", *cfg1.PtrField)
	require.NotNil(t, cfg1.NestedField)
	assert.Equal(t, "nested_default", cfg1.NestedField.NestedVal)
	require.NotNil(t, cfg1.NestedField.NestedInt)
	assert.Equal(t, 55, *cfg1.NestedField.NestedInt)
	assert.Equal(t, time.Duration(0), cfg1.ZeroDur)
	assert.Equal(t, "", cfg1.NoDefault)

	// Scenario 2: Apply defaults when some fields are already set (should not overwrite)
	initialStr := "initial"
	initialInt := 99
	cfg2 := &ConfigWithDefaults{
		IntField:  initialInt, // Pre-set
		PtrField:  &initialStr, // Pre-set
		NoDefault: "preset",
	}
	initializeNilPointers(cfg2)
	err = applyDefaultsToZeroFields(cfg2)
	require.NoError(t, err)

	assert.Equal(t, "default_string", cfg2.StrField) // Was zero, got default
	assert.Equal(t, initialInt, cfg2.IntField)       // Was set, kept initial
	assert.True(t, cfg2.BoolField)                   // Was zero, got default
	assert.Equal(t, 3*time.Minute, cfg2.DurField)    // Was zero, got default
	assert.Equal(t, []string{"one", "two"}, cfg2.SliceField) // Was zero, got default
	require.NotNil(t, cfg2.PtrField)
	assert.Equal(t, initialStr, *cfg2.PtrField) // Was set, kept initial
	require.NotNil(t, cfg2.NestedField)
	assert.Equal(t, "nested_default", cfg2.NestedField.NestedVal) // Was zero, got default
	require.NotNil(t, cfg2.NestedField.NestedInt)
	assert.Equal(t, 55, *cfg2.NestedField.NestedInt) // Was zero, got default
	assert.Equal(t, time.Duration(0), cfg2.ZeroDur)  // Still zero
	assert.Equal(t, "preset", cfg2.NoDefault)        // Was set, kept initial
}

// TestDefaultConfigLoader tests the loader for applying defaults from struct tags.
// 测试默认值加载器
func TestDefaultConfigLoader(t *testing.T) {
	type DefaultsStruct struct {
		Host string `default:"defaulthost.com"`
		Port int    `default:"8080"`
	}
	defaults := &DefaultsStruct{}
	loader := NewDefaultConfigLoader(defaults)
	require.NotNil(t, loader)
	assert.Equal(t, "DefaultTagLoader", loader.Name())

	v := viper.New()
	err := loader.Load(v) // Apply defaults to viper instance
	require.NoError(t, err)

	// Check if viper instance got the defaults
	assert.Equal(t, "defaulthost.com", v.GetString("Host")) // Viper keys might be case-insensitive depending on usage
	assert.Equal(t, 8080, v.GetInt("Port"))

	// Test panic on invalid input
	assert.Panics(t, func() { NewDefaultConfigLoader(DefaultsStruct{}) }, "Should panic if not a pointer")
	assert.Panics(t, func() { NewDefaultConfigLoader((*DefaultsStruct)(nil)) }, "Should panic if nil pointer")
	var notAStruct *int
	assert.Panics(t, func() { NewDefaultConfigLoader(&notAStruct) }, "Should panic if not pointer to struct")
}

// TestApplyDefaultsToZeroFields_MoreTypes tests more complex types and edge cases.
// 测试应用默认值到零值字段（更多类型和边缘情况）
func TestApplyDefaultsToZeroFields_MoreTypes(t *testing.T) {
	type MoreDefaults struct {
		// Basic Types (already tested in primary test)
		Int8Field   int8    `default:"-10"`
		Uint16Field uint16  `default:"65000"`
		Float32Field float32 `default:"3.14"`

		// Pointer to Basic Types
		IntPtr     *int     `default:"12345"`
		BoolPtr    *bool    `default:"true"`
		StringPtr  *string  // No default tag, should remain nil
		NilPtrWithDefault *int `default:"999"` // Pointer is nil, but has default

		// Slice (already tested)

		// Struct (already tested)

		// Time (Zero value handling)
		TimeField time.Time // No default, should remain zero

		// Fields already set (should not be overwritten)
		PresetInt int `default:"1"`
		PresetStr string `default:"abc"`
	}

	presetVal := 500
	cfg := &MoreDefaults{
		PresetInt: presetVal,
		PresetStr: "initial",
	}

	initializeNilPointers(cfg) // Ensure necessary pointers are created if applicable (like nested structs, though none here)
	err := applyDefaultsToZeroFields(cfg)
	require.NoError(t, err)

	// Check basic types
	assert.Equal(t, int8(-10), cfg.Int8Field)
	assert.Equal(t, uint16(65000), cfg.Uint16Field)
	assert.Equal(t, float32(3.14), cfg.Float32Field)

	// Check pointers to basic types
	require.NotNil(t, cfg.IntPtr)
	assert.Equal(t, 12345, *cfg.IntPtr)
	require.NotNil(t, cfg.BoolPtr)
	assert.Equal(t, true, *cfg.BoolPtr)
	assert.Nil(t, cfg.StringPtr, "StringPtr should remain nil as it had no default tag")
	require.NotNil(t, cfg.NilPtrWithDefault, "NilPtrWithDefault should be initialized by applyDefaults")
	assert.Equal(t, 999, *cfg.NilPtrWithDefault)

	// Check time zero value
	assert.True(t, cfg.TimeField.IsZero(), "TimeField without default should remain zero")

	// Check preset fields
	assert.Equal(t, presetVal, cfg.PresetInt, "PresetInt should not be overwritten by default")
	assert.Equal(t, "initial", cfg.PresetStr, "PresetStr should not be overwritten by default")

}

// TestSetDefaultsFromTags_EdgeCases tests edge cases for setting defaults in Viper.
// 测试从标签设置默认值到 Viper 的边缘情况
func TestSetDefaultsFromTags_EdgeCases(t *testing.T) {
	type NestedWithSkip struct {
		Keep   string `mapstructure:"keep" default:"nested_keep"`
		SkipMe string `mapstructure:"-"`
	}
	type SkipOmitEmpty struct {
		Field1    string `mapstructure:"field1,omitempty" default:"omit_default"`
		Field2    string `mapstructure:"-"`                        // Skipped field
		Field3    string `mapstructure:"field3" default:"no_omit"`
		NestedPtr *NestedWithSkip
		NestedVal NestedWithSkip
		NilPtr    *NestedWithSkip // Nil pointer to struct with defaults
	}

	cfg := &SkipOmitEmpty{
		// NestedPtr is nil initially
		// NilPtr is nil initially
	}

	v := viper.New()
	err := setDefaultsFromTags(v, cfg, "prefix") // Use a prefix
	require.NoError(t, err)

	// Check basic fields with prefix
	assert.Equal(t, "omit_default", v.GetString("prefix.field1"))
	assert.Equal(t, "no_omit", v.GetString("prefix.field3"))
	assert.False(t, v.IsSet("prefix.field2"), "Skipped field Field2 should not be set")
	assert.False(t, v.IsSet("prefix.SkipMe"), "Skipped field SkipMe should not be set") // mapstructure tag '-' takes precedence

	// Check nested struct values (value type)
	assert.Equal(t, "nested_keep", v.GetString("prefix.NestedVal.keep"))
	assert.False(t, v.IsSet("prefix.NestedVal.SkipMe"))

	// Check for NestedPtr - if it was nil, its defaults shouldn't be set
	assert.False(t, v.IsSet("prefix.NestedPtr.keep"), "Defaults for nil NestedPtr should not be set")

	// Check for NilPtr - if it was nil, its defaults shouldn't be set
	assert.False(t, v.IsSet("prefix.NilPtr.keep"), "Defaults for nil NilPtr should not be set")

	// Initialize NestedPtr and re-run to check its defaults
	cfg.NestedPtr = &NestedWithSkip{}
	v2 := viper.New()
	err = setDefaultsFromTags(v2, cfg, "prefix2") // Use a new viper instance and prefix
	require.NoError(t, err)
	assert.Equal(t, "nested_keep", v2.GetString("prefix2.NestedPtr.keep"))
}

func TestSetDefaultsFromTags_ErrorCase(t *testing.T) {
	type InvalidTagStruct struct {
		GoodField string    `default:"good_value"`
		BadField  int       `default:"this-is-not-an-integer"`
		NextField float64   `default:"3.14"`
	}

	cfg := &InvalidTagStruct{}
	v := viper.New()

	err := setDefaultsFromTags(v, cfg, "") // No prefix
	require.Error(t, err, "setDefaultsFromTags should return an error for invalid default tag")
	assert.True(t, stdErrors.Is(err, lmccerrors.ErrConfigDefaultTagParse), "Error code should be ErrConfigDefaultTagParse")

	// Verify that fields before the error were set, and fields after (and the error field) were not.
	assert.Equal(t, "good_value", v.GetString("GoodField"), "GoodField should be set in Viper")
	assert.False(t, v.IsSet("BadField"), "BadField with invalid default should not be set in Viper")
	assert.False(t, v.IsSet("NextField"), "NextField after error should not be set in Viper")
}