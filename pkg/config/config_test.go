/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Common setup and helper functions for config package tests.
 */

package config

import (
	"os"
	"path/filepath"
	"testing"

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
// or structs embedding it, like testAppConfig.
// (Helper function to initialize pointer fields in the embedded Config or structs embedding it)
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
	// CustomFeature is already a pointer in testAppConfig, let initializeNilPointers handle it if needed elsewhere
}

// Note: Actual test functions like TestLoadConfig_YAML_Valid, TestAccessors, etc.,
// have been moved to separate files like load_basic_test.go, accessors_test.go, etc.
// This file primarily contains shared setup code.
