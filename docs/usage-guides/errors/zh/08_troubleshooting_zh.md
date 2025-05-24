<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

## 问题排查 (Troubleshooting)

使用 `pkg/errors` 时遇到的常见问题及其解决方法。

(Common issues and how to address them when working with `pkg/errors`.)

**1. 堆栈跟踪未出现 (Stack Trace Not Appearing)**
   - **问题 (Issue)**: 您使用 `%+v` 打印错误，但没有看到堆栈跟踪。
     (You print an error with `%+v` but don't see a stack trace.)
   - **可能的原因和解决方案 (Possible Causes & Solutions)**:
     - 错误不是由 `pkg/errors` 创建的 (例如，它是标准库错误，如 `io.EOF` 或来自标准库的 `errors.New()`，并且后续没有被 `pkg/errors` 的函数如 `pkgErrors.Wrap` 或 `pkgErrors.WithCode` 包装)。
       (The error was not created by `pkg/errors` (e.g., it's a standard library error like `io.EOF` or from `errors.New()` from the standard library, and was not subsequently wrapped by `pkg/errors` functions like `pkgErrors.Wrap` or `pkgErrors.WithCode`).)
     - 您正在使用 `%v` 或 `%s` 而不是 `%+v` 打印错误。
       (You are printing the error with `%v` or `%s` instead of `%+v`.)
     - 错误是使用 `pkgErrors.WithCode(nil, someCoder)` 创建的。对 nil 错误使用 `WithCode` 会返回 nil。
       (The error was created with `pkgErrors.WithCode(nil, someCoder)`. `WithCode` on a nil error returns nil.)
     - 如果一个来自 `pkg/errors` 的错误被标准库的 `fmt.Errorf("... %w", pkgErr)` 包装，那么 `fmt.Printf("%+v", wrappedErr)` 可能不会显示来自 `pkgErr` 的堆栈跟踪，除非 `pkgErr` 本身实现了一个 `Format` 方法，`fmt.Errorf` 知道如何通过链为 `%+v` 动词调用该方法。我们的 `pkg/errors` `fundamental` 类型确实实现了 `Format`，这应该使其堆栈跟踪在作为原因时可见。
       (If an error from `pkg/errors` is wrapped by the standard library's `fmt.Errorf("... %w", pkgErr)`, then `fmt.Printf("%+v", wrappedErr)` might not display the stack trace from `pkgErr` unless `pkgErr` itself implements a `Format` method that `fmt.Errorf` knows how to call for the `%+v` verb through the chain. Our `pkg/errors` `fundamental` type does implement `Format`, which should make its stack trace visible when it's the cause.)

   ```go
   package main

   import (
   	"fmt"
   	standardErrors "errors"
   	pkgErrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
   	"io"
   )

   func main() {
   	// 原因1：标准库错误，未被 pkg/errors 包装
   	// (Cause 1: Standard library error, not wrapped by pkg/errors)
   	stdErr := io.EOF
   	fmt.Println("--- 标准错误 (io.EOF) --- (--- Standard Error (io.EOF) ---)")
   	fmt.Printf("%%v: %v\n", stdErr)
   	fmt.Printf("%%+v: %+v\n\n", stdErr) // io.EOF 本身没有堆栈跟踪 (No stack trace from io.EOF itself)

   	// 被 pkgErrors 包装 - 现在它有了堆栈 (从包装点开始)
   	// (Wrapped by pkgErrors - NOW it gets a stack (from the wrapping point))
   	wrappedStdErr := pkgErrors.Wrap(stdErr, "读取文件失败 (failed to read file)")
   	fmt.Println("--- 标准错误被 pkgErrors.Wrap 包装 --- (--- Standard Error Wrapped by pkgErrors.Wrap ---)")
   	fmt.Printf("%%v: %v\n", wrappedStdErr)
   	fmt.Printf("%%+v:\n%+v\n\n", wrappedStdErr) 

   	// 原因2：对 pkg/errors 错误使用 %v 而不是 %+v
   	// (Cause 2: Using %v instead of %+v for a pkg/errors error)
   	pkgErr := pkgErrors.New("一个 pkg/errors 错误 (a pkg/errors error)")
   	fmt.Println("--- pkg/errors 错误 --- (--- pkg/errors Error ---)")
   	fmt.Printf("%%v: %v\n", pkgErr)       // 没有堆栈跟踪 (No stack trace)
   	fmt.Printf("%%+v:\n%+v\n\n", pkgErr) // 出现堆栈跟踪 (Stack trace appears)
   	
   	// 原因3：fmt.Errorf 包装 pkgError
   	// (Cause 3: fmt.Errorf wrapping a pkgError)
   	fmtWrappedPkgErr := fmt.Errorf("由 fmt.Errorf 包装：%w (wrapped by fmt.Errorf: %w)", pkgErr)
   	fmt.Println("--- pkg/errors 错误被 fmt.Errorf 包装 --- (--- pkg/errors Error Wrapped by fmt.Errorf ---) ")
   	fmt.Printf("%%v: %v\n", fmtWrappedPkgErr)
   	// pkg/errors 中的 fundamental.Format 方法应该允许即使在这里也显示堆栈跟踪。
   	// (The fundamental.Format method in pkg/errors should allow stack trace to be shown even here.)
   	fmt.Printf("%%+v:\n%+v\n\n", fmtWrappedPkgErr) 
   }
   /* 输出 (Output):
   --- 标准错误 (io.EOF) --- (--- Standard Error (io.EOF) ---)
   %v: EOF
   %+v: EOF

   --- 标准错误被 pkgErrors.Wrap 包装 --- (--- Standard Error Wrapped by pkgErrors.Wrap ---)
   %v: 读取文件失败 (failed to read file): EOF
   %+v:
   读取文件失败 (failed to read file): EOF
   main.main
       /path/to/file.go:XX
   ...

   --- pkg/errors 错误 --- (--- pkg/errors Error ---)
   %v: 一个 pkg/errors 错误 (a pkg/errors error)
   %+v:
   一个 pkg/errors 错误 (a pkg/errors error)
   main.main
       /path/to/file.go:XX
   ...

   --- pkg/errors 错误被 fmt.Errorf 包装 --- (--- pkg/errors Error Wrapped by fmt.Errorf ---) 
   %v: 由 fmt.Errorf 包装：一个 pkg/errors 错误 (wrapped by fmt.Errorf: a pkg/errors error)
   %+v:
   由 fmt.Errorf 包装：一个 pkg/errors 错误 (wrapped by fmt.Errorf: a pkg/errors error)
   main.main
       /path/to/file.go:XX
   ...
   */
   ```

**2. `errors.Is` 或 `errors.As` 未按预期工作 ( `errors.Is` or `errors.As` Not Working as Expected)**
   - **问题 (Issue)**: `standardErrors.Is(err, target)` 返回 `false`，或者 `standardErrors.As(err, &targetType)` 没有找到类型，即使您认为它应该找到。
     ( `standardErrors.Is(err, target)` returns `false`, or `standardErrors.As(err, &targetType)` doesn't find the type, even though you think it should.)
   - **可能的原因和解决方案 (Possible Causes & Solutions)**:
     - **`Is` 与不同的实例 (`Is` with different instances)**: 如果 `target` 是在不同位置创建的非 nil 错误 (例如，`err = pkgErrors.New("msg"); target = pkgErrors.New("msg")`)，`errors.Is` 将为 `false`，因为它们是不同的实例。`errors.Is` 检查引用是否相等，或者链中的错误是否通过 `Is(error) bool` 方法报告自身等效。
       (If `target` is a non-nil error created at a different place (e.g., `err = pkgErrors.New("msg"); target = pkgErrors.New("msg")`), `errors.Is` will be `false` because they are different instances. `errors.Is` checks for reference equality or if an error in the chain reports itself as equivalent via an `Is(error) bool` method.)
       - 对于来自 `pkg/errors` 的 `Coder` 类型，检查 `errors.Is(err, pkgErrors.ErrNotFound)` 是有效的，因为 `pkgErrors.ErrNotFound` 是一个特定的全局变量 (哨兵)。
         (For `Coder` types from `pkg/errors`, checking `errors.Is(err, pkgErrors.ErrNotFound)` works because `pkgErrors.ErrNotFound` is a specific global variable (sentinel).)
       - 要按代码检查常规类别，请使用 `pkgErrors.IsCode(err, CoderWithTheCodeYouWant)`。
         (For checking general categories by code, use `pkgErrors.IsCode(err, CoderWithTheCodeYouWant)`.)
     - **`As` 类型错误 (`As` with wrong type)**: 确保传递给 `errors.As` 的 `targetType` 变量是指向接口类型的指针或指向实现 `error` 的具体类型的指针。
       (Ensure the `targetType` variable you pass to `errors.As` is a pointer to an interface type or a pointer to a concrete type that implements `error`.)
     - **错误不在链中 (Error not in chain)**: 您正在检查的特定错误实例或类型实际上可能不在 `err` 的包装错误链中。
       (The specific error instance or type you're checking for might not actually be in `err`'s chain of wrapped errors.)

   ```go
   package main

   import (
   	"fmt"
   	standardErrors "errors"
   	pkgErrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
   	"io"
   	"os" // 导入 os 包以供 os.PathError 使用 (Import os package for os.PathError)
   )

   type MyTroubleError struct{ Msg string }
   func (e *MyTroubleError) Error() string { return e.Msg }

   var SentinelError = pkgErrors.New("特定的哨兵错误实例 (Specific sentinel error instance)")

   func main() {
   	// --- errors.Is 问题 --- (--- errors.Is issues ---)
   	err1 := pkgErrors.Wrap(SentinelError, "包装的哨兵 (wrapped sentinel)")
   	fmt.Printf("err1 是 SentinelError 吗？%t\n", standardErrors.Is(err1, SentinelError)) // true (Is err1 SentinelError? %t\n)

   	nonSentinelPkgErr := pkgErrors.New("一些消息 (some message)")
   	anotherNonSentinel := pkgErrors.New("一些消息 (some message)")
   	fmt.Printf("nonSentinelPkgErr 是 anotherNonSentinel 吗？%t\n", standardErrors.Is(nonSentinelPkgErr, anotherNonSentinel)) // false，不同实例 (false, different instances)

   	wrappedIoEOF := pkgErrors.Wrap(io.EOF, "包装的 io.EOF (wrapped io.EOF)")
   	fmt.Printf("wrappedIoEOF 是 io.EOF 吗？%t\n", standardErrors.Is(wrappedIoEOF, io.EOF)) // true (Is wrappedIoEOF io.EOF? %t\n)

   	// --- errors.As 问题 --- (--- errors.As issues ---)
   	customErrInstance := &MyTroubleError{"自定义故障 (custom trouble)"}
   	errWithCustom := pkgErrors.Wrap(customErrInstance, "上下文 (context)")

   	var target *MyTroubleError
   	if standardErrors.As(errWithCustom, &target) {
   		fmt.Printf("As 找到 MyTroubleError：%s\n", target.Msg) // (As found MyTroubleError: %s\n)
   	} else {
   		fmt.Println("As 未找到 MyTroubleError (As did not find MyTroubleError)")
   	}

   	var targetIOErr *os.PathError // 链中没有的类型示例 (Example of a type not in the chain)
   	if standardErrors.As(errWithCustom, &targetIOErr) {
   		fmt.Printf("As 找到 os.PathError：%s\n", targetIOErr.Path) // (As found os.PathError: %s\n)
   	} else {
   		fmt.Println("As 在 errWithCustom 中未找到 os.PathError (As did not find os.PathError in errWithCustom)")
   	}
   }
   /* 输出 (Output):
   err1 是 SentinelError 吗？true
   nonSentinelPkgErr 是 anotherNonSentinel 吗？false
   wrappedIoEOF 是 io.EOF 吗？true
   As 找到 MyTroubleError：自定义故障 (custom trouble)
   As 在 errWithCustom 中未找到 os.PathError (As did not find os.PathError in errWithCustom)")
   */
   ```

**3. `Coder` 信息未被检索 (`Coder` Information Not Being Retrieved)**
   - **问题 (Issue)**: `errors.GetCoder(err)` 返回 `nil`。
     (`errors.GetCoder(err)` returns `nil`.)
   - **可能的原因和解决方案 (Possible Causes & Solutions)**:
     - 链中没有错误是使用 `NewWithCode`、`ErrorfWithCode` 创建的，或者没有通过 `WithCode` 或 `WrapWithCode` 附加 `Coder`。
       (No error in the chain was created with `NewWithCode`, `ErrorfWithCode`, or had a `Coder` attached via `WithCode` or `WrapWithCode`.)
     - 错误链可能只包含标准库错误或其他不使用 `pkg/errors` `Coder` 机制的包中的错误。
       (The error chain might only contain standard library errors or errors from other packages that don't use the `pkg/errors` `Coder` mechanism.)

   ```go
   package main

   import (
   	"fmt"
   	standardErrors "errors"
   	pkgErrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
   )

   func main() {
   	// 带有 Coder 的错误
   	// (Error with a Coder)
   	errWithCode := pkgErrors.NewWithCode(pkgErrors.ErrBadRequest, "确实是错误的请求 (bad request indeed)")
   	coder1 := pkgErrors.GetCoder(errWithCode)
   	if coder1 != nil {
   		fmt.Printf("找到 Coder：Code=%d，Message='%s'\n", coder1.Code(), coder1.String()) // (Coder found: Code=%d, Message='%s'\n)
   	} else {
   		fmt.Println("在 errWithCode 中未找到 coder (No coder found in errWithCode)")
   	}

   	// 标准错误，没有 Coder
   	// (Standard error, no Coder)
   	stdErr := standardErrors.New("一个普通的标准错误 (a plain standard error)")
   	coder2 := pkgErrors.GetCoder(stdErr)
   	if coder2 != nil {
   		fmt.Printf("在 stdErr 中找到 Coder：Code=%d，Message='%s'\n", coder2.Code(), coder2.String()) // (Coder found in stdErr: Code=%d, Message='%s'\n)
   	} else {
   		fmt.Println("在 stdErr 中未找到 coder (No coder found in stdErr)")
   	}

   	// 标准错误被 pkgErrors.Wrap 包装 (原始错误仍然没有 Coder)
   	// (Standard error wrapped by pkgErrors.Wrap (still no Coder from original error))
   	wrappedStdErr := pkgErrors.Wrap(stdErr, "添加了上下文 (context added)")
   	coder3 := pkgErrors.GetCoder(wrappedStdErr)
   	if coder3 != nil {
   		fmt.Printf("在 wrappedStdErr 中找到 Coder：Code=%d，Message='%s'\n", coder3.Code(), coder3.String()) // (Coder found in wrappedStdErr: Code=%d, Message='%s'\n)
   	} else {
   		fmt.Println("在 wrappedStdErr 中未找到 coder (原始错误没有 Coder) (No coder found in wrappedStdErr (original error had no Coder))")
   	}
   	
   	// 标准错误被 pkgErrors.WrapWithCode 包装
   	// (Standard error wrapped by pkgErrors.WrapWithCode)
   	wrappedStdErrWithCode := pkgErrors.WrapWithCode(stdErr, pkgErrors.ErrInternalServer, "用代码包装 (wrapped with code)")
   	coder4 := pkgErrors.GetCoder(wrappedStdErrWithCode)
   	if coder4 != nil {
   		fmt.Printf("在 wrappedStdErrWithCode 中找到 Coder：Code=%d，Message='%s'\n", coder4.Code(), coder4.String()) // (Coder found in wrappedStdErrWithCode: Code=%d, Message='%s'\n)
   	} else {
   		fmt.Println("在 wrappedStdErrWithCode 中未找到 coder (No coder found in wrappedStdErrWithCode)")
   	}
   }
   /* 输出 (Output):
   找到 Coder：Code=100003，Message='Bad request'
   在 stdErr 中未找到 coder (No coder found in stdErr)
   在 wrappedStdErr 中未找到 coder (原始错误没有 Coder) (No coder found in wrappedStdErr (original error had no Coder))
   找到 Coder：Code=100001，Message='Internal server error'
   */
   ```

如果您遇到其他问题，请确保您按预期使用 `pkg/errors` 模块中的函数，并检查特定函数的文档以了解行为详细信息 (例如，如何处理 `nil` 错误)。
(If you encounter other issues, ensure you are using the functions from the `pkg/errors` module as intended and check the specific function documentation for behavior details (e.g., how `nil` errors are handled).) 