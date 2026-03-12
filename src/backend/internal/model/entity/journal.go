package entity

import (
	"time"

	"gorm.io/gorm"
)

// Journal 期刊实体
type Journal struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	FullName     string         `gorm:"type:varchar(255);not null" json:"full_name"`
	ShortName    string         `gorm:"type:varchar(100)" json:"short_name"`
	ISSN         string         `gorm:"type:varchar(20);uniqueIndex;not null" json:"issn"`
	ImpactFactor float64        `json:"impact_factor"`
	Publisher    string         `gorm:"type:varchar(100)" json:"publisher"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Journal) TableName() string {
	return "journals"
}
