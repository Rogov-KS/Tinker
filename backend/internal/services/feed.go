package services

import (
	"errors"
	"fmt"

	"tinker-backend/internal/database"
	"tinker-backend/internal/dto"

	"gorm.io/gorm"
)

type FeedService struct {
	db *gorm.DB
}

func NewFeedService(db *gorm.DB) *FeedService {
	return &FeedService{
		db: db,
	}
}

func (s *FeedService) GetFeed(userID string) ([]dto.PostResponse, error) {
	// Получаем список пользователей, на которых подписан текущий пользователь
	var follows []database.Follow
	if err := s.db.Where(`"followerId" = ?`, userID).Find(&follows).Error; err != nil {
		return nil, fmt.Errorf("failed to get follows: %w", err)
	}

	// Если пользователь ни на кого не подписан, возвращаем пустой список
	if len(follows) == 0 {
		return []dto.PostResponse{}, nil
	}

	// Извлекаем ID пользователей, на которых подписан
	followingIDs := make([]string, len(follows))
	for i, f := range follows {
		followingIDs[i] = f.FollowingID
	}

	// Получаем посты от пользователей, на которых подписан
	var posts []database.Post
	if err := s.db.Where(`"userId" IN ?`, followingIDs).
		Order(`"createdAt" DESC`).
		Preload("User").
		Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("failed to get feed posts: %w", err)
	}

	// Строим ответы
	responses := make([]dto.PostResponse, len(posts))
	currentUserID := &userID
	for i := range posts {
		postResponse, err := s.buildPostResponse(&posts[i], currentUserID)
		if err != nil {
			return nil, err
		}
		responses[i] = *postResponse
	}

	return responses, nil
}

func (s *FeedService) buildPostResponse(post *database.Post, currentUserID *string) (*dto.PostResponse, error) {
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
		ID:                 post.ID,
		UserID:             post.UserID,
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

