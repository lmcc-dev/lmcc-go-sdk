/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package log

import (
	"context"
	"fmt"
	"os" // Needed for os.Stdout, os.Stderr
	"sync"
	"sync/atomic" // Added for atomic.Pointer

	merrors "github.com/marmotedu/errors" // 导入 marmotedu/errors 包 (Import marmotedu/errors package)
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

// 确保 logger 实现了 Logger 接口。
var _ Logger = (*logger)(nil)

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
// (Init initializes the global logger std with the given options.)
// (This function is thread-safe. If std is already initialized, it will be atomically overwritten with the new configuration.)
func Init(opts *Options) {
	// newLogger 是内部构造函数，确保它处理 opts 为 nil 的情况或返回错误
	// (newLogger is the internal constructor, ensure it handles nil opts or returns an error if appropriate)
	// 假设 newLogger 总是返回一个有效的 logger 或 panic (如果 opts 无效)
	// (Assuming newLogger always returns a valid logger or panics if opts are invalid)
	l := newLogger(opts) 
	std.Store(l)
}

// NewLogger 根据提供的选项创建一个新的 Logger 实例。
// (NewLogger creates a new Logger instance based on the provided options.)
func NewLogger(opts *Options) Logger {
	return newLogger(opts) // Calls the internal constructor
}

// Std 返回全局的 Logger 实例。
// 如果 std 未初始化，它将使用默认选项进行初始化。
// 此函数通过双重检查锁定和原子操作确保线程安全。
// (Std returns the global Logger instance.)
// (If std is not initialized, it will be initialized with default options.)
// (This function ensures thread-safety through double-checked locking and atomic operations.)
func Std() Logger {
	l := std.Load()
	if l == nil {
		mu.Lock()
		defer mu.Unlock()
		// Double-check in case another goroutine initialized it while we were waiting for the lock
		// (再次检查，以防在我们等待锁的时候，另一个 goroutine 已经初始化了它)
		l = std.Load()
		if l == nil {
			l = newLogger(NewOptions()) // Use default options if not initialized
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
		// 使用 merrors.New 创建错误，以包含堆栈跟踪。
		// (Use merrors.New to create an error with a stack trace.)
		return merrors.New("cannot reconfigure global logger with nil options")
	}
	// 假设 newLogger 会对 newOpts 进行校验，如果选项无效可能会 panic 或返回可区分的错误
	// (Assuming newLogger validates newOpts and might panic or return a distinguishable error for invalid options)
	// 为了简单起见，这里我们只创建新的 logger。如果 newLogger 可以返回错误，则应处理该错误。
	// (For simplicity, we just create the new logger here. If newLogger could return an error, it should be handled.)
	newL := newLogger(newOpts) // Use the internal constructor
	std.Store(newL)
	return nil // 假设 newLogger 成功或 panic
}

// newLogger 是 logger 的内部构造函数。
// (newLogger is the internal constructor for a logger.)
// 它接收 Options 并返回一个配置好的 *logger 实例。
// (It takes Options and returns a configured *logger instance.)
// 注意：这个函数现在返回 *logger 而不是 Logger 接口，以配合 atomic.Pointer[logger]。
// (Note: This function now returns *logger instead of the Logger interface to work with atomic.Pointer[logger].)
func newLogger(opts *Options) *logger { // Changed return type to *logger
	if opts == nil {
		opts = NewOptions() // 使用默认选项，如果提供的是 nil (Use default options if nil is provided)
	}

	// 将日志级别字符串转换为 zapcore.Level
	// (Convert log level string to zapcore.Level)
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		zapLevel = zapcore.InfoLevel // Default level
	}

	// 配置编码器 (Configure encoder)
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "timestamp",
		NameKey:        "logger", // NameKey is "logger", will be used if logger has a name
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder, // Default for console
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder, // Ensures the full name is encoded if present
	}

	var encoder zapcore.Encoder
	if opts.Format == FormatText {
		if opts.EnableColor {
			encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Colored level for text
		}
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else if opts.Format == FormatJSON {
		encoderConfig.LevelKey = "L" // Use "L" for level in JSON as per common practice
		encoderConfig.TimeKey = "T"
		encoderConfig.MessageKey = "M"
		encoderConfig.CallerKey = "C"
		encoderConfig.NameKey = "N" // Use "N" for name in JSON
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		// Fallback to JSON encoder if format is unknown
		// (如果格式未知，则回退到 JSON 编码器)
		encoderConfig.LevelKey = "L"
		encoderConfig.TimeKey = "T"
		encoderConfig.MessageKey = "M"
		encoderConfig.CallerKey = "C"
		encoderConfig.NameKey = "N"
		encoder = zapcore.NewJSONEncoder(encoderConfig)
		// Optionally log a warning that an unknown format was specified
		// (可以选择性地记录一个警告，说明指定了未知格式)
		// fmt.Fprintf(os.Stderr, "Unknown log format \\"%s\\", defaulting to JSON\\n", opts.Format)
	}

	// 获取写入器 (Get writer syncer)
	writeSyncer, err := getWriteSyncer(opts)
	if err != nil {
		// 如果获取写入器失败，打印错误到 stderr 并 panic
		// (If getting writer syncer fails, print error to stderr and panic)
		// This is a critical failure, as logging cannot proceed.
		// (这是一个严重错误，因为日志记录无法继续。)
		fmt.Fprintf(os.Stderr, "Failed to get write syncer: %v\\n", err)
		panic(fmt.Sprintf("Failed to get write syncer: %v", err))
	}

	core := zapcore.NewCore(encoder, writeSyncer, zapLevel)

	zapOptions := []zap.Option{
		zap.AddCaller(), // 总是添加调用者信息 (Always add caller info)
		// 根据 opts.DisableStacktrace 来决定是否添加 zap 自己的堆栈跟踪
		// (Decide whether to add zap's own stack trace based on opts.DisableStacktrace)
		// 注意：这与 marmotedu/errors 的 errorVerbose 堆栈是分开的
		// (Note: This is separate from marmotedu/errors' errorVerbose stack)
		// 如果 opts.DisableStacktrace 为 true，我们不希望 zap 自动添加堆栈
		// (If opts.DisableStacktrace is true, we don't want zap to add stack traces automatically)
	}
	if !opts.DisableStacktrace {
		zapOptions = append(zapOptions, zap.AddStacktrace(zapcore.ErrorLevel)) // ErrorLevel 及以上级别日志会自动附加堆栈跟踪 (ErrorLevel and above will automatically attach stack traces)
	}

	// 如果 Development 模式开启，配置特定的开发日志选项
	// (If Development mode is on, configure specific development logging options)
	if opts.Development {
		// Development 该选项会更改 DPanic 日志的行为，使其在开发模式下 panic
		// (The Development option changes the behavior of DPanic logs to panic in development mode)
		zapOptions = append(zapOptions, zap.Development())
	}

	// 为了让 "caller" 字段显示应用代码的调用位置，而不是 pkg/log 内部的位置
	// For the "caller" field to show the call site in application code, not within pkg/log
	// 需要跳过 sdklog 本身的封装层。通常是2层 (e.g., Debug -> sugar.Debug -> ...)
	// We need to skip sdklog's own wrapper layers. Usually 2 levels (e.g., Debug -> sugar.Debug -> ...)
	// 如果你的全局函数 (如 log.Info) 直接调用 logger 实例的方法 (如 l.Info),
	// 而实例方法又直接调用 zapLogger.Sugar().Info(), 那么从应用代码到 zap 核心是2层。
	// 实际跳过的层级可能需要根据具体调用链调整。
	// If your global functions (like log.Info) directly call logger instance methods (like l.Info),
	// and instance methods directly call zapLogger.Sugar().Info(), then it's 2 levels from app code to zap core.
	// The actual number of skipped levels might need adjustment based on the specific call chain.
	zapOptions = append(zapOptions, zap.AddCallerSkip(1)) // 先尝试跳过1层，通常是 sugar 封装。如果还不够，可以调整为2。


	finalLogger := zap.New(core, zapOptions...)

	// 如果 opts.Name 非空，则为 logger 设置名称
	// (If opts.Name is not empty, set the name for the logger)
	if opts.Name != "" {
		finalLogger = finalLogger.Named(opts.Name)
	}

	return &logger{
		zapLogger: finalLogger,
		opts:      opts, // Store the applied options
	}
}

// getWriteSyncer 根据选项创建 zapcore.WriteSyncer
// (getWriteSyncer creates a zapcore.WriteSyncer based on the options)
func getWriteSyncer(opts *Options) (zapcore.WriteSyncer, error) {
	var syncers []zapcore.WriteSyncer

	if len(opts.OutputPaths) == 0 {
		// 如果 OutputPaths 为空，则默认为 stdout
		// (Default to stdout if OutputPaths is empty)
		opts.OutputPaths = []string{"stdout"}
	}

	for _, path := range opts.OutputPaths {
		var syncer zapcore.WriteSyncer
		var err error // Declare err here to catch error from newRotateLogger
		switch path {
		case "stdout":
			syncer = zapcore.Lock(os.Stdout)
		case "stderr":
			syncer = zapcore.Lock(os.Stderr)
		default:
			// 视为文件路径，调用 rotation.go 中的 newRotateLogger
			// (Treat as file path, call newRotateLogger from rotation.go)
			syncer, err = newRotateLogger(path, opts) // Call the refactored function
			if err != nil {
				// 如果创建轮转日志失败，返回错误以便 NewLogger 处理
				// (If creating rotated log fails, return error for NewLogger to handle)
				// 使用 merrors.Wrapf 包装错误，以添加堆栈跟踪和上下文。
				// (Wrap the error with merrors.Wrapf to add stack trace and context.)
				return nil, merrors.Wrapf(err, "failed to create rotating logger for path %s", path)
			}
		}
		syncers = append(syncers, syncer)
	}

	if len(syncers) == 0 {
		// 理论上不应该发生，因为前面有默认 stdout 处理，但作为防御性编程
		// (Shouldn't happen due to default stdout, but defensive programming)
		return zapcore.Lock(os.Stdout), nil // Return stdout syncer as a safe default
	}

	// 合并所有 syncers
	// (Combine all syncers)
	return zapcore.NewMultiWriteSyncer(syncers...), nil
}

// getWriteSyncerForPaths is a helper specifically for creating a combined syncer for a given list of paths.
// It's used by NewLogger to handle ErrorOutputPaths separately.
func getWriteSyncerForPaths(paths []string, opts *Options) (zapcore.WriteSyncer, error) {
	if len(paths) == 0 {
		// 使用 merrors.New 创建错误，以包含堆栈跟踪。
		// (Use merrors.New to create an error with a stack trace.)
		return nil, merrors.New("no paths provided to getWriteSyncerForPaths")
	}
	syncers := make([]zapcore.WriteSyncer, 0, len(paths))
	for _, path := range paths {
		var syncer zapcore.WriteSyncer
		var err error
		switch path {
		case "stdout":
			syncer = zapcore.Lock(os.Stdout)
		case "stderr":
			syncer = zapcore.Lock(os.Stderr)
		default:
			syncer, err = newRotateLogger(path, opts)
			if err != nil {
				// 使用 merrors.Wrapf 包装错误，以添加堆栈跟踪和上下文。
				// (Wrap the error with merrors.Wrapf to add stack trace and context.)
				return nil, merrors.Wrapf(err, "failed to create rotating logger for path %s", path)
			}
		}
		syncers = append(syncers, syncer)
	}

	if len(syncers) == 1 {
		return syncers[0], nil
	}
	return zapcore.NewMultiWriteSyncer(syncers...), nil
}

// --- Logger interface method implementations ---

// Debug logs a message at DebugLevel using fmt.Sprint style.
func (l *logger) Debug(args ...any) { l.zapLogger.Sugar().Debug(args...) }

// Debugf logs a message at DebugLevel using fmt.Sprintf style.
func (l *logger) Debugf(template string, args ...any) { l.zapLogger.Sugar().Debugf(template, args...) }

// Debugw logs a message at DebugLevel with key-value pairs.
func (l *logger) Debugw(msg string, keysAndValues ...any) {
	l.zapLogger.Sugar().Debugw(msg, keysAndValues...)
}

// Info logs a message at InfoLevel using fmt.Sprint style.
func (l *logger) Info(args ...any) { l.zapLogger.Sugar().Info(args...) }

// Infof logs a message at InfoLevel using fmt.Sprintf style.
func (l *logger) Infof(template string, args ...any) { l.zapLogger.Sugar().Infof(template, args...) }

// Infow logs a message at InfoLevel with key-value pairs.
func (l *logger) Infow(msg string, keysAndValues ...any) {
	l.zapLogger.Sugar().Infow(msg, keysAndValues...)
}

// Warn logs a message at WarnLevel using fmt.Sprint style.
func (l *logger) Warn(args ...any) { l.zapLogger.Sugar().Warn(args...) }

// Warnf logs a message at WarnLevel using fmt.Sprintf style.
func (l *logger) Warnf(template string, args ...any) { l.zapLogger.Sugar().Warnf(template, args...) }

// Warnw logs a message at WarnLevel with key-value pairs.
func (l *logger) Warnw(msg string, keysAndValues ...any) {
	l.zapLogger.Sugar().Warnw(msg, keysAndValues...)
}

// Error logs a message at ErrorLevel using fmt.Sprint style.
func (l *logger) Error(args ...any) { l.zapLogger.Sugar().Error(args...) }

// Errorf logs a message at ErrorLevel using fmt.Sprintf style.
func (l *logger) Errorf(template string, args ...any) { l.zapLogger.Sugar().Errorf(template, args...) }

// Errorw logs a message at ErrorLevel with key-value pairs.
func (l *logger) Errorw(msg string, keysAndValues ...any) {
	l.zapLogger.Sugar().Errorw(msg, keysAndValues...)
}

// Fatal logs a message at FatalLevel using fmt.Sprint style, then calls os.Exit(1).
func (l *logger) Fatal(args ...any) { l.zapLogger.Sugar().Fatal(args...) } // Note: Fatal calls os.Exit

// Fatalf logs a message at FatalLevel using fmt.Sprintf style, then calls os.Exit(1).
func (l *logger) Fatalf(template string, args ...any) { l.zapLogger.Sugar().Fatalf(template, args...) } // Note: Fatal calls os.Exit

// Fatalw logs a message at FatalLevel with key-value pairs, then calls os.Exit(1).
func (l *logger) Fatalw(msg string, keysAndValues ...any) {
	l.zapLogger.Sugar().Fatalw(msg, keysAndValues...) // Note: Fatal calls os.Exit
}

// Ctx logs a message at InfoLevel using fmt.Sprint style, extracting fields from the context.
func (l *logger) Ctx(ctx context.Context, args ...any) {
	if l.zapLogger.Level() <= zapcore.InfoLevel {
		fields := extractContextFields(ctx, l.opts.ContextKeys)
		l.zapLogger.With(fields...).Sugar().Info(args...)
	}
}

// Ctxf logs a message at InfoLevel using fmt.Sprintf style, extracting fields from the context.
func (l *logger) Ctxf(ctx context.Context, template string, args ...any) {
	if l.zapLogger.Level() <= zapcore.InfoLevel {
		fields := extractContextFields(ctx, l.opts.ContextKeys)
		l.zapLogger.With(fields...).Sugar().Infof(template, args...)
	}
}

// Ctxw logs a message at InfoLevel with key-value pairs, also extracting fields from the context.
func (l *logger) Ctxw(ctx context.Context, msg string, keysAndValues ...any) {
	if l.zapLogger.Level() <= zapcore.InfoLevel {
		fields := extractContextFields(ctx, l.opts.ContextKeys)
		l.zapLogger.With(fields...).Sugar().Infow(msg, keysAndValues...)
	}
}

// Sync flushes any buffered log entries to the underlying writers.
func (l *logger) Sync() error { return l.zapLogger.Sync() }

// WithValues adds a set of key-value pairs context to the logger.
// It returns a new Logger instance with the added context.
func (l *logger) WithValues(keysAndValues ...any) Logger {
	newZapLogger := l.zapLogger.Sugar().With(keysAndValues...).Desugar()
	return &logger{
		zapLogger: newZapLogger,
		opts:      l.opts, // Inherit options
	}
}

// WithName adds a new element to the logger's name.
// It returns a new Logger instance with the extended name.
func (l *logger) WithName(name string) Logger {
	newZapLogger := l.zapLogger.Named(name)
	return &logger{
		zapLogger: newZapLogger,
		opts:      l.opts, // Inherit options
	}
}

// GetZapLogger returns the underlying zap.Logger.
// This might be useful for advanced integration or testing.
func (l *logger) GetZapLogger() *zap.Logger {
	return l.zapLogger
}

// --- Global convenience functions ---

// Sync flushes the global logger.
func Sync() error {
	return Std().Sync()
}

// Debug logs a message at DebugLevel using the global logger.
func Debug(args ...any) {
	Std().Debug(args...)
}

// Debugf logs a message at DebugLevel using the global logger.
func Debugf(template string, args ...any) {
	Std().Debugf(template, args...)
}

// Debugw logs a message at DebugLevel with key-value pairs using the global logger.
func Debugw(msg string, keysAndValues ...any) {
	Std().Debugw(msg, keysAndValues...)
}

// Info logs a message at InfoLevel using the global logger.
func Info(args ...any) {
	Std().Info(args...)
}

// Infof logs a message at InfoLevel using the global logger.
func Infof(template string, args ...any) {
	Std().Infof(template, args...)
}

// Infow logs a message at InfoLevel with key-value pairs using the global logger.
func Infow(msg string, keysAndValues ...any) {
	Std().Infow(msg, keysAndValues...)
}

// Warn logs a message at WarnLevel using the global logger.
func Warn(args ...any) {
	Std().Warn(args...)
}

// Warnf logs a message at WarnLevel using the global logger.
func Warnf(template string, args ...any) {
	Std().Warnf(template, args...)
}

// Warnw logs a message at WarnLevel with key-value pairs using the global logger.
func Warnw(msg string, keysAndValues ...any) {
	Std().Warnw(msg, keysAndValues...)
}

// Error logs a message at ErrorLevel using the global logger.
func Error(args ...any) {
	Std().Error(args...)
}

// Errorf logs a message at ErrorLevel using the global logger.
func Errorf(template string, args ...any) {
	Std().Errorf(template, args...)
}

// Errorw logs a message at ErrorLevel with key-value pairs using the global logger.
func Errorw(msg string, keysAndValues ...any) {
	Std().Errorw(msg, keysAndValues...)
}

// Fatal logs a message at FatalLevel using the global logger, then calls os.Exit(1).
func Fatal(args ...any) {
	Std().Fatal(args...)
}

// Fatalf logs a message at FatalLevel using the global logger, then calls os.Exit(1).
func Fatalf(template string, args ...any) {
	Std().Fatalf(template, args...)
}

// Fatalw logs a message at FatalLevel with key-value pairs using the global logger, then calls os.Exit(1).
func Fatalw(msg string, keysAndValues ...any) {
	Std().Fatalw(msg, keysAndValues...)
}

// Ctx logs a message at InfoLevel using the global logger, extracting fields from the context.
func Ctx(ctx context.Context, args ...any) {
	Std().Ctx(ctx, args...)
}

// Ctxf logs a message at InfoLevel using the global logger, extracting fields from the context.
func Ctxf(ctx context.Context, template string, args ...any) {
	Std().Ctxf(ctx, template, args...)
}

// Ctxw logs a message at InfoLevel with key-value pairs using the global logger, extracting fields from the context.
func Ctxw(ctx context.Context, msg string, keysAndValues ...any) {
	Std().Ctxw(ctx, msg, keysAndValues...)
}

// WithValues adds key-value pairs context to the global logger, returning a new logger.
func WithValues(keysAndValues ...any) Logger {
	return Std().WithValues(keysAndValues...)
}

// WithName adds a new element to the global logger's name, returning a new logger.
func WithName(name string) Logger {
	return Std().WithName(name)
}

// --- Global Contextual Logging Convenience Functions ---
// (全局上下文日志便捷函数)

// CtxDebugf 使用 fmt.Sprintf 风格记录一条 Debug 级别的消息，并从 context 中提取字段，通过全局记录器。
// (CtxDebugf logs a message at DebugLevel using fmt.Sprintf style, extracting fields from the context, via the global logger.)
func CtxDebugf(ctx context.Context, template string, args ...interface{}) {
	Std().CtxDebugf(ctx, template, args...)
}

// CtxInfof 使用 fmt.Sprintf 风格记录一条 Info 级别的消息，并从 context 中提取字段，通过全局记录器。
// (CtxInfof logs a message at InfoLevel using fmt.Sprintf style, extracting fields from the context, via the global logger.)
func CtxInfof(ctx context.Context, template string, args ...interface{}) {
	Std().CtxInfof(ctx, template, args...)
}

// CtxWarnf 使用 fmt.Sprintf 风格记录一条 Warn 级别的消息，并从 context 中提取字段，通过全局记录器。
// (CtxWarnf logs a message at WarnLevel using fmt.Sprintf style, extracting fields from the context, via the global logger.)
func CtxWarnf(ctx context.Context, template string, args ...interface{}) {
	Std().CtxWarnf(ctx, template, args...)
}

// CtxErrorf 使用 fmt.Sprintf 风格记录一条 Error 级别的消息，并从 context 中提取字段，通过全局记录器。
// (CtxErrorf logs a message at ErrorLevel using fmt.Sprintf style, extracting fields from the context, via the global logger.)
func CtxErrorf(ctx context.Context, template string, args ...interface{}) {
	Std().CtxErrorf(ctx, template, args...)
}

// CtxPanicf 使用 fmt.Sprintf 风格记录一条 Panic 级别的消息，然后发生 panic，并从 context 中提取字段，通过全局记录器。
// (CtxPanicf logs a message at PanicLevel using fmt.Sprintf style, then panics, extracting fields from the context, via the global logger.)
func CtxPanicf(ctx context.Context, template string, args ...interface{}) {
	Std().CtxPanicf(ctx, template, args...)
}

// CtxFatalf 使用 fmt.Sprintf 风格记录一条 Fatal 级别的消息，然后调用 os.Exit(1)，并从 context 中提取字段，通过全局记录器。
// (CtxFatalf logs a message at FatalLevel using fmt.Sprintf style, then calls os.Exit(1), extracting fields from the context, via the global logger.)
func CtxFatalf(ctx context.Context, template string, args ...interface{}) {
	Std().CtxFatalf(ctx, template, args...)
}

// --- Contextual Logging Implementation ---

// extractContextFields 从 context 中提取预定义的和用户指定的键。
// (extractContextFields extracts predefined and user-specified keys from the context.)
func extractContextFields(ctx context.Context, userKeys []any) []zap.Field {
	if ctx == nil {
		return nil
	}

	fields := make([]zap.Field, 0, 2+len(userKeys)) // Preallocate slice

	// 提取预定义的键 (Extract predefined keys)
	if traceID, ok := TraceIDFromContext(ctx); ok {
		// 注意：键的名称需要确定，这里暂时使用 "trace_id"
		// (Note: The key name needs to be decided, using "trace_id" for now)
		fields = append(fields, zap.String("trace_id", traceID))
	}
	if requestID, ok := RequestIDFromContext(ctx); ok {
		// 同样，键名需要确定，使用 "request_id"
		// (Similarly, key name needs to be decided, using "request_id")
		fields = append(fields, zap.String("request_id", requestID))
	}

	// 提取用户指定的键 (Extract user-specified keys)
	for _, key := range userKeys {
		if value := ctx.Value(key); value != nil {
			// 将键转换为字符串（如果可能），否则使用默认名称
			// (Convert key to string if possible, otherwise use default name)
			keyStr := fmt.Sprintf("%v", key)
			
			// 尝试将值转换为字符串（如果可能）
			// (Try to convert value to string if possible)
			if strVal, ok := value.(string); ok {
				fields = append(fields, zap.String(keyStr, strVal))
			} else {
				// 如果不是字符串，则使用 fmt.Sprintf 转换为字符串
				// (If not a string, use fmt.Sprintf to convert to string)
				strVal := fmt.Sprintf("%v", value)
				fields = append(fields, zap.String(keyStr, strVal))
			}
		}
	}

	return fields
}

func (l *logger) CtxDebugf(ctx context.Context, template string, args ...interface{}) {
	if l.zapLogger.Level() <= zapcore.DebugLevel {
		fields := extractContextFields(ctx, l.opts.ContextKeys)
		l.zapLogger.With(fields...).Sugar().Debugf(template, args...)
	}
}

func (l *logger) CtxInfof(ctx context.Context, template string, args ...interface{}) {
	if l.zapLogger.Level() <= zapcore.InfoLevel {
		fields := extractContextFields(ctx, l.opts.ContextKeys)
		l.zapLogger.With(fields...).Sugar().Infof(template, args...)
	}
}

func (l *logger) CtxWarnf(ctx context.Context, template string, args ...interface{}) {
	if l.zapLogger.Level() <= zapcore.WarnLevel {
		fields := extractContextFields(ctx, l.opts.ContextKeys)
		l.zapLogger.With(fields...).Sugar().Warnf(template, args...)
	}
}

func (l *logger) CtxErrorf(ctx context.Context, template string, args ...interface{}) {
	if l.zapLogger.Level() <= zapcore.ErrorLevel {
		fields := extractContextFields(ctx, l.opts.ContextKeys)
		l.zapLogger.With(fields...).Sugar().Errorf(template, args...)
	}
}

func (l *logger) CtxPanicf(ctx context.Context, template string, args ...interface{}) {
	// Panic level always logs regardless of configured level
	fields := extractContextFields(ctx, l.opts.ContextKeys)
	l.zapLogger.With(fields...).Sugar().Panicf(template, args...)
}

func (l *logger) CtxFatalf(ctx context.Context, template string, args ...interface{}) {
	// Fatal level always logs regardless of configured level
	fields := extractContextFields(ctx, l.opts.ContextKeys)
	l.zapLogger.With(fields...).Sugar().Fatalf(template, args...)
}
