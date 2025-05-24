<!--
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This document was collaboratively developed by Martin and AI Assistant.
-->

### 2. Wrap Errors for Context

When an error is returned from a function call, wrap it with `errors.Wrap` or `errors.Wrapf` to add contextual information about what the calling function was trying to do. This helps in understanding the error's path and origin when debugging.

- **Be Specific**: The wrapping message should clearly state the operation that failed at that level.
- **Avoid Redundancy**: Don't wrap if the error already has sufficient context from lower levels or if no new, useful information can be added.
- **Preserve Original Error**: Wrapping ensures the original error (and its `Coder` or stack trace) is preserved and accessible via `errors.Cause`, `errors.GetCoder`, or `standardErrors.Is/As`.

```go
package main

import (
	"fmt"
	"os"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
)

// Config represents some application configuration.

type Config struct {
	Port     int
	Hostname string
	Debug    bool
}

// loadConfigFile simulates loading a configuration file.
// It might return an error if the file doesn't exist or is unreadable.
func loadConfigFile(filePath string) ([]byte, error) {
	fmt.Printf("Attempting to load config file: %s\n", filePath)
	// Simulate file not found error
	if filePath == "nonexistent.yaml" {
		// This error comes from the os package, it does not have a Coder yet.
		return nil, os.ErrNotExist // A standard library sentinel error
	}
	// Simulate permission denied
	if filePath == "restricted.yaml" {
		// Using pkg/errors to create an error with a Coder from the start.
		return nil, errors.NewWithCode(errors.ErrPermissionDenied, "permission denied reading restricted.yaml")
	}
	
	// Simulate successful load
	fmt.Printf("Successfully loaded file: %s\n", filePath)
	return []byte("port: 8080\nhostname: localhost\ndebug: true"), nil
}

// parseConfigData simulates parsing the raw config data.
func parseConfigData(data []byte, filePath string) (*Config, error) {
	fmt.Printf("Attempting to parse config data from: %s\n", filePath)
	// Simulate a parsing error (e.g., invalid YAML)
	if string(data) == "invalid yaml content" {
		// Create a new error specific to parsing, with a Validation Coder.
		return nil, errors.NewWithCode(errors.ErrValidation, "config data is not valid YAML")
	}

	// Simulate successful parsing
	cfg := &Config{Port: 8080, Hostname: "localhost", Debug: true} // Dummy parsed data
	fmt.Printf("Successfully parsed config from: %s\n", filePath)
	return cfg, nil
}

// setupApplication loads and parses configuration to set up an application.
// This function demonstrates wrapping errors from its helper functions.
func setupApplication(configPath string) (*Config, error) {
	// Step 1: Load the config file
	configData, err := loadConfigFile(configPath)
	if err != nil {
		// Wrap the error from loadConfigFile to add context about *this* function's operation.
		// We also attach a Coder (ErrNotFound) if the original error was os.ErrNotExist.
		if os.IsNotExist(err) { // Check for standard library's os.ErrNotExist
			return nil, errors.WrapWithCode(err, errors.ErrNotFound, fmt.Sprintf("failed to setup application: config file '%s' not found", configPath))
		} 
		// For other errors from loadConfigFile (like our ErrPermissionDenied), 
		// they might already have a Coder. Wrap just adds context.
		return nil, errors.Wrapf(err, "failed to setup application while loading config '%s'", configPath)
	}

	// Step 2: Parse the config data
	appConfig, err := parseConfigData(configData, configPath)
	if err != nil {
		// Wrap the error from parseConfigData.
		// The error from parseConfigData (ErrValidation) already has a Coder.
		return nil, errors.Wrapf(err, "failed to setup application: could not parse config data from '%s'", configPath)
	}

	fmt.Printf("Application setup successful with config from '%s'!\n", configPath)
	return appConfig, nil
}

func main() {
	 scenarios := []struct {
		name       string
		configPath string
		expectCoder errors.Coder // Expected Coder to check for, nil if no specific Coder expected at top level
	}{
		{"File Not Found", "nonexistent.yaml", errors.ErrNotFound},
		{"Permission Denied", "restricted.yaml", errors.ErrPermissionDenied},
		{"Invalid YAML Data", "valid_file_invalid_data.yaml", errors.ErrValidation}, // Assume this file exists but contains "invalid yaml content"
		{"Successful Setup", "production.yaml", nil},
	}

	// Mocking a file that exists but has bad content for the "Invalid YAML Data" scenario
	// For a real test, you might create temporary files.
	_ = os.WriteFile("valid_file_invalid_data.yaml", []byte("invalid yaml content"), 0644)
	_ = os.WriteFile("production.yaml", []byte("port: 80\nhostname: prod.example.com\ndebug: false"), 0644)
	defer os.Remove("valid_file_invalid_data.yaml")
	defer os.Remove("production.yaml")

	for _, s := range scenarios {
		fmt.Printf("\n--- Scenario: %s ---\n", s.name)
		config, err := setupApplication(s.configPath)
		if err != nil {
			fmt.Printf("Error during setup: %+v\n", err) // Print with stack trace

			// Check if the error chain contains the expected Coder
			if s.expectCoder != nil {
				if errors.IsCode(err, s.expectCoder) {
					fmt.Printf("Verified: Error contains expected Coder (Code: %d, Message: %s).\n", s.expectCoder.Code(), s.expectCoder.String())
				} else {
					actualCoder := errors.GetCoder(err)
					if actualCoder != nil {
						fmt.Printf("Verification FAILED: Expected Coder %d (%s), but got %d (%s).\n",
							s.expectCoder.Code(), s.expectCoder.String(), actualCoder.Code(), actualCoder.String())
					} else {
						fmt.Printf("Verification FAILED: Expected Coder %d (%s), but got no Coder.\n",
							s.expectCoder.Code(), s.expectCoder.String())
					}
				}
			}
		} else {
			fmt.Printf("Setup successful. Config: %+v\n", config)
		}
	}
}

/*
Example Output (Stack traces and exact paths/line numbers will vary):

--- Scenario: File Not Found ---
Attempting to load config file: nonexistent.yaml
Error during setup: failed to setup application: config file 'nonexistent.yaml' not found: file does not exist
main.setupApplication
	/path/to/your/file.go:64
main.main
	/path/to/your/file.go:107
...
Not found
Verified: Error contains expected Coder (Code: 100002, Message: Not found).

--- Scenario: Permission Denied ---
Attempting to load config file: restricted.yaml
Error during setup: failed to setup application while loading config 'restricted.yaml': permission denied reading restricted.yaml
main.loadConfigFile
	/path/to/your/file.go:28
main.setupApplication
	/path/to/your/file.go:57
main.main
	/path/to/your/file.go:107
...
Permission denied
Verified: Error contains expected Coder (Code: 100004, Message: Permission denied).

--- Scenario: Invalid YAML Data ---
Attempting to load config file: valid_file_invalid_data.yaml
Successfully loaded file: valid_file_invalid_data.yaml
Attempting to parse config data from: valid_file_invalid_data.yaml
Error during setup: failed to setup application: could not parse config data from 'valid_file_invalid_data.yaml': config data is not valid YAML
main.parseConfigData
	/path/to/your/file.go:42
main.setupApplication
	/path/to/your/file.go:75
main.main
	/path/to/your/file.go:107
...
Validation failed
Verified: Error contains expected Coder (Code: 100006, Message: Validation failed).

--- Scenario: Successful Setup ---
Attempting to load config file: production.yaml
Successfully loaded file: production.yaml
Attempting to parse config data from: production.yaml
Successfully parsed config from: production.yaml
Application setup successful with config from 'production.yaml'!
Setup successful. Config: &{Port:80 Hostname:prod.example.com Debug:false}
*/ 