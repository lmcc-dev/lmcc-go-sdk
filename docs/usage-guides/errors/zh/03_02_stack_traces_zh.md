<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 堆栈跟踪 (Stack Traces)

所有由 `errors.New`、`errors.Errorf`、`errors.NewWithCode`、`errors.ErrorfWithCode`、`errors.Wrap` 和 `errors.Wrapf` 创建的错误都会在创建时自动捕获堆栈跟踪。由 `errors.WithCode` 包装的错误会保留原始错误的堆栈跟踪。

(All errors created by `errors.New`, `errors.Errorf`, `errors.NewWithCode`, `errors.ErrorfWithCode`, `errors.Wrap`, and `errors.Wrapf` automatically capture a stack trace at the point of their creation. Errors wrapped by `errors.WithCode` retain the stack trace of the original error.)

当使用 `fmt.Printf` 或类似函数以 `"%+v"` 格式化错误时，会打印堆栈跟踪。

(The stack trace is printed when the error is formatted with `"%+v"` using `fmt.Printf` or similar functions.)

```go
package main

import (
	"fmt"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
)

// functionA 调用 functionB
// (functionA calls functionB)
func functionA(triggerError bool) error {
	fmt.Println("正在执行 functionA... (Executing functionA...)")
	err := functionB(triggerError)
	if err != nil {
		// 包装来自 functionB 的错误，并添加来自 functionA 的上下文。
		// (Wrap the error from functionB, adding context from functionA.)
		// 保留 functionB 错误创建点的原始堆栈跟踪，
		// (The original stack trace from functionB\'s error creation point is preserved,)
		// 并且此包装点也是链的一部分。
		// (and this wrapping point is also part of the chain.)
		return errors.Wrap(err, "functionA 遇到一个问题 (functionA encountered an issue)")
	}
	return nil
}

// functionB 调用 functionC
// (functionB calls functionC)
func functionB(triggerError bool) error {
	fmt.Println("正在执行 functionB... (Executing functionB...)")
	err := functionC(triggerError)
	if err != nil {
		// 在这里，我们只是将错误向上传递。functionB 本身不会添加新的堆栈跟踪，
		// (Here, we are just passing the error up. No new stack trace is added by functionB itself)
		// 因为我们没有在这里创建新错误或使用 pkg/errors 函数包装它。
		// (because we are not creating a new error or wrapping it with pkg/errors functions here.)
		// 如果我们在这里使用 errors.Wrap()，那么这个点也会在跟踪中。
		// (If we were to use errors.Wrap() here, then this point would also be in the trace.)
		return err
	}
	return nil
}

// 如果 triggerError 为 true，functionC 会创建一个错误
// (functionC creates an error if triggerError is true)
func functionC(triggerError bool) error {
	fmt.Println("正在执行 functionC... (Executing functionC...)")
	if triggerError {
		// 使用 Coder 创建一个新错误。这是主堆栈跟踪的起源处。
		// (Create a new error with a Coder. This is where the primary stack trace will originate.)
		return errors.NewWithCode(errors.ErrInternalServer, "functionC 中模拟的关键故障 (simulated critical failure in functionC)")
	}
	fmt.Println("functionC 成功完成。(functionC completed successfully.)")
	return nil
}

func main() {
	fmt.Println("--- 场景1：未触发错误 --- (--- Scenario 1: No error triggered ---)")
	err1 := functionA(false)
	if err1 == nil {
		fmt.Println("场景1 无错误完成。(Scenario 1 completed without errors.)")
	}

	fmt.Println("\n--- 场景2：触发并传播错误 --- (--- Scenario 2: Error triggered and propagated ---)")
	err2 := functionA(true)
	if err2 != nil {
		fmt.Println("\n发生了一个错误：(An error occurred:)")
		fmt.Println("\n--- 使用 %v 打印错误 --- (--- Printing error with %v ---)")
		fmt.Printf("%v\n", err2)

		fmt.Println("\n--- 使用 %+v 打印错误 (包含堆栈跟踪) --- (--- Printing error with %+v (includes stack trace) ---)")
		fmt.Printf("%+v\n", err2) // 这将打印错误消息和堆栈跟踪。(This will print the error message and the stack trace.)
		
		// 演示 Cause() (Demonstrate Cause())
		originalCause := errors.Cause(err2)
		fmt.Println("\n--- 原始原因 (使用 errors.Cause) 与 %+v --- (--- Original cause (using errors.Cause) with %+v --- )")
		fmt.Printf("%+v\n", originalCause)
	}
}

/*
示例输出 (文件路径和行号会因您的环境而异)：
(Example Output (file paths and line numbers will vary based on your environment)):

--- 场景1：未触发错误 --- (--- Scenario 1: No error triggered ---)
正在执行 functionA... (Executing functionA...)
正在执行 functionB... (Executing functionB...)
正在执行 functionC... (Executing functionC...)
functionC 成功完成。(functionC completed successfully.)
场景1 无错误完成。(Scenario 1 completed without errors.)

--- 场景2：触发并传播错误 --- (--- Scenario 2: Error triggered and propagated ---)
正在执行 functionA... (Executing functionA...)
正在执行 functionB... (Executing functionB...)
正在执行 functionC... (Executing functionC...)

发生了一个错误：(An error occurred:)

--- 使用 %v 打印错误 --- (--- Printing error with %v ---)
functionA 遇到一个问题 (functionA encountered an issue): Internal server error: functionC 中模拟的关键故障 (simulated critical failure in functionC)

--- 使用 %+v 打印错误 (包含堆栈跟踪) --- (--- Printing error with %+v (includes stack trace) ---)
functionA 遇到一个问题 (functionA encountered an issue): functionC 中模拟的关键故障 (simulated critical failure in functionC)
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

--- 原始原因 (使用 errors.Cause) 与 %+v --- (--- Original cause (using errors.Cause) with %+v --- )
functionC 中模拟的关键故障 (simulated critical failure in functionC)
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

**关于堆栈跟踪和包装的说明：**
(Note on Stack Traces and Wrapping:)
- 当您使用 `errors.Wrap` 或 `errors.Wrapf` 包装错误时，会保留原始错误的堆栈跟踪 (如果它是由 `pkg/errors` 创建的)。`Wrap` 调用本身会向概念上的消息堆栈添加一个新帧，但不会生成一个*新的*完整堆栈跟踪；它将错误链接起来。
  (When you wrap an error using `errors.Wrap` or `errors.Wrapf`, the original error\'s stack trace (if it was created by `pkg/errors`) is preserved. The `Wrap` call itself adds a new frame to the conceptual stack of messages but doesn\'t generate a *new* full stack trace; it chains the errors.)
- `errors.WithCode` 也会保留原始错误的堆栈跟踪。
  (errors.WithCode also preserves the original error\'s stack trace.)
- 使用 `%+v` 格式化通常会显示从*最内层*错误 (原因) 由 `pkg/errors` 函数创建点开始的堆栈跟踪，然后是来自包装错误的消息。
  (Formatting with `%+v` will typically show the stack trace from the point where the *innermost* error (the cause) was created by a `pkg/errors` function, followed by the messages from wrapping errors.)

```go
</rewritten_file> 