# 配置管理模块

`pkg/config` 模块为 Go 应用程序提供了灵活且健壮的配置管理功能，其设计借鉴了 Marmotedu 等生态系统中的最佳实践。

## 快速链接

- [English Documentation](README.md)
- [快速开始指南](zh/01_quick_start_zh.md)
- [配置选项](zh/02_configuration_options_zh.md)
- [热重载](zh/03_hot_reload_zh.md)
- [最佳实践](zh/04_best_practices_zh.md)
- [集成示例](zh/05_integration_examples_zh.md)
- [故障排除](zh/06_troubleshooting_zh.md)

## 概述

配置模块利用 Viper 库来处理各种配置源，例如文件（YAML、JSON、TOML 等）、环境变量、命令行标志以及通过结构体标签定义的默认值。

### 主要功能

- **多配置源支持**：从文件、环境变量和默认值加载配置
- **热重载**：文件变更时自动重新加载配置
- **类型安全**：通过用户定义的结构体实现强类型
- **回调系统**：为配置变更注册回调函数
- **环境变量绑定**：支持前缀的自动绑定
- **默认值**：使用结构体标签设置默认值

## 快速示例

```go
package main

import (
    "log"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
)

type AppConfig struct {
    Server struct {
        Host string `mapstructure:"host" default:"localhost"`
        Port int    `mapstructure:"port" default:"8080"`
    } `mapstructure:"server"`
    Debug bool `mapstructure:"debug" default:"false"`
}

func main() {
    var cfg AppConfig
    
    cm, err := config.LoadConfigAndWatch(
        &cfg,
        config.WithConfigFile("config.yaml", ""),
        config.WithEnvPrefix("APP"),
        config.WithHotReload(true),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // 注册变更回调
    cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
        log.Println("配置已更新！")
        return nil
    })
    
    // 使用配置
    log.Printf("服务器运行在 %s:%d", cfg.Server.Host, cfg.Server.Port)
}
```

## 文档结构

### 英文文档
- [概述](en/00_overview.md) - 模块概述和架构
- [快速开始](en/01_quick_start.md) - 快速入门
- [配置选项](en/02_configuration_options.md) - 可用选项和设置
- [热重载](en/03_hot_reload.md) - 动态配置更新
- [最佳实践](en/04_best_practices.md) - 推荐模式
- [集成示例](en/05_integration_examples.md) - 实际应用示例
- [故障排除](en/06_troubleshooting.md) - 常见问题和解决方案

### 中文文档
- [概述](zh/00_overview_zh.md) - 模块概述和架构
- [快速开始](zh/01_quick_start_zh.md) - 快速入门
- [配置选项](zh/02_configuration_options_zh.md) - 可用选项和设置
- [热重载](zh/03_hot_reload_zh.md) - 动态配置更新
- [最佳实践](zh/04_best_practices_zh.md) - 推荐模式
- [集成示例](zh/05_integration_examples_zh.md) - 实际应用示例
- [故障排除](zh/06_troubleshooting_zh.md) - 常见问题和解决方案

## 贡献

在提交 pull request 之前，请阅读我们的[贡献指南](../../../CONTRIBUTING.md)。

## 许可证

本项目采用 MIT 许可证 - 详情请参阅 [LICENSE](../../../LICENSE) 文件。