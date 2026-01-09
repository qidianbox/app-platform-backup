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


// User 用户模型
type User struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	AppID       uint           `gorm:"index" json:"app_id"`
	OpenID      string         `gorm:"size:255;index" json:"open_id"`
	Nickname    string         `gorm:"size:255" json:"nickname"`
	Avatar      string         `gorm:"size:500" json:"avatar"`
	Phone       string         `gorm:"size:20" json:"phone"`
	Email       string         `gorm:"size:255" json:"email"`
	Status      int            `gorm:"default:1" json:"status"`
	LastLoginAt *time.Time     `json:"last_login_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// Message 消息模型
type Message struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	AppID     uint           `gorm:"index" json:"app_id"`
	UserID    *uint          `gorm:"index" json:"user_id"`
	Title     string         `gorm:"size:255" json:"title"`
	Content   string         `gorm:"type:text" json:"content"`
	Type      string         `gorm:"size:50;default:system" json:"type"`
	Status    int            `gorm:"default:0" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// PushRecord 推送记录模型
type PushRecord struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	AppID        uint           `gorm:"index" json:"app_id"`
	Title        string         `gorm:"size:255" json:"title"`
	Content      string         `gorm:"type:text" json:"content"`
	TargetType   string         `gorm:"size:50;default:all" json:"target_type"`
	TargetIDs    string         `gorm:"type:text" json:"target_ids"`
	Status       string         `gorm:"size:50;default:pending" json:"status"`
	SentCount    int            `gorm:"default:0" json:"sent_count"`
	SuccessCount int            `gorm:"default:0" json:"success_count"`
	FailedCount  int            `gorm:"default:0" json:"failed_count"`
	ScheduledAt  *time.Time     `json:"scheduled_at"`
	SentAt       *time.Time     `json:"sent_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// Event 事件模型
type Event struct {
	ID         uint64    `gorm:"primarykey" json:"id"`
	AppID      uint      `gorm:"index" json:"app_id"`
	UserID     *uint     `gorm:"index" json:"user_id"`
	EventCode  string    `gorm:"size:100;index" json:"event_code"`
	EventName  string    `gorm:"size:255" json:"event_name"`
	Properties string    `gorm:"type:json" json:"properties"`
	IP         string    `gorm:"size:50" json:"ip"`
	UserAgent  string    `gorm:"size:500" json:"user_agent"`
	CreatedAt  time.Time `gorm:"index" json:"created_at"`
}

// EventDefinition 事件定义模型
type EventDefinition struct {
	ID               uint           `gorm:"primarykey" json:"id"`
	AppID            uint           `gorm:"index" json:"app_id"`
	EventCode        string         `gorm:"size:100" json:"event_code"`
	EventName        string         `gorm:"size:255" json:"event_name"`
	Description      string         `gorm:"type:text" json:"description"`
	PropertiesSchema string         `gorm:"type:json" json:"properties_schema"`
	IsActive         int            `gorm:"default:1" json:"is_active"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// Log 日志模型
type Log struct {
	ID        uint64    `gorm:"primarykey" json:"id"`
	AppID     uint      `gorm:"index" json:"app_id"`
	Level     string    `gorm:"size:20;default:info;index" json:"level"`
	Module    string    `gorm:"size:100;index" json:"module"`
	Message   string    `gorm:"type:text" json:"message"`
	Context   string    `gorm:"type:json" json:"context"`
	IP        string    `gorm:"size:50" json:"ip"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}

// MonitorMetric 监控指标模型
type MonitorMetric struct {
	ID          uint64    `gorm:"primarykey" json:"id"`
	AppID       uint      `gorm:"index" json:"app_id"`
	MetricName  string    `gorm:"size:100;index" json:"metric_name"`
	MetricValue float64   `gorm:"type:decimal(20,4)" json:"metric_value"`
	Tags        string    `gorm:"type:json" json:"tags"`
	CreatedAt   time.Time `gorm:"index" json:"created_at"`
}

// MonitorAlert 告警模型
type MonitorAlert struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	AppID       uint           `gorm:"index" json:"app_id"`
	AlertName   string         `gorm:"size:255" json:"alert_name"`
	MetricName  string         `gorm:"size:100;index" json:"metric_name"`
	Condition   string         `gorm:"size:50" json:"condition"`
	Threshold   float64        `gorm:"type:decimal(20,4)" json:"threshold"`
	Status      string         `gorm:"size:50;default:normal" json:"status"`
	LastAlertAt *time.Time     `json:"last_alert_at"`
	IsActive    int            `gorm:"default:1" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// File 文件模型
type File struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	AppID     uint           `gorm:"index" json:"app_id"`
	Filename  string         `gorm:"size:255" json:"filename"`
	FilePath  string         `gorm:"size:500" json:"file_path"`
	FileSize  int64          `gorm:"default:0" json:"file_size"`
	MimeType  string         `gorm:"size:100" json:"mime_type"`
	UploadBy  *uint          `gorm:"index" json:"upload_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Config 配置模型
type Config struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	AppID       uint           `gorm:"index" json:"app_id"`
	ConfigKey   string         `gorm:"size:255" json:"config_key"`
	ConfigValue string         `gorm:"type:text" json:"config_value"`
	Description string         `gorm:"type:text" json:"description"`
	IsPublished int            `gorm:"default:0" json:"is_published"`
	PublishedAt *time.Time     `json:"published_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// ConfigHistory 配置历史模型
type ConfigHistory struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	ConfigID    uint      `gorm:"index" json:"config_id"`
	ConfigValue string    `gorm:"type:text" json:"config_value"`
	OperatorID  *uint     `json:"operator_id"`
	Operation   string    `gorm:"size:50" json:"operation"`
	CreatedAt   time.Time `gorm:"index" json:"created_at"`
}

// Version 版本模型
type Version struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	AppID         uint           `gorm:"index" json:"app_id"`
	VersionName   string         `gorm:"size:50" json:"version_name"`
	VersionCode   int            `gorm:"index" json:"version_code"`
	Description   string         `gorm:"type:text" json:"description"`
	DownloadURL   string         `gorm:"size:500" json:"download_url"`
	IsForceUpdate int            `gorm:"default:0" json:"is_force_update"`
	Status        string         `gorm:"size:50;default:draft" json:"status"`
	PublishedAt   *time.Time     `json:"published_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}
