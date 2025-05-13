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

	localAssert.Contains(lines[0], `"message":"Log with trace, request, and custom ID"`)
	localAssert.Contains(lines[0], fmt.Sprintf(`"trace_id":"%s"`, testTraceID))
	localAssert.Contains(lines[0], fmt.Sprintf(`"request_id":"%s"`, testRequestID))
	localAssert.Contains(lines[0], fmt.Sprintf(`"%s":"%s"`, string(customKey), customValue))

	localAssert.Contains(lines[1], `"message":"Log with only trace ID warning"`)
	localAssert.Contains(lines[1], fmt.Sprintf(`"trace_id":"%s"`, testTraceID))
	localAssert.NotContains(lines[1], `"request_id":`)
	localAssert.NotContains(lines[1], fmt.Sprintf(`"%s":`, string(customKey)))

	localAssert.Contains(lines[2], `"message":"Log with only custom key"`)
	localAssert.NotContains(lines[2], `"trace_id":`)
	localAssert.NotContains(lines[2], `"request_id":`)
	localAssert.Contains(lines[2], fmt.Sprintf(`"%s":"%s"`, string(customKey), customValue))

	localAssert.Contains(lines[3], `"message":"Log with no context values"`)
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

	localAssert.Contains(lines[0], `"message":"Global log with context via Ctx"`)
	localAssert.Contains(lines[0], fmt.Sprintf(`"trace_id":"%s"`, testTraceID))
	localAssert.Contains(lines[0], fmt.Sprintf(`"request_id":"%s"`, testRequestID))
	localAssert.Contains(lines[0], fmt.Sprintf(`"%s":"%s"`, string(globalCustomKey), customValue))

	localAssert.Contains(lines[1], `"message":"Global log with trace ID via Ctxf: formatted"`)
	localAssert.Contains(lines[1], fmt.Sprintf(`"trace_id":"%s"`, testTraceID))
	localAssert.NotContains(lines[1], `"request_id":`)
	localAssert.NotContains(lines[1], fmt.Sprintf(`"%s":`, string(globalCustomKey)))

	localAssert.Contains(lines[2], `"message":"Global log with custom key via Ctxw"`)
	localAssert.Contains(lines[2], `"extra":"field"`)
	localAssert.NotContains(lines[2], `"trace_id":`)
	localAssert.NotContains(lines[2], `"request_id":`)
	localAssert.Contains(lines[2], fmt.Sprintf(`"%s":"%s"`, string(globalCustomKey), customValue))
} 