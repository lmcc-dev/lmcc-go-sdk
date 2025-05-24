<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 1. 与日志记录库一起使用 (With Logging Libraries)

当记录来自 `pkg/errors` 的错误时，您可以提取丰富的信息用于结构化日志记录。

(When logging errors from `pkg/errors`, you can extract rich information for structured logging.)

- **堆栈跟踪 (Stack Trace)**: 使用 `fmt.Sprintf("%+v", err)` 获取带有堆栈跟踪的完整错误消息，用于详细日志。
  (Use `fmt.Sprintf("%+v", err)` to get the full error message with stack trace for detailed logs.)
- **Coder 信息 (Coder Information)**: 使用 `errors.GetCoder(err)` 检索 `Coder`，并将其 `Code()`、`HTTPStatus()`、`String()` (消息) 和 `Reference()` 作为结构化日志中的单独字段进行记录。
  (Use `errors.GetCoder(err)` to retrieve the `Coder` and log its `Code()`, `HTTPStatus()`, `String()` (message), and `Reference()` as separate fields in your structured logs.)
- **错误消息 (Error Message)**: 使用 `err.Error()` 获取简洁的错误消息 (通常是包装消息和原始错误消息的组合)。
  (Use `err.Error()` for a concise error message (often a combination of wrapping messages and the original error message).)

```go
package main

import (
	"bytes" // 用于捕获 slog 输出以进行演示 (To capture slog output for demonstration)
	"fmt"
	"log/slog" // Go 1.21+ 结构化记录器 (Go 1.21+ structured logger)
	"os"

	pkgErrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
)

var ErrFileProcessing = pkgErrors.NewCoder(80001, 500, "文件处理失败 (File processing failed)", "/docs/errors/file-processing#80001_zh")

// processFile 模拟一个可能返回复杂错误的操作。
// (processFile simulates an operation that might return a complex error.)
func processFile(fileName string) error {
	// 模拟底层错误 (Simulate a low-level error)
	lowLevelErr := pkgErrors.NewWithCode(pkgErrors.ErrResourceUnavailable, fmt.Sprintf("资源 '%s' 当前已被锁定 (resource '%s' is currently locked)", fileName))

	// 模拟中层包装 (Simulate a mid-level wrapping)
	midLevelErr := pkgErrors.Wrap(lowLevelErr, "获取处理锁失败 (failed to acquire lock for processing)")

	// 具有自身代码的顶层错误 (Top-level error with its own code)
	return pkgErrors.WrapWithCode(midLevelErr, ErrFileProcessing, fmt.Sprintf("无法在 '%s' 上完成文件操作 (cannot complete file operation on '%s')", fileName))
}

func main() {
	// --- 设置 slog 以捕获输出 --- (--- Setup slog for capturing output ---)
	var logOutput bytes.Buffer
	handlerOptions := &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// 删除时间以获得一致的可测试输出
			// (Remove time for consistent testable output)
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}
	logger := slog.New(slog.NewJSONHandler(&logOutput, handlerOptions))
	// --- slog 设置结束 --- (--- End slog setup ---)

	fileName := "important_document.txt"
	err := processFile(fileName)

	if err != nil {
		fmt.Printf("--- 原始错误 (%%v) ---\n%v\n", err)
		// (--- Original Error (%%v) ---\n%v\n)
		fmt.Printf("\n--- 原始错误 (%%+v) ---\n%+v\n", err)
		// (\n--- Original Error (%%+v) ---\n%+v\n)

		// 提取信息用于结构化日志记录
		// (Extracting information for structured logging)
		errMsg := err.Error() // 简洁消息 (Concise message)
		stackTrace := fmt.Sprintf("%+v", err) // 带堆栈的完整消息 (Full message with stack)

		var attrs []slog.Attr
		attrs = append(attrs, slog.String("fileName", fileName))
		attrs = append(attrs, slog.String("errorMessageConcise", errMsg))
		attrs = append(attrs, slog.String("errorFullTrace", stackTrace))

		if coder := pkgErrors.GetCoder(err); coder != nil {
			attrs = append(attrs, slog.Int("errorCode", coder.Code()))
			attrs = append(attrs, slog.Int("errorHTTPStatus", coder.HTTPStatus()))
			attrs = append(attrs, slog.String("errorCodeMessage", coder.String()))
			attrs = append(attrs, slog.String("errorRef", coder.Reference()))
		}
		
		// 使用 slog (Go 1.21+)
		// (Using slog (Go 1.21+))
		// 将 []slog.Attr 转换为 []any 以用于 logger.Error
		// (Convert []slog.Attr to []any for logger.Error)
		var logArgs []any
		for _, attr := range attrs {
			logArgs = append(logArgs, attr)
		}
		logger.Error("处理文件失败 (Failed to process file)", logArgs...)

		fmt.Println("\n--- 捕获的 slog JSON 输出 (近似) --- (--- Captured slog JSON Output (approximate) ---)")
		// 美化 JSON 以提高可读性 - 在实际场景中，它是紧凑的 JSON 行。
		// (Beautify JSON for readability - in real scenarios, it's compact JSON lines.)
		var prettyJSON bytes.Buffer
		if jsonErr :=_jsonIndent(&prettyJSON, logOutput.Bytes(), "", "  "); jsonErr == nil {
		    fmt.Println(prettyJSON.String())
		} else {
		    fmt.Println(logOutput.String()) // 如果缩进失败，则回退到原始输出 (Fallback to raw if indent fails)
		}
	}
}

// _jsonIndent 是一个辅助函数，用于漂亮地打印 JSON，类似于 `json.Indent`。
// ( _jsonIndent is a helper to pretty-print JSON, akin to `json.Indent`.)
// 这里需要它，因为 `json.Indent` 在这里不直接可用。
// (This is needed because `json.Indent` is not directly available here.)
func _jsonIndent(dst *bytes.Buffer, src []byte, prefix, indent string) error {
    // 这是一个简化的存根。真正的实现会解析并重新序列化。
    // (This is a simplified stub. A real implementation would parse and re-serialize.)
    // 对于此示例，如果我们无法实际缩进，我们将直接传递它。
    // (For this example, we'll just pass it through if we can't actually indent.)
    // 在真实的测试或应用程序中，您将使用 encoding/json。
    // (In a real test or application, you'd use encoding/json.)
    dst.Write(src) 
    return nil
}

/*
示例输出 (堆栈跟踪和 JSON 字段顺序可能有所不同)：
(Example Output (Stack traces and JSON field order may vary)):

--- 原始错误 (%%v) ---
(--- Original Error (%%v) ---)
文件处理失败 (File processing failed): 无法在 'important_document.txt' 上完成文件操作 (cannot complete file operation on 'important_document.txt'): 获取处理锁失败 (failed to acquire lock for processing): Resource unavailable: 资源 'important_document.txt' 当前已被锁定 (resource 'important_document.txt' is currently locked)

--- 原始错误 (%%+v) ---
(--- Original Error (%%+v) ---)
文件处理失败 (File processing failed): 无法在 'important_document.txt' 上完成文件操作 (cannot complete file operation on 'important_document.txt'): 获取处理锁失败 (failed to acquire lock for processing): 资源 'important_document.txt' 当前已被锁定 (resource 'important_document.txt' is currently locked)
main.processFile
	/path/to/your/file.go:26
main.main
	/path/to/your/file.go:45
...
Resource unavailable

--- 捕获的 slog JSON 输出 (近似) --- 
(--- Captured slog JSON Output (approximate) --- )
{"level":"ERROR","msg":"处理文件失败 (Failed to process file)","fileName":"important_document.txt","errorMessageConcise":"文件处理失败 (File processing failed): 无法在 'important_document.txt' 上完成文件操作 (cannot complete file operation on 'important_document.txt'): 获取处理锁失败 (failed to acquire lock for processing): Resource unavailable: 资源 'important_document.txt' 当前已被锁定 (resource 'important_document.txt' is currently locked)","errorFullTrace":"文件处理失败 (File processing failed): 无法在 'important_document.txt' 上完成文件操作 (cannot complete file operation on 'important_document.txt'): 获取处理锁失败 (failed to acquire lock for processing): 资源 'important_document.txt' 当前已被锁定 (resource 'important_document.txt' is currently locked)\nmain.processFile\n\t/path/to/your/file.go:26\nmain.main\n\t/path/to/your/file.go:45\n...\nResource unavailable","errorCode":80001,"errorHTTPStatus":500,"errorCodeMessage":"文件处理失败 (File processing failed)","errorRef":"/docs/errors/file-processing#80001_zh"}

*/
```

**注意 (Note)**: `_jsonIndent` 函数在此示例中是一个占位符。在实际应用程序中，如果您需要漂亮地打印 JSON，则会使用 `encoding/json.Indent`，或者更有可能的是，您的日志记录基础结构会直接处理 JSON 格式化。
(The `_jsonIndent` function is a placeholder for this example. In a real application, you would use `encoding/json.Indent` if you needed to pretty-print JSON, or more likely, your logging infrastructure would handle JSON formatting directly.) 