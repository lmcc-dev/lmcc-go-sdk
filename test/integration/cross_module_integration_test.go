/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package integration

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
	lmccerrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCrossModuleIntegration 跨模块集成测试套件
// (TestCrossModuleIntegration cross-module integration test suite)

// ApplicationConfig 应用程序配置，包含所有模块的配置
// (ApplicationConfig application configuration containing configs for all modules)
type ApplicationConfig struct {
	config.Config                 // 嵌入基础配置 (Embed base config)
	App           AppSettings     `mapstructure:"app"`
	Features      AppFeatureFlags `mapstructure:"features"`
}

type AppSettings struct {
	Name        string `mapstructure:"name" default:"CrossModuleTestApp"`
	Environment string `mapstructure:"environment" default:"testing"`
	Debug       bool   `mapstructure:"debug" default:"false"`
}

type AppFeatureFlags struct {
	EnableDetailedLogging bool `mapstructure:"enableDetailedLogging" default:"true"`
	LogErrorDetails       bool `mapstructure:"logErrorDetails" default:"true"`
}

// mockApplicationService 模拟应用服务，演示三个模块的集成使用
// (mockApplicationService mock application service demonstrating integration of three modules)
type mockApplicationService struct {
	name   string
	logger log.Logger
}

func newMockApplicationService(name string) *mockApplicationService {
	return &mockApplicationService{
		name:   name,
		logger: log.Std(),
	}
}

func (s *mockApplicationService) ProcessBusinessLogic(ctx context.Context, userID string, action string) error {
	// 记录请求开始 (Log request start)
	s.logger.Ctx(ctx, "Starting business logic processing",
		"user_id", userID,
		"action", action,
		"service", s.name,
	)

	// 模拟业务逻辑错误 (Simulate business logic errors)
	if userID == "" {
		err := lmccerrors.NewWithCode(
			lmccerrors.ErrValidation,
			"user ID cannot be empty",
		)

		// 使用配置决定是否记录错误详情 (Use config to decide whether to log error details)
		// 简化配置检查，直接记录错误 (Simplify config check, directly log error)
		s.logger.Errorw("Validation failed in business logic",
			"user_id", userID,
			"action", action,
			"error_code", lmccerrors.GetCoder(err).Code(),
			"error", err,
		)
		return err
	}

	if action == "forbidden" {
		err := lmccerrors.WithCode(
			lmccerrors.Errorf("action '%s' is not allowed for user %s", action, userID),
			lmccerrors.ErrForbidden,
		)

		s.logger.Errorw("Forbidden action attempted",
			"user_id", userID,
			"action", action,
			"error", err,
		)
		return err
	}

	if action == "fail" {
		return s.simulateInternalError(ctx, userID, action)
	}

	// 成功场景 (Success scenario)
	s.logger.Info("Business logic completed successfully",
		"user_id", userID,
		"action", action,
		"service", s.name,
	)

	return nil
}

func (s *mockApplicationService) simulateInternalError(ctx context.Context, userID, action string) error {
	// 模拟数据库操作失败 (Simulate database operation failure)
	dbErr := lmccerrors.New("database connection failed")

	// 包装错误并添加上下文 (Wrap error and add context)
	wrappedErr := lmccerrors.WithCode(
		lmccerrors.Wrapf(dbErr, "failed to process action '%s' for user '%s'", action, userID),
		lmccerrors.ErrInternalServer,
	)

	// 记录内部错误 (Log internal error)
	s.logger.Errorw("Internal server error occurred",
		"user_id", userID,
		"action", action,
		"service", s.name,
		"error", wrappedErr,
		"error_stack", fmt.Sprintf("%+v", wrappedErr),
	)

	return wrappedErr
}

// TestConfigLogIntegration 测试配置和日志模块的集成
// (TestConfigLogIntegration tests integration between config and log modules)
func TestConfigLogIntegration(t *testing.T) {
	tempDir := setupIntegrationTestDir(t)
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "app.log")
	configFile := filepath.Join(tempDir, "app.yaml")

	// 创建初始配置 (Create initial configuration)
	initialConfig := fmt.Sprintf(`
server:
  host: "localhost"
  port: 8080
log:
  level: "info"
  format: "json"
  outputPaths:
    - "%s"
  errorOutputPaths:
    - "%s"
app:
  name: "ConfigLogIntegrationTest"
  environment: "testing"
  debug: false
features:
  enableDetailedLogging: true
  logErrorDetails: true
`, logFile, logFile)

	err := os.WriteFile(configFile, []byte(initialConfig), 0644)
	require.NoError(t, err)

	// 加载配置并启动热重载 (Load config and start hot reload)
	var appCfg ApplicationConfig
	cfgManager, err := config.LoadConfigAndWatch(&appCfg,
		config.WithConfigFile(configFile, "yaml"),
		config.WithHotReload(true),
	)
	require.NoError(t, err)

	// 注册日志配置热重载 (Register log config hot reload)
	log.RegisterConfigHotReload(cfgManager)

	// 等待配置加载完成 (Wait for config loading to complete)
	time.Sleep(200 * time.Millisecond)

	// 验证初始日志级别 (Verify initial log level)
	assert.Equal(t, "info", log.Std().GetZapLogger().Level().String())

		// 记录一些测试日志 (Log some test messages)
	log.Debug("This debug message should not appear")
	log.Info("Initial config loaded successfully")
	log.Warn("This is a warning message")
	
	// 同步日志，忽略可能的stdout sync错误 (Sync logs, ignore possible stdout sync errors)
	if syncErr := log.Std().GetZapLogger().Sync(); syncErr != nil {
		// 忽略stdout/stderr的sync错误，这在某些系统上是正常的 (Ignore stdout/stderr sync errors, which are normal on some systems)
		t.Logf("Warning: log sync failed (this may be normal for stdout): %v", syncErr)
	}

	// 更新配置以改变日志级别 (Update config to change log level)
	updatedConfig := fmt.Sprintf(`
server:
  host: "localhost"
  port: 8080
log:
  level: "debug"
  format: "json"
  outputPaths:
    - "%s"
  errorOutputPaths:
    - "%s"
app:
  name: "ConfigLogIntegrationTest"
  environment: "testing"
  debug: true
features:
  enableDetailedLogging: true
  logErrorDetails: true
`, logFile, logFile)

	time.Sleep(100 * time.Millisecond)
	err = os.WriteFile(configFile, []byte(updatedConfig), 0644)
	require.NoError(t, err)

	// 等待热重载完成 (Wait for hot reload to complete)
	time.Sleep(600 * time.Millisecond)

	// 验证日志级别已更新 (Verify log level has been updated)
	assert.Equal(t, "debug", log.Std().GetZapLogger().Level().String())

		// 记录更多测试日志 (Log more test messages)
	log.Debug("This debug message should now appear")
	log.Info("Config hot reload completed successfully")
	
	// 同步日志，忽略可能的stdout sync错误 (Sync logs, ignore possible stdout sync errors)
	if syncErr := log.Std().GetZapLogger().Sync(); syncErr != nil {
		t.Logf("Warning: log sync failed (this may be normal for stdout): %v", syncErr)
	}

	// 验证日志文件内容 (Verify log file content)
	logContent, err := os.ReadFile(logFile)
	require.NoError(t, err)
	logStr := string(logContent)

	assert.NotContains(t, logStr, "debug message should not appear")
	// 检查重新配置后的消息 (Check messages after reconfiguration)
	assert.Contains(t, logStr, "debug message should now appear")
	assert.Contains(t, logStr, "Config hot reload completed successfully")
	// 检查日志系统重新配置的消息 (Check log system reconfiguration message)
	assert.Contains(t, logStr, "successfully reconfigured")
}

// TestErrorsLogIntegration 测试错误和日志模块的集成
// (TestErrorsLogIntegration tests integration between errors and log modules)
func TestErrorsLogIntegration(t *testing.T) {
	tempDir := setupIntegrationTestDir(t)
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "errors.log")

	// 配置日志记录器 (Configure logger)
	logOpts := log.NewOptions()
	logOpts.Level = "debug"
	logOpts.Format = "json"
	logOpts.OutputPaths = []string{logFile}
	logOpts.ErrorOutputPaths = []string{logFile}
	logOpts.Development = false
	logOpts.ContextKeys = []any{log.TraceIDKey, log.RequestIDKey} // 配置要提取的context keys (Configure context keys to extract)

	log.Init(logOpts)
	defer func() {
		defaultOpts := log.NewOptions()
		log.Init(defaultOpts)
	}()

		// 创建带上下文的错误并记录 (Create contextualized errors and log them)
	ctx := context.Background()
	ctx = log.ContextWithTraceID(ctx, "trace-12345")

	// 测试不同类型的错误日志记录 (Test different types of error logging)
	
	// 简单错误 (Simple error)
	simpleErr := lmccerrors.New("simple error occurred")
	log.Ctxw(ctx, "Simple error test", "error", simpleErr)

	// 带错误码的错误 (Error with code)
	validationErr := lmccerrors.NewWithCode(
		lmccerrors.ErrValidation,
		"field 'email' is required",
	)
	log.Errorw("Validation error occurred",
		"field", "email",
		"user_input", "invalid@",
		"error", validationErr,
	)

	// 包装的错误链 (Wrapped error chain)
	dbErr := lmccerrors.New("connection timeout")
	serviceErr := lmccerrors.Wrap(dbErr, "failed to save user")
	apiErr := lmccerrors.WithCode(
		lmccerrors.Wrapf(serviceErr, "API request failed for endpoint %s", "/api/users"),
		lmccerrors.ErrInternalServer,
	)

	log.Errorw("API error with full stack trace",
		"endpoint", "/api/users",
		"method", "POST",
		"trace_id", "trace-12345",
		"error", apiErr,
		"error_stack", fmt.Sprintf("%+v", apiErr),
	)

	// 同步日志 (Sync logs)
	if syncErr := log.Std().GetZapLogger().Sync(); syncErr != nil {
		t.Logf("Warning: log sync failed (this may be normal for stdout): %v", syncErr)
	}

	// 验证日志内容 (Verify log content)
	logContent, err := os.ReadFile(logFile)
	require.NoError(t, err)

	// 解析日志条目 (Parse log entries)
	scanner := bufio.NewScanner(strings.NewReader(string(logContent)))
	var logEntries []map[string]interface{}

	for scanner.Scan() {
		var entry map[string]interface{}
		err := json.Unmarshal([]byte(scanner.Text()), &entry)
		require.NoError(t, err)
		logEntries = append(logEntries, entry)
	}

	require.Len(t, logEntries, 3, "Should have 3 log entries")

	// 验证第一个条目 (Verify first entry)
	firstEntry := logEntries[0]
	assert.Contains(t, firstEntry["M"], "Simple error test")
	assert.Equal(t, "trace-12345", firstEntry["trace_id"])

	// 验证第二个条目 (Verify second entry)
	secondEntry := logEntries[1]
	assert.Contains(t, secondEntry["M"], "Validation error occurred")
	assert.Equal(t, "email", secondEntry["field"])

	// 验证第三个条目 (Verify third entry)
	thirdEntry := logEntries[2]
	assert.Contains(t, thirdEntry["M"], "API error with full stack trace")
	assert.Equal(t, "/api/users", thirdEntry["endpoint"])
	assert.Equal(t, "POST", thirdEntry["method"])
	assert.Contains(t, thirdEntry["error_stack"], "connection timeout")
	assert.Contains(t, thirdEntry["error_stack"], "failed to save user")
	assert.Contains(t, thirdEntry["error_stack"], "API request failed")
}

// TestFullApplicationIntegration 测试完整应用程序集成场景
// (TestFullApplicationIntegration tests full application integration scenario)
func TestFullApplicationIntegration(t *testing.T) {
	tempDir := setupIntegrationTestDir(t)
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "full_app.log")
	configFile := filepath.Join(tempDir, "full_app.yaml")

	// 创建完整的应用配置 (Create complete application config)
	appConfig := fmt.Sprintf(`
server:
  host: "0.0.0.0"
  port: 8080
  mode: "production"
log:
  level: "info"
  format: "json"
  outputPaths:
    - "%s"
  errorOutputPaths:
    - "%s"
  development: false
database:
  type: "postgres"
  host: "localhost"
  port: 5432
app:
  name: "FullIntegrationTestApp"
  environment: "production"
  debug: false
features:
  enableDetailedLogging: true
  logErrorDetails: true
`, logFile, logFile)

	err := os.WriteFile(configFile, []byte(appConfig), 0644)
	require.NoError(t, err)

	// 加载配置 (Load configuration)
	var cfg ApplicationConfig
	cfgManager, err := config.LoadConfigAndWatch(&cfg,
		config.WithConfigFile(configFile, "yaml"),
		config.WithHotReload(true),
	)
	require.NoError(t, err)

	// 注册日志配置热重载 (Register log config hot reload)
	log.RegisterConfigHotReload(cfgManager)
	
	// 初始化日志系统使用配置文件中的设置 (Initialize log system with config file settings)
	logOpts := log.NewOptions()
	logOpts.Level = cfg.Log.Level
	logOpts.Format = cfg.Log.Format
	logOpts.OutputPaths = cfg.Log.OutputPaths
	logOpts.ErrorOutputPaths = cfg.Log.ErrorOutputPaths
	logOpts.Development = cfg.Log.Development
	log.Init(logOpts)
	
	defer func() {
		defaultOpts := log.NewOptions()
		log.Init(defaultOpts)
	}()
	
	time.Sleep(200 * time.Millisecond)

	// 创建应用服务 (Create application service)
	service := newMockApplicationService("MainService")

	// 模拟各种业务场景 (Simulate various business scenarios)
	testCases := []struct {
		name        string
		userID      string
		action      string
		expectError bool
		errorCode   lmccerrors.Coder
	}{
		{
			name:        "Successful operation",
			userID:      "user123",
			action:      "create",
			expectError: false,
		},
		{
			name:        "Validation error",
			userID:      "",
			action:      "create",
			expectError: true,
			errorCode:   lmccerrors.ErrValidation,
		},
		{
			name:        "Forbidden action",
			userID:      "user456",
			action:      "forbidden",
			expectError: true,
			errorCode:   lmccerrors.ErrForbidden,
		},
		{
			name:        "Internal server error",
			userID:      "user789",
			action:      "fail",
			expectError: true,
			errorCode:   lmccerrors.ErrInternalServer,
		},
	}

	// 执行测试场景 (Execute test scenarios)
	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			ctx = log.ContextWithTraceID(ctx, fmt.Sprintf("trace-%d", i))

			err := service.ProcessBusinessLogic(ctx, tc.userID, tc.action)

			if tc.expectError {
				require.Error(t, err)
				if tc.errorCode != nil {
					assert.True(t, lmccerrors.IsCode(err, tc.errorCode),
						"Expected error code %s, got error: %v", tc.errorCode.String(), err)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}

	// 同步日志 (Sync logs)
	if syncErr := log.Std().GetZapLogger().Sync(); syncErr != nil {
		t.Logf("Warning: log sync failed (this may be normal for stdout): %v", syncErr)
	}

	// 验证日志文件内容 (Verify log file content)
	logContent, err := os.ReadFile(logFile)
	require.NoError(t, err)
	logStr := string(logContent)

	// 验证各种场景的日志记录 (Verify logging for various scenarios)
	assert.Contains(t, logStr, "Starting business logic processing")
	assert.Contains(t, logStr, "user123")                               // 成功场景 (Success scenario)
	assert.Contains(t, logStr, "Validation failed")                     // 验证错误 (Validation error)
	assert.Contains(t, logStr, "Forbidden action attempted")            // 禁止操作 (Forbidden action)
	assert.Contains(t, logStr, "Internal server error occurred")        // 内部错误 (Internal error)
	assert.Contains(t, logStr, "Business logic completed successfully") // 成功完成 (Successful completion)

	// 验证错误码在日志中的记录 (Verify error codes in logs)
	assert.Contains(t, logStr, fmt.Sprintf("\"error_code\":%d", lmccerrors.ErrValidation.Code()))
}

// TestConcurrentCrossModuleOperations 测试并发跨模块操作
// (TestConcurrentCrossModuleOperations tests concurrent cross-module operations)
func TestConcurrentCrossModuleOperations(t *testing.T) {
	tempDir := setupIntegrationTestDir(t)
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "concurrent.log")
	configFile := filepath.Join(tempDir, "concurrent.yaml")

	// 配置文件 (Configuration file)
	concurrentConfig := fmt.Sprintf(`
server:
  host: "localhost"
  port: 8080
log:
  level: "debug"
  format: "json"
  outputPaths:
    - "%s"
  errorOutputPaths:
    - "%s"
app:
  name: "ConcurrentTestApp"
  debug: true
features:
  enableDetailedLogging: true
  logErrorDetails: true
`, logFile, logFile)

	err := os.WriteFile(configFile, []byte(concurrentConfig), 0644)
	require.NoError(t, err)

	// 加载配置 (Load configuration)
	var cfg ApplicationConfig
	_, err = config.LoadConfigAndWatch(&cfg,
		config.WithConfigFile(configFile, "yaml"),
		config.WithHotReload(true),
	)
	require.NoError(t, err)

	// 配置日志 (Configure logging)
	logOpts := log.NewOptions()
	logOpts.Level = "debug"
	logOpts.Format = "json"
	logOpts.OutputPaths = []string{logFile}
	logOpts.ErrorOutputPaths = []string{logFile}
	logOpts.ContextKeys = []any{log.TraceIDKey, log.RequestIDKey} // 配置要提取的context keys (Configure context keys to extract)
	log.Init(logOpts)

	defer func() {
		defaultOpts := log.NewOptions()
		log.Init(defaultOpts)
	}()

	// 并发测试 (Concurrent testing)
	numGoroutines := 50
	var wg sync.WaitGroup
	errorGroup := *lmccerrors.NewErrorGroup("concurrent operations")
	var mu sync.Mutex

	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			ctx := context.Background()
			ctx = log.ContextWithTraceID(ctx, fmt.Sprintf("concurrent-trace-%d", id))

			service := newMockApplicationService(fmt.Sprintf("Service-%d", id))

			// 模拟不同的操作 (Simulate different operations)
			var err error
			switch id % 4 {
			case 0:
				err = service.ProcessBusinessLogic(ctx, fmt.Sprintf("user%d", id), "create")
			case 1:
				err = service.ProcessBusinessLogic(ctx, "", "create") // 验证错误 (Validation error)
			case 2:
				err = service.ProcessBusinessLogic(ctx, fmt.Sprintf("user%d", id), "forbidden") // 禁止操作 (Forbidden action)
			case 3:
				err = service.ProcessBusinessLogic(ctx, fmt.Sprintf("user%d", id), "fail") // 内部错误 (Internal error)
			}

			if err != nil {
				mu.Lock()
				errorGroup.Add(lmccerrors.Wrapf(err, "goroutine %d failed", id))
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()

	// 同步日志 (Sync logs)
	if syncErr := log.Std().GetZapLogger().Sync(); syncErr != nil {
		t.Logf("Warning: log sync failed (this may be normal for stdout): %v", syncErr)
	}

	// 验证错误收集 (Verify error collection)
	errors := errorGroup.Errors()
	expectedErrors := numGoroutines * 3 / 4 // 75% 的情况会产生错误 (75% cases will generate errors)
	assert.Equal(t, expectedErrors, len(errors))

	// 验证日志文件内容 (Verify log file content)
	logContent, err := os.ReadFile(logFile)
	require.NoError(t, err)
	logStr := string(logContent)

	// 验证并发日志记录 (Verify concurrent logging)
	assert.Contains(t, logStr, "Starting business logic processing")
	assert.Contains(t, logStr, "Business logic completed successfully")
	assert.Contains(t, logStr, "Validation failed")
	assert.Contains(t, logStr, "Forbidden action attempted")
	assert.Contains(t, logStr, "Internal server error occurred")

	// 验证trace ID的存在 (Verify presence of trace IDs)
	for i := 0; i < 10; i++ { // 检查前10个trace ID (Check first 10 trace IDs)
		assert.Contains(t, logStr, fmt.Sprintf("concurrent-trace-%d", i))
	}
}

// 辅助函数 (Helper functions)

func setupIntegrationTestDir(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "cross_module_integration_test")
	require.NoError(t, err)
	return tempDir
}
