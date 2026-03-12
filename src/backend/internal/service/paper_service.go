package service

import (
	stderrors "errors"
	"fmt"
	"regexp"
	"time"

	"biolitmanager/internal/model/entity"
	"biolitmanager/internal/repository"
	"biolitmanager/pkg/logger"

	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	// ErrPaperNotFound 论文不存在
	ErrPaperNotFound = stderrors.New("论文不存在")
	// ErrPaperDuplicate 论文重复
	ErrPaperDuplicate = stderrors.New("论文已存在重复")
	// ErrInvalidStatus 无效的状态
	ErrInvalidStatus = stderrors.New("无效的状态")
	// ErrPaperTitleRequired 标题为必填项
	ErrPaperTitleRequired = stderrors.New("标题为必填项")
	// ErrPaperAbstractRequired 摘要为必填项
	ErrPaperAbstractRequired = stderrors.New("摘要为必填项")
	// ErrPaperDOIInvalid DOI格式无效
	ErrPaperDOIInvalid = stderrors.New("DOI格式无效")
	// ErrPaperPublishDateInvalid 出版日期格式无效
	ErrPaperPublishDateInvalid = stderrors.New("出版日期格式无效")
	// ErrPaperImpactFactorInvalid 影响因子不能为负数
	ErrPaperImpactFactorInvalid = stderrors.New("影响因子不能为负数")
)

// PaperServiceInterface 论文服务接口
type PaperServiceInterface interface {
	CreatePaper(title, abstract string, journalID uint, doi string, impactFactor float64, publishDate *time.Time, submitterID uint, authors []*entity.Author, projectIDs []uint, operatorID uint, ipAddress string) (*entity.Paper, error)
	GetPaperByID(id uint) (*entity.Paper, error)
	ListPapers(page, size int) ([]*entity.Paper, int64, error)
	UpdatePaper(id uint, title, abstract string, journalID uint, doi string, impactFactor float64, publishDate *time.Time, authors []*entity.Author, projectIDs []uint, operatorID uint, ipAddress string) error
	DeletePaper(id, operatorID uint, ipAddress string) error
	SubmitForReview(id, operatorID uint, ipAddress string) error
	SaveDraft(id uint, title, abstract string, journalID uint, doi string, impactFactor float64, publishDate *time.Time, authors []*entity.Author, projectIDs []uint, operatorID uint, ipAddress string) error
	CheckDuplicate(title, doi string) ([]*entity.Paper, error)
	GetMyPapers(userID uint, page, size int) ([]*entity.Paper, int64, error)
	BatchImportPapers(file *excelize.File, submitterID uint) (int, int, []string)
}

// PaperService 论文服务
type PaperService struct {
	db                  *gorm.DB
	paperRepo           *repository.PaperRepository
	authorRepo          *repository.AuthorRepository
	attachmentRepo      *repository.AttachmentRepository
	paperProjectRepo    *repository.PaperProjectRepository
	operationLogService *OperationLogService
}

// NewPaperService 创建论文服务实例
func NewPaperService(
	db *gorm.DB,
	paperRepo *repository.PaperRepository,
	authorRepo *repository.AuthorRepository,
	attachmentRepo *repository.AttachmentRepository,
	paperProjectRepo *repository.PaperProjectRepository,
	operationLogService *OperationLogService,
) *PaperService {
	return &PaperService{
		db:                  db,
		paperRepo:           paperRepo,
		authorRepo:          authorRepo,
		attachmentRepo:      attachmentRepo,
		paperProjectRepo:    paperProjectRepo,
		operationLogService: operationLogService,
	}
}

// CreatePaper 创建论文（事务处理）
func (s *PaperService) CreatePaper(
	title, abstract string,
	journalID uint,
	doi string,
	impactFactor float64,
	publishDate *time.Time,
	submitterID uint,
	authors []*entity.Author,
	projectIDs []uint,
	operatorID uint,
	ipAddress string,
) (*entity.Paper, error) {
	// 检查论文是否重复
	duplicates, err := s.paperRepo.FindDuplicate(title, doi)
	if err != nil {
		logger.GetLogger().Error("Failed to check paper duplicate",
			zap.String("title", title),
			zap.String("doi", doi),
			zap.Error(err),
		)
		return nil, ErrSystemError
	}
	if len(duplicates) > 0 {
		logger.GetLogger().Warn("Paper already exists",
			zap.String("title", title),
			zap.String("doi", doi),
		)
		return nil, ErrPaperDuplicate
	}

	// 创建论文
	now := time.Now()
	paper := &entity.Paper{
		Title:        title,
		Abstract:     abstract,
		JournalID:    journalID,
		DOI:          doi,
		ImpactFactor: impactFactor,
		PublishDate:  publishDate,
		Status:       "draft",
		SubmitterID:  submitterID,
		SubmitTime:   now,
	}

	// 使用事务创建论文及其关联数据
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 创建论文
		if err := tx.Create(paper).Error; err != nil {
			logger.GetLogger().Error("Failed to create paper",
				zap.String("title", title),
				zap.Error(err),
			)
			return ErrSystemError
		}

		// 创建作者
		if len(authors) > 0 {
			for _, author := range authors {
				author.PaperID = paper.ID
			}
			if err := tx.Create(&authors).Error; err != nil {
				logger.GetLogger().Error("Failed to create authors",
					zap.Uint("paper_id", paper.ID),
					zap.Error(err),
				)
				return ErrSystemError
			}
		}

		// 创建课题关联
		if len(projectIDs) > 0 {
			paperProjects := make([]*entity.PaperProject, len(projectIDs))
			for i, projectID := range projectIDs {
				paperProjects[i] = &entity.PaperProject{
					PaperID:   paper.ID,
					ProjectID: projectID,
				}
			}
			if err := tx.Create(&paperProjects).Error; err != nil {
				logger.GetLogger().Error("Failed to create paper projects",
					zap.Uint("paper_id", paper.ID),
					zap.Error(err),
				)
				return ErrSystemError
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 记录操作日志
	s.operationLogService.LogOperation(operatorID, "create", "paper", fmt.Sprintf("%d", paper.ID), fmt.Sprintf("创建论文 %s", title), "成功", ipAddress)

	logger.GetLogger().Info("Paper created successfully",
		zap.Uint("paper_id", paper.ID),
		zap.String("title", title),
		zap.Uint("submitter_id", submitterID),
	)

	return paper, nil
}

// UpdatePaper 更新论文
func (s *PaperService) UpdatePaper(
	paperID uint,
	title, abstract string,
	journalID uint,
	doi string,
	impactFactor float64,
	publishDate *time.Time,
	authors []*entity.Author,
	projectIDs []uint,
	operatorID uint,
	ipAddress string,
) error {
	// 查询论文
	paper, err := s.paperRepo.FindByID(paperID)
	if err != nil {
		logger.GetLogger().Error("Failed to find paper",
			zap.Uint("paper_id", paperID),
			zap.Error(err),
		)
		return ErrSystemError
	}

	if paper == nil {
		return ErrPaperNotFound
	}

	// 检查论文是否重复（排除当前论文）
	duplicates, err := s.paperRepo.FindDuplicate(title, doi)
	if err != nil {
		logger.GetLogger().Error("Failed to check paper duplicate",
			zap.String("title", title),
			zap.String("doi", doi),
			zap.Error(err),
		)
		return ErrSystemError
	}
	for _, dup := range duplicates {
		if dup.ID != paperID {
			logger.GetLogger().Warn("Paper already exists",
				zap.String("title", title),
				zap.String("doi", doi),
			)
			return ErrPaperDuplicate
		}
	}

	// 使用事务更新论文及其关联数据
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 更新论文信息
		paper.Title = title
		paper.Abstract = abstract
		paper.JournalID = journalID
		paper.DOI = doi
		paper.ImpactFactor = impactFactor
		paper.PublishDate = publishDate

		if err := tx.Save(paper).Error; err != nil {
			logger.GetLogger().Error("Failed to update paper",
				zap.Uint("paper_id", paperID),
				zap.Error(err),
			)
			return ErrSystemError
		}

		// 删除原有作者
		if err := tx.Where("paper_id = ?", paperID).Delete(&entity.Author{}).Error; err != nil {
			logger.GetLogger().Error("Failed to delete authors",
				zap.Uint("paper_id", paperID),
				zap.Error(err),
			)
			return ErrSystemError
		}

		// 创建新作者
		if len(authors) > 0 {
			for _, author := range authors {
				author.PaperID = paperID
			}
			if err := tx.Create(&authors).Error; err != nil {
				logger.GetLogger().Error("Failed to create authors",
					zap.Uint("paper_id", paperID),
					zap.Error(err),
				)
				return ErrSystemError
			}
		}

		// 删除原有课题关联
		if err := tx.Where("paper_id = ?", paperID).Delete(&entity.PaperProject{}).Error; err != nil {
			logger.GetLogger().Error("Failed to delete paper projects",
				zap.Uint("paper_id", paperID),
				zap.Error(err),
			)
			return ErrSystemError
		}

		// 创建新课题关联
		if len(projectIDs) > 0 {
			paperProjects := make([]*entity.PaperProject, len(projectIDs))
			for i, projectID := range projectIDs {
				paperProjects[i] = &entity.PaperProject{
					PaperID:   paperID,
					ProjectID: projectID,
				}
			}
			if err := tx.Create(&paperProjects).Error; err != nil {
				logger.GetLogger().Error("Failed to create paper projects",
					zap.Uint("paper_id", paperID),
					zap.Error(err),
				)
				return ErrSystemError
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 记录操作日志
	s.operationLogService.LogOperation(operatorID, "update", "paper", fmt.Sprintf("%d", paperID), fmt.Sprintf("更新论文 %s", title), "成功", ipAddress)

	logger.GetLogger().Info("Paper updated successfully",
		zap.Uint("paper_id", paperID),
		zap.String("title", title),
		zap.Uint("operator_id", operatorID),
	)

	return nil
}

// GetPaperByID 获取论文详情（包含关联数据）
func (s *PaperService) GetPaperByID(paperID uint) (*entity.Paper, error) {
	var paper entity.Paper
	err := s.db.Preload("Journal").Preload("Submitter").
		Preload("Authors").
		Preload("Projects").
		Preload("Attachments").
		First(&paper, paperID).Error
	if err != nil {
		logger.GetLogger().Error("Failed to find paper",
			zap.Uint("paper_id", paperID),
			zap.Error(err),
		)
		return nil, ErrSystemError
	}

	return &paper, nil
}

// ListPapers 分页查询论文列表
func (s *PaperService) ListPapers(page, size int) ([]*entity.Paper, int64, error) {
	papers, total, err := s.paperRepo.List(page, size)
	if err != nil {
		logger.GetLogger().Error("Failed to list papers",
			zap.Int("page", page),
			zap.Int("size", size),
			zap.Error(err),
		)
		return nil, 0, ErrSystemError
	}

	return papers, total, nil
}

// SubmitForReview 提交审核
func (s *PaperService) SubmitForReview(paperID uint, operatorID uint, ipAddress string) error {
	// 查询论文
	paper, err := s.paperRepo.FindByID(paperID)
	if err != nil {
		logger.GetLogger().Error("Failed to find paper",
			zap.Uint("paper_id", paperID),
			zap.Error(err),
		)
		return ErrSystemError
	}

	if paper == nil {
		return ErrPaperNotFound
	}

	// 更新状态为"待业务审核"
	if err := s.paperRepo.UpdateStatus(paperID, "待业务审核"); err != nil {
		logger.GetLogger().Error("Failed to update paper status",
			zap.Uint("paper_id", paperID),
			zap.String("status", "待业务审核"),
			zap.Error(err),
		)
		return ErrSystemError
	}

	// 记录操作日志
	s.operationLogService.LogOperation(operatorID, "submit", "paper", fmt.Sprintf("%d", paperID), fmt.Sprintf("提交论文 %s 审核", paper.Title), "成功", ipAddress)

	logger.GetLogger().Info("Paper submitted for review",
		zap.Uint("paper_id", paperID),
		zap.String("title", paper.Title),
		zap.Uint("operator_id", operatorID),
	)

	return nil
}

// CheckDuplicate 查重
func (s *PaperService) CheckDuplicate(title, doi string) ([]*entity.Paper, error) {
	papers, err := s.paperRepo.FindDuplicate(title, doi)
	if err != nil {
		logger.GetLogger().Error("Failed to check paper duplicate",
			zap.String("title", title),
			zap.String("doi", doi),
			zap.Error(err),
		)
		return nil, ErrSystemError
	}

	return papers, nil
}

// DeletePaper 删除论文
func (s *PaperService) DeletePaper(paperID uint, operatorID uint, ipAddress string) error {
	// 查询论文
	paper, err := s.paperRepo.FindByID(paperID)
	if err != nil {
		logger.GetLogger().Error("Failed to find paper",
			zap.Uint("paper_id", paperID),
			zap.Error(err),
		)
		return ErrSystemError
	}

	if paper == nil {
		return ErrPaperNotFound
	}

	// 校验论文状态,仅草稿状态可删除
	if paper.Status != "draft" {
		logger.GetLogger().Warn("Paper status not allowed to delete",
			zap.Uint("paper_id", paperID),
			zap.String("status", paper.Status),
		)
		return ErrInvalidStatus
	}

	// 使用事务删除论文及其关联数据
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 删除作者
		if err := tx.Where("paper_id = ?", paperID).Delete(&entity.Author{}).Error; err != nil {
			logger.GetLogger().Error("Failed to delete authors",
				zap.Uint("paper_id", paperID),
				zap.Error(err),
			)
			return ErrSystemError
		}

		// 删除课题关联
		if err := tx.Where("paper_id = ?", paperID).Delete(&entity.PaperProject{}).Error; err != nil {
			logger.GetLogger().Error("Failed to delete paper projects",
				zap.Uint("paper_id", paperID),
				zap.Error(err),
			)
			return ErrSystemError
		}

		// 删除附件记录(物理文件删除由FileService处理)
		if err := tx.Where("paper_id = ?", paperID).Delete(&entity.Attachment{}).Error; err != nil {
			logger.GetLogger().Error("Failed to delete attachments",
				zap.Uint("paper_id", paperID),
				zap.Error(err),
			)
			return ErrSystemError
		}

		// 删除论文
		if err := tx.Delete(paper).Error; err != nil {
			logger.GetLogger().Error("Failed to delete paper",
				zap.Uint("paper_id", paperID),
				zap.Error(err),
			)
			return ErrSystemError
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 记录操作日志
	s.operationLogService.LogOperation(operatorID, "delete", "paper", fmt.Sprintf("%d", paperID), fmt.Sprintf("删除论文 %s", paper.Title), "成功", ipAddress)

	logger.GetLogger().Info("Paper deleted successfully",
		zap.Uint("paper_id", paperID),
		zap.String("title", paper.Title),
		zap.Uint("operator_id", operatorID),
	)

	return nil
}

// SaveDraft 保存草稿
func (s *PaperService) SaveDraft(
	paperID uint,
	title, abstract string,
	journalID uint,
	doi string,
	impactFactor float64,
	publishDate *time.Time,
	authors []*entity.Author,
	projectIDs []uint,
	operatorID uint,
	ipAddress string,
) error {
	// 查询论文
	paper, err := s.paperRepo.FindByID(paperID)
	if err != nil {
		logger.GetLogger().Error("Failed to find paper",
			zap.Uint("paper_id", paperID),
			zap.Error(err),
		)
		return ErrSystemError
	}

	if paper == nil {
		return ErrPaperNotFound
	}

	// 校验论文状态,仅草稿状态可保存
	if paper.Status != "draft" {
		logger.GetLogger().Warn("Paper status not allowed to save draft",
			zap.Uint("paper_id", paperID),
			zap.String("status", paper.Status),
		)
		return ErrInvalidStatus
	}

	// 使用事务更新论文及其关联数据(不进行严格校验和重复校验)
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 更新论文信息
		paper.Title = title
		paper.Abstract = abstract
		paper.JournalID = journalID
		paper.DOI = doi
		paper.ImpactFactor = impactFactor
		paper.PublishDate = publishDate

		if err := tx.Save(paper).Error; err != nil {
			logger.GetLogger().Error("Failed to save paper draft",
				zap.Uint("paper_id", paperID),
				zap.Error(err),
			)
			return ErrSystemError
		}

		// 删除原有作者
		if err := tx.Where("paper_id = ?", paperID).Delete(&entity.Author{}).Error; err != nil {
			logger.GetLogger().Error("Failed to delete authors",
				zap.Uint("paper_id", paperID),
				zap.Error(err),
			)
			return ErrSystemError
		}

		// 创建新作者
		if len(authors) > 0 {
			for _, author := range authors {
				author.PaperID = paperID
			}
			if err := tx.Create(&authors).Error; err != nil {
				logger.GetLogger().Error("Failed to create authors",
					zap.Uint("paper_id", paperID),
					zap.Error(err),
				)
				return ErrSystemError
			}
		}

		// 删除原有课题关联
		if err := tx.Where("paper_id = ?", paperID).Delete(&entity.PaperProject{}).Error; err != nil {
			logger.GetLogger().Error("Failed to delete paper projects",
				zap.Uint("paper_id", paperID),
				zap.Error(err),
			)
			return ErrSystemError
		}

		// 创建新课题关联
		if len(projectIDs) > 0 {
			paperProjects := make([]*entity.PaperProject, len(projectIDs))
			for i, projectID := range projectIDs {
				paperProjects[i] = &entity.PaperProject{
					PaperID:   paperID,
					ProjectID: projectID,
				}
			}
			if err := tx.Create(&paperProjects).Error; err != nil {
				logger.GetLogger().Error("Failed to create paper projects",
					zap.Uint("paper_id", paperID),
					zap.Error(err),
				)
				return ErrSystemError
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 记录操作日志
	s.operationLogService.LogOperation(operatorID, "save_draft", "paper", fmt.Sprintf("%d", paperID), fmt.Sprintf("保存论文草稿 %s", paper.Title), "成功", ipAddress)

	logger.GetLogger().Info("Paper draft saved successfully",
		zap.Uint("paper_id", paperID),
		zap.String("title", paper.Title),
		zap.Uint("operator_id", operatorID),
	)

	return nil
}

// GetMyPapers 获取我的论文列表
func (s *PaperService) GetMyPapers(submitterID uint, page, size int) ([]*entity.Paper, int64, error) {
	papers, total, err := s.paperRepo.ListBySubmitter(submitterID, page, size)
	if err != nil {
		logger.GetLogger().Error("Failed to list my papers",
			zap.Uint("submitter_id", submitterID),
			zap.Int("page", page),
			zap.Int("size", size),
			zap.Error(err),
		)
		return nil, 0, ErrSystemError
	}

	return papers, total, nil
}

// BatchImportPapers 批量导入论文
func (s *PaperService) BatchImportPapers(file *excelize.File, submitterID uint) (successCount, failedCount int, errors []string) {
	// 获取第一个工作表
	sheets := file.GetSheetList()
	if len(sheets) == 0 {
		logger.GetLogger().Error("Excel file has no sheets")
		errors = append(errors, "Excel文件没有工作表")
		return 0, 0, errors
	}

	sheetName := sheets[0]
	rows, err := file.GetRows(sheetName)
	if err != nil {
		logger.GetLogger().Error("Failed to get rows from Excel",
			zap.Error(err),
		)
		errors = append(errors, "读取Excel文件失败")
		return 0, 0, errors
	}

	// 跳过标题行,从第2行开始读取
	if len(rows) <= 1 {
		logger.GetLogger().Warn("Excel file has no data rows")
		return 0, 0, []string{"Excel文件没有数据行"}
	}

	// 逐行读取数据
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		rowNumber := i + 1

		// 校验列数(至少包含标题、摘要、期刊ID)
		if len(row) < 3 {
			logger.GetLogger().Warn("Row has insufficient columns",
				zap.Int("row", rowNumber),
				zap.Int("columns", len(row)),
			)
			errors = append(errors, fmt.Sprintf("第%d行:列数不足", rowNumber))
			failedCount++
			continue
		}

		// 解析期刊ID
		var journalID uint
		if _, err := fmt.Sscanf(row[2], "%d", &journalID); err != nil {
			logger.GetLogger().Warn("Invalid journal ID",
				zap.Int("row", rowNumber),
				zap.String("journal_id", row[2]),
			)
			errors = append(errors, fmt.Sprintf("第%d行:期刊ID格式错误", rowNumber))
			failedCount++
			continue
		}

		// 解析影响因子
		var impactFactor float64
		if len(row) > 4 && row[4] != "" {
			if _, err := fmt.Sscanf(row[4], "%f", &impactFactor); err != nil {
				logger.GetLogger().Warn("Invalid impact factor",
					zap.Int("row", rowNumber),
					zap.String("impact_factor", row[4]),
				)
				errors = append(errors, fmt.Sprintf("第%d行:影响因子格式错误", rowNumber))
				failedCount++
				continue
			}
		}

		// 解析出版日期
		var publishDate *time.Time
		if len(row) > 5 && row[5] != "" {
			if t, err := time.Parse("2006-01-02", row[5]); err == nil {
				publishDate = &t
			} else {
				logger.GetLogger().Warn("Invalid publish date",
					zap.Int("row", rowNumber),
					zap.String("publish_date", row[5]),
				)
				// 日期格式错误不是致命错误,继续处理
			}
		}

		// 获取DOI(第3列)
		doi := ""
		if len(row) > 3 {
			doi = row[3]
		}

		// 创建论文
		paper := &entity.Paper{
			Title:        row[0],
			Abstract:     row[1],
			JournalID:    journalID,
			DOI:          doi,
			ImpactFactor: impactFactor,
			PublishDate:  publishDate,
			Status:       "draft",
			SubmitterID:  submitterID,
			SubmitTime:   time.Now(),
		}

		if err := s.db.Create(paper).Error; err != nil {
			logger.GetLogger().Error("Failed to create paper",
				zap.Int("row", rowNumber),
				zap.String("title", row[0]),
				zap.Error(err),
			)
			errors = append(errors, fmt.Sprintf("第%d行:创建论文失败", rowNumber))
			failedCount++
			continue
		}

		successCount++

		logger.GetLogger().Info("Paper imported successfully",
			zap.Uint("paper_id", paper.ID),
			zap.String("title", row[0]),
			zap.Int("row", rowNumber),
		)
	}

	logger.GetLogger().Info("Batch import completed",
		zap.Int("success_count", successCount),
		zap.Int("failed_count", failedCount),
		zap.Int("total_errors", len(errors)),
	)

	return successCount, failedCount, errors
}

// ValidatePaperData 校验论文数据
func (s *PaperService) ValidatePaperData(title, abstract, doi, publishDate string, impactFactor float64) error {
	// 校验标题必填
	if title == "" {
		logger.GetLogger().Warn("Paper title is required")
		return ErrPaperTitleRequired
	}

	// 校验摘要必填
	if abstract == "" {
		logger.GetLogger().Warn("Paper abstract is required")
		return ErrPaperAbstractRequired
	}

	// 校验 DOI 格式（如果提供了）
	if doi != "" {
		doiPattern := regexp.MustCompile(`^10\.\d{4,9}/[-._;()/:A-Z0-9]+$`)
		if !doiPattern.MatchString(doi) {
			logger.GetLogger().Warn("Paper DOI format is invalid",
				zap.String("doi", doi),
			)
			return ErrPaperDOIInvalid
		}
	}

	// 校验出版日期格式（如果提供了）
	if publishDate != "" {
		if _, err := time.Parse("2006-01-02", publishDate); err != nil {
			logger.GetLogger().Warn("Paper publish date format is invalid",
				zap.String("publish_date", publishDate),
				zap.Error(err),
			)
			return ErrPaperPublishDateInvalid
		}
	}

	// 校验影响因子
	if impactFactor < 0 {
		logger.GetLogger().Warn("Paper impact factor cannot be negative",
			zap.Float64("impact_factor", impactFactor),
		)
		return ErrPaperImpactFactorInvalid
	}

	logger.GetLogger().Info("Paper data validation passed",
		zap.String("title", title),
		zap.String("doi", doi),
	)

	return nil
}
