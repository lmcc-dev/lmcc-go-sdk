/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * Error wrapping example demonstrating context addition and error chain navigation.
 */

package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

// DatabaseError 模拟数据库错误
// (DatabaseError simulates database errors)
type DatabaseError struct {
	Query    string
	Args     []interface{}
	Duration time.Duration
	Err      error
}

func (e DatabaseError) Error() string {
	return fmt.Sprintf("database query failed: %v", e.Err)
}

func (e DatabaseError) Unwrap() error {
	return e.Err
}

// UserRepository 用户仓库接口
// (UserRepository represents user repository interface)
type UserRepository interface {
	GetUserByID(id string) (*User, error)
	CreateUser(user *User) error
	UpdateUser(user *User) error
	DeleteUser(id string) error
}

// User 用户实体
// (User represents user entity)
type User struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Age      int       `json:"age"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	Active   bool      `json:"active"`
}

// MockUserRepository 模拟用户仓库实现
// (MockUserRepository simulates user repository implementation)
type MockUserRepository struct {
	users   map[string]*User
	latency time.Duration
}

// NewMockUserRepository 创建模拟用户仓库
// (NewMockUserRepository creates a mock user repository)
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:   make(map[string]*User),
		latency: 50 * time.Millisecond,
	}
}

// GetUserByID 通过ID获取用户（演示错误包装）
// (GetUserByID retrieves user by ID - demonstrates error wrapping)
func (r *MockUserRepository) GetUserByID(id string) (*User, error) {
	// 模拟网络延迟 (Simulate network latency)
	time.Sleep(r.latency)
	
	// 底层错误：输入验证 (Low-level error: input validation)
	if id == "" {
		baseErr := errors.New("user ID cannot be empty")
		return nil, errors.Wrap(baseErr, "input validation failed")
	}
	
	// 底层错误：格式验证 (Low-level error: format validation)
	if len(id) < 3 {
		baseErr := errors.New("user ID too short")
		return nil, errors.Wrapf(baseErr, "invalid ID format: %s", id)
	}
	
	// 模拟数据库连接错误 (Simulate database connection error)
	if id == "conn_error" {
		dbErr := DatabaseError{
			Query:    "SELECT * FROM users WHERE id = ?",
			Args:     []interface{}{id},
			Duration: r.latency,
			Err:      sql.ErrConnDone,
		}
		return nil, errors.Wrap(dbErr, "failed to execute database query")
	}
	
	// 模拟数据库超时 (Simulate database timeout)
	if id == "timeout" {
		baseErr := errors.New("context deadline exceeded")
		wrappedErr := errors.Wrap(baseErr, "database operation timed out")
		return nil, errors.Wrapf(wrappedErr, "failed to get user %s", id)
	}
	
	// 查找用户 (Find user)
	user, exists := r.users[id]
	if !exists {
		baseErr := errors.New("record not found")
		return nil, errors.Wrapf(baseErr, "user not found: %s", id)
	}
	
	// 检查用户状态 (Check user status)
	if !user.Active {
		baseErr := errors.New("user account deactivated")
		return nil, errors.Wrapf(baseErr, "cannot access deactivated user: %s", id)
	}
	
	return user, nil
}

// CreateUser 创建用户（演示多层错误包装）
// (CreateUser creates a user - demonstrates multi-layer error wrapping)
func (r *MockUserRepository) CreateUser(user *User) error {
	// 第一层：输入验证错误 (First layer: input validation errors)
	if err := r.validateUser(user); err != nil {
		return errors.Wrap(err, "user validation failed")
	}
	
	// 第二层：业务逻辑错误 (Second layer: business logic errors)
	if err := r.checkBusinessRules(user); err != nil {
		return errors.Wrapf(err, "business rule validation failed for user %s", user.ID)
	}
	
	// 第三层：数据库操作错误 (Third layer: database operation errors)
	if err := r.persistUser(user); err != nil {
		return errors.Wrapf(err, "failed to persist user %s to database", user.ID)
	}
	
	return nil
}

// validateUser 验证用户数据
// (validateUser validates user data)
func (r *MockUserRepository) validateUser(user *User) error {
	var validationErrors []string
	
	if user.ID == "" {
		validationErrors = append(validationErrors, "ID is required")
	}
	
	if user.Name == "" {
		validationErrors = append(validationErrors, "name is required")
	}
	
	if user.Email == "" {
		validationErrors = append(validationErrors, "email is required")
	}
	
	if user.Age < 0 || user.Age > 150 {
		validationErrors = append(validationErrors, fmt.Sprintf("invalid age: %d", user.Age))
	}
	
	if len(validationErrors) > 0 {
		baseErr := errors.New(strings.Join(validationErrors, "; "))
		return errors.Wrap(baseErr, "field validation errors")
	}
	
	return nil
}

// checkBusinessRules 检查业务规则
// (checkBusinessRules checks business rules)
func (r *MockUserRepository) checkBusinessRules(user *User) error {
	// 检查用户是否已存在 (Check if user already exists)
	if _, exists := r.users[user.ID]; exists {
		baseErr := errors.New("duplicate user ID")
		return errors.Wrapf(baseErr, "user %s already exists", user.ID)
	}
	
	// 检查邮箱是否已被使用 (Check if email is already used)
	for _, existingUser := range r.users {
		if existingUser.Email == user.Email {
			baseErr := errors.New("duplicate email address")
			return errors.Wrapf(baseErr, "email %s is already in use by user %s", 
				user.Email, existingUser.ID)
		}
	}
	
	// 模拟外部验证服务错误 (Simulate external validation service error)
	if user.Email == "invalid@blacklist.com" {
		baseErr := errors.New("email domain is blacklisted")
		externalErr := errors.Wrap(baseErr, "external email validation failed")
		return errors.Wrapf(externalErr, "email validation service rejected %s", user.Email)
	}
	
	return nil
}

// persistUser 持久化用户数据
// (persistUser persists user data)
func (r *MockUserRepository) persistUser(user *User) error {
	// 模拟不同的数据库错误 (Simulate different database errors)
	switch user.ID {
	case "db_error":
		dbErr := errors.New("table 'users' doesn't exist")
		return errors.Wrap(dbErr, "database schema error")
		
	case "constraint_error":
		dbErr := errors.New("foreign key constraint failed")
		constraintErr := errors.Wrap(dbErr, "referential integrity violation")
		return errors.Wrapf(constraintErr, "failed to insert user %s", user.ID)
		
	case "disk_full":
		dbErr := errors.New("no space left on device")
		storageErr := errors.Wrap(dbErr, "storage subsystem error")
		return errors.Wrapf(storageErr, "unable to write user %s to disk", user.ID)
	}
	
	// 设置时间戳 (Set timestamps)
	now := time.Now()
	user.Created = now
	user.Modified = now
	user.Active = true
	
	// 保存用户 (Save user)
	r.users[user.ID] = user
	return nil
}

// UpdateUser 更新用户（演示条件错误包装）
// (UpdateUser updates a user - demonstrates conditional error wrapping)
func (r *MockUserRepository) UpdateUser(user *User) error {
	// 首先检查用户是否存在 (First check if user exists)
	existing, err := r.GetUserByID(user.ID)
	if err != nil {
		// 包装来自GetUserByID的错误 (Wrap error from GetUserByID)
		return errors.Wrapf(err, "cannot update user %s: prerequisite check failed", user.ID)
	}
	
	// 检查用户是否可以更新 (Check if user can be updated)
	if !existing.Active {
		baseErr := errors.New("user account is deactivated")
		return errors.Wrapf(baseErr, "cannot update deactivated user %s", user.ID)
	}
	
	// 模拟并发修改检查 (Simulate concurrent modification check)
	if user.Modified.Before(existing.Modified) {
		baseErr := errors.New("concurrent modification detected")
		concurrencyErr := errors.Wrap(baseErr, "optimistic locking failed")
		return errors.Wrapf(concurrencyErr, "user %s was modified by another process", user.ID)
	}
	
	// 更新用户 (Update user)
	user.Modified = time.Now()
	r.users[user.ID] = user
	return nil
}

// DeleteUser 删除用户
// (DeleteUser deletes a user)
func (r *MockUserRepository) DeleteUser(id string) error {
	// 检查用户是否存在 (Check if user exists)
	user, err := r.GetUserByID(id)
	if err != nil {
		return errors.Wrapf(err, "cannot delete user %s: user lookup failed", id)
	}
	
	// 检查是否有依赖关系 (Check for dependencies)
	if err := r.checkDependencies(id); err != nil {
		return errors.Wrapf(err, "cannot delete user %s: dependency check failed", id)
	}
	
	// 软删除：标记为非活跃 (Soft delete: mark as inactive)
	user.Active = false
	user.Modified = time.Now()
	r.users[id] = user
	
	return nil
}

// checkDependencies 检查用户依赖关系
// (checkDependencies checks user dependencies)
func (r *MockUserRepository) checkDependencies(userID string) error {
	// 模拟依赖检查错误 (Simulate dependency check errors)
	if userID == "has_orders" {
		baseErr := errors.New("user has active orders")
		return errors.Wrap(baseErr, "dependency constraint violation")
	}
	
	if userID == "admin_user" {
		baseErr := errors.New("cannot delete admin user")
		return errors.Wrap(baseErr, "permission denied")
	}
	
	return nil
}

// UserService 用户服务（演示服务层错误包装）
// (UserService represents user service - demonstrates service layer error wrapping)
type UserService struct {
	repo UserRepository
}

// NewUserService 创建用户服务
// (NewUserService creates a user service)
func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

// RegisterUser 注册用户（演示服务层错误包装）
// (RegisterUser registers a user - demonstrates service layer error wrapping)
func (s *UserService) RegisterUser(id, name, email string, age int) (*User, error) {
	// 创建用户对象 (Create user object)
	user := &User{
		ID:    id,
		Name:  name,
		Email: email,
		Age:   age,
	}
	
	// 尝试创建用户 (Try to create user)
	if err := s.repo.CreateUser(user); err != nil {
		return nil, errors.Wrapf(err, "user registration failed for %s", id)
	}
	
	return user, nil
}

// GetUserProfile 获取用户资料（演示错误链分析）
// (GetUserProfile gets user profile - demonstrates error chain analysis)
func (s *UserService) GetUserProfile(id string) (*User, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		// 包装错误并添加服务层上下文 (Wrap error and add service layer context)
		return nil, errors.Wrapf(err, "failed to retrieve user profile for %s", id)
	}
	
	return user, nil
}

// analyzeErrorChain 分析错误链
// (analyzeErrorChain analyzes error chain)
func analyzeErrorChain(err error) {
	fmt.Println("=== Error Chain Analysis ===")
	
	if err == nil {
		fmt.Println("No error to analyze")
		return
	}
	
	fmt.Printf("Surface error: %v\n", err)
	fmt.Printf("Error type: %T\n", err)
	fmt.Println()
	
	// 逐层展开错误 (Unwrap errors layer by layer)
	current := err
	depth := 0
	
	for current != nil {
		fmt.Printf("Level %d: %v\n", depth, current)
		fmt.Printf("  Type: %T\n", current)
		
		// 检查是否有特殊方法 (Check for special methods)
		if coder := errors.GetCoder(current); coder != nil {
			fmt.Printf("  Code: %d (%s)\n", coder.Code(), coder.String())
		}
		
		// 尝试展开 (Try to unwrap)
		if unwrapper, ok := current.(interface{ Unwrap() error }); ok {
			current = unwrapper.Unwrap()
		} else {
			current = nil
		}
		depth++
		fmt.Println()
	}
	
	// 获取根本原因 (Get root cause)
	rootCause := errors.Cause(err)
	if rootCause != err {
		fmt.Printf("Root cause: %v\n", rootCause)
		fmt.Printf("Root cause type: %T\n", rootCause)
	}
	
	fmt.Println()
}

// demonstrateBasicWrapping 演示基础错误包装
// (demonstrateBasicWrapping demonstrates basic error wrapping)
func demonstrateBasicWrapping() {
	fmt.Println("=== Demonstrating Basic Error Wrapping ===")
	fmt.Println()
	
	// 1. 简单包装 (Simple wrapping)
	fmt.Println("1. Simple Error Wrapping:")
	baseErr := errors.New("file not found")
	wrappedErr := errors.Wrap(baseErr, "failed to read configuration")
	fmt.Printf("   Base error: %v\n", baseErr)
	fmt.Printf("   Wrapped error: %v\n", wrappedErr)
	fmt.Printf("   Root cause: %v\n", errors.Cause(wrappedErr))
	fmt.Println()
	
	// 2. 格式化包装 (Formatted wrapping)
	fmt.Println("2. Formatted Error Wrapping:")
	filename := "config.yaml"
	formattedErr := errors.Wrapf(baseErr, "failed to read file %s", filename)
	fmt.Printf("   Formatted wrapped: %v\n", formattedErr)
	fmt.Println()
	
	// 3. 多层包装 (Multi-layer wrapping)
	fmt.Println("3. Multi-layer Error Wrapping:")
	layer1 := errors.Wrap(baseErr, "file system error")
	layer2 := errors.Wrap(layer1, "configuration loading failed")
	layer3 := errors.Wrap(layer2, "application initialization failed")
	
	fmt.Printf("   Layer 3: %v\n", layer3)
	fmt.Printf("   Layer 2: %v\n", layer2)
	fmt.Printf("   Layer 1: %v\n", layer1)
	fmt.Printf("   Base: %v\n", baseErr)
	fmt.Println()
	
	analyzeErrorChain(layer3)
}

// demonstrateRepositoryErrors 演示仓库层错误
// (demonstrateRepositoryErrors demonstrates repository layer errors)
func demonstrateRepositoryErrors() {
	fmt.Println("=== Demonstrating Repository Error Wrapping ===")
	fmt.Println()
	
	repo := NewMockUserRepository()
	
	// 1. 测试获取用户错误 (Test get user errors)
	fmt.Println("1. GetUser Error Scenarios:")
	
	getUserTests := []struct {
		id   string
		desc string
	}{
		{"", "Empty ID"},
		{"xx", "Short ID"},
		{"conn_error", "Database connection error"},
		{"timeout", "Database timeout"},
		{"nonexistent", "Non-existent user"},
	}
	
	for _, test := range getUserTests {
		fmt.Printf("\nTesting %s (ID: %s):\n", test.desc, test.id)
		_, err := repo.GetUserByID(test.id)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			analyzeErrorChain(err)
		}
	}
}

// demonstrateServiceErrors 演示服务层错误
// (demonstrateServiceErrors demonstrates service layer errors)
func demonstrateServiceErrors() {
	fmt.Println("=== Demonstrating Service Layer Error Wrapping ===")
	fmt.Println()
	
	repo := NewMockUserRepository()
	service := NewUserService(repo)
	
	// 1. 测试用户注册错误 (Test user registration errors)
	fmt.Println("1. User Registration Error Scenarios:")
	
	registrationTests := []struct {
		id    string
		name  string
		email string
		age   int
		desc  string
	}{
		{"", "John Doe", "john@example.com", 25, "Empty ID"},
		{"user1", "", "john@example.com", 25, "Empty name"},
		{"user2", "Jane Doe", "", 30, "Empty email"},
		{"user3", "Bob Smith", "bob@example.com", -5, "Invalid age"},
		{"db_error", "Test User", "test@example.com", 25, "Database error"},
		{"constraint_error", "Constraint User", "constraint@example.com", 30, "Constraint error"},
	}
	
	for _, test := range registrationTests {
		fmt.Printf("\nTesting %s:\n", test.desc)
		_, err := service.RegisterUser(test.id, test.name, test.email, test.age)
		if err != nil {
			fmt.Printf("Registration failed: %v\n", err)
			fmt.Printf("Detailed error chain:\n")
			analyzeErrorChain(err)
		} else {
			fmt.Printf("✓ User registered successfully\n")
		}
	}
}

// demonstrateErrorFormatting 演示错误格式化
// (demonstrateErrorFormatting demonstrates error formatting)
func demonstrateErrorFormatting() {
	fmt.Println("=== Demonstrating Error Formatting with Wrapping ===")
	fmt.Println()
	
	// 创建复杂的包装错误 (Create complex wrapped error)
	baseErr := errors.New("connection refused")
	dbErr := errors.Wrap(baseErr, "database connection failed")
	repoErr := errors.Wrapf(dbErr, "user repository error: failed to query user %s", "user123")
	serviceErr := errors.Wrap(repoErr, "user service operation failed")
	
	fmt.Println("Error formatting with different verbs:")
	fmt.Printf("%%v:  %v\n", serviceErr)
	fmt.Printf("%%s:  %s\n", serviceErr)
	fmt.Printf("%%q:  %q\n", serviceErr)
	fmt.Printf("%%+v:\n%+v\n", serviceErr)
	fmt.Println()
}

func main() {
	fmt.Println("=== Error Wrapping Example ===")
	fmt.Println("This example demonstrates error wrapping patterns and error chain navigation.")
	fmt.Println()
	
	// 1. 初始化日志 (Initialize logging)
	logOpts := log.NewOptions()
	logOpts.Level = "info"
	logOpts.Format = "text"
	logOpts.EnableColor = true
	log.Init(logOpts)
	logger := log.Std()
	
	// 2. 演示基础错误包装 (Demonstrate basic error wrapping)
	demonstrateBasicWrapping()
	
	// 3. 演示仓库层错误 (Demonstrate repository layer errors)
	demonstrateRepositoryErrors()
	
	// 4. 演示服务层错误 (Demonstrate service layer errors)
	demonstrateServiceErrors()
	
	// 5. 演示错误格式化 (Demonstrate error formatting)
	demonstrateErrorFormatting()
	
	logger.Info("Error wrapping example completed successfully")
	fmt.Println("=== Example completed successfully ===")
} 