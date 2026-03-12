package service

import (
	"errors"
	"fmt"
	"time"

	"biolitmanager/internal/model/entity"
	"biolitmanager/internal/repository"
	"biolitmanager/pkg/logger"

	"go.uber.org/zap"
)

var (
	// ErrArchiveNotFound 归档记录不存在
	ErrArchiveNotFound = errors.New("归档记录不存在")
	// ErrPaperNotApproved 论文未通过审核，无法归档
	ErrPaperNotApproved = errors.New("论文未通过审核，无法归档")
)

// ArchiveService 归档服务
type ArchiveService struct {
	archiveRepo         *repository.ArchiveRepository
	paperRepo           *repository.PaperRepository
	operationLogService *OperationLogService
}

// NewArchiveService 创建归档服务实例
func NewArchiveService(
	archiveRepo *repository.ArchiveRepository,
	paperRepo *repository.PaperRepository,
	operationLogService *OperationLogService,
) *ArchiveService {
	return &ArchiveService{
		archiveRepo:         archiveRepo,
		paperRepo:           paperRepo,
		operationLogService: operationLogService,
	}
}

// ArchivePaper 归档论文
func (s *ArchiveService) ArchivePaper(paperID uint, archiverID uint, ipAddress string) (*entity.Archive, error) {
	// 查询论文
	paper, err := s.paperRepo.FindByID(paperID)
	if err != nil {
		logger.GetLogger().Error("Failed to find paper",
			zap.Uint("paper_id", paperID),
			zap.Error(err),
		)
		return nil, ErrSystemError
	}

	if paper == nil {
		return nil, ErrPaperNotFound
	}

	// 检查论文状态（只有审核通过的论文才能归档）
	if paper.Status != "审核通过" {
		return nil, ErrPaperNotApproved
	}

	// 检查是否已经归档
	existingArchive, err := s.archiveRepo.FindByPaperID(paperID)
	if err != nil {
		logger.GetLogger().Error("Failed to check archive",
			zap.Uint("paper_id", paperID),
			zap.Error(err),
		)
		return nil, ErrSystemError
	}

	if existingArchive != nil {
		return existingArchive, nil
	}

	// 生成归档编号
	archiveDate := time.Now()
	archiveNumber := s.generateArchiveNumber(archiveDate)

	// 创建归档记录
	archive := &entity.Archive{
		PaperID:       paperID,
		ArchiveNumber: archiveNumber,
		ArchiveDate:   archiveDate,
		ArchiverID:    archiverID,
	}

	if err := s.archiveRepo.Create(archive); err != nil {
		logger.GetLogger().Error("Failed to create archive",
			zap.Uint("paper_id", paperID),
			zap.Error(err),
		)
		return nil, ErrSystemError
	}

	// 记录操作日志
	s.operationLogService.LogOperation(archiverID, "archive", "paper", fmt.Sprintf("%d", paperID), fmt.Sprintf("归档论文 %s", paper.Title), "成功", ipAddress)

	logger.GetLogger().Info("Paper archived successfully",
		zap.Uint("paper_id", paperID),
		zap.String("archive_number", archiveNumber),
		zap.Uint("archiver_id", archiverID),
	)

	return archive, nil
}

// GetArchiveByPaperID 获取论文的归档记录
func (s *ArchiveService) GetArchiveByPaperID(paperID uint) (*entity.Archive, error) {
	archive, err := s.archiveRepo.FindByPaperID(paperID)
	if err != nil {
		logger.GetLogger().Error("Failed to find archive",
			zap.Uint("paper_id", paperID),
			zap.Error(err),
		)
		return nil, ErrSystemError
	}

	if archive == nil {
		return nil, ErrArchiveNotFound
	}

	return archive, nil
}

// generateArchiveNumber 生成归档编号
// 格式：ARCH-YYYYMMDD-XXXX
func (s *ArchiveService) generateArchiveNumber(archiveDate time.Time) string {
	// 格式：ARCH-20260311-0001
	dateStr := archiveDate.Format("20060102")
	sequence := "0001" // TODO: 实现根据当天的归档数量生成序号
	return fmt.Sprintf("ARCH-%s-%s", dateStr, sequence)
}
