/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

/*
Package config provides flexible and robust configuration management capabilities for Go applications,
inspired by best practices seen in ecosystems like Marmotedu.
(config 包为 Go 应用程序提供了灵活且健壮的配置管理功能，其设计借鉴了 Marmotedu 等生态系统中的最佳实践。)

It leverages the Viper library for handling various configuration sources such as files (YAML, JSON, TOML, etc.),
environment variables, command-line flags, and default values defined via struct tags.
(它利用 Viper 库来处理各种配置源，例如文件（YAML、JSON、TOML 等）、环境变量、命令行标志以及通过结构体标签定义的默认值。)

Key features include:
(主要功能包括：)

  - Loading configuration from multiple sources with a clear precedence order.
    (从多个来源加载配置，具有明确的优先级顺序。)
  - Setting default values using struct field tags (`default:"value"`).
    (使用结构体字段标签 (`default:"value"`) 设置默认值。)
  - Automatic binding of environment variables to struct fields (respecting prefixes and `mapstructure` tags).
    (自动将环境变量绑定到结构体字段（遵循前缀和 `mapstructure` 标签）。)
  - Optional hot-reloading of configuration files upon changes, allowing dynamic reconfiguration.
    (可选的配置文件变更时热重载，允许动态重新配置。)
  - A callback mechanism to notify application components about configuration changes during hot-reloading,
    enabling graceful updates.
    (回调机制，用于在热重载期间通知应用程序组件配置变更，从而实现平滑更新。)
  - Type-safe access to configuration values through user-defined structs.
    (通过用户定义的结构体进行类型安全的配置值访问。)
  - Designed with testability in mind, allowing for easier mocking or provision of test-specific configurations.
    (设计时考虑了可测试性，可以更轻松地模拟或提供特定于测试的配置。)

Error Handling:
(错误处理：)

It is crucial to meticulously handle errors returned by `LoadConfigAndWatch` and within the registered callback functions.
Proper error handling ensures the application behaves predictably and can recover gracefully from configuration issues.
(细致处理 `LoadConfigAndWatch` 及注册的回调函数返回的错误至关重要。正确的错误处理可确保应用程序行为可预测，并能从配置问题中平稳恢复。)

Basic Usage:
(基本用法：)

	type MyConfig struct {
		Server struct {
			Host string `mapstructure:"host" default:"localhost"`
			Port int    `mapstructure:"port" default:"8080"`
		} `mapstructure:"server"`
		DatabaseURL string `mapstructure:"database_url" default:"postgres://user:pass@host/db"`
		Debug       bool   `mapstructure:"debug" default:"false"`
	}

	var cfg MyConfig

	// Load config with defaults, file, env vars, and hot-reload
	// (加载配置，包含默认值、文件、环境变量和热重载)
	// Consider using a more specific path or allowing it to be configurable for production.
	// (在生产环境中，考虑使用更具体的路径或允许其可配置。)
	cm, err := config.LoadConfigAndWatch(
		&cfg,
		config.WithConfigFile("config/app.yaml"), // Specify config file (指定配置文件)
		config.WithEnvPrefix("APP"),             // Env var prefix (e.g., APP_SERVER_PORT) (环境变量前缀)
		config.WithHotReload(),                  // Enable hot-reload (启用热重载)
		// config.WithEnvKeyReplacer(strings.NewReplacer(".", "_")), // Optional: if env vars use '_' instead of '.'
	)
	if err != nil {
		// In a real application, consider more sophisticated error logging or a graceful shutdown.
		// (在实际应用中，考虑使用更完善的错误日志记录或优雅停机。)
		log.Fatalf("FATAL: Failed to load initial configuration: %v", err)
	}

	// Register a callback for changes (注册变更回调)
	// The callback ID allows for unregistering if needed, though not shown in this basic example.
	// (回调 ID 允许在需要时注销回调，尽管此基本示例中未展示。)
	callbackID := "myAppConfigUpdater"
	cm.RegisterCallback(callbackID, func(v *viper.Viper, currentCfgPtr any) error {
		appCfg, ok := currentCfgPtr.(*MyConfig) // Type assertion (类型断言)
		if !ok {
			// This should ideally not happen if LoadConfigAndWatch was called with *MyConfig
			// (如果 LoadConfigAndWatch 是用 *MyConfig 调用的，这理想情况下不应发生)
			log.Printf("ERROR: Config callback received unexpected type: %T", currentCfgPtr)
			return fmt.Errorf("config callback: unexpected type %T", currentCfgPtr)
		}
		log.Printf("INFO: Configuration reloaded! New Port: %d, Debug Mode: %t", appCfg.Server.Port, appCfg.Debug)

		// IMPORTANT: Update application components based on new config.
		// (重要：根据新配置更新应用组件。)
		// Ensure this logic is thread-safe if components are accessed concurrently.
		// (如果组件被并发访问，请确保此逻辑是线程安全的。)
		// Example: updateLogLevel(appCfg.LogLevel)
		// Example: reinitializeDatabaseConnection(appCfg.DatabaseURL)

		// Return nil if successful, or an error to indicate failure in processing the update.
		// (如果成功则返回 nil，或返回一个错误以指示处理更新失败。)
		// If an error is returned, the hot-reload might be logged as failed, but Viper continues watching.
		// (如果返回错误，热重载可能会被记录为失败，但 Viper 会继续监视。)
		return nil
	})
	if err != nil { // Check error from RegisterCallback, though in current implementation it's nil
		log.Printf("WARN: Failed to register config change callback '%s': %v", callbackID, err)
	}


	// Access configuration values (访问配置值)
	// These values reflect the initially loaded or most recently reloaded configuration.
	// (这些值反映了初始加载或最近一次热重载的配置。)
	fmt.Println("Initial/Current Server Port:", cfg.Server.Port)
	fmt.Println("Initial/Current Debug Mode:", cfg.Debug)

	// Application continues to run, configuration may be hot-reloaded in the background...
	// (应用程序继续运行，配置可能会在后台热重载...)
	// select {} // Block main goroutine to observe hot-reloading for testing
*/
package config