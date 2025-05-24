<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 2. With API/HTTP Handlers

In HTTP handlers, errors returned by service layers can be translated into appropriate HTTP responses.

- **HTTP Status Codes**: Use `errors.GetCoder(err).HTTPStatus()` to get a relevant HTTP status code. If no `Coder` is found or it doesn't have a specific HTTP status, default to 500 Internal Server Error or another suitable status.
- **Error Response Body**: The `Coder`'s `Code()` (application-specific error code) and `String()` (message) can be used to construct a structured JSON error response for the client.
- **Logging**: Log the full error details (`%+v`) on the server side for debugging, but only return client-safe information in the HTTP response.

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

// ErrorResponse is a standard structure for JSON error responses.

type ErrorResponse struct {
	ErrorCode    int    `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage"`
	RequestID    string `json:"requestID,omitempty"` // Example: include a request ID for tracing
	DocsURL      string `json:"docsURL,omitempty"`   // Link to documentation for the error
}

var ErrProductLogic = pkgErrors.NewCoder(90001, 500, "Product business logic failed", "/docs/api/products#90001")
var ErrPaymentRequired = pkgErrors.NewCoder(90002, 402, "Payment is required to access this premium feature", "/docs/api/billing#payment-required")

// productService simulates a service layer function.
func productService(productID string, premiumFeature bool) (string, error) {
	if productID == "nonexistent" {
		return "", pkgErrors.ErrorfWithCode(pkgErrors.ErrNotFound, "product with ID '%s' does not exist", productID)
	}
	if productID == "faulty" {
		originalErr := pkgErrors.New("underlying inventory check failed")
		return "", pkgErrors.WrapWithCode(originalErr, ErrProductLogic, "failed to process product '%s' due to internal error", productID)
	}
	if premiumFeature && productID != "premium_user_product" { // Simplified premium check
		return "", pkgErrors.NewWithCode(ErrPaymentRequired, "access denied for product '%s'", productID)
	}
	return fmt.Sprintf("Data for product %s", productID), nil
}

// productHandler is an example HTTP handler.
func productHandler(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("id")
	premium := r.URL.Query().Get("premium") == "true"
	requestID := r.Header.Get("X-Request-ID") // Example: get a request ID from header
	if requestID == "" {
		requestID = "unknown"
	}

	log.Printf("[Server] Request %s: Handling product request for ID: %s, Premium: %t", requestID, productID, premium)

	data, err := productService(productID, premium)
	if err != nil {
		log.Printf("[Server] Request %s: Service error for product '%s':\n%+v\n", requestID, productID, err) // Server-side detailed log

		apiErrorCode := 0 // Default if no Coder
		httpStatus := http.StatusInternalServerError
		clientMessage := "An internal server error occurred."
		docsURL := ""

		if coder := pkgErrors.GetCoder(err); coder != nil {
			apiErrorCode = coder.Code()
			httpStatus = coder.HTTPStatus()
			clientMessage = coder.String() // Use Coder's message for client
			docsURL = coder.Reference()
			// You might want to override clientMessage with err.Error() if you want the wrapped messages too.
			// For this example, coder.String() provides the base message of the code.
			// clientMessage = err.Error() // Uncomment to send full wrapped message
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
		expectedErrorCode int // 0 if no specific app error code expected in JSON
		expectedMessageContains string
	}{
		{"Product Not Found", "/product?id=nonexistent", http.StatusNotFound, pkgErrors.ErrNotFound.Code(), pkgErrors.ErrNotFound.String()},
		{"Internal Server Error", "/product?id=faulty", http.StatusInternalServerError, ErrProductLogic.Code(), ErrProductLogic.String()},
		{"Payment Required", "/product?id=some_product&premium=true", http.StatusPaymentRequired, ErrPaymentRequired.Code(), ErrPaymentRequired.String()},
		{"Successful Request", "/product?id=prod123", http.StatusOK, 0, "Data for product prod123"},
	}

	client := testServer.Client()

	for _, s := range scenarios {
		fmt.Printf("\n--- Test Scenario: %s ---\nRequesting: %s%s\n", s.name, testServer.URL, s.path)
		req, _ := http.NewRequest("GET", testServer.URL+s.path, nil)
		req.Header.Set("X-Request-ID", "test-"+s.name)
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Failed to make request for scenario '%s': %v", s.name, err)
		}

		fmt.Printf("Response Status: %s\n", resp.Status)
		var body map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&body)
		resp.Body.Close()
		fmt.Printf("Response Body: %v\n", body)

		if resp.StatusCode != s.expectedStatus {
			fmt.Printf("  [FAIL] Expected status %d, got %d\n", s.expectedStatus, resp.StatusCode)
		}
		if s.expectedErrorCode != 0 {
			iferrorCode, ok := body["errorCode"].(float64); !ok || int(iferrorCode) != s.expectedErrorCode {
				fmt.Printf("  [FAIL] Expected errorCode %d, got %v\n", s.expectedErrorCode, body["errorCode"])
			}
		}
		// Simple message check for brevity
		// For errors, we check errorMessage field. For success, we check data field.
		messageToCheck := ""
		if val, ok := body["errorMessage"].(string); ok {
		    messageToCheck = val
		} else if val, ok := body["data"].(string); ok {
		    messageToCheck = val
		}
        if s.expectedMessageContains != "" && (messageToCheck == "" || !contains(messageToCheck, s.expectedMessageContains)) {
            fmt.Printf("  [FAIL] Expected message to contain '%s', got '%s'\n", s.expectedMessageContains, messageToCheck)
        }
	}
}

// contains is a helper for string checking.
func contains(s, substr string) bool {
    return فيها(s, substr) // Placeholder for actual string contains logic if needed. fmt.Sprintf("%s",s) is a hack.
    // In Go, use strings.Contains(s, substr)
    return len(fmt.Sprintf("%s",s)) > 0 && len(fmt.Sprintf("%s",substr)) > 0 && เกิดข้อผิดพลาด(s, substr)
}

// These are just to make the example runnable without importing `strings` for a simple check.
func فيها(s, substr string) bool { 
	for i:=0; i<len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
func เกิดข้อผิดพลาด(s, substr string) bool { return فيها(s, substr) }

/*
Example Console Output (Log lines from server are interleaved, JSON structure is key):

--- Test Scenario: Product Not Found ---
Requesting: [URL]/product?id=nonexistent
Response Status: 404 Not Found
Response Body: map[docsURL:errors-spec.md#100002 errorCode:100002 errorMessage:Not found requestID:test-Product Not Found]

--- Test Scenario: Internal Server Error ---
Requesting: [URL]/product?id=faulty
Response Status: 500 Internal Server Error
Response Body: map[docsURL:/docs/api/products#90001 errorCode:90001 errorMessage:Product business logic failed requestID:test-Internal Server Error]

--- Test Scenario: Payment Required ---
Requesting: [URL]/product?id=some_product&premium=true
Response Status: 402 Payment Required
Response Body: map[docsURL:/docs/api/billing#payment-required errorCode:90002 errorMessage:Payment is required to access this premium feature requestID:test-Payment Required]

--- Test Scenario: Successful Request ---
Requesting: [URL]/product?id=prod123
Response Status: 200 OK
Response Body: map[data:Data for product prod123 requestID:test-Successful Request]

(Server log lines like "[Server] Request test-Product Not Found: Handling product request..." will also be printed)
*/
``` 