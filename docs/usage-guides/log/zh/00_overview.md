# 日志模块概述

`pkg/log` 模块为 Go 应用程序提供了全面且高性能的日志记录解决方案。基于经过验证的日志库构建，它提供结构化日志记录、多种输出格式、上下文感知日志记录和热重载功能。

## 核心特性

### 🚀 **高性能**
- 针对高吞吐量应用程序进行优化
- 最小化内存分配开销
- 可配置的日志级别以减少不必要的处理
- 针对不同输出格式的高效序列化

### 📊 **多种输出格式**
- **文本格式**：适用于开发的人类可读输出
- **JSON 格式**：适用于生产环境和日志聚合的结构化输出
- **键值格式**：适用于传统系统的传统格式

### 🎯 **上下文感知日志记录**
- 通过 Go 的 `context.Context` 自动传播上下文
- 请求跟踪和关联 ID
- 跨函数调用的结构化字段继承

### ⚙️ **灵活配置**
- 特定环境的配置
- 支持动态更新的热重载
- 与 `pkg/config` 模块集成
- 全面的配置验证

### 🔄 **日志轮转**
- 基于大小、时间或数量的自动日志文件轮转
- 归档日志文件的压缩
- 可配置的保留策略

## 架构

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   应用程序      │───▶│   日志模块       │───▶│   输出目标      │
│                 │    │                  │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                              │                         │
                              ▼                         ▼
                       ┌──────────────┐         ┌──────────────┐
                       │ 配置热重载   │         │ • stdout     │
                       │              │         │ • 文件       │
                       └──────────────┘         │ • 网络       │
                                               └──────────────┘
```

## 核心组件

### 日志接口
主要的日志接口提供不同日志级别和结构化日志记录的方法：

```go
// 基本日志记录方法
log.Debug("调试消息")
log.Info("信息消息")
log.Warn("警告消息")
log.Error("错误消息")

// 带字段的结构化日志记录
log.Infow("用户操作", "user_id", 123, "action", "login")

// 上下文感知日志记录
log.InfoContext(ctx, "请求已处理")
```

### 配置系统
日志记录各个方面的全面配置选项：

```go
type Options struct {
    Level            string   `mapstructure:"level"`
    Format           string   `mapstructure:"format"`
    OutputPaths      []string `mapstructure:"output_paths"`
    EnableColor      bool     `mapstructure:"enable_color"`
    DisableCaller    bool     `mapstructure:"disable_caller"`
    // ... 更多选项
}
```

### 上下文集成
与 Go 的上下文系统无缝集成，用于请求跟踪：

```go
ctx := log.WithValues(ctx, "request_id", "req-123")
log.InfoContext(ctx, "正在处理请求") // 自动包含 request_id
```

## 使用场景

### 开发环境
- 带颜色的人类可读文本格式
- 调试级别日志记录以获取详细信息
- 调用者信息便于调试

### 生产环境
- 用于结构化日志记录的 JSON 格式
- 用于性能的信息级别日志记录
- 日志轮转和压缩
- 与日志聚合系统集成

### 高性能应用程序
- 仅错误级别日志记录
- 禁用调用者信息和堆栈跟踪
- 用于最小开销的键值格式
- 异步日志记录模式

## 集成点

### 与 pkg/config 集成
```go
type AppConfig struct {
    Log log.Options `mapstructure:"log"`
}

// 自动配置加载和热重载
cm, err := config.LoadConfigAndWatch(&cfg, ...)
log.Init(&cfg.Log)
```

### 与 Web 框架集成
```go
// Gin 中间件示例
func LoggingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx := log.WithValues(c.Request.Context(),
            "request_id", generateRequestID(),
            "method", c.Request.Method,
            "path", c.Request.URL.Path,
        )
        c.Request = c.Request.WithContext(ctx)
        c.Next()
    }
}
```

### 与错误处理集成
```go
import "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"

if err := someOperation(); err != nil {
    log.ErrorContext(ctx, "操作失败",
        "error", err,
        "error_code", errors.GetCoder(err),
    )
}
```

## 性能特征

| 特性 | 影响 | 建议 |
|------|------|------|
| 日志级别 | 高 | 生产环境使用 "info" 或更高级别 |
| 输出格式 | 中等 | 生产环境使用 JSON，开发环境使用文本 |
| 调用者信息 | 中等 | 高性能场景中禁用 |
| 堆栈跟踪 | 高 | 仅用于错误级别及以上 |
| 上下文字段 | 低 | 对请求跟踪高效 |

## 快速开始

1. **[快速开始](01_quick_start.md)** - 几分钟内启动并运行
2. **[配置选项](02_configuration_options.md)** - 详细配置指南
3. **[输出格式](03_output_formats.md)** - 选择正确的格式
4. **[上下文日志记录](04_context_logging.md)** - 掌握请求跟踪
5. **[性能优化](05_performance.md)** - 针对您的用例进行优化
6. **[最佳实践](06_best_practices.md)** - 生产就绪模式

## 下一步

- 探索[快速开始指南](01_quick_start.md)以获得即时的实践体验
- 查看[配置选项](02_configuration_options.md)以进行详细设置
- 查看[最佳实践](06_best_practices.md)以获得生产部署指导 