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

// PostLikeAndUnlike adds a like or dislike for a post or comment
func PostLikeAndUnlike(userID, postID int) error {
	// Check if the like already exists
	    var like PostLike
    row := db.QueryRow("SELECT id, user_id, post_id, is_like FROM post_likes WHERE post_id = ? AND user_id = ?", postID, userID)

    err := row.Scan(&like.ID, &like.UserID, &like.PostID, &like.IsLike)
    if err != nil {
        if err == sql.ErrNoRows {
            // Insert the like
            _, err = db.Exec("INSERT INTO post_likes (user_id, post_id, is_like) VALUES (?, ?, ?)", userID, postID, true)
            RemoveDislikeFromPost(userID, postID)
            return err
        }
        return err
    }

    // If already liked then unlike/remove like but if it's been unliked then change to like
    if like.IsLike {
        _, err := db.Exec("UPDATE post_likes SET is_like = ? WHERE post_id = ? AND user_id = ?", false, postID, userID)
        return err
    } else {
        _, err := db.Exec("UPDATE post_likes SET is_like = ? WHERE post_id = ? AND user_id = ?", true, postID, userID)
        RemoveDislikeFromPost(userID, postID)
        return err
    }
}


// PostLikeAndUnlike adds a like or dislike for a post or comment
func PostDisLikeAndUndislike(userID, postID int) error {
	// Check if the like already exists
	    var dislike PostDislike
    row := db.QueryRow("SELECT id, user_id, post_id, is_dislike FROM post_dislikes WHERE post_id = ? AND user_id = ?", postID, userID)

    err := row.Scan(&dislike.ID, &dislike.UserID, &dislike.PostID, &dislike.IsDislike)
    if err != nil {
        if err == sql.ErrNoRows {
            // Insert the like
            _, err = db.Exec("INSERT INTO post_dislikes (user_id, post_id, is_dislike) VALUES (?, ?, ?)", userID, postID, true)
            RemoveDislikeFromPost(userID, postID)
            return err
        }
        return err
    }

    // If already liked then unlike/remove like but if it's been unliked then change to like
    if dislike.IsDislike {
        _, err := db.Exec("UPDATE post_dislikes SET is_dislike = ? WHERE post_id = ? AND user_id = ?", false, postID, userID)
        return err
    } else {
        _, err := db.Exec("UPDATE post_dislikes SET is_dislike = ? WHERE post_id = ? AND user_id = ?", true, postID, userID)
        RemoveLikeFromPost(userID, postID)
        return err
    }
}


// RemoveDislikeFromPost removes dislike for a post
func RemoveDislikeFromPost(userID, postID int) error {

    var dislike PostDislike
    row := db.QueryRow("SELECT id, user_id, post_id, is_dislike FROM post_dislikes WHERE post_id = ? AND user_id = ?", postID, userID)

    err := row.Scan(&dislike.ID, &dislike.UserID, &dislike.PostID, &dislike.IsDislike)
    if err != nil {
        if err == sql.ErrNoRows {
            // No need to insert the dislike
            return err
        }
        return err
    }

    // If already disliked then remove dislike
    if dislike.IsDislike {
        _, err := db.Exec("UPDATE post_dislikes SET is_dislike = ? WHERE post_id = ? AND user_id = ?", false, postID, userID)
        return err
    }
    return err
}

// RemoveLikeFromPost removes Like for a post
func RemoveLikeFromPost(userID, postID int) error {

    var like PostLike
    row := db.QueryRow("SELECT id, user_id, post_id, is_like FROM post_likes WHERE post_id = ? AND user_id = ?", postID, userID)

    err := row.Scan(&like.ID, &like.UserID, &like.PostID, &like.IsLike)
    if err != nil {
        if err == sql.ErrNoRows {
            // No need to insert the like
            return err
        }
        return err
    }

    // If already liked then remove like
    if like.IsLike {
        _, err := db.Exec("UPDATE post_likes SET is_like = ? WHERE post_id = ? AND user_id = ?", false, postID, userID)
        return err
    }
    return err
}

// CountPostLikes returns the total number of likes for a post
func CountPostLikes(postID int) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM post_likes WHERE post_id = ? AND is_like = 1"
	err := db.QueryRow(query, postID).Scan(&count)
	return count, err
}

// CountPostDislikes returns the total number of likes for a post
func CountPostDislikes(postID int) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM post_dislikes WHERE post_id = ? AND is_dislike = 1"
	err := db.QueryRow(query, postID).Scan(&count)
	return count, err
}
 
func CommentLikeAndUnlike(userID, commentID int) error {
	// Check if the like already exists
	    var like CommentLike
     row := db.QueryRow("SELECT id, user_id, comment_id, is_like FROM comment_likes WHERE comment_id = ? AND user_id = ?", commentID, userID)

    err := row.Scan(&like.ID, &like.UserID, &like.CommentID, &like.IsLike)
    if err != nil {
        if err == sql.ErrNoRows {
            // Insert the like
            _, err = db.Exec("INSERT INTO comment_likes (user_id, comment_id, is_like) VALUES (?, ?, ?)", userID, commentID, true)
            RemoveDislikeFromComment(userID, commentID)
            return err
        }
        return err
    }

    // If already liked then unlike/remove like but if it's been unliked then change to like
    if like.IsLike {
        _, err := db.Exec("UPDATE comment_likes SET is_like = ? WHERE comment_id = ? AND user_id = ?", false, commentID, userID)
        return err
    } else {
        _, err := db.Exec("UPDATE comment_likes SET is_like = ? WHERE comment_id = ? AND user_id = ?", true, commentID, userID)
        RemoveDislikeFromComment(userID, commentID)
        return err
    }

}

func CommentDisLikeAndUndislike(userID, commentID int) error {
	// Check if the like already exists
	    var dislike CommentDislike
    row := db.QueryRow("SELECT id, user_id, comment_id, is_dislike FROM comment_dislikes WHERE comment_id = ? AND user_id = ?", commentID, userID)

    err := row.Scan(&dislike.ID, &dislike.UserID, &dislike.CommentID, &dislike.IsDislike)
    if err != nil {
        if err == sql.ErrNoRows {
            // Insert the like
            _, err = db.Exec("INSERT INTO comment_dislikes (user_id, comment_id, is_dislike) VALUES (?, ?, ?)", userID, commentID, true)
            RemoveDislikeFromComment(userID, commentID)
            return err
        }
        return err
    }

    // If already liked then unlike/remove like but if it's been unliked then change to like
    if dislike.IsDislike {
        _, err := db.Exec("UPDATE comment_dislikes SET is_dislike = ? WHERE comment_id = ? AND user_id = ?", false, commentID, userID)
        return err
    } else {
        _, err := db.Exec("UPDATE comment_dislikes SET is_dislike = ? WHERE comment_id = ? AND user_id = ?", true, commentID, userID)
        RemoveLikeFromComment(userID, commentID)
        return err
    }
}

// RemoveDislikeFromComment removes dislike for a post
func RemoveDislikeFromComment(userID, commentID int) error {

    var dislike CommentDislike
    row := db.QueryRow("SELECT id, user_id, comment_id, is_dislike FROM comment_dislikes WHERE comment_id = ? AND user_id = ?", commentID, userID)

    err := row.Scan(&dislike.ID, &dislike.UserID, &dislike.PostID, &dislike.IsDislike)
    if err != nil {
        if err == sql.ErrNoRows {
            // No need to insert the dislike
            return err
        }
        return err
    }

    // If already disliked then remove dislike
    if dislike.IsDislike {
        _, err := db.Exec("UPDATE comment_dislikes SET is_dislike = ? WHERE comment_id = ? AND user_id = ?", false, commentID, userID)
        return err
    }
    return err
}

// RemoveLikeFromComment removes Like for a post
func RemoveLikeFromComment(userID, commentID int) error {

    var like CommentLike
    row := db.QueryRow("SELECT id, user_id, comment_id, is_like FROM comment_likes WHERE comment_id = ? AND user_id = ?", commentID, userID)

    err := row.Scan(&like.ID, &like.UserID, &like.CommentID, &like.IsLike)
    if err != nil {
        if err == sql.ErrNoRows {
            // No need to insert the like
            return err
        }
        return err
    }

    // If already liked then remove like
    if like.IsLike {
        _, err := db.Exec("UPDATE comment_likes SET is_like = ? WHERE comment_id = ? AND user_id = ?", false, commentID, userID)
        return err
    }
    return err
}