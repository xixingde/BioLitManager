package handler

import (
	"net/http"
	"strconv"

	"biolitmanager/internal/service"
	"biolitmanager/pkg/response"

	"github.com/gin-gonic/gin"
)

// StatsHandler 统计处理器
type StatsHandler struct {
	statsService service.StatsServiceInterface
}

// NewStatsHandler 创建统计处理器实例
func NewStatsHandler(statsService service.StatsServiceInterface) *StatsHandler {
	return &StatsHandler{
		statsService: statsService,
	}
}

// GetBasicStats 获取基础统计
// GET /api/stats/basic
func (h *StatsHandler) GetBasicStats(c *gin.Context) {
	// 获取操作者信息（用于权限控制）
	_, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	stats, err := h.statsService.GetBasicStats()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取统计数据失败")
		return
	}

	response.Success(c, stats)
}

// GetAuthorStats 获取作者统计
// GET /api/stats/author/:id
func (h *StatsHandler) GetAuthorStats(c *gin.Context) {
	// 获取操作者信息（用于权限控制）
	_, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	id := c.Param("id")
	authorID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "作者ID格式错误")
		return
	}

	stats, err := h.statsService.GetStatsByAuthor(uint(authorID))
	if err != nil {
		if err == service.ErrAuthorNotFound {
			response.Error(c, http.StatusNotFound, "作者不存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, "获取统计数据失败")
		return
	}

	response.Success(c, stats)
}

// GetProjectStats 获取课题统计
// GET /api/stats/project/:id
func (h *StatsHandler) GetProjectStats(c *gin.Context) {
	// 获取操作者信息（用于权限控制）
	_, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	id := c.Param("id")
	projectID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "课题ID格式错误")
		return
	}

	stats, err := h.statsService.GetStatsByProject(uint(projectID))
	if err != nil {
		if err == service.ErrProjectNotFound {
			response.Error(c, http.StatusNotFound, "课题不存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, "获取统计数据失败")
		return
	}

	response.Success(c, stats)
}

// GetDepartmentStats 获取单位统计
// GET /api/stats/department
func (h *StatsHandler) GetDepartmentStats(c *gin.Context) {
	// 获取操作者信息（用于权限控制）
	_, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	department := c.Query("name")
	if department == "" {
		response.Error(c, http.StatusBadRequest, "单位名称不能为空")
		return
	}

	stats, err := h.statsService.GetStatsByDepartment(department)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取统计数据失败")
		return
	}

	response.Success(c, stats)
}

// GetYearlyStats 获取年度统计
// GET /api/stats/yearly
func (h *StatsHandler) GetYearlyStats(c *gin.Context) {
	// 获取操作者信息（用于权限控制）
	_, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	stats, err := h.statsService.GetYearlyStats()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取统计数据失败")
		return
	}

	response.Success(c, stats)
}

// GetJournalStats 获取期刊统计
// GET /api/stats/journal
func (h *StatsHandler) GetJournalStats(c *gin.Context) {
	// 获取操作者信息（用于权限控制）
	_, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	stats, err := h.statsService.GetJournalStats()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取统计数据失败")
		return
	}

	response.Success(c, stats)
}
