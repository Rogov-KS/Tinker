package services

import (
	"errors"

	"tinker-backend/internal/database"
	"tinker-backend/internal/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostsService struct {
	db *gorm.DB
}

func NewPostsService(db *gorm.DB) *PostsService {
	return &PostsService{db: db}
}

func (s *PostsService) GetAllPosts(currentUserID *string) ([]dto.PostResponse, error) {
	var posts []database.Post
	if err := s.db.
		Preload("User").
		Preload("Comments.User").
		Preload("Likes").
		Order("createdAt DESC").
		Find(&posts).Error; err != nil {
		return nil, err
	}

	result := make([]dto.PostResponse, len(posts))
	for i := range posts {
		result[i] = s.mapPostToDTO(&posts[i], currentUserID)
	}

	return result, nil
}

func (s *PostsService) DeletePost(userID, postID string) error {
	var post database.Post
	if err := s.db.Where("id = ?", postID).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("post not found")
		}
		return err
	}

	if post.UserID != userID {
		return errors.New("you can only delete your own posts")
	}

	return s.db.Delete(&post).Error
}

func (s *PostsService) LikePost(userID, postID string) error {
	// Проверяем существование поста
	var post database.Post
	if err := s.db.Where("id = ?", postID).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("post not found")
		}
		return err
	}

	// Проверяем, не лайкнул ли уже
	var existing database.Like
	if err := s.db.Where("userId = ? AND postId = ?", userID, postID).First(&existing).Error; err == nil {
		// Уже лайкнул, просто возвращаем успех (идемпотентность)
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Создаем лайк
	like := database.Like{
		ID:     uuid.New().String(),
		UserID: userID,
		PostID: postID,
	}

	return s.db.Create(&like).Error
}

func (s *PostsService) UnlikePost(userID, postID string) error {
	return s.db.Where("userId = ? AND postId = ?", userID, postID).Delete(&database.Like{}).Error
}

func (s *PostsService) GetComments(postID string) ([]dto.CommentResponse, error) {
	var comments []database.Comment
	if err := s.db.
		Preload("User").
		Where("postId = ?", postID).
		Order("createdAt DESC").
		Find(&comments).Error; err != nil {
		return nil, err
	}

	result := make([]dto.CommentResponse, len(comments))
	for i, c := range comments {
		result[i] = dto.CommentResponse{
			ID:        c.ID,
			UserID:    c.UserID,
			Content:   c.Content,
			CreatedAt: c.CreatedAt,
			Author: dto.AuthorResponse{
				ID:       c.User.ID,
				Username: c.User.Username,
				Avatar:   c.User.Avatar,
			},
		}
	}

	return result, nil
}

func (s *PostsService) CreateComment(userID, postID string, content string) (*dto.CommentResponse, error) {
	// Проверяем существование поста
	var post database.Post
	if err := s.db.Where("id = ?", postID).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	comment := database.Comment{
		ID:      uuid.New().String(),
		UserID:  userID,
		PostID:  postID,
		Content: content,
	}

	if err := s.db.Create(&comment).Error; err != nil {
		return nil, err
	}

	// Загружаем с пользователем
	var createdComment database.Comment
	if err := s.db.Preload("User").Where("id = ?", comment.ID).First(&createdComment).Error; err != nil {
		return nil, err
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

func (s *PostsService) mapPostToDTO(post *database.Post, currentUserID *string) dto.PostResponse {
	comments := make([]dto.CommentResponse, len(post.Comments))
	for i, c := range post.Comments {
		comments[i] = dto.CommentResponse{
			ID:        c.ID,
			UserID:    c.UserID,
			Content:   c.Content,
			CreatedAt: c.CreatedAt,
			Author: dto.AuthorResponse{
				ID:       c.User.ID,
				Username: c.User.Username,
				Avatar:   c.User.Avatar,
			},
		}
	}

	likedByCurrentUser := false
	if currentUserID != nil {
		for _, like := range post.Likes {
			if like.UserID == *currentUserID {
				likedByCurrentUser = true
				break
			}
		}
	}

	return dto.PostResponse{
		ID:     post.ID,
		UserID: post.UserID,
		Author: dto.AuthorResponse{
			ID:       post.User.ID,
			Username: post.User.Username,
			Avatar:   post.User.Avatar,
		},
		Content:            post.Content,
		Likes:              len(post.Likes),
		LikedByCurrentUser: likedByCurrentUser,
		CreatedAt:          post.CreatedAt,
		Comments:           comments,
	}
}
