# Configuration Options

This document describes all available configuration options and how to use them effectively.

## Function Options

The configuration module uses the functional options pattern for clean and extensible API design.

### WithConfigFile

Specifies the configuration file path and optional search paths.

```go
config.WithConfigFile("config.yaml", "")
config.WithConfigFile("app.yaml", "/etc/myapp:/home/user/.config")
```

**Parameters:**
- `filename`: The configuration file name
- `searchPaths`: Colon-separated list of directories to search (optional)

**Supported file extensions:**
- `.yaml`, `.yml`
- `.json`
- `.toml`
- `.hcl`
- `.ini`
- `.properties`

### WithEnvPrefix

Sets a prefix for environment variable binding.

```go
config.WithEnvPrefix("APP")
```

With prefix "APP", the following mappings apply:
- `server.host` → `APP_SERVER_HOST`
- `database.url` → `APP_DATABASE_URL`
- `debug` → `APP_DEBUG`

### WithEnvVarOverride

Enables environment variables to override configuration file values.

```go
config.WithEnvVarOverride(true)  // Enable override
config.WithEnvVarOverride(false) // Disable override
```

### WithHotReload

Enables automatic configuration reloading when files change.

```go
config.WithHotReload(true)  // Enable hot reload
config.WithHotReload(false) // Disable hot reload
```

**Note:** Only available with `LoadConfigAndWatch()` function.

### WithEnvKeyReplacer

Customizes how configuration keys are transformed to environment variable names.

```go
import "strings"

config.WithEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
```

This replaces dots and hyphens with underscores in environment variable names.

## Configuration Structure Tags

### mapstructure Tag

Maps configuration keys to struct fields.

```go
type Config struct {
    ServerPort int `mapstructure:"server_port"`
    DebugMode  bool `mapstructure:"debug"`
}
```

### default Tag

Specifies default values for fields.

```go
type Config struct {
    Host    string `mapstructure:"host" default:"localhost"`
    Port    int    `mapstructure:"port" default:"8080"`
    Timeout string `mapstructure:"timeout" default:"30s"`
    Debug   bool   `mapstructure:"debug" default:"false"`
}
```

**Supported default value types:**
- String: `default:"localhost"`
- Integer: `default:"8080"`
- Boolean: `default:"true"` or `default:"false"`
- Duration: `default:"30s"`

## Advanced Configuration Patterns

### Nested Configuration

```go
type Config struct {
    Server struct {
        HTTP struct {
            Host string `mapstructure:"host" default:"localhost"`
            Port int    `mapstructure:"port" default:"8080"`
        } `mapstructure:"http"`
        
        GRPC struct {
            Host string `mapstructure:"host" default:"localhost"`
            Port int    `mapstructure:"port" default:"9090"`
        } `mapstructure:"grpc"`
    } `mapstructure:"server"`
    
    Database struct {
        Primary struct {
            URL      string `mapstructure:"url"`
            MaxConns int    `mapstructure:"max_connections" default:"10"`
        } `mapstructure:"primary"`
        
        Cache struct {
            URL string `mapstructure:"url"`
            TTL string `mapstructure:"ttl" default:"1h"`
        } `mapstructure:"cache"`
    } `mapstructure:"database"`
}
```

### Array and Slice Configuration

```go
type Config struct {
    Servers []ServerConfig `mapstructure:"servers"`
    Tags    []string       `mapstructure:"tags"`
}

type ServerConfig struct {
    Name string `mapstructure:"name"`
    URL  string `mapstructure:"url"`
}
```

YAML configuration:
```yaml
servers:
  - name: "primary"
    url: "https://primary.example.com"
  - name: "secondary"
    url: "https://secondary.example.com"

tags:
  - "production"
  - "web-server"
  - "critical"
```

### Map Configuration

```go
type Config struct {
    Features map[string]bool   `mapstructure:"features"`
    Limits   map[string]int    `mapstructure:"limits"`
    Metadata map[string]string `mapstructure:"metadata"`
}
```

YAML configuration:
```yaml
features:
  authentication: true
  rate_limiting: false
  metrics: true

limits:
  max_requests: 1000
  max_connections: 100

metadata:
  version: "1.0.0"
  environment: "production"
```

## Environment Variable Examples

### Basic Environment Variables

```bash
# With prefix "APP"
export APP_SERVER_HOST=production.example.com
export APP_SERVER_PORT=443
export APP_DEBUG=true
export APP_DATABASE_URL=postgres://user:pass@db.example.com/prod
```

### Nested Configuration

```bash
# For nested structures
export APP_SERVER_HTTP_HOST=web.example.com
export APP_SERVER_HTTP_PORT=80
export APP_SERVER_GRPC_HOST=grpc.example.com
export APP_SERVER_GRPC_PORT=9090
```

### Array Configuration

```bash
# Arrays can be specified as comma-separated values
export APP_TAGS=production,web-server,critical
```

## Configuration File Examples

### Complete YAML Example

```yaml
# Application configuration
app:
  name: "My Application"
  version: "1.0.0"
  environment: "production"

# Server configuration
server:
  http:
    host: "0.0.0.0"
    port: 8080
    timeout: "30s"
  grpc:
    host: "0.0.0.0"
    port: 9090
    timeout: "10s"

# Database configuration
database:
  primary:
    url: "postgres://user:pass@localhost/myapp"
    max_connections: 25
    timeout: "5s"
  cache:
    url: "redis://localhost:6379"
    ttl: "1h"

# Feature flags
features:
  authentication: true
  rate_limiting: true
  metrics: true
  tracing: false

# Logging configuration
logging:
  level: "info"
  format: "json"
  output: ["stdout", "/var/log/app.log"]

# Security settings
security:
  jwt_secret: "your-secret-key"
  cors_origins: ["https://example.com", "https://app.example.com"]
  rate_limit: 100

debug: false
```

### JSON Example

```json
{
  "server": {
    "host": "localhost",
    "port": 8080
  },
  "database": {
    "url": "postgres://localhost/myapp",
    "max_connections": 10
  },
  "features": {
    "authentication": true,
    "metrics": false
  },
  "debug": true
}
```

## Validation

### Built-in Validation

The module automatically validates:
- Type conversions (string to int, bool, etc.)
- Required fields (when no default is provided)
- File format syntax

### Custom Validation

You can add custom validation in your callback functions:

```go
cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
    cfg := currentCfg.(*AppConfig)
    
    // Validate port range
    if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
        return fmt.Errorf("invalid port: %d", cfg.Server.Port)
    }
    
    // Validate required fields
    if cfg.Database.URL == "" {
        return fmt.Errorf("database URL is required")
    }
    
    return nil
})
```

## Best Practices

1. **Use meaningful default values** that work for development
2. **Group related configuration** into nested structs
3. **Use environment variables** for sensitive data
4. **Validate configuration** in callbacks
5. **Document your configuration** structure
6. **Use consistent naming** conventions

## Next Steps

- [Hot Reload](03_hot_reload.md) - Implement dynamic configuration updates
- [Best Practices](04_best_practices.md) - Follow recommended patterns
- [Integration Examples](05_integration_examples.md) - See real-world examples 