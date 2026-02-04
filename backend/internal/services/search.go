package services

import (
	"errors"
	"fmt"
	"strings"

	"tinker-backend/internal/database"
	"tinker-backend/internal/dto"

	"gorm.io/gorm"
)

type SearchService struct {
	db *gorm.DB
}

func NewSearchService(db *gorm.DB) *SearchService {
	return &SearchService{
		db: db,
	}
}

func (s *SearchService) Search(query, searchType string) (interface{}, error) {
	if searchType != "users" && searchType != "posts" {
		return nil, errors.New("invalid type")
	}

	query = strings.TrimSpace(query)
	if query == "" {
		return nil, errors.New("query cannot be empty")
	}

	searchPattern := "%" + strings.ToLower(query) + "%"

	if searchType == "users" {
		var users []database.User
		if err := s.db.Where("LOWER(username) LIKE ?", searchPattern).Find(&users).Error; err != nil {
			return nil, fmt.Errorf("failed to search users: %w", err)
		}

		responses := make([]dto.UserResponse, len(users))
		for i, user := range users {
			followingIDs, err := s.getFollowingIDs(user.ID)
			if err != nil {
				return nil, err
			}
			responses[i] = dto.UserResponse{
				ID:        user.ID,
				Username:  user.Username,
				Avatar:    user.Avatar,
				Bio:       user.Bio,
				Following: followingIDs,
			}
		}

		// Return array directly to match NestJS format
		return responses, nil
	}

	// searchType == "posts"
	var posts []database.Post
	if err := s.db.Where("LOWER(content) LIKE ?", searchPattern).
		Order(`"createdAt" DESC`).
		Preload("User").
		Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("failed to search posts: %w", err)
	}

	responses := make([]dto.PostResponse, len(posts))
	for i := range posts {
		postResponse, err := s.buildPostResponse(&posts[i], nil)
		if err != nil {
			return nil, err
		}
		responses[i] = *postResponse
	}

	// Return array directly to match NestJS format
	return responses, nil
}

func (s *SearchService) getFollowingIDs(userID string) ([]string, error) {
	var follows []database.Follow
	if err := s.db.Where(`"followerId" = ?`, userID).Find(&follows).Error; err != nil {
		return nil, fmt.Errorf("failed to get follows: %w", err)
	}

	followingIDs := make([]string, len(follows))
	for i, f := range follows {
		followingIDs[i] = f.FollowingID
	}

	return followingIDs, nil
}

func (s *SearchService) buildPostResponse(post *database.Post, currentUserID *string) (*dto.PostResponse, error) {
	// Загружаем автора, если не загружен
	if post.User.ID == "" {
		if err := s.db.Model(post).Association("User").Find(&post.User); err != nil {
			// Логируем, но продолжаем
		}
	}

	// Подсчитываем лайки
	var likesCount int64
	if err := s.db.Model(&database.Like{}).Where(`"postId" = ?`, post.ID).Count(&likesCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count likes: %w", err)
	}

	// Проверяем, лайкнул ли текущий пользователь
	likedByCurrentUser := false
	if currentUserID != nil {
		var like database.Like
		if err := s.db.Where(`"postId" = ? AND "userId" = ?`, post.ID, *currentUserID).First(&like).Error; err == nil {
			likedByCurrentUser = true
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to check like: %w", err)
		}
	}

	// Загружаем комментарии
	var comments []database.Comment
	if err := s.db.Where(`"postId" = ?`, post.ID).Order(`"createdAt" ASC`).Preload("User").Find(&comments).Error; err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	commentResponses := make([]dto.CommentResponse, len(comments))
	for i, comment := range comments {
		commentResponses[i] = dto.CommentResponse{
			ID:        comment.ID,
			UserID:    comment.UserID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
			Author: dto.AuthorResponse{
				ID:       comment.User.ID,
				Username: comment.User.Username,
				Avatar:   comment.User.Avatar,
			},
		}
	}

	return &dto.PostResponse{
		ID:     post.ID,
		UserID: post.UserID,
		Author: dto.AuthorResponse{
			ID:       post.User.ID,
			Username: post.User.Username,
			Avatar:   post.User.Avatar,
		},
		Content:            post.Content,
		Likes:              int(likesCount),
		LikedByCurrentUser: likedByCurrentUser,
		CreatedAt:          post.CreatedAt,
		Comments:           commentResponses,
	}, nil
}
