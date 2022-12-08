package store

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already in storage")
	ErrNoContent     = errors.New("no content")
	ErrIsDeleted     = errors.New("is deleted")
	ErrNotAccessible = errors.New("not accessible")
	ErrAlreadyClosed = errors.New("storage is already closed")
)
