package file

import (
	"app-platform-backend/internal/model"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB
var uploadDir = "/tmp/uploads"

func InitDB(database *gorm.DB) {
	db = database
	// 确保上传目录存在
	os.MkdirAll(uploadDir, 0755)
}

// Upload 上传文件
func Upload(c *gin.Context) {
	appIDStr := c.PostForm("app_id")
	if appIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "app_id is required"})
		return
	}

	appID, _ := strconv.ParseUint(appIDStr, 10, 32)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "No file uploaded"})
		return
	}
	defer file.Close()

	// 生成唯一文件名
	ext := filepath.Ext(header.Filename)
	hash := md5.New()
	io.Copy(hash, file)
	file.Seek(0, 0)
	
	timestamp := time.Now().UnixNano()
	newFilename := fmt.Sprintf("%x_%d%s", hash.Sum(nil), timestamp, ext)
	
	// 按APP和日期组织目录
	dateDir := time.Now().Format("2006/01/02")
	fullDir := filepath.Join(uploadDir, fmt.Sprintf("%d", appID), dateDir)
	os.MkdirAll(fullDir, 0755)
	
	filePath := filepath.Join(fullDir, newFilename)
	
	// 保存文件
	out, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to save file"})
		return
	}
	defer out.Close()
	
	io.Copy(out, file)

	// 获取MIME类型
	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// 保存到数据库
	fileRecord := model.File{
		AppID:    uint(appID),
		Filename: header.Filename,
		FilePath: filePath,
		FileSize: header.Size,
		MimeType: mimeType,
	}

	if err := db.Create(&fileRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to save file record"})
		return
	}

	// 生成访问URL
	fileURL := fmt.Sprintf("/api/v1/files/download/%d", fileRecord.ID)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "File uploaded successfully",
		"data": gin.H{
			"id":        fileRecord.ID,
			"filename":  header.Filename,
			"size":      header.Size,
			"mime_type": mimeType,
			"url":       fileURL,
		},
	})
}

// List 文件列表
func List(c *gin.Context) {
	appID := c.Query("app_id")
	mimeType := c.Query("mime_type")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "app_id is required"})
		return
	}

	query := db.Model(&model.File{}).Where("app_id = ?", appID)

	if mimeType != "" {
		query = query.Where("mime_type LIKE ?", mimeType+"%")
	}

	var total int64
	query.Count(&total)

	var files []model.File
	offset := (page - 1) * size
	query.Offset(offset).Limit(size).Order("created_at DESC").Find(&files)

	// 添加URL
	var result []gin.H
	for _, f := range files {
		result = append(result, gin.H{
			"id":         f.ID,
			"filename":   f.Filename,
			"file_size":  f.FileSize,
			"mime_type":  f.MimeType,
			"url":        fmt.Sprintf("/api/v1/files/download/%d", f.ID),
			"created_at": f.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  result,
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}

// Detail 文件详情
func Detail(c *gin.Context) {
	id := c.Param("id")

	var file model.File
	if err := db.First(&file, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "File not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to query file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"id":         file.ID,
			"filename":   file.Filename,
			"file_size":  file.FileSize,
			"mime_type":  file.MimeType,
			"url":        fmt.Sprintf("/api/v1/files/download/%d", file.ID),
			"created_at": file.CreatedAt,
		},
	})
}

// Download 下载文件
func Download(c *gin.Context) {
	id := c.Param("id")

	var file model.File
	if err := db.First(&file, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "File not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to query file"})
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "File not found on disk"})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Filename))
	c.Header("Content-Type", file.MimeType)
	c.File(file.FilePath)
}

// Delete 删除文件
func Delete(c *gin.Context) {
	id := c.Param("id")

	var file model.File
	if err := db.First(&file, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "File not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to query file"})
		return
	}

	// 删除物理文件
	os.Remove(file.FilePath)

	// 删除数据库记录
	db.Delete(&file)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "File deleted successfully",
	})
}

// BatchDelete 批量删除文件
func BatchDelete(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	var files []model.File
	db.Find(&files, req.IDs)

	// 删除物理文件
	for _, file := range files {
		os.Remove(file.FilePath)
	}

	// 删除数据库记录
	result := db.Delete(&model.File{}, req.IDs)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Files deleted successfully",
		"data": gin.H{
			"affected": result.RowsAffected,
		},
	})
}

// Stats 文件统计
func Stats(c *gin.Context) {
	appID := c.Query("app_id")
	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "app_id is required"})
		return
	}

	var total int64
	var totalSize int64
	var todayCount int64

	db.Model(&model.File{}).Where("app_id = ?", appID).Count(&total)
	db.Model(&model.File{}).Where("app_id = ?", appID).Select("COALESCE(SUM(file_size), 0)").Scan(&totalSize)

	today := time.Now().Format("2006-01-02")
	db.Model(&model.File{}).Where("app_id = ? AND DATE(created_at) = ?", appID, today).Count(&todayCount)

	// 按类型统计
	var typeStats []struct {
		MimeType string `json:"mime_type"`
		Count    int64  `json:"count"`
		Size     int64  `json:"size"`
	}
	db.Model(&model.File{}).
		Select("SUBSTRING_INDEX(mime_type, '/', 1) as mime_type, COUNT(*) as count, SUM(file_size) as size").
		Where("app_id = ?", appID).
		Group("SUBSTRING_INDEX(mime_type, '/', 1)").
		Scan(&typeStats)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"total":       total,
			"total_size":  totalSize,
			"today_count": todayCount,
			"type_stats":  typeStats,
		},
	})
}
