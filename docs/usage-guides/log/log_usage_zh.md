# 日志 (`pkg/log`) 使用指南

[Switch to English (切换到英文)](./log_usage_en.md)

本指南解释了如何在 `lmcc-go-sdk` 中使用 `pkg/log` 模块以实现强大且可配置的日志记录功能。

## 1. 功能介绍

`pkg/log` 模块基于 `zap` 提供了一个灵活且强大的日志解决方案。主要特性包括：

-   **多日志级别：** 支持标准级别，如 Debug, Info, Warn, Error, Fatal。
-   **结构化日志：** 以 JSON 或人类可读的文本格式输出日志。
-   **可配置输出：** 将日志定向到 `stdout`、`stderr` 或一个或多个文件。通过 `outputPaths` 和 `errorOutputPaths` 进行配置。
-   **日志轮转：** 基于大小、时间和备份数量自动轮转日志文件。
-   **热重载：** 当应用程序配置发生变化时，通过与 `pkg/config` 集成，动态更新日志记录器配置（级别、格式、输出等）。
-   **上下文日志：** 自动在日志消息中包含来自 `context.Context` 的字段（如追踪 ID、请求 ID 以及通过 `contextKeys` 配置的自定义键）。
-   **调用者信息：** 可选择性地包含日志调用点的文件和行号（可通过 `disableCaller: false` 启用）。
-   **错误与堆栈跟踪：**
    -   对于 Error 及以上级别的日志，`zap` 默认会尝试附加堆栈跟踪（可通过 `disableStacktrace: true` 禁用 `zap` 自身的堆栈）。
    -   当与 `github.com/marmotedu/errors` 集成时，如果记录的错误是由 `marmotedu/errors` 包装的，其详细堆栈跟踪会包含在 JSON 日志的 `errorVerbose` 字段中。
-   **彩色输出：** 当格式为 `text` 时，可为不同的日志级别启用彩色输出（通过 `enableColor: true`）。
-   **开发模式：** `development: true` 可配置更易于开发时阅读的日志格式和行为。
-   **记录器命名：** 可以通过 `name` 字段为记录器实例指定一个名称。

## 2. 集成指南

本节演示了如何在典型的应用程序设置中将 `pkg/log` 与 `pkg/config` 集成。
完整的可运行示例请参考 `examples/simple-config-app` 目录。

### 2.1. 配置 (`config.yaml`)

首先，在您的 `config.yaml`（或其他支持的配置文件格式）中定义日志设置。配置文件中的 `log` 部分应对应于 `sdkconfig.LogConfig` 中的字段。

```yaml
# config.yaml 示例片段
server:
  port: 9091

log:
  level: "debug"       # 例如：debug, info, warn, error, fatal
  format: "json"      # \"json\" 或 \"text\"
  enableColor: true   # 当 format 为 \"text\" 且终端支持时，启用颜色输出
  outputPaths:        # 日志输出路径列表
    - "stdout"
    - "./logs/app.log"
  errorOutputPaths:   # 内部错误和 PANIC 日志的输出路径列表 (默认为 stderr)
    - "stderr"
    - "./logs/app_error.log"
  # filename: \"./logs/app.log\" # 旧的单文件输出方式，如果使用 outputPaths，则此项可忽略或用于特定轮转配置
  maxSize: 100        # 轮转前的最大大小 (MB)
  maxBackups: 3       # 保留的旧日志文件最大数量
  maxAge: 7           # 保留旧日志文件的最大天数
  compress: false     # 是否压缩轮转后的文件
  disableCaller: false # false 表示输出调用者信息 (文件名和行号)
  disableStacktrace: false # false 表示 zap 会尝试为 Error 及以上级别日志附加堆栈 (非 marmotedu/errors 的堆栈)
  development: false  # true 会启用更适合开发的日志配置 (例如，堆栈跟踪更易读)
  name: "example-app" # 日志记录器的名称
  contextKeys:        # 从 context.Context 中提取并包含在日志中的额外键列表
    - "customKey1"
    - "user_id"
```

**注意：**
-   `sdkconfig.LogConfig`（在 `pkg/config/types.go` 中）是 `pkg/config` 用来从您的配置文件中解析 `log` 部分的结构。
-   `examples/simple-config-app/main.go` 中的示例使用一个辅助函数 (`createLogOpts`) 将这些字段映射到 `sdklog.Init()` 所期望的 `sdklog.Options`。确保您的 `createLogOpts` 或类似的转换逻辑能够处理所有您希望从配置中读取的字段。

### 2.1.1. JSON 日志格式键名

当 `format` 设置为 `"json"` 时，为了优化性能和减小日志体积，日志记录器会对核心字段使用简洁的键名。在 JSON 输出中，这些核心字段的默认键名如下：

-   **`L`**: 日志级别 (例如："DEBUG", "INFO")
-   **`T`**: 时间戳 (例如："2023-10-27T10:00:00.123Z")
-   **`M`**: 日志消息
-   **`N`**: 日志记录器名称 (如果通过 `log.name` 配置)
-   **`C`**: 调用者信息 (例如："module/file.go:123")
-   **`stacktrace`**: 堆栈跟踪 (针对 ERROR, PANIC, FATAL 级别，或当 `errorVerbose` 存在时)

上下文相关的字段（如 `trace_id`, `request_id` 以及在 `contextKeys` 中指定的键）将保留其配置或定义时的名称。`errorVerbose` 字段在存在时也会保留其名称。

### 2.2. 应用程序代码 (`main.go` - 关键部分)

以下是如何在您的应用程序中初始化和使用日志记录器的关键部分。请参考 `examples/simple-config-app/main.go` 获取完整代码。

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
	merrors "github.com/marmotedu/errors" // 用于演示堆栈跟踪
)

// MyAppConfig 定义您的应用程序配置结构
type MyAppConfig struct {
	sdkconfig.Config // 嵌入 SDK 基础配置
	// 在此处添加其他自定义应用程序配置字段
}

var AppCfg MyAppConfig

// createLogOpts 将 sdkconfig.LogConfig 转换为 sdklog.Options
func createLogOpts(cfg *sdkconfig.LogConfig) *sdklog.Options {
	if cfg == nil {
		sdklog.Warn("Log configuration section is nil, creating default log options.")
		return sdklog.NewOptions()
	}
	opts := sdklog.NewOptions() // 从默认值开始

	opts.Level = cfg.Level
	opts.Format = cfg.Format
	opts.EnableColor = cfg.EnableColor // 新增：传递颜色配置

	// 输出路径
	if len(cfg.OutputPaths) > 0 {
		opts.OutputPaths = cfg.OutputPaths
	} else if cfg.Output != "" { // 向后兼容旧的 'output' 字段
	if cfg.Output == "stdout" {
		opts.OutputPaths = []string{"stdout"}
	} else if cfg.Output == "stderr" {
		opts.OutputPaths = []string{"stderr"}
		} else {
			// 如果 Filename 也存在，它可能在旧逻辑中优先于 Output
			// 这里简化为：如果 OutputPaths 为空，则 Output (如果非 stdout/stderr) 作为单文件路径
			filePath := cfg.Output
			if cfg.Filename != "" { // 如果 Filename 存在，则优先使用
				filePath = cfg.Filename
			}
			opts.OutputPaths = []string{filePath}
		}
	} else {
		opts.OutputPaths = []string{"stdout"} // 默认
	}
	
	// 错误输出路径
	if len(cfg.ErrorOutputPaths) > 0 {
		opts.ErrorOutputPaths = cfg.ErrorOutputPaths
	} else if cfg.ErrorOutput != "" { // 向后兼容旧的 'errorOutput' 字段
	    if cfg.ErrorOutput == "stdout" {
			opts.ErrorOutputPaths = []string{"stdout"}
		} else if cfg.ErrorOutput == "stderr" {
			opts.ErrorOutputPaths = []string{"stderr"}
		} else {
			opts.ErrorOutputPaths = []string{cfg.ErrorOutput}
		}
	} else {
	    opts.ErrorOutputPaths = []string{"stderr"} // 默认
	}

	opts.LogRotateMaxSize = cfg.MaxSize
	opts.LogRotateMaxBackups = cfg.MaxBackups
	opts.LogRotateMaxAge = cfg.MaxAge
	opts.LogRotateCompress = cfg.Compress

	opts.DisableCaller = cfg.DisableCaller         // 新增
	opts.DisableStacktrace = cfg.DisableStacktrace // 新增
	opts.Development = cfg.Development           // 新增
	opts.Name = cfg.Name                           // 新增

	if len(cfg.ContextKeys) > 0 { // 新增
		opts.ContextKeys = make([]any, len(cfg.ContextKeys))
		for i, k := range cfg.ContextKeys {
			opts.ContextKeys[i] = k
		}
	}
	return opts
}

// deeperErrorFunction 模拟产生错误的函数
func deeperErrorFunction() error {
    return merrors.Wrap(merrors.New("底层的数据库错误"), "服务层处理失败")
}

func main() {
	configFile := flag.String("config", "config.yaml", "配置文件路径")
	flag.Parse()

	configManager, err := sdkconfig.LoadConfigAndWatch(
		&AppCfg,
		sdkconfig.WithConfigFile(*configFile, "yaml"),
		sdkconfig.WithHotReload(true),
	)
	if err != nil {
		stdlog.Fatalf("FATAL: Failed to load initial configuration: %v\\n", err)
	}
	stdlog.Println("Initial configuration loaded successfully.")

	if AppCfg.Log == nil {
		stdlog.Fatalln("FATAL: Log configuration section is missing.")
	}

	logOpts := createLogOpts(AppCfg.Log)
	sdklog.Init(logOpts)
	sdklog.Info("SDK Logger initialized with initial config.")
	sdklog.Infof("Initial log settings: Level=%s, Format=%s, OutputPaths=%v, EnableColor=%t",
		logOpts.Level, logOpts.Format, logOpts.OutputPaths, logOpts.EnableColor)

	if configManager != nil {
		configManager.RegisterCallback(func(v *viper.Viper, currentCfgAny any) error {
			currentTypedCfg, ok := currentCfgAny.(*MyAppConfig)
			if !ok { /* ... error handling ... */ return fmt.Errorf("config type error")}
			if currentTypedCfg.Log == nil { /* ... error handling ... */ return fmt.Errorf("log config missing")}
			
			sdklog.Info("Configuration reloaded. Re-initializing logger...")
			newLogOpts := createLogOpts(currentTypedCfg.Log)
			sdklog.Init(newLogOpts) // Re-initialize
			sdklog.Infof("SDK Logger re-initialized. New settings: Level=%s, Format=%s, OutputPaths=%v, EnableColor=%t",
				newLogOpts.Level, newLogOpts.Format, newLogOpts.OutputPaths, newLogOpts.EnableColor)
			// Demo color output after re-init
			if newLogOpts.Format == sdklog.FormatText && newLogOpts.EnableColor {
				sdklog.Info("\\033[32mThis INFO message should be green.\\033[0m")
				sdklog.Warn("\\033[33mThis WARN message should be yellow.\\033[0m")
			}
			return nil
		})
		sdklog.Info("Callback for logger updates registered.")
	}

	// --- 演示日志记录 --- 
	sdklog.Debug("This is a debug message.")
	sdklog.Infow("User logged in", "username", "martin", "sessionID", 12345)
	
	// 演示错误和堆栈跟踪
	errWithStack := deeperErrorFunction()
	sdklog.Errorw("发生了一个错误，包含来自marmotedu/errors的堆栈跟踪", "error", errWithStack, "relevant_id", "id-123")

	// 上下文日志记录
	ctx := context.Background()
	ctx = sdklog.ContextWithTraceID(ctx, "trace-abc-123")
	ctx = sdklog.ContextWithRequestID(ctx, "req-def-456")
	// 假设 "customKey1" 在 ContextKeys 中配置
	ctx = context.WithValue(ctx, "customKey1", "customValueForLog")
	sdklog.Ctx(ctx, "Processing request with trace, request ID, and customKey1.")

	sdklog.Info("Application running. Modify config.yaml to test hot reload. Press Ctrl+C to exit.")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	sdklog.Info("Application shutting down.")
	if err := sdklog.Sync(); err != nil {
		stdlog.Printf("Error syncing logger: %v\\n", err)
	}
}
```

### 2.3. 运行示例

1.  参考 `examples/simple-config-app` 目录中的完整示例代码和 `config.yaml`。
2.  运行 `go run examples/simple-config-app/main.go -config examples/simple-config-app/config.yaml`。
3.  当它运行时，修改 `config.yaml` 中的 `log.level` 或 `log.format` 等并保存。观察日志输出以反映更改。

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
    -   **自定义键 (Custom Keys)**: 如果在记录器初始化期间通过 `Options.ContextKeys` (对应配置文件中的 `log.contextKeys`) 配置了键列表，这些键对应的值也将从上下文中提取并记录。

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

## 4. 高级特性

### 4.1. 错误堆栈跟踪

当您记录一个 `error` 类型的对象时：

-   **`zap` 默认行为**: 如果日志级别是 Error 或更高级别，并且 `log.disableStacktrace` 配置为 `false` (默认)，`zap` 会尝试附加一个堆栈跟踪。这个堆栈跟踪通常记录在名为 `stacktrace` 的字段中（JSON格式）。
-   **与 `marmotedu/errors` 集成**: 如果您记录的错误是使用 `github.com/marmotedu/errors` 包装的，那么由 `marmotedu/errors` 捕获的更详细的堆栈跟踪信息会在 JSON 日志中以 `errorVerbose` 字段名记录。这通常比 `zap` 自身的堆栈更易读，因为它专注于错误产生的路径。
    -   要利用此特性，确保您的错误是通过 `merrors.New`, `merrors.Errorf`, `merrors.Wrap` 等函数创建的。
    -   使用 `sdklog.Errorw("message", "error", yourMarmotError)` 或 `sdklog.Errorf("message: %+v", yourMarmotError)` (当格式为text时，`%+v`会打印堆栈) 来记录。
    -   `log.disableStacktrace: true` 配置**不会**禁用 `marmotedu/errors` 的 `errorVerbose` 堆栈。

**示例 (`config.yaml` 中 `log.format: "json"`)**:

```go
// 在你的代码中:
import merrors "github.com/marmotedu/errors"
// ...
func doSomething() error {
    return merrors.New("something bad happened")
}

err := doSomething()
if err != nil {
    sdklog.Errorw("Operation failed", "error", err, "user_id", 123)
}
```

**可能的 JSON 输出片段**:
```json
{
  "L": "ERROR", "T": "...", "C": "...", "N": "example-app",
  "M": "Operation failed",
  "user_id": 123,
  "error": "something bad happened",
  "errorVerbose": "something bad happened\\n    main.doSomething\\n        /path/to/your/app/main.go:XX\\n    main.main\\n        /path/to/your/app/main.go:YY\\n    ...",
  "stacktrace": "main.main\\n\\t/path/to/your/app/main.go:YY\\nruntime.main..." // zap 的堆栈 (如果 disableStacktrace=false)
}
```
如上所示，`errorVerbose` 提供了由 `marmotedu/errors` 生成的堆栈信息。

### 4.2. 彩色日志输出

为了在开发或本地调试时获得更好的可读性，您可以启用彩色日志输出。

-   **配置**:
    -   在 `config.yaml` 的 `log` 部分设置 `format: "text"`。
    -   设置 `enableColor: true`。
-   **效果**: 当日志输出到支持 ANSI 颜色转义序列的终端时，不同的日志级别会以不同的颜色显示（例如，Error 为红色，Warn 为黄色，Info 为绿色等）。
-   **注意**: JSON 格式的日志本质上是结构化数据，不包含颜色信息。此特性主要用于文本格式的控制台输出。

**示例 (`config.yaml`)**:
```yaml
log:
  format: "text"
  enableColor: true
  level: "debug"
  outputPaths: ["stdout"]
```
当您运行应用并查看终端输出时，您会看到彩色的日志条目。
`examples/simple-config-app` 中也演示了如何在配置热重载后检查颜色设置并打印特定颜色的消息。

## 5. 日志轮转

`pkg/log` 支持通过 `lumberjack.v2` 实现日志轮转。相关配置项：
-   `maxSize`: 单个日志文件的最大大小 (MB)。
-   `maxBackups`: 保留的旧日志文件的最大数量。
-   `maxAge`: 保留旧日志文件的最大天数 (天)。
-   `compress`: 是否压缩轮转后的日志文件 (例如，使用 gzip)。

当日志输出到文件时 (例如，`outputPaths: ["./logs/app.log"]`)，这些设置会自动生效。

## 6. 热重载

通过与 `pkg/config` 模块集成 (`sdkconfig.LoadConfigAndWatch` 和 `configManager.RegisterCallback`)，`pkg/log` 可以在运行时动态地响应配置文件的更改。这意味着您可以修改 `config.yaml` 中的 `log` 部分（例如，改变 `level`、`format`、`enableColor` 或输出路径），应用程序的日志行为会相应更新，无需重启。

完整的演示可以在 `examples/simple-config-app/main.go` 中找到。

## 7. 最佳实践

-   **同步日志 (Sync Logs)**: 在应用程序退出前，务必调用 `sdklog.Sync()` 来确保所有缓冲的日志都已写入。
-   **结构化日志优先**: 对于生产环境，优先使用 `json` 格式，因为它更易于机器解析和集成到日志管理系统中。
-   **适当的日志级别**: 根据信息的重要性和频率选择合适的日志级别。避免在生产中过多使用 Debug 级别。
-   **使用上下文日志**: 对于与请求相关的日志，始终传递 `context.Context` 以包含追踪信息。
-   **错误处理**: 使用 `Errorw` 或 `Errorf` 记录错误，并尽可能提供上下文信息。如果使用了 `marmotedu/errors`，其堆栈信息会被自动捕获和记录。
-   **配置管理**: 通过配置文件管理日志设置，并利用热重载功能进行动态调整。

本指南提供了 `pkg/log` 模块的全面概述。更多详细信息和高级用法，请参阅源代码和 `zap` 的官方文档。
