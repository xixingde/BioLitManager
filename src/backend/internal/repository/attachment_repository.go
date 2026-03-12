package repository

import (
	"errors"

	"biolitmanager/internal/model/entity"

	"gorm.io/gorm"
)

// AttachmentRepository 附件仓储接口
type AttachmentRepository struct {
	db *gorm.DB
}

// NewAttachmentRepository 创建附件仓储实例
func NewAttachmentRepository(db *gorm.DB) *AttachmentRepository {
	return &AttachmentRepository{
		db: db,
	}
}

// Create 创建附件记录
func (r *AttachmentRepository) Create(attachment *entity.Attachment) error {
	return r.db.Create(attachment).Error
}

// FindByID 根据ID查询附件
func (r *AttachmentRepository) FindByID(id uint) (*entity.Attachment, error) {
	var attachment entity.Attachment
	err := r.db.First(&attachment, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &attachment, nil
}

// Delete 删除附件记录（仅数据库记录）
func (r *AttachmentRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Attachment{}, id).Error
}

// ListByPaperID 根据论文ID查询该论文的所有附件
func (r *AttachmentRepository) ListByPaperID(paperID uint) ([]*entity.Attachment, error) {
	var attachments []*entity.Attachment
	err := r.db.Where("paper_id = ?", paperID).Order("created_at DESC").Find(&attachments).Error
	if err != nil {
		return nil, err
	}
	return attachments, nil
}

// DeleteByPaperID 根据论文ID删除该论文的所有附件记录
func (r *AttachmentRepository) DeleteByPaperID(paperID uint) error {
	return r.db.Where("paper_id = ?", paperID).Delete(&entity.Attachment{}).Error
}
