# 最佳实践

本文档汇总了在生产环境中使用日志模块的最佳实践，帮助您构建可靠、高性能和可维护的日志系统。

## 核心原则

### 1. 结构化日志优先

始终使用结构化日志而不是字符串格式化：

```go
// ❌ 避免：字符串格式化
log.Info(fmt.Sprintf("用户 %s 登录失败，尝试次数: %d", username, attempts))

// ✅ 推荐：结构化字段
log.Infow("用户登录失败",
    "username", username,
    "attempts", attempts,
    "ip", clientIP,
    "user_agent", userAgent,
)
```

### 2. 一致的字段命名

建立并遵循一致的字段命名约定：

```go
// 定义标准字段名称
const (
    FieldRequestID   = "request_id"
    FieldUserID      = "user_id"
    FieldOperation   = "operation"
    FieldDuration    = "duration"
    FieldError       = "error"
    FieldStatusCode  = "status_code"
    FieldMethod      = "method"
    FieldPath        = "path"
    FieldIP          = "ip"
    FieldUserAgent   = "user_agent"
)

// 使用标准字段名称
log.Infow("HTTP 请求",
    FieldRequestID, requestID,
    FieldMethod, r.Method,
    FieldPath, r.URL.Path,
    FieldIP, getClientIP(r),
    FieldUserAgent, r.UserAgent(),
)
```

### 3. 适当的日志级别

为不同类型的事件选择合适的日志级别：

```go
// DEBUG: 详细的调试信息
log.Debug("缓存查找", "key", cacheKey, "hit", hit)

// INFO: 一般信息事件
log.Info("用户登录成功", "user_id", userID)

// WARN: 警告但不影响功能
log.Warn("API 响应缓慢", "duration", duration, "threshold", threshold)

// ERROR: 错误但应用程序可以继续
log.Error("数据库连接失败", "error", err, "retry_count", retryCount)

// FATAL: 严重错误，应用程序无法继续
log.Fatal("配置文件无法读取", "file", configFile, "error", err)
```

## 环境配置策略

### 开发环境

```go
func developmentLogConfig() *log.Options {
    return &log.Options{
        Level:       "debug",
        Format:      "text",
        OutputPaths: []string{"stdout"},
        EnableColor: true,
        
        // 保留调试信息
        DisableCaller:     false,
        DisableStacktrace: false,
        StacktraceLevel:   "warn",
    }
}
```

### 测试环境

```go
func testLogConfig() *log.Options {
    return &log.Options{
        Level:       "warn",
        Format:      "text",
        OutputPaths: []string{"stdout"},
        EnableColor: false,
        
        // 提高测试性能
        DisableCaller:     true,
        DisableStacktrace: true,
    }
}
```

### 生产环境

```go
func productionLogConfig() *log.Options {
    return &log.Options{
        Level:  "info",
        Format: "json",
        
        OutputPaths: []string{
            "/var/log/app/app.log",
            "stdout", // 用于容器日志收集
        },
        ErrorOutputPaths: []string{
            "/var/log/app/error.log",
            "stderr",
        },
        
        // 平衡性能和可观测性
        DisableCaller:     false,
        DisableStacktrace: false,
        StacktraceLevel:   "error",
        
        // 日志轮转配置
        LogRotateMaxSize:    100, // 100MB
        LogRotateMaxBackups: 10,
        LogRotateMaxAge:     30, // 30天
        LogRotateCompress:   true,
    }
}
```

## 上下文使用模式

### HTTP 请求跟踪

```go
// 中间件：添加请求上下文
func LoggingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        
        // 生成或获取请求 ID
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = generateRequestID()
        }
        
        // 创建请求上下文
        ctx := log.WithValues(c.Request.Context(),
            FieldRequestID, requestID,
            FieldMethod, c.Request.Method,
            FieldPath, c.Request.URL.Path,
            FieldIP, c.ClientIP(),
            FieldUserAgent, c.Request.UserAgent(),
        )
        
        c.Request = c.Request.WithContext(ctx)
        c.Header("X-Request-ID", requestID)
        
        log.InfoContext(ctx, "HTTP 请求开始")
        
        c.Next()
        
        // 记录请求完成
        duration := time.Since(start)
        log.InfoContext(ctx, "HTTP 请求完成",
            FieldStatusCode, c.Writer.Status(),
            FieldDuration, duration,
            "response_size", c.Writer.Size(),
        )
    }
}
```

### 业务操作跟踪

```go
func ProcessOrder(ctx context.Context, orderID string) error {
    // 添加操作上下文
    ctx = log.WithValues(ctx,
        "order_id", orderID,
        FieldOperation, "process_order",
    )
    
    log.InfoContext(ctx, "开始处理订单")
    
    // 验证订单
    if err := validateOrder(ctx, orderID); err != nil {
        log.ErrorContext(ctx, "订单验证失败", FieldError, err)
        return err
    }
    
    // 处理支付
    if err := processPayment(ctx, orderID); err != nil {
        log.ErrorContext(ctx, "支付处理失败", FieldError, err)
        return err
    }
    
    // 更新库存
    if err := updateInventory(ctx, orderID); err != nil {
        log.ErrorContext(ctx, "库存更新失败", FieldError, err)
        return err
    }
    
    log.InfoContext(ctx, "订单处理完成")
    return nil
}

func validateOrder(ctx context.Context, orderID string) error {
    ctx = log.WithValues(ctx, "step", "validation")
    
    log.InfoContext(ctx, "开始验证订单")
    
    // 验证逻辑...
    
    log.InfoContext(ctx, "订单验证完成")
    return nil
}
```

## 错误处理模式

### 错误日志记录

```go
// 标准错误记录模式
func handleDatabaseError(ctx context.Context, operation string, err error) error {
    // 记录详细错误信息
    log.ErrorContext(ctx, "数据库操作失败",
        FieldOperation, operation,
        FieldError, err,
        "error_type", getErrorType(err),
        "retry_count", getRetryCount(ctx),
    )
    
    // 返回用户友好的错误
    return errors.New("internal server error")
}

// 错误分类
func getErrorType(err error) string {
    switch {
    case errors.Is(err, sql.ErrNoRows):
        return "not_found"
    case errors.Is(err, context.DeadlineExceeded):
        return "timeout"
    case isConnectionError(err):
        return "connection"
    default:
        return "unknown"
    }
}
```

### 错误恢复记录

```go
func recoverMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                ctx := c.Request.Context()
                
                log.ErrorContext(ctx, "应用程序恐慌",
                    "panic", err,
                    "stack", string(debug.Stack()),
                )
                
                c.JSON(http.StatusInternalServerError, gin.H{
                    "error": "Internal server error",
                })
                c.Abort()
            }
        }()
        
        c.Next()
    }
}
```

## 性能最佳实践

### 条件日志记录

```go
// 对于昂贵的操作，使用条件日志记录
func processLargeDataset(ctx context.Context, data []byte) {
    if log.IsDebugEnabled() {
        log.DebugContext(ctx, "处理大数据集",
            "size", len(data),
            "preview", string(data[:min(100, len(data))]),
        )
    }
    
    // 处理逻辑...
}

// 使用延迟计算避免不必要的序列化
func logUserDetails(ctx context.Context, user *User) {
    log.InfoContext(ctx, "用户操作",
        "user_summary", func() interface{} {
            return map[string]interface{}{
                "id":    user.ID,
                "name":  user.Name,
                "roles": user.Roles,
            }
        },
    )
}
```

### 批量操作日志

```go
// 对于批量操作，记录摘要而不是每个项目
func processBatch(ctx context.Context, items []Item) error {
    ctx = log.WithValues(ctx,
        FieldOperation, "batch_process",
        "total_items", len(items),
    )
    
    log.InfoContext(ctx, "开始批量处理")
    
    var processed, failed int
    var errors []error
    
    for i, item := range items {
        if err := processItem(ctx, item); err != nil {
            failed++
            errors = append(errors, err)
            
            // 仅记录前几个错误的详细信息
            if len(errors) <= 5 {
                log.WarnContext(ctx, "项目处理失败",
                    "item_index", i,
                    "item_id", item.ID,
                    FieldError, err,
                )
            }
        } else {
            processed++
        }
    }
    
    log.InfoContext(ctx, "批量处理完成",
        "processed", processed,
        "failed", failed,
        "success_rate", float64(processed)/float64(len(items)),
    )
    
    if failed > 0 {
        return fmt.Errorf("批量处理部分失败: %d/%d", failed, len(items))
    }
    
    return nil
}
```

## 安全和隐私

### 敏感信息处理

```go
// 敏感信息掩码
func maskSensitiveData(data string, maskChar rune, visibleChars int) string {
    if len(data) <= visibleChars {
        return strings.Repeat(string(maskChar), len(data))
    }
    
    masked := strings.Repeat(string(maskChar), len(data)-visibleChars)
    return masked + data[len(data)-visibleChars:]
}

// 安全的用户信息记录
func logUserAction(ctx context.Context, user *User, action string) {
    log.InfoContext(ctx, "用户操作",
        FieldUserID, user.ID,
        "username", user.Username,
        "action", action,
        // 不记录敏感信息
        // "password", user.Password, // ❌ 永远不要这样做
        // "ssn", user.SSN,           // ❌ 永远不要这样做
    )
}

// 支付信息记录
func logPaymentAttempt(ctx context.Context, payment *Payment) {
    log.InfoContext(ctx, "支付尝试",
        "payment_id", payment.ID,
        "amount", payment.Amount,
        "currency", payment.Currency,
        "card_last_four", payment.CardNumber[len(payment.CardNumber)-4:],
        // "card_number", payment.CardNumber, // ❌ 永远不要记录完整卡号
    )
}
```

### 审计日志

```go
// 审计日志结构
type AuditEvent struct {
    UserID    string    `json:"user_id"`
    Action    string    `json:"action"`
    Resource  string    `json:"resource"`
    Timestamp time.Time `json:"timestamp"`
    IP        string    `json:"ip"`
    UserAgent string    `json:"user_agent"`
    Success   bool      `json:"success"`
    Details   map[string]interface{} `json:"details,omitempty"`
}

func logAuditEvent(ctx context.Context, event AuditEvent) {
    log.InfoContext(ctx, "审计事件",
        "audit_user_id", event.UserID,
        "audit_action", event.Action,
        "audit_resource", event.Resource,
        "audit_success", event.Success,
        "audit_ip", event.IP,
        "audit_details", event.Details,
    )
}

// 使用示例
func deleteUser(ctx context.Context, userID string, operatorID string) error {
    // 执行删除操作
    err := userRepository.Delete(userID)
    
    // 记录审计事件
    logAuditEvent(ctx, AuditEvent{
        UserID:    operatorID,
        Action:    "delete_user",
        Resource:  fmt.Sprintf("user:%s", userID),
        Timestamp: time.Now(),
        IP:        getIPFromContext(ctx),
        UserAgent: getUserAgentFromContext(ctx),
        Success:   err == nil,
        Details: map[string]interface{}{
            "target_user_id": userID,
        },
    })
    
    return err
}
```

## 监控和告警

### 关键指标记录

```go
// 业务指标记录
func recordBusinessMetrics(ctx context.Context, operation string, duration time.Duration, success bool) {
    log.InfoContext(ctx, "业务指标",
        "metric_type", "business_operation",
        FieldOperation, operation,
        FieldDuration, duration,
        "success", success,
        "timestamp", time.Now().Unix(),
    )
}

// 性能指标记录
func recordPerformanceMetrics(ctx context.Context, component string, metrics map[string]interface{}) {
    log.InfoContext(ctx, "性能指标",
        "metric_type", "performance",
        "component", component,
        "metrics", metrics,
        "timestamp", time.Now().Unix(),
    )
}

// 使用示例
func processAPIRequest(ctx context.Context, endpoint string) error {
    start := time.Now()
    
    err := handleRequest(ctx, endpoint)
    
    duration := time.Since(start)
    success := err == nil
    
    // 记录业务指标
    recordBusinessMetrics(ctx, "api_request", duration, success)
    
    // 记录性能指标
    if duration > 1*time.Second {
        recordPerformanceMetrics(ctx, "api", map[string]interface{}{
            "endpoint": endpoint,
            "duration": duration,
            "slow_request": true,
        })
    }
    
    return err
}
```

### 错误率监控

```go
// 错误率跟踪
type ErrorTracker struct {
    mu           sync.RWMutex
    errorCounts  map[string]int64
    totalCounts  map[string]int64
    lastReset    time.Time
    resetInterval time.Duration
}

func NewErrorTracker(resetInterval time.Duration) *ErrorTracker {
    return &ErrorTracker{
        errorCounts:   make(map[string]int64),
        totalCounts:   make(map[string]int64),
        lastReset:     time.Now(),
        resetInterval: resetInterval,
    }
}

func (et *ErrorTracker) Record(operation string, isError bool) {
    et.mu.Lock()
    defer et.mu.Unlock()
    
    // 检查是否需要重置计数器
    if time.Since(et.lastReset) > et.resetInterval {
        et.errorCounts = make(map[string]int64)
        et.totalCounts = make(map[string]int64)
        et.lastReset = time.Now()
    }
    
    et.totalCounts[operation]++
    if isError {
        et.errorCounts[operation]++
    }
    
    // 检查错误率
    errorRate := float64(et.errorCounts[operation]) / float64(et.totalCounts[operation])
    if errorRate > 0.1 && et.totalCounts[operation] > 10 { // 10% 错误率阈值
        log.Warn("高错误率检测",
            FieldOperation, operation,
            "error_rate", errorRate,
            "error_count", et.errorCounts[operation],
            "total_count", et.totalCounts[operation],
        )
    }
}
```

## 日志聚合和分析

### 结构化查询友好

```go
// 设计查询友好的日志结构
func logUserEvent(ctx context.Context, event UserEvent) {
    log.InfoContext(ctx, "用户事件",
        // 标准字段用于过滤和聚合
        "event_type", event.Type,
        "event_category", event.Category,
        FieldUserID, event.UserID,
        "timestamp", event.Timestamp.Unix(),
        
        // 详细信息用于分析
        "event_details", event.Details,
        
        // 维度字段用于分组
        "user_segment", event.UserSegment,
        "feature_flag", event.FeatureFlag,
        "ab_test_group", event.ABTestGroup,
    )
}

// 电商事件示例
func logPurchaseEvent(ctx context.Context, purchase Purchase) {
    log.InfoContext(ctx, "购买事件",
        "event_type", "purchase",
        "event_category", "transaction",
        FieldUserID, purchase.UserID,
        "order_id", purchase.OrderID,
        "amount", purchase.Amount,
        "currency", purchase.Currency,
        "payment_method", purchase.PaymentMethod,
        "product_categories", purchase.ProductCategories,
        "discount_applied", purchase.DiscountApplied,
        "timestamp", purchase.Timestamp.Unix(),
    )
}
```

### 日志采样

```go
// 高频事件采样
type LogSampler struct {
    sampleRate float64
    counter    int64
}

func NewLogSampler(sampleRate float64) *LogSampler {
    return &LogSampler{sampleRate: sampleRate}
}

func (ls *LogSampler) ShouldLog() bool {
    count := atomic.AddInt64(&ls.counter, 1)
    return float64(count%100) < ls.sampleRate*100
}

var debugSampler = NewLogSampler(0.01) // 1% 采样率

func logHighFrequencyEvent(ctx context.Context, event string) {
    if debugSampler.ShouldLog() {
        log.DebugContext(ctx, "高频事件", "event", event)
    }
}
```

## 测试和验证

### 日志测试

```go
// 日志测试辅助函数
func captureLogOutput(t *testing.T, fn func()) []string {
    var buf bytes.Buffer
    
    // 临时重定向日志输出
    opts := &log.Options{
        Level:       "debug",
        Format:      "json",
        OutputPaths: []string{"stdout"},
    }
    
    // 保存原始配置
    originalLogger := log.GetLogger()
    defer log.SetLogger(originalLogger)
    
    // 设置测试配置
    log.Init(opts)
    
    // 执行测试函数
    fn()
    
    // 解析日志输出
    lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
    return lines
}

// 测试示例
func TestUserLoginLogging(t *testing.T) {
    logs := captureLogOutput(t, func() {
        ctx := context.Background()
        loginUser(ctx, "testuser", "password123")
    })
    
    // 验证日志内容
    assert.Contains(t, logs[0], "用户登录尝试")
    assert.Contains(t, logs[0], "testuser")
    assert.NotContains(t, logs[0], "password123") // 确保密码未被记录
}
```

### 日志格式验证

```go
// 验证日志格式的一致性
func validateLogFormat(t *testing.T, logLine string) {
    var logEntry map[string]interface{}
    err := json.Unmarshal([]byte(logLine), &logEntry)
    assert.NoError(t, err, "日志应该是有效的 JSON")
    
    // 验证必需字段
    requiredFields := []string{"level", "timestamp", "message"}
    for _, field := range requiredFields {
        assert.Contains(t, logEntry, field, "日志应该包含字段: %s", field)
    }
    
    // 验证时间戳格式
    timestamp, ok := logEntry["timestamp"].(string)
    assert.True(t, ok, "timestamp 应该是字符串")
    
    _, err = time.Parse(time.RFC3339, timestamp)
    assert.NoError(t, err, "timestamp 应该是有效的 RFC3339 格式")
}
```

## 部署和运维

### 容器化环境

```dockerfile
# Dockerfile 示例
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# 创建日志目录
RUN mkdir -p /var/log/app

COPY --from=builder /app/main .

# 设置环境变量
ENV LOG_LEVEL=info
ENV LOG_FORMAT=json
ENV LOG_OUTPUT=/var/log/app/app.log

CMD ["./main"]
```

```yaml
# docker-compose.yml 示例
version: '3.8'
services:
  app:
    build: .
    environment:
      - LOG_LEVEL=info
      - LOG_FORMAT=json
    volumes:
      - ./logs:/var/log/app
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

### Kubernetes 配置

```yaml
# k8s-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
      - name: app
        image: myapp:latest
        env:
        - name: LOG_LEVEL
          value: "info"
        - name: LOG_FORMAT
          value: "json"
        volumeMounts:
        - name: log-volume
          mountPath: /var/log/app
      volumes:
      - name: log-volume
        emptyDir: {}
```

### 日志收集配置

```yaml
# fluentd 配置示例
apiVersion: v1
kind: ConfigMap
metadata:
  name: fluentd-config
data:
  fluent.conf: |
    <source>
      @type tail
      path /var/log/app/*.log
      pos_file /var/log/fluentd-app.log.pos
      tag app.logs
      format json
      time_key timestamp
      time_format %Y-%m-%dT%H:%M:%S.%LZ
    </source>
    
    <match app.logs>
      @type elasticsearch
      host elasticsearch.logging.svc.cluster.local
      port 9200
      index_name app-logs
      type_name _doc
    </match>
```

## 故障排除清单

### 常见问题检查

1. **日志级别配置**
   ```bash
   # 检查当前日志级别
   curl http://localhost:8080/debug/log/level
   
   # 动态调整日志级别
   curl -X POST http://localhost:8080/debug/log/level -d '{"level":"debug"}'
   ```

2. **日志输出验证**
   ```bash
   # 检查日志文件是否正在写入
   tail -f /var/log/app/app.log
   
   # 检查日志轮转
   ls -la /var/log/app/
   ```

3. **性能影响评估**
   ```bash
   # 监控日志相关的系统调用
   strace -e write -p $(pgrep myapp)
   
   # 检查磁盘 I/O
   iostat -x 1
   ```

### 调试模式

```go
// 调试模式配置
func enableDebugMode() {
    opts := &log.Options{
        Level:             "debug",
        Format:            "text",
        OutputPaths:       []string{"stdout", "/tmp/debug.log"},
        EnableColor:       true,
        DisableCaller:     false,
        DisableStacktrace: false,
        StacktraceLevel:   "debug",
    }
    
    log.Init(opts)
    log.Info("调试模式已启用")
}
```

## 总结

遵循这些最佳实践将帮助您：

1. **提高可观测性** - 通过结构化日志和一致的字段命名
2. **保证性能** - 通过适当的配置和优化技术
3. **确保安全** - 通过正确处理敏感信息
4. **简化运维** - 通过标准化的部署和监控
5. **便于调试** - 通过详细的上下文信息和错误跟踪

记住，好的日志记录是一门艺术，需要在详细程度、性能和可用性之间找到平衡。根据您的具体需求调整这些实践，并持续改进您的日志策略。 