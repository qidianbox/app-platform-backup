package version

import (
"app-platform-backend/core/module"
versionapi "app-platform-backend/internal/api/v1/version"
"github.com/gin-gonic/gin"
)

func init() { module.Register(&VersionModule{}) }
type VersionModule struct{}
func (m *VersionModule) Meta() module.Meta {
return module.Meta{Code: "version_management", Name: "版本管理", Description: "版本管理模块", Icon: "git-branch", SortOrder: 9}
}
func (m *VersionModule) GetFunctions() []module.Function {
return []module.Function{
{Code: "version_list", Name: "版本列表", Type: "passive", Description: "获取版本列表"},
{Code: "version_create", Name: "创建版本", Type: "active", Description: "创建新版本"},
{Code: "version_publish", Name: "发布版本", Type: "active", Description: "发布版本"},
{Code: "version_offline", Name: "下线版本", Type: "active", Description: "下线版本"},
{Code: "version_check", Name: "更新检查", Type: "passive", Description: "检查版本更新"},
}
}
func (m *VersionModule) RegisterRoutes(group *gin.RouterGroup) {
group.GET("/versions", versionapi.List)
group.POST("/versions", versionapi.Create)
group.POST("/versions/:id/publish", versionapi.Publish)
group.POST("/versions/:id/offline", versionapi.Offline)
group.GET("/versions/check", versionapi.CheckUpdate)
}
func (m *VersionModule) Init() error { return nil }
