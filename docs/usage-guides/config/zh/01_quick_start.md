# 快速开始指南

本指南将帮助您在几分钟内开始使用配置模块。

## 安装

配置模块是 lmcc-go-sdk 的一部分。如果您还没有安装，请将其添加到您的项目中：

```bash
go mod init your-project
go get github.com/lmcc-dev/lmcc-go-sdk
```

## 基本用法

### 步骤 1：定义配置结构

创建一个表示应用程序配置的结构体：

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
    
    Database struct {
        URL      string `mapstructure:"url" default:"postgres://localhost/mydb"`
        MaxConns int    `mapstructure:"max_connections" default:"10"`
    } `mapstructure:"database"`
    
    Debug bool `mapstructure:"debug" default:"false"`
}
```

### 步骤 2：创建配置文件

在项目根目录创建 `config.yaml` 文件：

```yaml
server:
  host: "0.0.0.0"
  port: 3000

database:
  url: "postgres://user:pass@localhost/production_db"
  max_connections: 25

debug: false
```

### 步骤 3：加载配置

```go
func main() {
    var cfg AppConfig
    
    // 简单配置加载（不带热重载）
    err := config.LoadConfig(
        &cfg,
        config.WithConfigFile("config.yaml", ""),
        config.WithEnvPrefix("APP"),
        config.WithEnvVarOverride(true),
    )
    if err != nil {
        log.Fatalf("加载配置失败: %v", err)
    }
    
    // 使用您的配置
    log.Printf("服务器将运行在 %s:%d", cfg.Server.Host, cfg.Server.Port)
    log.Printf("数据库 URL: %s", cfg.Database.URL)
    log.Printf("调试模式: %t", cfg.Debug)
}
```

## 环境变量覆盖

您可以使用环境变量覆盖任何配置值。使用前缀 `APP`，环境变量将是：

```bash
export APP_SERVER_HOST=production.example.com
export APP_SERVER_PORT=443
export APP_DEBUG=true
```

## 热重载示例

对于需要动态配置更新的应用程序：

```go
func main() {
    var cfg AppConfig
    
    // 加载配置并启用热重载支持
    cm, err := config.LoadConfigAndWatch(
        &cfg,
        config.WithConfigFile("config.yaml", ""),
        config.WithEnvPrefix("APP"),
        config.WithHotReload(true),
        config.WithEnvVarOverride(true),
    )
    if err != nil {
        log.Fatalf("加载配置失败: %v", err)
    }
    
    // 注册配置变更回调
    cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
        appCfg := currentCfg.(*AppConfig)
        log.Printf("配置已更新！新端口: %d", appCfg.Server.Port)
        
        // 在这里更新您的应用程序组件
        // 例如：重启 HTTP 服务器、重新连接数据库等
        
        return nil
    })
    
    // 您的应用程序逻辑
    log.Printf("应用程序启动，配置: %+v", cfg)
    
    // 保持应用程序运行以观察热重载
    select {}
}
```

## 配置文件格式

模块支持多种文件格式。以下是示例：

### YAML（推荐）
```yaml
server:
  host: localhost
  port: 8080
debug: true
```

### JSON
```json
{
  "server": {
    "host": "localhost",
    "port": 8080
  },
  "debug": true
}
```

### TOML
```toml
debug = true

[server]
host = "localhost"
port = 8080
```

## 默认值

默认值通过结构体标签指定，在以下情况下使用：
- 未提供配置文件
- 配置文件中缺少字段
- 未设置环境变量

```go
type Config struct {
    // 如果未指定，将默认为 "localhost"
    Host string `mapstructure:"host" default:"localhost"`
    
    // 如果未指定，将默认为 8080
    Port int `mapstructure:"port" default:"8080"`
    
    // 如果未指定，将默认为 false
    Debug bool `mapstructure:"debug" default:"false"`
}
```

## 错误处理

始终适当处理配置加载错误：

```go
err := config.LoadConfig(&cfg, options...)
if err != nil {
    // 记录带有上下文的错误
    log.Printf("配置错误: %v", err)
    
    // 决定如何处理错误：
    // 1. 使用默认配置
    // 2. 退出应用程序
    // 3. 重试加载
    
    log.Fatal("没有有效配置无法继续")
}
```

## 下一步

- [配置选项](02_configuration_options.md) - 了解所有可用选项
- [热重载](03_hot_reload.md) - 实现动态配置更新
- [最佳实践](04_best_practices.md) - 遵循推荐模式
- [集成示例](05_integration_examples.md) - 查看实际示例 