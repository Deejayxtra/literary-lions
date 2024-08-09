package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResponseDetails struct {
	Success   bool
	Message   string
	Status    int
}

type AuthResponse struct {
	Success  bool
	Token    string `json:"token"`
	Message  string
	Username string `json:"username"`
}

// User struct represents a user in the system.
type User struct {
	ID       	int    		`json:"id"`
	Email    	string 		`json:"email"`
	Username 	string 		`json:"username"`
	Password 	string 		`json:"-"`
	ProfilePic 	string   	`json:"profile_pic"`
	CreatedAt  	time.Time 	`json:"created_at"`
	Role     	string 		`json:"role"`
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
	Category   string    `json:"category"`
	CategoryID int       `json:"category_id"`
	Title      string    `json:"title"`
	Username   string 	 `json:"username"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
	
}

// Post struct represents a post in the forum.
type PostDetails struct {
	Post     Post
	Comments []Comment
	Likes 	  int
	Dislikes  int
}

// Comment struct represents a comment on a post.
type Comment struct {
	ID        int       `json:"id"`
	PostID    string    `json:"post_id"`
	UserID    int       `json:"user_id"`
	Username  string 	`json:"username"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Likes 	  int		`json:"likes"`
	Dislikes  int		`json:"dislikes"`
}

// Like struct represents a like/dislike on a post or comment.
type Like struct {
	ID        int  `json:"id"`
	PostID    int  `json_id:"post_id"`
	CommentID int  `json_id:"comment_id"`
	UserID    int  `json_id:"user_id"`
	IsLike    bool `json:"is_like"`
}

// Message struct represents a message in a conversation room.
type Message struct {
	UserID  int       `json:"user_id"`
	Content string    `json:"content"`
	Time    time.Time `json:"time"`
}

// Data struct for template.
type Data struct {
	Posts         []Post
	Authenticated bool
}

// AuthenticateUser simulates user authentication.
func AuthenticateUser(email, password string) (*User, error) {
	// Example: hardcoded user for demonstration purposes
	if email == "test@example.com" {
		return &User{
			ID:       1,
			Email:    email,
			Username: "testuser",
			Password: "$2a$10$1yII3Pq/4FbDsZz5l4P2oOkKhCzI053GcP2LHKFvw1PeFNErc4Bd2", // bcrypt hash for "password123"
		}, nil
	}
	return nil, errors.New("user not found")
}

func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

// GetUserByEmail retrieves a user by email from your database.
func GetUserByEmail(email string) (*User, error) {
	// Implement this function to retrieve a user by email from your database
	// This is just a placeholder
	return &User{
		Email:    "test@example.com",
		Username: "testuser",
		// Password: "Ud21wC+n/y0I27JcwIEGRA==",
		Password: "$2a$10$1yII3Pq/4FbDsZz5l4P2oOkKhCzI053GcP2LHKFvw1PeFNErc4Bd2", // bcrypt hash of "password"
	}, nil
}
