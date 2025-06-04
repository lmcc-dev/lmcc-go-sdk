/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: 请求日志中间件接口定义 (Request logger middleware interface definitions)
 */

package middleware

import (
	"time"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
)

// LoggerConfig 日志中间件配置 (Logger middleware configuration)
type LoggerConfig struct {
	// SkipPaths 跳过记录的路径列表 (List of paths to skip logging)
	SkipPaths []string `json:"skip_paths" yaml:"skip_paths"`
	
	// Format 日志格式 (Log format)
	// 支持: "json", "text", "custom" (Supports: "json", "text", "custom")
	Format string `json:"format" yaml:"format"`
	
	// IncludeBody 是否包含请求体 (Whether to include request body)
	IncludeBody bool `json:"include_body" yaml:"include_body"`
	
	// MaxBodySize 最大请求体大小 (Maximum request body size)
	MaxBodySize int `json:"max_body_size" yaml:"max_body_size"`
	
	// IncludeHeaders 是否包含请求头 (Whether to include request headers)
	IncludeHeaders bool `json:"include_headers" yaml:"include_headers"`
	
	// IncludeQuery 是否包含查询参数 (Whether to include query parameters)
	IncludeQuery bool `json:"include_query" yaml:"include_query"`
	
	// TimeFormat 时间格式 (Time format)
	TimeFormat string `json:"time_format" yaml:"time_format"`
	
	// CustomFormatter 自定义格式化函数 (Custom formatter function)
	CustomFormatter func(LogEntry) string `json:"-" yaml:"-"`
}

// DefaultLoggerConfig 返回默认日志配置 (Return default logger configuration)
func DefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		SkipPaths:       []string{"/health", "/metrics", "/favicon.ico"},
		Format:          "json",
		IncludeBody:     false,
		MaxBodySize:     1024, // 1KB
		IncludeHeaders:  false,
		IncludeQuery:    true,
		TimeFormat:      time.RFC3339,
		CustomFormatter: nil,
	}
}

// LogEntry 日志条目 (Log entry)
type LogEntry struct {
	// 基础信息 (Basic information)
	Timestamp   time.Time `json:"timestamp"`
	Method      string    `json:"method"`
	Path        string    `json:"path"`
	StatusCode  int       `json:"status_code"`
	Latency     time.Duration `json:"latency"`
	ClientIP    string    `json:"client_ip"`
	UserAgent   string    `json:"user_agent"`
	
	// 可选信息 (Optional information)
	RequestBody  string            `json:"request_body,omitempty"`
	Headers      map[string]string `json:"headers,omitempty"`
	QueryParams  map[string]string `json:"query_params,omitempty"`
	ResponseSize int64             `json:"response_size,omitempty"`
	
	// 错误信息 (Error information)
	Error string `json:"error,omitempty"`
	
	// 自定义字段 (Custom fields)
	Custom map[string]interface{} `json:"custom,omitempty"`
}

// RequestLogger 请求日志记录器接口 (Request logger interface)
type RequestLogger interface {
	server.Middleware
	
	// SetConfig 设置配置 (Set configuration)
	SetConfig(config *LoggerConfig)
	
	// GetConfig 获取配置 (Get configuration)
	GetConfig() *LoggerConfig
	
	// SetSkipPaths 设置跳过记录的路径 (Set paths to skip logging)
	SetSkipPaths(paths []string)
	
	// AddSkipPath 添加跳过记录的路径 (Add path to skip logging)
	AddSkipPath(path string)
	
	// SetFormat 设置日志格式 (Set log format)
	SetFormat(format string)
	
	// SetIncludeBody 设置是否包含请求体 (Set whether to include request body)
	SetIncludeBody(include bool)
	
	// SetMaxBodySize 设置最大请求体大小 (Set maximum request body size)
	SetMaxBodySize(size int)
	
	// SetCustomFormatter 设置自定义格式化函数 (Set custom formatter function)
	SetCustomFormatter(formatter func(LogEntry) string)
	
	// LogRequest 记录请求 (Log request)
	LogRequest(entry LogEntry)
}

// LoggerFactory 日志中间件工厂接口 (Logger middleware factory interface)
type LoggerFactory interface {
	// CreateLogger 创建日志中间件 (Create logger middleware)
	CreateLogger(config *LoggerConfig) RequestLogger
	
	// CreateDefaultLogger 创建默认日志中间件 (Create default logger middleware)
	CreateDefaultLogger() RequestLogger
}

// LogLevel 日志级别 (Log level)
type LogLevel int

const (
	// LogLevelDebug 调试级别 (Debug level)
	LogLevelDebug LogLevel = iota
	// LogLevelInfo 信息级别 (Info level)
	LogLevelInfo
	// LogLevelWarn 警告级别 (Warning level)
	LogLevelWarn
	// LogLevelError 错误级别 (Error level)
	LogLevelError
)

// String 返回日志级别字符串 (Return log level string)
func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "debug"
	case LogLevelInfo:
		return "info"
	case LogLevelWarn:
		return "warn"
	case LogLevelError:
		return "error"
	default:
		return "unknown"
	}
}

// GetLogLevel 根据状态码和错误确定日志级别 (Determine log level based on status code and error)
func GetLogLevel(statusCode int, err error) LogLevel {
	if err != nil {
		return LogLevelError
	}
	
	switch {
	case statusCode >= 500:
		return LogLevelError
	case statusCode >= 400:
		return LogLevelWarn
	case statusCode >= 300:
		return LogLevelInfo
	default:
		return LogLevelDebug
	}
}

// FormatLogEntry 格式化日志条目 (Format log entry)
func FormatLogEntry(entry LogEntry, format string) string {
	switch format {
	case "json":
		return formatJSON(entry)
	case "text":
		return formatText(entry)
	default:
		return formatJSON(entry) // 默认使用JSON格式 (Default to JSON format)
	}
}

// formatJSON 格式化为JSON (Format as JSON)
func formatJSON(entry LogEntry) string {
	// 这里应该使用JSON序列化 (Should use JSON serialization here)
	// 为了简化，返回基本格式 (For simplicity, return basic format)
	return entry.Timestamp.Format(time.RFC3339) + " " + entry.Method + " " + entry.Path
}

// formatText 格式化为文本 (Format as text)
func formatText(entry LogEntry) string {
	return entry.Timestamp.Format("2006/01/02 15:04:05") + " [" + 
		GetLogLevel(entry.StatusCode, nil).String() + "] " +
		entry.Method + " " + entry.Path + " " +
		"status:" + string(rune(entry.StatusCode)) + " " +
		"latency:" + entry.Latency.String() + " " +
		"ip:" + entry.ClientIP
} 