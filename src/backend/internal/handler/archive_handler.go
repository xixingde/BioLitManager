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

// ListArchives 获取归档记录列表
// GET /api/archives
func (h *ArchiveHandler) ListArchives(c *gin.Context) {
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

	// TODO: 实现分页查询归档记录列表
	// 这里先返回空列表，需要在ArchiveService中添加ListArchives方法
	logger.GetLogger().Info("List archives called",
		zap.Int("page", page),
		zap.Int("size", size),
	)

	response.Success(c, gin.H{
		"list":  []interface{}{},
		"total": 0,
		"page":  page,
		"size":  size,
	})
}
