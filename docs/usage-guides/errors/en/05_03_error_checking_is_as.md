<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 3. Error Checking with `errors.Is` and `errors.As`

Code using `standardErrors.Is(err, target)` and `standardErrors.As(err, &targetType)` will continue to work as expected with errors created or wrapped by `pkg/errors`.

- **`standardErrors.Is`**: Checks if any error in the chain matches a target sentinel error. Works with `pkg/errors` sentinel-like `Coder` variables (e.g., `pkgErrors.ErrNotFound`).
- **`standardErrors.As`**: Checks if any error in the chain can be cast to a specific type and sets the target to that error. Useful for custom error types that might be part of the chain.
- **`pkgErrors.IsCode`**: Use this for checking against the numeric code of a `Coder` if you don't need to check for a specific `Coder` instance (see Best Practices).

```go
package main

import (
	"fmt"
	standardErrors "errors" // Go standard library errors
	pkgErrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors" // Our pkg/errors
	"os"
)

// A custom error type that might be used
type MyCustomErrorType struct {
	Msg     string
	Details string
}

func (e *MyCustomErrorType) Error() string {
	return fmt.Sprintf("%s (Details: %s)", e.Msg, e.Details)
}

// Function that can return various types of errors
func generateError(errorType string) error {
	switch errorType {
	case "std_is_target":
		// A standard library error that we might check with errors.Is
		return os.ErrPermission // Predefined standard error
	case "pkg_is_target":
		// A pkg/errors Coder that acts as a sentinel for errors.Is
		return pkgErrors.Wrap(pkgErrors.ErrUnauthorized, "user lacks permissions for this pkg operation")
	case "custom_as_target":
		// An instance of our custom error type, possibly wrapped
		customErr := &MyCustomErrorType{Msg: "Custom operation failed", Details: "Database constraint violation"}
		return pkgErrors.Wrap(customErr, "wrapping custom error type")
	case "pkg_with_code":
		// A pkg/errors error with a specific Coder, good for IsCode
		return pkgErrors.NewWithCode(pkgErrors.ErrConflict, "resource version mismatch")
	default:
		return pkgErrors.New("an unknown error occurred")
	}
}

func main() {
	testCases := []string{"std_is_target", "pkg_is_target", "custom_as_target", "pkg_with_code", "unknown"}

	for _, tc := range testCases {
		fmt.Printf("--- Testing error type: %s ---\n", tc)
		err := generateError(tc)

		if err == nil {
			fmt.Println("No error generated.")
			continue
		}
		fmt.Printf("Generated error: %v\n", err)

		// 1. Using standardErrors.Is
		if standardErrors.Is(err, os.ErrPermission) {
			fmt.Println("  [errors.Is] Matched os.ErrPermission.")
		}
		if standardErrors.Is(err, pkgErrors.ErrUnauthorized) { // pkgErrors.ErrUnauthorized is a Coder, also acts as a sentinel
			fmt.Println("  [errors.Is] Matched pkgErrors.ErrUnauthorized.")
		}

		// 2. Using standardErrors.As
		var customErrTarget *MyCustomErrorType
		if standardErrors.As(err, &customErrTarget) {
			fmt.Printf("  [errors.As] Matched MyCustomErrorType. Message: '%s', Details: '%s'\n", customErrTarget.Msg, customErrTarget.Details)
		} else {
			// fmt.Println("  [errors.As] Did not match MyCustomErrorType.")
		}

		// 3. Using pkgErrors.IsCode (for pkg/errors Coders)
		if pkgErrors.IsCode(err, pkgErrors.ErrUnauthorized) {
			fmt.Println("  [pkgErrors.IsCode] Matched code for pkgErrors.ErrUnauthorized.")
		}
		if pkgErrors.IsCode(err, pkgErrors.ErrConflict) {
			fmt.Println("  [pkgErrors.IsCode] Matched code for pkgErrors.ErrConflict.")
			coder := pkgErrors.GetCoder(err)
			if coder != nil {
				fmt.Printf("    Coder details: Code=%d, HTTPStatus=%d, Message='%s'\n", coder.Code(), coder.HTTPStatus(), coder.String())
			}
		}
		fmt.Println("") // Newline for clarity
	}
}

/*
Example Output:

--- Testing error type: std_is_target ---
Generated error: permission denied
  [errors.Is] Matched os.ErrPermission.

--- Testing error type: pkg_is_target ---
Generated error: user lacks permissions for this pkg operation: Unauthorized
  [errors.Is] Matched pkgErrors.ErrUnauthorized.
  [pkgErrors.IsCode] Matched code for pkgErrors.ErrUnauthorized.

--- Testing error type: custom_as_target ---
Generated error: wrapping custom error type: Custom operation failed (Details: Database constraint violation)
  [errors.As] Matched MyCustomErrorType. Message: 'Custom operation failed', Details: 'Database constraint violation'

--- Testing error type: pkg_with_code ---
Generated error: Conflict: resource version mismatch
  [pkgErrors.IsCode] Matched code for pkgErrors.ErrConflict.
    Coder details: Code=100005, HTTPStatus=409, Message='Conflict'

--- Testing error type: unknown ---
Generated error: an unknown error occurred

*/
``` 