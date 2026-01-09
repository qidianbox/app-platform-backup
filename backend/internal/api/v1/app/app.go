package app

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"

	"app-platform-backend/internal/model"
	"app-platform-backend/internal/pkg/database"

	"github.com/gin-gonic/gin"
)

func generateAppID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return "app_" + hex.EncodeToString(bytes)
}

func generateAppSecret() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func List(c *gin.Context) {
	var apps []model.App
	
	query := database.GetDB().Model(&model.App{})
	
	// 搜索
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("name LIKE ? OR app_id LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	
	// 状态筛选
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	
	// 分页
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	
	var total int64
	query.Count(&total)
	
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&apps)
	
	// 获取每个APP的模块数量
	type AppWithModules struct {
		model.App
		ModuleCount int64 `json:"module_count"`
		UserCount   int64 `json:"user_count"`
	}
	
	result := make([]AppWithModules, len(apps))
	for i, app := range apps {
		result[i].App = app
		database.GetDB().Model(&model.AppModule{}).Where("app_id = ? AND status = 1", app.ID).Count(&result[i].ModuleCount)
		result[i].UserCount = 0 // 暂时设为0
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":      result,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func Create(c *gin.Context) {
	var req struct {
		Name        string   `json:"name" binding:"required"`
		PackageName string   `json:"package_name"`
		Description string   `json:"description"`
		Icon        string   `json:"icon"`
		Modules     []string `json:"modules"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	app := model.App{
		Name:        req.Name,
		AppID:       generateAppID(),
		AppSecret:   generateAppSecret(),
		PackageName: req.PackageName,
		Description: req.Description,
		Icon:        req.Icon,
		Status:      1,
	}
	
	if err := database.GetDB().Create(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create app"})
		return
	}
	
	// 启用选中的模块
	for _, moduleCode := range req.Modules {
		appModule := model.AppModule{
			AppID:      app.ID,
			ModuleCode: moduleCode,
			Config:     "{}",
			Status:     1,
		}
		database.GetDB().Create(&appModule)
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": app,
	})
}

func Detail(c *gin.Context) {
	id := c.Param("id")
	
	var app model.App
	if err := database.GetDB().First(&app, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "App not found"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": app,
	})
}

func Update(c *gin.Context) {
	id := c.Param("id")
	
	var app model.App
	if err := database.GetDB().First(&app, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "App not found"})
		return
	}
	
	var req struct {
		Name        string `json:"name"`
		PackageName string `json:"package_name"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
		Status      *int   `json:"status"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.PackageName != "" {
		updates["package_name"] = req.PackageName
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	
	database.GetDB().Model(&app).Updates(updates)
	
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": app,
	})
}

func Delete(c *gin.Context) {
	id := c.Param("id")
	
	if err := database.GetDB().Delete(&model.App{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete app"})
		return
	}
	
	// 删除关联的模块
	database.GetDB().Where("app_id = ?", id).Delete(&model.AppModule{})
	
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "App deleted successfully",
	})
}

func ResetSecret(c *gin.Context) {
	id := c.Param("id")
	
	var app model.App
	if err := database.GetDB().First(&app, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "App not found"})
		return
	}
	
	newSecret := generateAppSecret()
	database.GetDB().Model(&app).Update("app_secret", newSecret)
	
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"app_secret": newSecret,
		},
	})
}
