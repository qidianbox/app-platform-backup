package monitor

import (
	"app-platform-backend/internal/model"
	"app-platform-backend/internal/response"
	"app-platform-backend/internal/validator"
	"encoding/json"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB(database *gorm.DB) {
	db = database
}

// ReportMetric 上报监控指标
func ReportMetric(c *gin.Context) {
	var req struct {
		AppID       uint              `json:"app_id" binding:"required"`
		MetricName  string            `json:"metric_name" binding:"required"`
		MetricValue float64           `json:"metric_value" binding:"required"`
		Tags        map[string]string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, "参数错误: "+err.Error())
		return
	}

	// 验证指标名称
	if len(req.MetricName) < 1 || len(req.MetricName) > 100 {
		response.ParamError(c, "指标名称长度应在1-100个字符之间")
		return
	}

	tagsJSON := "{}"
	if req.Tags != nil {
		if data, err := json.Marshal(req.Tags); err == nil {
			tagsJSON = string(data)
		}
	}

	metric := model.MonitorMetric{
		AppID:       req.AppID,
		MetricName:  req.MetricName,
		MetricValue: req.MetricValue,
		Tags:        tagsJSON,
	}

	if err := db.Create(&metric).Error; err != nil {
		response.DBError(c, err)
		return
	}

	// 检查是否触发告警
	checkAlerts(req.AppID, req.MetricName, req.MetricValue)

	response.SuccessWithMessage(c, nil, "指标上报成功")
}

// 检查告警规则
func checkAlerts(appID uint, metricName string, value float64) {
	var alerts []model.MonitorAlert
	db.Where("app_id = ? AND metric_name = ? AND is_active = 1", appID, metricName).Find(&alerts)

	for _, alert := range alerts {
		triggered := false
		switch alert.Condition {
		case "gt":
			triggered = value > alert.Threshold
		case "gte":
			triggered = value >= alert.Threshold
		case "lt":
			triggered = value < alert.Threshold
		case "lte":
			triggered = value <= alert.Threshold
		case "eq":
			triggered = value == alert.Threshold
		}

		if triggered {
			now := time.Now()
			db.Model(&alert).Updates(map[string]interface{}{
				"status":        "alerting",
				"last_alert_at": now,
			})
		}
	}
}

// Metrics 获取监控指标
func Metrics(c *gin.Context) {
	appID := c.Query("app_id")
	metricName := c.Query("metric_name")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "100"))

	if appID == "" {
		response.ParamError(c, "app_id 不能为空")
		return
	}

	// 验证分页参数
	page, size = validator.ValidatePagination(page, size)

	query := db.Model(&model.MonitorMetric{}).Where("app_id = ?", appID)

	if metricName != "" {
		query = query.Where("metric_name = ?", metricName)
	}
	if startTime != "" {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("created_at <= ?", endTime)
	}

	var total int64
	query.Count(&total)

	var metrics []model.MonitorMetric
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("created_at DESC").Find(&metrics).Error; err != nil {
		response.DBError(c, err)
		return
	}

	response.PageSuccess(c, metrics, total, page, size)
}

// MetricStats 指标统计
func MetricStats(c *gin.Context) {
	appID := c.Query("app_id")
	metricName := c.Query("metric_name")

	if appID == "" || metricName == "" {
		response.ParamError(c, "app_id 和 metric_name 不能为空")
		return
	}

	var stats struct {
		Avg   float64 `json:"avg"`
		Max   float64 `json:"max"`
		Min   float64 `json:"min"`
		Count int64   `json:"count"`
	}

	if err := db.Model(&model.MonitorMetric{}).
		Where("app_id = ? AND metric_name = ?", appID, metricName).
		Select("AVG(metric_value) as avg, MAX(metric_value) as max, MIN(metric_value) as min, COUNT(*) as count").
		Scan(&stats).Error; err != nil {
		response.DBError(c, err)
		return
	}

	// 获取最近的趋势数据
	var trends []struct {
		Time  time.Time `json:"time"`
		Value float64   `json:"value"`
	}
	db.Model(&model.MonitorMetric{}).
		Where("app_id = ? AND metric_name = ?", appID, metricName).
		Select("created_at as time, metric_value as value").
		Order("created_at DESC").
		Limit(100).
		Scan(&trends)

	response.Success(c, gin.H{
		"avg":    stats.Avg,
		"max":    stats.Max,
		"min":    stats.Min,
		"count":  stats.Count,
		"trends": trends,
	})
}

// Alerts 告警列表
func Alerts(c *gin.Context) {
	appID := c.Query("app_id")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	// 验证分页参数
	page, size = validator.ValidatePagination(page, size)

	query := db.Model(&model.MonitorAlert{})
	if appID != "" {
		query = query.Where("app_id = ?", appID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var alerts []model.MonitorAlert
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("created_at DESC").Find(&alerts).Error; err != nil {
		response.DBError(c, err)
		return
	}

	response.PageSuccess(c, alerts, total, page, size)
}

// CreateAlert 创建告警规则
func CreateAlert(c *gin.Context) {
	var req struct {
		AppID      uint    `json:"app_id" binding:"required"`
		AlertName  string  `json:"alert_name" binding:"required"`
		MetricName string  `json:"metric_name" binding:"required"`
		Condition  string  `json:"condition" binding:"required"`
		Threshold  float64 `json:"threshold" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, "参数错误: "+err.Error())
		return
	}

	// 验证告警名称长度
	if len(req.AlertName) < 2 || len(req.AlertName) > 100 {
		response.ParamError(c, "告警名称长度应在2-100个字符之间")
		return
	}

	// 验证条件
	validConditions := map[string]bool{"gt": true, "gte": true, "lt": true, "lte": true, "eq": true}
	if !validConditions[req.Condition] {
		response.ParamError(c, "无效的条件，请使用: gt, gte, lt, lte, eq")
		return
	}

	alert := model.MonitorAlert{
		AppID:      req.AppID,
		AlertName:  req.AlertName,
		MetricName: req.MetricName,
		Condition:  req.Condition,
		Threshold:  req.Threshold,
		Status:     "normal",
		IsActive:   1,
	}

	if err := db.Create(&alert).Error; err != nil {
		response.DBError(c, err)
		return
	}

	response.SuccessWithMessage(c, alert, "告警规则创建成功")
}

// UpdateAlert 更新告警规则
func UpdateAlert(c *gin.Context) {
	id := c.Param("id")
	appIDStr := c.Query("app_id")

	// 验证ID
	if _, err := validator.ValidateID(id); err != nil {
		response.ParamError(c, err.Error())
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

	var alert model.MonitorAlert
	// 同时验证id和app_id，防止越权操作
	if err := db.Where("id = ? AND app_id = ?", id, appID).First(&alert).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.NotFound(c, "告警规则不存在或无权限操作")
			return
		}
		response.DBError(c, err)
		return
	}

	var req struct {
		AlertName  string   `json:"alert_name"`
		MetricName string   `json:"metric_name"`
		Condition  string   `json:"condition"`
		Threshold  *float64 `json:"threshold"`
		IsActive   *int     `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, "参数错误: "+err.Error())
		return
	}

	// 验证条件
	if req.Condition != "" {
		validConditions := map[string]bool{"gt": true, "gte": true, "lt": true, "lte": true, "eq": true}
		if !validConditions[req.Condition] {
			response.ParamError(c, "无效的条件，请使用: gt, gte, lt, lte, eq")
			return
		}
	}

	updates := map[string]interface{}{}
	if req.AlertName != "" {
		updates["alert_name"] = req.AlertName
	}
	if req.MetricName != "" {
		updates["metric_name"] = req.MetricName
	}
	if req.Condition != "" {
		updates["condition"] = req.Condition
	}
	if req.Threshold != nil {
		updates["threshold"] = *req.Threshold
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if err := db.Model(&alert).Updates(updates).Error; err != nil {
		response.DBError(c, err)
		return
	}

	response.SuccessWithMessage(c, nil, "告警规则更新成功")
}

// DeleteAlert 删除告警规则
func DeleteAlert(c *gin.Context) {
	id := c.Param("id")
	appIDStr := c.Query("app_id")

	// 验证ID
	if _, err := validator.ValidateID(id); err != nil {
		response.ParamError(c, err.Error())
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

	var alert model.MonitorAlert
	// 同时验证id和app_id，防止越权删除
	if err := db.Where("id = ? AND app_id = ?", id, appID).First(&alert).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.NotFound(c, "告警规则不存在或无权限删除")
			return
		}
		response.DBError(c, err)
		return
	}

	if err := db.Delete(&alert).Error; err != nil {
		response.DBError(c, err)
		return
	}

	response.SuccessWithMessage(c, nil, "告警规则删除成功")
}

// ResolveAlert 解决告警
func ResolveAlert(c *gin.Context) {
	id := c.Param("id")
	appIDStr := c.Query("app_id")

	// 验证ID
	if _, err := validator.ValidateID(id); err != nil {
		response.ParamError(c, err.Error())
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

	var alert model.MonitorAlert
	// 同时验证id和app_id，防止越权操作
	if err := db.Where("id = ? AND app_id = ?", id, appID).First(&alert).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.NotFound(c, "告警规则不存在或无权限操作")
			return
		}
		response.DBError(c, err)
		return
	}

	if err := db.Model(&alert).Update("status", "normal").Error; err != nil {
		response.DBError(c, err)
		return
	}

	response.SuccessWithMessage(c, nil, "告警已解决")
}

// Rules 告警规则列表（兼容旧接口）
func Rules(c *gin.Context) {
	Alerts(c)
}

// Health 健康检查
func Health(c *gin.Context) {
	response.Success(c, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
		"uptime":    time.Since(startTime).Seconds(),
	})
}

// Stats 监控统计
func Stats(c *gin.Context) {
	appID := c.Query("app_id")
	if appID == "" {
		response.ParamError(c, "app_id 不能为空")
		return
	}

	var totalMetrics, totalAlerts, activeAlerts, alertingCount int64
	db.Model(&model.MonitorMetric{}).Where("app_id = ?", appID).Count(&totalMetrics)
	db.Model(&model.MonitorAlert{}).Where("app_id = ?", appID).Count(&totalAlerts)
	db.Model(&model.MonitorAlert{}).Where("app_id = ? AND is_active = 1", appID).Count(&activeAlerts)
	db.Model(&model.MonitorAlert{}).Where("app_id = ? AND status = ?", appID, "alerting").Count(&alertingCount)

	// 获取指标类型统计
	var metricStats []struct {
		MetricName string `json:"metric_name"`
		Count      int64  `json:"count"`
	}
	db.Model(&model.MonitorMetric{}).
		Select("metric_name, COUNT(*) as count").
		Where("app_id = ?", appID).
		Group("metric_name").
		Order("count DESC").
		Limit(10).
		Scan(&metricStats)

	response.Success(c, gin.H{
		"total_metrics":  totalMetrics,
		"total_alerts":   totalAlerts,
		"active_alerts":  activeAlerts,
		"alerting_count": alertingCount,
		"metric_stats":   metricStats,
	})
}

var startTime = time.Now()
