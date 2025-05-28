# 故障排除

本指南帮助您诊断和解决使用配置模块时遇到的常见问题。

## 配置加载问题

### 问题：配置文件未找到

**错误信息：**
```
配置文件未找到: config.yaml
```

**可能原因：**
1. 配置文件路径不正确
2. 配置文件不存在
3. 权限问题

**解决方案：**

```go
// 1. 检查文件是否存在
if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
    log.Printf("配置文件不存在: %v", err)
}

// 2. 使用绝对路径
err := config.LoadConfig(&cfg,
    config.WithConfigFile("/path/to/config.yaml", ""),
)

// 3. 使用搜索路径
err := config.LoadConfig(&cfg,
    config.WithConfigFile("config.yaml", "/etc/myapp:/home/user/.config"),
)

// 4. 提供回退配置
err := config.LoadConfig(&cfg,
    config.WithConfigFile("config.yaml", ""),
)
if err != nil {
    log.Printf("使用默认配置: %v", err)
    cfg = getDefaultConfig()
}
```

### 问题：配置文件格式错误

**错误信息：**
```
yaml: line 5: mapping values are not allowed in this context
```

**可能原因：**
1. YAML 语法错误
2. 缩进问题
3. 特殊字符未转义

**解决方案：**

```yaml
# 错误：缩进不一致
server:
  host: localhost
    port: 8080  # 错误的缩进

# 正确：一致的缩进
server:
  host: localhost
  port: 8080

# 错误：特殊字符未引用
password: my@password!

# 正确：特殊字符用引号包围
password: "my@password!"

# 错误：布尔值格式
debug: yes

# 正确：布尔值格式
debug: true
```

**验证配置文件：**

```bash
# 使用 yq 验证 YAML 语法
yq eval '.' config.yaml

# 使用 Python 验证
python -c "import yaml; yaml.safe_load(open('config.yaml'))"
```

### 问题：环境变量未生效

**错误信息：**
```
环境变量 APP_SERVER_PORT 未被识别
```

**可能原因：**
1. 环境变量名称不正确
2. 前缀配置错误
3. 环境变量覆盖未启用

**解决方案：**

```go
// 1. 确保正确的前缀和覆盖设置
err := config.LoadConfig(&cfg,
    config.WithEnvPrefix("APP"),           // 设置前缀
    config.WithEnvVarOverride(true),       // 启用覆盖
)

// 2. 检查环境变量名称映射
// 配置字段: server.port
// 环境变量: APP_SERVER_PORT

// 3. 调试环境变量
for _, env := range os.Environ() {
    if strings.HasPrefix(env, "APP_") {
        log.Printf("环境变量: %s", env)
    }
}

// 4. 手动设置环境变量进行测试
os.Setenv("APP_SERVER_PORT", "9090")
```

## 类型转换问题

### 问题：类型转换失败

**错误信息：**
```
无法将 "invalid" 转换为 int 类型
```

**可能原因：**
1. 配置值类型不匹配
2. 字符串格式错误
3. 默认值类型错误

**解决方案：**

```go
// 1. 确保配置值类型正确
type Config struct {
    // 错误：字符串类型但期望整数
    Port string `mapstructure:"port" default:"8080"`
    
    // 正确：整数类型
    Port int `mapstructure:"port" default:"8080"`
    
    // 错误：布尔值格式
    Debug string `mapstructure:"debug" default:"yes"`
    
    // 正确：布尔值格式
    Debug bool `mapstructure:"debug" default:"true"`
    
    // 时间持续时间
    Timeout time.Duration `mapstructure:"timeout" default:"30s"`
}

// 2. 添加验证
func (c *Config) Validate() error {
    if c.Port < 1 || c.Port > 65535 {
        return fmt.Errorf("无效端口: %d", c.Port)
    }
    return nil
}

// 3. 使用自定义解码钩子
func stringToTimeHookFunc() mapstructure.DecodeHookFunc {
    return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
        if f.Kind() != reflect.String {
            return data, nil
        }
        if t != reflect.TypeOf(time.Duration(0)) {
            return data, nil
        }
        
        return time.ParseDuration(data.(string))
    }
}
```

### 问题：默认值覆盖配置文件值

**错误信息：**
```
配置文件中设置了 enableMetrics: false，但实际加载的值是 true
```

**症状：**
- 配置文件中的布尔值 `false` 被更改为 `true`（如果默认值是 `true`）
- 配置文件中的零值（0、false、""）被结构体标签默认值替换
- 配置文件看似被忽略了某些字段

**示例：**
```yaml
# config.yaml
enableMetrics: false  # 这个值被改为了 true
```

```go
type Config struct {
    EnableMetrics bool `mapstructure:"enableMetrics" default:"true"`
}
```

**根本原因：**
默认值应用逻辑无法区分"用户在配置文件中显式设置为 false"和"字段未设置（零值）"。

**解决方案：**
此问题已在最新版本中通过改进的默认值处理逻辑得到修复：
- 记录配置文件中实际存在的键
- 只对真正未设置的字段应用默认值
- 正确处理用户显式设置的零值

**旧版本的解决方法：**
使用环境变量覆盖有问题的布尔字段：
```bash
export APP_ENABLE_METRICS=false
```

### 问题：默认值未应用

**错误信息：**
```
字段 'host' 为空，但应该有默认值
```

**可能原因：**
1. 默认值标签语法错误
2. 字段未导出
3. mapstructure 标签缺失

**解决方案：**

```go
// 错误：字段未导出
type Config struct {
    host string `mapstructure:"host" default:"localhost"`  // 小写，未导出
}

// 正确：字段导出
type Config struct {
    Host string `mapstructure:"host" default:"localhost"`  // 大写，导出
}

// 错误：缺少 mapstructure 标签
type Config struct {
    Host string `default:"localhost"`
}

// 正确：包含所有必要标签
type Config struct {
    Host string `mapstructure:"host" default:"localhost"`
}

// 错误：默认值类型不匹配
type Config struct {
    Port int `mapstructure:"port" default:"eight-thousand"`  // 字符串默认值用于 int 字段
}

// 正确：匹配的默认值类型
type Config struct {
    Port int `mapstructure:"port" default:"8000"`
}
```

## 热重载问题

### 问题：热重载不工作

**错误信息：**
```
配置文件已更改但回调未触发
```

**可能原因：**
1. 热重载未启用
2. 文件监视器问题
3. 回调注册错误

**解决方案：**

```go
// 1. 确保使用 LoadConfigAndWatch
cm, err := config.LoadConfigAndWatch(&cfg,  // 注意：使用 LoadConfigAndWatch
    config.WithConfigFile("config.yaml", ""),
    config.WithHotReload(true),  // 确保启用热重载
)

// 2. 检查文件权限
info, err := os.Stat("config.yaml")
if err != nil {
    log.Printf("文件状态错误: %v", err)
}
log.Printf("文件权限: %v", info.Mode())

// 3. 添加调试日志
cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
    log.Printf("配置变更回调触发")
    return nil
})

// 4. 检查文件系统事件
// 某些编辑器可能创建临时文件，导致多次触发
// 使用 inotify 工具监控文件事件（Linux）
```

### 问题：回调执行失败

**错误信息：**
```
配置回调执行失败: validation error
```

**可能原因：**
1. 回调中的验证失败
2. 回调中的错误处理不当
3. 回调执行时间过长

**解决方案：**

```go
// 1. 添加错误处理和日志
cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
    cfg := currentCfg.(*AppConfig)
    
    log.Printf("开始处理配置变更")
    
    // 验证新配置
    if err := cfg.Validate(); err != nil {
        log.Printf("配置验证失败: %v", err)
        return err  // 返回错误以阻止应用无效配置
    }
    
    // 应用配置变更
    if err := applyConfiguration(cfg); err != nil {
        log.Printf("应用配置失败: %v", err)
        return err
    }
    
    log.Printf("配置变更处理完成")
    return nil
})

// 2. 使用超时控制
cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    return applyConfigurationWithContext(ctx, currentCfg)
})

// 3. 异步处理长时间操作
cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
    // 快速验证
    if err := quickValidation(currentCfg); err != nil {
        return err
    }
    
    // 异步处理耗时操作
    go func() {
        if err := slowConfigurationUpdate(currentCfg); err != nil {
            log.Printf("异步配置更新失败: %v", err)
        }
    }()
    
    return nil
})
```

## 性能问题

### 问题：配置加载缓慢

**症状：**
- 应用启动时间长
- 配置加载占用大量时间

**可能原因：**
1. 配置文件过大
2. 网络配置源
3. 复杂的默认值处理

**解决方案：**

```go
// 1. 分析配置加载时间
start := time.Now()
err := config.LoadConfig(&cfg, options...)
duration := time.Since(start)
log.Printf("配置加载耗时: %v", duration)

// 2. 分割大配置文件
type AppConfig struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    // ... 其他配置
}

// 分别加载
var serverCfg ServerConfig
config.LoadConfig(&serverCfg, config.WithConfigFile("server.yaml", ""))

// 3. 缓存配置
var configCache *AppConfig
var configMutex sync.RWMutex

func GetConfig() *AppConfig {
    configMutex.RLock()
    defer configMutex.RUnlock()
    return configCache
}

func UpdateConfig(newConfig *AppConfig) {
    configMutex.Lock()
    defer configMutex.Unlock()
    configCache = newConfig
}

// 4. 延迟加载非关键配置
type Config struct {
    // 关键配置立即加载
    Server ServerConfig `mapstructure:"server"`
    
    // 非关键配置延迟加载
    Features map[string]interface{} `mapstructure:"features"`
}
```

### 问题：内存使用过高

**症状：**
- 内存使用持续增长
- 配置相关的内存泄漏

**可能原因：**
1. 配置对象未释放
2. 回调函数中的内存泄漏
3. 文件监视器资源未清理

**解决方案：**

```go
// 1. 正确清理资源
type ConfigManager struct {
    manager config.Manager
    cancel  context.CancelFunc
}

func (cm *ConfigManager) Close() error {
    if cm.cancel != nil {
        cm.cancel()
    }
    return cm.manager.Close()  // 如果有清理方法
}

// 2. 避免在回调中创建大对象
cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
    // 避免：创建大量临时对象
    // largeObject := createLargeObject()
    
    // 推荐：重用对象或使用对象池
    return updateExistingObjects(currentCfg)
})

// 3. 监控内存使用
import _ "net/http/pprof"

go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()

// 使用 go tool pprof 分析内存使用
// go tool pprof http://localhost:6060/debug/pprof/heap
```

## 并发问题

### 问题：并发访问配置

**错误信息：**
```
fatal error: concurrent map read and map write
```

**可能原因：**
1. 多个 goroutine 同时访问配置
2. 热重载期间的并发访问
3. 缺少同步机制

**解决方案：**

```go
// 1. 使用读写锁保护配置访问
type SafeConfig struct {
    mu     sync.RWMutex
    config *AppConfig
}

func (sc *SafeConfig) Get() *AppConfig {
    sc.mu.RLock()
    defer sc.mu.RUnlock()
    return sc.config
}

func (sc *SafeConfig) Update(newConfig *AppConfig) {
    sc.mu.Lock()
    defer sc.mu.Unlock()
    sc.config = newConfig
}

var safeConfig = &SafeConfig{}

// 2. 在回调中安全更新
cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
    cfg := currentCfg.(*AppConfig)
    safeConfig.Update(cfg)
    return nil
})

// 3. 使用原子操作进行简单值
type AtomicConfig struct {
    debugMode int64  // 使用 int64 表示 bool
}

func (ac *AtomicConfig) SetDebugMode(enabled bool) {
    var val int64
    if enabled {
        val = 1
    }
    atomic.StoreInt64(&ac.debugMode, val)
}

func (ac *AtomicConfig) IsDebugMode() bool {
    return atomic.LoadInt64(&ac.debugMode) == 1
}
```

## 调试技巧

### 启用详细日志

```go
// 1. 启用 viper 调试
import "github.com/spf13/viper"

viper.Debug()  // 启用 viper 内部调试日志

// 2. 添加配置加载日志
log.Printf("正在加载配置文件: %s", configFile)
err := config.LoadConfig(&cfg, options...)
if err != nil {
    log.Printf("配置加载失败: %v", err)
} else {
    log.Printf("配置加载成功")
}

// 3. 记录最终配置（隐藏敏感信息）
func logConfig(cfg *AppConfig) {
    safeCfg := *cfg
    safeCfg.Database.Password = "[REDACTED]"
    log.Printf("最终配置: %+v", safeCfg)
}
```

### 配置验证工具

```go
// 创建配置验证工具
func ValidateConfig(cfg *AppConfig) []error {
    var errors []error
    
    // 验证服务器配置
    if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
        errors = append(errors, fmt.Errorf("无效端口: %d", cfg.Server.Port))
    }
    
    // 验证数据库配置
    if cfg.Database.URL == "" {
        errors = append(errors, fmt.Errorf("数据库 URL 不能为空"))
    }
    
    // 验证日志级别
    validLevels := []string{"debug", "info", "warn", "error"}
    if !contains(validLevels, cfg.Log.Level) {
        errors = append(errors, fmt.Errorf("无效日志级别: %s", cfg.Log.Level))
    }
    
    return errors
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}
```

### 配置差异检测

```go
// 检测配置变更
func detectConfigChanges(old, new *AppConfig) []string {
    var changes []string
    
    if old.Server.Port != new.Server.Port {
        changes = append(changes, fmt.Sprintf("服务器端口: %d -> %d", old.Server.Port, new.Server.Port))
    }
    
    if old.Database.URL != new.Database.URL {
        changes = append(changes, "数据库 URL 已更改")
    }
    
    if old.Log.Level != new.Log.Level {
        changes = append(changes, fmt.Sprintf("日志级别: %s -> %s", old.Log.Level, new.Log.Level))
    }
    
    return changes
}

// 在回调中使用
var previousConfig *AppConfig

cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
    newCfg := currentCfg.(*AppConfig)
    
    if previousConfig != nil {
        changes := detectConfigChanges(previousConfig, newCfg)
        for _, change := range changes {
            log.Printf("配置变更: %s", change)
        }
    }
    
    previousConfig = newCfg
    return nil
})
```

## 常见错误模式

### 1. 忘记导出结构体字段

```go
// 错误：字段未导出
type Config struct {
    host string  // 小写，mapstructure 无法访问
    port int     // 小写，mapstructure 无法访问
}

// 正确：字段导出
type Config struct {
    Host string `mapstructure:"host"`
    Port int    `mapstructure:"port"`
}
```

### 2. 错误的标签语法

```go
// 错误：标签语法错误
type Config struct {
    Host string `mapstructure:host`        // 缺少引号
    Port int    `mapstructure:"port" default:8080`  // default 值缺少引号
}

// 正确：正确的标签语法
type Config struct {
    Host string `mapstructure:"host" default:"localhost"`
    Port int    `mapstructure:"port" default:"8080"`
}
```

### 3. 配置文件路径问题

```go
// 错误：硬编码路径
config.WithConfigFile("/home/user/config.yaml", "")

// 更好：使用相对路径和搜索路径
config.WithConfigFile("config.yaml", ".:./config:/etc/myapp")

// 最佳：使用环境变量
configFile := os.Getenv("CONFIG_FILE")
if configFile == "" {
    configFile = "config.yaml"
}
config.WithConfigFile(configFile, "")
```

## 获取帮助

如果您遇到本指南未涵盖的问题：

1. **检查日志**：启用详细日志记录以获取更多信息
2. **验证配置**：使用配置验证工具检查配置正确性
3. **简化测试**：创建最小复现示例
4. **查看源码**：检查配置模块源码以了解内部工作原理
5. **社区支持**：在项目仓库中提交 issue

## 预防措施

1. **始终验证配置**：在生产环境中使用配置验证
2. **使用类型安全**：利用 Go 的类型系统防止配置错误
3. **测试配置加载**：编写单元测试验证配置加载逻辑
4. **监控配置变更**：记录所有配置变更以便审计
5. **备份配置**：保持配置文件的版本控制和备份 