package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"biolitmanager/internal/model/dto/request"
	"biolitmanager/internal/repository"
	"biolitmanager/internal/security"
	"biolitmanager/internal/service"
	"biolitmanager/pkg/logger"
	"biolitmanager/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ExportHandler 导出处理器
type ExportHandler struct {
	exportService service.ExportServiceInterface
	paperService  *service.PaperService
	paperRepo     *repository.PaperRepository
	statsService  service.StatsServiceInterface
}

// NewExportHandler 创建导出处理器实例
func NewExportHandler(
	exportService service.ExportServiceInterface,
	paperService *service.PaperService,
	paperRepo *repository.PaperRepository,
	statsService service.StatsServiceInterface,
) *ExportHandler {
	return &ExportHandler{
		exportService: exportService,
		paperService:  paperService,
		paperRepo:     paperRepo,
		statsService:  statsService,
	}
}

// ExportPapers 处理查询结果导出请求
// POST /api/export/papers
func (h *ExportHandler) ExportPapers(c *gin.Context) {
	// 获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}
	userIDUint := userID.(uint)

	// 获取用户角色
	role, exists := c.Get("role")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户角色信息异常")
		return
	}
	roleStr, _ := role.(string)

	var req struct {
		SearchRequest request.SearchRequest `json:"search_request"`
		Fields        []string              `json:"fields" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Warn("Invalid export papers request",
			zap.Error(err),
		)
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 校验导出字段
	if len(req.Fields) == 0 {
		response.Error(c, http.StatusBadRequest, "导出字段不能为空")
		return
	}

	// 数据范围过滤：普通用户只能导出自己提交的论文
	// 管理员和超级管理员可以导出所有论文
	if roleStr != string(security.RoleSuperAdmin) && roleStr != string(security.RoleAdmin) {
		// 设置查询条件，只查询当前用户提交的论文
		req.SearchRequest.SubmitterID = &userIDUint
	}

	// 设置分页为最大，以便导出所有符合条件的数据
	if req.SearchRequest.Pagination.PageSize == 0 {
		req.SearchRequest.Pagination.PageSize = 10000
	}
	if req.SearchRequest.Pagination.Page == 0 {
		req.SearchRequest.Pagination.Page = 1
	}

	// 调用仓储查询数据
	papers, total, err := h.paperRepo.AdvancedSearch(&req.SearchRequest)
	if err != nil {
		logger.GetLogger().Error("Failed to search papers for export",
			zap.Error(err),
		)
		response.Error(c, http.StatusInternalServerError, "查询论文失败")
		return
	}

	if total == 0 {
		response.Error(c, http.StatusBadRequest, "没有可导出的数据")
		return
	}

	// 调用导出服务生成Excel
	filePath, err := h.exportService.ExportPapersToExcel(papers, req.Fields)
	if err != nil {
		logger.GetLogger().Error("Failed to export papers",
			zap.Error(err),
		)
		if err == service.ErrExportPermissionDenied {
			response.Error(c, http.StatusForbidden, "导出权限不足")
			return
		}
		response.Error(c, http.StatusInternalServerError, "导出失败")
		return
	}

	// 返回文件路径
	relativePath := strings.Replace(filePath, "uploads/", "", 1)
	response.Success(c, gin.H{
		"file_path": relativePath,
		"file_name": filepath.Base(filePath),
		"count":     len(papers),
	})
}

// ExportPaper 处理单篇论文导出请求
// GET /api/export/paper/:id?format=pdf|word
func (h *ExportHandler) ExportPaper(c *gin.Context) {
	// 获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}
	userIDUint := userID.(uint)

	// 获取用户角色
	role, exists := c.Get("role")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户角色信息异常")
		return
	}
	roleStr, _ := role.(string)

	// 解析论文ID
	id := c.Param("id")
	paperID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "论文ID格式错误")
		return
	}

	// 获取导出格式
	format := c.DefaultQuery("format", "pdf")
	if format != "pdf" && format != "word" {
		response.Error(c, http.StatusBadRequest, "不支持的导出格式，仅支持pdf和word")
		return
	}

	// 数据范围校验：普通用户只能导出自己提交的论文
	if roleStr != string(security.RoleSuperAdmin) && roleStr != string(security.RoleAdmin) {
		paper, err := h.paperService.GetPaperByID(uint(paperID))
		if err != nil {
			if err == service.ErrPaperNotFound {
				response.Error(c, http.StatusNotFound, "论文不存在")
				return
			}
			response.Error(c, http.StatusInternalServerError, "系统异常")
			return
		}
		if paper.SubmitterID != userIDUint {
			response.Error(c, http.StatusForbidden, "权限不足，只能导出自己提交的论文")
			return
		}
	}

	var filePath string
	if format == "pdf" {
		filePath, err = h.exportService.ExportPaperToPDF(uint(paperID), userIDUint)
	} else {
		filePath, err = h.exportService.ExportPaperToWord(uint(paperID), userIDUint)
	}

	if err != nil {
		logger.GetLogger().Error("Failed to export paper",
			zap.Uint("paper_id", uint(paperID)),
			zap.String("format", format),
			zap.Error(err),
		)
		if err == service.ErrPaperNotFound {
			response.Error(c, http.StatusNotFound, "论文不存在")
			return
		}
		if err == service.ErrExportPermissionDenied {
			response.Error(c, http.StatusForbidden, "导出权限不足")
			return
		}
		response.Error(c, http.StatusInternalServerError, "导出失败")
		return
	}

	// 返回文件路径
	relativePath := strings.Replace(filePath, "uploads/", "", 1)
	response.Success(c, gin.H{
		"file_path": relativePath,
		"file_name": filepath.Base(filePath),
		"format":    format,
	})
}

// ExportStats 处理统计结果导出请求
// POST /api/export/stats
func (h *ExportHandler) ExportStats(c *gin.Context) {
	// 获取用户角色
	role, exists := c.Get("role")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "用户角色信息异常")
		return
	}
	roleStr, _ := role.(string)

	// 权限校验：仅管理员可导出统计
	if roleStr != string(security.RoleSuperAdmin) && roleStr != string(security.RoleAdmin) {
		response.Error(c, http.StatusForbidden, "权限不足，仅管理员可导出统计结果")
		return
	}

	var req struct {
		StatsType string `json:"stats_type" binding:"required"` // basic, author, project, department, yearly, journal
		Format    string `json:"format" binding:"required"`     // excel, pdf
		ID        uint   `json:"id"`                            // 用于指定作者/课题/部门的ID
		Title     string `json:"title"`                         // 导出标题
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.GetLogger().Warn("Invalid export stats request",
			zap.Error(err),
		)
		response.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 校验格式
	if req.Format != "excel" && req.Format != "pdf" {
		response.Error(c, http.StatusBadRequest, "不支持的导出格式，仅支持excel和pdf")
		return
	}

	// 获取统计数据
	var stats interface{}
	var title string

	switch req.StatsType {
	case "basic":
		statsData, err := h.statsService.GetBasicStats()
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "获取统计数据失败")
			return
		}
		stats = map[string]interface{}{
			"论文总数":   statsData.TotalPapers,
			"总引用次数":  statsData.TotalCitations,
			"平均影响因子": fmt.Sprintf("%.2f", statsData.AvgImpactFactor),
		}
		title = "基础指标统计"
	case "author":
		if req.ID == 0 {
			response.Error(c, http.StatusBadRequest, "作者ID不能为空")
			return
		}
		statsData, err := h.statsService.GetStatsByAuthor(req.ID)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "获取作者统计数据失败")
			return
		}
		authorName := ""
		if statsData.Author != nil {
			authorName = statsData.Author.Name
		}
		stats = map[string]interface{}{
			"作者姓名":   authorName,
			"论文总数":   statsData.PaperCount,
			"第一作者论文": statsData.FirstAuthorCount,
			"通讯作者论文": statsData.CorrespondingCount,
			"总引用次数":  statsData.TotalCitations,
		}
		title = fmt.Sprintf("作者 %s 统计", authorName)
	case "project":
		if req.ID == 0 {
			response.Error(c, http.StatusBadRequest, "课题ID不能为空")
			return
		}
		statsData, err := h.statsService.GetStatsByProject(req.ID)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "获取课题统计数据失败")
			return
		}
		projectName := ""
		projectCode := ""
		if statsData.Project != nil {
			projectName = statsData.Project.Name
			projectCode = statsData.Project.Code
		}
		stats = map[string]interface{}{
			"课题名称":     projectName,
			"课题编号":     projectCode,
			"论文总数":     statsData.PaperCount,
			"高影响因子论文数": statsData.HighImpactCount,
			"SCI论文数":   statsData.SCIPaperCount,
		}
		title = fmt.Sprintf("课题 %s 统计", projectName)
	case "department":
		if req.ID == 0 {
			response.Error(c, http.StatusBadRequest, "部门ID不能为空")
			return
		}
		statsData, err := h.statsService.GetStatsByDepartment(fmt.Sprintf("%d", req.ID))
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "获取部门统计数据失败")
			return
		}
		stats = map[string]interface{}{
			"部门名称":  statsData.Department,
			"论文总数":  statsData.PaperCount,
			"总引用次数": statsData.TotalCitations,
			"总影响因子": fmt.Sprintf("%.2f", statsData.TotalImpactFactor),
		}
		title = fmt.Sprintf("部门 %s 统计", statsData.Department)
	case "yearly":
		statsData, err := h.statsService.GetYearlyStats()
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "获取年度统计数据失败")
			return
		}
		// 转换为导出格式
		yearlyStats := make([]map[string]interface{}, len(statsData))
		for i, ys := range statsData {
			yearlyStats[i] = map[string]interface{}{
				"年份":   ys.Year,
				"论文数量": ys.Count,
			}
		}
		stats = yearlyStats
		title = "年度统计"
	case "journal":
		statsData, err := h.statsService.GetJournalStats()
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "获取期刊统计数据失败")
			return
		}
		// 转换为导出格式
		journalStats := make([]map[string]interface{}, len(statsData))
		for i, js := range statsData {
			journalStats[i] = map[string]interface{}{
				"期刊名称":   js.JournalName,
				"论文数量":   js.PaperCount,
				"平均影响因子": fmt.Sprintf("%.2f", js.AvgImpactFactor),
			}
		}
		stats = journalStats
		title = "期刊统计"
	default:
		response.Error(c, http.StatusBadRequest, "不支持的统计类型")
		return
	}

	// 使用自定义标题
	if req.Title != "" {
		title = req.Title
	}

	// 调用导出服务
	var filePath string
	var err error
	if req.Format == "excel" {
		filePath, err = h.exportService.ExportStatsToExcel(stats, title)
	} else {
		filePath, err = h.exportService.ExportStatsToPDF(stats, title)
	}

	if err != nil {
		logger.GetLogger().Error("Failed to export stats",
			zap.String("stats_type", req.StatsType),
			zap.String("format", req.Format),
			zap.Error(err),
		)
		response.Error(c, http.StatusInternalServerError, "导出失败")
		return
	}

	// 返回文件路径
	relativePath := strings.Replace(filePath, "uploads/", "", 1)
	response.Success(c, gin.H{
		"file_path": relativePath,
		"file_name": filepath.Base(filePath),
		"format":    req.Format,
	})
}

// DownloadExportFile 处理文件下载
// GET /api/export/download?file=xxx
func (h *ExportHandler) DownloadExportFile(c *gin.Context) {
	// 获取文件路径
	fileName := c.Query("file")
	if fileName == "" {
		response.Error(c, http.StatusBadRequest, "文件路径不能为空")
		return
	}

	// 安全检查：防止路径遍历攻击
	if strings.Contains(fileName, "..") || strings.HasPrefix(fileName, "/") {
		response.Error(c, http.StatusBadRequest, "无效的文件路径")
		return
	}

	// 构建完整文件路径
	filePath := filepath.Join("uploads/exports", fileName)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		logger.GetLogger().Warn("Export file not found",
			zap.String("file_path", filePath),
		)
		response.Error(c, http.StatusNotFound, "文件不存在")
		return
	}

	// 设置响应头
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(filePath)))
	c.Header("Content-Type", "application/octet-stream")

	// 发送文件
	c.File(filePath)
}

// GetExportFields 获取可导出字段列表
// GET /api/export/fields
func (h *ExportHandler) GetExportFields(c *gin.Context) {
	fields := []map[string]string{
		{"key": "title", "name": "标题"},
		{"key": "authors", "name": "作者"},
		{"key": "journal_name", "name": "期刊名称"},
		{"key": "publish_date", "name": "发表日期"},
		{"key": "partition", "name": "收录类型"},
		{"key": "impact_factor", "name": "影响因子"},
		{"key": "citation_count", "name": "引用次数"},
		{"key": "doi", "name": "DOI"},
		{"key": "abstract", "name": "摘要"},
		{"key": "volume", "name": "卷"},
		{"key": "issue", "name": "期"},
		{"key": "start_page", "name": "起始页"},
		{"key": "end_page", "name": "结束页"},
	}

	response.Success(c, fields)
}
