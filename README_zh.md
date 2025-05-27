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
