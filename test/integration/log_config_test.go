/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * Contains integration tests for log and config packages.
 */

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 测试配置模板 (Test configuration template)
const configTemplate = `
log:
  level: "%s"
  format: "%s"
  outputPaths:
    - "%s"
  errorOutputPaths:
    - "%s"
  enableColor: %v
  # Add other log options if they become relevant for hot-reload tests
  # development: false 
  # name: ""
  # skipCaller: false
  # disableStacktrace: false 
`

// 为测试创建配置文件 (Create a configuration file for testing)
// Ensures outputPaths and errorOutputPaths always use the provided file paths.
// (确保 outputPaths 和 errorOutputPaths 始终使用提供的文件路径。)
func setupConfigFile(t *testing.T, level, format, outputPath, errorOutputPath string, enableColor bool) string {
	// 创建临时目录 (Create temporary directory)
	tempDir, err := os.MkdirTemp("", "config_log_integration_test")
	require.NoError(t, err, "Failed to create temp dir")

	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	configPath := filepath.Join(tempDir, "config.yaml")
	configContent := fmt.Sprintf(configTemplate,
		level, format, outputPath, errorOutputPath, enableColor)

	err = os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err, "Failed to write config file")

	return configPath
}

// 修改配置文件 (Modify configuration file)
// Ensures outputPaths and errorOutputPaths are preserved from the original file.
// (确保 outputPaths 和 errorOutputPaths 从原始文件中保留。)
func updateConfigFile(t *testing.T, configPath, level, format string, enableColor bool) {
	content, err := os.ReadFile(configPath)
	require.NoError(t, err, "Failed to read config file for update")

	v := viper.New()
	v.SetConfigType("yaml")
	err = v.ReadConfig(bytes.NewBuffer(content))
	require.NoError(t, err, "Failed to parse config for update")

	originalOutputPaths := v.GetStringSlice("log.outputPaths")
	originalErrorOutputPaths := v.GetStringSlice("log.errorOutputPaths")
	require.NotEmpty(t, originalOutputPaths, "log.outputPaths should not be empty in existing config")
	require.NotEmpty(t, originalErrorOutputPaths, "log.errorOutputPaths should not be empty in existing config")

	newConfigContent := fmt.Sprintf(configTemplate,
		level, format, originalOutputPaths[0], originalErrorOutputPaths[0], enableColor)

	time.Sleep(100 * time.Millisecond) // Shorter sleep before write, longer after

	err = os.WriteFile(configPath, []byte(newConfigContent), 0644)
	require.NoError(t, err, "Failed to update config file")
	time.Sleep(250 * time.Millisecond) // Increased delay after write for FS events
}

// Helper function to initialize logger with specific file paths for testing
func initTestLogger(t *testing.T, level, format, outputPath, errorOutputPath string, enableColor bool) *log.Options {
	opts := log.NewOptions()
	opts.Level = level
	opts.Format = format
	opts.OutputPaths = []string{outputPath}
	opts.ErrorOutputPaths = []string{errorOutputPath} // Explicitly set error paths
	opts.EnableColor = enableColor
	// Set other Options fields to match configTemplate if they are relevant

	errs := opts.Validate()
	require.Empty(t, errs, "Validation of initial log options failed")

	log.Init(opts) // Initialize the global logger (does not return error)
	// No error check needed here as Init likely handles errors internally (e.g., logs or panics)
	return opts
}

// 测试日志级别动态更新 (Test dynamic log level update)
func TestLogLevelDynamicUpdate(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "log_level_update_test")
	require.NoError(t, err, "Failed to create temp dir for logs")
	t.Cleanup(func() { os.RemoveAll(tempDir) })

	logFilePath := filepath.Join(tempDir, "test.log")
	errorLogPath := filepath.Join(tempDir, "error.log") // Dedicated error log file

	// 1. 初始化日志记录器 (Initialize the logger)
	initialLogLevel := "info"
	initialLogFormat := "json"
	_ = initTestLogger(t, initialLogLevel, initialLogFormat, logFilePath, errorLogPath, false)
	
	// 2. 创建并加载初始配置文件 (Create and load initial config file)
	configPath := setupConfigFile(t, initialLogLevel, initialLogFormat, logFilePath, errorLogPath, false)

	var appCfg config.Config // Use an empty config.Config struct
	cfgMgr, err := config.LoadConfigAndWatch(&appCfg,
		config.WithConfigFile(configPath, "yaml"),
		config.WithHotReload(true),
	)
	require.NoError(t, err, "Failed to load config and watch")

	// 3. 注册日志配置热重载 (Register log config hot-reload)
	log.RegisterConfigHotReload(cfgMgr)

	// 短暂等待，以确保 watcher 启动并可能已处理了初始配置（如果适用）
	// (Short wait to ensure watcher is started and potentially processed initial config if applicable)
	time.Sleep(250 * time.Millisecond)

	assert.Equal(t, initialLogLevel, log.Std().GetZapLogger().Level().String(), "Initial log level should be '%s'", initialLogLevel)

	log.Debug("This is a debug message that should NOT appear")
	log.Info("This is an info message that should appear")
	err = log.Std().GetZapLogger().Sync()
	require.NoError(t, err, "Failed to sync logger after initial writes")

	// 更新配置为 debug 级别 (Update config to debug level)
	updatedLogLevel := "debug"
	updateConfigFile(t, configPath, updatedLogLevel, initialLogFormat, false)

	time.Sleep(500 * time.Millisecond) // Wait for update

	assert.Equal(t, updatedLogLevel, log.Std().GetZapLogger().Level().String(), "Log level should be updated to '%s'", updatedLogLevel)

	log.Debug("This is a debug message that should now appear")
	log.Info("This is another info message")
	err = log.Std().GetZapLogger().Sync()
	require.NoError(t, err, "Failed to sync logger after debug writes")

	logContent, err := os.ReadFile(logFilePath)
	require.NoError(t, err, "Failed to read log file")
	logContentStr := string(logContent)

	assert.NotContains(t, logContentStr, "debug message that should NOT appear")
	assert.Contains(t, logContentStr, "info message that should appear")
	assert.Contains(t, logContentStr, "debug message that should now appear")
	assert.Contains(t, logContentStr, "another info message")
}

// 测试日志格式动态更新 (Test dynamic log format update)
func TestLogFormatDynamicUpdate(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "log_format_update_test")
	require.NoError(t, err, "Failed to create temp dir for logs")
	t.Cleanup(func() { os.RemoveAll(tempDir) })

	logFilePath := filepath.Join(tempDir, "format_test.log")
	errorLogPath := filepath.Join(tempDir, "format_error.log")

	initialLogLevel := "info"
	initialLogFormat := "json"
	_ = initTestLogger(t, initialLogLevel, initialLogFormat, logFilePath, errorLogPath, false)

	configPath := setupConfigFile(t, initialLogLevel, initialLogFormat, logFilePath, errorLogPath, false)
	
	var appCfg config.Config
	cfgMgr, err := config.LoadConfigAndWatch(&appCfg,
		config.WithConfigFile(configPath, "yaml"),
		config.WithHotReload(true),
	)
	require.NoError(t, err, "Failed to load config and watch")

	log.RegisterConfigHotReload(cfgMgr)
	time.Sleep(250 * time.Millisecond)

	log.Info("This is a message in JSON format")
	err = log.Std().GetZapLogger().Sync()
	require.NoError(t, err, "Failed to sync logger for JSON format")

	logContent, err := os.ReadFile(logFilePath)
	require.NoError(t, err, "Failed to read log file for JSON check")
	logContentStr := string(logContent)
	assert.True(t, isValidJSON(t, logContentStr), "Log should be in JSON format. Content:\n%s", logContentStr)

	updatedLogFormat := "text"
	updateConfigFile(t, configPath, initialLogLevel, updatedLogFormat, false)
	time.Sleep(500 * time.Millisecond)

	err = os.WriteFile(logFilePath, []byte{}, 0644) // Clear log for new format
	require.NoError(t, err, "Failed to clear log file for text format test")

	log.Info("This is a message in text format")
	err = log.Std().GetZapLogger().Sync()
	require.NoError(t, err, "Failed to sync logger for text format")

	logContent, err = os.ReadFile(logFilePath)
	require.NoError(t, err, "Failed to read log file for text check")
	logContentStr = string(logContent)
	assert.False(t, isValidJSON(t, logContentStr), "Log should be in text format after update. Content:\n%s", logContentStr)
	assert.Contains(t, logContentStr, "text format")
	assert.Contains(t, logContentStr, "INFO")
}

// 判断字符串是否为有效的 JSON (Determine if a string is valid JSON)
func isValidJSON(t *testing.T, s string) bool {
	s = strings.TrimSpace(s)
	if s == "" { // Handle empty string case after trim
		return false
	}
	if strings.Contains(s, "\n") {
		lines := strings.Split(s, "\n")
		 nonEmptyLines := 0
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			nonEmptyLines++
			var js map[string]interface{}
			if json.Unmarshal([]byte(line), &js) != nil {
				t.Logf("Invalid JSON line: [%s]", line)
				return false
			}
		}
		return nonEmptyLines > 0 // True if there was at least one valid JSON line
	}
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

// 测试并发配置更新 (Test concurrent configuration updates)
func TestConcurrentConfigUpdates(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "concurrent_config_updates_test")
	require.NoError(t, err, "Failed to create temp dir for logs")
	t.Cleanup(func() { os.RemoveAll(tempDir) })

	logFilePath := filepath.Join(tempDir, "concurrent.log")
	errorLogPath := filepath.Join(tempDir, "concurrent_error.log")

	initialLogLevel := "warn"
	initialLogFormat := "json"
	_ = initTestLogger(t, initialLogLevel, initialLogFormat, logFilePath, errorLogPath, false)

	configPath := setupConfigFile(t, initialLogLevel, initialLogFormat, logFilePath, errorLogPath, false)

	var appCfg config.Config
	cfgMgr, err := config.LoadConfigAndWatch(&appCfg,
		config.WithConfigFile(configPath, "yaml"),
		config.WithHotReload(true),
	)
	require.NoError(t, err, "Failed to load config and watch")

	log.RegisterConfigHotReload(cfgMgr)
	time.Sleep(250 * time.Millisecond) 

	assert.Equal(t, initialLogLevel, log.Std().GetZapLogger().Level().String(), "Initial log level should be '%s'", initialLogLevel)

	levels := []string{"debug", "info", "warn", "error", "info", "debug"}
	formats := []string{"json", "text", "json", "text", "json", "text"}

	for i, level := range levels {
		updateConfigFile(t, configPath, level, formats[i], false)
		time.Sleep(500 * time.Millisecond)

		currentLoggerLevel := log.Std().GetZapLogger().Level().String()
		assert.Equal(t, level, currentLoggerLevel,
			fmt.Sprintf("Log level should be updated to %s in iteration %d (after update, before logging). Expected %s, got %s", level, i, level, currentLoggerLevel))

		log.Debug(fmt.Sprintf("Debug message in %s format (iter %d)", formats[i], i))
		log.Info(fmt.Sprintf("Info message in %s format (iter %d)", formats[i], i))
		log.Warn(fmt.Sprintf("Warn message in %s format (iter %d)", formats[i], i))
		log.Error(fmt.Sprintf("Error message in %s format (iter %d)", formats[i], i))
		
		err = log.Std().GetZapLogger().Sync()
		require.NoError(t, err, fmt.Sprintf("Failed to sync logger in iteration %d", i))
	}

	logContent, err := os.ReadFile(logFilePath)
	require.NoError(t, err, "Failed to read final log file")
	logContentStr := string(logContent)

	finalExpectedLevel := levels[len(levels)-1]
	finalExpectedFormat := formats[len(formats)-1]
	assert.Equal(t, finalExpectedLevel, log.Std().GetZapLogger().Level().String(), "Final logger level should be correctly set")

	if finalExpectedLevel == "debug" {
		assert.Contains(t, logContentStr, fmt.Sprintf("Debug message in %s format (iter %d)", finalExpectedFormat, len(levels)-1))
	}
	assert.Contains(t, logContentStr, fmt.Sprintf("Info message in %s format (iter %d)", finalExpectedFormat, len(levels)-1))
} 