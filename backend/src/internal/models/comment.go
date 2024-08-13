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

// CreateComment inserts a new comment into the database.
// Parameters:
//   - postID: The ID of the post that the comment is associated with.
//   - userID: The ID of the user who is creating the comment.
//   - content: The text content of the comment.
//
// Returns:
//   - error: An error if the operation fails; otherwise, nil.
func CreateComment(postID, userID int, content string) error {
	// Execute the SQL command to insert a new comment into the 'comments' table.
	_, err := db.Exec("INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)", postID, userID, content)
	return err
}

// GetCommentsByPostID retrieves all comments associated with a specific post from the database.
// Parameters:
//   - postID: The ID of the post for which comments are being fetched.
//
// Returns:
//   - []Comment: A slice of Comment structs containing all comments for the specified post.
//   - error: An error if the operation fails; otherwise, nil.
func GetCommentsByPostID(postID int) ([]Comment, error) {
	// Query the database for all comments associated with the given post ID.
	rows, err := db.Query("SELECT id, post_id, user_id, content, created_at FROM comments WHERE post_id = ?", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Ensure rows are closed when done to avoid resource leaks.

	var comments []Comment
	for rows.Next() {
		var comment Comment
		// Scan the row into the Comment struct fields.
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}

		// Fetch the Username associated with the comment's user ID.
		user, err := GetUser(comment.UserID)
		if err != nil {
			return nil, err
		}
		comment.Username = user.Username

		// Fetch the total number of likes for the comment.
		likes, err := CountCommentLikes(comment.ID)
		if err != nil {
			return nil, err
		}
		comment.Likes = likes

		// Fetch the total number of dislikes for the comment.
		dislikes, err := CountCommentDislikes(comment.ID)
		if err != nil {
			return nil, err
		}
		comment.Dislikes = dislikes

		// Append the fully populated Comment struct to the slice.
		comments = append(comments, comment)
	}

	// Return the slice of comments and a nil error indicating success.
	return comments, nil
}
