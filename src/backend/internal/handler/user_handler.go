package handler

import (
	"net/http"
	"strconv"

	"biolitmanager/internal/service"
	"biolitmanager/pkg/logger"
	"biolitmanager/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Role       string `json:"role" binding:"required"`
	Department string `json:"department"`
	IDCard     string `json:"id_card"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Name       string `json:"name" binding:"required"`
	Department string `json:"department"`
	IDCard     string `json:"id_card"`
	Phone      string `json:"phone"`
	Email      string `json:"email"`
}

// UserHandler 用户处理器
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler 创建用户处理器实例
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUser 创建用户
// POST /api/system/users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Warn("Invalid create user request",
			zap.Error(err),
		)
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 获取操作者信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	ipAddress := c.ClientIP()

	user, err := h.userService.CreateUser(req.Username, req.Password, req.Name, req.Role, req.Department, req.IDCard, req.Phone, req.Email, userID.(uint), ipAddress)
	if err != nil {
		if err == service.ErrUsernameAlreadyExists {
			response.Error(c, http.StatusBadRequest, "用户名已存在")
			return
		}
		if err == service.ErrPasswordTooWeak {
			response.Error(c, http.StatusBadRequest, "密码强度不足")
			return
		}
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	// 返回用户信息
	userDTO := map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"name":       user.Name,
		"role":       user.Role,
		"department": user.Department,
		"email":      user.Email,
	}

	response.Success(c, userDTO)
}

// UpdateUser 更新用户信息
// PUT /api/system/users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "用户ID格式错误")
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Warn("Invalid update user request",
			zap.Error(err),
		)
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	// 获取操作者信息
	operatorID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	ipAddress := c.ClientIP()

	if err := h.userService.UpdateUser(uint(userID), req.Name, req.Department, req.IDCard, req.Phone, req.Email, operatorID.(uint), ipAddress); err != nil {
		if err == service.ErrNotFound {
			response.Error(c, http.StatusNotFound, "用户不存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	response.Success(c, nil)
}

// DisableUser 禁用用户账户
// PUT /api/system/users/:id/disable
func (h *UserHandler) DisableUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "用户ID格式错误")
		return
	}

	// 获取操作者信息
	operatorID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	ipAddress := c.ClientIP()

	if err := h.userService.DisableUser(uint(userID), operatorID.(uint), ipAddress); err != nil {
		if err == service.ErrNotFound {
			response.Error(c, http.StatusNotFound, "用户不存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	response.Success(c, nil)
}

// EnableUser 启用用户账户
// PUT /api/system/users/:id/enable
func (h *UserHandler) EnableUser(c *gin.Context) {
	id := c.Param("id")
	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "用户ID格式错误")
		return
	}

	// 获取操作者信息
	operatorID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	ipAddress := c.ClientIP()

	if err := h.userService.EnableUser(uint(userID), operatorID.(uint), ipAddress); err != nil {
		if err == service.ErrNotFound {
			response.Error(c, http.StatusNotFound, "用户不存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	response.Success(c, nil)
}

// ListUsers 分页查询用户列表
// GET /api/system/users
func (h *UserHandler) ListUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 || size > 100 {
		size = 10
	}

	users, total, err := h.userService.ListUsers(page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	// 转换为 DTO 格式
	userList := make([]map[string]interface{}, len(users))
	for i, user := range users {
		userList[i] = map[string]interface{}{
			"id":          user.ID,
			"username":    user.Username,
			"name":        user.Name,
			"role":        user.Role,
			"department":  user.Department,
			"email":       user.Email,
			"is_disabled": user.IsDisabled,
			"is_locked":   user.IsLocked,
		}
	}

	result := response.PageResult{
		List:  userList,
		Total: total,
		Page:  page,
		Size:  size,
	}

	response.Success(c, result)
}
