package models

import (
    "time"
)

type Comment1 struct {
    ID        int
    PostID    int
    UserID    int
    Content   string
    CreatedAt time.Time
}

func CreateComment(postID, userID int, content string) error {
    _, err := db.Exec("INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)", postID, userID, content)
    return err
}

func GetCommentsByPostID(postID int) ([]Comment1, error) {
    rows, err := db.Query("SELECT id, post_id, user_id, content, created_at FROM comments WHERE post_id = ?", postID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var comments []Comment1
    for rows.Next() {
        var comment Comment1
        err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt)
        if err != nil {
            return nil, err
        }
        comments = append(comments, comment)
    }

    return comments, nil
}
