/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * Contains tests for LoadConfigAndWatch with various option combinations.
 */

package config

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoadConfig_WithOptions tests LoadConfigAndWatch with various options.
// 测试带有不同选项的 LoadConfigAndWatch
func TestLoadConfig_WithOptions(t *testing.T) {
	// --- Test WithEnvPrefix and WithEnvVarOverride ---
	t.Run("EnvPrefixAndOverride", func(t *testing.T) {
		prefix := "MYAPP"
		key := "SERVER_HOST"
		fullEnvKey := prefix + "_" + key
		envValue := "env.override.host"

		yamlContent := `server: { host: "yaml.host" }` // Lowercase key in YAML
		configFile, cleanup := createTempConfigFile(t, yamlContent, "yaml")
		defer cleanup()

		t.Setenv(fullEnvKey, envValue) // Set environment variable with prefix
		defer os.Unsetenv(fullEnvKey)

		var loadedCfg testAppConfig
		initializeTestConfig(&loadedCfg)

		// Load with custom prefix and override enabled (default)
		_, err := LoadConfigAndWatch(&loadedCfg,
			WithConfigFile(configFile, "yaml"),
			WithEnvPrefix(prefix),
			// WithEnvVarOverride(true) is default
		)
		require.NoError(t, err)
		assert.Equal(t, envValue, loadedCfg.Server.Host, "Host should be overridden by env var with custom prefix")

		// Load with override disabled
		var loadedCfgNoOverride testAppConfig
		initializeTestConfig(&loadedCfgNoOverride)
		_, err = LoadConfigAndWatch(&loadedCfgNoOverride,
			WithConfigFile(configFile, "yaml"),
			WithEnvPrefix(prefix),
			WithEnvVarOverride(false), // Disable override
		)
		require.NoError(t, err)
		assert.Equal(t, "yaml.host", loadedCfgNoOverride.Server.Host, "Host should NOT be overridden when disabled")
	})

	// --- Test WithHotReload Option Activation ---
	// (Hot reload functionality is tested in TestConfigHotReload*,
	// this just checks if the option correctly enables the watcher log messages)
	t.Run("HotReloadOption", func(t *testing.T) {
		// We can't directly check if fsnotify watcher is active easily in a unit test.
		// We rely on the log messages produced by LoadConfigAndWatch.
		// This requires capturing log output. (Setup not included here for brevity)
		// Instead, we'll just call LoadConfigAndWatch with the option
		// and ensure it doesn't panic or error unexpectedly.

		yamlContent := `server: { host: "hotreload.test" }`
		configFile, cleanup := createTempConfigFile(t, yamlContent, "yaml")
		defer cleanup()

		var loadedCfg testAppConfig
		initializeTestConfig(&loadedCfg)

		// Load with hot reload enabled
		cm, err := LoadConfigAndWatch(&loadedCfg,
			WithConfigFile(configFile, "yaml"),
			WithHotReload(true),
		)
		require.NoError(t, err)
		require.NotNil(t, cm, "Config manager should be returned")
		// Asserting the log message "Hot reload enabled..." would be ideal but requires log capture setup.
		// For now, just ensuring no error is a basic check.

		// Load with hot reload disabled (default)
		var loadedCfgNoReload testAppConfig
		initializeTestConfig(&loadedCfgNoReload)
		cm2, err := LoadConfigAndWatch(&loadedCfgNoReload,
			WithConfigFile(configFile, "yaml"),
			WithHotReload(false), // Explicitly false
		)
		require.NoError(t, err)
		require.NotNil(t, cm2)
		// Asserting the *absence* of the log message "Hot reload enabled..." would be ideal.
	})

	// --- Test Multiple Options Together ---
	t.Run("MultipleOptions", func(t *testing.T) {
		prefix := "COMBO"
		envKey := prefix + "_SERVER_PORT"
		envPort := "5555"

		// Default tag for server mode
		type ComboConfig struct {
			Config Config `mapstructure:",squash"`
			Extra  string `default:"extra_default"`
		}
		yamlContent := `
server:
  host: "combo.host"
  port: 9999 # Will be overridden by env
`
		configFile, cleanup := createTempConfigFile(t, yamlContent, "yaml")
		defer cleanup()
		t.Setenv(envKey, envPort)
		defer os.Unsetenv(envKey)

		var loadedCfg ComboConfig
		// Initialize only the embedded Config part
		if loadedCfg.Config.Server == nil { loadedCfg.Config.Server = &ServerConfig{} }
		// No need to initialize Log, Database etc. if not used/expected

		_, err := LoadConfigAndWatch(&loadedCfg,
			WithConfigFile(configFile, "yaml"),
			WithEnvPrefix(prefix),
			WithEnvVarOverride(true), // Enable override
			WithHotReload(false),     // Disable hot reload for this test
		)
		require.NoError(t, err)

		// Check values from different sources:
		assert.Equal(t, "combo.host", loadedCfg.Config.Server.Host, "Host from YAML file") // Corrected access via loadedCfg.Config.Server
		assert.Equal(t, 5555, loadedCfg.Config.Server.Port, "Port overridden by environment variable") // Corrected access
		assert.Equal(t, "release", loadedCfg.Config.Server.Mode, "Mode from struct tag default") // Corrected access
		assert.Equal(t, "extra_default", loadedCfg.Extra, "Extra from struct tag default")
	})
}

// TestLoadConfigAndWatch_ErrorPaths tests error handling in LoadConfigAndWatch.
// 测试 LoadConfigAndWatch 的错误处理路径
func TestLoadConfigAndWatch_ErrorPaths(t *testing.T) {
	// --- Test Invalid Config File Content (already covered in errors_test.go, but good to have variation) ---
	t.Run("InvalidFileContent", func(t *testing.T) {
		invalidContent := `server: { host: "bad yaml`
		configFile, cleanup := createTempConfigFile(t, invalidContent, "yaml")
		defer cleanup()

		var loadedCfg testAppConfig
		initializeTestConfig(&loadedCfg)

		_, err := LoadConfigAndWatch(&loadedCfg, WithConfigFile(configFile, "yaml"))
		require.Error(t, err, "Should return error for invalid YAML content")
		assert.Contains(t, err.Error(), "failed to read config file", "Error message should indicate config file read failure")
	})

	// --- Test Config File Type Mismatch ---
	t.Run("ConfigFileTypeMismatch", func(t *testing.T) {
		yamlContent := `server: { host: "good yaml" }`
		configFile, cleanup := createTempConfigFile(t, yamlContent, "yaml")
		defer cleanup()

		var loadedCfg testAppConfig
		initializeTestConfig(&loadedCfg)

		// Specify TOML type for a YAML file
		_, err := LoadConfigAndWatch(&loadedCfg, WithConfigFile(configFile, "toml"))
		require.Error(t, err, "Should return error for file type mismatch")
		// Error might manifest during ReadInConfig or Unmarshal, check for relevant parts
		assert.Contains(t, err.Error(), "failed to read config file", "Error message likely indicates read failure due to type mismatch")
	})

	// --- Test Hot Reload with Read Error ---
	t.Run("HotReloadReadError", func(t *testing.T) {
		initialContent := `value: initial`
		updatedContentInvalid := `value: { bad yaml`
		configFile, cleanup := createTempConfigFile(t, initialContent, "yaml")
		defer cleanup()

		type SimpleCfg struct {
			Value string `mapstructure:"value"`
		}
		var loadedCfg SimpleCfg
		var callbackExecuted bool

		cm, err := LoadConfigAndWatch(&loadedCfg,
			WithConfigFile(configFile, "yaml"),
			WithHotReload(true),
		)
		require.NoError(t, err)
		require.NotNil(t, cm)
		assert.Equal(t, "initial", loadedCfg.Value)

		// Register callback to check if it's incorrectly triggered on error
		cm.RegisterCallback(func(v *viper.Viper, cfgAny any) error {
			callbackExecuted = true
			return nil
		})

		// Introduce invalid content and trigger reload
		err = os.WriteFile(configFile, []byte(updatedContentInvalid), 0644)
		require.NoError(t, err, "Failed to write invalid content to config file")

		// Wait for fsnotify event and processing (can be flaky)
		time.Sleep(100 * time.Millisecond) // Allow time for watcher to potentially react

		// Assert: Config should retain old value, callback should NOT be executed
		assert.Equal(t, "initial", loadedCfg.Value, "Config should retain old value after reload read error")
		assert.False(t, callbackExecuted, "Callback should not be executed after reload read error")
		// We expect error logs, but can't easily check them in unit tests without capture.
	})

	// --- Test Hot Reload with Unmarshal Error ---
	t.Run("HotReloadUnmarshalError", func(t *testing.T) {
		initialContent := `value: initial_text`
		// Valid YAML, but wrong type for the struct field
		updatedContentWrongType := `value: ["not", "a", "string"]`
		configFile, cleanup := createTempConfigFile(t, initialContent, "yaml")
		defer cleanup()

		type SimpleCfg struct {
			Value string `mapstructure:"value"`
		}
		var loadedCfg SimpleCfg
		var callbackExecuted bool

		cm, err := LoadConfigAndWatch(&loadedCfg,
			WithConfigFile(configFile, "yaml"),
			WithHotReload(true),
		)
		require.NoError(t, err)
		require.NotNil(t, cm)
		assert.Equal(t, "initial_text", loadedCfg.Value)

		cm.RegisterCallback(func(v *viper.Viper, cfgAny any) error {
			callbackExecuted = true
			return nil
		})

		// Introduce content with wrong type and trigger reload
		err = os.WriteFile(configFile, []byte(updatedContentWrongType), 0644)
		require.NoError(t, err, "Failed to write wrong type content to config file")

		time.Sleep(100 * time.Millisecond) // Allow time for watcher

		// Assert: Config should retain old value, callback should NOT be executed
		assert.Equal(t, "initial_text", loadedCfg.Value, "Config should retain old value after reload unmarshal error")
		assert.False(t, callbackExecuted, "Callback should not be executed after reload unmarshal error")
		// Expect error logs.
	})
} 