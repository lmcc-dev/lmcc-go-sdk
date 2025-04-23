# Configuration Management (`pkg/config`) Usage Guide

This guide details how to use the `pkg/config` module for managing application configuration within the `lmcc-go-sdk`.

## 1. Feature Introduction

The `pkg/config` package provides a flexible system for managing application configuration. Key features include:

- **Multiple Sources:** Load configuration from files (YAML, TOML, etc.), environment variables, and default values defined in struct tags.
- **Priority Order:** Merges sources with a defined priority (Env Vars > Config File > Defaults).
- **Type Safety:** Leverages Go structs and `mapstructure` for decoding configuration into typed fields.
- **Embedding:** Designed to be embedded into application-specific configuration structs using `sdkconfig.Config`.
- **Struct Tag Defaults:** Define default values directly within struct tags using the `default:"..."` tag.
- **Hot Reloading:** Optionally monitor configuration files for changes and reload automatically using the `WithHotReload()` option.

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

2.  **Load Configuration:** Use `sdkconfig.LoadConfig` at application startup.

    ```go
    import (
    	"flag"
    	"log"
    	sdkconfig "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
    	"strings"
    )

    func main() {
        configFile := flag.String("config", "config.yaml", "Path to configuration file")
        flag.Parse()

        err := sdkconfig.LoadConfig(
            &AppConfig,
            sdkconfig.WithConfigFile(*configFile, ""), // Infer type from extension
            // sdkconfig.WithEnvPrefix("MYAPP"), // Optional
            // sdkconfig.WithHotReload(),       // Optional
        )
        if err != nil {
            log.Fatalf("FATAL: Failed to load configuration: %v\n", err)
        }
        log.Println("Configuration loaded successfully.")
    }
    ```

3.  **Access Values:** Access loaded values via the struct instance.

    ```go
    port := AppConfig.Server.Port       // Access SDK field
    apiKey := AppConfig.CustomFeature.APIKey // Access custom field
    ```

## 3. API Reference

### 3.1. Core Function

-   **`LoadConfig[T any](cfg *T, opts ...Option) error`**:
    -   Loads configuration based on the provided options into the struct pointed to by `cfg`.
    -   `cfg`: Pointer to your application's config struct (must embed `sdkconfig.Config` or define compatible fields).
    -   `opts`: Functional options to control loading behavior.

### 3.2. Options (`config.Option`)

Use these functions to customize `LoadConfig`:

-   `WithConfigFile(filePath string, fileType string) Option`: Specifies the configuration file path. If `fileType` is empty (e.g., ""), it attempts to infer the type from the file extension. Supported types include "yaml", "toml", "json", etc. (handled by Viper).
-   `WithEnvPrefix(prefix string) Option`: Sets the prefix for environment variables (default: "LMCC"). Environment variable names are constructed as `PREFIX_SECTION_KEY` (e.g., `LMCC_SERVER_PORT`).
-   `WithEnvKeyReplacer(replacer *strings.Replacer) Option`: Provides a custom `strings.Replacer` to translate struct field keys (after `mapstructure` tag lookup) into environment variable segments (default replaces `.` and `-` with `_`).
-   `WithoutEnvVarOverride() Option`: Disables loading configuration settings from environment variables.
-   `WithHotReload() Option`: Enables monitoring the specified configuration file for changes and automatically reloading the configuration into the `cfg` struct. Requires `WithConfigFile` to be used.

### 3.3. Struct Tags

-   **`mapstructure:"<key_name>"`**: **Required** for fields you want to load from config files or environment variables. Defines the key name used in the file (e.g., YAML key) and forms the basis for the environment variable name.
-   **`default:"<string_value>"`**: **Optional**. Specifies the default value as a string if the key is not found in the config file or environment variables. The string will be parsed into the field's type (`int`, `bool`, `string`, `time.Duration`, etc.).

## 4. Relevant Makefile Commands

While there are no Makefile commands *specific* only to `pkg/config`, the following general commands are relevant for development and testing of this package:

-   `make test`: Runs unit tests for `pkg/config`.
-   `make cover`: Runs unit tests with coverage analysis for `pkg/config`.
-   `make lint`: Runs linters, which include checks on `pkg/config`.
-   `make format`: Formats the Go code in `pkg/config`. 