package models

type User struct {
    ID       int
    Email    string
    Username string
    Password string
}

type Category struct {
    ID   int
    Name string
}

type Post struct {
    ID         int
    UserID     int
    CategoryID int
    Title      string
    Content    string
    CreatedAt  string
}

type Comment struct {
    ID        int
    PostID    int
    UserID    int
    Content   string
    CreatedAt string
}

type Like struct {
    ID        int
    PostID    int
    CommentID int
    UserID    int
    IsLike    bool
}
