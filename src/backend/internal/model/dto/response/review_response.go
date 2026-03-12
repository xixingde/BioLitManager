package response

import "time"

// ReviewLogDTO 审核记录DTO
type ReviewLogDTO struct {
	ID         uint      `json:"id"`          // 审核记录ID
	PaperID    uint      `json:"paper_id"`    // 论文ID
	ReviewType string    `json:"review_type"` // 审核类型
	Result     string    `json:"result"`      // 审核结果
	Comment    string    `json:"comment"`     // 审核意见
	Reviewer   *UserDTO  `json:"reviewer"`    // 审核人信息
	ReviewTime time.Time `json:"review_time"` // 审核时间
	CreatedAt  time.Time `json:"created_at"`  // 创建时间
}

// PendingReviewDTO 待审核论文DTO
type PendingReviewDTO struct {
	ID              uint      `json:"id"`                // 论文ID
	Title           string    `json:"title"`             // 论文标题
	SubmitterName   string    `json:"submitter_name"`    // 提交人姓名
	SubmitTime      time.Time `json:"submit_time"`       // 提交时间
	Status          string    `json:"status"`            // 状态
	DaysSinceSubmit int       `json:"days_since_submit"` // 距离提交天数
}
