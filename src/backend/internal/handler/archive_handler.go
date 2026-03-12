package handler

import (
	"net/http"
	"strconv"

	"biolitmanager/internal/model/dto/request"
	"biolitmanager/internal/model/entity"
	"biolitmanager/internal/service"
	"biolitmanager/pkg/response"

	"github.com/gin-gonic/gin"
)

// ArchiveHandler 归档处理器
type ArchiveHandler struct {
	archiveService *service.ArchiveService
}

// NewArchiveHandler 创建归档处理器实例
func NewArchiveHandler(archiveService *service.ArchiveService) *ArchiveHandler {
	return &ArchiveHandler{
		archiveService: archiveService,
	}
}

// GetArchiveByPaper 根据论文ID获取归档记录
// GET /api/archives/paper/:paperId
func (h *ArchiveHandler) GetArchiveByPaper(c *gin.Context) {
	id := c.Param("paperId")
	paperID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "论文ID格式错误")
		return
	}

	archive, err := h.archiveService.GetArchiveByPaperID(uint(paperID))
	if err != nil {
		if err == service.ErrArchiveNotFound {
			response.Error(c, http.StatusNotFound, "归档记录不存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	// 转换为DTO格式
	archiveDTO := gin.H{
		"id":             archive.ID,
		"paper_id":       archive.PaperID,
		"archive_number": archive.ArchiveNumber,
		"archive_date":   archive.ArchiveDate,
		"archiver_id":    archive.ArchiverID,
		"created_at":     archive.CreatedAt,
	}

	response.Success(c, archiveDTO)
}

// GetArchiveList 获取归档列表（支持分类筛选）
// GET /api/archives
// 查询参数：year, paperType, author, projectCode, page, pageSize
func (h *ArchiveHandler) GetArchiveList(c *gin.Context) {
	// 解析查询参数
	yearStr := c.Query("year")
	paperType := c.Query("paperType")
	author := c.Query("author")
	projectCode := c.Query("projectCode")
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	// 解析分页参数
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	pagination := &request.Pagination{
		Page:     page,
		PageSize: pageSize,
	}

	var archives []entity.Archive
	var total int64

	// 根据筛选条件调用不同的服务方法
	if yearStr != "" {
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "年份参数格式错误")
			return
		}
		archives, total, err = h.archiveService.GetArchivedPapersByYear(year, pagination)
		if err != nil {
			if err == service.ErrInvalidParameter {
				response.Error(c, http.StatusBadRequest, "无效的年份参数")
				return
			}
			response.Error(c, http.StatusInternalServerError, "系统异常")
			return
		}
	} else if paperType != "" {
		archives, total, err = h.archiveService.GetArchivedPapersByType(paperType, pagination)
		if err != nil {
			if err == service.ErrInvalidParameter {
				response.Error(c, http.StatusBadRequest, "无效的收录类型")
				return
			}
			response.Error(c, http.StatusInternalServerError, "系统异常")
			return
		}
	} else if author != "" {
		archives, total, err = h.archiveService.GetArchivedPapersByAuthor(author, pagination)
		if err != nil {
			if err == service.ErrInvalidParameter {
				response.Error(c, http.StatusBadRequest, "无效的作者参数")
				return
			}
			response.Error(c, http.StatusInternalServerError, "系统异常")
			return
		}
	} else if projectCode != "" {
		archives, total, err = h.archiveService.GetArchivedPapersByProject(projectCode, pagination)
		if err != nil {
			if err == service.ErrInvalidParameter {
				response.Error(c, http.StatusBadRequest, "无效的课题参数")
				return
			}
			response.Error(c, http.StatusInternalServerError, "系统异常")
			return
		}
	} else {
		// 无筛选条件，获取全部归档列表
		archives, total, err = h.archiveService.GetArchivedPapers(pagination)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "系统异常")
			return
		}
	}

	// 转换为DTO格式
	var archiveList []gin.H
	for _, archive := range archives {
		archiveDTO := gin.H{
			"id":             archive.ID,
			"paper_id":       archive.PaperID,
			"archive_number": archive.ArchiveNumber,
			"archive_date":   archive.ArchiveDate,
			"archiver_id":    archive.ArchiverID,
			"status":         archive.Status,
			"is_hidden":      archive.IsHidden,
			"created_at":     archive.CreatedAt,
		}
		archiveList = append(archiveList, archiveDTO)
	}

	response.Success(c, response.PageResult{
		List:  archiveList,
		Total: total,
		Page:  page,
		Size:  pageSize,
	})
}

// HideArchive 隐藏归档论文
// PUT /api/archives/:paperId/hide
func (h *ArchiveHandler) HideArchive(c *gin.Context) {
	id := c.Param("paperId")
	paperID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "论文ID格式错误")
		return
	}

	// 获取操作者信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	ipAddress := c.ClientIP()

	// 先获取归档记录
	archive, err := h.archiveService.GetArchiveByPaperID(uint(paperID))
	if err != nil {
		if err == service.ErrArchiveNotFound {
			response.Error(c, http.StatusNotFound, "归档记录不存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	// 调用服务隐藏归档
	err = h.archiveService.HideArchive(archive.ID, userID.(uint), ipAddress)
	if err != nil {
		if err == service.ErrArchiveNotFound {
			response.Error(c, http.StatusNotFound, "归档记录不存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	response.Success(c, gin.H{
		"message":    "归档已隐藏",
		"paper_id":   paperID,
		"archive_id": archive.ID,
	})
}

// SubmitModifyRequest 提交修改申请
// POST /api/archives/:paperId/modify
func (h *ArchiveHandler) SubmitModifyRequest(c *gin.Context) {
	id := c.Param("paperId")
	paperID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "论文ID格式错误")
		return
	}

	// 绑定请求参数
	var req struct {
		RequestType   string                 `json:"request_type" binding:"required"`
		RequestReason string                 `json:"request_reason" binding:"required"`
		RequestData   map[string]interface{} `json:"request_data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	// 验证修改类型
	validTypes := map[string]bool{"update": true, "delete": true, "hide": true}
	if !validTypes[req.RequestType] {
		response.Error(c, http.StatusBadRequest, "无效的修改申请类型")
		return
	}

	// 获取操作者信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	ipAddress := c.ClientIP()

	// 先获取归档记录
	archive, err := h.archiveService.GetArchiveByPaperID(uint(paperID))
	if err != nil {
		if err == service.ErrArchiveNotFound {
			response.Error(c, http.StatusNotFound, "归档记录不存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	// 调用服务提交修改申请
	modifyRequest, err := h.archiveService.SubmitArchiveModifyRequest(
		archive.ID,
		req.RequestType,
		req.RequestReason,
		req.RequestData,
		userID.(uint),
		ipAddress,
	)
	if err != nil {
		if err == service.ErrInvalidModifyRequestType {
			response.Error(c, http.StatusBadRequest, "无效的修改申请类型")
			return
		}
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	response.Success(c, gin.H{
		"message":    "修改申请已提交",
		"paper_id":   paperID,
		"archive_id": archive.ID,
		"modify_request": gin.H{
			"id":             modifyRequest.ID,
			"request_type":   modifyRequest.RequestType,
			"request_reason": modifyRequest.RequestReason,
			"status":         modifyRequest.Status,
			"created_at":     modifyRequest.CreatedAt,
		},
	})
}
