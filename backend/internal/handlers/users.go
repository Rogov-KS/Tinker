package handlers

import (
	"net/http"

	"tinker-backend/internal/dto"
	"tinker-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UsersHandler struct {
	usersService *services.UsersService
	validator    *validator.Validate
}

func NewUsersHandler(usersService *services.UsersService) *UsersHandler {
	return &UsersHandler{
		usersService: usersService,
		validator:    validator.New(),
	}
}

// GetAllUsers godoc
// @Summary      Get all users
// @Description  Get list of all users
// @Tags         users
// @Produce      json
// @Success      200  {array}  dto.UserResponse
// @Router       /users [get]
func (h *UsersHandler) GetAllUsers(c *gin.Context) {
	users, err := h.usersService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// CreateUser godoc
// @Summary      Create a new user
// @Description  Register a new user account
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateUserRequest  true  "User data"
// @Success      200      {object}  dto.UserResponse
// @Failure      400      {object}  map[string]string
// @Router       /users [post]
func (h *UsersHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.usersService.CreateUser(req)
	if err != nil {
		if err.Error() == "username already exists" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUserById godoc
// @Summary      Get user by ID
// @Description  Get user information by user ID
// @Tags         users
// @Produce      json
// @Param        userId   path      string  true  "User ID"
// @Success      200      {object}  dto.UserResponse
// @Failure      404      {object}  map[string]string
// @Router       /users/{userId} [get]
func (h *UsersHandler) GetUserById(c *gin.Context) {
	userID := c.Param("userId")
	user, err := h.usersService.GetUserById(userID)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary      Update user profile
// @Description  Update user profile information (requires authentication)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        userId   path      string  true  "User ID"
// @Param        request  body      dto.UpdateUserRequest  true  "Update data"
// @Security     BearerAuth
// @Success      200      {object}  dto.UserResponse
// @Failure      400      {object}  map[string]string
// @Failure      403      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Router       /users/{userId} [patch]
func (h *UsersHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("userId")
	currentUserID, _ := c.Get("userId")

	if currentUserID.(string) != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own profile"})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.usersService.UpdateUser(userID, req)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		if err.Error() == "username already exists" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUserPosts godoc
// @Summary      Get user posts
// @Description  Get all posts by a specific user
// @Tags         users
// @Produce      json
// @Param        userId   path      string  true  "User ID"
// @Success      200      {array}   dto.PostResponse
// @Router       /users/{userId}/posts [get]
func (h *UsersHandler) GetUserPosts(c *gin.Context) {
	userID := c.Param("userId")
	posts, err := h.usersService.GetUserPosts(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, posts)
}

// CreatePost godoc
// @Summary      Create a new post
// @Description  Create a new post as a user (requires authentication)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        userId   path      string  true  "User ID"
// @Param        request  body      dto.CreatePostRequest  true  "Post data"
// @Security     BearerAuth
// @Success      200      {object}  dto.PostResponse
// @Failure      400      {object}  map[string]string
// @Failure      403      {object}  map[string]string
// @Router       /users/{userId}/posts [post]
func (h *UsersHandler) CreatePost(c *gin.Context) {
	userID := c.Param("userId")
	currentUserID, _ := c.Get("userId")

	if currentUserID.(string) != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only create posts as yourself"})
		return
	}

	var req dto.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := h.usersService.CreatePost(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// FollowUser godoc
// @Summary      Follow a user
// @Description  Follow another user (requires authentication)
// @Tags         users
// @Param        userId   path      string  true  "User ID to follow"
// @Security     BearerAuth
// @Success      200      {object}  map[string]string
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Router       /users/{userId}/follow [post]
func (h *UsersHandler) FollowUser(c *gin.Context) {
	userID := c.Param("userId")
	currentUserID, _ := c.Get("userId")

	err := h.usersService.FollowUser(currentUserID.(string), userID)
	if err != nil {
		if err.Error() == "cannot follow yourself" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot follow yourself"})
			return
		}
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User followed"})
}

// UnfollowUser godoc
// @Summary      Unfollow a user
// @Description  Unfollow a user (requires authentication)
// @Tags         users
// @Param        userId   path      string  true  "User ID to unfollow"
// @Security     BearerAuth
// @Success      204
// @Router       /users/{userId}/follow [delete]
func (h *UsersHandler) UnfollowUser(c *gin.Context) {
	userID := c.Param("userId")
	currentUserID, _ := c.Get("userId")

	err := h.usersService.UnfollowUser(currentUserID.(string), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetFollowing godoc
// @Summary      Get list of users being followed
// @Description  Get list of users that a specific user is following
// @Tags         users
// @Produce      json
// @Param        userId   path      string  true  "User ID"
// @Success      200      {array}   dto.UserResponse
// @Router       /users/{userId}/following [get]
func (h *UsersHandler) GetFollowing(c *gin.Context) {
	userID := c.Param("userId")
	following, err := h.usersService.GetFollowing(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, following)
}
