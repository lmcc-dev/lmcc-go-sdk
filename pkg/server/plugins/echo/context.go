/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Echo上下文适配器实现 (Echo context adapter implementation)
 */

package echo

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
)

// EchoContext Echo上下文适配器 (Echo context adapter)
// 实现server.Context接口 (Implements server.Context interface)
type EchoContext struct {
	echo             echo.Context                 // Echo上下文 (Echo context)
	serviceContainer services.ServiceContainer    // 服务容器 (Service container)
	store            map[string]interface{}       // 存储映射 (Storage map)
	storeMutex       sync.RWMutex                 // 存储互斥锁 (Storage mutex)
}

// NewEchoContext 创建Echo上下文适配器 (Create Echo context adapter)
func NewEchoContext(echoCtx echo.Context, serviceContainer services.ServiceContainer) server.Context {
	return &EchoContext{
		echo:             echoCtx,
		serviceContainer: serviceContainer,
		store:            make(map[string]interface{}),
	}
}

// Request 获取HTTP请求对象 (Get HTTP request object)
func (c *EchoContext) Request() *http.Request {
	return c.echo.Request()
}

// Response 获取HTTP响应写入器 (Get HTTP response writer)
func (c *EchoContext) Response() http.ResponseWriter {
	return c.echo.Response().Writer
}

// Param 获取路径参数 (Get path parameter)
func (c *EchoContext) Param(key string) string {
	return c.echo.Param(key)
}

// Query 获取查询参数 (Get query parameter)
func (c *EchoContext) Query(key string) string {
	return c.echo.QueryParam(key)
}

// Header 获取请求头 (Get request header)
func (c *EchoContext) Header(key string) string {
	return c.echo.Request().Header.Get(key)
}

// SetHeader 设置响应头 (Set response header)
func (c *EchoContext) SetHeader(key, value string) {
	c.echo.Response().Header().Set(key, value)
}

// JSON 返回JSON响应 (Return JSON response)
func (c *EchoContext) JSON(code int, obj interface{}) error {
	return c.echo.JSON(code, obj)
}

// String 返回字符串响应 (Return string response)
func (c *EchoContext) String(code int, format string, values ...interface{}) error {
	if len(values) > 0 {
		content := fmt.Sprintf(format, values...)
		return c.echo.String(code, content)
	}
	return c.echo.String(code, format)
}

// Data 返回原始数据响应 (Return raw data response)
func (c *EchoContext) Data(code int, contentType string, data []byte) error {
	return c.echo.Blob(code, contentType, data)
}

// Set 设置上下文值 (Set context value)
func (c *EchoContext) Set(key string, value interface{}) {
	c.storeMutex.Lock()
	defer c.storeMutex.Unlock()
	c.store[key] = value
	
	// 同时设置到Echo上下文中 (Also set to Echo context)
	c.echo.Set(key, value)
}

// Get 获取上下文值 (Get context value)
func (c *EchoContext) Get(key string) (interface{}, bool) {
	c.storeMutex.RLock()
	defer c.storeMutex.RUnlock()
	
	// 首先从本地存储获取 (First get from local storage)
	if value, exists := c.store[key]; exists {
		return value, true
	}
	
	// 然后从Echo上下文获取 (Then get from Echo context)
	value := c.echo.Get(key)
	if value != nil {
		return value, true
	}
	
	return nil, false
}

// GetString 获取字符串类型的上下文值 (Get string type context value)
func (c *EchoContext) GetString(key string) string {
	if value, exists := c.Get(key); exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

// GetInt 获取整数类型的上下文值 (Get integer type context value)
func (c *EchoContext) GetInt(key string) int {
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
func (c *EchoContext) GetBool(key string) bool {
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
func (c *EchoContext) Bind(obj interface{}) error {
	return c.echo.Bind(obj)
}

// ClientIP 获取客户端IP地址 (Get client IP address)
func (c *EchoContext) ClientIP() string {
	return c.echo.RealIP()
}

// UserAgent 获取用户代理 (Get user agent)
func (c *EchoContext) UserAgent() string {
	return c.echo.Request().UserAgent()
}

// Method 获取请求方法 (Get request method)
func (c *EchoContext) Method() string {
	return c.echo.Request().Method
}

// Path 获取请求路径 (Get request path)
func (c *EchoContext) Path() string {
	return c.echo.Request().URL.Path
}

// FullPath 获取完整路径模式 (Get full path pattern)
func (c *EchoContext) FullPath() string {
	return c.echo.Path()
}

// GetEchoContext 获取原生Echo上下文 (Get native Echo context)
// 提供类型安全的访问方法 (Provides type-safe access method)
func (c *EchoContext) GetEchoContext() echo.Context {
	return c.echo
}

// 验证EchoContext实现了server.Context接口 (Verify EchoContext implements server.Context interface)
var _ server.Context = (*EchoContext)(nil) 