package event

import (
"app-platform-backend/core/module"
eventapi "app-platform-backend/internal/api/v1/event"
"github.com/gin-gonic/gin"
)

func init() { module.Register(&EventModule{}) }
type EventModule struct{}
func (m *EventModule) Meta() module.Meta {
return module.Meta{Code: "event_tracking", Name: "埋点服务", Description: "埋点服务模块", Icon: "chart", SortOrder: 4}
}
func (m *EventModule) GetFunctions() []module.Function {
return []module.Function{
{Code: "event_report", Name: "事件上报", Type: "active", Description: "上报埋点事件"},
{Code: "event_list", Name: "事件列表", Type: "passive", Description: "获取事件列表"},
{Code: "event_stats", Name: "事件统计", Type: "passive", Description: "事件数据统计"},
{Code: "event_funnel", Name: "漏斗分析", Type: "passive", Description: "漏斗分析"},
{Code: "event_definition", Name: "事件定义", Type: "passive", Description: "管理事件定义"},
}
}
func (m *EventModule) RegisterRoutes(group *gin.RouterGroup) {
group.POST("/events", eventapi.Report)
group.GET("/events", eventapi.List)
group.GET("/events/stats", eventapi.Stats)
group.GET("/events/funnel", eventapi.Funnel)
group.GET("/events/definitions", eventapi.Definitions)
}
func (m *EventModule) Init() error { return nil }
