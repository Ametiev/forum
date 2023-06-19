package models

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Duplicates(u User) error {
	var count int

	stmt := "SELECT COUNT(*) FROM users WHERE email = ? OR name = ?"

	err := m.DB.QueryRow(stmt, u.Email, u.Name).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return ErrDuplicateEntry
	}

	return nil
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	user := User{Email: email, Name: name}
	if err := m.Duplicates(user); err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
		VALUES(?, ?, ?, datetime('now'))`
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"
	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (u *UserModel) GetUserNameByEmail(email string) (string, error) {
	stmt := `SELECT name FROM Users WHERE email = ?`

	var name string
	err := u.DB.QueryRow(stmt, email).Scan(&name)
	if err != nil {
		return "", err
	}

	return name, nil
}
