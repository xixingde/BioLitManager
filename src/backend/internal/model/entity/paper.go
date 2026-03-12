package entity

import (
	"time"

	"gorm.io/gorm"
)

// Paper 论文实体
type Paper struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Title        string         `gorm:"type:varchar(255);not null" json:"title"`
	Abstract     string         `gorm:"type:text" json:"abstract"`
	JournalID    uint           `gorm:"index;not null" json:"journal_id"`
	DOI          string         `gorm:"type:varchar(100)" json:"doi"`
	ImpactFactor float64        `json:"impact_factor"`
	PublishDate  *time.Time     `json:"publish_date"`
	Status       string         `gorm:"type:varchar(20);default:'draft'" json:"status"`
	SubmitterID  uint           `gorm:"index;not null" json:"submitter_id"`
	SubmitTime   time.Time      `gorm:"not null" json:"submit_time"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Journal     *Journal      `gorm:"foreignKey:JournalID" json:"journal,omitempty"`
	Submitter   *User         `gorm:"foreignKey:SubmitterID" json:"submitter,omitempty"`
	Authors     []*Author     `gorm:"foreignKey:PaperID" json:"authors,omitempty"`
	Projects    []*Project    `gorm:"foreignKey:PaperID;through:PaperProjects" json:"projects,omitempty"`
	Attachments []*Attachment `gorm:"foreignKey:PaperID" json:"attachments,omitempty"`
	ReviewLogs  []*ReviewLog  `gorm:"foreignKey:PaperID" json:"review_logs,omitempty"`
}

// TableName 指定表名
func (Paper) TableName() string {
	return "papers"
}
