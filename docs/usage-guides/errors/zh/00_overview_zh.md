<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

## 概述 (Overview)

`pkg/errors` 模块为 Go 应用程序提供了一种强大且结构化的错误处理方法。它通过以下方式扩展了标准库：

- **自动堆栈跟踪 (Automatic stack traces)**，便于调试
- **错误码 (Error codes)**，用于程序化错误处理
- 通过包装实现**丰富的错误上下文 (Rich error context through wrapping)**
- **错误聚合 (Error aggregation)**，用于收集多个故障
- **与标准库兼容 (Compatibility with standard library)** (`errors.Is`, `errors.As`) 