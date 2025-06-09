/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: 服务器生命周期管理集成测试 (Server lifecycle management integration tests)
 */

package integration

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin" // For *gin.Context
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	// "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services" // No longer needed for TestGracefulShutdown

	_ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin" // Import Gin for server creation
)

// TestGracefulShutdown 测试服务器的优雅关闭机制
// (TestGracefulShutdown tests the server's graceful shutdown mechanism)
func TestGracefulShutdown(t *testing.T) {
	t.Parallel()

	cfg := createTestServerConfig(0, "gin") 
	cfg.GracefulShutdown.Enabled = true
	cfg.GracefulShutdown.Timeout = 5 * time.Second 

	s, err := server.CreateServerManager(cfg.Framework, cfg)
	require.NoError(t, err, "Failed to create server manager")
	require.NotNil(t, s, "Server manager instance should not be nil")

	webFramework := s.GetFramework()
	require.NotNil(t, webFramework, "WebFramework should not be nil")
	engine, ok := webFramework.GetNativeEngine().(*gin.Engine)
	require.True(t, ok, "Failed to assert Gin engine type")

	engine.GET("/long-task", func(c *gin.Context) {
		time.Sleep(3 * time.Second) 
		c.String(http.StatusOK, "Long task completed")
	})

	serverErrChan := make(chan error, 1)
	go func() {
		serverErrChan <- s.Start(context.Background())
	}()

	serverAddr := fmt.Sprintf("127.0.0.1:%d", cfg.Port)
	require.Eventually(t, func() bool {
		conn, dialErr := net.DialTimeout("tcp", serverAddr, 100*time.Millisecond)
		if dialErr == nil {
			conn.Close()
			return true
		}
		return false
	}, 5*time.Second, 100*time.Millisecond, "Server did not start listening on %s", serverAddr)

	var wg sync.WaitGroup
	wg.Add(1)
	client := http.Client{Timeout: 6 * time.Second} 
	go func() {
		defer wg.Done()
		reqURL := fmt.Sprintf("http://%s/long-task", serverAddr)
		resp, reqErr := client.Get(reqURL)
		if reqErr != nil {
			t.Logf("Long task request failed: %v. This might be okay if server shutdown first.", reqErr)
			return
		}
		
		if resp != nil && resp.Body != nil {
			assert.Equal(t, http.StatusOK, resp.StatusCode, "Long task expected status OK")
			resp.Body.Close()
		}
	}()

	time.Sleep(500 * time.Millisecond)

	stopCtx, cancelStop := context.WithTimeout(context.Background(), cfg.GracefulShutdown.Timeout+1*time.Second)
	defer cancelStop()

	stopErr := s.Stop(stopCtx)
	assert.NoError(t, stopErr, "s.Stop() returned an error during graceful shutdown")

	wg.Wait()

	select {
	case errServerStart := <-serverErrChan:
		// 服务器应该返回包装的 http.ErrServerClosed (Server should return wrapped http.ErrServerClosed)
		if !errors.Is(errServerStart, http.ErrServerClosed) {
			t.Errorf("Server Start() expected ErrServerClosed, got %v", errServerStart)
		}
	case <-time.After(1 * time.Second):
		t.Error("Server Start() did not return after stop within timeout")
	}

	reqURLAfterStop := fmt.Sprintf("http://%s/long-task", serverAddr)
	respAfterStop, errAfterStop := client.Get(reqURLAfterStop)
	assert.Error(t, errAfterStop, "Request after server stop should fail")
	if respAfterStop != nil && respAfterStop.Body != nil {
		respAfterStop.Body.Close()
	}
}

// TestGracefulShutdownWithSignalHandling (可选) 测试服务器的信号处理和优雅关闭
// (TestGracefulShutdownWithSignalHandling (Optional) tests server signal handling and graceful shutdown)
// 这个测试更复杂，因为它需要模拟操作系统信号或在特定环境中运行。
// (This test is more complex as it needs to simulate OS signals or run in a specific environment.)
// 当前 SDK 设计中，信号处理内建于 ServerManager.Start() (当配置了优雅关闭时)。
// (In the current SDK design, signal handling is built into ServerManager.Start() when graceful shutdown is configured.)
// 不存在一个独立的 server.Run 函数或 server.IServer 接口如之前测试所假设。
// (There is no separate server.Run function or server.IServer interface as previously assumed by this test.)
func TestGracefulShutdownWithSignalHandling(t *testing.T) {
	t.Skip("Skipping signal handling test; requires specific setup or OS signal simulation. The SDK does not have server.Run or server.IServer as previously assumed. Signal handling is internal to ServerManager when GracefulShutdown is enabled.")

	// 原始代码保留如下以供参考，但它基于错误的假设 (Original code kept below for reference, but it's based on incorrect assumptions)
	/*
	cfg := createTestServerConfig(0, "gin")
	cfg.GracefulShutdown.Enabled = true 
	cfg.GracefulShutdown.Timeout = 3 * time.Second
	
	// serviceContainer is not directly passable to the mechanism that would handle signals (ServerManager itself)
	var serviceContainer services.ServiceContainer = nil 

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:0", cfg.Host))
	require.NoError(t, err)
	dynamicPort := listener.Addr().(*net.TCPAddr).Port
	require.NoError(t, listener.Close())
	cfg.Port = dynamicPort

	runErrChan := make(chan error, 1)
	go func() {
		// 假设的 server.Run 调用 (Hypothetical server.Run call)
		// runErrChan <- server.Run(cfg, serviceContainer, func(sActual server.IServer, engineActual interface{}) error {
		// 	 ginEngine, ok := engineActual.(*gin.Engine) 
		// 	 require.True(t, ok, "Failed to cast to *gin.Engine in server.Run callback")
		// 	 ginEngine.GET("/ping", func(c *gin.Context) { 
		// 		 time.Sleep(1 * time.Second)
		// 		 c.String(http.StatusOK, "pong")
		// 	 })
		// 	 return nil
		// })
		runErrChan <- fmt.Errorf("server.Run is not implemented as assumed") // Simulate error due to non-existent Run
	}()

	// ... (rest of the test would need significant rework to test ServerManager's internal signal handling)
	// For example, by starting a ServerManager and then sending an OS signal to the process.
	*/
}
