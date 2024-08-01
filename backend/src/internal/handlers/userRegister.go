package handlers

import (
	"net/http"

	"literary-lions/backend/src/internal/models"

	"github.com/gin-gonic/gin"
)

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
