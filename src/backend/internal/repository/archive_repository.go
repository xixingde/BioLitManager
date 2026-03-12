package repository

import (
	"errors"

	"biolitmanager/internal/model/dto/request"
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

// FindByID 根据ID查询归档记录
func (r *ArchiveRepository) FindByID(id uint) (*entity.Archive, error) {
	var archive entity.Archive
	err := r.db.Preload("Paper").Preload("Archiver").
		Where("id = ?", id).First(&archive).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &archive, nil
}

// ListByYear 按年份查询归档记录
func (r *ArchiveRepository) ListByYear(year int, pagination *request.Pagination) ([]entity.Archive, int64, error) {
	var archives []entity.Archive
	var total int64

	query := r.db.Model(&entity.Archive{}).
		Joins("JOIN papers ON papers.id = archives.paper_id").
		Where("EXTRACT(YEAR FROM papers.publish_date) = ?", year)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (pagination.Page - 1) * pagination.PageSize
	if err := query.Offset(offset).Limit(pagination.PageSize).
		Preload("Paper").Preload("Archiver").
		Find(&archives).Error; err != nil {
		return nil, 0, err
	}

	return archives, total, nil
}

// ListByType 按收录类型查询归档记录
func (r *ArchiveRepository) ListByType(paperType string, pagination *request.Pagination) ([]entity.Archive, int64, error) {
	var archives []entity.Archive
	var total int64

	// 构建类型条件
	var typeCondition string
	switch paperType {
	case "SCI":
		typeCondition = "papers.is_sci = true"
	case "EI":
		typeCondition = "papers.is_ei = true"
	case "CI":
		typeCondition = "papers.is_ci = true"
	case "DI":
		typeCondition = "papers.is_di = true"
	case "CORE":
		typeCondition = "papers.is_core = true"
	default:
		typeCondition = "1=1"
	}

	query := r.db.Model(&entity.Archive{}).
		Joins("JOIN papers ON papers.id = archives.paper_id").
		Where(typeCondition)

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (pagination.Page - 1) * pagination.PageSize
	if err := query.Offset(offset).Limit(pagination.PageSize).
		Preload("Paper").Preload("Archiver").
		Find(&archives).Error; err != nil {
		return nil, 0, err
	}

	return archives, total, nil
}

// ListByAuthor 按作者查询归档记录
func (r *ArchiveRepository) ListByAuthor(authorName string, pagination *request.Pagination) ([]entity.Archive, int64, error) {
	var archives []entity.Archive
	var total int64

	query := r.db.Model(&entity.Archive{}).
		Joins("JOIN papers ON papers.id = archives.paper_id").
		Joins("JOIN authors ON authors.paper_id = papers.id").
		Where("authors.name LIKE ?", "%"+authorName+"%")

	// 获取总数
	if err := query.Distinct("archives.*").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (pagination.Page - 1) * pagination.PageSize
	if err := query.Distinct("archives.*").Offset(offset).Limit(pagination.PageSize).
		Preload("Paper").Preload("Archiver").
		Find(&archives).Error; err != nil {
		return nil, 0, err
	}

	return archives, total, nil
}

// ListByProject 按课题查询归档记录
func (r *ArchiveRepository) ListByProject(projectCode string, pagination *request.Pagination) ([]entity.Archive, int64, error) {
	var archives []entity.Archive
	var total int64

	query := r.db.Model(&entity.Archive{}).
		Joins("JOIN papers ON papers.id = archives.paper_id").
		Joins("JOIN paper_projects ON paper_projects.paper_id = papers.id").
		Joins("JOIN projects ON projects.id = paper_projects.project_id").
		Where("projects.code = ?", projectCode)

	// 获取总数
	if err := query.Distinct("archives.*").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (pagination.Page - 1) * pagination.PageSize
	if err := query.Distinct("archives.*").Offset(offset).Limit(pagination.PageSize).
		Preload("Paper").Preload("Archiver").
		Find(&archives).Error; err != nil {
		return nil, 0, err
	}

	return archives, total, nil
}

// ListArchived 查询已归档论文（支持分页）
func (r *ArchiveRepository) ListArchived(pagination *request.Pagination) ([]entity.Archive, int64, error) {
	var archives []entity.Archive
	var total int64

	query := r.db.Model(&entity.Archive{})

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (pagination.Page - 1) * pagination.PageSize
	if err := query.Offset(offset).Limit(pagination.PageSize).
		Preload("Paper").Preload("Archiver").
		Find(&archives).Error; err != nil {
		return nil, 0, err
	}

	return archives, total, nil
}

// UpdateStatus 更新归档状态（公开/隐藏）
func (r *ArchiveRepository) UpdateStatus(archiveID uint, isHidden bool) error {
	return r.db.Model(&entity.Archive{}).
		Where("id = ?", archiveID).
		Update("is_hidden", isHidden).Error
}

// ListAll 查询所有归档记录（支持分页和筛选）
func (r *ArchiveRepository) ListAll(req *request.SearchRequest, pagination *request.Pagination) ([]entity.Archive, int64, error) {
	var archives []entity.Archive
	var total int64

	query := r.db.Model(&entity.Archive{}).
		Joins("JOIN papers ON papers.id = archives.paper_id").
		Group("archives.id")

	// 按关键词筛选
	if req.Keywords != "" {
		query = query.Where("papers.title LIKE ? OR papers.abstract LIKE ?",
			"%"+req.Keywords+"%", "%"+req.Keywords+"%")
	}

	// 按标题筛选
	if req.Title != "" {
		query = query.Where("papers.title LIKE ?", "%"+req.Title+"%")
	}

	// 按DOI筛选
	if req.DOI != "" {
		query = query.Where("papers.doi = ?", req.DOI)
	}

	// 按期刊名称筛选
	if req.JournalName != "" {
		query = query.Where("papers.journal_name LIKE ?", "%"+req.JournalName+"%")
	}

	// 按作者姓名筛选
	if req.AuthorName != "" {
		query = query.Joins("JOIN authors ON authors.paper_id = papers.id").
			Where("authors.name LIKE ?", "%"+req.AuthorName+"%")
	}

	// 按课题编号筛选
	if req.ProjectCode != "" {
		query = query.Joins("JOIN paper_projects ON paper_projects.paper_id = papers.id").
			Joins("JOIN projects ON projects.id = paper_projects.project_id").
			Where("projects.code = ?", req.ProjectCode)
	}

	// 按收录类型筛选
	if req.PaperType != "" {
		switch req.PaperType {
		case "SCI":
			query = query.Where("papers.is_sci = true")
		case "EI":
			query = query.Where("papers.is_ei = true")
		case "CI":
			query = query.Where("papers.is_ci = true")
		case "DI":
			query = query.Where("papers.is_di = true")
		case "CORE":
			query = query.Where("papers.is_core = true")
		}
	}

	// 按年份筛选
	if req.Year > 0 {
		query = query.Where("EXTRACT(YEAR FROM papers.publish_date) = ?", req.Year)
	}

	// 按日期范围筛选
	if req.PublishDateStart != "" {
		query = query.Where("papers.publish_date >= ?", req.PublishDateStart)
	}
	if req.PublishDateEnd != "" {
		query = query.Where("papers.publish_date <= ?", req.PublishDateEnd)
	}

	// 获取总数（需要先清除Group）
	countQuery := r.db.Model(&entity.Archive{}).
		Joins("JOIN papers ON papers.id = archives.paper_id")

	// 重新应用筛选条件获取总数
	if req.Keywords != "" {
		countQuery = countQuery.Where("papers.title LIKE ? OR papers.abstract LIKE ?",
			"%"+req.Keywords+"%", "%"+req.Keywords+"%")
	}
	if req.Title != "" {
		countQuery = countQuery.Where("papers.title LIKE ?", "%"+req.Title+"%")
	}
	if req.DOI != "" {
		countQuery = countQuery.Where("papers.doi = ?", req.DOI)
	}
	if req.JournalName != "" {
		countQuery = countQuery.Where("papers.journal_name LIKE ?", "%"+req.JournalName+"%")
	}
	if req.AuthorName != "" {
		countQuery = countQuery.Joins("JOIN authors ON authors.paper_id = papers.id").
			Where("authors.name LIKE ?", "%"+req.AuthorName+"%")
	}
	if req.ProjectCode != "" {
		countQuery = countQuery.Joins("JOIN paper_projects ON paper_projects.paper_id = papers.id").
			Joins("JOIN projects ON projects.id = paper_projects.project_id").
			Where("projects.code = ?", req.ProjectCode)
	}
	if req.PaperType != "" {
		switch req.PaperType {
		case "SCI":
			countQuery = countQuery.Where("papers.is_sci = true")
		case "EI":
			countQuery = countQuery.Where("papers.is_ei = true")
		case "CI":
			countQuery = countQuery.Where("papers.is_ci = true")
		case "DI":
			countQuery = countQuery.Where("papers.is_di = true")
		case "CORE":
			countQuery = countQuery.Where("papers.is_core = true")
		}
	}
	if req.Year > 0 {
		countQuery = countQuery.Where("EXTRACT(YEAR FROM papers.publish_date) = ?", req.Year)
	}
	if req.PublishDateStart != "" {
		countQuery = countQuery.Where("papers.publish_date >= ?", req.PublishDateStart)
	}
	if req.PublishDateEnd != "" {
		countQuery = countQuery.Where("papers.publish_date <= ?", req.PublishDateEnd)
	}

	if err := countQuery.Distinct("archives.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (pagination.Page - 1) * pagination.PageSize
	if err := query.Offset(offset).Limit(pagination.PageSize).
		Preload("Paper").Preload("Archiver").
		Find(&archives).Error; err != nil {
		return nil, 0, err
	}

	return archives, total, nil
}

// CreateModifyRequest 创建归档修改申请
func (r *ArchiveRepository) CreateModifyRequest(req *entity.ArchiveModifyRequest) error {
	return r.db.Create(req).Error
}

// FindModifyRequestByID 根据ID查询归档修改申请
func (r *ArchiveRepository) FindModifyRequestByID(id uint) (*entity.ArchiveModifyRequest, error) {
	var req entity.ArchiveModifyRequest
	err := r.db.Preload("Archive").Preload("Requester").Preload("Approver").
		Where("id = ?", id).First(&req).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &req, nil
}

// ListModifyRequestsByArchiveID 查询归档的所有修改申请
func (r *ArchiveRepository) ListModifyRequestsByArchiveID(archiveID uint) ([]entity.ArchiveModifyRequest, error) {
	var requests []entity.ArchiveModifyRequest
	err := r.db.Preload("Archive").Preload("Requester").Preload("Approver").
		Where("archive_id = ?", archiveID).
		Order("created_at DESC").
		Find(&requests).Error
	return requests, err
}
