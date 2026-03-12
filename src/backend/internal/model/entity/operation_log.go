package entity

import (
	"time"

	"gorm.io/gorm"
)

// 操作类型常量
const (
	// 操作日志类型 - 基础操作
	OperationTypeLogin  = "login"
	OperationTypeLogout = "logout"
	OperationTypeCreate = "create"
	OperationTypeUpdate = "update"
	OperationTypeDelete = "delete"
	OperationTypeView   = "view"
	OperationTypeExport = "export"
	OperationTypeImport = "import"

	// 操作日志类型 - 审核操作
	OperationTypeReviewBusiness  = "review:business"
	OperationTypeReviewPolitical = "review:political"
	OperationTypeReviewSubmit    = "review:submit"
	OperationTypeReviewReject    = "review:reject"
	OperationTypeReviewApprove   = "review:approve"
)

// OperationLog 操作日志实体
type OperationLog struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	UserID           uint           `gorm:"index;not null" json:"user_id"`
	OperationType    string         `gorm:"type:varchar(50);not null" json:"operation_type"`
	Module           string         `gorm:"type:varchar(50);not null" json:"module"`
	TargetID         string         `gorm:"type:varchar(100)" json:"target_id"`
	OperationContent string         `gorm:"type:text" json:"operation_content"`
	OperationResult  string         `gorm:"type:varchar(20)" json:"operation_result"`
	IPAddress        string         `gorm:"type:varchar(50)" json:"ip_address"`
	CreatedAt        time.Time      `json:"created_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (OperationLog) TableName() string {
	return "operation_logs"
}
