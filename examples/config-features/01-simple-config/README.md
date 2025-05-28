/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This documentation was collaboratively developed by Martin and AI Assistant.
 */

# Simple Configuration Loading Example

[中文版本](README_zh.md)

This example demonstrates the most basic usage of the configuration module, focusing on:
- Loading configuration from YAML files
- Using default values with struct tags
- Basic configuration validation
- Error handling for configuration issues

## What This Example Shows

### 1. Basic Configuration Loading
- Simple struct-based configuration definition
- YAML file parsing and mapping to Go structs
- Nested configuration structures

### 2. Default Value Mechanisms
- Using `default` struct tags for fallback values
- Handling different data types (string, int, bool, duration, slices)
- Default behavior for nested structs and pointers

### 3. Configuration Validation
- Custom validation logic for configuration values
- Port range validation
- Positive value checks for timeouts and connections
- Conditional validation based on feature flags

### 4. Error Handling Integration
- Configuration loading error handling
- Error codes and types from the errors module
- Graceful error reporting and debugging information

## File Structure

```
01-simple-config/
├── main.go         # Main application demonstrating config loading
├── config.yaml     # Example configuration file
├── README.md       # This file
└── README_zh.md    # Chinese version
```

## Configuration Structure

### Application Settings
```yaml
app_name: "Advanced Simple Config Example"
version: "2.0.0"
debug: false
port: 9080
timeout: "45s"
interval: "10m"
```

### Network Configuration
```yaml
allowed_ips:
  - "127.0.0.1"
  - "::1"
  - "192.168.1.0/24"
```

### Feature Flags
```yaml
features:
  - "authentication"
  - "authorization"
  - "logging"
  - "metrics"
```

### Database Settings
```yaml
database:
  driver: "postgresql"
  host: "db.example.com"
  port: 5432
  max_connections: 20
  enable_ssl: true
```

### Cache Configuration
```yaml
cache:
  enabled: true
  type: "redis"
  host: "cache.example.com"
  ttl: "2h"
```

## Running the Example

### Prerequisites
- Go 1.21 or later
- No external dependencies required

### Basic Run
```bash
cd examples/config-features/01-simple-config
go run main.go
```

### Test with Missing Config File
```bash
# Rename config file to see default behavior
mv config.yaml config.yaml.backup
go run main.go
# Restore config file
mv config.yaml.backup config.yaml
```

### Test with Invalid Config
```bash
# Create invalid config for testing error handling
cat > invalid-config.yaml << EOF
port: 99999  # Invalid port number
timeout: "-5s"  # Invalid timeout
EOF

# This would show validation errors
# go run main.go -config invalid-config.yaml
```

## Expected Output

The example produces structured output showing:

### 1. Default Values Demonstration
```
=== Demonstrating Default Values ===
Loading configuration with defaults only (no config file)...
Expected error (file not found): ...
Zero values (before applying defaults):
  AppName: ''
  Port: 0
  Debug: false
  Features: [] (length: 0)
```

### 2. Successful Configuration Loading
```
=== Loading Configuration from File ===
✓ Configuration loaded successfully!

=== Validating Configuration ===
✓ Configuration validation passed!
```

### 3. Configuration Summary
```
=== Configuration Summary ===
Application:
  Name: Advanced Simple Config Example
  Version: 2.0.0
  Debug: false
  Port: 9080
  Timeout: 45s
  Features: [authentication authorization logging metrics tracing caching]

Database:
  Driver: postgresql
  Host: db.example.com
  Max Connections: 20
  Enable SSL: true

Cache:
  Enabled: true
  Type: redis
  TTL: 2h
```

### 4. Configuration Usage
```
=== Using Configuration ===
Starting Advanced Simple Config Example version 2.0.0...
Server will listen on port 9080
Database: postgresql://db.example.com:5432/production_db
Cache: redis://cache.example.com:6379 (TTL: 2h)
```

## Key Learning Points

### 1. Struct Tag Defaults
```go
type Config struct {
    Port    int           `mapstructure:"port" default:"8080"`
    Timeout time.Duration `mapstructure:"timeout" default:"30s"`
    Debug   bool          `mapstructure:"debug" default:"true"`
    Features []string     `mapstructure:"features" default:"auth,logging"`
}
```

### 2. Configuration Loading Pattern
```go
var cfg SimpleAppConfig
err := config.LoadConfig(&cfg, 
    config.WithConfigFile("config.yaml", "yaml"))
if err != nil {
    // Handle error with detailed information
    if coder := errors.GetCoder(err); coder != nil {
        fmt.Printf("Config Error [%d]: %s", coder.Code(), coder.String())
    }
}
```

### 3. Validation Best Practices
```go
func validateConfig(cfg *SimpleAppConfig) error {
    if cfg.Port < 1 || cfg.Port > 65535 {
        return errors.Errorf("invalid port: %d", cfg.Port)
    }
    // More validation logic...
}
```

### 4. Safe Configuration Access
```go
// Check for nil pointers before accessing nested config
if cfg.Database != nil {
    connectionString := fmt.Sprintf("%s://%s:%d/%s", 
        cfg.Database.Driver, cfg.Database.Host, 
        cfg.Database.Port, cfg.Database.Database)
}
```

## Common Patterns Demonstrated

1. **Embedded Base Configuration**: Though not used in this simple example, shows how to structure config for larger applications

2. **Type-Safe Defaults**: Using struct tags for compile-time default definitions

3. **Nested Configuration**: Organizing related settings into logical groups

4. **Configuration Validation**: Implementing business rule validation separately from loading

5. **Error Handling**: Proper error propagation and user-friendly error messages

## Troubleshooting

### Issue: Configuration file not found
**Cause**: `config.yaml` doesn't exist in the working directory
**Solution**: Ensure the file exists or check the file path

### Issue: YAML parsing errors
**Cause**: Invalid YAML syntax in configuration file
**Solution**: Validate YAML syntax using online validators or `yamllint`

### Issue: Type conversion errors
**Cause**: Configuration values don't match expected Go types
**Solution**: Check data types in YAML match struct field types

### Issue: Default values not applied
**Cause**: Struct tags missing or incorrect syntax
**Solution**: Verify `default:"value"` tag syntax and field types

## Next Steps

After understanding this simple example:

1. Explore [02-hot-reload](../02-hot-reload/) for dynamic configuration updates
2. Try [03-env-override](../03-env-override/) for environment-based configuration
3. See [04-default-values](../04-default-values/) for advanced default handling
4. Check [05-multiple-formats](../05-multiple-formats/) for different file formats

## Related Documentation

- [Configuration Module Overview](../../../docs/usage-guides/config/en/)
- [Configuration Loading Guide](../../../docs/usage-guides/config/en/01_quick_start.md)
- [Default Values Documentation](../../../docs/usage-guides/config/en/02_configuration_options.md)