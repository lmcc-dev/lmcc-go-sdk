<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 4. Use `errors.IsCode` for Checking Specific Error Categories

When you need to check if an error belongs to a category defined by a `Coder` (e.g., any kind of "not found" error, regardless of the specific message), use `errors.IsCode(err, YourCoder)`.

- `errors.IsCode` checks the `Code()` of the `Coder` instances, not the `Coder` variable identity. This means if you have multiple `Coder` variables that happen to share the same integer code, `IsCode` would consider them matching for that code.
- This is generally more reliable than `standardErrors.Is(err, YourCoderVariable)` if `YourCoderVariable` might be wrapped or if you want to check solely based on the numeric code.
- However, standard `errors.Is(err, errors.ErrNotFound)` (where `errors.ErrNotFound` is a predefined `Coder` from the `pkg/errors` itself) works well because these predefined `Coder` variables act like sentinel values.

```go
package main

import (
	"fmt"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	standardLibErrors "errors" // Alias for standard library errors
)

// Predefined Coders (imagine these are from different parts of a larger system)
var ErrLegacyNotFound = errors.NewCoder(40400, 404, "Resource not found (legacy system)", "")
var ErrAPIV1NotFound = errors.NewCoder(40400, 404, "Endpoint not found (API v1)", "") // Same code as ErrLegacyNotFound
var ErrResourceGone = errors.NewCoder(41000, 410, "Resource permanently gone", "")

// Function that might return one of these errors
func fetchResource(resourceID string) error {
	switch resourceID {
	case "legacy_item_123":
		return errors.WrapWithCode(standardLibErrors.New("original db error: no rows"), ErrLegacyNotFound, "accessing legacy data store")
	case "api_v1_user_789":
		// Directly return an error with the Coder
		return errors.NewWithCode(ErrAPIV1NotFound, "user endpoint /users/789 not available")
	case "deleted_doc_456":
		return errors.NewWithCode(ErrResourceGone, "document 456 was deleted")
	case "other_issue":
		return errors.New("some other unrelated problem") // Error without a Coder
	default:
		fmt.Printf("Resource '%s' found successfully.\n", resourceID)
		return nil
	}
}

func main() {
	testCases := []struct {
		name           string
		resourceID     string
		expectIsCodeMatchWithLegacyNotFound bool
		expectIsMatchWithLegacyNotFound   bool
		expectIsMatchWithAPIV1NotFound bool
	}{
		{
			name: "Error with ErrLegacyNotFound",
			resourceID: "legacy_item_123",
			expectIsCodeMatchWithLegacyNotFound: true, // errors.IsCode checks the code number (40400)
			expectIsMatchWithLegacyNotFound:   true, // errors.Is checks for the specific Coder instance in the chain
			expectIsMatchWithAPIV1NotFound: false, // errors.Is would not match ErrAPIV1NotFound instance
		},
		{
			name: "Error with ErrAPIV1NotFound",
			resourceID: "api_v1_user_789",
			expectIsCodeMatchWithLegacyNotFound: true, // errors.IsCode matches because ErrAPIV1NotFound also has code 40400
			expectIsMatchWithLegacyNotFound:   false, // errors.Is would not match ErrLegacyNotFound instance
			expectIsMatchWithAPIV1NotFound: true,
		},
		{
			name: "Error with ErrResourceGone",
			resourceID: "deleted_doc_456",
			expectIsCodeMatchWithLegacyNotFound: false, // Different code (41000)
			expectIsMatchWithLegacyNotFound:   false,
			expectIsMatchWithAPIV1NotFound: false,
		},
		{
			name: "Uncoded error",
			resourceID: "other_issue",
			expectIsCodeMatchWithLegacyNotFound: false,
			expectIsMatchWithLegacyNotFound:   false,
			expectIsMatchWithAPIV1NotFound: false,
		},
		{
			name: "No error",
			resourceID: "existing_resource",
		},
	}

	for _, tc := range testCases {
		fmt.Printf("--- Test Case: %s (Resource ID: %s) ---\n", tc.name, tc.resourceID)
		err := fetchResource(tc.resourceID)

		if err == nil {
			fmt.Println("No error occurred.")
			continue
		}
		fmt.Printf("Received error: %v\n", err)

		// Using errors.IsCode (checks based on the Coder's Code() value)
		isCodeLegacy := errors.IsCode(err, ErrLegacyNotFound)
		fmt.Printf("  errors.IsCode(err, ErrLegacyNotFound (code %d)): %t. Expected: %t\n", ErrLegacyNotFound.Code(), isCodeLegacy, tc.expectIsCodeMatchWithLegacyNotFound)

		// Using standardErrors.Is (checks for specific Coder instance in the chain)
		isLegacy := standardLibErrors.Is(err, ErrLegacyNotFound)
		fmt.Printf("  standardErrors.Is(err, ErrLegacyNotFound): %t. Expected: %t\n", isLegacy, tc.expectIsMatchWithLegacyNotFound)

		isAPIV1 := standardLibErrors.Is(err, ErrAPIV1NotFound)
		fmt.Printf("  standardErrors.Is(err, ErrAPIV1NotFound): %t. Expected: %t\n", isAPIV1, tc.expectIsMatchWithAPIV1NotFound)

		coder := errors.GetCoder(err)
		if coder != nil {
			fmt.Printf("  Actual Coder in error: Code=%d, Message='%s'\n", coder.Code(), coder.String())
		} else {
			fmt.Println("  No Coder found in error.")
		}
	}
}

/*
Example Output:

--- Test Case: Error with ErrLegacyNotFound (Resource ID: legacy_item_123) ---
Received error: accessing legacy data store: Resource not found (legacy system): original db error: no rows
  errors.IsCode(err, ErrLegacyNotFound (code 40400)): true. Expected: true
  standardErrors.Is(err, ErrLegacyNotFound): true. Expected: true
  standardErrors.Is(err, ErrAPIV1NotFound): false. Expected: false
  Actual Coder in error: Code=40400, Message='Resource not found (legacy system)'

--- Test Case: Error with ErrAPIV1NotFound (Resource ID: api_v1_user_789) ---
Received error: Endpoint not found (API v1): user endpoint /users/789 not available
  errors.IsCode(err, ErrLegacyNotFound (code 40400)): true. Expected: true
  standardErrors.Is(err, ErrLegacyNotFound): false. Expected: false
  standardErrors.Is(err, ErrAPIV1NotFound): true. Expected: true
  Actual Coder in error: Code=40400, Message='Endpoint not found (API v1)'

--- Test Case: Error with ErrResourceGone (Resource ID: deleted_doc_456) ---
Received error: Resource permanently gone: document 456 was deleted
  errors.IsCode(err, ErrLegacyNotFound (code 40400)): false. Expected: false
  standardErrors.Is(err, ErrLegacyNotFound): false. Expected: false
  standardErrors.Is(err, ErrAPIV1NotFound): false. Expected: false
  Actual Coder in error: Code=41000, Message='Resource permanently gone'

--- Test Case: Uncoded error (Resource ID: other_issue) ---
Received error: some other unrelated problem
  errors.IsCode(err, ErrLegacyNotFound (code 40400)): false. Expected: false
  standardErrors.Is(err, ErrLegacyNotFound): false. Expected: false
  standardErrors.Is(err, ErrAPIV1NotFound): false. Expected: false
  No Coder found in error.

--- Test Case: No error (Resource ID: existing_resource) ---
Resource 'existing_resource' found successfully.
No error occurred.
*/
```

**When to use `errors.IsCode` vs `standardErrors.Is`:**
- Use `errors.IsCode(err, SpecificCoder)` when you want to check if the error is of a certain *category* identified by the `SpecificCoder.Code()`, regardless of the exact `Coder` instance.
- Use `standardErrors.Is(err, TargetErrorOrCoderInstance)` when checking for a specific sentinel error value or a specific `Coder` variable instance in the error chain. This is useful for predefined `Coder` variables from `pkg/errors` like `errors.ErrNotFound`, `errors.ErrInternalServer`, etc., which act as sentinels. It's also appropriate for your own package-level `Coder` variables if you want to check for that *exact* instance. 