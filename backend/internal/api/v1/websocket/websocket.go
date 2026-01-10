package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源，生产环境应该限制
	},
}

// Client 表示一个WebSocket客户端连接
type Client struct {
	ID       string
	AppID    uint
	UserID   string
	Conn     *ws.Conn
	Send     chan []byte
	Hub      *Hub
	mu       sync.Mutex
}

// Hub 管理所有WebSocket连接
type Hub struct {
	clients    map[*Client]bool
	appClients map[uint]map[*Client]bool // 按APP分组的客户端
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// Message WebSocket消息结构
type Message struct {
	Type      string      `json:"type"`      // 消息类型: monitor, alert, notification, log
	AppID     uint        `json:"app_id"`    // 目标APP ID，0表示广播
	UserID    string      `json:"user_id"`   // 目标用户ID，空表示广播
	Data      interface{} `json:"data"`      // 消息数据
	Timestamp int64       `json:"timestamp"` // 时间戳
}

// MonitorData 监控数据结构
type MonitorData struct {
	CPU        float64 `json:"cpu"`
	Memory     float64 `json:"memory"`
	Requests   int64   `json:"requests"`
	Errors     int64   `json:"errors"`
	ErrorRate  float64 `json:"error_rate"`
	AvgLatency float64 `json:"avg_latency"`
}

// AlertData 告警数据结构
type AlertData struct {
	ID        uint   `json:"id"`
	Level     string `json:"level"`     // critical, warning, info
	Title     string `json:"title"`
	Message   string `json:"message"`
	Source    string `json:"source"`
	Status    string `json:"status"`    // active, resolved
	CreatedAt int64  `json:"created_at"`
}

var hub *Hub

func init() {
	hub = NewHub()
	go hub.Run()
}

// NewHub 创建新的Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		appClients: make(map[uint]map[*Client]bool),
		broadcast:  make(chan *Message, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run 运行Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			if _, ok := h.appClients[client.AppID]; !ok {
				h.appClients[client.AppID] = make(map[*Client]bool)
			}
			h.appClients[client.AppID][client] = true
			h.mu.Unlock()
			log.Printf("[WebSocket] Client registered: %s (AppID: %d)", client.ID, client.AppID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				if appClients, ok := h.appClients[client.AppID]; ok {
					delete(appClients, client)
				}
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("[WebSocket] Client unregistered: %s", client.ID)

		case message := <-h.broadcast:
			h.mu.RLock()
			data, _ := json.Marshal(message)
			
			// 如果指定了AppID，只发送给该APP的客户端
			if message.AppID > 0 {
				if appClients, ok := h.appClients[message.AppID]; ok {
					for client := range appClients {
						select {
						case client.Send <- data:
						default:
							close(client.Send)
							delete(h.clients, client)
							delete(appClients, client)
						}
					}
				}
			} else {
				// 广播给所有客户端
				for client := range h.clients {
					select {
					case client.Send <- data:
					default:
						close(client.Send)
						delete(h.clients, client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// GetHub 获取全局Hub实例
func GetHub() *Hub {
	return hub
}

// Broadcast 广播消息
func (h *Hub) Broadcast(msg *Message) {
	msg.Timestamp = time.Now().UnixMilli()
	h.broadcast <- msg
}

// BroadcastToApp 向指定APP广播消息
func (h *Hub) BroadcastToApp(appID uint, msgType string, data interface{}) {
	h.Broadcast(&Message{
		Type:  msgType,
		AppID: appID,
		Data:  data,
	})
}

// BroadcastMonitorData 广播监控数据
func BroadcastMonitorData(appID uint, data *MonitorData) {
	hub.BroadcastToApp(appID, "monitor", data)
}

// BroadcastAlert 广播告警
func BroadcastAlert(appID uint, alert *AlertData) {
	hub.BroadcastToApp(appID, "alert", alert)
}

// BroadcastNotification 广播通知
func BroadcastNotification(appID uint, title, message string) {
	hub.BroadcastToApp(appID, "notification", map[string]string{
		"title":   title,
		"message": message,
	})
}

// HandleWebSocket WebSocket连接处理器
func HandleWebSocket(c *gin.Context) {
	appIDStr := c.Query("app_id")
	userID := c.Query("user_id")
	
	var appID uint
	if appIDStr != "" {
		var id uint64
		_, err := fmt.Sscanf(appIDStr, "%d", &id)
		if err == nil {
			appID = uint(id)
		}
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[WebSocket] Upgrade error: %v", err)
		return
	}

	client := &Client{
		ID:     generateClientID(),
		AppID:  appID,
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		Hub:    hub,
	}

	hub.register <- client

	// 启动读写协程
	go client.writePump()
	go client.readPump()
}

// readPump 读取客户端消息
func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512 * 1024) // 512KB
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if ws.IsUnexpectedCloseError(err, ws.CloseGoingAway, ws.CloseAbnormalClosure) {
				log.Printf("[WebSocket] Read error: %v", err)
			}
			break
		}

		// 处理客户端消息（如心跳、订阅等）
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err == nil {
			if msgType, ok := msg["type"].(string); ok {
				switch msgType {
				case "ping":
					c.Send <- []byte(`{"type":"pong"}`)
				case "subscribe":
					// 处理订阅请求
					log.Printf("[WebSocket] Client %s subscribed", c.ID)
				}
			}
		}
	}
}

// writePump 向客户端发送消息
func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(ws.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(ws.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 批量发送队列中的消息
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(ws.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// generateClientID 生成客户端ID
func generateClientID() string {
	return fmt.Sprintf("client_%d", time.Now().UnixNano())
}
