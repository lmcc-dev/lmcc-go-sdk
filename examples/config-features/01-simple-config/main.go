/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * Simple configuration loading example demonstrating basic config features.
 */

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
)

// SimpleAppConfig 简单应用配置结构体
// (SimpleAppConfig represents a simple application configuration structure)
type SimpleAppConfig struct {
	// 基本字段 (Basic fields)
	AppName    string `mapstructure:"app_name" default:"SimpleConfigExample"`
	Version    string `mapstructure:"version" default:"1.0.0"`
	Debug      bool   `mapstructure:"debug" default:"true"`
	Port       int    `mapstructure:"port" default:"8080"`
	
	// 时间相关字段 (Time-related fields)
	Timeout    time.Duration `mapstructure:"timeout" default:"30s"`
	Interval   time.Duration `mapstructure:"interval" default:"5m"`
	
	// 切片字段 (Slice fields)
	AllowedIPs []string `mapstructure:"allowed_ips" default:"127.0.0.1,::1"`
	Features   []string `mapstructure:"features" default:"auth,logging,metrics"`
	
	// 嵌套结构体 (Nested struct)
	Database *DatabaseConfig `mapstructure:"database"`
	Cache    *CacheConfig    `mapstructure:"cache"`
}

// DatabaseConfig 数据库配置
// (DatabaseConfig represents database configuration)
type DatabaseConfig struct {
	Driver          string        `mapstructure:"driver" default:"postgres"`
	Host            string        `mapstructure:"host" default:"localhost"`
	Port            int           `mapstructure:"port" default:"5432"`
	User            string        `mapstructure:"user" default:"app_user"`
	Database        string        `mapstructure:"database" default:"app_db"`
	Password        string        `mapstructure:"password"`  // No default for security
	MaxConnections  int           `mapstructure:"max_connections" default:"10"`
	ConnectTimeout  time.Duration `mapstructure:"connect_timeout" default:"10s"`
	EnableSSL       bool          `mapstructure:"enable_ssl" default:"false"`
}

// CacheConfig 缓存配置
// (CacheConfig represents cache configuration)
type CacheConfig struct {
	Enabled    bool          `mapstructure:"enabled" default:"true"`
	Type       string        `mapstructure:"type" default:"redis"`
	Host       string        `mapstructure:"host" default:"localhost"`
	Port       int           `mapstructure:"port" default:"6379"`
	TTL        time.Duration `mapstructure:"ttl" default:"1h"`
	MaxMemory  string        `mapstructure:"max_memory" default:"128mb"`
}

// validateConfig 验证配置的有效性
// (validateConfig validates the configuration validity)
func validateConfig(cfg *SimpleAppConfig) error {
	if cfg.Port < 1 || cfg.Port > 65535 {
		return errors.Errorf("invalid port number: %d (must be 1-65535)", cfg.Port)
	}
	
	if cfg.Timeout <= 0 {
		return errors.New("timeout must be positive")
	}
	
	if cfg.Database != nil {
		if cfg.Database.MaxConnections <= 0 {
			return errors.New("database max_connections must be positive")
		}
		
		if cfg.Database.ConnectTimeout <= 0 {
			return errors.New("database connect_timeout must be positive")
		}
	}
	
	if cfg.Cache != nil && cfg.Cache.Enabled {
		if cfg.Cache.Port < 1 || cfg.Cache.Port > 65535 {
			return errors.Errorf("invalid cache port: %d", cfg.Cache.Port)
		}
		
		if cfg.Cache.TTL <= 0 {
			return errors.New("cache TTL must be positive")
		}
	}
	
	return nil
}

// printConfig 打印配置信息
// (printConfig prints configuration information)
func printConfig(cfg *SimpleAppConfig) {
	fmt.Println("=== Configuration Summary ===")
	fmt.Printf("Application:\n")
	fmt.Printf("  Name: %s\n", cfg.AppName)
	fmt.Printf("  Version: %s\n", cfg.Version)
	fmt.Printf("  Debug: %t\n", cfg.Debug)
	fmt.Printf("  Port: %d\n", cfg.Port)
	fmt.Printf("  Timeout: %v\n", cfg.Timeout)
	fmt.Printf("  Interval: %v\n", cfg.Interval)
	fmt.Printf("  Allowed IPs: %v\n", cfg.AllowedIPs)
	fmt.Printf("  Features: %v\n", cfg.Features)
	fmt.Println()
	
	if cfg.Database != nil {
		fmt.Printf("Database:\n")
		fmt.Printf("  Driver: %s\n", cfg.Database.Driver)
		fmt.Printf("  Host: %s\n", cfg.Database.Host)
		fmt.Printf("  Port: %d\n", cfg.Database.Port)
		fmt.Printf("  User: %s\n", cfg.Database.User)
		fmt.Printf("  Database: %s\n", cfg.Database.Database)
		fmt.Printf("  Password: %s\n", maskPassword(cfg.Database.Password))
		fmt.Printf("  Max Connections: %d\n", cfg.Database.MaxConnections)
		fmt.Printf("  Connect Timeout: %v\n", cfg.Database.ConnectTimeout)
		fmt.Printf("  Enable SSL: %t\n", cfg.Database.EnableSSL)
		fmt.Println()
	}
	
	if cfg.Cache != nil {
		fmt.Printf("Cache:\n")
		fmt.Printf("  Enabled: %t\n", cfg.Cache.Enabled)
		fmt.Printf("  Type: %s\n", cfg.Cache.Type)
		fmt.Printf("  Host: %s\n", cfg.Cache.Host)
		fmt.Printf("  Port: %d\n", cfg.Cache.Port)
		fmt.Printf("  TTL: %v\n", cfg.Cache.TTL)
		fmt.Printf("  Max Memory: %s\n", cfg.Cache.MaxMemory)
		fmt.Println()
	}
}

// maskPassword 掩码密码显示
// (maskPassword masks password for display)
func maskPassword(password string) string {
	if password == "" {
		return "<not set>"
	}
	if len(password) <= 2 {
		return "***"
	}
	return password[:2] + "***"
}

// demonstrateDefaultValues 演示默认值的工作原理
// (demonstrateDefaultValues demonstrates how default values work)
func demonstrateDefaultValues() {
	fmt.Println("=== Demonstrating Default Values ===")
	fmt.Println("Loading configuration with defaults only (no config file)...")
	
	var cfg SimpleAppConfig
	
	// 尝试加载不存在的配置文件，这样只会使用默认值
	// (Try to load non-existent config file, so only defaults will be used)
	err := config.LoadConfig(&cfg, 
		config.WithConfigFile("non-existent.yaml", "yaml"))
	
	if err != nil {
		// 这是预期的错误，因为文件不存在
		// (This is expected error because file doesn't exist)
		fmt.Printf("Expected error (file not found): %v\n", err)
		
		// 但我们仍然可以展示默认值是如何工作的
		// (But we can still show how defaults work)
		fmt.Println("Demonstrating defaults without file loading...")
		
		// 手动初始化nil指针字段以触发默认值应用
		// (Manually initialize nil pointer fields to trigger default application)
		cfg.Database = &DatabaseConfig{}
		cfg.Cache = &CacheConfig{}
		
		// 这种情况下，我们展示什么是"零值"
		// (In this case, we show what are "zero values")
		fmt.Println("\nZero values (before applying defaults):")
		fmt.Printf("  AppName: '%s'\n", cfg.AppName)
		fmt.Printf("  Port: %d\n", cfg.Port)
		fmt.Printf("  Debug: %t\n", cfg.Debug)
		fmt.Printf("  Features: %v (length: %d)\n", cfg.Features, len(cfg.Features))
	}
	
	fmt.Println()
}

func main() {
	fmt.Println("=== Simple Configuration Loading Example ===")
	fmt.Println("This example demonstrates basic configuration loading with defaults.")
	fmt.Println()
	
	// 1. 演示默认值 (Demonstrate default values)
	demonstrateDefaultValues()
	
	// 2. 加载实际配置文件 (Load actual configuration file)
	fmt.Println("=== Loading Configuration from File ===")
	var cfg SimpleAppConfig
	
	err := config.LoadConfig(&cfg, 
		config.WithConfigFile("config.yaml", "yaml"))
	
	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		
		// 展示错误信息 (Show error information)
		if coder := errors.GetCoder(err); coder != nil {
			fmt.Printf("Error Code: %d\n", coder.Code())
			fmt.Printf("Error Type: %s\n", coder.String())
		}
		
		fmt.Println("\nTip: Make sure config.yaml exists in the current directory.")
		fmt.Println("You can create one based on the example in this directory.")
		os.Exit(1)
	}
	
	fmt.Println("✓ Configuration loaded successfully!")
	fmt.Println()
	
	// 3. 验证配置 (Validate configuration)
	fmt.Println("=== Validating Configuration ===")
	if err := validateConfig(&cfg); err != nil {
		fmt.Printf("Configuration validation failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Configuration validation passed!")
	fmt.Println()
	
	// 4. 显示配置信息 (Display configuration information)
	printConfig(&cfg)
	
	// 5. 演示配置的使用 (Demonstrate configuration usage)
	fmt.Println("=== Using Configuration ===")
	
	fmt.Printf("Starting %s version %s...\n", cfg.AppName, cfg.Version)
	
	if cfg.Debug {
		fmt.Println("Debug mode is enabled")
	}
	
	fmt.Printf("Server will listen on port %d\n", cfg.Port)
	fmt.Printf("Request timeout is set to %v\n", cfg.Timeout)
	
	if cfg.Database != nil {
		fmt.Printf("Database: %s://%s:%d/%s\n", 
			cfg.Database.Driver, cfg.Database.Host, cfg.Database.Port, cfg.Database.Database)
		fmt.Printf("Database max connections: %d\n", cfg.Database.MaxConnections)
	}
	
	if cfg.Cache != nil && cfg.Cache.Enabled {
		fmt.Printf("Cache: %s://%s:%d (TTL: %v)\n", 
			cfg.Cache.Type, cfg.Cache.Host, cfg.Cache.Port, cfg.Cache.TTL)
	} else {
		fmt.Println("Cache is disabled")
	}
	
	fmt.Printf("Enabled features: %v\n", cfg.Features)
	fmt.Printf("Allowed IPs: %v\n", cfg.AllowedIPs)
	
	fmt.Println()
	fmt.Println("=== Example completed successfully ===")
} 