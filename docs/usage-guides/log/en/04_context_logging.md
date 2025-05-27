# Context Logging

This document introduces how to use context-aware logging functionality, which is an essential tool for tracking requests and operations in modern applications.

## Context Logging Overview

Context logging allows you to attach relevant information to Go's `context.Context`, which is then automatically included throughout the request lifecycle.

### Key Benefits

1. **Request Tracing** - Track the complete lifecycle of individual requests
2. **Automatic Propagation** - Context information automatically propagates to child functions
3. **Structured Data** - Maintains structured and consistent logging
4. **Performance Optimization** - Avoids repeatedly passing the same log fields

## Basic Usage

### WithValues - Adding Context Fields

Use `log.WithValues()` to add key-value pairs to the context:

```go
package main

import (
    "context"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

func main() {
    log.Init(nil)
    
    // Create context with context fields
    ctx := context.Background()
    ctx = log.WithValues(ctx, "request_id", "req-123", "user_id", 456)
    
    // Use context for logging
    log.InfoContext(ctx, "Processing user request")
    // Output: {"level":"info","timestamp":"2024-01-15T10:30:45.123Z","message":"Processing user request","request_id":"req-123","user_id":456}
    
    processRequest(ctx)
}

func processRequest(ctx context.Context) {
    // Context fields are automatically included in logs
    log.InfoContext(ctx, "Starting request processing")
    
    // Can add more fields
    ctx = log.WithValues(ctx, "step", "validation")
    log.InfoContext(ctx, "Validating request parameters")
    
    validateRequest(ctx)
}

func validateRequest(ctx context.Context) {
    // All previous context fields are included
    log.InfoContext(ctx, "Parameter validation completed")
    // Output includes: request_id, user_id, step
}
```

### Context Logging Methods

The log module provides context-aware methods for all levels:

```go
// Basic context logging methods
log.DebugContext(ctx, "Debug information")
log.InfoContext(ctx, "Info message")
log.WarnContext(ctx, "Warning message")
log.ErrorContext(ctx, "Error message")

// Context logging methods with additional fields
log.DebugwContext(ctx, "Debug information", "key", "value")
log.InfowContext(ctx, "Info message", "key", "value")
log.WarnwContext(ctx, "Warning message", "key", "value")
log.ErrorwContext(ctx, "Error message", "key", "value")
```

## Real-World Use Cases

### HTTP Request Tracing

```go
package main

import (
    "context"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
    "github.com/google/uuid"
)

// Request ID middleware
func RequestIDMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }
        
        // Add request ID to context
        ctx := log.WithValues(c.Request.Context(), "request_id", requestID)
        c.Request = c.Request.WithContext(ctx)
        
        // Set response header
        c.Header("X-Request-ID", requestID)
        
        log.InfoContext(ctx, "Received HTTP request",
            "method", c.Request.Method,
            "path", c.Request.URL.Path,
            "remote_addr", c.ClientIP(),
        )
        
        c.Next()
        
        log.InfoContext(ctx, "HTTP request completed",
            "status", c.Writer.Status(),
            "response_size", c.Writer.Size(),
        )
    }
}

// User authentication middleware
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Assume getting user info from token
        userID := getUserIDFromToken(c.GetHeader("Authorization"))
        if userID != "" {
            // Add user info to context
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
    
    // Add operation-specific context
    ctx = log.WithValues(ctx, "operation", "get_user", "target_user_id", userID)
    
    log.InfoContext(ctx, "Starting to get user information")
    
    user, err := getUserFromDatabase(ctx, userID)
    if err != nil {
        log.ErrorContext(ctx, "Failed to get user", "error", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }
    
    log.InfoContext(ctx, "Successfully retrieved user information")
    c.JSON(http.StatusOK, user)
}

func createUserHandler(c *gin.Context) {
    ctx := c.Request.Context()
    ctx = log.WithValues(ctx, "operation", "create_user")
    
    log.InfoContext(ctx, "Starting to create user")
    
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        log.WarnContext(ctx, "Invalid request parameters", "error", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }
    
    createdUser, err := createUserInDatabase(ctx, &user)
    if err != nil {
        log.ErrorContext(ctx, "Failed to create user", "error", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
        return
    }
    
    log.InfoContext(ctx, "Successfully created user", "new_user_id", createdUser.ID)
    c.JSON(http.StatusCreated, createdUser)
}
```

### Database Operation Tracing

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
    // Add database operation context
    ctx = log.WithValues(ctx,
        "operation", "db_query",
        "table", "users",
        "query_type", "select",
    )
    
    log.InfoContext(ctx, "Starting database query")
    start := time.Now()
    
    query := "SELECT id, name, email FROM users WHERE id = ?"
    row := r.db.QueryRowContext(ctx, query, userID)
    
    var user User
    err := row.Scan(&user.ID, &user.Name, &user.Email)
    
    duration := time.Since(start)
    
    if err != nil {
        if err == sql.ErrNoRows {
            log.WarnContext(ctx, "User not found",
                "duration", duration,
                "query", query,
            )
            return nil, ErrUserNotFound
        }
        
        log.ErrorContext(ctx, "Database query failed",
            "error", err,
            "duration", duration,
            "query", query,
        )
        return nil, err
    }
    
    log.InfoContext(ctx, "Database query successful",
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
    
    log.InfoContext(ctx, "Starting to insert user record")
    start := time.Now()
    
    query := "INSERT INTO users (name, email) VALUES (?, ?)"
    result, err := r.db.ExecContext(ctx, query, user.Name, user.Email)
    
    duration := time.Since(start)
    
    if err != nil {
        log.ErrorContext(ctx, "Failed to insert user record",
            "error", err,
            "duration", duration,
        )
        return err
    }
    
    rowsAffected, _ := result.RowsAffected()
    lastInsertID, _ := result.LastInsertId()
    
    log.InfoContext(ctx, "Successfully inserted user record",
        "duration", duration,
        "rows_affected", rowsAffected,
        "last_insert_id", lastInsertID,
    )
    
    return nil
}
```

### Microservice Call Tracing

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
    // Add service call context
    ctx = log.WithValues(ctx,
        "operation", "service_call",
        "service", "user-service",
        "method", "GET",
        "endpoint", "/users/"+userID,
    )
    
    log.InfoContext(ctx, "Starting user service call")
    start := time.Now()
    
    req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/users/"+userID, nil)
    if err != nil {
        log.ErrorContext(ctx, "Failed to create request", "error", err)
        return nil, err
    }
    
    // Propagate request ID (if exists)
    if requestID := getRequestIDFromContext(ctx); requestID != "" {
        req.Header.Set("X-Request-ID", requestID)
    }
    
    resp, err := c.client.Do(req)
    duration := time.Since(start)
    
    if err != nil {
        log.ErrorContext(ctx, "Service call failed",
            "error", err,
            "duration", duration,
        )
        return nil, err
    }
    defer resp.Body.Close()
    
    log.InfoContext(ctx, "Service call completed",
        "status_code", resp.StatusCode,
        "duration", duration,
        "response_size", resp.ContentLength,
    )
    
    if resp.StatusCode != http.StatusOK {
        log.WarnContext(ctx, "Service returned non-success status code",
            "status_code", resp.StatusCode,
        )
        return nil, ErrServiceUnavailable
    }
    
    var profile UserProfile
    if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
        log.ErrorContext(ctx, "Failed to parse response", "error", err)
        return nil, err
    }
    
    return &profile, nil
}

// Extract request ID from context
func getRequestIDFromContext(ctx context.Context) string {
    // This needs to be adjusted based on actual context implementation
    // This is an example implementation
    if values := log.GetValuesFromContext(ctx); values != nil {
        if requestID, ok := values["request_id"].(string); ok {
            return requestID
        }
    }
    return ""
}
```

## Advanced Usage

### Nested Context

```go
func processOrder(ctx context.Context, orderID string) error {
    // Add order-related context
    ctx = log.WithValues(ctx, "order_id", orderID, "operation", "process_order")
    
    log.InfoContext(ctx, "Starting order processing")
    
    // Validate order
    if err := validateOrder(ctx, orderID); err != nil {
        return err
    }
    
    // Process payment
    if err := processPayment(ctx, orderID); err != nil {
        return err
    }
    
    // Update inventory
    if err := updateInventory(ctx, orderID); err != nil {
        return err
    }
    
    log.InfoContext(ctx, "Order processing completed")
    return nil
}

func validateOrder(ctx context.Context, orderID string) error {
    // Add validation step context
    ctx = log.WithValues(ctx, "step", "validation")
    
    log.InfoContext(ctx, "Starting order validation")
    
    // Validation logic...
    
    log.InfoContext(ctx, "Order validation completed")
    return nil
}

func processPayment(ctx context.Context, orderID string) error {
    // Add payment step context
    ctx = log.WithValues(ctx, "step", "payment")
    
    log.InfoContext(ctx, "Starting payment processing")
    
    // Payment logic...
    
    log.InfoContext(ctx, "Payment processing completed")
    return nil
}
```

### Conditional Context

```go
func handleRequest(ctx context.Context, req *Request) {
    // Base context
    ctx = log.WithValues(ctx, "operation", "handle_request")
    
    // Add different context information based on conditions
    if req.UserID != "" {
        ctx = log.WithValues(ctx, "user_id", req.UserID)
    }
    
    if req.IsAdmin {
        ctx = log.WithValues(ctx, "admin", true)
    }
    
    if req.ClientIP != "" {
        ctx = log.WithValues(ctx, "client_ip", req.ClientIP)
    }
    
    log.InfoContext(ctx, "Processing request")
    
    // Processing logic...
}
```

### Error Context

```go
func processWithErrorContext(ctx context.Context) error {
    ctx = log.WithValues(ctx, "operation", "critical_process")
    
    log.InfoContext(ctx, "Starting critical processing")
    
    if err := step1(ctx); err != nil {
        // Add error-related context
        ctx = log.WithValues(ctx, "failed_step", "step1", "error_type", "validation")
        log.ErrorContext(ctx, "Step 1 failed", "error", err)
        return err
    }
    
    if err := step2(ctx); err != nil {
        ctx = log.WithValues(ctx, "failed_step", "step2", "error_type", "database")
        log.ErrorContext(ctx, "Step 2 failed", "error", err)
        return err
    }
    
    log.InfoContext(ctx, "Critical processing completed")
    return nil
}
```

## Performance Considerations

### Context Field Count

```go
// Good practice: moderate context fields
ctx = log.WithValues(ctx,
    "request_id", requestID,
    "user_id", userID,
    "operation", "get_user",
)

// Avoid: too many context fields
ctx = log.WithValues(ctx,
    "field1", "value1",
    "field2", "value2",
    // ... 20+ fields
    "field25", "value25",
)
```

### Context Reuse

```go
// Good practice: reuse context
func processMultipleItems(ctx context.Context, items []Item) {
    baseCtx := log.WithValues(ctx, "operation", "process_items", "total_items", len(items))
    
    for i, item := range items {
        // Create specific context for each item
        itemCtx := log.WithValues(baseCtx, "item_index", i, "item_id", item.ID)
        processItem(itemCtx, item)
    }
}

// Avoid: repeatedly creating the same context
func processMultipleItemsBad(ctx context.Context, items []Item) {
    for i, item := range items {
        // Recreating complete context every time
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

## Best Practices

### 1. Consistent Field Naming

```go
// Use consistent field names
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

### 2. Layered Context

```go
// Application-level context
func withAppContext(ctx context.Context) context.Context {
    return log.WithValues(ctx,
        "app", "user-service",
        "version", "1.0.0",
    )
}

// Request-level context
func withRequestContext(ctx context.Context, requestID string) context.Context {
    return log.WithValues(ctx,
        "request_id", requestID,
        "timestamp", time.Now().Unix(),
    )
}

// Operation-level context
func withOperationContext(ctx context.Context, operation string) context.Context {
    return log.WithValues(ctx,
        "operation", operation,
    )
}
```

### 3. Sensitive Information Handling

```go
// Avoid logging sensitive information
func loginUser(ctx context.Context, username, password string) error {
    // Don't log passwords
    ctx = log.WithValues(ctx,
        "operation", "login",
        "username", username,
        // "password", password,  // Never do this
    )
    
    log.InfoContext(ctx, "User login attempt")
    
    // Login logic...
    
    return nil
}

// For sensitive fields, use masking or hashing
func processPayment(ctx context.Context, cardNumber string) error {
    maskedCard := maskCardNumber(cardNumber)
    
    ctx = log.WithValues(ctx,
        "operation", "process_payment",
        "card_number", maskedCard,  // Use masked version
    )
    
    log.InfoContext(ctx, "Processing payment")
    
    // Payment logic...
    
    return nil
}

func maskCardNumber(cardNumber string) string {
    if len(cardNumber) < 4 {
        return "****"
    }
    return "****-****-****-" + cardNumber[len(cardNumber)-4:]
}
```

## Debugging and Troubleshooting

### Viewing Context Fields

```go
// Debug: print all fields in current context
func debugContext(ctx context.Context) {
    if values := log.GetValuesFromContext(ctx); values != nil {
        log.InfoContext(ctx, "Current context fields", "context_values", values)
    }
}
```

### Context Propagation Verification

```go
func verifyContextPropagation(ctx context.Context) {
    // Verify that key fields exist
    values := log.GetValuesFromContext(ctx)
    if values == nil {
        log.Warn("No log fields in context")
        return
    }
    
    if _, ok := values["request_id"]; !ok {
        log.Warn("Missing request_id in context")
    }
    
    if _, ok := values["user_id"]; !ok {
        log.Warn("Missing user_id in context")
    }
}
```

## Next Steps

- [Performance](05_performance.md) - Optimize logging performance
- [Best Practices](06_best_practices.md) - Production-ready patterns 