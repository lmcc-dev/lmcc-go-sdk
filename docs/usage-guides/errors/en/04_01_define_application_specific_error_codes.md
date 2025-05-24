<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 1. Define Application-Specific Error Codes

For errors that your application might need to handle programmatically or that represent specific business logic failures, define custom `Coder` instances. This makes error handling more robust than relying on string comparisons of error messages.

- **Centralize Definitions**: Define your `Coder` instances as package-level variables, often in a dedicated `errors.go` or `codes.go` file within the relevant package (or a shared `apperrors` package).
- **Clear Naming**: Use descriptive names for your `Coder` variables (e.g., `ErrUserNotFound`, `ErrPaymentGatewayTimeout`).
- **Unique Codes**: Ensure application-specific error codes are unique within your application domain.
- **HTTP Status**: Assign appropriate HTTP status codes if your errors might be exposed via an API.
- **Reference Docs**: Optionally, provide a reference URL for more detailed documentation about the error.

```go
package main

import (
	"fmt"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
)

// It's good practice to define these in a central place for your application or module.

// ErrOrderNotFound indicates that an order with the specified ID could not be found.
var ErrOrderNotFound = errors.NewCoder(
	20001, // Application-specific code
	404,   // HTTP Status Not Found
	"Order not found", // Default message
	"/docs/errors/order-processing#20001", // Link to internal documentation
)

// ErrInventoryUnavailable indicates that an item is out of stock.
var ErrInventoryUnavailable = errors.NewCoder(
	20002, // Application-specific code
	409,   // HTTP Status Conflict (as the request cannot be processed due to current state)
	"Item inventory unavailable",
	"/docs/errors/order-processing#20002",
)

// ErrInvalidOrderData indicates that the provided order data is invalid.
var ErrInvalidOrderData = errors.NewCoder(
	20003, // Application-specific code
	400,   // HTTP Status Bad Request
	"Invalid order data provided",
	"/docs/errors/order-processing#20003",
)

// processOrder simulates processing an order.
func processOrder(orderID string, itemID string, quantity int) error {
	if orderID == "unknown_order" {
		return errors.ErrorfWithCode(ErrOrderNotFound, "order with ID '%s' could not be retrieved", orderID)
	}

	if itemID == "sold_out_item" && quantity > 0 {
		return errors.WrapfWithCode(ErrInventoryUnavailable, errors.New("stock level is zero"), "cannot fulfill order for item '%s'", itemID)
	}

	if quantity <= 0 {
		return errors.NewWithCode(ErrInvalidOrderData, "quantity must be positive")
	}

	fmt.Printf("Order '%s' for item '%s' (quantity %d) processed successfully.\n", orderID, itemID, quantity)
	return nil
}

func main() {
	fmt.Println("--- Scenario 1: Order Not Found ---")
	err1 := processOrder("unknown_order", "item123", 1)
	if err1 != nil {
		fmt.Printf("Error: %v\n", err1)
		if errors.IsCode(err1, ErrOrderNotFound) {
			fmt.Println("Programmatic check: This is an ErrOrderNotFound.")
			coder := errors.GetCoder(err1)
			fmt.Printf("  Coder details: Code=%d, HTTPStatus=%d, Message='%s', Ref='%s'\n",
				coder.Code(), coder.HTTPStatus(), coder.String(), coder.Reference())
		}
	}

	fmt.Println("\n--- Scenario 2: Inventory Unavailable ---")
	err2 := processOrder("order456", "sold_out_item", 5)
	if err2 != nil {
		fmt.Printf("Error: %v\n", err2)
		if errors.IsCode(err2, ErrInventoryUnavailable) {
			fmt.Println("Programmatic check: This is an ErrInventoryUnavailable.")
			coder := errors.GetCoder(err2)
			fmt.Printf("  Coder details: Code=%d, HTTPStatus=%d, Message='%s', Ref='%s'\n",
				coder.Code(), coder.HTTPStatus(), coder.String(), coder.Reference())
		}
		// Example of checking the cause if needed
		cause := errors.Cause(err2)
		fmt.Printf("  Underlying cause: %v\n", cause) // Note: cause here is the one wrapped by WrapfWithCode
	}

	fmt.Println("\n--- Scenario 3: Invalid Order Data ---")
	err3 := processOrder("order789", "item789", 0)
	if err3 != nil {
		fmt.Printf("Error: %v\n", err3)
		if errors.IsCode(err3, ErrInvalidOrderData) {
			fmt.Println("Programmatic check: This is an ErrInvalidOrderData.")
		}
	}

	fmt.Println("\n--- Scenario 4: Successful Order ---")
	err4 := processOrder("order101", "itemABC", 2)
	if err4 == nil {
		fmt.Println("Order processed successfully.")
	}
}

/*
Example Output:

--- Scenario 1: Order Not Found ---
Error: Order not found: order with ID 'unknown_order' could not be retrieved
Programmatic check: This is an ErrOrderNotFound.
  Coder details: Code=20001, HTTPStatus=404, Message='Order not found', Ref='/docs/errors/order-processing#20001'

--- Scenario 2: Inventory Unavailable ---
Error: Item inventory unavailable: cannot fulfill order for item 'sold_out_item': stock level is zero
Programmatic check: This is an ErrInventoryUnavailable.
  Coder details: Code=20002, HTTPStatus=409, Message='Item inventory unavailable', Ref='/docs/errors/order-processing#20002'
  Underlying cause: stock level is zero

--- Scenario 3: Invalid Order Data ---
Error: Invalid order data provided: quantity must be positive
Programmatic check: This is an ErrInvalidOrderData.

--- Scenario 4: Successful Order ---
Order 'order101' for item 'itemABC' (quantity 2) processed successfully.
Order processed successfully.
*/
``` 