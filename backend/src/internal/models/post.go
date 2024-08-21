package models

import (
	"database/sql"
	"strings"
	"time"
)

type Post struct {
	ID        int		`json:"id"`
	UserID    int		`json:"user_id"`
	Title     string	`json:"title"`
	Content   string	`json:"content"`
	Username  string	`json:"username"`
	Category  string	`json:"category"`
	CreatedAt time.Time `json:"created_at" db:"createdAt"`
}

// CreatePost inserts a new post into the database with the provided details.
// Parameters:
//   - userID: The ID of the user creating the post.
//   - title: The title of the post.
//   - content: The content of the post.
//   - category: The category of the post.
//
// Returns:
//   - error: An error if the operation fails; otherwise, nil.
func CreatePost(userID int, title, content, category string) error {
	_, err := db.Exec("INSERT INTO posts (user_id, title, content, category) VALUES (?, ?, ?, ?)", userID, title, content, category)
	return err
}

// GetAllPosts retrieves all posts from the database, sorted by the creation time in descending order.
// Parameters:
//   - db: The database connection to use for the query.
//
// Returns:
//   - []Post: A slice of Post structs containing the details of each post.
//   - error: An error if the operation fails; otherwise, nil.
func GetAllPosts(db *sql.DB) ([]Post, error) {
	// Query to select all posts, including the created_at field and sorting by it in descending order
	rows, err := db.Query("SELECT id, title, content, category, user_id, created_at FROM posts ORDER BY created_at DESC")
	if err != nil {
		if err == sql.ErrNoRows {

			return []Post{}, nil // Return an empty slice if no posts are found
		}

		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		// Scan each field into the Post struct, including created_at
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.UserID, &post.CreatedAt); err != nil {
			return nil, err
		}

		// Fetch the username for each post based on the user ID
        user, err := GetUser(post.UserID)
        if err != nil {
            return nil, err
        }
        post.Username = user.Username

		posts = append(posts, post)
	}

	// Check for errors encountered during the iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// GetPostByID retrieves a specific post by its ID from the database.
// Parameters:
//   - postID: The ID of the post to retrieve.
//
// Returns:
//   - Post: A Post struct containing the details of the requested post.
//   - error: An error if the operation fails; otherwise, nil.
func GetPostByID(postID int) (Post, error) {
	var post Post
	row := db.QueryRow("SELECT id, user_id, title, content, category, created_at FROM posts WHERE id = ?", postID)

	err := row.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Category, &post.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return Post{}, nil // Return an empty Post if not found
		}
		return Post{}, err
	}
	// Fetch the username for the post's author
	user, err := GetUser(post.UserID) 
	if err != nil {
		return Post{}, err
	}
	post.Username = user.Username
	
	return post, nil
}

// GetFilteredPosts retrieves posts from the database based on the provided filters.
func GetFilteredPosts(category, title string, startDate, endDate time.Time) ([]Post, error) {
	var posts []Post
	var filters []string
	var args []interface{}

	// Apply title filter
	if title != "" {
		filters = append(filters, "title LIKE ?")
		args = append(args, "%"+title+"%")
	}

	// Apply category filter
	if category != "" {
		filters = append(filters, "category LIKE ?")
		args = append(args, "%"+category+"%")
	}

	// Apply date range filter
	if !startDate.IsZero() {
		filters = append(filters, "created_at >= ?")
		args = append(args, startDate)
	}
	if !endDate.IsZero() {
		filters = append(filters, "created_at <= ?")
		args = append(args, endDate)
	}

	// Build the query
	query := "SELECT id, user_id, title, content, category, created_at FROM posts"
	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}
	query += " ORDER BY created_at DESC" // Order by most recent posts

	// Execute the query
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and scan the data into the posts slice
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Category, &post.CreatedAt); err != nil {
			return nil, err
		}

		// Fetch the username for each post based on the user ID
        user, err := GetUser(post.UserID)
        if err != nil {
            return nil, err
        }
        post.Username = user.Username

		posts = append(posts, post)
	}

	return posts, nil
}

// GetUserPosts fetches all posts created by the given user.
func GetUserPosts(userID int) ([]Post, error) {
	// Query to select posts by user ID
	rows, err := db.Query("SELECT id, title, content, category, user_id, created_at FROM posts WHERE user_id = ? ORDER BY created_at DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.UserID, &post.CreatedAt); err != nil {
			return nil, err
		}

		// Set the Username field directly, as all posts belong to the same user
		user, err := GetUser(post.UserID) // Assuming `GetUser` is a function that fetches user details based on `userID`
		if err != nil {
			return nil, err
		}
		post.Username = user.Username

		posts = append(posts, post)
	}

	// Check for any errors encountered during the iteration over the rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func GetLikedPostsByUserID(userID int) ([]Post, error) {
    // SQL query to select liked posts along with the username
    query := `
        SELECT p.id, p.title, p.content, p.category, p.user_id, p.created_at, u.username 
        FROM posts p
        INNER JOIN post_likes pl ON p.id = pl.post_id
        INNER JOIN users u ON p.user_id = u.id
        WHERE pl.user_id = ? AND pl.is_like = 1
    `

    // Execute the query
    rows, err := db.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Slice to hold the result
    var likedPosts []Post

    // Iterate over the rows
    for rows.Next() {
        var post Post
        // Assuming Post struct has a Username field
        if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.UserID, &post.CreatedAt, &post.Username); err != nil {
            return nil, err
        }
        likedPosts = append(likedPosts, post)
    }

    // Check for any errors encountered during iteration
    if err := rows.Err(); err != nil {
        return nil, err
    }

    return likedPosts, nil
}
