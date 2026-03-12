package entity

import (
	"time"

	"gorm.io/gorm"
)

// Author 作者实体
type Author struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	PaperID    uint           `gorm:"index;not null" json:"paper_id"`
	Name       string         `gorm:"type:varchar(100);not null" json:"name"`
	AuthorType string         `gorm:"type:varchar(20);not null" json:"author_type"`
	Rank       int            `gorm:"not null" json:"rank"`
	Department string         `gorm:"type:varchar(100)" json:"department"`
	UserID     *uint          `gorm:"index" json:"user_id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Paper *Paper `gorm:"foreignKey:PaperID" json:"paper,omitempty"`
	User  *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (Author) TableName() string {
	return "authors"
}
