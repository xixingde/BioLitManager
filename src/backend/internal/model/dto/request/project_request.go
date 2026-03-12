package request

// CreateProjectRequest 创建课题请求
type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required"`         // 课题名称
	Code        string `json:"code" binding:"required"`         // 课题编号
	ProjectType string `json:"project_type" binding:"required"` // 课题类型
	Source      string `json:"source"`                          // 课题来源
	Level       string `json:"level"`                           // 课题级别
}

// UpdateProjectRequest 更新课题请求
type UpdateProjectRequest struct {
	Name        string `json:"name"`         // 课题名称
	Code        string `json:"code"`         // 课题编号
	ProjectType string `json:"project_type"` // 课题类型
	Source      string `json:"source"`       // 课题来源
	Level       string `json:"level"`        // 课题级别
	Status      string `json:"status"`       // 课题状态
}
