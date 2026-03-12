package repository

import (
	"errors"

	"biolitmanager/internal/model/entity"

	"gorm.io/gorm"
)

// PaperRepository 论文仓储接口
type PaperRepository struct {
	db *gorm.DB
}

// NewPaperRepository 创建论文仓储实例
func NewPaperRepository(db *gorm.DB) *PaperRepository {
	return &PaperRepository{
		db: db,
	}
}

// Create 创建论文
func (r *PaperRepository) Create(paper *entity.Paper) error {
	return r.db.Create(paper).Error
}

// Update 更新论文
func (r *PaperRepository) Update(paper *entity.Paper) error {
	return r.db.Save(paper).Error
}

// Delete 删除论文
func (r *PaperRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Paper{}, id).Error
}

// FindByID 根据ID查询论文（包含关联数据）
func (r *PaperRepository) FindByID(id uint) (*entity.Paper, error) {
	var paper entity.Paper
	err := r.db.Preload("Journal").Preload("Submitter").First(&paper, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &paper, nil
}

// List 分页查询论文列表
func (r *PaperRepository) List(page, size int) ([]*entity.Paper, int64, error) {
	var papers []*entity.Paper
	var total int64

	if err := r.db.Model(&entity.Paper{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := r.db.Preload("Journal").Preload("Submitter").
		Offset(offset).Limit(size).
		Order("created_at DESC").
		Find(&papers).Error
	if err != nil {
		return nil, 0, err
	}

	return papers, total, nil
}

// ListByStatus 根据状态分页查询论文列表
func (r *PaperRepository) ListByStatus(status string, page, size int) ([]*entity.Paper, int64, error) {
	var papers []*entity.Paper
	var total int64

	if err := r.db.Model(&entity.Paper{}).Where("status = ?", status).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := r.db.Preload("Journal").Preload("Submitter").
		Where("status = ?", status).
		Offset(offset).Limit(size).
		Order("created_at DESC").
		Find(&papers).Error
	if err != nil {
		return nil, 0, err
	}

	return papers, total, nil
}

// FindDuplicate 根据标题或DOI查找重复论文
func (r *PaperRepository) FindDuplicate(title, doi string) ([]*entity.Paper, error) {
	var papers []*entity.Paper
	query := r.db.Where("title = ?", title)
	if doi != "" {
		query = query.Or("doi = ?", doi)
	}
	err := query.Find(&papers).Error
	if err != nil {
		return nil, err
	}
	return papers, nil
}

// ListBySubmitter 根据提交人分页查询论文列表
func (r *PaperRepository) ListBySubmitter(submitterID uint, page, size int) ([]*entity.Paper, int64, error) {
	var papers []*entity.Paper
	var total int64

	if err := r.db.Model(&entity.Paper{}).Where("submitter_id = ?", submitterID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := r.db.Preload("Journal").Preload("Submitter").
		Where("submitter_id = ?", submitterID).
		Offset(offset).Limit(size).
		Order("created_at DESC").
		Find(&papers).Error
	if err != nil {
		return nil, 0, err
	}

	return papers, total, nil
}

// UpdateStatus 更新论文状态
func (r *PaperRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&entity.Paper{}).Where("id = ?", id).Update("status", status).Error
}
