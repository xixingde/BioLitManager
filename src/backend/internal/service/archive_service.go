package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"biolitmanager/internal/model/dto/request"
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
	// ErrInvalidModifyRequestType 无效的修改申请类型
	ErrInvalidModifyRequestType = errors.New("无效的修改申请类型")
	// ErrModifyRequestNotFound 修改申请不存在
	ErrModifyRequestNotFound = errors.New("修改申请不存在")
	// ErrInvalidParameter 无效的参数
	ErrInvalidParameter = errors.New("无效的参数")
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
	archiveNumber := s.generateArchiveNumber(paperID)

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
// 格式：年份+论文ID+随机3位数字（例如：2026001001）
func (s *ArchiveService) generateArchiveNumber(paperID uint) string {
	// 格式：2026 + 001 + 001 = 2026001001
	currentYear := time.Now().Year()
	randomNum := rand.Intn(900) + 100 // 生成100-999的随机数
	return fmt.Sprintf("%d%03d%03d", currentYear, paperID, randomNum)
}

// GetArchivedPapers 获取已归档论文列表（分页）
func (s *ArchiveService) GetArchivedPapers(pagination *request.Pagination) ([]entity.Archive, int64, error) {
	archives, total, err := s.archiveRepo.ListArchived(pagination)
	if err != nil {
		logger.GetLogger().Error("Failed to get archived papers",
			zap.Error(err),
		)
		return nil, 0, ErrSystemError
	}

	logger.GetLogger().Info("Get archived papers successfully",
		zap.Int("count", len(archives)),
		zap.Int64("total", total),
	)

	return archives, total, nil
}

// GetArchivedPapersByYear 按年份获取归档论文
func (s *ArchiveService) GetArchivedPapersByYear(year int, pagination *request.Pagination) ([]entity.Archive, int64, error) {
	archives, total, err := s.archiveRepo.ListByYear(year, pagination)
	if err != nil {
		logger.GetLogger().Error("Failed to get archived papers by year",
			zap.Int("year", year),
			zap.Error(err),
		)
		return nil, 0, ErrSystemError
	}

	logger.GetLogger().Info("Get archived papers by year successfully",
		zap.Int("year", year),
		zap.Int("count", len(archives)),
		zap.Int64("total", total),
	)

	return archives, total, nil
}

// GetArchivedPapersByType 按收录类型获取归档论文（SCI/EI/CI/DI/CORE）
func (s *ArchiveService) GetArchivedPapersByType(paperType string, pagination *request.Pagination) ([]entity.Archive, int64, error) {
	// 验证类型参数
	validTypes := map[string]bool{"SCI": true, "EI": true, "CI": true, "DI": true, "CORE": true}
	if !validTypes[paperType] {
		logger.GetLogger().Warn("Invalid paper type",
			zap.String("paper_type", paperType),
		)
		return nil, 0, ErrInvalidParameter
	}

	archives, total, err := s.archiveRepo.ListByType(paperType, pagination)
	if err != nil {
		logger.GetLogger().Error("Failed to get archived papers by type",
			zap.String("paper_type", paperType),
			zap.Error(err),
		)
		return nil, 0, ErrSystemError
	}

	logger.GetLogger().Info("Get archived papers by type successfully",
		zap.String("paper_type", paperType),
		zap.Int("count", len(archives)),
		zap.Int64("total", total),
	)

	return archives, total, nil
}

// GetArchivedPapersByAuthor 按作者获取归档论文
func (s *ArchiveService) GetArchivedPapersByAuthor(authorName string, pagination *request.Pagination) ([]entity.Archive, int64, error) {
	if authorName == "" {
		logger.GetLogger().Warn("Author name is empty")
		return nil, 0, ErrInvalidParameter
	}

	archives, total, err := s.archiveRepo.ListByAuthor(authorName, pagination)
	if err != nil {
		logger.GetLogger().Error("Failed to get archived papers by author",
			zap.String("author_name", authorName),
			zap.Error(err),
		)
		return nil, 0, ErrSystemError
	}

	logger.GetLogger().Info("Get archived papers by author successfully",
		zap.String("author_name", authorName),
		zap.Int("count", len(archives)),
		zap.Int64("total", total),
	)

	return archives, total, nil
}

// GetArchivedPapersByProject 按课题获取归档论文
func (s *ArchiveService) GetArchivedPapersByProject(projectCode string, pagination *request.Pagination) ([]entity.Archive, int64, error) {
	if projectCode == "" {
		logger.GetLogger().Warn("Project code is empty")
		return nil, 0, ErrInvalidParameter
	}

	archives, total, err := s.archiveRepo.ListByProject(projectCode, pagination)
	if err != nil {
		logger.GetLogger().Error("Failed to get archived papers by project",
			zap.String("project_code", projectCode),
			zap.Error(err),
		)
		return nil, 0, ErrSystemError
	}

	logger.GetLogger().Info("Get archived papers by project successfully",
		zap.String("project_code", projectCode),
		zap.Int("count", len(archives)),
		zap.Int64("total", total),
	)

	return archives, total, nil
}

// HideArchive 隐藏归档论文
func (s *ArchiveService) HideArchive(archiveID uint, operatorID uint, ipAddress string) error {
	// 查询归档记录
	archive, err := s.archiveRepo.FindByID(archiveID)
	if err != nil {
		logger.GetLogger().Error("Failed to find archive",
			zap.Uint("archive_id", archiveID),
			zap.Error(err),
		)
		return ErrSystemError
	}

	if archive == nil {
		return ErrArchiveNotFound
	}

	// 更新隐藏状态
	if err := s.archiveRepo.UpdateStatus(archiveID, true); err != nil {
		logger.GetLogger().Error("Failed to hide archive",
			zap.Uint("archive_id", archiveID),
			zap.Error(err),
		)
		return ErrSystemError
	}

	// 记录操作日志
	s.operationLogService.LogOperation(operatorID, "hide", "archive", fmt.Sprintf("%d", archiveID), fmt.Sprintf("隐藏归档 %s", archive.ArchiveNumber), "成功", ipAddress)

	logger.GetLogger().Info("Archive hidden successfully",
		zap.Uint("archive_id", archiveID),
		zap.Uint("operator_id", operatorID),
	)

	return nil
}

// SubmitArchiveModifyRequest 提交归档修改申请
func (s *ArchiveService) SubmitArchiveModifyRequest(archiveID uint, requestType, requestReason string, requestData map[string]interface{}, requesterID uint, ipAddress string) (*entity.ArchiveModifyRequest, error) {
	// 验证修改类型
	validTypes := map[string]bool{"update": true, "delete": true, "hide": true}
	if !validTypes[requestType] {
		logger.GetLogger().Warn("Invalid modify request type",
			zap.String("request_type", requestType),
		)
		return nil, ErrInvalidModifyRequestType
	}

	// 查询归档记录
	archive, err := s.archiveRepo.FindByID(archiveID)
	if err != nil {
		logger.GetLogger().Error("Failed to find archive",
			zap.Uint("archive_id", archiveID),
			zap.Error(err),
		)
		return nil, ErrSystemError
	}

	if archive == nil {
		return nil, ErrArchiveNotFound
	}

	// 序列化请求数据
	requestDataJSON, err := json.Marshal(requestData)
	if err != nil {
		logger.GetLogger().Error("Failed to marshal request data",
			zap.Error(err),
		)
		return nil, ErrSystemError
	}

	// 创建修改申请
	modifyRequest := &entity.ArchiveModifyRequest{
		ArchiveID:     archiveID,
		RequestType:   requestType,
		RequestReason: requestReason,
		RequestData:   string(requestDataJSON),
		RequesterID:   requesterID,
		Status:        "pending",
	}

	if err := s.archiveRepo.CreateModifyRequest(modifyRequest); err != nil {
		logger.GetLogger().Error("Failed to create modify request",
			zap.Uint("archive_id", archiveID),
			zap.String("request_type", requestType),
			zap.Error(err),
		)
		return nil, ErrSystemError
	}

	// 记录操作日志
	s.operationLogService.LogOperation(requesterID, "submit", "archive_modify", fmt.Sprintf("%d", modifyRequest.ID), fmt.Sprintf("提交归档修改申请 %s-%s", archive.ArchiveNumber, requestType), "成功", ipAddress)

	logger.GetLogger().Info("Archive modify request submitted successfully",
		zap.Uint("archive_id", archiveID),
		zap.Uint("request_id", modifyRequest.ID),
		zap.String("request_type", requestType),
		zap.Uint("requester_id", requesterID),
	)

	return modifyRequest, nil
}
