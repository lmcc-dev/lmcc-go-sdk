/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This documentation was collaboratively developed by Martin and AI Assistant.
 */

# 配置功能示例

[English Version](README.md)

本目录包含展示配置模块各种功能的示例。

## 示例概览

### [01-simple-config](01-simple-config/)
**基础配置加载**
- 简单的YAML配置文件加载
- 结构体标签默认值
- 基础验证
- 错误处理

**学习内容**: 如何开始使用配置加载

### [02-hot-reload](02-hot-reload/)
**热重载配置**
- 实时配置更新
- 文件变化监控
- 回调注册
- 优雅更新

**学习内容**: 如何实现动态配置变更

### [03-env-override](03-env-override/)
**环境变量覆盖**
- 环境变量优先级
- 前缀配置
- 类型转换
- 生产部署模式

**学习内容**: 如何使用环境变量提供部署灵活性

### [04-default-values](04-default-values/)
**默认值演示**
- 结构体标签默认值
- 嵌套结构体默认值
- 指针字段处理
- 零值与显式值

**学习内容**: 如何实现健壮的默认值系统

### [05-multiple-formats](05-multiple-formats/)
**多文件格式支持**
- YAML配置
- JSON配置
- TOML配置
- 格式自动检测

**学习内容**: 如何支持不同的配置文件格式

## 运行示例

每个示例都是独立的，可以单独运行：

```bash
# 导航到特定示例
cd examples/config-features/01-simple-config

# 运行示例
go run main.go

# 某些示例支持额外的参数
go run main.go --help
```

## 演示的常用模式

### 1. 配置结构设计
```go
type AppConfig struct {
    config.Config                    // 嵌入SDK基础配置
    App    *AppSpecificConfig       `mapstructure:"app"`
    Feature *FeatureConfig         `mapstructure:"feature"`
}
```

### 2. 带选项的加载
```go
err := config.LoadConfig(&cfg,
    config.WithConfigFile("config.yaml", "yaml"),
    config.WithEnvPrefix("MYAPP"),
    config.WithEnvVarOverride(true),
)
```

### 3. 热重载设置
```go
cm, err := config.LoadConfigAndWatch(&cfg,
    config.WithConfigFile("config.yaml", "yaml"),
    config.WithHotReload(true),
)

cm.RegisterCallback(func(v *viper.Viper, cfg any) error {
    // 处理配置变更
    return nil
})
```

### 4. 环境变量模式
```bash
# 格式: {前缀}_{段}_{字段}
export MYAPP_SERVER_PORT=8080
export MYAPP_LOG_LEVEL=debug
export MYAPP_DATABASE_HOST=prod-db.example.com
```

## 展示的最佳实践

1. **结构组织**
   - 相关设置的逻辑分组
   - 基础配置和应用特定配置的清晰分离
   - 一致的命名约定

2. **默认值策略**
   - 开发环境的合理默认值
   - 生产就绪的配置
   - 优雅的备用机制

3. **环境变量使用**
   - 一致的前缀使用
   - 清晰的环境特定覆盖
   - 敏感数据的安全处理

4. **错误处理**
   - 全面的错误检查
   - 有意义的错误消息
   - 无效配置的优雅降级

5. **热重载考虑**
   - 线程安全的配置更新
   - 应用更改前的验证
   - 无效配置的回滚机制

## 集成技巧

### 与日志模块
```go
// 将配置转换为日志选项
logOpts := &log.Options{
    Level:  cfg.Log.Level,
    Format: cfg.Log.Format,
    // ... 其他映射
}
log.Init(logOpts)
```

### 与错误处理
```go
if err := config.LoadConfig(&cfg, opts...); err != nil {
    if coder := errors.GetCoder(err); coder != nil {
        log.Errorf("Config error [%d]: %s", coder.Code(), coder.String())
    }
    return errors.Wrap(err, "failed to initialize application")
}
```

## 下一步

探索这些配置示例后：

1. 尝试[错误处理示例](../error-handling/)学习健壮的错误管理
2. 探索[日志功能示例](../logging-features/)获得全面的日志记录
3. 查看[集成示例](../integration/)了解真实世界的使用模式

## 故障排除

### 常见问题

**问题**: 找不到配置文件
```
解决方案: 检查文件路径和工作目录
```

**问题**: 环境变量不生效
```
解决方案: 验证前缀格式并重启应用程序
```

**问题**: 热重载不工作
```
解决方案: 确保文件权限并检查回调注册
```

**问题**: 默认值未应用
```
解决方案: 验证结构体标签和字段类型
```

## 相关文档

- [配置模块概览](../../docs/usage-guides/config/zh/00_overview.md)
- [快速开始指南](../../docs/usage-guides/config/zh/01_quick_start.md)
- [配置选项](../../docs/usage-guides/config/zh/02_configuration_options.md)
- [热重载指南](../../docs/usage-guides/config/zh/03_hot_reload.md) 