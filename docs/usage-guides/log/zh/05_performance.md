# 性能优化

本文档介绍如何优化日志模块的性能，确保在高负载环境中保持应用程序的响应性。

## 性能概述

日志记录虽然对应用程序至关重要，但如果配置不当，可能会成为性能瓶颈。本指南将帮助您：

1. **识别性能瓶颈** - 了解日志记录的性能影响点
2. **优化配置** - 选择最佳的配置选项
3. **减少开销** - 最小化日志记录的性能开销
4. **监控性能** - 测量和监控日志性能

## 性能影响因素

### 1. 日志级别

日志级别是影响性能的最重要因素：

```go
// 高性能配置：仅记录错误
opts := &log.Options{
    Level: "error",  // 只记录 error、fatal、panic
}

// 开发配置：记录所有信息
opts := &log.Options{
    Level: "debug",  // 记录所有级别的日志
}
```

**性能对比：**
```go
// 基准测试结果（每秒操作数）
// Level "error": ~2,000,000 ops/sec
// Level "info":  ~1,500,000 ops/sec  
// Level "debug": ~1,000,000 ops/sec
```

### 2. 输出格式

不同格式的序列化性能差异：

```go
// 最快：KeyValue 格式
opts := &log.Options{
    Format: "keyvalue",
}

// 中等：Text 格式
opts := &log.Options{
    Format: "text",
}

// 最慢：JSON 格式
opts := &log.Options{
    Format: "json",
}
```

**性能对比：**
```go
// 序列化性能（纳秒/操作）
// KeyValue: ~200 ns/op
// Text:     ~300 ns/op
// JSON:     ~500 ns/op
```

### 3. 输出目标

输出目标对性能有显著影响：

```go
// 最快：内存缓冲
opts := &log.Options{
    OutputPaths: []string{"stdout"},
}

// 中等：本地文件
opts := &log.Options{
    OutputPaths: []string{"/var/log/app.log"},
}

// 最慢：网络输出
opts := &log.Options{
    OutputPaths: []string{"tcp://logserver:514"},
}
```

### 4. 调用者信息

调用者信息查找是昂贵的操作：

```go
// 高性能：禁用调用者信息
opts := &log.Options{
    DisableCaller: true,  // 节省 ~100ns/操作
}

// 调试友好：启用调用者信息
opts := &log.Options{
    DisableCaller: false,  // 额外开销但便于调试
}
```

### 5. 堆栈跟踪

堆栈跟踪是最昂贵的操作之一：

```go
// 高性能：禁用堆栈跟踪
opts := &log.Options{
    DisableStacktrace: true,
}

// 仅在严重错误时启用
opts := &log.Options{
    DisableStacktrace: false,
    StacktraceLevel:   "fatal",  // 仅 fatal 和 panic
}
```

## 高性能配置

### 生产环境高性能配置

```go
func highPerformanceConfig() *log.Options {
    return &log.Options{
        // 基本配置
        Level:  "warn",      // 仅记录警告和错误
        Format: "keyvalue",  // 最快的序列化格式
        
        // 输出配置
        OutputPaths:      []string{"/var/log/app.log"},
        ErrorOutputPaths: []string{"/var/log/error.log"},
        
        // 性能优化
        EnableColor:       false,  // 禁用颜色处理
        DisableCaller:     true,   // 禁用调用者查找
        DisableStacktrace: true,   // 禁用堆栈跟踪
        
        // 日志轮转（减少 I/O 开销）
        LogRotateMaxSize:    50,   // 较小的文件大小
        LogRotateMaxBackups: 3,    // 较少的备份文件
        LogRotateCompress:   true, // 压缩旧文件
    }
}
```

### 极致性能配置

```go
func extremePerformanceConfig() *log.Options {
    return &log.Options{
        Level:             "error",    // 仅记录错误
        Format:            "keyvalue", // 最快格式
        OutputPaths:       []string{"/dev/null"}, // 丢弃输出（仅用于测试）
        DisableCaller:     true,
        DisableStacktrace: true,
        EnableColor:       false,
    }
}
```

### 平衡配置

在性能和可观测性之间取得平衡：

```go
func balancedConfig() *log.Options {
    return &log.Options{
        Level:  "info",
        Format: "json",  // 结构化但性能可接受
        
        OutputPaths: []string{
            "stdout",           // 用于容器日志收集
            "/var/log/app.log", // 用于本地存储
        },
        
        DisableCaller:     false, // 保留调用者信息用于调试
        DisableStacktrace: false,
        StacktraceLevel:   "error", // 仅在错误时显示堆栈
        
        LogRotateMaxSize:    100,
        LogRotateMaxBackups: 5,
        LogRotateCompress:   true,
    }
}
```

## 代码级优化

### 1. 条件日志记录

避免在不必要时进行昂贵的操作：

```go
// 低效：总是计算昂贵的值
log.Debug("用户详情", "user", expensiveUserSerialization(user))

// 高效：仅在需要时计算
if log.IsDebugEnabled() {
    log.Debug("用户详情", "user", expensiveUserSerialization(user))
}

// 更好：使用延迟计算
log.Debug("用户详情", "user", func() interface{} {
    return expensiveUserSerialization(user)
})
```

### 2. 字段重用

重用常见的日志字段：

```go
// 低效：每次都创建新的字段
func processRequest(requestID string, userID int) {
    log.Info("开始处理", "request_id", requestID, "user_id", userID)
    log.Info("验证完成", "request_id", requestID, "user_id", userID)
    log.Info("处理完成", "request_id", requestID, "user_id", userID)
}

// 高效：使用上下文重用字段
func processRequest(requestID string, userID int) {
    ctx := context.Background()
    ctx = log.WithValues(ctx, "request_id", requestID, "user_id", userID)
    
    log.InfoContext(ctx, "开始处理")
    log.InfoContext(ctx, "验证完成")
    log.InfoContext(ctx, "处理完成")
}
```

### 3. 批量日志记录

对于大量日志，考虑批量处理：

```go
// 低效：逐个记录
func processItems(items []Item) {
    for _, item := range items {
        log.Info("处理项目", "item_id", item.ID, "status", item.Status)
    }
}

// 高效：批量记录
func processItems(items []Item) {
    log.Info("开始批量处理", "total_items", len(items))
    
    var processed, failed int
    for _, item := range items {
        if processItem(item) {
            processed++
        } else {
            failed++
        }
    }
    
    log.Info("批量处理完成",
        "total_items", len(items),
        "processed", processed,
        "failed", failed,
    )
}
```

### 4. 避免字符串格式化

直接使用结构化字段而不是格式化字符串：

```go
// 低效：字符串格式化
log.Info(fmt.Sprintf("用户 %d 执行了操作 %s", userID, action))

// 高效：结构化字段
log.Info("用户执行操作", "user_id", userID, "action", action)
```

## 异步日志记录

### 缓冲写入

使用缓冲减少 I/O 操作：

```go
// 配置缓冲写入
func bufferedConfig() *log.Options {
    return &log.Options{
        Level:       "info",
        Format:      "json",
        OutputPaths: []string{"/var/log/app.log"},
        
        // 注意：这些是示例配置，实际实现可能不同
        BufferSize:  64 * 1024, // 64KB 缓冲区
        FlushInterval: time.Second, // 每秒刷新
    }
}
```

### 非阻塞写入

避免日志写入阻塞主线程：

```go
// 示例：异步日志写入器
type AsyncLogger struct {
    logChan chan LogEntry
    done    chan struct{}
}

func NewAsyncLogger(bufferSize int) *AsyncLogger {
    al := &AsyncLogger{
        logChan: make(chan LogEntry, bufferSize),
        done:    make(chan struct{}),
    }
    
    go al.worker()
    return al
}

func (al *AsyncLogger) worker() {
    for {
        select {
        case entry := <-al.logChan:
            // 异步写入日志
            writeLogEntry(entry)
        case <-al.done:
            return
        }
    }
}

func (al *AsyncLogger) Log(entry LogEntry) {
    select {
    case al.logChan <- entry:
        // 非阻塞写入
    default:
        // 缓冲区满，可以选择丢弃或阻塞
        // 这里选择丢弃以保持性能
    }
}
```

## 内存优化

### 1. 对象池

重用日志对象以减少 GC 压力：

```go
import "sync"

var logEntryPool = sync.Pool{
    New: func() interface{} {
        return &LogEntry{
            Fields: make(map[string]interface{}, 8),
        }
    },
}

func getLogEntry() *LogEntry {
    return logEntryPool.Get().(*LogEntry)
}

func putLogEntry(entry *LogEntry) {
    // 清理字段
    for k := range entry.Fields {
        delete(entry.Fields, k)
    }
    entry.Message = ""
    entry.Level = ""
    
    logEntryPool.Put(entry)
}
```

### 2. 字符串优化

避免不必要的字符串分配：

```go
// 低效：字符串连接
func logUserAction(userID int, action string) {
    log.Info("用户操作: " + action, "user_id", userID)
}

// 高效：直接使用字段
func logUserAction(userID int, action string) {
    log.Info("用户操作", "user_id", userID, "action", action)
}
```

### 3. 预分配切片

对于已知大小的字段，预分配切片：

```go
// 低效：动态增长
func logMultipleFields() {
    fields := []interface{}{}
    fields = append(fields, "key1", "value1")
    fields = append(fields, "key2", "value2")
    // ...
    log.Infow("消息", fields...)
}

// 高效：预分配
func logMultipleFields() {
    fields := make([]interface{}, 0, 10) // 预分配容量
    fields = append(fields, "key1", "value1")
    fields = append(fields, "key2", "value2")
    // ...
    log.Infow("消息", fields...)
}
```

## 性能监控

### 1. 基准测试

创建日志性能基准测试：

```go
package main

import (
    "testing"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

func BenchmarkLogInfo(b *testing.B) {
    opts := &log.Options{
        Level:       "info",
        Format:      "json",
        OutputPaths: []string{"/dev/null"},
    }
    log.Init(opts)
    
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            log.Info("测试消息")
        }
    })
}

func BenchmarkLogInfoWithFields(b *testing.B) {
    opts := &log.Options{
        Level:       "info",
        Format:      "json",
        OutputPaths: []string{"/dev/null"},
    }
    log.Init(opts)
    
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            log.Infow("测试消息",
                "user_id", 123,
                "action", "test",
                "timestamp", time.Now(),
            )
        }
    })
}

func BenchmarkLogContext(b *testing.B) {
    opts := &log.Options{
        Level:       "info",
        Format:      "json",
        OutputPaths: []string{"/dev/null"},
    }
    log.Init(opts)
    
    ctx := context.Background()
    ctx = log.WithValues(ctx, "request_id", "req-123", "user_id", 456)
    
    b.ResetTimer()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            log.InfoContext(ctx, "测试消息")
        }
    })
}
```

运行基准测试：

```bash
go test -bench=. -benchmem
```

### 2. 性能指标

监控关键性能指标：

```go
import (
    "time"
    "sync/atomic"
)

type LogMetrics struct {
    TotalLogs     int64
    ErrorLogs     int64
    AvgLatency    int64
    MaxLatency    int64
}

var metrics LogMetrics

func logWithMetrics(level string, message string) {
    start := time.Now()
    
    // 执行日志记录
    log.Info(message)
    
    // 更新指标
    latency := time.Since(start).Nanoseconds()
    atomic.AddInt64(&metrics.TotalLogs, 1)
    atomic.StoreInt64(&metrics.AvgLatency, latency) // 简化版本
    
    if latency > atomic.LoadInt64(&metrics.MaxLatency) {
        atomic.StoreInt64(&metrics.MaxLatency, latency)
    }
}

func getMetrics() LogMetrics {
    return LogMetrics{
        TotalLogs:  atomic.LoadInt64(&metrics.TotalLogs),
        ErrorLogs:  atomic.LoadInt64(&metrics.ErrorLogs),
        AvgLatency: atomic.LoadInt64(&metrics.AvgLatency),
        MaxLatency: atomic.LoadInt64(&metrics.MaxLatency),
    }
}
```

### 3. 性能分析

使用 Go 的性能分析工具：

```go
import _ "net/http/pprof"

func main() {
    // 启用 pprof
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    
    // 应用程序逻辑...
}
```

分析命令：

```bash
# CPU 分析
go tool pprof http://localhost:6060/debug/pprof/profile

# 内存分析
go tool pprof http://localhost:6060/debug/pprof/heap

# 阻塞分析
go tool pprof http://localhost:6060/debug/pprof/block
```

## 环境特定优化

### 开发环境

```go
func developmentConfig() *log.Options {
    return &log.Options{
        Level:         "debug",
        Format:        "text",
        OutputPaths:   []string{"stdout"},
        EnableColor:   true,
        DisableCaller: false, // 保留调试信息
    }
}
```

### 测试环境

```go
func testConfig() *log.Options {
    return &log.Options{
        Level:             "warn",
        Format:            "text",
        OutputPaths:       []string{"stdout"},
        DisableCaller:     true,  // 提高测试速度
        DisableStacktrace: true,
    }
}
```

### 生产环境

```go
func productionConfig() *log.Options {
    return &log.Options{
        Level:  "info",
        Format: "json",
        OutputPaths: []string{
            "/var/log/app.log",
            "stdout", // 用于容器日志收集
        },
        DisableCaller:       false, // 保留用于故障排除
        DisableStacktrace:   false,
        StacktraceLevel:     "error",
        LogRotateMaxSize:    100,
        LogRotateMaxBackups: 10,
        LogRotateCompress:   true,
    }
}
```

## 故障排除

### 性能问题诊断

1. **识别瓶颈**：
   ```bash
   # 使用 strace 监控系统调用
   strace -c -p <pid>
   
   # 使用 iostat 监控 I/O
   iostat -x 1
   ```

2. **内存使用分析**：
   ```go
   import "runtime"
   
   func logMemoryUsage() {
       var m runtime.MemStats
       runtime.ReadMemStats(&m)
       
       log.Info("内存使用情况",
           "alloc", m.Alloc,
           "total_alloc", m.TotalAlloc,
           "sys", m.Sys,
           "num_gc", m.NumGC,
       )
   }
   ```

3. **日志延迟监控**：
   ```go
   func monitorLogLatency() {
       start := time.Now()
       log.Info("测试消息")
       latency := time.Since(start)
       
       if latency > 10*time.Millisecond {
           log.Warn("日志延迟过高", "latency", latency)
       }
   }
   ```

## 最佳实践总结

1. **选择合适的日志级别** - 生产环境使用 "info" 或 "warn"
2. **优化输出格式** - 高性能场景使用 "keyvalue"
3. **禁用不必要的功能** - 在性能关键路径上禁用调用者信息和堆栈跟踪
4. **使用上下文日志** - 避免重复传递相同字段
5. **监控性能指标** - 定期检查日志性能
6. **进行基准测试** - 验证配置更改的性能影响

## 下一步

- [最佳实践](06_best_practices.md) - 生产就绪模式和完整指南 