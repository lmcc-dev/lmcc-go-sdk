/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This documentation was collaboratively developed by Martin and AI Assistant.
 */

# Error Handling Examples

[中文版本](README_zh.md)

This directory contains examples demonstrating various features of the error handling module.

## Examples Overview

### [01-basic-errors](01-basic-errors/)
**Basic Error Creation and Formatting**
- Creating simple errors
- Error formatting with %v, %s, %+v
- Error message composition
- Stack trace basics

**Learn**: How to create and format basic errors

### [02-error-wrapping](02-error-wrapping/)
**Error Wrapping and Context**
- Wrapping errors with additional context
- Error chain navigation
- Cause extraction
- Context preservation

**Learn**: How to add context while preserving error information

### [03-error-codes](03-error-codes/)
**Error Codes and Types**
- Using predefined error codes
- Creating custom error codes
- Error code categorization
- HTTP status code mapping

**Learn**: How to implement structured error handling with codes

### [04-stack-traces](04-stack-traces/)
**Stack Trace Capture and Analysis**
- Automatic stack trace capture
- Stack trace formatting
- Performance considerations
- Debugging with stack traces

**Learn**: How to leverage stack traces for debugging

### [05-error-groups](05-error-groups/)
**Error Aggregation and Groups**
- Collecting multiple errors
- Error group operations
- Parallel processing errors
- Error filtering and categorization

**Learn**: How to handle multiple errors in complex operations

## Running the Examples

Each example is self-contained and can be run independently:

```bash
# Navigate to specific example
cd examples/error-handling/01-basic-errors

# Run the example
go run main.go

# Some examples support additional flags
go run main.go --help
```

## Common Patterns Demonstrated

### 1. Basic Error Creation
```go
// Simple error
err := errors.New("operation failed")

// Formatted error
err := errors.Errorf("invalid value: %d", value)

// Error with code
err := errors.WithCode(err, errors.ErrValidation)
```

### 2. Error Wrapping
```go
// Wrap with context
err = errors.Wrap(err, "failed to process user data")

// Wrap with formatted message
err = errors.Wrapf(err, "user %s: operation failed", userID)

// Extract original error
originalErr := errors.Cause(err)
```

### 3. Error Code Usage
```go
// Check error code
if coder := errors.GetCoder(err); coder != nil {
    switch coder.Code() {
    case errors.ErrValidation.Code():
        // Handle validation error
    case errors.ErrNotFound.Code():
        // Handle not found error
    }
}
```

### 4. Stack Trace Handling
```go
// Print detailed stack trace
fmt.Printf("Error: %+v\n", err)

// Get stack trace programmatically
if tracer := errors.GetStackTracer(err); tracer != nil {
    stack := tracer.StackTrace()
    // Process stack frames
}
```

## Best Practices Shown

1. **Error Creation**
   - Use descriptive error messages
   - Include relevant context information
   - Choose appropriate error codes
   - Preserve error chains

2. **Error Wrapping**
   - Add context at each abstraction layer
   - Maintain the original error information
   - Use consistent wrapping patterns
   - Avoid over-wrapping

3. **Error Handling**
   - Check for specific error types/codes
   - Handle errors at appropriate levels
   - Log errors with sufficient context
   - Provide meaningful user feedback

4. **Stack Traces**
   - Use stack traces for debugging
   - Be mindful of performance impact
   - Filter stack frames when needed
   - Include stack traces in error logs

5. **Error Groups**
   - Collect related errors together
   - Provide error summaries
   - Handle partial failures gracefully
   - Use parallel error processing

## Integration Tips

### With Logging Module
```go
// Log errors with context
logger.Errorf("Operation failed: %+v", err)

// Log with error code
if coder := errors.GetCoder(err); coder != nil {
    logger.WithFields(log.Fields{
        "error_code": coder.Code(),
        "error_type": coder.String(),
    }).Error("Structured error logging")
}
```

### With Configuration Module
```go
// Handle config errors specifically
if err := config.LoadConfig(&cfg); err != nil {
    if coder := errors.GetCoder(err); coder != nil {
        switch coder.Code() {
        case errors.ErrConfigFileRead.Code():
            log.Fatal("Configuration file not found or unreadable")
        case errors.ErrConfigValidation.Code():
            log.Fatal("Configuration validation failed")
        }
    }
}
```

## Next Steps

After exploring these error handling examples:

1. Try the [logging-features examples](../logging-features/) for comprehensive logging
2. Explore [integration examples](../integration/) for real-world usage patterns
3. Review [config-features examples](../config-features/) for configuration error handling

## Troubleshooting

### Common Issues

**Issue**: Stack traces not appearing
```
Solution: Ensure errors are created with errors.New() or errors.Errorf()
```

**Issue**: Error codes not recognized
```
Solution: Check if error implements the Coder interface
```

**Issue**: Wrapped errors not unwrapping properly
```
Solution: Use errors.Cause() to extract the root error
```

**Issue**: Performance impact from stack traces
```
Solution: Consider disabling stack traces in production if needed
```

## Related Documentation

- [Error Handling Module Overview](../../docs/usage-guides/errors/en/00_overview.md)
- [Quick Start Guide](../../docs/usage-guides/errors/en/01_quick_start.md)
- [Error Codes Reference](../../docs/usage-guides/errors/en/02_error_codes.md)
- [Stack Traces Guide](../../docs/usage-guides/errors/en/03_stack_traces.md) 