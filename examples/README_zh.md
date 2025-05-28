/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This documentation was collaboratively developed by Martin and AI Assistant.
 */

# LMCC Go SDK 示例集合

[English Version](README.md)

本目录包含了 LMCC Go SDK 三个核心模块的综合示例：

- **Config**: 配置管理（支持热重载）
- **Errors**: 增强错误处理（支持堆栈跟踪和错误码）
- **Log**: 结构化日志（支持上下文）

## 目录结构

```
examples/
├── README.md                       # 英文版文档
├── README_zh.md                    # 本文件
├── basic-usage/                    # 三个模块的基础用法
│   ├── main.go
│   ├── config.yaml
│   └── README.md
├── config-features/                # 配置模块功能演示
│   ├── 01-simple-config/           # 简单配置加载
│   ├── 02-hot-reload/              # 热重载演示
│   ├── 03-env-override/            # 环境变量覆盖
│   ├── 04-default-values/          # 默认值演示
│   ├── 05-multiple-formats/        # 多文件格式支持
│   └── README.md
├── error-handling/                 # 错误处理示例
│   ├── 01-basic-errors/            # 基础错误创建
│   ├── 02-error-wrapping/          # 错误包装
│   ├── 03-error-codes/             # 错误码使用
│   ├── 04-stack-traces/            # 堆栈跟踪演示
│   ├── 05-error-groups/            # 错误聚合
│   └── README.md
├── logging-features/               # 日志模块示例
│   ├── 01-basic-logging/           # 基础日志
│   ├── 02-structured-logging/      # 结构化日志
│   ├── 03-log-levels/              # 日志级别控制
│   ├── 04-custom-formatters/       # 自定义格式化器
│   ├── 05-integration-patterns/    # 集成模式
│   └── README.md
├── integration/                    # 集成示例
│   ├── web-app/                    # Web应用示例
│   ├── microservice/               # 微服务示例
│   ├── cli-tool/                   # 命令行工具示例
│   └── README.md

```

## 快速开始

### 1. 基础用法
从 `basic-usage/` 示例开始，了解三个模块如何协同工作：

```bash
cd examples/basic-usage
go run main.go
```

### 2. 模块专门示例
探索各个模块的特定功能：

```bash
# 配置示例
cd examples/config-features/01-simple-config
go run main.go

# 错误处理示例
cd examples/error-handling/01-basic-errors
go run main.go

# 日志示例
cd examples/logging-features/01-basic-logging
go run main.go
```

### 3. 集成示例
查看真实世界的使用模式：

```bash
# Web应用
cd examples/integration/web-app
go run main.go

# CLI工具
cd examples/integration/cli-tool
go run main.go --help
```

## 前置要求

- Go 1.21 或更高版本
- 对 Go modules 的基本了解

## 安装

每个示例都是独立的。运行任何示例：

1. 导航到示例目录
2. 运行 `go mod tidy`（如需要）
3. 运行 `go run main.go`

## 示例概览

### 基础用法 (`basic-usage/`)
演示三个模块的基本集成：
- 带默认值的配置加载
- 正确的错误包装处理
- 带配置的结构化日志

### 配置功能 (`config-features/`)
展示各种配置能力：
- **01-simple-config**: 基础配置文件加载
- **02-hot-reload**: 实时配置更新
- **03-env-override**: 环境变量优先级
- **04-default-values**: 默认值机制
- **05-multiple-formats**: YAML、JSON、TOML 支持

### 错误处理 (`error-handling/`)
演示错误管理模式：
- **01-basic-errors**: 创建和格式化错误
- **02-error-wrapping**: 为错误添加上下文
- **03-error-codes**: 使用类型化错误码
- **04-stack-traces**: 堆栈跟踪捕获和格式化
- **05-error-groups**: 聚合多个错误

### 日志功能 (`logging-features/`)
展示日志能力：
- **01-basic-logging**: 基础日志记录和不同日志级别
- **02-structured-logging**: JSON格式和结构化字段日志
- **03-log-levels**: 日志级别控制和动态调整
- **04-custom-formatters**: 自定义日志格式化器
- **05-integration-patterns**: 日志系统集成模式

### 集成示例 (`integration/`)
真实世界应用模式：
- **web-app**: 带中间件的HTTP服务器
- **microservice**: 带可观测性的gRPC服务
- **cli-tool**: 命令行应用程序

## 最佳实践演示

1. **配置管理**
   - 使用结构体标签设置默认值
   - 为动态服务实现热重载
   - 使用环境变量提供部署灵活性

2. **错误处理**
   - 始终为错误添加上下文
   - 使用类型化错误码用于API响应
   - 保留堆栈跟踪用于调试

3. **日志记录**
   - 使用结构化日志便于机器解析
   - 在所有日志消息中包含上下文
   - 为不同环境配置适当的日志级别

## 贡献

添加新示例：

1. 在适当的类别中创建新目录
2. 包含带有适当文档的 `main.go`
3. 如需要添加配置文件
4. 编写全面的 `README.md`
5. 更新此主README.md

## 支持

关于这些示例的问题或疑问，请参考：
- [配置模块文档](../docs/usage-guides/config/)
- [错误处理文档](../docs/usage-guides/errors/)
- [日志文档](../docs/usage-guides/log/)

## 许可证

这些示例是 LMCC Go SDK 的一部分，遵循相同的许可证条款。 