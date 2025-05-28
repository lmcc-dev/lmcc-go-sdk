/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * Stack traces example demonstrating capture and analysis of call stacks.
 */

package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

// FileProcessor 文件处理器接口
// (FileProcessor represents file processor interface)
type FileProcessor interface {
	ProcessFile(filename string) error
	ValidateFile(filename string) error
	BackupFile(filename string) error
}

// DocumentProcessor 文档处理器实现
// (DocumentProcessor implements document processor)
type DocumentProcessor struct {
	maxSize int64
	timeout time.Duration
}

// NewDocumentProcessor 创建文档处理器
// (NewDocumentProcessor creates a document processor)
func NewDocumentProcessor() *DocumentProcessor {
	return &DocumentProcessor{
		maxSize: 10 * 1024 * 1024, // 10MB
		timeout: 30 * time.Second,
	}
}

// ProcessFile 处理文件（演示深层调用栈）
// (ProcessFile processes a file - demonstrates deep call stack)
func (dp *DocumentProcessor) ProcessFile(filename string) error {
	// 第一层：文件处理入口 (First layer: file processing entry)
	if err := dp.validateFileExtension(filename); err != nil {
		return errors.Wrap(err, "file extension validation failed")
	}
	
	// 第二层：文件内容处理 (Second layer: file content processing)
	if err := dp.processFileContent(filename); err != nil {
		return errors.Wrapf(err, "failed to process content of file %s", filename)
	}
	
	return nil
}

// validateFileExtension 验证文件扩展名
// (validateFileExtension validates file extension)
func (dp *DocumentProcessor) validateFileExtension(filename string) error {
	if strings.HasSuffix(filename, ".tmp") {
		return errors.New("temporary files are not allowed")
	}
	
	if strings.HasSuffix(filename, ".lock") {
		return errors.New("lock files cannot be processed")
	}
	
	// 调用更深层的验证 (Call deeper validation)
	return dp.validateFilePermissions(filename)
}

// validateFilePermissions 验证文件权限
// (validateFilePermissions validates file permissions)
func (dp *DocumentProcessor) validateFilePermissions(filename string) error {
	// 模拟权限检查失败 (Simulate permission check failure)
	if strings.Contains(filename, "restricted") {
		return dp.createPermissionError(filename)
	}
	
	return nil
}

// createPermissionError 创建权限错误（更深层的调用）
// (createPermissionError creates permission error - deeper call)
func (dp *DocumentProcessor) createPermissionError(filename string) error {
	return dp.buildDetailedPermissionError(filename, "read access denied")
}

// buildDetailedPermissionError 构建详细的权限错误（最深层）
// (buildDetailedPermissionError builds detailed permission error - deepest layer)
func (dp *DocumentProcessor) buildDetailedPermissionError(filename, reason string) error {
	return errors.Errorf("permission denied for file %s: %s", filename, reason)
}

// processFileContent 处理文件内容
// (processFileContent processes file content)
func (dp *DocumentProcessor) processFileContent(filename string) error {
	// 第三层：内容解析 (Third layer: content parsing)
	if err := dp.parseFileContent(filename); err != nil {
		return errors.Wrap(err, "content parsing failed")
	}
	
	// 第四层：内容验证 (Fourth layer: content validation)
	if err := dp.validateContent(filename); err != nil {
		return errors.Wrap(err, "content validation failed")
	}
	
	return nil
}

// parseFileContent 解析文件内容
// (parseFileContent parses file content)
func (dp *DocumentProcessor) parseFileContent(filename string) error {
	// 模拟解析错误 (Simulate parsing error)
	if strings.Contains(filename, "corrupt") {
		return dp.handleParsingError(filename, "invalid format detected")
	}
	
	if strings.Contains(filename, "large") {
		return dp.handleSizeError(filename, dp.maxSize)
	}
	
	return nil
}

// handleParsingError 处理解析错误
// (handleParsingError handles parsing error)
func (dp *DocumentProcessor) handleParsingError(filename, details string) error {
	return dp.createDetailedParsingError(filename, details, getParserVersion())
}

// createDetailedParsingError 创建详细的解析错误
// (createDetailedParsingError creates detailed parsing error)
func (dp *DocumentProcessor) createDetailedParsingError(filename, details, version string) error {
	return errors.Errorf("parsing failed for %s with parser %s: %s", filename, version, details)
}

// handleSizeError 处理大小错误
// (handleSizeError handles size error)
func (dp *DocumentProcessor) handleSizeError(filename string, maxSize int64) error {
	return dp.createSizeError(filename, maxSize, getCurrentFileSize(filename))
}

// createSizeError 创建大小错误
// (createSizeError creates size error)
func (dp *DocumentProcessor) createSizeError(filename string, maxSize, actualSize int64) error {
	return errors.Errorf("file %s exceeds size limit: %d bytes > %d bytes", 
		filename, actualSize, maxSize)
}

// validateContent 验证内容
// (validateContent validates content)
func (dp *DocumentProcessor) validateContent(filename string) error {
	// 模拟内容验证错误 (Simulate content validation error)
	if strings.Contains(filename, "invalid") {
		return dp.validateContentFormat(filename)
	}
	
	return nil
}

// validateContentFormat 验证内容格式
// (validateContentFormat validates content format)
func (dp *DocumentProcessor) validateContentFormat(filename string) error {
	return dp.checkContentSchema(filename, "document-v1.0")
}

// checkContentSchema 检查内容模式
// (checkContentSchema checks content schema)
func (dp *DocumentProcessor) checkContentSchema(filename, schema string) error {
	return dp.performSchemaValidation(filename, schema, loadSchemaRules(schema))
}

// performSchemaValidation 执行模式验证
// (performSchemaValidation performs schema validation)
func (dp *DocumentProcessor) performSchemaValidation(filename, schema string, rules []string) error {
	return errors.Errorf("schema validation failed for %s against %s: missing required fields %v", 
		filename, schema, rules)
}

// ValidateFile 验证文件（演示简单调用栈）
// (ValidateFile validates a file - demonstrates simple call stack)
func (dp *DocumentProcessor) ValidateFile(filename string) error {
	if filename == "" {
		return errors.New("filename cannot be empty")
	}
	
	return errors.Errorf("file %s does not exist", filename)
}

// BackupFile 备份文件（演示包装的调用栈）
// (BackupFile backs up a file - demonstrates wrapped call stack)
func (dp *DocumentProcessor) BackupFile(filename string) error {
	if err := dp.createBackupDirectory(); err != nil {
		return errors.Wrap(err, "backup directory creation failed")
	}
	
	if err := dp.copyFileToBackup(filename); err != nil {
		return errors.Wrapf(err, "failed to copy %s to backup", filename)
	}
	
	return nil
}

// createBackupDirectory 创建备份目录
// (createBackupDirectory creates backup directory)
func (dp *DocumentProcessor) createBackupDirectory() error {
	return errors.New("insufficient disk space for backup directory")
}

// copyFileToBackup 复制文件到备份
// (copyFileToBackup copies file to backup)
func (dp *DocumentProcessor) copyFileToBackup(filename string) error {
	return errors.New("backup storage service unavailable")
}

// 辅助函数 (Helper functions)

// getParserVersion 获取解析器版本
// (getParserVersion gets parser version)
func getParserVersion() string {
	return "v2.1.3"
}

// getCurrentFileSize 获取当前文件大小
// (getCurrentFileSize gets current file size)
func getCurrentFileSize(filename string) int64 {
	return 50 * 1024 * 1024 // 50MB
}

// loadSchemaRules 加载模式规则
// (loadSchemaRules loads schema rules)
func loadSchemaRules(schema string) []string {
	return []string{"title", "author", "date"}
}

// StackAnalyzer 堆栈分析器
// (StackAnalyzer analyzes stack traces)
type StackAnalyzer struct {
	showFullPath bool
	maxDepth     int
}

// NewStackAnalyzer 创建堆栈分析器
// (NewStackAnalyzer creates a stack analyzer)
func NewStackAnalyzer() *StackAnalyzer {
	return &StackAnalyzer{
		showFullPath: false,
		maxDepth:     10,
	}
}

// AnalyzeError 分析错误的堆栈跟踪
// (AnalyzeError analyzes error stack trace)
func (sa *StackAnalyzer) AnalyzeError(err error, description string) {
	fmt.Printf("=== Stack Trace Analysis: %s ===\n", description)
	
	if err == nil {
		fmt.Println("No error to analyze")
		return
	}
	
	fmt.Printf("Error: %v\n", err)
	fmt.Printf("Error Type: %T\n", err)
	
	// 使用 %+v 显示完整的堆栈跟踪 (Use %+v to show full stack trace)
	fmt.Println("\nDetailed Stack Trace:")
	fmt.Printf("%+v\n", err)
	
	// 分析堆栈深度 (Analyze stack depth)
	sa.analyzeStackDepth(err)
	
	// 检查错误链 (Check error chain)
	sa.analyzeErrorChain(err)
	
	fmt.Println()
}

// analyzeStackDepth 分析堆栈深度
// (analyzeStackDepth analyzes stack depth)
func (sa *StackAnalyzer) analyzeStackDepth(err error) {
	stackStr := fmt.Sprintf("%+v", err)
	lines := strings.Split(stackStr, "\n")
	
	stackLines := 0
	for _, line := range lines {
		if strings.Contains(line, ".go:") {
			stackLines++
		}
	}
	
	fmt.Printf("\nStack Analysis:\n")
	fmt.Printf("  Total stack lines: %d\n", stackLines)
	
	if stackLines > 20 {
		fmt.Printf("  ⚠️  Deep stack detected (>20 frames)\n")
	} else if stackLines > 10 {
		fmt.Printf("  ⚠️  Medium stack depth (10-20 frames)\n")
	} else {
		fmt.Printf("  ✓ Normal stack depth (<%10 frames)\n")
	}
}

// analyzeErrorChain 分析错误链
// (analyzeErrorChain analyzes error chain)
func (sa *StackAnalyzer) analyzeErrorChain(err error) {
	fmt.Println("\nError Chain:")
	depth := 0
	current := err
	
	for current != nil {
		indent := strings.Repeat("  ", depth)
		fmt.Printf("%s%d. %v\n", indent, depth+1, current)
		
		// 尝试展开错误 (Try to unwrap error)
		if unwrapper, ok := current.(interface{ Unwrap() error }); ok {
			current = unwrapper.Unwrap()
		} else {
			current = nil
		}
		depth++
		
		if depth > sa.maxDepth {
			fmt.Printf("%s... (truncated at depth %d)\n", strings.Repeat("  ", depth), sa.maxDepth)
			break
		}
	}
	
	fmt.Printf("Chain depth: %d\n", depth)
}

// RuntimeAnalyzer 运行时分析器
// (RuntimeAnalyzer analyzes runtime information)
type RuntimeAnalyzer struct{}

// NewRuntimeAnalyzer 创建运行时分析器
// (NewRuntimeAnalyzer creates a runtime analyzer)
func NewRuntimeAnalyzer() *RuntimeAnalyzer {
	return &RuntimeAnalyzer{}
}

// PrintCurrentStack 打印当前调用栈
// (PrintCurrentStack prints current call stack)
func (ra *RuntimeAnalyzer) PrintCurrentStack() {
	fmt.Println("=== Current Runtime Stack ===")
	
	// 获取当前调用栈 (Get current call stack)
	pc := make([]uintptr, 32) // 最多32帧 (max 32 frames)
	n := runtime.Callers(1, pc)
	pc = pc[:n]
	
	frames := runtime.CallersFrames(pc)
	
	fmt.Println("Call stack from current location:")
	frameNum := 0
	for {
		frame, more := frames.Next()
		if !more {
			break
		}
		
		fmt.Printf("%d. %s\n", frameNum+1, frame.Function)
		fmt.Printf("   %s:%d\n", frame.File, frame.Line)
		
		frameNum++
		if frameNum >= 10 { // 限制显示的帧数 (limit displayed frames)
			fmt.Printf("   ... (showing first 10 frames)\n")
			break
		}
	}
	
	fmt.Printf("Total frames: %d\n\n", n)
}

// CompareStackTraces 比较堆栈跟踪
// (CompareStackTraces compares stack traces)
func (ra *RuntimeAnalyzer) CompareStackTraces(err1, err2 error, desc1, desc2 string) {
	fmt.Printf("=== Comparing Stack Traces ===\n")
	fmt.Printf("Error 1: %s\n", desc1)
	fmt.Printf("Error 2: %s\n", desc2)
	fmt.Println()
	
	// 获取两个错误的堆栈信息 (Get stack info for both errors)
	stack1 := fmt.Sprintf("%+v", err1)
	stack2 := fmt.Sprintf("%+v", err2)
	
	lines1 := strings.Split(stack1, "\n")
	lines2 := strings.Split(stack2, "\n")
	
	// 计算公共前缀 (Calculate common prefix)
	commonFrames := 0
	minLen := len(lines1)
	if len(lines2) < minLen {
		minLen = len(lines2)
	}
	
	for i := 0; i < minLen; i++ {
		if strings.Contains(lines1[i], ".go:") && strings.Contains(lines2[i], ".go:") {
			if lines1[i] == lines2[i] {
				commonFrames++
			} else {
				break
			}
		}
	}
	
	fmt.Printf("Stack comparison results:\n")
	fmt.Printf("  Error 1 stack depth: %d frames\n", countStackFrames(stack1))
	fmt.Printf("  Error 2 stack depth: %d frames\n", countStackFrames(stack2))
	fmt.Printf("  Common frames: %d\n", commonFrames)
	
	if commonFrames > 0 {
		fmt.Printf("  ✓ Errors share common call path\n")
	} else {
		fmt.Printf("  ⚠️  Errors have different call paths\n")
	}
	
	fmt.Println()
}

// countStackFrames 计算堆栈帧数
// (countStackFrames counts stack frames)
func countStackFrames(stack string) int {
	lines := strings.Split(stack, "\n")
	count := 0
	for _, line := range lines {
		if strings.Contains(line, ".go:") {
			count++
		}
	}
	return count
}

// demonstrateDeepStackTrace 演示深层堆栈跟踪
// (demonstrateDeepStackTrace demonstrates deep stack trace)
func demonstrateDeepStackTrace() {
	fmt.Println("=== Demonstrating Deep Stack Traces ===")
	fmt.Println()
	
	processor := NewDocumentProcessor()
	analyzer := NewStackAnalyzer()
	
	// 测试深层调用栈的错误 (Test errors with deep call stacks)
	testCases := []struct {
		filename    string
		description string
	}{
		{"restricted_document.pdf", "Permission Error (Deep Stack)"},
		{"corrupt_file.doc", "Parsing Error (Deep Stack)"},
		{"large_document.txt", "Size Error (Deep Stack)"},
		{"invalid_content.xml", "Schema Validation Error (Very Deep Stack)"},
	}
	
	for _, tc := range testCases {
		fmt.Printf("Testing: %s\n", tc.description)
		err := processor.ProcessFile(tc.filename)
		if err != nil {
			analyzer.AnalyzeError(err, tc.description)
		}
	}
}

// demonstrateSimpleStackTrace 演示简单堆栈跟踪
// (demonstrateSimpleStackTrace demonstrates simple stack trace)
func demonstrateSimpleStackTrace() {
	fmt.Println("=== Demonstrating Simple Stack Traces ===")
	fmt.Println()
	
	processor := NewDocumentProcessor()
	analyzer := NewStackAnalyzer()
	
	// 测试简单的错误 (Test simple errors)
	fmt.Println("Testing: Simple Validation Error")
	err := processor.ValidateFile("")
	if err != nil {
		analyzer.AnalyzeError(err, "Simple Validation Error")
	}
	
	fmt.Println("Testing: File Not Found Error")
	err = processor.ValidateFile("nonexistent.txt")
	if err != nil {
		analyzer.AnalyzeError(err, "File Not Found Error")
	}
}

// demonstrateWrappedStackTrace 演示包装的堆栈跟踪
// (demonstrateWrappedStackTrace demonstrates wrapped stack trace)
func demonstrateWrappedStackTrace() {
	fmt.Println("=== Demonstrating Wrapped Stack Traces ===")
	fmt.Println()
	
	processor := NewDocumentProcessor()
	analyzer := NewStackAnalyzer()
	
	// 测试包装的错误 (Test wrapped errors)
	fmt.Println("Testing: Backup Operation with Wrapping")
	err := processor.BackupFile("important_document.pdf")
	if err != nil {
		analyzer.AnalyzeError(err, "Wrapped Backup Error")
	}
}

// demonstrateRuntimeStack 演示运行时堆栈
// (demonstrateRuntimeStack demonstrates runtime stack)
func demonstrateRuntimeStack() {
	fmt.Println("=== Demonstrating Runtime Stack Analysis ===")
	fmt.Println()
	
	runtime := NewRuntimeAnalyzer()
	
	// 显示当前调用栈 (Show current call stack)
	runtime.PrintCurrentStack()
}

// demonstrateStackComparison 演示堆栈比较
// (demonstrateStackComparison demonstrates stack comparison)
func demonstrateStackComparison() {
	fmt.Println("=== Demonstrating Stack Trace Comparison ===")
	fmt.Println()
	
	processor := NewDocumentProcessor()
	runtime := NewRuntimeAnalyzer()
	
	// 创建两个不同的错误进行比较 (Create two different errors for comparison)
	err1 := processor.ProcessFile("restricted_document.pdf")
	err2 := processor.ProcessFile("corrupt_file.doc")
	
	if err1 != nil && err2 != nil {
		runtime.CompareStackTraces(err1, err2, "Permission Error", "Parsing Error")
	}
	
	// 比较包装错误和简单错误 (Compare wrapped error and simple error)
	err3 := processor.BackupFile("test.txt")
	err4 := processor.ValidateFile("")
	
	if err3 != nil && err4 != nil {
		runtime.CompareStackTraces(err3, err4, "Wrapped Backup Error", "Simple Validation Error")
	}
}

// demonstratePerformanceImpact 演示性能影响
// (demonstratePerformanceImpact demonstrates performance impact)
func demonstratePerformanceImpact() {
	fmt.Println("=== Demonstrating Stack Trace Performance Impact ===")
	fmt.Println()
	
	// 测试有堆栈跟踪的错误创建性能 (Test error creation performance with stack traces)
	iterations := 10000
	
	// 测试1：使用 errors.New（包含堆栈跟踪）
	// (Test 1: Using errors.New (includes stack traces))
	start := time.Now()
	for i := 0; i < iterations; i++ {
		_ = errors.New("test error with stack trace")
	}
	withStack := time.Since(start)
	
	// 测试2：使用标准库 fmt.Errorf（无堆栈跟踪）
	// (Test 2: Using standard library fmt.Errorf (no stack traces))
	start = time.Now()
	for i := 0; i < iterations; i++ {
		_ = fmt.Errorf("test error without stack trace")
	}
	withoutStack := time.Since(start)
	
	fmt.Printf("Performance comparison (%d iterations):\n", iterations)
	fmt.Printf("  With stack traces:    %v\n", withStack)
	fmt.Printf("  Without stack traces: %v\n", withoutStack)
	fmt.Printf("  Overhead ratio:       %.2fx\n", float64(withStack)/float64(withoutStack))
	
	if withStack > withoutStack*2 {
		fmt.Printf("  ⚠️  Significant overhead detected\n")
	} else {
		fmt.Printf("  ✓ Acceptable overhead\n")
	}
	
	fmt.Println()
}

func main() {
	fmt.Println("=== Stack Traces Example ===")
	fmt.Println("This example demonstrates stack trace capture and analysis.")
	fmt.Println()
	
	// 1. 初始化日志 (Initialize logging)
	logOpts := log.NewOptions()
	logOpts.Level = "info"
	logOpts.Format = "text"
	logOpts.EnableColor = true
	log.Init(logOpts)
	logger := log.Std()
	
	// 2. 演示深层堆栈跟踪 (Demonstrate deep stack traces)
	demonstrateDeepStackTrace()
	
	// 3. 演示简单堆栈跟踪 (Demonstrate simple stack traces)
	demonstrateSimpleStackTrace()
	
	// 4. 演示包装的堆栈跟踪 (Demonstrate wrapped stack traces)
	demonstrateWrappedStackTrace()
	
	// 5. 演示运行时堆栈分析 (Demonstrate runtime stack analysis)
	demonstrateRuntimeStack()
	
	// 6. 演示堆栈比较 (Demonstrate stack comparison)
	demonstrateStackComparison()
	
	// 7. 演示性能影响 (Demonstrate performance impact)
	demonstratePerformanceImpact()
	
	logger.Info("Stack traces example completed successfully")
	fmt.Println("=== Example completed successfully ===")
} 