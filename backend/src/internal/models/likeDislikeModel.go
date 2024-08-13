package models

import "database/sql"

type PostLike struct {
	ID        int  `json:"id"`
	PostID    int  `json:"post_id"`
	CommentID int  `json:"comment_id"`
	UserID    int  `json:"user_id"`
	IsLike    bool `json:"is_like"`
}

type PostDislike struct {
	ID        	int   `json:"id"`
	PostID    	int   `json:"post_id"`
	CommentID 	int   `json:"comment_id"`
	UserID    	int   `json:"user_id"`
	IsDislike   bool  `json:"is_dislike"`
}

type CommentLike struct {
	ID        int  `json:"id"`
	PostID    int  `json:"post_id"`
	CommentID int  `json:"comment_id"`
	UserID    int  `json:"user_id"`
	IsLike    bool `json:"is_like"`
}

type CommentDislike struct {
	ID        int  `json:"id"`
	PostID    int  `json:"post_id"`
	CommentID int  `json:"comment_id"`
	UserID    int  `json:"user_id"`
	IsDislike bool `json:"is_dislike"`
}

// PostLikeAndUnlike adds or updates a like for a post. 
// If the user has already liked the post, it will remove the like.
// If the user has disliked the post, it will remove the dislike and add a like instead.
// Parameters:
//   - userID: The ID of the user performing the action.
//   - postID: The ID of the post to like or unlike.
//
// Returns:
//   - error: An error if the operation fails; otherwise, nil.
func PostLikeAndUnlike(userID, postID int) error {
	// Check if the user has already liked the post.
	var like PostLike
	row := db.QueryRow("SELECT id, user_id, post_id, is_like FROM post_likes WHERE post_id = ? AND user_id = ?", postID, userID)

	err := row.Scan(&like.ID, &like.UserID, &like.PostID, &like.IsLike)
	if err != nil {
		if err == sql.ErrNoRows {
			// Insert a new like if no existing like is found.
			_, err = db.Exec("INSERT INTO post_likes (user_id, post_id, is_like) VALUES (?, ?, ?)", userID, postID, true)
			RemoveDislikeFromPost(userID, postID) // Remove any existing dislike
			return err
		}
		return err
	}

	// If the post is already liked, remove the like; otherwise, add a like.
	if like.IsLike {
		_, err := db.Exec("UPDATE post_likes SET is_like = ? WHERE post_id = ? AND user_id = ?", false, postID, userID)
		return err
	} else {
		_, err := db.Exec("UPDATE post_likes SET is_like = ? WHERE post_id = ? AND user_id = ?", true, postID, userID)
		RemoveDislikeFromPost(userID, postID) // Remove any existing dislike
		return err
	}
}

// PostDisLikeAndUndislike adds or updates a dislike for a post. 
// If the user has already disliked the post, it will remove the dislike.
// If the user has liked the post, it will remove the like and add a dislike instead.
// Parameters:
//   - userID: The ID of the user performing the action.
//   - postID: The ID of the post to dislike or undislike.
//
// Returns:
//   - error: An error if the operation fails; otherwise, nil.
func PostDisLikeAndUndislike(userID, postID int) error {
	// Check if the user has already disliked the post.
	var dislike PostDislike
	row := db.QueryRow("SELECT id, user_id, post_id, is_dislike FROM post_dislikes WHERE post_id = ? AND user_id = ?", postID, userID)

	err := row.Scan(&dislike.ID, &dislike.UserID, &dislike.PostID, &dislike.IsDislike)
	if err != nil {
		if err == sql.ErrNoRows {
			// Insert a new dislike if no existing dislike is found.
			_, err = db.Exec("INSERT INTO post_dislikes (user_id, post_id, is_dislike) VALUES (?, ?, ?)", userID, postID, true)
			RemoveLikeFromPost(userID, postID) // Remove any existing like
			return err
		}
		return err
	}

	// If the post is already disliked, remove the dislike; otherwise, add a dislike.
	if dislike.IsDislike {
		_, err := db.Exec("UPDATE post_dislikes SET is_dislike = ? WHERE post_id = ? AND user_id = ?", false, postID, userID)
		return err
	} else {
		_, err := db.Exec("UPDATE post_dislikes SET is_dislike = ? WHERE post_id = ? AND user_id = ?", true, postID, userID)
		RemoveLikeFromPost(userID, postID) // Remove any existing like
		return err
	}
}

// RemoveDislikeFromPost removes a dislike for a post.
// Parameters:
//   - userID: The ID of the user performing the action.
//   - postID: The ID of the post to remove dislike from.
//
// Returns:
//   - error: An error if the operation fails; otherwise, nil.
func RemoveDislikeFromPost(userID, postID int) error {
	// Check if the user has disliked the post.
	var dislike PostDislike
	row := db.QueryRow("SELECT id, user_id, post_id, is_dislike FROM post_dislikes WHERE post_id = ? AND user_id = ?", postID, userID)

	err := row.Scan(&dislike.ID, &dislike.UserID, &dislike.PostID, &dislike.IsDislike)
	if err != nil {
		if err == sql.ErrNoRows {
			// No need to remove dislike if it doesn't exist.
			return err
		}
		return err
	}

	// If the post is already disliked, update to remove the dislike.
	if dislike.IsDislike {
		_, err := db.Exec("UPDATE post_dislikes SET is_dislike = ? WHERE post_id = ? AND user_id = ?", false, postID, userID)
		return err
	}
	return err
}

// RemoveLikeFromPost removes a like for a post.
// Parameters:
//   - userID: The ID of the user performing the action.
//   - postID: The ID of the post to remove like from.
//
// Returns:
//   - error: An error if the operation fails; otherwise, nil.
func RemoveLikeFromPost(userID, postID int) error {
	// Check if the user has liked the post.
	var like PostLike
	row := db.QueryRow("SELECT id, user_id, post_id, is_like FROM post_likes WHERE post_id = ? AND user_id = ?", postID, userID)

	err := row.Scan(&like.ID, &like.UserID, &like.PostID, &like.IsLike)
	if err != nil {
		if err == sql.ErrNoRows {
			// No need to remove like if it doesn't exist.
			return err
		}
		return err
	}

	// If the post is already liked, update to remove the like.
	if like.IsLike {
		_, err := db.Exec("UPDATE post_likes SET is_like = ? WHERE post_id = ? AND user_id = ?", false, postID, userID)
		return err
	}
	return err
}

// CountPostLikes returns the total number of likes for a specific post.
// Parameters:
//   - postID: The ID of the post for which likes are being counted.
//
// Returns:
//   - int: The number of likes for the specified post.
//   - error: An error if the operation fails; otherwise, nil.
func CountPostLikes(postID int) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM post_likes WHERE post_id = ? AND is_like = 1"
	err := db.QueryRow(query, postID).Scan(&count)
	return count, err
}

// CountPostDislikes returns the total number of dislikes for a specific post.
// Parameters:
//   - postID: The ID of the post for which dislikes are being counted.
//
// Returns:
//   - int: The number of dislikes for the specified post.
//   - error: An error if the operation fails; otherwise, nil.
func CountPostDislikes(postID int) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM post_dislikes WHERE post_id = ? AND is_dislike = 1"
	err := db.QueryRow(query, postID).Scan(&count)
	return count, err
}

// CountCommentLikes returns the total number of likes for a specific comment.
// Parameters:
//   - commentID: The ID of the comment for which likes are being counted.
//
// Returns:
//   - int: The number of likes for the specified comment.
//   - error: An error if the operation fails; otherwise, nil.
func CountCommentLikes(commentID int) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM comment_likes WHERE comment_id = ? AND is_like = 1"
	err := db.QueryRow(query, commentID).Scan(&count)
	return count, err
}

// CountCommentDislikes returns the total number of dislikes for a specific comment.
// Parameters:
//   - commentID: The ID of the comment for which dislikes are being counted.
//
// Returns:
//   - int: The number of dislikes for the specified comment.
//   - error: An error if the operation fails; otherwise, nil.
func CountCommentDislikes(commentID int) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM comment_dislikes WHERE comment_id = ? AND is_dislike = 1"
	err := db.QueryRow(query, commentID).Scan(&count)
	return count, err
}

// CommentLikeAndUnlike adds or updates a like for a comment. 
// If the user has already liked the comment, it will remove the like.
// If the user has disliked the comment, it will remove the dislike and add a like instead.
// Parameters:
//   - userID: The ID of the user performing the action.
//   - commentID: The ID of the comment to like or unlike.
//
// Returns:
//   - error: An error if the operation fails; otherwise, nil.
func CommentLikeAndUnlike(userID, commentID int) error {
	// Check if the user has already liked the comment.
	var like CommentLike
	row := db.QueryRow("SELECT id, user_id, comment_id, is_like FROM comment_likes WHERE comment_id = ? AND user_id = ?", commentID, userID)

	err := row.Scan(&like.ID, &like.UserID, &like.CommentID, &like.IsLike)
	if err != nil {
		if err == sql.ErrNoRows {
			// Insert a new like if no existing like is found.
			_, err = db.Exec("INSERT INTO comment_likes (user_id, comment_id, is_like) VALUES (?, ?, ?)", userID, commentID, true)
			RemoveDislikeFromComment(userID, commentID) // Remove any existing dislike
			return err
		}
		return err
	}

	// If the comment is already liked, remove the like; otherwise, add a like.
	if like.IsLike {
		_, err := db.Exec("UPDATE comment_likes SET is_like = ? WHERE comment_id = ? AND user_id = ?", false, commentID, userID)
		return err
	} else {
		_, err := db.Exec("UPDATE comment_likes SET is_like = ? WHERE comment_id = ? AND user_id = ?", true, commentID, userID)
		RemoveDislikeFromComment(userID, commentID) // Remove any existing dislike
		return err
	}
}

// CommentDisLikeAndUndislike adds or updates a dislike for a comment. 
// If the user has already disliked the comment, it will remove the dislike.
// If the user has liked the comment, it will remove the like and add a dislike instead.
// Parameters:
//   - userID: The ID of the user performing the action.
//   - commentID: The ID of the comment to dislike or undislike.
//
// Returns:
//   - error: An error if the operation fails; otherwise, nil.
func CommentDisLikeAndUndislike(userID, commentID int) error {
	// Check if the user has already disliked the comment.
	var dislike CommentDislike
	row := db.QueryRow("SELECT id, user_id, comment_id, is_dislike FROM comment_dislikes WHERE comment_id = ? AND user_id = ?", commentID, userID)

	err := row.Scan(&dislike.ID, &dislike.UserID, &dislike.CommentID, &dislike.IsDislike)
	if err != nil {
		if err == sql.ErrNoRows {
			// Insert a new dislike if no existing dislike is found.
			_, err = db.Exec("INSERT INTO comment_dislikes (user_id, comment_id, is_dislike) VALUES (?, ?, ?)", userID, commentID, true)
			RemoveLikeFromComment(userID, commentID) // Remove any existing like
			return err
		}
		return err
	}

	// If the comment is already disliked, remove the dislike; otherwise, add a dislike.
	if dislike.IsDislike {
		_, err := db.Exec("UPDATE comment_dislikes SET is_dislike = ? WHERE comment_id = ? AND user_id = ?", false, commentID, userID)
		return err
	} else {
		_, err := db.Exec("UPDATE comment_dislikes SET is_dislike = ? WHERE comment_id = ? AND user_id = ?", true, commentID, userID)
		RemoveLikeFromComment(userID, commentID) // Remove any existing like
		return err
	}
}

// RemoveDislikeFromComment removes a dislike for a comment.
// Parameters:
//   - userID: The ID of the user performing the action.
//   - commentID: The ID of the comment to remove dislike from.
//
// Returns:
//   - error: An error if the operation fails; otherwise, nil.
func RemoveDislikeFromComment(userID, commentID int) error {
	// Check if the user has disliked the comment.
	var dislike CommentDislike
	row := db.QueryRow("SELECT id, user_id, comment_id, is_dislike FROM comment_dislikes WHERE comment_id = ? AND user_id = ?", commentID, userID)

	err := row.Scan(&dislike.ID, &dislike.UserID, &dislike.CommentID, &dislike.IsDislike)
	if err != nil {
		if err == sql.ErrNoRows {
			// No need to remove dislike if it doesn't exist.
			return err
		}
		return err
	}

	// If the comment is already disliked, update to remove the dislike.
	if dislike.IsDislike {
		_, err := db.Exec("UPDATE comment_dislikes SET is_dislike = ? WHERE comment_id = ? AND user_id = ?", false, commentID, userID)
		return err
	}
	return err
}

// RemoveLikeFromComment removes a like for a comment.
// Parameters:
//   - userID: The ID of the user performing the action.
//   - commentID: The ID of the comment to remove like from.
//
// Returns:
//   - error: An error if the operation fails; otherwise, nil.
func RemoveLikeFromComment(userID, commentID int) error {
	// Check if the user has liked the comment.
	var like CommentLike
	row := db.QueryRow("SELECT id, user_id, comment_id, is_like FROM comment_likes WHERE comment_id = ? AND user_id = ?", commentID, userID)

	err := row.Scan(&like.ID, &like.UserID, &like.CommentID, &like.IsLike)
	if err != nil {
		if err == sql.ErrNoRows {
			// No need to remove like if it doesn't exist.
			return err
		}
		return err
	}

	// If the comment is already liked, update to remove the like.
	if like.IsLike {
		_, err := db.Exec("UPDATE comment_likes SET is_like = ? WHERE comment_id = ? AND user_id = ?", false, commentID, userID)
		return err
	}
	return err
}
