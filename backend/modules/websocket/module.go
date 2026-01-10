package websocket

import (
	"app-platform-backend/core/module"
	wsapi "app-platform-backend/internal/api/v1/websocket"

	"github.com/gin-gonic/gin"
)

func init() { module.Register(&WebSocketModule{}) }

type WebSocketModule struct{}

func (m *WebSocketModule) Meta() module.Meta {
	return module.Meta{Code: "websocket", Name: "WebSocket服务", Description: "实时推送服务模块", Icon: "broadcast", SortOrder: 10}
}

func (m *WebSocketModule) GetFunctions() []module.Function {
	return []module.Function{
		{Code: "ws_connect", Name: "WebSocket连接", Type: "active", Description: "建立WebSocket连接"},
		{Code: "ws_monitor", Name: "监控数据推送", Type: "passive", Description: "实时推送监控数据"},
		{Code: "ws_alert", Name: "告警推送", Type: "passive", Description: "实时推送告警通知"},
	}
}

func (m *WebSocketModule) RegisterRoutes(group *gin.RouterGroup) {
	// WebSocket连接端点
	group.GET("/ws", wsapi.HandleWebSocket)
}

func (m *WebSocketModule) Init() error { return nil }
