/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package log_test

import (
	"reflect"
	"testing"

	"go.uber.org/zap/zapcore"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

// TestNewOptions 测试 NewOptions 函数是否返回带有正确默认值的 Options 实例。
// (TestNewOptions tests if the NewOptions function returns an Options instance with correct default values.)
func TestNewOptions(t *testing.T) {
	// 定义预期的默认值，包括日志轮转选项
	// (Define expected default values, including log rotation options)
	expected := &log.Options{
		Level:                zapcore.InfoLevel.String(),
		DisableCaller:        false,
		DisableStacktrace:    false,
		StacktraceLevel:      zapcore.ErrorLevel.String(),
		Format:               log.FormatJSON,
		EnableColor:          false,
		Development:          false,
		OutputPaths:          []string{"stdout"},
		ErrorOutputPaths:     []string{"stderr"},
		TimeFormat:           "",
		EncoderConfig:        nil,
		LogRotateMaxSize:     100, // Default: 100 MB
		LogRotateMaxBackups:  5,   // Default: 5 backups
		LogRotateMaxAge:      7,   // Default: 7 days
		LogRotateCompress:    false, // Default: false
	}

	actual := log.NewOptions()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("NewOptions() = %+v, want %+v", actual, expected)
	}
}

// TestOptions_Validate 测试 Options.Validate 方法。
// (TestOptions_Validate tests the Options.Validate method.)
func TestOptions_Validate(t *testing.T) {
	// 当前 Validate 方法为空，预期返回 nil 错误。
	// (Currently, the Validate method is empty, expecting nil error.)
	opts := log.NewOptions()
	errs := opts.Validate()

	if len(errs) != 0 {
		t.Errorf("Validate() returned unexpected errors: %v", errs)
	}

	// 在 Validate 实现验证逻辑后，需要添加更多测试用例。
	// (More test cases should be added after implementing validation logic in Validate.)
	// 例如：测试无效的 Level 值 (e.g., test invalid Level value)
	// opts.Level = "invalid-level"
	// errs = opts.Validate()
	// if len(errs) == 0 {
	// 	 t.Errorf("Validate() should return an error for invalid level, but got none")
	// }
}
