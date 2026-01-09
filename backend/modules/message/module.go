package message

import (
"app-platform-backend/core/module"
messageapi "app-platform-backend/internal/api/v1/message"
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
}
}

func (m *MessageModule) RegisterRoutes(group *gin.RouterGroup) {
group.POST("/messages", messageapi.Send)
group.GET("/messages", messageapi.List)
group.GET("/messages/templates", messageapi.Templates)
group.GET("/messages/unread", messageapi.UnreadCount)
}

func (m *MessageModule) Init() error { return nil }
