package repository

import (
	"biolitmanager/internal/model/entity"

	"gorm.io/gorm"
)

// PaperProjectRepository 论文课题关联仓储接口
type PaperProjectRepository struct {
	db *gorm.DB
}

// NewPaperProjectRepository 创建论文课题关联仓储实例
func NewPaperProjectRepository(db *gorm.DB) *PaperProjectRepository {
	return &PaperProjectRepository{
		db: db,
	}
}

// Create 创建论文-课题关联记录
func (r *PaperProjectRepository) Create(paperProject *entity.PaperProject) error {
	return r.db.Create(paperProject).Error
}

// FindByPaperID 根据论文ID查询该论文关联的所有课题
func (r *PaperProjectRepository) FindByPaperID(paperID uint) ([]*entity.PaperProject, error) {
	var paperProjects []*entity.PaperProject
	err := r.db.Where("paper_id = ?", paperID).Find(&paperProjects).Error
	if err != nil {
		return nil, err
	}
	return paperProjects, nil
}

// DeleteByPaperID 根据论文ID删除该论文的所有课题关联
func (r *PaperProjectRepository) DeleteByPaperID(paperID uint) error {
	return r.db.Where("paper_id = ?", paperID).Delete(&entity.PaperProject{}).Error
}

// CreateBatch 批量创建论文-课题关联记录（使用事务）
func (r *PaperProjectRepository) CreateBatch(paperProjects []*entity.PaperProject) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, pp := range paperProjects {
			if err := tx.Create(pp).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
