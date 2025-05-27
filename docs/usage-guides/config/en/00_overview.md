# Configuration Module Overview

The `pkg/config` module provides a comprehensive configuration management solution for Go applications. It's designed to handle complex configuration scenarios while maintaining simplicity for basic use cases.

## Architecture

The configuration module is built around several key components:

### Core Components

1. **Manager Interface**: Central configuration management interface
2. **Options System**: Flexible configuration options using the functional options pattern
3. **Callback System**: Event-driven configuration change notifications
4. **Type System**: Predefined configuration structures for common use cases

### Configuration Flow

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Config File   │    │  Environment     │    │   Defaults      │
│   (YAML/JSON)   │    │   Variables      │    │  (Struct Tags)  │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────▼──────────────┐
                    │        Viper Core         │
                    │    (Configuration         │
                    │     Aggregation)          │
                    └─────────────┬──────────────┘
                                 │
                    ┌─────────────▼──────────────┐
                    │     mapstructure          │
                    │    (Type Conversion)      │
                    └─────────────┬──────────────┘
                                 │
                    ┌─────────────▼──────────────┐
                    │    User Config Struct     │
                    │   (Strongly Typed)        │
                    └───────────────────────────┘
```

## Key Features

### 1. Multiple Configuration Sources

The module supports loading configuration from multiple sources with a clear precedence order:

1. **Command-line flags** (highest priority)
2. **Environment variables**
3. **Configuration files**
4. **Default values** (lowest priority)

### 2. Hot Reload Support

Configuration files can be monitored for changes, and the application can be notified automatically when updates occur.

### 3. Type Safety

All configuration is strongly typed through user-defined structs, providing compile-time safety and IDE support.

### 4. Flexible Options

The module uses the functional options pattern, allowing for clean and extensible API design.

## Supported File Formats

- **YAML** (recommended)
- **JSON**
- **TOML**
- **HCL**
- **INI**
- **Properties**

## Environment Variable Binding

Automatic binding of environment variables to configuration fields with:
- Configurable prefixes
- Automatic key transformation (dots to underscores)
- Type conversion

## Default Value System

Default values can be specified using struct tags:

```go
type Config struct {
    Port     int    `mapstructure:"port" default:"8080"`
    Host     string `mapstructure:"host" default:"localhost"`
    Debug    bool   `mapstructure:"debug" default:"false"`
    Timeout  string `mapstructure:"timeout" default:"30s"`
}
```

## Error Handling

The module provides comprehensive error handling with specific error codes for different failure scenarios:

- Configuration file not found
- Invalid configuration format
- Type conversion errors
- Validation failures

## Thread Safety

All operations are thread-safe, allowing for concurrent access to configuration data and safe hot-reload operations.

## Integration Points

The configuration module is designed to integrate seamlessly with:

- **Logging module** (`pkg/log`) - for configuration-driven log settings
- **Error handling** (`pkg/errors`) - for consistent error reporting
- **Application frameworks** - through callback mechanisms

## Next Steps

- [Quick Start Guide](01_quick_start.md) - Get started with basic configuration
- [Configuration Options](02_configuration_options.md) - Learn about all available options
- [Hot Reload](03_hot_reload.md) - Implement dynamic configuration updates