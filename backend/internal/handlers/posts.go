package handlers

import (
	"net/http"

	"tinker-backend/internal/dto"
	"tinker-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PostsHandler struct {
	postsService *services.PostsService
	validator    *validator.Validate
}

func NewPostsHandler(postsService *services.PostsService) *PostsHandler {
	return &PostsHandler{
		postsService: postsService,
		validator:    validator.New(),
	}
}

// GetAllPosts godoc
// @Summary      Get all posts
// @Description  Get list of all posts with optional like status if authenticated
// @Tags         posts
// @Produce      json
// @Success      200  {array}  dto.PostResponse
// @Router       /posts [get]
func (h *PostsHandler) GetAllPosts(c *gin.Context) {
	// Опциональная аутентификация - если пользователь авторизован, передаем его ID
	var currentUserID *string
	if userID, exists := c.Get("userId"); exists {
		id := userID.(string)
		currentUserID = &id
	}

	posts, err := h.postsService.GetAllPosts(currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, posts)
}

// DeletePost godoc
// @Summary      Delete a post
// @Description  Delete a post (requires authentication, only own posts)
// @Tags         posts
// @Param        postId   path      string  true  "Post ID"
// @Security     BearerAuth
// @Success      204
// @Failure      403      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Router       /posts/{postId} [delete]
func (h *PostsHandler) DeletePost(c *gin.Context) {
	postID := c.Param("postId")
	currentUserID, _ := c.Get("userId")

	err := h.postsService.DeletePost(currentUserID.(string), postID)
	if err != nil {
		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		if err.Error() == "you can only delete your own posts" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own posts"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}

// LikePost godoc
// @Summary      Like a post
// @Description  Add a like to a post (requires authentication)
// @Tags         posts
// @Param        postId   path      string  true  "Post ID"
// @Security     BearerAuth
// @Success      200      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Router       /posts/{postId}/likes [post]
func (h *PostsHandler) LikePost(c *gin.Context) {
	postID := c.Param("postId")
	currentUserID, _ := c.Get("userId")

	err := h.postsService.LikePost(currentUserID.(string), postID)
	if err != nil {
		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post liked"})
}

// UnlikePost godoc
// @Summary      Unlike a post
// @Description  Remove a like from a post (requires authentication)
// @Tags         posts
// @Param        postId   path      string  true  "Post ID"
// @Security     BearerAuth
// @Success      204
// @Router       /posts/{postId}/likes [delete]
func (h *PostsHandler) UnlikePost(c *gin.Context) {
	postID := c.Param("postId")
	currentUserID, _ := c.Get("userId")

	err := h.postsService.UnlikePost(currentUserID.(string), postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetComments godoc
// @Summary      Get comments for a post
// @Description  Get all comments for a specific post
// @Tags         posts
// @Produce      json
// @Param        postId   path      string  true  "Post ID"
// @Success      200      {array}   dto.CommentResponse
// @Router       /posts/{postId}/comments [get]
func (h *PostsHandler) GetComments(c *gin.Context) {
	postID := c.Param("postId")
	comments, err := h.postsService.GetComments(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, comments)
}

// CreateComment godoc
// @Summary      Create a comment
// @Description  Add a comment to a post (requires authentication)
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        postId   path      string  true  "Post ID"
// @Param        request  body      dto.CreateCommentRequest  true  "Comment data"
// @Security     BearerAuth
// @Success      200      {object}  dto.CommentResponse
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Router       /posts/{postId}/comments [post]
func (h *PostsHandler) CreateComment(c *gin.Context) {
	postID := c.Param("postId")
	currentUserID, _ := c.Get("userId")

	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := h.postsService.CreateComment(currentUserID.(string), postID, req.Content)
	if err != nil {
		if err.Error() == "post not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, comment)
}
