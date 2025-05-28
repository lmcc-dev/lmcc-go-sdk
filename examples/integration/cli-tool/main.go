/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * CLI tool example demonstrating command-line interface with subcommands and configuration.
 */

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
)

// CLIConfig CLI工具配置
// (CLIConfig represents CLI tool configuration)
type CLIConfig struct {
	App struct {
		Name        string `yaml:"name" default:"user-cli"`
		Version     string `yaml:"version" default:"v1.0.0"`
		Description string `yaml:"description" default:"User management CLI tool"`
	} `yaml:"app"`

	Database struct {
		Type string `yaml:"type" default:"file"`
		Path string `yaml:"path" default:"./users.json"`
	} `yaml:"database"`

	Output struct {
		Format string `yaml:"format" default:"table"`
		Quiet  bool   `yaml:"quiet" default:"false"`
		Color  bool   `yaml:"color" default:"true"`
	} `yaml:"output"`

	Logging struct {
		Level       string   `yaml:"level" default:"info"`
		Format      string   `yaml:"format" default:"text"`
		OutputPaths []string `yaml:"output_paths"`
	} `yaml:"logging"`
}

// User 用户模型
// (User represents user model)
type User struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Name     string    `json:"name"`
	Status   string    `json:"status"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

// Command 命令接口
// (Command represents command interface)
type Command interface {
	Name() string
	Description() string
	Usage() string
	Execute(ctx context.Context, args []string) error
}

// CLI CLI工具主结构
// (CLI represents the main CLI tool structure)
type CLI struct {
	config   *CLIConfig
	logger   log.Logger
	storage  *FileStorage
	commands map[string]Command
}

// NewCLI 创建CLI工具
// (NewCLI creates a new CLI tool)
func NewCLI(cfg *CLIConfig) *CLI {
	// 设置日志默认值 (Set logging defaults)
	if len(cfg.Logging.OutputPaths) == 0 {
		cfg.Logging.OutputPaths = []string{"stdout"}
	}

	// 初始化日志 (Initialize logging)
	opts := log.NewOptions()
	opts.Level = cfg.Logging.Level
	opts.Format = cfg.Logging.Format
	opts.EnableColor = cfg.Output.Color && cfg.Logging.Format == "text"
	opts.DisableCaller = true
	opts.DisableStacktrace = true
	opts.OutputPaths = cfg.Logging.OutputPaths

	log.Init(opts)

	logger := log.Std().WithValues(
		"app", cfg.App.Name,
		"version", cfg.App.Version,
		"component", "cli")

	// 初始化存储 (Initialize storage)
	storage := NewFileStorage(cfg.Database.Path, logger)

	cli := &CLI{
		config:   cfg,
		logger:   logger,
		storage:  storage,
		commands: make(map[string]Command),
	}

	// 注册命令 (Register commands)
	cli.registerCommands()

	if !cfg.Output.Quiet {
		logger.Infow("CLI tool initialized",
			"app_name", cfg.App.Name,
			"app_version", cfg.App.Version,
			"database_path", cfg.Database.Path,
			"output_format", cfg.Output.Format)
	}

	return cli
}

// registerCommands 注册所有命令
// (registerCommands registers all commands)
func (c *CLI) registerCommands() {
	c.commands["create"] = NewCreateCommand(c)
	c.commands["list"] = NewListCommand(c)
	c.commands["get"] = NewGetCommand(c)
	c.commands["update"] = NewUpdateCommand(c)
	c.commands["delete"] = NewDeleteCommand(c)
	c.commands["search"] = NewSearchCommand(c)
	c.commands["export"] = NewExportCommand(c)
	c.commands["import"] = NewImportCommand(c)
	c.commands["help"] = NewHelpCommand(c)
	c.commands["version"] = NewVersionCommand(c)
}

// Run 运行CLI工具
// (Run executes the CLI tool)
func (c *CLI) Run(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return c.commands["help"].Execute(ctx, []string{})
	}

	cmdName := args[0]
	cmdArgs := args[1:]

	if cmd, exists := c.commands[cmdName]; exists {
		return cmd.Execute(ctx, cmdArgs)
	}

	return errors.Errorf("unknown command: %s", cmdName)
}

// CreateCommand 创建用户命令
// (CreateCommand creates user command)
type CreateCommand struct {
	cli *CLI
}

// NewCreateCommand 创建CreateCommand
// (NewCreateCommand creates a new CreateCommand)
func NewCreateCommand(cli *CLI) *CreateCommand {
	return &CreateCommand{cli: cli}
}

func (c *CreateCommand) Name() string        { return "create" }
func (c *CreateCommand) Description() string { return "Create a new user" }
func (c *CreateCommand) Usage() string {
	return "create <username> <email> [--name <name>] [--status <status>]"
}

func (c *CreateCommand) Execute(ctx context.Context, args []string) error {
	if len(args) < 2 {
		return errors.New("usage: " + c.Usage())
	}

	username := args[0]
	email := args[1]

	// 解析可选参数 (Parse optional arguments)
	name := ""
	status := "active"

	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "--name":
			if i+1 < len(args) {
				name = args[i+1]
				i++
			}
		case "--status":
			if i+1 < len(args) {
				status = args[i+1]
				i++
			}
		}
	}

	if name == "" {
		name = username
	}

	c.cli.logger.Infow("Creating user",
		"username", username,
		"email", email,
		"name", name,
		"status", status)

	user := &User{
		ID:       generateID(),
		Username: username,
		Email:    email,
		Name:     name,
		Status:   status,
		Created:  time.Now(),
		Updated:  time.Now(),
	}

	if err := c.cli.storage.CreateUser(ctx, user); err != nil {
		c.cli.logger.Errorw("Failed to create user", "error", err)
		return err
	}

	if !c.cli.config.Output.Quiet {
		c.printUser(user)
		fmt.Printf("✅ User '%s' created successfully with ID: %s\n", username, user.ID)
	}

	return nil
}

// printUser 打印用户信息
// (printUser prints user information)
func (c *CreateCommand) printUser(user *User) {
	switch c.cli.config.Output.Format {
	case "json":
		data, _ := json.MarshalIndent(user, "", "  ")
		fmt.Println(string(data))
	case "table":
		fmt.Printf("┌────────────┬─────────────────────────────────────┐\n")
		fmt.Printf("│ Field      │ Value                               │\n")
		fmt.Printf("├────────────┼─────────────────────────────────────┤\n")
		fmt.Printf("│ ID         │ %-35s │\n", user.ID)
		fmt.Printf("│ Username   │ %-35s │\n", user.Username)
		fmt.Printf("│ Email      │ %-35s │\n", user.Email)
		fmt.Printf("│ Name       │ %-35s │\n", user.Name)
		fmt.Printf("│ Status     │ %-35s │\n", user.Status)
		fmt.Printf("│ Created    │ %-35s │\n", user.Created.Format("2006-01-02 15:04:05"))
		fmt.Printf("└────────────┴─────────────────────────────────────┘\n")
	default:
		fmt.Printf("ID: %s, Username: %s, Email: %s, Name: %s, Status: %s\n",
			user.ID, user.Username, user.Email, user.Name, user.Status)
	}
}

// HelpCommand 帮助命令
// (HelpCommand help command)
type HelpCommand struct {
	cli *CLI
}

func NewHelpCommand(cli *CLI) *HelpCommand { return &HelpCommand{cli: cli} }

func (c *HelpCommand) Name() string        { return "help" }
func (c *HelpCommand) Description() string { return "Show help information" }
func (c *HelpCommand) Usage() string       { return "help [command]" }

func (c *HelpCommand) Execute(ctx context.Context, args []string) error {
	if len(args) > 0 {
		// 显示特定命令帮助 (Show specific command help)
		cmdName := args[0]
		if cmd, exists := c.cli.commands[cmdName]; exists {
			fmt.Printf("Command: %s\n", cmd.Name())
			fmt.Printf("Description: %s\n", cmd.Description())
			fmt.Printf("Usage: %s %s\n", c.cli.config.App.Name, cmd.Usage())
			return nil
		}
		return errors.Errorf("unknown command: %s", cmdName)
	}

	// 显示总体帮助 (Show general help)
	fmt.Printf("%s - %s\n", c.cli.config.App.Name, c.cli.config.App.Description)
	fmt.Printf("Version: %s\n\n", c.cli.config.App.Version)
	fmt.Println("Available commands:")

	// 按字母顺序排序命令 (Sort commands alphabetically)
	var cmdNames []string
	for name := range c.cli.commands {
		cmdNames = append(cmdNames, name)
	}
	sort.Strings(cmdNames)

	for _, name := range cmdNames {
		cmd := c.cli.commands[name]
		fmt.Printf("  %-10s %s\n", cmd.Name(), cmd.Description())
	}

	fmt.Printf("\nUse '%s help <command>' for more information about a command.\n", c.cli.config.App.Name)
	return nil
}

// VersionCommand 版本命令
// (VersionCommand version command)
type VersionCommand struct {
	cli *CLI
}

func NewVersionCommand(cli *CLI) *VersionCommand { return &VersionCommand{cli: cli} }

func (c *VersionCommand) Name() string        { return "version" }
func (c *VersionCommand) Description() string { return "Show version information" }
func (c *VersionCommand) Usage() string       { return "version" }

func (c *VersionCommand) Execute(ctx context.Context, args []string) error {
	fmt.Printf("%s version %s\n", c.cli.config.App.Name, c.cli.config.App.Version)
	fmt.Printf("Built with lmcc-go-sdk\n")
	return nil
}

// ListCommand 列出用户命令
// (ListCommand lists users command)
type ListCommand struct {
	cli *CLI
}

func NewListCommand(cli *CLI) *ListCommand { return &ListCommand{cli: cli} }

func (c *ListCommand) Name() string        { return "list" }
func (c *ListCommand) Description() string { return "List all users" }
func (c *ListCommand) Usage() string       { return "list [--status <status>] [--limit <n>]" }

func (c *ListCommand) Execute(ctx context.Context, args []string) error {
	c.cli.logger.Infow("Listing users")
	fmt.Println("📋 No users found (demo implementation)")
	return nil
}

// GetCommand 获取用户命令
// (GetCommand gets user command)
type GetCommand struct {
	cli *CLI
}

func NewGetCommand(cli *CLI) *GetCommand { return &GetCommand{cli: cli} }

func (c *GetCommand) Name() string        { return "get" }
func (c *GetCommand) Description() string { return "Get user by ID or username" }
func (c *GetCommand) Usage() string       { return "get <id_or_username>" }

func (c *GetCommand) Execute(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("usage: " + c.Usage())
	}
	c.cli.logger.Infow("Getting user", "identifier", args[0])
	fmt.Printf("👤 User not found: %s (demo implementation)\n", args[0])
	return nil
}

// UpdateCommand 更新用户命令
// (UpdateCommand updates user command)
type UpdateCommand struct {
	cli *CLI
}

func NewUpdateCommand(cli *CLI) *UpdateCommand { return &UpdateCommand{cli: cli} }

func (c *UpdateCommand) Name() string        { return "update" }
func (c *UpdateCommand) Description() string { return "Update user information" }
func (c *UpdateCommand) Usage() string {
	return "update <id_or_username> [--email <email>] [--name <name>] [--status <status>]"
}

func (c *UpdateCommand) Execute(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("usage: " + c.Usage())
	}
	c.cli.logger.Infow("Updating user", "identifier", args[0])
	fmt.Printf("✏️  User update simulated: %s (demo implementation)\n", args[0])
	return nil
}

// DeleteCommand 删除用户命令
// (DeleteCommand deletes user command)
type DeleteCommand struct {
	cli *CLI
}

func NewDeleteCommand(cli *CLI) *DeleteCommand { return &DeleteCommand{cli: cli} }

func (c *DeleteCommand) Name() string        { return "delete" }
func (c *DeleteCommand) Description() string { return "Delete user" }
func (c *DeleteCommand) Usage() string       { return "delete <id_or_username> [--force]" }

func (c *DeleteCommand) Execute(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("usage: " + c.Usage())
	}
	c.cli.logger.Infow("Deleting user", "identifier", args[0])
	fmt.Printf("🗑️  User deletion simulated: %s (demo implementation)\n", args[0])
	return nil
}

// SearchCommand 搜索用户命令
// (SearchCommand searches users command)
type SearchCommand struct {
	cli *CLI
}

func NewSearchCommand(cli *CLI) *SearchCommand { return &SearchCommand{cli: cli} }

func (c *SearchCommand) Name() string        { return "search" }
func (c *SearchCommand) Description() string { return "Search users by keyword" }
func (c *SearchCommand) Usage() string       { return "search <keyword> [--field <field>]" }

func (c *SearchCommand) Execute(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("usage: " + c.Usage())
	}
	c.cli.logger.Infow("Searching users", "keyword", args[0])
	fmt.Printf("🔍 No users found matching: %s (demo implementation)\n", args[0])
	return nil
}

// ExportCommand 导出用户命令
// (ExportCommand exports users command)
type ExportCommand struct {
	cli *CLI
}

func NewExportCommand(cli *CLI) *ExportCommand { return &ExportCommand{cli: cli} }

func (c *ExportCommand) Name() string        { return "export" }
func (c *ExportCommand) Description() string { return "Export users to file" }
func (c *ExportCommand) Usage() string       { return "export <filename> [--format <json|csv>]" }

func (c *ExportCommand) Execute(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("usage: " + c.Usage())
	}
	c.cli.logger.Infow("Exporting users", "filename", args[0])
	fmt.Printf("📤 Export simulated to: %s (demo implementation)\n", args[0])
	return nil
}

// ImportCommand 导入用户命令
// (ImportCommand imports users command)
type ImportCommand struct {
	cli *CLI
}

func NewImportCommand(cli *CLI) *ImportCommand { return &ImportCommand{cli: cli} }

func (c *ImportCommand) Name() string        { return "import" }
func (c *ImportCommand) Description() string { return "Import users from file" }
func (c *ImportCommand) Usage() string       { return "import <filename> [--merge]" }

func (c *ImportCommand) Execute(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("usage: " + c.Usage())
	}
	c.cli.logger.Infow("Importing users", "filename", args[0])
	fmt.Printf("📥 Import simulated from: %s (demo implementation)\n", args[0])
	return nil
}

// FileStorage 文件存储 (简化版)
// (FileStorage provides file-based storage - simplified version)
type FileStorage struct {
	path   string
	logger log.Logger
}

// NewFileStorage 创建文件存储
// (NewFileStorage creates new file storage)
func NewFileStorage(path string, logger log.Logger) *FileStorage {
	return &FileStorage{
		path:   path,
		logger: logger.WithValues("component", "storage"),
	}
}

// CreateUser 创建用户
// (CreateUser creates a user)
func (fs *FileStorage) CreateUser(ctx context.Context, user *User) error {
	fs.logger.Infow("User created", "user_id", user.ID, "username", user.Username)
	return nil // 简化实现
}

// generateID 生成唯一ID
// (generateID generates unique ID)
func generateID() string {
	return fmt.Sprintf("user_%d", time.Now().Unix())
}

// runDemo 运行CLI演示
// (runDemo runs CLI demonstration)
func runDemo() {
	fmt.Println("=== Running CLI Tool Demonstration ===")
	fmt.Println()

	// 模拟命令行参数 (Simulate command line arguments)
	tests := []struct {
		name string
		args []string
	}{
		{"Show help", []string{"help"}},
		{"Create user alice", []string{"create", "alice", "alice@example.com", "--name", "Alice Smith"}},
		{"Show help for create command", []string{"help", "create"}},
	}

	// 加载配置 (Load configuration)
	cfg := &CLIConfig{}
	if err := config.LoadConfig(cfg); err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		// 使用默认配置继续 (Continue with default configuration)
	}

	// 设置演示配置 (Set demo configuration)
	cfg.Database.Path = "./demo_users.json"
	cfg.Output.Quiet = false

	// 创建CLI实例 (Create CLI instance)
	cli := NewCLI(cfg)
	ctx := context.Background()

	// 运行测试命令 (Run test commands)
	for i, test := range tests {
		fmt.Printf("%d. %s:\n", i+1, test.name)
		fmt.Printf("   Command: %s %s\n", cfg.App.Name, strings.Join(test.args, " "))

		if err := cli.Run(ctx, test.args); err != nil {
			fmt.Printf("   ❌ Error: %v\n", err)
		} else {
			fmt.Printf("   ✅ Success\n")
		}
		fmt.Println()
	}

	fmt.Println("=== CLI Tool Demonstration Completed ===")
}

func main() {
	// 如果没有参数，运行演示 (If no arguments, run demonstration)
	if len(os.Args) <= 1 {
		runDemo()
		return
	}

	fmt.Println("=== CLI Tool Integration Example ===")
	fmt.Println("This example demonstrates a command-line interface with subcommands and configuration.")
	fmt.Println()

	// 加载配置 (Load configuration)
	cfg := &CLIConfig{}
	if err := config.LoadConfig(cfg); err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		// 使用默认配置继续 (Continue with default configuration)
	}

	// 创建CLI实例 (Create CLI instance)
	cli := NewCLI(cfg)

	// 运行CLI (Run CLI)
	ctx := context.Background()
	if err := cli.Run(ctx, os.Args[1:]); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
} 