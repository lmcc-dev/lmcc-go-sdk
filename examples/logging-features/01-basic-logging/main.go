/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * Basic logging example demonstrating fundamental logging operations.
 */

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

// UserService 用户服务示例
// (UserService demonstrates user service example)
type UserService struct {
	logger log.Logger
}

// NewUserService 创建用户服务
// (NewUserService creates a user service)
func NewUserService(logger log.Logger) *UserService {
	return &UserService{
		logger: logger,
	}
}

// CreateUser 创建用户（演示基础日志记录）
// (CreateUser creates a user - demonstrates basic logging)
func (us *UserService) CreateUser(ctx context.Context, userID, name, email string) error {
	// 记录操作开始 (Log operation start)
	us.logger.Infow("Starting user creation process", "user_id", userID)
	
	// 验证输入参数 (Validate input parameters)
	if userID == "" {
		us.logger.Error("User creation failed: user ID is required")
		return fmt.Errorf("user ID cannot be empty")
	}
	
	if name == "" {
		us.logger.Warnw("User creation proceeding with empty name", "user_id", userID)
	}
	
	if email == "" {
		us.logger.Errorw("User creation failed: email is required", "user_id", userID)
		return fmt.Errorf("email cannot be empty")
	}
	
	// 记录详细信息 (Log detailed information)
	us.logger.Debugw("User data validation completed", 
		"user_id", userID,
		"name", name,
		"email", email)
	
	// 模拟数据库操作 (Simulate database operation)
	if err := us.simulateCreateUserInDB(ctx, userID, name, email); err != nil {
		us.logger.Errorw("Database operation failed", 
			"user_id", userID,
			"error", err)
		return err
	}
	
	// 记录成功完成 (Log successful completion)
	us.logger.Infow("User created successfully", 
		"user_id", userID,
		"name", name,
		"email", email)
	
	return nil
}

// simulateCreateUserInDB 模拟数据库创建用户
// (simulateCreateUserInDB simulates creating user in database)
func (us *UserService) simulateCreateUserInDB(ctx context.Context, userID, name, email string) error {
	us.logger.Debugw("Executing database query", 
		"operation", "INSERT",
		"table", "users")
	
	// 模拟数据库延迟 (Simulate database latency)
	time.Sleep(100 * time.Millisecond)
	
	// 模拟一些错误情况 (Simulate some error conditions)
	if userID == "error_user" {
		return fmt.Errorf("database constraint violation")
	}
	
	if email == "duplicate@example.com" {
		return fmt.Errorf("email already exists")
	}
	
	us.logger.Debug("Database query executed successfully")
	return nil
}

// GetUser 获取用户（演示带上下文的日志记录）
// (GetUser retrieves a user - demonstrates logging with context)
func (us *UserService) GetUser(ctx context.Context, userID string) error {
	// 使用带上下文的日志记录 (Use logging with context)
	logger := us.logger.WithValues("trace_id", ctx.Value("trace_id"))
	
	logger.Infow("Starting user retrieval", "user_id", userID)
	
	// 模拟查询操作 (Simulate query operation)
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		logger.Infow("User retrieval completed", 
			"user_id", userID,
			"duration", duration)
	}()
	
	if userID == "" {
		logger.Error("User retrieval failed: empty user ID")
		return fmt.Errorf("user ID cannot be empty")
	}
	
	// 模拟数据库查询 (Simulate database query)
	logger.Debugw("Executing user query", "query_type", "SELECT")
	time.Sleep(50 * time.Millisecond)
	
	if userID == "not_found" {
		logger.Warnw("User not found", "user_id", userID)
		return fmt.Errorf("user not found")
	}
	
	logger.Infow("User retrieved successfully", "user_id", userID)
	return nil
}

// BatchOperation 批量操作（演示批量日志记录）
// (BatchOperation batch operations - demonstrates bulk logging)
func (us *UserService) BatchOperation(ctx context.Context, userIDs []string) map[string]error {
	logger := us.logger.WithValues("trace_id", ctx.Value("trace_id"))
	
	logger.Infow("Starting batch operation", 
		"batch_size", len(userIDs),
		"user_ids", fmt.Sprintf("%v", userIDs))
	
	results := make(map[string]error)
	
	for i, userID := range userIDs {
		// 为每个操作创建带索引的日志记录器 (Create indexed logger for each operation)
		itemLogger := logger.WithValues(
			"user_id", userID,
			"batch_index", i)
		
		itemLogger.Debug("Processing batch item")
		
		// 模拟处理 (Simulate processing)
		if err := us.processBatchItem(ctx, userID); err != nil {
			itemLogger.Errorw("Batch item processing failed", "error", err)
			results[userID] = err
		} else {
			itemLogger.Debug("Batch item processed successfully")
		}
	}
	
	// 记录批量操作摘要 (Log batch operation summary)
	successCount := len(userIDs) - len(results)
	logger.Infow("Batch operation completed", 
		"total", len(userIDs),
		"success", successCount,
		"failed", len(results))
	
	return results
}

// processBatchItem 处理批量项目
// (processBatchItem processes batch item)
func (us *UserService) processBatchItem(ctx context.Context, userID string) error {
	// 模拟一些失败情况 (Simulate some failure cases)
	switch userID {
	case "fail_1":
		return fmt.Errorf("processing error for user %s", userID)
	case "fail_2":
		return fmt.Errorf("validation error for user %s", userID)
	default:
		time.Sleep(10 * time.Millisecond) // 模拟处理时间 (Simulate processing time)
		return nil
	}
}

// demonstrateBasicLogging 演示基础日志记录
// (demonstrateBasicLogging demonstrates basic logging)
func demonstrateBasicLogging() {
	fmt.Println("=== Demonstrating Basic Logging ===")
	fmt.Println()
	
	// 1. 创建基础日志记录器 (Create basic logger)
	logger := log.Std()
	
	// 2. 基本日志级别演示 (Basic log level demonstration)
	fmt.Println("1. Basic Log Levels:")
	logger.Debug("This is a debug message")
	logger.Info("Application started successfully")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")
	fmt.Println()
	
	// 3. 带字段的日志记录 (Logging with fields)
	fmt.Println("2. Logging with Fields:")
	logger.Infow("User operation", 
		"operation", "create",
		"user_id", "12345",
		"attempt", 1)
	
	logger.Errorw("Operation failed", 
		"operation", "delete",
		"reason", "permission denied",
		"elapsed", 150*time.Millisecond)
	fmt.Println()
	
	// 4. 创建带前缀的日志记录器 (Create logger with prefix)
	fmt.Println("3. Logger with Prefix:")
	componentLogger := logger.WithValues("component", "user-service")
	componentLogger.Info("Component initialized")
	componentLogger.Debugw("Processing request", "request_id", "req-001")
	fmt.Println()
}

// demonstrateUserServiceLogging 演示用户服务日志记录
// (demonstrateUserServiceLogging demonstrates user service logging)
func demonstrateUserServiceLogging() {
	fmt.Println("=== Demonstrating User Service Logging ===")
	fmt.Println()
	
	// 创建用户服务 (Create user service)
	logger := log.Std().WithValues("service", "user-service")
	userService := NewUserService(logger)
	ctx := context.Background()
	
	// 1. 测试成功的用户创建 (Test successful user creation)
	fmt.Println("1. Successful User Creation:")
	err := userService.CreateUser(ctx, "user_001", "John Doe", "john@example.com")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Println()
	
	// 2. 测试失败的用户创建 (Test failed user creation)
	fmt.Println("2. Failed User Creation (missing email):")
	err = userService.CreateUser(ctx, "user_002", "Jane Doe", "")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Println()
	
	// 3. 测试数据库错误 (Test database error)
	fmt.Println("3. Database Error:")
	err = userService.CreateUser(ctx, "error_user", "Error User", "error@example.com")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Println()
	
	// 4. 测试用户检索 (Test user retrieval)
	fmt.Println("4. User Retrieval:")
	err = userService.GetUser(ctx, "user_001")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	
	err = userService.GetUser(ctx, "not_found")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Println()
}

// demonstrateBatchLogging 演示批量操作日志记录
// (demonstrateBatchLogging demonstrates batch operation logging)
func demonstrateBatchLogging() {
	fmt.Println("=== Demonstrating Batch Operation Logging ===")
	fmt.Println()
	
	// 创建用户服务 (Create user service)
	logger := log.Std().WithValues("service", "user-service")
	userService := NewUserService(logger)
	ctx := context.Background()
	
	// 测试批量操作 (Test batch operation)
	userIDs := []string{"user_1", "user_2", "fail_1", "user_3", "fail_2", "user_4"}
	results := userService.BatchOperation(ctx, userIDs)
	
	fmt.Printf("Batch operation results:\n")
	for userID, err := range results {
		if err != nil {
			fmt.Printf("  %s: FAILED - %v\n", userID, err)
		}
	}
	
	successCount := len(userIDs) - len(results)
	fmt.Printf("Summary: %d/%d succeeded\n", successCount, len(userIDs))
	fmt.Println()
}

// demonstrateLoggingWithContext 演示带上下文的日志记录
// (demonstrateLoggingWithContext demonstrates logging with context)
func demonstrateLoggingWithContext() {
	fmt.Println("=== Demonstrating Logging with Context ===")
	fmt.Println()
	
	logger := log.Std()
	
	// 创建带追踪ID的上下文 (Create context with trace ID)
	ctx := context.WithValue(context.Background(), "trace_id", "trace-12345")
	
	// 使用带上下文的日志记录器 (Use context-aware logger)
	contextLogger := logger.WithValues("trace_id", ctx.Value("trace_id"))
	
	fmt.Println("1. Logging with Context:")
	contextLogger.Info("Processing request")
	contextLogger.Debug("Validating input parameters")
	contextLogger.Info("Request processed successfully")
	fmt.Println()
	
	// 演示嵌套上下文 (Demonstrate nested context)
	fmt.Println("2. Nested Context:")
	nestedCtx := context.WithValue(ctx, "user_id", "user-67890")
	nestedLogger := logger.WithValues("trace_id", nestedCtx.Value("trace_id"), "user_id", nestedCtx.Value("user_id"))
	
	nestedLogger.Info("User operation started")
	nestedLogger.Debug("Executing user validation")
	nestedLogger.Info("User operation completed")
	fmt.Println()
}

// demonstrateErrorLogging 演示错误日志记录
// (demonstrateErrorLogging demonstrates error logging)
func demonstrateErrorLogging() {
	fmt.Println("=== Demonstrating Error Logging ===")
	fmt.Println()
	
	logger := log.Std()
	
	// 1. 简单错误记录 (Simple error logging)
	fmt.Println("1. Simple Error Logging:")
	err := fmt.Errorf("connection timeout")
	logger.Errorw("Database connection failed", "error", err)
	fmt.Println()
	
	// 2. 带上下文的错误记录 (Error logging with context)
	fmt.Println("2. Error Logging with Context:")
	logger.Errorw("User operation failed", 
		"user_id", "user_123",
		"operation", "update",
		"error", err,
		"elapsed", 200*time.Millisecond)
	fmt.Println()
	
	// 3. 嵌套错误记录 (Nested error logging)
	fmt.Println("3. Nested Error Logging:")
	originalErr := fmt.Errorf("permission denied")
	wrappedErr := fmt.Errorf("failed to update user: %w", originalErr)
	
	logger.Errorw("Service operation failed", 
		"error", wrappedErr,
		"service", "user-management")
	fmt.Println()
}

// demonstratePerformanceLogging 演示性能日志记录
// (demonstratePerformanceLogging demonstrates performance logging)
func demonstratePerformanceLogging() {
	fmt.Println("=== Demonstrating Performance Logging ===")
	fmt.Println()
	
	logger := log.Std()
	
	// 1. 操作计时 (Operation timing)
	fmt.Println("1. Operation Timing:")
	
	start := time.Now()
	time.Sleep(100 * time.Millisecond) // 模拟操作 (Simulate operation)
	duration := time.Since(start)
	
	logger.Infow("Database query completed", 
		"query", "SELECT * FROM users",
		"duration", duration,
		"slow_query", duration > 50*time.Millisecond)
	fmt.Println()
	
	// 2. 带阈值的性能记录 (Performance logging with thresholds)
	fmt.Println("2. Performance Logging with Thresholds:")
	
	operations := []struct {
		name     string
		duration time.Duration
	}{
		{"fast_operation", 10 * time.Millisecond},
		{"normal_operation", 50 * time.Millisecond},
		{"slow_operation", 200 * time.Millisecond},
		{"very_slow_operation", 500 * time.Millisecond},
	}
	
	for _, op := range operations {
		var logMethod func(string, ...any)
		if op.duration > 300*time.Millisecond {
			logMethod = logger.Errorw
		} else if op.duration > 100*time.Millisecond {
			logMethod = logger.Warnw
		} else {
			logMethod = logger.Infow
		}
		
		logMethod("Operation completed", 
			"operation", op.name,
			"duration", op.duration)
	}
	fmt.Println()
}

func main() {
	fmt.Println("=== Basic Logging Example ===")
	fmt.Println("This example demonstrates fundamental logging operations using lmcc-go-sdk.")
	fmt.Println()
	
	// 1. 初始化日志系统 (Initialize logging system)
	opts := log.NewOptions()
	opts.Level = "debug"      // 设置日志级别 (Set log level)
	opts.Format = "text"      // 设置输出格式 (Set output format)
	opts.EnableColor = true   // 启用颜色 (Enable color)
	opts.DisableCaller = false // 显示调用者信息 (Show caller info)
	opts.DisableStacktrace = false // 显示堆栈跟踪 (Show stack trace)
	opts.OutputPaths = []string{"stdout"} // 输出到标准输出 (Output to stdout)
	
	// 初始化日志记录器 (Initialize logger)
	log.Init(opts)
	
	fmt.Println("Logger initialized successfully")
	fmt.Println()
	
	// 2. 演示基础日志记录 (Demonstrate basic logging)
	demonstrateBasicLogging()
	
	// 3. 演示用户服务日志记录 (Demonstrate user service logging)
	demonstrateUserServiceLogging()
	
	// 4. 演示批量操作日志记录 (Demonstrate batch operation logging)
	demonstrateBatchLogging()
	
	// 5. 演示带上下文的日志记录 (Demonstrate logging with context)
	demonstrateLoggingWithContext()
	
	// 6. 演示错误日志记录 (Demonstrate error logging)
	demonstrateErrorLogging()
	
	// 7. 演示性能日志记录 (Demonstrate performance logging)
	demonstratePerformanceLogging()
	
	log.Std().Info("Basic logging example completed successfully")
	fmt.Println("=== Example completed successfully ===")
} 