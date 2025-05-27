# 配置最佳实践

本指南概述了在生产应用程序中有效使用配置模块的推荐模式和实践。

## 配置结构设计

### 1. 使用分层组织

将您的配置组织成逻辑组：

```go
type AppConfig struct {
    // 应用程序元数据
    App struct {
        Name        string `mapstructure:"name" default:"my-app"`
        Version     string `mapstructure:"version" default:"1.0.0"`
        Environment string `mapstructure:"environment" default:"development"`
    } `mapstructure:"app"`
    
    // 服务器配置
    Server struct {
        HTTP struct {
            Host         string        `mapstructure:"host" default:"localhost"`
            Port         int           `mapstructure:"port" default:"8080"`
            ReadTimeout  time.Duration `mapstructure:"read_timeout" default:"30s"`
            WriteTimeout time.Duration `mapstructure:"write_timeout" default:"30s"`
        } `mapstructure:"http"`
        
        GRPC struct {
            Host string `mapstructure:"host" default:"localhost"`
            Port int    `mapstructure:"port" default:"9090"`
        } `mapstructure:"grpc"`
    } `mapstructure:"server"`
    
    // 数据库配置
    Database struct {
        Primary struct {
            URL          string        `mapstructure:"url"`
            MaxConns     int           `mapstructure:"max_connections" default:"25"`
            MaxIdleConns int           `mapstructure:"max_idle_connections" default:"5"`
            ConnTimeout  time.Duration `mapstructure:"connection_timeout" default:"5s"`
        } `mapstructure:"primary"`
        
        Cache struct {
            URL string        `mapstructure:"url"`
            TTL time.Duration `mapstructure:"ttl" default:"1h"`
        } `mapstructure:"cache"`
    } `mapstructure:"database"`
    
    // 可观测性
    Observability struct {
        Logging struct {
            Level  string `mapstructure:"level" default:"info"`
            Format string `mapstructure:"format" default:"json"`
        } `mapstructure:"logging"`
        
        Metrics struct {
            Enabled bool   `mapstructure:"enabled" default:"true"`
            Port    int    `mapstructure:"port" default:"9090"`
            Path    string `mapstructure:"path" default:"/metrics"`
        } `mapstructure:"metrics"`
        
        Tracing struct {
            Enabled    bool   `mapstructure:"enabled" default:"false"`
            Endpoint   string `mapstructure:"endpoint"`
            SampleRate float64 `mapstructure:"sample_rate" default:"0.1"`
        } `mapstructure:"tracing"`
    } `mapstructure:"observability"`
}
```

### 2. 使用有意义的默认值

提供适用于开发的合理默认值：

```go
type DatabaseConfig struct {
    // 好：为开发提供可用的默认值
    URL string `mapstructure:"url" default:"postgres://localhost:5432/myapp_dev"`
    
    // 好：在任何地方都能工作的保守默认值
    MaxConnections int `mapstructure:"max_connections" default:"10"`
    
    // 好：合理的超时时间
    Timeout time.Duration `mapstructure:"timeout" default:"30s"`
    
    // 不好：必需字段没有默认值
    // Password string `mapstructure:"password"`
    
    // 更好：对敏感数据使用环境变量
    Password string `mapstructure:"password" default:"${DB_PASSWORD}"`
}
```

### 3. 使用适当的数据类型

为您的配置选择正确的数据类型：

```go
type Config struct {
    // 对基于时间的值使用 time.Duration
    Timeout time.Duration `mapstructure:"timeout" default:"30s"`
    
    // 对更好的验证使用特定类型
    Port int `mapstructure:"port" default:"8080"`
    
    // 对功能标志使用 bool
    EnableMetrics bool `mapstructure:"enable_metrics" default:"true"`
    
    // 对列表使用切片
    AllowedOrigins []string `mapstructure:"allowed_origins"`
    
    // 对键值对使用映射
    Headers map[string]string `mapstructure:"headers"`
    
    // 对验证使用自定义类型
    LogLevel LogLevel `mapstructure:"log_level" default:"info"`
}

type LogLevel string

const (
    LogLevelDebug LogLevel = "debug"
    LogLevelInfo  LogLevel = "info"
    LogLevelWarn  LogLevel = "warn"
    LogLevelError LogLevel = "error"
)
```

## 环境变量策略

### 1. 使用一致的命名

建立清晰的命名约定：

```go
// 好：一致的前缀和结构
// APP_SERVER_HTTP_HOST
// APP_SERVER_HTTP_PORT
// APP_DATABASE_PRIMARY_URL
// APP_OBSERVABILITY_LOGGING_LEVEL

config.WithEnvPrefix("APP")
```

### 2. 分离敏感数据

将敏感数据保存在环境变量中：

```go
type Config struct {
    Database struct {
        // 非敏感：可以在配置文件中
        Host string `mapstructure:"host" default:"localhost"`
        Port int    `mapstructure:"port" default:"5432"`
        Name string `mapstructure:"name" default:"myapp"`
        
        // 敏感：应该在环境变量中
        Username string `mapstructure:"username"`
        Password string `mapstructure:"password"`
    } `mapstructure:"database"`
    
    Security struct {
        // 敏感：应该在环境变量中
        JWTSecret    string `mapstructure:"jwt_secret"`
        APIKey       string `mapstructure:"api_key"`
        EncryptionKey string `mapstructure:"encryption_key"`
    } `mapstructure:"security"`
}
```

环境变量：
```bash
export APP_DATABASE_USERNAME=myuser
export APP_DATABASE_PASSWORD=mysecretpassword
export APP_SECURITY_JWT_SECRET=my-jwt-secret-key
export APP_SECURITY_API_KEY=my-api-key
```

### 3. 使用特定环境的覆盖

```bash
# 开发环境
export APP_APP_ENVIRONMENT=development
export APP_OBSERVABILITY_LOGGING_LEVEL=debug
export APP_DATABASE_PRIMARY_URL=postgres://localhost:5432/myapp_dev

# 生产环境
export APP_APP_ENVIRONMENT=production
export APP_OBSERVABILITY_LOGGING_LEVEL=warn
export APP_DATABASE_PRIMARY_URL=postgres://prod-db:5432/myapp_prod
```

## 配置验证

### 1. 实现验证函数

```go
func (c *AppConfig) Validate() error {
    var errors []error
    
    // 验证服务器配置
    if c.Server.HTTP.Port < 1 || c.Server.HTTP.Port > 65535 {
        errors = append(errors, fmt.Errorf("无效的 HTTP 端口: %d", c.Server.HTTP.Port))
    }
    
    if c.Server.GRPC.Port < 1 || c.Server.GRPC.Port > 65535 {
        errors = append(errors, fmt.Errorf("无效的 GRPC 端口: %d", c.Server.GRPC.Port))
    }
    
    // 验证数据库配置
    if c.Database.Primary.URL == "" {
        errors = append(errors, fmt.Errorf("数据库 URL 是必需的"))
    }
    
    if c.Database.Primary.MaxConns < 1 {
        errors = append(errors, fmt.Errorf("最大连接数必须为正数"))
    }
    
    // 验证可观测性配置
    validLogLevels := map[string]bool{
        "debug": true, "info": true, "warn": true, "error": true,
    }
    if !validLogLevels[c.Observability.Logging.Level] {
        errors = append(errors, fmt.Errorf("无效的日志级别: %s", c.Observability.Logging.Level))
    }
    
    if len(errors) > 0 {
        return fmt.Errorf("配置验证失败: %v", errors)
    }
    
    return nil
}
```

### 2. 在回调中使用验证

```go
cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
    cfg := currentCfg.(*AppConfig)
    
    // 应用前验证配置
    if err := cfg.Validate(); err != nil {
        return fmt.Errorf("无效配置: %w", err)
    }
    
    // 应用配置
    return applyConfiguration(cfg)
})
```

## 错误处理模式

### 1. 优雅降级

```go
func loadConfiguration() (*AppConfig, error) {
    var cfg AppConfig
    
    err := config.LoadConfig(
        &cfg,
        config.WithConfigFile("config.yaml", ""),
        config.WithEnvPrefix("APP"),
        config.WithEnvVarOverride(true),
    )
    
    if err != nil {
        // 记录错误但继续使用默认值
        log.Printf("警告：加载配置文件失败: %v", err)
        log.Printf("使用默认配置")
        
        // 使用默认值初始化
        cfg = AppConfig{} // 这将使用结构体标签默认值
    }
    
    // 始终验证最终配置
    if err := cfg.Validate(); err != nil {
        return nil, fmt.Errorf("配置验证失败: %w", err)
    }
    
    return &cfg, nil
}
```

### 2. 配置回退链

```go
func loadConfigurationWithFallback() (*AppConfig, error) {
    var cfg AppConfig
    var err error
    
    // 尝试主配置文件
    err = config.LoadConfig(&cfg, 
        config.WithConfigFile("config.yaml", ""),
        config.WithEnvPrefix("APP"),
    )
    
    if err != nil {
        log.Printf("主配置失败: %v", err)
        
        // 尝试回退配置文件
        err = config.LoadConfig(&cfg,
            config.WithConfigFile("config.default.yaml", ""),
            config.WithEnvPrefix("APP"),
        )
        
        if err != nil {
            log.Printf("回退配置失败: %v", err)
            
            // 使用嵌入的默认值
            cfg = getDefaultConfiguration()
        }
    }
    
    return &cfg, nil
}
```

## 测试策略

### 1. 配置测试

```go
func TestConfigurationLoading(t *testing.T) {
    tests := []struct {
        name           string
        configContent  string
        envVars        map[string]string
        expectedConfig AppConfig
        expectError    bool
    }{
        {
            name: "有效配置",
            configContent: `
server:
  http:
    host: "0.0.0.0"
    port: 8080
database:
  primary:
    url: "postgres://localhost/test"
`,
            envVars: map[string]string{
                "APP_DATABASE_PRIMARY_URL": "postgres://test-db/test",
            },
            expectedConfig: AppConfig{
                Server: ServerConfig{
                    HTTP: HTTPConfig{
                        Host: "0.0.0.0",
                        Port: 8080,
                    },
                },
                Database: DatabaseConfig{
                    Primary: PrimaryDBConfig{
                        URL: "postgres://test-db/test", // 来自环境变量
                    },
                },
            },
            expectError: false,
        },
        {
            name: "无效端口",
            configContent: `
server:
  http:
    port: 99999
`,
            expectError: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 设置环境变量
            for key, value := range tt.envVars {
                os.Setenv(key, value)
                defer os.Unsetenv(key)
            }
            
            // 创建临时配置文件
            tmpFile, err := os.CreateTemp("", "config-*.yaml")
            require.NoError(t, err)
            defer os.Remove(tmpFile.Name())
            
            _, err = tmpFile.WriteString(tt.configContent)
            require.NoError(t, err)
            tmpFile.Close()
            
            // 加载配置
            var cfg AppConfig
            err = config.LoadConfig(&cfg,
                config.WithConfigFile(tmpFile.Name(), ""),
                config.WithEnvPrefix("APP"),
                config.WithEnvVarOverride(true),
            )
            
            if tt.expectError {
                assert.Error(t, err)
                return
            }
            
            require.NoError(t, err)
            assert.Equal(t, tt.expectedConfig.Server.HTTP.Host, cfg.Server.HTTP.Host)
            assert.Equal(t, tt.expectedConfig.Server.HTTP.Port, cfg.Server.HTTP.Port)
        })
    }
}
```

### 2. 热重载测试

```go
func TestHotReload(t *testing.T) {
    // 创建临时配置文件
    tmpFile, err := os.CreateTemp("", "config-*.yaml")
    require.NoError(t, err)
    defer os.Remove(tmpFile.Name())
    
    // 初始配置
    initialConfig := `
server:
  http:
    port: 8080
`
    _, err = tmpFile.WriteString(initialConfig)
    require.NoError(t, err)
    tmpFile.Close()
    
    var cfg AppConfig
    var callbackCalled bool
    var newPort int
    
    // 加载并启用热重载
    cm, err := config.LoadConfigAndWatch(&cfg,
        config.WithConfigFile(tmpFile.Name(), ""),
        config.WithHotReload(true),
    )
    require.NoError(t, err)
    
    // 注册回调
    cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
        appCfg := currentCfg.(*AppConfig)
        callbackCalled = true
        newPort = appCfg.Server.HTTP.Port
        return nil
    })
    
    // 验证初始配置
    assert.Equal(t, 8080, cfg.Server.HTTP.Port)
    
    // 更新配置文件
    updatedConfig := `
server:
  http:
    port: 9090
`
    err = os.WriteFile(tmpFile.Name(), []byte(updatedConfig), 0644)
    require.NoError(t, err)
    
    // 等待热重载（带超时）
    timeout := time.After(5 * time.Second)
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()
    
    for {
        select {
        case <-timeout:
            t.Fatal("热重载回调在超时内未被调用")
        case <-ticker.C:
            if callbackCalled && newPort == 9090 {
                return // 成功
            }
        }
    }
}
```

## 生产部署

### 1. 配置管理

```yaml
# config/production.yaml
app:
  name: "my-app"
  version: "1.2.3"
  environment: "production"

server:
  http:
    host: "0.0.0.0"
    port: 8080
    read_timeout: "30s"
    write_timeout: "30s"

database:
  primary:
    # URL 来自环境变量
    max_connections: 50
    max_idle_connections: 10
    connection_timeout: "10s"

observability:
  logging:
    level: "info"
    format: "json"
  metrics:
    enabled: true
    port: 9090
  tracing:
    enabled: true
    sample_rate: 0.1
```

### 2. Docker 配置

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o myapp .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/myapp .
COPY --from=builder /app/config/production.yaml ./config/

# 配置来自环境变量
ENV APP_APP_ENVIRONMENT=production
ENV APP_SERVER_HTTP_HOST=0.0.0.0
ENV APP_SERVER_HTTP_PORT=8080

CMD ["./myapp"]
```

### 3. Kubernetes 配置

```yaml
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  config.yaml: |
    app:
      name: "my-app"
      environment: "production"
    server:
      http:
        host: "0.0.0.0"
        port: 8080
    observability:
      logging:
        level: "info"
        format: "json"

---
# k8s/secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
type: Opaque
data:
  database-url: <base64-encoded-database-url>
  jwt-secret: <base64-encoded-jwt-secret>

---
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
spec:
  template:
    spec:
      containers:
      - name: my-app
        image: my-app:latest
        env:
        - name: APP_DATABASE_PRIMARY_URL
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: database-url
        - name: APP_SECURITY_JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: jwt-secret
        volumeMounts:
        - name: config
          mountPath: /app/config
      volumes:
      - name: config
        configMap:
          name: app-config
```

## 性能考虑

### 1. 最小化热重载影响

```go
// 好：只更新更改内容的高效回调
cm.RegisterSectionChangeCallback("server", func(v *viper.Viper) error {
    newPort := v.GetInt("server.http.port")
    if currentPort != newPort {
        return restartHTTPServer(newPort)
    }
    return nil // 没有更改，不需要操作
})

// 不好：总是重启所有内容的低效回调
cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
    // 即使是微小更改也会重启所有内容
    return restartAllServices(currentCfg)
})
```

### 2. 缓存配置值

```go
type ConfigCache struct {
    mu     sync.RWMutex
    config *AppConfig
}

func (cc *ConfigCache) GetConfig() *AppConfig {
    cc.mu.RLock()
    defer cc.mu.RUnlock()
    return cc.config
}

func (cc *ConfigCache) UpdateConfig(newConfig *AppConfig) {
    cc.mu.Lock()
    defer cc.mu.Unlock()
    cc.config = newConfig
}

var configCache = &ConfigCache{}

// 在回调中更新缓存
cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
    cfg := currentCfg.(*AppConfig)
    configCache.UpdateConfig(cfg)
    return nil
})
```

## 安全最佳实践

### 1. 保护敏感配置

```go
// 永远不要记录敏感配置
func logConfiguration(cfg *AppConfig) {
    // 好：创建安全副本用于记录
    safeCfg := *cfg
    safeCfg.Database.Primary.Password = "[REDACTED]"
    safeCfg.Security.JWTSecret = "[REDACTED]"
    
    log.Printf("配置已加载: %+v", safeCfg)
}
```

### 2. 验证文件权限

```go
func validateConfigFilePermissions(filename string) error {
    info, err := os.Stat(filename)
    if err != nil {
        return err
    }
    
    mode := info.Mode()
    if mode&0077 != 0 {
        return fmt.Errorf("配置文件 %s 权限过于宽松: %v", filename, mode)
    }
    
    return nil
}
```

## 下一步

- [集成示例](05_integration_examples.md) - 查看实际示例
- [故障排除](06_troubleshooting.md) - 常见问题和解决方案 