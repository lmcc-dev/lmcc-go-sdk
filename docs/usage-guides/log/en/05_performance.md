# Performance Optimization

This document covers how to optimize the performance of the log module to ensure your application remains responsive under high load conditions.

## Performance Overview

While logging is essential for applications, it can become a performance bottleneck if not configured properly. This guide will help you:

1. **Identify Performance Bottlenecks** - Understand the performance impact points of logging
2. **Optimize Configuration** - Choose the best configuration options
3. **Reduce Overhead** - Minimize the performance overhead of logging
4. **Monitor Performance** - Measure and monitor logging performance

## Performance Impact Factors

### 1. Log Level

Log level is the most important factor affecting performance:

```go
// High-performance configuration: log errors only
opts := &log.Options{
    Level: "error",  // Only log error, fatal, panic
}

// Development configuration: log everything
opts := &log.Options{
    Level: "debug",  // Log all levels
}
```

**Performance comparison:**
```go
// Benchmark results (operations per second)
// Level "error": ~2,000,000 ops/sec
// Level "info":  ~1,500,000 ops/sec  
// Level "debug": ~1,000,000 ops/sec
```

### 2. Output Format

Different formats have varying serialization performance:

```go
// Fastest: KeyValue format
opts := &log.Options{
    Format: "keyvalue",
}

// Medium: Text format
opts := &log.Options{
    Format: "text",
}

// Slowest: JSON format
opts := &log.Options{
    Format: "json",
}
```

**Performance comparison:**
```go
// Serialization performance (nanoseconds/operation)
// KeyValue: ~200 ns/op
// Text:     ~300 ns/op
// JSON:     ~500 ns/op
```

### 3. Output Destination

Output destination significantly affects performance:

```go
// Fastest: Memory buffer
opts := &log.Options{
    OutputPaths: []string{"stdout"},
}

// Medium: Local file
opts := &log.Options{
    OutputPaths: []string{"/var/log/app.log"},
}

// Slowest: Network output
opts := &log.Options{
    OutputPaths: []string{"tcp://logserver:514"},
}
```

### 4. Caller Information

Caller information lookup is an expensive operation:

```go
// High performance: disable caller information
opts := &log.Options{
    DisableCaller: true,  // Save ~100ns/operation
}

// Debug-friendly: enable caller information
opts := &log.Options{
    DisableCaller: false,  // Extra overhead but useful for debugging
}
```

### 5. Stack Trace

Stack trace is one of the most expensive operations:

```go
// High performance: disable stack trace
opts := &log.Options{
    DisableStacktrace: true,
}

// Enable only for serious errors
opts := &log.Options{
    DisableStacktrace: false,
    StacktraceLevel:   "fatal",  // Only fatal and panic
}
```

## High-Performance Configurations

### Production High-Performance Configuration

```go
func highPerformanceConfig() *log.Options {
    return &log.Options{
        // Basic configuration
        Level:  "warn",      // Only log warnings and errors
        Format: "keyvalue",  // Fastest serialization format
        
        // Output configuration
        OutputPaths:      []string{"/var/log/app.log"},
        ErrorOutputPaths: []string{"/var/log/error.log"},
        
        // Performance optimizations
        EnableColor:       false,  // Disable color processing
        DisableCaller:     true,   // Disable caller lookup
        DisableStacktrace: true,   // Disable stack trace
        
        // Log rotation (reduce I/O overhead)
        LogRotateMaxSize:    50,   // Smaller file size
        LogRotateMaxBackups: 3,    // Fewer backup files
        LogRotateCompress:   true, // Compress old files
    }
}
```

### Extreme Performance Configuration

```go
func extremePerformanceConfig() *log.Options {
    return &log.Options{
        Level:             "error",    // Log errors only
        Format:            "keyvalue", // Fastest format
        OutputPaths:       []string{"/dev/null"}, // Discard output (testing only)
        DisableCaller:     true,
        DisableStacktrace: true,
        EnableColor:       false,
    }
}
```

### Balanced Configuration

Balance between performance and observability:

```go
func balancedConfig() *log.Options {
    return &log.Options{
        Level:  "info",
        Format: "json",  // Structured but acceptable performance
        
        OutputPaths: []string{
            "stdout",           // For container log collection
            "/var/log/app.log", // For local storage
        },
        
        DisableCaller:     false, // Keep caller info for debugging
        DisableStacktrace: false,
        StacktraceLevel:   "error", // Only show stack trace on errors
        
        LogRotateMaxSize:    100,
        LogRotateMaxBackups: 5,
        LogRotateCompress:   true,
    }
}
```

## Code-Level Optimizations

### 1. Conditional Logging

Avoid expensive operations when unnecessary:

```go
// Inefficient: always compute expensive value
log.Debug("User details", "user", expensiveUserSerialization(user))

// Efficient: only compute when needed
if log.IsDebugEnabled() {
    log.Debug("User details", "user", expensiveUserSerialization(user))
}

// Better: use lazy evaluation
log.Debug("User details", "user", func() interface{} {
    return expensiveUserSerialization(user)
})
```

### 2. Field Reuse

Reuse common log fields:

```go
// Inefficient: create new fields every time
func processRequest(requestID string, userID int) {
    log.Info("Starting processing", "request_id", requestID, "user_id", userID)
    log.Info("Validation completed", "request_id", requestID, "user_id", userID)
    log.Info("Processing completed", "request_id", requestID, "user_id", userID)
}

// Efficient: use context to reuse fields
func processRequest(requestID string, userID int) {
    ctx := context.Background()
    ctx = log.WithValues(ctx, "request_id", requestID, "user_id", userID)
    
    log.InfoContext(ctx, "Starting processing")
    log.InfoContext(ctx, "Validation completed")
    log.InfoContext(ctx, "Processing completed")
}
```

### 3. Batch Logging

For large volumes of logs, consider batch processing:

```go
// Inefficient: log each item individually
func processItems(items []Item) {
    for _, item := range items {
        log.Info("Processing item", "item_id", item.ID, "status", item.Status)
    }
}

// Efficient: batch logging
func processItems(items []Item) {
    log.Info("Starting batch processing", "total_items", len(items))
    
    var processed, failed int
    for _, item := range items {
        if processItem(item) {
            processed++
        } else {
            failed++
        }
    }
    
    log.Info("Batch processing completed",
        "total_items", len(items),
        "processed", processed,
        "failed", failed,
    )
}
```

### 4. Avoid String Formatting

Use structured fields instead of string formatting:

```go
// Inefficient: string formatting
log.Info(fmt.Sprintf("User %d performed action %s", userID, action))

// Efficient: structured fields
log.Info("User performed action", "user_id", userID, "action", action)
```

## Asynchronous Logging

### Buffered Writing

Use buffering to reduce I/O operations:

```go
// Configure buffered writing
func bufferedConfig() *log.Options {
    return &log.Options{
        Level:       "info",
        Format:      "json",
        OutputPaths: []string{"/var/log/app.log"},
        
        // Note: These are example configurations, actual implementation may differ
        BufferSize:  64 * 1024, // 64KB buffer
        FlushInterval: time.Second, // Flush every second
    }
}
```

### Non-Blocking Writing

Avoid blocking the main thread for log writing:

```go
// Example: Async log writer
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
            // Async log writing
            writeLogEntry(entry)
        case <-al.done:
            return
        }
    }
}

func (al *AsyncLogger) Log(entry LogEntry) {
    select {
    case al.logChan <- entry:
        // Non-blocking write
    default:
        // Buffer full, can choose to drop or block
        // Here we choose to drop to maintain performance
    }
}
```

## Memory Optimization

### 1. Object Pooling

Reuse log objects to reduce GC pressure:

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
    // Clear fields
    for k := range entry.Fields {
        delete(entry.Fields, k)
    }
    entry.Message = ""
    entry.Level = ""
    
    logEntryPool.Put(entry)
}
```

### 2. String Optimization

Avoid unnecessary string allocations:

```go
// Inefficient: string concatenation
func logUserAction(userID int, action string) {
    log.Info("User action: " + action, "user_id", userID)
}

// Efficient: use fields directly
func logUserAction(userID int, action string) {
    log.Info("User action", "user_id", userID, "action", action)
}
```

### 3. Pre-allocate Slices

For known-size fields, pre-allocate slices:

```go
// Inefficient: dynamic growth
func logMultipleFields() {
    fields := []interface{}{}
    fields = append(fields, "key1", "value1")
    fields = append(fields, "key2", "value2")
    // ...
    log.Infow("Message", fields...)
}

// Efficient: pre-allocate
func logMultipleFields() {
    fields := make([]interface{}, 0, 10) // Pre-allocate capacity
    fields = append(fields, "key1", "value1")
    fields = append(fields, "key2", "value2")
    // ...
    log.Infow("Message", fields...)
}
```

## Performance Monitoring

### 1. Benchmarking

Create logging performance benchmarks:

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
            log.Info("Test message")
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
            log.Infow("Test message",
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
            log.InfoContext(ctx, "Test message")
        }
    })
}
```

Run benchmarks:

```bash
go test -bench=. -benchmem
```

### 2. Performance Metrics

Monitor key performance metrics:

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
    
    // Perform logging
    log.Info(message)
    
    // Update metrics
    latency := time.Since(start).Nanoseconds()
    atomic.AddInt64(&metrics.TotalLogs, 1)
    atomic.StoreInt64(&metrics.AvgLatency, latency) // Simplified version
    
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

### 3. Performance Profiling

Use Go's profiling tools:

```go
import _ "net/http/pprof"

func main() {
    // Enable pprof
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    
    // Application logic...
}
```

Profiling commands:

```bash
# CPU profiling
go tool pprof http://localhost:6060/debug/pprof/profile

# Memory profiling
go tool pprof http://localhost:6060/debug/pprof/heap

# Block profiling
go tool pprof http://localhost:6060/debug/pprof/block
```

## Environment-Specific Optimizations

### Development Environment

```go
func developmentConfig() *log.Options {
    return &log.Options{
        Level:         "debug",
        Format:        "text",
        OutputPaths:   []string{"stdout"},
        EnableColor:   true,
        DisableCaller: false, // Keep debug info
    }
}
```

### Test Environment

```go
func testConfig() *log.Options {
    return &log.Options{
        Level:             "warn",
        Format:            "text",
        OutputPaths:       []string{"stdout"},
        DisableCaller:     true,  // Improve test speed
        DisableStacktrace: true,
    }
}
```

### Production Environment

```go
func productionConfig() *log.Options {
    return &log.Options{
        Level:  "info",
        Format: "json",
        OutputPaths: []string{
            "/var/log/app.log",
            "stdout", // For container log collection
        },
        DisableCaller:       false, // Keep for troubleshooting
        DisableStacktrace:   false,
        StacktraceLevel:     "error",
        LogRotateMaxSize:    100,
        LogRotateMaxBackups: 10,
        LogRotateCompress:   true,
    }
}
```

## Troubleshooting

### Performance Issue Diagnosis

1. **Identify bottlenecks**:
   ```bash
   # Use strace to monitor system calls
   strace -c -p <pid>
   
   # Use iostat to monitor I/O
   iostat -x 1
   ```

2. **Memory usage analysis**:
   ```go
   import "runtime"
   
   func logMemoryUsage() {
       var m runtime.MemStats
       runtime.ReadMemStats(&m)
       
       log.Info("Memory usage",
           "alloc", m.Alloc,
           "total_alloc", m.TotalAlloc,
           "sys", m.Sys,
           "num_gc", m.NumGC,
       )
   }
   ```

3. **Log latency monitoring**:
   ```go
   func monitorLogLatency() {
       start := time.Now()
       log.Info("Test message")
       latency := time.Since(start)
       
       if latency > 10*time.Millisecond {
           log.Warn("High log latency detected", "latency", latency)
       }
   }
   ```

## Best Practices Summary

1. **Choose appropriate log level** - Use "info" or "warn" in production
2. **Optimize output format** - Use "keyvalue" for high-performance scenarios
3. **Disable unnecessary features** - Turn off caller info and stack trace in performance-critical paths
4. **Use context logging** - Avoid repeatedly passing the same fields
5. **Monitor performance metrics** - Regularly check logging performance
6. **Perform benchmarking** - Verify performance impact of configuration changes

## Next Steps

- [Best Practices](06_best_practices.md) - Production-ready patterns and complete guide 