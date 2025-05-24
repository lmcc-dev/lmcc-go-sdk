<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 1. Creating Errors

Creating basic errors is straightforward. Both `errors.New()` and `errors.Errorf()` automatically capture the call stack.

**`errors.New(message string) error`**: Creates a simple error with the given message.

**`errors.Errorf(format string, args ...interface{}) error`**: Creates an error with a formatted message, similar to `fmt.Errorf`.

```go
package main

import (
	"fmt"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
)

// simulateOperationA
// This function demonstrates creating a simple error using errors.New.
func simulateOperationA() error {
	// Some condition leads to an error
	return errors.New("connection timed out during operation A")
}

// simulateOperationB
// This function demonstrates creating a formatted error using errors.Errorf.
func simulateOperationB(userID string, itemID int) error {
	// Some condition leads to an error
	if userID == "" {
		return errors.Errorf("failed to process item %d: userID cannot be empty", itemID)
	}
	return errors.Errorf("failed to retrieve details for item %d for user '%s': resource not found", itemID, userID)
}

func main() {
	// --- Example for errors.New ---
	errA := simulateOperationA()
	if errA != nil {
		fmt.Println("Error from Operation A:")
		fmt.Printf("  Simple view (%%v):  %v\n", errA)
		fmt.Printf("  Detailed view (%%+v):\n%+v\n", errA) 
		// The %+v format verb will print the error message along with the stack trace.
	}

	fmt.Println("\n--------------------------\n")

	// --- Example for errors.Errorf ---
	userID := "user123"
	itemID := 789
	errB := simulateOperationB(userID, itemID)
	if errB != nil {
		fmt.Println("Error from Operation B:")
		fmt.Printf("  Simple view (%%v):  %v\n", errB)
		fmt.Printf("  Detailed view (%%+v):\n%+v\n", errB)
	}
	
	fmt.Println("\n--------------------------\n")

	// Example with an empty userID to trigger a different path in simulateOperationB
	errC := simulateOperationB("", 456)
		if errC != nil {
		fmt.Println("Error from Operation B (empty userID):")
		fmt.Printf("  Simple view (%%v):  %v\n", errC)
		fmt.Printf("  Detailed view (%%+v):\n%+v\n", errC)
	}
}
/*
Example Output (Stack traces will vary based on file paths and line numbers):

Error from Operation A:
  Simple view (%v):  connection timed out during operation A
  Detailed view (%+v):
connection timed out during operation A
main.simulateOperationA
	/path/to/your/file.go:13
main.main
	/path/to/your/file.go:34
runtime.main
	/usr/local/go/src/runtime/proc.go:250
runtime.goexit
	/usr/local/go/src/runtime/asm_amd64.s:1197

--------------------------

Error from Operation B:
  Simple view (%v):  failed to retrieve details for item 789 for user 'user123': resource not found
  Detailed view (%+v):
failed to retrieve details for item 789 for user 'user123': resource not found
main.simulateOperationB
	/path/to/your/file.go:26
main.main
	/path/to/your/file.go:46
runtime.main
	/usr/local/go/src/runtime/proc.go:250
runtime.goexit
	/usr/local/go/src/runtime/asm_amd64.s:1197

--------------------------

Error from Operation B (empty userID):
  Simple view (%v):  failed to process item 456: userID cannot be empty
  Detailed view (%+v):
failed to process item 456: userID cannot be empty
main.simulateOperationB
	/path/to/your/file.go:23
main.main
	/path/to/your/file.go:58
runtime.main
	/usr/local/go/src/runtime/proc.go:250
runtime.goexit
	/usr/local/go/src/runtime/asm_amd64.s:1197
*/
```
