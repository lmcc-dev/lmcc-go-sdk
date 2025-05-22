/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package errors_test // Use a different package name for black box testing

import (
	// Standard library errors
	"errors"
	"fmt"
	"testing"

	lmccerrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
)

// mockCoder is a simple Coder implementation for testing.
// mockCoder 是用于测试的简单 Coder 实现。
type mockCoder struct {
	C    int
	Ext  string
	HTTP int
	Ref  string
}

func (m *mockCoder) Code() int            { return m.C }
func (m *mockCoder) String() string       { return m.Ext }
func (m *mockCoder) HTTPStatus() int      { return m.HTTP }
func (m *mockCoder) Reference() string    { return m.Ref }

// Error implements the error interface for mockCoder.
// Error 为 mockCoder 实现 error 接口。
func (m *mockCoder) Error() string        { return m.Ext }

var (
	mc1 = &mockCoder{C: 1001, Ext: "Mock Coder 1001", HTTP: 500}
	mc2 = &mockCoder{C: 1002, Ext: "Mock Coder 1002", HTTP: 404}
)

// coderError wraps a Coder to make it an error, for testing with errors.Is.
// coderError 包装一个 Coder 使其成为一个错误，用于通过 errors.Is 进行测试。
type coderError struct {
	lmccerrors.Coder
}

// Error implements the error interface for coderError.
// Error 为 coderError 实现 error 接口。
func (ce coderError) Error() string {
	if ce.Coder == nil {
		return "<nil coder error>"
	}
	return fmt.Sprintf("CoderError: %s", ce.Coder.String())
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		msg     string
		wantMsg string
	}{
		{
			name:    "Simple error",
			msg:     "a new error occurred",
			wantMsg: "a new error occurred",
		},
		{
			name:    "Empty message",
			msg:     "",
			wantMsg: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := lmccerrors.New(tt.msg)
			if err == nil {
				t.Fatal("lmccerrors.New() returned nil, want error")
			}
			if err.Error() != tt.wantMsg {
				t.Errorf("err.Error() = %q, want %q", err.Error(), tt.wantMsg)
			}
			// Basic Sprintf checks (%s, %v) remain as they test fundamental error string representation.
			// Detailed %+v formatting is in format_test.go.
			formattedS := fmt.Sprintf("%s", err)
			if formattedS != tt.wantMsg {
			    t.Errorf("Output of fmt.Sprintf(\"%%s\", err) was %q, want %q", formattedS, tt.wantMsg)
			}
			formattedV := fmt.Sprintf("%v", err)
			if formattedV != tt.wantMsg {
			    t.Errorf("Output of fmt.Sprintf(\"%%v\", err) was %q, want %q", formattedV, tt.wantMsg)
			}
		})
	}
}

func TestErrorf(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []interface{}
		wantMsg  string
		wantErr  bool 
	}{
		{
			name:     "Simple formatted error",
			format:   "error: %s with value %d",
			args:     []interface{}{"something", 123},
			wantMsg:  "error: something with value 123",
			wantErr:  true,
		},
		{
			name:     "No arguments",
			format:   "a plain error",
			args:     []interface{}{},
			wantMsg:  "a plain error",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := lmccerrors.Errorf(tt.format, tt.args...)
			if tt.wantErr {
				if err == nil {
					t.Fatal("lmccerrors.Errorf() returned nil, want error")
				}
				if err.Error() != tt.wantMsg {
					t.Errorf("err.Error() = %q, want %q", err.Error(), tt.wantMsg)
				}
				// Basic Sprintf checks (%s, %v) remain.
				// Detailed %+v formatting is in format_test.go.
				formattedS := fmt.Sprintf("%s", err)
				if formattedS != tt.wantMsg {
				    t.Errorf("Output of fmt.Sprintf(\"%%s\", err) was %q, want %q", formattedS, tt.wantMsg)
				}
				formattedV := fmt.Sprintf("%v", err)
				if formattedV != tt.wantMsg {
				    t.Errorf("Output of fmt.Sprintf(\"%%v\", err) was %q, want %q", formattedV, tt.wantMsg)
				}

			} else {
				if err != nil {
					t.Errorf("lmccerrors.Errorf() returned error %v, want nil", err)
				}
			}
		})
	}
}

func TestWrap(t *testing.T) {
	originalErr := errors.New("original error")
	customOriginalErr := lmccerrors.New("custom original error") // This New already captures stack

	tests := []struct {
		name          string
		errToWrap     error
		wrapMsg       string
		wantOutputMsg string 
		wantCauseIs   error  
	}{
		{
			name:          "Wrap standard error",
			errToWrap:     originalErr,
			wrapMsg:       "failed to process",
			wantOutputMsg: "failed to process: original error",
			wantCauseIs:   originalErr,
		},
		{
			name:          "Wrap custom error",
			errToWrap:     customOriginalErr,
			wrapMsg:       "layer 2 wrapper",
			wantOutputMsg: "layer 2 wrapper: custom original error",
			wantCauseIs:   customOriginalErr,
		},
		{
			name:          "Wrap nil error",
			errToWrap:     nil,
			wrapMsg:       "this should not appear",
			wantOutputMsg: "", 
			wantCauseIs:   nil,
		},
		{
			name:          "Wrap with empty message",
			errToWrap:     originalErr,
			wrapMsg:       "",
			wantOutputMsg: ": original error", 
			wantCauseIs:   originalErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrappedErr := lmccerrors.Wrap(tt.errToWrap, tt.wrapMsg)

			if tt.errToWrap == nil {
				if wrappedErr != nil {
					t.Errorf("Wrap(nil, msg) = %v, want nil", wrappedErr)
				}
				return 
			}

			if wrappedErr == nil {
				t.Fatal("Wrap(err, msg) returned nil, want error")
			}

			if wrappedErr.Error() != tt.wantOutputMsg {
				t.Errorf("wrappedErr.Error() = %q, want %q", wrappedErr.Error(), tt.wantOutputMsg)
			}

			// Test Unwrap using standard library errors.Unwrap
			unwrapped := errors.Unwrap(wrappedErr)
			if unwrapped != tt.wantCauseIs {
				t.Errorf("errors.Unwrap(wrappedErr) = %v, want %v", unwrapped, tt.wantCauseIs)
			}

			// Test errors.Is
			if !errors.Is(wrappedErr, tt.wantCauseIs) {
				t.Errorf("errors.Is(wrappedErr, originalErr) = false, want true")
			}

			// %+v formatting tests are now in format_test.go
		})
	}
}

func TestWrapf(t *testing.T) {
	originalErr := errors.New("base error for wrapf")
	tests := []struct {
		name          string
		errToWrap     error
		format        string
		args          []interface{}
		wantOutputMsg string
		wantCauseIs   error
	}{
		{
			name:          "Wrapf standard error with formatting",
			errToWrap:     originalErr,
			format:        "failed with %s: %d",
			args:          []interface{}{"code red", 500},
			wantOutputMsg: "failed with code red: 500: base error for wrapf",
			wantCauseIs:   originalErr,
		},
		{
			name:          "Wrapf nil error",
			errToWrap:     nil,
			format:        "this should not appear %s",
			args:          []interface{}{"indeed"},
			wantOutputMsg: "",
			wantCauseIs:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrappedErr := lmccerrors.Wrapf(tt.errToWrap, tt.format, tt.args...)

			if tt.errToWrap == nil {
				if wrappedErr != nil {
					t.Errorf("Wrapf(nil, format, args...) = %v, want nil", wrappedErr)
				}
				return
			}

			if wrappedErr == nil {
				t.Fatal("Wrapf(err, format, args...) returned nil, want error")
			}

			if wrappedErr.Error() != tt.wantOutputMsg {
				t.Errorf("wrappedErr.Error() = %q, want %q", wrappedErr.Error(), tt.wantOutputMsg)
			}

			if !errors.Is(wrappedErr, tt.wantCauseIs) {
				t.Errorf("errors.Is(wrappedErr, originalErr) = false, want true")
			}

			// %+v formatting tests are now in format_test.go
		})
	}
}

func TestWithCode(t *testing.T) {
	baseErr := errors.New("database error")
	lmccBaseErr := lmccerrors.New("lmcc base error")

	tests := []struct {
		name            string
		errToWrap       error
		coder           lmccerrors.Coder
		wantErrorMsg    string
		wantCoderCode   int
		wantCoderString string
		shouldBeNil     bool
	}{
		{
			name:            "Wrap standard error with Coder",
			errToWrap:       baseErr,
			coder:           mc1,
			wantErrorMsg:    "Mock Coder 1001: database error",
			wantCoderCode:   1001,
			wantCoderString: "Mock Coder 1001",
		},
		{
			name:            "Wrap lmcc error with Coder",
			errToWrap:       lmccBaseErr,
			coder:           mc2,
			wantErrorMsg:    "Mock Coder 1002: lmcc base error",
			wantCoderCode:   1002,
			wantCoderString: "Mock Coder 1002",
		},
		{
			name:        "Wrap nil error with Coder",
			errToWrap:   nil,
			coder:       mc1,
			shouldBeNil: true,
		},
		{
			name:            "Wrap with nil Coder",
			errToWrap:       baseErr,
			coder:           nil, // Should use unknownCoder
			wantErrorMsg:    "An internal server error occurred: database error",
			wantCoderCode:   -1, // unknownCoder.Code()
			wantCoderString: "An internal server error occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errWithCode := lmccerrors.WithCode(tt.errToWrap, tt.coder)

			if tt.shouldBeNil {
				if errWithCode != nil {
					t.Errorf("WithCode(nil, coder) = %v, want nil", errWithCode)
				}
				return
			}

			if errWithCode == nil {
				t.Fatal("WithCode(err, coder) returned nil, want error")
			}

			if errWithCode.Error() != tt.wantErrorMsg {
				t.Errorf("err.Error() = %q, want %q", errWithCode.Error(), tt.wantErrorMsg)
			}

			// Test GetCoder
			retrievedCoder := lmccerrors.GetCoder(errWithCode)
			if retrievedCoder == nil {
				t.Fatal("GetCoder(errWithCode) returned nil")
			}
			if retrievedCoder.Code() != tt.wantCoderCode {
				t.Errorf("GetCoder().Code() = %d, want %d", retrievedCoder.Code(), tt.wantCoderCode)
			}
			if retrievedCoder.String() != tt.wantCoderString {
				t.Errorf("GetCoder().String() = %q, want %q", retrievedCoder.String(), tt.wantCoderString)
			}

			// Test errors.Is with the Coder itself (wrapped in coderError)
			expectedCoderForIs := tt.coder
			if tt.coder == nil {
				expectedCoderForIs = lmccerrors.GetUnknownCoder()
			}
			if !errors.Is(errWithCode, coderError{expectedCoderForIs}) {
				t.Errorf("errors.Is(errWithCode, CoderError(%s)) failed; Coder: %v", expectedCoderForIs.String(), expectedCoderForIs)
			}

			// Test errors.As for Coder
			var extractedCoder lmccerrors.Coder
			if !errors.As(errWithCode, &extractedCoder) {
				t.Error("errors.As(errWithCode, &Coder) failed")
			} else {
				if extractedCoder.Code() != tt.wantCoderCode {
					t.Errorf("extracted Coder.Code() = %d, want %d", extractedCoder.Code(), tt.wantCoderCode)
				}
			}

			// Test Unwrap
			if errors.Unwrap(errWithCode) != tt.errToWrap {
				t.Errorf("errors.Unwrap(errWithCode) did not return original wrapped error")
			}

			// %+v formatting tests are now in format_test.go
		})
	}
}

func TestNewWithCode(t *testing.T) {
	tests := []struct {
		name            string
		text            string
		coder           lmccerrors.Coder
		wantErrorMsg    string
		wantCoderCode   int
	// wantCoderString string // Redundant, checked by wantCoderCode and coder equality effectively
		wantCauseMsg    string 
	}{
		{
			name:            "New error with Coder",
			text:            "resource creation failed",
			coder:           mc1,
			wantErrorMsg:    "Mock Coder 1001: resource creation failed",
			wantCoderCode:   1001,
			wantCauseMsg:    "resource creation failed",
		},
		{
			name:            "New error with nil Coder",
			text:            "another issue",
			coder:           nil, // Should use unknownCoder
			wantErrorMsg:    "An internal server error occurred: another issue",
			wantCoderCode:   -1,
			wantCauseMsg:    "another issue",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := lmccerrors.NewWithCode(tt.coder, tt.text)
			if err == nil {
				t.Fatal("NewWithCode() returned nil, want error")
			}

			if err.Error() != tt.wantErrorMsg {
				t.Errorf("err.Error() = %q, want %q", err.Error(), tt.wantErrorMsg)
			}

			retrievedCoder := lmccerrors.GetCoder(err)
			if retrievedCoder == nil {
			    t.Fatalf("GetCoder(err) returned nil for NewWithCode error")
			}
			if retrievedCoder.Code() != tt.wantCoderCode {
				t.Errorf("GetCoder().Code() = %d, want %d", retrievedCoder.Code(), tt.wantCoderCode)
			}

			cause := errors.Unwrap(err)
			if cause == nil {
				t.Fatal("errors.Unwrap(err) returned nil for NewWithCode error")
			}
			if cause.Error() != tt.wantCauseMsg {
				t.Errorf("cause.Error() = %q, want %q", cause.Error(), tt.wantCauseMsg)
			}

			// Cause formatting and main error %+v formatting are in format_test.go
		})
	}
}

func TestErrorfWithCode(t *testing.T) {
	tests := []struct {
		name            string
		format          string
		args            []interface{}
		coder           lmccerrors.Coder
		wantErrorMsg    string
		wantCoderCode   int
		wantCauseMsg    string
	}{
		{
			name:            "Errorf with Coder and formatting",
			format:          "ID %d not found for user %s",
			args:            []interface{}{123, "john.doe"},
			coder:           mc2,
			wantErrorMsg:    "Mock Coder 1002: ID 123 not found for user john.doe",
			wantCoderCode:   1002,
			wantCauseMsg:    "ID 123 not found for user john.doe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := lmccerrors.ErrorfWithCode(tt.coder, tt.format, tt.args...)
			if err == nil {
				t.Fatal("ErrorfWithCode() returned nil")
			}
			if err.Error() != tt.wantErrorMsg {
				t.Errorf("err.Error() = %q, want %q", err.Error(), tt.wantErrorMsg)
			}
			retrievedCoder := lmccerrors.GetCoder(err)
			if retrievedCoder == nil {
			    t.Fatalf("GetCoder(err) returned nil for ErrorfWithCode error")
			}
			if retrievedCoder.Code() != tt.wantCoderCode {
				t.Errorf("GetCoder().Code() from ErrorfWithCode = %d, want %d", retrievedCoder.Code(), tt.wantCoderCode)
			}
			cause := errors.Unwrap(err)
			if cause == nil {
			    t.Fatalf("errors.Unwrap(err) returned nil for ErrorfWithCode error")
			}
			if cause.Error() != tt.wantCauseMsg {
				t.Errorf("Unwrapped cause error = %q, want %q", cause.Error(), tt.wantCauseMsg)
			}
			// Cause formatting and main error %+v formatting are in format_test.go
		})
	}
}

// Further tests for WithMessage, WithMessagef,
// Is, As (more complex cases), Cause will be added progressively.
// 후속 테스트는 WithMessage, WithMessagef,
// Is, As (더 복잡한 경우), Cause에 대해 점진적으로 추가됩니다. 