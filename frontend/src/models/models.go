package models

import (
    "errors"
    "time"
    "golang.org/x/crypto/bcrypt"

)

// User struct represents a user in the system.
type User struct {
    ID       int    `json:"id"`
    Email    string `json:"email"`
    Username string `json:"username"`
    Password string `json:"-"`
    Role     string `json:"role"`
}

// Category struct represents a category for posts.
type Category struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

// Post struct represents a post in the forum.
type Post struct {
    ID         int       `json:"id"`
    UserID     int       `json:"user_id"`
    CategoryID int       `json:"category_id"`
    Title      string    `json:"title"`
    Content    string    `json:"content"`
    CreatedAt  time.Time `json:"created_at"`
}

// Comment struct represents a comment on a post.
type Comment struct {
    ID        int       `json:"id"`
    PostID    int       `json:"post_id"`
    UserID    int       `json:"user_id"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
}

// Like struct represents a like/dislike on a post or comment.
type Like struct {
    ID        int  `json:"id"`
    PostID    int  `json:"post_id"`
    CommentID int  `json:"comment_id"`
    UserID    int  `json:"user_id"`
    IsLike    bool `json:"is_like"`
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func GetUserByEmail(email string) (*User, error) {
	// Implement this function to retrieve a user by email from your database
	// This is just a placeholder
	return &User{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "$2a$12$wJ89JKZa/nH/jf/Y0BZhKuGrOq1BF9N6ZOHYpDkqI9lRdfq9nWJ.e", // bcrypt hash of "password"
	}, nil
}

// AuthenticateUser simulates user authentication.
func AuthenticateUser(email, password string) (*User, error) {
    // Simulated authentication logic
    // Replace with actual authentication logic
    if email == "example@example.com" && password == "password" {
        return &User{
            ID:       1,
            Email:    email,
            Username: "example_user",
            Role:     "user",
        }, nil
    }
    return nil, errors.New("invalid credentials")
}

// CreateUser simulates user creation.
func CreateUser(email, username, password string) error {
    // Simulated user creation logic
    // Replace with actual user creation logic
    return nil
}

// CreatePost simulates post creation.
func PostComment(categoryID int, title, content string, userID int) error {
    // Simulated post creation logic
    // Replace with actual post creation logic
    return nil
}

// CreateComment simulates comment creation.
func CreateChannel(postID int, content string, userID int) error {
    // Simulated comment creation logic
    // Replace with actual comment creation logic
    return nil
}
