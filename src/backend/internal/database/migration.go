package database

import (
	"biolitmanager/internal/model/entity"

	"gorm.io/gorm"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&entity.User{},
		&entity.OperationLog{},
		&entity.Paper{},
		&entity.Author{},
		&entity.Project{},
		&entity.Journal{},
		&entity.PaperProject{},
		&entity.ReviewLog{},
		&entity.Archive{},
		&entity.Attachment{},
	)
	if err != nil {
		return err
	}
	return nil
}
