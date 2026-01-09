package message

import (
	"app-platform-backend/core/module"
	messageapi "app-platform-backend/internal/api/v1/message"
	"app-platform-backend/internal/pkg/database"

	"github.com/gin-gonic/gin"
)

func init() { module.Register(&MessageModule{}) }

type MessageModule struct{}

func (m *MessageModule) Meta() module.Meta {
	return module.Meta{Code: "message_center", Name: "消息中心", Description: "消息中心模块", Icon: "message", SortOrder: 2}
}

func (m *MessageModule) GetFunctions() []module.Function {
	return []module.Function{
		{Code: "message_send", Name: "发送消息", Type: "active", Description: "发送站内消息"},
		{Code: "message_list", Name: "消息列表", Type: "passive", Description: "获取消息列表"},
		{Code: "message_template", Name: "消息模板", Type: "passive", Description: "管理消息模板"},
		{Code: "message_unread", Name: "未读统计", Type: "passive", Description: "获取未读消息数"},
		{Code: "message_mark_read", Name: "标记已读", Type: "active", Description: "标记消息已读"},
		{Code: "message_batch_send", Name: "批量发送", Type: "active", Description: "批量发送消息"},
	}
}

func (m *MessageModule) RegisterRoutes(group *gin.RouterGroup) {
	messageapi.InitDB(database.GetDB())

	g := group.Group("/messages")
	{
		g.GET("", messageapi.List)
		g.POST("", messageapi.Send)
		g.GET("/templates", messageapi.Templates)
		g.GET("/unread", messageapi.UnreadCount)
		g.GET("/stats", messageapi.Stats)
		g.GET("/:id", messageapi.Detail)
		g.DELETE("/:id", messageapi.Delete)
		g.POST("/:id/read", messageapi.MarkRead)
		g.POST("/mark-all-read", messageapi.MarkAllRead)
		g.POST("/batch-delete", messageapi.BatchDelete)
		g.POST("/batch-send", messageapi.BatchSend)
	}
}

func (m *MessageModule) Init() error { return nil }
