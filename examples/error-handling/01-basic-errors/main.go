/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * Basic errors example demonstrating fundamental error creation and formatting.
 */

package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

// UserService 用户服务示例
// (UserService demonstrates error handling in service layer)
type UserService struct {
	users map[string]*User
}

// User 用户结构体
// (User represents a user entity)
type User struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Age      int       `json:"age"`
	Created  time.Time `json:"created"`
	Active   bool      `json:"active"`
}

// NewUserService 创建用户服务
// (NewUserService creates a new user service)
func NewUserService() *UserService {
	return &UserService{
		users: make(map[string]*User),
	}
}

// CreateUser 创建用户（演示基础错误创建）
// (CreateUser creates a user - demonstrates basic error creation)
func (us *UserService) CreateUser(id, name, email string, age int) (*User, error) {
	// 验证输入参数 (Validate input parameters)
	if id == "" {
		return nil, errors.New("user ID cannot be empty")
	}
	
	if name == "" {
		return nil, errors.New("user name cannot be empty")
	}
	
	if email == "" {
		return nil, errors.Errorf("user email is required for user %s", id)
	}
	
	if age < 0 || age > 150 {
		return nil, errors.Errorf("invalid age %d for user %s: must be between 0 and 150", age, id)
	}
	
	// 检查用户是否已存在 (Check if user already exists)
	if _, exists := us.users[id]; exists {
		return nil, errors.Errorf("user with ID %s already exists", id)
	}
	
	// 创建用户 (Create user)
	user := &User{
		ID:      id,
		Name:    name,
		Email:   email,
		Age:     age,
		Created: time.Now(),
		Active:  true,
	}
	
	us.users[id] = user
	return user, nil
}

// GetUser 获取用户（演示不同类型的错误）
// (GetUser retrieves a user - demonstrates different types of errors)
func (us *UserService) GetUser(id string) (*User, error) {
	if id == "" {
		return nil, errors.New("user ID cannot be empty")
	}
	
	user, exists := us.users[id]
	if !exists {
		return nil, errors.Errorf("user not found: %s", id)
	}
	
	if !user.Active {
		return nil, errors.Errorf("user %s is deactivated", id)
	}
	
	return user, nil
}

// ValidateUserData 验证用户数据（演示多重验证错误）
// (ValidateUserData validates user data - demonstrates multiple validation errors)
func ValidateUserData(data map[string]string) error {
	var validationErrors []string
	
	// 检查必填字段 (Check required fields)
	requiredFields := []string{"id", "name", "email"}
	for _, field := range requiredFields {
		if value, exists := data[field]; !exists || value == "" {
			validationErrors = append(validationErrors, fmt.Sprintf("field '%s' is required", field))
		}
	}
	
	// 检查年龄格式 (Check age format)
	if ageStr, exists := data["age"]; exists {
		if age, err := strconv.Atoi(ageStr); err != nil {
			validationErrors = append(validationErrors, 
				fmt.Sprintf("invalid age format '%s': must be a number", ageStr))
		} else if age < 0 || age > 150 {
			validationErrors = append(validationErrors, 
				fmt.Sprintf("invalid age %d: must be between 0 and 150", age))
		}
	}
	
	// 检查邮箱格式（简单验证）(Check email format - simple validation)
	if email, exists := data["email"]; exists && email != "" {
		if len(email) < 5 || !contains(email, "@") || !contains(email, ".") {
			validationErrors = append(validationErrors, 
				fmt.Sprintf("invalid email format: %s", email))
		}
	}
	
	// 如果有验证错误，返回聚合错误 (If validation errors exist, return aggregated error)
	if len(validationErrors) > 0 {
		errorMsg := "validation failed: " + joinStrings(validationErrors, "; ")
		return errors.New(errorMsg)
	}
	
	return nil
}

// contains 检查字符串是否包含子字符串
// (contains checks if string contains substring)
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// joinStrings 连接字符串数组
// (joinStrings joins string array)
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}
	
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}

// simulateFileOperation 模拟文件操作（演示系统错误处理）
// (simulateFileOperation simulates file operations - demonstrates system error handling)
func simulateFileOperation(filename string) error {
	if filename == "" {
		return errors.New("filename cannot be empty")
	}
	
	// 模拟不同类型的文件操作错误 (Simulate different types of file operation errors)
	switch filename {
	case "readonly.txt":
		return errors.New("permission denied: file is read-only")
	case "missing.txt":
		return errors.New("file not found: missing.txt")
	case "corrupt.txt":
		return errors.New("file corrupted: unable to read data")
	case "locked.txt":
		return errors.New("file is locked by another process")
	case "toolarge.txt":
		return errors.Errorf("file too large: %s exceeds maximum size limit", filename)
	default:
		// 模拟成功情况 (Simulate success case)
		return nil
	}
}

// demonstrateBasicErrors 演示基础错误创建和格式化
// (demonstrateBasicErrors demonstrates basic error creation and formatting)
func demonstrateBasicErrors() {
	fmt.Println("=== Demonstrating Basic Error Creation ===")
	fmt.Println()
	
	// 1. 简单错误创建 (Simple error creation)
	fmt.Println("1. Simple Error Creation:")
	err1 := errors.New("this is a simple error")
	fmt.Printf("   Simple error: %v\n", err1)
	fmt.Printf("   Error string: %s\n", err1.Error())
	fmt.Println()
	
	// 2. 格式化错误创建 (Formatted error creation)
	fmt.Println("2. Formatted Error Creation:")
	userID := "user123"
	operation := "delete"
	err2 := errors.Errorf("failed to %s user %s: insufficient permissions", operation, userID)
	fmt.Printf("   Formatted error: %v\n", err2)
	fmt.Println()
	
	// 3. 错误格式化展示 (Error formatting demonstration)
	fmt.Println("3. Error Formatting:")
	err3 := errors.Errorf("database connection failed: timeout after %d seconds", 30)
	fmt.Printf("   %%v format: %v\n", err3)
	fmt.Printf("   %%s format: %s\n", err3)
	fmt.Printf("   %%q format: %q\n", err3)
	fmt.Printf("   %%+v format: %+v\n", err3)
	fmt.Println()
}

// demonstrateUserServiceErrors 演示用户服务中的错误处理
// (demonstrateUserServiceErrors demonstrates error handling in user service)
func demonstrateUserServiceErrors() {
	fmt.Println("=== Demonstrating User Service Errors ===")
	fmt.Println()
	
	service := NewUserService()
	
	// 1. 成功创建用户 (Successfully create user)
	fmt.Println("1. Creating valid user:")
	user, err := service.CreateUser("user001", "John Doe", "john@example.com", 25)
	if err != nil {
		fmt.Printf("   ✗ Error: %v\n", err)
	} else {
		fmt.Printf("   ✓ User created: %s (%s)\n", user.Name, user.Email)
	}
	fmt.Println()
	
	// 2. 尝试创建无效用户 (Try to create invalid users)
	fmt.Println("2. Creating users with validation errors:")
	
	invalidUsers := []struct {
		id    string
		name  string
		email string
		age   int
		desc  string
	}{
		{"", "Jane", "jane@example.com", 30, "Empty ID"},
		{"user002", "", "jane@example.com", 30, "Empty name"},
		{"user003", "Bob", "", 25, "Empty email"},
		{"user004", "Alice", "alice@example.com", -5, "Negative age"},
		{"user005", "Charlie", "charlie@example.com", 200, "Age too high"},
		{"user001", "Duplicate", "dup@example.com", 35, "Duplicate ID"},
	}
	
	for _, u := range invalidUsers {
		fmt.Printf("   Testing %s: ", u.desc)
		_, err := service.CreateUser(u.id, u.name, u.email, u.age)
		if err != nil {
			fmt.Printf("✗ %v\n", err)
		} else {
			fmt.Printf("✓ Unexpected success\n")
		}
	}
	fmt.Println()
	
	// 3. 获取用户错误 (Get user errors)
	fmt.Println("3. Getting users with various scenarios:")
	
	getUserTests := []struct {
		id   string
		desc string
	}{
		{"user001", "Valid existing user"},
		{"", "Empty ID"},
		{"nonexistent", "Non-existent user"},
	}
	
	for _, test := range getUserTests {
		fmt.Printf("   Testing %s: ", test.desc)
		user, err := service.GetUser(test.id)
		if err != nil {
			fmt.Printf("✗ %v\n", err)
		} else {
			fmt.Printf("✓ Found user: %s\n", user.Name)
		}
	}
	fmt.Println()
}

// demonstrateValidationErrors 演示验证错误
// (demonstrateValidationErrors demonstrates validation errors)
func demonstrateValidationErrors() {
	fmt.Println("=== Demonstrating Validation Errors ===")
	fmt.Println()
	
	testCases := []struct {
		name string
		data map[string]string
	}{
		{
			name: "Valid data",
			data: map[string]string{
				"id":    "user123",
				"name":  "Test User",
				"email": "test@example.com",
				"age":   "25",
			},
		},
		{
			name: "Missing required fields",
			data: map[string]string{
				"name": "Incomplete User",
			},
		},
		{
			name: "Invalid age format",
			data: map[string]string{
				"id":    "user456",
				"name":  "Invalid Age User",
				"email": "invalid@example.com",
				"age":   "not-a-number",
			},
		},
		{
			name: "Multiple validation errors",
			data: map[string]string{
				"id":    "",
				"email": "invalid-email",
				"age":   "300",
			},
		},
	}
	
	for _, tc := range testCases {
		fmt.Printf("Testing %s:\n", tc.name)
		err := ValidateUserData(tc.data)
		if err != nil {
			fmt.Printf("   ✗ Validation failed: %v\n", err)
		} else {
			fmt.Printf("   ✓ Validation passed\n")
		}
		fmt.Println()
	}
}

// demonstrateFileOperationErrors 演示文件操作错误
// (demonstrateFileOperationErrors demonstrates file operation errors)
func demonstrateFileOperationErrors() {
	fmt.Println("=== Demonstrating File Operation Errors ===")
	fmt.Println()
	
	testFiles := []string{
		"normal.txt",
		"readonly.txt",
		"missing.txt",
		"corrupt.txt",
		"locked.txt",
		"toolarge.txt",
		"",
	}
	
	for _, filename := range testFiles {
		desc := filename
		if filename == "" {
			desc = "<empty filename>"
		}
		
		fmt.Printf("Processing %s: ", desc)
		err := simulateFileOperation(filename)
		if err != nil {
			fmt.Printf("✗ %v\n", err)
		} else {
			fmt.Printf("✓ Success\n")
		}
	}
	fmt.Println()
}

// demonstrateErrorFormatting 演示错误格式化的不同方式
// (demonstrateErrorFormatting demonstrates different ways of error formatting)
func demonstrateErrorFormatting() {
	fmt.Println("=== Demonstrating Error Formatting ===")
	fmt.Println()
	
	// 创建一个复杂的错误 (Create a complex error)
	err := errors.Errorf("database query failed: table 'users' not found in schema 'production'")
	
	fmt.Println("Error formatting examples:")
	fmt.Printf("%%v (value):           %v\n", err)
	fmt.Printf("%%s (string):          %s\n", err)
	fmt.Printf("%%q (quoted string):   %q\n", err)
	fmt.Printf("%%+v (verbose):        %+v\n", err)
	fmt.Println()
	
	// 演示带有更多上下文的错误 (Demonstrate error with more context)
	contextErr := errors.Errorf("user operation failed: %v", err)
	fmt.Println("Error with context:")
	fmt.Printf("%%+v format:\n%+v\n", contextErr)
	fmt.Println()
}

func main() {
	fmt.Println("=== Basic Errors Example ===")
	fmt.Println("This example demonstrates fundamental error creation and formatting patterns.")
	fmt.Println()
	
	// 1. 初始化日志 (Initialize logging)
	logOpts := log.NewOptions()
	logOpts.Level = "info"
	logOpts.Format = "text"
	logOpts.EnableColor = true
	log.Init(logOpts)
	logger := log.Std()
	
	// 2. 演示基础错误创建 (Demonstrate basic error creation)
	demonstrateBasicErrors()
	
	// 3. 演示用户服务错误 (Demonstrate user service errors)
	demonstrateUserServiceErrors()
	
	// 4. 演示验证错误 (Demonstrate validation errors)
	demonstrateValidationErrors()
	
	// 5. 演示文件操作错误 (Demonstrate file operation errors)
	demonstrateFileOperationErrors()
	
	// 6. 演示错误格式化 (Demonstrate error formatting)
	demonstrateErrorFormatting()
	
	logger.Info("Basic errors example completed successfully")
	fmt.Println("=== Example completed successfully ===")
} 