package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var auditDB *gorm.DB

// AuditLog 审计日志模型（与audit包中的定义一致）
type AuditLog struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	AppID         uint      `json:"app_id" gorm:"index"`
	UserID        string    `json:"user_id" gorm:"index"`
	UserName      string    `json:"user_name"`
	Action        string    `json:"action" gorm:"index"`
	Resource      string    `json:"resource" gorm:"index"`
	ResourceID    string    `json:"resource_id"`
	Description   string    `json:"description"`
	IPAddress     string    `json:"ip_address"`
	UserAgent     string    `json:"user_agent"`
	RequestPath   string    `json:"request_path"`
	RequestMethod string    `json:"request_method"`
	StatusCode    int       `json:"status_code"`
	Duration      int64     `json:"duration"`
	Extra         string    `json:"extra" gorm:"type:text"`
	CreatedAt     time.Time `json:"created_at" gorm:"index"`
}

// InitAuditDB 初始化审计数据库连接
func InitAuditDB(db *gorm.DB) {
	auditDB = db
}

// AuditMiddleware 审计日志中间件
func AuditMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过不需要审计的路径
		path := c.Request.URL.Path
		if shouldSkipAudit(path) {
			c.Next()
			return
		}

		startTime := time.Now()

		// 读取请求体（用于记录）
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// 处理请求
		c.Next()

		// 计算耗时
		duration := time.Since(startTime).Milliseconds()

		// 异步记录审计日志
		go recordAuditLog(c, duration)
	}
}

// shouldSkipAudit 判断是否跳过审计
func shouldSkipAudit(path string) bool {
	skipPaths := []string{
		"/api/v1/health",
		"/api/v1/ws",
		"/api/v1/monitor/metrics",
		"/api/v1/logs/report",
		"/api/v1/events",
		"/api/v1/audit",
	}
	
	for _, skip := range skipPaths {
		if path == skip || (len(path) > len(skip) && path[:len(skip)] == skip) {
			return true
		}
	}
	return false
}

// recordAuditLog 记录审计日志
func recordAuditLog(c *gin.Context, duration int64) {
	if auditDB == nil {
		return
	}

	// 从上下文获取用户信息
	userID, _ := c.Get("user_id")
	userName, _ := c.Get("user_name")
	
	userIDStr := ""
	if userID != nil {
		userIDStr = userID.(string)
	}
	userNameStr := ""
	if userName != nil {
		userNameStr = userName.(string)
	}

	// 解析操作类型和资源
	action, resource := parseActionAndResource(c.Request.Method, c.Request.URL.Path)

	// 获取资源ID
	resourceID := c.Param("id")
	if resourceID == "" {
		resourceID = c.Query("id")
	}

	// 生成描述
	description := generateDescription(action, resource, c.Request.Method, c.Request.URL.Path)

	log := &AuditLog{
		UserID:        userIDStr,
		UserName:      userNameStr,
		Action:        action,
		Resource:      resource,
		ResourceID:    resourceID,
		Description:   description,
		IPAddress:     c.ClientIP(),
		UserAgent:     c.Request.UserAgent(),
		RequestPath:   c.Request.URL.Path,
		RequestMethod: c.Request.Method,
		StatusCode:    c.Writer.Status(),
		Duration:      duration,
		CreatedAt:     time.Now(),
	}

	auditDB.Create(log)
}

// parseActionAndResource 解析操作类型和资源
func parseActionAndResource(method, path string) (action, resource string) {
	// 根据HTTP方法判断操作类型
	switch method {
	case "GET":
		action = "view"
	case "POST":
		action = "create"
	case "PUT", "PATCH":
		action = "update"
	case "DELETE":
		action = "delete"
	default:
		action = "unknown"
	}

	// 解析资源类型
	resourceMap := map[string]string{
		"/api/v1/users":    "user",
		"/api/v1/apps":     "app",
		"/api/v1/configs":  "config",
		"/api/v1/messages": "message",
		"/api/v1/push":     "push",
		"/api/v1/files":    "file",
		"/api/v1/versions": "version",
		"/api/v1/logs":     "log",
		"/api/v1/events":   "event",
		"/api/v1/monitor":  "monitor",
		"/api/v1/admin":    "admin",
		"/api/v1/modules":  "module",
	}

	resource = "unknown"
	for prefix, res := range resourceMap {
		if len(path) >= len(prefix) && path[:len(prefix)] == prefix {
			resource = res
			break
		}
	}

	// 特殊操作识别
	if contains(path, "/login") {
		action = "login"
	} else if contains(path, "/logout") {
		action = "logout"
	} else if contains(path, "/export") {
		action = "export"
	} else if contains(path, "/publish") {
		action = "publish"
	} else if contains(path, "/send") {
		action = "send"
	}

	return action, resource
}

// generateDescription 生成操作描述
func generateDescription(action, resource, method, path string) string {
	actionNames := map[string]string{
		"view":    "查看",
		"create":  "创建",
		"update":  "更新",
		"delete":  "删除",
		"login":   "登录",
		"logout":  "登出",
		"export":  "导出",
		"publish": "发布",
		"send":    "发送",
	}

	resourceNames := map[string]string{
		"user":    "用户",
		"app":     "应用",
		"config":  "配置",
		"message": "消息",
		"push":    "推送",
		"file":    "文件",
		"version": "版本",
		"log":     "日志",
		"event":   "事件",
		"monitor": "监控",
		"admin":   "管理员",
		"module":  "模块",
	}

	actionName := actionNames[action]
	if actionName == "" {
		actionName = action
	}
	resourceName := resourceNames[resource]
	if resourceName == "" {
		resourceName = resource
	}

	return actionName + resourceName
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
