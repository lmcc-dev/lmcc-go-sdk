/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Fiber上下文适配器测试 (Fiber context adapter tests)
 */

package fiber

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
)

func TestNewFiberContext(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		ctx := NewFiberContext(c, serviceContainer)
		assert.NotNil(t, ctx)
		
		// 验证接口实现
		fiberCtx, ok := ctx.(*FiberContext)
		assert.True(t, ok)
		assert.NotNil(t, fiberCtx.fiber)
		assert.Equal(t, serviceContainer, fiberCtx.serviceContainer)
		
		return c.SendString("ok")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFiberContext_Request(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		ctx := NewFiberContext(c, serviceContainer)
		req := ctx.Request()
		
		assert.NotNil(t, req)
		assert.Equal(t, "GET", req.Method)
		assert.Equal(t, "/test", req.URL.Path)
		
		return c.SendString("ok")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Test", "value")
	
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFiberContext_Param(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	
	app.Get("/users/:id", func(c *fiber.Ctx) error {
		ctx := NewFiberContext(c, serviceContainer)
		id := ctx.Param("id")
		
		assert.Equal(t, "123", id)
		
		return c.SendString("ok")
	})
	
	req := httptest.NewRequest("GET", "/users/123", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFiberContext_Query(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		ctx := NewFiberContext(c, serviceContainer)
		name := ctx.Query("name")
		
		assert.Equal(t, "test", name)
		
		return c.SendString("ok")
	})
	
	req := httptest.NewRequest("GET", "/test?name=test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFiberContext_Header(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		ctx := NewFiberContext(c, serviceContainer)
		value := ctx.Header("X-Test")
		
		assert.Equal(t, "test-value", value)
		
		return c.SendString("ok")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Test", "test-value")
	
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFiberContext_SetHeader(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		ctx := NewFiberContext(c, serviceContainer)
		ctx.SetHeader("X-Custom", "custom-value")
		
		return c.SendString("ok")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "custom-value", resp.Header.Get("X-Custom"))
}

func TestFiberContext_JSON(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		ctx := NewFiberContext(c, serviceContainer)
		return ctx.JSON(200, map[string]string{"message": "test"})
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	
	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	assert.Equal(t, "test", result["message"])
}

func TestFiberContext_String(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		ctx := NewFiberContext(c, serviceContainer)
		return ctx.String(200, "Hello %s", "World")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	
	buf := make([]byte, 1024)
	n, err := resp.Body.Read(buf)
	if err != nil && err.Error() != "EOF" {
		require.NoError(t, err)
	}
	assert.Equal(t, "Hello World", string(buf[:n]))
}

func TestFiberContext_Data(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		ctx := NewFiberContext(c, serviceContainer)
		data := []byte("test data")
		return ctx.Data(200, "text/plain", data)
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "text/plain", resp.Header.Get("Content-Type"))
	
	buf := make([]byte, 1024)
	n, err := resp.Body.Read(buf)
	if err != nil && err.Error() != "EOF" {
		require.NoError(t, err)
	}
	assert.Equal(t, "test data", string(buf[:n]))
}

func TestFiberContext_SetGet(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		ctx := NewFiberContext(c, serviceContainer)
		
		// 测试Set和Get
		ctx.Set("key", "value")
		value, exists := ctx.Get("key")
		assert.True(t, exists)
		assert.Equal(t, "value", value)
		
		// 测试不存在的键
		_, exists = ctx.Get("nonexistent")
		assert.False(t, exists)
		
		return c.SendString("ok")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFiberContext_GetString(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		ctx := NewFiberContext(c, serviceContainer)
		
		ctx.Set("string_key", "string_value")
		ctx.Set("int_key", 123)
		
		assert.Equal(t, "string_value", ctx.GetString("string_key"))
		assert.Equal(t, "", ctx.GetString("int_key"))
		assert.Equal(t, "", ctx.GetString("nonexistent"))
		
		return c.SendString("ok")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFiberContext_GetInt(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		ctx := NewFiberContext(c, serviceContainer)
		
		ctx.Set("int_key", 123)
		ctx.Set("int64_key", int64(456))
		ctx.Set("float_key", 78.9)
		ctx.Set("string_key", "999")
		ctx.Set("invalid_string", "abc")
		
		assert.Equal(t, 123, ctx.GetInt("int_key"))
		assert.Equal(t, 456, ctx.GetInt("int64_key"))
		assert.Equal(t, 78, ctx.GetInt("float_key"))
		assert.Equal(t, 999, ctx.GetInt("string_key"))
		assert.Equal(t, 0, ctx.GetInt("invalid_string"))
		assert.Equal(t, 0, ctx.GetInt("nonexistent"))
		
		return c.SendString("ok")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFiberContext_GetBool(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		ctx := NewFiberContext(c, serviceContainer)
		
		ctx.Set("bool_true", true)
		ctx.Set("bool_false", false)
		ctx.Set("string_true", "true")
		ctx.Set("string_false", "false")
		ctx.Set("invalid_string", "abc")
		
		assert.Equal(t, true, ctx.GetBool("bool_true"))
		assert.Equal(t, false, ctx.GetBool("bool_false"))
		assert.Equal(t, true, ctx.GetBool("string_true"))
		assert.Equal(t, false, ctx.GetBool("string_false"))
		assert.Equal(t, false, ctx.GetBool("invalid_string"))
		assert.Equal(t, false, ctx.GetBool("nonexistent"))
		
		return c.SendString("ok")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFiberContext_Bind(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	
	type TestData struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	
	app.Post("/test", func(c *fiber.Ctx) error {
		ctx := NewFiberContext(c, serviceContainer)
		
		var data TestData
		err := ctx.Bind(&data)
		assert.NoError(t, err)
		assert.Equal(t, "John", data.Name)
		assert.Equal(t, 30, data.Age)
		
		return c.SendString("ok")
	})
	
	testData := TestData{Name: "John", Age: 30}
	jsonData, _ := json.Marshal(testData)
	
	req := httptest.NewRequest("POST", "/test", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFiberContext_ClientIP(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		ctx := NewFiberContext(c, serviceContainer)
		ip := ctx.ClientIP()
		
		assert.NotEmpty(t, ip)
		
		return c.SendString("ok")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFiberContext_UserAgent(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		ctx := NewFiberContext(c, serviceContainer)
		ua := ctx.UserAgent()
		
		assert.Equal(t, "test-agent", ua)
		
		return c.SendString("ok")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("User-Agent", "test-agent")
	
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFiberContext_Method(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	
	app.Post("/test", func(c *fiber.Ctx) error {
		ctx := NewFiberContext(c, serviceContainer)
		method := ctx.Method()
		
		assert.Equal(t, "POST", method)
		
		return c.SendString("ok")
	})
	
	req := httptest.NewRequest("POST", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFiberContext_Path(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	
	app.Get("/test/path", func(c *fiber.Ctx) error {
		ctx := NewFiberContext(c, serviceContainer)
		path := ctx.Path()
		
		assert.Equal(t, "/test/path", path)
		
		return c.SendString("ok")
	})
	
	req := httptest.NewRequest("GET", "/test/path", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFiberContext_FullPath(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	
	app.Get("/users/:id", func(c *fiber.Ctx) error {
		ctx := NewFiberContext(c, serviceContainer)
		fullPath := ctx.FullPath()
		
		assert.Equal(t, "/users/:id", fullPath)
		
		return c.SendString("ok")
	})
	
	req := httptest.NewRequest("GET", "/users/123", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFiberContext_GetFiberContext(t *testing.T) {
	app := fiber.New()
	serviceContainer := services.NewServiceContainer()
	
	app.Get("/test", func(c *fiber.Ctx) error {
		ctx := NewFiberContext(c, serviceContainer)
		fiberCtx := ctx.(*FiberContext).GetFiberContext()
		
		assert.Equal(t, c, fiberCtx)
		
		return c.SendString("ok")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
} 