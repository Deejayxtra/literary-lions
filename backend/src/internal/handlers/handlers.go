package handlers

import (
	"database/sql"
	"literary-lions/backend/src/internal/models"
	"log"
	"net/http"
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

// InitHandlers initializes the handlers by setting up the database connection.
// It sets the global database variable and configures the models package to use this database.
func InitHandlers(database *sql.DB) {
	db = database               // Set the global db variable to the provided database connection
	models.SetDatabase(db)       // Configure the models package to use this database connection
}

// UpdatePost godoc
// @Summary Update a post
// @Description Update an existing post by ID
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
	id := c.Param("id")  // Retrieve the post ID from the URL path

	var post models.Post
	// Bind the incoming JSON payload to the post variable
	if err := c.ShouldBindJSON(&post); err != nil {
		// If the request payload is invalid, return a bad request error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Update the post in the database with the new data
	query := "UPDATE posts SET category_id = $1, title = $2, content = $3 WHERE id = $4"
	result, err := db.Exec(query, post.Category, post.Title, post.Content, id)
	if err != nil {
		// If the update fails, return an internal server error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update post"})
		return
	}

	// Check how many rows were affected by the update
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		// If no rows were affected, the post was not found, return a 404 error
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Return a success message if the post was updated successfully
	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
}

// DeletePost godoc
// @Summary Delete a post
// @Description Delete a post by ID
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
	id := c.Param("id")  // Retrieve the post ID from the URL path

	// Delete the post from the database
	query := "DELETE FROM posts WHERE id = $1"
	result, err := db.Exec(query, id)
	if err != nil {
		// If the deletion fails, return an internal server error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete post"})
		return
	}

	// Check how many rows were affected by the deletion
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		// If no rows were affected, the post was not found, return a 404 error
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Return a success message if the post was deleted successfully
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Retrieve all users from the database
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Failure 401 {object} gin.H
// @Router /api/users [get]
// @Security ApiKeyAuth
func GetAllUsers(c *gin.Context) {
	// Query the database to retrieve all users
	rows, err := db.Query("SELECT id, username, email, role FROM users")
	if err != nil {
		// If the query fails, return an internal server error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve users"})
		return
	}
	defer rows.Close()  // Ensure the rows are closed after use

	var users []models.User
	// Iterate over the rows and scan the user data into the user struct
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role); err != nil {
			// If scanning fails, return an internal server error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not scan user data"})
			return
		}
		// Append the user to the users slice
		users = append(users, user)
	}

	// Return the list of users in the response
	c.JSON(http.StatusOK, users)
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Retrieve a user from the database by ID
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
	id := c.Param("id")  // Retrieve the user ID from the URL path
	var user models.User

	// Query the database to retrieve the user by ID
	query := "SELECT id, username, email, role FROM users WHERE id = $1"
	err := db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			// If no rows were found, return a 404 error
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			// If the query fails, return an internal server error
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve user"})
		}
		return
	}

	// Return the user data in the response
	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update an existing user by ID
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
	id := c.Param("id")  // Retrieve the user ID from the URL path
	var user models.User

	// Bind the incoming JSON payload to the user variable
	if err := c.ShouldBindJSON(&user); err != nil {
		// If the request payload is invalid, return a bad request error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Update the user in the database with the new data
	query := "UPDATE users SET username = $1, email = $2, role = $3 WHERE id = $4"
	result, err := db.Exec(query, user.Username, user.Email, user.Role, id)
	if err != nil {
		// If the update fails, return an internal server error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user"})
		return
	}

	// Check how many rows were affected by the update
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		// If no rows were affected, the user was not found, return a 404 error
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Return a success message if the user was updated successfully
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by ID
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
	id := c.Param("id")  // Retrieve the user ID from the URL path

	// Delete the user from the database
	query := "DELETE FROM users WHERE id = $1"
	result, err := db.Exec(query, id)
	if err != nil {
		// If the deletion fails, return an internal server error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user"})
		return
	}

	// Check how many rows were affected by the deletion
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		// If no rows were affected, the user was not found, return a 404 error
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Return a success message if the user was deleted successfully
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
	id := c.Param("id")  // Retrieve the user ID from the URL path
	var requestBody struct {
		Role string `json:"role"`
	}

	// Bind the incoming JSON payload to the requestBody variable
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		// If the request payload is invalid, return a bad request error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Check if the provided role is valid
	if requestBody.Role != "admin" && requestBody.Role != "user" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	// Update the user's role in the database
	res, err := db.Exec("UPDATE users SET role = ? WHERE id = ?", requestBody.Role, id)
	if err != nil {
		// If the update fails, return an internal server error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user role"})
		return
	}

	// Check how many rows were affected by the update
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		// If there's an error retrieving affected rows, return an internal server error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve affected rows"})
		return
	}

	// If no rows were affected, the user was not found, return a 404 error
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Return a success message if the user's role was updated successfully
	c.JSON(http.StatusOK, gin.H{"message": "User role updated successfully"})
}

// UpdateUserProfile handles the user profile update request
func UpdateUserProfile(c *gin.Context) {
	// Retrieve the user ID from the context (assuming it's set by the middleware)
	userID, exists := c.Get("userID")
	if !exists {
		log.Println("Unauthorized access attempt: no userID found in context")
		// If the user ID is not found, return an unauthorized error
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Define the data structure for the expected input
	var data struct {
		Email     string `json:"email" binding:"required"`
		Username  string `json:"username" binding:"required"`
	}

	// Bind the incoming JSON data to the data struct
	if err := c.ShouldBindJSON(&data); err != nil {
		// If there's an error in binding, return a bad request status with the error message
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the ProfileUpdate function to update the user's profile in the database
	err := models.ProfileUpdate(userID.(int), data.Email, data.Username)
	if err != nil {
		// If the update fails, return an internal server error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the updated user data from the database
	user, err := models.GetUser(userID.(int))
	if err != nil {
		// If there's an error retrieving the user, return an internal server error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Return a success message with the updated user data
	c.JSON(http.StatusOK, gin.H{"username": user.Username, "email": user.Email})
}
