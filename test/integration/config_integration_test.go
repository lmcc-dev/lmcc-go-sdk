/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package integration

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
	lmccerrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConfigIntegration 配置模块集成测试套件
// (TestConfigIntegration config module integration test suite)
type TestConfigIntegration struct {
	tempDir    string
	configFile string
}

// ComplexAppConfig 复杂应用配置结构，用于测试各种配置场景
// (ComplexAppConfig complex application config structure for testing various config scenarios)
type ComplexAppConfig struct {
	config.Config                          // 嵌入基础配置 (Embed base config)
	App           *AppConfig    `mapstructure:"app"`
	Features      *FeatureFlags `mapstructure:"features"`
	External      ExternalAPIs  `mapstructure:"external"` // 值类型嵌套 (Value type nested)
	CustomValues  []string      `mapstructure:"customValues" default:"default1,default2"`
}

type AppConfig struct {
	Name        string `mapstructure:"name" default:"MyApp"`
	Version     string `mapstructure:"version" default:"1.0.0"`
	Environment string `mapstructure:"environment" default:"development"`
	Port        int    `mapstructure:"port" default:"3000"`
}

type FeatureFlags struct {
	EnableNewUI    bool   `mapstructure:"enableNewUI" default:"false"`
	EnableMetrics  bool   `mapstructure:"enableMetrics" default:"true"`
	ExperimentName string `mapstructure:"experimentName" default:"baseline"`
	MaxUsers       *int   `mapstructure:"maxUsers" default:"1000"` // 指针类型 (Pointer type)
}

type ExternalAPIs struct {
	PaymentAPI APIConfig `mapstructure:"paymentAPI"`
	NotifyAPI  APIConfig `mapstructure:"notifyAPI"`
}

type APIConfig struct {
	URL     string `mapstructure:"url" default:"http://localhost:8080"`
	Timeout string `mapstructure:"timeout" default:"30s"`
	APIKey  string `mapstructure:"apiKey"` // 无默认值，需要环境变量或配置文件 (No default, requires env var or config file)
}

// TestComplexConfigurationLoading 测试复杂配置结构的加载
// (TestComplexConfigurationLoading tests loading of complex config structures)
func TestComplexConfigurationLoading(t *testing.T) {
	suite := setupConfigTestSuite(t)
	defer suite.cleanup(t)

	// 创建复杂的配置文件 (Create complex config file)
	configContent := `
server:
  host: "0.0.0.0"
  port: 8080
  mode: "production"
log:
  level: "info"
  format: "json"
database:
  type: "postgres"
  host: "localhost"
  port: 5432
  user: "testuser"
  password: "testpass"
  dbName: "testdb"
app:
  name: "TestApp"
  version: "2.0.0"
  environment: "testing"
  port: 4000
features:
  enableNewUI: true
  enableMetrics: false  # 测试配置文件值能否覆盖默认值true (Test if config file value can override default value true)
  experimentName: "beta_test"
  maxUsers: 5000
external:
  paymentAPI:
    url: "https://api.payment.com"
    timeout: "60s"
    apiKey: "payment_key_123"
  notifyAPI:
    url: "https://api.notify.com"
    timeout: "45s"
customValues:
  - "custom1"
  - "custom2"
  - "custom3"
`

	err := os.WriteFile(suite.configFile, []byte(configContent), 0644)
	require.NoError(t, err)



	var cfg ComplexAppConfig
	err = config.LoadConfig(&cfg,
		config.WithConfigFile(suite.configFile, "yaml"),
		config.WithEnvPrefix("TEST_APP"),
		config.WithEnvVarOverride(true),
	)
	require.NoError(t, err)

	// 验证基础配置 (Verify base config)
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "info", cfg.Log.Level)
	assert.Equal(t, "postgres", cfg.Database.Type)

	// 验证嵌套配置 (Verify nested config)
	require.NotNil(t, cfg.App)
	assert.Equal(t, "TestApp", cfg.App.Name)
	assert.Equal(t, "2.0.0", cfg.App.Version)
	assert.Equal(t, "testing", cfg.App.Environment)
	assert.Equal(t, 4000, cfg.App.Port)

	// 验证特性标志 (Verify feature flags)
	require.NotNil(t, cfg.Features)
	assert.True(t, cfg.Features.EnableNewUI)
	
	// 关键测试：配置文件中设置为false应该覆盖默认值true (Key test: config file false should override default true)
	assert.False(t, cfg.Features.EnableMetrics, "Config file value 'false' should override default value 'true'")
	
	assert.Equal(t, "beta_test", cfg.Features.ExperimentName)
	require.NotNil(t, cfg.Features.MaxUsers)
	assert.Equal(t, 5000, *cfg.Features.MaxUsers)

	// 验证外部API配置 (Verify external API config)
	assert.Equal(t, "https://api.payment.com", cfg.External.PaymentAPI.URL)
	assert.Equal(t, "60s", cfg.External.PaymentAPI.Timeout)
	assert.Equal(t, "payment_key_123", cfg.External.PaymentAPI.APIKey)

	// 验证数组字段 (Verify array field)
	assert.Equal(t, []string{"custom1", "custom2", "custom3"}, cfg.CustomValues)
}

// TestEnvironmentVariableOverrides 测试环境变量覆盖功能
// (TestEnvironmentVariableOverrides tests environment variable override functionality)
func TestEnvironmentVariableOverrides(t *testing.T) {
	suite := setupConfigTestSuite(t)
	defer suite.cleanup(t)

	// 设置环境变量 (Set environment variables)
	envVars := map[string]string{
		"TEST_APP_SERVER_HOST":              "192.168.1.100",
		"TEST_APP_SERVER_PORT":              "9090",
		"TEST_APP_APP_NAME":                 "EnvApp",
		"TEST_APP_FEATURES_ENABLENEWUI":     "true",
		"TEST_APP_FEATURES_MAXUSERS":        "2000",
		"TEST_APP_EXTERNAL_PAYMENTAPI_URL":  "https://env.payment.com",
		"TEST_APP_EXTERNAL_PAYMENTAPI_APIKEY": "env_payment_key",
	}

	// 设置并清理环境变量 (Set and cleanup environment variables)
	for key, value := range envVars {
		t.Setenv(key, value)
	}

	// 创建基础配置文件 (Create base config file)
	configContent := `
server:
  host: "localhost"
  port: 8080
app:
  name: "FileApp"
  environment: "production"
features:
  enableNewUI: false
  maxUsers: 1000
external:
  paymentAPI:
    url: "https://file.payment.com"
    timeout: "30s"
`

	err := os.WriteFile(suite.configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	var cfg ComplexAppConfig
	err = config.LoadConfig(&cfg,
		config.WithConfigFile(suite.configFile, "yaml"),
		config.WithEnvPrefix("TEST_APP"),
		config.WithEnvVarOverride(true),
	)
	require.NoError(t, err)

	// 验证环境变量覆盖了配置文件值 (Verify env vars override config file values)
	assert.Equal(t, "192.168.1.100", cfg.Server.Host, "Environment variable should override config file")
	assert.Equal(t, 9090, cfg.Server.Port, "Environment variable should override config file")
	assert.Equal(t, "EnvApp", cfg.App.Name, "Environment variable should override config file")
	assert.True(t, cfg.Features.EnableNewUI, "Environment variable should override config file")
	require.NotNil(t, cfg.Features.MaxUsers)
	assert.Equal(t, 2000, *cfg.Features.MaxUsers, "Environment variable should override config file")
	assert.Equal(t, "https://env.payment.com", cfg.External.PaymentAPI.URL, "Environment variable should override config file")
	assert.Equal(t, "env_payment_key", cfg.External.PaymentAPI.APIKey, "Environment variable should set missing field")

	// 验证未被环境变量覆盖的值保持原样 (Verify values not overridden by env vars remain unchanged)
	assert.Equal(t, "production", cfg.App.Environment, "Non-overridden config should remain from file")
	assert.Equal(t, "30s", cfg.External.PaymentAPI.Timeout, "Non-overridden config should remain from file")
}

// TestConfigurationHotReload 测试配置热重载功能
// (TestConfigurationHotReload tests configuration hot reload functionality)
func TestConfigurationHotReload(t *testing.T) {
	suite := setupConfigTestSuite(t)
	defer suite.cleanup(t)

	// 初始配置 (Initial config)
	initialConfig := `
server:
  host: "localhost"
  port: 8080
app:
  name: "InitialApp"
  environment: "development"
features:
  enableNewUI: false
  experimentName: "initial"
`

	err := os.WriteFile(suite.configFile, []byte(initialConfig), 0644)
	require.NoError(t, err)

	var cfg ComplexAppConfig
	cfgManager, err := config.LoadConfigAndWatch(&cfg,
		config.WithConfigFile(suite.configFile, "yaml"),
		config.WithHotReload(true),
	)
	require.NoError(t, err)

	// 验证初始配置 (Verify initial config)
	assert.Equal(t, "localhost", cfg.Server.Host)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "InitialApp", cfg.App.Name)
	assert.False(t, cfg.Features.EnableNewUI)
	assert.Equal(t, "initial", cfg.Features.ExperimentName)

	// 设置回调计数器 (Setup callback counter)
	var callbackCount int
	var mu sync.Mutex
	var lastCallbackCfg *ComplexAppConfig

	cfgManager.RegisterCallback(func(v *viper.Viper, cfgPtr any) error {
		mu.Lock()
		defer mu.Unlock()
		callbackCount++
		if appCfg, ok := cfgPtr.(*ComplexAppConfig); ok {
			lastCallbackCfg = appCfg
		}
		return nil
	})

	// 更新配置文件 (Update config file)
	updatedConfig := `
server:
  host: "0.0.0.0"
  port: 9000
app:
  name: "UpdatedApp"
  environment: "production"
features:
  enableNewUI: true
  experimentName: "updated"
  maxUsers: 3000
`

	time.Sleep(100 * time.Millisecond) // 确保文件监视器准备就绪 (Ensure file watcher is ready)
	err = os.WriteFile(suite.configFile, []byte(updatedConfig), 0644)
	require.NoError(t, err)

	// 等待热重载完成 (Wait for hot reload to complete)
	time.Sleep(500 * time.Millisecond)

	// 验证配置已更新 (Verify config has been updated)
	mu.Lock()
	assert.Greater(t, callbackCount, 0, "Callback should have been triggered")
	require.NotNil(t, lastCallbackCfg, "Callback should have received updated config")
	assert.Equal(t, "0.0.0.0", lastCallbackCfg.Server.Host)
	assert.Equal(t, 9000, lastCallbackCfg.Server.Port)
	assert.Equal(t, "UpdatedApp", lastCallbackCfg.App.Name)
	assert.True(t, lastCallbackCfg.Features.EnableNewUI)
	assert.Equal(t, "updated", lastCallbackCfg.Features.ExperimentName)
	require.NotNil(t, lastCallbackCfg.Features.MaxUsers)
	assert.Equal(t, 3000, *lastCallbackCfg.Features.MaxUsers)
	mu.Unlock()

	// 验证原始配置对象也已更新 (Verify original config object is also updated)
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, 9000, cfg.Server.Port)
	assert.Equal(t, "UpdatedApp", cfg.App.Name)
}

// TestConfigurationErrorHandling 测试配置错误处理
// (TestConfigurationErrorHandling tests configuration error handling)
func TestConfigurationErrorHandling(t *testing.T) {
	suite := setupConfigTestSuite(t)
	defer suite.cleanup(t)

	tests := []struct {
		name        string
		configContent string
		expectError bool
		errorCode   lmccerrors.Coder
	}{
		{
			name: "Invalid YAML syntax",
			configContent: `
server:
  host: "localhost"
  port: invalid_port_format
  invalid_yaml_syntax
`,
			expectError: true,
			errorCode:   lmccerrors.ErrConfigFileRead,
		},
		{
			name: "Missing config file",
			configContent: "", // Will use non-existent file
			expectError: true,
			errorCode:   lmccerrors.ErrConfigFileRead,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cfg ComplexAppConfig
			var err error

			if tt.name == "Missing config file" {
				// 使用不存在的文件 (Use non-existent file)
				err = config.LoadConfig(&cfg,
					config.WithConfigFile("/non/existent/config.yaml", "yaml"),
				)
			} else {
				// 写入无效配置内容 (Write invalid config content)
				tempFile := filepath.Join(suite.tempDir, "invalid_config.yaml")
				writeErr := os.WriteFile(tempFile, []byte(tt.configContent), 0644)
				require.NoError(t, writeErr)

				err = config.LoadConfig(&cfg,
					config.WithConfigFile(tempFile, "yaml"),
				)
			}

			if tt.expectError {
				require.Error(t, err)
				assert.True(t, lmccerrors.IsCode(err, tt.errorCode),
					"Expected error code %s, got: %v", tt.errorCode.String(), err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestSectionSpecificCallbacks 测试特定节的回调功能
// (TestSectionSpecificCallbacks tests section-specific callback functionality)
func TestSectionSpecificCallbacks(t *testing.T) {
	suite := setupConfigTestSuite(t)
	defer suite.cleanup(t)

	initialConfig := `
server:
  host: "localhost"
  port: 8080
log:
  level: "info"
  format: "text"
app:
  name: "TestApp"
`

	err := os.WriteFile(suite.configFile, []byte(initialConfig), 0644)
	require.NoError(t, err)

	var cfg ComplexAppConfig
	cfgManager, err := config.LoadConfigAndWatch(&cfg,
		config.WithConfigFile(suite.configFile, "yaml"),
		config.WithHotReload(true),
	)
	require.NoError(t, err)

	// 注册节特定回调 (Register section-specific callbacks)
	var serverCallbackCount, logCallbackCount, appCallbackCount int
	var mu sync.Mutex

	cfgManager.RegisterSectionChangeCallback("server", func(v *viper.Viper) error {
		mu.Lock()
		serverCallbackCount++
		mu.Unlock()
		return nil
	})

	cfgManager.RegisterSectionChangeCallback("log", func(v *viper.Viper) error {
		mu.Lock()
		logCallbackCount++
		mu.Unlock()
		return nil
	})

	cfgManager.RegisterSectionChangeCallback("app", func(v *viper.Viper) error {
		mu.Lock()
		appCallbackCount++
		mu.Unlock()
		return nil
	})

	// 更新配置以触发回调 (Update config to trigger callbacks)
	updatedConfig := `
server:
  host: "0.0.0.0"
  port: 9000
log:
  level: "debug"
  format: "json"
app:
  name: "UpdatedApp"
  version: "2.0.0"
`

	time.Sleep(100 * time.Millisecond)
	err = os.WriteFile(suite.configFile, []byte(updatedConfig), 0644)
	require.NoError(t, err)
	time.Sleep(500 * time.Millisecond)

	// 验证所有节回调都被触发 (Verify all section callbacks were triggered)
	mu.Lock()
	assert.Greater(t, serverCallbackCount, 0, "Server section callback should be triggered")
	assert.Greater(t, logCallbackCount, 0, "Log section callback should be triggered")
	assert.Greater(t, appCallbackCount, 0, "App section callback should be triggered")
	mu.Unlock()
}

// TestConcurrentConfigAccess 测试并发配置访问
// (TestConcurrentConfigAccess tests concurrent config access)
func TestConcurrentConfigAccess(t *testing.T) {
	suite := setupConfigTestSuite(t)
	defer suite.cleanup(t)

	configContent := `
server:
  host: "localhost"
  port: 8080
app:
  name: "ConcurrentApp"
`

	err := os.WriteFile(suite.configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	var cfg ComplexAppConfig
	_, err = config.LoadConfigAndWatch(&cfg,
		config.WithConfigFile(suite.configFile, "yaml"),
		config.WithHotReload(true),
	)
	require.NoError(t, err)

	// 并发访问全局配置 (Concurrent access to global config)
	var wg sync.WaitGroup
	numGoroutines := 50
	results := make([]string, numGoroutines)

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			globalCfg := config.GetGlobalCfg()
			if globalCfg != nil && globalCfg.Server != nil {
				results[idx] = globalCfg.Server.Host
			}
		}(i)
	}

	wg.Wait()

	// 验证所有goroutine都获得了正确的配置值 (Verify all goroutines got correct config values)
	for i, result := range results {
		assert.Equal(t, "localhost", result, "Goroutine %d should get correct config value", i)
	}
}

// TestDefaultValueApplication 测试默认值应用
// (TestDefaultValueApplication tests default value application)
func TestDefaultValueApplication(t *testing.T) {
	suite := setupConfigTestSuite(t)
	defer suite.cleanup(t)

	// 创建只包含部分字段的配置文件 (Create config file with only partial fields)
	partialConfig := `
server:
  host: "custom.host"
# 其他字段将使用默认值 (Other fields will use default values)
`

	err := os.WriteFile(suite.configFile, []byte(partialConfig), 0644)
	require.NoError(t, err)

	var cfg ComplexAppConfig
	err = config.LoadConfig(&cfg,
		config.WithConfigFile(suite.configFile, "yaml"),
	)
	require.NoError(t, err)

	// 验证显式设置的值 (Verify explicitly set values)
	assert.Equal(t, "custom.host", cfg.Server.Host)

	// 验证默认值被正确应用 (Verify default values are correctly applied)
	assert.Equal(t, 8080, cfg.Server.Port) // 来自默认值 (From default value)
	assert.Equal(t, "MyApp", cfg.App.Name) // 来自默认值 (From default value)
	assert.Equal(t, "1.0.0", cfg.App.Version) // 来自默认值 (From default value)
	assert.Equal(t, 3000, cfg.App.Port) // 来自默认值 (From default value)
	assert.False(t, cfg.Features.EnableNewUI) // 来自默认值 (From default value)
	assert.True(t, cfg.Features.EnableMetrics) // 来自默认值 (From default value)
	assert.Equal(t, []string{"default1", "default2"}, cfg.CustomValues) // 来自默认值 (From default value)
	assert.Equal(t, "http://localhost:8080", cfg.External.PaymentAPI.URL) // 来自默认值 (From default value)
}

// 辅助函数 (Helper functions)

func setupConfigTestSuite(t *testing.T) *TestConfigIntegration {
	tempDir, err := os.MkdirTemp("", "config_integration_test")
	require.NoError(t, err)
	
	configFile := filepath.Join(tempDir, "config.yaml")
	
	return &TestConfigIntegration{
		tempDir:    tempDir,
		configFile: configFile,
	}
}

func (suite *TestConfigIntegration) cleanup(t *testing.T) {
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
} 