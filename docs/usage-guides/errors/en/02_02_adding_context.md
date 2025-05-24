<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 2. Adding Context

Often, an error from a lower-level function needs to be augmented with context from the calling function. `Wrap` and `Wrapf` are used for this. They preserve the original error (and its stack trace) while adding a new message and a new stack trace frame for the wrapping location.

**`errors.Wrap(err error, message string) error`**: Wraps `err` with a new message. If `err` is nil, `Wrap` returns nil.

**`errors.Wrapf(err error, format string, args ...interface{}) error`**: Wraps `err` with a new formatted message. If `err` is nil, `Wrapf` returns nil.

```go
package main

import (
	"fmt"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	standardErrors "errors" // For creating a sample original error
)

// simulateDatabaseConnection
// This function simulates a low-level operation that might fail.
func simulateDatabaseConnection(dsn string) error {
	if dsn == "" {
		return standardErrors.New("database DSN is empty") // Using standard error for variety
	}
	if dsn == "invalid_dsn" {
		return errors.New("failed to connect to database: invalid credentials") // Using lmcc-go-sdk error
	}
	// Simulate successful connection
	fmt.Printf("Successfully connected to database with DSN: %s\n", dsn)
	return nil
}

// initializeUserService
// This function calls the lower-level connection function and wraps its error if it occurs.
func initializeUserService(dsnConfig string) error {
	dbErr := simulateDatabaseConnection(dsnConfig)
	if dbErr != nil {
		// Wrap the error from simulateDatabaseConnection with more context.
		return errors.Wrap(dbErr, "user service initialization failed")
	}
	fmt.Println("User service initialized successfully.")
	return nil
}

// configureApplication
// This function calls initializeUserService and uses Wrapf for more dynamic context.
func configureApplication(serviceName string, dsnForService string) error {
	err := initializeUserService(dsnForService)
	if err != nil {
		// Wrap the error from initializeUserService with a formatted message.
		return errors.Wrapf(err, "failed to configure application component: %s", serviceName)
	}
	fmt.Printf("Application component '%s' configured successfully.\n", serviceName)
	return nil
}

func main() {
	// Scenario 1: Error originates, wrapped by Wrap
	fmt.Println("--- Scenario 1: Wrapping a standard error ---")
	err1 := initializeUserService("invalid_dsn") // This DSN will cause simulateDatabaseConnection to return an lmcc error
	if err1 != nil {
		fmt.Printf("Error (%%v):  %v\n", err1)
		fmt.Printf("Error (%%+v):\n%+v\n", err1)
	}

	fmt.Println("\n--- Scenario 2: Error originates, wrapped by Wrap, then Wrapf ---")
	// Scenario 2: Error originates, wrapped by Wrap, then wrapped again by Wrapf
	err2 := configureApplication("AuthService", "") // Empty DSN will cause simulateDatabaseConnection to return a standard error
	if err2 != nil {
		fmt.Printf("Error (%%v):  %v\n", err2)
		fmt.Printf("Error (%%+v):\n%+v\n", err2)
	}

	fmt.Println("\n--- Scenario 3: No error ---")
	// Scenario 3: No error occurs
	err3 := configureApplication("PaymentService", "valid_dsn_string")
	if err3 != nil { // This block should not be executed
		fmt.Printf("Unexpected error: %+v\n", err3)
	}
}

/*
Example Output (Stack traces will vary):

--- Scenario 1: Wrapping a standard error ---
Error (%v):  user service initialization failed: failed to connect to database: invalid credentials
Error (%+v):
failed to connect to database: invalid credentials
main.simulateDatabaseConnection
	/path/to/your/file.go:15
main.initializeUserService
	/path/to/your/file.go:29
main.main
	/path/to/your/file.go:50
runtime.main
	/usr/local/go/src/runtime/proc.go:250
runtime.goexit
	/usr/local/go/src/runtime/asm_amd64.s:1197
user service initialization failed
main.initializeUserService
	/path/to/your/file.go:30
main.main
	/path/to/your/file.go:50
runtime.main
	/usr/local/go/src/runtime/proc.go:250
runtime.goexit
	/usr/local/go/src/runtime/asm_amd64.s:1197

--- Scenario 2: Error originates, wrapped by Wrap, then Wrapf ---
Error (%v):  failed to configure application component: AuthService: user service initialization failed: database DSN is empty
Error (%+v):
database DSN is empty
main.simulateDatabaseConnection
	/path/to/your/file.go:12
main.initializeUserService
	/path/to/your/file.go:29
main.configureApplication
	/path/to/your/file.go:41
main.main
	/path/to/your/file.go:59
runtime.main
	/usr/local/go/src/runtime/proc.go:250
runtime.goexit
	/usr/local/go/src/runtime/asm_amd64.s:1197
user service initialization failed
main.initializeUserService
	/path/to/your/file.go:30
main.configureApplication
	/path/to/your/file.go:41
main.main
	/path/to/your/file.go:59
runtime.main
	/usr/local/go/src/runtime/proc.go:250
runtime.goexit
	/usr/local/go/src/runtime/asm_amd64.s:1197
failed to configure application component: AuthService
main.configureApplication
	/path/to/your/file.go:42
main.main
	/path/to/your/file.go:59
runtime.main
	/usr/local/go/src/runtime/proc.go:250
runtime.goexit
	/usr/local/go/src/runtime/asm_amd64.s:1197

--- Scenario 3: No error ---
Successfully connected to database with DSN: valid_dsn_string
User service initialized successfully.
Application component 'PaymentService' configured successfully.
*/
```
Notice how `Wrap` and `Wrapf` add layers to the error. When printed with `%+v`, all layers of messages and their respective stack traces are shown, providing a rich history of the error's propagation.
