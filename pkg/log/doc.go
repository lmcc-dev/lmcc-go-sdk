/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

/*
Package log provides a flexible, high-performance, and structured logging solution
for Go applications, built on top of Uber's zap library.
(log 包为 Go 应用程序提供了一个灵活、高性能且结构化的日志解决方案，它构建于 Uber 的 zap 库之上。)

It aims to offer a comprehensive set of features suitable for modern microservices
and applications, including dynamic configuration and context-aware logging.
(它旨在提供一套适用于现代微服务和应用程序的全面功能，包括动态配置和上下文感知日志记录。)

Core Features:
(核心特性：)

  - Structured Logging: Leverages zap for efficient structured logging.
    (结构化日志：利用 zap 实现高效的结构化日志记录。)
  - Multiple Log Levels: Supports standard levels (Debug, Info, Warn, Error, Fatal, Panic).
    (多日志级别：支持标准级别（Debug, Info, Warn, Error, Fatal, Panic）。)
  - Configurable Output Formats: Supports JSON and human-readable console (text) formats.
    (可配置的输出格式：支持 JSON 和人类可读的控制台（文本）格式。)
  - Multiple Output Paths: Can write logs to stdout, stderr, and one or more files simultaneously.
    (多输出路径：可以同时将日志写入 stdout、stderr 以及一个或多个文件。)
  - Log Rotation: Built-in support for log rotation based on size, age, and number of backups, with optional compression.
    (日志轮转：内置支持基于大小、保留时间、备份数量的日志轮转，并可选压缩。)
  - Dynamic Configuration: Integrates with `pkg/config` for hot-reloading of logging configurations (level, format, output paths, etc.)
    without application restart.
    (动态配置：与 `pkg/config` 集成，支持在不重启应用的情况下热重载日志配置（级别、格式、输出路径等）。)
  - Context-Aware Logging: Automatically extracts and logs predefined (TraceID, RequestID) and user-specified
    values from `context.Context`.
    (上下文感知日志记录：自动从 `context.Context` 中提取并记录预定义（TraceID、RequestID）和用户指定的值。)
  - Global and Instance Loggers: Provides a global singleton logger for convenience and the ability to create
    multiple independent logger instances with different configurations.
    (全局和实例记录器：为方便起见提供了全局单例记录器，并能创建具有不同配置的多个独立记录器实例。)
  - Helper functions for adding common IDs (TraceID, RequestID) to context.
    (辅助函数，用于将常用 ID (TraceID, RequestID) 添加到 context 中。)

Basic Usage with Global Logger:
(全局记录器基本用法：)

	import (
		"context" // Required for context examples (上下文示例需要)
		"fmt"     // Required for Sync error handling example (Sync 错误处理示例需要)
		"os"      // Required for Sync error handling example (Sync 错误处理示例需要)

		merrors "github.com/marmotedu/errors" // Use marmotedu/errors (使用 marmotedu/errors)
		"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
	)

	func main() {
		// 1. Initialize the global logger (usually once at application startup)
		// (1. 初始化全局记录器（通常在应用程序启动时执行一次）)
		opts := log.NewOptions()
		opts.Level = "debug" // Set desired level (设置所需级别)
		opts.Format = log.FormatText // Use text format for console (控制台使用文本格式)
		opts.EnableColor = true     // Enable color for readability (启用颜色以提高可读性)
		log.Init(opts)

		// Ensure logs are flushed before application exits
		// (确保在应用程序退出前刷写日志)
		defer func() {
			if err := log.Sync(); err != nil {
				// Handle sync error, e.g., print to stderr
				// (处理同步错误，例如，打印到 stderr)
				fmt.Fprintf(os.Stderr, "Failed to sync global logger: %v\n", err)
			}
		}()

		// 2. Use the global logger
		// (2. 使用全局记录器)
		log.Debug("This is a debug message.")
		log.Info("Service started successfully.", "port", 8080, "version", "1.0.2")
		log.Warnf("Configuration value %s is deprecated.", "old_setting")
		log.Errorw("Failed to process request",
			"request_id", "req-123",
			"error", merrors.New("something went wrong from marmotedu/errors"), // Correctly use merrors.New
		)

		// 3. Context-aware logging with the global logger
		// (3. 使用全局记录器进行上下文感知日志记录)
		ctx := context.Background()
		ctx = log.ContextWithTraceID(ctx, "trace-abc-123")
		ctx = log.ContextWithRequestID(ctx, "req-xyz-789")
		log.Ctx(ctx, "Processing incoming request with context.")
		log.Ctxf(ctx, "User %s performed action %s.", "user123", "update_profile")
		log.Ctxw(ctx, "Order processed", "order_id", "order-456", "customer_id", "cust-789")
	}

Creating and Using a Local Logger Instance:
(创建和使用本地记录器实例：)

	func doSomething() {
		// Create a specific logger for this component or task
		// (为此组件或任务创建一个特定的记录器)
		opts := log.NewOptions()
		opts.Name = "my-component"
		opts.Level = "info"
		opts.OutputPaths = []string{"./my-component.log"}
		opts.Format = log.FormatJSON
		logger := log.NewLogger(opts)
		defer func() { _ = logger.Sync() }() // Sync this specific logger instance (同步此特定的记录器实例)

		logger.Info("Component initialized.")
		// ... component logic ...
		logger.Warnw("Potential issue detected", "issue_code", 1001)
	}

Integration with pkg/config for Hot-Reloading:
(与 pkg/config 集成以实现热重载：)

To enable dynamic reconfiguration of logging, ensure your main application configuration struct
includes a section for logging that matches the structure of `log.Options` (using `mapstructure` tags),
and then register the log package's hot-reload handler with the `config.Manager`.
(要启用日志的动态重新配置，请确保您的主应用程序配置结构体包含一个与 `log.Options` 结构（使用 `mapstructure` 标签）
匹配的日志记录部分，然后向 `config.Manager` 注册日志包的热重载处理程序。)

Example structure in your `config.Config`:
(您的 `config.Config` 中的示例结构：)

	type Config struct {
	    // ... other config fields ...
	    Log log.Options `mapstructure:"log"` // Log configuration section (日志配置部分)
	}

In your main setup:
(在您的主要设置代码中：)

	var appConfig MyMainAppConfig // Your application's main config struct (您的应用程序主配置结构体)

	// Load initial config and start watching (加载初始配置并开始监视)
	// cm is your config.Manager instance (cm 是您的 config.Manager 实例)
	cm, err := sdkconfig.LoadConfigAndWatch(&appConfig, ...your config options...)
	if err != nil {
	    // Handle error
	}

	// Initialize global logger with settings from the loaded config
	// (使用从加载的配置中的设置初始化全局记录器)
	log.Init(&appConfig.Log) // Pass the log options from your loaded config (传递从加载的配置中获取的日志选项)

	// Register log package for hot-reloading with the config manager
	// (向配置管理器注册日志包以进行热重载)
	log.RegisterConfigHotReload(cm)

Now, when the "log" section of your configuration file changes, the global logger
will be automatically reconfigured.
(现在，当配置文件的 "log" 部分发生更改时，全局记录器将自动重新配置。)
*/
package log

// Import necessary packages for examples (if not already present globally for the package doc)
// (为示例导入必要的包（如果包文档的全局范围尚不存在）)

// For context examples (用于上下文示例)
// For error examples (用于错误示例)

// For os.Stderr in Sync example (用于 Sync 示例中的 os.Stderr)

// sdkconfig "github.com/lmcc-dev/lmcc-go-sdk/pkg/config" // Uncomment if sdkconfig alias is used in examples above
// and not defined package-wide.
// (如果上面的示例中使用了 sdkconfig 别名且未在包范围内定义，请取消注释。)