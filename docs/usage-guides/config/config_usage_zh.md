# 配置管理 (`pkg/config`) 使用指南

本指南详细说明了如何在 `lmcc-go-sdk` 中使用 `pkg/config` 模块来管理应用程序配置。

## 1. 功能介绍

`pkg/config` 包提供了一个灵活的系统来管理应用程序配置。主要特性包括：

- **多种来源:** 从文件（YAML、TOML 等）、环境变量和在结构体标签中定义的默认值加载配置。
- **优先级顺序:** 按定义的优先级合并配置源（环境变量 > 配置文件 > 默认值）。
- **类型安全:** 利用 Go 结构体和 `mapstructure` 将配置解码为类型化字段。
- **嵌入:** 设计为通过 `sdkconfig.Config` 嵌入到特定于应用程序的配置结构体中。
- **结构体标签默认值:** 使用 `default:"..."` 标签直接在结构体标签中定义默认值。
- **热加载:** 可选地使用 `WithHotReload()` 选项监控配置文件的更改并自动重新加载。

## 2. 接入指引

1.  **定义您的配置结构体:** 在您的应用程序中创建一个结构体，嵌入 `sdkconfig.Config` 并添加自定义字段。

    ```go
    package main

    import sdkconfig "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"

    type CustomFeatureConfig struct {
    	APIKey    string `mapstructure:"apiKey"`
    	RateLimit int    `mapstructure:"rateLimit" default:"100"`
    	Enabled   bool   `mapstructure:"enabled" default:"false"`
    }

    type MyAppConfig struct {
    	sdkconfig.Config                 // 嵌入 SDK 基础配置
    	CustomFeature *CustomFeatureConfig `mapstructure:"customFeature"`
    }

    var AppConfig MyAppConfig
    ```

2.  **加载配置:** 在应用程序启动时调用 `sdkconfig.LoadConfig`。

    ```go
    import (
    	"flag"
    	"log"
    	sdkconfig "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
    	"strings"
    )

    func main() {
        configFile := flag.String("config", "config.yaml", "配置文件路径")
        flag.Parse()

        err := sdkconfig.LoadConfig(
            &AppConfig, // 指向您的结构体的指针
            sdkconfig.WithConfigFile(*configFile, ""), // 从扩展名推断类型
            // sdkconfig.WithEnvPrefix("MYAPP"), // 可选：更改环境变量前缀
            // sdkconfig.WithHotReload(),       // 可选：启用热加载
        )
        if err != nil {
            log.Fatalf("致命错误：加载配置失败: %v\n", err)
        }
        log.Println("配置加载成功。")
    }
    ```

3.  **访问值:** 通过已填充的结构体实例访问配置值。

    ```go
    port := AppConfig.Server.Port       // 访问 SDK 字段
    apiKey := AppConfig.CustomFeature.APIKey // 访问自定义字段
    ```

## 3. API 参考

### 3.1. 核心函数

-   **`LoadConfig[T any](cfg *T, opts ...Option) error`**:
    -   根据提供的选项将配置加载到 `cfg` 指向的结构体中。
    -   `cfg`: 指向您应用程序配置结构体的指针（必须嵌入 `sdkconfig.Config` 或定义兼容字段）。
    -   `opts`: 用于控制加载行为的函数式选项。

### 3.2. 选项 (`config.Option`)

使用这些函数来自定义 `LoadConfig`：

-   `WithConfigFile(filePath string, fileType string) Option`: 指定配置文件的路径。如果 `fileType` 为空字符串（例如 `""`），它会尝试从文件扩展名推断类型。支持的类型包括 "yaml", "toml", "json" 等（由 Viper 处理）。
-   `WithEnvPrefix(prefix string) Option`: 设置环境变量的前缀（默认为 "LMCC"）。环境变量名称构造为 `PREFIX_SECTION_KEY`（例如 `LMCC_SERVER_PORT`）。
-   `WithEnvKeyReplacer(replacer *strings.Replacer) Option`: 提供一个自定义的 `strings.Replacer`，用于将结构体字段键（经过 `mapstructure` 标签查找后）转换为环境变量段（默认将 `.` 和 `-` 替换为 `_`）。
-   `WithoutEnvVarOverride() Option`: 禁用从环境变量加载配置设置。
-   `WithHotReload() Option`: 启用监控指定的配置文件以检测更改，并自动将配置重新加载到 `cfg` 结构体中。需要同时使用 `WithConfigFile`。

### 3.3. 结构体标签

-   **`mapstructure:"<key_name>"`**: **必需**，用于您希望从配置文件或环境变量加载的字段。定义文件中使用的键名（例如 YAML 键），并构成环境变量名称的基础。
-   **`default:"<string_value>"`**: **可选**。如果配置文件或环境变量中未找到该键，则指定默认值（以字符串形式）。该字符串将被解析为字段的类型（`int`、`bool`、`string`、`time.Duration` 等）。

## 4. 相关的 Makefile 命令

虽然没有 *仅* 针对 `pkg/config` 的 Makefile 命令，但以下通用命令与此包的开发和测试相关：

-   `make test`: 运行 `pkg/config` 的单元测试。
-   `make cover`: 对 `pkg/config` 运行带覆盖率分析的单元测试。
-   `make lint`: 运行 linter，其中包括对 `pkg/config` 的检查。
-   `make format`: 格式化 `pkg/config` 中的 Go 代码。 