/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package integration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
	lmccerrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPerformanceIntegration 性能集成测试套件
// (TestPerformanceIntegration performance integration test suite)

// BenchmarkConfigurationLoading 配置加载性能基准测试
// (BenchmarkConfigurationLoading configuration loading performance benchmark)
func BenchmarkConfigurationLoading(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "config_loading_benchmark")
	require.NoError(b, err)
	defer os.RemoveAll(tempDir)

	configFile := filepath.Join(tempDir, "bench_config.yaml")
	configContent := `
server:
  host: "localhost"
  port: 8080
  mode: "production"
log:
  level: "info"
  format: "json"
database:
  type: "postgres"
  host: "localhost"
  port: 5432
  user: "benchuser"
  password: "benchpass"
  dbName: "benchdb"
`

	err = os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(b, err)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var cfg config.Config
			err := config.LoadConfig(&cfg,
				config.WithConfigFile(configFile, "yaml"),
			)
			if err != nil {
				b.Fatalf("Failed to load config: %v", err)
			}
		}
	})
}

// BenchmarkErrorCreation 错误创建性能基准测试
// (BenchmarkErrorCreation error creation performance benchmark)
func BenchmarkErrorCreation(b *testing.B) {
	errorMessages := []string{
		"simple error message",
		"validation failed for field 'username'",
		"database connection timeout",
		"invalid request format",
		"authentication failed",
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			msg := errorMessages[i%len(errorMessages)]
			err := lmccerrors.NewWithCode(lmccerrors.ErrValidation, msg)
			_ = err // 避免优化掉 (Avoid optimization)
			i++
		}
	})
}

// BenchmarkErrorWrapping 错误包装性能基准测试
// (BenchmarkErrorWrapping error wrapping performance benchmark)
func BenchmarkErrorWrapping(b *testing.B) {
	baseErr := lmccerrors.New("base error")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := lmccerrors.Wrap(baseErr, "operation failed")
			err = lmccerrors.Wrapf(err, "in context %d", 123)
			err = lmccerrors.WithCode(err, lmccerrors.ErrInternalServer)
			_ = err // 避免优化掉 (Avoid optimization)
		}
	})
}

// BenchmarkLogging 日志记录性能基准测试
// (BenchmarkLogging logging performance benchmark)
func BenchmarkLogging(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "logging_benchmark")
	require.NoError(b, err)
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "bench.log")

	logOpts := log.NewOptions()
	logOpts.Level = "info"
	logOpts.Format = "json"
	logOpts.OutputPaths = []string{logFile}
	logOpts.Development = false

	log.Init(logOpts)
	defer func() {
		defaultOpts := log.NewOptions()
		log.Init(defaultOpts)
	}()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			log.Infow("benchmark log message",
				"iteration", i,
				"worker_id", runtime.NumGoroutine(),
				"timestamp", time.Now().Unix(),
			)
			i++
		}
	})
}

// BenchmarkCrossModuleOperations 跨模块操作性能基准测试
// (BenchmarkCrossModuleOperations cross-module operations performance benchmark)
func BenchmarkCrossModuleOperations(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "cross_module_benchmark")
	require.NoError(b, err)
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "cross_bench.log")
	configFile := filepath.Join(tempDir, "cross_bench.yaml")

	// 设置配置 (Setup configuration)
	configContent := fmt.Sprintf(`
server:
  host: "localhost"
  port: 8080
log:
  level: "info"
  format: "json"
  outputPaths:
    - "%s"
  errorOutputPaths:
    - "%s"
`, logFile, logFile)

	err = os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(b, err)

	var cfg config.Config
	_, err = config.LoadConfigAndWatch(&cfg,
		config.WithConfigFile(configFile, "yaml"),
		config.WithHotReload(false), // 关闭热重载以提高性能 (Disable hot reload for performance)
	)
	require.NoError(b, err)

	// 配置日志 (Configure logging)
	logOpts := log.NewOptions()
	logOpts.Level = "info"
	logOpts.Format = "json"
	logOpts.OutputPaths = []string{logFile}
	log.Init(logOpts)

	defer func() {
		defaultOpts := log.NewOptions()
		log.Init(defaultOpts)
	}()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// 模拟跨模块操作 (Simulate cross-module operations)
			ctx := context.Background()
			ctx = log.ContextWithTraceID(ctx, fmt.Sprintf("bench-trace-%d", i))

			// 创建错误 (Create error)
			err := lmccerrors.NewWithCode(
				lmccerrors.ErrValidation,
				fmt.Sprintf("benchmark error %d", i),
			)

			// 记录错误 (Log error)
			log.Ctxw(ctx, "Benchmark operation",
				"iteration", i,
				"error", err,
				"error_code", lmccerrors.GetCoder(err).Code(),
			)

			i++
		}
	})
}

// TestHighConcurrencyStress 高并发压力测试
// (TestHighConcurrencyStress high concurrency stress test)
func TestHighConcurrencyStress(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	tempDir, err := os.MkdirTemp("", "stress_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "stress.log")
	configFile := filepath.Join(tempDir, "stress.yaml")

	// 配置文件 (Configuration file)
	configContent := fmt.Sprintf(`
server:
  host: "localhost"
  port: 8080
log:
  level: "debug"
  format: "json"
  outputPaths:
    - "%s"
  errorOutputPaths:
    - "%s"
`, logFile, logFile)

	err = os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	// 加载配置 (Load configuration)
	var cfg config.Config
	_, err = config.LoadConfigAndWatch(&cfg,
		config.WithConfigFile(configFile, "yaml"),
		config.WithHotReload(false),
	)
	require.NoError(t, err)

	// 配置日志 (Configure logging)
	logOpts := log.NewOptions()
	logOpts.Level = "debug"
	logOpts.Format = "json"
	logOpts.OutputPaths = []string{logFile}
	log.Init(logOpts)

	defer func() {
		defaultOpts := log.NewOptions()
		log.Init(defaultOpts)
	}()

	// 压力测试参数 (Stress test parameters)
	numGoroutines := 200
	operationsPerGoroutine := 1000
	totalOperations := numGoroutines * operationsPerGoroutine

	var wg sync.WaitGroup
	errorGroup := *lmccerrors.NewErrorGroup("stress test operations")
	var mu sync.Mutex

	start := time.Now()

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(workerID int) {
			defer wg.Done()

			for j := 0; j < operationsPerGoroutine; j++ {
				operationID := workerID*operationsPerGoroutine + j
				ctx := context.Background()
				ctx = log.ContextWithTraceID(ctx, fmt.Sprintf("stress-%d-%d", workerID, j))

				// 模拟不同类型的操作 (Simulate different types of operations)
				switch operationID % 4 {
				case 0:
					// 成功操作 (Successful operation)
					log.Ctxw(ctx, "Successful operation",
						"worker_id", workerID,
						"operation_id", operationID,
						"type", "success",
					)
				case 1:
					// 验证错误 (Validation error)
					err := lmccerrors.NewWithCode(lmccerrors.ErrValidation, "validation failed")
					log.Ctxw(ctx, "Validation error occurred",
						"worker_id", workerID,
						"operation_id", operationID,
						"error", err,
					)
					mu.Lock()
					errorGroup.Add(err)
					mu.Unlock()
				case 2:
					// 内部错误 (Internal error)
					err := lmccerrors.NewWithCode(lmccerrors.ErrInternalServer, "internal server error")
					log.Ctxw(ctx, "Internal error occurred",
						"worker_id", workerID,
						"operation_id", operationID,
						"error", err,
					)
					mu.Lock()
					errorGroup.Add(err)
					mu.Unlock()
				case 3:
					// 复杂错误包装 (Complex error wrapping)
					baseErr := lmccerrors.New("database timeout")
					wrappedErr := lmccerrors.Wrap(baseErr, "query execution failed")
					finalErr := lmccerrors.WithCode(wrappedErr, lmccerrors.ErrInternalServer)
					log.Ctxw(ctx, "Complex error occurred",
						"worker_id", workerID,
						"operation_id", operationID,
						"error", finalErr,
						"error_stack", fmt.Sprintf("%+v", finalErr),
					)
					mu.Lock()
					errorGroup.Add(finalErr)
					mu.Unlock()
				}
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	// 同步日志 (Sync logs)
	err = log.Std().GetZapLogger().Sync()
	require.NoError(t, err)

	// 验证结果 (Verify results)
	t.Logf("Stress test completed:")
	t.Logf("  Total operations: %d", totalOperations)
	t.Logf("  Duration: %v", duration)
	t.Logf("  Operations per second: %.2f", float64(totalOperations)/duration.Seconds())
	t.Logf("  Average time per operation: %v", duration/time.Duration(totalOperations))

	// 验证错误收集 (Verify error collection)
	errors := errorGroup.Errors()
	expectedErrors := totalOperations * 3 / 4 // 75% 的操作会产生错误 (75% of operations generate errors)
	assert.Equal(t, expectedErrors, len(errors))

	// 验证日志文件存在且有内容 (Verify log file exists and has content)
	logInfo, err := os.Stat(logFile)
	require.NoError(t, err)
	assert.Greater(t, logInfo.Size(), int64(0))

	// 性能断言 (Performance assertions)
	opsPerSecond := float64(totalOperations) / duration.Seconds()
	assert.Greater(t, opsPerSecond, 1000.0, "Should handle at least 1000 operations per second")

	avgTimePerOp := duration / time.Duration(totalOperations)
	assert.Less(t, avgTimePerOp, 10*time.Millisecond, "Average time per operation should be less than 10ms")
}

// TestMemoryLeakDetection 内存泄漏检测测试
// (TestMemoryLeakDetection memory leak detection test)
func TestMemoryLeakDetection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping memory leak test in short mode")
	}

	tempDir, err := os.MkdirTemp("", "memory_leak_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "memory.log")

	// 配置日志 (Configure logging)
	logOpts := log.NewOptions()
	logOpts.Level = "info"
	logOpts.Format = "json"
	logOpts.OutputPaths = []string{logFile}
	log.Init(logOpts)

	defer func() {
		defaultOpts := log.NewOptions()
		log.Init(defaultOpts)
	}()

	// 测量初始内存 (Measure initial memory)
	runtime.GC()
	var initialStats runtime.MemStats
	runtime.ReadMemStats(&initialStats)

	// 执行大量操作 (Perform many operations)
	numIterations := 10000
	for i := 0; i < numIterations; i++ {
		ctx := context.Background()
		ctx = log.ContextWithTraceID(ctx, fmt.Sprintf("leak-test-%d", i))

		// 创建错误 (Create errors)
		err := lmccerrors.NewWithCode(lmccerrors.ErrValidation, fmt.Sprintf("test error %d", i))
		wrappedErr := lmccerrors.Wrap(err, "wrapped error")

		// 记录日志 (Log messages)
		log.Ctxw(ctx, "Memory leak test operation",
			"iteration", i,
			"error", wrappedErr,
		)

		// 周期性强制GC (Periodic forced GC)
		if i%1000 == 0 {
			runtime.GC()
		}
	}

	// 强制垃圾回收并测量最终内存 (Force GC and measure final memory)
	runtime.GC()
	runtime.GC() // 连续两次GC确保清理完成 (Two consecutive GCs to ensure cleanup)
	var finalStats runtime.MemStats
	runtime.ReadMemStats(&finalStats)

	// 同步日志 (Sync logs)
	err = log.Std().GetZapLogger().Sync()
	require.NoError(t, err)

	// 分析内存使用 (Analyze memory usage)
	var memoryGrowth uint64
	if finalStats.HeapAlloc > initialStats.HeapAlloc {
		memoryGrowth = finalStats.HeapAlloc - initialStats.HeapAlloc
	} else {
		memoryGrowth = 0 // 内存实际减少了 (Memory actually decreased)
	}
	memoryPerOperation := memoryGrowth / uint64(numIterations)

	t.Logf("Memory leak detection results:")
	t.Logf("  Initial heap alloc: %d bytes", initialStats.HeapAlloc)
	t.Logf("  Final heap alloc: %d bytes", finalStats.HeapAlloc)
	t.Logf("  Memory growth: %d bytes", memoryGrowth)
	t.Logf("  Memory per operation: %d bytes", memoryPerOperation)
	t.Logf("  Number of GC cycles: %d", finalStats.NumGC-initialStats.NumGC)

	// 内存泄漏断言 (Memory leak assertions)
	// 每个操作的内存使用应该小于1KB (Memory usage per operation should be less than 1KB)
	assert.Less(t, memoryPerOperation, uint64(1024),
		"Memory usage per operation should be reasonable (< 1KB)")

	// 总的内存增长应该小于50MB (Total memory growth should be less than 50MB)
	assert.Less(t, memoryGrowth, uint64(50*1024*1024),
		"Total memory growth should be reasonable (< 50MB)")
}
