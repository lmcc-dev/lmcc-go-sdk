/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package errors

import (
	"strings"
	"testing"
)

// Helper function to get a stack trace from a specific depth.
// 辅助函数，用于从特定深度获取堆栈跟踪。
func getTestStack(skip int) StackTrace {
	// skip: number of additional frames to skip besides getTestStack, callers, and runtime.Callers
	// We want the final skip for runtime.Callers to be 3 + skip
	// 3 to skip runtime.Callers, errors.callers, and errors.getTestStack
	return callers(skip + 3) 
}

func TestCallers(t *testing.T) {
	// This test is a bit tricky as file paths and line numbers can change.
	// We primarily check if we get a non-empty stack and if the top frames are reasonable.
	// 这个测试有点棘手，因为文件路径和行号可能会改变。
	// 我们主要检查是否获得非空堆栈以及顶层帧是否合理。

	stack := getTestStack(0) // We want the caller of getTestStack, which is TestCallers

	if stack == nil || len(stack) == 0 {
		t.Fatal("callers() returned nil or empty stack")
	}

	// Check the first frame (most recent call site within our package, excluding callers itself)
	// The actual function name might be tricky if this test is inlined or optimized.
	// Let's look for the current test function name or file.
	topFrame := stack[0]
	funcName := topFrame.name()
	fileName := topFrame.file()
	lineNumber := topFrame.line()

	t.Logf("Top frame (expected TestCallers): %s @ %s:%d", funcName, fileName, lineNumber)

	if !strings.Contains(funcName, "TestCallers") { // It might be TestCallers.funcN due to t.Run
		t.Errorf("Top frame function name unexpected: got %s, expected to contain TestCallers", funcName)
	}
	if !strings.HasSuffix(fileName, "stack_test.go") {
		t.Errorf("Top frame file name unexpected: got %s, want suffix stack_test.go", fileName)
	}
	if lineNumber == 0 {
		t.Errorf("Top frame line number is 0, expected a valid line number")
	}

	// Test with a different skip value
	// Note: The skip value here is relative to the direct call to callers()
	// So, if we call callers(0) from getTestStack, and getTestStack is called by TestCallers,
	// the first frame would be runtime.Callers, then getTestStack, then TestCallers.
	// Let's test callers(0) from a deeper function.
	func() {
		func() {
			stackInner := callers(0) // Skip 0 frames from this point (runtime.Callers is frame 0)
			if len(stackInner) < 3 { // Should have at least runtime.Callers, errors.callers, and this func
				t.Fatalf("callers(0) from inner func returned too few frames: %d, expected at least 3", len(stackInner))
			}
			
			// Frame 0: runtime.Callers
			// Frame 1: errors.callers
			// Frame 2: This anonymous function
			frame1 := stackInner[1] // errors.callers
			if !strings.Contains(frame1.name(), "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors.callers") {
				t.Errorf("callers(0) second frame (index 1) name expected to be errors.callers, got %s", frame1.name())
			}
			if !strings.HasSuffix(frame1.file(), "stack.go") {
			    t.Errorf("callers(0) second frame (index 1) file name unexpected: got %s, want suffix stack.go", frame1.file())
			}

			frame2 := stackInner[2] // This anonymous function
			if !strings.Contains(frame2.name(), "TestCallers.func") { // Go naming for closures
				t.Errorf("callers(0) third frame (index 2) name expected to contain TestCallers.func, got %s", frame2.name())
			}
			if !strings.HasSuffix(frame2.file(), "stack_test.go") {
			    t.Errorf("callers(0) third frame (index 2) file name unexpected: got %s, want suffix stack_test.go", frame2.file())
			}
		}()
	}()
}

func TestFrameMethods(t *testing.T) {
	// Get a stack where the top frame is known to be inside stack_test.go
	// We call a function that then calls getTestStack(0)
	var frame Frame
	func() {
		stack := getTestStack(0) // Top frame should be this anonymous func
		if len(stack) > 0 {
			frame = stack[0]
		} else {
			t.Fatal("getTestStack(0) returned empty stack in TestFrameMethods")
		}
	}()

	file := frame.file()
	line := frame.line()
	name := frame.name()

	t.Logf("Frame for TestFrameMethods: %s @ %s:%d", name, file, line)

	if !strings.HasSuffix(file, "stack_test.go") {
		t.Errorf("Frame.file() = %s, want suffix stack_test.go", file)
	}
	if line == 0 {
		t.Errorf("Frame.line() = 0, want non-zero line number")
	}
	// The name will be something like ...pkg/errors_test.TestFrameMethods.func1
	// Note: Since this test file is `package errors`, the path will be `pkg/errors.TestFrameMethods.func...`
	if !strings.Contains(name, "TestFrameMethods") { 
		t.Errorf("Frame.name() = %s, want to contain TestFrameMethods", name)
	}

	// Test unknown frame
	unknownFrame := Frame(0) // Invalid PC, should result in "unknown"
	if unknownFrame.file() != "unknown" {
		t.Errorf("Frame(0).file() = %s, want unknown", unknownFrame.file())
	}
	if unknownFrame.line() != 0 {
		t.Errorf("Frame(0).line() = %d, want 0", unknownFrame.line())
	}
	if unknownFrame.name() != "unknown" {
		t.Errorf("Frame(0).name() = %s, want unknown", unknownFrame.name())
	}
}

// TestStackTraceFormat and its helpers (aTestFunctionForStackTrace, anotherTestFunction)
// have been migrated to format_test.go as TestStackTrace_Format. 