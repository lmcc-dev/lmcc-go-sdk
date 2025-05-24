<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 5. Log Errors Effectively

Log errors at the highest appropriate level, typically where the error is handled and not propagated further (e.g., in your main function, HTTP handlers, or background job processors).

- **Include Stack Traces for Unexpected Errors**: For unexpected errors or those with severe impact, log the full stack trace using `fmt.Printf("%+v\n", err)` or a structured logger that supports it. This is crucial for debugging.
- **Structured Logging**: Use a structured logging library (like `log/slog` in Go 1.21+, or other third-party libraries) and include relevant context (e.g., request ID, user ID, Coder information) as structured fields. This makes logs easier to parse, search, and analyze.
- **Avoid Redundant Logging**: Don't log the same error at multiple levels of the call stack. This creates noisy logs and makes it harder to trace the original source of the error.

```go
package main

import (
	"fmt"
	"os"
	// "log/slog" // Uncomment for Go 1.21+ to use structured logging
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"time"
)

// For Go 1.21+ slog example (replace with your preferred structured logger if not using slog)
/*
var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))
*/

// Simple file-based logger for demonstration if not using slog
func simpleFileLogger(format string, args ...interface{}) {
	f, _ := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if f != nil {
		defer f.Close()
		fmt.Fprintf(f, "[%s] %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(format, args...))
	}
}

var ErrDataProcessing = errors.NewCoder(70001, 500, "Data processing pipeline failed", "")

// processData simulates a function that might fail.
func processData(dataID string, failStep string) error {
	fmt.Printf("[Worker] Starting to process data: %s\n", dataID)
	time.Sleep(5 * time.Millisecond) // Simulate work

	if failStep == "step1" {
		originalErr := errors.New("raw data validation failed: missing crucial field")
		return errors.WrapWithCode(originalErr, errors.ErrValidation, fmt.Sprintf("validation of data '%s' failed at step 1", dataID))
	}

	fmt.Printf("[Worker] Step 1 completed for data: %s\n", dataID)
	time.Sleep(5 * time.Millisecond)

	if failStep == "step2" {
		// More complex error with a custom Coder
		underlyingIssue := errors.New("external API timed out after 3 retries")
		wrappedErr := errors.Wrap(underlyingIssue, "failed to enrich data from external service")
		return errors.WrapWithCode(wrappedErr, ErrDataProcessing, fmt.Sprintf("processing pipeline for data '%s' failed at step 2", dataID))
	}

	fmt.Printf("[Worker] Data '%s' processed successfully.\n", dataID)
	return nil
}

// handleRequest simulates a top-level handler (e.g., an HTTP handler or a job consumer).
func handleRequest(requestID string, dataID string, failAtStep string) {
	fmt.Printf("[Handler] Received request %s to process data %s (simulating fail at '%s')\n", requestID, dataID, failAtStep)

	err := processData(dataID, failAtStep)

	if err != nil {
		// Log the error here, at the boundary.
		fmt.Printf("[Handler] Error occurred processing request %s for data %s.\n", requestID, dataID)

		// --- Simple Logging (fmt.Printf) ---
		// For user-facing message (less detail):
		fmt.Printf("  User-facing error message: %v\n", err)
		// For detailed internal logging (with stack trace):
		fmt.Printf("  Detailed internal log (%%+v):\n%+v\n", err)

		// --- Structured Logging Example (using simpleFileLogger, conceptual for slog) ---
		coder := errors.GetCoder(err)
		if coder != nil {
			simpleFileLogger("Error processing request. RequestID: %s, DataID: %s, ErrorCode: %d, ErrorHTTPStatus: %d, ErrorMessage: %s, ErrorRef: %s, StackTrace: \n%+v",
				requestID, dataID, coder.Code(), coder.HTTPStatus(), coder.String(), coder.Reference(), err)
			/* // slog equivalent for Go 1.21+
			logger.Error("Error processing request",
				slog.String("requestID", requestID),
				slog.String("dataID", dataID),
				slog.Int("errorCode", coder.Code()),
				slog.Int("errorHTTPStatus", coder.HTTPStatus()),
				slog.String("errorMessage", err.Error()), // or coder.String() for base message
				slog.String("errorRef", coder.Reference()),
				slog.String("stackTrace", fmt.Sprintf("%+v", err)),
			)
			*/
		} else {
			simpleFileLogger("Error processing request (uncoded). RequestID: %s, DataID: %s, ErrorMessage: %s, StackTrace: \n%+v",
				requestID, dataID, err.Error(), err)
			/* // slog equivalent for Go 1.21+
			logger.Error("Error processing request (uncoded)",
				slog.String("requestID", requestID),
				slog.String("dataID", dataID),
				slog.String("errorMessage", err.Error()),
				slog.String("stackTrace", fmt.Sprintf("%+v", err)),
			)
			*/
		}
		// Respond to client appropriately (e.g., HTTP 500 or specific error code)
		fmt.Printf("[Handler] Responding to client with an error status for request %s.\n", requestID)
		return
	}

	fmt.Printf("[Handler] Request %s for data %s processed successfully!\n", requestID, dataID)
	// Respond to client with success
}

func main() {
	// Clear log file at start for cleaner demo output
	_ = os.Remove("app.log")

	handleRequest("req-001", "data-alpha", "none")
	fmt.Println("---")
	handleRequest("req-002", "data-beta", "step1") 
	fmt.Println("---")
	handleRequest("req-003", "data-gamma", "step2")

	fmt.Println("\nLog file content (app.log):")
	content, _ := os.ReadFile("app.log")
	fmt.Print(string(content))
	_ = os.Remove("app.log") // Clean up log file
}

/*
Example Console Output (Stack traces in %+v are verbose and will vary):

[Handler] Received request req-001 to process data data-alpha (simulating fail at 'none')
[Worker] Starting to process data: data-alpha
[Worker] Step 1 completed for data: data-alpha
[Worker] Data 'data-alpha' processed successfully.
[Handler] Request req-001 for data data-alpha processed successfully!
---
[Handler] Received request req-002 to process data data-beta (simulating fail at 'step1')
[Worker] Starting to process data: data-beta
[Handler] Error occurred processing request req-002 for data data-beta.
  User-facing error message: Validation failed: validation of data 'data-beta' failed at step 1: raw data validation failed: missing crucial field
  Detailed internal log (%+v):
validation of data 'data-beta' failed at step 1: raw data validation failed: missing crucial field
main.processData
	/path/to/your/file.go:34
main.handleRequest
	/path/to/your/file.go:62
main.main
	/path/to/your/file.go:108
...
Validation failed
[Handler] Responding to client with an error status for request req-002.
---
[Handler] Received request req-003 to process data data-gamma (simulating fail at 'step2')
[Worker] Starting to process data: data-gamma
[Worker] Step 1 completed for data: data-gamma
[Handler] Error occurred processing request req-003 for data data-gamma.
  User-facing error message: Data processing pipeline failed: processing pipeline for data 'data-gamma' failed at step 2: failed to enrich data from external service: external API timed out after 3 retries
  Detailed internal log (%+v):
processing pipeline for data 'data-gamma' failed at step 2: failed to enrich data from external service: external API timed out after 3 retries
main.processData
	/path/to/your/file.go:45
main.handleRequest
	/path/to/your/file.go:62
main.main
	/path/to/your/file.go:110
...
Data processing pipeline failed
[Handler] Responding to client with an error status for request req-003.

Log file content (app.log):
[TIMESTAMP] Error processing request. RequestID: req-002, DataID: data-beta, ErrorCode: 100006, ErrorHTTPStatus: 400, ErrorMessage: Validation failed, ErrorRef: errors-spec.md#100006, StackTrace: 
validation of data 'data-beta' failed at step 1: raw data validation failed: missing crucial field
main.processData
	/path/to/your/file.go:34
main.handleRequest
	/path/to/your/file.go:62
main.main
	/path/to/your/file.go:108
...
Validation failed
[TIMESTAMP] Error processing request. RequestID: req-003, DataID: data-gamma, ErrorCode: 70001, ErrorHTTPStatus: 500, ErrorMessage: Data processing pipeline failed, ErrorRef: , StackTrace: 
processing pipeline for data 'data-gamma' failed at step 2: failed to enrich data from external service: external API timed out after 3 retries
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