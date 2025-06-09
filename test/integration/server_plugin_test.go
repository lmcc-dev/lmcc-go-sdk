/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: 服务器插件集成测试 (Server plugin integration tests)
 */

package integration

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	// Import Gin, Echo, Fiber plugins to ensure they are registered
	echoPlugin "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/echo"
	fiberPlugin "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/fiber"
	ginPlugin "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin"
)

func init() {
	// 显式注册插件以确保它们在测试环境中可用 (Explicitly register plugins to ensure they are available in test environment)
	_ = server.RegisterFramework(ginPlugin.NewPlugin())
	_ = server.RegisterFramework(echoPlugin.NewPlugin())
	_ = server.RegisterFramework(fiberPlugin.NewPlugin())
}

// TestServerPluginRegistration 测试服务器插件的注册和基本信息获取
// (TestServerPluginRegistration tests server plugin registration and basic information retrieval)
func TestServerPluginRegistration(t *testing.T) {
	t.Parallel()

	frameworks := []string{"gin", "echo", "fiber"}

	for _, frameworkName := range frameworks {
		frameworkName := frameworkName // Capture range variable
		t.Run(fmt.Sprintf("Framework_%s", frameworkName), func(t *testing.T) {
			t.Parallel()

			listener, err := net.Listen("tcp", "127.0.0.1:0")
			require.NoError(t, err, "Failed to listen on a free port for %s", frameworkName)
			dynamicPort := listener.Addr().(*net.TCPAddr).Port
			require.NoError(t, listener.Close(), "Failed to close listener for %s", frameworkName)

			cfg := createTestServerConfig(dynamicPort, frameworkName)

			s, err := server.CreateServerManager(frameworkName, cfg)
			require.NoError(t, err, "Failed to create server manager for %s", frameworkName)
			require.NotNil(t, s, "Server manager instance should not be nil for %s", frameworkName)

			plugin := s.GetFramework()
			require.NotNil(t, plugin, "Framework plugin for %s should not be nil", frameworkName)
			// The WebFramework interface does not have a Name() method.
			// The successful creation and retrieval of the framework is an implicit check.
			// assert.Equal(t, frameworkName, plugin.Name(), "Plugin name mismatch for %s", frameworkName) // Removed

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			serverErrChan := make(chan error, 1)
			go func() {
				errStart := s.Start(context.Background())
				if errStart != nil {
					serverErrChan <- errStart
				}
				close(serverErrChan)
			}()

			serverAddr := fmt.Sprintf("127.0.0.1:%d", dynamicPort)
			require.Eventually(t, func() bool {
				conn, errDial := net.DialTimeout("tcp", serverAddr, 100*time.Millisecond)
				if errDial == nil {
					conn.Close()
					return true
				}
				return false
			}, 5*time.Second, 100*time.Millisecond, "Server for %s did not start listening on %s", frameworkName, serverAddr)

			errStop := s.Stop(ctx)
			assert.NoError(t, errStop, "Failed to stop server %s gracefully: %v", frameworkName, errStop)

			// 等待服务器停止 (Wait for server to stop)
			select {
			case startErr := <-serverErrChan:
				// 服务器应该返回包装的 http.ErrServerClosed (Server should return wrapped http.ErrServerClosed)
				if startErr != nil && !errors.Is(startErr, http.ErrServerClosed) {
					t.Errorf("s.Start() for %s expected ErrServerClosed, got %v", frameworkName, startErr)
				} else if frameworkName == "fiber" && startErr == nil {
					t.Logf("s.Start() for %s returned cleanly, or channel closed before error could be sent. Assuming successful shutdown.", frameworkName)
				}
			case <-time.After(10 * time.Second):
				t.Errorf("s.Start() for %s did not return after stop within timeout", frameworkName)
			}
		})
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	dynamicPort := listener.Addr().(*net.TCPAddr).Port
	require.NoError(t, listener.Close())

	cfg := createTestServerConfig(dynamicPort, "gin")

	// Test creating a server with a non-existent framework.
	// The 's' variable from previous successful creation with "gin" is not needed here.
	// s, err := server.CreateServerManager("gin", cfg) // Removed
	// require.NoError(t, err) // Removed

	_, err = server.CreateServerManager("nonexistentframework", cfg)
	assert.Error(t, err, "Expected an error when creating server with a non-existent framework")
	if err != nil {
		assert.Contains(t, err.Error(), "is not registered", "Error message should indicate plugin not found")
	}
}
