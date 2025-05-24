<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 3. Using Error Codes

Error codes allow for programmatic handling of errors. The `pkg/errors` module provides a `Coder` interface and functions to create errors with codes or attach codes to existing errors.

**`errors.Coder` interface**: Defines methods `Code() int`, `HTTPStatus() int`, `String() string`, and `Reference() string`.

**`errors.NewCoder(code int, httpStatus int, message string, reference string) Coder`**: Creates a new `Coder`.

**`errors.NewWithCode(coder Coder, message string) error`**: Creates a new error with the given `Coder` and message.

**`errors.ErrorfWithCode(coder Coder, format string, args ...interface{}) error`**: Creates a new formatted error with the given `Coder`.

**`errors.WithCode(err error, coder Coder) error`**: Attaches a `Coder` to an existing error. If `err` is nil, it returns nil.


```go
package main

import (
	"fmt"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	standardErrors "errors" // For a sample pre-existing error
)

// Define a custom error code
// It's good practice to define these as constants or package-level variables.
var ErrCustomServiceUnavailable = errors.NewCoder(
	50001, // Unique application-specific code
	503,   // Corresponding HTTP status
	"Custom service is currently unavailable", // Default message for this code
	"doc/link/to/service_unavailable.md", // Optional reference link
)

// simulateFetchResource
// Demonstrates creating an error with a predefined Coder.
func simulateFetchResource(resourceID string) error {
	if resourceID == "nonexistent" {
		// Use a predefined Coder from the errors package (e.g., errors.ErrNotFound)
		return errors.NewWithCode(errors.ErrNotFound, fmt.Sprintf("resource with ID '%s' was not found", resourceID))
	}
	fmt.Printf("Resource '%s' fetched successfully.\n", resourceID)
	return nil
}

// simulateExternalServiceCall
// Demonstrates creating an error with a custom Coder.
func simulateExternalServiceCall(serviceName string) error {
	if serviceName == "unreliable_service" {
		// Use our custom Coder
		return errors.ErrorfWithCode(ErrCustomServiceUnavailable, "failed to call %s due to intermittent issues", serviceName)
	}
	fmt.Printf("External service '%s' called successfully.\n", serviceName)
	return nil
}

// processData
// Demonstrates attaching a Coder to an existing error.
func processData(data string) error {
	var existingError error
	if data == "corrupted" {
		existingError = standardErrors.New("underlying data corruption detected") // A pre-existing error
	} else if data == "forbidden_data" {
		existingError = errors.New("access to this data is forbidden by policy") // An lmcc error without a code yet
	}


	if existingError != nil {
		// Attach our custom Coder to this existing error.
		// This is useful if the Coder is determined based on the context where the error is handled.
		return errors.WithCode(existingError, errors.ErrBadRequest) // Using predefined ErrBadRequest
	}

	fmt.Printf("Data '%s' processed successfully.\n", data)
	return nil
}


func main() {
	fmt.Println("--- Example 1: NewWithCode (predefined Coder) ---")
	err1 := simulateFetchResource("nonexistent")
	if err1 != nil {
		fmt.Printf("Error (%%v): %v\n", err1)
		fmt.Printf("Error (%%+v):\n%+v\n", err1)
		if coder := errors.GetCoder(err1); coder != nil {
			fmt.Printf("  Code: %d, HTTP Status: %d, Message: %s, Reference: %s\n",
				coder.Code(), coder.HTTPStatus(), coder.String(), coder.Reference())
		}
	}

	fmt.Println("\n--- Example 2: ErrorfWithCode (custom Coder) ---")
	err2 := simulateExternalServiceCall("unreliable_service")
	if err2 != nil {
		fmt.Printf("Error (%%v): %v\n", err2)
		fmt.Printf("Error (%%+v):\n%+v\n", err2)
		if coder := errors.GetCoder(err2); coder != nil {
			fmt.Printf("  Code: %d, HTTP Status: %d, Message: %s, Reference: %s\n",
				coder.Code(), coder.HTTPStatus(), coder.String(), coder.Reference())
		}
	}
	
	fmt.Println("\n--- Example 3: WithCode (attaching Coder to existing error) ---")
	err3 := processData("corrupted")
	if err3 != nil {
		fmt.Printf("Error (%%v): %v\n", err3)
		fmt.Printf("Error (%%+v):\n%+v\n", err3)
		if coder := errors.GetCoder(err3); coder != nil {
			fmt.Printf("  Code: %d, HTTP Status: %d, Message: %s, Reference: %s\n",
				coder.Code(), coder.HTTPStatus(), coder.String(), coder.Reference())
		}
	}

	fmt.Println("\n--- Example 4: WithCode (attaching Coder to an lmcc error) ---")
	err4 := processData("forbidden_data")
	if err4 != nil {
		fmt.Printf("Error (%%v): %v\n", err4)
		fmt.Printf("Error (%%+v):\n%+v\n", err4)
		if coder := errors.GetCoder(err4); coder != nil {
			fmt.Printf("  Code: %d, HTTP Status: %d, Message: %s, Reference: %s\n",
				coder.Code(), coder.HTTPStatus(), coder.String(), coder.Reference())
		}
	}
	
	fmt.Println("\n--- Example 5: No error ---")
	_ = simulateFetchResource("existent_id")
	_ = simulateExternalServiceCall("reliable_service")
	_ = processData("clean_data")

}

/*
Example Output (Stack traces will vary):

--- Example 1: NewWithCode (predefined Coder) ---
Error (%v): Not found: resource with ID 'nonexistent' was not found
Error (%+v):
resource with ID 'nonexistent' was not found
main.simulateFetchResource
	/path/to/your/file.go:22
main.main
	/path/to/your/file.go:70
runtime.main
	...
runtime.goexit
	...
Not found
  Code: 100002, HTTP Status: 404, Message: Not found, Reference: errors-spec.md#100002

--- Example 2: ErrorfWithCode (custom Coder) ---
Error (%v): Custom service is currently unavailable: failed to call unreliable_service due to intermittent issues
Error (%+v):
failed to call unreliable_service due to intermittent issues
main.simulateExternalServiceCall
	/path/to/your/file.go:34
main.main
	/path/to/your/file.go:80
runtime.main
	...
runtime.goexit
	...
Custom service is currently unavailable
  Code: 50001, HTTP Status: 503, Message: Custom service is currently unavailable, Reference: doc/link/to/service_unavailable.md

--- Example 3: WithCode (attaching Coder to existing error) ---
Error (%v): Bad request: underlying data corruption detected
Error (%+v):
underlying data corruption detected
main.processData
    /path/to/your/file.go:54
main.main
    /path/to/your/file.go:90
runtime.main
    ...
runtime.goexit
    ...
Bad request
  Code: 100003, HTTP Status: 400, Message: Bad request, Reference: errors-spec.md#100003

--- Example 4: WithCode (attaching Coder to an lmcc error) ---
Error (%v): Bad request: access to this data is forbidden by policy
Error (%+v):
access to this data is forbidden by policy
main.processData
	/path/to/your/file.go:54
main.main
	/path/to/your/file.go:100
runtime.main
	...
runtime.goexit
	...
Bad request
  Code: 100003, HTTP Status: 400, Message: Bad request, Reference: errors-spec.md#100003

--- Example 5: No error ---
Resource 'existent_id' fetched successfully.
External service 'reliable_service' called successfully.
Data 'clean_data' processed successfully.
*/ 