/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Integration tests for CORS middleware with Gin framework / CORS中间件与Gin框架的集成测试
 */

package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	"github.com/stretchr/testify/assert"
)

// TestCORSIntegrationWithGin 测试CORS中间件与Gin框架的完整集成 (Test complete integration of CORS middleware with Gin framework)
func TestCORSIntegrationWithGin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// 创建CORS配置 (Create CORS configuration)
	config := &server.CORSConfig{
		Enabled:          true,
		AllowOrigins:     []string{"https://example.com", "https://app.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"X-Total-Count", "X-User-ID"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
	
	corsMiddleware := NewCORSMiddleware(config).(*CORSMiddleware)
	engine := gin.Default()
	
	// 应用CORS中间件 (Apply CORS middleware)
	engine.Use(corsMiddleware.GetGinHandler())
	
	// 添加其他中间件测试兼容性 (Add other middleware to test compatibility)
	engine.Use(func(c *gin.Context) {
		c.Header("X-Custom-Middleware", "applied")
		c.Next()
	})
	
	// 定义测试路由 (Define test routes)
	engine.GET("/api/users", func(c *gin.Context) {
		c.Header("X-Total-Count", "100")
		c.JSON(http.StatusOK, gin.H{
			"users": []gin.H{
				{"id": 1, "name": "Alice"},
				{"id": 2, "name": "Bob"},
			},
		})
	})
	
	engine.POST("/api/users", func(c *gin.Context) {
		var user gin.H
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		c.Header("X-User-ID", "123")
		c.JSON(http.StatusCreated, gin.H{
			"id":   123,
			"name": user["name"],
		})
	})
	
	engine.PUT("/api/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"id":      id,
			"updated": true,
		})
	})
	
	t.Run("GET请求不影响业务逻辑", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/users", nil)
		req.Header.Set("Origin", "https://example.com")
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		
		// 验证HTTP状态码 (Verify HTTP status code)
		assert.Equal(t, http.StatusOK, w.Code)
		
		// 验证CORS头 (Verify CORS headers)
		assert.Equal(t, "https://example.com", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "X-Total-Count,X-User-ID", w.Header().Get("Access-Control-Expose-Headers"))
		
		// 验证自定义中间件仍然工作 (Verify custom middleware still works)
		assert.Equal(t, "applied", w.Header().Get("X-Custom-Middleware"))
		assert.Equal(t, "100", w.Header().Get("X-Total-Count"))
		
		// 验证响应内容 (Verify response content)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "users")
	})
	
	t.Run("POST请求正常处理JSON", func(t *testing.T) {
		payload := `{"name": "Charlie"}`
		req := httptest.NewRequest("POST", "/api/users", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Origin", "https://app.com")
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		
		// 验证HTTP状态码 (Verify HTTP status code)
		assert.Equal(t, http.StatusCreated, w.Code)
		
		// 验证CORS头 (Verify CORS headers)
		assert.Equal(t, "https://app.com", w.Header().Get("Access-Control-Allow-Origin"))
		
		// 验证业务响应头 (Verify business response headers)
		assert.Equal(t, "123", w.Header().Get("X-User-ID"))
		
		// 验证响应内容 (Verify response content)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Charlie", response["name"])
	})
	
	t.Run("PUT请求路径参数正常解析", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/api/users/456", nil)
		req.Header.Set("Origin", "https://example.com")
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		
		// 验证HTTP状态码 (Verify HTTP status code)
		assert.Equal(t, http.StatusOK, w.Code)
		
		// 验证CORS头 (Verify CORS headers)
		assert.Equal(t, "https://example.com", w.Header().Get("Access-Control-Allow-Origin"))
		
		// 验证路径参数解析 (Verify path parameter parsing)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "456", response["id"])
		assert.True(t, response["updated"].(bool))
	})
	
	t.Run("预检请求不影响业务路由", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/api/users", nil)
		req.Header.Set("Origin", "https://example.com")
		req.Header.Set("Access-Control-Request-Method", "POST")
		req.Header.Set("Access-Control-Request-Headers", "Content-Type")
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		
		// 验证预检请求返回204 (Verify preflight request returns 204)
		assert.Equal(t, http.StatusNoContent, w.Code)
		
		// 验证CORS预检头 (Verify CORS preflight headers)
		assert.Equal(t, "https://example.com", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
		
		// 验证响应体为空（预检请求不应执行业务逻辑）(Verify empty response body - preflight should not execute business logic)
		assert.Empty(t, w.Body.String())
	})
	
	t.Run("非允许来源被正确拒绝", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/users", nil)
		req.Header.Set("Origin", "https://malicious.com")
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
		
		// 验证请求被拒绝 (Verify request is rejected)
		assert.Equal(t, http.StatusForbidden, w.Code)
		
		// 验证没有CORS头 (Verify no CORS headers)
		assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
		
		// 验证业务逻辑没有执行 (Verify business logic was not executed)
		assert.Empty(t, w.Header().Get("X-Total-Count"))
	})
	
	t.Run("中间件链顺序正确", func(t *testing.T) {
		// 创建新引擎测试中间件顺序 (Create new engine to test middleware order)
		testEngine := gin.Default()
		
		// 添加一个记录中间件执行顺序的中间件 (Add middleware to record execution order)
		var executionOrder []string
		
		testEngine.Use(func(c *gin.Context) {
			executionOrder = append(executionOrder, "before-cors")
			c.Next()
		})
		
		testEngine.Use(corsMiddleware.GetGinHandler())
		
		testEngine.Use(func(c *gin.Context) {
			executionOrder = append(executionOrder, "after-cors")
			c.Next()
		})
		
		testEngine.GET("/test", func(c *gin.Context) {
			executionOrder = append(executionOrder, "handler")
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})
		
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "https://example.com")
		w := httptest.NewRecorder()
		testEngine.ServeHTTP(w, req)
		
		// 验证中间件执行顺序 (Verify middleware execution order)
		expectedOrder := []string{"before-cors", "after-cors", "handler"}
		assert.Equal(t, expectedOrder, executionOrder)
		
		// 验证请求正常处理 (Verify request processed normally)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "https://example.com", w.Header().Get("Access-Control-Allow-Origin"))
	})
}

// TestCORSPerformanceImpact 测试CORS中间件对性能的影响 (Test performance impact of CORS middleware)
func TestCORSPerformanceImpact(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	config := &server.CORSConfig{
		Enabled:      true,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST"},
		MaxAge:       1 * time.Hour,
	}
	
	corsMiddleware := NewCORSMiddleware(config).(*CORSMiddleware)
	
	// 测试没有CORS中间件的性能 (Test performance without CORS middleware)
	engineWithoutCORS := gin.Default()
	engineWithoutCORS.GET("/bench", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "benchmark"})
	})
	
	// 测试有CORS中间件的性能 (Test performance with CORS middleware)
	engineWithCORS := gin.Default()
	engineWithCORS.Use(corsMiddleware.GetGinHandler())
	engineWithCORS.GET("/bench", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "benchmark"})
	})
	
	req := httptest.NewRequest("GET", "/bench", nil)
	req.Header.Set("Origin", "https://example.com")
	
	// 简单性能对比（不是正式基准测试）(Simple performance comparison - not formal benchmark)
	t.Run("响应时间对比", func(t *testing.T) {
		// 测试没有CORS的响应 (Test response without CORS)
		w1 := httptest.NewRecorder()
		start1 := time.Now()
		engineWithoutCORS.ServeHTTP(w1, req)
		duration1 := time.Since(start1)
		
		// 测试有CORS的响应 (Test response with CORS)
		w2 := httptest.NewRecorder()
		start2 := time.Now()
		engineWithCORS.ServeHTTP(w2, req)
		duration2 := time.Since(start2)
		
		// 验证两个响应都成功 (Verify both responses succeed)
		assert.Equal(t, http.StatusOK, w1.Code)
		assert.Equal(t, http.StatusOK, w2.Code)
		
		// CORS中间件的开销应该很小（通常小于10倍）(CORS middleware overhead should be minimal - usually less than 10x)
		assert.True(t, duration2 < duration1*10, 
			"CORS middleware overhead too high: %v vs %v", duration2, duration1)
		
		// 验证CORS头已添加 (Verify CORS headers added)
		assert.Empty(t, w1.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "*", w2.Header().Get("Access-Control-Allow-Origin"))
	})
} 