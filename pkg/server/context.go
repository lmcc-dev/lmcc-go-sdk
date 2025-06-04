/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: 统一的请求上下文抽象实现
 */

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

// BaseContext 基础上下文实现 (Base context implementation)
// 提供通用的上下文功能 (Provides common context functionality)
type BaseContext struct {
	request  *http.Request
	response http.ResponseWriter
	params   map[string]string
	store    map[string]interface{}
	mutex    sync.RWMutex
}

// NewBaseContext 创建基础上下文 (Create base context)
func NewBaseContext(req *http.Request, resp http.ResponseWriter) *BaseContext {
	return &BaseContext{
		request:  req,
		response: resp,
		params:   make(map[string]string),
		store:    make(map[string]interface{}),
	}
}

// Request 获取HTTP请求对象 (Get HTTP request object)
func (c *BaseContext) Request() *http.Request {
	return c.request
}

// Response 获取HTTP响应写入器 (Get HTTP response writer)
func (c *BaseContext) Response() http.ResponseWriter {
	return c.response
}

// Param 获取路径参数 (Get path parameter)
func (c *BaseContext) Param(key string) string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.params[key]
}

// SetParam 设置路径参数 (Set path parameter)
// 内部方法，用于框架适配器设置参数 (Internal method for framework adapters to set parameters)
func (c *BaseContext) SetParam(key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.params[key] = value
}

// Query 获取查询参数 (Get query parameter)
func (c *BaseContext) Query(key string) string {
	return c.request.URL.Query().Get(key)
}

// Header 获取请求头 (Get request header)
func (c *BaseContext) Header(key string) string {
	return c.request.Header.Get(key)
}

// SetHeader 设置响应头 (Set response header)
func (c *BaseContext) SetHeader(key, value string) {
	c.response.Header().Set(key, value)
}

// JSON 返回JSON响应 (Return JSON response)
func (c *BaseContext) JSON(code int, obj interface{}) error {
	c.SetHeader("Content-Type", "application/json; charset=utf-8")
	c.response.WriteHeader(code)
	
	encoder := json.NewEncoder(c.response)
	return encoder.Encode(obj)
}

// String 返回字符串响应 (Return string response)
func (c *BaseContext) String(code int, format string, values ...interface{}) error {
	c.SetHeader("Content-Type", "text/plain; charset=utf-8")
	c.response.WriteHeader(code)
	
	if len(values) > 0 {
		_, err := fmt.Fprintf(c.response, format, values...)
		return err
	}
	
	_, err := c.response.Write([]byte(format))
	return err
}

// Data 返回原始数据响应 (Return raw data response)
func (c *BaseContext) Data(code int, contentType string, data []byte) error {
	c.SetHeader("Content-Type", contentType)
	c.response.WriteHeader(code)
	
	_, err := c.response.Write(data)
	return err
}

// Set 设置上下文值 (Set context value)
func (c *BaseContext) Set(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.store[key] = value
}

// Get 获取上下文值 (Get context value)
func (c *BaseContext) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	value, exists := c.store[key]
	return value, exists
}

// GetString 获取字符串类型的上下文值 (Get string type context value)
func (c *BaseContext) GetString(key string) string {
	if value, exists := c.Get(key); exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

// GetInt 获取整数类型的上下文值 (Get integer type context value)
func (c *BaseContext) GetInt(key string) int {
	if value, exists := c.Get(key); exists {
		switch v := value.(type) {
		case int:
			return v
		case int64:
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
func (c *BaseContext) GetBool(key string) bool {
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
func (c *BaseContext) Bind(obj interface{}) error {
	contentType := c.Header("Content-Type")
	
	switch {
	case contentType == "application/json":
		decoder := json.NewDecoder(c.request.Body)
		return decoder.Decode(obj)
	default:
		return fmt.Errorf("unsupported content type: %s", contentType)
	}
}

// ClientIP 获取客户端IP地址 (Get client IP address)
func (c *BaseContext) ClientIP() string {
	// 检查X-Forwarded-For头 (Check X-Forwarded-For header)
	if xff := c.Header("X-Forwarded-For"); xff != "" {
		return xff
	}
	
	// 检查X-Real-IP头 (Check X-Real-IP header)
	if xri := c.Header("X-Real-IP"); xri != "" {
		return xri
	}
	
	// 返回远程地址 (Return remote address)
	return c.request.RemoteAddr
}

// UserAgent 获取用户代理 (Get user agent)
func (c *BaseContext) UserAgent() string {
	return c.Header("User-Agent")
}

// Method 获取请求方法 (Get request method)
func (c *BaseContext) Method() string {
	return c.request.Method
}

// Path 获取请求路径 (Get request path)
func (c *BaseContext) Path() string {
	return c.request.URL.Path
}

// FullPath 获取完整路径模式 (Get full path pattern)
// 默认实现返回请求路径，具体框架可以重写此方法 (Default implementation returns request path, specific frameworks can override this method)
func (c *BaseContext) FullPath() string {
	return c.request.URL.Path
}