/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: 服务器多框架支持集成测试 (Server multi-framework support integration tests)
 */

package integration

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"

	// Import Gin, Echo, Fiber plugins to ensure they are registered
	_ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/echo"
	_ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/fiber"
	_ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin"
)

// TestMultiFrameworkSupport 测试服务器对多种Web框架 (Gin, Echo, Fiber) 的支持
// (TestMultiFrameworkSupport tests server support for multiple web frameworks: Gin, Echo, Fiber)
func TestMultiFrameworkSupport(t *testing.T) {
	t.Parallel()

	frameworks := []struct {
		name             string
		setupRoute       func(instance interface{}, t *testing.T, expectedMessage string)
		expectedResponse string
	}{
		{
			name: "gin",
			setupRoute: func(instance interface{}, t *testing.T, expectedMessage string) {
				ginEngine, ok := instance.(*gin.Engine)
				require.True(t, ok, "Failed to cast to *gin.Engine for %s", "gin")
				ginEngine.GET("/hello", func(c *gin.Context) {
					c.String(http.StatusOK, expectedMessage)
				})
			},
			expectedResponse: "Hello from Gin",
		},
		{
			name: "echo",
			setupRoute: func(instance interface{}, t *testing.T, expectedMessage string) {
				echoEngine, ok := instance.(*echo.Echo)
				require.True(t, ok, "Failed to cast to *echo.Echo for %s", "echo")
				echoEngine.GET("/hello", func(c echo.Context) error {
					return c.String(http.StatusOK, expectedMessage)
				})
			},
			expectedResponse: "Hello from Echo",
		},
		{
			name: "fiber",
			setupRoute: func(instance interface{}, t *testing.T, expectedMessage string) {
				fiberApp, ok := instance.(*fiber.App)
				require.True(t, ok, "Failed to cast to *fiber.App for %s", "fiber")
				fiberApp.Get("/hello", func(c *fiber.Ctx) error {
					return c.SendString(expectedMessage) // Fiber uses SendString
				})
			},
			expectedResponse: "Hello from Fiber",
		},
	}

	var wg sync.WaitGroup
	for _, tc := range frameworks {
		tc := tc // Capture range variable
		wg.Add(1)
		go func() {
			defer wg.Done()
			t.Run(fmt.Sprintf("Framework_%s_Support", tc.name), func(t *testing.T) {
				listener, err := net.Listen("tcp", "127.0.0.1:0")
				require.NoError(t, err, "Failed to listen on a free port for %s", tc.name)
				dynamicPort := listener.Addr().(*net.TCPAddr).Port
				require.NoError(t, listener.Close(), "Failed to close listener for %s", tc.name)

				cfg := createTestServerConfig(dynamicPort, tc.name)
				// serviceContainer is not passed to CreateServerManager
				s, err := server.CreateServerManager(tc.name, cfg)
				require.NoError(t, err, "Failed to create server manager for %s", tc.name)
				require.NotNil(t, s, "Server manager for %s should not be nil", tc.name)

				webFramework := s.GetFramework()
				require.NotNil(t, webFramework, "WebFramework for %s should not be nil", tc.name)
				engine := webFramework.GetNativeEngine() // Corrected: GetNativeEngine instead of GetEngine
				require.NotNil(t, engine, "Engine for %s should not be nil", tc.name)

				tc.setupRoute(engine, t, tc.expectedResponse)

				// 在单独的goroutine中启动服务器 (Start server in separate goroutine)
				serverErrChan := make(chan error, 1)
				go func() {
					serverErrChan <- s.Start(context.Background())
				}()

				serverAddr := fmt.Sprintf("127.0.0.1:%d", dynamicPort)
				require.Eventually(t, func() bool {
					conn, errDial := net.DialTimeout("tcp", serverAddr, 100*time.Millisecond)
					if errDial == nil {
						conn.Close()
						return true
					}
					return false
				}, 5*time.Second, 100*time.Millisecond, "Server for %s did not start listening on %s", tc.name, serverAddr)

				client := http.Client{Timeout: 3 * time.Second}
				reqURL := fmt.Sprintf("http://%s:%d/hello", cfg.Host, cfg.Port)
				resp, err := client.Get(reqURL)
				require.NoError(t, err, "Failed to make GET request to %s for %s", reqURL, tc.name)
				
				var responseBody string
				if resp != nil && resp.Body != nil {
					bodyBytes, readErr := io.ReadAll(resp.Body)
					require.NoError(t, readErr, "Failed to read response body for %s", tc.name)
					responseBody = string(bodyBytes)
					resp.Body.Close()
				}
				require.Equal(t, http.StatusOK, resp.StatusCode, "Expected status OK for %s", tc.name)
				assert.Equal(t, tc.expectedResponse, responseBody, "Response body mismatch for %s", tc.name)

				stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer stopCancel()
				errStop := s.Stop(stopCtx)
				assert.NoError(t, errStop, "Failed to stop server %s gracefully: %v", tc.name, errStop)

				select {
				case startErr := <-serverErrChan:
					if startErr != nil && !errors.Is(startErr, http.ErrServerClosed) {
						t.Errorf("s.Start() for %s expected ErrServerClosed, got %v", tc.name, startErr)
					} else if tc.name == "fiber" && startErr == nil {
						t.Logf("s.Start() for %s returned cleanly or channel closed. Assuming successful shutdown.", tc.name)
					}
				case <-time.After(10 * time.Second):
					t.Errorf("s.Start() for %s did not return after stop within timeout", tc.name)
				}
			})
		}()
	}
	wg.Wait()
}
