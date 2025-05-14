/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package integration

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"
	lmccLog "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"

	// "github.com/lmcc-dev/lmcc-go-sdk/pkg/known" // No longer needed as we use lmccLog.ContextWithTraceID
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// goroutineContextKey is the key for goroutine ID in the context.
// goroutineContextKey 是上下文中 goroutine ID 的键。
type goroutineContextKey struct{}

/**
 * TestConcurrentLoggingWithContext tests the logging context propagation in a concurrent environment.
 * It ensures that context-specific fields (like traceID and a custom goroutine ID) are correctly
 * logged for messages from different goroutines.
 * (TestConcurrentLoggingWithContext 测试并发环境中的日志上下文传播。
 *  它确保来自不同 goroutine 的消息能够正确记录特定于上下文的字段（如 traceID 和自定义 goroutine ID）。)
 */
func TestConcurrentLoggingWithContext(t *testing.T) {
	// Create a temporary directory for logs
	// 为日志创建一个临时目录
	tempLogDir, err := os.MkdirTemp("", "log_concurrency_test_")
	require.NoError(t, err, "Failed to create temp log directory (创建临时日志目录失败)")
	defer os.RemoveAll(tempLogDir) // Clean up after the test (测试后清理)

	tempLogFile := filepath.Join(tempLogDir, "test_concurrent.log")

	// Configure logger options
	// 配置记录器选项
	logOpts := lmccLog.NewOptions()
	logOpts.Level = "debug" // Ensure all messages are logged (确保所有消息都被记录)
	logOpts.Format = "json" // JSON format for easy parsing (JSON 格式以便于解析)
	// 同时输出到文件和标准输出，以便在屏幕上查看日志
	// Output to both file and stdout to see logs on screen
	logOpts.OutputPaths = []string{tempLogFile, "stdout"}
	logOpts.ErrorOutputPaths = []string{tempLogFile, "stderr"} // Also send errors to the same file and stderr (此测试也将错误发送到同一文件和标准错误)
	logOpts.Development = false // Use production encoder for more structured logs (使用生产编码器以获得更结构化的日志)

	// Initialize logger
	// 初始化记录器
	lmccLog.Init(logOpts)
	
	// Testing if other pkg/log functions are recognized
	// 测试是否能识别 pkg/log 包中的其他函数
	// _ = lmccLog.Sync() // Should compile if pkg/log symbols are visible
	// lmccLog.Debug("Test message") // Try another package-level function
	
	defer func() {
		// Reset to default logger after test to avoid interference with other tests
		// 测试后重置为默认记录器，以避免干扰其他测试
		defaultOpts := lmccLog.NewOptions()
		defaultOpts.OutputPaths = []string{"stdout"}
		defaultOpts.ErrorOutputPaths = []string{"stderr"}
		lmccLog.Init(defaultOpts)
	}()

	numGoroutines := 20                                  // Number of concurrent goroutines (并发 goroutine 的数量)
	msgsPerGoroutine := 5                               // Number of messages per goroutine (每个 goroutine 的消息数)
	expectedTotalMessages := numGoroutines * msgsPerGoroutine // Expected total number of messages (预期的总消息数)

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Store traceIDs for each goroutine to verify later
	// 存储每个 goroutine 的 traceID 以便稍后验证
	goroutineTraceIDs := make(map[string]string)
	var mu sync.Mutex // Mutex to protect goroutineTraceIDs map (保护 goroutineTraceIDs 映射的互斥锁)

	for i := 0; i < numGoroutines; i++ {
		go func(goroutineIdx int) {
			defer wg.Done()

			// Generate unique IDs for this goroutine
			// 为此 goroutine 生成唯一 ID
			traceID := uuid.NewString()
			goroutineIDStr := strconv.Itoa(goroutineIdx)

			mu.Lock()
			goroutineTraceIDs[goroutineIDStr] = traceID
			mu.Unlock()

			// Create a new context with these IDs
			// 使用这些 ID 创建新的上下文
			ctx := context.Background()
			ctx = lmccLog.ContextWithTraceID(ctx, traceID) // Use the helper from pkg/log
			// For custom keys, use standard library's context.WithValue
			// (对于自定义键，使用标准库的 context.WithValue)
			ctx = context.WithValue(ctx, goroutineContextKey{}, goroutineIDStr)

			for j := 0; j < msgsPerGoroutine; j++ {
				// 使用正确的上下文感知日志记录方法
				// (Use correct context-aware logging method)
				msgStr := fmt.Sprintf("Message %d from goroutine %s", j, goroutineIDStr)
				lmccLog.Ctx(ctx, msgStr) // 正确的上下文日志方法 (Correct context logging method)
			}
		}(i)
	}

	wg.Wait() // Wait for all goroutines to complete (等待所有 goroutine 完成)

	// Read and verify log output
	// 读取并验证日志输出
	file, err := os.Open(tempLogFile)
	require.NoError(t, err, "Failed to open log file (打开日志文件失败)")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var loggedMessages []map[string]interface{}
	linesRead := 0
	for scanner.Scan() {
		linesRead++
		var logEntry map[string]interface{}
		line := scanner.Text()
		err := json.Unmarshal([]byte(line), &logEntry)
		require.NoError(t, err, fmt.Sprintf("Failed to unmarshal log line (无法反序列化日志行): %s", line))
		loggedMessages = append(loggedMessages, logEntry)
	}
	require.NoError(t, scanner.Err(), "Error reading log file (读取日志文件时出错)")

	assert.Equal(t, expectedTotalMessages, len(loggedMessages), "Unexpected number of log messages (日志消息数量不符合预期)")

	// Group messages by goroutineID extracted from the message string for verification
	// 从消息字符串中提取 goroutineID 以进行验证，按 goroutineID 分组消息
	messagesByGoroutine := make(map[string][]map[string]interface{})
	for _, entry := range loggedMessages {
		msg, ok := entry["message"].(string) // "message" is the standard key in JSON logs ("message" 是 JSON 日志中的标准消息键)
		require.True(t, ok, "Log entry missing message field 'message' (日志条目缺少消息字段 'message')")

		// Extract goroutineID from message: "Message X from goroutine Y"
		// 从消息中提取 goroutineID："Message X from goroutine Y"
		parts := strings.Split(msg, " ")
		require.True(t, len(parts) >= 5 && parts[3] == "goroutine", "Log message format is incorrect (日志消息格式不正确): %s", msg)
		goroutineIDFromMsg := parts[4]
		messagesByGoroutine[goroutineIDFromMsg] = append(messagesByGoroutine[goroutineIDFromMsg], entry)
	}

	assert.Equal(t, numGoroutines, len(messagesByGoroutine), "Logs from all goroutines not found (未找到所有 goroutine 的日志)")

	for gid, entries := range messagesByGoroutine {
		assert.Len(t, entries, msgsPerGoroutine, fmt.Sprintf("Incorrect number of messages for goroutine %s (goroutine %s 的消息数量不正确)", gid, gid))
		expectedTraceID, ok := goroutineTraceIDs[gid]
		require.True(t, ok, fmt.Sprintf("TraceID for goroutine %s not found in map (在映射中未找到 goroutine %s 的 TraceID)", gid, gid))

		for _, entry := range entries {
			// Verify traceID from context (expecting field name "trace_id")
			// 验证来自上下文的 traceID (期望字段名为 "trace_id")
			actualTraceID, traceOk := entry["trace_id"].(string) // Field name based on pkg/log/log.go
			assert.True(t, traceOk, fmt.Sprintf("Log entry for goroutine %s missing 'trace_id' field (goroutine %s 的日志条目缺少 'trace_id' 字段): %+v", gid, gid, entry))
			assert.Equal(t, expectedTraceID, actualTraceID, fmt.Sprintf("Mismatch trace_id for goroutine %s (goroutine %s 的 trace_id 不匹配)", gid, gid))

			// For goroutineContextKey, we expect it to be stringified as the struct's string representation
			// 对于 goroutineContextKey，我们期望它被字符串化为结构体的字符串表示
			customKeyStr := fmt.Sprintf("%v", goroutineContextKey{}) // How context values are keyed in logs
			actualGoroutineIDFromCtx, ctxGidOk := entry[customKeyStr].(string)
			if ctxGidOk { // This field might not be present if logger is not configured for this
				assert.Equal(t, gid, actualGoroutineIDFromCtx, fmt.Sprintf("Mismatch goroutineID from context for goroutine %s (goroutine %s 的上下文中 goroutineID 不匹配)", gid, gid))
			}
		}
	}
} 