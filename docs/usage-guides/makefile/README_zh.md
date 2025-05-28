# Makefile 文档

本目录包含 `lmcc-go-sdk` Makefile 系统的全面文档。

## 📖 文档文件

- **[English Guide](makefile_usage_en.md)** - 完整的英文 Makefile 使用指南
- **[中文指南](makefile_usage_zh.md)** - 完整的中文 Makefile 使用指南

## 🚀 快速参考

### 核心开发命令
```bash
make help              # 显示所有可用命令
make all               # 格式化、检查、测试和整理（默认）
make format            # 格式化 Go 源代码
make lint              # 运行代码检查器
make test-unit         # 运行单元测试
make test-integration  # 运行集成测试
make cover             # 生成覆盖率报告
make tidy              # 整理 go 模块
make clean             # 删除构建产物
```

### 示例管理（5 个分类中的 19 个示例）
```bash
make examples-list                           # 列出所有示例
make examples-build                          # 构建所有示例
make examples-run EXAMPLE=basic-usage       # 运行特定示例
make examples-run-all                       # 运行所有示例
make examples-test                          # 测试所有示例（检查 + 构建）
make examples-debug EXAMPLE=basic-usage     # 使用 delve 调试
make examples-analyze                       # 分析示例结构
make examples-category CATEGORY=config-features  # 按分类运行示例
make examples-clean                         # 清理示例二进制文件
```

### 文档命令
```bash
make doc-view PKG=./pkg/log    # 在终端查看包文档
make doc-serve                 # 启动本地文档服务器
```

### 工具管理
```bash
make tools                     # 安装所有必需工具
make tools.version            # 显示工具版本
make tools.help               # 显示工具帮助
```

## 📊 示例概览

项目包含 **19 个示例**，分为 **5 个分类**：

| 分类 | 数量 | 描述 |
|------|------|------|
| `basic-usage` | 1 | 基础集成示例 |
| `config-features` | 5 | 配置模块演示 |
| `error-handling` | 5 | 错误处理模式 |
| `integration` | 3 | 完整集成场景 |
| `logging-features` | 5 | 日志模块功能 |

## 🎯 常用工作流

### 开发工作流
```bash
# 1. 格式化和检查代码
make format lint

# 2. 运行测试
make test-unit

# 3. 检查覆盖率
make cover

# 4. 测试示例
make examples-test

# 5. 提交前的最终检查
make all
```

### 示例探索
```bash
# 1. 查看可用示例
make examples-list

# 2. 尝试基础示例
make examples-run EXAMPLE=basic-usage

# 3. 探索某个分类
make examples-category CATEGORY=config-features

# 4. 必要时进行调试
make examples-debug EXAMPLE=config-features/01-simple-config
```

### 文档生成
```bash
# 1. 查看特定包的文档
make doc-view PKG=./pkg/config

# 2. 启动文档服务器进行浏览
make doc-serve
# 在浏览器中打开 http://localhost:6060
```

## 🔧 变量

### 包选择
- `PKG=./path/to/package` - 指定测试/覆盖率的目标包

### 示例管理
- `EXAMPLE=<名称>` - 指定要运行/调试的示例
- `CATEGORY=<名称>` - 指定要运行的示例分类

### 工具配置
- `GOLANGCI_LINT_STRATEGY=stable|latest|auto` - 检查器版本策略
- `V=1` - 启用详细输出

## 📁 目录结构

```
scripts/make-rules/
├── common.mk          # 通用变量和函数
├── golang.mk          # Go 特定的构建规则
└── tools.mk           # 工具安装和管理

scripts/
├── examples-build.sh     # 构建所有示例
├── examples-run-all.sh   # 顺序运行所有示例
├── examples-analyze.sh   # 分析示例结构
├── examples-category.sh  # 按分类运行示例
└── format-func-coverage.sh  # 格式化覆盖率输出
```

## 🎪 功能特性

### ✅ 代码质量
- 自动代码格式化（`gofmt`、`goimports`）
- 全面的代码检查（`golangci-lint`）
- 竞态条件检测
- 代码覆盖率报告

### ✅ 测试
- 支持包选择的单元测试
- 集成测试
- 覆盖率报告（文本 + HTML）
- 函数级覆盖率详情

### ✅ 示例管理
- 自动发现示例
- 并行构建
- 基于分类的执行
- 交互式调试
- 全面分析

### ✅ 文档
- 终端文档查看
- 本地文档服务器
- 自动工具安装

### ✅ 开发工具
- 自动工具安装
- 版本管理策略
- 工具版本报告

## 📚 了解更多

详细的使用说明和示例，请参见：
- [English Documentation](makefile_usage_en.md)
- [中文文档](makefile_usage_zh.md)

## 🤝 贡献

添加新的 Makefile 目标时：
1. 使用 `## 注释` 添加适当的帮助文本
2. 遵循现有的命名约定
3. 更新文档
4. 全面测试

添加新示例时：
1. 在 `examples/` 下创建目录
2. 包含 `main.go` 文件
3. 添加分类相关的文档
4. 示例会自动被检测 