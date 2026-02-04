package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tinker-backend/internal/config"
	"tinker-backend/internal/database"
	"tinker-backend/internal/handlers"
	"tinker-backend/internal/middleware"
	"tinker-backend/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "tinker-backend/docs" // swagger docs
)

// @title           Tinker API
// @version         1.0
// @description     API for a simple social media application
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3001
// @BasePath  /api
// @schemes   http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// Настраиваем логирование
	setupLogging(cfg.LogLevel)

	// Подключаемся к БД
	if err := database.Connect(cfg.DatabaseURL); err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer database.Close()

	// Инициализируем сервисы
	authService := services.NewAuthService(database.DB, cfg.JWTSecret)
	usersService := services.NewUsersService(database.DB)
	postsService := services.NewPostsService(database.DB)
	commentsService := services.NewCommentsService(database.DB)
	feedService := services.NewFeedService(database.DB)
	searchService := services.NewSearchService(database.DB)

	// Инициализируем handlers
	authHandler := handlers.NewAuthHandler(authService)
	usersHandler := handlers.NewUsersHandler(usersService)
	postsHandler := handlers.NewPostsHandler(postsService)
	commentsHandler := handlers.NewCommentsHandler(commentsService)
	feedHandler := handlers.NewFeedHandler(feedService)
	searchHandler := handlers.NewSearchHandler(searchService)

	// Настраиваем Gin
	if cfg.LogLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Глобальные middleware
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Recovery())

	// API роуты
	api := router.Group("/api")
	{
		// Auth
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
		}

		// Users
		users := api.Group("/users")
		{
			users.GET("", usersHandler.GetAllUsers)
			users.POST("", usersHandler.CreateUser)
			users.GET("/:userId", usersHandler.GetUserById)
			users.PATCH("/:userId", middleware.AuthMiddleware(cfg.JWTSecret), usersHandler.UpdateUser)
			users.GET("/:userId/posts", usersHandler.GetUserPosts)
			users.POST("/:userId/posts", middleware.AuthMiddleware(cfg.JWTSecret), usersHandler.CreatePost)
			users.POST("/:userId/follow", middleware.AuthMiddleware(cfg.JWTSecret), usersHandler.FollowUser)
			users.DELETE("/:userId/follow", middleware.AuthMiddleware(cfg.JWTSecret), usersHandler.UnfollowUser)
			users.GET("/:userId/following", usersHandler.GetFollowing)
		}

		// Posts
		posts := api.Group("/posts")
		{
			posts.GET("", postsHandler.GetAllPosts)
			posts.DELETE("/:postId", middleware.AuthMiddleware(cfg.JWTSecret), postsHandler.DeletePost)
			posts.POST("/:postId/likes", middleware.AuthMiddleware(cfg.JWTSecret), postsHandler.LikePost)
			posts.DELETE("/:postId/likes", middleware.AuthMiddleware(cfg.JWTSecret), postsHandler.UnlikePost)
			posts.GET("/:postId/comments", postsHandler.GetComments)
			posts.POST("/:postId/comments", middleware.AuthMiddleware(cfg.JWTSecret), postsHandler.CreateComment)
		}

		// Comments
		comments := api.Group("/comments")
		{
			comments.DELETE("/:commentId", middleware.AuthMiddleware(cfg.JWTSecret), commentsHandler.DeleteComment)
		}

		// Feed
		feed := api.Group("/feed")
		{
			feed.GET("", middleware.AuthMiddleware(cfg.JWTSecret), feedHandler.GetFeed)
		}

		// Search
		search := api.Group("/search")
		{
			search.GET("", searchHandler.Search)
		}

		// Swagger documentation
		api.GET("/docs", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/api/docs/index.html")
		})
		api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Запускаем сервер
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		slog.Info("Starting server", "port", cfg.Port, "address", fmt.Sprintf("http://localhost:%s/api", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// Ожидаем сигнал для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exited")
}

func setupLogging(level string) {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
