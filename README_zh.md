# LMCC Go SDK

[![Go Report Card](https://goreportcard.com/badge/github.com/lmcc-dev/lmcc-go-sdk)](https://goreportcard.com/report/github.com/lmcc-dev/lmcc-go-sdk)
[![Go Reference](https://pkg.go.dev/badge/github.com/lmcc-dev/lmcc-go-sdk.svg)](https://pkg.go.dev/github.com/lmcc-dev/lmcc-go-sdk)
<!-- 后续添加其他徽章，如构建状态、覆盖率等 -->

[**English README**](./README.md)

`lmcc-go-sdk` 是一个 Go 语言软件开发工具包，旨在为构建健壮的应用程序提供基础组件和实用工具。

## ✨ 主要特性

*   **配置管理 (`pkg/config`):** 支持从文件（YAML、TOML）、环境变量和结构体标签默认值灵活加载配置，并具备热加载能力。
*   **(更多特性待添加)**

## 🚀 快速开始

### 安装

```bash
go get github.com/lmcc-dev/lmcc-go-sdk
```

### 快速入门示例 (配置管理)

```go
package main

import (
	"flag"
	"fmt"
	"log"
	sdkconfig "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
	"time"
)

// 定义您的应用程序配置结构体
type ServerConfig struct {
	Host string        `mapstructure:"host" default:"localhost"` // 主机地址
	Port int           `mapstructure:"port" default:"8080"`    // 端口
	Timeout time.Duration `mapstructure:"timeout" default:"5s"`   // 超时时间
}

type AppConfig struct {
	sdkconfig.Config // 嵌入 SDK 基础配置 (可选但推荐)
	Server *ServerConfig `mapstructure:"server"`               // 服务配置
	Debug  bool          `mapstructure:"debug" default:"false"` // 调试模式
}

var MyConfig AppConfig

func main() {
	// 使用 flag 获取配置文件路径，默认值为 "config.yaml"
	configFile := flag.String("config", "config.yaml", "配置文件路径 (例如 config.yaml)")
	flag.Parse()

	// 加载配置
	err := sdkconfig.LoadConfig(
		&MyConfig,                                // 指向您的配置结构体
		sdkconfig.WithConfigFile(*configFile, ""), // 从文件加载 (自动推断类型)
		// sdkconfig.WithEnvPrefix("MYAPP"),      // 可选: 覆盖默认的环境变量前缀 "LMCC"
		// sdkconfig.WithHotReload(),             // 可选: 启用热加载
	)
	if err != nil {
		// 可以根据需要处理特定错误类型，例如配置文件未找到
		log.Printf("警告: 从文件 '%s' 加载配置失败，将使用默认值和环境变量: %v\n", *configFile, err)
		// 在这里决定是应该视为致命错误，还是可以继续使用默认值运行
	} else {
		log.Printf("配置从 %s 加载成功\n", *configFile)
	}

	// 访问配置值
	fmt.Printf("服务器主机: %s\n", MyConfig.Server.Host)
	fmt.Printf("服务器端口: %d\n", MyConfig.Server.Port)
	fmt.Printf("服务器超时: %s\n", MyConfig.Server.Timeout)
	fmt.Printf("调试模式: %t\n", MyConfig.Debug)

	// 示例 config.yaml 文件内容:
	/*
	server:
	  host: "127.0.0.1"
	  port: 9090
	debug: true
	*/

	// 示例环境变量 (假设使用默认前缀 LMCC):
	// export LMCC_SERVER_PORT=9999
	// export LMCC_DEBUG=true
}

```

## 📚 使用指南

有关特定模块的详细信息，请参阅 [使用指南](./docs/usage-guides/index_zh.md)。

## 🤝 贡献

欢迎贡献！请参考 `CONTRIBUTING.md` 文件（待添加）获取贡献指南。

## 📄 许可证

本项目采用 MIT 许可证授权 - 详情请参阅 `LICENSE` 文件（待添加）。 