package request

// LogicType 逻辑类型枚举
type LogicType string

const (
	LogicAnd LogicType = "AND" // 且
	LogicOr  LogicType = "OR"  // 或
	LogicNot LogicType = "NOT" // 非
)

// AuthorType 作者类型枚举
type AuthorType string

const (
	AuthorTypeFirst         AuthorType = "first"         // 第一作者
	AuthorTypeCoFirst       AuthorType = "co_first"      // 共同第一作者
	AuthorTypeCorresponding AuthorType = "corresponding" // 通讯作者
)

// SortField 排序字段枚举
type SortField string

const (
	SortFieldPublishDate  SortField = "publish_date"  // 出版日期
	SortFieldImpactFactor SortField = "impact_factor" // 影响因子
	SortFieldCitation     SortField = "citation"      // 引用次数
)

// SortOrder 排序方向枚举
type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"  // 升序
	SortOrderDesc SortOrder = "desc" // 降序
)

// QueryCondition 单个查询条件
type QueryCondition struct {
	Field    string    `json:"field"`    // 查询字段
	Value    string    `json:"value"`    // 查询值
	Logic    LogicType `json:"logic"`    // 逻辑类型
	Operator string    `json:"operator"` // 操作符（eq, like, gt, lt, gte, lte, in）
}

// QueryGroup 查询条件组，支持嵌套组合
type QueryGroup struct {
	Logic      LogicType        `json:"logic"`      // 逻辑类型
	Conditions []QueryCondition `json:"conditions"` // 条件列表
	Groups     []QueryGroup     `json:"groups"`     // 嵌套条件组
}

// Pagination 分页参数
type Pagination struct {
	Page     int `json:"page"`      // 当前页码，默认1
	PageSize int `json:"page_size"` // 每页条数，默认20
}

// SearchRequest 搜索请求
type SearchRequest struct {
	// 基础查询条件
	Keywords    string `json:"keywords"`     // 关键词（标题、摘要）
	Title       string `json:"title"`        // 标题
	Abstract    string `json:"abstract"`     // 摘要
	DOI         string `json:"doi"`          // DOI号
	JournalID   uint   `json:"journal_id"`   // 期刊ID
	JournalName string `json:"journal_name"` // 期刊名称

	// 作者相关
	AuthorName string     `json:"author_name"` // 作者姓名
	AuthorType AuthorType `json:"author_type"` // 作者类型
	Department string     `json:"department"`  // 所在部门

	// 项目相关
	ProjectCode string `json:"project_code"` // 课题编号
	ProjectType string `json:"project_type"` // 课题类型

	// 时间范围
	PublishDateStart string `json:"publish_date_start"` // 发表日期开始
	PublishDateEnd   string `json:"publish_date_end"`   // 发表日期结束

	// 影响因子范围
	ImpactFactorMin float64 `json:"impact_factor_min"` // 影响因子最小值
	ImpactFactorMax float64 `json:"impact_factor_max"` // 影响因子最大值

	// 状态筛选
	Status string `json:"status"` // 论文状态

	// 年份和类型筛选
	Year      int    `json:"year"`       // 发表年份
	PaperType string `json:"paper_type"` // 收录类型（SCI/EI/CI/DI/CORE）

	// 提交人筛选（用于数据范围过滤）
	SubmitterID *uint `json:"submitter_id"` // 提交人ID

	// 高级查询
	QueryGroup *QueryGroup `json:"query_group"` // 复杂查询条件组

	// 排序
	SortField SortField `json:"sort_field"` // 排序字段
	SortOrder SortOrder `json:"sort_order"` // 排序方向

	// 分页
	Pagination Pagination `json:"pagination"` // 分页参数
}
