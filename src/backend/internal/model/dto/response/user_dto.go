package response

// UserDTO 用户信息DTO
type UserDTO struct {
	ID         uint   `json:"id"`
	Username   string `json:"username"`
	Name       string `json:"name"`
	Role       string `json:"role"`
	Department string `json:"department"`
	Email      string `json:"email"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}
