# lmcc-go-sdk Errors Module Usage Guide

**Authors:** Martin, AI Assistant

[中文版本 (Chinese Version)](USAGE_zh.md) | [Detailed Specification](SPECIFICATION.md)

## Overview

The `pkg/errors` module provides a powerful and structured approach to error handling in Go applications. It extends the standard library with:

- **Automatic stack traces** for easier debugging
- **Error codes** for programmatic error handling
- **Rich error context** through wrapping
- **Error aggregation** for collecting multiple failures
- **Standard library compatibility** with `errors.Is`, `errors.As`

## Quick Start

Replace your existing error handling:

```go
// Instead of:
import "errors"
err := errors.New("operation failed")

// Use:
import "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
err := errors.New("operation failed") // Now with stack trace!
```

## Basic Usage

### 1. Creating Errors

```go
// Simple error with stack trace
err := errors.New("connection failed")

// Formatted error with stack trace
err := errors.Errorf("failed to process user %s", userID)
```

### 2. Adding Context

```go
// Wrap errors to add context
dbErr := db.Connect()
if dbErr != nil {
    return errors.Wrap(dbErr, "failed to initialize database")
}

// Formatted wrapping
return errors.Wrapf(dbErr, "database connection failed for host %s", host)
```

### 3. Using Error Codes

```go
// Using predefined error codes
err := errors.NewWithCode(errors.ErrNotFound, "user not found")

// Creating custom error codes
var ErrInvalidInput = errors.NewCoder(40001, 400, "Invalid input", "")
err := errors.NewWithCode(ErrInvalidInput, "email format is invalid")
```

### 4. Error Checking

```go
// Check for specific error types
if errors.Is(err, errors.ErrNotFound) {
    // Handle not found case
}

// Extract error codes
if coder := errors.GetCoder(err); coder != nil {
    fmt.Printf("Error code: %d, HTTP status: %d", coder.Code(), coder.HTTPStatus())
}

// Check by error code
if errors.IsCode(err, errors.ErrValidation) {
    // Handle validation errors
}
```

## Advanced Usage

### Error Aggregation

Collect multiple errors into a single error:

```go
eg := errors.NewErrorGroup("Validation failed")

// Add individual errors
if username == "" {
    eg.Add(errors.New("username is required"))
}
if email == "" {
    eg.Add(errors.New("email is required"))
}

// Check if any errors occurred
if len(eg.Errors()) > 0 {
    return eg // Returns combined error message
}
```

### Stack Traces

Get detailed stack traces for debugging:

```go
// Print detailed error with stack trace
fmt.Printf("%+v\n", err)

// Example output:
// failed to save user: database connection failed
// github.com/yourapp/pkg/user.Save
//     /path/to/your/user.go:42
// github.com/yourapp/cmd/api.handleCreateUser
//     /path/to/your/api.go:123
```

### Error Chain Navigation

```go
// Get the root cause
rootErr := errors.Cause(wrappedErr)

// Extract specific error types
var validationErr *ValidationError
if errors.As(err, &validationErr) {
    // Handle validation error specifically
}
```

## Predefined Error Codes

The module includes common error codes for immediate use:

### General Errors
- `ErrInternalServer` (100001) - HTTP 500
- `ErrNotFound` (100002) - HTTP 404  
- `ErrBadRequest` (100003) - HTTP 400
- `ErrUnauthorized` (100004) - HTTP 401
- `ErrForbidden` (100005) - HTTP 403
- `ErrValidation` (100006) - HTTP 400
- `ErrTimeout` (100007) - HTTP 504
- `ErrTooManyRequests` (100008) - HTTP 429
- `ErrOperationFailed` (100009) - HTTP 500

### Config Package Errors (200001-200006)
- `ErrConfigFileRead`, `ErrConfigSetup`, `ErrConfigEnvBind`, etc.

### Log Package Errors (300001-300008) 
- `ErrLogInternal`, `ErrLogOptionInvalid`, `ErrLogReconfigure`, etc.

## Best Practices

### 1. Use Error Codes for Categorization

```go
// Good: Use specific error codes
return errors.NewWithCode(errors.ErrValidation, "invalid email format")

// Avoid: Generic errors without context
return errors.New("invalid input")
```

### 2. Add Context When Wrapping

```go
// Good: Add meaningful context
return errors.Wrap(err, "failed to create user account")

// Avoid: Redundant wrapping
return errors.Wrap(err, "error occurred")
```

### 3. Check Errors by Type, Not String

```go
// Good: Check by error code
if errors.IsCode(err, errors.ErrNotFound) {
    return http.StatusNotFound, "Resource not found"
}

// Avoid: String matching
if strings.Contains(err.Error(), "not found") {
    // Brittle and unreliable
}
```

### 4. Use ErrorGroup for Multiple Failures

```go
// Good: Collect validation errors
func ValidateUser(user *User) error {
    eg := errors.NewErrorGroup("User validation failed")
    
    if user.Email == "" {
        eg.Add(errors.NewWithCode(errors.ErrValidation, "email is required"))
    }
    if user.Age < 0 {
        eg.Add(errors.NewWithCode(errors.ErrValidation, "age must be positive"))
    }
    
    if len(eg.Errors()) > 0 {
        return eg
    }
    return nil
}
```

## Migration from Standard Library

### Simple Migration

```go
// Before
import "errors"
err := errors.New("something failed")

// After  
import "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
err := errors.New("something failed") // Now with stack trace!
```

### Enhanced Migration

```go
// Before
import "fmt"
return fmt.Errorf("failed to process %s: %w", id, err)

// After - More structured
return errors.Wrapf(err, "failed to process %s", id)

// Or with error codes
return errors.WithCode(err, errors.ErrOperationFailed)
```

## Integration Examples

### HTTP Handler Error Handling

```go
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
    user, err := h.parseUser(r)
    if err != nil {
        h.handleError(w, errors.WithCode(err, errors.ErrBadRequest))
        return
    }
    
    if err := h.userService.Create(user); err != nil {
        h.handleError(w, err)
        return
    }
    
    w.WriteHeader(http.StatusCreated)
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
    // Extract error code for HTTP status
    if coder := errors.GetCoder(err); coder != nil {
        w.WriteHeader(coder.HTTPStatus())
        json.NewEncoder(w).Encode(map[string]interface{}{
            "error": coder.String(),
            "code": coder.Code(),
        })
        return
    }
    
    // Fallback for unknown errors
    w.WriteHeader(http.StatusInternalServerError)
    json.NewEncoder(w).Encode(map[string]string{
        "error": "Internal server error",
    })
}
```

### Logging Integration

*Note: The `log` package used in the example below (e.g., `log.WithFields`, `log.Error`, `log.Debugf`, `log.GetLevel()`, `log.DebugLevel`) is a placeholder representing your project's chosen logging library (e.g., Logrus, Zap, or the `lmcc-go-sdk/pkg/log` module). You will need to adapt this example to your specific logging setup.*

```go
func logError(err error) {
    if coder := errors.GetCoder(err); coder != nil {
        log.WithFields(log.Fields{
            "error_code": coder.Code(),
            "http_status": coder.HTTPStatus(),
            "reference": coder.Reference(),
        }).Error(err.Error())
    } else {
        log.Error(err.Error())
    }
    
    // Log full stack trace in debug mode
    if log.GetLevel() == log.DebugLevel {
        log.Debugf("Stack trace: %+v", err)
    }
}
```

## Custom Error Codes

Define your own error codes for domain-specific errors:

```go
// Define custom error codes
var (
    ErrUserEmailExists = errors.NewCoder(50001, 409, "User email already exists", "")
    ErrInvalidPassword = errors.NewCoder(50002, 400, "Password does not meet requirements", "")
    ErrAccountLocked   = errors.NewCoder(50003, 423, "Account is locked", "")
)

// Use them in your application
func (s *UserService) CreateUser(email, password string) error {
    if s.userExists(email) {
        return errors.NewWithCode(ErrUserEmailExists, fmt.Sprintf("user with email %s already exists", email))
    }
    
    if !s.isValidPassword(password) {
        return errors.NewWithCode(ErrInvalidPassword, "password must be at least 8 characters")
    }
    
    // ... rest of creation logic
    return nil
}
```

## Troubleshooting

### Stack Traces Not Showing?

Make sure you're using `%+v` format specifier:

```go
// This shows only the error message
fmt.Printf("%v\n", err)

// This shows the full stack trace
fmt.Printf("%+v\n", err)
```

### Error Codes Not Being Detected?

Ensure you're using the error code checking functions:

```go
// Correct way to check error codes
if errors.IsCode(err, errors.ErrNotFound) {
    // handle not found
}

// This won't work as expected
if err == errors.ErrNotFound {
    // This compares different error instances
}
```

### Performance Considerations

Stack trace collection has minimal overhead, but if you're in a high-performance scenario:

- Use error codes to avoid string comparisons
- Consider error pooling for frequently created errors
- Use `ErrorGroup` to batch multiple errors instead of creating many individual errors

For detailed API documentation and complete specifications, see [SPECIFICATION.md](SPECIFICATION.md).