# 日志模块文档

日志模块为 Go 应用程序提供了强大、灵活且高性能的日志记录解决方案。基于 Zap 构建，提供结构化日志记录、多种输出格式、日志轮转和上下文感知日志功能。

## 快速链接

- **[English Documentation](README.md)** - 英文文档
- **[快速开始指南](zh/01_quick_start.md)** - 几分钟内上手
- **[配置选项](zh/02_configuration_options.md)** - 所有可用选项
- **[输出格式](zh/03_output_formats.md)** - JSON、文本和键值对格式
- **[上下文日志](zh/04_context_logging.md)** - 上下文感知日志
- **[日志轮转](zh/05_log_rotation.md)** - 文件轮转和管理
- **[最佳实践](zh/06_best_practices.md)** - 推荐模式
- **[集成示例](zh/07_integration_examples.md)** - 实际应用示例
- **[故障排除](zh/08_troubleshooting.md)** - 常见问题和解决方案

## 特性

### 🚀 高性能
- 基于 Uber 的 Zap 日志库，性能卓越
- 热路径零内存分配
- 高效的结构化日志记录，开销最小

### 📝 多种输出格式
- **JSON**: 机器可读的结构化日志
- **文本**: 人类可读的控制台输出，支持颜色
- **键值对**: 简单的键值对格式

### 🔄 灵活的输出目标
- 控制台输出 (stdout/stderr)
- 文件输出，支持自动轮转
- 多个同时输出目标
- 自定义输出目标

### 🎯 上下文感知日志
- 请求 ID 跟踪
- 用户上下文保持
- 自动字段继承
- 结构化上下文传播

### ⚙️ 简单配置
- YAML/JSON 配置文件
- 环境变量覆盖
- 热重载支持（配合 config 模块）
- 合理的默认值，快速设置

### 🔧 高级功能
- 按输出目标过滤日志级别
- 调用者信息（文件、行号、函数）
- 错误堆栈跟踪
- 高流量场景的日志采样
- 自定义字段编码器

## 快速示例

```go
package main

import (
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

func main() {
    // 使用默认设置初始化
    log.Init(nil)
    
    // 基础日志记录
    log.Info("应用程序已启动")
    log.Warn("这是一个警告")
    log.Error("出现了错误")
    
    // 结构化日志记录
    log.Infow("用户登录",
        "user_id", 12345,
        "username", "john_doe",
        "ip", "192.168.1.100",
    )
    
    // 上下文日志记录
    ctx := log.WithContext(context.Background())
    ctx = log.WithValues(ctx, "request_id", "req-123")
    
    log.InfoContext(ctx, "正在处理请求")
}
```

## 安装

日志模块是 lmcc-go-sdk 的一部分：

```bash
go get github.com/lmcc-dev/lmcc-go-sdk
```

## 基础配置

### 简单配置

```go
import "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"

// 使用默认值（控制台输出，info 级别）
log.Init(nil)
```

### 自定义配置

```go
opts := &log.Options{
    Level:            "debug",
    Format:           "json",
    OutputPaths:      []string{"stdout", "/var/log/app.log"},
    EnableColor:      false,
    DisableCaller:    false,
    DisableStacktrace: false,
}

log.Init(opts)
```

### YAML 配置

```yaml
# config.yaml
log:
  level: "info"
  format: "json"
  output_paths: ["stdout", "/var/log/app.log"]
  enable_color: true
  disable_caller: false
  disable_stacktrace: false
  stacktrace_level: "error"
  
  # 日志轮转设置
  log_rotate_max_size: 100      # MB
  log_rotate_max_backups: 5     # 文件数
  log_rotate_max_age: 30        # 天数
  log_rotate_compress: true
```

## 输出格式示例

### JSON 格式
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

### 文本格式
```
2024-01-15T10:30:45.123Z	INFO	main.go:25	用户登录	{"user_id": 12345, "username": "john_doe", "ip": "192.168.1.100"}
```

### 键值对格式
```
timestamp=2024-01-15T10:30:45.123Z level=info caller=main.go:25 message="用户登录" user_id=12345 username=john_doe ip=192.168.1.100
```

## 与配置模块集成

日志模块与配置模块无缝集成，支持热重载：

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
    
    // 加载配置并启用热重载
    cm, err := config.LoadConfigAndWatch(&cfg, 
        config.WithConfigFile("config.yaml", ""),
        config.WithHotReload(true),
    )
    if err != nil {
        panic(err)
    }
    
    // 初始化日志
    log.Init(&cfg.Log)
    
    // 注册热重载
    log.RegisterConfigHotReload(cm)
    
    log.Info("应用程序已启动，支持日志热重载")
}
```

## 性能特征

日志模块专为高性能应用程序设计：

- **零分配** 对于禁用的日志级别
- **最小分配** 对于启用的日志
- **高效 JSON 编码** 使用 Zap 的优化编码器
- **缓冲 I/O** 用于文件输出
- **异步日志** 选项，实现最大吞吐量

### 基准测试

```
BenchmarkLogInfo-8           	 5000000	       230 ns/op	       0 B/op	       0 allocs/op
BenchmarkLogInfow-8          	 3000000	       450 ns/op	      64 B/op	       1 allocs/op
BenchmarkLogJSON-8           	 2000000	       680 ns/op	     128 B/op	       2 allocs/op
```

## 架构

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│     应用程序     │───▶│     日志模块     │───▶│   Zap 日志器    │
│                 │    │                  │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                │                        │
                                ▼                        ▼
                       ┌──────────────────┐    ┌─────────────────┐
                       │    配置选项      │    │   输出接收器    │
                       │                  │    │ (控制台/文件)   │
                       └──────────────────┘    └─────────────────┘
```

## 使用场景

### Web 应用程序
- 带有关联 ID 的请求/响应日志
- 错误跟踪和调试
- 性能监控
- 安全审计跟踪

### 微服务
- 分布式跟踪关联
- 服务间通信日志
- 健康检查和指标日志
- 配置变更跟踪

### CLI 应用程序
- 用户操作日志
- 错误报告
- 调试信息
- 进度跟踪

### 后台服务
- 作业处理日志
- 定时任务执行
- 系统监控
- 数据处理管道

## 开始使用

1. **[快速开始指南](zh/01_quick_start.md)** - 基础设置和使用
2. **[配置选项](zh/02_configuration_options.md)** - 详细配置
3. **[输出格式](zh/03_output_formats.md)** - 选择合适的格式
4. **[上下文日志](zh/04_context_logging.md)** - 高级上下文功能
5. **[最佳实践](zh/06_best_practices.md)** - 生产环境建议

## 贡献

我们欢迎贡献！请查看我们的[贡献指南](../../CONTRIBUTING.md)了解详情。

## 许可证

本项目采用 MIT 许可证 - 详情请查看 [LICENSE](../../LICENSE) 文件。 