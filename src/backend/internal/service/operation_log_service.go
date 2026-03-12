package service

import (
	"biolitmanager/internal/model/entity"
	"biolitmanager/internal/repository"
	"biolitmanager/pkg/logger"

	"go.uber.org/zap"
)

// OperationLogService 操作日志服务
type OperationLogService struct {
	operationLogRepo *repository.OperationLogRepository
}

// NewOperationLogService 创建操作日志服务实例
func NewOperationLogService(operationLogRepo *repository.OperationLogRepository) *OperationLogService {
	return &OperationLogService{
		operationLogRepo: operationLogRepo,
	}
}

// LogOperation 记录操作日志
func (s *OperationLogService) LogOperation(userID uint, operationType, module, targetID, operationContent, operationResult, ipAddress string) error {
	log := &entity.OperationLog{
		UserID:           userID,
		OperationType:    operationType,
		Module:           module,
		TargetID:         targetID,
		OperationContent: operationContent,
		OperationResult:  operationResult,
		IPAddress:        ipAddress,
	}

	if err := s.operationLogRepo.Create(log); err != nil {
		logger.GetLogger().Error("Failed to create operation log",
			zap.Uint("user_id", userID),
			zap.String("operation_type", operationType),
			zap.Error(err),
		)
		return err
	}

	logger.GetLogger().Info("Operation log created",
		zap.Uint("user_id", userID),
		zap.String("operation_type", operationType),
		zap.String("module", module),
	)

	return nil
}
