# 集成示例

本文档提供了将配置模块与流行框架和库集成的实际示例。

## Web 框架集成

### Gin 框架集成

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
    "github.com/spf13/viper"
)

type AppConfig struct {
    App struct {
        Name        string `mapstructure:"name" default:"gin-app"`
        Environment string `mapstructure:"environment" default:"development"`
    } `mapstructure:"app"`
    
    Server struct {
        Host         string        `mapstructure:"host" default:"localhost"`
        Port         int           `mapstructure:"port" default:"8080"`
        ReadTimeout  time.Duration `mapstructure:"read_timeout" default:"30s"`
        WriteTimeout time.Duration `mapstructure:"write_timeout" default:"30s"`
    } `mapstructure:"server"`
    
    Database struct {
        URL      string `mapstructure:"url" default:"postgres://localhost/ginapp"`
        MaxConns int    `mapstructure:"max_connections" default:"25"`
    } `mapstructure:"database"`
    
    Log log.Options `mapstructure:"log"`
}

type Application struct {
    config *AppConfig
    server *http.Server
    router *gin.Engine
}

func main() {
    app := &Application{}
    
    // 加载配置
    if err := app.loadConfiguration(); err != nil {
        panic(fmt.Sprintf("加载配置失败: %v", err))
    }
    
    // 初始化日志
    log.Init(&app.config.Log)
    
    // 设置 Gin 模式
    if app.config.App.Environment == "production" {
        gin.SetMode(gin.ReleaseMode)
    }
    
    // 初始化路由
    app.setupRoutes()
    
    // 启动服务器
    app.startServer()
}

func (app *Application) loadConfiguration() error {
    var cfg AppConfig
    
    // 加载配置并启用热重载
    cm, err := config.LoadConfigAndWatch(
        &cfg,
        config.WithConfigFile("config.yaml", ""),
        config.WithEnvPrefix("GIN_APP"),
        config.WithHotReload(true),
        config.WithEnvVarOverride(true),
    )
    if err != nil {
        return err
    }
    
    app.config = &cfg
    
    // 注册配置变更回调
    cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
        newCfg := currentCfg.(*AppConfig)
        log.Infow("配置已更新",
            "server_port", newCfg.Server.Port,
            "environment", newCfg.App.Environment,
        )
        
        // 更新日志配置
        log.Init(&newCfg.Log)
        
        app.config = newCfg
        return nil
    })
    
    return nil
}

func (app *Application) setupRoutes() {
    app.router = gin.New()
    
    // 中间件
    app.router.Use(gin.Logger())
    app.router.Use(gin.Recovery())
    app.router.Use(app.configMiddleware())
    
    // 路由
    app.router.GET("/health", app.healthHandler)
    app.router.GET("/config", app.configHandler)
    
    api := app.router.Group("/api/v1")
    {
        api.GET("/users", app.getUsersHandler)
        api.POST("/users", app.createUserHandler)
    }
}

func (app *Application) configMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 将配置添加到上下文
        c.Set("config", app.config)
        c.Next()
    }
}

func (app *Application) healthHandler(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status": "healthy",
        "app":    app.config.App.Name,
        "env":    app.config.App.Environment,
    })
}

func (app *Application) configHandler(c *gin.Context) {
    // 返回非敏感配置信息
    c.JSON(http.StatusOK, gin.H{
        "app": gin.H{
            "name":        app.config.App.Name,
            "environment": app.config.App.Environment,
        },
        "server": gin.H{
            "host": app.config.Server.Host,
            "port": app.config.Server.Port,
        },
    })
}

func (app *Application) getUsersHandler(c *gin.Context) {
    // 从上下文获取配置
    cfg := c.MustGet("config").(*AppConfig)
    
    log.Infow("获取用户列表",
        "database_url", cfg.Database.URL,
        "max_connections", cfg.Database.MaxConns,
    )
    
    c.JSON(http.StatusOK, gin.H{"users": []string{}})
}

func (app *Application) createUserHandler(c *gin.Context) {
    c.JSON(http.StatusCreated, gin.H{"message": "用户已创建"})
}

func (app *Application) startServer() {
    addr := fmt.Sprintf("%s:%d", app.config.Server.Host, app.config.Server.Port)
    
    app.server = &http.Server{
        Addr:         addr,
        Handler:      app.router,
        ReadTimeout:  app.config.Server.ReadTimeout,
        WriteTimeout: app.config.Server.WriteTimeout,
    }
    
    // 在 goroutine 中启动服务器
    go func() {
        log.Infow("启动 HTTP 服务器",
            "address", addr,
            "environment", app.config.App.Environment,
        )
        
        if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Errorw("HTTP 服务器错误", "error", err)
        }
    }()
    
    // 等待中断信号以优雅关闭服务器
    app.gracefulShutdown()
}

func (app *Application) gracefulShutdown() {
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    log.Info("正在关闭服务器...")
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := app.server.Shutdown(ctx); err != nil {
        log.Errorw("服务器强制关闭", "error", err)
    }
    
    log.Info("服务器已退出")
}
```

配置文件 `config.yaml`：

```yaml
app:
  name: "gin-example-app"
  environment: "development"

server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: "30s"
  write_timeout: "30s"

database:
  url: "postgres://user:pass@localhost/ginapp"
  max_connections: 25

log:
  level: "info"
  format: "json"
  output_paths: ["stdout"]
```

### Echo 框架集成

```go
package main

import (
    "context"
    "net/http"
    "os"
    "os/signal"
    "time"
    
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
    "github.com/spf13/viper"
)

type EchoAppConfig struct {
    Server struct {
        Host string `mapstructure:"host" default:"localhost"`
        Port string `mapstructure:"port" default:"8080"`
    } `mapstructure:"server"`
    
    CORS struct {
        AllowOrigins []string `mapstructure:"allow_origins" default:"[\"*\"]"`
        AllowMethods []string `mapstructure:"allow_methods" default:"[\"GET\",\"POST\",\"PUT\",\"DELETE\"]"`
    } `mapstructure:"cors"`
    
    RateLimit struct {
        Enabled bool `mapstructure:"enabled" default:"true"`
        Rate    int  `mapstructure:"rate" default:"100"`
    } `mapstructure:"rate_limit"`
    
    Log log.Options `mapstructure:"log"`
}

type EchoApp struct {
    config *EchoAppConfig
    echo   *echo.Echo
}

func main() {
    app := &EchoApp{}
    
    // 加载配置
    if err := app.loadConfig(); err != nil {
        panic(err)
    }
    
    // 初始化日志
    log.Init(&app.config.Log)
    
    // 设置 Echo
    app.setupEcho()
    
    // 启动服务器
    app.start()
}

func (app *EchoApp) loadConfig() error {
    var cfg EchoAppConfig
    
    cm, err := config.LoadConfigAndWatch(
        &cfg,
        config.WithConfigFile("echo-config.yaml", ""),
        config.WithEnvPrefix("ECHO_APP"),
        config.WithHotReload(true),
    )
    if err != nil {
        return err
    }
    
    app.config = &cfg
    
    // 热重载回调
    cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
        newCfg := currentCfg.(*EchoAppConfig)
        log.Infow("Echo 配置已更新", "port", newCfg.Server.Port)
        app.config = newCfg
        return nil
    })
    
    return nil
}

func (app *EchoApp) setupEcho() {
    app.echo = echo.New()
    
    // 中间件
    app.echo.Use(middleware.Logger())
    app.echo.Use(middleware.Recover())
    
    // CORS 中间件
    app.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: app.config.CORS.AllowOrigins,
        AllowMethods: app.config.CORS.AllowMethods,
    }))
    
    // 速率限制中间件
    if app.config.RateLimit.Enabled {
        app.echo.Use(middleware.RateLimiter(
            middleware.NewRateLimiterMemoryStore(
                float64(app.config.RateLimit.Rate),
            ),
        ))
    }
    
    // 路由
    app.echo.GET("/", app.homeHandler)
    app.echo.GET("/health", app.healthHandler)
}

func (app *EchoApp) homeHandler(c echo.Context) error {
    return c.JSON(http.StatusOK, map[string]string{
        "message": "欢迎使用 Echo 应用",
        "version": "1.0.0",
    })
}

func (app *EchoApp) healthHandler(c echo.Context) error {
    return c.JSON(http.StatusOK, map[string]interface{}{
        "status":     "healthy",
        "rate_limit": app.config.RateLimit.Enabled,
        "cors":       len(app.config.CORS.AllowOrigins) > 0,
    })
}

func (app *EchoApp) start() {
    // 启动服务器
    go func() {
        addr := app.config.Server.Host + ":" + app.config.Server.Port
        log.Infow("启动 Echo 服务器", "address", addr)
        
        if err := app.echo.Start(addr); err != nil && err != http.ErrServerClosed {
            log.Errorw("Echo 服务器错误", "error", err)
        }
    }()
    
    // 优雅关闭
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt)
    <-quit
    
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    if err := app.echo.Shutdown(ctx); err != nil {
        log.Errorw("Echo 服务器关闭错误", "error", err)
    }
}
```

## 数据库集成

### GORM 集成

```go
package main

import (
    "fmt"
    "time"
    
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
    "github.com/spf13/viper"
)

type DatabaseConfig struct {
    Database struct {
        Host     string `mapstructure:"host" default:"localhost"`
        Port     int    `mapstructure:"port" default:"5432"`
        User     string `mapstructure:"user" default:"postgres"`
        Password string `mapstructure:"password"`
        DBName   string `mapstructure:"dbname" default:"myapp"`
        SSLMode  string `mapstructure:"sslmode" default:"disable"`
        
        Pool struct {
            MaxOpenConns    int           `mapstructure:"max_open_conns" default:"25"`
            MaxIdleConns    int           `mapstructure:"max_idle_conns" default:"5"`
            ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" default:"1h"`
            ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time" default:"30m"`
        } `mapstructure:"pool"`
        
        Log struct {
            Level string `mapstructure:"level" default:"warn"`
        } `mapstructure:"log"`
    } `mapstructure:"database"`
    
    Log log.Options `mapstructure:"log"`
}

type User struct {
    ID        uint      `gorm:"primarykey"`
    Name      string    `gorm:"size:100;not null"`
    Email     string    `gorm:"size:100;uniqueIndex;not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

type DatabaseManager struct {
    config *DatabaseConfig
    db     *gorm.DB
}

func main() {
    dm := &DatabaseManager{}
    
    // 加载配置
    if err := dm.loadConfig(); err != nil {
        panic(err)
    }
    
    // 初始化日志
    log.Init(&dm.config.Log)
    
    // 连接数据库
    if err := dm.connect(); err != nil {
        panic(err)
    }
    
    // 迁移数据库
    if err := dm.migrate(); err != nil {
        panic(err)
    }
    
    // 示例操作
    dm.exampleOperations()
    
    log.Info("数据库示例完成")
}

func (dm *DatabaseManager) loadConfig() error {
    var cfg DatabaseConfig
    
    cm, err := config.LoadConfigAndWatch(
        &cfg,
        config.WithConfigFile("db-config.yaml", ""),
        config.WithEnvPrefix("DB_APP"),
        config.WithHotReload(true),
    )
    if err != nil {
        return err
    }
    
    dm.config = &cfg
    
    // 注册数据库配置变更回调
    cm.RegisterSectionChangeCallback("database", func(v *viper.Viper) error {
        log.Info("数据库配置已更改，重新连接...")
        return dm.reconnect()
    })
    
    return nil
}

func (dm *DatabaseManager) connect() error {
    dsn := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        dm.config.Database.Host,
        dm.config.Database.Port,
        dm.config.Database.User,
        dm.config.Database.Password,
        dm.config.Database.DBName,
        dm.config.Database.SSLMode,
    )
    
    // 配置 GORM 日志级别
    var logLevel logger.LogLevel
    switch dm.config.Database.Log.Level {
    case "silent":
        logLevel = logger.Silent
    case "error":
        logLevel = logger.Error
    case "warn":
        logLevel = logger.Warn
    case "info":
        logLevel = logger.Info
    default:
        logLevel = logger.Warn
    }
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logLevel),
    })
    if err != nil {
        return fmt.Errorf("连接数据库失败: %w", err)
    }
    
    // 获取底层 sql.DB 以配置连接池
    sqlDB, err := db.DB()
    if err != nil {
        return fmt.Errorf("获取 sql.DB 失败: %w", err)
    }
    
    // 配置连接池
    sqlDB.SetMaxOpenConns(dm.config.Database.Pool.MaxOpenConns)
    sqlDB.SetMaxIdleConns(dm.config.Database.Pool.MaxIdleConns)
    sqlDB.SetConnMaxLifetime(dm.config.Database.Pool.ConnMaxLifetime)
    sqlDB.SetConnMaxIdleTime(dm.config.Database.Pool.ConnMaxIdleTime)
    
    dm.db = db
    
    log.Infow("数据库连接成功",
        "host", dm.config.Database.Host,
        "port", dm.config.Database.Port,
        "database", dm.config.Database.DBName,
        "max_open_conns", dm.config.Database.Pool.MaxOpenConns,
    )
    
    return nil
}

func (dm *DatabaseManager) reconnect() error {
    // 关闭现有连接
    if dm.db != nil {
        sqlDB, err := dm.db.DB()
        if err == nil {
            sqlDB.Close()
        }
    }
    
    // 重新连接
    return dm.connect()
}

func (dm *DatabaseManager) migrate() error {
    return dm.db.AutoMigrate(&User{})
}

func (dm *DatabaseManager) exampleOperations() {
    // 创建用户
    user := User{
        Name:  "张三",
        Email: "zhangsan@example.com",
    }
    
    result := dm.db.Create(&user)
    if result.Error != nil {
        log.Errorw("创建用户失败", "error", result.Error)
        return
    }
    
    log.Infow("用户创建成功", "user_id", user.ID, "name", user.Name)
    
    // 查询用户
    var foundUser User
    dm.db.First(&foundUser, user.ID)
    
    log.Infow("查询用户成功",
        "user_id", foundUser.ID,
        "name", foundUser.Name,
        "email", foundUser.Email,
    )
    
    // 更新用户
    dm.db.Model(&foundUser).Update("Name", "李四")
    
    log.Infow("用户更新成功", "user_id", foundUser.ID, "new_name", "李四")
}
```

### Redis 集成

```go
package main

import (
    "context"
    "time"
    
    "github.com/redis/go-redis/v9"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
    "github.com/spf13/viper"
)

type RedisConfig struct {
    Redis struct {
        Host     string `mapstructure:"host" default:"localhost"`
        Port     string `mapstructure:"port" default:"6379"`
        Password string `mapstructure:"password" default:""`
        DB       int    `mapstructure:"db" default:"0"`
        
        Pool struct {
            PoolSize     int           `mapstructure:"pool_size" default:"10"`
            MinIdleConns int           `mapstructure:"min_idle_conns" default:"5"`
            DialTimeout  time.Duration `mapstructure:"dial_timeout" default:"5s"`
            ReadTimeout  time.Duration `mapstructure:"read_timeout" default:"3s"`
            WriteTimeout time.Duration `mapstructure:"write_timeout" default:"3s"`
        } `mapstructure:"pool"`
    } `mapstructure:"redis"`
    
    Log log.Options `mapstructure:"log"`
}

type RedisManager struct {
    config *RedisConfig
    client *redis.Client
}

func main() {
    rm := &RedisManager{}
    
    // 加载配置
    if err := rm.loadConfig(); err != nil {
        panic(err)
    }
    
    // 初始化日志
    log.Init(&rm.config.Log)
    
    // 连接 Redis
    if err := rm.connect(); err != nil {
        panic(err)
    }
    
    // 示例操作
    rm.exampleOperations()
    
    log.Info("Redis 示例完成")
}

func (rm *RedisManager) loadConfig() error {
    var cfg RedisConfig
    
    cm, err := config.LoadConfigAndWatch(
        &cfg,
        config.WithConfigFile("redis-config.yaml", ""),
        config.WithEnvPrefix("REDIS_APP"),
        config.WithHotReload(true),
    )
    if err != nil {
        return err
    }
    
    rm.config = &cfg
    
    // Redis 配置变更回调
    cm.RegisterSectionChangeCallback("redis", func(v *viper.Viper) error {
        log.Info("Redis 配置已更改，重新连接...")
        return rm.reconnect()
    })
    
    return nil
}

func (rm *RedisManager) connect() error {
    addr := rm.config.Redis.Host + ":" + rm.config.Redis.Port
    
    rm.client = redis.NewClient(&redis.Options{
        Addr:         addr,
        Password:     rm.config.Redis.Password,
        DB:           rm.config.Redis.DB,
        PoolSize:     rm.config.Redis.Pool.PoolSize,
        MinIdleConns: rm.config.Redis.Pool.MinIdleConns,
        DialTimeout:  rm.config.Redis.Pool.DialTimeout,
        ReadTimeout:  rm.config.Redis.Pool.ReadTimeout,
        WriteTimeout: rm.config.Redis.Pool.WriteTimeout,
    })
    
    // 测试连接
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    _, err := rm.client.Ping(ctx).Result()
    if err != nil {
        return fmt.Errorf("Redis 连接失败: %w", err)
    }
    
    log.Infow("Redis 连接成功",
        "address", addr,
        "database", rm.config.Redis.DB,
        "pool_size", rm.config.Redis.Pool.PoolSize,
    )
    
    return nil
}

func (rm *RedisManager) reconnect() error {
    // 关闭现有连接
    if rm.client != nil {
        rm.client.Close()
    }
    
    // 重新连接
    return rm.connect()
}

func (rm *RedisManager) exampleOperations() {
    ctx := context.Background()
    
    // 设置键值
    err := rm.client.Set(ctx, "user:1:name", "张三", time.Hour).Err()
    if err != nil {
        log.Errorw("设置键值失败", "error", err)
        return
    }
    
    log.Info("键值设置成功")
    
    // 获取值
    val, err := rm.client.Get(ctx, "user:1:name").Result()
    if err != nil {
        log.Errorw("获取键值失败", "error", err)
        return
    }
    
    log.Infow("获取键值成功", "key", "user:1:name", "value", val)
    
    // 哈希操作
    err = rm.client.HMSet(ctx, "user:1", map[string]interface{}{
        "name":  "张三",
        "email": "zhangsan@example.com",
        "age":   30,
    }).Err()
    if err != nil {
        log.Errorw("设置哈希失败", "error", err)
        return
    }
    
    // 获取哈希
    userInfo := rm.client.HGetAll(ctx, "user:1")
    if userInfo.Err() != nil {
        log.Errorw("获取哈希失败", "error", userInfo.Err())
        return
    }
    
    log.Infow("获取用户信息成功", "user_info", userInfo.Val())
}
```

## 微服务集成

### gRPC 服务集成

```go
package main

import (
    "context"
    "fmt"
    "net"
    "os"
    "os/signal"
    "syscall"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
    "github.com/spf13/viper"
)

type GRPCConfig struct {
    Server struct {
        Host string `mapstructure:"host" default:"localhost"`
        Port int    `mapstructure:"port" default:"9090"`
    } `mapstructure:"server"`
    
    Features struct {
        Reflection bool `mapstructure:"reflection" default:"true"`
        Recovery   bool `mapstructure:"recovery" default:"true"`
    } `mapstructure:"features"`
    
    Log log.Options `mapstructure:"log"`
}

// 示例服务定义（通常从 .proto 文件生成）
type UserService struct {
    config *GRPCConfig
}

func (s *UserService) GetUser(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
    log.InfoContext(ctx, "获取用户请求", "user_id", req.UserId)
    
    return &GetUserResponse{
        User: &User{
            Id:    req.UserId,
            Name:  "张三",
            Email: "zhangsan@example.com",
        },
    }, nil
}

// 简化的消息类型（通常从 .proto 文件生成）
type GetUserRequest struct {
    UserId string
}

type GetUserResponse struct {
    User *User
}

type User struct {
    Id    string
    Name  string
    Email string
}

type GRPCServer struct {
    config     *GRPCConfig
    server     *grpc.Server
    userService *UserService
}

func main() {
    srv := &GRPCServer{}
    
    // 加载配置
    if err := srv.loadConfig(); err != nil {
        panic(err)
    }
    
    // 初始化日志
    log.Init(&srv.config.Log)
    
    // 设置 gRPC 服务器
    srv.setupServer()
    
    // 启动服务器
    srv.start()
}

func (srv *GRPCServer) loadConfig() error {
    var cfg GRPCConfig
    
    cm, err := config.LoadConfigAndWatch(
        &cfg,
        config.WithConfigFile("grpc-config.yaml", ""),
        config.WithEnvPrefix("GRPC_APP"),
        config.WithHotReload(true),
    )
    if err != nil {
        return err
    }
    
    srv.config = &cfg
    
    // 配置变更回调
    cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
        newCfg := currentCfg.(*GRPCConfig)
        log.Infow("gRPC 配置已更新",
            "port", newCfg.Server.Port,
            "reflection", newCfg.Features.Reflection,
        )
        srv.config = newCfg
        return nil
    })
    
    return nil
}

func (srv *GRPCServer) setupServer() {
    var opts []grpc.ServerOption
    
    // 添加中间件
    if srv.config.Features.Recovery {
        // 添加恢复中间件（实际实现需要导入相应包）
        log.Info("启用 gRPC 恢复中间件")
    }
    
    srv.server = grpc.NewServer(opts...)
    
    // 注册服务
    srv.userService = &UserService{config: srv.config}
    // RegisterUserServiceServer(srv.server, srv.userService) // 实际注册
    
    // 启用反射（用于调试）
    if srv.config.Features.Reflection {
        reflection.Register(srv.server)
        log.Info("gRPC 反射已启用")
    }
}

func (srv *GRPCServer) start() {
    addr := fmt.Sprintf("%s:%d", srv.config.Server.Host, srv.config.Server.Port)
    
    lis, err := net.Listen("tcp", addr)
    if err != nil {
        log.Errorw("监听失败", "error", err, "address", addr)
        return
    }
    
    // 在 goroutine 中启动服务器
    go func() {
        log.Infow("启动 gRPC 服务器", "address", addr)
        
        if err := srv.server.Serve(lis); err != nil {
            log.Errorw("gRPC 服务器错误", "error", err)
        }
    }()
    
    // 优雅关闭
    srv.gracefulShutdown()
}

func (srv *GRPCServer) gracefulShutdown() {
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    log.Info("正在关闭 gRPC 服务器...")
    srv.server.GracefulStop()
    log.Info("gRPC 服务器已关闭")
}
```

## 监控和可观测性集成

### Prometheus 指标集成

```go
package main

import (
    "net/http"
    
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
    "github.com/spf13/viper"
)

type MetricsConfig struct {
    Metrics struct {
        Enabled bool   `mapstructure:"enabled" default:"true"`
        Host    string `mapstructure:"host" default:"localhost"`
        Port    int    `mapstructure:"port" default:"9090"`
        Path    string `mapstructure:"path" default:"/metrics"`
    } `mapstructure:"metrics"`
    
    App struct {
        Name    string `mapstructure:"name" default:"my-app"`
        Version string `mapstructure:"version" default:"1.0.0"`
    } `mapstructure:"app"`
    
    Log log.Options `mapstructure:"log"`
}

type MetricsManager struct {
    config *MetricsConfig
    
    // Prometheus 指标
    requestsTotal    *prometheus.CounterVec
    requestDuration  *prometheus.HistogramVec
    activeConnections prometheus.Gauge
}

func main() {
    mm := &MetricsManager{}
    
    // 加载配置
    if err := mm.loadConfig(); err != nil {
        panic(err)
    }
    
    // 初始化日志
    log.Init(&mm.config.Log)
    
    // 设置指标
    mm.setupMetrics()
    
    // 启动指标服务器
    if mm.config.Metrics.Enabled {
        mm.startMetricsServer()
    }
    
    // 模拟应用程序工作
    mm.simulateWork()
}

func (mm *MetricsManager) loadConfig() error {
    var cfg MetricsConfig
    
    cm, err := config.LoadConfigAndWatch(
        &cfg,
        config.WithConfigFile("metrics-config.yaml", ""),
        config.WithEnvPrefix("METRICS_APP"),
        config.WithHotReload(true),
    )
    if err != nil {
        return err
    }
    
    mm.config = &cfg
    
    // 指标配置变更回调
    cm.RegisterSectionChangeCallback("metrics", func(v *viper.Viper) error {
        enabled := v.GetBool("metrics.enabled")
        log.Infow("指标配置已更新", "enabled", enabled)
        
        if enabled && !mm.config.Metrics.Enabled {
            // 启用指标
            go mm.startMetricsServer()
        }
        
        mm.config.Metrics.Enabled = enabled
        return nil
    })
    
    return nil
}

func (mm *MetricsManager) setupMetrics() {
    // 请求总数计数器
    mm.requestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "HTTP 请求总数",
            ConstLabels: prometheus.Labels{
                "app":     mm.config.App.Name,
                "version": mm.config.App.Version,
            },
        },
        []string{"method", "endpoint", "status"},
    )
    
    // 请求持续时间直方图
    mm.requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP 请求持续时间",
            Buckets: prometheus.DefBuckets,
            ConstLabels: prometheus.Labels{
                "app":     mm.config.App.Name,
                "version": mm.config.App.Version,
            },
        },
        []string{"method", "endpoint"},
    )
    
    // 活跃连接数量表
    mm.activeConnections = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_connections",
            Help: "当前活跃连接数",
            ConstLabels: prometheus.Labels{
                "app":     mm.config.App.Name,
                "version": mm.config.App.Version,
            },
        },
    )
    
    // 注册指标
    prometheus.MustRegister(mm.requestsTotal)
    prometheus.MustRegister(mm.requestDuration)
    prometheus.MustRegister(mm.activeConnections)
    
    log.Info("Prometheus 指标已设置")
}

func (mm *MetricsManager) startMetricsServer() {
    addr := fmt.Sprintf("%s:%d", mm.config.Metrics.Host, mm.config.Metrics.Port)
    
    http.Handle(mm.config.Metrics.Path, promhttp.Handler())
    
    log.Infow("启动指标服务器",
        "address", addr,
        "path", mm.config.Metrics.Path,
    )
    
    if err := http.ListenAndServe(addr, nil); err != nil {
        log.Errorw("指标服务器错误", "error", err)
    }
}

func (mm *MetricsManager) simulateWork() {
    // 模拟一些指标更新
    timer := prometheus.NewTimer(mm.requestDuration.WithLabelValues("GET", "/api/users"))
    defer timer.ObserveDuration()
    
    mm.requestsTotal.WithLabelValues("GET", "/api/users", "200").Inc()
    mm.activeConnections.Set(42)
    
    log.Info("指标已更新")
    
    // 保持程序运行
    select {}
}
```

配置文件示例：

```yaml
# metrics-config.yaml
metrics:
  enabled: true
  host: "0.0.0.0"
  port: 9090
  path: "/metrics"

app:
  name: "metrics-example"
  version: "1.0.0"

log:
  level: "info"
  format: "json"
```

这些集成示例展示了如何在实际应用程序中使用配置模块，包括：

1. **Web 框架集成**：Gin 和 Echo 框架的完整集成
2. **数据库集成**：GORM 和 Redis 的配置管理
3. **微服务集成**：gRPC 服务的配置
4. **监控集成**：Prometheus 指标的配置

每个示例都包含：
- 完整的配置结构定义
- 热重载支持
- 错误处理
- 日志集成
- 生产就绪的模式

## 下一步

- [故障排除](06_troubleshooting.md) - 常见问题和解决方案 