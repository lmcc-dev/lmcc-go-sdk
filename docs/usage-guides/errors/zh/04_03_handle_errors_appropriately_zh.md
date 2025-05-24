<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 3. 在不同层级适当地处理错误 (Handle Errors Appropriately at Different Layers)

并非所有错误都应以相同的方式处理。所采取的行动应取决于应用程序的层级和错误的性质。

(Not all errors should be handled in the same way. The action taken should depend on the layer of the application and the nature of the error.)

- **底层函数 (Low-Level Functions)**: 通常，底层函数 (例如，数据库访问、文件 I/O) 应简单地将错误 (可能已包装) 返回给其调用者。它们缺乏做出高级决策 (例如，重试、中止、响应用户) 的上下文。
  (Often, low-level functions (e.g., database access, file I/O) should simply return errors (possibly wrapped) to their callers. They lack the context to make high-level decisions (e.g., retry, abort, respond to user).)
- **服务/业务逻辑层 (Service/Business Logic Layer)**: 此层可能会通过重试操作、回退到替代策略或在返回错误之前用特定于业务的 `Coder` 实例包装错误来处理某些错误。
  (This layer might handle certain errors by retrying operations, falling back to alternative strategies, or by wrapping errors with business-specific `Coder` instances before returning them.)
- **顶层处理程序 (例如，HTTP 处理程序、主函数) (Top-Level Handlers (e.g., HTTP handlers, main function))**: 这通常是记录错误并为用户或客户端系统制定响应的地方。使用 `errors.GetCoder` 提取代码以生成适当的 HTTP 状态码或用户消息。
  (This is typically where errors are logged, and responses are formulated for the user or client system. Use `errors.GetCoder` to extract codes for generating appropriate HTTP status codes or user messages.)

```go
package main

import (
	"fmt"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"time"
)

// --- 第1层：底层数据访问 --- (--- Layer 1: Low-Level Data Access ---)

var ErrDBConnectionFailed = errors.NewCoder(30001, 500, "数据库连接失败 (Database connection failed)", "")
var ErrDBRecordNotFound = errors.NewCoder(30002, 404, "数据库记录未找到 (Database record not found)", "")

// fetchUserDataFromDB 模拟从数据库获取用户数据。
// (fetchUserDataFromDB simulates fetching user data from a database.)
func fetchUserDataFromDB(userID string) (string, error) {
	fmt.Printf("[数据库访问] 尝试获取用户 '%s'\n", userID)
	// ([DB Access] Attempting to fetch user '%s'\n)
	if userID == "user_transient_error" {
		// 模拟临时数据库问题 (Simulate a temporary DB issue)
		return "", errors.WrapWithCode(errors.New("连接数据库副本时网络超时 (network timeout connecting to db replica)"), ErrDBConnectionFailed, "fetchUserDataFromDB")
	}
	if userID == "user_not_found" {
		return "", errors.ErrorfWithCode(ErrDBRecordNotFound, "没有用户 ID '%s' 的记录 (no record for user ID '%s')", userID)
	}
	if userID == "user_critical_db_issue" {
		// 最初来自数据库层的更严重的、未编码的错误
		// (A more severe, non-coded error from DB layer initially)
		return "", errors.New("用户表上不可恢复的校验和错误 (unrecoverable checksum error on user table)")
	}
	fmt.Printf("[数据库访问] 成功获取用户 '%s' 的数据\n", userID)
	// ([DB Access] Successfully fetched user data for '%s'\n)
	return fmt.Sprintf("%s 的数据 (Data for %s)", userID, userID), nil
}

// --- 第2层：服务/业务逻辑 --- (--- Layer 2: Service/Business Logic ---)

var ErrUserServiceTemporarilyUnavailable = errors.NewCoder(40001, 503, "用户服务暂时不可用 (User service temporarily unavailable)", "docs/retry-policy_zh.md")

// getUserProfile 检索用户数据并执行一些业务逻辑。
// (getUserProfile retrieves user data and performs some business logic.)
func getUserProfile(userID string) (string, error) {
	fmt.Printf("[服务逻辑] 正在获取用户 '%s' 的配置文件\n", userID)
	// ([Service Logic] Getting profile for user '%s'\n)
	userData, err := fetchUserDataFromDB(userID)
	if err != nil {
		// 特别处理数据库连接错误：可能重试或返回服务不可用错误。
		// (Handle DB connection errors specifically: maybe retry or return a service unavailable error.)
		if errors.IsCode(err, ErrDBConnectionFailed) {
			// 对于此示例，我们将使用服务级别的"暂时不可用" Coder 包装它。
			// (For this example, we\'ll wrap it with a service-level "temporarily unavailable" Coder.)
			// 在实际应用程序中，您可能会在此处实现重试。
			// (In a real app, you might implement retries here.)
			return "", errors.WrapWithCode(err, ErrUserServiceTemporarilyUnavailable, "用户服务无法连接到数据库 (user service could not connect to database)")
		}
		// 对于其他错误 (例如 ErrDBRecordNotFound 或严重的数据库问题)，
		// (For other errors (like ErrDBRecordNotFound or the critical DB issue),)
		// 我们包装它们以添加服务级别的上下文，但让其原始/数据库级别的 Coder (如果存在) 传播。
		// (we wrap them to add service-level context but let their original/DB-level Coder (if any) propagate.)
		return "", errors.Wrapf(err, "获取用户 '%s' 的配置文件失败 (failed to get user profile for '%s')", userID)
	}

	// 业务逻辑：例如，丰富用户数据 (简化)
	// (Business logic: e.g., enrich user data (simplified))
	profile := fmt.Sprintf("配置文件：%s [由服务处理] (Profile: %s [Processed by Service])", userData, userData)
	fmt.Printf("[服务逻辑] 成功处理用户 '%s' 的配置文件\n", userID)
	// ([Service Logic] Successfully processed profile for '%s'\n)
	return profile, nil
}

// --- 第3层：顶层处理程序 (例如，HTTP 处理程序 / 主函数) --- (--- Layer 3: Top-Level Handler (e.g., HTTP Handler / Main Function) ---)

func handleAPIRequest(userID string) {
	fmt.Printf("[API 处理程序] 收到用户 '%s' 的请求\n", userID)
	// ([API Handler] Received request for user '%s'\n)
	profile, err := getUserProfile(userID)

	if err != nil {
		fmt.Printf("[API 处理程序] 处理请求时出错：%v\n", err) // 记录面向用户的消息 (Log the user-facing message)
		// ([API Handler] Error processing request: %v\n)
		// 对于详细日志记录，可以在实际应用程序中使用结构化记录器记录 %+v
		// (For detailed logging, one might use %+v with a structured logger in a real app)
		// log.Error("处理 API 请求失败 (Failed to handle API request)", "error", fmt.Sprintf("%+v", err), "userID", userID)

		coder := errors.GetCoder(err) // 获取链中最新的 Coder (Get the most specific Coder in the chain)
		if coder != nil {
			switch coder.Code() {
			case ErrUserServiceTemporarilyUnavailable.Code():
				fmt.Printf("[API 处理程序] 使用 HTTP %d 响应：%s (参考：%s)\n", coder.HTTPStatus(), coder.String(), coder.Reference())
				// ([API Handler] Responding with HTTP %d: %s (Reference: %s)\n)
				// 以 503 服务不可用响应客户端 (Respond to client with 503 Service Unavailable)
			case ErrDBRecordNotFound.Code(): // 此 Coder 来自数据库层，但已向上传播 (This Coder was from DB layer but propagated up)
				fmt.Printf("[API 处理程序] 使用 HTTP %d 响应：%s\n", coder.HTTPStatus(), coder.String())
				// ([API Handler] Responding with HTTP %d: %s\n)
				// 以 404 未找到响应客户端 (Respond to client with 404 Not Found)
			default:
				// 对于其他编码错误，使用其 HTTP 状态或通用的 500
				// (For other coded errors, use their HTTP status or a generic 500)
				fmt.Printf("[API 处理程序] 使用 HTTP %d 响应 (代码 %d 的默认值)：%s\n", coder.HTTPStatus(), coder.Code(), coder.String())
				// ([API Handler] Responding with HTTP %d (default for code %d): %s\n)
			}
		} else {
			// 未编码错误，以通用的 500 内部服务器错误响应
			// (Uncoded error, respond with a generic 500 Internal Server Error)
			fmt.Printf("[API 处理程序] 使用 HTTP 500 响应：发生意外的内部错误。\n")
			// ([API Handler] Responding with HTTP 500: An unexpected internal error occurred.\n)
			// 内部记录详细错误：fmt.Printf("%+v\n", err)
			// (Log the detailed error internally: fmt.Printf("%+v\n", err))
		}
		return
	}

	fmt.Printf("[API 处理程序] 成功处理请求。响应：%s\n", profile)
	// ([API Handler] Successfully processed request. Response: %s\n)
	// 以 200 OK 和配置文件响应客户端 (Respond to client with 200 OK and profile)
}

func main() {
	scenarios := []string{"user_transient_error", "user_not_found", "user_critical_db_issue", "valid_user"}

	for _, userID := range scenarios {
		fmt.Printf("\n--- 模拟用户 API 请求：'%s' ---\n", userID)
		// (--- Simulating API request for user: '%s' ---\n)
		handleAPIRequest(userID)
		time.Sleep(10 * time.Millisecond) // 为便于阅读输出，稍作延迟 (Small delay for readability of output)
	}
}

/*
示例输出 (简化，此摘要中省略了实际的堆栈跟踪)：
(Example Output (simplified, actual stack traces omitted from this summary)):

--- 模拟用户 API 请求：'user_transient_error' ---
(--- Simulating API request for user: 'user_transient_error' ---)
[API 处理程序] 收到用户 'user_transient_error' 的请求
([API Handler] Received request for user 'user_transient_error')
[服务逻辑] 正在获取用户 'user_transient_error' 的配置文件
([Service Logic] Getting profile for user 'user_transient_error')
[数据库访问] 尝试获取用户 'user_transient_error'
([DB Access] Attempting to fetch user 'user_transient_error')
[API 处理程序] 处理请求时出错：用户服务无法连接到数据库 (user service could not connect to database): fetchUserDataFromDB: 连接数据库副本时网络超时 (network timeout connecting to db replica)
([API Handler] Error processing request: user service could not connect to database: fetchUserDataFromDB: network timeout connecting to db replica)
[API 处理程序] 使用 HTTP 503 响应：用户服务暂时不可用 (User service temporarily unavailable) (参考：docs/retry-policy_zh.md)
([API Handler] Responding with HTTP 503: User service temporarily unavailable (Reference: docs/retry-policy_zh.md))

--- 模拟用户 API 请求：'user_not_found' ---
(--- Simulating API request for user: 'user_not_found' ---)
[API 处理程序] 收到用户 'user_not_found' 的请求
([API Handler] Received request for user 'user_not_found')
[服务逻辑] 正在获取用户 'user_not_found' 的配置文件
([Service Logic] Getting profile for user 'user_not_found')
[数据库访问] 尝试获取用户 'user_not_found'
([DB Access] Attempting to fetch user 'user_not_found')
[API 处理程序] 处理请求时出错：获取用户 'user_not_found' 的配置文件失败 (failed to get user profile for 'user_not_found'): 没有用户 ID 'user_not_found' 的记录 (no record for user ID 'user_not_found')
([API Handler] Error processing request: failed to get user profile for 'user_not_found': no record for user ID 'user_not_found')
[API 处理程序] 使用 HTTP 404 响应：数据库记录未找到 (Database record not found)
([API Handler] Responding with HTTP 404: Database record not found)

--- 模拟用户 API 请求：'user_critical_db_issue' ---
(--- Simulating API request for user: 'user_critical_db_issue' ---)
[API 处理程序] 收到用户 'user_critical_db_issue' 的请求
([API Handler] Received request for user 'user_critical_db_issue')
[服务逻辑] 正在获取用户 'user_critical_db_issue' 的配置文件
([Service Logic] Getting profile for user 'user_critical_db_issue')
[数据库访问] 尝试获取用户 'user_critical_db_issue'
([DB Access] Attempting to fetch user 'user_critical_db_issue')
[API 处理程序] 处理请求时出错：获取用户 'user_critical_db_issue' 的配置文件失败 (failed to get user profile for 'user_critical_db_issue'): 用户表上不可恢复的校验和错误 (unrecoverable checksum error on user table)
([API Handler] Error processing request: failed to get user profile for 'user_critical_db_issue': unrecoverable checksum error on user table)
[API 处理程序] 使用 HTTP 500 响应：发生意外的内部错误。
([API Handler] Responding with HTTP 500: An unexpected internal error occurred.)

--- 模拟用户 API 请求：'valid_user' ---
(--- Simulating API request for user: 'valid_user' ---)
[API 处理程序] 收到用户 'valid_user' 的请求
([API Handler] Received request for user 'valid_user')
[服务逻辑] 正在获取用户 'valid_user' 的配置文件
([Service Logic] Getting profile for user 'valid_user')
[数据库访问] 尝试获取用户 'valid_user'
([DB Access] Attempting to fetch user 'valid_user')
[数据库访问] 成功获取用户 'valid_user' 的数据
([DB Access] Successfully fetched user data for 'valid_user')
[服务逻辑] 成功处理用户 'valid_user' 的配置文件
([Service Logic] Successfully processed profile for 'valid_user')
[API 处理程序] 成功处理请求。响应：配置文件：valid_user 的数据 [由服务处理] (Profile: Data for valid_user [Processed by Service])
([API Handler] Successfully processed request. Response: Profile: Data for valid_user [Processed by Service])
*/
``` 