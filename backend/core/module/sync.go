// Package module 提供模块同步功能
// 在应用启动时，将所有已注册模块的功能同步到数据库
package module

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

// ModuleTemplateRecord 对应数据库中的 module_templates 表
type ModuleTemplateRecord struct {
	ID           uint      `gorm:"primaryKey"`
	ModuleCode   string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	ModuleName   string    `gorm:"type:varchar(100);not null"`
	Description  string    `gorm:"type:text"`
	Dependencies string    `gorm:"type:json"`
	Icon         string    `gorm:"type:varchar(100)"`
	ConfigSchema string    `gorm:"type:text"`
	SortOrder    int       `gorm:"default:0"`
	IsActive     bool      `gorm:"default:true"`
	SourceModule string    `gorm:"type:varchar(50)"` // 新增：来源模块Code
	FunctionType string    `gorm:"type:varchar(20)"` // 新增：功能类型 active/passive
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// TableName 指定表名
func (ModuleTemplateRecord) TableName() string {
	return "module_templates"
}

// Syncer 模块同步器
type Syncer struct {
	db *gorm.DB
}

// NewSyncer 创建一个新的同步器
func NewSyncer(db *gorm.DB) *Syncer {
	return &Syncer{db: db}
}

// SyncModulesToDB 将所有已注册模块的功能同步到数据库
// 使用 UPSERT 策略：存在则更新，不存在则插入
func (s *Syncer) SyncModulesToDB() error {
	modules := GetAllModules()
	log.Printf("[ModuleSync] Starting sync, found %d registered modules", len(modules))

	for _, m := range modules {
		meta := m.Meta()
		functions := m.GetFunctions()

		log.Printf("[ModuleSync] Syncing module: %s (%s) with %d functions",
			meta.Code, meta.Name, len(functions))

		for _, fn := range functions {
			if err := s.syncFunction(meta, fn); err != nil {
				return fmt.Errorf("failed to sync function %s: %w", fn.Code, err)
			}
		}
	}

	log.Printf("[ModuleSync] Sync completed successfully")
	return nil
}

// syncFunction 同步单个功能到数据库
func (s *Syncer) syncFunction(meta Meta, fn Function) error {
	// 序列化配置Schema
	configSchemaJSON := "{}"
	if fn.ConfigSchema != nil {
		bytes, err := json.Marshal(fn.ConfigSchema)
		if err != nil {
			return fmt.Errorf("failed to marshal config schema: %w", err)
		}
		configSchemaJSON = string(bytes)
	}

	// 序列化依赖关系
	dependenciesJSON := "[]"
	if len(fn.Dependencies) > 0 {
		bytes, err := json.Marshal(fn.Dependencies)
		if err != nil {
			return fmt.Errorf("failed to marshal dependencies: %w", err)
		}
		dependenciesJSON = string(bytes)
	}

	// 构建记录
	record := ModuleTemplateRecord{
		ModuleCode:   fn.Code,
		ModuleName:   fn.Name,
		Description:  fn.Description,
		Dependencies: dependenciesJSON,
		Icon:         meta.Icon,
		ConfigSchema: configSchemaJSON,
		SortOrder:    fn.SortOrder,
		IsActive:     true,
		SourceModule: meta.Code,
		FunctionType: fn.Type,
		UpdatedAt:    time.Now(),
	}

	// 使用 UPSERT：根据 ModuleCode 判断是更新还是插入
	var existing ModuleTemplateRecord
	result := s.db.Where("module_code = ?", fn.Code).First(&existing)

	if result.Error == gorm.ErrRecordNotFound {
		// 不存在，插入新记录
		record.CreatedAt = time.Now()
		if err := s.db.Create(&record).Error; err != nil {
			return fmt.Errorf("failed to create record: %w", err)
		}
		log.Printf("[ModuleSync] Created new function: %s", fn.Code)
	} else if result.Error != nil {
		return fmt.Errorf("failed to query existing record: %w", result.Error)
	} else {
		// 存在，更新记录
		updates := map[string]interface{}{
			"module_name":   record.ModuleName,
			"description":   record.Description,
			"dependencies":  record.Dependencies,
			"icon":          record.Icon,
			"config_schema": record.ConfigSchema,
			"sort_order":    record.SortOrder,
			"source_module": record.SourceModule,
			"function_type": record.FunctionType,
			"updated_at":    record.UpdatedAt,
		}
		if err := s.db.Model(&existing).Updates(updates).Error; err != nil {
			return fmt.Errorf("failed to update record: %w", err)
		}
		log.Printf("[ModuleSync] Updated existing function: %s", fn.Code)
	}

	return nil
}

// GetSyncStats 获取同步统计信息
func (s *Syncer) GetSyncStats() (created int, updated int, total int) {
	modules := GetAllModules()
	for _, m := range modules {
		total += len(m.GetFunctions())
	}
	// 这里简化处理，实际可以在 syncFunction 中统计
	return 0, 0, total
}
