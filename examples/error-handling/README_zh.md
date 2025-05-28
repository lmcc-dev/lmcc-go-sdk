/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This documentation was collaboratively developed by Martin and AI Assistant.
 */

# 错误处理示例

[English Version](README.md)

本目录包含展示错误处理模块各种功能的示例。

## 示例概览

### [01-basic-errors](01-basic-errors/)
**基础错误创建和格式化**
- 创建简单错误
- 使用 %v, %s, %+v 格式化错误
- 错误消息组合
- 堆栈跟踪基础

**学习内容**: 如何创建和格式化基础错误

### [02-error-wrapping](02-error-wrapping/)
**错误包装和上下文**
- 用额外上下文包装错误
- 错误链导航
- 原因提取
- 上下文保存

**学习内容**: 如何在保留错误信息的同时添加上下文

### [03-error-codes](03-error-codes/)
**错误码和类型**
- 使用预定义错误码
- 创建自定义错误码
- 错误码分类
- HTTP状态码映射

**学习内容**: 如何实现带有错误码的结构化错误处理

### [04-stack-traces](04-stack-traces/)
**堆栈跟踪捕获和分析**
- 自动堆栈跟踪捕获
- 堆栈跟踪格式化
- 性能考虑
- 用堆栈跟踪调试

**学习内容**: 如何利用堆栈跟踪进行调试

### [05-error-groups](05-error-groups/)
**错误聚合和分组**
- 收集多个错误
- 错误组操作
- 并行处理错误
- 错误过滤和分类

**学习内容**: 如何在复杂操作中处理多个错误

## 运行示例

每个示例都是独立的，可以单独运行：

```bash
# 导航到特定示例
cd examples/error-handling/01-basic-errors

# 运行示例
go run main.go

# 某些示例支持额外的参数
go run main.go --help
```

## 演示的常用模式

### 1. 基础错误创建
```go
// 简单错误
err := errors.New("operation failed")

// 格式化错误
err := errors.Errorf("invalid value: %d", value)

// 带错误码的错误
err := errors.WithCode(err, errors.ErrValidation)
```

### 2. 错误包装
```go
// 用上下文包装
err = errors.Wrap(err, "failed to process user data")

// 用格式化消息包装
err = errors.Wrapf(err, "user %s: operation failed", userID)

// 提取原始错误
originalErr := errors.Cause(err)
```

### 3. 错误码使用
```go
// 检查错误码
if coder := errors.GetCoder(err); coder != nil {
    switch coder.Code() {
    case errors.ErrValidation.Code():
        // 处理验证错误
    case errors.ErrNotFound.Code():
        // 处理未找到错误
    }
}
```

### 4. 堆栈跟踪处理
```go
// 打印详细堆栈跟踪
fmt.Printf("Error: %+v\n", err)

// 程序化获取堆栈跟踪
if tracer := errors.GetStackTracer(err); tracer != nil {
    stack := tracer.StackTrace()
    // 处理堆栈帧
}
```

## 展示的最佳实践

1. **错误创建**
   - 使用描述性错误消息
   - 包含相关上下文信息
   - 选择适当的错误码
   - 保留错误链

2. **错误包装**
   - 在每个抽象层添加上下文
   - 维护原始错误信息
   - 使用一致的包装模式
   - 避免过度包装

3. **错误处理**
   - 检查特定错误类型/代码
   - 在适当级别处理错误
   - 记录带有足够上下文的错误
   - 提供有意义的用户反馈

4. **堆栈跟踪**
   - 使用堆栈跟踪进行调试
   - 注意性能影响
   - 在需要时过滤堆栈帧
   - 在错误日志中包含堆栈跟踪

5. **错误组**
   - 将相关错误收集在一起
   - 提供错误摘要
   - 优雅地处理部分失败
   - 使用并行错误处理

## 集成技巧

### 与日志模块
```go
// 带上下文记录错误
logger.Errorf("Operation failed: %+v", err)

// 带错误码记录
if coder := errors.GetCoder(err); coder != nil {
    logger.WithFields(log.Fields{
        "error_code": coder.Code(),
        "error_type": coder.String(),
    }).Error("Structured error logging")
}
```

### 与配置模块
```go
// 特定处理配置错误
if err := config.LoadConfig(&cfg); err != nil {
    if coder := errors.GetCoder(err); coder != nil {
        switch coder.Code() {
        case errors.ErrConfigFileRead.Code():
            log.Fatal("Configuration file not found or unreadable")
        case errors.ErrConfigValidation.Code():
            log.Fatal("Configuration validation failed")
        }
    }
}
```

## 下一步

探索这些错误处理示例后：

1. 尝试[日志功能示例](../logging-features/)进行全面的日志记录
2. 探索[集成示例](../integration/)了解真实世界的使用模式
3. 查看[配置功能示例](../config-features/)了解配置错误处理

## 故障排除

### 常见问题

**问题**: 堆栈跟踪不显示
```
解决方案: 确保错误是用 errors.New() 或 errors.Errorf() 创建的
```

**问题**: 错误码不被识别
```
解决方案: 检查错误是否实现了 Coder 接口
```

**问题**: 包装的错误不能正确解包
```
解决方案: 使用 errors.Cause() 提取根错误
```

**问题**: 堆栈跟踪的性能影响
```
解决方案: 如果需要，考虑在生产环境中禁用堆栈跟踪
```

## 相关文档

- [错误处理模块概览](../../docs/usage-guides/errors/zh/00_overview.md)
- [快速开始指南](../../docs/usage-guides/errors/zh/01_quick_start.md)
- [错误码参考](../../docs/usage-guides/errors/zh/02_error_codes.md)
- [堆栈跟踪指南](../../docs/usage-guides/errors/zh/03_stack_traces.md) 