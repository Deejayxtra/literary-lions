// package handlers

// import (
// 	"net/http"
// 	"strconv"

// 	"literary-lions/backend/src/internal/models"

// 	"github.com/gin-gonic/gin"
// )

// // LikePost godoc
// // @Summary Like a post
// // @Description Like a post
// // @Tags likes
// // @Accept json
// // @Produce json
// // @Param id path int true "Post ID"
// // @Success 200 {object} gin.H
// // @Failure 400 {object} gin.H
// // @Failure 401 {object} gin.H
// // @Router /api/post/{id}/like [post]
// // @Security ApiKeyAuth
// func LikePost(c *gin.Context) {
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	postIDStr := c.Param("id")
// 	postID, err := strconv.Atoi(postIDStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
// 		return
// 	}

// 	// Add the like
// 	err = models.CreateLike(userID.(int), postID, 0, true)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Get the updated like count
// 	likeCount, err := models.CountLikes(postID, 0)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve like count"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Post liked successfully", "likes": likeCount})
// }

// // DislikePost godoc
// // @Summary Dislike a post
// // @Description Dislike a post
// // @Tags dislikes
// // @Accept json
// // @Produce json
// // @Param id path int true "Post ID"
// // @Success 200 {object} gin.H
// // @Failure 400 {object} gin.H
// // @Failure 401 {object} gin.H
// // @Router /api/post/{id}/dislike [post]
// // @Security ApiKeyAuth
// func DislikePost(c *gin.Context) {
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	postIDStr := c.Param("id")
// 	postID, err := strconv.Atoi(postIDStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
// 		return
// 	}

// 	// Add the dislike
// 	err = models.CreateLike(userID.(int), postID, 0, false)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Get the updated dislike count
// 	dislikeCount, err := models.CountDislikes(postID, 0)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve dislike count"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Post disliked successfully", "dislikes": dislikeCount})
// }

// // LikeComment godoc
// // @Summary Like a comment
// // @Description Like a comment
// // @Tags likes
// // @Accept json
// // @Produce json
// // @Param id path int true "Comment ID"
// // @Success 200 {object} gin.H
// // @Failure 400 {object} gin.H
// // @Failure 401 {object} gin.H
// // @Router /api/comment/{id}/like [post]
// // @Security ApiKeyAuth
// func LikeComment(c *gin.Context) {
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	commentIDStr := c.Param("id")
// 	commentID, err := strconv.Atoi(commentIDStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
// 		return
// 	}

// 	// Add the like
// 	err = models.CreateLike(userID.(int), 0, commentID, true)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Get the updated like count
// 	likeCount, err := models.CountLikes(0, commentID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve like count"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Comment liked successfully", "likes": likeCount})
// }

// // DislikeComment godoc
// // @Summary Dislike a comment
// // @Description Dislike a comment
// // @Tags dislikes
// // @Accept json
// // @Produce json
// // @Param id path int true "Comment ID"
// // @Success 200 {object} gin.H
// // @Failure 400 {object} gin.H
// // @Failure 401 {object} gin.H
// // @Router /api/comment/{id}/dislike [post]
// // @Security ApiKeyAuth
// func DislikeComment(c *gin.Context) {
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	commentIDStr := c.Param("id")
// 	commentID, err := strconv.Atoi(commentIDStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
// 		return
// 	}

// 	// Add the dislike
// 	err = models.CreateLike(userID.(int), 0, commentID, false)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Get the updated dislike count
// 	dislikeCount, err := models.CountDislikes(0, commentID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve dislike count"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Comment disliked successfully", "dislikes": dislikeCount})
// }

// // UnlikePost godoc
// // @Summary Remove like from a post
// // @Description Remove like from a post
// // @Tags likes
// // @Accept json
// // @Produce json
// // @Param id path int true "Post ID"
// // @Success 200 {object} gin.H
// // @Failure 400 {object} gin.H
// // @Failure 401 {object} gin.H
// // @Router /api/post/{id}/unlike [delete]
// // @Security ApiKeyAuth
// func UnlikePost(c *gin.Context) {
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	postIDStr := c.Param("id")
// 	postID, err := strconv.Atoi(postIDStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
// 		return
// 	}

// 	// Remove the like
// 	err = models.RemoveLike(userID.(int), postID, 0)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Get the updated like count
// 	likeCount, err := models.CountLikes(postID, 0)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve like count"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Post unliked successfully", "likes": likeCount})
// }

// // UndislikePost godoc
// // @Summary Remove dislike from a post
// // @Description Remove dislike from a post
// // @Tags dislikes
// // @Accept json
// // @Produce json
// // @Param id path int true "Post ID"
// // @Success 200 {object} gin.H
// // @Failure 400 {object} gin.H
// // @Failure 401 {object} gin.H
// // @Router /api/post/{id}/undislike [delete]
// // @Security ApiKeyAuth
// func UndislikePost(c *gin.Context) {
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	postIDStr := c.Param("id")
// 	postID, err := strconv.Atoi(postIDStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
// 		return
// 	}

// 	// Remove the dislike
// 	err = models.RemoveLike(userID.(int), postID, 0)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Get the updated dislike count
// 	dislikeCount, err := models.CountDislikes(postID, 0)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve dislike count"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Post undisliked successfully", "dislikes": dislikeCount})
// }

// // UnlikeComment godoc
// // @Summary Remove like from a comment
// // @Description Remove like from a comment
// // @Tags likes
// // @Accept json
// // @Produce json
// // @Param id path int true "Comment ID"
// // @Success 200 {object} gin.H
// // @Failure 400 {object} gin.H
// // @Failure 401 {object} gin.H
// // @Router /api/comment/{id}/unlike [delete]
// // @Security ApiKeyAuth
// func UnlikeComment(c *gin.Context) {
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	commentIDStr := c.Param("id")
// 	commentID, err := strconv.Atoi(commentIDStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
// 		return
// 	}

// 	// Remove the like
// 	err = models.RemoveLike(userID.(int), 0, commentID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Get the updated like count
// 	likeCount, err := models.CountLikes(0, commentID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve like count"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Comment unliked successfully", "likes": likeCount})
// }

// // UndislikeComment godoc
// // @Summary Remove dislike from a comment
// // @Description Remove dislike from a comment
// // @Tags dislikes
// // @Accept json
// // @Produce json
// // @Param id path int true "Comment ID"
// // @Success 200 {object} gin.H
// // @Failure 400 {object} gin.H
// // @Failure 401 {object} gin.H
// // @Router /api/comment/{id}/undislike [delete]
// // @Security ApiKeyAuth
// func UndislikeComment(c *gin.Context) {
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	commentIDStr := c.Param("id")
// 	commentID, err := strconv.Atoi(commentIDStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
// 		return
// 	}

// 	// Remove the dislike
// 	err = models.RemoveLike(userID.(int), 0, commentID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	// Get the updated dislike count
// 	dislikeCount, err := models.CountDislikes(0, commentID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve dislike count"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Comment undisliked successfully", "dislikes": dislikeCount})
// }

package handlers

import (
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

	// Add the like
	err = models.CreateLike(userID.(int), postID, 0, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the updated like count
	// likeCount, err := models.CountLikes(postID, 0)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve like count"})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{"message": "Post liked successfully", "likes": likeCount})
	c.JSON(http.StatusOK, gin.H{"message": "Request successfully"})
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

	// Add the dislike
	err = models.CreateLike(userID.(int), postID, 0, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the updated dislike count
	dislikeCount, err := models.CountDislikes(postID, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve dislike count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post disliked successfully", "dislikes": dislikeCount})
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
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	// Add the like
	err = models.CreateLike(userID.(int), 0, commentID, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the updated like count
	likeCount, err := models.CountLikes(0, commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve like count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment liked successfully", "likes": likeCount})
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
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	// Add the dislike
	err = models.CreateLike(userID.(int), 0, commentID, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the updated dislike count
	dislikeCount, err := models.CountDislikes(0, commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve dislike count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment disliked successfully", "dislikes": dislikeCount})
}

// UnlikePost godoc
// @Summary Remove like from a post
// @Description Remove like from a post
// @Tags likes
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /api/post/{id}/unlike [delete]
// @Security ApiKeyAuth
func UnlikePost(c *gin.Context) {
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

	// Remove the like
	err = models.RemoveLike(userID.(int), postID, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the updated like count
	likeCount, err := models.CountLikes(postID, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve like count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post unliked successfully", "likes": likeCount})
}

// UndislikePost godoc
// @Summary Remove dislike from a post
// @Description Remove dislike from a post
// @Tags dislikes
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /api/post/{id}/undislike [delete]
// @Security ApiKeyAuth
func UndislikePost(c *gin.Context) {
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

	// Remove the dislike
	err = models.RemoveLike(userID.(int), postID, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the updated dislike count
	dislikeCount, err := models.CountDislikes(postID, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve dislike count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post undisliked successfully", "dislikes": dislikeCount})
}

// UnlikeComment godoc
// @Summary Remove like from a comment
// @Description Remove like from a comment
// @Tags likes
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /api/comment/{id}/unlike [delete]
// @Security ApiKeyAuth
func UnlikeComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	commentIDStr := c.Param("id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	// Remove the like
	err = models.RemoveLike(userID.(int), 0, commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the updated like count
	likeCount, err := models.CountLikes(0, commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve like count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment unliked successfully", "likes": likeCount})
}

// UndislikeComment godoc
// @Summary Remove dislike from a comment
// @Description Remove dislike from a comment
// @Tags dislikes
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /api/comment/{id}/undislike [delete]
// @Security ApiKeyAuth
func UndislikeComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	commentIDStr := c.Param("id")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	// Remove the dislike
	err = models.RemoveLike(userID.(int), 0, commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the updated dislike count
	dislikeCount, err := models.CountDislikes(0, commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve dislike count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment undisliked successfully", "dislikes": dislikeCount})
}
