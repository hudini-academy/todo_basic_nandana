package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")
var ErrInvalidCredentials = errors.New("models: invalid credentials")
var ErrDuplicateEmail = errors.New("models: duplicate email")

type Todo struct {
	ID       int
	Name     string
	Created  time.Time
	Modified time.Time
	Errors   map[string]string
}

type User struct {
	ID             int
	NAME           string
	email          string
	HashedPassword []byte
	Created        time.Time
	Errors         map[string]string
}
