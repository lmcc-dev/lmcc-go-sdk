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

// ExampleWrapf demonstrates wrapping an existing error with a formatted additional message.
// ExampleWrapf 展示了如何用格式化的附加消息包装现有错误。
func ExampleWrapf() {
	originalErr := errors.New("database connection failed")
	operationID := "op-789"
	wrappedErr := errors.Wrapf(originalErr, "context: operation %s failed during commit", operationID)
	fmt.Println(wrappedErr)

	// Example of %+v for wrapped error to show stack trace
	// fmt.Printf("%+v\\n", wrappedErr)

	// Output: context: operation op-789 failed during commit: database connection failed
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

// ExampleGetCoder_notFound demonstrates GetCoder when no Coder is present.
// ExampleGetCoder_notFound 展示了当错误链中没有 Coder 时 GetCoder 的行为。
func ExampleGetCoder_notFound() {
	err := errors.New("an error without any coder")
	wrappedErr := errors.Wrap(err, "wrapped, but still no coder")

	coder := errors.GetCoder(wrappedErr)
	if coder != nil {
		fmt.Printf("Found Coder: %s (Code: %d)\\n", coder.String(), coder.Code())
	} else {
		fmt.Println("No Coder found in this chain.")
	}
	// Output: No Coder found in this chain.
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

// ExampleIsCode demonstrates checking if an error in a chain matches a Coder's code.
// ExampleIsCode 展示了如何检查错误链中的错误是否与特定 Coder 的代码匹配。
func ExampleIsCode() {
	errWithSpecificCoder := errors.NewWithCode(errors.ErrNotFound, "specific resource not found") // ErrNotFound Code: 100002

	// Another Coder, potentially with the same code but different instance or message
	// For this example, let's use a different Coder but with the same code as ErrNotFound for demonstration.
	// In practice, codes should be unique for distinct semantic meanings.
	// However, IsCode only cares about the integer value of Code().
	// 另一个 Coder，可能具有相同的代码但实例或消息不同
	// 在此示例中，为演示目的，我们使用一个与 ErrNotFound 代码相同但不同的 Coder。
	// 实际上，对于不同的语义，代码应该是唯一的。
	// 但是，IsCode 只关心 Code() 的整数值。
	anotherCoderWithSameCode := errors.NewCoder(errors.ErrNotFound.Code(), 404, "Another not found type", "")

	errWithAnotherCoder := errors.NewWithCode(anotherCoderWithSameCode, "something else not found")

	// Check using IsCode
	// 使用 IsCode 检查
	if errors.IsCode(errWithSpecificCoder, errors.ErrNotFound) {
		fmt.Println("IsCode: errWithSpecificCoder matches ErrNotFound's code.")
	}
	if errors.IsCode(errWithAnotherCoder, errors.ErrNotFound) {
		// This will match because IsCode checks the numerical Code(), and we made them the same.
		// 这会匹配，因为 IsCode 检查数字 Code()，而我们使它们相同。
		fmt.Println("IsCode: errWithAnotherCoder matches ErrNotFound's code (due to same numeric code).")
	}
	if !errors.IsCode(errWithSpecificCoder, errors.ErrBadRequest) { // ErrBadRequest has a different code // ErrBadRequest 有不同的代码
		fmt.Println("IsCode: errWithSpecificCoder does NOT match ErrBadRequest's code.")
	}

	// For comparison with stdErrors.Is:
	// stdErrors.Is checks for the specific instance (or an error that Is() that instance)
	// 与 stdErrors.Is 比较：
	// stdErrors.Is 检查特定实例 (或 Is() 该实例的错误)
	if stdErrors.Is(errWithSpecificCoder, errors.ErrNotFound) {
		fmt.Println("stdErrors.Is: errWithSpecificCoder IS errors.ErrNotFound instance.")
	}
	if !stdErrors.Is(errWithAnotherCoder, errors.ErrNotFound) {
		// This will be true because errWithAnotherCoder does not contain the *instance* errors.ErrNotFound,
		// even if their numerical codes are the same.
		// 这将为 true，因为 errWithAnotherCoder 不包含 *实例* errors.ErrNotFound，即使它们的数字代码相同。
		fmt.Println("stdErrors.Is: errWithAnotherCoder is NOT the errors.ErrNotFound instance.")
	}

	// Output:
	// IsCode: errWithSpecificCoder matches ErrNotFound's code.
	// IsCode: errWithAnotherCoder matches ErrNotFound's code (due to same numeric code).
	// IsCode: errWithSpecificCoder does NOT match ErrBadRequest's code.
	// stdErrors.Is: errWithSpecificCoder IS errors.ErrNotFound instance.
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

// MyCustomErrorForAs is a custom error type for ExampleAs.
// MyCustomErrorForAs 是 ExampleAs 的自定义错误类型。
type MyCustomErrorForAs struct {
	Msg     string
	Details string
}

func (e *MyCustomErrorForAs) Error() string {
	return fmt.Sprintf("%s (Details: %s)", e.Msg, e.Details)
}

// ExampleAs_customType demonstrates extracting a custom concrete error type using stdErrors.As.
// ExampleAs_customType 展示了如何使用 stdErrors.As 提取自定义具体错误类型。
func ExampleAs_customType() {
	originalCustomErr := &MyCustomErrorForAs{Msg: "Custom specific problem", Details: "Value out of range"}
	wrappedErr := errors.Wrap(originalCustomErr, "operation failed due to custom issue")

	var targetCustomError *MyCustomErrorForAs
	if stdErrors.As(wrappedErr, &targetCustomError) {
		fmt.Printf("Extracted Custom Error: Msg='%s', Details='%s'\n", targetCustomError.Msg, targetCustomError.Details)
		// Check if it's the same instance
		// 检查它是否是同一个实例
		if targetCustomError == originalCustomErr {
			fmt.Println("Extracted error is the original instance.")
		}
	} else {
		fmt.Println("Failed to extract MyCustomErrorForAs.")
	}

	// Demonstrate with a non-matching type
	// 使用不匹配的类型进行演示
	var targetCoder errors.Coder
	if !stdErrors.As(wrappedErr, &targetCoder) {
		fmt.Println("Correctly did not extract errors.Coder (as the immediate cause is MyCustomErrorForAs).")
	}

	// Output:
	// Extracted Custom Error: Msg='Custom specific problem', Details='Value out of range'
	// Extracted error is the original instance.
	// Correctly did not extract errors.Coder (as the immediate cause is MyCustomErrorForAs).
}

// ExampleErrorGroup demonstrates creating and using an ErrorGroup.
// ExampleErrorGroup 展示了如何创建和使用 ErrorGroup。
func ExampleErrorGroup() {
	err1 := errors.New("first underlying error from component A")
	errWithCode := errors.NewWithCode(errors.ErrBadRequest, "input validation failed in component B") // Using a predefined Coder // 使用预定义的Coder

	// Create an error group with an overarching message
	// 创建一个带有总体消息的错误组
	eg := errors.NewErrorGroup("Overall operation failed due to multiple issues")
	eg.Add(err1)
	eg.Add(nil) // Adding nil should be ignored // 添加 nil 应被忽略
	eg.Add(errWithCode)

	// Get the combined error message using Error() method
	// 使用 Error() 方法获取组合的错误消息
	fmt.Println("--- Combined Error Message (eg.Error()) ---")
	fmt.Println(eg.Error())

	// Iterate through individual errors using Errors() method
	// 使用 Errors() 方法遍历各个错误
	fmt.Println("\n--- Individual Errors (eg.Errors()) ---")
	for i, err := range eg.Errors() {
		fmt.Printf("Error %d: %v\n", i+1, err)
	}

	// Check for presence of individual errors using stdErrors.Is
	// 使用 stdErrors.Is 检查是否存在单个错误
	fmt.Println("\n--- Checking with stdErrors.Is ---")
	if stdErrors.Is(eg, err1) {
		fmt.Println("stdErrors.Is: ErrorGroup contains: first underlying error from component A")
	}
	// Check for a Coder within the group using stdErrors.Is
	// This works because ErrBadRequest is a specific Coder instance.
	// 使用 stdErrors.Is 检查组内是否存在 Coder (因为 ErrBadRequest 是特定的 Coder 实例)
	if stdErrors.Is(eg, errors.ErrBadRequest) {
		fmt.Println("stdErrors.Is: ErrorGroup contains an error matching ErrBadRequest Coder instance.")
	}

	// Check for a Coder's code within the group using errors.IsCode
	// 使用 errors.IsCode 检查组内是否存在特定错误码
	if errors.IsCode(eg, errors.ErrBadRequest) {
		fmt.Println("errors.IsCode: ErrorGroup contains an error with ErrBadRequest's code.")
	}

	// Extract a specific Coder type using stdErrors.As
	// 使用 stdErrors.As 提取特定的 Coder 类型
	fmt.Println("\n--- Extracting with stdErrors.As ---")
	var extractedCoder errors.Coder
	if stdErrors.As(eg, &extractedCoder) {
		// stdErrors.As with ErrorGroup will return the first error in the group
		// that matches the target type. The order of errors added matters here.
		// In this case, err1 is not a Coder, but errWithCode is. errors.As will check them in order.
		// stdErrors.As 与 ErrorGroup 一起使用时，将返回组中与目标类型匹配的第一个错误。
		// 此处添加错误的顺序很重要。 在这种情况下，err1 不是 Coder，但 errWithCode 是。
		// errors.As 会按顺序检查它们。
		fmt.Printf("stdErrors.As: Extracted Coder: %s (Code: %d)\n", extractedCoder.String(), extractedCoder.Code())
	}

	// Output:
	// --- Combined Error Message (eg.Error()) ---
	// Overall operation failed due to multiple issues: first underlying error from component A; Bad request: input validation failed in component B
	//
	// --- Individual Errors (eg.Errors()) ---
	// Error 1: first underlying error from component A
	// Error 2: Bad request: input validation failed in component B
	//
	// --- Checking with stdErrors.Is ---
	// stdErrors.Is: ErrorGroup contains: first underlying error from component A
	// stdErrors.Is: ErrorGroup contains an error matching ErrBadRequest Coder instance.
	// errors.IsCode: ErrorGroup contains an error with ErrBadRequest's code.
	//
	// --- Extracting with stdErrors.As ---
	// stdErrors.As: Extracted Coder: Bad request (Code: 100003)
}

// ExampleErrorfWithCode demonstrates creating a formatted error with a Coder.
// ExampleErrorfWithCode 展示了如何使用 Coder 创建一个格式化的错误。
func ExampleErrorfWithCode() {
	userID := "user456"
	err := errors.ErrorfWithCode(ErrExampleCoder, "operation failed for user %s", userID)
	fmt.Println(err)

	// To check the Coder
	if c := errors.GetCoder(err); c != nil {
		fmt.Printf("Coder: %s, Code: %d, Reference: %s\n", c.String(), c.Code(), c.Reference())
	}
	// Output:
	// Example Coder Error: operation failed for user user456
	// Coder: Example Coder Error, Code: 90001, Reference:
}

// ExampleIsUnknownCoder demonstrates using IsUnknownCoder.
// ExampleIsUnknownCoder 展示了如何使用 IsUnknownCoder。
func ExampleIsUnknownCoder() {
	err1 := errors.NewWithCode(errors.GetUnknownCoder(), "something truly unknown happened")
	err2 := errors.NewWithCode(errors.ErrBadRequest, "this is a bad request")

	if errors.IsUnknownCoder(errors.GetCoder(err1)) {
		fmt.Println("err1 is an Unknown Coder.")
	}
	if !errors.IsUnknownCoder(errors.GetCoder(err2)) {
		fmt.Println("err2 is NOT an Unknown Coder.")
	}
	if !errors.IsUnknownCoder(nil) {
		fmt.Println("nil is NOT an Unknown Coder.")
	}
	// Output:
	// err1 is an Unknown Coder.
	// err2 is NOT an Unknown Coder.
	// nil is NOT an Unknown Coder.
}

// ExampleGetUnknownCoder demonstrates using GetUnknownCoder.
// ExampleGetUnknownCoder 展示了如何使用 GetUnknownCoder。
func ExampleGetUnknownCoder() {
	unknown := errors.GetUnknownCoder()
	fmt.Printf("Default Unknown Coder: Code=%d, Message='%s'\n", unknown.Code(), unknown.String())

	// You can use this to check if a retrieved Coder is the unknown one
	// 您可以用它来检查检索到的 Coder 是否是未知的 Coder
	retrievedCoder := errors.GetCoder(errors.NewWithCode(errors.GetUnknownCoder(), "some message"))
	if retrievedCoder == errors.GetUnknownCoder() {
		fmt.Println("Retrieved coder is the default unknown coder.")
	}
	// Output:
	// Default Unknown Coder: Code=-1, Message='An internal server error occurred'
	// Retrieved coder is the default unknown coder.
}

// Note: For examples involving %+v, the stack trace can be non-deterministic
// due to test execution environment, go version, and file paths.
// It's often better to test for the presence of key functions/files in stack traces
// in unit tests rather than exact matches in example outputs.
// 注意：对于涉及 %+v 的示例，由于测试执行环境、Go 版本和文件路径的不同，
// 堆栈跟踪可能是不确定的。在单元测试中测试堆栈跟踪中关键函数/文件的存在，
// 通常比在示例输出中进行精确匹配更好。
