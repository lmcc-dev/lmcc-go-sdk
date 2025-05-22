/*
 * Author: Martin <lmcc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package errors_test // Use errors_test for black-box testing

import (
	"errors" // Standard library errors for Is/As and creating simple errors
	"fmt"
	"strings"
	"testing"

	lmccerrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// lmccErrNotFound is a predefined error for testing `errors.Is` with Coder types.
// It should be an actual error instance created with a Coder.
// 使用 pkg/errors/coder.go中已定义的 ErrNotFound Coder (Using the predefined ErrNotFound Coder from pkg/errors/coder.go)
var lmccErrNotFound = lmccerrors.NewWithCode(lmccerrors.ErrNotFound, "resource was not found for group test")

// MyCustomErrorType is a custom error type for As_Compatibility test.
// MyCustomErrorType 是用于 As_Compatibility 测试的自定义错误类型。
type MyCustomErrorType struct {
	Msg string
	Val int
}

// Error implements the error interface for MyCustomErrorType.
// Error 为 MyCustomErrorType 实现错误接口。
func (mce *MyCustomErrorType) Error() string { return fmt.Sprintf("MyCustomErrorType: %s: %d", mce.Msg, mce.Val) }

// NonExistentErrorType is a custom error type for As_Compatibility test, used to test 'As' for a type not in the group.
// NonExistentErrorType 是用于 As_Compatibility 测试的自定义错误类型，用于测试 'As' 获取组中不存在的类型。
type NonExistentErrorType struct{}

// Error implements the error interface for NonExistentErrorType.
// Error 为 NonExistentErrorType 实现错误接口。
func (ne *NonExistentErrorType) Error() string { return "non-existent" }

func TestNewErrorGroup(t *testing.T) {
	t.Parallel()

	t.Run("WithoutMessage", func(t *testing.T) {
		eg := lmccerrors.NewErrorGroup()
		if eg == nil {
			t.Fatal("NewErrorGroup() returned nil")
		}
		if len(eg.Errors()) != 0 {
			t.Errorf("Expected new group to have 0 errors, got %d", len(eg.Errors()))
		}
		// Default message when empty and Error() is called
		if eg.Error() != "no errors in group" {
		    t.Errorf("Expected default empty message, got %q", eg.Error())
		}
	})

	t.Run("WithMessage", func(t *testing.T) {
		msg := "group initialization failed"
		eg := lmccerrors.NewErrorGroup(msg)
		if eg == nil {
			t.Fatal("NewErrorGroup(message) returned nil")
		}
		if len(eg.Errors()) != 0 {
			t.Errorf("Expected new group with message to have 0 errors, got %d", len(eg.Errors()))
		}
		// Message should be returned by Error() if no sub-errors
		if eg.Error() != msg {
		    t.Errorf("Expected group message %q when empty, got %q", msg, eg.Error())
		}
	})
}

func TestErrorGroup_Add(t *testing.T) {
	t.Parallel()

	eg := lmccerrors.NewErrorGroup("test add")

	err1 := errors.New("first error")
	err2 := errors.New("second error")

	eg.Add(err1)
	if len(eg.Errors()) != 1 {
		t.Fatalf("Expected 1 error after adding one, got %d", len(eg.Errors()))
	}
	if eg.Errors()[0] != err1 {
		t.Errorf("Expected first added error to be %v, got %v", err1, eg.Errors()[0])
	}

	// Add nil error, should be ignored
	eg.Add(nil)
	if len(eg.Errors()) != 1 {
		t.Errorf("Expected 1 error after adding nil, got %d", len(eg.Errors()))
	}

	eg.Add(err2)
	if len(eg.Errors()) != 2 {
		t.Fatalf("Expected 2 errors after adding two, got %d", len(eg.Errors()))
	}
	if eg.Errors()[1] != err2 {
		t.Errorf("Expected second added error to be %v, got %v", err2, eg.Errors()[1])
	}
}

func TestErrorGroup_Error(t *testing.T) {
	t.Parallel()

	err1 := errors.New("DB connection failed")
	err2 := errors.New("validation rule broken")
	lmccErr1 := lmccerrors.NewWithCode(lmccerrors.ErrInternalServer, "internal service unavailable")

	tests := []struct {
		name        string
		groupMsg    string
		errsToAdd   []error
		wantErrorString string
	}{
		{
			name:        "NoErrors_NoGroupMessage",
			groupMsg:    "",
			errsToAdd:   []error{},
			wantErrorString: "no errors in group",
		},
		{
			name:        "NoErrors_WithGroupMessage",
			groupMsg:    "Process A failed",
			errsToAdd:   []error{},
			wantErrorString: "Process A failed",
		},
		{
			name:        "OneError_NoGroupMessage",
			groupMsg:    "",
			errsToAdd:   []error{err1},
			wantErrorString: "an error occurred: DB connection failed",
		},
		{
			name:        "OneError_WithGroupMessage",
			groupMsg:    "Request handling error",
			errsToAdd:   []error{err1},
			wantErrorString: "Request handling error: DB connection failed",
		},
		{
			name:        "MultipleErrors_NoGroupMessage",
			groupMsg:    "",
			errsToAdd:   []error{err1, err2},
			wantErrorString: "multiple errors occurred: DB connection failed; validation rule broken", // Corrected: err2.Error()
		},
		{
			name:        "MultipleErrors_WithGroupMessage",
			groupMsg:    "Task failed with errors",
			errsToAdd:   []error{err1, err2},
			wantErrorString: "Task failed with errors: DB connection failed; validation rule broken",
		},
		{
			name:        "MixOfStdAndLmccErrors_WithGroupMessage",
			groupMsg:    "Operation summary",
			errsToAdd:   []error{err1, lmccErr1, err2},
			wantErrorString: "Operation summary: DB connection failed; Internal server error: internal service unavailable; validation rule broken",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var eg *lmccerrors.ErrorGroup
			if tt.groupMsg == "" {
				eg = lmccerrors.NewErrorGroup()
			} else {
				eg = lmccerrors.NewErrorGroup(tt.groupMsg)
			}

			for _, err := range tt.errsToAdd {
				eg.Add(err)
			}

			if got := eg.Error(); got != tt.wantErrorString {
				t.Errorf("ErrorGroup.Error() = %q, want %q", got, tt.wantErrorString)
			}
		})
	}
	// Manual correction for one test case with err2
	t.Run("MultipleErrors_NoGroupMessage_CorrectedErr2", func(t *testing.T) {
		eg := lmccerrors.NewErrorGroup()
		eg.Add(err1)
		eg.Add(err2)
		want := "multiple errors occurred: DB connection failed; validation rule broken"
		if eg.Error() != want {
		    t.Errorf("ErrorGroup.Error() = %q, want %q", eg.Error(), want)
		}
	})
}

func TestErrorGroup_Unwrap(t *testing.T) {
	t.Parallel()

	err1 := errors.New("error one")
	err2 := lmccerrors.New("error two with lmcc stack")

	t.Run("EmptyGroup_UnwrapNil", func(t *testing.T) {
		eg := lmccerrors.NewErrorGroup()
		if unwrapped := eg.Unwrap(); unwrapped != nil {
			t.Errorf("Unwrap() on empty group should return nil, got %v", unwrapped)
		}
	})

	t.Run("GroupWithErrors_UnwrapReturnsAll", func(t *testing.T) {
		eg := lmccerrors.NewErrorGroup()
		eg.Add(err1)
		eg.Add(err2)

		unwrapped := eg.Unwrap()
		if len(unwrapped) != 2 {
			t.Fatalf("Unwrap() expected to return 2 errors, got %d", len(unwrapped))
		}
		if unwrapped[0] != err1 {
			t.Errorf("Unwrap()[0] = %v, want %v", unwrapped[0], err1)
		}
		if unwrapped[1] != err2 {
			t.Errorf("Unwrap()[1] = %v, want %v", unwrapped[1], err2)
		}
	})

	// Test compatibility with errors.Is and errors.As via Unwrap []error
	t.Run("Is_Compatibility", func(t *testing.T) {
		eg := lmccerrors.NewErrorGroup()
		eg.Add(err1)
		eg.Add(lmccErrNotFound) // lmccErrNotFound is a predefined Coder error

		if !errors.Is(eg, err1) {
			t.Errorf("errors.Is(group, err1) should be true")
		}
		if !errors.Is(eg, lmccErrNotFound) {
			t.Errorf("errors.Is(group, lmccErrNotFound) should be true")
		}
		stdErrNotFound := errors.New("not found std") // Different instance
		if errors.Is(eg, stdErrNotFound) {
			t.Errorf("errors.Is(group, differentErr) should be false")
		}
	})

	t.Run("As_Compatibility", func(t *testing.T) {
		// Define a custom error type for As
		// type MyCustomErrorType struct{
		// 	Msg string
		// 	Val int
		// }
		// func (mce *MyCustomErrorType) Error() string { return fmt.Sprintf("%s: %d", mce.Msg, mce.Val) }

		customErrInstance := &MyCustomErrorType{Msg: "specific issue", Val: 123}
		anotherStdErr := errors.New("another standard error")

		eg := lmccerrors.NewErrorGroup()
		eg.Add(anotherStdErr)
		eg.Add(customErrInstance)

		var targetAs *MyCustomErrorType
		if !errors.As(eg, &targetAs) {
			t.Fatal("errors.As(group, &MyCustomErrorType) should be true")
		}
		if targetAs.Msg != "specific issue" || targetAs.Val != 123 {
			t.Errorf("errors.As extracted wrong custom error: got %+v, want %+v", targetAs, customErrInstance)
		}

		// Try to extract a type that is not in the group
		// type NonExistentErrorType struct{}
		// func (ne *NonExistentErrorType) Error() string { return "non-existent" }
		var targetNonExistent *NonExistentErrorType
		if errors.As(eg, &targetNonExistent) {
			t.Errorf("errors.As(group, &NonExistentErrorType) should be false")
		}
	})
}

// mockError (模拟错误) is a simple error type for testing.
// mockError 是一个用于测试的简单错误类型。
type mockError struct {
	msg string
}

// Error (错误) implements the error interface.
// Error 实现了 error 接口。
func (e *mockError) Error() string {
	return e.msg
}

// mockErrorWithFormatter (模拟带格式化器的错误) is an error type that also implements fmt.Formatter.
// mockErrorWithFormatter 是一个同时实现 fmt.Formatter 的错误类型。
type mockErrorWithFormatter struct {
	msg       string
	detailMsg string
}

// Error (错误) implements the error interface.
// Error 实现了 error 接口。
func (e *mockErrorWithFormatter) Error() string {
	return e.msg
}

// Format (格式化) implements the fmt.Formatter interface.
// Format 实现了 fmt.Formatter 接口。
func (e *mockErrorWithFormatter) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%s - %s (detailed)", e.msg, e.detailMsg)
			return
		}
		fallthrough
	case 's', 'q':
		fmt.Fprintf(s, e.msg)
	}
}

// TestErrorGroup_Format (测试 ErrorGroup 的 Format 方法)
// Tests the custom formatting of ErrorGroup.
// (测试 ErrorGroup 的自定义格式化。)
func TestErrorGroup_Format(t *testing.T) {
	// Test setup
	// (测试设置)
	_ = lmccerrors.New("temporary test call") // Temporary call for diagnostics
	err1 := lmccerrors.New("original error 1") // This will be a *fundamental
	err2 := &mockError{msg: "mock error 2"}
	err3WithCode := lmccerrors.NewWithCode(lmccerrors.ErrOperationFailed, "error with code 3") // This will be a *withCode
	err4Formatted := &mockErrorWithFormatter{msg: "formatted error 4", detailMsg: "details for formatted error 4"}

	// Capture stack trace for err3WithCode to compare later
	// (捕获 err3WithCode 的堆栈跟踪以供后续比较)
	// Note: The actual stack trace will vary based on test execution, so we'll check for its presence
	// and the beginning of the format.
	// (注意：实际的堆栈跟踪会因测试执行而异，因此我们将检查其是否存在以及格式的开头。)
	var err3Stack string
	if fe, ok := err3WithCode.(fmt.Formatter); ok {
		err3Stack = fmt.Sprintf("%+v", fe)
		// We only need the part after the message for comparison, as the message is part of the "Error X of Y" line.
		// (我们只需要消息之后的部分进行比较，因为消息是 "Error X of Y" 行的一部分。)
		parts := strings.SplitN(err3Stack, "\n", 2) // Split by literal \n as stack trace contains it
		if len(parts) > 1 {
				// The message itself part of `err3WithCode.Error()` will be `ErrOperationFailed.String(): error with code 3`
				// The `fmt.Fprintf(s, "Error %d of %d: %+v", i+1, len(eg.errs), err)` in `ErrorGroup.Format`
				// will print this message. Then `err3WithCode.Format` will print its own message again, followed by stack.
				// So we need to find the beginning of the stack trace.
				// Let's look for the first line of the stack.
				err3Stack = parts[1] // Get the part with the stack
		} else {
			err3Stack = "stack trace expected but not fully captured for comparison"
		}
	}


	tests := []struct {
		name          string        // Test case name (测试用例名称)
		group         *lmccerrors.ErrorGroup   // The error group to test (要测试的错误组)
		format        string        // The format string (格式字符串)
		expectedParts []string      // Expected parts of the output string (预期输出字符串的部分内容)
		exactMatch    bool          // Whether to do an exact match or Contains check (是否进行精确匹配或包含检查)
		debugOutput   bool          // Flag to print output for debugging (打印输出以进行调试的标志)
	}{
		{
			name:   "empty group with %+v",
			group:  lmccerrors.NewErrorGroup(""),
			format: "%+v",
			expectedParts: []string{"empty error group"},
			exactMatch: true,
		},
		{
			name:   "empty group with message with %+v",
			group:  lmccerrors.NewErrorGroup("Group with no errors"),
			format: "%+v",
			expectedParts: []string{"Group with no errors\n"}, // Newline added after group message
			exactMatch: true,
		},
		{
			name:   "single error with %+v, no group message",
			group:  func() *lmccerrors.ErrorGroup {
				eg := lmccerrors.NewErrorGroup("")
				eg.Add(err1)
				return eg
			}(),
			format: "%+v",
			// Expected: "Error 1 of 1: " followed by err1's %+v output
			// Since err1 is *fundamental, its %+v will be its message then stack
			expectedParts: []string{
				"Error 1 of 1: original error 1",
				"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors.New", // Adjusted: Stack starts at New
			},
			debugOutput: true, // Enable debug output for this specific sub-test
		},
		{
			name:   "single error with %+v, with group message",
			group:  func() *lmccerrors.ErrorGroup {
				eg := lmccerrors.NewErrorGroup("Group A")
				eg.Add(err2)
				return eg
			}(),
			format: "%+v",
			expectedParts: []string{
				"Group A\n",
				"Error 1 of 1: mock error 2", // mockError doesn't have special %+v
			},
		},
		{
			name: "multiple errors with %+v, with group message",
			group: func() *lmccerrors.ErrorGroup {
				eg := lmccerrors.NewErrorGroup("Group B")
				eg.Add(err1) // fundamental
				eg.Add(err2) // mockError
				eg.Add(err3WithCode) // withCode (includes stack)
				eg.Add(err4Formatted) // mockErrorWithFormatter
				return eg
			}(),
			format: "%+v",
			expectedParts: []string{
				"Group B\n",
				"Error 1 of 4: original error 1", // err1 is *fundamental, so its %+v includes stack
				"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors.New", // Stack from err1 starts at New
				"\nError 2 of 4: mock error 2", // err2 is mockError
				fmt.Sprintf("\nError 3 of 4: %s: %s", lmccerrors.ErrOperationFailed.String(), "error with code 3"), // Corrected: single backslash for newline
				"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors.NewWithCode", // Stack from err3WithCode (wc.stack) starts at NewWithCode
				"\nError 4 of 4: formatted error 4 - details for formatted error 4 (detailed)", // Corrected: single backslash for newline
			},
			debugOutput: true, // Enable debug output for this specific sub-test
		},
		{
			name:   "group with message, no errors, using %s",
			group:  lmccerrors.NewErrorGroup("Group C msg"),
			format: "%s",
			expectedParts: []string{"Group C msg"},
			exactMatch: true,
		},
		{
			name:   "group with one error, using %s",
			group:  func() *lmccerrors.ErrorGroup {
				eg := lmccerrors.NewErrorGroup("Group D")
				eg.Add(err1)
				return eg
			}(),
			format: "%s",
			expectedParts: []string{"Group D: original error 1"},
			exactMatch: true,
		},
		{
			name:   "group with multiple errors, no group message, using %s",
			group:  func() *lmccerrors.ErrorGroup {
				eg := lmccerrors.NewErrorGroup("")
				eg.Add(err1)
				eg.Add(err2)
				return eg
			}(),
			format: "%s",
			expectedParts: []string{"multiple errors occurred: original error 1; mock error 2"},
			exactMatch: true,
		},
		{
			name:   "group with multiple errors, with group message, using %v", // %v falls through to %s
			group:  func() *lmccerrors.ErrorGroup {
				eg := lmccerrors.NewErrorGroup("Group E")
				eg.Add(err1)
				eg.Add(err2)
				return eg
			}(),
			format: "%v",
			expectedParts: []string{"Group E: original error 1; mock error 2"},
			exactMatch: true,
		},
		{
			name:   "group with multiple errors, with group message, using %q", // %q quotes %s
			group:  func() *lmccerrors.ErrorGroup {
				eg := lmccerrors.NewErrorGroup("Group F")
				eg.Add(err1)
				eg.Add(err2)
				return eg
			}(),
			format: "%q",
			expectedParts: []string{`"Group F: original error 1; mock error 2"`},
			exactMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formattedStr := fmt.Sprintf(tt.format, tt.group)
			// ALWAYS LOG FOR NOW TO DEBUG STACK TRACE
			t.Logf("Formatted string for '%s':\n%s", tt.name, formattedStr)

			if tt.debugOutput { // Keep original logic but also log above
				// t.Logf("Formatted string for '%s':\n%s", tt.name, formattedStr) // Already logged above
			}

			if tt.exactMatch {
				assert.Equal(t, strings.Join(tt.expectedParts, ""), formattedStr)
			} else {
				for _, part := range tt.expectedParts {
					assert.Contains(t, formattedStr, part)
				}
			}
		})
	}
} 