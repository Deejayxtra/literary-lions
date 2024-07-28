package models

import (
	"errors"
)

type Like struct {
	ID        int  `json:"id"`
	PostID    int  `json:"post_id"`
	CommentID int  `json:"comment_id"`
	UserID    int  `json:"user_id"`
	IsLike    bool `json:"is_like"`
}

// CreateLike adds a like or dislike for a post or comment
func CreateLike(userID, postID, commentID int, isLike bool) error {
	// Check if the like/dislike already exists
	var count int
	query := "SELECT COUNT(*) FROM likes WHERE user_id = ? AND post_id = ? AND comment_id = ?"
	err := db.QueryRow(query, userID, postID, commentID).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("already liked or disliked")
	}

	// Insert the like/dislike
	_, err = db.Exec("INSERT INTO likes (user_id, post_id, commentID, is_like) VALUES (?, ?, ?, ?)", userID, postID, commentID, isLike)
	return err
}

// RemoveLike removes a like or dislike for a post or comment
func RemoveLike(userID, postID, commentID int) error {
	_, err := db.Exec("DELETE FROM likes WHERE user_id = ? AND post_id = ? AND comment_id = ?", userID, postID, commentID)
	return err
}

// CountLikes returns the total number of likes for a post or comment
func CountLikes(postID, commentID int) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM likes WHERE post_id = ? AND comment_id = ? AND is_like = 1"
	err := db.QueryRow(query, postID, commentID).Scan(&count)
	return count, err
}

// CountDislikes returns the total number of dislikes for a post or comment
func CountDislikes(postID, commentID int) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM likes WHERE post_id = ? AND comment_id = ? AND is_like = 0"
	err := db.QueryRow(query, postID, commentID).Scan(&count)
	return count, err
}
