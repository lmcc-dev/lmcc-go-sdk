# Troubleshooting

This guide helps you diagnose and resolve common issues with the configuration module.

## Common Issues

### 1. Configuration File Not Found

**Problem:** Application fails to start with "configuration file not found" error.

**Symptoms:**
```
Error: failed to load configuration: config file not found
```

**Solutions:**

1. **Check file path:**
   ```go
   // Make sure the file exists at the specified path
   err := config.LoadConfig(&cfg, 
       config.WithConfigFile("config.yaml", ""), // Check this path
   )
   ```

2. **Use absolute path:**
   ```go
   err := config.LoadConfig(&cfg, 
       config.WithConfigFile("/app/config/config.yaml", ""),
   )
   ```

3. **Add search paths:**
   ```go
   err := config.LoadConfig(&cfg, 
       config.WithConfigFile("config.yaml", "/etc/myapp:/app/config:./config"),
   )
   ```

4. **Check file permissions:**
   ```bash
   ls -la config.yaml
   # Should be readable by the application user
   chmod 644 config.yaml
   ```

### 2. Environment Variables Not Working

**Problem:** Environment variables are not overriding configuration values.

**Symptoms:**
- Configuration values don't change despite setting environment variables
- Environment variables are ignored

**Solutions:**

1. **Enable environment variable override:**
   ```go
   err := config.LoadConfig(&cfg,
       config.WithConfigFile("config.yaml", ""),
       config.WithEnvPrefix("APP"),
       config.WithEnvVarOverride(true), // Make sure this is enabled
   )
   ```

2. **Check environment variable naming:**
   ```go
   // For this configuration structure:
   type Config struct {
       Server struct {
           Host string `mapstructure:"host"`
           Port int    `mapstructure:"port"`
       } `mapstructure:"server"`
   }
   
   // Environment variables should be:
   // APP_SERVER_HOST=localhost
   // APP_SERVER_PORT=8080
   ```

3. **Verify environment variables are set:**
   ```bash
   env | grep APP_
   echo $APP_SERVER_HOST
   ```

4. **Use custom key replacer if needed:**
   ```go
   import "strings"
   
   err := config.LoadConfig(&cfg,
       config.WithEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_")),
   )
   ```

### 3. Hot Reload Not Working

**Problem:** Configuration changes are not detected or applied.

**Symptoms:**
- File changes don't trigger callbacks
- Application doesn't respond to configuration updates

**Solutions:**

1. **Ensure hot reload is enabled:**
   ```go
   cm, err := config.LoadConfigAndWatch(&cfg,
       config.WithConfigFile("config.yaml", ""),
       config.WithHotReload(true), // Must be true
   )
   ```

2. **Check file system permissions:**
   ```bash
   # Application needs read access to config file and directory
   ls -la config.yaml
   ls -la $(dirname config.yaml)
   ```

3. **Verify callbacks are registered:**
   ```go
   cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
       log.Printf("Configuration changed") // Add logging
       return nil
   })
   ```

4. **Check for file system watcher issues:**
   ```go
   // Enable debug logging to see watcher events
   import "github.com/spf13/viper"
   viper.Debug()
   ```

5. **Handle editor-specific issues:**
   ```go
   // Some editors create temporary files that can interfere
   // Consider using atomic file operations or ignoring temp files
   ```

### 4. Configuration Validation Errors

**Problem:** Configuration fails validation with unclear error messages.

**Symptoms:**
- Validation errors without clear indication of the problem
- Application fails to start due to invalid configuration

**Solutions:**

1. **Add detailed validation:**
   ```go
   func (c *Config) Validate() error {
       var errors []error
       
       if c.Server.Port < 1 || c.Server.Port > 65535 {
           errors = append(errors, fmt.Errorf("server.port must be between 1 and 65535, got %d", c.Server.Port))
       }
       
       if c.Database.URL == "" {
           errors = append(errors, fmt.Errorf("database.url is required"))
       }
       
       if len(errors) > 0 {
           return fmt.Errorf("configuration validation failed:\n%s", 
               strings.Join(errorStrings(errors), "\n"))
       }
       
       return nil
   }
   
   func errorStrings(errors []error) []string {
       var strs []string
       for _, err := range errors {
           strs = append(strs, "  - "+err.Error())
       }
       return strs
   }
   ```

2. **Use struct tags for validation:**
   ```go
   import "github.com/go-playground/validator/v10"
   
   type Config struct {
       Server struct {
           Host string `mapstructure:"host" validate:"required,hostname"`
           Port int    `mapstructure:"port" validate:"required,min=1,max=65535"`
       } `mapstructure:"server" validate:"required"`
   }
   
   func validateConfig(cfg *Config) error {
       validate := validator.New()
       return validate.Struct(cfg)
   }
   ```

### 5. Default Values Overriding Configuration File Values

**Problem:** Configuration values from files are being overridden by default values defined in struct tags.

**Symptoms:**
- Boolean `false` values in config files are changed to `true` if the default is `true`
- Zero values (0, false, "") in config files are replaced by struct tag defaults
- Configuration file appears to be ignored for certain fields

**Example:**
```yaml
# config.yaml
enableMetrics: false  # This gets changed to true
```

```go
type Config struct {
    EnableMetrics bool `mapstructure:"enableMetrics" default:"true"`
}
```

**Root Cause:**
The default value application logic cannot distinguish between "user explicitly set to false in config file" and "field was not set (zero value)".

**Solution:**
This issue has been fixed in the latest version through improved default value handling that:
- Records which keys actually exist in the configuration file
- Only applies defaults to fields that are truly unset
- Properly handles zero values that are explicitly set by users

**Workaround for older versions:**
Use environment variables to override problematic boolean fields:
```bash
export APP_ENABLE_METRICS=false
```

### 6. Type Conversion Issues

**Problem:** Configuration values are not properly converted to expected types.

**Symptoms:**
- String values when integers are expected
- Boolean values not recognized
- Duration parsing errors

**Solutions:**

1. **Check mapstructure tags:**
   ```go
   type Config struct {
       // Correct
       Port    int           `mapstructure:"port"`
       Debug   bool          `mapstructure:"debug"`
       Timeout time.Duration `mapstructure:"timeout"`
       
       // Incorrect - missing mapstructure tag
       // Host string
   }
   ```

2. **Use proper YAML syntax:**
   ```yaml
   # Correct
   server:
     port: 8080        # Integer
     debug: true       # Boolean
     timeout: "30s"    # Duration string
   
   # Incorrect
   server:
     port: "8080"      # String instead of integer
     debug: "true"     # String instead of boolean
   ```

3. **Handle custom types:**
   ```go
   type LogLevel string
   
   func (l *LogLevel) UnmarshalText(text []byte) error {
       level := string(text)
       switch level {
       case "debug", "info", "warn", "error":
           *l = LogLevel(level)
           return nil
       default:
           return fmt.Errorf("invalid log level: %s", level)
       }
   }
   ```

### 7. Memory Leaks with Hot Reload

**Problem:** Application memory usage increases over time with hot reload enabled.

**Symptoms:**
- Gradual memory increase
- Performance degradation over time

**Solutions:**

1. **Properly clean up resources in callbacks:**
   ```go
   cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
       // Close old connections before creating new ones
       if oldDB != nil {
           oldDB.Close()
       }
       
       // Create new connections
       newDB, err := createDatabaseConnection(currentCfg)
       if err != nil {
           return err
       }
       
       oldDB = newDB
       return nil
   })
   ```

2. **Use context for cancellation:**
   ```go
   type Server struct {
       ctx    context.Context
       cancel context.CancelFunc
   }
   
   func (s *Server) reconfigure(cfg *Config) error {
       // Cancel previous context
       if s.cancel != nil {
           s.cancel()
       }
       
       // Create new context
       s.ctx, s.cancel = context.WithCancel(context.Background())
       
       // Start new services with new context
       return s.startServices(s.ctx, cfg)
   }
   ```

## Debugging Techniques

### 1. Enable Debug Logging

```go
import (
    "log"
    "github.com/spf13/viper"
)

func enableDebugLogging() {
    // Enable viper debug logging
    viper.Debug()
    
    // Enable application debug logging
    log.SetFlags(log.LstdFlags | log.Lshortfile)
}
```

### 2. Log Configuration Values

```go
func logConfiguration(cfg *Config) {
    // Create a safe copy for logging (remove sensitive data)
    safeCfg := *cfg
    safeCfg.Database.Password = "[REDACTED]"
    safeCfg.Security.JWTSecret = "[REDACTED]"
    
    log.Printf("Loaded configuration: %+v", safeCfg)
}
```

### 3. Validate Configuration Step by Step

```go
func debugConfigurationLoading() {
    var cfg Config
    
    log.Printf("Step 1: Loading configuration file...")
    err := config.LoadConfig(&cfg, config.WithConfigFile("config.yaml", ""))
    if err != nil {
        log.Printf("File loading failed: %v", err)
        return
    }
    log.Printf("File loaded successfully")
    
    log.Printf("Step 2: Applying environment variables...")
    err = config.LoadConfig(&cfg, 
        config.WithConfigFile("config.yaml", ""),
        config.WithEnvPrefix("APP"),
        config.WithEnvVarOverride(true),
    )
    if err != nil {
        log.Printf("Environment variable application failed: %v", err)
        return
    }
    log.Printf("Environment variables applied")
    
    log.Printf("Step 3: Validating configuration...")
    if err := cfg.Validate(); err != nil {
        log.Printf("Validation failed: %v", err)
        return
    }
    log.Printf("Configuration validated successfully")
}
```

### 4. Monitor File System Events

```go
import (
    "github.com/fsnotify/fsnotify"
    "log"
)

func monitorConfigFile(filename string) {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watcher.Close()
    
    err = watcher.Add(filename)
    if err != nil {
        log.Fatal(err)
    }
    
    for {
        select {
        case event, ok := <-watcher.Events:
            if !ok {
                return
            }
            log.Printf("File event: %s %s", event.Name, event.Op)
            
        case err, ok := <-watcher.Errors:
            if !ok {
                return
            }
            log.Printf("File watcher error: %v", err)
        }
    }
}
```

## Performance Issues

### 1. Slow Configuration Loading

**Problem:** Configuration loading takes too long.

**Solutions:**

1. **Optimize file I/O:**
   ```go
   // Use local files instead of network-mounted filesystems
   // Consider caching configuration in memory
   ```

2. **Reduce configuration complexity:**
   ```go
   // Split large configuration files into smaller, focused files
   // Load only necessary configuration sections
   ```

3. **Profile configuration loading:**
   ```go
   import (
       "time"
       "log"
   )
   
   func loadConfigWithProfiling() {
       start := time.Now()
       
       var cfg Config
       err := config.LoadConfig(&cfg, options...)
       
       duration := time.Since(start)
       log.Printf("Configuration loading took: %v", duration)
       
       if duration > time.Second {
           log.Printf("Warning: Configuration loading is slow")
       }
   }
   ```

### 2. High Memory Usage

**Problem:** Configuration module uses too much memory.

**Solutions:**

1. **Optimize configuration structure:**
   ```go
   // Use pointers for large nested structures
   type Config struct {
       LargeSection *LargeConfigSection `mapstructure:"large_section"`
   }
   ```

2. **Implement configuration lazy loading:**
   ```go
   type Config struct {
       loaded map[string]interface{}
       mu     sync.RWMutex
   }
   
   func (c *Config) GetSection(name string) interface{} {
       c.mu.RLock()
       if section, exists := c.loaded[name]; exists {
           c.mu.RUnlock()
           return section
       }
       c.mu.RUnlock()
       
       // Load section on demand
       c.mu.Lock()
       defer c.mu.Unlock()
       
       section := c.loadSection(name)
       c.loaded[name] = section
       return section
   }
   ```

## Error Messages Reference

### Common Error Patterns

| Error Message | Cause | Solution |
|---------------|-------|----------|
| `config file not found` | File doesn't exist or wrong path | Check file path and permissions |
| `yaml: unmarshal errors` | Invalid YAML syntax | Validate YAML syntax |
| `mapstructure: cannot decode` | Type mismatch | Check struct tags and types |
| `environment variable not found` | Missing env var with no default | Set environment variable or add default |
| `validation failed` | Configuration doesn't meet requirements | Fix configuration values |
| `permission denied` | Insufficient file permissions | Fix file permissions |
| `too many open files` | File descriptor leak | Close resources properly |

### Getting Help

1. **Enable verbose logging** to get more details about the issue
2. **Check the configuration file syntax** using online YAML validators
3. **Verify environment variables** are properly set and named
4. **Test with minimal configuration** to isolate the problem
5. **Review the documentation** for the specific feature you're using

## Next Steps

- [Quick Start Guide](01_quick_start.md) - Basic usage examples
- [Configuration Options](02_configuration_options.md) - All available options
- [Best Practices](04_best_practices.md) - Recommended patterns 