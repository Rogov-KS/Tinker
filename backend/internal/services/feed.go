package services

import (
	"tinker-backend/internal/database"
	"tinker-backend/internal/dto"

	"gorm.io/gorm"
)

type FeedService struct {
	db *gorm.DB
}

func NewFeedService(db *gorm.DB) *FeedService {
	return &FeedService{db: db}
}

func (s *FeedService) GetFeed(userID string) ([]dto.PostResponse, error) {
	// Получаем список подписок
	var follows []database.Follow
	if err := s.db.Where("followerId = ?", userID).Find(&follows).Error; err != nil {
		return nil, err
	}

	if len(follows) == 0 {
		return []dto.PostResponse{}, nil
	}

	followingIDs := make([]string, len(follows))
	for i, f := range follows {
		followingIDs[i] = f.FollowingID
	}

	// Получаем посты от подписок
	var posts []database.Post
	if err := s.db.
		Preload("User").
		Preload("Comments.User").
		Preload("Likes").
		Where("userId IN ?", followingIDs).
		Order("createdAt DESC").
		Find(&posts).Error; err != nil {
		return nil, err
	}

	result := make([]dto.PostResponse, len(posts))
	for i := range posts {
		result[i] = s.mapPostToDTO(&posts[i], &userID)
	}

	return result, nil
}

func (s *FeedService) mapPostToDTO(post *database.Post, currentUserID *string) dto.PostResponse {
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
