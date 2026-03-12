package entity

import (
	"time"
)

// PaperProject 论文课题关联实体
type PaperProject struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PaperID   uint      `gorm:"index:idx_paper_project;not null" json:"paper_id"`
	ProjectID uint      `gorm:"index:idx_paper_project;not null" json:"project_id"`
	CreatedAt time.Time `json:"created_at"`

	Paper   *Paper   `gorm:"foreignKey:PaperID;references:ID" json:"paper,omitempty"`
	Project *Project `gorm:"foreignKey:ProjectID;references:ID" json:"project,omitempty"`
}

// TableName 指定表名
func (PaperProject) TableName() string {
	return "paper_projects"
}
