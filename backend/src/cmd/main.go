package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"literary-lions/backend/src/docs"
	"literary-lions/backend/src/internal/handlers"
)

func main() {
	// Database connection
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	handlers.InitHandlers(db)

	// Swagger documentation
	docs.SwaggerInfo.Title = "API Documentation"
	docs.SwaggerInfo.Description = "This is the API documentation for the Gin API."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

	r := gin.Default()

	// Serve Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	api := r.Group("/api", func(c *gin.Context) {
		c.Next()
	})

	api.POST("/posts", handlers.IsAuthorized(handlers.CreatePost, "user"))
	api.PUT("/posts/:id", handlers.IsAuthorized(handlers.UpdatePost, "user"))
	api.DELETE("/posts/:id", handlers.IsAuthorized(handlers.DeletePost, "user"))
	api.GET("/posts/:id", handlers.GetPost)

	api.GET("/users", handlers.IsAuthorized(handlers.GetAllUsers, "admin"))
	api.DELETE("/users/:id", handlers.IsAuthorized(handlers.DeleteUser, "admin"))
	api.PUT("/users/:id/role", handlers.IsAuthorized(handlers.UpdateUserRole, "admin"))

	r.Run(":8080")
}
