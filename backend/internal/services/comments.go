package services

import (
	"errors"
	"fmt"

	"tinker-backend/internal/database"

	"gorm.io/gorm"
)

type CommentsService struct {
	db *gorm.DB
}

func NewCommentsService(db *gorm.DB) *CommentsService {
	return &CommentsService{
		db: db,
	}
}

func (s *CommentsService) DeleteComment(userID, commentID string) error {
	var comment database.Comment
	if err := s.db.Where("id = ?", commentID).First(&comment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("comment not found")
		}
		return fmt.Errorf("database error: %w", err)
	}

	if comment.UserID != userID {
		return errors.New("you can only delete your own comments")
	}

	if err := s.db.Delete(&comment).Error; err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	return nil
}

