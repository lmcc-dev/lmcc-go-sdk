<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

## 从标准库 `errors` 迁移 (Migration from Standard Library `errors`)

从 Go 的标准库 `errors` 包迁移到 `pkg/errors` 可以逐步进行。`pkg/errors` 模块的设计在很大程度上是兼容的，特别是在错误包装以及通过 `standardErrors.Is` 和 `standardErrors.As` 进行错误检查方面。

(Migrating from Go\'s standard library `errors` package to `pkg/errors` can be done incrementally. The `pkg/errors` module is designed to be largely compatible, especially with error wrapping and checking via `standardErrors.Is` and `standardErrors.As`.)

迁移的主要好处包括：

(Key benefits of migrating include:)

- **堆栈跟踪 (Stack Traces)**: 创建错误时自动捕获堆栈跟踪。
  (Automatic stack trace capture on error creation.)
- **错误码 (`Coder`) (Error Codes (`Coder`))**: 标准化的方式对错误进行分类，并将其与 HTTP 状态或其他元数据相关联。
  (Standardized way to categorize errors and associate them with HTTP statuses or other metadata.)
- **结构化格式化 (Structured Formatting)**: 使用 `%+v` 获得更丰富的错误消息。
  (Richer error messages with `%+v`.)
- **错误聚合 (Error Aggregation)**: 用于收集多个错误的 `ErrorGroup`。
  (`ErrorGroup` for collecting multiple errors.) 