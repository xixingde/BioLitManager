package repository

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"biolitmanager/internal/model/dto/request"
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

// FindByDOI 根据DOI精确查询论文
func (r *PaperRepository) FindByDOI(doi string) (*entity.Paper, error) {
	var paper entity.Paper
	err := r.db.Preload("Journal").Preload("Submitter").Where("doi = ?", doi).First(&paper).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &paper, nil
}

// FindByPubMedID 根据PubMedID精确查询论文
func (r *PaperRepository) FindByPubMedID(pubmedID string) (*entity.Paper, error) {
	var paper entity.Paper
	err := r.db.Preload("Journal").Preload("Submitter").Where("pubmed_id = ?", pubmedID).First(&paper).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &paper, nil
}

// FindByISSN 根据ISSN精确查询论文
func (r *PaperRepository) FindByISSN(issn string) ([]*entity.Paper, int64, error) {
	var papers []*entity.Paper
	var total int64

	if err := r.db.Model(&entity.Paper{}).Where("issn = ?", issn).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("Journal").Preload("Submitter").
		Where("issn = ?", issn).
		Order("publish_date DESC").
		Find(&papers).Error
	if err != nil {
		return nil, 0, err
	}

	return papers, total, nil
}

// FindByAuthorName 根据作者名查询论文
func (r *PaperRepository) FindByAuthorName(authorName string, page, size int) ([]*entity.Paper, int64, error) {
	var papers []*entity.Paper
	var total int64

	// 先查询包含该作者的论文ID
	subQuery := r.db.Model(&entity.Author{}).
		Select("paper_id").
		Where("name LIKE ?", "%"+authorName+"%")

	if err := r.db.Model(&entity.Paper{}).Where("id IN (?)", subQuery).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := r.db.Preload("Journal").Preload("Submitter").Preload("Authors").
		Where("id IN (?)", subQuery).
		Offset(offset).Limit(size).
		Order("publish_date DESC").
		Find(&papers).Error
	if err != nil {
		return nil, 0, err
	}

	return papers, total, nil
}

// FindByAuthorType 按作者类型查询论文
func (r *PaperRepository) FindByAuthorType(authorType request.AuthorType, page, size int) ([]*entity.Paper, int64, error) {
	var papers []*entity.Paper
	var total int64

	var db *gorm.DB
	switch authorType {
	case request.AuthorTypeFirst:
		db = r.db.Model(&entity.Paper{}).Where("is_first_author = ?", true)
	case request.AuthorTypeCoFirst:
		db = r.db.Model(&entity.Paper{}).Where("is_co_first_author = ?", true)
	case request.AuthorTypeCorresponding:
		db = r.db.Model(&entity.Paper{}).Where("is_corresponding_author = ?", true)
	default:
		return nil, 0, errors.New("invalid author type")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := r.db.Preload("Journal").Preload("Submitter").
		Offset(offset).Limit(size).
		Order("publish_date DESC").
		Find(&papers).Error
	if err != nil {
		return nil, 0, err
	}

	return papers, total, nil
}

// FindByProjectCode 根据课题编号查询论文
func (r *PaperRepository) FindByProjectCode(projectCode string, page, size int) ([]*entity.Paper, int64, error) {
	var papers []*entity.Paper
	var total int64

	// 通过课题编号关联查询论文
	subQuery := r.db.Model(&entity.Project{}).
		Select("id").
		Where("code = ?", projectCode)

	if err := r.db.Model(&entity.Paper{}).Where("id IN (SELECT paper_id FROM paper_projects WHERE project_id IN (?))", subQuery).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := r.db.Preload("Journal").Preload("Submitter").Preload("Projects").
		Where("id IN (SELECT paper_id FROM paper_projects WHERE project_id IN (?))", subQuery).
		Offset(offset).Limit(size).
		Order("publish_date DESC").
		Find(&papers).Error
	if err != nil {
		return nil, 0, err
	}

	return papers, total, nil
}

// ListByYear 按年份查询论文
func (r *PaperRepository) ListByYear(year int, page, size int) ([]*entity.Paper, int64, error) {
	var papers []*entity.Paper
	var total int64

	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year+1, 1, 1, 0, 0, 0, 0, time.UTC)

	if err := r.db.Model(&entity.Paper{}).Where("publish_date >= ? AND publish_date < ?", startDate, endDate).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := r.db.Preload("Journal").Preload("Submitter").
		Where("publish_date >= ? AND publish_date < ?", startDate, endDate).
		Offset(offset).Limit(size).
		Order("publish_date DESC").
		Find(&papers).Error
	if err != nil {
		return nil, 0, err
	}

	return papers, total, nil
}

// ListByType 按收录类型查询论文
func (r *PaperRepository) ListByType(paperType string, page, size int) ([]*entity.Paper, int64, error) {
	var papers []*entity.Paper
	var total int64

	var db *gorm.DB
	switch strings.ToUpper(paperType) {
	case "SCI":
		db = r.db.Model(&entity.Paper{}).Where("is_sci = ?", true)
	case "EI":
		db = r.db.Model(&entity.Paper{}).Where("is_ei = ?", true)
	case "CI":
		db = r.db.Model(&entity.Paper{}).Where("is_ci = ?", true)
	case "DI":
		db = r.db.Model(&entity.Paper{}).Where("is_di = ?", true)
	case "CORE":
		db = r.db.Model(&entity.Paper{}).Where("is_core = ?", true)
	default:
		return nil, 0, fmt.Errorf("invalid paper type: %s", paperType)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := r.db.Preload("Journal").Preload("Submitter").
		Offset(offset).Limit(size).
		Order("publish_date DESC").
		Find(&papers).Error
	if err != nil {
		return nil, 0, err
	}

	return papers, total, nil
}

// Count 统计论文总数
func (r *PaperRepository) Count() (int64, error) {
	var total int64
	err := r.db.Model(&entity.Paper{}).Count(&total).Error
	return total, err
}

// AdvancedSearch 高级搜索，支持多维度组合查询
func (r *PaperRepository) AdvancedSearch(req *request.SearchRequest) ([]*entity.Paper, int64, error) {
	var papers []*entity.Paper
	var total int64

	db := r.db.Model(&entity.Paper{})

	// 基础关键词搜索
	if req.Keywords != "" {
		keywords := "%" + req.Keywords + "%"
		db = db.Where("title LIKE ? OR abstract LIKE ?", keywords, keywords)
	}

	// 标题搜索
	if req.Title != "" {
		db = db.Where("title LIKE ?", "%"+req.Title+"%")
	}

	// 摘要搜索
	if req.Abstract != "" {
		db = db.Where("abstract LIKE ?", "%"+req.Abstract+"%")
	}

	// DOI精确查询
	if req.DOI != "" {
		db = db.Where("doi = ?", req.DOI)
	}

	// 期刊ID
	if req.JournalID > 0 {
		db = db.Where("journal_id = ?", req.JournalID)
	}

	// 期刊名称模糊查询
	if req.JournalName != "" {
		db = db.Where("journal_name LIKE ?", "%"+req.JournalName+"%")
	}

	// 作者名搜索
	if req.AuthorName != "" {
		subQuery := r.db.Model(&entity.Author{}).
			Select("paper_id").
			Where("name LIKE ?", "%"+req.AuthorName+"%")
		db = db.Where("id IN (?)", subQuery)
	}

	// 课题编号搜索
	if req.ProjectCode != "" {
		subQuery := r.db.Model(&entity.Project{}).
			Select("id").
			Where("code = ?", req.ProjectCode)
		db = db.Where("id IN (SELECT paper_id FROM paper_projects WHERE project_id IN (?))", subQuery)
	}

	// 发表日期范围
	if req.PublishDateStart != "" {
		startTime, err := time.Parse("2006-01-02", req.PublishDateStart)
		if err == nil {
			db = db.Where("publish_date >= ?", startTime)
		}
	}
	if req.PublishDateEnd != "" {
		endTime, err := time.Parse("2006-01-02", req.PublishDateEnd)
		if err == nil {
			db = db.Where("publish_date <= ?", endTime)
		}
	}

	// 影响因子范围
	if req.ImpactFactorMin > 0 {
		db = db.Where("impact_factor >= ?", req.ImpactFactorMin)
	}
	if req.ImpactFactorMax > 0 {
		db = db.Where("impact_factor <= ?", req.ImpactFactorMax)
	}

	// 状态筛选
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}

	// 提交人筛选（用于数据范围过滤）
	if req.SubmitterID != nil {
		db = db.Where("submitter_id = ?", *req.SubmitterID)
	}

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	page := req.Pagination.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.Pagination.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	db = db.Offset(offset).Limit(pageSize)

	// 排序
	sortField := "created_at"
	switch req.SortField {
	case request.SortFieldPublishDate:
		sortField = "publish_date"
	case request.SortFieldImpactFactor:
		sortField = "impact_factor"
	case request.SortFieldCitation:
		sortField = "citation_count"
	}
	sortOrder := "DESC"
	if req.SortOrder == request.SortOrderAsc {
		sortOrder = "ASC"
	}
	db = db.Order(sortField + " " + sortOrder)

	// 预加载关联数据
	db = db.Preload("Journal").Preload("Submitter").Preload("Authors").Preload("Projects")

	// 执行查询
	if err := db.Find(&papers).Error; err != nil {
		return nil, 0, err
	}

	return papers, total, nil
}
