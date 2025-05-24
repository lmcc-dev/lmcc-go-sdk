<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 1. 替换 `errors.New` 和 `fmt.Errorf` (Replacing `errors.New` and `fmt.Errorf`)

- 将 `errors.New("message")` 替换为 `pkgErrors.New("message")`。
  (Replace `errors.New("message")` with `pkgErrors.New("message")`.)
- 将 `fmt.Errorf("format %s", var)` 替换为 `pkgErrors.Errorf("format %s", var)`。
  (Replace `fmt.Errorf("format %s", var)` with `pkgErrors.Errorf("format %s", var)`.)

`pkg/errors` 版本会自动捕获堆栈跟踪。

(The `pkg/errors` versions will automatically capture stack traces.)

```go
package main

import (
	"fmt"
	standardErrors "errors" // Go 标准库 errors (Go standard library errors)
	pkgErrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors" // 我们的 pkg/errors (Our pkg/errors)
)

// oldFunction 使用标准库 errors
// (oldFunction uses standard library errors)
func oldFunction(succeed bool) error {
	if !succeed {
		// 标准库 errors.New
		// (Standard library errors.New)
		return standardErrors.New("oldFunction 因静态原因失败 (oldFunction failed due to a static reason)")
	}
	return nil
}

// anotherOldFunction 使用标准库 fmt.Errorf
// (anotherOldFunction uses standard library fmt.Errorf)
func anotherOldFunction(value int) error {
	if value < 0 {
		// 标准库 fmt.Errorf
		// (Standard library fmt.Errorf)
		return fmt.Errorf("anotherOldFunction 失败：无效值 %d (anotherOldFunction failed: invalid value %d)", value)
	}
	return nil
}

// newFunction 使用 pkg/errors.New
// (newFunction uses pkg/errors.New)
func newFunction(succeed bool) error {
	if !succeed {
		// pkg/errors.New - 捕获堆栈跟踪
		// (pkg/errors.New - captures stack trace)
		return pkgErrors.New("newFunction 因静态原因失败 (newFunction failed due to a static reason)")
	}
	return nil
}

// anotherNewFunction 使用 pkg/errors.Errorf
// (anotherNewFunction uses pkg/errors.Errorf)
func anotherNewFunction(value int) error {
	if value < 0 {
		// pkg/errors.Errorf - 捕获堆栈跟踪
		// (pkg/errors.Errorf - captures stack trace)
		return pkgErrors.Errorf("anotherNewFunction 失败：无效值 %d (anotherNewFunction failed: invalid value %d)", value)
	}
	return nil
}

func main() {
	fmt.Println("--- 标准库错误 (%%v 不会自动打印堆栈跟踪) --- (--- Standard Library Errors (No Automatic Stack Trace with %%v) ---)")
	errOld1 := oldFunction(false)
	if errOld1 != nil {
		fmt.Printf("oldFunction 错误 (%%v): %v\n", errOld1)
		// (oldFunction error (%%v): %v\n)
		// 标准错误除非实现了特定的 Formatter 接口 (errors.New 没有实现)，
		// (Standard errors don\'t automatically print stack trace with %%+v 
		// 否则使用 %%+v 不会自动打印堆栈跟踪。
		// unless they implement a specific Formatter interface, which errors.New doesn\'t.)
		fmt.Printf("oldFunction 错误 (%%+v): %+v\n", errOld1) 
		// (oldFunction error (%%+v): %+v\n)
	}

	errOld2 := anotherOldFunction(-10)
	if errOld2 != nil {
		fmt.Printf("anotherOldFunction 错误 (%%v): %v\n", errOld2)
		// (anotherOldFunction error (%%v): %v\n)
		fmt.Printf("anotherOldFunction 错误 (%%+v): %+v\n", errOld2)
		// (anotherOldFunction error (%%+v): %+v\n)
	}

	fmt.Println("\n--- pkg/errors (%%+v 会自动打印堆栈跟踪) --- (--- pkg/errors (With Automatic Stack Trace with %%+v) ---)")
	errNew1 := newFunction(false)
	if errNew1 != nil {
		fmt.Printf("newFunction 错误 (%%v): %v\n", errNew1)
		// (newFunction error (%%v): %v\n)
		fmt.Printf("newFunction 错误 (%%+v):\n%+v\n", errNew1) // 会显示堆栈跟踪 (Will show stack trace)
		// (newFunction error (%%+v):\n%+v\n)
	}

	errNew2 := anotherNewFunction(-20)
	if errNew2 != nil {
		fmt.Printf("anotherNewFunction 错误 (%%v): %v\n", errNew2)
		// (anotherNewFunction error (%%v): %v\n)
		fmt.Printf("anotherNewFunction 错误 (%%+v):\n%+v\n", errNew2) // 会显示堆栈跟踪 (Will show stack trace)
		// (anotherNewFunction error (%%+v):\n%+v\n)
	}
}

/*
示例输出 (堆栈跟踪会因您的环境而异)：
(Example Output (Stack traces will vary based on your environment)):

--- 标准库错误 (%%v 不会自动打印堆栈跟踪) --- (--- Standard Library Errors (No Automatic Stack Trace with %%v) ---)
oldFunction 错误 (%%v): oldFunction 因静态原因失败 (oldFunction failed due to a static reason)
(oldFunction error (%%v): oldFunction failed due to a static reason)
oldFunction 错误 (%%+v): oldFunction 因静态原因失败 (oldFunction failed due to a static reason)
(oldFunction error (%%+v): oldFunction failed due to a static reason)
anotherOldFunction 错误 (%%v): anotherOldFunction 失败：无效值 -10 (anotherOldFunction failed: invalid value -10)
(anotherOldFunction error (%%v): anotherOldFunction failed: invalid value -10)
anotherOldFunction 错误 (%%+v): anotherOldFunction 失败：无效值 -10 (anotherOldFunction failed: invalid value -10)
(anotherOldFunction error (%%+v): anotherOldFunction failed: invalid value -10)

--- pkg/errors (%%+v 会自动打印堆栈跟踪) --- (--- pkg/errors (With Automatic Stack Trace with %%+v) ---)
newFunction 错误 (%%v): newFunction 因静态原因失败 (newFunction failed due to a static reason)
(newFunction error (%%v): newFunction failed due to a static reason)
newFunction 错误 (%%+v):
(newFunction error (%%+v):)
newFunction 因静态原因失败 (newFunction failed due to a static reason)
main.newFunction
	/path/to/your/file.go:30
main.main
	/path/to/your/file.go:51
runtime.main
	...
runtime.goexit
	...
anotherNewFunction 错误 (%%v): anotherNewFunction 失败：无效值 -20 (anotherNewFunction failed: invalid value -20)
(anotherNewFunction error (%%v): anotherNewFunction failed: invalid value -20)
anotherNewFunction 错误 (%%+v):
(anotherNewFunction error (%%+v):)
anotherNewFunction 失败：无效值 -20 (anotherNewFunction failed: invalid value -20)
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

</rewritten_file> 