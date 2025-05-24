<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### Error Aggregation

Collect multiple errors that occur during a single broader operation (e.g., validating multiple fields of a form) into a single `ErrorGroup`.

**`errors.NewErrorGroup(message string) *ErrorGroup`**: Creates a new error group with an initial message.

An `ErrorGroup` has the following methods:
- **`Add(err error)`**: Adds an error to the group. Does nothing if `err` is nil.
- **`Errors() []error`**: Returns a slice of all errors added to the group. Returns `nil` if no errors were added.
- **`Error() string`**: Returns a string representation of all errors in the group, prefixed by the group's initial message.
- **`Unwrap() []error`**: (Go 1.20+) Allows `ErrorGroup` to be used with `standardErrors.Is` and `standardErrors.As` by exposing the collected errors.
- **`Format(s fmt.State, verb rune)`**: Implements `fmt.Formatter` for custom formatting (e.g., `%+v` for detailed output including individual error stack traces).

```go
package main

import (
	"fmt"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"strings"
)

// UserProfile
// A sample struct we want to validate.
type UserProfile struct {
	Username string
	Email    string
	Age      int
}

// validateUserProfile
// This function validates a UserProfile and collects all validation errors into an ErrorGroup.
func validateUserProfile(profile UserProfile) error {
	// Create a new error group with a general message for this validation operation.
	eg := errors.NewErrorGroup("User profile validation failed")

	// Validate Username
	if strings.TrimSpace(profile.Username) == "" {
		// Add an error for missing username. We can use errors.New or errors.NewWithCode.
		eg.Add(errors.New("username is required"))
	}

	// Validate Email
	if strings.TrimSpace(profile.Email) == "" {
		eg.Add(errors.NewWithCode(errors.ErrValidation, "email cannot be empty"))
	} else if !strings.Contains(profile.Email, "@") {
		eg.Add(errors.ErrorfWithCode(errors.ErrValidation, "email '%s' is not a valid format", profile.Email))
	}

	// Validate Age
	if profile.Age < 0 {
		eg.Add(errors.New("age cannot be negative"))
	} else if profile.Age < 18 {
		eg.Add(errors.Errorf("user must be at least 18 years old, got %d", profile.Age))
	}

	// Check if any errors were added to the group.
	if len(eg.Errors()) > 0 {
		return eg // Return the ErrorGroup itself, which is an error.
	}

	// If no errors were added, validation passed.
	return nil
}

func main() {
	// Scenario 1: Profile with multiple validation issues
	profile1 := UserProfile{Username: "", Email: "invalid-email", Age: 15}
	fmt.Println("--- Validating Profile 1 (multiple issues) ---")
	err1 := validateUserProfile(profile1)
	if err1 != nil {
		fmt.Printf("Validation Error (%%v):\n%v\n\n", err1)
		fmt.Printf("Validation Error (%%+v) - Detailed:\n%+v\n", err1)
		
			// You can iterate over individual errors if needed
	if eg, ok := err1.(*errors.ErrorGroup); ok {
		fmt.Println("\nIndividual errors in the group:")
		for i, individualErr := range eg.Errors() {
			fmt.Printf("  %d: %v\n", i+1, individualErr)
		}
		
		// Check if the ErrorGroup contains errors with specific codes
		fmt.Println("\nChecking for specific error codes in the group:")
		if errors.IsCode(eg, errors.ErrValidation) {
			fmt.Println("  - Contains validation errors")
		}
		if errors.IsCode(eg, errors.ErrNotFound) {
			fmt.Println("  - Contains not found errors")
		} else {
			fmt.Println("  - No not found errors")
		}
	}
	}

	fmt.Println("\n--- Validating Profile 2 (one issue) ---")
	// Scenario 2: Profile with a single validation issue
	profile2 := UserProfile{Username: "ValidUser", Email: "valid@example.com", Age: -5}
	err2 := validateUserProfile(profile2)
	if err2 != nil {
		fmt.Printf("Validation Error (%%v):\n%v\n\n", err2)
		// fmt.Printf("Validation Error (%%+v):\n%+v\n", err2) // %+v would also show stack traces for each error
	}

	fmt.Println("\n--- Validating Profile 3 (valid profile) ---")
	// Scenario 3: Valid profile, no errors
	profile3 := UserProfile{Username: "TestUser", Email: "test@example.com", Age: 30}
	err3 := validateUserProfile(profile3)
	if err3 == nil {
		fmt.Println("Profile 3 is valid!")
	} else {
		fmt.Printf("Unexpected validation error for Profile 3: %v\n", err3)
	}
}

/*
Example Output (Stack traces in %+v are omitted for brevity here but would be present):

--- Validating Profile 1 (multiple issues) ---
Validation Error (%v):
User profile validation failed: [username is required; email 'invalid-email' is not a valid format; user must be at least 18 years old, got 15]

Validation Error (%+v) - Detailed:
User profile validation failed: [username is required; email 'invalid-email' is not a valid format; user must be at least 18 years old, got 15]
Individual errors with stack traces:
1. username is required
   main.validateUserProfile
   	/path/to/your/file.go:XX
   ...
2. Validation failed: email 'invalid-email' is not a valid format
   main.validateUserProfile
   	/path/to/your/file.go:XX
   ...
3. user must be at least 18 years old, got 15
   main.validateUserProfile
   	/path/to/your/file.go:XX
   ...

Individual errors in the group:
  1: username is required
  2: Validation failed: email 'invalid-email' is not a valid format
  3: user must be at least 18 years old, got 15

--- Validating Profile 2 (one issue) ---
Validation Error (%v):
User profile validation failed: [age cannot be negative]

--- Validating Profile 3 (valid profile) ---
Profile 3 is valid!
*/
``` 