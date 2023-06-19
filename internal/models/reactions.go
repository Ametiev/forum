package models

import (
	"database/sql"
	"errors"
)

type ReactionModel struct {
	DB *sql.DB
}

func (r *ReactionModel) LikePost(userID, postID int) error {
	stmt := `SELECT like, dislike FROM post_reactions WHERE post_id = ? AND user_id = ?`
	var like, dislike int
	err := r.DB.QueryRow(stmt, postID, userID).Scan(&like, &dislike)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = r.DB.Exec(`INSERT INTO post_reactions (post_id, user_id, like, dislike) VALUES (?, ?, ?, ?)`, postID, userID, 1, 0)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if like == 1 {
		_, err := r.DB.Exec(`DELETE FROM post_reactions WHERE post_id = ? AND user_id = ?`, postID, userID)
		if err != nil {
			return err
		}
	} else if dislike == 1 {
		_, err := r.DB.Exec(`UPDATE post_reactions SET dislike = ? WHERE post_id = ? AND user_id = ?`, 0, postID, userID)
		if err != nil {
			return err
		}
		_, err = r.DB.Exec(`UPDATE post_reactions SET like = ? WHERE post_id = ? AND user_id = ?`, 1, postID, userID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ReactionModel) DislikePost(userID, postID int) error {
	stmt := `SELECT like, dislike FROM post_reactions WHERE post_id = ? AND user_id = ?`
	var like, dislike int
	err := r.DB.QueryRow(stmt, postID, userID).Scan(&like, &dislike)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = r.DB.Exec(`INSERT INTO post_reactions (post_id, user_id, like, dislike) VALUES (?, ?, ?, ?)`, postID, userID, 0, 1)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if dislike == 1 {
		_, err := r.DB.Exec(`DELETE FROM post_reactions WHERE post_id = ? AND user_id = ?`, postID, userID)
		if err != nil {
			return err
		}
	} else if like == 1 {
		_, err := r.DB.Exec(`UPDATE post_reactions SET like = ? WHERE post_id = ? AND user_id = ?`, 0, postID, userID)
		if err != nil {
			return err
		}
		_, err = r.DB.Exec(`UPDATE post_reactions SET dislike = ? WHERE post_id = ? AND user_id = ?`, 1, postID, userID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ReactionModel) LikeComment(userID, commentID int) error {
	stmt := `SELECT like, dislike FROM comment_reactions WHERE comment_id = ? AND user_id = ?`
	var like, dislike int
	err := r.DB.QueryRow(stmt, commentID, userID).Scan(&like, &dislike)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = r.DB.Exec(`INSERT INTO comment_reactions (comment_id, user_id, like, dislike) VALUES (?, ?, ?, ?)`, commentID, userID, 1, 0)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if like == 1 {
		_, err := r.DB.Exec(`DELETE FROM comment_reactions WHERE comment_id = ? AND user_id = ?`, commentID, userID)
		if err != nil {
			return err
		}
	} else if dislike == 1 {
		_, err := r.DB.Exec(`UPDATE comment_reactions SET dislike = ? WHERE comment_id = ? AND user_id = ?`, 0, commentID, userID)
		if err != nil {
			return err
		}
		_, err = r.DB.Exec(`UPDATE comment_reactions SET like = ? WHERE comment_id = ? AND user_id = ?`, 1, commentID, userID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ReactionModel) DislikeComment(userID, commentID int) error {
	stmt := `SELECT like, dislike FROM comment_reactions WHERE comment_id = ? AND user_id = ?`
	var like, dislike int
	err := r.DB.QueryRow(stmt, commentID, userID).Scan(&like, &dislike)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err = r.DB.Exec(`INSERT INTO comment_reactions (comment_id, user_id, like, dislike) VALUES (?, ?, ?, ?)`, commentID, userID, 0, 1)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if dislike == 1 {
		_, err := r.DB.Exec(`DELETE FROM comment_reactions WHERE comment_id = ? AND user_id = ?`, commentID, userID)
		if err != nil {
			return err
		}
	} else if like == 1 {
		_, err := r.DB.Exec(`UPDATE comment_reactions SET like = ? WHERE comment_id = ? AND user_id = ?`, 0, commentID, userID)
		if err != nil {
			return err
		}
		_, err = r.DB.Exec(`UPDATE comment_reactions SET dislike = ? WHERE comment_id = ? AND user_id = ?`, 1, commentID, userID)
		if err != nil {
			return err
		}
	}
	return nil
}
