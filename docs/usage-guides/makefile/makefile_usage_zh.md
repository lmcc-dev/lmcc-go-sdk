# Makefile 使用指南

本指南解释了如何使用 `lmcc-go-sdk` 项目中包含的 `Makefile` 来管理常见的开发任务。

## 1. 功能介绍

`Makefile` 提供了一种标准化的方式来构建、测试、执行代码检查、格式化和清理项目。它借鉴了 `marmotedu` 的实践，利用 `scripts/make-rules/` 目录下的 include 文件来实现模块化。

核心目标：
- 简化常见的开发工作流程。
- 确保构建和测试的一致性。
- 自动化代码质量检查和格式化。
- 管理开发工具依赖。
- 提供集中化的示例管理。

## 2. 常用命令

您可以在项目的根目录运行这些命令。

-   **`make help`**
    -   显示帮助信息，列出所有可用的目标及其描述。这是查看所有可用命令的最佳方式。

-   **`make all`** (默认目标)
    -   运行一系列常用任务：`format`, `lint`, `test-unit`, 和 `tidy`。这对于在提交代码前进行快速检查很有用。(注意: 默认运行 `test-unit`，而非所有测试)。

-   **`make format`**
    -   使用标准工具 (`gofmt`, `goimports`) 格式化 Go 源代码。
    -   自动检查 `goimports` 是否已安装，如果缺少则安装（通过 `make tools.verify.goimports`）。

-   **`make lint`**
    -   运行代码 linter 以检查风格问题和潜在错误。
    -   当前使用 `go vet` 和 `golangci-lint`。
    -   自动检查 `golangci-lint` 是否已安装，如果缺少则安装（通过 `make tools.verify.golangci-lint`）。
    -   *注意:* 您可能需要配置 `.golangci.yaml` 文件以适应项目特定的规则。

-   **`make test-unit [PKG=...] [RUN=...]`**
    -   运行单元测试。
    -   包含 `-race` 标志以检测竞态条件。
    -   默认情况下，运行所有包中的单元测试（不包括 `examples/`, `vendor/` 目录以及包含 `test/integration` 的路径）。
    -   **可选 `PKG`**: 指定要进行单元测试的包。使用相对路径（例如 `make test-unit PKG=./pkg/log`）。如果指定了 `PKG`，则仅运行该包内的测试。
    -   **可选 `RUN`**: 根据匹配测试函数名的正则表达式过滤要运行的测试（例如 `make test-unit RUN=TestMyFunction`，`make test-unit PKG=./pkg/log RUN=^TestLog`）。

-   **`make test-integration [RUN=...]`**
    -   运行集成测试。这些测试通常位于 `test/integration/` 目录下。
    -   包含 `-race` 标志。
    -   此目标通常不使用 `PKG` 参数，因为它旨在运行所有集成测试。
    -   **可选 `RUN`**: 根据正则表达式过滤要运行的集成测试。

-   **`make cover [PKG=...]`**
    -   运行单元测试（类似于 `test-unit`）并生成代码覆盖率报告。
    -   将文本格式的 profile 保存到 `_output/coverage/coverage.out`。
    -   将 HTML 格式的报告保存到 `_output/coverage/coverage.html`，可以在浏览器中打开以进行详细分析。
    -   **可选 `PKG`**: 指定要进行覆盖率检查的包。如果指定了 `PKG`，则仅为这些包生成覆盖率。否则，它将覆盖所有进行单元测试的包。
    -   *注意:* 此目标专注于单元测试的覆盖率。

-   **`make tidy`**
    -   运行 `go mod tidy` 以确保 `go.mod` 和 `go.sum` 文件与源代码依赖项一致。

-   **`make clean`**
    -   移除生成的构建产物，包括 `_output` 目录（包含构建工件、覆盖率报告）和下载的工具目录 (`_output/tools`)。

-   **`make tools [GOLANGCI_LINT_STRATEGY=...]`**
    -   安装 `scripts/make-rules/tools.mk` 中列出的所有必需的开发工具（当前是 `golangci-lint`、`goimports` 和 `godoc`）。这对于初次设置开发环境很有用。
    -   **golangci-lint 版本策略**：
        - `stable`（默认）：使用经过测试的稳定版本 v1.64.8，确保可重现构建
        - `latest`：始终使用最新可用版本，获得前沿特性
        - `auto`：根据 Go 版本自动选择（Go 1.24+ 使用最新版，较旧版本使用稳定版）
    -   示例：
        - `make tools`（使用稳定策略）
        - `make tools GOLANGCI_LINT_STRATEGY=latest`
        - `make tools GOLANGCI_LINT_STRATEGY=auto`
        - `make tools GOLANGCI_LINT_VERSION=v1.65.0`（覆盖特定版本）

-   **`make tools.version`**
    -   显示所有已安装工具的版本和当前策略设置。

-   **`make tools.help`**
    -   显示工具管理策略的详细帮助和使用示例。

### `make doc-view PKG=./path/to/package`

在终端中显示指定包的 Go 文档。`PKG` 变量是**必需的**。

*   **`PKG`**: 指定要显示文档的 Go 包的路径（例如，`./pkg/log`, `./pkg/config`）。

示例: `make doc-view PKG=./pkg/config`

### `make doc-serve`

启动一个本地 `godoc` HTTP 服务器（通常在 `http://localhost:6060`），用于在浏览器中浏览 Go 工作区中所有包的 HTML 文档，包括您当前的项目。如果 `godoc` 工具尚未安装，它会自动安装。

在终端中按 `Ctrl+C` 来停止服务。

## 3. 示例管理

本项目包含一个全面的示例管理系统，允许您构建、运行、测试和调试 5 个分类中的 19 个示例。

### 可用分类

- **basic-usage**: 基础集成示例（1 个示例）
- **config-features**: 配置模块演示（5 个示例）
- **error-handling**: 错误处理模式（5 个示例）
- **integration**: 完整集成场景（3 个示例）
- **logging-features**: 日志模块功能（5 个示例）

### 示例命令

-   **`make examples-list`**
    -   列出所有可用示例及其编号和分类。
    -   显示示例总数。

-   **`make examples-build`**
    -   将所有示例构建到 `_output/examples/` 目录中。
    -   每个示例都成为一个独立的可执行文件。
    -   支持并行构建以获得更好的性能。

-   **`make examples-clean`**
    -   从 `_output/examples/` 中删除所有构建的示例二进制文件。
    -   用于清理构建产物。

-   **`make examples-run EXAMPLE=<名称>`**
    -   按名称运行特定示例。
    -   **必需的 `EXAMPLE`**: 指定要运行的示例（例如 `basic-usage`, `config-features/01-simple-config`）。
    -   运行前验证示例是否存在。
    -   示例：
        - `make examples-run EXAMPLE=basic-usage`
        - `make examples-run EXAMPLE=config-features/01-simple-config`
        - `make examples-run EXAMPLE=error-handling/02-error-wrapping`

-   **`make examples-run-all`**
    -   顺序运行所有示例并显示进度跟踪。
    -   显示每个示例的执行状态。
    -   用于全面测试所有示例。

-   **`make examples-test`**
    -   对所有示例执行全面测试：
        - 步骤 1：使用 `golangci-lint` 检查所有示例代码
        - 步骤 2：构建所有示例以验证编译
    -   确保示例的代码质量和可构建性。

-   **`make examples-debug EXAMPLE=<名称>`**
    -   使用 `delve` 为特定示例启动交互式调试会话。
    -   **必需的 `EXAMPLE`**: 指定要调试的示例。
    -   如果不存在会自动安装 `dlv`。
    -   启动时显示有用的调试器命令。
    -   示例：
        - `make examples-debug EXAMPLE=basic-usage`
        - `make examples-debug EXAMPLE=config-features/02-hot-reload`

-   **`make examples-analyze`**
    -   提供示例结构的全面分析：
        - 目录树可视化
        - 按分类统计
        - 代码指标（总行数、Go 文件数）
    -   用于理解项目结构和范围。

-   **`make examples-category CATEGORY=<名称>`**
    -   运行特定分类中的所有示例。
    -   **必需的 `CATEGORY`**: 指定要运行的分类。
    -   可用分类：`basic-usage`, `config-features`, `error-handling`, `integration`, `logging-features`
    -   示例：
        - `make examples-category CATEGORY=config-features`
        - `make examples-category CATEGORY=error-handling`

-   **`make examples`**
    -   运行 `examples-test` 的默认示例目标。
    -   相当于对所有示例运行 lint 和构建验证。

### 示例使用技巧

1. **从列表开始**: 使用 `make examples-list` 查看所有可用示例。
2. **运行单个示例**: 使用 `make examples-run EXAMPLE=<名称>` 测试特定功能。
3. **分类探索**: 使用 `make examples-category CATEGORY=<名称>` 探索相关示例。
4. **开发工作流**: 使用 `make examples-test` 确保所有示例在提交前正常工作。
5. **调试**: 当需要逐步调试示例代码时，使用 `make examples-debug EXAMPLE=<名称>`。

## 4. 自定义 (可选)

-   **添加工具**: 编辑 `scripts/make-rules/tools.mk` 文件，将新工具添加到 `TOOLS_REQUIRED` 变量中，并提供相应的 `install.<toolname>` 规则。
-   **详细输出**: 运行任何命令时加上 `V=1` 可以获取更详细的输出（例如 `make test V=1`）。
-   **添加示例**: 在 `examples/` 目录中创建包含 `main.go` 文件的新示例，它们将自动被示例管理系统检测到。 