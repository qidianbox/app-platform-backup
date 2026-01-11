package main

import (
	"fmt"
	"log"
	"time"

	// 核心模块
	"app-platform-backend/core/module"

	// 内部包
	"app-platform-backend/internal/api/v1/admin"
	"app-platform-backend/internal/api/v1/app"
	moduleapi "app-platform-backend/internal/api/v1/module"
	statsapi "app-platform-backend/internal/api/v1/stats"
	"app-platform-backend/internal/api/v1/system"
	wsapi "app-platform-backend/internal/api/v1/websocket"
	"app-platform-backend/internal/config"
	"app-platform-backend/internal/middleware"
	"app-platform-backend/internal/pkg/database"
	"app-platform-backend/internal/scheduler"

	// 导入所有功能模块（通过 import 的副作用触发模块注册）
	_ "app-platform-backend/modules"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig("./configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	if err := database.InitDB(&cfg.Database); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	defer database.Close()

	// 初始化JWT
	middleware.InitJWT(&cfg.JWT)

	// 初始化审计日志数据库连接
	middleware.InitAuditDB(database.GetDB())
	log.Println("[Main] Audit logging initialized")

	// 初始化并启动审计日志清理调度器
	auditCleanupScheduler := scheduler.InitAuditCleanupScheduler(database.GetDB(), scheduler.AuditCleanupConfig{
		RetentionDays: 90,  // 保留最近90天的日志
		CleanupHour:   3,   // 每天凌晨3点执行清理
		BatchSize:     1000, // 每批删除1000条
	})
	auditCleanupScheduler.Start()
	log.Println("[Main] Audit log cleanup scheduler started (retention: 90 days, cleanup at 03:00)")

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 创建路由
	r := gin.Default()

	// 中间件
	r.Use(middleware.CORSMiddleware(&cfg.CORS))
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.SecurityHeadersMiddleware()) // 添加HTTP安全响应头

	// 初始化全局限流器 (100 QPS/IP, 突发200请求)
	middleware.InitRateLimiter(200, 100)
	r.Use(middleware.GlobalRateLimitMiddleware())
	log.Println("[Main] Global rate limiter initialized (200 burst, 100 QPS/IP)")

	// ========================================
	// 模块化架构：初始化和同步
	// ========================================
	log.Println("[Main] Starting modular architecture initialization...")

	// 1. 初始化所有模块
	if err := module.InitAllModules(); err != nil {
		log.Fatalf("Failed to init modules: %v", err)
	}
	log.Printf("[Main] %d modules initialized", module.GetModuleCount())

	// 2. 同步模块功能到数据库
	syncer := module.NewSyncer(database.GetDB())
	if err := syncer.SyncModulesToDB(); err != nil {
		log.Fatalf("Failed to sync modules to database: %v", err)
	}

	// ========================================
	// API路由组
	// ========================================
	v1 := r.Group("/api/v1")
	{
// 公开接口（无需认证）
				// 登录接口使用更严格的限流 (5次/分钟/IP，防止暴力破解)
				v1.POST("/admin/login", middleware.APIRateLimitMiddleware(5, time.Minute), admin.Login)
				
				// 错误报告接口（限流30次/分钟/IP）
				v1.POST("/system/error-report", middleware.APIRateLimitMiddleware(30, time.Minute), system.ErrorReportHandler)
			
			// WebSocket连接端点（无需JWT认证，通过URL参数传递token）
			v1.GET("/ws", wsapi.HandleWebSocket)

		// 需要认证的接口
		auth := v1.Group("")
		auth.Use(middleware.AuthMiddleware())
		auth.Use(middleware.AuditMiddleware()) // 审计日志中间件
		{
			// 管理员相关
			adminGroup := auth.Group("/admin")
			{
				adminGroup.GET("/info", admin.GetInfo)
				adminGroup.POST("/logout", admin.Logout)
				adminGroup.PUT("/password", admin.UpdatePassword)
			}

			// 统计数据
			statsHandler := statsapi.NewStatsHandler(database.GetDB())
			auth.GET("/stats", statsHandler.GetStats)

			// APP管理
			appGroup := auth.Group("/apps")
			{
				appGroup.GET("", app.List)
				appGroup.POST("", app.Create)
				appGroup.GET("/:id", app.Detail)
				appGroup.PUT("/:id", app.Update)
				appGroup.DELETE("/:id", app.Delete)
				appGroup.POST("/:id/reset-secret", app.ResetSecret)

				// APP模块管理
				appGroup.GET("/:id/modules", moduleapi.GetAppModules)
				appGroup.GET("/:id/modules/:module_code", moduleapi.GetAppModule)
				appGroup.POST("/:id/modules", moduleapi.EnableModule)
				appGroup.PUT("/:id/modules/:module_code", moduleapi.UpdateModule)
				appGroup.DELETE("/:id/modules/:module_code", moduleapi.DisableModule)
				appGroup.POST("/:id/modules/batch", moduleapi.BatchEnableModules)

				// APP模块配置管理
				appGroup.PUT("/:id/modules/:module_code/config", moduleapi.SaveModuleConfig)
				appGroup.GET("/:id/modules/:module_code/config", moduleapi.GetModuleConfig)
				appGroup.DELETE("/:id/modules/:module_code/config", moduleapi.ResetModuleConfig)
				appGroup.POST("/:id/modules/:module_code/config/test", moduleapi.TestModuleConfig)
				// 配置历史
				appGroup.GET("/:id/modules/:module_code/config/history", moduleapi.GetConfigHistory)
				appGroup.POST("/:id/modules/:module_code/config/rollback/:history_id", moduleapi.RollbackConfig)
				appGroup.GET("/:id/modules/:module_code/config/compare", moduleapi.CompareConfig)

				// 模块依赖管理
				appGroup.GET("/:id/modules/:module_code/dependencies/check", moduleapi.CheckModuleDependencies)
				appGroup.GET("/:id/modules/:module_code/dependencies/reverse", moduleapi.CheckModuleReverseDependencies)
				appGroup.POST("/:id/modules/:module_code/dependencies/auto-enable", moduleapi.AutoEnableModuleDependencies)

				// 批量配置导入导出 (暂时禁用)
				// appGroup.POST("/:id/config/export", moduleapi.ExportConfig)
				// appGroup.POST("/:id/config/import", moduleapi.ImportConfig)
				// appGroup.POST("/:id/config/import/preview", moduleapi.PreviewImportConfig)
			}

			// ========================================
			// 模块化架构：动态注册模块路由
			// ========================================
			log.Println("[Main] Registering module routes...")
			modules := module.GetAllModules()
			for _, m := range modules {
				meta := m.Meta()
				log.Printf("[Main] Registering routes for module: %s (%s)", meta.Code, meta.Name)
				m.RegisterRoutes(auth)
			}
			log.Printf("[Main] %d module routes registered", len(modules))

			// 模块模板管理（核心功能，不通过模块注册）
			moduleGroup := auth.Group("/modules")
			{
				moduleGroup.GET("/templates", moduleapi.GetAllTemplates)
				moduleGroup.GET("/dependencies/detect/:module_code", moduleapi.DetectCircularDependency)
			}
		}
	}

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":       "ok",
			"modules":      module.GetModuleCount(),
			"architecture": "modular",
		})
	})

	// 模块信息接口（用于调试）
	r.GET("/api/v1/system/modules", func(c *gin.Context) {
		modules := module.GetAllModules()
		result := make([]gin.H, 0, len(modules))
		for _, m := range modules {
			meta := m.Meta()
			functions := m.GetFunctions()
			funcList := make([]gin.H, 0, len(functions))
			for _, fn := range functions {
				funcList = append(funcList, gin.H{
					"code":        fn.Code,
					"name":        fn.Name,
					"type":        fn.Type,
					"description": fn.Description,
				})
			}
			result = append(result, gin.H{
				"code":        meta.Code,
				"name":        meta.Name,
				"description": meta.Description,
				"icon":        meta.Icon,
				"functions":   funcList,
			})
		}
		c.JSON(200, gin.H{
			"total":   len(modules),
			"modules": result,
		})
	})

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
