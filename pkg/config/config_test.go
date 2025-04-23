/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Unit tests for the config package.
 */

package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Test Setup ---

// testAppConfig 嵌入 SDK Config 并添加自定义字段用于测试
// (testAppConfig embeds SDK Config and adds custom fields for testing)
type testAppConfig struct {
	Config                             // 嵌入 SDK 基础配置 (Embed SDK base config)
	CustomFeature *customFeatureConfig `mapstructure:"customFeature"`
}

type customFeatureConfig struct {
	APIKey    string `mapstructure:"apiKey"`
	RateLimit int    `mapstructure:"rateLimit"`
	Enabled   bool   `mapstructure:"enabled"`
}

// Helper function to create a temporary config file
func createTempConfigFile(t *testing.T, content string, fileType string) (string, func()) {
	t.Helper()
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_config."+fileType)
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	require.NoError(t, err, "Failed to write temp config file")

	// Cleanup function
	cleanup := func() {
		// No explicit cleanup needed for t.TempDir() files
	}
	return tmpFile, cleanup
}

// Helper function to initialize pointer fields in the embedded Config
func initializeTestConfig(cfg *testAppConfig) {
	if cfg.Config.Server == nil {
		cfg.Config.Server = &ServerConfig{}
	}
	if cfg.Config.Log == nil {
		cfg.Config.Log = &LogConfig{}
	}
	if cfg.Config.Database == nil {
		cfg.Config.Database = &DatabaseConfig{}
	}
	if cfg.Config.Tracing == nil {
		cfg.Config.Tracing = &TracingConfig{}
	}
	if cfg.Config.Metrics == nil {
		cfg.Config.Metrics = &MetricsConfig{}
	}
	// CustomFeature is already a pointer in testAppConfig, Viper should handle it
}

// --- Test Cases ---

func TestLoadConfig_YAML_Valid(t *testing.T) {
	yamlContent := `
server:
  host: "127.0.0.1"
  port: 9999
  mode: "debug"
  readTimeout: 10s
  writeTimeout: 15s
  gracefulShutdownTimeout: 20s
log:
  level: "debug"
  format: "json"
  output: "file"
  filename: "test.log"
  maxSize: 200
  maxBackups: 10
  maxAge: 14
  compress: true
database:
  type: "postgres"
  host: "db.example.com"
  port: 5432
  user: "testuser"
  password: "testpassword"
  dbName: "testdb"
  maxIdleConns: 5
  maxOpenConns: 50
  connMaxLifetime: 30m
tracing:
  enabled: true
  provider: "jaeger"
  endpoint: "http://jaeger:14268/api/traces"
  samplerType: "probabilistic"
  samplerParam: 0.5
metrics:
  enabled: true
  provider: "prometheus"
  port: 9191
  path: "/testmetrics"
customFeature:
  apiKey: "test-api-key"
  rateLimit: 100
  enabled: true
`
	configFile, cleanup := createTempConfigFile(t, yamlContent, "yaml")
	defer cleanup()

	var loadedCfg testAppConfig
	initializeTestConfig(&loadedCfg) // Initialize before loading

	err := LoadConfig(&loadedCfg, WithConfigFile(configFile, "yaml"))
	require.NoError(t, err, "LoadConfig should not return an error for valid YAML")

	// Assert Server Config
	require.NotNil(t, loadedCfg.Server, "Server config should be loaded")
	assert.Equal(t, "127.0.0.1", loadedCfg.Server.Host)
	assert.Equal(t, 9999, loadedCfg.Server.Port)
	assert.Equal(t, "debug", loadedCfg.Server.Mode)
	assert.Equal(t, 10*time.Second, loadedCfg.Server.ReadTimeout)
	assert.Equal(t, 15*time.Second, loadedCfg.Server.WriteTimeout)
	assert.Equal(t, 20*time.Second, loadedCfg.Server.GracefulShutdownTimeout)

	// Assert Log Config
	require.NotNil(t, loadedCfg.Log, "Log config should be loaded")
	assert.Equal(t, "debug", loadedCfg.Log.Level)
	assert.Equal(t, "json", loadedCfg.Log.Format)
	assert.Equal(t, "file", loadedCfg.Log.Output)
	assert.Equal(t, "test.log", loadedCfg.Log.Filename)
	assert.Equal(t, 200, loadedCfg.Log.MaxSize)
	assert.Equal(t, 10, loadedCfg.Log.MaxBackups)
	assert.Equal(t, 14, loadedCfg.Log.MaxAge)
	assert.True(t, loadedCfg.Log.Compress)

	// Assert Database Config
	require.NotNil(t, loadedCfg.Database, "Database config should be loaded")
	assert.Equal(t, "postgres", loadedCfg.Database.Type)
	assert.Equal(t, "db.example.com", loadedCfg.Database.Host)
	assert.Equal(t, 5432, loadedCfg.Database.Port)
	assert.Equal(t, "testuser", loadedCfg.Database.User)
	assert.Equal(t, "testpassword", loadedCfg.Database.Password)
	assert.Equal(t, "testdb", loadedCfg.Database.DBName)
	assert.Equal(t, 5, loadedCfg.Database.MaxIdleConns)
	assert.Equal(t, 50, loadedCfg.Database.MaxOpenConns)
	assert.Equal(t, 30*time.Minute, loadedCfg.Database.ConnMaxLifetime)

	// Assert Tracing Config
	require.NotNil(t, loadedCfg.Tracing, "Tracing config should be loaded")
	assert.True(t, loadedCfg.Tracing.Enabled)
	assert.Equal(t, "jaeger", loadedCfg.Tracing.Provider)
	assert.Equal(t, "http://jaeger:14268/api/traces", loadedCfg.Tracing.Endpoint)
	assert.Equal(t, "probabilistic", loadedCfg.Tracing.SamplerType)
	assert.Equal(t, 0.5, loadedCfg.Tracing.SamplerParam)

	// Assert Metrics Config
	require.NotNil(t, loadedCfg.Metrics, "Metrics config should be loaded")
	assert.True(t, loadedCfg.Metrics.Enabled)
	assert.Equal(t, "prometheus", loadedCfg.Metrics.Provider)
	assert.Equal(t, 9191, loadedCfg.Metrics.Port)
	assert.Equal(t, "/testmetrics", loadedCfg.Metrics.Path)

	// Assert Custom Feature Config
	require.NotNil(t, loadedCfg.CustomFeature, "CustomFeature config should be loaded")
	assert.Equal(t, "test-api-key", loadedCfg.CustomFeature.APIKey)
	assert.Equal(t, 100, loadedCfg.CustomFeature.RateLimit)
	assert.True(t, loadedCfg.CustomFeature.Enabled)
}

func TestLoadConfig_TOML_Valid(t *testing.T) {
	tomlContent := `
[server]
host = "192.168.1.1"
port = 8888
mode = "test"
readTimeout = "8s"
writeTimeout = "12s"
gracefulShutdownTimeout = "18s"

[log]
level = "warn"
format = "text"
output = "stderr"

[database]
type = "mysql"
host = "mysql.internal"
port = 3307
user = "tomluser"
password = "tomlpass"
dbName = "tomldb"
maxIdleConns = 20
maxOpenConns = 150
connMaxLifetime = "2h"

[tracing]
enabled = true
provider = "zipkin"
endpoint = "http://zipkin:9411/api/v2/spans"
samplerType = "const"
samplerParam = 1.0

[metrics]
enabled = true
provider = "prometheus"
port = 9292
path = "/tomlmetrics"

[customFeature]
apiKey = "toml-api-key"
rateLimit = 50
enabled = false
`
	configFile, cleanup := createTempConfigFile(t, tomlContent, "toml")
	defer cleanup()

	var loadedCfg testAppConfig
	initializeTestConfig(&loadedCfg) // Initialize before loading

	err := LoadConfig(&loadedCfg, WithConfigFile(configFile, "toml"))
	require.NoError(t, err, "LoadConfig should not return an error for valid TOML")

	// Assert Server Config
	require.NotNil(t, loadedCfg.Server, "Server config should be loaded")
	assert.Equal(t, "192.168.1.1", loadedCfg.Server.Host)
	assert.Equal(t, 8888, loadedCfg.Server.Port)
	assert.Equal(t, "test", loadedCfg.Server.Mode)
	assert.Equal(t, 8*time.Second, loadedCfg.Server.ReadTimeout) // Viper reads TOML duration strings
	assert.Equal(t, 12*time.Second, loadedCfg.Server.WriteTimeout)
	assert.Equal(t, 18*time.Second, loadedCfg.Server.GracefulShutdownTimeout)

	// Assert Log Config (Only specified fields)
	require.NotNil(t, loadedCfg.Log, "Log config should be loaded")
	assert.Equal(t, "warn", loadedCfg.Log.Level)
	assert.Equal(t, "text", loadedCfg.Log.Format)
	assert.Equal(t, "stderr", loadedCfg.Log.Output)
	// Check defaults for unspecified fields
	assert.Equal(t, "app.log", loadedCfg.Log.Filename) // Default
	assert.Equal(t, 100, loadedCfg.Log.MaxSize)        // Default

	// Assert Database Config
	require.NotNil(t, loadedCfg.Database, "Database config should be loaded")
	assert.Equal(t, "mysql", loadedCfg.Database.Type)
	assert.Equal(t, "mysql.internal", loadedCfg.Database.Host)
	assert.Equal(t, 3307, loadedCfg.Database.Port)
	assert.Equal(t, "tomluser", loadedCfg.Database.User)
	assert.Equal(t, "tomlpass", loadedCfg.Database.Password)
	assert.Equal(t, "tomldb", loadedCfg.Database.DBName)
	assert.Equal(t, 20, loadedCfg.Database.MaxIdleConns)
	assert.Equal(t, 150, loadedCfg.Database.MaxOpenConns)
	assert.Equal(t, 2*time.Hour, loadedCfg.Database.ConnMaxLifetime) // Viper reads TOML duration strings

	// Assert Tracing Config
	require.NotNil(t, loadedCfg.Tracing, "Tracing config should be loaded")
	assert.True(t, loadedCfg.Tracing.Enabled)
	assert.Equal(t, "zipkin", loadedCfg.Tracing.Provider)
	assert.Equal(t, "http://zipkin:9411/api/v2/spans", loadedCfg.Tracing.Endpoint)
	assert.Equal(t, "const", loadedCfg.Tracing.SamplerType)
	assert.Equal(t, 1.0, loadedCfg.Tracing.SamplerParam)

	// Assert Metrics Config
	require.NotNil(t, loadedCfg.Metrics, "Metrics config should be loaded")
	assert.True(t, loadedCfg.Metrics.Enabled)
	assert.Equal(t, "prometheus", loadedCfg.Metrics.Provider)
	assert.Equal(t, 9292, loadedCfg.Metrics.Port)
	assert.Equal(t, "/tomlmetrics", loadedCfg.Metrics.Path)

	// Assert Custom Feature Config
	require.NotNil(t, loadedCfg.CustomFeature, "CustomFeature config should be loaded")
	assert.Equal(t, "toml-api-key", loadedCfg.CustomFeature.APIKey)
	assert.Equal(t, 50, loadedCfg.CustomFeature.RateLimit)
	assert.False(t, loadedCfg.CustomFeature.Enabled)
}

func TestLoadConfig_EnvOverride(t *testing.T) {
	// Setup: Create a base YAML file
	yamlContent := `
server:
  port: 8080 # This should be overridden by env var
  host: "localhost"
log:
  level: "info" # This should be overridden
database:
  password: "yamlpassword" # This should be overridden
customFeature:
  apiKey: "yaml-key" # This should be overridden
  rateLimit: 100
`
	configFile, cleanup := createTempConfigFile(t, yamlContent, "yaml")
	defer cleanup()

	// Setup: Set environment variables
	// Note the prefix "LMCC" and the replacer "." -> "_"
	t.Setenv("LMCC_SERVER_PORT", "9999")
	t.Setenv("LMCC_LOG_LEVEL", "error")
	t.Setenv("LMCC_DATABASE_PASSWORD", "envpassword")
	t.Setenv("LMCC_CUSTOMFEATURE_APIKEY", "env-key")
	t.Setenv("LMCC_TRACING_ENABLED", "true") // Set a value not in the file

	var loadedCfg testAppConfig
	initializeTestConfig(&loadedCfg) // Initialize before loading

	err := LoadConfig(&loadedCfg, WithConfigFile(configFile, "yaml"))
	require.NoError(t, err, "LoadConfig should succeed with environment overrides")

	// Assert overridden values
	require.NotNil(t, loadedCfg.Server, "Server config should be loaded")
	assert.Equal(t, 9999, loadedCfg.Server.Port, "Server port should be overridden by env var")
	assert.Equal(t, "localhost", loadedCfg.Server.Host, "Server host should be loaded from file") // Check non-overridden value

	require.NotNil(t, loadedCfg.Log, "Log config should be loaded")
	assert.Equal(t, "error", loadedCfg.Log.Level, "Log level should be overridden by env var")

	require.NotNil(t, loadedCfg.Database, "Database config should be loaded")
	assert.Equal(t, "envpassword", loadedCfg.Database.Password, "Database password should be overridden by env var")

	require.NotNil(t, loadedCfg.CustomFeature, "CustomFeature config should be loaded")
	assert.Equal(t, "env-key", loadedCfg.CustomFeature.APIKey, "CustomFeature apiKey should be overridden by env var")
	assert.Equal(t, 100, loadedCfg.CustomFeature.RateLimit, "CustomFeature rateLimit should be loaded from file") // Check non-overridden value

	// Assert value set only by environment variable (and default exists)
	require.NotNil(t, loadedCfg.Tracing, "Tracing config should be loaded even if only set by env")
	assert.True(t, loadedCfg.Tracing.Enabled, "Tracing enabled should be set by env var")
	assert.Equal(t, "jaeger", loadedCfg.Tracing.Provider) // Check default
}

func TestLoadConfig_Defaults(t *testing.T) {
	// Setup: Create an empty YAML file to force using defaults/env vars
	configFile, cleanup := createTempConfigFile(t, "", "yaml")
	defer cleanup()

	// Setup: Unset potentially interfering env vars (best effort for isolation)
	// Use t.Setenv to ensure cleanup after test
	unsetEnv := func(key string) {
		originalValue, wasSet := os.LookupEnv(key)
		t.Setenv(key, "")    // Attempt to unset using t.Setenv
		_ = os.Unsetenv(key) // Also try direct unset
		if wasSet {
			t.Cleanup(func() { // Ensure original value is restored if it was set
				_ = os.Setenv(key, originalValue)
			})
		} else {
			t.Cleanup(func() {
				_ = os.Unsetenv(key)
			})
		}
	}

	unsetEnv("LMCC_SERVER_HOST")
	unsetEnv("LMCC_SERVER_PORT")
	unsetEnv("LMCC_LOG_LEVEL")
	unsetEnv("LMCC_DATABASE_HOST") // Explicitly unset database host
	unsetEnv("LMCC_DATABASE_PORT")
	unsetEnv("LMCC_DATABASE_USER")
	unsetEnv("LMCC_DATABASE_PASSWORD")
	unsetEnv("LMCC_TRACING_ENABLED")
	unsetEnv("LMCC_CUSTOMFEATURE_APIKEY") // Unset custom feature env var too

	// Verify that the critical env var is indeed unset before LoadConfig
	hostEnvVal := os.Getenv("LMCC_DATABASE_HOST")
	require.Empty(t, hostEnvVal, "LMCC_DATABASE_HOST should be empty before LoadConfig in Defaults test")

	var loadedCfg testAppConfig
	initializeTestConfig(&loadedCfg) // Initialize before loading

	err := LoadConfig(&loadedCfg, WithConfigFile(configFile, "yaml"))
	// We expect a file not found message printed to stdout, but no error returned
	require.NoError(t, err, "LoadConfig should not return error when config file is empty/not found and defaults are used")

	// --- Mapstructure Isolation Debug ---
	/*
		log.Println("DEBUG: Isolating mapstructure decode with empty map...")
		emptyMap := make(map[string]interface{})
		tempDecoderCfg := &mapstructure.DecoderConfig{
			Result: &loadedCfg, // Decode into the existing loadedCfg
			WeaklyTypedInput: true, // Keep settings consistent with LoadConfig
			TagName: "mapstructure",
			Squash: true,
		}
		tempDecoder, errDecoderNew := mapstructure.NewDecoder(tempDecoderCfg)
		require.NoError(t, errDecoderNew, "Failed to create temp mapstructure decoder for isolation test")
		errDecode := tempDecoder.Decode(emptyMap)
		require.NoError(t, errDecode, "Manual empty map decode failed in isolation test")
		log.Printf("DEBUG: Value after manual empty map decode - Database.Host = '%s'", loadedCfg.Database.Host)
	*/
	// --- End Mapstructure Isolation Debug ---

	// Assert default values defined in setDefaults() or struct tags
	require.NotNil(t, loadedCfg.Server, "Server config should have defaults")
	assert.Equal(t, 8080, loadedCfg.Server.Port)
	assert.Equal(t, "0.0.0.0", loadedCfg.Server.Host)
	assert.Equal(t, "release", loadedCfg.Server.Mode)
	assert.Equal(t, 5*time.Second, loadedCfg.Server.ReadTimeout)
	assert.Equal(t, 10*time.Second, loadedCfg.Server.WriteTimeout)
	assert.Equal(t, 10*time.Second, loadedCfg.Server.GracefulShutdownTimeout)

	require.NotNil(t, loadedCfg.Log, "Log config should have defaults")
	assert.Equal(t, "info", loadedCfg.Log.Level)
	assert.Equal(t, "text", loadedCfg.Log.Format)
	assert.Equal(t, "stdout", loadedCfg.Log.Output)
	assert.Equal(t, "app.log", loadedCfg.Log.Filename)
	assert.Equal(t, 100, loadedCfg.Log.MaxSize)
	assert.Equal(t, 5, loadedCfg.Log.MaxBackups)
	assert.Equal(t, 7, loadedCfg.Log.MaxAge)
	assert.False(t, loadedCfg.Log.Compress)

	require.NotNil(t, loadedCfg.Database, "Database config should have defaults")
	assert.Equal(t, "mysql", loadedCfg.Database.Type)
	assert.Equal(t, 10, loadedCfg.Database.MaxIdleConns)
	assert.Equal(t, 100, loadedCfg.Database.MaxOpenConns)
	assert.Equal(t, time.Hour, loadedCfg.Database.ConnMaxLifetime)
	assert.Equal(t, "localhost", loadedCfg.Database.Host)
	assert.Empty(t, loadedCfg.Database.User)
	assert.Empty(t, loadedCfg.Database.Password)
	assert.Empty(t, loadedCfg.Database.DBName)

	require.NotNil(t, loadedCfg.Tracing, "Tracing config should have defaults")
	assert.False(t, loadedCfg.Tracing.Enabled)
	assert.Equal(t, "jaeger", loadedCfg.Tracing.Provider)
	assert.Equal(t, "const", loadedCfg.Tracing.SamplerType)
	assert.Equal(t, 1.0, loadedCfg.Tracing.SamplerParam)
	assert.Empty(t, loadedCfg.Tracing.Endpoint)

	require.NotNil(t, loadedCfg.Metrics, "Metrics config should have defaults")
	assert.False(t, loadedCfg.Metrics.Enabled)
	assert.Equal(t, "prometheus", loadedCfg.Metrics.Provider)
	assert.Equal(t, 9090, loadedCfg.Metrics.Port)
	assert.Equal(t, "/metrics", loadedCfg.Metrics.Path)

	// Custom feature has no defaults set in SDK, should be nil or zero struct
	// Depending on how viper initializes embedded structs, CustomFeature might be nil or a zero-value struct
	if loadedCfg.CustomFeature != nil {
		assert.Empty(t, loadedCfg.CustomFeature.APIKey, "Custom feature apiKey should be empty (zero value)")
		assert.Zero(t, loadedCfg.CustomFeature.RateLimit, "Custom feature rateLimit should be zero")
		assert.False(t, loadedCfg.CustomFeature.Enabled, "Custom feature enabled should be false (zero value)")
	} else {
		// If viper leaves it nil when no config section exists, that's also acceptable
		assert.Nil(t, loadedCfg.CustomFeature, "CustomFeature should be nil if not defined in config/env and no defaults exist")
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	// Setup: Provide a path that does not exist
	nonExistentPath := filepath.Join(t.TempDir(), "non_existent_config.yaml")

	var loadedCfg testAppConfig
	initializeTestConfig(&loadedCfg) // Initialize before loading

	// Execute: LoadConfig should not return an error if file is not found,
	// as it should fall back to defaults/env vars. It should print a message though.
	// We can't easily capture stdout here without more complex setup, so we just check the error.
	err := LoadConfig(&loadedCfg, WithConfigFile(nonExistentPath, "yaml"))
	require.NoError(t, err, "LoadConfig should not return error when config file is not found")

	// Assert: Check if defaults are loaded correctly (similar to TestLoadConfig_Defaults)
	require.NotNil(t, loadedCfg.Server, "Server config should have defaults even if file not found")
	assert.Equal(t, 8080, loadedCfg.Server.Port) // Check one default value
	require.NotNil(t, loadedCfg.Log, "Log config should have defaults even if file not found")
	assert.Equal(t, "info", loadedCfg.Log.Level) // Check one default value
}

func TestLoadConfig_InvalidFileContent(t *testing.T) {
	// Setup: Create a YAML file with invalid syntax
	invalidYamlContent := `
server:
  port: 8080
  host: "localhost
log: level: "info" # Invalid YAML syntax (missing indentation, extra quote)
`
	configFile, cleanup := createTempConfigFile(t, invalidYamlContent, "yaml")
	defer cleanup()

	var loadedCfg testAppConfig
	initializeTestConfig(&loadedCfg) // Initialize before loading

	// Execute & Assert: LoadConfig should return an error for invalid content
	err := LoadConfig(&loadedCfg, WithConfigFile(configFile, "yaml"))
	require.Error(t, err, "LoadConfig should return error for invalid YAML content")
	// Optionally check the error type or message if needed for more specific validation
	// assert.Contains(t, err.Error(), "failed to read config file")
}

func TestLoadConfig_InvalidConfigType(t *testing.T) {
	// Setup: Create a valid YAML file but try to load as unsupported type
	yamlContent := `
server:
  port: 8080
`
	configFile, cleanup := createTempConfigFile(t, yamlContent, "yaml")
	defer cleanup()

	var loadedCfg testAppConfig
	initializeTestConfig(&loadedCfg) // Initialize before loading

	// Execute & Assert: LoadConfig should return an error for unsupported config type
	err := LoadConfig(&loadedCfg, WithConfigFile(configFile, "unsupported"))
	require.Error(t, err, "LoadConfig should return error for unsupported config type")
	// Viper returns a specific error type in this case
	assert.Contains(t, err.Error(), "Unsupported Config Type")
}

// --- TODO: Add more test cases ---
// (All basic cases covered for now)
