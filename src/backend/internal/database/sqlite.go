package database

import (
	"database/sql"
	"fmt"
	"time"

	"biolitmanager/internal/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

var db *gorm.DB

// InitDB 初始化数据库连接
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s?_journal_mode=WAL&_timeout=5000&_cache_size=-2000", cfg.Database.Path)

	// 使用 modernc.org/sqlite 驱动（纯Go实现，不需要CGO）
	sqlDB, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db, err := gorm.Open(sqlite.Dialector{
		Conn: sqlDB,
	}, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	dbInstance, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	dbInstance.SetMaxIdleConns(10)
	dbInstance.SetMaxOpenConns(100)
	dbInstance.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return db
}

// SetDB 设置数据库实例
func SetDB(database *gorm.DB) {
	db = database
}
