package handlers

import (
	"literary-lions/backend/src/internal/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// AddComment godoc
// @Summary Add a new comment
// @Description Add a new comment to a specific post
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

    // Retrieve the user ID from the context (set by middleware)
    userID, exists := c.Get("userID")
    if !exists {
        // If the user is not authenticated, return an unauthorized error
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // Get the post ID from the URL parameter and convert it to an integer
    commentIDStr := c.Param("id")
    postID, err := strconv.Atoi(commentIDStr)
    if err != nil {
        // If the post ID is invalid, return a bad request error
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
        return
    }

    // Define a structure to bind the incoming JSON request
    var comment struct {
        Content string `json:"content"`
    }

    // Bind the JSON request body to the comment struct
    if err := c.ShouldBindJSON(&comment); err != nil {
        // If the request body is invalid, return a bad request error
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Call the function to create the comment in the database
    if err := models.CreateComment(postID, userID.(int), comment.Content); err != nil {
        // If the operation fails, return an internal server error
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Return a success message if the comment was added successfully
    c.JSON(http.StatusCreated, gin.H{"message": "Comment added successfully"})
}

// CreatePost godoc
// @Summary Create a new post
// @Description Create a new post with title, content, and category
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
    // Retrieve the user ID from the context (set by middleware)
    userID, exists := c.Get("userID")
    if !exists {
        // If the user is not authenticated, return an unauthorized error
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // Define a structure to bind the incoming JSON request
    var post struct {
        Title    string `json:"title" binding:"required"`
        Content  string `json:"content" binding:"required"`
        Category string `json:"category" binding:"required"`
    }

    // Bind the JSON request body to the post struct
    if err := c.ShouldBindJSON(&post); err != nil {
        // If the request body is invalid, return a bad request error
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
        return
    }

    // Call the function to create the post in the database
    err := models.CreatePost(userID.(int), post.Title, post.Content, post.Category)
    if err != nil {
        // If the operation fails, return an internal server error
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Return a success message if the post was created successfully
    c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully"})
}

// GetAllPosts godoc
// @Summary Get all posts
// @Description Retrieve all posts from the database
// @Tags posts
// @Accept json
// @Produce json
// @Success 200 {array} models.Post
// @Failure 401 {object} gin.H
// @Router /api/posts [get]
// @Security ApiKeyAuth
// GetAllPosts handles the retrieval of all posts using Gin
// GetAllPosts handles the retrieval of all posts with optional advanced search filters.
func GetAllPosts(c *gin.Context) {
	// Retrieve query parameters for filtering
	title := c.Query("keyword")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	category := c.Query("category")

	// Parse the start and end dates if provided
	var parsedStartDate, parsedEndDate time.Time
	var err error
	if startDate != "" {
		parsedStartDate, err = time.Parse("2006-01-02", startDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use YYYY-MM-DD."})
			return
		}
	}
	if endDate != "" {
		parsedEndDate, err = time.Parse("2006-01-02", endDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use YYYY-MM-DD."})
			return
		}
	}

	// Call the function to get all posts from the database with the provided filters
	posts, err := models.GetFilteredPosts(category, title, parsedStartDate, parsedEndDate)
	if err != nil {
        log.Print("Err: ", err.Error())
		// If the operation fails, return an internal server error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the list of posts as a JSON response
	c.JSON(http.StatusOK, posts)
}

// GetPost godoc
// @Summary Get a post by ID
// @Description Retrieve a single post by its ID along with comments, likes, and dislikes
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
    // Get the post ID from the URL parameter and convert it to an integer
    postIDStr := c.Param("id")
    postID, err := strconv.Atoi(postIDStr)
    if err != nil {
        // If the post ID is invalid, return a bad request error
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
        return
    }

    // Call the function to get the post by ID from the database
    post, err := models.GetPostByID(postID)
    if err != nil {
        // If the operation fails, return an internal server error
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Call the function to get comments associated with the post
    comments, err := models.GetCommentsByPostID(postID)
    if err != nil {
        // If the operation fails, return an internal server error
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Call the functions to count likes and dislikes for the post
    likes, err := models.CountPostLikes(postID)
    if err != nil {
        // If the operation fails, return an internal server error
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    dislikes, err := models.CountPostDislikes(postID)
    if err != nil {
        // If the operation fails, return an internal server error
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Prepare the response with the post, comments, likes, and dislikes
    response := struct {
        Post     models.Post     `json:"post"`
        Comments []models.Comment `json:"comments"`
        Likes    int              `json:"likes"`
        Dislikes int              `json:"dislikes"`
    }{
        Post:     post,
        Comments: comments,
        Likes:    likes,
        Dislikes: dislikes,
    }

    // Return the response as a JSON object
    c.JSON(http.StatusOK, response)
}

// GetPostsByCategory godoc
// @Summary Get posts by category
// @Description Retrieve all posts that belong to a specific category
// @Tags posts
// @Accept json
// @Produce json
// @Param category path string true "Category"
// @Success 200 {array} models.Post
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/posts/category/{category} [get]
// GetPostsByCategory handles the retrieval of posts filtered by category using Gin
func GetPostsByCategory(c *gin.Context) {
    // Call the function to get all posts from the database
    posts, err := models.GetAllPosts(db)
    if err != nil {
        // If the operation fails, return an internal server error
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Get the category from the URL parameter
    category := c.Param("category")
    var filteredPosts []models.Post

    // Filter the posts by the specified category
    for _, post := range posts {
        if post.Category == category {
            filteredPosts = append(filteredPosts, post)
        }
    }

    // If no posts are found in the specified category, return a not found error
    if len(filteredPosts) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"message": "No posts found for this category"})
        return
    }

    // Return the filtered posts as a JSON response
    c.JSON(http.StatusOK, filteredPosts)
}
