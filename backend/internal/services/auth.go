package services

import (
	"errors"
	"fmt"
	"time"

	"log/slog"
	"tinker-backend/internal/database"
	"tinker-backend/internal/dto"
	"tinker-backend/internal/utils"

	"github.com/golang-jwt/jwt/v5"

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
	slog.Info("Starting ValidateUserCredentials", "username", username)
	var user database.User
	slog.Info("Querying database for user", "username", username)
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		slog.Error("Database query failed", "error", err, "errorType", fmt.Sprintf("%T", err), "username", username)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Warn("Login failed: user not found", "username", username)
			return nil, errors.New("invalid credentials")
		}
		slog.Error("Database error in ValidateUserCredentials", "error", err, "username", username)
		return nil, fmt.Errorf("database error: %w", err)
	}
	slog.Info("User found in database", "userId", user.ID, "username", user.Username)

	slog.Info("Checking password hash", "username", username, "passwordLen", len(password), "hashLen", len(user.Password))
	isValid := utils.CheckPasswordHash(password, user.Password)
	slog.Info("Password check result", "username", username, "isValid", isValid)
	if !isValid {
		slog.Warn("Login failed: invalid password", "username", username)
		return nil, errors.New("invalid credentials")
	}

	slog.Info("Password validated successfully", "username", username)
	return &user, nil
}

func (s *AuthService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	slog.Info("Login attempt", "username", req.Username)

	user, err := s.ValidateUserCredentials(req.Username, req.Password)
	if err != nil {
		slog.Error("ValidateUserCredentials failed", "error", err, "username", req.Username)
		return nil, err
	}

	slog.Info("User validated", "userId", user.ID, "username", user.Username)

	// Получаем список подписок
	var follows []database.Follow
	slog.Info("Querying follows", "userId", user.ID)
	if err := s.db.Table("Follow").Where(`"followerId" = ?`, user.ID).Find(&follows).Error; err != nil {
		slog.Error("Failed to get follows", "error", err, "userId", user.ID, "errorType", fmt.Sprintf("%T", err), "errorString", err.Error())
		return nil, fmt.Errorf("failed to get follows: %w", err)
	}

	slog.Info("Got follows", "count", len(follows), "userId", user.ID)

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
