/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Tests for environment variable override and default value handling.
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

// Note: testAppConfig and helper functions like createTempConfigFile, initializeTestConfig
// are assumed to be defined in config_test.go or a shared test utility file.

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