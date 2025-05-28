/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * Log levels example demonstrating different logging levels and dynamic level adjustment.
 */

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

// LogLevelService 日志级别服务
// (LogLevelService provides log level functionality)
type LogLevelService struct {
	logger log.Logger
}

// NewLogLevelService 创建日志级别服务
// (NewLogLevelService creates a new log level service)
func NewLogLevelService(logger log.Logger) *LogLevelService {
	return &LogLevelService{
		logger: logger,
	}
}

// DemonstrateAllLevels 演示所有日志级别
// (DemonstrateAllLevels demonstrates all logging levels)
func (s *LogLevelService) DemonstrateAllLevels(ctx context.Context) {
	s.logger.Debug("This is a DEBUG message - detailed diagnostic information")
	s.logger.Debugw("Debug with fields", 
		"function", "DemonstrateAllLevels",
		"detail_level", "high",
		"execution_time", time.Now().Format(time.RFC3339))
		
	s.logger.Info("This is an INFO message - general information")
	s.logger.Infow("Info with fields",
		"operation", "user_login",
		"user_id", "user123",
		"status", "success")
		
	s.logger.Warn("This is a WARN message - potentially harmful situation")
	s.logger.Warnw("Warning with fields",
		"issue", "deprecated_api_usage",
		"api_version", "v1.0",
		"recommendation", "upgrade to v2.0")
		
	s.logger.Error("This is an ERROR message - error events but application continues")
	s.logger.Errorw("Error with fields",
		"error_type", "validation_failed",
		"field", "email",
		"value", "invalid-email",
		"recovery_action", "show_error_message")
}

// SimulateApplicationFlow 模拟应用程序流程中的不同日志级别
// (SimulateApplicationFlow simulates different log levels in application flow)
func (s *LogLevelService) SimulateApplicationFlow(ctx context.Context, userID string) error {
	// DEBUG: 详细的执行流程 (Detailed execution flow)
	s.logger.Debugw("Starting user authentication flow",
		"user_id", userID,
		"flow_type", "password_based",
		"timestamp", time.Now().Format(time.RFC3339))
	
	// INFO: 重要的业务事件 (Important business events)
	s.logger.Infow("User authentication initiated",
		"user_id", userID,
		"method", "username_password")
	
	// 模拟验证过程 (Simulate validation process)
	time.Sleep(50 * time.Millisecond)
	
	// DEBUG: 验证步骤详情 (Validation step details)
	s.logger.Debugw("Validating user credentials",
		"user_id", userID,
		"validation_steps", []string{"username_check", "password_hash", "account_status"})
	
	// 模拟可能的警告情况 (Simulate potential warning situations)
	if userID == "inactive_user" {
		s.logger.Warnw("Inactive user attempting login",
			"user_id", userID,
			"last_activity", "2023-06-01",
			"account_status", "inactive",
			"action_required", "account_reactivation")
		return fmt.Errorf("user account is inactive")
	}
	
	// 模拟错误情况 (Simulate error conditions)
	if userID == "blocked_user" {
		s.logger.Errorw("Blocked user authentication attempt",
			"user_id", userID,
			"block_reason", "multiple_failed_attempts",
			"blocked_until", time.Now().Add(30*time.Minute).Format(time.RFC3339),
			"security_alert", true)
		return fmt.Errorf("user account is blocked")
	}
	
	// 成功情况 (Success case)
	s.logger.Infow("User authentication successful",
		"user_id", userID,
		"session_id", fmt.Sprintf("sess_%d", time.Now().Unix()),
		"login_duration", "52ms")
	
	s.logger.Debugw("Authentication flow completed successfully",
		"user_id", userID,
		"total_duration", "102ms",
		"next_step", "redirect_to_dashboard")
	
	return nil
}

// DemonstrateConditionalLogging 演示条件日志记录
// (DemonstrateConditionalLogging demonstrates conditional logging)
func (s *LogLevelService) DemonstrateConditionalLogging(ctx context.Context, operations []string) {
	for i, operation := range operations {
		duration := time.Duration(i*30+20) * time.Millisecond
		
		// DEBUG: 每个操作的详细信息 (Detailed info for each operation)
		s.logger.Debugw("Processing operation",
			"operation", operation,
			"index", i,
			"estimated_duration", duration)
		
		// 模拟处理时间 (Simulate processing time)
		time.Sleep(duration)
		
		// 根据性能条件记录不同级别的日志 (Log different levels based on performance conditions)
		if duration > 80*time.Millisecond {
			s.logger.Warnw("Slow operation detected",
				"operation", operation,
				"duration", duration,
				"threshold", "80ms",
				"suggestion", "consider_optimization")
		} else if duration > 50*time.Millisecond {
			s.logger.Infow("Operation completed",
				"operation", operation,
				"duration", duration,
				"performance", "normal")
		} else {
			s.logger.Debugw("Fast operation completed",
				"operation", operation,
				"duration", duration,
				"performance", "excellent")
		}
	}
}

// DemonstrateLevelFiltering 演示日志级别过滤
// (DemonstrateLevelFiltering demonstrates log level filtering)
func DemonstrateLevelFiltering() {
	fmt.Println("=== Demonstrating Log Level Filtering ===")
	fmt.Println()
	
	levels := []string{"debug", "info", "warn", "error"}
	
	for _, level := range levels {
		fmt.Printf("--- Setting log level to: %s ---\n", level)
		
		// 创建新的日志选项 (Create new log options)
		opts := log.NewOptions()
		opts.Level = level
		opts.Format = "text"
		opts.EnableColor = true
		opts.DisableCaller = false
		opts.DisableStacktrace = true
		opts.OutputPaths = []string{"stdout"}
		
		// 重新初始化日志记录器 (Reinitialize logger)
		log.Init(opts)
		
		logger := log.Std().WithValues("level_demo", level)
		service := NewLogLevelService(logger)
		ctx := context.Background()
		
		// 尝试记录所有级别的日志 (Try to log all levels)
		service.DemonstrateAllLevels(ctx)
		fmt.Println()
	}
}

// DemonstrateComponentLevels 演示组件特定的日志级别
// (DemonstrateComponentLevels demonstrates component-specific log levels)
func DemonstrateComponentLevels() {
	fmt.Println("=== Demonstrating Component-Specific Logging ===")
	fmt.Println()
	
	// 创建基础配置 (Create base configuration)
	opts := log.NewOptions()
	opts.Level = "info"
	opts.Format = "text"
	opts.EnableColor = true
	opts.DisableCaller = false
	opts.DisableStacktrace = true
	opts.OutputPaths = []string{"stdout"}
	
	log.Init(opts)
	
	// 为不同组件创建不同的日志记录器 (Create different loggers for different components)
	components := map[string]log.Logger{
		"database":   log.Std().WithValues("component", "database"),
		"auth":       log.Std().WithValues("component", "auth"),
		"api":        log.Std().WithValues("component", "api"),
		"cache":      log.Std().WithValues("component", "cache"),
		"monitoring": log.Std().WithValues("component", "monitoring"),
	}
	
	ctx := context.Background()
	
	// 模拟不同组件的日志记录活动 (Simulate logging activity from different components)
	fmt.Println("Database component operations:")
	dbService := NewLogLevelService(components["database"])
	dbService.SimulateApplicationFlow(ctx, "db_user_001")
	fmt.Println()
	
	fmt.Println("Authentication component operations:")
	authService := NewLogLevelService(components["auth"])
	authService.SimulateApplicationFlow(ctx, "auth_user_002")
	fmt.Println()
	
	fmt.Println("API component operations:")
	apiService := NewLogLevelService(components["api"])
	apiService.DemonstrateConditionalLogging(ctx, []string{"validate_request", "process_data", "generate_response"})
	fmt.Println()
	
	fmt.Println("Cache component operations:")
	cacheService := NewLogLevelService(components["cache"])
	cacheService.DemonstrateConditionalLogging(ctx, []string{"cache_lookup", "cache_miss", "cache_update"})
	fmt.Println()
}

// DemonstrateProductionScenarios 演示生产环境场景
// (DemonstrateProductionScenarios demonstrates production environment scenarios)
func DemonstrateProductionScenarios() {
	fmt.Println("=== Demonstrating Production Environment Scenarios ===")
	fmt.Println()
	
	// 生产环境通常使用INFO级别 (Production environments typically use INFO level)
	opts := log.NewOptions()
	opts.Level = "info"
	opts.Format = "json"
	opts.EnableColor = false
	opts.DisableCaller = false
	opts.DisableStacktrace = true
	opts.OutputPaths = []string{"stdout"}
	
	log.Init(opts)
	
	logger := log.Std().WithValues("env", "production", "service", "user-api")
	service := NewLogLevelService(logger)
	ctx := context.Background()
	
	// 模拟生产环境中的典型操作 (Simulate typical operations in production)
	scenarios := []struct {
		name   string
		userID string
	}{
		{"Normal user login", "user_12345"},
		{"Inactive user login attempt", "inactive_user"},
		{"Blocked user login attempt", "blocked_user"},
		{"Admin user login", "admin_67890"},
	}
	
	for _, scenario := range scenarios {
		fmt.Printf("Scenario: %s\n", scenario.name)
		err := service.SimulateApplicationFlow(ctx, scenario.userID)
		if err != nil {
			fmt.Printf("Scenario result: %v\n", err)
		} else {
			fmt.Printf("Scenario result: Success\n")
		}
		fmt.Println()
	}
}

// DemonstratePerformanceImpact 演示日志级别对性能的影响
// (DemonstratePerformanceImpact demonstrates performance impact of log levels)
func DemonstratePerformanceImpact() {
	fmt.Println("=== Demonstrating Performance Impact of Log Levels ===")
	fmt.Println()
	
	levels := []string{"debug", "info", "warn", "error"}
	iterations := 1000
	
	for _, level := range levels {
		fmt.Printf("Testing performance with log level: %s\n", level)
		
		// 配置日志级别 (Configure log level)
		opts := log.NewOptions()
		opts.Level = level
		opts.Format = "text"
		opts.EnableColor = false
		opts.DisableCaller = true  // 禁用调用者信息以减少性能影响 (Disable caller info to reduce performance impact)
		opts.DisableStacktrace = true
		opts.OutputPaths = []string{"stdout"}
		
		log.Init(opts)
		
		logger := log.Std().WithValues("perf_test", level)
		
		// 测量日志记录性能 (Measure logging performance)
		start := time.Now()
		
		for i := 0; i < iterations; i++ {
			// 记录所有级别的日志 (Log all levels)
			logger.Debugw("Debug message", "iteration", i, "data", "debug_data")
			logger.Infow("Info message", "iteration", i, "data", "info_data")
			logger.Warnw("Warn message", "iteration", i, "data", "warn_data")
			logger.Errorw("Error message", "iteration", i, "data", "error_data")
		}
		
		duration := time.Since(start)
		
		fmt.Printf("Completed %d iterations in %v (avg: %v per iteration)\n", 
			iterations*4, duration, duration/time.Duration(iterations*4))
		fmt.Println()
	}
}

// DemonstrateAdvancedLevelUsage 演示高级日志级别用法
// (DemonstrateAdvancedLevelUsage demonstrates advanced log level usage)
func DemonstrateAdvancedLevelUsage() {
	fmt.Println("=== Demonstrating Advanced Log Level Usage ===")
	fmt.Println()
	
	// 配置基础日志 (Configure base logging)
	opts := log.NewOptions()
	opts.Level = "info"
	opts.Format = "text"
	opts.EnableColor = true
	opts.DisableCaller = false
	opts.DisableStacktrace = true
	opts.OutputPaths = []string{"stdout"}
	
	log.Init(opts)
	
	logger := log.Std().WithValues("advanced_demo", true)
	
	// 演示基于条件的日志级别选择 (Demonstrate conditional log level selection)
	fmt.Println("1. Conditional log level selection:")
	
	conditions := []struct {
		severity    string
		shouldAlert bool
		errorCount  int
	}{
		{"low", false, 1},
		{"medium", false, 5},
		{"high", true, 15},
		{"critical", true, 50},
	}
	
	for _, condition := range conditions {
		message := fmt.Sprintf("System condition detected: %s", condition.severity)
		
		if condition.severity == "critical" {
			logger.Errorw(message,
				"severity", condition.severity,
				"error_count", condition.errorCount,
				"alert_sent", condition.shouldAlert,
				"requires_immediate_attention", true)
		} else if condition.severity == "high" {
			logger.Warnw(message,
				"severity", condition.severity,
				"error_count", condition.errorCount,
				"alert_sent", condition.shouldAlert,
				"monitoring_required", true)
		} else {
			logger.Infow(message,
				"severity", condition.severity,
				"error_count", condition.errorCount,
				"status", "normal")
		}
	}
	fmt.Println()
	
	// 演示基于上下文的日志级别 (Demonstrate context-based log levels)
	fmt.Println("2. Context-based log levels:")
	
	contexts := []struct {
		userType string
		isDebug  bool
	}{
		{"regular", false},
		{"premium", false},
		{"admin", true},
		{"developer", true},
	}
	
	for _, userCtx := range contexts {
		contextLogger := logger.WithValues("user_type", userCtx.userType, "debug_enabled", userCtx.isDebug)
		
		if userCtx.isDebug {
			contextLogger.Debugw("Debug mode enabled for user",
				"detailed_info", "full_request_response_logging",
				"performance_metrics", true,
				"stack_traces", true)
		}
		
		contextLogger.Infow("User session started",
			"session_type", userCtx.userType,
			"features_enabled", []string{"basic", "advanced"}[:1+len(userCtx.userType)%2])
	}
	fmt.Println()
	
	// 演示日志聚合和批处理 (Demonstrate log aggregation and batching)
	fmt.Println("3. Log aggregation scenario:")
	
	errorCounts := map[string]int{
		"database_timeout":     3,
		"api_rate_limit":      12,
		"authentication_fail": 25,
		"permission_denied":    7,
	}
	
	for errorType, count := range errorCounts {
		if count > 20 {
			logger.Errorw("High error count detected",
				"error_type", errorType,
				"count", count,
				"threshold", 20,
				"action", "investigate_immediately")
		} else if count > 10 {
			logger.Warnw("Elevated error count",
				"error_type", errorType,
				"count", count,
				"threshold", 10,
				"action", "monitor_closely")
		} else {
			logger.Infow("Normal error count",
				"error_type", errorType,
				"count", count,
				"status", "within_normal_range")
		}
	}
	fmt.Println()
}

func main() {
	fmt.Println("=== Log Levels Example ===")
	fmt.Println("This example demonstrates different logging levels and dynamic level adjustment.")
	fmt.Println()
	
	// 1. 演示日志级别过滤 (Demonstrate log level filtering)
	DemonstrateLevelFiltering()
	
	// 2. 演示组件特定的日志级别 (Demonstrate component-specific log levels)
	DemonstrateComponentLevels()
	
	// 3. 演示生产环境场景 (Demonstrate production environment scenarios)
	DemonstrateProductionScenarios()
	
	// 4. 演示性能影响 (Demonstrate performance impact)
	DemonstratePerformanceImpact()
	
	// 5. 演示高级用法 (Demonstrate advanced usage)
	DemonstrateAdvancedLevelUsage()
	
	// 最终日志 (Final log)
	opts := log.NewOptions()
	opts.Level = "info"
	opts.Format = "text"
	opts.EnableColor = true
	opts.OutputPaths = []string{"stdout"}
	log.Init(opts)
	
	log.Std().Info("Log levels example completed successfully")
	fmt.Println("=== Example completed successfully ===")
} 