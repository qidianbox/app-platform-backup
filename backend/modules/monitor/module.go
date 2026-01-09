package monitor

import (
"app-platform-backend/core/module"
monitorapi "app-platform-backend/internal/api/v1/monitor"
"github.com/gin-gonic/gin"
)

func init() { module.Register(&MonitorModule{}) }
type MonitorModule struct{}
func (m *MonitorModule) Meta() module.Meta {
return module.Meta{Code: "monitor_service", Name: "监控服务", Description: "监控服务模块", Icon: "monitor", SortOrder: 6}
}
func (m *MonitorModule) GetFunctions() []module.Function {
return []module.Function{
{Code: "monitor_metrics", Name: "监控指标", Type: "passive", Description: "查看监控指标"},
{Code: "monitor_alerts", Name: "告警管理", Type: "passive", Description: "管理告警"},
{Code: "monitor_rules", Name: "告警规则", Type: "passive", Description: "管理告警规则"},
{Code: "monitor_health", Name: "健康检查", Type: "passive", Description: "系统健康检查"},
}
}
func (m *MonitorModule) RegisterRoutes(group *gin.RouterGroup) {
group.GET("/monitor/metrics", monitorapi.Metrics)
group.GET("/monitor/alerts", monitorapi.Alerts)
group.GET("/monitor/rules", monitorapi.Rules)
group.GET("/monitor/health", monitorapi.Health)
}
func (m *MonitorModule) Init() error { return nil }
