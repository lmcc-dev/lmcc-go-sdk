/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: 服务器服务容器集成测试 (Server service container integration tests)
 */

package integration

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	lmccerrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	lmcclog "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services" // Keep for NewConfigManagerImpl etc.

	_ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin" // Import Gin for default server creation if needed
)

// TestServiceContainerIntegration 测试服务容器 (日志、配置、错误处理) 的集成
// (TestServiceContainerIntegration tests the integration of the service container (logging, config, error handling))
func TestServiceContainerIntegration(t *testing.T) {
	t.Parallel()

	// 1. 日志测试 (Logging Test) - 使用全局日志替换
	var logOutput bytes.Buffer
	customLogOpts := lmcclog.NewOptions()
	customLogOpts.OutputPaths = []string{"stdout"} // Will be replaced by writer
	customLogOpts.Level = "debug"
	customLogOpts.Format = lmcclog.FormatText
	// Corrected: NewLoggerWithWriter returns only Logger, not (Logger, error)
	customLogger := lmcclog.NewLoggerWithWriter(customLogOpts, &logOutput)
	require.NotNil(t, customLogger, "Custom logger should not be nil")

	// 替换全局 logger (Replace global logger)
	// Corrected: Use Std() to get current global logger and SetGlobalLogger to replace
	originalGlobalLogger := lmcclog.Std()
	lmcclog.SetGlobalLogger(customLogger)
	defer lmcclog.SetGlobalLogger(originalGlobalLogger) // Restore original logger after test

	// 获取动态端口 (Get dynamic port)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err, "Failed to listen on a free port")
	dynamicPort := listener.Addr().(*net.TCPAddr).Port
	require.NoError(t, listener.Close(), "Failed to close listener")

	cfg := createTestServerConfig(dynamicPort, "gin")

	s, err := server.CreateServerManager(cfg.Framework, cfg)
	require.NoError(t, err, "Failed to create server manager")
	require.NotNil(t, s, "Server manager instance should not be nil")

	// 获取 Gin 引擎并注册路由 (Get Gin engine and register routes)
	webFramework := s.GetFramework()
	require.NotNil(t, webFramework, "WebFramework should not be nil")
	engine, ok := webFramework.GetNativeEngine().(*gin.Engine)
	require.True(t, ok, "Failed to cast to *gin.Engine")

	// 日志测试路由 (Log test route)
	engine.GET("/test-log", func(c *gin.Context) {
		// Corrected: Use Std() to get the global logger
		lmcclog.Std().Info("message from test log handler")
		c.String(http.StatusOK, "log test ok")
	})

	// 错误处理测试路由 (Error handling test route)
	engine.GET("/test-error", func(c *gin.Context) {
		coder := lmccerrors.ErrNotFound // Example error code
		c.JSON(coder.HTTPStatus(), gin.H{"error_code": coder.Code(), "message": coder.String()})
	})

	// 启动服务器 (Start the server)
	startCtx, startCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer startCancel()
	serverErrChan := make(chan error, 1)
	go func() {
		errStart := s.Start(context.Background()) // Use a background context for server run
		if errStart != nil {
			serverErrChan <- errStart
		}
		close(serverErrChan)
	}()

	serverAddr := fmt.Sprintf("127.0.0.1:%d", dynamicPort)
	require.Eventually(t, func() bool {
		conn, dialErr := net.DialTimeout("tcp", serverAddr, 100*time.Millisecond)
		if dialErr == nil {
			conn.Close()
			return true
		}
		return false
	}, 5*time.Second, 100*time.Millisecond, "Server did not start listening on %s", serverAddr)

	// 执行日志测试请求 (Execute log test request)
	client := http.Client{Timeout: 3 * time.Second}
	respLog, err := client.Get(fmt.Sprintf("http://%s/test-log", serverAddr))
	require.NoError(t, err, "Failed to make GET request to /test-log")
	if respLog != nil && respLog.Body != nil {
		assert.Equal(t, http.StatusOK, respLog.StatusCode)
		respLog.Body.Close()
	}
	assert.Contains(t, logOutput.String(), "message from test log handler", "Log output mismatch")

	// 执行错误处理测试请求 (Execute error handling test request)
	respErr, err := client.Get(fmt.Sprintf("http://%s/test-error", serverAddr))
	require.NoError(t, err, "Failed to make GET request to /test-error")
	if respErr != nil && respErr.Body != nil {
		bodyBytes, _ := io.ReadAll(respErr.Body)
		assert.Equal(t, lmccerrors.ErrNotFound.HTTPStatus(), respErr.StatusCode, "Status code for /test-error mismatch")
		assert.Contains(t, string(bodyBytes), lmccerrors.ErrNotFound.String(), "Error message in response mismatch")
		respErr.Body.Close()
	}

	// 测试服务器配置访问 (Test server config access - basic)
	retrievedCfg := s.GetConfig()
	require.NotNil(t, retrievedCfg, "Retrieved server config should not be nil")
	assert.Equal(t, cfg.Framework, retrievedCfg.Framework, "Framework in retrieved config mismatch")
	assert.Equal(t, dynamicPort, retrievedCfg.Port, "Port in retrieved config mismatch")

	// 优雅关闭 (Graceful shutdown)
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer stopCancel()
	errStop := s.Stop(stopCtx)
	assert.NoError(t, errStop, "Failed to stop server gracefully")

	select {
	case startErr, ok := <-serverErrChan:
		if ok && startErr != nil {
			// 服务器应该返回包装的 http.ErrServerClosed (Server should return wrapped http.ErrServerClosed)
			if !errors.Is(startErr, http.ErrServerClosed) {
				t.Errorf("s.Start() expected ErrServerClosed, got %v", startErr)
			}
		}
	case <-startCtx.Done(): // Use startCtx for overall server operation timeout including shutdown error check
		t.Errorf("Timeout waiting for s.Start() to return after stop: %v", startCtx.Err())
	}
}

// TestServiceContainer_GetConfigManager_NoUnderlyingManager tests config manager behavior when no underlying pkg/config.Manager is provided.
// (TestServiceContainer_GetConfigManager_NoUnderlyingManager 测试在未提供底层 pkg/config.Manager 时配置管理器的行为。)
func TestServiceContainer_GetConfigManager_NoUnderlyingManager(t *testing.T) {
	t.Parallel()
	cm := services.NewConfigManagerImpl(nil) // Create with nil pkg/config.Manager

	assert.Nil(t, cm.Get("somekey"), "Get should return nil when viper is nil")
	assert.Equal(t, "", cm.GetString("somekey"), "GetString should return empty string")
	assert.Equal(t, 0, cm.GetInt("somekey"), "GetInt should return 0")
	assert.False(t, cm.GetBool("somekey"), "GetBool should return false")
	assert.Equal(t, 0.0, cm.GetFloat64("somekey"), "GetFloat64 should return 0.0")
	assert.Nil(t, cm.GetStringSlice("somekey"), "GetStringSlice should return nil")
	assert.False(t, cm.IsSet("somekey"), "IsSet should return false")
	assert.Nil(t, cm.GetViperInstance(), "GetViperInstance should return nil")

	cm.Set("anotherkey", "value") // Should be a no-op and not panic
	assert.False(t, cm.IsSet("anotherkey"), "Set should not take effect")

	type DummyStruct struct{ Name string }
	var ds DummyStruct
	err := cm.Unmarshal(&ds)
	require.Error(t, err, "Unmarshal should return error")
	assert.Contains(t, err.Error(), "viper instance is nil")

	err = cm.UnmarshalKey("key", &ds)
	require.Error(t, err, "UnmarshalKey should return error")
	assert.Contains(t, err.Error(), "viper instance is nil")

	// RegisterCallback should be a no-op and not panic
	cm.RegisterCallback(func(v *viper.Viper, cfg any) error {
		t.Error("Callback should not be called when manager is nil")
		return nil
	})
}
