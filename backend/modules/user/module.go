package user

import (
"app-platform-backend/core/module"
userapi "app-platform-backend/internal/api/v1/user"
"github.com/gin-gonic/gin"
)

func init() { module.Register(&UserModule{}) }

type UserModule struct{}

func (m *UserModule) Meta() module.Meta {
return module.Meta{
Code: "user_management",
Name: "用户管理",
Description: "用户管理模块",
Icon: "user",
SortOrder: 1,
}
}

func (m *UserModule) GetFunctions() []module.Function {
return []module.Function{
{Code: "user_list", Name: "用户列表", Type: "passive", Description: "获取APP用户列表"},
{Code: "user_detail", Name: "用户详情", Type: "passive", Description: "获取用户详细信息"},
{Code: "user_status", Name: "用户状态管理", Type: "active", Description: "启用/禁用用户"},
{Code: "user_stats", Name: "用户统计", Type: "passive", Description: "用户数据统计"},
}
}

func (m *UserModule) RegisterRoutes(group *gin.RouterGroup) {
group.GET("/users", userapi.List)
group.GET("/users/:id", userapi.Detail)
group.PUT("/users/:id/status", userapi.UpdateStatus)
group.GET("/users/stats", userapi.Stats)
}

func (m *UserModule) Init() error { return nil }
