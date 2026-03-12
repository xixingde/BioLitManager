package service

import (
	"fmt"

	"biolitmanager/internal/repository"
	"biolitmanager/pkg/logger"

	"go.uber.org/zap"
)

// NotificationService 通知服务
type NotificationService struct {
	userRepo            *repository.UserRepository
	operationLogService *OperationLogService
}

// NewNotificationService 创建通知服务实例
func NewNotificationService(
	userRepo *repository.UserRepository,
	operationLogService *OperationLogService,
) *NotificationService {
	return &NotificationService{
		userRepo:            userRepo,
		operationLogService: operationLogService,
	}
}

// SendSubmitNotification 发送提交审核通知
func (s *NotificationService) SendSubmitNotification(paperID uint, submitterID uint, paperTitle string, ipAddress string) error {
	// 查询提交人信息
	submitter, err := s.userRepo.FindByID(submitterID)
	if err != nil {
		logger.GetLogger().Error("Failed to find submitter",
			zap.Uint("submitter_id", submitterID),
			zap.Error(err),
		)
		return ErrSystemError
	}

	if submitter == nil {
		return ErrSystemError
	}

	// TODO: 实现实际的通知发送逻辑（邮件、短信等）
	// 这里暂时只记录日志

	logger.GetLogger().Info("Submit notification sent",
		zap.Uint("paper_id", paperID),
		zap.String("paper_title", paperTitle),
		zap.Uint("submitter_id", submitterID),
		zap.String("submitter_name", submitter.Name),
	)

	// 记录操作日志
	s.operationLogService.LogOperation(submitterID, "submit", "paper", fmt.Sprintf("%d", paperID), fmt.Sprintf("提交论文 %s 审核", paperTitle), "成功", ipAddress)

	return nil
}

// SendReviewNotification 发送审核结果通知
func (s *NotificationService) SendReviewNotification(paperID uint, submitterID uint, paperTitle string, reviewType string, result string, comment string, ipAddress string) error {
	// 查询提交人信息
	submitter, err := s.userRepo.FindByID(submitterID)
	if err != nil {
		logger.GetLogger().Error("Failed to find submitter",
			zap.Uint("submitter_id", submitterID),
			zap.Error(err),
		)
		return ErrSystemError
	}

	if submitter == nil {
		return ErrSystemError
	}

	// TODO: 实现实际的通知发送逻辑（邮件、短信等）
	// 这里暂时只记录日志

	var reviewTypeName string
	if reviewType == "business" {
		reviewTypeName = "业务审核"
	} else if reviewType == "political" {
		reviewTypeName = "政工审核"
	} else {
		reviewTypeName = "审核"
	}

	var resultName string
	if result == "approved" {
		resultName = "通过"
	} else {
		resultName = "驳回"
	}

	logger.GetLogger().Info("Review result notification sent",
		zap.Uint("paper_id", paperID),
		zap.String("paper_title", paperTitle),
		zap.String("review_type", reviewType),
		zap.String("result", result),
		zap.Uint("submitter_id", submitterID),
		zap.String("submitter_name", submitter.Name),
	)

	// 记录操作日志
	s.operationLogService.LogOperation(submitterID, "review_result", "paper", fmt.Sprintf("%d", paperID), fmt.Sprintf("%s%s %s", paperTitle, reviewTypeName, resultName), "成功", ipAddress)

	return nil
}

// SendRejectNotification 发送驳回通知（包含驳回原因）
func (s *NotificationService) SendRejectNotification(paperID uint, submitterID uint, paperTitle string, reviewType string, reason string, ipAddress string) error {
	// 查询提交人信息
	submitter, err := s.userRepo.FindByID(submitterID)
	if err != nil {
		logger.GetLogger().Error("Failed to find submitter",
			zap.Uint("submitter_id", submitterID),
			zap.Error(err),
		)
		return ErrSystemError
	}

	if submitter == nil {
		return ErrSystemError
	}

	// TODO: 实现实际的通知发送逻辑（邮件、短信等）
	// 这里暂时只记录日志

	var reviewTypeName string
	if reviewType == "business" {
		reviewTypeName = "业务审核"
	} else if reviewType == "political" {
		reviewTypeName = "政工审核"
	} else {
		reviewTypeName = "审核"
	}

	logger.GetLogger().Info("Reject notification sent",
		zap.Uint("paper_id", paperID),
		zap.String("paper_title", paperTitle),
		zap.String("review_type", reviewType),
		zap.String("reason", reason),
		zap.Uint("submitter_id", submitterID),
		zap.String("submitter_name", submitter.Name),
	)

	// 记录操作日志
	s.operationLogService.LogOperation(submitterID, "reject", "paper", fmt.Sprintf("%d", paperID), fmt.Sprintf("%s %s 被驳回: %s", paperTitle, reviewTypeName, reason), "成功", ipAddress)

	return nil
}

// SendDeadlineReminder 发送审核时限提醒通知
func (s *NotificationService) SendDeadlineReminder(reviewerID uint, paperID uint, paperTitle string, daysSinceSubmit int, ipAddress string) error {
	// 查询审核人信息
	reviewer, err := s.userRepo.FindByID(reviewerID)
	if err != nil {
		logger.GetLogger().Error("Failed to find reviewer",
			zap.Uint("reviewer_id", reviewerID),
			zap.Error(err),
		)
		return ErrSystemError
	}

	if reviewer == nil {
		return ErrSystemError
	}

	// TODO: 实现实际的通知发送逻辑（邮件、短信等）
	// 这里暂时只记录日志

	logger.GetLogger().Info("Deadline reminder sent",
		zap.Uint("paper_id", paperID),
		zap.String("paper_title", paperTitle),
		zap.Int("days_since_submit", daysSinceSubmit),
		zap.Uint("reviewer_id", reviewerID),
		zap.String("reviewer_name", reviewer.Name),
	)

	// 记录操作日志
	s.operationLogService.LogOperation(reviewerID, "deadline_reminder", "paper", fmt.Sprintf("%d", paperID), fmt.Sprintf("审核时限提醒: %s（已提交%d天）", paperTitle, daysSinceSubmit), "成功", ipAddress)

	return nil
}

// SendApprovalNotification 发送审核通过通知
func (s *NotificationService) SendApprovalNotification(paperID uint, submitterID uint, paperTitle string, ipAddress string) error {
	// 查询提交人信息
	submitter, err := s.userRepo.FindByID(submitterID)
	if err != nil {
		logger.GetLogger().Error("Failed to find submitter",
			zap.Uint("submitter_id", submitterID),
			zap.Error(err),
		)
		return ErrSystemError
	}

	if submitter == nil {
		return ErrSystemError
	}

	// TODO: 实现实际的通知发送逻辑（邮件、短信等）
	// 这里暂时只记录日志

	logger.GetLogger().Info("Approval notification sent",
		zap.Uint("paper_id", paperID),
		zap.String("paper_title", paperTitle),
		zap.Uint("submitter_id", submitterID),
		zap.String("submitter_name", submitter.Name),
	)

	// 记录操作日志
	s.operationLogService.LogOperation(submitterID, "approval", "paper", fmt.Sprintf("%d", paperID), fmt.Sprintf("论文 %s 审核通过", paperTitle), "成功", ipAddress)

	return nil
}

// SendBusinessReviewPassedNotification 发送业务审核通过通知给政工审核人员
func (s *NotificationService) SendBusinessReviewPassedNotification(paperID uint, paperTitle string, ipAddress string) error {
	// TODO: 实现实际的通知发送逻辑（邮件、短信等）
	// 这里暂时只记录日志

	logger.GetLogger().Info("Business review passed notification sent",
		zap.Uint("paper_id", paperID),
		zap.String("paper_title", paperTitle),
	)

	// 记录操作日志
	s.operationLogService.LogOperation(0, "business_review_passed", "paper", fmt.Sprintf("%d", paperID), fmt.Sprintf("论文 %s 业务审核通过，进入政工审核", paperTitle), "成功", ipAddress)

	return nil
}
