# 配置管理模块

`pkg/config` 模块为 Go 应用程序提供了灵活且健壮的配置管理功能，其设计借鉴了 Marmotedu 等生态系统中的最佳实践。

## 快速链接

- **[English Documentation](README.md)** - 英文文档
- **[快速开始指南](zh/01_quick_start.md)** - 几分钟内上手
- **[配置选项](zh/02_configuration_options.md)** - 所有可用选项
- **[热重载](zh/03_hot_reload.md)** - 动态配置更新
- **[最佳实践](zh/04_best_practices.md)** - 推荐模式
- **[集成示例](zh/05_integration_examples.md)** - 实际应用示例
- **[故障排除](zh/06_troubleshooting.md)** - 常见问题和解决方案
- **[模块规范](zh/07_module_specification.md)** - 完整的 API 参考

## 特性

### 🚀 高性能
- 基于 Viper 库，高效的配置管理
- 配置访问的最小开销
- 针对高频配置读取进行优化

### 📝 多配置源支持
- **文件**: YAML、JSON、TOML 等格式
- **环境变量**: 支持前缀的自动绑定
- **默认值**: 使用结构体标签设置默认值
- **命令行**: 与 flag 包集成

### 🔄 动态配置
- **热重载**: 文件变更时自动重新加载配置
- **回调系统**: 为配置变更注册回调函数
- **监控模式**: 实时配置监控
- **优雅更新**: 无中断的配置更新

### 🎯 类型安全
- **强类型**: 通过用户定义的结构体实现强类型
- **验证**: 内置验证支持
- **自动解析**: 直接映射到 Go 结构体
- **错误处理**: 全面的错误报告

### ⚙️ 易于集成
- **简单 API**: 最小化设置要求
- **框架无关**: 适用于任何 Go 应用程序
- **中间件支持**: 易于与 Web 框架集成
- **测试友好**: 轻松模拟和测试配置

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

## 安装

配置模块是 lmcc-go-sdk 的一部分：

```bash
go get github.com/lmcc-dev/lmcc-go-sdk
```

## 基础配置

### 简单配置

```go
import "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"

var cfg MyConfig
err := config.LoadConfig(&cfg)
```

### 高级配置

```go
cm, err := config.LoadConfigAndWatch(
    &cfg,
    config.WithConfigFile("config.yaml", ""),
    config.WithEnvPrefix("APP"),
    config.WithHotReload(true),
)
```

### YAML 配置

```yaml
# config.yaml
server:
  host: "localhost"
  port: 8080
  timeout: "30s"
database:
  host: "localhost"
  port: 5432
  name: "myapp"
debug: false
```

## 与其他模块集成

配置模块与其他 SDK 模块无缝集成：

```go
import (
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

type AppConfig struct {
    Log    log.Options    `mapstructure:"log"`
    Server ServerConfig   `mapstructure:"server"`
}

func main() {
    var cfg AppConfig
    
    // 加载配置
    cm, err := config.LoadConfigAndWatch(&cfg, 
        config.WithConfigFile("config.yaml", ""),
        config.WithHotReload(true),
    )
    if err != nil {
        panic(err)
    }
    
    // 使用配置初始化日志
    log.Init(&cfg.Log)
    
    // 注册热重载
    log.RegisterConfigHotReload(cm)
    
    log.Info("应用程序已启动，配置集成完成")
}
```

## 快速开始

1. **[快速开始指南](zh/01_quick_start.md)** - 基础设置和使用
2. **[配置选项](zh/02_configuration_options.md)** - 详细配置
3. **[热重载](zh/03_hot_reload.md)** - 动态更新
4. **[最佳实践](zh/04_best_practices.md)** - 生产环境建议

## 贡献

在提交 pull request 之前，请阅读我们的[贡献指南](../../../CONTRIBUTING.md)。

## 许可证

本项目采用 MIT 许可证 - 详情请参阅 [LICENSE](../../../LICENSE) 文件。