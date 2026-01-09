package monitor

import (
	"app-platform-backend/core/module"
	monitorapi "app-platform-backend/internal/api/v1/monitor"
	"app-platform-backend/internal/pkg/database"

	"github.com/gin-gonic/gin"
)

func init() { module.Register(&MonitorModule{}) }

type MonitorModule struct{}

func (m *MonitorModule) Meta() module.Meta {
	return module.Meta{Code: "monitor_service", Name: "监控服务", Description: "监控服务模块", Icon: "monitor", SortOrder: 6}
}

func (m *MonitorModule) GetFunctions() []module.Function {
	return []module.Function{
		{Code: "monitor_report", Name: "上报指标", Type: "active", Description: "上报监控指标"},
		{Code: "monitor_metrics", Name: "监控指标", Type: "passive", Description: "查看监控指标"},
		{Code: "monitor_alerts", Name: "告警管理", Type: "passive", Description: "管理告警"},
		{Code: "monitor_stats", Name: "监控统计", Type: "passive", Description: "监控数据统计"},
		{Code: "monitor_health", Name: "健康检查", Type: "passive", Description: "系统健康检查"},
	}
}

func (m *MonitorModule) RegisterRoutes(group *gin.RouterGroup) {
	monitorapi.InitDB(database.GetDB())

	g := group.Group("/monitor")
	{
		g.GET("/metrics", monitorapi.Metrics)
		g.POST("/metrics", monitorapi.ReportMetric)
		g.GET("/metrics/stats", monitorapi.MetricStats)
		g.GET("/stats", monitorapi.Stats)
		g.GET("/health", monitorapi.Health)
		// 告警管理
		g.GET("/alerts", monitorapi.Alerts)
		g.POST("/alerts", monitorapi.CreateAlert)
		g.PUT("/alerts/:id", monitorapi.UpdateAlert)
		g.DELETE("/alerts/:id", monitorapi.DeleteAlert)
		g.POST("/alerts/:id/resolve", monitorapi.ResolveAlert)
		// 兼容旧接口
		g.GET("/rules", monitorapi.Rules)
	}
}

func (m *MonitorModule) Init() error { return nil }
