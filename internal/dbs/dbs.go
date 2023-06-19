package dbs

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "Forum.db")
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func CreatePosts(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created DATETIME NOT NULL,
			category TEXT NOT NULL,
			user_name TEXT NOT NULL
		);`,
		`CREATE INDEX IF NOT EXISTS idx_posts_created ON posts(created);`,
	}
	for _, stmt := range stmts {
		_, err := db.Exec(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateTables(b *sql.DB) error {
	var stmts []string = []string{Users, Comment, Session, PostReaction, CommentReaction}
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	for _, i := range stmts {

		_, err := b.Exec(i)
		if err != nil {
			errorLog.Fatal(err)
			return err
		}
	}
	return nil
}

const (
	Users = `CREATE TABLE IF NOT EXISTS Users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		hashed_password CHAR(60) NOT NULL,
		created DATETIME NOT NULL
	);`

	Comment = `CREATE TABLE IF NOT EXISTS comments (
		"Id"	INTEGER PRIMARY KEY AUTOINCREMENT,
		"CContent"	TEXT,
		"Author"	TEXT,
		"PostID" INTEGER
	);`
	Session = `CREATE TABLE IF NOT EXISTS Sessions (
		session_id INTEGER PRIMARY KEY,
		user_name TEXT,
		user_id INTEGER,
		token TEXT UNIQUE,
		expiration_date TIMESTAMP
	);`
	PostReaction = `CREATE TABLE IF NOT EXISTS post_reactions (
		user_id INTEGER,
		post_id INTEGER,
		like INTEGER,
		dislike INTEGER
	);`

	CommentReaction = `CREATE TABLE IF NOT EXISTS comment_reactions (
		user_id INTEGER,
		comment_id INTEGER,
		like INTEGER,
		dislike INTEGER
	);`
)
