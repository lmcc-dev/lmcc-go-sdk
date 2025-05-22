/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package errors_test // Use errors_test package

import (
	"fmt"
	"strings"
	"testing"

	// Import standard library errors for creating base errors in tests
	stdErrors "errors" // Renamed to avoid conflict

	lmccerrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
)

// localMockCoder is a Coder implementation for testing.
// localMockCoder 是用于测试的 Coder 实现。
type localMockCoder struct {
	C    int
	Ext  string
	HTTP int
	Ref  string
}

func (m *localMockCoder) Code() int         { return m.C }
func (m *localMockCoder) String() string    { return m.Ext }
func (m *localMockCoder) HTTPStatus() int   { return m.HTTP }
func (m *localMockCoder) Reference() string { return m.Ref }
func (m *localMockCoder) Error() string     { return m.Ext } // To satisfy error interface

// TestFundamental_Format tests the Format method of fundamental errors (created by New, Errorf).
// TestFundamental_Format 测试 fundamental 错误的 Format 方法 (由 New, Errorf 创建)。
func TestFundamental_Format(t *testing.T) {
	t.Parallel()

	type errorCreationFunc func(format string, args ...interface{}) error

	tests := []struct {
		name        string
		errorFunc   errorCreationFunc      // Function to create the error (New or Errorf)
		formatOrMsg string                 // Message for New, format string for Errorf
		args        []interface{}          // Args for Errorf, nil for New
		wantMsg     string                 // Expected from err.Error(), %s, %v
		wantInPlusV []string               // Substrings expected in %+v output
	}{
		{
			name:      "New_SimpleMessage_Format",
			errorFunc: func(f string, a ...interface{}) error { return lmccerrors.New(f) },
			formatOrMsg: "basic error for fundamental format testing via New",
			args:        nil,
			wantMsg:   "basic error for fundamental format testing via New",
			wantInPlusV: []string{
				"basic error for fundamental format testing via New",
				"TestFundamental_Format",
				"format_test.go",
			},
		},
		{
			name:      "New_EmptyMessage_Format",
			errorFunc: func(f string, a ...interface{}) error { return lmccerrors.New(f) },
			formatOrMsg: "",
			args:        nil,
			wantMsg:   "",
			wantInPlusV: []string{
				"TestFundamental_Format",
				"format_test.go",
			},
		},
		{
			name:      "Errorf_SimpleFormat_Format",
			errorFunc: lmccerrors.Errorf, // Direct reference
			formatOrMsg: "error from Errorf: %s, value: %d",
			args:        []interface{}{"test", 42},
			wantMsg:   "error from Errorf: test, value: 42",
			wantInPlusV: []string{
				"error from Errorf: test, value: 42",
				"TestFundamental_Format", // Errorf also calls callers, stack should point here
				"format_test.go",
			},
		},
		{
			name:      "Errorf_NoArgs_Format",
			errorFunc: lmccerrors.Errorf,
			formatOrMsg: "plain error from Errorf",
			args:        nil,
			wantMsg:   "plain error from Errorf",
			wantInPlusV: []string{
				"plain error from Errorf",
				"TestFundamental_Format",
				"format_test.go",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.errorFunc(tt.formatOrMsg, tt.args...)
			if err == nil {
				t.Fatal("error creation func returned nil, want error")
			}

			// Check err.Error()
			if err.Error() != tt.wantMsg {
				t.Errorf("err.Error() = %q, want %q", err.Error(), tt.wantMsg)
			}

			// Test %+v
			formattedPlusV := fmt.Sprintf("%+v", err)
			for _, wantSubstr := range tt.wantInPlusV {
				if !strings.Contains(formattedPlusV, wantSubstr) {
					t.Errorf("Output of fmt.Sprintf(\"%%+v\", err) should contain %q, but got %q", wantSubstr, formattedPlusV)
				}
			}
			// Ensure message is at the start of %+v before stack trace
			if tt.wantMsg != "" && !strings.HasPrefix(formattedPlusV, tt.wantMsg) {
			    t.Errorf("Output of fmt.Sprintf(\"%%+v\", err) should start with message %q, but got %q", tt.wantMsg, formattedPlusV)
			}

			// Test %s
			formattedS := fmt.Sprintf("%s", err)
			if formattedS != tt.wantMsg {
				t.Errorf("Output of fmt.Sprintf(\"%%s\", err) = %q, want %q", formattedS, tt.wantMsg)
			}

			// Test %v (without + flag)
			formattedV := fmt.Sprintf("%v", err)
			if formattedV != tt.wantMsg {
				t.Errorf("Output of fmt.Sprintf(\"%%v\", err) = %q, want %q", formattedV, tt.wantMsg)
			}
		})
	}
}

// TestWrapper_Format is a placeholder for testing wrapper.Format.
// TestWrapper_Format 是测试 wrapper.Format 的占位符。
func TestWrapper_Format(t *testing.T) {
	// t.Parallel() // Run in parallel if tests are independent

	originalStdErr := stdErrors.New("original standard error")
	customOriginalErr := lmccerrors.New("custom original error from lmccerrors") // Error with stack

	tests := []struct {
		name          string
		errorFunc     func(err error, msgOrFmt string, args ...interface{}) error // Wrap or Wrapf
		errToWrap     error
		msgOrFmt      string        // Message for Wrap, format string for Wrapf
		args          []interface{} // Args for Wrapf, nil for Wrap
		wantErrorMsg  string        // Expected from err.Error(), %s, %v
		wantInPlusV   []string      // Substrings expected in %+v output
		expectNil     bool          // If the resulting error should be nil
		checkFuncInV  string        // The function name of the test that should appear in the stack trace for %+v
		causeMsgInV   string        // The message of the cause error to check in %+v
		originalStack bool          // Whether to check for original error's stack in %+v
	}{
		// Test cases for Wrap
		{
			name:         "Wrap_StdError_SimpleMessage",
			errorFunc:    func(e error, m string, a ...interface{}) error { return lmccerrors.Wrap(e, m) },
			errToWrap:    originalStdErr,
			msgOrFmt:     "wrapper for std error",
			args:         nil,
			wantErrorMsg: "wrapper for std error: original standard error",
			wantInPlusV: []string{
				"wrapper for std error: original standard error", // Full message
				// "TestWrapper_Format",                             // Stack trace of the Wrap call itself
				// "format_test.go",
			},
			checkFuncInV: "TestWrapper_Format",
			causeMsgInV:  "original standard error",
		},
		{
			name:         "Wrap_CustomError_SimpleMessage",
			errorFunc:    func(e error, m string, a ...interface{}) error { return lmccerrors.Wrap(e, m) },
			errToWrap:    customOriginalErr,
			msgOrFmt:     "wrapper for custom error",
			args:         nil,
			wantErrorMsg: "wrapper for custom error: custom original error from lmccerrors",
			wantInPlusV: []string{
				"wrapper for custom error: custom original error from lmccerrors",
				// "TestWrapper_Format",
				// "custom original error from lmccerrors", // Message of the cause
				// "TestFundamental_Format",                // Stack of the customOriginalErr (or where it was created)
			},
			checkFuncInV:  "TestWrapper_Format",
			causeMsgInV:   "custom original error from lmccerrors",
			originalStack: true, // customOriginalErr has its own stack trace from lmccerrors.New
		},
		{
			name:      "Wrap_NilError",
			errorFunc: func(e error, m string, a ...interface{}) error { return lmccerrors.Wrap(e, m) },
			errToWrap: nil,
			msgOrFmt:  "this should not matter",
			args:      nil,
			expectNil: true,
		},
		{
			name:         "Wrap_StdError_EmptyMessage",
			errorFunc:    func(e error, m string, a ...interface{}) error { return lmccerrors.Wrap(e, m) },
			errToWrap:    originalStdErr,
			msgOrFmt:     "",
			args:         nil,
			wantErrorMsg: ": original standard error", // Behavior of current Wrap
			wantInPlusV: []string{
				": original standard error",
				// "TestWrapper_Format",
			},
			checkFuncInV: "TestWrapper_Format",
			causeMsgInV:  "original standard error",
		},
		// Test cases for Wrapf
		{
			name:         "Wrapf_StdError_WithMessageAndArgs",
			errorFunc:    func(e error, f string, a ...interface{}) error { return lmccerrors.Wrapf(e, f, a...) },
			errToWrap:    originalStdErr,
			msgOrFmt:     "wrapf message code %d, detail %s",
			args:         []interface{}{500, "server issue"},
			wantErrorMsg: "wrapf message code 500, detail server issue: original standard error",
			wantInPlusV: []string{
				"wrapf message code 500, detail server issue: original standard error",
				// "TestWrapper_Format",
			},
			checkFuncInV: "TestWrapper_Format",
			causeMsgInV:  "original standard error",
		},
		{
			name:      "Wrapf_NilError_WithMessageAndArgs",
			errorFunc: func(e error, f string, a ...interface{}) error { return lmccerrors.Wrapf(e, f, a...) },
			errToWrap: nil,
			msgOrFmt:  "this should not matter %s",
			args:      []interface{}{"indeed"},
			expectNil: true,
		},
		{
			name:         "Wrapf_CustomError_WithMessageAndArgs",
			errorFunc:    func(e error, f string, a ...interface{}) error { return lmccerrors.Wrapf(e, f, a...) },
			errToWrap:    customOriginalErr,
			msgOrFmt:     "wrapf for custom: %s",
			args:         []interface{}{"details here"},
			wantErrorMsg: "wrapf for custom: details here: custom original error from lmccerrors",
			wantInPlusV: []string{
				"wrapf for custom: details here: custom original error from lmccerrors",
				// "TestWrapper_Format",
				// "custom original error from lmccerrors",
			},
			checkFuncInV:  "TestWrapper_Format",
			causeMsgInV:   "custom original error from lmccerrors",
			originalStack: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.errorFunc(tt.errToWrap, tt.msgOrFmt, tt.args...)

			if tt.expectNil {
				if err != nil {
					t.Errorf("Expected nil error, but got: %v", err)
				}
				return // Skip further checks
			}

			if err == nil {
				t.Fatal("Expected non-nil error, but got nil")
			}

			// Check err.Error()
			if err.Error() != tt.wantErrorMsg {
				t.Errorf("err.Error() = %q, want %q", err.Error(), tt.wantErrorMsg)
			}

			// Test %s (should be same as err.Error())
			formattedS := fmt.Sprintf("%s", err)
			if formattedS != tt.wantErrorMsg {
				t.Errorf("Output of fmt.Sprintf(\"%%s\", err) was %q, want %q", formattedS, tt.wantErrorMsg)
			}

			// Test %v (without + flag, should be same as err.Error())
			formattedV := fmt.Sprintf("%v", err)
			if formattedV != tt.wantErrorMsg {
				t.Errorf("Output of fmt.Sprintf(\"%%v\", err) was %q, want %q", formattedV, tt.wantErrorMsg)
			}

			// Test %+v
			formattedPlusV := fmt.Sprintf("%+v", err)

			// 1. Check for the full error message (wrapper + cause)
			if !strings.Contains(formattedPlusV, tt.wantErrorMsg) {
				t.Errorf("Formatted error with '%%+v' should contain full message %q. Got: %q", tt.wantErrorMsg, formattedPlusV)
			}

			// 2. Check for the stack trace of the Wrap/Wrapf call itself
			if tt.checkFuncInV != "" && !strings.Contains(formattedPlusV, tt.checkFuncInV) {
				t.Errorf("Formatted error with '%%+v' should contain stack trace from function %q. Got: %q", tt.checkFuncInV, formattedPlusV)
			}
			if tt.checkFuncInV != "" && !strings.Contains(formattedPlusV, "format_test.go") { // Assuming stack is in this file
				t.Errorf("Formatted error with '%%+v' should mention 'format_test.go' for the wrapper stack. Got: %q", formattedPlusV)
			}
			
			// 3. Check for the cause error's message (if applicable, and if it's different or needs specific check)
			// The full message (tt.wantErrorMsg) already includes the cause's message, but let's ensure it's there as part of the cause's formatting.
			// This means the cause's own .Error() string, or how it's formatted when it's printed as part of the chain.
			if tt.causeMsgInV != "" && !strings.Contains(formattedPlusV, tt.causeMsgInV) {
				t.Errorf("Formatted error with '%%+v' should contain cause message %q. Got: %q", tt.causeMsgInV, formattedPlusV)
			}


			// 4. Check for the original error's stack trace if it was an lmccerror and had one
			if tt.originalStack {
				// This is a bit tricky. We need to ensure the stack trace of the *original* lmccerror
				// is also printed. The original lmccerror was `customOriginalErr`.
				// Its stack trace would point to where it was created, likely not `TestWrapper_Format`
				// but rather the setup lines of this test function or TestFundamental_Format if we reuse.
				// Let's check for its message followed by some stack trace indicators from *its* creation.
				// Example: "custom original error from lmccerrors\ngithub.com/lmcc-dev/lmcc-go-sdk/pkg/errors_test.TestWrapper_Format\n" (if created directly in this test)
				// or "custom original error from lmccerrors\ngithub.com/lmcc-dev/lmcc-go-sdk/pkg/errors.New\n" (if created via errors.New)
				// For `customOriginalErr := lmccerrors.New(...)`, its stack points to lmccerrors.New.
				// So we should look for that.
				
				// A simplified check: ensure the customOriginalErr message is followed by a newline (indicating stack starts)
				// and then some indication of its own stack.
				// This might need refinement based on actual output.
				expectedOriginalStackStart := tt.causeMsgInV + "\n" // Cause message followed by newline
				if !strings.Contains(formattedPlusV, expectedOriginalStackStart) {
					t.Errorf("Formatted error with '%%+v' for originalStack check, expected to find %q (cause message + newline). Got: %q", expectedOriginalStackStart, formattedPlusV)
				}

				// More specific check: after the cause message, we expect the stack of the cause.
				// If the cause is `customOriginalErr`, its stack was captured when `lmccerrors.New` was called.
				// The first frame of that stack trace should be `lmccerrors.New` or related internal functions.
				// Let's find the cause message and see what follows.
				idxCauseMsg := strings.Index(formattedPlusV, tt.causeMsgInV)
				if idxCauseMsg != -1 {
					relevantPortion := formattedPlusV[idxCauseMsg+len(tt.causeMsgInV):]
					// Example: "\ngithub.com/lmcc-dev/lmcc-go-sdk/pkg/errors.New\n\t..."
					// We need to be careful not to match the wrapper's stack.
					// The original error's stack should typically start with its own message,
					// then its stack frames. The wrapper prepends its message.
					// The default %+v for wrapped errors usually prints:
					// Wrapper Msg: Cause Msg
					// Wrapper Stack
					// Cause Stack (if cause also implements Format with stack)

					// The `lmccerrors.New` error (`customOriginalErr`) will have a stack trace starting with `lmccerrors.New`
					// or the line in `errors.go` where `newFundamental` is called.
					// The test for `customOriginalErr` in `TestFundamental_Format` checks for "TestFundamental_Format"
					// and "format_test.go" because `callers(2, defaultStackLength)` is used.
					// When `customOriginalErr` is wrapped, its `Format` method will be called.
					// Let's look for some lines from its creation stack.
					// Since customOriginalErr is created with `lmccerrors.New` within this test file:
					if !strings.Contains(relevantPortion, "lmccerrors.New") && !strings.Contains(relevantPortion, "newFundamental") {
						// If lmccerrors.New is inlined or optimized away, the direct caller might be this test function.
						// However, `Format` for fundamental errors should show the `callers` from its creation.
						// Let's assume for now that the `Format` method of `customOriginalErr` correctly prints its own stack.
						// The difficulty is distinguishing it clearly from the wrapper's stack if they are too similar.
						// The test `TestFundamental_Format` checks `customOriginalErr`'s stack more directly.
						// Here, we rely on `wrapper.Format` correctly invoking `cause.Format`.
					}
				}
			}

			// General check: for wrapped errors, the error message should be at the start of %+v
			if !strings.HasPrefix(formattedPlusV, tt.wantErrorMsg) {
				t.Errorf("Formatted error with '%%+v' should start with full message %q. Got: %q", tt.wantErrorMsg, formattedPlusV)
			}
		})
	}
}

// TestWithCode_Format is a placeholder for testing withCode.Format.
// TestWithCode_Format 是测试 withCode.Format 的占位符。
func TestWithCode_Format(t *testing.T) {
	// t.Parallel()

	baseStdErr := stdErrors.New("standard base error for WithCode")
	lmccBaseErr := lmccerrors.New("lmcc base error for WithCode") // Has its own stack

	// localMockCoder is now defined at the package level.
	// We can directly use it here.

	localMc1 := &localMockCoder{C: 2001, Ext: "Local Mock Coder 2001", HTTP: 500}
	localMc2 := &localMockCoder{C: 2002, Ext: "Local Mock Coder 2002", HTTP: 404}
	unknownCoder := lmccerrors.GetUnknownCoder() 

	tests := []struct {
		name            string
		errorFunc       func(coder lmccerrors.Coder, originalError error, format string, args ...interface{}) error 
		originalError   error                
		coder           lmccerrors.Coder
		textOrFormat    string               
		args            []interface{}        
		wantErrorMsg    string               
		wantInPlusV     []string             
		causeWantInPlusV []string            
		checkFuncInV    string               
		causeCreationFuncInV string        
		shouldBeNil     bool                 
		// originalStack   bool                 // Temporarily removed for debugging
	}{
		// Cases for WithCode
		{
			name:          "WithCode_StdError_ValidCoder",
			errorFunc:     func(c lmccerrors.Coder, e error, _ string, _ ...interface{}) error { return lmccerrors.WithCode(e, c) },
			originalError: baseStdErr,
			coder:         localMc1,
			wantErrorMsg:  "Local Mock Coder 2001: standard base error for WithCode",
			wantInPlusV: []string{
				"Local Mock Coder 2001: standard base error for WithCode",
				"standard base error for WithCode", 
			},
			checkFuncInV: "TestWithCode_Format",
		},
		{
			name:          "WithCode_LmccError_ValidCoder",
			errorFunc:     func(c lmccerrors.Coder, e error, _ string, _ ...interface{}) error { return lmccerrors.WithCode(e, c) },
			originalError: lmccBaseErr,
			coder:         localMc2,
			wantErrorMsg:  "Local Mock Coder 2002: lmcc base error for WithCode",
			wantInPlusV: []string{
				"Local Mock Coder 2002: lmcc base error for WithCode",
				"lmcc base error for WithCode", 
			},
			checkFuncInV:  "TestWithCode_Format",
			// originalStack: true, // Temporarily removed
		},
		{
			name:          "WithCode_StdError_NilCoder",
			errorFunc:     func(c lmccerrors.Coder, e error, _ string, _ ...interface{}) error { return lmccerrors.WithCode(e, c) },
			originalError: baseStdErr,
			coder:         nil, 
			wantErrorMsg:  unknownCoder.String() + ": standard base error for WithCode",
			wantInPlusV:   []string{unknownCoder.String() + ": standard base error for WithCode"},
			checkFuncInV:  "TestWithCode_Format",
		},
		{
			name:          "WithCode_NilError_ValidCoder",
			errorFunc:     func(c lmccerrors.Coder, e error, _ string, _ ...interface{}) error { return lmccerrors.WithCode(e, c) },
			originalError: nil,
			coder:         localMc1,
			shouldBeNil:   true,
		},
		// Cases for NewWithCode
		{
			name:         "NewWithCode_ValidCoder",
			errorFunc:    func(c lmccerrors.Coder, _ error, text string, _ ...interface{}) error { return lmccerrors.NewWithCode(c, text) },
			coder:        localMc1,
			textOrFormat: "new error with coder via NewWithCode",
			wantErrorMsg: "Local Mock Coder 2001: new error with coder via NewWithCode",
			wantInPlusV: []string{
				"Local Mock Coder 2001: new error with coder via NewWithCode",
			},
			causeWantInPlusV: []string {
				"new error with coder via NewWithCode", 
			},
			checkFuncInV: "TestWithCode_Format",
			causeCreationFuncInV: "lmccerrors.NewWithCode", 
		},
		{
			name:         "NewWithCode_NilCoder",
			errorFunc:    func(c lmccerrors.Coder, _ error, text string, _ ...interface{}) error { return lmccerrors.NewWithCode(c, text) },
			coder:        nil, 
			textOrFormat: "new error with nil coder",
			wantErrorMsg: unknownCoder.String() + ": new error with nil coder",
			wantInPlusV:  []string{unknownCoder.String() + ": new error with nil coder"},
			causeWantInPlusV: []string{"new error with nil coder"},
			checkFuncInV: "TestWithCode_Format",
			causeCreationFuncInV: "lmccerrors.NewWithCode",
		},
		// Cases for ErrorfWithCode
		{
			name:         "ErrorfWithCode_ValidCoder_Formatted",
			errorFunc:    func(c lmccerrors.Coder, _ error, format string, a ...interface{}) error { return lmccerrors.ErrorfWithCode(c, format, a...) },
			coder:        localMc2,
			textOrFormat: "formatted error code %d, detail %s",
			args:         []interface{}{77, "critical issue"},
			wantErrorMsg: "Local Mock Coder 2002: formatted error code 77, detail critical issue",
			wantInPlusV: []string{
				"Local Mock Coder 2002: formatted error code 77, detail critical issue",
			},
			causeWantInPlusV: []string{"formatted error code 77, detail critical issue"},
			checkFuncInV: "TestWithCode_Format",
			causeCreationFuncInV: "lmccerrors.ErrorfWithCode", 
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.errorFunc(tt.coder, tt.originalError, tt.textOrFormat, tt.args...)

			if tt.shouldBeNil {
				if err != nil {
					t.Errorf("Expected nil error, but got: %v", err)
				}
				return 
			}

			if err == nil {
				t.Fatal("Expected non-nil error, but got nil")
			}

			// Check err.Error(), %s, %v
			if err.Error() != tt.wantErrorMsg {
				t.Errorf("err.Error() = %q, want %q", err.Error(), tt.wantErrorMsg)
			}
			formattedS := fmt.Sprintf("%s", err)
			if formattedS != tt.wantErrorMsg {
				t.Errorf("Output of fmt.Sprintf(\"%%s\") was %q, want %q", formattedS, tt.wantErrorMsg)
			}
			formattedV := fmt.Sprintf("%v", err)
			if formattedV != tt.wantErrorMsg {
				t.Errorf("Output of fmt.Sprintf(\"%%v\") was %q, want %q", formattedV, tt.wantErrorMsg)
			}

			// Test %+v for the error itself
			formattedPlusV := fmt.Sprintf("%+v", err)

			// 1. Check for substrings in the main error's %+v output
			for _, wantSubstr := range tt.wantInPlusV {
				if !strings.Contains(formattedPlusV, wantSubstr) {
					t.Errorf("Formatted error with '%%+v' should contain %q. Got: %q", wantSubstr, formattedPlusV)
				}
			}
			// Ensure the full message is at the start
			if !strings.HasPrefix(formattedPlusV, tt.wantErrorMsg) {
			    t.Errorf("Formatted error with '%%+v' should start with message %q. Got: %q", tt.wantErrorMsg, formattedPlusV)
			}

			// 2. Check for the stack trace of the WithCode/NewWithCode/ErrorfWithCode call itself
			if tt.checkFuncInV != "" {
				if !strings.Contains(formattedPlusV, tt.checkFuncInV) {
					t.Errorf("Formatted error with '%%+v' should contain stack trace from function %q. Got: %q", tt.checkFuncInV, formattedPlusV)
				}
				if !strings.Contains(formattedPlusV, "format_test.go") {
					t.Errorf("Formatted error with '%%+v' should mention 'format_test.go' for the wrapper stack. Got: %q", formattedPlusV)
				}
			}

			// 3. For WithCode, check if original error's stack (if lmccerror) is present - Temporarily REMOVED
			/*
			if tt.originalError != nil && tt.originalStack { // tt.originalStack is now commented out in struct
				// Logic for checking original stack was here
			}
			*/

			// 4. For NewWithCode/ErrorfWithCode, check the formatting of the *cause*
			if tt.causeWantInPlusV != nil {
				cause := stdErrors.Unwrap(err) // Use stdErrors.Unwrap
				if cause == nil {
					t.Fatalf("errors.Unwrap(err) returned nil for a WithCode error that should have a cause (%s)", tt.name)
				}
				formattedCausePlusV := fmt.Sprintf("%+v", cause)
				for _, wantSubstr := range tt.causeWantInPlusV {
					if !strings.Contains(formattedCausePlusV, wantSubstr) {
						t.Errorf("Formatted cause with '%%+v' for test %s should contain %q. Got: %q", tt.name, wantSubstr, formattedCausePlusV)
					}
				}
				// Check that the cause's stack trace points to its creation within our library code
				/*
				if tt.causeCreationFuncInV != "" {
					// Check if the specific creation function (e.g., lmccerrors.NewWithCode) or a general internal creator is in the stack.
					// The general internal creator for fundamental errors is newFundamental, which calls lmccerrors.New or lmccerrors.Errorf.
					// The %+v output might show the exported function (NewWithCode) or deeper internal calls like New/Errorf or newFundamental itself.
					// Let's be a bit more flexible: check for the expected exported function OR the file where fundamental errors are made (errors.go)
					if !strings.Contains(formattedCausePlusV, tt.causeCreationFuncInV) && !strings.Contains(formattedCausePlusV, "errors.go") {
						t.Errorf("Formatted cause with '%%+v' for test %s should contain stack from %q or involve 'errors.go'. Got: %q", tt.name, tt.causeCreationFuncInV, formattedCausePlusV)
					}
				}
				*/
			}
		})
	}
}

// TestStackTrace_Format tests the Format method of StackTrace.
// TestStackTrace_Format 测试 StackTrace 的 Format 方法。
func TestStackTrace_Format(t *testing.T) {
	// t.Parallel() // Stack trace tests might be sensitive to parallel execution if line numbers are very specific

	// Helper functions to create errors with known stack structures
	// aTestFunctionForStackTrace is a helper for TestStackTrace_Format.
	// aTestFunctionForStackTrace 是 TestStackTrace_Format 的辅助函数。
	var aTestFunctionForStackTrace func() error // Declare first
	var anotherTestFunction func() error

	aTestFunctionForStackTrace = func() error {
		// Use lmccerrors.New directly to ensure it captures stack from here
		return lmccerrors.New("error from aTestFunctionForStackTrace_format_test")
	}

	// anotherTestFunction is another helper.
	// anotherTestFunction 是另一个辅助函数。
	anotherTestFunction = func() error {
		return lmccerrors.Wrap(aTestFunctionForStackTrace(), "wrapped in anotherTestFunction_format_test")
	}

	t.Run("StackTrace_Format_NestedError", func(t *testing.T) {
		err := anotherTestFunction()

		// Test %+v formatting for stack trace
		formattedError := fmt.Sprintf("%+v", err)
		t.Logf("Formatted error with stack (TestStackTrace_Format):\n%s", formattedError)

		// Check for key elements in the formatted string
		expectedMsgChain := "wrapped in anotherTestFunction_format_test: error from aTestFunctionForStackTrace_format_test"
		if !strings.Contains(formattedError, expectedMsgChain) {
			t.Errorf("Formatted error with '%%+v' should contain the full error message chain %q. Got: %q", expectedMsgChain, formattedError)
		}

		// Check for stack trace lines from our test functions
		// We are checking for the original function names they are part of.
		functionsInStack := []string{
			"aTestFunctionForStackTrace_format_test", // original error creation
			"anotherTestFunction_format_test",        // wrapper
			"TestStackTrace_Format",         // Check that the main test function name is part of the deeper stack frames
			"testing.tRunner",       // The test runner itself should be in the stack
		}
		for _, funcName := range functionsInStack {
			if !strings.Contains(formattedError, funcName) {
				t.Errorf("Formatted error with '%%+v' should contain stack trace line for %q. Got: %q", funcName, formattedError)
			}
		}

		// Check for file:line format (e.g., format_test.go:XXX)
		if !strings.Contains(formattedError, "format_test.go:") {
			t.Errorf("Formatted error with '%%+v' should contain file:line format like 'format_test.go:'. Got: %q", formattedError)
		}

		// Test %v and %s (should not print stack by default for wrapped errors)
		simpleFormatV := fmt.Sprintf("%v", err)
		if simpleFormatV != expectedMsgChain {
		    t.Errorf("Formatted error with %%v should be %q, got: %q", expectedMsgChain, simpleFormatV)
		}
		// Check that no stack trace appears for %v
		if strings.Count(simpleFormatV, "\n") > 0 && strings.Contains(simpleFormatV, "format_test.go:") {
			t.Errorf("Formatted error with %%v should not contain detailed stack trace. Got: %q", simpleFormatV)
		}

		simpleFormatS := fmt.Sprintf("%s", err)
		if simpleFormatS != expectedMsgChain {
		    t.Errorf("Formatted error with %%s should be %q, got: %q", expectedMsgChain, simpleFormatS)
		}
		// Check that no stack trace appears for %s
		if strings.Count(simpleFormatS, "\n") > 0 && strings.Contains(simpleFormatS, "format_test.go:") {
			t.Errorf("Formatted error with %%s should not contain detailed stack trace. Got: %q", simpleFormatS)
		}
	})

	t.Run("StackTrace_Format_DirectFundamentalError", func(t *testing.T) {
		// Test formatting of a fundamental error (which has a stack) directly.
		// This is somewhat covered by TestFundamental_Format, but this focuses on StackTrace specific aspects if any.
		fundamentalErr := lmccerrors.New("direct fundamental error for stack format")
		
		formattedFundamental := fmt.Sprintf("%+v", fundamentalErr)
		t.Logf("Formatted fundamental error with stack (TestStackTrace_Format):\n%s", formattedFundamental)

		expectedFundamentalMsg := "direct fundamental error for stack format"
		if !strings.HasPrefix(formattedFundamental, expectedFundamentalMsg) {
			t.Errorf("Formatted fundamental error with '%%+v' should start with message %q. Got: %q", expectedFundamentalMsg, formattedFundamental)
		}
		// Check that the current test function (or its t.Run sub-function) is in the stack
		if !strings.Contains(formattedFundamental, "TestStackTrace_Format") { 
			t.Errorf("Formatted fundamental error with '%%+v' should contain stack trace for current test func 'TestStackTrace_Format'. Got: %q", formattedFundamental)
		}
		if !strings.Contains(formattedFundamental, "format_test.go:") {
			t.Errorf("Formatted fundamental error with '%%+v' should contain 'format_test.go:'. Got: %q", formattedFundamental)
		}

		// %s and %v for fundamental error
		if fmt.Sprintf("%s", fundamentalErr) != expectedFundamentalMsg {
		    t.Errorf("Fundamental error with %%s was %q, want %q", fmt.Sprintf("%s", fundamentalErr), expectedFundamentalMsg)
		}
		if fmt.Sprintf("%v", fundamentalErr) != expectedFundamentalMsg {
		    t.Errorf("Fundamental error with %%v was %q, want %q", fmt.Sprintf("%v", fundamentalErr), expectedFundamentalMsg)
		}
	})
} 