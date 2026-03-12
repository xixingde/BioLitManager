package repository

import (
	"errors"

	"biolitmanager/internal/model/entity"

	"gorm.io/gorm"
)

// ArchiveRepository 归档仓储接口
type ArchiveRepository struct {
	db *gorm.DB
}

// NewArchiveRepository 创建归档仓储实例
func NewArchiveRepository(db *gorm.DB) *ArchiveRepository {
	return &ArchiveRepository{
		db: db,
	}
}

// Create 创建归档记录
func (r *ArchiveRepository) Create(archive *entity.Archive) error {
	return r.db.Create(archive).Error
}

// FindByPaperID 根据论文ID查询该论文的归档记录
func (r *ArchiveRepository) FindByPaperID(paperID uint) (*entity.Archive, error) {
	var archive entity.Archive
	err := r.db.Where("paper_id = ?", paperID).First(&archive).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &archive, nil
}
