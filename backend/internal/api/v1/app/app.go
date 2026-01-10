package app

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"

	"app-platform-backend/internal/model"
	"app-platform-backend/internal/pkg/database"
	"app-platform-backend/internal/response"
	"app-platform-backend/internal/validator"

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

// List 获取APP列表
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

	// 分页参数验证
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	page, pageSize = validator.ValidatePagination(page, pageSize)

	var total int64
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&apps).Error; err != nil {
		response.DBError(c, err)
		return
	}

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

	response.PageSuccess(c, result, total, page, pageSize)
}

// Create 创建APP
func Create(c *gin.Context) {
	var req validator.AppCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	// 验证请求参数
	if err := validator.ValidateAppCreate(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	// 获取实际的名称
	appName := req.Name
	if appName == "" {
		appName = req.AppName
	}

	app := model.App{
		Name:        appName,
		AppID:       generateAppID(),
		AppSecret:   generateAppSecret(),
		PackageName: req.PackageName,
		Description: req.Description,
		Icon:        req.Icon,
		Status:      1,
	}

	// 使用事务创建APP和关联模块
	err := database.WithTransaction(func(tx *database.DB) error {
		if err := tx.Create(&app).Error; err != nil {
			return err
		}

		// 启用选中的模块
		for _, sourceModule := range req.Modules {
			appModule := model.AppModule{
				AppID:        app.ID,
				ModuleCode:   sourceModule,
				SourceModule: sourceModule,
				Config:       "{}",
				Status:       1,
			}
			if err := tx.Create(&appModule).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		response.DBError(c, err)
		return
	}

	response.Success(c, app)
}

// Detail 获取APP详情
func Detail(c *gin.Context) {
	id := c.Param("id")

	// 验证ID
	if _, err := validator.ValidateID(id); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	var app model.App
	if err := database.GetDB().First(&app, id).Error; err != nil {
		response.NotFound(c, "应用不存在")
		return
	}

	response.Success(c, app)
}

// Update 更新APP
func Update(c *gin.Context) {
	id := c.Param("id")

	// 验证ID
	if _, err := validator.ValidateID(id); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	var app model.App
	if err := database.GetDB().First(&app, id).Error; err != nil {
		response.NotFound(c, "应用不存在")
		return
	}

	var req validator.AppUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	// 验证更新参数
	if err := validator.ValidateAppUpdate(&req); err != nil {
		response.ParamError(c, err.Error())
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

	if err := database.GetDB().Model(&app).Updates(updates).Error; err != nil {
		response.DBError(c, err)
		return
	}

	// 重新查询获取更新后的数据
	database.GetDB().First(&app, id)
	response.Success(c, app)
}

// Delete 删除APP
func Delete(c *gin.Context) {
	id := c.Param("id")

	// 验证ID
	if _, err := validator.ValidateID(id); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	// 使用事务删除APP和关联数据
	err := database.WithTransaction(func(tx *database.DB) error {
		if err := tx.Delete(&model.App{}, id).Error; err != nil {
			return err
		}
		// 删除关联的模块
		if err := tx.Where("app_id = ?", id).Delete(&model.AppModule{}).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		response.DBError(c, err)
		return
	}

	response.SuccessWithMessage(c, nil, "应用删除成功")
}

// ResetSecret 重置APP密钥
func ResetSecret(c *gin.Context) {
	id := c.Param("id")

	// 验证ID
	if _, err := validator.ValidateID(id); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	var app model.App
	if err := database.GetDB().First(&app, id).Error; err != nil {
		response.NotFound(c, "应用不存在")
		return
	}

	newSecret := generateAppSecret()
	if err := database.GetDB().Model(&app).Update("app_secret", newSecret).Error; err != nil {
		response.DBError(c, err)
		return
	}

	response.Success(c, gin.H{
		"app_secret": newSecret,
	})
}
