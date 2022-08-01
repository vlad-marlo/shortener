package store

import "errors"

var (
	ErrNotFound      error = errors.New("not found")
	ErrAlreadyExists error = errors.New("already in storage")
)
