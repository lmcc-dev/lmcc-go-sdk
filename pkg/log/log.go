/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package log

import (
	"context"
	"fmt"
	"io"      //确保导入 io 包
	"os"      // Needed for os.Stdout, os.Stderr
	"strings" // Added for strings.Contains
	"sync"
	"sync/atomic" // Added for atomic.Pointer

	lmccerrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors" // 导入 lmccerrors 包 (Import lmccerrors package)
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 定义了 SDK 的标准日志接口。
// (Logger defines the standard logging interface for the SDK.)
// 它包装了 zap.Logger，并提供了上下文感知和格式化的日志方法。
// (It wraps zap.Logger and provides context-aware and formatted logging methods.)
type Logger interface {
	// Debug 使用 fmt.Sprint 风格记录一条 Debug 级别的消息。
	// (Debug logs a message at DebugLevel using fmt.Sprint style.)
	Debug(args ...any)
	// Debugf 使用 fmt.Sprintf 风格记录一条 Debug 级别的消息。
	// (Debugf logs a message at DebugLevel using fmt.Sprintf style.)
	Debugf(template string, args ...any)
	// Debugw 记录一条 Debug 级别的消息，并附带键值对。
	// (Debugw logs a message at DebugLevel with key-value pairs.)
	Debugw(msg string, keysAndValues ...any)

	// Info 使用 fmt.Sprint 风格记录一条 Info 级别的消息。
	// (Info logs a message at InfoLevel using fmt.Sprint style.)
	Info(args ...any)
	// Infof 使用 fmt.Sprintf 风格记录一条 Info 级别的消息。
	// (Infof logs a message at InfoLevel using fmt.Sprintf style.)
	Infof(template string, args ...any)
	// Infow 记录一条 Info 级别的消息，并附带键值对。
	// (Infow logs a message at InfoLevel with key-value pairs.)
	Infow(msg string, keysAndValues ...any)

	// Warn 使用 fmt.Sprint 风格记录一条 Warn 级别的消息。
	// (Warn logs a message at WarnLevel using fmt.Sprint style.)
	Warn(args ...any)
	// Warnf 使用 fmt.Sprintf 风格记录一条 Warn 级别的消息。
	// (Warnf logs a message at WarnLevel using fmt.Sprintf style.)
	Warnf(template string, args ...any)
	// Warnw 记录一条 Warn 级别的消息，并附带键值对。
	// (Warnw logs a message at WarnLevel with key-value pairs.)
	Warnw(msg string, keysAndValues ...any)

	// Error 使用 fmt.Sprint 风格记录一条 Error 级别的消息。
	// (Error logs a message at ErrorLevel using fmt.Sprint style.)
	Error(args ...any)
	// Errorf 使用 fmt.Sprintf 风格记录一条 Error 级别的消息。
	// (Errorf logs a message at ErrorLevel using fmt.Sprintf style.)
	Errorf(template string, args ...any)
	// Errorw 记录一条 Error 级别的消息，并附带键值对。
	// (Errorw logs a message at ErrorLevel with key-value pairs.)
	Errorw(msg string, keysAndValues ...any)

	// Fatal 使用 fmt.Sprint 风格记录一条 Fatal 级别的消息，然后调用 os.Exit(1)。
	// (Fatal logs a message at FatalLevel using fmt.Sprint style, then calls os.Exit(1).)
	Fatal(args ...any)
	// Fatalf 使用 fmt.Sprintf 风格记录一条 Fatal 级别的消息，然后调用 os.Exit(1)。
	// (Fatalf logs a message at FatalLevel using fmt.Sprintf style, then calls os.Exit(1).)
	Fatalf(template string, args ...any)
	// Fatalw 记录一条 Fatal 级别的消息，并附带键值对，然后调用 os.Exit(1)。
	// (Fatalw logs a message at FatalLevel with key-value pairs, then calls os.Exit(1).)
	Fatalw(msg string, keysAndValues ...any)

	// DPanic 使用 fmt.Sprint 风格记录一条 DPanic 级别的消息。
	// 在开发模式下会 panic，在生产模式下则记录 Error 级别日志。
	// (DPanic logs a message at DPanicLevel using fmt.Sprint style.)
	// (It panics in development mode and logs at ErrorLevel in production.)
	DPanic(args ...any)
	// DPanicf 使用 fmt.Sprintf 风格记录一条 DPanic 级别的消息。
	// (DPanicf logs a message at DPanicLevel using fmt.Sprintf style.)
	DPanicf(template string, args ...any)
	// DPanicw 记录一条 DPanic 级别的消息，并附带键值对。
	// (DPanicw logs a message at DPanicLevel with key-value pairs.)
	DPanicw(msg string, keysAndValues ...any)

	// Ctx 使用 fmt.Sprint 风格记录一条 Info 级别的消息，并从 context 中提取字段。
	// (Ctx logs a message at InfoLevel using fmt.Sprint style, extracting fields from the context.)
	Ctx(ctx context.Context, args ...any)
	// Ctxf 使用 fmt.Sprintf 风格记录一条 Info 级别的消息，并从 context 中提取字段。
	// (Ctxf logs a message at InfoLevel using fmt.Sprintf style, extracting fields from the context.)
	Ctxf(ctx context.Context, template string, args ...any)
	// Ctxw 记录一条 Info 级别的消息，并附带键值对，同时从 context 中提取字段。
	// (Ctxw logs a message at InfoLevel with key-value pairs, also extracting fields from the context.)
	Ctxw(ctx context.Context, msg string, keysAndValues ...any)

	// Sync 将所有缓冲的日志条目刷新到底层写入器。
	// (Sync flushes any buffered log entries to the underlying writers.)
	Sync() error

	// WithValues 向日志记录器添加一组键值对上下文。
	// (WithValues adds a set of key-value pairs context to the logger.)
	WithValues(keysAndValues ...any) Logger

	// WithName 向日志记录器的名称添加一个新元素。
	// (WithName adds a new element to the logger's name.)
	WithName(name string) Logger

	// GetZapLogger 返回底层的 zap.Logger。
	// (GetZapLogger returns the underlying zap.Logger.)
	GetZapLogger() *zap.Logger

	// --- Contextual Logging ---
	// CtxDebugf uses fmt.Sprintf to log a templated message at DebugLevel.
	// It extracts fields from the context using pre-configured keys.
	CtxDebugf(ctx context.Context, template string, args ...interface{})
	// CtxInfof uses fmt.Sprintf to log a templated message at InfoLevel.
	// It extracts fields from the context using pre-configured keys.
	CtxInfof(ctx context.Context, template string, args ...interface{})
	// CtxWarnf uses fmt.Sprintf to log a templated message at WarnLevel.
	// It extracts fields from the context using pre-configured keys.
	CtxWarnf(ctx context.Context, template string, args ...interface{})
	// CtxErrorf uses fmt.Sprintf to log a templated message at ErrorLevel.
	// It extracts fields from the context using pre-configured keys.
	CtxErrorf(ctx context.Context, template string, args ...interface{})
	// CtxPanicf uses fmt.Sprintf to log a templated message at PanicLevel, then panics.
	// It extracts fields from the context using pre-configured keys.
	CtxPanicf(ctx context.Context, template string, args ...interface{})
	// CtxFatalf uses fmt.Sprintf to log a templated message at FatalLevel, then calls os.Exit(1).
	// It extracts fields from the context using pre-configured keys.
	CtxFatalf(ctx context.Context, template string, args ...interface{})
}

// logger 是 Logger 接口的 zap 实现。
// (logger is the zap implementation of the Logger interface.)
// 注意：保持 logger 结构体本身不导出，以封装实现细节。
// (Note: Keep the logger struct itself unexported to encapsulate implementation details.)
type logger struct {
	zapLogger *zap.Logger
	opts      *Options // Store applied options
}

// keyValueLogger 是一个包装器，用于在 key=value 格式下处理 WithValues
// (keyValueLogger is a wrapper for handling WithValues in key=value format)
type keyValueLogger struct {
	baseLogger *logger
	fields     []any
}

// 确保 logger 实现了 Logger 接口。
var _ Logger = (*logger)(nil)
var _ Logger = (*keyValueLogger)(nil)

// FormatJSON/FormatText are defined in options.go

var (
	// std 使用 atomic.Pointer 来存储全局 Logger 实例，以支持原子化更新。
	// (std uses atomic.Pointer to store the global Logger instance, enabling atomic updates for thread-safe reconfiguration.)
	std atomic.Pointer[logger] // Changed from Logger to *atomic.Pointer[logger]

	// mu 仅用于保护首次初始化时的竞态条件 (如果多个 goroutine 同时调用 Std() 且 std 未初始化)。
	// (mu is used only to protect against race conditions during the first initialization if multiple goroutines call Std() concurrently when std is uninitialized.)
	// 一旦 std 被设置，后续的 ReconfigureGlobalLogger 将通过原子操作进行。
	// (Once std is set, subsequent ReconfigureGlobalLogger calls will be atomic.)
	mu sync.Mutex
)

// Init 使用给定的选项初始化全局日志记录器 std。
// 此函数是线程安全的。如果 std 已经初始化，它将被新的配置原子地覆盖。
// 如果初始化失败，它将 panic。
// (Init initializes the global logger std with the given options.)
// (This function is thread-safe. If std is already initialized, it will be atomically overwritten with the new configuration.)
// (If initialization fails, it will panic.)
func Init(opts *Options) {
	l, err := newLogger(opts)
	if err != nil {
		// 对于 Init 失败，我们选择 panic，因为日志系统是基础组件。
		// (For Init failure, we choose to panic as the logging system is a fundamental component.)
		panic(lmccerrors.WithCode(
			lmccerrors.Wrap(err, "failed to initialize global logger with provided options"),
			lmccerrors.ErrLogInitialization,
		))
	}
	std.Store(l)
}

// NewLogger 根据提供的选项创建一个新的 Logger 实例。
// (NewLogger creates a new Logger instance based on the provided options.)
func NewLogger(opts *Options) (Logger, error) {
	l, err := newLogger(opts)
	if err != nil {
		return nil, err // 直接返回 newLogger 的错误 (Directly return error from newLogger)
	}
	return l, nil
}

// Std 返回全局的 Logger 实例。
// 如果 std 未初始化，它将使用默认选项进行初始化。
// 此函数通过双重检查锁定和原子操作确保线程安全。
// 如果在首次延迟初始化期间发生错误，它将 panic。
// (Std returns the global Logger instance.)
// (If std is not initialized, it will be initialized with default options.)
// (This function ensures thread-safety through double-checked locking and atomic operations.)
// (If an error occurs during the first lazy initialization, it will panic.)
func Std() Logger {
	l := std.Load()
	if l == nil {
		mu.Lock()
		defer mu.Unlock()
		// Double-check in case another goroutine initialized it while we were waiting for the lock
		// (再次检查，以防在我们等待锁的时候，另一个 goroutine 已经初始化了它)
		l = std.Load()
		if l == nil {
			var err error
			l, err = newLogger(NewOptions()) // Use default options if not initialized
			if err != nil {
				// 首次默认初始化失败，panic (First default initialization failed, panic)
				panic(lmccerrors.WithCode(
					lmccerrors.Wrap(err, "failed to initialize global logger with default options"),
					lmccerrors.ErrLogInitialization,
				))
			}
			std.Store(l)
		}
	}
	return l
}

// ReconfigureGlobalLogger 使用新的选项原子地重新配置全局日志记录器。
// 此操作是线程安全的。
// (ReconfigureGlobalLogger atomically reconfigures the global logger with new options.)
// (This operation is thread-safe.)
// Parameters:
//   newOpts: 指向新的日志选项的指针。不能为 nil。
//            (Pointer to the new logger options. Must not be nil.)
// Returns:
//   error: 如果 newOpts 为 nil，则返回错误。如果创建新的 logger 实例失败（例如，选项无效），也可能返回错误。
//          (Returns an error if newOpts is nil. May also return an error if creating the new logger instance fails (e.g., invalid options).)
func ReconfigureGlobalLogger(newOpts *Options) error {
	if newOpts == nil {
		// 使用 lmccerrors.NewWithCode 创建错误。
		// (Use lmccerrors.NewWithCode to create an error.)
		return lmccerrors.NewWithCode(lmccerrors.ErrLogOptionInvalid, "cannot reconfigure global logger with nil options")
	}

	newL, err := newLogger(newOpts)
	if err != nil {
		// 将 newLogger 返回的错误包装，以提供重新配置的上下文
		// (Wrap the error returned by newLogger to provide reconfiguration context)
		// 确保 %w 用于正确的错误链包装 (Ensure %w is used for correct error chain wrapping)
		return lmccerrors.WithCode(
			lmccerrors.Wrap(err, "failed to create new logger for reconfiguration"),
			lmccerrors.ErrLogReconfigure,
		)
	}
	std.Store(newL)
	return nil
}

// --- DPanic 系列方法实现 ---
// (DPanic series method implementation)
func (l *logger) DPanic(args ...any) {
	if l.opts.Development {
		l.zapLogger.Sugar().DPanic(args...)
	} else {
		l.zapLogger.Sugar().Error(args...)
	}
}
func (l *logger) DPanicf(template string, args ...any) {
	if l.opts.Development {
		l.zapLogger.Sugar().DPanicf(template, args...)
	} else {
		l.zapLogger.Sugar().Errorf(template, args...)
	}
}
func (l *logger) DPanicw(msg string, keysAndValues ...any) {
	if l.opts.Development {
		l.zapLogger.Sugar().DPanicw(msg, keysAndValues...)
	} else {
		l.zapLogger.Sugar().Errorw(msg, keysAndValues...)
	}
}

// --- 全局 DPanic 系列函数 ---
// DPanic 在全局 logger 上调用 DPanic。
// (DPanic calls DPanic on the global logger.)
func DPanic(args ...any) {
	Std().DPanic(args...)
}

// DPanicf 在全局 logger 上调用 DPanicf。
// (DPanicf calls DPanicf on the global logger.)
func DPanicf(template string, args ...any) {
	Std().DPanicf(template, args...)
}

// DPanicw 在全局 logger 上调用 DPanicw。
// (DPanicw calls DPanicw on the global logger.)
func DPanicw(msg string, keysAndValues ...any) {
	Std().DPanicw(msg, keysAndValues...)
}

// getEncoderConfig 根据选项创建并返回一个 zapcore.EncoderConfig。
// (getEncoderConfig creates and returns a zapcore.EncoderConfig based on the options.)
func getEncoderConfig(opts *Options) zapcore.EncoderConfig {
	// 如果用户提供了自定义的 EncoderConfig，则直接使用它。
	// (If the user provided a custom EncoderConfig, use it directly.)
	if opts.EncoderConfig != nil {
		return *opts.EncoderConfig
	}

	// 根据格式和测试的期望，定制 EncoderConfig
	// (Customize EncoderConfig based on format and test expectations)
	var encoderConfig zapcore.EncoderConfig
	
	if opts.Format == FormatText || opts.Format == FormatKeyValue {
		// 对于文本格式和 key=value 格式，使用开发配置作为基础
		// (For text format and key=value format, use development config as base)
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		// 覆盖一些字段以匹配测试期望
		// (Override some fields to match test expectations)
		encoderConfig.LevelKey = "L"
		encoderConfig.CallerKey = "C"
		encoderConfig.MessageKey = "M"
		encoderConfig.NameKey = "N"
	} else {
		// 对于 JSON 格式，使用生产配置
		// (For JSON format, use production config)
		encoderConfig = zap.NewProductionEncoderConfig()
		encoderConfig.LevelKey = "L"
		encoderConfig.CallerKey = "C"
		encoderConfig.MessageKey = "M"
		encoderConfig.NameKey = "N"
	}

	// 根据格式和颜色选项设置 LevelEncoder
	// (Set LevelEncoder based on format and color options)
	if opts.Format == FormatText || opts.Format == FormatKeyValue {
		if opts.EnableColor {
			encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // 文本格式且启用颜色时，使用大写带颜色的 LevelEncoder (For text format with color enabled, use capital color LevelEncoder)
		} else {
			encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 文本格式但不启用颜色时，使用大写 LevelEncoder (For text format without color, use capital LevelEncoder)
		}
	} else { // 默认为 JSON 格式或其他格式 (Default to JSON format or other formats)
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // JSON 格式使用大写 LevelEncoder (For JSON format, use capital LevelEncoder)
	}

	// 时间编码器 (Time encoder)
	if opts.TimeFormat != "" {
		encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(opts.TimeFormat)
	} else {
		// 默认为 ISO8601 格式，与 zap.NewDevelopmentEncoderConfig() 和当前测试输出一致
		// (Default to ISO8601 format, consistent with zap.NewDevelopmentEncoderConfig() and current test output)
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// 调用者编码器 (Caller encoder)
	if !opts.DisableCaller {
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder // 启用调用者信息时，使用短路径编码器 (When caller info is enabled, use short path encoder)
	} else {
		encoderConfig.EncodeCaller = nil // 明确禁用调用者信息 (Explicitly disable caller info)
	}

	return encoderConfig
}

// newLoggerInternal 是创建 zap.Logger 的核心逻辑，可被 NewLogger 和 NewLoggerWithWriter 复用。
// 它接收 Options 和一个已经构建好的 zapcore.WriteSyncer。
func newLoggerInternal(opts *Options, syncer zapcore.WriteSyncer) (*zap.Logger, *zap.AtomicLevel, error) {
	if opts == nil {
		return nil, nil, lmccerrors.NewWithCode(lmccerrors.ErrLogOptionInvalid, "options cannot be nil for newLoggerInternal")
	}

	// 验证选项 (Validate options first)
	if validationErrs := opts.Validate(); len(validationErrs) > 0 {
		// 将多个验证错误合并为一个 (Combine multiple validation errors into one)
		// 这里简单地取第一个错误，或者可以构造一个更复杂的错误消息
		// (Here, simply take the first error, or a more complex error message can be constructed)
		return nil, nil, lmccerrors.ErrorfWithCode(lmccerrors.ErrLogOptionInvalid, "invalid options: %v", validationErrs)
	}

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		return nil, nil, lmccerrors.WithCode(
			lmccerrors.Wrapf(err, "invalid log level '%s'", opts.Level),
			lmccerrors.ErrLogOptionInvalid,
		)
	}
	atomicLevel := zap.NewAtomicLevelAt(zapLevel)

	encoderConfig := getEncoderConfig(opts) // 使用修正后的 getEncoderConfig
	var encoder zapcore.Encoder
	if opts.Format == FormatJSON {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else if opts.Format == FormatText || opts.Format == FormatKeyValue {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		// Validate() 应该已经捕获了这个问题，但作为防御性检查
		// (Validate() should have caught this, but as a defensive check)
		return nil, nil, lmccerrors.ErrorfWithCode(lmccerrors.ErrLogOptionInvalid, "invalid log format: %s", opts.Format)
	}

	core := zapcore.NewCore(encoder, syncer, atomicLevel)

	var zapOpts []zap.Option
	if !opts.DisableCaller { // 使用 !opts.DisableCaller
		zapOpts = append(zapOpts, zap.AddCaller(), zap.AddCallerSkip(1)) // Skip our wrapper method
	}

	if opts.Development {
		zapOpts = append(zapOpts, zap.Development())
	}

	if !opts.DisableStacktrace {
		var stacktraceLevel zapcore.Level
		if err := stacktraceLevel.UnmarshalText([]byte(opts.StacktraceLevel)); err != nil {
			// 如果 StacktraceLevel 无效，则默认为 ErrorLevel 并记录一个内部警告
			// (If StacktraceLevel is invalid, default to ErrorLevel and log an internal warning)
			// 注意：这里不能使用全局 logger，因为它可能尚未初始化
			// (Note: Cannot use global logger here as it might not be initialized yet)
			fmt.Fprintf(os.Stderr, "Warning: Invalid StacktraceLevel '%s', defaulting to ErrorLevel. Error: %v\\n", opts.StacktraceLevel, err)
			stacktraceLevel = zapcore.ErrorLevel
		}
		zapOpts = append(zapOpts, zap.AddStacktrace(stacktraceLevel))
	}

	// Options 结构体中没有 ZapOptions 字段，移除相关代码
	// if len(opts.ZapOptions) > 0 {
	// zapOpts = append(zapOpts, opts.ZapOptions...)
	// }

	zapL := zap.New(core, zapOpts...)
	return zapL, &atomicLevel, nil
}

// newLogger 是 logger 的内部构造函数。
// (newLogger is the internal constructor for a logger.)
// 它接收 Options 并返回一个配置好的 *logger 实例或错误。
// (It takes Options and returns a configured *logger instance or an error.)
// 注意：这个函数现在返回 (*logger, error) 以便更好地处理错误。
// (Note: This function now returns (*logger, error) for better error handling.)
func newLogger(opts *Options) (*logger, error) { // Changed return type to (*logger, error)
	if opts == nil {
		opts = NewOptions() // 使用默认选项，如果提供的是 nil (Use default options if nil is provided)
	}

	// 获取写入同步器 (Get write syncer)
	writeSyncer, err := getWriteSyncer(opts) // getWriteSyncer will handle OutputPaths
	if err != nil {
		// 返回带有上下文的错误，而不是 panic (Return an error with context instead of panic)
		// 确保返回的错误是 ErrLogInitialization 类型 (Ensure the returned error is of type ErrLogInitialization)
		return nil, lmccerrors.WithCode(
			lmccerrors.Wrap(err, "failed to get write syncer for logger"),
			lmccerrors.ErrLogInitialization,
		)
	}

	zapL, _, err := newLoggerInternal(opts, writeSyncer) // Use newLoggerInternal
	if err != nil {
		// 如果 newLoggerInternal 返回错误，则将其包装并返回
		// (If newLoggerInternal returns an error, wrap and return it)
		return nil, lmccerrors.WithCode(
			lmccerrors.Wrap(err, "failed to create new zap logger"),
			lmccerrors.ErrLogInitialization,
		)
	}

	// 返回包装后的 logger (Return the wrapped logger)
	return &logger{
		zapLogger: zapL,
		opts:      opts, // 存储应用的选项 (Store applied options)
	}, nil
}

// NewLoggerWithWriter 创建一个新的 Logger 实例，将其输出写入提供的 io.Writer。
// 主要用于测试目的。
// (NewLoggerWithWriter creates a new Logger instance that writes its output to the provided io.Writer.)
// (This is primarily intended for testing purposes.)
func NewLoggerWithWriter(opts *Options, writer io.Writer) Logger {
	if opts == nil {
		opts = NewOptions() // 使用默认选项 (Use default options)
	}

	// 直接使用传入的 writer 创建 WriteSyncer
	writeSyncer := zapcore.AddSync(writer)

	zapL, _, err := newLoggerInternal(opts, writeSyncer) // Use newLoggerInternal
	if err != nil {
		// 这种情况理论上不应该发生，因为我们控制了 writer 且 newLoggerInternal 内部处理了其他选项错误
		// 但如果 newLoggerInternal 的其他部分失败了
		panic(lmccerrors.WithCode(
			lmccerrors.Wrap(err, "failed to create logger with writer"),
			lmccerrors.ErrLogInitialization,
		))
	}

	return &logger{
		zapLogger: zapL,
		opts:      opts,
	}
}

// GetGlobalLogger 返回全局的 Logger 实例。Std() 的别名。
// (GetGlobalLogger returns the global Logger instance. Alias for Std().)
func GetGlobalLogger() Logger {
	return Std()
}

// SetGlobalLogger 设置全局的 Logger 实例。
// 注意：这会直接替换全局 logger，主要用于测试。
// (SetGlobalLogger sets the global Logger instance.)
// (Note: This directly replaces the global logger, primarily for testing.)
func SetGlobalLogger(l Logger) {
	if l == nil {
		// 不允许设置 nil logger，重新初始化为默认值
		// (Cannot set nil logger, reinitialize to default)
		Init(NewOptions())
		return
	}
	internalLog, ok := l.(*logger)
	if !ok {
		// 如果传入的不是 *logger 类型，这是一个 API 使用错误，应该 panic。
		// (If the passed type is not *logger, it's an API usage error and should panic.)
		panic(fmt.Sprintf("SetGlobalLogger: incompatible logger type %T, expected *logger created by this package", l))
	}
	std.Store(internalLog)
}

// getWriteSyncer 根据提供的选项确定并返回一个 zapcore.WriteSyncer。
//它可以配置为写入标准输出、标准错误或一个或多个文件。
// (getWriteSyncer determines and returns a zapcore.WriteSyncer based on the provided options.)
// (It can be configured to write to stdout, stderr, or one or more files.)
func getWriteSyncer(opts *Options) (zapcore.WriteSyncer, error) {
	if len(opts.OutputPaths) == 0 {
		// 如果没有指定输出路径，则默认为 stdout (Default to stdout if no output paths are specified)
		// 但通常 NewOptions 会设置默认值，所以这里更多是防御性编程
		// (But usually NewOptions sets defaults, so this is more defensive)
		return zapcore.AddSync(os.Stdout), nil
	}
	return getWriteSyncerForPaths(opts.OutputPaths, opts)
}

// getWriteSyncerForPaths 为给定的路径列表创建一个 zapcore.WriteSyncer。
// 支持 "stdout", "stderr" 以及文件路径。
// (getWriteSyncerForPaths creates a zapcore.WriteSyncer for the given list of paths.)
// (Supports "stdout", "stderr", and file paths.)
func getWriteSyncerForPaths(paths []string, opts *Options) (zapcore.WriteSyncer, error) {
	var writers []zapcore.WriteSyncer
	for _, path := range paths {
		var ws zapcore.WriteSyncer
		// var err error // err is declared within the loop for file opening specifically
		switch strings.ToLower(path) {
		case "stdout":
			ws = zapcore.AddSync(os.Stdout)
		case "stderr":
			ws = zapcore.AddSync(os.Stderr)
		default:
			// 文件路径处理，包括轮转
			// (File path handling, including rotation)
			// Also handle cases like "http://", "tcp://", etc. as invalid file paths
			if strings.Contains(path, "://") {
				return nil, lmccerrors.NewWithCode(lmccerrors.ErrLogOptionInvalid, "unsupported output path scheme: "+path)
			}

			if opts.LogRotateMaxSize > 0 { // 使用 LogRotateMaxSize 判断是否启用轮转
				// 使用 newRotateLogger 函数，它包含了目录创建和错误处理逻辑
				// (Use newRotateLogger function which includes directory creation and error handling logic)
				var err error
				ws, err = newRotateLogger(path, opts)
				if err != nil {
					return nil, err
				}
			} else {
				// 普通文件写入 (Regular file writing)
				file, errOpen := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
				if errOpen != nil {
					return nil, lmccerrors.WithCode(
						lmccerrors.Wrapf(errOpen, "failed to open log file %s", path),
						lmccerrors.ErrLogInitialization,
					)
				}
				ws = zapcore.AddSync(file)
			}
		}
		// if err != nil { // This err check is problematic if err is not properly assigned in all paths within default
		// return nil, err
		// }
		if ws != nil {
			writers = append(writers, ws)
		}
	}

	if len(writers) == 0 {
		// This case should ideally be caught by opts.Validate() if OutputPaths becomes empty after processing
		// or if all paths are invalid but don't immediately error out above.
		// However, if paths contained only invalid schemes that returned early, writers could be empty.
		return nil, lmccerrors.NewWithCode(lmccerrors.ErrLogOptionInvalid, "no valid log output writers configured from paths")
	}
	return zapcore.NewMultiWriteSyncer(writers...), nil
}

// 以下是原有的 logger 方法和全局包装函数，保持不变
// ... (Debug, Debugf, Debugw methods for *logger) ...
// ... (Info, Infof, Infow methods for *logger) ...
// ... (Warn, Warnf, Warnw methods for *logger) ...
// ... (Error, Errorf, Errorw methods for *logger) ...
// ... (Fatal, Fatalf, Fatalw methods for *logger) ...
// ... (Ctx, Ctxf, Ctxw methods for *logger) ...
// ... (Sync, WithValues, WithName, GetZapLogger methods for *logger) ...
// ... (CtxDebugf, CtxInfof, etc. methods for *logger) ...

// ... (Global Debug, Info, Warn, Error, Fatal functions) ...
// ... (Global Ctx, Ctxf, Ctxw functions) ...
// ... (Global WithValues, WithName functions) ...
// ... (Global Sync function if it was separate, or rely on Std().Sync()) ...
// Global Sync function
func Sync() error {
	return Std().Sync()
}

// --- Global Logging Functions ---
// (全局日志记录函数)

// Debug 在全局 logger 上调用 Debug。
// (Debug calls Debug on the global logger.)
func Debug(args ...any) {
	Std().Debug(args...)
}

// Debugf 在全局 logger 上调用 Debugf。
// (Debugf calls Debugf on the global logger.)
func Debugf(template string, args ...any) {
	Std().Debugf(template, args...)
}

// Debugw 在全局 logger 上调用 Debugw。
// (Debugw calls Debugw on the global logger.)
func Debugw(msg string, keysAndValues ...any) {
	Std().Debugw(msg, keysAndValues...)
}

// Info 在全局 logger 上调用 Info。
// (Info calls Info on the global logger.)
func Info(args ...any) {
	Std().Info(args...)
}

// Infof 在全局 logger 上调用 Infof。
// (Infof calls Infof on the global logger.)
func Infof(template string, args ...any) {
	Std().Infof(template, args...)
}

// Infow 在全局 logger 上调用 Infow。
// (Infow calls Infow on the global logger.)
func Infow(msg string, keysAndValues ...any) {
	Std().Infow(msg, keysAndValues...)
}

// Warn 在全局 logger 上调用 Warn。
// (Warn calls Warn on the global logger.)
func Warn(args ...any) {
	Std().Warn(args...)
}

// Warnf 在全局 logger 上调用 Warnf。
// (Warnf calls Warnf on the global logger.)
func Warnf(template string, args ...any) {
	Std().Warnf(template, args...)
}

// Warnw 在全局 logger 上调用 Warnw。
// (Warnw calls Warnw on the global logger.)
func Warnw(msg string, keysAndValues ...any) {
	Std().Warnw(msg, keysAndValues...)
}

// Error 在全局 logger 上调用 Error。
// (Error calls Error on the global logger.)
func Error(args ...any) {
	Std().Error(args...)
}

// Errorf 在全局 logger 上调用 Errorf。
// (Errorf calls Errorf on the global logger.)
func Errorf(template string, args ...any) {
	Std().Errorf(template, args...)
}

// Errorw 在全局 logger 上调用 Errorw。
// (Errorw calls Errorw on the global logger.)
func Errorw(msg string, keysAndValues ...any) {
	Std().Errorw(msg, keysAndValues...)
}

// Fatal 在全局 logger 上调用 Fatal。
// (Fatal calls Fatal on the global logger.)
func Fatal(args ...any) {
	Std().Fatal(args...)
}

// Fatalf 在全局 logger 上调用 Fatalf。
// (Fatalf calls Fatalf on the global logger.)
func Fatalf(template string, args ...any) {
	Std().Fatalf(template, args...)
}

// Fatalw 在全局 logger 上调用 Fatalw。
// (Fatalw calls Fatalw on the global logger.)
func Fatalw(msg string, keysAndValues ...any) {
	Std().Fatalw(msg, keysAndValues...)
}

// Ctx 在全局 logger 上调用 Ctx。
// (Ctx calls Ctx on the global logger.)
func Ctx(ctx context.Context, args ...any) {
	Std().Ctx(ctx, args...)
}

// Ctxf 在全局 logger 上调用 Ctxf。
// (Ctxf calls Ctxf on the global logger.)
func Ctxf(ctx context.Context, template string, args ...any) {
	Std().Ctxf(ctx, template, args...)
}

// Ctxw 在全局 logger 上调用 Ctxw。
// (Ctxw calls Ctxw on the global logger.)
func Ctxw(ctx context.Context, msg string, keysAndValues ...any) {
	Std().Ctxw(ctx, msg, keysAndValues...)
}

// WithValues 在全局 logger 上调用 WithValues 并返回一个新的 Logger 实例。
// (WithValues calls WithValues on the global logger and returns a new Logger instance.)
func WithValues(keysAndValues ...any) Logger {
	return Std().WithValues(keysAndValues...)
}

// WithName 在全局 logger 上调用 WithName 并返回一个新的 Logger 实例。
// (WithName calls WithName on the global logger and returns a new Logger instance.)
func WithName(name string) Logger {
	return Std().WithName(name)
}

// --- 已有的 logger 方法实现 (示例，确保它们都存在) ---
func (l *logger) Debug(args ...any) { l.zapLogger.Sugar().Debug(args...) }
func (l *logger) Debugf(template string, args ...any) { l.zapLogger.Sugar().Debugf(template, args...) }
func (l *logger) Debugw(msg string, keysAndValues ...any) { l.zapLogger.Sugar().Debugw(msg, keysAndValues...) }

func (l *logger) Info(args ...any) { l.zapLogger.Sugar().Info(args...) }
func (l *logger) Infof(template string, args ...any) { l.zapLogger.Sugar().Infof(template, args...) }
func (l *logger) Infow(msg string, keysAndValues ...any) {
	if l.opts.Format == FormatKeyValue {
		if kvStr := formatKeyValuePairs(keysAndValues...); kvStr != "" {
			msg = msg + " " + kvStr
		}
		l.zapLogger.Sugar().Info(msg)
	} else {
		l.zapLogger.Sugar().Infow(msg, keysAndValues...)
	}
}

func (l *logger) Warn(args ...any) { l.zapLogger.Sugar().Warn(args...) }
func (l *logger) Warnf(template string, args ...any) { l.zapLogger.Sugar().Warnf(template, args...) }
func (l *logger) Warnw(msg string, keysAndValues ...any) {
	if l.opts.Format == FormatKeyValue {
		if kvStr := formatKeyValuePairs(keysAndValues...); kvStr != "" {
			msg = msg + " " + kvStr
		}
		l.zapLogger.Sugar().Warn(msg)
	} else {
		l.zapLogger.Sugar().Warnw(msg, keysAndValues...)
	}
}

func (l *logger) Error(args ...any) { l.zapLogger.Sugar().Error(args...) }
func (l *logger) Errorf(template string, args ...any) { l.zapLogger.Sugar().Errorf(template, args...) }
func (l *logger) Errorw(msg string, keysAndValues ...any) { l.zapLogger.Sugar().Errorw(msg, keysAndValues...) }

func (l *logger) Fatal(args ...any) { l.zapLogger.Sugar().Fatal(args...) }
func (l *logger) Fatalf(template string, args ...any) { l.zapLogger.Sugar().Fatalf(template, args...) }
func (l *logger) Fatalw(msg string, keysAndValues ...any) { l.zapLogger.Sugar().Fatalw(msg, keysAndValues...) }

func (l *logger) Ctx(ctx context.Context, args ...any) {
	fields := extractContextFields(ctx, l.opts.ContextKeys)
	l.zapLogger.With(fields...).Sugar().Info(args...)
}

func (l *logger) Ctxf(ctx context.Context, template string, args ...any) {
	fields := extractContextFields(ctx, l.opts.ContextKeys)
	l.zapLogger.With(fields...).Sugar().Infof(template, args...)
}

func (l *logger) Ctxw(ctx context.Context, msg string, keysAndValues ...any) {
	fields := extractContextFields(ctx, l.opts.ContextKeys)
	
	if l.opts.Format == FormatKeyValue {
		// 对于 key=value 格式，将字段格式化为字符串并附加到消息中
		// (For key=value format, format fields as string and append to message)
		var allParts []string
		
		// 添加上下文字段
		// (Add context fields)
		if contextStr := formatFieldsAsKeyValue(fields); contextStr != "" {
			allParts = append(allParts, contextStr)
		}
		
		// 添加额外的键值对
		// (Add additional key-value pairs)
		if kvStr := formatKeyValuePairs(keysAndValues...); kvStr != "" {
			allParts = append(allParts, kvStr)
		}
		
		// 如果有字段，将它们附加到消息中
		// (If there are fields, append them to the message)
		if len(allParts) > 0 {
			msg = msg + " " + strings.Join(allParts, " ")
		}
		
		l.zapLogger.Sugar().Info(msg)
	} else {
		// 对于其他格式，使用原有逻辑
		// (For other formats, use original logic)
		allFields := append(keysAndValues, fieldsToZapAny(fields)...)
		l.zapLogger.Sugar().Infow(msg, allFields...)
	}
}

func (l *logger) Sync() error { return l.zapLogger.Sync() }

func (l *logger) WithValues(keysAndValues ...any) Logger {
	if l.opts.Format == FormatKeyValue {
		// 对于 key=value 格式，我们需要特殊处理
		// 创建一个包装器，在日志时添加这些字段
		// (For key=value format, we need special handling)
		// (Create a wrapper that adds these fields when logging)
		return &keyValueLogger{
			baseLogger: l,
			fields:     keysAndValues,
		}
	} else {
		// 对于其他格式，使用原有逻辑
		// (For other formats, use original logic)
		return &logger{
			zapLogger: l.zapLogger.With(zapFields(keysAndValues...)...), // Ensure zapFields handles pairs correctly
			opts:      l.opts, // Options are typically immutable after logger creation or carried over
		}
	}
}

func (l *logger) WithName(name string) Logger {
	return &logger{
		zapLogger: l.zapLogger.Named(name),
		opts:      l.opts,
	}
}
func (l *logger) GetZapLogger() *zap.Logger {
	return l.zapLogger
}

// --- Contextual logging methods for *logger ---
func (l *logger) CtxDebugf(ctx context.Context, template string, args ...interface{}) {
		fields := extractContextFields(ctx, l.opts.ContextKeys)
		l.zapLogger.With(fields...).Sugar().Debugf(template, args...)
	}
func (l *logger) CtxInfof(ctx context.Context, template string, args ...interface{}) {
		fields := extractContextFields(ctx, l.opts.ContextKeys)
		l.zapLogger.With(fields...).Sugar().Infof(template, args...)
	}
func (l *logger) CtxWarnf(ctx context.Context, template string, args ...interface{}) {
		fields := extractContextFields(ctx, l.opts.ContextKeys)
		l.zapLogger.With(fields...).Sugar().Warnf(template, args...)
	}
func (l *logger) CtxErrorf(ctx context.Context, template string, args ...interface{}) {
		fields := extractContextFields(ctx, l.opts.ContextKeys)
		l.zapLogger.With(fields...).Sugar().Errorf(template, args...)
	}
func (l *logger) CtxPanicf(ctx context.Context, template string, args ...interface{}) {
	fields := extractContextFields(ctx, l.opts.ContextKeys)
	l.zapLogger.With(fields...).Sugar().Panicf(template, args...)
}
func (l *logger) CtxFatalf(ctx context.Context, template string, args ...interface{}) {
	fields := extractContextFields(ctx, l.opts.ContextKeys)
	l.zapLogger.With(fields...).Sugar().Fatalf(template, args...)
}

// Helper function to convert variadic key-value pairs to zap.Field array
// Ensure it handles non-string keys or odd number of arguments gracefully if necessary,
// though zap.Any typically handles pairs.
func zapFields(keysAndValues ...any) []zap.Field {
	if len(keysAndValues)%2 != 0 {
		// Log an internal error or handle mismatched pairs, e.g., by ignoring the last odd key
		// For simplicity, zap.Any might handle this, or we can enforce pair logging.
		// This could be a source of panic if not handled well by zap.Any or if keys are not strings.
		// Zap's SugaredLogger's ...w methods are more robust here.
		// For direct zap.Logger.With, fields must be constructed carefully.
		// Let's assume keysAndValues are proper Field pairs or zap.Any can handle them.
		// A more robust way for WithValues:
		// if sugar := l.zapLogger.Sugar(); sugar != nil {
		//     return sugar.With(keysAndValues...).Desugar().Core(). ... no, this is not right.
		// }
		// For now, we'll rely on zap.Any correctly interpreting these.
		// A better approach for WithValues might be to take []zap.Field directly or process carefully.
		var fields []zap.Field
		for i := 0; i < len(keysAndValues); i += 2 {
			if i+1 < len(keysAndValues) {
				if key, ok := keysAndValues[i].(string); ok {
					fields = append(fields, zap.Any(key, keysAndValues[i+1]))
				}
				// else: log an error, key is not a string
			}
			// else: log an error, odd number of arguments
		}
		return fields
	}
	// Simplified if assuming zap.Any can handle it directly (might not be true for zap.Logger.With)
	// Correct for zap.Logger.With is to pass []zap.Field.
	// zap.SugaredLogger.With correctly handles ...any.
	// Since our Logger interface resembles SugaredLogger, maybe logger.zapLogger should be *zap.SugaredLogger?
	// Or we properly construct []zap.Field here.
	var fields []zap.Field
	for i := 0; i < len(keysAndValues); i += 2 {
		if key, ok := keysAndValues[i].(string); ok {
			fields = append(fields, zap.Any(key, keysAndValues[i+1]))
		} else {
			// Handle non-string key, e.g., log an error or skip
			fmt.Fprintf(os.Stderr, "Warning: Non-string key provided to WithValues: %v (type %T)\\n", keysAndValues[i], keysAndValues[i])
		}
	}
	return fields
}

// fieldsToZapAny converts a slice of zap.Field to a slice of any for Infow/Errorw etc.
// This is needed if Ctxw appends zap.Field to a ...any slice.
func fieldsToZapAny(fields []zap.Field) []any {
	result := make([]any, 0, len(fields)*2)
	for _, f := range fields {
		// For Infow, Errorw, etc., which take ...any, we need key-value pairs.
		// Extract the actual value from the zap.Field
		// (从 zap.Field 中提取实际值)
		var value any
		if f.Interface != nil {
			value = f.Interface
		} else {
			// Fallback to type-specific fields
			// (回退到特定类型的字段)
			switch f.Type {
			case zapcore.StringType:
				value = f.String
			case zapcore.BoolType:
				value = f.Integer != 0
			default:
				value = f.Integer
			}
		}
		result = append(result, f.Key, value)
	}
	return result
}

// formatFieldsAsKeyValue 将字段格式化为 key=value 格式的字符串
// (formatFieldsAsKeyValue formats fields as key=value format string)
func formatFieldsAsKeyValue(fields []zap.Field) string {
	if len(fields) == 0 {
		return ""
	}
	
	var parts []string
	for _, f := range fields {
		var value string
		if f.Interface != nil {
			value = fmt.Sprintf("%v", f.Interface)
		} else {
			switch f.Type {
			case zapcore.StringType:
				value = f.String
			case zapcore.BoolType:
				if f.Integer != 0 {
					value = "true"
				} else {
					value = "false"
				}
			default:
				value = fmt.Sprintf("%v", f.Integer)
			}
		}
		// 如果值包含空格，用引号包围
		// (If value contains spaces, surround with quotes)
		if strings.Contains(value, " ") {
			value = fmt.Sprintf(`"%s"`, value)
		}
		parts = append(parts, fmt.Sprintf("%s=%s", f.Key, value))
	}
	return strings.Join(parts, " ")
}

// formatKeyValuePairs 将键值对格式化为 key=value 格式的字符串
// (formatKeyValuePairs formats key-value pairs as key=value format string)
func formatKeyValuePairs(keysAndValues ...any) string {
	if len(keysAndValues) == 0 {
		return ""
	}
	
	var parts []string
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 >= len(keysAndValues) {
			break // 跳过奇数个参数的最后一个
		}
		
		key := fmt.Sprintf("%v", keysAndValues[i])
		value := fmt.Sprintf("%v", keysAndValues[i+1])
		
		// 如果值包含空格，用引号包围
		// (If value contains spaces, surround with quotes)
		if strings.Contains(value, " ") {
			value = fmt.Sprintf(`"%s"`, value)
		}
		parts = append(parts, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(parts, " ")
}

// extractContextFields extracts configured keys from context and returns them as zap.Fields
func extractContextFields(ctx context.Context, contextKeys []any) []zap.Field { // Changed to []any
	if ctx == nil || len(contextKeys) == 0 {
		return nil
	}
	var fields []zap.Field
	for _, keyAny := range contextKeys {
		if keyAny == nil {
			continue
		} // Skip nil keys in the list

		value := ctx.Value(keyAny)
		if value != nil {
			var keyStr string
			switch typedKey := keyAny.(type) {
			case string:
				keyStr = typedKey
			case contextKey:
				switch typedKey {
				case TraceIDKey:
					keyStr = "trace_id"
				case RequestIDKey:
					keyStr = "request_id"
				default:
					keyStr = fmt.Sprintf("%v", typedKey) // Fallback for other contextKey values
				}
			case fmt.Stringer:
				keyStr = typedKey.String()
			default:
				// For custom string-based types (like type MyKey string), convert to string
				// (对于基于字符串的自定义类型（如 type MyKey string），转换为字符串)
				keyStr = fmt.Sprintf("%s", keyAny)
			}
			fields = append(fields, zap.Any(keyStr, value))
		}
	}
	return fields
}

// Make sure to import "gopkg.in/natefinch/lumberjack.v2" if RotationConfig is used.
// For example:
// import (
//     // ...
//     "gopkg.in/natefinch/lumberjack.v2"
// )
// This should be at the top of the file.

// Ensure RotationConfig fields are correctly mapped from options to lumberjack.Logger
// E.g., opts.RotationConfig.Enable, opts.RotationConfig.MaxSizeMB, etc.
// The current getWriteSyncerForPaths has a placeholder for this.

// --- Global Contextual Logging Functions ---
// (全局上下文日志记录函数)

// CtxDebugf 在全局 logger 上调用 CtxDebugf。
// (CtxDebugf calls CtxDebugf on the global logger.)
func CtxDebugf(ctx context.Context, template string, args ...interface{}) {
	Std().CtxDebugf(ctx, template, args...)
}

// CtxInfof 在全局 logger 上调用 CtxInfof。
// (CtxInfof calls CtxInfof on the global logger.)
func CtxInfof(ctx context.Context, template string, args ...interface{}) {
	Std().CtxInfof(ctx, template, args...)
}

// CtxWarnf 在全局 logger 上调用 CtxWarnf。
// (CtxWarnf calls CtxWarnf on the global logger.)
func CtxWarnf(ctx context.Context, template string, args ...interface{}) {
	Std().CtxWarnf(ctx, template, args...)
}

// CtxErrorf 在全局 logger 上调用 CtxErrorf。
// (CtxErrorf calls CtxErrorf on the global logger.)
func CtxErrorf(ctx context.Context, template string, args ...interface{}) {
	Std().CtxErrorf(ctx, template, args...)
}

// CtxPanicf 在全局 logger 上调用 CtxPanicf。
// (CtxPanicf calls CtxPanicf on the global logger.)
func CtxPanicf(ctx context.Context, template string, args ...interface{}) {
	Std().CtxPanicf(ctx, template, args...)
}

// CtxFatalf 在全局 logger 上调用 CtxFatalf。
// (CtxFatalf calls CtxFatalf on the global logger.)
func CtxFatalf(ctx context.Context, template string, args ...interface{}) {
	Std().CtxFatalf(ctx, template, args...)
}

// --- keyValueLogger 方法实现 ---
// (keyValueLogger method implementations)

func (kvl *keyValueLogger) Debug(args ...any) {
	msg := fmt.Sprint(args...)
	if kvStr := formatKeyValuePairs(kvl.fields...); kvStr != "" {
		msg = msg + " " + kvStr
	}
	kvl.baseLogger.Debug(msg)
}

func (kvl *keyValueLogger) Debugf(template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	if kvStr := formatKeyValuePairs(kvl.fields...); kvStr != "" {
		msg = msg + " " + kvStr
	}
	kvl.baseLogger.Debug(msg)
}

func (kvl *keyValueLogger) Debugw(msg string, keysAndValues ...any) {
	allFields := append(kvl.fields, keysAndValues...)
	if kvStr := formatKeyValuePairs(allFields...); kvStr != "" {
		msg = msg + " " + kvStr
	}
	kvl.baseLogger.Debug(msg)
}

func (kvl *keyValueLogger) Info(args ...any) {
	msg := fmt.Sprint(args...)
	if kvStr := formatKeyValuePairs(kvl.fields...); kvStr != "" {
		msg = msg + " " + kvStr
	}
	kvl.baseLogger.Info(msg)
}

func (kvl *keyValueLogger) Infof(template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	if kvStr := formatKeyValuePairs(kvl.fields...); kvStr != "" {
		msg = msg + " " + kvStr
	}
	kvl.baseLogger.Info(msg)
}

func (kvl *keyValueLogger) Infow(msg string, keysAndValues ...any) {
	allFields := append(kvl.fields, keysAndValues...)
	if kvStr := formatKeyValuePairs(allFields...); kvStr != "" {
		msg = msg + " " + kvStr
	}
	kvl.baseLogger.Info(msg)
}

func (kvl *keyValueLogger) Warn(args ...any) {
	msg := fmt.Sprint(args...)
	if kvStr := formatKeyValuePairs(kvl.fields...); kvStr != "" {
		msg = msg + " " + kvStr
	}
	kvl.baseLogger.Warn(msg)
}

func (kvl *keyValueLogger) Warnf(template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	if kvStr := formatKeyValuePairs(kvl.fields...); kvStr != "" {
		msg = msg + " " + kvStr
	}
	kvl.baseLogger.Warn(msg)
}

func (kvl *keyValueLogger) Warnw(msg string, keysAndValues ...any) {
	allFields := append(kvl.fields, keysAndValues...)
	if kvStr := formatKeyValuePairs(allFields...); kvStr != "" {
		msg = msg + " " + kvStr
	}
	kvl.baseLogger.Warn(msg)
}

func (kvl *keyValueLogger) Error(args ...any) {
	msg := fmt.Sprint(args...)
	if kvStr := formatKeyValuePairs(kvl.fields...); kvStr != "" {
		msg = msg + " " + kvStr
	}
	kvl.baseLogger.Error(msg)
}

func (kvl *keyValueLogger) Errorf(template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	if kvStr := formatKeyValuePairs(kvl.fields...); kvStr != "" {
		msg = msg + " " + kvStr
	}
	kvl.baseLogger.Error(msg)
}

func (kvl *keyValueLogger) Errorw(msg string, keysAndValues ...any) {
	allFields := append(kvl.fields, keysAndValues...)
	if kvStr := formatKeyValuePairs(allFields...); kvStr != "" {
		msg = msg + " " + kvStr
	}
	kvl.baseLogger.Error(msg)
}

func (kvl *keyValueLogger) Fatal(args ...any) {
	msg := fmt.Sprint(args...)
	if kvStr := formatKeyValuePairs(kvl.fields...); kvStr != "" {
		msg = msg + " " + kvStr
	}
	kvl.baseLogger.Fatal(msg)
}

func (kvl *keyValueLogger) Fatalf(template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	if kvStr := formatKeyValuePairs(kvl.fields...); kvStr != "" {
		msg = msg + " " + kvStr
	}
	kvl.baseLogger.Fatal(msg)
}

func (kvl *keyValueLogger) Fatalw(msg string, keysAndValues ...any) {
	allFields := append(kvl.fields, keysAndValues...)
	if kvStr := formatKeyValuePairs(allFields...); kvStr != "" {
		msg = msg + " " + kvStr
	}
	kvl.baseLogger.Fatal(msg)
}

func (kvl *keyValueLogger) DPanic(args ...any) {
	msg := fmt.Sprint(args...)
	if kvStr := formatKeyValuePairs(kvl.fields...); kvStr != "" {
		msg = msg + " " + kvStr
	}
	kvl.baseLogger.DPanic(msg)
}

func (kvl *keyValueLogger) DPanicf(template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	if kvStr := formatKeyValuePairs(kvl.fields...); kvStr != "" {
		msg = msg + " " + kvStr
	}
	kvl.baseLogger.DPanic(msg)
}

func (kvl *keyValueLogger) DPanicw(msg string, keysAndValues ...any) {
	allFields := append(kvl.fields, keysAndValues...)
	if kvStr := formatKeyValuePairs(allFields...); kvStr != "" {
		msg = msg + " " + kvStr
	}
	kvl.baseLogger.DPanic(msg)
}

func (kvl *keyValueLogger) Ctx(ctx context.Context, args ...any) {
	msg := fmt.Sprint(args...)
	fields := extractContextFields(ctx, kvl.baseLogger.opts.ContextKeys)
	
	var allParts []string
	if contextStr := formatFieldsAsKeyValue(fields); contextStr != "" {
		allParts = append(allParts, contextStr)
	}
	if kvStr := formatKeyValuePairs(kvl.fields...); kvStr != "" {
		allParts = append(allParts, kvStr)
	}
	
	if len(allParts) > 0 {
		msg = msg + " " + strings.Join(allParts, " ")
	}
	kvl.baseLogger.Info(msg)
}

func (kvl *keyValueLogger) Ctxf(ctx context.Context, template string, args ...any) {
	msg := fmt.Sprintf(template, args...)
	fields := extractContextFields(ctx, kvl.baseLogger.opts.ContextKeys)
	
	var allParts []string
	if contextStr := formatFieldsAsKeyValue(fields); contextStr != "" {
		allParts = append(allParts, contextStr)
	}
	if kvStr := formatKeyValuePairs(kvl.fields...); kvStr != "" {
		allParts = append(allParts, kvStr)
	}
	
	if len(allParts) > 0 {
		msg = msg + " " + strings.Join(allParts, " ")
	}
	kvl.baseLogger.Info(msg)
}

func (kvl *keyValueLogger) Ctxw(ctx context.Context, msg string, keysAndValues ...any) {
	fields := extractContextFields(ctx, kvl.baseLogger.opts.ContextKeys)
	
	var allParts []string
	if contextStr := formatFieldsAsKeyValue(fields); contextStr != "" {
		allParts = append(allParts, contextStr)
	}
	
	allFields := append(kvl.fields, keysAndValues...)
	if kvStr := formatKeyValuePairs(allFields...); kvStr != "" {
		allParts = append(allParts, kvStr)
	}
	
	if len(allParts) > 0 {
		msg = msg + " " + strings.Join(allParts, " ")
	}
	kvl.baseLogger.Info(msg)
}

func (kvl *keyValueLogger) Sync() error {
	return kvl.baseLogger.Sync()
}

func (kvl *keyValueLogger) WithValues(keysAndValues ...any) Logger {
	// 合并现有字段和新字段
	// (Merge existing fields with new fields)
	allFields := append(kvl.fields, keysAndValues...)
	return &keyValueLogger{
		baseLogger: kvl.baseLogger,
		fields:     allFields,
	}
}

func (kvl *keyValueLogger) WithName(name string) Logger {
	return &keyValueLogger{
		baseLogger: &logger{
			zapLogger: kvl.baseLogger.zapLogger.Named(name),
			opts:      kvl.baseLogger.opts,
		},
		fields: kvl.fields,
	}
}

func (kvl *keyValueLogger) GetZapLogger() *zap.Logger {
	return kvl.baseLogger.GetZapLogger()
}

func (kvl *keyValueLogger) CtxDebugf(ctx context.Context, template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	fields := extractContextFields(ctx, kvl.baseLogger.opts.ContextKeys)
	
	var allParts []string
	if contextStr := formatFieldsAsKeyValue(fields); contextStr != "" {
		allParts = append(allParts, contextStr)
	}
	if kvStr := formatKeyValuePairs(kvl.fields...); kvStr != "" {
		allParts = append(allParts, kvStr)
	}
	
	if len(allParts) > 0 {
		msg = msg + " " + strings.Join(allParts, " ")
	}
	kvl.baseLogger.Debug(msg)
}

func (kvl *keyValueLogger) CtxInfof(ctx context.Context, template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	fields := extractContextFields(ctx, kvl.baseLogger.opts.ContextKeys)
	
	var allParts []string
	if contextStr := formatFieldsAsKeyValue(fields); contextStr != "" {
		allParts = append(allParts, contextStr)
	}
	if kvStr := formatKeyValuePairs(kvl.fields...); kvStr != "" {
		allParts = append(allParts, kvStr)
	}
	
	if len(allParts) > 0 {
		msg = msg + " " + strings.Join(allParts, " ")
	}
	kvl.baseLogger.Info(msg)
}

func (kvl *keyValueLogger) CtxWarnf(ctx context.Context, template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	fields := extractContextFields(ctx, kvl.baseLogger.opts.ContextKeys)
	
	var allParts []string
	if contextStr := formatFieldsAsKeyValue(fields); contextStr != "" {
		allParts = append(allParts, contextStr)
	}
	if kvStr := formatKeyValuePairs(kvl.fields...); kvStr != "" {
		allParts = append(allParts, kvStr)
	}
	
	if len(allParts) > 0 {
		msg = msg + " " + strings.Join(allParts, " ")
	}
	kvl.baseLogger.Warn(msg)
}

func (kvl *keyValueLogger) CtxErrorf(ctx context.Context, template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	fields := extractContextFields(ctx, kvl.baseLogger.opts.ContextKeys)
	
	var allParts []string
	if contextStr := formatFieldsAsKeyValue(fields); contextStr != "" {
		allParts = append(allParts, contextStr)
	}
	if kvStr := formatKeyValuePairs(kvl.fields...); kvStr != "" {
		allParts = append(allParts, kvStr)
	}
	
	if len(allParts) > 0 {
		msg = msg + " " + strings.Join(allParts, " ")
	}
	kvl.baseLogger.Error(msg)
}

func (kvl *keyValueLogger) CtxPanicf(ctx context.Context, template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	fields := extractContextFields(ctx, kvl.baseLogger.opts.ContextKeys)
	
	var allParts []string
	if contextStr := formatFieldsAsKeyValue(fields); contextStr != "" {
		allParts = append(allParts, contextStr)
	}
	if kvStr := formatKeyValuePairs(kvl.fields...); kvStr != "" {
		allParts = append(allParts, kvStr)
	}
	
	if len(allParts) > 0 {
		msg = msg + " " + strings.Join(allParts, " ")
	}
	kvl.baseLogger.zapLogger.Sugar().Panicf(msg)
}

func (kvl *keyValueLogger) CtxFatalf(ctx context.Context, template string, args ...interface{}) {
	msg := fmt.Sprintf(template, args...)
	fields := extractContextFields(ctx, kvl.baseLogger.opts.ContextKeys)
	
	var allParts []string
	if contextStr := formatFieldsAsKeyValue(fields); contextStr != "" {
		allParts = append(allParts, contextStr)
	}
	if kvStr := formatKeyValuePairs(kvl.fields...); kvStr != "" {
		allParts = append(allParts, kvStr)
	}
	
	if len(allParts) > 0 {
		msg = msg + " " + strings.Join(allParts, " ")
	}
	kvl.baseLogger.Fatal(msg)
}



// getEncoder 根据 Options 配置返回一个 zapcore.Encoder。
// (getEncoder returns a zapcore.Encoder based on the Options configuration.)
// ... existing code ...

