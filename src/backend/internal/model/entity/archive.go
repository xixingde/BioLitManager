package entity

import (
	"time"
)

// Archive 归档实体
type Archive struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	PaperID       uint      `gorm:"index;not null" json:"paper_id"`
	ArchiveNumber string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"archive_number"`
	ArchiveDate   time.Time `gorm:"not null" json:"archive_date"`
	ArchiverID    uint      `gorm:"index;not null" json:"archiver_id"`
	Status        string    `gorm:"type:varchar(20);default:'public'" json:"status"`
	IsHidden      bool      `gorm:"default:false" json:"is_hidden"`
	CreatedAt     time.Time `json:"created_at"`

	// 关联关系
	Paper    *Paper `gorm:"foreignKey:PaperID" json:"paper,omitempty"`
	Archiver *User  `gorm:"foreignKey:ArchiverID" json:"archiver,omitempty"`
}

// TableName 指定表名
func (Archive) TableName() string {
	return "archives"
}
