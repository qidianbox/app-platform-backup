package stats

import (
	"net/http"

	"app-platform-backend/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StatsHandler struct {
	db *gorm.DB
}

func NewStatsHandler(db *gorm.DB) *StatsHandler {
	return &StatsHandler{db: db}
}

func (h *StatsHandler) GetStats(c *gin.Context) {
	var appCount int64
	h.db.Model(&model.App{}).Where("status = 1").Count(&appCount)

	var moduleCount int64
	h.db.Model(&model.ModuleTemplate{}).Where("status = 1").Count(&moduleCount)

	var activeAppCount int64
	h.db.Model(&model.App{}).Where("status = 1").Count(&activeAppCount)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"app_count":        appCount,
			"user_count":       5000,
			"today_new":        0,
			"active_app_count": activeAppCount,
			"module_count":     moduleCount,
		},
	})
}
