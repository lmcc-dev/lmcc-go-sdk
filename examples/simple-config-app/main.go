/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Example application demonstrating the usage of lmcc-go-sdk/pkg/config and pkg/log.
 */

package main

import (
	"context" // Import context for logger
	"flag"    // Import the flag package
	"fmt"
	"log" // Standard log for initial setup errors
	"os"
	"os/signal" // For graceful shutdown
	"syscall"   // For system signals
	"time"      // For demonstrating timed logging

	sdkconfig "github.com/lmcc-dev/lmcc-go-sdk/pkg/config" // Import the SDK config package
	sdklog "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"       // Import the SDK log package
	"github.com/spf13/viper"                               // Import viper for use in callbacks
)

// MyAppConfig 嵌入了 SDK 的 Config 并添加了自定义字段
// (MyAppConfig embeds the SDK's Config and adds custom fields)
type MyAppConfig struct {
	sdkconfig.Config                      // 嵌入 SDK 基础配置 (Embed SDK base config)
	CustomFeature    *CustomFeatureConfig `mapstructure:"customFeature"`
}

// CustomFeatureConfig 示例应用的自定义配置部分
// (Custom configuration section for the example application)
type CustomFeatureConfig struct {
	APIKey    string `mapstructure:"apiKey"`
	RateLimit int    `mapstructure:"rateLimit"`
	Enabled   bool   `mapstructure:"enabled"`
}

// AppCfg 是此应用的全局配置实例
// (AppCfg is the global configuration instance for this application)
var AppCfg MyAppConfig

// createLogOpts 从 sdkconfig.LogConfig 创建 sdklog.Options
// (createLogOpts creates sdklog.Options from sdkconfig.LogConfig)
func createLogOpts(cfg *sdkconfig.LogConfig) *sdklog.Options {
	if cfg == nil {
		return sdklog.NewOptions() // 返回默认选项 (Return default options)
	}
	
	// 使用 log.LogConfig 字段创建一个新的 sdklog.Options 实例
	// (Create a new sdklog.Options instance using log.LogConfig fields)
	opts := sdklog.NewOptions() // 从默认值开始 (Start with defaults)
	
	// 设置字段 (Set fields)
	opts.Level = cfg.Level  
	opts.Format = cfg.Format
	
	// 输出路径设置：从单个输出转换为输出路径数组
	// (Output path setting: convert from single output to output paths array)
	if cfg.Output == "stdout" {
		opts.OutputPaths = []string{"stdout"}
	} else if cfg.Output == "stderr" {
		opts.OutputPaths = []string{"stderr"}
	} else if cfg.Output != "" {
		// 如果指定了文件名并且不是 stdout/stderr，则使用该文件
		// (If filename specified and not stdout/stderr, use that file)
		if cfg.Filename != "" {
			opts.OutputPaths = []string{cfg.Filename}
		} else {
			// 否则使用默认文件名 (Otherwise use default filename)
			opts.OutputPaths = []string{cfg.Output}
		}
	}
	
	// 设置日志轮转设置 (Set log rotation settings)
	opts.LogRotateMaxSize = cfg.MaxSize
	opts.LogRotateMaxBackups = cfg.MaxBackups
	opts.LogRotateMaxAge = cfg.MaxAge
	opts.LogRotateCompress = cfg.Compress
	
	return opts
}

func main() {
	// 使用 flag 包处理命令行参数来指定配置文件路径
	// (Use the flag package to handle command-line arguments for specifying the config file path)
	configFile := flag.String("config", "config.yaml", "Path to the configuration file (e.g., -config /path/to/config.yaml)")
	flag.Parse() // 解析命令行参数 (Parse command-line arguments)

	configFilePath := *configFile

	fmt.Printf("Attempting to load configuration from: %s (use -config flag to override)\n", configFilePath)

	// 加载配置并启用热重载 (Load configuration and enable hot reload)
	configManager, err := sdkconfig.LoadConfigAndWatch(
		&AppCfg,
		sdkconfig.WithConfigFile(configFilePath, "yaml"), // "yaml" is explicit, "" would infer
		sdkconfig.WithHotReload(true),
	)
	if err != nil {
		log.Fatalf("FATAL: Failed to load initial configuration: %v\n", err) // Use standard log
		os.Exit(1)
	}
	fmt.Println("\n--- Initial Configuration Loaded Successfully ---")

	// --- 初始化 SDK 日志记录器 ---
	// (Initialize SDK Logger)
	// 使用加载的配置初始化日志记录器
	// (Initialize the logger using the loaded configuration)
	if AppCfg.Log == nil {
		log.Fatalln("FATAL: Log configuration section is missing in the loaded config.") // Use standard log
		os.Exit(1)
	}
	
	// 将 config.LogConfig 转换为 log.Options
	// (Convert config.LogConfig to log.Options)
	logOpts := createLogOpts(AppCfg.Log)
	
	// 初始化日志记录器 - 注意：sdklog.Init 不返回错误
	// (Initialize logger - Note: sdklog.Init does not return an error)
	sdklog.Init(logOpts) 
	sdklog.Info("SDK Logger initialized successfully with initial config.")

	// --- 注册配置更改回调以更新日志记录器 ---
	// (Register a config change callback to update the logger)
	if configManager != nil {
		// RegisterCallback 不返回错误 (RegisterCallback does not return an error)
		configManager.RegisterCallback(func(v *viper.Viper, currentCfgAny any) error {
			currentCfg, ok := currentCfgAny.(*MyAppConfig)
			if !ok {
				sdklog.Error("Config type assertion failed in callback")
				return fmt.Errorf("config type assertion failed in callback")
			}

			if currentCfg.Log == nil {
				sdklog.Warn("Log configuration section is missing after reload, logger not updated.")
				return fmt.Errorf("log configuration section is missing after reload")
			}

			sdklog.Info("Configuration reloaded. Attempting to re-initialize logger...")
			
			// 重新初始化日志记录器 (Re-initialize the logger)
			// 将 config.LogConfig 转换为 log.Options
			// (Convert config.LogConfig to log.Options)
			newLogOpts := createLogOpts(currentCfg.Log)
			
			// 初始化日志记录器 - 注意：sdklog.Init 不返回错误
			// (Initialize logger - Note: sdklog.Init does not return an error)
			sdklog.Init(newLogOpts)
			
			sdklog.Infof("SDK Logger re-initialized successfully. New level: %s, New format: %s", 
				currentCfg.Log.Level, currentCfg.Log.Format)
			return nil
		})
		
		sdklog.Info("Successfully registered callback for logger updates on config change.")
	}

	// --- 演示日志输出 ---
	// (Demonstrate Log Output)
	sdklog.Info("--- Starting Example Application ---")

	// 演示不同级别的日志 (Demonstrate different log levels)
	sdklog.Debug("This is a debug message.")
	sdklog.Infow("This is an info message.", "user", "Martin", "userID", 123)
	sdklog.Warn("This is a warning message.")
	sdklog.Error("This is an error message:", fmt.Errorf("something went wrong"))
	// sdklog.Fatal("This is a fatal message.") // This would exit

	// 携带上下文信息的日志 (Logging with context information)
	// 创建带有追踪ID和请求ID的上下文 (Create context with Trace ID and Request ID)
	ctx := context.Background()
	ctx = sdklog.ContextWithTraceID(ctx, "xyz123")
	ctx = sdklog.ContextWithRequestID(ctx, "req-987")
	
	// 使用 Ctx 记录带上下文的日志 (Log with context using Ctx)
	sdklog.Ctx(ctx, "This log message includes traceID and requestID from context.")

	// --- 访问 SDK 基础配置 ---
	// (Accessing SDK Base Configuration)
	sdklog.Info("\n[SDK Base Config]")
	if AppCfg.Server != nil {
		sdklog.Infof("Server Port (from file/env): %d", AppCfg.Server.Port)
		sdklog.Infof("Server Mode (from file/env): %s", AppCfg.Server.Mode)
		sdklog.Infof("Server Host (default): %s", AppCfg.Server.Host)
	} else {
		sdklog.Warn("Server configuration section not loaded.")
	}

	if AppCfg.Log != nil {
		sdklog.Infof("Log Level (from file/env): %s", AppCfg.Log.Level)
		sdklog.Infof("Log Format (default): %s", AppCfg.Log.Format)
		sdklog.Infof("Log Output (default): %s", AppCfg.Log.Output)
	} else {
		sdklog.Warn("Log configuration section not loaded.")
	}

	if AppCfg.Database != nil {
		sdklog.Infof("Database Host (from file/env): %s", AppCfg.Database.Host)
		sdklog.Infof("Database User (from file/env): %s", AppCfg.Database.User)
		sdklog.Infof("Database Name (from file/env): %s", AppCfg.Database.DBName)
		sdklog.Infof("Database Type (default): %s", AppCfg.Database.Type)
	} else {
		sdklog.Warn("Database configuration section not loaded.")
	}

	// --- 访问自定义配置 ---
	// (Accessing Custom Configuration)
	sdklog.Info("\n[Custom App Config]")
	if AppCfg.CustomFeature != nil {
		sdklog.Infof("Custom Feature API Key: %s", AppCfg.CustomFeature.APIKey)
		sdklog.Infof("Custom Feature Rate Limit: %d", AppCfg.CustomFeature.RateLimit)
		sdklog.Infof("Custom Feature Enabled: %t", AppCfg.CustomFeature.Enabled)
	} else {
		sdklog.Warn("CustomFeature configuration section not loaded or defined in config file.")
	}

	sdklog.Info("\n--- Example Application Running ---")
	sdklog.Info("Modify config.yaml (e.g., log.level, log.format) to see hot reload in action.")
	sdklog.Info("Press Ctrl+C to exit.")

	// 保持应用运行以观察热重载 (Keep the application running to observe hot reload)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// 模拟一些持续的日志活动 (Simulate some continuous logging activity)
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				sdklog.Debug("Periodic health check log (debug level).")
			case <-stop: // Avoid leaking goroutine if main exits
				return
			}
		}
	}()

	<-stop // 等待中断信号 (Wait for interrupt signal)

	sdklog.Info("--- Shutting down Example Application ---")
	// 任何清理代码 (Any cleanup code)
	err = sdklog.Sync() // sdklog.Sync() 会返回错误
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error syncing logger: %v\n", err) // Use fmt for final critical errors if logger is failing
	}
	fmt.Println("\n--- Example Finished ---")
}
