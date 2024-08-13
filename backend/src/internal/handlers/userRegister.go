package handlers

import (
	"net/http"
	"literary-lions/backend/src/internal/models"
	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email, username, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegistrationRequest true "User registration request"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Router /register [post]
// Register handles the user registration process
func Register(c *gin.Context) {
    // Define a structure to bind the incoming JSON request for registration
    var req RegistrationRequest

    // Bind the JSON request body to the RegistrationRequest struct
    // Returns a bad request error if the JSON is invalid or missing required fields
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }

    // Call the function to register the user with the provided details
    // If an error occurs during registration, return an internal server error
    if err := models.RegisterUser(req.Email, req.Username, req.Password); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Return a success message if the registration is successful
    c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}
