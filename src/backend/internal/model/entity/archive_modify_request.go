package entity

import (
	"time"
)

// ArchiveModifyRequest 归档修改申请实体
type ArchiveModifyRequest struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	ArchiveID     uint       `gorm:"index;not null" json:"archive_id"`
	RequestType   string     `gorm:"type:varchar(20);not null" json:"request_type"`    // 修改类型：update/delete/hide
	RequestReason string     `gorm:"type:text" json:"request_reason"`                  // 修改原因
	RequestData   string     `gorm:"type:jsonb" json:"request_data"`                   // 修改数据（JSON格式）
	RequesterID   uint       `gorm:"index;not null" json:"requester_id"`               // 申请人ID
	Status        string     `gorm:"type:varchar(20);default:'pending'" json:"status"` // 状态：pending/approved/rejected
	ApproverID    uint       `gorm:"index" json:"approver_id"`                         // 审批人ID
	ApproveTime   *time.Time `json:"approve_time"`                                     // 审批时间
	ApproveReason string     `gorm:"type:text" json:"approve_reason"`                  // 审批意见
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	// 关联关系
	Archive   *Archive `gorm:"foreignKey:ArchiveID" json:"archive,omitempty"`
	Requester *User    `gorm:"foreignKey:RequesterID" json:"requester,omitempty"`
	Approver  *User    `gorm:"foreignKey:ApproverID" json:"approver,omitempty"`
}

// TableName 指定表名
func (ArchiveModifyRequest) TableName() string {
	return "archive_modify_requests"
}
