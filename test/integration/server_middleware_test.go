/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: 服务器中间件集成测试 (Server middleware integration tests)
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

	"github.com/gin-gonic/gin" // Import Gin directly for H type
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	lmccerrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"

	// middleware "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/middleware"
	// services "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"

	// Ensure Gin plugin is registered
	_ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin"
)

// TestMiddlewareIntegration 测试中间件链的执行和自定义中间件
// (TestMiddlewareIntegration tests middleware chain execution and custom middleware)
func TestMiddlewareIntegration(t *testing.T) {
	t.Parallel()

	cfg := createTestServerConfig(0, "gin") // Dynamic port, ensure framework is set in config
	// serviceContainer is not used with CreateServerManager
	s, err := server.CreateServerManager(cfg.Framework, cfg)
	require.NoError(t, err, "Failed to create server manager")
	require.NotNil(t, s, "Server manager instance should not be nil")

	webFramework := s.GetFramework()
	require.NotNil(t, webFramework, "WebFramework should not be nil")
	engine, ok := webFramework.GetNativeEngine().(*gin.Engine)
	require.True(t, ok, "Failed to assert Gin engine type")

	var middlewareOrder []string

	// 自定义中间件1 (Custom Middleware 1)
	customMiddleware1 := func(c *gin.Context) {
		middlewareOrder = append(middlewareOrder, "mw1-before")
		c.Next()
		middlewareOrder = append(middlewareOrder, "mw1-after")
	}

	// 自定义中间件2 (Custom Middleware 2)
	customMiddleware2 := func(c *gin.Context) {
		middlewareOrder = append(middlewareOrder, "mw2-before")
		c.Next()
		middlewareOrder = append(middlewareOrder, "mw2-after")
	}

	// 自定义中间件3 (错误处理) (Custom Middleware 3 - error handling)
	customErrorMiddleware := func(c *gin.Context) {
		middlewareOrder = append(middlewareOrder, "errorMw-before")
		c.Next() // Process subsequent middleware and handlers
		middlewareOrder = append(middlewareOrder, "errorMw-after")

		if len(c.Errors) > 0 {
			lastError := c.Errors.Last()
			coder := lmccerrors.GetCoder(lastError.Err)
			if coder != nil {
				c.JSON(coder.HTTPStatus(), gin.H{"errorCode": coder.Code(), "message": coder.String()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"message": lastError.Err.Error()})
			}
			c.Abort() // Prevent other handlers from writing to the response
			return
		}
	}

	engine.Use(customErrorMiddleware, customMiddleware1, customMiddleware2)

	engine.GET("/test-middleware", func(c *gin.Context) {
		middlewareOrder = append(middlewareOrder, "handler")
		c.String(http.StatusOK, "Handler executed")
	})

	engine.GET("/test-error-middleware", func(c *gin.Context) {
		middlewareOrder = append(middlewareOrder, "errorHandler-before")
		_ = c.Error(lmccerrors.WithCode(lmccerrors.New("something went wrong"), lmccerrors.ErrInternalServer))
		middlewareOrder = append(middlewareOrder, "errorHandler-after")
	})

	serverErrChan := make(chan error, 1)
	go func() {
		// s.Start() will block until server shuts down
		serverErrChan <- s.Start(context.Background())
	}()

	serverAddr := fmt.Sprintf("127.0.0.1:%d", cfg.Port) // Use port from configuration
	require.Eventually(t, func() bool {
		conn, dialErr := net.DialTimeout("tcp", serverAddr, 50*time.Millisecond)
		if dialErr == nil {
			conn.Close()
			return true
		}
		return false
	}, 3*time.Second, 100*time.Millisecond, "Server did not start listening on %s", serverAddr)

	client := http.Client{Timeout: 3 * time.Second}
	reqURLNormal := fmt.Sprintf("http://%s/test-middleware", serverAddr)
	respNormal, err := client.Get(reqURLNormal)
	require.NoError(t, err, "Failed to make GET request to /test-middleware")
	require.NotNil(t, respNormal)
	var tbodyNormal string
	if respNormal != nil && respNormal.Body != nil {
		bodyBytes, _ := io.ReadAll(respNormal.Body)
		tbodyNormal = string(bodyBytes)
		assert.Equal(t, http.StatusOK, respNormal.StatusCode, "Expected status OK for /test-middleware")
		assert.Equal(t, "Handler executed", tbodyNormal, "Response body mismatch for /test-middleware")
		respNormal.Body.Close()
	}

	expectedOrderNormal := []string{"errorMw-before", "mw1-before", "mw2-before", "handler", "mw2-after", "mw1-after", "errorMw-after"}
	assert.Equal(t, expectedOrderNormal, middlewareOrder, "Middleware execution order mismatch for normal flow")

	middlewareOrder = []string{} // Reset for next test

	reqURLError := fmt.Sprintf("http://%s/test-error-middleware", serverAddr)
	respError, err := client.Get(reqURLError)
	require.NoError(t, err, "Failed to make GET request to /test-error-middleware")
	require.NotNil(t, respError)
	var tbodyError string
	if respError != nil && respError.Body != nil {
		bodyBytes, _ := io.ReadAll(respError.Body)
		tbodyError = string(bodyBytes)
		assert.Equal(t, lmccerrors.ErrInternalServer.HTTPStatus(), respError.StatusCode, "Expected status for /test-error-middleware")
		assert.Contains(t, tbodyError, lmccerrors.ErrInternalServer.String(), "Error message mismatch")
		assert.Contains(t, tbodyError, fmt.Sprintf(`"errorCode":%d`, lmccerrors.ErrInternalServer.Code()), "Error code mismatch")
		respError.Body.Close()
	}

	expectedOrderError := []string{"errorMw-before", "mw1-before", "mw2-before", "errorHandler-before", "errorHandler-after", "mw2-after", "mw1-after", "errorMw-after"}
	assert.Equal(t, expectedOrderError, middlewareOrder, "Middleware execution order mismatch for error flow")
	
	ctxStop, cancelStop := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancelStop() // cancelStop is called just before s.Stop for the first server instance
	errStop := s.Stop(ctxStop)
	assert.NoError(t, errStop, "Failed to stop server gracefully")
	cancelStop() // Call cancel after stop to release resources if Stop doesn't use the context for its full duration

	select {
	case startErr, ok := <-serverErrChan:
		if ok && startErr != nil {
			// 服务器应该返回包装的 http.ErrServerClosed (Server should return wrapped http.ErrServerClosed)
			if !errors.Is(startErr, http.ErrServerClosed) {
				t.Errorf("s.Start() expected ErrServerClosed, got %v", startErr)
			}
		}
	case <-time.After(2 * time.Second):
		t.Error("s.Start() did not return after stop within timeout")
	}


	middlewareOrder = []string{} // Reset for recovery test
	
	recoveryCfg := createTestServerConfig(0, "gin") // Dynamic port allocated by createTestServerConfig
	sRecovery, err := server.CreateServerManager(recoveryCfg.Framework, recoveryCfg)
	require.NoError(t, err, "Failed to create recovery server manager")
	require.NotNil(t, sRecovery, "Recovery server manager should not be nil")

	recoveryWebFramework := sRecovery.GetFramework()
	require.NotNil(t, recoveryWebFramework, "Recovery WebFramework should not be nil")
	recoveryEngine, ok := recoveryWebFramework.GetNativeEngine().(*gin.Engine)
	require.True(t, ok, "Failed to assert Gin engine type for recovery server")

	// Use Gin's standard recovery middleware
	recoveryEngine.Use(gin.Recovery())

	recoveryEngine.GET("/panic", func(c *gin.Context) {
		middlewareOrder = append(middlewareOrder, "panic-handler-before")
		panic("test panic in recovery")
	})

	serverErrChanRecovery := make(chan error, 1)
	go func() {
		// sRecovery.Start() will block until server shuts down
		serverErrChanRecovery <- sRecovery.Start(context.Background())
	}()

	serverAddrRecovery := fmt.Sprintf("127.0.0.1:%d", recoveryCfg.Port) // Use port from configuration
	require.Eventually(t, func() bool {
		conn, dialErr := net.DialTimeout("tcp", serverAddrRecovery, 50*time.Millisecond)
		if dialErr == nil {
			conn.Close()
			return true
		}
		return false
	}, 3*time.Second, 100*time.Millisecond, "Recovery server did not start listening on %s", serverAddrRecovery)

	reqURLPanic := fmt.Sprintf("http://%s/panic", serverAddrRecovery)
	respPanic, err := client.Get(reqURLPanic)
	require.NoError(t, err, "Failed to make GET request to /panic")
	require.NotNil(t, respPanic)
	var tbodyPanic string
	if respPanic != nil && respPanic.Body != nil {
		bodyBytes, _ := io.ReadAll(respPanic.Body)
		tbodyPanic = string(bodyBytes)
		assert.Equal(t, http.StatusInternalServerError, respPanic.StatusCode, "Expected status InternalServerError for /panic")
		// Gin's recovery middleware may not write a response body, so we just check the status code
		// (Gin的recovery中间件可能不会写入响应体，我们只检查状态码)
		t.Logf("Panic response body: %q", tbodyPanic) // Log the actual response body for debugging
		respPanic.Body.Close()
	}

	ctxStopRecovery, cancelStopRecovery := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancelStopRecovery() // Call explicitly after Stop
	errStopRecovery := sRecovery.Stop(ctxStopRecovery)
	assert.NoError(t, errStopRecovery, "Failed to stop recovery server gracefully")
	cancelStopRecovery() // Call cancel after stop

	select {
	case errStart := <-serverErrChanRecovery:
		if errStart != nil { // Check if an error was actually sent
			// 服务器应该返回包装的 http.ErrServerClosed (Server should return wrapped http.ErrServerClosed)
			if !errors.Is(errStart, http.ErrServerClosed) {
				t.Errorf("sRecovery.Start() expected ErrServerClosed, got %v", errStart)
			}
		}
	case <-time.After(2 * time.Second):
		t.Error("sRecovery.Start() did not return after stop within timeout")
	}
}
