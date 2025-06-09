# Configuration Guide

This guide covers all configuration options available in the lmcc-go-sdk server module. The server module provides comprehensive configuration capabilities for web server setup, middleware management, and framework-specific customization.

## Overview

The server module uses the `ServerConfig` structure to define all configuration options. Configuration can be loaded from YAML files, environment variables, or set programmatically. The module integrates seamlessly with the `pkg/config` module for advanced configuration management.

## Basic Configuration

### Default Configuration

The simplest way to get started is using the default configuration:

```go
package main

import (
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
)

func main() {
    // Get default configuration
    config := server.DefaultServerConfig()
    
    // Customize as needed
    config.Framework = "gin"
    config.Port = 8080
    config.Mode = "production"
    
    // Create server with configuration
    manager, err := server.CreateServerManager("gin", config)
    if err != nil {
        panic(err)
    }
    
    // Start server
    // ... rest of the application
}
```

### Complete Configuration Structure

Here's the complete structure with all available options:

```go
type ServerConfig struct {
    Framework        string                    // Framework: gin, echo, fiber
    Host            string                    // Listen host
    Port            int                       // Listen port
    Mode            string                    // Running mode: debug, release, test
    ReadTimeout     time.Duration             // HTTP read timeout
    WriteTimeout    time.Duration             // HTTP write timeout
    IdleTimeout     time.Duration             // HTTP idle timeout
    MaxHeaderBytes  int                       // Maximum request header size
    CORS            CORSConfig               // CORS configuration
    Middleware      MiddlewareConfig         // Middleware settings
    TLS             TLSConfig                // TLS/HTTPS settings
    Plugins         map[string]interface{}   // Plugin-specific config
    GracefulShutdown GracefulShutdownConfig  // Graceful shutdown settings
}
```

## Core Server Settings

### Framework Selection

Choose from supported web frameworks:

```go
config := server.DefaultServerConfig()

// Gin framework (high performance, rich ecosystem)
config.Framework = "gin"

// Echo framework (minimalist, fast)
config.Framework = "echo"

// Fiber framework (Express-inspired, ultra-fast)
config.Framework = "fiber"
```

### Network Configuration

Configure host and port settings:

```go
config := &server.ServerConfig{
    Framework: "gin",
    Host:     "0.0.0.0",    // Listen on all interfaces
    Port:     8080,         // Standard HTTP port
    Mode:     "production", // Running mode
}

// Alternative configurations
config.Host = "127.0.0.1"  // Localhost only
config.Host = "localhost"   // Same as 127.0.0.1
config.Port = 3000         // Alternative port
```

### Running Modes

The server supports three running modes:

```go
// Debug mode - detailed logging, hot reload
config.Mode = "debug"

// Release mode - optimized for production
config.Mode = "release"

// Test mode - minimal output, fast startup
config.Mode = "test"

// Check mode programmatically
if config.IsDebugMode() {
    // Debug-specific logic
}
```

### Timeout Configuration

Configure various timeout settings:

```go
config := &server.ServerConfig{
    Framework:      "gin",
    ReadTimeout:    30 * time.Second,  // Time to read request
    WriteTimeout:   30 * time.Second,  // Time to write response
    IdleTimeout:    60 * time.Second,  // Keep-alive timeout
    MaxHeaderBytes: 1 << 20,           // 1MB header limit
}

// Custom timeout example
config.ReadTimeout = 5 * time.Minute   // For large uploads
config.WriteTimeout = 2 * time.Minute  // For streaming responses
```

## CORS Configuration

Configure Cross-Origin Resource Sharing (CORS) settings:

```go
config := server.DefaultServerConfig()

// Basic CORS setup
config.CORS = server.CORSConfig{
    Enabled:          true,
    AllowOrigins:     []string{"https://example.com", "https://app.example.com"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
    ExposeHeaders:    []string{"Content-Length", "Content-Type"},
    AllowCredentials: true,
    MaxAge:          12 * time.Hour,
}

// Development CORS (allow all)
config.CORS = server.CORSConfig{
    Enabled:      true,
    AllowOrigins: []string{"*"},
    AllowMethods: []string{"*"},
    AllowHeaders: []string{"*"},
}

// Disable CORS
config.CORS.Enabled = false
```

## Middleware Configuration

Configure built-in middleware components:

### Logger Middleware

```go
config.Middleware.Logger = server.LoggerMiddlewareConfig{
    Enabled:     true,
    Format:      "json",        // json or text
    SkipPaths:   []string{"/health", "/metrics"},
    IncludeBody: false,         // Include request body in logs
    MaxBodySize: 1024,          // Max body size to log (bytes)
}

// Custom logger configuration
config.Middleware.Logger = server.LoggerMiddlewareConfig{
    Enabled:     true,
    Format:      "text",
    SkipPaths:   []string{"/favicon.ico", "/static/*"},
    IncludeBody: true,
    MaxBodySize: 2048,
}
```

### Recovery Middleware

```go
config.Middleware.Recovery = server.RecoveryMiddlewareConfig{
    Enabled:             true,
    PrintStack:          true,   // Print stack trace on panic
    DisableStackAll:     false,  // Disable all stack traces
    DisableColorConsole: false,  // Disable colored output
}

// Production recovery settings
config.Middleware.Recovery = server.RecoveryMiddlewareConfig{
    Enabled:             true,
    PrintStack:          false,  // Don't print stack in production
    DisableStackAll:     true,
    DisableColorConsole: true,
}
```

### Rate Limiting Middleware

```go
config.Middleware.RateLimit = server.RateLimitMiddlewareConfig{
    Enabled: true,
    Rate:    100.0,      // 100 requests per second
    Burst:   200,        // Allow burst of 200 requests
    KeyFunc: "ip",       // Rate limit by IP address
}

// Alternative rate limiting strategies
config.Middleware.RateLimit.KeyFunc = "user"   // Rate limit by user ID
config.Middleware.RateLimit.KeyFunc = "custom" // Custom key function
```

### Authentication Middleware

```go
config.Middleware.Auth = server.AuthMiddlewareConfig{
    Enabled:   true,
    Type:      "jwt",
    SkipPaths: []string{"/auth/login", "/auth/register", "/health"},
    JWT: server.JWTConfig{
        Secret:         "your-secret-key",
        Issuer:         "your-app",
        Audience:       "your-users",
        ExpirationTime: 24 * time.Hour,
        RefreshTime:    7 * 24 * time.Hour,
    },
}

// Basic authentication
config.Middleware.Auth.Type = "basic"

// Disable authentication
config.Middleware.Auth.Enabled = false
```

## TLS/HTTPS Configuration

Configure HTTPS and TLS settings:

### Basic TLS Setup

```go
config.TLS = server.TLSConfig{
    Enabled:  true,
    CertFile: "/path/to/certificate.crt",
    KeyFile:  "/path/to/private.key",
}
```

### Auto TLS (Let's Encrypt)

```go
config.TLS = server.TLSConfig{
    Enabled: true,
    AutoTLS: true,
    Domains: []string{"example.com", "www.example.com"},
}
```

### Development TLS

```go
// Self-signed certificate for development
config.TLS = server.TLSConfig{
    Enabled:  true,
    CertFile: "./dev-cert.pem",
    KeyFile:  "./dev-key.pem",
}
```

## Graceful Shutdown Configuration

Configure graceful shutdown behavior:

```go
config.GracefulShutdown = server.GracefulShutdownConfig{
    Enabled:  true,
    Timeout:  30 * time.Second,  // Maximum shutdown time
    WaitTime: 5 * time.Second,   // Wait before starting shutdown
}

// Quick shutdown for development
config.GracefulShutdown = server.GracefulShutdownConfig{
    Enabled:  true,
    Timeout:  5 * time.Second,
    WaitTime: 1 * time.Second,
}

// Disable graceful shutdown
config.GracefulShutdown.Enabled = false
```

## Plugin-Specific Configuration

Configure framework-specific options:

### Gin Plugin Configuration

```go
config.Plugins = map[string]interface{}{
    "gin": map[string]interface{}{
        "disable_color":        false,
        "trust_proxies":        []string{"127.0.0.1"},
        "max_multipart_memory": 32 << 20, // 32MB
        "enable_html_render":   true,
    },
}
```

### Echo Plugin Configuration

```go
config.Plugins = map[string]interface{}{
    "echo": map[string]interface{}{
        "hide_banner":      true,
        "hide_port":        false,
        "debug":           config.IsDebugMode(),
        "disable_startup_message": false,
    },
}
```

### Fiber Plugin Configuration

```go
config.Plugins = map[string]interface{}{
    "fiber": map[string]interface{}{
        "prefork":              false,
        "strict_routing":       false,
        "case_sensitive":       false,
        "immutable":           false,
        "unescaped_path":      false,
        "etag":                false,
        "body_limit":          4 * 1024 * 1024, // 4MB
        "read_buffer_size":    4096,
        "write_buffer_size":   4096,
    },
}
```

## Configuration from Files

### YAML Configuration

Create a `config.yaml` file:

```yaml
server:
  framework: gin
  host: 0.0.0.0
  port: 8080
  mode: production
  read-timeout: 30s
  write-timeout: 30s
  idle-timeout: 60s
  max-header-bytes: 1048576
  
  cors:
    enabled: true
    allow-origins:
      - "https://example.com"
      - "https://app.example.com"
    allow-methods:
      - GET
      - POST
      - PUT
      - DELETE
      - OPTIONS
    allow-headers:
      - Origin
      - Content-Type
      - Accept
      - Authorization
    allow-credentials: true
    max-age: 12h
  
  middleware:
    logger:
      enabled: true
      format: json
      skip-paths:
        - /health
        - /metrics
      include-body: false
      max-body-size: 1024
    
    recovery:
      enabled: true
      print-stack: false
      disable-stack-all: true
      disable-color-console: true
    
    rate-limit:
      enabled: true
      rate: 100.0
      burst: 200
      key-func: ip
    
    auth:
      enabled: true
      type: jwt
      skip-paths:
        - /auth/login
        - /auth/register
        - /health
      jwt:
        secret: your-secret-key
        issuer: your-app
        audience: your-users
        expiration-time: 24h
        refresh-time: 168h
  
  tls:
    enabled: false
    cert-file: ""
    key-file: ""
    auto-tls: false
    domains: []
  
  graceful-shutdown:
    enabled: true
    timeout: 30s
    wait-time: 5s
  
  plugins:
    gin:
      disable_color: false
      trust_proxies:
        - 127.0.0.1
      max_multipart_memory: 33554432
      enable_html_render: true
```

Load configuration using the config module:

```go
package main

import (
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
)

type AppConfig struct {
    Server server.ServerConfig `mapstructure:"server"`
}

func main() {
    var cfg AppConfig
    
    // Load configuration from file
    err := config.LoadConfig(&cfg,
        config.WithConfigFile("config.yaml", ""),
        config.WithHotReload(true),
    )
    if err != nil {
        panic(err)
    }
    
    // Create server with loaded configuration
    manager, err := server.CreateServerManager(cfg.Server.Framework, &cfg.Server)
    if err != nil {
        panic(err)
    }
    
    // Start server
    // ... rest of application
}
```

### Environment Variables

Configuration can also be loaded from environment variables:

```bash
# Server settings
SERVER_FRAMEWORK=gin
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_MODE=production

# CORS settings
SERVER_CORS_ENABLED=true
SERVER_CORS_ALLOW_ORIGINS=https://example.com,https://app.example.com

# Middleware settings
SERVER_MIDDLEWARE_LOGGER_ENABLED=true
SERVER_MIDDLEWARE_LOGGER_FORMAT=json
SERVER_MIDDLEWARE_RECOVERY_ENABLED=true

# TLS settings
SERVER_TLS_ENABLED=false
SERVER_TLS_CERT_FILE=/path/to/cert.pem
SERVER_TLS_KEY_FILE=/path/to/key.pem
```

Load from environment:

```go
func main() {
    var cfg AppConfig
    
    err := config.LoadConfig(&cfg,
        config.WithEnvPrefix("MYAPP"),
        config.WithEnvKeyReplacer(".", "_"),
    )
    if err != nil {
        panic(err)
    }
    
    // Use loaded configuration
    manager, err := server.CreateServerManager(cfg.Server.Framework, &cfg.Server)
    // ...
}
```

## Configuration Validation

The server module provides built-in configuration validation:

```go
config := &server.ServerConfig{
    Framework: "gin",
    Host:     "localhost",
    Port:     8080,
    Mode:     "production",
}

// Validate configuration
if err := config.Validate(); err != nil {
    panic(fmt.Sprintf("Invalid configuration: %v", err))
}

// Configuration is automatically validated when creating server
manager, err := server.CreateServerManager("gin", config)
if err != nil {
    // Handle validation errors
    panic(err)
}
```

### Custom Validation

Add custom validation logic:

```go
func validateCustomConfig(config *server.ServerConfig) error {
    // Custom business logic validation
    if config.Port == 80 && !config.TLS.Enabled {
        return fmt.Errorf("HTTP port 80 requires TLS in production")
    }
    
    if config.Mode == "production" && config.Middleware.Auth.Enabled == false {
        return fmt.Errorf("authentication required in production mode")
    }
    
    return nil
}

func main() {
    config := server.DefaultServerConfig()
    
    // Standard validation
    if err := config.Validate(); err != nil {
        panic(err)
    }
    
    // Custom validation
    if err := validateCustomConfig(config); err != nil {
        panic(err)
    }
    
    // Create server
    manager, err := server.CreateServerManager("gin", config)
    // ...
}
```

## Configuration Hot Reload

Enable hot reloading for configuration changes:

```go
package main

import (
    "context"
    "log"
    
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
    "github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
)

type AppConfig struct {
    Server server.ServerConfig `mapstructure:"server"`
}

func main() {
    var cfg AppConfig
    
    // Enable hot reload
    configManager, err := config.LoadConfig(&cfg,
        config.WithConfigFile("config.yaml", ""),
        config.WithHotReload(true),
    )
    if err != nil {
        panic(err)
    }
    
    // Create server
    manager, err := server.CreateServerManager(cfg.Server.Framework, &cfg.Server)
    if err != nil {
        panic(err)
    }
    
    // Listen for configuration changes
    go func() {
        for {
            select {
            case <-configManager.WatchConfig():
                log.Println("Configuration changed, reloading...")
                
                // Reload configuration
                var newCfg AppConfig
                if err := configManager.Unmarshal(&newCfg); err != nil {
                    log.Printf("Failed to reload config: %v", err)
                    continue
                }
                
                // Apply new configuration
                // Note: Some changes may require server restart
                if err := applyNewConfig(manager, &newCfg.Server); err != nil {
                    log.Printf("Failed to apply new config: %v", err)
                }
            }
        }
    }()
    
    // Start server
    if err := manager.Start(context.Background()); err != nil {
        panic(err)
    }
}

func applyNewConfig(manager *server.ServerManager, newConfig *server.ServerConfig) error {
    // Apply configuration changes that don't require restart
    // For changes that require restart, you would need to:
    // 1. Stop the current server
    // 2. Create a new server with new config
    // 3. Start the new server
    
    return nil
}
```

## Configuration Best Practices

### 1. Environment-Specific Configuration

Use different configurations for different environments:

```go
// config/development.yaml
server:
  mode: debug
  cors:
    allow-origins: ["*"]
  middleware:
    logger:
      format: text
      include-body: true

// config/production.yaml
server:
  mode: production
  cors:
    allow-origins: ["https://yourdomain.com"]
  middleware:
    logger:
      format: json
      include-body: false
  tls:
    enabled: true
    auto-tls: true
    domains: ["yourdomain.com"]
```

### 2. Secure Configuration

Keep sensitive information in environment variables:

```yaml
# config.yaml
server:
  framework: gin
  middleware:
    auth:
      jwt:
        secret: ${JWT_SECRET}  # From environment
        issuer: ${JWT_ISSUER}
  tls:
    cert-file: ${TLS_CERT_FILE}
    key-file: ${TLS_KEY_FILE}
```

### 3. Configuration Validation

Always validate configuration early:

```go
func main() {
    config := loadConfiguration()
    
    // Validate early
    if err := config.Validate(); err != nil {
        log.Fatalf("Invalid configuration: %v", err)
    }
    
    // Additional validation
    if err := validateBusinessRules(config); err != nil {
        log.Fatalf("Configuration validation failed: %v", err)
    }
    
    // Now safe to use configuration
    startServer(config)
}
```

### 4. Default Fallbacks

Provide sensible defaults:

```go
func loadConfiguration() *server.ServerConfig {
    config := server.DefaultServerConfig()
    
    // Override with file config if available
    if err := loadFromFile(config, "config.yaml"); err != nil {
        log.Printf("Config file not found, using defaults: %v", err)
    }
    
    // Override with environment variables
    loadFromEnv(config)
    
    return config
}
```

## Configuration Examples

### Microservice Configuration

```go
config := &server.ServerConfig{
    Framework: "gin",
    Host:     "0.0.0.0",
    Port:     8080,
    Mode:     "production",
    
    Middleware: server.MiddlewareConfig{
        Logger: server.LoggerMiddlewareConfig{
            Enabled: true,
            Format:  "json",
            SkipPaths: []string{"/health", "/metrics"},
        },
        Recovery: server.RecoveryMiddlewareConfig{
            Enabled:             true,
            PrintStack:          false,
            DisableColorConsole: true,
        },
        Auth: server.AuthMiddlewareConfig{
            Enabled:   true,
            Type:      "jwt",
            SkipPaths: []string{"/health", "/metrics"},
        },
    },
    
    GracefulShutdown: server.GracefulShutdownConfig{
        Enabled: true,
        Timeout: 30 * time.Second,
    },
}
```

### API Gateway Configuration

```go
config := &server.ServerConfig{
    Framework: "fiber",
    Host:     "0.0.0.0",
    Port:     8080,
    Mode:     "production",
    
    CORS: server.CORSConfig{
        Enabled:          true,
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"*"},
        AllowHeaders:     []string{"*"},
        AllowCredentials: false,
    },
    
    Middleware: server.MiddlewareConfig{
        Logger: server.LoggerMiddlewareConfig{
            Enabled: true,
            Format:  "json",
        },
        RateLimit: server.RateLimitMiddlewareConfig{
            Enabled: true,
            Rate:    1000,
            Burst:   2000,
            KeyFunc: "ip",
        },
    },
    
    Plugins: map[string]interface{}{
        "fiber": map[string]interface{}{
            "prefork":       true,
            "read_timeout":  "30s",
            "write_timeout": "30s",
        },
    },
}
```

### Development Configuration

```go
config := &server.ServerConfig{
    Framework: "gin",
    Host:     "localhost",
    Port:     3000,
    Mode:     "debug",
    
    CORS: server.CORSConfig{
        Enabled:      true,
        AllowOrigins: []string{"*"},
        AllowMethods: []string{"*"},
        AllowHeaders: []string{"*"},
    },
    
    Middleware: server.MiddlewareConfig{
        Logger: server.LoggerMiddlewareConfig{
            Enabled:     true,
            Format:      "text",
            IncludeBody: true,
        },
        Recovery: server.RecoveryMiddlewareConfig{
            Enabled:    true,
            PrintStack: true,
        },
        Auth: server.AuthMiddlewareConfig{
            Enabled: false, // Disabled for development
        },
    },
    
    GracefulShutdown: server.GracefulShutdownConfig{
        Enabled: false, // Quick shutdown for development
    },
}
```

## Configuration Reference

For a complete reference of all configuration options, see:

- **[Framework Plugins](03_framework_plugins.md)** - Framework-specific configuration
- **[Middleware System](04_middleware_system.md)** - Middleware configuration details
- **[Server Management](05_server_management.md)** - Advanced server configuration
- **[Module Specification](09_module_specification.md)** - Complete API reference

## Troubleshooting

### Common Configuration Issues

1. **Invalid Port**: Ensure port is between 1-65535
2. **Framework Not Found**: Verify framework plugin is imported
3. **TLS Certificate Issues**: Check certificate paths and permissions
4. **CORS Problems**: Verify origin patterns and allowed methods
5. **Middleware Conflicts**: Check middleware order and compatibility

See **[Troubleshooting Guide](08_troubleshooting.md)** for detailed solutions.
