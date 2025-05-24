<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 2. 添加上下文 (Adding Context)

当从一个函数返回错误时，使用 `errors.Wrap` 或 `errors.Wrapf` 来添加关于调用函数试图做什么的上下文信息。这有助于在调试时理解错误的路径和来源。

(When an error is returned from a function call, use `errors.Wrap` or `errors.Wrapf` to add contextual information about what the calling function was trying to do. This helps in understanding the error's path and origin when debugging.)

- **`errors.Wrap(err error, message string) error`**
  用给定的消息包装错误 `err`。如果 `err` 为 `nil`，则返回 `nil`。
  (Wraps error `err` with the given message. Returns `nil` if `err` is `nil`.)

- **`errors.Wrapf(err error, format string, args ...interface{}) error`**
  用格式化的消息包装错误 `err`。如果 `err` 为 `nil`，则返回 `nil`。
  (Wraps error `err` with a formatted message. Returns `nil` if `err` is `nil`.)

- **`errors.WithMessage(err error, message string) error`**
  `Wrap` 的别名。
  (Alias for `Wrap`.)

- **`errors.WithMessagef(err error, format string, args ...interface{}) error`**
  `Wrapf` 的别名。
  (Alias for `Wrapf`.)

```go
package main

import (
	"fmt"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"os" // 用于一个示例标准库错误 (For an example standard library error)
)

// readFileContent 模拟从文件读取内容，可能会失败。
// (readFileContent simulates reading content from a file, which might fail.)
func readFileContent(filePath string) (string, error) {
	if filePath == "nonexistent.txt" {
		// 模拟一个来自标准库的错误 (Simulate an error from the standard library)
		return "", os.ErrNotExist // 例如：文件未找到 (e.g., file not found)
	}
	if filePath == "unreadable.txt" {
		// 模拟另一个使用 pkg/errors 创建的底层错误
		// (Simulate another low-level error created with pkg/errors)
		return "", errors.New("底层读取权限被拒绝 (underlying read permission denied)")
	}
	return fmt.Sprintf("文件 '%s' 的内容 (Content of file '%s')", filePath, filePath), nil
}

// processDataFromFile 读取文件并处理其内容。
// (processDataFromFile reads a file and processes its content.)
func processDataFromFile(filePath string) error {
	content, err := readFileContent(filePath)
	if err != nil {
		// 使用 errors.Wrap 添加上下文到来自 readFileContent 的错误。
		// (Use errors.Wrap to add context to the error from readFileContent.)
		return errors.Wrap(err, fmt.Sprintf("处理文件数据失败 '%s' (failed to process data from file '%s')", filePath, filePath))
	}

	// 假设我们在这里处理内容 (Assume we process the content here)
	fmt.Printf("成功处理来自 '%s' 的数据: %s (Successfully processed data from '%s': %s)\n", filePath, content, filePath, content)
	return nil
}

// higherLevelOperation 调用 processDataFromFile 并可能添加更多上下文。
// (higherLevelOperation calls processDataFromFile and might add more context.)
func higherLevelOperation(filePath string, operationID string) error {
	err := processDataFromFile(filePath)
	if err != nil {
		// 使用 errors.Wrapf 添加更多格式化的上下文。
		// (Use errors.Wrapf to add more formatted context.)
		return errors.Wrapf(err, "高级操作 '%s' 在处理文件 '%s' 时失败 (High-level operation '%s' failed when processing file '%s')", operationID, filePath, operationID, filePath)
	}
	fmt.Printf("高级操作 '%s' 已成功完成。(High-level operation '%s' completed successfully.)\n", operationID, operationID)
	return nil
}

func main() {
	fmt.Println("--- 场景1：文件未找到 (os.ErrNotExist 被包装) --- (--- Scenario 1: File Not Found (os.ErrNotExist is wrapped) ---)")
	err1 := higherLevelOperation("nonexistent.txt", "OP001")
	if err1 != nil {
		fmt.Printf("错误 (Error) (%%v): %v\n", err1)
		fmt.Printf("错误 (Error) (%%+v):\n%+v\n", err1) // os.ErrNotExist 本身没有堆栈，但 Wrap 会添加 (os.ErrNotExist itself has no stack, but Wrap adds one)
		// 检查原始错误 (Check for the original error)
		if errors.Is(err1, os.ErrNotExist) {
			fmt.Println("根本原因确实是 os.ErrNotExist (The root cause is indeed os.ErrNotExist)")
		}
	}

	fmt.Println("\n--- 场景2：底层 pkg/errors 错误被包装 --- (--- Scenario 2: Underlying pkg/errors error is wrapped ---)")
	err2 := higherLevelOperation("unreadable.txt", "OP002")
	if err2 != nil {
		fmt.Printf("错误 (Error) (%%v): %v\n", err2)
		fmt.Printf("错误 (Error) (%%+v):\n%+v\n", err2) // 应该显示原始 pkg/errors 错误的堆栈跟踪 (Should show stack trace from original pkg/errors error)
		// 检查根本原因 (Check the root cause)
		cause := errors.Cause(err2)
		fmt.Printf("根本原因 (Root cause): %v\n", cause)
	}

	fmt.Println("\n--- 场景3：成功案例 --- (--- Scenario 3: Success Case ---)")
	err3 := higherLevelOperation("readable_file.txt", "OP003")
	if err3 == nil {
		fmt.Println("操作 OP003 成功完成！(Operation OP003 completed successfully!)")
	}
}

/*
预期输出 (Expected Output - 堆栈跟踪的路径和行号会因您的环境而异):

--- 场景1：文件未找到 (os.ErrNotExist 被包装) --- (--- Scenario 1: File Not Found (os.ErrNotExist is wrapped) ---)
错误 (Error) (%v): 高级操作 'OP001' 在处理文件 'nonexistent.txt' 时失败 (High-level operation 'OP001' failed when processing file 'nonexistent.txt'): 处理文件数据失败 'nonexistent.txt' (failed to process data from file 'nonexistent.txt'): file does not exist
错误 (Error) (%+v):
高级操作 'OP001' 在处理文件 'nonexistent.txt' 时失败 (High-level operation 'OP001' failed when processing file 'nonexistent.txt'): 处理文件数据失败 'nonexistent.txt' (failed to process data from file 'nonexistent.txt'): file does not exist
main.processDataFromFile
	/path/to/your/file.go:XX
main.higherLevelOperation
	/path/to/your/file.go:YY
main.main
	/path/to/your/file.go:ZZ
...
根本原因确实是 os.ErrNotExist (The root cause is indeed os.ErrNotExist)

--- 场景2：底层 pkg/errors 错误被包装 --- (--- Scenario 2: Underlying pkg/errors error is wrapped ---)
错误 (Error) (%v): 高级操作 'OP002' 在处理文件 'unreadable.txt' 时失败 (High-level operation 'OP002' failed when processing file 'unreadable.txt'): 处理文件数据失败 'unreadable.txt' (failed to process data from file 'unreadable.txt'): 底层读取权限被拒绝 (underlying read permission denied)
错误 (Error) (%+v):
高级操作 'OP002' 在处理文件 'unreadable.txt' 时失败 (High-level operation 'OP002' failed when processing file 'unreadable.txt'): 处理文件数据失败 'unreadable.txt' (failed to process data from file 'unreadable.txt'): 底层读取权限被拒绝 (underlying read permission denied)
main.readFileContent
	/path/to/your/file.go:AA
main.processDataFromFile
	/path/to/your/file.go:BB
main.higherLevelOperation
	/path/to/your/file.go:CC
main.main
	/path/to/your/file.go:DD
...
根本原因 (Root cause): 底层读取权限被拒绝 (underlying read permission denied)

--- 场景3：成功案例 --- (--- Scenario 3: Success Case ---)
成功处理来自 'readable_file.txt' 的数据: 文件 'readable_file.txt' 的内容 (Successfully processed data from 'readable_file.txt': Content of file 'readable_file.txt')
高级操作 'OP003' 已成功完成。(High-level operation 'OP003' completed successfully.)
操作 OP003 成功完成！(Operation OP003 completed successfully!)
*/
``` 