/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: 服务器路由和分组集成测试 (Server routing and grouping integration tests)
 */

package integration

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin" // Import Gin for gin.Context and gin.H
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"

	_ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin" // Import Gin for server creation
)

// TestRouteGrouping 测试路由分组和嵌套路由
// (TestRouteGrouping tests route grouping and nested routes)
func TestRouteGrouping(t *testing.T) {
	t.Parallel()

	cfg := createTestServerConfig(0, "gin") // Dynamic port allocated by createTestServerConfig
	s, err := server.CreateServerManager(cfg.Framework, cfg)
	require.NoError(t, err, "Failed to create server manager")
	require.NotNil(t, s, "Server manager instance should not be nil")

	webFramework := s.GetFramework()
	require.NotNil(t, webFramework, "WebFramework should not be nil")
	engine, ok := webFramework.GetNativeEngine().(*gin.Engine)
	require.True(t, ok, "Failed to assert Gin engine type")

	// 顶层路由 (Top-level route)
	engine.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// API V1 分组 (API V1 group)
	v1 := engine.Group("/api/v1")
	{
		v1.GET("/users", func(c *gin.Context) {
			c.String(http.StatusOK, "v1 users list")
		})
		v1.POST("/users", func(c *gin.Context) {
			c.String(http.StatusCreated, "v1 user created")
		})

		// 嵌套的 admin 分组 (Nested admin group)
		adminGroup := v1.Group("/admin")
		adminGroup.Use(func(c *gin.Context) { // Simulate a simple auth middleware
			if c.GetHeader("X-Admin-Token") != "secret-token" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				return
			}
			c.Next()
		})
		{
			adminGroup.GET("/dashboard", func(c *gin.Context) {
				c.String(http.StatusOK, "v1 admin dashboard")
			})
		}
	}

	// API V2 分组 (API V2 group)
	v2 := engine.Group("/api/v2")
	{
		v2.GET("/items", func(c *gin.Context) {
			c.String(http.StatusOK, "v2 items list")
		})
	}

	serverErrChan := make(chan error, 1)
	go func() {
		// s.Start() will block until server shuts down
		serverErrChan <- s.Start(context.Background())
	}()

	serverAddr := fmt.Sprintf("127.0.0.1:%d", cfg.Port) // Use the port from configuration
	require.Eventually(t, func() bool {
		conn, dialErr := net.DialTimeout("tcp", serverAddr, 100*time.Millisecond)
		if dialErr == nil {
			conn.Close()
			return true
		}
		return false
	}, 5*time.Second, 100*time.Millisecond, "Server did not start listening on %s", serverAddr)

	client := http.Client{Timeout: 3 * time.Second}

	testCases := []struct {
		name       string
		method     string
		path       string
		headers    map[string]string
		expectedSc int
		expectedBody string
	}{
		{"ping", "GET", "/ping", nil, http.StatusOK, "pong"},
		{"v1_users_get", "GET", "/api/v1/users", nil, http.StatusOK, "v1 users list"},
		{"v1_users_post", "POST", "/api/v1/users", nil, http.StatusCreated, "v1 user created"},
		{"v1_admin_dashboard_unauthorized", "GET", "/api/v1/admin/dashboard", nil, http.StatusUnauthorized, "Unauthorized"},
		{"v1_admin_dashboard_authorized", "GET", "/api/v1/admin/dashboard", map[string]string{"X-Admin-Token": "secret-token"}, http.StatusOK, "v1 admin dashboard"},
		{"v2_items_get", "GET", "/api/v2/items", nil, http.StatusOK, "v2 items list"},
		{"v1_nonexistent", "GET", "/api/v1/nonexistent", nil, http.StatusNotFound, ""}, // Gin returns 404 page, body might vary
	}

	for _, tc := range testCases {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			// 不要并行运行子测试，因为它们共享同一个服务器实例
			// (Don't run subtests in parallel as they share the same server instance)
			reqURL := fmt.Sprintf("http://%s%s", serverAddr, tc.path) // Use serverAddr
			req, err := http.NewRequest(tc.method, reqURL, nil)
			require.NoError(t, err)

			if tc.headers != nil {
				for k, v := range tc.headers {
					req.Header.Set(k, v)
				}
			}

			resp, err := client.Do(req)
			require.NoError(t, err, "Failed to make request for %s", tc.name)
			
			var responseBody string
			if resp != nil && resp.Body != nil {
				bodyBytes, readErr := io.ReadAll(resp.Body)
				require.NoError(t, readErr, "Failed to read response body for %s", tc.name)
				responseBody = string(bodyBytes)
				resp.Body.Close()
			}

			assert.Equal(t, tc.expectedSc, resp.StatusCode, "Status code mismatch for %s", tc.name)

			if tc.expectedBody != "" {
				if resp.StatusCode == http.StatusUnauthorized && strings.Contains(tc.path, "/admin/dashboard") {
					assert.Contains(t, responseBody, tc.expectedBody, "Response body mismatch for %s (auth error)", tc.name)
				} else {
					assert.Equal(t, tc.expectedBody, responseBody, "Response body mismatch for %s", tc.name)
				}
			} else if resp.StatusCode == http.StatusNotFound {
				// For 404, Gin often returns a default HTML page or a JSON error if customized.
				// We are expecting an empty body string from test case, so this is fine if it's empty or contains Gin's 404 message.
				t.Logf("Received 404 for %s, body: %s", tc.name, responseBody)
			}
		})
	}

	// Stop the server
	ctxStop, cancelStop := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelStop()
	errStop := s.Stop(ctxStop)
	assert.NoError(t, errStop, "Failed to stop server gracefully")

	select {
	case startErr, ok := <-serverErrChan:
		if ok && startErr != nil {
			// 服务器应该返回包装的 http.ErrServerClosed (Server should return wrapped http.ErrServerClosed)
			if !errors.Is(startErr, http.ErrServerClosed) {
				t.Errorf("s.Start() expected ErrServerClosed, got %v", startErr)
			}
		}
	case <-ctxStop.Done(): // Check against stopCtx as server start/stop is a unit of operation
		t.Errorf("Timeout waiting for s.Start() to return after stop: %v", ctxStop.Err())
	}
}