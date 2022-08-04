package httpserver

import "errors"

var (
	ErrIncorrectStoreType   = errors.New("incorrect storage type")
	ErrIncorrectUrlPath     = errors.New("incorrect url path")
	ErrIncorrectRequestBody = errors.New("incorrect request body")
)
