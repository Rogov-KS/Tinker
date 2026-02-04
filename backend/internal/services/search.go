package services

import (
	"errors"
	"strings"

	"gorm.io/gorm"
	"tinker-backend/internal/database"
	"tinker-backend/internal/dto"
)

type SearchService struct {
	db *gorm.DB
}

func NewSearchService(db *gorm.DB) *SearchService {
	return &SearchService{db: db}
}

func (s *SearchService) Search(query, searchType string) (interface{}, error) {
	if searchType == "users" {
		return s.searchUsers(query)
	} else if searchType == "posts" {
		return s.searchPosts(query)
	}
	return nil, errors.New("invalid type")
}

func (s *SearchService) searchUsers(query string) ([]dto.UserResponse, error) {
	var users []database.User
	queryLower := strings.ToLower(query)
	
	if err := s.db.
		Preload("Following").
		Where("LOWER(username) LIKE ?", "%"+queryLower+"%").
		Find(&users).Error; err != nil {
		return nil, err
	}

	result := make([]dto.UserResponse, len(users))
	for i, user := range users {
		followingIDs := make([]string, 0)
		if user.Following != nil {
			for _, f := range user.Following {
				followingIDs = append(followingIDs, f.FollowingID)
			}
		}
		result[i] = dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Avatar:    user.Avatar,
			Bio:       user.Bio,
			Following: followingIDs,
		}
	}

	return result, nil
}

func (s *SearchService) searchPosts(query string) ([]dto.PostResponse, error) {
	var posts []database.Post
	queryLower := strings.ToLower(query)
	
	if err := s.db.
		Preload("User").
		Preload("Comments.User").
		Preload("Likes").
		Where("LOWER(content) LIKE ?", "%"+queryLower+"%").
		Order("createdAt DESC").
		Find(&posts).Error; err != nil {
		return nil, err
	}

	result := make([]dto.PostResponse, len(posts))
	for i := range posts {
		// Для поиска постов не передаем currentUserID, поэтому likedByCurrentUser всегда false
		result[i] = s.mapPostToDTO(&posts[i], nil)
	}

	return result, nil
}

func (s *SearchService) mapPostToDTO(post *database.Post, currentUserID *string) dto.PostResponse {
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
		ID:               post.ID,
		UserID:           post.UserID,
		Author: dto.AuthorResponse{
			ID:       post.User.ID,
			Username: post.User.Username,
			Avatar:   post.User.Avatar,
		},
		Content:          post.Content,
		Likes:            len(post.Likes),
		LikedByCurrentUser: likedByCurrentUser,
		CreatedAt:        post.CreatedAt,
		Comments:         comments,
	}
}

