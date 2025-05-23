# lmcc-go-sdk Errors 模块规范

**版本:** 1.0.0
**作者:** Martin, AI Assistant

[English Version (英文版)](SPECIFICATION.md)

## 1. 引言

`lmcc-go-sdk` 中的 `pkg/errors` 模块为 Go 应用程序提供了一种灵活且结构化的错误处理方式。它通过引入用于分类错误信息的 `Coder` 系统、自动堆栈跟踪捕获、增强的格式化选项和错误聚合功能，扩展了标准库的错误处理能力。

该模块的主要目标是：
- 提供超越简单字符串消息的更丰富的错误上下文。
- 在应用程序的不同部分实现一致的错误处理和报告。
- 通过详细的堆栈跟踪和结构化的错误信息，方便调试。
- 与标准的 Go 错误处理习惯用法 (例如 `errors.Is`, `errors.As`) 保持兼容。

## 2. 核心概念

### 2.1. Coder (编码器)
`Coder` 是一个接口，代表一种特定类型的错误，具有唯一的代码、关联的 HTTP 状态、描述性字符串和可选的参考 URI。它允许根据错误分类进行程序化的检查和处理。

### 2.2. Stack Trace (堆栈跟踪)
模块在创建或包装错误时会自动捕获堆栈跟踪。此堆栈跟踪有助于查明错误源头。

### 2.3. Error Wrapping (错误包装)
可以包装错误以添加上下文信息。原始错误 ("cause") 会被保留并且可以解包，这与标准库中的 `errors.Unwrap` 兼容。

### 2.4. Error Aggregation (错误聚合)
可以将多个错误收集到单个 `ErrorGroup` 实例中。这对于一个操作可能导致多个需要一起报告的故障的场景非常有用。

## 3. 核心类型

### 3.1. `Coder` 接口

定义于 `pkg/errors/api.go`。

```go
type Coder interface {
    error // 嵌入标准错误接口

    Code() int
    String() string
    HTTPStatus() int
    Reference() string
}
```

- **`error`**: 嵌入 `error` 允许 `Coder` 实例直接用作标准错误。
- **`Code() int`**: 返回错误的唯一整数代码。
- **`String() string`**: 返回错误 Coder 的描述性字符串 (与 `Error()` 不同)。
- **`Error() string`**: (来自嵌入的 `error` 接口) 对于 `basicCoder` 实现，通常返回与 `String()` 相同的内容。
- **`HTTPStatus() int`**: 返回此错误建议的 HTTP 状态码 (例如 404, 500)。
- **`Reference() string`**: 返回有关此错误的更多信息的 URL 或文档链接。

主要实现是 `pkg/errors/coder.go` 中的 `basicCoder`。

### 3.2. `StackTrace` 类型

定义于 `pkg/errors/stack.go`。

```go
type StackTrace []Frame
type Frame uintptr
```

- **`Frame`**: 代表堆栈跟踪中的单个程序计数器。
- **`StackTrace`**: `Frame` 的切片，代表从最内层 (最新) 到最外层调用的调用堆栈。
- 它实现了 `fmt.Formatter` 接口，允许在使用 `%+v` 格式化错误时详细打印堆栈跟踪。

### 3.3. `fundamental` 错误类型

定义于 `pkg/errors/errors.go`。

这是模块创建的最基本的错误类型。它包含错误消息和堆栈跟踪。它不包装另一个错误。

- **字段**:
    - `msg string`: 错误消息。
    - `stack StackTrace`: 创建时捕获的堆栈跟踪。
- **关键方法**:
    - `Error() string`: 返回错误消息。
    - `Unwrap() error`: 返回 `nil`。
    - `Format(s fmt.State, verb rune)`: 实现 `fmt.Formatter` 接口。
        - `%s`, `%v`: 打印错误消息。
        - `%+v`: 打印错误消息和堆栈跟踪。

### 3.4. `wrapper` 错误类型

定义于 `pkg/errors/errors.go`。

此类型包装现有错误，添加新消息并在包装点捕获新的堆栈跟踪。

- **字段**:
    - `msg string`: 此包装器的附加消息。
    - `cause error`: 被包装的底层错误。
    - `stack StackTrace`: 包装点捕获的堆栈跟踪。
- **关键方法**:
    - `Error() string`: 返回 `msg + ": " + cause.Error()`。
    - `Unwrap() error`: 返回 `cause`。
    - `Format(s fmt.State, verb rune)`: 实现 `fmt.Formatter` 接口。
        - `%s`, `%v`: 打印组合的错误消息。
        - `%+v`: 打印组合的错误消息和此包装器的堆栈跟踪。

### 3.5. `withCode` 错误类型

定义于 `pkg/errors/errors.go`。

此类型将 `Coder` 与底层错误关联，并在关联点捕获堆栈跟踪。

- **字段**:
    - `cause error`: 底层错误。
    - `coder Coder`: 与此错误关联的 `Coder`。
    - `stack StackTrace`: 附加 `Coder` 时捕获的堆栈跟踪。
- **关键方法**:
    - `Error() string`: 返回一个结合了 `Coder` 字符串和 `cause` 错误消息的消息。
        - 如果 `cause` 不为 `nil`：格式为 `coder.String() + ": " + cause.Error()`。如果 `coder` 为 `nil` 或 `coder.String()` 为空，则返回 `cause.Error()`。
        - 如果 `cause` 为 `nil`：如果 `coder` 不为 `nil` 且 `coder.String()` 不为空，则返回 `coder.String()`。否则返回空字符串。(理想情况下，如果构造函数使用正确，则不应发生此情况)。
    - `Unwrap() error`: 返回 `cause`。
    - `Coder() Coder`: 返回关联的 `coder`。
    - `Format(s fmt.State, verb rune)`: 实现 `fmt.Formatter` 接口。
        - `%s`, `%v`: 打印包含 `Coder` 信息的错误消息。
        - `%+v`: 打印包含 `Coder` 信息的错误消息和此 `withCode` 包装器的堆栈跟踪。

### 3.6. `ErrorGroup` 类型

定义于 `pkg/errors/group.go`。

此类型允许将错误列表收集和管理为单个错误。

- **字段**:
    - `errs []error`: 错误切片。
    - `message string`: 错误组的可选总体消息。
- **关键方法**:
    - `Add(err error)`: 将非 `nil` 错误添加到组中。
    - `Errors() []error`: 返回包含的错误切片。
    - `Error() string`: 返回错误组的字符串表示形式。
        - 如果提供了组消息，它将作为消息的前缀。
        - 如果没有提供组消息且组中包含错误，则使用默认前缀 ("发生了一个错误: " 对于单个错误, "发生了多个错误: " 对于多个错误)。
        - 然后它列出所有包含的错误消息，以 "; " 分隔。
        - 如果组为空但有消息，则仅返回该消息。
        - 如果组为空且没有消息，则返回 "组内无错误"。
    - `Unwrap() []error`: 返回包含的错误切片，使其与 Go 1.20+ 的多错误解包兼容 (例如，用于 `errors.Is` 和 `errors.As`)。
    - `Format(s fmt.State, verb rune)`: 实现 `fmt.Formatter` 接口。
        - `%+v`: 打印组消息 (如果有)，然后是每个所含错误的详细、多行表示，如果可用并且使用 `%+v` 格式化，则包括其堆栈跟踪。如果组为空且没有消息，则打印 "空错误组"。如果组为空但有消息，则打印组消息后跟一个换行符。

## 4. 错误创建

### 4.1. `New(text string) error`
使用给定的文本消息创建一个新的 `fundamental` 错误，并捕获当前堆栈跟踪。

```go
err := errors.New("发生意外故障")
```

### 4.2. `Errorf(format string, args ...interface{}) error`
使用格式化的消息 (通过 `fmt.Sprintf`) 创建一个新的 `fundamental` 错误，并捕获当前堆栈跟踪。

```go
id := "123"
err := errors.Errorf("处理项目 %s 失败", id)
```

### 4.3. `Wrap(err error, message string) error`
使用附加 `message` 包装现有 `err`。它返回一个新的 `wrapper` 错误。如果 `err` 为 `nil`，`Wrap` 返回 `nil`。在包装点捕获新的堆栈跟踪。

```go
dbErr := // ... 来自数据库调用的某个错误
contextualErr := errors.Wrap(dbErr, "用户注册期间失败")
// contextualErr.Error() -> "用户注册期间失败: 原始 dbErr 消息"
```

### 4.4. `Wrapf(err error, format string, args ...interface{}) error`
使用附加的格式化 `message` 包装现有 `err`。它返回一个新的 `wrapper` 错误。如果 `err` 为 `nil`，`Wrapf` 返回 `nil`。捕获新的堆栈跟踪。

```go
id := "txn-007"
ioErr := // ... 某个 I/O 错误
contextualErr := errors.Wrapf(ioErr, "为事务 %s 写入数据失败", id)
```

### 4.5. `NewWithCode(coder Coder, text string) error`
创建一个将 `Coder` 与消息关联的新错误。底层错误 (`cause`) 是根据 `text` 创建的新的 `fundamental` 错误。返回一个 `withCode` 错误。如果 `coder` 为 `nil`，则使用预定义的 `unknownCoder`。

```go
var ErrUserNotFound = errors.NewCoder(10101, 404, "未找到用户", "")
err := errors.NewWithCode(ErrUserNotFound, "无法找到用户 'jane.doe'")
```

### 4.6. `ErrorfWithCode(coder Coder, format string, args ...interface{}) error`
创建一个将 `Coder` 与格式化消息关联的新错误。底层 `cause` 是一个新的 `fundamental` 错误。返回一个 `withCode` 错误。如果 `coder` 为 `nil`，则使用 `unknownCoder`。

```go
var ErrPaymentFailed = errors.NewCoder(20202, 400, "支付处理失败", "")
orderID := "order-456"
err := errors.ErrorfWithCode(ErrPaymentFailed, "由于资金不足，订单 %s 支付失败", orderID)
```

### 4.7. `WithCode(err error, coder Coder) error`
使用 `Coder` 注解现有 `err`。返回一个 `withCode` 错误。如果 `err` 为 `nil`，`WithCode` 返回 `nil`。如果 `coder` 为 `nil`，则使用 `unknownCoder`。

```go
var ErrPermissionDenied = errors.NewCoder(30303, 403, "权限被拒绝", "")
originalErr := errors.New("禁止访问资源")
errWithCode := errors.WithCode(originalErr, ErrPermissionDenied)
```

### 4.8. `WithMessage(err error, message string) error`
便捷函数，等同于 `Wrap(err, message)`。

### 4.9. `WithMessagef(err error, format string, args ...interface{}) error`
便捷函数，等同于 `Wrapf(err, format, args...)`。

### 4.10. `NewErrorGroup(message ...string) *ErrorGroup`
创建一个新的 `ErrorGroup`。可以为组提供一个可选的总体消息。

```go
eg := errors.NewErrorGroup("配置验证失败")
eg.Add(errors.New("缺少 API 密钥"))
eg.Add(errors.New("数据库 URL 无效"))
```

## 5. 错误检查与处理

### 5.1. `errors.Is(err, target error) bool` (标准库)
此模块的错误类型 (`fundamental`, `wrapper`, `withCode`, `ErrorGroup`) 与标准库的 `errors.Is` 函数兼容。
- 对于 `fundamental` 错误，`Is` 检查目标是否也是具有相同消息的 `fundamental` 错误。
- 对于 `wrapper` 错误，`Is` 将解包到 `cause`。
- 对于 `withCode` 错误，如果 `target` 是 `Coder` 或具有返回匹配代码的 `Coder()` 方法，`Is` 将首先尝试匹配 `Coder`。否则，它将解包到其 `cause`。
- 对于 `ErrorGroup`，`errors.Is` 将检查是否有任何包含的错误 (或其原因) 与 `target` 匹配。
- `Coder` 本身 (实现了 `error`) 可以用作 `target`。

```go
// 假设 ErrUserNotFound 是一个 Coder
if errors.Is(err, ErrUserNotFound) {
    //专门处理用户未找到的情况
}
```

### 5.2. `errors.As(err error, target interface{}) bool` (标准库)
与 `errors.As` 兼容。
- 对于 `fundamental` 错误，可以提取 `**fundamental`。
- 对于 `wrapper` 错误，可以提取 `**wrapper`，然后解包到 `cause`。
- 对于 `withCode` 错误，可以提取 `**withCode` 或 `*Coder` (以获取关联的 `Coder`)，然后解包到 `cause`。
- 对于 `ErrorGroup`，`errors.As` 将尝试查找可以分配给 `target` 的包含错误。

```go
var coder errors.Coder
if errors.As(err, &coder) {
    fmt.Printf("错误代码: %d, HTTP 状态: %d\n", coder.Code(), coder.HTTPStatus())
}
```

### 5.3. `Cause(err error) error`
递归地使用其 `Unwrap()` 方法 (或者如果中间错误类型上存在自定义的 `Cause()` 方法) 解包 `err`，直到找到原始的、根本的错误。如果错误没有原因，则返回错误本身。

```go
rootCause := errors.Cause(complexWrappedError)
```

### 5.4. `GetCoder(err error) Coder`
递归地解包 `err` 并返回错误链中找到的第一个非 `nil` 的 `Coder`。如果未找到 `Coder`，则返回 `nil`。

```go
coder := errors.GetCoder(err)
if coder != nil {
    // 使用 coder
}
```

### 5.5. `IsCode(err error, c Coder) bool`
定义于 `pkg/errors/api.go`。检查错误 `err` (或其链中的任何错误) 是否拥有一个 `Coder`，该 `Coder` 的代码与提供的 `Coder c` 的代码匹配。这对于检查由 `Coder` 代码表示的特定错误 *类别* 非常有用，即使 `Coder` 实例不同。

```go
if errors.IsCode(err, ErrPaymentFailed) { // 假设 ErrPaymentFailed 是一个 Coder
    // 处理任何支付失败的逻辑，无论具体消息如何
}
```

## 6. 格式化输出

此模块中的所有错误类型 (`fundamental`, `wrapper`, `withCode`, `ErrorGroup`) 都实现了 `fmt.Formatter` 接口。

- **`%s`, `%v`**:
    - 对于 `fundamental`: 打印错误消息。
    - 对于 `wrapper`: 打印组合消息 (`wrapper.msg: cause.Error()`)。
    - 对于 `withCode`: 打印包含 `Coder` 字符串的消息 (例如 `coder.String(): cause.Error()`)。
    - 对于 `ErrorGroup`: 打印组及其所有包含错误的组合消息。(行为已通过先前的 ErrorGroup.Error() 方法更新进行了优化)
- **`%+v`**:
    - 对于 `fundamental`: 打印错误消息及其堆栈跟踪。
    - 对于 `wrapper`: 打印组合消息，然后是在 *此包装点* 捕获的堆栈跟踪。
    - 对于 `withCode`: 打印消息 (包括 `Coder` 信息)，然后是在 *附加此 `Coder` 的点* 捕获的堆栈跟踪。
    - 对于 `ErrorGroup`: 打印组的主消息 (如果有)，然后是每个所含错误的详细、多行表示。如果包含的错误也支持 `%+v` (如此包中的错误)，则会打印其完整详细信息，包括其自己的堆栈跟踪。如果组为空且没有消息，则打印 "空错误组"。如果组为空但有消息，则打印组消息后跟一个换行符。

```go
// 打印带有完整堆栈跟踪的错误:
fmt.Printf("%+v\n", err)
```

## 7. 预定义 Coders

`pkg/errors/coder.go` 文件为方便起见定义了一组常见的 `Coder` 实例。示例包括：
- `unknownCoder`: 通用内部服务器错误。
- `ErrInternalServer`: HTTP 500。
- `ErrNotFound`: HTTP 404。
- `ErrBadRequest`: HTTP 400。
- `ErrUnauthorized`: HTTP 401。
- `ErrForbidden`: HTTP 403。
- `ErrValidation`: HTTP 400 (通常用于请求验证问题)。
- `ErrTimeout`: HTTP 504。
- `ErrTooManyRequests`: HTTP 429.
- `ErrOperationFailed`: 通用操作失败。
- `config` 包错误的特定编码器 (例如 `ErrConfigFileRead`)。
- `log` 包错误的特定编码器 (例如 `ErrLogOptionInvalid`)。

有关完整列表及其定义，请参阅 `pkg/errors/coder.go`。

## 8. 最佳实践和示例

- **对初始错误使用 `New` 或 `Errorf`**: 当错误源于您的代码时。
- **使用 `Wrap` 或 `Wrapf` 添加上下文**: 当来自较低层的错误在传播之前需要更多高级别信息时。
- **使用 `WithCode` 或 `NewWithCode`/`ErrorfWithCode` 对错误进行分类**: 关联预定义或自定义的 `Coder` 以启用程序化处理或特定日志记录。
- **使用 `errors.Is` 和 `errors.As` 检查错误**: 优先使用这些而不是类型断言，以获得灵活性和与包装错误的兼容性。使用 `IsCode` 检查特定的 `Coder` 代码。
- **使用 `Cause` 查找根本错误**: 如果错误被深度包装，这对于记录原始故障很有用。
- **使用 `GetCoder` 提取错误代码**: 当您需要根据 `Coder` 信息采取行动时。
- **使用 `%+v` 格式化详细的调试日志**: 这将提供消息和完整的堆栈跟踪。
- **定义特定于应用程序的 `Coders`**: 在您的应用程序或特定领域包中创建自己的 `Coder` 实例，以便更好地对错误进行分类。
- **对多个错误使用 `ErrorGroup`**: 当一个操作可能同时以多种方式失败时 (例如，验证请求的多个字段)。

有关可运行的示例，请参阅 `pkg/errors/example_test.go`。

## 9. Go 版本特定说明

- **错误包装**: `wrapper`、`withCode` 上的 `Unwrap()` 方法以及 `ErrorGroup` 上的 `Unwrap() []error` 方法确保了与 Go 1.13+ 中的 `errors.Is` 和 `errors.As` (用于单个错误解包) 以及 Go 1.20+ (用于 `ErrorGroup` 的多错误解包) 的兼容性。
- **堆栈跟踪**: 堆栈跟踪收集使用 `runtime` 包。

本规范旨在为使用 `pkg/errors` 模块提供全面的指南。有关具体的实现细节，请参阅源代码。
