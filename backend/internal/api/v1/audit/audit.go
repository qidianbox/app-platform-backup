package audit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"app-platform-backend/internal/scheduler"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

// AuditLog 审计日志模型
type AuditLog struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	AppID       uint      `json:"app_id" gorm:"index"`
	UserID      string    `json:"user_id" gorm:"index"`
	UserName    string    `json:"user_name"`
	Action      string    `json:"action" gorm:"index"`       // login, logout, create, update, delete, view, export, config
	Resource    string    `json:"resource" gorm:"index"`     // user, app, config, message, push, file, version, etc.
	ResourceID  string    `json:"resource_id"`
	Description string    `json:"description"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	RequestPath string    `json:"request_path"`
	RequestMethod string  `json:"request_method"`
	StatusCode  int       `json:"status_code"`
	Duration    int64     `json:"duration"` // 请求耗时（毫秒）
	Extra       string    `json:"extra" gorm:"type:text"` // JSON格式的额外信息
	CreatedAt   time.Time `json:"created_at" gorm:"index"`
}

// InitDB 初始化数据库
func InitDB(database *gorm.DB) {
	db = database
	db.AutoMigrate(&AuditLog{})
}

// RecordAudit 记录审计日志
func RecordAudit(c *gin.Context, action, resource, resourceID, description string, extra map[string]interface{}) {
	userID := c.GetString("user_id")
	userName := c.GetString("user_name")
	appIDStr := c.Param("app_id")
	if appIDStr == "" {
		appIDStr = c.Query("app_id")
	}
	
	var appID uint
	if id, err := strconv.ParseUint(appIDStr, 10, 32); err == nil {
		appID = uint(id)
	}

	extraJSON := ""
	if extra != nil {
		if data, err := json.Marshal(extra); err == nil {
			extraJSON = string(data)
		}
	}

	log := &AuditLog{
		AppID:         appID,
		UserID:        userID,
		UserName:      userName,
		Action:        action,
		Resource:      resource,
		ResourceID:    resourceID,
		Description:   description,
		IPAddress:     c.ClientIP(),
		UserAgent:     c.Request.UserAgent(),
		RequestPath:   c.Request.URL.Path,
		RequestMethod: c.Request.Method,
		StatusCode:    c.Writer.Status(),
		Extra:         extraJSON,
		CreatedAt:     time.Now(),
	}

	go func() {
		if err := db.Create(log).Error; err != nil {
			// 静默处理错误，不影响主流程
		}
	}()
}

// List 获取审计日志列表
func List(c *gin.Context) {
	appIDStr := c.Query("app_id")
	userID := c.Query("user_id")
	action := c.Query("action")
	resource := c.Query("resource")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := db.Model(&AuditLog{})

	if appIDStr != "" {
		if appID, err := strconv.ParseUint(appIDStr, 10, 32); err == nil {
			query = query.Where("app_id = ?", appID)
		}
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if resource != "" {
		query = query.Where("resource = ?", resource)
	}
	if startTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", startTime); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if endTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", endTime); err == nil {
			query = query.Where("created_at <= ?", t)
		}
	}
	if keyword != "" {
		query = query.Where("description LIKE ? OR user_name LIKE ? OR ip_address LIKE ?", 
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var logs []AuditLog
	query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&logs)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":      logs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// Stats 获取审计日志统计
func Stats(c *gin.Context) {
	appIDStr := c.Query("app_id")
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))

	query := db.Model(&AuditLog{})
	if appIDStr != "" {
		if appID, err := strconv.ParseUint(appIDStr, 10, 32); err == nil {
			query = query.Where("app_id = ?", appID)
		}
	}

	startTime := time.Now().AddDate(0, 0, -days)
	query = query.Where("created_at >= ?", startTime)

	// 总操作数
	var totalCount int64
	query.Count(&totalCount)

	// 按操作类型统计
	type ActionStat struct {
		Action string `json:"action"`
		Count  int64  `json:"count"`
	}
	var actionStats []ActionStat
	db.Model(&AuditLog{}).
		Select("action, COUNT(*) as count").
		Where("created_at >= ?", startTime).
		Group("action").
		Find(&actionStats)

	// 按资源类型统计
	type ResourceStat struct {
		Resource string `json:"resource"`
		Count    int64  `json:"count"`
	}
	var resourceStats []ResourceStat
	db.Model(&AuditLog{}).
		Select("resource, COUNT(*) as count").
		Where("created_at >= ?", startTime).
		Group("resource").
		Find(&resourceStats)

	// 按用户统计（Top 10）
	type UserStat struct {
		UserID   string `json:"user_id"`
		UserName string `json:"user_name"`
		Count    int64  `json:"count"`
	}
	var userStats []UserStat
	db.Model(&AuditLog{}).
		Select("user_id, user_name, COUNT(*) as count").
		Where("created_at >= ?", startTime).
		Group("user_id, user_name").
		Order("count DESC").
		Limit(10).
		Find(&userStats)

	// 按天统计趋势
	type DailyStat struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	var dailyStats []DailyStat
	db.Model(&AuditLog{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= ?", startTime).
		Group("DATE(created_at)").
		Order("date").
		Find(&dailyStats)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"total_count":    totalCount,
			"action_stats":   actionStats,
			"resource_stats": resourceStats,
			"user_stats":     userStats,
			"daily_stats":    dailyStats,
		},
	})
}

// Export 导出审计日志
func Export(c *gin.Context) {
	appIDStr := c.Query("app_id")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")
	format := c.DefaultQuery("format", "csv")

	query := db.Model(&AuditLog{})
	if appIDStr != "" {
		if appID, err := strconv.ParseUint(appIDStr, 10, 32); err == nil {
			query = query.Where("app_id = ?", appID)
		}
	}
	if startTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", startTime); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if endTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", endTime); err == nil {
			query = query.Where("created_at <= ?", t)
		}
	}

	var logs []AuditLog
	query.Order("created_at DESC").Limit(10000).Find(&logs)

	if format == "csv" {
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", "attachment; filename=audit_logs.csv")
		
		// CSV头
		c.Writer.WriteString("ID,时间,用户ID,用户名,操作,资源,资源ID,描述,IP地址,请求路径,状态码\n")
		
		for _, log := range logs {
			line := fmt.Sprintf("%d,%s,%s,%s,%s,%s,%s,%s,%s,%s,%d\n",
				log.ID,
				log.CreatedAt.Format("2006-01-02 15:04:05"),
				log.UserID,
				log.UserName,
				log.Action,
				log.Resource,
				log.ResourceID,
				log.Description,
				log.IPAddress,
				log.RequestPath,
				log.StatusCode,
			)
			c.Writer.WriteString(line)
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"data": logs,
		})
	}
}

// Cleanup 手动清理审计日志
func Cleanup(c *gin.Context) {
	retentionDays, _ := strconv.Atoi(c.DefaultQuery("retention_days", "90"))
	if retentionDays < 7 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "保留天数不能小于7天",
		})
		return
	}

	s := scheduler.GetScheduler()
	if s == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "清理调度器未初始化",
		})
		return
	}

	deletedRows, err := s.ManualCleanup(retentionDays)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "清理失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"deleted_rows":   deletedRows,
			"retention_days": retentionDays,
		},
		"message": "清理完成",
	})
}

// CleanupHistory 获取清理历史记录
func CleanupHistory(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	s := scheduler.GetScheduler()
	if s == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "清理调度器未初始化",
		})
		return
	}

	records, err := s.GetCleanupHistory(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取清理历史失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": records,
	})
}

// CleanupConfig 获取清理配置
func CleanupConfig(c *gin.Context) {
	s := scheduler.GetScheduler()
	if s == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "清理调度器未初始化",
		})
		return
	}

	config := s.GetConfig()
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"retention_days": config.RetentionDays,
			"cleanup_hour":   config.CleanupHour,
			"batch_size":     config.BatchSize,
		},
	})
}


