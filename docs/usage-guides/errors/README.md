<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

# `pkg/errors` - Enhanced Error Handling for Go

The `pkg/errors` module in the `lmcc-go-sdk` provides a robust and developer-friendly way to handle errors in Go applications. It builds upon the standard library's error handling capabilities by introducing features like automatic stack traces, error codes, contextual wrapping, and more, aiming to simplify debugging and improve error management.

**[中文版说明 (Chinese Version)](README_zh.md)**

## Overview

Get a high-level understanding of what `pkg/errors` offers and its main benefits.
*   [Module Overview](./en/00_overview.md)

## Quick Start

See a quick example of how to use `pkg/errors` and the immediate advantages it provides.
*   [Quick Start Guide](./en/01_quick_start.md)

## Table of Contents

### 1. Core Concepts
Detailed explanations of the fundamental features.
*   [Creating Errors](./en/02_01_creating_errors.md) - Learn how to create basic errors with `New` and `Errorf`.
*   [Adding Context](./en/02_02_adding_context.md) - Understand how to wrap errors to provide more context using `Wrap` and `Wrapf`.
*   [Using Error Codes](./en/02_03_using_error_codes.md) - Discover how to associate errors with specific codes using `Coder` objects.
*   [Error Checking](./en/02_04_error_checking.md) - Learn how to inspect errors using `standardErrors.Is`, `GetCoder`, and `IsCode`.

### 2. Advanced Features
Explore more advanced capabilities of the module.
*   [Error Aggregation](./en/03_01_error_aggregation.md) - Collect multiple errors into a single `ErrorGroup`.
*   [Stack Traces](./en/03_02_stack_traces.md) - Understand how stack traces are captured and displayed.

### 3. Best Practices
Recommendations for using `pkg/errors` effectively.
*   [Introduction to Best Practices](./en/04_00_best_practices_intro.md)
*   [Define Application-Specific Error Codes](./en/04_01_define_application_specific_error_codes.md)
*   [Wrap Errors for Context](./en/04_02_wrap_errors_for_context.md)
*   [Handle Errors Appropriately at Different Layers](./en/04_03_handle_errors_appropriately.md)
*   [Use `IsCode` for Checking Specific Error Categories](./en/04_04_use_iscode_for_checking.md)
*   [Log Errors Effectively](./en/04_05_log_errors_effectively.md)

### 4. Migration Guide
Guidance on migrating from the standard library's `errors` package.
*   [Introduction to Migration](./en/05_00_migration_intro.md)
*   [Replacing `errors.New` and `fmt.Errorf`](./en/05_01_replacing_new_and_errorf.md)
*   [Wrapping Errors](./en/05_02_wrapping_errors.md)
*   [Error Checking with `errors.Is` and `errors.As`](./en/05_03_error_checking_is_as.md)

### 5. Integration
Examples of how `pkg/errors` can be integrated with other parts of your application.
*   [With Logging Libraries](./en/06_01_with_logging_libraries.md)
*   [With API/HTTP Handlers](./en/06_02_with_api_http_handlers.md)

### 6. Custom Error Codes
Learn how to define and manage your own set of error codes.
*   [Defining and Using Custom Error Codes](./en/07_custom_error_codes.md)

### 7. Troubleshooting
Common issues and how to address them.
*   [Troubleshooting Guide](./en/08_troubleshooting.md)

### 8. Module Specification (API Reference)
A detailed specification of the `pkg/errors` module's public API.
*   [Module Specification](./en/09_module_specification.md)

## Contributing

This module aims to provide a comprehensive yet easy-to-use error handling mechanism. By leveraging its features, developers can build more resilient and maintainable Go applications. 