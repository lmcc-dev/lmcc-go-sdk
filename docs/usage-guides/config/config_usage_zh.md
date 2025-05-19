# 配置管理 (`pkg/config`) 使用指南

本指南详细说明了如何在 `lmcc-go-sdk` 中使用 `pkg/config` 模块来管理应用程序配置。

## 1. 功能介绍

`pkg/config` 包提供了一个灵活的系统来管理应用程序配置。主要特性包括：

- **多种来源:** 从文件（YAML、TOML 等）、环境变量和在结构体标签中定义的默认值加载配置。
- **优先级顺序:** 按定义的优先级合并配置源（环境变量 > 配置文件 > 默认值）。
- **类型安全:** 利用 Go 结构体和 `mapstructure` 将配置解码为类型化字段。
- **嵌入:** 设计为通过 `sdkconfig.Config` 嵌入到特定于应用程序的配置结构体中。
- **结构体标签默认值:** 使用 `default:"..."` 标签直接在结构体标签中定义默认值。
- **热重载 (Hot Reload):** 可选地使用 `WithHotReload(true)` 选项监控配置文件的更改并自动重新加载。
- **变更回调 (Change Callbacks):** 允许注册回调函数，在配置通过热重载更新后执行自定义逻辑（例如，重新配置日志级别）。

## 2. 接入指引

1.  **定义您的配置结构体:** 在您的应用程序中创建一个结构体，嵌入 `sdkconfig.Config` 并添加自定义字段。
    `sdkconfig.Config` 结构自身包含了预定义的配置节，如 `Server`, `Log`, `Database`, `Tracing`, `Metrics`。其中 `Log` 节 (`sdkconfig.LogConfig`) 的详细字段如下：

    ```go
    package main

    import (
        sdkconfig "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
        // "time" // 如果您的自定义配置部分使用了 time.Duration
    )

    // LogConfig 的详细结构 (通常定义在 sdkconfig 包内部，此处列出以供参考)
    /*
    type LogConfig struct {
        Level             string   `mapstructure:"level" default:"info"`
        Format            string   `mapstructure:"format" default:"text"`
        EnableColor       bool     `mapstructure:"enableColor" default:"false"`
        Output            string   `mapstructure:"output" default:"stdout"` // 旧版单输出，建议使用 outputPaths
        OutputPaths       []string `mapstructure:"outputPaths"`
        ErrorOutput       string   `mapstructure:"errorOutput" default:"stderr"` // 旧版单错误输出，建议使用 errorOutputPaths
        ErrorOutputPaths  []string `mapstructure:"errorOutputPaths"`
        Filename          string   `mapstructure:"filename" default:"app.log"` // 用于轮转的特定文件名
        MaxSize           int      `mapstructure:"maxSize" default:"100"`
        MaxBackups        int      `mapstructure:"maxBackups" default:"5"`
        MaxAge            int      `mapstructure:"maxAge" default:"7"`
        Compress          bool     `mapstructure:"compress" default:"false"`
        DisableCaller     bool     `mapstructure:"disableCaller" default:"false"`
        DisableStacktrace bool     `mapstructure:"disableStacktrace" default:"false"`
        Development       bool     `mapstructure:"development" default:"false"`
        Name              string   `mapstructure:"name"`
        ContextKeys       []string `mapstructure:"contextKeys"`
    }
    */

    type CustomFeatureConfig struct {
    	APIKey    string `mapstructure:"apiKey"`
    	RateLimit int    `mapstructure:"rateLimit" default:"100"`
    	Enabled   bool   `mapstructure:"enabled" default:"false"`
    }

    type MyAppConfig struct {
    	sdkconfig.Config                 // 嵌入 SDK 基础配置 (包含 Server, Log, Database 等)
    	CustomFeature *CustomFeatureConfig `mapstructure:"customFeature"`
    }

    var AppConfig MyAppConfig
    ```
    您可以在您的 `config.yaml` 文件中配置 `log` 节下的上述所有字段，例如：
    ```yaml
    # config.yaml 片段
    log:
      level: "debug"
      format: "json"
      enableColor: false # JSON 格式下通常不需要颜色
      outputPaths: ["stdout", "./app.log"]
      errorOutputPaths: ["stderr", "./app_error.log"]
      maxSize: 50
      # ... 其他 LogConfig 字段 ...
      disableCaller: false
      name: "my-awesome-app"
      contextKeys: ["user_id", "session_id"]
    ```

2.  **加载配置 (推荐: 使用 `LoadConfigAndWatch`)**: 在应用程序启动时调用 `sdkconfig.LoadConfigAndWatch`。这是推荐的方式，因为它提供了热重载和回调功能。

    ```go
    import (
    	"flag"
    	"fmt" // 引入 fmt
    	"log"
    	sdkconfig "github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
    	"github.com/spf13/viper" // 引入 viper 以便在回调中使用
    )

    func main() {
        configFile := flag.String("config", "config.yaml", "配置文件路径")
        flag.Parse()

        // 使用 LoadConfigAndWatch 加载配置，并选择性启用热重载
        // (Use LoadConfigAndWatch to load config and optionally enable hot reload)
        cm, err := sdkconfig.LoadConfigAndWatch(
            &AppConfig, // 指向您的结构体的指针
            sdkconfig.WithConfigFile(*configFile, ""), // 从扩展名推断类型
            sdkconfig.WithEnvPrefix("MYAPP"),          // 可选：更改环境变量前缀
            sdkconfig.WithHotReload(true),             // 启用热重载
            // sdkconfig.WithEnvVarOverride(true), // 默认启用环境变量覆盖
        )
        if err != nil {
            log.Fatalf("致命错误：加载配置失败: %v\n", err)
        }
        log.Println("配置加载成功。")

        // 访问配置管理器 (例如，注册回调)
        // (Access the config manager (e.g., to register callbacks))
        if cm != nil {
           // 在这里可以注册回调函数，参见下面的 "热重载与回调" 部分
           // (Callbacks can be registered here, see "Hot Reload and Callbacks" section below)
           log.Println("配置管理器已初始化，可以注册回调。")
        }


        // 稍后访问配置值
        // (Access config values later)
        fmt.Println("Server Port:", AppConfig.Server.Port)
        if AppConfig.CustomFeature != nil {
             fmt.Println("Custom Feature Enabled:", AppConfig.CustomFeature.Enabled)
        }
    }
    ```
    **注意:** 如果你**不需要**热重载和回调功能，可以使用简化的 `sdkconfig.LoadConfig(&AppConfig, ...)` 函数，它内部会调用 `LoadConfigAndWatch` 但会忽略 `WithHotReload(true)` 选项。

3.  **访问值:** 通过已填充的结构体实例 (`AppConfig`) 访问配置值。

    ```go
    port := AppConfig.Server.Port       // 访问 SDK 字段
    if AppConfig.CustomFeature != nil {
        apiKey := AppConfig.CustomFeature.APIKey // 访问自定义字段 (注意检查指针是否为 nil)
        isEnabled := AppConfig.CustomFeature.Enabled
        fmt.Printf("Custom Feature: APIKey=%s, Enabled=%t\n", apiKey, isEnabled)
    }
    ```

## 3. 热重载与回调 (Hot Reload and Callbacks)

如果你在调用 `LoadConfigAndWatch` 时使用了 `WithHotReload(true)`，当指定的配置文件发生更改时，配置会自动重新加载。你可以通过返回的 `configManager` 注册回调函数来响应这些更改。

```go
// ... 在 LoadConfigAndWatch 成功调用之后 ...
if cm != nil {
    err := cm.RegisterCallback(func(v *viper.Viper, currentCfgAny any) error {
        // 将 currentCfg 断言回你的配置类型
        // (Assert currentCfg back to your config type)
        currentCfg, ok := currentCfgAny.(*MyAppConfig)
        if !ok {
            // 这理论上不应发生，因为 LoadConfigAndWatch 保证了类型
            // (This should theoretically not happen as LoadConfigAndWatch ensures the type)
            return fmt.Errorf("回调中配置类型断言失败")
        }

        log.Printf("配置已重新加载！新端口: %d, 自定义功能启用状态: %t\n",
                   currentCfg.Server.Port, currentCfg.CustomFeature.Enabled)

        // 在这里根据新的配置执行操作，例如：
        // (Perform actions based on the new config here, e.g.:)
        // reconfigureLoggingLevel(currentCfg.Log.Level)
        // updateRateLimiter(currentCfg.CustomFeature.RateLimit)

        // 如果回调处理成功，返回 nil
        // (Return nil if the callback handles the change successfully)
        return nil
    })
    if err != nil {
        // 通常 RegisterCallback 不会返回错误，除非内部逻辑问题
        // (Usually RegisterCallback does not return error unless internal logic issue)
         log.Printf("警告：注册回调函数失败: %v", err)
    }
}
```
回调函数接收两个参数：
- `v *viper.Viper`: 当前的 Viper 实例。
- `cfg any`: 指向已更新的配置结构体的指针 (类型为 `any`，需要进行类型断言)。

## 4. API 参考

### 4.1. 核心函数

-   **`LoadConfigAndWatch[T any](cfg *T, opts ...Option) (*configManager[T], error)`**:
    -   **推荐使用。** 加载配置，并可选地启动文件监控和热重载。
    -   `cfg`: 指向您应用程序配置结构体的指针。
    -   `opts`: 用于控制加载行为的函数式选项。
    -   **返回:** 配置管理器实例 (`*configManager[T]`) 和错误 (`error`)。管理器可用于注册回调。
-   **`LoadConfig[T any](cfg *T, opts ...Option) error`**:
    -   加载配置，但不启动文件监控（会忽略 `WithHotReload(true)` 选项）。
    -   参数与 `LoadConfigAndWatch` 相同。
    -   **返回:** 错误 (`error`)。

### 4.2. 选项 (`config.Option`)

使用这些函数来自定义 `LoadConfig` 或 `LoadConfigAndWatch`：

-   `WithConfigFile(filePath string, fileType string) Option`: 指定配置文件的路径。如果 `fileType` 为空字符串（例如 `""`），它会尝试从文件扩展名推断类型。支持的类型包括 "yaml", "toml", "json" 等。
-   `WithEnvPrefix(prefix string) Option`: 设置环境变量的前缀（默认为 "LMCC"）。环境变量名称构造为 `PREFIX_SECTION_KEY`（例如 `LMCC_SERVER_PORT`）。
-   `WithEnvVarOverride(enable bool) Option`: 控制是否允许环境变量覆盖配置文件或默认值 (默认为 `true`)。
-   `WithHotReload(enable bool) Option`: 启用或禁用配置文件的热重载功能 (默认为 `false`)。如果启用 (`true`)，需要同时使用 `WithConfigFile` 指定要监控的文件。

### 4.3. 结构体标签

-   **`mapstructure:"<key_name>"`**: **必需**，用于您希望从配置文件或环境变量加载的字段。定义文件中使用的键名，并构成环境变量名称的基础。
-   **`default:"<string_value>"`**: **可选**。如果配置文件或环境变量中未找到该键，则指定默认值（以字符串形式）。该字符串将被解析为字段的类型（`int`、`bool`、`string`、`time.Duration` 等）。

## 5. 相关的 Makefile 命令

以下通用命令与此包的开发和测试相关：

-   `make test-unit PKG=./pkg/config`: 运行 `pkg/config` 的单元测试。
-   `make cover PKG=./pkg/config`: 对 `pkg/config` 运行带覆盖率分析的单元测试。
-   `make lint`: 运行 linter，其中包括对 `pkg/config` 的检查。
-   `make format`: 格式化 `pkg/config` 中的 Go 代码。
