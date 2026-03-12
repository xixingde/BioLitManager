package entity

import (
	"time"
)

// Attachment 附件实体
type Attachment struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PaperID    uint      `gorm:"index;not null" json:"paper_id"`
	FileType   string    `gorm:"type:varchar(20);not null" json:"file_type"`
	FileName   string    `gorm:"type:varchar(255);not null" json:"file_name"`
	FilePath   string    `gorm:"type:varchar(500);not null" json:"file_path"`
	FileSize   int64     `gorm:"not null" json:"file_size"`
	MimeType   string    `gorm:"type:varchar(100)" json:"mime_type"`
	UploaderID uint      `gorm:"index;not null" json:"uploader_id"`
	CreatedAt  time.Time `json:"created_at"`

	// 关联关系
	Paper    *Paper `gorm:"foreignKey:PaperID" json:"paper,omitempty"`
	Uploader *User  `gorm:"foreignKey:UploaderID" json:"uploader,omitempty"`
}

// TableName 指定表名
func (Attachment) TableName() string {
	return "attachments"
}
