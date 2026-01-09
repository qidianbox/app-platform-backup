package log

import (
"app-platform-backend/core/module"
logapi "app-platform-backend/internal/api/v1/log"
"github.com/gin-gonic/gin"
)

func init() { module.Register(&LogModule{}) }
type LogModule struct{}
func (m *LogModule) Meta() module.Meta {
return module.Meta{Code: "log_service", Name: "日志服务", Description: "日志服务模块", Icon: "file-text", SortOrder: 5}
}
func (m *LogModule) GetFunctions() []module.Function {
return []module.Function{
{Code: "log_system", Name: "系统日志", Type: "passive", Description: "查看系统日志"},
{Code: "log_operation", Name: "操作日志", Type: "passive", Description: "查看操作日志"},
{Code: "log_stats", Name: "日志统计", Type: "passive", Description: "日志数据统计"},
{Code: "log_clean", Name: "日志清理", Type: "active", Description: "清理历史日志"},
}
}
func (m *LogModule) RegisterRoutes(group *gin.RouterGroup) {
group.GET("/logs/system", logapi.System)
group.GET("/logs/operation", logapi.Operation)
group.GET("/logs/stats", logapi.Stats)
group.POST("/logs/clean", logapi.Clean)
}
func (m *LogModule) Init() error { return nil }
