package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"literary-lions/backend/src/internal/models"
)

func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("session_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
			c.Abort()
			return
		}

		token := cookie
		userID, err := models.ValidateSession(token)
		if err != nil || userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
			c.Abort()
			return
		}

		// Store user ID and role in context for further use
		c.Set("userID", userID)
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

	user, err := models.AuthenticateUser(creds.Email, creds.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	token, err := models.CreateSession(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create session"})
		return
	}

	// Set the session token as a cookie
	c.SetCookie("session_token", token, 24*3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"token": token, "username": user.Username, "email": user.Email})
}

// Logout godoc

// @Summary Logout a user
// @Description Logout a user
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body Credentials true "User credentials"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /logout [post]
// Logout handles user logout
func Logout(c *gin.Context) {
	// Retrieve the session token from the cookie
	token, err := c.Cookie("session_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Invalidate the session in the database
	if err := models.InvalidateSession(token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log out"})
		return
	}

	// Clear the session cookie
	c.SetCookie("session_token", "", -1, "/", "", false, true)

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
