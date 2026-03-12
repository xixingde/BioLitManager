package request

// BusinessReviewRequest 业务审核请求
type BusinessReviewRequest struct {
	Result  string `json:"result" binding:"required"` // 审核结果：通过/驳回
	Comment string `json:"comment"`                   // 审核意见
}

// PoliticalReviewRequest 政治审核请求
type PoliticalReviewRequest struct {
	Result  string `json:"result" binding:"required"` // 审核结果：通过/驳回
	Comment string `json:"comment"`                   // 审核意见
}
