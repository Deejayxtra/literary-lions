package models

import (
	"database/sql"
	"log"
	"time"
)

type Post struct {
	ID        int
	UserID    int
	Title     string
	Content   string
	Username  string
	Category  string
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
			log.Print("Empty post")
			return []Post{}, nil // Return an empty slice if no posts are found
		}
		log.Print("Error from db Post ", err.Error())
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

// ValidateSession retrieves the user ID from a session token and checks if the session is still valid.
// Parameters:
//   - token: The session token to validate.
//
// Returns:
//   - int: The user ID associated with the session if valid.
//   - error: An error if the session is invalid or if any other issue occurs; otherwise, nil.

// func ValidateSession(token string) (int, error) {
// 	var userID int
// 	var expiresAt time.Time

// 	// Query to get the user ID and expiration time for the given token
// 	err := db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE token = ?", token).Scan(&userID, &expiresAt)
// 	if err != nil {
// 		return 0, err
// 	}

// 	// Check if the current time is past the expiration time of the session
// 	if time.Now().After(expiresAt) {
// 		return 0, errors.New("session expired")
// 	}

// 	return userID, nil
// }
