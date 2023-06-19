package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	UserID         int
	UserName       string
	Token          string
	ExpirationDate time.Time
}

func (m *Model) CreateSession(userId int, userName string) (string, time.Time, error) {
	token := uuid.NewString()
	date := time.Now().Add(20 * time.Minute)
	stmt := `INSERT INTO Sessions ( user_id, user_name, token, expiration_date) 
			VALUES(?,?,?,?)`
	_, err := m.DB.Exec(stmt, userId, userName, token, date)
	if err != nil {
		return "", date, err
	}
	return token, date, nil
}

func (m *Model) GetSessionFromToken(token string) (*Session, error) {
	stmt := `SELECT user_id, user_name, token, expiration_date FROM Sessions  where token = ?`

	row := m.DB.QueryRow(stmt, token)
	session := &Session{}
	err := row.Scan(&session.UserID, &session.UserName, &session.Token, &session.ExpirationDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return session, nil
}

func (m *Model) DeleteSessionByUserId(userId int) error {
	_, err := m.DB.Exec("DELETE FROM Sessions WHERE user_id = $1", userId)
	if err != nil {
		return err
	}
	return nil
}

func (m *Model) DeleteExpiredSessions() error {
	_, err := m.DB.Exec("DELETE FROM Sessions WHERE expiration_date < $1", time.Now())
	if err != nil {
		return err
	}
	return nil
}
