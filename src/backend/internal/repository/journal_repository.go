package repository

import (
	"errors"

	"biolitmanager/internal/model/entity"

	"gorm.io/gorm"
)

// JournalRepository 期刊仓储接口
type JournalRepository struct {
	db *gorm.DB
}

// NewJournalRepository 创建期刊仓储实例
func NewJournalRepository(db *gorm.DB) *JournalRepository {
	return &JournalRepository{
		db: db,
	}
}

// Create 创建期刊
func (r *JournalRepository) Create(journal *entity.Journal) error {
	return r.db.Create(journal).Error
}

// Update 更新期刊
func (r *JournalRepository) Update(journal *entity.Journal) error {
	return r.db.Save(journal).Error
}

// FindByID 根据ID查询期刊
func (r *JournalRepository) FindByID(id uint) (*entity.Journal, error) {
	var journal entity.Journal
	err := r.db.First(&journal, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &journal, nil
}

// FindByName 根据期刊名称模糊查询
func (r *JournalRepository) FindByName(name string) ([]*entity.Journal, error) {
	var journals []*entity.Journal
	err := r.db.Where("full_name LIKE ? OR short_name LIKE ?", "%"+name+"%", "%"+name+"%").Find(&journals).Error
	if err != nil {
		return nil, err
	}
	return journals, nil
}

// FindByISSN 根据ISSN查询期刊
func (r *JournalRepository) FindByISSN(issn string) (*entity.Journal, error) {
	var journal entity.Journal
	err := r.db.Where("issn = ?", issn).First(&journal).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &journal, nil
}

// List 分页查询期刊列表
func (r *JournalRepository) List(page, size int) ([]*entity.Journal, int64, error) {
	var journals []*entity.Journal
	var total int64

	if err := r.db.Model(&entity.Journal{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := r.db.Offset(offset).Limit(size).Order("created_at DESC").Find(&journals).Error
	if err != nil {
		return nil, 0, err
	}

	return journals, total, nil
}
