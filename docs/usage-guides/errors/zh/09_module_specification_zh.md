<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

## `pkg/errors` 模块规范 (Module Specification)

本文档提供了 `pkg/errors` 模块的详细规范，概述了其公共 API、接口、函数和预定义组件。该模块通过提供堆栈跟踪、错误码、错误包装和聚合机制来增强 Go 的标准错误处理能力。

**[English Version (英文版说明)](../en/09_module_specification.md)**

### 1. 概述 (Overview)

`pkg/errors` 模块旨在使错误处理更具信息量且更易于管理。主要特性包括：
- 在错误创建和包装时自动捕获堆栈跟踪。
- `Coder` 接口，用于通过代码、HTTP 状态和参考文档对错误进行分类。
- 用于创建、包装和检查错误的函数，这些函数与标准库实用程序（如 `standardErrors.Is` 和 `standardErrors.As`）兼容。
- `ErrorGroup` 类型，用于将多个错误收集到一个错误中。

### 2. `Coder` 接口 (The `Coder` Interface)

`Coder` 接口是错误分类的核心。它允许错误携带机器可读的代码、关联的 HTTP 状态、人类可读的消息以及对更详细文档的引用。

```go
type Coder interface {
	error // 嵌入标准错误接口 (Embeds the standard error interface)

	// Code 返回此错误的整数代码。
	// (Code returns the integer code of this error.)
	Code() int

	// String 返回此错误代码的字符串表示形式
	// （通常是错误描述或类型）。
	// (String returns the string representation of this error code
	// (typically the error description or type).)
	String() string

	// HTTPStatus 返回关联的 HTTP 状态码（如果适用）。
	// 如果不适用，可能返回 0 或哨兵值。
	// (HTTPStatus returns the associated HTTP status code, if applicable.
	// May return 0 or a sentinel value if not applicable.)
	HTTPStatus() int

	// Reference 返回此错误代码的 URL 或文档引用。
	// 如果不适用，可能返回空字符串。
	// (Reference returns a URL or document reference for this error code.
	// May return an empty string if not applicable.)
	Reference() string
}
```

**工厂函数 (Factory Function):**
- **`NewCoder(code int, httpStatus int, description string, reference string) Coder`**: 创建 `Coder` 的实例。这是定义自定义错误类别的主要方式。
  (Creates an instance of `Coder`. This is the primary way to define custom error categories.)

### 3. 创建基本错误 (Creating Basic Errors)

这些函数创建新的错误实例，并在创建点自动捕获堆栈跟踪。

(These functions create new error instances, automatically capturing a stack trace at the point of creation.)

- **`New(text string) error`**: 返回一个格式化为给定文本的错误。
  (Returns an error that formats as the given text.)
- **`Errorf(format string, args ...interface{}) error`**: 根据格式说明符进行格式化，并将字符串作为错误返回。
  (Formats according to a format specifier and returns the string as an error.)

### 4. 包装错误 (添加上下文) (Wrapping Errors (Adding Context))

包装错误会向现有错误添加上下文信息，同时保留原始错误（"原因"）。如果包装器本身是新的错误实例，或者底层错误没有 `pkg/errors` 堆栈跟踪，则还会捕获包装站点的堆栈跟踪。

(Wrapping an error adds contextual information to an existing error while preserving the original error (the "cause"). The wrapping site's stack trace is also captured if the wrapper itself is a new error instance or if the underlying error did not have a `pkg/errors` stacktrace.)

- **`Wrap(err error, message string) error`**: 返回一个用新消息注释 `err` 的错误。如果 `err` 为 `nil`，`Wrap` 返回 `nil`。返回的错误将有一个返回 `err` 的 `Cause` 方法。
  (Returns an error annotating `err` with a new message. If `err` is `nil`, `Wrap` returns `nil`. The returned error will have a `Cause` method returning `err`.)
- **`Wrapf(err error, format string, args ...interface{}) error`**: 根据格式说明符进行格式化，并返回一个用该消息注释 `err` 的错误。如果 `err` 为 `nil`，`Wrapf` 返回 `nil`。
  (Formats according to a format specifier and returns an error annotating `err` with that message. If `err` is `nil`, `Wrapf` returns `nil`.)
- **`WithMessage(err error, message string) error`**: `Wrap` 的别名。
  (An alias for `Wrap`.)
- **`WithMessagef(err error, format string, args ...interface{}) error`**: `Wrapf` 的别名。
  (An alias for `Wrapf`.)

### 5. 使用错误码 (`Coder`)

这些函数允许创建与 `Coder` 关联的错误，或将 `Coder` 附加到现有错误。

- **`NewWithCode(coder Coder, text string) error`**: 创建一个新错误，其中包含提供的 `Coder` 和一个简单的文本消息。错误消息将是 `coder.String()` 和 `text` 的组合。
- **`ErrorfWithCode(coder Coder, format string, args ...interface{}) error`**: 创建一个新错误，其中包含提供的 `Coder` 和一个格式化的消息。错误消息将是 `coder.String()` 和格式化字符串的组合。
- **`WithCode(err error, coder Coder) error`**: 使用 `Coder` 注释现有错误 `err`。如果 `err` 为 `nil`，则返回 `nil`。原始错误 `err` 成为 `Cause`。错误消息将是 `coder.String()` 和 `err.Error()` 的组合。

### 6. 错误聚合 (`ErrorGroup`)

`ErrorGroup` 允许将多个错误收集到单个错误对象中。当一个操作涉及多个可能独立失败的子任务时（例如，验证表单的多个字段），这非常有用。

( `ErrorGroup` allows collecting multiple errors into a single error object. This is useful when an operation involves multiple sub-tasks that can fail independently (e.g., validating multiple fields of a form).)

- **`NewErrorGroup(message ...string) *ErrorGroup`**: 创建一个新的 `ErrorGroup`。可以为组提供一个可选的总体消息。
  (Creates a new `ErrorGroup`. An optional overarching message for the group can be provided.)

**`ErrorGroup` 方法 (Methods):**
- **`Add(err error)`**: 将错误添加到组中。如果 `err` 为 `nil`，则不执行任何操作。
  (Adds an error to the group. If `err` is `nil`, it does nothing.)
- **`Errors() []error`**: 返回添加到组中的所有错误的切片。如果未添加任何错误，则返回 `nil`。
  (Returns a slice of all errors added to the group. Returns `nil` if no errors were added.)
- **`Error() string`**: 返回组中所有错误的字符串表示形式，以组的总体消息（如果有）为前缀。单个错误用分号分隔。
  (Returns a string representation of all errors in the group, prefixed by the group's overarching message (if any). Individual errors are separated by semicolons.)
- **`Unwrap() []error`**: 实现 `Unwrap() []error` 模式 (Go 1.20+)，允许 `standardErrors.Is` 和 `standardErrors.As` 与收集到的错误一起工作。组中的每个错误都是匹配的潜在候选项。
  (Implements the `Unwrap() []error` pattern (Go 1.20+) allowing `standardErrors.Is` and `standardErrors.As` to work with the collected errors. Each error in the group is a potential candidate for matching.)
- **`Format(s fmt.State, verb rune)`**: 实现 `fmt.Formatter`。当与 `"%+v"` 一起使用时，它会打印组的消息，然后是每个包含错误的详细格式，包括它们各自的堆栈跟踪（如果可用）。
  (Implements `fmt.Formatter`. When used with `"%+v"`, it prints the group's message followed by detailed formatting of each contained error, including their individual stack traces if available.)

### 7. 检查错误 (Inspecting Errors)

这些函数有助于反思错误。

(These functions help in introspecting errors.)

- **`Cause(err error) error`**: 如果可能，返回错误的根本原因。如果错误实现了 `interface { Cause() error }` 或 `interface { Unwrap() error }` 接口，则它包装了另一个错误。如果 `err` 没有实现任一接口，`Cause` 返回 `err` 本身。
  (Returns the underlying cause of the error, if possible. An error wraps another error if it implements the `interface { Cause() error }` or `interface { Unwrap() error }` interface. If `err` does not implement either, `Cause` returns `err` itself.)
- **`GetCoder(err error) Coder`**: 遍历错误链（通过 `Unwrap` 或 `Cause`）并返回遇到的第一个 `Coder`。如果错误链中没有错误具有关联的 `Coder`，则返回 `nil` （或者，如果已配置，则返回默认的"未知"Coder，尽管当前实现似乎返回 `nil`）。
  (Traverses the error chain (via `Unwrap` or `Cause`) and returns the first `Coder` encountered. If no error in the chain has an associated `Coder`, it returns `nil` (or a default "unknown" Coder if configured, though current implementation seems to return `nil`).)
- **`IsCode(err error, c Coder) bool`**:报告 `err` 的链中是否有任何错误具有 `Coder`，其 `Code()` 与 `c.Code()` 匹配。这对于根据其数字代码检查错误的类别很有用。**注意**：此函数通过 `Unwrap() []error` 方法检查组内的所有错误，从而支持 `ErrorGroup`。
  (Reports whether any error in `err`'s chain has a `Coder` whose `Code()` matches `c.Code()`. This is useful for checking an error's category based on its numeric code. **Note**: This function supports `ErrorGroup` by checking all errors within the group through its `Unwrap() []error` method.)

**与标准库的兼容性 (Compatibility with Standard Library):**
- **`standardErrors.Is(err, target error) bool`**: 按预期工作。如果 `target` 是一个 `Coder` 实例（如预定义的 `ErrNotFound`），它会检查 `err` 或其任何原因是否是该特定的 `Coder` 实例。**重要**：对于使用 `WithCode` 创建的错误，`Is` 方法比较 `Coder` 代码而不是实例，这意味着具有相同代码的两个不同 `Coder` 实例将被认为是相等的。
  (Works as expected. If `target` is a `Coder` instance (like predefined `ErrNotFound`), it checks if `err` or any of its causes is that specific `Coder` instance. **Important**: For errors created with `WithCode`, the `Is` method compares `Coder` codes rather than instances, meaning two different `Coder` instances with the same code will be considered equal.)
- **`standardErrors.As(err, target interface{}) bool`**: 按预期工作。如果错误链中的错误嵌入了一个 `Coder` 并匹配 `Coder` 接口，或者要提取任何其他自定义错误类型，则可以使用它来提取 `Coder`。
  (Works as expected. It can be used to extract a `Coder` if an error in the chain embeds one and matches the `Coder` interface, or to extract any other custom error type.)

### 8. 格式化错误 (堆栈跟踪) (Formatting Errors (Stack Traces))

由 `pkg/errors` 函数创建或包装的错误会捕获堆栈跟踪。当使用 `fmt.Printf` 或类似函数的 `"%+v"` 动词格式化错误时，将打印此堆栈跟踪。

(Errors created or wrapped by `pkg/errors` functions capture a stack trace. This stack trace is printed when the error is formatted using the `"%+v"` verb with `fmt.Printf` or similar functions.)

- `pkg/errors` 中的 `fundamental`（用于 `New`、`Errorf`）、`wrapper`（用于 `Wrap`、`Wrapf`）和 `withCode` 错误类型都实现了 `fmt.Formatter` 接口，以便为 `"%+v"` 提供包括堆栈跟踪在内的详细输出。
  (The `fundamental` (for `New`, `Errorf`), `wrapper` (for `Wrap`, `Wrapf`), and `withCode` error types within `pkg/errors` all implement the `fmt.Formatter` interface to provide detailed output including stack traces for `"%+v"`.)
- `ErrorGroup` 也实现了 `fmt.Formatter` 以显示其所有包含错误的详细信息。
  (`ErrorGroup` also implements `fmt.Formatter` to show details of all its contained errors.)

### 9. 预定义的 `Coder` 实例 (Predefined `Coder` Instances)

`pkg/errors` 模块为常见的错误场景提供了几个预定义的 `Coder` 实例。这些是导出的变量。

(The `pkg/errors` module provides several predefined `Coder` instances for common error scenarios. These are exported variables.)

| 变量名 (Variable Name)   | 代码 (Code) | HTTP 状态 (HTTP Status) | 默认消息 (Default Message)             |
|--------------------------|-----------|-----------------------|--------------------------------------|
| `ErrUnknown`             | -1        | 500                   | 发生了内部服务器错误 (An internal server error occurred) |
| `ErrInternalServer`      | 100001    | 500                   | 内部服务器错误 (Internal server error)       |
| `ErrNotFound`            | 100002    | 404                   | 未找到 (Not found)                   |
| `ErrBadRequest`          | 100003    | 400                   | 错误的请求 (Bad request)                 |
| `ErrPermissionDenied`    | 100004    | 403                   | 权限被拒绝 (Permission denied)           |
| `ErrConflict`            | 100005    | 409                   | 冲突 (Conflict)                    |
| `ErrValidation`          | 100006    | 400                   | 验证失败 (Validation failed)           |
| `ErrResourceUnavailable` | 100007    | 503                   | 资源不可用 (Resource unavailable)        |
| `ErrTooManyRequests`     | 100008    | 429                   | 请求过多 (Too many requests)           |
| `ErrUnauthorized`        | 100009    | 401                   | 未授权 (Unauthorized)                |
| `ErrTimeout`             | 100010    | 504                   | 超时 (Timeout)                     |
| `ErrNotImplemented`      | 100011    | 501                   | 未实现 (Not implemented)             |
| `ErrNotSupported`        | 100012    | 501                   | 不支持 (Not supported)               |
| `ErrAlreadyExists`       | 100013    | 409                   | 已存在 (Already exists)              |
| `ErrDataLoss`            | 100014    | 500                   | 数据丢失 (Data loss)                   |
| `ErrDatabase`            | 100015    | 500                   | 数据库错误 (Database error)              |
| `ErrEncoding`            | 100016    | 500                   | 编码错误 (Encoding error)              |
| `ErrDecoding`            | 100017    | 500                   | 解码错误 (Decoding error)              |
| `ErrNetwork`             | 100018    | 500                   | 网络错误 (Network error)               |
| `ErrFilesystem`          | 100019    | 500                   | 文件系统错误 (Filesystem error)            |
| `ErrConfiguration`       | 100020    | 500                   | 配置错误 (Configuration error)         |
| `ErrAuthentication`      | 100021    | 401                   | 身份验证失败 (Authentication failed)       |
| `ErrAuthorization`       | 100022    | 403                   | 授权失败 (Authorization failed)        |
| `ErrRateLimitExceeded`   | 100023    | 429                   | 超出速率限制 (Rate limit exceeded)         |
| `ErrInvalidInput`        | 100024    | 400                   | 无效输入 (Invalid input)               |
| `ErrStateMismatch`       | 100025    | 409                   | 状态不匹配 (State mismatch)              |
| `ErrOperationAborted`    | 100026    | 499                   | 操作中止 (Operation aborted)           |
| `ErrResourceExhausted`   | 100027    | 507                   | 资源耗尽 (Resource exhausted)          |
| `ErrExternalService`     | 100028    | 502                   | 外部服务错误 (External service error)      |
| `ErrMaintenance`         | 100029    | 503                   | 维护中 (Maintenance)                 |
| `ErrLogOptionInvalid`    | 300001    | 500                   | 无效的日志选项 (Invalid log option)          |
| `ErrLogRotationSetup`    | 300002    | 500                   | 日志轮转设置失败 (Log rotation setup failed)   |
| `ErrLogWrite`            | 300003    | 500                   | 日志写入失败 (Log write failure)           |
| `ErrLogReconfigure`      | 300004    | 500                   | 日志重新配置失败 (Log reconfiguration failed)  |
| `ErrLogTargetCreate`     | 300005    | 500                   | 日志目标创建失败 (Log target creation failed)  |
| `ErrLogTargetNotSupported`| 300006   | 500                   | 不支持的日志目标 (Log target not supported)    |
| `ErrLogBufferFull`       | 300007    | 500                   | 日志缓冲区已满 (Log buffer full)             |
| `ErrLogRotationDirInvalid`| 300008   | 500                   | 无效的日志轮转目录 (Invalid log rotation directory)|

**Coder 的实用函数 (Utility functions for Coders):**
- **`IsUnknownCoder(coder Coder) bool`**: 检查给定的 `coder` 是否是预定义的 `ErrUnknown`。
  (Checks if the given `coder` is the predefined `ErrUnknown`.)
- **`GetUnknownCoder() Coder`**: 返回预定义的 `ErrUnknown` Coder。
  (Returns the predefined `ErrUnknown` coder.)

本规范应有助于很好地理解如何有效地使用 `pkg/errors` 模块。
有关更多示例和最佳实践，请参阅使用指南。

(This specification should provide a good understanding of how to use the `pkg/errors` module effectively.
Refer to the usage guides for more examples and best practices.) 