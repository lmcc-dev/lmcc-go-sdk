/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * Custom formatters example demonstrating different log output formats.
 */

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

// CustomFormatterService 自定义格式化器服务
// (CustomFormatterService provides custom formatter functionality)
type CustomFormatterService struct {
	logger log.Logger
}

// NewCustomFormatterService 创建自定义格式化器服务
// (NewCustomFormatterService creates a new custom formatter service)
func NewCustomFormatterService(logger log.Logger) *CustomFormatterService {
	return &CustomFormatterService{
		logger: logger,
	}
}

// LogSampleMessages 记录示例消息
// (LogSampleMessages logs sample messages)
func (s *CustomFormatterService) LogSampleMessages(ctx context.Context) {
	// 简单消息 (Simple messages)
	s.logger.Debug("Debug message for troubleshooting")
	s.logger.Info("Application started successfully")
	s.logger.Warn("Configuration file not found, using defaults")
	s.logger.Error("Failed to connect to database")
	
	// 带字段的消息 (Messages with fields)
	s.logger.Infow("User login successful",
		"user_id", "user123",
		"email", "user@example.com",
		"ip_address", "192.168.1.100",
		"user_agent", "Mozilla/5.0")
	
	s.logger.Errorw("Database operation failed",
		"operation", "INSERT",
		"table", "users",
		"error", "duplicate key violation",
		"duration", 150*time.Millisecond)
	
	// 复杂结构化数据 (Complex structured data)
	s.logger.Infow("API request processed",
		"request_id", "req-12345",
		"method", "POST",
		"path", "/api/v1/users",
		"status_code", 201,
		"response_time", 45*time.Millisecond,
		"request_size", 1024,
		"response_size", 512,
		"headers", map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer token123",
		},
		"metadata", map[string]interface{}{
			"version": "v1.2.3",
			"build":   "20231201-1234",
			"region":  "us-west-2",
		})
}

// DemonstrateTextFormat 演示文本格式
// (DemonstrateTextFormat demonstrates text format)
func DemonstrateTextFormat() {
	fmt.Println("=== Demonstrating Text Format ===")
	fmt.Println()
	
	opts := log.NewOptions()
	opts.Level = "debug"
	opts.Format = "text"
	opts.EnableColor = true
	opts.DisableCaller = false
	opts.DisableStacktrace = true
	opts.OutputPaths = []string{"stdout"}
	
	log.Init(opts)
	
	logger := log.Std().WithValues("format", "text", "component", "demo")
	service := NewCustomFormatterService(logger)
	ctx := context.Background()
	
	service.LogSampleMessages(ctx)
	fmt.Println()
}

// DemonstrateJSONFormat 演示JSON格式
// (DemonstrateJSONFormat demonstrates JSON format)
func DemonstrateJSONFormat() {
	fmt.Println("=== Demonstrating JSON Format ===")
	fmt.Println()
	
	opts := log.NewOptions()
	opts.Level = "debug"
	opts.Format = "json"
	opts.EnableColor = false // JSON格式不需要颜色 (JSON format doesn't need color)
	opts.DisableCaller = false
	opts.DisableStacktrace = true
	opts.OutputPaths = []string{"stdout"}
	
	log.Init(opts)
	
	logger := log.Std().WithValues("format", "json", "component", "demo")
	service := NewCustomFormatterService(logger)
	ctx := context.Background()
	
	service.LogSampleMessages(ctx)
	fmt.Println()
}

// DemonstrateKeyValueFormat 演示键值对格式
// (DemonstrateKeyValueFormat demonstrates key-value format)
func DemonstrateKeyValueFormat() {
	fmt.Println("=== Demonstrating Key-Value Format ===")
	fmt.Println()
	
	opts := log.NewOptions()
	opts.Level = "debug"
	opts.Format = "keyvalue"
	opts.EnableColor = true
	opts.DisableCaller = false
	opts.DisableStacktrace = true
	opts.OutputPaths = []string{"stdout"}
	
	log.Init(opts)
	
	logger := log.Std().WithValues("format", "keyvalue", "component", "demo")
	service := NewCustomFormatterService(logger)
	ctx := context.Background()
	
	service.LogSampleMessages(ctx)
	fmt.Println()
}

// DemonstrateFormatComparison 演示格式对比
// (DemonstrateFormatComparison demonstrates format comparison)
func DemonstrateFormatComparison() {
	fmt.Println("=== Demonstrating Format Comparison ===")
	fmt.Println()
	
	formats := []struct {
		name        string
		format      string
		enableColor bool
		description string
	}{
		{
			name:        "Text Format",
			format:      "text",
			enableColor: true,
			description: "Human-readable text format with color support",
		},
		{
			name:        "JSON Format",
			format:      "json",
			enableColor: false,
			description: "Structured JSON format for machine processing",
		},
		{
			name:        "Key-Value Format",
			format:      "keyvalue",
			enableColor: true,
			description: "Key-value pairs format for easy parsing",
		},
	}
	
	for _, fmt_config := range formats {
		fmt.Printf("--- %s ---\n", fmt_config.name)
		fmt.Printf("Description: %s\n", fmt_config.description)
		fmt.Println()
		
		opts := log.NewOptions()
		opts.Level = "info"
		opts.Format = fmt_config.format
		opts.EnableColor = fmt_config.enableColor
		opts.DisableCaller = false
		opts.DisableStacktrace = true
		opts.OutputPaths = []string{"stdout"}
		
		log.Init(opts)
		
		logger := log.Std().WithValues("format_demo", fmt_config.name)
		
		// 记录一条示例消息 (Log a sample message)
		logger.Infow("Format demonstration message",
			"timestamp", time.Now().Format(time.RFC3339),
			"level", "INFO",
			"message_type", "demo",
			"format_used", fmt_config.format,
			"sample_data", map[string]interface{}{
				"user_id":    "user123",
				"action":     "login",
				"success":    true,
				"duration":   "150ms",
				"ip_address": "192.168.1.100",
			})
		
		fmt.Println()
	}
}

// DemonstrateProductionFormats 演示生产环境格式
// (DemonstrateProductionFormats demonstrates production environment formats)
func DemonstrateProductionFormats() {
	fmt.Println("=== Demonstrating Production Environment Formats ===")
	fmt.Println()
	
	scenarios := []struct {
		name        string
		environment string
		format      string
		level       string
		enableColor bool
		description string
	}{
		{
			name:        "Development Environment",
			environment: "development",
			format:      "text",
			level:       "debug",
			enableColor: true,
			description: "Development environment with detailed logging",
		},
		{
			name:        "Staging Environment",
			environment: "staging",
			format:      "json",
			level:       "info",
			enableColor: false,
			description: "Staging environment with structured logging",
		},
		{
			name:        "Production Environment",
			environment: "production",
			format:      "json",
			level:       "warn",
			enableColor: false,
			description: "Production environment with minimal logging",
		},
		{
			name:        "Debug/Troubleshooting",
			environment: "debug",
			format:      "keyvalue",
			level:       "debug",
			enableColor: true,
			description: "Debug environment for troubleshooting",
		},
	}
	
	for _, scenario := range scenarios {
		fmt.Printf("--- %s ---\n", scenario.name)
		fmt.Printf("Environment: %s\n", scenario.environment)
		fmt.Printf("Description: %s\n", scenario.description)
		fmt.Println()
		
		opts := log.NewOptions()
		opts.Level = scenario.level
		opts.Format = scenario.format
		opts.EnableColor = scenario.enableColor
		opts.DisableCaller = false
		opts.DisableStacktrace = scenario.environment == "production"
		opts.OutputPaths = []string{"stdout"}
		
		log.Init(opts)
		
		logger := log.Std().WithValues(
			"environment", scenario.environment,
			"service", "user-api",
			"version", "v1.2.3")
		
		// 模拟不同环境的日志记录 (Simulate logging for different environments)
		if scenario.level == "debug" {
			logger.Debugw("Debug information available",
				"debug_mode", true,
				"trace_enabled", true)
		}
		
		logger.Infow("Service operational",
			"status", "healthy",
			"uptime", "2h30m",
			"requests_served", 15420)
		
		if scenario.environment != "production" {
			logger.Warnw("Non-production environment warning",
				"config_override", true,
				"test_data_enabled", true)
		}
		
		if scenario.environment == "production" {
			logger.Errorw("Production error example",
				"error_code", "DB_CONNECTION_TIMEOUT",
				"retry_attempts", 3,
				"fallback_enabled", true)
		}
		
		fmt.Println()
	}
}

// DemonstratePerformanceFormats 演示格式性能对比
// (DemonstratePerformanceFormats demonstrates format performance comparison)
func DemonstratePerformanceFormats() {
	fmt.Println("=== Demonstrating Format Performance Comparison ===")
	fmt.Println()
	
	formats := []string{"text", "json", "keyvalue"}
	iterations := 1000
	
	for _, format := range formats {
		fmt.Printf("Testing performance for format: %s\n", format)
		
		opts := log.NewOptions()
		opts.Level = "info"
		opts.Format = format
		opts.EnableColor = false  // 禁用颜色以获得一致的性能测试 (Disable color for consistent performance testing)
		opts.DisableCaller = true // 禁用调用者信息以减少性能影响 (Disable caller info to reduce performance impact)
		opts.DisableStacktrace = true
		opts.OutputPaths = []string{"stdout"}
		
		log.Init(opts)
		
		logger := log.Std().WithValues("perf_test", format)
		
		// 测量格式化性能 (Measure formatting performance)
		start := time.Now()
		
		for i := 0; i < iterations; i++ {
			logger.Infow("Performance test message",
				"iteration", i,
				"format", format,
				"timestamp", time.Now().Unix(),
				"data", map[string]interface{}{
					"user_id":     fmt.Sprintf("user_%d", i),
					"action":      "test_action",
					"duration":    time.Duration(i) * time.Millisecond,
					"success":     i%2 == 0,
					"error_count": i % 10,
				})
		}
		
		duration := time.Since(start)
		
		fmt.Printf("Completed %d iterations in %v (avg: %v per log)\n",
			iterations, duration, duration/time.Duration(iterations))
		fmt.Println()
	}
}

// DemonstrateAdvancedFormatting 演示高级格式化功能
// (DemonstrateAdvancedFormatting demonstrates advanced formatting features)
func DemonstrateAdvancedFormatting() {
	fmt.Println("=== Demonstrating Advanced Formatting Features ===")
	fmt.Println()
	
	// 1. 演示不同数据类型的格式化 (Demonstrate formatting of different data types)
	fmt.Println("1. Different Data Types Formatting:")
	
	opts := log.NewOptions()
	opts.Level = "info"
	opts.Format = "json"
	opts.EnableColor = false
	opts.DisableCaller = false
	opts.DisableStacktrace = true
	opts.OutputPaths = []string{"stdout"}
	
	log.Init(opts)
	
	logger := log.Std().WithValues("advanced_demo", true)
	
	// 不同数据类型的示例 (Examples of different data types)
	logger.Infow("Data types demonstration",
		"string_value", "hello world",
		"integer_value", 42,
		"float_value", 3.14159,
		"boolean_value", true,
		"nil_value", nil,
		"array_value", []string{"item1", "item2", "item3"},
		"map_value", map[string]interface{}{
			"nested_string": "nested value",
			"nested_number": 123,
			"nested_bool":   false,
		},
		"time_value", time.Now(),
		"duration_value", 5*time.Minute+30*time.Second)
	
	fmt.Println()
	
	// 2. 演示错误对象的格式化 (Demonstrate error object formatting)
	fmt.Println("2. Error Object Formatting:")
	
	sampleError := fmt.Errorf("connection failed: %w", fmt.Errorf("timeout after 30s"))
	
	logger.Errorw("Error demonstration",
		"error", sampleError,
		"error_type", "connection_error",
		"component", "database",
		"retry_possible", true,
		"context", map[string]interface{}{
			"host":         "localhost",
			"port":         5432,
			"database":     "userdb",
			"max_retries":  3,
			"current_retry": 1,
		})
	
	fmt.Println()
	
	// 3. 演示大数据对象的格式化 (Demonstrate large data object formatting)
	fmt.Println("3. Large Data Object Formatting:")
	
	largeObject := map[string]interface{}{
		"user_profile": map[string]interface{}{
			"id":       "user_12345",
			"username": "john_doe",
			"email":    "john@example.com",
			"settings": map[string]interface{}{
				"theme":             "dark",
				"notifications":     true,
				"language":          "en",
				"timezone":          "UTC",
				"privacy_settings": map[string]interface{}{
					"profile_visible":    true,
					"email_visible":      false,
					"activity_tracking":  true,
					"data_sharing":       false,
				},
			},
			"permissions": []string{"read", "write", "admin"},
			"last_login":  time.Now().Add(-2 * time.Hour),
			"login_count": 1567,
		},
		"session_info": map[string]interface{}{
			"session_id":   "sess_abc123",
			"ip_address":   "192.168.1.100",
			"user_agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
			"location":     "San Francisco, CA",
			"device_info": map[string]interface{}{
				"type":     "desktop",
				"os":       "macOS",
				"browser":  "Chrome",
				"version":  "91.0.4472.124",
			},
		},
	}
	
	logger.Infow("Large object demonstration",
		"object_type", "user_session",
		"object_size", "large",
		"data", largeObject)
	
	fmt.Println()
}

// DemonstrateCustomFieldFormatting 演示自定义字段格式化
// (DemonstrateCustomFieldFormatting demonstrates custom field formatting)
func DemonstrateCustomFieldFormatting() {
	fmt.Println("=== Demonstrating Custom Field Formatting ===")
	fmt.Println()
	
	// 使用不同的格式来演示字段格式化的差异 (Use different formats to demonstrate field formatting differences)
	formats := []string{"text", "json", "keyvalue"}
	
	for _, format := range formats {
		fmt.Printf("--- Custom Field Formatting with %s format ---\n", format)
		
		opts := log.NewOptions()
		opts.Level = "info"
		opts.Format = format
		opts.EnableColor = format != "json"
		opts.DisableCaller = false
		opts.DisableStacktrace = true
		opts.OutputPaths = []string{"stdout"}
		
		log.Init(opts)
		
		logger := log.Std().WithValues("formatter", format)
		
		// 演示各种字段组合 (Demonstrate various field combinations)
		logger.Infow("Custom field formatting example",
			// 基础字段 (Basic fields)
			"service_name", "user-api",
			"version", "v1.2.3",
			"environment", "production",
			
			// 性能相关字段 (Performance related fields)
			"request_id", "req-12345",
			"response_time_ms", 145,
			"memory_usage_mb", 256.5,
			"cpu_usage_percent", 23.7,
			
			// 业务相关字段 (Business related fields)
			"user_id", "user_67890",
			"operation", "update_profile",
			"success", true,
			"affected_records", 1,
			
			// 复杂结构字段 (Complex structure fields)
			"request_headers", map[string]string{
				"User-Agent":    "API-Client/1.0",
				"Content-Type":  "application/json",
				"Authorization": "Bearer ***",
			},
			"metrics", map[string]interface{}{
				"requests_per_minute": 1250,
				"error_rate":         0.02,
				"avg_response_time":  85.3,
			})
		
		fmt.Println()
	}
}

// demonstrateFileRotation 演示日志文件轮转
// (demonstrateFileRotation demonstrates log file rotation)
func demonstrateFileRotation() {
	fmt.Println("=== File Rotation Demonstration ===")
	fmt.Println()
	
	// 创建临时目录用于演示 (Create temporary directory for demonstration)
	tempDir := "/tmp/lmcc-log-rotation-demo"
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir) // 清理临时文件 (Clean up temporary files)
	
	// 1. 大小基于的轮转 (Size-based rotation)
	fmt.Println("1. Size-based rotation (5MB max, 3 backups):")
	opts1 := log.NewOptions()
	opts1.Level = "info"
	opts1.Format = "json"
	opts1.EnableColor = false
	opts1.DisableCaller = false
	opts1.DisableStacktrace = true
	opts1.OutputPaths = []string{fmt.Sprintf("%s/app-size.log", tempDir)}
	opts1.LogRotateMaxSize = 5    // 5 MB 最大文件大小 (5 MB max file size)
	opts1.LogRotateMaxBackups = 3 // 保留3个备份文件 (Keep 3 backup files)
	opts1.LogRotateCompress = true // 压缩旧文件 (Compress old files)
	
	log.Init(opts1)
	logger1 := log.Std().WithValues("component", "rotation-demo", "type", "size-based")
	
	logger1.Infow("Size-based rotation demo started",
		"max_size_mb", 5,
		"max_backups", 3,
		"compress", true)
	
	// 生成一些日志以触发轮转 (Generate some logs to trigger rotation)
	for i := 0; i < 10; i++ {
		logger1.Infow("Large log entry for size rotation test",
			"iteration", i,
			"timestamp", time.Now(),
			"large_data", generateLargeData(1000), // 生成1KB数据 (Generate 1KB data)
			"description", "This is a test entry to demonstrate size-based log rotation functionality")
	}
	
	// 检查生成的文件 (Check generated files)
	files1, _ := filepath.Glob(fmt.Sprintf("%s/app-size.log*", tempDir))
	fmt.Printf("Generated files: %v\n", files1)
	fmt.Println()
	
	// 2. 时间基于的轮转 (Time-based rotation)
	fmt.Println("2. Time-based rotation (7 days max age):")
	opts2 := log.NewOptions()
	opts2.Level = "info"
	opts2.Format = "text"
	opts2.EnableColor = false
	opts2.DisableCaller = false
	opts2.DisableStacktrace = true
	opts2.OutputPaths = []string{fmt.Sprintf("%s/app-time.log", tempDir)}
	opts2.LogRotateMaxSize = 100   // 100 MB (较大，主要依赖时间轮转)
	opts2.LogRotateMaxAge = 7      // 保留7天 (Keep for 7 days)
	opts2.LogRotateMaxBackups = 10 // 最多10个备份 (Max 10 backups)
	opts2.LogRotateCompress = false // 不压缩 (No compression)
	
	log.Init(opts2)
	logger2 := log.Std().WithValues("component", "rotation-demo", "type", "time-based")
	
	logger2.Infow("Time-based rotation demo started",
		"max_age_days", 7,
		"max_backups", 10,
		"compress", false)
	
	// 生成一些日志 (Generate some logs)
	for i := 0; i < 5; i++ {
		logger2.Infow("Time-based rotation test entry",
			"iteration", i,
			"current_time", time.Now().Format(time.RFC3339),
			"info", "This demonstrates time-based log rotation")
		time.Sleep(100 * time.Millisecond) // 短暂延迟 (Brief delay)
	}
	
	// 检查生成的文件 (Check generated files)
	files2, _ := filepath.Glob(fmt.Sprintf("%s/app-time.log*", tempDir))
	fmt.Printf("Generated files: %v\n", files2)
	fmt.Println()
	
	// 3. 混合轮转策略 (Combined rotation strategy)
	fmt.Println("3. Combined rotation strategy (size + time + count):")
	opts3 := log.NewOptions()
	opts3.Level = "debug"
	opts3.Format = "json"
	opts3.EnableColor = false
	opts3.DisableCaller = false
	opts3.DisableStacktrace = false
	opts3.OutputPaths = []string{fmt.Sprintf("%s/app-combined.log", tempDir)}
	opts3.LogRotateMaxSize = 1     // 1 MB (小文件用于演示)
	opts3.LogRotateMaxAge = 30     // 30天 (30 days)
	opts3.LogRotateMaxBackups = 5  // 5个备份 (5 backups)
	opts3.LogRotateCompress = true // 压缩旧文件 (Compress old files)
	
	log.Init(opts3)
	logger3 := log.Std().WithValues("component", "rotation-demo", "type", "combined")
	
	logger3.Infow("Combined rotation strategy demo started",
		"max_size_mb", 1,
		"max_age_days", 30,
		"max_backups", 5,
		"compress", true)
	
	// 生成更多日志以演示混合策略 (Generate more logs to demonstrate combined strategy)
	rotationService := NewRotationService(logger3)
	rotationService.ProcessLargeDataset()
	
	// 检查最终生成的文件 (Check final generated files)
	files3, _ := filepath.Glob(fmt.Sprintf("%s/app-combined.log*", tempDir))
	fmt.Printf("Generated files: %v\n", files3)
	
	// 显示文件详细信息 (Show file details)
	fmt.Println("\nFile details:")
	for _, file := range files3 {
		if info, err := os.Stat(file); err == nil {
			fmt.Printf("  %s: %d bytes, modified: %s\n", 
				filepath.Base(file), 
				info.Size(), 
				info.ModTime().Format("2006-01-02 15:04:05"))
		}
	}
	
	fmt.Println()
}

// RotationService 轮转演示服务
// (RotationService for demonstrating rotation)
type RotationService struct {
	logger log.Logger
}

// NewRotationService 创建轮转演示服务
// (NewRotationService creates rotation demonstration service)
func NewRotationService(logger log.Logger) *RotationService {
	return &RotationService{
		logger: logger.WithValues("service", "rotation"),
	}
}

// ProcessLargeDataset 处理大数据集以触发轮转
// (ProcessLargeDataset processes large dataset to trigger rotation)
func (rs *RotationService) ProcessLargeDataset() {
	rs.logger.Infow("Starting large dataset processing")
	
	datasets := []string{"users", "orders", "products", "analytics", "logs"}
	
	for _, dataset := range datasets {
		rs.processDataset(dataset)
	}
	
	rs.logger.Infow("Large dataset processing completed")
}

// processDataset 处理单个数据集
// (processDataset processes a single dataset)
func (rs *RotationService) processDataset(name string) {
	logger := rs.logger.WithValues("dataset", name)
	
	logger.Debugw("Starting dataset processing",
		"start_time", time.Now())
	
	// 模拟处理大量记录 (Simulate processing many records)
	recordCount := 50 + len(name)*10 // 根据数据集名称变化记录数
	
	for i := 0; i < recordCount; i++ {
		logger.Debugw("Processing record",
			"record_id", fmt.Sprintf("%s_%04d", name, i),
			"progress", fmt.Sprintf("%.1f%%", float64(i+1)/float64(recordCount)*100),
			"data", generateLargeData(200), // 200字节数据
			"metadata", map[string]interface{}{
				"created": time.Now().Add(-time.Duration(i) * time.Minute),
				"status":  []string{"pending", "processing", "completed"}[i%3],
				"priority": i%5 + 1,
			})
		
		// 每10条记录记录一次进度 (Log progress every 10 records)
		if (i+1)%10 == 0 {
			logger.Infow("Dataset processing progress",
				"completed", i+1,
				"total", recordCount,
				"percentage", fmt.Sprintf("%.1f%%", float64(i+1)/float64(recordCount)*100))
		}
	}
	
	logger.Infow("Dataset processing completed",
		"dataset", name,
		"total_records", recordCount,
		"duration", time.Since(time.Now().Add(-time.Duration(recordCount)*time.Millisecond)))
}

// generateLargeData 生成指定大小的测试数据
// (generateLargeData generates test data of specified size)
func generateLargeData(sizeBytes int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	data := make([]byte, sizeBytes)
	for i := range data {
		data[i] = charset[i%len(charset)]
	}
	return string(data)
}

// demonstrateRotationBestPractices 演示轮转最佳实践
// (demonstrateRotationBestPractices demonstrates rotation best practices)
func demonstrateRotationBestPractices() {
	fmt.Println("=== Log Rotation Best Practices ===")
	fmt.Println()
	
	fmt.Println("1. Production Environment Configuration:")
	fmt.Println("   - Use reasonable file sizes (50-100MB)")
	fmt.Println("   - Keep limited backups (5-10 files)")
	fmt.Println("   - Enable compression to save disk space")
	fmt.Println("   - Set appropriate retention period (7-30 days)")
	fmt.Println()
	
	fmt.Println("2. Development Environment Configuration:")
	fmt.Println("   - Use smaller file sizes (10-20MB)")
	fmt.Println("   - Keep fewer backups (2-3 files)")
	fmt.Println("   - Consider disabling compression for faster access")
	fmt.Println("   - Shorter retention period (1-3 days)")
	fmt.Println()
	
	fmt.Println("3. High-Volume Applications:")
	fmt.Println("   - Monitor disk usage regularly")
	fmt.Println("   - Consider external log aggregation")
	fmt.Println("   - Use structured logging for better parsing")
	fmt.Println("   - Implement log level filtering")
	fmt.Println()
	
	// 演示配置示例 (Demonstrate configuration examples)
	fmt.Println("Configuration Examples:")
	fmt.Println()
	
	fmt.Println("Production (High Volume):")
	prodOpts := log.NewOptions()
	prodOpts.Level = "info"
	prodOpts.Format = "json"
	prodOpts.LogRotateMaxSize = 100
	prodOpts.LogRotateMaxAge = 30
	prodOpts.LogRotateMaxBackups = 10
	prodOpts.LogRotateCompress = true
	
	fmt.Printf("  Level: %s\n", prodOpts.Level)
	fmt.Printf("  Format: %s\n", prodOpts.Format)
	fmt.Printf("  Max Size: %d MB\n", prodOpts.LogRotateMaxSize)
	fmt.Printf("  Max Age: %d days\n", prodOpts.LogRotateMaxAge)
	fmt.Printf("  Max Backups: %d\n", prodOpts.LogRotateMaxBackups)
	fmt.Printf("  Compress: %t\n", prodOpts.LogRotateCompress)
	fmt.Println()
	
	fmt.Println("Development (Low Volume):")
	devOpts := log.NewOptions()
	devOpts.Level = "debug"
	devOpts.Format = "text"
	devOpts.EnableColor = true
	devOpts.LogRotateMaxSize = 10
	devOpts.LogRotateMaxAge = 3
	devOpts.LogRotateMaxBackups = 3
	devOpts.LogRotateCompress = false
	
	fmt.Printf("  Level: %s\n", devOpts.Level)
	fmt.Printf("  Format: %s\n", devOpts.Format)
	fmt.Printf("  Color: %t\n", devOpts.EnableColor)
	fmt.Printf("  Max Size: %d MB\n", devOpts.LogRotateMaxSize)
	fmt.Printf("  Max Age: %d days\n", devOpts.LogRotateMaxAge)
	fmt.Printf("  Max Backups: %d\n", devOpts.LogRotateMaxBackups)
	fmt.Printf("  Compress: %t\n", devOpts.LogRotateCompress)
	fmt.Println()
}

func main() {
	fmt.Println("=== Custom Formatters Example ===")
	fmt.Println("This example demonstrates different log output formats and formatting options.")
	fmt.Println()
	
	// 1. 演示文本格式 (Demonstrate text format)
	DemonstrateTextFormat()
	
	// 2. 演示JSON格式 (Demonstrate JSON format)
	DemonstrateJSONFormat()
	
	// 3. 演示键值对格式 (Demonstrate key-value format)
	DemonstrateKeyValueFormat()
	
	// 4. 演示格式对比 (Demonstrate format comparison)
	DemonstrateFormatComparison()
	
	// 5. 演示生产环境格式 (Demonstrate production environment formats)
	DemonstrateProductionFormats()
	
	// 6. 演示格式性能对比 (Demonstrate format performance comparison)
	DemonstratePerformanceFormats()
	
	// 7. 演示高级格式化功能 (Demonstrate advanced formatting features)
	DemonstrateAdvancedFormatting()
	
	// 8. 演示自定义字段格式化 (Demonstrate custom field formatting)
	DemonstrateCustomFieldFormatting()
	
	// 9. 演示日志文件轮转 (Demonstrate log file rotation)
	demonstrateFileRotation()
	
	// 10. 演示轮转最佳实践 (Demonstrate rotation best practices)
	demonstrateRotationBestPractices()
	
	// 最终日志 (Final log)
	opts := log.NewOptions()
	opts.Level = "info"
	opts.Format = "text"
	opts.EnableColor = true
	opts.OutputPaths = []string{"stdout"}
	log.Init(opts)
	
	log.Std().Info("Custom formatters example completed successfully")
	fmt.Println("=== Example completed successfully ===")
} 