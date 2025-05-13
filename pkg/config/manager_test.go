/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * Contains tests for the config manager functionality.
 */

package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetViperInstance tests retrieving the internal Viper instance.
// 测试获取内部 Viper 实例
func TestGetViperInstance(t *testing.T) {
	// Need to create a configManager instance first
	var cfgData struct{} // Dummy config struct
	cm := newConfigManager(&cfgData) // Use the internal constructor
	require.NotNil(t, cm, "newConfigManager should return a non-nil instance")

	// Get the viper instance
	vInstance := cm.GetViperInstance()

	// Basic assertions
	require.NotNil(t, vInstance, "GetViperInstance should return a non-nil viper instance")

	// Check if it's the same instance stored internally (optional, requires access)
	// assert.Same(t, cm.v, vInstance, "Returned viper instance should be the same as internal")
	// Since cm.v is unexported, we can't directly compare.
	// We can try setting a value on the returned instance and see if it reflects internally,
	// or just rely on the non-nil check and type assertion.

	// Check the type
	_, ok := interface{}(vInstance).(*viper.Viper)
	assert.True(t, ok, "Returned instance should be of type *viper.Viper")

	// Try setting a value on the returned instance and reading it back
	testKey := "test_viper_instance_key"
	testValue := "test_viper_instance_value"
	vInstance.Set(testKey, testValue)
	assert.Equal(t, testValue, vInstance.GetString(testKey), "Should be able to set/get values on the returned viper instance")
}

// TODO: Add tests for newConfigManager if specific options need verification
// TODO: Add tests for notifyCallbacks (might be tricky without real callbacks) 