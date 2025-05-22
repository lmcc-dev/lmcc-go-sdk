/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package errors

import (
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"
)

// Frame represents a program counter inside a stack trace.
// Frame 表示堆栈跟踪中的一个程序计数器。
// For historical reasons if Frame is interpreted as a uintptr
// it represents the program counter + 1.
// 由于历史原因，如果将 Frame 解释为 uintptr，它表示程序计数器 + 1。
type Frame uintptr

// pc returns the program counter for this frame.
// pc 返回此帧的程序计数器。
func (f Frame) pc() uintptr { return uintptr(f) - 1 }

// file returns the full path to the file that contains the
// function for this Frame's pc.
// file 返回包含此 Frame 的 pc 所在函数的文件的完整路径。
func (f Frame) file() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	file, _ := fn.FileLine(f.pc())
	return file
}

// line returns the line number of source code of the
// function for this Frame's pc.
// line 返回此 Frame 的 pc 所在函数的源代码的行号。
func (f Frame) line() int {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(f.pc())
	return line
}

// name returns the name of this function, if known.
// name 返回此函数的名称 (如果已知)。
func (f Frame) name() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	return fn.Name()
}

// StackTrace is a stack of Frames from innermost (newest) to outermost (oldest).
// StackTrace 是一个从最内层 (最新) 到最外层 (最旧) 的 Frame 堆栈。
type StackTrace []Frame

const (
	maxStackDepth = 32
	// skipFrames    = 3 // Default skip for callers, New, Errorf etc. // 已在 errors.go 中定义 (Defined in errors.go)
)

// callers retrieves the current call stack.
// callers 检索当前的调用堆栈。
// It skips a number of frames specified by the 'skip' argument.
// 它会跳过 'skip' 参数指定的帧数。
func callers(skip int) StackTrace {
	pc := make([]uintptr, maxStackDepth)
	n := runtime.Callers(skip, pc)
	if n == 0 {
		return nil
	}

	var st StackTrace
	frames := runtime.CallersFrames(pc[:n])
	for {
		frame, more := frames.Next()
		// To avoid runtime errors from reporting non-Go sources we skip frames from non-Go files.
		// 为避免报告非 Go 源的运行时错误，我们跳过来自非 Go 文件的帧。
		if !strings.HasSuffix(frame.File, ".go") {
			if !more {
				break
			}
			continue
		}
		st = append(st, Frame(frame.PC))
		if !more {
			break
		}
	}
	return st
}

// Format formats the stack trace from the point of error creation.
// Format 格式化从错误创建点开始的堆栈跟踪。
// It will be ignored by `fmt.Fprintf` if the format string
// indicates that the result of `Error()` is to be used, e.g. `%s` or `%v`.
// 如果格式字符串指示要使用 `Error()` 的结果 (例如 `%s` 或 `%v`)，`fmt.Fprintf` 将忽略它。
// To print the stack trace, use format specifier `%+v` when printing the error.
// 要打印堆栈跟踪，请在打印错误时使用格式说明符 `%+v`。
//
// The output for each frame of the stack trace will be:
//   <function_name>
// 	<file>:<line>
// (堆栈跟踪的每个帧的输出将是：
//   <函数名>
// 	<文件>:<行号>)
func (st StackTrace) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			for _, f := range st {
				// Note: Using io.WriteString for potentially better performance
				// and to avoid issues if frame components contain formatting verbs.
				// 注意：为了潜在的性能提升和避免帧组件包含格式化动词时可能出现的问题，这里使用 io.WriteString。
				_, _ = io.WriteString(s, "\n")
				_, _ = io.WriteString(s, f.name())
				_, _ = io.WriteString(s, "\n\t")
				_, _ = io.WriteString(s, f.file())
				_, _ = io.WriteString(s, ":")
				_, _ = io.WriteString(s, strconv.Itoa(f.line()))
			}
		}
	}
} 