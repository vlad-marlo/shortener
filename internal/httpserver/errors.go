package httpserver

import (
	"errors"
)

var (
	// ErrIncorrectRequestBody ...
	ErrIncorrectRequestBody = errors.New("incorrect request body")
)
