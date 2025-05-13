/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * Contains tests for log file rotation functionality.
 */

package log_test

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

// TestLogRotationBySize tests log rotation triggered by size limit.
// (TestLogRotationBySize 测试由大小限制触发的日志轮转。)
func TestLogRotationBySize(t *testing.T) {
	localRequire := require.New(t)
	localAssert := assert.New(t)

	tempDir := t.TempDir()
	t.Logf("Test temporary directory: %s", tempDir)
	logBaseName := "rotate_size.log"
	logFilePath := filepath.Join(tempDir, logBaseName)

	opts := log.NewOptions()
	opts.OutputPaths = []string{logFilePath}
	opts.Format = log.FormatText 
	opts.Level = zapcore.InfoLevel.String()
	opts.EnableColor = false
	opts.LogRotateMaxSize = 1    
	opts.LogRotateMaxBackups = 3 
	opts.LogRotateMaxAge = 0     
	opts.LogRotateCompress = false

	logger := log.NewLogger(opts)
	localRequire.NotNil(logger)
	defer func() { _ = logger.Sync() }()

	payload := "payload-" + strings.Repeat("X", 50) 
	estimatedEntrySizeBytes := 150 + len(payload)
	
	targetFirstSizeBytes := int(math.Round(float64(1024*1024) * 1.1))
	lineCountFirstCheck := targetFirstSizeBytes / estimatedEntrySizeBytes
	
	t.Logf("First phase: Writing ~1.1MB (约 %d lines) to log file", lineCountFirstCheck)
	for i := 0; i < lineCountFirstCheck; i++ {
		logger.Infow("Rotation test entry", "phase", "first", "index", i, "payload", payload)
		if i > 0 && i%200 == 0 {
			_ = logger.Sync()
		}
	}
	
	err := logger.Sync()
	localRequire.NoError(err)
	
	time.Sleep(100 * time.Millisecond)
	
	fileInfo, err := os.Stat(logFilePath)
	localRequire.NoError(err)
	firstPhaseSize := fileInfo.Size()
	t.Logf("First phase complete: Log file size is %d bytes (~%.2f MB)", firstPhaseSize, float64(firstPhaseSize)/(1024*1024))
	
	files, err := filepath.Glob(filepath.Join(tempDir, logBaseName+"*"))
	localRequire.NoError(err)
	localAssert.Equal(1, len(files), "Before 1MB threshold, should only have 1 log file (no rotation yet)")
	
	additionalBytesNeeded := (1024*1024) - int(firstPhaseSize) + (300*1024) 
	additionalLinesNeeded := additionalBytesNeeded / estimatedEntrySizeBytes
	
	t.Logf("Second phase: Writing additional ~%d bytes (~%d lines) to trigger rotation", additionalBytesNeeded, additionalLinesNeeded)
	for i := 0; i < additionalLinesNeeded; i++ {
		logger.Infow("Rotation test entry", "phase", "second", "index", i, "payload", payload)
		if i > 0 && i%100 == 0 {
			_ = logger.Sync()
		}
	}
	
	err = logger.Sync()
	localRequire.NoError(err)
	
time.Sleep(1 * time.Second)
	
	fileInfo, err = os.Stat(logFilePath)
	if err != nil {
		t.Logf("Error checking main log file: %v (this might be expected if rotation occurred)", err)
	} else {
		secondPhaseSize := fileInfo.Size()
		t.Logf("After second phase: Main log file size is %d bytes (~%.2f MB)", secondPhaseSize, float64(secondPhaseSize)/(1024*1024))
	}
	
entries, err := os.ReadDir(tempDir)
	localRequire.NoError(err)
	
t.Logf("调试: 临时目录内容:")
	logFileCount := 0
	hasRotatedFiles := false
	
	for _, entry := range entries {
		entryName := entry.Name()
		info, entryInfoErr := entry.Info()
		
		if entryInfoErr != nil {
			t.Logf("  - %s (无法获取信息: %v)", entryName, entryInfoErr)
			continue
		}
		
t.Logf("  - %s (大小: %d bytes, 修改时间: %v)", 
			entryName, info.Size(), info.ModTime())
		
		baseNameWithoutExt := strings.TrimSuffix(logBaseName, filepath.Ext(logBaseName))
		rotationPattern := baseNameWithoutExt + "-"
		
t.Logf("    检测: 基础名称=%s, 无扩展名基础=%s, 轮转模式=%s, 当前文件=%s, 是否匹配基本名=%v, 是否包含轮转模式=%v",
			logBaseName, baseNameWithoutExt, rotationPattern, entryName, 
			entryName == logBaseName, strings.Contains(entryName, rotationPattern))
		
		if entryName == logBaseName {
			t.Logf("    匹配到主日志文件: %s", entryName)
			logFileCount++
		} else if strings.Contains(entryName, rotationPattern) {
			t.Logf("    匹配到轮转备份文件: %s", entryName)
			logFileCount++
			hasRotatedFiles = true
		}
	}
	
t.Logf("发现 %d 个日志文件，其中包含轮转文件: %v", logFileCount, hasRotatedFiles)
	
	localAssert.GreaterOrEqual(logFileCount, 2, "After exceeding 1MB, should have original log file and at least one backup")
	localAssert.True(hasRotatedFiles, "Should find at least one rotated backup file")
}

// TestLogRotationMaxBackups tests the max backups limit.
// (TestLogRotationMaxBackups 测试最大备份文件数限制。)
func TestLogRotationMaxBackups(t *testing.T) {
	localRequire := require.New(t)
	localAssert := assert.New(t)

	tempDir := t.TempDir()
	logBaseName := "rotate_backups.log"
	logFilePath := filepath.Join(tempDir, logBaseName)

	maxBackups := 2
	opts := log.NewOptions()
	opts.OutputPaths = []string{logFilePath}
	opts.Format = log.FormatText
	opts.Level = zapcore.InfoLevel.String()
	opts.EnableColor = false
	opts.LogRotateMaxSize = 1 
	opts.LogRotateMaxBackups = maxBackups 
	opts.LogRotateMaxAge = 0
	opts.LogRotateCompress = false

	logger := log.NewLogger(opts)
	localRequire.NotNil(logger)
	defer func() { _ = logger.Sync() }()

	logLine := "Testing max backups rotation. Log entry number: "
	logEntrySize := len(logLine) + 10
	linesPerMB := (1 * 1024 * 1024) / logEntrySize
	linesPerMB += 50 

	totalLines := linesPerMB * (maxBackups + 2)
	for i := 0; i < totalLines; i++ {
		logger.Infof("%s%d", logLine, i)
	}
	err := logger.Sync()
	localRequire.NoError(err)

	time.Sleep(500 * time.Millisecond)

	entries, err := os.ReadDir(tempDir)
	localRequire.NoError(err)
	
	logFileCount := 0
	originalFound := false
	backupCount := 0
	
	baseNameWithoutExt := strings.TrimSuffix(logBaseName, filepath.Ext(logBaseName))
	rotationPattern := baseNameWithoutExt + "-"
	
	for _, entry := range entries {
		entryName := entry.Name()
		
		if entryName == logBaseName {
			originalFound = true
			logFileCount++
		} else if strings.Contains(entryName, rotationPattern) {
			backupCount++
			logFileCount++
		}
	}

	expectedFileCount := maxBackups + 1
	localAssert.Equal(expectedFileCount, logFileCount, "Total number of log files should be current + maxBackups")

	localAssert.True(originalFound, "Original log file should exist")
	localAssert.Equal(maxBackups, backupCount, "Number of backup files should match maxBackups")
}

// TestLogRotationCompress tests if compressed backup files are named correctly.
// (TestLogRotationCompress 测试压缩后的备份文件名是否正确。)
func TestLogRotationCompress(t *testing.T) {
	localRequire := require.New(t)
	localAssert := assert.New(t)

	tempDir := t.TempDir()
	logBaseName := "rotate_compress.log"
	logFilePath := filepath.Join(tempDir, logBaseName)

	opts := log.NewOptions()
	opts.OutputPaths = []string{logFilePath}
	opts.Format = log.FormatText
	opts.Level = zapcore.InfoLevel.String()
	opts.EnableColor = false
	opts.LogRotateMaxSize = 1 
	opts.LogRotateMaxBackups = 1
	opts.LogRotateMaxAge = 0
	opts.LogRotateCompress = true

	logger := log.NewLogger(opts)
	localRequire.NotNil(logger)
	defer func() { _ = logger.Sync() }()

	logLine := "Testing compression rotation. Log entry number: "
	logEntrySize := len(logLine) + 10
	linesPerMB := (1 * 1024 * 1024) / logEntrySize
	linesPerMB += 50

	totalLines := linesPerMB * 2
	for i := 0; i < totalLines; i++ {
		logger.Infof("%s%d", logLine, i)
	}
	err := logger.Sync()
	localRequire.NoError(err)

	time.Sleep(500 * time.Millisecond)

	entries, err := os.ReadDir(tempDir)
	localRequire.NoError(err)
	
	baseNameWithoutExt := strings.TrimSuffix(logBaseName, filepath.Ext(logBaseName))
	rotationPattern := baseNameWithoutExt + "-"
	
	logFileCount := 0
	backupFound := false
	compressedBackupFound := false
	
	for _, entry := range entries {
		entryName := entry.Name()
		
		if entryName == logBaseName {
			logFileCount++
		} else if strings.Contains(entryName, rotationPattern) {
			backupFound = true
			logFileCount++
			
			if strings.HasSuffix(entryName, ".gz") {
				compressedBackupFound = true
			}
		}
	}
	
	localAssert.GreaterOrEqual(logFileCount, 2, "Should have original log and at least one backup")
	localAssert.True(backupFound, "Should find at least one backup file")
	localAssert.True(compressedBackupFound, "Backup file name should end with .gz when compression is enabled")

	fileInfo, err := os.Stat(logFilePath)
	if localAssert.NoError(err) {
		localAssert.Less(fileInfo.Size(), int64(1024*1024), "Current log file size should be less than 1MB after rotation")
	}
} 