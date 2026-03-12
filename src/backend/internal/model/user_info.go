package model

// UserInfo 用户信息结构体（用于缓存和 JWT Claims）
type UserInfo struct {
	UserID      uint     `json:"user_id"`
	Username    string   `json:"username"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}
