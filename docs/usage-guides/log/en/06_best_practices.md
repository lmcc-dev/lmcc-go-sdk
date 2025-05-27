# Best Practices

This document provides comprehensive best practices for using the log module in production environments, covering everything from configuration to monitoring.

## Overview

Following these best practices will help you:

1. **Maintain Performance** - Keep your application responsive
2. **Ensure Reliability** - Avoid logging-related failures
3. **Improve Observability** - Get meaningful insights from logs
4. **Simplify Maintenance** - Make logs easy to manage and analyze
5. **Enhance Security** - Protect sensitive information

## Configuration Best Practices

### 1. Environment-Specific Configuration

Use different configurations for different environments:

```go
package main

import (
    "os"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

func getLogConfig() *log.Options {
    env := os.Getenv("APP_ENV")
    
    switch env {
    case "development":
        return developmentConfig()
    case "test":
        return testConfig()
    case "staging":
        return stagingConfig()
    case "production":
        return productionConfig()
    default:
        return defaultConfig()
    }
}

func developmentConfig() *log.Options {
    return &log.Options{
        Level:         "debug",
        Format:        "text",
        OutputPaths:   []string{"stdout"},
        EnableColor:   true,
        DisableCaller: false,
    }
}

func productionConfig() *log.Options {
    return &log.Options{
        Level:  "info",
        Format: "json",
        OutputPaths: []string{
            "/var/log/app.log",
            "stdout", // For container environments
        },
        ErrorOutputPaths:    []string{"/var/log/error.log"},
        DisableCaller:       false,
        DisableStacktrace:   false,
        StacktraceLevel:     "error",
        LogRotateMaxSize:    100,
        LogRotateMaxBackups: 10,
        LogRotateCompress:   true,
    }
}

func testConfig() *log.Options {
    return &log.Options{
        Level:             "warn",
        Format:            "text",
        OutputPaths:       []string{"stdout"},
        DisableCaller:     true,
        DisableStacktrace: true,
    }
}
```

### 2. Configuration Validation

Always validate your configuration:

```go
func validateLogConfig(opts *log.Options) error {
    // Validate log level
    validLevels := []string{"debug", "info", "warn", "error", "fatal", "panic"}
    if !contains(validLevels, opts.Level) {
        return fmt.Errorf("invalid log level: %s", opts.Level)
    }
    
    // Validate format
    validFormats := []string{"text", "json", "keyvalue"}
    if !contains(validFormats, opts.Format) {
        return fmt.Errorf("invalid log format: %s", opts.Format)
    }
    
    // Validate output paths
    if len(opts.OutputPaths) == 0 {
        return fmt.Errorf("output paths cannot be empty")
    }
    
    // Validate file paths
    for _, path := range opts.OutputPaths {
        if path != "stdout" && path != "stderr" {
            if err := validateFilePath(path); err != nil {
                return fmt.Errorf("invalid output path %s: %v", path, err)
            }
        }
    }
    
    return nil
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}

func validateFilePath(path string) error {
    dir := filepath.Dir(path)
    if _, err := os.Stat(dir); os.IsNotExist(err) {
        return fmt.Errorf("directory does not exist: %s", dir)
    }
    return nil
}
```

### 3. Configuration from Environment Variables

Support configuration via environment variables:

```go
func configFromEnv() *log.Options {
    opts := &log.Options{
        Level:  getEnvOrDefault("LOG_LEVEL", "info"),
        Format: getEnvOrDefault("LOG_FORMAT", "json"),
    }
    
    // Parse output paths
    if paths := os.Getenv("LOG_OUTPUT_PATHS"); paths != "" {
        opts.OutputPaths = strings.Split(paths, ",")
    } else {
        opts.OutputPaths = []string{"stdout"}
    }
    
    // Parse boolean options
    opts.EnableColor = getEnvBool("LOG_ENABLE_COLOR", false)
    opts.DisableCaller = getEnvBool("LOG_DISABLE_CALLER", false)
    opts.DisableStacktrace = getEnvBool("LOG_DISABLE_STACKTRACE", false)
    
    // Parse rotation options
    opts.LogRotateMaxSize = getEnvInt("LOG_ROTATE_MAX_SIZE", 100)
    opts.LogRotateMaxBackups = getEnvInt("LOG_ROTATE_MAX_BACKUPS", 10)
    opts.LogRotateCompress = getEnvBool("LOG_ROTATE_COMPRESS", true)
    
    return opts
}

func getEnvOrDefault(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
    if value := os.Getenv(key); value != "" {
        if parsed, err := strconv.ParseBool(value); err == nil {
            return parsed
        }
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if parsed, err := strconv.Atoi(value); err == nil {
            return parsed
        }
    }
    return defaultValue
}
```

## Logging Patterns

### 1. Structured Logging

Always use structured logging with key-value pairs:

```go
// Good: Structured logging
log.Info("User login successful",
    "user_id", userID,
    "username", username,
    "ip_address", clientIP,
    "user_agent", userAgent,
    "login_time", time.Now(),
)

// Bad: String interpolation
log.Info(fmt.Sprintf("User %s (ID: %d) logged in from %s", username, userID, clientIP))
```

### 2. Consistent Field Names

Use consistent field names across your application:

```go
// Define constants for common field names
const (
    FieldUserID     = "user_id"
    FieldRequestID  = "request_id"
    FieldOperation  = "operation"
    FieldDuration   = "duration"
    FieldError      = "error"
    FieldStatusCode = "status_code"
    FieldMethod     = "method"
    FieldPath       = "path"
    FieldIPAddress  = "ip_address"
)

// Use consistent field names
log.Info("API request completed",
    FieldRequestID, requestID,
    FieldUserID, userID,
    FieldMethod, r.Method,
    FieldPath, r.URL.Path,
    FieldStatusCode, statusCode,
    FieldDuration, duration,
)
```

### 3. Context-Aware Logging

Use context to propagate common fields:

```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    // Create context with common fields
    ctx := r.Context()
    ctx = log.WithValues(ctx,
        FieldRequestID, generateRequestID(),
        FieldMethod, r.Method,
        FieldPath, r.URL.Path,
        FieldIPAddress, getClientIP(r),
    )
    
    log.InfoContext(ctx, "Request started")
    
    // Pass context to other functions
    user, err := authenticateUser(ctx, r)
    if err != nil {
        log.ErrorContext(ctx, "Authentication failed", FieldError, err)
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    
    // Add user info to context
    ctx = log.WithValues(ctx, FieldUserID, user.ID)
    
    processRequest(ctx, w, r, user)
}

func authenticateUser(ctx context.Context, r *http.Request) (*User, error) {
    log.InfoContext(ctx, "Authenticating user")
    
    // Authentication logic...
    
    return user, nil
}

func processRequest(ctx context.Context, w http.ResponseWriter, r *http.Request, user *User) {
    log.InfoContext(ctx, "Processing request")
    
    // Request processing logic...
    
    log.InfoContext(ctx, "Request completed")
}
```

### 4. Error Logging

Follow consistent error logging patterns:

```go
// Log errors with context and details
func processPayment(ctx context.Context, payment *Payment) error {
    log.InfoContext(ctx, "Starting payment processing",
        "payment_id", payment.ID,
        "amount", payment.Amount,
        "currency", payment.Currency,
    )
    
    if err := validatePayment(payment); err != nil {
        log.ErrorContext(ctx, "Payment validation failed",
            FieldError, err,
            "payment_id", payment.ID,
            "validation_errors", getValidationErrors(err),
        )
        return fmt.Errorf("payment validation failed: %w", err)
    }
    
    if err := chargeCard(ctx, payment); err != nil {
        log.ErrorContext(ctx, "Card charge failed",
            FieldError, err,
            "payment_id", payment.ID,
            "card_last_four", payment.Card.LastFour,
            "error_code", getErrorCode(err),
        )
        return fmt.Errorf("card charge failed: %w", err)
    }
    
    log.InfoContext(ctx, "Payment processed successfully",
        "payment_id", payment.ID,
        "transaction_id", payment.TransactionID,
    )
    
    return nil
}
```

### 5. Performance Logging

Log performance metrics consistently:

```go
func trackOperation(ctx context.Context, operation string, fn func() error) error {
    start := time.Now()
    
    log.InfoContext(ctx, "Operation started",
        FieldOperation, operation,
        "start_time", start,
    )
    
    err := fn()
    duration := time.Since(start)
    
    if err != nil {
        log.ErrorContext(ctx, "Operation failed",
            FieldOperation, operation,
            FieldDuration, duration,
            FieldError, err,
        )
        return err
    }
    
    log.InfoContext(ctx, "Operation completed",
        FieldOperation, operation,
        FieldDuration, duration,
    )
    
    // Log slow operations
    if duration > 5*time.Second {
        log.WarnContext(ctx, "Slow operation detected",
            FieldOperation, operation,
            FieldDuration, duration,
        )
    }
    
    return nil
}

// Usage
func processOrder(ctx context.Context, order *Order) error {
    return trackOperation(ctx, "process_order", func() error {
        // Order processing logic...
        return nil
    })
}
```

## Security Best Practices

### 1. Sensitive Data Protection

Never log sensitive information:

```go
// Sensitive fields that should never be logged
var sensitiveFields = map[string]bool{
    "password":     true,
    "credit_card":  true,
    "ssn":         true,
    "api_key":     true,
    "token":       true,
    "secret":      true,
}

func sanitizeFields(fields map[string]interface{}) map[string]interface{} {
    sanitized := make(map[string]interface{})
    
    for key, value := range fields {
        if sensitiveFields[strings.ToLower(key)] {
            sanitized[key] = "[REDACTED]"
        } else {
            sanitized[key] = value
        }
    }
    
    return sanitized
}

// Safe logging function
func logUserAction(ctx context.Context, action string, user *User, data map[string]interface{}) {
    sanitizedData := sanitizeFields(data)
    
    log.InfoContext(ctx, "User action",
        "action", action,
        FieldUserID, user.ID,
        "username", user.Username, // OK to log username
        "data", sanitizedData,
    )
}
```

### 2. Data Masking

Mask sensitive data when logging is necessary:

```go
func maskCreditCard(cardNumber string) string {
    if len(cardNumber) < 4 {
        return "****"
    }
    return "****-****-****-" + cardNumber[len(cardNumber)-4:]
}

func maskEmail(email string) string {
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return "[INVALID_EMAIL]"
    }
    
    username := parts[0]
    domain := parts[1]
    
    if len(username) <= 2 {
        return "**@" + domain
    }
    
    return username[:2] + "***@" + domain
}

func logPaymentAttempt(ctx context.Context, payment *Payment) {
    log.InfoContext(ctx, "Payment attempt",
        "payment_id", payment.ID,
        "amount", payment.Amount,
        "card_number", maskCreditCard(payment.CardNumber),
        "customer_email", maskEmail(payment.CustomerEmail),
    )
}
```

### 3. Log Access Control

Implement proper access controls for log files:

```go
func setupLogFiles() error {
    logDir := "/var/log/myapp"
    
    // Create log directory with restricted permissions
    if err := os.MkdirAll(logDir, 0750); err != nil {
        return fmt.Errorf("failed to create log directory: %w", err)
    }
    
    // Set ownership (example for Unix systems)
    if err := os.Chown(logDir, os.Getuid(), os.Getgid()); err != nil {
        return fmt.Errorf("failed to set log directory ownership: %w", err)
    }
    
    return nil
}
```

## Performance Best Practices

### 1. Conditional Logging

Use conditional logging for expensive operations:

```go
func processLargeDataset(ctx context.Context, data []DataItem) {
    if log.IsDebugEnabled() {
        log.DebugContext(ctx, "Processing dataset",
            "item_count", len(data),
            "sample_items", serializeItems(data[:min(5, len(data))]), // Only serialize if debug is enabled
        )
    }
    
    // Processing logic...
}

// Helper function to check if debug logging is enabled
func (l *Logger) IsDebugEnabled() bool {
    return l.level <= DebugLevel
}
```

### 2. Lazy Evaluation

Use lazy evaluation for expensive field computation:

```go
type LazyField struct {
    fn func() interface{}
}

func (lf LazyField) String() string {
    return fmt.Sprintf("%v", lf.fn())
}

func NewLazyField(fn func() interface{}) LazyField {
    return LazyField{fn: fn}
}

// Usage
func logComplexOperation(ctx context.Context, data *ComplexData) {
    log.InfoContext(ctx, "Complex operation",
        "data_summary", NewLazyField(func() interface{} {
            return computeExpensiveSummary(data) // Only computed if log level allows
        }),
    )
}
```

### 3. Batch Logging

For high-volume scenarios, consider batch logging:

```go
type BatchLogger struct {
    entries []LogEntry
    mutex   sync.Mutex
    ticker  *time.Ticker
    done    chan struct{}
}

func NewBatchLogger(flushInterval time.Duration) *BatchLogger {
    bl := &BatchLogger{
        entries: make([]LogEntry, 0, 100),
        ticker:  time.NewTicker(flushInterval),
        done:    make(chan struct{}),
    }
    
    go bl.flushLoop()
    return bl
}

func (bl *BatchLogger) Log(entry LogEntry) {
    bl.mutex.Lock()
    defer bl.mutex.Unlock()
    
    bl.entries = append(bl.entries, entry)
    
    // Flush if buffer is full
    if len(bl.entries) >= 100 {
        bl.flush()
    }
}

func (bl *BatchLogger) flushLoop() {
    for {
        select {
        case <-bl.ticker.C:
            bl.mutex.Lock()
            bl.flush()
            bl.mutex.Unlock()
        case <-bl.done:
            return
        }
    }
}

func (bl *BatchLogger) flush() {
    if len(bl.entries) == 0 {
        return
    }
    
    // Write all entries
    for _, entry := range bl.entries {
        writeLogEntry(entry)
    }
    
    // Clear buffer
    bl.entries = bl.entries[:0]
}
```

## Monitoring and Alerting

### 1. Log Metrics

Track important log metrics:

```go
import (
    "sync/atomic"
    "time"
)

type LogMetrics struct {
    TotalLogs    int64
    ErrorLogs    int64
    WarnLogs     int64
    InfoLogs     int64
    DebugLogs    int64
    LastLogTime  int64
}

var metrics LogMetrics

func incrementLogCounter(level string) {
    atomic.AddInt64(&metrics.TotalLogs, 1)
    atomic.StoreInt64(&metrics.LastLogTime, time.Now().Unix())
    
    switch level {
    case "error":
        atomic.AddInt64(&metrics.ErrorLogs, 1)
    case "warn":
        atomic.AddInt64(&metrics.WarnLogs, 1)
    case "info":
        atomic.AddInt64(&metrics.InfoLogs, 1)
    case "debug":
        atomic.AddInt64(&metrics.DebugLogs, 1)
    }
}

func GetLogMetrics() LogMetrics {
    return LogMetrics{
        TotalLogs:   atomic.LoadInt64(&metrics.TotalLogs),
        ErrorLogs:   atomic.LoadInt64(&metrics.ErrorLogs),
        WarnLogs:    atomic.LoadInt64(&metrics.WarnLogs),
        InfoLogs:    atomic.LoadInt64(&metrics.InfoLogs),
        DebugLogs:   atomic.LoadInt64(&metrics.DebugLogs),
        LastLogTime: atomic.LoadInt64(&metrics.LastLogTime),
    }
}
```

### 2. Health Checks

Include logging health in your health checks:

```go
func logHealthCheck() error {
    // Check if logging is working
    testMessage := fmt.Sprintf("Health check at %s", time.Now().Format(time.RFC3339))
    
    start := time.Now()
    log.Info(testMessage)
    duration := time.Since(start)
    
    // Check if logging is too slow
    if duration > 100*time.Millisecond {
        return fmt.Errorf("logging is too slow: %v", duration)
    }
    
    // Check if log files are writable
    if err := checkLogFileWritable(); err != nil {
        return fmt.Errorf("log file not writable: %w", err)
    }
    
    return nil
}

func checkLogFileWritable() error {
    // Implementation depends on your log configuration
    testFile := "/var/log/app.log.test"
    
    file, err := os.Create(testFile)
    if err != nil {
        return err
    }
    defer os.Remove(testFile)
    defer file.Close()
    
    _, err = file.WriteString("test")
    return err
}
```

### 3. Error Rate Monitoring

Monitor error rates and alert on anomalies:

```go
type ErrorRateMonitor struct {
    window     time.Duration
    threshold  float64
    errorCount int64
    totalCount int64
    lastReset  time.Time
    mutex      sync.RWMutex
}

func NewErrorRateMonitor(window time.Duration, threshold float64) *ErrorRateMonitor {
    return &ErrorRateMonitor{
        window:    window,
        threshold: threshold,
        lastReset: time.Now(),
    }
}

func (erm *ErrorRateMonitor) RecordLog(level string) {
    erm.mutex.Lock()
    defer erm.mutex.Unlock()
    
    // Reset counters if window has passed
    if time.Since(erm.lastReset) > erm.window {
        erm.errorCount = 0
        erm.totalCount = 0
        erm.lastReset = time.Now()
    }
    
    erm.totalCount++
    if level == "error" {
        erm.errorCount++
    }
    
    // Check if error rate exceeds threshold
    if erm.totalCount > 0 {
        errorRate := float64(erm.errorCount) / float64(erm.totalCount)
        if errorRate > erm.threshold {
            // Trigger alert
            log.Warn("High error rate detected",
                "error_rate", errorRate,
                "threshold", erm.threshold,
                "error_count", erm.errorCount,
                "total_count", erm.totalCount,
            )
        }
    }
}
```

## Testing Best Practices

### 1. Log Testing

Test your logging in unit tests:

```go
func TestUserLogin(t *testing.T) {
    // Capture logs for testing
    var logBuffer bytes.Buffer
    opts := &log.Options{
        Level:       "debug",
        Format:      "json",
        OutputPaths: []string{"memory://test"},
    }
    
    // Initialize logger with test configuration
    log.Init(opts)
    
    // Test the function
    err := loginUser("testuser", "password123")
    assert.NoError(t, err)
    
    // Verify logs
    logs := getTestLogs()
    assert.Contains(t, logs, "User login successful")
    assert.Contains(t, logs, "testuser")
}

func getTestLogs() []string {
    // Implementation depends on your test setup
    // This is a simplified example
    return []string{} // Return captured log messages
}
```

### 2. Mock Logging

Use mock loggers for testing:

```go
type MockLogger struct {
    logs []LogEntry
}

func (m *MockLogger) Info(msg string, fields ...interface{}) {
    m.logs = append(m.logs, LogEntry{
        Level:   "info",
        Message: msg,
        Fields:  parseFields(fields),
    })
}

func (m *MockLogger) GetLogs() []LogEntry {
    return m.logs
}

func (m *MockLogger) Clear() {
    m.logs = nil
}

// Test with mock logger
func TestWithMockLogger(t *testing.T) {
    mockLogger := &MockLogger{}
    
    // Inject mock logger
    originalLogger := log.GetLogger()
    log.SetLogger(mockLogger)
    defer log.SetLogger(originalLogger)
    
    // Test your function
    processUser(&User{ID: 123, Name: "Test"})
    
    // Verify logs
    logs := mockLogger.GetLogs()
    assert.Len(t, logs, 1)
    assert.Equal(t, "User processed", logs[0].Message)
}
```

## Deployment Best Practices

### 1. Log Rotation

Configure proper log rotation:

```go
func setupLogRotation() *log.Options {
    return &log.Options{
        Level:  "info",
        Format: "json",
        OutputPaths: []string{
            "/var/log/app.log",
        },
        
        // Rotation settings
        LogRotateMaxSize:    100,  // 100MB
        LogRotateMaxBackups: 30,   // Keep 30 files
        LogRotateMaxAge:     7,    // Keep for 7 days
        LogRotateCompress:   true, // Compress old files
    }
}
```

### 2. Container Logging

Best practices for containerized applications:

```go
func containerLogConfig() *log.Options {
    return &log.Options{
        Level:  "info",
        Format: "json",
        OutputPaths: []string{
            "stdout", // Container runtime will handle this
        },
        ErrorOutputPaths: []string{
            "stderr",
        },
        
        // Disable file-based features in containers
        DisableCaller: false, // Keep for debugging
        EnableColor:   false, // No color in container logs
    }
}
```

### 3. Graceful Shutdown

Ensure logs are flushed on shutdown:

```go
func main() {
    // Initialize logging
    log.Init(getLogConfig())
    
    // Setup graceful shutdown
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    
    go func() {
        <-c
        log.Info("Shutting down gracefully...")
        
        // Flush logs before exit
        log.Sync()
        
        os.Exit(0)
    }()
    
    // Application logic...
    
    log.Info("Application started")
    select {} // Keep running
}
```

## Common Pitfalls to Avoid

### 1. Over-Logging

```go
// Bad: Too much logging
func processItems(items []Item) {
    log.Info("Starting to process items")
    for i, item := range items {
        log.Info("Processing item", "index", i, "item_id", item.ID)
        // Process item...
        log.Info("Item processed", "index", i, "item_id", item.ID)
    }
    log.Info("Finished processing items")
}

// Good: Appropriate logging
func processItems(items []Item) {
    log.Info("Starting to process items", "total_count", len(items))
    
    var processed, failed int
    for _, item := range items {
        if err := processItem(item); err != nil {
            log.Error("Failed to process item", "item_id", item.ID, "error", err)
            failed++
        } else {
            processed++
        }
    }
    
    log.Info("Finished processing items",
        "total_count", len(items),
        "processed", processed,
        "failed", failed,
    )
}
```

### 2. Inconsistent Field Names

```go
// Bad: Inconsistent field names
log.Info("User login", "userId", 123)
log.Info("User logout", "user_id", 123)
log.Info("User action", "UserID", 123)

// Good: Consistent field names
const FieldUserID = "user_id"

log.Info("User login", FieldUserID, 123)
log.Info("User logout", FieldUserID, 123)
log.Info("User action", FieldUserID, 123)
```

### 3. Logging in Hot Paths

```go
// Bad: Logging in tight loops
func processLargeDataset(data []DataItem) {
    for _, item := range data { // This could be millions of items
        log.Debug("Processing item", "item_id", item.ID) // Don't do this
        processItem(item)
    }
}

// Good: Batch logging or conditional logging
func processLargeDataset(data []DataItem) {
    log.Info("Starting to process large dataset", "item_count", len(data))
    
    batchSize := 1000
    for i := 0; i < len(data); i += batchSize {
        end := min(i+batchSize, len(data))
        batch := data[i:end]
        
        log.Debug("Processing batch", "batch_start", i, "batch_size", len(batch))
        
        for _, item := range batch {
            processItem(item)
        }
    }
    
    log.Info("Finished processing dataset")
}
```

## Summary

Following these best practices will help you build robust, performant, and maintainable logging systems:

1. **Configure appropriately** for each environment
2. **Use structured logging** with consistent field names
3. **Protect sensitive data** through masking and sanitization
4. **Optimize performance** with conditional and batch logging
5. **Monitor and alert** on log metrics and error rates
6. **Test your logging** to ensure it works correctly
7. **Plan for deployment** with proper rotation and graceful shutdown

Remember that logging is a balance between observability and performance. Always consider the impact of your logging decisions on both system performance and operational visibility. 