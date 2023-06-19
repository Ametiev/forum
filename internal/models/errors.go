package models

import (
	"errors"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEntry     = errors.New("duplicate entry")
	ErrNoComments         = errors.New("models: no comments found for post")
)
