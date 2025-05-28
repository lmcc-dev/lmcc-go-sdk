# Web应用集成示例

[English Version](README.md)

此示例展示了使用lmcc-go-sdk包的综合Web应用程序实现。它演示了如何构建一个生产就绪的REST API，集成了日志记录、错误处理、配置管理和中间件模式。

## 功能特性

- **REST API**: 完整的HTTP API，包含CRUD操作
- **中间件架构**: 日志记录和错误处理中间件
- **请求链路追踪**: 请求ID生成和关联
- **配置管理**: 基于YAML的配置，带有默认值
- **结构化日志**: JSON和文本日志，支持请求关联
- **错误处理**: 全面的错误处理，带有适当的HTTP状态码
- **健康检查**: 服务健康监控端点
- **优雅关闭**: 适当的服务器生命周期管理
- **自动化测试**: 内置API端点测试套件

## 应用架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP服务器    │    │   中间件        │    │   API处理器     │
│                 │    │                 │    │                 │
│ • 路由          │───▶│ • 日志记录      │───▶│ • 请求处理      │
│ • 超时          │    │ • 错误处理      │    │ • 验证          │
│ • 优雅停止      │    │ • 请求ID        │    │ • 响应          │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  配置管理       │    │   用户服务      │    │   模拟数据库    │
│                 │    │                 │    │                 │
│ • YAML配置      │    │ • 业务逻辑      │    │ • 数据操作      │
│ • 默认值        │    │ • 验证          │    │ • 错误模拟      │
│ • 环境变量      │    │ • 错误处理      │    │ • 延迟模拟      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## API端点

### 核心端点
- **`GET /api/health`** - 健康检查端点
- **`GET /api/users/{id}`** - 通过ID检索用户
- **`POST /api/users`** - 创建新用户

### 响应格式

所有API响应都遵循一致的格式：

```json
{
  "success": true,
  "data": {
    "id": "user_123",
    "username": "john_doe",
    "email": "john@example.com",
    "created": "2024-11-23T10:15:30Z",
    "last_seen": "2024-11-23T12:30:45Z"
  },
  "request_id": "req_1700123456789",
  "timestamp": "2024-11-23T12:30:45Z"
}
```

错误响应包含错误详情：

```json
{
  "success": false,
  "error": "user not found",
  "request_id": "req_1700123456789",
  "timestamp": "2024-11-23T12:30:45Z"
}
```

## 配置

Web应用程序支持带有合理默认值的YAML配置：

```yaml
server:
  port: 8080
  host: "0.0.0.0"
  read_timeout: 30
  write_timeout: 30
  idle_timeout: 120

database:
  host: "localhost"
  port: 5432
  name: "webapp"
  user: "webapp_user"
  password: "webapp_pass"
  max_conns: 10
  conn_timeout: 5

logging:
  level: "info"
  format: "json"  # json, text
  output_paths: ["stdout"]
  enable_caller: true
  enable_stacktrace: false

api:
  rate_limit: 100
  auth_required: true
  api_version: "v1"
```

## 使用方法

### 运行应用程序

```bash
# 导航到Web应用程序目录
cd examples/integration/web-app

# 运行应用程序
go run main.go
```

### 预期输出

应用程序演示了：

1. **服务器启动**: 配置加载和服务器初始化
2. **端点测试**: 所有API端点的自动化测试
3. **请求日志**: 带有关联ID的详细请求/响应日志
4. **错误处理**: 各种错误场景和适当的HTTP状态码
5. **优雅关闭**: 清洁的服务器关闭过程

### 示例输出

```
=== Web应用集成示例 ===
此示例演示了一个完整的Web应用程序，集成了日志记录、错误处理和配置管理。

Web服务器运行在 http://0.0.0.0:8080
可用端点:
  GET  /api/health
  GET  /api/users/{id}
  POST /api/users

=== 运行自动化测试 ===

1. 测试健康检查端点:
   ✅ 通过: 健康检查端点工作正常

2. 测试获取用户端点（成功）:
   ✅ 通过: 获取用户端点工作正常

3. 测试获取用户端点（错误处理）:
   ✅ 通过: 错误处理工作正常

4. 测试创建用户端点（成功）:
   ✅ 通过: 创建用户端点工作正常

5. 测试创建用户端点（验证错误）:
   ✅ 通过: 验证错误处理工作正常

6. 测试不存在的端点:
   ✅ 通过: 404处理工作正常

=== 自动化测试完成 ===

正在关闭Web应用程序...
=== 示例成功完成 ===
```

## 关键学习要点

### 1. 中间件架构
- 带有关联ID的请求/响应日志
- 集中式错误处理和恢复
- 模块化中间件组合

### 2. 请求生命周期管理
- 用于链路追踪的请求ID生成
- 跨层上下文传播
- 带有请求关联的结构化日志

### 3. 错误处理模式
- 带HTTP状态映射的自定义错误类型
- 用户友好消息的优雅错误响应
- 带上下文和堆栈跟踪的错误日志

### 4. 配置管理
- 带结构标签的基于YAML配置
- 默认值处理
- 环境特定覆盖

### 5. 服务层设计
- HTTP和业务逻辑的清晰分离
- 可测试性的依赖注入
- 基于接口的服务设计

## 实现亮点

### 中间件实现

```go
type LoggingMiddleware struct {
    logger log.Logger
}

func (m *LoggingMiddleware) Handler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 请求ID生成和上下文设置
        requestID := fmt.Sprintf("req_%d", time.Now().UnixNano())
        ctx := context.WithValue(r.Context(), "request_id", requestID)
        
        // 请求日志和处理
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### 服务层模式

```go
type UserService struct {
    logger log.Logger
    config *AppConfig
}

func (s *UserService) GetUser(ctx context.Context, userID string) (*User, error) {
    // 带日志和错误处理的业务逻辑
    requestID := ctx.Value("request_id")
    logger := s.logger.WithValues("request_id", requestID)
    
    // 验证、数据库操作、错误处理
    return user, nil
}
```

### 响应处理

```go
func (h *APIHandler) writeJSONResponse(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}, err error) {
    response := APIResponse{
        Success:   err == nil,
        Data:      data,
        RequestID: r.Context().Value("request_id").(string),
        Timestamp: time.Now(),
    }
    
    if err != nil {
        response.Error = err.Error()
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(response)
}
```

## 测试功能

应用程序包含全面的自动化测试：

### 测试场景

1. **健康检查**: 验证服务健康端点
2. **用户检索（成功）**: 测试成功的用户数据检索
3. **用户检索（错误）**: 测试不存在用户的错误处理
4. **用户创建（成功）**: 测试成功的用户创建
5. **用户创建（验证）**: 测试验证错误处理
6. **404处理**: 测试不存在端点的响应

### 测试实现

每个测试验证：
- HTTP状态码
- 响应格式一致性
- 错误消息准确性
- 请求关联ID
- 日志输出

## 生产环境考虑

此示例演示了适用于生产Web应用程序的模式：

- **安全性**: 输入验证和清理
- **性能**: 请求超时和连接限制
- **可观察性**: 全面的日志记录和请求追踪
- **可靠性**: 错误处理和优雅关闭
- **可维护性**: 清洁架构和关注点分离
- **可测试性**: 自动化测试套件和基于接口的设计

## 扩展点

此Web应用程序可以通过以下方式扩展：

- **身份验证**: JWT令牌、OAuth集成
- **数据库集成**: PostgreSQL、MongoDB、Redis
- **缓存**: 响应缓存、会话存储
- **速率限制**: 请求节流和滥用防护
- **监控**: 指标收集、健康检查
- **API文档**: OpenAPI/Swagger集成
- **容器部署**: Docker、Kubernetes配置

## 性能特性

- **连接池**: 数据库连接管理
- **请求超时**: 可配置的超时设置
- **优雅关闭**: 清洁的服务器生命周期管理
- **内存管理**: 高效的请求处理

## 监控和可观察性

- **请求追踪**: 用于关联的唯一请求ID
- **结构化日志**: 用于日志聚合的JSON格式
- **健康端点**: 服务健康监控
- **错误分类**: HTTP状态码映射

## 安全考虑

- **输入验证**: 请求数据清理
- **错误清理**: 安全的错误消息暴露
- **请求大小限制**: 内容长度限制
- **超时保护**: 请求超时强制执行

此示例为使用Go和lmcc-go-sdk框架构建生产就绪的Web应用程序提供了坚实的基础，演示了API设计、错误处理和可观察性的最佳实践。