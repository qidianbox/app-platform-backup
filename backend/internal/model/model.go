package model

import (
	"time"

	"gorm.io/gorm"
)

// Admin 管理员模型
type Admin struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Username  string         `gorm:"uniqueIndex;size:50" json:"username"`
	Password  string         `gorm:"size:255" json:"-"`
	Nickname  string         `gorm:"size:100" json:"nickname"`
	Avatar    string         `gorm:"size:255" json:"avatar"`
	Status    int            `gorm:"default:1" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// App 应用模型
type App struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:100" json:"name"`
	AppID       string         `gorm:"uniqueIndex;size:50" json:"app_id"`
	AppSecret   string         `gorm:"size:100" json:"app_secret"`
	PackageName string         `gorm:"size:100" json:"package_name"`
	Description string         `gorm:"type:text" json:"description"`
	Icon        string         `gorm:"size:255" json:"icon"`
	Status      int            `gorm:"default:1" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// ModuleTemplate 模块模板
type ModuleTemplate struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	ModuleCode   string         `gorm:"uniqueIndex;size:50" json:"module_code"`
	ModuleName   string         `gorm:"size:100" json:"module_name"`
	Category     string         `gorm:"size:50" json:"category"`
	Description  string         `gorm:"type:text" json:"description"`
	Icon         string         `gorm:"size:100" json:"icon"`
	ConfigSchema string         `gorm:"type:json" json:"config_schema"`
	Dependencies string         `gorm:"type:json" json:"dependencies"`
	SourceModule string         `gorm:"size:50" json:"source_module"`
	FunctionType string         `gorm:"size:20" json:"function_type"`
	Status       int            `gorm:"default:1" json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// AppModule APP启用的模块
type AppModule struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	AppID      uint           `gorm:"index" json:"app_id"`
	ModuleCode string         `gorm:"size:50" json:"module_code"`
	Config     string         `gorm:"type:json" json:"config"`
	Status     int            `gorm:"default:1" json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// ModuleConfigHistory 模块配置历史
type ModuleConfigHistory struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	AppID      uint      `gorm:"index" json:"app_id"`
	ModuleCode string    `gorm:"size:50" json:"module_code"`
	Config     string    `gorm:"type:json" json:"config"`
	Version    int       `json:"version"`
	Operator   string    `gorm:"size:50" json:"operator"`
	Remark     string    `gorm:"size:255" json:"remark"`
	CreatedAt  time.Time `json:"created_at"`
}
