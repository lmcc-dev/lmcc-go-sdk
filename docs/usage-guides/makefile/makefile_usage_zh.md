# Makefile 使用指南

本指南解释了如何使用 `lmcc-go-sdk` 项目中包含的 `Makefile` 来管理常见的开发任务。

## 1. 功能介绍

`Makefile` 提供了一种标准化的方式来构建、测试、执行代码检查、格式化和清理项目。它借鉴了 `marmotedu` 的实践，利用 `scripts/make-rules/` 目录下的 include 文件来实现模块化。

核心目标：
- 简化常见的开发工作流程。
- 确保构建和测试的一致性。
- 自动化代码质量检查和格式化。
- 管理开发工具依赖。

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

-   **`make tools`**
    -   安装 `scripts/make-rules/tools.mk` 中列出的所有必需的开发工具（当前是 `golangci-lint` 和 `goimports`）。这对于初次设置开发环境很有用。

## 3. 自定义 (可选)

-   **添加工具**: 编辑 `scripts/make-rules/tools.mk` 文件，将新工具添加到 `TOOLS_REQUIRED` 变量中，并提供相应的 `install.<toolname>` 规则。
-   **详细输出**: 运行任何命令时加上 `V=1` 可以获取更详细的输出（例如 `make test V=1`）。 

### `make doc-view PKG=./path/to/package`

在终端中显示指定包的 Go 文档。`PKG` 变量是**必需的**。

*   **`PKG`**: 指定要显示文档的 Go 包的路径（例如，`./pkg/log`, `./pkg/config`）。

示例: `make doc-view PKG=./pkg/config`

### `make doc-serve`

启动一个本地 `godoc` HTTP 服务器（通常在 `http://localhost:6060`），用于在浏览器中浏览 Go 工作区中所有包的 HTML 文档，包括您当前的项目。如果 `godoc` 工具尚未安装，它会自动安装。

在终端中按 `Ctrl+C` 来停止服务。 