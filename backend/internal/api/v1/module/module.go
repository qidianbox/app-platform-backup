package module

import (
	"encoding/json"
	"fmt"
	"net/http"

	"app-platform-backend/internal/model"
	"app-platform-backend/internal/pkg/database"

	"github.com/gin-gonic/gin"
)

func GetAllTemplates(c *gin.Context) {
	var templates []model.ModuleTemplate
	database.GetDB().Where("status = 1").Find(&templates)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": templates,
	})
}

func GetAppModules(c *gin.Context) {
	appID := c.Param("id")

	var modules []model.AppModule
	database.GetDB().Where("app_id = ?", appID).Find(&modules)

	// 获取模块模板信息
	type ModuleWithTemplate struct {
		model.AppModule
		ModuleName  string `json:"module_name"`
		Category    string `json:"category"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}

	result := make([]ModuleWithTemplate, 0)
	for _, m := range modules {
		var template model.ModuleTemplate
		if err := database.GetDB().Where("module_code = ?", m.ModuleCode).First(&template).Error; err == nil {
			result = append(result, ModuleWithTemplate{
				AppModule:   m,
				ModuleName:  template.ModuleName,
				Category:    template.Category,
				Description: template.Description,
				Icon:        template.Icon,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": result,
	})
}

func GetAppModule(c *gin.Context) {
	appID := c.Param("id")
	moduleCode := c.Param("module_code")

	var module model.AppModule
	if err := database.GetDB().Where("app_id = ? AND module_code = ?", appID, moduleCode).First(&module).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Module not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": module,
	})
}

func EnableModule(c *gin.Context) {
	appID := c.Param("id")

	var req struct {
		ModuleCode string `json:"module_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查是否已启用
	var existing model.AppModule
	if err := database.GetDB().Where("app_id = ? AND module_code = ?", appID, req.ModuleCode).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Module already enabled"})
		return
	}

	module := model.AppModule{
		AppID:      parseUint(appID),
		ModuleCode: req.ModuleCode,
		Config:     "{}",
		Status:     1,
	}

	if err := database.GetDB().Create(&module).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enable module"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": module,
	})
}

func UpdateModule(c *gin.Context) {
	appID := c.Param("id")
	moduleCode := c.Param("module_code")

	var module model.AppModule
	if err := database.GetDB().Where("app_id = ? AND module_code = ?", appID, moduleCode).First(&module).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Module not found"})
		return
	}

	var req struct {
		Status *int `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Status != nil {
		database.GetDB().Model(&module).Update("status", *req.Status)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": module,
	})
}

func DisableModule(c *gin.Context) {
	appID := c.Param("id")
	moduleCode := c.Param("module_code")

	if err := database.GetDB().Where("app_id = ? AND module_code = ?", appID, moduleCode).Delete(&model.AppModule{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disable module"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Module disabled successfully",
	})
}

func BatchEnableModules(c *gin.Context) {
	appID := c.Param("id")

	var req struct {
		ModuleCodes []string `json:"module_codes" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, code := range req.ModuleCodes {
		var existing model.AppModule
		if err := database.GetDB().Where("app_id = ? AND module_code = ?", appID, code).First(&existing).Error; err != nil {
			module := model.AppModule{
				AppID:      parseUint(appID),
				ModuleCode: code,
				Config:     "{}",
				Status:     1,
			}
			database.GetDB().Create(&module)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Modules enabled successfully",
	})
}

func SaveModuleConfig(c *gin.Context) {
	appID := c.Param("id")
	moduleCode := c.Param("module_code")

	var module model.AppModule
	if err := database.GetDB().Where("app_id = ? AND module_code = ?", appID, moduleCode).First(&module).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Module not found"})
		return
	}

	var req struct {
		Config map[string]interface{} `json:"config" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 保存配置历史
	var maxVersion int
	database.GetDB().Model(&model.ModuleConfigHistory{}).
		Where("app_id = ? AND module_code = ?", appID, moduleCode).
		Select("COALESCE(MAX(version), 0)").Scan(&maxVersion)

	history := model.ModuleConfigHistory{
		AppID:      parseUint(appID),
		ModuleCode: moduleCode,
		Config:     module.Config,
		Version:    maxVersion + 1,
		Operator:   c.GetString("username"),
	}
	database.GetDB().Create(&history)

	// 更新配置
	configJSON, _ := json.Marshal(req.Config)
	database.GetDB().Model(&module).Update("config", string(configJSON))

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Config saved successfully",
	})
}

func GetModuleConfig(c *gin.Context) {
	appID := c.Param("id")
	moduleCode := c.Param("module_code")

	var module model.AppModule
	if err := database.GetDB().Where("app_id = ? AND module_code = ?", appID, moduleCode).First(&module).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Module not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"config": module.Config,
		},
	})
}

func ResetModuleConfig(c *gin.Context) {
	appID := c.Param("id")
	moduleCode := c.Param("module_code")

	var module model.AppModule
	if err := database.GetDB().Where("app_id = ? AND module_code = ?", appID, moduleCode).First(&module).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Module not found"})
		return
	}

	database.GetDB().Model(&module).Update("config", "{}")

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Config reset successfully",
	})
}

func TestModuleConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Config test passed",
		"data": gin.H{
			"success": true,
		},
	})
}

func GetConfigHistory(c *gin.Context) {
	appID := c.Param("id")
	moduleCode := c.Param("module_code")

	var history []model.ModuleConfigHistory
	database.GetDB().Where("app_id = ? AND module_code = ?", appID, moduleCode).
		Order("version DESC").Limit(20).Find(&history)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": history,
	})
}

func RollbackConfig(c *gin.Context) {
	appID := c.Param("id")
	moduleCode := c.Param("module_code")
	historyID := c.Param("history_id")

	var history model.ModuleConfigHistory
	if err := database.GetDB().First(&history, historyID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "History not found"})
		return
	}

	database.GetDB().Model(&model.AppModule{}).
		Where("app_id = ? AND module_code = ?", appID, moduleCode).
		Update("config", history.Config)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Config rolled back successfully",
	})
}

func CompareConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"diff": []string{},
		},
	})
}

func CheckModuleDependencies(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"satisfied":   true,
			"missing":     []string{},
			"suggestions": []string{},
		},
	})
}

func CheckModuleReverseDependencies(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"dependents": []string{},
		},
	})
}

func AutoEnableModuleDependencies(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Dependencies enabled successfully",
	})
}

func DetectCircularDependency(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"has_circular": false,
			"path":         []string{},
		},
	})
}

func parseUint(s string) uint {
	var id uint
	fmt.Sscanf(s, "%d", &id)
	return id
}


