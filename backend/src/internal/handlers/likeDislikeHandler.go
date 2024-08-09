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
// @Description Like a post
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
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    postIDStr := c.Param("id")
    postID, err := strconv.Atoi(postIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
        return
    }

    // Update the likes table for POST
    err = models.PostLikeAndUnlike(userID.(int), postID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Request successfully processed"})
}

// DislikePost godoc
// @Summary Dislike a post
// @Description Dislike a post
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
	    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    postIDStr := c.Param("id")
    postID, err := strconv.Atoi(postIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
        return
    }

    // Update the likes table for POST
    err = models.PostDisLikeAndUndislike(userID.(int), postID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Request successfully processed"})
}

// LikeComment godoc
// @Summary Like a comment
// @Description Like a comment
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
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    commentIDStr := c.Param("id")
    postID, err := strconv.Atoi(commentIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
        return
    }

    // Update the likes table for POST
    err = models.CommentLikeAndUnlike(userID.(int), postID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Request successfully processed"})
}

// DislikeComment godoc
// @Summary Dislike a comment
// @Description Dislike a comment
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
	    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    commentIDStr := c.Param("id")
    postID, err := strconv.Atoi(commentIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
        return
    }

    // Update the likes table for POST
    err = models.CommentDisLikeAndUndislike(userID.(int), postID)
    if err != nil {
        log.Print("error: ", err.Error())
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Request successfully processed"})
}

