package services

import (
	"errors"
	"fmt"

	"log/slog"
	"tinker-backend/internal/database"
	"tinker-backend/internal/dto"
	"tinker-backend/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UsersService struct {
	db *gorm.DB
}

func NewUsersService(db *gorm.DB) *UsersService {
	return &UsersService{
		db: db,
	}
}

func (s *UsersService) GetAllUsers() ([]dto.UserResponse, error) {
	var users []database.User
	if err := s.db.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
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

	return responses, nil
}

func (s *UsersService) CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Проверяем, существует ли пользователь с таким username
	var existingUser database.User
	if err := s.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("username already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Хешируем пароль
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Создаем пользователя
	user := database.User{
		ID:       uuid.New().String(),
		Username: req.Username,
		Password: hashedPassword,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Avatar:    user.Avatar,
		Bio:       user.Bio,
		Following: []string{},
	}, nil
}

func (s *UsersService) GetUserById(userID string) (*dto.UserResponse, error) {
	var user database.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	followingIDs, err := s.getFollowingIDs(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Avatar:    user.Avatar,
		Bio:       user.Bio,
		Following: followingIDs,
	}, nil
}

func (s *UsersService) UpdateUser(userID string, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	var user database.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Проверяем уникальность username, если он изменяется
	if req.Username != nil && *req.Username != user.Username {
		var existingUser database.User
		if err := s.db.Where("username = ? AND id != ?", *req.Username, userID).First(&existingUser).Error; err == nil {
			return nil, errors.New("username already exists")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("database error: %w", err)
		}
		user.Username = *req.Username
	}

	if req.Bio != nil {
		user.Bio = req.Bio
	}

	if req.Avatar != nil {
		user.Avatar = req.Avatar
	}

	if err := s.db.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	followingIDs, err := s.getFollowingIDs(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Avatar:    user.Avatar,
		Bio:       user.Bio,
		Following: followingIDs,
	}, nil
}

func (s *UsersService) GetUserPosts(userID string) ([]dto.PostResponse, error) {
	// Проверяем, существует ли пользователь
	var user database.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	var posts []database.Post
	if err := s.db.Where(`"userId" = ?`, userID).Order(`"createdAt" DESC`).Find(&posts).Error; err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}

	responses := make([]dto.PostResponse, len(posts))
	for i, post := range posts {
		postResponse, err := s.buildPostResponse(&post, nil)
		if err != nil {
			return nil, err
		}
		responses[i] = *postResponse
	}

	return responses, nil
}

func (s *UsersService) CreatePost(userID string, req dto.CreatePostRequest) (*dto.PostResponse, error) {
	// Проверяем, существует ли пользователь
	var user database.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	post := database.Post{
		ID:      uuid.New().String(),
		UserID:  userID,
		Content: req.Content,
	}

	if err := s.db.Create(&post).Error; err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	// Загружаем пост с автором для ответа
	var createdPost database.Post
	if err := s.db.Where("id = ?", post.ID).Preload("User").First(&createdPost).Error; err != nil {
		return nil, fmt.Errorf("failed to get created post: %w", err)
	}

	return s.buildPostResponse(&createdPost, &userID)
}

func (s *UsersService) FollowUser(followerID, followingID string) error {
	if followerID == followingID {
		return errors.New("cannot follow yourself")
	}

	// Проверяем, существует ли пользователь для подписки
	var followingUser database.User
	if err := s.db.Where("id = ?", followingID).First(&followingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("database error: %w", err)
	}

	// Проверяем, не подписан ли уже
	var existingFollow database.Follow
	if err := s.db.Where(`"followerId" = ? AND "followingId" = ?`, followerID, followingID).First(&existingFollow).Error; err == nil {
		// Уже подписан, ничего не делаем
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("database error: %w", err)
	}

	// Создаем подписку
	follow := database.Follow{
		ID:          uuid.New().String(),
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	if err := s.db.Create(&follow).Error; err != nil {
		return fmt.Errorf("failed to create follow: %w", err)
	}

	return nil
}

func (s *UsersService) UnfollowUser(followerID, followingID string) error {
	result := s.db.Where(`"followerId" = ? AND "followingId" = ?`, followerID, followingID).Delete(&database.Follow{})
	if result.Error != nil {
		return fmt.Errorf("failed to unfollow user: %w", result.Error)
	}
	return nil
}

func (s *UsersService) GetFollowing(userID string) ([]dto.UserResponse, error) {
	// Проверяем, существует ли пользователь
	var user database.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	var follows []database.Follow
	if err := s.db.Where(`"followerId" = ?`, userID).Preload("Following").Find(&follows).Error; err != nil {
		return nil, fmt.Errorf("failed to get follows: %w", err)
	}

	responses := make([]dto.UserResponse, len(follows))
	for i, follow := range follows {
		followingIDs, err := s.getFollowingIDs(follow.FollowingID)
		if err != nil {
			return nil, err
		}
		responses[i] = dto.UserResponse{
			ID:        follow.Following.ID,
			Username:  follow.Following.Username,
			Avatar:    follow.Following.Avatar,
			Bio:       follow.Following.Bio,
			Following: followingIDs,
		}
	}

	return responses, nil
}

// Вспомогательные методы

func (s *UsersService) getFollowingIDs(userID string) ([]string, error) {
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

func (s *UsersService) buildPostResponse(post *database.Post, currentUserID *string) (*dto.PostResponse, error) {
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

