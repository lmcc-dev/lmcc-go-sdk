/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package log

import (
	"os"
	"path/filepath" // Added for EnsureDir

	lmccerrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors" // SDK errors 包 (SDK errors package)
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// newRotateLogger 根据提供的选项为指定的文件路径创建一个启用了轮转的 zapcore.WriteSyncer。
// (newRotateLogger creates a zapcore.WriteSyncer with rotation enabled for the given file path, based on the provided options.)
func newRotateLogger(filePath string, opts *Options) (zapcore.WriteSyncer, error) {
	// 确保日志文件所在的目录存在 (Ensure the directory for the log file exists)
	if err := ensureDir(filePath); err != nil {
		// 使用 lmccerrors.ErrorfWithCode 包装错误，以添加堆栈跟踪、上下文和错误码。
		// (Wrap the error with lmccerrors.ErrorfWithCode to add stack trace, context, and error code.)
		// 确保错误链和 Coder 被正确保留 (Ensure error chain and Coder are properly preserved)
		return nil, lmccerrors.WithCode(
			lmccerrors.Wrapf(err, "failed to ensure directory for log file %s", filePath),
			lmccerrors.ErrLogRotationSetup,
		)
	}

	// 配置 lumberjack logger (Configure lumberjack logger)
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    opts.LogRotateMaxSize,    // MB
		MaxBackups: opts.LogRotateMaxBackups, // Number of backups
		MaxAge:     opts.LogRotateMaxAge,     // Days
		Compress:   opts.LogRotateCompress,   // Compress rotated files
		LocalTime:  true,                     // Use local time for timestamps in backup filenames
	}

	// 使用 AddSync 将其包装成 zapcore.WriteSyncer
	// (Wrap it with AddSync to make it a zapcore.WriteSyncer)
	return zapcore.AddSync(lumberjackLogger), nil
}

// ensureDir 确保给定文件路径的目录存在，如果不存在则创建它。
// (ensureDir ensures that the directory for the given file path exists, creating it if necessary.)
func ensureDir(filePath string) error {
	dir := filepath.Dir(filePath)
	// 使用 Stat 检查目录是否存在以及它是否确实是一个目录
	// (Use Stat to check if the directory exists and if it's actually a directory)
	if statInfo, err := os.Stat(dir); os.IsNotExist(err) {
		// 目录不存在，创建它 (Directory does not exist, create it)
		// 使用 0755 权限，允许所有者读/写/执行，组和其他人读/执行
		// (Use 0755 permissions: owner rwx, group rx, others rx)
		errMkdir := os.MkdirAll(dir, 0755)
		if errMkdir != nil {
			// 使用 lmccerrors.ErrorfWithCode 包装错误，以添加堆栈跟踪、上下文和错误码。
			// (Wrap the error with lmccerrors.ErrorfWithCode to add stack trace, context, and error code.)
			return lmccerrors.WithCode(
				lmccerrors.Wrapf(errMkdir, "failed to create directory %s", dir),
				lmccerrors.ErrLogRotationDirCreate,
			)
		}
	} else if err != nil {
		// Stat 返回了其他错误 (Stat returned some other error)
		// 使用 lmccerrors.ErrorfWithCode 包装错误。
		// (Wrap the error with lmccerrors.ErrorfWithCode.)
		return lmccerrors.WithCode(
			lmccerrors.Wrapf(err, "failed to stat directory %s", dir),
			lmccerrors.ErrLogRotationDirStat,
		)
	} else if !statInfo.IsDir() {
		// 路径存在但不是目录 (Path exists but is not a directory)
		return lmccerrors.ErrorfWithCode(lmccerrors.ErrLogRotationDirInvalid, "path %s exists but is not a directory", dir)
	}
	// 目录存在或已成功创建 (Directory exists or was successfully created)
	return nil
} 