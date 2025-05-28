# CLI工具集成示例

[English Version](README.md)

此示例展示了使用lmcc-go-sdk包的综合命令行界面(CLI)工具实现。它演示了如何构建一个专业的CLI应用程序，包含子命令、参数解析、配置管理和结构化日志。

## 功能特性

- **子命令架构**: 基于接口设计的模块化命令系统
- **参数解析**: 内置参数解析，支持标志和选项
- **配置管理**: 基于YAML的配置，带有默认值
- **帮助系统**: 命令和使用信息的综合帮助
- **多种输出格式**: 支持表格、JSON和纯文本输出
- **结构化日志**: 集成日志记录，带上下文和级别
- **错误处理**: 优雅的错误处理和详细消息
- **用户管理**: 用户管理的完整CRUD操作
- **数据导入/导出**: 基于文件的数据交换功能

## CLI架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CLI路由器     │    │   命令接口      │    │  文件存储       │
│                 │    │                 │    │                 │
│ • 参数解析      │───▶│ • 子命令        │───▶│ • JSON文件      │
│ • 帮助系统      │    │ • 验证          │    │ • CRUD操作      │
│ • 错误处理      │    │ • 执行          │    │ • 持久化        │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  配置管理       │    │   日志记录      │    │  输出格式       │
│                 │    │                 │    │                 │
│ • YAML配置      │    │ • 结构化        │    │ • 表格视图      │
│ • 默认值        │    │ • 上下文        │    │ • JSON导出      │
│ • 环境变量      │    │ • 级别          │    │ • 文本输出      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 可用命令

### 核心命令
- **`create`** - 创建新用户，包含验证
- **`list`** - 列出所有用户，支持过滤选项
- **`get`** - 通过ID或用户名检索用户
- **`update`** - 更新用户信息
- **`delete`** - 删除用户，支持强制标志
- **`search`** - 跨字段通过关键词搜索用户

### 实用命令
- **`export`** - 将用户导出到JSON或CSV文件
- **`import`** - 从文件导入用户，支持合并选项
- **`help`** - 显示命令帮助信息
- **`version`** - 显示版本和构建信息

## 配置

CLI工具支持带有合理默认值的YAML配置：

```yaml
app:
  name: "user-cli"
  version: "v1.0.0"
  description: "用户管理CLI工具"

database:
  type: "file"
  path: "./users.json"

output:
  format: "table"  # table, json, plain
  quiet: false
  color: true

logging:
  level: "info"
  format: "text"  # text, json
  output_paths: ["stdout"]
```

## 使用示例

### 基本操作

```bash
# 显示帮助
./cli-tool help

# 显示特定命令帮助
./cli-tool help create

# 创建用户
./cli-tool create alice alice@example.com --name "Alice Smith" --status active

# 列出所有用户
./cli-tool list

# 带过滤器列出用户
./cli-tool list --status active --limit 10

# 获取特定用户
./cli-tool get alice

# 更新用户
./cli-tool update alice --email newemail@example.com --status inactive

# 删除用户
./cli-tool delete alice --force

# 搜索用户
./cli-tool search smith --field name
```

### 数据管理

```bash
# 导出用户到JSON
./cli-tool export backup.json --format json

# 导出用户到CSV
./cli-tool export users.csv --format csv

# 从文件导入用户
./cli-tool import backup.json

# 带合并导入（更新现有用户）
./cli-tool import backup.json --merge
```

### 输出格式

```bash
# 表格格式（默认）
./cli-tool list

# JSON格式
./cli-tool list --format json

# 静默模式（最小输出）
./cli-tool create bob bob@example.com --quiet
```

## 示例输出

### 演示运行（无参数）

当无参数运行时，工具演示其功能：

```
=== 运行CLI工具演示 ===

1. 显示帮助:
   命令: user-cli help
user-cli - 用户管理CLI工具
版本: v1.0.0

可用命令:
  create     创建新用户
  delete     删除用户
  export     将用户导出到文件
  get        通过ID或用户名获取用户
  help       显示帮助信息
  import     从文件导入用户
  list       列出所有用户
  search     通过关键词搜索用户
  update     更新用户信息
  version    显示版本信息

使用 'user-cli help <command>' 获取命令的更多信息。
   ✅ 成功

2. 创建用户alice:
   命令: user-cli create alice alice@example.com --name Alice Smith
┌────────────┬─────────────────────────────────────┐
│ 字段       │ 值                                  │
├────────────┼─────────────────────────────────────┤
│ ID         │ user_1748425264                     │
│ Username   │ alice                               │
│ Email      │ alice@example.com                   │
│ Name       │ Alice Smith                         │
│ Status     │ active                              │
│ Created    │ 2025-05-28 17:41:04                 │
└────────────┴─────────────────────────────────────┘
✅ 用户 'alice' 创建成功，ID: user_1748425264
   ✅ 成功

=== CLI工具演示完成 ===
```

### 表格输出格式

```
┌──────────────┬──────────────┬─────────────────────────┬──────────────┬─────────┬─────────────────────┐
│ ID           │ 用户名       │ 邮箱                    │ 姓名         │ 状态    │ 创建时间            │
├──────────────┼──────────────┼─────────────────────────┼──────────────┼─────────┼─────────────────────┤
│ user_001     │ alice        │ alice@example.com       │ Alice Smith  │ active  │ 2024-11-23 10:15   │
│ user_002     │ bob          │ bob@example.com         │ Bob Johnson  │ active  │ 2024-11-23 10:16   │
└──────────────┴──────────────┴─────────────────────────┴──────────────┴─────────┴─────────────────────┘
```

## 关键学习要点

### 1. CLI架构设计
- 可扩展性的基于接口的命令系统
- 集中式参数解析和路由
- 模块化命令实现

### 2. 配置管理
- 带结构标签的基于YAML配置
- 默认值处理
- 环境特定覆盖

### 3. 错误处理模式
- 用户友好消息的优雅错误处理
- 输入验证和使用信息
- 脚本集成的适当退出代码

### 4. 输出格式化
- 多种输出格式（表格、JSON、纯文本）
- 跨命令的一致格式化
- 颜色和静默模式支持

### 5. 用户体验
- 综合帮助系统
- 直观的命令结构
- 清晰的反馈和错误消息

## 实现亮点

### 命令接口

```go
type Command interface {
    Name() string
    Description() string
    Usage() string
    Execute(ctx context.Context, args []string) error
}
```

### 配置结构

```go
type CLIConfig struct {
    App struct {
        Name        string `yaml:"name" default:"user-cli"`
        Version     string `yaml:"version" default:"v1.0.0"`
        Description string `yaml:"description" default:"用户管理CLI工具"`
    } `yaml:"app"`
    // ... 其他配置部分
}
```

### 参数解析

CLI演示了强大的参数解析模式：

- 必需参数的位置参数
- 基于`--flag value`语法的标志选项
- 开关的布尔标志
- 输入验证和错误处理

## 生产环境考虑

此示例演示了适用于生产CLI工具的模式：

- **错误处理**: 全面的错误分类和用户友好消息
- **配置**: 带合理默认值的外部化配置
- **日志**: 带上下文和可配置级别的结构化日志
- **文档**: 内置帮助系统和使用信息
- **数据验证**: 带清晰错误消息的输入验证
- **退出代码**: shell脚本集成的适当退出代码

## 扩展点

此CLI工具可以通过以下方式扩展：

- **数据库集成**: 用真实数据库替换文件存储
- **身份验证**: 添加用户身份验证和授权
- **API集成**: 连接到REST API或微服务
- **高级解析**: 与Cobra或CLI等库集成
- **Shell补全**: 添加bash/zsh补全支持
- **交互模式**: 添加交互式提示和向导

## 测试

示例包含内置演示，用于测试：

- 命令解析和执行
- 帮助系统功能
- 配置加载
- 输出格式化
- 错误处理场景

运行演示：
```bash
go run main.go
```

进行交互式测试，使用特定命令：
```bash
go run main.go help
go run main.go create testuser test@example.com
```

此示例为使用Go和lmcc-go-sdk框架构建专业命令行工具提供了坚实的基础。