/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Tests for config hot reload and callback functionality.
 */

package config

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Note: testAppConfig and helper functions like createTempConfigFile, initializeTestConfig
// are assumed to be defined in config_test.go or a shared test utility file.

func TestConfigHotReload_Callback(t *testing.T) {
	// Initial config content
	initialContent := `
log:
  level: "info"
server:
  port: 8080
customFeature:
  rateLimit: 100
`
	// Content after update
	updatedContent := `
log:
  level: "debug" # Updated level
server:
  port: 9090 # Updated port
customFeature:
  rateLimit: 200 # Updated rate limit
`

	configFile, cleanup := createTempConfigFile(t, initialContent, "yaml")
	defer cleanup()

	var loadedCfg testAppConfig
	initializeTestConfig(&loadedCfg) // Ensure pointers are not nil

	// Use LoadConfigAndWatch with hot reload enabled
	cm, err := LoadConfigAndWatch(&loadedCfg, WithConfigFile(configFile, "yaml"), WithHotReload(true))
	require.NoError(t, err, "LoadConfigAndWatch should succeed")
	require.NotNil(t, cm, "ConfigManager should be returned")

	// --- Verify Initial Load --- 
	require.NotNil(t, loadedCfg.Log)
	require.NotNil(t, loadedCfg.Server)
	require.NotNil(t, loadedCfg.CustomFeature)
	assert.Equal(t, "info", loadedCfg.Log.Level, "Initial log level should be info")
	assert.Equal(t, 8080, loadedCfg.Server.Port, "Initial server port should be 8080")
	assert.Equal(t, 100, loadedCfg.CustomFeature.RateLimit, "Initial rate limit should be 100")

	// --- Setup Callback --- 
	var callbackExecuted atomic.Bool // Use atomic for thread safety
	callbackChan := make(chan bool, 1)   // Channel to signal callback completion
	var callbackErr error
	var callbackMutex sync.Mutex

	cm.RegisterCallback(func(v *viper.Viper, cfg any) error {
		t.Log("Callback function executed!")
		callbackExecuted.Store(true)

		// Type assert the passed config object
		appCfg, ok := cfg.(*testAppConfig)
		if !ok {
			callbackMutex.Lock()
			callbackErr = fmt.Errorf("callback received unexpected config type: %T", cfg)
			callbackMutex.Unlock()
			callbackChan <- true // Signal completion even on error
			return callbackErr    // Return error from callback
		}

		// Verify values *inside* the callback
		callbackMutex.Lock()
		if v.GetString("log.level") != "debug" {
			callbackErr = fmt.Errorf("viper inside callback has wrong log level: %s", v.GetString("log.level"))
		} else if appCfg.Log.Level != "debug" {
			callbackErr = fmt.Errorf("cfg inside callback has wrong log level: %s", appCfg.Log.Level)
		} else if appCfg.Server.Port != 9090 {
			callbackErr = fmt.Errorf("cfg inside callback has wrong server port: %d", appCfg.Server.Port)
		} else if appCfg.CustomFeature.RateLimit != 200 {
			callbackErr = fmt.Errorf("cfg inside callback has wrong rate limit: %d", appCfg.CustomFeature.RateLimit)
		}
		callbackMutex.Unlock()
		
		callbackChan <- true // Signal completion
		return callbackErr
	})

	// --- Trigger Hot Reload --- 
	// Wait briefly to ensure the watcher is definitely set up
	time.Sleep(100 * time.Millisecond)

	// Modify the config file
	t.Logf("Writing updated content to %s", configFile)
	err = os.WriteFile(configFile, []byte(updatedContent), 0644)
	require.NoError(t, err, "Failed to write updated config file content")

	// --- Wait for Callback Execution --- 
	// Wait for the callback to signal completion or timeout
	select {
	case <-callbackChan:
		t.Log("Callback signal received.")
		// Check for errors recorded by the callback
		callbackMutex.Lock()
		assert.NoError(t, callbackErr, "Callback function reported an error")
		callbackMutex.Unlock()
	case <-time.After(5 * time.Second): // Generous timeout for file system events
		t.Fatal("Timeout waiting for config change callback to execute")
	}

	// --- Final Assertions --- 
	// Verify the callback flag was set
	assert.True(t, callbackExecuted.Load(), "Callback function should have been executed")

	// Verify the original config struct pointer was updated
	// Add a small delay to ensure the main goroutine sees the update after the callback finishes.
	// Although the callback signals, the update to loadedCfg might have a slight delay.
	time.Sleep(50 * time.Millisecond)
	require.NotNil(t, loadedCfg.Log)
	require.NotNil(t, loadedCfg.Server)
	require.NotNil(t, loadedCfg.CustomFeature)
	assert.Equal(t, "debug", loadedCfg.Log.Level, "Log level should be updated to debug after reload")
	assert.Equal(t, 9090, loadedCfg.Server.Port, "Server port should be updated to 9090 after reload")
	assert.Equal(t, 200, loadedCfg.CustomFeature.RateLimit, "Rate limit should be updated to 200 after reload")
}

func TestConfigHotReload_NoFile(t *testing.T) {
	var loadedCfg testAppConfig
	initializeTestConfig(&loadedCfg)

	// Load without a config file but with hot reload enabled
	// Expect a warning log, but no error and no panic
	_, err := LoadConfigAndWatch(&loadedCfg, WithHotReload(true))
	require.NoError(t, err, "LoadConfigAndWatch should succeed even if hot reload enabled with no file")

	// No easy way to assert the warning log without capturing logs, so we mainly check for absence of error/panic.
}

func TestConfigHotReload_FileDeleted(t *testing.T) {
	initialContent := `log:
  level: "warn"
`
	configFile, cleanup := createTempConfigFile(t, initialContent, "yaml")
	defer cleanup()

	var loadedCfg testAppConfig
	initializeTestConfig(&loadedCfg)

	cm, err := LoadConfigAndWatch(&loadedCfg, WithConfigFile(configFile, "yaml"), WithHotReload(true))
	require.NoError(t, err)
	require.NotNil(t, cm)
	assert.Equal(t, "warn", loadedCfg.Log.Level, "Initial log level should be warn")

	// Register a callback that should NOT be called
	callbackShouldNotRunChan := make(chan bool, 1)
	cm.RegisterCallback(func(v *viper.Viper, cfg any) error {
		t.Log("Unexpected callback execution after file deletion!")
		callbackShouldNotRunChan <- true
		return fmt.Errorf("callback should not have been called after file deletion")
	})

	// Delete the config file
	time.Sleep(100 * time.Millisecond) // Ensure watcher is running
	err = os.Remove(configFile)
	require.NoError(t, err, "Failed to delete config file")
	t.Logf("Deleted config file: %s", configFile)

	// Wait a bit to see if the callback gets triggered (it shouldn't)
	select {
	case <-callbackShouldNotRunChan:
		t.Fatal("Callback was executed after config file was deleted")
	case <-time.After(2 * time.Second): // Wait to ensure callback is not called
		t.Log("Callback correctly not executed after file deletion.")
	}

	// Verify the config struct retains the OLD value
	assert.Equal(t, "warn", loadedCfg.Log.Level, "Log level should remain 'warn' after config file deletion")
	// Viper/fsnotify might log an error internally when ReadInConfig fails, but LoadConfigAndWatch itself doesn't return it post-initial load.
}

func TestConfigCallback_Registration(t *testing.T) {
	var loadedCfg testAppConfig // Dummy config object
	cm := newConfigManager(&loadedCfg)
	require.NotNil(t, cm)

	var cb1Count, cb2Count atomic.Int32

	cb1 := func(v *viper.Viper, cfg any) error {
		cb1Count.Add(1)
		return nil
	}
	cb2 := func(v *viper.Viper, cfg any) error {
		cb2Count.Add(1)
		return fmt.Errorf("callback 2 error") // Simulate an error
	}

	cm.RegisterCallback(cb1)
	cm.RegisterCallback(cb2)
	cm.RegisterCallback(cb1) // Register cb1 again

	// Manually trigger notification to test callback execution
	cm.notifyCallbacks() // This function logs errors internally but does not return them.

	// Check if callbacks were called the correct number of times
	assert.Equal(t, int32(2), cb1Count.Load(), "Callback 1 should have been called twice")
	assert.Equal(t, int32(1), cb2Count.Load(), "Callback 2 should have been called once")
	// We can't easily assert the logged aggregate error without log capturing.
	// The primary check here is that callbacks are invoked as expected.
} 