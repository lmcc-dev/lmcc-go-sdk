# 模块规范

本文档提供了 `pkg/config` 模块的详细规范，概述了其公共 API、接口、函数和预定义组件。该模块通过提供热重载、环境变量集成和灵活的配置加载模式来增强 Go 的配置管理能力。

## 1. 概述

`pkg/config` 模块旨在使配置管理更加健壮且对开发者友好。主要特性包括：
- 从多种文件格式（YAML、JSON、TOML 等）自动加载配置
- 支持前缀的环境变量集成
- 动态配置更新的热重载功能
- 配置变更通知的回调系统
- 与流行配置库（Viper）的集成

## 2. 核心函数

### 配置加载函数

#### LoadConfig
```go
func LoadConfig(cfg interface{}, opts ...Option) error
```
从文件和环境变量加载配置到提供的结构体中。

**参数：**
- `cfg`：配置结构体的指针
- `opts`：配置加载的可变选项

**返回值：**
- `error`：配置加载失败时的错误

**示例：**
```go
type AppConfig struct {
    Server struct {
        Host string `mapstructure:"host" default:"localhost"`
        Port int    `mapstructure:"port" default:"8080"`
    } `mapstructure:"server"`
}

var cfg AppConfig
err := config.LoadConfig(&cfg,
    config.WithConfigFile("config.yaml", ""),
    config.WithEnvPrefix("APP"),
)
```

#### LoadConfigAndWatch
```go
func LoadConfigAndWatch(cfg interface{}, opts ...Option) (*ConfigManager, error)
```
加载配置并返回用于热重载功能的 ConfigManager。

**参数：**
- `cfg`：配置结构体的指针
- `opts`：配置加载的可变选项

**返回值：**
- `*ConfigManager`：处理配置更新的管理器
- `error`：配置加载失败时的错误

**示例：**
```go
cm, err := config.LoadConfigAndWatch(&cfg,
    config.WithConfigFile("config.yaml", ""),
    config.WithHotReload(true),
)
```

## 3. 配置选项

### Option 类型
```go
type Option func(*ConfigOptions)
```
用于配置加载行为的函数选项类型。

### 可用选项

#### WithConfigFile
```go
func WithConfigFile(filename, searchPaths string) Option
```
指定配置文件和可选的搜索路径。

**参数：**
- `filename`：配置文件名称
- `searchPaths`：要搜索的目录的冒号分隔列表

#### WithEnvPrefix
```go
func WithEnvPrefix(prefix string) Option
```
设置环境变量绑定的前缀。

**参数：**
- `prefix`：用于环境变量的前缀

#### WithEnvVarOverride
```go
func WithEnvVarOverride(enable bool) Option
```
启用或禁用环境变量覆盖配置值。

**参数：**
- `enable`：是否启用环境变量覆盖

#### WithHotReload
```go
func WithHotReload(enable bool) Option
```
启用或禁用热重载功能。

**参数：**
- `enable`：是否启用热重载

#### WithEnvKeyReplacer
```go
func WithEnvKeyReplacer(replacer *strings.Replacer) Option
```
为环境变量名称设置自定义键替换器。

**参数：**
- `replacer`：用于转换配置键的字符串替换器

## 4. ConfigManager 接口

`ConfigManager` 提供管理配置更新和回调的方法。

### 方法

#### RegisterCallback
```go
func (cm *ConfigManager) RegisterCallback(callback func(*viper.Viper, interface{}) error)
```
注册在任何配置更改时调用的全局回调函数。

**参数：**
- `callback`：配置更改时调用的函数
  - `*viper.Viper`：包含新配置的 Viper 实例
  - `interface{}`：更新的配置结构体
  - `error`：回调返回的错误

#### RegisterSectionChangeCallback
```go
func (cm *ConfigManager) RegisterSectionChangeCallback(section string, callback func(*viper.Viper) error)
```
为特定配置部分更改注册回调。

**参数：**
- `section`：要监视的配置部分（例如，"server"、"database"）
- `callback`：部分更改时调用的函数
  - `*viper.Viper`：包含新配置的 Viper 实例
  - `error`：回调返回的错误

#### Stop
```go
func (cm *ConfigManager) Stop()
```
停止配置监视器并清理资源。

## 5. 配置结构体标签

### mapstructure 标签
将配置键映射到结构体字段。

```go
type Config struct {
    ServerPort int `mapstructure:"server_port"`
    DebugMode  bool `mapstructure:"debug"`
}
```

### default 标签
为配置中未提供的字段指定默认值。

```go
type Config struct {
    Host    string `mapstructure:"host" default:"localhost"`
    Port    int    `mapstructure:"port" default:"8080"`
    Timeout string `mapstructure:"timeout" default:"30s"`
    Debug   bool   `mapstructure:"debug" default:"false"`
}
```

**支持的默认值类型：**
- 字符串：`default:"localhost"`
- 整数：`default:"8080"`
- 布尔值：`default:"true"` 或 `default:"false"`
- 持续时间：`default:"30s"`

## 6. 支持的文件格式

模块支持多种配置文件格式：

| 格式 | 扩展名 | 示例 |
|------|--------|------|
| YAML | `.yaml`, `.yml` | `config.yaml` |
| JSON | `.json` | `config.json` |
| TOML | `.toml` | `config.toml` |
| HCL | `.hcl` | `config.hcl` |
| INI | `.ini` | `config.ini` |
| Properties | `.properties` | `config.properties` |

## 7. 环境变量映射

### 自动映射
配置键自动映射到环境变量：

- 嵌套键使用下划线：`server.host` → `PREFIX_SERVER_HOST`
- 支持数组索引：`servers[0].port` → `PREFIX_SERVERS_0_PORT`
- 自定义键替换器可以修改映射

### 映射示例
```go
// 配置结构
type Config struct {
    Server struct {
        Host string `mapstructure:"host"`
        Port int    `mapstructure:"port"`
    } `mapstructure:"server"`
    Database struct {
        URL string `mapstructure:"url"`
    } `mapstructure:"database"`
}

// 使用前缀 "APP"，环境变量：
// APP_SERVER_HOST=localhost
// APP_SERVER_PORT=8080
// APP_DATABASE_URL=postgres://localhost/myapp
```

## 8. 热重载机制

### 文件监视
热重载机制使用文件系统监视器来监控配置文件：

1. **文件更改检测**：监控配置文件的修改
2. **验证**：在应用之前验证新配置
3. **回调执行**：使用新配置执行注册的回调
4. **错误处理**：提供错误处理和回滚功能

### 回调类型

#### 全局回调
对任何配置更改都会调用：
```go
cm.RegisterCallback(func(v *viper.Viper, cfg interface{}) error {
    // 处理任何配置更改
    return nil
})
```

#### 部分回调
仅在特定部分更改时调用：
```go
cm.RegisterSectionChangeCallback("server", func(v *viper.Viper) error {
    // 处理服务器配置更改
    return nil
})
```

## 9. 错误处理

### 配置加载错误
- **文件未找到**：如果配置文件不存在则返回错误
- **解析错误**：对于无效的文件格式或语法返回错误
- **验证错误**：对于无效的配置值返回错误
- **环境变量错误**：对于无效的环境变量值返回错误

### 热重载错误
- **回调错误**：如果回调返回错误，配置更新将被拒绝
- **验证错误**：无效的新配置被拒绝，保持之前的配置
- **文件系统错误**：文件系统监视器错误被记录但不会停止应用程序

## 10. 与 Viper 的集成

该模块基于 [Viper](https://github.com/spf13/viper) 构建，提供：

### Viper 特性
- 多种配置文件格式支持
- 环境变量集成
- 配置键大小写不敏感
- 配置值类型转换

### 扩展特性
- 使用函数选项的简化 API
- 热重载功能
- 配置更改的回调系统
- 通过结构体标签支持默认值

## 11. 最佳实践

### 配置结构设计
```go
// 好的：有组织的嵌套结构
type Config struct {
    App struct {
        Name    string `mapstructure:"name" default:"myapp"`
        Version string `mapstructure:"version" default:"1.0.0"`
    } `mapstructure:"app"`
    
    Server struct {
        Host string `mapstructure:"host" default:"localhost"`
        Port int    `mapstructure:"port" default:"8080"`
    } `mapstructure:"server"`
}

// 避免：复杂配置的扁平结构
type Config struct {
    AppName     string `mapstructure:"app_name"`
    AppVersion  string `mapstructure:"app_version"`
    ServerHost  string `mapstructure:"server_host"`
    ServerPort  int    `mapstructure:"server_port"`
}
```

### 错误处理
```go
// 始终处理配置加载错误
if err := config.LoadConfig(&cfg, opts...); err != nil {
    log.Fatalf("加载配置失败: %v", err)
}

// 适当处理回调错误
cm.RegisterCallback(func(v *viper.Viper, cfg interface{}) error {
    if err := validateConfig(cfg); err != nil {
        return fmt.Errorf("无效配置: %w", err)
    }
    return applyConfig(cfg)
})
```

### 特定环境配置
```go
func getConfigOptions() []config.Option {
    env := os.Getenv("APP_ENV")
    
    opts := []config.Option{
        config.WithEnvPrefix("APP"),
        config.WithEnvVarOverride(true),
    }
    
    switch env {
    case "development":
        opts = append(opts, config.WithConfigFile("config.dev.yaml", ""))
    case "production":
        opts = append(opts, config.WithConfigFile("config.prod.yaml", ""))
    default:
        opts = append(opts, config.WithConfigFile("config.yaml", ""))
    }
    
    return opts
}
```

## 12. 性能考虑

### 文件监视开销
- 文件系统监视器具有最小开销
- 仅在启用热重载时活跃
- ConfigManager 停止时自动清理

### 内存使用
- 配置加载到内存一次
- 热重载创建新实例但清理旧实例
- 回调函数应避免内存泄漏

### 回调性能
- 回调应该快速以避免阻塞配置更新
- 长时间运行的操作应异步执行
- 回调中的错误处理应该健壮

## 13. 线程安全

### 并发访问
- 初始化期间配置加载不是线程安全的
- 热重载回调按顺序执行
- 应用程序应为配置访问实现自己的同步

### 推荐模式
```go
type Application struct {
    config *AppConfig
    mutex  sync.RWMutex
}

func (app *Application) GetConfig() *AppConfig {
    app.mutex.RLock()
    defer app.mutex.RUnlock()
    return app.config
}

func (app *Application) updateConfig(newConfig *AppConfig) {
    app.mutex.Lock()
    defer app.mutex.Unlock()
    app.config = newConfig
}
```

本规范提供了如何有效使用 `pkg/config` 模块的全面理解。有关更多示例和最佳实践，请参阅使用指南。 