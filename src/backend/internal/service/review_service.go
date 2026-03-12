package service

import (
	"errors"
	"fmt"
	"time"

	"biolitmanager/internal/model/entity"
	"biolitmanager/internal/repository"
	"biolitmanager/internal/security"
	apperrors "biolitmanager/pkg/errors"
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

// ReviewServiceInterface 审核服务接口
type ReviewServiceInterface interface {
	BusinessReview(paperID uint, result, comment string, reviewerID uint, ipAddress string) error
	PoliticalReview(paperID uint, result, comment string, reviewerID uint, ipAddress string) error
	GetReviewLogsByPaperID(paperID uint) ([]*entity.ReviewLog, error)
	GetPendingPapersForBusinessReview() ([]*entity.Paper, error)
	GetPendingPapersForPoliticalReview() ([]*entity.Paper, error)
	GetMyReviews(reviewerID uint) ([]*entity.ReviewLog, error)
}

// AppError 应用错误类型
type AppError struct {
	Code int
	Msg  string
}

func (e *AppError) Error() string {
	return e.Msg
}

// ReviewService 审核服务
type ReviewService struct {
	db                  *gorm.DB
	reviewRepo          *repository.ReviewRepository
	paperRepo           *repository.PaperRepository
	userRepo            *repository.UserRepository
	operationLogService *OperationLogService
	notificationService *NotificationService
	archiveService      *ArchiveService
}

// NewReviewService 创建审核服务实例
func NewReviewService(
	db *gorm.DB,
	reviewRepo *repository.ReviewRepository,
	paperRepo *repository.PaperRepository,
	userRepo *repository.UserRepository,
	operationLogService *OperationLogService,
	notificationService *NotificationService,
	archiveService *ArchiveService,
) *ReviewService {
	return &ReviewService{
		db:                  db,
		reviewRepo:          reviewRepo,
		paperRepo:           paperRepo,
		userRepo:            userRepo,
		operationLogService: operationLogService,
		notificationService: notificationService,
		archiveService:      archiveService,
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

	// 审核通过后自动归档（归档失败不影响审核流程）
	if result == "通过" {
		go func() {
			err := s.ProcessApprovePaper(paperID, paper.Title, "审核通过", reviewerID)
			if err != nil {
				logger.GetLogger().Error("Auto archive failed after political review approval",
					zap.Uint("paper_id", paperID),
					zap.Error(err),
				)
			}
		}()
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

// ValidateReviewPermission 校验用户是否有审核权限
func (s *ReviewService) ValidateReviewPermission(reviewerID uint, reviewType string) error {
	user, err := s.userRepo.FindByID(reviewerID)
	if err != nil {
		logger.GetLogger().Error("Failed to find reviewer",
			zap.Uint("reviewer_id", reviewerID),
			zap.Error(err),
		)
		return ErrSystemError
	}

	if user == nil {
		logger.GetLogger().Warn("Reviewer not found",
			zap.Uint("reviewer_id", reviewerID),
		)
		return &AppError{
			Code: apperrors.ErrReviewerNotFound.Code,
			Msg:  apperrors.ErrReviewerNotFound.Msg,
		}
	}

	role := security.Role(user.Role)
	var requiredPermission security.Permission

	if reviewType == "业务审核" {
		requiredPermission = security.PermissionReviewBusiness
	} else if reviewType == "政工审核" {
		requiredPermission = security.PermissionReviewPolitical
	} else {
		logger.GetLogger().Warn("Invalid review type",
			zap.Uint("reviewer_id", reviewerID),
			zap.String("review_type", reviewType),
		)
		return &AppError{
			Code: apperrors.ParamError.Code,
			Msg:  fmt.Sprintf("无效的审核类型: %s", reviewType),
		}
	}

	if !security.HasPermission(role, requiredPermission) {
		logger.GetLogger().Warn("User has no review permission",
			zap.Uint("reviewer_id", reviewerID),
			zap.String("role", user.Role),
			zap.String("review_type", reviewType),
		)
		return &AppError{
			Code: apperrors.ErrNoReviewPermission.Code,
			Msg:  apperrors.ErrNoReviewPermission.Msg,
		}
	}

	return nil
}

// CheckReviewDeadline 计算距离提交的工作日数,判断是否超过3个工作日
func (s *ReviewService) CheckReviewDeadline(submitTime time.Time) (days int, overdue bool, err error) {
	if submitTime.After(time.Now()) {
		return 0, false, nil
	}

	days = 0
	current := submitTime
	now := time.Now()

	for current.Before(now) {
		current = current.Add(24 * time.Hour)
		weekday := current.Weekday()
		if weekday != time.Saturday && weekday != time.Sunday {
			days++
		}
	}

	overdue = days > 3

	logger.GetLogger().Debug("Review deadline check",
		zap.Time("submit_time", submitTime),
		zap.Int("days", days),
		zap.Bool("overdue", overdue),
	)

	return days, overdue, nil
}

// SendReviewReminder 发送审核时限提醒通知给审核人员
func (s *ReviewService) SendReviewReminder(reviewerID uint, paperID uint, reviewType string) error {
	paper, err := s.paperRepo.FindByID(paperID)
	if err != nil {
		logger.GetLogger().Error("Failed to find paper",
			zap.Uint("paper_id", paperID),
			zap.Error(err),
		)
		return ErrSystemError
	}

	if paper == nil {
		logger.GetLogger().Warn("Paper not found",
			zap.Uint("paper_id", paperID),
		)
		return &AppError{
			Code: apperrors.ErrPaperNotFound.Code,
			Msg:  apperrors.ErrPaperNotFound.Msg,
		}
	}

	daysSinceSubmit, _, err := s.CheckReviewDeadline(paper.SubmitTime)
	if err != nil {
		return err
	}

	err = s.notificationService.SendDeadlineReminder(reviewerID, paperID, paper.Title, daysSinceSubmit, "")
	if err != nil {
		logger.GetLogger().Error("Failed to send review reminder",
			zap.Uint("paper_id", paperID),
			zap.Uint("reviewer_id", reviewerID),
			zap.Error(err),
		)
		return err
	}

	logger.GetLogger().Info("Review reminder sent",
		zap.Uint("paper_id", paperID),
		zap.String("paper_title", paper.Title),
		zap.Uint("reviewer_id", reviewerID),
		zap.String("review_type", reviewType),
	)

	return nil
}

// ProcessRejectPaper 处理驳回论文,重置状态为草稿,发送驳回通知
func (s *ReviewService) ProcessRejectPaper(paperID uint, rejectReason string, paperTitle string) error {
	paper, err := s.paperRepo.FindByID(paperID)
	if err != nil {
		logger.GetLogger().Error("Failed to find paper",
			zap.Uint("paper_id", paperID),
			zap.Error(err),
		)
		return ErrSystemError
	}

	if paper == nil {
		logger.GetLogger().Warn("Paper not found",
			zap.Uint("paper_id", paperID),
		)
		return &AppError{
			Code: apperrors.ErrPaperNotFound.Code,
			Msg:  apperrors.ErrPaperNotFound.Msg,
		}
	}

	err = s.db.Model(&entity.Paper{}).Where("id = ?", paperID).Update("status", "草稿").Error
	if err != nil {
		logger.GetLogger().Error("Failed to update paper status to draft",
			zap.Uint("paper_id", paperID),
			zap.Error(err),
		)
		return ErrSystemError
	}

	determineReviewType := func(status string) string {
		switch status {
		case "业务审核驳回":
			return "business"
		case "政工审核驳回":
			return "political"
		default:
			return ""
		}
	}

	reviewType := determineReviewType(paper.Status)
	if reviewType != "" {
		err = s.notificationService.SendRejectNotification(paperID, paper.SubmitterID, paperTitle, reviewType, rejectReason, "")
		if err != nil {
			logger.GetLogger().Error("Failed to send reject notification",
				zap.Uint("paper_id", paperID),
				zap.Uint("submitter_id", paper.SubmitterID),
				zap.Error(err),
			)
			return err
		}
	}

	logger.GetLogger().Info("Paper reject processed",
		zap.Uint("paper_id", paperID),
		zap.String("paper_title", paperTitle),
		zap.String("reason", rejectReason),
	)

	return nil
}

// ProcessApprovePaper 处理审核通过论文,调用归档Service创建归档记录,发送通过通知
func (s *ReviewService) ProcessApprovePaper(paperID uint, paperTitle string, paperStatus string, reviewerID uint) error {
	paper, err := s.paperRepo.FindByID(paperID)
	if err != nil {
		logger.GetLogger().Error("Failed to find paper",
			zap.Uint("paper_id", paperID),
			zap.Error(err),
		)
		return ErrSystemError
	}

	if paper == nil {
		logger.GetLogger().Warn("Paper not found",
			zap.Uint("paper_id", paperID),
		)
		return &AppError{
			Code: apperrors.ErrPaperNotFound.Code,
			Msg:  apperrors.ErrPaperNotFound.Msg,
		}
	}

	// 归档失败不影响审核流程，记录日志但返回成功
	if paperStatus == "政工审核" || paperStatus == "审核通过" {
		_, err = s.archiveService.ArchivePaper(paperID, reviewerID, "")
		if err != nil {
			logger.GetLogger().Error("Failed to archive paper, but review process continues",
				zap.Uint("paper_id", paperID),
				zap.Error(err),
			)
		} else {
			logger.GetLogger().Info("Paper archived successfully during approval",
				zap.Uint("paper_id", paperID),
				zap.Uint("reviewer_id", reviewerID),
			)
		}
	}

	err = s.notificationService.SendApprovalNotification(paperID, paper.SubmitterID, paperTitle, "")
	if err != nil {
		logger.GetLogger().Error("Failed to send approval notification",
			zap.Uint("paper_id", paperID),
			zap.Uint("submitter_id", paper.SubmitterID),
			zap.Error(err),
		)
		return err
	}

	logger.GetLogger().Info("Paper approval processed",
		zap.Uint("paper_id", paperID),
		zap.String("paper_title", paperTitle),
		zap.String("status", paperStatus),
	)

	return nil
}

// GetMyReviews 获取审核员的审核记录
func (s *ReviewService) GetMyReviews(reviewerID uint) ([]*entity.ReviewLog, error) {
	reviews, err := s.reviewRepo.ListByReviewer(reviewerID)
	if err != nil {
		logger.GetLogger().Error("Failed to get my reviews",
			zap.Uint("reviewer_id", reviewerID),
			zap.Error(err),
		)
		return nil, err
	}

	return reviews, nil
}

// SendReviewReminderForOverdue 检查逾期审核并发送提醒
// 查询待审核且超过2个工作日的论文，发送提醒通知
func (s *ReviewService) SendReviewReminderForOverdue() error {
	logger.GetLogger().Info("Starting review reminder check for overdue papers")

	// 查询待业务审核的论文
	businessPapers, err := s.GetPendingPapersForBusinessReview()
	if err != nil {
		logger.GetLogger().Error("Failed to get pending business review papers",
			zap.Error(err),
		)
		return err
	}

	// 查询待政工审核的论文
	politicalPapers, err := s.GetPendingPapersForPoliticalReview()
	if err != nil {
		logger.GetLogger().Error("Failed to get pending political review papers",
			zap.Error(err),
		)
		return err
	}

	// 合并待审核论文
	allPapers := append(businessPapers, politicalPapers...)

	// 统计信息
	var sentCount int
	var skipCount int

	for _, paper := range allPapers {
		// 计算距离提交时间的工作日数
		days, overdue, err := s.CheckReviewDeadline(paper.SubmitTime)
		if err != nil {
			logger.GetLogger().Error("Failed to check review deadline",
				zap.Uint("paper_id", paper.ID),
				zap.Error(err),
			)
			continue
		}

		// 只处理超过2个工作日的论文
		if !overdue || days <= 2 {
			skipCount++
			continue
		}

		// 检查是否已发送过提醒（通过查询操作日志）
		hasReminder, err := s.hasSentReminder(paper.ID)
		if err != nil {
			logger.GetLogger().Error("Failed to check if reminder was sent",
				zap.Uint("paper_id", paper.ID),
				zap.Error(err),
			)
			continue
		}

		if hasReminder {
			skipCount++
			continue
		}

		// 确定审核类型和对应的审核人员
		reviewType := ""
		reviewerID := uint(0)

		if paper.Status == "待业务审核" {
			reviewType = "业务审核"
			// TODO: 查询业务审核人员ID，这里暂时使用提交者所在部门的审核人员
			reviewerID = s.getBusinessReviewerID(paper)
		} else if paper.Status == "待政工审核" {
			reviewType = "政工审核"
			// TODO: 查询政工审核人员ID
			reviewerID = s.getPoliticalReviewerID(paper)
		}

		if reviewerID == 0 {
			logger.GetLogger().Warn("No reviewer found for paper",
				zap.Uint("paper_id", paper.ID),
				zap.String("status", paper.Status),
			)
			continue
		}

		// 发送提醒通知
		err = s.SendReviewReminder(reviewerID, paper.ID, reviewType)
		if err != nil {
			logger.GetLogger().Error("Failed to send review reminder",
				zap.Uint("paper_id", paper.ID),
				zap.Uint("reviewer_id", reviewerID),
				zap.Error(err),
			)
			continue
		}

		sentCount++
		logger.GetLogger().Info("Review reminder sent successfully",
			zap.Uint("paper_id", paper.ID),
			zap.String("paper_title", paper.Title),
			zap.String("review_type", reviewType),
			zap.Int("days", days),
		)
	}

	logger.GetLogger().Info("Review reminder check completed",
		zap.Int("total_papers", len(allPapers)),
		zap.Int("sent_count", sentCount),
		zap.Int("skip_count", skipCount),
	)

	return nil
}

// hasSentReminder 检查是否已发送过提醒
func (s *ReviewService) hasSentReminder(paperID uint) (bool, error) {
	// 查询最近的操作日志，检查是否有 deadline_reminder 类型的记录
	var count int64
	err := s.db.Model(&entity.OperationLog{}).
		Where("target_type = ? AND target_id = ? AND operation_type = ?", "paper", paperID, "deadline_reminder").
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// getBusinessReviewerID 获取业务审核人员ID
// TODO: 实现具体的业务审核人员查询逻辑
func (s *ReviewService) getBusinessReviewerID(paper *entity.Paper) uint {
	// 暂时返回0，实际应该根据论文的提交者部门查询对应的业务审核人员
	// 可以通过查询具有业务审核权限的用户来确定
	return 0
}

// getPoliticalReviewerID 获取政工审核人员ID
// TODO: 实现具体的政工审核人员查询逻辑
func (s *ReviewService) getPoliticalReviewerID(paper *entity.Paper) uint {
	// 暂时返回0，实际应该查询政工审核人员
	return 0
}
