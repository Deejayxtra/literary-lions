package models

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"-"`
	Role     string `json:"role"`
}

type Category struct {
	ID   int
	Name string
}

type Post struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	CategoryID int    `json:"category_id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
}

type Comment struct {
	ID        int    `json:"id"`
	PostID    int    `json:"post_id"`
	UserID    int    `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type Like struct {
	ID        int  `json:"id"`
	PostID    int  `json:"post_id"`
	CommentID int  `json:"comment_id"`
	UserID    int  `json:"user_id"`
	IsLike    bool `json:"is_like"`
}
