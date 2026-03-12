package response

import "time"

// PaginationInfo 分页信息
type PaginationInfo struct {
	Page       int   `json:"page"`        // 当前页
	PageSize   int   `json:"page_size"`   // 每页条数
	TotalCount int64 `json:"total_count"` // 总条数
	TotalPages int   `json:"total_pages"` // 总页数
}

// AuthorInfo 作者信息
type AuthorInfo struct {
	Name       string `json:"name"`        // 作者姓名
	Department string `json:"department"`  // 所在部门
	AuthorType string `json:"author_type"` // 作者类型（first/co_first/corresponding）
	Rank       int    `json:"rank"`        // 作者排序
}

// ProjectInfo 课题信息
type ProjectInfo struct {
	Code        string `json:"code"`         // 课题编号
	ProjectType string `json:"project_type"` // 项目类型
	Name        string `json:"name"`         // 课题名称
}

// PaperSearchResult 论文搜索结果项
type PaperSearchResult struct {
	ID           uint          `json:"id"`            // 论文ID
	Title        string        `json:"title"`         // 论文标题
	Abstract     string        `json:"abstract"`      // 摘要
	JournalName  string        `json:"journal_name"`  // 期刊名称
	JournalShort string        `json:"journal_short"` // 期刊简称
	DOI          string        `json:"doi"`           // DOI号
	ImpactFactor float64       `json:"impact_factor"` // 影响因子
	PublishDate  *time.Time    `json:"publish_date"`  // 出版日期
	Partition    string        `json:"partition"`     // 分区
	Status       string        `json:"status"`        // 审核状态
	Authors      []AuthorInfo  `json:"authors"`       // 作者列表
	Projects     []ProjectInfo `json:"projects"`      // 课题列表
	Citation     int           `json:"citation"`      // 引用次数
	CreatedAt    time.Time     `json:"created_at"`    // 创建时间
}

// SearchResponse 搜索响应
type SearchResponse struct {
	Pagination PaginationInfo      `json:"pagination"` // 分页信息
	Results    []PaperSearchResult `json:"results"`    // 搜索结果列表
}

// ApiResponse 统一响应格式
type ApiResponse struct {
	Code    int         `json:"code"`    // 状态码
	Message string      `json:"message"` // 消息
	Data    interface{} `json:"data"`    // 数据
	Success bool        `json:"success"` // 是否成功
}

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(data interface{}) ApiResponse {
	return ApiResponse{
		Code:    200,
		Message: "success",
		Data:    data,
		Success: true,
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(code int, message string) ApiResponse {
	return ApiResponse{
		Code:    code,
		Message: message,
		Data:    nil,
		Success: false,
	}
}
