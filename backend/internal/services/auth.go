package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"tinker-backend/internal/database"
	"tinker-backend/internal/dto"
	"tinker-backend/internal/utils"

	"gorm.io/gorm"
)

type AuthService struct {
	db        *gorm.DB
	jwtSecret string
}

func NewAuthService(db *gorm.DB, jwtSecret string) *AuthService {
	return &AuthService{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) ValidateUserCredentials(username, password string) (*database.User, error) {
	var user database.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Warn("Login failed: user not found", "username", username)
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		slog.Warn("Login failed: invalid password", "username", username)
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}

func (s *AuthService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	slog.Info("Login attempt", "username", req.Username)

	user, err := s.ValidateUserCredentials(req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	// Получаем список подписок
	var follows []database.Follow
	if err := s.db.Where("follower_id = ?", user.ID).Find(&follows).Error; err != nil {
		return nil, err
	}

	followingIDs := make([]string, len(follows))
	for i, f := range follows {
		followingIDs[i] = f.FollowingID
	}

	// Генерируем JWT токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"exp":    time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	slog.Info("Successful login", "userId", user.ID, "username", user.Username)

	return &dto.AuthResponse{
		Token: tokenString,
		User: dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Avatar:    user.Avatar,
			Bio:       user.Bio,
			Following: followingIDs,
		},
	}, nil
}

