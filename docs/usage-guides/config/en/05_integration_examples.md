# Integration Examples

This document provides real-world examples of integrating the configuration module with various frameworks and libraries.

## Web Framework Integration

### Gin Framework Integration

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
    "github.com/spf13/viper"
)

type WebConfig struct {
    Server struct {
        Host         string        `mapstructure:"host" default:"localhost"`
        Port         int           `mapstructure:"port" default:"8080"`
        ReadTimeout  time.Duration `mapstructure:"read_timeout" default:"30s"`
        WriteTimeout time.Duration `mapstructure:"write_timeout" default:"30s"`
        Mode         string        `mapstructure:"mode" default:"debug"`
    } `mapstructure:"server"`
    
    Database struct {
        URL         string `mapstructure:"url" default:"postgres://localhost/myapp"`
        MaxConns    int    `mapstructure:"max_connections" default:"25"`
    } `mapstructure:"database"`
    
    Redis struct {
        URL string `mapstructure:"url" default:"redis://localhost:6379"`
        DB  int    `mapstructure:"db" default:"0"`
    } `mapstructure:"redis"`
}

type WebServer struct {
    config *WebConfig
    server *http.Server
    router *gin.Engine
}

func NewWebServer() *WebServer {
    return &WebServer{
        router: gin.New(),
    }
}

func (ws *WebServer) LoadConfig() error {
    var cfg WebConfig
    
    cm, err := config.LoadConfigAndWatch(
        &cfg,
        config.WithConfigFile("config.yaml", ""),
        config.WithEnvPrefix("APP"),
        config.WithHotReload(true),
        config.WithEnvVarOverride(true),
    )
    if err != nil {
        return fmt.Errorf("failed to load configuration: %w", err)
    }
    
    ws.config = &cfg
    
    // Register hot reload callback
    cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
        newCfg := currentCfg.(*WebConfig)
        return ws.reconfigure(newCfg)
    })
    
    return nil
}

func (ws *WebServer) reconfigure(newCfg *WebConfig) error {
    log.Printf("Reconfiguring web server...")
    
    // Update Gin mode
    gin.SetMode(newCfg.Server.Mode)
    
    // If server configuration changed, restart server
    if ws.config.Server.Host != newCfg.Server.Host || 
       ws.config.Server.Port != newCfg.Server.Port {
        
        log.Printf("Server address changed, restarting...")
        
        // Gracefully shutdown existing server
        if ws.server != nil {
            ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
            defer cancel()
            ws.server.Shutdown(ctx)
        }
        
        // Start new server
        go ws.startServer(newCfg)
    }
    
    ws.config = newCfg
    return nil
}

func (ws *WebServer) startServer(cfg *WebConfig) {
    ws.server = &http.Server{
        Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
        Handler:      ws.router,
        ReadTimeout:  cfg.Server.ReadTimeout,
        WriteTimeout: cfg.Server.WriteTimeout,
    }
    
    log.Printf("Starting server on %s", ws.server.Addr)
    if err := ws.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Printf("Server error: %v", err)
    }
}

func (ws *WebServer) setupRoutes() {
    ws.router.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status": "healthy",
            "config": gin.H{
                "host": ws.config.Server.Host,
                "port": ws.config.Server.Port,
                "mode": ws.config.Server.Mode,
            },
        })
    })
    
    ws.router.GET("/config", func(c *gin.Context) {
        // Return safe configuration (without sensitive data)
        c.JSON(http.StatusOK, gin.H{
            "server": ws.config.Server,
            "database": gin.H{
                "max_connections": ws.config.Database.MaxConns,
            },
        })
    })
}

func main() {
    server := NewWebServer()
    
    if err := server.LoadConfig(); err != nil {
        log.Fatal(err)
    }
    
    server.setupRoutes()
    server.startServer(server.config)
}
```

### Echo Framework Integration

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
)

type EchoConfig struct {
    Server struct {
        Host    string `mapstructure:"host" default:"localhost"`
        Port    int    `mapstructure:"port" default:"8080"`
        Debug   bool   `mapstructure:"debug" default:"true"`
    } `mapstructure:"server"`
    
    Middleware struct {
        CORS struct {
            Enabled bool     `mapstructure:"enabled" default:"true"`
            Origins []string `mapstructure:"origins" default:"*"`
        } `mapstructure:"cors"`
        
        RateLimit struct {
            Enabled bool `mapstructure:"enabled" default:"false"`
            Rate    int  `mapstructure:"rate" default:"100"`
        } `mapstructure:"rate_limit"`
    } `mapstructure:"middleware"`
}

type EchoServer struct {
    config *EchoConfig
    echo   *echo.Echo
    server *http.Server
}

func NewEchoServer() *EchoServer {
    return &EchoServer{
        echo: echo.New(),
    }
}

func (es *EchoServer) LoadConfig() error {
    var cfg EchoConfig
    
    cm, err := config.LoadConfigAndWatch(
        &cfg,
        config.WithConfigFile("echo-config.yaml", ""),
        config.WithEnvPrefix("ECHO"),
        config.WithHotReload(true),
    )
    if err != nil {
        return err
    }
    
    es.config = &cfg
    es.configureEcho()
    
    // Hot reload callback
    cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
        newCfg := currentCfg.(*EchoConfig)
        return es.reconfigure(newCfg)
    })
    
    return nil
}

func (es *EchoServer) configureEcho() {
    es.echo.Debug = es.config.Server.Debug
    
    // Configure middleware based on config
    if es.config.Middleware.CORS.Enabled {
        es.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
            AllowOrigins: es.config.Middleware.CORS.Origins,
        }))
    }
    
    if es.config.Middleware.RateLimit.Enabled {
        es.echo.Use(middleware.RateLimiter(
            middleware.NewRateLimiterMemoryStore(
                rate.Limit(es.config.Middleware.RateLimit.Rate),
            ),
        ))
    }
}

func (es *EchoServer) reconfigure(newCfg *EchoConfig) error {
    log.Printf("Reconfiguring Echo server...")
    
    // Update debug mode
    es.echo.Debug = newCfg.Server.Debug
    
    // If server address changed, restart
    if es.config.Server.Host != newCfg.Server.Host || 
       es.config.Server.Port != newCfg.Server.Port {
        
        if es.server != nil {
            ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
            defer cancel()
            es.server.Shutdown(ctx)
        }
        
        go es.start(newCfg)
    }
    
    es.config = newCfg
    return nil
}

func (es *EchoServer) start(cfg *EchoConfig) {
    addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
    es.server = &http.Server{Addr: addr}
    
    log.Printf("Starting Echo server on %s", addr)
    if err := es.echo.StartServer(es.server); err != nil && err != http.ErrServerClosed {
        log.Printf("Echo server error: %v", err)
    }
}
```

## Database Integration

### GORM Integration

```go
package main

import (
    "fmt"
    "log"
    "time"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
    "github.com/spf13/viper"
)

type DatabaseConfig struct {
    Primary struct {
        Host         string        `mapstructure:"host" default:"localhost"`
        Port         int           `mapstructure:"port" default:"5432"`
        User         string        `mapstructure:"user" default:"postgres"`
        Password     string        `mapstructure:"password"`
        Database     string        `mapstructure:"database" default:"myapp"`
        SSLMode      string        `mapstructure:"ssl_mode" default:"disable"`
        MaxConns     int           `mapstructure:"max_connections" default:"25"`
        MaxIdleConns int           `mapstructure:"max_idle_connections" default:"5"`
        ConnTimeout  time.Duration `mapstructure:"connection_timeout" default:"5s"`
        LogLevel     string        `mapstructure:"log_level" default:"warn"`
    } `mapstructure:"primary"`
    
    Replica struct {
        Enabled      bool          `mapstructure:"enabled" default:"false"`
        Host         string        `mapstructure:"host"`
        Port         int           `mapstructure:"port" default:"5432"`
        User         string        `mapstructure:"user"`
        Password     string        `mapstructure:"password"`
        Database     string        `mapstructure:"database"`
        MaxConns     int           `mapstructure:"max_connections" default:"10"`
        ConnTimeout  time.Duration `mapstructure:"connection_timeout" default:"5s"`
    } `mapstructure:"replica"`
}

type DatabaseManager struct {
    config      *DatabaseConfig
    primaryDB   *gorm.DB
    replicaDB   *gorm.DB
}

func NewDatabaseManager() *DatabaseManager {
    return &DatabaseManager{}
}

func (dm *DatabaseManager) LoadConfig() error {
    var cfg DatabaseConfig
    
    cm, err := config.LoadConfigAndWatch(
        &cfg,
        config.WithConfigFile("database.yaml", ""),
        config.WithEnvPrefix("DB"),
        config.WithHotReload(true),
    )
    if err != nil {
        return err
    }
    
    dm.config = &cfg
    
    // Initialize databases
    if err := dm.initializeDatabases(); err != nil {
        return err
    }
    
    // Register hot reload callback
    cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
        newCfg := currentCfg.(*DatabaseConfig)
        return dm.reconfigure(newCfg)
    })
    
    return nil
}

func (dm *DatabaseManager) initializeDatabases() error {
    // Initialize primary database
    if err := dm.initializePrimary(); err != nil {
        return fmt.Errorf("failed to initialize primary database: %w", err)
    }
    
    // Initialize replica database if enabled
    if dm.config.Replica.Enabled {
        if err := dm.initializeReplica(); err != nil {
            log.Printf("Warning: Failed to initialize replica database: %v", err)
        }
    }
    
    return nil
}

func (dm *DatabaseManager) initializePrimary() error {
    dsn := dm.buildDSN(dm.config.Primary.Host, dm.config.Primary.Port,
        dm.config.Primary.User, dm.config.Primary.Password,
        dm.config.Primary.Database, dm.config.Primary.SSLMode)
    
    logLevel := dm.getLogLevel(dm.config.Primary.LogLevel)
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logLevel),
    })
    if err != nil {
        return err
    }
    
    // Configure connection pool
    sqlDB, err := db.DB()
    if err != nil {
        return err
    }
    
    sqlDB.SetMaxOpenConns(dm.config.Primary.MaxConns)
    sqlDB.SetMaxIdleConns(dm.config.Primary.MaxIdleConns)
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    dm.primaryDB = db
    return nil
}

func (dm *DatabaseManager) initializeReplica() error {
    dsn := dm.buildDSN(dm.config.Replica.Host, dm.config.Replica.Port,
        dm.config.Replica.User, dm.config.Replica.Password,
        dm.config.Replica.Database, "disable")
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    })
    if err != nil {
        return err
    }
    
    sqlDB, err := db.DB()
    if err != nil {
        return err
    }
    
    sqlDB.SetMaxOpenConns(dm.config.Replica.MaxConns)
    sqlDB.SetMaxIdleConns(dm.config.Replica.MaxConns / 2)
    
    dm.replicaDB = db
    return nil
}

func (dm *DatabaseManager) buildDSN(host string, port int, user, password, database, sslmode string) string {
    return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        host, port, user, password, database, sslmode)
}

func (dm *DatabaseManager) getLogLevel(level string) logger.LogLevel {
    switch level {
    case "silent":
        return logger.Silent
    case "error":
        return logger.Error
    case "warn":
        return logger.Warn
    case "info":
        return logger.Info
    default:
        return logger.Warn
    }
}

func (dm *DatabaseManager) reconfigure(newCfg *DatabaseConfig) error {
    log.Printf("Reconfiguring database connections...")
    
    // Close existing connections
    if dm.primaryDB != nil {
        if sqlDB, err := dm.primaryDB.DB(); err == nil {
            sqlDB.Close()
        }
    }
    
    if dm.replicaDB != nil {
        if sqlDB, err := dm.replicaDB.DB(); err == nil {
            sqlDB.Close()
        }
    }
    
    // Update configuration
    dm.config = newCfg
    
    // Reinitialize databases
    return dm.initializeDatabases()
}

func (dm *DatabaseManager) GetPrimary() *gorm.DB {
    return dm.primaryDB
}

func (dm *DatabaseManager) GetReplica() *gorm.DB {
    if dm.replicaDB != nil {
        return dm.replicaDB
    }
    return dm.primaryDB // Fallback to primary
}
```

## Logging Integration

### Integration with pkg/log

```go
package main

import (
    "log"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
    sdklog "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
    "github.com/spf13/viper"
)

type AppConfig struct {
    App struct {
        Name string `mapstructure:"name" default:"my-app"`
    } `mapstructure:"app"`
    
    // Embed log options directly
    Log sdklog.Options `mapstructure:"log"`
    
    Server struct {
        Port int `mapstructure:"port" default:"8080"`
    } `mapstructure:"server"`
}

func main() {
    var cfg AppConfig
    
    // Load configuration with hot reload
    cm, err := config.LoadConfigAndWatch(
        &cfg,
        config.WithConfigFile("app-config.yaml", ""),
        config.WithEnvPrefix("APP"),
        config.WithHotReload(true),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Initialize logging with configuration
    sdklog.Init(&cfg.Log)
    
    // Register log package for hot reload
    sdklog.RegisterConfigHotReload(cm)
    
    // Register application callback
    cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
        appCfg := currentCfg.(*AppConfig)
        sdklog.Infof("Configuration reloaded for app: %s", appCfg.App.Name)
        return nil
    })
    
    // Use the logger
    sdklog.Info("Application started")
    sdklog.Infow("Server configuration", "port", cfg.Server.Port)
    
    // Keep running
    select {}
}
```

## Microservices Integration

### Service Discovery Integration

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
    "github.com/spf13/viper"
)

type ServiceConfig struct {
    Service struct {
        Name    string `mapstructure:"name" default:"my-service"`
        Version string `mapstructure:"version" default:"1.0.0"`
        Port    int    `mapstructure:"port" default:"8080"`
    } `mapstructure:"service"`
    
    Discovery struct {
        Enabled  bool   `mapstructure:"enabled" default:"true"`
        Endpoint string `mapstructure:"endpoint" default:"http://consul:8500"`
        Interval string `mapstructure:"interval" default:"30s"`
    } `mapstructure:"discovery"`
    
    Dependencies []ServiceDependency `mapstructure:"dependencies"`
}

type ServiceDependency struct {
    Name     string `mapstructure:"name"`
    Endpoint string `mapstructure:"endpoint"`
    Timeout  string `mapstructure:"timeout" default:"5s"`
}

type ServiceManager struct {
    config       *ServiceConfig
    dependencies map[string]*ServiceDependency
}

func NewServiceManager() *ServiceManager {
    return &ServiceManager{
        dependencies: make(map[string]*ServiceDependency),
    }
}

func (sm *ServiceManager) LoadConfig() error {
    var cfg ServiceConfig
    
    cm, err := config.LoadConfigAndWatch(
        &cfg,
        config.WithConfigFile("service.yaml", ""),
        config.WithEnvPrefix("SERVICE"),
        config.WithHotReload(true),
    )
    if err != nil {
        return err
    }
    
    sm.config = &cfg
    sm.updateDependencies()
    
    // Register for service discovery updates
    cm.RegisterSectionChangeCallback("dependencies", func(v *viper.Viper) error {
        return sm.updateServiceDependencies(v)
    })
    
    // Register for discovery configuration changes
    cm.RegisterSectionChangeCallback("discovery", func(v *viper.Viper) error {
        return sm.reconfigureDiscovery(v)
    })
    
    return nil
}

func (sm *ServiceManager) updateDependencies() {
    sm.dependencies = make(map[string]*ServiceDependency)
    for i := range sm.config.Dependencies {
        dep := &sm.config.Dependencies[i]
        sm.dependencies[dep.Name] = dep
        log.Printf("Registered dependency: %s -> %s", dep.Name, dep.Endpoint)
    }
}

func (sm *ServiceManager) updateServiceDependencies(v *viper.Viper) error {
    log.Printf("Service dependencies configuration changed")
    
    var newDeps []ServiceDependency
    if err := v.UnmarshalKey("dependencies", &newDeps); err != nil {
        return fmt.Errorf("failed to unmarshal dependencies: %w", err)
    }
    
    // Update dependencies map
    newDepsMap := make(map[string]*ServiceDependency)
    for i := range newDeps {
        dep := &newDeps[i]
        newDepsMap[dep.Name] = dep
    }
    
    // Check for removed dependencies
    for name := range sm.dependencies {
        if _, exists := newDepsMap[name]; !exists {
            log.Printf("Dependency removed: %s", name)
        }
    }
    
    // Check for new dependencies
    for name, dep := range newDepsMap {
        if _, exists := sm.dependencies[name]; !exists {
            log.Printf("New dependency added: %s -> %s", name, dep.Endpoint)
        }
    }
    
    sm.dependencies = newDepsMap
    sm.config.Dependencies = newDeps
    
    return nil
}

func (sm *ServiceManager) reconfigureDiscovery(v *viper.Viper) error {
    enabled := v.GetBool("discovery.enabled")
    endpoint := v.GetString("discovery.endpoint")
    
    log.Printf("Discovery configuration changed: enabled=%t, endpoint=%s", enabled, endpoint)
    
    if enabled {
        return sm.registerWithDiscovery(endpoint)
    } else {
        return sm.deregisterFromDiscovery()
    }
}

func (sm *ServiceManager) registerWithDiscovery(endpoint string) error {
    log.Printf("Registering service with discovery at %s", endpoint)
    // Implementation would register with actual service discovery
    return nil
}

func (sm *ServiceManager) deregisterFromDiscovery() error {
    log.Printf("Deregistering service from discovery")
    // Implementation would deregister from service discovery
    return nil
}

func (sm *ServiceManager) GetDependency(name string) (*ServiceDependency, bool) {
    dep, exists := sm.dependencies[name]
    return dep, exists
}
```

## Configuration File Examples

### Complete Application Configuration

```yaml
# app-config.yaml
app:
  name: "my-microservice"
  version: "1.2.3"
  environment: "production"

server:
  http:
    host: "0.0.0.0"
    port: 8080
    read_timeout: "30s"
    write_timeout: "30s"
  grpc:
    host: "0.0.0.0"
    port: 9090

database:
  primary:
    host: "postgres-primary"
    port: 5432
    user: "myapp"
    database: "myapp_prod"
    max_connections: 50
    max_idle_connections: 10
    connection_timeout: "10s"
    log_level: "warn"
  
  replica:
    enabled: true
    host: "postgres-replica"
    port: 5432
    user: "myapp_readonly"
    database: "myapp_prod"
    max_connections: 20

redis:
  url: "redis://redis-cluster:6379"
  db: 0
  pool_size: 10

# Logging configuration (integrated with pkg/log)
log:
  level: "info"
  format: "json"
  output_paths: ["stdout", "/var/log/app.log"]
  enable_color: false
  disable_caller: false
  disable_stacktrace: false
  stacktrace_level: "error"
  log_rotate_max_size: 100
  log_rotate_max_backups: 5
  log_rotate_max_age: 30
  log_rotate_compress: true

# Service discovery
discovery:
  enabled: true
  endpoint: "http://consul:8500"
  interval: "30s"

# Service dependencies
dependencies:
  - name: "user-service"
    endpoint: "http://user-service:8080"
    timeout: "5s"
  - name: "notification-service"
    endpoint: "http://notification-service:8080"
    timeout: "10s"

# Feature flags
features:
  authentication: true
  rate_limiting: true
  metrics: true
  tracing: true
  caching: true

# Observability
observability:
  metrics:
    enabled: true
    port: 9090
    path: "/metrics"
  
  tracing:
    enabled: true
    endpoint: "http://jaeger:14268/api/traces"
    sample_rate: 0.1
  
  health:
    enabled: true
    port: 8081
    path: "/health"

# Security
security:
  cors:
    enabled: true
    origins: ["https://app.example.com", "https://admin.example.com"]
  
  rate_limit:
    enabled: true
    requests_per_minute: 1000
    burst: 100
```

## Next Steps

- [Troubleshooting](06_troubleshooting.md) - Common issues and solutions
- [Best Practices](04_best_practices.md) - Follow recommended patterns 