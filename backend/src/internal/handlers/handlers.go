package handlers

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"fmt"
)

var db *sql.DB
var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
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

		if claims.Role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this resource"})
			c.Abort()
			return
		}

		// Pass the claims to the handler
		c.Set("claims", claims)

		handlerFunc(c)
	}
}

// @Summary Register a new user
// @Description Register a new user with email, username, and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   credentials  body      Credentials  true  "User credentials"
// @Success 200 {string} string "User registered successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /register [post]
func Register(c *gin.Context) {
	var creds Credentials
	err := c.BindJSON(&creds)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
		fmt.Printf("hashedpassword: %s\n", hashedPassword)
	// _, err = db.Exec("INSERT INTO users (email, username, password, role) VALUES (?, ?, ?, 'user')", creds.Email, creds.Username, hashedPassword)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// @Summary Login a user
// @Description Login a user with username and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   credentials  body      Credentials  true  "User credentials"
// @Success 200 {string} string "User logged in successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /login [post]
func Login(c *gin.Context) {
	var creds Credentials
	err := c.BindJSON(&creds)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var storedCreds Credentials
	var role string
	err = db.QueryRow("SELECT email, username, password, role FROM users WHERE username=?", creds.Username).Scan(&storedCreds.Email, &storedCreds.Username, &storedCreds.Password, &role)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: creds.Username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	c.JSON(http.StatusOK, gin.H{"message": "User logged in successfully"})
}

// @Summary Create a new post
// @Description Create a new post with title, content, and category_id
// @Tags posts
// @Accept  json
// @Produce  json
// @Param   post  body      map[string]interface{}  true  "Post details"
// @Success 201 {object} map[string]interface{} "Post created successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /posts [post]
// @Security ApiKeyAuth
func CreatePost(c *gin.Context) {
	claims := c.MustGet("claims").(*Claims)
	var post struct {
		CategoryID int    `json:"category_id"`
		Title      string `json:"title"`
		Content    string `json:"content"`
	}
	err := c.BindJSON(&post)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	_, err = db.Exec("INSERT INTO posts (category_id, title, content, user_id) VALUES (?, ?, ?, (SELECT id FROM users WHERE username = ?))", post.CategoryID, post.Title, post.Content, claims.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully"})
}

// @Summary Update a post by ID
// @Description Update a post by ID with title and content
// @Tags posts
// @Accept  json
// @Produce  json
// @Param   id    path      string                  true  "Post ID"
// @Param   post  body      map[string]interface{}  true  "Post details"
// @Success 200 {object} map[string]interface{} "Post updated successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /posts/{id} [put]
// @Security ApiKeyAuth
func UpdatePost(c *gin.Context) {
	claims := c.MustGet("claims").(*Claims)
	postID := c.Param("id")
	var post struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	err := c.BindJSON(&post)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	_, err = db.Exec("UPDATE posts SET title = ?, content = ? WHERE id = ? AND user_id = (SELECT id FROM users WHERE username = ?)", post.Title, post.Content, postID, claims.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
}

// @Summary Delete a post by ID
// @Description Delete a post by ID
// @Tags posts
// @Produce  json
// @Param   id  path      string  true  "Post ID"
// @Success 200 {object} map[string]interface{} "Post deleted successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /posts/{id} [delete]
// @Security ApiKeyAuth
func DeletePost(c *gin.Context) {
	claims := c.MustGet("claims").(*Claims)
	postID := c.Param("id")

	_, err := db.Exec("DELETE FROM posts WHERE id = ? AND user_id = (SELECT id FROM users WHERE username = ?)", postID, claims.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Post deleted successfully",
		"post_id":  postID,
		"username": claims.Username,
	})
}

// @Summary Get a post by ID
// @Description Get a post by ID
// @Tags posts
// @Produce  json
// @Param   id  path      string  true  "Post ID"
// @Success 200 {object} map[string]interface{} "Post details"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /posts/{id} [get]
func GetPost(c *gin.Context) {
	postID := c.Param("id")
	var post struct {
		ID         int    `json:"id"`
		CategoryID int    `json:"category_id"`
		Title      string `json:"title"`
		Content    string `json:"content"`
		Username   string `json:"username"`
	}

	err := db.QueryRow("SELECT p.id, p.category_id, p.title, p.content, u.username FROM posts p JOIN users u ON p.user_id = u.id WHERE p.id = ?", postID).Scan(&post.ID, &post.CategoryID, &post.Title, &post.Content, &post.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// @Summary Get all users
// @Description Get all users
// @Tags users
// @Produce  json
// @Success 200 {array} map[string]interface{} "List of users"
// @Failure 500 {string} string "Internal server error"
// @Router /users [get]
// @Security ApiKeyAuth
func GetAllUsers(c *gin.Context) {
	rows, err := db.Query("SELECT id, email, username, role FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	defer rows.Close()

	users := []map[string]interface{}{}
	for rows.Next() {
		var user struct {
			ID       int    `json:"id"`
			Email    string `json:"email"`
			Username string `json:"username"`
			Role     string `json:"role"`
		}
		err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		users = append(users, map[string]interface{}{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
			"role":     user.Role,
		})
	}

	c.JSON(http.StatusOK, users)
}

// @Summary Delete a user by ID
// @Description Delete a user by ID
// @Tags users
// @Produce  json
// @Param   id  path      string  true  "User ID"
// @Success 200 {object} map[string]interface{} "User deleted successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /users/{id} [delete]
// @Security ApiKeyAuth
func DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	_, err := db.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
		"user_id": userID,
	})
}

// @Summary Update a user's role
// @Description Update a user's role by ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param   id    path      string                    true  "User ID"
// @Param   role  body      map[string]interface{}    true  "User role"
// @Success 200 {object} map[string]interface{} "User role updated successfully"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /users/{id}/role [put]
// @Security ApiKeyAuth
func UpdateUserRole(c *gin.Context) {
	userID := c.Param("id")
	var requestBody struct {
		Role string `json:"role"`
	}

	err := c.BindJSON(&requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	_, err = db.Exec("UPDATE users SET role = ? WHERE id = ?", requestBody.Role, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User role updated successfully",
		"user_id": userID,
		"role":    requestBody.Role,
	})
}
