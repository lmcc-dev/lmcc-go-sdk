/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Gin上下文适配器单元测试 (Gin context adapter unit tests)
 */

package gin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// setupTestContext 创建测试用的Gin上下文 (Create test Gin context)
func setupTestContext(method, path string, body string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	
	w := httptest.NewRecorder()
	var req *http.Request
	
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	
	return c, w
}

// TestNewGinContext 测试创建Gin上下文适配器 (Test creating Gin context adapter)
func TestNewGinContext(t *testing.T) {
	ginCtx, _ := setupTestContext("GET", "/test", "")
	
	ctx := NewGinContext(ginCtx)
	
	if ctx == nil {
		t.Fatal("NewGinContext() returned nil")
	}
	
	// 类型断言获取GinContext (Type assertion to get GinContext)
	ginContextAdapter, ok := ctx.(*GinContext)
	if !ok {
		t.Fatal("NewGinContext() should return *GinContext")
	}
	
	if ginContextAdapter.GetGinContext() != ginCtx {
		t.Error("GetGinContext() should return the original gin context")
	}
}

// TestGinContextRequest 测试获取请求对象 (Test getting request object)
func TestGinContextRequest(t *testing.T) {
	ginCtx, _ := setupTestContext("GET", "/test", "")
	ctx := NewGinContext(ginCtx)
	
	req := ctx.Request()
	if req == nil {
		t.Fatal("Request() returned nil")
	}
	
	if req.Method != "GET" {
		t.Errorf("Expected method GET, got %s", req.Method)
	}
	
	if req.URL.Path != "/test" {
		t.Errorf("Expected path /test, got %s", req.URL.Path)
	}
}

// TestGinContextResponse 测试获取响应写入器 (Test getting response writer)
func TestGinContextResponse(t *testing.T) {
	ginCtx, w := setupTestContext("GET", "/test", "")
	ctx := NewGinContext(ginCtx)
	
	resp := ctx.Response()
	if resp == nil {
		t.Fatal("Response() returned nil")
	}
	
	// Gin的Writer是一个包装器，不会直接等于ResponseRecorder
	// 我们检查它是否实现了http.ResponseWriter接口 (Gin's Writer is a wrapper, won't directly equal ResponseRecorder. We check if it implements http.ResponseWriter interface)
	_, ok := resp.(http.ResponseWriter)
	if !ok {
		t.Error("Response() should return an http.ResponseWriter")
	}
	
	// 验证我们可以写入响应 (Verify we can write to the response)
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("test"))
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

// TestGinContextParam 测试获取路径参数 (Test getting path parameters)
func TestGinContextParam(t *testing.T) {
	ginCtx, _ := setupTestContext("GET", "/users/123", "")
	
	// 设置路径参数 (Set path parameters)
	ginCtx.Params = gin.Params{
		{Key: "id", Value: "123"},
		{Key: "name", Value: "test"},
	}
	
	ctx := NewGinContext(ginCtx)
	
	if ctx.Param("id") != "123" {
		t.Errorf("Expected param id=123, got %s", ctx.Param("id"))
	}
	
	if ctx.Param("name") != "test" {
		t.Errorf("Expected param name=test, got %s", ctx.Param("name"))
	}
	
	if ctx.Param("nonexistent") != "" {
		t.Error("Nonexistent param should return empty string")
	}
}

// TestGinContextQuery 测试获取查询参数 (Test getting query parameters)
func TestGinContextQuery(t *testing.T) {
	ginCtx, _ := setupTestContext("GET", "/test?name=john&age=30", "")
	ctx := NewGinContext(ginCtx)
	
	if ctx.Query("name") != "john" {
		t.Errorf("Expected query name=john, got %s", ctx.Query("name"))
	}
	
	if ctx.Query("age") != "30" {
		t.Errorf("Expected query age=30, got %s", ctx.Query("age"))
	}
	
	if ctx.Query("nonexistent") != "" {
		t.Error("Nonexistent query should return empty string")
	}
}

// TestGinContextHeader 测试获取请求头 (Test getting request headers)
func TestGinContextHeader(t *testing.T) {
	ginCtx, _ := setupTestContext("GET", "/test", "")
	ginCtx.Request.Header.Set("Authorization", "Bearer token123")
	ginCtx.Request.Header.Set("Content-Type", "application/json")
	
	ctx := NewGinContext(ginCtx)
	
	if ctx.Header("Authorization") != "Bearer token123" {
		t.Errorf("Expected Authorization header, got %s", ctx.Header("Authorization"))
	}
	
	if ctx.Header("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type header, got %s", ctx.Header("Content-Type"))
	}
	
	if ctx.Header("Nonexistent") != "" {
		t.Error("Nonexistent header should return empty string")
	}
}

// TestGinContextSetHeader 测试设置响应头 (Test setting response headers)
func TestGinContextSetHeader(t *testing.T) {
	ginCtx, w := setupTestContext("GET", "/test", "")
	ctx := NewGinContext(ginCtx)
	
	ctx.SetHeader("X-Custom-Header", "test-value")
	ctx.SetHeader("Cache-Control", "no-cache")
	
	if w.Header().Get("X-Custom-Header") != "test-value" {
		t.Error("Custom header should be set")
	}
	
	if w.Header().Get("Cache-Control") != "no-cache" {
		t.Error("Cache-Control header should be set")
	}
}

// TestGinContextJSON 测试JSON响应 (Test JSON response)
func TestGinContextJSON(t *testing.T) {
	ginCtx, w := setupTestContext("GET", "/test", "")
	ctx := NewGinContext(ginCtx)
	
	data := map[string]interface{}{
		"message": "hello",
		"code":    200,
	}
	
	err := ctx.JSON(http.StatusOK, data)
	if err != nil {
		t.Fatalf("JSON() failed: %v", err)
	}
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	
	if w.Header().Get("Content-Type") != "application/json; charset=utf-8" {
		t.Error("Content-Type should be application/json")
	}
	
	var result map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	if result["message"] != "hello" {
		t.Error("Response message should be 'hello'")
	}
}

// TestGinContextString 测试字符串响应 (Test string response)
func TestGinContextString(t *testing.T) {
	ginCtx, w := setupTestContext("GET", "/test", "")
	ctx := NewGinContext(ginCtx)
	
	err := ctx.String(http.StatusOK, "Hello %s", "World")
	if err != nil {
		t.Fatalf("String() failed: %v", err)
	}
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	
	if w.Body.String() != "Hello World" {
		t.Errorf("Expected 'Hello World', got %s", w.Body.String())
	}
}

// TestGinContextData 测试原始数据响应 (Test raw data response)
func TestGinContextData(t *testing.T) {
	ginCtx, w := setupTestContext("GET", "/test", "")
	ctx := NewGinContext(ginCtx)
	
	data := []byte("binary data")
	err := ctx.Data(http.StatusOK, "application/octet-stream", data)
	if err != nil {
		t.Fatalf("Data() failed: %v", err)
	}
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
	
	if w.Header().Get("Content-Type") != "application/octet-stream" {
		t.Error("Content-Type should be application/octet-stream")
	}
	
	if !bytes.Equal(w.Body.Bytes(), data) {
		t.Error("Response body should match the data")
	}
}

// TestGinContextSetGet 测试设置和获取上下文值 (Test setting and getting context values)
func TestGinContextSetGet(t *testing.T) {
	ginCtx, _ := setupTestContext("GET", "/test", "")
	ctx := NewGinContext(ginCtx)
	
	// 测试Set和Get (Test Set and Get)
	ctx.Set("key1", "value1")
	ctx.Set("key2", 42)
	ctx.Set("key3", true)
	
	value1, exists1 := ctx.Get("key1")
	if !exists1 || value1 != "value1" {
		t.Error("key1 should exist and equal 'value1'")
	}
	
	value2, exists2 := ctx.Get("key2")
	if !exists2 || value2 != 42 {
		t.Error("key2 should exist and equal 42")
	}
	
	_, exists3 := ctx.Get("nonexistent")
	if exists3 {
		t.Error("nonexistent key should not exist")
	}
}

// TestGinContextGetTyped 测试类型化的Get方法 (Test typed Get methods)
func TestGinContextGetTyped(t *testing.T) {
	ginCtx, _ := setupTestContext("GET", "/test", "")
	ctx := NewGinContext(ginCtx)
	
	ctx.Set("string_key", "test_string")
	ctx.Set("int_key", 123)
	ctx.Set("bool_key", true)
	
	// 测试GetString (Test GetString)
	if ctx.GetString("string_key") != "test_string" {
		t.Error("GetString should return 'test_string'")
	}
	
	if ctx.GetString("nonexistent") != "" {
		t.Error("GetString for nonexistent key should return empty string")
	}
	
	// 测试GetInt (Test GetInt)
	if ctx.GetInt("int_key") != 123 {
		t.Error("GetInt should return 123")
	}
	
	if ctx.GetInt("nonexistent") != 0 {
		t.Error("GetInt for nonexistent key should return 0")
	}
	
	// 测试GetBool (Test GetBool)
	if ctx.GetBool("bool_key") != true {
		t.Error("GetBool should return true")
	}
	
	if ctx.GetBool("nonexistent") != false {
		t.Error("GetBool for nonexistent key should return false")
	}
}

// TestGinContextBind 测试数据绑定 (Test data binding)
func TestGinContextBind(t *testing.T) {
	jsonData := `{"name":"john","age":30}`
	ginCtx, _ := setupTestContext("POST", "/test", jsonData)
	ctx := NewGinContext(ginCtx)
	
	type TestStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	
	var data TestStruct
	err := ctx.Bind(&data)
	if err != nil {
		t.Fatalf("Bind() failed: %v", err)
	}
	
	if data.Name != "john" {
		t.Errorf("Expected name 'john', got %s", data.Name)
	}
	
	if data.Age != 30 {
		t.Errorf("Expected age 30, got %d", data.Age)
	}
}

// TestGinContextClientIP 测试获取客户端IP (Test getting client IP)
func TestGinContextClientIP(t *testing.T) {
	ginCtx, _ := setupTestContext("GET", "/test", "")
	ginCtx.Request.RemoteAddr = "192.168.1.1:12345"
	
	ctx := NewGinContext(ginCtx)
	
	clientIP := ctx.ClientIP()
	if clientIP == "" {
		t.Error("ClientIP should not be empty")
	}
}

// TestGinContextUserAgent 测试获取用户代理 (Test getting user agent)
func TestGinContextUserAgent(t *testing.T) {
	ginCtx, _ := setupTestContext("GET", "/test", "")
	ginCtx.Request.Header.Set("User-Agent", "Test-Agent/1.0")
	
	ctx := NewGinContext(ginCtx)
	
	if ctx.UserAgent() != "Test-Agent/1.0" {
		t.Errorf("Expected User-Agent 'Test-Agent/1.0', got %s", ctx.UserAgent())
	}
}

// TestGinContextMethod 测试获取请求方法 (Test getting request method)
func TestGinContextMethod(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	
	for _, method := range methods {
		ginCtx, _ := setupTestContext(method, "/test", "")
		ctx := NewGinContext(ginCtx)
		
		if ctx.Method() != method {
			t.Errorf("Expected method %s, got %s", method, ctx.Method())
		}
	}
}

// TestGinContextPath 测试获取请求路径 (Test getting request path)
func TestGinContextPath(t *testing.T) {
	ginCtx, _ := setupTestContext("GET", "/api/users/123", "")
	ctx := NewGinContext(ginCtx)
	
	if ctx.Path() != "/api/users/123" {
		t.Errorf("Expected path '/api/users/123', got %s", ctx.Path())
	}
}

// TestGinContextFullPath 测试获取完整路径模式 (Test getting full path pattern)
func TestGinContextFullPath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// 创建一个真实的Gin引擎和路由 (Create a real Gin engine and route)
	engine := gin.New()
	w := httptest.NewRecorder()
	
	var capturedFullPath string
	engine.GET("/api/users/:id", func(c *gin.Context) {
		capturedFullPath = c.FullPath()
		ctx := NewGinContext(c)
		c.JSON(200, gin.H{"full_path": ctx.FullPath()})
	})
	
	req := httptest.NewRequest("GET", "/api/users/123", nil)
	engine.ServeHTTP(w, req)
	
	if capturedFullPath != "/api/users/:id" {
		t.Errorf("Expected full path '/api/users/:id', got %s", capturedFullPath)
	}
}

// TestGinContextQueryArray 测试获取查询参数数组 (Test getting query parameter array)
func TestGinContextQueryArray(t *testing.T) {
	ginCtx, _ := setupTestContext("GET", "/test?tags=go&tags=web&tags=api", "")
	ctx := NewGinContext(ginCtx)
	
	// 类型断言获取GinContext以访问扩展方法 (Type assertion to get GinContext for extended methods)
	ginContextAdapter, ok := ctx.(*GinContext)
	if !ok {
		t.Fatal("Context should be *GinContext")
	}
	
	tags := ginContextAdapter.GetQueryArray("tags")
	if len(tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(tags))
	}
	
	expected := []string{"go", "web", "api"}
	for i, tag := range tags {
		if tag != expected[i] {
			t.Errorf("Expected tag %s, got %s", expected[i], tag)
		}
	}
}

// TestGinContextPostForm 测试获取表单数据 (Test getting form data)
func TestGinContextPostForm(t *testing.T) {
	formData := url.Values{}
	formData.Set("username", "john")
	formData.Set("password", "secret")
	
	ginCtx, _ := setupTestContext("POST", "/login", formData.Encode())
	ginCtx.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	ctx := NewGinContext(ginCtx)
	
	// 类型断言获取GinContext以访问扩展方法 (Type assertion to get GinContext for extended methods)
	ginContextAdapter, ok := ctx.(*GinContext)
	if !ok {
		t.Fatal("Context should be *GinContext")
	}
	
	if ginContextAdapter.PostForm("username") != "john" {
		t.Errorf("Expected username 'john', got %s", ginContextAdapter.PostForm("username"))
	}
	
	if ginContextAdapter.PostForm("password") != "secret" {
		t.Errorf("Expected password 'secret', got %s", ginContextAdapter.PostForm("password"))
	}
}

// BenchmarkGinContextJSON 基准测试JSON响应 (Benchmark JSON response)
func BenchmarkGinContextJSON(b *testing.B) {
	ginCtx, _ := setupTestContext("GET", "/test", "")
	ctx := NewGinContext(ginCtx)
	
	data := map[string]string{"message": "hello"}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx.JSON(http.StatusOK, data)
	}
}

// BenchmarkGinContextSetGet 基准测试设置和获取值 (Benchmark setting and getting values)
func BenchmarkGinContextSetGet(b *testing.B) {
	ginCtx, _ := setupTestContext("GET", "/test", "")
	ctx := NewGinContext(ginCtx)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx.Set("key", "value")
		ctx.Get("key")
	}
} 