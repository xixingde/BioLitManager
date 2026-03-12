package service

import (
	"errors"

	"biolitmanager/internal/model/entity"
	"biolitmanager/internal/repository"
	"biolitmanager/pkg/logger"

	"go.uber.org/zap"
)

var (
	// ErrJournalNotFound 期刊不存在
	ErrJournalNotFound = errors.New("期刊不存在")
	// ErrISSNExists ISSN已存在
	ErrISSNExists = errors.New("ISSN已存在")
)

// JournalService 期刊服务
type JournalService struct {
	journalRepo         *repository.JournalRepository
	operationLogService *OperationLogService
}

// NewJournalService 创建期刊服务实例
func NewJournalService(
	journalRepo *repository.JournalRepository,
	operationLogService *OperationLogService,
) *JournalService {
	return &JournalService{
		journalRepo:         journalRepo,
		operationLogService: operationLogService,
	}
}

// CreateJournal 创建期刊
func (s *JournalService) CreateJournal(
	fullName string,
	shortName string,
	issn string,
	impactFactor float64,
	publisher string,
	operatorID uint,
	ipAddress string,
) (*entity.Journal, error) {
	// 校验ISSN唯一性
	existingJournal, err := s.journalRepo.FindByISSN(issn)
	if err != nil {
		logger.GetLogger().Error("Failed to check journal ISSN",
			zap.String("issn", issn),
			zap.Error(err),
		)
		return nil, ErrSystemError
	}

	if existingJournal != nil {
		return nil, ErrISSNExists
	}

	journal := &entity.Journal{
		FullName:     fullName,
		ShortName:    shortName,
		ISSN:         issn,
		ImpactFactor: impactFactor,
		Publisher:    publisher,
	}

	if err := s.journalRepo.Create(journal); err != nil {
		logger.GetLogger().Error("Failed to create journal",
			zap.String("full_name", fullName),
			zap.String("issn", issn),
			zap.Error(err),
		)
		return nil, ErrSystemError
	}

	// 记录操作日志
	s.operationLogService.LogOperation(operatorID, "create", "journal", issn, "创建期刊 "+fullName, "成功", ipAddress)

	logger.GetLogger().Info("Journal created successfully",
		zap.Uint("journal_id", journal.ID),
		zap.String("full_name", fullName),
		zap.String("issn", issn),
	)

	return journal, nil
}

// UpdateJournal 更新期刊
func (s *JournalService) UpdateJournal(
	id uint,
	fullName string,
	shortName string,
	issn string,
	impactFactor float64,
	publisher string,
	operatorID uint,
	ipAddress string,
) error {
	// 查询期刊
	journal, err := s.journalRepo.FindByID(id)
	if err != nil {
		logger.GetLogger().Error("Failed to find journal",
			zap.Uint("journal_id", id),
			zap.Error(err),
		)
		return ErrSystemError
	}

	if journal == nil {
		return ErrJournalNotFound
	}

	// 校验ISSN唯一性（排除当前期刊）
	existingJournal, err := s.journalRepo.FindByISSN(issn)
	if err != nil {
		logger.GetLogger().Error("Failed to check journal ISSN",
			zap.String("issn", issn),
			zap.Error(err),
		)
		return ErrSystemError
	}

	if existingJournal != nil && existingJournal.ID != id {
		return ErrISSNExists
	}

	// 更新期刊信息
	journal.FullName = fullName
	journal.ShortName = shortName
	journal.ISSN = issn
	journal.ImpactFactor = impactFactor
	journal.Publisher = publisher

	if err := s.journalRepo.Update(journal); err != nil {
		logger.GetLogger().Error("Failed to update journal",
			zap.Uint("journal_id", id),
			zap.Error(err),
		)
		return ErrSystemError
	}

	// 记录操作日志
	s.operationLogService.LogOperation(operatorID, "update", "journal", issn, "更新期刊 "+fullName, "成功", ipAddress)

	logger.GetLogger().Info("Journal updated successfully",
		zap.Uint("journal_id", id),
		zap.String("full_name", fullName),
	)

	return nil
}

// GetJournalByID 获取期刊详情
func (s *JournalService) GetJournalByID(id uint) (*entity.Journal, error) {
	journal, err := s.journalRepo.FindByID(id)
	if err != nil {
		logger.GetLogger().Error("Failed to find journal",
			zap.Uint("journal_id", id),
			zap.Error(err),
		)
		return nil, ErrSystemError
	}

	if journal == nil {
		return nil, ErrJournalNotFound
	}

	return journal, nil
}

// SearchJournals 搜索期刊（按名称或ISSN模糊匹配）
func (s *JournalService) SearchJournals(keyword string) ([]*entity.Journal, error) {
	journals, err := s.journalRepo.FindByName(keyword)
	if err != nil {
		logger.GetLogger().Error("Failed to search journals",
			zap.String("keyword", keyword),
			zap.Error(err),
		)
		return nil, ErrSystemError
	}

	return journals, nil
}

// ListJournals 分页查询期刊列表
func (s *JournalService) ListJournals(page, size int) ([]*entity.Journal, int64, error) {
	journals, total, err := s.journalRepo.List(page, size)
	if err != nil {
		logger.GetLogger().Error("Failed to list journals",
			zap.Int("page", page),
			zap.Int("size", size),
			zap.Error(err),
		)
		return nil, 0, ErrSystemError
	}

	return journals, total, nil
}

// UpdateImpactFactor 更新期刊影响因子
func (s *JournalService) UpdateImpactFactor(id uint, impactFactor float64, operatorID uint, ipAddress string) error {
	// 查询期刊
	journal, err := s.journalRepo.FindByID(id)
	if err != nil {
		logger.GetLogger().Error("Failed to find journal",
			zap.Uint("journal_id", id),
			zap.Error(err),
		)
		return ErrSystemError
	}

	if journal == nil {
		return ErrJournalNotFound
	}

	// 更新影响因子
	journal.ImpactFactor = impactFactor

	if err := s.journalRepo.Update(journal); err != nil {
		logger.GetLogger().Error("Failed to update journal impact factor",
			zap.Uint("journal_id", id),
			zap.Float64("impact_factor", impactFactor),
			zap.Error(err),
		)
		return ErrSystemError
	}

	// 记录操作日志
	s.operationLogService.LogOperation(operatorID, "update", "journal", journal.ISSN, "更新期刊影响因子 "+journal.FullName, "成功", ipAddress)

	logger.GetLogger().Info("Journal impact factor updated successfully",
		zap.Uint("journal_id", id),
		zap.String("full_name", journal.FullName),
		zap.Float64("impact_factor", impactFactor),
	)

	return nil
}
