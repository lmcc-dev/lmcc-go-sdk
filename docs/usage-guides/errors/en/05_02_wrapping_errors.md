<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 2. Wrapping Errors

If you were using `fmt.Errorf` with the `%w` verb to wrap errors, you can replace it with `pkgErrors.Wrap` or `pkgErrors.Wrapf`.

- `fmt.Errorf("context: %w", err)` becomes `pkgErrors.Wrap(err, "context")` or `pkgErrors.Wrapf(err, "context with %s", var)`.

`pkg/errors` wrapping functions also preserve the original error for `errors.Is` and `errors.As` and ensure the stack trace from the original `pkg/errors` error is maintained.

```go
package main

import (
	"fmt"
	standardErrors "errors" // Go standard library errors
	pkgErrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors" // Our pkg/errors
	"os"
)

var ErrOriginalStd = standardErrors.New("original standard error")
var ErrOriginalPkg = pkgErrors.New("original pkg/errors error")

// wrapWithStdFmt uses fmt.Errorf with %w
func wrapWithStdFmt(originalError error, context string) error {
	return fmt.Errorf("%s: %w", context, originalError)
}

// wrapWithPkgErrors uses pkgErrors.Wrap
func wrapWithPkgErrors(originalError error, context string) error {
	return pkgErrors.Wrap(originalError, context)
}

func main() {
	fmt.Println("--- Wrapping a standard library error ---")
	// Standard library error wrapped by fmt.Errorf
	wrappedStdByStd := wrapWithStdFmt(ErrOriginalStd, "std lib wrapped by fmt.Errorf")
	fmt.Printf("Std wrapped by fmt.Errorf (%%v): %v\n", wrappedStdByStd)
	fmt.Printf("Std wrapped by fmt.Errorf (%%+v): %+v\n", wrappedStdByStd) // No auto stack from original
	fmt.Printf("  Is ErrOriginalStd? %t\n", standardErrors.Is(wrappedStdByStd, ErrOriginalStd))

	fmt.Println("\n--- Wrapping a pkg/errors error ---")
	// pkg/errors error wrapped by pkgErrors.Wrap
	wrappedPkgByPkg := wrapWithPkgErrors(ErrOriginalPkg, "pkg/errors wrapped by pkgErrors.Wrap")
	fmt.Printf("pkg wrapped by pkgErrors.Wrap (%%v): %v\n", wrappedPkgByPkg)
	fmt.Printf("pkg wrapped by pkgErrors.Wrap (%%+v):\n%+v\n", wrappedPkgByPkg) // Stack trace from ErrOriginalPkg is shown
	fmt.Printf("  Is ErrOriginalPkg? %t\n", standardErrors.Is(wrappedPkgByPkg, ErrOriginalPkg))

	fmt.Println("\n--- Interoperability: Wrapping a pkg/errors error with fmt.Errorf %w ---")
	// pkg/errors error wrapped by fmt.Errorf %w
	wrappedPkgByStd := wrapWithStdFmt(ErrOriginalPkg, "pkg/errors wrapped by fmt.Errorf")
	fmt.Printf("pkg wrapped by fmt.Errorf (%%v): %v\n", wrappedPkgByStd)
	// fmt.Errorf does not know how to format pkgErrors for %+v to show stack trace from cause,
	// unless pkgErrors.fundamental itself implemented a specific Formatter logic for this.
	// Our pkgErrors.fundamental.Format handles its own stack trace, not necessarily when wrapped by external fmt.Errorf.
	fmt.Printf("pkg wrapped by fmt.Errorf (%%+v): %+v\n", wrappedPkgByStd) 
	fmt.Printf("  Is ErrOriginalPkg? %t\n", standardErrors.Is(wrappedPkgByStd, ErrOriginalPkg))


	fmt.Println("\n--- Interoperability: Wrapping a standard library error with pkgErrors.Wrap ---")
	// standard library error (os.ErrNotExist) wrapped by pkgErrors.Wrap
	originalStdLibError := os.ErrNotExist
	wrappedStdByPkg := wrapWithPkgErrors(originalStdLibError, "os.ErrNotExist wrapped by pkgErrors.Wrap")
	fmt.Printf("std (os.ErrNotExist) wrapped by pkg (%%v): %v\n", wrappedStdByPkg)
	// pkgErrors.Wrap adds a stack trace at the point of wrapping if the cause doesn't have one.
	// Since os.ErrNotExist doesn't have a stack trace that pkgErrors recognizes, pkgErrors.Wrap creates one.
	fmt.Printf("std (os.ErrNotExist) wrapped by pkg (%%+v):\n%+v\n", wrappedStdByPkg) 
	fmt.Printf("  Is os.ErrNotExist? %t\n", standardErrors.Is(wrappedStdByPkg, os.ErrNotExist))

	fmt.Println("\n--- Accessing underlying error (Cause) ---")
	fmt.Printf("Cause of wrappedPkgByPkg: %v\n", pkgErrors.Cause(wrappedPkgByPkg))
	// standardErrors.Unwrap can also be used
	fmt.Printf("Unwrap of wrappedPkgByPkg: %v\n", standardErrors.Unwrap(wrappedPkgByPkg))
}

/*
Example Output (Stack traces will vary):

--- Wrapping a standard library error ---
Std wrapped by fmt.Errorf (%v): std lib wrapped by fmt.Errorf: original standard error
Std wrapped by fmt.Errorf (%+v): std lib wrapped by fmt.Errorf: original standard error
  Is ErrOriginalStd? true

--- Wrapping a pkg/errors error ---
pkg wrapped by pkgErrors.Wrap (%v): pkg/errors wrapped by pkgErrors.Wrap: original pkg/errors error
pkg wrapped by pkgErrors.Wrap (%+v):
pkg/errors wrapped by pkgErrors.Wrap: original pkg/errors error
main.main
	/path/to/your/file.go:20
runtime.main
	...
runtime.goexit
	...
  Is ErrOriginalPkg? true

--- Interoperability: Wrapping a pkg/errors error with fmt.Errorf %w ---
pkg wrapped by fmt.Errorf (%v): pkg/errors wrapped by fmt.Errorf: original pkg/errors error
pkg wrapped by fmt.Errorf (%+v): pkg/errors wrapped by fmt.Errorf: original pkg/errors error
  Is ErrOriginalPkg? true

--- Interoperability: Wrapping a standard library error with pkgErrors.Wrap ---
std (os.ErrNotExist) wrapped by pkg (%v): os.ErrNotExist wrapped by pkgErrors.Wrap: file does not exist
std (os.ErrNotExist) wrapped by pkg (%+v):
os.ErrNotExist wrapped by pkgErrors.Wrap: file does not exist
main.main
	/path/to/your/file.go:40
runtime.main
	...
runtime.goexit
	...
  Is os.ErrNotExist? true

--- Accessing underlying error (Cause) ---
Cause of wrappedPkgByPkg: original pkg/errors error
Unwrap of wrappedPkgByPkg: original pkg/errors error
*/
```

**Key takeaway for wrapping:** `pkgErrors.Wrap` and `pkgErrors.Wrapf` are the preferred way to wrap errors when using this library, as they ensure proper stack trace handling and context addition while maintaining compatibility with `standardErrors.Is` and `standardErrors.As`. 