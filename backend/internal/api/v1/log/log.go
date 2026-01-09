package log

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

// List 日志列表
func List(c *gin.Context) {
	appID := c.Query("app_id")
	level := c.Query("level")
	module := c.Query("module")
	keyword := c.Query("keyword")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "app_id is required"})
		return
	}

	query := db.Model(&model.Log{}).Where("app_id = ?", appID)

	if level != "" {
		query = query.Where("level = ?", level)
	}
	if module != "" {
		query = query.Where("module = ?", module)
	}
	if keyword != "" {
		query = query.Where("message LIKE ?", "%"+keyword+"%")
	}
	if startTime != "" {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("created_at <= ?", endTime)
	}

	var total int64
	query.Count(&total)

	var logs []model.Log
	offset := (page - 1) * size
	query.Offset(offset).Limit(size).Order("created_at DESC").Find(&logs)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  logs,
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}

// Report 上报日志
func Report(c *gin.Context) {
	var req struct {
		AppID   uint   `json:"app_id" binding:"required"`
		Level   string `json:"level"`
		Module  string `json:"module"`
		Message string `json:"message" binding:"required"`
		Context string `json:"context"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	if req.Level == "" {
		req.Level = "info"
	}

	log := model.Log{
		AppID:   req.AppID,
		Level:   req.Level,
		Module:  req.Module,
		Message: req.Message,
		Context: req.Context,
		IP:      c.ClientIP(),
	}

	if err := db.Create(&log).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to report log"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Log reported successfully",
	})
}

// BatchReport 批量上报日志
func BatchReport(c *gin.Context) {
	var req struct {
		AppID uint `json:"app_id" binding:"required"`
		Logs  []struct {
			Level   string `json:"level"`
			Module  string `json:"module"`
			Message string `json:"message"`
			Context string `json:"context"`
		} `json:"logs" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	var logs []model.Log
	clientIP := c.ClientIP()
	for _, l := range req.Logs {
		level := l.Level
		if level == "" {
			level = "info"
		}
		logs = append(logs, model.Log{
			AppID:   req.AppID,
			Level:   level,
			Module:  l.Module,
			Message: l.Message,
			Context: l.Context,
			IP:      clientIP,
		})
	}

	if err := db.Create(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to report logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Logs reported successfully",
		"data": gin.H{
			"count": len(logs),
		},
	})
}

// Stats 日志统计
func Stats(c *gin.Context) {
	appID := c.Query("app_id")
	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "app_id is required"})
		return
	}

	var total, errorCount, warnCount, infoCount, debugCount, todayCount int64
	db.Model(&model.Log{}).Where("app_id = ?", appID).Count(&total)
	db.Model(&model.Log{}).Where("app_id = ? AND level = ?", appID, "error").Count(&errorCount)
	db.Model(&model.Log{}).Where("app_id = ? AND level = ?", appID, "warn").Count(&warnCount)
	db.Model(&model.Log{}).Where("app_id = ? AND level = ?", appID, "info").Count(&infoCount)
	db.Model(&model.Log{}).Where("app_id = ? AND level = ?", appID, "debug").Count(&debugCount)

	today := time.Now().Format("2006-01-02")
	db.Model(&model.Log{}).Where("app_id = ? AND DATE(created_at) = ?", appID, today).Count(&todayCount)

	// 获取最近7天的日志趋势
	var trends []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	db.Model(&model.Log{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("app_id = ? AND created_at >= ?", appID, time.Now().AddDate(0, 0, -7)).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&trends)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"total":       total,
			"error_count": errorCount,
			"warn_count":  warnCount,
			"info_count":  infoCount,
			"debug_count": debugCount,
			"today_count": todayCount,
			"trends":      trends,
		},
	})
}

// Export 导出日志
func Export(c *gin.Context) {
	appID := c.Query("app_id")
	level := c.Query("level")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "app_id is required"})
		return
	}

	query := db.Model(&model.Log{}).Where("app_id = ?", appID)

	if level != "" {
		query = query.Where("level = ?", level)
	}
	if startTime != "" {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("created_at <= ?", endTime)
	}

	var logs []model.Log
	query.Order("created_at DESC").Limit(10000).Find(&logs)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"logs":  logs,
			"count": len(logs),
		},
	})
}

// Clean 清理日志
func Clean(c *gin.Context) {
	var req struct {
		AppID      uint   `json:"app_id" binding:"required"`
		BeforeDate string `json:"before_date"`
		Level      string `json:"level"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	query := db.Where("app_id = ?", req.AppID)

	if req.BeforeDate != "" {
		query = query.Where("created_at < ?", req.BeforeDate)
	} else {
		// 默认清理30天前的日志
		query = query.Where("created_at < ?", time.Now().AddDate(0, 0, -30))
	}

	if req.Level != "" {
		query = query.Where("level = ?", req.Level)
	}

	result := query.Delete(&model.Log{})

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Logs cleaned successfully",
		"data": gin.H{
			"affected": result.RowsAffected,
		},
	})
}

// System 系统日志（兼容旧接口）
func System(c *gin.Context) {
	List(c)
}

// Operation 操作日志（兼容旧接口）
func Operation(c *gin.Context) {
	List(c)
}
