/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * This file implements the configuration watcher for the log package.
 * It is responsible for listening to configuration changes from pkg/config
 * and reconfiguring the global logger accordingly.
 */

package log

import (
	"sync"

	"fmt"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/config"
	lmccerrors "github.com/lmcc-dev/lmcc-go-sdk/pkg/errors" // SDK errors 包 (SDK errors package)
	"github.com/spf13/viper"
)

// ConfigChangeCallback defines the function signature for configuration change callbacks.
// (ConfigChangeCallback 定义了配置变更回调的函数签名。)
// These callbacks are specific to the log package and are triggered after the global logger has been reconfigured.
// (这些回调特定于日志包，并在全局记录器重新配置后触发。)
type ConfigChangeCallback func(newOptions *Options) error

var (
	// callbacks stores the registered configuration change callbacks.
	// (callbacks 存储已注册的配置变更回调。)
	callbacks = make(map[string]ConfigChangeCallback)
	// callbacksMu protects the callbacks map during concurrent access.
	// (callbacksMu 在并发访问期间保护 callbacks 映射。)
	callbacksMu sync.RWMutex
	// nextLogCallbackID is used to generate unique IDs for callbacks.
	// (nextLogCallbackID 用于为回调生成唯一的ID。)
	nextLogCallbackID int64

	// currentProcessLogConfigChange holds the current function to process log configuration changes.
	// This allows for swapping the implementation, e.g., for testing.
	// (currentProcessLogConfigChange 保存当前处理日志配置更改的函数。
	// 这允许交换实现，例如，用于测试。)
	currentProcessLogConfigChange func(v *viper.Viper) error
)

// init registers the log package's configuration update handler with the config package.
// (init 将日志包的配置更新处理函数注册到配置包中。)
func init() {
	// Initialize currentProcessLogConfigChange with the default implementation.
	// (使用默认实现初始化 currentProcessLogConfigChange。)
	currentProcessLogConfigChange = defaultHandleGlobalLogConfigChange
	// Placeholder for registering with pkg/config
	// This will be implemented once pkg/config provides the necessary interface.
	// For now, we can simulate a config change or set up a basic mechanism.
	// config.RegisterCallback(handleGlobalLogConfigChange) // Example of future integration
}

// RegisterCallback registers a new callback function to be called when the log configuration changes.
// It returns a unique ID for the registered callback, which can be used to unregister it later.
// (RegisterCallback 注册一个新的回调函数，当日志配置发生变化时调用。
// 它返回已注册回调的唯一ID，该ID可用于以后注销它。)
func RegisterCallback(callback ConfigChangeCallback) string {
	callbacksMu.Lock()
	defer callbacksMu.Unlock()

	nextLogCallbackID++
	callbackID := fmt.Sprintf("log-callback-%d", nextLogCallbackID)
	callbacks[callbackID] = callback
	// TODO: Log registration (use a local logger or print for now if global logger not ready)
	// Info("Registered new log configuration callback", "id", callbackID)
	return callbackID
}

// UnregisterCallback removes a previously registered callback function using its ID.
// (UnregisterCallback 使用其ID删除先前注册的回调函数。)
func UnregisterCallback(id string) {
	callbacksMu.Lock()
	defer callbacksMu.Unlock()

	// S1033: unnecessary guard around call to delete (gosimple) is addressed by removing the if.
	// In Go, deleting a non-existent key from a map is a no-op and safe.
	// (在 Go 中，从 map 中删除不存在的键是空操作且安全的。)
		delete(callbacks, id)
	// The empty else branch (previously SA9003) is removed as the preceding if is removed.
	// TODO: Log unregistration if the key existed, or warn if it didn't if that's desired behavior.
		// Info("Unregistered log configuration callback", "id", id)
		// Warn("Attempted to unregister non-existent log configuration callback", "id", id)
}

// RegisterConfigHotReload 注册日志配置的热重载回调到配置管理器。
// (RegisterConfigHotReload registers the log config hot-reload callback to the configuration manager.)
// 此函数应在应用程序初始化期间，加载配置后调用。
// (This function should be called during application initialization, after loading the configuration.)
//
// Parameters:
//   cfgManager: 配置管理器实例 (config.Manager)，它提供了 RegisterSectionChangeCallback 方法。
//               (cfgManager: The configuration manager instance (config.Manager), which provides the RegisterSectionChangeCallback method.)
func RegisterConfigHotReload(cfgManager config.Manager) {
	// 使用 RegisterSectionChangeCallback 注册一个只关心 \"log\" 配置节的回调。
	// (Use RegisterSectionChangeCallback to register a callback that only cares about the \"log\" configuration section.)
	cfgManager.RegisterSectionChangeCallback("log", func(v *viper.Viper) error {
		// 调用可替换的配置处理函数。
		// (Call the swappable configuration processing function.)
		if currentProcessLogConfigChange == nil {
			// 这是一个防御性检查，理论上 currentProcessLogConfigChange 总是在 init 中被设置。
			// (This is a defensive check; currentProcessLogConfigChange should always be set in init.)
			Error("Log configuration change processor is not initialized.") // Log using existing global logger (or its last known state)
			// 使用 lmccerrors.NewWithCode 创建错误，以包含堆栈跟踪和错误码。
			// (Use lmccerrors.NewWithCode to create an error with a stack trace and error code.)
			return lmccerrors.NewWithCode(lmccerrors.ErrLogInternal, "log configuration change processor is not initialized")
		}
		return currentProcessLogConfigChange(v)
	})
	// Consider logging the successful registration.
	// Example: Info("Successfully registered log configuration hot-reload handler with config manager.")
}

// defaultHandleGlobalLogConfigChange is the default function that will be called by pkg/config when
// the relevant section of the configuration changes.
// (defaultHandleGlobalLogConfigChange 是当配置的相关部分发生更改时，pkg/config 将调用的默认函数。)
func defaultHandleGlobalLogConfigChange(v *viper.Viper) error {
	callbacksMu.RLock()
	defer callbacksMu.RUnlock()

	// 1. Parse the new log options from the viper instance.
	//    (从 viper 实例解析新的日志选项。)
	opts := NewOptions()
	if err := v.UnmarshalKey("log", opts); err != nil {
		// 使用 lmccerrors.Wrapf 包装错误，以添加堆栈跟踪、上下文和错误码。
		// (Wrap the error with lmccerrors.Wrapf to add stack trace, context, and error code.)
		Error("Failed to unmarshal new log configuration", "error", err) // Log with existing logger
		return lmccerrors.WithCode(
			lmccerrors.Wrap(err, "failed to unmarshal new log configuration from viper"),
			lmccerrors.ErrLogOptionInvalid,
		)
	}

	// 2. Validate the new options (e.g., level, format are valid)
	//    (验证新选项（例如，级别、格式是否有效）)
	err_list := opts.Validate() // Renamed from errs to err_list to avoid conflict with err from ReconfigureGlobalLogger
	if len(err_list) > 0 {
		Error("Validation failed for new log options", "errors", err_list, "options", opts) // Log the original errors slice for detail
		// 使用 ErrorGroup 将多个错误合并，然后使用 lmccerrors.Wrapf 添加上下文、错误码和堆栈跟踪。
		// (Use ErrorGroup to combine multiple errors, then lmccerrors.Wrapf to add context, error code and stack trace.)
		// 创建一个新的 ErrorGroup 来包含所有验证错误。
		// (Create a new ErrorGroup to contain all validation errors.)
		eg := lmccerrors.NewErrorGroup("log options validation failed") // The message for the group itself
		for _, validationErr := range err_list {
			eg.Add(validationErr)
		}
		// 使用 ErrLogOptionInvalid Coder 包装整个错误组。
		// (Wrap the entire error group with the ErrLogOptionInvalid Coder.)
		return lmccerrors.WithCode(
			lmccerrors.Wrap(eg, "one or more log options are invalid"),
			lmccerrors.ErrLogOptionInvalid,
		)
	}

	// 3. Apply the new configuration to the global logger.
	if err := ReconfigureGlobalLogger(opts); err != nil {
		// ReconfigureGlobalLogger already returns an lmccerror, but we wrap it again to provide specific context from this function.
		// (ReconfigureGlobalLogger 已返回 lmccerror，但我们再次包装它以提供此函数的特定上下文。)
		Error("Failed to reconfigure global logger with new options", "error", err) // Log with existing logger
		return lmccerrors.WithCode(
			lmccerrors.Wrap(err, "failed to apply new options to global logger"),
			lmccerrors.ErrLogReconfigure,
		)
	}

	Info("Global logger successfully reconfigured with new options.", "options", opts)

	// 4. Notify all registered callbacks about the change.
	//    Iterate over a copy of the callbacks in case a callback tries to unregister itself.
	currentCallbacks := make(map[string]ConfigChangeCallback)
	for id, cb := range callbacks {
		currentCallbacks[id] = cb
	}

	for id, callback := range currentCallbacks {
		if err := callback(opts); err != nil {
			Error("Error executing log configuration change callback", "callbackID", id, "error", err)
			// Decide if one callback error should stop others or just be logged.
		}
	}

	return nil
}

// TODO: Add a function in pkg/config that pkg/log can call to register 'handleGlobalLogConfigChange'.
// For example: config.RegisterSectionChangeCallback("log", handleGlobalLogConfigChange)

// TODO: Consider how to handle initial configuration loading. Does pkg/log's Init
//       need to be aware of this watcher, or does the watcher get triggered after initial load? 

// ClearAllCallbacks removes all registered callbacks. Used primarily for testing.
// ... existing code ... 