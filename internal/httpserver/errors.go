package httpserver

import (
	"errors"
)

var (
	ErrIncorrectRequestBody = errors.New("incorrect request body")
)
