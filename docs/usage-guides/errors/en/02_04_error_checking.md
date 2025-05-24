<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 4. Error Checking

Once you have an error, you'll want to inspect it.

**`standardErrors.Is(err, target error) bool`**: (From the standard library) Reports whether any error in `err`'s chain matches `target`. This is useful for checking against sentinel errors like `errors.ErrNotFound`.

**`errors.GetCoder(err error) Coder`**: Traverses the error chain and returns the first `Coder` found. Returns `nil` if no `Coder` is found.

**`errors.IsCode(err error, coder Coder) bool`**: Reports whether any error in `err`'s chain has a `Coder` that matches the `Code()` of the provided `coder`.


```go
package main

import (
	"fmt"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	standardErrors "errors" // Standard library errors
)

// A custom sentinel error
var ErrCustomOperationFailed = errors.New("custom operation failed")

// A custom Coder
var ErrPaymentDeclined = errors.NewCoder(60001, 402, "Payment declined by provider", "")

// functionThatReturnsVariousErrors
func functionThatReturnsVariousErrors(condition string) error {
	switch condition {
	case "not_found":
		// Returns an error that wraps a predefined Coder (errors.ErrNotFound)
		originalErr := errors.New("database query returned no rows")
		return errors.WithCode(originalErr, errors.ErrNotFound)
	case "validation_error":
		// Returns an error with a different predefined Coder (errors.ErrValidation)
		return errors.NewWithCode(errors.ErrValidation, "input field 'email' is invalid")
	case "custom_sentinel":
		// Returns a wrapped custom sentinel error
		return errors.Wrap(ErrCustomOperationFailed, "attempting to process critical data")
	case "payment_declined":
		// Returns an error with our custom Coder ErrPaymentDeclined
		return errors.ErrorfWithCode(ErrPaymentDeclined, "transaction ID %s was declined", "txn_123abc")
	case "unclassified_error":
		// Returns an error without a Coder or known sentinel
		return errors.New("an unexpected issue occurred")
	default:
		return nil
	}
}

func main() {
	testConditions := []string{"not_found", "validation_error", "custom_sentinel", "payment_declined", "unclassified_error", "success"}

	for _, cond := range testConditions {
		fmt.Printf("--- Testing condition: %s ---\n", cond)
		err := functionThatReturnsVariousErrors(cond)

		if err == nil {
			fmt.Println("Operation successful, no error.")
			continue
		}

		fmt.Printf("Received error: %v\n", err)

		// 1. Check using standardErrors.Is for sentinel errors
		if standardErrors.Is(err, errors.ErrNotFound) { // Checking against a predefined Coder, which also acts like a sentinel
			fmt.Println("  [Is] This is a 'Not Found' error (checked via errors.ErrNotFound).")
		}
		if standardErrors.Is(err, ErrCustomOperationFailed) { // Checking against our custom sentinel
			fmt.Println("  [Is] This is our 'Custom Operation Failed' sentinel error.")
		}

		// 2. Extract Coder using errors.GetCoder
		if coder := errors.GetCoder(err); coder != nil {
			fmt.Printf("  [GetCoder] Extracted Coder: Code=%d, HTTPStatus=%d, Message='%s', Ref='%s'\n",
				coder.Code(), coder.HTTPStatus(), coder.String(), coder.Reference())

			// 3. Check by specific Coder instance using errors.IsCode (most precise for Coders)
			if errors.IsCode(err, errors.ErrNotFound) {
				fmt.Println("  [IsCode] This error has the 'ErrNotFound' Coder.")
			}
			if errors.IsCode(err, errors.ErrValidation) {
				fmt.Println("  [IsCode] This error has the 'ErrValidation' Coder.")
			}
			if errors.IsCode(err, ErrPaymentDeclined) {
				fmt.Println("  [IsCode] This error has our custom 'ErrPaymentDeclined' Coder.")
			}
		} else {
			fmt.Println("  [GetCoder] No Coder found in this error chain.")
		}
		fmt.Println("-----------------------------------\n")
	}
}
/*
Example Output:

--- Testing condition: not_found ---
Received error: Not found: database query returned no rows
  [Is] This is a 'Not Found' error (checked via errors.ErrNotFound).
  [GetCoder] Extracted Coder: Code=100002, HTTPStatus=404, Message='Not found', Ref='errors-spec.md#100002'
  [IsCode] This error has the 'ErrNotFound' Coder.
-----------------------------------

--- Testing condition: validation_error ---
Received error: Validation failed: input field 'email' is invalid
  [GetCoder] Extracted Coder: Code=100006, HTTPStatus=400, Message='Validation failed', Ref='errors-spec.md#100006'
  [IsCode] This error has the 'ErrValidation' Coder.
-----------------------------------

--- Testing condition: custom_sentinel ---
Received error: attempting to process critical data: custom operation failed
  [Is] This is our 'Custom Operation Failed' sentinel error.
  [GetCoder] No Coder found in this error chain.
-----------------------------------

--- Testing condition: payment_declined ---
Received error: Payment declined by provider: transaction ID txn_123abc was declined
  [GetCoder] Extracted Coder: Code=60001, HTTPStatus=402, Message='Payment declined by provider', Ref=''
  [IsCode] This error has our custom 'ErrPaymentDeclined' Coder.
-----------------------------------

--- Testing condition: unclassified_error ---
Received error: an unexpected issue occurred
  [GetCoder] No Coder found in this error chain.
-----------------------------------

--- Testing condition: success ---
Operation successful, no error.
*/
```

**Key points for error checking:**
- Use `standardErrors.Is` for checking against sentinel error values (like `io.EOF`, or custom ones like `ErrCustomOperationFailed`). Predefined `Coder` instances from `pkg/errors` (e.g., `errors.ErrNotFound`) can also be used with `standardErrors.Is` because they are, in effect, sentinel values.
- Use `errors.GetCoder` to retrieve a `Coder` if one exists in the error chain. You can then inspect its properties (`Code()`, `HTTPStatus()`, etc.).
- Use `errors.IsCode` to specifically check if an error in the chain matches a particular `Coder`'s `Code()`. This is the most direct way to check for an error category defined by a `Coder`.

</rewritten_file> 