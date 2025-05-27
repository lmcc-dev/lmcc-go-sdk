# 热重载配置

热重载允许您的应用程序自动检测和应用配置更改，而无需重启。这在生产环境中特别有用，可以最大限度地减少停机时间。

## 热重载工作原理

热重载机制使用文件系统监视器来监控配置文件的更改。当检测到更改时：

1. 重新读取配置文件
2. 验证新配置
3. 执行注册的回调函数
4. 更新应用程序状态

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│    配置文件     │    │   文件监视器     │    │     应用程序     │
│     修改        │───▶│    检测变更      │───▶│    执行回调     │
│                 │    │                  │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

## 启用热重载

### 基本设置

```go
package main

import (
    "log"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
    "github.com/spf13/viper"
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
    
    // 加载配置并启用热重载
    cm, err := config.LoadConfigAndWatch(
        &cfg,
        config.WithConfigFile("config.yaml", ""),
        config.WithEnvPrefix("APP"),
        config.WithHotReload(true), // 启用热重载
        config.WithEnvVarOverride(true),
    )
    if err != nil {
        log.Fatalf("加载配置失败: %v", err)
    }
    
    // 注册配置变更回调
    cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
        appCfg := currentCfg.(*AppConfig)
        log.Printf("配置已重新加载！新端口: %d", appCfg.Server.Port)
        return nil
    })
    
    // 您的应用程序逻辑
    log.Printf("应用程序启动，配置: %+v", cfg)
    
    // 保持应用程序运行
    select {}
}
```

## 回调管理

### 全局回调

全局回调在任何配置更改时执行：

```go
cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
    cfg := currentCfg.(*AppConfig)
    
    log.Printf("配置已更新:")
    log.Printf("  服务器: %s:%d", cfg.Server.Host, cfg.Server.Port)
    log.Printf("  调试: %t", cfg.Debug)
    
    // 在这里执行应用程序更新
    return updateApplicationState(cfg)
})
```

### 特定部分回调

您可以为特定配置部分注册回调：

```go
// 仅针对服务器配置更改的回调
cm.RegisterSectionChangeCallback("server", func(v *viper.Viper) error {
    host := v.GetString("server.host")
    port := v.GetInt("server.port")
    
    log.Printf("服务器配置已更改: %s:%d", host, port)
    
    // 使用新配置重启 HTTP 服务器
    return restartHTTPServer(host, port)
})

// 数据库配置更改回调
cm.RegisterSectionChangeCallback("database", func(v *viper.Viper) error {
    url := v.GetString("database.url")
    maxConns := v.GetInt("database.max_connections")
    
    log.Printf("数据库配置已更改")
    
    // 使用新设置重新连接数据库
    return reconnectDatabase(url, maxConns)
})
```

## 回调中的错误处理

### 优雅错误处理

```go
cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
    cfg := currentCfg.(*AppConfig)
    
    // 验证新配置
    if err := validateConfig(cfg); err != nil {
        log.Printf("检测到无效配置: %v", err)
        // 返回错误以防止应用无效配置
        return fmt.Errorf("配置验证失败: %w", err)
    }
    
    // 应用配置更改并处理错误
    if err := applyServerConfig(cfg.Server); err != nil {
        log.Printf("应用服务器配置失败: %v", err)
        // 决定是否返回错误或继续
        return err
    }
    
    if err := applyDatabaseConfig(cfg.Database); err != nil {
        log.Printf("应用数据库配置失败: %v", err)
        // 您可能希望即使数据库配置失败也继续
        log.Printf("继续使用之前的数据库配置")
    }
    
    return nil
})

func validateConfig(cfg *AppConfig) error {
    if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
        return fmt.Errorf("无效端口: %d", cfg.Server.Port)
    }
    
    if cfg.Server.Host == "" {
        return fmt.Errorf("服务器主机不能为空")
    }
    
    return nil
}
```

### 失败时回滚

```go
type ApplicationState struct {
    previousConfig *AppConfig
    currentConfig  *AppConfig
}

var appState = &ApplicationState{}

cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
    newCfg := currentCfg.(*AppConfig)
    
    // 存储之前的配置以便回滚
    appState.previousConfig = appState.currentConfig
    appState.currentConfig = newCfg
    
    // 尝试应用新配置
    if err := applyConfiguration(newCfg); err != nil {
        log.Printf("应用新配置失败: %v", err)
        
        // 回滚到之前的配置
        if appState.previousConfig != nil {
            log.Printf("回滚到之前的配置")
            if rollbackErr := applyConfiguration(appState.previousConfig); rollbackErr != nil {
                log.Printf("严重错误：回滚失败: %v", rollbackErr)
                return fmt.Errorf("配置更新失败且回滚失败: %w", rollbackErr)
            }
            appState.currentConfig = appState.previousConfig
        }
        
        return err
    }
    
    log.Printf("配置成功更新")
    return nil
})
```

## 实际应用示例

### HTTP 服务器重新配置

```go
import (
    "context"
    "net/http"
    "sync"
    "time"
)

type HTTPServerManager struct {
    server *http.Server
    mutex  sync.RWMutex
}

func (hsm *HTTPServerManager) UpdateServer(host string, port int) error {
    hsm.mutex.Lock()
    defer hsm.mutex.Unlock()
    
    // 优雅关闭现有服务器
    if hsm.server != nil {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        
        if err := hsm.server.Shutdown(ctx); err != nil {
            log.Printf("关闭服务器时出错: %v", err)
        }
    }
    
    // 使用更新的配置创建新服务器
    hsm.server = &http.Server{
        Addr:    fmt.Sprintf("%s:%d", host, port),
        Handler: createHandler(), // 您的 HTTP 处理器
    }
    
    // 在后台启动新服务器
    go func() {
        if err := hsm.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Printf("HTTP 服务器错误: %v", err)
        }
    }()
    
    log.Printf("HTTP 服务器在 %s:%d 重启", host, port)
    return nil
}

var serverManager = &HTTPServerManager{}

// 为服务器配置更改注册回调
cm.RegisterSectionChangeCallback("server", func(v *viper.Viper) error {
    host := v.GetString("server.host")
    port := v.GetInt("server.port")
    
    return serverManager.UpdateServer(host, port)
})
```

### 数据库连接池重新配置

```go
import (
    "database/sql"
    "sync"
)

type DatabaseManager struct {
    db    *sql.DB
    mutex sync.RWMutex
}

func (dm *DatabaseManager) UpdateDatabase(url string, maxConns int) error {
    dm.mutex.Lock()
    defer dm.mutex.Unlock()
    
    // 关闭现有连接
    if dm.db != nil {
        dm.db.Close()
    }
    
    // 使用更新的设置创建新连接
    newDB, err := sql.Open("postgres", url)
    if err != nil {
        return fmt.Errorf("打开数据库失败: %w", err)
    }
    
    newDB.SetMaxOpenConns(maxConns)
    newDB.SetMaxIdleConns(maxConns / 2)
    
    // 测试连接
    if err := newDB.Ping(); err != nil {
        newDB.Close()
        return fmt.Errorf("ping 数据库失败: %w", err)
    }
    
    dm.db = newDB
    log.Printf("数据库连接已更新: max_conns=%d", maxConns)
    return nil
}

var dbManager = &DatabaseManager{}

// 为数据库配置更改注册回调
cm.RegisterSectionChangeCallback("database", func(v *viper.Viper) error {
    url := v.GetString("database.url")
    maxConns := v.GetInt("database.max_connections")
    
    return dbManager.UpdateDatabase(url, maxConns)
})
```

## 最佳实践

### 1. 应用前验证

始终在应用更改前验证新配置：

```go
cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
    cfg := currentCfg.(*AppConfig)
    
    // 首先验证
    if err := validateConfiguration(cfg); err != nil {
        return fmt.Errorf("无效配置: %w", err)
    }
    
    // 然后应用
    return applyConfiguration(cfg)
})
```

### 2. 使用优雅关闭

重启服务时，始终使用优雅关闭：

```go
func gracefulRestart(server *http.Server, newAddr string) error {
    // 带超时的优雅关闭
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := server.Shutdown(ctx); err != nil {
        return fmt.Errorf("优雅关闭失败: %w", err)
    }
    
    // 启动新服务器
    return startNewServer(newAddr)
}
```

### 3. 记录配置更改

始终记录配置更改以便调试和审计：

```go
cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
    cfg := currentCfg.(*AppConfig)
    
    log.Printf("检测到配置更改:")
    log.Printf("  时间戳: %s", time.Now().Format(time.RFC3339))
    log.Printf("  服务器: %s:%d", cfg.Server.Host, cfg.Server.Port)
    log.Printf("  调试: %t", cfg.Debug)
    
    return applyConfiguration(cfg)
})
```

### 4. 处理部分失败

设计您的回调以优雅地处理部分失败：

```go
cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
    cfg := currentCfg.(*AppConfig)
    var errors []error
    
    // 尝试独立更新每个组件
    if err := updateServerConfig(cfg.Server); err != nil {
        errors = append(errors, fmt.Errorf("服务器配置: %w", err))
    }
    
    if err := updateDatabaseConfig(cfg.Database); err != nil {
        errors = append(errors, fmt.Errorf("数据库配置: %w", err))
    }
    
    if err := updateLoggingConfig(cfg.Logging); err != nil {
        errors = append(errors, fmt.Errorf("日志配置: %w", err))
    }
    
    // 如果有错误则返回组合错误
    if len(errors) > 0 {
        return fmt.Errorf("配置更新错误: %v", errors)
    }
    
    return nil
})
```

## 故障排除

### 常见问题

1. **文件权限**：确保应用程序对配置文件有读取权限
2. **文件锁定**：某些编辑器创建的临时文件可能触发错误重载
3. **快速更改**：多个快速更改可能导致回调泛滥

### 调试热重载

启用调试日志以排除热重载问题：

```go
import "github.com/spf13/viper"

// 启用 viper 调试日志
viper.SetConfigType("yaml")
viper.Debug() // 这将启用调试输出
```

## 下一步

- [最佳实践](04_best_practices.md) - 遵循推荐模式
- [集成示例](05_integration_examples.md) - 查看实际示例
- [故障排除](06_troubleshooting.md) - 常见问题和解决方案 