/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: 服务器配置管理集成测试 (Server configuration management integration tests)
 */

package integration

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	// For *gin.Context

	"github.com/stretchr/testify/require"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	// "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services" // No longer needed

	_ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin" // Import Gin for server creation
)

// TestConfigurationManagement 测试不同的服务器配置及其应用
// (TestConfigurationManagement tests different server configurations and their application)
func TestConfigurationManagement(t *testing.T) {
	t.Parallel()

	defaultServerTimeouts := server.DefaultServerConfig() // Added to get default values

	testCases := []struct {
		name                string
		configMutator       func(cfg *server.ServerConfig)
		expectedHost        string
		expectCORSEnabled   bool // Used to verify CORS.Enabled setting
		expectReadTimeout   time.Duration
		expectWriteTimeout  time.Duration
		checkCORSHeaders    bool // Flag to perform OPTIONS request for CORS check
	}{
		{
			name: "DefaultConfig",
			configMutator: func(cfg *server.ServerConfig) {},
			expectedHost:        "127.0.0.1",
			expectCORSEnabled:   true, // Assuming default is true from createTestServerConfig
			expectReadTimeout:   defaultServerTimeouts.ReadTimeout, // Corrected
			expectWriteTimeout:  defaultServerTimeouts.WriteTimeout, // Corrected
			checkCORSHeaders:    false, // No specific CORS check for default beyond Enabled status
		},
		{
			name: "CustomHost",
			configMutator: func(cfg *server.ServerConfig) {
				cfg.Host = "localhost"
			},
			expectedHost:        "localhost",
			expectCORSEnabled:   true,
			expectReadTimeout:   defaultServerTimeouts.ReadTimeout, // Corrected
			expectWriteTimeout:  defaultServerTimeouts.WriteTimeout, // Corrected
			checkCORSHeaders:    false,
		},
		{
			name: "DisableCORS",
			configMutator: func(cfg *server.ServerConfig) {
				cfg.CORS.Enabled = false
			},
			expectedHost:        "127.0.0.1",
			expectCORSEnabled:   false,
			expectReadTimeout:   defaultServerTimeouts.ReadTimeout, // Corrected
			expectWriteTimeout:  defaultServerTimeouts.WriteTimeout, // Corrected
			checkCORSHeaders:    true,
		},
		{
			name: "DifferentTimeouts",
			configMutator: func(cfg *server.ServerConfig) {
				cfg.ReadTimeout = 15 * time.Second
				cfg.WriteTimeout = 25 * time.Second
			},
			expectedHost:        "127.0.0.1",
			expectCORSEnabled:   true,
			expectReadTimeout:   15 * time.Second,
			expectWriteTimeout:  25 * time.Second,
			checkCORSHeaders:    false,
		},
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			
			// 创建基础配置并获取动态端口 (Create base config and get dynamic port)
			baseCfg := createTestServerConfig(0, "gin")
			
			// 获取动态端口 (Get dynamic port)
			listener, err := net.Listen("tcp", fmt.Sprintf("%s:0", baseCfg.Host))
			require.NoError(t, err, "Failed to listen for dynamic port for %s", tc.name)
			dynamicPort := listener.Addr().(*net.TCPAddr).Port
			require.NoError(t, listener.Close(), "Failed to close listener for %s", tc.name)
			baseCfg.Port = dynamicPort
			
			// 应用测试特定的配置变更 (Apply test-specific config mutations)
			tc.configMutator(baseCfg)
			
			// 创建服务器管理器 (Create server manager)
			s, err := server.CreateServerManager(baseCfg.Framework, baseCfg)
			require.NoError(t, err, "Failed to create server manager for %s", tc.name)
			require.NotNil(t, s, "Server manager instance should not be nil for %s", tc.name)
			
			// 验证实际使用的配置 (Verify actual config used by the server manager)
			actualCfg := s.GetConfig()
			require.NotNil(t, actualCfg, "Actual config from server manager should not be nil for %s", tc.name)
			require.Equal(t, tc.expectedHost, actualCfg.Host, "Host mismatch in actual config for %s", tc.name)
			require.Equal(t, dynamicPort, actualCfg.Port, "Port mismatch in actual config for %s", tc.name)
			require.Equal(t, tc.expectCORSEnabled, actualCfg.CORS.Enabled, "CORS.Enabled mismatch for %s", tc.name)
			require.Equal(t, tc.expectReadTimeout, actualCfg.ReadTimeout, "ReadTimeout mismatch for %s", tc.name)
			require.Equal(t, tc.expectWriteTimeout, actualCfg.WriteTimeout, "WriteTimeout mismatch for %s", tc.name)
			
			// 获取底层框架并设置测试路由 (Get underlying framework and set up test routes)
			webFramework := s.GetFramework()
			require.NotNil(t, webFramework, "WebFramework should not be nil for %s", tc.name)
			
			// 设置测试路由 (Set up test route)
			err = webFramework.RegisterRoute("GET", "/config-test-ping", server.HandlerFunc(func(ctx server.Context) error {
				return ctx.String(200, "pong")
			}))
			require.NoError(t, err, "Failed to set up test route for %s", tc.name)
			
			// 在单独的goroutine中启动服务器 (Start server in separate goroutine)
			serverErrChan := make(chan error, 1)
			go func() {
				// s.Start() 会阻塞直到服务器关闭 (s.Start() blocks until server shuts down)
				serverErrChan <- s.Start(context.Background())
			}()

			serverAddr := fmt.Sprintf("%s:%d", actualCfg.Host, actualCfg.Port)
			
			// 等待服务器启动 (Wait for server to start)
			require.Eventually(t, func() bool {
				conn, dialErr := net.DialTimeout("tcp", serverAddr, 100*time.Millisecond)
				if dialErr == nil {
					conn.Close()
					return true
				}
				return false
			}, 5*time.Second, 100*time.Millisecond, "Server for %s did not start listening on %s", tc.name, serverAddr)

			// 发送测试请求 (Send test request)
			client := http.Client{Timeout: 3 * time.Second}
			reqURL := fmt.Sprintf("http://%s/config-test-ping", serverAddr)
			resp, err := client.Get(reqURL)
			require.NoError(t, err, "Request to %s failed", reqURL)
			defer resp.Body.Close()
			
			require.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code for %s", tc.name)
			
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err, "Failed to read response body for %s", tc.name)
			require.Contains(t, string(body), "pong", "Unexpected response body for %s: %s", tc.name, string(body))
			
			// 检查CORS头部（如果需要）(Check CORS headers if needed)
			if tc.checkCORSHeaders {
				req, _ := http.NewRequest("OPTIONS", reqURL, nil)
				req.Header.Set("Origin", "http://example.com")
				req.Header.Set("Access-Control-Request-Method", "GET")
				corsResp, corsErr := client.Do(req)
				require.NoError(t, corsErr, "CORS OPTIONS request failed for %s", tc.name)
				defer corsResp.Body.Close()
				
				if !tc.expectCORSEnabled {
					require.Empty(t, corsResp.Header.Get("Access-Control-Allow-Origin"), 
						"Access-Control-Allow-Origin should be empty if CORS is disabled for %s", tc.name)
				}
			}
			
			// 停止服务器 (Stop server)
			stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			
			stopErr := s.Stop(stopCtx)
			require.NoError(t, stopErr, "Failed to stop server for %s", tc.name)
			
			// 等待服务器停止 (Wait for server to stop)
			select {
			case startErr := <-serverErrChan:
				// 服务器正常停止，应该返回http.ErrServerClosed或nil (Server stopped normally, should return http.ErrServerClosed or nil)
				if startErr != nil && !errors.Is(startErr, http.ErrServerClosed) {
					t.Fatalf("Server returned unexpected error for %s: %v", tc.name, startErr)
				}
			case <-time.After(5 * time.Second):
				t.Fatalf("Server did not stop within timeout for %s", tc.name)
			}
		})
	}
}
