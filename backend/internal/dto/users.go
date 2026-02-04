package dto

type CreateUserRequest struct {
	Username string `json:"username" binding:"required" validate:"required,min=3"`
	Password string `json:"password" binding:"required" validate:"required,min=6"`
}

type UpdateUserRequest struct {
	Username *string `json:"username,omitempty" validate:"omitempty,min=3"`
	Bio      *string `json:"bio,omitempty"`
	Avatar   *string `json:"avatar,omitempty"`
}

type UserResponse struct {
	ID        string   `json:"id"`
	Username  string   `json:"username"`
	Avatar    *string  `json:"avatar,omitempty"`
	Bio       *string  `json:"bio,omitempty"`
	Following []string `json:"following"`
}

