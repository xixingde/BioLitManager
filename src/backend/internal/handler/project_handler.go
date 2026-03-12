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

// ProjectHandler 课题处理器
type ProjectHandler struct {
	projectService *service.ProjectService
}

// NewProjectHandler 创建课题处理器实例
func NewProjectHandler(projectService *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

// CreateProject 创建课题
// POST /api/projects
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Code        string `json:"code" binding:"required"`
		ProjectType string `json:"project_type" binding:"required"`
		Source      string `json:"source"`
		Level       string `json:"level"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Warn("Invalid create project request",
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

	project, err := h.projectService.CreateProject(
		req.Name,
		req.Code,
		req.ProjectType,
		req.Source,
		req.Level,
		"进行中", // 默认状态
		operatorID.(uint),
		ipAddress,
	)

	if err != nil {
		if err == service.ErrProjectCodeExists {
			response.Error(c, http.StatusBadRequest, "课题编号已存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	response.Success(c, gin.H{
		"id": project.ID,
	})
}

// GetProject 获取课题详情
// GET /api/projects/:id
func (h *ProjectHandler) GetProject(c *gin.Context) {
	id := c.Param("id")
	projectID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "课题ID格式错误")
		return
	}

	project, err := h.projectService.GetProjectByID(uint(projectID))
	if err != nil {
		if err == service.ErrProjectNotFound {
			response.Error(c, http.StatusNotFound, "课题不存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	// 转换为DTO格式
	projectDTO := gin.H{
		"id":           project.ID,
		"name":         project.Name,
		"code":         project.Code,
		"project_type": project.ProjectType,
		"source":       project.Source,
		"level":        project.Level,
		"status":       project.Status,
		"created_at":   project.CreatedAt,
		"updated_at":   project.UpdatedAt,
	}

	response.Success(c, projectDTO)
}

// ListProjects 分页查询课题列表
// GET /api/projects
func (h *ProjectHandler) ListProjects(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "10")
	name := c.Query("name")
	code := c.Query("code")
	projectType := c.Query("project_type")
	level := c.Query("level")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 || size > 100 {
		size = 10
	}

	projects, total, err := h.projectService.ListProjects(page, size, name, code, projectType, level)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	// 转换为DTO格式
	projectList := make([]gin.H, len(projects))
	for i, project := range projects {
		projectList[i] = gin.H{
			"id":           project.ID,
			"name":         project.Name,
			"code":         project.Code,
			"project_type": project.ProjectType,
			"source":       project.Source,
			"level":        project.Level,
			"status":       project.Status,
			"created_at":   project.CreatedAt,
			"updated_at":   project.UpdatedAt,
		}
	}

	result := gin.H{
		"list":  projectList,
		"total": total,
		"page":  page,
		"size":  size,
	}

	response.Success(c, result)
}

// UpdateProject 更新课题
// PUT /api/projects/:id
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	id := c.Param("id")
	projectID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "课题ID格式错误")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		ProjectType string `json:"project_type"`
		Source      string `json:"source"`
		Level       string `json:"level"`
		Status      string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Warn("Invalid update project request",
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

	if err := h.projectService.UpdateProject(
		uint(projectID),
		req.Name,
		req.Code,
		req.ProjectType,
		req.Source,
		req.Level,
		req.Status,
		operatorID.(uint),
		ipAddress,
	); err != nil {
		if err == service.ErrProjectNotFound {
			response.Error(c, http.StatusNotFound, "课题不存在")
			return
		}
		if err == service.ErrProjectCodeExists {
			response.Error(c, http.StatusBadRequest, "课题编号已存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	response.Success(c, nil)
}

// DeleteProject 删除课题
// DELETE /api/projects/:id
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	id := c.Param("id")
	projectID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "课题ID格式错误")
		return
	}

	// 获取操作者信息
	operatorID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	ipAddress := c.ClientIP()

	if err := h.projectService.DeleteProject(uint(projectID), operatorID.(uint), ipAddress); err != nil {
		if err == service.ErrProjectNotFound {
			response.Error(c, http.StatusNotFound, "课题不存在")
			return
		}
		if err == service.ErrProjectLinked {
			response.Error(c, http.StatusBadRequest, "课题已关联论文，无法删除")
			return
		}
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	response.Success(c, nil)
}
