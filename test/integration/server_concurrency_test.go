/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: 服务器并发和性能测试 (Server concurrency and performance tests)
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
	"sync/atomic"
	"testing"
	"time"

	"github.com/gin-gonic/gin" // For *gin.Context and *gin.Engine
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	lmcclog "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"

	// services "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"

	_ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin" // Import Gin for server creation
)

// TestConcurrentRequests 测试服务器的并发请求处理能力
// (TestConcurrentRequests tests the server's ability to handle concurrent requests)
func TestConcurrentRequests(t *testing.T) {
	t.Parallel()

	cfg := createTestServerConfig(0, "gin") // Dynamic port allocated by createTestServerConfig
	s, err := server.CreateServerManager(cfg.Framework, cfg)
	require.NoError(t, err, "Failed to create server manager")
	require.NotNil(t, s, "Server manager instance should not be nil")

	webFramework := s.GetFramework()
	require.NotNil(t, webFramework, "WebFramework should not be nil")
	engine, ok := webFramework.GetNativeEngine().(*gin.Engine)
	require.True(t, ok, "Failed to assert Gin engine type")

	var requestCounter int32
	engine.GET("/concurrent-test", func(c *gin.Context) { // Corrected to *gin.Context
		atomic.AddInt32(&requestCounter, 1)
		time.Sleep(50 * time.Millisecond) // Simulate some work
		c.String(http.StatusOK, fmt.Sprintf("Request #%d processed", atomic.LoadInt32(&requestCounter)))
	})

	// Context for the overall test duration, including server start/stop
	ctxTest, cancelTest := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelTest()

	serverErrChan := make(chan error, 1)
	go func() {
		// s.Start() will block until server shuts down
		serverErrChan <- s.Start(context.Background())
	}()

	serverAddr := fmt.Sprintf("127.0.0.1:%d", cfg.Port) // Use the port from the configuration
	require.Eventually(t, func() bool {
		conn, dialErr := net.DialTimeout("tcp", serverAddr, 100*time.Millisecond)
		if dialErr == nil {
			conn.Close()
			return true
		}
		return false
	}, 5*time.Second, 100*time.Millisecond, "Server did not start listening on %s", serverAddr)

	numRequests := 100
	var wg sync.WaitGroup
	client := http.Client{Timeout: 5 * time.Second}

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(reqNum int) {
			defer wg.Done()
			reqURL := fmt.Sprintf("http://%s/concurrent-test", serverAddr)
			resp, reqErr := client.Get(reqURL)
			if reqErr != nil {
				t.Errorf("Request %d failed: %v", reqNum, reqErr)
				return
			}
			
			if resp != nil && resp.Body != nil {
				assert.Equal(t, http.StatusOK, resp.StatusCode, "Request %d expected status OK", reqNum)
				resp.Body.Close()
			}
		}(i)
	}

	wg.Wait()

	assert.Equal(t, int32(numRequests), atomic.LoadInt32(&requestCounter), "Expected %d requests to be processed", numRequests)

	// Stop the server
	stopCtx, stopCancel := context.WithTimeout(ctxTest, 5*time.Second) // Use test context as parent for stop
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
	case <-ctxTest.Done(): // Check if the overall test context timed out
		t.Errorf("Test timed out waiting for server operations: %v", ctxTest.Err())
	}
}

// BenchmarkServerPerformance 服务器性能基准测试
// (BenchmarkServerPerformance server performance benchmark test)
func BenchmarkServerPerformance(b *testing.B) {
	// Set up a null logger for benchmark to avoid I/O overhead
	nullLoggerOpts := lmcclog.NewOptions()
	nullLoggerOpts.Level = "fatal"
	nullLogger := lmcclog.NewLoggerWithWriter(nullLoggerOpts, io.Discard) // Use io.Discard
	require.NotNil(b, nullLogger, "Null logger should not be nil")

	originalGlobalLogger := lmcclog.Std()
	lmcclog.SetGlobalLogger(nullLogger)
	defer lmcclog.SetGlobalLogger(originalGlobalLogger)

	cfg := createTestServerConfig(0, "gin") // Dynamic port, ensure Framework is set
	// serviceContainer is not used with CreateServerManager
	s, err := server.CreateServerManager(cfg.Framework, cfg)
	require.NoError(b, err, "Failed to create server manager for benchmark")
	require.NotNil(b, s, "Server manager for benchmark should not be nil")

	webFramework := s.GetFramework()
	require.NotNil(b, webFramework, "WebFramework for benchmark should not be nil")
	engine, ok := webFramework.GetNativeEngine().(*gin.Engine)
	require.True(b, ok, "Failed to assert Gin engine type for benchmark")

	engine.GET("/benchmark-ping", func(c *gin.Context) { // Corrected to *gin.Context
		c.String(http.StatusOK, "pong")
	})

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:0", cfg.Host))
	require.NoError(b, err)
	dynamicPort := listener.Addr().(*net.TCPAddr).Port
	require.NoError(b, listener.Close())
	cfg.Port = dynamicPort

	go func() {
		// Corrected: s.Start() requires a context
		errStart := s.Start(context.Background())
		if errStart != nil && errStart != http.ErrServerClosed {
			b.Logf("s.Start() during benchmark returned an error: %v", errStart)
		}
	}()

	serverAddr := fmt.Sprintf("127.0.0.1:%d", dynamicPort)
	// Wait for server to be ready before resetting timer
	require.Eventually(b, func() bool {
		conn, dialErr := net.DialTimeout("tcp", serverAddr, 50*time.Millisecond)
		if dialErr == nil {
			conn.Close()
			return true
		}
		return false
	}, 3*time.Second, 100*time.Millisecond, "Benchmark server did not start listening on %s", serverAddr)

	client := http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 100,
		},
	}
	reqURL := fmt.Sprintf("http://%s/benchmark-ping", serverAddr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resp, reqErr := client.Get(reqURL)
		if reqErr != nil {
			b.Fatalf("Benchmark request failed: %v", reqErr)
		}
		if resp.StatusCode != http.StatusOK {
			b.Fatalf("Expected status OK, got %d", resp.StatusCode)
		}
		if resp.Body != nil {
			resp.Body.Close()
		}
	}
	b.StopTimer()

	stopCtx, cancelStop := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelStop()
	errStop := s.Stop(stopCtx)
	assert.NoError(b, errStop, "Failed to stop server gracefully after benchmark")
}
