# Configuration Best Practices

This guide outlines recommended patterns and practices for using the configuration module effectively in production applications.

## Configuration Structure Design

### 1. Use Hierarchical Organization

Organize your configuration into logical groups:

```go
type AppConfig struct {
    // Application metadata
    App struct {
        Name        string `mapstructure:"name" default:"my-app"`
        Version     string `mapstructure:"version" default:"1.0.0"`
        Environment string `mapstructure:"environment" default:"development"`
    } `mapstructure:"app"`
    
    // Server configuration
    Server struct {
        HTTP struct {
            Host         string        `mapstructure:"host" default:"localhost"`
            Port         int           `mapstructure:"port" default:"8080"`
            ReadTimeout  time.Duration `mapstructure:"read_timeout" default:"30s"`
            WriteTimeout time.Duration `mapstructure:"write_timeout" default:"30s"`
        } `mapstructure:"http"`
        
        GRPC struct {
            Host string `mapstructure:"host" default:"localhost"`
            Port int    `mapstructure:"port" default:"9090"`
        } `mapstructure:"grpc"`
    } `mapstructure:"server"`
    
    // Database configuration
    Database struct {
        Primary struct {
            URL          string        `mapstructure:"url"`
            MaxConns     int           `mapstructure:"max_connections" default:"25"`
            MaxIdleConns int           `mapstructure:"max_idle_connections" default:"5"`
            ConnTimeout  time.Duration `mapstructure:"connection_timeout" default:"5s"`
        } `mapstructure:"primary"`
        
        Cache struct {
            URL string        `mapstructure:"url"`
            TTL time.Duration `mapstructure:"ttl" default:"1h"`
        } `mapstructure:"cache"`
    } `mapstructure:"database"`
    
    // Observability
    Observability struct {
        Logging struct {
            Level  string `mapstructure:"level" default:"info"`
            Format string `mapstructure:"format" default:"json"`
        } `mapstructure:"logging"`
        
        Metrics struct {
            Enabled bool   `mapstructure:"enabled" default:"true"`
            Port    int    `mapstructure:"port" default:"9090"`
            Path    string `mapstructure:"path" default:"/metrics"`
        } `mapstructure:"metrics"`
        
        Tracing struct {
            Enabled    bool   `mapstructure:"enabled" default:"false"`
            Endpoint   string `mapstructure:"endpoint"`
            SampleRate float64 `mapstructure:"sample_rate" default:"0.1"`
        } `mapstructure:"tracing"`
    } `mapstructure:"observability"`
}
```

### 2. Use Meaningful Default Values

Provide sensible defaults that work for development:

```go
type DatabaseConfig struct {
    // Good: Provides a working default for development
    URL string `mapstructure:"url" default:"postgres://localhost:5432/myapp_dev"`
    
    // Good: Conservative default that works everywhere
    MaxConnections int `mapstructure:"max_connections" default:"10"`
    
    // Good: Reasonable timeout
    Timeout time.Duration `mapstructure:"timeout" default:"30s"`
    
    // Bad: No default for required field
    // Password string `mapstructure:"password"`
    
    // Better: Use environment variable for sensitive data
    Password string `mapstructure:"password" default:"${DB_PASSWORD}"`
}
```

### 3. Use Appropriate Data Types

Choose the right data types for your configuration:

```go
type Config struct {
    // Use time.Duration for time-based values
    Timeout time.Duration `mapstructure:"timeout" default:"30s"`
    
    // Use specific types for better validation
    Port int `mapstructure:"port" default:"8080"`
    
    // Use bool for feature flags
    EnableMetrics bool `mapstructure:"enable_metrics" default:"true"`
    
    // Use slices for lists
    AllowedOrigins []string `mapstructure:"allowed_origins"`
    
    // Use maps for key-value pairs
    Headers map[string]string `mapstructure:"headers"`
    
    // Use custom types for validation
    LogLevel LogLevel `mapstructure:"log_level" default:"info"`
}

type LogLevel string

const (
    LogLevelDebug LogLevel = "debug"
    LogLevelInfo  LogLevel = "info"
    LogLevelWarn  LogLevel = "warn"
    LogLevelError LogLevel = "error"
)
```

## Environment Variable Strategy

### 1. Use Consistent Naming

Establish a clear naming convention:

```go
// Good: Consistent prefix and structure
// APP_SERVER_HTTP_HOST
// APP_SERVER_HTTP_PORT
// APP_DATABASE_PRIMARY_URL
// APP_OBSERVABILITY_LOGGING_LEVEL

config.WithEnvPrefix("APP")
```

### 2. Separate Sensitive Data

Keep sensitive data in environment variables:

```go
type Config struct {
    Database struct {
        // Non-sensitive: can be in config file
        Host string `mapstructure:"host" default:"localhost"`
        Port int    `mapstructure:"port" default:"5432"`
        Name string `mapstructure:"name" default:"myapp"`
        
        // Sensitive: should be in environment variables
        Username string `mapstructure:"username"`
        Password string `mapstructure:"password"`
    } `mapstructure:"database"`
    
    Security struct {
        // Sensitive: should be in environment variables
        JWTSecret    string `mapstructure:"jwt_secret"`
        APIKey       string `mapstructure:"api_key"`
        EncryptionKey string `mapstructure:"encryption_key"`
    } `mapstructure:"security"`
}
```

Environment variables:
```bash
export APP_DATABASE_USERNAME=myuser
export APP_DATABASE_PASSWORD=mysecretpassword
export APP_SECURITY_JWT_SECRET=my-jwt-secret-key
export APP_SECURITY_API_KEY=my-api-key
```

### 3. Use Environment-Specific Overrides

```bash
# Development
export APP_APP_ENVIRONMENT=development
export APP_OBSERVABILITY_LOGGING_LEVEL=debug
export APP_DATABASE_PRIMARY_URL=postgres://localhost:5432/myapp_dev

# Production
export APP_APP_ENVIRONMENT=production
export APP_OBSERVABILITY_LOGGING_LEVEL=warn
export APP_DATABASE_PRIMARY_URL=postgres://prod-db:5432/myapp_prod
```

## Configuration Validation

### 1. Implement Validation Functions

```go
func (c *AppConfig) Validate() error {
    var errors []error
    
    // Validate server configuration
    if c.Server.HTTP.Port < 1 || c.Server.HTTP.Port > 65535 {
        errors = append(errors, fmt.Errorf("invalid HTTP port: %d", c.Server.HTTP.Port))
    }
    
    if c.Server.GRPC.Port < 1 || c.Server.GRPC.Port > 65535 {
        errors = append(errors, fmt.Errorf("invalid GRPC port: %d", c.Server.GRPC.Port))
    }
    
    // Validate database configuration
    if c.Database.Primary.URL == "" {
        errors = append(errors, fmt.Errorf("database URL is required"))
    }
    
    if c.Database.Primary.MaxConns < 1 {
        errors = append(errors, fmt.Errorf("max connections must be positive"))
    }
    
    // Validate observability configuration
    validLogLevels := map[string]bool{
        "debug": true, "info": true, "warn": true, "error": true,
    }
    if !validLogLevels[c.Observability.Logging.Level] {
        errors = append(errors, fmt.Errorf("invalid log level: %s", c.Observability.Logging.Level))
    }
    
    if len(errors) > 0 {
        return fmt.Errorf("configuration validation failed: %v", errors)
    }
    
    return nil
}
```

### 2. Use Validation in Callbacks

```go
cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
    cfg := currentCfg.(*AppConfig)
    
    // Validate configuration before applying
    if err := cfg.Validate(); err != nil {
        return fmt.Errorf("invalid configuration: %w", err)
    }
    
    // Apply configuration
    return applyConfiguration(cfg)
})
```

## Error Handling Patterns

### 1. Graceful Degradation

```go
func loadConfiguration() (*AppConfig, error) {
    var cfg AppConfig
    
    err := config.LoadConfig(
        &cfg,
        config.WithConfigFile("config.yaml", ""),
        config.WithEnvPrefix("APP"),
        config.WithEnvVarOverride(true),
    )
    
    if err != nil {
        // Log the error but continue with defaults
        log.Printf("Warning: Failed to load configuration file: %v", err)
        log.Printf("Using default configuration")
        
        // Initialize with defaults
        cfg = AppConfig{} // This will use struct tag defaults
    }
    
    // Always validate the final configuration
    if err := cfg.Validate(); err != nil {
        return nil, fmt.Errorf("configuration validation failed: %w", err)
    }
    
    return &cfg, nil
}
```

### 2. Configuration Fallback Chain

```go
func loadConfigurationWithFallback() (*AppConfig, error) {
    var cfg AppConfig
    var err error
    
    // Try primary configuration file
    err = config.LoadConfig(&cfg, 
        config.WithConfigFile("config.yaml", ""),
        config.WithEnvPrefix("APP"),
    )
    
    if err != nil {
        log.Printf("Primary config failed: %v", err)
        
        // Try fallback configuration file
        err = config.LoadConfig(&cfg,
            config.WithConfigFile("config.default.yaml", ""),
            config.WithEnvPrefix("APP"),
        )
        
        if err != nil {
            log.Printf("Fallback config failed: %v", err)
            
            // Use embedded defaults
            cfg = getDefaultConfiguration()
        }
    }
    
    return &cfg, nil
}
```

## Testing Strategies

### 1. Configuration Testing

```go
func TestConfigurationLoading(t *testing.T) {
    tests := []struct {
        name           string
        configContent  string
        envVars        map[string]string
        expectedConfig AppConfig
        expectError    bool
    }{
        {
            name: "valid configuration",
            configContent: `
server:
  http:
    host: "0.0.0.0"
    port: 8080
database:
  primary:
    url: "postgres://localhost/test"
`,
            envVars: map[string]string{
                "APP_DATABASE_PRIMARY_URL": "postgres://test-db/test",
            },
            expectedConfig: AppConfig{
                Server: ServerConfig{
                    HTTP: HTTPConfig{
                        Host: "0.0.0.0",
                        Port: 8080,
                    },
                },
                Database: DatabaseConfig{
                    Primary: PrimaryDBConfig{
                        URL: "postgres://test-db/test", // From env var
                    },
                },
            },
            expectError: false,
        },
        {
            name: "invalid port",
            configContent: `
server:
  http:
    port: 99999
`,
            expectError: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Set environment variables
            for key, value := range tt.envVars {
                os.Setenv(key, value)
                defer os.Unsetenv(key)
            }
            
            // Create temporary config file
            tmpFile, err := os.CreateTemp("", "config-*.yaml")
            require.NoError(t, err)
            defer os.Remove(tmpFile.Name())
            
            _, err = tmpFile.WriteString(tt.configContent)
            require.NoError(t, err)
            tmpFile.Close()
            
            // Load configuration
            var cfg AppConfig
            err = config.LoadConfig(&cfg,
                config.WithConfigFile(tmpFile.Name(), ""),
                config.WithEnvPrefix("APP"),
                config.WithEnvVarOverride(true),
            )
            
            if tt.expectError {
                assert.Error(t, err)
                return
            }
            
            require.NoError(t, err)
            assert.Equal(t, tt.expectedConfig.Server.HTTP.Host, cfg.Server.HTTP.Host)
            assert.Equal(t, tt.expectedConfig.Server.HTTP.Port, cfg.Server.HTTP.Port)
        })
    }
}
```

### 2. Hot Reload Testing

```go
func TestHotReload(t *testing.T) {
    // Create temporary config file
    tmpFile, err := os.CreateTemp("", "config-*.yaml")
    require.NoError(t, err)
    defer os.Remove(tmpFile.Name())
    
    // Initial configuration
    initialConfig := `
server:
  http:
    port: 8080
`
    _, err = tmpFile.WriteString(initialConfig)
    require.NoError(t, err)
    tmpFile.Close()
    
    var cfg AppConfig
    var callbackCalled bool
    var newPort int
    
    // Load with hot reload
    cm, err := config.LoadConfigAndWatch(&cfg,
        config.WithConfigFile(tmpFile.Name(), ""),
        config.WithHotReload(true),
    )
    require.NoError(t, err)
    
    // Register callback
    cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
        appCfg := currentCfg.(*AppConfig)
        callbackCalled = true
        newPort = appCfg.Server.HTTP.Port
        return nil
    })
    
    // Verify initial configuration
    assert.Equal(t, 8080, cfg.Server.HTTP.Port)
    
    // Update configuration file
    updatedConfig := `
server:
  http:
    port: 9090
`
    err = os.WriteFile(tmpFile.Name(), []byte(updatedConfig), 0644)
    require.NoError(t, err)
    
    // Wait for hot reload (with timeout)
    timeout := time.After(5 * time.Second)
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()
    
    for {
        select {
        case <-timeout:
            t.Fatal("Hot reload callback was not called within timeout")
        case <-ticker.C:
            if callbackCalled && newPort == 9090 {
                return // Success
            }
        }
    }
}
```

## Production Deployment

### 1. Configuration Management

```yaml
# config/production.yaml
app:
  name: "my-app"
  version: "1.2.3"
  environment: "production"

server:
  http:
    host: "0.0.0.0"
    port: 8080
    read_timeout: "30s"
    write_timeout: "30s"

database:
  primary:
    # URL comes from environment variable
    max_connections: 50
    max_idle_connections: 10
    connection_timeout: "10s"

observability:
  logging:
    level: "info"
    format: "json"
  metrics:
    enabled: true
    port: 9090
  tracing:
    enabled: true
    sample_rate: 0.1
```

### 2. Docker Configuration

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o myapp .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/myapp .
COPY --from=builder /app/config/production.yaml ./config/

# Configuration comes from environment variables
ENV APP_APP_ENVIRONMENT=production
ENV APP_SERVER_HTTP_HOST=0.0.0.0
ENV APP_SERVER_HTTP_PORT=8080

CMD ["./myapp"]
```

### 3. Kubernetes Configuration

```yaml
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  config.yaml: |
    app:
      name: "my-app"
      environment: "production"
    server:
      http:
        host: "0.0.0.0"
        port: 8080
    observability:
      logging:
        level: "info"
        format: "json"

---
# k8s/secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
type: Opaque
data:
  database-url: <base64-encoded-database-url>
  jwt-secret: <base64-encoded-jwt-secret>

---
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
spec:
  template:
    spec:
      containers:
      - name: my-app
        image: my-app:latest
        env:
        - name: APP_DATABASE_PRIMARY_URL
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: database-url
        - name: APP_SECURITY_JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: jwt-secret
        volumeMounts:
        - name: config
          mountPath: /app/config
      volumes:
      - name: config
        configMap:
          name: app-config
```

## Performance Considerations

### 1. Minimize Hot Reload Impact

```go
// Good: Efficient callback that only updates what changed
cm.RegisterSectionChangeCallback("server", func(v *viper.Viper) error {
    newPort := v.GetInt("server.http.port")
    if currentPort != newPort {
        return restartHTTPServer(newPort)
    }
    return nil // No change, no action needed
})

// Bad: Inefficient callback that always restarts everything
cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
    // This restarts everything even for minor changes
    return restartAllServices(currentCfg)
})
```

### 2. Cache Configuration Values

```go
type ConfigCache struct {
    mu     sync.RWMutex
    config *AppConfig
}

func (cc *ConfigCache) GetConfig() *AppConfig {
    cc.mu.RLock()
    defer cc.mu.RUnlock()
    return cc.config
}

func (cc *ConfigCache) UpdateConfig(newConfig *AppConfig) {
    cc.mu.Lock()
    defer cc.mu.Unlock()
    cc.config = newConfig
}

var configCache = &ConfigCache{}

// Update cache in callback
cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
    cfg := currentCfg.(*AppConfig)
    configCache.UpdateConfig(cfg)
    return nil
})
```

## Security Best Practices

### 1. Protect Sensitive Configuration

```go
// Never log sensitive configuration
func logConfiguration(cfg *AppConfig) {
    // Good: Create a safe copy for logging
    safeCfg := *cfg
    safeCfg.Database.Primary.Password = "[REDACTED]"
    safeCfg.Security.JWTSecret = "[REDACTED]"
    
    log.Printf("Configuration loaded: %+v", safeCfg)
}
```

### 2. Validate File Permissions

```go
func validateConfigFilePermissions(filename string) error {
    info, err := os.Stat(filename)
    if err != nil {
        return err
    }
    
    mode := info.Mode()
    if mode&0077 != 0 {
        return fmt.Errorf("config file %s has overly permissive permissions: %v", filename, mode)
    }
    
    return nil
}
```

## Next Steps

- [Integration Examples](05_integration_examples.md) - See real-world examples
- [Troubleshooting](06_troubleshooting.md) - Common issues and solutions 