package service

import (
	"errors"
	"fmt"
	"time"

	"biolitmanager/internal/cache"
	"biolitmanager/internal/model/entity"
	"biolitmanager/internal/repository"
	"biolitmanager/internal/security"
	"biolitmanager/pkg/logger"

	"go.uber.org/zap"
)

var (
	// ErrUsernamePasswordIncorrect 用户名或密码错误
	ErrUsernamePasswordIncorrect = errors.New("用户名或密码错误")
	// ErrAccountDisabled 账户已禁用
	ErrAccountDisabled = errors.New("账户已禁用")
	// ErrAccountLocked 账户已锁定
	ErrAccountLocked = errors.New("账户已锁定")
	// ErrSystemError 系统异常
	ErrSystemError = errors.New("系统异常")
	// ErrNotFound 资源不存在
	ErrNotFound = errors.New("资源不存在")

	// MaxLoginFailCount 最大登录失败次数
	MaxLoginFailCount = 5
	// LockDuration 锁定时长
	LockDuration = time.Hour
)

// AuthService 认证服务
type AuthService struct {
	userRepo            *repository.UserRepository
	cache               *cache.MemoryCache
	operationLogService *OperationLogService
}

// NewAuthService 创建认证服务实例
func NewAuthService(userRepo *repository.UserRepository, cache *cache.MemoryCache, operationLogService *OperationLogService) *AuthService {
	return &AuthService{
		userRepo:            userRepo,
		cache:               cache,
		operationLogService: operationLogService,
	}
}

// Login 用户登录
func (s *AuthService) Login(username, password, ipAddress string) (string, *entity.User, error) {
	// 查询用户
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		logger.GetLogger().Error("Failed to find user",
			zap.String("username", username),
			zap.Error(err),
		)
		s.operationLogService.LogOperation(0, "login", "auth", "", fmt.Sprintf("用户 %s 登录失败", username), "失败", ipAddress)
		return "", nil, ErrSystemError
	}

	if user == nil {
		logger.GetLogger().Warn("User not found",
			zap.String("username", username),
		)
		return "", nil, ErrUsernamePasswordIncorrect
	}

	// 检查账户是否被禁用
	if user.IsDisabled {
		logger.GetLogger().Warn("Account is disabled",
			zap.String("username", username),
			zap.Uint("user_id", user.ID),
		)
		s.operationLogService.LogOperation(user.ID, "login", "auth", fmt.Sprintf("%d", user.ID), "登录失败：账户已禁用", "失败", ipAddress)
		return "", nil, ErrAccountDisabled
	}

	// 检查账户是否被锁定
	if user.IsLocked && user.LockUntil != nil && time.Now().Before(*user.LockUntil) {
		logger.GetLogger().Warn("Account is locked",
			zap.String("username", username),
			zap.Uint("user_id", user.ID),
			zap.Time("lock_until", *user.LockUntil),
		)
		s.operationLogService.LogOperation(user.ID, "login", "auth", fmt.Sprintf("%d", user.ID), "登录失败：账户已锁定", "失败", ipAddress)
		return "", nil, ErrAccountLocked
	}

	// 验证密码
	if !security.CheckPassword(password, user.PasswordHash) {
		// 密码错误，增加失败计数
		s.incrementLoginFailCount(user, ipAddress)
		return "", nil, ErrUsernamePasswordIncorrect
	}

	// 登录成功，重置失败计数
	s.resetLoginFailCount(user)

	// 获取用户权限列表
	permissions := security.GetPermissionsByRole(security.Role(user.Role))
	permissionStrings := make([]string, len(permissions))
	for i, perm := range permissions {
		permissionStrings[i] = string(perm)
	}

	// 生成 JWT Token
	token, err := security.GenerateToken(user.ID, user.Username, user.Role, permissionStrings)
	if err != nil {
		logger.GetLogger().Error("Failed to generate token",
			zap.Uint("user_id", user.ID),
			zap.String("username", username),
			zap.Error(err),
		)
		return "", nil, ErrSystemError
	}

	// 存储会话到内存缓存（有效期2小时）
	claims := &security.Claims{
		UserID:      user.ID,
		Username:    user.Username,
		Role:        user.Role,
		Permissions: permissionStrings,
	}
	sessionKey := cache.SessionKey(token)
	s.cache.Set(sessionKey, claims, security.TokenExpireTime)
	logger.GetLogger().Info("Session stored in cache",
		zap.String("session_key", sessionKey),
		zap.Uint("user_id", user.ID),
		zap.Duration("ttl", security.TokenExpireTime),
	)

	// 更新最后登录时间和 IP
	now := time.Now()
	user.LastLoginAt = &now
	user.LastLoginIP = ipAddress
	if err := s.userRepo.Update(user); err != nil {
		logger.GetLogger().Error("Failed to update user last login",
			zap.Uint("user_id", user.ID),
			zap.Error(err),
		)
	}

	// 记录登录成功日志
	s.operationLogService.LogOperation(user.ID, "login", "auth", fmt.Sprintf("%d", user.ID), "登录成功", "成功", ipAddress)

	logger.GetLogger().Info("User logged in successfully",
		zap.Uint("user_id", user.ID),
		zap.String("username", username),
		zap.String("ip_address", ipAddress),
	)

	return token, user, nil
}

// Logout 用户登出
func (s *AuthService) Logout(token, ipAddress string, userID uint) error {
	// 从缓存中删除会话
	sessionKey := cache.SessionKey(token)
	s.cache.Delete(sessionKey)

	// 记录登出日志
	s.operationLogService.LogOperation(userID, "logout", "auth", fmt.Sprintf("%d", userID), "用户登出", "成功", ipAddress)

	logger.GetLogger().Info("User logged out",
		zap.Uint("user_id", userID),
		zap.String("ip_address", ipAddress),
	)

	return nil
}

// GetProfile 获取当前用户信息
func (s *AuthService) GetProfile(userID uint) (*entity.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		logger.GetLogger().Error("Failed to find user",
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
		return nil, ErrSystemError
	}

	if user == nil {
		return nil, ErrNotFound
	}

	return user, nil
}

// incrementLoginFailCount 增加登录失败计数
func (s *AuthService) incrementLoginFailCount(user *entity.User, ipAddress string) {
	user.LoginFailCount++

	// 如果达到最大失败次数，锁定账户
	if user.LoginFailCount >= MaxLoginFailCount {
		lockUntil := time.Now().Add(LockDuration)
		user.IsLocked = true
		user.LockUntil = &lockUntil

		logger.GetLogger().Warn("Account locked due to too many failed login attempts",
			zap.Uint("user_id", user.ID),
			zap.String("username", user.Username),
			zap.Int("fail_count", user.LoginFailCount),
		)
	}

	if err := s.userRepo.Update(user); err != nil {
		logger.GetLogger().Error("Failed to update user fail count",
			zap.Uint("user_id", user.ID),
			zap.Error(err),
		)
	}

	// 记录登录失败日志
	s.operationLogService.LogOperation(user.ID, "login", "auth", fmt.Sprintf("%d", user.ID), fmt.Sprintf("登录失败：密码错误（失败次数：%d）", user.LoginFailCount), "失败", ipAddress)
}

// resetLoginFailCount 重置登录失败计数
func (s *AuthService) resetLoginFailCount(user *entity.User) {
	if user.LoginFailCount > 0 || user.IsLocked {
		user.LoginFailCount = 0
		user.IsLocked = false
		user.LockUntil = nil

		if err := s.userRepo.Update(user); err != nil {
			logger.GetLogger().Error("Failed to reset user fail count",
				zap.Uint("user_id", user.ID),
				zap.Error(err),
			)
		}
	}
}
