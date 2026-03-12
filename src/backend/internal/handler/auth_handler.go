package handler

import (
	"net/http"

	"biolitmanager/internal/model/dto/response"
	"biolitmanager/internal/service"
	"biolitmanager/pkg/logger"
	resp "biolitmanager/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login 用户登录
// POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req response.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Warn("Invalid login request",
			zap.Error(err),
		)
		resp.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	ipAddress := c.ClientIP()

	token, user, err := h.authService.Login(req.Username, req.Password, ipAddress)
	if err != nil {
		if err == service.ErrUsernamePasswordIncorrect {
			resp.Error(c, http.StatusBadRequest, "用户名或密码错误")
			return
		}
		if err == service.ErrAccountDisabled {
			resp.Error(c, http.StatusForbidden, "账户已禁用")
			return
		}
		if err == service.ErrAccountLocked {
			resp.Error(c, http.StatusForbidden, "账户已锁定，请稍后再试")
			return
		}
		resp.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	// 返回登录响应
	loginResp := response.LoginResponse{
		Token: token,
		User: response.UserDTO{
			ID:         user.ID,
			Username:   user.Username,
			Name:       user.Name,
			Role:       user.Role,
			Department: user.Department,
			Email:      user.Email,
		},
	}

	resp.Success(c, loginResp)
}

// Logout 用户登出
// POST /api/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// 从上下文获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		resp.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	// 从请求头获取 Token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		resp.Error(c, http.StatusBadRequest, "Token 不能为空")
		return
	}

	// 解析 Bearer Token
	token := authHeader[7:] // 去掉 "Bearer " 前缀

	ipAddress := c.ClientIP()

	if err := h.authService.Logout(token, ipAddress, userID.(uint)); err != nil {
		logger.GetLogger().Error("Failed to logout",
			zap.Uint("user_id", userID.(uint)),
			zap.Error(err),
		)
		resp.Error(c, http.StatusInternalServerError, "登出失败")
		return
	}

	resp.Success(c, nil)
}

// GetProfile 获取当前用户信息
// GET /api/auth/profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// 从上下文获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		resp.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	user, err := h.authService.GetProfile(userID.(uint))
	if err != nil {
		if err == service.ErrNotFound {
			resp.Error(c, http.StatusNotFound, "用户不存在")
			return
		}
		resp.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	// 返回用户信息
	userDTO := response.UserDTO{
		ID:         user.ID,
		Username:   user.Username,
		Name:       user.Name,
		Role:       user.Role,
		Department: user.Department,
		Email:      user.Email,
	}

	resp.Success(c, userDTO)
}
