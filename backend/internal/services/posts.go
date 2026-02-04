package services

import (
	"errors"
	"fmt"

	"log/slog"
	"tinker-backend/internal/database"
	"tinker-backend/internal/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostsService struct {
	db *gorm.DB
}

func NewPostsService(db *gorm.DB) *PostsService {
	return &PostsService{
		db: db,
	}
}

func (s *PostsService) GetAllPosts(currentUserID *string) ([]dto.PostResponse, error) {
	var posts []database.Post
	if err := s.db.Order(`"createdAt" DESC`).Preload("User").Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}

	responses := make([]dto.PostResponse, len(posts))
	for i := range posts {
		postResponse, err := s.buildPostResponse(&posts[i], currentUserID)
		if err != nil {
			return nil, err
		}
		responses[i] = *postResponse
	}

	return responses, nil
}

func (s *PostsService) DeletePost(userID, postID string) error {
	var post database.Post
	if err := s.db.Where("id = ?", postID).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("post not found")
		}
		return fmt.Errorf("database error: %w", err)
	}

	if post.UserID != userID {
		return errors.New("you can only delete your own posts")
	}

	if err := s.db.Delete(&post).Error; err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

func (s *PostsService) LikePost(userID, postID string) error {
	// Проверяем, существует ли пост
	var post database.Post
	if err := s.db.Where("id = ?", postID).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("post not found")
		}
		return fmt.Errorf("database error: %w", err)
	}

	// Проверяем, не лайкнул ли уже
	var existingLike database.Like
	if err := s.db.Where(`"postId" = ? AND "userId" = ?`, postID, userID).First(&existingLike).Error; err == nil {
		// Уже лайкнул, ничего не делаем
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("database error: %w", err)
	}

	// Создаем лайк
	like := database.Like{
		ID:     uuid.New().String(),
		UserID: userID,
		PostID: postID,
	}

	if err := s.db.Create(&like).Error; err != nil {
		return fmt.Errorf("failed to create like: %w", err)
	}

	return nil
}

func (s *PostsService) UnlikePost(userID, postID string) error {
	result := s.db.Where(`"postId" = ? AND "userId" = ?`, postID, userID).Delete(&database.Like{})
	if result.Error != nil {
		return fmt.Errorf("failed to unlike post: %w", result.Error)
	}
	return nil
}

func (s *PostsService) GetComments(postID string) ([]dto.CommentResponse, error) {
	// Проверяем, существует ли пост
	var post database.Post
	if err := s.db.Where("id = ?", postID).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	var comments []database.Comment
	if err := s.db.Where(`"postId" = ?`, postID).Order(`"createdAt" ASC`).Preload("User").Find(&comments).Error; err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	responses := make([]dto.CommentResponse, len(comments))
	for i, comment := range comments {
		responses[i] = dto.CommentResponse{
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

	return responses, nil
}

func (s *PostsService) CreateComment(userID, postID, content string) (*dto.CommentResponse, error) {
	// Проверяем, существует ли пост
	var post database.Post
	if err := s.db.Where("id = ?", postID).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	comment := database.Comment{
		ID:      uuid.New().String(),
		UserID:  userID,
		PostID:  postID,
		Content: content,
	}

	if err := s.db.Create(&comment).Error; err != nil {
		return nil, fmt.Errorf("failed to create comment: %w", err)
	}

	// Загружаем комментарий с автором для ответа
	var createdComment database.Comment
	if err := s.db.Where("id = ?", comment.ID).Preload("User").First(&createdComment).Error; err != nil {
		return nil, fmt.Errorf("failed to get created comment: %w", err)
	}

	return &dto.CommentResponse{
		ID:        createdComment.ID,
		UserID:    createdComment.UserID,
		Content:   createdComment.Content,
		CreatedAt: createdComment.CreatedAt,
		Author: dto.AuthorResponse{
			ID:       createdComment.User.ID,
			Username: createdComment.User.Username,
			Avatar:   createdComment.User.Avatar,
		},
	}, nil
}

func (s *PostsService) buildPostResponse(post *database.Post, currentUserID *string) (*dto.PostResponse, error) {
	// Загружаем автора, если не загружен
	if post.User.ID == "" {
		if err := s.db.Model(post).Association("User").Find(&post.User); err != nil {
			slog.Error("Failed to load post author", "error", err, "postId", post.ID)
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

