<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 2. 包装错误以添加上下文 (Wrap Errors for Context)

当从函数调用返回错误时，用 `errors.Wrap` 或 `errors.Wrapf` 包装它，以添加关于调用函数试图做什么的上下文信息。这有助于在调试时理解错误的路径和来源。

(When an error is returned from a function call, wrap it with `errors.Wrap` or `errors.Wrapf` to add contextual information about what the calling function was trying to do. This helps in understanding the error\'s path and origin when debugging.)

- **具体说明 (Be Specific)**: 包装消息应清楚地说明在该级别失败的操作。
  (The wrapping message should clearly state the operation that failed at that level.)
- **避免冗余 (Avoid Redundancy)**: 如果错误已经从较低级别获得了足够的上下文，或者无法添加新的有用信息，则不要包装。
  (Don\'t wrap if the error already has sufficient context from lower levels or if no new, useful information can be added.)
- **保留原始错误 (Preserve Original Error)**: 包装可确保原始错误 (及其 `Coder` 或堆栈跟踪) 得以保留，并可通过 `errors.Cause`、`errors.GetCoder` 或 `standardErrors.Is/As` 进行访问。
  (Wrapping ensures the original error (and its `Coder` or stack trace) is preserved and accessible via `errors.Cause`, `errors.GetCoder`, or `standardErrors.Is/As`.)

```go
package main

import (
	"fmt"
	"os"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
)

// Config 表示某些应用程序配置。
// (Config represents some application configuration.)
type Config struct {
	Port     int
	Hostname string
	Debug    bool
}

// loadConfigFile 模拟加载配置文件。
// (loadConfigFile simulates loading a configuration file.)
// 如果文件不存在或不可读，它可能会返回错误。
// (It might return an error if the file doesn\'t exist or is unreadable.)
func loadConfigFile(filePath string) ([]byte, error) {
	fmt.Printf("尝试加载配置文件：%s\n", filePath)
	// (Attempting to load config file: %s\n)
	// 模拟文件未找到错误 (Simulate file not found error)
	if filePath == "nonexistent.yaml" {
		// 此错误来自 os 包，它尚无 Coder。
		// (This error comes from the os package, it does not have a Coder yet.)
		return nil, os.ErrNotExist // 标准库哨兵错误 (A standard library sentinel error)
	}
	// 模拟权限被拒绝 (Simulate permission denied)
	if filePath == "restricted.yaml" {
		// 从一开始就使用 pkg/errors 创建带有 Coder 的错误。
		// (Using pkg/errors to create an error with a Coder from the start.)
		return nil, errors.NewWithCode(errors.ErrPermissionDenied, "读取 restricted.yaml 时权限被拒绝 (permission denied reading restricted.yaml)")
	}
	
	// 模拟成功加载 (Simulate successful load)
	fmt.Printf("成功加载文件：%s\n", filePath)
	// (Successfully loaded file: %s\n)
	return []byte("port: 8080\nhostname: localhost\ndebug: true"), nil
}

// parseConfigData 模拟解析原始配置数据。
// (parseConfigData simulates parsing the raw config data.)
func parseConfigData(data []byte, filePath string) (*Config, error) {
	fmt.Printf("尝试从 %s 解析配置数据\n", filePath)
	// (Attempting to parse config data from: %s\n)
	// 模拟解析错误 (例如，无效的 YAML)
	// (Simulate a parsing error (e.g., invalid YAML))
	if string(data) == "invalid yaml content" {
		// 创建一个特定于解析的新错误，并带有 Validation Coder。
		// (Create a new error specific to parsing, with a Validation Coder.)
		return nil, errors.NewWithCode(errors.ErrValidation, "配置数据不是有效的 YAML (config data is not valid YAML)")
	}

	// 模拟成功解析 (Simulate successful parsing)
	cfg := &Config{Port: 8080, Hostname: "localhost", Debug: true} // 虚拟解析数据 (Dummy parsed data)
	fmt.Printf("成功从 %s 解析配置\n", filePath)
	// (Successfully parsed config from: %s\n)
	return cfg, nil
}

// setupApplication 加载并解析配置以设置应用程序。
// (setupApplication loads and parses configuration to set up an application.)
// 此函数演示了从其辅助函数包装错误。
// (This function demonstrates wrapping errors from its helper functions.)
func setupApplication(configPath string) (*Config, error) {
	// 步骤1：加载配置文件 (Step 1: Load the config file)
	configData, err := loadConfigFile(configPath)
	if err != nil {
		// 包装来自 loadConfigFile 的错误，以添加关于*此*函数操作的上下文。
		// (Wrap the error from loadConfigFile to add context about *this* function\'s operation.)
		// 如果原始错误是 os.ErrNotExist，我们还附加一个 Coder (ErrNotFound)。
		// (We also attach a Coder (ErrNotFound) if the original error was os.ErrNotExist.)
		if os.IsNotExist(err) { // 检查标准库的 os.ErrNotExist (Check for standard library\'s os.ErrNotExist)
			return nil, errors.WrapWithCode(err, errors.ErrNotFound, fmt.Sprintf("设置应用程序失败：未找到配置文件 '%s' (failed to setup application: config file '%s' not found)", configPath))
		} 
		// 对于来自 loadConfigFile 的其他错误 (例如我们的 ErrPermissionDenied)，
		// (For other errors from loadConfigFile (like our ErrPermissionDenied),)
		// 它们可能已经有一个 Coder。Wrap 仅添加上下文。
		// (they might already have a Coder. Wrap just adds context.)
		return nil, errors.Wrapf(err, "加载配置 '%s' 时设置应用程序失败 (failed to setup application while loading config '%s')", configPath)
	}

	// 步骤2：解析配置数据 (Step 2: Parse the config data)
	appConfig, err := parseConfigData(configData, configPath)
	if err != nil {
		// 包装来自 parseConfigData 的错误。
		// (Wrap the error from parseConfigData.)
		// 来自 parseConfigData 的错误 (ErrValidation) 已有一个 Coder。
		// (The error from parseConfigData (ErrValidation) already has a Coder.)
		return nil, errors.Wrapf(err, "设置应用程序失败：无法从 '%s' 解析配置数据 (failed to setup application: could not parse config data from '%s')", configPath)
	}

	fmt.Printf("使用来自 '%s' 的配置成功设置应用程序！\n", configPath)
	// (Application setup successful with config from '%s'!\n)
	return appConfig, nil
}

func main() {
	 scenarios := []struct {
		name       string
		configPath string
		expectCoder errors.Coder // 要检查的预期 Coder，如果顶层没有预期的特定 Coder，则为 nil
		                        // (Expected Coder to check for, nil if no specific Coder expected at top level)
	}{
		{"文件未找到 (File Not Found)", "nonexistent.yaml", errors.ErrNotFound},
		{"权限被拒绝 (Permission Denied)", "restricted.yaml", errors.ErrPermissionDenied},
		{"无效的 YAML 数据 (Invalid YAML Data)", "valid_file_invalid_data.yaml", errors.ErrValidation}, // 假设此文件存在但包含 "invalid yaml content"
		                                                                                             // (Assume this file exists but contains "invalid yaml content")
		{"成功设置 (Successful Setup)", "production.yaml", nil},
	}

	// 为"无效的 YAML 数据"场景模拟一个存在但内容错误的文件
	// (Mocking a file that exists but has bad content for the "Invalid YAML Data" scenario)
	// 对于真实测试，您可能会创建临时文件。
	// (For a real test, you might create temporary files.)
	_ = os.WriteFile("valid_file_invalid_data.yaml", []byte("invalid yaml content"), 0644)
	_ = os.WriteFile("production.yaml", []byte("port: 80\nhostname: prod.example.com\ndebug: false"), 0644)
	defer os.Remove("valid_file_invalid_data.yaml")
	defer os.Remove("production.yaml")

	for _, s := range scenarios {
		fmt.Printf("\n--- 场景：%s ---\n", s.name)
		// (--- Scenario: %s ---\n)
		config, err := setupApplication(s.configPath)
		if err != nil {
			fmt.Printf("设置期间出错：%+v\n", err) // 使用堆栈跟踪打印 (Print with stack trace)
			// (Error during setup: %+v\n)

			// 检查错误链是否包含预期的 Coder
			// (Check if the error chain contains the expected Coder)
			if s.expectCoder != nil {
				if errors.IsCode(err, s.expectCoder) {
					fmt.Printf("已验证：错误包含预期的 Coder (代码：%d，消息：%s)。\n", s.expectCoder.Code(), s.expectCoder.String())
					// (Verified: Error contains expected Coder (Code: %d, Message: %s).\n)
				} else {
					actualCoder := errors.GetCoder(err)
					if actualCoder != nil {
						fmt.Printf("验证失败：预期 Coder %d (%s)，但得到 %d (%s)。\n",
							s.expectCoder.Code(), s.expectCoder.String(), actualCoder.Code(), actualCoder.String())
						// (Verification FAILED: Expected Coder %d (%s), but got %d (%s).\n)
					} else {
						fmt.Printf("验证失败：预期 Coder %d (%s)，但没有得到 Coder。\n",
							s.expectCoder.Code(), s.expectCoder.String())
						// (Verification FAILED: Expected Coder %d (%s), but got no Coder.\n)
					}
				}
			}
		} else {
			fmt.Printf("设置成功。配置：%+v\n", config)
			// (Setup successful. Config: %+v\n)
		}
	}
}

/*
示例输出 (堆栈跟踪和确切的路径/行号将有所不同)：
(Example Output (Stack traces and exact paths/line numbers will vary)):

--- 场景：文件未找到 (File Not Found) ---
尝试加载配置文件：nonexistent.yaml
(Attempting to load config file: nonexistent.yaml)
设置期间出错：设置应用程序失败：未找到配置文件 'nonexistent.yaml'：file does not exist
(Error during setup: failed to setup application: config file 'nonexistent.yaml' not found: file does not exist)
main.setupApplication
	/path/to/your/file.go:64
main.main
	/path/to/your/file.go:107
...
Not found
已验证：错误包含预期的 Coder (代码：100002，消息：Not found)。
(Verified: Error contains expected Coder (Code: 100002, Message: Not found).)

--- 场景：权限被拒绝 (Permission Denied) ---
尝试加载配置文件：restricted.yaml
(Attempting to load config file: restricted.yaml)
设置期间出错：加载配置 'restricted.yaml' 时设置应用程序失败：读取 restricted.yaml 时权限被拒绝 (permission denied reading restricted.yaml)
(Error during setup: failed to setup application while loading config 'restricted.yaml': permission denied reading restricted.yaml)
main.loadConfigFile
	/path/to/your/file.go:28
main.setupApplication
	/path/to/your/file.go:57
main.main
	/path/to/your/file.go:107
...
Permission denied
已验证：错误包含预期的 Coder (代码：100004，消息：Permission denied)。
(Verified: Error contains expected Coder (Code: 100004, Message: Permission denied).)

--- 场景：无效的 YAML 数据 (Invalid YAML Data) ---
尝试加载配置文件：valid_file_invalid_data.yaml
(Attempting to load config file: valid_file_invalid_data.yaml)
成功加载文件：valid_file_invalid_data.yaml
(Successfully loaded file: valid_file_invalid_data.yaml)
尝试从 valid_file_invalid_data.yaml 解析配置数据
(Attempting to parse config data from: valid_file_invalid_data.yaml)
设置期间出错：设置应用程序失败：无法从 'valid_file_invalid_data.yaml' 解析配置数据：配置数据不是有效的 YAML (config data is not valid YAML)
(Error during setup: failed to setup application: could not parse config data from 'valid_file_invalid_data.yaml': config data is not valid YAML)
main.parseConfigData
	/path/to/your/file.go:42
main.setupApplication
	/path/to/your/file.go:75
main.main
	/path/to/your/file.go:107
...
Validation failed
已验证：错误包含预期的 Coder (代码：100006，消息：Validation failed)。
(Verified: Error contains expected Coder (Code: 100006, Message: Validation failed).)

--- 场景：成功设置 (Successful Setup) ---
尝试加载配置文件：production.yaml
(Attempting to load config file: production.yaml)
成功加载文件：production.yaml
(Successfully loaded file: production.yaml)
尝试从 production.yaml 解析配置数据
(Attempting to parse config data from: production.yaml)
成功从 production.yaml 解析配置
(Successfully parsed config from: production.yaml)
使用来自 'production.yaml' 的配置成功设置应用程序！
(Application setup successful with config from 'production.yaml'!)
设置成功。配置：&{Port:80 Hostname:prod.example.com Debug:false}
(Setup successful. Config: &{Port:80 Hostname:prod.example.com Debug:false})
*/
``` 