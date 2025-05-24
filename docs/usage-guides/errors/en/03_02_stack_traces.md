<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### Stack Traces

All errors created by `errors.New`, `errors.Errorf`, `errors.NewWithCode`, `errors.ErrorfWithCode`, `errors.Wrap`, and `errors.Wrapf` automatically capture a stack trace at the point of their creation. Errors wrapped by `errors.WithCode` retain the stack trace of the original error.

The stack trace is printed when the error is formatted with `"%+v"` using `fmt.Printf` or similar functions.

```go
package main

import (
	"fmt"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
)

// functionA calls functionB
func functionA(triggerError bool) error {
	fmt.Println("Executing functionA...")
	err := functionB(triggerError)
	if err != nil {
		// Wrap the error from functionB, adding context from functionA.
		// The original stack trace from functionB's error creation point is preserved,
		// and this wrapping point is also part of the chain.
		return errors.Wrap(err, "functionA encountered an issue")
	}
	return nil
}

// functionB calls functionC
func functionB(triggerError bool) error {
	fmt.Println("Executing functionB...")
	err := functionC(triggerError)
	if err != nil {
		// Here, we are just passing the error up. No new stack trace is added by functionB itself
		// because we are not creating a new error or wrapping it with pkg/errors functions here.
		// If we were to use errors.Wrap() here, then this point would also be in the trace.
		return err
	}
	return nil
}

// functionC creates an error if triggerError is true
func functionC(triggerError bool) error {
	fmt.Println("Executing functionC...")
	if triggerError {
		// Create a new error with a Coder. This is where the primary stack trace will originate.
		return errors.NewWithCode(errors.ErrInternalServer, "simulated critical failure in functionC")
	}
	fmt.Println("functionC completed successfully.")
	return nil
}

func main() {
	fmt.Println("--- Scenario 1: No error triggered ---")
	err1 := functionA(false)
	if err1 == nil {
		fmt.Println("Scenario 1 completed without errors.")
	}

	fmt.Println("\n--- Scenario 2: Error triggered and propagated ---")
	err2 := functionA(true)
	if err2 != nil {
		fmt.Println("\nAn error occurred:")
		fmt.Println("\n--- Printing error with %v ---")
		fmt.Printf("%v\n", err2)

		fmt.Println("\n--- Printing error with %+v (includes stack trace) ---")
		fmt.Printf("%+v\n", err2) // This will print the error message and the stack trace.
		
		// Demonstrate Cause()
		originalCause := errors.Cause(err2)
		fmt.Println("\n--- Original cause (using errors.Cause) with %+v --- ")
		fmt.Printf("%+v\n", originalCause)
	}
}

/*
Example Output (file paths and line numbers will vary based on your environment):

--- Scenario 1: No error triggered ---
Executing functionA...
Executing functionB...
Executing functionC...
functionC completed successfully.
Scenario 1 completed without errors.

--- Scenario 2: Error triggered and propagated ---
Executing functionA...
Executing functionB...
Executing functionC...

An error occurred:

--- Printing error with %v ---
functionA encountered an issue: Internal server error: simulated critical failure in functionC

--- Printing error with %+v (includes stack trace) ---
functionA encountered an issue: simulated critical failure in functionC
main.functionC
	/path/to/your/file.go:32
main.functionB
	/path/to/your/file.go:20
main.functionA
	/path/to/your/file.go:10
main.main
	/path/to/your/file.go:48
runtime.main
	/usr/local/go/src/runtime/proc.go:267
runtime.goexit
	/usr/local/go/src/runtime/asm_arm64.s:1197
Internal server error

--- Original cause (using errors.Cause) with %+v --- 
simulated critical failure in functionC
main.functionC
	/path/to/your/file.go:32
main.functionB
	/path/to/your/file.go:20
main.functionA
	/path/to/your/file.go:10
main.main
	/path/to/your/file.go:48
runtime.main
	/usr/local/go/src/runtime/proc.go:267
runtime.goexit
	/usr/local/go/src/runtime/asm_arm64.s:1197
Internal server error
*/
```

**Note on Stack Traces and Wrapping:**
- When you wrap an error using `errors.Wrap` or `errors.Wrapf`, the original error's stack trace (if it was created by `pkg/errors`) is preserved. The `Wrap` call itself adds a new frame to the conceptual stack of messages but doesn't generate a *new* full stack trace; it chains the errors.
- `errors.WithCode` also preserves the original error's stack trace.
- Formatting with `%+v` will typically show the stack trace from the point where the *innermost* error (the cause) was created by a `pkg/errors` function, followed by the messages from wrapping errors. 