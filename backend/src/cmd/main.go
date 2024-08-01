package main

import (
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "literary-lions/backend/docs"
	"literary-lions/backend/src/internal/db"
	"literary-lions/backend/src/internal/handlers"

	"github.com/gin-contrib/cors"
)

// @title Gin API
// @version 1.0
// @description API documentation for the Gin application.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// Initialize the database
	database, err := db.InitDB()
	if err != nil {
		log.Fatalf("Database initialization failed: %v\n", err)
	}
	defer database.Close()

	// Initialize handlers with the database connection
	handlers.InitHandlers(database)

	// Set up Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(cors.Default())

	// Serve Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.POST("/logout", handlers.Logout)

	api := r.Group("/api")
	api.GET("/posts", handlers.GetAllPosts)

	// Authorization middleware setup
	api.Use(handlers.AuthMiddleware("user")) // Apply middleware to the group

	{
		api.GET("/users", handlers.GetAllUsers)            // Apply middleware based on role in the function
		api.POST("/post", handlers.CreatePost)             // Create a new post
		api.GET("/post/:id", handlers.GetPostByID)         // Get a specific post by ID
		api.PUT("/post/:id", handlers.UpdatePost)          // Update a specific post by ID
		api.DELETE("/post/:id", handlers.DeletePost)       // Delete a specific post by ID
		api.POST("/post/:id/comment", handlers.AddComment) // Add a comment to a specific post by ID

		// Likes and dislikes for posts
		api.POST("/post/:id/like", handlers.LikePost)             // Like a specific post by ID
		api.POST("/post/:id/dislike", handlers.DislikePost)       // Dislike a specific post by ID
		api.DELETE("/post/:id/unlike", handlers.UnlikePost)       // Unlike a specific post by ID
		api.DELETE("/post/:id/undislike", handlers.UndislikePost) // Remove dislike from a specific post by ID

		// Likes and dislikes for comments
		api.POST("/comment/:id/like", handlers.LikeComment)             // Like a specific comment by ID
		api.POST("/comment/:id/dislike", handlers.DislikeComment)       // Dislike a specific comment by ID
		api.DELETE("/comment/:id/unlike", handlers.UnlikeComment)       // Unlike a specific comment by ID
		api.DELETE("/comment/:id/undislike", handlers.UndislikeComment) // Remove dislike from a specific comment by ID
	}

	// Start server on port 8080
	r.Run(":8080")
}
