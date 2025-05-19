/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package config

import (
	"time"

	"github.com/spf13/viper" // Import viper
)

// SectionChangeCallback defines the type for callback functions invoked on configuration change for a specific section.
// (SectionChangeCallback 定义了当特定配置节发生变更时调用的回调函数类型。)
type SectionChangeCallback func(v *viper.Viper) error

// Manager defines the interface for a configuration manager.
// (Manager 定义了配置管理器的接口。)
// It provides methods to access the underlying Viper instance and register callbacks for configuration changes.
type Manager interface {
	// GetViperInstance returns the underlying Viper instance used by the manager.
	// (GetViperInstance 返回管理器使用的底层 Viper 实例。)
	GetViperInstance() *viper.Viper

	// RegisterCallback registers a general configuration change callback.
	// This callback is triggered when any part of the configuration changes and receives the Viper instance and the decoded config object.
	// (RegisterCallback 注册一个通用的配置变更回调。
	// 当配置的任何部分发生更改时，将触发此回调，并接收 Viper 实例和解码后的配置对象。)
	RegisterCallback(callback func(v *viper.Viper, cfg any) error) // Matches existing ConfigChangeCallback

	// RegisterSectionChangeCallback registers a callback for changes in a specific configuration section.
	// The callback receives the Viper instance and is responsible for unmarshalling its specific section.
	// (RegisterSectionChangeCallback 注册特定配置节变更的回调。
	// 回调接收 Viper 实例，并负责解组其特定节。)
	RegisterSectionChangeCallback(sectionKey string, callback SectionChangeCallback)

	// TODO: Consider adding StopWatch() or similar to control the watcher lifecycle if needed.
}

// Config 是 SDK 提供的基础配置结构体 (Base configuration struct provided by the SDK)
// 用户可以通过嵌入此结构体来扩展自定义配置 (Users can extend this by embedding it in their own config struct)
type Config struct {
	Server   *ServerConfig   `mapstructure:"server"`
	Log      *LogConfig      `mapstructure:"log"`
	Database *DatabaseConfig `mapstructure:"database"`
	Tracing  *TracingConfig  `mapstructure:"tracing"`
	Metrics  *MetricsConfig  `mapstructure:"metrics"`
	// 注意：用户可以在他们自己的结构体中添加更多字段 (Note: Users can add more fields in their own struct)
}

// ServerConfig 服务器相关配置 (Server related configuration)
type ServerConfig struct {
	Host                    string        `mapstructure:"host" default:"0.0.0.0"`
	Port                    int           `mapstructure:"port" default:"8080"`
	Mode                    string        `mapstructure:"mode" default:"release"`
	ReadTimeout             time.Duration `mapstructure:"readTimeout" default:"5s"`
	WriteTimeout            time.Duration `mapstructure:"writeTimeout" default:"10s"`
	GracefulShutdownTimeout time.Duration `mapstructure:"gracefulShutdownTimeout" default:"10s"`
}

// LogConfig 日志相关配置 (Logging related configuration)
type LogConfig struct {
	Level             string   `mapstructure:"level" default:"info"`
	Format            string   `mapstructure:"format" default:"text"`
	EnableColor       bool     `mapstructure:"enableColor" default:"false"`
	Output            string   `mapstructure:"output" default:"stdout"`
	OutputPaths       []string `mapstructure:"outputPaths"`
	ErrorOutput       string   `mapstructure:"errorOutput" default:"stderr"`
	ErrorOutputPaths  []string `mapstructure:"errorOutputPaths"`
	Filename          string   `mapstructure:"filename" default:"app.log"`
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

// DatabaseConfig 数据库相关配置 (Database related configuration)
type DatabaseConfig struct {
	Type            string        `mapstructure:"type" default:"mysql"`
	Host            string        `mapstructure:"host" default:"localhost"`
	Port            int           `mapstructure:"port"`     // No default, should be specified if needed
	User            string        `mapstructure:"user"`     // No default
	Password        string        `mapstructure:"password"` // No default
	DBName          string        `mapstructure:"dbName"`   // No default
	MaxIdleConns    int           `mapstructure:"maxIdleConns" default:"10"`
	MaxOpenConns    int           `mapstructure:"maxOpenConns" default:"100"`
	ConnMaxLifetime time.Duration `mapstructure:"connMaxLifetime" default:"1h"`
}

// TracingConfig 链路追踪相关配置 (Tracing related configuration)
type TracingConfig struct {
	Enabled      bool    `mapstructure:"enabled" default:"false"`
	Provider     string  `mapstructure:"provider" default:"jaeger"`
	Endpoint     string  `mapstructure:"endpoint"` // No default
	SamplerType  string  `mapstructure:"samplerType" default:"const"`
	SamplerParam float64 `mapstructure:"samplerParam" default:"1"`
}

// MetricsConfig 监控指标相关配置 (Metrics related configuration)
type MetricsConfig struct {
	Enabled  bool   `mapstructure:"enabled" default:"false"`
	Provider string `mapstructure:"provider" default:"prometheus"`
	Port     int    `mapstructure:"port" default:"9090"`
	Path     string `mapstructure:"path" default:"/metrics"`
}
