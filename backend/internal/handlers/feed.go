package handlers

import (
	"net/http"

	"tinker-backend/internal/services"

	"github.com/gin-gonic/gin"

	_ "tinker-backend/internal/dto" // for Swagger documentation
)

type FeedHandler struct {
	feedService *services.FeedService
}

func NewFeedHandler(feedService *services.FeedService) *FeedHandler {
	return &FeedHandler{
		feedService: feedService,
	}
}

// GetFeed godoc
// @Summary      Get feed of posts
// @Description  Get feed of posts from followed users (requires authentication)
// @Tags         feed
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}  dto.PostResponse
// @Router       /feed [get]
func (h *FeedHandler) GetFeed(c *gin.Context) {
	currentUserID, _ := c.Get("userId")

	posts, err := h.feedService.GetFeed(currentUserID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, posts)
}
