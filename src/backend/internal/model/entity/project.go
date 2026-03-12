package entity

import (
	"time"

	"gorm.io/gorm"
)

// Project 课题实体
type Project struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Code        string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	ProjectType string         `gorm:"type:varchar(20);not null" json:"project_type"`
	Source      string         `gorm:"type:varchar(100)" json:"source"`
	Level       string         `gorm:"type:varchar(20)" json:"level"`
	Status      string         `gorm:"type:varchar(20);not null" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Papers []*Paper `gorm:"many2many:paper_projects;foreignKey:ID;references:ID" json:"papers,omitempty"`
}

// TableName 指定表名
func (Project) TableName() string {
	return "projects"
}
