/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * Multiple formats example demonstrating YAML, JSON, and TOML configuration support.
 */

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

// MultiFormatConfig 多格式配置结构体
// (MultiFormatConfig represents configuration supporting multiple file formats)
type MultiFormatConfig struct {
	config.Config                        // 嵌入SDK基础配置 (Embed SDK base configuration)
	
	// 应用配置 (Application configuration)
	App *AppConfig `mapstructure:"app"`
	
	// 服务器配置 (Server configuration)
	Server *ServerConfig `mapstructure:"server"`
	
	// 数据库配置 (Database configuration)
	Database *DatabaseConfig `mapstructure:"database"`
	
	// 功能开关 (Feature flags)
	Features *FeatureConfig `mapstructure:"features"`
}

// AppConfig 应用配置
// (AppConfig represents application configuration)
type AppConfig struct {
	Name        string        `mapstructure:"name" default:"MultiFormatExample"`
	Version     string        `mapstructure:"version" default:"1.0.0"`
	Environment string        `mapstructure:"environment" default:"development"`
	Debug       bool          `mapstructure:"debug" default:"true"`
	Timeout     time.Duration `mapstructure:"timeout" default:"30s"`
	Tags        []string      `mapstructure:"tags" default:"example,demo,multiformat"`
}

// ServerConfig 服务器配置
// (ServerConfig represents server configuration)
type ServerConfig struct {
	Host         string        `mapstructure:"host" default:"localhost"`
	Port         int          `mapstructure:"port" default:"8080"`
	TLS          bool         `mapstructure:"tls" default:"false"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout" default:"10s"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" default:"10s"`
}

// DatabaseConfig 数据库配置
// (DatabaseConfig represents database configuration)
type DatabaseConfig struct {
	Type         string `mapstructure:"type" default:"postgres"`
	Host         string `mapstructure:"host" default:"localhost"`
	Port         int    `mapstructure:"port" default:"5432"`
	User         string `mapstructure:"user" default:"app_user"`
	Password     string `mapstructure:"password" default:"app_password"`
	Database     string `mapstructure:"database" default:"app_db"`
	MaxConns     int    `mapstructure:"max_conns" default:"10"`
	SSLEnabled   bool   `mapstructure:"ssl_enabled" default:"false"`
}

// FeatureConfig 功能开关配置
// (FeatureConfig represents feature flags configuration)
type FeatureConfig struct {
	EnableAuth     bool `mapstructure:"enable_auth" default:"true"`
	EnableMetrics  bool `mapstructure:"enable_metrics" default:"true"`
	EnableLogging  bool `mapstructure:"enable_logging" default:"true"`
	EnableCaching  bool `mapstructure:"enable_caching" default:"false"`
	BetaFeatures   bool `mapstructure:"beta_features" default:"false"`
}

// FormatInfo 格式信息
// (FormatInfo represents configuration format information)
type FormatInfo struct {
	Name        string
	Extension   string
	Description string
	Supported   bool
}

// getSupportedFormats 获取支持的配置格式
// (getSupportedFormats returns supported configuration formats)
func getSupportedFormats() []FormatInfo {
	return []FormatInfo{
		{
			Name:        "YAML",
			Extension:   "yaml",
			Description: "YAML Ain't Markup Language - Human-readable data serialization",
			Supported:   true,
		},
		{
			Name:        "YAML (yml)",
			Extension:   "yml",
			Description: "YAML with .yml extension",
			Supported:   true,
		},
		{
			Name:        "JSON",
			Extension:   "json",
			Description: "JavaScript Object Notation - Lightweight data interchange",
			Supported:   true,
		},
		{
			Name:        "TOML",
			Extension:   "toml",
			Description: "Tom's Obvious, Minimal Language - Configuration file format",
			Supported:   true,
		},
	}
}

// printFormatInfo 打印格式信息
// (printFormatInfo prints format information)
func printFormatInfo() {
	fmt.Println("=== Supported Configuration Formats ===")
	formats := getSupportedFormats()
	
	for _, format := range formats {
		status := "✓"
		if !format.Supported {
			status = "✗"
		}
		fmt.Printf("%s %s (.%s): %s\n", status, format.Name, format.Extension, format.Description)
	}
	fmt.Println()
}

// loadConfigFromFormat 从指定格式加载配置
// (loadConfigFromFormat loads configuration from specified format)
func loadConfigFromFormat(filename, format string) (*MultiFormatConfig, error) {
	var cfg MultiFormatConfig
	
	err := config.LoadConfig(&cfg,
		config.WithConfigFile(filename, format),
		config.WithEnvPrefix("MULTIFORMAT"),
		config.WithEnvVarOverride(true),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load %s configuration from %s", format, filename)
	}
	
	return &cfg, nil
}

// printConfigSummary 打印配置摘要
// (printConfigSummary prints configuration summary)
func printConfigSummary(cfg *MultiFormatConfig, format string) {
	fmt.Printf("=== Configuration Summary (%s) ===\n", format)
	
	if cfg.App != nil {
		fmt.Printf("Application:\n")
		fmt.Printf("  Name: %s\n", cfg.App.Name)
		fmt.Printf("  Version: %s\n", cfg.App.Version)
		fmt.Printf("  Environment: %s\n", cfg.App.Environment)
		fmt.Printf("  Debug: %t\n", cfg.App.Debug)
		fmt.Printf("  Timeout: %v\n", cfg.App.Timeout)
		fmt.Printf("  Tags: %v\n", cfg.App.Tags)
	}
	
	if cfg.Server != nil {
		fmt.Printf("Server:\n")
		fmt.Printf("  Host: %s\n", cfg.Server.Host)
		fmt.Printf("  Port: %d\n", cfg.Server.Port)
		fmt.Printf("  TLS: %t\n", cfg.Server.TLS)
		fmt.Printf("  Read Timeout: %v\n", cfg.Server.ReadTimeout)
	}
	
	if cfg.Database != nil {
		fmt.Printf("Database:\n")
		fmt.Printf("  Type: %s\n", cfg.Database.Type)
		fmt.Printf("  Host: %s\n", cfg.Database.Host)
		fmt.Printf("  Port: %d\n", cfg.Database.Port)
		fmt.Printf("  Database: %s\n", cfg.Database.Database)
		fmt.Printf("  Max Conns: %d\n", cfg.Database.MaxConns)
		fmt.Printf("  SSL: %t\n", cfg.Database.SSLEnabled)
	}
	
	if cfg.Features != nil {
		fmt.Printf("Features:\n")
		fmt.Printf("  Auth: %t\n", cfg.Features.EnableAuth)
		fmt.Printf("  Metrics: %t\n", cfg.Features.EnableMetrics)
		fmt.Printf("  Logging: %t\n", cfg.Features.EnableLogging)
		fmt.Printf("  Caching: %t\n", cfg.Features.EnableCaching)
		fmt.Printf("  Beta: %t\n", cfg.Features.BetaFeatures)
	}
	
	fmt.Println()
}

// demonstrateFormatSupport 演示不同格式的支持
// (demonstrateFormatSupport demonstrates support for different formats)
func demonstrateFormatSupport() {
	fmt.Println("=== Demonstrating Multiple Format Support ===")
	fmt.Println()
	
	// 支持的格式列表 (List of supported formats)
	formats := []struct {
		filename string
		format   string
		desc     string
	}{
		{"config.yaml", "yaml", "YAML format"},
		{"config.yml", "yml", "YAML with .yml extension"},
		{"config.json", "json", "JSON format"},
		{"config.toml", "toml", "TOML format"},
	}
	
	for _, f := range formats {
		fmt.Printf("Testing %s (%s)...\n", f.desc, f.filename)
		
		// 检查文件是否存在 (Check if file exists)
		if _, err := os.Stat(f.filename); os.IsNotExist(err) {
			fmt.Printf("  ⚠ File %s not found, skipping\n", f.filename)
			fmt.Println()
			continue
		}
		
		// 尝试加载配置 (Try to load configuration)
		cfg, err := loadConfigFromFormat(f.filename, f.format)
		if err != nil {
			fmt.Printf("  ✗ Failed to load %s: %v\n", f.filename, err)
			if coder := errors.GetCoder(err); coder != nil {
				fmt.Printf("    Error Code: %d, Type: %s\n", coder.Code(), coder.String())
			}
		} else {
			fmt.Printf("  ✓ Successfully loaded %s\n", f.filename)
			printConfigSummary(cfg, f.format)
		}
	}
}

// demonstrateAutoDetection 演示格式自动检测
// (demonstrateAutoDetection demonstrates automatic format detection)
func demonstrateAutoDetection() {
	fmt.Println("=== Demonstrating Automatic Format Detection ===")
	fmt.Println()
	
	// 测试文件列表 (Test file list)
	testFiles := []string{
		"config.yaml",
		"config.yml", 
		"config.json",
		"config.toml",
	}
	
	for _, filename := range testFiles {
		fmt.Printf("Auto-detecting format for %s...\n", filename)
		
		// 检查文件是否存在 (Check if file exists)
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			fmt.Printf("  ⚠ File %s not found, skipping\n", filename)
			fmt.Println()
			continue
		}
		
		// 从文件扩展名推断格式 (Infer format from file extension)
		ext := filepath.Ext(filename)
		if ext != "" {
			ext = ext[1:] // 移除点号 (Remove dot)
		}
		
		// 特殊处理 yml 扩展名 (Special handling for yml extension)
		format := ext
		if ext == "yml" {
			format = "yaml"
		}
		
		fmt.Printf("  Detected format: %s\n", format)
		
		// 加载配置 (Load configuration)
		cfg, err := loadConfigFromFormat(filename, format)
		if err != nil {
			fmt.Printf("  ✗ Failed to load with auto-detected format: %v\n", err)
		} else {
			fmt.Printf("  ✓ Successfully loaded with auto-detected format\n")
			// 只打印应用名称作为验证 (Only print app name as verification)
			if cfg.App != nil {
				fmt.Printf("  App Name: %s\n", cfg.App.Name)
			}
		}
		fmt.Println()
	}
}

// printFormatComparisonGuide 打印格式对比指南
// (printFormatComparisonGuide prints format comparison guide)
func printFormatComparisonGuide() {
	fmt.Println("=== Configuration Format Comparison Guide ===")
	fmt.Println()
	
	fmt.Println("YAML:")
	fmt.Println("  ✓ Human-readable and writable")
	fmt.Println("  ✓ Supports comments")
	fmt.Println("  ✓ Supports multi-line strings")
	fmt.Println("  ✓ Widely used in DevOps")
	fmt.Println("  ⚠ Indentation-sensitive")
	fmt.Println("  ⚠ Can be verbose for simple configs")
	fmt.Println()
	
	fmt.Println("JSON:")
	fmt.Println("  ✓ Widely supported across languages")
	fmt.Println("  ✓ Compact and fast to parse")
	fmt.Println("  ✓ Well-defined standard")
	fmt.Println("  ✗ No comments support")
	fmt.Println("  ✗ No multi-line strings")
	fmt.Println("  ⚠ Requires quotes for all strings")
	fmt.Println()
	
	fmt.Println("TOML:")
	fmt.Println("  ✓ Human-readable and writable")
	fmt.Println("  ✓ Supports comments")
	fmt.Println("  ✓ Clear section structure")
	fmt.Println("  ✓ Type-safe (dates, numbers, booleans)")
	fmt.Println("  ⚠ Less common than YAML/JSON")
	fmt.Println("  ⚠ Can be verbose for nested structures")
	fmt.Println()
	
	fmt.Println("Recommendations:")
	fmt.Println("• Use YAML for:")
	fmt.Println("  - Complex configurations with deep nesting")
	fmt.Println("  - When comments and documentation are important")
	fmt.Println("  - DevOps and CI/CD configurations")
	fmt.Println()
	fmt.Println("• Use JSON for:")
	fmt.Println("  - API configurations and data exchange")
	fmt.Println("  - When parsing speed is critical")
	fmt.Println("  - Integration with web services")
	fmt.Println()
	fmt.Println("• Use TOML for:")
	fmt.Println("  - Application configurations")
	fmt.Println("  - When type safety is important")
	fmt.Println("  - Rust/Go ecosystem projects")
	fmt.Println()
}

// printMigrationTips 打印格式迁移技巧
// (printMigrationTips prints format migration tips)
func printMigrationTips() {
	fmt.Println("=== Format Migration Tips ===")
	fmt.Println()
	
	fmt.Println("Converting between formats:")
	fmt.Println("1. YAML ↔ JSON:")
	fmt.Println("   • Use online converters or tools like yq")
	fmt.Println("   • Watch out for: comments (lost in JSON), data types")
	fmt.Println()
	
	fmt.Println("2. YAML ↔ TOML:")
	fmt.Println("   • Manual conversion usually required")
	fmt.Println("   • TOML sections map to YAML nested objects")
	fmt.Println("   • Arrays and nested structures need attention")
	fmt.Println()
	
	fmt.Println("3. JSON ↔ TOML:")
	fmt.Println("   • Use specialized tools or libraries")
	fmt.Println("   • Consider data type preservation")
	fmt.Println("   • Test thoroughly after conversion")
	fmt.Println()
	
	fmt.Println("Best practices:")
	fmt.Println("• Keep format-specific examples in documentation")
	fmt.Println("• Use validation to ensure config correctness")
	fmt.Println("• Provide format-specific default config files")
	fmt.Println("• Test your application with all supported formats")
	fmt.Println()
}

func main() {
	fmt.Println("=== Multiple Configuration Formats Example ===")
	fmt.Println("This example demonstrates support for YAML, JSON, and TOML configuration files.")
	fmt.Println()
	
	// 1. 初始化日志 (Initialize logging)
	logOpts := log.NewOptions()
	logOpts.Level = "info"
	logOpts.Format = "text"
	logOpts.EnableColor = true
	log.Init(logOpts)
	logger := log.Std()
	
	// 2. 打印支持的格式信息 (Print supported format information)
	printFormatInfo()
	
	// 3. 演示多格式支持 (Demonstrate multiple format support)
	demonstrateFormatSupport()
	
	// 4. 演示自动格式检测 (Demonstrate automatic format detection)
	demonstrateAutoDetection()
	
	// 5. 打印格式对比指南 (Print format comparison guide)
	printFormatComparisonGuide()
	
	// 6. 打印迁移技巧 (Print migration tips)
	printMigrationTips()
	
	logger.Info("Multiple formats example completed successfully")
	fmt.Println("=== Example completed successfully ===")
} 