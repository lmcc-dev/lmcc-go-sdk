# 集成示例

本指南提供了真实世界的集成示例，展示如何在各种场景中使用 lmcc-go-sdk 服务器模块。

## 基本 REST API

### 简单 CRUD 操作

```go
package main

import (
    "context"
    "log"
    "strconv"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Email string `json:"email"`
}

var users = make(map[int]*User)
var nextID = 1

func main() {
    config := server.DefaultServerConfig()
    config.Framework = "gin"
    config.Port = 8080
    
    manager, err := server.CreateServerManager("gin", config)
    if err != nil {
        log.Fatal("Failed to create server:", err)
    }
    
    framework := manager.GetFramework()
    
    // 设置 CRUD 路由 (Setup CRUD routes)
    api := framework.Group("/api/v1")
    api.RegisterRoute("GET", "/users", server.HandlerFunc(getUsers))
    api.RegisterRoute("POST", "/users", server.HandlerFunc(createUser))
    api.RegisterRoute("GET", "/users/:id", server.HandlerFunc(getUser))
    api.RegisterRoute("PUT", "/users/:id", server.HandlerFunc(updateUser))
    api.RegisterRoute("DELETE", "/users/:id", server.HandlerFunc(deleteUser))
    
    if err := manager.Start(context.Background()); err != nil {
        log.Fatal("Server failed:", err)
    }
}

func getUsers(ctx server.Context) error {
    userList := make([]*User, 0, len(users))
    for _, user := range users {
        userList = append(userList, user)
    }
    return ctx.JSON(200, userList)
}

func createUser(ctx server.Context) error {
    var user User
    if err := ctx.Bind(&user); err != nil {
        return ctx.JSON(400, map[string]string{"error": "Invalid JSON"})
    }
    
    user.ID = nextID
    nextID++
    users[user.ID] = &user
    
    return ctx.JSON(201, user)
}

func getUser(ctx server.Context) error {
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        return ctx.JSON(400, map[string]string{"error": "Invalid ID"})
    }
    
    user, exists := users[id]
    if !exists {
        return ctx.JSON(404, map[string]string{"error": "User not found"})
    }
    
    return ctx.JSON(200, user)
}

func updateUser(ctx server.Context) error {
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        return ctx.JSON(400, map[string]string{"error": "Invalid ID"})
    }
    
    user, exists := users[id]
    if !exists {
        return ctx.JSON(404, map[string]string{"error": "User not found"})
    }
    
    var updatedUser User
    if err := ctx.Bind(&updatedUser); err != nil {
        return ctx.JSON(400, map[string]string{"error": "Invalid JSON"})
    }
    
    user.Name = updatedUser.Name
    user.Email = updatedUser.Email
    
    return ctx.JSON(200, user)
}

func deleteUser(ctx server.Context) error {
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        return ctx.JSON(400, map[string]string{"error": "Invalid ID"})
    }
    
    _, exists := users[id]
    if !exists {
        return ctx.JSON(404, map[string]string{"error": "User not found"})
    }
    
    delete(users, id)
    return ctx.JSON(204, nil)
}
```

## 数据库集成

### 使用 GORM 和 PostgreSQL

```go
package main

import (
    "context"
    "log"
    "os"
    
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin"
)

type User struct {
    ID    uint   `json:"id" gorm:"primaryKey"`
    Name  string `json:"name" gorm:"not null"`
    Email string `json:"email" gorm:"unique;not null"`
}

type UserService struct {
    db *gorm.DB
}

func main() {
    // 连接数据库 (Connect to database)
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        dsn = "host=localhost user=postgres password=postgres dbname=testdb port=5432 sslmode=disable"
    }
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect database:", err)
    }
    
    // 自动迁移模式 (Auto-migrate schema)
    db.AutoMigrate(&User{})
    
    userService := &UserService{db: db}
    
    // 设置服务器 (Setup server)
    config := server.DefaultServerConfig()
    config.Framework = "gin"
    config.Port = 8080
    
    manager, err := server.CreateServerManager("gin", config)
    if err != nil {
        log.Fatal("Failed to create server:", err)
    }
    
    framework := manager.GetFramework()
    
    // 设置路由 (Setup routes)
    api := framework.Group("/api/v1")
    api.RegisterRoute("GET", "/users", server.HandlerFunc(userService.GetUsers))
    api.RegisterRoute("POST", "/users", server.HandlerFunc(userService.CreateUser))
    api.RegisterRoute("GET", "/users/:id", server.HandlerFunc(userService.GetUser))
    
    if err := manager.Start(context.Background()); err != nil {
        log.Fatal("Server failed:", err)
    }
}

func (us *UserService) GetUsers(ctx server.Context) error {
    var users []User
    result := us.db.Find(&users)
    if result.Error != nil {
        return ctx.JSON(500, map[string]string{"error": "Database error"})
    }
    
    return ctx.JSON(200, users)
}

func (us *UserService) CreateUser(ctx server.Context) error {
    var user User
    if err := ctx.Bind(&user); err != nil {
        return ctx.JSON(400, map[string]string{"error": "Invalid JSON"})
    }
    
    result := us.db.Create(&user)
    if result.Error != nil {
        return ctx.JSON(500, map[string]string{"error": "Failed to create user"})
    }
    
    return ctx.JSON(201, user)
}

func (us *UserService) GetUser(ctx server.Context) error {
    id := ctx.Param("id")
    
    var user User
    result := us.db.First(&user, id)
    if result.Error != nil {
        return ctx.JSON(404, map[string]string{"error": "User not found"})
    }
    
    return ctx.JSON(200, user)
}
```

## JWT 身份验证

### JWT 中间件集成

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin"
)

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Password string `json:"-"`
    Role     string `json:"role"`
}

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

var (
    jwtSecret = []byte("your-secret-key")
    users     = map[string]*User{
        "admin": {ID: 1, Username: "admin", Password: hashPassword("admin123"), Role: "admin"},
        "user":  {ID: 2, Username: "user", Password: hashPassword("user123"), Role: "user"},
    }
)

func main() {
    config := server.DefaultServerConfig()
    config.Framework = "gin"
    config.Port = 8080
    
    manager, err := server.CreateServerManager("gin", config)
    if err != nil {
        log.Fatal("Failed to create server:", err)
    }
    
    framework := manager.GetFramework()
    
    // 公共路由 (Public routes)
    framework.RegisterRoute("POST", "/login", server.HandlerFunc(login))
    
    // 受保护的路由 (Protected routes)
    api := framework.Group("/api")
    api.RegisterMiddleware(JWTMiddleware())
    api.RegisterRoute("GET", "/profile", server.HandlerFunc(getProfile))
    
    if err := manager.Start(context.Background()); err != nil {
        log.Fatal("Server failed:", err)
    }
}

func login(ctx server.Context) error {
    var req LoginRequest
    if err := ctx.Bind(&req); err != nil {
        return ctx.JSON(400, map[string]string{"error": "Invalid JSON"})
    }
    
    user, exists := users[req.Username]
    if !exists || !checkPassword(req.Password, user.Password) {
        return ctx.JSON(401, map[string]string{"error": "Invalid credentials"})
    }
    
    token, err := generateJWT(user)
    if err != nil {
        return ctx.JSON(500, map[string]string{"error": "Failed to generate token"})
    }
    
    return ctx.JSON(200, map[string]interface{}{
        "token": token,
        "user":  user,
    })
}

func getProfile(ctx server.Context) error {
    user := ctx.Get("user").(*User)
    return ctx.JSON(200, user)
}

func JWTMiddleware() server.Middleware {
    return server.MiddlewareFunc(func(ctx server.Context, next func() error) error {
        authHeader := ctx.Header("Authorization")
        if authHeader == "" {
            return ctx.JSON(401, map[string]string{"error": "Authorization header required"})
        }
        
        tokenString := authHeader[len("Bearer "):]
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })
        
        if err != nil || !token.Valid {
            return ctx.JSON(401, map[string]string{"error": "Invalid token"})
        }
        
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            return ctx.JSON(401, map[string]string{"error": "Invalid token claims"})
        }
        
        username := claims["username"].(string)
        user, exists := users[username]
        if !exists {
            return ctx.JSON(401, map[string]string{"error": "User not found"})
        }
        
        ctx.Set("user", user)
        return next()
    })
}

func generateJWT(user *User) (string, error) {
    claims := jwt.MapClaims{
        "username": user.Username,
        "role":     user.Role,
        "exp":      time.Now().Add(24 * time.Hour).Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

func hashPassword(password string) string {
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(hash)
}

func checkPassword(password, hash string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
```

## WebSocket 聊天服务器

### 实时通信

```go
package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "sync"
    
    "github.com/gorilla/websocket"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin"
)

type ChatServer struct {
    clients    map[*websocket.Conn]*Client
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    mu         sync.RWMutex
}

type Client struct {
    conn     *websocket.Conn
    username string
    room     string
}

type Message struct {
    Type     string `json:"type"`
    Username string `json:"username"`
    Room     string `json:"room"`
    Content  string `json:"content"`
}

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // 允许所有来源 (Allow all origins)
    },
}

func NewChatServer() *ChatServer {
    return &ChatServer{
        clients:    make(map[*websocket.Conn]*Client),
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (cs *ChatServer) Run() {
    for {
        select {
        case client := <-cs.register:
            cs.mu.Lock()
            cs.clients[client.conn] = client
            cs.mu.Unlock()
            
        case client := <-cs.unregister:
            cs.mu.Lock()
            if _, ok := cs.clients[client.conn]; ok {
                delete(cs.clients, client.conn)
                client.conn.Close()
            }
            cs.mu.Unlock()
            
        case message := <-cs.broadcast:
            var msg Message
            json.Unmarshal(message, &msg)
            
            cs.mu.RLock()
            for conn, client := range cs.clients {
                if client.room == msg.Room {
                    if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
                        conn.Close()
                        delete(cs.clients, conn)
                    }
                }
            }
            cs.mu.RUnlock()
        }
    }
}

func main() {
    chatServer := NewChatServer()
    go chatServer.Run()
    
    config := server.DefaultServerConfig()
    config.Framework = "gin"
    config.Port = 8080
    
    manager, err := server.CreateServerManager("gin", config)
    if err != nil {
        log.Fatal("Failed to create server:", err)
    }
    
    framework := manager.GetFramework()
    
    // WebSocket 端点 (WebSocket endpoint)
    framework.RegisterRoute("GET", "/ws", server.HandlerFunc(func(ctx server.Context) error {
        conn, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
        if err != nil {
            return err
        }
        
        username := ctx.Query("username")
        room := ctx.Query("room")
        if username == "" || room == "" {
            conn.Close()
            return ctx.JSON(400, map[string]string{"error": "Username and room required"})
        }
        
        client := &Client{
            conn:     conn,
            username: username,
            room:     room,
        }
        
        chatServer.register <- client
        
        // 处理传入消息 (Handle incoming messages)
        go func() {
            defer func() {
                chatServer.unregister <- client
            }()
            
            for {
                var message Message
                err := conn.ReadJSON(&message)
                if err != nil {
                    break
                }
                
                message.Username = client.username
                message.Room = client.room
                
                data, _ := json.Marshal(message)
                chatServer.broadcast <- data
            }
        }()
        
        return nil
    }))
    
    if err := manager.Start(context.Background()); err != nil {
        log.Fatal("Server failed:", err)
    }
}
```

## 微服务集成

### 服务发现和负载均衡

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin"
)

type Service struct {
    Name     string    `json:"name"`
    Host     string    `json:"host"`
    Port     int       `json:"port"`
    Health   string    `json:"health"`
    LastSeen time.Time `json:"last_seen"`
}

type ServiceRegistry struct {
    services map[string][]*Service
    mu       sync.RWMutex
}

func NewServiceRegistry() *ServiceRegistry {
    return &ServiceRegistry{
        services: make(map[string][]*Service),
    }
}

func (sr *ServiceRegistry) Register(service *Service) {
    sr.mu.Lock()
    defer sr.mu.Unlock()
    
    if sr.services[service.Name] == nil {
        sr.services[service.Name] = make([]*Service, 0)
    }
    
    // 检查服务是否已存在 (Check if service already exists)
    for i, s := range sr.services[service.Name] {
        if s.Host == service.Host && s.Port == service.Port {
            sr.services[service.Name][i] = service
            return
        }
    }
    
    sr.services[service.Name] = append(sr.services[service.Name], service)
}

func (sr *ServiceRegistry) Discover(serviceName string) []*Service {
    sr.mu.RLock()
    defer sr.mu.RUnlock()
    
    services := sr.services[serviceName]
    healthy := make([]*Service, 0)
    
    for _, service := range services {
        if service.Health == "healthy" && time.Since(service.LastSeen) < 30*time.Second {
            healthy = append(healthy, service)
        }
    }
    
    return healthy
}

func (sr *ServiceRegistry) GetHealthyService(serviceName string) *Service {
    services := sr.Discover(serviceName)
    if len(services) == 0 {
        return nil
    }
    
    // 简单轮询 (Simple round-robin)
    return services[time.Now().Unix()%int64(len(services))]
}

func main() {
    registry := NewServiceRegistry()
    
    config := server.DefaultServerConfig()
    config.Framework = "gin"
    config.Port = 8080
    
    manager, err := server.CreateServerManager("gin", config)
    if err != nil {
        log.Fatal("Failed to create server:", err)
    }
    
    framework := manager.GetFramework()
    
    // 服务注册端点 (Service registration endpoint)
    framework.RegisterRoute("POST", "/register", server.HandlerFunc(func(ctx server.Context) error {
        var service Service
        if err := ctx.Bind(&service); err != nil {
            return ctx.JSON(400, map[string]string{"error": "Invalid JSON"})
        }
        
        service.LastSeen = time.Now()
        registry.Register(&service)
        
        return ctx.JSON(200, map[string]string{"status": "registered"})
    }))
    
    // 服务发现端点 (Service discovery endpoint)
    framework.RegisterRoute("GET", "/discover/:service", server.HandlerFunc(func(ctx server.Context) error {
        serviceName := ctx.Param("service")
        services := registry.Discover(serviceName)
        
        return ctx.JSON(200, services)
    }))
    
    // 健康检查代理 (Health check proxy)
    framework.RegisterRoute("GET", "/health/:service", server.HandlerFunc(func(ctx server.Context) error {
        serviceName := ctx.Param("service")
        service := registry.GetHealthyService(serviceName)
        
        if service == nil {
            return ctx.JSON(503, map[string]string{"error": "Service unavailable"})
        }
        
        // 代理健康检查到实际服务 (Proxy health check to actual service)
        healthURL := fmt.Sprintf("http://%s:%d/health", service.Host, service.Port)
        resp, err := http.Get(healthURL)
        if err != nil {
            service.Health = "unhealthy"
            return ctx.JSON(503, map[string]string{"error": "Service unhealthy"})
        }
        defer resp.Body.Close()
        
        service.Health = "healthy"
        service.LastSeen = time.Now()
        
        return ctx.JSON(resp.StatusCode, map[string]string{"status": "healthy"})
    }))
    
    // API 网关 - 路由到后端服务 (API Gateway - route to backend services)
    api := framework.Group("/api")
    api.RegisterRoute("GET", "/users", server.HandlerFunc(func(ctx server.Context) error {
        service := registry.GetHealthyService("user-service")
        if service == nil {
            return ctx.JSON(503, map[string]string{"error": "User service unavailable"})
        }
        
        // 代理请求到用户服务 (Proxy request to user service)
        url := fmt.Sprintf("http://%s:%d/users", service.Host, service.Port)
        resp, err := http.Get(url)
        if err != nil {
            return ctx.JSON(500, map[string]string{"error": "Failed to connect to user service"})
        }
        defer resp.Body.Close()
        
        var users []interface{}
        if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
            return ctx.JSON(500, map[string]string{"error": "Failed to decode response"})
        }
        
        return ctx.JSON(200, users)
    }))
    
    if err := manager.Start(context.Background()); err != nil {
        log.Fatal("Server failed:", err)
    }
}
```

## 中间件集成示例

### 综合中间件栈

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/middleware"
    _ "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin"
)

func main() {
    config := server.DefaultServerConfig()
    config.Framework = "gin"
    config.Port = 8080
    
    // 配置中间件 (Configure middleware)
    config.Middleware.Logger = server.LoggerMiddlewareConfig{
        Enabled: true,
        Format:  "json",
        SkipPaths: []string{"/health", "/metrics"},
    }
    
    config.Middleware.Recovery = server.RecoveryMiddlewareConfig{
        Enabled:    true,
        PrintStack: true,
    }
    
    config.CORS = server.CORSConfig{
        Enabled:      true,
        AllowOrigins: []string{"https://example.com"},
        AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders: []string{"Content-Type", "Authorization"},
    }
    
    config.Middleware.RateLimit = server.RateLimitMiddlewareConfig{
        Enabled: true,
        Rate:    100.0, // 每秒 100 个请求 (100 requests per second)
        Burst:   200,   // 突发 200 个请求 (Burst of 200 requests)
        KeyFunc: "ip",
    }
    
    manager, err := server.CreateServerManager("gin", config)
    if err != nil {
        log.Fatal("Failed to create server:", err)
    }
    
    framework := manager.GetFramework()
    
    // 注册全局中间件 (Register global middleware)
    framework.RegisterMiddleware(middleware.NewLogger(config.Middleware.Logger))
    framework.RegisterMiddleware(middleware.NewRecovery(config.Middleware.Recovery))
    framework.RegisterMiddleware(middleware.NewCORS(config.CORS))
    framework.RegisterMiddleware(middleware.NewSecurity())
    framework.RegisterMiddleware(middleware.NewCompression())
    framework.RegisterMiddleware(middleware.NewMetrics())
    
    // 受保护的 API 路由 (Protected API routes)
    api := framework.Group("/api/v1")
    api.RegisterMiddleware(middleware.NewAuth(config.Middleware.Auth))
    api.RegisterMiddleware(middleware.NewRateLimit(config.Middleware.RateLimit))
    
    api.RegisterRoute("GET", "/users", server.HandlerFunc(func(ctx server.Context) error {
        // 模拟数据库查询 (Simulate database query)
        time.Sleep(10 * time.Millisecond)
        
        users := []map[string]interface{}{
            {"id": 1, "name": "Alice", "email": "alice@example.com"},
            {"id": 2, "name": "Bob", "email": "bob@example.com"},
        }
        
        return ctx.JSON(200, users)
    }))
    
    // 健康检查端点 (Health check endpoint)
    framework.RegisterRoute("GET", "/health", server.HandlerFunc(func(ctx server.Context) error {
        return ctx.JSON(200, map[string]interface{}{
            "status":    "healthy",
            "timestamp": time.Now().Unix(),
            "version":   "1.0.0",
        })
    }))
    
    // 指标端点 (Metrics endpoint)
    framework.RegisterRoute("GET", "/metrics", server.HandlerFunc(func(ctx server.Context) error {
        metrics := map[string]interface{}{
            "requests_total":     1000,
            "requests_per_sec":   50.5,
            "response_time_avg":  15.2,
            "active_connections": 25,
        }
        
        return ctx.JSON(200, metrics)
    }))
    
    if err := manager.Start(context.Background()); err != nil {
        log.Fatal("Server failed:", err)
    }
}
```

## 下一步

- **[最佳实践](07_best_practices.md)** - 生产部署指南
- **[故障排除](08_troubleshooting.md)** - 常见问题和解决方案
- **[API 参考](09_api_reference.md)** - 完整 API 文档