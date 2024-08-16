package models

import (
	"time"
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
	Username  string
	Email     string
}

type AuthResponse struct {
	Success  bool	`json:"success"`
	Token    string `json:"token"`
	Message  string	`json:"message"`
	Username string `json:"username"`
	Email    string	`json:"email"`
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
	FormattedContent []string	`json:"formatted_content"`
	
}

// Post struct represents a post in the forum.
type PostDetails struct {
	Post     Post
	Comments []Comment
	Likes 	  int
	Dislikes  int
	Status	  int
	Content   string  
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
