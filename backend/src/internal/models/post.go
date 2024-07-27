package models

import (
	"database/sql"
	// "errors"
	"time"
)

type Post struct {
    ID        int
    UserID    int
    Title     string
    Content   string
    Category  string
    CreatedAt time.Time
}

func CreatePost(userID int, title, content, category string) error {
    _, err := db.Exec("INSERT INTO posts (user_id, title, content, category) VALUES (?, ?, ?, ?)", userID, title, content, category)
    return err
}

func GetAllPosts(db *sql.DB) ([]Post, error) {
    rows, err := db.Query("SELECT id, title, content, category, user_id FROM posts")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var post Post
        if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.UserID); err != nil {
            return nil, err
        }
        posts = append(posts, post)
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