package handler

import (
	"net/http"
	"os"
	"strconv"

	"biolitmanager/internal/service"
	"biolitmanager/pkg/logger"
	"biolitmanager/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// FileHandler 文件处理器
type FileHandler struct {
	fileService *service.FileService
}

// NewFileHandler 创建文件处理器实例
func NewFileHandler(fileService *service.FileService) *FileHandler {
	return &FileHandler{
		fileService: fileService,
	}
}

// UploadFile 上传文件
// POST /api/files/upload
func (h *FileHandler) UploadFile(c *gin.Context) {
	// 获取文件
	fileHeader, err := c.FormFile("file")
	if err != nil {
		logger.GetLogger().Warn("Failed to get file from request",
			zap.Error(err),
		)
		response.Error(c, http.StatusBadRequest, "文件上传失败")
		return
	}

	// 获取论文ID
	paperIDStr := c.PostForm("paper_id")
	paperID, err := strconv.ParseUint(paperIDStr, 10, 32)
	if err != nil || paperID == 0 {
		response.Error(c, http.StatusBadRequest, "论文ID格式错误")
		return
	}

	// 获取文件类型
	fileType := c.PostForm("file_type")
	if fileType == "" {
		response.Error(c, http.StatusBadRequest, "文件类型不能为空")
		return
	}

	// 获取上传人信息
	uploaderID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	ipAddress := c.ClientIP()

	// 调用Service上传文件
	attachment, err := h.fileService.UploadFile(
		uint(paperID),
		fileType,
		fileHeader,
		uploaderID.(uint),
		ipAddress,
	)

	if err != nil {
		if err == service.ErrFileTooLarge {
			response.Error(c, http.StatusBadRequest, "文件过大")
			return
		}
		if err == service.ErrInvalidFileType {
			response.Error(c, http.StatusBadRequest, "文件格式错误")
			return
		}
		response.Error(c, http.StatusInternalServerError, "文件上传失败")
		return
	}

	response.Success(c, gin.H{
		"id":        attachment.ID,
		"file_name": attachment.FileName,
		"file_size": attachment.FileSize,
		"file_type": attachment.FileType,
	})
}

// GetFile 获取文件信息
// GET /api/files/:id
func (h *FileHandler) GetFile(c *gin.Context) {
	id := c.Param("id")
	fileID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "文件ID格式错误")
		return
	}

	attachment, err := h.fileService.GetFileByID(uint(fileID))
	if err != nil {
		if err == service.ErrFileNotFound {
			response.Error(c, http.StatusNotFound, "文件不存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	// 转换为DTO格式
	fileDTO := gin.H{
		"id":         attachment.ID,
		"file_type":  attachment.FileType,
		"file_name":  attachment.FileName,
		"file_size":  attachment.FileSize,
		"mime_type":  attachment.MimeType,
		"created_at": attachment.CreatedAt,
	}

	response.Success(c, fileDTO)
}

// DownloadFile 下载文件
// GET /api/files/:id/download
func (h *FileHandler) DownloadFile(c *gin.Context) {
	id := c.Param("id")
	fileID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "文件ID格式错误")
		return
	}

	attachment, err := h.fileService.GetFileByID(uint(fileID))
	if err != nil {
		if err == service.ErrFileNotFound {
			response.Error(c, http.StatusNotFound, "文件不存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(attachment.FilePath); os.IsNotExist(err) {
		logger.GetLogger().Error("File not found on disk",
			zap.String("file_path", attachment.FilePath),
			zap.Error(err),
		)
		response.Error(c, http.StatusNotFound, "文件不存在")
		return
	}

	// 读取文件内容
	fileData, err := os.ReadFile(attachment.FilePath)
	if err != nil {
		logger.GetLogger().Error("Failed to read file",
			zap.String("file_path", attachment.FilePath),
			zap.Error(err),
		)
		response.Error(c, http.StatusInternalServerError, "读取文件失败")
		return
	}

	// 设置响应头
	c.Header("Content-Type", attachment.MimeType)
	c.Header("Content-Disposition", "attachment; filename=\""+attachment.FileName+"\"")
	c.Header("Content-Length", strconv.Itoa(len(fileData)))

	// 返回文件内容
	c.Data(http.StatusOK, attachment.MimeType, fileData)
}

// DeleteFile 删除文件
// DELETE /api/files/:id
func (h *FileHandler) DeleteFile(c *gin.Context) {
	id := c.Param("id")
	fileID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "文件ID格式错误")
		return
	}

	// 获取操作者信息
	operatorID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "未授权访问")
		return
	}

	ipAddress := c.ClientIP()

	if err := h.fileService.DeleteFile(uint(fileID), operatorID.(uint), ipAddress); err != nil {
		if err == service.ErrFileNotFound {
			response.Error(c, http.StatusNotFound, "文件不存在")
			return
		}
		response.Error(c, http.StatusInternalServerError, "系统异常")
		return
	}

	response.Success(c, nil)
}
