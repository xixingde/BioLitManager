package repository

import (
	"biolitmanager/internal/model/entity"

	"gorm.io/gorm"
)

// OperationLogRepository 操作日志仓储接口
type OperationLogRepository struct {
	db *gorm.DB
}

// NewOperationLogRepository 创建操作日志仓储实例
func NewOperationLogRepository(db *gorm.DB) *OperationLogRepository {
	return &OperationLogRepository{
		db: db,
	}
}

// Create 创建操作日志
func (r *OperationLogRepository) Create(log *entity.OperationLog) error {
	return r.db.Create(log).Error
}

// List 分页查询操作日志列表
func (r *OperationLogRepository) List(page, size int) ([]*entity.OperationLog, int64, error) {
	var logs []*entity.OperationLog
	var total int64

	if err := r.db.Model(&entity.OperationLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := r.db.Offset(offset).Limit(size).Order("created_at DESC").Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
