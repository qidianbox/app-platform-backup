// Package module 提供模块注册中心功能
// 所有模块在 init() 函数中调用 Register() 来注册自己
package module

import (
	"fmt"
	"sync"
)

var (
	// modules 存储所有已注册的模块
	modules = make(map[string]Module)
	// lock 保证并发安全
	lock = sync.RWMutex{}
	// initOrder 记录模块注册顺序，用于按顺序初始化
	initOrder []string
)

// Register 用于注册一个模块实例
// 通常在模块包的 init() 函数中调用
func Register(m Module) {
	lock.Lock()
	defer lock.Unlock()

	meta := m.Meta()
	if _, exists := modules[meta.Code]; exists {
		// 在应用启动时处理重复注册的错误
		panic(fmt.Sprintf("module already registered: %s", meta.Code))
	}
	modules[meta.Code] = m
	initOrder = append(initOrder, meta.Code)
}

// Get 根据模块Code获取模块实例
func Get(code string) (Module, bool) {
	lock.RLock()
	defer lock.RUnlock()

	m, exists := modules[code]
	return m, exists
}

// GetAllModules 返回所有已注册的模块
func GetAllModules() []Module {
	lock.RLock()
	defer lock.RUnlock()

	all := make([]Module, 0, len(modules))
	// 按注册顺序返回
	for _, code := range initOrder {
		if m, exists := modules[code]; exists {
			all = append(all, m)
		}
	}
	return all
}

// GetModuleCount 返回已注册模块的数量
func GetModuleCount() int {
	lock.RLock()
	defer lock.RUnlock()
	return len(modules)
}

// GetAllFunctions 返回所有模块提供的所有功能
func GetAllFunctions() []Function {
	lock.RLock()
	defer lock.RUnlock()

	var allFunctions []Function
	for _, code := range initOrder {
		if m, exists := modules[code]; exists {
			allFunctions = append(allFunctions, m.GetFunctions()...)
		}
	}
	return allFunctions
}

// InitAllModules 初始化所有已注册的模块
// 返回第一个遇到的错误
func InitAllModules() error {
	lock.RLock()
	defer lock.RUnlock()

	for _, code := range initOrder {
		if m, exists := modules[code]; exists {
			if err := m.Init(); err != nil {
				return fmt.Errorf("failed to init module %s: %w", code, err)
			}
		}
	}
	return nil
}

// Clear 清空所有已注册的模块（主要用于测试）
func Clear() {
	lock.Lock()
	defer lock.Unlock()
	modules = make(map[string]Module)
	initOrder = nil
}
