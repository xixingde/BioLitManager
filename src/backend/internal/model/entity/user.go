package entity

import (
	"time"

	"gorm.io/gorm"
)

// User 用户实体
type User struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Username       string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	PasswordHash   string         `gorm:"type:varchar(255);not null" json:"-"`
	Name           string         `gorm:"type:varchar(100)" json:"name"`
	Role           string         `gorm:"type:varchar(20);default:'user'" json:"role"`
	Department     string         `gorm:"type:varchar(100)" json:"department"`
	IDCard         string         `gorm:"type:varchar(18)" json:"id_card"`
	Phone          string         `gorm:"type:varchar(20)" json:"phone"`
	Email          string         `gorm:"type:varchar(100)" json:"email"`
	IsEmailNotify  bool           `gorm:"default:true" json:"is_email_notify"`
	IsLocked       bool           `gorm:"default:false" json:"is_locked"`
	LockUntil      *time.Time     `gorm:"index" json:"lock_until"`
	IsDisabled     bool           `gorm:"default:false" json:"is_disabled"`
	LoginFailCount int            `gorm:"default:0" json:"login_fail_count"`
	LastLoginAt    *time.Time     `json:"last_login_at"`
	LastLoginIP    string         `gorm:"type:varchar(50)" json:"last_login_ip"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
