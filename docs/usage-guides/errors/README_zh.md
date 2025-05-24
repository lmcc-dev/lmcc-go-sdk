<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

# `pkg/errors` - Go 增强错误处理机制

`lmcc-go-sdk` 中的 `pkg/errors` 模块为 Go 应用程序提供了一种健壮且对开发者友好的错误处理方式。它通过引入自动堆栈跟踪、错误码、上下文包装等功能，扩展了标准库的错误处理能力，旨在简化调试并改进错误管理。

**[English Version (英文版说明)](README.md)**

## 概述

宏观了解 `pkg/errors` 提供的功能及其主要优势。
*   [模块概述](./zh/00_overview_zh.md)

## 快速入门

查看一个快速示例，了解如何使用 `pkg/errors` 及其带来的直接好处。
*   [快速入门指南](./zh/01_quick_start_zh.md)

## 目录

### 1. 核心概念
对基础功能的详细说明。
*   [创建错误](./zh/02_01_creating_errors_zh.md) - 学习如何使用 `New` 和 `Errorf` 创建基本错误。
*   [添加上下文](./zh/02_02_adding_context_zh.md) - 理解如何使用 `Wrap` 和 `Wrapf` 包装错误以提供更多上下文。
*   [使用错误码](./zh/02_03_using_error_codes_zh.md) - 探索如何使用 `Coder` 对象将错误与特定代码关联。
*   [错误检查](./zh/02_04_error_checking_zh.md) - 学习如何使用 `standardErrors.Is`、`GetCoder` 和 `IsCode` 来检查错误。

### 2. 高级特性
探索模块更高级的功能。
*   [错误聚合](./zh/03_01_error_aggregation_zh.md) - 将多个错误收集到单个 `ErrorGroup` 中。
*   [堆栈跟踪](./zh/03_02_stack_traces_zh.md) - 理解堆栈跟踪是如何被捕获和显示的。

### 3. 最佳实践
有效使用 `pkg/errors` 的建议。
*   [最佳实践简介](./zh/04_00_best_practices_intro_zh.md)
*   [定义特定于应用的错误码](./zh/04_01_define_application_specific_error_codes_zh.md)
*   [包装错误以添加上下文](./zh/04_02_wrap_errors_for_context_zh.md)
*   [在不同层级恰当地处理错误](./zh/04_03_handle_errors_appropriately_zh.md)
*   [使用 `IsCode` 检查特定的错误类别](./zh/04_04_use_iscode_for_checking_zh.md)
*   [有效地记录错误](./zh/04_05_log_errors_effectively_zh.md)

### 4. 迁移指南
关于从标准库 `errors` 包迁移的指导。
*   [迁移简介](./zh/05_00_migration_intro_zh.md)
*   [替换 `errors.New` 和 `fmt.Errorf`](./zh/05_01_replacing_new_and_errorf_zh.md)
*   [包装错误](./zh/05_02_wrapping_errors_zh.md)
*   [使用 `errors.Is` 和 `errors.As` 进行错误检查](./zh/05_03_error_checking_is_as_zh.md)

### 5. 集成
`pkg/errors` 如何与应用程序其他部分集成的示例。
*   [与日志库集成](./zh/06_01_with_logging_libraries_zh.md)
*   [与 API/HTTP 处理程序集成](./zh/06_02_with_api_http_handlers_zh.md)

### 6. 自定义错误码
学习如何定义和管理您自己的错误码集合。
*   [定义和使用自定义错误码](./zh/07_custom_error_codes_zh.md)

### 7. 问题排查 (Troubleshooting)
使用 `pkg/errors` 时遇到的常见问题及其解决方法。
*   [问题排查指南](./zh/08_troubleshooting_zh.md)

### 8. 模块规范 (API 参考) (Module Specification (API Reference))
`pkg/errors` 模块公共 API 的详细规范。
*   [模块规范](./zh/09_module_specification_zh.md)

## 贡献 (Contributing)

该模块旨在提供一个全面且易于使用的错误处理机制。通过利用其功能，开发人员可以构建更具弹性和可维护性的 Go 应用程序。 