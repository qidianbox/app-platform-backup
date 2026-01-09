// Package module 定义了模块化架构的核心接口和类型
// 所有功能模块都必须实现这些接口才能被主程序识别和加载
package module

import "github.com/gin-gonic/gin"

// Meta 包含了模块的基本元数据
type Meta struct {
	Code        string // 模块唯一标识，例如 "user_management"
	Name        string // 人类可读的名称，例如 "用户管理"
	Description string // 模块功能描述
	Icon        string // 模块图标
	SortOrder   int    // 排序顺序
}

// Function 定义了一个具体的功能点，对应数据库中的一条记录
type Function struct {
	Code         string                 // 功能的唯一标识，例如 "get_user_profile"
	Name         string                 // 功能名称，例如 "获取用户详情"
	Description  string                 // 功能描述
	Type         string                 // 功能类型: "active" (工作台可见) 或 "passive" (后台运行)
	ConfigSchema map[string]interface{} // 功能的JSON Schema配置
	Dependencies []string               // 依赖的其他功能Code列表
	SortOrder    int                    // 排序顺序
}

// Module 是所有功能模块必须实现的接口
type Module interface {
	// Meta 返回模块的元数据
	Meta() Meta

	// RegisterRoutes 向主应用的Gin引擎注册HTTP路由
	// router 是一个已经带有 /api/v1 前缀的路由组
	RegisterRoutes(router *gin.RouterGroup)

	// GetFunctions 返回该模块提供的所有"功能"列表
	// 这些"功能"将与数据库中的 module_templates 表对应
	GetFunctions() []Function

	// Init 模块初始化方法，在应用启动时调用
	// 可以用于初始化模块内部状态、数据库连接等
	Init() error
}

// BaseModule 提供了 Module 接口的基础实现
// 模块可以嵌入此结构体来获得默认实现
type BaseModule struct {
	meta      Meta
	functions []Function
}

// NewBaseModule 创建一个基础模块实例
func NewBaseModule(meta Meta, functions []Function) *BaseModule {
	return &BaseModule{
		meta:      meta,
		functions: functions,
	}
}

// Meta 返回模块的元数据
func (m *BaseModule) Meta() Meta {
	return m.meta
}

// GetFunctions 返回该模块提供的所有功能
func (m *BaseModule) GetFunctions() []Function {
	return m.functions
}

// RegisterRoutes 默认空实现，子模块可以覆盖
func (m *BaseModule) RegisterRoutes(router *gin.RouterGroup) {
	// 默认不注册任何路由
}

// Init 默认空实现，子模块可以覆盖
func (m *BaseModule) Init() error {
	return nil
}
