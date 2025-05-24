<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

## 自定义错误码 (Custom Error Codes)

虽然 `pkg/errors` 为常见的错误场景提供了一组预定义的 `Coder` 实例 (例如 `ErrNotFound`、`ErrInternalServer` 等)，但您通常需要定义特定于应用程序的错误码。

(While `pkg/errors` provides a set of predefined `Coder` instances for common error scenarios (like `ErrNotFound`, `ErrInternalServer`, etc.), you will often need to define your own application-specific error codes.)

这可以通过创建您自己的 `errors.Coder` 接口实现来完成，或者更常见的是，通过使用 `errors.NewCoder` 工厂函数。

(This is done by creating your own implementations of the `errors.Coder` interface or, more commonly, by using the `errors.NewCoder` factory function.)

**`errors.NewCoder(code int, httpStatus int, message string, reference string) Coder`**

- **`code int`**: 您唯一的特定于应用程序的整数代码。建议系统地管理这些代码 (例如，在不同模块或服务的范围内)。
  (Your unique application-specific integer code. It's recommended to manage these codes systematically (e.g., within ranges for different modules or services).)
- **`httpStatus int`**: 如果通过 API 公开，此错误应映射到的相应 HTTP 状态码。
  (The corresponding HTTP status code that this error should map to if exposed via an API.)
- **`message string`**: 此错误代码的默认人类可读消息。
  (A default, human-readable message for this error code.)
- **`reference string`** (可选): 指向有关此错误的更详细文档的 URL 或路径，这对于 API 使用者或支持团队来说非常宝贵。
  ( (optional): A URL or path to more detailed documentation about this error, which can be invaluable for API consumers or support teams.)

### 示例：定义和使用自定义错误码 (Example: Defining and Using Custom Error Codes)

让我们为电子商务订单处理模块定义一些自定义错误码。

(Let's define some custom error codes for an e-commerce order processing module.)

```go
package main

import (
	"fmt"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
)

// 为订单处理系统定义自定义 Coder 实例
// (Define custom Coder instances for an order processing system)
var (
	// ErrOrderValidation 指示订单数据本身存在问题。
	// (ErrOrderValidation indicates a problem with the order data itself.)
	ErrOrderValidation = errors.NewCoder(
		20001, // 订单验证问题的特定于应用程序的代码 (Application-specific code for order validation issues)
		400,   // HTTP 错误请求 (HTTP Bad Request)
		"订单数据验证失败 (Order data validation failed)",
		"/docs/api/errors/order#20001_zh",
	)

	// ErrPaymentProcessing 指示付款期间发生故障。
	// (ErrPaymentProcessing indicates a failure during payment.)
	ErrPaymentProcessing = errors.NewCoder(
		20002, // 付款问题的特定于应用程序的代码 (Application-specific code for payment issues)
		402,   // HTTP 需要付款 (或 500，如果是内部网关错误) (HTTP Payment Required (or 500 if it's an internal gateway error))
		"付款处理失败 (Payment processing failed)",
		"/docs/api/errors/order#20002_zh",
	)

	// ErrInventoryCheck 指示商品缺货或数量不足。
	// (ErrInventoryCheck indicates an item is out of stock or insufficient.)
	ErrInventoryCheck = errors.NewCoder(
		20003, // 库存问题的特定于应用程序的代码 (Application-specific code for inventory issues)
		409,   // HTTP 冲突 (因为当前库存状态阻止履行) (HTTP Conflict (as the current state of inventory prevents fulfillment))
		"库存检查失败 (Inventory check failed)",
		"/docs/api/errors/order#20003_zh",
	)

	// ErrShippingUnavailable 指示向该地址发货时出现问题。
	// (ErrShippingUnavailable indicates a problem with shipping to the address.)
	ErrShippingUnavailable = errors.NewCoder(
		20004, // 发货问题的特定于应用程序的代码 (Application-specific code for shipping issues)
		400,   // HTTP 错误请求 (例如，无效地址) 或 503 (例如，没有可用的承运商) (HTTP Bad Request (e.g., invalid address) or 503 (e.g., no carriers available))
		"当前无法运送到提供的地址 (Shipping to the provided address is currently unavailable)",
		"/docs/api/errors/order#20004_zh",
	)
)

// Order 表示简化的订单结构
// (Order represents a simplified order structure)
type Order struct {
	OrderID    string
	UserID     string
	ItemID     string
	Quantity   int
	TotalPrice float64
	Address    string
	PaymentToken string
}

// processOrderStep 模拟可能失败的订单处理中的不同步骤。
// (processOrderStep simulates different steps in order processing that can fail.)
func processOrderStep(order *Order, step string) error {
	fmt.Printf("正在处理订单 '%s' 的步骤：%s\n", order.OrderID, step)
	// (Processing order '%s' at step: %s\n)
	switch step {
	case "validate_data":
		if order.Quantity <= 0 {
			// 使用 NewWithCode 创建带有我们自定义 Coder 的错误
			// (Use NewWithCode to create an error with our custom Coder)
			return errors.NewWithCode(ErrOrderValidation, "订单数量必须为正 (order quantity must be positive)")
		}
		if order.TotalPrice <= 0 {
			return errors.ErrorfWithCode(ErrOrderValidation, "订单总价 (%.2f) 必须为正 (order total price (%.2f) must be positive)", order.TotalPrice)
		}
	case "check_inventory":
		if order.ItemID == "ITEM_OUTOFSTOCK" {
			// 使用我们的自定义 Coder 和附加上下文包装底层 (可能是通用的) 错误
			// (Wrap an underlying (perhaps generic) error with our custom Coder and additional context)
			underlyingErr := errors.New("库存数据库报告商品数量为零 (stock database reports zero quantity for item)")
			return errors.WrapWithCode(underlyingErr, ErrInventoryCheck, fmt.Sprintf("商品 '%s' 已缺货 (item '%s' is out of stock)", order.ItemID))
		}
	case "process_payment":
		if order.PaymentToken == "TOKEN_DECLINED" {
			return errors.NewWithCode(ErrPaymentProcessing, "网关拒绝付款 (payment declined by gateway)")
		}
	case "arrange_shipping":
		if order.Address == "UNREACHABLE_LOCATION" {
			return errors.ErrorfWithCode(ErrShippingUnavailable, "无法运送到地址：%s (cannot ship to address: %s)", order.Address)
		}
	default:
		fmt.Printf("订单 '%s' 的步骤 '%s' 已成功完成。\n", step, order.OrderID)
		// (Step '%s' for order '%s' completed successfully.\n)
		return nil
	}
	fmt.Printf("订单 '%s' 的步骤 '%s' 已成功完成。\n", step, order.OrderID)
	// (Step '%s' for order '%s' completed successfully.\n)
	return nil
}

func main() {
	orders := []Order{
		{OrderID: "ORD001", ItemID: "ITEM_XYZ", Quantity: 0, TotalPrice: 10.00, Address: "123 Main St", PaymentToken: "TOKEN_VALID"}, // 数量无效 (Invalid quantity)
		{OrderID: "ORD002", ItemID: "ITEM_OUTOFSTOCK", Quantity: 1, TotalPrice: 25.00, Address: "456 Oak Ave", PaymentToken: "TOKEN_VALID"}, // 缺货 (Out of stock)
		{OrderID: "ORD003", ItemID: "ITEM_ABC", Quantity: 2, TotalPrice: 50.00, Address: "789 Pine Ln", PaymentToken: "TOKEN_DECLINED"}, // 付款被拒 (Payment declined)
		{OrderID: "ORD004", ItemID: "ITEM_DEF", Quantity: 1, TotalPrice: 30.00, Address: "UNREACHABLE_LOCATION", PaymentToken: "TOKEN_VALID"}, // 无法运送 (Shipping unavailable)
		{OrderID: "ORD005", ItemID: "ITEM_GHI", Quantity: 3, TotalPrice: 75.00, Address: "321 Elm Rd", PaymentToken: "TOKEN_VALID"}, // 成功订单 (Successful order)
	}

	steps := []string{"validate_data", "check_inventory", "process_payment", "arrange_shipping", "finalize"}

	for _, order := range orders {
		fmt.Printf("\n--- 正在处理订单：%s ---\n", order.OrderID)
		// (--- Processing Order: %s ---\n)
		var finalError error
		for _, step := range steps {
			err := processOrderStep(&order, step)
			if err != nil {
				fmt.Printf("  步骤 '%s' 出错：%v\n", step, err)
				// (  Error at step '%s': %v\n)
				// 记录带有堆栈跟踪的详细错误以进行内部调试
				// (Log detailed error with stack trace for internal debugging)
				// fmt.Printf("    内部日志：\n%+v\n", err)
				// (    Internal log: \n%+v\n)
				
				// 检查特定的自定义代码
				// (Check for specific custom codes)
				if errors.IsCode(err, ErrOrderValidation) {
					fmt.Printf("    [检查] 这是一个订单验证错误。HTTP 状态：%d。参考：%s\n",
						ErrOrderValidation.HTTPStatus(), ErrOrderValidation.Reference())
					// (    [Check] This is an Order Validation error. HTTP Status: %d. Ref: %s\n)
				}
				if errors.IsCode(err, ErrInventoryCheck) {
					fmt.Printf("    [检查] 这是一个库存检查错误。HTTP 状态：%d。参考：%s\n",
						ErrInventoryCheck.HTTPStatus(), ErrInventoryCheck.Reference())
					// (    [Check] This is an Inventory Check error. HTTP Status: %d. Ref: %s\n)
				}
				finalError = err // 保留工作流中遇到的第一个错误 (Keep the first error encountered in the workflow)
				break // 发生第一个错误时停止处理此订单 (Stop processing this order on first error)
			}
		}
		if finalError == nil {
			fmt.Println("  订单已成功处理！(Order processed successfully!)")
		}
	}
}

/*
示例输出 (此摘要中省略了 %+v 中的堆栈跟踪以保持简洁)：
(Example Output (Stack traces are omitted from %+v for brevity in this summary)):

--- 正在处理订单：ORD001 ---
(--- Processing Order: ORD001 ---)
正在处理订单 'ORD001' 的步骤：validate_data
(Processing order 'ORD001' at step: validate_data)
  步骤 'validate_data' 出错：订单数据验证失败 (Order data validation failed): 订单数量必须为正 (order quantity must be positive)
  (Error at step 'validate_data': Order data validation failed: order quantity must be positive)
    [检查] 这是一个订单验证错误。HTTP 状态：400。参考：/docs/api/errors/order#20001_zh
    ([Check] This is an Order Validation error. HTTP Status: 400. Ref: /docs/api/errors/order#20001_zh)

--- 正在处理订单：ORD002 ---
(--- Processing Order: ORD002 ---)
正在处理订单 'ORD002' 的步骤：validate_data
(Processing order 'ORD002' at step: validate_data)
订单 'ORD002' 的步骤 'validate_data' 已成功完成。
(Step 'validate_data' for order 'ORD002' completed successfully.)
正在处理订单 'ORD002' 的步骤：check_inventory
(Processing order 'ORD002' at step: check_inventory)
  步骤 'check_inventory' 出错：库存检查失败 (Inventory check failed): 商品 'ITEM_OUTOFSTOCK' 已缺货 (item 'ITEM_OUTOFSTOCK' is out of stock): 库存数据库报告商品数量为零 (stock database reports zero quantity for item)
  (Error at step 'check_inventory': Inventory check failed: item 'ITEM_OUTOFSTOCK' is out of stock: stock database reports zero quantity for item)
    [检查] 这是一个库存检查错误。HTTP 状态：409。参考：/docs/api/errors/order#20003_zh
    ([Check] This is an Inventory Check error. HTTP Status: 409. Ref: /docs/api/errors/order#20003_zh)

--- 正在处理订单：ORD003 ---
(--- Processing Order: ORD003 ---)
正在处理订单 'ORD003' 的步骤：validate_data
(Processing order 'ORD003' at step: validate_data)
订单 'ORD003' 的步骤 'validate_data' 已成功完成。
(Step 'validate_data' for order 'ORD003' completed successfully.)
正在处理订单 'ORD003' 的步骤：check_inventory
(Processing order 'ORD003' at step: check_inventory)
订单 'ORD003' 的步骤 'check_inventory' 已成功完成。
(Step 'check_inventory' for order 'ORD003' completed successfully.)
正在处理订单 'ORD003' 的步骤：process_payment
(Processing order 'ORD003' at step: process_payment)
  步骤 'process_payment' 出错：付款处理失败 (Payment processing failed): 网关拒绝付款 (payment declined by gateway)
  (Error at step 'process_payment': Payment processing failed: payment declined by gateway)

--- 正在处理订单：ORD004 ---
(--- Processing Order: ORD004 ---)
正在处理订单 'ORD004' 的步骤：validate_data
(Processing order 'ORD004' at step: validate_data)
订单 'ORD004' 的步骤 'validate_data' 已成功完成。
(Step 'validate_data' for order 'ORD004' completed successfully.)
正在处理订单 'ORD004' 的步骤：check_inventory
(Processing order 'ORD004' at step: check_inventory)
订单 'ORD004' 的步骤 'check_inventory' 已成功完成。
(Step 'check_inventory' for order 'ORD004' completed successfully.)
正在处理订单 'ORD004' 的步骤：process_payment
(Processing order 'ORD004' at step: process_payment)
订单 'ORD004' 的步骤 'process_payment' 已成功完成。
(Step 'process_payment' for order 'ORD004' completed successfully.)
正在处理订单 'ORD004' 的步骤：arrange_shipping
(Processing order 'ORD004' at step: arrange_shipping)
  步骤 'arrange_shipping' 出错：当前无法运送到提供的地址 (Shipping to the provided address is currently unavailable): 无法运送到地址：UNREACHABLE_LOCATION (cannot ship to address: UNREACHABLE_LOCATION)
  (Error at step 'arrange_shipping': Shipping to the provided address is currently unavailable: cannot ship to address: UNREACHABLE_LOCATION)

--- 正在处理订单：ORD005 ---
(--- Processing Order: ORD005 ---)
正在处理订单 'ORD005' 的步骤：validate_data
(Processing order 'ORD005' at step: validate_data)
订单 'ORD005' 的步骤 'validate_data' 已成功完成。
(Step 'validate_data' for order 'ORD005' completed successfully.)
正在处理订单 'ORD005' 的步骤：check_inventory
(Processing order 'ORD005' at step: check_inventory)
订单 'ORD005' 的步骤 'check_inventory' 已成功完成。
(Step 'check_inventory' for order 'ORD005' completed successfully.)
正在处理订单 'ORD005' 的步骤：process_payment
(Processing order 'ORD005' at step: process_payment)
订单 'ORD005' 的步骤 'process_payment' 已成功完成。
(Step 'process_payment' for order 'ORD005' completed successfully.)
正在处理订单 'ORD005' 的步骤：arrange_shipping
(Processing order 'ORD005' at step: arrange_shipping)
订单 'ORD005' 的步骤 'arrange_shipping' 已成功完成。
(Step 'arrange_shipping' for order 'ORD005' completed successfully.)
正在处理订单 'ORD005' 的步骤：finalize
(Processing order 'ORD005' at step: finalize)
订单 'ORD005' 的步骤 'finalize' 已成功完成。
(Step 'finalize' for order 'ORD005' completed successfully.)
  订单已成功处理！(Order processed successfully!)
*/
```

通过定义和使用自定义错误码，您可以使应用程序的错误处理更具可预测性、更易于测试，并且对开发人员和 API 使用者都更具信息性。
(By defining and using custom error codes, you make your application's error handling more predictable, easier to test, and more informative for both developers and API consumers.) 