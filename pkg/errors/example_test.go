/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package errors_test // Use errors_test package to avoid circular dependencies and test as a user

import (
	stdErrors "errors" // Import standard library errors as stdErrors
	"fmt"

	// "io" // Placeholder, might be needed for more complex examples

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
)

// ExampleNew demonstrates creating a basic error using New.
// ExampleNew 展示了如何使用 New 创建一个基本错误。
func ExampleNew() {
	err := errors.New("something went wrong")
	fmt.Println(err)
	// Output: something went wrong
}

// ExampleErrorf demonstrates creating a basic error with formatting using Errorf.
// ExampleErrorf 展示了如何使用 Errorf 创建一个带格式化的基本错误。
func ExampleErrorf() {
	requestID := "req-123"
	err := errors.Errorf("failed to process request %s", requestID)
	fmt.Println(err)
	// Output: failed to process request req-123
}

// Predefined coder for example purposes
// 用于示例的预定义 Coder
var ErrExampleCoder = errors.NewCoder(90001, 500, "Example Coder Error", "")

// ExampleNewWithCode demonstrates creating an error with a Coder.
// ExampleNewWithCode 展示了如何创建一个带有 Coder 的错误。
func ExampleNewWithCode() {
	err := errors.NewWithCode(ErrExampleCoder, "an operation failed with a specific code")
	fmt.Println(err)

	// To check the Coder
	if c := errors.GetCoder(err); c != nil {
		fmt.Printf("Coder: %s, Code: %d\n", c.String(), c.Code())
	}

	// Example of %+v (output order of stack trace might vary based on test execution)
	// We'll check for presence of key elements instead of exact match for stack.
	// fmt.Printf("%+v\n", err)

	// Output:
	// Example Coder Error: an operation failed with a specific code
	// Coder: Example Coder Error, Code: 90001
}

// ExampleWrap demonstrates wrapping an existing error with an additional message.
// ExampleWrap 展示了如何用附加消息包装现有错误。
func ExampleWrap() {
	originalErr := errors.New("original low-level error")
	wrappedErr := errors.Wrap(originalErr, "contextual information for wrapping")
	fmt.Println(wrappedErr)

	// Example of %+v for wrapped error
	// Stack trace output can be verbose and vary. We'll focus on the error messages.
	// fmt.Printf("%+v\n", wrappedErr)

	// Output: contextual information for wrapping: original low-level error
}

// ExampleWithCode demonstrates attaching a Coder to an existing error.
// ExampleWithCode 展示了如何给现有错误附加一个 Coder。
func ExampleWithCode() {
	originalErr := errors.New("some system error")
	errWithCode := errors.WithCode(originalErr, ErrExampleCoder)
	fmt.Println(errWithCode)

	if c := errors.GetCoder(errWithCode); c != nil {
		fmt.Printf("Coder: %s, Code: %d\n", c.String(), c.Code())
	}
	// Output:
	// Example Coder Error: some system error
	// Coder: Example Coder Error, Code: 90001
}

// ExampleCause demonstrates retrieving the root cause of an error.
// ExampleCause 展示了如何获取错误的根本原因。
func ExampleCause() {
	err1 := errors.New("root cause")
	err2 := errors.Wrap(err1, "intermediate wrapper")
	err3 := errors.Wrap(err2, "outer wrapper")

	fmt.Println(errors.Cause(err3))
	// Output: root cause
}

// ExampleGetCoder demonstrates retrieving a Coder from an error chain.
// ExampleGetCoder 展示了如何从错误链中检索 Coder。
func ExampleGetCoder() {
	nocoderErr := errors.New("no coder here")
	errWithCode := errors.WithCode(nocoderErr, ErrExampleCoder)
	wrappedWithCode := errors.Wrap(errWithCode, "wrapped with code")

	coder := errors.GetCoder(wrappedWithCode)
	if coder != nil {
		fmt.Printf("Found Coder: %s (Code: %d)", coder.String(), coder.Code())
	} else {
		fmt.Println("No Coder found")
	}
	// Output: Found Coder: Example Coder Error (Code: 90001)
}

// ExampleIs demonstrates checking if an error in a chain matches a Coder.
// ExampleIs 展示了如何检查链中的错误是否与 Coder 匹配。
func ExampleIs() {
	baseErr := errors.NewWithCode(ErrExampleCoder, "operation failed")
	wrapped := errors.Wrap(baseErr, "additional context")

	if stdErrors.Is(wrapped, ErrExampleCoder) { // Use stdErrors.Is
		fmt.Println("Error is of type ErrExampleCoder")
	} else {
		fmt.Println("Error is NOT of type ErrExampleCoder")
	}
	// Output: Error is of type ErrExampleCoder
}

// ExampleAs demonstrates extracting a Coder (or other types) from an error.
// ExampleAs 展示了如何从错误中提取 Coder (或其他类型)。
func ExampleAs() {
	baseErr := errors.NewWithCode(ErrExampleCoder, "specific failure")
	wrapped := errors.Wrap(baseErr, "context")

	var extractedCoder errors.Coder
	if stdErrors.As(wrapped, &extractedCoder) { // Use stdErrors.As
		fmt.Printf("Extracted Coder: %s, Code: %d, HTTP: %d", extractedCoder.String(), extractedCoder.Code(), extractedCoder.HTTPStatus())
	} else {
		fmt.Println("Failed to extract Coder")
	}
	// Output: Extracted Coder: Example Coder Error, Code: 90001, HTTP: 500
}

// ExampleErrorGroup demonstrates creating and using an ErrorGroup.
// ExampleErrorGroup 展示了如何创建和使用 ErrorGroup。
func ExampleErrorGroup() {
	err1 := errors.New("first underlying error")
	errWithCode := errors.NewWithCode(ErrExampleCoder, "error with a coder")

	// Create an error group with an overarching message
	eg := errors.NewErrorGroup("Operation failed due to multiple issues")
	eg.Add(err1)
	eg.Add(nil) // Adding nil should be ignored
	eg.Add(errWithCode)

	// Get the combined error message
	fmt.Println(eg.Error())

	// Check for individual errors using stdErrors.Is
	if stdErrors.Is(eg, err1) {
		fmt.Println("ErrorGroup contains: first underlying error")
	}

	// Check for a Coder within the group using stdErrors.Is
	if stdErrors.Is(eg, ErrExampleCoder) {
		fmt.Println("ErrorGroup contains an error matching ErrExampleCoder")
	}

	// Extract a specific error type using stdErrors.As (if a custom type was added)
	// For this example, we'll demonstrate extracting the Coder again via As
	var extractedCoder errors.Coder
	if stdErrors.As(eg, &extractedCoder) {
		fmt.Printf("Extracted Coder via As from ErrorGroup: %s\n", extractedCoder.String())
	}

	// Output:
	// Operation failed due to multiple issues: first underlying error; Example Coder Error: error with a coder
	// ErrorGroup contains: first underlying error
	// ErrorGroup contains an error matching ErrExampleCoder
	// Extracted Coder via As from ErrorGroup: Example Coder Error
}

// Note: For examples involving %+v, the stack trace can be non-deterministic
// due to test execution environment, go version, and file paths.
// It's often better to test for the presence of key functions/files in stack traces
// in unit tests rather than exact matches in example outputs.
// 注意：对于涉及 %+v 的示例，由于测试执行环境、Go 版本和文件路径的不同，
// 堆栈跟踪可能是不确定的。在单元测试中测试堆栈跟踪中关键函数/文件的存在，
// 通常比在示例输出中进行精确匹配更好。 