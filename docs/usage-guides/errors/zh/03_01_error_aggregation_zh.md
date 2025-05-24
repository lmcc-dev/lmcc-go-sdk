<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 错误聚合 (Error Aggregation)

将单个更广泛操作期间发生的多个错误 (例如，验证表单的多个字段) 收集到单个 `ErrorGroup` 中。

(Collect multiple errors that occur during a single broader operation (e.g., validating multiple fields of a form) into a single `ErrorGroup`.)

**`errors.NewErrorGroup(message string) *ErrorGroup`**: 创建一个带有初始消息的新错误组。

(Creates a new error group with an initial message.)

`ErrorGroup` 具有以下方法：

(An `ErrorGroup` has the following methods:)

- **`Add(err error)`**: 将错误添加到组中。如果 `err` 为 `nil`，则不执行任何操作。
  (Adds an error to the group. Does nothing if `err` is nil.)
- **`Errors() []error`**: 返回添加到组中的所有错误的切片。如果未添加任何错误，则返回 `nil`。
  (Returns a slice of all errors added to the group. Returns `nil` if no errors were added.)
- **`Error() string`**: 返回组中所有错误的字符串表示形式，以组的初始消息为前缀。
  (Returns a string representation of all errors in the group, prefixed by the group's initial message.)
- **`Unwrap() []error`**: (Go 1.20+) 通过公开收集到的错误，允许 `ErrorGroup` 与 `standardErrors.Is` 和 `standardErrors.As` 一起使用。
  (Allows `ErrorGroup` to be used with `standardErrors.Is` and `standardErrors.As` by exposing the collected errors.)
- **`Format(s fmt.State, verb rune)`**: 实现 `fmt.Formatter` 以进行自定义格式化 (例如，`%+v` 用于详细输出，包括单个错误的堆栈跟踪)。
  (Implements `fmt.Formatter` for custom formatting (e.g., `%+v` for detailed output including individual error stack traces).)

```go
package main

import (
	"fmt"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"strings"
)

// UserProfile
// 我们想要验证的示例结构体。
// (A sample struct we want to validate.)
type UserProfile struct {
	Username string
	Email    string
	Age      int
}

// validateUserProfile
// 此函数验证 UserProfile 并将所有验证错误收集到 ErrorGroup中。
// (This function validates a UserProfile and collects all validation errors into an ErrorGroup.)
func validateUserProfile(profile UserProfile) error {
	// 使用此验证操作的通用消息创建一个新的错误组。
	// (Create a new error group with a general message for this validation operation.)
	eg := errors.NewErrorGroup("用户配置文件验证失败 (User profile validation failed)")

	// 验证用户名 (Validate Username)
	if strings.TrimSpace(profile.Username) == "" {
		// 为缺少的用户名添加错误。我们可以使用 errors.New 或 errors.NewWithCode。
		// (Add an error for missing username. We can use errors.New or errors.NewWithCode.)
		eg.Add(errors.New("用户名是必需的 (username is required)"))
	}

	// 验证电子邮件 (Validate Email)
	if strings.TrimSpace(profile.Email) == "" {
		eg.Add(errors.NewWithCode(errors.ErrValidation, "电子邮件不能为空 (email cannot be empty)"))
	} else if !strings.Contains(profile.Email, "@") {
		eg.Add(errors.ErrorfWithCode(errors.ErrValidation, "电子邮件 '%s' 不是有效格式 (email '%s' is not a valid format)", profile.Email))
	}

	// 验证年龄 (Validate Age)
	if profile.Age < 0 {
		eg.Add(errors.New("年龄不能为负 (age cannot be negative)"))
	} else if profile.Age < 18 {
		eg.Add(errors.Errorf("用户必须年满18岁，实际年龄 %d (user must be at least 18 years old, got %d)", profile.Age))
	}

	// 检查是否有任何错误已添加到组中。
	// (Check if any errors were added to the group.)
	if len(eg.Errors()) > 0 {
		return eg // 返回 ErrorGroup 本身，它是一个错误。(Return the ErrorGroup itself, which is an error.)
	}

	// 如果没有添加错误，则验证通过。
	// (If no errors were added, validation passed.)
	return nil
}

func main() {
	// 场景1：具有多个验证问题的配置文件
	// (Scenario 1: Profile with multiple validation issues)
	profile1 := UserProfile{Username: "", Email: "invalid-email", Age: 15}
	fmt.Println("--- 验证配置文件1 (多个问题) --- (--- Validating Profile 1 (multiple issues) ---)")
	err1 := validateUserProfile(profile1)
	if err1 != nil {
		fmt.Printf("验证错误 (Validation Error) (%%v):\n%v\n\n", err1)
		fmt.Printf("验证错误 (Validation Error) (%%+v) - 详细信息 (Detailed):\n%+v\n", err1)
		
		// 如果需要，您可以遍历单个错误
		// (You can iterate over individual errors if needed)
		if eg, ok := err1.(*errors.ErrorGroup); ok {
			fmt.Println("\n组中的单个错误 (Individual errors in the group):")
			for i, individualErr := range eg.Errors() {
				fmt.Printf("  %d: %v\n", i+1, individualErr)
			}
			
			// 检查 ErrorGroup 是否包含特定代码的错误
			// (Check if the ErrorGroup contains errors with specific codes)
			fmt.Println("\n检查组中的特定错误代码 (Checking for specific error codes in the group):")
			if errors.IsCode(eg, errors.ErrValidation) {
				fmt.Println("  - 包含验证错误 (Contains validation errors)")
			}
			if errors.IsCode(eg, errors.ErrNotFound) {
				fmt.Println("  - 包含未找到错误 (Contains not found errors)")
			} else {
				fmt.Println("  - 没有未找到错误 (No not found errors)")
			}
		}
	}

	fmt.Println("\n--- 验证配置文件2 (一个问题) --- (--- Validating Profile 2 (one issue) ---)")
	// 场景2：具有单个验证问题的配置文件
	// (Scenario 2: Profile with a single validation issue)
	profile2 := UserProfile{Username: "ValidUser", Email: "valid@example.com", Age: -5}
	err2 := validateUserProfile(profile2)
	if err2 != nil {
		fmt.Printf("验证错误 (Validation Error) (%%v):\n%v\n\n", err2)
		// fmt.Printf("验证错误 (Validation Error) (%%+v):\n%+v\n", err2) // %%+v 也会显示每个错误的堆栈跟踪 (%+v would also show stack traces for each error)
	}

	fmt.Println("\n--- 验证配置文件3 (有效配置文件) --- (--- Validating Profile 3 (valid profile) ---)")
	// 场景3：有效的配置文件，无错误
	// (Scenario 3: Valid profile, no errors)
	profile3 := UserProfile{Username: "TestUser", Email: "test@example.com", Age: 30}
	err3 := validateUserProfile(profile3)
	if err3 == nil {
		fmt.Println("配置文件3有效！(Profile 3 is valid!)")
	} else {
		fmt.Printf("配置文件3出现意外的验证错误 (Unexpected validation error for Profile 3): %v\n", err3)
	}
}

/*
示例输出 (为简洁起见，此处省略了 %%+v 中的堆栈跟踪，但实际会存在)：
(Example Output (Stack traces in %+v are omitted for brevity here but would be present)):

--- 验证配置文件1 (多个问题) --- (--- Validating Profile 1 (multiple issues) ---)
验证错误 (Validation Error) (%v):
用户配置文件验证失败 (User profile validation failed): [用户名是必需的 (username is required); 电子邮件 'invalid-email' 不是有效格式 (email 'invalid-email' is not a valid format); 用户必须年满18岁，实际年龄 15 (user must be at least 18 years old, got 15)]

验证错误 (Validation Error) (%+v) - 详细信息 (Detailed):
用户配置文件验证失败 (User profile validation failed): [用户名是必需的 (username is required); 电子邮件 'invalid-email' 不是有效格式 (email 'invalid-email' is not a valid format); 用户必须年满18岁，实际年龄 15 (user must be at least 18 years old, got 15)]
带有堆栈跟踪的单个错误 (Individual errors with stack traces):
1. 用户名是必需的 (username is required)
   main.validateUserProfile
   	/path/to/your/file.go:XX
   ...
2. Validation failed: 电子邮件 'invalid-email' 不是有效格式 (email 'invalid-email' is not a valid format)
   main.validateUserProfile
   	/path/to/your/file.go:XX
   ...
3. 用户必须年满18岁，实际年龄 15 (user must be at least 18 years old, got 15)
   main.validateUserProfile
   	/path/to/your/file.go:XX
   ...

组中的单个错误 (Individual errors in the group):
  1: 用户名是必需的 (username is required)
  2: Validation failed: 电子邮件 'invalid-email' 不是有效格式 (email 'invalid-email' is not a valid format)
  3: 用户必须年满18岁，实际年龄 15 (user must be at least 18 years old, got 15)

--- 验证配置文件2 (一个问题) --- (--- Validating Profile 2 (one issue) ---)
验证错误 (Validation Error) (%v):
用户配置文件验证失败 (User profile validation failed): [年龄不能为负 (age cannot be negative)]

--- 验证配置文件3 (有效配置文件) --- (--- Validating Profile 3 (valid profile) ---)
配置文件3有效！(Profile 3 is valid!)
*/
``` 