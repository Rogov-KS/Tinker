package handlers

import (
	"net/http"

	"tinker-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	searchService *services.SearchService
}

func NewSearchHandler(searchService *services.SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

// Search godoc
// @Summary      Search for users or posts
// @Description  Search for users or posts by query string
// @Tags         search
// @Produce      json
// @Param        query   query     string  true  "Search query"
// @Param        type    query     string  true  "Search type (users or posts)"  Enums(users, posts)
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]string
// @Router       /search [get]
func (h *SearchHandler) Search(c *gin.Context) {
	query := c.Query("query")
	searchType := c.Query("type")

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter is required"})
		return
	}

	if searchType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type parameter is required"})
		return
	}

	if searchType != "users" && searchType != "posts" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type"})
		return
	}

	results, err := h.searchService.Search(query, searchType)
	if err != nil {
		if err.Error() == "invalid type" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, results)
}
