# 上下文日志记录

本文档介绍如何使用上下文感知的日志记录功能，这是现代应用程序中跟踪请求和操作的重要工具。

## 上下文日志概述

上下文日志记录允许您将相关信息附加到 Go 的 `context.Context` 中，然后在整个请求生命周期中自动包含这些信息。

### 主要优势

1. **请求跟踪** - 跟踪单个请求的完整生命周期
2. **自动传播** - 上下文信息自动传播到子函数
3. **结构化数据** - 保持日志的结构化和一致性
4. **性能优化** - 避免重复传递相同的日志字段

## 基本用法

### WithValues - 添加上下文字段

使用 `log.WithValues()` 向上下文添加键值对：

```go
package main

import (
    "context"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

func main() {
    log.Init(nil)
    
    // 创建带有上下文字段的上下文
    ctx := context.Background()
    ctx = log.WithValues(ctx, "request_id", "req-123", "user_id", 456)
    
    // 使用上下文记录日志
    log.InfoContext(ctx, "处理用户请求")
    // 输出：{"level":"info","timestamp":"2024-01-15T10:30:45.123Z","message":"处理用户请求","request_id":"req-123","user_id":456}
    
    processRequest(ctx)
}

func processRequest(ctx context.Context) {
    // 上下文字段会自动包含在日志中
    log.InfoContext(ctx, "开始处理请求")
    
    // 可以添加更多字段
    ctx = log.WithValues(ctx, "step", "validation")
    log.InfoContext(ctx, "验证请求参数")
    
    validateRequest(ctx)
}

func validateRequest(ctx context.Context) {
    // 所有之前的上下文字段都会包含
    log.InfoContext(ctx, "参数验证完成")
    // 输出包含：request_id, user_id, step
}
```

### 上下文日志方法

日志模块提供了所有级别的上下文感知方法：

```go
// 基本上下文日志方法
log.DebugContext(ctx, "调试信息")
log.InfoContext(ctx, "信息消息")
log.WarnContext(ctx, "警告消息")
log.ErrorContext(ctx, "错误消息")

// 带额外字段的上下文日志方法
log.DebugwContext(ctx, "调试信息", "key", "value")
log.InfowContext(ctx, "信息消息", "key", "value")
log.WarnwContext(ctx, "警告消息", "key", "value")
log.ErrorwContext(ctx, "错误消息", "key", "value")
```

## 实际应用场景

### HTTP 请求跟踪

```go
package main

import (
    "context"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
    "github.com/google/uuid"
)

// 请求 ID 中间件
func RequestIDMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        
        // 将请求 ID 添加到上下文
        ctx := log.WithValues(c.Request.Context(), "request_id", requestID)
        c.Request = c.Request.WithContext(ctx)
        
        // 设置响应头
        c.Header("X-Request-ID", requestID)
        
        log.InfoContext(ctx, "收到 HTTP 请求",
            "method", c.Request.Method,
            "path", c.Request.URL.Path,
            "remote_addr", c.ClientIP(),
        )
        
        c.Next()
        
        log.InfoContext(ctx, "HTTP 请求完成",
            "status", c.Writer.Status(),
            "response_size", c.Writer.Size(),
        )
    }
}

// 用户认证中间件
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 假设从 token 中获取用户信息
        userID := getUserIDFromToken(c.GetHeader("Authorization"))
        if userID != "" {
            // 添加用户信息到上下文
            ctx := log.WithValues(c.Request.Context(), "user_id", userID)
            c.Request = c.Request.WithContext(ctx)
        }
        
        c.Next()
    }
}

func main() {
    log.Init(nil)
    
    r := gin.New()
    r.Use(RequestIDMiddleware())
    r.Use(AuthMiddleware())
    
    r.GET("/users/:id", getUserHandler)
    r.POST("/users", createUserHandler)
    
    r.Run(":8080")
}

func getUserHandler(c *gin.Context) {
    ctx := c.Request.Context()
    userID := c.Param("id")
    
    // 添加特定操作的上下文
    ctx = log.WithValues(ctx, "operation", "get_user", "target_user_id", userID)
    
    log.InfoContext(ctx, "开始获取用户信息")
    
    user, err := getUserFromDatabase(ctx, userID)
    if err != nil {
        log.ErrorContext(ctx, "获取用户失败", "error", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }
    
    log.InfoContext(ctx, "成功获取用户信息")
    c.JSON(http.StatusOK, user)
}

func createUserHandler(c *gin.Context) {
    ctx := c.Request.Context()
    ctx = log.WithValues(ctx, "operation", "create_user")
    
    log.InfoContext(ctx, "开始创建用户")
    
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        log.WarnContext(ctx, "请求参数无效", "error", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }
    
    createdUser, err := createUserInDatabase(ctx, &user)
    if err != nil {
        log.ErrorContext(ctx, "创建用户失败", "error", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }
    
    log.InfoContext(ctx, "成功创建用户", "new_user_id", createdUser.ID)
    c.JSON(http.StatusCreated, createdUser)
}
```

### 数据库操作跟踪

```go
package main

import (
    "context"
    "database/sql"
    "time"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

type UserRepository struct {
    db *sql.DB
}

func (r *UserRepository) GetUser(ctx context.Context, userID string) (*User, error) {
    // 添加数据库操作上下文
    ctx = log.WithValues(ctx,
        "operation", "db_query",
        "table", "users",
        "query_type", "select",
    )
    
    log.InfoContext(ctx, "开始数据库查询")
    start := time.Now()
    
    query := "SELECT id, name, email FROM users WHERE id = ?"
    row := r.db.QueryRowContext(ctx, query, userID)
    
    var user User
    err := row.Scan(&user.ID, &user.Name, &user.Email)
    
    duration := time.Since(start)
    
    if err != nil {
        if err == sql.ErrNoRows {
            log.WarnContext(ctx, "用户不存在",
                "duration", duration,
                "query", query,
            )
            return nil, ErrUserNotFound
        }
        
        log.ErrorContext(ctx, "数据库查询失败",
            "error", err,
            "duration", duration,
            "query", query,
        )
        return nil, err
    }
    
    log.InfoContext(ctx, "数据库查询成功",
        "duration", duration,
        "rows_affected", 1,
    )
    
    return &user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *User) error {
    ctx = log.WithValues(ctx,
        "operation", "db_insert",
        "table", "users",
    )
    
    log.InfoContext(ctx, "开始插入用户记录")
    start := time.Now()
    
    query := "INSERT INTO users (name, email) VALUES (?, ?)"
    result, err := r.db.ExecContext(ctx, query, user.Name, user.Email)
    
    duration := time.Since(start)
    
    if err != nil {
        log.ErrorContext(ctx, "插入用户记录失败",
            "error", err,
            "duration", duration,
        )
        return err
    }
    
    rowsAffected, _ := result.RowsAffected()
    lastInsertID, _ := result.LastInsertId()
    
    log.InfoContext(ctx, "成功插入用户记录",
        "duration", duration,
        "rows_affected", rowsAffected,
        "last_insert_id", lastInsertID,
    )
    
    return nil
}
```

### 微服务调用跟踪

```go
package main

import (
    "context"
    "net/http"
    "time"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

type ServiceClient struct {
    baseURL string
    client  *http.Client
}

func (c *ServiceClient) CallUserService(ctx context.Context, userID string) (*UserProfile, error) {
    // 添加服务调用上下文
    ctx = log.WithValues(ctx,
        "operation", "service_call",
        "service", "user-service",
        "method", "GET",
        "endpoint", "/users/"+userID,
    )
    
    log.InfoContext(ctx, "开始调用用户服务")
    start := time.Now()
    
    req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/users/"+userID, nil)
    if err != nil {
        log.ErrorContext(ctx, "创建请求失败", "error", err)
        return nil, err
    }
    
    // 传播请求 ID（如果存在）
    if requestID := getRequestIDFromContext(ctx); requestID != "" {
        req.Header.Set("X-Request-ID", requestID)
    }
    
    resp, err := c.client.Do(req)
    duration := time.Since(start)
    
    if err != nil {
        log.ErrorContext(ctx, "服务调用失败",
            "error", err,
            "duration", duration,
        )
        return nil, err
    }
    defer resp.Body.Close()
    
    log.InfoContext(ctx, "服务调用完成",
        "status_code", resp.StatusCode,
        "duration", duration,
        "response_size", resp.ContentLength,
    )
    
    if resp.StatusCode != http.StatusOK {
        log.WarnContext(ctx, "服务返回非成功状态码",
            "status_code", resp.StatusCode,
        )
        return nil, ErrServiceUnavailable
    }
    
    var profile UserProfile
    if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
        log.ErrorContext(ctx, "解析响应失败", "error", err)
        return nil, err
    }
    
    return &profile, nil
}

// 从上下文中提取请求 ID
func getRequestIDFromContext(ctx context.Context) string {
    // 这需要根据实际的上下文实现来调整
    // 这里是一个示例实现
    if values := log.GetValuesFromContext(ctx); values != nil {
        if requestID, ok := values["request_id"].(string); ok {
            return requestID
        }
    }
    return ""
}
```

## 高级用法

### 嵌套上下文

```go
func processOrder(ctx context.Context, orderID string) error {
    // 添加订单相关上下文
    ctx = log.WithValues(ctx, "order_id", orderID, "operation", "process_order")
    
    log.InfoContext(ctx, "开始处理订单")
    
    // 验证订单
    if err := validateOrder(ctx, orderID); err != nil {
        return err
    }
    
    // 处理支付
    if err := processPayment(ctx, orderID); err != nil {
        return err
    }
    
    // 更新库存
    if err := updateInventory(ctx, orderID); err != nil {
        return err
    }
    
    log.InfoContext(ctx, "订单处理完成")
    return nil
}

func validateOrder(ctx context.Context, orderID string) error {
    // 添加验证步骤的上下文
    ctx = log.WithValues(ctx, "step", "validation")
    
    log.InfoContext(ctx, "开始验证订单")
    
    // 验证逻辑...
    
    log.InfoContext(ctx, "订单验证完成")
    return nil
}

func processPayment(ctx context.Context, orderID string) error {
    // 添加支付步骤的上下文
    ctx = log.WithValues(ctx, "step", "payment")
    
    log.InfoContext(ctx, "开始处理支付")
    
    // 支付逻辑...
    
    log.InfoContext(ctx, "支付处理完成")
    return nil
}
```

### 条件上下文

```go
func handleRequest(ctx context.Context, req *Request) {
    // 基础上下文
    ctx = log.WithValues(ctx, "operation", "handle_request")
    
    // 根据条件添加不同的上下文信息
    if req.UserID != "" {
        ctx = log.WithValues(ctx, "user_id", req.UserID)
    }
    
    if req.IsAdmin {
        ctx = log.WithValues(ctx, "admin", true)
    }
    
    if req.ClientIP != "" {
        ctx = log.WithValues(ctx, "client_ip", req.ClientIP)
    }
    
    log.InfoContext(ctx, "处理请求")
    
    // 处理逻辑...
}
```

### 错误上下文

```go
func processWithErrorContext(ctx context.Context) error {
    ctx = log.WithValues(ctx, "operation", "critical_process")
    
    log.InfoContext(ctx, "开始关键处理")
    
    if err := step1(ctx); err != nil {
        // 添加错误相关的上下文
        ctx = log.WithValues(ctx, "failed_step", "step1", "error_type", "validation")
        log.ErrorContext(ctx, "步骤1失败", "error", err)
        return err
    }
    
    if err := step2(ctx); err != nil {
        ctx = log.WithValues(ctx, "failed_step", "step2", "error_type", "database")
        log.ErrorContext(ctx, "步骤2失败", "error", err)
        return err
    }
    
    log.InfoContext(ctx, "关键处理完成")
    return nil
}
```

## 性能考虑

### 上下文字段数量

```go
// 好的做法：适量的上下文字段
ctx = log.WithValues(ctx,
    "request_id", requestID,
    "user_id", userID,
    "operation", "get_user",
)

// 避免：过多的上下文字段
ctx = log.WithValues(ctx,
    "field1", "value1",
    "field2", "value2",
    // ... 20+ 个字段
    "field25", "value25",
)
```

### 上下文重用

```go
// 好的做法：重用上下文
func processMultipleItems(ctx context.Context, items []Item) {
    baseCtx := log.WithValues(ctx, "operation", "process_items", "total_items", len(items))
    
    for i, item := range items {
        // 为每个项目创建特定的上下文
        itemCtx := log.WithValues(baseCtx, "item_index", i, "item_id", item.ID)
        processItem(itemCtx, item)
    }
}

// 避免：重复创建相同的上下文
func processMultipleItemsBad(ctx context.Context, items []Item) {
    for i, item := range items {
        // 每次都重新创建完整的上下文
        itemCtx := log.WithValues(ctx,
            "operation", "process_items",
            "total_items", len(items),
            "item_index", i,
            "item_id", item.ID,
        )
        processItem(itemCtx, item)
    }
}
```

## 最佳实践

### 1. 一致的字段命名

```go
// 使用一致的字段名称
const (
    FieldRequestID = "request_id"
    FieldUserID    = "user_id"
    FieldOperation = "operation"
    FieldDuration  = "duration"
    FieldError     = "error"
)

ctx = log.WithValues(ctx,
    FieldRequestID, requestID,
    FieldUserID, userID,
    FieldOperation, "create_user",
)
```

### 2. 分层上下文

```go
// 应用级别上下文
func withAppContext(ctx context.Context) context.Context {
    return log.WithValues(ctx,
        "app", "user-service",
        "version", "1.0.0",
    )
}

// 请求级别上下文
func withRequestContext(ctx context.Context, requestID string) context.Context {
    return log.WithValues(ctx,
        "request_id", requestID,
        "timestamp", time.Now().Unix(),
    )
}

// 操作级别上下文
func withOperationContext(ctx context.Context, operation string) context.Context {
    return log.WithValues(ctx,
        "operation", operation,
    )
}
```

### 3. 敏感信息处理

```go
// 避免记录敏感信息
func loginUser(ctx context.Context, username, password string) error {
    // 不要记录密码
    ctx = log.WithValues(ctx,
        "operation", "login",
        "username", username,
        // "password", password,  // 永远不要这样做
    )
    
    log.InfoContext(ctx, "用户登录尝试")
    
    // 登录逻辑...
    
    return nil
}

// 对于敏感字段，使用掩码或哈希
func processPayment(ctx context.Context, cardNumber string) error {
    maskedCard := maskCardNumber(cardNumber)
    
    ctx = log.WithValues(ctx,
        "operation", "process_payment",
        "card_number", maskedCard,  // 使用掩码版本
    )
    
    log.InfoContext(ctx, "处理支付")
    
    // 支付逻辑...
    
    return nil
}

func maskCardNumber(cardNumber string) string {
    if len(cardNumber) < 4 {
        return "****"
    }
    return "****-****-****-" + cardNumber[len(cardNumber)-4:]
}
```

## 调试和故障排除

### 查看上下文字段

```go
// 调试：打印当前上下文中的所有字段
func debugContext(ctx context.Context) {
    if values := log.GetValuesFromContext(ctx); values != nil {
        log.InfoContext(ctx, "当前上下文字段", "context_values", values)
    }
}
```

### 上下文传播验证

```go
func verifyContextPropagation(ctx context.Context) {
    // 验证关键字段是否存在
    values := log.GetValuesFromContext(ctx)
    if values == nil {
        log.Warn("上下文中没有日志字段")
        return
    }
    
    if _, ok := values["request_id"]; !ok {
        log.Warn("上下文中缺少 request_id")
    }
    
    if _, ok := values["user_id"]; !ok {
        log.Warn("上下文中缺少 user_id")
    }
}
```

## 下一步

- [性能优化](05_performance.md) - 优化日志性能
- [最佳实践](06_best_practices.md) - 生产就绪模式 