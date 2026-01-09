package log

import (
	"app-platform-backend/core/module"
	logapi "app-platform-backend/internal/api/v1/log"
	"app-platform-backend/internal/pkg/database"

	"github.com/gin-gonic/gin"
)

func init() { module.Register(&LogModule{}) }

type LogModule struct{}

func (m *LogModule) Meta() module.Meta {
	return module.Meta{Code: "log_service", Name: "日志服务", Description: "日志服务模块", Icon: "file-text", SortOrder: 5}
}

func (m *LogModule) GetFunctions() []module.Function {
	return []module.Function{
		{Code: "log_list", Name: "日志列表", Type: "passive", Description: "查看日志列表"},
		{Code: "log_report", Name: "上报日志", Type: "active", Description: "上报日志数据"},
		{Code: "log_stats", Name: "日志统计", Type: "passive", Description: "日志数据统计"},
		{Code: "log_export", Name: "导出日志", Type: "active", Description: "导出日志数据"},
		{Code: "log_clean", Name: "日志清理", Type: "active", Description: "清理历史日志"},
	}
}

func (m *LogModule) RegisterRoutes(group *gin.RouterGroup) {
	logapi.InitDB(database.GetDB())

	g := group.Group("/logs")
	{
		g.GET("", logapi.List)
		g.POST("/report", logapi.Report)
		g.POST("/batch-report", logapi.BatchReport)
		g.GET("/stats", logapi.Stats)
		g.GET("/export", logapi.Export)
		g.POST("/clean", logapi.Clean)
		// 兼容旧接口
		g.GET("/system", logapi.System)
		g.GET("/operation", logapi.Operation)
	}
}

func (m *LogModule) Init() error { return nil }
