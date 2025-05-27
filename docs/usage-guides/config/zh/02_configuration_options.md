# 配置选项

本文档描述了所有可用的配置选项以及如何有效使用它们。

## 函数选项

配置模块使用函数选项模式来实现清洁和可扩展的 API 设计。

### WithConfigFile

指定配置文件路径和可选的搜索路径。

```go
config.WithConfigFile("config.yaml", "")
config.WithConfigFile("app.yaml", "/etc/myapp:/home/user/.config")
```

**参数：**
- `filename`: 配置文件名
- `searchPaths`: 冒号分隔的目录搜索列表（可选）

**支持的文件扩展名：**
- `.yaml`, `.yml`
- `.json`
- `.toml`
- `.hcl`
- `.ini`
- `.properties`

### WithEnvPrefix

设置环境变量绑定的前缀。

```go
config.WithEnvPrefix("APP")
```

使用前缀 "APP"，以下映射适用：
- `server.host` → `APP_SERVER_HOST`
- `database.url` → `APP_DATABASE_URL`
- `debug` → `APP_DEBUG`

### WithEnvVarOverride

启用环境变量覆盖配置文件值。

```go
config.WithEnvVarOverride(true)  // 启用覆盖
config.WithEnvVarOverride(false) // 禁用覆盖
```

### WithHotReload

启用文件更改时的自动配置重新加载。

```go
config.WithHotReload(true)  // 启用热重载
config.WithHotReload(false) // 禁用热重载
```

**注意：** 仅在 `LoadConfigAndWatch()` 函数中可用。

### WithEnvKeyReplacer

自定义配置键如何转换为环境变量名称。

```go
import "strings"

config.WithEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
```

这会在环境变量名称中将点和连字符替换为下划线。

## 配置结构标签

### mapstructure 标签

将配置键映射到结构体字段。

```go
type Config struct {
    ServerPort int `mapstructure:"server_port"`
    DebugMode  bool `mapstructure:"debug"`
}
```

### default 标签

为字段指定默认值。

```go
type Config struct {
    Host    string `mapstructure:"host" default:"localhost"`
    Port    int    `mapstructure:"port" default:"8080"`
    Timeout string `mapstructure:"timeout" default:"30s"`
    Debug   bool   `mapstructure:"debug" default:"false"`
}
```

**支持的默认值类型：**
- 字符串: `default:"localhost"`
- 整数: `default:"8080"`
- 布尔值: `default:"true"` 或 `default:"false"`
- 持续时间: `default:"30s"`

## 高级配置模式

### 嵌套配置

```go
type Config struct {
    Server struct {
        HTTP struct {
            Host string `mapstructure:"host" default:"localhost"`
            Port int    `mapstructure:"port" default:"8080"`
        } `mapstructure:"http"`
        
        GRPC struct {
            Host string `mapstructure:"host" default:"localhost"`
            Port int    `mapstructure:"port" default:"9090"`
        } `mapstructure:"grpc"`
    } `mapstructure:"server"`
    
    Database struct {
        Primary struct {
            URL      string `mapstructure:"url"`
            MaxConns int    `mapstructure:"max_connections" default:"10"`
        } `mapstructure:"primary"`
        
        Cache struct {
            URL string `mapstructure:"url"`
            TTL string `mapstructure:"ttl" default:"1h"`
        } `mapstructure:"cache"`
    } `mapstructure:"database"`
}
```

### 数组和切片配置

```go
type Config struct {
    Servers []ServerConfig `mapstructure:"servers"`
    Tags    []string       `mapstructure:"tags"`
}

type ServerConfig struct {
    Name string `mapstructure:"name"`
    URL  string `mapstructure:"url"`
}
```

YAML 配置：
```yaml
servers:
  - name: "primary"
    url: "https://primary.example.com"
  - name: "secondary"
    url: "https://secondary.example.com"

tags:
  - "production"
  - "web-server"
  - "critical"
```

### 映射配置

```go
type Config struct {
    Features map[string]bool   `mapstructure:"features"`
    Limits   map[string]int    `mapstructure:"limits"`
    Metadata map[string]string `mapstructure:"metadata"`
}
```

YAML 配置：
```yaml
features:
  authentication: true
  rate_limiting: false
  metrics: true

limits:
  max_requests: 1000
  max_connections: 100

metadata:
  version: "1.0.0"
  environment: "production"
```

## 环境变量示例

### 基本环境变量

```bash
# 使用前缀 "APP"
export APP_SERVER_HOST=production.example.com
export APP_SERVER_PORT=443
export APP_DEBUG=true
export APP_DATABASE_URL=postgres://user:pass@db.example.com/prod
```

### 嵌套配置

```bash
# 对于嵌套结构
export APP_SERVER_HTTP_HOST=web.example.com
export APP_SERVER_HTTP_PORT=80
export APP_SERVER_GRPC_HOST=grpc.example.com
export APP_SERVER_GRPC_PORT=9090
```

### 数组配置

```bash
# 数组可以指定为逗号分隔的值
export APP_TAGS=production,web-server,critical
```

## 配置文件示例

### 完整的 YAML 示例

```yaml
# 应用程序配置
app:
  name: "我的应用程序"
  version: "1.0.0"
  environment: "production"

# 服务器配置
server:
  http:
    host: "0.0.0.0"
    port: 8080
    timeout: "30s"
  grpc:
    host: "0.0.0.0"
    port: 9090
    timeout: "10s"

# 数据库配置
database:
  primary:
    url: "postgres://user:pass@localhost/myapp"
    max_connections: 25
    timeout: "5s"
  cache:
    url: "redis://localhost:6379"
    ttl: "1h"

# 功能标志
features:
  authentication: true
  rate_limiting: true
  metrics: true
  tracing: false

# 日志配置
logging:
  level: "info"
  format: "json"
  output: ["stdout", "/var/log/app.log"]

# 安全设置
security:
  jwt_secret: "your-secret-key"
  cors_origins: ["https://example.com", "https://app.example.com"]
  rate_limit: 100

debug: false
```

### JSON 示例

```json
{
  "server": {
    "host": "localhost",
    "port": 8080
  },
  "database": {
    "url": "postgres://localhost/myapp",
    "max_connections": 10
  },
  "features": {
    "authentication": true,
    "metrics": false
  },
  "debug": true
}
```

## 验证

### 内置验证

模块自动验证：
- 类型转换（字符串到整数、布尔值等）
- 必需字段（当未提供默认值时）
- 文件格式语法

### 自定义验证

您可以在回调函数中添加自定义验证：

```go
cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
    cfg := currentCfg.(*AppConfig)
    
    // 验证端口范围
    if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
        return fmt.Errorf("无效端口: %d", cfg.Server.Port)
    }
    
    // 验证必需字段
    if cfg.Database.URL == "" {
        return fmt.Errorf("数据库 URL 是必需的")
    }
    
    return nil
})
```

## 最佳实践

1. **使用有意义的默认值** 适用于开发环境
2. **将相关配置分组** 到嵌套结构中
3. **使用环境变量** 处理敏感数据
4. **在回调中验证配置**
5. **记录您的配置** 结构
6. **使用一致的命名** 约定

## 下一步

- [热重载](03_hot_reload.md) - 实现动态配置更新
- [最佳实践](04_best_practices.md) - 遵循推荐模式
- [集成示例](05_integration_examples.md) - 查看实际示例 