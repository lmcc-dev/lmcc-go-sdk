# Configuration Management (`pkg/config`) Usage Guide

This guide details how to use the `pkg/config` module for managing application configuration within the `lmcc-go-sdk`.

## 1. Feature Introduction

The `pkg/config` package provides a flexible system for managing application configuration. Key features include:

- **Multiple Sources:** Load configuration from files (YAML, TOML, etc.), environment variables, and default values defined in struct tags.
- **Priority Order:** Merges sources with a defined priority (Env Vars > Config File > Defaults).
- **Type Safety:** Leverages Go structs and `mapstructure` for decoding configuration into typed fields.
- **Embedding:** Designed to be embedded into application-specific configuration structs using `sdkconfig.Config`.
- **Struct Tag Defaults:** Define default values directly within struct tags using the `default:"..."` tag.
- **Hot Reload:** Optionally monitor configuration files for changes and reload automatically using the `WithHotReload(true)` option.
- **Change Callbacks:** Allows registering callback functions to execute custom logic (e.g., reconfiguring log levels) after the configuration is updated via hot reload.

## 2. Integration Guide

1.  **Define Your Config Struct:** Embed `sdkconfig.Config` into your application's config struct and add custom fields.

    ```go
    package main

    import sdkconfig "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"

    type CustomFeatureConfig struct {
    	APIKey    string `mapstructure:"apiKey"`
    	RateLimit int    `mapstructure:"rateLimit" default:"100"`
    	Enabled   bool   `mapstructure:"enabled" default:"false"`
    }

    type MyAppConfig struct {
    	sdkconfig.Config                 // Embed SDK base config
    	CustomFeature *CustomFeatureConfig `mapstructure:"customFeature"`
    }

    var AppConfig MyAppConfig
    ```

2.  **Load Configuration (Recommended: Use `LoadConfigAndWatch`)**: Use `sdkconfig.LoadConfigAndWatch` at application startup. This is the recommended approach as it provides hot-reload and callback capabilities.

    ```go
    import (
    	"flag"
    	"fmt" // Import fmt
    	"log"
    	sdkconfig "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
    	"github.com/spf13/viper" // Import viper for use in callbacks
    )

    func main() {
        configFile := flag.String("config", "config.yaml", "Path to configuration file")
        flag.Parse()

        // Use LoadConfigAndWatch to load config and optionally enable hot reload
        cm, err := sdkconfig.LoadConfigAndWatch(
            &AppConfig, // Pointer to your struct
            sdkconfig.WithConfigFile(*configFile, ""), // Infer type from extension
            sdkconfig.WithEnvPrefix("MYAPP"),          // Optional: Change env var prefix
            sdkconfig.WithHotReload(true),             // Enable hot reload
            // sdkconfig.WithEnvVarOverride(true), // Env var override is enabled by default
        )
        if err != nil {
            log.Fatalf("FATAL: Failed to load configuration: %v\n", err)
        }
        log.Println("Configuration loaded successfully.")

        // Access the config manager (e.g., to register callbacks)
        if cm != nil {
           // Callbacks can be registered here, see "Hot Reload and Callbacks" section below
           log.Println("Config manager initialized, callbacks can be registered.")
        }

        // Access config values later
        fmt.Println("Server Port:", AppConfig.Server.Port)
         if AppConfig.CustomFeature != nil {
             fmt.Println("Custom Feature Enabled:", AppConfig.CustomFeature.Enabled)
        }
    }
    ```
    **Note:** If you do **not** need hot-reload and callbacks, you can use the simplified `sdkconfig.LoadConfig(&AppConfig, ...)` function. It calls `LoadConfigAndWatch` internally but ignores the `WithHotReload(true)` option.

3.  **Access Values:** Access loaded values via the populated struct instance (`AppConfig`).

    ```go
    port := AppConfig.Server.Port       // Access SDK field
    if AppConfig.CustomFeature != nil {
        apiKey := AppConfig.CustomFeature.APIKey // Access custom field (check for nil pointer)
        isEnabled := AppConfig.CustomFeature.Enabled
        fmt.Printf("Custom Feature: APIKey=%s, Enabled=%t\n", apiKey, isEnabled)
    }
    ```

## 3. Hot Reload and Callbacks

If you used `WithHotReload(true)` when calling `LoadConfigAndWatch`, the configuration will be automatically reloaded when the specified config file changes. You can respond to these changes by registering callback functions using the returned `configManager`.

```go
// ... After LoadConfigAndWatch succeeds ...
if cm != nil {
    err := cm.RegisterCallback(func(v *viper.Viper, currentCfgAny any) error {
        // Assert currentCfg back to your config type
        currentCfg, ok := currentCfgAny.(*MyAppConfig)
        if !ok {
            // This should theoretically not happen as LoadConfigAndWatch ensures the type
            return fmt.Errorf("config type assertion failed in callback")
        }

        log.Printf("Configuration reloaded! New Port: %d, Custom Feature Enabled: %t\n",
                   currentCfg.Server.Port, currentCfg.CustomFeature.Enabled)

        // Perform actions based on the new config here, e.g.:
        // reconfigureLoggingLevel(currentCfg.Log.Level)
        // updateRateLimiter(currentCfg.CustomFeature.RateLimit)

        // Return nil if the callback handles the change successfully
        return nil
    })
    if err != nil {
         // Usually RegisterCallback does not return error unless internal logic issue
         log.Printf("Warning: Failed to register callback: %v", err)
    }
}
```
The callback function receives two arguments:
- `v *viper.Viper`: The current Viper instance.
- `cfg any`: A pointer to the updated configuration struct (type `any`, requires type assertion).

## 4. API Reference

### 4.1. Core Functions

-   **`LoadConfigAndWatch[T any](cfg *T, opts ...Option) (*configManager[T], error)`**:
    -   **Recommended.** Loads configuration and optionally starts file watching and hot-reloading.
    -   `cfg`: Pointer to your application's config struct.
    -   `opts`: Functional options to control loading behavior.
    -   **Returns:** A config manager instance (`*configManager[T]`) and an error (`error`). The manager can be used to register callbacks.
-   **`LoadConfig[T any](cfg *T, opts ...Option) error`**:
    -   Loads configuration without starting file watching (ignores `WithHotReload(true)` option).
    -   Parameters are the same as `LoadConfigAndWatch`.
    -   **Returns:** An error (`error`).

### 4.2. Options (`config.Option`)

Use these functions to customize `LoadConfig` or `LoadConfigAndWatch`:

-   `WithConfigFile(filePath string, fileType string) Option`: Specifies the configuration file path. If `fileType` is empty (e.g., `""`), it attempts to infer the type from the file extension. Supported types include "yaml", "toml", "json", etc.
-   `WithEnvPrefix(prefix string) Option`: Sets the prefix for environment variables (default: "LMCC"). Environment variable names are constructed as `PREFIX_SECTION_KEY` (e.g., `LMCC_SERVER_PORT`).
-   `WithEnvVarOverride(enable bool) Option`: Controls whether environment variables are allowed to override config file or default values (default: `true`).
-   `WithHotReload(enable bool) Option`: Enables or disables the hot-reload feature for the configuration file (default: `false`). Requires `WithConfigFile` to specify the file to watch if enabled (`true`).

### 4.3. Struct Tags

-   **`mapstructure:"<key_name>"`**: **Required** for fields you want to load from config files or environment variables. Defines the key name used in the file and forms the basis for the environment variable name.
-   **`default:"<string_value>"`**: **Optional**. Specifies the default value as a string if the key is not found in the config file or environment variables. The string will be parsed into the field's type (`int`, `bool`, `string`, `time.Duration`, etc.).

## 5. Relevant Makefile Commands

The following general commands are relevant for development and testing of this package:

-   `make test-unit PKG=./pkg/config`: Runs unit tests for `pkg/config`.
-   `make cover PKG=./pkg/config`: Runs unit tests with coverage analysis for `pkg/config`.
-   `make lint`: Runs linters, which include checks on `pkg/config`.
-   `make format`: Formats the Go code in `pkg/config`.
