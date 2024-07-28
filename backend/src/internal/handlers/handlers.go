package handlers

import (
	"database/sql"
	// "fmt"
	"net/http"
	// "strings"

	// "time"

	// "literary-lions/backend/src/internal/auth"
	"literary-lions/backend/src/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var db *sql.DB

type Credentials struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegistrationRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func InitHandlers(database *sql.DB) {
	db = database
	models.SetDatabase(db)
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegistrationRequest true "User registration request"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Router /register [post]
func Register(c *gin.Context) {
	var req RegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// if err := auth.RegisterUser(db, req.Email, req.Username, req.Password); err != nil {
	if err := models.RegisterUser(req.Email, req.Username, req.Password); err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed"})
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

// UpdatePost godoc
// @Summary Update a post
// @Description Update a post
// @Tags post
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param post body models.Post true "Updated post object"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/post/{id} [put]
// @Security ApiKeyAuth
func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	query := "UPDATE posts SET category_id = $1, title = $2, content = $3 WHERE id = $4"
	result, err := db.Exec(query, post.Category, post.Title, post.Content, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update post"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
}

// DeletePost godoc
// @Summary Delete a post
// @Description Delete a post
// @Tags post
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/post/{id} [delete]
// @Security ApiKeyAuth
func DeletePost(c *gin.Context) {
	id := c.Param("id")

	query := "DELETE FROM posts WHERE id = $1"
	result, err := db.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete post"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Get all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Failure 401 {object} gin.H
// @Router /api/users [get]
// @Security ApiKeyAuth
func GetAllUsers(c *gin.Context) {
	rows, err := db.Query("SELECT id, username, email, role FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve users"})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not scan user data"})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Get a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/users/{id} [get]
// @Security ApiKeyAuth
func GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	query := "SELECT id, username, email, role FROM users WHERE id = $1"
	err := db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve user"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update a user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.User true "Updated user object"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/users/{id} [put]
// @Security ApiKeyAuth
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	query := "UPDATE users SET username = $1, email = $2, role = $3 WHERE id = $4"
	result, err := db.Exec(query, user.Username, user.Email, user.Role, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser godoc

// @Summary Delete a user
// @Description Delete a user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/users/{id} [delete]
// @Security ApiKeyAuth
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	query := "DELETE FROM users WHERE id = $1"
	result, err := db.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// UpdateUserRole godoc

// @Summary Update a user's role
// @Description Update a user's role by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param role body string true "New role"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/users/{id}/role [put]
// @Security ApiKeyAuth
func UpdateUserRole(c *gin.Context) {
	id := c.Param("id")
	var requestBody struct {
		Role string `json:"role"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if requestBody.Role != "admin" && requestBody.Role != "user" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	res, err := db.Exec("UPDATE users SET role = ? WHERE id = ?", requestBody.Role, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user role"})
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve affected rows"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User role updated successfully"})
}

// InitializeRoutes initializes all the routes
// func InitializeRoutes(router *gin.Engine) {
// 	// User authentication routes
// 	router.POST("/register", Register)
// 	router.POST("/login", Login)

// 	// Protected routes
// 	api := router.Group("/api")
// 	{
// 		api.GET("/posts/:id", GetPost)
// 		api.POST("/posts", AuthMiddleware("user", CreatePost))
// 		api.PUT("/posts/:id", AuthMiddleware("user", UpdatePost))
// 		api.DELETE("/posts/:id", AuthMiddleware("user", DeletePost))
// 		api.GET("/users", AuthMiddleware("admin", GetAllUsers))
// 		api.DELETE("/users/:id", AuthMiddleware("admin", DeleteUser))
// 		api.PUT("/users/:id/role", AuthMiddleware("admin", UpdateUserRole))
// 	}
// }

// func InitializeRoutes(router *gin.Engine) {
// 	// User authentication routes
// 	router.POST("/register", Register)
// 	router.POST("/login", Login)

// 	// Protected routes
// 	api := router.Group("/api")
// 	{
// 		// Apply AuthMiddleware directly as middleware
// 		api.POST("/posts", AuthMiddleware("user"), CreatePost)
// 		// api.POST("/posts", CreatePost)
// 		api.PUT("/posts/:id", AuthMiddleware("user"), UpdatePost)
// 		api.DELETE("/posts/:id", AuthMiddleware("user"), DeletePost)
// 		api.GET("/users", AuthMiddleware("admin"), GetAllUsers)
// 		api.DELETE("/users/:id", AuthMiddleware("admin"), DeleteUser)
// 		api.PUT("/users/:id/role", AuthMiddleware("admin"), UpdateUserRole)
// 	}
// }
