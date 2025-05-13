/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * Contains core logger tests and shared test helper functions for the log package.
 */

package log_test

import (
	// "os" // Potentially needed if helper functions for file ops are kept here
	// "path/filepath" // Potentially needed for file op helpers
	"sync"
	"testing"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
	"github.com/stretchr/testify/assert"

	// "github.com/stretchr/testify/require" // require might not be needed if only assert is used here
	"go.uber.org/zap/zapcore"
)

// --- Shared Test Helper Types/Functions (if any) ---
// Example: (If createTempLogFile was decided to be kept shared here)
/*
func createTempLogFile(t *testing.T, initialContent string) (string, func()) {
	t.Helper()
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "shared_test_log_file.log")
	if initialContent != "" {
		err := os.WriteFile(tmpFile, []byte(initialContent), 0644)
		require.New(t).NoError(err, "Failed to write initial temp log file content")
	}
	cleanup := func() {}
	return tmpFile, cleanup
}
*/

// --- Core Test Cases ---

func TestNewLogger(t *testing.T) {
	localAssert := assert.New(t)
	opts := log.NewOptions()
	localAssert.NotPanics(func() {
		logger := log.NewLogger(opts)
		localAssert.NotNil(logger)
		localAssert.NotNil(logger.GetZapLogger())
		logger.Info("Info from NewLogger default")
		logger.Debug("Debug from NewLogger default")
		_ = logger.Sync()
	})

	opts.Level = zapcore.DebugLevel.String()
	opts.Format = log.FormatText
	opts.Name = "TestLoggerInstance"
	localAssert.NotPanics(func() {
		logger := log.NewLogger(opts)
		localAssert.NotNil(logger)
		zapLogger := logger.GetZapLogger()
		localAssert.NotNil(zapLogger)
		logger.Info("Info from NewLogger specific")
		logger.Debug("Debug from NewLogger specific")
		_ = logger.Sync()
	})
}

func TestStd(t *testing.T) {
	localAssert := assert.New(t)
	stdLogger := log.Std()
	localAssert.NotNil(stdLogger, "Std() should return a non-nil logger")
	var wg sync.WaitGroup
	const numGoroutines = 10
	firstInstance := log.Std()
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			localAssert.Same(firstInstance, log.Std(), "Std() should return the same instance concurrently")
		}()
	}
	wg.Wait()
	localAssert.NotPanics(func() {
		log.Std().Info("Info via Std()")
		_ = log.Std().Sync()
	})
}

func TestInit(t *testing.T) {
	localAssert := assert.New(t)
	localAssert.NotPanics(func() {
		opts := log.NewOptions()
		opts.Level = zapcore.ErrorLevel.String()
		opts.Name = "TestInitLogger"
		log.Init(opts)
	})
	defer func() {
		log.Init(log.NewOptions())
	}()
}

func TestGlobalLoggingMethods(t *testing.T) {
	localAssert := assert.New(t)
	defer func() { _ = log.Sync() }()
	localAssert.NotPanics(func() {
		log.Debug("global debug")
		log.Debugf("global debugf %s", "test")
		log.Debugw("global debugw", "key", "value")
		log.Info("global info")
		log.Infof("global infof %s", "test")
		log.Infow("global infow", "key", "value")
		log.Warn("global warn")
		log.Warnf("global warnf %s", "test")
		log.Warnw("global warnw", "key", "value")
		log.Error("global error")
		log.Errorf("global errorf %s", "test")
		log.Errorw("global errorw", "key", "value")
	})
}

func TestWithName(t *testing.T) {
	localAssert := assert.New(t)
	baseLogger := log.NewLogger(log.NewOptions())
	defer func() { _ = baseLogger.Sync() }()
	localAssert.NotPanics(func() {
		logger := baseLogger.WithName("componentA")
		localAssert.NotNil(logger)
		logger.Info("Info from componentA")
		nestedLogger := logger.WithName("serviceB")
		localAssert.NotNil(nestedLogger)
		localAssert.NotSame(logger, nestedLogger)
		nestedLogger.Info("Info from componentA.serviceB")
	})
}

func TestWithValues(t *testing.T) {
	localAssert := assert.New(t)
	baseLogger := log.NewLogger(log.NewOptions())
	defer func() { _ = baseLogger.Sync() }()
	localAssert.NotPanics(func() {
		logger := baseLogger.WithValues("key1", "value1")
		localAssert.NotNil(logger)
		logger.Infow("Info with key1", "extraKey", "extraValue")
		nestedLogger := logger.WithValues("key2", "value2")
		localAssert.NotNil(nestedLogger)
		localAssert.NotSame(logger, nestedLogger)
		nestedLogger.Infow("Info with key1 and key2", "anotherKey", "anotherValue")
	})
}

func TestLogLevelFilteringBasic(t *testing.T) {
	localAssert := assert.New(t)
	localAssert.NotPanics(func() {
		opts := log.NewOptions()
		opts.Level = zapcore.DebugLevel.String()
		logger := log.NewLogger(opts)
		logger.Debug("debug should log")
		logger.Info("info should log")
		logger.Warn("warn should log")
		logger.Error("error should log")
		_ = logger.Sync()
	})
	localAssert.NotPanics(func() {
		opts := log.NewOptions()
		opts.Level = zapcore.WarnLevel.String()
		logger := log.NewLogger(opts)
		logger.Debug("debug should NOT log, but call shouldn't panic")
		logger.Info("info should NOT log, but call shouldn't panic")
		logger.Warn("warn should log")
		logger.Error("error should log")
		_ = logger.Sync()
	})
}

func TestLogFormatBasic(t *testing.T) {
	localAssert := assert.New(t)
	localAssert.NotPanics(func() {
		opts := log.NewOptions()
		opts.Format = log.FormatJSON
		logger := log.NewLogger(opts)
		logger.Info("json format test")
		_ = logger.Sync()
	})
	localAssert.NotPanics(func() {
		opts := log.NewOptions()
		opts.Format = log.FormatText
		logger := log.NewLogger(opts)
		logger.Info("text format test")
		_ = logger.Sync()
	})
}

// TestNewOptions and TestOptions_Validate are moved to options_test.go

// Any shared helper functions that were in the original log_test.go and are used by
// the tests remaining in this file, or by tests in *other* _test.go files in this package,
// should be kept or consolidated here.
// For example, if there was a common setup function or a more complex log capture helper.
// The simple createTempLogFile was moved to output_test.go as it was specific to file output tests.
// If more tests need it, it should be moved back here or to a shared test_utils.go.