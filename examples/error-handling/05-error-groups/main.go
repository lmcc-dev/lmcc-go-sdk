/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * Error groups example demonstrating aggregation and handling of multiple errors.
 */

package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

// ErrorGroup 错误组实现
// (ErrorGroup implementation for collecting multiple errors)
type ErrorGroup struct {
	mu     sync.Mutex
	errors []error
}

// NewErrorGroup 创建新的错误组
// (NewErrorGroup creates a new error group)
func NewErrorGroup() *ErrorGroup {
	return &ErrorGroup{
		errors: make([]error, 0),
	}
}

// Add 添加错误到组中
// (Add adds an error to the group)
func (eg *ErrorGroup) Add(err error) {
	if err == nil {
		return
	}
	
	eg.mu.Lock()
	defer eg.mu.Unlock()
	eg.errors = append(eg.errors, err)
}

// Errors 返回所有错误
// (Errors returns all collected errors)
func (eg *ErrorGroup) Errors() []error {
	eg.mu.Lock()
	defer eg.mu.Unlock()
	
	// 返回副本以避免竞态条件 (Return copy to avoid race conditions)
	result := make([]error, len(eg.errors))
	copy(result, eg.errors)
	return result
}

// HasErrors 检查是否有错误
// (HasErrors checks if there are any errors)
func (eg *ErrorGroup) HasErrors() bool {
	eg.mu.Lock()
	defer eg.mu.Unlock()
	return len(eg.errors) > 0
}

// Count 返回错误数量
// (Count returns the number of errors)
func (eg *ErrorGroup) Count() int {
	eg.mu.Lock()
	defer eg.mu.Unlock()
	return len(eg.errors)
}

// Clear 清空所有错误
// (Clear removes all errors)
func (eg *ErrorGroup) Clear() {
	eg.mu.Lock()
	defer eg.mu.Unlock()
	eg.errors = eg.errors[:0]
}

// Error 实现 error 接口
// (Error implements the error interface)
func (eg *ErrorGroup) Error() string {
	eg.mu.Lock()
	defer eg.mu.Unlock()
	
	if len(eg.errors) == 0 {
		return "no errors"
	}
	
	if len(eg.errors) == 1 {
		return eg.errors[0].Error()
	}
	
	result := fmt.Sprintf("%d errors occurred:", len(eg.errors))
	for i, err := range eg.errors {
		result += fmt.Sprintf("\n  [%d] %v", i+1, err)
	}
	return result
}

// First 返回第一个错误
// (First returns the first error)
func (eg *ErrorGroup) First() error {
	eg.mu.Lock()
	defer eg.mu.Unlock()
	
	if len(eg.errors) == 0 {
		return nil
	}
	return eg.errors[0]
}

// Last 返回最后一个错误
// (Last returns the last error)
func (eg *ErrorGroup) Last() error {
	eg.mu.Lock()
	defer eg.mu.Unlock()
	
	if len(eg.errors) == 0 {
		return nil
	}
	return eg.errors[len(eg.errors)-1]
}

// FilterByType 按类型过滤错误
// (FilterByType filters errors by type)
func (eg *ErrorGroup) FilterByType(errorType string) []error {
	eg.mu.Lock()
	defer eg.mu.Unlock()
	
	var filtered []error
	for _, err := range eg.errors {
		if coder := errors.GetCoder(err); coder != nil {
			if coder.String() == errorType {
				filtered = append(filtered, err)
			}
		}
	}
	return filtered
}

// GroupByType 按类型分组错误
// (GroupByType groups errors by type)
func (eg *ErrorGroup) GroupByType() map[string][]error {
	eg.mu.Lock()
	defer eg.mu.Unlock()
	
	groups := make(map[string][]error)
	
	for _, err := range eg.errors {
		var groupKey string
		if coder := errors.GetCoder(err); coder != nil {
			groupKey = coder.String()
		} else {
			groupKey = "unknown"
		}
		
		groups[groupKey] = append(groups[groupKey], err)
	}
	
	return groups
}

// MultiTaskProcessor 多任务处理器
// (MultiTaskProcessor processes multiple tasks)
type MultiTaskProcessor struct {
	workers     int
	timeout     time.Duration
	retryCount  int
	errorGroup  *ErrorGroup
}

// NewMultiTaskProcessor 创建多任务处理器
// (NewMultiTaskProcessor creates a multi-task processor)
func NewMultiTaskProcessor(workers int) *MultiTaskProcessor {
	return &MultiTaskProcessor{
		workers:    workers,
		timeout:    10 * time.Second,
		retryCount: 3,
		errorGroup: NewErrorGroup(),
	}
}

// Task 任务定义
// (Task represents a task to be processed)
type Task struct {
	ID          string
	Type        string
	Data        interface{}
	Priority    int
	MaxRetries  int
}

// TaskResult 任务结果
// (TaskResult represents the result of a task)
type TaskResult struct {
	TaskID    string
	Success   bool
	Error     error
	Duration  time.Duration
	Attempts  int
}

// ProcessTasks 处理多个任务（演示并行错误收集）
// (ProcessTasks processes multiple tasks - demonstrates parallel error collection)
func (mtp *MultiTaskProcessor) ProcessTasks(tasks []Task) []TaskResult {
	results := make([]TaskResult, len(tasks))
	var wg sync.WaitGroup
	
	// 使用工作池模式 (Use worker pool pattern)
	taskChan := make(chan int, len(tasks))
	
	// 启动工作协程 (Start worker goroutines)
	for i := 0; i < mtp.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for taskIndex := range taskChan {
				result := mtp.processTask(tasks[taskIndex])
				results[taskIndex] = result
				
				// 如果有错误，添加到错误组 (If there's an error, add to error group)
				if result.Error != nil {
					mtp.errorGroup.Add(result.Error)
				}
			}
		}()
	}
	
	// 分发任务 (Distribute tasks)
	for i := range tasks {
		taskChan <- i
	}
	close(taskChan)
	
	// 等待所有任务完成 (Wait for all tasks to complete)
	wg.Wait()
	
	return results
}

// processTask 处理单个任务
// (processTask processes a single task)
func (mtp *MultiTaskProcessor) processTask(task Task) TaskResult {
	start := time.Now()
	attempts := 0
	maxRetries := task.MaxRetries
	if maxRetries == 0 {
		maxRetries = mtp.retryCount
	}
	
	for attempts < maxRetries {
		attempts++
		
		// 模拟任务处理 (Simulate task processing)
		err := mtp.simulateTaskExecution(task)
		if err == nil {
			return TaskResult{
				TaskID:   task.ID,
				Success:  true,
				Duration: time.Since(start),
				Attempts: attempts,
			}
		}
		
		// 如果不是最后一次尝试，稍作等待 (If not the last attempt, wait a bit)
		if attempts < maxRetries {
			time.Sleep(time.Millisecond * 100)
		}
	}
	
	// 所有重试都失败了 (All retries failed)
	finalError := errors.Errorf("task %s failed after %d attempts", task.ID, attempts)
	return TaskResult{
		TaskID:   task.ID,
		Success:  false,
		Error:    finalError,
		Duration: time.Since(start),
		Attempts: attempts,
	}
}

// simulateTaskExecution 模拟任务执行
// (simulateTaskExecution simulates task execution)
func (mtp *MultiTaskProcessor) simulateTaskExecution(task Task) error {
	// 根据任务类型模拟不同的错误 (Simulate different errors based on task type)
	switch task.Type {
	case "network":
		if task.ID == "net_timeout" {
			return errors.New("network operation timed out")
		}
		if task.ID == "net_error" {
			return errors.New("network connection failed")
		}
	case "database":
		if task.ID == "db_lock" {
			return errors.New("database deadlock detected")
		}
		if task.ID == "db_constraint" {
			return errors.New("database constraint violation")
		}
	case "file":
		if task.ID == "file_perm" {
			return errors.New("file permission denied")
		}
		if task.ID == "file_notfound" {
			return errors.New("file not found")
		}
	case "computation":
		if task.ID == "compute_overflow" {
			return errors.New("arithmetic overflow")
		}
	}
	
	return nil // 任务成功 (Task succeeded)
}

// GetErrorSummary 获取错误摘要
// (GetErrorSummary gets error summary)
func (mtp *MultiTaskProcessor) GetErrorSummary() *ErrorGroup {
	return mtp.errorGroup
}

// ValidationProcessor 验证处理器
// (ValidationProcessor handles validation of multiple items)
type ValidationProcessor struct {
	rules      []ValidationRule
	errorGroup *ErrorGroup
}

// ValidationRule 验证规则
// (ValidationRule defines a validation rule)
type ValidationRule struct {
	Name        string
	Description string
	Validator   func(data interface{}) error
}

// NewValidationProcessor 创建验证处理器
// (NewValidationProcessor creates a validation processor)
func NewValidationProcessor() *ValidationProcessor {
	return &ValidationProcessor{
		rules:      make([]ValidationRule, 0),
		errorGroup: NewErrorGroup(),
	}
}

// AddRule 添加验证规则
// (AddRule adds a validation rule)
func (vp *ValidationProcessor) AddRule(rule ValidationRule) {
	vp.rules = append(vp.rules, rule)
}

// ValidateData 验证数据（收集所有验证错误）
// (ValidateData validates data - collects all validation errors)
func (vp *ValidationProcessor) ValidateData(data interface{}) error {
	vp.errorGroup.Clear()
	
	// 对每个规则执行验证 (Execute validation for each rule)
	for _, rule := range vp.rules {
		if err := rule.Validator(data); err != nil {
			wrappedErr := errors.Wrapf(err, "validation rule '%s' failed", rule.Name)
			vp.errorGroup.Add(wrappedErr)
		}
	}
	
	// 如果有错误，返回组合错误 (If there are errors, return combined error)
	if vp.errorGroup.HasErrors() {
		return vp.errorGroup
	}
	
	return nil
}

// GetValidationErrors 获取验证错误
// (GetValidationErrors gets validation errors)
func (vp *ValidationProcessor) GetValidationErrors() *ErrorGroup {
	return vp.errorGroup
}

// 验证规则实现 (Validation rule implementations)

// UserData 用户数据结构
// (UserData represents user data structure)
type UserData struct {
	Name  string
	Email string
	Age   int
	Phone string
}

// createUserValidationRules 创建用户验证规则
// (createUserValidationRules creates user validation rules)
func createUserValidationRules() []ValidationRule {
	return []ValidationRule{
		{
			Name:        "name_required",
			Description: "Name is required",
			Validator: func(data interface{}) error {
				if user, ok := data.(*UserData); ok {
					if user.Name == "" {
						return errors.New("name cannot be empty")
					}
				}
				return nil
			},
		},
		{
			Name:        "name_length",
			Description: "Name length validation",
			Validator: func(data interface{}) error {
				if user, ok := data.(*UserData); ok {
					if len(user.Name) < 2 || len(user.Name) > 50 {
						return errors.Errorf("name length must be between 2 and 50 characters, got %d", len(user.Name))
					}
				}
				return nil
			},
		},
		{
			Name:        "email_required",
			Description: "Email is required",
			Validator: func(data interface{}) error {
				if user, ok := data.(*UserData); ok {
					if user.Email == "" {
						return errors.New("email cannot be empty")
					}
				}
				return nil
			},
		},
		{
			Name:        "email_format",
			Description: "Email format validation",
			Validator: func(data interface{}) error {
				if user, ok := data.(*UserData); ok {
					if user.Email != "" && !isValidEmail(user.Email) {
						return errors.Errorf("invalid email format: %s", user.Email)
					}
				}
				return nil
			},
		},
		{
			Name:        "age_range",
			Description: "Age range validation",
			Validator: func(data interface{}) error {
				if user, ok := data.(*UserData); ok {
					if user.Age < 0 || user.Age > 150 {
						return errors.Errorf("age must be between 0 and 150, got %d", user.Age)
					}
				}
				return nil
			},
		},
		{
			Name:        "phone_format",
			Description: "Phone format validation",
			Validator: func(data interface{}) error {
				if user, ok := data.(*UserData); ok {
					if user.Phone != "" && !isValidPhone(user.Phone) {
						return errors.Errorf("invalid phone format: %s", user.Phone)
					}
				}
				return nil
			},
		},
	}
}

// isValidEmail 简单的邮箱验证
// (isValidEmail simple email validation)
func isValidEmail(email string) bool {
	return len(email) > 3 && 
		   fmt.Sprintf("%s", email) != "" && 
		   (email[0] != '@' && email[len(email)-1] != '@')
}

// isValidPhone 简单的电话验证
// (isValidPhone simple phone validation)
func isValidPhone(phone string) bool {
	return len(phone) >= 10 && len(phone) <= 15
}

// BatchProcessor 批处理器
// (BatchProcessor handles batch operations)
type BatchProcessor struct {
	batchSize  int
	errorGroup *ErrorGroup
}

// NewBatchProcessor 创建批处理器
// (NewBatchProcessor creates a batch processor)
func NewBatchProcessor(batchSize int) *BatchProcessor {
	return &BatchProcessor{
		batchSize:  batchSize,
		errorGroup: NewErrorGroup(),
	}
}

// ProcessBatch 处理批量操作
// (ProcessBatch processes batch operations)
func (bp *BatchProcessor) ProcessBatch(items []interface{}) []error {
	bp.errorGroup.Clear()
	var allErrors []error
	
	// 分批处理 (Process in batches)
	for i := 0; i < len(items); i += bp.batchSize {
		end := i + bp.batchSize
		if end > len(items) {
			end = len(items)
		}
		
		batch := items[i:end]
		batchErrors := bp.processBatchItems(batch, i/bp.batchSize+1)
		
		// 收集批次错误 (Collect batch errors)
		for _, err := range batchErrors {
			if err != nil {
				bp.errorGroup.Add(err)
				allErrors = append(allErrors, err)
			}
		}
	}
	
	return allErrors
}

// processBatchItems 处理批次中的项目
// (processBatchItems processes items in a batch)
func (bp *BatchProcessor) processBatchItems(items []interface{}, batchNumber int) []error {
	var batchErrors []error
	
	for i, item := range items {
		err := bp.processItem(item, batchNumber, i)
		batchErrors = append(batchErrors, err)
	}
	
	return batchErrors
}

// processItem 处理单个项目
// (processItem processes a single item)
func (bp *BatchProcessor) processItem(item interface{}, batchNumber, itemIndex int) error {
	// 模拟处理不同类型的项目 (Simulate processing different types of items)
	switch v := item.(type) {
	case string:
		if v == "error_item" {
			return errors.Errorf("failed to process string item at batch %d, index %d", batchNumber, itemIndex)
		}
		if v == "" {
			return errors.Errorf("empty string item at batch %d, index %d", batchNumber, itemIndex)
		}
	case int:
		if v < 0 {
			return errors.Errorf("negative number %d at batch %d, index %d", v, batchNumber, itemIndex)
		}
		if v > 1000 {
			return errors.Errorf("number too large %d at batch %d, index %d", v, batchNumber, itemIndex)
		}
	case nil:
		return errors.Errorf("nil item at batch %d, index %d", batchNumber, itemIndex)
	}
	
	return nil // 处理成功 (Processing succeeded)
}

// GetBatchErrors 获取批处理错误
// (GetBatchErrors gets batch processing errors)
func (bp *BatchProcessor) GetBatchErrors() *ErrorGroup {
	return bp.errorGroup
}

// demonstrateBasicErrorGroup 演示基本错误组操作
// (demonstrateBasicErrorGroup demonstrates basic error group operations)
func demonstrateBasicErrorGroup() {
	fmt.Println("=== Demonstrating Basic Error Group Operations ===")
	fmt.Println()
	
	// 创建错误组 (Create error group)
	errorGroup := NewErrorGroup()
	
	// 添加一些错误 (Add some errors)
	errorGroup.Add(errors.New("first error"))
	errorGroup.Add(errors.New("second error"))
	errorGroup.Add(errors.Errorf("third error with value: %d", 42))
	errorGroup.Add(nil) // 这个会被忽略 (This will be ignored)
	
	fmt.Printf("Error count: %d\n", errorGroup.Count())
	fmt.Printf("Has errors: %t\n", errorGroup.HasErrors())
	fmt.Printf("First error: %v\n", errorGroup.First())
	fmt.Printf("Last error: %v\n", errorGroup.Last())
	
	fmt.Println("\nAll errors:")
	for i, err := range errorGroup.Errors() {
		fmt.Printf("  [%d] %v\n", i+1, err)
	}
	
	fmt.Println("\nCombined error message:")
	fmt.Printf("%v\n", errorGroup)
	
	fmt.Println()
}

// demonstrateParallelTaskProcessing 演示并行任务处理
// (demonstrateParallelTaskProcessing demonstrates parallel task processing)
func demonstrateParallelTaskProcessing() {
	fmt.Println("=== Demonstrating Parallel Task Processing ===")
	fmt.Println()
	
	// 创建多任务处理器 (Create multi-task processor)
	processor := NewMultiTaskProcessor(3) // 3个工作协程 (3 worker goroutines)
	
	// 创建测试任务 (Create test tasks)
	tasks := []Task{
		{ID: "task_1", Type: "network", Data: "data1", Priority: 1},
		{ID: "net_timeout", Type: "network", Data: "data2", Priority: 2},
		{ID: "task_3", Type: "database", Data: "data3", Priority: 1},
		{ID: "db_lock", Type: "database", Data: "data4", Priority: 3},
		{ID: "file_perm", Type: "file", Data: "data5", Priority: 2},
		{ID: "task_6", Type: "computation", Data: "data6", Priority: 1},
		{ID: "compute_overflow", Type: "computation", Data: "data7", Priority: 2},
		{ID: "net_error", Type: "network", Data: "data8", Priority: 1},
	}
	
	fmt.Printf("Processing %d tasks with %d workers...\n", len(tasks), 3)
	
	// 处理任务 (Process tasks)
	results := processor.ProcessTasks(tasks)
	
	// 显示结果 (Show results)
	successCount := 0
	failureCount := 0
	
	fmt.Println("\nTask Results:")
	for _, result := range results {
		if result.Success {
			successCount++
			fmt.Printf("  ✓ %s completed in %v (attempts: %d)\n", 
				result.TaskID, result.Duration, result.Attempts)
		} else {
			failureCount++
			fmt.Printf("  ✗ %s failed after %v (attempts: %d): %v\n", 
				result.TaskID, result.Duration, result.Attempts, result.Error)
		}
	}
	
	fmt.Printf("\nSummary: %d succeeded, %d failed\n", successCount, failureCount)
	
	// 显示错误摘要 (Show error summary)
	errorSummary := processor.GetErrorSummary()
	if errorSummary.HasErrors() {
		fmt.Printf("\nError Summary (%d errors):\n", errorSummary.Count())
		fmt.Printf("%v\n", errorSummary)
	}
	
	fmt.Println()
}

// demonstrateValidationErrors 演示验证错误收集
// (demonstrateValidationErrors demonstrates validation error collection)
func demonstrateValidationErrors() {
	fmt.Println("=== Demonstrating Validation Error Collection ===")
	fmt.Println()
	
	// 创建验证处理器 (Create validation processor)
	validator := NewValidationProcessor()
	
	// 添加验证规则 (Add validation rules)
	rules := createUserValidationRules()
	for _, rule := range rules {
		validator.AddRule(rule)
	}
	
	// 测试有效数据 (Test valid data)
	fmt.Println("Testing valid user data:")
	validUser := &UserData{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
		Phone: "1234567890",
	}
	
	err := validator.ValidateData(validUser)
	if err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	} else {
		fmt.Println("✓ Validation passed")
	}
	
	// 测试无效数据 (Test invalid data)
	fmt.Println("\nTesting invalid user data:")
	invalidUser := &UserData{
		Name:  "", // 无效：空名称 (Invalid: empty name)
		Email: "@invalid", // 无效：邮箱格式 (Invalid: email format)
		Age:   -5, // 无效：年龄范围 (Invalid: age range)
		Phone: "123", // 无效：电话格式 (Invalid: phone format)
	}
	
	err = validator.ValidateData(invalidUser)
	if err != nil {
		fmt.Printf("Validation failed with multiple errors:\n%v\n", err)
		
		// 显示错误详情 (Show error details)
		validationErrors := validator.GetValidationErrors()
		fmt.Printf("\nDetailed validation errors (%d total):\n", validationErrors.Count())
		for i, validationErr := range validationErrors.Errors() {
			fmt.Printf("  [%d] %v\n", i+1, validationErr)
		}
	}
	
	fmt.Println()
}

// demonstrateBatchProcessing 演示批处理错误
// (demonstrateBatchProcessing demonstrates batch processing errors)
func demonstrateBatchProcessing() {
	fmt.Println("=== Demonstrating Batch Processing Errors ===")
	fmt.Println()
	
	// 创建批处理器 (Create batch processor)
	batchProcessor := NewBatchProcessor(3) // 批大小为3 (Batch size 3)
	
	// 创建测试数据 (Create test data)
	items := []interface{}{
		"valid_item_1",
		"valid_item_2", 
		"error_item", // 这个会导致错误 (This will cause error)
		42,
		-10, // 负数会导致错误 (Negative number will cause error)
		"valid_item_3",
		nil, // nil会导致错误 (nil will cause error)
		1500, // 过大的数会导致错误 (Too large number will cause error)
		"",   // 空字符串会导致错误 (Empty string will cause error)
		"valid_item_4",
	}
	
	fmt.Printf("Processing %d items in batches of %d...\n", len(items), 3)
	
	// 处理批次 (Process batches)
	allErrors := batchProcessor.ProcessBatch(items)
	
	// 显示结果 (Show results)
	successCount := len(items) - len(allErrors)
	fmt.Printf("\nProcessing completed: %d succeeded, %d failed\n", successCount, len(allErrors))
	
	// 显示批处理错误 (Show batch processing errors)
	batchErrors := batchProcessor.GetBatchErrors()
	if batchErrors.HasErrors() {
		fmt.Printf("\nBatch Processing Errors (%d total):\n", batchErrors.Count())
		fmt.Printf("%v\n", batchErrors)
	}
	
	fmt.Println()
}

// demonstrateErrorGrouping 演示错误分组
// (demonstrateErrorGrouping demonstrates error grouping)
func demonstrateErrorGrouping() {
	fmt.Println("=== Demonstrating Error Grouping ===")
	fmt.Println()
	
	// 创建带有不同类型错误的错误组 (Create error group with different types of errors)
	errorGroup := NewErrorGroup()
	
	// 添加不同类型的错误 (Add different types of errors)
	networkErr1 := errors.New("connection timeout")
	networkErr2 := errors.New("host unreachable")
	dbErr1 := errors.New("deadlock detected")
	dbErr2 := errors.New("constraint violation")
	unknownErr := fmt.Errorf("unknown error")
	
	errorGroup.Add(networkErr1)
	errorGroup.Add(networkErr2)
	errorGroup.Add(dbErr1)
	errorGroup.Add(dbErr2)
	errorGroup.Add(unknownErr)
	
	fmt.Printf("Total errors: %d\n", errorGroup.Count())
	
	// 按类型分组错误 (Group errors by type)
	groupedErrors := errorGroup.GroupByType()
	
	fmt.Println("\nErrors grouped by type:")
	for errorType, errorList := range groupedErrors {
		fmt.Printf("  %s (%d errors):\n", errorType, len(errorList))
		for i, err := range errorList {
			fmt.Printf("    [%d] %v\n", i+1, err)
		}
	}
	
	fmt.Println()
}

func main() {
	fmt.Println("=== Error Groups Example ===")
	fmt.Println("This example demonstrates aggregation and handling of multiple errors.")
	fmt.Println()
	
	// 1. 初始化日志 (Initialize logging)
	logOpts := log.NewOptions()
	logOpts.Level = "info"
	logOpts.Format = "text"
	logOpts.EnableColor = true
	log.Init(logOpts)
	logger := log.Std()
	
	// 2. 演示基本错误组操作 (Demonstrate basic error group operations)
	demonstrateBasicErrorGroup()
	
	// 3. 演示并行任务处理 (Demonstrate parallel task processing)
	demonstrateParallelTaskProcessing()
	
	// 4. 演示验证错误收集 (Demonstrate validation error collection)
	demonstrateValidationErrors()
	
	// 5. 演示批处理错误 (Demonstrate batch processing errors)
	demonstrateBatchProcessing()
	
	// 6. 演示错误分组 (Demonstrate error grouping)
	demonstrateErrorGrouping()
	
	logger.Info("Error groups example completed successfully")
	fmt.Println("=== Example completed successfully ===")
} 