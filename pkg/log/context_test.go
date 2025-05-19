/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * Contains tests for context-aware logging functionality.
 */

package log_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

// TestCtxFunctions tests logging with context extraction for logger instances.
// (TestCtxFunctions 测试 logger 实例的上下文提取日志记录。)
func TestCtxFunctions(t *testing.T) {
	localRequire := require.New(t)
	localAssert := assert.New(t)

	tempDir := t.TempDir()
	logFilePath := filepath.Join(tempDir, "test_ctx.log")
	opts := log.NewOptions()
	opts.OutputPaths = []string{logFilePath}
	opts.Format = log.FormatJSON 
	opts.Level = zapcore.DebugLevel.String()
	type customCtxKey string
	const customKey customCtxKey = "custom_key"
	opts.ContextKeys = []any{customKey}

	logger := log.NewLogger(opts)
	localRequire.NotNil(logger)
	defer func() { _ = logger.Sync() }()

	testTraceID := "trace-12345"
	testRequestID := "req-67890"
	customValue := "my-custom-data"

	ctx1 := context.Background()
	ctx1 = log.ContextWithTraceID(ctx1, testTraceID)
	ctx1 = log.ContextWithRequestID(ctx1, testRequestID)
	ctx1 = context.WithValue(ctx1, customKey, customValue)
	logger.CtxInfof(ctx1, "Log with trace, request, and custom ID")

	ctx2 := context.Background()
	ctx2 = log.ContextWithTraceID(ctx2, testTraceID)
	logger.CtxWarnf(ctx2, "Log with only trace ID %s", "warning")

	ctx3 := context.Background()
	ctx3 = context.WithValue(ctx3, customKey, customValue)
	logger.CtxErrorf(ctx3, "Log with only custom key")

	ctx4 := context.Background()
	logger.CtxDebugf(ctx4, "Log with no context values")

	err := logger.Sync()
	localRequire.NoError(err)

	contentBytes, err := os.ReadFile(logFilePath)
	localRequire.NoError(err)
	content := string(contentBytes)
	lines := strings.Split(strings.TrimSpace(content), "\n")
	localRequire.Equal(4, len(lines), "Should have 4 log lines")

	localAssert.Contains(lines[0], `"M":"Log with trace, request, and custom ID"`)
	localAssert.Contains(lines[0], fmt.Sprintf(`"trace_id":"%s"`, testTraceID))
	localAssert.Contains(lines[0], fmt.Sprintf(`"request_id":"%s"`, testRequestID))
	localAssert.Contains(lines[0], fmt.Sprintf(`"%s":"%s"`, string(customKey), customValue))

	localAssert.Contains(lines[1], `"M":"Log with only trace ID warning"`)
	localAssert.Contains(lines[1], fmt.Sprintf(`"trace_id":"%s"`, testTraceID))
	localAssert.NotContains(lines[1], `"request_id":`)
	localAssert.NotContains(lines[1], fmt.Sprintf(`"%s":`, string(customKey)))

	localAssert.Contains(lines[2], `"M":"Log with only custom key"`)
	localAssert.NotContains(lines[2], `"trace_id":`)
	localAssert.NotContains(lines[2], `"request_id":`)
	localAssert.Contains(lines[2], fmt.Sprintf(`"%s":"%s"`, string(customKey), customValue))

	localAssert.Contains(lines[3], `"M":"Log with no context values"`)
	localAssert.NotContains(lines[3], `"trace_id":`)
	localAssert.NotContains(lines[3], `"request_id":`)
	localAssert.NotContains(lines[3], fmt.Sprintf(`"%s":`, string(customKey)))
}

// TestGlobalCtxFunctions tests global logging functions with context extraction.
// (TestGlobalCtxFunctions 测试带有上下文提取的全局日志记录函数。)
func TestGlobalCtxFunctions(t *testing.T) {
	localRequire := require.New(t)
	localAssert := assert.New(t)

	tempDir := t.TempDir()
	logFilePath := filepath.Join(tempDir, "test_global_ctx.log")
	opts := log.NewOptions()
	opts.OutputPaths = []string{logFilePath}
	opts.Format = log.FormatJSON 
	opts.Level = zapcore.InfoLevel.String()
	type globalCustomCtxKey string
	const globalCustomKey globalCustomCtxKey = "global_custom_key"
	opts.ContextKeys = []any{globalCustomKey}

	log.Init(opts)
	defer func() {
		log.Init(log.NewOptions()) 
		_ = log.Sync()
	}()

	testTraceID := "global-trace-id"
	testRequestID := "global-req-id"
	customValue := "global-custom-data"

	ctx1 := context.Background()
	ctx1 = log.ContextWithTraceID(ctx1, testTraceID)
	ctx1 = log.ContextWithRequestID(ctx1, testRequestID)
	ctx1 = context.WithValue(ctx1, globalCustomKey, customValue)
	log.Ctx(ctx1, "Global log with context via Ctx")

	ctx2 := context.Background()
	ctx2 = log.ContextWithTraceID(ctx2, testTraceID)
	log.Ctxf(ctx2, "Global log with trace ID via Ctxf: %s", "formatted")

	ctx3 := context.Background()
	ctx3 = context.WithValue(ctx3, globalCustomKey, customValue)
	log.Ctxw(ctx3, "Global log with custom key via Ctxw", "extra", "field")

	err := log.Sync()
	localRequire.NoError(err)

	contentBytes, err := os.ReadFile(logFilePath)
	localRequire.NoError(err)
	content := string(contentBytes)
	lines := strings.Split(strings.TrimSpace(content), "\n")
	localRequire.Len(lines, 3, "Should have 3 log lines from global context functions")

	localAssert.Contains(lines[0], `"M":"Global log with context via Ctx"`)
	localAssert.Contains(lines[0], fmt.Sprintf(`"trace_id":"%s"`, testTraceID))
	localAssert.Contains(lines[0], fmt.Sprintf(`"request_id":"%s"`, testRequestID))
	localAssert.Contains(lines[0], fmt.Sprintf(`"%s":"%s"`, string(globalCustomKey), customValue))

	localAssert.Contains(lines[1], `"M":"Global log with trace ID via Ctxf: formatted"`)
	localAssert.Contains(lines[1], fmt.Sprintf(`"trace_id":"%s"`, testTraceID))
	localAssert.NotContains(lines[1], `"request_id":`)
	localAssert.NotContains(lines[1], fmt.Sprintf(`"%s":`, string(globalCustomKey)))

	localAssert.Contains(lines[2], `"M":"Global log with custom key via Ctxw"`)
	localAssert.Contains(lines[2], `"extra":"field"`)
	localAssert.NotContains(lines[2], `"trace_id":`)
	localAssert.NotContains(lines[2], `"request_id":`)
	localAssert.Contains(lines[2], fmt.Sprintf(`"%s":"%s"`, string(globalCustomKey), customValue))
}

// TestGlobalCtxLevelFunctions tests global contextual logging functions for various levels.
// (TestGlobalCtxLevelFunctions 测试不同级别的全局上下文日志记录函数。)
func TestGlobalCtxLevelFunctions(t *testing.T) {
	localRequire := require.New(t)
	localAssert := assert.New(t)

	tempDir := t.TempDir()
	logFilePath := filepath.Join(tempDir, "test_global_ctx_levels.log")
	opts := log.NewOptions()
	opts.OutputPaths = []string{logFilePath}
	opts.Format = log.FormatJSON
	opts.Level = zapcore.DebugLevel.String() // Ensure all levels are logged
	type globalCtxLevelKey string
	const ctxKeyForLevels globalCtxLevelKey = "ctx_level_key"
	opts.ContextKeys = []any{ctxKeyForLevels}

	log.Init(opts)
	defer func() {
		// Restore global logger to a known default state
		log.Init(log.NewOptions()) 
		_ = log.Sync()
	}()

	testTraceID := "global-ctx-level-trace"
	testRequestID := "global-ctx-level-req"
	customLevelValue := "ctx-level-data"

	ctx := context.Background()
	ctx = log.ContextWithTraceID(ctx, testTraceID)
	ctx = log.ContextWithRequestID(ctx, testRequestID)
	ctx = context.WithValue(ctx, ctxKeyForLevels, customLevelValue)

	// Test each level
	log.CtxDebugf(ctx, "Global CtxDebugf message: %s", "debug_val")
	log.CtxInfof(ctx, "Global CtxInfof message: %s", "info_val")
	log.CtxWarnf(ctx, "Global CtxWarnf message: %s", "warn_val")
	log.CtxErrorf(ctx, "Global CtxErrorf message: %s", "error_val")

	// Test Panic and Fatal separately due to their nature
	// We can't easily test os.Exit in Fatal, so we'll focus on Panic
	// and trust Fatal calls the underlying logger's Fatal.

	// Test CtxPanicf
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("log.CtxPanicf did not panic")
			}
		}()
		log.CtxPanicf(ctx, "Global CtxPanicf message: %s", "panic_val")
	}()

	// After panic, the logger might be in an inconsistent state or synced. 
	// It's good practice to re-initialize or ensure state if tests continue.
	// For this test, we expect logs up to CtxErrorf and the Panic log to be written.
	// Re-init for safety if further global log operations were planned *after* a panic test in a real scenario.
	// However, since we defer log.Init(log.NewOptions()), this is handled at test end.
	log.Init(opts) // Re-initialize with the same options to ensure the logger is active after panic for sync
	// Ensure that the instance of logger used by `log.Std()` is the same as `originalStdLogger` or a new one created by `log.Init`
	// This is to make sure the global logger `std` is not unexpectedly changed.
	currentStdLogger := log.Std()
	// Depending on internal Init behavior, it might be the same or a new instance.
	// What matters is that it's functional and using `opts`.
	localAssert.NotNil(currentStdLogger, "Global logger should be non-nil after re-init")


	err := log.Sync() // Sync logs
	localRequire.NoError(err, "Failed to sync logger for global context level functions")

	contentBytes, err := os.ReadFile(logFilePath)
	localRequire.NoError(err, "Failed to read log file for global context level functions")
	content := string(contentBytes)
	lines := strings.Split(strings.TrimSpace(content), "\n")

	// Expected number of lines: Debug, Info, Warn, Error, Panic = 5 lines
	// Fatal is not easily testable here due to os.Exit
	localRequire.Len(lines, 5, "Should have 5 log lines from global context level functions (excluding Fatal)")

	expectedMessages := []string{
		"Global CtxDebugf message: debug_val",
		"Global CtxInfof message: info_val",
		"Global CtxWarnf message: warn_val",
		"Global CtxErrorf message: error_val",
		"Global CtxPanicf message: panic_val",
	}
	expectedLevels := []string{"DEBUG", "INFO", "WARN", "ERROR", "PANIC"}

	for i, line := range lines {
		localAssert.Contains(line, fmt.Sprintf(`"M":"%s"`, expectedMessages[i]))
		localAssert.Contains(line, fmt.Sprintf(`"L":"%s"`, expectedLevels[i]))
		localAssert.Contains(line, fmt.Sprintf(`"trace_id":"%s"`, testTraceID))
		localAssert.Contains(line, fmt.Sprintf(`"request_id":"%s"`, testRequestID))
		localAssert.Contains(line, fmt.Sprintf(`"%s":"%s"`, string(ctxKeyForLevels), customLevelValue))
		if expectedLevels[i] == "ERROR" || expectedLevels[i] == "PANIC" {
			localAssert.Contains(line, `"stacktrace":`)
		}
	}
} 