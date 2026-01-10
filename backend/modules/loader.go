// Package modules 提供模块加载功能
// 通过导入此包，可以自动加载所有已注册的功能模块
package modules

import (
	// 导入所有功能模块，触发其 init() 函数执行注册
	_ "app-platform-backend/modules/config"
	_ "app-platform-backend/modules/event"
	_ "app-platform-backend/modules/file"
	_ "app-platform-backend/modules/log"
	_ "app-platform-backend/modules/message"
	_ "app-platform-backend/modules/monitor"
	_ "app-platform-backend/modules/push"
	_ "app-platform-backend/modules/user"
	_ "app-platform-backend/modules/version"
	_ "app-platform-backend/modules/websocket"
	_ "app-platform-backend/modules/audit"
)

// LoadAllModules 是一个空函数，其唯一目的是确保此包被导入
// 当此包被导入时，上面的 import 语句会自动触发各模块的 init() 函数
// 从而完成所有模块的自动注册
func LoadAllModules() {
	// 此函数体为空，模块加载通过 import 的副作用完成
}
