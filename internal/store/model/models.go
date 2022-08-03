package model

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrURLContainSpace = errors.New("url must have no spaces in it")
)

type URL struct {
	ID      string
	BaseURL string
}

func NewURL(url string) (URL, error) {
	if strings.Contains(url, string(" ")) {
		return URL{}, ErrURLContainSpace
	}
	return URL{
		ID:      uuid.New().String(),
		BaseURL: url,
	}, nil
}
