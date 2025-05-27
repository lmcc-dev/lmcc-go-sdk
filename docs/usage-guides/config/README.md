# Configuration Management Module

The `pkg/config` module provides flexible and robust configuration management capabilities for Go applications, inspired by best practices seen in ecosystems like Marmotedu.

## Quick Links

- [中文文档 (Chinese Documentation)](README_zh.md)
- [Quick Start Guide](en/01_quick_start.md)
- [Configuration Options](en/02_configuration_options.md)
- [Hot Reload](en/03_hot_reload.md)
- [Best Practices](en/04_best_practices.md)
- [Integration Examples](en/05_integration_examples.md)
- [Troubleshooting](en/06_troubleshooting.md)

## Overview

The config module leverages the Viper library for handling various configuration sources such as files (YAML, JSON, TOML, etc.), environment variables, command-line flags, and default values defined via struct tags.

### Key Features

- **Multiple Configuration Sources**: Load from files, environment variables, and defaults
- **Hot Reload**: Automatic configuration reloading when files change
- **Type Safety**: Strong typing through user-defined structs
- **Callback System**: Register callbacks for configuration changes
- **Environment Variable Binding**: Automatic binding with prefix support
- **Default Values**: Set defaults using struct tags

## Quick Example

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
    
    // Register callback for changes
    cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
        log.Println("Configuration updated!")
        return nil
    })
    
    // Use configuration
    log.Printf("Server running on %s:%d", cfg.Server.Host, cfg.Server.Port)
}
```

## Documentation Structure

### English Documentation
- [Overview](en/00_overview.md) - Module overview and architecture
- [Quick Start](en/01_quick_start.md) - Get started quickly
- [Configuration Options](en/02_configuration_options.md) - Available options and settings
- [Hot Reload](en/03_hot_reload.md) - Dynamic configuration updates
- [Best Practices](en/04_best_practices.md) - Recommended patterns
- [Integration Examples](en/05_integration_examples.md) - Real-world examples
- [Troubleshooting](en/06_troubleshooting.md) - Common issues and solutions

### Chinese Documentation (中文文档)
- [概述](zh/00_overview_zh.md) - 模块概述和架构
- [快速开始](zh/01_quick_start_zh.md) - 快速入门
- [配置选项](zh/02_configuration_options_zh.md) - 可用选项和设置
- [热重载](zh/03_hot_reload_zh.md) - 动态配置更新
- [最佳实践](zh/04_best_practices_zh.md) - 推荐模式
- [集成示例](zh/05_integration_examples_zh.md) - 实际应用示例
- [故障排除](zh/06_troubleshooting_zh.md) - 常见问题和解决方案

## Contributing

Please read our [contributing guidelines](../../../CONTRIBUTING.md) before submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](../../../LICENSE) file for details.