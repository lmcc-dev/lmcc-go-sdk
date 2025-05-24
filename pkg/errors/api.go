/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package errors

import (
	"errors"
)

// Coder defines the interface for an error code.
// It also embeds the standard error interface, so Coder instances can be used as errors directly.
// Coder 定义了错误码的接口。
// 它也嵌入了标准错误接口，因此 Coder 实例可以直接用作错误。
type Coder interface {
	error // Embed the standard error interface

	// Code returns the integer code of this error.
	// Code 返回错误的整数码。
	Code() int

	// String returns the string representation of this error code.
	// String 返回错误码的字符串表示 (通常是错误描述或类型)。
	// (Usually the error description or type).
	// Note: This might be shadowed by the embedded error.Error() if not careful.
	// For Coder's specific string representation, ensure usage is clear.
	String() string

	// HTTPStatus returns the associated HTTP status code, if applicable.
	// HTTPStatus 返回关联的 HTTP 状态码 (如果适用)。
	// May return 0 or a sentinel value if not applicable.
	// 如果不适用，可能返回 0 或哨兵值。
	HTTPStatus() int

	// Reference returns a URL or document reference for this error code.
	// Reference 返回此错误码的 URL 或文档参考。
	// May return an empty string if not applicable.
	// 如果不适用，可能返回空字符串。
	Reference() string
}

// IsCode checks if the error (or any error in its chain) has a Coder
// that matches the code of the provided Coder `c`.
// (IsCode 检查错误（或其链中的任何错误）是否拥有一个 Coder，
// 该 Coder 的代码与提供的 Coder `c` 的代码匹配。)
func IsCode(err error, c Coder) bool {
	if err == nil || c == nil {
		return false
	}

	targetCode := c.Code()

	for {
		if coderHolder, ok := err.(interface{ Coder() Coder }); ok {
			if currentCoder := coderHolder.Coder(); currentCoder != nil {
				if currentCoder.Code() == targetCode {
					return true
				}
			}
		}
		// Also check if the error itself is a Coder that matches
		if currentAsCoder, ok := err.(Coder); ok {
			if currentAsCoder.Code() == targetCode {
				return true
			}
		}

		// Check for multi-error unwrapping (Go 1.20+ style, like ErrorGroup)
		if multiUnwrapper, ok := err.(interface{ Unwrap() []error }); ok {
			for _, subErr := range multiUnwrapper.Unwrap() {
				if IsCode(subErr, c) {
					return true
				}
			}
		}

		unwrappedError := errors.Unwrap(err) // Use standard library errors.Unwrap
		if unwrappedError == nil {
			return false
		}
		err = unwrappedError
	}
}
