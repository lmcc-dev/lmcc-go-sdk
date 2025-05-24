<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 4. 使用 `errors.IsCode` 检查特定的错误类别 (Use `errors.IsCode` for Checking Specific Error Categories)

当您需要检查错误是否属于由 `Coder` 定义的类别时 (例如，任何类型的"未找到"错误，无论具体消息如何)，请使用 `errors.IsCode(err, YourCoder)`。

(When you need to check if an error belongs to a category defined by a `Coder` (e.g., any kind of "not found" error, regardless of the specific message), use `errors.IsCode(err, YourCoder)`.)

- `errors.IsCode` 检查 `Coder` 实例的 `Code()`，而不是 `Coder` 变量的身份。这意味着如果您有多个碰巧共享相同整数代码的 `Coder` 变量，`IsCode` 会认为它们对于该代码是匹配的。
  ( `errors.IsCode` checks the `Code()` of the `Coder` instances, not the `Coder` variable identity. This means if you have multiple `Coder` variables that happen to share the same integer code, `IsCode` would consider them matching for that code.)
- 如果 `YourCoderVariable` 可能被包装，或者您只想根据数字代码进行检查，这通常比 `standardErrors.Is(err, YourCoderVariable)` 更可靠。
  (This is generally more reliable than `standardErrors.Is(err, YourCoderVariable)` if `YourCoderVariable` might be wrapped or if you want to check solely based on the numeric code.)
- 然而，标准的 `errors.Is(err, errors.ErrNotFound)` (其中 `errors.ErrNotFound` 是 `pkg/errors` 本身预定义的 `Coder`) 工作得很好，因为这些预定义的 `Coder` 变量就像哨兵值一样。
  (However, standard `errors.Is(err, errors.ErrNotFound)` (where `errors.ErrNotFound` is a predefined `Coder` from the `pkg/errors` itself) works well because these predefined `Coder` variables act like sentinel values.)

```go
package main

import (
	"fmt"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	standardLibErrors "errors" // 标准库错误的别名 (Alias for standard library errors)
)

// 预定义的 Coder (假设这些来自大型系统的不同部分)
// (Predefined Coders (imagine these are from different parts of a larger system))
var ErrLegacyNotFound = errors.NewCoder(40400, 404, "资源未找到 (遗留系统) (Resource not found (legacy system))", "")
var ErrAPIV1NotFound = errors.NewCoder(40400, 404, "端点未找到 (API v1) (Endpoint not found (API v1))", "") // 与 ErrLegacyNotFound 代码相同 (Same code as ErrLegacyNotFound)
var ErrResourceGone = errors.NewCoder(41000, 410, "资源已永久删除 (Resource permanently gone)", "")

// 可能返回这些错误之一的函数
// (Function that might return one of these errors)
func fetchResource(resourceID string) error {
	switch resourceID {
	case "legacy_item_123":
		return errors.WrapWithCode(standardLibErrors.New("原始数据库错误：无行 (original db error: no rows)"), ErrLegacyNotFound, "访问遗留数据存储 (accessing legacy data store)")
	case "api_v1_user_789":
		// 直接返回带有 Coder 的错误
		// (Directly return an error with the Coder)
		return errors.NewWithCode(ErrAPIV1NotFound, "用户端点 /users/789 不可用 (user endpoint /users/789 not available)")
	case "deleted_doc_456":
		return errors.NewWithCode(ErrResourceGone, "文档 456 已被删除 (document 456 was deleted)")
	case "other_issue":
		return errors.New("其他一些不相关的问题 (some other unrelated problem)") // 没有 Coder 的错误 (Error without a Coder)
	default:
		fmt.Printf("资源 '%s' 已成功找到。(Resource '%s' found successfully.)\n", resourceID, resourceID)
		return nil
	}
}

func main() {
	testCases := []struct {
		name           string
		resourceID     string
		expectIsCodeMatchWithLegacyNotFound bool // errors.IsCode 检查代码编号 (40400)
		                                        // (errors.IsCode checks the code number (40400))
		expectIsMatchWithLegacyNotFound   bool // errors.Is 检查链中特定的 Coder 实例
		                                        // (errors.Is checks for the specific Coder instance in the chain)
		expectIsMatchWithAPIV1NotFound bool // errors.Is 不会匹配 ErrAPIV1NotFound 实例
		                                        // (errors.Is would not match ErrAPIV1NotFound instance)
	}{
		{
			name: "带有 ErrLegacyNotFound 的错误 (Error with ErrLegacyNotFound)",
			resourceID: "legacy_item_123",
			expectIsCodeMatchWithLegacyNotFound: true, 
			expectIsMatchWithLegacyNotFound:   true, 
			expectIsMatchWithAPIV1NotFound: false, 
		},
		{
			name: "带有 ErrAPIV1NotFound 的错误 (Error with ErrAPIV1NotFound)",
			resourceID: "api_v1_user_789",
			expectIsCodeMatchWithLegacyNotFound: true, // errors.IsCode 匹配，因为 ErrAPIV1NotFound 也有代码 40400
			                                        // (errors.IsCode matches because ErrAPIV1NotFound also has code 40400)
			expectIsMatchWithLegacyNotFound:   false, 
			expectIsMatchWithAPIV1NotFound: true,
		},
		{
			name: "带有 ErrResourceGone 的错误 (Error with ErrResourceGone)",
			resourceID: "deleted_doc_456",
			expectIsCodeMatchWithLegacyNotFound: false, // 不同的代码 (41000) (Different code (41000))
			expectIsMatchWithLegacyNotFound:   false,
			expectIsMatchWithAPIV1NotFound: false,
		},
		{
			name: "未编码的错误 (Uncoded error)",
			resourceID: "other_issue",
			expectIsCodeMatchWithLegacyNotFound: false,
			expectIsMatchWithLegacyNotFound:   false,
			expectIsMatchWithAPIV1NotFound: false,
		},
		{
			name: "没有错误 (No error)",
			resourceID: "existing_resource",
		},
	}

	for _, tc := range testCases {
		fmt.Printf("--- 测试用例：%s (资源 ID：%s) ---\n", tc.name, tc.resourceID)
		// (--- Test Case: %s (Resource ID: %s) ---\n)
		err := fetchResource(tc.resourceID)

		if err == nil {
			fmt.Println("没有发生错误。(No error occurred.)")
			continue
		}
		fmt.Printf("收到的错误：%v\n", err)
		// (Received error: %v\n)

		// 使用 errors.IsCode (基于 Coder 的 Code() 值进行检查)
		// (Using errors.IsCode (checks based on the Coder's Code() value))
		isCodeLegacy := errors.IsCode(err, ErrLegacyNotFound)
		fmt.Printf("  errors.IsCode(err, ErrLegacyNotFound (代码 %d)): %t。预期：%t\n", ErrLegacyNotFound.Code(), isCodeLegacy, tc.expectIsCodeMatchWithLegacyNotFound)
		// (errors.IsCode(err, ErrLegacyNotFound (code %d)): %t. Expected: %t\n)

		// 使用 standardErrors.Is (检查链中特定的 Coder 实例)
		// (Using standardErrors.Is (checks for specific Coder instance in the chain))
		isLegacy := standardLibErrors.Is(err, ErrLegacyNotFound)
		fmt.Printf("  standardErrors.Is(err, ErrLegacyNotFound): %t。预期：%t\n", isLegacy, tc.expectIsMatchWithLegacyNotFound)
		// (standardErrors.Is(err, ErrLegacyNotFound): %t. Expected: %t\n)

		isAPIV1 := standardLibErrors.Is(err, ErrAPIV1NotFound)
		fmt.Printf("  standardErrors.Is(err, ErrAPIV1NotFound): %t。预期：%t\n", isAPIV1, tc.expectIsMatchWithAPIV1NotFound)
		// (standardErrors.Is(err, ErrAPIV1NotFound): %t. Expected: %t\n)

		coder := errors.GetCoder(err)
		if coder != nil {
			fmt.Printf("  错误中的实际 Coder：Code=%d，Message='%s'\n", coder.Code(), coder.String())
			// (Actual Coder in error: Code=%d, Message='%s'\n)
		} else {
			fmt.Println("  错误中未找到 Coder。(No Coder found in error.)")
		}
	}
}

/*
示例输出 (Example Output):

--- 测试用例：带有 ErrLegacyNotFound 的错误 (资源 ID：legacy_item_123) ---
(--- Test Case: Error with ErrLegacyNotFound (Resource ID: legacy_item_123) ---)
收到的错误：访问遗留数据存储 (accessing legacy data store): 资源未找到 (遗留系统) (Resource not found (legacy system)): 原始数据库错误：无行 (original db error: no rows)
(Received error: accessing legacy data store: Resource not found (legacy system): original db error: no rows)
  errors.IsCode(err, ErrLegacyNotFound (代码 40400)): true。预期：true
  (errors.IsCode(err, ErrLegacyNotFound (code 40400)): true. Expected: true)
  standardErrors.Is(err, ErrLegacyNotFound): true。预期：true
  (standardErrors.Is(err, ErrLegacyNotFound): true. Expected: true)
  standardErrors.Is(err, ErrAPIV1NotFound): false。预期：false
  (standardErrors.Is(err, ErrAPIV1NotFound): false. Expected: false)
  错误中的实际 Coder：Code=40400，Message='资源未找到 (遗留系统) (Resource not found (legacy system))'
  (Actual Coder in error: Code=40400, Message='Resource not found (legacy system)')

--- 测试用例：带有 ErrAPIV1NotFound 的错误 (资源 ID：api_v1_user_789) ---
(--- Test Case: Error with ErrAPIV1NotFound (Resource ID: api_v1_user_789) ---)
收到的错误：端点未找到 (API v1) (Endpoint not found (API v1)): 用户端点 /users/789 不可用 (user endpoint /users/789 not available)
(Received error: Endpoint not found (API v1): user endpoint /users/789 not available)
  errors.IsCode(err, ErrLegacyNotFound (代码 40400)): true。预期：true
  (errors.IsCode(err, ErrLegacyNotFound (code 40400)): true. Expected: true)
  standardErrors.Is(err, ErrLegacyNotFound): false。预期：false
  (standardErrors.Is(err, ErrLegacyNotFound): false. Expected: false)
  standardErrors.Is(err, ErrAPIV1NotFound): true。预期：true
  (standardErrors.Is(err, ErrAPIV1NotFound): true. Expected: true)
  错误中的实际 Coder：Code=40400，Message='端点未找到 (API v1) (Endpoint not found (API v1))'
  (Actual Coder in error: Code=40400, Message='Endpoint not found (API v1)')

--- 测试用例：带有 ErrResourceGone 的错误 (资源 ID：deleted_doc_456) ---
(--- Test Case: Error with ErrResourceGone (Resource ID: deleted_doc_456) ---)
收到的错误：资源已永久删除 (Resource permanently gone): 文档 456 已被删除 (document 456 was deleted)
(Received error: Resource permanently gone: document 456 was deleted)
  errors.IsCode(err, ErrLegacyNotFound (代码 40400)): false。预期：false
  (errors.IsCode(err, ErrLegacyNotFound (code 40400)): false. Expected: false)
  standardErrors.Is(err, ErrLegacyNotFound): false。预期：false
  (standardErrors.Is(err, ErrLegacyNotFound): false. Expected: false)
  standardErrors.Is(err, ErrAPIV1NotFound): false。预期：false
  (standardErrors.Is(err, ErrAPIV1NotFound): false. Expected: false)
  错误中的实际 Coder：Code=41000，Message='资源已永久删除 (Resource permanently gone)'
  (Actual Coder in error: Code=41000, Message='Resource permanently gone')

--- 测试用例：未编码的错误 (资源 ID：other_issue) ---
(--- Test Case: Uncoded error (Resource ID: other_issue) ---)
收到的错误：其他一些不相关的问题 (some other unrelated problem)
(Received error: some other unrelated problem)
  errors.IsCode(err, ErrLegacyNotFound (代码 40400)): false。预期：false
  (errors.IsCode(err, ErrLegacyNotFound (code 40400)): false. Expected: false)
  standardErrors.Is(err, ErrLegacyNotFound): false。预期：false
  (standardErrors.Is(err, ErrLegacyNotFound): false. Expected: false)
  standardErrors.Is(err, ErrAPIV1NotFound): false。预期：false
  (standardErrors.Is(err, ErrAPIV1NotFound): false. Expected: false)
  错误中未找到 Coder。(No Coder found in error.)

--- 测试用例：没有错误 (资源 ID：existing_resource) ---
(--- Test Case: No error (Resource ID: existing_resource) ---)
资源 'existing_resource' 已成功找到。(Resource 'existing_resource' found successfully.)
没有发生错误。(No error occurred.)
*/
```

**何时使用 `errors.IsCode` 与 `standardErrors.Is`：(When to use `errors.IsCode` vs `standardErrors.Is`)**

- 当您想检查错误是否属于由 `SpecificCoder.Code()` 标识的某个*类别*，而不管确切的 `Coder` 实例时，请使用 `errors.IsCode(err, SpecificCoder)`。
  (Use `errors.IsCode(err, SpecificCoder)` when you want to check if the error is of a certain *category* identified by the `SpecificCoder.Code()`, regardless of the exact `Coder` instance.)
- 当检查错误链中特定的哨兵错误值或特定的 `Coder` 变量实例时，请使用 `standardErrors.Is(err, TargetErrorOrCoderInstance)`。这对于 `pkg/errors` 中预定义的 `Coder` 变量 (如 `errors.ErrNotFound`、`errors.ErrInternalServer` 等，它们充当哨兵) 非常有用。如果您想检查那个*确切的*实例，它也适用于您自己的包级 `Coder` 变量。
  (Use `standardErrors.Is(err, TargetErrorOrCoderInstance)` when checking for a specific sentinel error value or a specific `Coder` variable instance in the error chain. This is useful for predefined `Coder` variables from `pkg/errors` like `errors.ErrNotFound`, `errors.ErrInternalServer`, etc., which act as sentinels. It's also appropriate for your own package-level `Coder` variables if you want to check for that *exact* instance.)