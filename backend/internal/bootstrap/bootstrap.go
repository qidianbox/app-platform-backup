// Package bootstrap 提供应用启动和初始化功能
package bootstrap

import (
	"log"

	"app-platform-backend/core/module"
	"app-platform-backend/internal/middleware"
	"app-platform-backend/internal/pkg/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Bootstrap 应用启动器
type Bootstrap struct {
	db     *gorm.DB
	router *gin.Engine
}

// New 创建一个新的启动器实例
func New(db *gorm.DB, router *gin.Engine) *Bootstrap {
	return &Bootstrap{
		db:     db,
		router: router,
	}
}

// InitModules 初始化所有已注册的模块
func (b *Bootstrap) InitModules() error {
	log.Println("[Bootstrap] Initializing modules...")

	// 1. 初始化所有模块
	if err := module.InitAllModules(); err != nil {
		return err
	}

	log.Printf("[Bootstrap] %d modules initialized", module.GetModuleCount())
	return nil
}

// SyncModulesToDB 将所有模块的功能同步到数据库
func (b *Bootstrap) SyncModulesToDB() error {
	log.Println("[Bootstrap] Syncing modules to database...")

	syncer := module.NewSyncer(b.db)
	if err := syncer.SyncModulesToDB(); err != nil {
		return err
	}

	log.Println("[Bootstrap] Module sync completed")
	return nil
}

// RegisterModuleRoutes 注册所有模块的路由
func (b *Bootstrap) RegisterModuleRoutes(authMiddleware gin.HandlerFunc) {
	log.Println("[Bootstrap] Registering module routes...")

	// 创建需要认证的路由组
	v1 := b.router.Group("/api/v1")
	authGroup := v1.Group("")
	authGroup.Use(authMiddleware)

	// 遍历所有模块，注册路由
	modules := module.GetAllModules()
	for _, m := range modules {
		meta := m.Meta()
		log.Printf("[Bootstrap] Registering routes for module: %s", meta.Code)
		m.RegisterRoutes(authGroup)
	}

	log.Printf("[Bootstrap] %d module routes registered", len(modules))
}

// Run 执行完整的启动流程
func (b *Bootstrap) Run() error {
	// 1. 初始化模块
	if err := b.InitModules(); err != nil {
		return err
	}

	// 2. 同步模块到数据库
	if err := b.SyncModulesToDB(); err != nil {
		return err
	}

	// 3. 注册模块路由
	b.RegisterModuleRoutes(middleware.AuthMiddleware())

	return nil
}

// RunWithDB 使用全局数据库连接执行启动流程
func RunWithDB(router *gin.Engine) error {
	db := database.GetDB()
	bootstrap := New(db, router)
	return bootstrap.Run()
}
