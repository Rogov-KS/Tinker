package handlers

import (
	"net/http"

	"tinker-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type CommentsHandler struct {
	commentsService *services.CommentsService
}

func NewCommentsHandler(commentsService *services.CommentsService) *CommentsHandler {
	return &CommentsHandler{
		commentsService: commentsService,
	}
}

// DeleteComment godoc
// @Summary      Delete a comment
// @Description  Delete a comment (requires authentication, only own comments)
// @Tags         comments
// @Param        commentId   path      string  true  "Comment ID"
// @Security     BearerAuth
// @Success      204
// @Failure      403         {object}  map[string]string
// @Failure      404         {object}  map[string]string
// @Router       /comments/{commentId} [delete]
func (h *CommentsHandler) DeleteComment(c *gin.Context) {
	commentID := c.Param("commentId")
	currentUserID, _ := c.Get("userId")

	err := h.commentsService.DeleteComment(currentUserID.(string), commentID)
	if err != nil {
		if err.Error() == "comment not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		if err.Error() == "you can only delete your own comments" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own comments"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}
