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

// JournalHandler 期刊处理器
type JournalHandler struct {
	journalService *service.JournalService
}

// NewJournalHandler 创建期刊处理器实例
func NewJournalHandler(journalService *service.JournalService) *JournalHandler {
	return &JournalHandler{
		journalService: journalService,
	}
}

// CreateJournal 创建期刊
// POST /api/journals
func (h *JournalHandler) CreateJournal(c *gin.Context) {
	var req struct {
		FullName     string  `json:"full_name" binding:"required"`
		ShortName    string  `json:"short_name"`
		ISSN         string  `json:"issn" binding:"required"`
		ImpactFactor float64 `json:"impact_factor"`
		Publisher    string  `json:"publisher"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Warn("Invalid create journal request",
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

	journal, err := h.journalService.CreateJournal(
		req.FullName,
		req.ShortName,
		req.ISSN,
		req.ImpactFactor,
		req.Publisher,
		operatorID.(uint),
		ipAddress,
	)

	if err != nil {
		if err == service.ErrISSNExists {
			response.Error(c, http.StatusBadRequest, "ISSN已存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	response.Success(c, gin.H{
		"id": journal.ID,
	})
}

// GetJournal 获取期刊详情
// GET /api/journals/:id
func (h *JournalHandler) GetJournal(c *gin.Context) {
	id := c.Param("id")
	journalID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "期刊ID格式错误")
		return
	}

	journal, err := h.journalService.GetJournalByID(uint(journalID))
	if err != nil {
		if err == service.ErrJournalNotFound {
			response.Error(c, http.StatusNotFound, "期刊不存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	// 转换为DTO格式
	journalDTO := gin.H{
		"id":            journal.ID,
		"full_name":     journal.FullName,
		"short_name":    journal.ShortName,
		"issn":          journal.ISSN,
		"impact_factor": journal.ImpactFactor,
		"publisher":     journal.Publisher,
		"created_at":    journal.CreatedAt,
		"updated_at":    journal.UpdatedAt,
	}

	response.Success(c, journalDTO)
}

// ListJournals 分页查询期刊列表
// GET /api/journals
func (h *JournalHandler) ListJournals(c *gin.Context) {
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

	journals, total, err := h.journalService.ListJournals(page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	// 转换为DTO格式
	journalList := make([]gin.H, len(journals))
	for i, journal := range journals {
		journalList[i] = gin.H{
			"id":            journal.ID,
			"full_name":     journal.FullName,
			"short_name":    journal.ShortName,
			"issn":          journal.ISSN,
			"impact_factor": journal.ImpactFactor,
			"publisher":     journal.Publisher,
			"created_at":    journal.CreatedAt,
			"updated_at":    journal.UpdatedAt,
		}
	}

	result := gin.H{
		"list":  journalList,
		"total": total,
		"page":  page,
		"size":  size,
	}

	response.Success(c, result)
}

// UpdateJournal 更新期刊
// PUT /api/journals/:id
func (h *JournalHandler) UpdateJournal(c *gin.Context) {
	id := c.Param("id")
	journalID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "期刊ID格式错误")
		return
	}

	var req struct {
		FullName     string  `json:"full_name"`
		ShortName    string  `json:"short_name"`
		ISSN         string  `json:"issn"`
		ImpactFactor float64 `json:"impact_factor"`
		Publisher    string  `json:"publisher"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Warn("Invalid update journal request",
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

	if err := h.journalService.UpdateJournal(
		uint(journalID),
		req.FullName,
		req.ShortName,
		req.ISSN,
		req.ImpactFactor,
		req.Publisher,
		operatorID.(uint),
		ipAddress,
	); err != nil {
		if err == service.ErrJournalNotFound {
			response.Error(c, http.StatusNotFound, "期刊不存在")
			return
		}
		if err == service.ErrISSNExists {
			response.Error(c, http.StatusBadRequest, "ISSN已存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	response.Success(c, nil)
}

// UpdateImpactFactor 更新期刊影响因子
// PUT /api/journals/:id/impact-factor
func (h *JournalHandler) UpdateImpactFactor(c *gin.Context) {
	id := c.Param("id")
	journalID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "期刊ID格式错误")
		return
	}

	var req struct {
		ImpactFactor float64 `json:"impact_factor" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Warn("Invalid update impact factor request",
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

	if err := h.journalService.UpdateImpactFactor(uint(journalID), req.ImpactFactor, operatorID.(uint), ipAddress); err != nil {
		if err == service.ErrJournalNotFound {
			response.Error(c, http.StatusNotFound, "期刊不存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	response.Success(c, nil)
}

// SearchJournals 搜索期刊
// GET /api/journals/search
func (h *JournalHandler) SearchJournals(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		response.Error(c, http.StatusBadRequest, "搜索关键字不能为空")
		return
	}

	journals, err := h.journalService.SearchJournals(keyword)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	// 转换为DTO格式
	journalList := make([]gin.H, len(journals))
	for i, journal := range journals {
		journalList[i] = gin.H{
			"id":            journal.ID,
			"full_name":     journal.FullName,
			"short_name":    journal.ShortName,
			"issn":          journal.ISSN,
			"impact_factor": journal.ImpactFactor,
			"publisher":     journal.Publisher,
		}
	}

	response.Success(c, gin.H{
		"list": journalList,
	})
}
