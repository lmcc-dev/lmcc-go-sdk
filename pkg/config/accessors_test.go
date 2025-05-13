/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * Contains tests for configuration accessor functions.
 */

package config

import (
	"testing"

	// "github.com/spf13/viper" // viper import might not be needed directly in this file anymore
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetGlobalCfg tests the GetGlobalCfg function after LoadConfig.
// 测试 GetGlobalCfg 函数
func TestGetGlobalCfg(t *testing.T) {
	originalCfg := Cfg // Save original global Cfg
	defer func() { Cfg = originalCfg }() // Restore original global Cfg

	Cfg = nil // Ensure Cfg starts nil

	yamlContent := `server: { host: "global.test" }`
	configFile, cleanup := createTempConfigFile(t, yamlContent, "yaml")
	defer cleanup()

	var loadedCfg testAppConfig
	initializeTestConfig(&loadedCfg)

	err := LoadConfig(&loadedCfg, WithConfigFile(configFile, "yaml"))
	require.NoError(t, err)

	// After LoadConfig, the global Cfg should be updated (pointing to loadedCfg.Config)
	globalCfg := GetGlobalCfg()
	require.NotNil(t, globalCfg, "Global Cfg should not be nil after LoadConfig")
	require.Equal(t, "global.test", globalCfg.Server.Host, "Global Cfg should reflect loaded value")

	// Also check if the returned pointer is the same as the embedded one
	assert.Same(t, &loadedCfg.Config, globalCfg, "GetGlobalCfg should return the same instance set by LoadConfig")

	// Test when LoadConfig fails or doesn't update Cfg (difficult to simulate perfectly)
	Cfg = nil // Reset Cfg
	// Simulate a scenario where updateGlobalCfg might fail (e.g., wrong struct type)
	var simpleStruct struct{ Name string }
	err = LoadConfig(&simpleStruct, WithConfigFile(configFile, "yaml")) // LoadConfig might error or just not update Cfg
	if err != nil {
		t.Logf("Ignoring expected LoadConfig error for simpleStruct: %v", err)
	}
	assert.Nil(t, GetGlobalCfg(), "Global Cfg should remain nil if LoadConfig fails to update it")
}

// TestUpdateGlobalCfg tests the updateGlobalCfg helper function directly.
// 测试 updateGlobalCfg 辅助函数
func TestUpdateGlobalCfg(t *testing.T) {
	originalCfg := Cfg // Save global Cfg
	defer func() { Cfg = originalCfg }() // Restore global Cfg

	// Test Case 1: Input is *config.Config
	t.Run("DirectConfigPointer", func(t *testing.T) {
		Cfg = nil // Reset global Cfg
		testCfg := &Config{Server: &ServerConfig{Host: "direct_test"}}
		updateGlobalCfg(testCfg)
		require.NotNil(t, Cfg, "Global Cfg should be updated")
		assert.Same(t, testCfg, Cfg, "Global Cfg should point to the input config")
		assert.Equal(t, "direct_test", Cfg.Server.Host)
	})

	// Test Case 2: Input embeds config.Config (value receiver)
	t.Run("EmbedConfigValue", func(t *testing.T) {
		Cfg = nil // Reset global Cfg
		type EmbedValue struct {
			Config         // Embed by value
			Custom string
		}
		testEmbed := &EmbedValue{
			Config: Config{Server: &ServerConfig{Host: "embed_value_test"}},
			Custom: "field",
		}
		initializeNilPointers(&testEmbed.Config) // Ensure nested pointers are initialized

		updateGlobalCfg(testEmbed)
		require.NotNil(t, Cfg, "Global Cfg should be updated")
		assert.Same(t, &testEmbed.Config, Cfg, "Global Cfg should point to the embedded config's address")
		assert.Equal(t, "embed_value_test", Cfg.Server.Host)
	})

	// Test Case 3: Input has pointer field *config.Config (initialized)
	t.Run("PointerFieldInitialized", func(t *testing.T) {
		Cfg = nil // Reset global Cfg
		type PtrField struct {
			CfgPtr *Config
			Other  string
		}
		innerCfg := &Config{Server: &ServerConfig{Host: "ptr_field_test"}}
		testPtr := &PtrField{
			CfgPtr: innerCfg,
			Other:  "data",
		}
		updateGlobalCfg(testPtr)
		require.NotNil(t, Cfg, "Global Cfg should be updated")
		assert.Same(t, innerCfg, Cfg, "Global Cfg should point to the inner config pointer")
		assert.Equal(t, "ptr_field_test", Cfg.Server.Host)
	})

	// Test Case 4: Input has pointer field *config.Config (nil)
	t.Run("PointerFieldNil", func(t *testing.T) {
		Cfg = &Config{Server: &ServerConfig{Host: "initial_global"}} // Set a known global Cfg
		initialGlobalAddr := Cfg // Store the address
		type PtrFieldNil struct {
			CfgPtr *Config // Nil pointer
			Other  string
		}
		testPtrNil := &PtrFieldNil{Other: "data"}

		updateGlobalCfg(testPtrNil)
		assert.NotNil(t, Cfg, "Global Cfg should NOT be updated")
		assert.Same(t, initialGlobalAddr, Cfg, "Global Cfg should remain unchanged")
		assert.Equal(t, "initial_global", Cfg.Server.Host)
	})

	// Test Case 5: Input is unrelated struct
	t.Run("UnrelatedStruct", func(t *testing.T) {
		Cfg = &Config{Server: &ServerConfig{Host: "initial_global2"}} // Set a known global Cfg
		initialGlobalAddr := Cfg // Store the address
		type Unrelated struct {
			Name string
		}
		unrelated := &Unrelated{Name: "test"}

		updateGlobalCfg(unrelated)
		assert.NotNil(t, Cfg, "Global Cfg should NOT be updated")
		assert.Same(t, initialGlobalAddr, Cfg, "Global Cfg should remain unchanged")
		assert.Equal(t, "initial_global2", Cfg.Server.Host)
	})

	// Test Case 6: Embedded field but cannot take address (simulated - hard to create in Go)
	// Go usually allows taking address of fields in addressable structs.
	// We'll skip this explicit test case as it's difficult to construct reliably.
	// The warning log in updateGlobalCfg covers this scenario conceptually.
} 