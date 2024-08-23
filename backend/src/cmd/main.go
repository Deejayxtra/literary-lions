package main

import (
	// "fmt"
	"fmt"
	_ "literary-lions/backend/docs"
	"literary-lions/backend/src/internal/db"
	"literary-lions/backend/src/internal/handlers"
	"literary-lions/backend/src/internal/middleware"
	"literary-lions/backend/src/internal/models"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Gin API
// @version 1.0
// @description API documentation for the Gin application.
// @host localhost:8080/api/v1.0
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

	// Set the database for the models package
	models.SetDatabase(database)

	// Initialize handlers with the database connection
	handlers.InitHandlers(database)

	// Set up Gin router
	r := gin.Default()

	// Apply the NoCache middleware globally
	r.Use(middleware.NoCache())

	// Serve Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Slice to store routes
	var routes []string

	note := "<h4>Only the GET endpoints will display, the POST endpoints will give Page not found as there is no payload to send</h4>"
	// Root route with clickable links
	r.GET("/", func(c *gin.Context) {
		routeList := note + "<h1>Available Endpoints:</h1><ul>"
		for _, route := range routes {
			routeList += fmt.Sprintf(`<li><a href="%s">%s</a></li>`, route, route)
		}
		routeList += "</ul>"
		c.Header("Content-Type", "text/html")
		c.String(200, routeList)
	})

	// Function to add routes to the list and define them
	addRoute := func(method, path string, handler gin.HandlerFunc) {
		routes = append(routes, path)
		r.Handle(method, path, handler)
	}

	// Define some example routes
	addRoute("POST", "/login", handlers.Login)
	addRoute("POST", "/logout", handlers.Logout)
	addRoute("POST", "/register", handlers.Register)
	addRoute("GET", "/posts", handlers.GetAllPosts)

	api := r.Group("/api/v1.0")

	// Public routes
	api.POST("/register", handlers.Register)
	api.POST("/login", handlers.Login)
	api.POST("/logout", handlers.Logout)
	api.GET("/posts", handlers.GetAllPosts)

	api.GET("/post/:id", handlers.GetPostByID) // Get a specific post by ID

	// Authorization middleware setup
	api.Use(handlers.AuthMiddleware("user")) // Apply middleware to the group

	{
		addRoute("GET", "/filtered-posts", handlers.GetAllPosts)
		addRoute("GET", "/users", handlers.GetAllUsers)
		addRoute("POST", "/post", handlers.CreatePost)
		addRoute("PUT", "/post/:id", handlers.UpdatePost)
		addRoute("DELETE", "/post/:id", handlers.DeletePost)
		addRoute("POST", "/post/:id/comment", handlers.AddComment)
		addRoute("PUT", "/userprofile-update", handlers.UpdateUserProfile)
		addRoute("POST", "/post/:id/like", handlers.LikePost)
		addRoute("POST", "/post/:id/dislike", handlers.DislikePost)
		addRoute("POST", "/comment/:id/like", handlers.LikeComment)
		addRoute("POST", "/comment/:id/dislike", handlers.DislikeComment)
	}

	{
		api.GET("/filtered-posts", handlers.GetAllPosts)           // This is the endpoint to be called when filter query is set
		api.GET("/users", handlers.GetAllUsers)                    // Apply middleware based on role in the function
		api.POST("/post", handlers.CreatePost)                     // Create a new post
		api.PUT("/post/:id", handlers.UpdatePost)                  // Update a specific post by ID
		api.DELETE("/post/:id", handlers.DeletePost)               // Delete a specific post by ID
		api.POST("/post/:id/comment", handlers.AddComment)         // Add a comment to a specific post by ID
		api.PUT("/userprofile-update", handlers.UpdateUserProfile) // Update user profile

		// Likes and dislikes for posts
		api.POST("/post/:id/like", handlers.LikePost)       // Like a specific post by ID
		api.POST("/post/:id/dislike", handlers.DislikePost) // Dislike a specific post by ID

		// // Likes and dislikes for comments
		api.POST("/comment/:id/like", handlers.LikeComment)       // Like a specific comment by ID
		api.POST("/comment/:id/dislike", handlers.DislikeComment) // Dislike a specific comment by ID

	}

	// Start server on port 8080
	r.Run(":8080")
}
