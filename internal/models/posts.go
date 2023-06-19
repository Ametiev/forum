package models

import (
	"database/sql"
	"errors"
	"time"
)

type Post struct {
	ID              int
	Title           string
	Content         string
	Created         time.Time
	Category        string
	UserName        string
	Comments        []Comment
	Likes           int
	Dislikes        int
	IsAuthenticated bool
}

type Comment struct {
	Id              int
	Author          string
	CContent        string
	PostID          int
	Likes           int
	Dislikes        int
	IsAuthenticated bool
}

type Model struct {
	DB *sql.DB
}

func (m *Model) Insert(title string, content string, category string, userName string) (int, error) {
	stmt := `INSERT INTO posts (title, content, created, category, user_name)
    VALUES (?, ?, datetime('now'), ?, ?)`
	result, err := m.DB.Exec(stmt, title, content, category, userName)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// func (m *Model) Get(id int) (*Post, error) {
// 	stmt := `SELECT id, title, content, created, category, user_name FROM posts
// 	WHERE id = ?`
// 	row := m.DB.QueryRow(stmt, id)
// 	s := &Post{}
// 	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Category, &s.UserName)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return nil, ErrNoRecord
// 		} else {
// 			return nil, err
// 		}
// 	}

// 	comments, err := m.GetComments(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	s.Comments = comments

// 	stmt = `SELECT like, dislike FROM post_reactions WHERE post_id = ?`
// 	rows, err := m.DB.Query(stmt, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var like, dislike int
// 		err = rows.Scan(&like, &dislike)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if like == 1 {
// 			s.Likes += like
// 		} else if dislike == 1 {
// 			s.Dislikes += dislike
// 		}
// 	}

// 	if err = rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return s, nil
// }

func (m *Model) Get(id int) (*Post, error) {
	stmt := `SELECT id, title, content, created, category, user_name FROM posts WHERE id = ?`
	row := m.DB.QueryRow(stmt, id)
	post := &Post{}
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.Created, &post.Category, &post.UserName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	stmt = `SELECT COALESCE(SUM(like), 0), COALESCE(SUM(dislike), 0) FROM post_reactions WHERE post_id = ?`
	row = m.DB.QueryRow(stmt, id)
	err = row.Scan(&post.Likes, &post.Dislikes)
	if err != nil {
		return nil, err
	}

	comments, err := m.GetComments(id)
	if err != nil {
		return nil, err
	}
	post.Comments = comments

	return post, nil
}

func (m *Model) Latest() ([]*Post, error) {
	stmt := `SELECT id, title, content, created, category, user_name FROM posts ORDER BY id DESC LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*Post{}
	for rows.Next() {
		s := &Post{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Category, &s.UserName)
		if err != nil {
			return nil, err
		}
		posts = append(posts, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (m *Model) ByCategory(category string) ([]*Post, error) {
	stmt := `SELECT id, title, content, created, category, user_name FROM posts WHERE category = ? ORDER BY id DESC LIMIT 10`
	rows, err := m.DB.Query(stmt, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*Post{}
	for rows.Next() {
		s := &Post{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Category, &s.UserName)
		if err != nil {
			return nil, err
		}
		posts = append(posts, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (m *Model) GetPostsByUser(name string) ([]*Post, error) {
	stmt := `SELECT id, title, content, created, category, user_name FROM posts WHERE user_name = ? ORDER BY id DESC`
	rows, err := m.DB.Query(stmt, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*Post{}
	for rows.Next() {
		p := &Post{}
		err = rows.Scan(&p.ID, &p.Title, &p.Content, &p.Created, &p.Category, &p.UserName)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (m *Model) GetComments(postID int) ([]Comment, error) {
	commentsQuery := `
		SELECT c.Id, c.CContent, c.Author, c.PostID, COALESCE(SUM(r.like), 0) AS Likes, COALESCE(SUM(r.dislike), 0) AS Dislikes
		FROM comments AS c
		LEFT JOIN comment_reactions AS r ON c.Id = r.comment_id
		WHERE c.PostID = ?
		GROUP BY c.Id, c.CContent, c.Author, c.PostID
	`

	stmt, err := m.DB.Prepare(commentsQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}

	for rows.Next() {
		comment := Comment{}
		err = rows.Scan(&comment.Id, &comment.CContent, &comment.Author, &comment.PostID, &comment.Likes, &comment.Dislikes)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (m *Model) PostComment(CommentInput Comment) error {
	if _, err := m.DB.Exec("INSERT INTO comments (CContent, Author, PostID) VALUES ($1,$2,$3)", CommentInput.CContent, CommentInput.Author, CommentInput.PostID); err != nil {
		return err
	}

	return nil
}

func (m *Model) GetPostsByUserReaction(userID int) ([]*Post, error) {
	stmt := `SELECT p.id, p.title, p.content, p.created, p.category, p.user_name
		FROM posts p
		INNER JOIN post_reactions pr ON p.id = pr.post_id
		WHERE pr.user_id = ? AND pr.like = 1
		ORDER BY p.id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*Post{}
	for rows.Next() {
		p := &Post{}
		err = rows.Scan(&p.ID, &p.Title, &p.Content, &p.Created, &p.Category, &p.UserName)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
