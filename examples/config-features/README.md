/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This documentation was collaboratively developed by Martin and AI Assistant.
 */

# Configuration Features Examples

[中文版本](README_zh.md)

This directory contains examples demonstrating various features of the configuration module.

## Examples Overview

### [01-simple-config](01-simple-config/)
**Basic Configuration Loading**
- Simple YAML configuration file loading
- Struct tag defaults
- Basic validation
- Error handling

**Learn**: How to get started with configuration loading

### [02-hot-reload](02-hot-reload/)
**Hot-Reload Configuration**
- Real-time configuration updates
- File change monitoring
- Callback registration
- Graceful updates

**Learn**: How to implement dynamic configuration changes

### [03-env-override](03-env-override/)
**Environment Variable Override**
- Environment variable precedence
- Prefix configuration
- Type conversion
- Production deployment patterns

**Learn**: How to use environment variables for deployment flexibility

### [04-default-values](04-default-values/)
**Default Values Demonstration**
- Struct tag defaults
- Nested struct defaults
- Pointer field handling
- Zero value vs explicit values

**Learn**: How to implement robust default value systems

### [05-multiple-formats](05-multiple-formats/)
**Multiple File Format Support**
- YAML configuration
- JSON configuration
- TOML configuration
- Format auto-detection

**Learn**: How to support different configuration file formats

## Running the Examples

Each example is self-contained and can be run independently:

```bash
# Navigate to specific example
cd examples/config-features/01-simple-config

# Run the example
go run main.go

# Some examples support additional flags
go run main.go --help
```

## Common Patterns Demonstrated

### 1. Configuration Structure Design
```go
type AppConfig struct {
    config.Config                    // Embed SDK base config
    App    *AppSpecificConfig       `mapstructure:"app"`
    Feature *FeatureConfig         `mapstructure:"feature"`
}
```

### 2. Loading with Options
```go
err := config.LoadConfig(&cfg,
    config.WithConfigFile("config.yaml", "yaml"),
    config.WithEnvPrefix("MYAPP"),
    config.WithEnvVarOverride(true),
)
```

### 3. Hot-Reload Setup
```go
cm, err := config.LoadConfigAndWatch(&cfg,
    config.WithConfigFile("config.yaml", "yaml"),
    config.WithHotReload(true),
)

cm.RegisterCallback(func(v *viper.Viper, cfg any) error {
    // Handle configuration changes
    return nil
})
```

### 4. Environment Variable Patterns
```bash
# Format: {PREFIX}_{SECTION}_{FIELD}
export MYAPP_SERVER_PORT=8080
export MYAPP_LOG_LEVEL=debug
export MYAPP_DATABASE_HOST=prod-db.example.com
```

## Best Practices Shown

1. **Structure Organization**
   - Logical grouping of related settings
   - Clear separation between base and app-specific config
   - Consistent naming conventions

2. **Default Value Strategy**
   - Sensible defaults for development
   - Production-ready configurations
   - Graceful fallback mechanisms

3. **Environment Variable Usage**
   - Consistent prefix usage
   - Clear environment-specific overrides
   - Secure handling of sensitive data

4. **Error Handling**
   - Comprehensive error checking
   - Meaningful error messages
   - Graceful degradation

5. **Hot-Reload Considerations**
   - Thread-safe configuration updates
   - Validation before applying changes
   - Rollback mechanisms for invalid configs

## Integration Tips

### With Logging Module
```go
// Convert config to log options
logOpts := &log.Options{
    Level:  cfg.Log.Level,
    Format: cfg.Log.Format,
    // ... other mappings
}
log.Init(logOpts)
```

### With Error Handling
```go
if err := config.LoadConfig(&cfg, opts...); err != nil {
    if coder := errors.GetCoder(err); coder != nil {
        log.Errorf("Config error [%d]: %s", coder.Code(), coder.String())
    }
    return errors.Wrap(err, "failed to initialize application")
}
```

## Next Steps

After exploring these configuration examples:

1. Try the [error-handling examples](../error-handling/) to learn about robust error management
2. Explore [logging-features examples](../logging-features/) for comprehensive logging
3. See [integration examples](../integration/) for real-world usage patterns

## Troubleshooting

### Common Issues

**Issue**: Configuration file not found
```
Solution: Check file path and working directory
```

**Issue**: Environment variables not taking effect
```
Solution: Verify prefix format and restart application
```

**Issue**: Hot-reload not working
```
Solution: Ensure file permissions and check callback registration
```

**Issue**: Default values not applied
```
Solution: Verify struct tags and field types
```

## Related Documentation

- [Configuration Module Overview](../../docs/usage-guides/config/en/00_overview.md)
- [Quick Start Guide](../../docs/usage-guides/config/en/01_quick_start.md)
- [Configuration Options](../../docs/usage-guides/config/en/02_configuration_options.md)
- [Hot Reload Guide](../../docs/usage-guides/config/en/03_hot_reload.md) 