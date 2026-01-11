package scheduler

import (
	"log"
	"sync"
	"time"

	"gorm.io/gorm"
)

// AuditCleanupConfig 审计日志清理配置
type AuditCleanupConfig struct {
	RetentionDays int           // 日志保留天数，默认90天
	CleanupHour   int           // 每天清理的小时（0-23），默认凌晨3点
	BatchSize     int           // 每批删除的记录数，默认1000
	Interval      time.Duration // 清理检查间隔，默认24小时
}

// DefaultAuditCleanupConfig 默认配置
var DefaultAuditCleanupConfig = AuditCleanupConfig{
	RetentionDays: 90,
	CleanupHour:   3,
	BatchSize:     1000,
	Interval:      24 * time.Hour,
}

// AuditCleanupScheduler 审计日志清理调度器
type AuditCleanupScheduler struct {
	db       *gorm.DB
	config   AuditCleanupConfig
	stopChan chan struct{}
	running  bool
	mu       sync.Mutex
}

// CleanupRecord 清理记录
type CleanupRecord struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CleanupTime time.Time `json:"cleanup_time"`
	DeletedRows int64     `json:"deleted_rows"`
	CutoffDate  time.Time `json:"cutoff_date"`
	Duration    int64     `json:"duration"` // 毫秒
	Status      string    `json:"status"`   // success, failed
	ErrorMsg    string    `json:"error_msg"`
	CreatedAt   time.Time `json:"created_at"`
}

var (
	scheduler *AuditCleanupScheduler
	once      sync.Once
)

// InitAuditCleanupScheduler 初始化审计日志清理调度器
func InitAuditCleanupScheduler(db *gorm.DB, config ...AuditCleanupConfig) *AuditCleanupScheduler {
	once.Do(func() {
		cfg := DefaultAuditCleanupConfig
		if len(config) > 0 {
			cfg = config[0]
		}

		scheduler = &AuditCleanupScheduler{
			db:       db,
			config:   cfg,
			stopChan: make(chan struct{}),
		}

		// 自动创建清理记录表
		db.AutoMigrate(&CleanupRecord{})

		log.Printf("[AuditCleanup] Scheduler initialized with config: RetentionDays=%d, CleanupHour=%d, BatchSize=%d",
			cfg.RetentionDays, cfg.CleanupHour, cfg.BatchSize)
	})

	return scheduler
}

// GetScheduler 获取调度器实例
func GetScheduler() *AuditCleanupScheduler {
	return scheduler
}

// Start 启动定时清理任务
func (s *AuditCleanupScheduler) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	go s.run()
	log.Printf("[AuditCleanup] Scheduler started")
}

// Stop 停止定时清理任务
func (s *AuditCleanupScheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	close(s.stopChan)
	s.running = false
	log.Printf("[AuditCleanup] Scheduler stopped")
}

// run 运行清理任务
func (s *AuditCleanupScheduler) run() {
	// 计算下一次清理时间
	nextCleanup := s.calculateNextCleanupTime()
	log.Printf("[AuditCleanup] Next cleanup scheduled at: %s", nextCleanup.Format("2006-01-02 15:04:05"))

	timer := time.NewTimer(time.Until(nextCleanup))
	defer timer.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-timer.C:
			// 执行清理
			s.executeCleanup()

			// 计算下一次清理时间
			nextCleanup = s.calculateNextCleanupTime()
			log.Printf("[AuditCleanup] Next cleanup scheduled at: %s", nextCleanup.Format("2006-01-02 15:04:05"))
			timer.Reset(time.Until(nextCleanup))
		}
	}
}

// calculateNextCleanupTime 计算下一次清理时间
func (s *AuditCleanupScheduler) calculateNextCleanupTime() time.Time {
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), s.config.CleanupHour, 0, 0, 0, now.Location())

	// 如果今天的清理时间已过，则安排到明天
	if next.Before(now) {
		next = next.Add(24 * time.Hour)
	}

	return next
}

// executeCleanup 执行清理操作
func (s *AuditCleanupScheduler) executeCleanup() {
	startTime := time.Now()
	cutoffDate := startTime.AddDate(0, 0, -s.config.RetentionDays)

	log.Printf("[AuditCleanup] Starting cleanup, deleting logs before: %s", cutoffDate.Format("2006-01-02"))

	var totalDeleted int64
	var lastErr error

	// 分批删除，避免长时间锁表
	for {
		result := s.db.Exec(
			"DELETE FROM audit_logs WHERE created_at < ? LIMIT ?",
			cutoffDate, s.config.BatchSize,
		)

		if result.Error != nil {
			lastErr = result.Error
			log.Printf("[AuditCleanup] Error during cleanup: %v", result.Error)
			break
		}

		totalDeleted += result.RowsAffected

		// 如果删除的行数小于批次大小，说明已经删除完毕
		if result.RowsAffected < int64(s.config.BatchSize) {
			break
		}

		// 短暂休眠，避免对数据库造成过大压力
		time.Sleep(100 * time.Millisecond)
	}

	duration := time.Since(startTime).Milliseconds()

	// 记录清理结果
	record := &CleanupRecord{
		CleanupTime: startTime,
		DeletedRows: totalDeleted,
		CutoffDate:  cutoffDate,
		Duration:    duration,
		Status:      "success",
		CreatedAt:   time.Now(),
	}

	if lastErr != nil {
		record.Status = "failed"
		record.ErrorMsg = lastErr.Error()
	}

	if err := s.db.Create(record).Error; err != nil {
		log.Printf("[AuditCleanup] Failed to save cleanup record: %v", err)
	}

	log.Printf("[AuditCleanup] Cleanup completed: deleted %d rows in %dms, status: %s",
		totalDeleted, duration, record.Status)
}

// ManualCleanup 手动执行清理
func (s *AuditCleanupScheduler) ManualCleanup(retentionDays int) (int64, error) {
	if retentionDays <= 0 {
		retentionDays = s.config.RetentionDays
	}

	startTime := time.Now()
	cutoffDate := startTime.AddDate(0, 0, -retentionDays)

	log.Printf("[AuditCleanup] Manual cleanup started, deleting logs before: %s", cutoffDate.Format("2006-01-02"))

	var totalDeleted int64

	// 分批删除
	for {
		result := s.db.Exec(
			"DELETE FROM audit_logs WHERE created_at < ? LIMIT ?",
			cutoffDate, s.config.BatchSize,
		)

		if result.Error != nil {
			return totalDeleted, result.Error
		}

		totalDeleted += result.RowsAffected

		if result.RowsAffected < int64(s.config.BatchSize) {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}

	duration := time.Since(startTime).Milliseconds()

	// 记录清理结果
	record := &CleanupRecord{
		CleanupTime: startTime,
		DeletedRows: totalDeleted,
		CutoffDate:  cutoffDate,
		Duration:    duration,
		Status:      "success",
		CreatedAt:   time.Now(),
	}
	s.db.Create(record)

	log.Printf("[AuditCleanup] Manual cleanup completed: deleted %d rows in %dms", totalDeleted, duration)

	return totalDeleted, nil
}

// GetCleanupHistory 获取清理历史记录
func (s *AuditCleanupScheduler) GetCleanupHistory(limit int) ([]CleanupRecord, error) {
	if limit <= 0 {
		limit = 20
	}

	var records []CleanupRecord
	err := s.db.Order("created_at DESC").Limit(limit).Find(&records).Error
	return records, err
}

// GetConfig 获取当前配置
func (s *AuditCleanupScheduler) GetConfig() AuditCleanupConfig {
	return s.config
}

// UpdateConfig 更新配置
func (s *AuditCleanupScheduler) UpdateConfig(config AuditCleanupConfig) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config = config
	log.Printf("[AuditCleanup] Config updated: RetentionDays=%d, CleanupHour=%d, BatchSize=%d",
		config.RetentionDays, config.CleanupHour, config.BatchSize)
}
