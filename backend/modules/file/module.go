package file

import (
"app-platform-backend/core/module"
fileapi "app-platform-backend/internal/api/v1/file"
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
{Code: "file_list", Name: "文件列表", Type: "passive", Description: "获取文件列表"},
{Code: "file_delete", Name: "文件删除", Type: "active", Description: "删除文件"},
{Code: "file_stats", Name: "存储统计", Type: "passive", Description: "存储使用统计"},
}
}
func (m *FileModule) RegisterRoutes(group *gin.RouterGroup) {
group.POST("/files", fileapi.Upload)
group.GET("/files", fileapi.List)
group.DELETE("/files/:id", fileapi.Delete)
group.GET("/files/stats", fileapi.Stats)
}
func (m *FileModule) Init() error { return nil }
