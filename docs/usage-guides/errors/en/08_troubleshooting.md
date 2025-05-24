<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

## Troubleshooting

Common issues and how to address them when working with `pkg/errors`.

**1. Stack Trace Not Appearing**
   - **Issue**: You print an error with `%+v` but don't see a stack trace.
   - **Possible Causes & Solutions**:
     - The error was not created by `pkg/errors` (e.g., it's a standard library error like `io.EOF` or from `errors.New()` from the standard library, and was not subsequently wrapped by `pkg/errors` functions like `pkgErrors.Wrap` or `pkgErrors.WithCode`).
     - You are printing the error with `%v` or `%s` instead of `%+v`.
     - The error was created with `pkgErrors.WithCode(nil, someCoder)`. `WithCode` on a nil error returns nil.
     - If an error from `pkg/errors` is wrapped by the standard library's `fmt.Errorf("... %w", pkgErr)`, then `fmt.Printf("%+v", wrappedErr)` might not display the stack trace from `pkgErr` unless `pkgErr` itself implements a `Format` method that `fmt.Errorf` knows how to call for the `%+v` verb through the chain. Our `pkg/errors` `fundamental` type does implement `Format`, which should make its stack trace visible when it's the cause.

   ```go
   package main

   import (
   	"fmt"
   	standardErrors "errors"
   	pkgErrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
   	"io"
   )

   func main() {
   	// Cause 1: Standard library error, not wrapped by pkg/errors
   	stdErr := io.EOF
   	fmt.Println("--- Standard Error (io.EOF) ---")
   	fmt.Printf("%%v: %v\n", stdErr)
   	fmt.Printf("%%+v: %+v\n\n", stdErr) // No stack trace from io.EOF itself

   	// Wrapped by pkgErrors - NOW it gets a stack (from the wrapping point)
   	wrappedStdErr := pkgErrors.Wrap(stdErr, "failed to read file")
   	fmt.Println("--- Standard Error Wrapped by pkgErrors.Wrap ---")
   	fmt.Printf("%%v: %v\n", wrappedStdErr)
   	fmt.Printf("%%+v:\n%+v\n\n", wrappedStdErr) 

   	// Cause 2: Using %v instead of %+v for a pkg/errors error
   	pkgErr := pkgErrors.New("a pkg/errors error")
   	fmt.Println("--- pkg/errors Error ---")
   	fmt.Printf("%%v: %v\n", pkgErr)       // No stack trace
   	fmt.Printf("%%+v:\n%+v\n\n", pkgErr) // Stack trace appears
   	
   	// Cause 3: fmt.Errorf wrapping a pkgError
   	fmtWrappedPkgErr := fmt.Errorf("wrapped by fmt.Errorf: %w", pkgErr)
   	fmt.Println("--- pkg/errors Error Wrapped by fmt.Errorf --- ")
   	fmt.Printf("%%v: %v\n", fmtWrappedPkgErr)
   	// The fundamental.Format method in pkg/errors should allow stack trace to be shown even here.
   	fmt.Printf("%%+v:\n%+v\n\n", fmtWrappedPkgErr) 
   }
   /* Output:
   --- Standard Error (io.EOF) ---
   %v: EOF
   %+v: EOF

   --- Standard Error Wrapped by pkgErrors.Wrap ---
   %v: failed to read file: EOF
   %+v:
   failed to read file: EOF
   main.main
       /path/to/file.go:XX
   ...

   --- pkg/errors Error ---
   %v: a pkg/errors error
   %+v:
   a pkg/errors error
   main.main
       /path/to/file.go:XX
   ...

   --- pkg/errors Error Wrapped by fmt.Errorf --- 
   %v: wrapped by fmt.Errorf: a pkg/errors error
   %+v:
   wrapped by fmt.Errorf: a pkg/errors error
   main.main
       /path/to/file.go:XX
   ...
   */
   ```

**2. `errors.Is` or `errors.As` Not Working as Expected**
   - **Issue**: `standardErrors.Is(err, target)` returns `false`, or `standardErrors.As(err, &targetType)` doesn't find the type, even though you think it should.
   - **Possible Causes & Solutions**:
     - **`Is` with different instances**: If `target` is a non-nil error created at a different place (e.g., `err = pkgErrors.New("msg"); target = pkgErrors.New("msg")`), `errors.Is` will be `false` because they are different instances. `errors.Is` checks for reference equality or if an error in the chain reports itself as equivalent via an `Is(error) bool` method.
       - For `Coder` types from `pkg/errors`, checking `errors.Is(err, pkgErrors.ErrNotFound)` works because `pkgErrors.ErrNotFound` is a specific global variable (sentinel).
       - For checking general categories by code, use `pkgErrors.IsCode(err, CoderWithTheCodeYouWant)`.
     - **`As` with wrong type**: Ensure the `targetType` variable you pass to `errors.As` is a pointer to an interface type or a pointer to a concrete type that implements `error`.
     - **Error not in chain**: The specific error instance or type you're checking for might not actually be in `err`'s chain of wrapped errors.

   ```go
   package main

   import (
   	"fmt"
   	standardErrors "errors"
   	pkgErrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
   	"io"
   )

   type MyTroubleError struct{ Msg string }
   func (e *MyTroubleError) Error() string { return e.Msg }

   var SentinelError = pkgErrors.New("Specific sentinel error instance")

   func main() {
   	// --- errors.Is issues ---
   	err1 := pkgErrors.Wrap(SentinelError, "wrapped sentinel")
   	fmt.Printf("Is err1 SentinelError? %t\n", standardErrors.Is(err1, SentinelError)) // true

   	nonSentinelPkgErr := pkgErrors.New("some message")
   	anotherNonSentinel := pkgErrors.New("some message")
   	fmt.Printf("Is nonSentinelPkgErr anotherNonSentinel? %t\n", standardErrors.Is(nonSentinelPkgErr, anotherNonSentinel)) // false, different instances

   	wrappedIoEOF := pkgErrors.Wrap(io.EOF, "wrapped io.EOF")
   	fmt.Printf("Is wrappedIoEOF io.EOF? %t\n", standardErrors.Is(wrappedIoEOF, io.EOF)) // true

   	// --- errors.As issues ---
   	customErrInstance := &MyTroubleError{"custom trouble"}
   	errWithCustom := pkgErrors.Wrap(customErrInstance, "context")

   	var target *MyTroubleError
   	if standardErrors.As(errWithCustom, &target) {
   		fmt.Printf("As found MyTroubleError: %s\n", target.Msg)
   	} else {
   		fmt.Println("As did not find MyTroubleError")
   	}

   	var targetIOErr *os.PathError // Example of a type not in the chain
   	if standardErrors.As(errWithCustom, &targetIOErr) {
   		fmt.Printf("As found os.PathError: %s\n", targetIOErr.Path)
   	} else {
   		fmt.Println("As did not find os.PathError in errWithCustom")
   	}
   }
   /* Output:
   Is err1 SentinelError? true
   Is nonSentinelPkgErr anotherNonSentinel? false
   Is wrappedIoEOF io.EOF? true
   As found MyTroubleError: custom trouble
   As did not find os.PathError in errWithCustom
   */
   ```

**3. `Coder` Information Not Being Retrieved**
   - **Issue**: `errors.GetCoder(err)` returns `nil`.
   - **Possible Causes & Solutions**:
     - No error in the chain was created with `NewWithCode`, `ErrorfWithCode`, or had a `Coder` attached via `WithCode` or `WrapWithCode`.
     - The error chain might only contain standard library errors or errors from other packages that don't use the `pkg/errors` `Coder` mechanism.

   ```go
   package main

   import (
   	"fmt"
   	standardErrors "errors"
   	pkgErrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
   )

   func main() {
   	// Error with a Coder
   	errWithCode := pkgErrors.NewWithCode(pkgErrors.ErrBadRequest, "bad request indeed")
   	coder1 := pkgErrors.GetCoder(errWithCode)
   	if coder1 != nil {
   		fmt.Printf("Coder found: Code=%d, Message='%s'\n", coder1.Code(), coder1.String())
   	} else {
   		fmt.Println("No coder found in errWithCode")
   	}

   	// Standard error, no Coder
   	stdErr := standardErrors.New("a plain standard error")
   	coder2 := pkgErrors.GetCoder(stdErr)
   	if coder2 != nil {
   		fmt.Printf("Coder found in stdErr: Code=%d, Message='%s'\n", coder2.Code(), coder2.String())
   	} else {
   		fmt.Println("No coder found in stdErr")
   	}

   	// Standard error wrapped by pkgErrors.Wrap (still no Coder from original error)
   	wrappedStdErr := pkgErrors.Wrap(stdErr, "context added")
   	coder3 := pkgErrors.GetCoder(wrappedStdErr)
   	if coder3 != nil {
   		fmt.Printf("Coder found in wrappedStdErr: Code=%d, Message='%s'\n", coder3.Code(), coder3.String())
   	} else {
   		fmt.Println("No coder found in wrappedStdErr (original error had no Coder)")
   	}
   	
   	// Standard error wrapped by pkgErrors.WrapWithCode
   	wrappedStdErrWithCode := pkgErrors.WrapWithCode(stdErr, pkgErrors.ErrInternalServer, "wrapped with code")
   	coder4 := pkgErrors.GetCoder(wrappedStdErrWithCode)
   	if coder4 != nil {
   		fmt.Printf("Coder found in wrappedStdErrWithCode: Code=%d, Message='%s'\n", coder4.Code(), coder4.String())
   	} else {
   		fmt.Println("No coder found in wrappedStdErrWithCode")
   	}
   }
   /* Output:
   Coder found: Code=100003, Message='Bad request'
   No coder found in stdErr
   No coder found in wrappedStdErr (original error had no Coder)
   Coder found in wrappedStdErrWithCode: Code=100001, Message='Internal server error'
   */
   ```

If you encounter other issues, ensure you are using the functions from the `pkg/errors` module as intended and check the specific function documentation for behavior details (e.g., how `nil` errors are handled). 