package httpserver

import "errors"

var (
	IncorrectStoreType error = errors.New("Incorrect storage type.")
)
