/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Echo上下文适配器测试 (Echo context adapter tests)
 */

package echo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
)

func TestNewEchoContext(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	echoCtx := e.NewContext(req, rec)
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	
	ctx := NewEchoContext(echoCtx, serviceContainer)
	assert.NotNil(t, ctx)
	
	// 验证接口实现 (Verify interface implementation)
	_, ok := ctx.(server.Context)
	assert.True(t, ok)
	
	// 验证类型转换 (Verify type conversion)
	echoContext, ok := ctx.(*EchoContext)
	assert.True(t, ok)
	assert.Equal(t, echoCtx, echoContext.GetEchoContext())
}

func TestEchoContext_Request(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/test?param=value", strings.NewReader("body"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	echoCtx := e.NewContext(req, rec)
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	ctx := NewEchoContext(echoCtx, serviceContainer)
	
	// 测试Request方法 (Test Request method)
	request := ctx.Request()
	assert.Equal(t, req, request)
	assert.Equal(t, http.MethodPost, request.Method)
	assert.Equal(t, "/test", request.URL.Path)
	assert.Equal(t, "application/json", request.Header.Get("Content-Type"))
}

func TestEchoContext_Response(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	echoCtx := e.NewContext(req, rec)
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	ctx := NewEchoContext(echoCtx, serviceContainer)
	
	// 测试Response方法 (Test Response method)
	response := ctx.Response()
	assert.NotNil(t, response)
	assert.Equal(t, rec, response)
}

func TestEchoContext_Param(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users/123", nil)
	rec := httptest.NewRecorder()
	echoCtx := e.NewContext(req, rec)
	
	// 设置路径参数 (Set path parameters)
	echoCtx.SetParamNames("id")
	echoCtx.SetParamValues("123")
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	ctx := NewEchoContext(echoCtx, serviceContainer)
	
	// 测试Param方法 (Test Param method)
	id := ctx.Param("id")
	assert.Equal(t, "123", id)
	
	// 测试不存在的参数 (Test non-existent parameter)
	name := ctx.Param("name")
	assert.Equal(t, "", name)
}

func TestEchoContext_Query(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test?name=john&age=30", nil)
	rec := httptest.NewRecorder()
	echoCtx := e.NewContext(req, rec)
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	ctx := NewEchoContext(echoCtx, serviceContainer)
	
	// 测试Query方法 (Test Query method)
	name := ctx.Query("name")
	assert.Equal(t, "john", name)
	
	age := ctx.Query("age")
	assert.Equal(t, "30", age)
	
	// 测试不存在的查询参数 (Test non-existent query parameter)
	email := ctx.Query("email")
	assert.Equal(t, "", email)
}

func TestEchoContext_Header(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer token123")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	echoCtx := e.NewContext(req, rec)
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	ctx := NewEchoContext(echoCtx, serviceContainer)
	
	// 测试Header方法 (Test Header method)
	auth := ctx.Header("Authorization")
	assert.Equal(t, "Bearer token123", auth)
	
	contentType := ctx.Header("Content-Type")
	assert.Equal(t, "application/json", contentType)
	
	// 测试不存在的头部 (Test non-existent header)
	custom := ctx.Header("X-Custom")
	assert.Equal(t, "", custom)
}

func TestEchoContext_SetHeader(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	echoCtx := e.NewContext(req, rec)
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	ctx := NewEchoContext(echoCtx, serviceContainer)
	
	// 测试SetHeader方法 (Test SetHeader method)
	ctx.SetHeader("X-Custom", "custom-value")
	ctx.SetHeader("Cache-Control", "no-cache")
	
	// 验证头部已设置 (Verify headers are set)
	assert.Equal(t, "custom-value", rec.Header().Get("X-Custom"))
	assert.Equal(t, "no-cache", rec.Header().Get("Cache-Control"))
}

func TestEchoContext_JSON(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	echoCtx := e.NewContext(req, rec)
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	ctx := NewEchoContext(echoCtx, serviceContainer)
	
	// 测试JSON方法 (Test JSON method)
	data := map[string]interface{}{
		"message": "hello",
		"code":    200,
	}
	
	err := ctx.JSON(200, data)
	assert.NoError(t, err)
	
	// 验证响应 (Verify response)
	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	
	var result map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, "hello", result["message"])
	assert.Equal(t, float64(200), result["code"]) // JSON numbers are float64
}

func TestEchoContext_String(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	echoCtx := e.NewContext(req, rec)
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	ctx := NewEchoContext(echoCtx, serviceContainer)
	
	// 测试String方法 - 无格式化 (Test String method - no formatting)
	err := ctx.String(200, "hello world")
	assert.NoError(t, err)
	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "hello world", rec.Body.String())
	
	// 重置recorder (Reset recorder)
	rec = httptest.NewRecorder()
	echoCtx = e.NewContext(req, rec)
	ctx = NewEchoContext(echoCtx, serviceContainer)
	
	// 测试String方法 - 带格式化 (Test String method - with formatting)
	err = ctx.String(201, "Hello %s, you are %d years old", "John", 30)
	assert.NoError(t, err)
	assert.Equal(t, 201, rec.Code)
	assert.Equal(t, "Hello John, you are 30 years old", rec.Body.String())
}

func TestEchoContext_Data(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	echoCtx := e.NewContext(req, rec)
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	ctx := NewEchoContext(echoCtx, serviceContainer)
	
	// 测试Data方法 (Test Data method)
	data := []byte("binary data")
	err := ctx.Data(200, "application/octet-stream", data)
	assert.NoError(t, err)
	
	// 验证响应 (Verify response)
	assert.Equal(t, 200, rec.Code)
	assert.Equal(t, "application/octet-stream", rec.Header().Get("Content-Type"))
	assert.Equal(t, data, rec.Body.Bytes())
}

func TestEchoContext_SetGet(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	echoCtx := e.NewContext(req, rec)
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	ctx := NewEchoContext(echoCtx, serviceContainer)
	
	// 测试Set和Get方法 (Test Set and Get methods)
	ctx.Set("user_id", 123)
	ctx.Set("username", "john")
	ctx.Set("is_admin", true)
	
	// 获取值 (Get values)
	userID, exists := ctx.Get("user_id")
	assert.True(t, exists)
	assert.Equal(t, 123, userID)
	
	username, exists := ctx.Get("username")
	assert.True(t, exists)
	assert.Equal(t, "john", username)
	
	isAdmin, exists := ctx.Get("is_admin")
	assert.True(t, exists)
	assert.Equal(t, true, isAdmin)
	
	// 测试不存在的键 (Test non-existent key)
	_, exists = ctx.Get("non_existent")
	assert.False(t, exists)
}

func TestEchoContext_GetTyped(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	echoCtx := e.NewContext(req, rec)
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	ctx := NewEchoContext(echoCtx, serviceContainer)
	
	// 设置不同类型的值 (Set different types of values)
	ctx.Set("string_val", "hello")
	ctx.Set("int_val", 42)
	ctx.Set("int64_val", int64(100))
	ctx.Set("int32_val", int32(50))
	ctx.Set("float64_val", 3.14)
	ctx.Set("bool_val", true)
	ctx.Set("string_int", "123")
	ctx.Set("string_bool", "true")
	
	// 测试GetString (Test GetString)
	assert.Equal(t, "hello", ctx.GetString("string_val"))
	assert.Equal(t, "", ctx.GetString("int_val"))
	assert.Equal(t, "", ctx.GetString("non_existent"))
	
	// 测试GetInt (Test GetInt)
	assert.Equal(t, 42, ctx.GetInt("int_val"))
	assert.Equal(t, 100, ctx.GetInt("int64_val"))
	assert.Equal(t, 50, ctx.GetInt("int32_val"))
	assert.Equal(t, 3, ctx.GetInt("float64_val"))
	assert.Equal(t, 123, ctx.GetInt("string_int"))
	assert.Equal(t, 0, ctx.GetInt("string_val"))
	assert.Equal(t, 0, ctx.GetInt("non_existent"))
	
	// 测试GetBool (Test GetBool)
	assert.Equal(t, true, ctx.GetBool("bool_val"))
	assert.Equal(t, true, ctx.GetBool("string_bool"))
	assert.Equal(t, false, ctx.GetBool("string_val"))
	assert.Equal(t, false, ctx.GetBool("non_existent"))
}

func TestEchoContext_Bind(t *testing.T) {
	e := echo.New()
	
	// 测试JSON绑定 (Test JSON binding)
	jsonData := `{"name":"john","age":30}`
	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	echoCtx := e.NewContext(req, rec)
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	ctx := NewEchoContext(echoCtx, serviceContainer)
	
	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	
	var user User
	err := ctx.Bind(&user)
	assert.NoError(t, err)
	assert.Equal(t, "john", user.Name)
	assert.Equal(t, 30, user.Age)
}

func TestEchoContext_ClientIP(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Real-IP", "192.168.1.100")
	rec := httptest.NewRecorder()
	echoCtx := e.NewContext(req, rec)
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	ctx := NewEchoContext(echoCtx, serviceContainer)
	
	// 测试ClientIP方法 (Test ClientIP method)
	ip := ctx.ClientIP()
	assert.Equal(t, "192.168.1.100", ip)
}

func TestEchoContext_UserAgent(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Test Browser")
	rec := httptest.NewRecorder()
	echoCtx := e.NewContext(req, rec)
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	ctx := NewEchoContext(echoCtx, serviceContainer)
	
	// 测试UserAgent方法 (Test UserAgent method)
	userAgent := ctx.UserAgent()
	assert.Equal(t, "Mozilla/5.0 Test Browser", userAgent)
}

func TestEchoContext_Method(t *testing.T) {
	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}
	
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(method, "/test", nil)
			rec := httptest.NewRecorder()
			echoCtx := e.NewContext(req, rec)
			
			serviceContainer := services.NewServiceContainerWithDefaults()
			ctx := NewEchoContext(echoCtx, serviceContainer)
			
			// 测试Method方法 (Test Method method)
			assert.Equal(t, method, ctx.Method())
		})
	}
}

func TestEchoContext_Path(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users?id=123", nil)
	rec := httptest.NewRecorder()
	echoCtx := e.NewContext(req, rec)
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	ctx := NewEchoContext(echoCtx, serviceContainer)
	
	// 测试Path方法 (Test Path method)
	path := ctx.Path()
	assert.Equal(t, "/api/v1/users", path)
}

func TestEchoContext_FullPath(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/123", nil)
	rec := httptest.NewRecorder()
	echoCtx := e.NewContext(req, rec)
	
	// 设置路由路径 (Set route path)
	echoCtx.SetPath("/api/v1/users/:id")
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	ctx := NewEchoContext(echoCtx, serviceContainer)
	
	// 测试FullPath方法 (Test FullPath method)
	fullPath := ctx.FullPath()
	assert.Equal(t, "/api/v1/users/:id", fullPath)
}

func TestEchoContext_ConcurrentAccess(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	echoCtx := e.NewContext(req, rec)
	
	serviceContainer := services.NewServiceContainerWithDefaults()
	ctx := NewEchoContext(echoCtx, serviceContainer)
	
	// 测试并发访问 (Test concurrent access)
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func(id int) {
			key := fmt.Sprintf("key_%d", id)
			value := fmt.Sprintf("value_%d", id)
			
			ctx.Set(key, value)
			
			retrieved, exists := ctx.Get(key)
			assert.True(t, exists)
			assert.Equal(t, value, retrieved)
			
			done <- true
		}(i)
	}
	
	// 等待所有goroutine完成 (Wait for all goroutines to complete)
	for i := 0; i < 10; i++ {
		<-done
	}
} 