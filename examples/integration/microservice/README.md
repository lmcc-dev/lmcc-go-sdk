# Microservice Integration Example

[中文版本](README_zh.md)

This example demonstrates a comprehensive microservice implementation using the lmcc-go-sdk package. It showcases how to build a production-ready microservice with observability features, health checks, and proper error handling.

## Features

- **User Service**: Complete CRUD operations for user management
- **gRPC-style Service Structure**: Clean service layer architecture
- **HTTP REST API**: RESTful endpoints for user operations
- **Health Checks**: Comprehensive health monitoring
- **Metrics Collection**: Request counting, error tracking, and latency measurement
- **Distributed Tracing**: Request tracing across service layers
- **Database Layer**: Simulated database operations with proper error handling
- **Configuration Management**: YAML-based configuration with defaults
- **Structured Logging**: JSON and text logging formats with trace correlation

## Service Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Server   │    │   User Service  │    │ Database Layer  │
│                 │    │                 │    │                 │
│ • REST API      │───▶│ • Business Logic│───▶│ • Data Storage  │
│ • Health Check  │    │ • Validation    │    │ • Query Logic   │
│ • Metrics       │    │ • Error Handling│    │ • Transactions  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ Metrics Service │    │ Tracing Service │    │ Logging Service │
│                 │    │                 │    │                 │
│ • Request Count │    │ • Trace ID Gen  │    │ • Structured    │
│ • Error Rates   │    │ • Context Pass  │    │ • Correlation   │
│ • Latency Stats │    │ • Span Tracking │    │ • Multi-level   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Components

### User Service
- **GetUser**: Retrieve user by ID with validation and error handling
- **CreateUser**: Create new users with conflict detection
- **Request/Response Models**: Structured data transfer objects

### Observability Stack
- **MetricsCollector**: Records requests, errors, and latency statistics
- **TracingService**: Generates and propagates trace IDs across operations
- **HealthChecker**: Multi-layer health monitoring (database, memory, response time)

### HTTP Layer
- **REST API Endpoints**: 
  - `GET /api/users/{id}` - Get user by ID
  - `POST /api/users` - Create new user
  - `GET /health` - Health check
  - `GET /metrics` - Service metrics

## Configuration

The service supports YAML configuration with sensible defaults:

```yaml
service:
  name: "user-service"
  version: "v1.0.0"
  port: 50051

http:
  port: 8080
  health_check_path: "/health"
  metrics_path: "/metrics"

database:
  host: "localhost"
  port: 5432
  name: "userdb"
  max_conns: 10
  conn_timeout: 5

logging:
  level: "info"
  format: "json"
  output_paths: ["stdout"]

observability:
  tracing_enabled: true
  metrics_enabled: true
  service_mesh: "istio"
```

## Usage

### Running the Example

```bash
# Navigate to the microservice example directory
cd examples/integration/microservice

# Run the microservice
go run main.go
```

### Expected Output

The example demonstrates:

1. **Service Initialization**: Configuration loading and component setup
2. **User Operations**: CRUD operations with validation and error handling
3. **Health Monitoring**: Database, memory, and response time checks
4. **Metrics Collection**: Request counting and performance tracking
5. **HTTP API Testing**: REST endpoint validation
6. **Distributed Tracing**: Request correlation across service layers

### Sample Output

```
=== Microservice Integration Example ===
This example demonstrates a microservice with gRPC, observability, and health checks.

=== Running Microservice Demonstration ===

1. Testing User Operations:
   ✅ Get user successful: alice (trace_20241123_12345)
   ✅ Get non-existent user correctly failed: user not found
   ✅ Create user successful: user_1700123456 (trace_20241123_12346)

2. Testing Health Check:
   Service Status: healthy
   Timestamp: 2024-11-23T10:15:30Z

3. Service Metrics:
   get_user: 2 requests
   create_user: 1 requests

Starting HTTP server for additional testing...
HTTP server running on port 8080
Available endpoints:
  GET  http://localhost:8080/health
  GET  http://localhost:8080/metrics
  GET  http://localhost:8080/api/users/user_001

=== Running HTTP Endpoint Tests ===

Testing health check endpoint:
   ✅ Health check successful: 200
Testing metrics endpoint:
   ✅ Metrics endpoint successful: 200
Testing user API:
   ✅ User API successful: 200

=== HTTP Tests Completed ===
=== Example completed successfully ===
```

## Key Learning Points

### 1. Service Layer Architecture
- Clean separation between HTTP, business logic, and data layers
- Dependency injection pattern for component composition
- Interface-based design for testability

### 2. Observability Integration
- **Metrics**: Request counting, error classification, latency tracking
- **Tracing**: Request correlation across service boundaries
- **Logging**: Structured logging with trace context

### 3. Error Handling Patterns
- Custom error types with context
- Error propagation across service layers
- Graceful error responses with client-friendly messages

### 4. Configuration Management
- YAML-based configuration with struct tags
- Default value handling
- Environment-specific overrides

### 5. Health Check Design
- Multi-layer health monitoring
- Dependency health verification
- Performance threshold validation

## Production Considerations

This example demonstrates patterns suitable for production microservices:

- **Configuration**: Externalized configuration with defaults
- **Logging**: Structured logging with correlation IDs
- **Monitoring**: Health checks and metrics collection
- **Error Handling**: Comprehensive error classification and handling
- **Performance**: Latency tracking and optimization opportunities
- **Maintainability**: Clean architecture with separated concerns

## Testing

The example includes comprehensive testing scenarios:

- Successful operations
- Error conditions (validation, not found)
- Health check validation
- HTTP endpoint verification
- Metrics collection validation

## Integration Patterns

This microservice example can be extended with:

- **Database Integration**: Replace simulated database with real persistence
- **Message Queues**: Add async communication patterns
- **Service Discovery**: Integrate with service mesh or discovery systems
- **API Gateway**: Add authentication and rate limiting
- **Container Deployment**: Docker and Kubernetes configurations

For more advanced patterns, see the other integration examples in this directory. 