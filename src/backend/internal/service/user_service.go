package service

import (
	"fmt"

	"biolitmanager/internal/model/entity"
	"biolitmanager/internal/repository"
	"biolitmanager/internal/security"
	"biolitmanager/pkg/logger"

	"go.uber.org/zap"
)

var (
	// ErrUsernameAlreadyExists 用户名已存在
	ErrUsernameAlreadyExists = fmt.Errorf("用户名已存在")
	// ErrPasswordTooWeak 密码强度不足
	ErrPasswordTooWeak = fmt.Errorf("密码强度不足")
)

// UserService 用户服务
type UserService struct {
	userRepo            *repository.UserRepository
	operationLogService *OperationLogService
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo *repository.UserRepository, operationLogService *OperationLogService) *UserService {
	return &UserService{
		userRepo:            userRepo,
		operationLogService: operationLogService,
	}
}

// CreateUser 创建用户
func (s *UserService) CreateUser(username, password, name, role, department, idCard, phone, email string, operatorID uint, ipAddress string) (*entity.User, error) {
	// 校验用户名是否已存在
	existingUser, err := s.userRepo.FindByUsername(username)
	if err != nil {
		logger.GetLogger().Error("Failed to check username existence",
			zap.String("username", username),
			zap.Error(err),
		)
		return nil, ErrSystemError
	}

	if existingUser != nil {
		logger.GetLogger().Warn("Username already exists",
			zap.String("username", username),
		)
		return nil, ErrUsernameAlreadyExists
	}

	// 校验密码复杂度
	if err := security.ValidatePasswordComplexity(password); err != nil {
		logger.GetLogger().Warn("Password complexity validation failed",
			zap.String("username", username),
			zap.Error(err),
		)
		return nil, ErrPasswordTooWeak
	}

	// 密码哈希
	passwordHash, err := security.HashPassword(password)
	if err != nil {
		logger.GetLogger().Error("Failed to hash password",
			zap.String("username", username),
			zap.Error(err),
		)
		return nil, ErrSystemError
	}

	// 创建用户记录
	user := &entity.User{
		Username:     username,
		PasswordHash: passwordHash,
		Name:         name,
		Role:         role,
		Department:   department,
		IDCard:       idCard,
		Phone:        phone,
		Email:        email,
	}

	if err := s.userRepo.Create(user); err != nil {
		logger.GetLogger().Error("Failed to create user",
			zap.String("username", username),
			zap.Error(err),
		)
		return nil, ErrSystemError
	}

	// 记录操作日志
	s.operationLogService.LogOperation(operatorID, "create", "user", fmt.Sprintf("%d", user.ID), fmt.Sprintf("创建用户 %s", username), "成功", ipAddress)

	logger.GetLogger().Info("User created successfully",
		zap.Uint("user_id", user.ID),
		zap.String("username", username),
		zap.Uint("operator_id", operatorID),
	)

	return user, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(userID uint, name, department, idCard, phone, email string, operatorID uint, ipAddress string) error {
	// 查询用户
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		logger.GetLogger().Error("Failed to find user",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return ErrSystemError
	}

	if user == nil {
		return ErrNotFound
	}

	// 更新用户信息
	user.Name = name
	user.Department = department
	user.IDCard = idCard
	user.Phone = phone
	user.Email = email

	if err := s.userRepo.Update(user); err != nil {
		logger.GetLogger().Error("Failed to update user",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return ErrSystemError
	}

	// 记录操作日志
	s.operationLogService.LogOperation(operatorID, "update", "user", fmt.Sprintf("%d", userID), fmt.Sprintf("更新用户 %s", user.Username), "成功", ipAddress)

	logger.GetLogger().Info("User updated successfully",
		zap.Uint("user_id", userID),
		zap.String("username", user.Username),
		zap.Uint("operator_id", operatorID),
	)

	return nil
}

// DisableUser 禁用用户账户
func (s *UserService) DisableUser(userID uint, operatorID uint, ipAddress string) error {
	// 查询用户
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		logger.GetLogger().Error("Failed to find user",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return ErrSystemError
	}

	if user == nil {
		return ErrNotFound
	}

	// 禁用用户
	user.IsDisabled = true

	if err := s.userRepo.Update(user); err != nil {
		logger.GetLogger().Error("Failed to disable user",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return ErrSystemError
	}

	// 记录操作日志
	s.operationLogService.LogOperation(operatorID, "disable", "user", fmt.Sprintf("%d", userID), fmt.Sprintf("禁用用户 %s", user.Username), "成功", ipAddress)

	logger.GetLogger().Info("User disabled successfully",
		zap.Uint("user_id", userID),
		zap.String("username", user.Username),
		zap.Uint("operator_id", operatorID),
	)

	return nil
}

// EnableUser 启用用户账户
func (s *UserService) EnableUser(userID uint, operatorID uint, ipAddress string) error {
	// 查询用户
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		logger.GetLogger().Error("Failed to find user",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return ErrSystemError
	}

	if user == nil {
		return ErrNotFound
	}

	// 启用用户
	user.IsDisabled = false

	if err := s.userRepo.Update(user); err != nil {
		logger.GetLogger().Error("Failed to enable user",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return ErrSystemError
	}

	// 记录操作日志
	s.operationLogService.LogOperation(operatorID, "enable", "user", fmt.Sprintf("%d", userID), fmt.Sprintf("启用用户 %s", user.Username), "成功", ipAddress)

	logger.GetLogger().Info("User enabled successfully",
		zap.Uint("user_id", userID),
		zap.String("username", user.Username),
		zap.Uint("operator_id", operatorID),
	)

	return nil
}

// ListUsers 分页查询用户列表
func (s *UserService) ListUsers(page, size int) ([]*entity.User, int64, error) {
	users, total, err := s.userRepo.List(page, size)
	if err != nil {
		logger.GetLogger().Error("Failed to list users",
			zap.Int("page", page),
			zap.Int("size", size),
			zap.Error(err),
		)
		return nil, 0, ErrSystemError
	}

	return users, total, nil
}
