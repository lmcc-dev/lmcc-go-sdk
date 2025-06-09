# Integration Examples

This guide provides real-world integration examples showing how to use the lmcc-go-sdk server module in various scenarios.

## Basic REST API

### Simple CRUD Operations

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
    
    // Setup CRUD routes
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

## Database Integration

### Using GORM with PostgreSQL

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
    // Connect to database
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        dsn = "host=localhost user=postgres password=postgres dbname=testdb port=5432 sslmode=disable"
    }
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect database:", err)
    }
    
    // Auto-migrate schema
    db.AutoMigrate(&User{})
    
    userService := &UserService{db: db}
    
    // Setup server
    config := server.DefaultServerConfig()
    config.Framework = "gin"
    config.Port = 8080
    
    manager, err := server.CreateServerManager("gin", config)
    if err != nil {
        log.Fatal("Failed to create server:", err)
    }
    
    framework := manager.GetFramework()
    
    // Setup routes
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

## Authentication with JWT

### JWT Middleware Integration

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
    
    // Public routes
    framework.RegisterRoute("POST", "/login", server.HandlerFunc(login))
    
    // Protected routes
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

## WebSocket Chat Server

### Real-time Communication

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
        return true // Allow all origins
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
    
    // WebSocket endpoint
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
        
        // Handle incoming messages
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

## Next Steps

- **[Best Practices](07_best_practices.md)** - Production deployment guidelines
- **[Troubleshooting](08_troubleshooting.md)** - Common issues and solutions
- **[API Reference](09_api_reference.md)** - Complete API documentation
