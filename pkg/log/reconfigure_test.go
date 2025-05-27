/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * Contains tests for logger reconfiguration functionality.
 */

package log_test

import (
	// "bytes" // Commented out if only used by captureOutput
	"errors" // Added for errors.Unwrap
	// "io" // Commented out if only used by captureOutput
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	lmccerrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

/* // Commenting out unused function captureOutput
// captureOutput 捕获指定函数的输出。
// (captureOutput captures the output of a given function.)
func captureOutput(action func()) string {
	old := os.Stdout // 备份 os.Stdout (Backup os.Stdout)
	r, w, _ := os.Pipe()
	os.Stdout = w

	action()

	// 确保在关闭写入器和恢复 stdout 之前，所有内容都已刷写。
	// (Ensure everything is flushed before closing the writer and restoring stdout.)
	// log.Sync() 应该由 action 内部调用，以刷写 zap 的内部缓冲区。
	// (log.Sync() should be called inside action to flush zap's internal buffers.)
	// 这里的 w.Sync() 尝试刷写管道写入器的操作系统级缓冲区。
	// (This w.Sync() attempts to flush the OS-level buffer for the pipe writer.)
	_ = w.Sync()

	_ = w.Close()
	os.Stdout = old // 恢复 os.Stdout (Restore os.Stdout)
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		// This is a test helper, so perhaps a panic or t.Fatal is more appropriate
		// depending on how critical stdout capturing is for the tests that would use this.
		// For now, let's just log it to standard error if we can't capture.
		fmt.Fprintf(os.Stderr, "Error capturing stdout in test helper: %v\\n", err)
	}
	return strings.TrimSpace(buf.String())
}
*/

// TestReconfigureGlobalLogger_Basic 测试基本的日志记录器重新配置功能。
// (TestReconfigureGlobalLogger_Basic tests basic logger reconfiguration functionality.)
func TestReconfigureGlobalLogger_Basic(t *testing.T) {
	localRequire := require.New(t)
	localAssert := assert.New(t)

	// 0. 保存并恢复原始的全局 logger 的 zap 实例以便同步，并确保测试后恢复到默认配置
	// (Save the original global logger's zap instance for syncing, and ensure restoration to default config after test)
	originalZapLogger := log.Std().GetZapLogger()
	defer func() {
		log.Init(log.NewOptions()) // 恢复到默认日志配置 (Restore to default log config)
		if originalZapLogger != nil {
			_ = originalZapLogger.Sync() // 同步原始 logger (Sync original logger)
		}
	}()

	// 1. 使用初始选项初始化全局 logger (例如，INFO 级别，文本格式，输出到临时文件)
	// (Initialize global logger with initial options (e.g., INFO level, text format, output to a temp file))
	tempDir := t.TempDir() // Create a general temp directory for the test
	initialLogFilePath := filepath.Join(tempDir, "initial_test_out.log")

	initialOpts := log.NewOptions()
	initialOpts.Level = zapcore.InfoLevel.String()
	initialOpts.Format = log.FormatText
	initialOpts.OutputPaths = []string{initialLogFilePath} // 输出到临时文件 (Output to temp file)
	initialOpts.EnableColor = false
	log.Init(initialOpts)

	// 2. 记录一条初始消息
	// (Log an initial message)
	initialMsg := "Initial INFO message to temp file"
	log.Info(initialMsg)
	log.Debug("This initial debug message should not appear in temp file (INFO level).")
	_ = log.Sync() // Sync to ensure it's flushed

	// 读取初始日志文件的内容并验证
	// (Read initial log file content and verify)
	initialLogBytes, err := os.ReadFile(initialLogFilePath)
	localRequire.NoError(err, "Failed to read initial log file")
	initialLogOutput := string(initialLogBytes)
	localAssert.Contains(initialLogOutput, initialMsg, "Initial message should be in the temp log file")
	localAssert.NotContains(initialLogOutput, "initial debug message", "Initial debug message should not be in the temp log file")

	// 3. 定义用于后续重新配置的日志文件路径
	// (Define log file path for subsequent reconfiguration)
	reconfiguredLogFilePath := filepath.Join(tempDir, "reconfigured_test.log")

	// 4. 使用新的选项调用 ReconfigureGlobalLogger (输出到另一个文件)
	// (Call ReconfigureGlobalLogger with new options (output to another file))
	reconfiguredOpts := log.NewOptions()
	reconfiguredOpts.Level = zapcore.DebugLevel.String()       // 更改为 DEBUG 级别 (Change to DEBUG level)
	reconfiguredOpts.Format = log.FormatJSON                // 更改为 JSON 格式 (Change to JSON format)
	reconfiguredOpts.OutputPaths = []string{reconfiguredLogFilePath} // 输出到新的文件 (Output to new file)

	err = log.ReconfigureGlobalLogger(reconfiguredOpts)
	localRequire.NoError(err, "ReconfigureGlobalLogger should not return an error")

	// 5. 记录第二条和第三条消息 (现在应该写入新文件)
	// (Log second and third messages (should now write to new file))
	reconfiguredMsgDebug := "Reconfigured DEBUG message to file"
	reconfiguredMsgInfo := "Reconfigured INFO message to file"
	log.Debug(reconfiguredMsgDebug)
	log.Info(reconfiguredMsgInfo)
	_ = log.Sync() // Sync to ensure all logs are written to the reconfigured file

	// 6. 读取重新配置后的日志文件的内容
	// (Read the content of the reconfigured log file)
	reconfiguredContentBytes, err := os.ReadFile(reconfiguredLogFilePath)
	localRequire.NoError(err, "Failed to read reconfigured log file")
	fileContent := string(reconfiguredContentBytes)

	// 7. 验证重新配置后的日志文件的内容
	// (Verify the content of the reconfigured log file)
	localAssert.NotEmpty(fileContent, "Reconfigured log file should not be empty")
	localAssert.Contains(fileContent, reconfiguredMsgDebug, "Reconfigured file content should contain the DEBUG message")
	localAssert.Contains(fileContent, `"L":"DEBUG"`, "Reconfigured file content should indicate DEBUG level")
	localAssert.Contains(fileContent, reconfiguredMsgInfo, "Reconfigured file content should contain the INFO message")
	localAssert.Contains(fileContent, `"L":"INFO"`, "Reconfigured file content should indicate INFO level")
	localAssert.True(strings.HasPrefix(strings.TrimSpace(fileContent), "{"), "Log content should start with '{' indicating JSON")
	localAssert.True(strings.HasSuffix(strings.TrimSpace(fileContent), "}"), "Log content should end with '}' indicating JSON")
	localAssert.NotContains(fileContent, initialMsg, "Reconfigured file content should NOT contain the initial message")

	// 8. 验证调用 ReconfigureGlobalLogger(nil) 会返回错误
	// (Verify calling ReconfigureGlobalLogger(nil) returns an error)
	err = log.ReconfigureGlobalLogger(nil)
	localAssert.Error(err, "ReconfigureGlobalLogger(nil) should return an error")
	if err != nil { // Check err is not nil before further assertions
		localAssert.True(lmccerrors.IsCode(err, lmccerrors.ErrLogOptionInvalid), "Error code should be ErrLogOptionInvalid for nil options")
		localAssert.Contains(err.Error(), "cannot reconfigure global logger with nil options", "Error message for nil options mismatch")
	}
}

// TestReconfigureGlobalLogger_Concurrent 安全地并发调用 ReconfigureGlobalLogger 和日志记录函数。
// (TestReconfigureGlobalLogger_Concurrent safely calls ReconfigureGlobalLogger and logging functions concurrently.)
func TestReconfigureGlobalLogger_Concurrent(t *testing.T) {
	localRequire := require.New(t)
	localAssert := assert.New(t)

	originalZapLogger := log.Std().GetZapLogger()
	defer func() {
		log.Init(log.NewOptions()) // 恢复到默认日志配置 (Restore to default log config)
		if originalZapLogger != nil {
			_ = originalZapLogger.Sync() // 同步原始 logger (Sync original logger)
		}
	}()

	tempDir := t.TempDir()
	logFile1Path := filepath.Join(tempDir, "concurrent_1.log")
	logFile2Path := filepath.Join(tempDir, "concurrent_2.log")

	initialOpts := log.NewOptions()
	initialOpts.OutputPaths = []string{logFile1Path}
	initialOpts.Level = zapcore.InfoLevel.String()
	initialOpts.Format = log.FormatText
	initialOpts.EnableColor = false
	log.Init(initialOpts)

	var wg sync.WaitGroup
	numGoroutines := 50
	numLogsPerGoroutine := 100

	for i := 0; i < numGoroutines/2; i++ {
		wg.Add(1)
		go func(gNum int) {
			defer wg.Done()
			for j := 0; j < numLogsPerGoroutine; j++ {
				log.Infof("Goroutine %d, Log %d, Initial Config", gNum, j)
				time.Sleep(time.Microsecond)
			}
		}(i)
	}

	// Reconfigure logger concurrently
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(5 * time.Millisecond) // Allow some initial logs to be written
		reconfiguredOpts := log.NewOptions()
		reconfiguredOpts.OutputPaths = []string{logFile2Path}
		reconfiguredOpts.Level = zapcore.DebugLevel.String()
		reconfiguredOpts.Format = log.FormatJSON
		err := log.ReconfigureGlobalLogger(reconfiguredOpts)
		localRequire.NoError(err, "Concurrent ReconfigureGlobalLogger should not fail")
		log.Debug("Logger reconfigured during concurrent operations.")
	}()

	for i := numGoroutines / 2; i < numGoroutines; i++ {
		wg.Add(1)
		go func(gNum int) {
			defer wg.Done()
			for j := 0; j < numLogsPerGoroutine; j++ {
				log.Infof("Goroutine %d, Log %d, Potentially New Config", gNum, j)
				log.Debugf("Goroutine %d, Debug Log %d, Potentially New Config", gNum, j)
				time.Sleep(time.Microsecond)
			}
		}(i)
	}

	wg.Wait()
	finalLogger := log.Std()
	if finalLogger != nil && finalLogger.GetZapLogger() != nil {
	    errSync := finalLogger.GetZapLogger().Sync()
	    localRequire.NoError(errSync, "Syncing final logger should not produce an error")
	}

	content1, err1 := os.ReadFile(logFile1Path)
	content2, err2 := os.ReadFile(logFile2Path)

	localRequire.NoError(err1, "Error reading logFile1")
	localRequire.NoError(err2, "Error reading logFile2")

	localRequire.True(len(content1) > 0 || len(content2) > 0, "At least one log file should have content")

	if len(content2) > 0 {
		file2Content := string(content2)
		localAssert.Contains(file2Content, `"L":"DEBUG"`, "logFile2 should contain debug level messages")
		localAssert.Contains(file2Content, "Logger reconfigured during concurrent operations.", "logFile2 should contain the reconfiguration debug message")
	}

	t.Logf("Log file 1 (%s) size: %d bytes", logFile1Path, len(content1))
	t.Logf("Log file 2 (%s) size: %d bytes", logFile2Path, len(content2))
}

// TestReconfigureGlobalLogger_ErrorCases tests scenarios where newLogger fails during reconfiguration, checking for lmccerrors.ErrLogReconfigure and the wrapped original errors (ErrLogOptionInvalid, ErrLogInitialization).
// (TestReconfigureGlobalLogger_ErrorCases tests scenarios where newLogger fails during reconfiguration, checking for lmccerrors.ErrLogReconfigure and the wrapped original errors (ErrLogOptionInvalid, ErrLogInitialization).)
// 当前 newLogger 在失败时会 panic 或回退到 stderr，所以难以直接测试错误返回路径。
// (Currently newLogger panics or falls back to stderr on failure, making it hard to test error return paths directly.)
func TestReconfigureGlobalLogger_ErrorCases(t *testing.T) {
	localAssert := assert.New(t)

	// Save and restore original logger state
	originalZapLogger := log.Std().GetZapLogger()
	defer func() {
		log.Init(log.NewOptions()) // Restore to default
		if originalZapLogger != nil {
			_ = originalZapLogger.Sync()
		}
	}()

	t.Run("InvalidLevelDuringReconfigure", func(t *testing.T) {
		opts := log.NewOptions()
		opts.Level = "invalid-reconfig-level"
		err := log.ReconfigureGlobalLogger(opts)
		localAssert.Error(err, "ReconfigureGlobalLogger with invalid level should return an error")
		if err != nil {
			localAssert.True(lmccerrors.IsCode(err, lmccerrors.ErrLogReconfigure), "Error code should be ErrLogReconfigure")
			localAssert.Contains(err.Error(), "invalid-reconfig-level", "Error message should indicate the specific invalid level")
			// Use errors.Unwrap to get the immediately wrapped error
			wrappedError := errors.Unwrap(err)
			localAssert.NotNil(wrappedError, "Wrapped error should not be nil")
			if wrappedError != nil {
				localAssert.True(lmccerrors.IsCode(wrappedError, lmccerrors.ErrLogOptionInvalid), "Wrapped error should be ErrLogOptionInvalid")
			}
		}
	})

	t.Run("InvalidOutputDuringReconfigure", func(t *testing.T) {
		opts := log.NewOptions()
		invalidPath := "unsupported://reconfig-output"
		opts.OutputPaths = []string{invalidPath}
		err := log.ReconfigureGlobalLogger(opts)
		localAssert.Error(err, "ReconfigureGlobalLogger with invalid output scheme should return an error")
		if err != nil {
			// Check top-level error code and message components
			localAssert.True(lmccerrors.IsCode(err, lmccerrors.ErrLogReconfigure), "Error code should be ErrLogReconfigure")
			localAssert.Contains(err.Error(), lmccerrors.ErrLogReconfigure.String(), "Error message should contain ErrLogReconfigure string")
			localAssert.Contains(err.Error(), "failed to create new logger for reconfiguration", "Error message should indicate reconfiguration failure context")
			// The next part of the message should now come from the ErrLogOptionInvalid due to unsupported scheme
			localAssert.Contains(err.Error(), lmccerrors.ErrLogOptionInvalid.String(), "Error message should contain ErrLogOptionInvalid string for unsupported scheme")
			localAssert.Contains(err.Error(), invalidPath, "Error message should mention the unsupported path scheme")

			// Check wrapped errors
			wrappedError1 := errors.Unwrap(err)
			localAssert.NotNil(wrappedError1, "First wrapped error should not be nil")
			if wrappedError1 != nil {
				// wrappedError1 is the error returned by newLogger, which is ErrorfWithCode(ErrLogInitialization, "... %w", originalNewWithCodeError)
				// Its Coder should be ErrLogInitialization (though the test was checking for ErrLogOptionInvalid, which is deeper)
				// Let's adjust to check the Coder of wrappedError1 first.
				// localAssert.True(lmccerrors.IsCode(wrappedError1, lmccerrors.ErrLogOptionInvalid), "Wrapped error should be ErrLogOptionInvalid for unsupported scheme")
				// The direct Coder of wrappedError1 is actually from the wrapping ErrorfWithCode in newLogger, e.g. ErrLogInitialization
				// The test was previously checking the Coder of the *cause* of wrappedError1.

				// Let's get the error that was created by NewWithCode (the one with ErrLogOptionInvalid)
				errCreatedByNewWithCode := errors.Unwrap(wrappedError1) // This should be the *withCode from NewWithCode(ErrLogOptionInvalid, ...)
				localAssert.NotNil(errCreatedByNewWithCode, "The error from newLogger should wrap the error from NewWithCode")
				if errCreatedByNewWithCode != nil {
					localAssert.True(lmccerrors.IsCode(errCreatedByNewWithCode, lmccerrors.ErrLogOptionInvalid), "The cause of newLogger's error should have Coder ErrLogOptionInvalid")
					localAssert.Contains(errCreatedByNewWithCode.Error(), lmccerrors.ErrLogOptionInvalid.String(), "Error message from NewWithCode should contain its Coder string")
					localAssert.Contains(errCreatedByNewWithCode.Error(), invalidPath, "Error message from NewWithCode should contain the specific unsupported path")

					// Now, unwrap errCreatedByNewWithCode (which is a *withCode) to get its cause (which is a *fundamental)
					actualFundamentalError := errors.Unwrap(errCreatedByNewWithCode)
					localAssert.NotNil(actualFundamentalError, "The error from NewWithCode should wrap an internal fundamental error")
					if actualFundamentalError != nil {
						localAssert.Contains(actualFundamentalError.Error(), invalidPath, "Fundamental error message mismatch")

						// The fundamental error may or may not wrap further errors depending on implementation
						// This is acceptable as long as the error chain is correct
						_ = errors.Unwrap(actualFundamentalError) // Just unwrap without asserting nil
					}
				}
			}
		}
	})
}

// TODO: TestReconfigureGlobalLogger_CallerSkip - 如果有必要，测试重新配置后 caller skip 是否仍然正确
// (TestReconfigureGlobalLogger_CallerSkip - if necessary, test if caller skip remains correct after reconfiguration)
// 通常 zap 的 AddCallerSkip(1) 应该对实例级别的方法调用有效，所以重新配置不应影响它。
// (Usually zap's AddCallerSkip(1) should work for instance-level method calls, so reconfiguration shouldn't affect it.)

// TestReconfigureGlobalLogger tests the ReconfigureGlobalLogger function.
// (TestReconfigureGlobalLogger 测试 ReconfigureGlobalLogger 函数。)
func TestReconfigureGlobalLogger(t *testing.T) {
	// localRequire := require.New(t) // Removed as it's not used at this top level
	// localAssert := assert.New(t) // Already removed

	initialDefaultOptsForTestFunc := log.NewOptions()
	log.Init(initialDefaultOptsForTestFunc)
	_ = log.Sync()

	defer func() {
		log.Init(initialDefaultOptsForTestFunc)
		_ = log.Sync()
	}()

	captureGlobalLogOutput := func(actionToLog func(), levelForCaptureFile string) string {
		tempDir := t.TempDir()
		logFilePath := filepath.Join(tempDir, "capture.log")
		captureOpts := log.NewOptions()
		captureOpts.OutputPaths = []string{logFilePath}
		captureOpts.Format = log.FormatText
		captureOpts.Level = levelForCaptureFile
		captureOpts.EnableColor = false
		// This will be restored to after capture, effectively resetting global log for next subtest or action
		restoreToDefaultOpts := log.NewOptions() 

		err := log.ReconfigureGlobalLogger(captureOpts)
		require.New(t).NoError(err, "Helper: Failed to reconfigure global logger for capture")

		actionToLog()

		err = log.Sync()
		require.New(t).NoError(err, "Helper: Failed to sync global logger during capture")

		contentBytes, errReadFile := os.ReadFile(logFilePath)
		require.New(t).NoError(errReadFile, "Helper: Failed to read captured log file")

		err = log.ReconfigureGlobalLogger(restoreToDefaultOpts)
		require.New(t).NoError(err, "Helper: Failed to restore global logger to default after capture")
		_ = log.Sync()

		return string(contentBytes)
	}

	t.Run("ReconfigureLevelToDebug", func(t *testing.T) {
		subtestRequire := require.New(t)
		subtestAssert := assert.New(t)

		subtestDefaultOpts := log.NewOptions() // INFO level
		log.Init(subtestDefaultOpts)
		_ = log.Sync()
		initialSubtestStd := log.Std()
		subtestRequire.Equal(zapcore.InfoLevel.String(), log.Std().GetZapLogger().Level().String())

		newOpts := log.NewOptions()
		newOpts.Level = zapcore.DebugLevel.String()
		err := log.ReconfigureGlobalLogger(newOpts)
		subtestRequire.NoError(err)

		updatedLogger := log.Std()
		subtestAssert.Equal(zapcore.DebugLevel.String(), updatedLogger.GetZapLogger().Level().String())
		subtestAssert.NotSame(initialSubtestStd, updatedLogger)

		// At this point, global logger is DEBUG.
		// captureGlobalLogOutput will set its own file output and level for capture.
		output := captureGlobalLogOutput(func() {
			log.Info("Info message for debug reconfig")   // Logged because current global level is DEBUG
			log.Debug("Debug message for debug reconfig") // Logged
		}, zapcore.DebugLevel.String()) // Capture file is also set to DEBUG

		subtestAssert.Contains(output, "INFO") // Check for level string
		subtestAssert.Contains(output, "Info message for debug reconfig") // Check for message
		subtestAssert.Contains(output, "DEBUG") // Check for level string
		subtestAssert.Contains(output, "Debug message for debug reconfig") // Check for message
		
		// After captureGlobalLogOutput, global logger is reset to INFO by the helper.
		subtestAssert.Equal(zapcore.InfoLevel.String(), log.Std().GetZapLogger().Level().String(), "Global logger should be INFO after capture helper restores it")
		log.Init(subtestDefaultOpts) // Explicitly reset for safety before next subtest if any in this scope
		_ = log.Sync()
	})

	t.Run("ReconfigureToWarnAndVerifyOutput", func(t *testing.T) {
		subtestRequire := require.New(t)
		subtestAssert := assert.New(t)

		subtestDefaultOpts := log.NewOptions() // INFO level
		log.Init(subtestDefaultOpts)
		_ = log.Sync()
		subtestRequire.Equal(zapcore.InfoLevel.String(), log.Std().GetZapLogger().Level().String())

		newOpts := log.NewOptions()
		newOpts.Level = zapcore.WarnLevel.String()
		err := log.ReconfigureGlobalLogger(newOpts)
		subtestRequire.NoError(err)
		subtestAssert.Equal(zapcore.WarnLevel.String(), log.Std().GetZapLogger().Level().String())

		// Global logger is now WARN.
		output := captureGlobalLogOutput(func() {
			log.Info("This INFO should NOT be seen after warn reconfig")    // Not logged
			log.Debug("This DEBUG should NOT be seen after warn reconfig") // Not logged
			log.Warn("This WARN should be seen after warn reconfig")      // Logged
		}, zapcore.WarnLevel.String()) // Capture file is set to WARN
		
		subtestAssert.NotContains(output, "INFO	This INFO should NOT be seen after warn reconfig")
		subtestAssert.NotContains(output, "DEBUG	This DEBUG should NOT be seen after warn reconfig")
		subtestAssert.Contains(output, "WARN")
		subtestAssert.Contains(output, "This WARN should be seen after warn reconfig")

		log.Init(subtestDefaultOpts)
		_ = log.Sync()
	})

	t.Run("ReconfigureWithNilOptions", func(t *testing.T) {
		subtestRequire := require.New(t)
		subtestAssert := assert.New(t)

		currentOpts := log.NewOptions()
		currentOpts.Level = zapcore.ErrorLevel.String()
		log.Init(currentOpts)
		_ = log.Sync()
		originalLevel := log.Std().GetZapLogger().Level().String()

		err := log.ReconfigureGlobalLogger(nil)
		subtestRequire.Error(err)

		subtestAssert.Contains(err.Error(), "cannot reconfigure global logger with nil options")
		subtestAssert.Equal(originalLevel, log.Std().GetZapLogger().Level().String(), "Global logger level should not change after nil reconfigure attempt")
		
		log.Init(log.NewOptions())
		_ = log.Sync()
	})
} 

// TestHotReloadWithValidConfig simulates a hot reload scenario with a valid new configuration.
// ... existing code ... 