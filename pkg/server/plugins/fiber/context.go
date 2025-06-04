/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Fiber上下文适配器实现 (Fiber context adapter implementation)
 */

package fiber

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
)

// FiberContext Fiber上下文适配器 (Fiber context adapter)
// 实现server.Context接口 (Implements server.Context interface)
type FiberContext struct {
	fiber            *fiber.Ctx                   // Fiber上下文 (Fiber context)
	serviceContainer services.ServiceContainer    // 服务容器 (Service container)
	store            map[string]interface{}       // 存储映射 (Storage map)
	storeMutex       sync.RWMutex                 // 存储互斥锁 (Storage mutex)
}

// NewFiberContext 创建Fiber上下文适配器 (Create Fiber context adapter)
func NewFiberContext(fiber *fiber.Ctx, serviceContainer services.ServiceContainer) server.Context {
	return &FiberContext{
		fiber:            fiber,
		serviceContainer: serviceContainer,
		store:            make(map[string]interface{}),
		storeMutex:       sync.RWMutex{},
	}
}

// Request 获取HTTP请求对象 (Get HTTP request object)
func (c *FiberContext) Request() *http.Request {
	// Fiber使用fasthttp，需要转换为标准http.Request
	// 这里返回一个基本的http.Request结构
	uri := c.fiber.Request().URI()
	url := &url.URL{
		Scheme:   string(uri.Scheme()),
		Host:     string(uri.Host()),
		Path:     string(uri.Path()),
		RawQuery: string(uri.QueryString()),
	}
	
	req := &http.Request{
		Method: c.fiber.Method(),
		URL:    url,
		Header: make(http.Header),
	}
	
	// 复制头部 (Copy headers)
	c.fiber.Request().Header.VisitAll(func(key, value []byte) {
		req.Header.Add(string(key), string(value))
	})
	
	return req
}

// Response 获取HTTP响应写入器 (Get HTTP response writer)
func (c *FiberContext) Response() http.ResponseWriter {
	// Fiber使用fasthttp，这里返回一个适配器
	return &fiberResponseWriter{ctx: c.fiber}
}

// Param 获取路径参数 (Get path parameter)
func (c *FiberContext) Param(key string) string {
	return c.fiber.Params(key)
}

// Query 获取查询参数 (Get query parameter)
func (c *FiberContext) Query(key string) string {
	return c.fiber.Query(key)
}

// Header 获取请求头 (Get request header)
func (c *FiberContext) Header(key string) string {
	return c.fiber.Get(key)
}

// SetHeader 设置响应头 (Set response header)
func (c *FiberContext) SetHeader(key, value string) {
	c.fiber.Set(key, value)
}

// JSON 返回JSON响应 (Return JSON response)
func (c *FiberContext) JSON(code int, obj interface{}) error {
	c.fiber.Status(code)
	return c.fiber.JSON(obj)
}

// String 返回字符串响应 (Return string response)
func (c *FiberContext) String(code int, format string, values ...interface{}) error {
	c.fiber.Status(code)
	if len(values) > 0 {
		return c.fiber.SendString(fmt.Sprintf(format, values...))
	}
	return c.fiber.SendString(format)
}

// Data 返回原始数据响应 (Return raw data response)
func (c *FiberContext) Data(code int, contentType string, data []byte) error {
	c.fiber.Status(code)
	c.fiber.Set("Content-Type", contentType)
	return c.fiber.Send(data)
}

// Set 设置上下文值 (Set context value)
func (c *FiberContext) Set(key string, value interface{}) {
	c.storeMutex.Lock()
	defer c.storeMutex.Unlock()
	c.store[key] = value
}

// Get 获取上下文值 (Get context value)
func (c *FiberContext) Get(key string) (interface{}, bool) {
	c.storeMutex.RLock()
	defer c.storeMutex.RUnlock()
	value, exists := c.store[key]
	return value, exists
}

// GetString 获取字符串类型的上下文值 (Get string type context value)
func (c *FiberContext) GetString(key string) string {
	if value, exists := c.Get(key); exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

// GetInt 获取整数类型的上下文值 (Get integer type context value)
func (c *FiberContext) GetInt(key string) int {
	if value, exists := c.Get(key); exists {
		switch v := value.(type) {
		case int:
			return v
		case int64:
			return int(v)
		case int32:
			return int(v)
		case float64:
			return int(v)
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				return i
			}
		}
	}
	return 0
}

// GetBool 获取布尔类型的上下文值 (Get boolean type context value)
func (c *FiberContext) GetBool(key string) bool {
	if value, exists := c.Get(key); exists {
		switch v := value.(type) {
		case bool:
			return v
		case string:
			if b, err := strconv.ParseBool(v); err == nil {
				return b
			}
		}
	}
	return false
}

// Bind 绑定请求数据到结构体 (Bind request data to struct)
func (c *FiberContext) Bind(obj interface{}) error {
	return c.fiber.BodyParser(obj)
}

// ClientIP 获取客户端IP地址 (Get client IP address)
func (c *FiberContext) ClientIP() string {
	return c.fiber.IP()
}

// UserAgent 获取用户代理 (Get user agent)
func (c *FiberContext) UserAgent() string {
	return c.fiber.Get("User-Agent")
}

// Method 获取请求方法 (Get request method)
func (c *FiberContext) Method() string {
	return c.fiber.Method()
}

// Path 获取请求路径 (Get request path)
func (c *FiberContext) Path() string {
	return c.fiber.Path()
}

// FullPath 获取完整路径模式 (Get full path pattern)
func (c *FiberContext) FullPath() string {
	return c.fiber.Route().Path
}

// GetFiberContext 获取原生Fiber上下文 (Get native Fiber context)
func (c *FiberContext) GetFiberContext() *fiber.Ctx {
	return c.fiber
}

// fiberResponseWriter Fiber响应写入器适配器 (Fiber response writer adapter)
type fiberResponseWriter struct {
	ctx *fiber.Ctx
}

// Header 实现http.ResponseWriter接口 (Implement http.ResponseWriter interface)
func (w *fiberResponseWriter) Header() http.Header {
	header := make(http.Header)
	// Fiber的响应头处理方式不同，这里提供基本实现
	return header
}

// Write 实现http.ResponseWriter接口 (Implement http.ResponseWriter interface)
func (w *fiberResponseWriter) Write(data []byte) (int, error) {
	return len(data), w.ctx.Send(data)
}

// WriteHeader 实现http.ResponseWriter接口 (Implement http.ResponseWriter interface)
func (w *fiberResponseWriter) WriteHeader(statusCode int) {
	w.ctx.Status(statusCode)
} 