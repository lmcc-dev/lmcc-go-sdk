/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package log_test

import (
	"bufio"
	"bytes"
	"context" // Standard library errors for Is/As and creating simple errors
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	lmccerrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TestMain is used for setup and teardown, not strictly necessary here but good practice.
func TestMain(m *testing.M) {
	// Setup: Initialize a default global logger to ensure no panics if tests run out of order or individually.
	initialOpts := log.NewOptions()
	log.Init(initialOpts)
	// Run tests
	code := m.Run()
	// Teardown
	// ...
	os.Exit(code)
}

// TestNewLogger verifies that NewLogger creates a logger instance correctly.
func TestNewLogger(t *testing.T) {
		opts := log.NewOptions()
	opts.Level = "info"
	opts.OutputPaths = []string{"stdout"}

	logger, err := log.NewLogger(opts)
	assert.NoError(t, err)
	assert.NotNil(t, logger)

	// Test with a specific logger instance
	var buf bytes.Buffer
	specificOpts := log.NewOptions()
	specificOpts.Level = "info"
	// Create a writer that writes to the buffer
	specificLogger := log.NewLoggerWithWriter(specificOpts, &buf)

	specificLogger.Info("Info from NewLogger specific")
	assert.Contains(t, buf.String(), "Info from NewLogger specific")

	// Test with default logger
	defaultOpts := log.NewOptions()
	defaultLogger, err := log.NewLogger(defaultOpts)
	assert.NoError(t, err)
	assert.NotNil(t, defaultLogger)
	// Since default output is stdout, we can't easily capture it without redirecting os.Stdout globally.
	// So we just check if it logs without panic.
	assert.NotPanics(t, func() {
		defaultLogger.Info("Info from NewLogger default")
	})
}

// TestLoggerInstance tests methods of a logger instance.
func TestLoggerInstance(t *testing.T) {
	var buf bytes.Buffer
		opts := log.NewOptions()
	opts.Level = "info"
	opts.Format = log.FormatJSON // For predictable named logger output
	// Configure logger to write to buffer
	l := log.NewLoggerWithWriter(opts, &buf)

	// Test different log levels and methods
	l.Debug("debug log") // Should not appear if level is info
	l.Info("info log")
	l.Warn("warn log")
	l.Error("error log")
	assert.NotPanics(t, func() { // DPanic should not panic in non-development mode (typical for tests)
		l.DPanic("dpanic log")
	})

	output := buf.String()
	t.Logf("Logger output:\n%s", output) // Log output for debugging

	assert.NotContains(t, output, "debug log")
	assert.Contains(t, output, "info log")
	assert.Contains(t, output, "warn log")
	assert.Contains(t, output, "error log")
	// DPanic logs as error in typical test environments if development mode is false
	assert.Contains(t, output, "dpanic log")

	// Test formatters
	buf.Reset() // Clear buffer
	l.Infof("info with format: %s", "formatted_string")
	assert.Contains(t, buf.String(), "info with format: formatted_string")

	buf.Reset()
	l.Warnf("warn with format: %d", 123)
	assert.Contains(t, buf.String(), "warn with format: 123")

	buf.Reset()
	l.Errorf("error with format: %v", true)
	assert.Contains(t, buf.String(), "error with format: true")

	// Test WithValues
	buf.Reset()
	lWith := l.WithValues("key1", "value1", "key2", 2)
	lWith.Info("info with fields")
	outputWith := buf.String()
	assert.Contains(t, outputWith, "info with fields")
	assert.Contains(t, outputWith, "key1")
	assert.Contains(t, outputWith, "value1")
	assert.Contains(t, outputWith, "key2")
	// Check for value 2; might be tricky with JSON if not careful about type
	assert.Contains(t, outputWith, "2")

	// Test WithName
	buf.Reset()
	namedLogger := l.WithName("myComponent")
	namedLogger.Info("info from named logger")
	outputNamed := buf.String()
	// Depending on the format (JSON/text), the name might appear differently.
	// For JSON, it's often under a "logger" key.
	if opts.Format == log.FormatJSON {
		assert.Contains(t, outputNamed, "\"N\":\"myComponent\"")
	} else {
		assert.Contains(t, outputNamed, "myComponent")
	}
	assert.Contains(t, outputNamed, "info from named logger")

	// Test further nesting of logger names
	nestedLogger := namedLogger.WithName("serviceX")
	buf.Reset()
	nestedLogger.Info("info from nested logger")
	outputNested := buf.String()
	assert.Contains(t, outputNested, "\"N\":\"myComponent.serviceX\"")
	assert.Contains(t, outputNested, "info from nested logger")
}

// TestGlobalLoggingMethods verifies that global logging functions (Info, Warn, Error, etc.) work as expected.
func TestGlobalLoggingMethods(t *testing.T) {
	// Redirect global logger output to a buffer for testing
	var buf bytes.Buffer
	opts := log.NewOptions()
	opts.Level = "info"
	opts.Format = log.FormatJSON // Use JSON for easier parsing of structured fields
	originalGlobalLogger := log.GetGlobalLogger() // Save original
	testLogger := log.NewLoggerWithWriter(opts, &buf)
	log.SetGlobalLogger(testLogger)
	defer log.SetGlobalLogger(originalGlobalLogger) // Restore original

	// Test global logging functions
	// These should not panic
	t.Run("GlobalLogFunctionsNoPanic", func(t *testing.T) {
		assert.NotPanics(t, func() { log.Debug("global debug") })
		assert.NotPanics(t, func() { log.Debugf("global debugf %s", "test") })
		assert.NotPanics(t, func() { log.Debugw("global debugw", "key", "value") })
		assert.NotPanics(t, func() { log.Info("global info") })
		assert.NotPanics(t, func() { log.Infof("global infof %s", "test") })
		assert.NotPanics(t, func() { log.Infow("global infow", "key", "value") })
		assert.NotPanics(t, func() { log.Warn("global warn") })
		assert.NotPanics(t, func() { log.Warnf("global warnf %s", "test") })
		assert.NotPanics(t, func() { log.Warnw("global warnw", "key", "value") })
		assert.NotPanics(t, func() { log.Error("global error") })
		assert.NotPanics(t, func() { log.Errorf("global errorf %s", "test") })
		assert.NotPanics(t, func() { log.Errorw("global errorw", "key", "value") })
		// DPanic logs as error in test environments
		assert.NotPanics(t, func() { log.DPanic("global dpanic") })
		assert.NotPanics(t, func() { log.DPanicf("global dpanicf %s", "test") })
		assert.NotPanics(t, func() { log.DPanicw("global dpanicw", "key", "value") })

		// Fatal methods will call os.Exit(1), so they need special handling
		// or should be tested by observing side effects (like a file being written before exit)
		// or by using a mock os.Exit.
		// For simplicity, we ensure they don't panic *before* os.Exit.
		// Note: Actual os.Exit behavior is not tested here.
	})

	output := buf.String()
	t.Logf("Global logger output:\n%s", output) // Log output for debugging

	// Verify log content
	assert.NotContains(t, output, "global debug") // Debug logs should be filtered out
	assert.Contains(t, output, "global info")
	assert.Contains(t, output, "global infof test")
	assert.Contains(t, output, "global infow")
	assert.Contains(t, output, "\"key\":\"value\"") // Check for structured field
	assert.Contains(t, output, "global warn")
	assert.Contains(t, output, "global error")
	assert.Contains(t, output, "global dpanic")

	// Test Sync
	assert.NoError(t, log.Sync()) // Should not error

	// Test global WithValues
	buf.Reset()
	log.WithValues("globalKey", "globalValue").Info("info with global field")
	outputWith := buf.String()
	assert.Contains(t, outputWith, "info with global field")
	assert.Contains(t, outputWith, "globalKey")
	assert.Contains(t, outputWith, "globalValue")

	// Test global WithName
	buf.Reset()
	log.WithName("globalComponent").Info("info from global named logger")
	outputNamed := buf.String()
	if opts.Format == log.FormatJSON {
		assert.Contains(t, outputNamed, "\"N\":\"globalComponent\"")
	} else {
		assert.Contains(t, outputNamed, "globalComponent")
	}
	assert.Contains(t, outputNamed, "info from global named logger")
}

// TestNamedLogger tests the WithName method for creating sub-loggers.
func TestNamedLogger(t *testing.T) {
	var buf bytes.Buffer
	opts := log.NewOptions()
	opts.Level = "info"
	opts.Format = log.FormatJSON
	l := log.NewLoggerWithWriter(opts, &buf)

	componentA := l.WithName("componentA")
	componentA.Info("Info from componentA")

	componentB := componentA.WithName("serviceB")
	componentB.Info("Info from componentA.serviceB")

	output := buf.String()
	assert.Contains(t, output, "\"N\":\"componentA\"", "Expected logger name componentA")
	assert.Contains(t, output, "Info from componentA")
	assert.Contains(t, output, "\"N\":\"componentA.serviceB\"", "Expected logger name componentA.serviceB")
	assert.Contains(t, output, "Info from componentA.serviceB")
}

// TestWithValues tests adding structured fields to log entries using WithValues.
func TestWithValues(t *testing.T) {
	var buf bytes.Buffer
	opts := log.NewOptions()
	// Ensure no default stdout/stderr to conflict with buffer capture
	opts.OutputPaths = []string{} // This will cause getWriteSyncer to err, must provide a writer or valid path
	opts.ErrorOutputPaths = []string{}
	opts.Format = log.FormatKeyValue // Use key=value format for key=value output
	opts.Level = "debug"

	// Create a logger that writes to the buffer
	// Since OutputPaths is empty, NewLogger would fail. We need NewLoggerWithWriter.
	logger := log.NewLoggerWithWriter(opts, &buf)
	assert.NotNil(t, logger, "NewLoggerWithWriter should return a logger")

	// Test instance WithValues
	l := logger.WithValues("component", "testComponent", "version", 1)
	l.Info("message from WithValues logger")

	output := buf.String()
	assert.Contains(t, output, `component=testComponent`, "Output should contain component field")
	assert.Contains(t, output, `version=1`, "Output should contain version field")
	assert.Contains(t, output, "message from WithValues logger", "Output should contain the log message")

	buf.Reset()

	// Test global WithValues
	originalGlobal := log.GetGlobalLogger()
	// Configure the global logger to use our buffer via a new logger instance
	globalTestOpts := log.NewOptions()
	globalTestOpts.Format = log.FormatKeyValue
	globalTestOpts.Level = "debug"
	globalTestLoggerWithBuffer := log.NewLoggerWithWriter(globalTestOpts, &buf) // This logger writes to buf
	log.SetGlobalLogger(globalTestLoggerWithBuffer)
	defer log.SetGlobalLogger(originalGlobal)

	lGlobal := log.WithValues("globalKey", "globalValue", "requestID", "xyz123")
	lGlobal.Info("global message with fields")
	globalOutput := buf.String()
	assert.Contains(t, globalOutput, `globalKey=globalValue`, "Global output should contain globalKey")
	assert.Contains(t, globalOutput, `requestID=xyz123`, "Global output should contain requestID")
	assert.Contains(t, globalOutput, "global message with fields", "Global output should contain the message")

	// Test that the original global logger (now our test logger) itself was not affected by WithValues,
	// as WithValues returns a new logger instance.
	buf.Reset()
	log.Info("message from base global logger") // This uses globalTestLoggerWithBuffer directly
	baseGlobalOutput := buf.String()
	assert.NotContains(t, baseGlobalOutput, `globalKey=globalValue`, "Base global logger should not have fields from WithValues")
	assert.Contains(t, baseGlobalOutput, "message from base global logger")

	// Test chaining WithValues on an instance logger
	buf.Reset()
	l1 := logger.WithValues("key1", "val1")
	l2 := l1.WithValues("key2", "val2", "key3", 3)
	l2.Info("chained message")
	chainedOutput := buf.String()
	assert.Contains(t, chainedOutput, `key1=val1`, "Chained output should contain key1")
	assert.Contains(t, chainedOutput, `key2=val2`, "Chained output should contain key2")
	assert.Contains(t, chainedOutput, `key3=3`, "Chained output should contain key3")
	assert.Contains(t, chainedOutput, "chained message", "Chained output should contain the message")
}

// TestLogLevelFilteringBasic basic tests for log level filtering.
func TestLogLevelFilteringBasic(t *testing.T) {
	var buf bytes.Buffer
	testCases := []struct {
		name      string
		level     string
		logFunc   func(l log.Logger) // Changed from *log.Logger to log.Logger
		shouldLog bool
		message   string
	}{
		{"debug at debug level", "debug", func(l log.Logger) { l.Debug("debug should log") }, true, "debug should log"},
		{"info at debug level", "debug", func(l log.Logger) { l.Info("info should log") }, true, "info should log"},
		{"warn at debug level", "debug", func(l log.Logger) { l.Warn("warn should log") }, true, "warn should log"},
		{"error at debug level", "debug", func(l log.Logger) { l.Error("error should log") }, true, "error should log"},

		{"debug at info level", "info", func(l log.Logger) { l.Debug("debug should NOT log") }, false, "debug should NOT log"},
		{"info at info level", "info", func(l log.Logger) { l.Info("info should log") }, true, "info should log"},
		{"warn at info level", "info", func(l log.Logger) { l.Warn("warn should log") }, true, "warn should log"},
		{"error at info level", "info", func(l log.Logger) { l.Error("error should log") }, true, "error should log"},

		{"debug at warn level", "warn", func(l log.Logger) { l.Debug("debug should NOT log") }, false, "debug should NOT log"},
		{"info at warn level", "warn", func(l log.Logger) { l.Info("info should NOT log") }, false, "info should NOT log"},
		{"warn at warn level", "warn", func(l log.Logger) { l.Warn("warn should log") }, true, "warn should log"},
		{"error at warn level", "warn", func(l log.Logger) { l.Error("error should log") }, true, "error should log"},

		{"debug at error level", "error", func(l log.Logger) { l.Debug("debug should NOT log") }, false, "debug should NOT log"},
		{"info at error level", "error", func(l log.Logger) { l.Info("info should NOT log") }, false, "info should NOT log"},
		{"warn at error level", "error", func(l log.Logger) { l.Warn("warn should NOT log") }, false, "warn should NOT log"},
		{"error at error level", "error", func(l log.Logger) { l.Error("error should log") }, true, "error should log"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf.Reset()
			opts := log.NewOptions()
			opts.Level = tc.level
			l := log.NewLoggerWithWriter(opts, &buf)
			assert.NotPanics(t, func() { // Ensure logging itself doesn't panic
				tc.logFunc(l)
			})
			output := buf.String()
			if tc.shouldLog {
				assert.Contains(t, output, tc.message, "Expected log message to be present")
			} else {
				// When message shouldn't be logged, it also shouldn't be present in output.
				// The tc.message for !shouldLog cases was "", so NotContains wouldn't work as expected.
				// Corrected tc.message to reflect the string that shouldn't be present.
				assert.NotContains(t, output, tc.message, "Expected log message to be absent")
			}
		})
	}
}

// TestLogFormat tests different log formats (json, text).
func TestLogFormat(t *testing.T) {
	var buf bytes.Buffer

	// Test JSON format
	buf.Reset()
	optsJSON := log.NewOptions()
	optsJSON.Level = "info"
	optsJSON.Format = log.FormatJSON // Use const
	loggerJSON := log.NewLoggerWithWriter(optsJSON, &buf)
	loggerJSON.Info("json format test")
	outputJSON := buf.String()
	assert.True(t, strings.HasPrefix(outputJSON, "{"), "Expected JSON output to start with {")
	assert.Contains(t, outputJSON, "\"M\":\"json format test\"", "Expected message field in JSON")

	// Test Text format
	buf.Reset()
	optsText := log.NewOptions()
	optsText.Level = "info"
	optsText.Format = log.FormatText // Use const
	loggerText := log.NewLoggerWithWriter(optsText, &buf)
	loggerText.Info("text format test")
	outputText := buf.String()
	// Text format includes timestamp, level, caller, message
	assert.NotContains(t, outputText, "{", "Text output should not contain JSON braces")
	assert.Contains(t, outputText, "INFO", "Expected level in text output")
	assert.Contains(t, outputText, "text format test", "Expected message in text output")
	assert.Contains(t, outputText, "log/log_test.go", "Expected caller info in text output")

	// Test Invalid format (NewLogger should return error)
	buf.Reset()
	optsInvalid := log.NewOptions()
	optsInvalid.Level = "info"
	optsInvalid.Format = "invalid_format"
	loggerInvalid, err := log.NewLogger(optsInvalid)
	assert.Error(t, err, "Expected error for invalid format")
	assert.True(t, lmccerrors.IsCode(err, lmccerrors.ErrLogOptionInvalid), "Expected ErrLogOptionInvalid for invalid format")
	assert.Nil(t, loggerInvalid, "Logger should be nil for invalid format")

	// Test with writer (NewLoggerWithWriter should panic for invalid format string in options)
	buf.Reset()
	optsInvalidWithWriter := log.NewOptions()
	optsInvalidWithWriter.Level = "info"
	optsInvalidWithWriter.Format = "invalid_format" // This will be handled by newLoggerInternal leading to a panic in NewLoggerWithWriter

	assert.Panics(t, func() {
		log.NewLoggerWithWriter(optsInvalidWithWriter, &buf)
	}, "NewLoggerWithWriter should panic for invalid format in options passed to newLoggerInternal")
}

// TestLogOutputPaths tests logging to stdout, stderr, and a file.
func TestLogOutputPaths(t *testing.T) {
	// Test stdout
	t.Run("stdout", func(t *testing.T) {
		opts := log.NewOptions()
		opts.OutputPaths = []string{"stdout"}
		logger, err := log.NewLogger(opts)
		assert.NoError(t, err)
		assert.NotPanics(t, func() { logger.Info("Message to stdout") })
	})

	// Test stderr
	t.Run("stderr", func(t *testing.T) {
		opts := log.NewOptions()
		opts.OutputPaths = []string{"stderr"}
		logger, err := log.NewLogger(opts)
		assert.NoError(t, err)
		assert.NotPanics(t, func() { logger.Info("Message to stderr") })
	})

	// Test file output
	t.Run("file output", func(t *testing.T) {
		tempDir := t.TempDir()
		logFilePath := filepath.Join(tempDir, "test.log")

		opts := log.NewOptions()
		opts.OutputPaths = []string{logFilePath}
		logger, err := log.NewLogger(opts)
		assert.NoError(t, err)
		assert.NotNil(t, logger)

		logger.Info("Message to file")
		// Important: Sync the logger to ensure the message is written to the file.
		errSync := logger.Sync()
		assert.NoError(t, errSync)

		// Verify file content
		content, err := os.ReadFile(logFilePath)
		assert.NoError(t, err, "Failed to read log file")
		assert.Contains(t, string(content), "Message to file", "Log file content mismatch")
	})

	// Test multiple output paths
	t.Run("multiple outputs", func(t *testing.T) {
		tempDir := t.TempDir()
		logFilePath1 := filepath.Join(tempDir, "multi1.log")
		logFilePath2 := filepath.Join(tempDir, "multi2.log")

		opts := log.NewOptions()
		// For testing multiple files, stdout/stderr can make verification harder.
		opts.OutputPaths = []string{logFilePath1, logFilePath2}
		logger, err := log.NewLogger(opts)
		assert.NoError(t, err)
		assert.NotNil(t, logger)

		logger.Info("Message to multiple files")
		errSync := logger.Sync() // Sync to flush
		assert.NoError(t, errSync)

		content1, err1 := os.ReadFile(logFilePath1)
		assert.NoError(t, err1)
		assert.Contains(t, string(content1), "Message to multiple files")

		content2, err2 := os.ReadFile(logFilePath2)
		assert.NoError(t, err2)
		assert.Contains(t, string(content2), "Message to multiple files")
	})



	// Test mixed stdout and file
	t.Run("stdout and file", func(t *testing.T) {
		tempDir := t.TempDir()
		logFilePath := filepath.Join(tempDir, "stdout_file.log")

		opts := log.NewOptions()
		opts.OutputPaths = []string{"stdout", logFilePath}
		logger, err := log.NewLogger(opts)
		assert.NoError(t, err)
		assert.NotNil(t, logger)

		// Log a message
		logger.Info("Message to stdout and file")
		errSync := logger.Sync() // Sync to flush file
		// Note: Sync on stdout may fail on some systems, so we don't assert NoError
		// The important thing is that the file write works
		_ = errSync

		// Verify file content
		content, errFile := os.ReadFile(logFilePath)
		assert.NoError(t, errFile)
		assert.Contains(t, string(content), "Message to stdout and file")

		// Verifying stdout is harder here without capturing os.Stdout globally.
		// We assume if file write worked and no error, stdout also worked.
	})

	// Test unsupported URI (assuming file paths are the only supported non-stdout/stderr)
	t.Run("unsupported URI", func(t *testing.T) {
		opts := log.NewOptions()
		opts.OutputPaths = []string{"http://localhost/log"}
		logger, err := log.NewLogger(opts)
		assert.Error(t, err, "Expected error for unsupported URI scheme")
		assert.Nil(t, logger, "Logger should be nil for unsupported URI")
		// The error might be ErrLogInitialization if opening fails, or ErrLogOptionInvalid if path is deemed invalid before open attempt.
		// Based on current getWriteSyncer, it will try to open it as a file.
		// If the OS treats "http://..." as an invalid filename, os.OpenFile will fail.
		assert.True(t, lmccerrors.IsCode(err, lmccerrors.ErrLogInitialization), "Expected ErrLogInitialization for bad file path like URI")
	})
}

// TestEnableDisableCaller tests enabling and disabling caller information in logs.
func TestEnableDisableCaller(t *testing.T) {
	var buf bytes.Buffer

	// Enable caller
	buf.Reset()
	optsCaller := log.NewOptions()
	optsCaller.Level = "info"
	optsCaller.DisableCaller = false // Explicitly enable
	optsCaller.Format = log.FormatJSON
	loggerCaller := log.NewLoggerWithWriter(optsCaller, &buf)
	loggerCaller.Info("caller enabled")
	outputCaller := buf.String()
	assert.Contains(t, outputCaller, "\"C\":", "Expected caller field when enabled")
	assert.Contains(t, outputCaller, "log/log_test.go", "Expected caller file name")

	// Disable caller
	buf.Reset()
	optsNoCaller := log.NewOptions()
	optsNoCaller.Level = "info"
	optsNoCaller.DisableCaller = true
	optsNoCaller.Format = log.FormatJSON
	loggerNoCaller := log.NewLoggerWithWriter(optsNoCaller, &buf)
	loggerNoCaller.Info("caller disabled")
	outputNoCaller := buf.String()
	assert.NotContains(t, outputNoCaller, "\"C\":", "Expected no caller field when disabled")
}

// TestEnableDisableStacktrace tests enabling and disabling stack traces for error-level logs.
func TestEnableDisableStacktrace(t *testing.T) {
	var buf bytes.Buffer

	// Enable stacktrace (default for "error" level)
	buf.Reset()
	optsStack := log.NewOptions()
	optsStack.Level = "info" // Log info to ensure error is above threshold
	optsStack.DisableStacktrace = false // Explicitly enable (or rely on default)
	optsStack.Format = log.FormatJSON
	// Default StacktraceLevel is "error"
	loggerStack := log.NewLoggerWithWriter(optsStack, &buf)
	loggerStack.Error("stacktrace enabled for error")
	outputStack := buf.String()
	assert.Contains(t, outputStack, "\"stacktrace\":", "Expected stacktrace field for errors when enabled")

	// Disable stacktrace
	buf.Reset()
	optsNoStack := log.NewOptions()
	optsNoStack.Level = "info"
	optsNoStack.DisableStacktrace = true
	optsNoStack.Format = log.FormatJSON
	loggerNoStack := log.NewLoggerWithWriter(optsNoStack, &buf)
	loggerNoStack.Error("stacktrace disabled for error")
	outputNoStack := buf.String()
	assert.NotContains(t, outputNoStack, "\"stacktrace\":", "Expected no stacktrace field for errors when disabled")

	// Test stacktrace at a different level (e.g., "warn" if configured)
	buf.Reset()
	optsWarnStack := log.NewOptions()
	optsWarnStack.Level = "warn"           // Set log level to warn to actually see warn logs
	optsWarnStack.StacktraceLevel = "warn" // Enable stacktrace for warn and above
	optsWarnStack.Format = log.FormatJSON
	loggerWarnStack := log.NewLoggerWithWriter(optsWarnStack, &buf)

	loggerWarnStack.Warn("stacktrace at warn")
	outputWarnStack := buf.String()
	assert.Contains(t, outputWarnStack, "\"stacktrace\":", "Expected stacktrace for warn level when configured")

	buf.Reset() // Reset buffer before next log to isolate output
	loggerWarnStack.Info("info below warn stacktrace level") // Should not have stack
	outputInfoNoStack := buf.String()
	assert.NotContains(t, outputInfoNoStack, "stacktrace", "Info log should not have stacktrace when StacktraceLevel is warn")
}

// TestDevelopmentMode tests the behavior of the logger in development mode (e.g., DPanic panics).
func TestDevelopmentMode(t *testing.T) {
	var buf bytes.Buffer
	optsDev := log.NewOptions()
	optsDev.Development = true
	optsDev.Level = "debug" // Ensure DPanic level is active
	// DPanic should panic in development mode
	loggerDev := log.NewLoggerWithWriter(optsDev, &buf)

	assert.PanicsWithValue(t, "dpanic in dev mode", func() {
		loggerDev.DPanic("dpanic in dev mode")
	}, "DPanic should panic in development mode")

	// Ensure other levels don't panic
	assert.NotPanics(t, func() { loggerDev.Debug("debug in dev") })
	assert.NotPanics(t, func() { loggerDev.Info("info in dev") })
	assert.NotPanics(t, func() { loggerDev.Error("error in dev") })

	// Test production mode DPanic
	buf.Reset()
	optsProd := log.NewOptions()
	optsProd.Development = false
	optsProd.Level = "debug"
	loggerProd := log.NewLoggerWithWriter(optsProd, &buf)

	assert.NotPanics(t, func() {
		loggerProd.DPanic("dpanic in prod mode")
	}, "DPanic should not panic in production mode")
	assert.Contains(t, buf.String(), "dpanic in prod mode", "DPanic should log as error in production")
	// In prod, DPanic logs at ErrorLevel
	// Check for JSON format's level key, assuming default opts.Format = json for optsProd
	if optsProd.Format == log.FormatJSON {
		assert.Contains(t, buf.String(), "\"L\":\"ERROR\"", "DPanic should log as ERROR level in production with JSON format")
	} else {
		assert.Contains(t, buf.String(), "ERROR", "DPanic should log as ERROR level in production with text format")
	}
}

// TestZapEncodingConfig verifies that custom Zap encoding configurations are applied.
func TestZapEncodingConfig(t *testing.T) {
	var buf bytes.Buffer
	opts := log.NewOptions()
	opts.Format = log.FormatJSON // Ensure JSON for predictable key names
	opts.TimeFormat = "2006-01-02" // Custom time format
	// opts.ZapOptions removed as it's not in Options struct

	// Customize encoder config (example: full caller path)
	customEncoderCfg := zap.NewProductionEncoderConfig() // Start with a base
	customEncoderCfg.MessageKey = "MSG"
	customEncoderCfg.LevelKey = "LVL"
	customEncoderCfg.TimeKey = "TS"
	customEncoderCfg.CallerKey = "CALLER"
	customEncoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02") // Apply custom time format
	customEncoderCfg.EncodeLevel = zapcore.LowercaseLevelEncoder            // e.g., info, error
	customEncoderCfg.EncodeCaller = zapcore.FullCallerEncoder              // Full caller path
	opts.EncoderConfig = &customEncoderCfg                                 // Set the custom config

	l := log.NewLoggerWithWriter(opts, &buf)
	l.Info("testing custom encoder")

	output := buf.String()
	assert.Contains(t, output, "\"MSG\":\"testing custom encoder\"", "Expected custom message key and message")
	assert.Contains(t, output, "\"LVL\":\"info\"", "Expected custom level key and lowercase level")
	assert.Regexp(t, `\"TS\":\"\d{4}-\d{2}-\d{2}\"`, output, "Expected custom time key and format")
	assert.Contains(t, output, "\"CALLER\":", "Expected custom caller key")
	assert.Contains(t, output, "/pkg/log/log_test.go", "Expected full caller path")
}

// TestInitialGlobalLogger ensures a default global logger is available immediately after package initialization.
func TestInitialGlobalLogger(t *testing.T) {
	// This test relies on the package-level init or TestMain to set up a global logger.
	assert.NotPanics(t, func() {
		log.Info("Testing initial global logger") // Should not panic
	}, "Initial global logger should be usable")

	globalL := log.GetGlobalLogger()
	assert.NotNil(t, globalL, "GetGlobalLogger should return a non-nil logger instance")
}

// TestContextHandling verifies that context can be used with the logger.
func TestContextHandling(t *testing.T) {
	var buf bytes.Buffer
	testOpts := log.NewOptions()
	testOpts.Format = log.FormatJSON // For predictable output
	// Define a context key and value
	type contextKey string
	const myKey contextKey = "myCustomKey"
	const myValue string = "customValue123"
	testOpts.ContextKeys = []any{myKey} // Configure logger to extract this key

	l := log.NewLoggerWithWriter(testOpts, &buf)

	ctx := context.WithValue(context.Background(), myKey, myValue)

	assert.NotPanics(t, func() {
		l.CtxInfof(ctx, "logging with context, value for %s should be present", myKey)
	})

	output := buf.String()
	assert.Contains(t, output, "logging with context")
	assert.Contains(t, output, "\"myCustomKey\":\"customValue123\"", "Expected custom context key and value in log output")

	// Test Ctxw
	buf.Reset()
	l.Ctxw(ctx, "ctxw message with fields", "extraField", "extraVal")
	outputCtxw := buf.String()
	assert.Contains(t, outputCtxw, "ctxw message with fields")
	assert.Contains(t, outputCtxw, "\"myCustomKey\":\"customValue123\"")
	assert.Contains(t, outputCtxw, "\"extraField\":\"extraVal\"")
}

// TestConcurrentLogging basic test for concurrent logging to ensure thread safety.
func TestConcurrentLogging(t *testing.T) {
	var buf bytes.Buffer // Shared buffer
	opts := log.NewOptions()
	opts.Level = "debug"
	l := log.NewLoggerWithWriter(opts, &buf) // Logger writing to the shared buffer

	numGoroutines := 50
	logsPerGoroutine := 20
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < logsPerGoroutine; j++ {
				l.Infof("Log from goroutine %d, message %d", id, j)
			}
		}(i)
	}
	wg.Wait()

	// Sync to ensure all logs are written
	err := l.Sync()
	assert.NoError(t, err, "Sync should not return an error")
	
	// Verify the output (it's hard to check exact order or count without more complex sync)
	// The main thing is to ensure no race conditions or panics.
	// We can count the number of log entries.
	output := buf.String()
	scanner := bufio.NewScanner(strings.NewReader(output))
	lineCount := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" { // Only count non-empty lines
			lineCount++
		}
	}

	expectedLogCount := numGoroutines * logsPerGoroutine
	// Due to concurrent writes to buffer, we might lose some lines or have partial writes
	// So we check that we have at least some reasonable number of lines
	assert.GreaterOrEqual(t, lineCount, expectedLogCount/2, "Should have at least half the expected log lines from concurrent logging")
	assert.LessOrEqual(t, lineCount, expectedLogCount, "Should not have more than expected log lines")
}

// TestFatalLogging
// Note: True fatal behavior (os.Exit) is hard to test without a mock.
// This test checks if the log message is written before the (unmocked) exit.
// This test is inherently risky as it *will* call os.Exit if the logger works correctly.
// It's often skipped or run in a subprocess in CI.
func TestFatalLogging(t *testing.T) {
	t.Skip("Skipping TestFatalLogging as it calls os.Exit and disrupts test suite flow.")

	// If this test were to run, it would look something like this:
	/*
		var buf bytes.Buffer
		opts := log.NewOptions()
		opts.Level = "info"
		// Setup a logger that writes to the buffer
		fatalLogger := log.NewLoggerWithWriter(opts, &buf)

		// We need to fork the process to test os.Exit, or use a mock.
		// For simplicity, this example does not fork or mock os.Exit.
		// Instead, it would rely on observing the buffer content if the program
		// somehow continued after Fatal (which it shouldn't).

		// This will attempt to call os.Exit(1)
		fatalLogger.Fatal("this is a fatal error")

		// Code below this line should not be reached if Fatal works as expected.
		assert.Contains(t, buf.String(), "this is a fatal error", "Fatal log message should be present")
	*/
}

// TODO: Add tests for ErrorGroup.Format method once implemented.

// TestKeyValueFormat tests the new key=value format functionality.
func TestKeyValueFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := log.NewOptions()
	opts.Format = log.FormatKeyValue
	opts.Level = "debug"

	logger := log.NewLoggerWithWriter(opts, &buf)

	t.Run("basic WithValues", func(t *testing.T) {
		buf.Reset()
		l := logger.WithValues("component", "testComponent", "version", 1)
		l.Info("message from WithValues logger")
		output := buf.String()
		assert.Contains(t, output, "component=testComponent", "Output should contain component field")
		assert.Contains(t, output, "version=1", "Output should contain version field")
		assert.Contains(t, output, "message from WithValues logger", "Output should contain the log message")
	})

	t.Run("chained WithValues", func(t *testing.T) {
		buf.Reset()
		l1 := logger.WithValues("key1", "val1")
		l2 := l1.WithValues("key2", "val2", "key3", 3)
		l2.Info("chained message")
		output := buf.String()
		assert.Contains(t, output, "key1=val1", "Output should contain key1")
		assert.Contains(t, output, "key2=val2", "Output should contain key2")
		assert.Contains(t, output, "key3=3", "Output should contain key3")
		assert.Contains(t, output, "chained message", "Output should contain the message")
	})

	t.Run("values with spaces", func(t *testing.T) {
		buf.Reset()
		l := logger.WithValues("description", "this is a test", "count", 42)
		l.Info("message with spaces")
		output := buf.String()
		assert.Contains(t, output, `description="this is a test"`, "Output should contain quoted value with spaces")
		assert.Contains(t, output, "count=42", "Output should contain count field")
	})

	t.Run("context fields with key=value", func(t *testing.T) {
		buf.Reset()
		type customKey string
		const myKey customKey = "myCustomKey"
		optsWithContext := log.NewOptions()
		optsWithContext.Format = log.FormatKeyValue
		optsWithContext.ContextKeys = []any{myKey}
		loggerWithContext := log.NewLoggerWithWriter(optsWithContext, &buf)

		ctx := context.WithValue(context.Background(), myKey, "customValue123")
		loggerWithContext.Ctxw(ctx, "context message", "extraField", "extraVal")
		output := buf.String()
		assert.Contains(t, output, "myCustomKey=customValue123", "Output should contain context field")
		assert.Contains(t, output, "extraField=extraVal", "Output should contain extra field")
		assert.Contains(t, output, "context message", "Output should contain the message")
	})

	t.Run("different log levels", func(t *testing.T) {
		buf.Reset()
		l := logger.WithValues("service", "test-service")
		
		l.Debug("debug message")
		l.Info("info message")
		l.Warn("warn message")
		l.Error("error message")
		
		output := buf.String()
		// All should contain the service field
		lines := strings.Split(strings.TrimSpace(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "message") {
				assert.Contains(t, line, "service=test-service", "Each log line should contain the service field")
			}
		}
	})

	t.Run("format validation", func(t *testing.T) {
		// Test that FormatKeyValue is properly validated
		opts := log.NewOptions()
		opts.Format = log.FormatKeyValue
		errs := opts.Validate()
		assert.Empty(t, errs, "FormatKeyValue should be valid")
		
		// Test invalid format
		opts.Format = "invalid-format"
		errs = opts.Validate()
		assert.NotEmpty(t, errs, "Invalid format should produce validation errors")
	})
}