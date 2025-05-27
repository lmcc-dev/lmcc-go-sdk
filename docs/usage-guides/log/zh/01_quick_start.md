# 快速开始指南

本指南将帮助您在几分钟内开始使用日志模块。

## 安装

日志模块是 lmcc-go-sdk 的一部分。如果您还没有安装，请将其添加到您的项目中：

```bash
go mod init your-project
go get github.com/lmcc-dev/lmcc-go-sdk
```

## 基本用法

### 步骤 1：初始化日志器

最简单的开始方式是使用默认配置：

```go
package main

import (
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

func main() {
    // 使用默认设置初始化
    // - 级别: info
    // - 格式: text（控制台中带颜色）
    // - 输出: stdout
    log.Init(nil)
    
    log.Info("日志器初始化成功！")
}
```

### 步骤 2：基本日志记录

```go
func basicLogging() {
    // 不同的日志级别
    log.Debug("这是调试消息")   // 默认 'info' 级别下不会显示
    log.Info("应用程序已启动")        // 会显示
    log.Warn("这是一个警告")          // 会显示
    log.Error("出现了错误")      // 会显示
    
    // Fatal 和 Panic（谨慎使用）
    // log.Fatal("严重错误")  // 调用 os.Exit(1)
    // log.Panic("恐慌消息")   // 调用 panic()
}
```

### 步骤 3：结构化日志记录

真正的威力来自使用键值对的结构化日志记录：

```go
func structuredLogging() {
    // 使用 'w' 变体进行结构化日志记录
    log.Infow("用户登录",
        "user_id", 12345,
        "username", "john_doe",
        "ip", "192.168.1.100",
        "timestamp", time.Now(),
    )
    
    log.Errorw("数据库连接失败",
        "error", "连接超时",
        "database", "users_db",
        "retry_count", 3,
        "duration", "5.2s",
    )
    
    log.Warnw("检测到高内存使用",
        "memory_usage", "85%",
        "threshold", "80%",
        "process_id", 1234,
    )
}
```

## 自定义配置

### 步骤 1：配置输出格式

```go
func customConfiguration() {
    opts := &log.Options{
        Level:  "debug",  // 显示调试消息
        Format: "json",   // 使用 JSON 格式而不是文本
    }
    
    log.Init(opts)
    
    log.Debug("这个调试消息现在会出现")
    log.Infow("用户操作",
        "action", "login",
        "user_id", 123,
    )
}
```

### 步骤 2：配置输出目标

```go
func multipleOutputs() {
    opts := &log.Options{
        Level:       "info",
        Format:      "json",
        OutputPaths: []string{
            "stdout",              // 控制台输出
            "/var/log/app.log",    // 文件输出
        },
    }
    
    log.Init(opts)
    
    log.Info("此消息同时输出到控制台和文件")
}
```

### 步骤 3：启用日志轮转

```go
func withLogRotation() {
    opts := &log.Options{
        Level:       "info",
        Format:      "json",
        OutputPaths: []string{"stdout", "/var/log/app.log"},
        
        // 日志轮转设置
        LogRotateMaxSize:    100,  // 每个文件 100 MB
        LogRotateMaxBackups: 5,    // 保留 5 个备份文件
        LogRotateMaxAge:     30,   // 保留文件 30 天
        LogRotateCompress:   true, // 压缩旧文件
    }
    
    log.Init(opts)
    
    log.Info("启用轮转的日志记录")
}
```

## 上下文感知日志记录

### 基本上下文使用

```go
import (
    "context"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

func contextLogging() {
    // 创建带有日志字段的上下文
    ctx := context.Background()
    ctx = log.WithValues(ctx, "request_id", "req-123", "user_id", 456)
    
    // 使用上下文记录日志 - 字段会自动包含
    log.InfoContext(ctx, "正在处理用户请求")
    log.WarnContext(ctx, "请求耗时超过预期")
    
    // 向上下文添加更多字段
    ctx = log.WithValues(ctx, "operation", "database_query")
    log.InfoContext(ctx, "执行数据库查询")
}
```

### HTTP 请求上下文

```go
import (
    "net/http"
    "github.com/google/uuid"
)

func httpHandler(w http.ResponseWriter, r *http.Request) {
    // 创建带有关联 ID 的请求上下文
    requestID := uuid.New().String()
    ctx := log.WithValues(r.Context(), 
        "request_id", requestID,
        "method", r.Method,
        "path", r.URL.Path,
        "remote_addr", r.RemoteAddr,
    )
    
    log.InfoContext(ctx, "请求开始")
    
    // 处理请求...
    processRequest(ctx)
    
    log.InfoContext(ctx, "请求完成")
}

func processRequest(ctx context.Context) {
    // 此函数中的所有日志都将包含请求上下文
    log.InfoContext(ctx, "验证请求")
    log.InfoContext(ctx, "查询数据库")
    log.InfoContext(ctx, "生成响应")
}
```

## 输出格式示例

### JSON 格式输出

使用 `Format: "json"` 时，日志如下所示：

```json
{
  "level": "info",
  "timestamp": "2024-01-15T10:30:45.123Z",
  "caller": "main.go:25",
  "message": "用户登录",
  "user_id": 12345,
  "username": "john_doe",
  "ip": "192.168.1.100"
}
```

### 文本格式输出

使用 `Format: "text"`（默认）时，日志如下所示：

```
2024-01-15T10:30:45.123Z	INFO	main.go:25	用户登录	{"user_id": 12345, "username": "john_doe", "ip": "192.168.1.100"}
```

### 键值对格式输出

使用 `Format: "keyvalue"` 时，日志如下所示：

```
timestamp=2024-01-15T10:30:45.123Z level=info caller=main.go:25 message="用户登录" user_id=12345 username=john_doe ip=192.168.1.100
```

## 完整示例

这是一个演示各种功能的完整示例：

```go
package main

import (
    "context"
    "time"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

func main() {
    // 配置日志器
    opts := &log.Options{
        Level:            "debug",
        Format:           "json",
        OutputPaths:      []string{"stdout", "app.log"},
        EnableColor:      false,  // 文件输出禁用颜色
        DisableCaller:    false,  // 包含调用者信息
        DisableStacktrace: false, // 包含错误堆栈跟踪
        StacktraceLevel:  "error", // 仅为错误及以上级别显示堆栈跟踪
    }
    
    log.Init(opts)
    
    // 基本日志记录
    log.Info("应用程序启动中")
    
    // 结构化日志记录
    log.Infow("配置已加载",
        "config_file", "app.yaml",
        "log_level", opts.Level,
        "output_format", opts.Format,
    )
    
    // 上下文日志记录
    ctx := context.Background()
    ctx = log.WithValues(ctx, "component", "database", "version", "1.0.0")
    
    log.InfoContext(ctx, "连接到数据库")
    
    // 模拟一些工作
    time.Sleep(100 * time.Millisecond)
    
    log.InfoContext(ctx, "数据库连接已建立")
    
    // 带堆栈跟踪的错误日志记录
    log.Errorw("处理用户数据失败",
        "user_id", 123,
        "error", "验证失败",
        "field", "email",
    )
    
    log.Info("应用程序启动完成")
}
```

## 与配置模块集成

对于生产应用程序，您通常希望从文件加载日志配置：

```go
import (
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

type AppConfig struct {
    Log log.Options `mapstructure:"log"`
    // ... 其他配置
}

func main() {
    var cfg AppConfig
    
    // 加载配置
    err := config.LoadConfig(&cfg,
        config.WithConfigFile("config.yaml", ""),
        config.WithEnvPrefix("APP"),
    )
    if err != nil {
        panic(err)
    }
    
    // 使用加载的配置初始化日志记录
    log.Init(&cfg.Log)
    
    log.Info("应用程序启动，使用基于配置的日志记录")
}
```

配合 `config.yaml` 文件：

```yaml
log:
  level: "info"
  format: "json"
  output_paths: ["stdout", "/var/log/app.log"]
  enable_color: false
  disable_caller: false
  log_rotate_max_size: 100
  log_rotate_max_backups: 5
  log_rotate_max_age: 30
  log_rotate_compress: true
```

## 入门最佳实践

1. **从简单开始**：从 `log.Init(nil)` 开始，根据需要添加配置
2. **使用结构化日志记录**：优先使用 `log.Infow()` 而不是 `log.Info()` 以获得更好的可搜索性
3. **包含上下文**：使用上下文日志记录进行请求跟踪和关联
4. **选择正确的级别**：使用适当的日志级别（开发用 debug，生产用 info）
5. **配置轮转**：在生产环境中始终为文件输出设置日志轮转

## 下一步

- [配置选项](02_configuration_options.md) - 了解所有可用选项
- [输出格式](03_output_formats.md) - 理解不同的输出格式
- [上下文日志记录](04_context_logging.md) - 掌握上下文感知日志记录
- [最佳实践](06_best_practices.md) - 生产就绪模式 