/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * Structured logging example demonstrating JSON format and structured fields.
 */

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

// RequestInfo 请求信息结构
// (RequestInfo represents request information structure)
type RequestInfo struct {
	RequestID    string            `json:"request_id"`
	Method       string            `json:"method"`
	Path         string            `json:"path"`
	UserAgent    string            `json:"user_agent"`
	RemoteIP     string            `json:"remote_ip"`
	Headers      map[string]string `json:"headers"`
	QueryParams  map[string]string `json:"query_params"`
	StartTime    time.Time         `json:"start_time"`
	Duration     time.Duration     `json:"duration,omitempty"`
	StatusCode   int               `json:"status_code,omitempty"`
	ResponseSize int64             `json:"response_size,omitempty"`
}

// UserInfo 用户信息结构
// (UserInfo represents user information structure)
type UserInfo struct {
	UserID       string   `json:"user_id"`
	Username     string   `json:"username"`
	Email        string   `json:"email"`
	Roles        []string `json:"roles"`
	IsActive     bool     `json:"is_active"`
	LastLoginAt  string   `json:"last_login_at,omitempty"`
	CreatedAt    string   `json:"created_at"`
	SessionID    string   `json:"session_id,omitempty"`
	Organization string   `json:"organization,omitempty"`
}

// DatabaseOperation 数据库操作结构
// (DatabaseOperation represents database operation structure)
type DatabaseOperation struct {
	Operation    string        `json:"operation"`
	Table        string        `json:"table"`
	Query        string        `json:"query,omitempty"`
	Parameters   []interface{} `json:"parameters,omitempty"`
	Duration     time.Duration `json:"duration"`
	RowsAffected int64         `json:"rows_affected,omitempty"`
	Error        string        `json:"error,omitempty"`
}

// APIResponse API响应结构
// (APIResponse represents API response structure)
type APIResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data,omitempty"`
	Error      string      `json:"error,omitempty"`
	ErrorCode  string      `json:"error_code,omitempty"`
	Message    string      `json:"message,omitempty"`
	Timestamp  string      `json:"timestamp"`
	RequestID  string      `json:"request_id"`
	Duration   string      `json:"duration"`
}

// StructuredLogService 结构化日志服务
// (StructuredLogService provides structured logging functionality)
type StructuredLogService struct {
	logger log.Logger
}

// NewStructuredLogService 创建结构化日志服务
// (NewStructuredLogService creates a new structured log service)
func NewStructuredLogService(logger log.Logger) *StructuredLogService {
	return &StructuredLogService{
		logger: logger,
	}
}

// LogHTTPRequest 记录HTTP请求（结构化日志）
// (LogHTTPRequest logs HTTP request with structured logging)
func (s *StructuredLogService) LogHTTPRequest(ctx context.Context, req *RequestInfo) {
	s.logger.Infow("HTTP request started",
		"request_id", req.RequestID,
		"method", req.Method,
		"path", req.Path,
		"user_agent", req.UserAgent,
		"remote_ip", req.RemoteIP,
		"headers", req.Headers,
		"query_params", req.QueryParams,
		"start_time", req.StartTime,
	)
}

// LogHTTPResponse 记录HTTP响应（结构化日志）
// (LogHTTPResponse logs HTTP response with structured logging)
func (s *StructuredLogService) LogHTTPResponse(ctx context.Context, req *RequestInfo, resp *APIResponse) {
	// 计算请求处理时间 (Calculate request processing time)
	req.Duration = time.Since(req.StartTime)
	
	logLevel := "info"
	if req.StatusCode >= 400 && req.StatusCode < 500 {
		logLevel = "warn"
	} else if req.StatusCode >= 500 {
		logLevel = "error"
	}
	
	logMethod := s.logger.Infow
	if logLevel == "warn" {
		logMethod = s.logger.Warnw
	} else if logLevel == "error" {
		logMethod = s.logger.Errorw
	}
	
	logMethod("HTTP request completed",
		"request_id", req.RequestID,
		"method", req.Method,
		"path", req.Path,
		"status_code", req.StatusCode,
		"duration", req.Duration,
		"response_size", req.ResponseSize,
		"success", resp.Success,
		"error_code", resp.ErrorCode,
		"remote_ip", req.RemoteIP,
	)
}

// LogUserOperation 记录用户操作（结构化日志）
// (LogUserOperation logs user operation with structured logging)
func (s *StructuredLogService) LogUserOperation(ctx context.Context, user *UserInfo, operation, resource string, success bool, details map[string]interface{}) {
	logData := map[string]interface{}{
		"user_id":      user.UserID,
		"username":     user.Username,
		"email":        user.Email,
		"roles":        user.Roles,
		"is_active":    user.IsActive,
		"session_id":   user.SessionID,
		"organization": user.Organization,
		"operation":    operation,
		"resource":     resource,
		"success":      success,
		"timestamp":    time.Now().Format(time.RFC3339),
	}
	
	// 添加详细信息 (Add details)
	for k, v := range details {
		logData[k] = v
	}
	
	// 转换为键值对 (Convert to key-value pairs)
	var keyValues []interface{}
	for k, v := range logData {
		keyValues = append(keyValues, k, v)
	}
	
	if success {
		s.logger.Infow("User operation completed", keyValues...)
	} else {
		s.logger.Errorw("User operation failed", keyValues...)
	}
}

// LogDatabaseOperation 记录数据库操作（结构化日志）
// (LogDatabaseOperation logs database operation with structured logging)
func (s *StructuredLogService) LogDatabaseOperation(ctx context.Context, dbOp *DatabaseOperation) {
	if dbOp.Error != "" {
		s.logger.Errorw("Database operation failed",
			"operation", dbOp.Operation,
			"table", dbOp.Table,
			"query", dbOp.Query,
			"parameters", dbOp.Parameters,
			"duration", dbOp.Duration,
			"error", dbOp.Error,
		)
	} else {
		s.logger.Infow("Database operation completed",
			"operation", dbOp.Operation,
			"table", dbOp.Table,
			"query", dbOp.Query,
			"parameters", dbOp.Parameters,
			"duration", dbOp.Duration,
			"rows_affected", dbOp.RowsAffected,
		)
	}
}

// LogBusinessEvent 记录业务事件（结构化日志）
// (LogBusinessEvent logs business event with structured logging)
func (s *StructuredLogService) LogBusinessEvent(ctx context.Context, eventType, eventName string, metadata map[string]interface{}) {
	event := map[string]interface{}{
		"event_type": eventType,
		"event_name": eventName,
		"timestamp":  time.Now().Format(time.RFC3339),
		"trace_id":   ctx.Value("trace_id"),
		"user_id":    ctx.Value("user_id"),
	}
	
	// 添加元数据 (Add metadata)
	for k, v := range metadata {
		event[k] = v
	}
	
	// 转换为键值对 (Convert to key-value pairs)
	var keyValues []interface{}
	for k, v := range event {
		keyValues = append(keyValues, k, v)
	}
	
	s.logger.Infow("Business event occurred", keyValues...)
}

// LogPerformanceMetrics 记录性能指标（结构化日志）
// (LogPerformanceMetrics logs performance metrics with structured logging)
func (s *StructuredLogService) LogPerformanceMetrics(ctx context.Context, component string, metrics map[string]interface{}) {
	s.logger.Infow("Performance metrics",
		"component", component,
		"timestamp", time.Now().Format(time.RFC3339),
		"cpu_usage", metrics["cpu_usage"],
		"memory_usage", metrics["memory_usage"],
		"goroutines", metrics["goroutines"],
		"gc_pause", metrics["gc_pause"],
		"requests_per_second", metrics["requests_per_second"],
		"response_time_p95", metrics["response_time_p95"],
		"error_rate", metrics["error_rate"],
	)
}

// demonstrateHTTPRequestLogging 演示HTTP请求日志记录
// (demonstrateHTTPRequestLogging demonstrates HTTP request logging)
func demonstrateHTTPRequestLogging() {
	fmt.Println("=== Demonstrating HTTP Request Logging ===")
	fmt.Println()
	
	logger := log.Std().WithValues("component", "http-server")
	service := NewStructuredLogService(logger)
	ctx := context.Background()
	
	// 模拟HTTP请求 (Simulate HTTP requests)
	requests := []*RequestInfo{
		{
			RequestID: "req-001",
			Method:    "GET",
			Path:      "/api/users",
			UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X)",
			RemoteIP:  "192.168.1.100",
			Headers: map[string]string{
				"Authorization": "Bearer token123",
				"Content-Type":  "application/json",
			},
			QueryParams: map[string]string{
				"page":  "1",
				"limit": "10",
			},
			StartTime: time.Now(),
		},
		{
			RequestID: "req-002",
			Method:    "POST",
			Path:      "/api/users",
			UserAgent: "curl/7.68.0",
			RemoteIP:  "10.0.0.50",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			StartTime: time.Now(),
		},
	}
	
	for i, req := range requests {
		// 记录请求开始 (Log request start)
		service.LogHTTPRequest(ctx, req)
		
		// 模拟请求处理 (Simulate request processing)
		time.Sleep(time.Duration(50+i*30) * time.Millisecond)
		
		// 模拟响应 (Simulate response)
		var resp *APIResponse
		if i == 0 {
			// 成功响应 (Successful response)
			req.StatusCode = 200
			req.ResponseSize = 1024
			resp = &APIResponse{
				Success:   true,
				Data:      []string{"user1", "user2"},
				Timestamp: time.Now().Format(time.RFC3339),
				RequestID: req.RequestID,
				Duration:  time.Since(req.StartTime).String(),
			}
		} else {
			// 错误响应 (Error response)
			req.StatusCode = 400
			req.ResponseSize = 256
			resp = &APIResponse{
				Success:   false,
				Error:     "Invalid request payload",
				ErrorCode: "INVALID_PAYLOAD",
				Message:   "Request validation failed",
				Timestamp: time.Now().Format(time.RFC3339),
				RequestID: req.RequestID,
				Duration:  time.Since(req.StartTime).String(),
			}
		}
		
		// 记录响应 (Log response)
		service.LogHTTPResponse(ctx, req, resp)
		fmt.Println()
	}
}

// demonstrateUserOperationLogging 演示用户操作日志记录
// (demonstrateUserOperationLogging demonstrates user operation logging)
func demonstrateUserOperationLogging() {
	fmt.Println("=== Demonstrating User Operation Logging ===")
	fmt.Println()
	
	logger := log.Std().WithValues("component", "user-service")
	service := NewStructuredLogService(logger)
	ctx := context.WithValue(context.Background(), "trace_id", "trace-12345")
	
	// 模拟用户信息 (Simulate user information)
	user := &UserInfo{
		UserID:       "user-67890",
		Username:     "john_doe",
		Email:        "john@example.com",
		Roles:        []string{"user", "editor"},
		IsActive:     true,
		LastLoginAt:  "2023-12-01T10:00:00Z",
		CreatedAt:    "2023-01-15T08:30:00Z",
		SessionID:    "session-abc123",
		Organization: "example-org",
	}
	
	// 模拟不同的用户操作 (Simulate different user operations)
	operations := []struct {
		operation string
		resource  string
		success   bool
		details   map[string]interface{}
	}{
		{
			operation: "create",
			resource:  "document",
			success:   true,
			details: map[string]interface{}{
				"document_id":   "doc-123",
				"document_type": "report",
				"size_bytes":    2048,
				"permissions":   []string{"read", "write"},
			},
		},
		{
			operation: "update",
			resource:  "profile",
			success:   true,
			details: map[string]interface{}{
				"fields_changed": []string{"email", "phone"},
				"old_email":      "old@example.com",
				"new_email":      "john@example.com",
			},
		},
		{
			operation: "delete",
			resource:  "file",
			success:   false,
			details: map[string]interface{}{
				"file_id":     "file-456",
				"error_code":  "PERMISSION_DENIED",
				"error_msg":   "User does not have delete permissions",
				"file_size":   1024000,
				"file_type":   "image/jpeg",
			},
		},
	}
	
	for _, op := range operations {
		service.LogUserOperation(ctx, user, op.operation, op.resource, op.success, op.details)
	}
	fmt.Println()
}

// demonstrateDatabaseLogging 演示数据库操作日志记录
// (demonstrateDatabaseLogging demonstrates database operation logging)
func demonstrateDatabaseLogging() {
	fmt.Println("=== Demonstrating Database Operation Logging ===")
	fmt.Println()
	
	logger := log.Std().WithValues("component", "database")
	service := NewStructuredLogService(logger)
	ctx := context.Background()
	
	// 模拟数据库操作 (Simulate database operations)
	operations := []*DatabaseOperation{
		{
			Operation:    "SELECT",
			Table:        "users",
			Query:        "SELECT id, name, email FROM users WHERE active = ? LIMIT ?",
			Parameters:   []interface{}{true, 10},
			Duration:     25 * time.Millisecond,
			RowsAffected: 10,
		},
		{
			Operation:    "INSERT",
			Table:        "user_logs",
			Query:        "INSERT INTO user_logs (user_id, action, timestamp) VALUES (?, ?, ?)",
			Parameters:   []interface{}{"user-123", "login", time.Now()},
			Duration:     15 * time.Millisecond,
			RowsAffected: 1,
		},
		{
			Operation:  "UPDATE",
			Table:      "users",
			Query:      "UPDATE users SET last_login = ? WHERE id = ?",
			Parameters: []interface{}{time.Now(), "user-123"},
			Duration:   120 * time.Millisecond,
			Error:      "deadlock detected",
		},
		{
			Operation:    "DELETE",
			Table:        "sessions",
			Query:        "DELETE FROM sessions WHERE expires_at < ?",
			Parameters:   []interface{}{time.Now().Add(-24 * time.Hour)},
			Duration:     45 * time.Millisecond,
			RowsAffected: 156,
		},
	}
	
	for _, op := range operations {
		service.LogDatabaseOperation(ctx, op)
	}
	fmt.Println()
}

// demonstrateBusinessEventLogging 演示业务事件日志记录
// (demonstrateBusinessEventLogging demonstrates business event logging)
func demonstrateBusinessEventLogging() {
	fmt.Println("=== Demonstrating Business Event Logging ===")
	fmt.Println()
	
	logger := log.Std().WithValues("component", "business-events")
	service := NewStructuredLogService(logger)
	
	// 创建带用户信息的上下文 (Create context with user information)
	ctx := context.WithValue(context.Background(), "trace_id", "trace-98765")
	ctx = context.WithValue(ctx, "user_id", "user-54321")
	
	// 模拟业务事件 (Simulate business events)
	events := []struct {
		eventType string
		eventName string
		metadata  map[string]interface{}
	}{
		{
			eventType: "user_activity",
			eventName: "user_registered",
			metadata: map[string]interface{}{
				"registration_method": "email",
				"user_agent":         "Chrome/91.0",
				"referrer":           "google.com",
				"plan_type":          "free",
				"verification_sent":  true,
			},
		},
		{
			eventType: "payment",
			eventName: "payment_processed",
			metadata: map[string]interface{}{
				"payment_id":     "pay-789",
				"amount":         29.99,
				"currency":       "USD",
				"payment_method": "credit_card",
				"gateway":        "stripe",
				"subscription":   "premium",
				"billing_cycle":  "monthly",
			},
		},
		{
			eventType: "security",
			eventName: "suspicious_login_attempt",
			metadata: map[string]interface{}{
				"ip_address":     "203.0.113.42",
				"country":        "Unknown",
				"failed_attempts": 5,
				"account_locked": true,
				"detection_rule": "brute_force",
				"alert_sent":     true,
			},
		},
		{
			eventType: "data_processing",
			eventName: "batch_job_completed",
			metadata: map[string]interface{}{
				"job_id":         "job-456",
				"job_type":       "user_export",
				"records_processed": 10000,
				"duration_seconds":  300,
				"output_file":    "export_20231201.csv",
				"file_size_mb":   15.2,
			},
		},
	}
	
	for _, event := range events {
		service.LogBusinessEvent(ctx, event.eventType, event.eventName, event.metadata)
	}
	fmt.Println()
}

// demonstratePerformanceLogging 演示性能指标日志记录
// (demonstratePerformanceLogging demonstrates performance metrics logging)
func demonstratePerformanceLogging() {
	fmt.Println("=== Demonstrating Performance Metrics Logging ===")
	fmt.Println()
	
	logger := log.Std().WithValues("component", "monitoring")
	service := NewStructuredLogService(logger)
	ctx := context.Background()
	
	// 模拟不同组件的性能指标 (Simulate performance metrics for different components)
	components := []struct {
		name    string
		metrics map[string]interface{}
	}{
		{
			name: "web_server",
			metrics: map[string]interface{}{
				"cpu_usage":           45.2,
				"memory_usage":        68.5,
				"goroutines":          150,
				"gc_pause":            "2.3ms",
				"requests_per_second": 1250,
				"response_time_p95":   "95ms",
				"error_rate":          0.02,
			},
		},
		{
			name: "database",
			metrics: map[string]interface{}{
				"cpu_usage":           72.8,
				"memory_usage":        85.3,
				"goroutines":          45,
				"gc_pause":            "1.8ms",
				"requests_per_second": 800,
				"response_time_p95":   "25ms",
				"error_rate":          0.001,
			},
		},
		{
			name: "cache_server",
			metrics: map[string]interface{}{
				"cpu_usage":           15.6,
				"memory_usage":        42.1,
				"goroutines":          20,
				"gc_pause":            "0.5ms",
				"requests_per_second": 3000,
				"response_time_p95":   "2ms",
				"error_rate":          0.0001,
			},
		},
	}
	
	for _, comp := range components {
		service.LogPerformanceMetrics(ctx, comp.name, comp.metrics)
	}
	fmt.Println()
}

// demonstrateComplexStructuredLogging 演示复杂结构化日志记录
// (demonstrateComplexStructuredLogging demonstrates complex structured logging)
func demonstrateComplexStructuredLogging() {
	fmt.Println("=== Demonstrating Complex Structured Logging ===")
	fmt.Println()
	
	logger := log.Std().WithValues("service", "order-processing")
	
	// 模拟复杂的订单处理流程 (Simulate complex order processing flow)
	orderData := map[string]interface{}{
		"order_id":      "order-789123",
		"customer_id":   "cust-456",
		"customer_tier": "premium",
		"items": []map[string]interface{}{
			{
				"product_id": "prod-001",
				"name":       "Premium Widget",
				"quantity":   2,
				"price":      29.99,
				"category":   "electronics",
			},
			{
				"product_id": "prod-002",
				"name":       "Deluxe Gadget",
				"quantity":   1,
				"price":      79.99,
				"category":   "accessories",
			},
		},
		"shipping": map[string]interface{}{
			"address": map[string]interface{}{
				"street":  "123 Main St",
				"city":    "Anytown",
				"state":   "CA",
				"zip":     "12345",
				"country": "US",
			},
			"method":      "express",
			"cost":        9.99,
			"tracking_id": "TRK123456789",
		},
		"payment": map[string]interface{}{
			"method":        "credit_card",
			"last_four":     "1234",
			"amount":        149.97,
			"currency":      "USD",
			"transaction_id": "txn_abc123",
			"processor":     "stripe",
		},
		"timestamps": map[string]interface{}{
			"created_at":   time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
			"processed_at": time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
			"shipped_at":   time.Now().Format(time.RFC3339),
		},
		"metadata": map[string]interface{}{
			"source":         "web",
			"campaign_id":    "holiday2023",
			"discount_code":  "SAVE20",
			"discount_amount": 30.00,
			"notes":          "Gift wrapping requested",
		},
	}
	
	// 转换为键值对进行日志记录 (Convert to key-value pairs for logging)
	var keyValues []interface{}
	for k, v := range orderData {
		keyValues = append(keyValues, k, v)
	}
	
	logger.Infow("Order processing completed", keyValues...)
	fmt.Println()
}

func main() {
	fmt.Println("=== Structured Logging Example ===")
	fmt.Println("This example demonstrates JSON format and structured field logging.")
	fmt.Println()
	
	// 1. 初始化结构化日志系统 (Initialize structured logging system)
	opts := log.NewOptions()
	opts.Level = "debug"
	opts.Format = "json"         // 使用JSON格式 (Use JSON format)
	opts.EnableColor = false     // JSON格式不需要颜色 (JSON format doesn't need color)
	opts.DisableCaller = false   // 显示调用者信息 (Show caller info)
	opts.DisableStacktrace = true // 禁用堆栈跟踪以保持JSON整洁 (Disable stack trace to keep JSON clean)
	opts.OutputPaths = []string{"stdout"}
	
	// 初始化日志记录器 (Initialize logger)
	log.Init(opts)
	
	fmt.Println("Structured logger initialized with JSON format")
	fmt.Println()
	
	// 2. 演示HTTP请求日志记录 (Demonstrate HTTP request logging)
	demonstrateHTTPRequestLogging()
	
	// 3. 演示用户操作日志记录 (Demonstrate user operation logging)
	demonstrateUserOperationLogging()
	
	// 4. 演示数据库操作日志记录 (Demonstrate database operation logging)
	demonstrateDatabaseLogging()
	
	// 5. 演示业务事件日志记录 (Demonstrate business event logging)
	demonstrateBusinessEventLogging()
	
	// 6. 演示性能指标日志记录 (Demonstrate performance metrics logging)
	demonstratePerformanceLogging()
	
	// 7. 演示复杂结构化日志记录 (Demonstrate complex structured logging)
	demonstrateComplexStructuredLogging()
	
	log.Std().Info("Structured logging example completed successfully")
	fmt.Println("=== Example completed successfully ===")
} 