package handlers

import (
	"database/sql"
	"literary-lions/backend/src/internal/models"
	"literary-lions/backend/src/internal/utils"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Validate the email format
	if !isValidEmail(creds.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Open a database connection
	db, err := sql.Open("sqlite3", "literary_lions.db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the database"})
		return
	}
	defer db.Close()

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin database transaction"})
		return
	}
	defer tx.Rollback()

	// Check if the user exists
	user, err := models.FindUserByEmail(tx, creds.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User does not exist"})
		return
	}

	// Check if the password is correct using the utils package
	if !utils.CheckPassword(user.Password, creds.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
		return
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Create a session token for the authenticated user
	token, err := models.CreateSession(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create session"})
		return
	}

	// Set the session token as a cookie in the response
	c.SetCookie("session_token", token, 24*3600, "/", "", false, true)

	// Respond with the session token and user details
	c.JSON(http.StatusOK, gin.H{"token": token, "username": user.Username, "email": user.Email})
}

func isValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
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
