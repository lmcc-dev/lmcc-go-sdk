/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package errors

import (
	"fmt"
	"io"
	"strings"
)

// ErrorGroup holds a list of errors. It implements the error interface
// and is compatible with Go 1.20's errors.Join mechanics via an Unwrap method.
// ErrorGroup 包含一个错误列表。它实现了 error 接口，并通过 Unwrap 方法与 Go 1.20 的 errors.Join 机制兼容。
type ErrorGroup struct {
	errs    []error
	message string // Optional: An overarching message for the group (主要信息)
	// stack StackTrace // Optional: Stack trace for the creation of the group itself (聚合错误自身的堆栈)
	// TODO: Consider if a stack trace for the group itself is needed, or if relying on individual error stacks is sufficient.
}

// NewErrorGroup creates a new ErrorGroup.
// NewErrorGroup 创建一个新的 ErrorGroup。
//
// Parameters:
//   message: An optional overarching message for the error group. (错误组的可选总体消息。)
//
// Returns:
//   *ErrorGroup: A pointer to the newly created ErrorGroup. (指向新创建的 ErrorGroup 的指针。)
func NewErrorGroup(message ...string) *ErrorGroup {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}
	return &ErrorGroup{
		errs:    []error{},
		message: msg,
		// stack: callers(1, defaultStackLength), // Capture stack if we decide to include it for the group
	}
}

// Add adds a non-nil error to the group. Nil errors are ignored.
// Add 将一个非 nil 错误添加到组中。Nil 错误将被忽略。
//
// Parameters:
//   err: The error to add. (要添加的错误。)
func (eg *ErrorGroup) Add(err error) {
	if err == nil {
		return
	}
	eg.errs = append(eg.errs, err)
}

// Errors returns the list of errors in the group.
// Errors 返回组中的错误列表。
//
// Returns:
//   []error: A slice of errors. (错误切片。)
func (eg *ErrorGroup) Errors() []error {
	// Return a copy to prevent external modification if eg.errs is not pointer-based.
	// However, []error itself is a reference type.
	// To be safe and clear, if direct modification of the internal slice is a concern,
	// one might copy it. For now, direct return is fine as errors are typically immutable or their state is not changed via this.
	return eg.errs
}

// Error implements the error interface. It returns a string representation of the error group.
// Error 实现了 error 接口。它返回错误组的字符串表示形式。
//
// Returns:
//   string: A string describing all errors in the group. (描述组中所有错误的字符串。)
func (eg *ErrorGroup) Error() string {
	if len(eg.errs) == 0 {
		if eg.message != "" {
			return eg.message // Return just the group message if no errors but message exists
		}
		return "no errors in group" // (组内没有错误)
	}

	var b strings.Builder
	if eg.message != "" {
		b.WriteString(eg.message)
		b.WriteString(": ")
	} else {
		// Default prefix if no custom message for the group itself
		// (如果没有为组本身提供自定义消息，则使用默认前缀)
		if len(eg.errs) > 1 {
			b.WriteString("multiple errors occurred: ") // (发生多个错误：)
		} else {
			b.WriteString("an error occurred: ") // (发生一个错误：)
		}
	}

	for i, err := range eg.errs {
		if i > 0 {
			b.WriteString("; ")
		}
		b.WriteString(err.Error())
	}
	return b.String()
}

// Unwrap returns the list of contained errors, making it compatible with errors.Is and errors.As
// for unwrapping multiple errors (Go 1.20+ behavior).
// Unwrap 返回包含的错误列表，使其与 errors.Is 和 errors.As 兼容，
// 用于解包多个错误 (Go 1.20+ 的行为)。
//
// Returns:
//   []error: The slice of errors contained within the group. (组中包含的错误切片。)
//            Returns nil if the group contains no errors, to align with errors.Join behavior
//            where Join(nil...) is nil.
//            (如果组中没有错误，则返回 nil，以与 errors.Join 的行为保持一致，其中 Join(nil...) 为 nil。)
func (eg *ErrorGroup) Unwrap() []error {
	if len(eg.errs) == 0 {
		return nil
	}
	// Return a copy to prevent modification of the internal slice.
	// Although []error is a reference type, the standard library's Join also effectively
	// works with a new slice internally when it constructs the joined error.
	// For consistency and safety, returning a new slice is good practice here.
	// However, the standard library's multi-unwrapper pattern simply returns the slice.
	// Let's stick to returning the direct slice as per typical Unwrap []error pattern,
	// or a copy if modification is a high concern. For now, direct is simple.
	// Reconsidering: errors.Join returns a new error instance. The Unwrap() []error method
	// on such an error would return the slice it was constructed with.
	// We should return eg.errs directly.
	return eg.errs
}

// Format implements fmt.Formatter to provide custom formatting for ErrorGroup.
// Format 实现了 fmt.Formatter 接口，为 ErrorGroup 提供自定义格式化。
// When the verb is 'v' and the '+' flag is used (e.g., "%+v"),
// it prints the group's message (if any) followed by a detailed, multi-line
// representation of each contained error, including their stack traces if available.
// (当动词是 'v' 且使用了 '+' 标志 (例如 "%+v") 时，)
// (它会打印错误组的消息 (如果有)，然后是每个所含错误的详细、多行表示，)
// (如果可用，则包括其堆栈跟踪。)
// For other verbs ('s', 'q') or when '+' is not used with 'v',
// it defaults to the output of the Error() method.
// (对于其他动词 ('s', 'q') 或当 'v' 未与 '+' 一起使用时，它默认为 Error() 方法的输出。)
func (eg *ErrorGroup) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			if eg.message != "" {
				io.WriteString(s, eg.message)
				io.WriteString(s, "\n") // Add a newline after the group message
			}
			if len(eg.errs) == 0 && eg.message == "" { // Handle case where group is empty and has no message
				io.WriteString(s, "empty error group") // (空错误组)
				return
			} else if len(eg.errs) == 0 && eg.message != "" { // Group has message but no errors
				// The message was already printed if it exists.
				// We can add a note that there are no sub-errors if desired.
				// io.WriteString(s, " (contains no sub-errors)")
				return // Avoids printing "Error X of 0"
			}

			for i, err := range eg.errs {
				if i > 0 {
					io.WriteString(s, "\n") // Add a separator line between errors
				}
				// Use Fprintf to format each sub-error with its details using %+v
				// This will recursively call Format on sub-errors if they implement fmt.Formatter
				// (使用 Fprintf 通过 %+v 格式化每个子错误的详细信息)
				// (如果子错误实现了 fmt.Formatter，这将递归调用其 Format 方法)
				fmt.Fprintf(s, "Error %d of %d: %+v", i+1, len(eg.errs), err)
			}
			return
		}
		fallthrough // For '%v' without '+', fall through to '%s'
	case 's':
		io.WriteString(s, eg.Error())
	case 'q':
		fmt.Fprintf(s, "%q", eg.Error())
	}
}

// TODO: (已移除关于 Format 的 TODO, 因为已实现)
// The existing TODO for stack trace on the group itself can remain if that feature is still desired.
// func (eg *ErrorGroup) Format(s fmt.State, verb rune) {
// ... (此处是旧的注释掉的 Format 方法，将被新的实现替换)
// } 