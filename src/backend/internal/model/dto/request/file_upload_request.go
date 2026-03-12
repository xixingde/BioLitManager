package request

// UploadFileRequest 文件上传请求
type UploadFileRequest struct {
	PaperID  uint        `json:"paper_id" binding:"required"`  // 论文ID
	FileType string      `json:"file_type" binding:"required"` // 文件类型
	File     interface{} `json:"file" binding:"required"`      // 文件对象
}
