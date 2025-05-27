/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package log

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

const (
	// FormatJSON 表示 JSON 输出格式。(FormatJSON represents the JSON output format.)
	FormatJSON = "json"
	// FormatText 表示 Text 输出格式。(FormatText represents the Text output format.)
	FormatText = "text"
	// FormatKeyValue 表示 key=value 输出格式。(FormatKeyValue represents the key=value output format.)
	FormatKeyValue = "keyvalue"
)

// Options 定义了日志配置选项。(Options defines configuration options for the logger.)
// 它遵循选项模式，允许用户自定义日志行为。
// (It follows the options pattern, allowing users to customize logging behavior.)
type Options struct {
	// OutputPaths 指定了日志的输出路径，可以是 stdout、stderr 或文件路径。
	// (OutputPaths specifies the log output paths. It can be stdout, stderr, or file paths.)
	OutputPaths []string `json:"output-paths" mapstructure:"outputPaths"`

	// ErrorOutputPaths 指定了内部错误日志的输出路径。
	// (ErrorOutputPaths specifies the output paths for internal error logs.)
	ErrorOutputPaths []string `json:"error-output-paths" mapstructure:"errorOutputPaths"`

	// Level 指定了日志级别，例如 "debug", "info", "warn", "error", "fatal"。
	// (Level specifies the log level, e.g., "debug", "info", "warn", "error", "fatal".)
	Level string `json:"level" mapstructure:"level"`

	// Format 指定了日志的输出格式，"json" 或 "text"。
	// (Format specifies the log output format, either "json" or "text".)
	Format string `json:"format" mapstructure:"format"`

	// DisableCaller 禁用在日志条目中包含调用者信息（文件和行号）。
	// (DisableCaller disables including caller information (file and line number) in log entries.)
	DisableCaller bool `json:"disable-caller" mapstructure:"disable-caller"`

	// DisableStacktrace 禁用在 Error 级别及以上的日志中自动记录堆栈跟踪。
	// (DisableStacktrace disables automatic stacktrace recording on logs at Error level and above.)
	DisableStacktrace bool `json:"disable-stacktrace" mapstructure:"disable-stacktrace"`

	// StacktraceLevel 指定了开始记录堆栈跟踪的最低日志级别。例如 "warn" 表示在 "warn", "error", "fatal" 级别记录。
	// (StacktraceLevel specifies the minimum log level to start recording stacktraces. e.g., "warn" means record at "warn", "error", "fatal".)
	StacktraceLevel string `json:"stacktrace-level" mapstructure:"stacktrace-level"`

	// EnableColor 在 Text 格式的日志输出中启用颜色。
	// (EnableColor enables color in Text format log output.)
	EnableColor bool `json:"enable-color" mapstructure:"enable-color"`

	// Development 设置日志库为开发模式，会改变 DPanicLevel 的行为并在日志中包含更多调试信息。
	// (Development sets the logger to development mode, which changes DPanicLevel behavior and includes more debugging info in logs.)
	Development bool `json:"development" mapstructure:"development"`

	// Name 是日志记录器的名称。
	// (Name is the name of the logger.)
	Name string `json:"name" mapstructure:"name"`

	// TimeFormat 定义了日志中时间戳的格式。如果为空，则使用 zap 的默认格式。
	// (TimeFormat defines the format for timestamps in logs. If empty, zap's default is used.)
	TimeFormat string `json:"time-format" mapstructure:"time-format"`

	// EncoderConfig 允许用户提供自定义的 zapcore.EncoderConfig。
	// 如果为 nil，将根据其他选项（如 Format, EnableColor, TimeFormat）自动生成配置。
	// (EncoderConfig allows the user to provide a custom zapcore.EncoderConfig.
	// If nil, a config will be generated automatically based on other options like Format, EnableColor, TimeFormat.)
	EncoderConfig *zapcore.EncoderConfig `json:"-" mapstructure:"-"` // 通常不通过配置文件设置，而是代码中设置

	// --- 日志轮转选项 (Log Rotation Options) ---

	// LogRotateMaxSize 是日志文件的最大大小（以 MB 为单位）。
	// (LogRotateMaxSize is the maximum size in megabytes of the log file before it gets rotated.)
	LogRotateMaxSize int `json:"log-rotate-max-size" mapstructure:"log-rotate-max-size"`

	// LogRotateMaxBackups 是要保留的旧日志文件的最大数量。
	// (LogRotateMaxBackups is the maximum number of old log files to retain.)
	LogRotateMaxBackups int `json:"log-rotate-max-backups" mapstructure:"log-rotate-max-backups"`

	// LogRotateMaxAge 是根据文件名中的时间戳保留旧日志文件的最大天数。
	// (LogRotateMaxAge is the maximum number of days to retain old log files based on the timestamp encoded in their filename.)
	LogRotateMaxAge int `json:"log-rotate-max-age" mapstructure:"log-rotate-max-age"`

	// LogRotateCompress 决定是否压缩（gzip）轮转的日志文件。
	// (LogRotateCompress determines if the rotated log files should be compressed (gzip).)
	LogRotateCompress bool `json:"log-rotate-compress" mapstructure:"log-rotate-compress"`

	// ContextKeys 是用户希望从 context 中自动提取并添加到日志字段的额外键列表。
	// 这些键的类型应该与 context.WithValue 中使用的键类型完全匹配。
	// (ContextKeys is a list of additional keys that the user wants to automatically extract
	// from the context and add to the log fields. The type of these keys should exactly match
	// the type of keys used in context.WithValue.)
	ContextKeys []any `json:"context-keys" mapstructure:"context-keys"`
}

// NewOptions 创建具有默认值的日志选项 (creates logging options with default values)
func NewOptions() *Options {
	return &Options{
		Level:               zapcore.InfoLevel.String(),     // 默认级别 info (Default level info)
		DisableCaller:       false,
		DisableStacktrace:   false,                          // 默认启用堆栈跟踪（但由 StacktraceLevel 控制）
		StacktraceLevel:     zapcore.ErrorLevel.String(),  // 默认从 Error 级别开始记录堆栈
		Format:              FormatJSON,                     // 默认格式 json (Default format json)
		EnableColor:         false,                          // 默认禁用颜色 (Color disabled by default)
		Development:         false,                          // 默认生产模式 (Production mode by default)
		OutputPaths:         []string{"stdout"},             // 默认输出到 stdout (Default output to stdout)
		ErrorOutputPaths:    []string{"stderr"},             // 默认错误输出到 stderr (Default error output to stderr)
		TimeFormat:          "",                             // 默认为空，JSON 为 ISO8601，Text 有特定默认值 (Default empty, JSON is ISO8601, Text has specific default)
		EncoderConfig:       nil,                            // 默认为 nil，将自动生成
		LogRotateMaxSize:    100,                            // 默认 100 MB (Default 100 MB)
		LogRotateMaxAge:     7,                              // 默认保留 7 天 (Default retention 7 days)
		LogRotateMaxBackups: 5,                             // 默认保留 5 个备份 (Default retain 5 backups)
		LogRotateCompress:   false,                          // 默认不压缩 (No compression by default)
		ContextKeys:         nil,                            // 默认不提取额外键 (No extra keys by default, now type is []any)
	}
}

// Validate 验证日志选项是否有效。
// (Validate validates if the logging options are valid.)
func (o *Options) Validate() []error {
	var errs []error

	// 验证 Level
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
		errs = append(errs, fmt.Errorf("invalid log level '%s': %w", o.Level, err))
	}

	// 验证 Format
	if o.Format != FormatJSON && o.Format != FormatText && o.Format != FormatKeyValue {
		errs = append(errs, fmt.Errorf("invalid log format '%s', must be '%s', '%s', or '%s'", o.Format, FormatJSON, FormatText, FormatKeyValue))
	}

	// 验证 StacktraceLevel
	var stacktraceZapLevel zapcore.Level
	if err := stacktraceZapLevel.UnmarshalText([]byte(o.StacktraceLevel)); err != nil {
		errs = append(errs, fmt.Errorf("invalid stacktrace level '%s': %w", o.StacktraceLevel, err))
	}

	// 其他验证可以根据需要添加，例如 OutputPaths 是否有效等。

	return errs
}

// AddFlags 将日志选项相关的标志添加到指定的 pflag.FlagSet
// func (o *Options) AddFlags(fs *pflag.FlagSet) {
// 	fs.StringVar(&o.Level, "log.level", o.Level, "Minimum log output level.")
// 	fs.BoolVar(&o.DisableCaller, "log.disable-caller", o.DisableCaller, "Disable output caller info.")
// 	fs.BoolVar(&o.DisableStacktrace, "log.disable-stacktrace", o.DisableStacktrace, "Disable stacktrace for error logs.")
// 	fs.StringVar(&o.StacktraceLevel, "log.stacktrace-level", o.StacktraceLevel, "Minimum log level to record stacktrace.")
// 	fs.StringVar(&o.Format, "log.format", o.Format, "Log output format, can be 'json' or 'text'.")
// 	fs.BoolVar(&o.EnableColor, "log.enable-color", o.EnableColor, "Enable output ansi colors.")
// 	fs.StringVar(&o.TimeFormat, "log.time-format", o.TimeFormat, "Format for timestamps in logs.")
// 	fs.StringSliceVar(&o.OutputPaths, "log.output-paths", o.OutputPaths, "Log output paths.")
// 	fs.StringSliceVar(&o.ErrorOutputPaths, "log.error-output-paths", o.ErrorOutputPaths, "Error log output paths.")
// 	fs.BoolVar(&o.Development, "log.development", o.Development, "Enable development mode for logging.")
// 	fs.StringVar(&o.Name, "log.name", o.Name, "The name of the logger.")
// // 添加轮转相关的标志 (Add rotation related flags)
// 	fs.IntVar(&o.LogRotateMaxSize, "log.rotate.max-size", o.LogRotateMaxSize, "Maximum size in megabytes of the log file before rotation.")
// 	fs.IntVar(&o.LogRotateMaxBackups, "log.rotate.max-backups", o.LogRotateMaxBackups, "Maximum number of old log files to retain.")
// 	fs.IntVar(&o.LogRotateMaxAge, "log.rotate.max-age", o.LogRotateMaxAge, "Maximum number of days to retain old log files.")
// 	fs.BoolVar(&o.LogRotateCompress, "log.rotate.compress", o.LogRotateCompress, "Compress rotated log files using gzip.")
// }
