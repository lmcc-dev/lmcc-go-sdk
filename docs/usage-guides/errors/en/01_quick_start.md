<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

## Quick Start

This section shows a very basic comparison of using the standard `errors` package versus `lmcc-go-sdk/pkg/errors`.

**Scenario:** Create a simple error message.

```go
package main

import (
	"fmt"
	standardErrors "errors" // Standard library errors
	customErrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors" // Our custom errors package
)

func main() {
	// Example using the standard library's errors.New
	stdErr := standardErrors.New("standard library: operation failed")
	fmt.Printf("Standard error: %v\n", stdErr)
	// Output: Standard error: standard library: operation failed
	// Note: No automatic stack trace here.

	fmt.Println("---")

	// Example using github.com/lmcc-dev/lmcc-go-sdk/pkg/errors
	lmccErr := customErrors.New("lmcc-go-sdk: operation failed")
	fmt.Printf("lmcc-go-sdk error (simple view): %v\n", lmccErr)
	// Output: lmcc-go-sdk error (simple view): lmcc-go-sdk: operation failed

	// To see the stack trace, you would use %+v with fmt.Printf
	fmt.Printf("lmcc-go-sdk error (detailed view with stack trace):\n%+v\n", lmccErr)
	// Example Output (will vary based on where you run it):
	// lmcc-go-sdk error (detailed view with stack trace):
	// lmcc-go-sdk: operation failed
	// main.main
	//	/path/to/your/main.go:XX
	// ... other stack frames ...
}
```
The key takeaway is that `lmcc-go-sdk/pkg/errors` automatically captures a stack trace when an error is created with `New` or `Errorf`, which is invaluable for debugging.
