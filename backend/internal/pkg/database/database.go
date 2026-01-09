package database

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"app-platform-backend/internal/config"

	"github.com/go-sql-driver/mysql"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func InitDB(cfg *config.DatabaseConfig) error {
	var dsn string
	
	// 检查是否有DATABASE_URL环境变量（Manus平台数据库）
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" {
		log.Println("[Database] Using DATABASE_URL from environment")
		// 解析DATABASE_URL格式: mysql://user:pass@host:port/dbname?tls=...
		var err error
		dsn, err = parseDatabaseURL(databaseURL)
		if err != nil {
			return fmt.Errorf("failed to parse DATABASE_URL: %v", err)
		}
		log.Printf("[Database] DSN: %s", maskPassword(dsn))
	} else {
		log.Println("[Database] Using config file database settings")
		// 使用配置文件中的数据库配置
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	}

	var err error
	
	// 配置重试机制
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(gormMysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err == nil {
			// 测试连接
			sqlDB, _ := db.DB()
			if pingErr := sqlDB.Ping(); pingErr == nil {
				log.Println("[Database] Connected successfully")
				return nil
			}
		}
		log.Printf("[Database] Connection attempt %d failed: %v, retrying...", i+1, err)
		time.Sleep(time.Second * 2)
	}

	return err
}

// parseDatabaseURL 解析DATABASE_URL格式并返回DSN
func parseDatabaseURL(databaseURL string) (string, error) {
	// 注册TLS配置 - 使用skip-verify模式
	err := mysql.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true,
	})
	if err != nil {
		// 如果已经注册过，忽略错误
		log.Printf("[Database] TLS config registration: %v", err)
	}
	
	// 解析URL
	// 格式: mysql://user:pass@host:port/dbname?tls=true
	u, err := url.Parse(databaseURL)
	if err != nil {
		return "", fmt.Errorf("invalid DATABASE_URL: %v", err)
	}
	
	// 获取用户名和密码
	username := u.User.Username()
	password, _ := u.User.Password()
	
	// 获取主机和端口
	host := u.Host
	
	// 获取数据库名（去掉开头的/）
	dbName := strings.TrimPrefix(u.Path, "/")
	
	// 构建DSN
	// 格式: user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local&tls=tidb
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=tidb&timeout=30s&readTimeout=30s&writeTimeout=30s",
		username, password, host, dbName)
	
	return dsn, nil
}

// maskPassword 隐藏DSN中的密码用于日志输出
func maskPassword(dsn string) string {
	// 简单的密码隐藏
	colonIdx := strings.Index(dsn, ":")
	atIdx := strings.Index(dsn, "@")
	if colonIdx > 0 && atIdx > colonIdx {
		return dsn[:colonIdx+1] + "****" + dsn[atIdx:]
	}
	return dsn
}

func GetDB() *gorm.DB {
	return db
}

func Close() {
	if db != nil {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}
}
