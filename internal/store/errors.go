package store

import "errors"

// vars ...
var (
	// ErrNotFound ...
	ErrNotFound = errors.New("not found")
	// ErrAlreadyExists ...
	ErrAlreadyExists = errors.New("already in storage")
	// ErrNoContent ...
	ErrNoContent = errors.New("no content")
	// ErrIsDeleted ...
	ErrIsDeleted = errors.New("is deleted")
	// ErrNotAccessible ...
	ErrNotAccessible = errors.New("not accessible")
	// ErrAlreadyClosed ...
	ErrAlreadyClosed = errors.New("storage is already closed")
)
