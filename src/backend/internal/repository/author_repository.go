package repository

import (
	"errors"

	"biolitmanager/internal/model/entity"

	"gorm.io/gorm"
)

// AuthorRepository 作者仓储接口
type AuthorRepository struct {
	db *gorm.DB
}

// NewAuthorRepository 创建作者仓储实例
func NewAuthorRepository(db *gorm.DB) *AuthorRepository {
	return &AuthorRepository{
		db: db,
	}
}

// Create 创建作者
func (r *AuthorRepository) Create(author *entity.Author) error {
	return r.db.Create(author).Error
}

// Update 更新作者
func (r *AuthorRepository) Update(author *entity.Author) error {
	return r.db.Save(author).Error
}

// Delete 删除作者
func (r *AuthorRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Author{}, id).Error
}

// FindByID 根据ID查询作者
func (r *AuthorRepository) FindByID(id uint) (*entity.Author, error) {
	var author entity.Author
	err := r.db.First(&author, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &author, nil
}

// ListByPaperID 根据论文ID查询该论文的所有作者
func (r *AuthorRepository) ListByPaperID(paperID uint) ([]*entity.Author, error) {
	var authors []*entity.Author
	err := r.db.Where("paper_id = ?", paperID).Order("rank ASC").Find(&authors).Error
	if err != nil {
		return nil, err
	}
	return authors, nil
}

// DeleteByPaperID 根据论文ID删除该论文的所有作者
func (r *AuthorRepository) DeleteByPaperID(paperID uint) error {
	return r.db.Where("paper_id = ?", paperID).Delete(&entity.Author{}).Error
}

// CreateBatch 批量创建作者
func (r *AuthorRepository) CreateBatch(authors []*entity.Author) error {
	if len(authors) == 0 {
		return nil
	}
	return r.db.Create(&authors).Error
}
