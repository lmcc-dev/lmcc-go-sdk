/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Example application demonstrating the usage of lmcc-go-sdk/pkg/config.
 */

package main

import (
	"flag" // Import the flag package
	"fmt"
	"log"
	"os"

	sdkconfig "github.com/lmcc-dev/lmcc-go-sdk/pkg/config" // Import the SDK config package
)

// MyAppConfig 嵌入了 SDK 的 Config 并添加了自定义字段
// (MyAppConfig embeds the SDK's Config and adds custom fields)
type MyAppConfig struct {
	sdkconfig.Config                      // 嵌入 SDK 基础配置 (Embed SDK base config)
	CustomFeature    *CustomFeatureConfig `mapstructure:"customFeature"`
}

// CustomFeatureConfig 示例应用的自定义配置部分
// (Custom configuration section for the example application)
type CustomFeatureConfig struct {
	APIKey    string `mapstructure:"apiKey"`
	RateLimit int    `mapstructure:"rateLimit"`
	Enabled   bool   `mapstructure:"enabled"`
}

// AppCfg 是此应用的全局配置实例
// (AppCfg is the global configuration instance for this application)
var AppCfg MyAppConfig

func main() {
	// 使用 flag 包处理命令行参数来指定配置文件路径
	// (Use the flag package to handle command-line arguments for specifying the config file path)
	configFile := flag.String("config", "config.yaml", "Path to the configuration file (e.g., -config /path/to/config.yaml)")
	flag.Parse() // 解析命令行参数 (Parse command-line arguments)

	// 确定配置文件路径 (Determine config file path)
	// 通常，我们会从命令行标志或环境变量获取路径，这里为了简单直接指定
	// (Usually, we'd get the path from command-line flags or env vars, here we specify directly for simplicity)
	configFilePath := *configFile // 使用 flag 获取的路径 (Use the path obtained from the flag)

	fmt.Printf("Attempting to load configuration from: %s (use -config flag to override)\n", configFilePath)

	// 加载配置，注意传递 AppCfg 的指针
	// (Load configuration, note passing the pointer of AppCfg)
	err := sdkconfig.LoadConfig(&AppCfg, sdkconfig.WithConfigFile(configFilePath, "yaml"))
	if err != nil {
		log.Fatalf("FATAL: Failed to load configuration: %v\n", err)
		os.Exit(1) // 确保在致命错误后退出 (Ensure exit after fatal error)
	}

	fmt.Println("\n--- Configuration Loaded Successfully ---")

	// --- 访问 SDK 基础配置 ---
	// (Accessing SDK Base Configuration)
	fmt.Println("\n[SDK Base Config]")
	if AppCfg.Server != nil {
		fmt.Printf("Server Port (from file/env): %d\n", AppCfg.Server.Port) // Should be 9091 from config.yaml
		fmt.Printf("Server Mode (from file/env): %s\n", AppCfg.Server.Mode) // Should be debug
		fmt.Printf("Server Host (default): %s\n", AppCfg.Server.Host)       // Should be default 0.0.0.0
	} else {
		fmt.Println("Server configuration section not loaded.")
	}

	if AppCfg.Log != nil {
		fmt.Printf("Log Level (from file/env): %s\n", AppCfg.Log.Level) // Should be trace
		fmt.Printf("Log Format (default): %s\n", AppCfg.Log.Format)     // Should be default text
		fmt.Printf("Log Output (default): %s\n", AppCfg.Log.Output)     // Should be default stdout
	} else {
		fmt.Println("Log configuration section not loaded.")
	}

	if AppCfg.Database != nil {
		fmt.Printf("Database Host (from file/env): %s\n", AppCfg.Database.Host)   // Should be example-db-host
		fmt.Printf("Database User (from file/env): %s\n", AppCfg.Database.User)   // Should be appuser
		fmt.Printf("Database Name (from file/env): %s\n", AppCfg.Database.DBName) // Should be app_db
		fmt.Printf("Database Type (default): %s\n", AppCfg.Database.Type)         // Should be default mysql
	} else {
		fmt.Println("Database configuration section not loaded.")
	}

	// --- 访问自定义配置 ---
	// (Accessing Custom Configuration)
	fmt.Println("\n[Custom App Config]")
	if AppCfg.CustomFeature != nil {
		fmt.Printf("Custom Feature API Key: %s\n", AppCfg.CustomFeature.APIKey)
		fmt.Printf("Custom Feature Rate Limit: %d\n", AppCfg.CustomFeature.RateLimit)
		fmt.Printf("Custom Feature Enabled: %t\n", AppCfg.CustomFeature.Enabled)
	} else {
		fmt.Println("CustomFeature configuration section not loaded or defined in config file.")
	}

	fmt.Println("\n--- Example Finished ---")
}
