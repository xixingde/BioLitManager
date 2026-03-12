package request

// CreateAuthorRequest 创建作者请求
type CreateAuthorRequest struct {
	PaperID    uint   `json:"paper_id"`                       // 论文ID
	Name       string `json:"name" binding:"required"`        // 作者姓名
	AuthorType string `json:"author_type" binding:"required"` // 作者类型
	Rank       int    `json:"rank" binding:"required"`        // 作者排序
	Department string `json:"department"`                     // 所在部门
	UserID     *uint  `json:"user_id"`                        // 关联用户ID
}

// UpdateAuthorRequest 更新作者请求
type UpdateAuthorRequest struct {
	Name       string `json:"name"`        // 作者姓名
	AuthorType string `json:"author_type"` // 作者类型
	Rank       int    `json:"rank"`        // 作者排序
	Department string `json:"department"`  // 所在部门
	UserID     *uint  `json:"user_id"`     // 关联用户ID
}
