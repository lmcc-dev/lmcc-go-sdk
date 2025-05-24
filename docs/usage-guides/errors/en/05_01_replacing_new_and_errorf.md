<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 1. Replacing `errors.New` and `fmt.Errorf`

- Replace `errors.New("message")` with `pkgErrors.New("message")`.
- Replace `fmt.Errorf("format %s", var)` with `pkgErrors.Errorf("format %s", var)`.

The `pkg/errors` versions will automatically capture stack traces.

```go
package main

import (
	"fmt"
	standardErrors "errors" // Go standard library errors
	pkgErrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors" // Our pkg/errors
)

// oldFunction uses standard library errors
func oldFunction(succeed bool) error {
	if !succeed {
		// Standard library errors.New
		return standardErrors.New("oldFunction failed due to a static reason")
	}
	return nil
}

// anotherOldFunction uses standard library fmt.Errorf
func anotherOldFunction(value int) error {
	if value < 0 {
		// Standard library fmt.Errorf
		return fmt.Errorf("anotherOldFunction failed: invalid value %d", value)
	}
	return nil
}

// newFunction uses pkg/errors.New
func newFunction(succeed bool) error {
	if !succeed {
		// pkg/errors.New - captures stack trace
		return pkgErrors.New("newFunction failed due to a static reason")
	}
	return nil
}

// anotherNewFunction uses pkg/errors.Errorf
func anotherNewFunction(value int) error {
	if value < 0 {
		// pkg/errors.Errorf - captures stack trace
		return pkgErrors.Errorf("anotherNewFunction failed: invalid value %d", value)
	}
	return nil
}

func main() {
	fmt.Println("--- Standard Library Errors (No Automatic Stack Trace with %v) ---")
	errOld1 := oldFunction(false)
	if errOld1 != nil {
		fmt.Printf("oldFunction error (%%v): %v\n", errOld1)
		// Standard errors don't automatically print stack trace with %+v 
		// unless they implement a specific Formatter interface, which errors.New doesn't.
		fmt.Printf("oldFunction error (%%+v): %+v\n", errOld1) 
	}

	errOld2 := anotherOldFunction(-10)
	if errOld2 != nil {
		fmt.Printf("anotherOldFunction error (%%v): %v\n", errOld2)
		fmt.Printf("anotherOldFunction error (%%+v): %+v\n", errOld2)
	}

	fmt.Println("\n--- pkg/errors (With Automatic Stack Trace with %+v) ---")
	errNew1 := newFunction(false)
	if errNew1 != nil {
		fmt.Printf("newFunction error (%%v): %v\n", errNew1)
		fmt.Printf("newFunction error (%%+v):\n%+v\n", errNew1) // Will show stack trace
	}

	errNew2 := anotherNewFunction(-20)
	if errNew2 != nil {
		fmt.Printf("anotherNewFunction error (%%v): %v\n", errNew2)
		fmt.Printf("anotherNewFunction error (%%+v):\n%+v\n", errNew2) // Will show stack trace
	}
}

/*
Example Output (Stack traces will vary based on your environment):

--- Standard Library Errors (No Automatic Stack Trace with %v) ---
oldFunction error (%v): oldFunction failed due to a static reason
oldFunction error (%+v): oldFunction failed due to a static reason
anotherOldFunction error (%v): anotherOldFunction failed: invalid value -10
anotherOldFunction error (%+v): anotherOldFunction failed: invalid value -10

--- pkg/errors (With Automatic Stack Trace with %+v) ---
newFunction error (%v): newFunction failed due to a static reason
newFunction error (%+v):
newFunction failed due to a static reason
main.newFunction
	/path/to/your/file.go:30
main.main
	/path/to/your/file.go:51
runtime.main
	...
runtime.goexit
	...
anotherNewFunction error (%v): anotherNewFunction failed: invalid value -20
anotherNewFunction error (%+v):
anotherNewFunction failed: invalid value -20
main.anotherNewFunction
	/path/to/your/file.go:39
main.main
	/path/to/your/file.go:57
runtime.main
	...
runtime.goexit
	...
*/
``` 