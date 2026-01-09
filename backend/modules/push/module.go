package push

import (
"app-platform-backend/core/module"
pushapi "app-platform-backend/internal/api/v1/push"
"github.com/gin-gonic/gin"
)

func init() { module.Register(&PushModule{}) }
type PushModule struct{}
func (m *PushModule) Meta() module.Meta {
return module.Meta{Code: "push_service", Name: "推送服务", Description: "推送服务模块", Icon: "bell", SortOrder: 3}
}
func (m *PushModule) GetFunctions() []module.Function {
return []module.Function{
{Code: "push_send", Name: "发送推送", Type: "active", Description: "发送推送通知"},
{Code: "push_tasks", Name: "推送任务", Type: "passive", Description: "推送任务列表"},
{Code: "push_stats", Name: "推送统计", Type: "passive", Description: "推送数据统计"},
{Code: "push_template", Name: "推送模板", Type: "passive", Description: "管理推送模板"},
}
}
func (m *PushModule) RegisterRoutes(group *gin.RouterGroup) {
group.POST("/push", pushapi.Send)
group.GET("/push/tasks", pushapi.Tasks)
group.GET("/push/stats", pushapi.Stats)
group.GET("/push/templates", pushapi.Templates)
}
func (m *PushModule) Init() error { return nil }
