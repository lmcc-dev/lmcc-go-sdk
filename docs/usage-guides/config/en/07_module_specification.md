# Module Specification

This document provides a detailed specification of the `pkg/config` module, outlining its public API, interfaces, functions, and predefined components. The module enhances Go's configuration management capabilities by providing hot-reload, environment variable integration, and flexible configuration loading patterns.

## 1. Overview

The `pkg/config` module is designed to make configuration management more robust and developer-friendly. Key features include:
- Automatic configuration loading from multiple file formats (YAML, JSON, TOML, etc.)
- Environment variable integration with prefix support
- Hot-reload capabilities for dynamic configuration updates
- Callback system for configuration change notifications
- Integration with popular configuration libraries (Viper)

## 2. Core Functions

### Configuration Loading Functions

#### LoadConfig
```go
func LoadConfig(cfg interface{}, opts ...Option) error
```
Loads configuration from files and environment variables into the provided struct.

**Parameters:**
- `cfg`: Pointer to the configuration struct
- `opts`: Variadic options for configuration loading

**Returns:**
- `error`: Error if configuration loading fails

**Example:**
```go
type AppConfig struct {
    Server struct {
        Host string `mapstructure:"host" default:"localhost"`
        Port int    `mapstructure:"port" default:"8080"`
    } `mapstructure:"server"`
}

var cfg AppConfig
err := config.LoadConfig(&cfg,
    config.WithConfigFile("config.yaml", ""),
    config.WithEnvPrefix("APP"),
)
```

#### LoadConfigAndWatch
```go
func LoadConfigAndWatch(cfg interface{}, opts ...Option) (*ConfigManager, error)
```
Loads configuration and returns a ConfigManager for hot-reload capabilities.

**Parameters:**
- `cfg`: Pointer to the configuration struct
- `opts`: Variadic options for configuration loading

**Returns:**
- `*ConfigManager`: Manager for handling configuration updates
- `error`: Error if configuration loading fails

**Example:**
```go
cm, err := config.LoadConfigAndWatch(&cfg,
    config.WithConfigFile("config.yaml", ""),
    config.WithHotReload(true),
)
```

## 3. Configuration Options

### Option Type
```go
type Option func(*ConfigOptions)
```
Functional option type for configuring the loading behavior.

### Available Options

#### WithConfigFile
```go
func WithConfigFile(filename, searchPaths string) Option
```
Specifies the configuration file and optional search paths.

**Parameters:**
- `filename`: Name of the configuration file
- `searchPaths`: Colon-separated list of directories to search

#### WithEnvPrefix
```go
func WithEnvPrefix(prefix string) Option
```
Sets a prefix for environment variable binding.

**Parameters:**
- `prefix`: Prefix to use for environment variables

#### WithEnvVarOverride
```go
func WithEnvVarOverride(enable bool) Option
```
Enables or disables environment variable override of configuration values.

**Parameters:**
- `enable`: Whether to enable environment variable override

#### WithHotReload
```go
func WithHotReload(enable bool) Option
```
Enables or disables hot-reload functionality.

**Parameters:**
- `enable`: Whether to enable hot-reload

#### WithEnvKeyReplacer
```go
func WithEnvKeyReplacer(replacer *strings.Replacer) Option
```
Sets a custom key replacer for environment variable names.

**Parameters:**
- `replacer`: String replacer for transforming configuration keys

## 4. ConfigManager Interface

The `ConfigManager` provides methods for managing configuration updates and callbacks.

### Methods

#### RegisterCallback
```go
func (cm *ConfigManager) RegisterCallback(callback func(*viper.Viper, interface{}) error)
```
Registers a global callback function that is called when any configuration changes.

**Parameters:**
- `callback`: Function to call on configuration changes
  - `*viper.Viper`: Viper instance with new configuration
  - `interface{}`: Updated configuration struct
  - `error`: Error returned by callback

#### RegisterSectionChangeCallback
```go
func (cm *ConfigManager) RegisterSectionChangeCallback(section string, callback func(*viper.Viper) error)
```
Registers a callback for specific configuration section changes.

**Parameters:**
- `section`: Configuration section to watch (e.g., "server", "database")
- `callback`: Function to call when the section changes
  - `*viper.Viper`: Viper instance with new configuration
  - `error`: Error returned by callback

#### Stop
```go
func (cm *ConfigManager) Stop()
```
Stops the configuration watcher and cleans up resources.

## 5. Configuration Structure Tags

### mapstructure Tag
Maps configuration keys to struct fields.

```go
type Config struct {
    ServerPort int `mapstructure:"server_port"`
    DebugMode  bool `mapstructure:"debug"`
}
```

### default Tag
Specifies default values for fields when not provided in configuration.

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

## 6. Supported File Formats

The module supports multiple configuration file formats:

| Format | Extensions | Example |
|--------|------------|---------|
| YAML | `.yaml`, `.yml` | `config.yaml` |
| JSON | `.json` | `config.json` |
| TOML | `.toml` | `config.toml` |
| HCL | `.hcl` | `config.hcl` |
| INI | `.ini` | `config.ini` |
| Properties | `.properties` | `config.properties` |

## 7. Environment Variable Mapping

### Automatic Mapping
Configuration keys are automatically mapped to environment variables:

- Nested keys use underscores: `server.host` → `PREFIX_SERVER_HOST`
- Array indices are supported: `servers[0].port` → `PREFIX_SERVERS_0_PORT`
- Custom key replacers can modify the mapping

### Example Mapping
```go
// Configuration structure
type Config struct {
    Server struct {
        Host string `mapstructure:"host"`
        Port int    `mapstructure:"port"`
    } `mapstructure:"server"`
    Database struct {
        URL string `mapstructure:"url"`
    } `mapstructure:"database"`
}

// With prefix "APP", environment variables:
// APP_SERVER_HOST=localhost
// APP_SERVER_PORT=8080
// APP_DATABASE_URL=postgres://localhost/myapp
```

## 8. Hot-Reload Mechanism

### File Watching
The hot-reload mechanism uses file system watchers to monitor configuration files:

1. **File Change Detection**: Monitors configuration file for modifications
2. **Validation**: Validates new configuration before applying
3. **Callback Execution**: Executes registered callbacks with new configuration
4. **Error Handling**: Provides error handling and rollback capabilities

### Callback Types

#### Global Callbacks
Called for any configuration change:
```go
cm.RegisterCallback(func(v *viper.Viper, cfg interface{}) error {
    // Handle any configuration change
    return nil
})
```

#### Section Callbacks
Called only when specific sections change:
```go
cm.RegisterSectionChangeCallback("server", func(v *viper.Viper) error {
    // Handle server configuration changes
    return nil
})
```

## 9. Error Handling

### Configuration Loading Errors
- **File Not Found**: Returns error if configuration file doesn't exist
- **Parse Errors**: Returns error for invalid file format or syntax
- **Validation Errors**: Returns error for invalid configuration values
- **Environment Variable Errors**: Returns error for invalid environment variable values

### Hot-Reload Errors
- **Callback Errors**: If a callback returns an error, the configuration update is rejected
- **Validation Errors**: Invalid new configuration is rejected, keeping the previous configuration
- **File System Errors**: File system watcher errors are logged but don't stop the application

## 10. Integration with Viper

The module is built on top of [Viper](https://github.com/spf13/viper), providing:

### Viper Features
- Multiple configuration file format support
- Environment variable integration
- Configuration key case-insensitivity
- Configuration value type conversion

### Extended Features
- Simplified API with functional options
- Hot-reload capabilities
- Callback system for configuration changes
- Default value support through struct tags

## 11. Best Practices

### Configuration Structure Design
```go
// Good: Organized, nested structure
type Config struct {
    App struct {
        Name    string `mapstructure:"name" default:"myapp"`
        Version string `mapstructure:"version" default:"1.0.0"`
    } `mapstructure:"app"`
    
    Server struct {
        Host string `mapstructure:"host" default:"localhost"`
        Port int    `mapstructure:"port" default:"8080"`
    } `mapstructure:"server"`
}

// Avoid: Flat structure for complex configurations
type Config struct {
    AppName     string `mapstructure:"app_name"`
    AppVersion  string `mapstructure:"app_version"`
    ServerHost  string `mapstructure:"server_host"`
    ServerPort  int    `mapstructure:"server_port"`
}
```

### Error Handling
```go
// Always handle configuration loading errors
if err := config.LoadConfig(&cfg, opts...); err != nil {
    log.Fatalf("Failed to load configuration: %v", err)
}

// Handle callback errors appropriately
cm.RegisterCallback(func(v *viper.Viper, cfg interface{}) error {
    if err := validateConfig(cfg); err != nil {
        return fmt.Errorf("invalid configuration: %w", err)
    }
    return applyConfig(cfg)
})
```

### Environment-Specific Configuration
```go
func getConfigOptions() []config.Option {
    env := os.Getenv("APP_ENV")
    
    opts := []config.Option{
        config.WithEnvPrefix("APP"),
        config.WithEnvVarOverride(true),
    }
    
    switch env {
    case "development":
        opts = append(opts, config.WithConfigFile("config.dev.yaml", ""))
    case "production":
        opts = append(opts, config.WithConfigFile("config.prod.yaml", ""))
    default:
        opts = append(opts, config.WithConfigFile("config.yaml", ""))
    }
    
    return opts
}
```

## 12. Performance Considerations

### File Watching Overhead
- File system watchers have minimal overhead
- Only active when hot-reload is enabled
- Automatically cleaned up when ConfigManager is stopped

### Memory Usage
- Configuration is loaded into memory once
- Hot-reload creates new instances but cleans up old ones
- Callback functions should avoid memory leaks

### Callback Performance
- Callbacks should be fast to avoid blocking configuration updates
- Long-running operations should be performed asynchronously
- Error handling in callbacks should be robust

## 13. Thread Safety

### Concurrent Access
- Configuration loading is not thread-safe during initialization
- Hot-reload callbacks are executed sequentially
- Applications should implement their own synchronization for configuration access

### Recommended Pattern
```go
type Application struct {
    config *AppConfig
    mutex  sync.RWMutex
}

func (app *Application) GetConfig() *AppConfig {
    app.mutex.RLock()
    defer app.mutex.RUnlock()
    return app.config
}

func (app *Application) updateConfig(newConfig *AppConfig) {
    app.mutex.Lock()
    defer app.mutex.Unlock()
    app.config = newConfig
}
```

This specification provides a comprehensive understanding of how to use the `pkg/config` module effectively. For more examples and best practices, refer to the usage guides. 