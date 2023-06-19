package models

import (
	"database/sql"
	"fmt"
)

type ReactionCommentModel struct {
	DB *sql.DB
}

func (r *ReactionCommentModel) LikeComment(user_id, comment_id int) error {
	// stmt := `SELECT like_comm, dislike_comm FROM reactionsComment WHERE comment_id = ? AND user_id = ?`

	fmt.Printf("userid=%d\n", user_id)
	stmt := `SELECT name FROM Users WHERE id =?`
	var name string
	fmt.Println(1)
	fmt.Println(2)
	err := r.DB.QueryRow(stmt, user_id).Scan(&name)
	fmt.Println(3)
	fmt.Println(4)
	fmt.Println(err)
	fmt.Println(name)

	// if err != nil {
	// 	if errors.Is(err, sql.ErrNoRows) {
	// 		_, err = r.DB.Exec(`INSERT INTO reactionsComment (comment_id, user_id, like_comm, dislike_comm) VALUES (?, ?, ?, ?)`, comment_id, user_id, 1, 0)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	} else {
	// 		return err
	// 	}
	// }
	// if like == 1 {
	// 	_, err := r.DB.Exec(`DELETE FROM reactions WHERE post_id = ? AND user_id = ?`, comment_id, user_id)
	// 	if err != nil {
	// 		return err
	// 	}
	// } else if dislike == 1 {
	// 	_, err := r.DB.Exec(`UPDATE reactions SET dislike = ? WHERE post_id = ? AND user_id = ?`, 0, comment_id, user_id)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	_, err = r.DB.Exec(`UPDATE reactions SET like = ? WHERE post_id = ? AND user_id = ?`, 1, comment_id, user_id)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return nil
}
