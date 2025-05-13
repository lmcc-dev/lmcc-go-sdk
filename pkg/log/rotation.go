/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package log

import (
	"os"
	"path/filepath" // Added for EnsureDir

	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// newRotateLogger 根据提供的选项为指定的文件路径创建一个启用了轮转的 zapcore.WriteSyncer。
// (newRotateLogger creates a zapcore.WriteSyncer with rotation enabled for the given file path, based on the provided options.)
func newRotateLogger(filePath string, opts *Options) (zapcore.WriteSyncer, error) {
	// 确保日志文件所在的目录存在 (Ensure the directory for the log file exists)
	if err := ensureDir(filePath); err != nil {
		return nil, err // Return error if directory creation fails
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
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 目录不存在，创建它 (Directory does not exist, create it)
		// 使用 0755 权限，允许所有者读/写/执行，组和其他人读/执行
		// (Use 0755 permissions: owner rwx, group rx, others rx)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err // Return error if MkdirAll fails
		}
	} else if err != nil {
		// Stat 返回了其他错误 (Stat returned some other error)
		return err
	}
	// 目录存在或已成功创建 (Directory exists or was successfully created)
	return nil
} 