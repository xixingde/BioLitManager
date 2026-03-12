package response

import "time"

// JournalDTO 期刊信息DTO
type JournalDTO struct {
	ID           uint    `json:"id"`            // 期刊ID
	FullName     string  `json:"full_name"`     // 期刊全称
	ShortName    string  `json:"short_name"`    // 期刊简称
	ISSN         string  `json:"issn"`          // ISSN号
	ImpactFactor float64 `json:"impact_factor"` // 影响因子
	Publisher    string  `json:"publisher"`     // 出版商
}

// AuthorDTO 作者信息DTO
type AuthorDTO struct {
	ID         uint   `json:"id"`          // 作者ID
	Name       string `json:"name"`        // 作者姓名
	AuthorType string `json:"author_type"` // 作者类型
	Rank       int    `json:"rank"`        // 作者排序
	Department string `json:"department"`  // 所在部门
	UserID     *uint  `json:"user_id"`     // 关联用户ID
}

// ProjectDTO 课题信息DTO
type ProjectDTO struct {
	ID          uint   `json:"id"`           // 课题ID
	Name        string `json:"name"`         // 课题名称
	Code        string `json:"code"`         // 课题编号
	ProjectType string `json:"project_type"` // 课题类型
	Source      string `json:"source"`       // 课题来源
	Level       string `json:"level"`        // 课题级别
}

// AttachmentDTO 附件信息DTO
type AttachmentDTO struct {
	ID        uint      `json:"id"`         // 附件ID
	FileType  string    `json:"file_type"`  // 文件类型
	FileName  string    `json:"file_name"`  // 文件名
	FilePath  string    `json:"file_path"`  // 文件路径
	FileSize  int64     `json:"file_size"`  // 文件大小
	MimeType  string    `json:"mime_type"`  // MIME类型
	CreatedAt time.Time `json:"created_at"` // 上传时间
}

// PaperDTO 论文信息DTO
type PaperDTO struct {
	ID           uint            `json:"id"`            // 论文ID
	Title        string          `json:"title"`         // 论文标题
	Abstract     string          `json:"abstract"`      // 摘要
	Journal      *JournalDTO     `json:"journal"`       // 期刊信息
	DOI          string          `json:"doi"`           // DOI号
	ImpactFactor float64         `json:"impact_factor"` // 影响因子
	PublishDate  *time.Time      `json:"publish_date"`  // 发表日期
	Status       string          `json:"status"`        // 论文状态
	Submitter    *UserDTO        `json:"submitter"`     // 提交人信息
	SubmitTime   time.Time       `json:"submit_time"`   // 提交时间
	Authors      []AuthorDTO     `json:"authors"`       // 作者列表
	Projects     []ProjectDTO    `json:"projects"`      // 课题列表
	Attachments  []AttachmentDTO `json:"attachments"`   // 附件列表
	CreatedAt    time.Time       `json:"created_at"`    // 创建时间
	UpdatedAt    time.Time       `json:"updated_at"`    // 更新时间
}

// PaperListResponse 论文列表响应
type PaperListResponse struct {
	List  []PaperDTO `json:"list"`  // 论文列表
	Total int64      `json:"total"` // 总数
	Page  int        `json:"page"`  // 当前页
	Size  int        `json:"size"`  // 每页大小
}
