package services

import (
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"tinker-backend/internal/database"
	"tinker-backend/internal/dto"
	"tinker-backend/internal/utils"
)

type UsersService struct {
	db *gorm.DB
}

func NewUsersService(db *gorm.DB) *UsersService {
	return &UsersService{db: db}
}

func (s *UsersService) GetAllUsers() ([]dto.UserResponse, error) {
	var users []database.User
	if err := s.db.Preload("Following").Order("created_at ASC").Find(&users).Error; err != nil {
		return nil, err
	}

	result := make([]dto.UserResponse, len(users))
	for i, user := range users {
		result[i] = s.mapUserToDTO(&user)
	}

	return result, nil
}

func (s *UsersService) GetUserById(userID string) (*dto.UserResponse, error) {
	var user database.User
	if err := s.db.Preload("Following").Where("id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	result := s.mapUserToDTO(&user)
	return &result, nil
}

func (s *UsersService) CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	slog.Info("Creating new user", "username", req.Username)

	// Проверяем существование пользователя
	var existing database.User
	if err := s.db.Where("username = ?", req.Username).First(&existing).Error; err == nil {
		slog.Warn("Username already exists", "username", req.Username)
		return nil, errors.New("username already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Хешируем пароль
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Создаем пользователя
	user := database.User{
		ID:       uuid.New().String(),
		Username: req.Username,
		Password: hashedPassword,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	slog.Info("User created successfully", "userId", user.ID, "username", user.Username)

	result := s.mapUserToDTO(&user)
	return &result, nil
}

func (s *UsersService) UpdateUser(userID string, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	updates := make(map[string]interface{})

	if req.Username != nil {
		// Проверяем уникальность username
		var existing database.User
		if err := s.db.Where("username = ? AND id != ?", *req.Username, userID).First(&existing).Error; err == nil {
			return nil, errors.New("username already exists")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		updates["username"] = *req.Username
	}

	if req.Bio != nil {
		updates["bio"] = *req.Bio
	}

	if req.Avatar != nil {
		updates["avatar"] = *req.Avatar
	}

	if len(updates) == 0 {
		return s.GetUserById(userID)
	}

	if err := s.db.Model(&database.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		return nil, err
	}

	return s.GetUserById(userID)
}

func (s *UsersService) FollowUser(followerID, followingID string) error {
	if followerID == followingID {
		return errors.New("cannot follow yourself")
	}

	// Проверяем существование пользователя для подписки
	var following database.User
	if err := s.db.Where("id = ?", followingID).First(&following).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Проверяем, не подписан ли уже
	var existing database.Follow
	if err := s.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).First(&existing).Error; err == nil {
		// Уже подписан, просто возвращаем успех (идемпотентность)
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Создаем подписку
	follow := database.Follow{
		ID:          uuid.New().String(),
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	return s.db.Create(&follow).Error
}

func (s *UsersService) UnfollowUser(followerID, followingID string) error {
	return s.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&database.Follow{}).Error
}

func (s *UsersService) GetFollowing(userID string) ([]dto.UserResponse, error) {
	var follows []database.Follow
	if err := s.db.Preload("Following").Preload("Following.Following").Where("follower_id = ?", userID).Find(&follows).Error; err != nil {
		return nil, err
	}

	result := make([]dto.UserResponse, len(follows))
	for i, f := range follows {
		result[i] = s.mapUserToDTO(&f.Following)
	}

	return result, nil
}

func (s *UsersService) GetUserPosts(userID string) ([]dto.PostResponse, error) {
	var posts []database.Post
	if err := s.db.
		Preload("User").
		Preload("Comments.User").
		Preload("Likes").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&posts).Error; err != nil {
		return nil, err
	}

	result := make([]dto.PostResponse, len(posts))
	for i := range posts {
		result[i] = s.mapPostToDTO(&posts[i], &userID)
	}

	return result, nil
}

func (s *UsersService) CreatePost(userID string, req dto.CreatePostRequest) (*dto.PostResponse, error) {
	slog.Info("Creating post for user", "userId", userID)

	post := database.Post{
		ID:      uuid.New().String(),
		UserID:  userID,
		Content: req.Content,
	}

	if err := s.db.Create(&post).Error; err != nil {
		return nil, err
	}

	// Загружаем с отношениями
	var createdPost database.Post
	if err := s.db.
		Preload("User").
		Preload("Comments.User").
		Preload("Likes").
		Where("id = ?", post.ID).
		First(&createdPost).Error; err != nil {
		return nil, err
	}

	slog.Info("Post created", "postId", post.ID, "userId", userID)

	result := s.mapPostToDTO(&createdPost, &userID)
	return &result, nil
}

func (s *UsersService) mapUserToDTO(user *database.User) dto.UserResponse {
	followingIDs := make([]string, 0)
	if user.Following != nil {
		for _, f := range user.Following {
			followingIDs = append(followingIDs, f.FollowingID)
		}
	}

	return dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Avatar:    user.Avatar,
		Bio:       user.Bio,
		Following: followingIDs,
	}
}

func (s *UsersService) mapPostToDTO(post *database.Post, currentUserID *string) dto.PostResponse {
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

