/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * Basic usage example demonstrating integration of config, errors, and log modules.
 */

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

// AppConfig 应用程序配置结构体
// (AppConfig represents the application configuration structure)
type AppConfig struct {
	config.Config                        // 嵌入SDK基础配置 (Embed SDK base configuration)
	App           *AppSpecificConfig     `mapstructure:"app"`
}

// AppSpecificConfig 应用特定的配置
// (AppSpecificConfig represents application-specific configuration)
type AppSpecificConfig struct {
	Name        string        `mapstructure:"name" default:"BasicUsageExample"`
	Version     string        `mapstructure:"version" default:"1.0.0"`
	Environment string        `mapstructure:"environment" default:"development"`
	Timeout     time.Duration `mapstructure:"timeout" default:"30s"`
	MaxRetries  int          `mapstructure:"max_retries" default:"3"`
}

// BusinessService 模拟业务服务
// (BusinessService simulates a business service)
type BusinessService struct {
	config *AppConfig
	logger log.Logger
}

// NewBusinessService 创建新的业务服务实例
// (NewBusinessService creates a new business service instance)
func NewBusinessService(cfg *AppConfig, logger log.Logger) *BusinessService {
	return &BusinessService{
		config: cfg,
		logger: logger,
	}
}

// ProcessData 模拟数据处理操作，演示错误处理
// (ProcessData simulates data processing operation, demonstrating error handling)
func (bs *BusinessService) ProcessData(ctx context.Context, data string) error {
	bs.logger.CtxInfof(ctx, "Starting data processing for: %s", data)

	// 模拟一些业务逻辑和可能的错误
	// (Simulate some business logic and potential errors)
	if data == "invalid" {
		err := errors.New("invalid data format detected")
		return errors.WithCode(
			errors.Wrap(err, "failed to process data"),
			errors.ErrConfigFileRead, // 使用示例错误码 (Using example error code)
		)
	}

	if data == "timeout" {
		// 模拟超时错误 (Simulate timeout error)
		timeoutErr := errors.Errorf("operation timeout after %v", bs.config.App.Timeout)
		return errors.WithMessage(timeoutErr, "data processing failed due to timeout")
	}

	// 模拟处理时间 (Simulate processing time)
	time.Sleep(100 * time.Millisecond)

	bs.logger.CtxInfof(ctx, "Successfully processed data: %s", data)
	return nil
}

// RetryOperation 演示带重试的操作
// (RetryOperation demonstrates operation with retry logic)
func (bs *BusinessService) RetryOperation(ctx context.Context, operation string) error {
	bs.logger.CtxInfof(ctx, "Attempting operation: %s (max retries: %d)", 
		operation, bs.config.App.MaxRetries)

	var lastErr error
	for i := 0; i < bs.config.App.MaxRetries; i++ {
		bs.logger.CtxInfof(ctx, "Attempt %d/%d for operation: %s", 
			i+1, bs.config.App.MaxRetries, operation)

		// 模拟可能失败的操作 (Simulate potentially failing operation)
		if operation == "flaky" && i < 2 {
			lastErr = errors.Errorf("attempt %d failed", i+1)
			bs.logger.CtxWarnf(ctx, "Operation failed, will retry: %v", lastErr)
			continue
		}

		bs.logger.CtxInfof(ctx, "Operation succeeded on attempt %d", i+1)
		return nil
	}

	// 所有重试都失败了 (All retries failed)
	finalErr := errors.WithCode(
		errors.Wrapf(lastErr, "operation failed after %d retries", bs.config.App.MaxRetries),
		errors.ErrConfigSetup, // 使用示例错误码 (Using example error code)
	)
	bs.logger.CtxErrorf(ctx, "Operation ultimately failed: %v", finalErr)
	return finalErr
}

// createLogOptions 从配置创建日志选项
// (createLogOptions creates log options from configuration)
func createLogOptions(logCfg *config.LogConfig) *log.Options {
	if logCfg == nil {
		return log.NewOptions() // 返回默认选项 (Return default options)
	}

	opts := log.NewOptions()
	opts.Level = logCfg.Level
	opts.Format = logCfg.Format
	opts.EnableColor = logCfg.EnableColor
	opts.Development = logCfg.Development
	opts.Name = logCfg.Name

	// 设置输出路径 (Set output paths)
	if len(logCfg.OutputPaths) > 0 {
		opts.OutputPaths = logCfg.OutputPaths
	} else if logCfg.Output != "" && logCfg.Output != "stdout" {
		opts.OutputPaths = []string{logCfg.Output}
	}

	// 设置错误输出路径 (Set error output paths)
	if len(logCfg.ErrorOutputPaths) > 0 {
		opts.ErrorOutputPaths = logCfg.ErrorOutputPaths
	}

	return opts
}

func main() {
	fmt.Println("=== LMCC Go SDK Basic Usage Example ===")
	fmt.Println("This example demonstrates the integration of config, errors, and log modules.")
	fmt.Println()

	// 1. 加载配置 (Load configuration)
	fmt.Println("1. Loading configuration...")
	var cfg AppConfig

	err := config.LoadConfig(
		&cfg,
		config.WithConfigFile("config.yaml", "yaml"),
		config.WithEnvPrefix("BASIC_EXAMPLE"),
		config.WithEnvVarOverride(true),
	)
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		
		// 演示错误处理 (Demonstrate error handling)
		if coder := errors.GetCoder(err); coder != nil {
			fmt.Printf("Error code: %d, Message: %s\n", coder.Code(), coder.String())
		}
		os.Exit(1)
	}

	fmt.Printf("✓ Configuration loaded successfully\n")
	fmt.Printf("  App Name: %s\n", cfg.App.Name)
	fmt.Printf("  Version: %s\n", cfg.App.Version)
	fmt.Printf("  Environment: %s\n", cfg.App.Environment)
	fmt.Printf("  Timeout: %v\n", cfg.App.Timeout)
	fmt.Printf("  Max Retries: %d\n", cfg.App.MaxRetries)
	fmt.Println()

	// 2. 初始化日志记录器 (Initialize logger)
	fmt.Println("2. Initializing logger...")
	logOpts := createLogOptions(cfg.Log)
	log.Init(logOpts)

	logger := log.Std()
	logger.Info("Logger initialized successfully")
	logger.Infow("Application starting",
		"name", cfg.App.Name,
		"version", cfg.App.Version,
		"environment", cfg.App.Environment,
	)
	fmt.Println("✓ Logger initialized")
	fmt.Println()

	// 3. 创建业务服务 (Create business service)
	fmt.Println("3. Creating business service...")
	service := NewBusinessService(&cfg, logger)
	fmt.Println("✓ Business service created")
	fmt.Println()

	// 4. 演示正常操作 (Demonstrate normal operations)
	fmt.Println("4. Demonstrating normal operations...")
	ctx := context.Background()
	ctx = log.ContextWithRequestID(ctx, "req-12345")

	// 成功的数据处理 (Successful data processing)
	err = service.ProcessData(ctx, "valid-data")
	if err != nil {
		logger.CtxErrorf(ctx, "Unexpected error in normal operation: %v", err)
	} else {
		fmt.Println("✓ Data processing completed successfully")
	}
	fmt.Println()

	// 5. 演示错误处理 (Demonstrate error handling)
	fmt.Println("5. Demonstrating error handling...")
	
	// 无效数据错误 (Invalid data error)
	err = service.ProcessData(ctx, "invalid")
	if err != nil {
		fmt.Printf("✓ Expected error caught: %v\n", err)
		
		// 展示错误的详细信息 (Show detailed error information)
		if coder := errors.GetCoder(err); coder != nil {
			fmt.Printf("  Error Code: %d\n", coder.Code())
			fmt.Printf("  Error Type: %s\n", coder.String())
		}
		
		// 展示堆栈跟踪（简化版）(Show stack trace (simplified))
		fmt.Printf("  Full error: %+v\n", err)
	}
	fmt.Println()

	// 超时错误 (Timeout error)
	err = service.ProcessData(ctx, "timeout")
	if err != nil {
		fmt.Printf("✓ Timeout error caught: %v\n", err)
	}
	fmt.Println()

	// 6. 演示重试逻辑 (Demonstrate retry logic)
	fmt.Println("6. Demonstrating retry logic...")
	
	// 成功的操作 (Successful operation)
	err = service.RetryOperation(ctx, "simple")
	if err != nil {
		fmt.Printf("Unexpected error: %v\n", err)
	} else {
		fmt.Println("✓ Simple operation succeeded")
	}

	// 需要重试的操作 (Operation that needs retries)
	err = service.RetryOperation(ctx, "flaky")
	if err != nil {
		fmt.Printf("Operation failed after retries: %v\n", err)
	} else {
		fmt.Println("✓ Flaky operation succeeded after retries")
	}
	fmt.Println()

	// 7. 演示结构化日志 (Demonstrate structured logging)
	fmt.Println("7. Demonstrating structured logging...")
	logger.Infow("Application metrics",
		"processed_items", 3,
		"errors_count", 2,
		"success_rate", 0.67,
		"avg_processing_time", "100ms",
	)

	logger.Warnw("Resource usage warning",
		"memory_usage", "85%",
		"cpu_usage", "78%",
		"threshold", "80%",
	)
	fmt.Println("✓ Structured logging demonstrated")
	fmt.Println()

	// 8. 清理和结束 (Cleanup and finish)
	fmt.Println("8. Cleaning up...")
	logger.Info("Application shutting down gracefully")
	
	// 同步日志缓冲 (Sync log buffers)
	_ = logger.Sync()
	
	fmt.Println("✓ Application completed successfully")
	fmt.Println()
	fmt.Println("=== Example completed ===")
} 