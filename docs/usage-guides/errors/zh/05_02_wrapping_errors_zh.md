<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 2. 包装错误 (Wrapping Errors)

如果您之前使用带有 `%w` 动词的 `fmt.Errorf` 来包装错误，您可以将其替换为 `pkgErrors.Wrap` 或 `pkgErrors.Wrapf`。

(If you were using `fmt.Errorf` with the `%w` verb to wrap errors, you can replace it with `pkgErrors.Wrap` or `pkgErrors.Wrapf`.)

- `fmt.Errorf("context: %w", err)` 变为 `pkgErrors.Wrap(err, "context")` 或 `pkgErrors.Wrapf(err, "context with %s", var)`。
  (`fmt.Errorf("context: %w", err)` becomes `pkgErrors.Wrap(err, "context")` or `pkgErrors.Wrapf(err, "context with %s", var)`.)

`pkg/errors` 包装函数还会保留原始错误以供 `errors.Is` 和 `errors.As` 使用，并确保维护原始 `pkg/errors` 错误的堆栈跟踪。

(`pkg/errors` wrapping functions also preserve the original error for `errors.Is` and `errors.As` and ensure the stack trace from the original `pkg/errors` error is maintained.)

```go
package main

import (
	"fmt"
	standardErrors "errors" // Go 标准库 errors (Go standard library errors)
	pkgErrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors" // 我们的 pkg/errors (Our pkg/errors)
	"os"
)

var ErrOriginalStd = standardErrors.New("原始标准错误 (original standard error)")
var ErrOriginalPkg = pkgErrors.New("原始 pkg/errors 错误 (original pkg/errors error)")

// wrapWithStdFmt 使用 fmt.Errorf 和 %w
// (wrapWithStdFmt uses fmt.Errorf with %w)
func wrapWithStdFmt(originalError error, context string) error {
	return fmt.Errorf("%s: %w", context, originalError)
}

// wrapWithPkgErrors 使用 pkgErrors.Wrap
// (wrapWithPkgErrors uses pkgErrors.Wrap)
func wrapWithPkgErrors(originalError error, context string) error {
	return pkgErrors.Wrap(originalError, context)
}

func main() {
	fmt.Println("--- 包装标准库错误 --- (--- Wrapping a standard library error ---)")
	// 由 fmt.Errorf 包装的标准库错误
	// (Standard library error wrapped by fmt.Errorf)
	wrappedStdByStd := wrapWithStdFmt(ErrOriginalStd, "由 fmt.Errorf 包装的标准库 (std lib wrapped by fmt.Errorf)")
	fmt.Printf("由 fmt.Errorf 包装的标准库 (%%v): %v\n", wrappedStdByStd)
	// (Std wrapped by fmt.Errorf (%%v): %v\n)
	fmt.Printf("由 fmt.Errorf 包装的标准库 (%%+v): %+v\n", wrappedStdByStd) // 原始错误没有自动堆栈跟踪 (No auto stack from original)
	// (Std wrapped by fmt.Errorf (%%+v): %+v\n)
	fmt.Printf("  是 ErrOriginalStd 吗？%t\n", standardErrors.Is(wrappedStdByStd, ErrOriginalStd))
	// (Is ErrOriginalStd? %t\n)

	fmt.Println("\n--- 包装 pkg/errors 错误 --- (--- Wrapping a pkg/errors error ---)")
	// 由 pkgErrors.Wrap 包装的 pkg/errors 错误
	// (pkg/errors error wrapped by pkgErrors.Wrap)
	wrappedPkgByPkg := wrapWithPkgErrors(ErrOriginalPkg, "由 pkgErrors.Wrap 包装的 pkg/errors (pkg/errors wrapped by pkgErrors.Wrap)")
	fmt.Printf("由 pkgErrors.Wrap 包装的 pkg (%%v): %v\n", wrappedPkgByPkg)
	// (pkg wrapped by pkgErrors.Wrap (%%v): %v\n)
	fmt.Printf("由 pkgErrors.Wrap 包装的 pkg (%%+v):\n%+v\n", wrappedPkgByPkg) // 显示来自 ErrOriginalPkg 的堆栈跟踪 (Stack trace from ErrOriginalPkg is shown)
	// (pkg wrapped by pkgErrors.Wrap (%%+v):\n%+v\n)
	fmt.Printf("  是 ErrOriginalPkg 吗？%t\n", standardErrors.Is(wrappedPkgByPkg, ErrOriginalPkg))
	// (Is ErrOriginalPkg? %t\n)

	fmt.Println("\n--- 互操作性：使用 fmt.Errorf %%w 包装 pkg/errors 错误 --- (--- Interoperability: Wrapping a pkg/errors error with fmt.Errorf %%w ---)")
	// 由 fmt.Errorf %%w 包装的 pkg/errors 错误
	// (pkg/errors error wrapped by fmt.Errorf %%w)
	wrappedPkgByStd := wrapWithStdFmt(ErrOriginalPkg, "由 fmt.Errorf 包装的 pkg/errors (pkg/errors wrapped by fmt.Errorf)")
	fmt.Printf("由 fmt.Errorf 包装的 pkg (%%v): %v\n", wrappedPkgByStd)
	// (pkg wrapped by fmt.Errorf (%%v): %v\n)
	// fmt.Errorf 不知道如何格式化 pkgErrors 以便 %%+v 显示来自原因的堆栈跟踪，
	// (fmt.Errorf does not know how to format pkgErrors for %%+v to show stack trace from cause,)
	// 除非 pkgErrors.fundamental 本身为此实现了特定的 Formatter 逻辑。
	// (unless pkgErrors.fundamental itself implemented a specific Formatter logic for this.)
	// 我们的 pkgErrors.fundamental.Format 处理其自身的堆栈跟踪，而不一定在由外部 fmt.Errorf 包装时处理。
	// (Our pkgErrors.fundamental.Format handles its own stack trace, not necessarily when wrapped by external fmt.Errorf.)
	fmt.Printf("由 fmt.Errorf 包装的 pkg (%%+v): %+v\n", wrappedPkgByStd) 
	// (pkg wrapped by fmt.Errorf (%%+v): %+v\n)
	fmt.Printf("  是 ErrOriginalPkg 吗？%t\n", standardErrors.Is(wrappedPkgByStd, ErrOriginalPkg))
	// (Is ErrOriginalPkg? %t\n)


	fmt.Println("\n--- 互操作性：使用 pkgErrors.Wrap 包装标准库错误 --- (--- Interoperability: Wrapping a standard library error with pkgErrors.Wrap ---)")
	// 由 pkgErrors.Wrap 包装的标准库错误 (os.ErrNotExist)
	// (standard library error (os.ErrNotExist) wrapped by pkgErrors.Wrap)
	originalStdLibError := os.ErrNotExist
	wrappedStdByPkg := wrapWithPkgErrors(originalStdLibError, "由 pkgErrors.Wrap 包装的 os.ErrNotExist (os.ErrNotExist wrapped by pkgErrors.Wrap)")
	fmt.Printf("由 pkg 包装的标准库 (os.ErrNotExist) (%%v): %v\n", wrappedStdByPkg)
	// (std (os.ErrNotExist) wrapped by pkg (%%v): %v\n)
	// 如果原因没有堆栈跟踪，pkgErrors.Wrap 会在包装点添加一个堆栈跟踪。
	// (pkgErrors.Wrap adds a stack trace at the point of wrapping if the cause doesn\'t have one.)
	// 由于 os.ErrNotExist 没有 pkgErrors 可识别的堆栈跟踪，因此 pkgErrors.Wrap 会创建一个。
	// (Since os.ErrNotExist doesn\'t have a stack trace that pkgErrors recognizes, pkgErrors.Wrap creates one.)
	fmt.Printf("由 pkg 包装的标准库 (os.ErrNotExist) (%%+v):\n%+v\n", wrappedStdByPkg) 
	// (std (os.ErrNotExist) wrapped by pkg (%%+v):\n%+v\n)
	fmt.Printf("  是 os.ErrNotExist 吗？%t\n", standardErrors.Is(wrappedStdByPkg, os.ErrNotExist))
	// (Is os.ErrNotExist? %t\n)

	fmt.Println("\n--- 访问底层错误 (Cause) --- (--- Accessing underlying error (Cause) ---)")
	fmt.Printf("wrappedPkgByPkg 的原因：%v\n", pkgErrors.Cause(wrappedPkgByPkg))
	// (Cause of wrappedPkgByPkg: %v\n)
	// 也可以使用 standardErrors.Unwrap
	// (standardErrors.Unwrap can also be used)
	fmt.Printf("wrappedPkgByPkg 的 Unwrap：%v\n", standardErrors.Unwrap(wrappedPkgByPkg))
	// (Unwrap of wrappedPkgByPkg: %v\n)
}

/*
示例输出 (堆栈跟踪会因情况而异)：
(Example Output (Stack traces will vary)):

--- 包装标准库错误 --- (--- Wrapping a standard library error ---)
由 fmt.Errorf 包装的标准库 (%%v): 由 fmt.Errorf 包装的标准库 (std lib wrapped by fmt.Errorf): 原始标准错误 (original standard error)
(Std wrapped by fmt.Errorf (%%v): std lib wrapped by fmt.Errorf: original standard error)
由 fmt.Errorf 包装的标准库 (%%+v): 由 fmt.Errorf 包装的标准库 (std lib wrapped by fmt.Errorf): 原始标准错误 (original standard error)
(Std wrapped by fmt.Errorf (%%+v): std lib wrapped by fmt.Errorf: original standard error)
  是 ErrOriginalStd 吗？true
  (Is ErrOriginalStd? true)

--- 包装 pkg/errors 错误 --- (--- Wrapping a pkg/errors error ---)
由 pkgErrors.Wrap 包装的 pkg (%%v): 由 pkgErrors.Wrap 包装的 pkg/errors (pkg/errors wrapped by pkgErrors.Wrap): 原始 pkg/errors 错误 (original pkg/errors error)
(pkg wrapped by pkgErrors.Wrap (%%v): pkg/errors wrapped by pkgErrors.Wrap: original pkg/errors error)
由 pkgErrors.Wrap 包装的 pkg (%%+v):
(pkg wrapped by pkgErrors.Wrap (%%+v):)
由 pkgErrors.Wrap 包装的 pkg/errors (pkg/errors wrapped by pkgErrors.Wrap): 原始 pkg/errors 错误 (original pkg/errors error)
main.main
	/path/to/your/file.go:20
runtime.main
	...
runtime.goexit
	...
  是 ErrOriginalPkg 吗？true
  (Is ErrOriginalPkg? true)

--- 互操作性：使用 fmt.Errorf %%w 包装 pkg/errors 错误 --- (--- Interoperability: Wrapping a pkg/errors error with fmt.Errorf %%w ---)
由 fmt.Errorf 包装的 pkg (%%v): 由 fmt.Errorf 包装的 pkg/errors (pkg/errors wrapped by fmt.Errorf): 原始 pkg/errors 错误 (original pkg/errors error)
(pkg wrapped by fmt.Errorf (%%v): pkg/errors wrapped by fmt.Errorf: original pkg/errors error)
由 fmt.Errorf 包装的 pkg (%%+v): 由 fmt.Errorf 包装的 pkg/errors (pkg/errors wrapped by fmt.Errorf): 原始 pkg/errors 错误 (original pkg/errors error)
(pkg wrapped by fmt.Errorf (%%+v): pkg/errors wrapped by fmt.Errorf: original pkg/errors error)
  是 ErrOriginalPkg 吗？true
  (Is ErrOriginalPkg? true)

--- 互操作性：使用 pkgErrors.Wrap 包装标准库错误 --- (--- Interoperability: Wrapping a standard library error with pkgErrors.Wrap ---)
由 pkg 包装的标准库 (os.ErrNotExist) (%%v): 由 pkgErrors.Wrap 包装的 os.ErrNotExist (os.ErrNotExist wrapped by pkgErrors.Wrap): file does not exist
(std (os.ErrNotExist) wrapped by pkg (%%v): os.ErrNotExist wrapped by pkgErrors.Wrap: file does not exist)
由 pkg 包装的标准库 (os.ErrNotExist) (%%+v):
(std (os.ErrNotExist) wrapped by pkg (%%+v):)
由 pkgErrors.Wrap 包装的 os.ErrNotExist (os.ErrNotExist wrapped by pkgErrors.Wrap): file does not exist
main.main
	/path/to/your/file.go:40
runtime.main
	...
runtime.goexit
	...
  是 os.ErrNotExist 吗？true
  (Is os.ErrNotExist? true)

--- 访问底层错误 (Cause) --- (--- Accessing underlying error (Cause) ---)
wrappedPkgByPkg 的原因：原始 pkg/errors 错误 (original pkg/errors error)
(Cause of wrappedPkgByPkg: original pkg/errors error)
wrappedPkgByPkg 的 Unwrap：原始 pkg/errors 错误 (original pkg/errors error)
(Unwrap of wrappedPkgByPkg: original pkg/errors error)
*/
```

**包装的关键要点：(Key takeaway for wrapping)** `pkgErrors.Wrap` 和 `pkgErrors.Wrapf` 是使用此库时包装错误的首选方法，因为它们可以确保正确的堆栈跟踪处理和上下文添加，同时保持与 `standardErrors.Is` 和 `standardErrors.As` 的兼容性。
( `pkgErrors.Wrap` and `pkgErrors.Wrapf` are the preferred way to wrap errors when using this library, as they ensure proper stack trace handling and context addition while maintaining compatibility with `standardErrors.Is` and `standardErrors.As`.) 