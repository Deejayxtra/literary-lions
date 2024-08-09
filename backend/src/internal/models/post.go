package models

import (
	"database/sql"
	"log"

	// "errors"
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

func CreatePost(userID int, title, content, category string) error {
	_, err := db.Exec("INSERT INTO posts (user_id, title, content, category) VALUES (?, ?, ?, ?)", userID, title, content, category)
	return err
}

// GetAllPosts retrieves all posts sorted by the createdAt field in descending order
func GetAllPosts(db *sql.DB) ([]Post, error) {
	// Updated SQL query to include createdAt and order by createdAt DESC
	rows, err := db.Query("SELECT id, title, content, category, user_id, created_at FROM posts ORDER BY created_at DESC")
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print("Empty post")
			return []Post{}, nil //Retrun array of empty object
		}
		log.Print("Error from db Post ", err.Error())
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		// Scan the createdAt field into the Post struct
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.UserID, &post.CreatedAt); err != nil {
			return nil, err
		}

		// Fetch the username for each post
        user, err := GetUser(post.UserID)
        if err != nil {
            return nil, err
        }
        post.Username = user.Username

		posts = append(posts, post)
	}

	// Check for errors after looping through rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// GetPostByID retrieves a post by its ID from the database
func GetPostByID(postID int) (Post, error) {
	var post Post
	row := db.QueryRow("SELECT id, user_id, title, content, category, created_at FROM posts WHERE id = ?", postID)

	err := row.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Category, &post.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return Post{}, nil // Post not found
		}
		return Post{}, err
	}
	user, err := GetUser(post.UserID) 
	if err != nil {
		return Post{}, err
	}
	post.Username = user.Username
	

	return post, nil
}

// // ValidateSession retrieves the user ID from a session token
// func ValidateSession(token string) (int, error) {
//     var userID int
//     var expiresAt time.Time

//     err := db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE token = ?", token).Scan(&userID, &expiresAt)
//     if err != nil {
//         return 0, err
//     }

//     if time.Now().After(expiresAt) {
//         return 0, errors.New("session expired")
//     }

//     return userID, nil
// }
