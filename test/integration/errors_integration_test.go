/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package integration

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"testing"

	lmccerrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestErrorsIntegration 错误模块集成测试套件
// (TestErrorsIntegration errors module integration test suite)

// mockService 模拟服务，用于测试复杂的错误场景
// (mockService mock service for testing complex error scenarios)
type mockService struct {
	name string
}

func (s *mockService) ProcessRequest(requestID string) error {
	if requestID == "" {
		return lmccerrors.NewWithCode(
			lmccerrors.ErrValidation,
			"request ID cannot be empty",
		)
	}
	
	if requestID == "invalid" {
		return lmccerrors.ErrorfWithCode(
			lmccerrors.ErrBadRequest,
			"invalid request ID format: %s",
			requestID,
		)
	}
	
	if requestID == "fail" {
		return s.simulateDeepError()
	}
	
	return nil
}

func (s *mockService) simulateDeepError() error {
	err := s.databaseOperation()
	if err != nil {
		return lmccerrors.WithCode(
			lmccerrors.Wrapf(err, "failed to process request in service %s", s.name),
			lmccerrors.ErrInternalServer,
		)
	}
	return nil
}

func (s *mockService) databaseOperation() error {
	err := s.lowLevelDbCall()
	if err != nil {
		return lmccerrors.Wrap(err, "database operation failed")
	}
	return nil
}

func (s *mockService) lowLevelDbCall() error {
	return lmccerrors.New("connection timeout")
}

// TestErrorCodeSystem 测试错误码系统的完整性
// (TestErrorCodeSystem tests the completeness of error code system)
func TestErrorCodeSystem(t *testing.T) {
	tests := []struct {
		name         string
		coder        lmccerrors.Coder
		expectedCode int
		expectedHTTP int
		hasReference bool
	}{
		{
			name:         "Internal Server Error",
			coder:        lmccerrors.ErrInternalServer,
			expectedCode: 100001,
			expectedHTTP: http.StatusInternalServerError,
			hasReference: false,
		},
		{
			name:         "Not Found Error",
			coder:        lmccerrors.ErrNotFound,
			expectedCode: 100002,
			expectedHTTP: http.StatusNotFound,
			hasReference: false,
		},
		{
			name:         "Bad Request Error",
			coder:        lmccerrors.ErrBadRequest,
			expectedCode: 100003,
			expectedHTTP: http.StatusBadRequest,
			hasReference: false,
		},
		{
			name:         "Validation Error",
			coder:        lmccerrors.ErrValidation,
			expectedCode: 100006,
			expectedHTTP: http.StatusBadRequest,
			hasReference: false,
		},
		{
			name:         "Config File Read Error",
			coder:        lmccerrors.ErrConfigFileRead,
			expectedCode: 200001,
			expectedHTTP: http.StatusInternalServerError,
			hasReference: true,
		},
		{
			name:         "Config Setup Error",
			coder:        lmccerrors.ErrConfigSetup,
			expectedCode: 200002,
			expectedHTTP: http.StatusInternalServerError,
			hasReference: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 测试错误码 (Test error code)
			assert.Equal(t, tt.expectedCode, tt.coder.Code(),
				"Error code should match expected value")

			// 测试HTTP状态码 (Test HTTP status code)
			assert.Equal(t, tt.expectedHTTP, tt.coder.HTTPStatus(),
				"HTTP status code should match expected value")

			// 测试错误描述 (Test error description)
			assert.NotEmpty(t, tt.coder.String(),
				"Error description should not be empty")

			// 测试参考文档 (Test reference documentation)
			reference := tt.coder.Reference()
			if tt.hasReference {
				assert.NotEmpty(t, reference,
					"Error should have reference documentation")
			}

			// 测试错误实现了error接口 (Test error implements error interface)
			assert.NotEmpty(t, tt.coder.Error(),
				"Error should implement error interface")
		})
	}
}

// createErrorWithCode 创建带错误码的错误的辅助函数 (Helper function to create error with code)
func createErrorWithCode() error {
	return lmccerrors.NewWithCode(
		lmccerrors.ErrValidation,
		"validation failed",
	)
}

// createFormattingTestError 创建用于格式化测试的错误的辅助函数 (Helper function to create error for formatting test)
func createFormattingTestError() error {
	return lmccerrors.NewWithCode(
		lmccerrors.ErrValidation,
		"validation failed for field 'name'",
	)
}

// TestErrorStackTraces 测试错误堆栈跟踪功能
// (TestErrorStackTraces tests error stack trace functionality)
func TestErrorStackTraces(t *testing.T) {
	// 测试基础错误创建 (Test basic error creation)
	t.Run("Basic Error Creation", func(t *testing.T) {
		err := lmccerrors.New("test error message")
		require.Error(t, err)
		
		stackTrace := fmt.Sprintf("%+v", err)
		assert.Contains(t, stackTrace, "test error message")
		assert.Contains(t, stackTrace, "TestErrorStackTraces") // 应该包含当前函数名 (Should contain current function name)
		assert.Contains(t, stackTrace, ".go:") // 应该包含文件名和行号 (Should contain filename and line number)
	})

	// 测试错误包装 (Test error wrapping)
	t.Run("Error Wrapping", func(t *testing.T) {
		originalErr := lmccerrors.New("original error")
		wrappedErr := lmccerrors.Wrap(originalErr, "wrapped error")
		
		stackTrace := fmt.Sprintf("%+v", wrappedErr)
		assert.Contains(t, stackTrace, "original error")
		assert.Contains(t, stackTrace, "wrapped error")
		
		// 测试错误解包 (Test error unwrapping)
		unwrappedErr := lmccerrors.Cause(wrappedErr)
		assert.Equal(t, originalErr, unwrappedErr)
	})

	// 测试带错误码的错误 (Test errors with codes)
	t.Run("Errors With Codes", func(t *testing.T) {
		err := createErrorWithCode() // 通过辅助函数创建错误以确保正确的堆栈跟踪 (Create error through helper function to ensure correct stack trace)
		
		stackTrace := fmt.Sprintf("%+v", err)
		assert.Contains(t, stackTrace, "validation failed")
		// 堆栈跟踪应包含辅助函数名而不是测试函数名 (Stack trace should contain helper function name instead of test function name)
		assert.Contains(t, stackTrace, "createErrorWithCode")
		
		// 测试错误码提取 (Test error code extraction)
		coder := lmccerrors.GetCoder(err)
		require.NotNil(t, coder)
		assert.Equal(t, lmccerrors.ErrValidation.Code(), coder.Code())
	})
}

// TestErrorWrappingChains 测试错误包装链
// (TestErrorWrappingChains tests error wrapping chains)
func TestErrorWrappingChains(t *testing.T) {
	service := &mockService{name: "TestService"}
	
	// 模拟深层错误传播 (Simulate deep error propagation)
	err := service.ProcessRequest("fail")
	require.Error(t, err)
	
	// 验证错误包装链 (Verify error wrapping chain)
	stackTrace := fmt.Sprintf("%+v", err)
	
	// 应该包含所有层级的错误信息 (Should contain error messages from all levels)
	assert.Contains(t, stackTrace, "connection timeout")
	assert.Contains(t, stackTrace, "database operation failed")
	assert.Contains(t, stackTrace, "failed to process request in service TestService")
	
	// 应该包含正确的错误码 (Should contain correct error code)
	coder := lmccerrors.GetCoder(err)
	require.NotNil(t, coder)
	assert.Equal(t, lmccerrors.ErrInternalServer.Code(), coder.Code())
	
	// 测试错误链的遍历 (Test error chain traversal)
	rootCause := lmccerrors.Cause(err)
	assert.Contains(t, rootCause.Error(), "connection timeout")
}

// TestErrorGroupFunctionality 测试错误分组功能
// (TestErrorGroupFunctionality tests error group functionality)
func TestErrorGroupFunctionality(t *testing.T) {
	// 测试空错误组 (Test empty error group)
	t.Run("Empty Error Group", func(t *testing.T) {
		eg := lmccerrors.NewErrorGroup()
		assert.Len(t, eg.Errors(), 0)
		assert.Equal(t, "no errors", eg.Error())
	})

	// 测试单个错误 (Test single error)
	t.Run("Single Error", func(t *testing.T) {
		eg := lmccerrors.NewErrorGroup("operation failed")
		err := lmccerrors.New("test error")
		eg.Add(err)
		
		errors := eg.Errors()
		assert.Len(t, errors, 1)
		assert.Equal(t, err, errors[0])
		assert.Contains(t, eg.Error(), "operation failed")
		assert.Contains(t, eg.Error(), "test error")
	})

	// 测试多个错误 (Test multiple errors)
	t.Run("Multiple Errors", func(t *testing.T) {
		eg := lmccerrors.NewErrorGroup("batch operation failed")
		
		// 添加多个不同类型的错误 (Add multiple different types of errors)
		eg.Add(lmccerrors.NewWithCode(lmccerrors.ErrValidation, "validation error"))
		eg.Add(lmccerrors.New("simple error"))
		eg.Add(lmccerrors.ErrorfWithCode(lmccerrors.ErrBadRequest, "parameter %s is invalid", "id"))
		
		errors := eg.Errors()
		assert.Len(t, errors, 3)
		
		errorMessage := eg.Error()
		assert.Contains(t, errorMessage, "batch operation failed")
		assert.Contains(t, errorMessage, "validation error")
		assert.Contains(t, errorMessage, "simple error")
		assert.Contains(t, errorMessage, "parameter id is invalid")
	})

	// 测试错误组的格式化输出 (Test error group formatting)
	t.Run("Error Group Formatting", func(t *testing.T) {
		eg := lmccerrors.NewErrorGroup("formatting test")
		eg.Add(lmccerrors.New("error 1"))
		eg.Add(lmccerrors.New("error 2"))
		
		// 测试详细格式化 (Test detailed formatting)
		detailedOutput := fmt.Sprintf("%+v", eg)
		assert.Contains(t, detailedOutput, "formatting test")
		assert.Contains(t, detailedOutput, "error 1")
		assert.Contains(t, detailedOutput, "error 2")
	})
}

// TestConcurrentErrorHandling 测试并发错误处理
// (TestConcurrentErrorHandling tests concurrent error handling)
func TestConcurrentErrorHandling(t *testing.T) {
	numGoroutines := 100
	eg := lmccerrors.NewErrorGroup("concurrent operation")
	service := &mockService{name: "ConcurrentService"}
	
	var wg sync.WaitGroup
	var mu sync.Mutex
	
	wg.Add(numGoroutines)
	
	// 并发生成错误 (Generate errors concurrently)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			
			var err error
			switch id % 4 {
			case 0:
				err = service.ProcessRequest("")
			case 1:
				err = service.ProcessRequest("invalid")
			case 2:
				err = service.ProcessRequest("fail")
			case 3:
				err = nil // 成功情况 (Success case)
			}
			
			if err != nil {
				mu.Lock()
				eg.Add(lmccerrors.Wrapf(err, "goroutine %d failed", id))
				mu.Unlock()
			}
		}(i)
	}
	
	wg.Wait()
	
	// 验证错误收集 (Verify error collection)
	errors := eg.Errors()
	expectedErrors := numGoroutines * 3 / 4 // 75% 的情况会产生错误 (75% cases will generate errors)
	assert.Equal(t, expectedErrors, len(errors))
	
	// 验证错误类型分布 (Verify error type distribution)
	validationErrors := 0
	badRequestErrors := 0
	internalErrors := 0
	
	for _, err := range errors {
		coder := lmccerrors.GetCoder(err)
		if coder != nil {
			switch coder.Code() {
			case lmccerrors.ErrValidation.Code():
				validationErrors++
			case lmccerrors.ErrBadRequest.Code():
				badRequestErrors++
			case lmccerrors.ErrInternalServer.Code():
				internalErrors++
			}
		}
	}
	
	expectedEachType := numGoroutines / 4
	assert.Equal(t, expectedEachType, validationErrors)
	assert.Equal(t, expectedEachType, badRequestErrors)
	assert.Equal(t, expectedEachType, internalErrors)
}

// TestErrorFormatting 测试错误格式化输出
// (TestErrorFormatting tests error formatting output)
func TestErrorFormatting(t *testing.T) {
	// 测试简单错误格式化 (Test simple error formatting)
	t.Run("Simple Error Formatting", func(t *testing.T) {
		err := lmccerrors.New("simple error")
		
		// %s 格式 (Simple string format)
		simpleStr := fmt.Sprintf("%s", err)
		assert.Equal(t, "simple error", simpleStr)
		
		// %v 格式 (Default format)
		defaultStr := fmt.Sprintf("%v", err)
		assert.Equal(t, "simple error", defaultStr)
		
		// %+v 格式应包含堆栈信息 (%+v format should include stack info)
		detailedStr := fmt.Sprintf("%+v", err)
		assert.Contains(t, detailedStr, "simple error")
		assert.Contains(t, detailedStr, "TestErrorFormatting")
	})

	// 测试带错误码的错误格式化 (Test error with code formatting)
	t.Run("Error With Code Formatting", func(t *testing.T) {
		err := createFormattingTestError() // 通过辅助函数创建错误以确保正确的堆栈跟踪 (Create error through helper function to ensure correct stack trace)
		
		// 简单格式应包含错误码信息 (Simple format should include error code info)
		simpleStr := fmt.Sprintf("%s", err)
		assert.Contains(t, simpleStr, "validation failed for field 'name'")
		assert.Contains(t, simpleStr, lmccerrors.ErrValidation.String())
		
		// 详细格式应包含堆栈信息 (Detailed format should include stack info)
		detailedStr := fmt.Sprintf("%+v", err)
		assert.Contains(t, detailedStr, "validation failed for field 'name'")
		assert.Contains(t, detailedStr, lmccerrors.ErrValidation.String())
		// 堆栈跟踪应包含辅助函数名 (Stack trace should contain helper function name)
		assert.Contains(t, detailedStr, "createFormattingTestError")
	})

	// 测试错误包装格式化 (Test error wrapping formatting)
	t.Run("Error Wrapping Formatting", func(t *testing.T) {
		originalErr := lmccerrors.NewWithCode(lmccerrors.ErrBadRequest, "invalid user ID")
		wrappedErr := lmccerrors.Wrap(originalErr, "user validation failed")
		doubleWrappedErr := lmccerrors.WithCode(
			lmccerrors.Wrapf(wrappedErr, "request processing failed for endpoint %s", "/api/users"),
			lmccerrors.ErrInternalServer,
		)
		
		detailedStr := fmt.Sprintf("%+v", doubleWrappedErr)
		
		// 应该包含所有层级的错误信息 (Should contain error info from all levels)
		assert.Contains(t, detailedStr, "invalid user ID")
		assert.Contains(t, detailedStr, "user validation failed")
		assert.Contains(t, detailedStr, "request processing failed for endpoint /api/users")
		assert.Contains(t, detailedStr, lmccerrors.ErrInternalServer.String())
		
		// 应该包含多个堆栈跟踪点 (Should contain multiple stack trace points)
		stackFrameCount := strings.Count(detailedStr, ".go:")
		assert.Greater(t, stackFrameCount, 2, "Should have multiple stack frames")
	})
}

// TestErrorCodeCompatibility 测试错误码兼容性
// (TestErrorCodeCompatibility tests error code compatibility)
func TestErrorCodeCompatibility(t *testing.T) {
	// 测试IsCode函数 (Test IsCode function)
	t.Run("IsCode Function", func(t *testing.T) {
		err := lmccerrors.NewWithCode(lmccerrors.ErrValidation, "test validation error")
		
		assert.True(t, lmccerrors.IsCode(err, lmccerrors.ErrValidation))
		assert.False(t, lmccerrors.IsCode(err, lmccerrors.ErrBadRequest))
		assert.False(t, lmccerrors.IsCode(err, lmccerrors.ErrInternalServer))
		
		// 测试包装后的错误 (Test wrapped errors)
		wrappedErr := lmccerrors.Wrap(err, "wrapper message")
		assert.True(t, lmccerrors.IsCode(wrappedErr, lmccerrors.ErrValidation))
		
		// 测试nil错误 (Test nil error)
		assert.False(t, lmccerrors.IsCode(nil, lmccerrors.ErrValidation))
		
		// 测试标准错误 (Test standard error)
		stdErr := fmt.Errorf("standard error")
		assert.False(t, lmccerrors.IsCode(stdErr, lmccerrors.ErrValidation))
	})

	// 测试GetCoder函数 (Test GetCoder function)
	t.Run("GetCoder Function", func(t *testing.T) {
		err := lmccerrors.NewWithCode(lmccerrors.ErrBadRequest, "test parameter error")
		
		coder := lmccerrors.GetCoder(err)
		require.NotNil(t, coder)
		assert.Equal(t, lmccerrors.ErrBadRequest.Code(), coder.Code())
		
		// 测试包装后的错误 (Test wrapped errors)
		wrappedErr := lmccerrors.Wrap(err, "wrapper")
		wrappedCoder := lmccerrors.GetCoder(wrappedErr)
		require.NotNil(t, wrappedCoder)
		assert.Equal(t, lmccerrors.ErrBadRequest.Code(), wrappedCoder.Code())
		
		// 测试没有错误码的错误 (Test error without code)
		simpleErr := lmccerrors.New("simple error")
		simpleCoder := lmccerrors.GetCoder(simpleErr)
		assert.Nil(t, simpleCoder)
		
		// 测试nil错误 (Test nil error)
		nilCoder := lmccerrors.GetCoder(nil)
		assert.Nil(t, nilCoder)
	})
}

// TestMemoryUsageAndPerformance 测试内存使用和性能
// (TestMemoryUsageAndPerformance tests memory usage and performance)
func TestMemoryUsageAndPerformance(t *testing.T) {
	// 性能测试：创建大量错误 (Performance test: create many errors)
	t.Run("Error Creation Performance", func(t *testing.T) {
		const numErrors = 10000
		
		// 测量内存使用前的状态 (Measure memory usage before)
		var memStatsBefore runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&memStatsBefore)
		
		// 创建大量错误 (Create many errors)
		errors := make([]error, numErrors)
		for i := 0; i < numErrors; i++ {
			errors[i] = lmccerrors.NewWithCode(
				lmccerrors.ErrValidation,
				fmt.Sprintf("error %d", i),
			)
		}
		
		// 测量内存使用后的状态 (Measure memory usage after)
		var memStatsAfter runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&memStatsAfter)
		
		// 验证错误数量 (Verify error count)
		assert.Len(t, errors, numErrors)
		
		// 验证每个错误都是有效的 (Verify each error is valid)
		for i, err := range errors {
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), fmt.Sprintf("error %d", i))
			assert.True(t, lmccerrors.IsCode(err, lmccerrors.ErrValidation))
		}
		
		// 计算内存增长 (Calculate memory growth)
		memoryGrowth := memStatsAfter.HeapAlloc - memStatsBefore.HeapAlloc
		averageMemoryPerError := memoryGrowth / numErrors
		
		// 验证内存使用合理 (Verify reasonable memory usage)
		// 每个错误的平均内存使用应该小于1KB (Average memory per error should be less than 1KB)
		assert.Less(t, averageMemoryPerError, uint64(1024),
			"Average memory per error should be reasonable")
		
		t.Logf("Created %d errors, memory growth: %d bytes, average per error: %d bytes",
			numErrors, memoryGrowth, averageMemoryPerError)
	})

	// 测试错误包装性能 (Test error wrapping performance)
	t.Run("Error Wrapping Performance", func(t *testing.T) {
		const wrapLevels = 100
		
		// 创建基础错误 (Create base error)
		baseErr := lmccerrors.New("base error")
		err := baseErr
		
		// 多层包装 (Multiple levels of wrapping)
		for i := 0; i < wrapLevels; i++ {
			err = lmccerrors.Wrapf(err, "wrapper level %d", i)
		}
		
		// 验证包装深度 (Verify wrapping depth)
		stackTrace := fmt.Sprintf("%+v", err)
		wrapperCount := strings.Count(stackTrace, "wrapper level")
		assert.Equal(t, wrapLevels, wrapperCount)
		
		// 验证根本原因提取性能 (Verify root cause extraction performance)
		rootCause := lmccerrors.Cause(err)
		assert.Equal(t, baseErr, rootCause)
	})
} 