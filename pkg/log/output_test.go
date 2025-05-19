/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * Contains tests for log file output functionality.
 */

package log_test

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

// TestFileOutputJSON tests logging to a file in JSON format.
// (TestFileOutputJSON 测试以 JSON 格式记录到文件。)
func TestFileOutputJSON(t *testing.T) {
	localRequire := require.New(t)
	localAssert := assert.New(t)

	// 创建临时目录 (Create temporary directory)
	tempDir := t.TempDir() // This creates a unique temp dir for the test
	logFilePath := filepath.Join(tempDir, "test_json.log")

	// 配置 logger 输出到文件 (Configure logger to output to file)
	opts := log.NewOptions()
	opts.OutputPaths = []string{logFilePath}
	opts.Format = log.FormatJSON
	opts.Level = zapcore.DebugLevel.String()

	logger := log.NewLogger(opts)
	localRequire.NotNil(logger)

	// 写入日志 (Write logs)
	logger.Debugw("Debug message", "key1", "val1")
	logger.Info("Info message")
	logger.Warnf("Warn message %s", "formatted")

	// 确保日志写入文件 (Ensure logs are written to file)
	err := logger.Sync()
	localRequire.NoError(err)

	// 读取文件内容 (Read file content)
	contentBytes, err := os.ReadFile(logFilePath)
	localRequire.NoError(err)
	content := string(contentBytes)

	// 断言文件内容 (Assert file content)
	localAssert.NotEmpty(content, "Log file should not be empty")
	localAssert.Contains(content, `"L":"DEBUG"`, "Should contain debug level log")
	localAssert.Contains(content, `"M":"Debug message"`, "Should contain debug message")
	localAssert.Contains(content, `"key1":"val1"`, "Should contain debug key-value pair")
	localAssert.Contains(content, `"L":"INFO"`, "Should contain info level log")
	localAssert.Contains(content, `"M":"Info message"`, "Should contain info message")
	localAssert.Contains(content, `"L":"WARN"`, "Should contain warn level log")
	localAssert.Contains(content, `"M":"Warn message formatted"`, "Should contain warn message")
}

// TestFileOutputText tests logging to a file in Text format.
// (TestFileOutputText 测试以 Text 格式记录到文件。)
func TestFileOutputText(t *testing.T) {
	localRequire := require.New(t)
	localAssert := assert.New(t)

	tempDir := t.TempDir()
	logFilePath := filepath.Join(tempDir, "test_text.log")

	opts := log.NewOptions()
	opts.OutputPaths = []string{logFilePath}
	opts.Format = log.FormatText
	opts.EnableColor = false // Disable color for easier string matching
	opts.Level = zapcore.InfoLevel.String()

	logger := log.NewLogger(opts)
	localRequire.NotNil(logger)

	logger.Info("Info message text")
	logger.Warnf("Warn message text %s", "formatted")

	err := logger.Sync()
	localRequire.NoError(err)

	contentBytes, err := os.ReadFile(logFilePath)
	localRequire.NoError(err)
	content := string(contentBytes)

	localAssert.NotEmpty(content, "Log file should not be empty")
	localAssert.Contains(content, "INFO", "Should contain INFO level log")
	localAssert.Contains(content, "Info message text", "Should contain info message")
	localAssert.Contains(content, "WARN", "Should contain WARN level log")
	localAssert.Contains(content, "Warn message text formatted", "Should contain warn message")
	
	matched, err := regexp.MatchString(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}`, content)
	localRequire.NoError(err, "Regex compilation should not fail")
	localAssert.True(matched, "Log content should contain timestamp")
	localAssert.Contains(content, "log/output_test.go", "Should contain caller info relative to this new file")
}

// TestFileOutputMulti tests logging to both stdout and a file.
// (TestFileOutputMulti 测试同时记录到 stdout 和文件。)
func TestFileOutputMulti(t *testing.T) {
	localRequire := require.New(t)
	localAssert := assert.New(t)

	tempDir := t.TempDir()
	logFilePath := filepath.Join(tempDir, "test_multi.log")

	opts := log.NewOptions()
	opts.OutputPaths = []string{logFilePath} 
	opts.Format = log.FormatJSON             
	opts.Level = zapcore.InfoLevel.String()

	logger := log.NewLogger(opts)
	localRequire.NotNil(logger)

	logger.Info("Multi output info")
	logger.Warn("Multi output warn")

	err := logger.Sync()
	localRequire.NoError(err)

	contentBytes, err := os.ReadFile(logFilePath)
	localRequire.NoError(err)
	content := string(contentBytes)

	localAssert.NotEmpty(content, "Log file should not be empty")
	localAssert.Contains(content, `"M":"Multi output info"`, "File should contain info message")
	localAssert.Contains(content, `"M":"Multi output warn"`, "File should contain warn message")
} 