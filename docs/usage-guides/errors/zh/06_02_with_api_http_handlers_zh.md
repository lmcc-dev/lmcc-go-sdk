<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 2. 与 API/HTTP 处理程序一起使用 (With API/HTTP Handlers)

在 HTTP 处理程序中，服务层返回的错误可以转换为适当的 HTTP 响应。

(In HTTP handlers, errors returned by service layers can be translated into appropriate HTTP responses.)

- **HTTP 状态码 (HTTP Status Codes)**: 使用 `errors.GetCoder(err).HTTPStatus()` 获取相关的 HTTP 状态码。如果未找到 `Coder` 或其没有特定的 HTTP 状态，则默认为 500 内部服务器错误或其他适当的状态。
  (Use `errors.GetCoder(err).HTTPStatus()` to get a relevant HTTP status code. If no `Coder` is found or it doesn't have a specific HTTP status, default to 500 Internal Server Error or another suitable status.)
- **错误响应体 (Error Response Body)**: `Coder` 的 `Code()` (特定于应用程序的错误码) 和 `String()` (消息) 可用于为客户端构建结构化的 JSON 错误响应。
  (The `Coder`'s `Code()` (application-specific error code) and `String()` (message) can be used to construct a structured JSON error response for the client.)
- **日志记录 (Logging)**: 在服务器端记录完整的错误详细信息 (`%+v`) 以进行调试，但仅在 HTTP 响应中返回客户端安全的信息。
  (Log the full error details (`%+v`) on the server side for debugging, but only return client-safe information in the HTTP response.)

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	pkgErrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
)

// ErrorResponse 是 JSON 错误响应的标准结构。
// (ErrorResponse is a standard structure for JSON error responses.)
type ErrorResponse struct {
	ErrorCode    int    `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage"`
	RequestID    string `json:"requestID,omitempty"` // 示例：包含请求 ID 以进行跟踪 (Example: include a request ID for tracing)
	DocsURL      string `json:"docsURL,omitempty"`   // 指向错误文档的链接 (Link to documentation for the error)
}

var ErrProductLogic = pkgErrors.NewCoder(90001, 500, "产品业务逻辑失败 (Product business logic failed)", "/docs/api/products#90001_zh")
var ErrPaymentRequired = pkgErrors.NewCoder(90002, 402, "访问此高级功能需要付款 (Payment is required to access this premium feature)", "/docs/api/billing#payment-required_zh")

// productService 模拟服务层函数。
// (productService simulates a service layer function.)
func productService(productID string, premiumFeature bool) (string, error) {
	if productID == "nonexistent" {
		return "", pkgErrors.ErrorfWithCode(pkgErrors.ErrNotFound, "ID 为 '%s' 的产品不存在 (product with ID '%s' does not exist)", productID)
	}
	if productID == "faulty" {
		originalErr := pkgErrors.New("底层库存检查失败 (underlying inventory check failed)")
		return "", pkgErrors.WrapWithCode(originalErr, ErrProductLogic, "由于内部错误，处理产品 '%s' 失败 (failed to process product '%s' due to internal error)", productID)
	}
	if premiumFeature && productID != "premium_user_product" { // 简化的付费检查 (Simplified premium check)
		return "", pkgErrors.NewWithCode(ErrPaymentRequired, "拒绝访问产品 '%s' (access denied for product '%s')", productID)
	}
	return fmt.Sprintf("产品 %s 的数据 (Data for product %s)", productID, productID), nil
}

// productHandler 是一个示例 HTTP 处理程序。
// (productHandler is an example HTTP handler.)
func productHandler(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("id")
	premium := r.URL.Query().Get("premium") == "true"
	requestID := r.Header.Get("X-Request-ID") // 示例：从标头获取请求 ID (Example: get a request ID from header)
	if requestID == "" {
		requestID = "unknown"
	}

	log.Printf("[服务器] 请求 %s：正在处理产品 ID：%s，高级功能：%t 的产品请求 ([Server] Request %s: Handling product request for ID: %s, Premium: %t)", requestID, productID, premium, requestID, productID, premium)

	data, err := productService(productID, premium)
	if err != nil {
		log.Printf("[服务器] 请求 %s：产品 '%s' 的服务错误：
%+v
 ([Server] Request %s: Service error for product '%s':\n%+v\n)", requestID, productID, err, requestID, productID, err) // 服务器端详细日志 (Server-side detailed log)

		apiErrorCode := 0 // 如果没有 Coder，则默认为 0 (Default if no Coder)
		httpStatus := http.StatusInternalServerError
		clientMessage := "发生内部服务器错误。(An internal server error occurred.)"
		docsURL := ""

		if coder := pkgErrors.GetCoder(err); coder != nil {
			apiErrorCode = coder.Code()
			httpStatus = coder.HTTPStatus()
			clientMessage = coder.String() // 使用 Coder 的消息给客户端 (Use Coder's message for client)
			docsURL = coder.Reference()
			// 如果您也想要包装的消息，您可能希望使用 err.Error() 覆盖 clientMessage。
			// (You might want to override clientMessage with err.Error() if you want the wrapped messages too.)
			// 对于此示例，coder.String() 提供了代码的基本消息。
			// (For this example, coder.String() provides the base message of the code.)
			// clientMessage = err.Error() // 取消注释以发送完整的包装消息 (Uncomment to send full wrapped message)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpStatus)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrorCode:    apiErrorCode,
			ErrorMessage: clientMessage,
			RequestID:    requestID,
			DocsURL:      docsURL,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"data": data, "requestID": requestID})
}

func main() {
	tmux := http.NewServeMux()
	tmux.HandleFunc("/product", productHandler)

	testServer := httptest.NewServer(tmux)
	defer testServer.Close()

	scenarios := []struct {
		name           string
		path           string
		expectedStatus int
		expectedErrorCode int // 如果 JSON 中没有预期的特定应用程序错误代码，则为 0 (0 if no specific app error code expected in JSON)
		expectedMessageContains string
	}{
		{"产品未找到 (Product Not Found)", "/product?id=nonexistent", http.StatusNotFound, pkgErrors.ErrNotFound.Code(), pkgErrors.ErrNotFound.String()},
		{"内部服务器错误 (Internal Server Error)", "/product?id=faulty", http.StatusInternalServerError, ErrProductLogic.Code(), ErrProductLogic.String()},
		{"需要付款 (Payment Required)", "/product?id=some_product&premium=true", http.StatusPaymentRequired, ErrPaymentRequired.Code(), ErrPaymentRequired.String()},
		{"成功请求 (Successful Request)", "/product?id=prod123", http.StatusOK, 0, "产品 prod123 的数据 (Data for product prod123)"},
	}

	client := testServer.Client()

	for _, s := range scenarios {
		fmt.Printf("\n--- 测试场景：%s ---\n请求：%s%s\n", s.name, testServer.URL, s.path)
		// (--- Test Scenario: %s ---\nRequesting: %s%s\n)
		req, _ := http.NewRequest("GET", testServer.URL+s.path, nil)
		req.Header.Set("X-Request-ID", "test-"+s.name)
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("为场景 '%s' 发出请求失败：%v (Failed to make request for scenario '%s': %v)", s.name, err, s.name, err)
		}

		fmt.Printf("响应状态：%s\n", resp.Status)
		// (Response Status: %s\n)
		var body map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&body)
		resp.Body.Close()
		fmt.Printf("响应体：%v\n", body)
		// (Response Body: %v\n)

		if resp.StatusCode != s.expectedStatus {
			fmt.Printf("  [失败] 预期状态 %d，得到 %d\n", s.expectedStatus, resp.StatusCode)
			// (  [FAIL] Expected status %d, got %d\n)
		}
		if s.expectedErrorCode != 0 {
			iferrorCode, ok := body["errorCode"].(float64); !ok || int(iferrorCode) != s.expectedErrorCode {
				fmt.Printf("  [失败] 预期 errorCode %d，得到 %v\n", s.expectedErrorCode, body["errorCode"])
				// (  [FAIL] Expected errorCode %d, got %v\n)
			}
		}
		// 为简洁起见，进行简单的消息检查
		// (Simple message check for brevity)
		// 对于错误，我们检查 errorMessage 字段。对于成功，我们检查 data 字段。
		// (For errors, we check errorMessage field. For success, we check data field.)
		messageToCheck := ""
		if val, ok := body["errorMessage"].(string); ok {
		    messageToCheck = val
		} else if val, ok := body["data"].(string); ok {
		    messageToCheck = val
		}
        if s.expectedMessageContains != "" && (messageToCheck == "" || !contains(messageToCheck, s.expectedMessageContains)) {
            fmt.Printf("  [失败] 预期消息包含 '%s'，得到 '%s'\n", s.expectedMessageContains, messageToCheck)
            // (  [FAIL] Expected message to contain '%s', got '%s'\n)
        }
	}
}

// contains 是用于字符串检查的辅助函数。
// (contains is a helper for string checking.)
func contains(s, substr string) bool {
    // 在 Go 中，使用 strings.Contains(s, substr)
    // (In Go, use strings.Contains(s, substr))
	// 这是一个简化的实现，以避免导入 "strings" 包
    // (This is a simplified implementation to avoid importing "strings" package)
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}


/*
示例控制台输出 (服务器的日志行是交错的，JSON 结构是关键)：
(Example Console Output (Log lines from server are interleaved, JSON structure is key)):

--- 测试场景：产品未找到 (Product Not Found) ---
请求：[URL]/product?id=nonexistent
(Requesting: [URL]/product?id=nonexistent)
响应状态：404 Not Found
(Response Status: 404 Not Found)
响应体：map[docsURL:errors-spec.md#100002_zh errorCode:100002 errorMessage:Not found requestID:test-产品未找到 (Product Not Found)]
(Response Body: map[docsURL:errors-spec.md#100002_zh errorCode:100002 errorMessage:Not found requestID:test-Product Not Found])

--- 测试场景：内部服务器错误 (Internal Server Error) ---
请求：[URL]/product?id=faulty
(Requesting: [URL]/product?id=faulty)
响应状态：500 Internal Server Error
(Response Status: 500 Internal Server Error)
响应体：map[docsURL:/docs/api/products#90001_zh errorCode:90001 errorMessage:产品业务逻辑失败 (Product business logic failed) requestID:test-内部服务器错误 (Internal Server Error)]
(Response Body: map[docsURL:/docs/api/products#90001_zh errorCode:90001 errorMessage:Product business logic failed requestID:test-Internal Server Error])

--- 测试场景：需要付款 (Payment Required) ---
请求：[URL]/product?id=some_product&premium=true
(Requesting: [URL]/product?id=some_product&premium=true)
响应状态：402 Payment Required
(Response Status: 402 Payment Required)
响应体：map[docsURL:/docs/api/billing#payment-required_zh errorCode:90002 errorMessage:访问此高级功能需要付款 (Payment is required to access this premium feature) requestID:test-需要付款 (Payment Required)]
(Response Body: map[docsURL:/docs/api/billing#payment-required_zh errorCode:90002 errorMessage:Payment is required to access this premium feature requestID:test-Payment Required])

--- 测试场景：成功请求 (Successful Request) ---
请求：[URL]/product?id=prod123
(Requesting: [URL]/product?id=prod123)
响应状态：200 OK
(Response Status: 200 OK)
响应体：map[data:产品 prod123 的数据 (Data for product prod123) requestID:test-成功请求 (Successful Request)]
(Response Body: map[data:Data for product prod123 requestID:test-Successful Request])

(服务器日志行，如 "[服务器] 请求 test-产品未找到 (Product Not Found)：正在处理产品请求..." 也会被打印)
((Server log lines like "[Server] Request test-Product Not Found: Handling product request..." will also be printed))
*/
``` 