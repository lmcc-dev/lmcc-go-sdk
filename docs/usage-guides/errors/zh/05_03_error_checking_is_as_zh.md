<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 3. 使用 `errors.Is` 和 `errors.As` 进行错误检查 (Error Checking with `errors.Is` and `errors.As`)

使用 `standardErrors.Is(err, target)` 和 `standardErrors.As(err, &targetType)` 的代码将继续按预期与由 `pkg/errors` 创建或包装的错误一起工作。

(Code using `standardErrors.Is(err, target)` and `standardErrors.As(err, &targetType)` will continue to work as expected with errors created or wrapped by `pkg/errors`.)

- **`standardErrors.Is`**: 检查链中的任何错误是否与目标哨兵错误匹配。适用于 `pkg/errors` 类似哨兵的 `Coder` 变量 (例如，`pkgErrors.ErrNotFound`)。
  (`standardErrors.Is`: Checks if any error in the chain matches a target sentinel error. Works with `pkg/errors` sentinel-like `Coder` variables (e.g., `pkgErrors.ErrNotFound`).)
- **`standardErrors.As`**: 检查链中的任何错误是否可以转换为特定类型，并将目标设置为该错误。适用于可能是链一部分的自定义错误类型。
  (`standardErrors.As`: Checks if any error in the chain can be cast to a specific type and sets the target to that error. Useful for custom error types that might be part of the chain.)
- **`pkgErrors.IsCode`**: 如果您不需要检查特定的 `Coder` 实例，请使用此方法检查 `Coder` 的数字代码 (请参阅最佳实践)。
  (`pkgErrors.IsCode`: Use this for checking against the numeric code of a `Coder` if you don\'t need to check for a specific `Coder` instance (see Best Practices).)

```go
package main

import (
	"fmt"
	standardErrors "errors" // Go 标准库 errors (Go standard library errors)
	pkgErrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors" // 我们的 pkg/errors (Our pkg/errors)
	"os"
)

// 可能使用的自定义错误类型
// (A custom error type that might be used)
type MyCustomErrorType struct {
	Msg     string
	Details string
}

func (e *MyCustomErrorType) Error() string {
	return fmt.Sprintf("%s (详情：%s) (%s (Details: %s))", e.Msg, e.Details, e.Msg, e.Details)
}

// 可以返回各种类型错误的函数
// (Function that can return various types of errors)
func generateError(errorType string) error {
	switch errorType {
	case "std_is_target":
		// 我们可能使用 errors.Is 检查的标准库错误
		// (A standard library error that we might check with errors.Is)
		return os.ErrPermission // 预定义的标准错误 (Predefined standard error)
	case "pkg_is_target":
		// 一个 pkg/errors Coder，充当 errors.Is 的哨兵
		// (A pkg/errors Coder that acts as a sentinel for errors.Is)
		return pkgErrors.Wrap(pkgErrors.ErrUnauthorized, "用户缺乏此 pkg 操作的权限 (user lacks permissions for this pkg operation)")
	case "custom_as_target":
		// 我们自定义错误类型的实例，可能已包装
		// (An instance of our custom error type, possibly wrapped)
		customErr := &MyCustomErrorType{Msg: "自定义操作失败 (Custom operation failed)", Details: "数据库约束冲突 (Database constraint violation)"}
		return pkgErrors.Wrap(customErr, "包装自定义错误类型 (wrapping custom error type)")
	case "pkg_with_code":
		// 带有特定 Coder 的 pkg/errors 错误，适用于 IsCode
		// (A pkg/errors error with a specific Coder, good for IsCode)
		return pkgErrors.NewWithCode(pkgErrors.ErrConflict, "资源版本不匹配 (resource version mismatch)")
	default:
		return pkgErrors.New("发生未知错误 (an unknown error occurred)")
	}
}

func main() {
	testCases := []string{"std_is_target", "pkg_is_target", "custom_as_target", "pkg_with_code", "unknown"}

	for _, tc := range testCases {
		fmt.Printf("--- 测试错误类型：%s ---\n", tc)
		// (--- Testing error type: %s ---\n)
		err := generateError(tc)

		if err == nil {
			fmt.Println("未生成错误。(No error generated.)")
			continue
		}
		fmt.Printf("生成的错误：%v\n", err)
		// (Generated error: %v\n)

		// 1. 使用 standardErrors.Is (Using standardErrors.Is)
		if standardErrors.Is(err, os.ErrPermission) {
			fmt.Println("  [errors.Is] 匹配 os.ErrPermission。(Matched os.ErrPermission.)")
		}
		if standardErrors.Is(err, pkgErrors.ErrUnauthorized) { // pkgErrors.ErrUnauthorized 是一个 Coder，也充当哨兵 (pkgErrors.ErrUnauthorized is a Coder, also acts as a sentinel)
			fmt.Println("  [errors.Is] 匹配 pkgErrors.ErrUnauthorized。(Matched pkgErrors.ErrUnauthorized.)")
		}

		// 2. 使用 standardErrors.As (Using standardErrors.As)
		var customErrTarget *MyCustomErrorType
		if standardErrors.As(err, &customErrTarget) {
			fmt.Printf("  [errors.As] 匹配 MyCustomErrorType。消息：'%s'，详情：'%s'\n", customErrTarget.Msg, customErrTarget.Details)
			// (Matched MyCustomErrorType. Message: '%s', Details: '%s'\n)
		} else {
			// fmt.Println("  [errors.As] 未匹配 MyCustomErrorType。(Did not match MyCustomErrorType.)")
		}

		// 3. 使用 pkgErrors.IsCode (适用于 pkg/errors Coders) (Using pkgErrors.IsCode (for pkg/errors Coders))
		if pkgErrors.IsCode(err, pkgErrors.ErrUnauthorized) {
			fmt.Println("  [pkgErrors.IsCode] 匹配 pkgErrors.ErrUnauthorized 的代码。(Matched code for pkgErrors.ErrUnauthorized.)")
		}
		if pkgErrors.IsCode(err, pkgErrors.ErrConflict) {
			fmt.Println("  [pkgErrors.IsCode] 匹配 pkgErrors.ErrConflict 的代码。(Matched code for pkgErrors.ErrConflict.)")
			coder := pkgErrors.GetCoder(err)
			if coder != nil {
				fmt.Printf("    Coder 详情：Code=%d, HTTPStatus=%d, Message='%s'\n", coder.Code(), coder.HTTPStatus(), coder.String())
				// (Coder details: Code=%d, HTTPStatus=%d, Message='%s'\n)
			}
		}
		fmt.Println("") // 换行以便清晰 (Newline for clarity)
	}
}

/*
示例输出 (Example Output):

--- 测试错误类型：std_is_target ---
(--- Testing error type: std_is_target ---)
生成的错误：permission denied
(Generated error: permission denied)
  [errors.Is] 匹配 os.ErrPermission。(Matched os.ErrPermission.)

--- 测试错误类型：pkg_is_target ---
(--- Testing error type: pkg_is_target ---)
生成的错误：用户缺乏此 pkg 操作的权限 (user lacks permissions for this pkg operation): Unauthorized
(Generated error: user lacks permissions for this pkg operation: Unauthorized)
  [errors.Is] 匹配 pkgErrors.ErrUnauthorized。(Matched pkgErrors.ErrUnauthorized.)
  [pkgErrors.IsCode] 匹配 pkgErrors.ErrUnauthorized 的代码。(Matched code for pkgErrors.ErrUnauthorized.)

--- 测试错误类型：custom_as_target ---
(--- Testing error type: custom_as_target ---)
生成的错误：包装自定义错误类型 (wrapping custom error type): 自定义操作失败 (详情：数据库约束冲突) (Custom operation failed (Details: Database constraint violation))
(Generated error: wrapping custom error type: Custom operation failed (Details: Database constraint violation))
  [errors.As] 匹配 MyCustomErrorType。消息：'自定义操作失败 (Custom operation failed)'，详情：'数据库约束冲突 (Database constraint violation)'
  (Matched MyCustomErrorType. Message: 'Custom operation failed', Details: 'Database constraint violation')

--- 测试错误类型：pkg_with_code ---
(--- Testing error type: pkg_with_code ---)
生成的错误：Conflict: 资源版本不匹配 (resource version mismatch)
(Generated error: Conflict: resource version mismatch)
  [pkgErrors.IsCode] 匹配 pkgErrors.ErrConflict 的代码。(Matched code for pkgErrors.ErrConflict.)
    Coder 详情：Code=100005, HTTPStatus=409, Message='Conflict'
    (Coder details: Code=100005, HTTPStatus=409, Message='Conflict')

--- 测试错误类型：unknown ---
(--- Testing error type: unknown ---)
生成的错误：发生未知错误 (an unknown error occurred)
(Generated error: an unknown error occurred)

*/
``` 