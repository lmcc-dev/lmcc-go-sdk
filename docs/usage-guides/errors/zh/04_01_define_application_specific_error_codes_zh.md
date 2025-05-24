<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 1. 定义特定于应用程序的错误码 (Define Application-Specific Error Codes)

对于应用程序可能需要以编程方式处理的错误，或表示特定业务逻辑失败的错误，请定义自定义 `Coder` 实例。这使得错误处理比依赖错误消息的字符串比较更加健壮。

(For errors that your application might need to handle programmatically or that represent specific business logic failures, define custom `Coder` instances. This makes error handling more robust than relying on string comparisons of error messages.)

- **集中定义 (Centralize Definitions)**: 将您的 `Coder` 实例定义为包级变量，通常放在相关包内一个专用的 `errors.go` 或 `codes.go` 文件中 (或一个共享的 `apperrors` 包中)。
  (Define your `Coder` instances as package-level variables, often in a dedicated `errors.go` or `codes.go` file within the relevant package (or a shared `apperrors` package).)
- **清晰命名 (Clear Naming)**: 为您的 `Coder` 变量使用描述性名称 (例如，`ErrUserNotFound`、`ErrPaymentGatewayTimeout`)。
  (Use descriptive names for your `Coder` variables (e.g., `ErrUserNotFound`, `ErrPaymentGatewayTimeout`).)
- **唯一代码 (Unique Codes)**: 确保特定于应用程序的错误码在您的应用程序域内是唯一的。
  (Ensure application-specific error codes are unique within your application domain.)
- **HTTP 状态 (HTTP Status)**: 如果您的错误可能通过 API 公开，请分配适当的 HTTP 状态码。
  (Assign appropriate HTTP status codes if your errors might be exposed via an API.)
- **参考文档 (Reference Docs)**: 可选地，提供一个指向有关该错误的更详细文档的参考 URL。
  (Optionally, provide a reference URL for more detailed documentation about the error.)

```go
package main

import (
	"fmt"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
)

// 良好的做法是在您的应用程序或模块的中心位置定义这些。
// (It's good practice to define these in a central place for your application or module.)

// ErrOrderNotFound 表示找不到具有指定 ID 的订单。
// (ErrOrderNotFound indicates that an order with the specified ID could not be found.)
var ErrOrderNotFound = errors.NewCoder(
	20001, // 特定于应用程序的代码 (Application-specific code)
	404,   // HTTP 状态未找到 (HTTP Status Not Found)
	"订单未找到 (Order not found)", // 默认消息 (Default message)
	"/docs/errors/order-processing#20001_zh", // 指向内部文档的链接 (Link to internal documentation)
)

// ErrInventoryUnavailable 表示商品缺货。
// (ErrInventoryUnavailable indicates that an item is out of stock.)
var ErrInventoryUnavailable = errors.NewCoder(
	20002, // 特定于应用程序的代码 (Application-specific code)
	409,   // HTTP 状态冲突 (HTTP Status Conflict) (因为请求由于当前状态而无法处理)
	"商品库存不可用 (Item inventory unavailable)",
	"/docs/errors/order-processing#20002_zh",
)

// ErrInvalidOrderData 表示提供的订单数据无效。
// (ErrInvalidOrderData indicates that the provided order data is invalid.)
var ErrInvalidOrderData = errors.NewCoder(
	20003, // 特定于应用程序的代码 (Application-specific code)
	400,   // HTTP 状态错误请求 (HTTP Status Bad Request)
	"提供的订单数据无效 (Invalid order data provided)",
	"/docs/errors/order-processing#20003_zh",
)

// processOrder 模拟处理订单。
// (processOrder simulates processing an order.)
func processOrder(orderID string, itemID string, quantity int) error {
	if orderID == "unknown_order" {
		return errors.ErrorfWithCode(ErrOrderNotFound, "无法检索到 ID 为 '%s' 的订单 (order with ID '%s' could not be retrieved)", orderID)
	}

	if itemID == "sold_out_item" && quantity > 0 {
		return errors.WrapfWithCode(ErrInventoryUnavailable, errors.New("库存水平为零 (stock level is zero)"), "无法完成商品 '%s' 的订单 (cannot fulfill order for item '%s')", itemID)
	}

	if quantity <= 0 {
		return errors.NewWithCode(ErrInvalidOrderData, "数量必须为正 (quantity must be positive)")
	}

	fmt.Printf("商品 '%s' 的订单 '%s' (数量 %d) 已成功处理。(Order '%s' for item '%s' (quantity %d) processed successfully.)\n", orderID, itemID, quantity, orderID, itemID, quantity) // Duplicate arguments removed for clarity
	return nil
}

func main() {
	fmt.Println("--- 场景1：订单未找到 --- (--- Scenario 1: Order Not Found ---)")
	err1 := processOrder("unknown_order", "item123", 1)
	if err1 != nil {
		fmt.Printf("错误 (Error): %v\n", err1)
		if errors.IsCode(err1, ErrOrderNotFound) {
			fmt.Println("程序化检查：这是一个 ErrOrderNotFound。(Programmatic check: This is an ErrOrderNotFound.)")
			coder := errors.GetCoder(err1)
			fmt.Printf("  Coder 详细信息 (Coder details): Code=%d, HTTPStatus=%d, Message='%s', Ref='%s'\n",
				coder.Code(), coder.HTTPStatus(), coder.String(), coder.Reference())
		}
	}

	fmt.Println("\n--- 场景2：库存不可用 --- (--- Scenario 2: Inventory Unavailable ---)")
	err2 := processOrder("order456", "sold_out_item", 5)
	if err2 != nil {
		fmt.Printf("错误 (Error): %v\n", err2)
		if errors.IsCode(err2, ErrInventoryUnavailable) {
			fmt.Println("程序化检查：这是一个 ErrInventoryUnavailable。(Programmatic check: This is an ErrInventoryUnavailable.)")
			coder := errors.GetCoder(err2)
			fmt.Printf("  Coder 详细信息 (Coder details): Code=%d, HTTPStatus=%d, Message='%s', Ref='%s'\n",
				coder.Code(), coder.HTTPStatus(), coder.String(), coder.Reference())
		}
		// 如果需要，检查原因的示例 (Example of checking the cause if needed)
		cause := errors.Cause(err2)
		fmt.Printf("  根本原因 (Underlying cause): %v\n", cause) // 注意：这里的原因是 WrapfWithCode 包装的那个 (Note: cause here is the one wrapped by WrapfWithCode)
	}

	fmt.Println("\n--- 场景3：订单数据无效 --- (--- Scenario 3: Invalid Order Data ---)")
	err3 := processOrder("order789", "item789", 0)
	if err3 != nil {
		fmt.Printf("错误 (Error): %v\n", err3)
		if errors.IsCode(err3, ErrInvalidOrderData) {
			fmt.Println("程序化检查：这是一个 ErrInvalidOrderData。(Programmatic check: This is an ErrInvalidOrderData.)")
		}
	}

	fmt.Println("\n--- 场景4：订单成功 --- (--- Scenario 4: Successful Order ---)")
	err4 := processOrder("order101", "itemABC", 2)
	if err4 == nil {
		fmt.Println("订单已成功处理。(Order processed successfully.)")
	}
}

/*
示例输出 (Example Output):

--- 场景1：订单未找到 --- (--- Scenario 1: Order Not Found ---)
错误 (Error): 订单未找到 (Order not found): 无法检索到 ID 为 'unknown_order' 的订单 (order with ID 'unknown_order' could not be retrieved)
程序化检查：这是一个 ErrOrderNotFound。(Programmatic check: This is an ErrOrderNotFound.)
  Coder 详细信息 (Coder details): Code=20001, HTTPStatus=404, Message='订单未找到 (Order not found)', Ref='/docs/errors/order-processing#20001_zh'

--- 场景2：库存不可用 --- (--- Scenario 2: Inventory Unavailable ---)
错误 (Error): 商品库存不可用 (Item inventory unavailable): 无法完成商品 'sold_out_item' 的订单 (cannot fulfill order for item 'sold_out_item'): 库存水平为零 (stock level is zero)
程序化检查：这是一个 ErrInventoryUnavailable。(Programmatic check: This is an ErrInventoryUnavailable.)
  Coder 详细信息 (Coder details): Code=20002, HTTPStatus=409, Message='商品库存不可用 (Item inventory unavailable)', Ref='/docs/errors/order-processing#20002_zh'
  根本原因 (Underlying cause): 库存水平为零 (stock level is zero)

--- 场景3：订单数据无效 --- (--- Scenario 3: Invalid Order Data ---)
错误 (Error): 提供的订单数据无效 (Invalid order data provided): 数量必须为正 (quantity must be positive)
程序化检查：这是一个 ErrInvalidOrderData。(Programmatic check: This is an ErrInvalidOrderData.)

--- 场景4：订单成功 --- (--- Scenario 4: Successful Order ---)
商品 'order101' 的订单 'itemABC' (数量 2) 已成功处理。(Order 'order101' for item 'itemABC' (quantity 2) processed successfully.)
订单已成功处理。(Order processed successfully.)
*/
``` 