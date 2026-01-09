package config

import (
"app-platform-backend/core/module"
configapi "app-platform-backend/internal/api/v1/config"
"github.com/gin-gonic/gin"
)

func init() { module.Register(&ConfigModule{}) }
type ConfigModule struct{}
func (m *ConfigModule) Meta() module.Meta {
return module.Meta{Code: "config_management", Name: "配置管理", Description: "配置管理模块", Icon: "settings", SortOrder: 8}
}
func (m *ConfigModule) GetFunctions() []module.Function {
return []module.Function{
{Code: "config_list", Name: "配置列表", Type: "passive", Description: "获取配置列表"},
{Code: "config_create", Name: "创建配置", Type: "active", Description: "创建新配置"},
{Code: "config_update", Name: "更新配置", Type: "active", Description: "更新配置"},
{Code: "config_publish", Name: "发布配置", Type: "active", Description: "发布配置"},
{Code: "config_history", Name: "配置历史", Type: "passive", Description: "查看配置历史"},
}
}
func (m *ConfigModule) RegisterRoutes(group *gin.RouterGroup) {
group.GET("/configs", configapi.List)
group.POST("/configs", configapi.Create)
group.PUT("/configs/:id", configapi.Update)
group.POST("/configs/:id/publish", configapi.Publish)
group.GET("/configs/:id/history", configapi.History)
}
func (m *ConfigModule) Init() error { return nil }
