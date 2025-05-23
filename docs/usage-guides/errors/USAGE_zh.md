# lmcc-go-sdk Errors 模块使用指南

**作者:** Martin, AI Assistant

[English Version (英文版)](USAGE.md) | [详细规范文档 (Detailed Specification)](SPECIFICATION.md)

## 概述

`pkg/errors` 模块为 Go 应用程序提供了一种强大且结构化的错误处理方法。它通过以下方式扩展了标准库：

- **自动堆栈跟踪**，便于调试
- **错误码**，用于程序化错误处理
- 通过包装实现**丰富的错误上下文**
- **错误聚合**，用于收集多个故障
- **与标准库兼容** (`errors.Is`, `errors.As`)

## 快速入门

替换您现有的错误处理方式：

```go
// 之前：
import "errors"
err := errors.New("操作失败")

// 之后：
import "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
err := errors.New("操作失败") // 现在带有堆栈跟踪！
```

## 基本用法

### 1. 创建错误

```go
// 带堆栈跟踪的简单错误
err := errors.New("连接失败")

// 带堆栈跟踪的格式化错误
err := errors.Errorf("处理用户 %s 失败", userID)
```

### 2. 添加上下文

```go
// 包装错误以添加上下文
dbErr := db.Connect()
if dbErr != nil {
    return errors.Wrap(dbErr, "初始化数据库失败")
}

// 格式化包装
return errors.Wrapf(dbErr, "主机 %s 的数据库连接失败", host)
```

### 3. 使用错误码

```go
// 使用预定义的错误码
err := errors.NewWithCode(errors.ErrNotFound, "未找到用户")

// 创建自定义错误码
var ErrInvalidInput = errors.NewCoder(40001, 400, "无效输入", "")
err := errors.NewWithCode(ErrInvalidInput, "电子邮件格式无效")
```

### 4. 错误检查

```go
// 检查特定错误类型
if errors.Is(err, errors.ErrNotFound) {
    // 处理未找到的情况
}

// 提取错误码
if coder := errors.GetCoder(err); coder != nil {
    fmt.Printf("错误码: %d, HTTP 状态: %d", coder.Code(), coder.HTTPStatus())
}

// 按错误码检查
if errors.IsCode(err, errors.ErrValidation) {
    // 处理验证错误
}
```

## 高级用法

### 错误聚合

将多个错误收集到一个错误中：

```go
eg := errors.NewErrorGroup("验证失败")

// 添加单个错误
if username == "" {
    eg.Add(errors.New("用户名是必需的"))
}
if email == "" {
    eg.Add(errors.New("电子邮件是必需的"))
}

// 检查是否发生任何错误
if len(eg.Errors()) > 0 {
    return eg // 返回组合的错误消息
}
```

### 堆栈跟踪

获取详细的堆栈跟踪以进行调试：

```go
// 打印带堆栈跟踪的详细错误
fmt.Printf("%+v\n", err)

// 示例输出：
// failed to save user: database connection failed
// github.com/yourapp/pkg/user.Save
//     /path/to/your/user.go:42
// github.com/yourapp/cmd/api.handleCreateUser
//     /path/to/your/api.go:123
```

### 错误链导航

```go
// 获取根本原因
rootErr := errors.Cause(wrappedErr)

// 提取特定错误类型
var validationErr *ValidationError
if errors.As(err, &validationErr) {
    //专门处理验证错误
}
```

## 预定义错误码

模块包含常用的错误码供立即使用：

### 通用错误
- `ErrInternalServer` (100001) - HTTP 500
- `ErrNotFound` (100002) - HTTP 404  
- `ErrBadRequest` (100003) - HTTP 400
- `ErrUnauthorized` (100004) - HTTP 401
- `ErrForbidden` (100005) - HTTP 403
- `ErrValidation` (100006) - HTTP 400
- `ErrTimeout` (100007) - HTTP 504
- `ErrTooManyRequests` (100008) - HTTP 429
- `ErrOperationFailed` (100009) - HTTP 500

### Config 包错误 (200001-200006)
- `ErrConfigFileRead`, `ErrConfigSetup`, `ErrConfigEnvBind` 等。

### Log 包错误 (300001-300008) 
- `ErrLogInternal`, `ErrLogOptionInvalid`, `ErrLogReconfigure` 等。

## 最佳实践

### 1. 使用错误码进行分类

```go
// 推荐：使用特定的错误码
return errors.NewWithCode(errors.ErrValidation, "无效的电子邮件格式")

// 避免：没有上下文的通用错误
return errors.New("无效输入")
```

### 2. 包装时添加上下文

```go
// 推荐：添加有意义的上下文
return errors.Wrap(err, "创建用户帐户失败")

// 避免：冗余包装
return errors.Wrap(err, "发生错误")
```

### 3. 按类型检查错误，而不是字符串

```go
// 推荐：按错误码检查
if errors.IsCode(err, errors.ErrNotFound) {
    return http.StatusNotFound, "未找到资源"
}

// 避免：字符串匹配
if strings.Contains(err.Error(), "not found") {
    // 不可靠且易出错
}
```

### 4. 对多个故障使用 ErrorGroup

```go
// 推荐：收集验证错误
func ValidateUser(user *User) error {
    eg := errors.NewErrorGroup("用户验证失败")
    
    if user.Email == "" {
        eg.Add(errors.NewWithCode(errors.ErrValidation, "电子邮件是必需的"))
    }
    if user.Age < 0 {
        eg.Add(errors.NewWithCode(errors.ErrValidation, "年龄必须为正数"))
    }
    
    if len(eg.Errors()) > 0 {
        return eg
    }
    return nil
}
```

## 从标准库迁移

### 简单迁移

```go
// 之前
import "errors"
err := errors.New("某些操作失败")

// 之后  
import "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
err := errors.New("某些操作失败") // 现在带有堆栈跟踪！
```

### 增强迁移

```go
// 之前
import "fmt"
return fmt.Errorf("处理 %s 失败: %w", id, err)

// 之后 - 更结构化
return errors.Wrapf(err, "处理 %s 失败", id)

// 或使用错误码
return errors.WithCode(err, errors.ErrOperationFailed)
```

## 集成示例

### HTTP 处理程序错误处理

```go
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
    user, err := h.parseUser(r)
    if err != nil {
        h.handleError(w, errors.WithCode(err, errors.ErrBadRequest))
        return
    }
    
    if err := h.userService.Create(user); err != nil {
        h.handleError(w, err)
        return
    }
    
    w.WriteHeader(http.StatusCreated)
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
    // 提取错误码以获取 HTTP 状态
    if coder := errors.GetCoder(err); coder != nil {
        w.WriteHeader(coder.HTTPStatus())
        json.NewEncoder(w).Encode(map[string]interface{}{
            "error": coder.String(),
            "code": coder.Code(),
        })
        return
    }
    
    // 未知错误的回退处理
    w.WriteHeader(http.StatusInternalServerError)
    json.NewEncoder(w).Encode(map[string]string{
        "error": "Internal server error",
    })
}
```

### 日志集成

*请注意：下面示例中使用的 `log` 包（例如 `log.WithFields`、`log.Error`、`log.Debugf`、`log.GetLevel()`、`log.DebugLevel`）是一个占位符，代表您项目中选用的日志库（例如 Logrus、Zap 或 `lmcc-go-sdk/pkg/log` 模块）。您需要根据您的具体日志设置调整此示例。*

```go
func logError(err error) {
    if coder := errors.GetCoder(err); coder != nil {
        log.WithFields(log.Fields{
            "error_code": coder.Code(),
            "http_status": coder.HTTPStatus(),
            "reference": coder.Reference(),
        }).Error(err.Error())
    } else {
        log.Error(err.Error())
    }
    
    // 在调试模式下记录完整的堆栈跟踪
    if log.GetLevel() == log.DebugLevel {
        log.Debugf("Stack trace: %+v", err)
    }
}
```

## 自定义错误码

为领域特定的错误定义您自己的错误码：

```go
// 定义自定义错误码
var (
    ErrUserEmailExists = errors.NewCoder(50001, 409, "用户电子邮件已存在", "")
    ErrInvalidPassword = errors.NewCoder(50002, 400, "密码不符合要求", "")
    ErrAccountLocked   = errors.NewCoder(50003, 423, "帐户已锁定", "")
)

// 在您的应用程序中使用它们
func (s *UserService) CreateUser(email, password string) error {
    if s.userExists(email) {
        return errors.NewWithCode(ErrUserEmailExists, fmt.Sprintf("电子邮件 %s 的用户已存在", email))
    }
    
    if !s.isValidPassword(password) {
        return errors.NewWithCode(ErrInvalidPassword, "密码长度至少为8个字符")
    }
    
    // ... 其余创建逻辑
    return nil
}
```

## 故障排除

### 堆栈跟踪未显示？

确保您使用的是 `%+v` 格式说明符：

```go
// 这仅显示错误消息
fmt.Printf("%v\n", err)

// 这显示完整的堆栈跟踪
fmt.Printf("%+v\n", err)
```

### 错误码未被检测到？

确保您使用的是错误码检查函数：

```go
// 正确的错误码检查方式
if errors.IsCode(err, errors.ErrNotFound) {
    // 处理未找到的情况
}

// 这样将无法按预期工作
if err == errors.ErrNotFound {
    // 这比较的是不同的错误实例
}
```

### 性能考虑

堆栈跟踪收集的开销很小，但如果您处于高性能场景中：

- 使用错误码以避免字符串比较
- 对频繁创建的错误考虑使用错误池
- 使用 `ErrorGroup` 批量处理多个错误，而不是创建许多单个错误

有关详细的 API 文档和完整的规范，请参阅 [SPECIFICATION.md](SPECIFICATION.md)。