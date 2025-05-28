# Web Application Integration Example

[中文版本](README_zh.md)

This example demonstrates a comprehensive web application implementation using the lmcc-go-sdk package. It showcases how to build a production-ready REST API with integrated logging, error handling, configuration management, and middleware patterns.

## Features

- **REST API**: Complete HTTP API with CRUD operations
- **Middleware Architecture**: Logging and error handling middleware
- **Request Tracing**: Request ID generation and correlation
- **Configuration Management**: YAML-based configuration with defaults
- **Structured Logging**: JSON and text logging with request correlation
- **Error Handling**: Comprehensive error handling with proper HTTP status codes
- **Health Checks**: Service health monitoring endpoint
- **Graceful Shutdown**: Proper server lifecycle management
- **Automated Testing**: Built-in test suite for API endpoints

## Application Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Server   │    │   Middleware    │    │   API Handler   │
│                 │    │                 │    │                 │
│ • Routing       │───▶│ • Logging       │───▶│ • Request Proc  │
│ • Timeouts      │    │ • Error Handling│    │ • Validation    │
│ • Graceful Stop │    │ • Request ID    │    │ • Response      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Configuration  │    │   User Service  │    │   Simulated DB  │
│                 │    │                 │    │                 │
│ • YAML Config   │    │ • Business Logic│    │ • Data Operations│
│ • Defaults      │    │ • Validation    │    │ • Error Simulation│
│ • Environment   │    │ • Error Handling│    │ • Latency Sim   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## API Endpoints

### Core Endpoints
- **`GET /api/health`** - Health check endpoint
- **`GET /api/users/{id}`** - Retrieve user by ID
- **`POST /api/users`** - Create new user

### Response Format

All API responses follow a consistent format:

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

Error responses include error details:

```json
{
  "success": false,
  "error": "user not found",
  "request_id": "req_1700123456789",
  "timestamp": "2024-11-23T12:30:45Z"
}
```

## Configuration

The web application supports YAML configuration with sensible defaults:

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

## Usage

### Running the Application

```bash
# Navigate to the web application directory
cd examples/integration/web-app

# Run the application
go run main.go
```

### Expected Output

The application demonstrates:

1. **Server Startup**: Configuration loading and server initialization
2. **Endpoint Testing**: Automated tests for all API endpoints
3. **Request Logging**: Detailed request/response logging with correlation IDs
4. **Error Handling**: Various error scenarios and proper HTTP status codes
5. **Graceful Shutdown**: Clean server shutdown process

### Sample Output

```
=== Web Application Integration Example ===
This example demonstrates a complete web application with integrated logging, error handling, and configuration management.

Web server is running on http://0.0.0.0:8080
Available endpoints:
  GET  /api/health
  GET  /api/users/{id}
  POST /api/users

=== Running Automated Tests ===

1. Testing Health Check Endpoint:
   ✅ PASSED: Health check endpoint working

2. Testing Get User Endpoint (Success):
   ✅ PASSED: Get user endpoint working

3. Testing Get User Endpoint (Error Handling):
   ✅ PASSED: Error handling working correctly

4. Testing Create User Endpoint (Success):
   ✅ PASSED: Create user endpoint working

5. Testing Create User Endpoint (Validation Error):
   ✅ PASSED: Validation error handling working

6. Testing Non-existent Endpoint:
   ✅ PASSED: 404 handling working correctly

=== Automated Tests Completed ===

Shutting down web application...
=== Example completed successfully ===
```

## Key Learning Points

### 1. Middleware Architecture
- Request/response logging with correlation IDs
- Centralized error handling and recovery
- Modular middleware composition

### 2. Request Lifecycle Management
- Request ID generation for tracing
- Context propagation across layers
- Structured logging with request correlation

### 3. Error Handling Patterns
- Custom error types with HTTP status mapping
- Graceful error responses with user-friendly messages
- Error logging with context and stack traces

### 4. Configuration Management
- YAML-based configuration with struct tags
- Default value handling
- Environment-specific overrides

### 5. Service Layer Design
- Clean separation between HTTP and business logic
- Dependency injection for testability
- Interface-based service design

## Implementation Highlights

### Middleware Implementation

```go
type LoggingMiddleware struct {
    logger log.Logger
}

func (m *LoggingMiddleware) Handler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Request ID generation and context setup
        requestID := fmt.Sprintf("req_%d", time.Now().UnixNano())
        ctx := context.WithValue(r.Context(), "request_id", requestID)
        
        // Request logging and processing
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### Service Layer Pattern

```go
type UserService struct {
    logger log.Logger
    config *AppConfig
}

func (s *UserService) GetUser(ctx context.Context, userID string) (*User, error) {
    // Business logic with logging and error handling
    requestID := ctx.Value("request_id")
    logger := s.logger.WithValues("request_id", requestID)
    
    // Validation, database operations, error handling
    return user, nil
}
```

### Response Handling

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

## Testing Features

The application includes comprehensive automated testing:

### Test Scenarios

1. **Health Check**: Validates service health endpoint
2. **User Retrieval (Success)**: Tests successful user data retrieval
3. **User Retrieval (Error)**: Tests error handling for non-existent users
4. **User Creation (Success)**: Tests successful user creation
5. **User Creation (Validation)**: Tests validation error handling
6. **404 Handling**: Tests non-existent endpoint responses

### Test Implementation

Each test validates:
- HTTP status codes
- Response format consistency
- Error message accuracy
- Request correlation IDs
- Logging output

## Production Considerations

This example demonstrates patterns suitable for production web applications:

- **Security**: Input validation and sanitization
- **Performance**: Request timeouts and connection limits
- **Observability**: Comprehensive logging and request tracing
- **Reliability**: Error handling and graceful shutdown
- **Maintainability**: Clean architecture and separation of concerns
- **Testability**: Automated test suite and interface-based design

## Extension Points

This web application can be extended with:

- **Authentication**: JWT tokens, OAuth integration
- **Database Integration**: PostgreSQL, MongoDB, Redis
- **Caching**: Response caching, session storage
- **Rate Limiting**: Request throttling and abuse prevention
- **Monitoring**: Metrics collection, health checks
- **API Documentation**: OpenAPI/Swagger integration
- **Container Deployment**: Docker, Kubernetes configurations

## Performance Features

- **Connection Pooling**: Database connection management
- **Request Timeouts**: Configurable timeout settings
- **Graceful Shutdown**: Clean server lifecycle management
- **Memory Management**: Efficient request processing

## Monitoring and Observability

- **Request Tracing**: Unique request IDs for correlation
- **Structured Logging**: JSON format for log aggregation
- **Health Endpoints**: Service health monitoring
- **Error Classification**: HTTP status code mapping

## Security Considerations

- **Input Validation**: Request data sanitization
- **Error Sanitization**: Safe error message exposure
- **Request Size Limits**: Content length restrictions
- **Timeout Protection**: Request timeout enforcement

This example provides a solid foundation for building production-ready web applications with Go and the lmcc-go-sdk framework, demonstrating best practices for API design, error handling, and observability.