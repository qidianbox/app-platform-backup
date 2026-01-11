package version

import (
	"app-platform-backend/internal/model"
	"app-platform-backend/internal/response"
	"app-platform-backend/internal/validator"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(database *gorm.DB) {
	db = database
}

// List 版本列表
func List(c *gin.Context) {
	appID := c.Query("app_id")
	status := c.Query("status")
	page, size := validator.ParsePagination(c.DefaultQuery("page", "1"), c.DefaultQuery("size", "20"))

	if appID == "" {
		response.ParamError(c, "app_id 不能为空")
		return
	}

	// 验证分页参数
	page, size = validator.ValidatePagination(page, size)

	query := db.Model(&model.Version{}).Where("app_id = ?", appID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		response.DBError(c, err)
		return
	}

	var versions []model.Version
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("created_at DESC").Find(&versions).Error; err != nil {
		response.DBError(c, err)
		return
	}

	// 转换为前端需要的格式
	var result []map[string]interface{}
	for _, v := range versions {
		result = append(result, map[string]interface{}{
			"id":           v.ID,
			"app_id":       v.AppID,
			"version":      v.VersionName,
			"version_code": v.VersionCode,
			"platform":     "android", // 默认平台
			"description":  v.Description,
			"download_url": v.DownloadURL,
			"force_update": v.IsForceUpdate == 1,
			"status":       v.Status,
			"published_at": v.PublishedAt,
			"created_at":   v.CreatedAt,
		})
	}

	response.PageSuccess(c, result, total, page, size)
}

// Create 创建版本
func Create(c *gin.Context) {
	var req struct {
		AppID       uint   `json:"app_id" binding:"required"`
		Version     string `json:"version" binding:"required"`
		Platform    string `json:"platform"`
		Description string `json:"description"`
		DownloadURL string `json:"download_url"`
		ForceUpdate bool   `json:"force_update"`
		GrayRelease bool   `json:"gray_release"`
		GrayPercent int    `json:"gray_percent"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, "参数错误: "+err.Error())
		return
	}

	// 验证版本号格式
	if len(req.Version) < 1 || len(req.Version) > 20 {
		response.ParamError(c, "版本号长度应在1-20个字符之间")
		return
	}

	// 验证下载URL
	if req.DownloadURL != "" {
		if err := validator.ValidateURL(req.DownloadURL); err != nil {
			response.ParamError(c, "下载URL格式错误")
			return
		}
	}

	// 获取最大版本号
	var maxVersionCode int
	db.Model(&model.Version{}).Where("app_id = ?", req.AppID).Select("COALESCE(MAX(version_code), 0)").Scan(&maxVersionCode)

	forceUpdate := 0
	if req.ForceUpdate {
		forceUpdate = 1
	}

	version := model.Version{
		AppID:         req.AppID,
		VersionName:   req.Version,
		VersionCode:   maxVersionCode + 1,
		Description:   req.Description,
		DownloadURL:   req.DownloadURL,
		IsForceUpdate: forceUpdate,
		Status:        "draft",
	}

	if err := db.Create(&version).Error; err != nil {
		response.DBError(c, err)
		return
	}

	response.SuccessWithMessage(c, version, "版本创建成功")
}

// Update 更新版本
func Update(c *gin.Context) {
	idStr := c.Param("id")
	appIDStr := c.Query("app_id")

	// 验证ID
	id, err := validator.ValidateID(idStr)
	if err != nil {
		response.ParamError(c, "无效的版本ID")
		return
	}

	// 验证app_id
	if appIDStr == "" {
		response.ParamError(c, "app_id 不能为空")
		return
	}
	appID, err := strconv.ParseUint(appIDStr, 10, 32)
	if err != nil {
		response.ParamError(c, "无效的 app_id")
		return
	}

	// 检查版本是否存在且属于该APP
	var existingVersion model.Version
	if err := db.Where("id = ? AND app_id = ?", id, appID).First(&existingVersion).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.NotFound(c, "版本不存在或无权限操作")
			return
		}
		response.DBError(c, err)
		return
	}

	var req struct {
		Version     string `json:"version"`
		Description string `json:"description"`
		DownloadURL string `json:"download_url"`
		ForceUpdate bool   `json:"force_update"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, "参数错误: "+err.Error())
		return
	}

	// 验证下载URL
	if req.DownloadURL != "" {
		if err := validator.ValidateURL(req.DownloadURL); err != nil {
			response.ParamError(c, "下载URL格式错误")
			return
		}
	}

	updates := map[string]interface{}{}
	if req.Version != "" {
		updates["version_name"] = req.Version
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.DownloadURL != "" {
		updates["download_url"] = req.DownloadURL
	}
	updates["is_force_update"] = 0
	if req.ForceUpdate {
		updates["is_force_update"] = 1
	}

	if err := db.Model(&model.Version{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		response.DBError(c, err)
		return
	}

	response.SuccessWithMessage(c, nil, "版本更新成功")
}

// Publish 发布版本
func Publish(c *gin.Context) {
	idStr := c.Param("id")
	appIDStr := c.Query("app_id")

	// 验证ID
	id, err := validator.ValidateID(idStr)
	if err != nil {
		response.ParamError(c, "无效的版本ID")
		return
	}

	// 验证app_id
	if appIDStr == "" {
		response.ParamError(c, "app_id 不能为空")
		return
	}
	appID, err := strconv.ParseUint(appIDStr, 10, 32)
	if err != nil {
		response.ParamError(c, "无效的 app_id")
		return
	}

	// 检查版本是否存在且属于该APP
	var existingVersion model.Version
	if err := db.Where("id = ? AND app_id = ?", id, appID).First(&existingVersion).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.NotFound(c, "版本不存在或无权限操作")
			return
		}
		response.DBError(c, err)
		return
	}

	if existingVersion.Status == "published" {
		response.ParamError(c, "版本已发布")
		return
	}

	now := time.Now()
	if err := db.Model(&model.Version{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       "published",
		"published_at": now,
	}).Error; err != nil {
		response.DBError(c, err)
		return
	}

	response.SuccessWithMessage(c, nil, "版本发布成功")
}

// Offline 下线版本
func Offline(c *gin.Context) {
	idStr := c.Param("id")
	appIDStr := c.Query("app_id")

	// 验证ID
	id, err := validator.ValidateID(idStr)
	if err != nil {
		response.ParamError(c, "无效的版本ID")
		return
	}

	// 验证app_id
	if appIDStr == "" {
		response.ParamError(c, "app_id 不能为空")
		return
	}
	appID, err := strconv.ParseUint(appIDStr, 10, 32)
	if err != nil {
		response.ParamError(c, "无效的 app_id")
		return
	}

	// 检查版本是否存在且属于该APP
	var existingVersion model.Version
	if err := db.Where("id = ? AND app_id = ?", id, appID).First(&existingVersion).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.NotFound(c, "版本不存在或无权限操作")
			return
		}
		response.DBError(c, err)
		return
	}

	if existingVersion.Status == "offline" {
		response.ParamError(c, "版本已下线")
		return
	}

	if err := db.Model(&model.Version{}).Where("id = ?", id).Update("status", "offline").Error; err != nil {
		response.DBError(c, err)
		return
	}

	response.SuccessWithMessage(c, nil, "版本下线成功")
}

// Delete 删除版本
func Delete(c *gin.Context) {
	idStr := c.Param("id")
	appIDStr := c.Query("app_id")

	// 验证ID
	id, err := validator.ValidateID(idStr)
	if err != nil {
		response.ParamError(c, "无效的版本ID")
		return
	}

	// 验证app_id
	if appIDStr == "" {
		response.ParamError(c, "app_id 不能为空")
		return
	}
	appID, err := strconv.ParseUint(appIDStr, 10, 32)
	if err != nil {
		response.ParamError(c, "无效的 app_id")
		return
	}

	// 检查版本是否存在且属于该APP
	var existingVersion model.Version
	if err := db.Where("id = ? AND app_id = ?", id, appID).First(&existingVersion).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.NotFound(c, "版本不存在或无权限删除")
			return
		}
		response.DBError(c, err)
		return
	}

	// 不允许删除已发布的版本
	if existingVersion.Status == "published" {
		response.ParamError(c, "不能删除已发布的版本，请先下线")
		return
	}

	if err := db.Delete(&model.Version{}, id).Error; err != nil {
		response.DBError(c, err)
		return
	}

	response.SuccessWithMessage(c, nil, "版本删除成功")
}

// CheckUpdate 检查更新
func CheckUpdate(c *gin.Context) {
	appID := c.Query("app_id")
	currentVersion := c.Query("version")

	if appID == "" {
		response.ParamError(c, "app_id 不能为空")
		return
	}

	// 获取最新发布的版本
	var latestVersion model.Version
	err := db.Where("app_id = ? AND status = ?", appID, "published").
		Order("version_code DESC").
		First(&latestVersion).Error

	if err == gorm.ErrRecordNotFound {
		response.Success(c, gin.H{
			"has_update": false,
		})
		return
	}

	if err != nil {
		response.DBError(c, err)
		return
	}

	hasUpdate := latestVersion.VersionName != currentVersion

	response.Success(c, gin.H{
		"has_update":   hasUpdate,
		"version":      latestVersion.VersionName,
		"version_code": latestVersion.VersionCode,
		"description":  latestVersion.Description,
		"download_url": latestVersion.DownloadURL,
		"force_update": latestVersion.IsForceUpdate == 1,
	})
}

// Stats 版本统计
func Stats(c *gin.Context) {
	appID := c.Query("app_id")
	if appID == "" {
		response.ParamError(c, "app_id 不能为空")
		return
	}

	var total, published, draft, offline int64
	db.Model(&model.Version{}).Where("app_id = ?", appID).Count(&total)
	db.Model(&model.Version{}).Where("app_id = ? AND status = ?", appID, "published").Count(&published)
	db.Model(&model.Version{}).Where("app_id = ? AND status = ?", appID, "draft").Count(&draft)
	db.Model(&model.Version{}).Where("app_id = ? AND status = ?", appID, "offline").Count(&offline)

	response.Success(c, gin.H{
		"total":     total,
		"published": published,
		"draft":     draft,
		"offline":   offline,
	})
}
