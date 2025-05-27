# 输出格式

本文档详细介绍了日志模块支持的三种输出格式：文本、JSON 和键值对格式。

## 格式概览

日志模块支持三种主要的输出格式：

1. **Text（文本）** - 人类可读的格式，适合开发和调试
2. **JSON** - 结构化格式，适合生产环境和日志聚合
3. **KeyValue（键值对）** - 传统的键值对格式，适合传统日志系统

## Text 格式

### 基本结构

文本格式提供人类可读的日志输出，特别适合开发环境。

```go
opts := &log.Options{
    Format: "text",
    EnableColor: true,  // 在终端中启用颜色
}
log.Init(opts)
```

### 输出示例

```go
log.Info("应用程序启动")
log.Infow("用户登录", "user_id", 123, "username", "john_doe")
log.Errorw("数据库连接失败", "error", "连接超时", "retry_count", 3)
```

**输出：**
```
2024-01-15T10:30:45.123Z	INFO	main.go:25	应用程序启动
2024-01-15T10:30:45.124Z	INFO	main.go:26	用户登录	{"user_id": 123, "username": "john_doe"}
2024-01-15T10:30:45.125Z	ERROR	main.go:27	数据库连接失败	{"error": "连接超时", "retry_count": 3}
```

### 颜色支持

在支持颜色的终端中，不同级别的日志会显示不同颜色：

- **DEBUG** - 蓝色
- **INFO** - 绿色
- **WARN** - 黄色
- **ERROR** - 红色
- **FATAL** - 红色（粗体）
- **PANIC** - 红色（粗体）

```go
opts := &log.Options{
    Format:      "text",
    EnableColor: true,
}
```

### 调用者信息

文本格式可以包含调用者信息（文件名和行号）：

```go
opts := &log.Options{
    Format:        "text",
    DisableCaller: false,  // 显示调用者信息
}
```

**输出：**
```
2024-01-15T10:30:45.123Z	INFO	main.go:25	消息内容
```

## JSON 格式

### 基本结构

JSON 格式提供结构化的日志输出，非常适合生产环境和自动化日志处理。

```go
opts := &log.Options{
    Format: "json",
}
log.Init(opts)
```

### 输出示例

```go
log.Info("应用程序启动")
log.Infow("用户登录", "user_id", 123, "username", "john_doe", "ip", "192.168.1.100")
log.Errorw("数据库连接失败", "error", "连接超时", "retry_count", 3, "duration", "5.2s")
```

**输出：**
```json
{"level":"info","timestamp":"2024-01-15T10:30:45.123Z","caller":"main.go:25","message":"应用程序启动"}
{"level":"info","timestamp":"2024-01-15T10:30:45.124Z","caller":"main.go:26","message":"用户登录","user_id":123,"username":"john_doe","ip":"192.168.1.100"}
{"level":"error","timestamp":"2024-01-15T10:30:45.125Z","caller":"main.go:27","message":"数据库连接失败","error":"连接超时","retry_count":3,"duration":"5.2s"}
```

### JSON 字段说明

标准 JSON 日志包含以下字段：

- **level** - 日志级别（debug, info, warn, error, fatal, panic）
- **timestamp** - ISO 8601 格式的时间戳
- **caller** - 调用者信息（文件名:行号）
- **message** - 日志消息
- **其他字段** - 通过 `log.Infow()` 等方法添加的结构化字段

### 上下文字段

使用上下文日志记录时，上下文字段会自动包含在 JSON 输出中：

```go
ctx := context.Background()
ctx = log.WithValues(ctx, "request_id", "req-123", "user_id", 456)

log.InfoContext(ctx, "处理请求")
```

**输出：**
```json
{"level":"info","timestamp":"2024-01-15T10:30:45.123Z","caller":"main.go:25","message":"处理请求","request_id":"req-123","user_id":456}
```

### 错误堆栈跟踪

在 JSON 格式中，堆栈跟踪会作为字符串字段包含：

```go
opts := &log.Options{
    Format:            "json",
    DisableStacktrace: false,
    StacktraceLevel:   "error",
}
```

**输出：**
```json
{"level":"error","timestamp":"2024-01-15T10:30:45.123Z","caller":"main.go:25","message":"严重错误","stacktrace":"goroutine 1 [running]:\nmain.main()\n\t/path/to/main.go:25 +0x123"}
```

## KeyValue 格式

### 基本结构

键值对格式提供传统的日志格式，适合与现有日志系统集成。

```go
opts := &log.Options{
    Format: "keyvalue",
}
log.Init(opts)
```

### 输出示例

```go
log.Info("应用程序启动")
log.Infow("用户登录", "user_id", 123, "username", "john_doe")
log.Errorw("数据库连接失败", "error", "连接超时", "retry_count", 3)
```

**输出：**
```
timestamp=2024-01-15T10:30:45.123Z level=info caller=main.go:25 message="应用程序启动"
timestamp=2024-01-15T10:30:45.124Z level=info caller=main.go:26 message="用户登录" user_id=123 username=john_doe
timestamp=2024-01-15T10:30:45.125Z level=error caller=main.go:27 message="数据库连接失败" error="连接超时" retry_count=3
```

### 字段格式规则

在键值对格式中：

- **字符串值** - 如果包含空格或特殊字符，会用引号包围
- **数字值** - 直接输出，不加引号
- **布尔值** - 输出为 `true` 或 `false`
- **特殊字符** - 在值中的特殊字符会被转义

### 复杂值处理

对于复杂的数据类型，键值对格式会进行适当的序列化：

```go
log.Infow("用户信息",
    "user", map[string]interface{}{
        "id":   123,
        "name": "张三",
        "tags": []string{"admin", "active"},
    },
    "timestamp", time.Now(),
)
```

**输出：**
```
timestamp=2024-01-15T10:30:45.123Z level=info caller=main.go:25 message="用户信息" user="{\"id\":123,\"name\":\"张三\",\"tags\":[\"admin\",\"active\"]}" timestamp=2024-01-15T10:30:45.123Z
```

## 格式比较

### 可读性

| 格式 | 人类可读性 | 机器可读性 | 适用场景 |
|------|------------|------------|----------|
| Text | ⭐⭐⭐⭐⭐ | ⭐⭐ | 开发、调试 |
| JSON | ⭐⭐ | ⭐⭐⭐⭐⭐ | 生产、日志聚合 |
| KeyValue | ⭐⭐⭐ | ⭐⭐⭐⭐ | 传统系统、监控 |

### 性能

| 格式 | 序列化速度 | 文件大小 | 解析速度 |
|------|------------|----------|----------|
| Text | 快 | 中等 | 慢 |
| JSON | 中等 | 大 | 快 |
| KeyValue | 快 | 小 | 中等 |

### 存储效率

```go
// 相同日志消息的不同格式大小比较
message := "用户登录"
fields := map[string]interface{}{
    "user_id": 123,
    "username": "john_doe",
    "ip": "192.168.1.100",
    "timestamp": time.Now(),
}

// Text: ~120 字节
// JSON: ~180 字节  
// KeyValue: ~100 字节
```

## 自定义格式化

### 时间戳格式

虽然不能直接自定义时间戳格式，但可以通过后处理来调整：

```go
// 所有格式都使用 ISO 8601 格式
// 2024-01-15T10:30:45.123Z
```

### 字段顺序

在 JSON 和 KeyValue 格式中，字段顺序是固定的：

1. 标准字段（level, timestamp, caller, message）
2. 上下文字段（按添加顺序）
3. 结构化字段（按添加顺序）

## 格式选择指南

### 开发环境

```go
opts := &log.Options{
    Level:       "debug",
    Format:      "text",
    OutputPaths: []string{"stdout"},
    EnableColor: true,
}
```

**优势：**
- 易于阅读和调试
- 颜色编码提高可读性
- 快速识别问题

### 生产环境

```go
opts := &log.Options{
    Level:        "info",
    Format:       "json",
    OutputPaths:  []string{"/var/log/app.log"},
    EnableColor:  false,
}
```

**优势：**
- 结构化数据便于查询
- 与日志聚合系统兼容
- 支持复杂的过滤和分析

### 传统系统集成

```go
opts := &log.Options{
    Level:       "warn",
    Format:      "keyvalue",
    OutputPaths: []string{"/var/log/syslog"},
}
```

**优势：**
- 与现有工具兼容
- 紧凑的输出格式
- 易于解析和处理

## 格式转换

### 运行时切换格式

```go
// 可以在运行时动态切换格式
func switchToJSONFormat() {
    opts := &log.Options{
        Level:  "info",
        Format: "json",
    }
    log.Init(opts)
}

func switchToTextFormat() {
    opts := &log.Options{
        Level:       "debug",
        Format:      "text",
        EnableColor: true,
    }
    log.Init(opts)
}
```

### 多格式输出

虽然单个日志器实例只能使用一种格式，但可以配置多个输出：

```go
// 这需要在应用程序级别实现
// 例如：同时输出到控制台（text）和文件（json）
```

## 最佳实践

### 1. 环境特定格式

```go
func getLogFormat() string {
    env := os.Getenv("APP_ENV")
    switch env {
    case "development":
        return "text"
    case "production":
        return "json"
    case "testing":
        return "keyvalue"
    default:
        return "text"
    }
}
```

### 2. 性能优化

```go
// 高性能场景：使用 keyvalue 格式
opts := &log.Options{
    Level:             "error",
    Format:            "keyvalue",
    DisableCaller:     true,
    DisableStacktrace: true,
}
```

### 3. 调试友好

```go
// 调试场景：使用 text 格式
opts := &log.Options{
    Level:       "debug",
    Format:      "text",
    EnableColor: true,
}
```

### 4. 生产监控

```go
// 生产监控：使用 json 格式
opts := &log.Options{
    Level:  "info",
    Format: "json",
    OutputPaths: []string{
        "/var/log/app.log",
        "stdout",  // 用于容器日志收集
    },
}
```

## 下一步

- [上下文日志记录](04_context_logging.md) - 掌握上下文感知日志记录
- [性能优化](05_performance.md) - 优化日志性能
- [最佳实践](06_best_practices.md) - 生产就绪模式 