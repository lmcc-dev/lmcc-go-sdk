/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Gin Context适配器实现
 */

package gin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
)

// GinContext Gin上下文适配器 (Gin context adapter)
// 将Gin的Context适配到统一的Context接口 (Adapts Gin's Context to unified Context interface)
type GinContext struct {
	ginCtx *gin.Context
}

// NewGinContext 创建Gin上下文适配器 (Create Gin context adapter)
func NewGinContext(ginCtx *gin.Context) server.Context {
	return &GinContext{
		ginCtx: ginCtx,
	}
}

// Request 获取HTTP请求对象 (Get HTTP request object)
func (c *GinContext) Request() *http.Request {
	return c.ginCtx.Request
}

// Response 获取HTTP响应写入器 (Get HTTP response writer)
func (c *GinContext) Response() http.ResponseWriter {
	return c.ginCtx.Writer
}

// Param 获取路径参数 (Get path parameter)
func (c *GinContext) Param(key string) string {
	return c.ginCtx.Param(key)
}

// Query 获取查询参数 (Get query parameter)
func (c *GinContext) Query(key string) string {
	return c.ginCtx.Query(key)
}

// Header 获取请求头 (Get request header)
func (c *GinContext) Header(key string) string {
	return c.ginCtx.GetHeader(key)
}

// SetHeader 设置响应头 (Set response header)
func (c *GinContext) SetHeader(key, value string) {
	c.ginCtx.Header(key, value)
}

// JSON 返回JSON响应 (Return JSON response)
func (c *GinContext) JSON(code int, obj interface{}) error {
	c.ginCtx.JSON(code, obj)
	return nil
}

// String 返回字符串响应 (Return string response)
func (c *GinContext) String(code int, format string, values ...interface{}) error {
	c.ginCtx.String(code, format, values...)
	return nil
}

// Data 返回原始数据响应 (Return raw data response)
func (c *GinContext) Data(code int, contentType string, data []byte) error {
	c.ginCtx.Data(code, contentType, data)
	return nil
}

// Set 设置上下文值 (Set context value)
func (c *GinContext) Set(key string, value interface{}) {
	c.ginCtx.Set(key, value)
}

// Get 获取上下文值 (Get context value)
func (c *GinContext) Get(key string) (interface{}, bool) {
	return c.ginCtx.Get(key)
}

// GetString 获取字符串类型的上下文值 (Get string type context value)
func (c *GinContext) GetString(key string) string {
	return c.ginCtx.GetString(key)
}

// GetInt 获取整数类型的上下文值 (Get integer type context value)
func (c *GinContext) GetInt(key string) int {
	return c.ginCtx.GetInt(key)
}

// GetBool 获取布尔类型的上下文值 (Get boolean type context value)
func (c *GinContext) GetBool(key string) bool {
	return c.ginCtx.GetBool(key)
}

// Bind 绑定请求数据到结构体 (Bind request data to struct)
func (c *GinContext) Bind(obj interface{}) error {
	return c.ginCtx.ShouldBind(obj)
}

// ClientIP 获取客户端IP地址 (Get client IP address)
func (c *GinContext) ClientIP() string {
	return c.ginCtx.ClientIP()
}

// UserAgent 获取用户代理 (Get user agent)
func (c *GinContext) UserAgent() string {
	return c.ginCtx.GetHeader("User-Agent")
}

// Method 获取请求方法 (Get request method)
func (c *GinContext) Method() string {
	return c.ginCtx.Request.Method
}

// Path 获取请求路径 (Get request path)
func (c *GinContext) Path() string {
	return c.ginCtx.Request.URL.Path
}

// FullPath 获取完整路径模式 (Get full path pattern)
func (c *GinContext) FullPath() string {
	return c.ginCtx.FullPath()
}

// GetGinContext 获取原生Gin上下文 (Get native Gin context)
// 用于访问Gin特定功能 (Used to access Gin-specific features)
func (c *GinContext) GetGinContext() *gin.Context {
	return c.ginCtx
}

// 扩展方法 - Gin特有功能 (Extended methods - Gin-specific features)

// PostForm 获取POST表单参数 (Get POST form parameter)
func (c *GinContext) PostForm(key string) string {
	return c.ginCtx.PostForm(key)
}

// DefaultPostForm 获取POST表单参数，带默认值 (Get POST form parameter with default value)
func (c *GinContext) DefaultPostForm(key, defaultValue string) string {
	return c.ginCtx.DefaultPostForm(key, defaultValue)
}

// FormFile 获取上传文件 (Get uploaded file)
func (c *GinContext) FormFile(name string) (*http.Request, error) {
	file, header, err := c.ginCtx.Request.FormFile(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	// 创建新的请求对象包含文件信息 (Create new request object containing file info)
	// 这里简化处理，实际使用中可能需要更复杂的文件处理 (Simplified handling, actual use may require more complex file processing)
	c.ginCtx.Set("uploaded_file_header", header)
	return c.ginCtx.Request, nil
}

// Cookie 获取Cookie值 (Get cookie value)
func (c *GinContext) Cookie(name string) (string, error) {
	return c.ginCtx.Cookie(name)
}

// SetCookie 设置Cookie (Set cookie)
func (c *GinContext) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	c.ginCtx.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
}

// Redirect 重定向 (Redirect)
func (c *GinContext) Redirect(code int, location string) {
	c.ginCtx.Redirect(code, location)
}

// Status 设置响应状态码 (Set response status code)
func (c *GinContext) Status(code int) {
	c.ginCtx.Status(code)
}

// GetRawData 获取原始请求数据 (Get raw request data)
func (c *GinContext) GetRawData() ([]byte, error) {
	return c.ginCtx.GetRawData()
}

// IsWebsocket 检查是否为WebSocket请求 (Check if it's a WebSocket request)
func (c *GinContext) IsWebsocket() bool {
	return c.ginCtx.IsWebsocket()
}

// ContentType 获取请求内容类型 (Get request content type)
func (c *GinContext) ContentType() string {
	return c.ginCtx.ContentType()
}

// GetQueryArray 获取查询参数数组 (Get query parameter array)
func (c *GinContext) GetQueryArray(key string) []string {
	return c.ginCtx.QueryArray(key)
}

// GetPostFormArray 获取POST表单参数数组 (Get POST form parameter array)
func (c *GinContext) GetPostFormArray(key string) []string {
	return c.ginCtx.PostFormArray(key)
}

// GetQueryMap 获取查询参数映射 (Get query parameter map)
func (c *GinContext) GetQueryMap(key string) map[string]string {
	return c.ginCtx.QueryMap(key)
}

// GetPostFormMap 获取POST表单参数映射 (Get POST form parameter map)
func (c *GinContext) GetPostFormMap(key string) map[string]string {
	return c.ginCtx.PostFormMap(key)
}

// Abort 中止请求处理 (Abort request processing)
func (c *GinContext) Abort() {
	c.ginCtx.Abort()
}

// AbortWithStatus 中止请求并设置状态码 (Abort request with status code)
func (c *GinContext) AbortWithStatus(code int) {
	c.ginCtx.AbortWithStatus(code)
}

// AbortWithStatusJSON 中止请求并返回JSON (Abort request with JSON response)
func (c *GinContext) AbortWithStatusJSON(code int, jsonObj interface{}) {
	c.ginCtx.AbortWithStatusJSON(code, jsonObj)
}

// IsAborted 检查请求是否已中止 (Check if request is aborted)
func (c *GinContext) IsAborted() bool {
	return c.ginCtx.IsAborted()
}

// Next 调用下一个处理器 (Call next handler)
func (c *GinContext) Next() {
	c.ginCtx.Next()
}

// 类型转换辅助方法 (Type conversion helper methods)

// GetInt64 获取int64类型的上下文值 (Get int64 type context value)
func (c *GinContext) GetInt64(key string) int64 {
	return c.ginCtx.GetInt64(key)
}

// GetUint 获取uint类型的上下文值 (Get uint type context value)
func (c *GinContext) GetUint(key string) uint {
	return c.ginCtx.GetUint(key)
}

// GetUint64 获取uint64类型的上下文值 (Get uint64 type context value)
func (c *GinContext) GetUint64(key string) uint64 {
	return c.ginCtx.GetUint64(key)
}

// GetFloat64 获取float64类型的上下文值 (Get float64 type context value)
func (c *GinContext) GetFloat64(key string) float64 {
	return c.ginCtx.GetFloat64(key)
}

// GetTime 获取时间类型的上下文值 (Get time type context value)
func (c *GinContext) GetTime(key string) (time.Time, error) {
	value := c.ginCtx.GetTime(key)
	return value, nil
}

// GetDuration 获取持续时间类型的上下文值 (Get duration type context value)
func (c *GinContext) GetDuration(key string) (time.Duration, error) {
	value := c.ginCtx.GetDuration(key)
	return value, nil
}

// GetStringSlice 获取字符串切片类型的上下文值 (Get string slice type context value)
func (c *GinContext) GetStringSlice(key string) []string {
	return c.ginCtx.GetStringSlice(key)
}

// GetStringMap 获取字符串映射类型的上下文值 (Get string map type context value)
func (c *GinContext) GetStringMap(key string) map[string]interface{} {
	return c.ginCtx.GetStringMap(key)
}

// GetStringMapString 获取字符串到字符串映射类型的上下文值 (Get string to string map type context value)
func (c *GinContext) GetStringMapString(key string) map[string]string {
	return c.ginCtx.GetStringMapString(key)
}

// GetStringMapStringSlice 获取字符串到字符串切片映射类型的上下文值 (Get string to string slice map type context value)
func (c *GinContext) GetStringMapStringSlice(key string) map[string][]string {
	return c.ginCtx.GetStringMapStringSlice(key)
}