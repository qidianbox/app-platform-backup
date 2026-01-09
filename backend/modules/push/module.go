package push

import (
	"app-platform-backend/core/module"
	pushapi "app-platform-backend/internal/api/v1/push"
	"app-platform-backend/internal/pkg/database"

	"github.com/gin-gonic/gin"
)

func init() { module.Register(&PushModule{}) }

type PushModule struct{}

func (m *PushModule) Meta() module.Meta {
	return module.Meta{Code: "push_service", Name: "推送服务", Description: "推送服务模块", Icon: "bell", SortOrder: 3}
}

func (m *PushModule) GetFunctions() []module.Function {
	return []module.Function{
		{Code: "push_create", Name: "创建推送", Type: "active", Description: "创建推送任务"},
		{Code: "push_send", Name: "发送推送", Type: "active", Description: "发送推送通知"},
		{Code: "push_list", Name: "推送列表", Type: "passive", Description: "推送任务列表"},
		{Code: "push_stats", Name: "推送统计", Type: "passive", Description: "推送数据统计"},
		{Code: "push_template", Name: "推送模板", Type: "passive", Description: "管理推送模板"},
		{Code: "push_cancel", Name: "取消推送", Type: "active", Description: "取消推送任务"},
	}
}

func (m *PushModule) RegisterRoutes(group *gin.RouterGroup) {
	pushapi.InitDB(database.GetDB())

	g := group.Group("/push")
	{
		g.GET("", pushapi.List)
		g.POST("", pushapi.Create)
		g.GET("/stats", pushapi.Stats)
		g.GET("/templates", pushapi.Templates)
		g.GET("/:id", pushapi.Detail)
		g.POST("/:id/send", pushapi.Send)
		g.POST("/:id/cancel", pushapi.Cancel)
		g.DELETE("/:id", pushapi.Delete)
		// 兼容旧接口
		g.GET("/tasks", pushapi.Tasks)
	}
}

func (m *PushModule) Init() error { return nil }
