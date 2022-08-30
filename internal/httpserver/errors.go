package httpserver

import (
	"errors"
)

var (
	ErrIncorrectStoreType   = errors.New("incorrect storage type")
	ErrIncorrectRequestBody = errors.New("incorrect request body")
)
