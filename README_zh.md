# LMCC Go SDK

[![Go Report Card](https://goreportcard.com/badge/github.com/lmcc-dev/lmcc-go-sdk)](https://goreportcard.com/report/github.com/lmcc-dev/lmcc-go-sdk)
[![Go Reference](https://pkg.go.dev/badge/github.com/lmcc-dev/lmcc-go-sdk.svg)](https://pkg.go.dev/github.com/lmcc-dev/lmcc-go-sdk)

为构建健壮应用程序提供基础组件和实用工具的综合 Go SDK。

## 快速链接

- **[English Documentation](README.md)** - 英文文档
- **[📚 使用指南](./docs/usage-guides/)** - 全面的模块文档
- **[API 参考](https://pkg.go.dev/github.com/lmcc-dev/lmcc-go-sdk)** - Go 包文档
- **[示例](./examples/)** - 可运行的代码示例

## 特性

### 📦 核心模块
- **配置管理**: 支持热重载的多源配置
- **结构化日志**: 高性能多格式日志记录
- **错误处理**: 增强的错误处理与错误码和堆栈跟踪

### 🚀 开发体验
- **类型安全**: 通过用户定义结构体实现强类型
- **热重载**: 无需重启的动态配置更新
- **多种格式**: 支持 JSON、YAML、TOML 配置
- **环境集成**: 自动环境变量绑定

## 快速示例

```go
package main

import (
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

func main() {
	// 初始化日志
	log.Init(nil)
	log.Info("你好，LMCC Go SDK！")
	
	// 加载配置
	var cfg MyConfig
	err := config.LoadConfig(&cfg)
	if err != nil {
		log.Error("加载配置失败", "error", err)
	}
}
```

## 安装

```bash
go get github.com/lmcc-dev/lmcc-go-sdk
```

## 可用模块

| 模块 | 描述 | 文档 |
|------|------|------|
| **config** | 支持热重载的配置管理 | [📖 指南](./docs/usage-guides/config/) |
| **log** | 高性能结构化日志记录 | [📖 指南](./docs/usage-guides/log/) |
| **errors** | 增强的错误处理与错误码 | [📖 指南](./docs/usage-guides/errors/) |

## 开发工具

本项目包含一个用于开发工作流和示例管理的全面 Makefile。

### 快速命令

```bash
# 开发工作流
make help              # 显示所有可用命令
make all               # 格式化、检查、测试和整理（提交前推荐）
make format            # 格式化 Go 源代码
make lint              # 运行代码检查器
make test-unit         # 运行单元测试
make cover             # 生成覆盖率报告

# 示例管理（5 个分类中的 19 个示例）
make examples-list                        # 列出所有可用示例
make examples-run EXAMPLE=basic-usage    # 运行特定示例
make examples-test                       # 测试所有示例
make examples-build                      # 构建所有示例
make examples-debug EXAMPLE=basic-usage  # 使用 delve 调试

# 文档
make doc-serve         # 启动本地文档服务器
make doc-view PKG=./pkg/log  # 在终端查看包文档
```

### 示例分类

项目包含 **19 个实用示例**，分为 **5 个分类**：

- **basic-usage** (1): 基础集成模式
- **config-features** (5): 配置管理演示
- **error-handling** (5): 错误处理模式
- **integration** (3): 完整集成场景
- **logging-features** (5): 日志功能

**📖 完整的 Makefile 文档**: [docs/usage-guides/makefile/](./docs/usage-guides/makefile/)

## 快速开始

1. **[浏览所有模块](./docs/usage-guides/)** 在使用指南目录中
2. **选择一个模块** 符合你的需求
3. **按照该模块的快速开始指南** 进行操作
4. **使用详细文档探索高级功能**
5. **查看最佳实践** 以获得生产就绪的模式

## 贡献

欢迎贡献！请查看我们的[贡献指南](./CONTRIBUTING.md)。

## 许可证

本项目采用 MIT 许可证 - 详情请参阅 [LICENSE](./LICENSE) 文件。 
