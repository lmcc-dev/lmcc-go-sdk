/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package errors

import (
	"net/http"
	"testing"
)

// TestNewCoder tests the NewCoder function and the methods of basicCoder.
// TestNewCoder 测试 NewCoder 函数和 basicCoder 的方法。
func TestNewCoder(t *testing.T) {
	tests := []struct {
		name        string
		code        int
		httpStatus  int
		description string
		reference   []string
		wantCode    int
		wantHTTP    int
		wantString  string
		wantRef     string
	}{
		{
			name:        "Full arguments",
			code:        101,
			httpStatus:  http.StatusInternalServerError,
			description: "Internal Error",
			reference:   []string{"http://example.com/docs/101"},
			wantCode:    101,
			wantHTTP:    http.StatusInternalServerError,
			wantString:  "Internal Error",
			wantRef:     "http://example.com/docs/101",
		},
		{
			name:        "No reference",
			code:        102,
			httpStatus:  http.StatusNotFound,
			description: "Not Found",
			reference:   []string{},
			wantCode:    102,
			wantHTTP:    http.StatusNotFound,
			wantString:  "Not Found",
			wantRef:     "",
		},
		{
			name:        "Empty description",
			code:        103,
			httpStatus:  http.StatusBadRequest,
			description: "",
			wantCode:    103,
			wantHTTP:    http.StatusBadRequest,
			wantString:  "",
			wantRef:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			refArg := ""
			if len(tt.reference) > 0 {
				refArg = tt.reference[0]
			}
			coder := NewCoder(tt.code, tt.httpStatus, tt.description, refArg)

			if got := coder.Code(); got != tt.wantCode {
				t.Errorf("Coder.Code() = %v, want %v", got, tt.wantCode)
			}
			if got := coder.HTTPStatus(); got != tt.wantHTTP {
				t.Errorf("Coder.HTTPStatus() = %v, want %v", got, tt.wantHTTP)
			}
			if got := coder.String(); got != tt.wantString {
				t.Errorf("Coder.String() = %v, want %v", got, tt.wantString)
			}
			if got := coder.Reference(); got != tt.wantRef {
				t.Errorf("Coder.Reference() = %v, want %v", got, tt.wantRef)
			}
		})
	}
}

// TestPredefinedCoders tests the properties of predefined Coder instances.
// TestPredefinedCoders 测试预定义 Coder 实例的属性。
func TestPredefinedCoders(t *testing.T) {
	tests := []struct {
		name       string
		coder      Coder
		wantCode   int
		wantHTTP   int
		wantString string
	}{
		{"unknownCoder", unknownCoder, -1, 500, "An internal server error occurred"},
		{"ErrInternalServer", ErrInternalServer, 100001, 500, "Internal server error"},
		{"ErrNotFound", ErrNotFound, 100002, 404, "Resource not found"},
		{"ErrBadRequest", ErrBadRequest, 100003, 400, "Bad request"},
		{"ErrUnauthorized", ErrUnauthorized, 100004, 401, "Unauthorized"},
		{"ErrForbidden", ErrForbidden, 100005, 403, "Forbidden"},
		{"ErrValidation", ErrValidation, 100006, 400, "Validation error"},
		{"ErrTimeout", ErrTimeout, 100007, 504, "Request timeout"},
		{"ErrTooManyRequests", ErrTooManyRequests, 100008, 429, "Too many requests"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.coder.Code(); got != tt.wantCode {
				t.Errorf("%s.Code() = %v, want %v", tt.name, got, tt.wantCode)
			}
			if got := tt.coder.HTTPStatus(); got != tt.wantHTTP {
				t.Errorf("%s.HTTPStatus() = %v, want %v", tt.name, got, tt.wantHTTP)
			}
			if got := tt.coder.String(); got != tt.wantString {
				t.Errorf("%s.String() = %v, want %v", tt.name, got, tt.wantString)
			}
			// Reference is expected to be empty for these predefined coders
			if ref := tt.coder.Reference(); ref != "" {
				t.Errorf("%s.Reference() = %v, want \"\"", tt.name, ref)
			}
		})
	}
}

// TestIsUnknownCoder tests the IsUnknownCoder helper function.
// TestIsUnknownCoder 测试 IsUnknownCoder 辅助函数。
func TestIsUnknownCoder(t *testing.T) {
	if !IsUnknownCoder(unknownCoder) {
		t.Errorf("IsUnknownCoder(unknownCoder) = false, want true")
	}
	if IsUnknownCoder(ErrInternalServer) {
		t.Errorf("IsUnknownCoder(ErrInternalServer) = true, want false")
	}
	coder := NewCoder(123, 200, "custom", "")
	if IsUnknownCoder(coder) {
		t.Errorf("IsUnknownCoder(customCoder) = true, want false")
	}
	if IsUnknownCoder(nil) {
		t.Errorf("IsUnknownCoder(nil) = true, want false")
	}
}

// TestGetUnknownCoder tests the GetUnknownCoder helper function.
// TestGetUnknownCoder 测试 GetUnknownCoder 辅助函数。
func TestGetUnknownCoder(t *testing.T) {
	got := GetUnknownCoder()
	if got == nil {
		t.Fatalf("GetUnknownCoder() returned nil")
	}
	if got.Code() != -1 || got.String() != "An internal server error occurred" || got.HTTPStatus() != 500 {
		t.Errorf("GetUnknownCoder() returned unexpected Coder: Code=%d, String=\"%s\", HTTP=%d",
			got.Code(), got.String(), got.HTTPStatus())
	}
	// Check if it's indeed the same instance (or an equivalent one)
	if !IsUnknownCoder(got) {
		t.Errorf("GetUnknownCoder() did not return the unknownCoder instance")
	}
}
