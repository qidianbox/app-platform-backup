package audit

import (
	"app-platform-backend/core/module"
	auditapi "app-platform-backend/internal/api/v1/audit"
	"app-platform-backend/internal/pkg/database"

	"github.com/gin-gonic/gin"
)

func init() { module.Register(&AuditModule{}) }

type AuditModule struct{}

func (m *AuditModule) Meta() module.Meta {
	return module.Meta{Code: "audit_log", Name: "审计日志", Description: "操作审计日志模块", Icon: "shield", SortOrder: 11}
}

func (m *AuditModule) GetFunctions() []module.Function {
	return []module.Function{
		{Code: "audit_list", Name: "审计日志列表", Type: "passive", Description: "查看审计日志列表"},
		{Code: "audit_stats", Name: "审计统计", Type: "passive", Description: "审计日志统计分析"},
		{Code: "audit_export", Name: "导出审计日志", Type: "active", Description: "导出审计日志"},
	}
}

func (m *AuditModule) RegisterRoutes(group *gin.RouterGroup) {
	auditapi.InitDB(database.GetDB())

	g := group.Group("/audit")
	{
		g.GET("", auditapi.List)
		g.GET("/logs", auditapi.List)  // 别名路由，兼容前端请求
		g.GET("/stats", auditapi.Stats)
		g.GET("/export", auditapi.Export)
		g.POST("/cleanup", auditapi.Cleanup)        // 手动清理
		g.GET("/cleanup/history", auditapi.CleanupHistory) // 清理历史
		g.GET("/cleanup/config", auditapi.CleanupConfig)   // 清理配置
	}
}

func (m *AuditModule) Init() error { return nil }
