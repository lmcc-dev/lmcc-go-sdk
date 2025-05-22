
# lmcc-go-sdk Errors Module Specification

**Version:** 1.0.0
**Authors:** Martin, AI Assistant

[中文版 (Chinese Version)](SPECIFICATION_zh.md)

## 1. Introduction

The `pkg/errors` module in the `lmcc-go-sdk` provides a flexible and structured way to handle errors in Go applications. It extends the standard library's error handling capabilities by introducing a `Coder` system for categorized error information, automatic stack trace capturing, enhanced formatting options, and error aggregation.

The primary goals of this module are:
- To provide more context with errors beyond a simple string message.
- To enable consistent error handling and reporting across different parts of an application.
- To facilitate easier debugging through detailed stack traces and structured error information.
- To remain compatible with standard Go error handling idioms (e.g., `errors.Is`, `errors.As`).

## 2. Core Concepts

### 2.1. Coder
A `Coder` is an interface that represents a specific type of error with a unique code, an associated HTTP status, a descriptive string, and an optional reference URI. It allows for programmatic inspection and handling of errors based on their classification.

### 2.2. Stack Trace
The module automatically captures a stack trace whenever an error is created or wrapped. This stack trace helps pinpoint the origin of an error.

### 2.3. Error Wrapping
Errors can be wrapped to add contextual information. The original error (the "cause") is preserved and can be unwrapped, which is compatible with `errors.Unwrap` from the standard library.

### 2.4. Error Aggregation
Multiple errors can be collected into a single `ErrorGroup` instance. This is useful for scenarios where an operation might result in several failures that need to be reported together.

## 3. Core Types

### 3.1. `Coder` Interface

Defined in `pkg/errors/api.go`.

```go
type Coder interface {
    error // Embeds the standard error interface

    Code() int
    String() string
    HTTPStatus() int
    Reference() string
}
```

- **`error`**: Embedding `error` allows `Coder` instances to be used directly as standard errors.
- **`Code() int`**: Returns the unique integer code of the error.
- **`String() string`**: Returns the descriptive string of the error Coder (distinct from `Error()`).
- **`Error() string`**: (From embedded `error` interface) Typically returns the same as `String()` for `basicCoder` implementations.
- **`HTTPStatus() int`**: Returns the suggested HTTP status code for this error (e.g., 404, 500).
- **`Reference() string`**: Returns a URL or document link for more information about the error.

The primary implementation is `basicCoder` in `pkg/errors/coder.go`.

### 3.2. `StackTrace` Type

Defined in `pkg/errors/stack.go`.

```go
type StackTrace []Frame
type Frame uintptr
```

- **`Frame`**: Represents a single program counter in a stack trace.
- **`StackTrace`**: A slice of `Frame`s, representing the call stack from the innermost (most recent) to the outermost call.
- It implements `fmt.Formatter` to allow detailed printing of the stack trace when an error is formatted with `%+v`.

### 3.3. `fundamental` Error Type

Defined in `pkg/errors/errors.go`.

This is the most basic error type created by the module. It contains an error message and a stack trace. It does not wrap another error.

- **Fields**:
    - `msg string`: The error message.
    - `stack StackTrace`: The stack trace captured at the point of creation.
- **Key Methods**:
    - `Error() string`: Returns the error message.
    - `Unwrap() error`: Returns `nil`.
    - `Format(s fmt.State, verb rune)`: Implements `fmt.Formatter`.
        - `%s`, `%v`: Print the error message.
        - `%+v`: Print the error message and the stack trace.

### 3.4. `wrapper` Error Type

Defined in `pkg/errors/errors.go`.

This type wraps an existing error, adding a new message and capturing a new stack trace at the point of wrapping.

- **Fields**:
    - `msg string`: The additional message for this wrapper.
    - `cause error`: The underlying error being wrapped.
    - `stack StackTrace`: The stack trace captured at the point of wrapping.
- **Key Methods**:
    - `Error() string`: Returns `msg + ": " + cause.Error()`.
    - `Unwrap() error`: Returns `cause`.
    - `Format(s fmt.State, verb rune)`: Implements `fmt.Formatter`.
        - `%s`, `%v`: Print the combined error message.
        - `%+v`: Print the combined error message and the stack trace of this wrapper.

### 3.5. `withCode` Error Type

Defined in `pkg/errors/errors.go`.

This type associates a `Coder` with an underlying error and captures a stack trace at the point of association.

- **Fields**:
    - `cause error`: The underlying error.
    - `coder Coder`: The `Coder` associated with this error.
    - `stack StackTrace`: The stack trace captured at the point of attaching the `Coder`.
- **Key Methods**:
    - `Error() string`: Returns a message including the `Coder`'s string and the `cause`'s error message (e.g., `coder.String() + ": " + cause.Error()`). If the `coder` or its string is empty, it returns `cause.Error()`.
    - `Unwrap() error`: Returns `cause`.
    - `Coder() Coder`: Returns the associated `coder`.
    - `Format(s fmt.State, verb rune)`: Implements `fmt.Formatter`.
        - `%s`, `%v`: Print the error message (including `Coder` info).
        - `%+v`: Print the error message (including `Coder` info) and the stack trace of this `withCode` wrapper.

### 3.6. `ErrorGroup` Type

Defined in `pkg/errors/group.go`.

This type allows for collecting and managing a list of errors as a single error.

- **Fields**:
    - `errs []error`: A slice of errors.
    - `message string`: An optional overarching message for the group.
- **Key Methods**:
    - `Add(err error)`: Adds a non-nil error to the group.
    - `Errors() []error`: Returns the slice of contained errors.
    - `Error() string`: Returns a string representation combining the group message (if any) and all contained error messages.
    - `Unwrap() []error`: Returns the slice of contained errors, making it compatible with Go 1.20+ multi-error unwrapping (e.g., for `errors.Is` and `errors.As`).
    - `Format(s fmt.State, verb rune)`: Implements `fmt.Formatter`.
        - `%+v`: Prints the group message (if any) followed by a detailed, multi-line representation of each contained error, including their stack traces if available and formatted with `%+v`.

## 4. Error Creation

### 4.1. `New(text string) error`
Creates a new `fundamental` error with the given text message and captures the current stack trace.

```go
err := errors.New("an unexpected failure occurred")
```

### 4.2. `Errorf(format string, args ...interface{}) error`
Creates a new `fundamental` error with a formatted message (using `fmt.Sprintf`) and captures the current stack trace.

```go
id := "123"
err := errors.Errorf("failed to process item %s", id)
```

### 4.3. `Wrap(err error, message string) error`
Wraps an existing `err` with an additional `message`. It returns a new `wrapper` error. If `err` is `nil`, `Wrap` returns `nil`. A new stack trace is captured at the point of wrapping.

```go
dbErr := // ... some error from a database call
contextualErr := errors.Wrap(dbErr, "failed during user registration")
// contextualErr.Error() -> "failed during user registration: original dbErr message"
```

### 4.4. `Wrapf(err error, format string, args ...interface{}) error`
Wraps an existing `err` with an additional formatted `message`. It returns a new `wrapper` error. If `err` is `nil`, `Wrapf` returns `nil`. A new stack trace is captured.

```go
id := "txn-007"
ioErr := // ... some I/O error
contextualErr := errors.Wrapf(ioErr, "failed to write data for transaction %s", id)
```

### 4.5. `NewWithCode(coder Coder, text string) error`
Creates a new error that associates a `Coder` with a message. The underlying error (the `cause`) is a new `fundamental` error created from `text`. A `withCode` error is returned. If `coder` is `nil`, a predefined `unknownCoder` is used.

```go
var ErrUserNotFound = errors.NewCoder(10101, 404, "User not found", "")
err := errors.NewWithCode(ErrUserNotFound, "could not find user 'jane.doe'")
```

### 4.6. `ErrorfWithCode(coder Coder, format string, args ...interface{}) error`
Creates a new error that associates a `Coder` with a formatted message. The underlying `cause` is a new `fundamental` error. A `withCode` error is returned. If `coder` is `nil`, `unknownCoder` is used.

```go
var ErrPaymentFailed = errors.NewCoder(20202, 400, "Payment processing failed", "")
orderID := "order-456"
err := errors.ErrorfWithCode(ErrPaymentFailed, "payment failed for order %s due to insufficient funds", orderID)
```

### 4.7. `WithCode(err error, coder Coder) error`
Annotates an existing `err` with a `Coder`. Returns a `withCode` error. If `err` is `nil`, `WithCode` returns `nil`. If `coder` is `nil`, `unknownCoder` is used.

```go
var ErrPermissionDenied = errors.NewCoder(30303, 403, "Permission denied", "")
originalErr := errors.New("access to resource forbidden")
errWithCode := errors.WithCode(originalErr, ErrPermissionDenied)
```

### 4.8. `WithMessage(err error, message string) error`
Convenience function, equivalent to `Wrap(err, message)`.

### 4.9. `WithMessagef(err error, format string, args ...interface{}) error`
Convenience function, equivalent to `Wrapf(err, format, args...)`.

### 4.10. `NewErrorGroup(message ...string) *ErrorGroup`
Creates a new `ErrorGroup`. An optional overarching message for the group can be provided.

```go
eg := errors.NewErrorGroup("Configuration validation failed")
eg.Add(errors.New("missing API key"))
eg.Add(errors.New("database URL is invalid"))
```

## 5. Error Inspection and Handling

### 5.1. `errors.Is(err, target error) bool` (Standard Library)
This module's error types (`fundamental`, `wrapper`, `withCode`, `ErrorGroup`) are compatible with the standard library's `errors.Is` function.
- For `fundamental` errors, `Is` checks if the target is also a `fundamental` error with the same message.
- For `wrapper` errors, `Is` will unwrap to the `cause`.
- For `withCode` errors, `Is` will first try to match the `Coder` if the `target` is a `Coder` or has a `Coder()` method that returns a matching code. If not, it unwraps to its `cause`.
- For `ErrorGroup`, `errors.Is` will check if any of the contained errors (or their causes) match the `target`.
- A `Coder` itself (which implements `error`) can be used as a `target`.

```go
// Assuming ErrUserNotFound is a Coder
if errors.Is(err, ErrUserNotFound) {
    // Handle user not found specifically
}
```

### 5.2. `errors.As(err error, target interface{}) bool` (Standard Library)
Compatible with `errors.As`.
- For `fundamental` errors, can extract `**fundamental`.
- For `wrapper` errors, can extract `**wrapper` and then unwraps to `cause`.
- For `withCode` errors, can extract `**withCode`, or `*Coder` (to get the associated `Coder`), and then unwraps to `cause`.
- For `ErrorGroup`, `errors.As` will attempt to find a contained error that can be assigned to `target`.

```go
var coder errors.Coder
if errors.As(err, &coder) {
    fmt.Printf("Error Code: %d, HTTP Status: %d\n", coder.Code(), coder.HTTPStatus())
}
```

### 5.3. `Cause(err error) error`
Recursively unwraps `err` using its `Unwrap()` method (or a custom `Cause()` method if available on an intermediate error type) until the original, underlying error is found. If the error does not have a cause, the error itself is returned.

```go
rootCause := errors.Cause(complexWrappedError)
```

### 5.4. `GetCoder(err error) Coder`
Recursively unwraps `err` and returns the first non-nil `Coder` found in the error chain. If no `Coder` is found, it returns `nil`.

```go
coder := errors.GetCoder(err)
if coder != nil {
    // Use the coder
}
```

### 5.5. `IsCode(err error, c Coder) bool`
Defined in `pkg/errors/api.go`. Checks if the error `err` (or any error in its chain) has a `Coder` whose code matches the code of the provided `Coder c`. This is useful for checking against a specific error *category* represented by a `Coder`'s code, even if the `Coder` instances are different.

```go
if errors.IsCode(err, ErrPaymentFailed) { // Assuming ErrPaymentFailed is a Coder
    // Logic for handling any payment failure, regardless of specific message
}
```

## 6. Formatted Output

All error types in this module (`fundamental`, `wrapper`, `withCode`, `ErrorGroup`) implement the `fmt.Formatter` interface.

- **`%s`, `%v`**:
    - For `fundamental`: Prints the error message.
    - For `wrapper`: Prints the combined message (`wrapper.msg: cause.Error()`).
    - For `withCode`: Prints the message including the `Coder`'s string (e.g., `coder.String(): cause.Error()`).
    - For `ErrorGroup`: Prints the combined message of the group and all its contained errors.
- **`%+v`**:
    - For `fundamental`: Prints the error message followed by its stack trace.
    - For `wrapper`: Prints the combined message, followed by the stack trace captured *at the point this wrapping*.
    - For `withCode`: Prints the message (including `Coder` info), followed by the stack trace captured *at the point this `Coder` was attached*.
    - For `ErrorGroup`: Prints the group's main message (if any), followed by a detailed, multi-line representation of each contained error. If a contained error also supports `%+v` (like those from this package), its full details including its own stack trace will be printed.

```go
// To print an error with its full stack trace(s):
fmt.Printf("%+v\n", err)
```

## 7. Predefined Coders

The `pkg/errors/coder.go` file defines a set of common `Coder` instances for convenience. Examples include:
- `unknownCoder`: A generic internal server error.
- `ErrInternalServer`: HTTP 500.
- `ErrNotFound`: HTTP 404.
- `ErrBadRequest`: HTTP 400.
- `ErrUnauthorized`: HTTP 401.
- `ErrForbidden`: HTTP 403.
- `ErrValidation`: HTTP 400 (typically for request validation issues).
- `ErrTimeout`: HTTP 504.
- `ErrTooManyRequests`: HTTP 429.
- `ErrOperationFailed`: A generic operation failure.
- Specific coders for `config` package errors (e.g., `ErrConfigFileRead`).
- Specific coders for `log` package errors (e.g., `ErrLogOptionInvalid`).

Refer to `pkg/errors/coder.go` for the complete list and their definitions.

## 8. Best Practices and Examples

- **Use `New` or `Errorf` for initial errors**: When an error originates in your code.
- **Use `Wrap` or `Wrapf` to add context**: When an error from a lower layer needs more high-level information before being propagated.
- **Use `WithCode` or `NewWithCode`/`ErrorfWithCode` to categorize errors**: Associate a predefined or custom `Coder` to enable programmatic handling or specific logging.
- **Check errors with `errors.Is` and `errors.As`**: Prefer these over type assertions for flexibility and compatibility with wrapped errors. Use `IsCode` for checking against specific `Coder` codes.
- **Use `Cause` to find the root error**: Useful for logging the original failure if it's deeply wrapped.
- **Use `GetCoder` to extract error codes**: When you need to act based on the `Coder` information.
- **Format with `%+v` for detailed debugging logs**: This provides the message and full stack trace(s).
- **Define application-specific `Coders`**: Create your own `Coder` instances in your application or domain-specific packages for better error classification.
- **Use `ErrorGroup` for multiple errors**: When an operation can fail in multiple ways simultaneously (e.g., validating multiple fields of a request).

For runnable examples, please see `pkg/errors/example_test.go`.

## 9. Go Version Specifics

- **Error Wrapping**: The `Unwrap()` method on `wrapper`, `withCode`, and the `Unwrap() []error` on `ErrorGroup` ensure compatibility with `errors.Is` and `errors.As` in Go 1.13+ (for single error unwrapping) and Go 1.20+ (for multi-error unwrapping with `ErrorGroup`).
- **Stack Traces**: Stack trace collection uses the `runtime` package.

This specification aims to provide a comprehensive guide to using the `pkg/errors` module. For specific implementation details, refer to the source code. 