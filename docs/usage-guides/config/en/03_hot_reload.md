# Hot Reload Configuration

Hot reload allows your application to automatically detect and apply configuration changes without requiring a restart. This is particularly useful for production environments where downtime should be minimized.

## How Hot Reload Works

The hot reload mechanism uses file system watchers to monitor configuration files for changes. When a change is detected:

1. The configuration file is re-read
2. The new configuration is validated
3. Registered callbacks are executed
4. The application state is updated

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│  Config File    │    │   File Watcher   │    │   Application   │
│   Modified      │───▶│    Detects       │───▶│    Callbacks    │
│                 │    │    Change        │    │    Executed     │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

## Enabling Hot Reload

### Basic Setup

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
    
    // Load configuration with hot reload enabled
    cm, err := config.LoadConfigAndWatch(
        &cfg,
        config.WithConfigFile("config.yaml", ""),
        config.WithEnvPrefix("APP"),
        config.WithHotReload(true), // Enable hot reload
        config.WithEnvVarOverride(true),
    )
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }
    
    // Register callback for configuration changes
    cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
        appCfg := currentCfg.(*AppConfig)
        log.Printf("Configuration reloaded! New port: %d", appCfg.Server.Port)
        return nil
    })
    
    // Your application logic here
    log.Printf("Application started with config: %+v", cfg)
    
    // Keep the application running
    select {}
}
```

## Callback Management

### Global Callbacks

Global callbacks are executed for any configuration change:

```go
cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
    cfg := currentCfg.(*AppConfig)
    
    log.Printf("Configuration updated:")
    log.Printf("  Server: %s:%d", cfg.Server.Host, cfg.Server.Port)
    log.Printf("  Debug: %t", cfg.Debug)
    
    // Perform application updates here
    return updateApplicationState(cfg)
})
```

### Section-Specific Callbacks

You can register callbacks for specific configuration sections:

```go
// Callback for server configuration changes only
cm.RegisterSectionChangeCallback("server", func(v *viper.Viper) error {
    host := v.GetString("server.host")
    port := v.GetInt("server.port")
    
    log.Printf("Server configuration changed: %s:%d", host, port)
    
    // Restart HTTP server with new configuration
    return restartHTTPServer(host, port)
})

// Callback for database configuration changes
cm.RegisterSectionChangeCallback("database", func(v *viper.Viper) error {
    url := v.GetString("database.url")
    maxConns := v.GetInt("database.max_connections")
    
    log.Printf("Database configuration changed")
    
    // Reconnect to database with new settings
    return reconnectDatabase(url, maxConns)
})
```

## Error Handling in Callbacks

### Graceful Error Handling

```go
cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
    cfg := currentCfg.(*AppConfig)
    
    // Validate the new configuration
    if err := validateConfig(cfg); err != nil {
        log.Printf("Invalid configuration detected: %v", err)
        // Return error to prevent applying invalid config
        return fmt.Errorf("configuration validation failed: %w", err)
    }
    
    // Apply configuration changes with error handling
    if err := applyServerConfig(cfg.Server); err != nil {
        log.Printf("Failed to apply server config: %v", err)
        // Decide whether to return error or continue
        return err
    }
    
    if err := applyDatabaseConfig(cfg.Database); err != nil {
        log.Printf("Failed to apply database config: %v", err)
        // You might want to continue even if database config fails
        log.Printf("Continuing with previous database configuration")
    }
    
    return nil
})

func validateConfig(cfg *AppConfig) error {
    if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
        return fmt.Errorf("invalid port: %d", cfg.Server.Port)
    }
    
    if cfg.Server.Host == "" {
        return fmt.Errorf("server host cannot be empty")
    }
    
    return nil
}
```

### Rollback on Failure

```go
type ApplicationState struct {
    previousConfig *AppConfig
    currentConfig  *AppConfig
}

var appState = &ApplicationState{}

cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
    newCfg := currentCfg.(*AppConfig)
    
    // Store previous configuration for rollback
    appState.previousConfig = appState.currentConfig
    appState.currentConfig = newCfg
    
    // Try to apply new configuration
    if err := applyConfiguration(newCfg); err != nil {
        log.Printf("Failed to apply new configuration: %v", err)
        
        // Rollback to previous configuration
        if appState.previousConfig != nil {
            log.Printf("Rolling back to previous configuration")
            if rollbackErr := applyConfiguration(appState.previousConfig); rollbackErr != nil {
                log.Printf("CRITICAL: Rollback failed: %v", rollbackErr)
                return fmt.Errorf("configuration update failed and rollback failed: %w", rollbackErr)
            }
            appState.currentConfig = appState.previousConfig
        }
        
        return err
    }
    
    log.Printf("Configuration successfully updated")
    return nil
})
```

## Real-World Examples

### HTTP Server Reconfiguration

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
    
    // Gracefully shutdown existing server
    if hsm.server != nil {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        
        if err := hsm.server.Shutdown(ctx); err != nil {
            log.Printf("Error shutting down server: %v", err)
        }
    }
    
    // Create new server with updated configuration
    hsm.server = &http.Server{
        Addr:    fmt.Sprintf("%s:%d", host, port),
        Handler: createHandler(), // Your HTTP handler
    }
    
    // Start new server in background
    go func() {
        if err := hsm.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Printf("HTTP server error: %v", err)
        }
    }()
    
    log.Printf("HTTP server restarted on %s:%d", host, port)
    return nil
}

var serverManager = &HTTPServerManager{}

// Register callback for server configuration changes
cm.RegisterSectionChangeCallback("server", func(v *viper.Viper) error {
    host := v.GetString("server.host")
    port := v.GetInt("server.port")
    
    return serverManager.UpdateServer(host, port)
})
```

### Database Connection Pool Reconfiguration

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
    
    // Close existing connection
    if dm.db != nil {
        dm.db.Close()
    }
    
    // Create new connection with updated settings
    newDB, err := sql.Open("postgres", url)
    if err != nil {
        return fmt.Errorf("failed to open database: %w", err)
    }
    
    newDB.SetMaxOpenConns(maxConns)
    newDB.SetMaxIdleConns(maxConns / 2)
    
    // Test the connection
    if err := newDB.Ping(); err != nil {
        newDB.Close()
        return fmt.Errorf("failed to ping database: %w", err)
    }
    
    dm.db = newDB
    log.Printf("Database connection updated: max_conns=%d", maxConns)
    return nil
}

var dbManager = &DatabaseManager{}

// Register callback for database configuration changes
cm.RegisterSectionChangeCallback("database", func(v *viper.Viper) error {
    url := v.GetString("database.url")
    maxConns := v.GetInt("database.max_connections")
    
    return dbManager.UpdateDatabase(url, maxConns)
})
```

## Best Practices

### 1. Validate Before Applying

Always validate new configuration before applying changes:

```go
cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
    cfg := currentCfg.(*AppConfig)
    
    // Validate first
    if err := validateConfiguration(cfg); err != nil {
        return fmt.Errorf("invalid configuration: %w", err)
    }
    
    // Then apply
    return applyConfiguration(cfg)
})
```

### 2. Use Graceful Shutdowns

When restarting services, always use graceful shutdowns:

```go
func gracefulRestart(server *http.Server, newAddr string) error {
    // Graceful shutdown with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := server.Shutdown(ctx); err != nil {
        return fmt.Errorf("graceful shutdown failed: %w", err)
    }
    
    // Start new server
    return startNewServer(newAddr)
}
```

### 3. Log Configuration Changes

Always log configuration changes for debugging and auditing:

```go
cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
    cfg := currentCfg.(*AppConfig)
    
    log.Printf("Configuration change detected:")
    log.Printf("  Timestamp: %s", time.Now().Format(time.RFC3339))
    log.Printf("  Server: %s:%d", cfg.Server.Host, cfg.Server.Port)
    log.Printf("  Debug: %t", cfg.Debug)
    
    return applyConfiguration(cfg)
})
```

### 4. Handle Partial Failures

Design your callbacks to handle partial failures gracefully:

```go
cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
    cfg := currentCfg.(*AppConfig)
    var errors []error
    
    // Try to update each component independently
    if err := updateServerConfig(cfg.Server); err != nil {
        errors = append(errors, fmt.Errorf("server config: %w", err))
    }
    
    if err := updateDatabaseConfig(cfg.Database); err != nil {
        errors = append(errors, fmt.Errorf("database config: %w", err))
    }
    
    if err := updateLoggingConfig(cfg.Logging); err != nil {
        errors = append(errors, fmt.Errorf("logging config: %w", err))
    }
    
    // Return combined errors if any
    if len(errors) > 0 {
        return fmt.Errorf("configuration update errors: %v", errors)
    }
    
    return nil
})
```

## Troubleshooting

### Common Issues

1. **File permissions**: Ensure the application has read access to configuration files
2. **File locks**: Some editors create temporary files that might trigger false reloads
3. **Rapid changes**: Multiple rapid changes might cause callback flooding

### Debugging Hot Reload

Enable debug logging to troubleshoot hot reload issues:

```go
import "github.com/spf13/viper"

// Enable viper debug logging
viper.SetConfigType("yaml")
viper.Debug() // This will enable debug output
```

## Next Steps

- [Best Practices](04_best_practices.md) - Follow recommended patterns
- [Integration Examples](05_integration_examples.md) - See real-world examples
- [Troubleshooting](06_troubleshooting.md) - Common issues and solutions 