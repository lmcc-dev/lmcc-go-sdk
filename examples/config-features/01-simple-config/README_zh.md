/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This documentation was collaboratively developed by Martin and AI Assistant.
 */

# 简单配置加载示例

[English Version](README.md)

本示例演示配置模块最基本的用法，重点关注：
- 从YAML文件加载配置
- 使用结构体标签设置默认值
- 基础配置验证
- 配置问题的错误处理

## 本示例展示的功能

### 1. 基础配置加载
- 简单的基于结构体的配置定义
- YAML文件解析并映射到Go结构体
- 嵌套配置结构

### 2. 默认值机制
- 使用`default`结构体标签设置备用值
- 处理不同数据类型（字符串、整数、布尔值、时长、切片）
- 嵌套结构体和指针的默认行为

### 3. 配置验证
- 配置值的自定义验证逻辑
- 端口范围验证
- 超时和连接数的正值检查
- 基于功能开关的条件验证

### 4. 错误处理集成
- 配置加载错误处理
- 来自错误模块的错误码和类型
- 优雅的错误报告和调试信息

## 文件结构

```
01-simple-config/
├── main.go         # 演示配置加载的主应用程序
├── config.yaml     # 示例配置文件
├── README.md       # 英文版文档
└── README_zh.md    # 本文件
```

## 配置结构

### 应用设置
```yaml
app_name: "Advanced Simple Config Example"
version: "2.0.0"
debug: false
port: 9080
timeout: "45s"
interval: "10m"
```

### 网络配置
```yaml
allowed_ips:
  - "127.0.0.1"
  - "::1"
  - "192.168.1.0/24"
```

### 功能开关
```yaml
features:
  - "authentication"
  - "authorization"
  - "logging"
  - "metrics"
```

### 数据库设置
```yaml
database:
  driver: "postgresql"
  host: "db.example.com"
  port: 5432
  max_connections: 20
  enable_ssl: true
```

### 缓存配置
```yaml
cache:
  enabled: true
  type: "redis"
  host: "cache.example.com"
  ttl: "2h"
```

## 运行示例

### 前置要求
- Go 1.21 或更高版本
- 无需外部依赖

### 基础运行
```bash
cd examples/config-features/01-simple-config
go run main.go
```

### 测试缺少配置文件的情况
```bash
# 重命名配置文件以查看默认行为
mv config.yaml config.yaml.backup
go run main.go
# 恢复配置文件
mv config.yaml.backup config.yaml
```

### 测试无效配置
```bash
# 创建无效配置用于测试错误处理
cat > invalid-config.yaml << EOF
port: 99999  # 无效端口号
timeout: "-5s"  # 无效超时时间
EOF

# 这将显示验证错误
# go run main.go -config invalid-config.yaml
```

## 预期输出

示例会产生结构化输出，显示：

### 1. 默认值演示
```
=== Demonstrating Default Values ===
Loading configuration with defaults only (no config file)...
Expected error (file not found): ...
Zero values (before applying defaults):
  AppName: ''
  Port: 0
  Debug: false
  Features: [] (length: 0)
```

### 2. 成功配置加载
```
=== Loading Configuration from File ===
✓ Configuration loaded successfully!

=== Validating Configuration ===
✓ Configuration validation passed!
```

### 3. 配置摘要
```
=== Configuration Summary ===
Application:
  Name: Advanced Simple Config Example
  Version: 2.0.0
  Debug: false
  Port: 9080
  Timeout: 45s
  Features: [authentication authorization logging metrics tracing caching]

Database:
  Driver: postgresql
  Host: db.example.com
  Max Connections: 20
  Enable SSL: true

Cache:
  Enabled: true
  Type: redis
  TTL: 2h
```

### 4. 配置使用
```
=== Using Configuration ===
Starting Advanced Simple Config Example version 2.0.0...
Server will listen on port 9080
Database: postgresql://db.example.com:5432/production_db
Cache: redis://cache.example.com:6379 (TTL: 2h)
```

## 关键学习要点

### 1. 结构体标签默认值
```go
type Config struct {
    Port    int           `mapstructure:"port" default:"8080"`
    Timeout time.Duration `mapstructure:"timeout" default:"30s"`
    Debug   bool          `mapstructure:"debug" default:"true"`
    Features []string     `mapstructure:"features" default:"auth,logging"`
}
```

### 2. 配置加载模式
```go
var cfg SimpleAppConfig
err := config.LoadConfig(&cfg, 
    config.WithConfigFile("config.yaml", "yaml"))
if err != nil {
    // 处理错误并提供详细信息
    if coder := errors.GetCoder(err); coder != nil {
        fmt.Printf("Config Error [%d]: %s", coder.Code(), coder.String())
    }
}
```

### 3. 验证最佳实践
```go
func validateConfig(cfg *SimpleAppConfig) error {
    if cfg.Port < 1 || cfg.Port > 65535 {
        return errors.Errorf("invalid port: %d", cfg.Port)
    }
    // 更多验证逻辑...
}
```

### 4. 安全配置访问
```go
// 访问嵌套配置前检查nil指针
if cfg.Database != nil {
    connectionString := fmt.Sprintf("%s://%s:%d/%s", 
        cfg.Database.Driver, cfg.Database.Host, 
        cfg.Database.Port, cfg.Database.Database)
}
```

## 演示的常用模式

1. **嵌入基础配置**: 虽然在这个简单示例中未使用，但展示了如何为大型应用程序构建配置结构

2. **类型安全的默认值**: 使用结构体标签进行编译时默认定义

3. **嵌套配置**: 将相关设置组织成逻辑组

4. **配置验证**: 独立于加载过程实现业务规则验证

5. **错误处理**: 正确的错误传播和用户友好的错误消息

## 故障排除

### 问题：找不到配置文件
**原因**: `config.yaml`在工作目录中不存在
**解决方案**: 确保文件存在或检查文件路径

### 问题：YAML解析错误
**原因**: 配置文件中存在无效的YAML语法
**解决方案**: 使用在线验证器或`yamllint`验证YAML语法

### 问题：类型转换错误
**原因**: 配置值与期望的Go类型不匹配
**解决方案**: 检查YAML中的数据类型是否与结构体字段类型匹配

### 问题：默认值未应用
**原因**: 结构体标签缺失或语法不正确
**解决方案**: 验证`default:"value"`标签语法和字段类型

## 下一步

理解这个简单示例后：

1. 探索[02-hot-reload](../02-hot-reload/)了解动态配置更新
2. 尝试[03-env-override](../03-env-override/)学习基于环境的配置
3. 查看[04-default-values](../04-default-values/)了解高级默认值处理
4. 检查[05-multiple-formats](../05-multiple-formats/)了解不同文件格式

## 相关文档

- [配置模块概览](../../../docs/usage-guides/config/zh/)
- [配置加载指南](../../../docs/usage-guides/config/zh/01_quick_start.md)
- [默认值文档](../../../docs/usage-guides/config/zh/02_configuration_options.md) 