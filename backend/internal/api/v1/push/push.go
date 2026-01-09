package push

import (
	"app-platform-backend/internal/model"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(database *gorm.DB) {
	db = database
}

// List 推送列表
func List(c *gin.Context) {
	appID := c.Query("app_id")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "app_id is required"})
		return
	}

	query := db.Model(&model.PushRecord{}).Where("app_id = ?", appID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var records []model.PushRecord
	offset := (page - 1) * size
	query.Offset(offset).Limit(size).Order("created_at DESC").Find(&records)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  records,
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}

// Create 创建推送任务
func Create(c *gin.Context) {
	var req struct {
		AppID       uint     `json:"app_id" binding:"required"`
		Title       string   `json:"title" binding:"required"`
		Content     string   `json:"content" binding:"required"`
		TargetType  string   `json:"target_type"`
		TargetIDs   []string `json:"target_ids"`
		ScheduledAt string   `json:"scheduled_at"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	if req.TargetType == "" {
		req.TargetType = "all"
	}

	record := model.PushRecord{
		AppID:      req.AppID,
		Title:      req.Title,
		Content:    req.Content,
		TargetType: req.TargetType,
		TargetIDs:  strings.Join(req.TargetIDs, ","),
		Status:     "pending",
	}

	if req.ScheduledAt != "" {
		scheduledTime, err := time.Parse("2006-01-02 15:04:05", req.ScheduledAt)
		if err == nil {
			record.ScheduledAt = &scheduledTime
		}
	}

	if err := db.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create push task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    record,
		"message": "Push task created successfully",
	})
}

// Detail 推送详情
func Detail(c *gin.Context) {
	id := c.Param("id")

	var record model.PushRecord
	if err := db.First(&record, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Push record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to query push record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": record,
	})
}

// Send 立即发送推送
func Send(c *gin.Context) {
	id := c.Param("id")

	var record model.PushRecord
	if err := db.First(&record, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Push record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to query push record"})
		return
	}

	// 模拟发送推送
	now := time.Now()
	sentCount := 100
	successCount := 95
	failedCount := 5

	db.Model(&record).Updates(map[string]interface{}{
		"status":        "sent",
		"sent_at":       now,
		"sent_count":    sentCount,
		"success_count": successCount,
		"failed_count":  failedCount,
	})

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Push sent successfully",
		"data": gin.H{
			"sent_count":    sentCount,
			"success_count": successCount,
			"failed_count":  failedCount,
		},
	})
}

// Cancel 取消推送任务
func Cancel(c *gin.Context) {
	id := c.Param("id")

	var record model.PushRecord
	if err := db.First(&record, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Push record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to query push record"})
		return
	}

	if record.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Only pending push can be cancelled"})
		return
	}

	db.Model(&record).Update("status", "cancelled")

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Push cancelled successfully",
	})
}

// Delete 删除推送记录
func Delete(c *gin.Context) {
	id := c.Param("id")

	var record model.PushRecord
	if err := db.First(&record, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Push record not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to query push record"})
		return
	}

	db.Delete(&record)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Push record deleted successfully",
	})
}

// Stats 推送统计
func Stats(c *gin.Context) {
	appID := c.Query("app_id")
	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "app_id is required"})
		return
	}

	var total, pending, sent, cancelled int64
	var totalSent, totalSuccess, totalFailed int64

	db.Model(&model.PushRecord{}).Where("app_id = ?", appID).Count(&total)
	db.Model(&model.PushRecord{}).Where("app_id = ? AND status = ?", appID, "pending").Count(&pending)
	db.Model(&model.PushRecord{}).Where("app_id = ? AND status = ?", appID, "sent").Count(&sent)
	db.Model(&model.PushRecord{}).Where("app_id = ? AND status = ?", appID, "cancelled").Count(&cancelled)

	db.Model(&model.PushRecord{}).Where("app_id = ?", appID).
		Select("COALESCE(SUM(sent_count), 0)").Scan(&totalSent)
	db.Model(&model.PushRecord{}).Where("app_id = ?", appID).
		Select("COALESCE(SUM(success_count), 0)").Scan(&totalSuccess)
	db.Model(&model.PushRecord{}).Where("app_id = ?", appID).
		Select("COALESCE(SUM(failed_count), 0)").Scan(&totalFailed)

	successRate := float64(0)
	if totalSent > 0 {
		successRate = float64(totalSuccess) / float64(totalSent) * 100
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"total":         total,
			"pending":       pending,
			"sent":          sent,
			"cancelled":     cancelled,
			"total_sent":    totalSent,
			"total_success": totalSuccess,
			"total_failed":  totalFailed,
			"success_rate":  successRate,
		},
	})
}

// Tasks 推送任务列表（兼容旧接口）
func Tasks(c *gin.Context) {
	List(c)
}

// Templates 推送模板（兼容旧接口）
func Templates(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": []gin.H{
			{"id": 1, "name": "系统通知", "title_template": "【系统通知】{{title}}", "content_template": "{{content}}"},
			{"id": 2, "name": "活动推送", "title_template": "【活动】{{title}}", "content_template": "{{content}}，点击查看详情"},
			{"id": 3, "name": "订单通知", "title_template": "订单{{order_id}}状态更新", "content_template": "您的订单{{order_id}}{{status}}"},
		},
	})
}
