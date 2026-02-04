package dto

import "time"

type CreatePostRequest struct {
	Content string `json:"content" binding:"required" validate:"required"`
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required" validate:"required"`
}

type AuthorResponse struct {
	ID       string  `json:"id"`
	Username string  `json:"username"`
	Avatar   *string `json:"avatar,omitempty"`
}

type CommentResponse struct {
	ID        string         `json:"id"`
	UserID    string         `json:"userId"`
	Content   string         `json:"content"`
	CreatedAt time.Time      `json:"createdAt"`
	Author    AuthorResponse `json:"author"`
}

type PostResponse struct {
	ID               string           `json:"id"`
	UserID           string           `json:"userId"`
	Author           AuthorResponse   `json:"author"`
	Content          string           `json:"content"`
	Likes            int              `json:"likes"`
	LikedByCurrentUser bool           `json:"likedByCurrentUser"`
	CreatedAt        time.Time        `json:"createdAt"`
	Comments         []CommentResponse `json:"comments"`
}

