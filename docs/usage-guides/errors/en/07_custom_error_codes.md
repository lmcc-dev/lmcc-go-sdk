<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

## Custom Error Codes

While `pkg/errors` provides a set of predefined `Coder` instances for common error scenarios (like `ErrNotFound`, `ErrInternalServer`, etc.), you will often need to define your own application-specific error codes.

This is done by creating your own implementations of the `errors.Coder` interface or, more commonly, by using the `errors.NewCoder` factory function.

**`errors.NewCoder(code int, httpStatus int, message string, reference string) Coder`**

- **`code int`**: Your unique application-specific integer code. It's recommended to manage these codes systematically (e.g., within ranges for different modules or services).
- **`httpStatus int`**: The corresponding HTTP status code that this error should map to if exposed via an API.
- **`message string`**: A default, human-readable message for this error code.
- **`reference string`** (optional): A URL or path to more detailed documentation about this error, which can be invaluable for API consumers or support teams.

### Example: Defining and Using Custom Error Codes

Let's define some custom error codes for an e-commerce order processing module.

```go
package main

import (
	"fmt"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
)

// Define custom Coder instances for an order processing system
var (
	// ErrOrderValidation indicates a problem with the order data itself.
	ErrOrderValidation = errors.NewCoder(
		20001, // Application-specific code for order validation issues
		400,   // HTTP Bad Request
		"Order data validation failed",
		"/docs/api/errors/order#20001",
	)

	// ErrPaymentProcessing indicates a failure during payment.
	ErrPaymentProcessing = errors.NewCoder(
		20002, // Application-specific code for payment issues
		402,   // HTTP Payment Required (or 500 if it's an internal gateway error)
		"Payment processing failed",
		"/docs/api/errors/order#20002",
	)

	// ErrInventoryCheck indicates an item is out of stock or insufficient.
	ErrInventoryCheck = errors.NewCoder(
		20003, // Application-specific code for inventory issues
		409,   // HTTP Conflict (as the current state of inventory prevents fulfillment)
		"Inventory check failed",
		"/docs/api/errors/order#20003",
	)

	// ErrShippingUnavailable indicates a problem with shipping to the address.
	ErrShippingUnavailable = errors.NewCoder(
		20004, // Application-specific code for shipping issues
		400,   // HTTP Bad Request (e.g., invalid address) or 503 (e.g., no carriers available)
		"Shipping to the provided address is currently unavailable",
		"/docs/api/errors/order#20004",
	)
)

// Order represents a simplified order structure
type Order struct {
	OrderID    string
	UserID     string
	ItemID     string
	Quantity   int
	TotalPrice float64
	Address    string
	PaymentToken string
}

// processOrderStep simulates different steps in order processing that can fail.
func processOrderStep(order *Order, step string) error {
	fmt.Printf("Processing order '%s' at step: %s\n", order.OrderID, step)
	switch step {
	case "validate_data":
		if order.Quantity <= 0 {
			// Use NewWithCode to create an error with our custom Coder
			return errors.NewWithCode(ErrOrderValidation, "order quantity must be positive")
		}
		if order.TotalPrice <= 0 {
			return errors.ErrorfWithCode(ErrOrderValidation, "order total price (%.2f) must be positive", order.TotalPrice)
		}
	case "check_inventory":
		if order.ItemID == "ITEM_OUTOFSTOCK" {
			// Wrap an underlying (perhaps generic) error with our custom Coder and additional context
			underlyingErr := errors.New("stock database reports zero quantity for item")
			return errors.WrapWithCode(underlyingErr, ErrInventoryCheck, fmt.Sprintf("item '%s' is out of stock", order.ItemID))
		}
	case "process_payment":
		if order.PaymentToken == "TOKEN_DECLINED" {
			return errors.NewWithCode(ErrPaymentProcessing, "payment declined by gateway")
		}
	case "arrange_shipping":
		if order.Address == "UNREACHABLE_LOCATION" {
			return errors.ErrorfWithCode(ErrShippingUnavailable, "cannot ship to address: %s", order.Address)
		}
	default:
		fmt.Printf("Step '%s' for order '%s' completed successfully.\n", step, order.OrderID)
		return nil
	}
	fmt.Printf("Step '%s' for order '%s' completed successfully.\n", step, order.OrderID)
	return nil
}

func main() {
	orders := []Order{
		{OrderID: "ORD001", ItemID: "ITEM_XYZ", Quantity: 0, TotalPrice: 10.00, Address: "123 Main St", PaymentToken: "TOKEN_VALID"}, // Invalid quantity
		{OrderID: "ORD002", ItemID: "ITEM_OUTOFSTOCK", Quantity: 1, TotalPrice: 25.00, Address: "456 Oak Ave", PaymentToken: "TOKEN_VALID"}, // Out of stock
		{OrderID: "ORD003", ItemID: "ITEM_ABC", Quantity: 2, TotalPrice: 50.00, Address: "789 Pine Ln", PaymentToken: "TOKEN_DECLINED"}, // Payment declined
		{OrderID: "ORD004", ItemID: "ITEM_DEF", Quantity: 1, TotalPrice: 30.00, Address: "UNREACHABLE_LOCATION", PaymentToken: "TOKEN_VALID"}, // Shipping unavailable
		{OrderID: "ORD005", ItemID: "ITEM_GHI", Quantity: 3, TotalPrice: 75.00, Address: "321 Elm Rd", PaymentToken: "TOKEN_VALID"}, // Successful order
	}

	steps := []string{"validate_data", "check_inventory", "process_payment", "arrange_shipping", "finalize"}

	for _, order := range orders {
		fmt.Printf("\n--- Processing Order: %s ---\n", order.OrderID)
		var finalError error
		for _, step := range steps {
			err := processOrderStep(&order, step)
			if err != nil {
				fmt.Printf("  Error at step '%s': %v\n", step, err)
				// Log detailed error with stack trace for internal debugging
				// fmt.Printf("    Internal log: \n%+v\n", err)
				
				// Check for specific custom codes
				if errors.IsCode(err, ErrOrderValidation) {
					fmt.Printf("    [Check] This is an Order Validation error. HTTP Status: %d. Ref: %s\n",
						ErrOrderValidation.HTTPStatus(), ErrOrderValidation.Reference())
				}
				if errors.IsCode(err, ErrInventoryCheck) {
					fmt.Printf("    [Check] This is an Inventory Check error. HTTP Status: %d. Ref: %s\n",
						ErrInventoryCheck.HTTPStatus(), ErrInventoryCheck.Reference())
				}
				finalError = err // Keep the first error encountered in the workflow
				break // Stop processing this order on first error
			}
		}
		if finalError == nil {
			fmt.Println("  Order processed successfully!")
		}
	}
}

/*
Example Output (Stack traces are omitted from %+v for brevity in this summary):

--- Processing Order: ORD001 ---
Processing order 'ORD001' at step: validate_data
  Error at step 'validate_data': Order data validation failed: order quantity must be positive
    [Check] This is an Order Validation error. HTTP Status: 400. Ref: /docs/api/errors/order#20001

--- Processing Order: ORD002 ---
Processing order 'ORD002' at step: validate_data
Step 'validate_data' for order 'ORD002' completed successfully.
Processing order 'ORD002' at step: check_inventory
  Error at step 'check_inventory': Inventory check failed: item 'ITEM_OUTOFSTOCK' is out of stock: stock database reports zero quantity for item
    [Check] This is an Inventory Check error. HTTP Status: 409. Ref: /docs/api/errors/order#20003

--- Processing Order: ORD003 ---
Processing order 'ORD003' at step: validate_data
Step 'validate_data' for order 'ORD003' completed successfully.
Processing order 'ORD003' at step: check_inventory
Step 'check_inventory' for order 'ORD003' completed successfully.
Processing order 'ORD003' at step: process_payment
  Error at step 'process_payment': Payment processing failed: payment declined by gateway

--- Processing Order: ORD004 ---
Processing order 'ORD004' at step: validate_data
Step 'validate_data' for order 'ORD004' completed successfully.
Processing order 'ORD004' at step: check_inventory
Step 'check_inventory' for order 'ORD004' completed successfully.
Processing order 'ORD004' at step: process_payment
Step 'process_payment' for order 'ORD004' completed successfully.
Processing order 'ORD004' at step: arrange_shipping
  Error at step 'arrange_shipping': Shipping to the provided address is currently unavailable: cannot ship to address: UNREACHABLE_LOCATION

--- Processing Order: ORD005 ---
Processing order 'ORD005' at step: validate_data
Step 'validate_data' for order 'ORD005' completed successfully.
Processing order 'ORD005' at step: check_inventory
Step 'check_inventory' for order 'ORD005' completed successfully.
Processing order 'ORD005' at step: process_payment
Step 'process_payment' for order 'ORD005' completed successfully.
Processing order 'ORD005' at step: arrange_shipping
Step 'arrange_shipping' for order 'ORD005' completed successfully.
Processing order 'ORD005' at step: finalize
Step 'finalize' for order 'ORD005' completed successfully.
  Order processed successfully!
*/
```

By defining and using custom error codes, you make your application's error handling more predictable, easier to test, and more informative for both developers and API consumers.
 