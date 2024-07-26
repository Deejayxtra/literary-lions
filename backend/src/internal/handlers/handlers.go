// package handlers

// import (
// 	"database/sql"
// 	"net/http"
// 	"strings"
// 	"time"

// 	"literary-lions/backend/src/internal/auth"
// 	"literary-lions/backend/src/internal/models"

// 	"github.com/dgrijalva/jwt-go"
// 	"github.com/gin-gonic/gin"
// )

// var db *sql.DB
// var jwtKey = []byte("my_secret_key")

// type Credentials struct {
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

// type Claims struct {
// 	Username string `json:"username"`
// 	Role     string `json:"role"`
// 	jwt.StandardClaims
// }

// func InitHandlers(database *sql.DB) {
// 	db = database
// }

// // IsAuthorized checks the token validity and role
// func IsAuthorized(handlerFunc gin.HandlerFunc, requiredRole string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.Request.Header.Get("Authorization")
// 		if authHeader == "" {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
// 			c.Abort()
// 			return
// 		}

// 		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
// 		claims := &Claims{}

// 		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
// 			return jwtKey, nil
// 		})

// 		if err != nil || !token.Valid {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 			c.Abort()
// 			return
// 		}

// 		if claims.Role != requiredRole && claims.Role != "admin" {
// 			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
// 			c.Abort()
// 			return
// 		}

// 		c.Set("username", claims.Username)
// 		c.Set("role", claims.Role)

// 		handlerFunc(c)
// 	}
// }

// // AuthMiddleware is a middleware to check authorization token and role
// func AuthMiddleware(requiredRole string, handlerFunc gin.HandlerFunc) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if !strings.HasPrefix(authHeader, "Bearer ") {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
// 			c.Abort()
// 			return
// 		}

// 		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

// 		claims := &Claims{}
// 		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
// 			return jwtKey, nil
// 		})

// 		if err != nil || !token.Valid {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 			c.Abort()
// 			return
// 		}

// 		if claims.Role != requiredRole && claims.Role != "admin" {
// 			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
// 			c.Abort()
// 			return
// 		}

// 		c.Set("username", claims.Username)
// 		c.Set("role", claims.Role)

// 		handlerFunc(c)
// 	}
// }

// // Register godoc
// // @Summary Register a new user
// // @Description Register a new user
// // @Tags auth
// // @Accept json
// // @Produce json
// // @Param user body Credentials true "User credentials"
// // @Success 200 {object} gin.H
// // @Failure 400 {object} gin.H
// // @Router /register [post]
// func Register(c *gin.Context) {
// 	var creds Credentials
// 	if err := c.ShouldBindJSON(&creds); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
// 		return
// 	}

// 	if err := auth.RegisterUser(db, creds.Email, creds.Username, creds.Password); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
// }

// // Login godoc
// // @Summary Login a user
// // @Description Login a user
// // @Tags auth
// // @Accept json
// // @Produce json
// // @Param credentials body Credentials true "User credentials"
// // @Success 200 {object} gin.H
// // @Failure 400 {object} gin.H
// // @Failure 401 {object} gin.H
// // @Router /login [post]
// func Login(c *gin.Context) {
// 	var creds Credentials
// 	if err := c.ShouldBindJSON(&creds); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
// 		return
// 	}

// 	user, err := auth.AuthenticateUser(db, creds.Email, creds.Password)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
// 		return
// 	}

// 	expirationTime := time.Now().Add(24 * time.Hour)
// 	claims := &Claims{
// 		Username: user.Username,
// 		Role:     user.Role,
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: expirationTime.Unix(),
// 		},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tokenString, err := token.SignedString(jwtKey)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"token": tokenString})
// }

// // CreatePost godoc
// // @Summary Create a post
// // @Description Create a new post
// // @Tags posts
// // @Accept json
// // @Produce json
// // @Param post body models.Post true "Post content"
// // @Success 200 {object} gin.H
// // @Failure 400 {object} gin.H
// // @Failure 401 {object} gin.H
// // @Router /api/posts [post]
// // @Security ApiKeyAuth
// func CreatePost(c *gin.Context) {
// 	var post models.Post
// 	if err := c.ShouldBindJSON(&post); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
// 		return
// 	}

// 	username := c.MustGet("username").(string)
// 	row := db.QueryRow("SELECT id FROM users WHERE username = ?", username)
// 	err := row.Scan(&post.UserID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find user"})
// 		return
// 	}

// 	result, err := db.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)", post.UserID, post.Title, post.Content)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create post"})
// 		return
// 	}

// 	postID, err := result.LastInsertId()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve post ID"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"post_id": postID})
// }

// // UpdatePost godoc
// // @Summary Update a post
// // @Description Update a post by ID
// // @Tags posts
// // @Accept json
// // @Produce json
// // @Param id path int true "Post ID"
// // @Param post body models.Post true "Post content"
// // @Success 200 {object} gin.H
// // @Failure 400 {object} gin.H
// // @Failure 401 {object} gin.H
// // @Failure 403 {object} gin.H
// // @Failure 404 {object} gin.H
// // @Router /api/posts/{id} [put]
// // @Security ApiKeyAuth
// func UpdatePost(c *gin.Context) {
// 	id := c.Param("id")
// 	var post models.Post
// 	if err := c.ShouldBindJSON(&post); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
// 		return
// 	}

// 	username := c.MustGet("username").(string)
// 	row := db.QueryRow("SELECT id FROM users WHERE username = ?", username)
// 	err := row.Scan(&post.UserID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find user"})
// 		return
// 	}

// 	res, err := db.Exec("UPDATE posts SET title = ?, content = ? WHERE id = ? AND user_id = ?", post.Title, post.Content, id, post.UserID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update post"})
// 		return
// 	}

// 	rowsAffected, err := res.RowsAffected()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve affected rows"})
// 		return
// 	}

// 	if rowsAffected == 0 {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found or unauthorized"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
// }

// // DeletePost godoc
// // @Summary Delete a post
// // @Description Delete a post by ID
// // @Tags posts
// // @Accept json
// // @Produce json
// // @Param id path int true "Post ID"
// // @Success 200 {object} gin.H
// // @Failure 400 {object} gin.H
// // @Failure 401 {object} gin.H
// // @Failure 403 {object} gin.H
// // @Failure 404 {object} gin.H
// // @Router /api/posts/{id} [delete]
// // @Security ApiKeyAuth
// func DeletePost(c *gin.Context) {
// 	id := c.Param("id")
// 	username := c.MustGet("username").(string)

// 	var userID int
// 	row := db.QueryRow("SELECT id FROM users WHERE username = ?", username)
// 	err := row.Scan(&userID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find user"})
// 		return
// 	}

// 	res, err := db.Exec("DELETE FROM posts WHERE id = ? AND user_id = ?", id, userID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete post"})
// 		return
// 	}

// 	rowsAffected, err := res.RowsAffected()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve affected rows"})
// 		return
// 	}

// 	if rowsAffected == 0 {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found or unauthorized"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
// }

// // GetPost godoc
// // @Summary Get a post by ID
// // @Description Get post details by ID
// // @Tags posts
// // @Accept json
// // @Produce json
// // @Param id path int true "Post ID"
// // @Success 200 {object} models.Post
// // @Failure 400 {object} gin.H
// // @Failure 404 {object} gin.H
// // @Router /api/posts/{id} [get]
// func GetPost(c *gin.Context) {
// 	id := c.Param("id")

// 	var post models.Post
// 	row := db.QueryRow("SELECT id, user_id, title, content, created_at FROM posts WHERE id = ?", id)
// 	err := row.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve post"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, post)
// }

// // GetAllUsers godoc
// // @Summary Get all users
// // @Description Get all registered users
// // @Tags users
// // @Accept json
// // @Produce json
// // @Success 200 {array} models.User
// // @Failure 401 {object} gin.H
// // @Failure 403 {object} gin.H
// // @Router /api/users [get]
// // @Security ApiKeyAuth
// func GetAllUsers(c *gin.Context) {
// 	rows, err := db.Query("SELECT id, username, role FROM users")
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve users"})
// 		return
// 	}
// 	defer rows.Close()

// 	var users []models.User
// 	for rows.Next() {
// 		var user models.User
// 		if err := rows.Scan(&user.ID, &user.Username, &user.Role); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not scan user data"})
// 			return
// 		}
// 		users = append(users, user)
// 	}

// 	c.JSON(http.StatusOK, users)
// }

// // DeleteUser godoc
// // @Summary Delete a user
// // @Description Delete a user by ID
// // @Tags users
// // @Accept json
// // @Produce json
// // @Param id path int true "User ID"
// // @Success 200 {object} gin.H
// // @Failure 400 {object} gin.H
// // @Failure 401 {object} gin.H
// // @Failure 403 {object} gin.H
// // @Failure 404 {object} gin.H
// // @Router /api/users/{id} [delete]
// // @Security ApiKeyAuth
// func DeleteUser(c *gin.Context) {
// 	id := c.Param("id")

// 	res, err := db.Exec("DELETE FROM users WHERE id = ?", id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user"})
// 		return
// 	}

// 	rowsAffected, err := res.RowsAffected()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve affected rows"})
// 		return
// 	}

// 	if rowsAffected == 0 {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
// }

// // UpdateUserRole godoc
// // @Summary Update a user's role
// // @Description Update a user's role by ID
// // @Tags users
// // @Accept json
// // @Produce json
// // @Param id path int true "User ID"
// // @Param role body string true "New role"
// // @Success 200 {object} gin.H
// // @Failure 400 {object} gin.H
// // @Failure 401 {object} gin.H
// // @Failure 403 {object} gin.H
// // @Failure 404 {object} gin.H
// // @Router /api/users/{id}/role [put]
// // @Security ApiKeyAuth
// func UpdateUserRole(c *gin.Context) {
// 	id := c.Param("id")
// 	var requestBody struct {
// 		Role string `json:"role"`
// 	}

// 	if err := c.ShouldBindJSON(&requestBody); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
// 		return
// 	}

// 	if requestBody.Role != "admin" && requestBody.Role != "user" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
// 		return
// 	}

// 	res, err := db.Exec("UPDATE users SET role = ? WHERE id = ?", requestBody.Role, id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user role"})
// 		return
// 	}

// 	rowsAffected, err := res.RowsAffected()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve affected rows"})
// 		return
// 	}

// 	if rowsAffected == 0 {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "User role updated successfully"})
// }

// // InitializeRoutes initializes all the routes
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

package handlers

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"literary-lions/backend/src/internal/auth"
	"literary-lions/backend/src/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var db *sql.DB
var jwtKey = []byte("my_secret_key")

// type Credentials struct {
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

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
}

// IsAuthorized checks the token validity and role
func IsAuthorized(handlerFunc gin.HandlerFunc, requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims.Role != requiredRole && claims.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		handlerFunc(c)
	}
}

// AuthMiddleware is a middleware to check authorization token and role
// func AuthMiddleware(requiredRole string, handlerFunc gin.HandlerFunc) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if !strings.HasPrefix(authHeader, "Bearer ") {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
// 			c.Abort()
// 			return
// 		}

// 		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

// 		claims := &Claims{}
// 		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
// 			return jwtKey, nil
// 		})

// 		if err != nil || !token.Valid {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 			c.Abort()
// 			return
// 		}

// 		if claims.Role != requiredRole && claims.Role != "admin" {
// 			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
// 			c.Abort()
// 			return
// 		}

// 		c.Set("username", claims.Username)
// 		c.Set("role", claims.Role)

//			handlerFunc(c)
//		}
//	}
//
// AuthMiddleware checks for token validity and role
// func AuthMiddleware(requiredRole string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.GetHeader("Authorization")
// 		if !strings.HasPrefix(authHeader, "Bearer ") {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
// 			c.Abort()
// 			return
// 		}

// 		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
// 		claims := &Claims{}
// 		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
// 			return jwtKey, nil
// 		})

// 		if err != nil || !token.Valid {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 			c.Abort()
// 			return
// 		}

// 		if claims.Role != requiredRole && claims.Role != "admin" {
// 			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
// 			c.Abort()
// 			return
// 		}

// 		c.Set("username", claims.Username)
// 		c.Set("role", claims.Role)

//			c.Next()
//		}
//	}
func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims.Role != requiredRole && claims.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// Login godoc
// @Summary Login a user
// @Description Login a user
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body Credentials true "User credentials"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /login [post]
func Login(c *gin.Context) {
	var creds Credentials
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	user, err := auth.AuthenticateUser(db, creds.Email, creds.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: user.Username,
		Role:     user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
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

	if err := auth.RegisterUser(db, req.Email, req.Username, req.Password); err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed"})
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

// CreatePost godoc
// @Summary Create a new post
// @Description Create a new post
// @Tags posts
// @Accept json
// @Produce json
// @Param post body models.Post true "Post object"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /api/posts [post]
// @Security ApiKeyAuth
func CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	username := c.GetString("username")
	userID, err := auth.GetUserIDByUsername(db, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve user information"})
		return
	}

	query := "INSERT INTO posts (user_id, category_id, title, content, created_at) VALUES ($1, $2, $3, $4, $5)"
	if _, err := db.Exec(query, userID, post.CategoryID, post.Title, post.Content, time.Now()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post created successfully"})
}

// GetAllPosts godoc
// @Summary Get all posts
// @Description Get all posts
// @Tags posts
// @Accept json
// @Produce json
// @Success 200 {array} models.Post
// @Failure 401 {object} gin.H
// @Router /api/posts [get]
// @Security ApiKeyAuth
func GetAllPosts(c *gin.Context) {
	rows, err := db.Query("SELECT id, user_id, category_id, title, content, created_at FROM posts")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve posts"})
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.CategoryID, &post.Title, &post.Content, &post.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not scan post data"})
			return
		}
		posts = append(posts, post)
	}

	c.JSON(http.StatusOK, posts)
}

// GetPost godoc
// @Summary Get a post by ID
// @Description Get a post by ID
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} models.Post
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/posts/{id} [get]
func GetPost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post

	query := "SELECT id, user_id, category_id, title, content, created_at FROM posts WHERE id = $1"
	err := db.QueryRow(query, id).Scan(&post.ID, &post.UserID, &post.CategoryID, &post.Title, &post.Content, &post.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve post"})
		}
		return
	}

	c.JSON(http.StatusOK, post)
}

// GetPostByID godoc
// @Summary Get a post by ID
// @Description Get a post by ID
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} models.Post
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/posts/{id} [get]
// @Security ApiKeyAuth
func GetPostByID(c *gin.Context) {
	id := c.Param("id")
	var post models.Post

	query := "SELECT id, user_id, category_id, title, content, created_at FROM posts WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&post.ID, &post.UserID, &post.CategoryID, &post.Title, &post.Content, &post.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve post"})
		}
		return
	}

	c.JSON(http.StatusOK, post)
}

// UpdatePost godoc
// @Summary Update a post
// @Description Update a post
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param post body models.Post true "Updated post object"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/posts/{id} [put]
// @Security ApiKeyAuth
func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	query := "UPDATE posts SET category_id = $1, title = $2, content = $3 WHERE id = $4"
	result, err := db.Exec(query, post.CategoryID, post.Title, post.Content, id)
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
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/posts/{id} [delete]
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

func InitializeRoutes(router *gin.Engine) {
	// User authentication routes
	router.POST("/register", Register)
	router.POST("/login", Login)

	// Protected routes
	api := router.Group("/api")
	{
		// Apply AuthMiddleware directly as middleware
		api.POST("/posts", AuthMiddleware("user"), CreatePost)
		api.PUT("/posts/:id", AuthMiddleware("user"), UpdatePost)
		api.DELETE("/posts/:id", AuthMiddleware("user"), DeletePost)
		api.GET("/users", AuthMiddleware("admin"), GetAllUsers)
		api.DELETE("/users/:id", AuthMiddleware("admin"), DeleteUser)
		api.PUT("/users/:id/role", AuthMiddleware("admin"), UpdateUserRole)
	}
}
