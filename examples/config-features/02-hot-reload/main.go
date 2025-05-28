/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * Hot-reload configuration example demonstrating real-time config updates.
 */

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/errors"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/log"
	"github.com/spf13/viper"
)

// HotReloadConfig 支持热重载的配置结构体
// (HotReloadConfig represents a configuration structure that supports hot-reload)
type HotReloadConfig struct {
	config.Config                     // 嵌入SDK基础配置 (Embed SDK base configuration)
	
	// 应用配置 (Application configuration)
	App *AppConfig `mapstructure:"app"`
	
	// 服务器配置 (Server configuration)  
	Server *ServerConfig `mapstructure:"server"`
	
	// 功能开关 (Feature flags)
	Features *FeatureConfig `mapstructure:"features"`
	
	// 限流配置 (Rate limiting configuration)
	RateLimit *RateLimitConfig `mapstructure:"rate_limit"`
}

// AppConfig 应用配置
// (AppConfig represents application configuration)
type AppConfig struct {
	Name        string        `mapstructure:"name" default:"HotReloadExample"`
	Version     string        `mapstructure:"version" default:"1.0.0"`
	Environment string        `mapstructure:"environment" default:"development"`
	Debug       bool          `mapstructure:"debug" default:"true"`
	Timeout     time.Duration `mapstructure:"timeout" default:"30s"`
	Workers     int          `mapstructure:"workers" default:"4"`
}

// ServerConfig 服务器配置
// (ServerConfig represents server configuration)
type ServerConfig struct {
	Host           string        `mapstructure:"host" default:"localhost"`
	Port           int          `mapstructure:"port" default:"8080"`
	ReadTimeout    time.Duration `mapstructure:"read_timeout" default:"10s"`
	WriteTimeout   time.Duration `mapstructure:"write_timeout" default:"10s"`
	MaxHeaderBytes int          `mapstructure:"max_header_bytes" default:"1048576"`
}

// FeatureConfig 功能开关配置
// (FeatureConfig represents feature flags configuration)
type FeatureConfig struct {
	EnableAuth     bool `mapstructure:"enable_auth" default:"true"`
	EnableMetrics  bool `mapstructure:"enable_metrics" default:"true"`
	EnableTracing  bool `mapstructure:"enable_tracing" default:"false"`
	EnableCaching  bool `mapstructure:"enable_caching" default:"true"`
	MaintenanceMode bool `mapstructure:"maintenance_mode" default:"false"`
}

// RateLimitConfig 限流配置
// (RateLimitConfig represents rate limiting configuration)
type RateLimitConfig struct {
	Enabled     bool          `mapstructure:"enabled" default:"true"`
	RequestsPerSecond int    `mapstructure:"requests_per_second" default:"100"`
	BurstSize   int          `mapstructure:"burst_size" default:"200"`
	WindowSize  time.Duration `mapstructure:"window_size" default:"1m"`
}

// ConfigWatcher 配置监视器
// (ConfigWatcher manages configuration hot-reload)
type ConfigWatcher struct {
	cfg      *HotReloadConfig
	manager  config.Manager
	logger   log.Logger
	mu       sync.RWMutex
	callbacks []ConfigCallback
}

// ConfigCallback 配置变更回调函数类型
// (ConfigCallback is the type for configuration change callbacks)
type ConfigCallback func(oldCfg, newCfg *HotReloadConfig) error

// NewConfigWatcher 创建新的配置监视器
// (NewConfigWatcher creates a new configuration watcher)
func NewConfigWatcher(configFile string, logger log.Logger) (*ConfigWatcher, error) {
	var cfg HotReloadConfig
	
	// 加载配置并启用监视 (Load configuration and enable watching)
	manager, err := config.LoadConfigAndWatch(&cfg,
		config.WithConfigFile(configFile, "yaml"),
		config.WithEnvPrefix("HOTRELOAD"),
		config.WithEnvVarOverride(true),
		config.WithHotReload(true),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load and watch configuration")
	}
	
	watcher := &ConfigWatcher{
		cfg:     &cfg,
		manager: manager,
		logger:  logger,
	}
	
	// 注册配置变更处理器 (Register configuration change handler)
	manager.RegisterCallback(watcher.handleConfigChange)
	
	logger.Info("Configuration watcher initialized successfully")
	return watcher, nil
}

// GetConfig 获取当前配置（线程安全）
// (GetConfig returns current configuration thread-safely)
func (cw *ConfigWatcher) GetConfig() *HotReloadConfig {
	cw.mu.RLock()
	defer cw.mu.RUnlock()
	
	// 返回配置的副本以避免并发修改 (Return copy to avoid concurrent modification)
	return cw.copyConfig(cw.cfg)
}

// copyConfig 深拷贝配置
// (copyConfig creates a deep copy of configuration)
func (cw *ConfigWatcher) copyConfig(cfg *HotReloadConfig) *HotReloadConfig {
	newCfg := *cfg
	if cfg.App != nil {
		app := *cfg.App
		newCfg.App = &app
	}
	if cfg.Server != nil {
		server := *cfg.Server
		newCfg.Server = &server
	}
	if cfg.Features != nil {
		features := *cfg.Features
		newCfg.Features = &features
	}
	if cfg.RateLimit != nil {
		rateLimit := *cfg.RateLimit
		newCfg.RateLimit = &rateLimit
	}
	return &newCfg
}

// RegisterCallback 注册配置变更回调
// (RegisterCallback registers a configuration change callback)
func (cw *ConfigWatcher) RegisterCallback(callback ConfigCallback) {
	cw.mu.Lock()
	defer cw.mu.Unlock()
	cw.callbacks = append(cw.callbacks, callback)
}

// handleConfigChange 处理配置变更
// (handleConfigChange handles configuration changes)
func (cw *ConfigWatcher) handleConfigChange(v *viper.Viper, newConfig any) error {
	cw.logger.Info("Configuration change detected, processing...")
	
	newCfg, ok := newConfig.(*HotReloadConfig)
	if !ok {
		return errors.New("invalid configuration type")
	}
	
	// 验证新配置 (Validate new configuration)
	if err := cw.validateConfig(newCfg); err != nil {
		cw.logger.Errorf("Configuration validation failed: %v", err)
		return errors.Wrap(err, "new configuration validation failed")
	}
	
	// 获取旧配置 (Get old configuration)
	cw.mu.Lock()
	oldCfg := cw.copyConfig(cw.cfg)
	cw.cfg = newCfg
	cw.mu.Unlock()
	
	// 执行回调函数 (Execute callbacks)
	for i, callback := range cw.callbacks {
		if err := callback(oldCfg, newCfg); err != nil {
			cw.logger.Errorf("Config callback %d failed: %v", i, err)
			// 继续执行其他回调，不因为一个失败而停止 (Continue with other callbacks)
		}
	}
	
	cw.logger.Info("Configuration updated successfully")
	return nil
}

// validateConfig 验证配置
// (validateConfig validates configuration)
func (cw *ConfigWatcher) validateConfig(cfg *HotReloadConfig) error {
	if cfg.Server != nil {
		if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
			return errors.Errorf("invalid server port: %d", cfg.Server.Port)
		}
	}
	
	if cfg.App != nil {
		if cfg.App.Workers < 1 {
			return errors.Errorf("invalid worker count: %d", cfg.App.Workers)
		}
	}
	
	if cfg.RateLimit != nil && cfg.RateLimit.Enabled {
		if cfg.RateLimit.RequestsPerSecond <= 0 {
			return errors.New("requests_per_second must be positive")
		}
		if cfg.RateLimit.BurstSize <= 0 {
			return errors.New("burst_size must be positive")
		}
	}
	
	return nil
}

// Close 关闭配置监视器
// (Close closes the configuration watcher)
func (cw *ConfigWatcher) Close() error {
	// Config manager doesn't need explicit closing
	return nil
}

// MockService 模拟服务，用于演示配置热重载的效果
// (MockService simulates a service to demonstrate hot-reload effects)
type MockService struct {
	watcher *ConfigWatcher
	logger  log.Logger
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewMockService 创建模拟服务
// (NewMockService creates a mock service)
func NewMockService(watcher *ConfigWatcher, logger log.Logger) *MockService {
	ctx, cancel := context.WithCancel(context.Background())
	
	service := &MockService{
		watcher: watcher,
		logger:  logger,
		ctx:     ctx,
		cancel:  cancel,
	}
	
	// 注册配置变更回调 (Register configuration change callback)
	watcher.RegisterCallback(service.onConfigChange)
	
	return service
}

// onConfigChange 处理配置变更
// (onConfigChange handles configuration changes)
func (ms *MockService) onConfigChange(oldCfg, newCfg *HotReloadConfig) error {
	ms.logger.Info("Service received configuration change notification")
	
	// 比较和报告变更 (Compare and report changes)
	changes := ms.detectChanges(oldCfg, newCfg)
	for _, change := range changes {
		ms.logger.Infof("Configuration change: %s", change)
	}
	
	// 模拟服务重新配置 (Simulate service reconfiguration)
	ms.logger.Info("Reconfiguring service with new settings...")
	
	return nil
}

// detectChanges 检测配置变更
// (detectChanges detects configuration changes)
func (ms *MockService) detectChanges(oldCfg, newCfg *HotReloadConfig) []string {
	var changes []string
	
	if oldCfg.App != nil && newCfg.App != nil {
		if oldCfg.App.Debug != newCfg.App.Debug {
			changes = append(changes, fmt.Sprintf("Debug mode: %t → %t", oldCfg.App.Debug, newCfg.App.Debug))
		}
		if oldCfg.App.Workers != newCfg.App.Workers {
			changes = append(changes, fmt.Sprintf("Worker count: %d → %d", oldCfg.App.Workers, newCfg.App.Workers))
		}
		if oldCfg.App.Timeout != newCfg.App.Timeout {
			changes = append(changes, fmt.Sprintf("Timeout: %v → %v", oldCfg.App.Timeout, newCfg.App.Timeout))
		}
	}
	
	if oldCfg.Server != nil && newCfg.Server != nil {
		if oldCfg.Server.Port != newCfg.Server.Port {
			changes = append(changes, fmt.Sprintf("Server port: %d → %d", oldCfg.Server.Port, newCfg.Server.Port))
		}
	}
	
	if oldCfg.Features != nil && newCfg.Features != nil {
		if oldCfg.Features.MaintenanceMode != newCfg.Features.MaintenanceMode {
			changes = append(changes, fmt.Sprintf("Maintenance mode: %t → %t", 
				oldCfg.Features.MaintenanceMode, newCfg.Features.MaintenanceMode))
		}
	}
	
	if oldCfg.RateLimit != nil && newCfg.RateLimit != nil {
		if oldCfg.RateLimit.RequestsPerSecond != newCfg.RateLimit.RequestsPerSecond {
			changes = append(changes, fmt.Sprintf("Rate limit RPS: %d → %d", 
				oldCfg.RateLimit.RequestsPerSecond, newCfg.RateLimit.RequestsPerSecond))
		}
	}
	
	return changes
}

// Start 启动模拟服务
// (Start starts the mock service)
func (ms *MockService) Start() {
	ms.logger.Info("Mock service started")
	
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ms.ctx.Done():
			ms.logger.Info("Mock service stopped")
			return
		case <-ticker.C:
			// 定期打印当前配置状态 (Periodically print current configuration status)
			cfg := ms.watcher.GetConfig()
			ms.printCurrentStatus(cfg)
		}
	}
}

// printCurrentStatus 打印当前状态
// (printCurrentStatus prints current status)
func (ms *MockService) printCurrentStatus(cfg *HotReloadConfig) {
	ms.logger.Infow("Current service status",
		"app_name", cfg.App.Name,
		"debug_mode", cfg.App.Debug,
		"workers", cfg.App.Workers,
		"server_port", cfg.Server.Port,
		"maintenance_mode", cfg.Features.MaintenanceMode,
		"rate_limit_rps", cfg.RateLimit.RequestsPerSecond,
	)
}

// Stop 停止模拟服务
// (Stop stops the mock service)
func (ms *MockService) Stop() {
	ms.cancel()
}

func main() {
	fmt.Println("=== Hot-Reload Configuration Example ===")
	fmt.Println("This example demonstrates real-time configuration updates.")
	fmt.Println("Try modifying config.yaml while the program is running!")
	fmt.Println()
	
	// 1. 初始化日志 (Initialize logging)
	logOpts := log.NewOptions()
	logOpts.Level = "info"
	logOpts.Format = "text"
	logOpts.EnableColor = true
	log.Init(logOpts)
	logger := log.Std()
	
	// 2. 创建配置监视器 (Create configuration watcher)
	fmt.Println("Initializing configuration watcher...")
	watcher, err := NewConfigWatcher("config.yaml", logger)
	if err != nil {
		logger.Errorf("Failed to create config watcher: %v", err)
		if coder := errors.GetCoder(err); coder != nil {
			fmt.Printf("Error Code: %d, Type: %s\n", coder.Code(), coder.String())
		}
		os.Exit(1)
	}
	defer watcher.Close()
	
	// 3. 创建模拟服务 (Create mock service)
	fmt.Println("Creating mock service...")
	service := NewMockService(watcher, logger)
	
	// 4. 启动服务 (Start service)
	fmt.Println("Starting mock service...")
	go service.Start()
	
	// 5. 打印初始配置 (Print initial configuration)
	fmt.Println()
	fmt.Println("=== Initial Configuration ===")
	cfg := watcher.GetConfig()
	printConfigSummary(cfg)
	
	// 6. 提示用户如何测试热重载 (Prompt user how to test hot-reload)
	fmt.Println()
	fmt.Println("=== Hot-Reload Testing ===")
	fmt.Println("The service is now running and watching for configuration changes.")
	fmt.Println("To test hot-reload:")
	fmt.Println("1. Open another terminal")
	fmt.Println("2. Edit config.yaml in this directory")
	fmt.Println("3. Save the file")
	fmt.Println("4. Watch this program detect and apply the changes!")
	fmt.Println()
	fmt.Println("Example changes to try:")
	fmt.Println("- Change app.debug from true to false")
	fmt.Println("- Modify server.port (e.g., 8080 → 8081)")
	fmt.Println("- Toggle features.maintenance_mode")
	fmt.Println("- Adjust rate_limit.requests_per_second")
	fmt.Println()
	fmt.Println("Press Ctrl+C to exit")
	fmt.Println()
	
	// 7. 等待信号 (Wait for signals)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	<-sigChan
	fmt.Println("\nShutdown signal received, stopping...")
	
	// 8. 优雅关闭 (Graceful shutdown)
	service.Stop()
	logger.Info("Application shutdown completed")
}

// printConfigSummary 打印配置摘要
// (printConfigSummary prints configuration summary)
func printConfigSummary(cfg *HotReloadConfig) {
	fmt.Printf("Application:\n")
	fmt.Printf("  Name: %s\n", cfg.App.Name)
	fmt.Printf("  Debug: %t\n", cfg.App.Debug)
	fmt.Printf("  Workers: %d\n", cfg.App.Workers)
	fmt.Printf("  Timeout: %v\n", cfg.App.Timeout)
	
	fmt.Printf("Server:\n")
	fmt.Printf("  Host: %s\n", cfg.Server.Host)
	fmt.Printf("  Port: %d\n", cfg.Server.Port)
	
	fmt.Printf("Features:\n")
	fmt.Printf("  Auth: %t\n", cfg.Features.EnableAuth)
	fmt.Printf("  Metrics: %t\n", cfg.Features.EnableMetrics)
	fmt.Printf("  Maintenance Mode: %t\n", cfg.Features.MaintenanceMode)
	
	fmt.Printf("Rate Limit:\n")
	fmt.Printf("  Enabled: %t\n", cfg.RateLimit.Enabled)
	fmt.Printf("  RPS: %d\n", cfg.RateLimit.RequestsPerSecond)
	fmt.Printf("  Burst: %d\n", cfg.RateLimit.BurstSize)
} 