/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 */

package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

// ConfigChangeCallback 定义了配置变更时调用的回调函数类型。
// (ConfigChangeCallback defines the type for callback functions invoked on configuration change.)
// 这个回调接收 Viper 实例和解码后的整个配置对象。
// (This callback receives the Viper instance and the decoded entire configuration object.)
type ConfigChangeCallback func(v *viper.Viper, cfg any) error

// configManager 封装了 Viper 实例、配置对象和回调逻辑。
// 这是包内部使用的结构体。
// (configManager encapsulates the Viper instance, config object, and callback logic.)
// (This is an internal struct used within the package.)
type configManager[T any] struct {
	v                   *viper.Viper
	cfg                 *T
	callbacks           []ConfigChangeCallback // 通用回调 (General callbacks)
	callbackMux         sync.RWMutex
	sectionCallbacks    map[string][]SectionChangeCallback // 特定节回调 (Section-specific callbacks)
	sectionCallbacksMux sync.RWMutex
	options             Options // Use the Options type defined in options.go
	// watcher             *fsnotify.Watcher // 保持对 watcher 的引用，以便可以停止它 (Keep a reference to the watcher so it can be stopped)
	// watchStopper      chan struct{}     // 用于停止监视 goroutine 的通道 (Channel to stop the watch goroutine)
}

// newConfigManager 创建并初始化一个新的 configManager 实例。
// (newConfigManager creates and initializes a new configManager instance.)
// Parameters:
//   cfg: 指向用户配置结构体的指针。
//        (Pointer to the user's configuration struct.)
//   opts: 应用于配置加载的选项。
//         (Options to apply for configuration loading.)
// Returns:
//   *configManager[T]: 指向新创建的配置管理器的指针。
//                      (Pointer to the newly created config manager.)
func newConfigManager[T any](cfg *T, opts ...Option) *configManager[T] {
	// 应用默认选项和用户提供的选项
	// (Apply default options and user-provided options)
	appliedOptions := defaultOptions // Start with defaults (defined in options.go)
	for _, opt := range opts {
		opt(&appliedOptions)
	}
	return &configManager[T]{
		v:                viper.New(),
		cfg:              cfg,
		options:          appliedOptions, // Use the processed options
		sectionCallbacks: make(map[string][]SectionChangeCallback),
		// watchStopper:     make(chan struct{}), // 初始化停止通道 (Initialize stop channel)
	}
}

// RegisterCallback 注册一个配置变更回调函数，该函数将在配置通过热重载更新时被调用。
// (RegisterCallback registers a configuration change callback function, which will be invoked when the configuration is updated via hot-reload.)
// Parameters:
//   callback: 要注册的回调函数 (ConfigChangeCallback)。
//             (The callback function (ConfigChangeCallback) to register.)
func (cm *configManager[T]) RegisterCallback(callback func(v *viper.Viper, cfg any) error) { // Ensure signature matches interface
	cm.callbackMux.Lock()
	defer cm.callbackMux.Unlock()
	cm.callbacks = append(cm.callbacks, callback)
	log.Printf("Info: Registered a general configuration change callback.")
}

// RegisterSectionChangeCallback 注册一个针对特定配置节的变更回调函数。
// (RegisterSectionChangeCallback registers a callback function for changes in a specific configuration section.)
// Parameters:
//   sectionKey: 要监视的配置节的键名 (例如 "log", "database")。
//               (The key of the configuration section to watch (e.g., "log", "database").)
//   callback:   当配置节变更时调用的回调函数 (SectionChangeCallback)。
//               (The callback function (SectionChangeCallback) to invoke when the section changes.)
func (cm *configManager[T]) RegisterSectionChangeCallback(sectionKey string, callback SectionChangeCallback) {
	cm.sectionCallbacksMux.Lock()
	defer cm.sectionCallbacksMux.Unlock()
	cm.sectionCallbacks[sectionKey] = append(cm.sectionCallbacks[sectionKey], callback)
	log.Printf("Info: Registered a configuration change callback for section [%s].", sectionKey)
}

// notifyCallbacks 在配置变更后通知所有注册的回调函数。
// (notifyCallbacks notifies all registered callback functions after a configuration change.)
func (cm *configManager[T]) notifyCallbacks() {
	// 通知通用回调 (Notify general callbacks)
	cm.callbackMux.RLock()
	// 创建副本以避免在回调执行期间持有锁 (Create a copy to avoid holding lock during callback execution)
	currentCallbacks := make([]ConfigChangeCallback, len(cm.callbacks))
	copy(currentCallbacks, cm.callbacks)
	cm.callbackMux.RUnlock()

	if len(currentCallbacks) > 0 {
		log.Printf("Info: Notifying %d general callback(s) about configuration change...", len(currentCallbacks))
		for i, callback := range currentCallbacks {
			if err := callback(cm.v, cm.cfg); err != nil {
				log.Printf("Error executing general configuration change callback %d: %v", i+1, err)
			}
		}
	}

	// 通知特定节回调 (Notify section-specific callbacks)
	cm.sectionCallbacksMux.RLock()
	// 创建深拷贝以避免在回调执行期间持有锁或回调修改映射 (Create a deep copy to avoid holding lock or callback modifying map during execution)
	currentSectionCallbacks := make(map[string][]SectionChangeCallback)
	for key, callbacksSlice := range cm.sectionCallbacks {
		copiedSlice := make([]SectionChangeCallback, len(callbacksSlice))
		copy(copiedSlice, callbacksSlice)
		currentSectionCallbacks[key] = copiedSlice
	}
	cm.sectionCallbacksMux.RUnlock()

	if len(currentSectionCallbacks) > 0 {
		log.Printf("Info: Notifying section-specific callback(s) about configuration change...")
		for sectionKey, callbacksSlice := range currentSectionCallbacks {
			// 理论上，我们应该只在特定节实际更改时调用这些回调。
			// 但 Viper 的 OnConfigChange 不提供此类粒度信息。
			// 因此，我们通知所有已注册的特定节回调，并让它们自行处理。
			// (Theoretically, we should only call these if the specific section actually changed.
			// However, Viper's OnConfigChange doesn't provide that granular info.
			// So, we notify all registered section callbacks and let them handle it.)
			log.Printf("Info: Notifying %d callback(s) for section [%s]...", len(callbacksSlice), sectionKey)
			for i, callback := range callbacksSlice {
				if err := callback(cm.v); err != nil {
					log.Printf("Error executing configuration change callback for section [%s], callback %d: %v", sectionKey, i+1, err)
				}
			}
		}
	}
}

// GetViperInstance 返回 configManager 内部使用的 Viper 实例。
// (GetViperInstance returns the internal Viper instance used by the configManager.)
// Returns:
//   *viper.Viper: 内部的 Viper 实例。
//                 (The internal Viper instance.)
func (cm *configManager[T]) GetViperInstance() *viper.Viper {
	return cm.v
}

// StopWatch 停止文件监视器。如果未启用热重载，则此操作无效。
// (StopWatch stops the file watcher. No-op if hot-reloading was not enabled.)
// func (cm *configManager[T]) StopWatch() {
// 	if cm.options.hotReloadEnabled && cm.watcher != nil {
// 		log.Println("Info: Stopping configuration file watcher...")
// 		close(cm.watchStopper) // Signal the watch goroutine to stop
// 		cm.watcher.Close()
// 		log.Println("Info: Configuration file watcher stopped.")
// 	}
// } 