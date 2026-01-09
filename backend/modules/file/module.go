package file

import (
	"app-platform-backend/core/module"
	fileapi "app-platform-backend/internal/api/v1/file"
	"app-platform-backend/internal/pkg/database"

	"github.com/gin-gonic/gin"
)

func init() { module.Register(&FileModule{}) }

type FileModule struct{}

func (m *FileModule) Meta() module.Meta {
	return module.Meta{Code: "file_storage", Name: "文件存储", Description: "文件存储模块", Icon: "folder", SortOrder: 7}
}

func (m *FileModule) GetFunctions() []module.Function {
	return []module.Function{
		{Code: "file_upload", Name: "文件上传", Type: "active", Description: "上传文件"},
		{Code: "file_download", Name: "文件下载", Type: "active", Description: "下载文件"},
		{Code: "file_list", Name: "文件列表", Type: "passive", Description: "获取文件列表"},
		{Code: "file_delete", Name: "文件删除", Type: "active", Description: "删除文件"},
		{Code: "file_stats", Name: "存储统计", Type: "passive", Description: "存储使用统计"},
	}
}

func (m *FileModule) RegisterRoutes(group *gin.RouterGroup) {
	fileapi.InitDB(database.GetDB())

	g := group.Group("/files")
	{
		g.GET("", fileapi.List)
		g.POST("", fileapi.Upload)
		g.GET("/stats", fileapi.Stats)
		g.GET("/:id", fileapi.Detail)
		g.GET("/download/:id", fileapi.Download)
		g.DELETE("/:id", fileapi.Delete)
		g.POST("/batch-delete", fileapi.BatchDelete)
	}
}

func (m *FileModule) Init() error { return nil }
