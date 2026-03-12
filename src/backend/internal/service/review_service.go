package service

import (
	"errors"
	"fmt"
	"time"

	"biolitmanager/internal/model/entity"
	"biolitmanager/internal/repository"
	"biolitmanager/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	// ErrReviewNotFound 审核记录不存在
	ErrReviewNotFound = errors.New("审核记录不存在")
	// ErrInvalidReviewResult 无效的审核结果
	ErrInvalidReviewResult = errors.New("无效的审核结果")
)

// ReviewService 审核服务
type ReviewService struct {
	db                  *gorm.DB
	reviewRepo          *repository.ReviewRepository
	paperRepo           *repository.PaperRepository
	operationLogService *OperationLogService
}

// NewReviewService 创建审核服务实例
func NewReviewService(
	db *gorm.DB,
	reviewRepo *repository.ReviewRepository,
	paperRepo *repository.PaperRepository,
	operationLogService *OperationLogService,
) *ReviewService {
	return &ReviewService{
		db:                  db,
		reviewRepo:          reviewRepo,
		paperRepo:           paperRepo,
		operationLogService: operationLogService,
	}
}

// BusinessReview 业务审核
func (s *ReviewService) BusinessReview(paperID uint, result, comment string, reviewerID uint, ipAddress string) error {
	// 校验审核结果
	if result != "通过" && result != "驳回" {
		logger.GetLogger().Warn("Invalid review result",
			zap.Uint("reviewer_id", reviewerID),
			zap.String("result", result),
		)
		return ErrInvalidReviewResult
	}

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

	// 使用事务创建审核记录并更新论文状态
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 创建审核记录
		reviewLog := &entity.ReviewLog{
			PaperID:    paperID,
			ReviewerID: reviewerID,
			ReviewType: "业务审核",
			Result:     result,
			Comment:    comment,
			ReviewTime: time.Now(),
		}

		if err := tx.Create(reviewLog).Error; err != nil {
			logger.GetLogger().Error("Failed to create review log",
				zap.Uint("paper_id", paperID),
				zap.Error(err),
			)
			return ErrSystemError
		}

		// 更新论文状态
		var newStatus string
		if result == "通过" {
			newStatus = "待政工审核"
		} else {
			newStatus = "业务审核驳回"
		}

		if err := tx.Model(&entity.Paper{}).Where("id = ?", paperID).Update("status", newStatus).Error; err != nil {
			logger.GetLogger().Error("Failed to update paper status",
				zap.Uint("paper_id", paperID),
				zap.String("status", newStatus),
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
	s.operationLogService.LogOperation(reviewerID, "review", "paper", fmt.Sprintf("%d", paperID), fmt.Sprintf("业务审核论文 %s：%s", paper.Title, result), "成功", ipAddress)

	logger.GetLogger().Info("Business review completed",
		zap.Uint("paper_id", paperID),
		zap.String("title", paper.Title),
		zap.Uint("reviewer_id", reviewerID),
		zap.String("result", result),
	)

	return nil
}

// PoliticalReview 政工审核
func (s *ReviewService) PoliticalReview(paperID uint, result, comment string, reviewerID uint, ipAddress string) error {
	// 校验审核结果
	if result != "通过" && result != "驳回" {
		logger.GetLogger().Warn("Invalid review result",
			zap.Uint("reviewer_id", reviewerID),
			zap.String("result", result),
		)
		return ErrInvalidReviewResult
	}

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

	// 使用事务创建审核记录并更新论文状态
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 创建审核记录
		reviewLog := &entity.ReviewLog{
			PaperID:    paperID,
			ReviewerID: reviewerID,
			ReviewType: "政工审核",
			Result:     result,
			Comment:    comment,
			ReviewTime: time.Now(),
		}

		if err := tx.Create(reviewLog).Error; err != nil {
			logger.GetLogger().Error("Failed to create review log",
				zap.Uint("paper_id", paperID),
				zap.Error(err),
			)
			return ErrSystemError
		}

		// 更新论文状态
		var newStatus string
		if result == "通过" {
			newStatus = "审核通过"
		} else {
			newStatus = "政工审核驳回"
		}

		if err := tx.Model(&entity.Paper{}).Where("id = ?", paperID).Update("status", newStatus).Error; err != nil {
			logger.GetLogger().Error("Failed to update paper status",
				zap.Uint("paper_id", paperID),
				zap.String("status", newStatus),
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
	s.operationLogService.LogOperation(reviewerID, "review", "paper", fmt.Sprintf("%d", paperID), fmt.Sprintf("政工审核论文 %s：%s", paper.Title, result), "成功", ipAddress)

	logger.GetLogger().Info("Political review completed",
		zap.Uint("paper_id", paperID),
		zap.String("title", paper.Title),
		zap.Uint("reviewer_id", reviewerID),
		zap.String("result", result),
	)

	return nil
}

// GetReviewLogsByPaperID 获取论文的审核记录
func (s *ReviewService) GetReviewLogsByPaperID(paperID uint) ([]*entity.ReviewLog, error) {
	reviewLogs, err := s.reviewRepo.FindByPaperID(paperID)
	if err != nil {
		logger.GetLogger().Error("Failed to find review logs",
			zap.Uint("paper_id", paperID),
			zap.Error(err),
		)
		return nil, ErrSystemError
	}

	return reviewLogs, nil
}

// GetPendingPapersForBusinessReview 获取待业务审核的论文列表
func (s *ReviewService) GetPendingPapersForBusinessReview() ([]*entity.Paper, error) {
	var papers []*entity.Paper
	err := s.db.Preload("Submitter").
		Where("status = ?", "待业务审核").
		Order("submit_time ASC").
		Find(&papers).Error
	if err != nil {
		logger.GetLogger().Error("Failed to get pending papers for business review",
			zap.Error(err),
		)
		return nil, ErrSystemError
	}

	return papers, nil
}

// GetPendingPapersForPoliticalReview 获取待政工审核的论文列表
func (s *ReviewService) GetPendingPapersForPoliticalReview() ([]*entity.Paper, error) {
	var papers []*entity.Paper
	err := s.db.Preload("Submitter").
		Preload("ReviewLogs", "review_type = ?", "业务审核").
		Where("status = ?", "待政工审核").
		Order("submit_time ASC").
		Find(&papers).Error
	if err != nil {
		logger.GetLogger().Error("Failed to get pending papers for political review",
			zap.Error(err),
		)
		return nil, ErrSystemError
	}

	return papers, nil
}
