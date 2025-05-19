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
	merrors "github.com/marmotedu/errors"                  // Import marmotedu/errors for stack trace demo
	"github.com/spf13/viper"                               // Import viper for use in callbacks
	"go.uber.org/zap/zapcore"                              // Import zapcore for log levels
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
		sdklog.Warn("Log configuration section is nil, creating default log options.")
		return sdklog.NewOptions() // 返回默认选项 (Return default options)
	}

	opts := sdklog.NewOptions() // 从默认值开始 (Start with defaults)

	opts.Level = cfg.Level
	opts.Format = cfg.Format
	opts.EnableColor = cfg.EnableColor // Pass enableColor

	// 输出路径设置 (Output path setting)
	if cfg.Output == "stdout" {
		opts.OutputPaths = []string{"stdout"}
	} else if cfg.Output == "stderr" {
		opts.OutputPaths = []string{"stderr"}
	} else if cfg.Output != "" {
		// 如果指定了文件名并且不是 stdout/stderr，则使用该文件
		// (If filename specified and not stdout/stderr, use that file)
		// 优先使用 Filename 字段 (Prioritize Filename field)
		if cfg.Filename != "" {
			opts.OutputPaths = []string{cfg.Filename}
		} else {
			// 否则使用 Output 作为文件名 (Otherwise use Output as filename)
			opts.OutputPaths = []string{cfg.Output}
		}
	} else if len(cfg.OutputPaths) > 0 { // 支持直接在配置文件中定义 OutputPaths 数组
		opts.OutputPaths = cfg.OutputPaths
	}


	// 错误输出路径设置 (Error output path setting)
	// 优先使用 ErrorOutputPaths 数组 (Prioritize ErrorOutputPaths array)
	if len(cfg.ErrorOutputPaths) > 0 {
		opts.ErrorOutputPaths = cfg.ErrorOutputPaths
	} else if cfg.ErrorOutput != "" { // 向下兼容旧的单个 ErrorOutput 字段
		if cfg.ErrorOutput == "stdout" {
			opts.ErrorOutputPaths = []string{"stdout"}
		} else if cfg.ErrorOutput == "stderr" {
			opts.ErrorOutputPaths = []string{"stderr"}
		} else {
			opts.ErrorOutputPaths = []string{cfg.ErrorOutput}
		}
	}
	// 如果两者都为空，NewOptions() 中已设置默认 stderr

	opts.DisableCaller = cfg.DisableCaller
	opts.DisableStacktrace = cfg.DisableStacktrace
	opts.Development = cfg.Development
	opts.Name = cfg.Name

	// 设置日志轮转设置 (Set log rotation settings)
	opts.LogRotateMaxSize = cfg.MaxSize
	opts.LogRotateMaxBackups = cfg.MaxBackups
	opts.LogRotateMaxAge = cfg.MaxAge
	opts.LogRotateCompress = cfg.Compress
	
	// 设置上下文键 (Set context keys)
	if len(cfg.ContextKeys) > 0 {
		opts.ContextKeys = make([]any, len(cfg.ContextKeys))
		for i, k := range cfg.ContextKeys {
			opts.ContextKeys[i] = k
		}
	}

	return opts
}

// simulateErrorCreation 模拟一个会产生嵌套错误的操作
// (simulateErrorCreation simulates an operation that produces a nested error)
func simulateErrorCreation() error {
	err := merrors.New("database connection failed")
	return merrors.Wrap(err, "failed to query user data")
}

// deeperErrorFunction 模拟更深层次的函数调用
// (deeperErrorFunction simulates a deeper function call)
func deeperErrorFunction() error {
	return simulateErrorCreation()
}

// entryPointFunction 模拟错误产生的入口点
// (entryPointFunction simulates the entry point for error generation)
func entryPointFunction() error {
	err := deeperErrorFunction()
	return merrors.Wrap(err, "entry point encountered an issue")
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
	sdklog.Infof("Initial log settings: Level=%s, Format=%s, OutputPaths=%v, EnableColor=%v",
		logOpts.Level, logOpts.Format, logOpts.OutputPaths, logOpts.EnableColor)


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

			sdklog.Infof("SDK Logger re-initialized successfully. New settings: Level=%s, Format=%s, OutputPaths=%v, EnableColor=%v",
				newLogOpts.Level, newLogOpts.Format, newLogOpts.OutputPaths, newLogOpts.EnableColor)
			
			// 根据新配置演示颜色输出 (Demonstrate color output based on new config)
			if newLogOpts.Format == sdklog.FormatText && newLogOpts.EnableColor {
				sdklog.Info("\033[32mThis INFO message should be green if text format and color enabled.\033[0m")
				sdklog.Warn("\033[33mThis WARN message should be yellow if text format and color enabled.\033[0m")
				sdklog.Error("\033[31mThis ERROR message should be red if text format and color enabled.\033[0m")
			} else {
				sdklog.Info("Color output demo: Not text format or color not enabled.")
			}
			return nil
		})

		sdklog.Info("Successfully registered callback for logger updates on config change.")
	}

	// --- 演示日志输出 ---
	// (Demonstrate Log Output)
	sdklog.Info("--- Starting Example Application ---")

	// 演示不同级别的日志 (Demonstrate different log levels)
	// 注意：颜色代码仅在文本格式和 EnableColor: true 时在终端中可见
	// (Note: Color codes are only visible in the terminal with text format and EnableColor: true)
	sdklog.Debug("This is a regular debug message.")
	sdklog.Info("This is a regular info message. If format is text & color enabled, it might be \033[32mgreen\033[0m (or default).")
	sdklog.Warn("This is a regular warning message. If format is text & color enabled, it should be \033[33myellow\033[0m.")
	sdklog.Error("This is a regular error message. If format is text & color enabled, it should be \033[31mred\033[0m.")


	// 演示 marmotedu/errors 的堆栈跟踪
	// (Demonstrate stack trace with marmotedu/errors)
	sdklog.Info("\n--- Demonstrating Stack Trace Logging ---")
	errWithStack := entryPointFunction()
	if errWithStack != nil {
		// 使用 Errorw 记录错误，它会自动处理堆栈信息（如果存在）
		// (Use Errorw to log the error, it handles stack info automatically if present)
		sdklog.Errorw("An error with stack trace occurred", "error", errWithStack, "detail", "simulated from entryPointFunction")
		// 如果是 JSON 格式，堆栈信息会包含在 "errorVerbose" 字段中
		// (If JSON format, stack trace will be in "errorVerbose" field)
		// 你也可以使用 Errorf 并手动格式化:
		// (You can also use Errorf and format manually:)
		// sdklog.Errorf("Detailed error: %+v", errWithStack) // %+v 会打印堆栈
	}


	// 携带上下文信息的日志 (Logging with context information)
	// 创建带有追踪ID和请求ID的上下文 (Create context with Trace ID and Request ID)
	sdklog.Info("\n--- Demonstrating Contextual Logging ---")
	ctx := context.Background()
	ctx = sdklog.ContextWithTraceID(ctx, "trace-abc-123-example")
	ctx = sdklog.ContextWithRequestID(ctx, "req-xyz-789-example")

	// 添加自定义键到上下文中，以便 sdklog 根据 config.yaml 中的 contextKeys 配置进行提取
	// (Add custom key to context so that sdklog can extract it based on contextKeys config in config.yaml)
	ctx = context.WithValue(ctx, "myCustomKey", "This is the value for myCustomKey")
	ctx = context.WithValue(ctx, "user_id", "martin_guo_123") // 演示 config.yaml 中的另一个 contextKey

	// 使用 Ctx 记录带上下文的日志 (Log with context using Ctx)
	sdklog.Ctx(ctx, "This log message includes traceID, requestID and custom keys from context.")
	sdklog.CtxDebugf(ctx, "Contextual debug message with value: %s and custom keys.", "debug_ctx_val")


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
		sdklog.Infof("Log Format (from file/env): %s", AppCfg.Log.Format) // Format is now from file
		sdklog.Infof("Log OutputPaths (from file/env): %v", AppCfg.Log.OutputPaths) // Displaying OutputPaths
		sdklog.Infof("Log EnableColor (from file/env): %t", AppCfg.Log.EnableColor)
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
	sdklog.Info("Modify config.yaml (e.g., log.level, log.format, log.enableColor) to see hot reload in action.")
	sdklog.Info("Press Ctrl+C to exit.")

	// 保持应用运行以观察热重载 (Keep the application running to observe hot reload)
	stopCh := make(chan os.Signal, 1) // Renamed to avoid conflict with stop in ticker goroutine
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	// 模拟一些持续的日志活动 (Simulate some continuous logging activity)
	go func() {
		ticker := time.NewTicker(15 * time.Second) // Increased interval
		defer ticker.Stop()
		looping := true
		for looping {
			select {
			case <-ticker.C:
				sdklog.Debug("Periodic health check log (debug level).")
				// 演示颜色（如果适用）
				// (Demonstrate color (if applicable))
				if sdklog.Std().GetZapLogger().Core().Enabled(zapcore.InfoLevel) { // Check if Info level is enabled (Corrected to zapcore.InfoLevel)
					currentFormat := AppCfg.Log.Format // Access current format from AppCfg
					enableColor := AppCfg.Log.EnableColor // Access current color setting
					if currentFormat == sdklog.FormatText && enableColor {
						sdklog.Info("Periodic \033[36mCYAN\033[0m info tick.")
					} else {
						sdklog.Info("Periodic info tick.")
					}
				}

			case <-stopCh: // Listen to the same channel as main goroutine for shutdown
				sdklog.Debug("Ticker goroutine stopping due to main shutdown signal.")
				looping = false // Exit loop
				return
			}
		}
		sdklog.Debug("Ticker goroutine finished.")
	}()

	<-stopCh // 等待中断信号 (Wait for interrupt signal)

	sdklog.Info("--- Shutting down Example Application ---")
	// 任何清理代码 (Any cleanup code)
	err = sdklog.Sync() // sdklog.Sync() 会返回错误
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error syncing logger: %v\n", err) // Use fmt for final critical errors if logger is failing
	}
	fmt.Println("\n--- Example Finished ---")
}
