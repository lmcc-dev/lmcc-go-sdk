# CLI Tool Integration Example

[中文版本](README_zh.md)

This example demonstrates a comprehensive command-line interface (CLI) tool implementation using the lmcc-go-sdk package. It showcases how to build a professional CLI application with subcommands, argument parsing, configuration management, and structured logging.

## Features

- **Subcommand Architecture**: Modular command system with interface-based design
- **Argument Parsing**: Built-in argument parsing with flags and options
- **Configuration Management**: YAML-based configuration with defaults
- **Help System**: Comprehensive help for commands and usage information
- **Multiple Output Formats**: Support for table, JSON, and plain text output
- **Structured Logging**: Integrated logging with context and levels
- **Error Handling**: Graceful error handling with detailed messages
- **User Management**: Complete CRUD operations for user management
- **Data Export/Import**: File-based data exchange capabilities

## CLI Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CLI Router    │    │   Command       │    │  File Storage   │
│                 │    │   Interface     │    │                 │
│ • Arg Parsing   │───▶│ • Subcommands   │───▶│ • JSON Files    │
│ • Help System   │    │ • Validation    │    │ • CRUD Ops      │
│ • Error Handler │    │ • Execution     │    │ • Persistence   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Configuration  │    │   Logging       │    │  Output Format  │
│                 │    │                 │    │                 │
│ • YAML Config   │    │ • Structured    │    │ • Table View    │
│ • Defaults      │    │ • Context       │    │ • JSON Export   │
│ • Environment   │    │ • Levels        │    │ • Text Output   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Available Commands

### Core Commands
- **`create`** - Create a new user with validation
- **`list`** - List all users with filtering options
- **`get`** - Retrieve user by ID or username
- **`update`** - Update user information
- **`delete`** - Delete user with optional force flag
- **`search`** - Search users by keyword across fields

### Utility Commands
- **`export`** - Export users to JSON or CSV files
- **`import`** - Import users from files with merge options
- **`help`** - Show help information for commands
- **`version`** - Display version and build information

## Configuration

The CLI tool supports YAML configuration with sensible defaults:

```yaml
app:
  name: "user-cli"
  version: "v1.0.0"
  description: "User management CLI tool"

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

## Usage Examples

### Basic Operations

```bash
# Show help
./cli-tool help

# Show help for specific command
./cli-tool help create

# Create a user
./cli-tool create alice alice@example.com --name "Alice Smith" --status active

# List all users
./cli-tool list

# List users with filters
./cli-tool list --status active --limit 10

# Get specific user
./cli-tool get alice

# Update user
./cli-tool update alice --email newemail@example.com --status inactive

# Delete user
./cli-tool delete alice --force

# Search users
./cli-tool search smith --field name
```

### Data Management

```bash
# Export users to JSON
./cli-tool export backup.json --format json

# Export users to CSV
./cli-tool export users.csv --format csv

# Import users from file
./cli-tool import backup.json

# Import with merge (update existing)
./cli-tool import backup.json --merge
```

### Output Formats

```bash
# Table format (default)
./cli-tool list

# JSON format
./cli-tool list --format json

# Quiet mode (minimal output)
./cli-tool create bob bob@example.com --quiet
```

## Sample Output

### Demo Run (No Arguments)

When run without arguments, the tool demonstrates its capabilities:

```
=== Running CLI Tool Demonstration ===

1. Show help:
   Command: user-cli help
user-cli - User management CLI tool
Version: v1.0.0

Available commands:
  create     Create a new user
  delete     Delete user
  export     Export users to file
  get        Get user by ID or username
  help       Show help information
  import     Import users from file
  list       List all users
  search     Search users by keyword
  update     Update user information
  version    Show version information

Use 'user-cli help <command>' for more information about a command.
   ✅ Success

2. Create user alice:
   Command: user-cli create alice alice@example.com --name Alice Smith
┌────────────┬─────────────────────────────────────┐
│ Field      │ Value                               │
├────────────┼─────────────────────────────────────┤
│ ID         │ user_1748425264                     │
│ Username   │ alice                               │
│ Email      │ alice@example.com                   │
│ Name       │ Alice Smith                         │
│ Status     │ active                              │
│ Created    │ 2025-05-28 17:41:04                 │
└────────────┴─────────────────────────────────────┘
✅ User 'alice' created successfully with ID: user_1748425264
   ✅ Success

=== CLI Tool Demonstration Completed ===
```

### Table Output Format

```
┌──────────────┬──────────────┬─────────────────────────┬──────────────┬─────────┬─────────────────────┐
│ ID           │ Username     │ Email                   │ Name         │ Status  │ Created             │
├──────────────┼──────────────┼─────────────────────────┼──────────────┼─────────┼─────────────────────┤
│ user_001     │ alice        │ alice@example.com       │ Alice Smith  │ active  │ 2024-11-23 10:15   │
│ user_002     │ bob          │ bob@example.com         │ Bob Johnson  │ active  │ 2024-11-23 10:16   │
└──────────────┴──────────────┴─────────────────────────┴──────────────┴─────────┴─────────────────────┘
```

## Key Learning Points

### 1. CLI Architecture Design
- Interface-based command system for extensibility
- Centralized argument parsing and routing
- Modular command implementations

### 2. Configuration Management
- YAML-based configuration with struct tags
- Default value handling
- Environment-specific overrides

### 3. Error Handling Patterns
- Graceful error handling with user-friendly messages
- Input validation and usage information
- Proper exit codes for script integration

### 4. Output Formatting
- Multiple output formats (table, JSON, plain)
- Consistent formatting across commands
- Color and quiet mode support

### 5. User Experience
- Comprehensive help system
- Intuitive command structure
- Clear feedback and error messages

## Implementation Highlights

### Command Interface

```go
type Command interface {
    Name() string
    Description() string
    Usage() string
    Execute(ctx context.Context, args []string) error
}
```

### Configuration Structure

```go
type CLIConfig struct {
    App struct {
        Name        string `yaml:"name" default:"user-cli"`
        Version     string `yaml:"version" default:"v1.0.0"`
        Description string `yaml:"description" default:"User management CLI tool"`
    } `yaml:"app"`
    // ... additional configuration sections
}
```

### Argument Parsing

The CLI demonstrates robust argument parsing patterns:

- Positional arguments for required parameters
- Flag-based options with `--flag value` syntax
- Boolean flags for switches
- Input validation and error handling

## Production Considerations

This example demonstrates patterns suitable for production CLI tools:

- **Error Handling**: Comprehensive error classification and user-friendly messages
- **Configuration**: Externalized configuration with sensible defaults
- **Logging**: Structured logging with context and configurable levels
- **Documentation**: Built-in help system and usage information
- **Data Validation**: Input validation with clear error messages
- **Exit Codes**: Proper exit codes for shell script integration

## Extension Points

This CLI tool can be extended with:

- **Database Integration**: Replace file storage with real databases
- **Authentication**: Add user authentication and authorization
- **API Integration**: Connect to REST APIs or microservices
- **Advanced Parsing**: Integrate with libraries like Cobra or CLI
- **Shell Completion**: Add bash/zsh completion support
- **Interactive Mode**: Add interactive prompts and wizards

## Testing

The example includes built-in demonstration that tests:

- Command parsing and execution
- Help system functionality
- Configuration loading
- Output formatting
- Error handling scenarios

Run the demonstration with:
```bash
go run main.go
```

For interactive testing, use specific commands:
```bash
go run main.go help
go run main.go create testuser test@example.com
```

This example provides a solid foundation for building professional command-line tools with Go and the lmcc-go-sdk framework. 