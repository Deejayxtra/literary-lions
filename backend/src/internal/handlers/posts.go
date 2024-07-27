package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"

	"literary-lions/backend/src/internal/models"
)

// AddComment godoc
// @Summary Add a new comment
// @Description Add a new comment
// @Tags comments
// @Accept json
// @Produce json
// @Param post body models.Post true "Post object"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /api/posts [post]
// @Security ApiKeyAuth
// AddComment handles adding a comment to a post using Gin
func AddComment(c *gin.Context) {
	
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

    var comment struct {
        Content string `json:"content"`
    }

    if err := c.ShouldBindJSON(&comment); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := models.CreateComment(postID, userID.(int), comment.Content); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Comment added successfully"})
}

// CreatePost godoc
// @Summary Create a new post
// @Description Create a new post
// @Tags posts
// @Accept json
// @Produce json
// @Param post body models.Post true "Post object"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /api/post [post]
// @Security ApiKeyAuth
func CreatePost(c *gin.Context) {
	
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    var post struct {
        Title    string `json:"title" binding:"required"`
        Content  string `json:"content" binding:"required"`
        Category string `json:"category" binding:"required"`
    }

    if err := c.ShouldBindJSON(&post); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }

    err := models.CreatePost(userID.(int), post.Title, post.Content, post.Category)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully"})
}


// GetAllPosts godoc
// @Summary Get all posts
// @Description Get all posts
// @Tags posts
// @Accept json
// @Produce json
// @Success 200 {array} models.Post
// @Failure 401 {object} gin.H
// @Router /api/posts [get]
// @Security ApiKeyAuth
// GetAllPosts handles the retrieval of all posts using Gin
func GetAllPosts(c *gin.Context) {
    posts, err := models.GetAllPosts(db)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, posts)
}

// GetPost godoc
// @Summary Get a post by ID
// @Description Get a post by ID
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} models.Post
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /api/post/{id} [get]
// GetPostByID handles the retrieval of a single post by ID using Gin
func GetPostByID(c *gin.Context) {
    postIDStr := c.Param("id")
    postID, err := strconv.Atoi(postIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
        return
    }

    post, err := models.GetPostByID(postID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    comments, err := models.GetCommentsByPostID(postID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    response := struct {
        Post     models.Post      `json:"post"`
        Comments []models.Comment1 `json:"comments"`
    }{
        Post:     post,
        Comments: comments,
    }

    c.JSON(http.StatusOK, response)
}