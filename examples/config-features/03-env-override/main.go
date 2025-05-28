/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * Environment variable override example demonstrating configuration precedence.
 */

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

// EnvOverrideConfig 环境变量覆盖配置结构体
// (EnvOverrideConfig represents configuration with environment variable override support)
type EnvOverrideConfig struct {
	config.Config                    // 嵌入SDK基础配置 (Embed SDK base configuration)
	
	// 应用配置 (Application configuration)
	App *AppConfig `mapstructure:"app"`
	
	// 数据库配置 (Database configuration)
	Database *DatabaseConfig `mapstructure:"database"`
	
	// 缓存配置 (Cache configuration)
	Cache *CacheConfig `mapstructure:"cache"`
	
	// API配置 (API configuration)
	API *APIConfig `mapstructure:"api"`
}

// AppConfig 应用配置
// (AppConfig represents application configuration)
type AppConfig struct {
	Name        string        `mapstructure:"name" default:"EnvOverrideExample"`
	Version     string        `mapstructure:"version" default:"1.0.0"`
	Environment string        `mapstructure:"environment" default:"development"`
	Debug       bool          `mapstructure:"debug" default:"true"`
	Port        int          `mapstructure:"port" default:"8080"`
	Timeout     time.Duration `mapstructure:"timeout" default:"30s"`
	MaxWorkers  int          `mapstructure:"max_workers" default:"10"`
}

// DatabaseConfig 数据库配置（常用于环境变量覆盖）
// (DatabaseConfig represents database configuration, commonly overridden by env vars)
type DatabaseConfig struct {
	Host         string `mapstructure:"host" default:"localhost"`
	Port         int    `mapstructure:"port" default:"5432"`
	User         string `mapstructure:"user" default:"app_user"`
	Password     string `mapstructure:"password"`  // 通常从环境变量获取 (Usually from env vars)
	Database     string `mapstructure:"database" default:"app_db"`
	SSLMode      string `mapstructure:"ssl_mode" default:"disable"`
	MaxPoolSize  int    `mapstructure:"max_pool_size" default:"20"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout" default:"5m"`
}

// CacheConfig 缓存配置
// (CacheConfig represents cache configuration)
type CacheConfig struct {
	Type        string        `mapstructure:"type" default:"redis"`
	Host        string        `mapstructure:"host" default:"localhost"`
	Port        int          `mapstructure:"port" default:"6379"`
	Password    string        `mapstructure:"password"`  // 敏感信息，通常从环境变量获取
	Database    int          `mapstructure:"database" default:"0"`
	MaxRetries  int          `mapstructure:"max_retries" default:"3"`
	PoolSize    int          `mapstructure:"pool_size" default:"10"`
	DefaultTTL  time.Duration `mapstructure:"default_ttl" default:"1h"`
}

// APIConfig API配置
// (APIConfig represents API configuration)
type APIConfig struct {
	BaseURL        string        `mapstructure:"base_url" default:"https://api.example.com"`
	APIKey         string        `mapstructure:"api_key"`        // 敏感信息 (Sensitive data)
	Secret         string        `mapstructure:"secret"`         // 敏感信息 (Sensitive data)
	Timeout        time.Duration `mapstructure:"timeout" default:"10s"`
	MaxRetries     int          `mapstructure:"max_retries" default:"3"`
	RateLimitRPS   int          `mapstructure:"rate_limit_rps" default:"100"`
}

// ConfigSource 配置来源枚举
// (ConfigSource represents the source of configuration values)
type ConfigSource int

const (
	SourceDefault ConfigSource = iota
	SourceFile
	SourceEnv
)

func (s ConfigSource) String() string {
	switch s {
	case SourceDefault:
		return "Default"
	case SourceFile:
		return "File"
	case SourceEnv:
		return "Environment"
	default:
		return "Unknown"
	}
}

// ConfigAnalyzer 配置分析器，用于展示配置来源
// (ConfigAnalyzer analyzes configuration sources)
type ConfigAnalyzer struct {
	envPrefix string
}

// NewConfigAnalyzer 创建配置分析器
// (NewConfigAnalyzer creates a configuration analyzer)
func NewConfigAnalyzer(envPrefix string) *ConfigAnalyzer {
	return &ConfigAnalyzer{envPrefix: envPrefix}
}

// AnalyzeConfigSources 分析配置值的来源
// (AnalyzeConfigSources analyzes the sources of configuration values)
func (ca *ConfigAnalyzer) AnalyzeConfigSources(cfg *EnvOverrideConfig) map[string]ConfigSource {
	sources := make(map[string]ConfigSource)
	
	// 分析应用配置 (Analyze app configuration)
	if cfg.App != nil {
		sources["app.name"] = ca.getSource("APP_NAME", cfg.App.Name, "EnvOverrideExample")
		sources["app.environment"] = ca.getSource("APP_ENVIRONMENT", cfg.App.Environment, "development")
		sources["app.debug"] = ca.getSource("APP_DEBUG", strconv.FormatBool(cfg.App.Debug), "true")
		sources["app.port"] = ca.getSource("APP_PORT", strconv.Itoa(cfg.App.Port), "8080")
		sources["app.max_workers"] = ca.getSource("APP_MAX_WORKERS", strconv.Itoa(cfg.App.MaxWorkers), "10")
	}
	
	// 分析数据库配置 (Analyze database configuration)
	if cfg.Database != nil {
		sources["database.host"] = ca.getSource("DATABASE_HOST", cfg.Database.Host, "localhost")
		sources["database.port"] = ca.getSource("DATABASE_PORT", strconv.Itoa(cfg.Database.Port), "5432")
		sources["database.user"] = ca.getSource("DATABASE_USER", cfg.Database.User, "app_user")
		sources["database.password"] = ca.getSource("DATABASE_PASSWORD", cfg.Database.Password, "")
		sources["database.ssl_mode"] = ca.getSource("DATABASE_SSL_MODE", cfg.Database.SSLMode, "disable")
	}
	
	// 分析缓存配置 (Analyze cache configuration)
	if cfg.Cache != nil {
		sources["cache.host"] = ca.getSource("CACHE_HOST", cfg.Cache.Host, "localhost")
		sources["cache.port"] = ca.getSource("CACHE_PORT", strconv.Itoa(cfg.Cache.Port), "6379")
		sources["cache.password"] = ca.getSource("CACHE_PASSWORD", cfg.Cache.Password, "")
		sources["cache.pool_size"] = ca.getSource("CACHE_POOL_SIZE", strconv.Itoa(cfg.Cache.PoolSize), "10")
	}
	
	// 分析API配置 (Analyze API configuration)
	if cfg.API != nil {
		sources["api.base_url"] = ca.getSource("API_BASE_URL", cfg.API.BaseURL, "https://api.example.com")
		sources["api.api_key"] = ca.getSource("API_API_KEY", cfg.API.APIKey, "")
		sources["api.secret"] = ca.getSource("API_SECRET", cfg.API.Secret, "")
		sources["api.rate_limit_rps"] = ca.getSource("API_RATE_LIMIT_RPS", strconv.Itoa(cfg.API.RateLimitRPS), "100")
	}
	
	return sources
}

// getSource 确定配置值的来源
// (getSource determines the source of a configuration value)
func (ca *ConfigAnalyzer) getSource(envKey, currentValue, defaultValue string) ConfigSource {
	fullEnvKey := ca.envPrefix + "_" + envKey
	envValue := os.Getenv(fullEnvKey)
	
	if envValue != "" && envValue == currentValue {
		return SourceEnv
	}
	
	if currentValue != defaultValue {
		return SourceFile
	}
	
	return SourceDefault
}

// printEnvironmentGuide 打印环境变量使用指南
// (printEnvironmentGuide prints environment variable usage guide)
func printEnvironmentGuide(prefix string) {
	fmt.Println("=== Environment Variable Override Guide ===")
	fmt.Printf("Prefix: %s_\n", prefix)
	fmt.Println()
	
	fmt.Println("Format: {PREFIX}_{SECTION}_{FIELD}")
	fmt.Println("Examples:")
	fmt.Printf("  %s_APP_NAME=\"My Custom App\"\n", prefix)
	fmt.Printf("  %s_APP_DEBUG=false\n", prefix)
	fmt.Printf("  %s_APP_PORT=9090\n", prefix)
	fmt.Printf("  %s_DATABASE_HOST=prod-db.example.com\n", prefix)
	fmt.Printf("  %s_DATABASE_PASSWORD=super_secret_password\n", prefix)
	fmt.Printf("  %s_CACHE_HOST=cache-cluster.example.com\n", prefix)
	fmt.Printf("  %s_API_API_KEY=abc123def456\n", prefix)
	fmt.Println()
	
	fmt.Println("Common use cases:")
	fmt.Println("1. Production deployment with different database")
	fmt.Println("2. Development with local cache server")
	fmt.Println("3. CI/CD pipeline with test environment settings")
	fmt.Println("4. Secure handling of passwords and API keys")
	fmt.Println()
}

// demonstrateEnvironmentOverride 演示环境变量覆盖
// (demonstrateEnvironmentOverride demonstrates environment variable override)
func demonstrateEnvironmentOverride() {
	fmt.Println("=== Demonstrating Environment Variable Override ===")
	fmt.Println()
	
	// 设置一些示例环境变量 (Set some example environment variables)
	envVars := map[string]string{
		"ENVDEMO_APP_NAME": "Environment Override Demo",
		"ENVDEMO_APP_DEBUG": "false",
		"ENVDEMO_APP_PORT": "9090",
		"ENVDEMO_DATABASE_HOST": "prod-database.example.com",
		"ENVDEMO_DATABASE_PASSWORD": "secure_production_password",
		"ENVDEMO_CACHE_HOST": "redis-cluster.example.com",
		"ENVDEMO_CACHE_PASSWORD": "redis_secret_key",
		"ENVDEMO_API_API_KEY": "demo_api_key_12345",
		"ENVDEMO_API_SECRET": "demo_secret_abcdef",
	}
	
	fmt.Println("Setting example environment variables:")
	for key, value := range envVars {
		os.Setenv(key, value)
		if strings.Contains(strings.ToLower(key), "password") || 
		   strings.Contains(strings.ToLower(key), "secret") || 
		   strings.Contains(strings.ToLower(key), "key") {
			fmt.Printf("  %s=%s\n", key, maskSensitiveValue(value))
		} else {
			fmt.Printf("  %s=%s\n", key, value)
		}
	}
	fmt.Println()
	
	// 清理函数 (Cleanup function)
	defer func() {
		for key := range envVars {
			os.Unsetenv(key)
		}
	}()
	
	// 加载配置 (Load configuration)
	var cfg EnvOverrideConfig
	err := config.LoadConfig(&cfg,
		config.WithConfigFile("config.yaml", "yaml"),
		config.WithEnvPrefix("ENVDEMO"),
		config.WithEnvVarOverride(true),
	)
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		return
	}
	
	// 分析配置来源 (Analyze configuration sources)
	analyzer := NewConfigAnalyzer("ENVDEMO")
	sources := analyzer.AnalyzeConfigSources(&cfg)
	
	fmt.Println("=== Configuration Analysis ===")
	printConfigWithSources(&cfg, sources)
}

// printConfigWithSources 打印配置及其来源
// (printConfigWithSources prints configuration with their sources)
func printConfigWithSources(cfg *EnvOverrideConfig, sources map[string]ConfigSource) {
	fmt.Println("Configuration values and their sources:")
	fmt.Println()
	
	fmt.Println("Application:")
	printValueWithSource("  Name", cfg.App.Name, sources["app.name"])
	printValueWithSource("  Environment", cfg.App.Environment, sources["app.environment"])
	printValueWithSource("  Debug", strconv.FormatBool(cfg.App.Debug), sources["app.debug"])
	printValueWithSource("  Port", strconv.Itoa(cfg.App.Port), sources["app.port"])
	printValueWithSource("  Max Workers", strconv.Itoa(cfg.App.MaxWorkers), sources["app.max_workers"])
	fmt.Println()
	
	fmt.Println("Database:")
	printValueWithSource("  Host", cfg.Database.Host, sources["database.host"])
	printValueWithSource("  Port", strconv.Itoa(cfg.Database.Port), sources["database.port"])
	printValueWithSource("  User", cfg.Database.User, sources["database.user"])
	printValueWithSource("  Password", maskSensitiveValue(cfg.Database.Password), sources["database.password"])
	printValueWithSource("  SSL Mode", cfg.Database.SSLMode, sources["database.ssl_mode"])
	fmt.Println()
	
	fmt.Println("Cache:")
	printValueWithSource("  Host", cfg.Cache.Host, sources["cache.host"])
	printValueWithSource("  Port", strconv.Itoa(cfg.Cache.Port), sources["cache.port"])
	printValueWithSource("  Password", maskSensitiveValue(cfg.Cache.Password), sources["cache.password"])
	printValueWithSource("  Pool Size", strconv.Itoa(cfg.Cache.PoolSize), sources["cache.pool_size"])
	fmt.Println()
	
	fmt.Println("API:")
	printValueWithSource("  Base URL", cfg.API.BaseURL, sources["api.base_url"])
	printValueWithSource("  API Key", maskSensitiveValue(cfg.API.APIKey), sources["api.api_key"])
	printValueWithSource("  Secret", maskSensitiveValue(cfg.API.Secret), sources["api.secret"])
	printValueWithSource("  Rate Limit RPS", strconv.Itoa(cfg.API.RateLimitRPS), sources["api.rate_limit_rps"])
}

// printValueWithSource 打印值及其来源
// (printValueWithSource prints a value with its source)
func printValueWithSource(label, value string, source ConfigSource) {
	sourceColor := getSourceColor(source)
	fmt.Printf("%s: %s %s[%s]%s\n", label, value, sourceColor, source.String(), "\033[0m")
}

// getSourceColor 获取来源的颜色代码
// (getSourceColor gets color code for source)
func getSourceColor(source ConfigSource) string {
	switch source {
	case SourceDefault:
		return "\033[90m"  // 灰色 (Gray)
	case SourceFile:
		return "\033[34m"  // 蓝色 (Blue)
	case SourceEnv:
		return "\033[32m"  // 绿色 (Green)
	default:
		return "\033[0m"   // 重置 (Reset)
	}
}

// maskSensitiveValue 掩码敏感值
// (maskSensitiveValue masks sensitive values)
func maskSensitiveValue(value string) string {
	if value == "" {
		return "<not set>"
	}
	if len(value) <= 4 {
		return "****"
	}
	return value[:2] + "****" + value[len(value)-2:]
}

// printSecurityBestPractices 打印安全最佳实践
// (printSecurityBestPractices prints security best practices)
func printSecurityBestPractices() {
	fmt.Println("=== Security Best Practices ===")
	fmt.Println()
	
	fmt.Println("1. Sensitive Data in Environment Variables:")
	fmt.Println("   ✓ Passwords, API keys, secrets")
	fmt.Println("   ✓ Database connection strings")
	fmt.Println("   ✓ Encryption keys and certificates")
	fmt.Println()
	
	fmt.Println("2. Non-Sensitive Data in Configuration Files:")
	fmt.Println("   ✓ Application names and versions")
	fmt.Println("   ✓ Timeout values and worker counts")
	fmt.Println("   ✓ Feature flags and debug settings")
	fmt.Println()
	
	fmt.Println("3. Deployment Recommendations:")
	fmt.Println("   • Use container orchestration for env var management")
	fmt.Println("   • Leverage secret management systems (K8s secrets, AWS Secrets Manager)")
	fmt.Println("   • Never commit .env files with real secrets to version control")
	fmt.Println("   • Rotate secrets regularly")
	fmt.Println("   • Use principle of least privilege for environment access")
	fmt.Println()
}

func main() {
	fmt.Println("=== Environment Variable Override Example ===")
	fmt.Println("This example demonstrates configuration precedence and environment variable override.")
	fmt.Println()
	
	// 1. 初始化日志 (Initialize logging)
	logOpts := log.NewOptions()
	logOpts.Level = "info"
	logOpts.Format = "text"
	logOpts.EnableColor = true
	log.Init(logOpts)
	logger := log.Std()
	
	// 2. 打印环境变量使用指南 (Print environment variable guide)
	printEnvironmentGuide("ENVDEMO")
	
	// 3. 加载基础配置（仅文件和默认值）(Load base configuration - file and defaults only)
	fmt.Println("=== Loading Base Configuration (File + Defaults) ===")
	var baseCfg EnvOverrideConfig
	err := config.LoadConfig(&baseCfg,
		config.WithConfigFile("config.yaml", "yaml"),
		config.WithEnvPrefix("ENVDEMO"),
		config.WithEnvVarOverride(false), // 禁用环境变量覆盖 (Disable env var override)
	)
	if err != nil {
		logger.Errorf("Failed to load base configuration: %v", err)
		if coder := errors.GetCoder(err); coder != nil {
			fmt.Printf("Error Code: %d, Type: %s\n", coder.Code(), coder.String())
		}
		// 继续演示，使用默认值 (Continue demo with defaults)
	} else {
		fmt.Println("✓ Base configuration loaded successfully")
	}
	
	// 4. 分析基础配置来源 (Analyze base configuration sources)
	baseAnalyzer := NewConfigAnalyzer("ENVDEMO")
	baseSources := baseAnalyzer.AnalyzeConfigSources(&baseCfg)
	
	fmt.Println()
	fmt.Println("Base configuration (before environment override):")
	printConfigWithSources(&baseCfg, baseSources)
	fmt.Println()
	
	// 5. 演示环境变量覆盖 (Demonstrate environment variable override)
	demonstrateEnvironmentOverride()
	
	// 6. 打印安全最佳实践 (Print security best practices)
	fmt.Println()
	printSecurityBestPractices()
	
	fmt.Println("=== Example completed successfully ===")
} 