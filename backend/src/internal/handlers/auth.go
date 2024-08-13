package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"literary-lions/backend/src/internal/models"
)

// AuthMiddleware is a middleware function that checks if the user is authenticated
// and optionally checks if the user has the required role.
// If the user is not authenticated or doesn't have the required role, the request is aborted.
func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve the session token from the cookie
		cookie, err := c.Cookie("session_token")
		if err != nil {
			// If the session token is missing or invalid, return an unauthorized error
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
			c.Abort() // Abort the request, no further handlers will be called
			return
		}

		// Validate the session token and retrieve the associated user ID
		token := cookie
		userID, err := models.ValidateSession(token)
		if err != nil || userID == 0 {
			// If the session is invalid or the user ID is 0, return an unauthorized error
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
			c.Abort() // Abort the request, no further handlers will be called
			return
		}

		// Store the user ID in the context for further use in the request lifecycle
		c.Set("userID", userID)
		c.Next() // Continue to the next handler in the chain
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
// Login handles user authentication and session creation.
func Login(c *gin.Context) {
	var creds Credentials
	// Bind the incoming JSON payload to the creds variable
	if err := c.ShouldBindJSON(&creds); err != nil {
		// If the request payload is invalid, return a bad request error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Authenticate the user with the provided email and password
	user, err := models.AuthenticateUser(creds.Email, creds.Password)
	if err != nil {
		// If authentication fails, return an unauthorized error
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	// Create a session token for the authenticated user
	token, err := models.CreateSession(user.ID)
	if err != nil {
		// If session creation fails, return an internal server error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create session"})
		return
	}

	// Set the session token as a cookie in the response
	c.SetCookie("session_token", token, 24*3600, "/", "", false, true)

	// Respond with the session token and user details
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
// Logout handles user logout and session invalidation.
func Logout(c *gin.Context) {
	// Retrieve the session token from the cookie
	token, err := c.Cookie("session_token")
	if err != nil {
		// If the session token is missing, return an unauthorized error
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Invalidate the session in the database
	if err := models.InvalidateSession(token); err != nil {
		// If session invalidation fails, return an internal server error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log out"})
		return
	}

	// Clear the session token cookie
	c.SetCookie("session_token", "", -1, "/", "", false, true)

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
