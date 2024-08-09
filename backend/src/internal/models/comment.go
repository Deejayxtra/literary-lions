package models

import (
	"time"
)

type Comment struct {
	ID        int
	PostID    int
	UserID    int
	Username  string
	Likes	  int
	Dislikes  int
	Content   string
	CreatedAt time.Time `json:"created_at" db:"createdAt"`
}

func CreateComment(postID, userID int, content string) error {
	_, err := db.Exec("INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)", postID, userID, content)
	return err
}

func GetCommentsByPostID(postID int) ([]Comment, error) {
	rows, err := db.Query("SELECT id, post_id, user_id, content, created_at FROM comments WHERE post_id = ?", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}

		// To fetch the Username for the comment
		user, err := GetUser(comment.UserID) 
		if err != nil {
			return nil, err
		}
		comment.Username = user.Username

		// To fetch the total number of likes for the comment
		likes, err := CountCommentLikes(comment.ID)
		if err != nil {
			return nil, err
		}
		comment.Likes = likes

		// To fetch the total number of dislikes for the comment
		dislikes, err := CountCommentDislikes(comment.ID)
		if err != nil {
			return nil, err
		}
		comment.Dislikes = dislikes

		comments = append(comments, comment)

	}


	return comments, nil
}
