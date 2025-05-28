/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * Default values example demonstrating comprehensive default handling patterns.
 */

package main

import (
	"fmt"
	"reflect"
	"time"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

// DefaultValuesConfig 默认值演示配置结构体
// (DefaultValuesConfig demonstrates various default value patterns)
type DefaultValuesConfig struct {
	config.Config                        // 嵌入SDK基础配置 (Embed SDK base configuration)
	
	// 基础类型默认值 (Basic type defaults)
	BasicTypes *BasicTypesConfig `mapstructure:"basic_types"`
	
	// 集合类型默认值 (Collection type defaults)
	Collections *CollectionsConfig `mapstructure:"collections"`
	
	// 嵌套结构默认值 (Nested struct defaults)
	Nested *NestedConfig `mapstructure:"nested"`
	
	// 条件默认值 (Conditional defaults)
	Conditional *ConditionalConfig `mapstructure:"conditional"`
	
	// 复杂默认值 (Complex defaults)
	Complex *ComplexConfig `mapstructure:"complex"`
}

// BasicTypesConfig 基础类型默认值配置
// (BasicTypesConfig demonstrates basic type default values)
type BasicTypesConfig struct {
	// 字符串类型 (String types)
	StringValue     string `mapstructure:"string_value" default:"default_string"`
	EmptyString     string `mapstructure:"empty_string" default:""`
	MultiWordString string `mapstructure:"multi_word_string" default:"Hello World Default"`
	
	// 数值类型 (Numeric types)
	IntValue     int     `mapstructure:"int_value" default:"42"`
	Int32Value   int32   `mapstructure:"int32_value" default:"32"`
	Int64Value   int64   `mapstructure:"int64_value" default:"64"`
	FloatValue   float32 `mapstructure:"float_value" default:"3.14"`
	DoubleValue  float64 `mapstructure:"double_value" default:"2.718"`
	
	// 布尔类型 (Boolean types)
	BoolTrue     bool `mapstructure:"bool_true" default:"true"`
	BoolFalse    bool `mapstructure:"bool_false" default:"false"`
	
	// 时间类型 (Time types)
	Duration     time.Duration `mapstructure:"duration" default:"30s"`
	LongDuration time.Duration `mapstructure:"long_duration" default:"1h30m"`
}

// CollectionsConfig 集合类型默认值配置
// (CollectionsConfig demonstrates collection type default values)
type CollectionsConfig struct {
	// 字符串切片 (String slices)
	StringSlice    []string `mapstructure:"string_slice" default:"item1,item2,item3"`
	EmptySlice     []string `mapstructure:"empty_slice" default:""`
	SingleItem     []string `mapstructure:"single_item" default:"single"`
	
	// 数值切片 (Numeric slices)
	IntSlice       []int    `mapstructure:"int_slice" default:"1,2,3,4,5"`
	FloatSlice     []float64 `mapstructure:"float_slice" default:"1.1,2.2,3.3"`
	
	// 映射类型 (Map types) - 注意：映射的默认值处理比较复杂
	StringMap      map[string]string `mapstructure:"string_map"`
	IntMap         map[string]int    `mapstructure:"int_map"`
}

// NestedConfig 嵌套结构默认值配置
// (NestedConfig demonstrates nested struct default values)
type NestedConfig struct {
	Level1 *Level1Config `mapstructure:"level1"`
}

// Level1Config 第一层嵌套配置
// (Level1Config represents first level nested configuration)
type Level1Config struct {
	Name    string        `mapstructure:"name" default:"Level1Default"`
	Value   int          `mapstructure:"value" default:"100"`
	Enabled bool         `mapstructure:"enabled" default:"true"`
	Level2  *Level2Config `mapstructure:"level2"`
}

// Level2Config 第二层嵌套配置
// (Level2Config represents second level nested configuration)
type Level2Config struct {
	Name        string        `mapstructure:"name" default:"Level2Default"`
	Timeout     time.Duration `mapstructure:"timeout" default:"15s"`
	MaxRetries  int          `mapstructure:"max_retries" default:"5"`
	Features    []string      `mapstructure:"features" default:"feature1,feature2"`
}

// ConditionalConfig 条件默认值配置
// (ConditionalConfig demonstrates conditional default values)
type ConditionalConfig struct {
	// 基础字段 (Base fields)
	Mode        string `mapstructure:"mode" default:"development"`
	Environment string `mapstructure:"environment" default:"dev"`
	
	// 开发环境配置 (Development configuration)
	DevConfig   *DevEnvironmentConfig  `mapstructure:"dev_config"`
	
	// 生产环境配置 (Production configuration)
	ProdConfig  *ProdEnvironmentConfig `mapstructure:"prod_config"`
}

// DevEnvironmentConfig 开发环境配置
// (DevEnvironmentConfig represents development environment configuration)
type DevEnvironmentConfig struct {
	Debug           bool          `mapstructure:"debug" default:"true"`
	LogLevel        string        `mapstructure:"log_level" default:"debug"`
	HotReload       bool          `mapstructure:"hot_reload" default:"true"`
	MockExternalAPI bool          `mapstructure:"mock_external_api" default:"true"`
	TestTimeout     time.Duration `mapstructure:"test_timeout" default:"5s"`
}

// ProdEnvironmentConfig 生产环境配置
// (ProdEnvironmentConfig represents production environment configuration)
type ProdEnvironmentConfig struct {
	Debug           bool          `mapstructure:"debug" default:"false"`
	LogLevel        string        `mapstructure:"log_level" default:"info"`
	HotReload       bool          `mapstructure:"hot_reload" default:"false"`
	MockExternalAPI bool          `mapstructure:"mock_external_api" default:"false"`
	RequestTimeout  time.Duration `mapstructure:"request_timeout" default:"30s"`
	MaxConnections  int          `mapstructure:"max_connections" default:"1000"`
}

// ComplexConfig 复杂默认值配置
// (ComplexConfig demonstrates complex default value patterns)
type ComplexConfig struct {
	// 指针类型默认值 (Pointer type defaults)
	OptionalString *string `mapstructure:"optional_string"`
	OptionalInt    *int    `mapstructure:"optional_int"`
	OptionalBool   *bool   `mapstructure:"optional_bool"`
	
	// 接口类型（需要特殊处理）(Interface types - need special handling)
	CustomInterface interface{} `mapstructure:"custom_interface"`
	
	// 自定义类型 (Custom types)
	Priority       Priority       `mapstructure:"priority" default:"medium"`
	Protocol       Protocol       `mapstructure:"protocol" default:"https"`
}

// Priority 优先级枚举
// (Priority represents priority enumeration)
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

// Protocol 协议枚举
// (Protocol represents protocol enumeration)
type Protocol string

const (
	ProtocolHTTP  Protocol = "http"
	ProtocolHTTPS Protocol = "https"
	ProtocolTCP   Protocol = "tcp"
	ProtocolUDP   Protocol = "udp"
)

// DefaultValueAnalyzer 默认值分析器
// (DefaultValueAnalyzer analyzes default values)
type DefaultValueAnalyzer struct{}

// NewDefaultValueAnalyzer 创建默认值分析器
// (NewDefaultValueAnalyzer creates a default value analyzer)
func NewDefaultValueAnalyzer() *DefaultValueAnalyzer {
	return &DefaultValueAnalyzer{}
}

// AnalyzeDefaults 分析默认值的应用
// (AnalyzeDefaults analyzes the application of default values)
func (dva *DefaultValueAnalyzer) AnalyzeDefaults(cfg *DefaultValuesConfig) {
	fmt.Println("=== Default Values Analysis ===")
	
	// 分析基础类型 (Analyze basic types)
	if cfg.BasicTypes != nil {
		fmt.Println("\nBasic Types:")
		dva.analyzeBasicTypes(cfg.BasicTypes)
	}
	
	// 分析集合类型 (Analyze collections)
	if cfg.Collections != nil {
		fmt.Println("\nCollections:")
		dva.analyzeCollections(cfg.Collections)
	}
	
	// 分析嵌套结构 (Analyze nested structures)
	if cfg.Nested != nil {
		fmt.Println("\nNested Structures:")
		dva.analyzeNested(cfg.Nested)
	}
	
	// 分析条件默认值 (Analyze conditional defaults)
	if cfg.Conditional != nil {
		fmt.Println("\nConditional Defaults:")
		dva.analyzeConditional(cfg.Conditional)
	}
	
	// 分析复杂类型 (Analyze complex types)
	if cfg.Complex != nil {
		fmt.Println("\nComplex Types:")
		dva.analyzeComplex(cfg.Complex)
	}
}

// analyzeBasicTypes 分析基础类型默认值
// (analyzeBasicTypes analyzes basic type default values)
func (dva *DefaultValueAnalyzer) analyzeBasicTypes(bt *BasicTypesConfig) {
	fmt.Printf("  String Value: '%s' (expected: 'default_string')\n", bt.StringValue)
	fmt.Printf("  Empty String: '%s' (expected: '')\n", bt.EmptyString)
	fmt.Printf("  Multi Word: '%s' (expected: 'Hello World Default')\n", bt.MultiWordString)
	fmt.Printf("  Int Value: %d (expected: 42)\n", bt.IntValue)
	fmt.Printf("  Float Value: %.2f (expected: 3.14)\n", bt.FloatValue)
	fmt.Printf("  Bool True: %t (expected: true)\n", bt.BoolTrue)
	fmt.Printf("  Bool False: %t (expected: false)\n", bt.BoolFalse)
	fmt.Printf("  Duration: %v (expected: 30s)\n", bt.Duration)
	fmt.Printf("  Long Duration: %v (expected: 1h30m)\n", bt.LongDuration)
}

// analyzeCollections 分析集合类型默认值
// (analyzeCollections analyzes collection type default values)
func (dva *DefaultValueAnalyzer) analyzeCollections(c *CollectionsConfig) {
	fmt.Printf("  String Slice: %v (expected: [item1 item2 item3])\n", c.StringSlice)
	fmt.Printf("  Empty Slice: %v (expected: [])\n", c.EmptySlice)
	fmt.Printf("  Single Item: %v (expected: [single])\n", c.SingleItem)
	fmt.Printf("  Int Slice: %v (expected: [1 2 3 4 5])\n", c.IntSlice)
	fmt.Printf("  Float Slice: %v (expected: [1.1 2.2 3.3])\n", c.FloatSlice)
	fmt.Printf("  String Map: %v (maps need special handling)\n", c.StringMap)
	fmt.Printf("  Int Map: %v (maps need special handling)\n", c.IntMap)
}

// analyzeNested 分析嵌套结构默认值
// (analyzeNested analyzes nested structure default values)
func (dva *DefaultValueAnalyzer) analyzeNested(n *NestedConfig) {
	if n.Level1 != nil {
		fmt.Printf("  Level1 Name: '%s' (expected: 'Level1Default')\n", n.Level1.Name)
		fmt.Printf("  Level1 Value: %d (expected: 100)\n", n.Level1.Value)
		fmt.Printf("  Level1 Enabled: %t (expected: true)\n", n.Level1.Enabled)
		
		if n.Level1.Level2 != nil {
			fmt.Printf("  Level2 Name: '%s' (expected: 'Level2Default')\n", n.Level1.Level2.Name)
			fmt.Printf("  Level2 Timeout: %v (expected: 15s)\n", n.Level1.Level2.Timeout)
			fmt.Printf("  Level2 Max Retries: %d (expected: 5)\n", n.Level1.Level2.MaxRetries)
			fmt.Printf("  Level2 Features: %v (expected: [feature1 feature2])\n", n.Level1.Level2.Features)
		} else {
			fmt.Println("  Level2: nil (nested pointers need initialization)")
		}
	} else {
		fmt.Println("  Level1: nil (nested pointers need initialization)")
	}
}

// analyzeConditional 分析条件默认值
// (analyzeConditional analyzes conditional default values)
func (dva *DefaultValueAnalyzer) analyzeConditional(c *ConditionalConfig) {
	fmt.Printf("  Mode: '%s' (expected: 'development')\n", c.Mode)
	fmt.Printf("  Environment: '%s' (expected: 'dev')\n", c.Environment)
	
	if c.DevConfig != nil {
		fmt.Printf("  Dev Debug: %t (expected: true)\n", c.DevConfig.Debug)
		fmt.Printf("  Dev Log Level: '%s' (expected: 'debug')\n", c.DevConfig.LogLevel)
		fmt.Printf("  Dev Hot Reload: %t (expected: true)\n", c.DevConfig.HotReload)
	} else {
		fmt.Println("  DevConfig: nil")
	}
	
	if c.ProdConfig != nil {
		fmt.Printf("  Prod Debug: %t (expected: false)\n", c.ProdConfig.Debug)
		fmt.Printf("  Prod Log Level: '%s' (expected: 'info')\n", c.ProdConfig.LogLevel)
		fmt.Printf("  Prod Max Connections: %d (expected: 1000)\n", c.ProdConfig.MaxConnections)
	} else {
		fmt.Println("  ProdConfig: nil")
	}
}

// analyzeComplex 分析复杂类型默认值
// (analyzeComplex analyzes complex type default values)
func (dva *DefaultValueAnalyzer) analyzeComplex(c *ComplexConfig) {
	if c.OptionalString != nil {
		fmt.Printf("  Optional String: '%s'\n", *c.OptionalString)
	} else {
		fmt.Println("  Optional String: nil (pointers are nil by default)")
	}
	
	if c.OptionalInt != nil {
		fmt.Printf("  Optional Int: %d\n", *c.OptionalInt)
	} else {
		fmt.Println("  Optional Int: nil (pointers are nil by default)")
	}
	
	if c.OptionalBool != nil {
		fmt.Printf("  Optional Bool: %t\n", *c.OptionalBool)
	} else {
		fmt.Println("  Optional Bool: nil (pointers are nil by default)")
	}
	
	fmt.Printf("  Priority: '%s' (expected: 'medium')\n", c.Priority)
	fmt.Printf("  Protocol: '%s' (expected: 'https')\n", c.Protocol)
	
	if c.CustomInterface != nil {
		fmt.Printf("  Custom Interface: %v (type: %s)\n", c.CustomInterface, reflect.TypeOf(c.CustomInterface))
	} else {
		fmt.Println("  Custom Interface: nil (interfaces are nil by default)")
	}
}

// demonstrateDefaultBehaviors 演示不同的默认值行为
// (demonstrateDefaultBehaviors demonstrates different default value behaviors)
func demonstrateDefaultBehaviors() {
	fmt.Println("=== Demonstrating Default Value Behaviors ===")
	fmt.Println()
	
	// 1. 无配置文件，仅默认值 (No config file, defaults only)
	fmt.Println("1. Loading with defaults only (no config file):")
	var cfg1 DefaultValuesConfig
	err := config.LoadConfig(&cfg1, 
		config.WithConfigFile("non-existent.yaml", "yaml"))
	if err != nil {
		fmt.Printf("   Expected error (file not found): %v\n", err)
	}
	
	// 手动初始化一些嵌套结构以演示默认值 (Manually initialize some nested structs to demo defaults)
	cfg1.BasicTypes = &BasicTypesConfig{}
	cfg1.Collections = &CollectionsConfig{}
	
	analyzer := NewDefaultValueAnalyzer()
	analyzer.AnalyzeDefaults(&cfg1)
	
	fmt.Println()
	
	// 2. 带配置文件 (With config file)
	fmt.Println("2. Loading with config file (if exists):")
	var cfg2 DefaultValuesConfig
	err = config.LoadConfig(&cfg2,
		config.WithConfigFile("config.yaml", "yaml"))
	if err != nil {
		fmt.Printf("   Config file error: %v\n", err)
		fmt.Println("   Using defaults for demonstration...")
		// 继续使用默认值演示 (Continue with defaults for demo)
		cfg2 = cfg1
	} else {
		fmt.Println("   ✓ Config file loaded successfully")
		analyzer.AnalyzeDefaults(&cfg2)
	}
}

// printDefaultValueBestPractices 打印默认值最佳实践
// (printDefaultValueBestPractices prints default value best practices)
func printDefaultValueBestPractices() {
	fmt.Println("=== Default Value Best Practices ===")
	fmt.Println()
	
	fmt.Println("1. Struct Tag Defaults:")
	fmt.Println("   ✓ Use `default:\"value\"` for simple types")
	fmt.Println("   ✓ Use comma-separated values for slices: `default:\"a,b,c\"`")
	fmt.Println("   ✓ Use ISO 8601 durations: `default:\"30s\"`, `default:\"1h30m\"`")
	fmt.Println("   ✓ Use string representation for all types")
	fmt.Println()
	
	fmt.Println("2. Nested Structures:")
	fmt.Println("   ⚠ Pointer fields (*struct) are nil by default")
	fmt.Println("   ✓ Initialize nested pointers explicitly if needed")
	fmt.Println("   ✓ Consider using embedded structs instead of pointers")
	fmt.Println("   ✓ Use factory functions for complex initialization")
	fmt.Println()
	
	fmt.Println("3. Collections and Maps:")
	fmt.Println("   ✓ Slices can use comma-separated default values")
	fmt.Println("   ⚠ Maps need special handling (no direct tag support)")
	fmt.Println("   ✓ Consider using factory functions for map initialization")
	fmt.Println("   ✓ Document expected map structure clearly")
	fmt.Println()
	
	fmt.Println("4. Optional Values:")
	fmt.Println("   ✓ Use pointers for truly optional fields")
	fmt.Println("   ✓ Use empty string/zero values for fields with defaults")
	fmt.Println("   ✓ Document whether nil/zero means 'disabled' or 'use default'")
	fmt.Println("   ✓ Provide helper functions to check if optional values are set")
	fmt.Println()
	
	fmt.Println("5. Environment-Specific Defaults:")
	fmt.Println("   ✓ Use different default values for dev/staging/prod")
	fmt.Println("   ✓ Consider using conditional logic in initialization")
	fmt.Println("   ✓ Document environment-specific behaviors")
	fmt.Println("   ✓ Validate defaults make sense for target environment")
	fmt.Println()
}

func main() {
	fmt.Println("=== Default Values Example ===")
	fmt.Println("This example demonstrates comprehensive default value handling patterns.")
	fmt.Println()
	
	// 1. 初始化日志 (Initialize logging)
	logOpts := log.NewOptions()
	logOpts.Level = "info"
	logOpts.Format = "text"
	logOpts.EnableColor = true
	log.Init(logOpts)
	logger := log.Std()
	
	// 2. 演示默认值行为 (Demonstrate default value behaviors)
	demonstrateDefaultBehaviors()
	
	// 3. 打印最佳实践 (Print best practices)
	fmt.Println()
	printDefaultValueBestPractices()
	
	logger.Info("Default values example completed successfully")
	fmt.Println("=== Example completed successfully ===")
} 