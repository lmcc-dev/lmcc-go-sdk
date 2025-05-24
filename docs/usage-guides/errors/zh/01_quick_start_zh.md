<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

## 快速入门 (Quick Start)

替换您现有的错误处理方式：

```go
// 之前 (Before)：
import "errors"
// import "fmt" // 如果使用 fmt.Errorf

func doSomethingOld(fail bool) error {
    if fail {
        return errors.New("操作失败 (old)")
        // 或者 (or): return fmt.Errorf("操作 %s 失败 (old)", "something") 
    }
    return nil
}

// 之后 (After)：
import "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"

func doSomethingNew(fail bool) error {
    if fail {
        // 现在带有堆栈跟踪！(Now with stack trace!)
        return errors.New("操作失败 (new)") 
        // 或者 (or): return errors.Errorf("操作 %s 失败 (new)", "something") 
    }
    return nil
}

func main() {
    // 演示旧方法 (Demonstrating old method)
    errOld := doSomethingOld(true)
    if errOld != nil {
        fmt.Printf("旧错误 (Old error): %v\n", errOld)
        // fmt.Printf("旧错误详情 (Old error details %%+v): %+v\n", errOld) // 标准库 errors.New 不会自动提供堆栈跟踪给 %+v (Standard library errors.New doesn't automatically provide stack trace to %+v)
    }

    fmt.Println("---")

    // 演示新方法 (Demonstrating new method)
    errNew := doSomethingNew(true)
    if errNew != nil {
        fmt.Printf("新错误 (New error): %v\n", errNew)
        fmt.Printf("新错误详情 (New error details %%+v):\n%+v\n", errNew) // pkg/errors 会提供堆栈跟踪 (pkg/errors provides stack trace)
    }
}

/*
预期输出 (Expected Output - 堆栈跟踪路径和行号会变化 (Stack trace paths and line numbers will vary)):

旧错误 (Old error): 操作失败 (old)
---
新错误 (New error): 操作失败 (new)
新错误详情 (New error details %%+v):
操作失败 (new)
main.doSomethingNew
	/path/to/your/file.go:XX
main.main
	/path/to/your/file.go:YY
runtime.main
	...
runtime.goexit
	...
*/
```

**主要变化 (Key Change):**
只需将标准库 `errors` (以及 `fmt.Errorf` 用于错误创建的场景) 替换为 `github.com/lmcc-dev/lmcc-go-sdk/pkg/errors`，即可自动获得带堆栈跟踪的错误对象。

(The primary change is simply replacing the standard library `errors` (and `fmt.Errorf` for error creation) with `github.com/lmcc-dev/lmcc-go-sdk/pkg/errors` to automatically get errors with stack traces.) 