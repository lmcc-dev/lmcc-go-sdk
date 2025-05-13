/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Tests for error handling during config loading.
 */

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Note: testAppConfig and helper functions like createTempConfigFile, initializeTestConfig
// are assumed to be defined in config_test.go or a shared test utility file.

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