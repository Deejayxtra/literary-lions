package main

import (
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "literary-lions/backend/docs" // Import generated docs
	"literary-lions/backend/src/internal/db"
	"literary-lions/backend/src/internal/handlers"
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
	database, err := db.InitDB() // Renamed variable to avoid shadowing package
	if err != nil {
		log.Fatalf("Database initialization failed: %v\n", err)
	}
	defer database.Close()

	// Initialize handlers with the database connection
	handlers.InitHandlers(database)

	// Set up Gin router
	r := gin.Default()

	// Serve Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	api := r.Group("/api")
	{
		api.POST("/posts", handlers.IsAuthorized(handlers.CreatePost, "user"))
		api.PUT("/posts/:id", handlers.IsAuthorized(handlers.UpdatePost, "user"))
		api.DELETE("/posts/:id", handlers.IsAuthorized(handlers.DeletePost, "user"))
		api.GET("/posts/:id", handlers.GetPost)

		api.GET("/users", handlers.IsAuthorized(handlers.GetAllUsers, "admin"))
		api.DELETE("/users/:id", handlers.IsAuthorized(handlers.DeleteUser, "admin"))
		api.PUT("/users/:id/role", handlers.IsAuthorized(handlers.UpdateUserRole, "admin"))
	}

	// Start server on port 8080
	r.Run(":8080")
}
