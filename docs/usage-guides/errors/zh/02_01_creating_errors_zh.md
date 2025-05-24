<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 1. 创建错误 (Creating Errors)

使用 `errors.New` 或 `errors.Errorf` 来创建错误。这些函数会自动捕获堆栈跟踪。

(Use `errors.New` or `errors.Errorf` to create errors. These functions automatically capture stack traces.)

- **`errors.New(message string) error`**
  创建一个包含给定消息的新错误。
  (Creates a new error with the given message.)

- **`errors.Errorf(format string, args ...interface{}) error`**
  根据格式说明符格式化并返回一个新错误。
  (Formats according to a format specifier and returns the string as a new error.)

```go
package main

import (
	"fmt"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
)

// simulateOperationA 演示 errors.New
// (simulateOperationA demonstrates errors.New)
func simulateOperationA(succeed bool) error {
	if !succeed {
		// 创建一个简单的错误，它将自动包含堆栈跟踪。
		// (Create a simple error, which will automatically include a stack trace.)
		return errors.New("操作 A 连接目标服务失败 (Operation A failed to connect to target service)")
	}
	fmt.Println("操作 A 成功完成 (Operation A completed successfully)")
	return nil
}

// simulateOperationB 演示 errors.Errorf
// (simulateOperationB demonstrates errors.Errorf)
func simulateOperationB(userID string, resourceID string) error {
	if userID == "" {
		// 创建一个格式化的错误，它也将自动包含堆栈跟踪。
		// (Create a formatted error, which will also automatically include a stack trace.)
		return errors.Errorf("处理资源 '%s' 失败：用户ID不能为空 (Failed to process resource '%s': userID cannot be empty)", resourceID, resourceID)
	}
	if resourceID == "forbidden_resource" {
		return errors.Errorf("用户 '%s' 无权访问资源 '%s' (User '%s' is not authorized to access resource '%s')", userID, resourceID, userID, resourceID)
	}
	fmt.Printf("操作 B：用户 '%s' 已成功处理资源 '%s' (Operation B: User '%s' processed resource '%s' successfully)\n", userID, resourceID, userID, resourceID)
	return nil
}

func main() {
	fmt.Println("--- 演示 errors.New --- (--- Demonstrating errors.New ---)")
	errA := simulateOperationA(false) // 触发错误 (Trigger an error)
	if errA != nil {
		fmt.Printf("接收到的错误 (Received error) (%%v): %v\n", errA)
		fmt.Printf("接收到的错误 (Received error) (%%+v):\n%+v\n", errA) // %+v 将打印堆栈跟踪 (%+v will print the stack trace)
	}

	fmt.Println("\n--- 演示 errors.Errorf --- (--- Demonstrating errors.Errorf ---)")
	errB1 := simulateOperationB("", "data_file.txt") // 用户ID为空，触发错误 (Empty userID, triggers an error)
	if errB1 != nil {
		fmt.Printf("接收到的错误1 (Received error 1) (%%v): %v\n", errB1)
		fmt.Printf("接收到的错误1 (Received error 1) (%%+v):\n%+v\n", errB1)
	}

	fmt.Println("")
	errB2 := simulateOperationB("user123", "forbidden_resource") // 访问禁止的资源，触发错误 (Accessing forbidden resource, triggers an error)
	if errB2 != nil {
		fmt.Printf("接收到的错误2 (Received error 2) (%%v): %v\n", errB2)
		fmt.Printf("接收到的错误2 (Received error 2) (%%+v):\n%+v\n", errB2)
	}

	fmt.Println("\n--- 成功案例 --- (--- Success Cases ---)")
	_ = simulateOperationA(true)
	_ = simulateOperationB("user456", "allowed_document.pdf")
}

/*
预期输出 (Expected Output - 堆栈跟踪的路径和行号会因您的环境而异 (Stack trace paths and line numbers will vary based on your environment)):

--- 演示 errors.New --- (--- Demonstrating errors.New ---)
接收到的错误 (Received error) (%v): 操作 A 连接目标服务失败 (Operation A failed to connect to target service)
接收到的错误 (Received error) (%+v):
操作 A 连接目标服务失败 (Operation A failed to connect to target service)
main.simulateOperationA
	/path/to/your/file.go:XX
main.main
	/path/to/your/file.go:YY
...

--- 演示 errors.Errorf --- (--- Demonstrating errors.Errorf ---)
接收到的错误1 (Received error 1) (%v): 处理资源 'data_file.txt' 失败：用户ID不能为空 (Failed to process resource 'data_file.txt': userID cannot be empty)
接收到的错误1 (Received error 1) (%+v):
处理资源 'data_file.txt' 失败：用户ID不能为空 (Failed to process resource 'data_file.txt': userID cannot be empty)
main.simulateOperationB
	/path/to/your/file.go:ZZ
main.main
	/path/to/your/file.go:WW
...

接收到的错误2 (Received error 2) (%v): 用户 'user123' 无权访问资源 'forbidden_resource' (User 'user123' is not authorized to access resource 'forbidden_resource')
接收到的错误2 (Received error 2) (%+v):
用户 'user123' 无权访问资源 'forbidden_resource' (User 'user123' is not authorized to access resource 'forbidden_resource')
main.simulateOperationB
	/path/to/your/file.go:QQ
main.main
	/path/to/your/file.go:PP
...

--- 成功案例 --- (--- Success Cases ---)
操作 A 成功完成 (Operation A completed successfully)
操作 B：用户 'user456' 已成功处理资源 'allowed_document.pdf' (Operation B: User 'user456' processed resource 'allowed_document.pdf' successfully)
*/
``` 