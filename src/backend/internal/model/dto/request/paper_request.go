package request

// CreatePaperRequest 创建论文请求
type CreatePaperRequest struct {
	Title        string                `json:"title" binding:"required"`      // 论文标题
	Abstract     string                `json:"abstract"`                      // 摘要
	JournalID    uint                  `json:"journal_id" binding:"required"` // 期刊ID
	DOI          string                `json:"doi"`                           // DOI号
	ImpactFactor float64               `json:"impact_factor"`                 // 影响因子
	PublishDate  string                `json:"publish_date"`                  // 发表日期
	Authors      []CreateAuthorRequest `json:"authors"`                       // 作者列表
	Projects     []uint                `json:"projects"`                      // 课题ID列表
	Attachments  []UploadFileRequest   `json:"attachments"`                   // 附件信息
}

// UpdatePaperRequest 更新论文请求
type UpdatePaperRequest struct {
	Title        string                `json:"title"`         // 论文标题
	Abstract     string                `json:"abstract"`      // 摘要
	JournalID    uint                  `json:"journal_id"`    // 期刊ID
	DOI          string                `json:"doi"`           // DOI号
	ImpactFactor float64               `json:"impact_factor"` // 影响因子
	PublishDate  string                `json:"publish_date"`  // 发表日期
	Authors      []CreateAuthorRequest `json:"authors"`       // 作者列表
	Projects     []uint                `json:"projects"`      // 课题ID列表
}

// SubmitForReviewRequest 提交审核请求
type SubmitForReviewRequest struct {
	PaperID uint `json:"paper_id" binding:"required"` // 论文ID
}

// SaveDraftRequest 保存草稿请求
type SaveDraftRequest struct {
	PaperID uint `json:"paper_id" binding:"required"` // 论文ID
}

// CheckDuplicateRequest 查重请求
type CheckDuplicateRequest struct {
	Title string `json:"title"` // 论文标题
	DOI   string `json:"doi"`   // DOI号
}

// BatchImportRequest 批量导入请求
type BatchImportRequest struct {
	File interface{} `json:"file" binding:"required"` // Excel文件
}
