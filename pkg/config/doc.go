/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

/*
Package config provides flexible configuration management capabilities for Go applications.
(config 包为 Go 应用程序提供了灵活的配置管理功能。)

It leverages the Viper library for handling various configuration sources like files (YAML, JSON, TOML, etc.),
environment variables, and default values defined via struct tags.
(它利用 Viper 库来处理各种配置源，例如文件（YAML、JSON、TOML 等）、环境变量以及通过结构体标签定义的默认值。)

Key features include:
(主要功能包括：)

  - Loading configuration from files and environment variables.
    (从文件和环境变量加载配置。)
  - Setting default values using struct field tags (`default:"value"`).
    (使用结构体字段标签 (`default:"value"`) 设置默认值。)
  - Automatic binding of environment variables to struct fields (respecting prefixes and `mapstructure` tags).
    (自动将环境变量绑定到结构体字段（遵循前缀和 `mapstructure` 标签）。)
  - Optional hot-reloading of configuration files upon changes.
    (可选的配置文件变更时热重载。)
  - Callback mechanism to notify application components about configuration changes during hot-reloading.
    (回调机制，用于在热重载期间通知应用程序组件配置变更。)
  - Type-safe access to configuration values through user-defined structs.
    (通过用户定义的结构体进行类型安全的配置值访问。)

Basic Usage:
(基本用法：)

	type MyConfig struct {
		Server struct {
			Host string `mapstructure:"host" default:"localhost"`
			Port int    `mapstructure:"port" default:"8080"`
		} `mapstructure:"server"`
		DatabaseURL string `mapstructure:"database_url" default:"postgres://user:pass@host/db"`
	}

	var cfg MyConfig

	// Load config with defaults, file, env vars, and hot-reload
	// (加载配置，包含默认值、文件、环境变量和热重载)
	cm, err := config.LoadConfigAndWatch(
		&cfg,
		config.WithConfigFile("config/app.yaml"), // Specify config file (指定配置文件)
		config.WithEnvPrefix("APP"),             // Env var prefix (e.g., APP_SERVER_PORT) (环境变量前缀)
		config.WithHotReload(),                  // Enable hot-reload (启用热重载)
	)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Register a callback for changes (注册变更回调)
	cm.RegisterCallback(func(v *viper.Viper, currentCfg any) error {
		appCfg := currentCfg.(*MyConfig) // Type assertion (类型断言)
		log.Printf("Config reloaded! New Port: %d", appCfg.Server.Port)
		// Update application components based on new config... (根据新配置更新应用组件...)
		return nil
	})


	// Access configuration values (访问配置值)
	fmt.Println("Server Port:", cfg.Server.Port)
*/
package config