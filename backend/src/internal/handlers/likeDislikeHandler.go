package handlers

import (
	"log"
	"net/http"
	"strconv"
	"literary-lions/backend/src/internal/models"
	"github.com/gin-gonic/gin"
)

// LikePost godoc
// @Summary Like a post
// @Description Like a post by its ID
// @Tags likes
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /api/post/{id}/like [post]
// @Security ApiKeyAuth
func LikePost(c *gin.Context) {
    // Retrieve the user ID from the context (assuming it's set by the middleware)
    userID, exists := c.Get("userID")
    if !exists {
        // If the user ID is not found, return an unauthorized error
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // Get the post ID from the URL parameter and convert it to an integer
    postIDStr := c.Param("id")
    postID, err := strconv.Atoi(postIDStr)
    if err != nil {
        // If the post ID is invalid, return a bad request error
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
        return
    }

    // Call the function to like or unlike the post
    err = models.PostLikeAndUnlike(userID.(int), postID)
    if err != nil {
        // If the operation fails, return an internal server error
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Return a success message if the operation was successful
    c.JSON(http.StatusOK, gin.H{"message": "Request successfully processed"})
}

// DislikePost godoc
// @Summary Dislike a post
// @Description Dislike a post by its ID
// @Tags dislikes
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /api/post/{id}/dislike [post]
// @Security ApiKeyAuth
func DislikePost(c *gin.Context) {
    // Retrieve the user ID from the context (assuming it's set by the middleware)
    userID, exists := c.Get("userID")
    if !exists {
        // If the user ID is not found, return an unauthorized error
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // Get the post ID from the URL parameter and convert it to an integer
    postIDStr := c.Param("id")
    postID, err := strconv.Atoi(postIDStr)
    if err != nil {
        // If the post ID is invalid, return a bad request error
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
        return
    }

    // Call the function to dislike or undislike the post
    err = models.PostDisLikeAndUndislike(userID.(int), postID)
    if err != nil {
        // If the operation fails, return an internal server error
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Return a success message if the operation was successful
    c.JSON(http.StatusOK, gin.H{"message": "Request successfully processed"})
}

// LikeComment godoc
// @Summary Like a comment
// @Description Like a comment by its ID
// @Tags likes
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /api/comment/{id}/like [post]
// @Security ApiKeyAuth
func LikeComment(c *gin.Context) {
    // Retrieve the user ID from the context (assuming it's set by the middleware)
    userID, exists := c.Get("userID")
    if !exists {
        // If the user ID is not found, return an unauthorized error
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // Get the comment ID from the URL parameter and convert it to an integer
    commentIDStr := c.Param("id")
    commentID, err := strconv.Atoi(commentIDStr)
    if err != nil {
        // If the comment ID is invalid, return a bad request error
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
        return
    }

    // Call the function to like or unlike the comment
    err = models.CommentLikeAndUnlike(userID.(int), commentID)
    if err != nil {
        // If the operation fails, return an internal server error
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Return a success message if the operation was successful
    c.JSON(http.StatusOK, gin.H{"message": "Request successfully processed"})
}

// DislikeComment godoc
// @Summary Dislike a comment
// @Description Dislike a comment by its ID
// @Tags dislikes
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /api/comment/{id}/dislike [post]
// @Security ApiKeyAuth
func DislikeComment(c *gin.Context) {
    // Retrieve the user ID from the context (assuming it's set by the middleware)
    userID, exists := c.Get("userID")
    if !exists {
        // If the user ID is not found, return an unauthorized error
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // Get the comment ID from the URL parameter and convert it to an integer
    commentIDStr := c.Param("id")
    commentID, err := strconv.Atoi(commentIDStr)
    if err != nil {
        // If the comment ID is invalid, return a bad request error
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
        return
    }

    // Call the function to dislike or undislike the comment
    err = models.CommentDisLikeAndUndislike(userID.(int), commentID)
    if err != nil {
        // If the operation fails, log the error and return an internal server error
        log.Print("error: ", err.Error())
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Return a success message if the operation was successful
    c.JSON(http.StatusOK, gin.H{"message": "Request successfully processed"})
}
