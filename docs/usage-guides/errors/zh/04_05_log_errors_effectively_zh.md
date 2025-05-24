<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 5. 有效地记录错误 (Log Errors Effectively)

在适当的最高级别记录错误，通常是在错误得到处理且不再进一步传播的地方 (例如，在您的主函数、HTTP 处理程序或后台作业处理器中)。

(Log errors at the highest appropriate level, typically where the error is handled and not propagated further (e.g., in your main function, HTTP handlers, or background job processors).)

- **为意外错误包含堆栈跟踪 (Include Stack Traces for Unexpected Errors)**: 对于意外错误或具有严重影响的错误，请使用 `fmt.Printf("%+v\n", err)` 或支持它的结构化记录器记录完整的堆栈跟踪。这对于调试至关重要。
  (For unexpected errors or those with severe impact, log the full stack trace using `fmt.Printf("%+v\n", err)` or a structured logger that supports it. This is crucial for debugging.)
- **结构化日志记录 (Structured Logging)**: 使用结构化日志记录库 (例如 Go 1.21+ 中的 `log/slog`，或其他第三方库) 并包含相关上下文 (例如，请求 ID、用户 ID、Coder 信息) 作为结构化字段。这使日志更易于解析、搜索和分析。
  (Use a structured logging library (like `log/slog` in Go 1.21+, or other third-party libraries) and include relevant context (e.g., request ID, user ID, Coder information) as structured fields. This makes logs easier to parse, search, and analyze.)
- **避免冗余日志记录 (Avoid Redundant Logging)**: 不要在调用堆栈的多个级别记录相同的错误。这会产生嘈杂的日志，并使追踪错误的原始来源更加困难。
  (Don't log the same error at multiple levels of the call stack. This creates noisy logs and makes it harder to trace the original source of the error.)

```go
package main

import (
	"fmt"
	"os"
	// "log/slog" // 如果使用 Go 1.21+，取消注释以使用结构化日志记录
	// (// Uncomment for Go 1.21+ to use structured logging)
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"time"
)

// 对于 Go 1.21+ slog 示例 (如果不使用 slog，请替换为您首选的结构化记录器)
// (For Go 1.21+ slog example (replace with your preferred structured logger if not using slog))
/*
var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))
*/

// 如果不使用 slog，用于演示的简单基于文件的记录器
// (Simple file-based logger for demonstration if not using slog)
func simpleFileLogger(format string, args ...interface{}) {
	f, _ := os.OpenFile("app_zh.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if f != nil {
		defer f.Close()
		fmt.Fprintf(f, "[%s] %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(format, args...))
	}
}

var ErrDataProcessing = errors.NewCoder(70001, 500, "数据处理管道失败 (Data processing pipeline failed)", "")

// processData 模拟一个可能失败的函数。
// (processData simulates a function that might fail.)
func processData(dataID string, failStep string) error {
	fmt.Printf("[工作程序] 开始处理数据：%s\n", dataID)
	// ([Worker] Starting to process data: %s\n)
	time.Sleep(5 * time.Millisecond) // 模拟工作 (Simulate work)

	if failStep == "step1" {
		originalErr := errors.New("原始数据验证失败：缺少关键字段 (raw data validation failed: missing crucial field)")
		return errors.WrapWithCode(originalErr, errors.ErrValidation, fmt.Sprintf("数据 '%s' 的验证在步骤1失败 (validation of data '%s' failed at step 1)", dataID))
	}

	fmt.Printf("[工作程序] 数据 %s 的步骤1已完成\n", dataID)
	// ([Worker] Step 1 completed for data: %s\n)
	time.Sleep(5 * time.Millisecond)

	if failStep == "step2" {
		// 带有自定义 Coder 的更复杂错误
		// (More complex error with a custom Coder)
		underlyingIssue := errors.New("外部 API 在3次重试后超时 (external API timed out after 3 retries)")
		wrappedErr := errors.Wrap(underlyingIssue, "从外部服务丰富数据失败 (failed to enrich data from external service)")
		return errors.WrapWithCode(wrappedErr, ErrDataProcessing, fmt.Sprintf("数据 '%s' 的处理管道在步骤2失败 (processing pipeline for data '%s' failed at step 2)", dataID))
	}

	fmt.Printf("[工作程序] 数据 '%s' 已成功处理。\n", dataID)
	// ([Worker] Data '%s' processed successfully.\n)
	return nil
}

// handleRequest 模拟一个顶层处理程序 (例如，HTTP 处理程序或作业消费者)。
// (handleRequest simulates a top-level handler (e.g., an HTTP handler or a job consumer).)
func handleRequest(requestID string, dataID string, failAtStep string) {
	fmt.Printf("[处理程序] 收到请求 %s 以处理数据 %s (模拟在 '%s' 失败)\n", requestID, dataID, failAtStep)
	// ([Handler] Received request %s to process data %s (simulating fail at '%s')\n)

	err := processData(dataID, failAtStep)

	if err != nil {
		// 在此处，在边界处记录错误。
		// (Log the error here, at the boundary.)
		fmt.Printf("[处理程序] 处理请求 %s (数据 %s) 时发生错误。\n", requestID, dataID)
		// ([Handler] Error occurred processing request %s for data %s.\n)

		// --- 简单日志记录 (fmt.Printf) --- (--- Simple Logging (fmt.Printf) ---)
		// 面向用户的消息 (较少细节)：
		// (For user-facing message (less detail):)
		fmt.Printf("  面向用户的错误消息：%v\n", err)
		// (User-facing error message: %v\n)
		// 详细的内部日志记录 (带堆栈跟踪)：
		// (For detailed internal logging (with stack trace):)
		fmt.Printf("  详细的内部日志 (%%+v)：\n%+v\n", err)
		// (Detailed internal log (%%+v):\n%+v\n)

		// --- 结构化日志记录示例 (使用 simpleFileLogger，slog 的概念性示例) --- 
		// (--- Structured Logging Example (using simpleFileLogger, conceptual for slog) ---)
		coder := errors.GetCoder(err)
		if coder != nil {
			simpleFileLogger("处理请求时出错。请求ID：%s，数据ID：%s，错误代码：%d，错误HTTP状态：%d，错误消息：%s，错误参考：%s，堆栈跟踪：\n%+v",
				requestID, dataID, coder.Code(), coder.HTTPStatus(), coder.String(), coder.Reference(), err)
			// (Error processing request. RequestID: %s, DataID: %s, ErrorCode: %d, ErrorHTTPStatus: %d, ErrorMessage: %s, ErrorRef: %s, StackTrace: \n%+v)
			/* // Go 1.21+ 的 slog 等效代码
			logger.Error("处理请求时出错 (Error processing request)",
				slog.String("requestID", requestID),
				slog.String("dataID", dataID),
				slog.Int("errorCode", coder.Code()),
				slog.Int("errorHTTPStatus", coder.HTTPStatus()),
				slog.String("errorMessage", err.Error()), // 或 coder.String() 获取基本消息 (or coder.String() for base message)
				slog.String("errorRef", coder.Reference()),
				slog.String("stackTrace", fmt.Sprintf("%+v", err)),
			)
			*/
		} else {
			simpleFileLogger("处理请求时出错 (未编码)。请求ID：%s，数据ID：%s，错误消息：%s，堆栈跟踪：\n%+v",
				requestID, dataID, err.Error(), err)
			// (Error processing request (uncoded). RequestID: %s, DataID: %s, ErrorMessage: %s, StackTrace: \n%+v)
			/* // Go 1.21+ 的 slog 等效代码
			logger.Error("处理请求时出错 (未编码) (Error processing request (uncoded))",
				slog.String("requestID", requestID),
				slog.String("dataID", dataID),
				slog.String("errorMessage", err.Error()),
				slog.String("stackTrace", fmt.Sprintf("%+v", err)),
			)
			*/
		}
		// 适当地响应客户端 (例如，HTTP 500 或特定的错误代码)
		// (Respond to client appropriately (e.g., HTTP 500 or specific error code))
		fmt.Printf("[处理程序] 正在向客户端响应请求 %s 的错误状态。\n", requestID)
		// ([Handler] Responding to client with an error status for request %s.\n)
		return
	}

	fmt.Printf("[处理程序] 请求 %s (数据 %s) 已成功处理！\n", requestID, dataID)
	// ([Handler] Request %s for data %s processed successfully!\n)
	// 成功响应客户端 (Respond to client with success)
}

func main() {
	// 开始时清除日志文件以获得更清晰的演示输出
	// (Clear log file at start for cleaner demo output)
	_ = os.Remove("app_zh.log")

	handleRequest("req-001", "data-alpha", "none")
	fmt.Println("---")
	handleRequest("req-002", "data-beta", "step1") 
	fmt.Println("---")
	handleRequest("req-003", "data-gamma", "step2")

	fmt.Println("\n日志文件内容 (app_zh.log)：(Log file content (app_zh.log):)")
	content, _ := os.ReadFile("app_zh.log")
	fmt.Print(string(content))
	_ = os.Remove("app_zh.log") // 清理日志文件 (Clean up log file)
}

/*
示例控制台输出 (%%+v 中的堆栈跟踪非常详细，并且会因情况而异)：
(Example Console Output (Stack traces in %%+v are verbose and will vary)):

[处理程序] 收到请求 req-001 以处理数据 data-alpha (模拟在 'none' 失败)
([Handler] Received request req-001 to process data data-alpha (simulating fail at 'none'))
[工作程序] 开始处理数据：data-alpha
([Worker] Starting to process data: data-alpha)
[工作程序] 数据 data-alpha 的步骤1已完成
([Worker] Step 1 completed for data: data-alpha)
[工作程序] 数据 'data-alpha' 已成功处理。
([Worker] Data 'data-alpha' processed successfully.)
[处理程序] 请求 req-001 (数据 data-alpha) 已成功处理！
([Handler] Request req-001 for data data-alpha processed successfully!)
---
[处理程序] 收到请求 req-002 以处理数据 data-beta (模拟在 'step1' 失败)
([Handler] Received request req-002 to process data data-beta (simulating fail at 'step1'))
[工作程序] 开始处理数据：data-beta
([Worker] Starting to process data: data-beta)
[处理程序] 处理请求 req-002 (数据 data-beta) 时发生错误。
([Handler] Error occurred processing request req-002 for data data-beta.)
  面向用户的错误消息：Validation failed: 数据 'data-beta' 的验证在步骤1失败 (validation of data 'data-beta' failed at step 1): 原始数据验证失败：缺少关键字段 (raw data validation failed: missing crucial field)
  (User-facing error message: Validation failed: validation of data 'data-beta' failed at step 1: raw data validation failed: missing crucial field)
  详细的内部日志 (%%+v)：
  (Detailed internal log (%%+v):)
数据 'data-beta' 的验证在步骤1失败 (validation of data 'data-beta' failed at step 1): 原始数据验证失败：缺少关键字段 (raw data validation failed: missing crucial field)
main.processData
	/path/to/your/file.go:34
main.handleRequest
	/path/to/your/file.go:62
main.main
	/path/to/your/file.go:108
...
Validation failed
[处理程序] 正在向客户端响应请求 req-002 的错误状态。
([Handler] Responding to client with an error status for request req-002.)
---
[处理程序] 收到请求 req-003 以处理数据 data-gamma (模拟在 'step2' 失败)
([Handler] Received request req-003 to process data data-gamma (simulating fail at 'step2'))
[工作程序] 开始处理数据：data-gamma
([Worker] Starting to process data: data-gamma)
[工作程序] 数据 data-gamma 的步骤1已完成
([Worker] Step 1 completed for data: data-gamma)
[处理程序] 处理请求 req-003 (数据 data-gamma) 时发生错误。
([Handler] Error occurred processing request req-003 for data data-gamma.)
  面向用户的错误消息：数据处理管道失败 (Data processing pipeline failed): 数据 'data-gamma' 的处理管道在步骤2失败 (processing pipeline for data 'data-gamma' failed at step 2): 从外部服务丰富数据失败 (failed to enrich data from external service): 外部 API 在3次重试后超时 (external API timed out after 3 retries)
  (User-facing error message: Data processing pipeline failed: processing pipeline for data 'data-gamma' failed at step 2: failed to enrich data from external service: external API timed out after 3 retries)
  详细的内部日志 (%%+v)：
  (Detailed internal log (%%+v):)
数据 'data-gamma' 的处理管道在步骤2失败 (processing pipeline for data 'data-gamma' failed at step 2): 从外部服务丰富数据失败 (failed to enrich data from external service): 外部 API 在3次重试后超时 (external API timed out after 3 retries)
main.processData
	/path/to/your/file.go:45
main.handleRequest
	/path/to/your/file.go:62
main.main
	/path/to/your/file.go:110
...
Data processing pipeline failed
[处理程序] 正在向客户端响应请求 req-003 的错误状态。
([Handler] Responding to client with an error status for request req-003.)

日志文件内容 (app_zh.log)：(Log file content (app_zh.log):)
[TIMESTAMP] 处理请求时出错。请求ID：req-002，数据ID：data-beta，错误代码：100006，错误HTTP状态：400，错误消息：Validation failed，错误参考：errors-spec.md#100006，堆栈跟踪： 
(Error processing request. RequestID: req-002, DataID: data-beta, ErrorCode: 100006, ErrorHTTPStatus: 400, ErrorMessage: Validation failed, ErrorRef: errors-spec.md#100006, StackTrace: )
数据 'data-beta' 的验证在步骤1失败 (validation of data 'data-beta' failed at step 1): 原始数据验证失败：缺少关键字段 (raw data validation failed: missing crucial field)
main.processData
	/path/to/your/file.go:34
main.handleRequest
	/path/to/your/file.go:62
main.main
	/path/to/your/file.go:108
...
Validation failed
[TIMESTAMP] 处理请求时出错。请求ID：req-003，数据ID：data-gamma，错误代码：70001，错误HTTP状态：500，错误消息：数据处理管道失败 (Data processing pipeline failed)，错误参考：，堆栈跟踪： 
(Error processing request. RequestID: req-003, DataID: data-gamma, ErrorCode: 70001, ErrorHTTPStatus: 500, ErrorMessage: Data processing pipeline failed, ErrorRef: , StackTrace: )
数据 'data-gamma' 的处理管道在步骤2失败 (processing pipeline for data 'data-gamma' failed at step 2): 从外部服务丰富数据失败 (failed to enrich data from external service): 外部 API 在3次重试后超时 (external API timed out after 3 retries)
main.processData
	/path/to/your/file.go:45
main.handleRequest
	/path/to/your/file.go:62
main.main
	/path/to/your/file.go:110
...
Data processing pipeline failed
*/
``` 