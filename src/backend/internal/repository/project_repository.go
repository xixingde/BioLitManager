package repository

import (
	"errors"

	"biolitmanager/internal/model/entity"

	"gorm.io/gorm"
)

// ProjectRepository 课题仓储接口
type ProjectRepository struct {
	db *gorm.DB
}

// NewProjectRepository 创建课题仓储实例
func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{
		db: db,
	}
}

// Create 创建课题
func (r *ProjectRepository) Create(project *entity.Project) error {
	return r.db.Create(project).Error
}

// Update 更新课题
func (r *ProjectRepository) Update(project *entity.Project) error {
	return r.db.Save(project).Error
}

// Delete 删除课题
func (r *ProjectRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Project{}, id).Error
}

// FindByID 根据ID查询课题
func (r *ProjectRepository) FindByID(id uint) (*entity.Project, error) {
	var project entity.Project
	err := r.db.First(&project, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

// FindByCode 根据课题编号查询课题
func (r *ProjectRepository) FindByCode(code string) (*entity.Project, error) {
	var project entity.Project
	err := r.db.Where("code = ?", code).First(&project).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

// List 分页查询课题列表
func (r *ProjectRepository) List(page, size int) ([]*entity.Project, int64, error) {
	var projects []*entity.Project
	var total int64

	if err := r.db.Model(&entity.Project{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := r.db.Offset(offset).Limit(size).Order("created_at DESC").Find(&projects).Error
	if err != nil {
		return nil, 0, err
	}

	return projects, total, nil
}

// CheckIsLinked 检查课题是否已关联论文
func (r *ProjectRepository) CheckIsLinked(projectID uint) (int64, error) {
	var count int64
	err := r.db.Model(&entity.PaperProject{}).Where("project_id = ?", projectID).Count(&count).Error
	return count, err
}
