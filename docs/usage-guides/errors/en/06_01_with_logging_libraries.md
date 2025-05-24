<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 1. With Logging Libraries

When logging errors from `pkg/errors`, you can extract rich information for structured logging.

- **Stack Trace**: Use `fmt.Sprintf("%+v", err)` to get the full error message with stack trace for detailed logs.
- **Coder Information**: Use `errors.GetCoder(err)` to retrieve the `Coder` and log its `Code()`, `HTTPStatus()`, `String()` (message), and `Reference()` as separate fields in your structured logs.
- **Error Message**: Use `err.Error()` for a concise error message (often a combination of wrapping messages and the original error message).

```go
package main

import (
	"fmt"
	"log/slog" // Go 1.21+ structured logger
	"os"
	"bytes" // To capture slog output for demonstration

	pkgErrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
)

var ErrFileProcessing = pkgErrors.NewCoder(80001, 500, "File processing failed", "/docs/errors/file-processing#80001")

// processFile simulates an operation that might return a complex error.
func processFile(fileName string) error {
	// Simulate a low-level error
	lowLevelErr := pkgErrors.NewWithCode(pkgErrors.ErrResourceUnavailable, fmt.Sprintf("resource '%s' is currently locked", fileName))

	// Simulate a mid-level wrapping
	midLevelErr := pkgErrors.Wrap(lowLevelErr, "failed to acquire lock for processing")

	// Top-level error with its own code
	return pkgErrors.WrapWithCode(midLevelErr, ErrFileProcessing, fmt.Sprintf("cannot complete file operation on '%s'", fileName))
}

func main() {
	// --- Setup slog for capturing output --- 
	var logOutput bytes.Buffer
	handlerOptions := &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove time for consistent testable output
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}
	logger := slog.New(slog.NewJSONHandler(&logOutput, handlerOptions))
	// --- End slog setup ---

	fileName := "important_document.txt"
	err := processFile(fileName)

	if err != nil {
		fmt.Printf("--- Original Error (%%v) ---\n%v\n", err)
		fmt.Printf("\n--- Original Error (%%+v) ---\n%+v\n", err)

		// Extracting information for structured logging
		errMsg := err.Error() // Concise message
		stackTrace := fmt.Sprintf("%+v", err) // Full message with stack

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
		
		// Using slog (Go 1.21+)
		// Convert []slog.Attr to []any for logger.Error
		var logArgs []any
		for _, attr := range attrs {
			logArgs = append(logArgs, attr)
		}
		logger.Error("Failed to process file", logArgs...)

		fmt.Println("\n--- Captured slog JSON Output (approximate) --- ")
		// Beautify JSON for readability - in real scenarios, it's compact JSON lines.
		var prettyJSON bytes.Buffer
		if jsonErr :=_jsonIndent(&prettyJSON, logOutput.Bytes(), "", "  "); jsonErr == nil {
		    fmt.Println(prettyJSON.String())
		} else {
		    fmt.Println(logOutput.String()) // Fallback to raw if indent fails
		}
	}
}

// _jsonIndent is a helper to pretty-print JSON, akin to `json.Indent`.
// This is needed because `json.Indent` is not directly available here.
func _jsonIndent(dst *bytes.Buffer, src []byte, prefix, indent string) error {
    // This is a simplified stub. A real implementation would parse and re-serialize.
    // For this example, we'll just pass it through if we can't actually indent.
    // In a real test or application, you'd use encoding/json.
    dst.Write(src) 
    return nil
}

/*
Example Output (Stack traces and JSON field order may vary):

--- Original Error (%v) ---
File processing failed: cannot complete file operation on 'important_document.txt': failed to acquire lock for processing: Resource unavailable: resource 'important_document.txt' is currently locked

--- Original Error (%+v) ---
File processing failed: cannot complete file operation on 'important_document.txt': failed to acquire lock for processing: resource 'important_document.txt' is currently locked
main.processFile
	/path/to/your/file.go:26
main.main
	/path/to/your/file.go:45
...
Resource unavailable

--- Captured slog JSON Output (approximate) --- 
{"level":"ERROR","msg":"Failed to process file","fileName":"important_document.txt","errorMessageConcise":"File processing failed: cannot complete file operation on 'important_document.txt': failed to acquire lock for processing: Resource unavailable: resource 'important_document.txt' is currently locked","errorFullTrace":"File processing failed: cannot complete file operation on 'important_document.txt': failed to acquire lock for processing: resource 'important_document.txt' is currently locked\nmain.processFile\n\t/path/to/your/file.go:26\nmain.main\n\t/path/to/your/file.go:45\n...\nResource unavailable","errorCode":80001,"errorHTTPStatus":500,"errorCodeMessage":"File processing failed","errorRef":"/docs/errors/file-processing#80001"}

*/
```

**Note**: The `_jsonIndent` function is a placeholder for this example. In a real application, you would use `encoding/json.Indent` if you needed to pretty-print JSON, or more likely, your logging infrastructure would handle JSON formatting directly. 