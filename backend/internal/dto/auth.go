package dto

type LoginRequest struct {
	Username string `json:"username" binding:"required" validate:"required,min=3"`
	Password string `json:"password" binding:"required" validate:"required,min=6"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  UserResponse `json:"user"`
}

