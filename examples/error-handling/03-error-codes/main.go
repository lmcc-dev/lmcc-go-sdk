/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * Error codes example demonstrating structured error handling with codes.
 */

package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

// 自定义错误码 (Custom error codes)
var (
	// 业务逻辑错误码 (Business logic error codes)
	ErrUserNotFound      = errors.NewCoder(1001, 404, "UserNotFound", "User not found")
	ErrUserAlreadyExists = errors.NewCoder(1002, 409, "UserAlreadyExists", "User already exists")
	ErrInvalidUserData   = errors.NewCoder(1003, 400, "InvalidUserData", "Invalid user data")
	ErrUserDeactivated   = errors.NewCoder(1004, 403, "UserDeactivated", "User account is deactivated")
	ErrInsufficientPermissions = errors.NewCoder(1005, 403, "InsufficientPermissions", "Insufficient permissions")
	
	// 数据访问错误码 (Data access error codes)
	ErrDatabaseConnection = errors.NewCoder(2001, 503, "DatabaseConnection", "Database connection failed")
	ErrDatabaseTimeout    = errors.NewCoder(2002, 504, "DatabaseTimeout", "Database operation timeout")
	ErrDatabaseConstraint = errors.NewCoder(2003, 409, "DatabaseConstraint", "Database constraint violation")
	ErrDatabaseSchema     = errors.NewCoder(2004, 500, "DatabaseSchema", "Database schema error")
	
	// 外部服务错误码 (External service error codes)
	ErrExternalService    = errors.NewCoder(3001, 502, "ExternalService", "External service error")
	ErrServiceUnavailable = errors.NewCoder(3002, 503, "ServiceUnavailable", "Service temporarily unavailable")
	ErrRateLimitExceeded  = errors.NewCoder(3003, 429, "RateLimitExceeded", "Rate limit exceeded")
	ErrAPIQuotaExceeded   = errors.NewCoder(3004, 402, "APIQuotaExceeded", "API quota exceeded")
	
	// 系统错误码 (System error codes)
	ErrInternalServer     = errors.NewCoder(5001, 500, "InternalServer", "Internal server error")
	ErrConfigurationError = errors.NewCoder(5002, 500, "ConfigurationError", "Configuration error")
	ErrResourceExhausted  = errors.NewCoder(5003, 503, "ResourceExhausted", "Resource exhausted")
	ErrMaintenanceMode    = errors.NewCoder(5004, 503, "MaintenanceMode", "System in maintenance mode")
)

// User 用户实体
// (User represents user entity)
type User struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Age      int       `json:"age"`
	Status   string    `json:"status"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

// UserRepository 用户仓库接口
// (UserRepository represents user repository interface)
type UserRepository struct {
	users       map[string]*User
	maintenance bool
}

// NewUserRepository 创建用户仓库
// (NewUserRepository creates a user repository)
func NewUserRepository() *UserRepository {
	return &UserRepository{
		users:       make(map[string]*User),
		maintenance: false,
	}
}

// GetUser 根据ID获取用户（演示错误码使用）
// (GetUser retrieves user by ID - demonstrates error code usage)
func (r *UserRepository) GetUser(id string) (*User, error) {
	// 检查维护模式 (Check maintenance mode)
	if r.maintenance {
		return nil, errors.WithCode(
			errors.New("system is under maintenance"),
			ErrMaintenanceMode,
		)
	}
	
	// 输入验证 (Input validation)
	if id == "" {
		return nil, errors.WithCode(
			errors.New("user ID cannot be empty"),
			ErrInvalidUserData,
		)
	}
	
	// 模拟数据库连接错误 (Simulate database connection error)
	if id == "db_error" {
		return nil, errors.WithCode(
			errors.New("failed to connect to database"),
			ErrDatabaseConnection,
		)
	}
	
	// 模拟数据库超时 (Simulate database timeout)
	if id == "timeout" {
		return nil, errors.WithCode(
			errors.New("database query timed out after 30s"),
			ErrDatabaseTimeout,
		)
	}
	
	// 查找用户 (Find user)
	user, exists := r.users[id]
	if !exists {
		return nil, errors.WithCode(
			errors.Errorf("user with ID %s does not exist", id),
			ErrUserNotFound,
		)
	}
	
	// 检查用户状态 (Check user status)
	if user.Status == "deactivated" {
		return nil, errors.WithCode(
			errors.Errorf("user %s is deactivated", id),
			ErrUserDeactivated,
		)
	}
	
	return user, nil
}

// CreateUser 创建用户（演示错误码和包装）
// (CreateUser creates a user - demonstrates error codes and wrapping)
func (r *UserRepository) CreateUser(user *User) error {
	// 检查维护模式 (Check maintenance mode)
	if r.maintenance {
		return errors.WithCode(
			errors.New("cannot create user during maintenance"),
			ErrMaintenanceMode,
		)
	}
	
	// 验证用户数据 (Validate user data)
	if err := r.validateUser(user); err != nil {
		return errors.WithCode(err, ErrInvalidUserData)
	}
	
	// 检查用户是否已存在 (Check if user already exists)
	if _, exists := r.users[user.ID]; exists {
		return errors.WithCode(
			errors.Errorf("user with ID %s already exists", user.ID),
			ErrUserAlreadyExists,
		)
	}
	
	// 模拟数据库约束错误 (Simulate database constraint error)
	if user.Email == "constraint@error.com" {
		return errors.WithCode(
			errors.New("email violates unique constraint"),
			ErrDatabaseConstraint,
		)
	}
	
	// 模拟数据库架构错误 (Simulate database schema error)
	if user.ID == "schema_error" {
		return errors.WithCode(
			errors.New("table 'users' does not exist"),
			ErrDatabaseSchema,
		)
	}
	
	// 设置默认值 (Set default values)
	user.Status = "active"
	user.Created = time.Now()
	user.Modified = time.Now()
	
	// 保存用户 (Save user)
	r.users[user.ID] = user
	return nil
}

// validateUser 验证用户数据
// (validateUser validates user data)
func (r *UserRepository) validateUser(user *User) error {
	if user.ID == "" {
		return errors.New("user ID is required")
	}
	
	if user.Name == "" {
		return errors.New("user name is required")
	}
	
	if user.Email == "" {
		return errors.New("user email is required")
	}
	
	if user.Age < 0 || user.Age > 150 {
		return errors.Errorf("invalid age: %d", user.Age)
	}
	
	return nil
}

// UpdateUser 更新用户
// (UpdateUser updates a user)
func (r *UserRepository) UpdateUser(user *User) error {
	// 检查用户是否存在 (Check if user exists)
	existing, err := r.GetUser(user.ID)
	if err != nil {
		return err // 错误码已经设置 (Error code already set)
	}
	
	// 验证数据 (Validate data)
	if err := r.validateUser(user); err != nil {
		return errors.WithCode(err, ErrInvalidUserData)
	}
	
	// 更新时间戳 (Update timestamp)
	user.Status = existing.Status
	user.Created = existing.Created
	user.Modified = time.Now()
	
	// 保存更新 (Save update)
	r.users[user.ID] = user
	return nil
}

// DeleteUser 删除用户
// (DeleteUser deletes a user)
func (r *UserRepository) DeleteUser(id string) error {
	// 检查用户是否存在 (Check if user exists)
	user, err := r.GetUser(id)
	if err != nil {
		return err // 错误码已经设置 (Error code already set)
	}
	
	// 检查权限 (Check permissions)
	if user.Email == "admin@example.com" {
		return errors.WithCode(
			errors.New("cannot delete admin user"),
			ErrInsufficientPermissions,
		)
	}
	
	// 软删除 (Soft delete)
	user.Status = "deleted"
	user.Modified = time.Now()
	r.users[id] = user
	
	return nil
}

// SetMaintenanceMode 设置维护模式
// (SetMaintenanceMode sets maintenance mode)
func (r *UserRepository) SetMaintenanceMode(enabled bool) {
	r.maintenance = enabled
}

// UserService 用户服务（演示服务层错误码处理）
// (UserService represents user service - demonstrates service layer error code handling)
type UserService struct {
	repo         *UserRepository
	rateLimiter  map[string]time.Time
	rateLimit    time.Duration
	quotaTracker map[string]int
	dailyQuota   int
}

// NewUserService 创建用户服务
// (NewUserService creates a user service)
func NewUserService(repo *UserRepository) *UserService {
	return &UserService{
		repo:         repo,
		rateLimiter:  make(map[string]time.Time),
		rateLimit:    time.Second,
		quotaTracker: make(map[string]int),
		dailyQuota:   100,
	}
}

// GetUserProfile 获取用户资料（演示速率限制）
// (GetUserProfile gets user profile - demonstrates rate limiting)
func (s *UserService) GetUserProfile(userID, clientID string) (*User, error) {
	// 检查速率限制 (Check rate limit)
	if err := s.checkRateLimit(clientID); err != nil {
		return nil, err
	}
	
	// 检查配额 (Check quota)
	if err := s.checkQuota(clientID); err != nil {
		return nil, err
	}
	
	// 获取用户 (Get user)
	user, err := s.repo.GetUser(userID)
	if err != nil {
		return nil, err // 错误码已经设置 (Error code already set)
	}
	
	return user, nil
}

// checkRateLimit 检查速率限制
// (checkRateLimit checks rate limit)
func (s *UserService) checkRateLimit(clientID string) error {
	lastRequest, exists := s.rateLimiter[clientID]
	if exists && time.Since(lastRequest) < s.rateLimit {
		return errors.WithCode(
			errors.Errorf("rate limit exceeded for client %s", clientID),
			ErrRateLimitExceeded,
		)
	}
	
	s.rateLimiter[clientID] = time.Now()
	return nil
}

// checkQuota 检查配额
// (checkQuota checks quota)
func (s *UserService) checkQuota(clientID string) error {
	count := s.quotaTracker[clientID]
	if count >= s.dailyQuota {
		return errors.WithCode(
			errors.Errorf("daily quota exceeded for client %s", clientID),
			ErrAPIQuotaExceeded,
		)
	}
	
	s.quotaTracker[clientID] = count + 1
	return nil
}

// CreateUserWithValidation 创建用户并进行全面验证
// (CreateUserWithValidation creates user with comprehensive validation)
func (s *UserService) CreateUserWithValidation(user *User, clientID string) error {
	// 检查速率限制 (Check rate limit)
	if err := s.checkRateLimit(clientID); err != nil {
		return err
	}
	
	// 模拟外部验证服务调用 (Simulate external validation service call)
	if err := s.validateWithExternalService(user); err != nil {
		return err
	}
	
	// 创建用户 (Create user)
	return s.repo.CreateUser(user)
}

// validateWithExternalService 模拟外部服务验证
// (validateWithExternalService simulates external service validation)
func (s *UserService) validateWithExternalService(user *User) error {
	// 模拟服务不可用 (Simulate service unavailable)
	if user.Email == "service@unavailable.com" {
		return errors.WithCode(
			errors.New("email validation service is unavailable"),
			ErrServiceUnavailable,
		)
	}
	
	// 模拟外部服务错误 (Simulate external service error)
	if user.Email == "external@error.com" {
		return errors.WithCode(
			errors.New("external validation service returned error"),
			ErrExternalService,
		)
	}
	
	return nil
}

// HTTPErrorMapper HTTP错误映射器
// (HTTPErrorMapper maps errors to HTTP status codes)
type HTTPErrorMapper struct{}

// MapErrorToHTTPStatus 将错误码映射到HTTP状态码
// (MapErrorToHTTPStatus maps error codes to HTTP status codes)
func (m *HTTPErrorMapper) MapErrorToHTTPStatus(err error) (int, string) {
	if err == nil {
		return http.StatusOK, "OK"
	}
	
	coder := errors.GetCoder(err)
	if coder == nil {
		return http.StatusInternalServerError, "Internal Server Error"
	}
	
	switch coder.Code() {
	// 业务逻辑错误 (Business logic errors)
	case ErrUserNotFound.Code():
		return http.StatusNotFound, "User Not Found"
	case ErrUserAlreadyExists.Code():
		return http.StatusConflict, "User Already Exists"
	case ErrInvalidUserData.Code():
		return http.StatusBadRequest, "Invalid User Data"
	case ErrUserDeactivated.Code():
		return http.StatusForbidden, "User Account Deactivated"
	case ErrInsufficientPermissions.Code():
		return http.StatusForbidden, "Insufficient Permissions"
		
	// 数据访问错误 (Data access errors)
	case ErrDatabaseConnection.Code(), ErrDatabaseTimeout.Code():
		return http.StatusServiceUnavailable, "Database Service Unavailable"
	case ErrDatabaseConstraint.Code():
		return http.StatusConflict, "Data Constraint Violation"
	case ErrDatabaseSchema.Code():
		return http.StatusInternalServerError, "Database Schema Error"
		
	// 外部服务错误 (External service errors)
	case ErrExternalService.Code(), ErrServiceUnavailable.Code():
		return http.StatusBadGateway, "External Service Error"
	case ErrRateLimitExceeded.Code():
		return http.StatusTooManyRequests, "Rate Limit Exceeded"
	case ErrAPIQuotaExceeded.Code():
		return http.StatusPaymentRequired, "API Quota Exceeded"
		
	// 系统错误 (System errors)
	case ErrInternalServer.Code():
		return http.StatusInternalServerError, "Internal Server Error"
	case ErrConfigurationError.Code():
		return http.StatusInternalServerError, "Configuration Error"
	case ErrResourceExhausted.Code():
		return http.StatusServiceUnavailable, "Resource Exhausted"
	case ErrMaintenanceMode.Code():
		return http.StatusServiceUnavailable, "Service Under Maintenance"
		
	default:
		return http.StatusInternalServerError, "Unknown Error"
	}
}

// analyzeError 分析错误码信息
// (analyzeError analyzes error code information)
func analyzeError(err error) {
	fmt.Println("=== Error Analysis ===")
	
	if err == nil {
		fmt.Println("No error to analyze")
		return
	}
	
	fmt.Printf("Error: %v\n", err)
	fmt.Printf("Error Type: %T\n", err)
	
	// 检查错误码 (Check error code)
	coder := errors.GetCoder(err)
	if coder != nil {
		fmt.Printf("Error Code: %d\n", coder.Code())
		fmt.Printf("Error String: %s\n", coder.String())
		fmt.Printf("Error Message: %s\n", coder.Error())
		
		// HTTP映射 (HTTP mapping)
		mapper := &HTTPErrorMapper{}
		httpStatus, httpMessage := mapper.MapErrorToHTTPStatus(err)
		fmt.Printf("HTTP Status: %d %s\n", httpStatus, httpMessage)
	} else {
		fmt.Println("No error code information available")
	}
	
	fmt.Println()
}

// demonstrateBasicErrorCodes 演示基础错误码使用
// (demonstrateBasicErrorCodes demonstrates basic error code usage)
func demonstrateBasicErrorCodes() {
	fmt.Println("=== Demonstrating Basic Error Code Usage ===")
	fmt.Println()
	
	repo := NewUserRepository()
	
	// 1. 测试各种错误场景 (Test various error scenarios)
	errorScenarios := []struct {
		id   string
		desc string
	}{
		{"", "Empty ID (validation error)"},
		{"db_error", "Database connection error"},
		{"timeout", "Database timeout error"},
		{"nonexistent", "User not found error"},
	}
	
	for _, scenario := range errorScenarios {
		fmt.Printf("Testing %s:\n", scenario.desc)
		_, err := repo.GetUser(scenario.id)
		if err != nil {
			analyzeError(err)
		}
	}
}

// demonstrateUserOperations 演示用户操作中的错误码
// (demonstrateUserOperations demonstrates error codes in user operations)
func demonstrateUserOperations() {
	fmt.Println("=== Demonstrating Error Codes in User Operations ===")
	fmt.Println()
	
	repo := NewUserRepository()
	
	// 1. 测试创建用户错误 (Test user creation errors)
	fmt.Println("1. User Creation Error Scenarios:")
	
	users := []*User{
		{ID: "", Name: "Invalid User", Email: "invalid@example.com", Age: 25},
		{ID: "user1", Name: "", Email: "noname@example.com", Age: 30},
		{ID: "user2", Name: "Valid User", Email: "", Age: 28},
		{ID: "user3", Name: "Invalid Age", Email: "invalid@example.com", Age: -5},
		{ID: "schema_error", Name: "Schema Error", Email: "schema@example.com", Age: 25},
		{ID: "valid_user", Name: "Valid User", Email: "valid@example.com", Age: 25},
	}
	
	for _, user := range users {
		fmt.Printf("\nCreating user %s:\n", user.ID)
		err := repo.CreateUser(user)
		if err != nil {
			analyzeError(err)
		} else {
			fmt.Printf("✓ User created successfully\n\n")
		}
	}
	
	// 2. 测试重复创建 (Test duplicate creation)
	fmt.Println("2. Duplicate User Creation:")
	validUser := &User{ID: "valid_user", Name: "Duplicate", Email: "dup@example.com", Age: 30}
	err := repo.CreateUser(validUser)
	if err != nil {
		analyzeError(err)
	}
}

// demonstrateServiceLayerErrors 演示服务层错误码
// (demonstrateServiceLayerErrors demonstrates service layer error codes)
func demonstrateServiceLayerErrors() {
	fmt.Println("=== Demonstrating Service Layer Error Codes ===")
	fmt.Println()
	
	repo := NewUserRepository()
	service := NewUserService(repo)
	
	// 创建一个测试用户 (Create a test user)
	testUser := &User{ID: "test_user", Name: "Test User", Email: "test@example.com", Age: 25}
	repo.CreateUser(testUser)
	
	// 1. 测试速率限制 (Test rate limiting)
	fmt.Println("1. Rate Limiting Test:")
	for i := 0; i < 3; i++ {
		fmt.Printf("Request %d:\n", i+1)
		_, err := service.GetUserProfile("test_user", "client1")
		if err != nil {
			analyzeError(err)
		} else {
			fmt.Printf("✓ Request successful\n\n")
		}
	}
	
	// 2. 测试外部服务错误 (Test external service errors)
	fmt.Println("2. External Service Error Test:")
	
	externalUsers := []*User{
		{ID: "svc_unavail", Name: "Service Unavailable", Email: "service@unavailable.com", Age: 25},
		{ID: "ext_error", Name: "External Error", Email: "external@error.com", Age: 30},
	}
	
	for _, user := range externalUsers {
		fmt.Printf("Creating user with email %s:\n", user.Email)
		err := service.CreateUserWithValidation(user, "client2")
		if err != nil {
			analyzeError(err)
		}
	}
}

// demonstrateMaintenanceMode 演示维护模式错误
// (demonstrateMaintenanceMode demonstrates maintenance mode errors)
func demonstrateMaintenanceMode() {
	fmt.Println("=== Demonstrating Maintenance Mode Errors ===")
	fmt.Println()
	
	repo := NewUserRepository()
	
	// 启用维护模式 (Enable maintenance mode)
	repo.SetMaintenanceMode(true)
	fmt.Println("Maintenance mode enabled")
	
	// 尝试各种操作 (Try various operations)
	fmt.Println("\nTrying operations in maintenance mode:")
	
	// 获取用户 (Get user)
	fmt.Println("1. Get user:")
	_, err := repo.GetUser("any_user")
	if err != nil {
		analyzeError(err)
	}
	
	// 创建用户 (Create user)
	fmt.Println("2. Create user:")
	user := &User{ID: "maint_user", Name: "Maintenance User", Email: "maint@example.com", Age: 25}
	err = repo.CreateUser(user)
	if err != nil {
		analyzeError(err)
	}
	
	// 禁用维护模式 (Disable maintenance mode)
	repo.SetMaintenanceMode(false)
	fmt.Println("Maintenance mode disabled")
	
	// 再次尝试创建用户 (Try creating user again)
	fmt.Println("\n3. Create user after maintenance:")
	err = repo.CreateUser(user)
	if err != nil {
		analyzeError(err)
	} else {
		fmt.Printf("✓ User created successfully after maintenance\n\n")
	}
}

// demonstrateErrorCodeComparison 演示错误码比较
// (demonstrateErrorCodeComparison demonstrates error code comparison)
func demonstrateErrorCodeComparison() {
	fmt.Println("=== Demonstrating Error Code Comparison ===")
	fmt.Println()
	
	repo := NewUserRepository()
	
	// 生成不同类型的错误 (Generate different types of errors)
	errorList := []error{}
	
	// 用户不存在错误 (User not found error)
	_, err1 := repo.GetUser("nonexistent")
	if err1 != nil {
		errorList = append(errorList, err1)
	}
	
	// 验证错误 (Validation error)
	err2 := repo.CreateUser(&User{ID: "", Name: "Invalid", Email: "test@example.com", Age: 25})
	if err2 != nil {
		errorList = append(errorList, err2)
	}
	
	// 数据库连接错误 (Database connection error)
	_, err3 := repo.GetUser("db_error")
	if err3 != nil {
		errorList = append(errorList, err3)
	}
	
	// 比较错误码 (Compare error codes)
	fmt.Println("Error code comparison:")
	for i, err := range errorList {
		fmt.Printf("Error %d:\n", i+1)
		coder := errors.GetCoder(err)
		if coder != nil {
			fmt.Printf("  Code: %d\n", coder.Code())
			fmt.Printf("  String: %s\n", coder.String())
			
			// 检查特定错误类型 (Check specific error types)
			switch coder.Code() {
			case ErrUserNotFound.Code():
				fmt.Printf("  Type: User not found error\n")
			case ErrInvalidUserData.Code():
				fmt.Printf("  Type: Validation error\n")
			case ErrDatabaseConnection.Code():
				fmt.Printf("  Type: Database error\n")
			}
		}
		fmt.Println()
	}
}

func main() {
	fmt.Println("=== Error Codes Example ===")
	fmt.Println("This example demonstrates structured error handling with error codes.")
	fmt.Println()
	
	// 1. 初始化日志 (Initialize logging)
	logOpts := log.NewOptions()
	logOpts.Level = "info"
	logOpts.Format = "text"
	logOpts.EnableColor = true
	log.Init(logOpts)
	logger := log.Std()
	
	// 2. 演示基础错误码使用 (Demonstrate basic error code usage)
	demonstrateBasicErrorCodes()
	
	// 3. 演示用户操作中的错误码 (Demonstrate error codes in user operations)
	demonstrateUserOperations()
	
	// 4. 演示服务层错误码 (Demonstrate service layer error codes)
	demonstrateServiceLayerErrors()
	
	// 5. 演示维护模式错误 (Demonstrate maintenance mode errors)
	demonstrateMaintenanceMode()
	
	// 6. 演示错误码比较 (Demonstrate error code comparison)
	demonstrateErrorCodeComparison()
	
	logger.Info("Error codes example completed successfully")
	fmt.Println("=== Example completed successfully ===")
} 