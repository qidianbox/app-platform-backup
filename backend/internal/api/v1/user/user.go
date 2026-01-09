package user

import (
	"app-platform-backend/internal/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitDB 初始化数据库连接
func InitDB(database *gorm.DB) {
	db = database
}

// ListRequest 用户列表请求参数
type ListRequest struct {
	AppID  uint   `form:"app_id" binding:"required"`
	Page   int    `form:"page" binding:"min=1"`
	Size   int    `form:"size" binding:"min=1,max=100"`
	Status *int   `form:"status"`
	Search string `form:"search"`
}

// List 获取用户列表
func List(c *gin.Context) {
	var req ListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// 默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Size == 0 {
		req.Size = 20
	}

	query := db.Model(&model.User{}).Where("app_id = ?", req.AppID)

	// 状态筛选
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	// 搜索
	if req.Search != "" {
		query = query.Where("nickname LIKE ? OR phone LIKE ? OR email LIKE ?",
			"%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to count users"})
		return
	}

	// 分页查询
	var users []model.User
	offset := (req.Page - 1) * req.Size
	if err := query.Offset(offset).Limit(req.Size).Order("created_at DESC").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to query users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  users,
			"total": total,
			"page":  req.Page,
			"size":  req.Size,
		},
	})
}

// Detail 获取用户详情
func Detail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid user ID"})
		return
	}

	var user model.User
	if err := db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to query user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": user,
	})
}

// UpdateStatusRequest 更新用户状态请求参数
type UpdateStatusRequest struct {
	Status int `json:"status" binding:"required,oneof=0 1"`
}

// UpdateStatus 更新用户状态
func UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid user ID"})
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// 更新状态
	if err := db.Model(&model.User{}).Where("id = ?", id).Update("status", req.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to update user status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "User status updated successfully",
	})
}

// Stats 用户统计
func Stats(c *gin.Context) {
	appIDStr := c.Query("app_id")
	if appIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "app_id is required"})
		return
	}

	appID, err := strconv.ParseUint(appIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Invalid app_id"})
		return
	}

	// 总用户数
	var total int64
	db.Model(&model.User{}).Where("app_id = ?", appID).Count(&total)

	// 活跃用户数（最近7天登录）
	var active int64
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	db.Model(&model.User{}).Where("app_id = ? AND last_login_at > ?", appID, sevenDaysAgo).Count(&active)

	// 今日新增
	var todayNew int64
	today := time.Now().Format("2006-01-02")
	db.Model(&model.User{}).Where("app_id = ? AND DATE(created_at) = ?", appID, today).Count(&todayNew)

	// 禁用用户数
	var disabled int64
	db.Model(&model.User{}).Where("app_id = ? AND status = 0", appID).Count(&disabled)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"total":      total,
			"active":     active,
			"today_new":  todayNew,
			"disabled":   disabled,
			"normal":     total - disabled,
		},
	})
}
