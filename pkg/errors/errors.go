/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package errors

import (
	"errors" // Import standard errors package
	"fmt"
	// Ensure runtime is imported for callers()
	// Added for runtime.Callers
)

// skipFrames is the number of frames to skip to get to the caller of the errors package functions.
// skipFrames 是要跳过的帧数，以到达 errors 包函数的调用者。
const skipFrames = 2 // 0: runtime.Callers, 1: errors.callers, 2: New/Errorf/Wrap/Wrapf/etc.

// errFundamentalTestSentinel is a sentinel error for testing fundamental.Unwrap behavior.
// errFundamentalTestSentinel 是用于测试 fundamental.Unwrap 行为的哨兵错误。
// var errFundamentalTestSentinel = errors.New("fundamental unwrap test sentinel") // Reverted

// fundamental is an error that has a message and a stack trace.
// fundamental 是一个包含消息和堆栈跟踪的错误。
type fundamental struct {
	// msg is the error message.
	// msg 是错误消息。
	msg string

	// stack is the stack trace from the point where the error was created.
	// stack 是从错误创建点开始的堆栈跟踪。
	stack StackTrace
}

// Error returns the message of the fundamental error.
// Error 返回 fundamental 错误的消息。
func (f *fundamental) Error() string {
	return f.msg
}

// Unwrap returns nil for a fundamental error, as it does not wrap another error.
// Unwrap 为 fundamental 错误返回 nil，因为它不包装另一个错误。
func (f *fundamental) Unwrap() error {
	return nil // Original behavior restored
}

// Format implements the fmt.Formatter interface for fundamental errors.
// Format 为 fundamental 错误实现 fmt.Formatter 接口。
//
// Supported verbs:
//
//	%s, %v: Print the error message. (打印错误消息。)
//	%+v:    Print the error message followed by the stack trace. (打印错误消息和堆栈跟踪。)
func (f *fundamental) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			// %+v: message and stack trace
			// %+v: 消息和堆栈跟踪
			fmt.Fprint(s, f.msg)
			f.stack.Format(s, verb) // Delegate to StackTrace's Formatter
			return
		}
		fallthrough
	case 's':
		// %s, %v: message only
		// %s, %v: 仅消息
		fmt.Fprint(s, f.msg)
	}
}

// New creates a new fundamental error with the given text message.
// New 使用给定的文本消息创建一个新的 fundamental 错误。
// It captures the stack trace at the point of creation.
// 它在创建点捕获堆栈跟踪。
func New(text string) error {
	return &fundamental{
		msg:   text,
		stack: callers(skipFrames), // skip New itself and runtime.Callers
	}
}

// Errorf creates a new fundamental error with a formatted message.
// Errorf 使用格式化的消息创建一个新的 fundamental 错误。
// It captures the stack trace at the point of creation.
// 它在创建点捕获堆栈跟踪。
func Errorf(format string, args ...interface{}) error {
	return &fundamental{
		msg:   fmt.Sprintf(format, args...),
		stack: callers(skipFrames), // skip Errorf itself and runtime.Callers
	}
}

// Is checks if the fundamental error is equivalent to the target error.
// Is 检查 fundamental 错误是否等同于目标错误。
// For fundamental errors, this primarily means checking if the target is also a fundamental
// error with the same message. Or if target is a Coder, it will not match by default.
// 对于 fundamental 错误，这主要意味着检查目标是否也是具有相同消息的 fundamental 错误。
// 或者如果 target 是一个 Coder，默认情况下它不会匹配。
func (f *fundamental) Is(target error) bool {
	if target == nil {
		return false
	}
	targetF, ok := target.(*fundamental)
	if ok {
		return f.msg == targetF.msg
	}
	// A fundamental error doesn't inherently carry a Coder for direct Is comparison
	// unless we decide to change its structure or Is logic for Coders specifically.
	return false
}

// As checks if the fundamental error can be represented as the target type.
// As 检查 fundamental 错误是否可以表示为目标类型。
func (f *fundamental) As(target interface{}) bool {
	if target == nil {
		return false
	}
	if fundamentalTarget, ok := target.(**fundamental); ok {
		*fundamentalTarget = f
		return true
	}
	// No Coder to extract from a plain fundamental error.
	return false
}

// wrapper is an error that wraps another error, adding a message and a stack trace.
// wrapper 是一个包装另一个错误的错误，添加了消息和堆栈跟踪。
type wrapper struct {
	// msg is the message for this error wrapper.
	// msg 是此错误包装器的消息。
	msg string

	// cause is the underlying error that is being wrapped.
	// cause 是被包装的底层错误。
	cause error

	// stack is the stack trace from the point where the error was wrapped.
	// stack 是从错误包装点开始的堆栈跟踪。
	stack StackTrace
}

// Error returns the message of the wrapper and the underlying error.
// Error 返回包装器及其底层错误的消息。
func (w *wrapper) Error() string {
	// We need to handle the case where cause is nil, although Wrap/Wrapf should prevent this.
	// 我们需要处理 cause 为 nil 的情况，尽管 Wrap/Wrapf 应该防止这种情况。
	if w.cause == nil {
		return w.msg
	}
	return w.msg + ": " + w.cause.Error()
}

// Unwrap returns the underlying error for compatibility with errors.Is and errors.As.
// Unwrap 返回底层错误，以兼容 errors.Is 和 errors.As。
func (w *wrapper) Unwrap() error {
	return w.cause
}

// Format implements the fmt.Formatter interface for wrapper errors.
// Format 为 wrapper 错误实现 fmt.Formatter 接口。
//
// Supported verbs:
//
//	%s, %v: Print the wrapper's message and the underlying error's message. (打印包装器的消息和底层错误的消息。)
//	        Format: "wrapper.msg: cause.Error()"
//	%+v:    Print the wrapper's message, the underlying error's message, and the wrapper's stack trace.
//	        (打印包装器的消息、底层错误的消息以及包装器的堆栈跟踪。)
func (w *wrapper) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			// %+v: message, cause, and stack trace for this wrapper
			// %+v: 此包装器的消息、原因和堆栈跟踪
			fmt.Fprint(s, w.Error()) // Prints "msg: cause.Error()"
			w.stack.Format(s, verb)  // Delegate to StackTrace's Formatter for this wrapper's stack
			return
		}
		fallthrough
	case 's':
		// %s, %v: message and cause only
		// %s, %v: 仅消息和原因
		fmt.Fprint(s, w.Error()) // Prints "msg: cause.Error()"
	}
}

// Wrap annotates err with a new message and a stack trace.
// Wrap 使用新消息和堆栈跟踪来注解错误 err。
// If err is nil, Wrap returns nil.
// 如果 err 为 nil，Wrap 返回 nil。
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return &wrapper{
		msg:   message,
		cause: err,
		stack: callers(skipFrames), // skip Wrap itself and runtime.Callers
	}
}

// Wrapf annotates err with a new formatted message and a stack trace.
// Wrapf 使用新的格式化消息和堆栈跟踪来注解错误 err。
// If err is nil, Wrapf returns nil.
// 如果 err 为 nil，Wrapf 返回 nil。
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &wrapper{
		msg:   fmt.Sprintf(format, args...),
		cause: err,
		stack: callers(skipFrames), // skip Wrapf itself and runtime.Callers
	}
}

// Is checks if the wrapper error or its cause is equivalent to the target error.
// Is 检查 wrapper 错误或其 cause 是否等同于目标错误。
func (w *wrapper) Is(target error) bool {
	if target == nil {
		return false
	}
	// Check if the wrapper itself matches (e.g. if target is a *wrapper with same msg)
	targetW, ok := target.(*wrapper)
	if ok {
		// This simple comparison might not be what's usually desired for Is.
		// errors.Is will typically use this method to unwrap.
		// A direct match of wrapper messages is less common for Is.
		return w.msg == targetW.msg && errors.Is(w.cause, targetW.cause)
	}
	// Delegate to the cause
	return errors.Is(w.cause, target)
}

// As checks if the wrapper error or its cause can be represented as the target type.
// As 检查 wrapper 错误或其 cause 是否可以表示为目标类型。
func (w *wrapper) As(target interface{}) bool {
	if target == nil {
		return false
	}
	if wrapperTarget, ok := target.(**wrapper); ok {
		*wrapperTarget = w
		return true
	}
	// Delegate to the cause
	return errors.As(w.cause, target)
}

// withCode is an error that associates a Coder with an underlying error.
// withCode 是一个将 Coder 与底层错误关联起来的错误。
type withCode struct {
	// cause is the underlying error.
	// cause 是底层错误。
	cause error

	// coder is the Coder associated with this error.
	// coder 是与此错误关联的 Coder。
	coder Coder

	// stack is the stack trace from the point where the Coder was attached.
	// stack 是从附加 Coder 的点开始的堆栈跟踪。
	stack StackTrace
}

// Error returns a string representation of the error, including the Coder's message.
// Error 返回错误的字符串表示，包括 Coder 的消息。
// Format: "coder.String(): cause.Error()"
// 格式："coder.String(): cause.Error()"
func (wc *withCode) Error() string {
	if wc.cause == nil {
		// This case should ideally not happen if constructors are used correctly.
		// 如果构造函数使用正确，理想情况下不应发生这种情况。
		if wc.coder != nil {
			return wc.coder.String()
		}
		return ""
	}
	// Handle nil coder or empty coder string
	// If the coder is nil or its string representation is empty, just return the cause's error string.
	// 如果 coder 为 nil，或者其字符串表示为空，则仅返回 cause 的错误字符串。
	if wc.coder == nil || wc.coder.String() == "" {
		return wc.cause.Error()
	}
	// Otherwise, include the coder string.
	// 否则，包含 coder 字符串。
	return wc.coder.String() + ": " + wc.cause.Error()
}

// Unwrap returns the underlying error.
// Unwrap 返回底层错误。
func (wc *withCode) Unwrap() error {
	return wc.cause
}

// Coder returns the Coder associated with this error.
// Coder 返回与此错误关联的 Coder。
func (wc *withCode) Coder() Coder {
	return wc.coder
}

// Is checks if the withCode error or its cause is equivalent to the target error.
// Is 检查 withCode 错误或其 cause 是否等同于目标错误。
// It gives priority to Coder comparison if the target is a Coder.
// 如果目标是 Coder，它优先进行 Coder 比较。
func (wc *withCode) Is(target error) bool {
	if target == nil {
		return false
	}

	// Attempt to compare by Coder first if the target has a Coder.
	// 如果目标具有 Coder，则首先尝试按 Coder 进行比较。
	targetAsCoder, ok := target.(Coder) // Simpler: target is directly a Coder
	if ok && wc.coder != nil {
		return wc.coder.Code() == targetAsCoder.Code()
	}

	// If target is an error that has a Coder() method (like our withCode or a Coder itself that is also an error)
	// (如果目标是具有 Coder() 方法的错误（例如我们的 withCode 或本身也是错误的 Coder）)
	if typeWithCoder, ok := target.(interface{ Coder() Coder }); ok {
		if wc.coder != nil && typeWithCoder.Coder() != nil {
			return wc.coder.Code() == typeWithCoder.Coder().Code()
		}
	}

	// If direct Coder comparison isn't applicable or doesn't match, fall back to errors.Is on the cause.
	// 如果直接 Coder 比较不适用或不匹配，则回退到对 cause 使用 errors.Is。
	return errors.Is(wc.cause, target)
}

// As checks if the withCode error or its cause can be represented as the target type.
// As 检查 withCode 错误或其 cause 是否可以表示为目标类型。
// It also allows extracting the Coder if the target is of type *Coder.
// 如果目标类型为 *Coder，它还允许提取 Coder。
func (wc *withCode) As(target interface{}) bool {
	if target == nil {
		return false
	}

	// Check if target is *Coder
	// 检查目标是否为 *Coder
	if targetCoder, ok := target.(*Coder); ok {
		*targetCoder = wc.coder
		return true // Successfully assigned the Coder
	}

	// Check if target is **withCode (to extract the withCode instance itself)
	// _检查目标是否为_ **withCode （以提取 withCode 实例本身）
	if withCodeTarget, ok := target.(**withCode); ok {
		*withCodeTarget = wc
		return true
	}

	// Delegate to the cause for other types
	// 对于其他类型，委托给 cause
	return errors.As(wc.cause, target)
}

// Format implements the fmt.Formatter interface for withCode errors.
// Format 为 withCode 错误实现 fmt.Formatter 接口。
//
// Supported verbs:
//
//	%s, %v: Print the error message, including Coder string. (打印错误消息，包括 Coder 字符串。)
//	        Format: "coder.String(): cause.Error()" or "cause.Error()" if Coder or its string is empty.
//	%+v:    Print the error message (as above) and the stack trace of where the code was attached.
//	        (打印错误消息（如上）以及附加代码处的堆栈跟踪。)
func (wc *withCode) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			// %+v: message (including Coder) and stack trace for this withCode wrapper
			// %+v: 此 withCode 包装器的消息（包括 Coder）和堆栈跟踪
			fmt.Fprint(s, wc.Error()) // Prints "coder.String(): cause.Error()" or similar
			wc.stack.Format(s, verb)  // Delegate to StackTrace's Formatter for this withCode's stack
			return
		}
		fallthrough
	case 's':
		// %s, %v: message (including Coder) only
		// %s, %v: 仅消息（包括 Coder）
		fmt.Fprint(s, wc.Error())
	}
}

// NewWithCode creates a new error that associates a Coder with a message.
// NewWithCode 创建一个将 Coder 与消息关联的新错误。
// The underlying error is a new fundamental error created from the text.
// 底层错误是根据文本创建的新的 fundamental 错误。
func NewWithCode(coder Coder, text string) error {
	if coder == nil {
		coder = unknownCoder // Default to unknownCoder if nil Coder is provided
	}
	return &withCode{
		cause: &fundamental{
			msg: text,
			// No separate stack for fundamental here, stack is for withCode
			// fundamental 在这里没有单独的堆栈，堆栈用于 withCode
		},
		coder: coder,
		stack: callers(skipFrames), // skip NewWithCode itself and runtime.Callers
	}
}

// ErrorfWithCode creates a new error that associates a Coder with a formatted message.
// ErrorfWithCode 创建一个将 Coder 与格式化消息关联的新错误。
// The underlying error is a new fundamental error created from the formatted message.
// 底层错误是根据格式化消息创建的新的 fundamental 错误。
func ErrorfWithCode(coder Coder, format string, args ...interface{}) error {
	if coder == nil {
		coder = unknownCoder // Default to unknownCoder if nil Coder is provided
	}
	return &withCode{
		cause: &fundamental{
			msg: fmt.Sprintf(format, args...),
			// No separate stack for fundamental here, stack is for withCode
		},
		coder: coder,
		stack: callers(skipFrames), // skip ErrorfWithCode itself and runtime.Callers
	}
}

// WithCode annotates an existing error with a Coder.
// WithCode 使用 Coder 注解现有错误。
// If err is nil, it returns nil.
// 如果 err 为 nil，则返回 nil。
// If coder is nil, it defaults to unknownCoder.
// 如果 coder 为 nil，则默认为 unknownCoder。
func WithCode(err error, coder Coder) error {
	if err == nil {
		return nil // If err is nil, return nil as per test expectation
	}
	if coder == nil {
		coder = unknownCoder
	}
	// If err is nil, create a new error with the coder and an empty message.
	// This ensures that we always return an error that has the specified Coder.
	// 如果 err 为 nil，则使用 coder 和空消息创建一个新错误。
	// 这确保我们始终返回具有指定 Coder 的错误。
	// return &withCode{ // Previous logic that caused test failure
	// 	cause: &fundamental{msg: ""}, // Use an empty fundamental error as cause
	// 	coder: coder,
	// 	stack: callers(skipFrames), // skip WithCode itself and runtime.Callers
	// }
	return &withCode{
		cause: err,
		coder: coder,
		stack: callers(skipFrames), // skip WithCode itself and runtime.Callers
	}
}

// Cause returns the underlying cause of the error, if possible.
// Cause 返回错误的根本原因（如果可能）。
// An error is considered to have a cause if it implements the causer interface.
// 如果错误实现了 causer 接口，则认为它具有原因。
// If the error does not have a cause, the error itself is returned.
// 如果错误没有原因，则返回错误本身。
func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer) // Check if err implements causer (our internal type)
		if !ok {
			// If not our causer, check for the standard library's Unwrap method
			// 如果不是我们的 causer，则检查标准库的 Unwrap 方法
			unwrapper, okUnwrap := err.(interface{ Unwrap() error })
			if !okUnwrap {
				break // No more unwrapping possible
			}
			unwrappedErr := unwrapper.Unwrap()
			if unwrappedErr == nil {
				break // Reached the end of unwrapping chain
			}
			err = unwrappedErr
			continue
		}
		// If it is our causer, get its cause
		// 如果是我们的 causer，则获取其原因
		c := cause.Cause()
		if c == nil {
			// This should not happen if Cause() is implemented correctly (e.g. not returning nil cause)
			// 如果 Cause() 实现正确（例如，不返回 nil 原因），则不应发生这种情况
			break
		}
		err = c
	}
	return err
}

// GetCoder recursively unwraps an error and returns the first Coder found.
// GetCoder 递归地解包错误并返回找到的第一个 Coder。
// If no Coder is found, it returns nil.
// 如果未找到 Coder，则返回 nil。
func GetCoder(err error) Coder {
	if err == nil {
		return nil
	}

	type coderError interface {
		Coder() Coder
		error // Ensure it is an error
	}

	var currentErr = err
	for currentErr != nil {
		// Check if the current error itself implements Coder() Coder
		// 检查当前错误本身是否实现 Coder() Coder
		if ce, ok := currentErr.(coderError); ok {
			if c := ce.Coder(); c != nil {
				return c // Found a Coder
			}
		}

		// Check if the current error can be unwrapped
		// 检查当前错误是否可以解包
		unwrapper, okUnwrap := currentErr.(interface{ Unwrap() error })
		if !okUnwrap {
			break // No more unwrapping possible
		}
		currentErr = unwrapper.Unwrap()
	}

	return nil // No Coder found in the chain
}

// WithMessage annotates err with a new message.
// WithMessage 使用新消息注解错误 err。
// If err is nil, WithMessage returns nil.
// 如果 err 为 nil，WithMessage 返回 nil。
// This is a convenience function and is equivalent to Wrap(err, message).
// 这是一个方便的函数，等同于 Wrap(err, message)。
func WithMessage(err error, message string) error {
	return Wrap(err, message)
}

// WithMessagef annotates err with a new formatted message.
// WithMessagef 使用新的格式化消息注解错误 err。
// If err is nil, WithMessagef returns nil.
// 如果 err 为 nil，WithMessagef 返回 nil。
// This is a convenience function and is equivalent to Wrapf(err, format, args...).
// 这是一个方便的函数，等同于 Wrapf(err, format, args...)。
func WithMessagef(err error, format string, args ...interface{}) error {
	return Wrapf(err, format, args...)
}

// ... existing code ...
// ... existing code ...
