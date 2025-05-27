# 配置选项

本文档详细描述了日志模块的所有配置选项以及如何有效使用它们。

## Options 结构

日志模块使用 `log.Options` 结构来配置所有日志行为：

```go
type Options struct {
    // 基本配置
    Level            string   `mapstructure:"level" default:"info"`
    Format           string   `mapstructure:"format" default:"text"`
    OutputPaths      []string `mapstructure:"output_paths" default:"[\"stdout\"]"`
    ErrorOutputPaths []string `mapstructure:"error_output_paths" default:"[\"stderr\"]"`
    
    // 显示选项
    EnableColor      bool `mapstructure:"enable_color" default:"true"`
    DisableCaller    bool `mapstructure:"disable_caller" default:"false"`
    DisableStacktrace bool `mapstructure:"disable_stacktrace" default:"false"`
    StacktraceLevel  string `mapstructure:"stacktrace_level" default:"error"`
    
    // 日志轮转配置
    LogRotateMaxSize    int  `mapstructure:"log_rotate_max_size" default:"100"`
    LogRotateMaxBackups int  `mapstructure:"log_rotate_max_backups" default:"5"`
    LogRotateMaxAge     int  `mapstructure:"log_rotate_max_age" default:"30"`
    LogRotateCompress   bool `mapstructure:"log_rotate_compress" default:"true"`
}
```

## 基本配置选项

### Level（日志级别）

控制哪些日志消息会被输出。

**可用值：**
- `"debug"` - 显示所有日志消息
- `"info"` - 显示 info、warn、error、fatal、panic
- `"warn"` - 显示 warn、error、fatal、panic
- `"error"` - 显示 error、fatal、panic
- `"fatal"` - 显示 fatal、panic
- `"panic"` - 仅显示 panic

**示例：**
```go
opts := &log.Options{
    Level: "debug",  // 开发环境
}

opts := &log.Options{
    Level: "warn",   // 生产环境
}
```

### Format（输出格式）

控制日志消息的输出格式。

**可用值：**
- `"text"` - 人类可读的文本格式（默认）
- `"json"` - 结构化 JSON 格式
- `"keyvalue"` - 键值对格式

**示例：**

```go
// 文本格式（适合开发）
opts := &log.Options{
    Format: "text",
}
// 输出：2024-01-15T10:30:45.123Z	INFO	main.go:25	用户登录	{"user_id": 123}

// JSON 格式（适合生产和日志聚合）
opts := &log.Options{
    Format: "json",
}
// 输出：{"level":"info","timestamp":"2024-01-15T10:30:45.123Z","caller":"main.go:25","message":"用户登录","user_id":123}

// 键值对格式（适合传统日志系统）
opts := &log.Options{
    Format: "keyvalue",
}
// 输出：timestamp=2024-01-15T10:30:45.123Z level=info caller=main.go:25 message="用户登录" user_id=123
```

### OutputPaths（输出路径）

指定日志消息的输出目标。

**可用值：**
- `"stdout"` - 标准输出
- `"stderr"` - 标准错误
- 文件路径 - 如 `"/var/log/app.log"`

**示例：**
```go
// 仅输出到控制台
opts := &log.Options{
    OutputPaths: []string{"stdout"},
}

// 同时输出到控制台和文件
opts := &log.Options{
    OutputPaths: []string{"stdout", "/var/log/app.log"},
}

// 仅输出到文件
opts := &log.Options{
    OutputPaths: []string{"/var/log/app.log"},
}
```

### ErrorOutputPaths（错误输出路径）

指定错误级别日志的输出目标。

**示例：**
```go
opts := &log.Options{
    OutputPaths:      []string{"stdout", "/var/log/app.log"},
    ErrorOutputPaths: []string{"stderr", "/var/log/error.log"},
}
```

## 显示选项

### EnableColor（启用颜色）

在终端输出中启用颜色编码。

**示例：**
```go
// 开发环境：启用颜色
opts := &log.Options{
    Format:      "text",
    EnableColor: true,
}

// 生产环境或文件输出：禁用颜色
opts := &log.Options{
    Format:      "json",
    EnableColor: false,
}
```

### DisableCaller（禁用调用者信息）

控制是否在日志中包含调用者信息（文件名和行号）。

**示例：**
```go
// 包含调用者信息（默认）
opts := &log.Options{
    DisableCaller: false,
}
// 输出：2024-01-15T10:30:45.123Z	INFO	main.go:25	消息

// 不包含调用者信息
opts := &log.Options{
    DisableCaller: true,
}
// 输出：2024-01-15T10:30:45.123Z	INFO	消息
```

### DisableStacktrace（禁用堆栈跟踪）

控制是否在错误日志中包含堆栈跟踪。

**示例：**
```go
// 启用堆栈跟踪（默认）
opts := &log.Options{
    DisableStacktrace: false,
    StacktraceLevel:   "error",
}

// 禁用堆栈跟踪
opts := &log.Options{
    DisableStacktrace: true,
}
```

### StacktraceLevel（堆栈跟踪级别）

指定从哪个级别开始包含堆栈跟踪。

**可用值：**
- `"debug"`, `"info"`, `"warn"`, `"error"`, `"fatal"`, `"panic"`

**示例：**
```go
// 仅在 error 及以上级别显示堆栈跟踪
opts := &log.Options{
    StacktraceLevel: "error",
}

// 在 warn 及以上级别显示堆栈跟踪
opts := &log.Options{
    StacktraceLevel: "warn",
}
```

## 日志轮转配置

### LogRotateMaxSize（最大文件大小）

单个日志文件的最大大小（MB）。

**示例：**
```go
opts := &log.Options{
    OutputPaths:         []string{"/var/log/app.log"},
    LogRotateMaxSize:    100,  // 100 MB
}
```

### LogRotateMaxBackups（最大备份数量）

保留的旧日志文件数量。

**示例：**
```go
opts := &log.Options{
    LogRotateMaxBackups: 5,  // 保留 5 个备份文件
}
```

### LogRotateMaxAge（最大保留天数）

日志文件的最大保留天数。

**示例：**
```go
opts := &log.Options{
    LogRotateMaxAge: 30,  // 保留 30 天
}
```

### LogRotateCompress（压缩旧文件）

是否压缩轮转的日志文件。

**示例：**
```go
opts := &log.Options{
    LogRotateCompress: true,  // 压缩旧文件以节省空间
}
```

## 配置示例

### 开发环境配置

```go
func developmentConfig() *log.Options {
    return &log.Options{
        Level:             "debug",
        Format:            "text",
        OutputPaths:       []string{"stdout"},
        EnableColor:       true,
        DisableCaller:     false,
        DisableStacktrace: false,
        StacktraceLevel:   "error",
    }
}
```

### 生产环境配置

```go
func productionConfig() *log.Options {
    return &log.Options{
        Level:               "info",
        Format:              "json",
        OutputPaths:         []string{"/var/log/app.log"},
        ErrorOutputPaths:    []string{"/var/log/error.log"},
        EnableColor:         false,
        DisableCaller:       false,
        DisableStacktrace:   false,
        StacktraceLevel:     "error",
        LogRotateMaxSize:    100,
        LogRotateMaxBackups: 10,
        LogRotateMaxAge:     30,
        LogRotateCompress:   true,
    }
}
```

### 测试环境配置

```go
func testConfig() *log.Options {
    return &log.Options{
        Level:             "warn",
        Format:            "text",
        OutputPaths:       []string{"stdout"},
        EnableColor:       false,
        DisableCaller:     true,
        DisableStacktrace: true,
    }
}
```

### 高性能配置

```go
func highPerformanceConfig() *log.Options {
    return &log.Options{
        Level:             "error",  // 仅记录错误
        Format:            "json",   // 更快的序列化
        OutputPaths:       []string{"/var/log/app.log"},
        EnableColor:       false,    // 禁用颜色处理
        DisableCaller:     true,     // 禁用调用者查找
        DisableStacktrace: true,     // 禁用堆栈跟踪
    }
}
```

## 从配置文件加载

### YAML 配置

```yaml
# config.yaml
log:
  level: "info"
  format: "json"
  output_paths: ["stdout", "/var/log/app.log"]
  error_output_paths: ["stderr", "/var/log/error.log"]
  enable_color: false
  disable_caller: false
  disable_stacktrace: false
  stacktrace_level: "error"
  log_rotate_max_size: 100
  log_rotate_max_backups: 5
  log_rotate_max_age: 30
  log_rotate_compress: true
```

```go
import (
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

type AppConfig struct {
    Log log.Options `mapstructure:"log"`
}

func main() {
    var cfg AppConfig
    
    err := config.LoadConfig(&cfg,
        config.WithConfigFile("config.yaml", ""),
    )
    if err != nil {
        panic(err)
    }
    
    log.Init(&cfg.Log)
}
```

### JSON 配置

```json
{
  "log": {
    "level": "info",
    "format": "json",
    "output_paths": ["stdout", "/var/log/app.log"],
    "enable_color": false,
    "log_rotate_max_size": 100
  }
}
```

### 环境变量覆盖

```bash
# 环境变量可以覆盖配置文件设置
export APP_LOG_LEVEL=debug
export APP_LOG_FORMAT=text
export APP_LOG_OUTPUT_PATHS=stdout
export APP_LOG_ENABLE_COLOR=true
```

## 动态配置更新

### 使用配置热重载

```go
type AppConfig struct {
    Log log.Options `mapstructure:"log"`
}

func main() {
    var cfg AppConfig
    
    cm, err := config.LoadConfigAndWatch(&cfg,
        config.WithConfigFile("config.yaml", ""),
        config.WithHotReload(true),
    )
    if err != nil {
        panic(err)
    }
    
    // 初始化日志
    log.Init(&cfg.Log)
    
    // 注册日志配置变更回调
    cm.RegisterSectionChangeCallback("log", func(v *viper.Viper) error {
        var newLogOpts log.Options
        if err := v.UnmarshalKey("log", &newLogOpts); err != nil {
            return err
        }
        
        log.Info("更新日志配置")
        log.Init(&newLogOpts)
        log.Info("日志配置已更新")
        
        return nil
    })
    
    // 应用程序逻辑...
}
```

## 配置验证

### 验证函数

```go
func (opts *log.Options) Validate() error {
    // 验证日志级别
    validLevels := []string{"debug", "info", "warn", "error", "fatal", "panic"}
    if !contains(validLevels, opts.Level) {
        return fmt.Errorf("无效的日志级别: %s", opts.Level)
    }
    
    // 验证输出格式
    validFormats := []string{"text", "json", "keyvalue"}
    if !contains(validFormats, opts.Format) {
        return fmt.Errorf("无效的日志格式: %s", opts.Format)
    }
    
    // 验证输出路径
    if len(opts.OutputPaths) == 0 {
        return fmt.Errorf("至少需要一个输出路径")
    }
    
    // 验证轮转配置
    if opts.LogRotateMaxSize <= 0 {
        return fmt.Errorf("日志轮转最大大小必须大于 0")
    }
    
    if opts.LogRotateMaxBackups < 0 {
        return fmt.Errorf("日志轮转最大备份数不能为负数")
    }
    
    return nil
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}
```

### 使用验证

```go
func initializeLogging(opts *log.Options) error {
    // 验证配置
    if err := opts.Validate(); err != nil {
        return fmt.Errorf("日志配置验证失败: %w", err)
    }
    
    // 初始化日志
    log.Init(opts)
    
    return nil
}
```

## 最佳实践

1. **环境特定配置**：为不同环境使用不同的配置
2. **敏感信息**：避免在日志中记录敏感信息
3. **性能考虑**：在高负载环境中调整日志级别和选项
4. **轮转配置**：在生产环境中始终配置日志轮转
5. **监控**：监控日志文件大小和磁盘使用情况

## 下一步

- [输出格式](03_output_formats.md) - 了解不同的输出格式
- [上下文日志记录](04_context_logging.md) - 掌握上下文感知日志记录
- [性能优化](05_performance.md) - 优化日志性能
- [最佳实践](06_best_practices.md) - 生产就绪模式 