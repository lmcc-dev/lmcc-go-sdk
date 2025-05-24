<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 3. Handle Errors Appropriately at Different Layers

Not all errors should be handled in the same way. The action taken should depend on the layer of the application and the nature of the error.

- **Low-Level Functions**: Often, low-level functions (e.g., database access, file I/O) should simply return errors (possibly wrapped) to their callers. They lack the context to make high-level decisions (e.g., retry, abort, respond to user).
- **Service/Business Logic Layer**: This layer might handle certain errors by retrying operations, falling back to alternative strategies, or by wrapping errors with business-specific `Coder` instances before returning them.
- **Top-Level Handlers (e.g., HTTP handlers, main function)**: This is typically where errors are logged, and responses are formulated for the user or client system. Use `errors.GetCoder` to extract codes for generating appropriate HTTP status codes or user messages.

```go
package main

import (
	"fmt"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"time"
)

// --- Layer 1: Low-Level Data Access ---

var ErrDBConnectionFailed = errors.NewCoder(30001, 500, "Database connection failed", "")
var ErrDBRecordNotFound = errors.NewCoder(30002, 404, "Database record not found", "")

// fetchUserDataFromDB simulates fetching user data from a database.
func fetchUserDataFromDB(userID string) (string, error) {
	fmt.Printf("[DB Access] Attempting to fetch user '%s'\n", userID)
	if userID == "user_transient_error" {
		// Simulate a temporary DB issue
		return "", errors.WrapWithCode(errors.New("network timeout connecting to db replica"), ErrDBConnectionFailed, "fetchUserDataFromDB")
	}
	if userID == "user_not_found" {
		return "", errors.ErrorfWithCode(ErrDBRecordNotFound, "no record for user ID '%s'", userID)
	}
	if userID == "user_critical_db_issue" {
		// A more severe, non-coded error from DB layer initially
		return "", errors.New("unrecoverable checksum error on user table")
	}
	fmt.Printf("[DB Access] Successfully fetched user data for '%s'\n", userID)
	return fmt.Sprintf("Data for %s", userID), nil
}

// --- Layer 2: Service/Business Logic ---

var ErrUserServiceTemporarilyUnavailable = errors.NewCoder(40001, 503, "User service temporarily unavailable", "docs/retry-policy.md")

// getUserProfile retrieves user data and performs some business logic.
func getUserProfile(userID string) (string, error) {
	fmt.Printf("[Service Logic] Getting profile for user '%s'\n", userID)
	userData, err := fetchUserDataFromDB(userID)
	if err != nil {
		// Handle DB connection errors specifically: maybe retry or return a service unavailable error.
		if errors.IsCode(err, ErrDBConnectionFailed) {
			// For this example, we'll wrap it with a service-level "temporarily unavailable" Coder.
			// In a real app, you might implement retries here.
			return "", errors.WrapWithCode(err, ErrUserServiceTemporarilyUnavailable, "user service could not connect to database")
		}
		// For other errors (like ErrDBRecordNotFound or the critical DB issue),
		// we wrap them to add service-level context but let their original/DB-level Coder (if any) propagate.
		return "", errors.Wrapf(err, "failed to get user profile for '%s'", userID)
	}

	// Business logic: e.g., enrich user data (simplified)
	profile := fmt.Sprintf("Profile: %s [Processed by Service]", userData)
	fmt.Printf("[Service Logic] Successfully processed profile for '%s'\n", userID)
	return profile, nil
}

// --- Layer 3: Top-Level Handler (e.g., HTTP Handler / Main Function) ---

func handleAPIRequest(userID string) {
	fmt.Printf("[API Handler] Received request for user '%s'\n", userID)
	profile, err := getUserProfile(userID)

	if err != nil {
		fmt.Printf("[API Handler] Error processing request: %v\n", err) // Log the user-facing message
		// For detailed logging, one might use %+v with a structured logger in a real app
		// log.Error("Failed to handle API request", "error", fmt.Sprintf("%+v", err), "userID", userID)

		coder := errors.GetCoder(err) // Get the most specific Coder in the chain
		if coder != nil {
			switch coder.Code() {
			case ErrUserServiceTemporarilyUnavailable.Code():
				fmt.Printf("[API Handler] Responding with HTTP %d: %s (Reference: %s)\n", coder.HTTPStatus(), coder.String(), coder.Reference())
				// Respond to client with 503 Service Unavailable
			case ErrDBRecordNotFound.Code(): // This Coder was from DB layer but propagated up
				fmt.Printf("[API Handler] Responding with HTTP %d: %s\n", coder.HTTPStatus(), coder.String())
				// Respond to client with 404 Not Found
			default:
				// For other coded errors, use their HTTP status or a generic 500
				fmt.Printf("[API Handler] Responding with HTTP %d (default for code %d): %s\n", coder.HTTPStatus(), coder.Code(), coder.String())
			}
		} else {
			// Uncoded error, respond with a generic 500 Internal Server Error
			fmt.Printf("[API Handler] Responding with HTTP 500: An unexpected internal error occurred.\n")
			// Log the detailed error internally: fmt.Printf("%+v\n", err)
		}
		return
	}

	fmt.Printf("[API Handler] Successfully processed request. Response: %s\n", profile)
	// Respond to client with 200 OK and profile
}

func main() {
	scenarios := []string{"user_transient_error", "user_not_found", "user_critical_db_issue", "valid_user"}

	for _, userID := range scenarios {
		fmt.Printf("\n--- Simulating API request for user: '%s' ---\n", userID)
		handleAPIRequest(userID)
		time.Sleep(10 * time.Millisecond) // Small delay for readability of output
	}
}

/*
Example Output (simplified, actual stack traces omitted from this summary):

--- Simulating API request for user: 'user_transient_error' ---
[API Handler] Received request for user 'user_transient_error'
[Service Logic] Getting profile for user 'user_transient_error'
[DB Access] Attempting to fetch user 'user_transient_error'
[API Handler] Error processing request: user service could not connect to database: fetchUserDataFromDB: network timeout connecting to db replica
[API Handler] Responding with HTTP 503: User service temporarily unavailable (Reference: docs/retry-policy.md)

--- Simulating API request for user: 'user_not_found' ---
[API Handler] Received request for user 'user_not_found'
[Service Logic] Getting profile for user 'user_not_found'
[DB Access] Attempting to fetch user 'user_not_found'
[API Handler] Error processing request: failed to get user profile for 'user_not_found': no record for user ID 'user_not_found'
[API Handler] Responding with HTTP 404: Database record not found

--- Simulating API request for user: 'user_critical_db_issue' ---
[API Handler] Received request for user 'user_critical_db_issue'
[Service Logic] Getting profile for user 'user_critical_db_issue'
[DB Access] Attempting to fetch user 'user_critical_db_issue'
[API Handler] Error processing request: failed to get user profile for 'user_critical_db_issue': unrecoverable checksum error on user table
[API Handler] Responding with HTTP 500: An unexpected internal error occurred.

--- Simulating API request for user: 'valid_user' ---
[API Handler] Received request for user 'valid_user'
[Service Logic] Getting profile for user 'valid_user'
[DB Access] Attempting to fetch user 'valid_user'
[DB Access] Successfully fetched user data for 'valid_user'
[Service Logic] Successfully processed profile for 'valid_user'
[API Handler] Successfully processed request. Response: Profile: Data for valid_user [Processed by Service]
*/
``` 