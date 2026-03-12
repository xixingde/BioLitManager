package entity

import (
	"time"
)

// ReviewLog 审核记录实体
type ReviewLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PaperID    uint      `gorm:"index;not null" json:"paper_id"`
	ReviewType string    `gorm:"type:varchar(20);not null" json:"review_type"`
	Result     string    `gorm:"type:varchar(20);not null" json:"result"`
	Comment    string    `gorm:"type:text" json:"comment"`
	ReviewerID uint      `gorm:"index;not null" json:"reviewer_id"`
	ReviewTime time.Time `gorm:"not null" json:"review_time"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// 关联关系
	Paper    *Paper `gorm:"foreignKey:PaperID" json:"paper,omitempty"`
	Reviewer *User  `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
}

// TableName 指定表名
func (ReviewLog) TableName() string {
	return "review_logs"
}
