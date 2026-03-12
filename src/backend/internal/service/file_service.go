package service

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"biolitmanager/internal/model/entity"
	"biolitmanager/internal/repository"
	"biolitmanager/pkg/logger"

	"go.uber.org/zap"
)

var (
	// ErrFileTooLarge 文件过大
	ErrFileTooLarge = errors.New("文件过大")
	// ErrInvalidFileType 无效的文件类型
	ErrInvalidFileType = errors.New("无效的文件类型")
	// ErrFileNotFound 文件不存在
	ErrFileNotFound = errors.New("文件不存在")
	// ErrUploadFailed 上传失败
	ErrUploadFailed = errors.New("上传失败")
)

const (
	// MaxFileSize 最大文件大小：50MB
	MaxFileSize = 50 * 1024 * 1024
	// UploadDir 上传目录
	UploadDir = "./uploads"
)

// FileService 文件服务
type FileService struct {
	attachmentRepo *repository.AttachmentRepository
}

// NewFileService 创建文件服务实例
func NewFileService(attachmentRepo *repository.AttachmentRepository) *FileService {
	return &FileService{
		attachmentRepo: attachmentRepo,
	}
}

// UploadFile 文件上传（校验大小、格式，保存文件）
func (s *FileService) UploadFile(
	paperID uint,
	fileType string,
	file *multipart.FileHeader,
	operatorID uint,
	ipAddress string,
) (*entity.Attachment, error) {
	// 校验文件大小
	if file.Size > MaxFileSize {
		logger.GetLogger().Warn("File size exceeds limit",
			zap.Uint("paper_id", paperID),
			zap.String("filename", file.Filename),
			zap.Int64("size", file.Size),
			zap.Int64("max_size", MaxFileSize),
		)
		return nil, ErrFileTooLarge
	}

	// 校验文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !s.isValidFileType(ext) {
		logger.GetLogger().Warn("Invalid file type",
			zap.Uint("paper_id", paperID),
			zap.String("filename", file.Filename),
			zap.String("ext", ext),
		)
		return nil, ErrInvalidFileType
	}

	// 确保上传目录存在
	if err := os.MkdirAll(UploadDir, 0755); err != nil {
		logger.GetLogger().Error("Failed to create upload directory",
			zap.String("dir", UploadDir),
			zap.Error(err),
		)
		return nil, ErrUploadFailed
	}

	// 生成唯一文件名
	uniqueFileName := fmt.Sprintf("%d_%d%s", paperID, file.Size, ext)
	filePath := filepath.Join(UploadDir, uniqueFileName)

	// 打开源文件
	src, err := file.Open()
	if err != nil {
		logger.GetLogger().Error("Failed to open uploaded file",
			zap.String("filename", file.Filename),
			zap.Error(err),
		)
		return nil, ErrUploadFailed
	}
	defer src.Close()

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		logger.GetLogger().Error("Failed to create file",
			zap.String("path", filePath),
			zap.Error(err),
		)
		return nil, ErrUploadFailed
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, src); err != nil {
		logger.GetLogger().Error("Failed to copy file",
			zap.String("filename", file.Filename),
			zap.String("path", filePath),
			zap.Error(err),
		)
		// 删除已创建的文件
		os.Remove(filePath)
		return nil, ErrUploadFailed
	}

	// 创建附件记录
	attachment := &entity.Attachment{
		PaperID:    paperID,
		FileName:   file.Filename,
		FileType:   fileType,
		FilePath:   filePath,
		FileSize:   file.Size,
		UploaderID: operatorID,
	}

	if err := s.attachmentRepo.Create(attachment); err != nil {
		logger.GetLogger().Error("Failed to create attachment record",
			zap.Uint("paper_id", paperID),
			zap.String("filename", file.Filename),
			zap.Error(err),
		)
		// 删除已上传的文件
		os.Remove(filePath)
		return nil, ErrUploadFailed
	}

	logger.GetLogger().Info("File uploaded successfully",
		zap.Uint("attachment_id", attachment.ID),
		zap.Uint("paper_id", paperID),
		zap.String("filename", file.Filename),
		zap.Uint("operator_id", operatorID),
	)

	return attachment, nil
}

// GetFileByID 获取附件信息
func (s *FileService) GetFileByID(id uint) (*entity.Attachment, error) {
	attachment, err := s.attachmentRepo.FindByID(id)
	if err != nil {
		logger.GetLogger().Error("Failed to find attachment",
			zap.Uint("attachment_id", id),
			zap.Error(err),
		)
		return nil, ErrSystemError
	}

	if attachment == nil {
		return nil, ErrFileNotFound
	}

	return attachment, nil
}

// DeleteFile 删除文件
func (s *FileService) DeleteFile(id uint, operatorID uint, ipAddress string) error {
	// 查询附件
	attachment, err := s.attachmentRepo.FindByID(id)
	if err != nil {
		logger.GetLogger().Error("Failed to find attachment",
			zap.Uint("attachment_id", id),
			zap.Error(err),
		)
		return ErrSystemError
	}

	if attachment == nil {
		return ErrFileNotFound
	}

	// 删除物理文件
	if err := os.Remove(attachment.FilePath); err != nil {
		// 如果文件不存在，忽略错误
		if !os.IsNotExist(err) {
			logger.GetLogger().Error("Failed to delete file",
				zap.String("path", attachment.FilePath),
				zap.Error(err),
			)
			return ErrSystemError
		}
	}

	// 删除数据库记录
	if err := s.attachmentRepo.Delete(id); err != nil {
		logger.GetLogger().Error("Failed to delete attachment record",
			zap.Uint("attachment_id", id),
			zap.Error(err),
		)
		return ErrSystemError
	}

	logger.GetLogger().Info("File deleted successfully",
		zap.Uint("attachment_id", id),
		zap.String("filename", attachment.FileName),
		zap.Uint("operator_id", operatorID),
	)

	return nil
}

// isValidFileType 校验文件类型是否合法
func (s *FileService) isValidFileType(ext string) bool {
	// 允许的文件类型
	allowedExtensions := map[string]bool{
		".pdf":  true,
		".doc":  true,
		".docx": true,
		".xls":  true,
		".xlsx": true,
		".ppt":  true,
		".pptx": true,
		".zip":  true,
		".rar":  true,
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	return allowedExtensions[ext]
}
