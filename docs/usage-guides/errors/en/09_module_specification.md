<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

## `pkg/errors` Module Specification

This document provides a detailed specification of the `pkg/errors` module, outlining its public API, interfaces, functions, and predefined components. This module enhances Go's standard error handling by providing mechanisms for stack traces, error codes, error wrapping, and aggregation.

**[中文版说明 (Chinese Version)](../zh/09_module_specification_zh.md)**

### 1. Overview

The `pkg/errors` module is designed to make error handling more informative and manageable. Key features include:
- Automatic stack trace capture on error creation and wrapping.
- A `Coder` interface for categorizing errors with codes, HTTP statuses, and reference documentation.
- Functions for creating, wrapping, and inspecting errors that are compatible with standard library utilities like `standardErrors.Is` and `standardErrors.As`.
- An `ErrorGroup` type for collecting multiple errors into a single error.

### 2. The `Coder` Interface

The `Coder` interface is central to categorizing errors. It allows errors to carry a machine-readable code, an associated HTTP status, a human-readable message, and a reference to further documentation.

```go
type Coder interface {
	error // Embeds the standard error interface

	// Code returns the integer code of this error.
	Code() int

	// String returns the string representation of this error code
	// (typically the error description or type).
	String() string

	// HTTPStatus returns the associated HTTP status code, if applicable.
	// May return 0 or a sentinel value if not applicable.
	HTTPStatus() int

	// Reference returns a URL or document reference for this error code.
	// May return an empty string if not applicable.
	Reference() string
}
```

**Factory Function:**
- **`NewCoder(code int, httpStatus int, description string, reference string) Coder`**: Creates an instance of `Coder`. This is the primary way to define custom error categories.

### 3. Creating Basic Errors

These functions create new error instances, automatically capturing a stack trace at the point of creation.

- **`New(text string) error`**: Returns an error that formats as the given text.
- **`Errorf(format string, args ...interface{}) error`**: Formats according to a format specifier and returns the string as an error.

### 4. Wrapping Errors (Adding Context)

Wrapping an error adds contextual information to an existing error while preserving the original error (the "cause"). The wrapping site's stack trace is also captured if the wrapper itself is a new error instance or if the underlying error did not have a `pkg/errors` stacktrace.

- **`Wrap(err error, message string) error`**: Returns an error annotating `err` with a new message. If `err` is `nil`, `Wrap` returns `nil`. The returned error will have a `Cause` method returning `err`.
- **`Wrapf(err error, format string, args ...interface{}) error`**: Formats according to a format specifier and returns an error annotating `err` with that message. If `err` is `nil`, `Wrapf` returns `nil`.
- **`WithMessage(err error, message string) error`**: An alias for `Wrap`.
- **`WithMessagef(err error, format string, args ...interface{}) error`**: An alias for `Wrapf`.

### 5. Working with Error Codes (`Coder`)

These functions allow creating errors that are associated with a `Coder` or attaching a `Coder` to an existing error.

- **`NewWithCode(coder Coder, text string) error`**: Creates a new error that includes the provided `Coder` and a simple text message. The error message will be a combination of the `coder.String()` and the `text`.
- **`ErrorfWithCode(coder Coder, format string, args ...interface{}) error`**: Creates a new error that includes the provided `Coder` and a formatted message. The error message will be a combination of the `coder.String()` and the formatted string.
- **`WithCode(err error, coder Coder) error`**: Annotates an existing error `err` with a `Coder`. If `err` is `nil`, it returns `nil`. The original error `err` becomes the `Cause`. The error message will be a combination of `coder.String()` and `err.Error()`.

### 6. Error Aggregation (`ErrorGroup`)

`ErrorGroup` allows collecting multiple errors into a single error object. This is useful when an operation involves multiple sub-tasks that can fail independently (e.g., validating multiple fields of a form).

- **`NewErrorGroup(message ...string) *ErrorGroup`**: Creates a new `ErrorGroup`. An optional overarching message for the group can be provided.

**`ErrorGroup` Methods:**
- **`Add(err error)`**: Adds an error to the group. If `err` is `nil`, it does nothing.
- **`Errors() []error`**: Returns a slice of all errors added to the group. Returns `nil` if no errors were added.
- **`Error() string`**: Returns a string representation of all errors in the group, prefixed by the group's overarching message (if any). Individual errors are separated by semicolons.
- **`Unwrap() []error`**: Implements the `Unwrap() []error` pattern (Go 1.20+) allowing `standardErrors.Is` and `standardErrors.As` to work with the collected errors. Each error in the group is a potential candidate for matching.
- **`Format(s fmt.State, verb rune)`**: Implements `fmt.Formatter`. When used with `"%+v"`, it prints the group's message followed by detailed formatting of each contained error, including their individual stack traces if available.

### 7. Inspecting Errors

These functions help in introspecting errors.

- **`Cause(err error) error`**: Returns the underlying cause of the error, if possible. An error wraps another error if it implements the `interface { Cause() error }` or `interface { Unwrap() error }` interface. If `err` does not implement either, `Cause` returns `err` itself.
- **`GetCoder(err error) Coder`**: Traverses the error chain (via `Unwrap` or `Cause`) and returns the first `Coder` encountered. If no error in the chain has an associated `Coder`, it returns `nil` (or a default "unknown" Coder if configured, though current implementation seems to return `nil`).
- **`IsCode(err error, c Coder) bool`**: Reports whether any error in `err`'s chain has a `Coder` whose `Code()` matches `c.Code()`. This is useful for checking an error's category based on its numeric code. **Note**: This function supports `ErrorGroup` by checking all errors within the group through its `Unwrap() []error` method.

**Compatibility with Standard Library:**
- **`standardErrors.Is(err, target error) bool`**: Works as expected. If `target` is a `Coder` instance (like predefined `ErrNotFound`), it checks if `err` or any of its causes is that specific `Coder` instance. **Important**: For errors created with `WithCode`, the `Is` method compares `Coder` codes rather than instances, meaning two different `Coder` instances with the same code will be considered equal.
- **`standardErrors.As(err, target interface{}) bool`**: Works as expected. It can be used to extract a `Coder` if an error in the chain embeds one and matches the `Coder` interface, or to extract any other custom error type.

### 8. Formatting Errors (Stack Traces)

Errors created or wrapped by `pkg/errors` functions capture a stack trace. This stack trace is printed when the error is formatted using the `"%+v"` verb with `fmt.Printf` or similar functions.

- The `fundamental` (for `New`, `Errorf`), `wrapper` (for `Wrap`, `Wrapf`), and `withCode` error types within `pkg/errors` all implement the `fmt.Formatter` interface to provide detailed output including stack traces for `"%+v"`.
- `ErrorGroup` also implements `fmt.Formatter` to show details of all its contained errors.

### 9. Predefined `Coder` Instances

The `pkg/errors` module provides several predefined `Coder` instances for common error scenarios. These are exported variables.

| Variable Name          | Code   | HTTP Status | Default Message             |
|------------------------|--------|-------------|-----------------------------|
| `ErrUnknown`           | -1     | 500         | An internal server error occurred |
| `ErrInternalServer`    | 100001 | 500         | Internal server error       |
| `ErrNotFound`          | 100002 | 404         | Not found                   |
| `ErrBadRequest`        | 100003 | 400         | Bad request                 |
| `ErrPermissionDenied`  | 100004 | 403         | Permission denied           |
| `ErrConflict`          | 100005 | 409         | Conflict                    |
| `ErrValidation`        | 100006 | 400         | Validation failed           |
| `ErrResourceUnavailable`|100007 | 503         | Resource unavailable        |
| `ErrTooManyRequests`   | 100008 | 429         | Too many requests           |
| `ErrUnauthorized`      | 100009 | 401         | Unauthorized                |
| `ErrTimeout`           | 100010 | 504         | Timeout                     |
| `ErrNotImplemented`    | 100011 | 501         | Not implemented             |
| `ErrNotSupported`      | 100012 | 501         | Not supported               |
| `ErrAlreadyExists`     | 100013 | 409         | Already exists              |
| `ErrDataLoss`          | 100014 | 500         | Data loss                   |
| `ErrDatabase`          | 100015 | 500         | Database error              |
| `ErrEncoding`          | 100016 | 500         | Encoding error              |
| `ErrDecoding`          | 100017 | 500         | Decoding error              |
| `ErrNetwork`           | 100018 | 500         | Network error               |
| `ErrFilesystem`        | 100019 | 500         | Filesystem error            |
| `ErrConfiguration`     | 100020 | 500         | Configuration error         |
| `ErrAuthentication`    | 100021 | 401         | Authentication failed       |
| `ErrAuthorization`     | 100022 | 403         | Authorization failed        |
| `ErrRateLimitExceeded` | 100023 | 429         | Rate limit exceeded         |
| `ErrInvalidInput`      | 100024 | 400         | Invalid input               |
| `ErrStateMismatch`     | 100025 | 409         | State mismatch              |
| `ErrOperationAborted`  | 100026 | 499         | Operation aborted           |
| `ErrResourceExhausted` | 100027 | 507         | Resource exhausted          |
| `ErrExternalService`   | 100028 | 502         | External service error      |
| `ErrMaintenance`       | 100029 | 503         | Maintenance                 |
| `ErrLogOptionInvalid`  | 300001 | 500         | Invalid log option          |
| `ErrLogRotationSetup`  | 300002 | 500         | Log rotation setup failed   |
| `ErrLogWrite`          | 300003 | 500         | Log write failure           |
| `ErrLogReconfigure`    | 300004 | 500         | Log reconfiguration failed  |
| `ErrLogTargetCreate`   | 300005 | 500         | Log target creation failed  |
| `ErrLogTargetNotSupported`| 300006 | 500     | Log target not supported    |
| `ErrLogBufferFull`     | 300007 | 500         | Log buffer full             |
| `ErrLogRotationDirInvalid`| 300008 | 500     | Invalid log rotation directory|

**Utility functions for Coders:**
- **`IsUnknownCoder(coder Coder) bool`**: Checks if the given `coder` is the predefined `ErrUnknown`.
- **`GetUnknownCoder() Coder`**: Returns the predefined `ErrUnknown` coder.

This specification should provide a good understanding of how to use the `pkg/errors` module effectively.
Refer to the usage guides for more examples and best practices. 