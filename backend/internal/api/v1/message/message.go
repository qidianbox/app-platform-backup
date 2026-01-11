package message

import (
	"app-platform-backend/internal/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(database *gorm.DB) {
	db = database
}

// List 消息列表
func List(c *gin.Context) {
	appID := c.Query("app_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "app_id is required"})
		return
	}

	query := db.Model(&model.Message{}).Where("app_id = ?", appID)

	var total int64
	query.Count(&total)

	var messages []model.Message
	offset := (page - 1) * size
	query.Offset(offset).Limit(size).Order("created_at DESC").Find(&messages)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  messages,
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}

// Send 发送消息
func Send(c *gin.Context) {
	var req struct {
		AppID   uint   `json:"app_id" binding:"required"`
		UserID  *uint  `json:"user_id"`
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
		Type    string `json:"type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	if req.Type == "" {
		req.Type = "system"
	}

	message := model.Message{
		AppID:   req.AppID,
		UserID:  req.UserID,
		Title:   req.Title,
		Content: req.Content,
		Type:    req.Type,
		Status:  0,
	}

	if err := db.Create(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to send message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    message,
		"message": "Message sent successfully",
	})
}

func Templates(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
}

func UnreadCount(c *gin.Context) {
	appID := c.Query("app_id")
	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "app_id is required"})
		return
	}

	var count int64
	db.Model(&model.Message{}).Where("app_id = ? AND status = 0", appID).Count(&count)

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"count": count}})
}

// Detail 消息详情
func Detail(c *gin.Context) {
	id := c.Param("id")
	appID := c.Query("app_id")

	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "app_id is required"})
		return
	}

	var message model.Message
	// 同时验证id和app_id，防止越权访问
	if err := db.Where("id = ? AND app_id = ?", id, appID).First(&message).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Message not found or no permission"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to query message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": message,
	})
}

// Stats 消息统计
func Stats(c *gin.Context) {
	appID := c.Query("app_id")
	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "app_id is required"})
		return
	}

	var total, unread, todayCount int64
	db.Model(&model.Message{}).Where("app_id = ?", appID).Count(&total)
	db.Model(&model.Message{}).Where("app_id = ? AND status = 0", appID).Count(&unread)

	today := time.Now().Format("2006-01-02")
	db.Model(&model.Message{}).Where("app_id = ? AND DATE(created_at) = ?", appID, today).Count(&todayCount)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"total":       total,
			"unread":      unread,
			"read":        total - unread,
			"today_count": todayCount,
		},
	})
}

// MarkRead 标记消息已读
func MarkRead(c *gin.Context) {
	id := c.Param("id")
	appID := c.Query("app_id")

	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "app_id is required"})
		return
	}

	var message model.Message
	// 同时验证id和app_id，防止越权操作
	if err := db.Where("id = ? AND app_id = ?", id, appID).First(&message).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Message not found or no permission"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to query message"})
		return
	}

	db.Model(&message).Update("status", 1)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Message marked as read",
	})
}

// MarkAllRead 标记所有消息已读
func MarkAllRead(c *gin.Context) {
	var req struct {
		AppID  uint  `json:"app_id" binding:"required"`
		UserID *uint `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	query := db.Model(&model.Message{}).Where("app_id = ? AND status = 0", req.AppID)
	if req.UserID != nil {
		query = query.Where("user_id = ?", *req.UserID)
	}

	result := query.Update("status", 1)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "All messages marked as read",
		"data": gin.H{
			"affected": result.RowsAffected,
		},
	})
}

// Delete 删除消息
func Delete(c *gin.Context) {
	id := c.Param("id")
	appID := c.Query("app_id")

	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "app_id is required"})
		return
	}

	var message model.Message
	// 同时验证id和app_id，防止越权删除
	if err := db.Where("id = ? AND app_id = ?", id, appID).First(&message).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Message not found or no permission"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to query message"})
		return
	}

	db.Delete(&message)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Message deleted successfully",
	})
}

// BatchDelete 批量删除消息
func BatchDelete(c *gin.Context) {
	var req struct {
		AppID uint   `json:"app_id" binding:"required"`
		IDs   []uint `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// 只删除属于该APP的消息，防止越权删除
	result := db.Where("id IN ? AND app_id = ?", req.IDs, req.AppID).Delete(&model.Message{})

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Messages deleted successfully",
		"data": gin.H{
			"affected": result.RowsAffected,
		},
	})
}

// BatchSend 批量发送消息
func BatchSend(c *gin.Context) {
	var req struct {
		AppID   uint    `json:"app_id" binding:"required"`
		UserIDs []uint  `json:"user_ids"`
		Title   string  `json:"title" binding:"required"`
		Content string  `json:"content" binding:"required"`
		Type    string  `json:"type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	if req.Type == "" {
		req.Type = "system"
	}

	var messages []model.Message
	if len(req.UserIDs) == 0 {
		// 发送给所有用户（广播）
		messages = append(messages, model.Message{
			AppID:   req.AppID,
			UserID:  nil,
			Title:   req.Title,
			Content: req.Content,
			Type:    req.Type,
			Status:  0,
		})
	} else {
		// 发送给指定用户
		for _, userID := range req.UserIDs {
			uid := userID
			messages = append(messages, model.Message{
				AppID:   req.AppID,
				UserID:  &uid,
				Title:   req.Title,
				Content: req.Content,
				Type:    req.Type,
				Status:  0,
			})
		}
	}

	if err := db.Create(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to send messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Messages sent successfully",
		"data": gin.H{
			"count": len(messages),
		},
	})
}
