<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 3. 使用错误码 (Using Error Codes)

错误码允许对错误进行程序化处理。`pkg/errors` 模块提供了一个 `Coder` 接口以及用于创建带错误码的错误或将错误码附加到现有错误的函数。

(Error codes allow for programmatic handling of errors. The `pkg/errors` module provides a `Coder` interface and functions to create errors with codes or attach codes to existing errors.)

- **`errors.Coder` 接口 (interface)**: 定义了 `Code() int`, `HTTPStatus() int`, `String() string`, 和 `Reference() string` 方法。
  (Defines methods `Code() int`, `HTTPStatus() int`, `String() string`, and `Reference() string`.)

- **`errors.NewCoder(code int, httpStatus int, message string, reference string) Coder`**: 创建一个新的 `Coder`。
  (Creates a new `Coder`.)

- **`errors.NewWithCode(coder Coder, message string) error`**: 使用给定的 `Coder` 和消息创建一个新错误。
  (Creates a new error with the given `Coder` and message.)

- **`errors.ErrorfWithCode(coder Coder, format string, args ...interface{}) error`**: 使用给定的 `Coder` 创建一个新的格式化错误。
  (Creates a new formatted error with the given `Coder`.)

- **`errors.WithCode(err error, coder Coder) error`**: 将一个 `Coder` 附加到一个现有的错误上。如果 `err` 为 `nil`，则返回 `nil`。
  (Attaches a `Coder` to an existing error. If `err` is `nil`, it returns `nil`.)


```go
package main

import (
	"fmt"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	standardErrors "errors" // 用于示例预先存在的错误 (For a sample pre-existing error)
)

// 定义一个自定义错误码 (Define a custom error code)
// 良好的做法是将这些定义为常量或包级变量。
// (It's good practice to define these as constants or package-level variables.)
var ErrCustomServiceUnavailable = errors.NewCoder(
	50001, // 应用程序特定的唯一代码 (Unique application-specific code)
	503,   // 对应的HTTP状态 (Corresponding HTTP status)
	"自定义服务当前不可用 (Custom service is currently unavailable)", // 此代码的默认消息 (Default message for this code)
	"doc/link/to/service_unavailable_zh.md", // 可选的参考链接 (Optional reference link)
)

// simulateFetchResource
// 演示创建一个带有预定义 Coder 的错误。
// (Demonstrates creating an error with a predefined Coder.)
func simulateFetchResource(resourceID string) error {
	if resourceID == "nonexistent" {
		// 使用 errors 包中预定义的 Coder (例如：errors.ErrNotFound)
		// (Use a predefined Coder from the errors package (e.g., errors.ErrNotFound))
		return errors.NewWithCode(errors.ErrNotFound, fmt.Sprintf("未找到ID为 '%s' 的资源 (resource with ID '%s' was not found)", resourceID, resourceID))
	}
	fmt.Printf("资源 '%s' 获取成功。(Resource '%s' fetched successfully.)\n", resourceID, resourceID)
	return nil
}

// simulateExternalServiceCall
// 演示创建一个带有自定义 Coder 的错误。
// (Demonstrates creating an error with a custom Coder.)
func simulateExternalServiceCall(serviceName string) error {
	if serviceName == "unreliable_service" {
		// 使用我们的自定义 Coder (Use our custom Coder)
		return errors.ErrorfWithCode(ErrCustomServiceUnavailable, "由于间歇性问题调用 %s 失败 (failed to call %s due to intermittent issues)", serviceName, serviceName)
	}
	fmt.Printf("外部服务 '%s' 调用成功。(External service '%s' called successfully.)\n", serviceName, serviceName)
	return nil
}

// processData
// 演示将 Coder 附加到现有错误。
// (Demonstrates attaching a Coder to an existing error.)
func processData(data string) error {
	var existingError error
	if data == "corrupted" {
		existingError = standardErrors.New("检测到底层数据损坏 (underlying data corruption detected)") // 一个预先存在的错误 (A pre-existing error)
	} else if data == "forbidden_data" {
		existingError = errors.New("策略禁止访问此数据 (access to this data is forbidden by policy)") // 一个尚无代码的 lmcc 错误 (An lmcc error without a code yet)
	}


	if existingError != nil {
		// 将我们的自定义 Coder 附加到此现有错误。
		// (Attach our custom Coder to this existing error.)
		// 如果 Coder 是根据处理错误时的上下文确定的，这将非常有用。
		// (This is useful if the Coder is determined based on the context where the error is handled.)
		return errors.WithCode(existingError, errors.ErrBadRequest) // 使用预定义的 ErrBadRequest (Using predefined ErrBadRequest)
	}

	fmt.Printf("数据 '%s' 处理成功。(Data '%s' processed successfully.)\n", data, data)
	return nil
}


func main() {
	fmt.Println("--- 示例1: NewWithCode (预定义 Coder) --- (--- Example 1: NewWithCode (predefined Coder) ---)")
	err1 := simulateFetchResource("nonexistent")
	if err1 != nil {
		fmt.Printf("错误 (Error) (%%v): %v\n", err1)
		fmt.Printf("错误 (Error) (%%+v):\n%+v\n", err1)
		if coder := errors.GetCoder(err1); coder != nil {
			fmt.Printf("  错误码 (Code): %d, HTTP状态 (HTTP Status): %d, 消息 (Message): %s, 参考 (Reference): %s\n",
				coder.Code(), coder.HTTPStatus(), coder.String(), coder.Reference())
		}
	}

	fmt.Println("\n--- 示例2: ErrorfWithCode (自定义 Coder) --- (--- Example 2: ErrorfWithCode (custom Coder) ---)")
	err2 := simulateExternalServiceCall("unreliable_service")
	if err2 != nil {
		fmt.Printf("错误 (Error) (%%v): %v\n", err2)
		fmt.Printf("错误 (Error) (%%+v):\n%+v\n", err2)
		if coder := errors.GetCoder(err2); coder != nil {
			fmt.Printf("  错误码 (Code): %d, HTTP状态 (HTTP Status): %d, 消息 (Message): %s, 参考 (Reference): %s\n",
				coder.Code(), coder.HTTPStatus(), coder.String(), coder.Reference())
		}
	}
	
	fmt.Println("\n--- 示例3: WithCode (将Coder附加到现有错误) --- (--- Example 3: WithCode (attaching Coder to existing error) ---)")
	err3 := processData("corrupted")
	if err3 != nil {
		fmt.Printf("错误 (Error) (%%v): %v\n", err3)
		fmt.Printf("错误 (Error) (%%+v):\n%+v\n", err3)
		if coder := errors.GetCoder(err3); coder != nil {
			fmt.Printf("  错误码 (Code): %d, HTTP状态 (HTTP Status): %d, 消息 (Message): %s, 参考 (Reference): %s\n",
				coder.Code(), coder.HTTPStatus(), coder.String(), coder.Reference())
		}
	}

	fmt.Println("\n--- 示例4: WithCode (将Coder附加到 lmcc 错误) --- (--- Example 4: WithCode (attaching Coder to an lmcc error) ---)")
	err4 := processData("forbidden_data")
	if err4 != nil {
		fmt.Printf("错误 (Error) (%%v): %v\n", err4)
		fmt.Printf("错误 (Error) (%%+v):\n%+v\n", err4)
		if coder := errors.GetCoder(err4); coder != nil {
			fmt.Printf("  错误码 (Code): %d, HTTP状态 (HTTP Status): %d, 消息 (Message): %s, 参考 (Reference): %s\n",
				coder.Code(), coder.HTTPStatus(), coder.String(), coder.Reference())
		}
	}
	
	fmt.Println("\n--- 示例5: 无错误 --- (--- Example 5: No error ---)")
	_ = simulateFetchResource("existent_id")
	_ = simulateExternalServiceCall("reliable_service")
	_ = processData("clean_data")

}

/*
预期输出 (堆栈跟踪会变化) (Expected Output (Stack traces will vary)):

--- 示例1: NewWithCode (预定义 Coder) --- (--- Example 1: NewWithCode (predefined Coder) ---)
错误 (Error) (%v): Not found: 未找到ID为 'nonexistent' 的资源 (resource with ID 'nonexistent' was not found)
错误 (Error) (%+v):
未找到ID为 'nonexistent' 的资源 (resource with ID 'nonexistent' was not found)
main.simulateFetchResource
	/path/to/your/file.go:XX
main.main
	/path/to/your/file.go:YY
...
Not found
  错误码 (Code): 100002, HTTP状态 (HTTP Status): 404, 消息 (Message): Not found, 参考 (Reference): errors-spec.md#100002

--- 示例2: ErrorfWithCode (自定义 Coder) --- (--- Example 2: ErrorfWithCode (custom Coder) ---)
错误 (Error) (%v): 自定义服务当前不可用 (Custom service is currently unavailable): 由于间歇性问题调用 unreliable_service 失败 (failed to call unreliable_service due to intermittent issues)
错误 (Error) (%+v):
由于间歇性问题调用 unreliable_service 失败 (failed to call unreliable_service due to intermittent issues)
main.simulateExternalServiceCall
	/path/to/your/file.go:XX
main.main
	/path/to/your/file.go:YY
...
自定义服务当前不可用 (Custom service is currently unavailable)
  错误码 (Code): 50001, HTTP状态 (HTTP Status): 503, 消息 (Message): 自定义服务当前不可用 (Custom service is currently unavailable), 参考 (Reference): doc/link/to/service_unavailable_zh.md

--- 示例3: WithCode (将Coder附加到现有错误) --- (--- Example 3: WithCode (attaching Coder to existing error) ---)
错误 (Error) (%v): Bad request: 检测到底层数据损坏 (underlying data corruption detected)
错误 (Error) (%+v):
检测到底层数据损坏 (underlying data corruption detected)
main.processData
    /path/to/your/file.go:XX
main.main
    /path/to/your/file.go:YY
...
Bad request
  错误码 (Code): 100003, HTTP状态 (HTTP Status): 400, 消息 (Message): Bad request, 参考 (Reference): errors-spec.md#100003

--- 示例4: WithCode (将Coder附加到 lmcc 错误) --- (--- Example 4: WithCode (attaching Coder to an lmcc error) ---)
错误 (Error) (%v): Bad request: 策略禁止访问此数据 (access to this data is forbidden by policy)
错误 (Error) (%+v):
策略禁止访问此数据 (access to this data is forbidden by policy)
main.processData
	/path/to/your/file.go:XX
main.main
	/path/to/your/file.go:YY
...
Bad request
  错误码 (Code): 100003, HTTP状态 (HTTP Status): 400, 消息 (Message): Bad request, 参考 (Reference): errors-spec.md#100003

--- 示例5: 无错误 --- (--- Example 5: No error ---)
资源 'existent_id' 获取成功。(Resource 'existent_id' fetched successfully.)
外部服务 'reliable_service' 调用成功。(External service 'reliable_service' called successfully.)
数据 'clean_data' 处理成功。(Data 'clean_data' processed successfully.)
*/ 