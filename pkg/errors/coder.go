/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package errors

// basicCoder is a basic implementation of the Coder interface.
// basicCoder 是 Coder 接口的一个基础实现。
type basicCoder struct {
	C    int    // Code, 错误码 (Error code)
	HTTP int    // HTTPStatus, HTTP 状态码 (HTTP status code)
	Ext  string // String, 错误描述 (Error description)
	Ref  string // Reference, 参考文档 (Reference document)
}

// Code returns the integer code of this error.
// Code 返回错误的整数码。
func (c *basicCoder) Code() int {
	return c.C
}

// String returns the string representation of this error code.
// String 返回错误码的字符串表示。
func (c *basicCoder) String() string {
	return c.Ext
}

// Error returns the string representation of this error code, satisfying the error interface.
// Error 返回错误码的字符串表示，满足 error 接口。
func (c *basicCoder) Error() string {
	return c.Ext // Same as String() for convenience with error interface
}

// HTTPStatus returns the associated HTTP status code.
// HTTPStatus 返回关联的 HTTP 状态码。
func (c *basicCoder) HTTPStatus() int {
	return c.HTTP
}

// Reference returns a URL or document reference for this error code.
// Reference 返回此错误码的 URL 或文档参考。
func (c *basicCoder) Reference() string {
	return c.Ref
}

// NewCoder creates a new Coder instance.
// NewCoder 创建一个新的 Coder 实例。
func NewCoder(code int, httpStatus int, description string, reference string) Coder {
	return &basicCoder{
		C:    code,
		HTTP: httpStatus,
		Ext:  description,
		Ref:  reference,
	}
}

// --- Predefined Coders --- (预定义的 Coder)

var (
	// unknownCoder is an unknown error (usually wrapped from an external error source).
	// unknownCoder 是一个未知错误 (通常是从外部错误源包装而来)。
	unknownCoder = NewCoder(-1, 500, "An internal server error occurred", "")

	// ErrInternalServer represents an internal server error (500).
	// ErrInternalServer 表示内部服务器错误 (500)。
	ErrInternalServer = NewCoder(100001, 500, "Internal server error", "")

	// ErrNotFound represents a resource not found error (404).
	// ErrNotFound 表示资源未找到错误 (404)。
	ErrNotFound = NewCoder(100002, 404, "Resource not found", "")

	// ErrBadRequest represents a bad request error (400).
	// ErrBadRequest 表示错误的请求错误 (400)。
	ErrBadRequest = NewCoder(100003, 400, "Bad request", "")

	// ErrUnauthorized represents an unauthorized access error (401).
	// ErrUnauthorized 表示未经授权的访问错误 (401)。
	ErrUnauthorized = NewCoder(100004, 401, "Unauthorized", "")

	// ErrForbidden represents a forbidden access error (403).
	// ErrForbidden 表示禁止访问错误 (403)。
	ErrForbidden = NewCoder(100005, 403, "Forbidden", "")

	// ErrValidation represents a data validation error (400 or 422).
	// ErrValidation 表示数据验证错误 (400 或 422)。
	ErrValidation = NewCoder(100006, 400, "Validation error", "")

	// ErrTimeout represents a request timeout error (504).
	// ErrTimeout 表示请求超时错误 (504)。
	ErrTimeout = NewCoder(100007, 504, "Request timeout", "")

	// ErrTooManyRequests represents a too many requests error (429).
	// ErrTooManyRequests 表示请求过多错误 (429).
	ErrTooManyRequests = NewCoder(100008, 429, "Too many requests", "")

	// ErrOperationFailed represents a generic operation failure.
	// ErrOperationFailed 表示通用操作失败。
	ErrOperationFailed = NewCoder(100009, 500, "Operation failed", "")

	// ErrConfigFileRead represents an error encountered while reading a configuration file.
	// ErrConfigFileRead 表示读取配置文件时遇到的错误。
	ErrConfigFileRead = NewCoder(200001, 500, "Config file read error", "https://lmcc-go-sdk.dev/docs/errors/config#file-read")

	// ErrConfigSetup represents an error encountered during configuration setup.
	// ErrConfigSetup 表示配置设置过程中遇到的错误。
	ErrConfigSetup = NewCoder(200002, 500, "Config setup error", "https://lmcc-go-sdk.dev/docs/errors/config#setup")

	// ErrConfigEnvBind represents an error encountered during environment variable binding for configuration.
	// ErrConfigEnvBind 表示配置的环境变量绑定过程中遇到的错误。
	ErrConfigEnvBind = NewCoder(200003, 500, "Config environment variable binding error", "")

	// ErrConfigDefaultTagParse represents an error encountered while parsing a 'default' struct tag for configuration.
	// ErrConfigDefaultTagParse 表示解析配置的 'default' 结构体标签时遇到的错误。
	ErrConfigDefaultTagParse = NewCoder(200004, 500, "Config default tag parsing error", "")

	// ErrConfigInternal represents an internal error within the configuration logic.
	// ErrConfigInternal 表示配置逻辑内部的错误。
	ErrConfigInternal = NewCoder(200005, 500, "Config internal error", "")

	// ErrConfigHotReload represents an error encountered during configuration hot-reloading.
	// ErrConfigHotReload 表示配置热重载过程中遇到的错误。
	ErrConfigHotReload = NewCoder(200006, 500, "Config hot-reload error", "")

	// --- Log Package Errors (pkg/log) ---

	// ErrLogInternal represents an internal error within the logging system.
	// ErrLogInternal 表示日志系统内部的错误。
	ErrLogInternal = NewCoder(300001, 500, "Log internal error", "")

	// ErrLogOptionInvalid represents an invalid option provided for logger configuration.
	// ErrLogOptionInvalid 表示为日志记录器配置提供了无效选项。
	ErrLogOptionInvalid = NewCoder(300002, 400, "Log option invalid", "")

	// ErrLogReconfigure represents an error encountered during logger reconfiguration.
	// ErrLogReconfigure 表示日志记录器重新配置期间遇到的错误。
	ErrLogReconfigure = NewCoder(300003, 500, "Log reconfiguration error", "")

	// ErrLogInitialization represents an error encountered during logger initialization.
	// ErrLogInitialization 表示日志记录器初始化期间遇到的错误。
	ErrLogInitialization = NewCoder(300004, 500, "Log initialization error", "")

	// ErrLogRotationSetup represents an error encountered during log rotation setup.
	// ErrLogRotationSetup 表示日志轮转设置期间遇到的错误。
	ErrLogRotationSetup = NewCoder(300005, 500, "Log rotation setup error", "")

	// ErrLogRotationDirCreate represents an error when creating a directory for log rotation.
	// ErrLogRotationDirCreate 表示为日志轮转创建目录时出错。
	ErrLogRotationDirCreate = NewCoder(300006, 500, "Log rotation directory creation error", "")

	// ErrLogRotationDirStat represents an error when stating a directory for log rotation.
	// ErrLogRotationDirStat 表示为日志轮转获取目录状态时出错。
	ErrLogRotationDirStat = NewCoder(300007, 500, "Log rotation directory stat error", "")

	// ErrLogRotationDirInvalid represents that the log rotation path exists but is not a directory.
	// ErrLogRotationDirInvalid 表示日志轮转路径存在但不是一个目录。
	ErrLogRotationDirInvalid = NewCoder(300008, 500, "Log rotation path exists but is not a directory", "")
)

// IsUnknownCoder checks if the Coder is the predefined unknownCoder.
// IsUnknownCoder 检查 Coder 是否是预定义的 unknownCoder。
func IsUnknownCoder(coder Coder) bool {
	// Compare by instance pointer, or by values if guaranteed unique & constant.
	// 通过实例指针比较，或者如果保证唯一和常量，则通过值比较。
	// For now, direct comparison should work as unknownCoder is a package-level var.
	// 目前，直接比较应该可行，因为 unknownCoder 是一个包级变量。
	return coder == unknownCoder
}

// GetUnknownCoder returns the predefined unknown Coder.
// GetUnknownCoder 返回预定义的未知 Coder。
func GetUnknownCoder() Coder {
	return unknownCoder
}
