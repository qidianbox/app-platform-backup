package event

import (
	"app-platform-backend/internal/model"
	"encoding/json"
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

// Report 上报事件
func Report(c *gin.Context) {
	var req struct {
		AppID      uint                   `json:"app_id" binding:"required"`
		UserID     *uint                  `json:"user_id"`
		EventCode  string                 `json:"event_code" binding:"required"`
		EventName  string                 `json:"event_name"`
		Properties map[string]interface{} `json:"properties"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	propertiesJSON := "{}"
	if req.Properties != nil {
		if data, err := json.Marshal(req.Properties); err == nil {
			propertiesJSON = string(data)
		}
	}

	event := model.Event{
		AppID:      req.AppID,
		UserID:     req.UserID,
		EventCode:  req.EventCode,
		EventName:  req.EventName,
		Properties: propertiesJSON,
		IP:         c.ClientIP(),
		UserAgent:  c.GetHeader("User-Agent"),
	}

	if err := db.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to report event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Event reported successfully",
	})
}

// BatchReport 批量上报事件
func BatchReport(c *gin.Context) {
	var req struct {
		AppID  uint `json:"app_id" binding:"required"`
		Events []struct {
			UserID     *uint                  `json:"user_id"`
			EventCode  string                 `json:"event_code"`
			EventName  string                 `json:"event_name"`
			Properties map[string]interface{} `json:"properties"`
		} `json:"events" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	var events []model.Event
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	for _, e := range req.Events {
		propertiesJSON := "{}"
		if e.Properties != nil {
			if data, err := json.Marshal(e.Properties); err == nil {
				propertiesJSON = string(data)
			}
		}

		events = append(events, model.Event{
			AppID:      req.AppID,
			UserID:     e.UserID,
			EventCode:  e.EventCode,
			EventName:  e.EventName,
			Properties: propertiesJSON,
			IP:         clientIP,
			UserAgent:  userAgent,
		})
	}

	if err := db.Create(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to report events"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Events reported successfully",
		"data": gin.H{
			"count": len(events),
		},
	})
}

// List 事件列表
func List(c *gin.Context) {
	appID := c.Query("app_id")
	eventCode := c.Query("event_code")
	userID := c.Query("user_id")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "app_id is required"})
		return
	}

	query := db.Model(&model.Event{}).Where("app_id = ?", appID)

	if eventCode != "" {
		query = query.Where("event_code = ?", eventCode)
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if startTime != "" {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("created_at <= ?", endTime)
	}

	var total int64
	query.Count(&total)

	var events []model.Event
	offset := (page - 1) * size
	query.Offset(offset).Limit(size).Order("created_at DESC").Find(&events)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  events,
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}

// Stats 事件统计
func Stats(c *gin.Context) {
	appID := c.Query("app_id")
	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "app_id is required"})
		return
	}

	var total, todayCount, uniqueUsers int64
	db.Model(&model.Event{}).Where("app_id = ?", appID).Count(&total)

	today := time.Now().Format("2006-01-02")
	db.Model(&model.Event{}).Where("app_id = ? AND DATE(created_at) = ?", appID, today).Count(&todayCount)
	db.Model(&model.Event{}).Where("app_id = ?", appID).Distinct("user_id").Count(&uniqueUsers)

	// 获取事件类型统计
	var eventStats []struct {
		EventCode string `json:"event_code"`
		Count     int64  `json:"count"`
	}
	db.Model(&model.Event{}).
		Select("event_code, COUNT(*) as count").
		Where("app_id = ?", appID).
		Group("event_code").
		Order("count DESC").
		Limit(10).
		Scan(&eventStats)

	// 获取最近7天的事件趋势
	var trends []struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	db.Model(&model.Event{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("app_id = ? AND created_at >= ?", appID, time.Now().AddDate(0, 0, -7)).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&trends)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"total":        total,
			"today_count":  todayCount,
			"unique_users": uniqueUsers,
			"event_stats":  eventStats,
			"trends":       trends,
		},
	})
}

// Funnel 漏斗分析
func Funnel(c *gin.Context) {
	appID := c.Query("app_id")
	steps := c.QueryArray("steps")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "app_id is required"})
		return
	}

	if len(steps) == 0 {
		// 返回默认漏斗示例
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": gin.H{
				"steps": []gin.H{
					{"name": "页面浏览", "count": 1000, "rate": 100},
					{"name": "点击按钮", "count": 500, "rate": 50},
					{"name": "提交表单", "count": 200, "rate": 20},
					{"name": "完成转化", "count": 100, "rate": 10},
				},
			},
		})
		return
	}

	var funnelData []gin.H
	var prevCount int64 = 0

	for i, step := range steps {
		query := db.Model(&model.Event{}).Where("app_id = ? AND event_code = ?", appID, step)
		if startTime != "" {
			query = query.Where("created_at >= ?", startTime)
		}
		if endTime != "" {
			query = query.Where("created_at <= ?", endTime)
		}

		var count int64
		query.Count(&count)

		rate := float64(100)
		if i > 0 && prevCount > 0 {
			rate = float64(count) / float64(prevCount) * 100
		}

		funnelData = append(funnelData, gin.H{
			"name":  step,
			"count": count,
			"rate":  rate,
		})

		prevCount = count
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"steps": funnelData,
		},
	})
}

// Definitions 事件定义列表
func Definitions(c *gin.Context) {
	appID := c.Query("app_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	query := db.Model(&model.EventDefinition{})
	if appID != "" {
		query = query.Where("app_id = ?", appID)
	}

	var total int64
	query.Count(&total)

	var definitions []model.EventDefinition
	offset := (page - 1) * size
	query.Offset(offset).Limit(size).Order("created_at DESC").Find(&definitions)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  definitions,
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}

// CreateDefinition 创建事件定义
func CreateDefinition(c *gin.Context) {
	var req struct {
		AppID            uint   `json:"app_id" binding:"required"`
		EventCode        string `json:"event_code" binding:"required"`
		EventName        string `json:"event_name" binding:"required"`
		Description      string `json:"description"`
		PropertiesSchema string `json:"properties_schema"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	definition := model.EventDefinition{
		AppID:            req.AppID,
		EventCode:        req.EventCode,
		EventName:        req.EventName,
		Description:      req.Description,
		PropertiesSchema: req.PropertiesSchema,
		IsActive:         1,
	}

	if err := db.Create(&definition).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to create event definition"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    definition,
		"message": "Event definition created successfully",
	})
}

// UpdateDefinition 更新事件定义
func UpdateDefinition(c *gin.Context) {
	id := c.Param("id")

	var definition model.EventDefinition
	if err := db.First(&definition, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Event definition not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to query event definition"})
		return
	}

	var req struct {
		EventName        string `json:"event_name"`
		Description      string `json:"description"`
		PropertiesSchema string `json:"properties_schema"`
		IsActive         *int   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.EventName != "" {
		updates["event_name"] = req.EventName
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.PropertiesSchema != "" {
		updates["properties_schema"] = req.PropertiesSchema
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	db.Model(&definition).Updates(updates)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Event definition updated successfully",
	})
}

// DeleteDefinition 删除事件定义
func DeleteDefinition(c *gin.Context) {
	id := c.Param("id")

	var definition model.EventDefinition
	if err := db.First(&definition, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Event definition not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to query event definition"})
		return
	}

	db.Delete(&definition)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Event definition deleted successfully",
	})
}
