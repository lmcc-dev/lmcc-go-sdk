/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This documentation was collaboratively developed by Martin and AI Assistant.
 */

# 日志功能示例

[English Version](README.md)

本目录包含了展示 LMCC Go SDK 日志功能的综合示例。每个示例专注于日志系统的不同方面，从基础用法到高级集成模式。

## 示例概览

### 01-basic-logging
**基础日志操作**

演示基础日志操作，包括：
- 不同的日志级别（Debug、Info、Warn、Error）
- 带键值对的结构化日志
- 上下文感知日志
- 性能日志和计时
- 错误日志模式

```bash
cd 01-basic-logging
go run main.go
```

**主要特性：**
- 带认证流程的用户服务
- 批量操作日志
- 上下文传播
- 性能测量

### 02-structured-logging
**JSON格式和结构化字段**

展示如何实现机器可读的结构化日志：
- JSON输出格式
- 复杂数据结构日志
- HTTP请求/响应日志
- 业务事件跟踪
- 性能指标日志

```bash
cd 02-structured-logging
go run main.go
```

**主要特性：**
- RequestInfo和ResponseInfo结构
- 数据库操作日志
- 业务事件模式
- 大对象序列化

### 03-log-levels
**日志级别控制和动态调整**

演示日志级别管理和过滤：
- 动态日志级别调整
- 组件特定的日志级别
- 性能影响分析
- 生产环境场景
- 基于阈值的条件日志

```bash
cd 03-log-levels
go run main.go
```

**主要特性：**
- 级别过滤演示
- 性能基准测试
- 组件隔离
- 高级级别使用模式

### 04-custom-formatters
**自定义日志格式化器和日志轮转**

探索不同的日志输出格式、自定义和文件轮转：
- 带颜色的文本格式
- 结构化数据的JSON格式
- 便于解析的键值格式
- 格式间的性能比较
- 生产环境格式化
- **文件轮转策略（基于大小、时间、混合）**
- **不同环境的日志轮转最佳实践**

```bash
cd 04-custom-formatters
go run main.go
```

**主要特性：**
- 格式比较和基准测试
- 环境特定配置
- 自定义字段格式化
- 大数据对象处理
- **基于大小、时间和混合的日志轮转**
- **轮转配置最佳实践**
- **磁盘空间管理演示**

### 05-integration-patterns
**日志系统集成模式**

展示真实世界集成模式的综合示例：
- HTTP中间件日志
- 配置驱动的日志设置
- 服务层集成
- 错误处理集成
- 横切关注点

```bash
cd 05-integration-patterns
go run main.go
```

**主要特性：**
- HTTP请求的LoggingMiddleware
- 配置集成
- 多层错误处理
- 跨服务的上下文传播
- 真实世界服务模式

## 演示的通用模式

### 1. 上下文感知日志
所有示例展示如何通过日志系统传播上下文信息（如请求ID、用户ID）：

```go
logger := log.Std().WithValues("request_id", requestID, "user_id", userID)
logger.Infow("Processing request", "operation", "user_creation")
```

### 2. 结构化字段
示例演示一致使用结构化字段以便更好的日志解析：

```go
logger.Infow("Database operation completed",
    "operation", "INSERT",
    "table", "users",
    "duration", duration,
    "rows_affected", rowsAffected)
```

### 3. 错误集成
展示与lmcc-go-sdk错误包的集成：

```go
if err != nil {
    wrappedErr := errors.Wrap(err, "operation failed")
    logger.Errorw("Database error", "error", wrappedErr)
    return wrappedErr
}
```

### 4. 性能日志
演示计时和性能测量模式：

```go
start := time.Now()
defer func() {
    duration := time.Since(start)
    logger.Infow("Operation completed", "duration", duration)
}()
```

## 配置示例

每个示例展示不同的配置方法：

### 开发环境配置
```go
opts := log.NewOptions()
opts.Level = "debug"
opts.Format = "text"
opts.EnableColor = true
opts.DisableCaller = false
```

### 生产环境配置
```go
opts := log.NewOptions()
opts.Level = "info"
opts.Format = "json"
opts.EnableColor = false
opts.DisableCaller = false
opts.DisableStacktrace = true
```

### 与配置包集成
```go
type AppConfig struct {
    Logging struct {
        Level  string `yaml:"level" default:"info"`
        Format string `yaml:"format" default:"json"`
        Output string `yaml:"output" default:"stdout"`
    } `yaml:"logging"`
}
```

## 展示的最佳实践

1. **使用结构化日志**：始终优先使用键值对而非字符串格式化
2. **包含上下文**：添加请求ID、用户ID和其他上下文信息
3. **测量性能**：记录操作持续时间用于监控
4. **正确处理错误**：使用错误包进行错误包装和上下文
5. **为环境配置**：开发环境与生产环境使用不同设置
6. **使用适当级别**：Debug用于开发，Info用于正常操作，Warn用于关注情况，Error用于实际问题

## 运行所有示例

按顺序运行所有示例：

```bash
# 运行每个示例
for dir in 01-basic-logging 02-structured-logging 03-log-levels 04-custom-formatters 05-integration-patterns; do
    echo "=== 运行 $dir ==="
    cd $dir
    go run main.go
    cd ..
    echo
done
```

## 与其他模块的集成

这些示例演示与以下模块的集成：
- **config包**：用于配置驱动的日志设置
- **errors包**：用于增强错误处理和日志
- **HTTP中间件**：用于请求/响应日志
- **数据库操作**：用于查询和事务日志

## 下一步

探索这些示例后，考虑：
1. 在您的应用程序中实现类似模式
2. 为您的特定需求自定义日志格式
3. 设置日志聚合和监控
4. 与可观测性平台集成
5. 创建您自己的中间件和服务模式

更多信息，请参阅：
- [日志包文档](../../docs/usage-guides/log/)
- [配置包文档](../../docs/usage-guides/config/)
- [错误包文档](../../docs/usage-guides/errors/) 