# 日志 (`pkg/log`) 使用指南

[Switch to English (切换到英文)](./log_usage_en.md)

本指南解释了如何在 `lmcc-go-sdk` 中使用 `pkg/log` 模块以实现强大且可配置的日志记录功能。

## 1. 功能介绍

`pkg/log` 模块基于 `zap` 提供了一个灵活且强大的日志解决方案。主要特性包括：

-   **多日志级别：** 支持标准级别，如 Debug, Info, Warn, Error, Fatal。
-   **结构化日志：** 以 JSON 或人类可读的文本格式输出日志。
-   **可配置输出：** 将日志定向到 `stdout`、`stderr` 或一个或多个文件。
-   **日志轮转：** 基于大小、时间和备份数量自动轮转日志文件。
-   **热重载：** 当应用程序配置发生变化时，通过与 `pkg/config` 集成，动态更新日志记录器配置（级别、格式、输出）。
-   **上下文日志：** 自动在日志消息中包含来自 `context.Context` 的字段（如追踪 ID、请求 ID）。
-   **调用者信息：** 可选择性地包含日志调用点的文件和行号。
-   **堆栈跟踪：** 自动为错误级别的日志包含堆栈跟踪。

## 2. 集成指南

本节演示了如何在典型的应用程序设置中将 `pkg/log` 与 `pkg/config` 集成。

### 2.1. 配置 (`config.yaml`)

首先，在您的 `config.yaml`（或其他支持的配置文件格式）中定义日志设置。配置文件中的 `log` 部分应对应于 `sdkconfig.LogConfig` 中的字段（该结构体本身映射到 `sdklog.Options` 中的字段）。

```yaml
# config.yaml 示例片段
server:
  port: 9091

log:
  level: "info"       # 例如：debug, info, warn, error, fatal
  format: "json"      # "json" 或 "text"
  output: "stdout"    # "stdout", "stderr", 或文件路径如 "./logs/app.log"
  # filename: "./logs/app.log" # 如果您只指定文件并希望进行轮转，则可替代 'output'
  maxSize: 100        # 轮转前的最大大小 (MB)
  maxBackups: 3       # 保留的旧日志文件最大数量
  maxAge: 7           # 保留旧日志文件的最大天数
  compress: false     # 是否压缩轮转后的文件
  # disableCaller: false
  # disableStacktrace: false
  # enableColor: true # 如果 format 是 text
  # development: false
  # name: "my-app-logger"
  # errorOutputPaths: ["stderr", "./logs/app_error.log"]
  # contextKeys: ["customKey1", "customKey2"] # 如果您有自定义键需要从上下文中提取
```

**注意：** `sdkconfig.LogConfig`（在 `pkg/config/types.go` 中）是 `pkg/config` 用来从您的配置文件中解析 `log` 部分的结构。然后，`examples/simple-config-app/main.go` 中的示例使用一个辅助函数 (`createLogOpts`) 将这些字段映射到 `sdklog.Init()` 所期望的 `sdklog.Options`。

### 2.2. 应用程序代码 (`main.go`)

以下是如何在您的应用程序中初始化和使用日志记录器：

```go
package main

import (
	"context"
	"flag"
	"fmt"
	stdlog "log" // 用于初始设置错误的标准日志库
	"os"
	"os/signal"
	"syscall"
	"time"

	sdkconfig "github.com/lmcc-dev/lmcc-go-sdk/pkg/config" // 导入 SDK 配置包
	sdklog "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"       // 导入 SDK 日志包
	"github.com/spf13/viper"
)

// MyAppConfig 定义您的应用程序配置结构
// (MyAppConfig defines your application's config structure)
type MyAppConfig struct {
	sdkconfig.Config // 嵌入 SDK 基础配置
	// 在此处添加其他自定义应用程序配置字段
}

var AppCfg MyAppConfig

// createLogOpts 将 sdkconfig.LogConfig 转换为 sdklog.Options
// (createLogOpts converts sdkconfig.LogConfig to sdklog.Options)
func createLogOpts(cfg *sdkconfig.LogConfig) *sdklog.Options {
	if cfg == nil {
		return sdklog.NewOptions() // 返回默认选项
	}
	opts := sdklog.NewOptions() // 从默认值开始
	opts.Level = cfg.Level
	opts.Format = cfg.Format
	if cfg.Output == "stdout" {
		opts.OutputPaths = []string{"stdout"}
	} else if cfg.Output == "stderr" {
		opts.OutputPaths = []string{"stderr"}
	} else if cfg.Output != "" {
		if cfg.Filename != "" {
			opts.OutputPaths = []string{cfg.Filename}
		} else {
			opts.OutputPaths = []string{cfg.Output} // 如果 Filename 为空，则使用 Output 作为文件名
		}
	}
	opts.LogRotateMaxSize = cfg.MaxSize
	opts.LogRotateMaxBackups = cfg.MaxBackups
	opts.LogRotateMaxAge = cfg.MaxAge
	opts.LogRotateCompress = cfg.Compress
	// 要使用 sdklog.Options 中的其他字段，请确保它们存在于
	// sdkconfig.LogConfig 中并在此处进行映射。例如：
	// opts.DisableCaller = cfg.DisableCaller // 假设 DisableCaller 存在于 sdkconfig.LogConfig 中
	// opts.Name = cfg.Name                 // 假设 Name 存在于 sdkconfig.LogConfig 中
	return opts
}

func main() {
	configFile := flag.String("config", "config.yaml", "配置文件路径")
	flag.Parse()

	// 加载配置并监视更改
	// (Load configuration and watch for changes)
	configManager, err := sdkconfig.LoadConfigAndWatch(
		&AppCfg,
		sdkconfig.WithConfigFile(*configFile, "yaml"),
		sdkconfig.WithHotReload(true),
	)
	if err != nil {
		stdlog.Fatalf("FATAL: Failed to load initial configuration: %v\n", err)
	}
	stdlog.Println("Initial configuration loaded successfully.")

	if AppCfg.Log == nil {
		stdlog.Fatalln("FATAL: Log configuration section is missing.")
	}

	// 使用加载的配置初始化日志记录器
	// (Initialize logger with loaded config)
	logOpts := createLogOpts(AppCfg.Log)
	sdklog.Init(logOpts) // sdklog.Init 不返回错误
	sdklog.Info("SDK Logger initialized with initial config.")

	// 注册日志配置更改的回调
	// (Register callback for log config changes)
	if configManager != nil {
		configManager.RegisterCallback(func(v *viper.Viper, currentCfgAny any) error {
			currentTypedCfg, ok := currentCfgAny.(*MyAppConfig)
			if !ok {
				sdklog.Error("Config type assertion failed in callback")
				return fmt.Errorf("config type assertion error")
			}
			if currentTypedCfg.Log == nil {
				sdklog.Warn("Log configuration section missing after reload.")
				return fmt.Errorf("log config missing after reload")
			}
			sdklog.Info("Configuration reloaded. Re-initializing logger...")
			newLogOpts := createLogOpts(currentTypedCfg.Log)
			sdklog.Init(newLogOpts)
			sdklog.Infof("SDK Logger re-initialized. New level: %s, New format: %s", 
				newLogOpts.Level, newLogOpts.Format)
			return nil
		})
		sdklog.Info("Callback for logger updates registered.")
	}

	// --- 演示日志记录 --- 
	sdklog.Debug("This is a debug message.")
	sdklog.Infow("User logged in", "username", "martin", "sessionID", 12345)
	sdklog.Warn("Potential issue detected.")
	sdklog.Error("An error occurred", "errorDetails", fmt.Errorf("database connection failed"))

	// 上下文日志记录
	// (Contextual logging)
	ctx := context.Background()
	ctx = sdklog.ContextWithTraceID(ctx, "trace-abc-123")
	ctx = sdklog.ContextWithRequestID(ctx, "req-def-456")
	sdklog.Ctx(ctx, "Processing request with trace and request ID.")

	sdklog.Info("Application running. Modify config.yaml to test hot reload. Press Ctrl+C to exit.")

	// 保持应用运行
	// (Keep app running)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	sdklog.Info("Application shutting down.")
	if err := sdklog.Sync(); err != nil {
		stdlog.Printf("Error syncing logger: %v\n", err)
	}
}
```

### 2.3. 运行示例

1.  将上面的 Go 代码保存为新目录中的 `main.go`。
2.  在同一目录中创建一个包含日志设置的 `config.yaml`。
3.  确保 `lmcc-go-sdk` 在您的 `go.mod` 中。
4.  运行 `go mod tidy`。
5.  运行 `go run main.go -config config.yaml`。
6.  当它运行时，修改 `config.yaml` 中的 `log.level` 或 `log.format` 并保存。观察日志输出以反映更改。

## 3. API 参考

### 3.1.核心函数

-   **`Init(opts *Options)`**: 使用给定的选项初始化或重新初始化全局日志记录器。此函数是线程安全的。
-   **`NewOptions() *Options`**: 返回一个填充了默认值的新 `Options` 结构体。
-   **`NewLogger(opts *Options) Logger`**: 使用给定的选项创建一个新的日志记录器实例。如果您需要多个不同的记录器实例（尽管通常全局记录器 `Std()` 就足够了），这很有用。
-   **`Std() Logger`**: 返回全局单例日志记录器实例。如果尚未调用 `Init`，则使用默认选项进行初始化。
-   **`Sync() error`**: 从全局日志记录器中刷新所有缓冲的日志条目。在应用程序退出之前调用此函数非常重要。

### 3.2. 日志记录方法 (全局和 `Logger` 实例上的方法)

这些可用作全局函数（例如 `sdklog.Info(...)`），它们使用全局记录器；也可用作 `Logger` 实例上的方法（例如 `myLogger.Info(...)`）。

-   `Debug(args ...any)` / `Debugf(template string, args ...any)` / `Debugw(msg string, keysAndValues ...any)`
-   `Info(args ...any)` / `Infof(template string, args ...any)` / `Infow(msg string, keysAndValues ...any)`
-   `Warn(args ...any)` / `Warnf(template string, args ...any)` / `Warnw(msg string, keysAndValues ...any)`
-   `Error(args ...any)` / `Errorf(template string, args ...any)` / `Errorw(msg string, keysAndValues ...any)`
-   `Fatal(args ...any)` / `Fatalf(template string, args ...any)` / `Fatalw(msg string, keysAndValues ...any)` (注意: Fatal 日志记录后会调用 `os.Exit(1)`)

### 3.3. 上下文日志记录

`pkg/log` 模块允许您使用来自 `context.Context` 的数据来丰富您的日志消息。这对于请求追踪以及将日志与特定操作关联起来特别有用。

-   **`ContextWithTraceID(ctx context.Context, traceID string) context.Context`**: 返回设置了 TraceID 的新上下文。`pkg/log` 模块内部使用一个未导出的键类型 (`pkg/log.TraceIDKey`) 来存储此值。
-   **`ContextWithRequestID(ctx context.Context, requestID string) context.Context`**: 返回设置了 RequestID 的新上下文。与 TraceID 类似，使用一个内部键 (`pkg/log.RequestIDKey`)。
-   **`TraceIDFromContext(ctx context.Context) (string, bool)`**: 如果存在，则从上下文中提取 TraceID。
-   **`RequestIDFromContext(ctx context.Context) (string, bool)`**: 如果存在，则从上下文中提取 RequestID。

-   **`Ctx(ctx context.Context, args ...any)`**: 以 InfoLevel 级别记录消息。它会自动从上下文中提取已识别的字段：
    -   **追踪 ID (Trace ID)**: 使用 `TraceIDFromContext` 提取，并以字段键 **`"trace_id"`** 记录。
    -   **请求 ID (Request ID)**: 使用 `RequestIDFromContext` 提取，并以字段键 **`"request_id"`** 记录。
    -   **自定义键 (Custom Keys)**: 如果在记录器初始化期间填充了 `Options.ContextKeys`，这些键也将被提取。对于非简单字符串的自定义上下文键（例如结构体类型），`pkg/log` 通常使用 `fmt.Sprintf("%v", key)` 来生成日志输出中的字段名。

    **示例:**
    ```go
    package main

    import (
    	"context"
    	"fmt"
    	sdklog "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
    	"github.com/google/uuid" // 用于生成唯一ID
    )

    // 定义一个自定义键类型 (作为高级上下文用法示例)
    type myCustomKey struct{}

    func main() {
    	// 初始化记录器 (假设已设置opts，例如输出JSON到stdout)
    	opts := sdklog.NewOptions()
    	opts.Format = "json"
    	opts.OutputPaths = []string{"stdout"}
    	opts.Level = "info"
    	// 要提取 'myCustomKey'，您需要将其添加到 opts.ContextKeys:
    	// opts.ContextKeys = []any{myCustomKey{}} 
    	sdklog.Init(opts)

    	traceID := uuid.NewString()
    	requestID := uuid.NewString()
    	customValue := "my_custom_context_data"

    	ctx := context.Background()
    	ctx = sdklog.ContextWithTraceID(ctx, traceID)
    	ctx = sdklog.ContextWithRequestID(ctx, requestID)
    	ctx = context.WithValue(ctx, myCustomKey{}, customValue) // 存储自定义值

    	// 使用上下文记录日志。`pkg/log` 将自动提取 trace_id 和 request_id。
    	// 如果 myCustomKey{} 已添加到 Options.ContextKeys，它也将被提取。
    	// 日志中 myCustomKey{} 的字段名将是其字符串表示形式，例如 "{}"
    	sdklog.Ctx(ctx, "正在处理带有上下文数据的用户请求。")
    
    	// 您可能如何验证这一点的示例 (概念性的，用于测试):
    	// 日志输出 (JSON):
    	// {
    	//   "level": "info",
    	//   "timestamp": "...",
    	//   "caller": "...",
    	//   "message": "正在处理带有上下文数据的用户请求。",
    	//   "trace_id": "generated-trace-id",  // 自动添加
    	//   "request_id": "generated-request-id", // 自动添加
    	//   "{}": "my_custom_context_data"  // 如果 myCustomKey{} 在 ContextKeys 中，且其字符串形式为 "{}"
    	// }
    }
    ```

-   其他级别也提供了类似的上下文感知日志记录函数：
    -   `CtxDebugf(ctx context.Context, template string, args ...interface{})`
    -   `CtxInfof(ctx context.Context, template string, args ...interface{})`
    -   `CtxWarnf(ctx context.Context, template string, args ...interface{})`
    -   `CtxErrorf(ctx context.Context, template string, args ...interface{})`
    -   `CtxFatalf(ctx context.Context, template string, args ...interface{})` (也会退出程序)
    -   `CtxPanicf(ctx context.Context, template string, args ...interface{})` (也会触发panic)

### 3.4. 日志记录器操作

-   **`WithName(name string) Logger`**: 返回一个新的日志记录器实例，其名称附加了指定的名称。
-   **`WithValues(keysAndValues ...any) Logger`**: 返回一个新的日志记录器实例，其结构化上下文中添加了给定的键值对。

### 3.5. `Options` 结构体

`pkg/log/Options` 中的关键字段（有关所有字段和默认值，请参阅 `pkg/log/options.go`）：

-   `Level`: `string` (例如 "debug", "info")
-   `Format`: `string` ("json" 或 "text")
-   `OutputPaths`: `[]string` (例如 `["stdout"]`, `["./logs/app.log"]`)
-   `ErrorOutputPaths`: `[]string` (用于内部记录器错误)
-   `DisableCaller`: `bool`
-   `DisableStacktrace`: `bool`
-   `EnableColor`: `bool` (用于文本格式)
-   `Development`: `bool`
-   `Name`: `string` (记录器名称)
-   `LogRotateMaxSize`: `int` (MB)
-   `LogRotateMaxBackups`: `int`
-   `LogRotateMaxAge`: `int` (天)
-   `LogRotateCompress`: `bool`
-   `ContextKeys`: `[]any` (从上下文中提取的自定义键列表)

## 4. 相关的 Makefile 命令

-   `make test-unit PKG=./pkg/log`: 运行 `pkg/log` 的单元测试。
-   `make cover PKG=./pkg/log`: 运行 `pkg/log` 的带覆盖率的单元测试。
-   `make test-integration`: 运行所有集成测试 (如果 `pkg/log` 有特定的集成测试，或作为更广泛集成测试的一部分进行测试，则相关)。
-   `make lint`: 对代码库（包括 `pkg/log`）进行 lint 检查。

本指南全面概述了 `pkg/log` 模块的使用。有关更多详细信息，请参阅 `pkg/log` 目录中的源代码和特定函数文档。
