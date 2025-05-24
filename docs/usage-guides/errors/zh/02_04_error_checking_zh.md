<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 4. 错误检查 (Error Checking)

当您捕获到一个错误后，通常需要检查它的具体内容。

(Once you have an error, you'll want to inspect it.)

**`standardErrors.Is(err, target error) bool`**: (来自标准库) 判断 `err` 的错误链中是否有任何错误匹配 `target`。这对于检查哨兵错误 (sentinel errors) 非常有用，例如 `errors.ErrNotFound`。
(Reports whether any error in `err`'s chain matches `target`. This is useful for checking against sentinel errors like `errors.ErrNotFound`.)

**`errors.GetCoder(err error) Coder`**: 遍历错误链并返回找到的第一个 `Coder`。如果未找到 `Coder`，则返回 `nil`。
(Traverses the error chain and returns the first `Coder` found. Returns `nil` if no `Coder` is found.)

**`errors.IsCode(err error, coder Coder) bool`**: 判断 `err` 的错误链中是否有任何错误的 `Coder` 与所提供的 `coder` 的 `Code()` 相匹配。
(Reports whether any error in `err`'s chain has a `Coder` that matches the `Code()` of the provided `coder`.)


```go
package main

import (
	"fmt"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	standardErrors "errors" // 标准库 errors (Standard library errors)
)

// 自定义哨兵错误 (A custom sentinel error)
var ErrCustomOperationFailed = errors.New("自定义操作失败 (custom operation failed)")

// 自定义 Coder (A custom Coder)
var ErrPaymentDeclined = errors.NewCoder(60001, 402, "提供方拒绝付款 (Payment declined by provider)", "")

// functionThatReturnsVariousErrors
// 一个返回多种不同类型错误的函数
// (A function that returns various types of errors)
func functionThatReturnsVariousErrors(condition string) error {
	switch condition {
	case "not_found":
		// 返回一个包装了预定义 Coder (errors.ErrNotFound) 的错误
		// (Returns an error that wraps a predefined Coder (errors.ErrNotFound))
		originalErr := errors.New("数据库查询未返回任何行 (database query returned no rows)")
		return errors.WithCode(originalErr, errors.ErrNotFound)
	case "validation_error":
		// 返回一个带有不同预定义 Coder (errors.ErrValidation) 的错误
		// (Returns an error with a different predefined Coder (errors.ErrValidation))
		return errors.NewWithCode(errors.ErrValidation, "输入字段 'email' 无效 (input field 'email' is invalid)")
	case "custom_sentinel":
		// 返回一个包装后的自定义哨兵错误
		// (Returns a wrapped custom sentinel error)
		return errors.Wrap(ErrCustomOperationFailed, "尝试处理关键数据 (attempting to process critical data)")
	case "payment_declined":
		// 返回一个带有我们自定义 Coder ErrPaymentDeclined 的错误
		// (Returns an error with our custom Coder ErrPaymentDeclined)
		return errors.ErrorfWithCode(ErrPaymentDeclined, "交易ID %s 被拒绝 (transaction ID %s was declined)", "txn_123abc", "txn_123abc")
	case "unclassified_error":
		// 返回一个没有 Coder 或已知哨兵错误的错误
		// (Returns an error without a Coder or known sentinel)
		return errors.New("发生意外问题 (an unexpected issue occurred)")
	default:
		return nil
	}
}

func main() {
	testConditions := []string{"not_found", "validation_error", "custom_sentinel", "payment_declined", "unclassified_error", "success"}

	for _, cond := range testConditions {
		fmt.Printf("--- 测试条件 (Testing condition): %s ---\n", cond)
		err := functionThatReturnsVariousErrors(cond)

		if err == nil {
			fmt.Println("操作成功，无错误。(Operation successful, no error.)")
			continue
		}

		fmt.Printf("收到的错误 (Received error): %v\n", err)

		// 1. 使用 standardErrors.Is 检查哨兵错误
		// (Check using standardErrors.Is for sentinel errors)
		if standardErrors.Is(err, errors.ErrNotFound) { // 检查预定义的 Coder，它也像一个哨兵错误
											// (Checking against a predefined Coder, which also acts like a sentinel)
			fmt.Println("  [Is] 这是一个 'Not Found' 错误 (通过 errors.ErrNotFound 检查)。([Is] This is a 'Not Found' error (checked via errors.ErrNotFound).)")
		}
		if standardErrors.Is(err, ErrCustomOperationFailed) { // 检查我们的自定义哨兵错误
														// (Checking against our custom sentinel)
			fmt.Println("  [Is] 这是我们的 'Custom Operation Failed' 哨兵错误。([Is] This is our 'Custom Operation Failed' sentinel error.)")
		}

		// 2. 使用 errors.GetCoder 提取 Coder
		// (Extract Coder using errors.GetCoder)
		if coder := errors.GetCoder(err); coder != nil {
			fmt.Printf("  [GetCoder] 提取到的 Coder (Extracted Coder): Code=%d, HTTPStatus=%d, Message='%s', Ref='%s'\n",
				coder.Code(), coder.HTTPStatus(), coder.String(), coder.Reference())

			// 3. 使用 errors.IsCode 通过特定的 Coder 实例进行检查 (对于 Coder 而言最精确)
			// (Check by specific Coder instance using errors.IsCode (most precise for Coders))
			if errors.IsCode(err, errors.ErrNotFound) {
				fmt.Println("  [IsCode] 此错误包含 'ErrNotFound' Coder。([IsCode] This error has the 'ErrNotFound' Coder.)")
			}
			if errors.IsCode(err, errors.ErrValidation) {
				fmt.Println("  [IsCode] 此错误包含 'ErrValidation' Coder。([IsCode] This error has the 'ErrValidation' Coder.)")
			}
			if errors.IsCode(err, ErrPaymentDeclined) {
				fmt.Println("  [IsCode] 此错误包含我们自定义的 'ErrPaymentDeclined' Coder。([IsCode] This error has our custom 'ErrPaymentDeclined' Coder.)")
			}
		} else {
			fmt.Println("  [GetCoder] 此错误链中未找到 Coder。([GetCoder] No Coder found in this error chain.)")
		}
		fmt.Println("-----------------------------------\n")
	}
}
/*
示例输出 (Example Output):

--- 测试条件 (Testing condition): not_found ---
收到的错误 (Received error): Not found: 数据库查询未返回任何行 (database query returned no rows)
  [Is] 这是一个 'Not Found' 错误 (通过 errors.ErrNotFound 检查)。([Is] This is a 'Not Found' error (checked via errors.ErrNotFound).)
  [GetCoder] 提取到的 Coder (Extracted Coder): Code=100002, HTTPStatus=404, Message='Not found', Ref='errors-spec.md#100002'
  [IsCode] 此错误包含 'ErrNotFound' Coder。([IsCode] This error has the 'ErrNotFound' Coder.)
-----------------------------------

--- 测试条件 (Testing condition): validation_error ---
收到的错误 (Received error): Validation failed: 输入字段 'email' 无效 (input field 'email' is invalid)
  [GetCoder] 提取到的 Coder (Extracted Coder): Code=100006, HTTPStatus=400, Message='Validation failed', Ref='errors-spec.md#100006'
  [IsCode] 此错误包含 'ErrValidation' Coder。([IsCode] This error has the 'ErrValidation' Coder.)
-----------------------------------

--- 测试条件 (Testing condition): custom_sentinel ---
收到的错误 (Received error): 尝试处理关键数据 (attempting to process critical data): 自定义操作失败 (custom operation failed)
  [Is] 这是我们的 'Custom Operation Failed' 哨兵错误。([Is] This is our 'Custom Operation Failed' sentinel error.)
  [GetCoder] 此错误链中未找到 Coder。([GetCoder] No Coder found in this error chain.)
-----------------------------------

--- 测试条件 (Testing condition): payment_declined ---
收到的错误 (Received error): 提供方拒绝付款 (Payment declined by provider): 交易ID txn_123abc 被拒绝 (transaction ID txn_123abc was declined)
  [GetCoder] 提取到的 Coder (Extracted Coder): Code=60001, HTTPStatus=402, Message='提供方拒绝付款 (Payment declined by provider)', Ref=''
  [IsCode] 此错误包含我们自定义的 'ErrPaymentDeclined' Coder。([IsCode] This error has our custom 'ErrPaymentDeclined' Coder.)
-----------------------------------

--- 测试条件 (Testing condition): unclassified_error ---
收到的错误 (Received error): 发生意外问题 (an unexpected issue occurred)
  [GetCoder] 此错误链中未找到 Coder。([GetCoder] No Coder found in this error chain.)
-----------------------------------

--- 测试条件 (Testing condition): success ---
操作成功，无错误。(Operation successful, no error.)
*/
```

**错误检查的关键点 (Key points for error checking):**
- 使用 `standardErrors.Is` 检查哨兵错误值 (例如 `io.EOF`，或自定义的如 `ErrCustomOperationFailed`)。`pkg/errors` 中预定义的 `Coder` 实例 (例如 `errors.ErrNotFound`) 也可以与 `standardErrors.Is` 一起使用，因为它们实际上就是哨兵值。
  (Use `standardErrors.Is` for checking against sentinel error values (like `io.EOF`, or custom ones like `ErrCustomOperationFailed`). Predefined `Coder` instances from `pkg/errors` (e.g., `errors.ErrNotFound`) can also be used with `standardErrors.Is` because they are, in effect, sentinel values.)
- 如果错误链中存在 `Coder`，使用 `errors.GetCoder` 来检索它。然后您可以检查其属性 (`Code()`, `HTTPStatus()` 等)。
  (Use `errors.GetCoder` to retrieve a `Coder` if one exists in the error chain. You can then inspect its properties (`Code()`, `HTTPStatus()`, etc.).)
- 使用 `errors.IsCode` 来专门检查错误链中的错误是否与特定 `Coder` 的 `Code()` 匹配。这是检查由 `Coder` 定义的错误类别的最直接方法。
  (Use `errors.IsCode` to specifically check if an error in the chain matches a particular `Coder`'s `Code()`. This is the most direct way to check for an error category defined by a `Coder`.)

</rewritten_file> 